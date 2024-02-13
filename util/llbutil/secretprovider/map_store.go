package secretprovider

import (
	"context"
	"errors"
	"net/url"

	"github.com/moby/buildkit/session/secrets"
)

type mapStore map[string][]byte

// NewMapStore returns a new map-based secret store
func NewMapStore(m map[string][]byte) secrets.SecretStore {
	return mapStore(m)
}

// GetSecret gets a secret from the map store
func (m mapStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	q, err := url.ParseQuery(id)
	if err != nil {
		return nil, errors.New("failed to parse secret ID")
	}

	if q.Get("name") == "" {
		return nil, errors.New("name parameter not found")
	}
	v, ok := m[q.Get("name")]
	if !ok {
		return nil, secrets.ErrNotFound
	}
	return v, nil
}
