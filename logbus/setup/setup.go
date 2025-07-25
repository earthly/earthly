package setup

import (
	"context"
	"os"
	"strings"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/logbus/writersub"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/earthly/earthly/util/execstatssummary"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// BusSetup is a helper for setting up a logbus.Bus.
type BusSetup struct {
	Bus              *logbus.Bus
	ConsoleWriter    *writersub.WriterSub
	Formatter        *formatter.Formatter
	SolverMonitor    *solvermon.SolverMonitor
	BusDebugWriter   *writersub.RawWriterSub
	InitialManifest  *logstream.RunManifest
	execStatsTracker *execstatssummary.Tracker

	logStreamerStarted bool
	verbose            bool
}

// New creates a new BusSetup.
func New(ctx context.Context, bus *logbus.Bus, debug, verbose, displayStats, forceColor, noColor, disableOngoingUpdates bool, busDebugFile, buildID string, execStatsTracker *execstatssummary.Tracker, isGitHubActions bool) (*BusSetup, error) {
	bs := &BusSetup{
		Bus:           bus,
		ConsoleWriter: writersub.New(os.Stderr, "_full"),
		Formatter:     nil, // set below
		SolverMonitor: nil, // set below
		InitialManifest: &logstream.RunManifest{
			BuildId:            buildID,
			Version:            deltautil.Version,
			CreatedAtUnixNanos: uint64(bus.CreatedAt().UnixNano()),
		},
		execStatsTracker: execStatsTracker,
		verbose:          verbose,
	}
	bs.Formatter = formatter.New(ctx, bs.Bus, debug, verbose, displayStats, forceColor, noColor, disableOngoingUpdates, execStatsTracker, isGitHubActions)
	bs.Bus.AddRawSubscriber(bs.Formatter)
	bs.Bus.AddFormattedSubscriber(bs.ConsoleWriter)
	bs.SolverMonitor = solvermon.New(bs.Bus)
	if busDebugFile != "" {
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

// SetDefaultPlatform sets the default platform of the build.
func (bs *BusSetup) SetDefaultPlatform(platform string) {
	bs.Formatter.SetDefaultPlatform(platform)
}

// SetGitAuthor records the Git author information on the initial manifest.
func (bs *BusSetup) SetGitAuthor(gitAuthor, gitCommitEmail string) {
	bs.InitialManifest.GitAuthor = gitAuthor
	bs.InitialManifest.GitConfigEmail = gitCommitEmail
}

// SetCI tracks whether this build is being run in a CI environment.
func (bs *BusSetup) SetCI(isCI bool) {
	bs.InitialManifest.IsCi = isCI
}

// DumpManifestToFile dumps the manifest to the given file.
func (bs *BusSetup) DumpManifestToFile(path string) error {
	m := bs.Formatter.Manifest()
	proto.Merge(m, bs.InitialManifest)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open bus manifest debug file %s", path)
	}
	useJson := strings.HasSuffix(path, ".json")
	var dt []byte
	if useJson {
		jsonOpts := protojson.MarshalOptions{
			Multiline:       true,
			Indent:          "  ",
			UseProtoNames:   false,
			EmitUnpopulated: true,
		}
		dt, err = jsonOpts.Marshal(m)
	} else {
		dt, err = proto.Marshal(m)
	}
	if err != nil {
		return errors.Wrapf(err, "failed to marshal manifest")
	}
	_, err = f.Write(dt)
	if err != nil {
		return errors.Wrapf(err, "failed to write manifest")
	}
	return nil
}

// Close the bus setup & gather all errors.
func (bs *BusSetup) Close(ctx context.Context) error {
	var ret error

	if bs.execStatsTracker != nil {
		err := bs.execStatsTracker.Close(ctx)
		if err != nil {
			ret = multierror.Append(ret, errors.Wrap(err, "exec stats summary"))
		}
	}

	if errs := bs.ConsoleWriter.Errors(); len(errs) > 0 {
		multi := &multierror.Error{Errors: errs}
		ret = multierror.Append(ret, errors.Wrap(multi, "console writer"))
	}

	if err := bs.Formatter.Close(); err != nil {
		ret = multierror.Append(ret, errors.Wrap(err, "formatter"))
	}

	if bs.BusDebugWriter != nil {
		if errs := bs.BusDebugWriter.Errors(); len(errs) > 0 {
			multi := &multierror.Error{Errors: errs}
			ret = multierror.Append(ret, errors.Wrap(multi, "bus debug writer"))
		}
	}

	return ret
}
