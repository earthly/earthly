package cloud

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

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
	mu         sync.Mutex
	sent       []*pb.Delta
	attempted  []*pb.Delta
	sendCalls  int
	recv       chan *pb.StreamLogResponse
	done       chan struct{}
	failsAfter int
	streamNops
}

func (f *flakyStream) Send(r *pb.StreamLogRequest) error {
	f.sendCalls++

	f.attempted = append(f.attempted, r.GetDeltas()...)

	if f.failsAfter > -1 && len(f.sent) > f.failsAfter {
		close(f.done)
		return status.Error(codes.Unavailable, "unavailable")
	}

	f.mu.Lock()
	f.sent = append(f.sent, r.GetDeltas()...)
	f.mu.Unlock()

	if r.GetEof() {
		f.recv <- &pb.StreamLogResponse{EofAck: true}
	}

	return nil
}

func (f *flakyStream) Recv() (*pb.StreamLogResponse, error) {
	select {
	case r := <-f.recv:
		return r, nil
	case <-f.done:
		return nil, status.Error(codes.Unknown, "unknown")
	}
}

func newFlakyStream(failsAfter int) *flakyStream {
	return &flakyStream{
		failsAfter: failsAfter,
		recv:       make(chan *pb.StreamLogResponse),
		done:       make(chan struct{}),
	}
}

func TestStreamLogs(t *testing.T) {
	stream := newFlakyStream(-1) // -1 means never fail

	testClient := &testClient{
		stream: func() pb.LogStream_StreamLogsClient {
			return stream
		},
	}

	cl := &Client{
		logstream:        testClient,
		logstreamBackoff: 10 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	man := &pb.RunManifest{
		BuildId: uuid.NewString(),
	}

	ch := make(chan *pb.Delta)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- logDelta("log")
		}
		close(ch)
	}()

	errCh := cl.StreamLogs(ctx, man, ch)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	require.Empty(t, errs)
	require.Equal(t, stream.sendCalls, 12, "expected 10 Sends plus first manifest & EOF (12)")
	require.Len(t, stream.sent, 11, "expected 10 deltas sent and 1 manifest (11)")
}

func TestStreamLogsResume(t *testing.T) {
	streams := []*flakyStream{newFlakyStream(4), newFlakyStream(0), newFlakyStream(-1)}
	idx := 0
	testClient := &testClient{
		stream: func() pb.LogStream_StreamLogsClient {
			defer func() {
				idx++
			}()
			return streams[idx]
		},
	}

	cl := &Client{
		logstream:        testClient,
		logstreamBackoff: 10 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	man := &pb.RunManifest{
		BuildId: uuid.NewString(),
	}

	ch := make(chan *pb.Delta)

	go func() {
		for i := 0; i < 15; i++ {
			ch <- logDelta(fmt.Sprintf("log %d", i))
		}
		close(ch)
	}()

	errCh := cl.StreamLogs(ctx, man, ch)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	require.Len(t, errs, 2)

	// This is the second stream.
	last := streams[len(streams)-1]
	require.Greater(t, last.sendCalls, 1)
	require.NotNil(t, last.sent[0].GetDeltaManifest().GetResume())

	// There should be a duplicate in the "attempted" set as 1 will have failed once.
	counts := map[string]int{}
	for _, stream := range streams {
		for _, delta := range stream.attempted {
			k := string(delta.GetDeltaFormattedLog().GetData())
			if k != "" {
				counts[k]++
			}
		}
	}

	// There ought to exist a log delta that was attempted 3 times.
	var found bool
	for _, attempts := range counts {
		if attempts == 3 {
			found = true
			break
		}
	}
	require.True(t, found)

	// The "sent" stream should not contain any duplicates as the dropped delta
	// will be present in the second stream.
	counts = map[string]int{}
	for _, stream := range streams {
		for _, delta := range stream.sent {
			k := string(delta.GetDeltaFormattedLog().GetData())
			if k != "" {
				counts[k]++
			}
		}
	}
	require.Len(t, counts, 15)

}

func logDelta(message string) *pb.Delta {
	return &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaFormattedLog{
			DeltaFormattedLog: &pb.DeltaFormattedLog{
				TargetId:           "target-1",
				TimestampUnixNanos: 0,
				Data:               []byte(message),
			},
		},
	}
}

func Test_recoverableError(t *testing.T) {
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
			got := recoverableError(c.err)
			if got != c.want {
				t.Errorf("wanted %+v, got %+v", c.want, got)
			}
		})
	}
}

func Test_calcBackoff(t *testing.T) {
	base := 250 * time.Millisecond

	b := calcBackoff(base, []int{0, 0, 0, 0})
	require.Equal(t, int64(4000), b.Milliseconds())

	b = calcBackoff(base, []int{})
	require.Equal(t, int64(250), b.Milliseconds())

	b = calcBackoff(base, []int{0, 51, 0, 0, 0})
	require.Equal(t, int64(2000), b.Milliseconds())

	b = calcBackoff(base, []int{0, 0, 0, 0, 0})
	require.Equal(t, int64(8000), b.Milliseconds())
}
