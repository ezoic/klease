package klease

import "github.com/nu7hatch/gouuid"
import "errors"
import "sync"
import "time"
import "fmt"

type Renewer struct {
	leaseManager       *Manager
	ownedLeases        map[string]*KLease
	workerId           string
	leaseDurationNanos int64
	ownedLeasesMutex   map[string]*sync.Mutex
}

type renewResult struct {
	success bool
	err     error
}

const RenewalRetries = 2
const MaxWorkersForRenewing = 20

func NewKLeaseRenewer(leaseManager *Manager, workerId string, leaseDurationNanos int64) *Renewer {
	return &Renewer{
		leaseManager:       leaseManager,
		workerId:           workerId,
		leaseDurationNanos: leaseDurationNanos,
		ownedLeasesMutex:   map[string]*sync.Mutex{},
	}
}

//todo
func (r *Renewer) RenewLeases() error {
	var lostLeases, leasesInUnknownState int64
	var lastErr error
	renewLeaseTasks := make(chan renewResult, len(r.ownedLeases))
	sem := make(chan bool, MaxWorkersForRenewing)
	for _, lease := range r.ownedLeases {
		sem <- true
		go r.renewLease(lease, renewLeaseTasks, sem)
	}

	for i := 0; i < len(r.ownedLeases); i++ {
		result := <-renewLeaseTasks
		if result.err != nil {
			leasesInUnknownState++
			lastErr = result.err
		} else if !result.success {
			lostLeases++
		}

	}
	if leasesInUnknownState > 0 {
		return fmt.Errorf("Encountered an exception while renewing leases. The number of leases which "+
			"might not have been renewed is %d. Last exception was: %s", leasesInUnknownState, lastErr)
	}
	return nil
}

func (r *Renewer) renewLease(lease *KLease, renewLeaseTasks chan renewResult, sem chan bool) {
	defer func() { <-sem }()
	renewed, err := r.renewLeaseInner(lease, false)
	renewLeaseTasks <- renewResult{
		success: renewed,
		err:     err,
	}
}

func (r *Renewer) renewLeaseInner(lease *KLease, renewEvenIfExpired bool) (bool, error) {
	var err error
	leaseKey := lease.GetLeaseKey()
	renewedLease := false

	var locker *sync.Mutex
	if _, ok := r.ownedLeasesMutex[leaseKey]; ok {
		locker = r.ownedLeasesMutex[leaseKey]
	} else {
		//this won't do much, but keeps us from having to check if it exists later
		locker = &sync.Mutex{}
	}

	for i := 1; i <= RenewalRetries; i++ {
		locker.Lock()
		// Don't renew expired lease during regular renewals. getCopyOfHeldLease may have returned null
		// triggering the application processing to treat this as a lost lease.
		if renewEvenIfExpired || !lease.IsExpired(r.leaseDurationNanos, time.Now().UnixNano()) {
			renewedLease, err = r.leaseManager.RenewLease(lease)
			if err != nil {
				if err.Error() == "ProvisionedThroughputExceededException" {
					locker.Unlock()
					continue
				}
				locker.Unlock()
				return false, err
			}
		}
		if renewedLease {
			err = lease.SetLastCounterIncrementNanos(time.Now().UnixNano())
			if err != nil {
				locker.Unlock()
				return false, err
			}
		} else {
			delete(r.ownedLeases, leaseKey)
			delete(r.ownedLeasesMutex, leaseKey)
		}
		locker.Unlock()
		break
	}

	return renewedLease, nil
}

//GetCurrentlyHeldLeases returns a map of  lease key to leases we hold
func (r *Renewer) GetCurrentlyHeldLeases() map[string]*KLease {
	result := map[string]*KLease{}
	now := time.Now().UnixNano()

	for leaseKey := range r.ownedLeases {
		copyLease := r.getCopyOfHeldLease(leaseKey, now)
		if copyLease != nil {
			result[copyLease.GetLeaseKey()] = copyLease
		}
	}

	return result
}

//GetCurrentlyHeldLease returns a copy of a lease if we hold it or nil if we don't
func (r *Renewer) GetCurrentlyHeldLease(leaseKey string) *KLease {
	return r.getCopyOfHeldLease(leaseKey, time.Now().UnixNano())
}

//getCopyOfHeldLease is an internal method to return a lease with a specific lease key only if we currently hold it.
func (r *Renewer) getCopyOfHeldLease(leaseKey string, now int64) *KLease {
	if _, ok := r.ownedLeases[leaseKey]; !ok {
		return nil
	}

	authoritativeLease := r.ownedLeases[leaseKey]
	r.ownedLeasesMutex[leaseKey].Lock()
	leaseCopy := *authoritativeLease
	r.ownedLeasesMutex[leaseKey].Unlock()

	if leaseCopy.IsExpired(r.leaseDurationNanos, now) {
		return nil
	}

	return &leaseCopy
}

