package synccache

import (
	"sync"
)

// Constructor is a func that is used to construct a cache value, given a key.
type Constructor func(interface{}) (interface{}, error)

type entry struct {
	mu          sync.Mutex
	value       interface{}
	err         error
	constructed bool
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
func (sc *SyncCache) Do(key interface{}, c Constructor) (interface{}, error) {
	e := sc.getEntry(key)
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.constructed {
		e.value, e.err = c(key)
		e.constructed = true
	}
	return e.value, e.err
}

func (sc *SyncCache) getEntry(key interface{}) *entry {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	e, ok := sc.store[key]
	if !ok {
		e = new(entry)
		sc.store[key] = e
	}
	return e
}
