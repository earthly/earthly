package cloud

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *client) withAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", c.authToken))
}

func (c *client) withReqID(ctx context.Context) context.Context {
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

// UnaryInterceptor is a unary middleware function for the Earthly gRPC client which
// handle re-authentication when necessary, and automatically
// prints requestIDs to errors when errors are received from the server.
func (c *client) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = c.withReqID(ctx)
		ctx, err := c.reAuthIfExpired(ctx)
		if err != nil {
			return errors.Wrap(err, "failed refreshing expired token")
		}
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				return err
			}
			err = status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}", s.Err(), getReqID(ctx)))
			if s.Code() == codes.Unauthenticated {
				ctx, err = c.reAuthCtx(ctx)
				if err != nil {
					return err
				}
				return invoker(ctx, method, req, reply, cc, opts...)
			}
			return err
		}
		return nil
	}
}

// StreamInterceptor is a stream middleware function for the Earthly gRPC client which
// handle re-authentication when necessary, and automatically
// prints requestIDs to errors when errors are received from the server.
func (c *client) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = c.withReqID(ctx)
		ctx, err := c.reAuthIfExpired(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed refreshing expired token")
		}
		newStreamer, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				return newStreamer, err
			}
			err = status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}", s.Err(), getReqID(ctx)))
			if s.Code() == codes.Unauthenticated {
				ctx, err = c.reAuthCtx(ctx)
				if err != nil {
					return newStreamer, err
				}
				return streamer(ctx, desc, cc, method, opts...)
			}
			return newStreamer, err
		}
		return newStreamer, nil
	}
}

func (c *client) reAuthIfExpired(ctx context.Context) (context.Context, error) {
	if time.Now().UTC().After(c.authTokenExpiry) {
		err := c.Authenticate(ctx)
		if err != nil {
			return ctx, errors.Wrap(err, "failed refreshing expired token")
		}
		ctx = c.resetAuthMeta(ctx)
	}
	return ctx, nil
}

func (c *client) reAuthCtx(ctx context.Context) (context.Context, error) {
	err := c.Authenticate(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed re-authenticating")
	}
	return c.resetAuthMeta(ctx), nil
}

func (c *client) resetAuthMeta(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md.Delete("authorization")
	} else {
		md = metadata.New(map[string]string{})
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return c.withAuth(ctx)
}
