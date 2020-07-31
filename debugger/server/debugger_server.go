package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/debugger/common"

	"github.com/creack/pty"
	"github.com/hashicorp/yamux"
	"golang.org/x/crypto/ssh/terminal"
)

type session struct {
	yaSession *yamux.Session

	ctx    context.Context
	cancel context.CancelFunc

	ttyCon     net.Conn
	resizeConn net.Conn

	server *DebugServer
}

// DebugServer provides a server that accepts a remote debugging shell
type DebugServer struct {
	session *session
	console conslogging.ConsoleLogger

	ctx    context.Context
	cancel context.CancelFunc

	sigs chan os.Signal
	addr string
}

func (s *session) sendNewWindowSize(size *pty.Winsize) error {
	b, err := json.Marshal(size)
	if err != nil {
		return err
	}
	return common.WriteUint16PrefixedData(s.resizeConn, b)
}

func (s *session) handle() error {
	for {
		stream, err := s.yaSession.Accept()
		if err != nil {
			return err
		}

		buf := make([]byte, 1)
		stream.Read(buf)

		switch buf[0] {
		case common.PtyStream:
			go s.handle1(stream)

		case common.WinChangeStream:
			s.resizeConn = stream
			s.server.sigs <- syscall.SIGWINCH
		default:
			return fmt.Errorf("unsupported stream code %v", buf[0])
		}
	}
}

func (s *session) handle1(conn net.Conn) error {
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		s.cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		s.cancel()
	}()

	<-s.ctx.Done()
	return nil
}

func (ds *DebugServer) handleRequest(conn net.Conn) error {
	defer conn.Close()

	yaSession, err := yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
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
		return err
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
func (ds *DebugServer) Start() (string, error) {
	l, err := net.Listen("tcp", ds.addr)
	if err != nil {
		return "", err
	}

	go ds.windowResizeHandler()

	go func() {
		ds.console.Printf("interactive debugger listening\n")
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
	return l.Addr().String(), nil
}

// Stop stops the server
func (ds *DebugServer) Stop() {
	ds.cancel()
}

// NewDebugServer returns a new deubgging server
func NewDebugServer(console conslogging.ConsoleLogger) *DebugServer {
	sigs := make(chan os.Signal, 100)
	signal.Notify(sigs, syscall.SIGWINCH)

	ctx, cancel := context.WithCancel(context.Background())
	srv := &DebugServer{
		console: console,
		sigs:    sigs,
		ctx:     ctx,
		cancel:  cancel,
		addr:    "127.0.0.1:0",
	}

	return srv
}
