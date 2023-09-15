package cloud

import (
	"context"
	"sync/atomic"

	pb "github.com/earthly/cloud-api/logstream"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (c *Client) StreamLogs(ctx context.Context, buildID string, ch <-chan *pb.Delta) error {
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
			select {
			case <-ctx.Done():
				return ctx.Err()
			case delta, ok := <-ch:
				if ok {
					msg := &pb.StreamLogRequest{
						BuildId: buildID,
						Deltas:  []*pb.Delta{delta},
					}
					if err := streamClient.Send(msg); err != nil {
						return errors.Wrap(err, "failed to send log delta")
					}
				} else {
					msg := &pb.StreamLogRequest{
						BuildId: buildID,
						Eof:     true,
					}
					err := streamClient.Send(msg)
					if err != nil {
						return errors.Wrap(err, "failed to send EOF to log stream")
					}
					finished.Store(true)
					return nil
				}
			}

		}
	})

	return eg.Wait()
}
