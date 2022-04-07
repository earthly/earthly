package secretprovider

import (
	"context"
	"strings"

	"github.com/earthly/earthly/cloud"
	"github.com/moby/buildkit/session/secrets"
)

type cloudStore struct {
	client cloud.Client
}

// NewCloudStore returns a new cloud secret store
func NewCloudStore(client cloud.Client) secrets.SecretStore {
	return &cloudStore{
		client: client,
	}
}

// GetSecret returns a secret.
// secrets are referenced via +secret/name or +secret/org/name (or +secret/org/subdir1/.../name)
// however by the time GetSecret is called, the "+secret" prefix is removed.
func (cs *cloudStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	if !strings.HasPrefix(id, "/") {
		return nil, secrets.ErrNotFound
	}
	dt, err := cs.client.Get(id)
	if err != nil {
		return nil, err
	}
	return dt, nil
}
