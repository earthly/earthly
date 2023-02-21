package logstreamer

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/earthly/cloud-api/logstream"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

type deltasIter struct {
	mu                   sync.Mutex
	ch                   chan []*logstream.Delta
	closed               bool
	manifestsWritten     atomic.Int32
	formattedLogsWritten atomic.Int32
	verbose              bool
}

func newDeltasIter(bufferSize int, initialManifest *logstream.RunManifest, verbose bool) *deltasIter {
	d := &deltasIter{
		mu:      sync.Mutex{},
		ch:      make(chan []*logstream.Delta, bufferSize),
		verbose: verbose,
	}
	d.ch <- []*logstream.Delta{{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{ResetAll: initialManifest},
			},
		},
	}}
	return d
}

func (d *deltasIter) deltas() (chan []*logstream.Delta, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return d.ch, false
	}
	return d.ch, true
}

func (d *deltasIter) Write(delta *logstream.Delta) {
	ch, ok := d.deltas()
	if !ok {
		//  (vladaionescu): If these messages show up, we need to rethink
		//					the closing sequence.
		if d.verbose {
			dt, err := protojson.Marshal(delta)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Log streamer closed, but failed to marshal log delta: %v", err)
			}
			fmt.Fprintf(os.Stderr, "Log streamer closed, dropping delta %v\n", string(dt))
		}
		return
	}
	if delta.GetDeltaFormattedLog() != nil {
		d.formattedLogsWritten.Add(1)
	} else if delta.GetDeltaManifest() != nil {
		d.manifestsWritten.Add(1)
	}
	ch <- []*logstream.Delta{delta}
}

func (d *deltasIter) close() (int32, int32) {
	d.mu.Lock()
	defer d.mu.Unlock()
	close(d.ch)
	d.closed = true
	return d.manifestsWritten.Load(), d.formattedLogsWritten.Load()
}

func (d *deltasIter) Next(ctx context.Context) ([]*logstream.Delta, error) {
	deltas, _ := d.deltas()
	if deltas == nil {
		return nil, errors.Wrap(io.EOF, "logstreamer: buffer not yet allocated")
	}
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "logstreamer: context closed while waiting on next delta")
	case delta, ok := <-deltas:
		if !ok {
			return nil, errors.Wrap(io.EOF, "logstreamer: channel closed")
		}
		return delta, nil
	}
}
