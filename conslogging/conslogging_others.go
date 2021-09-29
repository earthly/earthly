// +build !windows

package conslogging

import (
	"os"

	"github.com/fatih/color"
)

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int, verbose bool) ConsoleLogger {
	return ConsoleLogger{
		errW:           os.Stderr,
		colorMode:      colorMode,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		prefixPadding:  prefixPadding,
		mu:             &currentConsoleMutex,
		verbose:        verbose,
	}
}
