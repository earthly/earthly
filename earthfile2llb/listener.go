package earthfile2llb

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/earthly/earthly/earthfile2llb/parser"
	"github.com/pkg/errors"
)

var _ parser.EarthParserListener = &listener{}

type listener struct {
	*parser.BaseEarthParserListener
	interpreter commandInterpreter
	ctx         context.Context

	executeTarget   string
	currentTarget   string
	targetFound     bool
	saveImageExists bool
	pushOnlyAllowed bool

	envArgKey   string
	envArgValue string
	labelKeys   []string
	labelValues []string

	execMode  bool
	stmtWords []string

	err error
}

func newListener(ctx context.Context, interpreter commandInterpreter, executeTarget string) *listener {
	return &listener{
		ctx:           ctx,
		interpreter:   interpreter,
		executeTarget: executeTarget,
		currentTarget: "base",
		targetFound:   (executeTarget == "base"),
	}
}

func (l *listener) Err() error {
	if l.err != nil {
		return l.err
	}
	if !l.targetFound {
		return fmt.Errorf("target %s not defined", l.executeTarget)
	}
	return nil
}

func (l *listener) EnterTargetHeader(c *parser.TargetHeaderContext) {
	// Apply implicit SAVE IMAGE for +base.
	if l.executeTarget == "base" {
		if !l.saveImageExists {
			l.interpreter.SaveImage(l.ctx, []string{}, false)
		}
		l.saveImageExists = true
	}

	l.currentTarget = strings.TrimSuffix(c.GetText(), ":")
	if l.currentTarget == l.executeTarget {
		if l.targetFound {
			l.err = fmt.Errorf("target %s is declared twice", l.currentTarget)
			return
		}
		l.targetFound = true
	}
	if l.shouldSkip() {
		return
	}
	if l.currentTarget == "base" {
		l.err = errors.New("target name cannot be base")
		return
	}
	// Apply implicit FROM +base
	err := l.interpreter.From(l.ctx, "+base", nil)
	if err != nil {
		l.err = errors.Wrap(err, "apply implicit FROM +base")
		return
	}
	l.pushOnlyAllowed = false
}

//
// Commands.

func (l *listener) EnterStmt(c *parser.StmtContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = nil
	l.envArgKey = ""
	l.envArgValue = ""
	l.labelKeys = nil
	l.labelValues = nil
	l.execMode = false
}

func (l *listener) ExitFromStmt(c *parser.FromStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("FROM", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid FROM arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 1 {
		if fs.NArg() == 3 && fs.Arg(1) == "AS" {
			l.err = errors.New("AS not supported, use earthly targets instead")
		} else {
			l.err = fmt.Errorf("invalid number of arguments for FROM: %s", l.stmtWords)
		}
		return
	}
	imageName := fs.Arg(0)
	err = l.interpreter.From(l.ctx, imageName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrapf(err, "apply FROM %s", imageName)
		return
	}
}

func (l *listener) ExitCopyStmt(c *parser.CopyStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("COPY", flag.ContinueOnError)
	from := fs.String("from", "", "")
	isDirCopy := fs.Bool("dir", false, "")
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid COPY arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() < 2 {
		l.err = fmt.Errorf("not enough COPY arguments %v", l.stmtWords)
		return
	}
	if *from != "" {
		l.err = errors.New("COPY --from not implemented. Use COPY artifacts form instead")
		return
	}
	srcs := fs.Args()[:fs.NArg()-1]
	dest := fs.Arg(fs.NArg() - 1)
	allClassical := true
	allArtifacts := true
	for _, src := range srcs {
		if strings.Contains(src, "+") {
			allClassical = false
		} else {
			allArtifacts = false
		}
	}
	if !allClassical && !allArtifacts {
		l.err = fmt.Errorf("Combining artifacts and build context arguments in a single COPY command is not allowed: %v", srcs)
		return
	}
	if allArtifacts {
		for _, src := range srcs {
			err = l.interpreter.CopyArtifact(l.ctx, src, dest, buildArgs.Args, *isDirCopy)
			if err != nil {
				l.err = errors.Wrapf(err, "copy artifact")
				return
			}
		}
	} else {
		if len(buildArgs.Args) != 0 {
			l.err = fmt.Errorf("build args not supported for non +artifact arguments case %v", l.stmtWords)
			return
		}
		l.interpreter.CopyClassical(l.ctx, srcs, dest, *isDirCopy)
	}
}

func (l *listener) ExitRunStmt(c *parser.RunStmtContext) {
	if l.shouldSkip() {
		return
	}
	if len(l.stmtWords) < 1 {
		l.err = errors.New("not enough arguments for RUN")
		return
	}

	fs := flag.NewFlagSet("RUN", flag.ContinueOnError)
	pushFlag := fs.Bool("push", false, "")
	privileged := fs.Bool("privileged", false, "")
	withEntrypoint := fs.Bool("entrypoint", false, "")
	withDocker := fs.Bool("with-docker", false, "")
	secrets := new(StringSliceFlag)
	fs.Var(secrets, "secret", "")
	mounts := new(StringSliceFlag)
	fs.Var(mounts, "mount", "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid RUN arguments %v", l.stmtWords)
		return
	}
	withShell := !l.execMode
	if *withDocker {
		*privileged = true
	}
	if !*pushFlag && l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	// TODO: In the bracket case, should flags be outside of the brackets?

	err = l.interpreter.Run(l.ctx, fs.Args(), mounts.Args, secrets.Args, *privileged, *withEntrypoint, *withDocker, withShell, *pushFlag)
	if err != nil {
		l.err = errors.Wrap(err, "run")
		return
	}
	if *pushFlag {
		l.pushOnlyAllowed = true
	}
}

