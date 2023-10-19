package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithAdditionalHeaders adds metadata to all outgoing gRPC calls. It acts just like metadata.Pairs does, and calls *will* fail if
// you provide an odd number of items. Values should be in alternating key-value order - ("k1", "v1", "k1", "v2", "k2", "v3").
func WithAdditionalMetadataContext(kv ...string) ClientOpt {
	return &withAdditionalHeaders{kv}
}

type withAdditionalHeaders struct {
	kv []string
}

func (*withAdditionalHeaders) isClientOpt() {}

func headersUnaryInterceptor(kv ...string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, kv...)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func headersStreamInterceptor(kv ...string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, kv...)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
