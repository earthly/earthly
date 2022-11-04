package logbus

import (
	"context"
	"os"
	"strings"

	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/logbus/writersub"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// BusSetup is a helper for setting up a logbus.Bus.
type BusSetup struct {
	Bus            *logbus.Bus
	ConsoleWriter  *writersub.WriterSub
	Formatter      *formatter.Formatter
	SolverMonitor  *solvermon.SolverMonitor
	BusDebugWriter *writersub.RawWriterSub
}

// New creates a new BusSetup.
func New(ctx context.Context, bus *logbus.Bus, debug bool, verbose bool, disableOngoingUpdates bool, busDebugFile string) (*BusSetup, error) {
	bs := &BusSetup{
		Bus:           bus,
		ConsoleWriter: writersub.New(os.Stderr, "_full"),
		Formatter:     nil, // set below
		SolverMonitor: nil, // set below
	}
	bs.Formatter = formatter.New(ctx, bs.Bus, verbose, disableOngoingUpdates)
	bs.Bus.AddRawSubscriber(bs.Formatter)
	bs.Bus.AddFormattedSubscriber(bs.ConsoleWriter)
	bs.SolverMonitor = solvermon.New(bs.Bus)
	if busDebugFile != "" {
		// Open file for writing.
		f, err := os.OpenFile(busDebugFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open bus debug file %s", busDebugFile)
		}
		useJson := strings.HasSuffix(busDebugFile, ".json")
		bs.BusDebugWriter = writersub.NewRaw(f, useJson)
		bs.Bus.AddSubscriber(bs.BusDebugWriter)
	}
	return bs, nil
}

// Close closes the BusSetup.
func (bs *BusSetup) Close() error {
	var retErr error
	errs := bs.ConsoleWriter.Errors()
	var cwErr error
	for _, err := range errs {
		cwErr = multierror.Append(cwErr, err)
	}
	if cwErr != nil {
		retErr = multierror.Append(retErr, errors.Wrap(cwErr, "console writer"))
	}
	fErr := bs.Formatter.Close()
	if fErr != nil {
		retErr = multierror.Append(retErr, errors.Wrap(fErr, "formatter"))
	}
	if bs.BusDebugWriter != nil {
		errs := bs.BusDebugWriter.Errors()
		var bdwErr error
		for _, err := range errs {
			bdwErr = multierror.Append(bdwErr, err)
		}
		if bdwErr != nil {
			retErr = multierror.Append(retErr, errors.Wrap(bdwErr, "bus debug writer"))
		}
	}
	return retErr
}
