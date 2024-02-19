package earthfile2llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sort"
	"sync"

	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

type withDockerRunTar struct {
	*withDockerRunBase
	c *Converter

	enableParallel bool
	tarLoads       []tarLoad
	sem            semutil.Semaphore
	mu             sync.Mutex
}

func newWithDockerRunTar(c *Converter, enableParallel bool) *withDockerRunTar {
	// This semaphore ensures that there is at least one thread allowed to progress,
	// even if parallelism is completely starved.
	sem := semutil.NewMultiSem(c.opt.Parallelism, semutil.NewWeighted(1))

	return &withDockerRunTar{
		c:                 c,
		withDockerRunBase: &withDockerRunBase{c},
		enableParallel:    enableParallel,
		sem:               sem,
	}
}

type tarLoad struct {
	imgName  string
	platform platutil.Platform
	state    pllb.State
}

func (w *withDockerRunTar) prepareImages(ctx context.Context, opt *WithDockerOpt) error {
	// Grab relevant images from compose file(s).
	composePulls, err := w.getComposePulls(ctx, *opt)
	if err != nil {
		return err
	}

	type setKey struct {
		imageName   string
		platformStr string
	}

	composeImagesSet := make(map[setKey]bool)
	for _, pull := range composePulls {
		pull.Platform = w.c.platr.SubPlatform(pull.Platform)
		platformStr := w.c.platr.Materialize(pull.Platform).String()
		composeImagesSet[setKey{
			imageName:   pull.ImageName,
			platformStr: platformStr,
		}] = true
	}

	// Loads.
	loadOptPromises := make([]chan DockerLoadOpt, 0, len(opt.Loads))
	for _, loadOpt := range opt.Loads {
		loadOpt.Platform = w.c.platr.SubPlatform(loadOpt.Platform)
		optPromise, err := w.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
		loadOptPromises = append(loadOptPromises, optPromise)
	}
	for _, loadOptPromise := range loadOptPromises {
		select {
		case loadOpt := <-loadOptPromise:
			// Make sure we don't pull a compose image which is loaded.
			platformStr := w.c.platr.Materialize(loadOpt.Platform).String()
			key := setKey{
				imageName:   loadOpt.ImageName, // may have changed
				platformStr: platformStr,
			}
			if composeImagesSet[key] {
				delete(composeImagesSet, key)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Add compose images (what's left of them) to the pull list.
	for _, pull := range composePulls {
		pull.Platform = w.c.platr.SubPlatform(pull.Platform)
		platformStr := w.c.platr.Materialize(pull.Platform).String()
		key := setKey{
			imageName:   pull.ImageName,
			platformStr: platformStr,
		}
		if composeImagesSet[key] {
			opt.Pulls = append(opt.Pulls, pull)
		}
	}

	// Pulls.
	pullPromises := make([]chan struct{}, 0, len(opt.Pulls))
	for _, pullImageName := range opt.Pulls {
		pullPromise, err := w.pull(ctx, pullImageName)
		if err != nil {
			return errors.Wrap(err, "pull")
		}
		pullPromises = append(pullPromises, pullPromise)
	}
	for _, pullPromise := range pullPromises {
		select {
		case <-pullPromise:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (w *withDockerRunTar) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := w.c.checkAllowed(runCmd)
	if err != nil {
		return err
	}

	w.c.nonSaveCommand()

	err = w.installDeps(ctx, opt)
	if err != nil {
		return err
	}

	err = w.prepareImages(ctx, &opt)
	if err != nil {
		return err
	}

	// Sort the tar list, to make the operation consistent.
	sort.Slice(w.tarLoads, func(i, j int) bool {
		if w.tarLoads[i].imgName == w.tarLoads[j].imgName {
			return w.tarLoads[i].platform.String() < w.tarLoads[j].platform.String()
		}
		return w.tarLoads[i].imgName < w.tarLoads[j].imgName
	})

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

	// TODO: /tmp/earthly should not be hard-coded here. It should match whatever
	//       buildkit's image EARTHLY_TMP_DIR is set to.
	crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(
		"/var/earthly/dind", pllb.Scratch(), llb.HostBind(), llb.SourcePath("/tmp/earthly/dind")))
	crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(
		dockerdWrapperPath, pllb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))
	crOpts.extraRunOpts = append(crOpts.extraRunOpts, opt.extraRunOpts...)

	var tarPaths []string
	for index, tl := range w.tarLoads {
		loadDir := fmt.Sprintf("/var/earthly/load-%d", index)
		crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(loadDir, tl.state, llb.Readonly))
		tarPaths = append(tarPaths, path.Join(loadDir, "image.tar"))
	}

	dindID, err := w.c.mts.Final.TargetInput().Hash()
	if err != nil {
		return errors.Wrap(err, "make dind ID")
	}
	crOpts.shellWrap = makeWithDockerdWrapFun(dindID, tarPaths, nil, opt)

	platformIncompatible := !w.c.platr.PlatformEquals(w.c.platr.Current(), platutil.NativePlatform)
	if platformIncompatible {
		w.c.opt.Console.Warnf("Error: " + platformIncompatMsg(w.c.platr))
		return errors.New("platform incompatible")
	}

	_, err = w.c.internalRun(ctx, crOpts)
	if err != nil {
		return err
	}

	return nil
}

func (w *withDockerRunTar) pull(ctx context.Context, opt DockerPullOpt) (chan struct{}, error) {
	promise := make(chan struct{})
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
	solveFun := func() error {
		err := w.solveImage(
			ctx, mts, opt.ImageName, opt.ImageName,
			llb.WithCustomNamef("%sDOCKER LOAD (PULL %s)", w.c.imageVertexPrefix(opt.ImageName, opt.Platform), opt.ImageName))
		if err != nil {
			return err
		}
		close(promise)
		return nil
	}
	if w.enableParallel {
		w.c.opt.ErrorGroup.Go(func() error {
			release, err := w.sem.Acquire(ctx, 1)
			if err != nil {
				return errors.Wrapf(err, "acquiring parallelism semaphore for pull load %s", opt.ImageName)
			}
			defer release()
			return solveFun()
		})
	} else {
		err = solveFun()
		if err != nil {
			return nil, err
		}
	}
	return promise, nil
}

func (w *withDockerRunTar) load(ctx context.Context, opt DockerLoadOpt) (chan DockerLoadOpt, error) {
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
		err := w.solveImage(
			ctx, mts, depTarget.String(), opt.ImageName,
			llb.WithCustomNamef(
				"%sDOCKER LOAD %s %s", w.c.imageVertexPrefix(opt.ImageName, mts.Final.PlatformResolver.Current()), depTarget.String(), opt.ImageName))
		if err != nil {
			return err
		}
		optPromise <- opt
		return nil
	}
	if w.enableParallel {
		err = w.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, loadCmd, afterFun, w.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := w.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.PassArgs, opt.BuildArgs, false, loadCmd, "")
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

func (w *withDockerRunTar) solveImage(ctx context.Context, mts *states.MultiTarget, opName string, dockerTag string, opts ...llb.RunOption) error {
	solveID, err := states.KeyFromHashAndTag(mts.Final, dockerTag)
	if err != nil {
		return errors.Wrap(err, "state key func")
	}
	tarContext, err := w.c.opt.SolveCache.Do(ctx, solveID, func(ctx context.Context, _ states.StateKey) (pllb.State, error) {
		// Use a builder to create docker .tar file, mount it via a local build
		// context, then docker load it within the current side effects state.
		outDir, err := os.MkdirTemp(os.TempDir(), "earthly-docker-load")
		if err != nil {
			return pllb.State{}, errors.Wrap(err, "mk temp dir for docker load")
		}
		w.c.opt.CleanCollection.Add(func() error {
			return os.RemoveAll(outDir)
		})
		outFile := path.Join(outDir, "image.tar")
		err = w.c.opt.DockerImageSolverTar.SolveImage(ctx, mts, dockerTag, outFile, !w.c.ftrs.NoTarBuildOutput)
		if err != nil {
			return pllb.State{}, errors.Wrapf(err, "build target %s for docker load", opName)
		}
		dockerImageID, err := dockertar.GetID(outFile)
		if err != nil {
			return pllb.State{}, errors.Wrap(err, "inspect docker tar after build")
		}
		// Use the docker image ID + dockerTag as sessionID. This will cause
		// buildkit to use cache when these are the same as before (eg a docker image
		// that is identical as before).
		// Note that this is achieved by causing CacheKey() to return a consistent result under https://github.com/moby/buildkit/blob/b3e8c63a48ad8c015f5631fc1947945b229b3919/source/local/local.go#L78
		// However, this introduces a bug during the Snapshot() call where `ls.sm.Get(timeoutCtx, sessionID, false)` is called on a non-existant sessionID, which then causes
		// buildkit to wait until a context-cancelled error occurs, and then ultimately fallsback to using any available session from the group.
		sessionIDKey := fmt.Sprintf("%s-%s", dockerTag, dockerImageID)
		sha256SessionIDKey := sha256.Sum256([]byte(sessionIDKey))
		sessionID := hex.EncodeToString(sha256SessionIDKey[:])

		prefix, _, err := w.c.newVertexMeta(ctx, false, false, true, nil, true)
		if err != nil {
			return pllb.State{}, err
		}
		tarContext := pllb.Local(
			string(solveID),
			llb.SessionID(sessionID),
			llb.Platform(w.c.platr.LLBNative()),
			llb.WithCustomNamef("%sdocker tar context %s %s", prefix, opName, sessionID),
		)
		// Add directly to build context so that if a later statement forces execution, the images are available.
		w.c.opt.BuildContextProvider.AddDir(string(solveID), outDir)
		return tarContext, nil
	})
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.tarLoads = append(w.tarLoads, tarLoad{
		imgName:  dockerTag,
		platform: mts.Final.PlatformResolver.Current(),
		state:    tarContext,
	})
	return nil
}
