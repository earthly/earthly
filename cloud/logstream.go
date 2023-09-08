package cloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

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
	var finished atomic.Bool
	eg.Go(func() error {
		for {
			resp, err := streamClient.Recv()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[%s] CLIENT FAILED %+v", time.Now(), err)
				return errors.Wrap(err, "failed to read log stream response")
			}
			if resp.GetEofAck() {
				fmt.Fprint(os.Stderr, "GOT EOF ACK\n")
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
			dl, err := deltas.Next(ctx)
			switch {
			case errors.Is(err, io.EOF):
				fmt.Fprint(os.Stderr, "SENDING FINAL EOF\n")
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
			case errors.Is(err, io.ErrNoProgress):
				time.Sleep(time.Millisecond * 250)
				continue
			case err != nil:
				return errors.Wrap(err, "cloud: error getting next delta")
			}

			fmt.Fprintf(os.Stderr, "SENDING %d LOGS TO SERVER\n", len(dl))

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
