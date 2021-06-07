// +build windows

package terminal

import (
	"context"

	"github.com/earthly/earthly/conslogging"

	"github.com/pkg/errors"
)

func ConnectTerm(ctx context.Context, addr string, console conslogging.ConsoleLogger) error {
	return errors.New("debugger not supported on Windows yet")
}
