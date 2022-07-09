package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) orgCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "create",
			Usage:     "Create a new organization",
			UsageText: "earthly [options] org create <org-name>",
			Action:    app.actionOrgCreate,
		},
		{
			Name:      "ls",
			Aliases:   []string{"list"},
			Usage:     "List organizations you belong to",
			UsageText: "earthly [options] org ls",
			Action:    app.actionOrgList,
		},
		{
			Name:      "list-permissions",
			Usage:     "List permissions and membership of an organization",
			UsageText: "earthly [options] org list-permissions <org-name>",
			Action:    app.actionOrgListPermissions,
		},
		{
			Name:      "invite",
			Usage:     "Invite accounts to your organization",
			UsageText: "earthly [options] org invite [options] <path> <email> [<email> ...]",
			Action:    app.actionOrgInvite,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "write",
					Usage:       "Grant write permissions in addition to read",
					Destination: &app.writePermission,
				},
			},
		},
		{
			Name:      "revoke",
			Usage:     "Remove accounts from your organization",
			UsageText: "earthly [options] org revoke <path> <email> [<email> ...]",
			Action:    app.actionOrgRevoke,
		},
	}
}

func (app *earthlyApp) orgCmdsPreview() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "invite",
			Usage:       "Invite accounts to your organization",
			Description: "Invite accounts to your organization",
			UsageText:   "earthly org [--org <organization-name>] invite <email> [--name <recipient-name>] [--permission <permission>] [--message <message>]",
			Action:      app.actionOrgInviteEmail,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "permission",
					Usage:    "The access level the new organization member will have. Can be one of: read, write, or admin.",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "message",
					Usage:    "An optional message to send with the invitation email.",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "name",
					Usage:    "The invite recipient's name (optional).",
					Required: false,
				},
			},
		},
	}
}

func (app *earthlyApp) actionOrgCreate(cliCtx *cli.Context) error {
	app.commandName = "orgCreate"
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	org := cliCtx.Args().Get(0)
	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	err = cloudClient.CreateOrg(cliCtx.Context, org)
	if err != nil {
		return errors.Wrap(err, "failed to create org")
	}
	return nil
}

func (app *earthlyApp) actionOrgList(cliCtx *cli.Context) error {
	app.commandName = "orgList"
	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	orgs, err := cloudClient.ListOrgs(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to list orgs")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, org := range orgs {
		fmt.Fprintf(w, "/%s/", org.Name)
		if org.Admin {
			fmt.Fprintf(w, "\tadmin")
		} else {
			fmt.Fprintf(w, "\tmember")
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionOrgListPermissions(cliCtx *cli.Context) error {
	app.commandName = "orgListPermissions"
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	orgs, err := cloudClient.ListOrgPermissions(cliCtx.Context, path)
	if err != nil {
		return errors.Wrap(err, "failed to list org permissions")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, org := range orgs {
		fmt.Fprintf(w, "%s\t%s", org.Path, org.User)
		if org.Write {
			fmt.Fprintf(w, "\trw")
		} else {
			fmt.Fprintf(w, "\tr")
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}

func (app *earthlyApp) actionOrgInvite(cliCtx *cli.Context) error {
	app.commandName = "orgInvite"
	if cliCtx.NArg() < 2 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		return errors.New("invitation paths must end with a slash (/)")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	userEmail := cliCtx.Args().Get(1)
	err = cloudClient.Invite(cliCtx.Context, path, userEmail, app.writePermission)
	if err != nil {
		return errors.Wrap(err, "failed to invite user into org")
	}
	return nil
}

func (app *earthlyApp) actionOrgInviteEmail(cliCtx *cli.Context) error {
	app.commandName = "orgInviteEmail"

	if cliCtx.NArg() != 1 {
		return errors.New("user email address required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := projectOrgName(cliCtx, cloudClient)
	if err != nil {
		return err
	}

	userEmail := cliCtx.Args().Get(0)
	if !strings.Contains(userEmail, "@") {
		return errors.New("invalid email address")
	}

	invite := &cloud.OrgInvitation{
		Email:   userEmail,
		OrgName: orgName,
	}

	name := cliCtx.String("name")
	if name == "" {
		name, err = promptInput(cliCtx.Context, "New user's name: ")
		if err != nil {
			return errors.Wrap(err, "failed to read name")
		}
	}
	invite.Name = name

	permission := cliCtx.String("permission")
	if permission == "" {
		permission, err = promptInput(cliCtx.Context, "New user's permission (read, write, admin): ")
		if err != nil {
			return errors.Wrap(err, "failed to read permission")
		}
	}
	invite.Permission = permission

	message := cliCtx.String("message")
	if message == "" {
		message, err = promptInput(cliCtx.Context, "Message to user: ")
		if err != nil {
			return errors.Wrap(err, "failed to read message")
		}
	}
	invite.Message = message

	_, err = cloudClient.InviteToOrg(cliCtx.Context, invite)
	if err != nil {
		return errors.Wrap(err, "failed to invite user into org")
	}

	app.console.Printf("Invite sent!\n")

	return nil
}

func (app *earthlyApp) actionOrgRevoke(cliCtx *cli.Context) error {
	app.commandName = "orgRevoke"
	if cliCtx.NArg() < 2 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		return errors.New("revoked paths must end with a slash (/)")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	userEmail := cliCtx.Args().Get(1)
	err = cloudClient.RevokePermission(cliCtx.Context, path, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to revoke user from org")
	}
	return nil
}
