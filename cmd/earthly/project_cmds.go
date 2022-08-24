package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/earthly/earthly/cloud"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const dateFormat = "2006-01-02"

func (app *earthlyApp) projectCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "ls",
			Usage:       "List all projects that belong to the specified organization",
			Description: "List all projects that belong to the specified organization",
			UsageText:   "earthly project [--org <organization-name>] ls",
			Action:      app.actionProjectList,
		},
		{
			Name:        "rm",
			Usage:       "Remove an existing project from the organization",
			Description: "Remove an existing project from the organization",
			UsageText:   "earthly project [--org <organization-name>] rm <project-name>",
			Action:      app.actionProjectRemove,
		},
		{
			Name:        "create",
			Usage:       "Create a new project in the specified organization",
			Description: "Create a new project in the specified organization",
			UsageText:   "earthly project [--org <organization-name>] create <project-name>",
			Action:      app.actionProjectCreate,
		},
		{
			Name:        "member",
			Aliases:     []string{"members"},
			Usage:       "Manage project members",
			Description: "Manage project members",
			UsageText:   "earthly project member (ls|rm|add|update)",
			Subcommands: []*cli.Command{
				{
					Name:        "add",
					Usage:       "Add a new member to the specified project",
					Description: "Add a new member to the specified project",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name> member add <user-email> <permission>",
					Action:      app.actionProjectMemberAdd,
				},
				{
					Name:        "rm",
					Usage:       "Remove a member from the specified project",
					Description: "Remove a member from the specified project",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name member rm <user-email>",
					Action:      app.actionProjectMemberRemove,
				},
				{
					Name:        "ls",
					Usage:       "List all members in the specified project",
					Description: "List all members in the specified project",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name> member ls",
					Action:      app.actionProjectMemberList,
				},
				{
					Name:        "update",
					Usage:       "Update the project member's permission",
					Description: "Update the project member's permission",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name> member update <user-email> <permission>",
					Action:      app.actionProjectMemberUpdate,
				},
			},
		},
	}
}

func (app *earthlyApp) actionProjectList(cliCtx *cli.Context) error {
	app.commandName = "projectList"

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	projects, err := cloudClient.ListProjects(cliCtx.Context, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to list projects")
	}

	for _, project := range projects {
		app.console.Printf("%s\n", project.Name)
	}

	return nil
}

func (app *earthlyApp) actionProjectRemove(cliCtx *cli.Context) error {
	app.commandName = "projectRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("project name is required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	projectName := cliCtx.Args().Get(0)
	if projectName == "" {
		return errors.New("project name is required")
	}

	err = cloudClient.DeleteProject(cliCtx.Context, orgName, projectName)
	if err != nil {
		return errors.Wrap(err, "failed to remove project")
	}

	return nil
}

func (app *earthlyApp) actionProjectCreate(cliCtx *cli.Context) error {
	app.commandName = "projectCreate"

	if cliCtx.NArg() != 1 {
		return errors.New("project name is required")
	}

	if app.projectName == "" {
		return errors.New("project name is required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	_, err = cloudClient.CreateProject(cliCtx.Context, app.projectName, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}

	return nil
}

func (app *earthlyApp) actionProjectMemberList(cliCtx *cli.Context) error {
	app.commandName = "projectMemberList"

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	if app.projectName == "" {
		return errors.New("project name is required")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	members, err := cloudClient.ListProjectMembers(cliCtx.Context, orgName, app.projectName)
	if err != nil {
		return errors.Wrap(err, "failed to list project members")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "User Email\tPermission\tCreated\n")
	for _, m := range members {
		fmt.Fprintf(w, "%s\t%s\t%s\n", m.UserEmail, m.Permission, m.CreatedAt.Format(dateFormat))
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionProjectMemberRemove(cliCtx *cli.Context) error {
	app.commandName = "projectMemberRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("user email are required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if app.projectName == "" {
		return errors.New("project name is required")
	}

	userEmail := cliCtx.Args().Get(0)
	if userEmail == "" {
		return errors.New("user email is required")
	}

	err = cloudClient.RemoveProjectMember(cliCtx.Context, orgName, app.projectName, userEmail)
	if err != nil {
		return errors.Wrap(err, "failed to remove project member")
	}

	return nil
}

func (app *earthlyApp) actionProjectMemberAdd(cliCtx *cli.Context) error {
	app.commandName = "projectMemberAdd"

	if cliCtx.NArg() != 2 {
		return errors.New("user email and permission arguments are required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if app.projectName == "" {
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

	err = cloudClient.AddProjectMember(cliCtx.Context, orgName, app.projectName, userEmail, permission)
	if err != nil {
		return errors.Wrap(err, "failed to add project member")
	}

	return nil
}

func (app *earthlyApp) actionProjectMemberUpdate(cliCtx *cli.Context) error {
	app.commandName = "projectMemberUpdate"

	if cliCtx.NArg() != 2 {
		return errors.New("user email and permission arguments are required")
	}

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	if app.projectName == "" {
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

	err = cloudClient.UpdateProjectMember(cliCtx.Context, orgName, app.projectName, userEmail, permission)
	if err != nil {
		return errors.Wrap(err, "failed to update project member")
	}

	return nil
}

// projectOrgName returns the specified org or retrieves the default org from the API.
func (app *earthlyApp) projectOrgName(ctx context.Context, cloudClient cloud.Client) (string, error) {

	if app.orgName != "" {
		return app.orgName, nil
	}

	userOrgs, err := cloudClient.ListOrgs(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to list organizations")
	}

	if len(userOrgs) == 0 {
		return "", errors.New("no organizations found, please specify with --org")
	} else if len(userOrgs) > 1 {
		return "", errors.New("multiple organizations found, please specify with --org")
	}

	return userOrgs[0].Name, nil
}
