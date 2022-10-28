package delta2cons

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/conslogging"
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

type cmdKey struct {
	targetID string
	index    int32
}
type command struct {
	lastProgress   time.Time
	lastPercentage int32
	// Line of output that has not yet been terminated with a \n.
	openLine            []byte
	lastOpenLineUpdate  time.Time
	lastOpenLineSkipped bool
}

// Delta2Cons is a delta to console logger.
type Delta2Cons struct {
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
	commands                   map[cmdKey]*command
}

// New creates a new Delta2Cons.
func New(console conslogging.ConsoleLogger, verbose bool, disableOngoingUpdates bool) *Delta2Cons {
	ongoingTick := durationBetweenOngoingUpdatesNoAnsi
	if ansiSupported {
		ongoingTick = durationBetweenOngoingUpdates
	}
	ongoingTicker := time.NewTicker(ongoingTick)
	ongoingTicker.Stop()
	return &Delta2Cons{
		console:               console,
		verbose:               verbose,
		disableOngoingUpdates: disableOngoingUpdates,
		timingTable:           make(map[string]time.Duration),
		startTime:             time.Now(),
		ongoingTicker:         ongoingTicker,
		ongoingTick:           ongoingTick,
		manifest:              &logstream.RunManifest{},
		commands:              make(map[cmdKey]*command),
	}
}

