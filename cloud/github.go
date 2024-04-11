package cloud

import (
	"context"
	"fmt"

	pb "github.com/earthly/cloud-api/compute"
)

func (c *Client) SetGithubToken(ctx context.Context, orgName string, ghOrg string, ghRepo string, token string) error {
	orgID, err := c.GetOrgID(ctx, orgName)
	if err != nil {
		return fmt.Errorf("failed getting org ID: %w", err)
	}
	_, err = c.compute.SetGithubToken(c.withAuth(ctx), &pb.SetGithubTokenRequest{
		OrgId:          orgID,
		GithubOrgName:  ghOrg,
		GithubRepoName: ghRepo,
		GithubToken:    token,
	})
	if err != nil {
		return fmt.Errorf("failed setting Github token: %w", err)
	}
	return nil
}
