package main

import (
	"context"
	"math/rand"
	"net/url"
	"path/filepath"
	"time"

	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
)

func (app *earthlyApp) initFrontend(cliCtx *cli.Context) error {
	console := app.console.WithPrefix("frontend")
	feConfig := &containerutil.FrontendConfig{
		BuildkitHostCLIValue:       app.buildkitHost,
		BuildkitHostFileValue:      app.cfg.Global.BuildkitHost,
		DebuggerHostCLIValue:       app.debuggerHost,
		DebuggerHostFileValue:      app.cfg.Global.DebuggerHost,
		DebuggerPortFileValue:      app.cfg.Global.DebuggerPort,
		LocalRegistryHostFileValue: app.cfg.Global.LocalRegistryHost,
		Console:                    console,
	}
	fe, err := containerutil.FrontendForSetting(cliCtx.Context, app.cfg.Global.ContainerFrontend, feConfig)
	if err != nil {
		origErr := err
		fe, err = containerutil.NewStubFrontend(cliCtx.Context, feConfig)
		if err != nil {
			return errors.Wrap(err, "failed stub frontend initialization")
		}

		if !app.verbose {
			console.Printf("No frontend initialized. Use --verbose to see details\n")
		}
		console.VerbosePrintf("%s frontend initialization failed due to %s", app.cfg.Global.ContainerFrontend, origErr.Error())
	} else {
		console.VerbosePrintf("%s frontend initialized.\n", fe.Config().Setting)
	}
	app.containerFrontend = fe

	// command line option overrides the config which overrides the default value
	if !cliCtx.IsSet("buildkit-image") && app.cfg.Global.BuildkitImage != "" {
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
		app.handleTLSCertificateSettings(cliCtx)
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
	app.buildkitdSettings.EnableProfiler = app.enableProfiler
	app.buildkitdSettings.NoUpdate = app.noBuildkitUpdate

	// ensure the MTU is something allowable in IPv4, cap enforced by type. Zero is autodetect.
	if app.cfg.Global.CniMtu != 0 && app.cfg.Global.CniMtu < 68 {
		return errors.New("invalid overridden MTU size")
	}
	app.buildkitdSettings.CniMtu = app.cfg.Global.CniMtu

	if app.cfg.Global.IPTables != "" && app.cfg.Global.IPTables != "iptables-legacy" && app.cfg.Global.IPTables != "iptables-nft" {
		return errors.New(`invalid overridden iptables name. Valid values are "iptables-legacy" or "iptables-nft"`)
	}
	app.buildkitdSettings.IPTables = app.cfg.Global.IPTables
	earthlyDir, err := cliutil.GetOrCreateEarthlyDir()
	if err != nil {
		return errors.Wrap(err, "failed to get earthly dir")
	}
	app.buildkitdSettings.StartUpLockPath = filepath.Join(earthlyDir, "buildkitd-startup.lock")
	return nil
}

func (app *earthlyApp) getBuildkitClient(cliCtx *cli.Context, cloudClient cloud.Client) (*client.Client, error) {
	err := app.initFrontend(cliCtx)
	if err != nil {
		return nil, err
	}
	err = app.configureSatellite(cliCtx, cloudClient)
	if err != nil {
		return nil, errors.Wrapf(err, "could not construct new buildkit client")
	}

	return buildkitd.NewClient(cliCtx.Context, app.console, app.buildkitdImage, app.containerName, app.containerFrontend, Version, app.buildkitdSettings)
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

func (app *earthlyApp) configureSatellite(cliCtx *cli.Context, cloudClient cloud.Client) error {
	if cliCtx.IsSet("buildkit-host") && cliCtx.IsSet("satellite") {
		return errors.New("cannot specify both buildkit-host and satellite")
	}
	if cliCtx.IsSet("satellite") && app.noSatellite {
		return errors.New("cannot specify both no-satellite and satellite")
	}
	if !app.isUsingSatellite(cliCtx) || cloudClient == nil {
		// If the app is not using a cloud client, or the command doesn't interact with the cloud (prune, bootstrap)
		// then pretend its all good and use your regular configuration.
		return nil
	}

	// Set up extra settings needed for buildkit RPC metadata
	if app.satelliteName == "" {
		app.satelliteName = app.cfg.Satellite.Name
	}
	if app.orgName == "" {
		app.orgName = app.cfg.Satellite.Org
	}
	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}
	app.buildkitdSettings.SatelliteName = app.satelliteName
	app.buildkitdSettings.SatelliteOrgID = orgID
	if app.satelliteAddress != "" {
		app.buildkitdSettings.BuildkitAddress = app.satelliteAddress
	} else {
		app.buildkitdSettings.BuildkitAddress = containerutil.SatelliteAddress
	}
	app.analyticsMetadata.isSatellite = true
	app.analyticsMetadata.satelliteVersion = "" // TODO

	app.console.Printf("") // newline
	app.console.Printf("The following feature flag is recommended for use with Satellites and will be auto-enabled:")
	app.console.Printf("  --new-platform")
	app.console.Printf("") // newline

	if app.featureFlagOverrides != "" {
		app.featureFlagOverrides += ","
	}
	app.featureFlagOverrides += "new-platform"

	token, err := cloudClient.GetAuthToken(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	app.buildkitdSettings.SatelliteToken = token

	// Reserve the satellite for the upcoming build.
	// This operation can take a moment if the satellite is asleep.
	err = app.reserveSatellite(cliCtx.Context, cloudClient, app.satelliteName, orgID)
	if err != nil {
		return err
	}

	// TODO (dchw) what other settings might we want to override here?
	return nil
}

func (app *earthlyApp) isUsingSatellite(cliCtx *cli.Context) bool {
	if app.noSatellite {
		return false
	}
	if cliCtx.IsSet("buildkit-host") {
		// buildkit-host takes precedence
		return false
	}
	return app.cfg.Satellite.Name != "" || app.satelliteName != ""
}

func (app *earthlyApp) reserveSatellite(ctx context.Context, cloudClient cloud.Client, name, orgID string) error {
	console := app.console.WithPrefix("satellite")
	out := make(chan string)
	var reserveErr error
	go func() { reserveErr = cloudClient.ReserveSatellite(ctx, name, orgID, out) }()
	loadingMsgs := getSatelliteLoadingMessages()
	var (
		loggedSleep      bool
		loggedStop       bool
		loggedStart      bool
		shouldLogLoading bool
	)
	for status := range out {
		shouldLogLoading = true
		switch status {
		case cloud.SatelliteStatusSleep:
			if !loggedSleep {
				console.Printf("%s is waking up. Please wait...", name)
				loggedSleep = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusStopping:
			if !loggedStop {
				console.Printf("%s is falling asleep. Please wait...", name)
				loggedStop = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusStarting:
			if !loggedStart && !loggedSleep {
				console.Printf("%s is starting. Please wait...", name)
				loggedStart = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusOperational:
			// Should be last update received at this point.
			console.Printf("...System online.")
			shouldLogLoading = false
		default:
			// In case there's a new state later which we didn't expect here,
			// we'll still try to inform the user as best we can.
			// Note the state might just be "Unknown" if it maps to an gRPC enum we don't know about.
			console.Printf("%s state is: %s", name, status)
			shouldLogLoading = false
		}
		if shouldLogLoading {
			var msg string
			msg, loadingMsgs = nextSatelliteLoadingMessage(loadingMsgs)
			console.Printf("...%s...", msg)
		}
	}
	if reserveErr != nil {
		return errors.Wrap(reserveErr, "failed reserving satellite for build")
	}
	return nil
}

func nextSatelliteLoadingMessage(msgs []string) (nextMsg string, remainingMsgs []string) {
	if len(msgs) == 0 {
		msgs = getSatelliteLoadingMessages()
	}
	return msgs[0], msgs[1:]
}

func getSatelliteLoadingMessages() []string {
	baseMessages := []string{
		"tracking orbit",
		"adjusting course",
		"deploying solar array",
		"aligning solar panels",
		"calibrating guidance system",
		"establishing transponder uplink",
		"testing signal quality",
		"fueling thrusters",
		"amplifying transmission signal",
		"checking thermal controls",
		"stabilizing trajectory",
		"contacting mission control",
		"testing antennas",
		"reporting fuel levels",
		"scanning surroundings",
		"avoiding debris",
		"taking solar reading",
		"reporting thermal conditions",
		"testing system integrity",
		"checking battery levels",
		"calibrating transponders",
		"modifying downlink frequency",
	}
	msgs := baseMessages
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(msgs), func(i, j int) { msgs[i], msgs[j] = msgs[j], msgs[i] })
	return msgs
}
