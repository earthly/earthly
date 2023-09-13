package ship

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/hashicorp/go-multierror"
)

const maxDeltasPerIter = 200

type bus interface {
	AddSubscriber(sub logbus.Subscriber)
}

type streamer interface {
	StreamLogs(ctx context.Context, buildID string, deltas cloud.DeltaIterator) error
}

type LogShipper struct {
	buildID      string
	initManifest *pb.RunManifest
	iter         *deltaIter
	bus          bus
	cl           streamer
	done         chan struct{}
	errs         []error
}

func NewLogShipper(bus bus, cl *cloud.Client, initManifest *pb.RunManifest) *LogShipper {
	return &LogShipper{
		bus:          bus,
		cl:           cl,
		initManifest: initManifest,
		done:         make(chan struct{}),
	}
}

func (l *LogShipper) Start(ctx context.Context) {
	go l.attempt(ctx)
}

func (l *LogShipper) attempt(ctx context.Context) {
	l.iter = &deltaIter{ctx: ctx}
	l.iter.init(l.initManifest)
	l.bus.AddSubscriber(l.iter)
	err := l.cl.StreamLogs(ctx, l.buildID, l.iter)
	if err != nil {
		l.errs = append(l.errs, err)
	}
	close(l.done)
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
	ctx          context.Context
	buf          []*pb.Delta
	mu           sync.RWMutex
	started      atomic.Bool
	closed       atomic.Bool
	initManifest *pb.Delta
}

func (d *deltaIter) init(initManifest *pb.RunManifest) {
	d.initManifest = &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaManifest{
			DeltaManifest: &pb.DeltaManifest{
				DeltaManifestOneof: &pb.DeltaManifest_ResetAll{
					ResetAll: initManifest,
				},
			},
		},
	}
}

func (d *deltaIter) close() {
	d.closed.Store(true)
}

func (d *deltaIter) Next(ctx context.Context) ([]*pb.Delta, error) {

	// The first returned value must be a single item slice with the
	// initial/reset manifest.
	if !d.started.Load() {
		d.started.Store(true)
		return []*pb.Delta{d.initManifest}, nil
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
