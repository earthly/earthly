package main

import (
	"fmt"
	"io"
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

func main() {
	remoteConsoleAddr := os.Getenv("EARTHLY_REMOTE_CONSOLE_ADDR")
	if remoteConsoleAddr == "" {
		remoteConsoleAddr = "127.0.0.1:8543"
	}

	args := os.Args[1:]

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %v; entering debug mode (debugger version %v)\n", err, Version)
		interactiveMode(remoteConsoleAddr)
		os.Exit(1)
	}
}
