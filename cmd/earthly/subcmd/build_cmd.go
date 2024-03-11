package subcmd

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/cli/cli/config"
	billingpb "github.com/earthly/cloud-api/billing"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	bkclient "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth"
	dockerauthprovider "github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/localhost/localhostprovider"
	"github.com/moby/buildkit/session/socketforward/socketprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
	buildkitgitutil "github.com/moby/buildkit/util/gitutil"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/billing"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/bk"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/cmd/earthly/helper"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/terminal"
	"github.com/earthly/earthly/docker2earthly"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/inputgraph"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/util/gatewaycrafter"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil/authprovider"
	"github.com/earthly/earthly/util/llbutil/authprovider/cloudauth"
	"github.com/earthly/earthly/util/llbutil/secretprovider"
	"github.com/earthly/earthly/util/params"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/termutil"
	"github.com/earthly/earthly/variables"
)

const autoSkipPrefix = "auto-skip"

type Build struct {
	cli CLI

	buildArgs    cli.StringSlice
	platformsStr cli.StringSlice
	secrets      cli.StringSlice
	secretFiles  cli.StringSlice
	cacheFrom    cli.StringSlice
	dockerTags   cli.StringSlice
	dockerTarget string
}

func NewBuild(cli CLI) *Build {
	return &Build{
		cli: cli,
	}
}

func (a *Build) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "build",
			Usage:       "Build an Earthly target",
			Description: "Build an Earthly target.",
			Action:      a.Action,
			Flags:       a.buildFlags(),
			Hidden:      true, // Meant to be used mainly for help output.
		},
		{
			Name:        "docker-build",
			Usage:       "*beta* Build a Dockerfile without an Earthfile",
			UsageText:   "earthly [options] docker-build [--dockerfile <dockerfile-path>] [--tag=<image-tag>] [--target=<target-name>] [--platform <platform1[,platform2,...]>] <build-context-dir> [--arg1=arg-value]",
			Description: "*beta* Builds a Dockerfile without an Earthfile.",
			Action:      a.actionDockerBuild,
			Flags: append(a.buildFlags(),
				&cli.StringFlag{
					Name:        "dockerfile",
					Aliases:     []string{"f"},
					EnvVars:     []string{"EARTHLY_DOCKER_FILE"},
					Usage:       "Path to dockerfile input",
					Value:       "Dockerfile",
					Destination: &a.cli.Flags().DockerfilePath,
				},
				&cli.StringSliceFlag{
					Name:        "tag",
					Aliases:     []string{"t"},
					EnvVars:     []string{"EARTHLY_DOCKER_TAGS"},
					Usage:       "Name and tag for the built image; formatted as 'name:tag'",
					Destination: &a.dockerTags,
				},
				&cli.StringFlag{
					Name:        "target",
					EnvVars:     []string{"EARTHLY_DOCKER_TARGET"},
					Usage:       "The docker target to build in the specified dockerfile",
					Destination: &a.dockerTarget,
				},
			),
		},
	}
}

func (a *Build) Action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("build")

	if a.cli.Flags().CI {
		a.cli.Flags().NoOutput = !a.cli.Flags().Output && !a.cli.Flags().ArtifactMode && !a.cli.Flags().ImageMode
		a.cli.Flags().Strict = true

		if a.cli.Flags().InteractiveDebugging {
			return params.Errorf("unable to use --ci flag in combination with --interactive flag")
		}
	}

	if a.cli.Flags().ImageMode && a.cli.Flags().ArtifactMode {
		return params.Errorf("both image and artifact modes cannot be active at the same time")
	}
	if (a.cli.Flags().ImageMode && a.cli.Flags().NoOutput) || (a.cli.Flags().ArtifactMode && a.cli.Flags().NoOutput) {
		if a.cli.Flags().CI {
			a.cli.Flags().NoOutput = false
		} else {
			return params.Errorf("cannot use --no-output with image or artifact modes")
		}
	}
	if a.cli.Flags().InteractiveDebugging && !termutil.IsTTY() {
		return params.Errorf("A tty-terminal must be present in order to use the --interactive flag")
	}

	flagArgs, nonFlagArgs, err := variables.ParseFlagArgsWithNonFlags(cliCtx.Args().Slice())
	if err != nil {
		return errors.Wrapf(err, "parse args %s", strings.Join(cliCtx.Args().Slice(), " "))
	}

	return a.ActionBuildImp(cliCtx, flagArgs, nonFlagArgs)
}

// warnIfArgContainsBuildArg will issue a warning if a flag is incorrectly prefixed with build-arg.
// TODO this check should be replaced with a warning if an arg was given but never used.
func (a *Build) warnIfArgContainsBuildArg(flagArgs []string) {
	for _, flag := range flagArgs {
		if strings.HasPrefix(flag, "build-arg=") || strings.HasPrefix(flag, "buildarg=") {
			a.cli.Console().Warnf("Found a flag named %q; flags after the build target should be specified as --KEY=VAL\n", flag)
		}
	}
}

func (a *Build) gitLogLevel() buildkitgitutil.GitLogLevel {
	if a.cli.Flags().Debug {
		return buildkitgitutil.GitLogLevelTrace
	}
	if a.cli.Flags().Verbose {
		return buildkitgitutil.GitLogLevelDebug
	}
	return buildkitgitutil.GitLogLevelDefault
}

