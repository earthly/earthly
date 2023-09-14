package ship

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/hashicorp/go-multierror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxDeltasPerIter = 200

type bus interface {
	AddSubscriber(sub logbus.Subscriber, replay bool)
	RemoveSubscriber(sub logbus.Subscriber)
}

type streamer interface {
	StreamLogs(ctx context.Context, buildID string, deltas cloud.DeltaIterator) error
}

type LogShipper struct {
	man       *pb.RunManifest
	iter      *deltaIter
	bus       bus
	cl        streamer
	done      chan struct{}
	first     bool
	errs      []error
	retryWait time.Duration
}

func NewLogShipper(bus bus, cl *cloud.Client, man *pb.RunManifest) *LogShipper {
	return &LogShipper{
		bus:       bus,
		cl:        cl,
		man:       man,
		done:      make(chan struct{}),
		retryWait: time.Millisecond * 200,
		first:     true,
	}
}

func (l *LogShipper) Start(ctx context.Context) {
	go func() {
		tick := time.NewTicker(l.retryWait)
		defer func() {
			tick.Stop()
			close(l.done)
		}()
		for {
			select {
			case <-tick.C:
				err := l.attempt(ctx)
				if err != nil {
					l.errs = append(l.errs, err)
					if !retryable(err) {
						return
					}
				} else {
					return
				}
			case <-ctx.Done():
				l.errs = append(l.errs, ctx.Err())
				return
			}
		}
	}()
}

func (l *LogShipper) attempt(ctx context.Context) error {
	l.iter = &deltaIter{ctx: ctx}
	if l.first {
		l.first = false
		l.iter.init(l.man)
	} else {
		l.iter.resume(l.man)
	}
	l.bus.AddSubscriber(l.iter, false)
	defer l.bus.RemoveSubscriber(l.iter)
	return l.cl.StreamLogs(ctx, l.man.GetBuildId(), l.iter)
}

func (l *LogShipper) Err() error {
	if len(l.errs) > 0 {
		return &multierror.Error{Errors: l.errs}
	}
	return nil
}

func (l *LogShipper) Close() {
	l.iter.close()
	<-l.done
}

type deltaIter struct {
	ctx     context.Context
	buf     []*pb.Delta
	mu      sync.RWMutex
	started atomic.Bool
	closed  atomic.Bool
}

func (d *deltaIter) init(man *pb.RunManifest) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buf = append(d.buf, &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaManifest{
			DeltaManifest: &pb.DeltaManifest{
				DeltaManifestOneof: &pb.DeltaManifest_ResetAll{
					ResetAll: man,
				},
			},
		},
	})
}

func (d *deltaIter) resume(man *pb.RunManifest) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buf = append(d.buf, &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaManifest{
			DeltaManifest: &pb.DeltaManifest{
				DeltaManifestOneof: &pb.DeltaManifest_Resume{
					Resume: &pb.DeltaManifest_ResumeBuild{
						BuildId:     man.GetBuildId(),
						Token:       man.GetResumeToken(),
						OrgName:     man.GetOrgName(),
						ProjectName: man.GetProjectName(),
					},
				},
			},
		},
	})
}

func (d *deltaIter) close() {
	d.closed.Store(true)
}

func (d *deltaIter) Next(ctx context.Context) ([]*pb.Delta, error) {

	// The first returned value must be a single item slice with the
	// initial reset or resume delta.
	if !d.started.Load() {
		d.started.Store(true)
		d.mu.Lock()
		send := d.buf[0:1]
		d.buf = []*pb.Delta{}
		d.mu.Unlock()
		return send, nil
	}

	d.mu.RLock()
	l := len(d.buf)
	d.mu.RUnlock()

	if l == 0 {
		if d.closed.Load() {
			return nil, io.EOF
		}
		return nil, cloud.ErrNoDeltas
	}

	var (
		ret []*pb.Delta
		buf []*pb.Delta
	)

	if l > maxDeltasPerIter {
		ret = d.buf[0:maxDeltasPerIter]
		buf = d.buf[maxDeltasPerIter:l]
	} else {
		ret = d.buf
		buf = []*pb.Delta{}
	}

	d.mu.Lock()
	d.buf = buf
	d.mu.Unlock()

	return ret, nil
}

func (d *deltaIter) Write(del *pb.Delta) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buf = append(d.buf, del)
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
