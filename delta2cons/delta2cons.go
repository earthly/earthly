package delta2cons

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/mattn/go-isatty"
	"github.com/moby/buildkit/client"
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

// Delta2Cons is a delta to console logger.
type Delta2Cons struct {
	console                    conslogging.ConsoleLogger
	verbose                    bool
	disableOngoingUpdates      bool
	lastOutputWasProgress      bool
	lastOutputWasOngoingUpdate bool
	timingTable                map[string]time.Duration // targetID -> duration
	startTime                  time.Time
	ongoingTicker              *time.Ticker
	ongoingTick                time.Duration
	manifest                   *logstream.RunManifest
}

// NewDelta2Cons creates a new Delta2Cons.
func NewDelta2Cons(console conslogging.ConsoleLogger, verbose bool, disableOngoingUpdates bool) *Delta2Cons {
	ongoingTick := durationBetweenOngoingUpdatesNoAnsi
	if ansiSupported {
		ongoingTick = durationBetweenOngoingUpdates
	}
	return &Delta2Cons{
		console:               console,
		verbose:               verbose,
		disableOngoingUpdates: disableOngoingUpdates,
		timingTable:           make(map[string]time.Duration),
		startTime:             time.Now(),
		ongoingTicker:         time.NewTicker(ongoingTick),
		ongoingTick:           ongoingTick,
		manifest:              &logstream.RunManifest{},
	}
}

// PipeDeltasToConsole takes a channel of deltas interprets them and
// writes them to the console.
func (d2c *Delta2Cons) PipeDeltasToConsole(ctx context.Context, ch chan *logstream.Delta, bkClient *client.Client) error {
	for {
		select {
		case <-ctx.Done():
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
			err := d2c.processOngoingTick(ctx, bkClient)
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
			if cmd.GetStatus() == logstream.RunStatus_RUN_STATUS_IN_PROGRESS {
				tm, ok := d2c.manifest.GetTargets()[targetID]
				if !ok {
					return fmt.Errorf("target %s not found in manifest", targetID)
				}
				cm := tm.GetCommands()[index]
				d2c.printHeader(targetID, index, tm, cm)
			}
			if cmd.GetHasHasProgress() && cmd.GetHasProgress() {
				d2c.printProgress(targetID, index, cmd.GetProgress())
			}
		}
	}
	return nil
}

func (d2c *Delta2Cons) handleDeltaLog(ctx context.Context, dl *logstream.DeltaLog) error {
	// TODO
	// Lookup target for the log in the manifest to create color,
	// and write it to the console.

	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = false
	return fmt.Errorf("not implemented")
}

func (d2c *Delta2Cons) processOngoingTick(ctx context.Context, bkClient *client.Client) error {
	if d2c.disableOngoingUpdates {
		return nil
	}
	d2c.console.WithPrefix("ongoing").Printf("ongoing TODO\n")
	// Go through all the commands and find which one is ongoing.
	// Print their targets on the console.
	// TODO
	d2c.lastOutputWasOngoingUpdate = true
	d2c.lastOutputWasProgress = false
	return fmt.Errorf("not implemented")
}

func (d2c *Delta2Cons) printHeader(targetID string, index int32, tm *logstream.TargetManifest, cm *logstream.CommandManifest) {
	c := d2c.console.WithPrefixAndSalt(tm.GetName(), targetID)
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
}

func (d2c *Delta2Cons) printProgress(targetID string, index int32, progress int32) {
	// TODO
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = (progress != 100)
}

func (d2c *Delta2Cons) shouldPrintProgress(targetID string, index int32, progress int32, verbose bool, sameAsLast bool) bool {
	// minDelta := durationBetweenOngoingUpdates
	// if sameAsLast && ansiSupported {
	// 	minDelta = durationBetweenProgressUpdateIfSame
	// } else if strings.HasPrefix(id, "sha256:") || strings.HasPrefix(id, "extracting sha256:") {
	// 	minDelta = durationBetweenSha256ProgressUpdate
	// }

	// TODO
	return true
}

// TODO: What to do with interactive mode? We need a way for an external
//       process to signal interactive.
