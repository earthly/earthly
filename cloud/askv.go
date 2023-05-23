package cloud

import (
	"context"

	"github.com/earthly/cloud-api/askv"
	"github.com/pkg/errors"
)

func (c *Client) AutoSkipExists(ctx context.Context, org, project, pipeline string, hash []byte) (bool, error) {
	response, err := c.askv.Exists(c.withAuth(ctx), &askv.ExistsRequest{
		OrgName:      org,
		ProjectName:  project,
		PipelineName: pipeline,
		Hash:         hash,
	})
	if err != nil {
		return false, errors.Wrap(err, "failed querying auto-skip service")
	}
	return response.Exists, nil
}

func (c *Client) AutoSkipAdd(ctx context.Context, org, project, pipeline string, hash []byte) error {
	_, err := c.askv.Add(c.withAuth(ctx), &askv.AddRequest{
		OrgName:      org,
		ProjectName:  project,
		PipelineName: pipeline,
		Hash:         hash,
	})
	if err != nil {
		return errors.Wrap(err, "failed adding to auto-skip service")
	}
	return nil
}
