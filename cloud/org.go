package cloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	secretsapi "github.com/earthly/cloud-api/secrets"
	"github.com/pkg/errors"
)

// OrgDetail contains an organization and details
type OrgDetail struct {
	ID       string
	Name     string
	Admin    bool
	Personal bool
}

// OrgPermissions contains permission details within an org
type OrgPermissions struct {
	User  string
	Path  string
	Write bool
}

// OrgMember represents a user that belongs to an org
type OrgMember struct {
	UserEmail  string
	Permission string
	OrgName    string
}

// ListOrgs lists all orgs a user has permission to view.
func (c *Client) ListOrgs(ctx context.Context) ([]*OrgDetail, error) {
	status, body, err := c.doCall(ctx, "GET", "/api/v0/admin/organizations?includePersonalOrg=true", withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list orgs: %s", msg)
	}

	var listOrgsResponse secretsapi.ListOrgsResponse
	err = c.jum.Unmarshal(body, &listOrgsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal list orgs response")
	}

	res := []*OrgDetail{}
	for _, org := range listOrgsResponse.Details {
		res = append(res, &OrgDetail{
			ID:       org.Id,
			Name:     org.Name,
			Admin:    org.Admin,
			Personal: org.Type == secretsapi.OrgType_PERSONAL,
		})
		c.orgIDCache.Store(org.Name, org.Id)
	}

	return res, nil
}

// Invite a user to an org.
func (c *Client) Invite(ctx context.Context, path, user string, write bool) error {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return errors.Errorf("invalid path")
	}

	permission := secretsapi.OrgPermissions{
		Path:  path,
		Email: user,
		Write: write,
	}

	u := fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName))

	status, body, err := c.doCall(ctx, "PUT", u, withAuth(), withJSONBody(&permission))
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to invite user into org: %s", msg)
	}
	return nil
}

// RevokePermission removes the org permission from the user.
func (c *Client) RevokePermission(ctx context.Context, path, user string) error {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return errors.Errorf("invalid path")
	}

	permission := secretsapi.OrgPermissions{
		Path:  path,
		Email: user,
	}

	status, body, err := c.doCall(ctx, "DELETE", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth(), withJSONBody(&permission))
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to revoke user from org: %s", msg)
	}
	return nil
}

// ListOrgPermissions returns all configured permissions for the org.
func (c *Client) ListOrgPermissions(ctx context.Context, path string) ([]*OrgPermissions, error) {
	orgName, ok := getOrgFromPath(path)
	if !ok {
		return nil, errors.Errorf("invalid path")
	}

	status, body, err := c.doCall(ctx, "GET", fmt.Sprintf("/api/v0/admin/organizations/%s/permissions", url.QueryEscape(orgName)), withAuth())
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return nil, errors.Errorf("failed to list org permissions: %s", msg)
	}

	var listOrgPermissionsResponse secretsapi.ListOrgPermissionsResponse
	err = c.jum.Unmarshal(body, &listOrgPermissionsResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal list org permissions response")
	}

	res := []*OrgPermissions{}
	for _, perm := range listOrgPermissionsResponse.Permissions {
		if strings.Contains(perm.Path, path) {
			res = append(res, &OrgPermissions{
				Path:  perm.Path,
				User:  perm.Email,
				Write: perm.Write,
			})
		}
	}

	return res, nil
}

// CreateOrg creates a new org by name.
func (c *Client) CreateOrg(ctx context.Context, org string) error {
	status, body, err := c.doCall(ctx, "PUT", fmt.Sprintf("/api/v0/admin/organizations/%s", url.QueryEscape(org)), withAuth())
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		msg, err := getMessageFromJSON(bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode response body (status code: %d)", status))
		}
		return errors.Errorf("failed to create org: %s", msg)
	}
	return nil
}

// GetOrgID retrieves the org ID for a named org.
func (c *Client) GetOrgID(ctx context.Context, orgName string) (string, error) {
	if orgID, ok := c.orgIDCache.Load(orgName); ok {
		return orgID.(string), nil
	}
	orgs, err := c.ListOrgs(ctx)
	if err != nil {
		return "", err
	}
	for _, o := range orgs {
		if o.Name == orgName {
			return o.ID, nil
		}
	}
	return "", errors.Errorf("org not found: %s", orgName)
}

// GuessOrgMembership returns an org name and ID if the user belongs to a single org
// Deprecated: we should stop "guessing" org membership and have the user always specify they want to use.
// A future `org select` command would make specifying the org easier.
func (c *Client) GuessOrgMembership(ctx context.Context) (orgName, orgID string, err error) {
	orgs, err := c.ListOrgs(ctx)
	if err != nil {
		return "", "", err
	}

	if len(orgs) == 2 {
		for _, o := range orgs {
			if !o.Personal {
				return o.Name, o.ID, nil
			}
		}
	}
	return "", "", errors.New("please specify the name of the organization using `--org`")
}

func getOrgFromPath(path string) (string, bool) {
	if path == "" || path[0] != '/' {
		return "", false
	}

	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 {
		return "", false
	}
	return parts[1], true
}
