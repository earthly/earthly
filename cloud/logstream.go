package cloud

import (
	"context"
	"io"
	"sync/atomic"

	"github.com/earthly/cloud-api/logstream"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Deltas is a type for iterating over logstream deltas.
type Deltas interface {
	Next(ctx context.Context) ([]*logstream.Delta, error)
}

func (c *Client) StreamLogs(ctx context.Context, buildID string, deltas Deltas) error {
	streamClient, err := c.logstream.StreamLogs(c.withAuth(ctx), grpc_retry.Disable())
	if err != nil {
		return errors.Wrap(err, "failed to create log stream client")
	}
	eg, ctx := errgroup.WithContext(ctx)
	var finished int32
	eg.Go(func() error {
		for {
			resp, err := streamClient.Recv()
			if err != nil {
				return errors.Wrap(err, "failed to read log stream response")
			}
			if resp.GetEofAck() {
				if atomic.LoadInt32(&finished) == 0 {
					return errors.New("unexpected EOF ack")
				}
				err := streamClient.CloseSend()
				if err != nil {
					return errors.Wrap(err, "failed to close log stream")
				}
				return nil
			}
		}
	})
	eg.Go(func() error {
		for {
			dl, err := deltas.Next(ctx)
			if errors.Is(err, io.EOF) {
				msg := &logstream.StreamLogRequest{
					BuildId: buildID,
					Eof:     true,
				}
				if err := streamClient.Send(msg); err != nil {
					return errors.Wrap(err, "failed to send EOF to log stream")
				}
				atomic.StoreInt32(&finished, 1)
				return nil
			}
			if err != nil {
				return errors.Wrap(err, "cloud: error getting next delta")
			}

			msg := &logstream.StreamLogRequest{
				BuildId: buildID,
				Deltas:  dl,
			}
			if err := streamClient.Send(msg); err != nil {
				return errors.Wrap(err, "failed to send log delta")
			}
		}
	})
	return eg.Wait()
}
