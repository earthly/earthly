package secretprovider

import (
	"context"

	"github.com/moby/buildkit/session/secrets"
)

type mapStore map[string][]byte

// NewMapStore returns a new map-based secret store
func NewMapStore(m map[string][]byte) secrets.SecretStore {
	return mapStore(m)
}

// GetSecret gets a secret from the map store
func (m mapStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	v, ok := m[id]
	if !ok {
		return nil, secrets.ErrNotFound
	}
	return v, nil
}
