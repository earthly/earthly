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
	ch              chan *logstream.Delta
	done            chan struct{}

	mu     sync.Mutex
	errors []error
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

func (ls *LogStreamer) retryLoop(ctx context.Context) {
	const maxRetry = 10
	for i := 0; i < maxRetry; i++ {
		ls.bus.RemoveSubscriber(ls) // no-op if not added yet
		ls.ch = make(chan *logstream.Delta, 10240)
		ls.ch <- &logstream.Delta{
			DeltaTypeOneof: &logstream.Delta_DeltaManifest{
				DeltaManifest: &logstream.DeltaManifest{
					DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{
						ResetAll: ls.initialManifest,
					},
				},
			},
		}
		ls.bus.AddSubscriber(ls)

		err := ls.c.StreamLogs(ctx, ls.buildID, ls.ch)
		if err == nil {
			// Success.
			close(ls.done)
			return
		}
		s, ok := status.FromError(errors.Cause(err))
		if !ok {
			ls.mu.Lock()
			ls.errors = append(ls.errors, err)
			ls.mu.Unlock()
			return
		}
		switch s.Code() {
		case codes.Unavailable, codes.Internal, codes.DeadlineExceeded:
			// Retry,
			// TODO (vladaionescu): Log this in a nicer way.
			fmt.Fprintf(os.Stderr, "transient error streaming logs: %v", err)
		default:
			ls.mu.Lock()
			ls.errors = append(ls.errors, err)
			ls.mu.Unlock()
			return
		}
	}
}

// Write writes the given delta to the log streamer.
func (ls *LogStreamer) Write(delta *logstream.Delta) {
	ls.ch <- delta
}

// Close closes the log streamer.
func (ls *LogStreamer) Close() error {
	close(ls.ch)
	<-ls.done
	var retErr error
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for _, err := range ls.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}
