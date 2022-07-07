package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) secretCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "set",
			Usage: "Stores a secret in the secrets store",
			UsageText: "earthly [options] secret set <path> <value>\n" +
				"   earthly [options] secret set --file <local-path> <path>\n" +
				"   earthly [options] secret set --file <local-path> <path>",
			Action: app.actionSecretsSet,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "file",
					Aliases:     []string{"f"},
					Usage:       "Stores secret stored in file",
					Destination: &app.secretFile,
				},
				&cli.BoolFlag{
					Name:        "stdin",
					Aliases:     []string{"i"},
					Usage:       "Stores secret read from stdin",
					Destination: &app.secretStdin,
				},
			},
		},
		{
			Name:      "get",
			Action:    app.actionSecretsGet,
			Usage:     "Retrieve a secret from the secrets store",
			UsageText: "earthly [options] secret get [options] <path>",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Aliases:     []string{"n"},
					Usage:       "Disable newline at the end of the secret",
					Destination: &app.disableNewLine,
				},
			},
		},
		{
			Name:      "ls",
			Usage:     "List secrets in the secrets store",
			UsageText: "earthly [options] secret ls [<path>]",
			Action:    app.actionSecretsList,
		},
		{
			Name:      "rm",
			Usage:     "Removes a secret from the secrets store",
			UsageText: "earthly [options] secret rm <path>",
			Action:    app.actionSecretsRemove,
		},
	}
}

func (app *earthlyApp) actionSecretsList(cliCtx *cli.Context) error {
	app.commandName = "secretsList"

	path := "/"
	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	} else if cliCtx.NArg() == 1 {
		path = cliCtx.Args().Get(0)
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	paths, err := cloudClient.List(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to list secret")
	}
	for _, path := range paths {
		fmt.Println(path)
	}
	return nil
}

func (app *earthlyApp) actionSecretsGet(cliCtx *cli.Context) error {
	app.commandName = "secretsGet"
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	data, err := cloudClient.Get(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to get secret")
	}
	fmt.Printf("%s", data)
	if !app.disableNewLine {
		fmt.Printf("\n")
	}
	return nil
}

func (app *earthlyApp) actionSecretsRemove(cliCtx *cli.Context) error {
	app.commandName = "secretsRemove"
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cloudClient.Remove(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to remove secret")
	}
	return nil
}

func (app *earthlyApp) actionSecretsSet(cliCtx *cli.Context) error {
	app.commandName = "secretsSet"
	var path string
	var value string
	if app.secretFile == "" && !app.secretStdin {
		if cliCtx.NArg() != 2 {
			return errors.New("invalid number of arguments provided")
		}
		path = cliCtx.Args().Get(0)
		value = cliCtx.Args().Get(1)
	} else if app.secretStdin {
		if app.secretFile != "" {
			return errors.New("only one of --file or --stdin can be used at a time")
		}
		if cliCtx.NArg() != 1 {
			return errors.New("invalid number of arguments provided")
		}
		path = cliCtx.Args().Get(0)
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read from stdin")
		}
		value = string(data)
	} else {
		if cliCtx.NArg() != 1 {
			return errors.New("invalid number of arguments provided")
		}
		path = cliCtx.Args().Get(0)
		data, err := os.ReadFile(app.secretFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read secret from %s", app.secretFile)
		}
		value = string(data)
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cloudClient.Set(cliCtx.Context, path, []byte(value))
	if err != nil {
		return errors.Wrap(err, "failed to set secret")
	}
	return nil
}
