package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/outmon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/pullping"
	"github.com/moby/buildkit/util/entitlements"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type onImageFunc func(context.Context, *errgroup.Group, string) (io.WriteCloser, error)
type onArtifactFunc func(context.Context, int, domain.Artifact, string, string) (string, error)
type onFinalArtifactFunc func(context.Context) (string, error)
type onReadyForPullFunc func(context.Context, []string) error

type solver struct {
	sm              *outmon.SolverMonitor
	bkClient        *client.Client
	attachables     []session.Attachable
	enttlmnts       []entitlements.Entitlement
	cacheImports    *states.CacheImports
	cacheExport     string
	maxCacheExport  string
	saveInlineCache bool
}

// solveDockerWithRegistry uses buildkit to create a new image. The image is
// then pushed to a local Docker registry instead of having buildkit export a
// local tar file.
func (s *solver) solveDockerWithRegistry(ctx context.Context, state pllb.State, platform specs.Platform, img *image.Image, dockerTag string, outFile string, printOutput bool) error {
	fmt.Println("USING NEW CODE PATH")
	return nil
}

// solveDockerTar has buildkit export a Docker image as a local tar file.
func (s *solver) solveDockerTar(ctx context.Context, state pllb.State, platform specs.Platform, img *image.Image, dockerTag string, outFile string, printOutput bool) error {
	dt, err := state.Marshal(ctx, llb.Platform(platform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	pipeR, pipeW := io.Pipe()
	solveOpt, err := s.newSolveOptDocker(img, dockerTag, pipeW)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Solve(ctx, dt, *solveOpt, ch)
		if err != nil {
			return errors.Wrap(err, "solve")
		}
		return nil
	})
	var vertexFailureOutput string
	eg.Go(func() error {
		var err error
		if printOutput {
			vertexFailureOutput, err = s.sm.MonitorProgress(ctx, ch, "", true)
			return err
		}
		// Silent case.
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-ch:
				if !ok {
					return nil
				}
				// Do nothing - just consume the status updates silently.
			}
		}
	})
	eg.Go(func() error {
		file, err := os.Create(outFile)
		if err != nil {
			return errors.Wrapf(err, "open file %s for writing", outFile)
		}
		defer file.Close()
		bufFile := bufio.NewWriter(file)
		defer bufFile.Flush()
		buf := make([]byte, 1024)
		for {
			n, err := pipeR.Read(buf)
			if err != nil && err != io.EOF {
				return errors.Wrap(err, "pipe read")
			}
			if err == io.EOF {
				break
			}
			_, err = bufFile.Write(buf[:n])
			if err != nil {
				return errors.Wrap(err, "write chunk to file")
			}
		}
		return nil
	})
	go func() {
		<-ctx.Done()
		// Close read pipe on cancels, otherwise the whole thing hangs.
		pipeR.Close()
	}()
	err = eg.Wait()
	if err != nil {
		return NewBuildError(err, vertexFailureOutput)
	}
	return nil
}

func (s *solver) buildMainMulti(ctx context.Context, bf gwclient.BuildFunc, onImage onImageFunc, onArtifact onArtifactFunc, onFinalArtifact onFinalArtifactFunc, onPullCallback onReadyForPullFunc, phaseText string) error {
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	solveOpt, err := s.newSolveOptMulti(ctx, eg, onImage, onArtifact, onFinalArtifact, onPullCallback)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	var buildErr error
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Build(ctx, *solveOpt, "", bf, ch)
		if err != nil {
			// The actual error from bkClient.Build sometimes races with
			// a context cancelled in the solver monitor.
			buildErr = err
			return err
		}
		return nil
	})
	var vertexFailureOutput string
	eg.Go(func() error {
		var err error
		vertexFailureOutput, err = s.sm.MonitorProgress(ctx, ch, phaseText, false)
		return err
	})
	err = eg.Wait()
	if buildErr != nil {
		return NewBuildError(buildErr, vertexFailureOutput)
	}
	if err != nil {
		return NewBuildError(err, vertexFailureOutput)
	}
	return nil
}

func (s *solver) newSolveOptDocker(img *image.Image, dockerTag string, w io.WriteCloser) (*client.SolveOpt, error) {
	imgJSON, err := json.Marshal(img)
	if err != nil {
		return nil, errors.Wrap(err, "image json marshal")
	}
	var cacheImports []client.CacheOptionsEntry
	for ci := range s.cacheImports.AsMap() {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
	}
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: client.ExporterDocker,
				Attrs: map[string]string{
					"name":                  dockerTag,
					"containerimage.config": string(imgJSON),
				},
				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return w, nil
				},
			},
		},
		CacheImports:        cacheImports,
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

func (s *solver) newSolveOptMulti(ctx context.Context, eg *errgroup.Group, onImage onImageFunc, onArtifact onArtifactFunc, onFinalArtifact onFinalArtifactFunc, onPullCallback onReadyForPullFunc) (*client.SolveOpt, error) {
	var cacheImports []client.CacheOptionsEntry
	for ci := range s.cacheImports.AsMap() {
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
					return onImage(ctx, eg, imageName)
				},
				OutputDirFunc: func(md map[string]string) (string, error) {
					if md["export-dir"] != "true" {
						// Use the other fun for images.
						return "", nil
					}
					if md["final-artifact"] == "true" {
						return onFinalArtifact(ctx)
					}
					indexStr := md["dir-index"]
					index, err := strconv.Atoi(indexStr)
					if err != nil {
						return "", errors.Wrapf(err, "parse dir-index %s", indexStr)
					}
					artifactStr := md["artifact"]
					srcPath := md["src-path"]
					destPath := md["dest-path"]
					artifact, err := domain.ParseArtifact(artifactStr)
					if err != nil {
						return "", errors.Wrapf(err, "parse artifact %s", artifactStr)
					}
					return onArtifact(ctx, index, artifact, srcPath, destPath)
				},
				OutputPullCallback: pullping.PullCallback(onPullCallback),
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
