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
