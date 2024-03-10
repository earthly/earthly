package subcmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/earthly/earthly/billing"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/earthly/earthly/util/stringutil"

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
					Usage:       "The name of the Earthly organization to view billing info for",
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

	orgName := a.cli.OrgName()

	if orgName == "" {
		return errors.New("organization name must be specified")
	}
	if err := a.cli.CollectBillingInfo(cliCtx.Context, cloudClient, orgName); err != nil {
		return fmt.Errorf("failed to get billing info: %w", err)
	}

	plan := billing.Plan()
	allowedArches := strings.Join(stringutil.EnumToStringArray(plan.GetAllowedArchs(), stringutil.Lower), ",")
	allowedInstances := strings.Join(stringutil.EnumToStringArray(plan.GetAllowedInstances(), stringutil.Lower), ",")

	w := new(tabwriter.Writer)
	buf := new(bytes.Buffer)
	w.Init(buf, 0, 8, 0, '\t', 0)
	fmt.Fprintf(w, "Tier:\t%s\n", stringutil.Title(plan.GetTier()))
	fmt.Fprintf(w, "Plan Type:\t%s\n", stringutil.Title(plan.GetType()))
	fmt.Fprintf(w, "Started At:\t%s\n", plan.GetStartedAt().AsTime().UTC().Format("January 2, 2006"))
	fmt.Fprintf(w, "Used Build Time:\t%s (%d minutes)\n", billing.UsedBuildTime(), int(billing.UsedBuildTime().Minutes()))
	fmt.Fprintf(w, "Max Builds Minutes:\t%s\n", valueOrUnlimited(plan.GetMaxBuildMinutes()))
	fmt.Fprintf(w, "Max Minutes per Build:\t%d\n", plan.GetMaxMinutesPerBuild())
	fmt.Fprintf(w, "Included Minutes:\t%d\n", plan.GetIncludedMinutes())
	fmt.Fprintf(w, "Max Satellites:\t%d\n", plan.GetMaxSatellites())
	fmt.Fprintf(w, "Max Hours Cache Retention:\t%s\n", valueOrUnlimited(plan.GetMaxHoursCacheRetention()))
	fmt.Fprintf(w, "Allowed Architectures:\t%s\n", allowedArches)
	fmt.Fprintf(w, "Allowed Instances:\t%s\n", allowedInstances)
	fmt.Fprintf(w, "Default Instance:\t%s\n", stringutil.Lower(plan.GetDefaultInstanceType()))
	w.Flush()
	a.cli.Console().Printf(buf.String())

	return nil
}

func valueOrUnlimited(value int32) string {
	if value != 0 {
		return fmt.Sprintf("%d", value)
	}
	return "Unlimited"
}
