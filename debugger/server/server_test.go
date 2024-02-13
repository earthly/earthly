package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/slog"

	"github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.TODO()
	log := slog.GetLogger(ctx).With("test.name", t.Name())

	addr := "127.0.0.1:9834"
	s := NewServer(addr, log)
	go func() {
		err := s.Start()
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	const numRetries = 3
	attempts := 0
	var termConn net.Conn
	var err error

	for attempts < numRetries {
		termConn, err = net.Dial("tcp", addr)

		if err != nil {
			// Retry since the connection is rejected sometimes.
			fmt.Printf("Dial failed. Attempt: %v/%v, Error: %s", attempts, numRetries, err.Error())
			time.Sleep(time.Second)

			attempts++
			err = nil
			continue
		}
		break
	}
	if attempts >= numRetries {
		t.Fatal("Retries exhausted")
	}

	_, err = termConn.Write([]byte{common.TermID})
	if err != nil {
		t.Fatal(err)
	}

	// then the shell terminal
	shellConn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	_, err = shellConn.Write([]byte{common.ShellID})
	if err != nil {
		t.Fatal(err)
	}

	inputStr := "hello world"

	// send data from shell to term
	_, err = shellConn.Write([]byte(inputStr))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 100)
	n, err := termConn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	outputStr := string(buf[:n])

	if inputStr != outputStr {
		t.Fatalf("want %v; got %v", inputStr, outputStr)
	}
}