// PipeDeltasToConsole takes a channel of deltas interprets them and
// writes them to the console.
func (d2c *Delta2Cons) PipeDeltasToConsole(ctx context.Context, ch chan *logstream.Delta) error {
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
	d2c.ongoingTicker.Reset(d2c.ongoingTick)
	defer d2c.ongoingTicker.Stop()
	for {
		select {
		case <-closeCh:
			return ctx.Err()
		case delta, ok := <-ch:
			if !ok {
				return nil
			}
			var err error
			d2c.manifest, err = deltautil.ApplyDeltaManifest(d2c.manifest, delta)
			if err != nil {
				return errors.Wrap(err, "failed to apply delta")
			}
			switch d := delta.GetDeltaTypeOneof().(type) {
			case *logstream.Delta_DeltaManifest:
				err := d2c.handleDeltaManifest(ctx, d.DeltaManifest)
				if err != nil {
					return errors.Wrap(err, "failed to handle delta manifest")
				}
			case *logstream.Delta_DeltaLog:
				err := d2c.handleDeltaLog(ctx, d.DeltaLog)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown delta type %T", d)
			}
		case <-d2c.ongoingTicker.C:
			err := d2c.processOngoingTick(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (d2c *Delta2Cons) handleDeltaManifest(ctx context.Context, dm *logstream.DeltaManifest) error {
	if dm.GetFields() == nil {
		return nil
	}
	if dm.GetFields().GetTargets() == nil {
		return nil
	}
	for targetID, t := range dm.GetFields().GetTargets() {
		for index, cmd := range t.GetCommands() {
			tm, ok := d2c.manifest.GetTargets()[targetID]
			if !ok {
				return fmt.Errorf("target %s not found in manifest", targetID)
			}
			cm := tm.GetCommands()[index]
			if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_IN_PROGRESS {
				d2c.printHeader(d2c.console, targetID, index, tm, cm)
			}
			if cmd.GetHasHasProgress() && d2c.shouldPrintProgress(targetID, index, cm) {
				d2c.printProgress(targetID, index, cm)
			}
			if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_FAILURE {
				d2c.printError(targetID, index, tm, cm)
			}
		}
	}
	if dm.GetFields().GetHasFailure() {
		d2c.printBuildFailure()
	}
	return nil
}

func (d2c *Delta2Cons) getCommand(targetID string, index int32) *command {
	key := cmdKey{targetID: targetID, index: index}
	cmd, ok := d2c.commands[key]
	if !ok {
		cmd = &command{}
		d2c.commands[key] = cmd
	}
	return cmd
}

func (d2c *Delta2Cons) handleDeltaLog(ctx context.Context, dl *logstream.DeltaLog) error {
	c := d2c.console
	if strings.HasPrefix("_generic:", dl.GetTargetId()) {
		prefix := strings.TrimPrefix(dl.GetTargetId(), "_generic:")
		c = c.WithPrefixAndSalt(prefix, prefix)
	} else if strings.HasPrefix("_internal:", dl.GetTargetId()) {
		prefix := strings.TrimPrefix(dl.GetTargetId(), "_internal:")
		c = c.WithPrefixAndSalt(fmt.Sprintf("internal %s", prefix), prefix)
	} else {
		tm, ok := d2c.manifest.GetTargets()[dl.GetTargetId()]
		if !ok {
			return fmt.Errorf("target %s not found in manifest", dl.GetTargetId())
		}
		c = c.WithPrefixAndSalt(tm.GetName(), dl.GetTargetId())
	}
	cmd := d2c.getCommand(dl.GetTargetId(), dl.GetCommandIndex())

	sameAsLast := (!d2c.lastOutputWasOngoingUpdate &&
		!d2c.lastOutputWasProgress &&
		d2c.lastCommandOutput == cmd)
	output := dl.GetLog()
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
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = false
	d2c.lastCommandOutput = cmd
	return nil
}

func (d2c *Delta2Cons) processOngoingTick(ctx context.Context) error {
	if d2c.disableOngoingUpdates {
		return nil
	}
	d2c.console.WithPrefix("ongoing").Printf("ongoing TODO\n")
	// TODO(vladaionescu): Go through all the commands and find which one is ongoing.
	// Print their targets on the console.
	d2c.lastOutputWasOngoingUpdate = true
	d2c.lastOutputWasProgress = false
	d2c.lastCommandOutput = nil
	return nil
}

func (d2c *Delta2Cons) printHeader(c conslogging.ConsoleLogger, targetID string, index int32, tm *logstream.TargetManifest, cm *logstream.CommandManifest) {
	c = c.WithPrefixAndSalt(tm.GetName(), targetID)
	var metaParts []string
	if tm.GetPlatform() != "" {
		metaParts = append(metaParts, tm.GetPlatform())
	}
	if tm.GetOverrideArgs() != nil {
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

	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = false
	d2c.lastCommandOutput = nil
}

func (d2c *Delta2Cons) printProgress(targetID string, index int32, cm *logstream.CommandManifest) {
	builder := make([]string, 0, 2)
	if d2c.lastOutputWasProgress {
		builder = append(builder, string(ansiUp))
	}
	progressBar := progressbar.ProgressBar(int(cm.GetProgress()), 10)
	builder = append(builder, fmt.Sprintf(
		"[%s] %3d%% %s%s\n",
		progressBar, cm.GetProgress(), cm.GetName(), string(ansiEraseRestLine)))
	d2c.console.PrintBytes([]byte(strings.Join(builder, "")))
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = (cm.GetProgress() != 100)
	d2c.lastCommandOutput = nil
}

func (d2c *Delta2Cons) shouldPrintProgress(targetID string, index int32, cm *logstream.CommandManifest) bool {
	if !cm.GetHasProgress() {
		return false
	}
	// TODO(vladaionescu): Skip some internal progress for non-ansi.
	minDelta := durationBetweenOngoingUpdates
	if d2c.lastOutputWasProgress && ansiSupported {
		minDelta = durationBetweenProgressUpdateIfSame
	}
	// TODO(vladaionescu): Handle sha256 progress in a special manner.
	// } else if strings.HasPrefix(id, "sha256:") || strings.HasPrefix(id, "extracting sha256:") {
	// 	minDelta = durationBetweenSha256ProgressUpdate
	// }
	cmd := d2c.getCommand(targetID, index)
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

func (d2c *Delta2Cons) printError(targetID string, index int32, tm *logstream.TargetManifest, cm *logstream.CommandManifest) {
	c := d2c.console.WithPrefixAndSalt(tm.GetName(), targetID)
	c.Printf("%s\n", cm.GetErrorMessage())
	c.VerbosePrintf("Overriding args used: %s\n", strings.Join(tm.GetOverrideArgs(), " "))
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = false
	d2c.lastCommandOutput = nil
}

func (d2c *Delta2Cons) printBuildFailure() {
	if d2c.manifest.GetFailure() == nil {
		return
	}
	f := d2c.manifest.GetFailure()
	var tm *logstream.TargetManifest
	var cm *logstream.CommandManifest
	if f.GetTargetId() != "" {
		tm = d2c.manifest.GetTargets()[f.GetTargetId()]
		if f.GetHasCommandIndex() {
			cm = tm.GetCommands()[f.GetCommandIndex()]
		}
	}
	c := d2c.console.WithFailed(true)
	c.Printf("Repeating the failure error...\n")
	if tm != nil {
		c = c.WithPrefixAndSalt(tm.GetName(), f.GetTargetId())
		if cm != nil {
			d2c.printHeader(c, f.GetTargetId(), f.GetCommandIndex(), tm, cm)
		}
	}
	if len(f.GetOutput()) > 0 {
		c.PrintBytes(f.GetOutput())
	} else {
		c.Printf("[no output]\n")
	}
	if f.GetSummary() != "" {
		c.Printf("%s\n", f.GetSummary())
	}
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = false
	d2c.lastCommandOutput = nil
}
