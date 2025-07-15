package base

import (
	"context"
	"fmt"

	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
)

func (cli *CLI) GetBuildkitClient(cliCtx *cli.Context) (c *client.Client, err error) {
	err = cli.InitFrontend(cliCtx)
	if err != nil {
		return nil, err
	}
	c, err = buildkitd.NewClient(cliCtx.Context, cli.Console(), cli.Flags().BuildkitdImage, cli.Flags().ContainerName, cli.Flags().InstallationName, cli.Flags().ContainerFrontend, cli.Version(), cli.Flags().BuildkitdSettings)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct new buildkit client")
	}
	return c, nil
}

func (c *CLI) OrgName() string {
	if c.Flags().OrgName != "" {
		return c.Flags().OrgName
	}
	return c.Cfg().Global.Org
}

func (c *CLI) GetSatelliteOrg(ctx context.Context, cloudClient *cloud.Client) (orgName, orgID string, err error) {
	// We are cheating here and forcing a re-auth before running any satellite commands.
	// This is because there is an issue on the backend where the token might be outdated
	// if a user was invited to an org recently after already logging-in.
	// TODO Eventually we should be able to remove this cheat.
	_, err = cloudClient.Authenticate(ctx)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to authenticate")
	}
	if c.Flags().OrgName != "" {
		orgID, err = cloudClient.GetOrgID(ctx, c.Flags().OrgName)
		if err != nil {
			return "", "", errors.Wrap(err, "invalid org provided")
		}
		return c.Flags().OrgName, orgID, nil
	}
	if c.Cfg().Global.Org != "" {
		orgID, err = cloudClient.GetOrgID(ctx, c.Cfg().Global.Org)
		if err != nil {
			return "", "", errors.Wrapf(err, "failed resolving ID for org '%s'", c.Cfg().Global.Org)
		}
		return c.Cfg().Global.Org, orgID, nil
	}
	orgName, orgID, err = cloudClient.GuessOrgMembership(ctx)
	if err != nil {
		return "", "", errors.Wrap(err, "could not guess default org")
	}
	c.Console().Warnf("Auto-selecting the default org will no longer be supported in the future.\n" +
		"You can select a default org using the command 'earthly org select',\n" +
		"or otherwise specify an org using the --org flag or EARTHLY_ORG environment variable.")
	return orgName, orgID, nil
}

func GetSatelliteName(ctx context.Context, orgName, satelliteName string, cloudClient *cloud.Client) (string, error) {
	satellites, err := cloudClient.ListSatellites(ctx, orgName)
	if err != nil {
		return "", err
	}
	for _, s := range satellites {
		if satelliteName == s.Name {
			return s.Name, nil
		}
	}

	return "", fmt.Errorf("satellite %q not found", satelliteName)
}

func (cli *CLI) reserveSatellite(ctx context.Context, cloudClient *cloud.Client, name, displayName, orgName, gitAuthor, gitConfigEmail string) error {
	console := cli.Console().WithPrefix("satellite")
	out := cloudClient.ReserveSatellite(ctx, name, orgName, gitAuthor, gitConfigEmail, false)
	err := ShowSatelliteLoading(console, displayName, out)
	if err != nil {
		return errors.Wrap(err, "failed reserving satellite for build")
	}
	return nil
}
