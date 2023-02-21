package logstreamer_test

import (
	"context"
	"io"
	"runtime"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus/logstreamer"
	"github.com/pkg/errors"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/v2"
)

const testTimeout = time.Second

func TestLogstreamer(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
		ctx    context.Context

		mockClient *mockCloudClient
		mockBus    *mockLogBus

		streamer *logstreamer.Orchestrator
	}

	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		expect := expect.New(t)
		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		t.Cleanup(cancel)

		bus := newMockLogBus(t, testTimeout)
		client := newMockCloudClient(t, testTimeout)

		initManifest := &logstream.RunManifest{
			BuildId: "foo",
			Version: 1,
		}

		streamer := logstreamer.NewOrchestrator(bus, client, initManifest)
		streamer.Start(ctx)
		return testCtx{
			t:          t,
			expect:     expect,
			ctx:        ctx,
			mockClient: client,
			mockBus:    bus,
			streamer:   streamer,
		}
	})

	o.Spec("after close, the deltas continually return EOF", func(tt testCtx) {
		var (
			ctx     context.Context
			buildID string
			deltas  cloud.Deltas
		)
		tt.expect(tt.mockClient).To(haveMethodExecuted(
			"StreamLogs",
			within(testTimeout),
			storeArgs(&ctx, &buildID, &deltas),
		))
		_, _ = deltas.Next(tt.ctx) // ignore the initial manifest

		go tt.streamer.Close()

		_, err := deltas.Next(tt.ctx)
		tt.expect(err).To(beErr(io.EOF))

		// If the code doesn't check for a nil channel and tries to read off of
		// it, then it will block until the context is cancelled instead of
		// returning EOF.
		_, err = deltas.Next(tt.ctx)
		tt.expect(err).To(beErr(io.EOF))

		pers.Return(tt.mockClient.StreamLogsOutput, nil)
	})

	o.Spec("it finishes sending deltas", func(tt testCtx) {
		var (
			ctx     context.Context
			buildID string
			deltas  cloud.Deltas
		)
		tt.expect(tt.mockClient).To(haveMethodExecuted(
			"StreamLogs",
			within(testTimeout),
			storeArgs(&ctx, &buildID, &deltas),
		))
		_, _ = deltas.Next(tt.ctx) // ignore the initial manifest

		const toSend = 5
		for i := 0; i < toSend; i++ {
			tt.streamer.WriteToDeltaIter(&logstream.Delta{})
		}
		go tt.streamer.Close()

		runtime.Gosched() // ensure that closeLastLogstreamer() has a chance to close the deltas

		for i := 0; i < toSend; i++ {
			_, err := deltas.Next(tt.ctx)
			tt.expect(err).To(not(haveOccurred()))
		}
		_, err := deltas.Next(tt.ctx)
		tt.expect(err).To(beErr(io.EOF))

		pers.Return(tt.mockClient.StreamLogsOutput, nil)
	})
}

type beErrMatcher struct {
	err error
}

func beErr(err error) beErrMatcher {
	return beErrMatcher{err: err}
}

func (m beErrMatcher) Match(actual any) (any, error) {
	err, ok := actual.(error)
	if !ok {
		return nil, errors.Errorf("expected type %T to implement error", actual)
	}
	if !errors.Is(err, m.err) {
		return nil, errors.Errorf("expected errors.Is([%v], [%v]) to be true", err, m.err)
	}
	return actual, nil
}
