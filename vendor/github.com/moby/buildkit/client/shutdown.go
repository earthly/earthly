package client

import (
	"context"

	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/pkg/errors"
)

func (c *Client) ShutdownIfIdle(ctx context.Context) (bool, int, error) {
	res, err := c.ControlClient().ShutdownIfIdle(
		ctx, &controlapi.ShutdownIfIdleRequest{})
	if err != nil {
		return false, 0, errors.Wrap(err, "failed to call shutdown if idle")
	}
	return res.GetWillShutdown(), int(res.GetNumSessions()), nil
}
