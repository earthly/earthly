package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/tabwriter"

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
				"   earthly [options] secret set --stdin <path>",
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

func (app *earthlyApp) secretCmdsPreview() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "set",
			Usage: "Stores a secret in the secrets store",
			UsageText: "earthly [options] secret set <path> <value>\n" +
				"   earthly [options] secret set --file <local-path> <path>\n" +
				"   earthly [options] secret set --stdin <path>",
			Action: app.actionSecretsSetV2,
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
			Action:    app.actionSecretsGetV2,
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
			Action:    app.actionSecretsListV2,
		},
		{
			Name:      "rm",
			Usage:     "Removes a secret from the secrets store",
			UsageText: "earthly [options] secret rm <path>",
			Action:    app.actionSecretsRemoveV2,
		},
		{
			Name:      "migrate",
			Usage:     "Migrate existing secrets into the new project-based structure",
			UsageText: "earthly [options] secret --org <organization> --project <project> migrate <source-organization>",
			Action:    app.actionSecretsMigrate,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "dry-run",
					Aliases:     []string{"d"},
					Usage:       "Output what the command will do without actually doing it",
					Destination: &app.dryRun,
				},
			},
		},
		{
			Name:      "permission",
			Aliases:   []string{"permissions"},
			Usage:     "Manage user-level secret permissions.",
			UsageText: "earthly [options] secret permission (ls|set|rm)",
			Subcommands: []*cli.Command{
				{
					Name:      "ls",
					Usage:     "List any user secret permissions.",
					UsageText: "earthly [options] secret permission ls <path>",
					Action:    app.actionSecretPermsList,
				},
				{
					Name:      "rm",
					Usage:     "Remove a user secret permission.",
					UsageText: "earthly [options] secret permission rm <path> <user-email>",
					Action:    app.actionSecretPermsRemove,
				},
				{
					Name:      "set",
					Usage:     "Create or update a user secret permission.",
					UsageText: "earthly [options] secret permission set <path> <user-email> <permission>",
					Action:    app.actionSecretPermsSet,
				},
			},
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

func (app *earthlyApp) actionSecretsListV2(cliCtx *cli.Context) error {
	app.commandName = "secretsList"

	path := "/"

	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	} else if cliCtx.NArg() == 1 {
		path = cliCtx.Args().Get(0)
	}

	if app.orgName == "" {
		return errors.New("invalid organization name")
	}

	if app.projectName == "" {
		return errors.New("invalid project name")
	}

	path = app.fullSecretPath(path)

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}

	for _, secret := range secrets {
		display := secret.Path
		if !strings.HasPrefix(display, "/user/") {
			prefix := fmt.Sprintf("/%s/%s/", app.orgName, app.projectName)
			display = strings.TrimPrefix(display, prefix)
		}
		app.console.Printf(display)
	}

	return nil
}

func (app *earthlyApp) actionSecretsGetV2(cliCtx *cli.Context) error {
	app.commandName = "secretsGet"

	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	path := cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	path = app.fullSecretPath(path)

	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to get secret")
	}

	if len(secrets) == 0 {
		return errors.New("no secret found for that path")
	}

	fmt.Print(secrets[0].Value)
	if !app.disableNewLine {
		fmt.Printf("\n")
	}

	return nil
}

func (app *earthlyApp) actionSecretsRemoveV2(cliCtx *cli.Context) error {
	app.commandName = "secretsRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	path := cliCtx.Args().Get(0)

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	path = app.fullSecretPath(path)

	err = cloudClient.RemoveSecret(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to remove secret")
	}

	return nil
}

func (app *earthlyApp) actionSecretsSetV2(cliCtx *cli.Context) error {
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

	path = app.fullSecretPath(path)

	err = cloudClient.SetSecret(cliCtx.Context, path, []byte(value))
	if err != nil {
		return errors.Wrap(err, "failed to set secret")
	}

	return nil
}

func (app *earthlyApp) fullSecretPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if strings.HasPrefix(path, "/user") {
		return path
	}

	// TODO: These values will eventually come from the new PROJECT command (if
	// one is present). For now, we can use the flag/env values as a temporary
	// measure.
	return fmt.Sprintf("/%s/%s%s", app.orgName, app.projectName, path)
}

