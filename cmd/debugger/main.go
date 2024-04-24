package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/slog"

	"github.com/alessio/shellescape"
	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// Version is the version of the debugger
	Version string

	// GitSha is the git sha used to build the debugger
	GitSha string

	// ErrNoShellFound occurs when the container has no shell
	ErrNoShellFound = errors.New("no shell found")

	errInteractiveModeWaitFailed = errors.New("interactive mode wait failed")
)

type waitError struct {
	set bool
	err error
}

func newWaitError(err error, set bool) *waitError {
	return &waitError{
		set: set,
		err: err,
	}
}

func getShellPath() (string, bool) {
	for _, sh := range []string{
		"bash", "ksh", "zsh", "ash", "sh",
	} {
		if path, err := exec.LookPath(sh); err == nil {
			return path, true
		}
	}
	return "", false
}

func handlePtyData(ptmx *os.File, data []byte) error {
	_, err := ptmx.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write to ptmx")
	}
	return nil
}

func handleWinChangeData(ptmx *os.File, data []byte) error {
	var size pty.Winsize
	err := json.Unmarshal(data, &size)
	if err != nil {
		return errors.Wrap(err, "failed unmarshal data")
	}

	err = pty.Setsize(ptmx, &size)
	if err != nil {
		return errors.Wrap(err, "failed to set window size")
	}
	return nil
}

func populateShellHistory(cmd string) error {
	var result error
	for _, f := range []string{
		"/root/.ash_history",
		"/root/.bash_history",
	} {

		f, err := os.Create(f)
		if err != nil {
			result = multierror.Append(result, err)
		}
		defer f.Close()
		_, err = f.Write([]byte(cmd + "\n"))
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

func sendFile(ctx context.Context, sockAddr, src, dst string) error {
	log := slog.GetLogger(ctx)

	conn, err := net.Dial("unix", sockAddr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to remote debugger")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "earthly debugger: error closing"))
		}
	}()

	// send a protocol version
	err = common.WriteDataPacket(conn, 0x02, nil)
	if err != nil {
		return err
	}

	err = common.WriteUint16PrefixedData(conn, []byte(dst))
	if err != nil {
		return err
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	b := make([]byte, 0, math.MaxUint16)
	for {
		n, err := r.Read(b[:cap(b)])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = common.WriteUint16PrefixedData(conn, b[:n])
		if err != nil {
			return err
		}
	}

	// send end of file packet
	return common.WriteUint16PrefixedData(conn, nil)
}

func interactiveMode(ctx context.Context, remoteConsoleAddr string, cmdBuilder func() (*exec.Cmd, error), conslogger conslogging.ConsoleLogger) error {
	log := slog.GetLogger(ctx)

	conn, err := net.Dial("unix", remoteConsoleAddr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to remote debugger")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "earthly debugger: error closing"))
		}
	}()

	err = common.WriteDataPacket(conn, common.StartShellSession, nil)
	if err != nil {
		return err
	}

	c, err := cmdBuilder()
	if err != nil {
		return err
	}

	// once the command completes, waitErr will be set to true, indicating the command finished (or failed)
	// if it is not set, then that means the interactive debugger has exited before the wrapped command has finished.
	waitErr := atomic.Pointer[waitError]{}
	waitErr.Store(newWaitError(nil, false))

	hasCommandFinished := func() bool {
		// give c.Wait() time to acquire the lock first, to detect if the command closed as expected
		time.Sleep(time.Millisecond * 10)
		return waitErr.Load().set
	}

	logErrorIfNonCleanExit := func(err error) {
		if hasCommandFinished() {
			return
		}
		conslogger.Warnf("%v\n", errors.Wrap(err, "failed to start pty"))
	}

	ptmx, err := pty.Start(c)
	if err != nil {
		conslogger.Warnf("%v\n", errors.Wrap(err, "failed to start pty"))
		return err
	}
	defer func() { _ = ptmx.Close() }() // Best effort.

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		for {
			connDataType, data, err := common.ReadDataPacket(conn)
			if err != nil {
				logErrorIfNonCleanExit(errors.Wrap(err, "failed to read data from conn"))
				return
			}
			switch connDataType {
			case common.PtyData:
				err = handlePtyData(ptmx, data)
				if err != nil {
					logErrorIfNonCleanExit(errors.Wrap(err, "failed to handle pty data"))
					return
				}
			case common.WinSizeData:
				err = handleWinChangeData(ptmx, data)
				if err != nil {
					logErrorIfNonCleanExit(errors.Wrap(err, "failed to handle win change data"))
					return
				}
			default:
				conslogger.Warnf("unhandled data type (%v)\n", connDataType)
			}
		}
	}()
	go func() {
		defer cancel()
		initialData := true
		for {
			buf := make([]byte, 100)
			n, err := ptmx.Read(buf)
			if err != nil {
				logErrorIfNonCleanExit(errors.Wrap(err, "failed to read from ptmx"))
				return
			}
			buf = buf[:n]
			if initialData {
				buf = append([]byte("\r\n"), buf...)
				initialData = false
			}
			err = common.WriteDataPacket(conn, common.PtyData, buf)
			if err != nil {
				logErrorIfNonCleanExit(errors.Wrap(err, "failed to write data to conn"))
				return
			}
		}
	}()

	go func() {
		err := c.Wait()
		waitErr.Store(newWaitError(err, true))
		cancel()
	}()

	<-ctx.Done()

	err = common.WriteDataPacket(conn, common.EndShellSession, nil)
	if err != nil {
		return errors.Wrap(err, "failed to send end shell session")
	}

	if !waitErr.Load().set {
		return errInteractiveModeWaitFailed
	}
	return waitErr.Load().err
}

