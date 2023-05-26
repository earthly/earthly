package main

import (
	"context"
	"net/url"
	"path/filepath"
	"time"

	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/containerutil"
)

func (app *earthlyApp) initFrontend(cliCtx *cli.Context) error {
	// command line option overrides the config which overrides the default value
	if !cliCtx.IsSet("buildkit-image") && app.cfg.Global.BuildkitImage != "" {
		app.buildkitdImage = app.cfg.Global.BuildkitImage
	}

	bkURL, err := url.Parse(app.buildkitHost) // Not validated because we already did that when we calculated it.
	if err != nil {
		return errors.Wrap(err, "failed to parse generated buildkit URL")
	}

	if bkURL.Scheme == "tcp" && app.cfg.Global.TLSEnabled {
		app.buildkitdSettings.ClientTLSCert = app.cfg.Global.ClientTLSCert
		app.buildkitdSettings.ClientTLSKey = app.cfg.Global.ClientTLSKey
		app.buildkitdSettings.TLSCA = app.cfg.Global.TLSCACert
		app.buildkitdSettings.ServerTLSCert = app.cfg.Global.ServerTLSCert
		app.buildkitdSettings.ServerTLSKey = app.cfg.Global.ServerTLSKey
	}

	app.buildkitdSettings.AdditionalArgs = app.cfg.Global.BuildkitAdditionalArgs
	app.buildkitdSettings.AdditionalConfig = app.cfg.Global.BuildkitAdditionalConfig
	app.buildkitdSettings.Timeout = time.Duration(app.cfg.Global.BuildkitRestartTimeoutS) * time.Second
	app.buildkitdSettings.Debug = app.debug
	app.buildkitdSettings.BuildkitAddress = app.buildkitHost
	app.buildkitdSettings.LocalRegistryAddress = app.localRegistryHost
	app.buildkitdSettings.UseTCP = bkURL.Scheme == "tcp"
	app.buildkitdSettings.UseTLS = app.cfg.Global.TLSEnabled
	app.buildkitdSettings.MaxParallelism = app.cfg.Global.BuildkitMaxParallelism
	app.buildkitdSettings.CacheSizeMb = app.cfg.Global.BuildkitCacheSizeMb
	app.buildkitdSettings.CacheSizePct = app.cfg.Global.BuildkitCacheSizePct
	app.buildkitdSettings.CacheKeepDuration = app.cfg.Global.BuildkitCacheKeepDurationS
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
	earthlyDir, err := cliutil.GetOrCreateEarthlyDir(app.installationName)
	if err != nil {
		return errors.Wrap(err, "failed to get earthly dir")
	}
	app.buildkitdSettings.StartUpLockPath = filepath.Join(earthlyDir, "buildkitd-startup.lock")
	return nil
}

func (app *earthlyApp) getBuildkitClient(cliCtx *cli.Context, cloudClient *cloud.Client) (*client.Client, error) {
	err := app.initFrontend(cliCtx)
	if err != nil {
		return nil, err
	}
	err = app.configureSatellite(cliCtx, cloudClient, "", "") // no gitAuthor/gitConfigEmail is passed for non-build commands (e.g. debug_cmds.go or root_cmds.go code)
	if err != nil {
		return nil, errors.Wrapf(err, "could not construct new buildkit client")
	}

	return buildkitd.NewClient(cliCtx.Context, app.console, app.buildkitdImage, app.containerName, app.installationName, app.containerFrontend, Version, app.buildkitdSettings)
}

func (app *earthlyApp) configureSatellite(cliCtx *cli.Context, cloudClient *cloud.Client, gitAuthor, gitConfigEmail string) error {
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

	app.buildkitdSettings.UseTCP = true
	if app.cfg.Global.TLSEnabled {
		// satellite connection with tls enabled does not use configuration certificates
		app.buildkitdSettings.ClientTLSCert = ""
		app.buildkitdSettings.ClientTLSKey = ""
		app.buildkitdSettings.TLSCA = ""
		app.buildkitdSettings.ServerTLSCert = ""
		app.buildkitdSettings.ServerTLSKey = ""
	}

	_, orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return errors.Wrap(err, "failed getting org")
	}
	satelliteName, err := app.getSatelliteName(cliCtx.Context, orgID, app.satelliteName, cloudClient)
	if err != nil {
		return errors.Wrap(err, "failed getting satellite name")
	}
	app.buildkitdSettings.SatelliteName = satelliteName
	app.buildkitdSettings.SatelliteDisplayName = app.satelliteName
	app.buildkitdSettings.SatelliteOrgID = orgID
	if app.satelliteAddress != "" {
		app.buildkitdSettings.BuildkitAddress = app.satelliteAddress
	} else {
		app.buildkitdSettings.BuildkitAddress = containerutil.SatelliteAddress
	}
	app.analyticsMetadata.isSatellite = true
	app.analyticsMetadata.satelliteCurrentVersion = "" // TODO

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
	err = app.reserveSatellite(cliCtx.Context, cloudClient, satelliteName, app.satelliteName, orgID, gitAuthor, gitConfigEmail)
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

func (app *earthlyApp) reserveSatellite(ctx context.Context, cloudClient *cloud.Client, name, displayName, orgID, gitAuthor, gitConfigEmail string) error {
	console := app.console.WithPrefix("satellite")
	_, isCI := analytics.DetectCI(app.earthlyCIRunner)
	out := cloudClient.ReserveSatellite(ctx, name, orgID, gitAuthor, gitConfigEmail, isCI)
	err := showSatelliteLoading(console, displayName, out)
	if err != nil {
		return errors.Wrap(err, "failed reserving satellite for build")
	}
	return nil
}
