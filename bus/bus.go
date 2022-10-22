package bus

import (
	"context"

	"github.com/earthly/cloud-api/logstream"
)

// Bus is a build log data bus.
// It listens for events from BuildKit and forwards them to the console and
// possibly also to a remote logstream service.
type Bus struct {
	ch chan *logstream.Delta
	bp *BuildPrinter
	sm *SolverMonitor
}

// New creates a new Bus.
func New(ctx context.Context) *Bus {
	b := &Bus{
		ch: make(chan *logstream.Delta, 10240),
		bp: nil, // set below
		sm: nil, // set below
	}
	b.bp = newBuildPrinter(b)
	b.sm = newSolverMonitor(b)
	go func() {
		<-ctx.Done()
		b.Close()
	}()
	go b.messageLoop(ctx)
	return b
}

// GetPrinter returns the underlying BuildPrinter.
func (b *Bus) GetPrinter() *BuildPrinter {
	return b.bp
}

// GetSolverMonitor returns the underlying SolverMonitor.
func (b *Bus) GetSolverMonitor() *SolverMonitor {
	return b.sm
}

// RawDelta sends a raw delta on the bus.
func (b *Bus) RawDelta(delta *logstream.Delta) {
	b.ch <- delta
}

func (b *Bus) messageLoop(ctx context.Context) {
	for delta := range b.ch {
		b.handleDelta(ctx, delta)
	}
}

func (b *Bus) handleDelta(ctx context.Context, delta *logstream.Delta) {
	// TODO
}

// Close closes the bus.
func (b *Bus) Close() {
	close(b.ch)
}
