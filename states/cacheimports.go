package states

import "sync"

// CacheImports is a synchronized set of cache imports.
type CacheImports struct {
	mu    sync.RWMutex
	store []string
}

// NewCacheImports creates a new cache imports structure.
func NewCacheImports(imports []string) *CacheImports {
	clone := make([]string, len(imports))
	copy(clone, imports)

	return &CacheImports{
		store: clone,
	}
}

// Add adds imports to the set.
func (ci *CacheImports) Add(tags ...string) {
	ci.mu.Lock()
	defer ci.mu.Unlock()

	ci.store = append(ci.store, tags...)
}

// AsSlice returns the cache imports contents as a slice.
func (ci *CacheImports) AsSlice() []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	clone := make([]string, len(ci.store))
	copy(clone, ci.store)

	return clone
}
