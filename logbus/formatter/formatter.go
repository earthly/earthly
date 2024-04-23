package formatter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	runc "github.com/containerd/go-runc"
	humanize "github.com/dustin/go-humanize"
	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/earthly/earthly/util/execstatssummary"
	"github.com/earthly/earthly/util/progressbar"
	"github.com/hashicorp/go-multierror"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const (
	durationBetweenSha256ProgressUpdate = 5 * time.Second
	durationBetweenProgressUpdate       = 3 * time.Second
	durationBetweenProgressUpdateIfSame = 5 * time.Millisecond
	durationBetweenOngoingUpdates       = 5 * time.Second
	durationBetweenOngoingUpdatesNoAnsi = 60 * time.Second

	// BuildkitStatsStream is the stream number associated with runc stats
	BuildkitStatsStream = 99 // TODO move to a common location in buildkit
)

const esc = 27

const (
	genericPrefix = "_generic:"
)

var (
	ansiUp            = []byte(fmt.Sprintf("%c[A", esc))
	ansiEraseRestLine = []byte(fmt.Sprintf("%c[K", esc))
	ansiSupported     = os.Getenv("TERM") != "dumb" &&
		(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

// TODO(vladaionescu): What to do with interactive mode? We need a way for an external
//                     process to signal interactive.

type command struct {
	lastProgress   time.Time
	lastPercentage int32
	// openLine is the line of output that has not yet been terminated with a \n.
	openLine []byte
}

// Formatter is a delta to console logger.
type Formatter struct {
	bus              *logbus.Bus
	console          conslogging.ConsoleLogger
	verbose          bool
	displayStats     bool
	execStatsTracker *execstatssummary.Tracker
	ongoingTick      time.Duration
	ongoingTicker    *time.Ticker
	startTime        time.Time
	closedCh         chan struct{}
	defaultPlatform  string

	mu                         sync.Mutex
	interactives               map[string]struct{} // set of command IDs
	lastOutputWasProgress      bool
	lastOutputWasOngoingUpdate bool
	lastCommandOutput          *command
	timingTable                map[string]time.Duration // targetID -> duration
	manifest                   *logstream.RunManifest
	commands                   map[string]*command
	errors                     []error
}

// New creates a new Formatter.
func New(ctx context.Context, b *logbus.Bus, debug, verbose, displayStats, forceColor, noColor, disableOngoingUpdates bool, execStatsTracker *execstatssummary.Tracker) *Formatter {
	ongoingTick := durationBetweenOngoingUpdatesNoAnsi
	if ansiSupported {
		ongoingTick = durationBetweenOngoingUpdates
	}
	ongoingTicker := time.NewTicker(ongoingTick)
	ongoingTicker.Stop()
	var logLevel conslogging.LogLevel
	switch {
	case debug:
		logLevel = conslogging.Debug
	case verbose:
		logLevel = conslogging.Verbose
	default:
		logLevel = conslogging.Info
	}
	var colorMode conslogging.ColorMode
	switch {
	case forceColor:
		colorMode = conslogging.ForceColor
	case noColor:
		colorMode = conslogging.NoColor
	default:
		colorMode = conslogging.AutoColor
	}
	f := &Formatter{
		bus:              b,
		console:          conslogging.New(nil, nil, colorMode, conslogging.DefaultPadding, logLevel),
		verbose:          verbose,
		displayStats:     displayStats,
		execStatsTracker: execStatsTracker,
		timingTable:      make(map[string]time.Duration),
		startTime:        time.Now(),
		closedCh:         make(chan struct{}),
		ongoingTicker:    ongoingTicker,
		ongoingTick:      ongoingTick,
		manifest:         &logstream.RunManifest{},
		commands:         make(map[string]*command),
		interactives:     make(map[string]struct{}),
	}
	if !disableOngoingUpdates {
		go f.ongoingTickLoop(ctx)
	}
	return f
}

// Write writes a delta to the console.
func (f *Formatter) Write(delta *logstream.Delta) {
	f.mu.Lock()
	defer f.mu.Unlock()
	err := f.processDelta(delta)
	if err != nil {
		f.errors = append(f.errors, err)
	}
}

// SetDefaultPlatform sets the default platform.
func (f *Formatter) SetDefaultPlatform(platform string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.defaultPlatform = platform
}

// Close stops the formatter and returns any errors encountered during
// formatting.
func (f *Formatter) Close() error {
	close(f.closedCh)
	f.mu.Lock()
	defer f.mu.Unlock()
	var retErr error
	for _, err := range f.errors {
		retErr = multierror.Append(retErr, err)
	}
	return retErr
}

// Manifest returns a copy of the manifest.
func (f *Formatter) Manifest() *logstream.RunManifest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return proto.Clone(f.manifest).(*logstream.RunManifest)
}

func (f *Formatter) processDelta(delta *logstream.Delta) error {
	err := deltautil.ApplyDelta(f.manifest, delta)
	if err != nil {
		return errors.Wrap(err, "failed to apply delta")
	}
	switch d := delta.GetDeltaTypeOneof().(type) {
	case *logstream.Delta_DeltaManifest:
		err := f.handleDeltaManifest(d.DeltaManifest)
		if err != nil {
			return errors.Wrap(err, "failed to handle delta manifest")
		}
	case *logstream.Delta_DeltaLog:
		err := f.handleDeltaLog(d.DeltaLog)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown delta type %T", d)
	}
	return nil
}

func (f *Formatter) ongoingTickLoop(ctx context.Context) {
	f.ongoingTicker.Reset(f.ongoingTick)
	defer f.ongoingTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-f.closedCh:
			return
		case <-f.ongoingTicker.C:
			f.mu.Lock()
			err := f.processOngoingTick(ctx)
			if err != nil {
				f.errors = append(f.errors, err)
			}
			f.mu.Unlock()
		}
	}
}

