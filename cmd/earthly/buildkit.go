package main

import (
	"net/url"
	"time"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) initFrontend(cliCtx *cli.Context) error {
	feConfig := &containerutil.FrontendConfig{
		BuildkitHostCLIValue:       app.buildkitHost,
		BuildkitHostFileValue:      app.cfg.Global.BuildkitHost,
		DebuggerHostCLIValue:       app.debuggerHost,
		DebuggerHostFileValue:      app.cfg.Global.DebuggerHost,
		DebuggerPortFileValue:      app.cfg.Global.DebuggerPort,
		LocalRegistryHostFileValue: app.cfg.Global.LocalRegistryHost,
		Console:                    app.console,
	}
	fe, err := containerutil.FrontendForSetting(cliCtx.Context, app.cfg.Global.ContainerFrontend, feConfig)
	if err != nil {
		origErr := err
		fe, err = containerutil.NewStubFrontend(cliCtx.Context, feConfig)
		if err != nil {
			return errors.Wrap(err, "failed frontend initialization")
		}
		app.console.Warnf("%s frontend initialization failed due to %s; but will try anyway", app.cfg.Global.ContainerFrontend, origErr.Error())
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
	return nil
}

func (app *earthlyApp) getBuildkitClient(cliCtx *cli.Context, cloudClient cloud.Client) (*client.Client, error) {
	if !app.isUsingSatellite(cliCtx) {
		err := app.initFrontend(cliCtx)
		if err != nil {
			return nil, err
		}
	}
	err := app.configureSatellite(cliCtx, cloudClient)
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

	// When using a satellite, interactive and local do not work; as they are not SSL nor routable yet.
	app.console.Warnf("Note: the interactive debugger, interactive RUN commands, and inline caching do not yet work on Earthly Satellites.")

	// Set up extra settings needed for buildkit RPC metadata
	if app.satelliteName == "" {
		app.satelliteName = app.cfg.Satellite.Name
	}
	if app.satelliteOrg == "" {
		app.satelliteOrg = app.cfg.Satellite.Org
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

	app.console.Warnf("") // newline
	app.console.Warnf("The following feature flags are recommended for use with Satellites and will be auto-enabled:")
	app.console.Warnf("  --new-platform, --use-registry-for-with-docker")
	app.console.Warnf("") // newline

	if app.featureFlagOverrides != "" {
		app.featureFlagOverrides += ","
	}
	app.featureFlagOverrides += "new-platform,use-registry-for-with-docker"

	token, err := cloudClient.GetAuthToken(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	app.buildkitdSettings.SatelliteToken = token

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
