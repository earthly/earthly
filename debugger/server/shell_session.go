package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"

	"github.com/earthly/earthly/debugger/common"

	"github.com/creack/pty"
	"github.com/hashicorp/yamux"
)

// session is used to track a single reverse shell's session
type session struct {
	yaSession *yamux.Session

	ctx    context.Context
	cancel context.CancelFunc

	ttyCon     net.Conn
	resizeConn net.Conn

	server *DebugServer
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
			go s.handlePtyStream(stream)

		case common.WinChangeStream:
			s.resizeConn = stream
			s.server.sigs <- syscall.SIGWINCH
		default:
			return fmt.Errorf("unsupported stream code %v", buf[0])
		}
	}
}

func (s *session) handlePtyStream(conn net.Conn) error {
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "failed copying stdout to ptyStream: %v\n", err)
		}
		s.cancel()
	}()
	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "failed copying ptyStream to stdin: %v\n", err)
		}
		s.cancel()
	}()

	<-s.ctx.Done()
	return nil
}
