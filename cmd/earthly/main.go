package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" // enable pprof handlers on net/http listener
	"net/url"
	"os"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/localhost/localhostprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wille/osutil"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/autocomplete"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildcontext/provider"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/terminal"
	"github.com/earthly/earthly/docker2earthly"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/llbutil/secretprovider"
	"github.com/earthly/earthly/util/platutil"
	"github.com/earthly/earthly/util/reflectutil"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/termutil"
	"github.com/earthly/earthly/variables"
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
	satelliteOrg              string
}

var (
	// DefaultBuildkitdImage is the default buildkitd image to use.
	DefaultBuildkitdImage string

	// Version is the version of this CLI app.
	Version string

	// GitSha contains the git sha used to build this app
	GitSha string
)

var (
	errLoginFlagsHaveNoEffect            = errors.New("account login flags have no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")
	errLogoutHasNoEffectWhenAuthTokenSet = errors.New("account logout has no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")
)

func profhandler() {
	addr := "127.0.0.1:6060"
	fmt.Printf("listening for pprof on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

func main() {
	startTime := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		receivedSignal := false
		for sig := range c {
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
		ctxTimeout, cancel := context.WithTimeout(ctx, time.Millisecond*500)
		defer cancel()
		displayErrors := app.verbose
		cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
		if err != nil && displayErrors {
			app.console.Warnf("unable to start cloud client: %s", err)
		} else if err == nil {
			analytics.CollectAnalytics(
				ctxTimeout, cc, displayErrors, Version, getPlatform(),
				GitSha, app.commandName, exitCode, time.Since(startTime),
			)
		}
	}
	os.Exit(exitCode)
}

func getVersionPlatform() string {
	return fmt.Sprintf("%s %s %s", Version, GitSha, getPlatform())
}

func getPlatform() string {
	// Work-around for windows panics; this can be removed once https://github.com/wille/osutil/pull/10 is merged
	showOSInfo := func() (info string) {
		defer func() {
			if err := recover(); err != nil {
				// skipped
				info = "unknown"
				return
			}
		}()
		return osutil.GetDisplay()
	}
	return fmt.Sprintf("%s/%s; %s", runtime.GOOS, runtime.GOARCH, showOSInfo())
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
	sessionIDBytes := make([]byte, 64)
	_, err := rand.Read(sessionIDBytes)
	if err != nil {
		panic(err)
	}
	app := &earthlyApp{
		cliApp:    cli.NewApp(),
		console:   console,
		sessionID: base64.StdEncoding.EncodeToString(sessionIDBytes),
		cliFlags: cliFlags{
			buildkitdSettings: buildkitd.Settings{},
		},
	}

	earthly := getBinaryName()

	app.cliApp.Usage = "A build automation tool for the container era"
	app.cliApp.UsageText = "\t" + earthly + " [options] <target-ref>\n" +
		"\n" +
		"   \t" + earthly + " [options] --image <target-ref>\n" +
		"\n" +
		"   \t" + earthly + " [options] --artifact <target-ref>/<artifact-path> [<dest-path>]\n" +
		"\n" +
		"   \t" + earthly + " [options] command [command options]\n" +
		"\n" +
		"Executes Earthly builds. For more information see https://docs.earthly.dev/earthly-command.\n" +
		"To get started with using Earthly, check out the getting started guide at https://docs.earthly.dev/guides/basics."
	app.cliApp.UseShortOptionHandling = true
	app.cliApp.Action = app.actionBuild
	app.cliApp.Version = getVersionPlatform()
	app.cliApp.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "platform",
			EnvVars: []string{"EARTHLY_PLATFORMS"},
			Usage:   "Specify the target platform to build for",
			Value:   &app.platformsStr,
		},
		&cli.StringSliceFlag{
			Name:    "build-arg",
			EnvVars: []string{"EARTHLY_BUILD_ARGS"},
			Usage:   "A build arg override, specified as <key>=[<value>]",
			Value:   &app.buildArgs,
		},
		&cli.StringSliceFlag{
			Name:    "secret",
			Aliases: []string{"s"},
			EnvVars: []string{"EARTHLY_SECRETS"},
			Usage:   "A secret override, specified as <key>=[<value>]",
			Value:   &app.secrets,
		},
		&cli.StringSliceFlag{
			Name:    "secret-file",
			EnvVars: []string{"EARTHLY_SECRET_FILES"},
			Usage:   "A secret override, specified as <key>=<path>",
			Value:   &app.secretFiles,
		},
		&cli.BoolFlag{
			Name:        "artifact",
			Aliases:     []string{"a"},
			Usage:       "Output specified artifact; a wildcard (*) can be used to output all artifacts",
			Destination: &app.artifactMode,
		},
		&cli.BoolFlag{
			Name:        "image",
			Usage:       "Output only docker image of the specified target",
			Destination: &app.imageMode,
		},
		&cli.BoolFlag{
			Name:        "pull",
			EnvVars:     []string{"EARTHLY_PULL"},
			Usage:       "Force pull any referenced Docker images",
			Destination: &app.pull,
		},
		&cli.BoolFlag{
			Name:        "push",
			EnvVars:     []string{"EARTHLY_PUSH"},
			Usage:       "Push docker images and execute RUN --push commands",
			Destination: &app.push,
		},
		&cli.BoolFlag{
			Name:        "ci",
			EnvVars:     []string{"EARTHLY_CI"},
			Usage:       "Execute in CI mode (implies --use-inline-cache --save-inline-cache --no-output --strict)",
			Destination: &app.ci,
		},
		&cli.BoolFlag{
			Name:        "output",
			EnvVars:     []string{"EARTHLY_OUTPUT"},
			Usage:       "Allow artifacts or images to be output, even when running under --ci mode",
			Destination: &app.output,
		},
		&cli.BoolFlag{
			Name:        "no-output",
			EnvVars:     []string{"EARTHLY_NO_OUTPUT"},
			Usage:       wrap("Do not output artifacts or images", "(using --push is still allowed)"),
			Destination: &app.noOutput,
		},
		&cli.BoolFlag{
			Name:        "no-cache",
			EnvVars:     []string{"EARTHLY_NO_CACHE"},
			Usage:       "Do not use cache while building",
			Destination: &app.noCache,
		},
		&cli.StringFlag{
			Name:        "config",
			Value:       defaultConfigPath(),
			EnvVars:     []string{"EARTHLY_CONFIG"},
			Usage:       "Path to config file",
			Destination: &app.configPath,
		},
		&cli.StringFlag{
			Name:        "ssh-auth-sock",
			Value:       os.Getenv("SSH_AUTH_SOCK"),
			EnvVars:     []string{"EARTHLY_SSH_AUTH_SOCK"},
			Usage:       wrap("The SSH auth socket to use for ssh-agent forwarding", ""),
			Destination: &app.sshAuthSock,
		},
		&cli.StringFlag{
			Name:        "auth-token",
			EnvVars:     []string{"EARTHLY_TOKEN"},
			Usage:       "Force Earthly account login to authenticate with supplied token",
			Destination: &app.authToken,
		},
		&cli.StringFlag{
			Name:        "git-username",
			EnvVars:     []string{"GIT_USERNAME"},
			Usage:       "The git username to use for git HTTPS authentication",
			Destination: &app.gitUsernameOverride,
		},
		&cli.StringFlag{
			Name:        "git-password",
			EnvVars:     []string{"GIT_PASSWORD"},
			Usage:       "The git password to use for git HTTPS authentication",
			Destination: &app.gitPasswordOverride,
		},
		&cli.BoolFlag{
			Name:        "allow-privileged",
			Aliases:     []string{"P"},
			EnvVars:     []string{"EARTHLY_ALLOW_PRIVILEGED"},
			Usage:       "Allow build to use the --privileged flag in RUN commands",
			Destination: &app.allowPrivileged,
		},
		&cli.BoolFlag{
			Name:        "profiler",
			EnvVars:     []string{"EARTHLY_PROFILER"},
			Usage:       "Enable the profiler",
			Destination: &app.enableProfiler,
			Hidden:      true, // Dev purposes only.
		},
		&cli.StringFlag{
			Name:        "buildkit-host",
			Value:       "",
			EnvVars:     []string{"EARTHLY_BUILDKIT_HOST"},
			Usage:       wrap("The URL to use for connecting to a buildkit host. ", "If empty, earthly will attempt to start a buildkitd instance via docker run"),
			Destination: &app.buildkitHost,
		},
		&cli.StringFlag{
			Name:        "debugger-host",
			EnvVars:     []string{"EARTHLY_DEBUGGER_HOST"},
			Usage:       wrap("The URL to use for connecting to a debugger host. ", "If empty, earthly uses the default debugger port, combined with the desired buildkit host."),
			Destination: &app.debuggerHost,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "tlscert",
			Value:       "./certs/earthly_cert.pem",
			EnvVars:     []string{"EARTHLY_TLS_CERT"},
			Usage:       wrap("The path to the client TLS cert", "If relative, will be interpreted as relative to the ~/.earthly folder."),
			Destination: &app.certPath,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "tlskey",
			Value:       "./certs/earthly_key.pem",
			EnvVars:     []string{"EARTHLY_TLS_KEY"},
			Usage:       wrap("The path to the client TLS key.", "If relative, will be interpreted as relative to the ~/.earthly folder."),
			Destination: &app.keyPath,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "buildkit-image",
			Value:       DefaultBuildkitdImage,
			EnvVars:     []string{"EARTHLY_BUILDKIT_IMAGE"},
			Usage:       "The docker image to use for the buildkit daemon",
			Destination: &app.buildkitdImage,
		},
		&cli.StringFlag{
			Name:        "buildkit-container-name",
			Value:       DefaultBuildkitdContainerName,
			EnvVars:     []string{"EARTHLY_CONTAINER_NAME"},
			Usage:       "The docker container name to use for the buildkit daemon",
			Destination: &app.containerName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "buildkit-volume-name",
			Value:       DefaultBuildkitdVolumeName,
			EnvVars:     []string{"EARTHLY_VOLUME_NAME"},
			Usage:       "The docker volume name to use for the buildkit daemon cache",
			Destination: &app.buildkitdSettings.VolumeName,
			Hidden:      true,
		},
		&cli.StringFlag{
			Name:        "remote-cache",
			EnvVars:     []string{"EARTHLY_REMOTE_CACHE"},
			Usage:       "A remote docker image tag use as explicit cache",
			Destination: &app.remoteCache,
		},
		&cli.BoolFlag{
			Name:        "max-remote-cache",
			EnvVars:     []string{"EARTHLY_MAX_REMOTE_CACHE"},
			Usage:       "Saves all intermediate images too in the remote cache",
			Destination: &app.maxRemoteCache,
		},
		&cli.BoolFlag{
			Name:        "save-inline-cache",
			EnvVars:     []string{"EARTHLY_SAVE_INLINE_CACHE"},
			Usage:       "Enable cache inlining when pushing images",
			Destination: &app.saveInlineCache,
		},
		&cli.BoolFlag{
			Name:        "use-inline-cache",
			EnvVars:     []string{"EARTHLY_USE_INLINE_CACHE"},
			Usage:       wrap("Attempt to use any inline cache that may have been previously pushed ", "uses image tags referenced by SAVE IMAGE --push or SAVE IMAGE --cache-from"),
			Destination: &app.useInlineCache,
		},
		&cli.BoolFlag{
			Name:        "interactive",
			Aliases:     []string{"i"},
			EnvVars:     []string{"EARTHLY_INTERACTIVE"},
			Usage:       "Enable interactive debugging",
			Destination: &app.interactiveDebugging,
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"V"},
			EnvVars:     []string{"EARTHLY_VERBOSE"},
			Usage:       "Enable verbose logging",
			Destination: &app.verbose,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"D"},
			EnvVars:     []string{"EARTHLY_DEBUG"},
			Usage:       "Enable debug mode. This flag also turns on the debug mode of buildkitd, which may cause it to restart",
			Destination: &app.debug,
			Hidden:      true, // For development purposes only.
		},
		&cli.StringFlag{
			Name:        "server",
			Value:       "https://api.earthly.dev",
			EnvVars:     []string{"EARTHLY_SERVER"},
			Usage:       "API server override for dev purposes",
			Destination: &app.apiServer,
			Hidden:      true, // Internal.
		},
		&cli.StringFlag{
			Name:        "satellite",
			Value:       containerutil.SatelliteAddress,
			EnvVars:     []string{"EARTHLY_SATELLITE"},
			Usage:       "Satellite address override for dev purposes",
			Destination: &app.satelliteAddress,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "no-fake-dep",
			EnvVars:     []string{"EARTHLY_NO_FAKE_DEP"},
			Usage:       "Internal feature flag for fake-dep",
			Destination: &app.noFakeDep,
			Hidden:      true, // Internal.
		},
		&cli.BoolFlag{
			Name:        "strict",
			EnvVars:     []string{"EARTHLY_STRICT"},
			Usage:       "Disallow usage of features that may create unrepeatable builds",
			Destination: &app.strict,
		},
		// TODO: completely remove conversion-parallelism in some future release
		&cli.IntFlag{
			Name:        "conversion-parallelism",
			EnvVars:     []string{"EARTHLY_CONVERSION_PARALLELISM"},
			Usage:       "This flag is obsolete, use 'earthly config global.conversion_parallelism <parallelism>' instead'",
			Destination: &app.conversionParllelism,
			Hidden:      true, // obsolete in favor of config
		},
		&cli.BoolFlag{
			EnvVars:     []string{"EARTHLY_DISABLE_ANALYTICS", "DO_NOT_TRACK"},
			Usage:       "Disable collection of analytics",
			Destination: &app.disableAnalytics,
		},
		&cli.StringFlag{
			Name:        "version-flag-overrides",
			EnvVars:     []string{"EARTHLY_VERSION_FLAG_OVERRIDES"},
			Usage:       "Apply additional flags after each VERSION command across all Earthfiles, multiple flags can be seperated by commas",
			Destination: &app.featureFlagOverrides,
			Hidden:      true, // used for feature-flipping from ./earthly dev script
		},
		&cli.StringFlag{
			Name:        envFileFlag,
			EnvVars:     []string{"EARTHLY_ENV_FILE"},
			Usage:       "Use values from this file as earthly environment variables, buildargs, or secrets",
			Value:       defaultEnvFile,
			Destination: &app.envFile,
		},
	}

	app.cliApp.Commands = []*cli.Command{
		{
			Name:        "bootstrap",
			Usage:       "Bootstraps earthly installation including shell autocompletion and buildkit image download",
			Description: "Bootstraps earthly installation including shell autocompletion and buildkit image download",
			Action:      app.actionBootstrap,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "source",
					Usage:       "Output source file (for use in homebrew install)",
					Hidden:      true, // only meant for use with homebrew formula
					Destination: &app.homebrewSource,
				},
				&cli.BoolFlag{
					Name:        "no-buildkit",
					Usage:       "Do not bootstrap buildkit",
					Destination: &app.bootstrapNoBuildkit,
				},
				&cli.BoolFlag{
					Name:        "with-autocomplete",
					Usage:       "Add earthly autocompletions",
					Destination: &app.bootstrapWithAutocomplete,
				},
			},
		},
		{
			Name:        "docker",
			Usage:       "Build a Dockerfile without converting to an Earthfile *experimental*",
			Description: "Builds a dockerfile",
			Hidden:      true, // Experimental.
			Action:      app.actionDocker,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dockerfile",
					Usage:       "Path to dockerfile input, or - for stdin",
					Value:       "Dockerfile",
					Destination: &app.dockerfilePath,
				},
				&cli.StringFlag{
					Name:        "tag",
					Usage:       "Name and tag for the built image; formatted as 'name:tag'",
					Destination: &app.earthfileFinalImage,
				},
			},
		},
		{
			Name:        "docker2earthly",
			Usage:       "Convert a Dockerfile into Earthfile",
			Description: "Converts an existing dockerfile into an Earthfile",
			Hidden:      true, // Experimental.
			Action:      app.actionDocker2Earthly,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dockerfile",
					Usage:       "Path to dockerfile input, or - for stdin",
					Value:       "Dockerfile",
					Destination: &app.dockerfilePath,
				},
				&cli.StringFlag{
					Name:        "earthfile",
					Usage:       "Path to Earthfile output, or - for stdout",
					Value:       "Earthfile",
					Destination: &app.earthfilePath,
				},
				&cli.StringFlag{
					Name:        "tag",
					Usage:       "Name and tag for the built image; formatted as 'name:tag'",
					Destination: &app.earthfileFinalImage,
				},
			},
		},
		{
			Name:  "org",
			Usage: "Earthly organization administration *experimental*",
			Subcommands: []*cli.Command{
				{
					Name:      "create",
					Usage:     "Create a new organization",
					UsageText: "earthly [options] org create <org-name>",
					Action:    app.actionOrgCreate,
				},
				{
					Name:      "list",
					Usage:     "List organizations you belong to",
					UsageText: "earthly [options] org list",
					Action:    app.actionOrgList,
				},
				{
					Name:      "list-permissions",
					Usage:     "List permissions and membership of an organization",
					UsageText: "earthly [options] org list-permissions <org-name>",
					Action:    app.actionOrgListPermissions,
				},
				{
					Name:      "invite",
					Usage:     "Invite accounts to your organization",
					UsageText: "earthly [options] org invite [options] <path> <email> [<email> ...]",
					Action:    app.actionOrgInvite,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "write",
							Usage:       "Grant write permissions in addition to read",
							Destination: &app.writePermission,
						},
					},
				},
				{
					Name:      "revoke",
					Usage:     "Remove accounts from your organization",
					UsageText: "earthly [options] org revoke <path> <email> [<email> ...]",
					Action:    app.actionOrgRevoke,
				},
			},
		},
		{
			Name:      "ls",
			Usage:     "List targets from an Earthfile *experimental*",
			UsageText: "earthly [options] ls [<project-ref>]",
			Action:    app.actionListTargets,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "args",
					Aliases:     []string{"a"},
					Usage:       "Show Arguments",
					Destination: &app.lsShowArgs,
				},
				&cli.BoolFlag{
					Name:        "long",
					Aliases:     []string{"l"},
					Usage:       "Show full target-ref",
					Destination: &app.lsShowLong,
				},
			},
		},
		{
			Name:        "secrets",
			Usage:       "Earthly secrets",
			Description: "Manage cloud secrets *experimental*",
			Subcommands: []*cli.Command{
				{
					Name:  "set",
					Usage: "Stores a secret in the secrets store",
					UsageText: "earthly [options] secrets set <path> <value>\n" +
						"   earthly [options] secrets set --file <local-path> <path>\n" +
						"   earthly [options] secrets set --file <local-path> <path>",
					Action: app.actionSecretsSet,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "file",
							Aliases:     []string{"f"},
							Usage:       "Stores secret stored in file",
							Destination: &app.secretFile,
						},
						&cli.BoolFlag{
							Name:        "stdin",
							Aliases:     []string{"i"},
							Usage:       "Stores secret read from stdin",
							Destination: &app.secretStdin,
						},
					},
				},
				{
					Name:      "get",
					Action:    app.actionSecretsGet,
					Usage:     "Retrieve a secret from the secrets store",
					UsageText: "earthly [options] secrets get [options] <path>",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Aliases:     []string{"n"},
							Usage:       "Disable newline at the end of the secret",
							Destination: &app.disableNewLine,
						},
					},
				},
				{
					Name:      "ls",
					Usage:     "List secrets in the secrets store",
					UsageText: "earthly [options] secrets ls [<path>]",
					Action:    app.actionSecretsList,
				},
				{
					Name:      "rm",
					Usage:     "Removes a secret from the secrets store",
					UsageText: "earthly [options] secrets rm <path>",
					Action:    app.actionSecretsRemove,
				},
			},
		},
		{
			Name:  "account",
			Usage: "Create or manage an Earthly account *experimental*",
			Subcommands: []*cli.Command{
				{
					Name:        "register",
					Usage:       "Register for an Earthly account",
					Description: "Register for an Earthly account",
					UsageText: "You may register using GitHub OAuth, by visiting https://ci.earthly.dev\n" +
						"   Once authenticated, a login token will be displayed which can be used to login:\n" +
						"\n" +
						"       earthly [options] account login --token <token>\n" +
						"\n" +
						"   Alternatively, you can register using an email:\n" +
						"       first, request a token with:\n" +
						"\n" +
						"           earthly [options] account register --email <email>\n" +
						"\n" +
						"       then check your email to retrieve the token, then continue by running:\n" +
						"\n" +
						"           earthly [options] account register --token <token>\n",
					Action: app.actionRegister,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "email",
							Usage:       "Email address",
							Destination: &app.email,
						},
						&cli.StringFlag{
							Name:        "token",
							Usage:       "Email verification token",
							Destination: &app.token,
						},
						&cli.StringFlag{
							Name:        "password",
							EnvVars:     []string{"EARTHLY_PASSWORD"},
							Usage:       "Specify password on the command line instead of interactively being asked",
							Destination: &app.password,
						},
						&cli.StringFlag{
							Name:        "public-key",
							EnvVars:     []string{"EARTHLY_PUBLIC_KEY"},
							Usage:       "Path to public key to register",
							Destination: &app.registrationPublicKey,
						},
						&cli.BoolFlag{
							Name:        "accept-terms-of-service-privacy",
							EnvVars:     []string{"EARTHLY_ACCEPT_TERMS_OF_SERVICE_PRIVACY"},
							Usage:       "Accept the Terms & Conditions, and Privacy Policy",
							Destination: &app.termsConditionsPrivacy,
						},
					},
				},
				{
					Name:        "login",
					Usage:       "Login to an Earthly account",
					Description: "Login to an Earthly account",
					UsageText: "earthly [options] account login\n" +
						"   earthly [options] account login --email <email>\n" +
						"   earthly [options] account login --email <email> --password <password>\n" +
						"   earthly [options] account login --token <token>\n",
					Action: app.actionAccountLogin,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "email",
							Usage:       "Email address",
							Destination: &app.email,
						},
						&cli.StringFlag{
							Name:        "token",
							Usage:       "Authentication token",
							Destination: &app.token,
						},
						&cli.StringFlag{
							Name:        "password",
							EnvVars:     []string{"EARTHLY_PASSWORD"},
							Usage:       "Specify password on the command line instead of interactively being asked",
							Destination: &app.password,
						},
					},
				},
				{
					Name:        "logout",
					Usage:       "Logout of an Earthly account",
					Description: "Logout of an Earthly account; this has no effect for ssh-based authentication",
					Action:      app.actionAccountLogout,
				},
				{
					Name:      "list-keys",
					Usage:     "List associated public keys used for authentication",
					UsageText: "earthly [options] account list-keys",
					Action:    app.actionAccountListKeys,
				},
				{
					Name:      "add-key",
					Usage:     "Associate a new public key with your account",
					UsageText: "earthly [options] add-key [<key>]",
					Action:    app.actionAccountAddKey,
				},
				{
					Name:      "remove-key",
					Usage:     "Removes an existing public key from your account",
					UsageText: "earthly [options] remove-key <key>",
					Action:    app.actionAccountRemoveKey,
				},
				{
					Name:      "list-tokens",
					Usage:     "List associated tokens used for authentication",
					UsageText: "earthly [options] account list-tokens",
					Action:    app.actionAccountListTokens,
				},
				{
					Name:      "create-token",
					Usage:     "Create a new authentication token for your account",
					UsageText: "earthly [options] account create-token [options] <token name>",
					Action:    app.actionAccountCreateToken,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "write",
							Usage:       "Grant write permissions in addition to read",
							Destination: &app.writePermission,
						},
						&cli.StringFlag{
							Name:        "expiry",
							Usage:       "Set token expiry date in the form YYYY-MM-DD or never (default 1year)",
							Destination: &app.expiry,
						},
					},
				},
				{
					Name:      "remove-token",
					Usage:     "Remove an authentication token from your account",
					UsageText: "earthly [options] account remove-token <token>",
					Action:    app.actionAccountRemoveToken,
				},
			},
		},
		{
			Name:        "debug",
			Usage:       "Print debug information about an Earthfile",
			Description: "Print debug information about an Earthfile",
			ArgsUsage:   "[<path>]",
			Hidden:      true, // Dev purposes only.
			Subcommands: []*cli.Command{
				{
					Name:      "ast",
					Usage:     "Output the AST",
					UsageText: "earthly [options] debug ast",
					Action:    app.actionDebugAst,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "source-map",
							Usage:       "Enable outputting inline sourcemap",
							Destination: &app.enableSourceMap,
						},
					},
				},
			},
		},
		{
			Name:        "prune",
			Usage:       "Prune Earthly build cache",
			Description: "Prune Earthly build cache",
			Action:      app.actionPrune,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "all",
					Aliases:     []string{"a"},
					EnvVars:     []string{"EARTHLY_PRUNE_ALL"},
					Usage:       "Prune all cache",
					Destination: &app.pruneAll,
				},
				&cli.BoolFlag{
					Name:        "reset",
					EnvVars:     []string{"EARTHLY_PRUNE_RESET"},
					Usage:       "Reset cache entirely by wiping cache dir",
					Destination: &app.pruneReset,
				},
			},
		},
		{
			Name:   "config",
			Usage:  "Edits your Earthly configuration file",
			Action: app.actionConfig,
			UsageText: `This command takes a path, and a value and sets it in your configuration file.

	 As the configuration file is YAML, the key must be a valid key within the file. You can specify sub-keys by using "." to separate levels.
	 If the sub-key you wish to use has a "." in it, you can quote that subsection, like this: git."github.com".

	 Values must be valid YAML, and also be deserializable into the key you wish to assign them to.
	 This means you can set higher level objects using a compact style, or single values with simple values.

	 Only one key/value can be set per invocation.

	 To get help with a specific key, do "config [key] --help". Or, visit https://docs.earthly.dev/earthly-config for more details.`,
			Description: `Set your cache size:

	config global.cache_size_mb 1234

Set additional buildkit args, using a YAML array:

	config global.buildkit_additional_args '["userns", "--host"]'

Set a key containing a period:

	config 'git."example.com".password' hunter2

Set up a whole custom git repository for a server called example.com, using a single-line YAML literal:
	* which stores git repos under /var/git/repos/name-of-repo.git
	* allows access over ssh
	* using port 2222
	* sets the username to git
	* is recognized to earthly as example.com/name-of-repo

	config git "{example: {pattern: 'example.com/([^/]+)', substitute: 'ssh://git@example.com:2222/var/git/repos/\$1.git', auth: ssh}}"
			`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "dry-run",
					Usage:       "Print the changed config file to the console instead of writing it out",
					Destination: &app.configDryRun,
				},
			},
		},
		{
			Name: "satellite",
			Usage: "Launch and use a Satellite runner as remote backend for Earthly builds.\n" +
				"	Satellites can be used to optimize and share cache between multiple builds and users,\n" +
				"	as well as run builds in native architectures independent of where the Earthly client is invoked.\n" +
				"	Note: this feature is currently experimental.\n" +
				"	If you'd like to try it out, please contact us via Slack to be added to the beta testers group.",
			UsageText:   "earthly satellite (launch|list|destroy|unselect)",
			Description: "Create and manage Earthly build Satellites",
			Subcommands: []*cli.Command{
				{
					Name:        "launch",
					Description: "Launch a new Earthly Satellite",
					UsageText: "earthly satellite launch <satellite-name>\n" +
						"	earthly satellite launch --org <organization-name> <satellite-name>",
					Action: app.actionSatelliteLaunch,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the organization the satellite belongs to. Required when user is a member of multiple.",
							Required:    false,
							Destination: &app.satelliteOrg,
						},
					},
				},
				{
					Name:        "destroy",
					Description: "Destroy an Earthly Satellite",
					UsageText: "earthly satellite destroy <satellite-name>\n" +
						"	earthly satellite destroy --org <organization-name> <satellite-name>",
					Action: app.actionSatelliteDestroy,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the organization the satellite belongs to. Required when user is a member of multiple.",
							Required:    false,
							Destination: &app.satelliteOrg,
						},
					},
				},
				{
					Name:        "list",
					Description: "List your Earthly Satellites",
					UsageText: "earthly satellite list\n" +
						"	earthly satellite list --org <organization-name>",
					Action: app.actionSatelliteList,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the organization the satellite belongs to. Required when user is a member of multiple.",
							Required:    false,
							Destination: &app.satelliteOrg,
						},
					},
				},
				{
					Name:        "describe",
					Description: "Show additional details about a Satellite instance",
					UsageText: "earthly satellite describe <satellite-name>\n" +
						"	earthly satellite list --org <organization-name> <satellite-name>",
					Action: app.actionSatelliteDescribe,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the organization the satellite belongs to. Required when user is a member of multiple.",
							Required:    false,
							Destination: &app.satelliteOrg,
						},
					},
				},
				{
					Name:        "select",
					Description: "Choose which satellite to use to build your app.",
					UsageText: "earthly satellite select <satellite-name>\n" +
						"	earthly satellite select --org <organization-name> <satellite-name>",
					Action: app.actionSatelliteSelect,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the organization the satellite belongs to. Required when user is a member of multiple.",
							Required:    false,
							Destination: &app.satelliteOrg,
						},
					},
				},
				{
					Name:        "unselect",
					Description: "Remove any currently selected Satellite instance from your Earthly configuration.",
					UsageText:   "earthly satellite unselect",
					Action:      app.actionSatelliteUnselect,
				},
			},
		},
	}

	app.cliApp.Before = app.before
	return app
}

