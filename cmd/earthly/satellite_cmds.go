package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/util/containerutil"
)

const selectUsageText = "earthly satellite select <satellite-name>\n" +
	"	earthly satellite [--org <organization-name>] select <satellite-name>"

func (app *earthlyApp) satelliteCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "launch",
			Usage:       "Launch a new Earthly Satellite",
			Description: "Launch a new Earthly Satellite",
			UsageText: "earthly satellite launch <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] launch <satellite-name>",
			Action: app.actionSatelliteLaunch,
		},
		{
			Name:        "rm",
			Usage:       "Destroy an Earthly Satellite",
			Description: "Destroy an Earthly Satellite",
			UsageText: "earthly satellite rm <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] rm <satellite-name>",
			Action: app.actionSatelliteRemove,
		},
		{
			Name:        "ls",
			Description: "List your Earthly Satellites",
			Usage:       "List your Earthly Satellites",
			UsageText: "earthly satellite ls\n" +
				"	earthly satellite [--org <organization-name>] ls",
			Action: app.actionSatelliteList,
		},
		{
			Name:        "inspect",
			Description: "Show additional details about a Satellite instance",
			Usage:       "Show additional details about a Satellite instance",
			UsageText: "earthly satellite inspect <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] inspect <satellite-name>",
			Action: app.actionSatelliteInspect,
		},
		{
			Name:        "select",
			Aliases:     []string{"s"},
			Usage:       "Choose which satellite to use to build your app",
			Description: "Choose which satellite to use to build your app",
			UsageText:   selectUsageText,
			Action:      app.actionSatelliteSelect,
		},
		{
			Name:        "unselect",
			Aliases:     []string{"uns"},
			Usage:       "Remove any currently selected Satellite instance from your Earthly configuration",
			Description: "Remove any currently selected Satellite instance from your Earthly configuration",
			UsageText:   "earthly satellite unselect",
			Action:      app.actionSatelliteUnselect,
		},
	}
}

func (app *earthlyApp) useSatellite(cliCtx *cli.Context, satelliteName, orgName string) error {
	inConfig, err := config.ReadConfigFile(app.configPath)
	if err != nil {
		if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "read config")
		}
	}

	newConfig, err := config.Upsert(inConfig, "satellite.name", satelliteName)
	if err != nil {
		return errors.Wrap(err, "could not update satellite name")
	}
	// Update in-place so we can use it later, assuming the config change was successful.
	app.cfg.Satellite.Name = satelliteName

	newConfig, err = config.Upsert(newConfig, "satellite.org", orgName)
	if err != nil {
		return errors.Wrap(err, "could not update satellite name")
	}
	app.cfg.Satellite.Org = orgName
	err = config.WriteConfigFile(app.configPath, newConfig)
	if err != nil {
		return errors.Wrap(err, "could not save config")
	}
	app.console.Printf("Updated selected satellite in %s", app.configPath)

	return nil
}

func (app *earthlyApp) printSatellites(satellites []cloud.SatelliteInstance, orgID string) {
	for _, satellite := range satellites {
		if satellite.Name == app.cfg.Satellite.Name && satellite.Org == orgID {
			fmt.Printf("* %s\n", satellite.Name)
		} else {
			fmt.Printf("  %s\n", satellite.Name)
		}
	}
}

func (app *earthlyApp) getSatelliteOrgID(ctx context.Context, cloudClient cloud.Client) (string, error) {
	// We are cheating here and forcing a re-auth before running any satellite commands.
	// This is because there is an issue on the backend where the token might be outdated
	// if a user was invited to an org recently after already logging-in.
	// TODO Eventually we should be able to remove this cheat.
	err := cloudClient.Authenticate(ctx)
	if err != nil {
		return "", errors.New("unable to authenticate")
	}
	var orgID string
	if app.orgName == "" {
		orgs, err := cloudClient.ListOrgs(ctx)
		if err != nil {
			return "", errors.Wrap(err, "failed finding org")
		}
		if len(orgs) == 0 {
			return "", errors.New("not a member of any organizations - satellites only work within an org")
		}
		if len(orgs) > 1 {
			return "", errors.New("more than one organizations available - please specify the name of the organization using `--org`")
		}
		app.orgName = orgs[0].Name
		orgID = orgs[0].ID
	} else {
		var err error
		orgID, err = cloudClient.GetOrgID(ctx, app.orgName)
		if err != nil {
			return "", errors.Wrap(err, "invalid org provided")
		}
	}
	return orgID, nil
}