func (a *Build) parseTarget(cliCtx *cli.Context, nonFlagArgs []string) (domain.Target, domain.Artifact, string, error) {
	var (
		target   domain.Target
		artifact domain.Artifact
		destPath = "./"
	)

	switch {
	case a.cli.Flags().ImageMode:
		if len(nonFlagArgs) == 0 {
			_ = cli.ShowAppHelp(cliCtx)
			return target, artifact, "", params.Errorf(
				"no image reference provided. Try %s --image +<target-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) != 1 {
			_ = cli.ShowAppHelp(cliCtx)
			return target, artifact, "", params.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		targetName := nonFlagArgs[0]
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return target, artifact, "", params.Wrapf(err, "invalid target name %s", targetName)
		}
	case a.cli.Flags().ArtifactMode:
		if len(nonFlagArgs) == 0 {
			_ = cli.ShowAppHelp(cliCtx)
			return target, artifact, "", params.Errorf(
				"no artifact reference provided. Try %s --artifact +<target-name>/<artifact-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) > 2 {
			_ = cli.ShowAppHelp(cliCtx)
			return target, artifact, "", params.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		artifactName := nonFlagArgs[0]
		if len(nonFlagArgs) == 2 {
			destPath = nonFlagArgs[1]
		}
		var err error
		artifact, err = domain.ParseArtifact(artifactName)
		if err != nil {
			return target, artifact, "", params.Wrapf(err, "invalid artifact name %s", artifactName)
		}
		target = artifact.Target
	default:
		if len(nonFlagArgs) == 0 {
			_ = cli.ShowAppHelp(cliCtx)
			return target, artifact, "", params.Errorf(
				"no target reference provided. Try %s +<target-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) != 1 {
			_ = cli.ShowAppHelp(cliCtx)
			return target, artifact, "", params.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		targetName := nonFlagArgs[0]
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return target, artifact, "", params.Errorf("invalid target %s", targetName)
		}
	}

	return target, artifact, destPath, nil
}

func (a *Build) ActionBuildImp(cliCtx *cli.Context, flagArgs, nonFlagArgs []string) error {

	target, artifact, destPath, err := a.parseTarget(cliCtx, nonFlagArgs)
	if err != nil {
		return err
	}

	a.cli.SetAnaMetaTarget(target)

	var (
		gitCommitAuthor string
		gitConfigEmail  string
	)
	if !target.IsRemote() {
		meta, _ := gitutil.Metadata(cliCtx.Context, target.GetLocalPath(), a.cli.Flags().GitBranchOverride)
		if meta != nil {
			// Git commit detection here is best effort
			gitCommitAuthor = meta.Author
		}
		if email, err := gitutil.ConfigEmail(cliCtx.Context); err == nil {
			gitConfigEmail = email
		}
	}
	cleanCollection := cleanup.NewCollection()
	defer cleanCollection.Close()

	cloudClient, err := helper.NewCloudClient(a.cli, cloud.WithLogstreamGRPCAddressOverride(a.cli.Flags().LogstreamAddressOverride))
	if err != nil {
		return err
	}

	// Determine if Logstream is enabled and create log sharing link in either case.
	logstreamURL, doLogstreamUpload, printLinkFn := a.logShareLink(cliCtx.Context, cloudClient, target, cleanCollection)

	a.cli.Console().PrintPhaseHeader(builder.PhaseInit, false, "")
	a.warnIfArgContainsBuildArg(flagArgs)

	showUnexpectedEnvWarnings := true
	dotEnvMap, err := godotenv.Read(a.cli.Flags().EnvFile)
	if err != nil {
		// ignore ErrNotExist when using default .env file
		if cliCtx.IsSet(flag.EnvFileFlag) || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "read %s", a.cli.Flags().EnvFile)
		}
	}
	argMap, err := godotenv.Read(a.cli.Flags().ArgFile)
	if err == nil {
		showUnexpectedEnvWarnings = false
	} else {
		// ignore ErrNotExist when using default .env file
		if cliCtx.IsSet(flag.ArgFileFlag) || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "read %s", a.cli.Flags().ArgFile)
		}
	}
	secretsFileMap, err := godotenv.Read(a.cli.Flags().SecretFile)
	if err == nil {
		showUnexpectedEnvWarnings = false
	} else {
		// ignore ErrNotExist when using default .env file
		if cliCtx.IsSet(flag.SecretFileFlag) || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "read %s", a.cli.Flags().SecretFile)
		}
	}

	if showUnexpectedEnvWarnings {
		validEnvNames := cliutil.GetValidEnvNames(a.cli.App())
		for k := range dotEnvMap {
			if _, found := validEnvNames[k]; !found {
				a.cli.Console().Warnf("unexpected env \"%s\": as of v0.7.0, --build-arg values must be defined in .arg (and --secret values in .secret)", k)
			}
		}
	}

	secretsMap, err := common.ProcessSecrets(a.secrets.Value(), a.secretFiles.Value(), secretsFileMap, a.cli.Flags().SecretFile)
	if err != nil {
		return err
	}
	for secretKey := range secretsMap {
		if !ast.IsValidEnvVarName(secretKey) {
			// TODO If the year is 2024 or later, please move this check into processSecrets, and turn it into an error; see https://github.com/earthly/earthly/issues/2883
			a.cli.Console().Warnf("Deprecation: secret key %q does not follow the recommended naming convention (a letter followed by alphanumeric characters or underscores); this will become an error in a future version of earthly.", secretKey)
		}
	}

	overridingVars, err := common.CombineVariables(argMap, flagArgs, a.buildArgs.Value())
	if err != nil {
		return err
	}

	skipDB, err := bk.NewBuildkitSkipper(a.cli.Flags().LocalSkipDB, cloudClient)
	if err != nil {
		a.cli.Console().WithPrefix(autoSkipPrefix).Warnf("Failed to initialize auto-skip database: %v", err)
	}

	addHashFn, doSkip, err := a.initAutoSkip(cliCtx.Context, skipDB, target, overridingVars)
	if err != nil {
		a.cli.Console().PrintFailure("auto-skip")
		return err
	}
	if doSkip {
		return nil
	}

	// Output log sharing link after build. Invoked after auto-skip is checked (above).
	a.cli.AddDeferredFunc(printLinkFn)

	err = a.cli.InitFrontend(cliCtx)
	if err != nil {
		return errors.Wrapf(err, "could not init frontend")
	}

	cleanupTLS, err := a.cli.ConfigureSatellite(cliCtx, cloudClient, gitCommitAuthor, gitConfigEmail)
	if err != nil {
		return errors.Wrapf(err, "could not configure satellite")
	}
	defer cleanupTLS()

	// Collect info to help with printing a richer message in the beginning of the build or on failure to reserve satellite due to missing build minutes.
	if err = a.cli.CollectBillingInfo(cliCtx.Context, cloudClient, a.cli.OrgName()); err != nil {
		a.cli.Console().DebugPrintf("failed to get billing plan info, error is %v\n", err)
	}

	// After configuring frontend and satellites, buildkit address should not be empty.
	// It should be set to a local container, remote address, or satellite address at this point.
	if a.cli.Flags().BuildkitdSettings.BuildkitAddress == "" {
		return errors.New("could not determine buildkit address - is Docker or Podman running?")
	}

	bkClient, err := buildkitd.NewClient(
		cliCtx.Context,
		a.cli.Console(),
		a.cli.Flags().BuildkitdImage,
		a.cli.Flags().ContainerName,
		a.cli.Flags().InstallationName,
		a.cli.Flags().ContainerFrontend,
		a.cli.Version(),
		a.cli.Flags().BuildkitdSettings,
	)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()

	platr, err := a.platformResolver(cliCtx.Context, bkClient, target)
	if err != nil {
		return err
	}

	runnerName, isLocal, err := a.runnerName(cliCtx.Context)
	if err != nil {
		return err
	}

	a.cli.SetAnaMetaIsRemoteBK(!isLocal)

	localhostProvider, err := localhostprovider.NewLocalhostProvider()
	if err != nil {
		return errors.Wrap(err, "failed to create localhostprovider")
	}

	cacheLocalDir, err := os.MkdirTemp("", "earthly-cache")
	if err != nil {
		return errors.Wrap(err, "make temp dir for cache")
	}
	defer os.RemoveAll(cacheLocalDir)
	defaultLocalDirs := make(map[string]string)
	defaultLocalDirs["earthly-cache"] = cacheLocalDir
	buildContextProvider := provider.NewBuildContextProvider(a.cli.Console())
	buildContextProvider.AddDirs(defaultLocalDirs)

	internalSecretStore := secretprovider.NewMutableMapStore(nil)
	customSecretProviderCmd, err := secretprovider.NewSecretProviderCmd(a.cli.Cfg().Global.SecretProvider)
	if err != nil {
		return errors.Wrap(err, "NewSecretProviderCmd")
	}
	secretProvider := secretprovider.New(
		internalSecretStore,
		secretprovider.NewMapStore(secretsMap),
		customSecretProviderCmd,
		secretprovider.NewCloudStore(cloudClient),
	)

	attachables := []session.Attachable{
		secretProvider,
		buildContextProvider,
		localhostProvider,
	}

	cfg := config.LoadDefaultConfigFile(os.Stderr)
	cloudStoredAuthProvider := cloudauth.NewProvider(cfg, cloudClient, a.cli.Console())

	var authChildren []authprovider.Child
	if _, _, _, err := cloudClient.WhoAmI(cliCtx.Context); err == nil {
		// only add cloud-based auth provider when logged in
		authChildren = append(authChildren, cloudStoredAuthProvider.(auth.AuthServer))
	}

	switch a.cli.Flags().ContainerFrontend.Config().Setting {
	case containerutil.FrontendPodman, containerutil.FrontendPodmanShell:
		authChildren = append(authChildren, authprovider.NewPodman(os.Stderr).(auth.AuthServer))

	default:
		// includes containerutil.FrontendDocker, containerutil.FrontendDockerShell:
		authChildren = append(authChildren, dockerauthprovider.NewDockerAuthProvider(cfg, nil).(auth.AuthServer))
	}

	authProvider := authprovider.New(a.cli.Console(), authChildren)
	attachables = append(attachables, authProvider)

	gitLookup := buildcontext.NewGitLookup(a.cli.Console(), a.cli.Flags().SSHAuthSock)
	err = a.updateGitLookupConfig(gitLookup)
	if err != nil {
		return err
	}

	if a.cli.Flags().SSHAuthSock != "" {
		ssh, err := sshprovider.NewSSHAgentProvider([]sshprovider.AgentConfig{{
			Paths: []string{a.cli.Flags().SSHAuthSock},
		}})
		if err != nil {
			return errors.Wrap(err, "ssh agent provider")
		}
		attachables = append(attachables, ssh)
	}

	localArtifactWhiteList := gatewaycrafter.NewLocalArtifactWhiteList()

	socketProvider, err := socketprovider.NewSocketProvider(map[string]socketprovider.SocketAcceptCb{
		"earthly_save_file": getTryCatchSaveFileHandler(localArtifactWhiteList),
		"earthly_interactive": func(ctx context.Context, conn io.ReadWriteCloser) error {
			if !termutil.IsTTY() {
				return fmt.Errorf("interactive mode unavailable due to terminal not being tty")
			}
			debugTermConsole := a.cli.Console().WithPrefix("internal-term")
			err := terminal.ConnectTerm(cliCtx.Context, conn, debugTermConsole)
			if err != nil {
				return errors.Wrap(err, "interactive terminal")
			}
			return nil
		},
	})
	if err != nil {
		return errors.Wrap(err, "ssh agent provider")
	}
	attachables = append(attachables, socketProvider)

	var enttlmnts []entitlements.Entitlement
	if a.cli.Flags().AllowPrivileged {
		enttlmnts = append(enttlmnts, entitlements.EntitlementSecurityInsecure)
	}

	imageResolveMode := llb.ResolveModePreferLocal
	if a.cli.Flags().Pull {
		imageResolveMode = llb.ResolveModeForcePull
	}

	cacheImports := make([]string, 0)
	var cacheImportImageName string
	if a.cli.Flags().RemoteCache != "" {
		cacheImportImageName, _, err = a.parseImageNameAndAttrs(a.cli.Flags().RemoteCache)
		cacheImports = append(cacheImports, cacheImportImageName)
	}
	if err != nil {
		return errors.Wrapf(err, "parse remote cache error: %s", a.cli.Flags().RemoteCache)
	}
	if len(a.cacheFrom.Value()) > 0 {
		cacheImports = append(cacheImports, a.cacheFrom.Value()...)
	}
	var cacheExport string
	var maxCacheExport string
	var cacheExportAttrs map[string]string
	var maxCacheExportAttrs map[string]string
	if a.cli.Flags().RemoteCache != "" && a.cli.Flags().Push {
		if a.cli.Flags().MaxRemoteCache {
			maxCacheExport, maxCacheExportAttrs, err = a.parseImageNameAndAttrs(a.cli.Flags().RemoteCache)
		} else {
			cacheExport, cacheExportAttrs, err = a.parseImageNameAndAttrs(a.cli.Flags().RemoteCache)
		}
	}
	if err != nil {
		return errors.Wrapf(err, "parse remote cache error: %s", a.cli.Flags().RemoteCache)
	}

	if a.cli.Cfg().Global.ConversionParallelism <= 0 {
		return fmt.Errorf("configuration error: \"conversion_parallelism\" must be larger than zero")
	}
	parallelism := semutil.NewWeighted(int64(a.cli.Cfg().Global.ConversionParallelism))

	localRegistryAddr := ""
	if isLocal && a.cli.Flags().LocalRegistryHost != "" {
		u, err := url.Parse(a.cli.Flags().LocalRegistryHost)
		if err != nil {
			return errors.Wrapf(err, "parse local registry host %s", a.cli.Flags().LocalRegistryHost)
		}
		localRegistryAddr = u.Host
	}

	var logbusSM *solvermon.SolverMonitor
	if a.cli.Flags().Logstream {
		logbusSM = a.cli.LogbusSetup().SolverMonitor
	} else if a.cli.Flags().DisplayExecStats {
		return fmt.Errorf("the --exec-stats feature is only available when --logstream is enabled")
	}

	builderOpts := builder.Opt{
		BkClient:                              bkClient,
		LogBusSolverMonitor:                   logbusSM,
		UseLogstream:                          a.cli.Flags().Logstream,
		Console:                               a.cli.Console(),
		Verbose:                               a.cli.Flags().Verbose,
		Attachables:                           attachables,
		Enttlmnts:                             enttlmnts,
		NoCache:                               a.cli.Flags().NoCache,
		CacheImports:                          states.NewCacheImports(cacheImports),
		CacheExport:                           cacheExport,
		CacheExportAttributes:                 cacheExportAttrs,
		MaxCacheExport:                        maxCacheExport,
		MaxCacheExportAttributes:              maxCacheExportAttrs,
		UseInlineCache:                        a.cli.Flags().UseInlineCache,
		SaveInlineCache:                       a.cli.Flags().SaveInlineCache,
		ImageResolveMode:                      imageResolveMode,
		CleanCollection:                       cleanCollection,
		OverridingVars:                        overridingVars,
		BuildContextProvider:                  buildContextProvider,
		GitLookup:                             gitLookup,
		GitBranchOverride:                     a.cli.Flags().GitBranchOverride,
		UseFakeDep:                            !a.cli.Flags().NoFakeDep,
		Strict:                                a.cli.Flags().Strict,
		DisableNoOutputUpdates:                a.cli.Flags().InteractiveDebugging,
		ParallelConversion:                    (a.cli.Cfg().Global.ConversionParallelism != 0),
		Parallelism:                           parallelism,
		LocalRegistryAddr:                     localRegistryAddr,
		DarwinProxyImage:                      a.cli.Cfg().Global.DarwinProxyImage,
		DarwinProxyWait:                       a.cli.Cfg().Global.DarwinProxyWait,
		FeatureFlagOverrides:                  a.cli.Flags().FeatureFlagOverrides,
		ContainerFrontend:                     a.cli.Flags().ContainerFrontend,
		InternalSecretStore:                   internalSecretStore,
		InteractiveDebugging:                  a.cli.Flags().InteractiveDebugging,
		InteractiveDebuggingDebugLevelLogging: a.cli.Flags().Debug,
		GitImage:                              a.cli.Cfg().Global.GitImage,
		GitLFSInclude:                         a.cli.Flags().GitLFSPullInclude,
		GitLogLevel:                           a.gitLogLevel(),
		DisableRemoteRegistryProxy:            a.cli.Flags().DisableRemoteRegistryProxy,
		BuildkitSkipper:                       skipDB,
		NoAutoSkip:                            a.cli.Flags().NoAutoSkip,
	}

	b, err := builder.NewBuilder(cliCtx.Context, builderOpts)
	if err != nil {
		return errors.Wrap(err, "new builder")
	}

	a.cli.Console().PrintPhaseFooter(builder.PhaseInit, false, "")

	builtinArgs := variables.DefaultArgs{
		EarthlyVersion:  a.cli.Version(),
		EarthlyBuildSha: a.cli.GitSHA(),
	}
	buildOpts := builder.BuildOpt{
		PrintPhases:                true,
		Push:                       a.cli.Flags().Push,
		CI:                         a.cli.Flags().CI,
		EarthlyCIRunner:            a.cli.Flags().EarthlyCIRunner,
		NoOutput:                   a.cli.Flags().NoOutput,
		OnlyFinalTargetImages:      a.cli.Flags().ImageMode,
		PlatformResolver:           platr,
		EnableGatewayClientLogging: a.cli.Flags().Debug,
		BuiltinArgs:                builtinArgs,
		LocalArtifactWhiteList:     localArtifactWhiteList,
		Logbus:                     a.cli.Logbus(),
		Runner:                     runnerName,

		// feature-flip the removal of builder.go code
		// once VERSION 0.7 is released AND support for 0.6 is dropped,
		// we can remove this flag along with code from builder.go.
		GlobalWaitBlockFtr: a.cli.Flags().GlobalWaitEnd,

		// explicitly set this to true at the top level (without granting the entitlements.EntitlementSecurityInsecure buildkit option),
		// to differentiate between a user forgetting to run earthly -P, versus a remotely referencing an earthfile that requires privileged.
		AllowPrivileged: true,

		ProjectAdder: authProvider,
	}
	if a.cli.Flags().ArtifactMode {
		buildOpts.OnlyArtifact = &artifact
		buildOpts.OnlyArtifactDestPath = destPath
	}

	// Kick off log streaming upload only when we've passed the necessary
	// information to logbusSetup. This function will be called right at the
	// beginning of the build within earthfile2llb.
	buildOpts.MainTargetDetailsFunc = func(d earthfile2llb.TargetDetails) error {
		// Use the flag/env values for org & project if specified, but fallback to
		// the PROJECT command if provided.
		orgName := d.EarthlyOrgName
		if a.cli.OrgName() != "" {
			orgName = a.cli.OrgName()
		}
		projectName := d.EarthlyProjectName
		if a.cli.Flags().ProjectName != "" {
			projectName = a.cli.Flags().ProjectName
		}
		a.cli.Console().WithPrefix("logbus").Printf("Setting organization %q and project %q", orgName, projectName)
		analytics.AddEarthfileProject(orgName, projectName)
		if !a.cli.Flags().Logstream {
			return nil
		}
		setup := a.cli.LogbusSetup()
		setup.SetOrgAndProject(orgName, projectName)
		setup.SetGitAuthor(gitCommitAuthor, gitConfigEmail)
		_, isCI := analytics.DetectCI(a.cli.Flags().EarthlyCIRunner)
		setup.SetCI(isCI)
		if doLogstreamUpload {
			setup.StartLogStreamer(cliCtx.Context, cloudClient)
		}
		return nil
	}

	if a.cli.Flags().Logstream && doLogstreamUpload && !a.cli.LogbusSetup().LogStreamerStarted() {
		a.cli.Console().ColorPrintf(color.New(color.FgHiYellow), "Streaming logs to %s\n\n", logstreamURL)
	}

	a.maybePrintBuildMinutesInfo(cliCtx)

	_, err = b.BuildTarget(cliCtx.Context, target, buildOpts)
	if err != nil {
		return errors.Wrap(err, "build target")
	}

	if a.cli.Flags().SkipBuildkit && addHashFn != nil {
		addHashFn()
	}

	return nil
}

