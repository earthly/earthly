package cloud

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/earthly/cloud-api/logstream"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var ErrNoDeltas = errors.New("no deltas")

// DeltaIterator is a type for iterating over logstream deltas.
type DeltaIterator interface {
	Next(ctx context.Context) ([]*logstream.Delta, error)
}

func (c *Client) StreamLogs(ctx context.Context, buildID string, iter DeltaIterator) error {
	streamClient, err := c.logstream.StreamLogs(c.withAuth(ctx), grpc_retry.Disable())
	if err != nil {
		return errors.Wrap(err, "failed to create log stream client")
	}
	eg, ctx := errgroup.WithContext(ctx)
	var finished atomic.Bool
	eg.Go(func() error {
		for {
			resp, err := streamClient.Recv()
			if err != nil {
				return errors.Wrap(err, "failed to read log stream response")
			}
			if resp.GetEofAck() {
				if !finished.Load() {
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
			dl, err := iter.Next(ctx)
			switch {
			case errors.Is(err, io.EOF):
				msg := &logstream.StreamLogRequest{
					BuildId: buildID,
					Eof:     true,
				}
				err := streamClient.Send(msg)
				if err != nil {
					return errors.Wrap(err, "failed to send EOF to log stream")
				}
				finished.Store(true)
				return nil
			case errors.Is(err, ErrNoDeltas):
				// Not yet finished, but no logs written since last
				// Next. Back-off for a bit.
				time.Sleep(250 * time.Millisecond)
				// Pass through and send a request with no deltas (heartbeat).
			case err != nil:
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
