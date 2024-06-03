package subcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/base"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
)

type Satellite struct {
	cli CLI

	platform               string
	size                   string
	featureFlags           cli.StringSlice
	maintenanceWindow      string
	maintenaceWeekendsOnly bool
	version                string
	cloudName              string
	force                  bool
	printJSON              bool
	dropCache              bool
}

func NewSatellite(cli CLI) *Satellite {
	return &Satellite{
		cli: cli,
	}
}

func (a *Satellite) Cmds() []*cli.Command {
	return []*cli.Command{
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
					Destination: &a.cli.Flags().OrgName,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "launch",
					Usage:       "Launch a new Earthly Satellite",
					Description: "Launch a new Earthly Satellite.",
					UsageText: "earthly satellite launch <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] launch <satellite-name>",
					Action: a.actionLaunch,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "platform",
							Usage: `The platform to use when launching a new satellite
					Supported values: linux/amd64, linux/arm64`,
							Required:    false,
							Value:       cloud.SatellitePlatformAMD64,
							Destination: &a.platform,
						},
						&cli.StringFlag{
							Name: "size",
							Usage: `The size of the satellite. See https://earthly.dev/pricing for details on each size
					Supported values: xsmall, small, medium, large, xlarge, 2xlarge, 3xlarge, 4xlarge`,
							Required:    false,
							Destination: &a.size,
						},
						&cli.StringSliceFlag{
							Name:        "feature-flag",
							EnvVars:     []string{"EARTHLY_SATELLITE_FEATURE_FLAGS"},
							Usage:       "One or more of experimental features to enable on a new satellite",
							Required:    false,
							Hidden:      true,
							Destination: &a.featureFlags,
						},
						&cli.StringFlag{
							Name:    "maintenance-window",
							Aliases: []string{"mw"},
							Usage: `Sets a maintenance window for satellite auto-updates
					If there is a new satellite version available, the satellite will update within 2 hrs of the time specified.
					Format must be in HH:MM (24 hr) and will be automatically converted from your current local time to UTC.
					Default value is 02:00 in your local time.`,
							Required:    false,
							Destination: &a.maintenanceWindow,
						},
						&cli.BoolFlag{
							Name:        "maintenance-weekends-only",
							Aliases:     []string{"wo"},
							Usage:       "When set, satellite auto-updates will only occur on Saturday or Sunday during the specified maintenance window",
							Required:    false,
							Destination: &a.maintenaceWeekendsOnly,
						},
						&cli.StringFlag{
							Name:        "version",
							Usage:       "Launch and pin a satellite at a specific version (disables auto-updates)",
							Required:    false,
							Destination: &a.version,
						},
						&cli.StringFlag{
							Name:        "cloud",
							Usage:       "Launch the satellite within a configured cloud installation (only applies to accounts subscribed to the BYOC plan)",
							Required:    false,
							Destination: &a.cloudName,
						},
					},
				},
				{
					Name:        "rm",
					Usage:       "Destroy an Earthly Satellite",
					Description: "Destroy an Earthly Satellite.",
					UsageText: "earthly satellite rm <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] rm <satellite-name>",
					Action: a.actionRemove,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "force",
							Aliases:     []string{"f"},
							Usage:       "Forces the removal of the satellite, even if it's running",
							Required:    false,
							Destination: &a.force,
						},
					},
				},
				{
					Name:        "ls",
					Usage:       "List your Earthly Satellites",
					Description: "List your Earthly Satellites.",
					UsageText: "earthly satellite ls\n" +
						"	earthly satellite [--org <organization-name>] ls",
					Action: a.actionSatelliteList,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "json",
							Usage:       "Prints the output in JSON format",
							Required:    false,
							Destination: &a.printJSON,
						},
					},
				},
				{
					Name:        "inspect",
					Usage:       "Show additional details about an Earthly Satellite instance",
					Description: "Show additional details about an Earthly Satellite instance.",
					UsageText: "earthly satellite inspect <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] inspect <satellite-name>",
					Action: a.actionInspect,
				},
				{
					Name:        "select",
					Aliases:     []string{"s"},
					Usage:       "Choose which Earthly Satellite to use to build your app",
					Description: "Choose which Earthly Satellite to use to build your a.",
					UsageText: "earthly satellite select <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] select <satellite-name>",
					Action: a.actionSelect,
				},
				{
					Name:        "unselect",
					Aliases:     []string{"uns"},
					Usage:       "Remove any currently selected Earthly Satellite instance from your Earthly configuration",
					Description: "Remove any currently selected Earthly Satellite instance from your Earthly configuration.",
					UsageText:   "earthly satellite unselect",
					Action:      a.actionUnselect,
				},
				{
					Name:        "wake",
					Usage:       "Manually force an Earthly Satellite to wake up from a sleep state",
					Description: "Manually force an Earthly Satellite to wake up from a sleep state.",
					UsageText: "earthly satellite wake <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] wake <satellite-name>",
					Action: a.actionWake,
				},
				{
					Name:  "sleep",
					Usage: "Manually force an Earthly Satellite to sleep from an operational state",
					Description: "Manually force an Earthly Satellite to sleep from an operational state.\n" +
						"Note that this may interrupt ongoing builds.",
					UsageText: "earthly satellite sleep <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] sleep <satellite-name>",
					Action: a.actionSleep,
				},
				{
					Name:        "update",
					Usage:       "Manually update an Earthly Satellite to the latest version (may cause downtime)",
					Description: "Manually update an Earthly Satellite to the latest version (may cause downtime).",
					UsageText: "earthly satellite update <satellite-name>\n" +
						"	earthly satellite [--org <organization-name>] update <satellite-name>",
					Action: a.actionUpdate,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "platform",
							Usage: `The platform to use when launching a new satellite
					Supported values: linux/amd64, linux/arm64`,
							Required:    false,
							Value:       cloud.SatellitePlatformAMD64,
							Destination: &a.platform,
						},
						&cli.StringFlag{
							Name: "size",
							Usage: `Change the size of the satellite. See https://earthly.dev/pricing for details on each size.
					Supported values: xsmall, small, medium, large, xlarge, 2xlarge, 3xlarge, 4xlarge`,
							Required:    false,
							Destination: &a.size,
						},
						&cli.StringFlag{
							Name:        "maintenance-window",
							Aliases:     []string{"mw"},
							Usage:       "Set a new custom maintenance window for future satellite auto-updates",
							Required:    false,
							Destination: &a.maintenanceWindow,
						},
						&cli.BoolFlag{
							Name:        "maintenance-weekends-only",
							Aliases:     []string{"wo"},
							Usage:       "When set, satellite auto-updates will only occur on Saturday or Sunday during the specified maintenance window",
							Required:    false,
							Destination: &a.maintenaceWeekendsOnly,
						},
						&cli.BoolFlag{
							Name:        "drop-cache",
							Usage:       "Drop existing cache as part of the update operation",
							Required:    false,
							Destination: &a.dropCache,
						},
						&cli.StringSliceFlag{
							Name:        "feature-flag",
							EnvVars:     []string{"EARTHLY_SATELLITE_FEATURE_FLAGS"},
							Usage:       "One or more of experimental features to enable on the updated satellite",
							Required:    false,
							Hidden:      true,
							Destination: &a.featureFlags,
						},
						&cli.StringFlag{
							Name:        "version",
							Usage:       "Launch a specific satellite version (disables auto-updates)",
							Required:    false,
							Hidden:      true,
							Destination: &a.version,
						},
						&cli.BoolFlag{
							Name:        "force",
							Aliases:     []string{"f"},
							Usage:       "Forces the satellite to sleep (if necessary) before starting the updating",
							Required:    false,
							Destination: &a.force,
						},
					},
				},
			},
		},
	}
}

