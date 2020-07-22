package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
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
	buildArgs           cli.StringSlice
	secrets             cli.StringSlice
	artifactMode        bool
	imageMode           bool
	push                bool
	noOutput            bool
	noCache             bool
	pruneAll            bool
	pruneReset          bool
	buildkitdSettings   buildkitd.Settings
	allowPrivileged     bool
	buildkitHost        string
	buildkitdImage      string
	remoteCache         string
	configPath          string
	gitUsernameOverride string
	gitPasswordOverride string
}

var (
	// DefaultBuildkitdImage is the default buildkitd image to use.
	DefaultBuildkitdImage string
	// Version is the version of this CLI app.
	Version string
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
		fmt.Printf("Warning: cannot create dir %s\n", logDir)
	} else {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Printf("Warning: cannot open log file for writing %s\n", logFile)
		} else {
			logrus.SetOutput(f)
		}
	}

	ctx := context.Background()
	os.Exit(newEarthApp(ctx, conslogging.Current(false)).run(ctx, os.Args))
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
	app.cliApp.Version = Version
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
			Aliases:     []string{"i"},
			Usage:       "Output only docker image of the specified target",
			Destination: &app.imageMode,
		},
		&cli.BoolFlag{
			Name:        "push",
			EnvVars:     []string{"EARTHLY_PUSH"},
			Usage:       "Push docker images and execute RUN --push commmands",
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
			Value:       "https://gitlab.com/=git@gitlab.com:,https://github.com/=git@github.com:",
			EnvVars:     []string{"GIT_URL_INSTEAD_OF"},
			Usage:       "Rewrite git URLs of a certain pattern. Similar to git-config url.<base>.insteadOf (https://git-scm.com/docs/git-config#Documentation/git-config.txt-urlltbasegtinsteadOf). Multiple values can be separated by commas. Format: <base>=<instead-of>[,...]. For example: ''git@github.com:=https://github.com/'",
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

	// command line overrides the config file
	if app.gitUsernameOverride != "" || app.gitPasswordOverride != "" {
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

	gitConfig, gitCredentials, err := config.CreateGitConfig(cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to create git config from %s", app.configPath)
	}

	app.buildkitdSettings.GitConfig = gitConfig
	app.buildkitdSettings.GitCredentials = gitCredentials
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
		app.console.Printf("Error: %v\n", err)
		if strings.Contains(err.Error(), "failed to fetch remote") {
			app.console.Printf(
				"Check your git auth settings.\n" +
					"Did you ssh-add today? Need to configure ~/earthly/config.yaml?\n" +
					"For more information see https://docs.earthly.dev/guides/auth\n")
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
		if c.NArg() != 1 {
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
		if c.NArg() != 1 && c.NArg() != 2 {
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
		if c.NArg() != 1 {
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
	secrets := processSecrets(app.secrets.Value())
	attachables := []session.Attachable{
		secrets,
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
		nil, varCollection)
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

func processSecrets(secrets []string) session.Attachable {
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
	return secretsprovider.FromMap(finalSecrets)
}

func defaultSSHAuthSock() string {
	if runtime.GOOS == "darwin" {
		return "/run/host-services/ssh-auth.sock"
	}
	return os.Getenv("SSH_AUTH_SOCK")
}
