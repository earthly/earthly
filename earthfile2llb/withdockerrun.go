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
	"github.com/earthly/earthly/llbutil"
	"github.com/earthly/earthly/logging"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

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
	var loadCmds []string
	for index, tarContext := range wdr.tarLoads {
		loadDir := fmt.Sprintf("/var/earthly/load-%d", index)
		runOpts = append(runOpts, llb.AddMount(loadDir, tarContext, llb.Readonly))
		loadTar := path.Join(loadDir, "image.tar")
		loadCmds = append(loadCmds, fmt.Sprintf("docker load -i %s", loadTar))
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
	dindID := cacheKey(wdr.c.mts.FinalTarget())
	shellWrap := makeWithDockerdWrapFun(dindID, loadCmds)
	return wdr.c.internalRun(ctx, finalArgs, opt.Secrets, opt.WithShell, shellWrap, false, runStr, runOpts...)
}

func (wdr *withDockerRun) pull(ctx context.Context, dockerTag string) error {
	dockerTag = wdr.c.expandArgs(dockerTag)
	logging.GetLogger(ctx).With("dockerTag", dockerTag).Info("Applying DOCKER PULL")
	state, image, _, err := wdr.c.internalFromClassical(
		ctx, dockerTag,
		llb.WithCustomNamef("%sDOCKER PULL %s", wdr.c.vertexPrefix(), dockerTag),
	)
	if err != nil {
		return err
	}
	mts := &MultiTargetStates{
		FinalStates: &SingleTargetStates{
			SideEffectsState: state,
			SideEffectsImage: image,
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
		llb.WithCustomNamef("%sDOCKER LOAD (PULL %s)", wdr.c.vertexPrefix(), dockerTag))
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
			"%sDOCKER LOAD %s %s", wdr.c.vertexPrefix(), depTarget.String(), dockerTag))
}

func (wdr *withDockerRun) solveImage(ctx context.Context, mts *MultiTargetStates, opName string, dockerTag string, opts ...llb.RunOption) error {
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
	tarContext := llb.Local(
		opName,
		llb.SharedKeyHint(opName),
		llb.SessionID(sessionID),
		llb.Platform(llbutil.TargetPlatform),
		llb.WithCustomNamef("[internal] docker tar context %s %s", opName, sessionID),
	)
	wdr.tarLoads = append(wdr.tarLoads, tarContext)
	wdr.c.mts.FinalStates.LocalDirs[opName] = outDir
	return nil
}

func makeWithDockerdWrapFun(dindID string, loadCmds []string) shellWrapFun {
	return func(args []string, envVars []string, isWithShell bool) []string {
		return []string{
			"/bin/sh", "-c",
			fmt.Sprintf(
				"/bin/sh <<EOF\n%s\nEOF", dockerdWrapCmds(args, envVars, isWithShell, dindID, loadCmds)),
		}
	}
}

func dockerdWrapCmds(args []string, envVars []string, isWithShell bool, dindID string, loadCmds []string) string {
	dockerRoot := path.Join("/var/earthly/dind", dindID)
	var cmds []string
	cmds = append(cmds, "#!/bin/sh")
	cmds = append(cmds, startDockerdCmds(dockerRoot)...)
	cmds = append(cmds, loadCmds...)
	cmds = append(cmds, strWithEnvVars(args, envVars, isWithShell))
	cmds = append(cmds, "exit_code=\"\\$?\"")
	cmds = append(cmds, stopDockerdCmds(dockerRoot)...)
	cmds = append(cmds, "exit \"\\$exit_code\"")
	return strings.Join(cmds, "\n")
}

func startDockerdCmds(dockerRoot string) []string {
	return []string{
		// Uncomment this line for debugging.
		// "set -x",
		fmt.Sprintf("mkdir -p %s", dockerRoot),
		fmt.Sprintf("dockerd-entrypoint.sh dockerd --data-root=%s &>/var/log/docker.log &", dockerRoot),
		"dockerd_pid=\"\\$!\"",
		"let i=1",
		"while ! docker ps &>/dev/null ; do",
		"sleep 1",
		"if [ \"\\$i\" -gt \"30\" ] ; then",
		// Print logs dockerd start failure.
		"cat /var/log/docker.log",
		"exit 1",
		"fi",
		"let i+=1",
		"done",
	}
}

func stopDockerdCmds(dockerRoot string) []string {
	return []string{
		"kill \"\\$dockerd_pid\" &>/dev/null",
		"let i=1",
		"while kill -0 \"\\$dockerd_pid\" &>/dev/null ; do",
		"sleep 1",
		"let i+=1",
		"if [ \"\\$i\" -gt \"10\" ]; then",
		"kill -9 \"\\$dockerd_pid\" &>/dev/null",
		"sleep 1",
		"fi",
		// Wipe the WITH DOCKER docker data after each run.
		fmt.Sprintf("rm -rf %s", dockerRoot),
		"done",
	}
}