func getTryCatchSaveFileHandler(localArtifactWhiteList *gatewaycrafter.LocalArtifactWhiteList) func(ctx context.Context, conn io.ReadWriteCloser) error {
	return func(ctx context.Context, conn io.ReadWriteCloser) error {
		// version
		protocolVersion, _, err := debuggercommon.ReadDataPacket(conn)
		if err != nil {
			return err
		}

		switch protocolVersion {
		case 1:
			return receiveFileVersion1(ctx, conn, localArtifactWhiteList)
		case 2:
			return receiveFileVersion2(ctx, conn, localArtifactWhiteList)
		default:
			return fmt.Errorf("unexpected version %d", protocolVersion)
		}
	}
}

func (a *Build) updateGitLookupConfig(gitLookup *buildcontext.GitLookup) error {
	for k, v := range a.cli.Cfg().Git {
		if k == "github" || k == "gitlab" || k == "bitbucket" {
			a.cli.Console().Warnf("git configuration for %q found, did you mean %q?\n", k, k+".com")
		}
		pattern := v.Pattern
		if pattern == "" {
			// if empty, assume it will be of the form host.com/user/repo.git
			host := k
			if !strings.Contains(host, ".") {
				host += ".com"
			}
			pattern = host + "/[^/]+/[^/]+"
		}
		auth := v.Auth
		suffix := v.Suffix
		if suffix == "" {
			suffix = ".git"
		}
		err := gitLookup.AddMatcher(k, pattern, v.Substitute, v.User, v.Password, v.Prefix, suffix, auth, v.ServerKey, common.IfNilBoolDefault(v.StrictHostKeyChecking, true), v.Port, v.SSHCommand)
		if err != nil {
			return errors.Wrap(err, "gitlookup")
		}
	}
	return nil
}

