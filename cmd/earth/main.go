package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/fatih/color"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/server"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/earthfile2llb/variables"
	"github.com/earthly/earthly/logging"

	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

type earthApp struct {
	cliApp    *cli.App
	console   conslogging.ConsoleLogger
	sessionID string
	cliFlags
}

type cliFlags struct {
	buildArgs            cli.StringSlice
	secrets              cli.StringSlice
	artifactMode         bool
	imageMode            bool
	push                 bool
	noOutput             bool
	noCache              bool
	pruneAll             bool
	pruneReset           bool
	buildkitdSettings    buildkitd.Settings
	allowPrivileged      bool
	buildkitHost         string
	buildkitdImage       string
	debuggerImage        string
	remoteCache          string
	configPath           string
	gitUsernameOverride  string
	gitPasswordOverride  string
	interactiveDebugging bool
}

var (
	// DefaultBuildkitdImage is the default buildkitd image to use.
	DefaultBuildkitdImage string

	// DefaultDebuggerImage is the default debugger image to use.
	DefaultDebuggerImage string

	// Version is the version of this CLI app.
	Version string

	// GitSha contains the git sha used to build this app
	GitSha string
)

func main() {
	// Set up file-based logging.
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	logDir := filepath.Join(os.Getenv("HOME"), ".earthly")
	logFile := filepath.Join(logDir, "earth.log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cannot create dir %s\n", logDir)
	} else {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot open log file for writing %s\n", logFile)
		} else {
			logrus.SetOutput(f)
		}
	}

	ctx := context.Background()
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
	os.Exit(newEarthApp(ctx, conslogging.Current(colorMode)).run(ctx, os.Args))
}

func getVersion() string {
	var isRelease = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
	if isRelease.MatchString(Version) {
		return Version
	}
	return fmt.Sprintf("%s-%s", Version, GitSha)
}

