package subcmd

import (
	"fmt"
	"github.com/earthly/earthly/util/stringutil"
	"strings"

	"github.com/earthly/earthly/billing"
	"github.com/earthly/earthly/cmd/earthly/helper"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Billing struct {
	cli CLI
}

func NewBilling(cli CLI) *Billing {
	return &Billing{
		cli: cli,
	}
}

func (a *Billing) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "billing",
			Aliases:     []string{"bill"},
			Description: `*experimental* View Earthly billing info.`,
			Usage:       `*experimental* View Earthly billing info`,
			UsageText:   "earthly billing (view)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the Earthly organization to to view billing info for",
					Required:    false,
					Destination: &a.cli.Flags().OrgName,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "view",
					Usage:       "View billing information for the specified organization",
					Description: "View billing information for the specified organization.",
					UsageText:   "earthly billing [--org <organization-name>] view",
					Action:      a.actionView,
				},
			},
		},
	}
}

func (a *Billing) actionView(cliCtx *cli.Context) error {
	a.cli.SetCommandName("billingView")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	if !cloudClient.IsLoggedIn(cliCtx.Context) {
		return errors.New("user must be logged in")
	}

	orgName := verifyOrg(a.cli)

	if orgName == "" {
		return errors.New("organization name must be specified")
	}
	if err := a.cli.CollectBillingInfo(cliCtx.Context, cloudClient, orgName); err != nil {
		return fmt.Errorf("failed to get billing info: %w", err)
	}

	allowedArches := strings.Join(stringutil.EnumToStringArray(billing.Plan().AllowedArchs, stringutil.Lower), ",")
	allowedInstances := strings.Join(stringutil.EnumToStringArray(billing.Plan().AllowedInstances, stringutil.Lower), ",")

	items := []*tuple{
		{key: "Tier", value: stringutil.Title(billing.Plan().Tier)},
		{key: "Plan Type", value: stringutil.Title(billing.Plan().Type)},
		{key: "Started At", value: billing.Plan().StartedAt.AsTime().UTC().Format("January 2, 2006")},
		{key: "Used Builds Minutes", value: fmt.Sprintf("%d", int(billing.UsedBuildTime().Minutes()))},
		{key: "Max Builds Minutes", value: valueOrUnlimited(billing.Plan().MaxBuildMinutes)},
		{key: "Max Minutes per Build", value: fmt.Sprintf("%d", billing.Plan().MaxMinutesPerBuild)},
		{key: "Included Minutes", value: fmt.Sprintf("%d", billing.Plan().IncludedMinutes)},
		{key: "Max Satellites", value: fmt.Sprintf("%d", billing.Plan().MaxSatellites)},
		{key: "Max Hours Cache Retention", value: valueOrUnlimited(billing.Plan().MaxHoursCacheRetention)},
		{key: "Allowed Architectures", value: allowedArches},
		{key: "Allowed Instances", value: allowedInstances},
		{key: "Default Instance", value: stringutil.Lower(billing.Plan().DefaultInstanceType)},
	}

	a.cli.Console().Printf(planInfoText(items, 5))

	return nil
}

func planInfoText(items []*tuple, padding int) string {
	maxKeyLen := 0
	for _, item := range items {
		if len(item.key) > maxKeyLen {
			maxKeyLen = len(item.key)
		}
	}

	sb := &strings.Builder{}
	for i, item := range items {
		if i != 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(fmt.Sprintf("%-*s", maxKeyLen+padding, fmt.Sprintf("%s:", item.key)))
		sb.WriteString(item.value)
	}
	return sb.String()
}

func valueOrUnlimited(value int32) string {
	if value != 0 {
		return fmt.Sprintf("%d", value)
	}
	return "Unlimited"
}

type tuple struct {
	key   string
	value string
}
