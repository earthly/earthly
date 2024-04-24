package app

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/earthly/earthly/cmd/earthly/subcmd"

	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	logbussetup "github.com/earthly/earthly/logbus/setup"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/envutil"
	"github.com/earthly/earthly/util/execstatssummary"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *EarthlyApp) before(cliCtx *cli.Context) error {
	if app.BaseCLI.Flags().EnableProfiler {
		go profhandler()
	}

	if app.BaseCLI.Flags().InstallationName != "" {
		if !cliCtx.IsSet("config") {
			app.BaseCLI.Flags().ConfigPath = defaultConfigPath(app.BaseCLI.Flags().InstallationName)
		}
		if !cliCtx.IsSet("buildkit-container-name") {
			app.BaseCLI.Flags().ContainerName = fmt.Sprintf("%s-buildkitd", app.BaseCLI.Flags().InstallationName)
		}
		if !cliCtx.IsSet("buildkit-volume-name") {
			app.BaseCLI.Flags().BuildkitdSettings.VolumeName = fmt.Sprintf("%s-cache", app.BaseCLI.Flags().InstallationName)
		}
	}
	if app.BaseCLI.Flags().Debug {
		app.BaseCLI.SetConsole(app.BaseCLI.Console().WithLogLevel(conslogging.Debug))
	} else if app.BaseCLI.Flags().Verbose {
		app.BaseCLI.SetConsole(app.BaseCLI.Console().WithLogLevel(conslogging.Verbose))
	}
	if app.BaseCLI.Flags().LogstreamUpload {
		app.BaseCLI.Flags().Logstream = true
	}
	if app.BaseCLI.Flags().Logstream {
		app.BaseCLI.SetConsole(app.BaseCLI.Console().WithPrefixWriter(app.BaseCLI.Logbus().Run().Generic()))
		if app.BaseCLI.Flags().BuildID == "" {
			app.BaseCLI.Flags().BuildID = uuid.NewString()
		}
		var execStatsTracker *execstatssummary.Tracker
		if app.BaseCLI.Flags().ExecStatsSummary != "" {
			execStatsTracker = execstatssummary.NewTracker(app.BaseCLI.Flags().ExecStatsSummary)
		}
		disableOngoingUpdates := !app.BaseCLI.Flags().Logstream || app.BaseCLI.Flags().InteractiveDebugging
		forceColor := envutil.IsTrue("FORCE_COLOR")
		noColor := envutil.IsTrue("NO_COLOR")
		var err error
		newSetup, err := logbussetup.New(
			cliCtx.Context, app.BaseCLI.Logbus(), app.BaseCLI.Flags().Debug, app.BaseCLI.Flags().Verbose, app.BaseCLI.Flags().DisplayExecStats, forceColor, noColor,
			disableOngoingUpdates, app.BaseCLI.Flags().LogstreamDebugFile, app.BaseCLI.Flags().BuildID, execStatsTracker, app.BaseCLI.Flags().IsGithubAnnotations)
		app.BaseCLI.SetLogbusSetup(newSetup)
		if err != nil {
			return errors.Wrap(err, "logbus setup")
		}
	}

	if cliCtx.IsSet("config") {
		app.BaseCLI.Console().Printf("loading config values from %q\n", app.BaseCLI.Flags().ConfigPath)
	}

	var yamlData []byte
	if app.BaseCLI.Flags().ConfigPath != "" {
		var err error
		yamlData, err = config.ReadConfigFile(app.BaseCLI.Flags().ConfigPath)
		if err != nil {
			if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
				return errors.Wrapf(err, "read config")
			}
		}
	}

	cfg, err := config.ParseYAML(yamlData, app.BaseCLI.Flags().InstallationName)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", app.BaseCLI.Flags().ConfigPath)
	}
	app.BaseCLI.SetCfg(&cfg)

	err = app.processDeprecatedCommandOptions(cliCtx, app.BaseCLI.Cfg())
	if err != nil {
		return err
	}

	err = app.parseFrontend(cliCtx, app.BaseCLI.Cfg())
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

	if !isBootstrapCmd && !cliutil.IsBootstrapped(app.BaseCLI.Flags().InstallationName) {
		app.BaseCLI.Flags().BootstrapNoBuildkit = true // Docker may not be available, for instance... like our integration tests.
		newBootstrap := subcmd.NewBootstrap(app.BaseCLI)
		err = newBootstrap.Action(cliCtx)
		if err != nil {
			return errors.Wrap(err, "bootstrap unbootstrclied installation")
		}
	}

	if !cliCtx.IsSet("org") {
		app.BaseCLI.Flags().OrgName = app.BaseCLI.Cfg().Global.Org
	}
	return nil
}

