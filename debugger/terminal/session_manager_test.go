package terminal

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/stretchr/testify/require"
)

type blockingTermConnector struct {
	started     chan struct{}
	finish      chan struct{}
	releaseOnce sync.Once
	calls       atomic.Int32
	err         error
}

func newBlockingTermConnector(err error) *blockingTermConnector {
	return &blockingTermConnector{
		started: make(chan struct{}),
		finish:  make(chan struct{}),
		err:     err,
	}
}

func newTestSessionManager(connect connectTermFunc) *SessionManager {
	return &SessionManager{connect: connect}
}

func (c *blockingTermConnector) connect(context.Context, io.ReadWriteCloser, conslogging.ConsoleLogger) error {
	c.calls.Add(1)
	close(c.started)
	<-c.finish
	return c.err
}

func (c *blockingTermConnector) unblock() {
	c.releaseOnce.Do(func() {
		close(c.finish)
	})
}

func connectAsync(manager *SessionManager, ctx context.Context) <-chan error {
	result := make(chan error, 1)
	go func() {
		result <- manager.ConnectTerm(ctx, nil, conslogging.ConsoleLogger{})
	}()
	return result
}

func requireBlocked(t *testing.T, result <-chan error) {
	t.Helper()

	select {
	case err := <-result:
		t.Fatalf("debugger connection returned while the first session was active: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestSessionManagerConnectsOnlyFirstSession(t *testing.T) {
	connector := newBlockingTermConnector(nil)
	t.Cleanup(connector.unblock)
	manager := newTestSessionManager(connector.connect)

	firstResult := connectAsync(manager, context.Background())
	<-connector.started

	const waiterCount = 3
	waiterResults := make(chan error, waiterCount)
	for i := 0; i < waiterCount; i++ {
		go func() {
			waiterResults <- manager.ConnectTerm(context.Background(), nil, conslogging.ConsoleLogger{})
		}()
	}

	requireBlocked(t, waiterResults)
	connector.unblock()

	require.NoError(t, <-firstResult)
	for i := 0; i < waiterCount; i++ {
		require.NoError(t, <-waiterResults)
	}
	require.Equal(t, int32(1), connector.calls.Load())
}

func TestSessionManagerStopsWaitingWhenContextIsCanceled(t *testing.T) {
	connector := newBlockingTermConnector(nil)
	t.Cleanup(connector.unblock)
	manager := newTestSessionManager(connector.connect)

	firstResult := connectAsync(manager, context.Background())
	<-connector.started

	ctx, cancel := context.WithCancel(context.Background())
	waiterResult := connectAsync(manager, ctx)
	requireBlocked(t, waiterResult)
	cancel()

	require.ErrorIs(t, <-waiterResult, context.Canceled)
	require.Equal(t, int32(1), connector.calls.Load())

	connector.unblock()
	require.NoError(t, <-firstResult)
}

func TestSessionManagerCanceledConnectionDoesNotInitializeSession(t *testing.T) {
	var callCount atomic.Int32
	manager := newTestSessionManager(func(context.Context, io.ReadWriteCloser, conslogging.ConsoleLogger) error {
		callCount.Add(1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, manager.ConnectTerm(ctx, nil, conslogging.ConsoleLogger{}), context.Canceled)
	require.Zero(t, callCount.Load())
	require.Nil(t, manager.sessionDone)

	require.NoError(t, manager.ConnectTerm(context.Background(), nil, conslogging.ConsoleLogger{}))
	require.Equal(t, int32(1), callCount.Load())
	require.NotNil(t, manager.sessionDone)
}

func TestSessionManagerReturnsFirstConnectionError(t *testing.T) {
	wantErr := errors.New("connect terminal")
	connector := newBlockingTermConnector(wantErr)
	manager := newTestSessionManager(connector.connect)

	firstResult := connectAsync(manager, context.Background())
	<-connector.started
	connector.unblock()

	require.ErrorIs(t, <-firstResult, wantErr)
	require.NoError(t, manager.ConnectTerm(context.Background(), nil, conslogging.ConsoleLogger{}))
	require.Equal(t, int32(1), connector.calls.Load())
}
