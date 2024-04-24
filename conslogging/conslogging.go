package conslogging

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	salt            string
	colorMode       ColorMode
	isCached        bool
	isFailed        bool
	isGitHubActions bool
	logLevel        LogLevel

	// The following are shared between instances and are protected by the mutex.
	mu             *sync.Mutex
	saltColors     map[string]*color.Color
	nextColorIndex *int
	errW           io.Writer
	consoleErrW    io.Writer
	prefixWriter   PrefixWriter
	trailingLine   bool
	prefixPadding  int
	bb             *BundleBuilder
}

// Current returns the current console.
func Current(colorMode ColorMode, prefixPadding int, logLevel LogLevel, isGitHubActions bool) ConsoleLogger {
	return New(getCompatibleStderr(), &currentConsoleMutex, colorMode, prefixPadding, logLevel, isGitHubActions)
}

// New returns a new ConsoleLogger with a predefined target writer.
func New(w io.Writer, mu *sync.Mutex, colorMode ColorMode, prefixPadding int, logLevel LogLevel, isGitHubActions bool) ConsoleLogger {
	if mu == nil {
		mu = &sync.Mutex{}
	}
	return ConsoleLogger{
		consoleErrW:     w,
		errW:            w,
		colorMode:       colorMode,
		saltColors:      make(map[string]*color.Color),
		nextColorIndex:  new(int),
		prefixPadding:   prefixPadding,
		mu:              mu,
		logLevel:        logLevel,
		isGitHubActions: isGitHubActions,
	}
}

func (cl ConsoleLogger) clone() ConsoleLogger {
	return ConsoleLogger{
		consoleErrW:     cl.consoleErrW,
		errW:            cl.errW,
		prefixWriter:    cl.prefixWriter,
		prefix:          cl.prefix,
		metadataMode:    cl.metadataMode,
		isLocal:         cl.isLocal,
		logLevel:        cl.logLevel,
		salt:            cl.salt,
		isCached:        cl.isCached,
		isFailed:        cl.isFailed,
		isGitHubActions: cl.isGitHubActions,
		saltColors:      cl.saltColors,
		colorMode:       cl.colorMode,
		nextColorIndex:  cl.nextColorIndex,
		prefixPadding:   cl.prefixPadding,
		mu:              cl.mu,
		bb:              cl.bb,
	}
}

// WithPrefix returns a ConsoleLogger with a prefix added.
func (cl ConsoleLogger) WithPrefix(prefix string) ConsoleLogger {
	return cl.WithPrefixAndSalt(prefix, prefix)
}

