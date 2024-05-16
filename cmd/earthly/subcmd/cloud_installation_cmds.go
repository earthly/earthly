package subcmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/fatih/color"
	"os"
	"text/tabwriter"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/helper"
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
			Aliases:     []string{"clouds"},
			Usage:       "Configure Cloud Installations for BYOC plans",
			UsageText:   "earthly cloud (install|use|ls|rm)",
			Description: "Configure Cloud Installations for BYOC plans.",
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
					UsageText:   "earthly cloud install <cloud-name>",
					Action:      a.install,
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
					Aliases:     []string{"list"},
					Usage:       "List available Cloud Installation",
					Description: "List available Cloud Installation.",
					UsageText:   "earthly cloud ls",
					Action:      a.list,
				},
				{
					Name:        "rm",
					Aliases:     []string{"remove"},
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
		return errors.New("cloud name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single cloud name is supported")
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

	installation, err := c.getInstallationDataFromCloudFormation(ctx, cloudName)
	if err != nil {
		return err
	}

	c.cli.Console().Printf("Configuring new Cloud Installation: %s. Please wait...", cloudName)

	install, err := cloudClient.ConfigureCloud(ctx, orgID, installation)
	if err != nil {
		return errors.Wrap(err, "failed installing cloud")
	}

	if install.Status == cloud.CloudStatusRed || install.Status == cloud.CloudStatusYellow {
		c.cli.Console().Warnf("There is a problem with the cloud installation.")
		c.cli.Console().Warnf(install.StatusMessage)
		return errors.New("cloud installation failed validation")
	}

	c.cli.Console().Printf("...Done\n")
	c.cli.Console().Printf("Cloud Installation was successful. Current status of cloud is: %s", install.Status)
	c.cli.Console().Printf("")
	c.cli.Console().Printf("To make your new cloud the default destination for future satellite launches, run the following:")
	c.cli.Console().Printf("  earthly cloud use %s", cloudName)
	c.cli.Console().Printf("")

	return nil
}

func (c *CloudInstallation) use(cliCtx *cli.Context) error {
	c.cli.SetCommandName("useCloud")
	ctx := cliCtx.Context

	if cliCtx.NArg() == 0 {
		return errors.New("cloud name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single cloud name is supported")
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

	install, err := cloudClient.UseCloud(ctx, orgID, &cloud.CloudConfigurationOpt{
		Name:       cloudName,
		SetDefault: true,
	})
	if err != nil {
		return errors.Wrap(err, "could not select cloud")
	}

	if install.Status == cloud.CloudStatusRed || install.Status == cloud.CloudStatusYellow {
		c.cli.Console().Warnf("There is a problem with the cloud installation.")
		c.cli.Console().Warnf(install.StatusMessage)
		return errors.New("cloud Installation failed validation")
	}

	c.cli.Console().Printf("...Done\n")
	c.cli.Console().Printf("Current status is: %s", install.Status)
	c.cli.Console().Printf("The cloud '%s' will now be used as the default for all satellite launches within the org '%s'.", cloudName, orgName)

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

	c.printTable(installs)
	return nil
}

func (c *CloudInstallation) remove(cliCtx *cli.Context) error {
	c.cli.SetCommandName("removeCloud")
	ctx := cliCtx.Context

	if cliCtx.NArg() == 0 {
		return errors.New("cloud name is required")
	}
	if cliCtx.NArg() > 1 {
		return errors.New("only a single cloud name is supported")
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

func (c *CloudInstallation) printTable(installations []cloud.Installation) {
	t := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	fmt.Fprintln(t, " \tNAME\tSATELLITES\tSTATUS\t")
	for _, i := range installations {
		selected := " "
		if i.IsDefault {
			selected = "*"
		}
		var coloredStatus string
		switch i.Status {
		case cloud.CloudStatusGreen:
			coloredStatus = color.GreenString(i.Status)
		case cloud.CloudStatusYellow:
			coloredStatus = color.YellowString(i.Status)
		case cloud.CloudStatusRed:
			coloredStatus = color.RedString(i.Status)
		default:
			coloredStatus = color.HiRedString(i.Status)
		}
		suffix := ""
		if i.StatusMessage != "" {
			suffix = fmt.Sprintf(": %s", i.StatusMessage)
		}
		fullStatus := fmt.Sprintf("%s%s", coloredStatus, suffix)
		fmt.Fprintf(t, "%s\t%s\t%d\t%s\t\n", selected, i.Name, i.NumSatellites, fullStatus)
	}
	if err := t.Flush(); err != nil {
		fmt.Printf("failed to print satellites: %s", err.Error())
	}
}

func (c *CloudInstallation) getInstallationDataFromCloudFormation(ctx context.Context, stackName string) (*cloud.CloudConfigurationOpt, error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not load aws config")
	}

	client := cloudformation.NewFromConfig(awsConfig)

	describeStacksOutput, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not describe stack %s", stackName))
	}

	if len(describeStacksOutput.Stacks) != 1 {
		return nil, fmt.Errorf("unexpected number of stacks(%v) found with name %q", len(describeStacksOutput.Stacks), stackName)
	}

	stack := describeStacksOutput.Stacks[0]
	params := &cloud.CloudConfigurationOpt{}

	for _, output := range stack.Outputs {
		if output.OutputKey == nil {
			return nil, fmt.Errorf("specified stack %s has nil output key", stackName)
		}
		if output.OutputValue == nil {
			return nil, fmt.Errorf("specified stack %s has nil value for key %s", stackName, *output.OutputKey)
		}

		switch *output.OutputKey {
		case "InstallationName":
			params.Name = *output.OutputValue
		case "SshKeyName":
			params.SshKeyName = *output.OutputValue
		case "ComputeRoleArn":
			params.ComputeRoleArn = *output.OutputValue
		case "AccountId":
			params.AccountId = *output.OutputValue
		case "AllowedSubnetIds":
			params.AllowedSubnetIds = []string{*output.OutputValue}
		case "SecurityGroupId":
			params.SecurityGroupId = *output.OutputValue
		case "Region":
			params.Region = *output.OutputValue
		case "InstanceProfileArn":
			params.InstanceProfileArn = *output.OutputValue
		}
	}

	return params, nil
}