func (f *Formatter) handleDeltaManifest(dm *logstream.DeltaManifest) error {
	for commandID, cmd := range dm.GetFields().GetCommands() {
		cm, ok := f.manifest.GetCommands()[commandID]
		if !ok {
			return fmt.Errorf("command %q not found in manifest", commandID)
		}
		var tm *logstream.TargetManifest
		if cm.GetTargetId() != "" {
			var ok bool
			tm, ok = f.manifest.GetTargets()[cm.GetTargetId()]
			if !ok {
				return fmt.Errorf("target %s not found in manifest", cm.GetTargetId())
			}
		}
		if cmd.GetHasInteractive() && cmd.GetIsInteractive() {
			if cm.GetEndedAtUnixNanos() == 0 {
				if len(f.interactives) == 0 {
					f.ongoingTicker.Stop()
				}
				f.interactives[commandID] = struct{}{}
			} else if cmd.GetEndedAtUnixNanos() != 0 {
				delete(f.interactives, commandID)
				if len(f.interactives) == 0 {
					f.ongoingTicker.Reset(f.ongoingTick)
				}
			}
		}
		if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_IN_PROGRESS {
			f.printHeader(cm.GetTargetId(), commandID, tm, cm, false)
		}
		if cmd.GetHasHasProgress() && f.shouldPrintProgress(cm.GetTargetId(), commandID, cm) {
			f.printProgress(cm.GetTargetId(), commandID, cm)
		}
		if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_FAILURE && cm.GetTargetId() != "" {
			f.printError(cm.GetTargetId(), commandID, tm, cm)
		}
	}
	if dm.GetFields().GetHasFailure() {
		f.printBuildFailure()
	}
	return nil
}

func (f *Formatter) getCommand(commandID string) *command {
	cmd, ok := f.commands[commandID]
	if !ok {
		cmd = &command{}
		f.commands[commandID] = cmd
	}
	return cmd
}

