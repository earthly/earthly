package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/moby/buildkit/client"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/docker2earthly"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/termutil"
)

func (app *earthlyApp) rootCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "build",
			Usage:       "Build an Earthly target",
			Description: "Build an Earthly target",
			Action:      app.actionBuild,
			Flags:       app.buildFlags(),
			Hidden:      true, // Meant to be used mainly for help output.
		},
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
				&cli.StringFlag{
					Name:        "certs-hostname",
					Usage:       "hostname to generate certificates for",
					EnvVars:     []string{"EARTHLY_CERTS_HOSTNAME"},
					Hidden:      true,
					Value:       "localhost",
					Destination: &app.certsHostName,
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
			Name:        "org",
			Aliases:     []string{"orgs"},
			Usage:       "Earthly organization administration *beta*",
			Subcommands: app.orgCmds(),
		},
		{
			Name:      "doc",
			Usage:     "Document targets from an Earthfile *beta*",
			UsageText: "earthly [options] doc [<project-ref>[+<target-ref>]]",
			Action:    app.actionDocumentTarget,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "long",
					Aliases:     []string{"l"},
					Usage:       "Show full details for all target inputs and outputs",
					Destination: &app.docShowLong,
				},
			},
		},
		{
			Name:      "ls",
			Usage:     "List targets from an Earthfile *beta*",
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
			Name:        "account",
			Usage:       "Create or manage an Earthly account *beta*",
			Subcommands: app.accountCmds(),
		},
		{
			Name:        "debug",
			Usage:       "Print debug information about an Earthfile",
			Description: "Print debug information about an Earthfile",
			ArgsUsage:   "[<path>]",
			Hidden:      true, // Dev purposes only.
			Subcommands: app.debugCmds(),
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
				&cli.GenericFlag{
					Name:  "size",
					Usage: "Prune cache to specified size, starting from oldest",
					Value: &app.pruneTargetSize,
				},
				&cli.DurationFlag{
					Name:        "age",
					Usage:       "Prune cache older than the specified duration",
					Destination: &app.pruneKeepDuration,
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
			Name:    "satellite",
			Aliases: []string{"satellites", "sat"},
			Usage: "Launch and use a Satellite runner as remote backend for Earthly builds. *experimental*\n" +
				"	Satellites can be used to optimize and share cache between multiple builds and users,\n" +
				"	as well as run builds in native architectures independent of where the Earthly client is invoked.\n" +
				"	Note: This feature is currently experimental.\n" +
				"	If you'd like to try it out, please contact us at support@earthly.dev or by visiting https://earthly.dev/slack.",
			UsageText:   "earthly satellite (launch|ls|inspect|select|unselect|rm)",
			Description: "Create and manage Earthly Satellites *beta*",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the organization the satellite belongs to. Required when user is a member of multiple.",
					Required:    false,
					Destination: &app.orgName,
				},
			},
			Subcommands: app.satelliteCmds(),
		},
		{
			Name:        "project",
			Aliases:     []string{"projects"},
			Description: "Manage Earthly projects *beta*",
			Usage:       "Manage Earthly projects *beta*",
			UsageText:   "earthly project (ls|rm|create|member)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the organization to which the project belongs. Required when user is a member of multiple.",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					Aliases:     []string{"p"},
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The project to act on.",
					Required:    false,
					Destination: &app.projectName,
				},
			},
			Subcommands: app.projectCmds(),
		},
		{
			Name:        "secret",
			Aliases:     []string{"secrets"},
			Description: "Manage cloud secrets *beta*",
			Usage:       "Manage cloud secrets *beta*",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs.",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store secrets.",
					Required:    false,
					Destination: &app.projectName,
				},
			},
			Subcommands: app.secretCmds(),
		},
		{
			Name:        "registry",
			Aliases:     []string{"registries"},
			Description: "Manage registry access *beta*",
			Usage:       "Manage registry access *beta*",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs.",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store registry credentials.",
					Required:    false,
					Destination: &app.projectName,
				},
			},
			Subcommands: app.registryCmds(),
		},
		{
			Name:        "web",
			Description: "Access the web UI *beta*",
			Usage:       "Access the web UI via your default browser and print the url *beta*",
			UsageText:   "earthly web (--provider=github)",
			Hidden:      true,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "provider",
					EnvVars:     []string{"EARTHLY_LOGIN_PROVIDER"},
					Usage:       "The provider to use when logging into the web ui",
					Required:    false,
					Destination: &app.loginProvider,
				},
			},
			Action: app.webUI,
		},
	}
}

func (app *earthlyApp) actionBootstrap(cliCtx *cli.Context) error {
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

	return app.bootstrap(cliCtx)
}

