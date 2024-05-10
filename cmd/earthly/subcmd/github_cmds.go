package subcmd

import (
	"fmt"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cmd/earthly/helper"
)

type Github struct {
	cli CLI

	Org     string
	GHOrg   string
	GHRepo  string
	GHToken string
}

func NewGithub(cli CLI) *Github {
	return &Github{
		cli: cli,
	}
}

func (a *Github) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "github",
			Usage:       "*experimental* Manage GitHub integration",
			Description: "*experimental* Manage GitHub integration.",
			Hidden:      true,
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "Add GHA integration",
					Description: `This command sets the configuration to create a new Github-Earthly integration, to trigger satellite builds from GHA (GitHub Actions).
From the Github side, integration can be done at two levels: organization-wide and per repository. 
The provided token must have enough permissions to register webhooks and to create Github self hosted runners in those two scenarios.`,
					UsageText: "earthly github add --org <earthly_org> --gh-org <github_org> [--gh-repo <github_repo>] --gh-token <github_token>",
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
							Usage:       "The name of the GitHub repository to set an integration with",
							Destination: &a.GHRepo,
						},
						&cli.StringFlag{
							Name:        "gh-token",
							Usage:       "The GitHub token used for the integration",
							Destination: &a.GHToken,
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
	err = cloudClient.SetGithubToken(cliCtx.Context, a.Org, a.GHOrg, a.GHRepo, a.GHToken)
	if err != nil {
		return fmt.Errorf("error found running github add: %w", err)
	}
	a.cli.Console().Printf("GitHub integration successfully created")
	return nil
}

func promptToken() (string, error) {
	return promptHiddenText("Enter GH token")
}
