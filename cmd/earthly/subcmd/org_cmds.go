package subcmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/helper"

	"github.com/earthly/earthly/config"
)

type Org struct {
	cli CLI

	invitePermission string
	inviteMessage    string
	userPermission   string

	Cfg *config.Config
}

func NewOrg(cli CLI) *Org {
	return &Org{
		cli: cli,
	}
}

func (a *Org) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "org",
			Aliases:     []string{"orgs"},
			Usage:       "Create or manage your Earthly orgs",
			Description: "Create or manage your Earthly orgs.",
			Subcommands: []*cli.Command{
				{
					Name:        "create",
					Usage:       "Create a new organization",
					UsageText:   "earthly [options] org create <org-name>",
					Description: "Create a new organization.",
					Action:      a.actionCreate,
				},
				{
					Name:        "ls",
					Aliases:     []string{"list"},
					Usage:       "List organizations you are a member or administrator of",
					UsageText:   "earthly [options] org ls",
					Description: "List organizations you are a member or administrator of.",
					Action:      a.actionList,
				},
				{
					Name:        "list-permissions",
					Usage:       "List all accounts and the paths they have permission to access under a particular organization",
					UsageText:   "earthly [options] org list-permissions <org-name>",
					Description: "List all accounts and the paths they have permission to access under a particular organization.",
					Action:      a.actionListPermissions,
				},
				{
					Name:        "revoke",
					Usage:       "Revokes a previously invited user from an organization",
					UsageText:   "earthly [options] org revoke <path> <email> [<email> ...]",
					Description: "Revokes a previously invited user from an organization.",
					Action:      a.actionRevoke,
				},
				{
					Name:        "invite",
					Usage:       "Invite users to your org",
					Description: "Invite users to your org.",
					UsageText:   "earthly org [--org <organization-name>] invite [--name <recipient-name>] [--permission <permission>] [--message <message>] <email>",
					Action:      a.actionInviteEmail,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name: "permission",
							Usage: `The access level the new organization member will have
					Can be one of: read, write, or admin.`,
							Required:    false,
							Destination: &a.invitePermission,
						},
						&cli.StringFlag{
							Name:        "message",
							Usage:       "An optional message to send with the invitation email",
							Required:    false,
							Destination: &a.inviteMessage,
						},
					},
					Subcommands: []*cli.Command{
						{
							Name:        "accept",
							Usage:       "Accept an invitation to join an organization",
							Description: "Accept an invitation to join an organization.",
							UsageText:   "earthly org invite accept <invite-code>",
							Action:      a.actionInviteAccept,
						},
						{
							Name:        "ls",
							Aliases:     []string{"list"},
							Usage:       "List all sent invitations (both pending and accepted)",
							Description: "List all pending and accepted invitations.",
							UsageText:   "earthly org [--org <organization>] invite ls",
							Action:      a.actionInviteList,
						},
					},
				},
				{
					Name:        "member",
					Aliases:     []string{"members"},
					Usage:       "Manage organization members",
					Description: "Manage organization members.",
					UsageText:   "earthly org [--org <organization-name>] members (ls|update|rm)",
					Subcommands: []*cli.Command{
						{
							Name:        "ls",
							Aliases:     []string{"list"},
							Usage:       "List organization members and their permission level",
							Description: "List organization members and their permission level.",
							UsageText:   "earthly org [--org organization] members ls",
							Action:      a.actionMemberList,
						},
						{
							Name:        "update",
							Usage:       "Update an organization member's permission",
							Description: "Update an organization member's permission.",
							UsageText:   "earthly org [--org organization] members update --permission <permission> <user-email>",
							Action:      a.actionMemberUpdate,
							Flags: []cli.Flag{
								&cli.StringFlag{
									Name: "permission",
									Usage: `Update an organization member's permission.
					Can be one of: read, write, or admin.`,
									Destination: &a.userPermission,
								},
							},
						},
						{
							Name:        "rm",
							Usage:       "Remove a user from the organization",
							Description: "Remove a user from the organization.",
							UsageText:   "earthly org [--org organization] members rm <user-email>",
							Action:      a.actionMemberRemove,
						},
					},
				},
				{
					Name:        "select",
					Usage:       "Select the default organization",
					UsageText:   "earthly [options] org select <org-name>",
					Description: "Select the default organization.",
					Aliases:     []string{"s"},
					Action:      a.actionSelect,
				},
				{
					Name:        "unselect",
					Usage:       "Unselects the default organization",
					UsageText:   "earthly [options] org select <org-name>",
					Description: "Unselects the default organization.",
					Aliases:     []string{"uns"},
					Action:      a.actionUnselect,
				},
			},
		},
	}
}

func (a *Org) actionCreate(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgCreate")
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	org := cliCtx.Args().Get(0)
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	err = cloudClient.CreateOrg(cliCtx.Context, org)
	if err != nil {
		return errors.Wrap(err, "failed to create org")
	}
	return nil
}

func (a *Org) actionList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgList")
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	orgs, err := cloudClient.ListOrgs(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to list orgs")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, org := range orgs {
		selected := " "
		if org.Name == a.cli.Cfg().Global.Org {
			selected = "*"
		}
		fmt.Fprintf(w, "%s\t%s", selected, org.Name)
		var orgPermission string
		if org.Admin {
			orgPermission = "\tadmin"
		} else {
			orgPermission = "\tmember"
		}
		if org.Personal {
			orgPermission += " (personal)"
		}
		fmt.Fprintf(w, "%s\n", orgPermission)
	}
	w.Flush()

	return nil
}