func (r *Renewer) UpdateLease(lease *KLease, concurrencyToken *uuid.UUID) (bool, error) {
	if lease == nil || lease.GetLeaseKey() == "" || concurrencyToken == nil {
		return false, errors.New("lease, leasekey, and concurrencyToken cannot be nil / empty")
	}

	leaseKey := lease.GetLeaseKey()
	var authoritativeLease *KLease
	if al, ok := r.ownedLeases[leaseKey]; ok {
		authoritativeLease = al
	} else {
		return false, nil
	}

	//If the passed-in concurrency token doesn't match the concurrency token of the authoritative lease, it means
	//the lease was lost and regained between when the caller acquired his concurrency token and when the caller
	//called update.
	if authoritativeLease.GetConcurrencyToken().String() != concurrencyToken.String() {
		return false, nil
	}

	r.ownedLeasesMutex[leaseKey].Lock()
	defer r.ownedLeasesMutex[leaseKey].Unlock()

	authoritativeLease.Update(lease)
	updatedLease, err := r.leaseManager.UpdateLease(authoritativeLease)
	if err != nil {
		return false, err
	}

	if updatedLease {
		// Updates increment the counter
		err := authoritativeLease.SetLastCounterIncrementNanos(time.Now().UnixNano())
		if err != nil {
			return false, err
		}
	} else {
		/*
		* If updateLease returns false, it means someone took the lease from us. Remove the lease
		* from our set of owned leases pro-actively rather than waiting for a run of renewLeases().
		 */

		/*
		* Remove only if the value currently in the map is the same as the authoritative lease. We're
		* guarding against a pause after the concurrency token check above. It plays out like so:
		*
		* 1) Concurrency token check passes
		* 2) Pause. Lose lease, re-acquire lease. This requires at least one lease counter update.
		* 3) Unpause. leaseManager.updateLease fails conditional write due to counter updates, returns
		* false.
		* 4) ownedLeases.remove(key, value) doesn't do anything because authoritativeLease does not
		* .equals() the re-acquired version in the map on the basis of lease counter. This is what we want.
		* If we just used ownedLease.remove(key), we would have pro-actively removed a lease incorrectly.
		*
		* Note that there is a subtlety here - Lease.equals() deliberately does not check the concurrency
		* token, but it does check the lease counter, so this scheme works.
		 */
		if r.ownedLeases[leaseKey].Equals(authoritativeLease) {
			delete(r.ownedLeases, leaseKey)
			delete(r.ownedLeasesMutex, leaseKey)
		}
	}
	return updatedLease, nil
}

func (r *Renewer) AddLeasesToRenew(newLeases []*KLease) error {

	for _, lease := range newLeases {
		if lease.GetLastCounterIncrementNanos() == 0 {
			continue
		}

		authoritativeLease := *lease
		token, err := uuid.NewV4()
		if err != nil {
			return err
		}
		authoritativeLease.SetConcurrencyToken(token)
		r.ownedLeases[authoritativeLease.GetLeaseKey()] = &authoritativeLease
		r.ownedLeasesMutex[authoritativeLease.GetLeaseKey()] = &sync.Mutex{}
	}

	return nil
}

func (r *Renewer) ClearCurentHeldLeases() {
	r.ownedLeases = map[string]*KLease{}
	r.ownedLeasesMutex = map[string]*sync.Mutex{}
}

func (r *Renewer) DropLease(lease *KLease) {
	if _, ok := r.ownedLeases[lease.GetLeaseKey()]; ok {
		delete(r.ownedLeases, lease.GetLeaseKey())
		delete(r.ownedLeasesMutex, lease.GetLeaseKey())
	}
}

func (r *Renewer) Init() error {
	leases, err := r.leaseManager.ListLeases()
	if err != nil {
		return err
	}
	myLeases := []*KLease{}
	renewEvenIfExpired := true

	for _, lease := range leases {
		if r.workerId == lease.GetLeaseOwner() {
			// Okay to renew even if lease is expired, because we start with an empty list and we add the lease to
			// our list only after a successful renew. So we don't need to worry about the edge case where we could
			// continue renewing a lease after signaling a lease loss to the application.
			renewed, err := r.renewLeaseInner(lease, renewEvenIfExpired)
			if err != nil {
				return err
			}
			if renewed {
				myLeases = append(myLeases, lease)
			}
		}
	}
	return r.AddLeasesToRenew(myLeases)
}
