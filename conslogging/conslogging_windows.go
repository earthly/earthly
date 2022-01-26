//go:build windows
// +build windows

package conslogging

import (
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int, verbose bool) ConsoleLogger {
	return ConsoleLogger{
		consoleErrW:    colorable.NewColorable(os.Stderr),
		errW:           colorable.NewColorable(os.Stderr),
		colorMode:      colorMode,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		prefixPadding:  prefixPadding,
		mu:             &currentConsoleMutex,
		verbose:        verbose,
	}
}
