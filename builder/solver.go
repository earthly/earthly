package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/states/image"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type onImageFunc func(context.Context, *errgroup.Group, int, string, string) (io.WriteCloser, error)
type onArtifactFunc func(context.Context, int, domain.Artifact, string, string) (string, error)

type solver struct {
	sm          *solverMonitor
	bkClient    *client.Client
	attachables []session.Attachable
	enttlmnts   []entitlements.Entitlement
	remoteCache string
}

func (s *solver) solveDocker(ctx context.Context, state llb.State, img *image.Image, dockerTag string, push bool) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
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
	eg.Go(func() error {
		return s.sm.monitorProgress(ctx, ch)
	})
	eg.Go(func() error {
		defer pipeR.Close()
		err := loadDockerTar(ctx, pipeR)
		if err != nil {
			return errors.Wrapf(err, "load docker tar for %s", dockerTag)
		}
		if push {
			err := pushDockerImage(ctx, dockerTag)
			if err != nil {
				return err
			}
		}
		return nil
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Close read pipe on cancels, otherwise the whole thing hangs.
				pipeR.Close()
			}
		}
	}()
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *solver) solveDockerTar(ctx context.Context, state llb.State, img *image.Image, dockerTag string, outFile string) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
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
	eg.Go(func() error {
		return s.sm.monitorProgress(ctx, ch)
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
		select {
		case <-ctx.Done():
			// Close read pipe on cancels, otherwise the whole thing hangs.
			pipeR.Close()
		}
	}()
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *solver) solveArtifacts(ctx context.Context, state llb.State, outDir string) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	solveOpt, err := s.newSolveOptArtifacts(outDir)
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
	eg.Go(func() error {
		return s.sm.monitorProgress(ctx, ch)
	})
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *solver) buildMainMulti(ctx context.Context, bf gwclient.BuildFunc, onImage onImageFunc, onArtifact onArtifactFunc) error {
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	solveOpt, err := s.newSolveOptMulti(ctx, eg, onImage, onArtifact)
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
	eg.Go(func() error {
		return s.sm.monitorProgress(ctx, ch)
	})
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *solver) solveMain(ctx context.Context, state llb.State) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
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
	eg.Go(func() error {
		return s.sm.monitorProgress(ctx, ch)
	})
	err = eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *solver) newSolveOptDocker(img *image.Image, dockerTag string, w io.WriteCloser) (*client.SolveOpt, error) {
	imgJSON, err := json.Marshal(img)
	if err != nil {
		return nil, errors.Wrap(err, "image json marshal")
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
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

func (s *solver) newSolveOptArtifacts(outDir string) (*client.SolveOpt, error) {
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:      client.ExporterLocal,
				OutputDir: outDir,
			},
		},
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

func (s *solver) newSolveOptMulti(ctx context.Context, eg *errgroup.Group, onImage onImageFunc, onArtifact onArtifactFunc) (*client.SolveOpt, error) {
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:  client.ExporterEarthly,
				Attrs: map[string]string{},
				Output: func(md map[string]string) (io.WriteCloser, error) {
					if md["export-image"] != "true" {
						return nil, nil
					}
					indexStr := md["image-index"]
					index, err := strconv.Atoi(indexStr)
					if err != nil {
						return nil, errors.Wrapf(err, "parse image-index %s", indexStr)
					}
					imageName := md["image.name"]
					digest := md["containerimage.digest"]
					return onImage(ctx, eg, index, imageName, digest)
				},
				OutputDirFunc: func(md map[string]string) (string, error) {
					if md["export-dir"] != "true" {
						// Use the other fun for images.
						return "", nil
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
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

func (s *solver) newSolveOptMain() (*client.SolveOpt, error) {
	var cacheImportExport []client.CacheOptionsEntry
	if s.remoteCache != "" {
		cacheImportExport = append(cacheImportExport, newRegistryCacheOpt(s.remoteCache))
	}
	return &client.SolveOpt{
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		CacheImports:        cacheImportExport,
		CacheExports:        cacheImportExport,
	}, nil
}

func newRegistryCacheOpt(ref string) client.CacheOptionsEntry {
	registryCacheOptAttrs := make(map[string]string)
	registryCacheOptAttrs["ref"] = ref
	registryCacheOptAttrs["mode"] = "max"
	return client.CacheOptionsEntry{
		Type:  "registry",
		Attrs: registryCacheOptAttrs,
	}
}
