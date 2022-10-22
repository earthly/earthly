package delta2cons

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/delta2cons/deltautil"
	"github.com/mattn/go-isatty"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
)

const (
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
	manifest                   *logstream.BuildManifest
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
		manifest:              &logstream.BuildManifest{},
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
			_, err := deltautil.ApplyDelta(d2c.manifest, delta)
			if err != nil {
				return errors.Wrap(err, "failed to apply delta")
			}
			for _, dm := range delta.DeltaManifests {
				err := d2c.handleDeltaManifest(ctx, dm)
				if err != nil {
					return errors.Wrap(err, "failed to handle delta manifest")
				}
			}
			for _, dl := range delta.GetDeltaLogs() {
				err := d2c.handleDeltaLog(ctx, dl)
				if err != nil {
					return err
				}
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
			if cmd.GetStatus() == logstream.BuildStatus_BUILD_STATUS_IN_PROGRESS {
				d2c.printHeader(targetID, index)
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
	// Go through all the commands and find which one is ongoing.
	// Print their targets on the console.
	// TODO
	d2c.lastOutputWasOngoingUpdate = true
	d2c.lastOutputWasProgress = false
	return fmt.Errorf("not implemented")
}

func (d2c *Delta2Cons) printHeader(targetID string, index int32) {
	// TODO
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = false
}

func (d2c *Delta2Cons) printProgress(targetID string, index int32, progress int32) {
	// TODO
	d2c.lastOutputWasOngoingUpdate = false
	d2c.lastOutputWasProgress = true
}

// TODO: What to do with interactive mode? We need a way for an external
//       process to signal interactive.
