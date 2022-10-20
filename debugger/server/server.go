package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/earthly/earthly/slog"

	"github.com/pkg/errors"
)

// Server provides a debugger server
type Server struct {
	shellConn    net.Conn
	terminalConn net.Conn
	mux          sync.Mutex

	dataForShell    chan []byte
	dataForTerminal chan []byte

	addr string
	log  slog.Logger
}

func (s *Server) handleConn(conn net.Conn, readFrom, writeTo chan []byte) {
	ctx, cancel := context.WithCancel(context.TODO())

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				s.log.Error(errors.Wrap(err, "reading from connection failed"))
				break
			}
			writeTo <- buf[:n]
		}
		cancel()
	}()

	go func() {
		defer wg.Done()

	outer:
		for {
			select {
			case <-ctx.Done():
				break outer
			case data := <-readFrom:
				_, err := conn.Write(data)
				if err != nil {
					s.log.Error(errors.Wrap(err, "writing to connection failed"))
				}
			}
		}
	}()

	wg.Wait()
}

func (s *Server) handleTermConn(conn net.Conn) {
	s.handleConn(conn, s.dataForShell, s.dataForTerminal)
}

func (s *Server) handleShellConn(conn net.Conn) {
	s.handleConn(conn, s.dataForTerminal, s.dataForShell)
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()

	connLog := s.log.With("remote.addr", conn.RemoteAddr().String())

	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		connLog.Error(errors.Wrap(err, "reading from connection failed"))
		return
	}

	var isShellConn bool
	switch buf[0] {
	case 0x01:
		isShellConn = true
	case 0x02:
		isShellConn = false
	default:
		fmt.Fprintf(os.Stderr, "unexpected data: %v", buf[0])
		return
	}

	func() {
		s.mux.Lock()
		defer s.mux.Unlock()
		if isShellConn {
			connLog.Debug("received shell connection")
			if s.shellConn != nil {
				connLog.Debug("closing existing shell connection")
				s.shellConn.Close()
				s.shellConn = nil
			}
			s.shellConn = conn
		} else {
			connLog.Debug("received term connection")
			if s.terminalConn != nil {
				connLog.Debug("closing existing term connection")
				s.terminalConn.Close()
				s.terminalConn = nil
			}
			s.terminalConn = conn
		}
	}()

	if isShellConn {
		s.handleShellConn(conn)
	} else {
		s.handleTermConn(conn)
	}
}

// Start starts the debug server listener
func (s *Server) Start() error {
	s.log.With("addr", s.addr).Debug("starting debugger server")
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting: %v", err.Error())
			break
		}
		go s.handleRequest(conn)
	}
	return nil
}

// NewServer returns a new server
func NewServer(addr string, log slog.Logger) *Server {
	return &Server{
		addr:            addr,
		dataForShell:    make(chan []byte, 100),
		dataForTerminal: make(chan []byte, 100),
		log:             log,
	}
}
