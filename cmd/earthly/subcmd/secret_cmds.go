package subcmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/earthly/earthly/util/termutil"
)

type Secret struct {
	cli CLI

	secretStdin    bool
	disableNewLine bool
	dryRun         bool
}

func NewSecret(cli CLI) *Secret {
	return &Secret{
		cli: cli,
	}
}

func (a *Secret) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "secret",
			Aliases:     []string{"secrets"},
			Description: "*beta* Manage cloud secrets.",
			Usage:       "*beta* Manage cloud secrets",
			UsageText:   "earthly [options] secrets [--org <organization-name>, --project <project>] (set|get|ls|rm|migrate|permission)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs",
					Required:    false,
					Destination: &a.cli.Flags().OrgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store secrets",
					Required:    false,
					Destination: &a.cli.Flags().ProjectName,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:  "set",
					Usage: "*beta* Stores a secret in the secrets store",
					UsageText: "earthly [options] secret set <path>\n" +
						"   earthly [options] secrets set <path> <value>\n" +
						"   earthly [options] secrets set --file <local-path> <path>\n" +
						"   earthly [options] secrets set --stdin <path>\n" +
						"\n" +
						"Security Recommendation: avoid specifying the secret <value> on the command line (to prevent storing secrets in your shell's history);\n" +
						"instead simply omit it, which will cause earthly to interactively prompt for the secret value, or use the --file or --stdin options.",
					Description: "*beta* Stores a secret in the secrets store.",
					Action:      a.actionSetV2,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "file",
							Aliases:     []string{"f"},
							Usage:       "Stores secret stored in file",
							Destination: &a.cli.Flags().SecretFile,
						},
						&cli.BoolFlag{
							Name:        "stdin",
							Aliases:     []string{"i"},
							Usage:       "Stores secret read from stdin",
							Destination: &a.secretStdin,
						},
					},
				},
				{
					Name:        "get",
					Action:      a.actionGetV2,
					Usage:       "*beta* Retrieve a secret from the secrets store",
					UsageText:   "earthly [options] secrets get [options] <path>",
					Description: "*beta* Retriece a secret from the secrets store.",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Aliases:     []string{"n"},
							Usage:       "Disable newline at the end of the secret",
							Destination: &a.disableNewLine,
						},
					},
				},
				{
					Name:        "ls",
					Usage:       "*beta* List secrets in the secrets store",
					UsageText:   "earthly [options] secrets ls [<path>]",
					Description: "*beta* List secrets in the secrets store.",
					Action:      a.actionListV2,
				},
				{
					Name:        "rm",
					Usage:       "*beta* Removes a secret from the secrets store",
					UsageText:   "earthly [options] secrets rm <path>",
					Description: "*beta* Removes a secret from the secrets store.",
					Action:      a.actionRemoveV2,
				},
				{
					Name:        "migrate",
					Usage:       "*beta* Migrate existing secrets into the new project-based structure",
					UsageText:   "earthly [options] secrets --org <organization> --project <project> migrate <source-organization>",
					Description: "*beta* Migrate existing secrets into the new project-based structure.",
					Action:      a.actionMigrate,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "dry-run",
							Aliases:     []string{"d"},
							Usage:       "Output what the command will do without actually doing it",
							Destination: &a.dryRun,
						},
					},
				},
				{
					Name:        "permission",
					Aliases:     []string{"permissions"},
					Usage:       "*beta* Manage user-level secret permissions",
					UsageText:   "earthly [options] secrets permission (ls|set|rm)",
					Description: "*beta* Manage user-level secret permissions.",
					Subcommands: []*cli.Command{
						{
							Name:        "ls",
							Usage:       "List any user secret permissions",
							UsageText:   "earthly [options] secret permission ls <path>",
							Description: "List any user secret permissions.",
							Action:      a.actionPermsList,
						},
						{
							Name:        "rm",
							Usage:       "Remove a user secret permission",
							UsageText:   "earthly [options] secret permission rm <path> <user-email>",
							Description: "Remove a user secret permission.",
							Action:      a.actionPermsRemove,
						},
						{
							Name:        "set",
							Usage:       "Create or update a user secret permission",
							UsageText:   "earthly [options] secret permission set <path> <user-email> <permission>",
							Description: "Create or update a user secret permission.",
							Action:      a.actionPermsSet,
						},
					},
				},
			},
		},
	}
}

func (a *Secret) actionListV2(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretsList")

	path := "/"

	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	} else if cliCtx.NArg() == 1 {
		path = cliCtx.Args().Get(0)
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}

	if len(secrets) == 0 {
		a.cli.Console().Printf("No secrets found")
		return nil
	}

	orgName, projectName, isPersonal, err := getOrgAndProject(cliCtx.Context, a.cli.OrgName(), a.cli.Flags().ProjectName, cloudClient, path)
	if err != nil {
		return err
	}

	for _, secret := range secrets {
		fmt.Println(secretDisplay(isPersonal, orgName, projectName, secret))
	}

	return nil
}

func secretDisplay(personal bool, org, proj string, secret *cloud.Secret) string {
	if personal && proj == "" {
		return strings.TrimPrefix(secret.Path, "/user/")
	}
	return strings.TrimPrefix(secret.Path, fmt.Sprintf("/%s/%s/", org, proj))
}

func (a *Secret) actionGetV2(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretsGet")

	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	path := cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

	secret, err := cloudClient.GetUserOrProjectSecret(cliCtx.Context, path)
	if err != nil {
		if errors.Is(err, cloud.ErrNotFound) {
			return errors.New("no secret found for that path")
		}
		return errors.Wrap(err, "failed to get secret")
	}

	fmt.Print(secret.Value)
	if !a.disableNewLine {
		fmt.Printf("\n")
	}

	return nil
}

