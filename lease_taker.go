package klease

import (
	"math"
	"math/rand"
	"time"

	l4g "github.com/ezoic/log4go"
)

const takeRetries = 3
const scanRetries = 1

//Taker is used by LeaseCoordinator to take new leases, or leases that other workers fail to renew. Each
//LeaseCoordinator instance corresponds to one worker and uses exactly one ILeaseTaker to take leases for that worker.
type Taker struct {
	leaseManager           *Manager
	workerId               string
	allLeases              map[string]*KLease
	leaseDurationNanos     int64
	maxLeasesForWorker     int64
	maxLeasesToStealAtOnce int64
	lastScanTimeNanos      int64
}

//NewKLeaseTaker initializes and returns a Taker
func NewKLeaseTaker(manager *Manager, workerId string, leaseDurationNanos int64) *Taker {
	rand.Seed(time.Now().UnixNano())
	return &Taker{
		leaseManager:           manager,
		workerId:               workerId,
		leaseDurationNanos:     leaseDurationNanos,
		maxLeasesForWorker:     math.MaxInt64,
		maxLeasesToStealAtOnce: 1,
		allLeases:              map[string]*KLease{},
	}
}

//WithMaxLeasesForWorker will allow us to set a maximum. This is usually not desired as it can lead to
//data loss if there are not enough workers. Some shards may never be processed. Useful for testing or
//if absolutely necessary due to resource constraints
func (t *Taker) WithMaxLeasesForWorker(maxLeasesForWorker int64) *Taker {
	//must be able to get at least 1 lease
	if maxLeasesForWorker <= 0 {
		maxLeasesForWorker = 1
	}
	t.maxLeasesForWorker = maxLeasesForWorker
	return t
}

//WithMaxLeasesToStealAtOneTime sets a max to take from a more loaded Worker at one time (for load balancing).
//Setting this to a higher number can allow for faster load convergence (e.g. during deployments, cold starts),
//but can cause higher churn in the system.
func (t *Taker) WithMaxLeasesToStealAtOneTime(maxLeasesToStealAtOnce int64) *Taker {
	//must be able to steal at least 1 lease at a time
	if maxLeasesToStealAtOnce <= 0 {
		maxLeasesToStealAtOnce = 1
	}
	t.maxLeasesToStealAtOnce = maxLeasesToStealAtOnce
	return t
}

//TakeLeases compute the set of leases available to be taken and attempt to take them. Lease taking rules are:
//1) If a lease's counter hasn't changed in long enough, try to take it.
//2) If we see a lease we've never seen before, take it only if owner == null. If it's owned, odds are the owner is
//holding it. We can't tell until we see it more than once.
//3) For load balancing purposes, you may violate rules 1 and 2 for EXACTLY ONE lease per call of takeLeases().
func (t *Taker) TakeLeases() map[string]*KLease {
	takenLeases := map[string]*KLease{}

	var lastErr error

	//dynamodb.ErrCodeProvisionedThroughputExceededException

	for i := 1; i <= scanRetries; i++ {
		err := t.updateAllLeases()
		if err != nil {
			lastErr = err
		}
	}

	if lastErr != nil {
		return takenLeases
	}

	expiredLeases := t.getExpiredLeases()
	//l4g.Debug("Worker %s - Expired Leases %d", t.workerId, len(expiredLeases))
	leasesToTake := t.computeLeasesToTake(expiredLeases)
	untakenLeaseKeys := []string{}

	for _, lease := range leasesToTake {
		leaseKey := lease.GetLeaseKey()

		for i := 1; i <= takeRetries; i++ {
			isTaken, err := t.leaseManager.TakeLease(lease, t.workerId)
			if err != nil {
				l4g.Error("Worker %s - Error taking leases: %s", t.workerId, err.Error())
			} else {
				//l4g.Debug("Worker %s - Did we succeed in taking lease %s? : %t", t.workerId, leaseKey, isTaken)
				if isTaken {
					l4g.Info("Worker %s - Stole lease in table %s, key %s", t.workerId, t.leaseManager.table, leaseKey)
					lease.SetLastCounterIncrementNanos(time.Now().UnixNano())
					takenLeases[leaseKey] = lease
				} else {
					untakenLeaseKeys = append(untakenLeaseKeys, leaseKey)
				}
				//didn't error so don't need to try again
				break
			}
		}

	}

	return takenLeases
}

