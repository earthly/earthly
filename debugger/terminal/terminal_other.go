//go:build !windows
// +build !windows

package terminal

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/debugger/common"

	"github.com/creack/pty"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

func handlePtyData(data []byte) error {
	_, err := os.Stdout.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write data to stdout")
	}
	return nil
}

func getWindowSizePayload() ([]byte, error) {
	size, err := pty.GetsizeFull(os.Stdin)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(size)
	if err != nil {
		return nil, err
	}
	return common.SerializeDataPacket(common.WinSizeData, b)
}

// ConnectTerm presents a terminal to the shell repeater
func ConnectTerm(ctx context.Context, conn io.ReadWriteCloser, console conslogging.ConsoleLogger) error {
	sigs := make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGWINCH)

	writeCh := make(chan []byte, 10)

	ctx, cancel := context.WithCancel(ctx)

	ts := &termState{}
	go func() {
	outer:
		for {
			connDataType, data, err := common.ReadDataPacket(conn)
			if err != nil {
				console.VerbosePrintf("ReadDataPacket failed: %s\n", err.Error())
				break
			}
			switch connDataType {
			case common.StartShellSession:
				console.VerbosePrintf("starting new interactive shell pseudo terminal\n")
				err := ts.makeRaw()
				if err != nil {
					console.VerbosePrintf("makeRaw failed: %s\n", err.Error())
					break outer
				}
				sigs <- syscall.SIGWINCH
			case common.EndShellSession:
				err := ts.restore()
				if err != nil {
					console.VerbosePrintf("restore failed: %s\n", err.Error())
					break outer
				}
			case common.PtyData:
				err := handlePtyData(data)
				if err != nil {
					console.VerbosePrintf("handlePtyData failed: %s\n", err.Error())
					break outer
				}
			default:
				console.VerbosePrintf("unhandled terminal data type: %q\n", connDataType)
				break outer
			}
		}
		cancel()
	}()

	go func() {
		for range sigs {
			if len(sigs) > 0 {
				continue
			}
			data, err := getWindowSizePayload()
			if err != nil {
				console.VerbosePrintf("failed to get window size payload: %s\n", err.Error())
				break
			}
			writeCh <- data
		}
		cancel()
	}()

	go func() {
		for {
			buf := <-writeCh
			_, err := conn.Write(buf)
			if err != nil {
				console.VerbosePrintf("failed to send term data to shell: %s\n", err.Error())
				break
			}
		}
		cancel()
	}()
	go func() {
		for {
			buf := make([]byte, 100)
			n, err := os.Stdin.Read(buf)
			if err != nil {
				console.VerbosePrintf("failed to read from stdin: %s\n", err.Error())
				break
			}
			buf = buf[:n]
			buf2, err := common.SerializeDataPacket(common.PtyData, buf)
			if err != nil {
				console.VerbosePrintf("failed to serialize data: %s\n", err.Error())
				break
			}

			writeCh <- buf2
		}
		cancel()
	}()

	<-ctx.Done()

	console.VerbosePrintf("exiting interactive debugger shell\n")
	err := ts.restore()
	if err != nil {
		return err
	}
	return nil
}

type termState struct {
	oldState *term.State
	mu       sync.Mutex
}

func (ts *termState) makeRaw() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.oldState == nil {
		var err error
		ts.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return errors.Wrap(err, "failed to initialize terminal in raw mode")
		}
	}
	return nil
}

func (ts *termState) restore() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.oldState != nil {
		err := term.Restore(int(os.Stdin.Fd()), ts.oldState)
		if err != nil {
			return errors.Wrap(err, "failed to restore terminal mode")
		}
		ts.oldState = nil
	}
	return nil
}
