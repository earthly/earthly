package states

import "sync"

// CacheImports is a synchronized set of cache imports.
type CacheImports struct {
	mu    sync.RWMutex
	slice []string
	store map[string]bool
}

// NewCacheImports creates a new cache imports structure.
func NewCacheImports(imports []string) *CacheImports {
	clone := make([]string, len(imports))
	copy(clone, imports)

	store := make(map[string]bool)
	for _, tag := range imports {
		store[tag] = true
	}

	return &CacheImports{
		slice: clone,
		store: store,
	}
}

// Add adds imports to the set.
func (ci *CacheImports) Add(tags ...string) {
	ci.mu.Lock()
	defer ci.mu.Unlock()

	ci.slice = append(ci.slice, tags...)

	for _, tag := range tags {
		ci.store[tag] = true
	}
}

// Has checks if a passed tag is added.
func (ci *CacheImports) Has(tag string) bool {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	_, ok := ci.store[tag]

	return ok
}

// AsSlice returns the cache imports contents as a slice.
func (ci *CacheImports) AsSlice() []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	clone := make([]string, len(ci.slice))
	copy(clone, ci.slice)

	return clone
}