func (a *Satellite) useSatellite(cliCtx *cli.Context, satelliteName, orgName string) error {
	inConfig, err := config.ReadConfigFile(a.cli.Flags().ConfigPath)
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
	a.cli.Cfg().Satellite.Name = satelliteName

	err = config.WriteConfigFile(a.cli.Flags().ConfigPath, newConfig)
	if err != nil {
		return errors.Wrap(err, "could not save config")
	}
	a.cli.Console().Printf("Updated selected satellite in %s", a.cli.Flags().ConfigPath)

	return nil
}

func printRow(t *tabwriter.Writer, c []color.Attribute, items []string) {
	sprint := color.New(c...).SprintFunc()

	for idx, item := range items {
		items[idx] = sprint(item)
	}
	line := strings.Join(items, "\t")
	line += "\n"

	fmt.Fprint(t, line)
}

func (a *Satellite) printSatellitesTable(satellites []cloud.SatelliteInstance, isOrgSelected bool) {
	slices.SortStableFunc(satellites, func(a, b cloud.SatelliteInstance) int {
		return strings.Compare(a.Name, b.Name)
	})

	includeCloudColumn := false
	cloudNames := make(map[string]bool)
	for _, s := range satellites {
		if n := s.CloudName; n != "" {
			cloudNames[n] = true
		}
	}
	if len(cloudNames) > 1 {
		includeCloudColumn = true
	}

	t := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	headerRow := []string{" ", "NAME", "PLATFORM", "SIZE", "VERSION", "STATE"} // The leading space is for the selection marker, leave it alone
	if includeCloudColumn {
		headerRow = append(headerRow, "CLOUD")
	}
	printRow(t, []color.Attribute{color.Reset}, headerRow)

	for _, s := range satellites {
		var selected = ""
		if s.Name == a.cli.Cfg().Satellite.Name && isOrgSelected {
			selected = "*"
		}

		row := []string{selected, s.Name, s.Platform, s.Size, s.Version, strings.ToLower(s.State)}
		c := []color.Attribute{color.Reset}
		if includeCloudColumn {
			name := s.CloudName
			if name == "" {
				name = "-"
			}
			row = append(row, name)
		}

		printRow(t, c, row)
	}
	err := t.Flush()
	if err != nil {
		fmt.Printf("failed to print satellites: %s", err.Error())
	}
}

