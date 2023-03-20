package cloud

import (
	"context"
	pb "github.com/earthly/cloud-api/pipelines"
	"github.com/pkg/errors"
)

type PipelineTrigger struct {
	Type     string
	Modifier string
}

type PipelineArg struct {
	Name         string
	DefaultValue string
}

type Pipeline struct {
	Repo          string
	Path          string
	Name          string
	Org           string
	Triggers      []*PipelineTrigger
	Args          []*PipelineArg
	RepoId        string
	Project       string
	IsPush        bool
	Id            string
	PathHash      string
	ProviderOrg   string
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
		triggers := make([]*PipelineTrigger, len(p.Triggers))
		for i, t := range p.Triggers {
			triggers[i] = &PipelineTrigger{
				Type:     t.Type.String(),
				Modifier: t.Modifier,
			}
		}

		args := make([]*PipelineArg, len(p.Triggers))
		for i, a := range p.Args {
			args[i] = &PipelineArg{
				Name:         a.Name,
				DefaultValue: a.DefaultValue,
			}
		}

		pipelines[i] = Pipeline{
			Repo:          p.Repo,
			Path:          p.Path,
			Name:          p.Name,
			Org:           p.Org,
			Triggers:      triggers,
			Args:          args,
			RepoId:        p.RepoId,
			Project:       p.Project,
			IsPush:        p.IsPush,
			Id:            p.Id,
			PathHash:      p.PathHash,
			ProviderOrg:   p.ProviderOrg,
			SatelliteName: p.SatelliteName,
		}
	}

	return pipelines, nil
}
