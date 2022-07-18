package secretprovider

import (
	"context"
	"errors"
	"net/url"
	"path"
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

	q, err := url.ParseQuery(id)
	if err != nil {
		return nil, errors.New("failed to parse secret ID")
	}

	var data []byte

	name := q.Get("name")
	if name == "" {
		return nil, errors.New("name parameter not found")
	}

	switch q.Get("v") {
	case "0": // Legacy secret ID format includes the name only

		// This format requires the secret to be prefixed with either <org-name>
		// or 'user'.
		if !strings.Contains(name, "/") {
			return nil, secrets.ErrNotFound
		}

		// For the old secret format, there should never be a secret
		// that starts with a forward slash.
		if strings.HasPrefix(name, "/") {
			return nil, errors.New("secret name starts with '/'; this should never happen")
		}

		name = "/" + name
		data, err = cs.client.Get(ctx, name)
		if err != nil {
			return nil, err
		}

	case "1": // Project-based secret style includes the org and project name
		if !strings.HasPrefix(name, "user/") {
			name = path.Join(q.Get("org"), q.Get("project"), name)
		}
		name = "/" + name
		res, err := cs.client.ListSecrets(ctx, name)
		if err != nil {
			return nil, err
		}
		var match *cloud.Secret
		for _, sec := range res {
			if sec.Path == name {
				match = sec
				break
			}
		}
		if match == nil {
			return nil, secrets.ErrNotFound
		}
		data = []byte(match.Value)
	default:
		return nil, errors.New("invalid secret ID format")
	}

	return data, nil
}