// WithPrefixAndSalt returns a ConsoleLogger with a prefix and a seed added.
func (cl ConsoleLogger) WithPrefixAndSalt(prefix string, salt string) ConsoleLogger {
	ret := cl.clone()
	if cl.bb != nil {
		ret.errW = io.MultiWriter(cl.consoleErrW, cl.bb.PrefixWriter(prefix))
	}
	ret.prefix = prefix
	ret.salt = salt
	if ret.prefixWriter != nil {
		ret.prefixWriter = ret.prefixWriter.WithPrefix(prefix)
		ret.errW = ret.prefixWriter
	}
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

// WithPrefixWriter returns a ConsoleLogger with a prefix writer.
func (cl ConsoleLogger) WithPrefixWriter(w PrefixWriter) ConsoleLogger {
	ret := cl.clone()
	ret.prefixWriter = w
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
	w := new(bytes.Buffer)
	cl.mu.Lock()
	defer func() {
		_, _ = w.WriteTo(cl.errW)
		cl.mu.Unlock()
	}()
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
	cl.printGithubActionsControl(groupCommand, msg)
	c.Fprintf(w, " %s", msg)
	fmt.Fprintf(w, "\n")
	c.Fprintf(w, "%s", strings.Repeat("—", underlineLength))
	fmt.Fprintf(w, "\n\n")
}

// PrintPhaseFooter prints the phase footer.
func (cl ConsoleLogger) PrintPhaseFooter(phase string, disabled bool, special string) {
	w := new(bytes.Buffer)
	cl.mu.Lock()
	defer func() {
		_, _ = w.WriteTo(cl.errW)
		cl.mu.Unlock()
	}()
	c := cl.color(noColor)
	cl.printGithubActionsControl(endGroupCommand, phase)
	c.Fprintf(w, "\n")
}

// PrintSuccess prints the success message.
func (cl ConsoleLogger) PrintSuccess() {
	cl.PrintBar(successColor, "🌍 Earthly Build  ✅ SUCCESS", "")
}

// PrintFailure prints the failure message.
func (cl ConsoleLogger) PrintFailure(phase string) {
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

// Prints a GitHub Actions summary message to GITHUB_STEP_SUMMARY
func (cl *ConsoleLogger) PrintGHASummary(message string) {
	if !cl.isGitHubActions {
		return
	}

	path := os.Getenv("GITHUB_STEP_SUMMARY")
	if path == "" {
		return
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.WriteString(message + "\n")
}

// PrintGHAError constructs a GitHub Actions error message.
// The `file`, `line`, and `col` parameters are optional.
func (cl *ConsoleLogger) PrintGHAError(message string, details ...string) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	file := ""
	line := ""
	col := ""

	if len(details) >= 3 {
		file, line, col = details[0], details[1], details[2]
	}

	if file != "" && line != "" && col != "" {
		cl.printGithubActionsControl(errorCommand, "file=%s,line=%s,col=%s,title=Error::%s", file, line, col, message)
	} else {
		cl.printGithubActionsControl(errorCommand, "title=Error::%s", message)
	}
}

type ghHeader string

const (
	errorCommand    ghHeader = "::error"
	groupCommand    ghHeader = "::group::"
	endGroupCommand ghHeader = "::endgroup::"
)

// Print GHA control messages like ::group and ::error
func (cl ConsoleLogger) printGithubActionsControl(header ghHeader, format string, a ...any) {
	if !cl.isGitHubActions {
		return
	}
	// Assumes mu locked.
	w := new(bytes.Buffer)
	defer func() {
		_, _ = w.WriteTo(cl.errW)
	}()

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fullFormat := string(header) + " " + format

	fmt.Fprintf(w, fullFormat, a...)
}

// PrintBar prints an earthly message bar.
func (cl ConsoleLogger) PrintBar(c *color.Color, msg, phase string) {
	w := new(bytes.Buffer)
	cl.mu.Lock()
	defer func() {
		_, _ = w.WriteTo(cl.errW)
		cl.mu.Unlock()
	}()
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

	fmt.Fprintf(w, "\n")
	c.Fprintf(w, "%s%s%s", leftBar, center, rightBar)
	fmt.Fprintf(w, "\n\n")
}

// Warnf prints a warning message in red to errWriter.
func (cl ConsoleLogger) Warnf(format string, args ...interface{}) {
	c := cl.color(warnColor)
	cl.colorPrintf(Warn, c, format, args...)
}

// VerboseWarnf prints a warning message in red to errWriter when verbose flag is set.
func (cl ConsoleLogger) VerboseWarnf(format string, args ...interface{}) {
	if cl.logLevel < Verbose {
		return
	}
	cl.Warnf(format, args...)
}

// HelpPrintf prints formatted text to the console with `Help:` prefix in a specific color
func (cl ConsoleLogger) HelpPrintf(format string, args ...interface{}) {
	cl.ColorPrintf(cl.color(helpColor), fmt.Sprintf("\nHelp: %s\n", format), args...)
}

// Printf prints formatted text to the console.
func (cl ConsoleLogger) Printf(format string, args ...interface{}) {
	c := cl.color(noColor)
	if cl.metadataMode {
		c = cl.color(metadataModeColor)
	}
	cl.ColorPrintf(c, format, args...)
}

func (cl ConsoleLogger) colorPrintf(level LogLevel, c *color.Color, format string, args ...interface{}) {
	if cl.logLevel < level {
		return
	}
	w := new(bytes.Buffer)
	cl.mu.Lock()
	defer func() {
		_, _ = w.WriteTo(cl.errW)
		cl.mu.Unlock()
	}()

	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")
	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix(w)
		c.Fprintf(w, "%s", line)

		// Don't use a background color for \n.
		noColor.Fprintf(w, "\n")
	}
}

func (cl ConsoleLogger) ColorPrintf(c *color.Color, format string, args ...interface{}) {
	cl.colorPrintf(Info, c, format, args...)
}

// PrintBytes prints bytes directly to the console.
func (cl ConsoleLogger) PrintBytes(data []byte) {
	w := new(bytes.Buffer)
	w.Grow(len(data) + len(data)/4)
	cl.mu.Lock()
	defer func() {
		_, _ = w.WriteTo(cl.errW)
		cl.mu.Unlock()
	}()
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
					c.Fprintf(w, "%s", string(output))
					output = output[:0]
				}
				cl.printPrefix(w)
				cl.trailingLine = true
			}
			output = append(output, ch...)
		}
	}
	if len(output) > 0 {
		c.Fprintf(w, "%s", string(output))
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

func (cl ConsoleLogger) printPrefix(w io.Writer) {
	// Assumes mu locked.
	if cl.prefixWriter != nil {
		// When the prefix writer is in use, we don't need to print the prefix.
		return
	}
	if cl.prefix == "" {
		return
	}
	c := cl.PrefixColor()
	c.Fprintf(w, "%s", prettyPrefix(cl.prefixPadding, cl.prefix))
	if cl.isLocal {
		fmt.Fprintf(w, " *")
		cl.color(localColor).Fprintf(w, "local")
		fmt.Fprintf(w, "*")
	}
	if cl.isFailed {
		fmt.Fprintf(w, " *")
		cl.color(warnColor).Fprintf(w, "failed")
		fmt.Fprintf(w, "*")
	}
	fmt.Fprintf(w, " | ")
	if cl.isCached {
		fmt.Fprintf(w, "*")
		cl.color(cachedColor).Fprintf(w, "cached")
		fmt.Fprintf(w, "* ")
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
	return formatter.Format(prefix, prefixPadding)
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
