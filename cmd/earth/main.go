package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof" // enable pprof handlers on net/http listener
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/earthly/earthly/autocomplete"
	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cleanup"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	debuggercommon "github.com/earthly/earthly/debugger/common"
	"github.com/earthly/earthly/debugger/terminal"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/earthfile2llb/variables"
	"github.com/earthly/earthly/logging"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/moby/buildkit/client"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var dotEnvPath = ".env"

type earthApp struct {
	cliApp    *cli.App
	console   conslogging.ConsoleLogger
	cfg       *config.Config
	sessionID string
	cliFlags
}

type cliFlags struct {
	buildArgs            cli.StringSlice
	secrets              cli.StringSlice
	artifactMode         bool
	imageMode            bool
	pull                 bool
	push                 bool
	noOutput             bool
	noCache              bool
	pruneAll             bool
	pruneReset           bool
	buildkitdSettings    buildkitd.Settings
	allowPrivileged      bool
	enableProfiler       bool
	buildkitHost         string
	buildkitdImage       string
	remoteCache          string
	configPath           string
	gitUsernameOverride  string
	gitPasswordOverride  string
	interactiveDebugging bool
	sshAuthSock          string
	homebrewSource       string
}

var (
	// DefaultBuildkitdImage is the default buildkitd image to use.
	DefaultBuildkitdImage string

	// Version is the version of this CLI app.
	Version string

	// GitSha contains the git sha used to build this app
	GitSha string
)

