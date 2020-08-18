package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
)

// Server provides a debugger server
type Server struct {
	shellConn    net.Conn
	terminalConn net.Conn
	mux          sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	dataForShell    chan []byte
	dataForTerminal chan []byte

	sigs chan os.Signal
	addr string
}

func (s *Server) handleConn(conn net.Conn, readFrom, writeTo chan []byte) {
	ctx, cancel := context.WithCancel(s.ctx)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("err %v\n", err)
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
					fmt.Printf("failed %v\n", err)
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

	buf := make([]byte, 1)
	conn.Read(buf)

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
			if s.shellConn != nil {
				s.shellConn.Close()
				s.shellConn = nil
			}
			s.shellConn = conn
		} else {
			if s.terminalConn != nil {
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
func NewServer(addr string) *Server {
	return &Server{
		addr:            addr,
		dataForShell:    make(chan []byte, 100),
		dataForTerminal: make(chan []byte, 100),
	}
}
