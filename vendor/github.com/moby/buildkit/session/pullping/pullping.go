package pullping

import (
	context "context"

	"github.com/moby/buildkit/session"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type PullCallback func(ctx context.Context, images []string, resp map[string]string) error

type pullPing struct {
	callback PullCallback
}

// NewPullPing creates a new PullPing attachable, given a callback. The callback will be issued
// when the remote registry is ready for pulls. The callback should return when the pulls
// have completed, in order to inform the remote exporter that the pulls have completed.
func NewPullPing(cb PullCallback) session.Attachable {
	return &pullPing{callback: cb}
}

// Register registers the object with a grpc server.
func (pp *pullPing) Register(server *grpc.Server) {
	RegisterPullPingServer(server, pp)
}

// Pull implements the gRPC Pull message. It calls the callback and sends a
// message to the client when the callback has completed.
func (pp *pullPing) Pull(pr *PullRequest, ps PullPing_PullServer) error {
	err := pp.callback(ps.Context(), pr.GetImages(), pr.GetResp())
	if err != nil {
		return err
	}
	err = ps.Send(&PullResponse{})
	if err != nil {
		return err
	}
	return nil
}

// PullPingChannel returns a channel which signals an error when the pull
// operation has completed from the client-side. The error is nil
// if the operation has been completed successfully.
func PullPingChannel(ctx context.Context, images []string, resp map[string]string, c session.Caller) chan error {
	respChan := make(chan error, 1)
	ppc := NewPullPingClient(c.Conn())
	pc, err := ppc.Pull(ctx, &PullRequest{
		Images: images,
		Resp:   resp,
	})
	if err != nil {
		respChan <- errors.Wrap(err, "pull ping request")
		return respChan
	}
	go func() {
		_, err := pc.Recv()
		respChan <- errors.Wrap(err, "pull ping response")
	}()
	return respChan
}