func receiveFileVersion1(ctx context.Context, conn io.ReadWriteCloser, localArtifactWhiteList *gatewaycrafter.LocalArtifactWhiteList) error {
	// dst path
	_, dst, err := debuggercommon.ReadDataPacket(conn)
	if err != nil {
		return err
	}

	if !localArtifactWhiteList.Exists(string(dst)) {
		return fmt.Errorf("file %s does not appear in the white list", dst)
	}

	// data
	_, data, err := debuggercommon.ReadDataPacket(conn)
	if err != nil {
		return err
	}

	// EOF
	n, _, err := debuggercommon.ReadDataPacket(conn)
	if err != nil {
		return err
	}
	if n != 0 {
		return fmt.Errorf("expected EOF, but got more data")
	}

	f, err := os.Create(string(dst))
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Close()
}

func receiveFileVersion2(ctx context.Context, conn io.ReadWriteCloser, localArtifactWhiteList *gatewaycrafter.LocalArtifactWhiteList) (retErr error) {
	// dst path
	dst, err := debuggercommon.ReadUint16PrefixedData(conn)
	if err != nil {
		return err
	}

	if !localArtifactWhiteList.Exists(string(dst)) {
		return fmt.Errorf("file %s does not appear in the white list", dst)
	}
	err = os.MkdirAll(path.Dir(string(dst)), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(string(dst))
	if err != nil {
		return err
	}

	defer func() {
		if retErr != nil {
			// don't output incomplete data
			_ = f.Close()
			_ = os.Remove(string(dst))
		}
	}()

	// data
	for {
		data, err := debuggercommon.ReadUint16PrefixedData(conn)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			break
		}
		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}

	return f.Close()
}

