package format

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/earthly/earthly/util/progressbar"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

const (
	durationBetweenSha256ProgressUpdate = 5 * time.Second
	durationBetweenProgressUpdate       = 3 * time.Second
	durationBetweenProgressUpdateIfSame = 5 * time.Millisecond
	durationBetweenOpenLineUpdate       = time.Second
	durationBetweenOngoingUpdates       = 5 * time.Second
	durationBetweenOngoingUpdatesNoAnsi = 60 * time.Second
)

const esc = 27

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
	// Line of output that has not yet been terminated with a \n.
	openLine            []byte
	lastOpenLineUpdate  time.Time
	lastOpenLineSkipped bool
}

// Formatter is a delta to console logger.
type Formatter struct {
	bus                        *logbus.Bus
	console                    conslogging.ConsoleLogger
	verbose                    bool
	disableOngoingUpdates      bool
	lastOutputWasProgress      bool
	lastOutputWasOngoingUpdate bool
	lastCommandOutput          *command
	timingTable                map[string]time.Duration // targetID -> duration
	startTime                  time.Time
	ongoingTicker              *time.Ticker
	ongoingTick                time.Duration
	manifest                   *logstream.RunManifest
	commands                   map[string]*command
}

// New creates a new Formatter.
func New(b *logbus.Bus, verbose bool, disableOngoingUpdates bool) *Formatter {
	ongoingTick := durationBetweenOngoingUpdatesNoAnsi
	if ansiSupported {
		ongoingTick = durationBetweenOngoingUpdates
	}
	ongoingTicker := time.NewTicker(ongoingTick)
	ongoingTicker.Stop()
	return &Formatter{
		bus: b,
		// TODO (vladaionescu): Pass in color detection and log level.
		console:               conslogging.New(nil, nil, conslogging.AutoColor, conslogging.DefaultPadding, conslogging.Info),
		verbose:               verbose,
		disableOngoingUpdates: disableOngoingUpdates,
		timingTable:           make(map[string]time.Duration),
		startTime:             time.Now(),
		ongoingTicker:         ongoingTicker,
		ongoingTick:           ongoingTick,
		manifest:              &logstream.RunManifest{},
		commands:              make(map[string]*command),
	}
}

