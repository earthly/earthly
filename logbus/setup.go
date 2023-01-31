package logbus

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/logbus/logbus"
	"github.com/earthly/earthly/logbus/logstreamer"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/logbus/writersub"
	"github.com/earthly/earthly/util/deltautil"
	"github.com/hashicorp/go-multierror"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// TODO: Document interface
type Logstream interface {
	// From LogBus setup
	SetDefaultPlatform(platform string)
	GetBuildID() string
	GetBuildURL() string
	StartLogStreamer(ctx context.Context, c cloud.Client, orgName, projectName string)

	Close() error

	// TODO: Consider whether we can delegate - do we need the full structure?

	// TODO: Rename?
	MonitorProgress(ctx context.Context, ch chan *client.SolveStatus) error

	// From logbus
	GenericWriter() *logbus.Generic
	SetFatalError(end time.Time, targetID string, commandID string, failureType logstream.FailureType, errString string)
	SetEnd(end time.Time, status logstream.RunStatus)
	SetStart(start time.Time)

	StartNewTarget(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, initialPlatform string, runner string) (*logbus.Target, error)
}

type logstreamFacade struct {
	args            *LogstreamArgs
	bus             *logbus.Bus
	consoleWriter   *writersub.WriterSub
	formatter       *formatter.Formatter
	solverMonitor   *solvermon.SolverMonitor
	busDebugWriter  *writersub.RawWriterSub
	logStreamer     *logstreamer.LogStreamer
	initialManifest *logstream.RunManifest
}

func (l *logstreamFacade) SetFatalError(end time.Time, targetID string, commandID string, failureType logstream.FailureType, errString string) {
	l.Run().SetFatalError(end, targetID, commandID, failureType, errString)
}

func (l *logstreamFacade) SetEnd(end time.Time, status logstream.RunStatus) {
	l.Run().SetEnd(end, status)
}

func (l *logstreamFacade) SetStart(start time.Time) {
	l.Run().SetStart(start)
}

func (l *logstreamFacade) MonitorProgress(ctx context.Context, ch chan *client.SolveStatus) error {
	return l.solverMonitor.MonitorProgress(ctx, ch)
}

// SetDefaultPlatform sets the default platform of the build.
func (l *logstreamFacade) SetDefaultPlatform(platform string) {
	l.formatter.SetDefaultPlatform(platform)
}

// GetBuildID returns the buildID logstream was initialized with
func (l *logstreamFacade) GetBuildID() string {
	return l.args.BuildID
}

func (l *logstreamFacade) GetBuildURL() string {
	return fmt.Sprintf(path.Join(l.args.CIHost, "builds", l.GetBuildID()))
}

// StartLogStreamer starts a LogStreamer for the given build.
// The LogStreamer streams logs to the cloud - only if upload streaming is enabled
func (l *logstreamFacade) StartLogStreamer(ctx context.Context, c cloud.Client, orgName, projectName string) {
	if l.args.UploadLogstream {
		l.initialManifest.OrgName = orgName
		l.initialManifest.ProjectName = projectName
		l.logStreamer.StartStreaming(ctx, c)
		l.args.ConsolePrinter.Printf("Streaming logs to %s\n", l.GetBuildURL())
	}
}

func (l *logstreamFacade) GetSolverMonitor() *solvermon.SolverMonitor {
	return l.solverMonitor
}

func (l *logstreamFacade) Run() *logbus.Run {
	return l.bus.Run()
}

func (l *logstreamFacade) GenericWriter() *logbus.Generic {
	return l.bus.Run().Generic()
}

func (l *logstreamFacade) StartNewTarget(targetID, shortTargetName, canonicalTargetName string, overrideArgs []string, initialPlatform string, runner string) (*logbus.Target, error) {
	target, err := l.bus.Run().NewTarget(targetID, shortTargetName, canonicalTargetName, overrideArgs, initialPlatform, runner)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new target")
	}
	target.SetStart(time.Now())
	return target, nil
}

type Printer interface {
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}

type LogstreamArgs struct {
	BuildID                    string
	CIHost                     string
	Debug                      bool
	Verbose                    bool
	ForceColor                 bool
	NoColor                    bool
	DisableOngoingUpdates      bool
	UseLogstream               bool
	UploadLogstream            bool
	LogstreamDebugFile         string
	LogstreamDebugManifestFile string
	ConsolePrinter             Printer
}

// LogstreamFactory sets up all dependencies necessary to run Logstream
func LogstreamFactory(ctx context.Context, args *LogstreamArgs) (Logstream, error) {
	bus := logbus.New()
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
