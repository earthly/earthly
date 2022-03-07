package semutil

import "context"

// ReleaseFun is a function that needs to be called to release the semaphore.
type ReleaseFun func()

// Semaphore is a generic semaphore.
type Semaphore interface {
	// Acquire acquires a resource on the semaphore. If no resource is available,
	// the call blocks until one is made available.
	Acquire(ctx context.Context, n int64) (ReleaseFun, error)
	// TryAcquire acquires a resource on the semaphore. If no resource is
	// available, the call returns immediately with false.
	TryAcquire(n int64) (ReleaseFun, bool)
}
