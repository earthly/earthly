package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/vladaionescu/earthly/buildcontext"
	"github.com/vladaionescu/earthly/builder"
	"github.com/vladaionescu/earthly/buildkitd"
	"github.com/vladaionescu/earthly/cleanup"
	"github.com/vladaionescu/earthly/conslogging"
	"github.com/vladaionescu/earthly/domain"
	"github.com/vladaionescu/earthly/earthfile2llb"
	"github.com/vladaionescu/earthly/earthfile2llb/variables"
	"github.com/vladaionescu/earthly/logging"
	"golang.org/x/sync/errgroup"
)

type earthApp struct {
	cliApp    *cli.App
	ctx       context.Context
	console   conslogging.ConsoleLogger
	sessionID string
	cliFlags
}

type cliFlags struct {
	buildArgs         cli.StringSlice
	secrets           cli.StringSlice
	artifactMode      bool
	imageMode         bool
	push              bool
	noOutput          bool
	noCache           bool
	pruneAll          bool
	buildkitdSettings buildkitd.Settings
	allowPrivileged   bool
	buildkitHost      string
	buildkitdImage    string
}

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

	os.Exit(newEarthApp(context.Background(), conslogging.Current(false)).run(os.Args))
}

func newEarthApp(ctx context.Context, console conslogging.ConsoleLogger) *earthApp {
	sessionIDBytes := make([]byte, 64)
	_, err := rand.Read(sessionIDBytes)
	if err != nil {
		panic(err)
	}
	app := &earthApp{
		cliApp:    cli.NewApp(),
		ctx:       ctx,
		console:   console,
		sessionID: base64.StdEncoding.EncodeToString(sessionIDBytes),
		cliFlags: cliFlags{
			buildkitdSettings: buildkitd.Settings{
				// Add one empty entry for git settings.
				GitSettings: []buildkitd.GitSetting{buildkitd.GitSetting{}},
			},
		},
	}

	app.cliApp.Usage = "The build system for mere mortals"
	app.cliApp.ArgsUsage = "+<target> | <command>"
	app.cliApp.UseShortOptionHandling = true
	app.cliApp.Action = app.actionBuild
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
			Name:        "ssh-auth-sock",
			EnvVars:     []string{"SSH_AUTH_SOCK"},
			Usage:       "The SSH auth socket to use for ssh-agent forwarding",
			Destination: &app.buildkitdSettings.SSHAuthSock,
		},
		&cli.StringFlag{
			Name:        "git-username",
			EnvVars:     []string{"GIT_USERNAME"},
			Usage:       "The git username to use for git HTTPS authentication",
			Destination: &app.buildkitdSettings.GitSettings[0].Username,
		},
		&cli.StringFlag{
			Name:        "git-password",
			EnvVars:     []string{"GIT_PASSWORD"},
			Usage:       "The git password to use for git HTTPS authentication",
			Destination: &app.buildkitdSettings.GitSettings[0].Password,
		},
		&cli.StringFlag{
			Name:        "git-url-instead-of",
			EnvVars:     []string{"GIT_URL_INSTEAD_OF"},
			Usage:       "Rewrite git URLs of a certain pattern. Similar to git-config url.<base>.insteadOf (https://git-scm.com/docs/git-config#Documentation/git-config.txt-urlltbasegtinsteadOf). Format: <base>=<replacement>. For example: 'git@github.com:=https://github.com/'",
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
			EnvVars:     []string{"EARTHLY_BUILDKIT_CACHE_SIZE_MB"},
			Usage:       "The total size of the buildkit cache, in MB",
			Destination: &app.buildkitdSettings.CacheSizeMb,
		},
		&cli.StringFlag{
			Name:        "buildkit-image",
			Value:       "earthly/buildkitd:latest",
			EnvVars:     []string{"EARTHLY_BUILDKITD_IMAGE"},
			Usage:       "The docker image to use for the buildkit daemon",
			Destination: &app.buildkitdImage,
		},
	}

	app.cliApp.Commands = []*cli.Command{
		{
			Name:        "debug",
			Usage:       "Print debug information about a build.earth file",
			Description: "Print debug information about a build.earth file",
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
			},
		},
	}
	return app
}

