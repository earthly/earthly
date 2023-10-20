package leaselock

import "sync"

// LeaseLock is a lock that returns a lease when locked. The lease owner can unlock the lock.
type LeaseLock interface {
	// Lock locks the lock and returns a lease. If the caller is already the lease owner, locking
	// panics.
	Lock() LeaseLock
	// Unlock unlocks the lock and returns the lock. If the caller is not the lease owner, unlocking
	// panics.
	Unlock() LeaseLock
	// HaveLease returns true if the caller is the lease owner.
	HaveLease() bool
}

type lock struct {
	mu sync.Mutex
}

// New returns a new lease lock.
func New() LeaseLock {
	return &lock{}
}

// Lock locks the lock and returns a lease.
func (l *lock) Lock() LeaseLock {
	l.mu.Lock()
	return &lease{lock: l}
}

// Unlock fails on a lock without a lease.
func (l *lock) Unlock() LeaseLock {
	panic("unlock called without lease")
}

// HaveLease returns false on a lock without a lease.
func (l *lock) HaveLease() bool {
	return false
}

type lease struct {
	lock *lock
}

// Lock fails on a lease.
func (l *lease) Lock() LeaseLock {
	panic("lock called with lease")
}

// Unlock unlocks the lock and returns the lock.
func (l *lease) Unlock() LeaseLock {
	l.lock.mu.Unlock()
	return l.lock
}

// HaveLease returns true on a lease.
func (l *lease) HaveLease() bool {
	return true
}