func (app *earthlyApp) bootstrap(cliCtx *cli.Context) error {
	var err error
	console := app.console.WithPrefix("bootstrap")
	defer func() {
		// cliutil.IsBootstrapped() determines if bootstrapping was done based
		// on the existence of ~/.earthly; therefore we must ensure it's created.
		_, err := cliutil.GetOrCreateEarthlyDir(app.installationName)
		if err != nil {
			console.Warnf("Warning: Failed to create Earthly Dir: %v", err)
			// Keep going.
		}
		err = cliutil.EnsurePermissions(app.installationName)
		if err != nil {
			console.Warnf("Warning: Failed to ensure permissions: %v", err)
			// Keep going.
		}
	}()

	if app.bootstrapWithAutocomplete {
		// Because this requires sudo, it should warn and not fail the rest of it.
		err = app.insertBashCompleteEntry()
		if err != nil {
			console.Warnf("Warning: %s\n", err.Error())
			// Keep going.
		}
		err = app.insertZSHCompleteEntry()
		if err != nil {
			console.Warnf("Warning: %s\n", err.Error())
			// Keep going.
		}

		console.Printf("You may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")
	}

	err = symlinkEarthlyToEarth()
	if err != nil {
		console.Warnf("Warning: %s\n", err.Error())
		err = nil
	}

	if !app.bootstrapNoBuildkit && !app.isUsingSatellite(cliCtx) {
		bkURL, err := url.Parse(app.buildkitHost)
		if err != nil {
			return errors.Wrapf(err, "invalid buildkit_host: %s", app.buildkitHost)
		}
		if bkURL.Scheme == "tcp" && app.cfg.Global.TLSEnabled {
			err := buildkitd.GenCerts(*app.cfg, app.certsHostName)
			if err != nil {
				return errors.Wrap(err, "failed to generate TLS certs")
			}
		}

		// Bootstrap buildkit - pulls image and starts daemon.
		bkClient, err := app.getBuildkitClient(cliCtx, nil)
		if err != nil {
			console.Warnf("Warning: Bootstrapping buildkit failed: %v", err)
			// Keep going.
		} else {
			defer bkClient.Close()
		}
	}

	console.Printf("Bootstrapping successful.\n")
	return nil
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

func (app *earthlyApp) actionDocker(cliCtx *cli.Context) error {
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

	return app.actionBuildImp(cliCtx, flagArgs, nonFlagArgs)
}

func (app *earthlyApp) actionDocker2Earthly(cliCtx *cli.Context) error {
	app.commandName = "docker2earthly"
	err := docker2earthly.Docker2Earthly(app.dockerfilePath, app.earthfilePath, app.earthfileFinalImage)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "An Earthfile has been generated; to run it use: earthly +build; then run with docker run -ti %s\n", app.earthfileFinalImage)
	return nil
}

func (app *earthlyApp) actionConfig(cliCtx *cli.Context) error {
	app.commandName = "config"
	if cliCtx.NArg() != 2 {
		return errors.New("invalid number of arguments provided")
	}

	args := cliCtx.Args().Slice()
	inConfig, err := config.ReadConfigFile(app.configPath)
	if err != nil {
		if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
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
	app.console.Printf("Updated config file %s", app.configPath)

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

func (app *earthlyApp) actionListTargets(cliCtx *cli.Context) error {
	app.commandName = "listTargets"

	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}
	var targetToParse string
	if cliCtx.NArg() > 0 {
		targetToParse = cliCtx.Args().Get(0)
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
	resolver := buildcontext.NewResolver(nil, gitLookup, app.console, "", app.gitBranchOverride)
	var gwClient gwclient.Client // TODO this is a nil pointer which causes a panic if we try to expand a remotely referenced earthfile
	// it's expensive to create this gwclient, so we need to implement a lazy eval which returns it when required.

	target, err := domain.ParseTarget(fmt.Sprintf("%s+base", targetToParse)) // the +base is required to make ParseTarget work; however is ignored by GetTargets
	if errors.Is(err, buildcontext.ErrEarthfileNotExist{}) {
		return errors.Errorf("unable to locate Earthfile under %s", targetToDisplay)
	} else if err != nil {
		return err
	}

	targets, err := earthfile2llb.GetTargets(cliCtx.Context, resolver, gwClient, target)
	if err != nil {
		switch err := errors.Cause(err).(type) {
		case *buildcontext.ErrEarthfileNotExist:
			return errors.Errorf("unable to locate Earthfile under %s", targetToDisplay)
		default:
			return err
		}
	}

	targets = append(targets, "base")
	sort.Strings(targets)
	for _, t := range targets {
		var args []string
		if t != "base" {
			target.Target = t
			args, err = earthfile2llb.GetTargetArgs(cliCtx.Context, resolver, gwClient, target)
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

func (app *earthlyApp) actionPrune(cliCtx *cli.Context) error {
	app.commandName = "prune"
	if cliCtx.NArg() != 0 {
		return errors.New("invalid arguments")
	}
	if app.pruneReset {
		if app.isUsingSatellite(cliCtx) {
			return errors.New("Cannot prune --reset when using a satellite. Try without --reset")
		}
		err := app.initFrontend(cliCtx)
		if err != nil {
			return err
		}
		err = buildkitd.ResetCache(cliCtx.Context, app.console, app.buildkitdImage, app.containerName, app.installationName, app.containerFrontend, app.buildkitdSettings)
		if err != nil {
			return errors.Wrap(err, "reset cache")
		}
		return nil
	}

	// Prune via API.
	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}
	bkClient, err := app.getBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "prune new buildkitd client")
	}
	defer bkClient.Close()
	var opts []client.PruneOption

	if app.pruneAll {
		opts = append(opts, client.PruneAll)
	}

	if app.pruneKeepDuration > 0 || app.pruneTargetSize > 0 {
		opts = append(opts, client.WithKeepOpt(app.pruneKeepDuration, int64(app.pruneTargetSize)))
	}

	ch := make(chan client.UsageInfo, 1)
	eg, ctx := errgroup.WithContext(cliCtx.Context)
	eg.Go(func() error {
		err = bkClient.Prune(ctx, ch, opts...)
		if err != nil {
			return errors.Wrap(err, "buildkit prune")
		}
		close(ch)
		return nil
	})

	total := uint64(0)
	eg.Go(func() error {
		for {
			select {
			case usageInfo, ok := <-ch:
				if !ok {
					return nil
				}
				app.console.Printf("%s\t%s\n", usageInfo.ID, humanize.Bytes(uint64(usageInfo.Size)))
				total += uint64(usageInfo.Size)
			case <-ctx.Done():
				return nil
			}
		}
	})
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "err group")
	}
	app.console.Printf("Freed %s\n", humanize.Bytes(total))
	return nil
}
