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

const maxDeltasPerIter = 200

type deltasIter struct {
	mu                   sync.Mutex
	manifestsWritten     atomic.Int32
	formattedLogsWritten atomic.Int32
	verbose              bool
	closed               bool
	ready                atomic.Int32
	initialDelta         *logstream.Delta
	allDeltas            []*logstream.Delta
}

func newDeltasIter(initialManifest *logstream.RunManifest, verbose bool) *deltasIter {
	d := &deltasIter{
		mu:      sync.Mutex{},
		verbose: verbose,
	}

	d.initialDelta = &logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{ResetAll: initialManifest},
			},
		},
	}

	return d
}

func (d *deltasIter) deltas() []*logstream.Delta {

	if d.initialDelta != nil {
		defer func() {
			d.initialDelta = nil
		}()
		return []*logstream.Delta{d.initialDelta}
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var deltas []*logstream.Delta

	// Take up to 100 log entries for sending.
	if len(d.allDeltas) > maxDeltasPerIter {
		deltas = d.allDeltas[0:maxDeltasPerIter]
		d.allDeltas = d.allDeltas[maxDeltasPerIter:len(d.allDeltas)]
	} else {
		deltas = d.allDeltas
		d.allDeltas = []*logstream.Delta{}
	}

	return deltas
}

func (d *deltasIter) Write(delta *logstream.Delta) {
	if d.closed {
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

	d.mu.Lock()
	defer d.mu.Unlock()

	d.allDeltas = append(d.allDeltas, delta)
}

func (d *deltasIter) close() (int32, int32) {
	d.closed = true
	return d.manifestsWritten.Load(), d.formattedLogsWritten.Load()
}

func (d *deltasIter) Next(ctx context.Context) ([]*logstream.Delta, error) {
	deltas := d.deltas()
	if d.closed && len(deltas) == 0 {
		return nil, errors.Wrap(io.EOF, "logstreamer: closed with no remaining deltas")
	}
	if len(deltas) == 0 {
		return nil, io.ErrNoProgress
	}
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "logstreamer: context closed while waiting on next delta")
	default:
		return deltas, nil
	}
}