func (app *earthlyApp) actionSecretPermsList(cliCtx *cli.Context) error {
	app.commandName = "secretPermissionList"

	if cliCtx.NArg() != 1 {
		return errors.New("secret path is required")
	}

	path := cliCtx.Args().Get(0)
	path = app.fullSecretPath(path)

	if strings.Contains(path, "/user") {
		return errors.New("user secrets don't support permissions")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	perms, err := cloudClient.ListSecretPermissions(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to list permissions")
	}

	if len(perms) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "User Email\tPermission\tCreated\n")
	for _, perm := range perms {
		fmt.Fprintf(w, "%s\t%s\t%s\n", perm.UserEmail, perm.Permission, perm.CreatedAt.Format(dateFormat))
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionSecretPermsRemove(cliCtx *cli.Context) error {
	app.commandName = "secretPermissionRemove"

	if cliCtx.NArg() != 2 {
		return errors.New("secret path and user email are required")
	}

	path := cliCtx.Args().Get(0)
	path = app.fullSecretPath(path)

	if strings.Contains(path, "/user") {
		return errors.New("user secrets don't support permissions")
	}

	userEmail := cliCtx.Args().Get(1)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	err = cloudClient.RemoveSecretPermission(cliCtx.Context, path, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to remove permission")
	}

	return nil
}

func (app *earthlyApp) actionSecretPermsSet(cliCtx *cli.Context) error {
	app.commandName = "secretPermissionSet"

	if cliCtx.NArg() != 3 {
		return errors.New("secret path, user email, and permission are required")
	}

	path := cliCtx.Args().Get(0)
	path = app.fullSecretPath(path)

	if strings.Contains(path, "/user") {
		return errors.New("user secrets don't support permissions")
	}

	userEmail := cliCtx.Args().Get(1)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	perm := cliCtx.Args().Get(2)
	if perm == "" {
		return errors.New("permission is required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	err = cloudClient.SetSecretPermission(cliCtx.Context, path, userEmail, perm)
	if err != nil {
		return errors.Wrap(err, "failed to set permission")
	}

	return nil
}

func (app *earthlyApp) actionSecretsMigrate(cliCtx *cli.Context) error {
	app.commandName = "secretMigrate"

	if cliCtx.NArg() != 1 {
		return errors.New("source organization required")
	}

	srcOrg := cliCtx.Args().Get(0)
	if srcOrg == "" {
		return errors.New("source organization is required")
	}

	destOrg := cliCtx.String("org")
	if destOrg == "" {
		return errors.New("destination organization is required")
	}

	destProject := cliCtx.String("project")
	if destProject == "" {
		return errors.New("destination project is required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	_, err = cloudClient.GetProject(cliCtx.Context, destOrg, destProject)
	if err != nil {
		return errors.Wrap(err, "failed to load destination project")
	}

	secretPaths, err := cloudClient.List(cliCtx.Context, fmt.Sprintf("/%s/", srcOrg))
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}

	app.console.Printf("Copying %d secrets to %s.\n", len(secretPaths), destProject)

	for _, secretPath := range secretPaths {
		val, err := cloudClient.Get(cliCtx.Context, secretPath)
		if err != nil {
			return errors.Wrapf(err, "failed to load secret %q", secretPath)
		}

		parts := strings.Split(secretPath, "/")
		newPath := "/" + path.Join(destOrg, destProject, path.Join(parts[2:]...))

		if app.verbose {
			app.console.Printf("Copying secret %q to %q\n", secretPath, newPath)
		} else {
			app.console.PrintBytes([]byte("."))
		}

		if app.dryRun {
			continue
		}

		err = cloudClient.SetSecret(cliCtx.Context, newPath, val)
		if err != nil {
			return errors.Wrap(err, "failed to set secret")
		}
	}

	if !app.verbose {
		app.console.Printf("\n")
	}

	if !app.dryRun {
		app.console.Printf("%d secrets migrated successfully!\n", len(secretPaths))
	}

	return nil
}
