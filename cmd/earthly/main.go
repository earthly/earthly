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
	"github.com/joho/godotenv"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/reflectutil"
	"github.com/earthly/earthly/util/stringutil"
)

const (
	// DefaultBuildkitdContainerName is the name of the buildkitd container.
	DefaultBuildkitdContainerName = "earthly-buildkitd"
	// DefaultBuildkitdVolumeName is the name of the docker volume used for storing the cache.
	DefaultBuildkitdVolumeName = "earthly-cache"

	defaultEnvFile = ".env"
	envFileFlag    = "env-file"
)

type earthlyApp struct {
	cliApp      *cli.App
	console     conslogging.ConsoleLogger
	cfg         *config.Config
	sessionID   string
	commandName string
	cliFlags
	analyticsMetadata
}

type cliFlags struct {
	platformsStr              cli.StringSlice
	buildArgs                 cli.StringSlice
	secrets                   cli.StringSlice
	secretFiles               cli.StringSlice
	artifactMode              bool
	imageMode                 bool
	pull                      bool
	push                      bool
	ci                        bool
	output                    bool
	noOutput                  bool
	noCache                   bool
	pruneAll                  bool
	pruneReset                bool
	buildkitdSettings         buildkitd.Settings
	allowPrivileged           bool
	enableProfiler            bool
	buildkitHost              string
	buildkitdImage            string
	containerName             string
	volumeName                string
	cacheFrom                 cli.StringSlice
	remoteCache               string
	maxRemoteCache            bool
	saveInlineCache           bool
	useInlineCache            bool
	configPath                string
	gitUsernameOverride       string
	gitPasswordOverride       string
	interactiveDebugging      bool
	sshAuthSock               string
	verbose                   bool
	dryRun                    bool
	debug                     bool
	homebrewSource            string
	bootstrapNoBuildkit       bool
	bootstrapWithAutocomplete bool
	email                     string
	token                     string
	password                  string
	disableNewLine            bool
	secretFile                string
	secretStdin               bool
	apiServer                 string
	satelliteAddress          string
	writePermission           bool
	registrationPublicKey     string
	dockerfilePath            string
	earthfilePath             string
	earthfileFinalImage       string
	expiry                    string
	termsConditionsPrivacy    bool
	authToken                 string
	noFakeDep                 bool
	enableSourceMap           bool
	configDryRun              bool
	strict                    bool
	conversionParllelism      int
	debuggerHost              string
	certPath                  string
	keyPath                   string
	disableAnalytics          bool
	featureFlagOverrides      string
	localRegistryHost         string
	envFile                   string
	lsShowLong                bool
	lsShowArgs                bool
	containerFrontend         containerutil.ContainerFrontend
	logSharing                bool
	satelliteName             string
	noSatellite               bool
	satelliteFeatureFlags     cli.StringSlice
	userPermission            string
	noBuildkitUpdate          bool
	globalWaitEnd             bool // for feature-flipping builder.go code removal
	projectName               string
	orgName                   string
	invitePermission          string
	inviteMessage             string
}

