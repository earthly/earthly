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
	"time"

	"github.com/earthly/earthly/cmd/earthly/bk"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/earthly/earthly/docker2earthly"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/cli/cli/config"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth"
	dockerauthprovider "github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/localhost/localhostprovider"
	"github.com/moby/buildkit/session/socketforward/socketprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/cloud"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/terminal"
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
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/termutil"
	"github.com/earthly/earthly/variables"
)

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
			return errors.New("unable to use --ci flag in combination with --interactive flag")
		}
	}

	if a.cli.Flags().ImageMode && a.cli.Flags().ArtifactMode {
		return errors.New("both image and artifact modes cannot be active at the same time")
	}
	if (a.cli.Flags().ImageMode && a.cli.Flags().NoOutput) || (a.cli.Flags().ArtifactMode && a.cli.Flags().NoOutput) {
		if a.cli.Flags().CI {
			a.cli.Flags().NoOutput = false
		} else {
			return errors.New("cannot use --no-output with image or artifact modes")
		}
	}
	if a.cli.Flags().InteractiveDebugging && !termutil.IsTTY() {
		return errors.New("A tty-terminal must be present in order to use the --interactive flag")
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

func (a *Build) gitLogLevel() llb.GitLogLevel {
	if a.cli.Flags().Debug {
		return llb.GitLogLevelTrace
	}
	if a.cli.Flags().Verbose {
		return llb.GitLogLevelDebug
	}
	return llb.GitLogLevelDefault
}

func (a *Build) ActionBuildImp(cliCtx *cli.Context, flagArgs, nonFlagArgs []string) error {
	var target domain.Target
	var artifact domain.Artifact
	destPath := "./"
	if a.cli.Flags().ImageMode {
		if len(nonFlagArgs) == 0 {
			_ = cli.ShowAppHelp(cliCtx)
			return errors.Errorf(
				"no image reference provided. Try %s --image +<target-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) != 1 {
			_ = cli.ShowAppHelp(cliCtx)
			return errors.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		targetName := nonFlagArgs[0]
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return errors.Wrapf(err, "parse target name %s", targetName)
		}
	} else if a.cli.Flags().ArtifactMode {
		if len(nonFlagArgs) == 0 {
			_ = cli.ShowAppHelp(cliCtx)
			return errors.Errorf(
				"no artifact reference provided. Try %s --artifact +<target-name>/<artifact-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) > 2 {
			_ = cli.ShowAppHelp(cliCtx)
			return errors.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		artifactName := nonFlagArgs[0]
		if len(nonFlagArgs) == 2 {
			destPath = nonFlagArgs[1]
		}
		var err error
		artifact, err = domain.ParseArtifact(artifactName)
		if err != nil {
			return errors.Wrapf(err, "parse artifact name %s", artifactName)
		}
		target = artifact.Target
	} else {
		if len(nonFlagArgs) == 0 {
			_ = cli.ShowAppHelp(cliCtx)
			return errors.Errorf(
				"no target reference provided. Try %s +<target-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) != 1 {
			_ = cli.ShowAppHelp(cliCtx)
			return errors.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		targetName := nonFlagArgs[0]
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return errors.Wrapf(err, "parse target name %s", targetName)
		}
	}
	a.cli.SetAnaMetaTarget(target)

	var (
		gitCommitAuthor string
		gitConfigEmail  string
	)
	if !target.IsRemote() {
		if meta, err := gitutil.Metadata(cliCtx.Context, target.GetLocalPath(), a.cli.Flags().GitBranchOverride); err == nil {
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

	// Default upload logs, unless explicitly configured
	doLogstreamUpload := false
	var logstreamURL string
	if !a.cli.Cfg().Global.DisableLogSharing {
		if cloudClient.IsLoggedIn(cliCtx.Context) {
			if a.cli.Flags().LogstreamUpload {
				doLogstreamUpload = true
				logstreamURL = fmt.Sprintf("%s/builds/%s", a.cli.CIHost(), a.cli.LogbusSetup().InitialManifest.GetBuildId())
				defer func() {
					a.cli.Console().ColorPrintf(color.New(color.FgHiYellow), "View logs at %s\n", logstreamURL)
				}()
			} else {
				// If you are logged in, then add the bundle builder code, and configure cleanup and post-build messages.
				a.cli.SetConsole(a.cli.Console().WithLogBundleWriter(target.String(), cleanCollection))

				defer func() { // Defer this to keep log upload code together
					logPath, err := a.cli.Console().WriteBundleToDisk()
					if err != nil {
						err := errors.Wrapf(err, "failed to write log to disk")
						a.cli.Console().Warnf(err.Error())
						return
					}

					id, err := cloudClient.UploadLog(cliCtx.Context, logPath)
					if err != nil {
						err := errors.Wrapf(err, "failed to upload log")
						a.cli.Console().Warnf(err.Error())
						return
					}
					a.cli.Console().ColorPrintf(color.New(color.FgHiYellow), "Shareable link: %s\n", id)
				}()
			}
		} else {
			defer func() { // Defer this to keep log upload code together
				a.cli.Console().Printf(
					"🛰️ Reuse cache between CI runs with Earthly Satellites! " +
						"2-20X faster than without cache. Generous free tier " +
						"https://cloud.earthly.dev\n")
			}()
		}
	}

	a.cli.Console().PrintPhaseHeader(builder.PhaseInit, false, "")
	a.warnIfArgContainsBuildArg(flagArgs)

	var skipDB bk.BuildkitSkipper
	var targetHash []byte
	if a.cli.Flags().SkipBuildkit {
		var orgName string
		var projectName string
		orgName, projectName, targetHash, err = inputgraph.HashTarget(cliCtx.Context, target, a.cli.Console())
		if err != nil {
			a.cli.Console().Warnf("unable to calculate hash for %s: %s", target.String(), err.Error())
		} else {
			skipDB, err = bk.NewBuildkitSkipper(a.cli.Flags().LocalSkipDB, orgName, projectName, target.GetName(), cloudClient)
			if err != nil {
				return err
			}
			exists, err := skipDB.Exists(cliCtx.Context, targetHash)
			if err != nil {
				a.cli.Console().Warnf("unable to check if target %s (hash %x) has already been run: %s", target.String(), targetHash, err.Error())
			}
			if exists {
				a.cli.Console().Printf("target %s (hash %x) has already been run; exiting", target.String(), targetHash)
				return nil
			}
		}
	}

	err = a.cli.InitFrontend(cliCtx)
	if err != nil {
		return errors.Wrapf(err, "could not init frontend")
	}

	err = a.cli.ConfigureSatellite(cliCtx, cloudClient, gitCommitAuthor, gitConfigEmail)
	if err != nil {
		return errors.Wrapf(err, "could not configure satellite")
	}

	// After configuring frontend and satellites, buildkit address should not be empty.
	// It should be set to a local container, remote address, or satellite address at this point.
	if a.cli.Flags().BuildkitdSettings.BuildkitAddress == "" {
		return errors.New("could not determine buildkit address - is Docker or Podman running?")
	}

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
			runnerName = fmt.Sprintf("sat:%s/%s", a.cli.Flags().OrgName, a.cli.Flags().SatelliteName)
		} else {
			runnerName = fmt.Sprintf("bk:%s", a.cli.Flags().BuildkitdSettings.BuildkitAddress)
		}
	}
	if !isLocal && (a.cli.Flags().UseInlineCache || a.cli.Flags().SaveInlineCache) {
		a.cli.Console().Warnf("Note that inline cache (--use-inline-cache and --save-inline-cache) occasionally cause builds to get stuck at 100%% CPU on Satellites and remote Buildkit.")
		a.cli.Console().Warnf("") // newline
	}
	if isLocal && !a.cli.Flags().ContainerFrontend.IsAvailable(cliCtx.Context) {
		return errors.New("Frontend is not available to perform the build. Is Docker installed and running?")
	}

	bkclient, err := buildkitd.NewClient(
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
	defer bkclient.Close()
	a.cli.SetAnaMetaIsRemoteBK(!isLocal)

	nativePlatform, err := platutil.GetNativePlatformViaBkClient(cliCtx.Context, bkclient)
	if err != nil {
		return errors.Wrap(err, "get native platform via buildkit client")
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
			return errors.Wrapf(err, "parse platform %s", p)
		}
		platformsSlice = append(platformsSlice, platform)
	}
	switch len(platformsSlice) {
	case 0:
	case 1:
		platr.UpdatePlatform(platformsSlice[0])
	default:
		return errors.Errorf("multi-platform builds are not yet supported on the command line. You may, however, create a target with the instruction BUILD --platform ... --platform ... %s", target)
	}

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

	secretsMap, err := common.ProcessSecrets(a.secrets.Value(), a.secretFiles.Value(), secretsFileMap)
	if err != nil {
		return err
	}
	for secretKey := range secretsMap {
		if !ast.IsValidEnvVarName(secretKey) {
			// TODO If the year is 2024 or later, please move this check into processSecrets, and turn it into an error; see https://github.com/earthly/earthly/issues/2883
			a.cli.Console().Warnf("Deprecation: secret key %q does not follow the recommended naming convention (a letter followed by alphanumeric characters or underscores); this will become an error in a future version of earthly.", secretKey)
		}
	}

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
		authChildren = append(authChildren, dockerauthprovider.NewDockerAuthProvider(cfg).(auth.AuthServer))
	}

	authProvider := authprovider.New(authChildren)
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

	overridingVars, err := common.CombineVariables(argMap, flagArgs, a.buildArgs.Value())
	if err != nil {
		return err
	}

	imageResolveMode := llb.ResolveModePreferLocal
	if a.cli.Flags().Pull {
		imageResolveMode = llb.ResolveModeForcePull
	}

	cacheImports := make([]string, 0)
	if a.cli.Flags().RemoteCache != "" {
		cacheImports = append(cacheImports, a.cli.Flags().RemoteCache)
	}
	if len(a.cacheFrom.Value()) > 0 {
		cacheImports = append(cacheImports, a.cacheFrom.Value()...)
	}
	var cacheExport string
	var maxCacheExport string
	if a.cli.Flags().RemoteCache != "" && a.cli.Flags().Push {
		if a.cli.Flags().MaxRemoteCache {
			maxCacheExport = a.cli.Flags().RemoteCache
		} else {
			cacheExport = a.cli.Flags().RemoteCache
		}
	}
	if a.cli.Cfg().Global.ConversionParallelism <= 0 {
		return fmt.Errorf("configuration error: \"conversion_parallelism\" must be larger than zero")
	}
	parallelism := semutil.NewWeighted(int64(a.cli.Cfg().Global.ConversionParallelism))
	localRegistryAddr := ""
	if isLocal && a.cli.Flags().LocalRegistryHost != "" {
		lrURL, err := url.Parse(a.cli.Flags().LocalRegistryHost)
		if err != nil {
			return errors.Wrapf(err, "parse local registry host %s", a.cli.Flags().LocalRegistryHost)
		}
		localRegistryAddr = lrURL.Host
	}
	var logbusSM *solvermon.SolverMonitor
	if a.cli.Flags().Logstream {
		logbusSM = a.cli.LogbusSetup().SolverMonitor
	}
	builderOpts := builder.Opt{
		BkClient:                              bkclient,
		LogBusSolverMonitor:                   logbusSM,
		UseLogstream:                          a.cli.Flags().Logstream,
		Console:                               a.cli.Console(),
		Verbose:                               a.cli.Flags().Verbose,
		Attachables:                           attachables,
		Enttlmnts:                             enttlmnts,
		NoCache:                               a.cli.Flags().NoCache,
		CacheImports:                          states.NewCacheImports(cacheImports),
		CacheExport:                           cacheExport,
		MaxCacheExport:                        maxCacheExport,
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
		FeatureFlagOverrides:                  a.cli.Flags().FeatureFlagOverrides,
		ContainerFrontend:                     a.cli.Flags().ContainerFrontend,
		InternalSecretStore:                   internalSecretStore,
		InteractiveDebugging:                  a.cli.Flags().InteractiveDebugging,
		InteractiveDebuggingDebugLevelLogging: a.cli.Flags().Debug,
		GitImage:                              a.cli.Cfg().Global.GitImage,
		GitLFSInclude:                         a.cli.Flags().GitLFSPullInclude,
		GitLogLevel:                           a.gitLogLevel(),
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

	// Kick off logstream upload only when we've passed the necessary
	// information to logbusSetup. This function will be called right at the
	// beginning of the build within earthfile2llb.
	buildOpts.MainTargetDetailsFunc = func(d earthfile2llb.TargetDetails) error {
		if a.cli.LogbusSetup().LogStreamerStarted() {
			// If the org & project have been provided by envs, let's verify
			// that they're correct once we've parsed them from the Earthfile.
			if a.cli.Flags().OrgName != d.EarthlyOrgName || a.cli.Flags().ProjectName != d.EarthlyProjectName {
				return fmt.Errorf("organization or project do not match PROJECT statement")
			}
			a.cli.Console().VerbosePrintf("Organization and project already set via environmental")
			return nil
		}
		a.cli.Console().VerbosePrintf("Logbus: setting organization %q and project %q at %s", d.EarthlyOrgName, d.EarthlyProjectName, time.Now().Format(time.RFC3339Nano))
		analytics.AddEarthfileProject(d.EarthlyOrgName, d.EarthlyProjectName)
		if a.cli.Flags().Logstream {
			a.cli.LogbusSetup().SetOrgAndProject(d.EarthlyOrgName, d.EarthlyProjectName)
			if doLogstreamUpload {
				a.cli.LogbusSetup().StartLogStreamer(cliCtx.Context, cloudClient)
			}
		}
		return nil
	}

	if a.cli.Flags().Logstream && doLogstreamUpload && !a.cli.LogbusSetup().LogStreamerStarted() {
		a.cli.Console().ColorPrintf(color.New(color.FgHiYellow), "Streaming logs to %s\n\n", logstreamURL)
	}

	_, err = b.BuildTarget(cliCtx.Context, target, buildOpts)
	if err != nil {
		return errors.Wrap(err, "build target")
	}

	if a.cli.Flags().SkipBuildkit && targetHash != nil {
		err := skipDB.Add(cliCtx.Context, targetHash)
		if err != nil {
			a.cli.Console().Warnf("failed to record %s (hash %x) as completed: %s", target.String(), target, err)
		}
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
		err := gitLookup.AddMatcher(k, pattern, v.Substitute, v.User, v.Password, v.Prefix, suffix, auth, v.ServerKey, common.IfNilBoolDefault(v.StrictHostKeyChecking, true), v.Port)
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