func durationWithDaysPart(d time.Duration) string {
	sd := d.Round(time.Second)
	remainder := sd % humanize.Day
	days := int((sd - remainder) / humanize.Day)

	durStr := fmt.Sprintf("%vd%s", days, remainder.String())

	// trim zero suffixes since they are distracting
	durStr = strings.TrimSuffix(durStr, "0s")
	durStr = strings.TrimSuffix(durStr, "0m")
	durStr = strings.TrimSuffix(durStr, "0h")

	return durStr
}

type satelliteJSON struct {
	Name           string `json:"name"`
	State          string `json:"state"`
	Platform       string `json:"platform"`
	Size           string `json:"size"`
	Version        string `json:"version"`
	Selected       bool   `json:"selected"`
	Project        string `json:"project"`
	LastUsed       string `json:"last_used"`
	CacheRetention string `json:"cache_retention"`
	CloudName      string `json:"cloud_name"`
}

func (a *Satellite) printSatellitesJSON(satellites []cloud.SatelliteInstance, isOrgSelected bool) {
	jsonSats := make([]satelliteJSON, len(satellites))
	for i, s := range satellites {
		selected := s.Name == a.cli.Cfg().Satellite.Name && isOrgSelected
		jsonSats[i] = satelliteJSON{
			Name:           s.Name,
			Size:           s.Size,
			State:          s.State,
			Platform:       s.Platform,
			Version:        s.Version,
			Selected:       selected,
			LastUsed:       s.LastUsed.String(),
			CacheRetention: s.CacheRetention.String(),
			CloudName:      s.CloudName,
		}
	}
	b, err := json.MarshalIndent(jsonSats, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal json: %s", err.Error()) // unlikely
	}
	fmt.Println(string(b))
}

