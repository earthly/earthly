package cloud

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *client) withAuth(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", c.authToken))
}

func withReqID(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "request-id", uuid.NewString())
}

func getReqID(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		vals := md.Get("request-id")
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return ""
}

func (c *client) UnaryAuthInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = withReqID(ctx)
		if time.Now().UTC().After(c.authTokenExpiry) {
			err := c.Authenticate(ctx)
			if err != nil {
				return errors.Wrap(err, "failed refreshing expired token")
			}
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				err = status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}", s.Err(), getReqID(ctx)))
				if s.Code() == codes.Unauthenticated {
					err = c.Authenticate(ctx)
					if err != nil {
						return errors.Wrap(err, "failed re-authenticating")
					}
					md, ok := metadata.FromOutgoingContext(ctx)
					if ok {
						md.Delete("authorization")
					} else {
						md = metadata.New(map[string]string{})
					}
					ctx = metadata.NewOutgoingContext(ctx, md)
					return invoker(c.withAuth(ctx), method, req, reply, cc, opts...)
				}
			}
		}
		return err
	}
}

func (c *client) StreamAuthInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = withReqID(ctx)
		if time.Now().UTC().After(c.authTokenExpiry) {
			err := c.Authenticate(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "failed refreshing expired token")
			}
		}
		newStreamer, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			s, ok := status.FromError(err)
			if ok {
				err = status.Errorf(s.Code(), fmt.Sprintf("%s {reqID: %s}", s.Err(), getReqID(ctx)))
				if s.Code() == codes.Unauthenticated {
					err = c.Authenticate(ctx)
					if err != nil {
						return nil, errors.Wrap(err, "failed re-authenticating")
					}
					md, ok := metadata.FromOutgoingContext(ctx)
					if ok {
						md.Delete("authorization")
					} else {
						md = metadata.New(map[string]string{})
					}
					ctx = metadata.NewOutgoingContext(ctx, md)
					// TODO not sure if newStreamer(...) should be called here instead
					return streamer(c.withAuth(ctx), desc, cc, method, opts...)
				}
			}
			return newStreamer, err
		}
		return newStreamer, nil
	}
}
