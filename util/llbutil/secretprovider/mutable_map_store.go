package secretprovider

import (
	"context"
	"sync"

	"github.com/moby/buildkit/session/secrets"
)

var _ secrets.SecretStore = &MutableMapStore{}

// MutableMapStore is a secret store which can be mutated.
type MutableMapStore struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewMutableMapStore returns a new map-based mutable secret store.
func NewMutableMapStore(m map[string][]byte) *MutableMapStore {
	if m == nil {
		m = make(map[string][]byte)
	}
	return &MutableMapStore{
		store: m,
	}
}

// GetSecret gets a secret from the map store.
func (m *MutableMapStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.store[id]
	if !ok {
		return nil, secrets.ErrNotFound
	}
	return v, nil
}

// SetSecret sets a secret in the map store.
func (m *MutableMapStore) SetSecret(ctx context.Context, id string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[id] = data
	return nil
}

// DeleteSecret deletes a secret from the map store.
func (m *MutableMapStore) DeleteSecret(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, id)
	return nil
}
