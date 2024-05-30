package earthfile2llb

import (
	"context"
	"fmt"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/pkg/errors"
)

type withDockerRunLocalReg struct {
	c              *Converter
	sem            semutil.Semaphore
	enableParallel bool
}

func newWithDockerRunLocalReg(c *Converter, enableParallel bool) *withDockerRunLocalReg {
	// This semaphore ensures that there is at least one thread allowed to progress,
	// even if parallelism is completely starved.
	sem := semutil.NewMultiSem(c.opt.Parallelism, semutil.NewWeighted(1))

	return &withDockerRunLocalReg{
		c:              c,
		enableParallel: enableParallel,
		sem:            sem,
	}
}

func (w *withDockerRunLocalReg) Run(ctx context.Context, args []string, opt WithDockerOpt) (retErr error) {
	err := w.c.checkAllowed(runCmd)
	if err != nil {
		return err
	}
	w.c.nonSaveCommand()

	cmdID, cmd, err := w.c.newLogbusCommand(ctx, "WITH DOCKER RUN")
	if err != nil {
		return errors.Wrap(err, "failed to create command")
	}

	defer func() {
		cmd.SetEndError(retErr)
	}()

	var imagesToBuild []*states.ImageDef

	// Build and solve images to be loaded.
	imageDefChans := make([]chan *states.ImageDef, 0, len(opt.Loads))
	for _, loadOpt := range opt.Loads {
		imageDefChan, err := w.load(ctx, cmdID, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
		imageDefChans = append(imageDefChans, imageDefChan)
	}
	for _, imageDefChan := range imageDefChans {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case imageDef := <-imageDefChan:
			imagesToBuild = append(imagesToBuild, imageDef)
		}
	}

	res, err := w.c.opt.MultiImageSolver.SolveImages(ctx, imagesToBuild)
	if err != nil {
		return errors.Wrap(err, "solving images")
	}
	defer res.ReleaseFunc()

	// Forward any build errors to the existing ErrGroup, which will handle display.
	w.c.opt.ErrorGroup.Go(func() error {
		for {
			select {
			case err, ok := <-res.ErrChan:
				if !ok {
					return nil
				}
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	// Wait for all images to build (channel will be closed when finished).
	for result := range res.ResultChan {
		// Pull and then retag all images with expected tags.
		pullImage := fmt.Sprintf("%s/%s", w.c.opt.LocalRegistryAddr, result.IntermediateImageName)
		err := w.c.containerFrontend.ImagePull(ctx, pullImage)
		if err != nil {
			return err
		}
		err = w.c.containerFrontend.ImageTag(ctx, containerutil.ImageTag{
			SourceRef: pullImage,
			TargetRef: result.FinalImageName,
		})
		if err != nil {
			return errors.Wrapf(err, "tag image %q", result.FinalImageName)
		}
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

	_, err = w.c.internalRun(ctx, crOpts)
	if err != nil {
		return err
	}

	// Force synchronous command execution if we're using the local registry for
	// loads and pulls.
	err = w.c.forceExecution(ctx, w.c.mts.Final.MainState, w.c.platr)
	if err != nil && errors.Is(err, ErrUnlazyForceExecution) {
		// The forced error will be returned elsewhere via magic I don't understand
		// So swallowing the error here keeps error messages consistent.
		return nil
	}
	return err
}

func (w *withDockerRunLocalReg) load(ctx context.Context, cmdID string, opt DockerLoadOpt) (chan *states.ImageDef, error) {
	imageDefChan := make(chan *states.ImageDef, 1)

	depTarget, err := domain.ParseTarget(opt.Target)
	if err != nil {
		return nil, errors.Wrapf(err, "parse target %s", opt.Target)
	}

	afterFun := func(ctx context.Context, mts *states.MultiTarget) error {
		if opt.ImageName == "" {
			// Infer image name from the SAVE IMAGE statement.
			if len(mts.Final.SaveImages) == 0 || mts.Final.SaveImages[0].DockerTag == "" {
				return errNoImageTag
			}
			if len(mts.Final.SaveImages) > 1 {
				return errors.Wrap(errNoImageTag, "multiple tags mentioned in SAVE IMAGE")
			}
			opt.ImageName = mts.Final.SaveImages[0].DockerTag
		}

		imageDefChan <- &states.ImageDef{
			MTS:       mts,
			ImageName: opt.ImageName,
			Platform:  opt.Platform,
		}

		return nil
	}

	if w.enableParallel {
		err = w.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, loadCmd, afterFun, w.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := w.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, false, loadCmd, "", nil)
		if err != nil {
			return nil, err
		}
		err = afterFun(ctx, mts)
		if err != nil {
			return nil, err
		}
	}

	return imageDefChan, nil
}
