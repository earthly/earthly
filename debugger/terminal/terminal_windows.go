// +build windows

package terminal

import (
	"context"

	"github.com/pkg/errors"
)

func ConnectTerm(ctx context.Context, addr string) error {
	return errors.New("debugger not supported on Windows yet")
}