func (app *earthlyApp) actionSatelliteLaunch(cliCtx *cli.Context) error {
	app.commandName = "satelliteLaunch"

	if cliCtx.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	app.console.Printf("Launching Satellite. This could take a moment...\n")
	err = cloudClient.LaunchSatellite(cliCtx.Context, app.satelliteName, orgID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			app.console.Printf("Operation interrupted. Satellite should finish launching in background (if server received request).\n")
			return nil
		}
		return errors.Wrapf(err, "failed to create satellite %s", app.satelliteName)
	}
	app.console.Printf("...Done\n")

	err = app.useSatellite(cliCtx, app.satelliteName, app.orgName)
	if err != nil {
		return errors.Wrap(err, "could not configure satellite for use")
	}
	app.console.Printf("The satellite %s has been automatically selected for use. To go back to using local builds you can use\n\n\tearthly satellite unselect\n\n", app.satelliteName)

	return nil
}

func (app *earthlyApp) actionSatelliteList(cliCtx *cli.Context) error {
	app.commandName = "satelliteList"

	if cliCtx.NArg() != 0 {
		return errors.New("command does not accept any arguments")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satellites, err := cloudClient.ListSatellites(cliCtx.Context, orgID)
	if err != nil {
		return err
	}

	app.printSatellites(satellites, orgID)
	return nil
}

func (app *earthlyApp) actionSatelliteRemove(cliCtx *cli.Context) error {
	app.commandName = "satelliteRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	app.console.Printf("Destroying Satellite. This could take a moment...\n")
	err = cloudClient.DeleteSatellite(cliCtx.Context, app.satelliteName, orgID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			app.console.Printf("Operation interrupted. Satellite should finish destroying in background (if server received request).\n")
			return nil
		}
		return errors.Wrapf(err, "failed to delete satellite %s", app.satelliteName)
	}
	app.console.Printf("...Done\n")

	if app.satelliteName == app.cfg.Satellite.Name {
		err = app.useSatellite(cliCtx, "", "")
		if err != nil {
			return errors.Wrapf(err, "failed unselecting satellite")
		}
		app.console.Printf("Satellite has also been unselected\n")
	}
	return nil
}

func (app *earthlyApp) actionSatelliteInspect(cliCtx *cli.Context) error {
	app.commandName = "satelliteInspect"

	if cliCtx.NArg() != 1 {
		return errors.New("satellite name is required")
	}

	satelliteToInspect := cliCtx.Args().Get(0)
	selectedSatellite := app.cfg.Satellite.Name

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satellite, err := cloudClient.GetSatellite(cliCtx.Context, satelliteToInspect, orgID)
	if err != nil {
		return err
	}

	token, err := cloudClient.GetAuthToken(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}

	app.buildkitdSettings.SatelliteToken = token
	app.buildkitdSettings.SatelliteName = satelliteToInspect
	app.buildkitdSettings.SatelliteOrgID = orgID
	if app.satelliteAddress != "" {
		app.buildkitdSettings.BuildkitAddress = app.satelliteAddress
	} else {
		app.buildkitdSettings.BuildkitAddress = containerutil.SatelliteAddress
	}

	err = buildkitd.PrintSatelliteInfo(cliCtx.Context, app.console, Version, app.buildkitdSettings)
	if err != nil {
		return errors.Wrap(err, "failed checking buildkit info")
	}

	selected := "No"
	if selectedSatellite == satellite.Name {
		selected = "Yes"
	}
	app.console.Printf("Instance state: %s", satellite.Status)
	app.console.Printf("Currently selected: %s", selected)
	return nil
}

func (app *earthlyApp) actionSatelliteSelect(cliCtx *cli.Context) error {
	app.commandName = "satelliteSelect"

	if cliCtx.NArg() != 1 {
		if app.cfg.Satellite.Name == "" {
			app.console.Printf("No satellite selected\n\n")
			app.console.Printf("To select a satellite:\n\t%s\n", selectUsageText)
		} else {
			app.console.Printf("Selected satellite: %s\n\n", app.cfg.Satellite.Name)
			app.console.Printf("To select a different satellite:\n\t%s\n", selectUsageText)
		}
		return errors.New("satellite name is required")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satellites, err := cloudClient.ListSatellites(cliCtx.Context, orgID)
	if err != nil {
		return err
	}

	found := false
	for _, s := range satellites {
		if app.satelliteName == s.Name {
			err = app.useSatellite(cliCtx, s.Name, app.orgName)
			if err != nil {
				return errors.Wrapf(err, "could not select satellite %s", app.satelliteName)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%s is not a valid satellite", app.satelliteName)
	}

	app.printSatellites(satellites, orgID)
	return nil
}

func (app *earthlyApp) actionSatelliteUnselect(cliCtx *cli.Context) error {
	app.commandName = "satelliteUnselect"

	if cliCtx.NArg() != 0 {
		return errors.New("command does not accept any arguments")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	err := app.useSatellite(cliCtx, "", "")
	if err != nil {
		return errors.Wrap(err, "could not unselect satellite")
	}

	return nil
}
