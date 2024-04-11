package subcmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cmd/earthly/helper"
)

type Github struct {
	cli CLI

	Token  string
	GHOrg  string
	GHRepo string
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
					Description: `This command sets the configuration to create a new Github-Earthly integration to trigger satellite builds from GHA (GitHub Actions).
Integration can be done at two levels: org-wide and per repository. 
The provided token must have enough permissions to register webhook and to create Github self hosted runners in those two different scenarios.`,
					UsageText: "earthly github add --org <org> [--repo <repo>] --token <token>",
					Action:    a.actionAdd,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "org",
							Usage:       "The name of the GitHub organization to set an integration with",
							Required:    true,
							Destination: &a.GHOrg,
						},
						&cli.StringFlag{
							Name:        "repo",
							Usage:       "The name of the GitHub repository to set an integration with",
							Destination: &a.GHRepo,
						},
						&cli.StringFlag{
							Name:        "token",
							Usage:       "The GitHub token used for the integration",
							Destination: &a.Token,
						},
					},
				},
			},
		},
	}
}

func (a *Github) actionAdd(cliCtx *cli.Context) error {
	a.cli.SetCommandName("githubAdd")
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	err = cloudClient.SetGithubToken(cliCtx.Context, a.cli.OrgName(), a.GHOrg, a.GHRepo, a.Token)
	if err != nil {
		return fmt.Errorf("error found running github add: %w", err)
	}
	a.cli.Console().Printf("GitHub integration successfully created")
	return nil
}
