package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	_ "net/http/pprof" // enable pprof handlers on net/http listener
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	gsysinfo "github.com/elastic/go-sysinfo"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/earthly/cloud-api/logstream"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/logbus"
	logbussetup "github.com/earthly/earthly/logbus/setup"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/reflectutil"
)

const (
	// DefaultBuildkitdContainerSuffix is the suffix of the buildkitd container.
	DefaultBuildkitdContainerSuffix = "-buildkitd"
	// DefaultBuildkitdVolumeSuffix is the suffix of the docker volume used for storing the cache.
	DefaultBuildkitdVolumeSuffix = "-cache"

	defaultEnvFile = ".env"
	envFileFlag    = "env-file-path"

	defaultArgFile = ".arg"
	argFileFlag    = "arg-file-path"

	defaultSecretFile = ".secret"
	secretFileFlag    = "secret-file-path"
)

var runExitCodeRegexp = regexp.MustCompile(`did not complete successfully: exit code: [^0][0-9]*$`)

type earthlyApp struct {
	cliApp      *cli.App
	console     conslogging.ConsoleLogger
	cfg         *config.Config
	logbusSetup *logbussetup.BusSetup
	logbus      *logbus.Bus
	commandName string
	cliFlags
	analyticsMetadata
}

type cliFlags struct {
	platformsStr                    cli.StringSlice
	buildArgs                       cli.StringSlice
	secrets                         cli.StringSlice
	secretFiles                     cli.StringSlice
	artifactMode                    bool
	imageMode                       bool
	pull                            bool
	push                            bool
	ci                              bool
	output                          bool
	noOutput                        bool
	noCache                         bool
	pruneAll                        bool
	pruneReset                      bool
	pruneTargetSize                 byteSizeValue
	pruneKeepDuration               time.Duration
	buildkitdSettings               buildkitd.Settings
	allowPrivileged                 bool
	enableProfiler                  bool
	buildkitHost                    string
	buildkitdImage                  string
	containerName                   string
	installationName                string
	cacheFrom                       cli.StringSlice
	remoteCache                     string
	maxRemoteCache                  bool
	saveInlineCache                 bool
	useInlineCache                  bool
	configPath                      string
	gitUsernameOverride             string
	gitPasswordOverride             string
	interactiveDebugging            bool
	sshAuthSock                     string
	verbose                         bool
	dryRun                          bool
	debug                           bool
	homebrewSource                  string
	bootstrapNoBuildkit             bool
	bootstrapWithAutocomplete       bool
	email                           string
	token                           string
	password                        string
	disableNewLine                  bool
	secretStdin                     bool
	cloudHTTPAddr                   string
	cloudGRPCAddr                   string
	cloudGRPCInsecure               bool
	satelliteAddress                string
	writePermission                 bool
	registrationPublicKey           string
	dockerfilePath                  string
	earthfilePath                   string
	earthfileFinalImage             string
	expiry                          string
	termsConditionsPrivacy          bool
	authToken                       string
	authJWT                         string
	noFakeDep                       bool
	enableSourceMap                 bool
	configDryRun                    bool
	strict                          bool
	conversionParallelism           int
	certPath                        string
	keyPath                         string
	caPath                          string
	tlsEnabled                      bool
	disableAnalytics                bool
	featureFlagOverrides            string
	localRegistryHost               string
	envFile                         string
	argFile                         string
	secretFile                      string
	lsShowLong                      bool
	lsShowArgs                      bool
	docShowLong                     bool
	containerFrontend               containerutil.ContainerFrontend
	satelliteName                   string
	noSatellite                     bool
	satelliteFeatureFlags           cli.StringSlice
	satellitePlatform               string
	satelliteSize                   string
	satellitePrintJSON              bool
	satelliteMaintenanceWindow      string
	satelliteMaintenaceWeekendsOnly bool
	satelliteDropCache              bool
	satelliteVersion                string
	satelliteIncludeHidden          bool
	userPermission                  string
	noBuildkitUpdate                bool
	globalWaitEnd                   bool // for feature-flipping builder.go code removal
	projectName                     string
	forceRemoveProject              bool
	orgName                         string
	invitePermission                string
	inviteMessage                   string
	logstream                       bool
	logstreamUpload                 bool
	logstreamDebugFile              string
	logstreamDebugManifestFile      string
	logstreamAddressOverride        string
	requestID                       string
	buildID                         string
	loginProvider                   string
	registryUsername                string
	registryPassword                string
	registryPasswordStdin           bool
	registryCredHelper              string
	awsAccessKeyID                  string
	awsSecretAccessKey              string
	gcpServiceAccountKeyPath        string
	gcpServiceAccountKey            string
	gcpServiceAccountKeyStdin       bool
}

