package client

import (
	"context"

	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/pkg/errors"
)

func (c *Client) Reserve(ctx context.Context) error {
	_, err := c.ControlClient().Reserve(ctx, &controlapi.ReserveRequest{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
