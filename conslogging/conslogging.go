package conslogging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/fatih/color"
)

var currentConsoleMutex sync.Mutex

// ConsoleLogger is a writer for consoles.
type ConsoleLogger struct {
	prefix        string
	disableColors bool
	isCached      bool

	// The following are shared between instances and are protected by the mutex.
	mu             *sync.Mutex
	prefixColors   map[string]*color.Color
	nextColorIndex *int
	w              io.Writer
	trailingLine   bool
}

// Current returns the current console.
func Current(disableColors bool) ConsoleLogger {
	return ConsoleLogger{
		w:              os.Stdout,
		disableColors:  disableColors || color.NoColor,
		prefixColors:   make(map[string]*color.Color),
		nextColorIndex: new(int),
		mu:             &currentConsoleMutex,
	}
}

func (cl ConsoleLogger) clone() ConsoleLogger {
	return ConsoleLogger{
		w:              cl.w,
		prefix:         cl.prefix,
		isCached:       cl.isCached,
		prefixColors:   cl.prefixColors,
		disableColors:  cl.disableColors,
		nextColorIndex: cl.nextColorIndex,
		mu:             cl.mu,
	}
}

// WithPrefix returns a ConsoleLogger with a prefix added.
func (cl ConsoleLogger) WithPrefix(prefix string) ConsoleLogger {
	ret := cl.clone()
	ret.prefix = prefix
	return ret
}

// Prefix returns the console's prefix.
func (cl ConsoleLogger) Prefix() string {
	return cl.prefix
}

// WithCached returns a ConsoleLogger with isCached flag set accordingly.
func (cl ConsoleLogger) WithCached(isCached bool) ConsoleLogger {
	ret := cl.clone()
	ret.isCached = isCached
	return ret
}

// PrintSuccess prints the success message.
func (cl ConsoleLogger) PrintSuccess() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	successColor.Fprintf(cl.w, "=========================== SUCCESS ===========================\n")
}

// Printf prints formatted text to the console.
func (cl ConsoleLogger) Printf(format string, args ...interface{}) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")
	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		cl.w.Write([]byte(line))
		cl.w.Write([]byte("\n"))
	}
}

// PrintBytes prints bytes directly to the console.
func (cl ConsoleLogger) PrintBytes(data []byte) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	output := []byte{}
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		ch := data[:size]
		data = data[size:]
		switch r {
		case '\r':
			output = append(output, ch...)
			cl.trailingLine = false
		case '\n':
			output = append(output, ch...)
			cl.trailingLine = false
		default:
			if !cl.trailingLine {
				if len(output) > 0 {
					cl.w.Write(output)
					output = []byte{}
				}
				cl.printPrefix()
				cl.trailingLine = true
			}
			output = append(output, ch...)
		}
	}
	if len(output) > 0 {
		cl.w.Write(output)
		output = []byte{}
	}
}

func (cl ConsoleLogger) printPrefix() {
	// Assumes mu locked.
	if cl.prefix == "" {
		return
	}
	c := noColor
	if !cl.disableColors {
		var found bool
		c, found = cl.prefixColors[cl.prefix]
		if !found {
			c = availablePrefixColors[*cl.nextColorIndex]
			cl.prefixColors[cl.prefix] = c
			*cl.nextColorIndex = (*cl.nextColorIndex + 1) % len(availablePrefixColors)
		}
	}
	c.Fprintf(cl.w, "%s", cl.prefix)
	cl.w.Write([]byte(" | "))
	if cl.isCached {
		cl.w.Write([]byte("*"))
		cachedColor.Fprintf(cl.w, "cached")
		cl.w.Write([]byte("* "))
	}
}