type analyticsMetadata struct {
	isSatellite             bool
	isRemoteBuildkit        bool
	satelliteCurrentVersion string
	buildkitPlatform        string
	userPlatform            string
}

var (
	// DefaultBuildkitdImage is the default buildkitd image to use.
	DefaultBuildkitdImage string
	// Version is the version of this CLI app.
	Version string
	// GitSha contains the git sha used to build this app
	GitSha string
	// BuiltBy contains information on which build-system was used (e.g. official earthly binaries, homebrew, etc)
	BuiltBy string
	// DefaultInstallationName is the name included in the various earthly global resources on the system,
	// such as the ~/.earthly dir name, the buildkitd container name, the docker volume name, etc.
	// This should be set to "earthly" for official releases.
	DefaultInstallationName string
)

func main() {
	startTime := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()
	var lastSignal os.Signal
	go func() {
		for sig := range sigChan {
			cancel()
			if lastSignal != nil {
				// This is the second time we have received a signal. Quit immediately.
				fmt.Printf("Received second signal %s. Forcing exit.\n", sig.String())
				os.Exit(9)
			}
			lastSignal = sig
			fmt.Printf("Received signal %s. Cleaning up before exiting...\n", sig.String())
			go func() {
				// Wait for 30 seconds before forcing an exit.
				time.Sleep(30 * time.Second)
				fmt.Printf("Timed out cleaning up. Forcing exit.\n")
				os.Exit(9)
			}()
		}
	}()
	// Occasional spurious warnings show up - these are coming from imported libraries. Discard them.
	logrus.StandardLogger().Out = io.Discard

	// Load .env into current global env's. This is mainly for applying Earthly settings.
	// Separate call is made for build args and secrets.
	envFile := defaultEnvFile
	envFileOverride := false
	if envFileFromEnv, ok := os.LookupEnv("EARTHLY_ENV_FILE"); ok {
		envFile = envFileFromEnv
		envFileOverride = true
	}
	envFileFromArgOK := true
	flagSet := flag.NewFlagSet(getBinaryName(), flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	for _, f := range newEarthlyApp(ctx, conslogging.ConsoleLogger{}).cliApp.Flags {
		if err := f.Apply(flagSet); err != nil {
			envFileFromArgOK = false
			break
		}
	}
	if envFileFromArgOK {
		if err := flagSet.Parse(os.Args[1:]); err == nil {
			if envFileFlag := flagSet.Lookup(envFileFlag); envFileFlag != nil {
				envFile = envFileFlag.Value.String()
				envFileOverride = envFile != defaultEnvFile // flag lib doesn't expose if a value was set or not
			}
		}
	}
	err := godotenv.Load(envFile)
	if err != nil {
		// ignore ErrNotExist when using default .env file
		if envFileOverride || !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Error loading dot-env file %s: %s\n", envFile, err.Error())
			os.Exit(1)
		}
	}

	colorMode := conslogging.AutoColor
	_, forceColor := os.LookupEnv("FORCE_COLOR")
	if forceColor {
		colorMode = conslogging.ForceColor
		color.NoColor = false
	}
	_, noColor := os.LookupEnv("NO_COLOR")
	if noColor {
		colorMode = conslogging.NoColor
		color.NoColor = true
	}

	padding := conslogging.DefaultPadding
	customPadding, ok := os.LookupEnv("EARTHLY_TARGET_PADDING")
	if ok {
		targetPadding, err := strconv.Atoi(customPadding)
		if err == nil {
			padding = targetPadding
		}
	}

	_, fullTarget := os.LookupEnv("EARTHLY_FULL_TARGET")
	if fullTarget {
		padding = conslogging.NoPadding
	}

	app := newEarthlyApp(ctx, conslogging.Current(colorMode, padding, conslogging.Info))
	err = app.unhideFlags(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error un-hiding flags %v", err)
		os.Exit(1)
	}
	app.autoComplete(ctx)

	exitCode := app.run(ctx, os.Args)
	// app.cfg will be nil when a user runs `earthly --version`;
	// however in all other regular commands app.cfg will be set in app.Before
	if !app.disableAnalytics && app.cfg != nil && !app.cfg.Global.DisableAnalytics {
		// Use a new context, in case the original context is cancelled due to sigint.
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		displayErrors := app.verbose
		cloudClient, err := app.newCloudClient()
		if err != nil && displayErrors {
			app.console.Warnf("unable to start cloud client: %s", err)
		} else if err == nil {
			analytics.AddCLIProject(app.orgName, app.projectName)
			org, project := analytics.ProjectDetails()
			analytics.CollectAnalytics(
				ctxTimeout, cloudClient, displayErrors, analytics.Meta{
					Version:          Version,
					Platform:         getPlatform(),
					BuildkitPlatform: app.analyticsMetadata.buildkitPlatform,
					UserPlatform:     app.analyticsMetadata.userPlatform,
					GitSHA:           GitSha,
					CommandName:      app.commandName,
					ExitCode:         exitCode,
					IsSatellite:      app.analyticsMetadata.isSatellite,
					SatelliteVersion: app.analyticsMetadata.satelliteCurrentVersion,
					IsRemoteBuildkit: app.analyticsMetadata.isRemoteBuildkit,
					Realtime:         time.Since(startTime),
					OrgName:          org,
					ProjectName:      project,
				},
				app.installationName,
			)
		}
	}
	if lastSignal != nil {
		app.console.PrintBar(color.New(color.FgHiYellow), fmt.Sprintf("WARNING: exiting due to received %s signal", lastSignal.String()), "")
	}
	os.Exit(exitCode)
}

