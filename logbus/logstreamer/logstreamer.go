//go:generate hel --output helheim_mocks_test.go

package logstreamer

import (
	"context"
	"sync"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// DefaultBufferSize is the default size of the buffer in a LogStreamer.
	DefaultBufferSize = 10240
)

// CloudClient is the type of client that a LogStreamer needs to connect to
// cloud and stream logs.
type CloudClient interface {
	StreamLogs(ctx context.Context, buildID string, deltas cloud.Deltas) error
}

// LogBus is a type that LogStreamer subscribes to.
type LogBus interface {
	AddSubscriber(logbus.Subscriber)
	RemoveSubscriber(logbus.Subscriber)
}

// LogStreamer is a log streamer. It uses the cloud client to Write
// log deltas to the cloud. It retries on transient errors.
type LogStreamer struct {
	c       CloudClient
	buildID string
	errors  []error

	mu        sync.Mutex
	cancelled bool
	deltas    cloud.Deltas
}

// New creates a new LogStreamer.
func New(c CloudClient, buildID string, deltas *deltasIter) *LogStreamer {
	ls := &LogStreamer{
		c:       c,
		buildID: buildID,
		deltas:  deltas,
	}
	return ls
}

// Stream will attempt to stream deltas from the deltasIter to the Cloud
// This also returns a bool indicating whether the error that occurred should be considered retryable
func (ls *LogStreamer) Stream(ctx context.Context) (bool, error) {
	ctxTry, cancelTry := context.WithCancel(ctx)
	defer cancelTry()
	ls.mu.Lock()
	if ls.cancelled {
		ls.mu.Unlock()
		return false, errors.New("log streamer closed")
	}
	ls.mu.Unlock()
	if err := ls.c.StreamLogs(ctxTry, ls.buildID, ls.deltas); err != nil {
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

// Close closes the log streamer.
func (ls *LogStreamer) Close() {
	ls.mu.Lock()
	ls.cancelled = true
	ls.mu.Unlock()
}
