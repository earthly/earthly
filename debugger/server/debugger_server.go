package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"runtime"
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

func (ds *DebugServer) getIP() string {
	if runtime.GOOS == "darwin" {
		// macOS doesn't have a docker0 bridge
		return "127.0.0.1"
	}

	iface, err := net.InterfaceByName("docker0")
	if err != nil {
		ds.console.Warnf("falling back to 0.0.0.0 due to docker0 lookup error: %v", err.Error())
		return "0.0.0.0"
	}
	addrs, err := iface.Addrs()
	if err != nil {
		ds.console.Warnf("falling back to 0.0.0.0 due to docker0 addrs error: %v", err.Error())
		return "0.0.0.0"
	}
	for _, a := range addrs {
		switch v := a.(type) {
		case *net.IPNet:
			if x := v.IP.To4(); x != nil {
				return x.String()
			}
		}
	}

	ds.console.Warnf("falling back to 0.0.0.0 due to docker0 addrs being empty")
	return "0.0.0.0"
}

// Start starts the debug server listener
func (ds *DebugServer) Start() (string, error) {
	addr := fmt.Sprintf("%s:0", ds.getIP())

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed listening on %s", addr))
	}

	go ds.windowResizeHandler()

	go func() {
		ds.console.Printf("Interactive debugger listening on %v\n", l.Addr())
		defer l.Close()
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

	assignedAddr := l.Addr()

	if runtime.GOOS == "darwin" {
		tcpAddr, ok := assignedAddr.(*net.TCPAddr)
		if !ok {
			panic("failed to cast to TCPAddr (shouldn't happen)")
		}

		// under macOS this dns points back to the host (but doesn't work under linux)
		return fmt.Sprintf("host.docker.internal:%d", tcpAddr.Port), nil
	}

	return assignedAddr.String(), nil
}

// Stop stops the server
func (ds *DebugServer) Stop() {
	ds.cancel()
}

// NewDebugServer returns a new deubgging server
func NewDebugServer(ctx context.Context, console conslogging.ConsoleLogger) *DebugServer {
	sigs := make(chan os.Signal, 100)
	signal.Notify(sigs, syscall.SIGWINCH)

	ctx, cancel := context.WithCancel(ctx)
	srv := &DebugServer{
		console: console,
		sigs:    sigs,
		ctx:     ctx,
		cancel:  cancel,
	}

	return srv
}
