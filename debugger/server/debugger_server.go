package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/earthly/earthly/conslogging"
)

// DebugServer is provides a console which reverse shells connect to
type DebugServer struct {
	addr    string
	console conslogging.ConsoleLogger
}

// NewDebugServer creates a new debug server
func NewDebugServer(console conslogging.ConsoleLogger) *DebugServer {
	return &DebugServer{
		addr:    "127.0.0.1:8543", // TODO make this configurable (and support port 0 which auto assigns a free port)
		console: console,
	}
}

func (ds *DebugServer) handleRequest(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 256)
	for {
		n, err := os.Stdin.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("err: %v\n", err)
			}
			break
		}
		conn.Write(b[:n])
	}
}

// Start starts the debug server listener
func (ds *DebugServer) Start() (string, error) {
	l, err := net.Listen("tcp", ds.addr)
	if err != nil {
		return "", err
	}
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
			ds.handleRequest(conn)
		}
	}()
	return l.Addr().String(), nil
}
