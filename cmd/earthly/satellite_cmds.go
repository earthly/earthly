package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/containerutil"
)

func (app *earthlyApp) satelliteCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "launch",
			Usage:       "Launch a new Earthly Satellite",
			Description: "Launch a new Earthly Satellite",
			UsageText: "earthly satellite launch <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] launch <satellite-name>",
			Action: app.actionSatelliteLaunch,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "platform",
					Usage:       "The platform to use when launching a new satellite. Supported values: linux/amd64, linux/arm64.",
					Required:    false,
					Destination: &app.satellitePlatform,
				},
				&cli.StringFlag{
					Name:        "size",
					Usage:       "The size of the satellite. See https://earthly.dev/pricing#compute for details on each size. Supported values: small, medium, large.",
					Required:    false,
					Destination: &app.satelliteSize,
				},
				&cli.StringSliceFlag{
					Name:        "feature-flag",
					EnvVars:     []string{"EARTHLY_SATELLITE_FEATURE_FLAGS"},
					Usage:       "One or more of experimental features to enable on a new satellite",
					Required:    false,
					Hidden:      true,
					Destination: &app.satelliteFeatureFlags,
				},
			},
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
			UsageText: "earthly satellite select <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] select <satellite-name>",
			Action: app.actionSatelliteSelect,
		},
		{
			Name:        "unselect",
			Aliases:     []string{"uns"},
			Usage:       "Remove any currently selected Satellite instance from your Earthly configuration",
			Description: "Remove any currently selected Satellite instance from your Earthly configuration",
			UsageText:   "earthly satellite unselect",
			Action:      app.actionSatelliteUnselect,
		},
		{
			Name:        "wake",
			Usage:       "Manually force a Satellite to wake up from a sleep state",
			Description: "Manually force a Satellite to wake up from a sleep state",
			UsageText: "earthly satellite wake <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] wake <satellite-name>",
			Action: app.actionSatelliteWake,
		},
		{
			Name:        "sleep",
			Usage:       "Manually force a Satellite to sleep from an operational state",
			Description: "Manually force a Satellite to sleep from an operational state",
			UsageText: "earthly satellite sleep <satellite-name>\n" +
				"	earthly satellite [--org <organization-name>] sleep <satellite-name>",
			Action: app.actionSatelliteSleep,
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
	t := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	fmt.Fprintf(t, " \tNAME\tPLATFORM\tSIZE\n") // The leading space is for the selection marker, leave it alone
	for _, s := range satellites {
		var selected = ""
		if s.Name == app.cfg.Satellite.Name && s.Org == orgID {
			selected = "*"
		}
		fmt.Fprintf(t, "%s\t%s\t%s\t%s\n", selected, s.Name, s.Platform, s.Size)
	}
	err := t.Flush()
	if err != nil {
		fmt.Printf("failed to print satellites: %s", err.Error())
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

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	app.console.Printf("Launching Satellite. This could take a moment...\n")
	err = cloudClient.LaunchSatellite(cliCtx.Context, app.satelliteName, orgID, app.satellitePlatform, app.satelliteSize, app.satelliteFeatureFlags.Value())
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

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
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

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
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

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	satelliteToInspect := cliCtx.Args().Get(0)
	selectedSatellite := app.cfg.Satellite.Name

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
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

	app.buildkitdSettings.Timeout = 30 * time.Second
	app.buildkitdSettings.SatelliteToken = token
	app.buildkitdSettings.SatelliteName = satelliteToInspect
	app.buildkitdSettings.SatelliteOrgID = orgID
	if app.satelliteAddress != "" {
		app.buildkitdSettings.BuildkitAddress = app.satelliteAddress
	} else {
		app.buildkitdSettings.BuildkitAddress = containerutil.SatelliteAddress
	}

	selected := "No"
	if selectedSatellite == satellite.Name {
		selected = "Yes"
	}

	app.console.Printf("Instance state: %s", satellite.Status)
	app.console.Printf("Instance platform: %s", satellite.Platform)
	app.console.Printf("Instance size: %s", satellite.Size)
	app.console.Printf("Currently selected: %s", selected)
	app.console.Printf("")

	if satellite.Status == cloud.SatelliteStatusOperational {
		err = buildkitd.PrintSatelliteInfo(cliCtx.Context, app.console, Version, app.buildkitdSettings)
		if err != nil {
			return errors.Wrap(err, "failed checking buildkit info")
		}
	} else {
		app.console.Printf("More info available when Satellite is awake:")
		app.console.Printf("")
		app.console.Printf("    earthly satellite --org %s wake %s", app.orgName, satelliteToInspect)
		app.console.Printf("")
	}
	return nil
}

func (app *earthlyApp) actionSatelliteSelect(cliCtx *cli.Context) error {
	app.commandName = "satelliteSelect"

	if cliCtx.NArg() == 0 {
		if app.cfg.Satellite.Name == "" {
			app.console.Printf("No satellite selected\n\n")
		} else {
			app.console.Printf("Selected satellite: %s\n\n", app.cfg.Satellite.Name)
		}
		_ = cli.ShowCommandHelp(cliCtx, cliCtx.Command.Name)
		return errors.New("satellite name is required")
	}

	if cliCtx.NArg() > 1 {
		_ = cli.ShowCommandHelp(cliCtx, cliCtx.Command.Name)
		return errors.New(fmt.Sprintf("can only provide 1 satellite name, %d provided", cliCtx.NArg()))
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
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

func (app *earthlyApp) actionSatelliteWake(cliCtx *cli.Context) error {
	app.commandName = "satelliteWake"

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	sat, err := cloudClient.GetSatellite(cliCtx.Context, app.satelliteName, orgID)
	if err != nil {
		return err
	}

	if sat.Status == cloud.SatelliteStatusOperational {
		app.console.Printf("%s is already awake.", app.satelliteName)
	}

	out := cloudClient.WakeSatellite(cliCtx.Context, app.satelliteName, orgID)
	err = showSatelliteLoading(app.console, app.satelliteName, out)
	if err != nil {
		return errors.Wrap(err, "failed waiting for satellite wake")
	}

	return nil
}

func (app *earthlyApp) actionSatelliteSleep(cliCtx *cli.Context) error {
	app.commandName = "satelliteSleep"

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	app.satelliteName = cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgID, err := app.getSatelliteOrgID(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	out := cloudClient.SleepSatellite(cliCtx.Context, app.satelliteName, orgID)
	err = showSatelliteStopping(app.console, app.satelliteName, out)
	if err != nil {
		return errors.Wrap(err, "failed waiting for satellite wake")
	}

	return nil
}

func showSatelliteLoading(console conslogging.ConsoleLogger, satName string, out chan cloud.SatelliteStatusUpdate) error {
	loadingMsgs := getSatelliteLoadingMessages()
	var (
		loggedSleep      bool
		loggedStop       bool
		loggedStart      bool
		shouldLogLoading bool
	)
	for o := range out {
		if o.Err != nil {
			return errors.Wrap(o.Err, "failed processing satellite status")
		}
		shouldLogLoading = true
		switch o.State {
		case cloud.SatelliteStatusSleep:
			if !loggedSleep {
				console.Printf("%s is waking up. Please wait...", satName)
				loggedSleep = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusStopping:
			if !loggedStop {
				console.Printf("%s is currently falling asleep. Waiting to send wake up signal...", satName)
				loggedStop = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusStarting:
			if !loggedStart && !loggedSleep {
				console.Printf("%s is starting. Please wait...", satName)
				loggedStart = true
				shouldLogLoading = false
			}
		case cloud.SatelliteStatusOperational:
			if loggedSleep || loggedStop || loggedStart {
				// Satellite was in a different state previously but is now online
				console.Printf("...System online.")
			}
			shouldLogLoading = false
		default:
			// In case there's a new state later which we didn't expect here,
			// we'll still try to inform the user as best we can.
			// Note the state might just be "Unknown" if it maps to an gRPC enum we don't know about.
			console.Printf("%s state is: %s", satName, o)
			shouldLogLoading = false
		}
		if shouldLogLoading {
			var msg string
			msg, loadingMsgs = nextSatelliteLoadingMessage(loadingMsgs)
			console.Printf("...%s...", msg)
		}
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

func showSatelliteStopping(console conslogging.ConsoleLogger, satName string, out chan cloud.SatelliteStatusUpdate) error {
	loggedStopping := false
	for o := range out {
		if o.Err != nil {
			return errors.Wrap(o.Err, "failed processing satellite status")
		}
		switch o.State {
		case cloud.SatelliteStatusSleep:
			if !loggedStopping {
				console.Printf("%s is already asleep", satName)
			} else {
				console.Printf("...Done.")
			}
		case cloud.SatelliteStatusOperational:
			console.Printf("%s is going to sleep. Please wait...", satName)
		case cloud.SatelliteStatusStopping:
			loggedStopping = true
			console.Printf("...still shutting down...")
		}
	}
	return nil
}
