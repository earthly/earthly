package setup

import (
	"context"
	"os"
	"strings"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/logbus/ship"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/logbus/writersub"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// BusSetup is a helper for setting up a logbus.Bus.
type BusSetup struct {
	Bus             *logbus.Bus
	ConsoleWriter   *writersub.WriterSub
	Formatter       *formatter.Formatter
	SolverMonitor   *solvermon.SolverMonitor
	BusDebugWriter  *writersub.RawWriterSub
	LogStreamer     *ship.LogShipper
	InitialManifest *logstream.RunManifest

	logStreamerStarted bool
	verbose            bool
}

// New creates a new BusSetup.
func New(ctx context.Context, bus *logbus.Bus, debug, verbose, forceColor, noColor, disableOngoingUpdates bool, busDebugFile string, buildID string) (*BusSetup, error) {
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
		verbose: verbose,
	}
	bs.Formatter = formatter.New(ctx, bs.Bus, debug, verbose, forceColor, noColor, disableOngoingUpdates)
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

// SetOrgAndProject sets the org and project for the manifest.
func (bs *BusSetup) SetOrgAndProject(orgName, projectName string) {
	bs.InitialManifest.OrgName = orgName
	bs.InitialManifest.ProjectName = projectName
}

// LogStreamerStarted returns true if the log streamer has been started.
func (bs *BusSetup) LogStreamerStarted() bool {
	if bs == nil {
		return false
	}
	return bs.logStreamerStarted
}

// StartLogStreamer starts a LogStreamer for the given build. The
// LogStreamer streams logs to the cloud.
func (bs *BusSetup) StartLogStreamer(ctx context.Context, c *cloud.Client) {
	bs.LogStreamer = ship.NewLogShipper(c, bs.InitialManifest, bs.verbose)
	bs.LogStreamer.Start(ctx)
	bs.Bus.AddSubscriber(bs.LogStreamer)
	bs.logStreamerStarted = true
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
func (bs *BusSetup) Close() error {
	var ret error

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

	if bs.LogStreamer != nil {
		bs.LogStreamer.Close()
		if errs := bs.LogStreamer.Errs(); len(errs) > 0 {
			multi := &multierror.Error{}
			for _, err := range errs {
				streamErr := &cloud.StreamError{}
				if ok := errors.As(err, &streamErr); bs.verbose || !ok || (ok && !streamErr.Recoverable) {
					multi = multierror.Append(multi, err)
				}
			}
			ret = multierror.Append(ret, errors.Wrap(multi, "log streamer"))
		}
	}

	return ret
}
