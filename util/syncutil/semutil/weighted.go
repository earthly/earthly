package semutil

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// Weighted is a weighted semaphore.
type Weighted struct {
	sem *semaphore.Weighted
}

// NewWeighted creates a new weighted semaphore.
func NewWeighted(n int64) Semaphore {
	return &Weighted{sem: semaphore.NewWeighted(n)}
}

// Acquire acquires a semaphore.
func (w *Weighted) Acquire(ctx context.Context, n int64) (ReleaseFun, error) {
	err := w.sem.Acquire(ctx, n)
	if err != nil {
		return nil, err
	}
	return func() {
		w.sem.Release(n)
	}, nil
}

// TryAcquire acquires a semaphore, but does not block.
func (w *Weighted) TryAcquire(n int64) (ReleaseFun, bool) {
	ok := w.sem.TryAcquire(n)
	if !ok {
		return nil, false
	}
	return func() {
		w.sem.Release(n)
	}, true
}
