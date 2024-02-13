//go:build windows
// +build windows

package terminal

import (
	"context"
	"io"

	"github.com/earthly/earthly/conslogging"

	"github.com/pkg/errors"
)

func ConnectTerm(ctx context.Context, addr io.ReadWriteCloser, console conslogging.ConsoleLogger) error {
	return errors.New("debugger not supported on Windows yet")
}
