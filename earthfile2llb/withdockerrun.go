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

	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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
	Platform        *specs.Platform
	BuildArgs       []string
	AllowPrivileged bool
}

// DockerPullOpt holds parameters for the WITH DOCKER --pull parameter.
type DockerPullOpt struct {
	ImageName string
	Platform  *specs.Platform
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

type withDockerRun struct {
	c *Converter

	enableParallel bool

	sem      semutil.Semaphore
	mu       sync.Mutex
	tarLoads []tarLoad
}

type tarLoad struct {
	imgName  string
	platform *specs.Platform
	state    pllb.State
}

func (wdr *withDockerRun) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := wdr.c.checkAllowed(runCmd)
	if err != nil {
		return err
	}
	wdr.c.nonSaveCommand()
	// This semaphore ensures that there is at least one thread allowed to progress,
	// even if parallelism is completely starved.
	wdr.sem = semutil.NewMultiSem(wdr.c.opt.Parallelism, semutil.NewWeighted(1))

	err = wdr.installDeps(ctx, opt)
	if err != nil {
		return err
	}
	// Grab relevant images from compose file(s).
	composePulls, err := wdr.getComposePulls(ctx, opt)
	if err != nil {
		return err
	}
	type setKey struct {
		imageName   string
		platformStr string
	}
	composeImagesSet := make(map[setKey]bool)
	for _, pull := range composePulls {
		composeImagesSet[setKey{
			imageName:   pull.ImageName,
			platformStr: llbutil.PlatformToString(pull.Platform),
		}] = true
	}

	// Loads.
	loadOptPromises := make([]chan DockerLoadOpt, 0, len(opt.Loads))
	for _, loadOpt := range opt.Loads {
		optPromise, err := wdr.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
		loadOptPromises = append(loadOptPromises, optPromise)
	}
	for _, loadOptPromise := range loadOptPromises {
		select {
		case loadOpt := <-loadOptPromise:
			// Make sure we don't pull a compose image which is loaded.
			key := setKey{
				imageName:   loadOpt.ImageName, // may have changed
				platformStr: llbutil.PlatformToString(loadOpt.Platform),
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
		key := setKey{
			imageName:   pull.ImageName,
			platformStr: llbutil.PlatformToString(pull.Platform),
		}
		if composeImagesSet[key] {
			opt.Pulls = append(opt.Pulls, pull)
		}
	}
	// Pulls.
	pullPromises := make([]chan struct{}, 0, len(opt.Pulls))
	for _, pullImageName := range opt.Pulls {
		pullPromise, err := wdr.pull(ctx, pullImageName)
		if err != nil {
			return errors.Wrap(err, "pull")
		}
		pullPromises = append(pullPromises, pullPromise)
	}
	for _, pullPromise := range pullPromises {
		select {
		case <-pullPromise:
		case <-ctx.Done():
		}
	}
	// Sort the tar list, to make the operation consistent.
	sort.Slice(wdr.tarLoads, func(i, j int) bool {
		if wdr.tarLoads[i].imgName == wdr.tarLoads[j].imgName {
			return llbutil.PlatformToString(wdr.tarLoads[i].platform) < llbutil.PlatformToString(wdr.tarLoads[j].platform)
		}
		return wdr.tarLoads[i].imgName < wdr.tarLoads[j].imgName
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
	for index, tl := range wdr.tarLoads {
		loadDir := fmt.Sprintf("/var/earthly/load-%d", index)
		crOpts.extraRunOpts = append(crOpts.extraRunOpts, pllb.AddMount(loadDir, tl.state, llb.Readonly))
		tarPaths = append(tarPaths, path.Join(loadDir, "image.tar"))
	}

	dindID, err := wdr.c.mts.Final.TargetInput().Hash()
	if err != nil {
		return errors.Wrap(err, "compute dind id")
	}
	crOpts.shellWrap = makeWithDockerdWrapFun(dindID, tarPaths, opt)

	_, err = wdr.c.internalRun(ctx, crOpts)
	return err
}

func (wdr *withDockerRun) installDeps(ctx context.Context, opt WithDockerOpt) error {
	params := composeParams(opt)
	args := []string{
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s %s",
			strings.Join(params, " "),
			dockerAutoInstallScriptPath),
	}
	runOpts := []llb.RunOption{
		llb.AddMount(
			dockerAutoInstallScriptPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerAutoInstallScriptPath)),
		llb.Args(args),
		llb.WithCustomNamef("%sWITH DOCKER (install deps)", wdr.c.vertexPrefix(false, false, false)),
	}
	wdr.c.mts.Final.MainState = wdr.c.mts.Final.MainState.Run(runOpts...).Root()
	return nil
}

func (wdr *withDockerRun) getComposePulls(ctx context.Context, opt WithDockerOpt) ([]DockerPullOpt, error) {
	if len(opt.ComposeFiles) == 0 {
		// Quick way out. Compose not used.
		return nil, nil
	}
	// Get compose images from compose config.
	composeConfigDt, err := wdr.getComposeConfig(ctx, opt)
	if err != nil {
		return nil, err
	}
	type composeService struct {
		Image    string `yaml:"image"`
		Platform string `yaml:"platform"`
	}
	type composeData struct {
		Services map[string]composeService `yaml:"services"`
	}
	var config composeData
	err = yaml.Unmarshal(composeConfigDt, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "parse compose config for %v", opt.ComposeFiles)
	}

	// Collect relevant images from the compose config.
	composeServicesSet := make(map[string]bool)
	for _, composeService := range opt.ComposeServices {
		composeServicesSet[composeService] = true
	}
	var pulls []DockerPullOpt
	for serviceName, serviceInfo := range config.Services {
		if serviceInfo.Image == "" {
			// Image not specified in yaml.
			continue
		}
		platform := wdr.c.opt.Platform
		if serviceInfo.Platform != "" {
			p, err := platforms.Parse(serviceInfo.Platform)
			if err != nil {
				return nil, errors.Wrapf(
					err, "parse platform for image %s: %s", serviceInfo.Image, serviceInfo.Platform)
			}
			platform = &p
		}
		if len(opt.ComposeServices) > 0 {
			if composeServicesSet[serviceName] {
				pulls = append(pulls, DockerPullOpt{
					ImageName: serviceInfo.Image,
					Platform:  platform,
				})
			}
		} else {
			// No services specified. Special case: collect all.
			pulls = append(pulls, DockerPullOpt{
				ImageName: serviceInfo.Image,
				Platform:  platform,
			})
		}
	}
	return pulls, nil
}

func (wdr *withDockerRun) pull(ctx context.Context, opt DockerPullOpt) (chan struct{}, error) {
	promise := make(chan struct{})
	plat := llbutil.PlatformWithDefault(opt.Platform)
	state, image, _, err := wdr.c.internalFromClassical(
		ctx, opt.ImageName, plat,
		llb.WithCustomNamef("%sDOCKER PULL %s", wdr.c.imageVertexPrefix(opt.ImageName, opt.Platform), opt.ImageName),
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
			Platform: opt.Platform,
		},
	}
	solveFun := func() error {
		err := wdr.solveImage(
			ctx, mts, opt.ImageName, opt.ImageName,
			llb.WithCustomNamef("%sDOCKER LOAD (PULL %s)", wdr.c.imageVertexPrefix(opt.ImageName, opt.Platform), opt.ImageName))
		if err != nil {
			return err
		}
		close(promise)
		return nil
	}
	if wdr.enableParallel {
		wdr.c.opt.ErrorGroup.Go(func() error {
			release, err := wdr.sem.Acquire(ctx, 1)
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

func (wdr *withDockerRun) load(ctx context.Context, opt DockerLoadOpt) (chan DockerLoadOpt, error) {
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
		err := wdr.solveImage(
			ctx, mts, depTarget.String(), opt.ImageName,
			llb.WithCustomNamef(
				"%sDOCKER LOAD %s %s", wdr.c.imageVertexPrefix(opt.ImageName, mts.Final.Platform), depTarget.String(), opt.ImageName))
		if err != nil {
			return err
		}
		optPromise <- opt
		return nil
	}
	if wdr.enableParallel {
		err = wdr.c.BuildAsync(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, loadCmd, afterFun, wdr.sem)
		if err != nil {
			return nil, err
		}
	} else {
		mts, err := wdr.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, false, loadCmd)
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

func (wdr *withDockerRun) solveImage(ctx context.Context, mts *states.MultiTarget, opName string, dockerTag string, opts ...llb.RunOption) error {
	solveID, err := states.KeyFromHashAndTag(mts.Final, dockerTag)
	if err != nil {
		return errors.Wrap(err, "state key func")
	}
	tarContext, err := wdr.c.opt.SolveCache.Do(ctx, solveID, func(ctx context.Context, _ states.StateKey) (pllb.State, error) {
		// Use a builder to create docker .tar file, mount it via a local build context,
		// then docker load it within the current side effects state.
		outDir, err := os.MkdirTemp(os.TempDir(), "earthly-docker-load")
		if err != nil {
			return pllb.State{}, errors.Wrap(err, "mk temp dir for docker load")
		}
		wdr.c.opt.CleanCollection.Add(func() error {
			return os.RemoveAll(outDir)
		})
		outFile := path.Join(outDir, "image.tar")
		err = wdr.c.opt.DockerBuilderFun(ctx, mts, dockerTag, outFile)
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
			llb.Platform(llbutil.DefaultPlatform()),
			llb.WithCustomNamef("%sdocker tar context %s %s", wdr.c.vertexPrefix(false, false, true), opName, sessionID),
		)
		// Add directly to build context so that if a later statement forces execution, the images are available.
		wdr.c.opt.BuildContextProvider.AddDir(string(solveID), outDir)
		return tarContext, nil
	})
	if err != nil {
		return err
	}
	wdr.mu.Lock()
	defer wdr.mu.Unlock()
	wdr.tarLoads = append(wdr.tarLoads, tarLoad{
		imgName:  dockerTag,
		platform: mts.Final.Platform,
		state:    tarContext,
	})
	return nil
}

func makeWithDockerdWrapFun(dindID string, tarPaths []string, opt WithDockerOpt) shellWrapFun {
	dockerRoot := path.Join("/var/earthly/dind", dindID)
	params := []string{
		fmt.Sprintf("EARTHLY_DOCKERD_DATA_ROOT=\"%s\"", dockerRoot),
		fmt.Sprintf("EARTHLY_DOCKER_LOAD_FILES=\"%s\"", strings.Join(tarPaths, " ")),
	}
	params = append(params, composeParams(opt)...)
	return func(args []string, envVars []string, isWithShell, withDebugger, forceDebugger bool) []string {
		envVars2 := append(params, envVars...)
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars2, isWithShell, withDebugger, forceDebugger, true, "", ""),
		}
	}
}

func (wdr *withDockerRun) getComposeConfig(ctx context.Context, opt WithDockerOpt) ([]byte, error) {
	// Add the right run to fetch the docker compose config.
	params := composeParams(opt)
	args := []string{
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s %s get-compose-config",
			strings.Join(params, " "),
			dockerdWrapperPath),
	}
	runOpts := []llb.RunOption{
		llb.AddMount(
			dockerdWrapperPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)),
		llb.Args(args),
		llb.WithCustomNamef("%sWITH DOCKER (docker-compose config)", wdr.c.vertexPrefix(false, false, false)),
	}
	state := wdr.c.mts.Final.MainState.Run(runOpts...).Root()
	ref, err := llbutil.StateToRef(ctx, wdr.c.opt.GwClient, state, wdr.c.opt.NoCache, wdr.c.opt.Platform, wdr.c.opt.CacheImports.AsMap())
	if err != nil {
		return nil, errors.Wrap(err, "state to ref compose config")
	}
	composeConfigDt, err := ref.ReadFile(ctx, gwclient.ReadRequest{
		Filename: fmt.Sprintf("/tmp/earthly/%s", composeConfigFile),
	})
	if err != nil {
		return nil, errors.Wrap(err, "read compose config file")
	}
	return composeConfigDt, nil
}

func composeParams(opt WithDockerOpt) []string {
	return []string{
		fmt.Sprintf("EARTHLY_START_COMPOSE=\"%t\"", (len(opt.ComposeFiles) > 0)),
		fmt.Sprintf("EARTHLY_COMPOSE_FILES=\"%s\"", strings.Join(opt.ComposeFiles, " ")),
		fmt.Sprintf("EARTHLY_COMPOSE_SERVICES=\"%s\"", strings.Join(opt.ComposeServices, " ")),
		// fmt.Sprintf("EARTHLY_DEBUG=\"true\""),
	}
}
