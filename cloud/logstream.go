package cloud

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/earthly/cloud-api/logstream"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (c *client) StreamLogs(ctx context.Context, buildID string, deltasCh chan []*logstream.Delta) error {
	printDebugLog(fmt.Sprintf("starting StreamLogs"))
	streamClient, err := c.logstream.StreamLogs(c.withAuth(ctx), grpc_retry.Disable())
	if err != nil {
		return errors.Wrap(err, "failed to create log stream client")
	}
	printDebugLog(fmt.Sprintf("created StreamClient"))
	eg, ctx := errgroup.WithContext(ctx)
	var mu sync.Mutex
	finished := false
	printDebugLog(fmt.Sprintf("starting errgroup"))
	eg.Go(func() error {
		for {
			printDebugLog(fmt.Sprintf("group 1: before recv"))
			resp, err := streamClient.Recv()
			if err != nil {
				printDebugLog(fmt.Sprintf("group 1: before recv, got err> %s", err.Error()))
				return errors.Wrap(err, "failed to read log stream response")
			}
			printDebugLog(fmt.Sprintf("group 1: after recv"))
			if resp.GetEofAck() {
				printDebugLog(fmt.Sprintf("group 1: got EOF ack recv"))
				mu.Lock()
				defer mu.Unlock()
				printDebugLog(fmt.Sprintf("group 1: obtained lock"))
				if !finished {
					printDebugLog(fmt.Sprintf("group 1: did not finish, err"))
					return errors.New("unexpected EOF ack")
				}
				err := streamClient.CloseSend()
				if err != nil {
					printDebugLog(fmt.Sprintf("group 1: failed on closeSend> %s", err))
					return errors.Wrap(err, "failed to close log stream")
				}
				printDebugLog(fmt.Sprintf("group 1: returning"))
				return nil
			}
		}
	})
	eg.Go(func() error {
		for {
			printDebugLog(fmt.Sprintf("waiting for delta or ctx done"))
			select {
			case <-ctx.Done():
				printDebugLog(fmt.Sprintf("group 2: got ctx done > %s", ctx.Err()))
				return ctx.Err()
			case deltas, ok := <-deltasCh:
				printDebugLog(fmt.Sprintf("group 2: got delta %v, ok %v", deltas, ok))
				if !ok {
					printDebugLog(fmt.Sprintf("group 2: sending eof"))
					err := streamClient.Send(&logstream.StreamLogRequest{
						BuildId: buildID,
						Eof:     true,
					})
					printDebugLog(fmt.Sprintf("group 2: sent eof"))
					if err != nil {
						printDebugLog(fmt.Sprintf("group 2: failed to send eof> %s", err))
						return errors.Wrap(err, "failed to send EOF to log stream")
					}
					printDebugLog(fmt.Sprintf("group 2: obtaining lock"))
					mu.Lock()
					printDebugLog(fmt.Sprintf("group 2: got lock"))
					finished = true
					mu.Unlock()
					printDebugLog(fmt.Sprintf("group 2: unlocked"))
					return nil
				}
				printDebugLog(fmt.Sprintf("group 2: sending deltas"))
				err := streamClient.Send(&logstream.StreamLogRequest{
					BuildId: buildID,
					Deltas:  deltas,
				})
				if err != nil {
					printDebugLog(fmt.Sprintf("group 2: failed to send deltas > %s", err))
					return errors.Wrap(err, "failed to send log delta")
				}
				printDebugLog(fmt.Sprintf("group 2: sent deltas"))
			}
		}
	})
	printDebugLog(fmt.Sprintf("waiting on errgroup"))
	err = eg.Wait()
	printDebugLog(fmt.Sprintf("errgroup returned"))
	return err
}

func printDebugLog(log string) {
	_, err := fmt.Fprintf(os.Stdout, fmt.Sprintf("DEBUG: ========== %s ========== END DEBUG\n", log))
	if err != nil {
		msg := fmt.Sprintf("failed to log debug log > %s", err)
		panic(msg)
	}
}