func (l *listener) ExitSaveArtifact(c *parser.SaveArtifactContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) == 0 {
		l.err = fmt.Errorf("no arguments provided to the SAVE ARTIFACT command")
		return
	}
	if len(l.stmtWords) > 5 {
		l.err = fmt.Errorf("too many arguments provided to the SAVE ARTIFACT command: %v", l.stmtWords)
		return
	}
	saveAsLocalTo := ""
	saveTo := "./"
	if len(l.stmtWords) >= 4 {
		if strings.Join(l.stmtWords[len(l.stmtWords)-3:len(l.stmtWords)-1], " ") == "AS LOCAL" {
			saveAsLocalTo = l.stmtWords[len(l.stmtWords)-1]
			if len(l.stmtWords) == 5 {
				saveTo = l.stmtWords[1]
			}
		} else {
			l.err = fmt.Errorf("invalid arguments for SAVE ARTIFACT command: %v", l.stmtWords)
			return
		}
	} else if len(l.stmtWords) == 2 {
		saveTo = l.stmtWords[1]
	} else if len(l.stmtWords) == 3 {
		l.err = fmt.Errorf("invalid arguments for SAVE ARTIFACT command: %v", l.stmtWords)
		return
	}
	saveFrom := l.stmtWords[0]

	err := l.interpreter.SaveArtifact(l.ctx, saveFrom, saveTo, saveAsLocalTo)
	if err != nil {
		l.err = errors.Wrap(err, "apply SAVE ARTIFACT")
		return
	}
}

func (l *listener) ExitSaveImage(c *parser.SaveImageContext) {
	if l.shouldSkip() {
		return
	}
	if l.saveImageExists {
		l.err = fmt.Errorf(
			"more than one SAVE IMAGE statement per target not allowed: %s", c.GetText())
		return
	}
	l.saveImageExists = true

	fs := flag.NewFlagSet("SAVE IMAGE", flag.ContinueOnError)
	pushFlag := fs.Bool("push", false, "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid SAVE IMAGE arguments %v", l.stmtWords)
		return
	}
	if !*pushFlag && l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if *pushFlag && fs.NArg() == 0 {
		l.err = errors.Wrapf(err, "invalid number of arguments for SAVE IMAGE --push: %v", l.stmtWords)
		return
	}
	imageNames := fs.Args()

	l.interpreter.SaveImage(l.ctx, imageNames, *pushFlag)
	if *pushFlag {
		l.pushOnlyAllowed = true
	}
}

func (l *listener) ExitBuildStmt(c *parser.BuildStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("BUILD", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid BUILD arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 1 {
		l.err = fmt.Errorf("invalid number of arguments for BUILD: %s", l.stmtWords)
		return
	}
	fullTargetName := fs.Arg(0)
	_, err = l.interpreter.Build(l.ctx, fullTargetName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrapf(err, "apply BUILD %s", fullTargetName)
		return
	}
}

func (l *listener) ExitWorkdirStmt(c *parser.WorkdirStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for WORKDIR: %v", l.stmtWords)
		return
	}
	workdirPath := l.stmtWords[0]
	l.interpreter.Workdir(l.ctx, workdirPath)
}

