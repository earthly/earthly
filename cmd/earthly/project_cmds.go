package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
)

const dateFormat = "2006-01-02"

func (app *earthlyApp) projectCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "ls",
			Usage:       "List all projects that belong to the specified organization *beta*",
			Description: "List all projects that belong to the specified organization *beta*",
			UsageText:   "earthly project [--org <organization-name>] ls",
			Action:      app.actionProjectList,
		},
		{
			Name:        "rm",
			Usage:       "Remove an existing project and all of its associated pipelines and secrets *beta*",
			Description: "Remove an existing project and all of its associated pipelines and secrets *beta*",
			UsageText:   "earthly project [--org <organization-name>] rm <project-name>",
			Action:      app.actionProjectRemove,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "force",
					Aliases:     []string{"f"},
					Usage:       "Force removal without asking permission",
					Destination: &app.forceRemoveProject,
				},
			},
		},
		{
			Name:        "create",
			Usage:       "Create a new project in the specified organization *beta*",
			Description: "Create a new project in the specified organization *beta*",
			UsageText:   "earthly project [--org <organization-name>] create <project-name>",
			Action:      app.actionProjectCreate,
		},
		{
			Name:        "member",
			Aliases:     []string{"members"},
			Usage:       "Manage project members *beta*",
			Description: "Manage project members *beta*",
			UsageText:   "earthly project member (ls|rm|add|update)",
			Subcommands: []*cli.Command{
				{
					Name:        "add",
					Usage:       "Add a new member to the specified project *beta*",
					Description: "Add a new member to the specified project *beta*",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name> member add <user-email> <permission>",
					Action:      app.actionProjectMemberAdd,
				},
				{
					Name:        "rm",
					Usage:       "Remove a member from the specified project *beta*",
					Description: "Remove a member from the specified project *beta*",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name member rm <user-email>",
					Action:      app.actionProjectMemberRemove,
				},
				{
					Name:        "ls",
					Usage:       "List all members in the specified project *beta*",
					Description: "List all members in the specified project *beta*",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name> member ls",
					Action:      app.actionProjectMemberList,
				},
				{
					Name:        "update",
					Usage:       "Update the project member's permission *beta*",
					Description: "Update the project member's permission *beta*",
					UsageText:   "earthly project [--org <organization-name>] --project <project-name> member update <user-email> <permission>",
					Action:      app.actionProjectMemberUpdate,
				},
			},
		},
	}
}

func (app *earthlyApp) actionProjectList(cliCtx *cli.Context) error {
	app.commandName = "projectList"

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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
		fmt.Println(project.Name)
	}

	return nil
}

func (app *earthlyApp) actionProjectRemove(cliCtx *cli.Context) error {
	app.commandName = "projectRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("project name is required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	projectName := cliCtx.Args().Get(0)
	if projectName == "" {
		return errors.New("project name is required")
	}

	if !app.forceRemoveProject {
		answer, err := promptInput(cliCtx.Context,
			"WARNING: you are about to permanently delete this project and all of its associated pipelines, build history and secrets.\n"+
				"Would you like to continue?\n"+
				"Type 'y' or 'yes': ")
		if err != nil {
			return errors.Wrap(err, "failed requesting user input")
		}
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			app.console.Printf("Operation aborted.")
			return nil
		}
	}

	err = cloudClient.DeleteProject(cliCtx.Context, orgName, projectName)
	if err != nil {
		return errors.Wrap(err, "failed to remove project")
	}

	app.console.Printf("Project %s removed from %s", projectName, orgName)

	return nil
}

func (app *earthlyApp) actionProjectCreate(cliCtx *cli.Context) error {
	app.commandName = "projectCreate"

	if cliCtx.NArg() != 1 {
		return errors.New("project name is required")
	}

	projectName := cliCtx.Args().Get(0)
	if projectName == "" {
		return errors.New("project name is required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	orgName, err := app.projectOrgName(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	_, err = cloudClient.CreateProject(cliCtx.Context, projectName, orgName)
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}

	app.console.Printf("Project %s created in %s", projectName, orgName)

	return nil
}

func (app *earthlyApp) actionProjectMemberList(cliCtx *cli.Context) error {
	app.commandName = "projectMemberList"

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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

	if len(members) == 0 {
		app.console.Printf("No permissions found for this secret")
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

func (app *earthlyApp) actionProjectMemberRemove(cliCtx *cli.Context) error {
	app.commandName = "projectMemberRemove"

	if cliCtx.NArg() != 1 {
		return errors.New("user email are required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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

	app.console.Printf("%s was removed from %s", userEmail, orgName)

	return nil
}

func (app *earthlyApp) actionProjectMemberAdd(cliCtx *cli.Context) error {
	app.commandName = "projectMemberAdd"

	if cliCtx.NArg() != 2 {
		return errors.New("user email and permission arguments are required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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

	app.console.Printf("%s has been added to %s with %s permission", userEmail, orgName, permission)

	return nil
}

func (app *earthlyApp) actionProjectMemberUpdate(cliCtx *cli.Context) error {
	app.commandName = "projectMemberUpdate"

	if cliCtx.NArg() != 2 {
		return errors.New("user email and permission arguments are required")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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

	app.console.Printf("%s now has %s permission in %s", userEmail, permission, orgName)

	return nil
}

// projectOrgName returns the specified org or retrieves the default org from the API.
func (app *earthlyApp) projectOgirgName(ctx context.Context, cloudClient *cloud.Client) (string, error) {

	if app.orgName != "" {
		return app.orgName, nil
	} else if app.cfg.Global.Org != "" {
		return app.cfg.Global.Org, nil
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
