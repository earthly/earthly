package ship

import (
	"context"
	"time"

	pb "github.com/earthly/cloud-api/logstream"
)

type streamer interface {
	StreamLogs(ctx context.Context, man *pb.RunManifest, ch <-chan *pb.Delta, verbose bool) error
}

type LogShipper struct {
	cl      streamer
	ch      chan *pb.Delta
	man     *pb.RunManifest
	err     error
	cancel  context.CancelFunc
	done    chan struct{}
	verbose bool
}

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
	l.ch <- delta
}

func (l *LogShipper) Start(ctx context.Context) {
	go func() {
		ctx, l.cancel = context.WithCancel(ctx)
		defer l.cancel()
		out := bufferedDeltaChan(ctx, l.ch)
		l.err = l.cl.StreamLogs(ctx, l.man, out, l.verbose)
		l.done <- struct{}{}
	}()
}

func (l *LogShipper) Close() {
	close(l.ch)
	// Graceful clean-up then hard-close after 30s.
	t := time.NewTimer(30 * time.Second)
	select {
	case <-t.C:
		l.cancel()
	case <-l.done:
		return
	}
}

func (l *LogShipper) Err() error {
	return l.err
}

// bufferedDeltaChan emulates a dynamically-resizing buffered channel by
// maintaining a buffer between two blocking channels. The code is also meant to
// respect context cancellations.
func bufferedDeltaChan(ctx context.Context, in <-chan *pb.Delta) <-chan *pb.Delta {
	out := make(chan *pb.Delta)
	var buf []*pb.Delta
	go func() {
		defer close(out)
		for {
			// If the buffer is empty, wait for the first item and append it to
			// the buffer. If the input channel is closed here we can safely
			// return as the buffer has been drained.
			if len(buf) == 0 {
				delta, ok := <-in
				if !ok {
					return
				}
				buf = append(buf, delta)
				continue
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
