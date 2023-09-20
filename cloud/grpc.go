package cloud

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var RPCErrRegex = regexp.MustCompile(`(?U)rpc error: code = .+ desc = `)

func (c *Client) withAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", c.authToken))
}

func (c *Client) withReqID(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, requestID, c.getRequestID())
}

func getReqID(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		vals := md.Get(requestID)
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

type interceptorOpts struct {
	skipAuth map[string]struct{}
}

type InterceptorOpt func(opt *interceptorOpts)

func WithSkipAuth(methods ...string) InterceptorOpt {
	return func(opts *interceptorOpts) {
		if opts.skipAuth == nil {
			opts.skipAuth = map[string]struct{}{}
		}
		for _, method := range methods {
			opts.skipAuth[method] = struct{}{}
		}
	}
}

// UnaryInterceptor is a unary middleware function for the Earthly gRPC client which
// handle re-authentication when necessary, and automatically
// prints requestIDs to errors when errors are received from the server.
func (c *Client) UnaryInterceptor(opts ...InterceptorOpt) grpc.UnaryClientInterceptor {
	interceptorOpts := &interceptorOpts{}
	for _, opt := range opts {
		opt(interceptorOpts)
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = c.withReqID(ctx)
		if _, ok := interceptorOpts.skipAuth[method]; ok {
			// It would probably be better to break this interceptor into multiple so that skipping auth doesn't affect anything else that
			// might be added here in the future
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		ctx, err := c.reAuthIfExpired(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed refreshing expired token: %s", getReqID(ctx))
		}
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				return fmt.Errorf("%s {reqID: %s}", err.Error(), getReqID(ctx))
			}
			if s.Code() == codes.Unauthenticated {
				ctx, err = c.reAuthCtx(ctx)
				if err != nil {
					return fmt.Errorf("%s {reqID: %s}", err.Error(), getReqID(ctx))
				}
				return invoker(ctx, method, req, reply, cc, opts...)
			}
			return status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}", cleanStatusError(err.Error()), getReqID(ctx)))
		}
		return nil
	}
}

// StreamInterceptor is a stream middleware function for the Earthly gRPC client which
// handle re-authentication when necessary, and automatically
// prints requestIDs to errors when errors are received from the server.
func (c *Client) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = c.withReqID(ctx)
		ctx, err := c.reAuthIfExpired(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "failed refreshing expired token {reqID: %s}", getReqID(ctx))
		}
		newStreamer, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				return nil, fmt.Errorf("%s {reqID: %s}", err.Error(), getReqID(ctx))
			}
			if s.Code() == codes.Unauthenticated {
				ctx, err = c.reAuthCtx(ctx)
				if err != nil {
					return nil, fmt.Errorf("%s {reqID: %s}", err.Error(), getReqID(ctx))
				}
				return streamer(ctx, desc, cc, method, opts...)
			}
			return nil, status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}",
				cleanStatusError(err.Error()), getReqID(ctx)))
		}
		return newWrappedStream(ctx, newStreamer), nil
	}
}

// wrappedStream  wraps around the embedded grpc.ClientStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ClientStream
	ctx context.Context
}

func (w *wrappedStream) RecvMsg(m any) error {
	if err := w.ClientStream.RecvMsg(m); err != nil {
		if err == io.EOF {
			return err
		}
		s, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("%s {reqID: %s}", err.Error(), getReqID(w.ctx))
		}
		return status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}",
			cleanStatusError(err.Error()), getReqID(w.ctx)))
	}
	return nil
}

func (w *wrappedStream) SendMsg(m any) error {
	if err := w.ClientStream.SendMsg(m); err != nil {
		if err == io.EOF {
			return err
		}
		s, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("%s {reqID: %s}", err.Error(), getReqID(w.ctx))
		}
		return status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}",
			cleanStatusError(err.Error()), getReqID(w.ctx)))
	}
	return nil
}

func newWrappedStream(ctx context.Context, s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s, ctx}
}

// cleanStatusError returns the underlying error message from a gRPC status error
func cleanStatusError(errStr string) string {
	return RPCErrRegex.ReplaceAllString(errStr, "")
}

func (c *Client) reAuthIfExpired(ctx context.Context) (context.Context, error) {
	if time.Now().UTC().After(c.authTokenExpiry) {
		_, err := c.Authenticate(ctx)
		if err != nil {
			return ctx, errors.Wrap(err, "failed refreshing expired token")
		}
		ctx = c.resetAuthMeta(ctx)
	}
	return ctx, nil
}

func (c *Client) reAuthCtx(ctx context.Context) (context.Context, error) {
	_, err := c.Authenticate(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed re-authenticating")
	}
	return c.resetAuthMeta(ctx), nil
}

func (c *Client) resetAuthMeta(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Delete("authorization")
	} else {
		md = metadata.New(map[string]string{})
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return c.withAuth(ctx)
}
