package conslogging

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/earthly/earthly/cleanup"
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

// LogLevel defines which types of log messages are displayed (e.g. warning, info, verbose)
type LogLevel int

const (
	// Silent silences logging
	Silent LogLevel = iota
	// Warn only display warning log messages
	Warn
	// Info displays info and higher priority log messages
	Info
	// Verbose displays verbose and higher priority log messages
	Verbose
	// Debug displays all log messages
	Debug
)

const barWidth = 80

var currentConsoleMutex sync.Mutex

// ConsoleLogger is a writer for consoles.
type ConsoleLogger struct {
	prefix string
	// metadataMode are printed in a different color.
	metadataMode bool
	// isLocal has a special prefix *local* added.
	isLocal bool
	// salt is a salt used for color consistency
	// (the same salt will get the same color).
	salt      string
	colorMode ColorMode
	isCached  bool
	isFailed  bool
	logLevel  LogLevel

	// The following are shared between instances and are protected by the mutex.
	mu             *sync.Mutex
	saltColors     map[string]*color.Color
	nextColorIndex *int
	errW           io.Writer
	consoleErrW    io.Writer
	trailingLine   bool
	prefixPadding  int
	bb             *BundleBuilder
}

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int, logLevel LogLevel) ConsoleLogger {
	return ConsoleLogger{
		consoleErrW:    getCompatibleStderr(),
		errW:           getCompatibleStderr(),
		colorMode:      colorMode,
		saltColors:     make(map[string]*color.Color),
		nextColorIndex: new(int),
		prefixPadding:  prefixPadding,
		mu:             &currentConsoleMutex,
		logLevel:       logLevel,
	}
}

func (cl ConsoleLogger) clone() ConsoleLogger {
	return ConsoleLogger{
		consoleErrW:    cl.consoleErrW,
		errW:           cl.errW,
		prefix:         cl.prefix,
		metadataMode:   cl.metadataMode,
		isLocal:        cl.isLocal,
		logLevel:       cl.logLevel,
		salt:           cl.salt,
		isCached:       cl.isCached,
		isFailed:       cl.isFailed,
		saltColors:     cl.saltColors,
		colorMode:      cl.colorMode,
		nextColorIndex: cl.nextColorIndex,
		prefixPadding:  cl.prefixPadding,
		mu:             cl.mu,
		bb:             cl.bb,
	}
}

// WithPrefix returns a ConsoleLogger with a prefix added.
func (cl ConsoleLogger) WithPrefix(prefix string) ConsoleLogger {
	ret := cl.clone()
	if cl.bb != nil {
		ret.errW = io.MultiWriter(cl.consoleErrW, cl.bb.PrefixWriter(prefix))
	}
	ret.prefix = prefix
	ret.salt = prefix
	return ret
}

// WithMetadataMode returns a ConsoleLogger with metadata printing mode set.
func (cl ConsoleLogger) WithMetadataMode(metadataMode bool) ConsoleLogger {
	ret := cl.clone()
	ret.metadataMode = metadataMode
	return ret
}

// WithLocal returns a ConsoleLogger with local set.
func (cl ConsoleLogger) WithLocal(isLocal bool) ConsoleLogger {
	ret := cl.clone()
	ret.isLocal = isLocal
	return ret
}

// WithPrefixAndSalt returns a ConsoleLogger with a prefix and a seed added.
func (cl ConsoleLogger) WithPrefixAndSalt(prefix string, salt string) ConsoleLogger {
	ret := cl.clone()
	if cl.bb != nil {
		ret.errW = io.MultiWriter(cl.consoleErrW, cl.bb.PrefixWriter(prefix))
	}
	ret.prefix = prefix
	ret.salt = salt
	return ret
}

// Prefix returns the console's prefix.
func (cl ConsoleLogger) Prefix() string {
	return cl.prefix
}