// runnerName returns the name of the local or remote BK "runner"; which is a
// representation of what BuildKit instance is being used,
// e.g. local:<hostname>, sat:<org>/<name>, or bk:<remote-address>
func (a *Build) runnerName(ctx context.Context) (string, bool, error) {
	var runnerName string
	isLocal := containerutil.IsLocal(a.cli.Flags().BuildkitdSettings.BuildkitAddress)
	if isLocal {
		hostname, err := os.Hostname()
		if err != nil {
			a.cli.Console().Warnf("failed to get hostname: %v", err)
			hostname = "unknown"
		}
		runnerName = fmt.Sprintf("local:%s", hostname)
	} else {
		if a.cli.Flags().SatelliteName != "" {
			runnerName = fmt.Sprintf("sat:%s/%s", a.cli.OrgName(), a.cli.Flags().SatelliteName)
		} else {
			runnerName = fmt.Sprintf("bk:%s", a.cli.Flags().BuildkitdSettings.BuildkitAddress)
		}
	}
	if !isLocal && (a.cli.Flags().UseInlineCache || a.cli.Flags().SaveInlineCache) {
		a.cli.Console().Warnf("Note that inline cache (--use-inline-cache and --save-inline-cache) occasionally cause builds to get stuck at 100%% CPU on Satellites and remote Buildkit.")
		a.cli.Console().Warnf("") // newline
	}
	if isLocal && !a.cli.Flags().ContainerFrontend.IsAvailable(ctx) {
		return "", false, errors.New("Frontend is not available to perform the build. Is Docker installed and running?")
	}
	return runnerName, isLocal, nil
}

