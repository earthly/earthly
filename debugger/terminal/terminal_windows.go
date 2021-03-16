// +build windows

package terminal

import (
	"github.com/pkg/errors"
)

func ConnectTerm(ctx context.Context, addr string) error {
	return errors.New("Debugger not supported on Windows yet")
}
