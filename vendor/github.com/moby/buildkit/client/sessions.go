package client

import (
	"context"

	"github.com/pkg/errors"

	controlapi "github.com/moby/buildkit/api/services/control"
)

func (c *Client) SessionHistory(ctx context.Context) ([]*controlapi.SessionHistoryResponse_History, error) {
	res, err := c.ControlClient().SessionHistory(ctx, &controlapi.SessionHistoryRequest{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res.History, nil
}
