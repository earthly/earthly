package inputgraph

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"strings"
	"time"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/features"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
)

var errNoArgSub = fmt.Errorf("arg sub supported is not yet supported by stub converter")

type StubConverter struct {
	conslog conslogging.ConsoleLogger
	target  domain.Target
}

func (c *StubConverter) From(ctx context.Context, imageName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) error {
	if strings.Contains(imageName, "+") {
		return c.propigateTarget(ctx, imageName, platform, allowPrivileged, buildArgs)
	}
	fmt.Printf("TODO hash FROM %s\n", imageName)
	return nil
}

func (c *StubConverter) FromDockerfile(ctx context.Context, contextPath string, dfPath string, dfTarget string, platform platutil.Platform, buildArgs []string) error {
	panic("not supported")
	return nil
}

func (c *StubConverter) Locally(ctx context.Context) error {
	return nil
}

func (c *StubConverter) CopyArtifactLocal(ctx context.Context, artifactName string, dest string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, isDir bool) error {
	return nil
}

func (c *StubConverter) CopyArtifact(ctx context.Context, artifactName string, dest string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, isDir bool, keepTs bool, keepOwn bool, chown string, chmod *fs.FileMode, ifExists, symlinkNoFollow bool) error {
	return nil
}

func (c *StubConverter) CopyClassical(ctx context.Context, srcs []string, dest string, isDir bool, keepTs bool, keepOwn bool, chown string, chmod *fs.FileMode, ifExists bool) error {
	fmt.Printf("TODO hash COPY %v\n", srcs)
	return nil
}

func (c *StubConverter) Run(ctx context.Context, opts earthfile2llb.ConvertRunOpts) error {
	fmt.Printf("TODO hash RUN %+v\n", opts)
	return nil
}

func (c *StubConverter) RunExitCode(ctx context.Context, opts earthfile2llb.ConvertRunOpts) (int, error) {
	return 1, nil
}
func (c *StubConverter) RunExpression(ctx context.Context, expressionName string, opts earthfile2llb.ConvertRunOpts) (string, error) {
	return "", nil
}
func (c *StubConverter) RunCommand(ctx context.Context, commandName string, opts earthfile2llb.ConvertRunOpts) (string, error) {
	return "", nil
}
func (c *StubConverter) SaveArtifact(ctx context.Context, saveFrom string, saveTo string, saveAsLocalTo string, keepTs bool, keepOwn bool, ifExists, symlinkNoFollow, force bool, isPush bool) error {
	return nil
}
func (c *StubConverter) SaveArtifactFromLocal(ctx context.Context, saveFrom, saveTo string, keepTs, keepOwn bool, chown string) error {
	return nil
}
func (c *StubConverter) PushWaitBlock(ctx context.Context) error {
	return nil
}
func (c *StubConverter) PopWaitBlock(ctx context.Context) error {
	return nil
}

func (c *StubConverter) SaveImage(ctx context.Context, imageNames []string, pushImages bool, insecurePush bool, cacheHint bool, cacheFrom []string, noManifestList bool) error {
	return nil
}

func (c *StubConverter) Build(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) error {
	return c.propigateTarget(ctx, fullTargetName, platform, allowPrivileged, buildArgs)
}

