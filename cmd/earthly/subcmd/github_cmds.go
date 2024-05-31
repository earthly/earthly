package subcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	pb "github.com/earthly/cloud-api/compute"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cmd/earthly/helper"
)

type Github struct {
	cli CLI

	Org     string
	GHOrg   string
	GHRepo  string
	GHToken string

	printJSON bool
}

type integrationJSON struct {
	GithubOrgName  string `json:"githubOrg"`
	GitHubRepoName string `json:"githubRepo"`
	CreatedBy      string `json:"createdBy"`
	CreatedAt      string `json:"createdAt"`
}

func NewGithub(cli CLI) *Github {
	return &Github{
		cli: cli,
	}
}

func (a *Github) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "gha",
			Usage:       "*experimental* Manage GitHub Actions integrations",
			Description: "*experimental* Manage GitHub Actions integrations. See https://docs.earthly.dev/earthly-cloud/satellites/gha-runners for detailed information.",
			Subcommands: []*cli.Command{
				{
					Name:        "add",
					Usage:       "Add a GHA integration",
					Description: `Creates a new integration to trigger satellite builds directly from GitHub Actions without the need of an intermediate runner. See https://docs.earthly.dev/earthly-cloud/satellites/gha-runners for detailed information.`,

					UsageText: "earthly gha add --org <earthly-org> --gh-org <github-org> [--gh-repo <github-repo>] --gh-token <github-token>",
					Action:    a.actionAdd,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the Earthly organization to set an integration with. Defaults to selected organization",
							Destination: &a.Org,
						},
						&cli.StringFlag{
							Name:        "gh-org",
							Usage:       "The name of the GitHub organization to set an integration with",
							Required:    true,
							Destination: &a.GHOrg,
						},
						&cli.StringFlag{
							Name:        "gh-repo",
							Usage:       "The name of the GitHub repository to set an integration with (blank for an organization-wide integration)",
							Destination: &a.GHRepo,
						},
						&cli.StringFlag{
							Name:        "gh-token",
							Usage:       "The GitHub token used for the integration",
							Destination: &a.GHToken,
						},
					},
				},
				{
					Name:  "remove",
					Usage: "Remove a GHA integration",
					Description: `Removes a GitHub-Earthly integration, to trigger satellite builds from GHA (GitHub Actions).
From the GitHub side, integration can be done at two levels: organization-wide and per repository. 
The provided token must have enough permissions to register webhooks and to create GitHub self hosted runners in those two scenarios.`,
					UsageText: "earthly gha remove --org <earthly_org> --gh-org <github_org> [--gh-repo <github_repo>]",
					Action:    a.actionRemove,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the Earthly organization to remove an integration from. Defaults to selected organization",
							Destination: &a.Org,
						},
						&cli.StringFlag{
							Name:        "gh-org",
							Usage:       "The name of the GitHub organization of the integration",
							Required:    true,
							Destination: &a.GHOrg,
						},
						&cli.StringFlag{
							Name:        "gh-repo",
							Usage:       "The name of the GitHub repository of the integration (blank if the integration is organization-wide)",
							Destination: &a.GHRepo,
						},
					},
				},
				{
					Name:        "ls",
					Usage:       "List GHA integrations",
					Description: `List the GitHub-Earthly integrations of a given organization.`,
					UsageText:   "earthly gha ls --org <earthly_org>",
					Action:      a.actionList,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the Earthly organization to set an integration with. Defaults to selected organization",
							Destination: &a.Org,
						},
						&cli.BoolFlag{
							Name:        "json",
							Usage:       "Prints the output in JSON format",
							Required:    false,
							Destination: &a.printJSON,
						},
					},
				},
			},
		},
	}
}

func (a *Github) actionAdd(cliCtx *cli.Context) error {
	a.cli.SetCommandName("githubAdd")
	if a.Org == "" {
		if a.cli.OrgName() == "" {
			return fmt.Errorf("coudn't determine Earthly organization")
		}
		a.Org = a.cli.OrgName()
	}
	if a.GHToken == "" {
		// Our signal handling under main() doesn't cause reading from stdin to cancel
		// as there's no way to pass app.ctx to stdin read calls.
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)
		token, err := promptToken()
		if err != nil {
			return err
		}
		a.GHToken = token
	}
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	err = cloudClient.CreateGHAIntegration(cliCtx.Context, a.Org, a.GHOrg, a.GHRepo, a.GHToken)
	if err != nil {
		return fmt.Errorf("error found running github add: %w", err)
	}
	a.cli.Console().Printf("GitHub integration successfully created")
	return nil
}

func (a *Github) actionRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("githubRemove")
	if a.Org == "" {
		if a.cli.OrgName() == "" {
			return fmt.Errorf("coudn't determine Earthly organization")
		}
		a.Org = a.cli.OrgName()
	}
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	err = cloudClient.RemoveGHAIntegration(cliCtx.Context, a.Org, a.GHOrg, a.GHRepo)
	if err != nil {
		return fmt.Errorf("error found running github remove: %w", err)
	}
	a.cli.Console().Printf("GitHub integration successfully removed")
	return nil
}

func (a *Github) actionList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("githubList")
	if a.Org == "" {
		if a.cli.OrgName() == "" {
			return fmt.Errorf("coudn't determine Earthly organization")
		}
		a.Org = a.cli.OrgName()
	}
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	integrations, err := cloudClient.ListGHAIntegrations(cliCtx.Context, a.Org)
	if err != nil {
		return fmt.Errorf("error found running github list: %w", err)
	}
	if a.printJSON {
		if err := a.printIntegrationsJSON(integrations); err != nil {
			return err
		}
	} else {
		if err := a.printIntegrationsTable(integrations); err != nil {
			return err
		}
	}
	return nil
}

func (a *Github) printIntegrationsJSON(integrations *pb.ListGHAIntegrationsResponse) error {
	jsonInts := make([]integrationJSON, len(integrations.Integrations))
	for i, integration := range integrations.Integrations {
		jsonInts[i] = integrationJSON{
			GithubOrgName:  integration.GithubOrgName,
			GitHubRepoName: integration.GithubRepoName,
			CreatedBy:      integration.CreatedBy,
			CreatedAt:      integration.CreatedAt.AsTime().Format(time.DateTime),
		}
	}
	b, err := json.MarshalIndent(jsonInts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	a.cli.Console().Printf(string(b))
	return nil
}

func (a *Github) printIntegrationsTable(integrations *pb.ListGHAIntegrationsResponse) error {
	t := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	headerRow := []string{"GH_ORG", "GH_REPO", "CREATED_BY", "CREATED_AT"}
	printRow(t, []color.Attribute{color.Reset}, headerRow)

	for _, integration := range integrations.Integrations {
		row := []string{integration.GithubOrgName, integration.GithubRepoName, integration.CreatedBy, integration.CreatedAt.AsTime().Format(time.DateTime)}
		c := []color.Attribute{color.Reset}
		printRow(t, c, row)
	}
	err := t.Flush()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}
	return nil
}

func promptToken() (string, error) {
	return promptHiddenText("Enter GH token")
}
