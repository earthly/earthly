package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/logging"

	"github.com/creack/pty"
	"github.com/fatih/color"
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

const remoteConsoleAddr = "127.0.0.1:5000"

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

func interactiveMode(ctx context.Context, remoteConsoleAddr string) error {
	log := logging.GetLogger(ctx)

	conn, err := net.Dial("tcp", remoteConsoleAddr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to remote debugger")
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "error closing"))
		}
	}()

	_, err = conn.Write([]byte{common.ShellID})
	if err != nil {
		return err
	}

	err = common.WriteDataPacket(conn, common.StartShellSession, nil)
	if err != nil {
		return err
	}

	shellPath, ok := getShellPath()
	if !ok {
		return ErrNoShellFound
	}
	log.With("shell", shellPath).Debug("found shell")
	c := exec.Command(shellPath)

	ptmx, err := pty.Start(c)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to start pty"))
		return err
	}
	defer func() { _ = ptmx.Close() }() // Best effort.

	ctx, cancel := context.WithCancel(context.Background())

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
		for {
			buf := make([]byte, 100)
			n, err := ptmx.Read(buf)
			if err != nil {
				log.Error(errors.Wrap(err, "failed to read from ptmx"))
				break
			}
			buf = buf[:n]
			common.WriteDataPacket(conn, common.PtyData, buf)

		}
		cancel()
	}()

	go func() {
		c.Wait()
		cancel()
	}()

	<-ctx.Done()

	common.WriteDataPacket(conn, common.EndShellSession, nil)

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

	if args[0] == "--version" {
		fmt.Printf("version: %v-%v\n", Version, GitSha)
		return
	}

	conslogger := conslogging.Current(conslogging.ForceColor)
	color.NoColor = false

	debuggerSettings, err := getSettings(fmt.Sprintf("/run/secrets/%s", common.DebuggerSettingsSecretsKey))
	if err != nil {
		conslogger.Warnf("failed to read settings: %v\n", debuggerSettings)
		os.Exit(1)
	}

	if debuggerSettings.DebugLevelLogging {
		logrus.SetLevel(logrus.DebugLevel)
	}

	ctx := context.Background()

	log := logging.GetLogger(ctx)

	log.With("command", args).With("version", Version).Debug("running command")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			conslogger.Warnf("Command %v failed with exit code %d\n", strings.Join(args, " "), exitCode)
		} else {
			conslogger.Warnf("Command %v failed with unexpected execution error %v\n", strings.Join(args, " "), err)
		}

		if debuggerSettings.Enabled {
			c := color.New(color.FgYellow)
			c.Println("Entering interactive debugger (**Warning: only a single debugger per host is supported**)")

			// Sometimes the interactive shell doesn't correctly get a newline
			// Take a brief pause and issue a new line as a work around.
			time.Sleep(time.Millisecond * 5)

			err := interactiveMode(ctx, remoteConsoleAddr)
			if err != nil {
				log.Error(err)
			}

			// ensure that this always exits with an error status; otherwise it will be cached by earthly
			if exitCode == 0 {
				exitCode = 1
			}
			os.Exit(exitCode)
		}
	}

}
