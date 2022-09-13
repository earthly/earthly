package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/slog"

	"github.com/alessio/shellescape"
	"github.com/creack/pty"
	"github.com/fatih/color"
	multierror "github.com/hashicorp/go-multierror"
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
		f.Write([]byte(cmd + "\n"))
	}
	return result
}

func interactiveMode(ctx context.Context, remoteConsoleAddr string, cmdBuilder func() (*exec.Cmd, error)) error {
	log := slog.GetLogger(ctx)

	conn, err := net.Dial("unix", remoteConsoleAddr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to remote debugger")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "error closing"))
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

	ptmx, err := pty.Start(c)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to start pty"))
		return err
	}
	defer func() { _ = ptmx.Close() }() // Best effort.

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			connDataType, data, err := common.ReadDataPacket(conn)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to read data from conn"))
				break
			}
			switch connDataType {
			case common.PtyData:
				handlePtyData(ptmx, data)
			case common.WinSizeData:
				handleWinChangeData(ptmx, data)
			default:
				log.With("datatype", connDataType).Warning("unhandled data type")
			}
		}
		cancel()
	}()
	go func() {
		initialData := true
		for {
			buf := make([]byte, 100)
			n, err := ptmx.Read(buf)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to read from ptmx"))
				break
			}
			buf = buf[:n]
			if initialData {
				buf = append([]byte("\r\n"), buf...)
				initialData = false
			}
			common.WriteDataPacket(conn, common.PtyData, buf)

		}
		cancel()
	}()

	var waitErr error
	var waitErrSet bool
	go func() {
		waitErr = c.Wait()
		waitErrSet = true
		cancel()
	}()

	<-ctx.Done()

	common.WriteDataPacket(conn, common.EndShellSession, nil)

	if !waitErrSet {
		return errInteractiveModeWaitFailed
	}
	return waitErr
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
		fmt.Printf("version: %v-%v\n", Version, GitSha)
		return
	}

	forceInteractive := false
	if args[0] == "--force" {
		args = args[1:]
		forceInteractive = true
	}

	conslogger := conslogging.Current(conslogging.ForceColor, conslogging.NoPadding, conslogging.Info)
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

	log := slog.GetLogger(ctx)

	if forceInteractive {
		quotedCmd := shellescape.QuoteCommand(args)

		conslogger.PrintBar(color.New(color.FgHiMagenta), "🌍 Earthly Build Interactive Session", quotedCmd)

		// Sometimes the interactive shell doesn't correctly get a newline
		// Take a brief pause and issue a new line as a work around.
		time.Sleep(time.Millisecond * 5)

		err := os.Setenv("TERM", debuggerSettings.Term)
		if err != nil {
			conslogger.Warnf("Failed to set term: %v", err)
		}

		cmdBuilder := func() (*exec.Cmd, error) {
			return exec.Command(args[0], args[1:]...), nil
		}

		exitCode := 0
		err = interactiveMode(ctx, debuggerSettings.SocketPath, cmdBuilder)
		if err != nil {
			log.Error(err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 127
			}
		}

		conslogger.PrintBar(color.New(color.FgHiMagenta), " End Interactive Session ", "")

		os.Exit(exitCode)
	}

	log.With("command", args).With("version", Version).Debug("running command")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {

		quotedCmd := shellescape.QuoteCommand(args)

		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			conslogger.Warnf("Command %s failed with exit code %d\n", quotedCmd, exitCode)
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
				conslogger.Warnf("Failed to set term: %v", err)
			}

			cmdBuilder := func() (*exec.Cmd, error) {
				_ = populateShellHistory(quotedCmd) // best effort

				shellPath, ok := getShellPath()
				if !ok {
					return nil, ErrNoShellFound
				}
				log.With("shell", shellPath).Debug("found shell")
				return exec.Command(shellPath), nil
			}

			err = interactiveMode(ctx, debuggerSettings.SocketPath, cmdBuilder)
			if err != nil {
				log.Error(err)
			}
		}
		// ensure that this always exits with an error status; otherwise it will be cached by earthly
		if exitCode == 0 {
			exitCode = 1
		}
		os.Exit(exitCode)
	}
}
