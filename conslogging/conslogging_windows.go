// +build windows

package conslogging

import (
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int) ConsoleLogger {
	return ConsoleLogger{
		outW:           colorable.NewColorable(os.Stderr), // So logs dont sully any intended outputs of commands.
		errW:           colorable.NewColorable(os.Stderr),
		colorMode:      colorMode,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		prefixPadding:  prefixPadding,
		mu:             &currentConsoleMutex,
	}
}
