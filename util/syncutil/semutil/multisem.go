package semutil

import (
	"context"
)

// MultiSem is a semaphore that is made out of multiple underlying semaphores.
// The semaphores are attempted one at a time.
type MultiSem struct {
	sems []Semaphore
}

// NewMultiSem creates a new MultiSem.
func NewMultiSem(sems ...Semaphore) Semaphore {
	if len(sems) == 0 {
		panic("no semaphores passed")
	}
	return &MultiSem{
		sems: sems,
	}
}

// Acquire acquires a semaphore. If all semaphores are starved, it only blocks
// on the last semaphore.
func (ms *MultiSem) Acquire(ctx context.Context, n int64) (ReleaseFun, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for _, sem := range ms.sems {
		rel, ok := sem.TryAcquire(n)
		if ok {
			return rel, nil
		}
	}
	lastSem := ms.sems[len(ms.sems)-1]
	rel, err := lastSem.Acquire(ctx, n)
	if err != nil {
		return nil, err
	}
	return rel, nil
}

// TryAcquire acquires a semaphore, but does not block.
func (ms *MultiSem) TryAcquire(n int64) (ReleaseFun, bool) {
	for _, sem := range ms.sems {
		rel, ok := sem.TryAcquire(n)
		if ok {
			return rel, true
		}
	}
	return nil, false
}
