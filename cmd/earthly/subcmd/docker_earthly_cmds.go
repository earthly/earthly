package subcmd

import (
	"fmt"
	"os"

	"github.com/earthly/earthly/docker2earthly"

	"github.com/urfave/cli/v2"
)

type Doc2Earth struct {
	cli CLI

	earthfilePath       string
	earthfileFinalImage string
}

func NewDoc2Earth(cli CLI) *Doc2Earth {
	return &Doc2Earth{
		cli: cli,
	}
}

func (a *Doc2Earth) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "docker2earthly",
			Usage:       "Convert a Dockerfile into Earthfile",
			Description: "Converts an existing dockerfile into an Earthfile.",
			Hidden:      true, // Experimental.
			Action:      a.action,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dockerfile",
					Usage:       "Path to dockerfile input, or - for stdin",
					Value:       "Dockerfile",
					Destination: &a.cli.Flags().DockerfilePath,
				},
				&cli.StringFlag{
					Name:        "earthfile",
					Usage:       "Path to Earthfile output, or - for stdout",
					Value:       "Earthfile",
					Destination: &a.earthfilePath,
				},
				&cli.StringFlag{
					Name:        "tag",
					Usage:       "Name and tag for the built image; formatted as 'name:tag'",
					Destination: &a.earthfileFinalImage,
				},
			},
		},
	}
}

func (a *Doc2Earth) action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("docker2earthly")
	err := docker2earthly.Docker2Earthly(a.cli.Flags().DockerfilePath, a.earthfilePath, a.earthfileFinalImage)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "An Earthfile has been generated; to run it use: earthly +build; then run with docker run -ti %s\n", a.earthfileFinalImage)
	return nil
}
