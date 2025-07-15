package base

import (
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cmd/earthly/flag"
	"github.com/earthly/earthly/logbus"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
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
	deferredFuncs           []func()
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

func (c *CLI) AddDeferredFunc(f func()) {
	c.deferredFuncs = append([]func(){f}, c.deferredFuncs...)
}

func (c *CLI) ExecuteDeferredFuncs() {
	for _, f := range c.deferredFuncs {
		f()
	}
}
