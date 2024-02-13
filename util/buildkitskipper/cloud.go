package buildkitskipper

import (
	"context"
	"errors"
	"strings"
)

// ASKVClient provides methods which allow for adding & checking the existence
// of an auto-skip hash. The current implementations are a local BoltDB & a
// remote gRPC service.
type ASKVClient interface {
	AutoSkipExists(ctx context.Context, org string, hash []byte) (bool, error)
	AutoSkipAdd(ctx context.Context, org, path, target string, hash []byte) error
}

// NewCloud creates and returns a new Cloud API implementation.
func NewCloud(client ASKVClient) (*CloudClient, error) {
	return &CloudClient{
		client: client,
	}, nil
}

// CloudClient implements the Cloud API version of the ASKVClient.
type CloudClient struct {
	client ASKVClient
}

// Add a new hash to the Cloud DB.
func (c *CloudClient) Add(ctx context.Context, org, target string, data []byte) error {
	parts := strings.Split(target, "+")
	if len(parts) < 2 {
		return errors.New("invalid target format")
	}
	return c.client.AutoSkipAdd(ctx, org, parts[0], parts[1], data)
}

// Exists checks if an auto-skip hash exists.
func (c *CloudClient) Exists(ctx context.Context, org string, data []byte) (bool, error) {
	return c.client.AutoSkipExists(ctx, org, data)
}
