package socketforward

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/session"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
)

const SocketIDKey = "socket_id"

type server struct {
	caller session.Caller
}

func (s *server) run(ctx context.Context, l net.Listener, id string) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	eg.Go(func() error {
		for {
			conn, err := l.Accept()
			if err != nil {
				return err
			}

			client := NewSocketClient(s.caller.Conn())

			opts := make(map[string][]string)
			opts[SocketIDKey] = []string{id}
			ctx = metadata.NewOutgoingContext(ctx, opts)

			stream, err := client.Proxy(ctx)
			if err != nil {
				conn.Close()
				return err
			}

			go Copy(ctx, conn, stream, stream.CloseSend)
		}
	})

	return eg.Wait()
}

type SocketOpt struct {
	ID   string
	UID  int
	GID  int
	Mode int
}

func MountSocket(ctx context.Context, c session.Caller, opt SocketOpt) (sockPath string, closer func() error, err error) {
	dir, err := os.MkdirTemp("", ".buildkit-socketforward")
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	if err := os.Chmod(dir, 0711); err != nil {
		return "", nil, errors.WithStack(err)
	}

	sockPath = filepath.Join(dir, "sock")

	l, err := net.Listen("unix", sockPath)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if err := os.Chown(sockPath, opt.UID, opt.GID); err != nil {
		l.Close()
		return "", nil, errors.WithStack(err)
	}
	if err := os.Chmod(sockPath, os.FileMode(opt.Mode)); err != nil {
		l.Close()
		return "", nil, errors.WithStack(err)
	}

	s := &server{caller: c}

	id := opt.ID
	if id == "" {
		return "", nil, fmt.Errorf("socket ID is missing")
	}

	go s.run(ctx, l, id) // erroring per connection allowed

	return sockPath, func() error {
		err := l.Close()
		os.RemoveAll(sockPath)
		return errors.WithStack(err)
	}, nil
}
