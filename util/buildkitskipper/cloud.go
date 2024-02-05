package buildkitskipper

import (
	"context"
	"errors"
	"strings"
)

type ASKVClient interface {
	AutoSkipExists(ctx context.Context, org string, hash []byte) (bool, error)
	AutoSkipAdd(ctx context.Context, org, path, target string, hash []byte) error
}

func NewCloud(client ASKVClient) (*CloudClient, error) {
	return &CloudClient{
		client: client,
	}, nil
}

type CloudClient struct {
	org    string
	client ASKVClient
}

func (c *CloudClient) Add(ctx context.Context, org, target string, data []byte) error {
	parts := strings.Split(target, "+")
	if len(parts) < 2 {
		return errors.New("invalid target format")
	}
	return c.client.AutoSkipAdd(ctx, org, parts[0], parts[1], data)
}

func (c *CloudClient) Exists(ctx context.Context, org string, data []byte) (bool, error) {
	return c.client.AutoSkipExists(ctx, org, data)
}
