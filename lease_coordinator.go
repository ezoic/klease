package klease

import (
	"math"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

/**
 * Coordinator abstracts away LeaseTaker and LeaseRenewer from the application code that's using leasing. It owns
 * the scheduling of the two previously mentioned components as well as informing LeaseRenewer when LeaseTaker takes new
 * leases.
 *
 */
type Coordinator struct {
	renewer               *Renewer
	taker                 *Taker
	renewerIntervalMillis int64
	takerIntervalMillis   int64
	shutdownLock          *sync.Mutex
	running               bool
	isTakerCancelled      chan bool
	isRenewerCancelled    chan bool
	cancelTaker           bool
	cancelRenewer         bool
}

const WorkerIdentifierMetric = "WorkerIdentifier"
const StopWaitTimeMillis = 2000
const DefaultMaxLeasesForWorker = math.MaxInt64
const DefaultMaxLeasesToStealAtOneTime = 1

//NewLeaseCoordinator generates and returns a new Coordinator.
//epsilonNanos Allows for some variance when calculating lease expirations
func NewLeaseCoordinator(leaseManager *Manager, workerId string, leaseDurationMillis, epsilonMillis int64) *Coordinator {
	return &Coordinator{
		renewer: NewKLeaseRenewer(leaseManager, workerId, leaseDurationMillis),
		taker:   NewKLeaseTaker(leaseManager, workerId, leaseDurationMillis).WithMaxLeasesForWorker(DefaultMaxLeasesForWorker).WithMaxLeasesToStealAtOneTime(DefaultMaxLeasesToStealAtOneTime),
		renewerIntervalMillis: leaseDurationMillis/3 - epsilonMillis,
		takerIntervalMillis:   (leaseDurationMillis + epsilonMillis) * 2,
		shutdownLock:          &sync.Mutex{},
	}
}

//RunTaker runs a single iteration of the lease taker - used by integration tests.
func (c *Coordinator) RunTaker() {

	takenLeases := c.taker.TakeLeases()

	takenLeasesSlice := make([]*KLease, len(takenLeases))
	idx := 0
	for _, value := range takenLeases {
		takenLeasesSlice[idx] = value
		idx++
	}

	c.shutdownLock.Lock()
	if c.running {
		err := c.renewer.AddLeasesToRenew(takenLeasesSlice)
		if err != nil {
			//log here
		}
	}
	c.shutdownLock.Unlock()
}

//RunRenewer runs a single iteration of the lease renewer - used by integration tests.
func (c *Coordinator) RunRenewer() error {
	err := c.renewer.RenewLeases()
	if err != nil {
		//log here
	}
	return err
}

//Start background LeaseHolder and LeaseTaker threads.
func (c *Coordinator) Start() error {
	err := c.renewer.Init()
	if err != nil {
		return err
	}

	c.cancelRenewer = false
	c.cancelTaker = false

	// Taker runs with fixed DELAY because we want it to run slower in the event of performance degredation.
	go c.loopTakerWithFixedDelay(c.takerIntervalMillis)

	// Renewer runs at fixed INTERVAL because we want it to run at the same rate in the event of degredation.
	go c.loopRenewerWithFixedInterval(c.renewerIntervalMillis)

	c.running = true

	return nil
}

//this one will wait delayInMillis since the last one finished before firing off the next
func (c *Coordinator) loopTakerWithFixedDelay(delayInMillis int64) {
	c.isTakerCancelled = make(chan bool, 1)
	for !c.cancelTaker {
		c.RunTaker()
		if !c.cancelTaker {
			time.Sleep(time.Duration(delayInMillis) * time.Millisecond)
		}

	}
	c.isTakerCancelled <- true
}

//this one will wait intervalInMillis since the last one started before firing off the next
//unless the last one is not done. Then it will wait until it is
func (c *Coordinator) loopRenewerWithFixedInterval(intervalInMillis int64) {
	c.isRenewerCancelled = make(chan bool, 1)
	for !c.cancelRenewer {
		startTime := time.Now().UnixNano() / 1000000
		c.RunRenewer()
		elapsedTime := startTime - (time.Now().UnixNano() / 1000000)
		if elapsedTime < intervalInMillis && !c.cancelRenewer {
			time.Sleep(time.Duration(intervalInMillis-elapsedTime) * time.Millisecond)
		}
	}
	c.isRenewerCancelled <- true
}

//GetAssignments returns currently held leases
func (c *Coordinator) GetAssignments() []*KLease {
	leases := c.renewer.GetCurrentlyHeldLeases()
	leasesSlice := make([]*KLease, len(leases))

	idx := 0
	for _, value := range leases {
		leasesSlice[idx] = value
		idx++
	}
	return leasesSlice
}

// GetCurrentlyHeldLease returns deep copy of currently held Lease for given key, or null if we don't hold the lease for that key
func (c *Coordinator) GetCurrentlyHeldLease(leaseKey string) *KLease {
	return c.renewer.GetCurrentlyHeldLease(leaseKey)
}

//GetWorkerId returns the worker id as a string
func (c *Coordinator) GetWorkerId() string {
	return c.taker.GetWorkerId()
}

//Stop stops background threads and waits for all background tasks to complete.
//it should force stop it after a certain amount of time but not sure how to implement that atm
//i'll leave that part as a todo
func (c *Coordinator) Stop() {

	c.cancelRenewer = true
	c.cancelTaker = true
	//check if they are running and wait for them to stop
	if c.isRenewerCancelled != nil && c.isTakerCancelled != nil {
		<-c.isRenewerCancelled
		<-c.isTakerCancelled
	} else if c.isRenewerCancelled != nil {
		<-c.isRenewerCancelled
	} else if c.isTakerCancelled != nil {
		<-c.isTakerCancelled
	}
	c.isRenewerCancelled = nil
	c.isTakerCancelled = nil

	c.shutdownLock.Lock()
	c.renewer.ClearCurentHeldLeases()
	c.running = false
	c.shutdownLock.Unlock()
}

// StopLeaseTaker requests the cancellation of the lease taker.
func (c *Coordinator) StopLeaseTaker() {
	c.cancelTaker = true
}

// DropLease requests that renewals for the given lease are stopped.
func (c *Coordinator) DropLease(lease *KLease) {
	if lease != nil {
		c.shutdownLock.Lock()
		c.renewer.DropLease(lease)
		c.shutdownLock.Unlock()
	}
}

//IsRunning returns true if this LeaseCoordinator is running
func (c *Coordinator) IsRunning() bool {
	return c.running
}

//UpdateLease updates application specific lease values in DynamoDB
func (c *Coordinator) UpdateLease(lease *KLease, concurrencyToken *uuid.UUID) (bool, error) {
	return c.renewer.UpdateLease(lease, concurrencyToken)
}