func (a *Satellite) actionLaunch(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteLaunch")

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)
	ffs := a.featureFlags.Value()
	size := a.size
	platform := a.platform
	window := a.maintenanceWindow
	version := a.version

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if platform != "" && !cloud.ValidSatellitePlatform(platform) {
		return errors.Errorf("not a valid platform: %q", platform)
	}
	if size != "" && !cloud.ValidSatelliteSize(size) {
		return errors.Errorf("not a valid size: %q", size)
	}

	if window == "" {
		window = "02:00"
	}

	zone, offset := time.Now().Zone()
	localWindow := window
	if window != "" {
		window, err = cloud.LocalMaintenanceWindowToUTC(window, time.FixedZone(zone, offset))
		if err != nil {
			return err
		}
	}

	a.cli.Console().Printf("Launching Satellite %q with auto-updates set to run at %s (%s)\n",
		a.cli.Flags().SatelliteName, localWindow, zone)
	a.cli.Console().Printf("This may take a few minutes...\n")

	// Collect info to help with printing a richer message in the beginning of the build or on failure to reserve satellite due to missing build minutes.
	if err = a.cli.CollectBillingInfo(cliCtx.Context, cloudClient, orgName); err != nil {
		a.cli.Console().DebugPrintf("failed to get billing plan info, error is %v\n", err)
	}

	err = cloudClient.LaunchSatellite(cliCtx.Context, cloud.LaunchSatelliteOpt{
		Name:                    a.cli.Flags().SatelliteName,
		OrgName:                 orgName,
		Platform:                platform,
		Size:                    size,
		PinnedVersion:           version,
		MaintenanceWindowStart:  window,
		FeatureFlags:            ffs,
		MaintenanceWeekendsOnly: a.maintenaceWeekendsOnly,
		CloudName:               a.cloudName,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			a.cli.Console().Printf("Operation interrupted. Satellite should finish launching in background (if server received request).\n")
			return nil
		}
		return errors.Wrapf(err, "failed to create satellite %s", a.cli.Flags().SatelliteName)
	}
	a.cli.Console().Printf("...Done\n")

	err = a.useSatellite(cliCtx, a.cli.Flags().SatelliteName, orgName)
	if err != nil {
		return errors.Wrap(err, "could not configure satellite for use")
	}
	a.cli.Console().Printf("The satellite %s has been automatically selected for use. To go back to using local builds you can use\n\n\tearthly satellite unselect\n\n", a.cli.Flags().SatelliteName)

	return nil
}

func (a *Satellite) actionSatelliteList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteList")

	if cliCtx.NArg() != 0 {
		return errors.New("command does not accept any arguments")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satellites, err := cloudClient.ListSatellites(cliCtx.Context, orgName)
	if err != nil {
		return err
	}

	// a.cli.Cfg().Satellite.Org is deprecated, but we can still check it here for compatability
	// with config files that may still have it set
	isOrgSelected := a.cli.Cfg().Satellite.Org == orgName || a.cli.Cfg().Global.Org == orgName

	if a.printJSON {
		a.printSatellitesJSON(satellites, isOrgSelected)
	} else {
		a.printSatellitesTable(satellites, isOrgSelected)
	}

	return nil
}

