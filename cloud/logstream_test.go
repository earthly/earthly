package cloud_test

import (
	"context"
	"io"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/v2"
	"google.golang.org/protobuf/proto"
)

const testTimeout = time.Second

func TestStreamLogs(topT *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation

		mockConn   *mockGRPCConn
		mockDeltas *mockDeltas

		cloudClient *cloud.Client
	}

	o := onpar.BeforeEach(onpar.New(topT), func(t *testing.T) testCtx {
		expect := expect.New(t)
		mockConn := newMockGRPCConn(t, testTimeout)
		client, err := cloud.NewClientOpts(context.Background(), cloud.WithGRPCConn(mockConn))
		expect(err).To(not(haveOccurred()))
		return testCtx{
			t:          t,
			expect:     expect,
			mockConn:   mockConn,
			mockDeltas: newMockDeltas(t, testTimeout),

			cloudClient: client,
		}
	})

	o.Spec("it errors on unexpected EOF", func(tt testCtx) {
		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		errs := make(chan error)
		go func() {
			errs <- tt.cloudClient.StreamLogs(ctx, "foo", tt.mockDeltas)
		}()

		stream := newMockClientStream(tt.t, testTimeout)
		tt.expect(tt.mockConn).To(haveMethodExecuted(
			"NewStream",
			within(testTimeout),
			withArgs(pers.Any, pers.Any, "/api.public.logstream.LogStream/StreamLogs", pers.VariadicAny),
			returning(stream, nil),
		))

		var respInter any
		tt.expect(stream).To(haveMethodExecuted(
			"RecvMsg",
			within(testTimeout),
			storeArgs(&respInter),
		))

		resp := respInter.(*logstream.StreamLogResponse)
		resp.EofAck = true
		pers.Return(stream.RecvMsgOutput, nil)

		// We also need the sending goroutine to exit for us to get the error,
		// but we need to give the reading goroutine time to blow up first,
		// since that's the error we expect.
		time.Sleep(100 * time.Millisecond)

		pers.Return(tt.mockDeltas.NextOutput, nil, errors.New("boom"))

		select {
		case err := <-errs:
			tt.expect(err).To(haveOccurred())
			tt.expect(err.Error()).To(equal("unexpected EOF ack"))
		case <-time.After(testTimeout):
			tt.t.Fatalf("timed out waiting for StreamLogs to exit")
		}
	})

	o.Spec("it forwards deltas", func(tt testCtx) {
		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		errs := make(chan error)
		go func() {
			errs <- tt.cloudClient.StreamLogs(ctx, "foo", tt.mockDeltas)
		}()

		stream := newMockClientStream(tt.t, testTimeout)

		tt.expect(tt.mockConn).To(haveMethodExecuted(
			"NewStream",
			within(testTimeout),
			withArgs(pers.Any, pers.Any, "/api.public.logstream.LogStream/StreamLogs", pers.VariadicAny),
			returning(stream, nil),
		))

		delta := &logstream.Delta{
			Version: 1,
			DeltaTypeOneof: &logstream.Delta_DeltaLog{
				DeltaLog: &logstream.DeltaLog{
					TargetId:  "foo",
					CommandId: "bar",
				},
			},
		}
		deltas := []*logstream.Delta{delta}
		pers.Return(tt.mockDeltas.NextOutput, deltas, nil)

		tt.expect(stream).To(haveMethodExecuted(
			"SendMsg",
			within(testTimeout),
			withArgs(protoEqual(&logstream.StreamLogRequest{
				BuildId: "foo",
				Deltas:  deltas,
				Eof:     false,
			})),
			returning(nil),
		))

		tt.t.Cleanup(func() {
			pers.Return(tt.mockDeltas.NextOutput, nil, io.EOF)
			tt.expect(stream).To(haveMethodExecuted(
				"SendMsg",
				within(testTimeout),
				withArgs(protoEqual(&logstream.StreamLogRequest{
					BuildId: "foo",
					Eof:     true,
				})),
				returning(nil),
			))

			var respInter any
			tt.expect(stream).To(haveMethodExecuted(
				"RecvMsg",
				within(testTimeout),
				storeArgs(&respInter),
			))
			resp := respInter.(*logstream.StreamLogResponse)
			resp.EofAck = true
			pers.Return(stream.RecvMsgOutput, nil)
			pers.Return(stream.CloseSendOutput, nil)
			select {
			case err := <-errs:
				tt.expect(err).To(not(haveOccurred()))
			case <-time.After(testTimeout):
				tt.t.Fatalf("timed out waiting for StreamLogs to exit")
			}
		})
	})
}

type protoMatcher struct {
	expected proto.Message
}

func protoEqual(m proto.Message) *protoMatcher {
	return &protoMatcher{expected: m}
}

func (m *protoMatcher) Match(actual any) (any, error) {
	msg, ok := actual.(proto.Message)
	if !ok {
		return nil, errors.Errorf("expected type %T to implement proto.Message, but it does not", actual)
	}
	if !proto.Equal(msg, m.expected) {
		return nil, errors.Errorf("expected message [%v] to equal [%v]", msg, m.expected)
	}
	return actual, nil
}
