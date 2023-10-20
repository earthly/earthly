package leaselock

import "sync"

// LeaseLock is a lock that returns a lease when locked. The lease owner can unlock the lock.
type LeaseLock interface {
	Lock() LeaseLock
	Unlock() LeaseLock
	MaybeUnlock() LeaseLock
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

// MaybeUnlock has no effect on a lock without a lease.
func (l *lock) MaybeUnlock() LeaseLock {
	return l
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

// MaybeUnlock unlocks the lock and returns the lock.
func (l *lease) MaybeUnlock() LeaseLock {
	return l.Unlock()
}

// HaveLease returns true on a lease.
func (l *lease) HaveLease() bool {
	return true
}