func getSettings(path string) (*common.DebuggerSettings, error) {
	s, err := os.ReadFile(path)
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

	if args[0] == "--version" {
		return
	}

	forceInteractive := false
	if args[0] == "--force" {
		args = args[1:]
		forceInteractive = true
	}

	conslogger := conslogging.Current(conslogging.ForceColor, conslogging.NoPadding, conslogging.Info, false).
		WithPrefix("earthly debugger")

	color.NoColor = false

	debuggerSettings, err := getSettings(fmt.Sprintf("/run/secrets/%s", common.DebuggerSettingsSecretsKey))
	if err != nil {
		conslogger.Warnf("failed to read settings: %v\n", debuggerSettings)
		os.Exit(1)
	}

	if debuggerSettings.DebugLevelLogging {
		logrus.SetLevel(logrus.DebugLevel)
		conslogger = conslogger.WithLogLevel(conslogging.Verbose)
	}

	ctx := context.Background()

	if forceInteractive {
		quotedCmd := shellescape.QuoteCommand(args)

		conslogger.PrintBar(color.New(color.FgHiMagenta), "🌍 Earthly Build Interactive Session", quotedCmd)

		// Sometimes the interactive shell doesn't correctly get a newline
		// Take a brief pause and issue a new line as a workaround.
		time.Sleep(time.Millisecond * 5)

		err := os.Setenv("TERM", debuggerSettings.Term)
		if err != nil {
			conslogger.Warnf("Failed to set term: %v\n", err)
		}

		cmdBuilder := func() (*exec.Cmd, error) {
			return exec.Command(args[0], args[1:]...), nil
		}

		exitCode := 0
		err = interactiveMode(ctx, debuggerSettings.SocketPath, cmdBuilder, conslogger)
		if err != nil {
			conslogger.Warnf("%v\n", err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 127
			}
		}

		conslogger.PrintBar(color.New(color.FgHiMagenta), " End Interactive Session ", "")

		os.Exit(exitCode)
	}

	conslogger.VerbosePrintf("running command: (%s); version: %s\n", args, Version)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {

		quotedCmd := shellescape.QuoteCommand(args)

		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			if debuggerSettings.Enabled {
				conslogger.Warnf("Command %s failed with exit code %d\n", quotedCmd, exitCode)
			}
		} else {
			conslogger.Warnf("Command %s failed with unexpected execution error %v\n", quotedCmd, err)
		}

		if debuggerSettings.Enabled {
			c := color.New(color.FgYellow)
			c.Println("Entering interactive debugger")
			// Sometimes the interactive shell doesn't correctly get a newline
			// Take a brief pause and issue a new line as a work around.
			time.Sleep(time.Millisecond * 5)

			err := os.Setenv("TERM", debuggerSettings.Term)
			if err != nil {
				conslogger.Warnf("Failed to set term: %v\n", err)
			}

			cmdBuilder := func() (*exec.Cmd, error) {
				_ = populateShellHistory(quotedCmd) // best effort

				shellPath, ok := getShellPath()
				if !ok {
					return nil, ErrNoShellFound
				}
				conslogger.VerbosePrintf("found shell: (%s)\n", shellPath)
				return exec.Command(shellPath), nil
			}

			err = interactiveMode(ctx, debuggerSettings.SocketPath, cmdBuilder, conslogger)
			if err != nil {
				conslogger.Warnf("%v\n", err)
			}
		}

		for _, saveFile := range debuggerSettings.SaveFiles {
			err = sendFile(ctx, common.DefaultSaveFileSocketPath, saveFile.Src, saveFile.Dst)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) || !saveFile.IfExists {
					// treat it as a warning (we will exit due to RUN failure)
					conslogger.Warnf("failed to save %s: %s\n", saveFile.Src, err)
				}
			}
		}

		// ensure that this always exits with an error status; otherwise it will be cached by earthly
		if exitCode == 0 {
			exitCode = 1
		}
		os.Exit(exitCode)
	}
}
