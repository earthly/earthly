package metacontext

import (
	"context"
	"sync"
	"time"
)

var _ context.Context = &MetaContext{}

// MetaContext is an object which implements context.Context and which holds multiple
// contexts within it. The MetaContext is considered canceled only when ALL of the
// underlying contexts have been canceled.
//
// Once canceled, it cannot be uncancelled, so it is an error to keep adding contexts
// once the meta context is considered cancelled.
type MetaContext struct {
	subDoneCh chan int // index

	mu      sync.Mutex
	doneCh  chan struct{}
	numDone int
	sub     []context.Context

	firstDoneMu  sync.Mutex
	firstDoneErr error
}

// New returns a new metacontext.
func New(ctx context.Context) *MetaContext {
	mc := &MetaContext{
		doneCh:    make(chan struct{}),
		subDoneCh: make(chan int),
	}
	_ = mc.Add(ctx)
	go mc.monitor()
	return mc
}

func (mc *MetaContext) monitor() {
	for index := range mc.subDoneCh {
		mc.mu.Lock()
		mc.numDone++
		if mc.numDone == 1 {
			firstDoneCtx := mc.sub[index]
			mc.firstDoneMu.Lock()
			go func() {
				// Call .Err() outside of our lock. Also, use a different lock
				// to block a caller to our .Err if it'll take a long time.
				defer mc.firstDoneMu.Unlock()
				err := firstDoneCtx.Err()
				mc.firstDoneErr = err
			}()
		}
		if mc.numDone == len(mc.sub) {
			close(mc.doneCh)
			mc.mu.Unlock()
			return
		}
		mc.mu.Unlock()
	}
}

// Add adds a new context to the metacontext.
func (mc *MetaContext) Add(ctx context.Context) error {
	mc.mu.Lock()
	select {
	case <-mc.doneCh:
		mc.mu.Unlock()
		mc.firstDoneMu.Lock()
		defer mc.firstDoneMu.Unlock()
		return mc.firstDoneErr
	default:
	}
	mc.sub = append(mc.sub, ctx)
	index := len(mc.sub) - 1
	mc.mu.Unlock()
	go func() {
		<-ctx.Done()
		mc.subDoneCh <- index
	}()
	return nil
}

// Deadline returns the earliest Deadline in the pool.
func (mc *MetaContext) Deadline() (deadline time.Time, ok bool) {
	mc.mu.Lock()
	copy := append([]context.Context{}, mc.sub...)
	mc.mu.Unlock()
	if len(copy) == 0 {
		return time.Time{}, false
	}
	min := time.Time{}
	hasDl := false
	for _, ctx := range copy {
		dl, ok := ctx.Deadline() // don't hold lock for this call
		if ok {
			if !hasDl || dl.Before(min) {
				min = dl
			}
			hasDl = true
		}
	}
	return min, hasDl
}

// Done returns the done channel. The MetaContext is done only when ALL of the
// contained contexts are done.
func (mc *MetaContext) Done() <-chan struct{} {
	return mc.doneCh
}

// Err returns the first done error reported by any context, if the whole
// context is done. Nil otherwise.
func (mc *MetaContext) Err() error {
	select {
	case <-mc.doneCh:
		mc.firstDoneMu.Lock()
		defer mc.firstDoneMu.Unlock()
		if mc.firstDoneErr != nil {
			return mc.firstDoneErr
		}
		return context.Canceled
	default:
		return nil
	}
}

// Value calls context.Value on the first not-done context, or on the first one,
// if all are done.
func (mc *MetaContext) Value(key interface{}) interface{} {
	mc.mu.Lock()
	if len(mc.sub) == 0 {
		mc.mu.Unlock()
		return nil
	}
	// Find the first not-done ctx. If none found, use the first one.
	var selectedCtx context.Context
	for _, ctx := range mc.sub {
		select {
		case <-mc.doneCh:
			continue
		default:
		}
		selectedCtx = ctx
		break
	}
	if selectedCtx == nil {
		selectedCtx = mc.sub[0]
	}
	mc.mu.Unlock()
	return selectedCtx.Value(key) // don't hold lock for this call
}
