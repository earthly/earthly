package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
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
	UserID     string
	Permission string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func (c *client) Remove(ctx context.Context, path string) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "DELETE", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to remove secret: %s", msg)
	}
	return nil
}

func (c *client) List(ctx context.Context, path string) ([]string, error) {
	if path != "" && !strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list secrets: %s", msg)
	}
	if body == "" {
		return []string{}, nil
	}
	return strings.Split(body, "\n"), nil
}

func (c *client) Get(ctx context.Context, path string) ([]byte, error) {
	if path == "" || path[0] != '/' || strings.HasSuffix(path, "/") {
		return nil, errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/secrets%s", path), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to get secret: %s", msg)
	}
	return []byte(body), nil
}

func (c *client) Set(ctx context.Context, path string, data []byte) error {
	if path == "" || path[0] != '/' {
		return errors.Errorf("invalid path")
	}
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/secrets%s", path), withAuth(), withBody(string(data)))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader([]byte(body)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to set secret: %s", msg)
	}
	return nil
}

// ListSecrets returns a list of secrets base on the given path.
func (c *client) ListSecrets(ctx context.Context, path string) ([]*Secret, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets" + path

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to list secrets: %s", body)
	}

	var secrets []*Secret

	res := &secretsapi.ListSecretsResponse{}
	err = jsonpb.UnmarshalString(body, res)
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

// SetSecret adds or updates the given path and secret combination.
func (c *client) SetSecret(ctx context.Context, path string, secret []byte) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets" + path

	status, body, err := c.doCall(ctx, http.MethodPut, u, withAuth(), withBody(string(secret)))
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to set secret: %s", body)
	}

	return nil
}

// RemoveSecret deletes a secret by path name.
func (c *client) RemoveSecret(ctx context.Context, path string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets" + path

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to delete secret: %s", body)
	}

	return nil
}

// ListSecretPermissions returns a set of user permissions for project secrets.
func (c *client) ListSecretPermissions(ctx context.Context, path string) ([]*SecretPermission, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets/permissions" + path

	status, body, err := c.doCall(ctx, http.MethodPut, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to set secret: %s", body)
	}

	var secretPerms []*SecretPermission

	res := &secretsapi.ListSecretPermissionsResponse{}
	err = jsonpb.UnmarshalString(body, res)
	if err != nil {
		return nil, err
	}

	for _, perm := range res.SecretPermissions {
		secretPerms = append(secretPerms, &SecretPermission{
			UserID:     perm.UserId,
			Path:       perm.Path,
			Permission: perm.Permission,
			CreatedAt:  perm.CreatedAt.AsTime(),
			ModifiedAt: perm.ModifiedAt.AsTime(),
		})
	}

	return secretPerms, nil
}

// SetSecretPermission is used to set a user permission on a given secret path.
func (c *client) SetSecretPermission(ctx context.Context, path, userID, permission string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets/permissions" + path

	req := &secretsapi.UpdateSecretPermissionRequest{
		UserId:     userID,
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
func (c *client) RemoveSecretPermission(ctx context.Context, path, userID string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := "/api/v1/secrets/permissions" + path + "/" + userID

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to delete secret permission: %s", body)
	}

	return nil
}
