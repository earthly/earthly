package earthfile2llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/pllb"
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
	NoCache         bool
	Interactive     bool
	interactiveKeep bool
	Pulls           []DockerPullOpt
	Loads           []DockerLoadOpt
	ComposeFiles    []string
	ComposeServices []string
}

type withDockerRun struct {
	c        *Converter
	tarLoads []pllb.State
}

func (wdr *withDockerRun) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := wdr.installDeps(ctx, opt)
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
	for _, loadOpt := range opt.Loads {
		// Make sure we don't pull a compose image which is loaded.
		key := setKey{
			imageName:   loadOpt.ImageName,
			platformStr: llbutil.PlatformToString(loadOpt.Platform),
		}
		if composeImagesSet[key] {
			delete(composeImagesSet, key)
		}
		// Load.
		err := wdr.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
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
	// Sort to make the operation consistent.
	sort.Slice(opt.Pulls, func(i, j int) bool {
		if opt.Pulls[i].ImageName == opt.Pulls[j].ImageName {
			return llbutil.PlatformToString(opt.Pulls[i].Platform) < llbutil.PlatformToString(opt.Pulls[j].Platform)
		}
		return opt.Pulls[i].ImageName < opt.Pulls[j].ImageName
	})
	for _, pullImageName := range opt.Pulls {
		err := wdr.pull(ctx, pullImageName)
		if err != nil {
			return errors.Wrap(err, "pull")
		}
	}
	var runOpts []llb.RunOption
	mountRunOpts, err := parseMounts(
		opt.Mounts, wdr.c.mts.Final.Target, wdr.c.mts.Final.TargetInput(), wdr.c.cacheContext)
	if err != nil {
		return errors.Wrap(err, "parse mounts")
	}
	runOpts = append(runOpts, mountRunOpts...)
	runOpts = append(runOpts, pllb.AddMount(
		"/var/earthly/dind", pllb.Scratch(), llb.HostBind(), llb.SourcePath("/tmp/earthly/dind")))
	runOpts = append(runOpts, pllb.AddMount(
		dockerdWrapperPath, pllb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))
	// This seems to make earthly-in-earthly work
	// (and docker run --privileged, together with -v /sys/fs/cgroup:/sys/fs/cgroup),
	// however, it breaks regular cases.
	//runOpts = append(runOpts, pllb.AddMount(
	//"/sys/fs/cgroup", pllb.Scratch(), llb.HostBind(), llb.SourcePath("/sys/fs/cgroup")))
	var tarPaths []string
	for index, tarContext := range wdr.tarLoads {
		loadDir := fmt.Sprintf("/var/earthly/load-%d", index)
		runOpts = append(runOpts, pllb.AddMount(loadDir, tarContext, llb.Readonly))
		tarPaths = append(tarPaths, path.Join(loadDir, "image.tar"))
	}

	finalArgs := args
	if opt.WithEntrypoint {
		if len(args) == 0 {
			// No args provided. Use the image's CMD.
			args := make([]string, len(wdr.c.mts.Final.MainImage.Config.Cmd))
			copy(args, wdr.c.mts.Final.MainImage.Config.Cmd)
		}
		finalArgs = append(wdr.c.mts.Final.MainImage.Config.Entrypoint, args...)
		opt.WithShell = false // Don't use shell when --entrypoint is passed.
	}
	runOpts = append(runOpts, llb.Security(llb.SecurityModeInsecure))
	runStr := fmt.Sprintf(
		"WITH DOCKER RUN %s%s",
		strIf(opt.WithEntrypoint, "--entrypoint "),
		strings.Join(finalArgs, " "))
	runOpts = append(
		runOpts,
		llb.WithCustomNamef("%s%s", wdr.c.vertexPrefix(false, opt.Interactive || opt.interactiveKeep), runStr))
	dindID, err := wdr.c.mts.Final.TargetInput().Hash()
	if err != nil {
		return errors.Wrap(err, "compute dind id")
	}
	shellWrap := makeWithDockerdWrapFun(dindID, tarPaths, opt)
	_, err = wdr.c.internalRun(
		ctx, finalArgs, opt.Secrets, opt.WithShell, shellWrap,
		false, false, false, opt.NoCache, opt.Interactive, opt.interactiveKeep, runStr, runOpts...)
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
		llb.WithCustomNamef("%sWITH DOCKER (install deps)", wdr.c.vertexPrefix(false, false)),
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

	// Collect relevant images from the comopose config.
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

func (wdr *withDockerRun) pull(ctx context.Context, opt DockerPullOpt) error {
	plat := llbutil.PlatformWithDefault(opt.Platform)
	state, image, _, err := wdr.c.internalFromClassical(
		ctx, opt.ImageName, plat,
		llb.WithCustomNamef("%sDOCKER PULL %s", wdr.c.imageVertexPrefix(opt.ImageName), opt.ImageName),
	)
	if err != nil {
		return err
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
		},
	}
	return wdr.solveImage(
		ctx, mts, opt.ImageName, opt.ImageName,
		llb.WithCustomNamef("%sDOCKER LOAD (PULL %s)", wdr.c.imageVertexPrefix(opt.ImageName), opt.ImageName))
}

func (wdr *withDockerRun) load(ctx context.Context, opt DockerLoadOpt) error {
	depTarget, err := domain.ParseTarget(opt.Target)
	if err != nil {
		return errors.Wrapf(err, "parse target %s", opt.Target)
	}
	mts, err := wdr.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, false, "LOAD")
	if err != nil {
		return err
	}
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
	return wdr.solveImage(
		ctx, mts, depTarget.String(), opt.ImageName,
		llb.WithCustomNamef(
			"%sDOCKER LOAD %s %s", wdr.c.imageVertexPrefix(depTarget.String()), depTarget.String(), opt.ImageName))
}

func (wdr *withDockerRun) solveImage(ctx context.Context, mts *states.MultiTarget, opName string, dockerTag string, opts ...llb.RunOption) error {
	solveID, err := states.KeyFromHashAndTag(mts.Final, dockerTag)
	if err != nil {
		return errors.Wrap(err, "state key func")
	}
	tarContext, err := wdr.c.opt.SolveCache.Do(ctx, solveID, func(ctx context.Context, _ states.StateKey) (pllb.State, error) {
		// Use a builder to create docker .tar file, mount it via a local build context,
		// then docker load it within the current side effects state.
		outDir, err := ioutil.TempDir(os.TempDir(), "earthly-docker-load")
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
			llb.WithCustomNamef("[internal] docker tar context %s %s", opName, sessionID),
		)
		wdr.c.mts.Final.LocalDirs[string(solveID)] = outDir
		return tarContext, nil
	})
	if err != nil {
		return err
	}
	wdr.tarLoads = append(wdr.tarLoads, tarContext)
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
			strWithEnvVarsAndDocker(args, envVars2, isWithShell, withDebugger, forceDebugger, true, ""),
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
		llb.WithCustomNamef("%sWITH DOCKER (docker-compose config)", wdr.c.vertexPrefix(false, false)),
	}
	state := wdr.c.mts.Final.MainState.Run(runOpts...).Root()
	ref, err := llbutil.StateToRef(ctx, wdr.c.opt.GwClient, state, wdr.c.opt.Platform, wdr.c.opt.CacheImports.AsMap())
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
