package autocomplete

import (
	"context"
	"testing"

	"github.com/earthly/earthly/buildcontext"
	"github.com/earthly/earthly/conslogging"

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
		&cli.StringFlag{
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

func getPotentials(cmd string) ([]string, error) {
	logger := conslogging.Current(conslogging.NoColor, 0, conslogging.Info, false)
	gitLookup := buildcontext.NewGitLookup(logger, "")
	resolver := buildcontext.NewResolver(nil, gitLookup, logger, "", "", "", 0, "")
	return GetPotentials(context.TODO(), resolver, nil, cmd, len(cmd), getApp())
}

func TestFlagCompletion(t *testing.T) {
	matches, err := getPotentials("earthly --fl")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag ", "--fleet "}, matches)
}

func TestFlagCompletionWithPreviousFlags(t *testing.T) {
	matches, err := getPotentials("earthly --fig desertking --fla")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag "}, matches)
}

func TestFlagCompletionWithPreviousFlags2(t *testing.T) {
	matches, err := getPotentials("earthly --fig ")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{}, matches)
}

func TestFlagCompletionWithPreviousFlagsContainingEqual(t *testing.T) {
	matches, err := getPotentials("earthly --fig=desertking --fla")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag "}, matches)
}

func TestCommandCompletion(t *testing.T) {
	matches, err := getPotentials("earthly pru")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"prune "}, matches)
}

func TestCommandCompletionHidden(t *testing.T) {
	matches, err := getPotentials("earthly hid")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{}, matches)
}

func TestCommandSubCompletion(t *testing.T) {
	matches, err := getPotentials("earthly sub -")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subflag "}, matches)
}

func TestCommandSubCompletion2(t *testing.T) {
	matches, err := getPotentials("earthly sub --subflag abba --s")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subsubflag ", "--surf-the-internet "}, matches)
}

func TestCommandSubSubCompletion(t *testing.T) {
	matches, err := getPotentials("earthly sub --subflag abba --sub")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subsubflag "}, matches)
}

func TestCommandSubSubCompletion2(t *testing.T) {
	matches, err := getPotentials("earthly sub --subflag abba ")
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"dancing-queen "}, matches)
}
