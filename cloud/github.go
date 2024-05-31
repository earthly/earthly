package cloud

import (
	"context"
	"fmt"

	pb "github.com/earthly/cloud-api/compute"
)

func (c *Client) CreateGHAIntegration(ctx context.Context, orgName string, ghOrg string, ghRepo string, token string) error {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return fmt.Errorf("failed getting org ID: %w", err)
	}
	_, err = c.compute.CreateGHAIntegration(c.withAuth(ctx), &pb.CreateGHAIntegrationRequest{
		OrgId:          orgID,
		GithubOrgName:  ghOrg,
		GithubRepoName: ghRepo,
		GithubToken:    token,
	})
	if err != nil {
		return fmt.Errorf("failed to create GitHub integration: %w", err)
	}
	return nil
}

func (c *Client) RemoveGHAIntegration(ctx context.Context, orgName string, ghOrg string, ghRepo string) error {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return fmt.Errorf("failed getting org ID: %w", err)
	}
	_, err = c.compute.RemoveGHAIntegration(c.withAuth(ctx), &pb.RemoveGHAIntegrationRequest{
		OrgId:          orgID,
		GithubOrgName:  ghOrg,
		GithubRepoName: ghRepo,
	})
	if err != nil {
		return fmt.Errorf("failed remove GitHub integration: %w", err)
	}
	return nil
}

func (c *Client) ListGHAIntegrations(ctx context.Context, orgName string) (*pb.ListGHAIntegrationsResponse, error) {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return nil, fmt.Errorf("failed getting org ID: %w", err)
	}
	integrations, err := c.compute.ListGHAIntegrations(c.withAuth(ctx), &pb.ListGHAIntegrationsRequest{
		OrgId: orgID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed remove GitHub integration: %w", err)
	}
	return integrations, nil
}
