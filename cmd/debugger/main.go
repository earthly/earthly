package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/earthly/earthly/debugger/common"

	"github.com/creack/pty"
	"github.com/hashicorp/yamux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// Version is the version of the debugger
	Version string

	// GitSha is the git sha used to build the debugger
	GitSha string

	// ErrNoShellFound occurs when the container has no shell
	ErrNoShellFound = fmt.Errorf("no shell found")
)

func getShellPath() (string, bool) {
	for _, sh := range []string{
		"bash", "ksh", "zsh", "sh",
	} {
		if path, err := exec.LookPath(sh); err == nil {
			return path, true
		}
	}
	return "", false
}

func interactiveMode(ctx context.Context, addr string, log *logrus.Logger) error {
	var conn net.Conn
	var err error

	log.WithField("addr", addr).Debug("attempting connection")
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to connect to %v", addr))
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error closing: %v\n", err)
		}
	}()

	session, err := yamux.Client(conn, nil)
	if err != nil {
		return errors.Wrap(err, "failed creating yamux client")
	}

	ptyStream, err := session.Open()
	if err != nil {
		return errors.Wrap(err, "failed openning ptyStream session")
	}
	ptyStream.Write([]byte{common.PtyStream})

	winChangeStream, err := session.Open()
	if err != nil {
		return errors.Wrap(err, "failed openning winChangeStream session")
	}
	winChangeStream.Write([]byte{common.WinChangeStream})

	shellPath, ok := getShellPath()
	if !ok {
		return ErrNoShellFound
	}
	c := exec.CommandContext(ctx, shellPath)

	ptmx, err := pty.Start(c)
	if err != nil {
		return errors.Wrap(err, "failed to start pty")
	}
	defer func() { _ = ptmx.Close() }() // Best effort.

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		_, err := io.Copy(ptmx, ptyStream)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "failed copying pty to ptyStream: %v\n", err)
		}
		cancel()
	}()
	go func() {
		_, err := io.Copy(ptyStream, ptmx)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "failed copying pty to ptyStream: %v\n", err)
		}
		cancel()
	}()
	go func() {
		_ = c.Wait()
		cancel()
	}()

	go func() {
		for {
			data, err := common.ReadUint16PrefixedData(winChangeStream)
			if err == io.EOF {
				return
			} else if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read data: %v\n", err)
				break
			}

			var size pty.Winsize
			err = json.Unmarshal(data, &size)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to unmarshal data: %v\n", err)
				break
			}

			err = pty.Setsize(ptmx, &size)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to set window size: %v\n", err)
				break
			}

		}
		cancel()
	}()

	<-ctx.Done()

	fmt.Fprintf(os.Stderr, "exiting interactive debugger shell\n")
	return nil
}

func getSettings(path string) (*common.DebuggerSettings, error) {
	s, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read %s", path))
	}
	var data common.DebuggerSettings
	err = json.Unmarshal(s, &data)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to unmarshal %s", path))
	}
	return &data, nil
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "No command to run\n")
		os.Exit(1)
	}

	if args[0] == "--version" {
		fmt.Printf("version: %v-%v\n", Version, GitSha)
		return
	}

	ctx := context.Background()

	var log = logrus.New()

	debuggerSettings, err := getSettings("/run/secrets/earthly_debugger_settings")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read settings: %v\n", debuggerSettings)
		os.Exit(1)
	}

	if debuggerSettings.DebugLevelLogging {
		log.SetLevel(logrus.DebugLevel)
	}

	log.WithField("command", args).Debug("running command")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			fmt.Fprintf(os.Stderr, "Command %v failed with exit code %d\n", args, exitCode)
		} else {
			fmt.Fprintf(os.Stderr, "Command %v failed with unexpected execution error %v\n", args, err)
		}

		if debuggerSettings.Enabled {
			// Sometimes the interactive shell doesn't correctly get a newline
			// Take a brief pause and issue a new line as a work around.
			time.Sleep(time.Millisecond * 5)
			fmt.Printf("\n")
			interactiveMode(ctx, debuggerSettings.RemoteConsoleAddr, log)
		}

		// ensure that this always exits with an error status; otherwise it will be cached by earthly
		if exitCode == 0 {
			exitCode = 1
		}
		os.Exit(exitCode)
	}
}
