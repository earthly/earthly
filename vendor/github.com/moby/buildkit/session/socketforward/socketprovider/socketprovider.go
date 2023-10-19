package socketprovider

import (
	"context"
	"fmt"
	"io"

	"github.com/moby/buildkit/session"
	socketforward "github.com/moby/buildkit/session/socketforward"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// SocketAcceptCb is called when the unix socket is connected to
type SocketAcceptCb func(context.Context, io.ReadWriteCloser) error

// NewSocketProvider creates a session provider that proxies a readwriter to a unix socket
func NewSocketProvider(m map[string]SocketAcceptCb) (session.Attachable, error) {
	return &socketProvider{m: m}, nil
}

type socketProvider struct {
	m map[string]SocketAcceptCb
}

func (sp *socketProvider) Register(server *grpc.Server) {
	socketforward.RegisterSocketServer(server, sp)
}

func (sp *socketProvider) Proxy(stream socketforward.Socket_ProxyServer) error {
	opts, _ := metadata.FromIncomingContext(stream.Context()) // if no metadata continue with empty object
	var id string
	if v, ok := opts[socketforward.SocketIDKey]; ok && len(v) > 0 && v[0] != "" {
		id = v[0]
	}

	cb, ok := sp.m[id]
	if !ok {
		return fmt.Errorf("no callback registered for socket ID: %s", id)
	}

	s1, s2 := sockPair()
	eg, ctx := errgroup.WithContext(stream.Context())
	eg.Go(func() error {
		return cb(ctx, s1)
	})
	eg.Go(func() error {
		defer s1.Close()
		return socketforward.Copy(ctx, s2, stream, nil)
	})
	return eg.Wait()
}

func sockPair() (io.ReadWriteCloser, io.ReadWriteCloser) {
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()
	return &sock{pr1, pw2, pw1}, &sock{pr2, pw1, pw2}
}

type sock struct {
	io.Reader
	io.Writer
	io.Closer
}