func getVersionPlatform() string {
	s := fmt.Sprintf("%s %s %s", Version, GitSha, getPlatform())
	if BuiltBy != "" {
		s += " " + BuiltBy
	}
	return s
}

func getPlatform() string {
	h, err := gsysinfo.Host()
	if err != nil {
		return "unknown"
	}
	info := h.Info()
	return fmt.Sprintf("%s/%s; %s %s", runtime.GOOS, runtime.GOARCH, info.OS.Name, info.OS.Version)
}

func getBinaryName() string {
	if len(os.Args) == 0 {
		return "earthly"
	}
	binPath := os.Args[0] // can't use os.Executable() here; because it will give us earthly if executed via the earth symlink
	baseName := path.Base(binPath)
	return baseName
}

func newEarthlyApp(ctx context.Context, console conslogging.ConsoleLogger) *earthlyApp {
	app := &earthlyApp{
		cliApp:  cli.NewApp(),
		console: console,
		cliFlags: cliFlags{
			buildkitdSettings: buildkitd.Settings{},
		},
		logbus: logbus.New(),
	}

	earthly := getBinaryName()

	app.cliApp.Usage = "The CI/CD framework that runs anywhere"
	app.cliApp.UsageText = "\t" + earthly + " [options] <target-ref>\n" +
		"\n" +
		"   \t" + earthly + " [options] --image <target-ref>\n" +
		"\n" +
		"   \t" + earthly + " [options] --artifact <target-ref>/<artifact-path> [<dest-path>]\n" +
		"\n" +
		"   \t" + earthly + " [options] command [command options]\n" +
		"\n" +
		"Executes Earthly builds. For more information see https://docs.earthly.dev/docs/earthly-command.\n" +
		"To get started with using Earthly, check out the getting started guide at https://docs.earthly.dev/basics.\n" +
		"\n" +
		"For help on build-specific flags try \n" +
		"\n" +
		"\t" + earthly + " build --help"
	app.cliApp.UseShortOptionHandling = true
	app.cliApp.Action = app.actionBuild
	app.cliApp.Version = getVersionPlatform()

	app.cliApp.Flags = app.rootFlags()                                     // These will show up in help.
	app.cliApp.Flags = append(app.cliApp.Flags, app.hiddenBuildFlags()...) // These will not.

	app.cliApp.Commands = app.rootCmds()

	app.cliApp.Before = app.before
	return app
}

