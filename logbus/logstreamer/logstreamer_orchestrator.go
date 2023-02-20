package logstreamer

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/earthly/cloud-api/logstream"
	"github.com/hashicorp/go-multierror"
)

type LogstreamOrchestrator struct {
	buildID         string
	bus             LogBus
	c               CloudClient
	initialManifest *logstream.RunManifest

	errMu    sync.Mutex
	errors   []error
	started  atomic.Bool
	retries  atomic.Int32
	streamer *LogStreamer
	deltas   *deltasIter
}

type LOpt func(orchestrator *LogstreamOrchestrator) *LogstreamOrchestrator

func NewLogstreamOrchestrator(bus LogBus, c CloudClient, initialManifest *logstream.RunManifest, opts ...LOpt) *LogstreamOrchestrator {
	ls := &LogstreamOrchestrator{
		bus:     bus,
		c:       c,
		buildID: initialManifest.GetBuildId(),
	}
	ls.retries.Store(10)
	for _, o := range opts {
		ls = o(ls)
	}
	return ls
}

func (l *LogstreamOrchestrator) StartLogstreamer(ctx context.Context) {
	if l.started.Swap(true) {
		// Can only start once
		return
	}
	go func() {
		for l.retries.Add(-1) > 0 {
			l.CloseLastLogstreamer()
			l.deltas = newDeltasIter(DefaultBufferSize, l.initialManifest)
			go l.bus.AddSubscriber(l.deltas)
			l.streamer = New(l.c, l.buildID, l.deltas)
			shouldRetry, err := l.streamer.tryStream(ctx)
			l.addError(err)
			if !shouldRetry {
				return
			}
		}
	}()
}

func (l *LogstreamOrchestrator) CloseLastLogstreamer() {
	if l.deltas != nil {
		l.bus.RemoveSubscriber(l.deltas)
		l.deltas.close()
	}
	if l.streamer != nil {
		go func(streamer *LogStreamer) {
			l.addError(streamer.Close())
		}(l.streamer)
	}
}

func (l *LogstreamOrchestrator) addError(err error) {
	l.errMu.Lock()
	defer l.errMu.Unlock()
	l.errors = append(l.errors, err)
}

func (l *LogstreamOrchestrator) getError() error {
	l.errMu.Lock()
	defer l.errMu.Unlock()
	var retErr error
	for _, err := range l.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}

func (l *LogstreamOrchestrator) Close() error {
	if l.deltas != nil {
		l.bus.RemoveSubscriber(l.deltas)
		l.deltas.close()
	}
	if l.streamer != nil {
		l.addError(l.streamer.Close())
		return l.getError()
	}
	return nil
}
