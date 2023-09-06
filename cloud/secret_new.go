package cloud

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	secretsapi "github.com/earthly/cloud-api/secrets"
	httperror "github.com/earthly/earthly/error/http"

	"github.com/pkg/errors"
)

// Secret represents a Cloud secret with a path key and a string value.
type Secret struct {
	Path       string
	Value      string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// SecretPermission contains information about a user-specific secret
// permission override.
type SecretPermission struct {
	Path       string
	UserEmail  string
	Permission string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// ListSecrets returns a list of secrets base on the given path.
func (c *Client) ListSecrets(ctx context.Context, path string) ([]*Secret, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets" + path

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("failed to list secrets: %w", httperror.New(status, string(body)))
	}

	var secrets []*Secret

	res := &secretsapi.ListSecretsResponse{}
	err = c.jum.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	for _, secret := range res.Secrets {
		secrets = append(secrets, &Secret{
			Path:       secret.Path,
			Value:      secret.Value,
			CreatedAt:  secret.CreatedAt.AsTime(),
			ModifiedAt: secret.ModifiedAt.AsTime(),
		})
	}

	return secrets, nil
}

// GetProjectSecret gets a secret from a project secret store
func (c *Client) GetProjectSecret(ctx context.Context, org, project, secretName string) (*Secret, error) {
	if org == "" {
		return nil, fmt.Errorf("GetProjectSecret called with empty org")
	}
	if project == "" {
		return nil, fmt.Errorf("GetProjectSecret called with empty project")
	}
	if secretName == "" {
		return nil, fmt.Errorf("GetProjectSecret called with empty secretName")
	}
	return c.getSecretV2(ctx, fmt.Sprintf("/%s/%s/%s", org, project, secretName))
}

// GetUserSecret gets a secret from the current user's personal secret store
func (c *Client) GetUserSecret(ctx context.Context, secretName string) (*Secret, error) {
	if secretName == "" {
		return nil, fmt.Errorf("GetUserSecret called with empty secretName")
	}
	return c.getSecretV2(ctx, fmt.Sprintf("/user/%s", secretName))
}

// GetUserOrProjectSecret gets a secret from a project or the current user's personal secret store,
// depending on the path being structured as /org/project/... or /user/... respectively.
func (c *Client) GetUserOrProjectSecret(ctx context.Context, path string) (*Secret, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, ErrMalformedSecretPath
	}
	// path must be /user/..., or /org/project/...
	if !strings.HasPrefix(path, "/user/") && strings.Count(path, "/") < 3 {
		return nil, ErrMalformedSecretPath
	}
	return c.getSecretV2(ctx, path)
}

func (c *Client) getSecretV2(ctx context.Context, path string) (*Secret, error) {
	res, err := c.ListSecrets(ctx, path)
	if err != nil {
		if httperror.Code(err) == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	for _, sec := range res {
		if sec.Path == path {
			return sec, nil
		}
	}
	return nil, ErrNotFound
}

// SetSecret adds or updates the given path and secret combination.
func (c *Client) SetSecret(ctx context.Context, path string, secret []byte) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets" + path

	status, body, err := c.doCall(ctx, http.MethodPut, u, withAuth(), withBody(secret))
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to set secret: %s", body)
	}

	return nil
}

// RemoveSecret deletes a secret by path name.
func (c *Client) RemoveSecret(ctx context.Context, path string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets" + path

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	switch status {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return errors.Errorf("failed to delete secret: %s", body)
	}

	return nil
}

// ListSecretPermissions returns a set of user permissions for project secrets.
func (c *Client) ListSecretPermissions(ctx context.Context, path string) ([]*SecretPermission, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets/permissions" + path

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to list secret permissions: %s", body)
	}

	var secretPerms []*SecretPermission

	res := &secretsapi.ListSecretPermissionsResponse{}
	err = c.jum.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	for _, perm := range res.SecretPermissions {
		secretPerms = append(secretPerms, &SecretPermission{
			UserEmail:  perm.UserEmail,
			Path:       perm.Path,
			Permission: perm.Permission,
			CreatedAt:  perm.CreatedAt.AsTime(),
			ModifiedAt: perm.ModifiedAt.AsTime(),
		})
	}

	return secretPerms, nil
}

// SetSecretPermission is used to set a user permission on a given secret path.
func (c *Client) SetSecretPermission(ctx context.Context, path, userEmail, permission string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets/permissions" + path

	req := &secretsapi.UpdateSecretPermissionRequest{
		UserEmail:  userEmail,
		Permission: permission,
	}

	status, body, err := c.doCall(ctx, http.MethodPut, u, withAuth(), withJSONBody(req))
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to set secret permission: %s", body)
	}

	return nil
}

// RemoveSecretPermission removes a secret permission for the user and path.
func (c *Client) RemoveSecretPermission(ctx context.Context, path, userEmail string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets/permissions" + path + "/" + userEmail

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to delete secret permission: %s", body)
	}

	return nil
}
