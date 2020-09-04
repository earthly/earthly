package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/earthly/earthly/earthfile2llb/image"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/logging"
	"github.com/golang/protobuf/proto"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type solver struct {
	sm          *solverMonitor
	bkClient    *client.Client
	attachables []session.Attachable
	enttlmnts   []entitlements.Entitlement
	remoteCache string
}

func (s *solver) solveDocker(ctx context.Context, localDirs map[string]string, state llb.State, img *image.Image, dockerTag string, push bool) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	pipeR, pipeW := io.Pipe()
	solveOpt, err := s.newSolveOptDocker(img, dockerTag, localDirs, pipeW)
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
		logging.GetLogger(ctx).Info("Solve successful")
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
		logging.GetLogger(ctx).Info("Docker load success")
		if push {
			err := pushDockerImage(ctx, dockerTag)
			if err != nil {
				return err
			}
			logging.GetLogger(ctx).Info("Docker push success")
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

func (s *solver) solveDockerTar(ctx context.Context, localDirs map[string]string, state llb.State, img *image.Image, dockerTag string, outFile string) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	pipeR, pipeW := io.Pipe()
	solveOpt, err := s.newSolveOptDocker(img, dockerTag, localDirs, pipeW)
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
		logging.GetLogger(ctx).Info("Solve successful")
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

func (s *solver) solveArtifacts(ctx context.Context, localDirs map[string]string, state llb.State, outDir string) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	solveOpt, err := s.newSolveOptArtifacts(outDir, localDirs)
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
		logging.GetLogger(ctx).Info("Solve successful")
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

// when printDetailed is false, we only print non-cached items
func (s *solver) solveSideEffects(ctx context.Context, localDirs map[string]string, state llb.State) error {
	dt, err := state.Marshal(ctx, llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	dtBytes, err := json.Marshal(dt)
	if err != nil {
		return errors.Wrap(err, "json marshal of state")
	}
	var ops []*pb.Op
	for _, opDef := range dt.Def {
		var op pb.Op
		err = proto.Unmarshal(opDef, &op)
		if err != nil {
			return errors.Wrap(err, "proto unmarshal of op")
		}
		ops = append(ops, &op)
	}
	opsBytes, err := json.Marshal(&ops)
	if err != nil {
		return errors.Wrap(err, "json marshal of ops")
	}
	logging.GetLogger(ctx).
		With("ops", string(opsBytes)).
		With("dt", string(dtBytes)).
		Debug("Side effectsLLB")
	solveOpt, err := s.newSolveOptSideEffects(localDirs)
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
		logging.GetLogger(ctx).Info("Solve successful")
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

func (s *solver) newSolveOptDocker(img *image.Image, dockerTag string, localDirs map[string]string, w io.WriteCloser) (*client.SolveOpt, error) {
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
		LocalDirs:           localDirs,
	}, nil
}

func (s *solver) newSolveOptArtifacts(outDir string, localDirs map[string]string) (*client.SolveOpt, error) {
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:      client.ExporterLocal,
				OutputDir: outDir,
			},
		},
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		LocalDirs:           localDirs,
	}, nil
}

func (s *solver) newSolveOptSideEffects(localDirs map[string]string) (*client.SolveOpt, error) {
	var cacheImportExport []client.CacheOptionsEntry
	if s.remoteCache != "" {
		cacheImportExport = append(cacheImportExport, newRegistryCacheOpt(s.remoteCache))
	}
	return &client.SolveOpt{
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		LocalDirs:           localDirs,
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

func loadDockerTar(ctx context.Context, r io.ReadCloser) error {
	// TODO: This is a gross hack - should use proper docker client.
	cmd := exec.CommandContext(ctx, "docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "docker load")
	}
	return nil
}

func pushDockerImage(ctx context.Context, imageName string) error {
	cmd := exec.CommandContext(ctx, "docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker push %s", imageName)
	}
	return nil
}