func (app *earthlyApp) before(cliCtx *cli.Context) error {
	if app.enableProfiler {
		go profhandler()
	}

	if app.installationName != "" {
		if !cliCtx.IsSet("config") {
			app.configPath = defaultConfigPath(app.installationName)
		}
		if !cliCtx.IsSet("buildkit-container-name") {
			app.containerName = fmt.Sprintf("%s-buildkitd", app.installationName)
		}
		if !cliCtx.IsSet("buildkit-volume-name") {
			app.buildkitdSettings.VolumeName = fmt.Sprintf("%s-cache", app.installationName)
		}
	}
	if app.debug {
		app.console = app.console.WithLogLevel(conslogging.Debug)
	} else if app.verbose {
		app.console = app.console.WithLogLevel(conslogging.Verbose)
	}
	if app.logstreamUpload {
		app.logstream = true
	}
	if app.logstream {
		app.console = app.console.WithPrefixWriter(app.logbus.Run().Generic())
		if app.buildID == "" {
			app.buildID = uuid.NewString()
		}
		disableOngoingUpdates := !app.logstream || app.interactiveDebugging
		_, forceColor := os.LookupEnv("FORCE_COLOR")
		_, noColor := os.LookupEnv("NO_COLOR")
		var err error
		app.logbusSetup, err = logbussetup.New(
			cliCtx.Context, app.logbus, app.debug, app.verbose, forceColor, noColor,
			disableOngoingUpdates, app.logstreamDebugFile, app.buildID)
		if err != nil {
			return errors.Wrap(err, "logbus setup")
		}
	}

	if cliCtx.IsSet("config") {
		app.console.Printf("loading config values from %q\n", app.configPath)
	}

	var yamlData []byte
	if app.configPath != "" {
		var err error
		yamlData, err = config.ReadConfigFile(app.configPath)
		if err != nil {
			if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
				return errors.Wrapf(err, "read config")
			}
		}
	}

	var err error
	app.cfg, err = config.ParseConfigFile(yamlData, app.installationName)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", app.configPath)
	}

	if app.cfg.Git == nil {
		app.cfg.Git = map[string]config.GitConfig{}
	}

	err = app.processDeprecatedCommandOptions(cliCtx, app.cfg)
	if err != nil {
		return err
	}

	// Make a small attempt to check if we are not bootstrapped. If not, then do that before we do anything else.
	isBootstrapCmd := false
	for _, f := range cliCtx.Args().Slice() {
		isBootstrapCmd = f == "bootstrap"

		if isBootstrapCmd {
			break
		}
	}

	if !isBootstrapCmd && !cliutil.IsBootstrapped(app.installationName) {
		app.bootstrapNoBuildkit = true // Docker may not be available, for instance... like our integration tests.
		err = app.bootstrap(cliCtx)
		if err != nil {
			return errors.Wrap(err, "bootstrap unbootstrapped installation")
		}
	}

	return nil
}

func (app *earthlyApp) warnIfEarth() {
	if len(os.Args) == 0 {
		return
	}
	binPath := os.Args[0] // can't use os.Executable() here; because it will give us earthly if executed via the earth symlink

	baseName := path.Base(binPath)
	if baseName == "earth" {
		app.console.Warnf("Warning: the earth binary has been renamed to earthly; the earth command is currently symlinked, but is deprecated and will one day be removed.")

		absPath, err := filepath.Abs(binPath)
		if err != nil {
			return
		}
		earthlyPath := path.Join(path.Dir(absPath), "earthly")
		earthlyPathExists, _ := fileutil.FileExists(earthlyPath)
		if earthlyPathExists {
			app.console.Warnf("Once you are ready to switch over to earthly, you can `rm %s`", absPath)
		}
	}
}

func (app *earthlyApp) processDeprecatedCommandOptions(cliCtx *cli.Context, cfg *config.Config) error {
	app.warnIfEarth()

	if cfg.Global.CachePath != "" {
		app.console.Warnf("Warning: the setting cache_path is now obsolete and will be ignored")
	}

	if app.conversionParallelism != 0 {
		app.console.Warnf("Warning: --conversion-parallelism and EARTHLY_CONVERSION_PARALLELISM is obsolete, please use 'earthly config global.conversion_parallelism <parallelism>' instead")
	}

	// command line overrides the config file
	if app.gitUsernameOverride != "" || app.gitPasswordOverride != "" {
		app.console.Warnf("Warning: the --git-username and --git-password command flags are deprecated and are now configured in the ~/.earthly/config.yml file under the git section; see https://docs.earthly.dev/earthly-config for reference.\n")
		if _, ok := cfg.Git["github.com"]; !ok {
			cfg.Git["github.com"] = config.GitConfig{}
		}
		if _, ok := cfg.Git["gitlab.com"]; !ok {
			cfg.Git["gitlab.com"] = config.GitConfig{}
		}

		for k, v := range cfg.Git {
			v.Auth = "https"
			if app.gitUsernameOverride != "" {
				v.User = app.gitUsernameOverride
			}
			if app.gitPasswordOverride != "" {
				v.Password = app.gitPasswordOverride
			}
			cfg.Git[k] = v
		}
	}

	return nil
}

