//go:build chaos

package logstreamer_test

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	"git.sr.ht/~nelsam/hel/v4/pkg/pers"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/logstreamer"
	"github.com/pkg/errors"
	"github.com/poy/onpar/v2/expect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testTimeout = 10 * time.Second

func newTestLogStreamer(ctx context.Context, client *mockCloudClient, initMani *logstream.RunManifest, opts ...logstreamer.Opt) (*logstreamer.LogStreamer, func(*testing.T)) {
	str := logstreamer.New(ctx, logbus.New(), client, initMani, opts...)
	cleanup := func(t *testing.T) {
		pers.ConsistentlyReturn(t, client.StreamLogsOutput, errors.New("done"))
		str.Close()
	}
	return str, cleanup
}

// TestDataRace_InitManifest ensures that there's no race condition on sending
// the initial manifest when there are calls to Write blocked on buffer
// allocation.
func TestDataRace_InitManifest(t *testing.T) {
	const runs = 5
	for i := 0; i < runs; i++ {
		// While this test fails a little more consistently than the deadlock
		// test, it still passes sometimes, depending on scheduler luck.
		// Exercising it multiple times gives us a pretty good chance of seeing
		// a failure every time. These don't have to run concurrently though.
		t.Run(fmt.Sprintf("run %v", i), exerciseInitManifest)
	}
}

func exerciseInitManifest(t *testing.T) {
	client := newMockCloudClient(t, testTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	initManifest := &logstream.RunManifest{
		BuildId: "foo",
		Version: 1,
	}
	initDelta := &logstream.Delta{
		DeltaTypeOneof: &logstream.Delta_DeltaManifest{
			DeltaManifest: &logstream.DeltaManifest{
				DeltaManifestOneof: &logstream.DeltaManifest_ResetAll{
					ResetAll: initManifest,
				},
			},
		},
	}
	str, cleanup := newTestLogStreamer(ctx, client, initManifest)
	defer cleanup(t)

	// To trigger this data race, all we need is for a call to `Write` to block
	// on access to a lock while the buffer is being allocated. If the lock is
	// released before the initial manifest is sent, then the scheduler can
	// schedule the `Write` goroutine before the initial manifest is added to
	// the buffer.
	//
	// To improve our chances of getting the scheduler to schedule one of those
	// calls, we ensure that there are plenty of them for it to choose from.
	const concurrentWrites = 100
	for i := 0; i < concurrentWrites; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					str.Write(&logstream.Delta{})
				}
			}
		}()
	}

	var (
		streamCtx context.Context
		buildID   string
		deltas    cloud.Deltas
	)
	expect.Expect(t, client).To(
		pers.HaveMethodExecuted(
			"StreamLogs",
			pers.Within(testTimeout),
			pers.StoreArgs(&streamCtx, &buildID, &deltas),
		),
	)

	dl, err := deltas.Next(ctx)
	if err != nil {
		t.Fatalf("expected a <nil> error from Next(); got '%v'", err)
	}
	if len(dl) != 1 {
		t.Fatalf("expected initial deltas to have length 1; got %d", len(dl))
	}
	if !reflect.DeepEqual(dl[0], initDelta) {
		t.Fatalf("expected to receive initial manifest first; got '%v'", dl[0])
	}
}

// TestTransientError_Deadlock tests for a specific (but uncommon) deadlock
// scenario. See the comments in the test code for detail, but the gist:
//
// In a scenario with a lot of congestion or a slow reader, when the cloud log
// stream RPC errors with a retryable code, we can leave `Write` calls in a
// deadlocked state.
func TestTransientError_Deadlock(t *testing.T) {
	const concurrency = 10

	for i := 0; i < concurrency; i++ {
		// Running this test multiple times concurrently exercises the scheduler
		// enough to get it to fail consistently.
		t.Run(fmt.Sprintf("goroutine %v", i), exerciseDeadlock)
	}
}

func exerciseDeadlock(t *testing.T) {
	t.Parallel()

	// This is a rare edge case, and getting it to trigger on fast CPUs with
	// many cores in a unit test can be almost impossible. So our first step is
	// to change a few settings to make the case trigger more consistently in
	// the unit test.
	//
	// GOMAXPROCS: simulate a system where most goroutines are _not_ being used
	// for the LogStreamer. In a unit test, most of the processing power is
	// dedicated to the test; in a real world scenario, the scheduler would be
	// scheduling a lot of non-logging work.
	//
	// chSize: shrink the size of the channel to make everything fill up and
	// drain more quickly. We expect this edge case to trigger maybe once every
	// several days in practice - in a unit test, we need to shorten that time
	// period by quite a lot.
	runtime.GOMAXPROCS(2)
	const chSize = 10

	client := newMockCloudClient(t, testTimeout)
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// We need a heavily congested log streamer, with a full channel.
	str, cleanup := newTestLogStreamer(ctx, client, &logstream.RunManifest{
		BuildId: "foo",
		Version: 1,
	}, logstreamer.WithBuffer(chSize))
	defer cleanup(t)

	expect.Expect(t, client).To(
		pers.HaveMethodExecuted(
			"StreamLogs",
			pers.Within(testTimeout),
			pers.WithArgs(pers.Any, "foo", pers.Any),
		),
	)

	remaining := chSize - 1 // initialManifest is already on the channel
	for i := 0; i < remaining; i++ {
		str.Write(&logstream.Delta{})
	}

	// If we congested the streamer correctly, then new calls to `Write` will
	// gain access to the _current_ channel and then block on writing to the
	// channel. We should be able to handle _effectively_ infinite blocked
	// `Write` calls when we de-congest the channel later.

	// NOTE: `WriteAsync` was written specifically for this test case. There was
	// no way to get a consistent test result when calling `Write` in a
	// goroutine, because we had no way to guarantee that all of the calls to
	// `Write` had gotten access to the current channel before we trigger the
	// decongestion logic. WriteAsync gets access to the channel synchronously
	// and then writes to the channel asynchronously.

	var unblocked []<-chan struct{}

	blocked := chSize
	for i := 0; i < chSize; i++ {
		unblocked = append(unblocked, str.WriteAsync(&logstream.Delta{}))
	}

	// At this point, we have cap(deltaCh) goroutines all blocked on sending to
	// the _current_ deltaCh. But we should be able to handle quite a lot
	// more...

	var extrasUnblocked []<-chan struct{}

	beyondCap := 1000 * chSize
	blocked += beyondCap
	for i := 0; i < beyondCap; i++ {
		extrasUnblocked = append(extrasUnblocked, str.WriteAsync(&logstream.Delta{}))
	}

	select {
	case <-unblocked[0]:
		t.Fatalf("setup failed: Write goroutines were not blocked")
	default:
	}

	// And now we cause the log streamer to retry, which means it reallocates
	// the channel.
	pers.Return(client.StreamLogsOutput, status.Error(codes.DeadlineExceeded, "BOOM!"))

	// At this point, our decongestion goroutine is running - all of our calls
	// to `WriteAsync` should eventually unblock, so all done channels should
	// close.

	timeout := time.After(testTimeout)

	for i, ch := range append(unblocked, extrasUnblocked...) {
		select {
		case <-ch:
		case <-timeout:
			t.Fatalf("timed out waiting for all WriteAsync calls to unblock; failed on call %d (out of %d)", i, len(unblocked)+len(extrasUnblocked))
		}
	}
}
