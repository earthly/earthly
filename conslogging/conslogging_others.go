//go:build !windows
// +build !windows

package conslogging

import (
	"io"
	"os"
)

func getCompatibleStderr() io.Writer {
	return os.Stderr
}
