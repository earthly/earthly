// +build !windows

package conslogging

import (
	"os"

	"github.com/fatih/color"
)

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int) ConsoleLogger {
	return ConsoleLogger{
		outW:           os.Stderr, // So logs dont sully any intended outputs of commands.
		errW:           os.Stderr,
		colorMode:      colorMode,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		prefixPadding:  prefixPadding,
		mu:             &currentConsoleMutex,
	}
}