// updateAllLeases scans all leases and updates lastRenewalTime. Also adds new leases and deletes old leases.
func (t *Taker) updateAllLeases() error {
	freshList, err := t.leaseManager.ListLeases()
	if err != nil {
		return err
	}

	t.lastScanTimeNanos = time.Now().UnixNano()

	// This set will hold the lease keys not updated by the previous listLeases call.
	notUpdated := map[string]bool{}
	for k := range t.allLeases {
		notUpdated[k] = true
	}

	// Iterate over all leases, finding ones to try to acquire that haven't changed since the last iteration
	for _, lease := range freshList {
		leaseKey := lease.GetLeaseKey()
		var oldLease *KLease
		if ol, ok := t.allLeases[leaseKey]; ok {
			oldLease = ol
		} else {
			oldLease = nil
		}
		t.allLeases[leaseKey] = lease
		if _, ok := notUpdated[leaseKey]; ok {
			delete(notUpdated, leaseKey)
		}

		if oldLease != nil {
			// If we've seen this lease before...
			if oldLease.GetLeaseCounter() == lease.GetLeaseCounter() {
				// ...and the counter hasn't changed, propagate the lastRenewalNanos time from the old lease
				lease.SetLastCounterIncrementNanos(oldLease.GetLastCounterIncrementNanos())
			} else {
				// ...and the counter has changed, set lastRenewalNanos to the time of the scan.
				lease.SetLastCounterIncrementNanos(t.lastScanTimeNanos)
			}
		} else {
			if lease.GetLeaseOwner() == "" {
				// if this new lease is unowned, it's never been renewed.
				lease.SetLastCounterIncrementNanos(0)
			} else {
				// if this new lease is owned, treat it as renewed as of the scan
				lease.SetLastCounterIncrementNanos(t.lastScanTimeNanos)
			}
		}

	}

	// Remove dead leases from allLeases
	for key := range notUpdated {
		if _, ok := t.allLeases[key]; ok {
			delete(t.allLeases, key)
		}
	}
	return nil

}

func (t *Taker) getExpiredLeases() []*KLease {
	expiredLeases := []*KLease{}
	for _, lease := range t.allLeases {
		if lease.IsExpired(t.leaseDurationNanos, t.lastScanTimeNanos) {
			//l4g.Debug("Worker %s - Found Expired Lease %s", t.workerId, lease.GetLeaseKey())
			expiredLeases = append(expiredLeases, lease)
		}
	}
	return expiredLeases
}

//Compute the number of leases I should try to take based on the state of the system.
func (t *Taker) computeLeasesToTake(expiredLeases []*KLease) []*KLease {
	leaseCounts := t.computeLeaseCounts(expiredLeases)
	leasesToTake := []*KLease{}
	var target, numLeases, numWorkers int64
	numLeases = int64(len(t.allLeases))
	numWorkers = int64(len(leaseCounts))

	//l4g.Debug("Worker %s - NumWorkers: %d", t.workerId, numWorkers)

	//no leases, no take-y
	if numLeases <= 0 {
		return leasesToTake
	}

	if numWorkers >= numLeases {
		target = 1
	} else {
		var addOn int64
		if numLeases%numWorkers > 0 {
			addOn = 1
		}
		target = numLeases/numWorkers + addOn

		// Spill over is the number of leases this worker should have claimed, but did not because it would
		// exceed the max allowed for this worker.
		var leaseSpillover int64
		if target > t.maxLeasesForWorker {
			leaseSpillover = target - t.maxLeasesForWorker
			target = t.maxLeasesForWorker
			if leaseSpillover > 0 {
				//l4g.Debug("Worker %s - Spillover: %d", t.workerId, leaseSpillover)
			}
		}
	}

	myCount := leaseCounts[t.workerId]
	numLeasesToReachTarget := target - myCount

	//l4g.Debug("Worker %s - Target: %d", t.workerId, target)
	//l4g.Debug("Worker %s - numLeasesToReachTarget: %d", t.workerId, numLeasesToReachTarget)
	if numLeasesToReachTarget <= 0 {
		// If we don't need anything, return empty.
		return leasesToTake
	}

	// Shuffle expiredLeases so workers don't all try to contend for the same leases.
	t.shuffle(expiredLeases)

	originalExpiredLeasesSize := len(expiredLeases)

	if originalExpiredLeasesSize > 0 {
		//l4g.Debug("Worker %s - Going to try to take free leases", t.workerId)
		for numLeasesToReachTarget > 0 && len(expiredLeases) > 0 {
			//fancy pants pop from slice
			var eLease *KLease
			eLease, expiredLeases = expiredLeases[len(expiredLeases)-1], expiredLeases[:len(expiredLeases)-1]
			leasesToTake = append(leasesToTake, eLease)
			numLeasesToReachTarget--
		}
	} else {
		// If there are no expired leases and we need a lease, consider stealing.
		//l4g.Debug("Worker %s - Going to try to steal leases", t.workerId)
		leasesToSteal := t.chooseLeasesToSteal(leaseCounts, numLeasesToReachTarget, target)
		for _, leaseToSteal := range leasesToSteal {
			leasesToTake = append(leasesToTake, leaseToSteal)
		}
	}
	//l4g.Debug("Worker %s - Leases I'm going to try to take: %d", t.workerId, len(leasesToTake))
	return leasesToTake
}

