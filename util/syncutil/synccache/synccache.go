package synccache

import (
	"context"
	"sync"

	"github.com/earthly/earthly/util/syncutil/metacontext"
	"github.com/pkg/errors"
)

// Constructor is a func that is used to construct a cache value, given a key.
type Constructor func(ctx context.Context, key interface{}) (value interface{}, err error)

type entry struct {
	metaCtx *metacontext.MetaContext

	constructed chan struct{}
	err         error
	value       interface{}
}

// SyncCache is an object which can be used to create singletons stored in a key-value store.
type SyncCache struct {
	mu    sync.Mutex
	store map[interface{}]*entry
}

// New creates an empty SyncCache.
func New() *SyncCache {
	return &SyncCache{
		store: make(map[interface{}]*entry),
	}
}

// Do executes the constructor, if a value for key hasn't already been constructed.
func (sc *SyncCache) Do(ctx context.Context, key interface{}, c Constructor) (interface{}, error) {
	e, found := sc.getEntry(ctx, key)
	if !found {
		// We need to construct this.
		go func() {
			// The metaCtx will ensure that this stays alive even if the original Do has
			// been canceled, thanks to the metaCtx. This is canceled only when ALL of
			// the Do's are canceled.
			e.value, e.err = c(e.metaCtx, key)
			// Don't cache context canceled. Whoever is currently waiting will still get this,
			// but no future callers to Do will.
			if errors.Is(e.err, context.Canceled) {
				sc.deleteEntry(key)
			}
			close(e.constructed)
		}()
	} else {
		go func() {
			select {
			case <-e.constructed:
			default:
				// Add our context to metaCtx in case all others get canceled.
				_ = e.metaCtx.Add(ctx)
			}
		}()
	}
	<-e.constructed
	return e.value, e.err
}

// Add adds a readily constructed value for a given key.
func (sc *SyncCache) Add(ctx context.Context, key interface{}, value interface{}, valueErr error) error {
	e, found := sc.getEntry(ctx, key)
	if found {
		return errors.New("already exists")
	}
	e.value = value
	e.err = valueErr
	close(e.constructed)
	return nil
}

func (sc *SyncCache) getEntry(ctx context.Context, key interface{}) (*entry, bool) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	e, ok := sc.store[key]
	if !ok {
		e = &entry{
			metaCtx:     metacontext.New(ctx),
			constructed: make(chan struct{}),
		}
		sc.store[key] = e
	}
	return e, ok
}

func (sc *SyncCache) deleteEntry(key interface{}) {
	// note; this does not cancel any ongoing construction.
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.store, key)
}