func newEarthApp(ctx context.Context, console conslogging.ConsoleLogger) *earthApp {
	sessionIDBytes := make([]byte, 64)
	_, err := rand.Read(sessionIDBytes)
	if err != nil {
		panic(err)
	}
	app := &earthApp{
		cliApp:    cli.NewApp(),
		console:   console,
		sessionID: base64.StdEncoding.EncodeToString(sessionIDBytes),
		cliFlags: cliFlags{
			buildkitdSettings: buildkitd.Settings{},
		},
	}

	app.cliApp.Usage = "A build automation tool for the container era"
	app.cliApp.UsageText = "\tearth [options] <target-ref>\n" +
		"\n" +
		"   \tearth [options] --image <target-ref>\n" +
		"\n" +
		"   \tearth [options] --artifact <artifact-ref> [<dest-path>]\n" +
		"\n" +
		"   \tearth [options] command [command options]\n" +
		"\n" +
		"Executes Earthly builds. For more information see https://docs.earthly.dev/earth-command.\n" +
		"To get started with using Earthly, check out the getting started guide at https://docs.earthly.dev/guides/basics."
	app.cliApp.UseShortOptionHandling = true
	app.cliApp.Action = app.actionBuild
	app.cliApp.Version = getVersion()
	app.cliApp.Flags = []cli.Flag{
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
		&cli.BoolFlag{
			Name:        "artifact",
			Aliases:     []string{"a"},
			Usage:       "Output only specified artifact",
			Destination: &app.artifactMode,
		},
		&cli.BoolFlag{
			Name:        "image",
			Usage:       "Output only docker image of the specified target",
			Destination: &app.imageMode,
		},
		&cli.BoolFlag{
			Name:        "push",
			EnvVars:     []string{"EARTHLY_PUSH"},
			Usage:       "Push docker images and execute RUN --push commands",
			Destination: &app.push,
		},
		&cli.BoolFlag{
			Name:        "no-output",
			EnvVars:     []string{"EARTHLY_NO_OUTPUT"},
			Usage:       "Do not output artifacts or images",
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
			Value:       filepath.Join(os.Getenv("HOME"), ".earthly", "config.yaml"),
			EnvVars:     []string{"EARTHLY_CONFIG"},
			Usage:       "Path to config file for",
			Destination: &app.configPath,
		},
		&cli.StringFlag{
			Name:        "ssh-auth-sock",
			Value:       defaultSSHAuthSock(),
			EnvVars:     []string{"EARTHLY_SSH_AUTH_SOCK"},
			Usage:       "The SSH auth socket to use for ssh-agent forwarding",
			Destination: &app.buildkitdSettings.SSHAuthSock,
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
		&cli.StringFlag{
			Name:        "git-url-instead-of",
			Value:       "",
			EnvVars:     []string{"GIT_URL_INSTEAD_OF"},
			Usage:       "Rewrite git URLs of a certain pattern. Similar to git-config url.<base>.insteadOf (https://git-scm.com/docs/git-config#Documentation/git-config.txt-urlltbasegtinsteadOf). Multiple values can be separated by commas. Format: <base>=<instead-of>[,...]. For example: 'https://github.com/=git@github.com:'",
			Destination: &app.buildkitdSettings.GitURLInsteadOf,
		},
		&cli.BoolFlag{
			Name:        "allow-privileged",
			Aliases:     []string{"P"},
			EnvVars:     []string{"EARTHLY_ALLOW_PRIVILEGED"},
			Usage:       "Allow build to use the --privileged flag in RUN commands",
			Destination: &app.allowPrivileged,
		},
		&cli.StringFlag{
			Name:        "buildkit-host",
			EnvVars:     []string{"EARTHLY_BUILDKIT_HOST"},
			Usage:       "The URL to use for connecting to a buildkit host. If empty, earth will attempt to start a buildkitd instance via docker run",
			Destination: &app.buildkitHost,
		},
		&cli.IntFlag{
			Name:        "buildkit-cache-size-mb",
			Value:       10000,
			EnvVars:     []string{"EARTHLY_BUILDKIT_CACHE_SIZE_MB"},
			Usage:       "The total size of the buildkit cache, in MB",
			Destination: &app.buildkitdSettings.CacheSizeMb,
		},
		&cli.StringFlag{
			Name:        "buildkit-image",
			Value:       DefaultBuildkitdImage,
			EnvVars:     []string{"EARTHLY_BUILDKIT_IMAGE"},
			Usage:       "The docker image to use for the buildkit daemon",
			Destination: &app.buildkitdImage,
		},
		&cli.StringFlag{
			Name:        "debugger-image",
			Value:       DefaultDebuggerImage,
			EnvVars:     []string{"EARTHLY_DEBUGGER_IMAGE"},
			Usage:       "The docker image to use for the interactive debugger process",
			Destination: &app.debuggerImage,
			Hidden:      true,
		},
		&cli.BoolFlag{
			Name:        "no-loop-device",
			EnvVars:     []string{"EARTHLY_NO_LOOP_DEVICE"},
			Usage:       "Disables the use of a loop device for storing the cache contents",
			Destination: &app.buildkitdSettings.DisableLoopDevice,
		},
		&cli.StringFlag{
			Name:        "remote-cache",
			EnvVars:     []string{"EARTHLY_REMOTE_CACHE"},
			Usage:       "A remote docker image repository to be used as build cache",
			Destination: &app.remoteCache,
			Hidden:      true, // Experimental.
		},
		&cli.BoolFlag{
			Name:        "interactive",
			Aliases:     []string{"i"},
			EnvVars:     []string{"EARTHLY_INTERACTIVE"},
			Usage:       "Enable interactive debugging",
			Destination: &app.interactiveDebugging,
			Hidden:      true, // Experimental.
		},
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"V"},
			EnvVars:     []string{"EARTHLY_VERBOSE"},
			Usage:       "enable verbose logging of the earthly-buildkitd container",
			Destination: &app.buildkitdSettings.Debug,
		},
	}

	app.cliApp.Commands = []*cli.Command{
		{
			Name:        "debug",
			Usage:       "Print debug information about an Earthfile",
			Description: "Print debug information about an Earthfile",
			ArgsUsage:   "[<path>]",
			Hidden:      true,
			Action:      app.actionDebug,
		},
		{
			Name:        "prune",
			Usage:       "Prune earthly build cache",
			Description: "Prune earthly build cache",
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
	}

	app.cliApp.Before = app.parseConfigFile
	return app
}