func (f *Formatter) handleDeltaLog(dl *logstream.DeltaLog) error {
	commandID := dl.GetCommandId()
	targetID := dl.GetTargetId()

	//lookup --raw-output from the command manifest
	cm := f.manifest.GetCommands()[dl.GetCommandId()]
	rawOutput := false
	if cm != nil {
		//commandStr building order in converter.internalRun means
		// --raw-output is always first after run
		rawOutput = strings.Contains(cm.Name, "RUN --raw-output")
	}
	c, verboseOnly := f.targetConsole(targetID, commandID, rawOutput)
	if verboseOnly && !f.verbose {
		return nil
	}
	cmd := f.getCommand(dl.GetCommandId())

	sameAsLast := (!f.lastOutputWasOngoingUpdate &&
		!f.lastOutputWasProgress &&
		f.lastCommandOutput == cmd)

	output := dl.GetData()

	if dl.Stream == BuildkitStatsStream {
		var stats runc.Stats
		err := json.Unmarshal(output, &stats)
		if err != nil {
			return errors.Wrap(err, "failed to parse stats")
		}
		totalCPU := time.Duration(stats.Cpu.Usage.Total) // Total is reported in nanoseconds
		totalMem := stats.Memory.Usage.Usage             // in bytes
		output = []byte(fmt.Sprintf("[stats] total CPU: %s; total memory: %s\n", totalCPU, humanize.Bytes(totalMem)))
		if f.execStatsTracker != nil {
			f.execStatsTracker.Observe(f.targetName(dl.GetTargetId()), f.commandName(dl.GetCommandId()), totalMem, totalCPU)
		}
		if !f.displayStats {
			return nil
		}
	}

	printOutput := make([]byte, 0, len(cmd.openLine)+len(output)+10)
	if bytes.HasPrefix(output, []byte{'\n'}) && len(cmd.openLine) > 0 {
		// Optimization for cases where ansi control sequences are not supported:
		// if the output starts with a \n, then treat the open line as closed and
		// just keep going after that.
		cmd.openLine = nil
		output = output[1:]
	}
	if sameAsLast && len(cmd.openLine) > 0 {
		// Prettiness optimization: if there is an open line and the previous print out
		// was of the same vertex, then use ANSI control sequence to go up one line and
		// keep writing there.
		printOutput = append(printOutput, ansiUp...)
	}
	// Prepend the open line to the output.
	printOutput = append(printOutput, cmd.openLine...)
	printOutput = append(printOutput, output...)
	// Look for the last \n to update the open line.
	lastNewLine := bytes.LastIndexByte(printOutput, '\n')
	if lastNewLine != -1 {
		// Ends up being empty slice if output ends in \n.
		cmd.openLine = printOutput[(lastNewLine + 1):]
	} else {
		// No \n found - update cmd.openLine to append the new output.
		cmd.openLine = append(cmd.openLine, output...)
	}
	if !bytes.HasSuffix(printOutput, []byte{'\n'}) {
		// If output doesn't terminate in \n, add our own.
		printOutput = append(printOutput, '\n')
	}

	c.PrintBytes(printOutput)
	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = cmd
	return nil
}

func (f *Formatter) processOngoingTick(ctx context.Context) error {
	c := f.console.WithWriter(f.bus.FormattedWriter("ongoing", "")).WithPrefix("ongoing")
	c.VerbosePrintf("ongoing TODO\n")
	// TODO(vladaionescu): Go through all the commands and find which one is ongoing.
	// Print their targets on the console.
	f.lastOutputWasOngoingUpdate = true
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
	return nil
}

func (f *Formatter) printHeader(targetID string, commandID string, tm *logstream.TargetManifest, cm *logstream.CommandManifest, failure bool) {
	c, verboseOnly := f.targetConsole(targetID, commandID, false)
	if verboseOnly && !f.verbose {
		return
	}
	if failure {
		c = c.WithFailed(true)
	}
	var metaParts []string
	if cm.GetPlatform() != "" && cm.GetPlatform() != f.defaultPlatform {
		metaParts = append(metaParts, cm.GetPlatform())
	}
	if tm != nil && tm.GetOverrideArgs() != nil {
		metaParts = append(metaParts, strings.Join(tm.GetOverrideArgs(), " "))
	}
	if len(metaParts) > 0 {
		c.WithMetadataMode(true).Printf("%s\n", strings.Join(metaParts, " | "))
	}
	out := []string{}
	out = append(out, "-->")
	out = append(out, cm.GetName())
	if cm.GetIsCached() {
		c = c.WithCached(true)
	}
	c.Printf("%s\n", strings.Join(out, " "))

	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
}

func (f *Formatter) printProgress(targetID string, commandID string, cm *logstream.CommandManifest) {
	c, verboseOnly := f.targetConsole(targetID, commandID, false)
	if verboseOnly && !f.verbose {
		return
	}
	builder := make([]string, 0, 2)
	if f.lastOutputWasProgress {
		builder = append(builder, string(ansiUp))
	}
	progressBar := progressbar.ProgressBar(int(cm.GetProgress()), 10)
	builder = append(builder, fmt.Sprintf(
		"[%s] %3d%% %s%s\n",
		progressBar, cm.GetProgress(), cm.GetName(), string(ansiEraseRestLine)))
	c.PrintBytes([]byte(strings.Join(builder, "")))
	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = (cm.GetProgress() != 100)
	f.lastCommandOutput = nil
}

func (f *Formatter) shouldPrintProgress(targetID string, commandID string, cm *logstream.CommandManifest) bool {
	if !cm.GetHasProgress() {
		return false
	}
	// TODO(vladaionescu): Skip some internal progress for non-ansi.
	minDelta := durationBetweenOngoingUpdates
	if f.lastOutputWasProgress && ansiSupported {
		minDelta = durationBetweenProgressUpdateIfSame
	}
	// TODO(vladaionescu): Handle sha256 progress in a special manner.
	// } else if strings.HasPrefix(id, "sha256:") || strings.HasPrefix(id, "extracting sha256:") {
	// 	minDelta = durationBetweenSha256ProgressUpdate
	// }
	cmd := f.getCommand(commandID)
	lastProgress := cmd.lastProgress
	if time.Since(lastProgress) < minDelta && cm.GetProgress() < 100 {
		return false
	}
	if cmd.lastPercentage == cm.GetProgress() {
		return false
	}
	cmd.lastPercentage = cm.GetProgress()
	cmd.lastProgress = time.Now()
	return true
}

