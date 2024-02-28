package secretprovider

import (
	"context"
	"net/url"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/hint"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

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

	v, err := url.ParseQuery(req.ID)
	if err != nil {
		return nil, errors.New("failed to parse secret ID")
	}

	for _, store := range sp.stores {
		data, err := store.GetSecret(ctx, req.ID)
		if err != nil {
			if errors.Is(err, secrets.ErrNotFound) || errors.Is(err, cloud.ErrNotFound) {
				continue
			}
			return nil, err
		}
		return &secrets.GetSecretResponse{
			Data: data,
		}, nil
	}

	return nil, hint.Wrap(errors.WithStack(errors.Wrapf(secrets.ErrNotFound, "unable to lookup secret %q", v.Get("name"))),
		"Make sure to set the project at the top of the Earthfile by using the PROJECT command.",
		"Note, if this secret was called from a FUNCTION, the project needs to be set in the Earthfile that calls the FUNCTION.")

}

// New returns a new secrets provider which looks up secrets
// in each supplied secret store (ordered by argument ordering) and returns the first found secret
func New(stores ...secrets.SecretStore) session.Attachable {
	return &secretProvider{
		stores: stores,
	}
}
