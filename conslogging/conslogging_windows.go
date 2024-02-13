//go:build windows
// +build windows

package conslogging

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
)

func getCompatibleStderr() io.Writer {
	return colorable.NewColorable(os.Stderr)
}
