package ship

import (
	"context"
	"testing"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testClient struct {
	count int
}

func (t *testClient) StreamLogs(ctx context.Context, man *pb.RunManifest, ch <-chan *pb.Delta) <-chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
				t.count++
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()
	return errCh
}

func TestLogShipper(t *testing.T) {
	cl := &testClient{}
	buildID := uuid.NewString()

	man := &pb.RunManifest{
		BuildId: buildID,
	}

	s := &LogShipper{
		cl:   cl,
		man:  man,
		ch:   make(chan *pb.Delta),
		done: make(chan struct{}),
	}

	s.Start()

	n := 50
	for i := 0; i < n; i++ {
		s.Write(logDelta())
	}
	s.Close()

	require.Equal(t, cl.count, n)
	require.Empty(t, s.Errs())
}

func Test_bufferedDeltaChan(t *testing.T) {
	in := make(chan *pb.Delta)
	ctx := context.Background()
	out := bufferedDeltaChan(ctx, in)

	n := 50
	go func() {
		for i := 0; i < n; i++ {
			in <- logDelta()
		}
		close(in)
	}()

	got := []*pb.Delta{}
	for d := range out {
		got = append(got, d)
	}

	require.Len(t, got, n)
	_, ok := <-out
	require.False(t, ok, "expected output chan to be closed")
}

func Test_bufferedDeltaChan_cancel(t *testing.T) {
	in := make(chan *pb.Delta)
	ctx, cancel := context.WithCancel(context.Background())
	out := bufferedDeltaChan(ctx, in)

	in <- logDelta()

	// A cancel here will close the output channel & exit the internal loop.
	cancel()

	_, ok := <-out
	require.False(t, ok, "expected output chan to be closed")
}

func Test_bufferedDeltaChan_drain(t *testing.T) {
	in := make(chan *pb.Delta)
	ctx := context.Background()
	out := bufferedDeltaChan(ctx, in)

	// Run in same thread to ensure we buffer.
	n := 50
	for i := 0; i < n; i++ {
		in <- logDelta()
	}
	close(in)

	// Consume from buffer after input channel is closed.
	got := []*pb.Delta{}
	for d := range out {
		got = append(got, d)
	}

	require.Len(t, got, n)
	_, ok := <-out
	require.False(t, ok, "expected output chan to be closed")
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
