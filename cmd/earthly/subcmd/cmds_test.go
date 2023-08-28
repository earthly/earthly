package subcmd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/earthly/earthly/conslogging"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
	"github.com/poy/onpar/matchers"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cmd/earthly/app"
	"github.com/earthly/earthly/cmd/earthly/base"
	"github.com/earthly/earthly/cmd/earthly/subcmd"
)

func TestRootCmdsHelp(t *testing.T) {
	type testCtx struct {
		t      *testing.T
		expect expect.Expectation
	}

	o := onpar.BeforeEach(onpar.New(t), func(t *testing.T) testCtx {
		return testCtx{
			t:      t,
			expect: expect.New(t),
		}
	})
	defer o.Run()

	ctx := context.TODO()
	newCLI := base.NewCLI(conslogging.ConsoleLogger{},
		base.WithVersion(""),
		base.WithGitSHA(""),
		base.WithBuiltBy(""),
		base.WithDefaultBuildkitdImage(""),
		base.WithDefaultInstallationName(""),
	)
	buildApp := subcmd.NewBuild(newCLI)
	rootApp := subcmd.NewRoot(newCLI, buildApp)
	app := app.NewEarthlyApp(newCLI, rootApp, buildApp, ctx)

	rootCLI := app.BaseCLI.App().Commands

	usageChecks := onpar.TableSpec(o, func(tt testCtx, cmd *cli.Command) {
		tt.expect(cmd.Usage).To(matchers.Not(matchers.EndWith(".")))
	})
	descChecks := onpar.TableSpec(o, func(tt testCtx, cmd *cli.Command) {
		tt.expect(cmd.Description).To(matchers.EndWith("."))
	})
	for _, subCmd := range checkSubCommands(rootCLI) {
		usageChecks.Entry(fmt.Sprintf("Help usage for %s should not end with '.'", subCmd.Name), subCmd)
		descChecks.Entry(fmt.Sprintf("Help description for %s should end with '.'", subCmd.Name), subCmd)
	}
}

// Check if command has any subCommands to verify
func checkSubCommands(commands []*cli.Command) []*cli.Command {
	var allCommands []*cli.Command
	for _, command := range commands {
		allCommands = append(allCommands, command)
		if len(command.Subcommands) != 0 {
			allCommands = append(allCommands, checkSubCommands(command.Subcommands)...)
		}
	}
	return allCommands
}
