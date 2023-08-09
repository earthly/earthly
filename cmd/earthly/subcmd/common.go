package subcmd

import (
	"context"

	"github.com/dustin/go-humanize"
	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type byteSizeValue uint64

func (b *byteSizeValue) Set(s string) error {
	v, err := humanize.ParseBytes(s)
	if err != nil {
		return err
	}
	*b = byteSizeValue(v)
	return nil
}

func (b *byteSizeValue) String() string { return humanize.Bytes(uint64(*b)) }

// projectOrgName returns the specified org or retrieves the default org from the API.
func projectOrgName(cli CLI, ctx context.Context, cloudClient *cloud.Client) (string, error) {

	if configuredOrg := verifyOrg(cli); configuredOrg != "" {
		return configuredOrg, nil
	}

	userOrgs, err := cloudClient.ListOrgs(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to list organizations")
	}

	if len(userOrgs) == 0 {
		return "", errors.New("no organizations found, please specify with --org")
	} else if len(userOrgs) > 1 {
		return "", errors.New("multiple organizations found, please specify with --org")
	}

	return userOrgs[0].Name, nil
}

// verifyOrg returns orgName from cli Flags if set, else from config
func verifyOrg(cli CLI) string {
	if cli.Flags().OrgName != "" {
		return cli.Flags().OrgName
	}
	return cli.Cfg().Global.Org
}

func concatCmds(slices [][]*cli.Command) []*cli.Command {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]*cli.Command, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s)
	}

	return result
}

func getOrgAndProject(cli CLI, ctx context.Context, client *cloud.Client) (org, project string, isPersonal bool, err error) {
	org = verifyOrg(cli)
	if org == "" {
		return org, project, isPersonal, errors.Errorf("provide an org using the --org flag or `org select` command")
	}
	allOrgs, err := client.ListOrgs(ctx)
	if err != nil {
		return org, project, isPersonal, errors.Wrap(err, "failed listing orgs from cloud")
	}
	var cloudOrg *cloud.OrgDetail
	for _, o := range allOrgs {
		if o.Name == org {
			cloudOrg = o
			break
		}
	}
	if cloudOrg == nil {
		return org, project, isPersonal, errors.Errorf("not a member of org %q", org)
	}
	isPersonal = cloudOrg.Personal
	project = cli.Flags().ProjectName
	if project == "" && !cloudOrg.Personal {
		return org, project, isPersonal, errors.Errorf("the --project flag is required")
	}
	return org, project, isPersonal, nil
}