func (app *earthApp) run(args []string) int {
	joinedArgs := ""
	if len(args) > 2 {
		joinedArgs = strings.Join(args[2:], " ")
	}
	app.ctx = logging.With(app.ctx, logging.COMMAND, fmt.Sprintf("earth %s", joinedArgs))
	err := app.cliApp.Run(args)
	if err != nil {
		logging.GetLogger(app.ctx).Error(err)
		app.console.Printf("Error: %v\n", err)
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
	path = filepath.Join(path, "build.earth")

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
	bkClient, err := app.newBuildkitdClient()
	if err != nil {
		return errors.Wrap(err, "buildkitd new client")
	}
	defer bkClient.Close()
	var opts []client.PruneOption
	if app.pruneAll {
		opts = append(opts, client.PruneAll)
	}
	ch := make(chan client.UsageInfo, 1)
	eg, ctx := errgroup.WithContext(app.ctx)
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
			return errors.New("invalid number of args")
		}
		targetName := c.Args().Get(0)
		var err error
		target, err = domain.ParseTarget(targetName)
		if err != nil {
			return errors.Wrapf(err, "parse target name %s", targetName)
		}
	}
	bkClient, err := app.newBuildkitdClient()
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
	b, err := builder.NewBuilder(app.ctx, bkClient, app.console, attachables, enttlmnts, app.noCache)
	if err != nil {
		return errors.Wrap(err, "new builder")
	}

	buildArgsSlice := appendBuiltinBuildArgs(app.buildArgs.Value(), target)
	buildArgs, err := parseBuildArgs(buildArgsSlice)
	if err != nil {
		return errors.Wrap(err, "parse build args")
	}
	cleanCollection := cleanup.NewCollection()
	defer cleanCollection.Close()
	mts, err := earthfile2llb.Earthfile2LLB(
		app.ctx, target, resolver, b.BuildOnlyLastImageAsTar, cleanCollection,
		nil, buildArgs)
	if err != nil {
		return errors.Wrap(err, "parse build")
	}

	if app.imageMode {
		err = b.BuildOnlyImages(app.ctx, mts, app.push)
		if err != nil {
			return errors.Wrap(err, "build only image")
		}
	} else if app.artifactMode {
		err = b.BuildOnlyArtifact(app.ctx, mts, artifact, destPath)
		if err != nil {
			return errors.Wrap(err, "build only artifact")
		}
	} else {
		err = b.Build(app.ctx, mts, app.noOutput, app.push)
		if err != nil {
			return errors.Wrap(err, "build")
		}
	}
	return nil
}

func (app *earthApp) newBuildkitdClient(opts ...client.ClientOpt) (*client.Client, error) {
	if app.buildkitHost == "" {
		// Start our own.
		bkClient, err := buildkitd.NewClient(app.ctx, app.console, app.buildkitdImage, app.buildkitdSettings)
		if err != nil {
			return nil, errors.Wrap(err, "buildkitd new client (own)")
		}
		return bkClient, nil
	}

	// Use provided.
	bkClient, err := client.New(app.ctx, app.buildkitHost, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "buildkitd new client (provided)")
	}
	return bkClient, nil
}

func appendBuiltinBuildArgs(buildArgs []string, target domain.Target) []string {
	return buildArgs
}

func parseBuildArgs(buildArgs []string) (map[string]variables.Variable, error) {
	out := make(map[string]variables.Variable)
	for _, arg := range buildArgs {
		splitArg := strings.SplitN(arg, "=", 2)
		if len(splitArg) < 1 {
			return nil, fmt.Errorf("Invalid build arg %s", splitArg)
		}
		key := splitArg[0]
		value := ""
		if len(splitArg) == 2 {
			value = splitArg[1]
		}
		if value == "" {
			value = os.Getenv(key)
		}
		out[key] = variables.NewConstant(value)
	}
	return out, nil
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