func (a *Secret) actionRemoveV2(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretsRemove")

	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	path := cliCtx.Args().Get(0)

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

	err = cloudClient.RemoveSecret(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to remove secret")
	}

	a.cli.Console().Printf("Secret successfully deleted")

	return nil
}

func (a *Secret) actionSetV2(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretsSet")
	var path string
	var value string
	if a.cli.Flags().SecretFile == "" && !a.secretStdin {
		switch cliCtx.NArg() {
		case 1:
			path = cliCtx.Args().Get(0)
			var err error
			value, err = promptHiddenText("secret value")
			if err != nil {
				return err
			}
		case 2:
			path = cliCtx.Args().Get(0)
			value = cliCtx.Args().Get(1)
		default:
			return errors.New("invalid number of arguments provided")
		}
	} else if a.secretStdin {
		if termutil.IsTTY() {
			a.cli.Console().Printf("Reading secret from stdin; waiting for eof (ctrl-d); if you are running this interactively consider running \"earthy secret set <path>\" (without a value) instead.\n")
		}
		if a.cli.Flags().SecretFile != "" {
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
		data, err := os.ReadFile(a.cli.Flags().SecretFile)
		if err != nil {
			return errors.Wrapf(err, "failed to read secret from %s", a.cli.Flags().SecretFile)
		}
		value = string(data)
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

	err = cloudClient.SetSecret(cliCtx.Context, path, []byte(value))
	if err != nil {
		return errors.Wrap(err, "failed to set secret")
	}

	return nil
}

func (a *Secret) fullSecretPath(ctx context.Context, cloudClient *cloud.Client, path string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if strings.HasPrefix(path, "/user") {
		return path, nil
	}

	orgName, projectName, isPersonal, err := getOrgAndProject(ctx, a.cli.OrgName(), a.cli.Flags().ProjectName, cloudClient, "")
	if err != nil {
		return "", err
	}

	if isPersonal && projectName == "" && !strings.HasPrefix(path, "/user") {
		if path == "/" {
			return "/user", nil
		}
		return fmt.Sprintf("/user%s", path), nil
	}

	// TODO: These values will eventually come from the new PROJECT command (if
	//   one is present). For now, we can use the flag/env values as a temporary
	//   measure.
	return fmt.Sprintf("/%s/%s%s", orgName, projectName, path), nil
}

func (a *Secret) actionPermsList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretPermissionList")

	if cliCtx.NArg() != 1 {
		return errors.New("secret path is required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path := cliCtx.Args().Get(0)
	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

	if strings.Contains(path, "/user") {
		return errors.New("user secrets don't support permissions")
	}
	perms, err := cloudClient.ListSecretPermissions(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to list permissions")
	}

	if len(perms) == 0 {
		a.cli.Console().Printf("No permissions found for this secret")
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

func (a *Secret) actionPermsRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretPermissionRemove")

	if cliCtx.NArg() != 2 {
		return errors.New("secret path and user email are required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path := cliCtx.Args().Get(0)
	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

	if strings.Contains(path, "/user") {
		return errors.New("user secrets don't support permissions")
	}

	userEmail := cliCtx.Args().Get(1)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	err = cloudClient.RemoveSecretPermission(cliCtx.Context, path, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to remove permission")
	}

	a.cli.Console().Printf("Permission removed successfully")

	return nil
}

func (a *Secret) actionPermsSet(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretPermissionSet")

	if cliCtx.NArg() != 3 {
		return errors.New("secret path, user email, and permission are required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path := cliCtx.Args().Get(0)
	path, err = a.fullSecretPath(cliCtx.Context, cloudClient, path)
	if err != nil {
		return err
	}

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

	err = cloudClient.SetSecretPermission(cliCtx.Context, path, userEmail, perm)
	if err != nil {
		return errors.Wrap(err, "failed to set permission")
	}

	a.cli.Console().Printf("%s was granted %s permission on the secret", userEmail, perm)

	return nil
}

func (a *Secret) actionMigrate(cliCtx *cli.Context) error {
	a.cli.SetCommandName("secretMigrate")

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

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	_, err = cloudClient.GetProject(cliCtx.Context, destOrg, destProject)
	if err != nil {
		return errors.Wrap(err, "failed to load destination project")
	}

	secretPaths, err := cloudClient.List(cliCtx.Context, fmt.Sprintf("/%s/", srcOrg))
	if err != nil {
		return errors.Wrap(err, "failed to list secrets")
	}

	a.cli.Console().Printf("Copying %d secrets to %s.\n", len(secretPaths), destProject)

	for _, secretPath := range secretPaths {
		val, err := cloudClient.Get(cliCtx.Context, secretPath)
		if err != nil {
			return errors.Wrapf(err, "failed to load secret %q", secretPath)
		}

		parts := strings.Split(secretPath, "/")
		newPath := "/" + path.Join(destOrg, destProject, path.Join(parts[2:]...))

		if a.cli.Flags().Verbose {
			a.cli.Console().Printf("Copying secret %q to %q\n", secretPath, newPath)
		} else {
			a.cli.Console().PrintBytes([]byte("."))
		}

		if a.dryRun {
			continue
		}

		err = cloudClient.SetSecret(cliCtx.Context, newPath, val)
		if err != nil {
			return errors.Wrap(err, "failed to set secret")
		}
	}

	if !a.cli.Flags().Verbose {
		a.cli.Console().Printf("\n")
	}

	if !a.dryRun {
		a.cli.Console().Printf("%d secrets migrated successfully!\n", len(secretPaths))
	}

	return nil
}
