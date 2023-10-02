package regproxy

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	registry "github.com/moby/buildkit/api/services/registry"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const readDeadline = 50 * time.Millisecond

// NewRegistryProxy creates and returns a new RegistryProxy that streams Docker
// container image data from the BK embedded Docker registry.
func NewRegistryProxy(ln net.Listener, cl registry.RegistryClient) *RegistryProxy {
	return &RegistryProxy{ln: ln, cl: cl, errCh: make(chan error)}
}

// RegistryProxy uses a gRPC stream to translate incoming Docker image requests
// into a gRPC byte stream and back out into a valid HTTP response. The data is
// streamed over the gRPC connection rather than buffered as the images can be
// rather large.
type RegistryProxy struct {
	ln    net.Listener
	cl    registry.RegistryClient
	errCh chan error
	done  atomic.Bool
}

// Serve waits for TCP connections and pipes data received from the connection
// to BK via the gRPC server.
func (r *RegistryProxy) Serve(ctx context.Context) {
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

func (r *RegistryProxy) Close() {
	r.done.Store(true)
	r.ln.Close()
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
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err = copyWithDeadline(conn, rw)
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

func copyWithDeadline(conn net.Conn, w io.Writer) (int64, error) {
	var t int64
	for {
		err := conn.SetReadDeadline(time.Now().Add(readDeadline))
		if err != nil {
			return 0, err
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) || isNetTimeout(err) {
				break
			}
			return 0, err
		}
		buf = buf[0:n]
		n, err = w.Write(buf)
		if err != nil {
			return 0, err
		}
		t += int64(n)
	}
	return t, nil
}

func isNetTimeout(err error) bool {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}