func (app *earthApp) parseConfigFile(context *cli.Context) error {
	if context.IsSet("config") {
		app.console.Printf("loading config values from %q\n", app.configPath)
	}

	yamlData, err := ioutil.ReadFile(app.configPath)
	if os.IsNotExist(err) && !context.IsSet("config") {
		yamlData = []byte{}
	} else if err != nil {
		return errors.Wrapf(err, "failed to read from %s", app.configPath)
	}

	cfg, err := config.ParseConfigFile(yamlData)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", app.configPath)
	}

	if cfg.Git == nil {
		cfg.Git = map[string]config.GitConfig{}
	}

	err = app.processDeprecatedCommandOptions(context, cfg)
	if err != nil {
		return err
	}

	gitConfig, gitCredentials, err := config.CreateGitConfig(cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to create git config from %s", app.configPath)
	}

	// command line option overrides the config which overrides the default value
	if !context.IsSet("buildkit-image") && cfg.Global.BuildkitImage != "" {
		app.buildkitdImage = cfg.Global.BuildkitImage
	}
	if !context.IsSet("debugger-image") && cfg.Global.DebuggerImage != "" {
		app.debuggerImage = cfg.Global.DebuggerImage
	}

	if app.buildkitdSettings.SSHAuthSock != "" {
		// EvalSymlinks evaluates "" as "." which then breaks docker volume mounting
		realSSHSocketPath, err := filepath.EvalSymlinks(app.buildkitdSettings.SSHAuthSock)
		if err != nil {
			if runtime.GOOS != "darwin" {
				app.console.Warnf("failed to evaluate potential symbolic links in ssh auth socket %q: %v\n", app.buildkitdSettings.SSHAuthSock, err)
			} // else ignore the error on mac
		} else {
			app.buildkitdSettings.SSHAuthSock = realSSHSocketPath
		}
	}

	app.buildkitdSettings.TempDir = cfg.Global.CachePath
	app.buildkitdSettings.GitConfig = gitConfig
	app.buildkitdSettings.GitCredentials = gitCredentials
	return nil

}