func profhandler() {
	addr := "127.0.0.1:6060"
	fmt.Printf("listening for pprof on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

func main() {
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
		for {
			select {
			case sig := <-c:
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
		}
	}()

	// Load .env into current global env's. This is mainly for applying Earthly settings.
	// Separate call is made for build args and secrets.
	if fileExists(dotEnvPath) {
		err := godotenv.Load(dotEnvPath)
		if err != nil {
			fmt.Printf("Error loading dot-env file %s: %s\n", dotEnvPath, err.Error())
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

	app := newEarthApp(ctx, conslogging.Current(colorMode))
	app.autoComplete()

	// Set up file-based logging.
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error looking up current user: %s\n", err.Error())
		os.Exit(1)
	}
	logDir := filepath.Join(homeDir, ".earthly")
	logFile := filepath.Join(logDir, "earth.log")
	err = os.MkdirAll(logDir, 0755)
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

	os.Exit(app.run(ctx, os.Args))
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
			Name:        "pull",
			EnvVars:     []string{"EARTHLY_PULL"},
			Usage:       "Force pull any referenced Docker images",
			Destination: &app.pull,
			Hidden:      true, // Experimental.
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
			Value:       defaultConfigPath(),
			EnvVars:     []string{"EARTHLY_CONFIG"},
			Usage:       "Path to config file for",
			Destination: &app.configPath,
		},
		&cli.StringFlag{
			Name:        "ssh-auth-sock",
			Value:       os.Getenv("SSH_AUTH_SOCK"),
			EnvVars:     []string{"EARTHLY_SSH_AUTH_SOCK"},
			Usage:       "The SSH auth socket to use for ssh-agent forwarding",
			Destination: &app.sshAuthSock,
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
		&cli.BoolFlag{
			Name:        "profiler",
			EnvVars:     []string{"EARTHLY_PROFILER"},
			Usage:       "Enable the profiler",
			Destination: &app.enableProfiler,
			Hidden:      true, // for use in dev debugging
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
			Usage:       "enable verbose logging of the earthly-buildkitd container",
			Destination: &app.buildkitdSettings.Debug,
		},
	}

	app.cliApp.Commands = []*cli.Command{
		{
			Name:        "bootstrap",
			Usage:       "Bootstraps earth bash autocompletion",
			Description: "Performs initial earth bootstrapping for bash autocompletion",
			Hidden:      false,
			Action:      app.actionBootstrap,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "source",
					Usage:       "output source file (for use in homebrew install)",
					Hidden:      true, // only meant for use with homebrew formula
					Destination: &app.homebrewSource,
				},
			},
		},
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

	app.cliApp.Before = app.before
	return app
}

func (app *earthApp) before(context *cli.Context) error {
	if app.enableProfiler {
		go profhandler()
	}

	if context.IsSet("config") {
		app.console.Printf("loading config values from %q\n", app.configPath)
	}

	yamlData, err := ioutil.ReadFile(app.configPath)
	if os.IsNotExist(err) && !context.IsSet("config") {
		yamlData = []byte{}
	} else if err != nil {
		return errors.Wrapf(err, "failed to read from %s", app.configPath)
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

	gitConfig, gitCredentials, err := config.CreateGitConfig(app.cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to create git config from %s", app.configPath)
	}

	// command line option overrides the config which overrides the default value
	if !context.IsSet("buildkit-image") && app.cfg.Global.BuildkitImage != "" {
		app.buildkitdImage = app.cfg.Global.BuildkitImage
	}

	if runtime.GOOS == "darwin" {
		// on darwin buildkit is running inside a docker container and must reference this sock instead
		app.buildkitdSettings.SSHAuthSock = "/run/host-services/ssh-auth.sock"
	} else {
		app.buildkitdSettings.SSHAuthSock = app.sshAuthSock
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

	if !dirExists(app.cfg.Global.RunPath) {
		err := os.MkdirAll(app.cfg.Global.RunPath, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to create run directory %s", app.cfg.Global.RunPath)
		}
	}

	app.buildkitdSettings.DebuggerPort = app.cfg.Global.DebuggerPort
	app.buildkitdSettings.RunDir = app.cfg.Global.RunPath
	app.buildkitdSettings.GitConfig = gitConfig
	app.buildkitdSettings.GitCredentials = gitCredentials
	return nil
}

func (app *earthApp) processDeprecatedCommandOptions(context *cli.Context, cfg *config.Config) error {
	if cfg.Global.CachePath != "" {
		app.console.Warnf("Warning: the setting cache_path is now obsolete and will be ignored")
	}

	// command line overrides the config file
	if app.gitUsernameOverride != "" || app.gitPasswordOverride != "" {
		app.console.Warnf("Warning: the --git-username and --git-password command flags are deprecated and are now configured in the ~/.earthly/config.yml file under the git section; see https://docs.earthly.dev/earth-config for reference.\n")
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
		app.console.Warnf("Warning: the --git-url-instead-of command flag is deprecated and is now configured in the ~/.earthly/config.yml file under the git global url_instead_of setting; see https://docs.earthly.dev/earth-config for reference.\n")
	} else {
		if gitGlobal, ok := cfg.Git["global"]; ok {
			if gitGlobal.GitURLInsteadOf != "" {
				app.buildkitdSettings.GitURLInsteadOf = gitGlobal.GitURLInsteadOf
			}
		}
	}

	if context.IsSet("no-loop-device") {
		app.console.Warnf("Warning: the --no-loop-device command flag is deprecated and is now configured in the ~/.earthly/config.yml file under the no_loop_device setting; see https://docs.earthly.dev/earth-config for reference.\n")
	} else {
		app.buildkitdSettings.DisableLoopDevice = cfg.Global.DisableLoopDevice
	}

	if context.IsSet("buildkit-cache-size-mb") {
		app.console.Warnf("Warning: the --buildkit-cache-size-mb command flag is deprecated and is now configured in the ~/.earthly/config.yml file under the buildkit_cache_size setting; see https://docs.earthly.dev/earth-config for reference.\n")
	} else {
		app.buildkitdSettings.CacheSizeMb = cfg.Global.BuildkitCacheSizeMb
	}

	return nil
}

// to enable autocomplete, enter
// complete -o nospace -C "/path/to/earth" earth
func (app *earthApp) autoComplete() {
	_, found := os.LookupEnv("COMP_LINE")
	if !found {
		return
	}

	err := app.autoCompleteImp()
	if err != nil {
		errToLog := err
		homeDir, err := os.UserHomeDir()
		if err != nil {
			os.Exit(1)
		}
		logDir := filepath.Join(homeDir, ".earthly")
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

func (app *earthApp) autoCompleteImp() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered panic in autocomplete %s: %s", r, debug.Stack())
		}
	}()

	compLine := os.Getenv("COMP_LINE")   // full command line
	compPoint := os.Getenv("COMP_POINT") // where the cursor is

	compPointInt, err := strconv.ParseUint(compPoint, 10, 64)
	if err != nil {
		return err
	}

	flags := []string{}
	for _, f := range app.cliApp.Flags {
		for _, n := range f.Names() {
			if len(n) > 1 {
				flags = append(flags, n)
			}
		}
	}

	commands := []string{}
	for _, cmd := range app.cliApp.Commands {
		commands = append(commands, cmd.Name)
	}

	potentials, err := autocomplete.GetPotentials(compLine, int(compPointInt), flags, commands)
	if err != nil {
		return err
	}
	for _, p := range potentials {
		fmt.Printf("%s\n", p)
	}

	return err
}

const bashCompleteEntry = "complete -o nospace -C '/usr/local/bin/earth' earth\n"

func (app *earthApp) insertBashCompleteEntry() error {
	var path string
	if runtime.GOOS == "darwin" {
		path = "/usr/local/etc/bash_completion.d/earth"
	} else {
		path = "/usr/share/bash-completion/completions/earth"
	}
	dirPath := filepath.Dir(path)

	if !dirExists(dirPath) {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s does not exist\n", dirPath)
		return nil // bash-completion isn't available, silently fail.
	}

	if fileExists(path) {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(bashCompleteEntry))
	return err
}

func (app *earthApp) deleteZcompdump() error {
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
	files, err := ioutil.ReadDir(homeDir)
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

const zshCompleteEntry = `#compdef _earth earth

function _earth {
    autoload -Uz bashcompinit
    bashcompinit
    complete -o nospace -C '/usr/local/bin/earth' earth
}
`

// If debugging this, it might be required to run `rm ~/.zcompdump*` to remove the cache
func (app *earthApp) insertZSHCompleteEntry() error {
	// should be the same on linux and macOS
	path := "/usr/local/share/zsh/site-functions/_earth"
	dirPath := filepath.Dir(path)

	if !dirExists(dirPath) {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable zsh-completion: %s does not exist\n", dirPath)
		return nil // zsh-completion isn't available, silently fail.
	}

	if fileExists(path) {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(zshCompleteEntry))
	if err != nil {
		return err
	}

	return app.deleteZcompdump()
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
		if strings.Contains(err.Error(), "security.insecure is not allowed") {
			app.console.Warnf("Error: --allow-privileged (-P) flag is required\n")
		} else if strings.Contains(err.Error(), "failed to fetch remote") {
			app.console.Warnf("Error: %v\n", err)
			app.console.Printf(
				"Check your git auth settings.\n" +
					"Did you ssh-add today? Need to configure ~/.earthly/config.yml?\n" +
					"For more information see https://docs.earthly.dev/guides/auth\n")
		} else {
			app.console.Warnf("Error: %v\n", err)
		}
		if errors.Is(err, context.Canceled) {
			return 2
		}
		return 1
	}
	return 0
}

func (app *earthApp) actionBootstrap(c *cli.Context) error {

	switch app.homebrewSource {
	case "bash":
		fmt.Printf(bashCompleteEntry)
		return nil
	case "zsh":
		fmt.Printf(zshCompleteEntry)
		return nil
	case "":
		break
	default:
		return fmt.Errorf("unhandled source %q", app.homebrewSource)
	}

	err := app.insertBashCompleteEntry()
	if err != nil {
		return err
	}

	err = app.insertZSHCompleteEntry()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Bootstrapping successful; you may have to restart your shell for autocomplete to get initialized (e.g. run \"exec $SHELL\")\n")

	return nil
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
		// Use twice the restart timeout for reset operations
		// (needs extra time to also remove the files).
		opTimeout := 2 * time.Duration(app.cfg.Global.BuildkitRestartTimeoutS) * time.Second
		err := buildkitd.ResetCache(
			c.Context, app.console, app.buildkitdImage, app.buildkitdSettings,
			opTimeout)
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
	sockName := fmt.Sprintf("debugger.sock.%d", time.Now().UnixNano())

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

	dotEnvMap := make(map[string]string)
	if fileExists(dotEnvPath) {
		dotEnvMap, err = godotenv.Read(dotEnvPath)
		if err != nil {
			return errors.Wrapf(err, "read %s", dotEnvPath)
		}
	}
	secretsMap, err := processSecrets(secrets, dotEnvMap)
	if err != nil {
		return err
	}

	debuggerSettings := debuggercommon.DebuggerSettings{
		DebugLevelLogging: app.buildkitdSettings.Debug,
		Enabled:           app.interactiveDebugging,
		SockPath:          fmt.Sprintf("/run/earthly/%s", sockName),
		Term:              os.Getenv("TERM"),
	}

	debuggerSettingsData, err := json.Marshal(&debuggerSettings)
	if err != nil {
		return errors.Wrap(err, "debugger settings json marshal")
	}
	secretsMap[debuggercommon.DebuggerSettingsSecretsKey] = debuggerSettingsData

	attachables := []session.Attachable{
		secretsprovider.FromMap(secretsMap),
		authprovider.NewDockerAuthProvider(os.Stderr),
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
	b, err := builder.NewBuilder(
		c.Context, bkClient, app.console, attachables, enttlmnts, app.noCache, app.remoteCache)
	if err != nil {
		return errors.Wrap(err, "new builder")
	}

	if app.interactiveDebugging {
		go terminal.ConnectTerm(c.Context, fmt.Sprintf("127.0.0.1:%d", app.buildkitdSettings.DebuggerPort))
	}

	varCollection, err := variables.ParseCommandLineBuildArgs(app.buildArgs.Value(), dotEnvMap)
	if err != nil {
		return errors.Wrap(err, "parse build args")
	}
	imageResolveMode := llb.ResolveModePreferLocal
	if app.pull {
		imageResolveMode = llb.ResolveModeForcePull
	}
	cleanCollection := cleanup.NewCollection()
	defer cleanCollection.Close()
	mts, err := earthfile2llb.Earthfile2LLB(
		c.Context, target, earthfile2llb.ConvertOpt{
			Resolver:           resolver,
			ImageResolveMode:   imageResolveMode,
			DockerBuilderFun:   b.MakeImageAsTarBuilderFun(),
			ArtifactBuilderFun: b.MakeArtifactBuilderFun(),
			CleanCollection:    cleanCollection,
			VarCollection:      varCollection,
		})
	if err != nil {
		return err
	}

	opts := builder.BuildOpt{
		PrintSuccess: true,
		Push:         app.push,
		NoOutput:     app.noOutput,
	}
	if app.imageMode {
		err = b.BuildOnlyImages(c.Context, mts, opts)
		if err != nil {
			return err
		}
	} else if app.artifactMode {
		err = b.BuildOnlyArtifact(c.Context, mts, artifact, destPath, opts)
		if err != nil {
			return err
		}
	} else {
		err = b.Build(c.Context, mts, opts)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *earthApp) newBuildkitdClient(ctx context.Context, opts ...client.ClientOpt) (*client.Client, error) {
	if app.buildkitHost == "" {
		// Start our own.
		opTimeout := time.Duration(app.cfg.Global.BuildkitRestartTimeoutS) * time.Second
		bkClient, err := buildkitd.NewClient(
			ctx, app.console, app.buildkitdImage, app.buildkitdSettings, opTimeout)
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

func processSecrets(secrets []string, dotEnvMap map[string]string) (map[string][]byte, error) {
	finalSecrets := make(map[string][]byte)
	for k, v := range dotEnvMap {
		finalSecrets[k] = []byte(v)
	}
	for _, secret := range secrets {
		parts := strings.SplitN(secret, "=", 2)
		if len(parts) == 2 {
			// Already set.
			finalSecrets[parts[0]] = []byte(parts[1])
		} else {
			// Not set. Use environment to fetch it.
			value, found := os.LookupEnv(secret)
			if !found {
				return nil, fmt.Errorf("env var %s not set", secret)
			}
			finalSecrets[secret] = []byte(value)
		}
	}
	return finalSecrets, nil
}

func defaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	oldConfig := filepath.Join(homeDir, ".earthly", "config.yaml")
	newConfig := filepath.Join(homeDir, ".earthly", "config.yml")
	if fileExists(oldConfig) && !fileExists(newConfig) {
		return oldConfig
	}
	return newConfig
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