// Salt returns the console's salt.
func (cl ConsoleLogger) Salt() string {
	return cl.salt
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

// WithWriter returns a ConsoleLogger with stderr pointed at the provided io.Writer.
func (cl ConsoleLogger) WithWriter(w io.Writer) ConsoleLogger {
	ret := cl.clone()
	ret.errW = w
	return ret
}

// WithLogBundleWriter returns a ConsoleLogger with a BundleWriter attached to capture output into a log bundle, for upload to log sharing.
func (cl ConsoleLogger) WithLogBundleWriter(entrypoint string, collection *cleanup.Collection) ConsoleLogger {
	ret := cl.clone()
	ret.bb = NewBundleBuilder(entrypoint, collection)
	fullW := ret.bb.PrefixWriter(fullLog)
	ret.consoleErrW = io.MultiWriter(ret.consoleErrW, fullW)
	ret.errW = ret.consoleErrW
	return ret
}

// PrintPhaseHeader prints the phase header.
func (cl ConsoleLogger) PrintPhaseHeader(phase string, disabled bool, special string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	msg := phase
	c := cl.color(phaseColor)
	if disabled {
		c = cl.color(disabledPhaseColor)
		msg = fmt.Sprintf("%s (disabled)", msg)
	} else if special != "" {
		c = cl.color(specialPhaseColor)
		msg = fmt.Sprintf("%s (%s)", msg, special)
	}
	underlineLength := utf8.RuneCountInString(msg) + 2
	if underlineLength < barWidth {
		underlineLength = barWidth
	}
	c.Fprintf(cl.errW, " %s", msg)
	cl.errW.Write([]byte("\n"))
	c.Fprintf(cl.errW, "%s", strings.Repeat("—", underlineLength))
	cl.errW.Write([]byte("\n\n"))
}

// PrintPhaseFooter prints the phase footer.
func (cl ConsoleLogger) PrintPhaseFooter(phase string, disabled bool, special string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	c := cl.color(noColor)
	c.Fprintf(cl.errW, "\n")
}

// PrintSuccess prints the success message.
func (cl ConsoleLogger) PrintSuccess() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cl.PrintBar(successColor, "🌍 Earthly Build  ✅ SUCCESS", "")
}

// PrintFailure prints the failure message.
func (cl ConsoleLogger) PrintFailure(phase string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	cl.PrintBar(warnColor, "❌ FAILURE", phase)
}

// PrefixColor returns the color used for the prefix.
func (cl ConsoleLogger) PrefixColor() *color.Color {
	c, found := cl.saltColors[cl.salt]
	if !found {
		c = availablePrefixColors[*cl.nextColorIndex]
		cl.saltColors[cl.salt] = c
		*cl.nextColorIndex = (*cl.nextColorIndex + 1) % len(availablePrefixColors)
	}
	return cl.color(c)
}

// PrintBar prints an earthly message bar.
func (cl ConsoleLogger) PrintBar(c *color.Color, msg, phase string) {
	c = cl.color(c)
	center := msg
	if phase != "" {
		center = fmt.Sprintf("%s [%s]", msg, phase)
	}
	center = fmt.Sprintf(" %s ", center)

	sideWidth := (barWidth - utf8.RuneCountInString(center)) / 2
	if sideWidth < 0 {
		sideWidth = 0
	}
	eqBar := strings.Repeat("=", sideWidth)
	leftBar := eqBar
	rightBar := eqBar
	if utf8.RuneCountInString(center)%2 == 1 && sideWidth > 0 {
		// Ensure the width is always barWidth
		rightBar += "="
	}

	cl.errW.Write([]byte("\n"))
	c.Fprintf(cl.errW, "%s%s%s", leftBar, center, rightBar)
	cl.errW.Write([]byte("\n\n"))
}

// Warnf prints a warning message in red to errWriter.
func (cl ConsoleLogger) Warnf(format string, args ...interface{}) {
	if cl.logLevel < Warn {
		return
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()

	c := cl.color(warnColor)
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")

	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		c.Fprintf(cl.errW, "%s\n", line)
	}
}

// Printf prints formatted text to the console.
func (cl ConsoleLogger) Printf(format string, args ...interface{}) {
	if cl.logLevel < Info {
		return
	}
	cl.mu.Lock()
	defer cl.mu.Unlock()
	c := cl.color(noColor)
	if cl.metadataMode {
		c = cl.color(metadataModeColor)
	}
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")
	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		c.Fprintf(cl.errW, "%s", line)

		// Don't use a background color for \n.
		noColor.Fprintf(cl.errW, "\n")
	}
}

// PrintBytes prints bytes directly to the console.
func (cl ConsoleLogger) PrintBytes(data []byte) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	c := cl.color(noColor)
	if cl.metadataMode {
		c = cl.color(metadataModeColor)
	}

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
					c.Fprintf(cl.errW, "%s", string(output))
					output = output[:0]
				}
				cl.printPrefix()
				cl.trailingLine = true
			}
			output = append(output, ch...)
		}
	}
	if len(output) > 0 {
		c.Fprintf(cl.errW, "%s", string(output))
		// output = output[:0] // needed if output is used further in the future
	}
}

