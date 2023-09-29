package regproxy

import (
	"context"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	registry "github.com/moby/buildkit/api/services/registry"
	"github.com/pkg/errors"
)

// NewRegistryProxy creates and returns a new RegistryProxy that streams Docker
// container image data from the BK embedded Docker registry.
func NewRegistryProxy(cl registry.RegistryClient) *RegistryProxy {
	return &RegistryProxy{cl: cl, errCh: make(chan error)}
}

// RegistryProxy uses a gRPC stream to translate incoming Docker image requests
// into a gRPC byte stream and back out into a valid HTTP response. The data is
// streamed over the gRPC connection rather than buffered as the images can be
// rather large.
type RegistryProxy struct {
	cl    registry.RegistryClient
	errCh chan error
}

// Serve waits for TCP connections and pipes data received from the connection
// to BK via the gRPC server.
func (r *RegistryProxy) Serve(ctx context.Context, ln net.Listener) {
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		close(r.errCh)
	}()
	for {
		select {
		case <-ctx.Done():
			r.errCh <- ctx.Err()
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				r.errCh <- errors.Wrap(err, "failed to accept")
				continue
			}
			go func() {
				wg.Add(1)
				defer wg.Done()
				r.errCh <- r.handle(ctx, conn)
			}()
		}
	}
}

func (r *RegistryProxy) Err() <-chan error {
	return r.errCh
}

func (r *RegistryProxy) handle(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	stream, err := r.cl.Proxy(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create proxy stream")
	}

	rw := registry.NewStreamRW(stream)

	connR := &httpConnReader{conn: conn}

	_, err = io.Copy(rw, connR)
	if err != nil {
		return errors.Wrap(err, "failed to write to stream")
	}

	err = stream.CloseSend()
	if err != nil {
		return errors.Wrap(err, "failed to close stream")
	}

	_, err = io.Copy(conn, rw)
	if err != nil {
		return errors.Wrap(err, "failed to read from stream")
	}

	return nil
}

type httpConnReader struct {
	conn net.Conn
	done bool
}

func (h *httpConnReader) Read(p []byte) (int, error) {
	if h.done {
		return 0, io.EOF
	}
	buf := make([]byte, 512)
	n, err := h.conn.Read(buf)
	if err != nil {
		return 0, err
	}
	buf = buf[0:n]
	if strings.HasSuffix(string(buf), "\r\n\r\n") {
		h.done = true
	}
	copy(p, buf)
	return n, nil
}