func wrap(s ...string) string {
	return strings.Join(s, "\n\t")
}

func (app *earthlyApp) before(context *cli.Context) error {
	if app.enableProfiler {
		go profhandler()
	}

	if app.verbose {
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

	feConfig := &containerutil.FrontendConfig{
		BuildkitHostCLIValue:       app.buildkitHost,
		BuildkitHostFileValue:      app.cfg.Global.BuildkitHost,
		DebuggerHostCLIValue:       app.debuggerHost,
		DebuggerHostFileValue:      app.cfg.Global.DebuggerHost,
		DebuggerPortFileValue:      app.cfg.Global.DebuggerPort,
		LocalRegistryHostFileValue: app.cfg.Global.LocalRegistryHost,
		Console:                    app.console,
	}
	fe, err := containerutil.FrontendForSetting(context.Context, app.cfg.Global.ContainerFrontend, feConfig)
	if err != nil {
		origErr := err
		fe, err = containerutil.NewStubFrontend(context.Context, feConfig)
		if err != nil {
			return errors.Wrap(err, "failed frontend initialization")
		}
		app.console.Warnf("%s frontend initialization failed due to %s; but will try anyway", app.cfg.Global.ContainerFrontend, origErr.Error())
	}
	app.containerFrontend = fe

	// command line option overrides the config which overrides the default value
	if !context.IsSet("buildkit-image") && app.cfg.Global.BuildkitImage != "" {
		app.buildkitdImage = app.cfg.Global.BuildkitImage
	}

	// These URLs were calculated relative to the configured frontend. In the case of an automatically detected frontend,
	// they are calculated according to the first selected one in order of precedence.
	buildkitURLs := fe.Config().FrontendURLs
	app.buildkitHost = buildkitURLs.BuildkitHost.String()
	app.debuggerHost = buildkitURLs.DebuggerHost.String()
	app.localRegistryHost = buildkitURLs.LocalRegistryHost.String()

	bkURL, err := url.Parse(app.buildkitHost) // Not validated because we already did that when we calculated it.
	if err != nil {
		return errors.Wrap(err, "failed to parse generated buildkit URL")
	}

	if bkURL.Scheme == "tcp" {
		app.handleTLSCertificateSettings(context)
	}

	app.buildkitdSettings.AdditionalArgs = app.cfg.Global.BuildkitAdditionalArgs
	app.buildkitdSettings.AdditionalConfig = app.cfg.Global.BuildkitAdditionalConfig
	app.buildkitdSettings.Timeout = time.Duration(app.cfg.Global.BuildkitRestartTimeoutS) * time.Second
	app.buildkitdSettings.Debug = app.debug
	app.buildkitdSettings.BuildkitAddress = app.buildkitHost
	app.buildkitdSettings.DebuggerAddress = app.debuggerHost
	app.buildkitdSettings.LocalRegistryAddress = app.localRegistryHost
	app.buildkitdSettings.UseTCP = bkURL.Scheme == "tcp"
	app.buildkitdSettings.UseTLS = app.cfg.Global.TLSEnabled
	app.buildkitdSettings.MaxParallelism = app.cfg.Global.BuildkitMaxParallelism
	app.buildkitdSettings.CacheSizeMb = app.cfg.Global.BuildkitCacheSizeMb
	app.buildkitdSettings.CacheSizePct = app.cfg.Global.BuildkitCacheSizePct

	// ensure the MTU is something allowable in IPv4, cap enforced by type. Zero is autodetect.
	if app.cfg.Global.CniMtu != 0 && app.cfg.Global.CniMtu < 68 {
		return errors.New("invalid overridden MTU size")
	}
	app.buildkitdSettings.CniMtu = app.cfg.Global.CniMtu

	if app.cfg.Global.IPTables != "" && app.cfg.Global.IPTables != "iptables-legacy" && app.cfg.Global.IPTables != "iptables-nft" {
		return errors.New(`invalid overridden iptables name. Valid values are "iptables-legacy" or "iptables-nft"`)
	}
	app.buildkitdSettings.IPTables = app.cfg.Global.IPTables

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

func (app *earthlyApp) configureSatellite(cc cloud.Client) error {
	if !app.isUsingSatellite() || cc == nil {
		// If the app is not using a cloud client, or the command doesn't interact with the cloud (prune, bootstrap)
		// then pretend its all good and use your regular configuration.
		return nil
	}

	app.console.Printf("Using Satellite: %s", app.satelliteName)

	// When using a satellite, interactive and local do not work; as they are not SSL nor routable yet.
	app.console.Warnf("Note: the Interactive Debugger, Interactive RUN commands, and Local Registries do not yet work on Earthly Satellites.")

	// Set up extra settings needed for buildkit RPC metadata
	if app.cfg.Satellite.Name != "" {
		app.buildkitdSettings.SatelliteName = app.cfg.Satellite.Name
		app.buildkitdSettings.SatelliteOrg = app.cfg.Satellite.Org
		if app.satelliteAddress != "" {
			app.buildkitdSettings.BuildkitAddress = app.satelliteAddress
		} else {
			app.buildkitdSettings.BuildkitAddress = containerutil.SatelliteAddress
		}
	}
	token, err := cc.GetAuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	app.buildkitdSettings.SatelliteToken = token

	// TODO (dchw) what other settings might we want to override here?

	return nil
}

func (app *earthlyApp) isUsingSatellite() bool {
	return len(app.cfg.Satellite.Name) > 0
}

func (app *earthlyApp) GetBuildkitClient(c *cli.Context, cc cloud.Client) (*client.Client, error) {
	err := app.configureSatellite(cc)
	if err != nil {
		return nil, errors.Wrapf(err, "could not construct new buildkit client")
	}

	return buildkitd.NewClient(c.Context, app.console, app.buildkitdImage, app.containerName, app.containerFrontend, app.buildkitdSettings)
}

func (app *earthlyApp) handleTLSCertificateSettings(context *cli.Context) {
	if !app.cfg.Global.TLSEnabled {
		return
	}

	app.buildkitdSettings.TLSCA = app.cfg.Global.TLSCA

	if !context.IsSet("tlscert") && app.cfg.Global.ClientTLSCert != "" {
		app.certPath = app.cfg.Global.ClientTLSCert
	}

	if !context.IsSet("tlskey") && app.cfg.Global.ClientTLSKey != "" {
		app.keyPath = app.cfg.Global.ClientTLSKey
	}

	app.buildkitdSettings.ClientTLSCert = app.certPath
	app.buildkitdSettings.ClientTLSKey = app.keyPath

	app.buildkitdSettings.ServerTLSCert = app.cfg.Global.ServerTLSCert
	app.buildkitdSettings.ServerTLSKey = app.cfg.Global.ServerTLSKey
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

// to enable autocomplete, enter
// complete -o nospace -C "/path/to/earthly" earthly
func (app *earthlyApp) autoComplete(ctx context.Context) {
	_, found := os.LookupEnv("COMP_LINE")
	if !found {
		return
	}

	app.console = app.console.WithLogLevel(conslogging.Silent)

	err := app.autoCompleteImp(ctx)
	if err != nil {
		errToLog := err
		logDir, err := cliutil.GetOrCreateEarthlyDir()
		if err != nil {
			os.Exit(1)
		}
		logFile := filepath.Join(logDir, "autocomplete.log")
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			os.Exit(1)
		}
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			os.Exit(1)
		}
		fmt.Fprintf(f, "error during autocomplete: %s\n", errToLog)
		os.Exit(1)
	}
	os.Exit(0)
}

func (app *earthlyApp) autoCompleteImp(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("recovered panic in autocomplete %s: %s", r, debug.Stack())
		}
	}()

	compLine := os.Getenv("COMP_LINE")   // full command line
	compPoint := os.Getenv("COMP_POINT") // where the cursor is

	compPointInt, err := strconv.ParseUint(compPoint, 10, 64)
	if err != nil {
		return err
	}

	gitLookup := buildcontext.NewGitLookup(app.console, app.sshAuthSock)
	resolver := buildcontext.NewResolver("", nil, gitLookup, app.console, "")
	var gwClient gwclient.Client // TODO this is a nil pointer which causes a panic if we try to expand a remotely referenced earthfile
	// it's expensive to create this gwclient, so we need to implement a lazy eval which returns it when required.

	potentials, err := autocomplete.GetPotentials(ctx, resolver, gwClient, compLine, int(compPointInt), app.cliApp)
	if err != nil {
		return err
	}
	for _, p := range potentials {
		fmt.Printf("%s\n", p)
	}

	return err
}
func (app *earthlyApp) insertBashCompleteEntry() error {
	var path string
	if runtime.GOOS == "darwin" {
		path = "/usr/local/etc/bash_completion.d/earthly"
	} else {
		path = "/usr/share/bash-completion/completions/earthly"
	}
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", dirPath)
	}
	if !dirPathExists {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s does not exist\n", dirPath)
		return nil // bash-completion isn't available, silently fail.
	}

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", path)
	}
	if pathExists {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bashEntry, err := bashCompleteEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s\n", err)
		return nil // bash-completion isn't available, silently fail.
	}

	_, err = f.Write([]byte(bashEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}
	return nil
}

