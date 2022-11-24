package logbus

import (
	"context"
	"os"
	"strings"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/logbus/logstreamer"
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
	LogStreamer     *logstreamer.LogStreamer
	InitialManifest *logstream.RunManifest
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

// SetOrgAndProject sets the org and project for the manifest.
func (bs *BusSetup) SetOrgAndProject(orgName, projectName string) {
	bs.InitialManifest.OrgName = orgName
	bs.InitialManifest.ProjectName = projectName
}

// StartLogStreamer starts a LogStreamer for the given build. The
// LogStreamer streams logs to the cloud.
func (bs *BusSetup) StartLogStreamer(ctx context.Context, c cloud.Client) {
	bs.LogStreamer = logstreamer.New(ctx, bs.Bus, c, bs.InitialManifest)
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
	if bs.LogStreamer != nil {
		err := bs.LogStreamer.Close()
		if err != nil {
			retErr = multierror.Append(retErr, errors.Wrap(err, "log streamer"))
		}
	}
	return retErr
}