func (app *earthlyApp) unhideFlags(ctx context.Context) error {
	var err error
	if os.Getenv("EARTHLY_AUTOCOMPLETE_HIDDEN") != "" && os.Getenv("COMP_POINT") == "" { // TODO delete this check after 2022-03-01
		// only display warning when NOT under complete mode (otherwise we break auto completion)
		app.console.Warnf("Warning: EARTHLY_AUTOCOMPLETE_HIDDEN has been renamed to EARTHLY_SHOW_HIDDEN\n")
	}
	showHidden := false
	showHiddenStr := os.Getenv("EARTHLY_SHOW_HIDDEN")
	if showHiddenStr != "" {
		showHidden, err = strconv.ParseBool(showHiddenStr)
		if err != nil {
			return err
		}
	}
	if !showHidden {
		return nil
	}

	for _, fl := range app.cliApp.Flags {
		reflectutil.SetBool(fl, "Hidden", false)
	}

	unhideFlagsCommands(ctx, app.cliApp.Commands)

	return nil
}

func unhideFlagsCommands(ctx context.Context, cmds []*cli.Command) {
	for _, cmd := range cmds {
		reflectutil.SetBool(cmd, "Hidden", false)
		for _, flg := range cmd.Flags {
			reflectutil.SetBool(flg, "Hidden", false)
		}
		unhideFlagsCommands(ctx, cmd.Subcommands)
	}
}

