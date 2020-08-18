package terminal

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/logging"

	"github.com/creack/pty"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

// Terminal provides a terminal for a user to type commands into
// and to display the output of the shell.
// The terminal does not run commands, but rather passes them to the shell
// via the shell repeater.
type Terminal struct {
	conn net.Conn
}

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
func ConnectTerm(ctx context.Context, addr string) error {
	log := logging.GetLogger(ctx)

	var d net.Dialer

	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte{common.TermID})
	if err != nil {
		return errors.Wrap(err, "failed to write TermID connection")
	}

	sigs := make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGWINCH)

	writeCh := make(chan []byte, 10)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		var oldState *terminal.State

	outer:
		for {
			connDataType, data, err := common.ReadDataPacket(conn)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to read from connection"))
				break
			}
			switch connDataType {
			case common.StartShellSession:
				if oldState == nil {
					var err error
					oldState, err = terminal.MakeRaw(int(os.Stdin.Fd()))
					if err != nil {
						log.Error(errors.Wrap(err, "failed to initialize terminal in raw mode"))
						break outer
					}
				}
				sigs <- syscall.SIGWINCH
			case common.EndShellSession:
				if oldState != nil {
					err := terminal.Restore(int(os.Stdin.Fd()), oldState)
					if err != nil {
						log.Error(errors.Wrap(err, "failed to restore terminal mode"))
						break outer
					}
					oldState = nil
				}
			case common.PtyData:
				err := handlePtyData(data)
				if err != nil {
					log.Error(errors.Wrap(err, "failed to handle pty data"))
					break outer
				}
			default:
				log.With("datatype", connDataType).Warning("unhandled data type")
				break outer
			}
		}
		cancel()
	}()

	go func() {
	outer:
		for {
			select {
			case _ = <-sigs:
				if len(sigs) > 0 {
					continue
				}
				data, err := getWindowSizePayload()
				if err != nil {
					log.Error(errors.Wrap(err, "failed to restore terminal mode"))
					break outer
				}
				writeCh <- data
			}
		}
		cancel()
	}()

	go func() {
		for {
			buf := <-writeCh
			_, err := conn.Write(buf)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to restore terminal mode"))
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
				log.Error(errors.Wrap(err, "failed to read from stdin"))
				break
			}
			buf = buf[:n]
			buf2, err := common.SerializeDataPacket(common.PtyData, buf)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to serialize data"))
				break
			}

			writeCh <- buf2
		}
		cancel()
	}()

	<-ctx.Done()

	fmt.Fprintf(os.Stderr, "exiting interactive debugger shell\n")
	return nil
}
