package subcmd

import (
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type AutoSkip struct {
	cli CLI

	path   string
	target string
	deep   bool
}

func NewAutoSkip(cli CLI) *AutoSkip {
	return &AutoSkip{
		cli: cli,
	}
}

func (a *AutoSkip) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "prune-auto-skip",
			Usage:       "Prune Earthly auto-skip data",
			Description: `Prune Earthly auto-skip hash data by hash value, organization, project, or target name.`,
			Action:      a.action,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "path",
					Usage:       "Prune auto-skip data by the specified path",
					Destination: &a.path,
				},
				&cli.StringFlag{
					Name:        "target",
					Usage:       "Prune auto-skip data by the specified target name",
					Destination: &a.target,
				},
				&cli.BoolFlag{
					Name:        "deep",
					Usage:       "Prune all sub-directories that begin with the provided path",
					Destination: &a.deep,
				},
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs",
					Required:    true,
					Destination: &a.cli.Flags().OrgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store secrets",
					Required:    true,
					Destination: &a.cli.Flags().ProjectName,
				},
			},
		},
	}
}

func (a *AutoSkip) action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("prune-auto-skip")
	if cliCtx.NArg() != 0 {
		return errors.New("invalid arguments")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	if a.path == "" && a.target == "" {
		return errors.New("no target or path specified")
	}

	count, err := cloudClient.AutoSkipPrune(
		cliCtx.Context,
		a.cli.Flags().OrgName,
		a.path,
		a.target,
		a.deep,
	)
	if err != nil {
		return errors.Wrap(err, "failed to prune auto-skip hashes")
	}

	pluralForm := "hashes"
	if count == 1 {
		pluralForm = "hash"
	}

	a.cli.Console().Printf("Deleted %d matching auto-skip %s", count, pluralForm)
	return nil
}
