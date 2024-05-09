package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	_ "net/http/pprof" // enable pprof handlers on net/http listener
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/earthly/earthly/internal/version"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/earthly/earthly/cmd/earthly/app"
	"github.com/earthly/earthly/cmd/earthly/base"
	"github.com/earthly/earthly/cmd/earthly/common"
	eFlag "github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/cmd/earthly/subcmd"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/envutil"
	"github.com/earthly/earthly/util/syncutil"
)

// These vars are set by ldflags
var (
	// Version is the version of this CLI app.
	Version string
	// GitSha contains the git sha used to build this app
	GitSha string
	// BuiltBy contains information on which build-system was used (e.g. official earthly binaries, homebrew, etc)
	BuiltBy string

	// DefaultBuildkitdImage is the default buildkitd image to use.
	DefaultBuildkitdImage string

	// DefaultInstallationName is the name included in the various earthly global resources on the system,
	// such as the ~/.earthly dir name, the buildkitd container name, the docker volume name, etc.
	// This should be set to "earthly" for official releases.
	DefaultInstallationName string

	Foo string
)

func setExportableVars() {
	version.Version = Version
	version.GitSha = GitSha
	version.BuiltBy = BuiltBy
}

func main() {
	setExportableVars()
	startTime := time.Now()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()
	lastSignal := &syncutil.Signal{}
	go func() {
		for sig := range sigChan {
			if lastSignal.Get() != nil {
				// This is the second time we have received a signal. Quit immediately.
				fmt.Printf("Received second signal %s. Forcing exit.\n", sig.String())
				os.Exit(9)
			}
			lastSignal.Set(sig)
			cancel()
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
	envFile := eFlag.DefaultEnvFile
	envFileOverride := false
	if envFileFromEnv, ok := os.LookupEnv("EARTHLY_ENV_FILE"); ok {
		envFile = envFileFromEnv
		envFileOverride = true
	}
	envFileFromArgOK := true
	flagSet := flag.NewFlagSet(common.GetBinaryName(), flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	cli := base.NewCLI(conslogging.ConsoleLogger{},
		base.WithVersion(Version),
		base.WithGitSHA(GitSha),
		base.WithBuiltBy(BuiltBy),
		base.WithDefaultBuildkitdImage(DefaultBuildkitdImage),
		base.WithDefaultInstallationName(DefaultInstallationName),
	)
	buildApp := subcmd.NewBuild(cli)
	rootApp := subcmd.NewRoot(cli, buildApp)

	for _, f := range cli.Flags().RootFlags(DefaultInstallationName, DefaultBuildkitdImage) {
		if err := f.Apply(flagSet); err != nil {
			envFileFromArgOK = false
			break
		}
	}

	if envFileFromArgOK {
		if err := flagSet.Parse(os.Args[1:]); err == nil {
			if envFileFlag := flagSet.Lookup(eFlag.EnvFileFlag); envFileFlag != nil {
				envFile = envFileFlag.Value.String()
				envFileOverride = envFile != eFlag.DefaultEnvFile // flag lib doesn't expose if a value was set or not
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
	if envutil.IsTrue("FORCE_COLOR") {
		colorMode = conslogging.ForceColor
		color.NoColor = false
	}
	if envutil.IsTrue("NO_COLOR") {
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

	if envutil.IsTrue("EARTHLY_FULL_TARGET") {
		padding = conslogging.NoPadding
	}
	logging := conslogging.Current(colorMode, padding, conslogging.Info, cli.Flags().GithubAnnotations)

	cli.SetConsole(logging)
	earthly := app.NewEarthlyApp(cli, rootApp, buildApp, ctx)
	exitCode := earthly.Run(ctx, logging, startTime, lastSignal)
	os.Exit(exitCode)
}
