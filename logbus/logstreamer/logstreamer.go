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
	"google.golang.org/protobuf/encoding/protojson"
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
	ch        chan []*logstream.Delta
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
	const chSize = 10240
	if ls.ch != nil {
		// In case the channel is congested, this frees up any stuck writers.
		prevCh := ls.ch
		go func() {
			for range prevCh {
			}
		}()
	}
	ls.ch = make(chan []*logstream.Delta, chSize)
	ls.mu.Unlock()
	ls.ch <- []*logstream.Delta{
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
	ls.bus.AddSubscriber(ls)
	defer ls.bus.RemoveSubscriber(ls)
	err := ls.c.StreamLogs(ctxTry, ls.buildID, ls.ch)
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

// Write writes the given delta to the log streamer.
func (ls *LogStreamer) Write(delta *logstream.Delta) {
	ls.mu.Lock()
	ch := ls.ch // ls.ch may get replaced on retry
	ls.mu.Unlock()
	if ch != nil {
		ch <- []*logstream.Delta{delta}
	} else {
		// TODO (vladaionescu): If these messages show up, we need to rethink
		//						the closing sequence.
		// TODO (vladaionescu): We should only log this if verbose is enabled.
		dt, err := protojson.Marshal(delta)
		if err != nil {
			fmt.Fprintf(os.Stderr, "log streamer closed, but failed to marshal log delta: %v", err)
			return
		}
		fmt.Fprintf(os.Stderr, "log streamer closed, dropping delta %v\n", string(dt))
	}
}

// Close closes the log streamer.
func (ls *LogStreamer) Close() error {
	ls.mu.Lock()
	if ls.ch != nil {
		close(ls.ch)
		ls.ch = nil
	}
	ls.cancelled = true
	ls.mu.Unlock()
	// wait for all messages to be sent
	timedOut := false
	select {
	case <-ls.doneCh:
	case <-time.After(60 * time.Second):
		timedOut = true
	}
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if timedOut {
		ls.errors = append(ls.errors, errors.New("timed out waiting for log streamer to close"))
	}
	var retErr error
	for _, err := range ls.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}