func (app *earthlyApp) deleteZcompdump() error {
	var homeDir string
	sudoUser, found := os.LookupEnv("SUDO_USER")
	if !found {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return errors.Wrapf(err, "failed to lookup current user home dir")
		}
	} else {
		currentUser, err := user.Lookup(sudoUser)
		if err != nil {
			return errors.Wrapf(err, "failed to lookup user %s", sudoUser)
		}
		homeDir = currentUser.HomeDir
	}
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read dir %s", homeDir)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".zcompdump") {
			path := filepath.Join(homeDir, f.Name())
			err := os.Remove(path)
			if err != nil {
				return errors.Wrapf(err, "failed to remove %s", path)
			}
		}
	}
	return nil
}

func bashCompleteEntry() (string, error) {
	template := "complete -o nospace -C '__earthly__' earthly\n"
	return renderEntryTemplate(template)
}

func zshCompleteEntry() (string, error) {
	template := `#compdef _earthly earthly

function _earthly {
    autoload -Uz bashcompinit
    bashcompinit
    complete -o nospace -C '__earthly__' earthly
}
`
	return renderEntryTemplate(template)
}

func renderEntryTemplate(template string) (string, error) {
	earthlyPath, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to determine earthly path: %s", err)
	}
	return strings.ReplaceAll(template, "__earthly__", earthlyPath), nil
}

