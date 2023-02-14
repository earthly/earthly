package cloud

import (
	"context"
	"fmt"
	"net/http"
	"time"

	secretsapi "github.com/earthly/cloud-api/secrets"
	"github.com/pkg/errors"
)

// Project contains information about the org project.
type Project struct {
	ID         string
	Name       string
	OrgName    string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// ProjectMember contains information about the project member.
type ProjectMember struct {
	UserID     string
	UserEmail  string
	UserName   string
	Permission string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// CreateProject creates a new project within the specified organization.
func (c *Client) CreateProject(ctx context.Context, name, orgName string) (*Project, error) {
	u := "/api/v0/projects"

	req := &secretsapi.CreateProjectRequest{
		Project: &secretsapi.Project{
			OrgName: orgName,
			Name:    name,
		},
	}

	status, body, err := c.doCall(ctx, http.MethodPost, u, withAuth(), withJSONBody(req))
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		return nil, errors.Errorf("failed to create project: %s", body)
	}

	res := &secretsapi.CreateProjectResponse{}

	err = c.jum.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return &Project{
		Name:       res.Project.Name,
		OrgName:    res.Project.OrgName,
		CreatedAt:  res.Project.CreatedAt.AsTime(),
		ModifiedAt: res.Project.ModifiedAt.AsTime(),
	}, nil
}

// ListProjects returns all projects in the organization that are visible to the
// logged-in user.
func (c *Client) ListProjects(ctx context.Context, orgName string) ([]*Project, error) {
	u := fmt.Sprintf("/api/v0/projects/%s", orgName)

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to list projects: %s", body)
	}

	res := &secretsapi.ListProjectsResponse{}

	err = c.jum.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	var projects []*Project

	for _, pro := range res.Projects {
		projects = append(projects, &Project{
			Name:       pro.Name,
			OrgName:    pro.OrgName,
			CreatedAt:  pro.CreatedAt.AsTime(),
			ModifiedAt: pro.ModifiedAt.AsTime(),
		})
	}

	return projects, nil
}

// GetProject loads a single project from the projects endpoint.
func (c *Client) GetProject(ctx context.Context, orgName, name string) (*Project, error) {
	u := fmt.Sprintf("/api/v0/projects/%s/%s", orgName, name)

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to get project: %s", body)
	}

	res := &secretsapi.GetProjectResponse{}

	err = c.jum.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return &Project{
		Name:       res.Project.Name,
		OrgName:    res.Project.OrgName,
		CreatedAt:  res.Project.CreatedAt.AsTime(),
		ModifiedAt: res.Project.ModifiedAt.AsTime(),
	}, nil
}

// DeleteProject deletes a given project by name.
func (c *Client) DeleteProject(ctx context.Context, orgName, name string) error {
	u := fmt.Sprintf("/api/v0/projects/%s/%s", orgName, name)

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to delete project: %s", body)
	}

	return nil
}

// AddProjectMember adds a new member to the project by email or user ID.
func (c *Client) AddProjectMember(ctx context.Context, orgName, name, userEmail, permission string) error {
	u := fmt.Sprintf("/api/v0/projects/%s/%s/members", orgName, name)

	req := &secretsapi.AddProjectMemberRequest{
		Permission: permission,
		UserEmail:  userEmail,
	}

	status, body, err := c.doCall(ctx, http.MethodPost, u, withAuth(), withJSONBody(req))
	if err != nil {
		return err
	}

	if status != http.StatusCreated {
		return errors.Errorf("failed to add member to project: %s", body)
	}

	return nil
}

// UpdateProjectMember updates an existing member with the new permission
func (c *Client) UpdateProjectMember(ctx context.Context, orgName, name, userEmail, permission string) error {
	u := fmt.Sprintf("/api/v0/projects/%s/%s/members/%s", orgName, name, userEmail)

	req := &secretsapi.AddProjectMemberRequest{
		Permission: permission,
		UserEmail:  userEmail,
	}

	status, body, err := c.doCall(ctx, http.MethodPut, u, withAuth(), withJSONBody(req))
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to update member: %s", body)
	}

	return nil
}

// ListProjectMembers will return all project members if the user has permission to do so.
func (c *Client) ListProjectMembers(ctx context.Context, orgName, name string) ([]*ProjectMember, error) {
	u := fmt.Sprintf("/api/v0/projects/%s/%s/members", orgName, name)

	status, body, err := c.doCall(ctx, http.MethodGet, u, withAuth())
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, errors.Errorf("failed to list project members: %s", body)
	}

	var members []*ProjectMember

	res := &secretsapi.ListProjectMembersResponse{}

	err = c.jum.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	for _, m := range res.Members {
		members = append(members, &ProjectMember{
			UserName:   m.UserName,
			UserEmail:  m.UserEmail,
			Permission: m.Permission,
			CreatedAt:  m.CreatedAt.AsTime(),
			ModifiedAt: m.ModifiedAt.AsTime(),
		})
	}

	return members, nil
}

// RemoveProjectMember will remove a member from a project.
func (c *Client) RemoveProjectMember(ctx context.Context, orgName, name, userEmail string) error {
	u := fmt.Sprintf("/api/v0/projects/%s/%s/members/%s", orgName, name, userEmail)

	status, body, err := c.doCall(ctx, http.MethodDelete, u, withAuth())
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.Errorf("failed to remove a project member: %s", body)
	}

	return nil
}
