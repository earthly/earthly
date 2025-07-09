package subcmd

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
)

type orgLister interface {
	ListOrgs(ctx context.Context) ([]*cloud.OrgDetail, error)
}

// projectOrgName returns the specified org or retrieves the default org from the API.
func projectOrgName(cli CLI, ctx context.Context, cloudClient *cloud.Client) (string, error) {

	if configuredOrg := cli.OrgName(); configuredOrg != "" {
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

func getOrgAndProject(ctx context.Context, orgFlag, projectFlag string, client orgLister, path string) (org, project string, isPersonal bool, err error) {

	allOrgs, err := client.ListOrgs(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed listing orgs from cloud")
		return
	}

	org, project = orgFlag, projectFlag

	if org == "" || strings.HasPrefix(path, "/user") {
		for _, o := range allOrgs {
			if o.Personal {
				org = o.Name
				project = ""
				isPersonal = true
				break
			}
		}
	} else {
		var found bool
		for _, o := range allOrgs {
			if o.Name == org {
				isPersonal = o.Personal
				found = true
				break
			}
		}
		if !found {
			err = errors.Errorf("not a member of org %q", org)
			return
		}
	}

	if org == "" {
		err = errors.New("provide an org using the --org flag or `org select` command")
		return
	}

	if project == "" && !isPersonal {
		err = errors.Errorf("the --project flag is required")
		return
	}

	return
}