// PipeDeltasToConsole takes a channel of deltas interprets them and
// writes them to the console.
func (f *Formatter) PipeDeltasToConsole(ctx context.Context, ch chan *logstream.Delta) error {
	closeCh := make(chan struct{})
	returnedCh := make(chan struct{})
	defer close(returnedCh)
	go func() {
		<-ctx.Done()
		// Don't close immediately, as we want to print any
		// final messages that might be coming in.
		select {
		case <-returnedCh:
		case <-time.After(5 * time.Second):
		}
		close(closeCh)
	}()
	f.ongoingTicker.Reset(f.ongoingTick)
	defer f.ongoingTicker.Stop()
	for {
		select {
		case <-closeCh:
			return ctx.Err()
		case delta, ok := <-ch:
			if !ok {
				return nil
			}
			var err error
			f.manifest, err = deltautil.ApplyDeltaManifest(f.manifest, delta)
			if err != nil {
				return errors.Wrap(err, "failed to apply delta")
			}
			// TODO(vladaionescu): Make debugging a flag.
			// switch d := delta.GetDeltaTypeOneof().(type) {
			// case *logstream.Delta_DeltaManifest:
			// 	fmt.Printf("@# delta manifest: %+v\n", d)
			// case *logstream.Delta_DeltaLog:
			// 	fmt.Printf("@# delta log: %+v\n", d)
			// default:
			// }
			switch d := delta.GetDeltaTypeOneof().(type) {
			case *logstream.Delta_DeltaManifest:
				err := f.handleDeltaManifest(ctx, d.DeltaManifest)
				if err != nil {
					return errors.Wrap(err, "failed to handle delta manifest")
				}
			case *logstream.Delta_DeltaLog:
				err := f.handleDeltaLog(ctx, d.DeltaLog)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown delta type %T", d)
			}
		case <-f.ongoingTicker.C:
			err := f.processOngoingTick(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (f *Formatter) handleDeltaManifest(ctx context.Context, dm *logstream.DeltaManifest) error {
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
		if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_IN_PROGRESS {
			f.printHeader(cm.GetTargetId(), commandID, tm, cm, false)
		}
		if cmd.GetHasHasProgress() && f.shouldPrintProgress(cm.GetTargetId(), commandID, cm) {
			f.printProgress(cm.GetTargetId(), commandID, cm)
		}
		if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_FAILURE {
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

func (f *Formatter) handleDeltaLog(ctx context.Context, dl *logstream.DeltaLog) error {
	c := f.console
	if dl.GetTargetId() == "" && strings.HasPrefix("_generic:", dl.GetCommandId()) {
		prefix := strings.TrimPrefix(dl.GetCommandId(), "_generic:")
		c = c.WithPrefixAndSalt(prefix, prefix)
	} else if dl.GetTargetId() != "" {
		tm, ok := f.manifest.GetTargets()[dl.GetTargetId()]
		if !ok {
			return fmt.Errorf("target %s not found in manifest", dl.GetTargetId())
		}
		c = c.WithPrefixAndSalt(tm.GetName(), dl.GetTargetId())
	} else {
		c = c.WithPrefixAndSalt(dl.GetCommandId(), dl.GetCommandId())
	}
	cmd := f.getCommand(dl.GetCommandId())

	sameAsLast := (!f.lastOutputWasOngoingUpdate &&
		!f.lastOutputWasProgress &&
		f.lastCommandOutput == cmd)
	output := dl.GetData()
	printOutput := make([]byte, 0, len(cmd.openLine)+len(output)+10)
	if bytes.HasPrefix(output, []byte{'\n'}) && len(cmd.openLine) > 0 && !cmd.lastOpenLineSkipped {
		// Optimization for cases where ansi control sequences are not supported:
		// if the output starts with a \n, then treat the open line as closed and
		// just keep going after that.
		cmd.openLine = nil
		output = output[1:]
		cmd.lastOpenLineUpdate = time.Time{}
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
		// A \n exists - reset the open line timer.
		cmd.lastOpenLineUpdate = time.Time{}
	} else {
		// No \n found - update cmd.openLine to append the new output.
		cmd.openLine = append(cmd.openLine, output...)
	}
	if !bytes.HasSuffix(printOutput, []byte{'\n'}) {
		if time.Since(cmd.lastOpenLineUpdate) > durationBetweenOpenLineUpdate {
			// Skip printing if trying to update the same line too frequently.
			cmd.lastOpenLineSkipped = true
			return nil
		}
		cmd.lastOpenLineUpdate = time.Now()
		// If output doesn't terminate in \n, add our own.
		printOutput = append(printOutput, '\n')
	}

	cmd.lastOpenLineSkipped = false
	c.PrintBytes(printOutput)
	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = cmd
	return nil
}

func (f *Formatter) processOngoingTick(ctx context.Context) error {
	if f.disableOngoingUpdates {
		return nil
	}
	f.console.WithPrefix("ongoing").Printf("ongoing TODO\n")
	// TODO(vladaionescu): Go through all the commands and find which one is ongoing.
	// Print their targets on the console.
	f.lastOutputWasOngoingUpdate = true
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
	return nil
}

func (f *Formatter) printHeader(targetID string, commandID string, tm *logstream.TargetManifest, cm *logstream.CommandManifest, failure bool) {
	c := f.targetConsole(targetID, commandID)
	if failure {
		c = c.WithFailed(true)
	}
	var metaParts []string
	if cm.GetPlatform() != "" {
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
	c := f.targetConsole(targetID, commandID)
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
	c := f.targetConsole(targetID, commandID)
	c.Printf("%s\n", cm.GetErrorMessage())
	c.VerbosePrintf("Overriding args used: %s\n", strings.Join(tm.GetOverrideArgs(), " "))
	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
}

func (f *Formatter) printBuildFailure() {
	if f.manifest.GetFailure() == nil {
		return
	}
	failure := f.manifest.GetFailure()
	var tm *logstream.TargetManifest
	var cm *logstream.CommandManifest
	if failure.GetTargetId() != "" {
		tm = f.manifest.GetTargets()[failure.GetTargetId()]
	}
	if failure.GetCommandId() != "" {
		cm = f.manifest.GetCommands()[failure.GetCommandId()]
	}
	c := f.targetConsole(failure.GetTargetId(), failure.GetCommandId())
	c = c.WithFailed(true)
	c.Printf("Repeating the failure error...\n")
	f.printHeader(failure.GetTargetId(), failure.GetCommandId(), tm, cm, true)
	if len(failure.GetOutput()) > 0 {
		c.PrintBytes(failure.GetOutput())
	} else {
		c.Printf("[no output]\n")
	}
	if failure.GetErrorMessage() != "" {
		c.Printf("%s\n", failure.GetErrorMessage())
	}
	f.lastOutputWasOngoingUpdate = false
	f.lastOutputWasProgress = false
	f.lastCommandOutput = nil
}

func (f *Formatter) targetConsole(targetID string, commandID string) conslogging.ConsoleLogger {
	var targetName string
	writerTargetID := targetID
	if targetID != "" {
		tm := f.manifest.GetTargets()[targetID]
		targetName = tm.GetName()
	} else {
		writerTargetID = "_internal"
		targetName = "internal"
		if commandID != "" {
			writerTargetID = fmt.Sprintf("_internal:%s", commandID)
			targetName = writerTargetID
		}
	}
	return f.console.
		WithWriter(f.bus.FormattedWriter(writerTargetID)).
		WithPrefixAndSalt(targetName, targetID)
}
