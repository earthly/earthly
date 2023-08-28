package subcmd

import (
	"fmt"
	"os"

	"github.com/earthly/earthly/config"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

type Config struct {
	cli CLI

	dryRun bool
}

func NewConfig(cli CLI) *Config {
	return &Config{
		cli: cli,
	}
}

func (a *Config) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "config",
			Usage:  "Edits your Earthly configuration file",
			Action: a.action,
			UsageText: `Examples of common settings:
		
						Set your cache size:
		
							config global.cache_size_mb 1234
		
						Set additional buildkit args, using a YAML array:
		
							config global.buildkit_additional_args '["userns", "--host"]'
		
						Set a key containing a period:
		
							config 'git."example.com".password' hunter2
		
						Set up a whole custom git repository for a server called example.com, using a single-line YAML literal:
							* which stores git repos under /var/git/repos/name-of-repo.git
							* allows access over ssh
							* using port 2222
							* sets the username to git
							* is recognized to earthly as example.com/name-of-repo
		
							config git "{example: {pattern: 'example.com/([^/]+)', substitute: 'ssh://git@example.com:2222/var/git/repos/\$1.git', auth: ssh}}`,
			Description: `This command takes both a path and a value. It then sets them in your configuration file.
		
						As the configuration file is YAML, the key must be a valid key within the file. You can specify sub-keys by using "." to separate levels.
						If the sub-key you wish to use has a "." in it, you can quote that subsection, like this: git."github.com".
		
						Values must be valid YAML, and also be deserializable into the key you wish to assign them to.
						This means you can set higher level objects using a compact style, or single values with simple values.
		
						Only one key/value can be set per invocation.
		
						To get help with a specific key, do "config [key] --help". Or, visit https://docs.earthly.dev/earthly-config for more details.`,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "dry-run",
					Usage:       "Print the changed config file to the console instead of writing it to file",
					Destination: &a.dryRun,
				},
			},
		},
	}
}

func (a *Config) action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("config")
	if cliCtx.NArg() != 2 {
		return errors.New("invalid number of arguments provided")
	}

	args := cliCtx.Args().Slice()
	inConfig, err := config.ReadConfigFile(a.cli.Flags().ConfigPath)
	if err != nil {
		if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrapf(err, "read config")
		}
	}

	var outConfig []byte

	switch args[1] {
	case "-h", "--help":
		if err = config.PrintHelp(args[0]); err != nil {
			return errors.Wrap(err, "help")
		}
		return nil // exit now without writing any changes to config
	case "--delete":
		outConfig, err = config.Delete(inConfig, args[0])
		if err != nil {
			return errors.Wrap(err, "delete config")
		}
	default:
		// args are key/value pairs, e.g. ["global.conversion_parallelism","5"]
		outConfig, err = config.Upsert(inConfig, args[0], args[1])
		if err != nil {
			return errors.Wrap(err, "upsert config")
		}
	}

	if a.dryRun {
		fmt.Println(string(outConfig))
		return nil
	}

	err = config.WriteConfigFile(a.cli.Flags().ConfigPath, outConfig)
	if err != nil {
		return errors.Wrap(err, "write config")
	}
	a.cli.Console().Printf("Updated config file %s", a.cli.Flags().ConfigPath)

	return nil
}
