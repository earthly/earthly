package secretprovider

import (
	"context"
	"strings"

	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// InternalSecretPrefix is a prefix used for Earthly-internal secrets.
const InternalSecretPrefix = "52804da5-2787-46ad-8478-80c50f305e76/"

type secretProvider struct {
	stores []secrets.SecretStore
}

// Register registers the secret provider
func (sp *secretProvider) Register(server *grpc.Server) {
	secrets.RegisterSecretsServer(server, sp)
}

// GetSecret returns a secret.
// secrets are referenced via +secret/name or +secret/org/name (or +secret/org/subdir1/.../name)
// however by the time GetSecret is called, the "+secret/" prefix is removed.
// if the name contains a /, then we can infer that it references the shared secret service.
func (sp *secretProvider) GetSecret(ctx context.Context, req *secrets.GetSecretRequest) (*secrets.GetSecretResponse, error) {

	// shared secrets will be of the form org/path
	// and must be transformed into /org/path
	secretName := req.ID
	if strings.Contains(secretName, "/") {
		if req.ID[0] == '/' {
			panic("secret name starts with '/'; this should never happen")
		}
		secretName = "/" + secretName
	}

	for _, store := range sp.stores {
		dt, err := store.GetSecret(ctx, secretName)
		if err != nil {
			if errors.Is(err, secrets.ErrNotFound) {
				continue
			}
			return nil, err
		}
		return &secrets.GetSecretResponse{
			Data: dt,
		}, nil
	}

	return nil, errors.WithStack(errors.Wrapf(secrets.ErrNotFound, "unable to lookup secret %s", secretName))
}

// New returns a new secrets provider which looks up secrets
// in each supplied secret store (ordered by argument ordering) and returns the first found secret
func New(stores ...secrets.SecretStore) session.Attachable {
	return &secretProvider{
		stores: stores,
	}
}