func (l *listener) ExitUserStmt(c *parser.UserStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for USER: %v", l.stmtWords)
		return
	}
	user := l.stmtWords[0]
	l.interpreter.User(l.ctx, user)
}

func (l *listener) ExitCmdStmt(c *parser.CmdStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	withShell := !l.execMode
	l.interpreter.Cmd(l.ctx, l.stmtWords, withShell)
}

func (l *listener) ExitEntrypointStmt(c *parser.EntrypointStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	withShell := !l.execMode
	l.interpreter.Entrypoint(l.ctx, l.stmtWords, withShell)
}

func (l *listener) ExitExposeStmt(c *parser.ExposeStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) == 0 {
		l.err = fmt.Errorf("no arguments provided to the EXPOSE command")
		return
	}
	l.interpreter.Expose(l.ctx, l.stmtWords)
}

func (l *listener) ExitVolumeStmt(c *parser.VolumeStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) == 0 {
		l.err = fmt.Errorf("no arguments provided to the VOLUME command")
		return
	}
	l.interpreter.Volume(l.ctx, l.stmtWords)
}

func (l *listener) ExitEnvStmt(c *parser.EnvStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	l.interpreter.Env(l.ctx, l.envArgKey, l.envArgValue)
}

func (l *listener) ExitArgStmt(c *parser.ArgStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	l.interpreter.Arg(l.ctx, l.envArgKey, l.envArgValue)
}

func (l *listener) ExitLabelStmt(c *parser.LabelStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.labelKeys) == 0 {
		l.err = fmt.Errorf("no labels provided in LABEL command: %s", c.GetText())
		return
	}
	if len(l.labelKeys) != len(l.labelValues) {
		l.err = fmt.Errorf("label keys and values do not match: %s", c.GetText())
		return
	}
	labels := make(map[string]string)
	for i := range l.labelKeys {
		labels[l.labelKeys[i]] = l.labelValues[i]
	}
	l.interpreter.Label(l.ctx, labels)
}

func (l *listener) ExitGitCloneStmt(c *parser.GitCloneStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("GIT CLONE", flag.ContinueOnError)
	branch := fs.String("branch", "", "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid GIT CLONE arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 2 {
		l.err = fmt.Errorf("invalid number of arguments for GIT CLONE: %s", l.stmtWords)
		return
	}
	gitURL := fs.Arg(0)
	gitCloneDest := fs.Arg(1)
	err = l.interpreter.GitClone(l.ctx, gitURL, *branch, gitCloneDest)
	if err != nil {
		l.err = errors.Wrap(err, "git clone")
		return
	}
}

func (l *listener) ExitDockerLoadStmt(c *parser.DockerLoadStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("DOCKER LOAD", flag.ContinueOnError)
	buildArgs := new(StringSliceFlag)
	fs.Var(buildArgs, "build-arg", "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid DOCKER LOAD arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() != 2 {
		l.err = fmt.Errorf("invalid number of arguments for DOCKER LOAD: %s", l.stmtWords)
		return
	}
	fullTargetName := fs.Arg(0)
	imageName := fs.Arg(1)
	err = l.interpreter.DockerLoad(l.ctx, fullTargetName, imageName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrap(err, "docker load")
		return
	}
}

func (l *listener) ExitDockerPullStmt(c *parser.DockerPullStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for DOCKER PULL: %s", l.stmtWords)
		return
	}
	imageName := l.stmtWords[0]
	err := l.interpreter.DockerPull(l.ctx, imageName)
	if err != nil {
		l.err = errors.Wrap(err, "docker pull")
		return
	}
}

