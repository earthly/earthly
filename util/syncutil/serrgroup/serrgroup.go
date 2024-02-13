// Package serrgroup is an error group that does not crash if work is added after
// Wait has returned after an error. In addition, the error is exposed via the
// API.
// This is based on
// https://cs.opensource.google/go/x/sync/+/036812b2:errgroup/errgroup.go
// with some modifications.
package serrgroup

import (
	"context"
	"sync"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
	errMu   sync.Mutex
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Err returns the first non-nil error (if any) of all goroutines started by Go.
func (g *Group) Err() error {
	g.errMu.Lock()
	defer g.errMu.Unlock()
	return g.err
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.Err()
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func() error) {
	g.errMu.Lock()
	if g.err != nil {
		g.errMu.Unlock()
		// Don't add more work if there has been an error.
		return
	}
	g.errMu.Unlock()

	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.errMu.Lock()
				defer g.errMu.Unlock()
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
