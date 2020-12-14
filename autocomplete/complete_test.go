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
		{
			Name:   "hide",
			Hidden: true,
		},
		{
			Name: "sub",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "subflag",
				},
			},
			Subcommands: []*cli.Command{
				{
					Name: "abc",
				},
				{
					Name: "abba",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "subsubflag",
						},
						&cli.BoolFlag{
							Name: "surf-the-internet",
						},
					},
					Subcommands: []*cli.Command{
						{
							Name: "dancing-queen",
						},
					},
				},
				{
					Name:   "hide",
					Hidden: true,
				},
			},
		},
	}
	return app
}

func TestFlagCompletion(t *testing.T) {

	matches, err := GetPotentials("earthly --fl", 12, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag ", "--fleet "}, matches)
}

func TestCommandCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly pru", 11, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"prune "}, matches)
}

func TestCommandCompletionHidden(t *testing.T) {
	matches, err := GetPotentials("earthly hid", 11, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{}, matches)
}

func TestCommandSubCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly sub -", 13, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subflag "}, matches)
}

func TestCommandSubCompletion2(t *testing.T) {
	matches, err := GetPotentials("earthly sub --subflag abba --s", 30, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subsubflag ", "--surf-the-internet "}, matches)
}

func TestCommandSubSubCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly sub --subflag abba --sub", 32, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subsubflag "}, matches)
}

func TestCommandSubSubCompletion2(t *testing.T) {
	matches, err := GetPotentials("earthly sub --subflag abba ", 27, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"dancing-queen "}, matches)
}

func TestPathCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly .", 9, getApp())
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"./", "../"}, matches)
}
