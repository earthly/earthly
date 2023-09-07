//go:generate hel --output helheim_mocks_test.go

package logstreamer

import (
	"context"
	"sync/atomic"

	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CloudClient is the type of client that a LogStreamer needs to connect to
// cloud and stream logs.
type CloudClient interface {
	StreamLogs(ctx context.Context, buildID string, deltas cloud.Deltas) error
}

// LogStreamer is a log streamer. It uses the cloud client to Write
// log deltas to the cloud. It retries on transient errors.
type LogStreamer struct {
	c       CloudClient
	buildID string

	cancelled atomic.Bool
	deltas    cloud.Deltas
}

// NewLogStreamer creates a new LogStreamer.
func NewLogStreamer(c CloudClient, buildID string, deltas *deltasIter) *LogStreamer {
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
	if ls.cancelled.Load() {
		return false, errors.New("log streamer closed")
	}
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
	ls.cancelled.Store(true)
}