type analyticsMetadata struct {
	isSatellite      bool
	isRemoteBuildkit bool
	satelliteVersion string
	buildkitPlatform string
	userPlatform     string
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
	go func() {
		receivedSignal := false
		for sig := range sigChan {
			cancel()
			if receivedSignal {
				// This is the second time we have received a signal. Quit immediately.
				fmt.Printf("Received second signal %s. Forcing exit.\n", sig.String())
				os.Exit(9)
			}
			receivedSignal = true
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
	app.unhideFlags(ctx)
	app.autoComplete(ctx)

	exitCode := app.run(ctx, os.Args)
	// app.cfg will be nil when a user runs `earthly --version`;
	// however in all other regular commands app.cfg will be set in app.Before
	if !app.disableAnalytics && app.cfg != nil && !app.cfg.Global.DisableAnalytics {
		// Use a new context, in case the original context is cancelled due to sigint.
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		displayErrors := app.verbose
		cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
		if err != nil && displayErrors {
			app.console.Warnf("unable to start cloud client: %s", err)
		} else if err == nil {
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
					SatelliteVersion: app.analyticsMetadata.satelliteVersion,
					IsRemoteBuildkit: app.analyticsMetadata.isRemoteBuildkit,
					Realtime:         time.Since(startTime),
				},
			)
		}
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
		cliApp:    cli.NewApp(),
		console:   console,
		sessionID: stringutil.RandomAlphanumeric(64),
		cliFlags: cliFlags{
			buildkitdSettings: buildkitd.Settings{},
		},
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

func (app *earthlyApp) before(context *cli.Context) error {
	if app.enableProfiler {
		go profhandler()
	}

	if app.debug {
		app.console = app.console.WithLogLevel(conslogging.Debug)
	} else if app.verbose {
		app.console = app.console.WithLogLevel(conslogging.Verbose)
	}

	if context.IsSet("config") {
		app.console.Printf("loading config values from %q\n", app.configPath)
	}

	if app.containerName != DefaultBuildkitdContainerName {
		// TODO remove this once the debugger port value is randomly supplied by the OS
		// which is required to run multiple instances in parallel. Currently it attempts to bind the same
		// port and fails.
		return fmt.Errorf("buildkit-container-name is not currently supported")
	}

	var yamlData []byte
	var err error
	if app.configPath != "" {
		yamlData, err = config.ReadConfigFile(app.configPath)
		if err != nil {
			if context.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
				return errors.Wrapf(err, "read config")
			}
		}
	}

	app.cfg, err = config.ParseConfigFile(yamlData)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", app.configPath)
	}

	if app.cfg.Git == nil {
		app.cfg.Git = map[string]config.GitConfig{}
	}

	err = app.processDeprecatedCommandOptions(context, app.cfg)
	if err != nil {
		return err
	}

	// Make a small attempt to check if we are not bootstrapped. If not, then do that before we do anything else.
	isBootstrapCmd := false
	for _, f := range context.Args().Slice() {
		isBootstrapCmd = f == "bootstrap"

		if isBootstrapCmd {
			break
		}
	}

	if !isBootstrapCmd && !cliutil.IsBootstrapped() {
		app.bootstrapNoBuildkit = true // Docker may not be available, for instance... like our integration tests.
		err = app.bootstrap(context)
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

func (app *earthlyApp) processDeprecatedCommandOptions(context *cli.Context, cfg *config.Config) error {
	app.warnIfEarth()

	if cfg.Global.CachePath != "" {
		app.console.Warnf("Warning: the setting cache_path is now obsolete and will be ignored")
	}

	if app.conversionParllelism != 0 {
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

	if cfg.Global.DebuggerPort != config.DefaultDebuggerPort {
		app.console.Warnf("Warning: specifying the port using the debugger-port setting is deprecated. Set it in ~/.earthly/config.yml as part of the debugger_host variable; see https://docs.earthly.dev/earthly-config for reference.\n")
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
	rpcRegex := regexp.MustCompile(`(?U)rpc error: code = .+ desc = `)
	err := app.cliApp.RunContext(ctx, args)
	if err != nil {
		ie, isInterpereterError := earthfile2llb.GetInterpreterError(err)

		var failedOutput string
		var buildErr *builder.BuildError
		if errors.As(err, &buildErr) {
			failedOutput = buildErr.VertexLog()
		}
		if app.verbose {
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

		if strings.Contains(err.Error(), "security.insecure is not allowed") {
			app.console.Warnf("Error: --allow-privileged (-P) flag is required\n")
		} else if strings.Contains(err.Error(), "failed to fetch remote") {
			app.console.Warnf("Error: %v\n", err)
			app.console.Printf(
				"Check your git auth settings.\n" +
					"Did you ssh-add today? Need to configure ~/.earthly/config.yml?\n" +
					"For more information see https://docs.earthly.dev/guides/auth\n")
		} else if strings.Contains(err.Error(), "failed to compute cache key") && strings.Contains(err.Error(), ": not found") {
			re := regexp.MustCompile(`("[^"]*"): not found`)
			var matches = re.FindStringSubmatch(err.Error())
			if len(matches) == 2 {
				app.console.Warnf("Error: File not found %v\n", matches[1])
			} else {
				app.console.Warnf("Error: File not found: %v\n", err.Error())
			}
		} else if strings.Contains(failedOutput, "Invalid ELF image for this architecture") {
			app.console.Warnf("Error: %v\n", err)
			app.console.Printf(
				"Are you using --platform to target a different architecture? You may have to manually install QEMU.\n" +
					"For more information see https://docs.earthly.dev/guides/multi-platform\n")
		} else if !app.verbose && rpcRegex.MatchString(err.Error()) {
			baseErr := errors.Cause(err)
			baseErrMsg := rpcRegex.ReplaceAll([]byte(baseErr.Error()), []byte(""))
			app.console.Warnf("Error: %s\n", string(baseErrMsg))
			if bytes.Contains(baseErrMsg, []byte("transport is closing")) {
				app.console.Warnf(
					"It seems that buildkitd is shutting down or it has crashed. " +
						"You can report crashes at https://github.com/earthly/earthly/issues/new.")
				app.printCrashLogs(ctx)
				return 7
			}
		} else if errors.Is(err, buildkitd.ErrBuildkitCrashed) {
			app.console.Warnf("Error: %v\n", err)
			app.console.Warnf(
				"It seems that buildkitd is shutting down or it has crashed. " +
					"You can report crashes at https://github.com/earthly/earthly/issues/new.")
			app.printCrashLogs(ctx)
			return 7
		} else if errors.Is(err, buildkitd.ErrBuildkitStartFailure) {
			app.console.Warnf("Error: %v\n", err)
			app.console.Warnf(
				"It seems that buildkitd had an issue. " +
					"You can report crashes at https://github.com/earthly/earthly/issues/new.")
			app.printCrashLogs(ctx)
			return 6
		} else if isInterpereterError {
			app.console.Warnf("Error: %s\n", ie.Error())
		} else {
			app.console.Warnf("Error: %v\n", err)
		}
		if errors.Is(err, context.Canceled) {
			app.console.Warnf("Context canceled: %v\n", err)
			return 2
		}
		if status.Code(err) == codes.Canceled {
			app.console.Warnf("Context canceled from buildkitd: %v\n", err)
			app.printCrashLogs(ctx)
			return 2
		}
		return 1
	}
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
