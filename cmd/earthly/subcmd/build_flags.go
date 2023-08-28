package subcmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

func (a *Build) buildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "platform",
			EnvVars: []string{"EARTHLY_PLATFORMS"},
			Usage:   "Specify the target platform to build for or this can be read from ENV VAR",
			Value:   &a.platformsStr,
		},
		&cli.StringSliceFlag{
			Name:    "build-arg",
			EnvVars: []string{"EARTHLY_BUILD_ARGS"},
			Usage:   "A build arg override, specified as <key>=[<value>]",
			Value:   &a.buildArgs,
			Hidden:  true, // Deprecated
		},
		&cli.StringSliceFlag{
			Name:    "secret",
			Aliases: []string{"s"},
			EnvVars: []string{"EARTHLY_SECRETS"},
			Usage:   "A secret override, specified as <key>=[<value>]",
			Value:   &a.secrets,
		},
		&cli.StringSliceFlag{
			Name:    "secret-file",
			EnvVars: []string{"EARTHLY_SECRET_FILES"},
			Usage:   "A secret override, specified as <key>=<path>",
			Value:   &a.secretFiles,
		},
		&cli.StringSliceFlag{
			Name:    "cache-from",
			EnvVars: []string{"EARTHLY_CACHE_FROM"},
			Usage:   "Remote docker image tags to use as readonly explicit cache (experimental)",
			Value:   &a.cacheFrom,
			Hidden:  true, // Experimental
		},
	}
}

func (a *Build) HiddenFlags() []cli.Flag {
	_, isAutocomplete := os.LookupEnv("COMP_LINE")
	flags := a.buildFlags()
	if isAutocomplete {
		// Don't hide the build flags for autocomplete.
		return flags
	}
	for _, flag := range flags {
		switch f := flag.(type) {
		case *cli.StringSliceFlag:
			f.Hidden = true
		case *cli.StringFlag:
			f.Hidden = true
		case *cli.BoolFlag:
			f.Hidden = true
		case *cli.IntFlag:
			f.Hidden = true
		}
	}

	return flags
}
