package klease

import (
	"crypto/md5"
	"errors"
	"io"
	"strconv"

	"github.com/nu7hatch/gouuid"
)

//KCheckpoint implements the Checkpont interface
type KCheckpoint struct {
	AppName        string
	StreamName     string
	SequenceNumber string
}

// KLease contains data pertaining to a Lease. Distributed systems may use leases to partition work across a
// fleet of workers. Each shard (identified by a leaseKey) has a corresponding Lease. Every worker will contend
// for all leases - only one worker will successfully take each one. The worker should hold the lease until it is ready to stop
// processing the corresponding unit of work, or until it fails. When the worker stops holding the lease, another worker will
// take and hold the lease.
type KLease struct {
	checkpoint                   *KCheckpoint
	ownerSwitchesSinceCheckpoint int64
	//parentShardIds               map[string]string //might not be used

	leaseKey     string //this is shard id
	leaseOwner   string
	leaseCounter int64

	concurrencyToken          *uuid.UUID
	lastCounterIncrementNanos int64

	status string // contains status information about the lease, useful for storing "stolen"
}

// NewKLeaseFromLease creates and returns a copy of other
func NewKLeaseFromLease(other KLease) *KLease {
	lease := KLease{
		checkpoint:                   other.checkpoint,
		ownerSwitchesSinceCheckpoint: other.ownerSwitchesSinceCheckpoint,
		//parentShardIds:               other.parentShardIds,
		leaseKey:                  other.leaseKey,
		leaseOwner:                other.leaseOwner,
		leaseCounter:              other.leaseCounter,
		concurrencyToken:          other.concurrencyToken,
		lastCounterIncrementNanos: other.lastCounterIncrementNanos,
	}
	return &lease
}

func NewKLease(leaseKey, leaseOwner string, checkpoint *KCheckpoint) *KLease {
	lease := KLease{
		checkpoint: checkpoint,
		leaseKey:   leaseKey,
		leaseOwner: leaseOwner,
	}
	return &lease
}

func (l *KLease) GetLeaseKey() string {
	return l.leaseKey
}

func (l *KLease) GetLeaseCounter() int64 {
	return l.leaseCounter
}

func (l *KLease) GetLeaseOwner() string {
	return l.leaseOwner
}

func (l *KLease) GetConcurrencyToken() *uuid.UUID {
	return l.concurrencyToken
}

func (l *KLease) GetLastCounterIncrementNanos() int64 {
	return l.lastCounterIncrementNanos
}

func (l *KLease) GetCheckpoint() *KCheckpoint {
	return l.checkpoint
}

func (l *KLease) GetOwnerSwitchesSinceCheckpoint() int64 {
	return l.ownerSwitchesSinceCheckpoint
}

// func (l *KLease) GetParentShardIds() map[string]string {
// 	return l.parentShardIds
// }

func (l *KLease) IsExpired(leaseDurationNanos, asOfNanos int64) bool {
	if l.lastCounterIncrementNanos < 0 {
		return true
	}

	age := asOfNanos - l.lastCounterIncrementNanos
	return age > leaseDurationNanos

}

func (l *KLease) SetLeaseKey(leaseKey string) error {
	if l.leaseKey != "" {
		return errors.New("Cannot change leaseKey once set")
	}
	l.leaseKey = leaseKey
	return nil
}

func (l *KLease) SetLeaseCounter(leaseCounter int64) error {
	l.leaseCounter = leaseCounter
	return nil
}

func (l *KLease) SetLeaseOwner(leaseOwner string) error {
	l.leaseOwner = leaseOwner
	return nil
}

func (l *KLease) SetConcurrencyToken(concurrencyToken *uuid.UUID) error {
	if concurrencyToken == nil {
		return errors.New("concurrencyToken Cannot be nil")
	}
	l.concurrencyToken = concurrencyToken
	return nil
}

func (l *KLease) SetLastCounterIncrementNanos(lastCounterIncrementNanos int64) error {
	l.lastCounterIncrementNanos = lastCounterIncrementNanos
	return nil
}

func (l *KLease) SetCheckpoint(checkpoint *KCheckpoint) error {
	l.checkpoint = checkpoint
	return nil
}

func (l *KLease) SetOwnerSwitchesSinceCheckpoint(ownerSwitchesSinceCheckpoint int64) error {
	l.ownerSwitchesSinceCheckpoint = ownerSwitchesSinceCheckpoint
	return nil
}

// func (l *KLease) SetParentShardIds(parentShardIds map[string]string) error {
// 	l.parentShardIds = parentShardIds
// 	return nil
// }

//Update updates this Lease's mutable, application-specific fields based on the passed-in lease object. Does not update
//fields that are internal to the leasing library (leaseKey, leaseOwner, leaseCounter)
func (l *KLease) Update(other *KLease) {
	l.SetOwnerSwitchesSinceCheckpoint(other.ownerSwitchesSinceCheckpoint)
	l.SetCheckpoint(other.checkpoint)
	//l.SetParentShardIds(other.parentShardIds)
}

func (l *KLease) HashCode() string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(l.leaseCounter, 10)+l.leaseOwner+l.leaseKey)

	return string(h.Sum(nil)[:h.Size()])
}

//Equals returns true if lease other equals this lease
func (l *KLease) Equals(other *KLease) bool {
	if other == nil ||
		l.leaseCounter != other.leaseCounter ||
		l.leaseOwner != other.leaseOwner ||
		l.leaseKey != other.leaseKey ||
		l.checkpoint != other.checkpoint ||
		l.ownerSwitchesSinceCheckpoint != other.ownerSwitchesSinceCheckpoint {
		return false
	}

	return true
}
