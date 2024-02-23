package cloud

import (
	"context"
	"io"
	"math"
	"sync/atomic"
	"time"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamError struct {
	Recoverable bool
	Err         error
}

func (s *StreamError) Error() string {
	return s.Err.Error()
}

func (c *Client) StreamLogs(ctx context.Context, man *pb.RunManifest, ch <-chan *pb.Delta) <-chan error {
	if man.GetResumeToken() == "" {
		man.ResumeToken = stringutil.RandomAlphanumeric(40)
	}
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		var (
			retry  bool
			counts []int
			last   *pb.Delta
		)
		for {
			first := []*pb.Delta{firstDelta(man, retry)}
			if last != nil {
				first = append(first, last)
			}
			var (
				err       error
				sendCount int
			)
			sendCount, last, err = c.streamLogsAttempt(ctx, man.GetBuildId(), first, ch)
			if err != nil {
				recoverable := recoverableError(err)
				errCh <- &StreamError{Err: err, Recoverable: recoverable}
				if recoverable {
					retry = true
					counts = append(counts, sendCount)
					select {
					case <-time.After(calcBackoff(c.logstreamBackoff, counts)):
					case <-ctx.Done():
						errCh <- ctx.Err()
						return
					}
					continue
				}
			}
			break
		}
	}()
	return errCh
}

// calcBackoff calculates the time to wait before attempting another stream. The
// backoff is calculated by inspecting the passed attempts for outright failures
// versus partial streams.
func calcBackoff(base time.Duration, counts []int) time.Duration {
	fails := 0
	for i := 0; i < len(counts); i++ {
		// If a previous stream failed before sending 50 messages, we'll
		// increment the backoff multiplier. Otherwise reset it.
		if counts[i] < 50 {
			fails++
		} else {
			fails = 0
		}
	}
	// min(backoff + backoff^(1 + (fails / 5)), 40s)
	b := float64(base.Milliseconds())
	t := b * math.Pow(2, float64(fails))
	return time.Millisecond * time.Duration(math.Min(float64(40_000), t))
}

func (c *Client) streamLogsAttempt(ctx context.Context, buildID string, first []*pb.Delta, ch <-chan *pb.Delta) (int, *pb.Delta, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := c.logstream.StreamLogs(c.withAuth(ctx))
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to create log stream client")
	}

	eg, ctx := errgroup.WithContext(ctx)
	var finished atomic.Bool

	eg.Go(func() error {
		for {
			resp, err := stream.Recv()
			if err != nil {
				return errors.Wrap(err, "failed waiting for log stream server")
			}
			if resp.GetEofAck() {
				if !finished.Load() {
					return errors.New("unexpected EOF ack")
				}
				err := stream.CloseSend()
				if err != nil {
					return errors.Wrap(err, "failed to close log stream")
				}
				return nil
			}
		}
	})

	// If an error occurs, we can't assume that the last delta was correctly
	// sent. Let's return it for subsequent attempts. Note that writing log
	// entries is idempotent on the server side.
	var last *pb.Delta

	sendSingle := func(delta *pb.Delta) error {
		last = delta
		msg := &pb.StreamLogRequest{
			BuildId: buildID,
			Deltas:  []*pb.Delta{delta},
		}
		err := stream.Send(msg)
		if err != nil {
			return errors.Wrap(err, "failed to send log data")
		}
		return nil
	}

	var count int

	eg.Go(func() (err error) {
		defer func() {
			// Ensure that the receive side correctly exits when an error occurs.
			if err != nil {
				cancel()
			}
		}()
		// Send the reset or resume delta first. Also sends any dropped deltas
		// from a previous attempt.
		for _, d := range first {
			err = sendSingle(d)
			if err != nil {
				return
			}
		}
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case delta, ok := <-ch:
				if !ok {
					last = nil // Clear the last sent delta so it's not sent again.
					finished.Store(true)
					msg := &pb.StreamLogRequest{
						BuildId: buildID,
						Eof:     true,
					}
					err = stream.Send(msg)
					if err != nil {
						err = errors.Wrap(err, "failed to send EOF to log stream")
						return
					}
					return nil
				}
				err = sendSingle(delta)
				if err != nil {
					return
				}
				// We track & return the total count of successfully sent messages.
				count++
			}
		}
	})

	return count, last, eg.Wait()
}

func recoverableError(err error) bool {
	if errors.Is(err, io.EOF) {
		return true
	}
	for {
		if err == nil {
			return false
		}
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unavailable, codes.Unknown:
				return true
			default:
				return false
			}
		}
		err = errors.Unwrap(err)
	}
}

func firstDelta(man *pb.RunManifest, retry bool) *pb.Delta {
	var deltaMan *pb.DeltaManifest

	if retry {
		deltaMan = &pb.DeltaManifest{
			DeltaManifestOneof: &pb.DeltaManifest_Resume{
				Resume: &pb.DeltaManifest_ResumeBuild{
					BuildId:     man.GetBuildId(),
					Token:       man.GetResumeToken(),
					OrgName:     man.GetOrgName(),
					ProjectName: man.GetProjectName(),
				},
			},
		}
	} else {
		deltaMan = &pb.DeltaManifest{
			DeltaManifestOneof: &pb.DeltaManifest_ResetAll{
				ResetAll: man,
			},
		}
	}

	return &pb.Delta{
		DeltaTypeOneof: &pb.Delta_DeltaManifest{
			DeltaManifest: deltaMan,
		},
	}
}
