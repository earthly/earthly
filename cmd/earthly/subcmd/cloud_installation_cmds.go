package subcmd

import (
	"fmt"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type CloudInstallation struct {
	cli CLI
}

func NewCloudInstallation(cli CLI) *CloudInstallation {
	return &CloudInstallation{
		cli: cli,
	}
}

func (a *CloudInstallation) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "cloud",
			Aliases:     []string{"cloud"},
			Usage:       "Configure Cloud Installations for BYOC plans",
			UsageText:   "earthly cloud (install|use|ls|rm)",
			Description: "Configure Cloud Installations for BYOC plans",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the organization the cloud installation belongs to",
					Required:    false,
					Destination: &a.cli.Flags().OrgName,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "install",
					Usage:       "Configure a new Cloud Installation",
					Description: "Configure a new Cloud Installation.",
					UsageText: "earthly cloud install <cloud-name>\n" +
						"	earthly cloud [--org <organization-name>] install <cloud-name>",
					Action: a.install,
				},
				{
					Name:        "use",
					Usage:       "Select a Cloud Installation to use for satellites",
					Description: "Select a cLoud Installation to use for satellites.",
					UsageText:   "earthly cloud use <cloud-name>",
					Action:      a.use,
				},
				{
					Name:        "ls",
					Usage:       "List available Cloud Installation",
					Description: "List available Cloud Installation.",
					UsageText:   "earthly cloud ls",
					Action:      a.list,
				},
				{
					Name:        "rm",
					Usage:       "Remove a previously installed Cloud Installation",
					Description: "Remove a previously installed Cloud Installation.",
					UsageText:   "earthly cloud rm",
					Action:      a.remove,
				},
			},
		},
	}
}

func (c *CloudInstallation) install(cliCtx *cli.Context) error {
	c.cli.SetCommandName("installCloud")
	ctx := cliCtx.Context

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	cloudName := cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(c.cli)
	if err != nil {
		return err
	}

	_, orgID, err := c.cli.GetSatelliteOrg(ctx, cloudClient)
	if err != nil {
		return err
	}

	c.cli.Console().Printf("Configuring new Cloud Installation: %s. Please wait...", cloudName)

	// TODO should this set default or no?
	install, err := cloudClient.ConfigureCloud(ctx, orgID, cloudName, false)
	if err != nil {
		return errors.Wrap(err, "could not install cloud")
	}

	if err != nil {
		return errors.Wrap(err, "failed installing cloud")
	}

	if install.Status == cloud.CloudStatusProblem {
		c.cli.Console().Warnf("There is a problem with the cloud installation. Please contact Earthly team for support.")
		return errors.New("cloud Installation failed validation")
	}

	c.cli.Console().Printf("...Done\n")
	c.cli.Console().Printf("Cloud Installation was successful. Current status is: %s", install.Status)

	return nil
}

func (c *CloudInstallation) use(cliCtx *cli.Context) error {
	c.cli.SetCommandName("useCloud")
	ctx := cliCtx.Context

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	cloudName := cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(c.cli)
	if err != nil {
		return err
	}

	orgName, orgID, err := c.cli.GetSatelliteOrg(ctx, cloudClient)
	if err != nil {
		return err
	}

	c.cli.Console().Printf("Validating Cloud Installation: %s. Please wait...", cloudName)

	install, err := cloudClient.ConfigureCloud(ctx, orgID, cloudName, true)
	if err != nil {
		return errors.Wrap(err, "could not select cloud")
	}

	if err != nil {
		return errors.Wrap(err, "failed selecting cloud")
	}

	if install.Status == cloud.CloudStatusProblem {
		c.cli.Console().Warnf("There is a problem with the cloud installation. Please contact Earthly team for support.")
		return errors.New("cloud Installation failed validation")
	}

	c.cli.Console().Printf("...Done\n")
	c.cli.Console().Printf("Current status is: %s", install.Status)
	c.cli.Console().Printf("%s will be used as the cloud for all future satellite operations across %s.", cloudName, orgName)

	return nil
}

func (c *CloudInstallation) list(cliCtx *cli.Context) error {
	c.cli.SetCommandName("listClouds")
	ctx := cliCtx.Context

	cloudClient, err := helper.NewCloudClient(c.cli)
	if err != nil {
		return err
	}

	_, orgID, err := c.cli.GetSatelliteOrg(ctx, cloudClient)
	if err != nil {
		return err
	}

	installs, err := cloudClient.ListClouds(ctx, orgID)
	if err != nil {
		return errors.Wrap(err, "could not select cloud")
	}

	for _, i := range installs {
		fmt.Printf("%+v\n", i) // TODO pretty print
	}

	return nil
}

func (c *CloudInstallation) remove(cliCtx *cli.Context) error {
	c.cli.SetCommandName("removeCloud")
	ctx := cliCtx.Context

	if cliCtx.NArg() == 0 {
		return errors.New("satellite name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single satellite name is supported")
	}

	cloudName := cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(c.cli)
	if err != nil {
		return err
	}

	_, orgID, err := c.cli.GetSatelliteOrg(ctx, cloudClient)
	if err != nil {
		return err
	}

	c.cli.Console().Printf("Removing Cloud Installation: %s. Please wait...", cloudName)

	err = cloudClient.DeleteCloud(ctx, orgID, cloudName)
	if err != nil {
		return errors.Wrap(err, "could not select cloud")
	}

	c.cli.Console().Printf("...Done\n")
	return nil
}
