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
	prefix string
	// salt is a salt used for color consistency
	// (the same salt will get the same color).
	salt          string
	disableColors bool
	isCached      bool
	isFailed      bool

	// The following are shared between instances and are protected by the mutex.
	mu             *sync.Mutex
	saltColors     map[string]*color.Color
	nextColorIndex *int
	w              io.Writer
	trailingLine   bool
}

// Current returns the current console.
func Current(disableColors bool) ConsoleLogger {
	return ConsoleLogger{
		w:              os.Stdout,
		disableColors:  disableColors || color.NoColor,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		mu:             &currentConsoleMutex,
	}
}

func (cl ConsoleLogger) clone() ConsoleLogger {
	return ConsoleLogger{
		w:              cl.w,
		prefix:         cl.prefix,
		salt:           cl.salt,
		isCached:       cl.isCached,
		isFailed:       cl.isFailed,
		saltColors:     cl.saltColors,
		disableColors:  cl.disableColors,
		nextColorIndex: cl.nextColorIndex,
		mu:             cl.mu,
	}
}

// WithPrefix returns a ConsoleLogger with a prefix added.
func (cl ConsoleLogger) WithPrefix(prefix string) ConsoleLogger {
	ret := cl.clone()
	ret.prefix = prefix
	ret.salt = prefix
	return ret
}

// WithPrefixAndSalt returns a ConsoleLogger with a prefix and a seed added.
func (cl ConsoleLogger) WithPrefixAndSalt(prefix string, salt string) ConsoleLogger {
	ret := cl.clone()
	ret.prefix = prefix
	ret.salt = salt
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

// WithFailed returns a ConsoleLogger with isFailed flag set accordingly.
func (cl ConsoleLogger) WithFailed(isFailed bool) ConsoleLogger {
	ret := cl.clone()
	ret.isFailed = isFailed
	return ret
}

// PrintSuccess prints the success message.
func (cl ConsoleLogger) PrintSuccess() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	successColor.Fprintf(cl.w, "=========================== SUCCESS ===========================\n")
}

// PrintFailure prints the failure message.
func (cl ConsoleLogger) PrintFailure() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	warnColor.Fprintf(cl.w, "=========================== FAILURE ===========================\n")
}

// Warnf prints a warning message in red
func (cl ConsoleLogger) Warnf(format string, args ...interface{}) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	c := noColor
	if !cl.disableColors {
		c = warnColor
	}

	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")

	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		c.Fprintf(cl.w, "%s\n", line)
	}
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

	output := make([]byte, 0, len(data))
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
					output = output[:0]
				}
				cl.printPrefix()
				cl.trailingLine = true
			}
			output = append(output, ch...)
		}
	}
	if len(output) > 0 {
		cl.w.Write(output)
		output = output[:0]
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
		c, found = cl.saltColors[cl.salt]
		if !found {
			c = availablePrefixColors[*cl.nextColorIndex]
			cl.saltColors[cl.salt] = c
			*cl.nextColorIndex = (*cl.nextColorIndex + 1) % len(availablePrefixColors)
		}
	}
	c.Fprintf(cl.w, "%s", cl.prefix)
	if cl.isFailed {
		cl.w.Write([]byte(" *"))
		warnColor.Fprintf(cl.w, "failed")
		cl.w.Write([]byte("*"))
	}
	cl.w.Write([]byte(" | "))
	if cl.isCached {
		cl.w.Write([]byte("*"))
		cachedColor.Fprintf(cl.w, "cached")
		cl.w.Write([]byte("* "))
	}
}
