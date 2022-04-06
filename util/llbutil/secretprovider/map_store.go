package secretprovider

import (
	"context"

	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
)

type mapStore map[string][]byte

// GetSecret gets a secret from the map store
func (m mapStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	v, ok := m[id]
	if !ok {
		return nil, errors.WithStack(errors.Wrapf(secrets.ErrNotFound, "unable to lookup secret %s", id))
	}
	return v, nil
}
