package conslogging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/cheggaaa/pb/v3"

	"github.com/fatih/color"
)

// ColorMode is the mode in which colors are represented in the output.
type ColorMode int

const (
	// AutoColor automatically detects the presence of a TTY to decide if
	// color should be used.
	AutoColor ColorMode = iota
	// NoColor disables use of color.
	NoColor
	// ForceColor forces use of color.
	ForceColor
)

var currentConsoleMutex sync.Mutex

// ConsoleLogger is a writer for consoles.
type ConsoleLogger struct {
	prefix string
	// salt is a salt used for color consistency
	// (the same salt will get the same color).
	salt      string
	colorMode ColorMode
	isCached  bool
	isFailed  bool

	// The following are shared between instances and are protected by the mutex.
	mu             *sync.Mutex
	saltColors     map[string]*color.Color
	nextColorIndex *int
	w              io.Writer
	trailingLine   bool

	// Progress bar variables
	progressBar   *pb.ProgressBar
	progressIsSet bool
}

// Current returns the current console.
func Current(colorMode ColorMode) ConsoleLogger {
	return ConsoleLogger{
		w:              os.Stdout,
		colorMode:      colorMode,
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
		colorMode:      cl.colorMode,
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
	cl.color(successColor).Fprintf(cl.w, "=========================== SUCCESS ===========================\n")
	cl.mu.Unlock()
}

// PrintFailure prints the failure message.
func (cl ConsoleLogger) PrintFailure() {
	cl.mu.Lock()
	cl.color(warnColor).Fprintf(cl.w, "=========================== FAILURE ===========================\n")
	cl.mu.Unlock()
}

// Warnf prints a warning message in red
func (cl ConsoleLogger) Warnf(format string, args ...interface{}) {
	cl.mu.Lock()

	c := cl.color(warnColor)
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")

	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		c.Fprintf(cl.w, "%s\n", line)
	}
	cl.mu.Unlock()
}

// Printf prints formatted text to the console.
func (cl ConsoleLogger) Printf(format string, args ...interface{}) {
	cl.mu.Lock()
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")
	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		cl.w.Write([]byte(line))
		cl.w.Write([]byte("\n"))
	}
	cl.mu.Unlock()
}

// PrintProgress attempts to initialise a progress bar. If already started, it updates the progress value.
func (cl ConsoleLogger) PrintProgress(progress int64) {
	cl.mu.Lock()
	if !cl.progressIsSet {
		cl.progressBar = pb.New(100)
		cl.progressBar.Start()
		cl.progressIsSet = true
	}
	cl.progressBar.SetCurrent(progress)
	if progress == 100 {
		cl.progressBar.Finish()
	}
	cl.mu.Unlock()
}

// PrintBytes prints bytes directly to the console.
func (cl ConsoleLogger) PrintBytes(data []byte) {
	cl.mu.Lock()

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
	cl.mu.Unlock()
}

func (cl ConsoleLogger) printPrefix() {
	// Assumes mu locked.

	if cl.prefix == "" {
		return
	}
	c, found := cl.saltColors[cl.salt]
	if !found {
		c = availablePrefixColors[*cl.nextColorIndex]
		cl.saltColors[cl.salt] = c
		*cl.nextColorIndex = (*cl.nextColorIndex + 1) % len(availablePrefixColors)
	}
	c = cl.color(c)
	c.Fprintf(cl.w, "%s", cl.prefix)
	if cl.isFailed {
		cl.w.Write([]byte(" *"))
		cl.color(warnColor).Fprintf(cl.w, "failed")
		cl.w.Write([]byte("*"))
	}
	cl.w.Write([]byte(" | "))
	if cl.isCached {
		cl.w.Write([]byte("*"))
		cl.color(cachedColor).Fprintf(cl.w, "cached")
		cl.w.Write([]byte("* "))
	}
}

func (cl ConsoleLogger) color(c *color.Color) *color.Color {
	switch cl.colorMode {
	case NoColor:
		return noColor
	case ForceColor:
		return c
	case AutoColor:
		if color.NoColor {
			return noColor
		}
		return c
	}
	return noColor
}
