package earthfile2llb

import (
	"context"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

type withDockerRunRegistry struct {
	*withDockerRunBase
	c *Converter

	enableParallel bool
	sem            semutil.Semaphore
}

func newWithDockerRunRegistry(c *Converter, enableParallel bool) *withDockerRunRegistry {
	// This semaphore ensures that there is at least one thread allowed to progress,
	// even if parallelism is completely starved.
	sem := semutil.NewMultiSem(c.opt.Parallelism, semutil.NewWeighted(1))

	return &withDockerRunRegistry{
		withDockerRunBase: &withDockerRunBase{c},
		enableParallel:    enableParallel,
		c:                 c,
		sem:               sem,
	}
}

func (w *withDockerRunRegistry) prepareImages(ctx context.Context, opt *WithDockerOpt) ([]*states.ImageDef, error) {
	// Grab relevant images from compose file(s).
	composePulls, err := w.getComposePulls(ctx, *opt)
	if err != nil {
		return nil, err
	}

	var imagesToBuild []*states.ImageDef

	type setKey struct {
		imageName string
		platform  string
	}

	composeImagesSet := make(map[setKey]bool)
	for _, pull := range composePulls {
		pull.Platform = w.c.platr.SubPlatform(pull.Platform)
		composeImagesSet[setKey{
			imageName: pull.ImageName,
			platform:  w.c.platr.Materialize(pull.Platform).String(),
		}] = true
	}

	// Loads.
	imageDefChans := make([]chan *states.ImageDef, 0, len(opt.Loads))
	for _, loadOpt := range opt.Loads {
		loadOpt.Platform = w.c.platr.SubPlatform(loadOpt.Platform)
		imageDefChan, err := w.load(ctx, loadOpt)
		if err != nil {
			return nil, errors.Wrap(err, "load")
		}
		imageDefChans = append(imageDefChans, imageDefChan)
	}
	for _, imageDefChan := range imageDefChans {
		select {
		case imageDef := <-imageDefChan:
			imagesToBuild = append(imagesToBuild, imageDef)
			// Make sure we don't pull a compose image which is loaded.
			key := setKey{
				imageName: imageDef.ImageName, // may have changed
				platform:  w.c.platr.Materialize(imageDef.Platform).String(),
			}
			if composeImagesSet[key] {
				delete(composeImagesSet, key)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Add compose images (what's left of them) to the pull list.
	for _, pull := range composePulls {
		pull.Platform = w.c.platr.SubPlatform(pull.Platform)
		key := setKey{
			imageName: pull.ImageName,
			platform:  w.c.platr.Materialize(pull.Platform).String(),
		}
		if composeImagesSet[key] {
			opt.Pulls = append(opt.Pulls, pull)
		}
	}

	// Pulls.
	for _, pullOpt := range opt.Pulls {
		imageDef, err := w.pull(ctx, pullOpt)
		if err != nil {
			return nil, errors.Wrap(err, "pull")
		}
		imagesToBuild = append(imagesToBuild, imageDef)
	}

	return imagesToBuild, nil
}

func (w *withDockerRunRegistry) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := w.c.checkAllowed(runCmd)
	if err != nil {
		return err
	}

	w.c.nonSaveCommand()

	err = w.installDeps(ctx, opt)
	if err != nil {
		return err
	}

	imagesToBuild, err := w.prepareImages(ctx, &opt)
	if err != nil {
		return err
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
	var pullImages []string
	for imageName := range res.ResultChan {
		pullImages = append(pullImages, imageName)
	}

	// Construct run command with all options and images.
	crOpts := ConvertRunOpts{
		CommandName:     "WITH DOCKER RUN",
		Args:            args,
		Mounts:          opt.Mounts,
		Secrets:         opt.Secrets,
		WithEntrypoint:  opt.WithEntrypoint,
		WithShell:       opt.WithShell,
		Privileged:      true, // needed for dockerd
		WithSSH:         opt.WithSSH,
		NoCache:         opt.NoCache,
		Interactive:     opt.Interactive,
		InteractiveKeep: opt.interactiveKeep,
	}

	crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(
		"/var/earthly/dind", pllb.Scratch(), llb.HostBind(), llb.SourcePath("/tmp/earthly/dind")))
	crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(
		dockerdWrapperPath, pllb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))

	dindID, err := w.c.mts.Final.TargetInput().Hash()
	if err != nil {
		return errors.Wrap(err, "compute dind id")
	}
	crOpts.shellWrap = makeWithDockerdWrapFun(dindID, nil, pullImages, opt)

	platformIncompatible := !w.c.platr.PlatformEquals(w.c.platr.Current(), platutil.NativePlatform)
	if platformIncompatible {
		w.c.opt.Console.Warnf("Error: " + platformIncompatMsg(w.c.platr))
		return errors.New("platform incompatible")
	}

	_, err = w.c.internalRun(ctx, crOpts)
	if err != nil {
		return err
	}

	// Force synchronous command execution if we're using the local registry for
	// loads and pulls.
	return w.c.forceExecution(ctx, w.c.mts.Final.MainState, w.c.platr)
}

func (w *withDockerRunRegistry) pull(ctx context.Context, opt DockerPullOpt) (*states.ImageDef, error) {
	state, image, _, err := w.c.internalFromClassical(
		ctx, opt.ImageName, opt.Platform,
		llb.WithCustomNamef("%sDOCKER PULL %s", w.c.imageVertexPrefix(opt.ImageName, opt.Platform), opt.ImageName),
	)
	if err != nil {
		return nil, err
	}

	mts := &states.MultiTarget{
		Final: &states.SingleTarget{
			MainState: state,
			MainImage: image,
			SaveImages: []states.SaveImage{
				{
					State:     state,
					Image:     image,
					DockerTag: opt.ImageName,
				},
			},
			PlatformResolver: w.c.platr.SubResolver(opt.Platform),
		},
	}

	return &states.ImageDef{
		MTS:       mts,
		ImageName: opt.ImageName,
		Platform:  opt.Platform,
	}, nil
}

var errNoImageTag = errors.New("no docker image tag specified in load and it cannot be inferred from the SAVE IMAGE statement")

func (w *withDockerRunRegistry) load(ctx context.Context, opt DockerLoadOpt) (chan *states.ImageDef, error) {
	imageDefChan := make(chan *states.ImageDef, 1)

	depTarget, err := domain.ParseTarget(opt.Target)
	if err != nil {
		return nil, errors.Wrapf(err, "parse target %s", opt.Target)
	}

	afterFn := func(ctx context.Context, mts *states.MultiTarget) error {

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
		err = w.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, loadCmd, afterFn, w.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := w.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, false, loadCmd)
		if err != nil {
			return nil, err
		}
		err = afterFn(ctx, mts)
		if err != nil {
			return nil, err
		}
	}

	return imageDefChan, nil
}
