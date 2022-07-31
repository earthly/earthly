package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/cli/cli/config"
	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/cloud"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/terminal"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/llbutil/secretprovider"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/termutil"
	"github.com/earthly/earthly/variables"
	"github.com/joho/godotenv"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/localhost/localhostprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
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
	if app.interactiveDebugging {
		if !termutil.IsTTY() {
			return errors.New("A tty-terminal must be present in order to the --interactive flag")
		}
		if !containerutil.IsLocal(app.buildkitHost) {
			return errors.New("the --interactive flag is not currently supported with non-local buildkit servers")
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
			cli.ShowAppHelp(cliCtx)
			return errors.Errorf(
				"no image reference provided. Try %s --image +<target-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) != 1 {
			cli.ShowAppHelp(cliCtx)
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
			cli.ShowAppHelp(cliCtx)
			return errors.Errorf(
				"no artifact reference provided. Try %s --artifact +<target-name>/<artifact-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) > 2 {
			cli.ShowAppHelp(cliCtx)
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
			cli.ShowAppHelp(cliCtx)
			return errors.Errorf(
				"no target reference provided. Try %s +<target-name>", cliCtx.App.Name)
		} else if len(nonFlagArgs) != 1 {
			cli.ShowAppHelp(cliCtx)
			return errors.Errorf("invalid arguments %s", strings.Join(nonFlagArgs, " "))
		}
		targetName := nonFlagArgs[0]
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return errors.Wrapf(err, "parse target name %s", targetName)
		}
	}

	cleanCollection := cleanup.NewCollection()
	defer cleanCollection.Close()

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	// Default upload logs, unless explicitly configured
	if !app.cfg.Global.DisableLogSharing {
		if cloudClient.IsLoggedIn(cliCtx.Context) {
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
		} else {
			defer func() { // Defer this to keep log upload code together
				app.console.Printf("Share your logs with an Earthly account (experimental)! Register for one at https://ci.earthly.dev.")
			}()
		}
	}

	app.console.PrintPhaseHeader(builder.PhaseInit, false, "")
	app.warnIfArgContainsBuildArg(flagArgs)

	err = app.configureSatellite(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrapf(err, "could not construct new buildkit client")
	}

	isLocal := containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress)

	if !isLocal && app.ci {
		app.console.Warnf("Please note that --use-inline-cache and --save-inline-cache are currently disabled when using --ci on Satellites or remote Buildkit.")
		app.console.Warnf("") // newline
		app.useInlineCache = false
		app.saveInlineCache = false
	}

	bkClient, err := buildkitd.NewClient(cliCtx.Context, app.console, app.buildkitdImage, app.containerName, app.containerFrontend, Version, app.buildkitdSettings)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	app.analyticsMetadata.isRemoteBuildkit = !isLocal

	bkIP, err := buildkitd.GetContainerIP(cliCtx.Context, app.containerName, app.containerFrontend, app.buildkitdSettings)
	if err != nil {
		return errors.Wrap(err, "get buildkit container IP")
	}

	nativePlatform, err := platutil.GetNativePlatformViaBkClient(cliCtx.Context, bkClient)
	if err != nil {
		return errors.Wrap(err, "get native platform via buildkit client")
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
		return errors.Errorf("multi-platform builds are not yet supported on the command line. You may, however, create a target with the instruction BUILD --plaform ... --platform ... %s", target)
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

	debuggerSettings := debuggercommon.DebuggerSettings{
		DebugLevelLogging: app.debug,
		Enabled:           app.interactiveDebugging,
		RepeaterAddr:      fmt.Sprintf("%s:8373", bkIP),
		Term:              os.Getenv("TERM"),
	}
	if app.interactiveDebugging {
		analytics.Count("features", "interactive-debugging")
	}

	debuggerSettingsData, err := json.Marshal(&debuggerSettings)
	if err != nil {
		return errors.Wrap(err, "debugger settings json marshal")
	}
	secretsMap[debuggercommon.DebuggerSettingsSecretsKey] = debuggerSettingsData

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

	var enttlmnts []entitlements.Entitlement
	if app.allowPrivileged {
		enttlmnts = append(enttlmnts, entitlements.EntitlementSecurityInsecure)
	}

	if termutil.IsTTY() {
		go func() {
			// Dialing does not accept URLs, it accepts an address and a "network". These cannot be handled as URL schemes.
			// Since Shellrepeater hard-codes TCP, we drop it here and log the error if we fail to connect.

			u, err := url.Parse(app.debuggerHost)
			if err != nil {
				panic("debugger host was not a URL")
			}

			debugTermConsole := app.console.WithPrefix("internal-term")
			err = terminal.ConnectTerm(cliCtx.Context, u.Host, debugTermConsole)
			if err != nil {
				debugTermConsole.VerbosePrintf("unable to connect to terminal: %s", err.Error())
			}
		}()
	}

	overridingVars, err := app.combineVariables(dotEnvMap, flagArgs)
	if err != nil {
		return err
	}

	imageResolveMode := llb.ResolveModePreferLocal
	if app.pull {
		imageResolveMode = llb.ResolveModeForcePull
	}

	cacheImports := make(map[string]bool)
	if app.remoteCache != "" || len(app.cacheFrom.Value()) > 0 {
		cacheImports[app.remoteCache] = true

		for _, c := range app.cacheFrom.Value() {
			cacheImports[c] = true
		}
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
	builderOpts := builder.Opt{
		BkClient:               bkClient,
		Console:                app.console,
		Verbose:                app.verbose,
		Attachables:            attachables,
		Enttlmnts:              enttlmnts,
		NoCache:                app.noCache,
		CacheImports:           states.NewCacheImports(cacheImports),
		CacheExport:            cacheExport,
		MaxCacheExport:         maxCacheExport,
		UseInlineCache:         app.useInlineCache,
		SaveInlineCache:        app.saveInlineCache,
		SessionID:              app.sessionID,
		ImageResolveMode:       imageResolveMode,
		CleanCollection:        cleanCollection,
		OverridingVars:         overridingVars,
		BuildContextProvider:   buildContextProvider,
		GitLookup:              gitLookup,
		UseFakeDep:             !app.noFakeDep,
		Strict:                 app.strict,
		DisableNoOutputUpdates: app.interactiveDebugging,
		ParallelConversion:     (app.cfg.Global.ConversionParallelism != 0),
		Parallelism:            parallelism,
		LocalRegistryAddr:      localRegistryAddr,
		FeatureFlagOverrides:   app.featureFlagOverrides,
		ContainerFrontend:      app.containerFrontend,
		InternalSecretStore:    internalSecretStore,
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
		NoOutput:                   app.noOutput,
		OnlyFinalTargetImages:      app.imageMode,
		PlatformResolver:           platr,
		EnableGatewayClientLogging: app.debug,
		BuiltinArgs:                builtinArgs,

		// explicitly set this to true at the top level (without granting the entitlements.EntitlementSecurityInsecure buildkit option),
		// to differentiate between a user forgetting to run earthly -P, versus a remotely referening an earthfile that requires privileged.
		AllowPrivileged: true,
	}
	if app.artifactMode {
		buildOpts.OnlyArtifact = &artifact
		buildOpts.OnlyArtifactDestPath = destPath
	}
	_, err = b.BuildTarget(cliCtx.Context, target, buildOpts)
	if err != nil {
		return errors.Wrap(err, "build target")
	}

	return nil
}
