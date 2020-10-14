package earthfile2llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/logging"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/dedup"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	dockerdWrapperPath          = "/var/earthly/dockerd-wrapper.sh"
	dockerAutoInstallScriptPath = "/var/earthly/docker-auto-install.sh"
	composeConfigFile           = "compose-config.yml"
)

// DockerLoadOpt holds parameters for WITH DOCKER --load parameter.
type DockerLoadOpt struct {
	Target    string
	ImageName string
	BuildArgs []string
}

// WithDockerOpt holds parameters for WITH DOCKER run.
type WithDockerOpt struct {
	Mounts          []string
	Secrets         []string
	WithShell       bool
	WithEntrypoint  bool
	Pulls           []string
	Loads           []DockerLoadOpt
	ComposeFiles    []string
	ComposeServices []string
}

type withDockerRun struct {
	c        *Converter
	tarLoads []llb.State
}

func (wdr *withDockerRun) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	err := wdr.installDeps(ctx, opt)
	if err != nil {
		return err
	}
	// Grab relevant images from compose file(s).
	composeImages, err := wdr.getComposeImages(ctx, opt)
	if err != nil {
		return err
	}
	composeImagesSet := make(map[string]bool)
	for _, composeImg := range composeImages {
		composeImagesSet[composeImg] = true
	}
	for _, loadOpt := range opt.Loads {
		// Make sure we don't pull a compose image which is loaded.
		if composeImagesSet[loadOpt.ImageName] {
			delete(composeImagesSet, loadOpt.ImageName)
		}
		// Load.
		err := wdr.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
	}
	// Add compose images (what's left of them) to the pull list.
	for composeImg := range composeImagesSet {
		opt.Pulls = append(opt.Pulls, composeImg)
	}
	// Sort to make the operation consistent.
	sort.Strings(opt.Pulls)
	for _, pullImageName := range opt.Pulls {
		err := wdr.pull(ctx, pullImageName)
		if err != nil {
			return errors.Wrap(err, "pull")
		}
	}
	logging.GetLogger(ctx).
		With("args", args).
		With("mounts", opt.Mounts).
		With("secrets", opt.Secrets).
		With("privileged", true).
		With("withEntrypoint", opt.WithEntrypoint).
		With("push", false).
		Info("Applying WITH DOCKER RUN")
	var runOpts []llb.RunOption
	mountRunOpts, err := parseMounts(
		opt.Mounts, wdr.c.mts.Final.Target, wdr.c.mts.Final.TargetInput, wdr.c.cacheContext)
	if err != nil {
		return errors.Wrap(err, "parse mounts")
	}
	runOpts = append(runOpts, mountRunOpts...)
	runOpts = append(runOpts, llb.AddMount(
		"/var/earthly/dind", llb.Scratch(), llb.HostBind(), llb.SourcePath("/tmp/earthly/dind")))
	runOpts = append(runOpts, llb.AddMount(
		dockerdWrapperPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))
	// This seems to make earthly-in-earthly work
	// (and docker run --privileged, together with -v /sys/fs/cgroup:/sys/fs/cgroup),
	// however, it breaks regular cases.
	//runOpts = append(runOpts, llb.AddMount(
	//"/sys/fs/cgroup", llb.Scratch(), llb.HostBind(), llb.SourcePath("/sys/fs/cgroup")))
	var tarPaths []string
	for index, tarContext := range wdr.tarLoads {
		loadDir := fmt.Sprintf("/var/earthly/load-%d", index)
		runOpts = append(runOpts, llb.AddMount(loadDir, tarContext, llb.Readonly))
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
	runOpts = append(runOpts, llb.WithCustomNamef("%s%s", wdr.c.vertexPrefix(), runStr))
	dindID, err := wdr.c.mts.Final.TargetInput.Hash()
	if err != nil {
		return errors.Wrap(err, "compute dind id")
	}
	shellWrap := makeWithDockerdWrapFun(dindID, tarPaths, opt)
	return wdr.c.internalRun(ctx, finalArgs, opt.Secrets, opt.WithShell, shellWrap, false, false, runStr, runOpts...)
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
		llb.WithCustomNamef("%sWITH DOCKER (install deps)", wdr.c.vertexPrefix()),
	}
	wdr.c.mts.Final.MainState = wdr.c.mts.Final.MainState.Run(runOpts...).Root()
	return nil
}

func (wdr *withDockerRun) getComposeImages(ctx context.Context, opt WithDockerOpt) ([]string, error) {
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
		Image string `yaml:"image"`
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
	var images []string
	for serviceName, serviceInfo := range config.Services {
		if serviceInfo.Image == "" {
			// Image not specified in yaml.
			continue
		}
		if len(opt.ComposeServices) > 0 {
			if composeServicesSet[serviceName] {
				images = append(images, serviceInfo.Image)
			}
		} else {
			// No services specified. Special case: collect all.
			images = append(images, serviceInfo.Image)
		}
	}
	return images, nil
}

func (wdr *withDockerRun) pull(ctx context.Context, dockerTag string) error {
	logging.GetLogger(ctx).With("dockerTag", dockerTag).Info("Applying DOCKER PULL")
	state, image, _, err := wdr.c.internalFromClassical(
		ctx, dockerTag,
		llb.WithCustomNamef("%sDOCKER PULL %s", wdr.c.imageVertexPrefix(dockerTag), dockerTag),
	)
	if err != nil {
		return err
	}
	mts := &states.MultiTarget{
		Final: &states.SingleTarget{
			MainState: state,
			MainImage: image,
			TargetInput: dedup.TargetInput{
				TargetCanonical: fmt.Sprintf("+@docker-pull:%s", dockerTag),
			},
			SaveImages: []states.SaveImage{
				{
					State:     state,
					Image:     image,
					DockerTag: dockerTag,
				},
			},
		},
	}
	return wdr.solveImage(
		ctx, mts, dockerTag, dockerTag,
		llb.WithCustomNamef("%sDOCKER LOAD (PULL %s)", wdr.c.imageVertexPrefix(dockerTag), dockerTag))
}

