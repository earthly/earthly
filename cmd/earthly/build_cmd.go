package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/cli/cli/config"
	"github.com/joho/godotenv"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/localhost/localhostprovider"
	"github.com/moby/buildkit/session/socketforward/socketprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/terminal"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/logbus/solvermon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/gatewaycrafter"
	"github.com/earthly/earthly/util/gitutil"
	"github.com/earthly/earthly/util/llbutil/secretprovider"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/termutil"
	"github.com/earthly/earthly/variables"
)

func (app *earthlyApp) actionBuild(cliCtx *cli.Context) error {
	app.commandName = "build"

	if app.ci {
		app.useInlineCache = true
		app.noOutput = !app.output && !app.artifactMode && !app.imageMode
		app.strict = true
		if app.remoteCache == "" && app.push {
			app.saveInlineCache = true
		}

		if app.interactiveDebugging {
			return errors.New("unable to use --ci flag in combination with --interactive flag")
		}
	}

	if app.imageMode && app.artifactMode {
		return errors.New("both image and artifact modes cannot be active at the same time")
	}
	if (app.imageMode && app.noOutput) || (app.artifactMode && app.noOutput) {
		if app.ci {
			app.noOutput = false
		} else {
			return errors.New("cannot use --no-output with image or artifact modes")
		}
	}
	if app.interactiveDebugging && !termutil.IsTTY() {
		return errors.New("A tty-terminal must be present in order to use the --interactive flag")
	}

	flagArgs, nonFlagArgs, err := variables.ParseFlagArgsWithNonFlags(cliCtx.Args().Slice())
	if err != nil {
		return errors.Wrapf(err, "parse args %s", strings.Join(cliCtx.Args().Slice(), " "))
	}

	return app.actionBuildImp(cliCtx, flagArgs, nonFlagArgs)
}

// warnIfArgContainsBuildArg will issue a warning if a flag is incorrectly prefixed with build-arg.
// TODO this check should be replaced with a warning if an arg was given but never used.
func (app *earthlyApp) warnIfArgContainsBuildArg(flagArgs []string) {
	for _, flag := range flagArgs {
		if strings.HasPrefix(flag, "build-arg=") || strings.HasPrefix(flag, "buildarg=") {
			app.console.Warnf("Found a flag named %q; flags after the build target should be specified as --KEY=VAL\n", flag)
		}
	}
}

func (app *earthlyApp) combineVariables(dotEnvMap map[string]string, flagArgs []string) (*variables.Scope, error) {
	dotEnvVars := variables.NewScope()
	for k, v := range dotEnvMap {
		dotEnvVars.AddInactive(k, v)
	}
	buildArgs := append([]string{}, app.buildArgs.Value()...)
	buildArgs = append(buildArgs, flagArgs...)
	overridingVars, err := variables.ParseCommandLineArgs(buildArgs)
	if err != nil {
		return nil, errors.Wrap(err, "parse build args")
	}
	return variables.CombineScopes(overridingVars, dotEnvVars), nil
}

