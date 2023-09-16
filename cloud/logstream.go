package cloud

import (
	"context"
	"io"
	"sync/atomic"

	pb "github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) StreamLogs(ctx context.Context, man *pb.RunManifest, ch <-chan *pb.Delta, verbose bool) error {
	verbose = true // Debug
	if man.GetResumeToken() == "" {
		man.ResumeToken = stringutil.RandomAlphanumeric(40)
	}
	var (
		errs  []error
		retry bool
	)
	for {
		first := firstDelta(man, retry)
		err := c.streamLogsAttempt(ctx, man.GetBuildId(), first, ch)
		if err != nil {
			if retryable(err) {
				retry = true
				if verbose {
					errs = append(errs, err)
				}
				continue
			} else {
				errs = append(errs, err)
			}
		}
		break
	}
	if len(errs) > 0 {
		return &multierror.Error{Errors: errs}
	}
	return nil
}

func (c *Client) streamLogsAttempt(ctx context.Context, buildID string, first *pb.Delta, ch <-chan *pb.Delta) error {
	stream, err := c.logstream.StreamLogs(c.withAuth(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to create log stream client")
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

	sendSingle := func(delta *pb.Delta) error {
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

	eg.Go(func() error {
		// Send the reset or resume delta first.
		err := sendSingle(first)
		if err != nil {
			return err
		}
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case delta, ok := <-ch:
				if !ok {
					finished.Store(true)
					msg := &pb.StreamLogRequest{
						BuildId: buildID,
						Eof:     true,
					}
					err := stream.Send(msg)
					if err != nil {
						return errors.Wrap(err, "failed to send EOF to log stream")
					}
					return nil
				}
				err := sendSingle(delta)
				if err != nil {
					return err
				}
			}
		}
	})

	return eg.Wait()
}

func retryable(err error) bool {
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
