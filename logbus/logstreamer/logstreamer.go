//go:generate hel --output helheim_mocks_test.go

package logstreamer

import (
	"context"
	"fmt"
	"io"
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

const (
	// DefaultBufferSize is the default size of the buffer in a LogStreamer.
	DefaultBufferSize = 10240

	// DrainTimeout is the time that a LogStreamer.Close() call will wait for
	// any remaining deltas to drain from its buffer.
	DrainTimeout = 60 * time.Second
)

// CloudClient is the type of client that a LogStreamer needs to connect to
// cloud and stream logs.
type CloudClient interface {
	StreamLogs(ctx context.Context, buildID string, deltas cloud.Deltas) error
}

// Opt is an option function, used to adjust optional attributes of a
// LogStreamer during New().
type Opt func(*LogStreamer) *LogStreamer

// WithBuffer ensures that there is a buffer for at least size messages on the
// cloud.Deltas value that the LogStreamer sends to StreamLogs.
func WithBuffer(size int) Opt {
	return func(l *LogStreamer) *LogStreamer {
		l.deltas.chSize = size
		return l
	}
}

// LogStreamer is a log streamer. It uses the cloud client to send
// log deltas to the cloud. It retries on transient errors.
type LogStreamer struct {
	bus     *logbus.Bus
	c       CloudClient
	buildID string
	doneCh  chan struct{}
	errors  []error

	mu        sync.Mutex
	cancelled bool
	deltas    deltasIter
}

// New creates a new LogStreamer.
func New(ctx context.Context, bus *logbus.Bus, c CloudClient, initialManifest *logstream.RunManifest, opts ...Opt) *LogStreamer {
	ls := &LogStreamer{
		bus:     bus,
		c:       c,
		buildID: initialManifest.GetBuildId(),
		doneCh:  make(chan struct{}),
		deltas: deltasIter{
			chSize: DefaultBufferSize,
			init:   initialManifest,
		},
	}
	for _, o := range opts {
		ls = o(ls)
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
	ls.deltas.reset()
	ls.mu.Unlock()
	ls.bus.AddSubscriber(ls)
	defer ls.bus.RemoveSubscriber(ls)
	if err := ls.c.StreamLogs(ctxTry, ls.buildID, &ls.deltas); err != nil {
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

// WriteAsync queues delta up for writing and returns immediately. The returned
// channel will be closed when delta has been passed off to the buffer.
//
// This is different than Write in two important ways:
//
// 1. It will never block.
//
// 2. If the buffer is full, it cannot guarantee order the order that delta will
// be added to the buffer in the order in which WriteAsync was called.
func (ls *LogStreamer) WriteAsync(delta *logstream.Delta) <-chan struct{} {
	done := make(chan struct{})
	if ls.deltas.sendAsync(done, delta) {
		return done
	}

	// TODO (vladaionescu): If these messages show up, we need to rethink
	//						the closing sequence.
	// TODO (vladaionescu): We should only log this if verbose is enabled.
	dt, err := protojson.Marshal(delta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Log streamer closed, but failed to marshal log delta: %v", err)
		return done
	}
	fmt.Fprintf(os.Stderr, "Log streamer closed, dropping delta %v\n", string(dt))
	return done
}

// Write writes the given delta to the log streamer. If there is no room in ls's
// buffer, Write will block until room has been freed.
func (ls *LogStreamer) Write(delta *logstream.Delta) {
	if ls.deltas.send(delta) {
		return
	}

	// TODO (vladaionescu): If these messages show up, we need to rethink
	//						the closing sequence.
	// TODO (vladaionescu): We should only log this if verbose is enabled.
	dt, err := protojson.Marshal(delta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Log streamer closed, but failed to marshal log delta: %v", err)
		return
	}
	fmt.Fprintf(os.Stderr, "Log streamer closed, dropping delta %v\n", string(dt))
}

// Close closes the log streamer.
func (ls *LogStreamer) Close() error {
	ls.mu.Lock()
	ls.deltas.close()
	ls.cancelled = true
	ls.mu.Unlock()

	return ls.drain()
}

func (ls *LogStreamer) drain() error {
	timedOut := false
	select {
	case <-ls.doneCh:
	case <-time.After(DrainTimeout):
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

func decongest(ch <-chan []*logstream.Delta) {
	const decongestTimeout = 100 * time.Millisecond
	t := time.NewTimer(0)
	defer t.Stop()
	for {
		if !t.Stop() {
			<-t.C
		}
		t.Reset(decongestTimeout)
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-t.C:
			return
		}
	}
}

type deltasIter struct {
	mu     sync.Mutex
	chSize int
	ch     chan []*logstream.Delta
	init   *logstream.RunManifest
}

func (d *deltasIter) deltas() chan []*logstream.Delta {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.ch
}

func (d *deltasIter) reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ch != nil {
		go decongest(d.ch)
	}
	d.ch = make(chan []*logstream.Delta, d.chSize)
	d.ch <- []*logstream.Delta{{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{
					ResetAll: d.init,
				},
			},
		},
	}}
}

func (d *deltasIter) sendAsync(sent chan<- struct{}, deltas ...*logstream.Delta) bool {
	ch := d.deltas()
	if ch == nil {
		close(sent)
		return false
	}
	go func() {
		defer close(sent)
		ch <- deltas
	}()
	return true
}

func (d *deltasIter) send(deltas ...*logstream.Delta) bool {
	ch := d.deltas()
	if ch == nil {
		return false
	}
	ch <- deltas
	return true
}

func (d *deltasIter) close() {
	d.mu.Lock()
	defer d.mu.Unlock()
	decongest(d.ch)
	close(d.ch)
	d.ch = nil
}

func (d *deltasIter) Next(ctx context.Context) ([]*logstream.Delta, error) {
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "logstreamer: context closed while waiting on next delta")
	case delta, ok := <-d.deltas():
		if !ok {
			return nil, errors.Wrap(io.EOF, "logstreamer: channel closed")
		}
		return delta, nil
	}
}