func (a *Org) actionListPermissions(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgListPermissions")
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	cloudClient, err := helper.NewCloudClient(a.cli)
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

func (a *Org) actionInviteAccept(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgInviteAccept")

	if cliCtx.NArg() != 1 {
		return errors.New("invite code is required")
	}

	code := cliCtx.Args().Get(0)
	if code == "" {
		return errors.New("invite code is required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	err = cloudClient.AcceptInvite(cliCtx.Context, code)
	if err != nil {
		return errors.Wrap(err, "failed to accept invite")
	}

	a.cli.Console().Printf("Invite accepted!")

	return nil
}

func (a *Org) actionInviteList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgInviteList")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
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

func (a *Org) actionInviteEmail(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgInviteEmail")
	if cliCtx.NArg() == 0 {
		return errors.New("user email address required")
	}
	emails := cliCtx.Args().Slice()

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}
	permission := a.invitePermission
	if permission == "" {
		permission, err = common.PromptInput(cliCtx.Context, "New user's permission [read/write/admin] (default=read): ")
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
			Message:    a.inviteMessage,
		}
		_, err = cloudClient.InviteToOrg(cliCtx.Context, invite)
		if err != nil {
			return errors.Wrapf(err, "failed to invite user %s into org", userEmail)
		}
		a.cli.Console().Printf("Invite sent to %s", userEmail)
	}

	return nil
}

func (a *Org) actionRevoke(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgRevoke")
	if cliCtx.NArg() < 2 {
		return errors.New("invalid number of arguments provided")
	}
	path := cliCtx.Args().Get(0)
	if !strings.HasSuffix(path, "/") {
		return errors.New("revoked paths must end with a slash (/)")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
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

func (a *Org) actionMemberList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgMemberList")

	if cliCtx.NArg() != 0 {
		return errors.New("expected no arguments")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	members, err := cloudClient.ListOrgMembers(cliCtx.Context, orgName)
	if err != nil {
		return err
	}

	if len(members) == 0 {
		a.cli.Console().Printf("No members in %s", orgName)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, member := range members {
		fmt.Fprintf(w, "%s\t%s\n", member.UserEmail, member.Permission)
	}
	w.Flush()

	return nil
}

func (a *Org) actionMemberUpdate(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgMemberUpdate")

	if cliCtx.NArg() < 1 {
		return errors.New("member email required")
	}

	if cliCtx.NArg() > 1 {
		return errors.New("too many arguments provided")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("member email required")
	}

	if a.userPermission == "" {
		return errors.New("permission required")
	}

	err = cloudClient.UpdateOrgMember(cliCtx.Context, orgName, userEmail, a.userPermission)
	if err != nil {
		return err
	}

	a.cli.Console().Printf("Member %q updated successfully", userEmail)

	return nil
}

func (a *Org) actionMemberRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgMemberRemove")

	if cliCtx.NArg() != 1 {
		return errors.New("member email required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
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

	a.cli.Console().Printf("Member %q removed successfully", userEmail)

	return nil
}

func (a *Org) actionSelect(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgSelect")
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	org := cliCtx.Args().Get(0)
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	_, err = cloudClient.GetOrgID(cliCtx.Context, org)
	if err != nil {
		return errors.Wrap(err, "failed to get org")
	}

	inConfig, err := config.ReadConfigFile(a.cli.Flags().ConfigPath)
	if err != nil {
		if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "read config")
		}
	}

	newConfig, err := config.Upsert(inConfig, "global.org", org)
	if err != nil {
		return errors.Wrap(err, "could not update default org")
	}
	a.cli.Cfg().Global.Org = org

	if a.cli.Cfg().Satellite.Org != "" {
		newConfig, err = config.Upsert(newConfig, "satellite.org", "")
		if err != nil {
			return errors.Wrap(err, "could not remove deprecated setting")
		}
		a.cli.Cfg().Satellite.Org = ""
	}

	err = config.WriteConfigFile(a.cli.Flags().ConfigPath, newConfig)
	if err != nil {
		return errors.Wrap(err, "could not save config")
	}
	a.cli.Console().Printf("Updated selected org in %s", a.cli.Flags().ConfigPath)

	return nil
}

func (a *Org) actionUnselect(cliCtx *cli.Context) error {
	a.cli.SetCommandName("orgSelect")
	if cliCtx.NArg() != 0 {
		return errors.New("invalid number of arguments provided")
	}

	inConfig, err := config.ReadConfigFile(a.cli.Flags().ConfigPath)
	if err != nil {
		if cliCtx.IsSet("config") || !errors.Is(err, os.ErrNotExist) {
			return errors.Wrap(err, "read config")
		}
	}

	newConfig, err := config.Delete(inConfig, "global.org")
	if err != nil {
		return errors.Wrap(err, "could not unselect default org")
	}

	err = config.WriteConfigFile(a.cli.Flags().ConfigPath, newConfig)
	if err != nil {
		return errors.Wrap(err, "could not save config")
	}
	a.cli.Console().Printf("Unselected org in %s", a.cli.Flags().ConfigPath)

	return nil
}
