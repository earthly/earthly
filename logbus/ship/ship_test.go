package ship

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	pb "github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testBus struct {
	sub     logbus.Subscriber
	subChan chan struct{}
}

func (b *testBus) AddSubscriber(sub logbus.Subscriber, replay bool) {
	b.sub = sub
	b.subChan <- struct{}{}
}

func (b *testBus) RemoveSubscriber(sub logbus.Subscriber) {
	b.sub = nil
}

func (b *testBus) Write(d *pb.Delta) {
	b.sub.Write(d)
}

type testClient struct {
	done   chan struct{}
	deltas []*pb.Delta
}

func (c *testClient) StreamLogs(ctx context.Context, buildID string, deltas cloud.DeltaIterator) error {
	for {
		d, err := deltas.Next(ctx)
		if errors.Is(err, io.EOF) {
			c.done <- struct{}{}
			return nil
		} else if errors.Is(err, cloud.ErrNoDeltas) {
			time.Sleep(time.Millisecond * 10) // To prevent a tight loop.
			continue
		} else if err != nil {
			return err
		}
		c.deltas = append(c.deltas, d...)
	}
}

type failClient struct {
	done   chan struct{}
	deltas []*pb.Delta
}

func (c *failClient) StreamLogs(ctx context.Context, buildID string, deltas cloud.DeltaIterator) error {
	fmt.Println("STREAM LOGS")
	r := rand.New(rand.NewSource(1337)) // Consistent seed for deterministic tests.
	var count int
	for {
		d, err := deltas.Next(ctx)
		spew.Dump(d, err)
		if errors.Is(err, io.EOF) {
			fmt.Println("done")
			c.done <- struct{}{}
			return nil
		} else if errors.Is(err, cloud.ErrNoDeltas) {
			fmt.Println("none")
			time.Sleep(time.Millisecond * 10) // To prevent a tight loop.
			continue
		} else if err != nil {
			return err
		}
		c.deltas = append(c.deltas, d...)
		count++
		fmt.Println(count)
		if r.Intn(4) == 0 {
			fmt.Println("error ")
			return status.Error(codes.Unavailable, "unavailable")
		}
	}
}

func TestLogShipper(t *testing.T) {
	ctx := context.Background()

	bus := &testBus{
		subChan: make(chan struct{}),
	}

	cl := &testClient{
		done: make(chan struct{}, 1),
	}

	shipper := &LogShipper{
		man: &pb.RunManifest{
			BuildId:     uuid.NewString(),
			UserId:      uuid.NewString(),
			OrgName:     "my-org",
			ProjectName: "my-project",
		},
		bus:       bus,
		cl:        cl,
		done:      make(chan struct{}),
		retryWait: time.Duration(1), // No delay.
		first:     true,
		chunkSize: 10,
	}

	shipper.Start(ctx)

	<-bus.subChan

	bus.Write(&pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaLog{
			DeltaLog: &pb.DeltaLog{
				TargetId:  "target",
				CommandId: "command",
				Data:      []byte("info"),
			},
		},
	})

	shipper.Close()

	fmt.Println("done")

	<-cl.done

	if got := len(cl.deltas); got != 2 {
		t.Fatalf("expected 2 deltas, got %d", got)
	}

	if cl.deltas[0].GetDeltaManifest().GetResetAll() == nil {
		t.Fatal("expected first delta to be reset")
	}
}

func TestLogShipperResume(t *testing.T) {
	ctx := context.Background()

	bus := &testBus{
		subChan: make(chan struct{}, 100),
	}

	cl := &failClient{
		done: make(chan struct{}, 1),
	}

	shipper := &LogShipper{
		man: &pb.RunManifest{
			BuildId:     uuid.NewString(),
			UserId:      uuid.NewString(),
			OrgName:     "my-org",
			ProjectName: "my-project",
		},
		bus:       bus,
		cl:        cl,
		done:      make(chan struct{}),
		retryWait: time.Duration(1), // No delay.
		first:     true,
		chunkSize: 10,
	}

	shipper.Start(ctx)

	<-bus.subChan

	for i := 0; i < 60; i++ {
		bus.Write(&pb.Delta{
			DeltaTypeOneof: &pb.Delta_DeltaLog{
				DeltaLog: &pb.DeltaLog{
					TargetId:  "target",
					CommandId: "command",
					Data:      []byte("info"),
				},
			},
		})
	}

	shipper.Close()

	fmt.Println("closed")

	<-cl.done

	spew.Dump(cl.deltas)

	if got := len(cl.deltas); got != 2 {
		t.Fatalf("expected 2 deltas, got %d", got)
	}

	if cl.deltas[0].GetDeltaManifest().GetResetAll() == nil {
		t.Fatal("expected first delta to be reset")
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
