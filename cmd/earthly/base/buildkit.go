package base

import (
	"context"
	"fmt"

	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (cli *CLI) GetBuildkitClient(cliCtx *cli.Context, cloudClient *cloud.Client) (*client.Client, error) {
	err := cli.InitFrontend(cliCtx)
	if err != nil {
		return nil, err
	}
	err = cli.ConfigureSatellite(cliCtx, cloudClient, "", "") // no gitAuthor/gitConfigEmail is passed for non-build commands (e.g. debug_cmds.go or root_cmds.go code)
	if err != nil {
		return nil, errors.Wrapf(err, "could not construct new buildkit client")
	}

	return buildkitd.NewClient(cliCtx.Context, cli.Console(), cli.Flags().BuildkitdImage, cli.Flags().ContainerName, cli.Flags().InstallationName, cli.Flags().ContainerFrontend, cli.Version(), cli.Flags().BuildkitdSettings)
}

func (cli *CLI) ConfigureSatellite(cliCtx *cli.Context, cloudClient *cloud.Client, gitAuthor, gitConfigEmail string) error {
	if cliCtx.IsSet("buildkit-host") && cliCtx.IsSet("satellite") {
		return errors.New("cannot specify both buildkit-host and satellite")
	}
	if cliCtx.IsSet("satellite") && cli.Flags().NoSatellite {
		return errors.New("cannot specify both no-satellite and satellite")
	}
	if !cli.IsUsingSatellite(cliCtx) || cloudClient == nil {
		// If the app is not using a cloud client, or the command doesn't interact with the cloud (prune, bootstrap)
		// then pretend its all good and use your regular configuration.
		return nil
	}

	// Set up extra settings needed for buildkit RPC metadata
	if cli.Flags().SatelliteName == "" {
		cli.Flags().SatelliteName = cli.Cfg().Satellite.Name
	}
	if cli.Flags().OrgName == "" {
		cli.Flags().OrgName = cli.Cfg().Satellite.Org
	}

	cli.Flags().BuildkitdSettings.UseTCP = true
	if cli.Cfg().Global.TLSEnabled {
		// satellite connection with tls enabled does not use configuration certificates
		cli.Flags().BuildkitdSettings.ClientTLSCert = ""
		cli.Flags().BuildkitdSettings.ClientTLSKey = ""
		cli.Flags().BuildkitdSettings.TLSCA = ""
		cli.Flags().BuildkitdSettings.ServerTLSCert = ""
		cli.Flags().BuildkitdSettings.ServerTLSKey = ""
	} else {
		cli.Console().Warnf("TLS has been disabled; this should never be done when connecting to Earthly's production API\n")
	}

	orgName, orgID, err := cli.GetSatelliteOrg(cliCtx.Context, cloudClient)
	if err != nil {
		return errors.Wrap(err, "failed getting org")
	}
	satelliteName, err := GetSatelliteName(cliCtx.Context, orgName, cli.Flags().SatelliteName, cloudClient)
	if err != nil {
		return errors.Wrap(err, "failed getting satellite name")
	}
	cli.Flags().BuildkitdSettings.SatelliteName = satelliteName
	cli.Flags().BuildkitdSettings.SatelliteDisplayName = cli.Flags().SatelliteName
	cli.Flags().BuildkitdSettings.SatelliteOrgID = orgID
	if cli.Flags().SatelliteAddress != "" {
		cli.Flags().BuildkitdSettings.BuildkitAddress = cli.Flags().SatelliteAddress
	} else {
		cli.Flags().BuildkitdSettings.BuildkitAddress = containerutil.SatelliteAddress
	}
	cli.SetAnaMetaIsSat(true)
	cli.SetAnaMetaSatCurrentVersion("") // TODO

	if cli.Flags().FeatureFlagOverrides != "" {
		cli.Flags().FeatureFlagOverrides += ","
	}
	cli.Flags().FeatureFlagOverrides += "new-platform"

	token, err := cloudClient.GetAuthToken(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	cli.Flags().BuildkitdSettings.SatelliteToken = token

	// Reserve the satellite for the upcoming build.
	// This operation can take a moment if the satellite is asleep.
	err = cli.reserveSatellite(cliCtx.Context, cloudClient, satelliteName, cli.Flags().SatelliteName, orgName, gitAuthor, gitConfigEmail)
	if err != nil {
		return err
	}

	// TODO (dchw) what other settings might we want to override here?
	return nil
}

func (c *CLI) IsUsingSatellite(cliCtx *cli.Context) bool {
	if c.Flags().NoSatellite {
		return false
	}
	if cliCtx.IsSet("buildkit-host") {
		// buildkit-host takes precedence
		return false
	}
	return c.Cfg().Satellite.Name != "" || c.Flags().SatelliteName != ""
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
	satellites, err := cloudClient.ListSatellites(ctx, orgName, true)
	if err != nil {
		return "", err
	}
	for _, s := range satellites {
		if satelliteName == s.Name {
			return s.Name, nil
		}
	}

	pipelines, err := GetAllPipelinesForAllProjects(ctx, orgName, cloudClient)
	if err != nil {
		return "", err
	}
	for _, p := range pipelines {
		if satelliteName == PipelineSatelliteName(&p) {
			return p.SatelliteName, nil
		}
	}

	return "", fmt.Errorf("satellite %q not found", satelliteName)
}

func (cli *CLI) reserveSatellite(ctx context.Context, cloudClient *cloud.Client, name, displayName, orgName, gitAuthor, gitConfigEmail string) error {
	console := cli.Console().WithPrefix("satellite")
	_, isCI := analytics.DetectCI(cli.Flags().EarthlyCIRunner)
	out := cloudClient.ReserveSatellite(ctx, name, orgName, gitAuthor, gitConfigEmail, isCI)
	err := ShowSatelliteLoading(console, displayName, out)
	if err != nil {
		return errors.Wrap(err, "failed reserving satellite for build")
	}
	return nil
}

func GetAllPipelinesForAllProjects(ctx context.Context, orgName string, cloudClient *cloud.Client) ([]cloud.Pipeline, error) {
	projects, err := cloudClient.ListProjects(ctx, orgName)
	if err != nil {
		return nil, err
	}

	allPipelines := make([]cloud.Pipeline, 0)
	for _, pr := range projects {
		pipelines, err := cloudClient.ListPipelines(ctx, pr.Name, orgName, "")
		if err != nil {
			return nil, err
		}

		allPipelines = append(allPipelines, pipelines...)
	}

	return allPipelines, nil
}

func PipelineSatelliteName(p *cloud.Pipeline) string {
	return fmt.Sprintf("%s/%s", p.Project, p.Name)
}
