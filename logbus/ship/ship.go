package ship

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"time"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNoDeltas = errors.New("no deltas")

type streamer interface {
	StreamLogs(ctx context.Context, buildID string, deltaChan <-chan *pb.Delta) error
}

type LogShipper struct {
	cl     streamer
	man    *pb.RunManifest
	buf    []*pb.Delta
	n      int
	ch     chan *pb.Delta
	done   chan struct{}
	closed atomic.Bool
	retry  bool
	first  bool
	errs   []error
}

func NewLogShipper(cl streamer, man *pb.RunManifest) *LogShipper {
	if man.GetResumeToken() == "" {
		man.ResumeToken = stringutil.RandomAlphanumeric(40)
	}
	return &LogShipper{
		cl:    cl,
		man:   man,
		first: true,
		ch:    make(chan *pb.Delta),
		done:  make(chan struct{}),
	}
}

func (l *LogShipper) Write(delta *pb.Delta) {
	l.buf = append(l.buf, delta)
}

func (l *LogShipper) Start(ctx context.Context) {
	go func() {
		defer func() {
			l.done <- struct{}{}
		}()
		var retry bool
		for {
			select {
			default:
				err := l.attempt(ctx, retry)
				if err != nil {
					l.errs = append(l.errs, err)
					if !retryable(err) {
						return
					}
				} else {
					return
				}
				retry = true
			case <-ctx.Done():
				l.errs = append(l.errs, ctx.Err())
				return
			}
		}
	}()
}

func (l *LogShipper) Close() {
	l.closed.Store(true)
	<-l.done
}

func (l *LogShipper) attempt(ctx context.Context, retry bool) error {

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		l.ch <- l.firstDelta(retry)
		for {
			select {
			default:
				delta, err := l.next()
				if errors.Is(err, io.EOF) {
					close(l.ch)
					return nil
				} else if errors.Is(err, errNoDeltas) {
					time.Sleep(time.Millisecond * 50)
				}
				l.ch <- delta
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	wg.Go(func() error {
		return l.cl.StreamLogs(ctx, l.man.GetBuildId(), l.ch)
	})

	return wg.Wait()
}

func (l *LogShipper) next() (*pb.Delta, error) {
	curLen := len(l.buf)
	if l.closed.Load() && l.n == curLen-1 {
		return nil, io.EOF
	}
	if l.n > curLen-1 {
		return nil, errNoDeltas
	}
	delta := l.buf[l.n]
	l.n++
	return delta, nil
}

func (l *LogShipper) firstDelta(retry bool) *pb.Delta {
	var manifest *pb.DeltaManifest

	if retry {
		manifest = &pb.DeltaManifest{
			DeltaManifestOneof: &pb.DeltaManifest_Resume{
				Resume: &pb.DeltaManifest_ResumeBuild{
					BuildId:     l.man.GetBuildId(),
					Token:       l.man.GetResumeToken(),
					OrgName:     l.man.GetOrgName(),
					ProjectName: l.man.GetProjectName(),
				},
			},
		}
	} else {
		manifest = &pb.DeltaManifest{
			DeltaManifestOneof: &pb.DeltaManifest_ResetAll{
				ResetAll: l.man,
			},
		}
	}

	return &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaManifest{
			DeltaManifest: manifest,
		},
	}
}

func (l *LogShipper) Err() error {
	if len(l.errs) > 0 {
		return &multierror.Error{Errors: l.errs}
	}
	return nil
}

func retryable(err error) bool {
	for {
		if err == nil {
			return false
		}
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unavailable, codes.Unknown:
				return true
			default:
				return false
			}
		}
		err = errors.Unwrap(err)
	}
}
