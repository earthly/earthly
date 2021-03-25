package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states/image"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/util/entitlements"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type onImageFunc func(context.Context, *errgroup.Group, string) (io.WriteCloser, error)
type onArtifactFunc func(context.Context, int, domain.Artifact, string, string) (string, error)
type onFinalArtifactFunc func(context.Context) (string, error)

type solver struct {
	sm              *solverMonitor
	bkClient        *client.Client
	attachables     []session.Attachable
	enttlmnts       []entitlements.Entitlement
	cacheImports    map[string]bool
	cacheExport     string
	maxCacheExport  string
	saveInlineCache bool
}

func (s *solver) solveDockerTar(ctx context.Context, state llb.State, platform specs.Platform, img *image.Image, dockerTag string, outFile string) error {
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
		vertexFailureOutput, err = s.sm.monitorProgress(ctx, ch, "")
		return err
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

func (s *solver) buildMainMulti(ctx context.Context, bf gwclient.BuildFunc, onImage onImageFunc, onArtifact onArtifactFunc, onFinalArtifact onFinalArtifactFunc, phaseText string) error {
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	solveOpt, err := s.newSolveOptMulti(ctx, eg, onImage, onArtifact, onFinalArtifact)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Build(ctx, *solveOpt, "", bf, ch)
		if err != nil {
			return errors.Wrap(err, "bkClient.Build")
		}
		return nil
	})
	var vertexFailureOutput string
	eg.Go(func() error {
		var err error
		vertexFailureOutput, err = s.sm.monitorProgress(ctx, ch, phaseText)
		return err
	})
	err = eg.Wait()
	if err != nil {
		return NewBuildError(err, vertexFailureOutput)
	}
	return nil
}

func (s *solver) solveMain(ctx context.Context, state llb.State, platform specs.Platform) error {
	dt, err := state.Marshal(ctx, llb.Platform(platform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	solveOpt, err := s.newSolveOptMain()
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
		vertexFailureOutput, err = s.sm.monitorProgress(ctx, ch, "")
		return err
	})
	err = eg.Wait()
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
	for ci := range s.cacheImports {
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

func (s *solver) newSolveOptMulti(ctx context.Context, eg *errgroup.Group, onImage onImageFunc, onArtifact onArtifactFunc, onFinalArtifact onFinalArtifactFunc) (*client.SolveOpt, error) {
	var cacheImports []client.CacheOptionsEntry
	for ci := range s.cacheImports {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
	}
	var cacheExports []client.CacheOptionsEntry
	if s.cacheExport != "" {
		cacheExports = append(cacheExports, newCacheExportOpt(s.cacheExport, false))
	}
	if s.maxCacheExport != "" {
		cacheExports = append(cacheExports, newCacheExportOpt(s.cacheExport, true))
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
			},
		},
		CacheImports:        cacheImports,
		CacheExports:        cacheExports,
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

func (s *solver) newSolveOptMain() (*client.SolveOpt, error) {
	var cacheImports []client.CacheOptionsEntry
	for ci := range s.cacheImports {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
	}
	return &client.SolveOpt{
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		CacheImports:        cacheImports,
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