func (a *Build) platformResolver(ctx context.Context, bkClient *bkclient.Client, target domain.Target) (*platutil.Resolver, error) {
	nativePlatform, err := platutil.GetNativePlatformViaBkClient(ctx, bkClient)
	if err != nil {
		return nil, errors.Wrap(err, "get native platform via buildkit client")
	}
	if a.cli.Flags().Logstream {
		a.cli.LogbusSetup().SetDefaultPlatform(platforms.Format(nativePlatform))
	}
	platr := platutil.NewResolver(nativePlatform)
	a.cli.SetAnaMetaBKPlatform(platforms.Format(nativePlatform))
	a.cli.SetAnaMetaUserPlatform(platforms.Format(platr.LLBUser()))
	platr.AllowNativeAndUser = true
	platformsSlice := make([]platutil.Platform, 0, len(a.platformsStr.Value()))
	for _, p := range a.platformsStr.Value() {
		platform, err := platr.Parse(p)
		if err != nil {
			return nil, errors.Wrapf(err, "parse platform %s", p)
		}
		platformsSlice = append(platformsSlice, platform)
	}
	switch len(platformsSlice) {
	case 0:
	case 1:
		platr.UpdatePlatform(platformsSlice[0])
	default:
		return nil, errors.Errorf("multi-platform builds are not yet supported on the command line. You may, however, create a target with the instruction BUILD --platform ... --platform ... %s", target)
	}

	return platr, nil
}

