package llbutil

import (
	"context"
	"fmt"
	"strings"

	"github.com/earthly/earthly/cloud"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// ErrNoCloudClient occurs when the secrets client is referenced but was never provided
var ErrNoCloudClient = errors.Errorf("no secrets client provided")

type secretProvider struct {
	store  secrets.SecretStore
	client cloud.Client
}

// Register registers the secret provider
func (sp *secretProvider) Register(server *grpc.Server) {
	secrets.RegisterSecretsServer(server, sp)
}

func (sp *secretProvider) getSecretFromServer(path string) ([]byte, error) {
	if sp.client == nil {
		return nil, ErrNoCloudClient
	}
	data, err := sp.client.Get(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to lookup secret %q from secrets server", path))
	}
	return data, nil
}

// GetSecret returns a secret.
// secrets are referenced via +secret/name or +secret/org/name (or +secret/org/subdir1/.../name)
// however by the time GetSecret is called, the "+secret/" prefix is removed.
// if the name contains a /, then we can infer that it references the shared secret service.
func (sp *secretProvider) GetSecret(ctx context.Context, req *secrets.GetSecretRequest) (*secrets.GetSecretResponse, error) {
	isSharedSecret := false
	secretName := req.ID
	if strings.Contains(req.ID, "/") {
		isSharedSecret = true
		if req.ID[0] == '/' {
			panic("secret name starts with '/'; this should never happen")
		}
		secretName = "/" + req.ID
	}

	dt, err := sp.store.GetSecret(ctx, secretName)
	if err != nil {
		if errors.Is(err, secrets.ErrNotFound) && isSharedSecret {
			dt, err = sp.getSecretFromServer(secretName)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &secrets.GetSecretResponse{
		Data: dt,
	}, nil
}

// NewSecretProvider returns a new secrets provider
func NewSecretProvider(client cloud.Client, overrides map[string][]byte) session.Attachable {
	return &secretProvider{
		store:  mapStore(overrides),
		client: client,
	}
}

type mapStore map[string][]byte

// GetSecret gets a secret from the map store
func (m mapStore) GetSecret(ctx context.Context, id string) ([]byte, error) {
	v, ok := m[id]
	if !ok {
		return nil, errors.WithStack(errors.Wrapf(secrets.ErrNotFound, "unable to lookup secret %s", id))
	}
	return v, nil
}
