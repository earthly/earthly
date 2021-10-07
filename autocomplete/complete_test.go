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
	matches, err := GetPotentials("earthly --fl", 12, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--flag ", "--fleet "}, matches)
}

func TestCommandCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly pru", 11, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"prune "}, matches)
}

func TestCommandCompletionHidden(t *testing.T) {
	matches, err := GetPotentials("earthly hid", 11, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{}, matches)
}

func TestCommandCompletionShowHidden(t *testing.T) {
	matches, err := GetPotentials("earthly hid", 11, getApp(), true)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"hide "}, matches)
}

func TestCommandSubCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly sub -", 13, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subflag "}, matches)
}

func TestCommandSubCompletion2(t *testing.T) {
	matches, err := GetPotentials("earthly sub --subflag abba --s", 30, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subsubflag ", "--surf-the-internet "}, matches)
}

func TestCommandSubSubCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly sub --subflag abba --sub", 32, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--subsubflag "}, matches)
}

func TestCommandSubSubCompletion2(t *testing.T) {
	matches, err := GetPotentials("earthly sub --subflag abba ", 27, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"dancing-queen "}, matches)
}

func TestPathCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly .", 9, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"./", "../"}, matches)
}

func TestTargetCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly +tar", 12, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"+target "}, matches)
}

func TestTargetEnvArgCompletion(t *testing.T) {
	matches, err := GetPotentials("earthly +target ", 16, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--foo="}, matches)
}

func TestTargetEnvArgCompletionForFlagPrefix(t *testing.T) {
	matches, err := GetPotentials("earthly +target -", 17, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--foo=", "--flag ", "--fleet ", "--fig ", "--version ", "--help "}, matches)
}

func TestTargetEnvArgCompletionForArgPrefix(t *testing.T) {
	matches, err := GetPotentials("earthly +code --", 16, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--BUILDKIT_PROJECT=", "--flag ", "--fleet ", "--fig ", "--version ", "--help "}, matches)
}

func TestTargetEnvArgCompletionForArgDefault(t *testing.T) {
	matches, err := GetPotentials("earthly +earthly ", 17, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--GOOS=", "--TARGETARCH=", "--TARGETVARIANT=", "--GOARCH=", "--VARIANT=", "--GO_EXTRA_LDFLAGS=", "--EXECUTABLE_NAME=", "--EARTHLY_TARGET_TAG_DOCKER=", "--VERSION=", "--EARTHLY_GIT_HASH="}, matches)
}

func TestTargetEnvArgCompletionForMultiArgs(t *testing.T) {
	matches, err := GetPotentials("earthly +earthly --GOOS=linux --", 32, getApp(), false)
	NoError(t, err, "GetPotentials failed")
	Equal(t, []string{"--GOOS=", "--TARGETARCH=", "--TARGETVARIANT=", "--GOARCH=", "--VARIANT=", "--GO_EXTRA_LDFLAGS=", "--EXECUTABLE_NAME=", "--EARTHLY_TARGET_TAG_DOCKER=", "--VERSION=", "--EARTHLY_GIT_HASH=", "--flag ", "--fleet ", "--fig ", "--version ", "--help "}, matches)
}
