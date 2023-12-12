package subcmd

import (
	"strings"

	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type AutoSkip struct {
	cli CLI

	prefix string
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
					Name:        "prefix",
					Usage:       "Prune auto-skip data by path prefix and/org target name",
					Destination: &a.prefix,
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

	if a.prefix == "" {
		return errors.New("no target specified")
	}

	// Supports the following forms:
	// * ./path/name+target
	// * ./path/n+target (where n is a prefix)
	// * +target

	var pathPrefix, target string

	if strings.HasPrefix(a.prefix, "+") {
		target = a.prefix[1:]
	} else {
		parts := strings.Split(a.prefix, "+")
		pathPrefix = parts[0]
		if len(parts) > 1 {
			target = parts[1]
		}
	}

	err = cloudClient.AutoSkipPrune(cliCtx.Context, a.cli.Flags().OrgName, a.cli.Flags().ProjectName, pathPrefix, target)
	if err != nil {
		return errors.Wrap(err, "failed to prune auto-skip hashes")
	}

	a.cli.Console().Printf("Deleted hashes for pattern: %q", a.prefix)
	return nil
}