func (wdr *withDockerRun) load(ctx context.Context, opt DockerLoadOpt) error {
	logging.GetLogger(ctx).With("target-name", opt.Target).With("dockerTag", opt.ImageName).Info("Applying DOCKER LOAD")
	depTarget, err := domain.ParseTarget(opt.Target)
	if err != nil {
		return errors.Wrapf(err, "parse target %s", opt.Target)
	}
	mts, err := wdr.c.Build(ctx, depTarget.String(), opt.BuildArgs)
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
	solveID, err := mts.Final.TargetInput.Hash()
	if err != nil {
		return errors.Wrap(err, "target input hash")
	}
	tarContext, found := wdr.c.solveCache[solveID]
	if found {
		wdr.tarLoads = append(wdr.tarLoads, tarContext)
		return nil
	}
	// Use a builder to create docker .tar file, mount it via a local build context,
	// then docker load it within the current side effects state.
	outDir, err := ioutil.TempDir("/tmp", "earthly-docker-load")
	if err != nil {
		return errors.Wrap(err, "mk temp dir for docker load")
	}
	wdr.c.cleanCollection.Add(func() error {
		return os.RemoveAll(outDir)
	})
	outFile := path.Join(outDir, "image.tar")
	err = wdr.c.dockerBuilderFun(ctx, mts, dockerTag, outFile)
	if err != nil {
		return errors.Wrapf(err, "build target %s for docker load", opName)
	}
	dockerImageID, err := dockertar.GetID(outFile)
	if err != nil {
		return errors.Wrap(err, "inspect docker tar after build")
	}
	// Use the docker image ID + dockerTag as sessionID. This will cause
	// buildkit to use cache when these are the same as before (eg a docker image
	// that is identical as before).
	sessionIDKey := fmt.Sprintf("%s-%s", dockerTag, dockerImageID)
	sha256SessionIDKey := sha256.Sum256([]byte(sessionIDKey))
	sessionID := hex.EncodeToString(sha256SessionIDKey[:])
	// Add the tar to the local context.
	tarContext = llb.Local(
		solveID,
		llb.SharedKeyHint(opName),
		llb.SessionID(sessionID),
		llb.Platform(llbutil.TargetPlatform),
		llb.WithCustomNamef("[internal] docker tar context %s %s", opName, sessionID),
	)
	wdr.tarLoads = append(wdr.tarLoads, tarContext)
	wdr.c.mts.Final.LocalDirs[solveID] = outDir
	wdr.c.solveCache[solveID] = tarContext
	return nil
}

func makeWithDockerdWrapFun(dindID string, tarPaths []string, opt WithDockerOpt) shellWrapFun {
	dockerRoot := path.Join("/var/earthly/dind", dindID)
	params := []string{
		fmt.Sprintf("EARTHLY_DOCKERD_DATA_ROOT=\"%s\"", dockerRoot),
		fmt.Sprintf("EARTHLY_DOCKER_LOAD_FILES=\"%s\"", strings.Join(tarPaths, " ")),
	}
	params = append(params, composeParams(opt)...)
	return func(args []string, envVars []string, isWithShell bool, withDebugger bool) []string {
		envVars2 := append(params, envVars...)
		return []string{
			"/bin/sh", "-c",
			strWithEnvVarsAndDocker(args, envVars2, isWithShell, withDebugger, true),
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
		llb.WithCustomNamef("%sWITH DOCKER (docker-compose config)", wdr.c.vertexPrefix()),
	}
	state := wdr.c.mts.Final.MainState.Run(runOpts...).Root()

	// Perform solve to output compose config. We will use that compose config to read in images.
	composeConfigState := llbutil.CopyOp(
		state, []string{fmt.Sprintf("/tmp/earthly/%s", composeConfigFile)},
		llb.Scratch().Platform(llbutil.TargetPlatform), fmt.Sprintf("/%s", composeConfigFile),
		false, false, "",
		llb.WithCustomNamef("[internal] copy %s", composeConfigFile))
	mts := &states.MultiTarget{
		Visited: wdr.c.mts.Visited,
		Final: &states.SingleTarget{
			Target:         wdr.c.mts.Final.Target,
			MainImage:      wdr.c.mts.Final.MainImage,
			MainState:      state,
			ArtifactsState: composeConfigState,
			LocalDirs:      wdr.c.mts.Final.LocalDirs,
		},
	}
	composeConfigArtifact := domain.Artifact{
		Target:   wdr.c.mts.Final.Target,
		Artifact: composeConfigFile,
	}
	outDir, err := ioutil.TempDir("/tmp", "earthly-compose-config")
	if err != nil {
		return nil, errors.Wrap(err, "mk temp dir for solve compose config")
	}
	wdr.c.cleanCollection.Add(func() error {
		return os.RemoveAll(outDir)
	})
	err = wdr.c.artifactBuilderFun(ctx, mts, composeConfigArtifact, fmt.Sprintf("%s/", outDir))
	if err != nil {
		return nil, errors.Wrapf(err, "build artifact %s", composeConfigArtifact.String())
	}
	outComposeConfig := filepath.Join(outDir, composeConfigFile)
	composeConfigDt, err := ioutil.ReadFile(outComposeConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "read %s", outComposeConfig)
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
