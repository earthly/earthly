package base

import (
	"strings"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/logbus"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/logbus/setup"
)

type CLI struct {
	app                     *cli.App
	console                 conslogging.ConsoleLogger
	cfg                     *config.Config
	logbusSetup             *setup.BusSetup
	logbus                  *logbus.Bus
	commandName             string
	version                 string
	gitSHA                  string
	builtBy                 string
	defaultBuildkitdImage   string
	defaultInstallationName string
	flags                   flag.Global
	analyticsMetadata
}

type analyticsMetadata struct {
	isSatellite             bool
	isRemoteBuildkit        bool
	satelliteCurrentVersion string
	buildkitPlatform        string
	userPlatform            string
	target                  domain.Target
}

type CLIOpt func(CLI) CLI

func WithVersion(version string) CLIOpt {
	return func(c CLI) CLI {
		c.version = version
		return c
	}
}

func WithGitSHA(sha string) CLIOpt {
	return func(c CLI) CLI {
		c.gitSHA = sha
		return c
	}
}

func WithBuiltBy(builtby string) CLIOpt {
	return func(c CLI) CLI {
		c.builtBy = builtby
		return c
	}
}

func WithDefaultBuildkitdImage(image string) CLIOpt {
	return func(c CLI) CLI {
		c.defaultBuildkitdImage = image
		return c
	}
}

func WithDefaultInstallationName(name string) CLIOpt {
	return func(c CLI) CLI {
		c.defaultInstallationName = name
		return c
	}
}

func NewCLI(console conslogging.ConsoleLogger, opts ...CLIOpt) *CLI {
	cli := CLI{
		app:     cli.NewApp(),
		console: console,
		logbus:  logbus.New(),
		flags: flag.Global{
			BuildkitdSettings: buildkitd.Settings{},
		},
	}

	for _, opt := range opts {
		cli = opt(cli)
	}

	return &cli
}

func (c *CLI) App() *cli.App {
	return c.app
}

func (c *CLI) SetAppUsage(usage string) {
	c.app.Usage = usage
}
func (c *CLI) SetAppUsageText(usageText string) {
	c.app.UsageText = usageText
}
func (c *CLI) SetAppUseShortOptionHandling(use bool) {
	c.app.UseShortOptionHandling = use
}
func (c *CLI) SetAction(action cli.ActionFunc) {
	c.app.Action = action
}
func (c *CLI) SetVersion(version string) {
	c.app.Version = version
}
func (c *CLI) SetFlags(flags []cli.Flag) {
	c.app.Flags = flags
}
func (c *CLI) SetCommands(commands []*cli.Command) {
	c.app.Commands = commands
}
func (c *CLI) SetBefore(before cli.BeforeFunc) {
	c.app.Before = before
}

func (c *CLI) Console() conslogging.ConsoleLogger {
	return c.console
}
func (c *CLI) SetConsole(cons conslogging.ConsoleLogger) {
	c.console = cons
}

func (c *CLI) Cfg() *config.Config {
	return c.cfg
}
func (c *CLI) SetCfg(cfg *config.Config) {
	c.cfg = cfg
}

func (c *CLI) CommandName() string {
	return c.commandName
}
func (c *CLI) SetCommandName(commandName string) {
	c.commandName = commandName
}

func (c *CLI) Version() string {
	return c.version
}
func (c *CLI) GitSHA() string {
	return c.gitSHA
}
func (c *CLI) BuiltBy() string {
	return c.builtBy
}
func (c *CLI) DefaultBuildkitdImage() string {
	return c.defaultBuildkitdImage
}
func (c *CLI) DefaultInstallationName() string {
	return c.defaultInstallationName
}

func (c *CLI) LogbusSetup() *setup.BusSetup {
	return c.logbusSetup
}
func (c *CLI) SetLogbusSetup(setup *setup.BusSetup) {
	c.logbusSetup = setup
}

func (c *CLI) Logbus() *logbus.Bus {
	return c.logbus
}
func (c *CLI) SetLogbus(logbus *logbus.Bus) {
	c.logbus = logbus
}

func (c *CLI) Flags() *flag.Global {
	return &c.flags
}

func (c *CLI) AnaMetaIsSat() bool {
	return c.analyticsMetadata.isSatellite
}
func (c *CLI) AnaMetaIsRemoteBK() bool {
	return c.analyticsMetadata.isRemoteBuildkit
}
func (c *CLI) AnaMetaSatCurrentVersion() string {
	return c.analyticsMetadata.satelliteCurrentVersion
}
func (c *CLI) AnaMetaBKPlatform() string {
	return c.analyticsMetadata.buildkitPlatform
}
func (c *CLI) AnaMetaUserPlatform() string {
	return c.analyticsMetadata.userPlatform
}
func (c *CLI) AnaMetaTarget() domain.Target {
	return c.analyticsMetadata.target
}
func (c *CLI) SetAnaMetaIsSat(isSat bool) {
	c.analyticsMetadata.isSatellite = isSat
}
func (c *CLI) SetAnaMetaIsRemoteBK(isRBK bool) {
	c.analyticsMetadata.isRemoteBuildkit = isRBK
}
func (c *CLI) SetAnaMetaSatCurrentVersion(currentVersion string) {
	c.analyticsMetadata.satelliteCurrentVersion = currentVersion
}
func (c *CLI) SetAnaMetaBKPlatform(platform string) {
	c.analyticsMetadata.buildkitPlatform = platform
}
func (c *CLI) SetAnaMetaUserPlatform(platform string) {
	c.analyticsMetadata.userPlatform = platform
}
func (c *CLI) SetAnaMetaTarget(target domain.Target) {
	c.analyticsMetadata.target = target
}

// CIHost returns protocol://hostname
func (c *CLI) CIHost() string {
	switch {
	case strings.Contains(c.Flags().CloudGRPCAddr, "staging"):
		return "https://cloud.staging.earthly.dev"
	case strings.Contains(c.Flags().CloudGRPCAddr, "earthly.local"):
		return "http://earthly.local:3000"
	}
	return "https://cloud.earthly.dev"
}