// VerbosePrintf prints formatted text to the console when verbose flag is set.
func (cl ConsoleLogger) VerbosePrintf(format string, args ...interface{}) {
	if cl.logLevel < Verbose {
		return
	}
	cl.WithMetadataMode(true).Printf(format, args...)
}

// VerboseBytes prints bytes directly to the console when verbose flag is set.
func (cl ConsoleLogger) VerboseBytes(data []byte) {
	if cl.logLevel < Verbose {
		return
	}
	cl.WithMetadataMode(true).PrintBytes(data)
}

// DebugPrintf prints formatted text to the console when debug flag is set.
func (cl ConsoleLogger) DebugPrintf(format string, args ...interface{}) {
	if cl.logLevel < Debug {
		return
	}
	cl.WithMetadataMode(true).Printf(format, args...)
}

// DebugBytes prints bytes directly to the console when debug flag is set.
func (cl ConsoleLogger) DebugBytes(data []byte) {
	if cl.logLevel < Debug {
		return
	}
	cl.WithMetadataMode(true).PrintBytes(data)
}

func (cl ConsoleLogger) printPrefix() {
	// Assumes mu locked.
	if cl.prefix == "" {
		return
	}
	c := cl.PrefixColor()
	c.Fprintf(cl.errW, prettyPrefix(cl.prefixPadding, cl.prefix))
	if cl.isLocal {
		cl.errW.Write([]byte(" *"))
		cl.color(localColor).Fprintf(cl.errW, "local")
		cl.errW.Write([]byte("*"))
	}
	if cl.isFailed {
		cl.errW.Write([]byte(" *"))
		cl.color(warnColor).Fprintf(cl.errW, "failed")
		cl.errW.Write([]byte("*"))
	}
	cl.errW.Write([]byte(" | "))
	if cl.isCached {
		cl.errW.Write([]byte("*"))
		cl.color(cachedColor).Fprintf(cl.errW, "cached")
		cl.errW.Write([]byte("* "))
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

func prettyPrefix(prefixPadding int, prefix string) string {
	if prefixPadding == NoPadding {
		return prefix
	}

	var brackets string
	bracketParts := strings.SplitN(prefix, "(", 2)
	if len(bracketParts) > 1 {
		brackets = fmt.Sprintf("(%s", bracketParts[1])
	}
	prettyPrefix := bracketParts[0]
	if len(prefix) > prefixPadding {
		parts := strings.Split(prefix, "/")
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

	formatString := fmt.Sprintf("%%%vv", prefixPadding)
	return fmt.Sprintf(formatString, fmt.Sprintf("%s%s", prettyPrefix, brackets))
}

// WithLogLevel changes the log level
func (cl ConsoleLogger) WithLogLevel(logLevel LogLevel) ConsoleLogger {
	ret := cl.clone()
	ret.logLevel = logLevel
	return ret
}

// WriteBundleToDisk makes an attached bundle writer (if any) write the collected bundle to disk.
func (cl ConsoleLogger) WriteBundleToDisk() (string, error) {
	if cl.bb == nil {
		return "", nil
	}

	return cl.bb.WriteToDisk()
}

// MarkBundleBuilderResult marks the current targets result in a log bundle for a given prefix with the current result.
func (cl ConsoleLogger) MarkBundleBuilderResult(isError, isCanceled bool) {
	if cl.bb == nil {
		return
	}

	var result string
	if isCanceled {
		result = ResultCancelled
	} else {
		if isError {
			result = ResultFailure
		} else {
			result = ResultSuccess
		}
	}

	cl.bb.PrefixResult(cl.Prefix(), result)
}

// MarkBundleBuilderStatus marks the current targets status in a log bundle for a given prefix with the current status.
func (cl ConsoleLogger) MarkBundleBuilderStatus(isStarted, isFinished, isCanceled bool) {
	if cl.bb == nil {
		return
	}

	var status string
	if isCanceled {
		status = StatusCancelled
	} else {
		if isStarted {
			if isFinished {
				status = StatusComplete
			} else {
				status = StatusInProgress
			}
		} else {
			status = StatusWaiting
		}
	}

	cl.bb.PrefixStatus(cl.Prefix(), status)
}
