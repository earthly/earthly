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
	"github.com/joho/godotenv"
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
	"github.com/earthly/earthly/util/flagutil"
	"github.com/earthly/earthly/util/termutil"
	"github.com/earthly/earthly/variables"
)

func (app *earthlyApp) rootCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "build",
			Usage:       "Build an Earthly target",
			Description: "Build an Earthly target.",
			Action:      app.actionBuild,
			Flags:       app.buildFlags(),
			Hidden:      true, // Meant to be used mainly for help output.
		},
		{
			Name:        "bootstrap",
			Usage:       "Bootstraps earthly installation including buildkit image download and optionally shell autocompletion",
			UsageText:   "earthly [options] bootstrap [--no-buildkit, --with-autocomplete, --certs-hostname]",
			Description: "Bootstraps earthly installation including buildkit image download and optionally shell autocompletion.",
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
					Usage:       "Skips setting up the BuildKit container",
					Destination: &app.bootstrapNoBuildkit,
				},
				&cli.BoolFlag{
					Name:        "with-autocomplete",
					Usage:       "Install shell autocompletions during bootstrap",
					Destination: &app.bootstrapWithAutocomplete,
				},
				&cli.StringFlag{
					Name:        "certs-hostname",
					Usage:       "Hostname to generate certificates for",
					EnvVars:     []string{"EARTHLY_CERTS_HOSTNAME"},
					Value:       "localhost",
					Destination: &app.certsHostName,
				},
			},
		},
		{
			Name:        "docker-build",
			Usage:       "*beta* Build a Dockerfile without an Earthfile",
			UsageText:   "earthly [options] docker-build [--dockerfile <dockerfile-path>] [--tag=<image-tag>] [--target=<target-name>] [--platform <platform1[,platform2,...]>] <build-context-dir> [--arg1=arg-value]",
			Description: "*beta* Builds a Dockerfile without an Earthfile.",
			Action:      app.actionDockerBuild,
			Flags: append(app.buildFlags(),
				&cli.StringFlag{
					Name:        "dockerfile",
					Aliases:     []string{"f"},
					EnvVars:     []string{"EARTHLY_DOCKER_FILE"},
					Usage:       "Path to dockerfile input",
					Value:       "Dockerfile",
					Destination: &app.dockerfilePath,
				},
				&cli.StringSliceFlag{
					Name:        "tag",
					Aliases:     []string{"t"},
					EnvVars:     []string{"EARTHLY_DOCKER_TAGS"},
					Usage:       "Name and tag for the built image; formatted as 'name:tag'",
					Destination: &app.dockerTags,
				},
				&cli.StringFlag{
					Name:        "target",
					EnvVars:     []string{"EARTHLY_DOCKER_TARGET"},
					Usage:       "The docker target to build in the specified dockerfile",
					Destination: &app.dockerTarget,
				},
			),
		},
		{
			Name:        "docker2earthly",
			Usage:       "Convert a Dockerfile into Earthfile",
			Description: "Converts an existing dockerfile into an Earthfile.",
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
			Usage:       "Create or manage your Earthly orgs",
			Description: "Create or manage your Earthly orgs.",
			Subcommands: app.orgCmds(),
		},
		{
			Name:        "doc",
			Usage:       "Document targets from an Earthfile",
			UsageText:   "earthly [options] doc [<project-ref>[+<target-ref>]]",
			Description: "Document targets from an Earthfile by reading in line comments.",
			Action:      app.actionDocumentTarget,
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
			Name:        "ls",
			Usage:       "List targets from an Earthfile",
			UsageText:   "earthly [options] ls [<project-ref>]",
			Description: "List targets from an Earthfile.",
			Action:      app.actionListTargets,
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
			Usage:       "Create or manage an Earthly account",
			Description: "Create or manage an Earthly account.",
			Subcommands: app.accountCmds(),
		},
		{
			Name:        "debug",
			Usage:       "Print debug information about an Earthfile",
			Description: "Print debug information about an Earthfile.",
			ArgsUsage:   "[<path>]",
			Hidden:      true, // Dev purposes only.
			Subcommands: app.debugCmds(),
		},
		{
			Name:  "prune",
			Usage: "Prune Earthly build cache",
			Description: `Prune Earthly build cache in one of two forms.

Standard Form:
	Issues a prune command on the BuildKit daemon.
Reset Form:
	Restarts the BuildKit daemon and instructs it to complete delete the cache
	directory on startup.`,
			Action: app.actionPrune,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "all",
					Aliases:     []string{"a"},
					EnvVars:     []string{"EARTHLY_PRUNE_ALL"},
					Usage:       "Prune all cache via BuildKit daemon",
					Destination: &app.pruneAll,
				},
				&cli.BoolFlag{
					Name:    "reset",
					EnvVars: []string{"EARTHLY_PRUNE_RESET"},
					Usage: `Reset cache entirely by restarting BuildKit daemon and wiping cache dir.
			This option is not available when using satellites.`,
					Destination: &app.pruneReset,
				},
				&cli.DurationFlag{
					Name: "age",
					Usage: `Prune cache older than the specified duration passed in as a string; 
					duration is specified with an integer value followed by a m, h, or d suffix which represents minutes, hours, or days respectively, e.g. 24h, or 1d`,
					Destination: &app.pruneKeepDuration,
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
			UsageText: `Examples of common settings:
			
			Set your cache size:
			
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
			
				config git "{example: {pattern: 'example.com/([^/]+)', substitute: 'ssh://git@example.com:2222/var/git/repos/\$1.git', auth: ssh}}`,
			Description: `This command takes both a path and a value. It then sets them in your configuration file.

			As the configuration file is YAML, the key must be a valid key within the file. You can specify sub-keys by using "." to separate levels.
			If the sub-key you wish to use has a "." in it, you can quote that subsection, like this: git."github.com".
		
			Values must be valid YAML, and also be deserializable into the key you wish to assign them to.
			This means you can set higher level objects using a compact style, or single values with simple values.
		
			Only one key/value can be set per invocation.
		
			To get help with a specific key, do "config [key] --help". Or, visit https://docs.earthly.dev/earthly-config for more details.`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "dry-run",
					Usage:       "Print the changed config file to the console instead of writing it to file",
					Destination: &app.configDryRun,
				},
			},
		},
		{
			Name:      "satellite",
			Aliases:   []string{"satellites", "sat"},
			Usage:     "Create and manage Earthly Satellites",
			UsageText: "earthly satellite (launch|ls|inspect|select|unselect|rm)",
			Description: `Launch and use a Satellite runner as remote backend for Earthly builds.

- Read more about satellites here: https://docs.earthly.dev/earthly-cloud/satellites 
- Sign up for satellites here: https://cloud.earthly.dev/login

Satellites can be used to share cache between multiple builds and users,
as well as run builds in native architectures independent of where the Earthly client is invoked.`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the organization the satellite belongs to",
					Required:    false,
					Destination: &app.orgName,
				},
			},
			Subcommands: app.satelliteCmds(),
		},
		{
			Name:    "project",
			Aliases: []string{"projects"},
			Description: `Manage Earthly projects which are shared resources of Earthly orgs. 

Within Earthly projects users can be invited and granted different access levels including: read, read+secrets, write, and admin.`,
			Usage:     "Manage Earthly projects",
			UsageText: "earthly project (ls|rm|create|member)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the Earthly organization to which the Earthly project belongs",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					Aliases:     []string{"p"},
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The Earthly project to act on",
					Required:    false,
					Destination: &app.projectName,
				},
			},
			Subcommands: app.projectCmds(),
		},
		{
			Name:        "secret",
			Aliases:     []string{"secrets"},
			Description: "*beta* Manage cloud secrets.",
			Usage:       "*beta* Manage cloud secrets",
			UsageText:   "earthly [options] secrets [--org <organization-name>, --project <project>] (set|get|ls|rm|migrate|permission)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store secrets",
					Required:    false,
					Destination: &app.projectName,
				},
			},
			Subcommands: app.secretCmds(),
		},
		{
			Name:        "registry",
			Aliases:     []string{"registries"},
			Description: "*beta* Manage registry access.",
			Usage:       "*beta* Manage registry access",
			UsageText:   "earthly [options] registry [--org <organization-name>, --project <project>] (setup|list|remove) [<flags>]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store registry credentials",
					Required:    false,
					Destination: &app.projectName,
				},
			},
			Subcommands: app.registryCmds(),
		},
		{
			Name:      "web",
			Usage:     "*beta* Access the web UI via your default browser and print the url",
			UsageText: "earthly web (--provider=github)",
			Description: `*beta* Prints a url for entering the CI application and attempts to open your default browser with that url.

	If the provider argument is given the CI application will automatically begin an OAuth flow with the given provider.
	If you are logged into the CLI the url will contain a token used to link your OAuth credentials to your Earthly user.`,
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
		{
			Name:        "init",
			Description: "*experimental* Initialize a project.",
			Usage:       "*experimental* Initialize an Earthfile for the current project",
			Action:      app.actionInit,
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
		return errors.Wrapf(err, "failed to check if %q exists", earthPath)
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

func (app *earthlyApp) actionDockerBuild(cliCtx *cli.Context) error {
	app.commandName = "docker-build"

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

	argMap, err := godotenv.Read(app.argFile)
	if err != nil && (cliCtx.IsSet(argFileFlag) || !errors.Is(err, os.ErrNotExist)) {
		return errors.Wrapf(err, "read %q", app.argFile)
	}

	buildArgs, err := app.combineVariables(argMap, flagArgs)
	if err != nil {
		return errors.Wrapf(err, "combining build args")
	}

	platforms := flagutil.SplitFlagString(app.platformsStr)
	content, err := docker2earthly.GenerateEarthfile(buildContextPath, app.dockerfilePath, app.dockerTags.Value(), buildArgs.Sorted(), platforms, app.dockerTarget)
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
	app.imageMode = false
	app.artifactMode = false
	app.platformsStr = cli.StringSlice{}
	app.dockerTarget = ""
	app.dockerfilePath = ""
	app.dockerTags = cli.StringSlice{}

	nonFlagArgs = []string{tempDir + "+build"}
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
	resolver := buildcontext.NewResolver(nil, gitLookup, app.console, "", app.gitBranchOverride, app.gitLFSPullInclude, 0, "")
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
