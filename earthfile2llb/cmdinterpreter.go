package earthfile2llb

import (
	"context"
	"time"
)

// DockerLoadParams holds parameters for DOCKER LOAD commands.
type DockerLoadOpt struct {
	Target    string
	ImageName string
	BuildArgs []string
}

// WithDockerOpt holds metadata related to a WITH DOCKER run.
type WithDockerOpt struct {
	Mounts         []string
	Secrets        []string
	WithShell      bool
	WithEntrypoint bool
	Pulls          []string
	Loads          []DockerLoadOpt
}

type commandInterpreter interface {
	From(ctx context.Context, imageName string, buildArgs []string) error
	FromDockerfile(ctx context.Context, path string, dfPath string, dfTarget string, buildArgs []string) error
	CopyArtifact(ctx context.Context, artifactName string, dest string, buildArgs []string, isDir bool, chown string) error
	CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool, chown string)
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
	WithDockerRun(ctx context.Context, args []string, opt WithDockerOpt) error
	DockerLoadOld(ctx context.Context, targetName string, dockerTag string, buildArgs []string) error
	DockerPullOld(ctx context.Context, dockerTag string) error
	Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int)
}
