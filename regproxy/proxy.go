package regproxy

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"

	registry "github.com/moby/buildkit/api/services/registry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// newRegistryProxy creates and returns a new registry proxy that streams Docker
// container image data from the BK embedded Docker registry.
func newRegistryProxy(ln net.Listener, cl registry.RegistryClient) *registryProxy {
	return &registryProxy{ln: ln, cl: cl, errCh: make(chan error)}
}

// registryProxy uses a gRPC stream to translate incoming Docker image requests
// into a gRPC byte stream and back out into a valid HTTP response. The data is
// streamed over the gRPC connection rather than buffered as the images can be
// rather large.
type registryProxy struct {
	ln    net.Listener
	cl    registry.RegistryClient
	errCh chan error
	done  atomic.Bool
}

// Serve waits for TCP connections and pipes data received from the connection
// to BK via the gRPC server.
func (r *registryProxy) serve(ctx context.Context) {
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
			conn, err := r.ln.Accept()
			if err != nil {
				if !r.done.Load() {
					r.errCh <- errors.Wrap(err, "failed to accept")
				}
				return
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.errCh <- r.handle(ctx, conn)
			}()
		}
	}
}

func (r *registryProxy) close() {
	r.done.Store(true)
	r.ln.Close()
}

func (r *registryProxy) err() <-chan error {
	return r.errCh
}

func (r *registryProxy) handle(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	stream, err := r.cl.Proxy(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create proxy stream")
	}

	rw := registry.NewStreamRW(stream)
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err = registry.CopyWithDeadline(conn, rw)
		if err != nil {
			return errors.Wrap(err, "failed to write to stream")
		}
		err = stream.CloseSend()
		if err != nil {
			return errors.Wrap(err, "failed to close stream")
		}
		return nil
	})

	eg.Go(func() error {
		_, err = io.Copy(conn, rw)
		if err != nil {
			return errors.Wrap(err, "failed to read from stream")
		}
		return nil
	})

	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait")
	}

	return nil
}
