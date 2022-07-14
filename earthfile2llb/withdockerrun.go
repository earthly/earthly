package earthfile2llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

const (
	dockerdWrapperPath          = "/var/earthly/dockerd-wrapper.sh"
	dockerAutoInstallScriptPath = "/var/earthly/docker-auto-install.sh"
	composeConfigFile           = "compose-config.yml"
)

// DockerLoadOpt holds parameters for WITH DOCKER --load parameter.
type DockerLoadOpt struct {
	Target          string
	ImageName       string
	Platform        platutil.Platform
	BuildArgs       []string
	AllowPrivileged bool
}

// DockerPullOpt holds parameters for the WITH DOCKER --pull parameter.
type DockerPullOpt struct {
	ImageName string
	Platform  platutil.Platform
}

// WithDockerOpt holds parameters for WITH DOCKER run.
type WithDockerOpt struct {
	Mounts          []string
	Secrets         []string
	WithShell       bool
	WithEntrypoint  bool
	WithSSH         bool
	NoCache         bool
	Interactive     bool
	interactiveKeep bool
	Pulls           []DockerPullOpt
	Loads           []DockerLoadOpt
	ComposeFiles    []string
	ComposeServices []string
}

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

	var tarPaths []string
	for index, tl := range w.tarLoads {
		loadDir := fmt.Sprintf("/var/earthly/load-%d", index)
		crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(loadDir, tl.state, llb.Readonly))
		tarPaths = append(tarPaths, path.Join(loadDir, "image.tar"))
	}

	dindID, err := makeDindID(w.c.mts.Final.TargetInput(), w.c.opt.GwClient.BuildOpts().SessionID)
	if err != nil {
		return err
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
		err = w.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, loadCmd, afterFun, w.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := w.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, false, loadCmd)
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
		sessionIDKey := fmt.Sprintf("%s-%s", dockerTag, dockerImageID)
		sha256SessionIDKey := sha256.Sum256([]byte(sessionIDKey))
		sessionID := hex.EncodeToString(sha256SessionIDKey[:])

		tarContext := pllb.Local(
			string(solveID),
			llb.SessionID(sessionID),
			llb.Platform(w.c.platr.LLBNative()),
			llb.WithCustomNamef("%sdocker tar context %s %s", w.c.vertexPrefix(false, false, true), opName, sessionID),
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

func makeWithDockerdWrapFun(dindID string, tarPaths []string, pullImages []string, opt WithDockerOpt) shellWrapFun {
	dockerRoot := path.Join("/var/earthly/dind", dindID)
	params := []string{
		fmt.Sprintf("EARTHLY_DOCKERD_DATA_ROOT=\"%s\"", dockerRoot),
		fmt.Sprintf("EARTHLY_DOCKER_LOAD_FILES=\"%s\"", strings.Join(tarPaths, " ")),
		fmt.Sprintf("EARTHLY_DOCKER_LOAD_REGISTRY=\"%s\"", strings.Join(pullImages, " ")),
	}
	params = append(params, composeParams(opt)...)
	return func(args []string, envVars []string, isWithShell, withDebugger, forceDebugger bool) []string {
		envVars2 := append(params, envVars...)
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars2, isWithShell, withDebugger, forceDebugger, true, false, "", ""),
		}
	}
}

func composeParams(opt WithDockerOpt) []string {
	return []string{
		fmt.Sprintf("EARTHLY_START_COMPOSE=\"%t\"", (len(opt.ComposeFiles) > 0)),
		fmt.Sprintf("EARTHLY_COMPOSE_FILES=\"%s\"", strings.Join(opt.ComposeFiles, " ")),
		fmt.Sprintf("EARTHLY_COMPOSE_SERVICES=\"%s\"", strings.Join(opt.ComposeServices, " ")),
		// fmt.Sprintf("EARTHLY_DEBUG=\"true\""),
	}
}

func platformIncompatMsg(platr *platutil.Resolver) string {
	currentPlatStr := platr.Materialize(platr.Current()).String()
	nativePlatStr := platr.Materialize(platutil.NativePlatform).String()
	return "running WITH DOCKER as a non-native CPU architecture. This is not supported.\n" +
		fmt.Sprintf("Current platform: %s\n", currentPlatStr) +
		fmt.Sprintf("Native platform of the worker: %s\n", nativePlatStr) +
		"Try using\n\n\tFROM --platform=native earthly/dind:alpine\n\ninstead.\n" +
		"You may still --load and --pull images of a different platform.\n"
}

func makeDindID(ti dedup.TargetInput, sessionID string) (string, error) {
	hash, err := ti.Hash()
	if err != nil {
		return "", errors.Wrap(err, "hash target input")
	}
	return fmt.Sprintf("%s-%s", hash, sessionID), nil
}
