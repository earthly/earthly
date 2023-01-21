package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
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
			Name:      "revoke",
			Usage:     "Remove accounts from your organization",
			UsageText: "earthly [options] org revoke <path> <email> [<email> ...]",
			Action:    app.actionOrgRevoke,
		},
		{
			Name:        "invite",
			Usage:       "Invite users",
			Description: "Invite users",
			UsageText:   "earthly org [--org <organization-name>] invite [--name <recipient-name>] [--permission <permission>] [--message <message>] <email>",
			Action:      app.actionOrgInviteEmail,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "permission",
					Usage:       "The access level the new organization member will have. Can be one of: read, write, or admin.",
					Required:    false,
					Destination: &app.invitePermission,
				},
				&cli.StringFlag{
					Name:        "message",
					Usage:       "An optional message to send with the invitation email.",
					Required:    false,
					Destination: &app.inviteMessage,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "accept",
					Usage:       "Accept an invitation to join an organization",
					Description: "Accept an invitation to join an organization",
					UsageText:   "earthly org invite accept <invite-code>",
					Action:      app.actionOrgInviteAccept,
				},
				{
					Name:        "ls",
					Aliases:     []string{"list"},
					Usage:       "List all sent invitations (both pending and accepted)",
					Description: "List all pending and accepted invitations",
					UsageText:   "earthly org [--org <organization>] invite ls",
					Action:      app.actionOrgInviteList,
				},
			},
		},
		{
			Name:        "member",
			Aliases:     []string{"members"},
			Usage:       "Manage organization members",
			Description: "Manage organization members",
			UsageText:   "earthly org [--org <organization-name>] members (ls|update|rm)",
			Subcommands: []*cli.Command{
				{
					Name:        "ls",
					Aliases:     []string{"list"},
					Usage:       "List organization members and their permission level",
					Description: "List organization members and their permission level",
					UsageText:   "earthly org [--org organization] members ls",
					Action:      app.actionOrgMemberList,
				},
				{
					Name:        "update",
					Usage:       "Update an organization member's permission",
					Description: "Update an organization member's permission",
					UsageText:   "earthly org [--org organization] members update --permission <permission> <user-email>",
					Action:      app.actionOrgMemberUpdate,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "permission",
							Usage:       "Update an organization member's permission.",
							Destination: &app.userPermission,
						},
					},
				},
				{
					Name:        "rm",
					Usage:       "Remove a user from the organization",
					Description: "Remove a user from the organization",
					UsageText:   "earthly org [--org organization] members rm <user-email>",
					Action:      app.actionOrgMemberRemove,
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
	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}
	err = cloudClient.CreateOrg(cliCtx.Context, org)
	if err != nil {
		return errors.Wrap(err, "failed to create org")
	}
	return nil
}

func (app *earthlyApp) actionOrgList(cliCtx *cli.Context) error {
	app.commandName = "orgList"
	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}
	orgs, err := cloudClient.ListOrgs(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to list orgs")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, org := range orgs {
		fmt.Fprintf(w, "%s", org.Name)
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
	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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

func (app *earthlyApp) actionOrgInviteAccept(cliCtx *cli.Context) error {
	app.commandName = "orgInviteAccept"

	if cliCtx.NArg() != 1 {
		return errors.New("invite code is required")
	}

	code := cliCtx.Args().Get(0)
	if code == "" {
		return errors.New("invite code is required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	err = cloudClient.AcceptInvite(cliCtx.Context, code)
	if err != nil {
		return errors.Wrap(err, "failed to accept invite")
	}

	app.console.Printf("Invite accepted!")

	return nil
}

func (app *earthlyApp) actionOrgInviteList(cliCtx *cli.Context) error {
	app.commandName = "orgInviteList"

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	invites, err := cloudClient.ListInvites(cliCtx.Context, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to list invites")
	}

	if len(invites) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "User Email\tPermission\tCreated\tAccepted\n")
	for _, invite := range invites {
		accepted := "No"
		if !invite.AcceptedAt.IsZero() {
			accepted = invite.AcceptedAt.Format(dateFormat)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", invite.Email, invite.Permission, invite.CreatedAt.Format(dateFormat), accepted)
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionOrgInviteEmail(cliCtx *cli.Context) error {
	app.commandName = "orgInviteEmail"
	if cliCtx.NArg() == 0 {
		return errors.New("user email address required")
	}
	emails := cliCtx.Args().Slice()

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}
	permission := app.invitePermission
	if permission == "" {
		permission, err = promptInput(cliCtx.Context, "New user's permission [read/write/admin] (default=read): ")
		if err != nil {
			return errors.Wrap(err, "failed to read permission")
		}
	}
	permission = strings.ToLower(permission)
	switch permission {
	case "":
		permission = "read"
	case "r":
		permission = "read"
	case "w":
		permission = "write"
	case "a":
		permission = "admin"
	default:
	}
	switch permission {
	case "read", "write", "admin":
	default:
		return fmt.Errorf("invalid permission %s", permission)
	}

	for _, userEmail := range emails {
		if !strings.Contains(userEmail, "@") {
			return fmt.Errorf("invalid email address %s", userEmail)
		}
	}
	for _, userEmail := range emails {
		invite := &cloud.OrgInvitation{
			Email:      userEmail,
			OrgName:    orgName,
			Permission: permission,
			Message:    app.inviteMessage,
		}
		_, err = cloudClient.InviteToOrg(cliCtx.Context, invite)
		if err != nil {
			return errors.Wrapf(err, "failed to invite user %s into org", userEmail)
		}
		app.console.Printf("Invite sent to %s", userEmail)
	}

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

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}
	userEmail := cliCtx.Args().Get(1)
	err = cloudClient.RevokePermission(cliCtx.Context, path, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to revoke user from org")
	}
	return nil
}

func (app *earthlyApp) actionOrgMemberList(cliCtx *cli.Context) error {
	app.commandName = "orgMemberList"

	if cliCtx.NArg() != 0 {
		return errors.New("expected no arguments")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	members, err := cloudClient.ListOrgMembers(cliCtx.Context, orgName)
	if err != nil {
		return err
	}

	if len(members) == 0 {
		app.console.Printf("No members in %s", orgName)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, member := range members {
		fmt.Fprintf(w, "%s\t%s\n", member.UserEmail, member.Permission)
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionOrgMemberUpdate(cliCtx *cli.Context) error {
	app.commandName = "orgMemberUpdate"

	if cliCtx.NArg() < 1 {
		return errors.New("member email required")
	}

	if cliCtx.NArg() > 1 {
		return errors.New("too many arguments provided")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("member email required")
	}

	if app.userPermission == "" {
		return errors.New("permission required")
	}

	err = cloudClient.UpdateOrgMember(cliCtx.Context, orgName, userEmail, app.userPermission)
	if err != nil {
		return err
	}

	app.console.Printf("Member %q updated successfully", userEmail)

	return nil
}

func (app *earthlyApp) actionOrgMemberRemove(cliCtx *cli.Context) error {
	app.commandName = "orgMemberRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("member email required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("member email required")
	}

	err = cloudClient.RemoveOrgMember(cliCtx.Context, orgName, userEmail)
	if err != nil {
		return err
	}

	app.console.Printf("Member %q removed successfully", userEmail)

	return nil
}
