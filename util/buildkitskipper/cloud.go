package buildkitskipper

import (
	"context"
	"strings"

	"github.com/earthly/earthly/domain"
)

type ASKVClient interface {
	AutoSkipExists(ctx context.Context, org string, hash []byte) (bool, error)
	AutoSkipAdd(ctx context.Context, org, path, target string, hash []byte) error
}

func NewCloud(org string, target domain.Target, client ASKVClient) (*CloudClient, error) {
	parts := strings.Split(target.StringCanonical(), "+")
	return &CloudClient{
		org:    org,
		path:   parts[0],
		target: parts[1],
		client: client,
	}, nil
}

type CloudClient struct {
	org    string
	target string
	path   string
	client ASKVClient
}

func (c *CloudClient) Add(ctx context.Context, data []byte) error {
	return c.client.AutoSkipAdd(ctx, c.org, c.path, c.target, data)
}

func (c *CloudClient) Exists(ctx context.Context, data []byte) (bool, error) {
	return c.client.AutoSkipExists(ctx, c.org, data)
}
