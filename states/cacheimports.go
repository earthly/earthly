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
	slice := make([]string, 0, len(imports))
	store := make(map[string]bool)

	for _, tag := range imports {
		if _, exists := store[tag]; exists {
			continue
		}

		store[tag] = true
		slice = append(slice, tag)
	}

	return &CacheImports{
		slice: slice,
		store: store,
	}
}

// Add adds import to the set.
func (ci *CacheImports) Add(tag string) {
	ci.mu.Lock()
	defer ci.mu.Unlock()

	if _, exists := ci.store[tag]; exists {
		return
	}

	ci.store[tag] = true
	ci.slice = append(ci.slice, tag)
}

// Has checks if a passed tag is added.
func (ci *CacheImports) Has(tag string) bool {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	_, exists := ci.store[tag]

	return exists
}

// AsSlice returns the cache imports contents as a slice.
func (ci *CacheImports) AsSlice() []string {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	clone := make([]string, len(ci.slice))
	copy(clone, ci.slice)

	return clone
}
