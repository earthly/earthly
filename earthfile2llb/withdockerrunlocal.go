package earthfile2llb

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session/localhost"
	"github.com/pkg/errors"
)

type withDockerRunLocal struct {
	c        *Converter
	tarLoads []llb.State
}

func (wdrl *withDockerRunLocal) Run(ctx context.Context, args []string, opt WithDockerOpt) error {
	for _, loadOpt := range opt.Loads {
		// Load.
		localImageTarPath, err := wdrl.load(ctx, loadOpt)
		if err != nil {
			return errors.Wrap(err, "load")
		}
		// then issue docker load
		runOpts := []llb.RunOption{
			llb.IgnoreCache,
			llb.Args([]string{localhost.RunOnLocalHostMagicStr, "/bin/sh", "-c", fmt.Sprintf("cat %s | docker load", localImageTarPath)}),
		}
		wdrl.c.mts.Final.MainState = wdrl.c.mts.Final.MainState.Run(runOpts...).Root()
	}

	// then finally run the command
	return wdrl.c.RunLocal(ctx, args, false)
}

func (wdrl *withDockerRunLocal) load(ctx context.Context, opt DockerLoadOpt) (string, error) {
	depTarget, err := domain.ParseTarget(opt.Target)
	if err != nil {
		return "", errors.Wrapf(err, "parse target %s", opt.Target)
	}
	mts, err := wdrl.c.buildTarget(ctx, depTarget.String(), opt.Platform, opt.AllowPrivileged, opt.BuildArgs, false, loadCmd)
	if err != nil {
		return "", err
	}
	if opt.ImageName == "" {
		// Infer image name from the SAVE IMAGE statement.
		if len(mts.Final.SaveImages) == 0 || mts.Final.SaveImages[0].DockerTag == "" {
			return "", errors.New(
				"no docker image tag specified in load and it cannot be inferred from the SAVE IMAGE statement")
		}
		if len(mts.Final.SaveImages) > 1 {
			return "", errors.New(
				"no docker image tag specified in load and it cannot be inferred from the SAVE IMAGE statement: " +
					"multiple tags mentioned in SAVE IMAGE")
		}
		opt.ImageName = mts.Final.SaveImages[0].DockerTag
	}
	return wdrl.solveImage(
		ctx, mts, depTarget.String(), opt.ImageName,
		llb.WithCustomNamef(
			"%sDOCKER LOAD %s %s", wdrl.c.imageVertexPrefix(depTarget.String()), depTarget.String(), opt.ImageName))
}

func (wdrl *withDockerRunLocal) solveImage(ctx context.Context, mts *states.MultiTarget, opName string, dockerTag string, opts ...llb.RunOption) (string, error) {
	outDir, err := ioutil.TempDir(os.TempDir(), "earthly-docker-load")
	if err != nil {
		return "", errors.Wrap(err, "mk temp dir for docker load")
	}
	wdrl.c.opt.CleanCollection.Add(func() error {
		return os.RemoveAll(outDir)
	})
	outFile := path.Join(outDir, "image.tar")
	err = wdrl.c.opt.DockerBuilderFun(ctx, mts, dockerTag, outFile)
	if err != nil {
		return "", errors.Wrapf(err, "build target %s for docker load", opName)
	}
	return outFile, nil
}
