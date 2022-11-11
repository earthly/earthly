package logstreamer

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogStreamer is a log streamer. It uses the cloud client to send
// log deltas to the cloud. It retries on transient errors.
type LogStreamer struct {
	bus             *logbus.Bus
	c               cloud.Client
	buildID         string
	initialManifest *logstream.RunManifest
	doneCh          chan struct{}
	errors          []error

	mu        sync.Mutex
	cancelled bool
	ch        chan *logstream.Delta
	batchedCh chan []*logstream.Delta
}

// New creates a new LogStreamer.
func New(ctx context.Context, bus *logbus.Bus, c cloud.Client, initialManifest *logstream.RunManifest) *LogStreamer {
	ls := &LogStreamer{
		bus:             bus,
		c:               c,
		buildID:         initialManifest.GetBuildId(),
		initialManifest: initialManifest,
		doneCh:          make(chan struct{}),
	}
	go ls.retryLoop(ctx)
	return ls
}

const maxBatchSize = 128

func (ls *LogStreamer) retryLoop(ctx context.Context) {
	defer close(ls.doneCh)
	const maxRetry = 10
	for i := 0; i < maxRetry; i++ {
		retry, err := ls.tryStream(ctx)
		if err == nil {
			return
		}
		if i == maxRetry-1 {
			retry = false
		}
		if !retry {
			ls.errors = append(ls.errors, err)
			return
		}
		fmt.Fprintf(os.Stderr, "transient error streaming logs: %v\n", err)
	}
}

func (ls *LogStreamer) tryStream(ctx context.Context) (bool, error) {
	ctxTry, cancelTry := context.WithCancel(ctx)
	defer cancelTry()
	ls.mu.Lock()
	if ls.cancelled {
		// TODO (vladaionescu): It would be nice if on cancellation we could
		// 						still go through the entire retry loop.
		//						This would require that we close ls.ch on each
		//						retry somehow safely.
		ls.mu.Unlock()
		return false, errors.New("log streamer closed")
	}
	const chSize = 128
	ls.ch = make(chan *logstream.Delta, chSize)
	ls.batchedCh = make(chan []*logstream.Delta, chSize*maxBatchSize)
	ls.mu.Unlock()
	ls.batchedCh <- []*logstream.Delta{
		{
			DeltaTypeOneof: &logstream.Delta_DeltaManifest{
				DeltaManifest: &logstream.DeltaManifest{
					DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{
						ResetAll: ls.initialManifest,
					},
				},
			},
		},
	}
	go ls.batcherLoop(ctxTry)
	ls.bus.AddSubscriber(ls)
	defer ls.bus.RemoveSubscriber(ls)
	err := ls.c.StreamLogs(ctxTry, ls.buildID, ls.batchedCh)
	if err != nil {
		s, ok := status.FromError(errors.Cause(err))
		if !ok {
			return false, err
		}
		switch s.Code() {
		case codes.Unavailable, codes.Internal, codes.DeadlineExceeded:
			return true, err
		default:
			return false, err
		}
	}
	return false, nil
}

// batcherLoop batches deltas into batches of deltas, based on a maximum
// batch size and a maximum delay interval.
func (ls *LogStreamer) batcherLoop(ctx context.Context) {
	ls.mu.Lock()
	ch := ls.ch
	batchedCh := ls.batchedCh
	ls.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return
		case delta, ok := <-ch:
			if !ok {
				close(batchedCh)
				return
			}
			batchedCh <- []*logstream.Delta{delta}
		}
	}

	// TODO @#
	// ticker := time.NewTicker(maxBatchDuration)
	// defer ticker.Stop()
	// batch := make([]*logstream.Delta, 0, maxBatchSize)
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	case <-ticker.C:
	// 		if len(batch) > 0 {
	// 			batchedCh <- batch
	// 			batch = make([]*logstream.Delta, 0, maxBatchSize)
	// 		}
	// 	case delta, ok := <-ch:
	// 		if !ok {
	// 			if len(batch) > 0 {
	// 				batchedCh <- batch
	// 			}
	// 			close(batchedCh)
	// 			return
	// 		}
	// 		batch = append(batch, delta)
	// 		if len(batch) >= maxBatchSize {
	// 			batchedCh <- batch
	// 			batch = make([]*logstream.Delta, 0, maxBatchSize)
	// 		}
	// 	}
	// }
}

// Write writes the given delta to the log streamer.
func (ls *LogStreamer) Write(delta *logstream.Delta) {
	ls.mu.Lock()
	ch := ls.ch // ls.ch may get replaced on retry
	ls.mu.Unlock()
	ch <- delta
}

// Close closes the log streamer.
func (ls *LogStreamer) Close() error {
	ls.mu.Lock()
	if ls.ch != nil {
		close(ls.ch)
	}
	ls.cancelled = true
	ls.mu.Unlock()
	<-ls.doneCh // wait for all messages to be sent (or for log streamer to error-out).
	var retErr error
	for _, err := range ls.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}
