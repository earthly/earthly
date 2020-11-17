package conslogging

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

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

const (
	// NoPadding means the old behavior of printing the full target only.
	NoPadding int = -1
	// DefaultPadding always prints 20 characters for the target, right
	// justified. If it is longer, it prints the right 20 characters.
	DefaultPadding int = 20
)

var currentConsoleMutex sync.Mutex

// ConsoleLogger is a writer for consoles.
type ConsoleLogger struct {
	prefix string
	// params are printed right after the prefix delimiter.
	params string
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
	prefixPadding  int
}

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int) ConsoleLogger {
	return ConsoleLogger{
		w:              os.Stdout,
		colorMode:      colorMode,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		prefixPadding:  prefixPadding,
		mu:             &currentConsoleMutex,
	}
}

func (cl ConsoleLogger) clone() ConsoleLogger {
	return ConsoleLogger{
		w:              cl.w,
		prefix:         cl.prefix,
		params:         cl.params,
		salt:           cl.salt,
		isCached:       cl.isCached,
		isFailed:       cl.isFailed,
		saltColors:     cl.saltColors,
		colorMode:      cl.colorMode,
		nextColorIndex: cl.nextColorIndex,
		prefixPadding:  cl.prefixPadding,
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

// WithParams returns a ConsoleLogger with params added.
func (cl ConsoleLogger) WithParams(params string) ConsoleLogger {
	ret := cl.clone()
	ret.params = params
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
	cl.color(successColor).Fprintf(cl.w, "=========================== SUCCESS ===========================\n")
}

// PrintFailure prints the failure message.
func (cl ConsoleLogger) PrintFailure() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.color(warnColor).Fprintf(cl.w, "=========================== FAILURE ===========================\n")
}

// Warnf prints a warning message in red
func (cl ConsoleLogger) Warnf(format string, args ...interface{}) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	c := cl.color(warnColor)
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
	c, found := cl.saltColors[cl.salt]
	if !found {
		c = availablePrefixColors[*cl.nextColorIndex]
		cl.saltColors[cl.salt] = c
		*cl.nextColorIndex = (*cl.nextColorIndex + 1) % len(availablePrefixColors)
	}
	c = cl.color(c)
	c.Fprintf(cl.w, cl.prettyPrefix())
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
	if cl.params != "" {
		cl.color(paramsColor).Fprintf(cl.w, cl.params)
		cl.w.Write([]byte(" "))
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

var bracketsRegexp = regexp.MustCompile("\\(([^\\]]*)\\)")

func (cl ConsoleLogger) prettyPrefix() string {
	if cl.prefixPadding == NoPadding {
		return cl.prefix
	}

	var brackets string
	bracketParts := strings.SplitN(cl.prefix, "(", 2)
	if len(bracketParts) > 1 {
		brackets = fmt.Sprintf("(%s", bracketParts[1])
	}
	prettyPrefix := bracketParts[0]
	if len(cl.prefix) > cl.prefixPadding {
		parts := strings.Split(cl.prefix, "/")
		target := parts[len(parts)-1]

		truncated := ""
		for _, part := range parts[:len(parts)-1] {
			letter := part
			if len(part) > 0 && part != ".." {
				letter = string(part[0])
			}

			truncated += letter + "/"
		}

		prettyPrefix = truncated + target
	}

	formatString := fmt.Sprintf("%%%vv", cl.prefixPadding)
	return fmt.Sprintf(formatString, fmt.Sprintf("%s%s", prettyPrefix, brackets))
}
