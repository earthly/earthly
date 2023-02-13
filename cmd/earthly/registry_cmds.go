package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) registryCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "login",
			Usage:       "Login and store credentials in earthly-cloud *beta*",
			Description: "Login and store credentials in earthly-cloud *beta*",
			UsageText: "earthly registry login --username <username> --password <password> [<host>]\n" +
				"	earthly registry login --org <org> --project <project> --username <username> --password <password> [<host>]\n",
			Action: app.actionRegistryLogin,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					Usage:       "The organization to which the project belongs.",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					Usage:       "The organization project in which to store registry credentials, if empty credentials will be stored under the user's secret storage.",
					Required:    false,
					Destination: &app.projectName,
				},
				&cli.StringFlag{
					Name:        "username",
					EnvVars:     []string{"EARTHLY_REGISTRY_USERNAME"},
					Usage:       "The username to use when logging into the registry.",
					Required:    true,
					Destination: &app.registryUsername,
				},
				&cli.StringFlag{
					Name:        "password",
					EnvVars:     []string{"EARTHLY_REGISTRY_PASSWORD"},
					Usage:       "The password to use when logging into the registry (use --password-stdin to prevent leaking your password via your shell history).",
					Required:    false,
					Destination: &app.registryPassword,
				},
				&cli.BoolFlag{
					Name:        "password-stdin",
					EnvVars:     []string{"EARTHLY_REGISTRY_PASSWORD_STDIN"},
					Usage:       "Read the password from stdin (recommended).",
					Required:    false,
					Destination: &app.registryPasswordStdin,
				},
			},
		},
		{
			Name:  "list",
			Usage: "List configured registries *beta*",
			UsageText: "earthly registry list\n" +
				"	earthly registry list --org <org> --project <project>\n",
			Action: app.actionRegistryList,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					Usage:       "The organization to which the project belongs.",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					Usage:       "The organization project in which to store registry credentials, if empty credentials will be stored under the user's secret storage.",
					Required:    false,
					Destination: &app.projectName,
				},
			},
		},
		{
			Name:  "logout",
			Usage: "Logout of a registry (that has credentials stored in earthly-cloud) *beta*",
			UsageText: "earthly registry logout [<host>]\n" +
				"	earthly registry login --org <org> --project <project> [<host>]\n",
			Action: app.actionRegistryLogout,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					Usage:       "The organization to which the project belongs.",
					Required:    false,
					Destination: &app.orgName,
				},
				&cli.StringFlag{
					Name:        "project",
					Usage:       "The organization project in which to store registry credentials, if empty credentials will be stored under the user's secret storage.",
					Required:    false,
					Destination: &app.projectName,
				},
			},
		},
	}
}

func (app *earthlyApp) isUserRegistryLocation() (bool, error) {
	if app.orgName == "" && app.projectName == "" {
		return true, nil
	}
	if app.orgName == "" {
		return false, fmt.Errorf("--project was specified without an --org value")
	}
	if app.projectName == "" {
		return false, fmt.Errorf("--org was specified without a --project value")
	}
	return false, nil
}

func (app *earthlyApp) getRegistriesPath() (string, error) {
	user, err := app.isUserRegistryLocation()
	if err != nil {
		return "", err
	}
	if user {
		return "/user/std/registry/", nil
	}
	return fmt.Sprintf("/%s/%s/std/registry/", app.orgName, app.projectName), nil
}

func (app *earthlyApp) actionRegistryLogin(cliCtx *cli.Context) error {
	app.commandName = "registryLogin"

	path, err := app.getRegistriesPath()
	if err != nil {
		return err
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	if cliCtx.NArg() > 1 {
		return fmt.Errorf("only a single host can be given")
	}

	host := cliCtx.Args().Get(0)
	if host == "" {
		host = "registry-1.docker.io"
	}

	if strings.Contains(host, "/") {
		return fmt.Errorf("hosts is malformed")
	}

	var password []byte
	if app.registryPasswordStdin {
		if app.registryPassword != "" {
			return fmt.Errorf("only one of  --password or --password-stdin")
		}
		password, err = io.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read from stdin")
		}
	} else {
		password = []byte(app.registryPassword)
	}
	if len(password) == 0 {
		return fmt.Errorf("password can not be empty")
	}

	err = cloudClient.SetSecret(cliCtx.Context, path+host+"/username", []byte(app.registryUsername))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, path+host+"/password", password)
	if err != nil {
		return err
	}

	return nil
}

type registryCredentials struct {
	host     string
	username string
}

func secretsToRegistryLogins(pathPrefix string, secrets []*cloud.Secret) []*registryCredentials {
	logins := []*registryCredentials{}
	for _, secret := range secrets {
		parts := strings.Split(strings.TrimPrefix(secret.Path, pathPrefix), "/")
		if len(parts) != 2 {
			continue
		}
		host := parts[0]
		key := parts[1]

		if key == "username" {
			logins = append(logins, &registryCredentials{
				host:     host,
				username: secret.Value,
			})
		}
	}
	return logins
}

func (app *earthlyApp) actionRegistryList(cliCtx *cli.Context) error {
	app.commandName = "registryList"

	path, err := app.getRegistriesPath()
	if err != nil {
		return err
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\n", "Registry", "username")
	for _, registry := range secretsToRegistryLogins(path, secrets) {
		fmt.Fprintf(w, "%s\t%s\n", registry.host, registry.username)
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionRegistryLogout(cliCtx *cli.Context) error {
	app.commandName = "registryRemove"
	path, err := app.getRegistriesPath()
	if err != nil {
		return err
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	host := cliCtx.Args().Get(0)
	if host == "" {
		host = "registry-1.docker.io"
	}

	if cliCtx.NArg() > 1 {
		return fmt.Errorf("only a single registry host can be given")
	}

	fmt.Printf("Removing login credentials for %s\n", host)
	for _, secretName := range []string{"username", "password"} {
		err = cloudClient.RemoveSecret(cliCtx.Context, path+host+"/"+secretName)
		if err != nil && !errors.Is(err, cloud.ErrNotFound) {
			return err
		}
	}

	return nil
}