// If debugging this, it might be required to run `rm ~/.zcompdump*` to remove the cache
func (app *earthlyApp) insertZSHCompleteEntry() error {
	// should be the same on linux and macOS
	path := "/usr/local/share/zsh/site-functions/_earthly"
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", dirPath)
	}
	if !dirPathExists {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: %s does not exist\n", dirPath)
		return nil // zsh-completion isn't available, silently fail.
	}

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", path)
	}
	if pathExists {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	compEntry, err := zshCompleteEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: %s\n", err)
		return nil // zsh-completion isn't available, silently fail.
	}

	_, err = f.Write([]byte(compEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}

	return app.deleteZcompdump()
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

func isEarthlyBinary(path string) bool {
	// apply heuristics to see if binary is a version of earthly
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	if !bytes.Contains(data, []byte("docs.earthly.dev")) {
		return false
	}
	if !bytes.Contains(data, []byte("api.earthly.dev")) {
		return false
	}
	if !bytes.Contains(data, []byte("Earthfile")) {
		return false
	}
	return true
}

func symlinkEarthlyToEarth() error {
	binPath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to get current executable path")
	}

	baseName := path.Base(binPath)
	if baseName != "earthly" {
		return nil
	}

	earthPath := path.Join(path.Dir(binPath), "earth")

	earthPathExists, err := fileutil.FileExists(earthPath)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", earthPath)
	}
	if !earthPathExists && termutil.IsTTY() {
		return nil // legacy earth binary doesn't exist, don't create it (unless we're under a non-tty system e.g. CI)
	}

	if !isEarthlyBinary(earthPath) {
		return nil // file exists but is not an earthly binary, leave it alone.
	}

	// otherwise legacy earth command has been detected, remove it and symlink
	// to the new earthly command.
	err = os.Remove(earthPath)
	if err != nil {
		return errors.Wrapf(err, "failed to remove old install at %s", earthPath)
	}
	err = os.Symlink(binPath, earthPath)
	if err != nil {
		return errors.Wrapf(err, "failed to symlink %s to %s", binPath, earthPath)
	}
	return nil
}