func (a *Build) initAutoSkip(ctx context.Context, skipDB bk.BuildkitSkipper, target domain.Target, overridingVars *variables.Scope) (func(), bool, error) {

	if !a.cli.Flags().SkipBuildkit {
		return nil, false, nil
	}

	console := a.cli.Console().WithPrefix(autoSkipPrefix)

	if skipDB == nil {
		return nil, false, nil
	}

	consoleNoPrefix := a.cli.Console()

	if a.cli.Flags().NoCache {
		return nil, false, errors.New("--no-cache cannot be used with --auto-skip")
	}

	if a.cli.Flags().NoAutoSkip {
		return nil, false, errors.New("--no-auto-skip cannot be used with --auto-skip")
	}

	orgName := a.cli.Flags().OrgName

	targetHash, stats, err := inputgraph.HashTarget(ctx, inputgraph.HashOpt{
		Target:         target,
		Console:        a.cli.Console(),
		CI:             a.cli.Flags().CI,
		BuiltinArgs:    variables.DefaultArgs{EarthlyVersion: a.cli.Version(), EarthlyBuildSha: a.cli.GitSHA()},
		OverridingVars: overridingVars,
	})
	if err != nil {
		return nil, false, errors.Wrapf(err, "auto-skip is unable to calculate hash for %s", target)
	}

	console.VerbosePrintf("targets visited: %d; targets hashed: %d; target cache hits: %d", stats.TargetsVisited, stats.TargetsHashed, stats.TargetCacheHits)
	console.VerbosePrintf("hash calculation took %s", stats.Duration)

	if a.cli.Flags().LocalSkipDB == "" && orgName == "" {
		orgName, _, err = inputgraph.ParseProjectCommand(ctx, target, console)
		if err != nil {
			return nil, false, errors.New("organization not found in Earthfile, command flag or environmental variables")
		}
	}

	if !target.IsRemote() {
		meta, err := gitutil.Metadata(ctx, target.GetLocalPath(), a.cli.Flags().GitBranchOverride)
		if err != nil {
			console.VerboseWarnf("unable to detect all git metadata: %v", err.Error())
		}
		target = gitutil.ReferenceWithGitMeta(target, meta).(domain.Target)
		target.Tag = ""
	}

	targetConsole := a.cli.Console().WithPrefix(target.String())
	targetStr := targetConsole.PrefixColor().Sprintf("%s", target.StringCanonical())

	exists, err := skipDB.Exists(ctx, orgName, targetHash)
	if err != nil {
		console.Warnf("Unable to check if target %s (hash %x) has already been run: %s", targetStr, targetHash, err.Error())
		return nil, false, nil
	}

	if exists {
		console.Printf("Target %s (hash %x) has already been run. Skipping.", targetStr, targetHash)
		consoleNoPrefix.PrintSuccess()
		return nil, true, nil
	}

	addHashFn := func() {
		err := skipDB.Add(ctx, orgName, target.StringCanonical(), targetHash)
		if err != nil {
			a.cli.Console().WithPrefix(autoSkipPrefix).Warnf("failed to record %s (hash %x) as completed: %s", target.String(), target, err)
		}
	}

	return addHashFn, false, nil
}

func (a *Build) logShareLink(ctx context.Context, cloudClient *cloud.Client, target domain.Target, clean *cleanup.Collection) (string, bool, func()) {

	if a.cli.Cfg().Global.DisableLogSharing {
		return "", false, func() {}
	}

	if !cloudClient.IsLoggedIn(ctx) {
		printLinkFn := func() {
			a.cli.Console().Printf(
				"🛰️ Reuse cache between CI runs with Earthly Satellites! " +
					"2-20X faster than without cache. Generous free tier " +
					"https://cloud.earthly.dev\n")
		}
		return "", false, printLinkFn
	}

	if !a.cli.Flags().LogstreamUpload {
		// If you are logged in, then add the bundle builder code, and
		// configure cleanup and post-build messages.
		a.cli.SetConsole(a.cli.Console().WithLogBundleWriter(target.String(), clean))
		printLinkFn := func() { // Defer this to keep log upload code together
			logPath, err := a.cli.Console().WriteBundleToDisk()
			if err != nil {
				err := errors.Wrapf(err, "failed to write log to disk")
				a.cli.Console().Warnf(err.Error())
				return
			}

			id, err := cloudClient.UploadLog(ctx, logPath)
			if err != nil {
				err := errors.Wrapf(err, "failed to upload log")
				a.cli.Console().Warnf(err.Error())
				return
			}
			a.cli.Console().ColorPrintf(color.New(color.FgHiYellow), "Shareable link: %s\n", id)
		}
		return "", false, printLinkFn
	}

	logstreamURL := fmt.Sprintf("%s/builds/%s", a.cli.CIHost(), a.cli.LogbusSetup().InitialManifest.GetBuildId())

	printLinkFn := func() {
		a.cli.Console().ColorPrintf(color.New(color.FgHiYellow), "View logs at %s\n", logstreamURL)
	}

	return logstreamURL, true, printLinkFn
}