//Choose leases to steal by randomly selecting one or more (up to max) from the most loaded worker.
//Stealing rules:
//
//Steal up to maxLeasesToStealAtOneTime leases from the most loaded worker if
//a) he has > target leases and I need >= 1 leases : steal min(leases needed, maxLeasesToStealAtOneTime)
//b) he has == target leases and I need > 1 leases : steal 1
func (t *Taker) chooseLeasesToSteal(leaseCounts map[string]int64, needed, target int64) []*KLease {
	var mostLoadedWorker string
	var mostLoadedWorkerLoad int64

	//find most loaded worker
	for worker, load := range leaseCounts {
		if mostLoadedWorker == "" || mostLoadedWorkerLoad < load {
			mostLoadedWorker = worker
			mostLoadedWorkerLoad = load
		}
	}

	var numLeasesToSteal int64
	if mostLoadedWorkerLoad >= target && needed > 0 {
		leasesOverTarget := mostLoadedWorkerLoad - target
		if needed > leasesOverTarget {
			numLeasesToSteal = leasesOverTarget
		} else {
			numLeasesToSteal = needed
		}

		// steal 1 if we need > 1 and max loaded worker has target leases.
		if needed > 1 && numLeasesToSteal == 0 {
			numLeasesToSteal = 1
		}

		if numLeasesToSteal > t.maxLeasesToStealAtOnce {
			numLeasesToSteal = t.maxLeasesToStealAtOnce
		}
	}

	if numLeasesToSteal <= 0 {
		return []*KLease{}
	}

	//get leases belonging to mostLoadedWorker
	candidates := []*KLease{}
	for _, lease := range t.allLeases {
		if lease.GetLeaseOwner() == mostLoadedWorker {
			candidates = append(candidates, lease)
		}
	}

	t.shuffle(candidates)
	if int64(len(candidates)) < numLeasesToSteal {
		numLeasesToSteal = int64(len(candidates))
	}

	return candidates[:numLeasesToSteal]
}

//Count leases by host. Always includes myself, but otherwise only includes hosts that are currently holding leases.
func (t *Taker) computeLeaseCounts(expiredLeases []*KLease) map[string]int64 {
	leaseCounts := map[string]int64{}
	// Compute the number of leases per worker by looking through allLeases and ignoring leases that have expired or have no owner.
	for _, lease := range t.allLeases {
		//l4g.Debug("Worker %s - Checking Lease %s and its expired state is %t", t.workerId, lease.GetLeaseKey(), t.contains(lease, expiredLeases))
		if !t.contains(lease, expiredLeases) && lease.GetLeaseOwner() != "" {
			leaseOwner := lease.GetLeaseOwner()
			leaseCounts[leaseOwner]++
		}
	}

	//If I have no leases, I wasn't represented in leaseCounts. Let's fix that
	//if map value does not exist go will return a falsy value by default,
	//but will not actually exists for things like len() unless we actually set it
	if leaseCounts[t.workerId] == 0 {
		leaseCounts[t.workerId] = 0
	}

	return leaseCounts
}

//GetWorkerId returns the worker id for the current taker
func (t *Taker) GetWorkerId() string {
	return t.workerId
}

func (t *Taker) contains(a *KLease, b []*KLease) bool {
	for _, lease := range b {
		if a.Equals(lease) {
			return true
		}
	}
	return false
}

func (t *Taker) shuffle(a []*KLease) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
