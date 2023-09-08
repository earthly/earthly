package logstreamer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/logbus"
	"github.com/hashicorp/go-multierror"
)

// LogBus is a type that LogStreamer subscribes to.
type LogBus interface {
	AddSubscriber(logbus.Subscriber)
	RemoveSubscriber(logbus.Subscriber)
}

type Orchestrator struct {
	buildID         string
	bus             LogBus
	c               CloudClient
	initialManifest *logstream.RunManifest

	mu     sync.Mutex
	prevMU sync.Mutex
	errMu  sync.Mutex

	doneCH chan struct{}
	subCH  chan struct{}

	errors    []error
	startOnce sync.Once
	closed    atomic.Bool
	retries   int
	streamer  *LogStreamer
	deltas    *deltasIter
	verbose   bool
}

type OrchestratorOpt func(*Orchestrator) *Orchestrator

// WithVerbose sets the verbose option on the Orchestrator
func WithVerbose(verbose bool) func(orchestrator *Orchestrator) *Orchestrator {
	return func(orchestrator *Orchestrator) *Orchestrator {
		orchestrator.verbose = verbose
		return orchestrator
	}
}

func NewOrchestrator(bus LogBus, c CloudClient, initialManifest *logstream.RunManifest, opts ...OrchestratorOpt) *Orchestrator {
	ls := &Orchestrator{
		buildID:         initialManifest.GetBuildId(),
		bus:             bus,
		c:               c,
		initialManifest: initialManifest,
		retries:         10,
		doneCH:          make(chan struct{}),
		subCH:           make(chan struct{}),
	}
	for _, o := range opts {
		ls = o(ls)
	}
	return ls
}

// Start will restart streaming to the cloud retrying up the retry count
// Callers should listen to Done to be notified when the streaming contract completes
// Start may only be called once
func (l *Orchestrator) Start(ctx context.Context) {
	l.startOnce.Do(func() {
		go func() {
			defer l.markDone()
			for i := 0; i < l.retries; i++ {
				l.subscribe()
				fmt.Fprintf(os.Stderr, "ATTEMPTING STREAM RETRY %d\n", i)
				shouldRetry, err := l.streamer.Stream(ctx)
				l.addError(err)
				if !shouldRetry || l.closed.Load() {
					return
				}
			}
		}()
	})
}

func (l *Orchestrator) addError(err error) {
	if err == nil {
		return
	}
	l.errMu.Lock()
	defer l.errMu.Unlock()
	l.errors = append(l.errors, err)
}

func (l *Orchestrator) getError() error {
	l.errMu.Lock()
	defer l.errMu.Unlock()
	var retErr error
	for _, err := range l.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}

// Close will mark the deltas and streamer as closed and prevent further writes to the cloud
// Callers should listen to Done() to know when it is safe to exit.
func (l *Orchestrator) Close() (manifestsWritten int32, logsWritten int32, _ error) {
	l.closed.Store(true)
	return l.closePreviousStreamer()
}

func (l *Orchestrator) WriteToDeltaIter(delta *logstream.Delta) {
	l.mu.Lock()
	defer l.mu.Unlock()
	d := l.deltas
	d.Write(delta)
}

// Done returns a channel that is closed once the Logstreamer has finished
// communicating with the server.
// Callers should listen to the Done channel before exiting
func (l *Orchestrator) Done() chan struct{} {
	return l.doneCH
}

func (l *Orchestrator) markDone() {
	close(l.doneCH)
}

func (l *Orchestrator) subscribe() {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, _, _ = l.closePreviousStreamer()
	l.deltas = newDeltasIter(l.initialManifest, l.verbose)
	l.streamer = NewLogStreamer(l.c, l.buildID, l.deltas)
	l.subCH = make(chan struct{})
	go func(subCH chan struct{}) {
		defer close(subCH)
		l.bus.AddSubscriber(l.deltas)
	}(l.subCH)
}

func (l *Orchestrator) closePreviousStreamer() (manifestsWritten int32, logsWritten int32, _ error) {
	l.prevMU.Lock()
	defer l.prevMU.Unlock()

	hasPreviouslySubscribed := l.deltas != nil && l.streamer != nil
	if !hasPreviouslySubscribed {
		return
	}

	<-l.subCH

	l.bus.RemoveSubscriber(l.deltas)
	manifestsWritten, logsWritten = l.deltas.close()
	l.streamer.Close()

	l.deltas = nil
	l.streamer = nil

	return manifestsWritten, logsWritten, l.getError()
}
