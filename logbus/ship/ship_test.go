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

func (t *testClient) StreamLogs(ctx context.Context, man *pb.RunManifest, ch <-chan *pb.Delta) error {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return nil
			}
			t.count++
		case <-ctx.Done():
			return ctx.Err()
		}
	}
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

	ctx := context.Background()

	s.Start(ctx)

	n := 0
	for i := 0; i < n; i++ {
		s.Write(logDelta())
	}
	s.Close()

	require.Equal(t, cl.count, n)
	require.NoError(t, s.Err())
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
