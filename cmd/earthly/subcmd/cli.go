package subcmd

import (
	"context"

	"github.com/moby/buildkit/client"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/logbus/setup"
)

type CLI interface {
	App() *cli.App

	Version() string
	GitSHA() string

	Flags() *flag.Global
	Console() conslogging.ConsoleLogger
	SetConsole(conslogging.ConsoleLogger)

	InitFrontend(*cli.Context) error
	Cfg() *config.Config
	SetCommandName(name string)

	SetAnaMetaTarget(domain.Target)
	SetAnaMetaIsRemoteBK(bool)
	SetAnaMetaBKPlatform(string)
	SetAnaMetaUserPlatform(string)
	IsUsingSatellite(*cli.Context) bool
	OrgName() string

	GetBuildkitClient(*cli.Context, *cloud.Client) (client *client.Client, cleanupTLS func(), err error)
	GetSatelliteOrg(context.Context, *cloud.Client) (string, string, error)

	ConfigureSatellite(*cli.Context, *cloud.Client, string, string) (cleanupTLS func(), err error)
	CIHost() string
	LogbusSetup() *setup.BusSetup
	Logbus() *logbus.Bus

	AddDeferredFunc(f func())

	CollectBillingInfo(ctx context.Context, cloudClient *cloud.Client, orgName string) error
}
