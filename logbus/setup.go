package logbus

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
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

type Logstream interface {
	// From LogBus setup
	SetDefaultPlatform(platform string)
	GetBuildID() string
	SetOrgAndProject(orgName, projectName string)
	StartLogStreamer(ctx context.Context, c cloud.Client)

	DumpManifestToFile(path string) error
	Close() error

	// TODO: Consider whether we can delegate - do we need the full structure?
	GetSolverMonitor() *solvermon.SolverMonitor

	// From logbus
	Run() *Run
	StartNewTarget(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, initialPlatform string, runner string) (*Target, error)
}

type logstreamFacade struct {
	args            *LogstreamArgs
	bus             *Bus
	consoleWriter   *writersub.WriterSub
	formatter       *formatter.Formatter
	solverMonitor   *solvermon.SolverMonitor
	busDebugWriter  *writersub.RawWriterSub
	logStreamer     *logstreamer.LogStreamer
	initialManifest *logstream.RunManifest
}

// SetDefaultPlatform sets the default platform of the build.
func (l *logstreamFacade) SetDefaultPlatform(platform string) {
	l.formatter.SetDefaultPlatform(platform)
}

// GetBuildID returns the buildID logstream was initialized with
func (l *logstreamFacade) GetBuildID() string {
	return l.args.BuildID
}

// SetOrgAndProject sets the org and project for the manifest.
func (l *logstreamFacade) SetOrgAndProject(orgName, projectName string) {
	l.initialManifest.OrgName = orgName
	l.initialManifest.ProjectName = projectName
}

// StartLogStreamer starts a LogStreamer for the given build. The
// LogStreamer streams logs to the cloud.
func (l *logstreamFacade) StartLogStreamer(ctx context.Context, c cloud.Client) {
	l.logStreamer.StartStreaming(ctx, c)
}

func (l *logstreamFacade) GetSolverMonitor() *solvermon.SolverMonitor {
	return l.solverMonitor
}

func (l *logstreamFacade) Run() *Run {
	return l.bus.Run()
}

func (l *logstreamFacade) StartNewTarget(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, initialPlatform string, runner string) (*Target, error) {
	target, err := l.bus.Run().NewTarget(targetID, shortTargetName, canonicalTargetName, overrideArgs, initialPlatform, runner)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new target")
	}
	target.SetStart(time.Now())
	return target, nil
}

type LogstreamArgs struct {
	BuildID                    string
	Debug                      bool
	Verbose                    bool
	ForceColor                 bool
	NoColor                    bool
	DisableOngoingUpdates      bool
	UseLogstream               bool
	UploadLogstream            bool
	LogstreamDebugFile         string
	LogstreamDebugManifestFile string
}

// LogstreamFactory sets up all dependencies necessary to run Logstream
func LogstreamFactory(ctx context.Context, args *LogstreamArgs) (Logstream, error) {
	bus := New()
	l := &logstreamFacade{
		args:          args,
		bus:           bus,
		consoleWriter: writersub.New(os.Stderr, "_full"),
		formatter:     nil, // set below
		solverMonitor: nil, // set below
		initialManifest: &logstream.RunManifest{
			BuildId:            args.BuildID,
			Version:            deltautil.Version,
			CreatedAtUnixNanos: uint64(bus.CreatedAt().UnixNano()),
		},
	}
	l.logStreamer = logstreamer.New(bus, l.initialManifest)
	l.formatter = formatter.New(ctx, l.bus, args.Debug, args.Verbose, args.ForceColor, args.NoColor, args.DisableOngoingUpdates)
	l.bus.AddRawSubscriber(l.formatter)
	l.bus.AddFormattedSubscriber(l.consoleWriter)
	l.solverMonitor = solvermon.New(l.bus)
	if args.LogstreamDebugFile != "" {
		f, err := os.OpenFile(args.LogstreamDebugFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open bus debug file %s", args.LogstreamDebugFile)
		}
		useJson := strings.HasSuffix(args.LogstreamDebugFile, ".json")
		l.busDebugWriter = writersub.NewRaw(f, useJson)
		l.bus.AddSubscriber(l.busDebugWriter)
	}

	return l, nil
}

// DumpManifestToFile dumps the manifest to the given file.
func (l *logstreamFacade) DumpManifestToFile(path string) error {
	m := l.formatter.Manifest()
	proto.Merge(m, l.initialManifest)
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

// Close closes Logstream.
func (l *logstreamFacade) Close() error {
	var retErr error
	cw := l.consoleWriter
	errs := cw.Errors()
	var cwErr error
	for _, err := range errs {
		cwErr = multierror.Append(cwErr, err)
	}
	if cwErr != nil {
		retErr = multierror.Append(retErr, errors.Wrap(cwErr, "console writer"))
	}
	fErr := l.formatter.Close()
	if fErr != nil {
		retErr = multierror.Append(retErr, errors.Wrap(fErr, "formatter"))
	}
	if l.busDebugWriter != nil {
		errs := l.busDebugWriter.Errors()
		var bdwErr error
		for _, err := range errs {
			bdwErr = multierror.Append(bdwErr, err)
		}
		if bdwErr != nil {
			retErr = multierror.Append(retErr, errors.Wrap(bdwErr, "bus debug writer"))
		}
	}
	if l.logStreamer != nil {
		err := l.logStreamer.Close()
		if err != nil {
			retErr = multierror.Append(retErr, errors.Wrap(err, "log streamer"))
		}
	}
	if l.args.LogstreamDebugManifestFile != "" {
		err := l.DumpManifestToFile(l.args.LogstreamDebugManifestFile)
		if err != nil {
			retErr = multierror.Append(retErr, errors.Wrap(err, "error dumping manifest"))
		}
	}
	return retErr
}