func (f *Formatter) printError(targetID string, commandID string, tm *logstream.TargetManifest, cm *logstream.CommandManifest) {
	c, _ := f.targetConsole(targetID, commandID, false)
	c.Printf("%s\n", cm.GetErrorMessage())
	c.VerbosePrintf("Overriding args used: %s\n", strings.Join(tm.GetOverrideArgs(), " "))
	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
}

func (f *Formatter) printBuildFailure() {
	failure := f.manifest.GetFailure()
	if failure.GetErrorMessage() == "" {
		return
	}
	var tm *logstream.TargetManifest
	var cm *logstream.CommandManifest
	if failure.GetTargetId() != "" {
		tm = f.manifest.GetTargets()[failure.GetTargetId()]
	}
	if failure.GetCommandId() != "" {
		cm = f.manifest.GetCommands()[failure.GetCommandId()]
	}
	c, _ := f.targetConsole(failure.GetTargetId(), failure.GetCommandId(), false)
	c = c.WithFailed(true)
	msgPrefix := "Error: " // print this prefix only when the command id is set
	if failure.GetCommandId() != "" && failure.GetCommandId() != logbus.GenericDefault {
		msgPrefix = ""
		c.PrintFailure("")
		c.Printf("Repeating the failure error...\n")
		f.printHeader(failure.GetTargetId(), failure.GetCommandId(), tm, cm, true)
		if len(failure.GetOutput()) > 0 {
			c.PrintBytes(failure.GetOutput())
		} else {
			c.Printf("[no output]\n")
		}
	}
	c.Printf("%s%s\n", msgPrefix, failure.GetErrorMessage())

	pattern := `ERROR (/.+) line (\d+):(\d+)\s+([\s\S]+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(failure.GetErrorMessage())
	if len(matches) == 5 {
		file := matches[1]
		line := matches[2]
		col := matches[3]
		message := strings.Join(strings.Fields(matches[4]), " ")
		c.PrintGHAError(message, file, line, col)
	} else {
		c.Printf("%s%s\n", msgPrefix, failure.GetErrorMessage())
		c.PrintGHAError(failure.GetErrorMessage())
	}

	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
}

func (f *Formatter) targetName(targetID string) string {
	if tm, ok := f.manifest.GetTargets()[targetID]; ok {
		return tm.GetName()
	}
	return "unknown"
}

func (f *Formatter) commandName(commandID string) string {
	if cm, ok := f.manifest.GetCommands()[commandID]; ok {
		return cm.GetName()
	}
	return "unknown"
}

func (f *Formatter) targetConsole(targetID string, commandID string, rawOutput bool) (conslogging.ConsoleLogger, bool) {
	var targetName string
	var writerTargetID string
	verboseOnly := false
	switch {
	case targetID != "":
		tm := f.manifest.GetTargets()[targetID]
		targetName = tm.GetName()
		writerTargetID = targetID
	case commandID == logbus.GenericDefault:
		targetName = ""
		writerTargetID = commandID
	case strings.HasPrefix(commandID, genericPrefix):
		targetName = strings.TrimPrefix(commandID, genericPrefix)
		writerTargetID = commandID
		switch targetName {
		case "context":
			verboseOnly = true
		default:
		}
	case commandID != "":
		cm, ok := f.manifest.GetCommands()[commandID]
		if ok {
			targetName = cm.GetCategory()
			if targetName == "" {
				targetName = cm.GetName()
			}
		}
		switch {
		case strings.HasPrefix(targetName, "internal "):
			verboseOnly = true
			targetName = strings.TrimPrefix(targetName, "internal ")
		case targetName == "internal":
			verboseOnly = true
		case targetName == "context":
			verboseOnly = true
		case targetName == "":
			verboseOnly = true
			targetName = fmt.Sprintf("_internal:%s", commandID)
		default:
		}
		writerTargetID = commandID
	default:
		targetName = "_unknown"
		writerTargetID = "_unknown"
	}
	if rawOutput {
		return f.console.
			WithWriter(f.bus.FormattedWriter(writerTargetID, commandID)), verboseOnly
	}

	return f.console.
		WithWriter(f.bus.FormattedWriter(writerTargetID, commandID)).
		WithPrefixAndSalt(targetName, writerTargetID), verboseOnly
}
