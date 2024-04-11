package subcmd

import (
	"github.com/urfave/cli/v2"
)

type Root struct {
	cli CLI

	buildCmd *Build
}

func NewRoot(cli CLI, buildCmd *Build) *Root {
	return &Root{
		cli:      cli,
		buildCmd: buildCmd,
	}
}

func (a *Root) Cmds() []*cli.Command {
	cmds := concatCmds([][]*cli.Command{
		NewDebug(a.cli).Cmds(),
		NewBootstrap(a.cli).Cmds(),
		a.buildCmd.Cmds(),
		NewAccount(a.cli).Cmds(),
		NewConfig(a.cli).Cmds(),
		NewDoc(a.cli).Cmds(),
		NewDoc2Earth(a.cli).Cmds(),
		NewInit(a.cli).Cmds(),
		NewList(a.cli).Cmds(),
		NewOrg(a.cli).Cmds(),
		NewProject(a.cli).Cmds(),
		NewPrune(a.cli).Cmds(),
		NewAutoSkip(a.cli).Cmds(),
		NewRegistry(a.cli).Cmds(),
		NewSatellite(a.cli).Cmds(),
		NewCloudInstallation(a.cli).Cmds(),
		NewSecret(a.cli).Cmds(),
		NewWeb(a.cli).Cmds(),
		NewBilling(a.cli).Cmds(),
	})

	return cmds

}
