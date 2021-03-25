package states

import "sync"

// CacheImports is a synchronized set of cache imports.
type CacheImports struct {
	mu    sync.Mutex
	store map[string]bool
}

// NewCacheImports creates a new cache imports structure.
func NewCacheImports(imports map[string]bool) *CacheImports {
	store := make(map[string]bool)
	for k, v := range imports {
		store[k] = v
	}
	return &CacheImports{
		store: store,
	}
}

// Add adds an import to the set.
func (ci *CacheImports) Add(tag string) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.store[tag] = true
}

// AsMap returns the cache imports contents as a map.
func (ci *CacheImports) AsMap() map[string]bool {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	copy := make(map[string]bool)
	for k, v := range ci.store {
		copy[k] = v
	}
	return copy
}