func (a *Satellite) actionRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteRemove")

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satellites, err := cloudClient.ListSatellites(cliCtx.Context, orgName)
	if err != nil {
		return err
	}

	var sat *cloud.SatelliteInstance
	for _, s := range satellites {
		if a.cli.Flags().SatelliteName == s.Name {
			sat = &s
			break
		}
	}
	if sat == nil {
		return fmt.Errorf("could not find %q for deletion", a.cli.Flags().SatelliteName)
	}

	isOffline := sat.State == cloud.SatelliteStatusSleep || sat.State == cloud.SatelliteStatusOffline
	if !a.force && !isOffline {
		a.cli.Console().Printf("")
		a.cli.Console().Printf("Cannot destroy a running satellite.")
		a.cli.Console().Printf("Please sleep the satellite first, or use the --force flag.")
		a.cli.Console().Printf("Note that force removing a satellite may interrupt ongoing builds.")
		a.cli.Console().Printf("")
		return errors.New("satellite is running")
	}

	a.cli.Console().Printf("Destroying Satellite %q. This may take a few minutes...\n", a.cli.Flags().SatelliteName)
	err = cloudClient.DeleteSatellite(cliCtx.Context, a.cli.Flags().SatelliteName, orgName, a.force)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			a.cli.Console().Printf("Operation interrupted. Satellite should finish destroying in background (if server received request).\n")
			return nil
		}
		return errors.Wrapf(err, "failed to delete satellite %s", a.cli.Flags().SatelliteName)
	}
	a.cli.Console().Printf("...Done\n")

	if a.cli.Flags().SatelliteName == a.cli.Cfg().Satellite.Name {
		err = a.useSatellite(cliCtx, "", "")
		if err != nil {
			return errors.Wrapf(err, "failed unselecting satellite")
		}
		a.cli.Console().Printf("Satellite has also been unselected\n")
	}
	return nil
}

func (a *Satellite) actionInspect(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteInspect")

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	satelliteToInspect := cliCtx.Args().Get(0)
	selectedSatellite := a.cli.Cfg().Satellite.Name

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, orgID, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satelliteToInspectName, err := base.GetSatelliteName(cliCtx.Context, orgName, satelliteToInspect, cloudClient)
	if err != nil {
		return err
	}

	satellite, err := cloudClient.GetSatellite(cliCtx.Context, satelliteToInspectName, orgName)
	if err != nil {
		return err
	}

	token, err := cloudClient.GetAuthToken(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}

	a.cli.Flags().BuildkitdSettings.UseTCP = true
	a.cli.Flags().BuildkitdSettings.UseTLS = a.cli.Cfg().Global.TLSEnabled
	a.cli.Flags().BuildkitdSettings.Timeout = 30 * time.Second
	a.cli.Flags().BuildkitdSettings.SatelliteToken = token
	a.cli.Flags().BuildkitdSettings.SatelliteName = satelliteToInspectName
	a.cli.Flags().BuildkitdSettings.SatelliteIsManaged = satellite.IsManaged
	a.cli.Flags().BuildkitdSettings.SatelliteDisplayName = satelliteToInspect
	a.cli.Flags().BuildkitdSettings.SatelliteOrgID = orgID // must be the ID and not name, due to satellite-proxy requirements
	satelliteAddress := a.cli.Flags().SatelliteAddress
	if satelliteAddress == "" {
		// A self-hosted satellite uses its own address
		satelliteAddress = fmt.Sprintf("tcp://%s", satellite.Address)
	}
	a.cli.Flags().BuildkitdSettings.BuildkitAddress = satelliteAddress

	selected := "No"
	if selectedSatellite == satelliteToInspect {
		selected = "Yes"
	}

	size := satellite.Size
	if !satellite.IsManaged {
		size = "self-hosted"
	}

	a.cli.Console().Printf("Address: %s", satellite.Address)
	a.cli.Console().Printf("State: %s", satellite.State)
	a.cli.Console().Printf("Platform: %s", satellite.Platform)
	a.cli.Console().Printf("Size: %s", size)
	a.cli.Console().Printf("Last Used: %s", satellite.LastUsed.In(time.Local))
	a.cli.Console().Printf("Cache Duration: %s", durationWithDaysPart(satellite.CacheRetention))
	pinned := ""
	if satellite.VersionPinned {
		pinned = " (pinned)"
	}
	a.cli.Console().Printf("Version: %s%s", satellite.Version, pinned)
	if satellite.RevisionID > 0 {
		a.cli.Console().Printf("Revision: %d", satellite.RevisionID)
	}
	if len(satellite.FeatureFlags) > 0 {
		a.cli.Console().Printf("Feature Flags: %+v", satellite.FeatureFlags)
	}
	if satellite.MaintenanceWindowStart != "" {
		zone := time.FixedZone(time.Now().Zone()) // Important not to use this instead of time.Local
		mwStart, err := cloud.UTCMaintenanceWindowToLocal(satellite.MaintenanceWindowStart, zone)
		if err != nil {
			return errors.Wrap(err, "failed converting maintenance window start to local time")
		}
		mwEnd, err := cloud.UTCMaintenanceWindowToLocal(satellite.MaintenanceWindowEnd, zone)
		if err != nil {
			return errors.Wrap(err, "failed converting maintenance window end to local time")
		}
		var weekends string
		if satellite.MaintenanceWeekendsOnly {
			weekends = " (weekends only)"
		}
		a.cli.Console().Printf("Maintenance Window: [%s - %s]%s", mwStart, mwEnd, weekends)
	}
	a.cli.Console().Printf("Currently selected: %s", selected)
	a.cli.Console().Printf("")

	if satellite.State == cloud.SatelliteStatusOperational {
		if satellite.Certificate != nil {
			cleanup, err := buildkitd.ConfigureSatelliteTLS(&a.cli.Flags().BuildkitdSettings, satellite)
			if err != nil {
				return errors.Wrap(err, "failed configuring satellite tls")
			}
			defer cleanup()
		}
		err = buildkitd.PrintSatelliteInfo(cliCtx.Context, a.cli.Console(), a.cli.App().Version, a.cli.Flags().BuildkitdSettings, a.cli.Flags().InstallationName)
		if err != nil {
			return errors.Wrap(err, "failed checking buildkit info")
		}
	} else {
		a.cli.Console().Printf("More info available when Satellite is awake.")
		if satellite.State == cloud.SatelliteStatusSleep || satellite.State == cloud.SatelliteStatusOffline {
			// Only instruct the user to run this if the satellite is asleep or offline.
			// Otherwise, satellite may be updating, still starting, etc.
			a.cli.Console().Printf("")
			a.cli.Console().Printf("    earthly satellite --org %s wake %s", orgName, satelliteToInspect)
			a.cli.Console().Printf("")
		}
	}
	return nil
}