func (app *EarthlyApp) parseFrontend(cliCtx *cli.Context, cfg *config.Config) error {
	console := app.BaseCLI.Console().WithPrefix("frontend")
	feCfg := &containerutil.FrontendConfig{
		BuildkitHostCLIValue:       app.BaseCLI.Flags().BuildkitHost,
		BuildkitHostFileValue:      app.BaseCLI.Cfg().Global.BuildkitHost,
		LocalRegistryHostFileValue: app.BaseCLI.Cfg().Global.LocalRegistryHost,
		InstallationName:           app.BaseCLI.Flags().InstallationName,
		DefaultPort:                8372 + config.PortOffset(app.BaseCLI.Flags().InstallationName),
		Console:                    console,
	}
	fe, err := containerutil.FrontendForSetting(cliCtx.Context, app.BaseCLI.Cfg().Global.ContainerFrontend, feCfg)
	if err != nil {
		origErr := err
		stub, err := containerutil.NewStubFrontend(cliCtx.Context, feCfg)
		if err != nil {
			return errors.Wrap(err, "failed stub frontend initialization")
		}
		app.BaseCLI.Flags().ContainerFrontend = stub

		if !app.BaseCLI.Flags().Verbose {
			console.Printf("Unable to detect Docker or Podman. Use --verbose to see details (or errors)\n")
		}
		console.VerbosePrintf("%s frontend initialization failed due to %s", app.BaseCLI.Cfg().Global.ContainerFrontend, origErr.Error())
		return nil
	}

	console.VerbosePrintf("%s frontend initialized.\n", fe.Config().Setting)
	app.BaseCLI.Flags().ContainerFrontend = fe

	// These URLs were calculated relative to the configured frontend. In the
	// case of an automatically detected frontend, they are calculated according
	// to the first selected one in order of precedence.
	buildkitURLs := app.BaseCLI.Flags().ContainerFrontend.Config().FrontendURLs
	app.BaseCLI.Flags().BuildkitHost = buildkitURLs.BuildkitHost.String()
	app.BaseCLI.Flags().LocalRegistryHost = buildkitURLs.LocalRegistryHost.String()

	return nil
}

func (app *EarthlyApp) processDeprecatedCommandOptions(cliCtx *cli.Context, cfg *config.Config) error {
	app.warnIfEarth()

	if cfg.Global.CachePath != "" {
		app.BaseCLI.Console().Warnf("Warning: the setting cache_path is now obsolete and will be ignored")
	}

	if app.BaseCLI.Flags().ConversionParallelism != 0 {
		app.BaseCLI.Console().Warnf("Warning: --conversion-parallelism and EARTHLY_CONVERSION_PARALLELISM is obsolete, please use 'earthly config global.conversion_parallelism <parallelism>' instead")
	}

	// command line overrides the config file
	if app.BaseCLI.Flags().GitUsernameOverride != "" || app.BaseCLI.Flags().GitPasswordOverride != "" {
		app.BaseCLI.Console().Warnf("Warning: the --git-username and --git-password command flags are deprecated and are now configured in the ~/.earthly/config.yml file under the git section; see https://docs.earthly.dev/earthly-config for reference.\n")
		if _, ok := cfg.Git["github.com"]; !ok {
			cfg.Git["github.com"] = config.GitConfig{}
		}
		if _, ok := cfg.Git["gitlab.com"]; !ok {
			cfg.Git["gitlab.com"] = config.GitConfig{}
		}

		for k, v := range cfg.Git {
			v.Auth = "https"
			if app.BaseCLI.Flags().GitUsernameOverride != "" {
				v.User = app.BaseCLI.Flags().GitUsernameOverride
			}
			if app.BaseCLI.Flags().GitPasswordOverride != "" {
				v.Password = app.BaseCLI.Flags().GitPasswordOverride
			}
			cfg.Git[k] = v
		}
	}

	if cfg.Satellite.Org != "" {
		if cfg.Global.Org != "" {
			app.BaseCLI.Console().Warnf("Two default organizations were specified.\n" +
				"You can remove the deprecated value by running 'earthly config satellite.org \"\"'\n" +
				"Earthly will use the global value.")
		} else {
			app.BaseCLI.Console().Warnf("Auto-selecting the default org will no longer be supported in the future.\n" +
				"You can select a default org using the command 'earthly org select',\n" +
				"or otherwise specify an org using the --org flag or EARTHLY_ORG environment variable.")
			cfg.Global.Org = cfg.Satellite.Org
		}
	}

	return nil
}

func (app *EarthlyApp) warnIfEarth() {
	if len(os.Args) == 0 {
		return
	}
	binPath := os.Args[0] // can't use os.Executable() here; because it will give us earthly if executed via the earth symlink

	baseName := path.Base(binPath)
	if baseName == "earth" {
		app.BaseCLI.Console().Warnf("Warning: the earth binary has been renamed to earthly; the earth command is currently symlinked, but is deprecated and will one day be removed.")

		absPath, err := filepath.Abs(binPath)
		if err != nil {
			return
		}
		earthlyPath := path.Join(path.Dir(absPath), "earthly")
		earthlyPathExists, _ := fileutil.FileExists(earthlyPath)
		if earthlyPathExists {
			app.BaseCLI.Console().Warnf("Once you are ready to switch over to earthly, you can `rm %s`", absPath)
		}
	}
}

func profhandler() {
	addr := "127.0.0.1:6060"
	fmt.Printf("listening for pprof on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("error listening for pprof: %v", err)
	}
}

func defaultConfigPath(installName string) string {
	earthlyDir := cliutil.GetEarthlyDir(installName)
	oldConfig := filepath.Join(earthlyDir, "config.yaml")
	newConfig := filepath.Join(earthlyDir, "config.yml")
	oldConfigExists, _ := fileutil.FileExists(oldConfig)
	newConfigExists, _ := fileutil.FileExists(newConfig)
	if oldConfigExists && !newConfigExists {
		return oldConfig
	}
	return newConfig
}