func (c *StubConverter) propigateTarget(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string) error {
	relTarget, err := domain.ParseTarget(fullTargetName)
	if err != nil {
		panic(err)
	}
	if relTarget.IsRemote() {
		panic("remote not supported")
	}

	targetRef, err := domain.JoinReferences(c.target, relTarget)
	if err != nil {
		panic(err)
	}
	target := targetRef.(domain.Target)

	fmt.Printf("TODO pull in %s deps\n", target.String())
	Load(ctx, target, c.conslog)

	return nil
}
func (c *StubConverter) BuildAsync(ctx context.Context, fullTargetName string, platform platutil.Platform, allowPrivileged bool, buildArgs []string, cmdT earthfile2llb.CmdType, apf earthfile2llb.AfterParallelFunc, sem semutil.Semaphore) error {
	return nil
}
func (c *StubConverter) Workdir(ctx context.Context, workdirPath string) error {
	return nil
}
func (c *StubConverter) User(ctx context.Context, user string) error {
	return nil
}
func (c *StubConverter) Cmd(ctx context.Context, cmdArgs []string, isWithShell bool) error {
	return nil
}
func (c *StubConverter) Entrypoint(ctx context.Context, entrypointArgs []string, isWithShell bool) error {
	return nil
}
func (c *StubConverter) Expose(ctx context.Context, ports []string) error {
	return nil
}
func (c *StubConverter) Volume(ctx context.Context, volumes []string) error {
	return nil
}
func (c *StubConverter) Env(ctx context.Context, envKey string, envValue string) error {
	return nil
}
func (c *StubConverter) Arg(ctx context.Context, argKey string, defaultArgValue string, opts earthfile2llb.ArgOpts) error {
	return nil
}
func (c *StubConverter) Let(ctx context.Context, key string, value string) error {
	return nil
}
func (c *StubConverter) UpdateArg(ctx context.Context, argKey string, argValue string, isBase bool) error {
	return nil
}
func (c *StubConverter) SetArg(ctx context.Context, argKey string, argValue string) error {
	return nil
}
func (c *StubConverter) UnsetArg(ctx context.Context, argKey string) error {
	return nil
}
func (c *StubConverter) Label(ctx context.Context, labels map[string]string) error {
	return nil
}
func (c *StubConverter) GitClone(ctx context.Context, gitURL string, branch string, dest string, keepTs bool) error {
	return nil
}
func (c *StubConverter) WithDockerRun(ctx context.Context, args []string, opt earthfile2llb.WithDockerOpt, allowParallel bool) error {
	return nil
}
func (c *StubConverter) WithDockerRunLocal(ctx context.Context, args []string, opt earthfile2llb.WithDockerOpt, allowParallel bool) error {
	return nil
}
func (c *StubConverter) Healthcheck(ctx context.Context, isNone bool, cmdArgs []string, interval time.Duration, timeout time.Duration, startPeriod time.Duration, retries int) error {
	return nil
}
func (c *StubConverter) Import(ctx context.Context, importStr, as string, isGlobal, currentlyPrivileged, allowPrivilegedFlag bool) error {
	return nil
}
func (c *StubConverter) Cache(ctx context.Context, mountTarget string, sharing string) error {
	return nil
}
func (c *StubConverter) Host(ctx context.Context, hostname string, ip net.IP) error {
	return nil
}
func (c *StubConverter) Project(ctx context.Context, org, project string) error {
	return nil
}
func (c *StubConverter) Pipeline(ctx context.Context) error {
	return nil
}
func (c *StubConverter) ResolveReference(ctx context.Context, ref domain.Reference) (bc *buildcontext.Data, allowPrivileged, allowPrivilegedSet bool, err error) {
	return nil, false, false, nil
}
func (c *StubConverter) EnterScopeDo(ctx context.Context, command domain.Command, baseTarget domain.Target, allowPrivileged bool, scopeName string, buildArgs []string) error {
	return nil
}
func (c *StubConverter) ExitScope(ctx context.Context) error {
	return nil
}
func (c *StubConverter) StackString() string {
	return ""
}
func (c *StubConverter) FinalizeStates(ctx context.Context) (*states.MultiTarget, error) {
	return nil, nil
}
func (c *StubConverter) RecordTargetFailure(ctx context.Context, err error) {
}
func (c *StubConverter) ExpandArgs(ctx context.Context, runOpts earthfile2llb.ConvertRunOpts, word string, allowShellOut bool) (string, error) {
	if strings.Contains(word, "$") {
		return "", errNoArgSub
	}
	return word, nil
}

func (c *StubConverter) PlatrParse(s string) (platutil.Platform, error) {
	//panic("PlatrParse not yet working")
	return platutil.Platform{}, nil
}

func (c *StubConverter) Ftrs() *features.Features {
	return nil
}
