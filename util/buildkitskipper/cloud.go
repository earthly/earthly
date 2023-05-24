package buildkitskipper

import (
	"context"
)

type ASKVClient interface {
	AutoSkipExists(ctx context.Context, org, project, pipeline string, hash []byte) (bool, error)
	AutoSkipAdd(ctx context.Context, org, project, pipeline string, hash []byte) error
}

func NewCloud(org, project, pipeline string, client ASKVClient) (*CloudClient, error) {
	return &CloudClient{
		org:      org,
		project:  project,
		pipeline: pipeline,
		c:        client,
	}, nil
}

type CloudClient struct {
	org      string
	project  string
	pipeline string
	c        ASKVClient
}

func (cc *CloudClient) Add(ctx context.Context, data []byte) error {
	return cc.c.AutoSkipAdd(ctx, cc.org, cc.project, cc.pipeline, data)
}

func (cc *CloudClient) Exists(ctx context.Context, data []byte) (bool, error) {
	return cc.c.AutoSkipExists(ctx, cc.org, cc.project, cc.pipeline, data)
}