func (app *earthlyApp) actionBuildImp(cliCtx *cli.Context, flagArgs, nonFlagArgs []string) error {
	var target domain.Target
	var artifact domain.Artifact
	destPath := "./"
	if app.imageMode {
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
	} else if app.artifactMode {
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

	var (
		gitCommitAuthor string
		gitConfigEmail  string
	)
	if !target.IsRemote() {
		if meta, err := gitutil.Metadata(cliCtx.Context, target.GetLocalPath()); err == nil {
			// Git commit detection here is best effort
			gitCommitAuthor = meta.Author
		}
		if email, err := gitutil.ConfigEmail(cliCtx.Context); err == nil {
			gitConfigEmail = email
		}
	}

	cleanCollection := cleanup.NewCollection()
	defer cleanCollection.Close()

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	// Default upload logs, unless explicitly configured
	doLogstreamUpload := false
	var logstreamURL string
	if !app.cfg.Global.DisableLogSharing {
		if cloudClient.IsLoggedIn(cliCtx.Context) {
			if app.logstreamUpload {
				doLogstreamUpload = true
				logstreamURL := fmt.Sprintf("%s/builds/%s", app.getCIHost(), app.logbusSetup.InitialManifest.GetBuildId())
				defer func() {
					app.console.Printf("View logs at %s\n", logstreamURL)
				}()
			} else {
				// If you are logged in, then add the bundle builder code, and configure cleanup and post-build messages.
				app.console = app.console.WithLogBundleWriter(target.String(), cleanCollection)

				defer func() { // Defer this to keep log upload code together
					logPath, err := app.console.WriteBundleToDisk()
					if err != nil {
						err := errors.Wrapf(err, "failed to write log to disk")
						app.console.Warnf(err.Error())
						return
					}

					id, err := cloudClient.UploadLog(cliCtx.Context, logPath)
					if err != nil {
						err := errors.Wrapf(err, "failed to upload log")
						app.console.Warnf(err.Error())
						return
					}
					app.console.Printf("Shareable link: %s\n", id)
				}()
			}
		} else {
			defer func() { // Defer this to keep log upload code together
				app.console.Printf("Share your logs with an Earthly account (experimental)! Register for one at https://ci.earthly.dev.")
			}()
		}
	}

	app.console.PrintPhaseHeader(builder.PhaseInit, false, "")
	app.warnIfArgContainsBuildArg(flagArgs)

	err = app.initFrontend(cliCtx)
	if err != nil {
		return err
	}

	err = app.configureSatellite(cliCtx, cloudClient, gitCommitAuthor, gitConfigEmail)
	if err != nil {
		return errors.Wrapf(err, "could not construct new buildkit client")
	}

	var runnerName string
	isLocal := containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress)
	if isLocal {
		hostname, err := os.Hostname()
		if err != nil {
			app.console.Warnf("failed to get hostname: %v", err)
			hostname = "unknown"
		}
		runnerName = fmt.Sprintf("local:%s", hostname)
	} else {
		if app.satelliteName != "" {
			runnerName = fmt.Sprintf("sat:%s/%s", app.orgName, app.satelliteName)
		} else {
			runnerName = fmt.Sprintf("bk:%s", app.buildkitdSettings.BuildkitAddress)
		}
	}
	if !isLocal && app.ci {
		app.console.Printf("Please note that --use-inline-cache and --save-inline-cache are currently disabled when using --ci on Satellites or remote Buildkit.")
		app.console.Printf("") // newline
		app.useInlineCache = false
		app.saveInlineCache = false
	}
	if isLocal && !app.containerFrontend.IsAvailable(cliCtx.Context) {
		return errors.New("Frontend is not available to perform the build. Is Docker installed and running?")
	}

	bkClient, err := buildkitd.NewClient(cliCtx.Context, app.console, app.buildkitdImage, app.containerName, app.installationName, app.containerFrontend, Version, app.buildkitdSettings)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	app.analyticsMetadata.isRemoteBuildkit = !isLocal

	nativePlatform, err := platutil.GetNativePlatformViaBkClient(cliCtx.Context, bkClient)
	if err != nil {
		return errors.Wrap(err, "get native platform via buildkit client")
	}
	if app.logstream {
		app.logbusSetup.SetDefaultPlatform(platforms.Format(nativePlatform))
	}
	platr := platutil.NewResolver(nativePlatform)
	app.analyticsMetadata.buildkitPlatform = platforms.Format(nativePlatform)
	app.analyticsMetadata.userPlatform = platforms.Format(platr.LLBUser())
	platr.AllowNativeAndUser = true
	platformsSlice := make([]platutil.Platform, 0, len(app.platformsStr.Value()))
	for _, p := range app.platformsStr.Value() {
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

	dotEnvMap, err := godotenv.Read(app.envFile)
	if err != nil {
		// ignore ErrNotExist when using default .env file
		if app.envFile != defaultEnvFile || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "read %s", app.envFile)
		}
	}

	secretsMap, err := processSecrets(app.secrets.Value(), app.secretFiles.Value(), dotEnvMap)
	if err != nil {
		return err
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
	buildContextProvider := provider.NewBuildContextProvider(app.console)
	buildContextProvider.AddDirs(defaultLocalDirs)

	internalSecretStore := secretprovider.NewMutableMapStore(nil)
	customSecretProviderCmd, err := secretprovider.NewSecretProviderCmd(app.cfg.Global.SecretProvider)
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

	switch app.containerFrontend.Config().Setting {
	case containerutil.FrontendDocker, containerutil.FrontendDockerShell:
		cfg := config.LoadDefaultConfigFile(os.Stderr)
		attachables = append(attachables, authprovider.NewDockerAuthProvider(cfg))

	case containerutil.FrontendPodman, containerutil.FrontendPodmanShell:
		attachables = append(attachables, authprovider.NewPodmanAuthProvider(os.Stderr))

	default:
		// Old default behavior
		cfg := config.LoadDefaultConfigFile(os.Stderr)
		attachables = append(attachables, authprovider.NewDockerAuthProvider(cfg))
	}

	gitLookup := buildcontext.NewGitLookup(app.console, app.sshAuthSock)
	err = app.updateGitLookupConfig(gitLookup)
	if err != nil {
		return err
	}

	if app.sshAuthSock != "" {
		ssh, err := sshprovider.NewSSHAgentProvider([]sshprovider.AgentConfig{{
			Paths: []string{app.sshAuthSock},
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
			debugTermConsole := app.console.WithPrefix("internal-term")
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
	if app.allowPrivileged {
		enttlmnts = append(enttlmnts, entitlements.EntitlementSecurityInsecure)
	}

	overridingVars, err := app.combineVariables(dotEnvMap, flagArgs)
	if err != nil {
		return err
	}

	imageResolveMode := llb.ResolveModePreferLocal
	if app.pull {
		imageResolveMode = llb.ResolveModeForcePull
	}

	cacheImports := make([]string, 0)
	if app.remoteCache != "" {
		cacheImports = append(cacheImports, app.remoteCache)
	}
	if len(app.cacheFrom.Value()) > 0 {
		cacheImports = append(cacheImports, app.cacheFrom.Value()...)
	}
	var cacheExport string
	var maxCacheExport string
	if app.remoteCache != "" && app.push {
		if app.maxRemoteCache {
			maxCacheExport = app.remoteCache
		} else {
			cacheExport = app.remoteCache
		}
	}
	var parallelism semutil.Semaphore
	if app.cfg.Global.ConversionParallelism != 0 {
		parallelism = semutil.NewWeighted(int64(app.cfg.Global.ConversionParallelism))
	}
	localRegistryAddr := ""
	if isLocal && app.localRegistryHost != "" {
		lrURL, err := url.Parse(app.localRegistryHost)
		if err != nil {
			return errors.Wrapf(err, "parse local registry host %s", app.localRegistryHost)
		}
		localRegistryAddr = lrURL.Host
	}
	var logbusSM *solvermon.SolverMonitor
	if app.logstream {
		logbusSM = app.logbusSetup.SolverMonitor
	}
	builderOpts := builder.Opt{
		BkClient:                              bkClient,
		LogBusSolverMonitor:                   logbusSM,
		UseLogstream:                          app.logstream,
		Console:                               app.console,
		Verbose:                               app.verbose,
		Attachables:                           attachables,
		Enttlmnts:                             enttlmnts,
		NoCache:                               app.noCache,
		CacheImports:                          states.NewCacheImports(cacheImports),
		CacheExport:                           cacheExport,
		MaxCacheExport:                        maxCacheExport,
		UseInlineCache:                        app.useInlineCache,
		SaveInlineCache:                       app.saveInlineCache,
		ImageResolveMode:                      imageResolveMode,
		CleanCollection:                       cleanCollection,
		OverridingVars:                        overridingVars,
		BuildContextProvider:                  buildContextProvider,
		GitLookup:                             gitLookup,
		UseFakeDep:                            !app.noFakeDep,
		Strict:                                app.strict,
		DisableNoOutputUpdates:                app.interactiveDebugging,
		ParallelConversion:                    (app.cfg.Global.ConversionParallelism != 0),
		Parallelism:                           parallelism,
		LocalRegistryAddr:                     localRegistryAddr,
		FeatureFlagOverrides:                  app.featureFlagOverrides,
		ContainerFrontend:                     app.containerFrontend,
		InternalSecretStore:                   internalSecretStore,
		InteractiveDebugging:                  app.interactiveDebugging,
		InteractiveDebuggingDebugLevelLogging: app.debug,
	}
	b, err := builder.NewBuilder(cliCtx.Context, builderOpts)
	if err != nil {
		return errors.Wrap(err, "new builder")
	}

	app.console.PrintPhaseFooter(builder.PhaseInit, false, "")

	builtinArgs := variables.DefaultArgs{
		EarthlyVersion:  Version,
		EarthlyBuildSha: GitSha,
	}
	buildOpts := builder.BuildOpt{
		PrintPhases:                true,
		Push:                       app.push,
		CI:                         app.ci,
		NoOutput:                   app.noOutput,
		OnlyFinalTargetImages:      app.imageMode,
		PlatformResolver:           platr,
		EnableGatewayClientLogging: app.debug,
		BuiltinArgs:                builtinArgs,
		LocalArtifactWhiteList:     localArtifactWhiteList,
		Logbus:                     app.logbus,
		MainTargetDetailsFuture:    make(chan earthfile2llb.TargetDetails, 1),
		Runner:                     runnerName,

		// feature-flip the removal of builder.go code
		// once VERSION 0.7 is released AND support for 0.6 is dropped,
		// we can remove this flag along with code from builder.go.
		GlobalWaitBlockFtr: app.globalWaitEnd,

		// explicitly set this to true at the top level (without granting the entitlements.EntitlementSecurityInsecure buildkit option),
		// to differentiate between a user forgetting to run earthly -P, versus a remotely referencing an earthfile that requires privileged.
		AllowPrivileged: true,
	}
	if app.artifactMode {
		buildOpts.OnlyArtifact = &artifact
		buildOpts.OnlyArtifactDestPath = destPath
	}
	// Kick off logstream upload only when we've passed the necessary information to logbusSetup.
	// This information is passed back right at the beginning of the build within earthfile2llb.
	go func() {
		select {
		case <-cliCtx.Context.Done():
			return
		case details := <-buildOpts.MainTargetDetailsFuture:
			if app.logstream {
				app.logbusSetup.SetOrgAndProject(details.EarthlyOrgName, details.EarthlyProjectName)
				if doLogstreamUpload {
					app.logbusSetup.StartLogStreamer(cliCtx.Context, cloudClient)
					app.console.Printf("Streaming logs to %s\n", logstreamURL)
				}
			}
		}
	}()
	_, err = b.BuildTarget(cliCtx.Context, target, buildOpts)
	if err != nil {
		return errors.Wrap(err, "build target")
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

func receiveFileVersion2(ctx context.Context, conn io.ReadWriteCloser, localArtifactWhiteList *gatewaycrafter.LocalArtifactWhiteList) error {
	// dst path
	dst, err := debuggercommon.ReadUint16PrefixedData(conn)
	if err != nil {
		return err
	}

	if !localArtifactWhiteList.Exists(string(dst)) {
		return fmt.Errorf("file %s does not appear in the white list", dst)
	}

	f, err := os.Create(string(dst))
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
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
