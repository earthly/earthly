package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/earthly/earthly/conslogging"

	"github.com/creack/pty"
	"github.com/hashicorp/yamux"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

// DebugServer provides a server that accepts a remote debugging shell
type DebugServer struct {
	session *session
	console conslogging.ConsoleLogger

	ctx    context.Context
	cancel context.CancelFunc

	sigs chan os.Signal
	addr string
}

func (ds *DebugServer) handleRequest(conn net.Conn) error {
	defer conn.Close()

	yaSession, err := yamux.Server(conn, nil)
	if err != nil {
		return errors.Wrap(err, "failed creating yamux server")
	}

	ctx, cancel := context.WithCancel(ds.ctx)
	ds.session = &session{
		yaSession: yaSession,
		ctx:       ctx,
		cancel:    cancel,
		server:    ds,
	}
	defer cancel()
	defer func() { ds.session = nil }()

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Wrap(err, "failed initializing raw terminal mode")
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

	return ds.session.handle()
}

func (ds *DebugServer) windowResizeHandler() error {
	for {
		select {
		case _ = <-ds.sigs:
			break

		case <-ds.ctx.Done():
			return nil
		}
		if len(ds.sigs) > 0 {
			continue
		}
		size, err := pty.GetsizeFull(os.Stdin)
		if err != nil {
			ds.console.Warnf("failed to get size: %v\n", err)
		} else {
			if ds.session != nil {
				ds.session.sendNewWindowSize(size)
			}
		}
	}
}

// Start starts the debug server listener
func (ds *DebugServer) Start() error {
	l, err := net.Listen("unix", ds.addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed listening on %s", ds.addr))
	}

	go ds.windowResizeHandler()

	go func() {
		ds.console.Printf("interactive debugger listening on %v\n", ds.addr)
		defer l.Close()
		defer fmt.Printf("deleting %v", ds.addr)
		defer os.Remove(ds.addr)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				ds.console.Warnf("Error accepting: %v", err.Error())
				os.Exit(1)
			}
			// Handle connections in a new goroutine.
			err = ds.handleRequest(conn)
			if err != nil && err != io.EOF {
				ds.console.Warnf("lost connection to interactive debugger: %v\n", err)
			} else {
				ds.console.Printf("interactive debugger closed\n")
			}
		}
	}()
	return nil
}

// Stop stops the server
func (ds *DebugServer) Stop() {
	ds.cancel()
}

// NewDebugServer returns a new deubgging server
func NewDebugServer(ctx context.Context, console conslogging.ConsoleLogger, socketPath string) *DebugServer {
	sigs := make(chan os.Signal, 100)
	signal.Notify(sigs, syscall.SIGWINCH)

	ctx, cancel := context.WithCancel(ctx)
	srv := &DebugServer{
		console: console,
		sigs:    sigs,
		ctx:     ctx,
		cancel:  cancel,
		addr:    socketPath,
	}

	return srv
}