func (app *earthApp) processDeprecatedCommandOptions(context *cli.Context, cfg *config.Config) error {

	// command line overrides the config file
	if app.gitUsernameOverride != "" || app.gitPasswordOverride != "" {
		app.console.Warnf("Warning: the --git-username and --git-password command flags are deprecated and are now configured in the ~/.earthly/config.yaml file under the git section; see https://docs.earthly.dev/earth-config for reference.\n")
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

	if context.IsSet("git-url-instead-of") {
		app.console.Warnf("Warning: the --git-url-instead-of command flag is deprecated and is now configured in the ~/.earthly/config.yaml file under the git global url_instead_of setting; see https://docs.earthly.dev/earth-config for reference.\n")
	} else {
		if gitGlobal, ok := cfg.Git["global"]; ok {
			if gitGlobal.GitURLInsteadOf != "" {
				app.buildkitdSettings.GitURLInsteadOf = gitGlobal.GitURLInsteadOf
			}
		}
	}

	if context.IsSet("no-loop-device") {
		app.console.Warnf("Warning: the --no-loop-device command flag is deprecated and is now configured in the ~/.earthly/config.yaml file under the no_loop_device setting; see https://docs.earthly.dev/earth-config for reference.\n")
	} else {
		app.buildkitdSettings.DisableLoopDevice = cfg.Global.DisableLoopDevice
	}

	if context.IsSet("buildkit-cache-size-mb") {
		app.console.Warnf("Warning: the --buildkit-cache-size-mb command flag is deprecated and is now configured in the ~/.earthly/config.yaml file under the buildkit_cache_size setting; see https://docs.earthly.dev/earth-config for reference.\n")
	} else {
		app.buildkitdSettings.CacheSizeMb = cfg.Global.BuildkitCacheSizeMb
	}

	return nil

}

func (app *earthApp) run(ctx context.Context, args []string) int {
	joinedArgs := ""
	if len(args) > 2 {
		joinedArgs = strings.Join(args[2:], " ")
	}
	ctx = logging.With(ctx, logging.COMMAND, fmt.Sprintf("earth %s", joinedArgs))
	err := app.cliApp.RunContext(ctx, args)
	if err != nil {
		logging.GetLogger(ctx).Error(err)
		if errors.Is(err, builder.ErrSolve) {
			// Do not print error if it's a solve error. It has already been printed somewhere else.
		} else if strings.Contains(err.Error(), "failed to fetch remote") {
			app.console.Printf("Error: %v\n", err)
			app.console.Printf(
				"Check your git auth settings.\n" +
					"Did you ssh-add today? Need to configure ~/earthly/config.yaml?\n" +
					"For more information see https://docs.earthly.dev/guides/auth\n")
		} else {
			app.console.Printf("Error: %v\n", err)
		}
		return 1
	}
	return 0
}

func (app *earthApp) actionDebug(c *cli.Context) error {
	if c.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := "."
	if c.NArg() == 1 {
		path = c.Args().First()
	}
	path = filepath.Join(path, "Earthfile")

	err := earthfile2llb.ParseDebug(path)
	if err != nil {
		return errors.Wrap(err, "parse debug")
	}
	return nil
}

func (app *earthApp) actionPrune(c *cli.Context) error {
	if c.NArg() != 0 {
		return errors.New("invalid arguments")
	}
	if app.pruneReset {
		// Prune by resetting container.
		if app.buildkitHost != "" {
			return errors.New("Cannot use prune --reset on non-default buildkit-host setting")
		}
		err := buildkitd.ResetCache(
			c.Context, app.console, app.buildkitdImage, app.buildkitdSettings)
		if err != nil {
			return errors.Wrap(err, "reset cache")
		}
		return nil
	}

	// Prune via API.
	bkClient, err := app.newBuildkitdClient(c.Context)
	if err != nil {
		return errors.Wrap(err, "buildkitd new client")
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

func (app *earthApp) actionBuild(c *cli.Context) error {
	var remoteConsoleAddr string
	var err error
	if app.interactiveDebugging {
		debugServer := server.NewDebugServer(c.Context, app.console)
		remoteConsoleAddr, err = debugServer.Start()
		if err != nil {
			app.console.Warnf("failed to open remote console listener: %v; interactive debugging disabled\n", err)
			app.interactiveDebugging = false
		}
		defer debugServer.Stop()
	}

	if app.imageMode && app.artifactMode {
		return errors.New("both image and artifact modes cannot be active at the same time")
	}
	if (app.imageMode && app.noOutput) || (app.artifactMode && app.noOutput) {
		return errors.New("cannot use --no-output with image or artifact modes")
	}
	if app.push && app.noOutput {
		return errors.New("cannot use --no-output with --push")
	}
	var target domain.Target
	var artifact domain.Artifact
	destPath := "./"
	if app.imageMode {
		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			return fmt.Errorf(
				"no image reference provided. Try %s --image +<target-name>", c.App.Name)
		} else if c.NArg() != 1 {
			cli.ShowAppHelp(c)
			return errors.New("invalid number of args")
		}
		targetName := c.Args().Get(0)
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return errors.Wrapf(err, "parse target name %s", targetName)
		}
	} else if app.artifactMode {
		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			return fmt.Errorf(
				"no artifact reference provided. Try %s --artifact +<target-name>/<artifact-name>", c.App.Name)
		} else if c.NArg() != 1 && c.NArg() != 2 {
			cli.ShowAppHelp(c)
			return errors.New("invalid number of args")
		}
		artifactName := c.Args().Get(0)
		if c.NArg() == 2 {
			destPath = c.Args().Get(1)
		}
		var err error
		artifact, err = domain.ParseArtifact(artifactName)
		if err != nil {
			return errors.Wrapf(err, "parse artifact name %s", artifactName)
		}
		target = artifact.Target
	} else {
		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			return fmt.Errorf(
				"no target reference provided. Try %s +<target-name>", c.App.Name)
		} else if c.NArg() != 1 {
			cli.ShowAppHelp(c)
			return errors.New("invalid number of args")
		}
		targetName := c.Args().Get(0)
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return errors.Wrapf(err, "parse target name %s", targetName)
		}
	}
	bkClient, err := app.newBuildkitdClient(c.Context)
	if err != nil {
		return errors.Wrap(err, "buildkitd new client")
	}
	defer bkClient.Close()
	resolver := buildcontext.NewResolver(bkClient, app.console, app.sessionID)
	defer resolver.Close()
	secrets := app.secrets.Value()
	//interactive debugger settings are passed as secrets to avoid having it affect the cache hash

	secretsMap := processSecrets(secrets)

	debuggerSettings := debuggercommon.DebuggerSettings{
		Debug:   true,
		Addrs:   []string{remoteConsoleAddr},
		Enabled: app.interactiveDebugging,
	}

	debuggerSettingsData, err := json.Marshal(&debuggerSettings)
	if err != nil {
		return errors.Wrap(err, "debugger settings json marshal")
	}
	secretsMap["earthly_debugger_settings"] = debuggerSettingsData

	attachables := []session.Attachable{
		secretsprovider.FromMap(secretsMap),
		authprovider.NewDockerAuthProvider(os.Stderr),
	}
	var enttlmnts []entitlements.Entitlement
	if app.allowPrivileged {
		enttlmnts = append(enttlmnts, entitlements.EntitlementSecurityInsecure)
	}
	b, err := builder.NewBuilder(
		c.Context, bkClient, app.console, attachables, enttlmnts, app.noCache, app.remoteCache)
	if err != nil {
		return errors.Wrap(err, "new builder")
	}

	varCollection, err := variables.ParseCommandLineBuildArgs(app.buildArgs.Value())
	if err != nil {
		return errors.Wrap(err, "parse build args")
	}
	cleanCollection := cleanup.NewCollection()
	defer cleanCollection.Close()
	mts, err := earthfile2llb.Earthfile2LLB(
		c.Context, target, resolver, b.BuildOnlyLastImageAsTar, cleanCollection,
		nil, varCollection, app.interactiveDebugging, app.debuggerImage, remoteConsoleAddr)
	if err != nil {
		return err
	}

	if app.imageMode {
		err = b.BuildOnlyImages(c.Context, mts, app.push)
		if err != nil {
			return err
		}
	} else if app.artifactMode {
		err = b.BuildOnlyArtifact(c.Context, mts, artifact, destPath)
		if err != nil {
			return err
		}
	} else {
		err = b.Build(c.Context, mts, app.noOutput, app.push)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *earthApp) newBuildkitdClient(ctx context.Context, opts ...client.ClientOpt) (*client.Client, error) {
	if app.buildkitHost == "" {
		// Start our own.
		bkClient, err := buildkitd.NewClient(ctx, app.console, app.buildkitdImage, app.buildkitdSettings)
		if err != nil {
			return nil, errors.Wrap(err, "buildkitd new client (own)")
		}
		return bkClient, nil
	}

	// Use provided.
	bkClient, err := client.New(ctx, app.buildkitHost, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "buildkitd new client (provided)")
	}
	return bkClient, nil
}

func processSecrets(secrets []string) map[string][]byte {
	finalSecrets := make(map[string][]byte)
	for _, secret := range secrets {
		parts := strings.SplitN(secret, "=", 2)
		if len(parts) == 2 {
			// Already set.
			finalSecrets[parts[0]] = []byte(parts[1])
		} else {
			// Not set. Use environment to fetch it.
			value := os.Getenv(secret)
			finalSecrets[secret] = []byte(value)
		}
	}
	return finalSecrets
}

func defaultSSHAuthSock() string {
	if runtime.GOOS == "darwin" {
		return "/run/host-services/ssh-auth.sock"
	}
	return os.Getenv("SSH_AUTH_SOCK")
}
