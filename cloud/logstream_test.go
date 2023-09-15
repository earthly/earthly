package cloud

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type testClient struct {
	stream func() pb.LogStream_StreamLogsClient
}

func (t *testClient) StreamLogs(ctx context.Context, opts ...grpc.CallOption) (pb.LogStream_StreamLogsClient, error) {
	return t.stream(), nil
}

type streamNops struct{}

func (s *streamNops) Header() (metadata.MD, error) { return nil, nil }
func (s *streamNops) Trailer() metadata.MD         { return nil }
func (s *streamNops) CloseSend() error             { return nil }
func (s *streamNops) Context() context.Context     { return nil }
func (s *streamNops) SendMsg(m interface{}) error  { return nil }
func (s *streamNops) RecvMsg(m interface{}) error  { return nil }

type flakyStream struct {
	mu    sync.Mutex
	fail  chan struct{}
	err   error
	calls map[string]int
	sent  []*pb.Delta
	recv  chan *pb.StreamLogResponse
	streamNops
}

func (f *flakyStream) Send(r *pb.StreamLogRequest) error {
	select {
	case <-f.fail:
		return f.err
	default:
	}

	f.mu.Lock()
	f.calls["Send"]++
	f.sent = append(f.sent, r.GetDeltas()...)
	f.mu.Unlock()

	if r.GetEof() {
		f.recv <- &pb.StreamLogResponse{EofAck: true}
	}
	return nil
}

func (f *flakyStream) Recv() (*pb.StreamLogResponse, error) {
	f.mu.Lock()
	f.calls["Recv"]++
	f.mu.Unlock()

	select {
	case r := <-f.recv:
		return r, nil
	case <-f.fail:
		return nil, f.err
	}
}

func newFlakyStream() *flakyStream {
	return &flakyStream{
		calls: map[string]int{},
		recv:  make(chan *pb.StreamLogResponse),
		fail:  make(chan struct{}),
	}
}

func TestStreamLogs(t *testing.T) {
	stream := newFlakyStream()

	testClient := &testClient{
		stream: func() pb.LogStream_StreamLogsClient {
			return stream
		},
	}

	cl := &Client{logstream: testClient}

	ctx := context.Background()

	man := &pb.RunManifest{
		BuildId: uuid.NewString(),
	}

	ch := make(chan *pb.Delta)
	errCh := make(chan error)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- logDelta()
		}
		close(ch)
	}()

	go func() {
		errCh <- cl.StreamLogs(ctx, man, ch, false)
	}()

	err := <-errCh
	require.NoError(t, err)
	require.Equal(t, 12, stream.calls["Send"], "expected 10 Sends plus first manifest & EOF (12)")
	require.Equal(t, 1, stream.calls["Recv"], "expected 1 Recv")
	require.Len(t, stream.sent, 11, "expected 10 deltas sent and 1 manifest (11)")
}

func TestStreamLogsResume(t *testing.T) {
	stream := newFlakyStream()

	r := rand.New(rand.NewSource(1337)) // Deterministic.

	testClient := &testClient{
		stream: func() pb.LogStream_StreamLogsClient {
			return stream
		},
	}

	cl := &Client{logstream: testClient}

	ctx := context.Background()

	man := &pb.RunManifest{
		BuildId: uuid.NewString(),
	}

	ch := make(chan *pb.Delta)
	errCh := make(chan error)

	go func() {
		for i := 0; i < 15; i++ {
			ch <- logDelta()
			if r.Intn(5) == 0 {
				stream.err = status.Error(codes.Unavailable, "unavailable")
				close(stream.fail)
				stream = newFlakyStream()
			}
		}
		close(ch)
	}()

	go func() {
		errCh <- cl.StreamLogs(ctx, man, ch, false)
	}()

	err := <-errCh
	require.NoError(t, err)

	// This is the second stream.
	require.True(t, stream.calls["Send"] > 1)
	require.Equal(t, 1, stream.calls["Recv"], "expected 1 Recv")
	require.NotNil(t, stream.sent[0].GetDeltaManifest().GetResume())
}

func logDelta() *pb.Delta {
	return &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaFormattedLog{
			DeltaFormattedLog: &pb.DeltaFormattedLog{
				TargetId:           "target-1",
				TimestampUnixNanos: 0,
				Data:               []byte("message"),
			},
		},
	}
}

func Test_retryable(t *testing.T) {
	cases := []struct {
		note string
		err  error
		want bool
	}{
		{
			note: "not status error",
			err:  errors.New("fail"),
			want: false,
		},
		{
			note: "unavailable status error",
			err:  status.Error(codes.Unavailable, "unavailable"),
			want: true,
		},
		{
			note: "unknown error",
			err:  status.Error(codes.Unknown, "unknown"),
			want: true,
		},
		{
			note: "wrapped unknown error",
			err:  fmt.Errorf("error: %w", status.Error(codes.Unknown, "unknown")),
			want: true,
		},
		{
			note: "wrapped non-status error",
			err:  fmt.Errorf("error: %w", errors.New("failed")),
			want: false,
		},
		{
			note: "double-wrapped unknown error",
			err:  fmt.Errorf("error: %w", fmt.Errorf("error: %w", status.Error(codes.Unknown, "unknown"))),
			want: true,
		},
	}

	for _, c := range cases {
		t.Run(c.note, func(t *testing.T) {
			got := retryable(c.err)
			if got != c.want {
				t.Errorf("wanted %+v, got %+v", c.want, got)
			}
		})
	}
}
