package earthfile2llb

import (
	"context"
	"time"
)

type commandInterpreter interface {
	From(ctx context.Context, imageName string, buildArgs []string) error
	FromDockerfile(ctx context.Context, path string, dfPath string, dfTarget string, buildArgs []string) error
	CopyArtifact(ctx context.Context, artifactName string, dest string, buildArgs []string, isDir bool) error
	CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool)
	Run(ctx context.Context, args []string, mounts []string, secretKeyValues []string, privileged bool, withEntrypoint bool, withDocker bool, isWithShell bool, pushFlag bool) error
	SaveArtifact(ctx context.Context, saveFrom string, saveTo string, saveAsLocalTo string) error
	SaveImage(ctx context.Context, imageNames []string, pushImages bool)
	Build(ctx context.Context, fullTargetName string, buildArgs []string) (*MultiTargetStates, error)
	Workdir(ctx context.Context, workdirPath string)
	User(ctx context.Context, user string)
	Cmd(ctx context.Context, cmdArgs []string, isWithShell bool)
	Entrypoint(ctx context.Context, entrypointArgs []string, isWithShell bool)
	Expose(ctx context.Context, ports []string)
	Volume(ctx context.Context, volumes []string)
	Env(ctx context.Context, envKey string, envValue string)
	Arg(ctx context.Context, argKey string, defaultArgValue string)
	Label(ctx context.Context, labels map[string]string)
	GitClone(ctx context.Context, gitURL string, branch string, dest string) error
	DockerLoad(ctx context.Context, targetName string, dockerTag string, buildArgs []string) error
	DockerPull(ctx context.Context, dockerTag string) error
	Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int)
}
