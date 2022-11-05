package logstreamer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

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
	done            chan struct{}

	mu        sync.Mutex
	ch        chan *logstream.Delta
	batchedCh chan []*logstream.Delta
	errors    []error
}

// New creates a new LogStreamer.
func New(ctx context.Context, bus *logbus.Bus, c cloud.Client, initialManifest *logstream.RunManifest) *LogStreamer {
	ls := &LogStreamer{
		bus:             bus,
		c:               c,
		buildID:         initialManifest.GetBuildId(),
		initialManifest: initialManifest,
		done:            make(chan struct{}),
	}
	go ls.retryLoop(ctx)
	return ls
}

// batcherLoop batches deltas into batches of deltas, based on a maximum
// batch size and a maximum delay interval.
func (ls *LogStreamer) batcherLoop(ctx context.Context) {
	const maxBatchSize = 128
	const maxBatchDuration = 200 * time.Millisecond
	ls.mu.Lock()
	ch := ls.ch
	batchedCh := ls.batchedCh
	ls.mu.Unlock()
	ticker := time.NewTicker(maxBatchDuration)
	defer ticker.Stop()
	batch := make([]*logstream.Delta, 0, maxBatchSize)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(batch) > 0 {
				batchedCh <- batch
				batch = make([]*logstream.Delta, 0, maxBatchSize)
			}
		case delta, ok := <-ch:
			if !ok {
				return
			}
			batch = append(batch, delta)
			if len(batch) >= maxBatchSize {
				batchedCh <- batch
				batch = make([]*logstream.Delta, 0, maxBatchSize)
			}
		}
	}
}

func (ls *LogStreamer) retryLoop(ctx context.Context) {
	defer close(ls.done)
	const maxRetry = 10
	for i := 0; i < maxRetry; i++ {
		err := ls.tryStream(ctx)
		if err == nil {
			return
		}
		if err != errRetry {
			ls.mu.Lock()
			ls.errors = append(ls.errors, err)
			ls.mu.Unlock()
			return
		}
	}
}

var errRetry = errors.New("retry")

func (ls *LogStreamer) tryStream(ctx context.Context) error {
	ctxTry, cancelTry := context.WithCancel(ctx)
	defer cancelTry()
	ls.bus.RemoveSubscriber(ls) // no-op if not added yet
	ls.mu.Lock()
	ls.ch = make(chan *logstream.Delta, 10240)
	ls.batchedCh = make(chan []*logstream.Delta, 10240)
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

	err := ls.c.StreamLogs(ctxTry, ls.buildID, ls.batchedCh)
	if err != nil {
		s, ok := status.FromError(errors.Cause(err))
		if !ok {
			return err
		}
		switch s.Code() {
		case codes.Unavailable, codes.Internal, codes.DeadlineExceeded:
			fmt.Fprintf(os.Stderr, "transient error streaming logs: %v", err)
			return errRetry // will cause a retry.
		default:
			return err
		}
	}
	return nil
}

// Write writes the given delta to the log streamer.
func (ls *LogStreamer) Write(delta *logstream.Delta) {
	ls.mu.Lock()
	ch := ls.ch
	ls.mu.Unlock()
	ch <- delta
}

// Close closes the log streamer.
func (ls *LogStreamer) Close() error {
	close(ls.ch)
	<-ls.done // wait for all messages to be sent (or for log streamer to error-out).
	ls.mu.Lock()
	defer ls.mu.Unlock()
	var retErr error
	for _, err := range ls.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}
