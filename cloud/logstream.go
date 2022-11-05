package cloud

import (
	"context"
	"sync"

	"github.com/earthly/cloud-api/logstream"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (c *client) StreamLogs(ctx context.Context, buildID string, deltas chan *logstream.Delta) error {
	streamClient, err := c.logstream.StreamLogs(c.withAuth(ctx))
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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case delta, ok := <-deltas:
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
			// TODO (vladaionescu): Do some batching for the deltas.
			err := streamClient.Send(&logstream.StreamLogRequest{
				BuildId: buildID,
				Deltas:  []*logstream.Delta{delta},
			})
			if err != nil {
				return errors.Wrap(err, "failed to send log delta")
			}
		}
	}
}
