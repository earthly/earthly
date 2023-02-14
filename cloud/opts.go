package cloud

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Opt is an option function for declaring a Client with option functions.
type Opt func(context.Context, *Client) (*Client, error)

// WithHTTPAddr uses addr as its connection to cloud for http requests.
func WithHTTPAddr(addr string) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.httpAddr = addr
		return c, nil
	}
}

// WithGRPCDial dials gRPC using the provided address and options. This also
// takes care of adding the Client's middleware (interceptors and the like) to
// the dial options.
func WithGRPCDial(addr string, opts ...grpc.DialOption) Opt {
	// TODO: this dependency is kind of cyclic. The *Client needs a gRPC
	// connection to perform most of its methods, but the gRPC connection needs
	// an allocated and configured client to construct its middleware.
	//
	// It probably makes sense to spin off the StreamInterceptor and
	// UnaryInterceptor methods onto their own Interceptor type to break the
	// cycle. But they depend on the auth methods on Client - so it probably
	// makes sense to spin off _those_ methods to an Authenticator type. And
	// there may be other pieces of the puzzle as we go deeper in the rabbit
	// hole.
	dialOpt := func(ctx context.Context, c *Client) (*Client, error) {
		retryOpts := []grpc_retry.CallOption{
			grpc_retry.WithMax(10),
			grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
			grpc_retry.WithCodes(codes.Internal, codes.Unavailable),
		}
		opts = append(opts,
			// c.StreamInterceptor() and c.UnaryInterceptor() depend on several
			// fields in c being set properly. These are the reason that this
			// option uses postAllocationOpts.
			grpc.WithChainStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...), c.StreamInterceptor()),
			grpc.WithChainUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...), c.UnaryInterceptor()),
			grpc.WithChainUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...), c.UnaryInterceptor(WithSkipAuth("/api.public.analytics.Analytics/SendAnalytics"))),
		)
		conn, err := grpc.DialContext(ctx, addr, opts...)
		if err != nil {
			return nil, errors.Wrapf(err, "cloud: failed to dial gRPC address [%v]", addr)
		}
		c.grpcConn = conn
		return c, nil
	}
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.postAllocationOpts = append(c.postAllocationOpts, dialOpt)
		return c, nil
	}
}

// WithGRPCConn uses conn as the grpc connection for gRPC requests against cloud.
func WithGRPCConn(conn grpc.ClientConnInterface) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.grpcConn = conn
		return c, nil
	}
}

// WithAgentSockPath sets the SSH agent socket path to use when communicating
// with cloud.
func WithAgentSockPath(path string) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.sshAgent = &lazySSHAgent{sockPath: path}
		return c, nil
	}
}

// AuthOverrideOpt is an option that may be passed to WithAuthOverride to
// determine how the auth will be overridden.
type AuthOverrideOpt func(ctx context.Context, c *Client) *Client

// JWT tells WithAuthOverride to override the JWT with the given token.
func JWT(tok string) AuthOverrideOpt {
	return func(ctx context.Context, c *Client) *Client {
		c.authToken = tok
		c.authTokenExpiry = time.Now().Add(24 * 365 * time.Hour) // Never expire when using JWT.
		return c
	}
}

// Creds tells WithAuthOverride to override the creds with the given token.
func Creds(tok string) AuthOverrideOpt {
	return func(ctx context.Context, c *Client) *Client {
		c.authCredToken = tok
		return c
	}
}

// WithAuthCredsOverride overrides the auth creds token. o must not be nil.
func WithAuthOverride(o AuthOverrideOpt) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c = o(ctx, c)
		return c, nil
	}
}

// WithInstallation sets the name of the installation for this client.
func WithInstallation(name string) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.installationName = name
		return c, nil
	}
}

// WithRequestID sets the ID that will be used when this client makes requests.
func WithRequestID(id string) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.requestID = id
		return c, nil
	}
}

// WithWarnCallback sets the callback function that will be called when
// transient or otherwise recoverable errors occur.
func WithWarnCallback(f func(string, ...any)) Opt {
	return func(ctx context.Context, c *Client) (*Client, error) {
		c.warnFunc = f
		return c, nil
	}
}