func (app *earthlyApp) run(ctx context.Context, args []string) int {
	defer func() {
		if app.logstream {
			err := app.logbusSetup.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error(s) in logbus: %v", err)
			}
			if app.logstreamDebugManifestFile != "" {
				err := app.logbusSetup.DumpManifestToFile(app.logstreamDebugManifestFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error dumping manifest: %v", err)
				}
			}
		}
	}()
	app.logbus.Run().SetStart(time.Now())
	defer func() {
		// Just in case this is forgotten somewhere else.
		app.logbus.Run().SetFatalError(
			time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER,
			"No SetFatalError called appropriately. This should never happen.")
	}()
	rpcRegex := regexp.MustCompile(`(?U)rpc error: code = .+ desc = `)
	err := app.cliApp.RunContext(ctx, args)
	if err != nil {
		ie, isInterpreterError := earthfile2llb.GetInterpreterError(err)

		var failedOutput string
		var buildErr *builder.BuildError
		if errors.As(err, &buildErr) {
			failedOutput = buildErr.VertexLog()
		}
		if app.debug {
			// Get the stack trace from the deepest error that has it and print it.
			type stackTracer interface {
				StackTrace() errors.StackTrace
			}
			errChain := []error{}
			for it := err; it != nil; it = errors.Unwrap(it) {
				errChain = append(errChain, it)
			}
			for index := len(errChain) - 1; index > 0; index-- {
				it := errChain[index]
				errWithStack, ok := it.(stackTracer)
				if ok {
					app.console.Warnf("Error stack trace:%+v\n", errWithStack.StackTrace())
					break
				}
			}
		}

		switch {
		case runExitCodeRegexp.MatchString(err.Error()):
			// error has already been displayed in console, don't display it again
			return 1
		case strings.Contains(err.Error(), "security.insecure is not allowed"):
			app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_NEEDS_PRIVILEGED, err.Error())
			app.console.Warnf("Error: earthly --allow-privileged (earthly -P) flag is required\n")
			return 9
		case strings.Contains(err.Error(), "failed to fetch remote"):
			app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_GIT, err.Error())
			app.console.Warnf("Error: %v\n", err)
			app.console.Printf(
				"Check your git auth settings.\n" +
					"Did you ssh-add today? Need to configure ~/.earthly/config.yml?\n" +
					"For more information see https://docs.earthly.dev/guides/auth\n")
			return 1
		case strings.Contains(err.Error(), "failed to compute cache key") && strings.Contains(err.Error(), ": not found"):
			app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND, err.Error())
			re := regexp.MustCompile(`("[^"]*"): not found`)
			var matches = re.FindStringSubmatch(err.Error())
			if len(matches) == 2 {
				app.console.Warnf("Error: File not found %v\n", matches[1])
			} else {
				app.console.Warnf("Error: File not found: %v\n", err.Error())
			}
			return 1
		case strings.Contains(failedOutput, "Invalid ELF image for this architecture"):
			app.console.Warnf("Error: %v\n", err)
			app.console.Printf(
				"Are you using --platform to target a different architecture? You may have to manually install QEMU.\n" +
					"For more information see https://docs.earthly.dev/guides/multi-platform\n")
			return 1
		case !app.verbose && rpcRegex.MatchString(err.Error()):
			baseErr := errors.Cause(err)
			baseErrMsg := rpcRegex.ReplaceAll([]byte(baseErr.Error()), []byte(""))
			app.console.Warnf("Error: %s\n", string(baseErrMsg))
			if bytes.Contains(baseErrMsg, []byte("transport is closing")) {
				app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED, baseErr.Error())
				app.console.Warnf(
					"It seems that buildkitd is shutting down or it has crashed. " +
						"You can report crashes at https://github.com/earthly/earthly/issues/new.")
				if containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress) {
					app.printCrashLogs(ctx)
				}
				return 7
			}
			return 1
		case errors.Is(err, buildkitd.ErrBuildkitCrashed):
			app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED, err.Error())
			app.console.Warnf("Error: %v\n", err)
			app.console.Warnf(
				"It seems that buildkitd is shutting down or it has crashed. " +
					"You can report crashes at https://github.com/earthly/earthly/issues/new.")
			if containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress) {
				app.printCrashLogs(ctx)
			}
			return 7
		case errors.Is(err, buildkitd.ErrBuildkitConnectionFailure):
			app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_CONNECTION_FAILURE, err.Error())
			app.console.Warnf("Error: %v\n", err)
			if containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress) {
				app.console.Warnf(
					"It seems that buildkitd had an issue. " +
						"You can report crashes at https://github.com/earthly/earthly/issues/new.")
				app.printCrashLogs(ctx)
			}
			return 6
		case errors.Is(err, context.Canceled):
			app.logbus.Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_CANCELED)
			app.console.Warnf("Canceled\n")
			return 2
		case status.Code(errors.Cause(err)) == codes.Canceled:
			app.logbus.Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_CANCELED)
			app.console.Warnf("Canceled\n")
			if containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress) {
				app.printCrashLogs(ctx)
			}
			return 2
		case isInterpreterError:
			app.logbus.Run().SetFatalError(time.Now(), ie.TargetID, "", logstream.FailureType_FAILURE_TYPE_SYNTAX, ie.Error())
			app.console.Warnf("Error: %s\n", ie.Error())
			return 1
		default:
			app.logbus.Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER, err.Error())
			app.console.Warnf("Error: %v\n", err)
			return 1
		}
	}
	app.logbus.Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_SUCCESS)
	return 0
}

func (app *earthlyApp) printCrashLogs(ctx context.Context) {
	app.console.PrintBar(color.New(color.FgHiRed), "System Info", "")
	fmt.Fprintf(os.Stderr, "version: %s\n", Version)
	fmt.Fprintf(os.Stderr, "build-sha: %s\n", GitSha)
	fmt.Fprintf(os.Stderr, "platform: %s\n", getPlatform())

	dockerVersion, err := buildkitd.GetDockerVersion(ctx, app.containerFrontend)
	if err != nil {
		app.console.Warnf("failed querying docker version: %s\n", err.Error())
	} else {
		app.console.PrintBar(color.New(color.FgHiRed), "Docker Version", "")
		fmt.Fprintln(os.Stderr, dockerVersion)
	}

	logs, err := buildkitd.GetLogs(ctx, app.containerName, app.containerFrontend, app.buildkitdSettings)
	if err != nil {
		app.console.Warnf("failed fetching earthly-buildkit logs: %s\n", err.Error())
	} else {
		app.console.PrintBar(color.New(color.FgHiRed), "Buildkit Logs", "")
		fmt.Fprintln(os.Stderr, logs)
	}
}
