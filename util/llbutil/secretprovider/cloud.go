package secretprovider

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/moby/buildkit/session/secrets"

	"github.com/earthly/earthly/cloud"
)

type cloudStore struct {
	client *cloud.Client
}

// NewCloudStore returns a new cloud secret store
func NewCloudStore(client *cloud.Client) secrets.SecretStore {
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
		return cs.client.Get(ctx, name)

	case "1": // Project-based secret style includes the org and project name
		if strings.HasPrefix(name, "user/") {
			secret, err := cs.client.GetUserSecret(ctx, strings.TrimPrefix(name, "user/"))
			if err != nil {
				return nil, err
			}
			return []byte(secret.Value), nil
		}

		org := q.Get("org")
		project := q.Get("project")
		if org == "" || project == "" {
			return nil, secrets.ErrNotFound
		}

		secret, err := cs.client.GetProjectSecret(ctx, org, project, name)
		if err != nil {
			return nil, err
		}
		return []byte(secret.Value), nil

	default:
		return nil, errors.New("invalid secret ID format")
	}
}