func (a *Build) maybePrintBuildMinutesInfo(cliCtx *cli.Context) {
	orgName := a.cli.OrgName()
	settings := a.cli.Flags().BuildkitdSettings
	if !a.cli.IsUsingSatellite(cliCtx) || !settings.SatelliteIsManaged {
		return
	}

	plan := billing.Plan()
	if plan.GetMaxBuildMinutes() == 0 {
		return
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Build Minutes: %d out of %d used\n", int(billing.UsedBuildTime().Minutes()), plan.GetMaxBuildMinutes()))
	if plan.GetTier() == billingpb.BillingPlan_TIER_LIMITED_FREE_TIER {
		sb.WriteString(fmt.Sprintf("Visit your organization settings to verify your account\nand get 6000 free build minutes per month: %s\n", billing.GetBillingURL(a.cli.CIHost(), orgName)))
	}
	sb.WriteRune('\n')
	if plan.GetType() == billingpb.BillingPlan_PLAN_TYPE_FREE {
		a.cli.Console().ColorPrintf(color.New(color.FgGreen), sb.String())
	} else {
		a.cli.Console().VerbosePrintf(sb.String())
	}
}

func (a *Build) actionDockerBuild(cliCtx *cli.Context) error {
	a.cli.SetCommandName("docker-build")

	flagArgs, nonFlagArgs, err := variables.ParseFlagArgsWithNonFlags(cliCtx.Args().Slice())
	if err != nil {
		return errors.Wrapf(err, "parse args %s", strings.Join(cliCtx.Args().Slice(), " "))
	}
	if len(nonFlagArgs) == 0 {
		_ = cli.ShowAppHelp(cliCtx)
		return errors.Errorf(
			"no build context path provided. Try %s docker-build <path>", cliCtx.App.Name)
	}
	if len(nonFlagArgs) != 1 {
		_ = cli.ShowAppHelp(cliCtx)
		return errors.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
	}

	buildContextPath, err := filepath.Abs(nonFlagArgs[0])
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path for build context")
	}

	tempDir, err := os.MkdirTemp("", "docker-build")
	if err != nil {
		return errors.Wrap(err, "docker-build: failed to create temporary dir for Earthfile")
	}
	defer os.RemoveAll(tempDir)

	argMap, err := godotenv.Read(a.cli.Flags().ArgFile)
	if err != nil && (cliCtx.IsSet(flag.ArgFileFlag) || !errors.Is(err, os.ErrNotExist)) {
		return errors.Wrapf(err, "read %q", a.cli.Flags().ArgFile)
	}

	buildArgs, err := common.CombineVariables(argMap, flagArgs, a.buildArgs.Value())
	if err != nil {
		return errors.Wrapf(err, "combining build args")
	}

	platforms := flagutil.SplitFlagString(a.platformsStr)
	content, err := docker2earthly.GenerateEarthfile(buildContextPath, a.cli.Flags().DockerfilePath, a.dockerTags.Value(), buildArgs.Sorted(), platforms, a.dockerTarget)
	if err != nil {
		return errors.Wrap(err, "docker-build: failed to wrap Dockerfile with an Earthfile")
	}

	earthfilePath := filepath.Join(tempDir, "Earthfile")

	out, err := os.Create(earthfilePath)
	if err != nil {
		return errors.Wrapf(err, "docker-build: failed to create Earthfile %q", earthfilePath)
	}
	defer out.Close()

	_, err = out.WriteString(content)
	if err != nil {
		return errors.Wrapf(err, "docker-build: failed to write to %q", earthfilePath)
	}

	// The following should not be set in the context of executing the build from the generated Earthfile:
	a.cli.Flags().DockerfilePath = ""
	a.cli.Flags().ImageMode = false
	a.cli.Flags().ArtifactMode = false
	a.dockerTarget = ""
	a.dockerTags = cli.StringSlice{}
	a.platformsStr = cli.StringSlice{}

	nonFlagArgs = []string{tempDir + "+build"}
	return a.ActionBuildImp(cliCtx, flagArgs, nonFlagArgs)
}

func (a *Build) parseImageNameAndAttrs(s string) (string, map[string]string, error) {
	entries := strings.Split(s, ",")
	imageName := entries[0]
	attrs := make(map[string]string)
	var err error
	for _, entry := range entries[1:] {
		pair := strings.Split(strings.TrimSpace(entry), "=")
		if len(pair) != 2 {
			return "", attrs, errors.Errorf("failed to parse export attribute: expected a key=value pair while parsing %q", entry)
		}
		attrs[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}
	return imageName, attrs, err
}
