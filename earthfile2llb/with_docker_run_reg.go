package earthfile2llb

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/logbus/solvermon"
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

const internalWithDockerSecretPrefix = "52804da5-2787-46ad-8478-80c50f305e76"

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

func (w *withDockerRunRegistry) prepareImages(ctx context.Context, cmdID string, opt *WithDockerOpt) ([]*states.ImageDef, error) {
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
		imageDefChan, err := w.load(ctx, cmdID, loadOpt)
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

func (w *withDockerRunRegistry) Run(ctx context.Context, args []string, opt WithDockerOpt) (retErr error) {
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
		if retErr != nil {
			message := solvermon.FormatError("WITH DOCKER RUN", retErr.Error())
			cmd.SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_FAILURE, message)
		}
	}()

	err = w.installDeps(ctx, opt)
	if err != nil {
		return err
	}

	imagesToBuild, err := w.prepareImages(ctx, cmdID, &opt)
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
	results, err := readImgResults(ctx, res.ResultChan)
	if err != nil {
		return errors.Wrap(err, "error while preparing WITH DOCKER images")
	}

	// Sort the results for LLB consistency.
	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalImageNameWithDigest < results[j].FinalImageNameWithDigest
	})
	var pullImages []string
	var imgsWithDigests []string
	for _, result := range results {
		// This will be decoded in the wrapper.
		if result.NewInterImgFormat {
			pullImages = append(
				pullImages, fmt.Sprintf("%s|%s", result.IntermediateImageName, result.FinalImageName))
		} else {
			pullImages = append(pullImages, result.IntermediateImageName)
		}
		imgsWithDigests = append(imgsWithDigests, result.FinalImageNameWithDigest)
	}

	// Construct run command with all options and images.
	crOpts := ConvertRunOpts{
		CommandName:          "WITH DOCKER RUN",
		Args:                 args,
		Mounts:               opt.Mounts,
		Secrets:              opt.Secrets,
		WithEntrypoint:       opt.WithEntrypoint,
		WithShell:            opt.WithShell,
		Privileged:           true, // needed for dockerd
		WithSSH:              opt.WithSSH,
		NoCache:              opt.NoCache,
		Interactive:          opt.Interactive,
		InteractiveKeep:      opt.interactiveKeep,
		InteractiveSaveFiles: opt.TryCatchSaveArtifacts,
	}

	crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(
		"/var/earthly/dind", pllb.Scratch(), llb.HostBind(), llb.SourcePath("/tmp/earthly/dind")))
	crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(
		dockerdWrapperPath, pllb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))
	crOpts.extraRunOpts = append(crOpts.extraRunOpts, opt.extraRunOpts...)

	var dindID string
	if opt.CacheID == "" {
		dindID, err = w.c.mts.Final.TargetInput().Hash()
		if err != nil {
			return errors.Wrap(err, "make dind ID")
		}
	} else {
		// Note that the "cache_" prefix here is used to prevent auto-cleanup
		dindID = "cache_" + opt.CacheID
	}
	// We will pass along the variable EARTHLY_DOCKER_LOAD_REGISTRY via a secret
	// to prevent busting the cache, as the intermediate image names are
	// different every time.
	dockerLoadRegistrySecretID := fmt.Sprintf(
		"%s-%s-EARTHLY_DOCKER_LOAD_REGISTRY", internalWithDockerSecretPrefix, dindID)
	crOpts.extraRunOpts = append(
		crOpts.extraRunOpts,
		llb.AddSecret(
			"EARTHLY_DOCKER_LOAD_REGISTRY",
			llb.SecretID(dockerLoadRegistrySecretID),
			llb.SecretAsEnv(true),
		))
	err = w.c.opt.InternalSecretStore.SetSecret(
		ctx, dockerLoadRegistrySecretID, []byte(strings.Join(pullImages, " ")))
	if err != nil {
		return errors.Wrap(err, "set docker load registry secret")
	}
	w.c.opt.CleanCollection.Add(func() error {
		return w.c.opt.InternalSecretStore.DeleteSecret(
			context.TODO(), dockerLoadRegistrySecretID)
	})

	crOpts.shellWrap = makeWithDockerdWrapFun(dindID, nil, imgsWithDigests, opt)

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

func (w *withDockerRunRegistry) load(ctx context.Context, cmdID string, opt DockerLoadOpt) (chan *states.ImageDef, error) {
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
		err = w.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, loadCmd, afterFn, w.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := w.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, false, loadCmd, cmdID, nil)
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

func readImgResults(ctx context.Context, ch chan *states.ImageResult) ([]*states.ImageResult, error) {
	var results []*states.ImageResult
	for {
		select {
		case result, ok := <-ch:
			if !ok {
				return results, nil
			}
			results = append(results, result)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
