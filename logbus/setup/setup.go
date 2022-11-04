package logbus

import (
	"context"
	"os"

	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/logbus/writersub"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// BusSetup is a helper for setting up a logbus.Bus.
type BusSetup struct {
	Bus           *logbus.Bus
	ConsoleWriter *writersub.WriterSub
	Formatter     *formatter.Formatter
	SolverMonitor *solvermon.SolverMonitor
}

// New creates a new BusSetup.
func New(ctx context.Context, bus *logbus.Bus, debug bool, verbose bool, disableOngoingUpdates bool) *BusSetup {
	bs := &BusSetup{
		Bus:           bus,
		ConsoleWriter: writersub.New(os.Stderr),
		Formatter:     nil, // set below
		SolverMonitor: nil, // set below
	}
	bs.Formatter = formatter.New(ctx, bs.Bus, verbose, disableOngoingUpdates)
	bs.Bus.AddRawSubscriber(bs.Formatter.Write)
	bs.Bus.AddFormattedSubscriber(bs.ConsoleWriter.Write)
	bs.SolverMonitor = solvermon.New(bs.Bus)
	return bs
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
	return retErr
}
