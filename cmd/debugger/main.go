package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

var (
	// Version is the version of the debugger
	Version string
)

func interactiveMode(remoteConsoleAddr string) {
	c, err := net.Dial("tcp", remoteConsoleAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to earth remote console: %v\n", err)
		return
	}

	cmd := exec.Command("sh", "-i")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to attach to shell stdin: %v\n", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to exec debugger shell: %v\n", err)
		return
	}

	b := make([]byte, 256)
	for {
		n, err := c.Read(b)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("err: %v\n", err)
			}
			break
		}
		stdin.Write(b[:n])
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to wait for debugger shell: %v\n", err)
		return
	}
}

func getRemoteDebuggerAddr() string {
	remoteConsoleAddr, err := ioutil.ReadFile("/run/secrets/earthly_remote_console_addr")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to read earthly_remote_console_addr: %v", err)
	}
	return string(remoteConsoleAddr)
}

func main() {
	args := os.Args[1:]

	remoteConsoleAddr := getRemoteDebuggerAddr()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			fmt.Fprintf(os.Stderr, "Command %v failed with exit code %d", args, exitCode)
		} else {
			fmt.Fprintf(os.Stderr, "Command %v failed with unexpected execution error %v", args, err)
		}

		if remoteConsoleAddr != "" {
			interactiveMode(remoteConsoleAddr)
		}

		os.Exit(exitCode)
	}
}