func (app *earthlyApp) actionBootstrap(c *cli.Context) error {
	app.commandName = "bootstrap"

	switch app.homebrewSource {
	case "bash":
		compEntry, err := bashCompleteEntry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to enable bash-completion: %s\n", err)
			return nil // zsh-completion isn't available, silently fail.
		}
		fmt.Print(compEntry)
		return nil
	case "zsh":
		compEntry, err := zshCompleteEntry()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to bootstrap zsh-completion: %s\n", err)
			return nil // zsh-completion isn't available, silently fail.
		}
		fmt.Print(compEntry)
		return nil
	case "":
		break
	default:
		return errors.Errorf("unhandled source %q", app.homebrewSource)
	}

	return app.bootstrap(c)
}

func (app *earthlyApp) bootstrap(c *cli.Context) error {
	var err error
	console := app.console.WithPrefix("bootstrap")
	defer func() {
		// cliutil.IsBootstrapped() determines if bootstrapping was done based
		// on the existance of ~/.earthly; therefore we must ensure it's created.
		cliutil.GetOrCreateEarthlyDir()
		cliutil.EnsurePermissions()
	}()

	if app.bootstrapWithAutocomplete {
		// Because this requires sudo, it should warn and not fail the rest of it.
		err = app.insertBashCompleteEntry()
		if err != nil {
			console.Warnf("Warning: %s\n", err.Error())
			err = nil
		}
		err = app.insertZSHCompleteEntry()
		if err != nil {
			console.Warnf("Warning: %s\n", err.Error())
			err = nil
		}

		console.Printf("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	}

	err = symlinkEarthlyToEarth()
	if err != nil {
		console.Warnf("Warning: %s\n", err.Error())
		err = nil
	}

	if !app.bootstrapNoBuildkit && !app.isUsingSatellite() {
		bkURL, err := url.Parse(app.buildkitHost)
		if err != nil {
			return errors.Wrapf(err, "invalid buildkit_host: %s", app.cfg.Global.BuildkitHost)
		}
		if bkURL.Scheme == "tcp" && app.cfg.Global.TLSEnabled {
			root, err := cliutil.GetOrCreateEarthlyDir()
			if err != nil {
				return err
			}

			certsDir := filepath.Join(root, "certs")
			err = buildkitd.GenerateCertificates(certsDir)
			if err != nil {
				return errors.Wrap(err, "setup TLS")
			}
		}

		// Bootstrap buildkit - pulls image and starts daemon.
		bkClient, err := app.GetBuildkitClient(c, nil)
		if err != nil {
			return errors.Wrap(err, "bootstrap new buildkitd client")
		}
		defer bkClient.Close()
	}

	console.Printf("Bootstrapping successful.\n")
	return nil
}

func promptInput(question string) string {
	fmt.Printf("%s", question)
	rbuf := bufio.NewReader(os.Stdin)
	line, err := rbuf.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimRight(line, "\n")
}

func (app *earthlyApp) actionOrgCreate(c *cli.Context) error {
	app.commandName = "orgCreate"
	if c.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	org := c.Args().Get(0)
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cc.CreateOrg(org)
	if err != nil {
		return errors.Wrap(err, "failed to create org")
	}
	return nil
}

func (app *earthlyApp) actionOrgList(c *cli.Context) error {
	app.commandName = "orgList"
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	orgs, err := cc.ListOrgs()
	if err != nil {
		return errors.Wrap(err, "failed to list orgs")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, org := range orgs {
		fmt.Fprintf(w, "/%s/", org.Name)
		if org.Admin {
			fmt.Fprintf(w, "\tadmin")
		} else {
			fmt.Fprintf(w, "\tmember")
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionOrgListPermissions(c *cli.Context) error {
	app.commandName = "orgListPermissions"
	if c.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := c.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	orgs, err := cc.ListOrgPermissions(path)
	if err != nil {
		return errors.Wrap(err, "failed to list org permissions")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, org := range orgs {
		fmt.Fprintf(w, "%s\t%s", org.Path, org.User)
		if org.Write {
			fmt.Fprintf(w, "\trw")
		} else {
			fmt.Fprintf(w, "\tr")
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}

func (app *earthlyApp) actionOrgInvite(c *cli.Context) error {
	app.commandName = "orgInvite"
	if c.NArg() < 2 {
		return errors.New("invalid number of arguments provided")
	}
	path := c.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		return errors.New("invitation paths must end with a slash (/)")
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	userEmail := c.Args().Get(1)
	err = cc.Invite(path, userEmail, app.writePermission)
	if err != nil {
		return errors.Wrap(err, "failed to invite user into org")
	}
	return nil
}

func (app *earthlyApp) actionOrgRevoke(c *cli.Context) error {
	app.commandName = "orgRevoke"
	if c.NArg() < 2 {
		return errors.New("invalid number of arguments provided")
	}
	path := c.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		return errors.New("revoked paths must end with a slash (/)")
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	userEmail := c.Args().Get(1)
	err = cc.RevokePermission(path, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to revoke user from org")
	}
	return nil
}

func (app *earthlyApp) actionSecretsList(c *cli.Context) error {
	app.commandName = "secretsList"

	path := "/"
	if c.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	} else if c.NArg() == 1 {
		path = c.Args().Get(0)
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	paths, err := cc.List(path)
	if err != nil {
		return errors.Wrap(err, "failed to list secret")
	}
	for _, path := range paths {
		fmt.Println(path)
	}
	return nil
}

func (app *earthlyApp) actionSecretsGet(c *cli.Context) error {
	app.commandName = "secretsGet"
	if c.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := c.Args().Get(0)
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	data, err := cc.Get(path)
	if err != nil {
		return errors.Wrap(err, "failed to get secret")
	}
	fmt.Printf("%s", data)
	if !app.disableNewLine {
		fmt.Printf("\n")
	}
	return nil
}

func (app *earthlyApp) actionSecretsRemove(c *cli.Context) error {
	app.commandName = "secretsRemove"
	if c.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := c.Args().Get(0)
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cc.Remove(path)
	if err != nil {
		return errors.Wrap(err, "failed to remove secret")
	}
	return nil
}

func (app *earthlyApp) actionSecretsSet(c *cli.Context) error {
	app.commandName = "secretsSet"
	var path string
	var value string
	if app.secretFile == "" && !app.secretStdin {
		if c.NArg() != 2 {
			return errors.New("invalid number of arguments provided")
		}
		path = c.Args().Get(0)
		value = c.Args().Get(1)
	} else if app.secretStdin {
		if app.secretFile != "" {
			return errors.New("only one of --file or --stdin can be used at a time")
		}
		if c.NArg() != 1 {
			return errors.New("invalid number of arguments provided")
		}
		path = c.Args().Get(0)
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read from stdin")
		}
		value = string(data)
	} else {
		if c.NArg() != 1 {
			return errors.New("invalid number of arguments provided")
		}
		path = c.Args().Get(0)
		data, err := os.ReadFile(app.secretFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read secret from %s", app.secretFile)
		}
		value = string(data)
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cc.Set(path, []byte(value))
	if err != nil {
		return errors.Wrap(err, "failed to set secret")
	}
	return nil
}

func (app *earthlyApp) actionRegister(c *cli.Context) error {
	app.commandName = "secretsRegister"
	if app.email == "" {
		return errors.New("no email given")
	}

	if !strings.Contains(app.email, "@") {
		return errors.New("email is invalid")
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	if app.token == "" {
		err := cc.RegisterEmail(app.email)
		if err != nil {
			return errors.Wrap(err, "failed to register email")
		}
		fmt.Printf("An email has been sent to %q containing a registration token\n", app.email)
		return nil
	}

	var publicKeys []*agent.Key
	if app.sshAuthSock != "" {
		var err error
		publicKeys, err = cc.GetPublicKeys()
		if err != nil {
			return err
		}
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass app.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	pword := app.password
	if app.password == "" {
		fmt.Printf("pick a password: ")
		enteredPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println("")
		fmt.Printf("confirm password: ")
		enteredPassword2, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println("")
		if string(enteredPassword) != string(enteredPassword2) {
			return errors.Errorf("passwords do not match")
		}
		pword = string(enteredPassword)
	}

	var interactiveAccept bool
	if !app.termsConditionsPrivacy {
		rawAccept := promptInput("I acknowledge Earthly Technologies’ Privacy Policy (https://earthly.dev/privacy-policy) and agree to Earthly Technologies Terms of Service (https://earthly.dev/tos) [y/N]: ")
		if rawAccept == "" {
			rawAccept = "n"
		}
		accept := strings.ToLower(rawAccept)[0]

		interactiveAccept = accept == 'y'
	}
	termsConditionsPrivacy := app.termsConditionsPrivacy || interactiveAccept

	var publicKey string
	if app.registrationPublicKey == "" {
		if len(publicKeys) > 0 {
			fmt.Printf("Which of the following keys do you want to register?\n")
			fmt.Printf("0) none\n")
			for i, key := range publicKeys {
				fmt.Printf("%d) %s\n", i+1, key.String())
			}
			keyNum := promptInput("enter key number (1=default): ")
			if keyNum == "" {
				keyNum = "1"
			}
			i, err := strconv.Atoi(keyNum)
			if err != nil {
				return errors.Wrap(err, "invalid key number")
			}
			if i < 0 || i > len(publicKeys) {
				return errors.Errorf("invalid key number")
			}
			if i > 0 {
				publicKey = publicKeys[i-1].String()
			}
		}
	} else {
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(app.registrationPublicKey))
		if err == nil {
			// supplied public key is valid
			publicKey = app.registrationPublicKey
		} else {
			// otherwise see if it matches the name (Comment) of a key known by the ssh agent
			for _, key := range publicKeys {
				if key.Comment == app.registrationPublicKey {
					publicKey = key.String()
					break
				}
			}
			if publicKey == "" {
				return errors.Errorf("failed to find key in ssh agent's known keys")
			}
		}
	}

	err = cc.CreateAccount(app.email, app.token, pword, publicKey, termsConditionsPrivacy)
	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	fmt.Println("Account registration complete")
	return nil
}

func (app *earthlyApp) actionAccountListKeys(c *cli.Context) error {
	app.commandName = "accountListKeys"
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	keys, err := cc.ListPublicKeys()
	if err != nil {
		return errors.Wrap(err, "failed to list account keys")
	}
	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
	return nil
}

func (app *earthlyApp) actionAccountAddKey(c *cli.Context) error {
	app.commandName = "accountAddKey"
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	if c.NArg() > 1 {
		for _, k := range c.Args().Slice() {
			err := cc.AddPublickKey(k)
			if err != nil {
				return errors.Wrap(err, "failed to add public key to account")
			}
		}
		return nil
	}

	publicKeys, err := cc.GetPublicKeys()
	if err != nil {
		return err
	}
	if len(publicKeys) == 0 {
		return errors.Errorf("unable to list available public keys, is ssh-agent running?; do you need to run ssh-add?")
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass app.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Which of the following keys do you want to register?\n")
	for i, key := range publicKeys {
		fmt.Printf("%d) %s\n", i+1, key.String())
	}
	keyNum := promptInput("enter key number (1=default): ")
	if keyNum == "" {
		keyNum = "1"
	}
	i, err := strconv.Atoi(keyNum)
	if err != nil {
		return errors.Wrap(err, "invalid key number")
	}
	if i <= 0 || i > len(publicKeys) {
		return errors.Errorf("invalid key number")
	}
	publicKey := publicKeys[i-1].String()

	err = cc.AddPublickKey(publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to add public key to account")
	}

	// switch over to new key if the user is currently using password-based auth
	email, authType, _, err := cc.WhoAmI()
	if err != nil {
		return errors.Wrap(err, "failed to validate auth token")
	}
	if authType == "password" {
		err = cc.SetSSHCredentials(email, publicKey)
		if err != nil {
			app.console.Warnf("failed to authenticate using newly added public key: %s", err.Error())
			return nil
		}
		fmt.Printf("Switching from password-based login to ssh-based login\n")
	}

	return nil
}

func (app *earthlyApp) actionAccountRemoveKey(c *cli.Context) error {
	app.commandName = "accountRemoveKey"
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	for _, k := range c.Args().Slice() {
		err := cc.RemovePublickKey(k)
		if err != nil {
			return errors.Wrap(err, "failed to add public key to account")
		}
	}
	return nil
}
func (app *earthlyApp) actionAccountListTokens(c *cli.Context) error {
	app.commandName = "accountListTokens"
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	tokens, err := cc.ListTokens()
	if err != nil {
		return errors.Wrap(err, "failed to list account tokens")
	}
	if len(tokens) == 0 {
		return nil // avoid printing header columns when there are no tokens
	}

	now := time.Now()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Token Name\tRead/Write\tExpiry\n")
	for _, token := range tokens {
		expired := now.After(token.Expiry)
		fmt.Fprintf(w, "%s", token.Name)
		if token.Write {
			fmt.Fprintf(w, "\trw")
		} else {
			fmt.Fprintf(w, "\tr")
		}
		fmt.Fprintf(w, "\t%s UTC", token.Expiry.UTC().Format("2006-01-02T15:04"))
		if expired {
			fmt.Fprintf(w, " *expired*")
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}
func (app *earthlyApp) actionAccountCreateToken(c *cli.Context) error {
	app.commandName = "accountCreateToken"
	if c.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	var expiry time.Time
	if app.expiry == "" {
		expiry = time.Now().Add(time.Hour * 24 * 365)
	} else if app.expiry == "never" {
		expiry = time.Now().Add(time.Hour * 24 * 365 * 100) // TODO save this some other way
	} else {
		layouts := []string{
			"2006-01-02",
			time.RFC3339,
		}

		var err error
		for _, layout := range layouts {
			expiry, err = time.Parse(layout, app.expiry)
			if err == nil {
				break
			}
		}
		if err != nil {
			return errors.Errorf("failed to parse expiry %q", app.expiry)
		}
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	name := c.Args().First()
	token, err := cc.CreateToken(name, app.writePermission, &expiry)
	if err != nil {
		return errors.Wrap(err, "failed to create token")
	}
	expiryStr := humanize.Time(expiry)
	fmt.Printf("created token %q which will expire in %s; save this token somewhere, it can't be viewed again (only reset)\n", token, expiryStr)
	return nil
}
func (app *earthlyApp) actionAccountRemoveToken(c *cli.Context) error {
	app.commandName = "accountRemoveToken"
	if c.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	name := c.Args().First()
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cc.RemoveToken(name)
	if err != nil {
		return errors.Wrap(err, "failed to remove account tokens")
	}
	return nil
}

func (app *earthlyApp) actionAccountLogin(c *cli.Context) error {
	app.commandName = "accountLogin"
	email := app.email
	token := app.token
	pass := app.password

	if c.NArg() == 1 {
		emailOrToken := c.Args().First()
		if token == "" && email == "" {
			if cloud.IsValidEmail(emailOrToken) {
				email = emailOrToken
			} else {
				token = emailOrToken
			}

		} else {
			return errors.New("invalid number of arguments provided")
		}
	} else if c.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}

	if token != "" && (email != "" || pass != "") {
		return errors.New("--token cannot be used in conjuction with --email or --password")
	}
	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	// special case where global auth token overrides login logic
	if app.authToken != "" {
		if email != "" || token != "" || pass != "" {
			return errLoginFlagsHaveNoEffect
		}
		loggedInEmail, authType, writeAccess, err := cc.WhoAmI()
		if err != nil {
			return errors.Wrap(err, "failed to validate auth token")
		}
		if !writeAccess {
			authType = "read-only-" + authType
		}
		app.console.Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		app.printLogSharingMessage()
		return nil
	}

	if err = cc.DeleteCachedToken(); err != nil {
		return err
	}

	if token != "" || pass != "" {
		err := cc.DeleteAuthCache()
		if err != nil {
			return errors.Wrap(err, "failed to clear cached credentials")
		}
		cc.DisableSSHKeyGuessing()
	} else if email != "" {
		if err = cc.FindSSHCredentials(email); err == nil {
			// if err is not nil, we will try again below via cc.WhoAmI()

			if err = cc.Authenticate(); err != nil {
				return errors.Wrap(err, "authentication with cloud server failed")
			}
			app.console.Printf("Logged in as %q using ssh auth\n", email)
			app.printLogSharingMessage()
			return nil
		}
	}

	loggedInEmail, authType, writeAccess, err := cc.WhoAmI()
	switch errors.Cause(err) {
	case cloud.ErrUnauthorized:
		break
	case nil:
		if email != "" && email != loggedInEmail {
			break // case where a user has multiple emails and wants to switch
		}
		if !writeAccess {
			authType = "read-only-" + authType
		}
		app.console.Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		app.printLogSharingMessage()
		return nil
	default:
		return err
	}

	if email == "" && token == "" {
		if app.sshAuthSock == "" {
			app.console.Warnf("No ssh auth socket detected; falling back to password-based login\n")
		}

		emailOrToken := promptInput("enter your email or auth token: ")
		if strings.Contains(emailOrToken, "@") {
			email = emailOrToken
		} else {
			token = emailOrToken
		}
	}

	if email != "" && pass == "" {
		app.console.Printf("enter your password: \n")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		pass = string(passwordBytes)
		if pass == "" {
			return errors.Errorf("no password entered")
		}
	}

	if token != "" {
		email, err = cc.SetTokenCredentials(token)
		if err != nil {
			return err
		}
		app.console.Printf("Logged in as %q using token auth\n", email) // TODO display if using read-only token
		app.printLogSharingMessage()
	} else {
		err = cc.SetPasswordCredentials(email, string(pass))
		if err != nil {
			return err
		}
		app.console.Printf("Logged in as %q using password auth\n", email)
		app.console.Printf("Warning unencrypted password has been stored under ~/.earthly/auth.credentials; consider using ssh-based auth to prevent this.\n")
		app.printLogSharingMessage()
	}
	if err = cc.Authenticate(); err != nil {
		return errors.Wrap(err, "authentication with cloud server failed")
	}
	return nil
}

func (app *earthlyApp) printLogSharingMessage() {
	app.console.Printf("Log sharing is enabled by default. If you would like to disable it, run:\n" +
		"\n" +
		"\tearthly config global.disable_log_sharing true")
}

func (app *earthlyApp) actionAccountLogout(c *cli.Context) error {
	app.commandName = "accountLogout"

	if app.authToken != "" {
		return errLogoutHasNoEffectWhenAuthTokenSet
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return err
	}
	err = cc.DeleteAuthCache()
	if err != nil {
		return errors.Wrap(err, "failed to logout")
	}
	return nil
}

func (app *earthlyApp) actionDebugAst(c *cli.Context) error {
	app.commandName = "debugAst"
	if c.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := "./Earthfile"
	if c.NArg() == 1 {
		path = c.Args().First()
	}

	ef, err := ast.Parse(c.Context, path, app.enableSourceMap)
	if err != nil {
		return err
	}
	efDt, err := json.Marshal(ef)
	if err != nil {
		return errors.Wrap(err, "marshal ast")
	}
	fmt.Print(string(efDt))
	return nil
}

func (app *earthlyApp) actionPrune(c *cli.Context) error {
	app.commandName = "prune"
	if c.NArg() != 0 {
		return errors.New("invalid arguments")
	}
	if app.pruneReset {
		err := buildkitd.ResetCache(c.Context, app.console, app.buildkitdImage, app.containerName, app.containerFrontend, app.buildkitdSettings)
		if err != nil {
			return errors.Wrap(err, "reset cache")
		}
		return nil
	}

	if app.isUsingSatellite() {
		return errors.New("Cannot prune when using a satellite")
	}

	// Prune via API.
	bkClient, err := app.GetBuildkitClient(c, nil)
	if err != nil {
		return errors.Wrap(err, "prune new buildkitd client")
	}
	defer bkClient.Close()
	var opts []client.PruneOption
	if app.pruneAll {
		opts = append(opts, client.PruneAll)
	}
	ch := make(chan client.UsageInfo, 1)
	eg, ctx := errgroup.WithContext(c.Context)
	eg.Go(func() error {
		err = bkClient.Prune(ctx, ch, opts...)
		if err != nil {
			return errors.Wrap(err, "buildkit prune")
		}
		close(ch)
		return nil
	})
	eg.Go(func() error {
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return nil
				}
				// TODO: Print some progress info.
			case <-ctx.Done():
				return nil
			}
		}
	})
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "err group")
	}
	return nil
}

func (app *earthlyApp) actionDocker(c *cli.Context) error {
	app.commandName = "docker"

	dir := filepath.Dir(app.dockerfilePath)
	earthfilePath := filepath.Join(dir, "Earthfile")
	earthfilePathExists, err := fileutil.FileExists(earthfilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to check if %s exists", earthfilePath)
	}
	if earthfilePathExists {
		return errors.Errorf("earthfile already exists; please delete it if you wish to continue")
	}
	defer os.Remove(earthfilePath)

	err = docker2earthly.Docker2Earthly(app.dockerfilePath, earthfilePath, app.earthfileFinalImage)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Warning: earthly does not support all dockerfile commands and is highly experimental as a result, use with caution.\n")

	app.imageMode = false
	app.artifactMode = false
	app.interactiveDebugging = true
	flagArgs := []string{}
	nonFlagArgs := []string{"+build"}

	return app.actionBuildImp(c, flagArgs, nonFlagArgs)
}

func (app *earthlyApp) actionDocker2Earthly(c *cli.Context) error {
	app.commandName = "docker2earthly"
	err := docker2earthly.Docker2Earthly(app.dockerfilePath, app.earthfilePath, app.earthfileFinalImage)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "An Earthfile has been generated; to run it use: earthly +build; then run with docker run -ti %s\n", app.earthfileFinalImage)
	return nil
}

func (app *earthlyApp) actionConfig(c *cli.Context) error {
	app.commandName = "config"
	if c.NArg() != 2 {
		return errors.New("invalid number of arguments provided")
	}

	args := c.Args().Slice()
	inConfig, err := config.ReadConfigFile(app.configPath)
	if err != nil {
		if c.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "read config")
		}
	}

	var outConfig []byte

	switch args[1] {
	case "-h", "--help":
		if err = config.PrintHelp(args[0]); err != nil {
			return errors.Wrap(err, "help")
		}
		return nil // exit now without writing any changes to config
	case "--delete":
		outConfig, err = config.Delete(inConfig, args[0])
		if err != nil {
			return errors.Wrap(err, "delete config")
		}
	default:
		// args are key/value pairs, e.g. ["global.conversion_parallelism","5"]
		outConfig, err = config.Upsert(inConfig, args[0], args[1])
		if err != nil {
			return errors.Wrap(err, "upsert config")
		}
	}

	if app.configDryRun {
		fmt.Println(string(outConfig))
		return nil
	}

	err = config.WriteConfigFile(app.configPath, outConfig)
	if err != nil {
		return errors.Wrap(err, "write config")
	}

	return nil
}

func (app *earthlyApp) actionBuild(c *cli.Context) error {
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

	flagArgs, nonFlagArgs, err := variables.ParseFlagArgsWithNonFlags(c.Args().Slice())
	if err != nil {
		return errors.Wrapf(err, "parse args %s", strings.Join(c.Args().Slice(), " "))
	}

	return app.actionBuildImp(c, flagArgs, nonFlagArgs)
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

func (app *earthlyApp) actionBuildImp(c *cli.Context, flagArgs, nonFlagArgs []string) error {
	var target domain.Target
	var artifact domain.Artifact
	destPath := "./"
	if app.imageMode {
		if len(nonFlagArgs) == 0 {
			cli.ShowAppHelp(c)
			return errors.Errorf(
				"no image reference provided. Try %s --image +<target-name>", c.App.Name)
		} else if len(nonFlagArgs) != 1 {
			cli.ShowAppHelp(c)
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
			cli.ShowAppHelp(c)
			return errors.Errorf(
				"no artifact reference provided. Try %s --artifact +<target-name>/<artifact-name>", c.App.Name)
		} else if len(nonFlagArgs) > 2 {
			cli.ShowAppHelp(c)
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
			cli.ShowAppHelp(c)
			return errors.Errorf(
				"no target reference provided. Try %s +<target-name>", c.App.Name)
		} else if len(nonFlagArgs) != 1 {
			cli.ShowAppHelp(c)
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

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	// Default upload logs, unless explicitly configured
	if !app.cfg.Global.DisableLogSharing {
		if cc.IsLoggedIn() {
			// If you are logged in, then add the bundle builder code, and configure cleanup and post-build messages.
			app.console = app.console.WithLogBundleWriter(target.String(), cleanCollection)

			defer func() { // Defer this to keep log upload code together
				logPath, err := app.console.WriteBundleToDisk()
				if err != nil {
					err := errors.Wrapf(err, "failed to write log to disk")
					app.console.Warnf(err.Error())
					return
				}

				id, err := cc.UploadLog(logPath)
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

	bkClient, err := app.GetBuildkitClient(c, cc)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	isLocal := containerutil.IsLocal(app.buildkitdSettings.BuildkitAddress)

	bkIP, err := buildkitd.GetContainerIP(c.Context, app.containerName, app.containerFrontend, app.buildkitdSettings)
	if err != nil {
		return errors.Wrap(err, "get buildkit container IP")
	}

	nativePlatform, err := platutil.GetNativePlatformViaBkClient(c.Context, bkClient)
	if err != nil {
		return errors.Wrap(err, "get native platform via buildkit client")
	}
	platr := platutil.NewResolver(nativePlatform)
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

	customSecretProviderCmd, err := secretprovider.NewSecretProviderCmd(app.cfg.Global.SecretProvider)
	if err != nil {
		return errors.Wrap(err, "NewSecretProviderCmd")
	}
	secretProvider := secretprovider.New(
		customSecretProviderCmd,
		secretprovider.NewMapStore(secretsMap),
		secretprovider.NewCloudStore(cc),
	)

	attachables := []session.Attachable{
		secretProvider,
		buildContextProvider,
		localhostProvider,
	}

	switch app.containerFrontend.Config().Setting {
	case containerutil.FrontendDocker, containerutil.FrontendDockerShell:
		attachables = append(attachables, authprovider.NewDockerAuthProvider(os.Stderr))

	case containerutil.FrontendPodman, containerutil.FrontendPodmanShell:
		attachables = append(attachables, authprovider.NewPodmanAuthProvider(os.Stderr))

	default:
		// Old default behavior
		attachables = append(attachables, authprovider.NewDockerAuthProvider(os.Stderr))
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
			err = terminal.ConnectTerm(c.Context, u.Host, debugTermConsole)
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
	if app.remoteCache != "" {
		cacheImports[app.remoteCache] = true
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
	}
	b, err := builder.NewBuilder(c.Context, builderOpts)
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
	_, err = b.BuildTarget(c.Context, target, buildOpts)
	if err != nil {
		return errors.Wrap(err, "build target")
	}

	return nil
}

func (app *earthlyApp) updateGitLookupConfig(gitLookup *buildcontext.GitLookup) error {
	for k, v := range app.cfg.Git {
		if k == "github" || k == "gitlab" || k == "bitbucket" {
			app.console.Warnf("git configuration for %q found, did you mean %q?\n", k, k+".com")
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
		err := gitLookup.AddMatcher(k, pattern, v.Substitute, v.User, v.Password, v.Prefix, suffix, auth, v.ServerKey, ifNilBoolDefault(v.StrictHostKeyChecking, true), v.Port)
		if err != nil {
			return errors.Wrap(err, "gitlookup")
		}
	}
	return nil
}

func ifNilBoolDefault(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func (app *earthlyApp) actionListTargets(c *cli.Context) error {
	app.commandName = "listTargets"

	if c.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}
	var targetToParse string
	if c.NArg() > 0 {
		targetToParse = c.Args().Get(0)
		if !(strings.HasPrefix(targetToParse, "/") || strings.HasPrefix(targetToParse, ".")) {
			return errors.New("remote-paths are not currently supported; local paths must start with \"/\" or \".\"")
		}
		if strings.Contains(targetToParse, "+") {
			return errors.New("path cannot contain a +")
		}
		targetToParse = strings.TrimSuffix(targetToParse, "/Earthfile")
	}

	targetToDisplay := targetToParse
	if targetToParse == "" {
		targetToDisplay = "current directory"
	}

	gitLookup := buildcontext.NewGitLookup(app.console, app.sshAuthSock)
	resolver := buildcontext.NewResolver("", nil, gitLookup, app.console, "")
	var gwClient gwclient.Client // TODO this is a nil pointer which causes a panic if we try to expand a remotely referenced earthfile
	// it's expensive to create this gwclient, so we need to implement a lazy eval which returns it when required.

	target, err := domain.ParseTarget(fmt.Sprintf("%s+base", targetToParse)) // the +base is required to make ParseTarget work; however is ignored by GetTargets
	if err != nil {
		return errors.Errorf("unable to locate Earthfile under %s", targetToDisplay)
	}

	targets, err := earthfile2llb.GetTargets(c.Context, resolver, gwClient, target)
	if err != nil {
		return errors.Errorf("unable to locate Earthfile under %s", targetToDisplay)
	}
	targets = append(targets, "base")
	sort.Strings(targets)
	for _, t := range targets {
		var args []string
		if t != "base" {
			target.Target = t
			args, err = earthfile2llb.GetTargetArgs(c.Context, resolver, gwClient, target)
			if err != nil {
				return err
			}
		}
		if app.lsShowLong {
			fmt.Printf("%s+%s\n", targetToParse, t)
		} else {
			fmt.Printf("+%s\n", t)
		}
		if app.lsShowArgs {
			for _, arg := range args {
				fmt.Printf("  --%s\n", arg)
			}
		}
	}
	return nil
}

func (app *earthlyApp) useSatellite(c *cli.Context, satelliteName, orgID string) error {
	inConfig, err := config.ReadConfigFile(app.configPath)
	if err != nil {
		if c.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "read config")
		}
	}

	newConfig, err := config.Upsert(inConfig, "satellite.name", satelliteName)
	if err != nil {
		return errors.Wrap(err, "could not update satellite name")
	}
	// Update in-place so we can use it later, assuming the config change was successful.
	app.cfg.Satellite.Name = satelliteName

	newConfig, err = config.Upsert(newConfig, "satellite.org", orgID)
	if err != nil {
		return errors.Wrap(err, "could not update satellite name")
	}
	app.cfg.Satellite.Org = orgID
	err = config.WriteConfigFile(app.configPath, newConfig)
	if err != nil {
		return errors.Wrap(err, "could not save config")
	}

	return nil
}

func (app *earthlyApp) printSatellites(satellites []cloud.SatelliteInstance) {
	for _, satellite := range satellites {
		app.console.Printf("name: %s, selected: %t", satellite.Name, app.cfg.Satellite.Name == satellite.Name)
	}
}

func (app *earthlyApp) getSatelliteOrgID(cc cloud.Client) (string, error) {
	var orgID string
	if app.satelliteOrg == "" {
		orgs, err := cc.ListOrgs()
		if err != nil {
			return "", errors.Wrap(err, "failed finding org")
		}
		if len(orgs) != 1 {
			return "", errors.New("more than one organizations available - please specify the name of the organization using `--org`")
		}
		app.satelliteOrg = orgs[0].Name
		orgID = orgs[0].ID
	} else {
		var err error
		orgID, err = cc.GetOrgID(app.satelliteOrg)
		if err != nil {
			return "", errors.Wrap(err, "invalid org provided")
		}
	}
	return orgID, nil
}

func (app *earthlyApp) actionSatelliteLaunch(c *cli.Context) error {
	app.commandName = "launch"

	if c.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	app.satelliteName = c.Args().Get(0)

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cc)
	if err != nil {
		return err
	}

	satellite, err := cc.LaunchSatellite(app.satelliteName, orgID)
	if err != nil {
		return errors.Wrapf(err, "failed to create satellite %s", app.satelliteName)
	}

	err = app.useSatellite(c, app.satelliteName, orgID)
	if err != nil {
		return errors.Wrap(err, "could not configure satellite for use")
	}

	app.printSatellites([]cloud.SatelliteInstance{*satellite})
	return nil
}

func (app *earthlyApp) actionSatelliteList(c *cli.Context) error {
	app.commandName = "list"

	if c.NArg() != 0 {
		return errors.New("command does not accept any arguments")
	}

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cc)
	if err != nil {
		return err
	}

	satellites, err := cc.ListSatellites(orgID)
	if err != nil {
		return err
	}

	app.printSatellites(satellites)
	return nil
}

func (app *earthlyApp) actionSatelliteDestroy(c *cli.Context) error {
	app.commandName = "launch"

	if c.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	app.satelliteName = c.Args().Get(0)

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cc)
	if err != nil {
		return err
	}

	err = cc.DeleteSatellite(app.satelliteName, orgID)
	if err != nil {
		return errors.Wrapf(err, "failed to delete satellite %s", app.satelliteName)
	}

	if app.satelliteName == app.cfg.Satellite.Name {
		// TODO what strategy do we want to use if you delete your current satellite?
		if err = app.useSatellite(c, "", ""); err != nil {
			return errors.Wrapf(err, "failed unselecting satellite")
		}
	}
	return nil
}

func (app *earthlyApp) actionSatelliteDescribe(c *cli.Context) error {
	app.commandName = "describe"

	if c.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	app.satelliteName = c.Args().Get(0)

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cc)
	if err != nil {
		return err
	}

	satellite, err := cc.GetSatellite(app.satelliteName, orgID)
	if err != nil {
		return err
	}

	app.console.Printf("name: %s", satellite.Name)
	app.console.Printf("version: %s", satellite.Version)
	app.console.Printf("platform: %s", satellite.Platform)
	app.console.Printf("status: %s", satellite.Status)
	app.console.Printf("selected: %t", app.satelliteName == satellite.Name)
	return nil
}

func (app *earthlyApp) actionSatelliteSelect(c *cli.Context) error {
	app.commandName = "select"

	if c.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	app.satelliteName = c.Args().Get(0)

	cc, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cc)
	if err != nil {
		return err
	}

	satellites, err := cc.ListSatellites(orgID)
	if err != nil {
		return err
	}

	found := false
	for _, s := range satellites {
		if app.satelliteName == s.Name {
			err = app.useSatellite(c, s.Name, orgID)
			if err != nil {
				return errors.Wrapf(err, "could not select satellite %s", app.satelliteName)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%s is not a valid satellite", app.satelliteName)
	}

	app.printSatellites(satellites)
	return nil
}

func (app *earthlyApp) actionSatelliteUnselect(c *cli.Context) error {
	app.commandName = "unselect"

	if c.NArg() != 0 {
		return errors.New("command does not accept any arguments")
	}

	app.satelliteName = c.Args().Get(0)

	if err := app.useSatellite(c, "", ""); err != nil {
		return errors.Wrap(err, "could not unselect satellite")
	}

	return nil
}

func processSecrets(secrets, secretFiles []string, dotEnvMap map[string]string) (map[string][]byte, error) {
	finalSecrets := make(map[string][]byte)
	for k, v := range dotEnvMap {
		finalSecrets[k] = []byte(v)
	}
	for _, secret := range secrets {
		parts := strings.SplitN(secret, "=", 2)
		key := parts[0]
		var data []byte
		if len(parts) == 2 {
			// secret value passed as argument
			data = []byte(parts[1])
		} else {
			// Not set. Use environment to fetch it.
			value, found := os.LookupEnv(secret)
			if !found {
				return nil, errors.Errorf("env var %s not set", secret)
			}
			data = []byte(value)
		}
		if _, ok := finalSecrets[key]; ok {
			return nil, errors.Errorf("secret %q already contains a value", key)
		}
		finalSecrets[key] = data
	}
	for _, secret := range secretFiles {
		parts := strings.SplitN(secret, "=", 2)
		if len(parts) != 2 {
			return nil, errors.Errorf("unable to parse --secret-file argument: %q", secret)
		}
		k := parts[0]
		path := fileutil.ExpandPath(parts[1])
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open %q", path)
		}
		if _, ok := finalSecrets[k]; ok {
			return nil, errors.Errorf("secret %q already contains a value", k)
		}
		finalSecrets[k] = []byte(data)
	}
	return finalSecrets, nil
}

func defaultConfigPath() string {
	earthlyDir := cliutil.GetEarthlyDir()
	oldConfig := filepath.Join(earthlyDir, "config.yaml")
	newConfig := filepath.Join(earthlyDir, "config.yml")
	oldConfigExists, _ := fileutil.FileExists(oldConfig)
	newConfigExists, _ := fileutil.FileExists(newConfig)
	if oldConfigExists && !newConfigExists {
		return oldConfig
	}
	return newConfig
}
