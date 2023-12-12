package cloud

import (
	"context"

	"github.com/earthly/cloud-api/askv"
	"github.com/pkg/errors"
)

func (c *Client) AutoSkipExists(ctx context.Context, org, project string, hash []byte) (bool, error) {
	response, err := c.askv.Exists(c.withAuth(ctx), &askv.ExistsRequest{
		OrgName:     org,
		ProjectName: project,
		Hash:        hash,
	})
	if err != nil {
		return false, errors.Wrap(err, "auto-skip existence check failed")
	}
	return response.Exists, nil
}

func (c *Client) AutoSkipAdd(ctx context.Context, org, project, path, target string, hash []byte) error {
	_, err := c.askv.Add(c.withAuth(ctx), &askv.AddRequest{
		OrgName:     org,
		ProjectName: project,
		TargetName:  target,
		TargetPath:  path,
		Hash:        hash,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add auto-skip hash")
	}
	return nil
}

func (c *Client) AutoSkipPrune(ctx context.Context, org, project, pathPrefix, target string) error {
	_, err := c.askv.PruneTarget(c.withAuth(ctx), &askv.PruneTargetRequest{
		OrgName:     org,
		ProjectName: project,
		TargetPath:  pathPrefix,
		TargetName:  target,
	})
	if err != nil {
		return errors.Wrap(err, "failed to prune auto-skip data")
	}
	return nil
}
