package subcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/helper"
)

const dateFormat = "2006-01-02"

type Project struct {
	cli CLI

	forceRemoveProject bool
}

func NewProject(cli CLI) *Project {
	return &Project{
		cli: cli,
	}
}

func (a *Project) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "project",
			Aliases: []string{"projects"},
			Description: `Manage Earthly projects which are shared resources of Earthly orgs.

Within Earthly projects users can be invited and granted different access levels including: read, read+secrets, write, and admin.`,
			Usage:     "Manage Earthly projects",
			UsageText: "earthly project (ls|rm|create|member)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The name of the Earthly organization to which the Earthly project belongs",
					Required:    false,
					Destination: &a.cli.Flags().OrgName,
				},
				&cli.StringFlag{
					Name:        "project",
					Aliases:     []string{"p"},
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The Earthly project to act on",
					Required:    false,
					Destination: &a.cli.Flags().ProjectName,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "ls",
					Usage:       "List all projects that belong to the specified organization",
					Description: "List all projects that belong to the specified organization.",
					UsageText:   "earthly project [--org <organization-name>] ls",
					Action:      a.actionList,
				},
				{
					Name:        "rm",
					Usage:       "Remove an existing project and all of its associated pipelines and secrets",
					Description: "Remove an existing project and all of its associated pipelines and secrets.",
					UsageText:   "earthly project [--org <organization-name>] rm <project-name>",
					Action:      a.actionRemove,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "force",
							Aliases:     []string{"f"},
							Usage:       "Force removal without asking permission",
							Destination: &a.forceRemoveProject,
						},
					},
				},
				{
					Name:        "create",
					Usage:       "Create a new project in the specified organization",
					Description: "Create a new project in the specified organization.",
					UsageText:   "earthly project [--org <organization-name>] create <project-name>",
					Action:      a.actionCreate,
				},
				{
					Name:        "member",
					Aliases:     []string{"members"},
					Usage:       "Manage project members",
					Description: "Manage project members.",
					UsageText:   "earthly project member (ls|rm|add|update)",
					Subcommands: []*cli.Command{
						{
							Name:        "add",
							Usage:       "Add a new member to the specified project",
							Description: "Add a new member to the specified project.",
							UsageText:   "earthly project [--org <organization-name>] --project <project-name> member add <user-email> <permission>",
							Action:      a.actionMemberAdd,
						},
						{
							Name:        "rm",
							Usage:       "Remove a member from the specified project",
							Description: "Remove a member from the specified project.",
							UsageText:   "earthly project [--org <organization-name>] --project <project-name member rm <user-email>",
							Action:      a.actionMemberRemove,
						},
						{
							Name:        "ls",
							Usage:       "List all members in the specified project",
							Description: "List all members in the specified project.",
							UsageText:   "earthly project [--org <organization-name>] --project <project-name> member ls",
							Action:      a.actionMemberList,
						},
						{
							Name:        "update",
							Usage:       "Update the project member's permission",
							Description: "Update the project member's permission.",
							UsageText:   "earthly project [--org <organization-name>] --project <project-name> member update <user-email> <permission>",
							Action:      a.actionMemberUpdate,
						},
					},
				},
			},
		},
	}
}

func (a *Project) actionList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectList")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	projects, err := cloudClient.ListProjects(cliCtx.Context, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to list projects")
	}

	for _, project := range projects {
		fmt.Println(project.Name)
	}

	return nil
}

func (a *Project) actionRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectRemove")

	if cliCtx.NArg() != 1 {
		return errors.New("project name is required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	projectName := cliCtx.Args().Get(0)
	if projectName == "" {
		return errors.New("project name is required")
	}

	if !a.forceRemoveProject {
		answer, err := common.PromptInput(cliCtx.Context,
			"WARNING: you are about to permanently delete this project and all of its associated pipelines, build history and secrets.\n"+
				"Would you like to continue?\n"+
				"Type 'y' or 'yes': ")
		if err != nil {
			return errors.Wrap(err, "failed requesting user input")
		}
		if !isResponseYes(answer) {
			a.cli.Console().Printf("Operation aborted.")
			return nil
		}
	}

	err = cloudClient.DeleteProject(cliCtx.Context, orgName, projectName)
	if err != nil {
		return errors.Wrap(err, "failed to remove project")
	}

	a.cli.Console().Printf("Project %s removed from %s", projectName, orgName)

	return nil
}

func (a *Project) actionCreate(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectCreate")

	if cliCtx.NArg() != 1 {
		return errors.New("project name is required")
	}

	projectName := cliCtx.Args().Get(0)
	if projectName == "" {
		return errors.New("project name is required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	_, err = cloudClient.CreateProject(cliCtx.Context, projectName, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}

	a.cli.Console().Printf("Project %s created in %s", projectName, orgName)

	return nil
}

func (a *Project) actionMemberList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectMemberList")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	if a.cli.Flags().ProjectName == "" {
		return errors.New("project name is required")
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	members, err := cloudClient.ListProjectMembers(cliCtx.Context, orgName, a.cli.Flags().ProjectName)
	if err != nil {
		return errors.Wrap(err, "failed to list project members")
	}

	if len(members) == 0 {
		a.cli.Console().Printf("No permissions found for this secret")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "User Email\tPermission\tCreated\n")
	for _, m := range members {
		fmt.Fprintf(w, "%s\t%s\t%s\n", m.UserEmail, m.Permission, m.CreatedAt.Format(dateFormat))
	}
	w.Flush()

	return nil
}

func (a *Project) actionMemberRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectMemberRemove")

	if cliCtx.NArg() != 1 {
		return errors.New("user email are required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if a.cli.Flags().ProjectName == "" {
		return errors.New("project name is required")
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	err = cloudClient.RemoveProjectMember(cliCtx.Context, orgName, a.cli.Flags().ProjectName, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to remove project member")
	}

	a.cli.Console().Printf("%s was removed from %s", userEmail, orgName)

	return nil
}

func (a *Project) actionMemberAdd(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectMemberAdd")

	if cliCtx.NArg() != 2 {
		return errors.New("user email and permission arguments are required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if a.cli.Flags().ProjectName == "" {
		return errors.New("project name is required")
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	permission := cliCtx.Args().Get(1)
	if permission == "" {
		return errors.New("permission is required")
	}

	err = cloudClient.AddProjectMember(cliCtx.Context, orgName, a.cli.Flags().ProjectName, userEmail, permission)
	if err != nil {
		return errors.Wrap(err, "failed to add project member")
	}

	a.cli.Console().Printf("%s has been added to %s with %s permission", userEmail, orgName, permission)

	return nil
}

func (a *Project) actionMemberUpdate(cliCtx *cli.Context) error {
	a.cli.SetCommandName("projectMemberUpdate")

	if cliCtx.NArg() != 2 {
		return errors.New("user email and permission arguments are required")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	orgName, err := projectOrgName(a.cli, cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if a.cli.Flags().ProjectName == "" {
		return errors.New("project name is required")
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	permission := cliCtx.Args().Get(1)
	if permission == "" {
		return errors.New("permission is required")
	}

	err = cloudClient.UpdateProjectMember(cliCtx.Context, orgName, a.cli.Flags().ProjectName, userEmail, permission)
	if err != nil {
		return errors.Wrap(err, "failed to update project member")
	}

	a.cli.Console().Printf("%s now has %s permission in %s", userEmail, permission, orgName)

	return nil
}
