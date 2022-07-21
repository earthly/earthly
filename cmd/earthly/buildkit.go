package main

import (
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/containerutil"
)

func (app *earthlyApp) getBuildkitClient(cliCtx *cli.Context, cloudClient cloud.Client) (*client.Client, error) {
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
	app.console.Warnf("Note: the interactive debugger, interactive RUN commands, and fast output via embedded registry do not yet work on Earthly Satellites.")

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

	if app.ci {
		app.console.Warnf("When using Satellites, the --ci flag is an alias for:")
		app.console.Warnf("  --strict --no-output when running earthly in target mode.")
		app.console.Warnf("  --strict when running earthly in artifact or image mode.")
		app.console.Warnf("") // newline
		if !app.imageMode && !app.artifactMode {
			app.noOutput = true
		}
		app.strict = true
		app.ci = false
	}

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
