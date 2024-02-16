package earthfile2llb

import (
	"context"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session/localhost"
	"github.com/pkg/errors"
)

type withDockerRunLocalTar struct {
	c   *Converter
	sem semutil.Semaphore

	enableParallel bool

	mu       sync.Mutex
	tarLoads []tarLoadLocal
}

func newWithDockerRunLocal(c *Converter, enableParallel bool) *withDockerRunLocalTar {
	// This semaphore ensures that there is at least one thread allowed to progress,
	// even if parallelism is completely starved.
	sem := semutil.NewMultiSem(c.opt.Parallelism, semutil.NewWeighted(1))

	return &withDockerRunLocalTar{
		c:              c,
		sem:            sem,
		enableParallel: enableParallel,
	}
}

type tarLoadLocal struct {
	imgName string
	imgFile string
}

func (wdrl *withDockerRunLocalTar) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := wdrl.c.checkAllowed(runCmd)
	if err != nil {
		return err
	}
	wdrl.c.nonSaveCommand()

	// Build and solve images to be loaded.
	loadPromises := make([]chan DockerLoadOpt, 0, len(opt.Loads))
	for _, loadOpt := range opt.Loads {
		lp, err := wdrl.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
		loadPromises = append(loadPromises, lp)
	}
	for _, lp := range loadPromises {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-lp:
		}
	}
	// Sort the tar list, to make the operation consistent.
	sort.Slice(wdrl.tarLoads, func(i, j int) bool {
		return wdrl.tarLoads[i].imgName < wdrl.tarLoads[j].imgName
	})
	// Issue docker load.
	for _, tl := range wdrl.tarLoads {
		runOpts := []llb.RunOption{
			llb.IgnoreCache,
			llb.Args([]string{localhost.RunOnLocalHostMagicStr, "/bin/sh", "-c", wdrl.c.containerFrontend.ImageLoadFromFileCommand(tl.imgFile)}),
		}
		wdrl.c.mts.Final.MainState = wdrl.c.mts.Final.MainState.Run(runOpts...).Root()
	}

	crOpts := ConvertRunOpts{
		CommandName:          "WITH DOCKER RUN",
		Locally:              true,
		Args:                 args,
		Mounts:               opt.Mounts,
		Secrets:              opt.Secrets,
		WithEntrypoint:       opt.WithEntrypoint,
		WithShell:            opt.WithShell,
		NoCache:              opt.NoCache,
		Interactive:          opt.Interactive,
		InteractiveKeep:      opt.interactiveKeep,
		InteractiveSaveFiles: opt.TryCatchSaveArtifacts,
	}

	// then finally run the command
	_, err = wdrl.c.internalRun(ctx, crOpts)
	if err != nil {
		return err
	}
	return nil
}

func (wdrl *withDockerRunLocalTar) load(ctx context.Context, opt DockerLoadOpt) (chan DockerLoadOpt, error) {
	optPromise := make(chan DockerLoadOpt, 1)
	depTarget, err := domain.ParseTarget(opt.Target)
	if err != nil {
		return nil, errors.Wrapf(err, "parse target %s", opt.Target)
	}
	afterFun := func(ctx context.Context, mts *states.MultiTarget) error {
		if opt.ImageName == "" {
			// Infer image name from the SAVE IMAGE statement.
			if len(mts.Final.SaveImages) == 0 || mts.Final.SaveImages[0].DockerTag == "" {
				return errors.New(
					"no docker image tag specified in load and it cannot be inferred from the SAVE IMAGE statement")
			}
			if len(mts.Final.SaveImages) > 1 {
				return errors.New(
					"no docker image tag specified in load and it cannot be inferred from the SAVE IMAGE statement: " +
						"multiple tags mentioned in SAVE IMAGE")
			}
			opt.ImageName = mts.Final.SaveImages[0].DockerTag
		}
		err := wdrl.solveImage(
			ctx, mts, depTarget.String(), opt.ImageName,
			llb.WithCustomNamef(
				"%sDOCKER LOAD %s %s", wdrl.c.imageVertexPrefix(opt.ImageName, opt.Platform), depTarget.String(), opt.ImageName))
		if err != nil {
			return err
		}
		optPromise <- opt
		return nil
	}
	if wdrl.enableParallel {
		err = wdrl.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, loadCmd, afterFun, wdrl.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := wdrl.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, false, loadCmd, "")
		if err != nil {
			return nil, err
		}
		err = afterFun(ctx, mts)
		if err != nil {
			return nil, err
		}
	}
	return optPromise, nil
}

func (wdrl *withDockerRunLocalTar) solveImage(ctx context.Context, mts *states.MultiTarget, opName string, dockerTag string, opts ...llb.RunOption) error {
	outDir, err := os.MkdirTemp(os.TempDir(), "earthly-docker-load")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for docker load")
	}
	wdrl.c.opt.CleanCollection.Add(func() error {
		return os.RemoveAll(outDir)
	})
	outFile := path.Join(outDir, "image.tar")
	err = wdrl.c.opt.DockerImageSolverTar.SolveImage(ctx, mts, dockerTag, outFile, !wdrl.c.ftrs.NoTarBuildOutput)
	if err != nil {
		return errors.Wrapf(err, "build target %s for docker load", opName)
	}
	wdrl.mu.Lock()
	defer wdrl.mu.Unlock()
	wdrl.tarLoads = append(wdrl.tarLoads, tarLoadLocal{
		imgName: dockerTag,
		imgFile: outFile,
	})
	return nil
}
