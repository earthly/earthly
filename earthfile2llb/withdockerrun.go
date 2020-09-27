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

	"github.com/earthly/earthly/dockertar"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb/dedup"
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/logging"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

const dockerdWrapperPath = "/var/earthly/dockerd-wrapper.sh"

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
	if len(opt.ComposeFiles) > 0 {
		err := wdr.installDeps(ctx)
		if err != nil {
			return err
		}
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
		opt.Mounts, wdr.c.mts.FinalStates.Target, wdr.c.mts.FinalStates.TargetInput, wdr.c.cacheContext)
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
			args := make([]string, len(wdr.c.mts.FinalStates.SideEffectsImage.Config.Cmd))
			copy(args, wdr.c.mts.FinalStates.SideEffectsImage.Config.Cmd)
		}
		finalArgs = append(wdr.c.mts.FinalStates.SideEffectsImage.Config.Entrypoint, args...)
		opt.WithShell = false // Don't use shell when --entrypoint is passed.
	}
	runOpts = append(runOpts, llb.Security(llb.SecurityModeInsecure))
	runStr := fmt.Sprintf(
		"WITH DOCKER RUN %s%s",
		strIf(opt.WithEntrypoint, "--entrypoint "),
		strings.Join(finalArgs, " "))
	runOpts = append(runOpts, llb.WithCustomNamef("%s%s", wdr.c.vertexPrefix(), runStr))
	dindID, err := wdr.c.mts.FinalStates.TargetInput.Hash()
	if err != nil {
		return errors.Wrap(err, "compute dind id")
	}
	shellWrap := makeWithDockerdWrapFun(dindID, tarPaths, opt)
	return wdr.c.internalRun(ctx, finalArgs, opt.Secrets, opt.WithShell, shellWrap, false, false, runStr, runOpts...)
}

func (wdr *withDockerRun) installDeps(ctx context.Context) error {
	var runOpts []llb.RunOption
	runOpts = append(runOpts, llb.AddMount(
		dockerdWrapperPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))
	args := []string{dockerdWrapperPath, "install-deps"}
	runStr := fmt.Sprintf("WITH DOCKER (install deps)")
	return wdr.c.internalRun(ctx, args, nil, true, withShellAndEnvVars, false, false, runStr, runOpts...)
}

func (wdr *withDockerRun) getComposeImages(ctx context.Context, opt WithDockerOpt) ([]string, error) {
	var runOpts []llb.RunOption
	runOpts = append(runOpts, llb.AddMount(
		dockerdWrapperPath, llb.Scratch(), llb.HostBind(), llb.SourcePath(dockerdWrapperPath)))
	params := composeParams(opt)
	args := []string{
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s %s get-compose-config",
			strings.Join(params, " "),
			dockerdWrapperPath),
	}
	runStr := fmt.Sprintf("WITH DOCKER (detect docker-compose images)")
	// TODO: internalRun is not necessary here. Could use state.Run directly.
	err := wdr.c.internalRun(ctx, args, nil, true, withShellAndEnvVars, false, false, runStr, runOpts...)
	if err != nil {
		return nil, err
	}
	// TODO: Solve.
	// TODO: Parse & fetch service -> image map.
	// TODO: Filter out images from services not specified.
	return nil, nil
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
	mts := &MultiTargetStates{
		FinalStates: &SingleTargetStates{
			SideEffectsState: state,
			SideEffectsImage: image,
			TargetInput: dedup.TargetInput{
				TargetCanonical: fmt.Sprintf("+@docker-pull:%s", dockerTag),
			},
			SaveImages: []SaveImage{
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
	return wdr.solveImage(
		ctx, mts, depTarget.String(), opt.ImageName,
		llb.WithCustomNamef(
			"%sDOCKER LOAD %s %s", wdr.c.imageVertexPrefix(depTarget.String()), depTarget.String(), opt.ImageName))
}

func (wdr *withDockerRun) solveImage(ctx context.Context, mts *MultiTargetStates, opName string, dockerTag string, opts ...llb.RunOption) error {
	solveID, err := mts.FinalStates.TargetInput.Hash()
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
	wdr.c.mts.FinalStates.LocalDirs[solveID] = outDir
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
		return []string{
			"/bin/sh", "-c",
			fmt.Sprintf(
				"%s %s execute %s",
				strings.Join(params, " "),
				dockerdWrapperPath,
				strWithEnvVars(args, envVars, isWithShell, withDebugger)),
		}
	}
}

func composeParams(opt WithDockerOpt) []string {
	return []string{
		fmt.Sprintf("EARTHLY_START_COMPOSE=\"%t\"", (len(opt.ComposeFiles) > 0)),
		fmt.Sprintf("EARTHLY_COMPOSE_FILES=\"%s\"", strings.Join(opt.ComposeFiles, " ")),
		fmt.Sprintf("EARTHLY_COMPOSE_SERVICES=\"%s\"", strings.Join(opt.ComposeServices, " ")),
	}
}
