package base

import (
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/buildkitd"
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
