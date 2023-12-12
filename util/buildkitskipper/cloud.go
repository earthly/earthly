package buildkitskipper

import (
	"context"

	"github.com/earthly/earthly/domain"
)

type ASKVClient interface {
	AutoSkipExists(ctx context.Context, org, project string, hash []byte) (bool, error)
	AutoSkipAdd(ctx context.Context, org, project, path, target string, hash []byte) error
}

func NewCloud(org, project string, target domain.Target, client ASKVClient) (*CloudClient, error) {
	return &CloudClient{
		org:     org,
		project: project,
		path:    target.GetLocalPath(),
		target:  target.GetName(),
		client:  client,
	}, nil
}

type CloudClient struct {
	org     string
	project string
	target  string
	path    string
	client  ASKVClient
}

func (c *CloudClient) Add(ctx context.Context, data []byte) error {
	return c.client.AutoSkipAdd(ctx, c.org, c.project, c.path, c.target, data)
}

func (c *CloudClient) Exists(ctx context.Context, data []byte) (bool, error) {
	return c.client.AutoSkipExists(ctx, c.org, c.project, data)
}
