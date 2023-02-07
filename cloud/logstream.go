package cloud

import (
	"context"
	"sync"

	"github.com/earthly/cloud-api/logstream"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (c *Client) StreamLogs(ctx context.Context, buildID string, deltasCh <-chan []*logstream.Delta) error {
	streamClient, err := c.logstream.StreamLogs(c.withAuth(ctx), grpc_retry.Disable())
	if err != nil {
		return errors.Wrap(err, "failed to create log stream client")
	}
	eg, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	finished := false
	eg.Go(func() error {
		for {
			resp, err := streamClient.Recv()
			if err != nil {
				return errors.Wrap(err, "failed to read log stream response")
			}
			if resp.GetEofAck() {
				mu.Lock()
				defer mu.Unlock()
				if !finished {
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
			case deltas, ok := <-deltasCh:
				if !ok {
					err := streamClient.Send(&logstream.StreamLogRequest{
						BuildId: buildID,
						Eof:     true,
					})
					if err != nil {
						return errors.Wrap(err, "failed to send EOF to log stream")
					}
					mu.Lock()
					finished = true
					mu.Unlock()
					return nil
				}
				err := streamClient.Send(&logstream.StreamLogRequest{
					BuildId: buildID,
					Deltas:  deltas,
				})
				if err != nil {
					return errors.Wrap(err, "failed to send log delta")
				}
			}
		}
	})
	return eg.Wait()
}
