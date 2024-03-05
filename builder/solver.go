package builder

import (
	"context"
	"io"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/outmon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/fsutilprogress"

	"github.com/moby/buildkit/client"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/pullping"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/moby/buildkit/util/grpcerrors"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// statusChanSize is used to ensure we consume all BK status messages without
// causing back-pressure that forces BK to cancel.
const statusChanSize = 500

type onImageFunc func(context.Context, *errgroup.Group, string, string, string) (io.WriteCloser, error)
type onArtifactFunc func(context.Context, string, domain.Artifact, string, string) (string, error)
type onFinalArtifactFunc func(context.Context) (string, error)

type solver struct {
	sm              *outmon.SolverMonitor
	logbusSM        *solvermon.SolverMonitor
	useLogstream    bool
	bkClient        *client.Client
	attachables     []session.Attachable
	enttlmnts       []entitlements.Entitlement
	cacheImports    *states.CacheImports
	cacheExport     string
	maxCacheExport  string
	saveInlineCache bool
}

func (s *solver) buildMainMulti(ctx context.Context, bf gwclient.BuildFunc, onImage onImageFunc, onArtifact onArtifactFunc, onFinalArtifact onFinalArtifactFunc, onPullCallback pullping.PullCallback, phaseText string, console conslogging.ConsoleLogger) error {
	ch := make(chan *client.SolveStatus, statusChanSize)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	solveOpt, err := s.newSolveOptMulti(ctx, eg, onImage, onArtifact, onFinalArtifact, onPullCallback, console)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	var buildErr error
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Build(ctx, *solveOpt, "", bf, ch)
		if err != nil {
			if grpcErr, ok := grpcerrors.AsGRPCStatus(err); ok {
				if ie, ok := earthfile2llb.FromError(errors.New(grpcErr.Message())); ok {
					err = ie
				}
			}
			// The actual error from bkClient.Build sometimes races with
			// a context cancelled in the solver monitor.
			buildErr = err
			return err
		}
		return nil
	})
	var vertexFailureOutput string
	if s.useLogstream {
		eg.Go(func() error {
			return s.logbusSM.MonitorProgress(ctx, ch)
		})
	} else {
		eg.Go(func() error {
			var err error
			vertexFailureOutput, err = s.sm.MonitorProgress(ctx, ch, phaseText, false, s.bkClient)
			return err
		})
	}
	err = eg.Wait()
	if buildErr != nil {
		return NewBuildError(buildErr, vertexFailureOutput)
	}
	if err != nil {
		return NewBuildError(err, vertexFailureOutput)
	}
	return nil
}

func (s *solver) newSolveOptMulti(ctx context.Context, eg *errgroup.Group, onImage onImageFunc, onArtifact onArtifactFunc, onFinalArtifact onFinalArtifactFunc, onPullCallback pullping.PullCallback, console conslogging.ConsoleLogger) (*client.SolveOpt, error) {
	var cacheImports []client.CacheOptionsEntry
	for _, ci := range s.cacheImports.AsSlice() {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
	}
	var cacheExports []client.CacheOptionsEntry
	if s.cacheExport != "" {
		cacheExports = append(cacheExports, newCacheExportOpt(s.cacheExport, false))
	}
	if s.maxCacheExport != "" {
		cacheExports = append(cacheExports, newCacheExportOpt(s.maxCacheExport, true))
	}
	if s.saveInlineCache {
		cacheExports = append(cacheExports, newInlineCacheOpt())
	}

	progressCB := fsutilprogress.New("", console.WithPrefix("output"))

	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:  client.ExporterEarthly,
				Attrs: map[string]string{},
				Output: func(md map[string]string) (io.WriteCloser, error) {
					if md["export-image"] != "true" {
						return nil, nil
					}
					imageName := md["image.name"]
					waitFor := md["export-image-wait-for"]
					manifestKey := md["export-image-manifest-key"]
					return onImage(ctx, eg, imageName, waitFor, manifestKey)
				},
				OutputDirFunc: func(md map[string]string) (string, error) {
					if md["export-dir"] != "true" {
						// Use the other fun for images.
						return "", nil
					}
					if md["final-artifact"] == "true" {
						return onFinalArtifact(ctx)
					}
					indexStr := md["dir-id"]
					artifactStr := md["artifact"]
					srcPath := md["src-path"]
					destPath := md["dest-path"]
					artifact, err := domain.ParseArtifact(artifactStr)
					if err != nil {
						return "", errors.Wrapf(err, "parse artifact %s", artifactStr)
					}
					return onArtifact(ctx, indexStr, artifact, srcPath, destPath)
				},
				OutputPullCallback: pullping.PullCallback(onPullCallback),
				VerboseProgressCB:  progressCB.Verbose,
			},
		},
		CacheImports:        cacheImports,
		CacheExports:        cacheExports,
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

func newCacheImportOpt(ref string) client.CacheOptionsEntry {
	registryCacheOptAttrs := make(map[string]string)
	registryCacheOptAttrs["ref"] = ref
	return client.CacheOptionsEntry{
		Type:  "registry",
		Attrs: registryCacheOptAttrs,
	}
}

func newCacheExportOpt(ref string, max bool) client.CacheOptionsEntry {
	registryCacheOptAttrs := make(map[string]string)
	registryCacheOptAttrs["ref"] = ref
	if max {
		registryCacheOptAttrs["mode"] = "max"
		registryCacheOptAttrs["image-manifest"] = "true"
		registryCacheOptAttrs["oci-mediatypes"] = "true"
	}
	return client.CacheOptionsEntry{
		Type:  "registry",
		Attrs: registryCacheOptAttrs,
	}
}

func newInlineCacheOpt() client.CacheOptionsEntry {
	return client.CacheOptionsEntry{
		Type: "inline",
	}
}
