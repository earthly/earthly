package ship

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/earthly/cloud-api/logstream"
)

type streamer interface {
	StreamLogs(ctx context.Context, man *pb.RunManifest, ch <-chan *pb.Delta) <-chan error
}

// LogShipper subscribes to the Log Bus & streams log entries up to the remote
// Logstream service. It uses a non-blocking, dynamically resizing buffer to
// reliably stream logs to the server.
type LogShipper struct {
	cl      streamer
	ch      chan *pb.Delta
	man     *pb.RunManifest
	errs    []error
	cancel  context.CancelFunc
	done    chan struct{}
	verbose bool
	mu      sync.Mutex
	closed  atomic.Bool
}

// NewLogShipper creates and returns a new LogShipper.
func NewLogShipper(cl streamer, man *pb.RunManifest, verbose bool) *LogShipper {
	return &LogShipper{
		cl:      cl,
		man:     man,
		ch:      make(chan *pb.Delta),
		done:    make(chan struct{}),
		verbose: verbose,
	}
}

func (l *LogShipper) Write(delta *pb.Delta) {
	if l.closed.Load() {
		fmt.Fprintf(os.Stderr, "WARNING: Message sent to closed log stream: %s\n", delta.String())
		return
	}
	l.ch <- delta
}

// Start the log streaming process and begin writing logs to the server.
func (l *LogShipper) Start(ctx context.Context) {
	go func() {
		ctx, l.cancel = context.WithCancel(context.Background())
		defer l.cancel()
		out := bufferedDeltaChan(ctx, l.ch)
		errCh := l.cl.StreamLogs(ctx, l.man, out)
		for err := range errCh {
			l.mu.Lock()
			l.errs = append(l.errs, err)
			l.mu.Unlock()
		}
		l.done <- struct{}{}
	}()
}

// Close the process and allow for a 10s grace period where lagging messages
// will be drained.
func (l *LogShipper) Close() {
	if l.closed.Load() {
		return
	}
	l.closed.Store(true)
	close(l.ch)
	// Graceful attempt to drain any in-flight logs then force-quit after delay.
	t := time.NewTimer(10 * time.Second)
	defer t.Stop()
	select {
	case <-t.C:
		l.cancel()
	case <-l.done:
		return
	}
}

// Errs returns all errors that were encountered during the process.
func (l *LogShipper) Errs() []error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.errs
}

// bufferedDeltaChan emulates a dynamically resized buffered channel by
// maintaining a slice buffer between two blocking channels. The buffer grows
// and shrinks based on the rate of input and consumption. The code is also
// meant to respect context cancellations.
func bufferedDeltaChan(ctx context.Context, in <-chan *pb.Delta) <-chan *pb.Delta {
	out := make(chan *pb.Delta)
	var buf []*pb.Delta
	go func() {
		defer close(out)
		for {
			// If the buffer is empty, wait for the first item and append it to
			// the buffer. If the input channel is closed here we can safely
			// return as the buffer has been drained (0 items). We also need to
			// respect context cancellations.
			if len(buf) == 0 {
				select {
				case <-ctx.Done():
					return
				case delta, ok := <-in:
					if !ok {
						return
					}
					buf = append(buf, delta)
					continue
				}
			}
			select {
			case <-ctx.Done():
				return
			case delta, ok := <-in:
				if !ok {
					// If input is closed, attempt to drain the buffer while
					// respecting any cancellations.
					for _, delta := range buf {
						select {
						case <-ctx.Done():
							return
						case out <- delta:
						}
					}
					return
				}
				buf = append(buf, delta)
			case out <- buf[0]:
				if len(buf) == 1 {
					buf = nil // This should help with GC.
				} else {
					buf = buf[1:]
				}
			}
		}
	}()
	return out
}