func (l *listener) ExitHealthcheckStmt(c *parser.HealthcheckStmtContext) {
	if l.shouldSkip() {
		return
	}
	if l.pushOnlyAllowed {
		l.err = fmt.Errorf("no non-push commands allowed after a --push: %s", c.GetText())
		return
	}
	fs := flag.NewFlagSet("HEALTHCHECK", flag.ContinueOnError)
	interval := fs.Duration("interval", 30*time.Second, "")
	timeout := fs.Duration("timeout", 30*time.Second, "")
	startPeriod := fs.Duration("start-period", 0, "")
	retries := fs.Int("retries", 3, "")
	err := fs.Parse(l.stmtWords)
	if err != nil {
		l.err = errors.Wrapf(err, "invalid HEALTHCHECK arguments %v", l.stmtWords)
		return
	}
	if fs.NArg() == 0 {
		l.err = fmt.Errorf("invalid number of arguments for HEALTHCHECK: %s", l.stmtWords)
		return
	}
	isNone := false
	var cmdArgs []string
	switch fs.Arg(0) {
	case "NONE":
		if fs.NArg() != 1 {
			l.err = fmt.Errorf("invalid arguments for HEALTHCHECK: %s", l.stmtWords)
			return
		}
		isNone = true
	case "CMD":
		if fs.NArg() == 1 {
			l.err = fmt.Errorf("invalid number of arguments for HEALTHCHECK CMD: %s", l.stmtWords)
			return
		}
		cmdArgs = fs.Args()[1:]
	default:
		if strings.HasPrefix(fs.Arg(0), "[") {
			l.err = fmt.Errorf("exec form not yet supported for HEALTHCHECK CMD: %s", l.stmtWords)
			return
		}
		l.err = fmt.Errorf("invalid arguments for HEALTHCHECK: %s", l.stmtWords)
		return
	}
	l.interpreter.Healthcheck(l.ctx, isNone, cmdArgs, *interval, *timeout, *startPeriod, *retries)
}

func (l *listener) ExitAddStmt(c *parser.AddStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command ADD not yet supported")
}

func (l *listener) ExitStopsignalStmt(c *parser.StopsignalStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command STOPSIGNAL not yet supported")
}

func (l *listener) ExitOnbuildStmt(c *parser.OnbuildStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command ONBUILD not supported")
}

func (l *listener) ExitShellStmt(c *parser.ShellStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Command SHELL not yet supported")
}

func (l *listener) ExitGenericCommandStmt(c *parser.GenericCommandStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.err = fmt.Errorf("Invalid command %s", c.GetText())
}

//
// Variables.

func (l *listener) EnterEnvArgKey(c *parser.EnvArgKeyContext) {
	if l.shouldSkip() {
		return
	}
	l.envArgKey = c.GetText()
}

func (l *listener) EnterEnvArgValue(c *parser.EnvArgValueContext) {
	if l.shouldSkip() {
		return
	}
	l.envArgValue = c.GetText()
}

func (l *listener) EnterLabelKey(c *parser.LabelKeyContext) {
	if l.shouldSkip() {
		return
	}
	l.labelKeys = append(l.labelKeys, c.GetText())
}

func (l *listener) EnterLabelValue(c *parser.LabelValueContext) {
	if l.shouldSkip() {
		return
	}
	l.labelValues = append(l.labelValues, c.GetText())
}

func (l *listener) ExitStmtWordsMaybeJSON(c *parser.StmtWordsMaybeJSONContext) {
	if l.shouldSkip() {
		return
	}
	// Try to parse as JSON. If parse works, override the already collected stmtWords.
	var words []string
	err := json.Unmarshal([]byte(c.GetText()), &words)
	if err == nil {
		l.stmtWords = words
		l.execMode = true
	}
}

func (l *listener) EnterStmtWord(c *parser.StmtWordContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = append(l.stmtWords, replaceEscape(c.GetText()))
}

func (l *listener) shouldSkip() bool {
	return l.err != nil || l.currentTarget != l.executeTarget
}

// StringSliceFlag is a flag backed by a string slice.
type StringSliceFlag struct {
	Args []string
}

// String returns a string representation of the flag.
func (ssf *StringSliceFlag) String() string {
	if ssf == nil {
		return ""
	}
	return strings.Join(ssf.Args, ",")
}

// Set adds a flag value to the string slice.
func (ssf *StringSliceFlag) Set(arg string) error {
	ssf.Args = append(ssf.Args, arg)
	return nil
}

var lineContinuationRegexp = regexp.MustCompile("\\\\(\\n|(\\r\\n))[\\t ]*")

func replaceEscape(str string) string {
	return lineContinuationRegexp.ReplaceAllString(str, "")
}
