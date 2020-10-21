package autocomplete

import (
	"testing"

	. "github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "flag",
		},
		&cli.BoolFlag{
			Name: "fleet",
		},
		&cli.BoolFlag{
			Name: "fig",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name: "prune",
		},
		{
			Name: "foo",
		},
	}
	return app
}

func TestFlagCompletion(t *testing.T) {

	matches, err := GetPotentials("earth --fl", 10, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag ", "--fleet "}, matches)
}

func TestCommandCompletion(t *testing.T) {
	matches, err := GetPotentials("earth pru", 9, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"prune "}, matches)
}

func TestPathCompletion(t *testing.T) {
	matches, err := GetPotentials("earth .", 7, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"./", "../"}, matches)
}
