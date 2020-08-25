package earthfile2llb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

type withDockerRun struct {
	c        *Converter
	tarLoads []llb.State
}

func (wdr *withDockerRun) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	for _, pullImageName := range opt.Pulls {
		err := wdr.pull(ctx, pullImageName)
		if err != nil {
			return errors.Wrap(err, "pull")
		}
	}
	for _, loadOpt := range opt.Loads {
		err := wdr.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
	}

	// TODO: This does not work, because it strips away some quotes, which are valuable to the shell.
	//       In any case, this is probably working as intended as is.
	// if !isWithShell {
	// 	for i := range args {
	// 		args[i] = c.expandArgs(args[i])
	// 	}
	// }
	for i := range opt.Mounts {
		opt.Mounts[i] = wdr.c.expandArgs(opt.Mounts[i])
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
	shellWrap := makeWithDockerdWrapFun(dindID, tarPaths)
	return wdr.c.internalRun(ctx, finalArgs, opt.Secrets, opt.WithShell, shellWrap, false, runStr, runOpts...)
}

func (wdr *withDockerRun) pull(ctx context.Context, dockerTag string) error {
	dockerTag = wdr.c.expandArgs(dockerTag)
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
	targetName := wdr.c.expandArgs(opt.Target)
	dockerTag := wdr.c.expandArgs(opt.ImageName)
	for i := range opt.BuildArgs {
		opt.BuildArgs[i] = wdr.c.expandArgs(opt.BuildArgs[i])
	}
	logging.GetLogger(ctx).With("target-name", targetName).With("dockerTag", dockerTag).Info("Applying DOCKER LOAD")
	depTarget, err := domain.ParseTarget(targetName)
	if err != nil {
		return errors.Wrapf(err, "parse target %s", targetName)
	}
	mts, err := wdr.c.Build(ctx, depTarget.String(), opt.BuildArgs)
	if err != nil {
		return err
	}
	return wdr.solveImage(
		ctx, mts, depTarget.String(), dockerTag,
		llb.WithCustomNamef(
			"%sDOCKER LOAD %s %s", wdr.c.imageVertexPrefix(depTarget.String()), depTarget.String(), dockerTag))
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

func makeWithDockerdWrapFun(dindID string, tarPaths []string) shellWrapFun {
	dockerRoot := path.Join("/var/earthly/dind", dindID)
	params := []string{
		fmt.Sprintf("EARTHLY_DOCKERD_DATA_ROOT=\"%s\"", dockerRoot),
		fmt.Sprintf("EARTHLY_DOCKER_LOAD_IMAGES=\"%s\"", strings.Join(tarPaths, " ")),
	}
	return func(args []string, envVars []string, isWithShell bool, withDebugger bool) []string {
		return []string{
			"/bin/sh", "-c",
			fmt.Sprintf(
				"%s %s %s",
				strings.Join(params, " "),
				dockerdWrapperPath,
				strWithEnvVars(args, envVars, isWithShell, withDebugger)),
		}
	}
}
