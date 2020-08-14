package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/earthly/earthly/debugger/common"
)

func TestServer(t *testing.T) {
	addr := "127.0.0.1:9834"
	s := NewServer(addr)
	go s.Start()

	time.Sleep(10 * time.Millisecond)

	// first open terminal
	termConn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}
	_, err = termConn.Write([]byte{common.TermID})
	if err != nil {
		t.Error(err)
	}

	// then the shell terminal
	shellConn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}

	_, err = shellConn.Write([]byte{common.ShellID})
	if err != nil {
		t.Error(err)
	}

	inputStr := "hello world"

	// send data from shell to term
	_, err = shellConn.Write([]byte(inputStr))
	if err != nil {
		t.Error(err)
	}

	buf := make([]byte, 100)
	n, err := termConn.Read(buf)
	if err != nil {
		t.Error(err)
	}
	outputStr := string(buf[:n])

	if inputStr != outputStr {
		t.Error(fmt.Sprintf("want %v; got %v", inputStr, outputStr))
	}

}