func (a *Satellite) actionSelect(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteSelect")

	if cliCtx.NArg() == 0 {
		if a.cli.Cfg().Satellite.Name == "" {
			a.cli.Console().Printf("No satellite selected\n\n")
		} else {
			a.cli.Console().Printf("Selected satellite: %s\n\n", a.cli.Cfg().Satellite.Name)
		}
		_ = cli.ShowCommandHelp(cliCtx, cliCtx.Command.Name)
		return errors.New("satellite name is required")
	}

	if cliCtx.NArg() > 1 {
		_ = cli.ShowCommandHelp(cliCtx, cliCtx.Command.Name)
		return errors.New(fmt.Sprintf("can only provide 1 satellite name, %d provided", cliCtx.NArg()))
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	// This could be replaced with an base.GetSatelliteName() call, if we did not care about printing a list
	// after the command was run. Its done this way to save some API calls.
	found := false
	satelliteName := ""
	satellites, err := cloudClient.ListSatellites(cliCtx.Context, orgName)
	if err != nil {
		return err
	}
	for _, s := range satellites {
		if a.cli.Flags().SatelliteName == s.Name {
			found = true
			satelliteName = s.Name
		}
	}

	if !found {
		return fmt.Errorf("no satellite named %q found", a.cli.Flags().SatelliteName)
	}

	err = a.useSatellite(cliCtx, satelliteName, orgName)
	if err != nil {
		return errors.Wrapf(err, "could not select satellite %s", a.cli.Flags().SatelliteName)
	}

	a.printSatellitesTable(satellites, true)
	return nil
}

func (a *Satellite) actionUnselect(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteUnselect")

	if cliCtx.NArg() != 0 {
		return errors.New("command does not accept any arguments")
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)

	err := a.useSatellite(cliCtx, "", "")
	if err != nil {
		return errors.Wrap(err, "could not unselect satellite")
	}

	return nil
}

func (a *Satellite) actionWake(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteWake")

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satName, err := base.GetSatelliteName(cliCtx.Context, orgName, a.cli.Flags().SatelliteName, cloudClient)
	if err != nil {
		return err
	}

	sat, err := cloudClient.GetSatellite(cliCtx.Context, satName, orgName)
	if err != nil {
		return err
	}

	if sat.State == cloud.SatelliteStatusOperational {
		a.cli.Console().Printf("%s is already awake.", a.cli.Flags().SatelliteName)
	}

	out := cloudClient.WakeSatellite(cliCtx.Context, satName, orgName)
	err = base.ShowSatelliteLoading(a.cli.Console(), a.cli.Flags().SatelliteName, out)
	if err != nil {
		return errors.Wrap(err, "failed waiting for satellite wake")
	}

	return nil
}

func (a *Satellite) actionSleep(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteSleep")

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	satName, err := base.GetSatelliteName(cliCtx.Context, orgName, a.cli.Flags().SatelliteName, cloudClient)
	if err != nil {
		return err
	}

	out := cloudClient.SleepSatellite(cliCtx.Context, satName, orgName)
	err = showSatelliteStopping(a.cli.Console(), a.cli.Flags().SatelliteName, out)
	if err != nil {
		return errors.Wrap(err, "failed waiting for satellite wake")
	}

	return nil
}

func (a *Satellite) actionUpdate(cliCtx *cli.Context) error {
	a.cli.SetCommandName("satelliteUpdate")

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	a.cli.Flags().SatelliteName = cliCtx.Args().Get(0)
	window := a.maintenanceWindow
	ffs := a.featureFlags.Value()
	dropCache := a.dropCache
	version := a.version
	size := a.size
	platform := a.platform

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, _, err := a.cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	sat, err := cloudClient.GetSatellite(cliCtx.Context, a.cli.Flags().SatelliteName, orgName)
	if err != nil {
		return errors.Wrap(err, "failed getting satellite")
	}

	if window != "" {
		window, err = cloud.LocalMaintenanceWindowToUTC(window, time.Local)
		if err != nil {
			return err
		}
		z, _ := time.Now().Zone()
		a.cli.Console().Printf("Auto-update maintenance window set to %s (%s)\n", a.maintenanceWindow, z)
	}

	if size != "" && !cloud.ValidSatelliteSize(size) {
		return errors.Errorf("not a valid size: %q", size)
	}

	if platform != "" && !cloud.ValidSatellitePlatform(platform) {
		return errors.Errorf("not a valid platform: %q", platform)
	}

	if sat.State != cloud.SatelliteStatusSleep {
		if !a.force {
			a.cli.Console().Printf("")
			a.cli.Console().Printf("The satellite must be asleep to start the update.")
			a.cli.Console().Printf("You can re-run this command with the `--force` flag to force the satellite asleep and start the update now.")
			a.cli.Console().Printf("Note that Putting the satellite to sleep will interrupt any running builds.")
			a.cli.Console().Printf("")
			return errors.New("update aborted: satellite is not asleep.")
		}
		out := cloudClient.SleepSatellite(cliCtx.Context, sat.Name, orgName)
		err = showSatelliteStopping(a.cli.Console(), a.cli.Flags().SatelliteName, out)
		if err != nil {
			return errors.Wrap(err, "failed waiting for satellite to sleep")
		}
	}

	err = cloudClient.UpdateSatellite(cliCtx.Context, cloud.UpdateSatelliteOpt{
		Name:                    sat.Name,
		OrgName:                 orgName,
		PinnedVersion:           version,
		MaintenanceWindowStart:  window,
		MaintenanceWeekendsOnly: a.maintenaceWeekendsOnly,
		DropCache:               dropCache,
		FeatureFlags:            ffs,
		Size:                    size,
		Platform:                platform,
	})
	if err != nil {
		return errors.Wrap(err, "failed starting satellite update")
	}

	a.cli.Console().Printf("Update now running on satellite %q...\n", a.cli.Flags().SatelliteName)
	return nil
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
