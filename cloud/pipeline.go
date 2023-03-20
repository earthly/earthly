package cloud

import (
	"context"
	pb "github.com/earthly/cloud-api/pipelines"
	"github.com/pkg/errors"
)

type Pipeline struct {
	Name          string
	SatelliteName string
}

func (c *Client) ListPipelines(ctx context.Context, project, org, earthfileHash string) ([]Pipeline, error) {
	resp, err := c.pipelines.ListPipelines(c.withAuth(ctx), &pb.ListPipelinesRequest{
		Project:       project,
		Org:           org,
		EarthfileHash: earthfileHash,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed listing pipelines")
	}

	pipelines := make([]Pipeline, len(resp.Pipelines))
	for i, p := range resp.Pipelines {
		pipelines[i] = Pipeline{
			Name:          p.Name,
			SatelliteName: p.SatelliteName,
		}
	}

	return pipelines, nil
}
