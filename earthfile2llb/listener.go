package earthfile2llb

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vladaionescu/earthly/earthfile2llb/parser"
)

var _ parser.EarthParserListener = &listener{}

type listener struct {
	*parser.BaseEarthParserListener
	converter *Converter
	ctx       context.Context

	executeTarget   string
	currentTarget   string
	saveImageExists bool

	envArgKey   string
	envArgValue string

	execMode  bool
	stmtWords []string

	err error
}

func newListener(ctx context.Context, converter *Converter, executeTarget string) *listener {
	return &listener{
		ctx:           ctx,
		converter:     converter,
		executeTarget: executeTarget,
		currentTarget: "base",
	}
}

func (l *listener) EnterTargetHeader(c *parser.TargetHeaderContext) {
	// Apply implicit SAVE IMAGE for +base.
	if l.executeTarget == "base" {
		if !l.saveImageExists {
			l.converter.SaveImage(l.ctx, []string{})
		}
		l.saveImageExists = true
	}

	l.currentTarget = strings.TrimSuffix(c.GetText(), ":")
	if l.shouldSkip() {
		return
	}
	if l.currentTarget == "base" {
		l.err = errors.New("Target name cannot be base")
		return
	}
	// Apply implicit FROM +base
	err := l.converter.From(l.ctx, "+base", nil)
	if err != nil {
		l.err = errors.Wrap(err, "apply implicit FROM +base")
		return
	}
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
	l.execMode = false
}

func (l *listener) ExitFromStmt(c *parser.FromStmtContext) {
	if l.shouldSkip() {
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
	err = l.converter.From(l.ctx, imageName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrapf(err, "apply FROM %s", imageName)
		return
	}
}

func (l *listener) ExitCopyStmt(c *parser.CopyStmtContext) {
	if l.shouldSkip() {
		return
	}
	fs := flag.NewFlagSet("COPY", flag.ContinueOnError)
	isArtifactCopy := fs.Bool("artifact", false, "")
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
	if *from != "" && *isArtifactCopy {
		l.err = fmt.Errorf("invalid COPY flags %v: . The flags --from and --artifact cannot both be specified at the same time", l.stmtWords)
		return
	}
	if *from != "" {
		l.err = errors.New("COPY --from not implemented. Use COPY --artifact instead")
		return
	}
	if *isArtifactCopy {
		if fs.NArg() != 2 {
			l.err = errors.New("more than 2 COPY arguments not yet supported for --artifact")
			return
		}
		artifactName := fs.Arg(0)
		dest := fs.Arg(1)
		err = l.converter.CopyArtifact(l.ctx, artifactName, dest, buildArgs.Args, *isDirCopy)
		if err != nil {
			l.err = errors.Wrapf(err, "copy artifact")
			return
		}
	} else {
		if len(buildArgs.Args) != 0 {
			l.err = fmt.Errorf("build args not supported for non --artifact case %v", l.stmtWords)
			return
		}
		srcs := fs.Args()[:fs.NArg()-1]
		dest := fs.Arg(fs.NArg() - 1)
		l.converter.CopyClassical(l.ctx, srcs, dest, *isDirCopy)
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
	// TODO: In the bracket case, should flags be outside of the brackets?

	err = l.converter.Run(l.ctx, fs.Args(), mounts.Args, secrets.Args, *privileged, *withEntrypoint, *withDocker, withShell)
	if err != nil {
		l.err = errors.Wrap(err, "run")
		return
	}
}

func (l *listener) ExitSaveArtifact(c *parser.SaveArtifactContext) {
	if l.shouldSkip() {
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

	l.converter.SaveArtifact(l.ctx, saveFrom, saveTo, saveAsLocalTo)
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

	l.converter.SaveImage(l.ctx, l.stmtWords)
}

func (l *listener) ExitBuildStmt(c *parser.BuildStmtContext) {
	if l.shouldSkip() {
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
	_, err = l.converter.Build(l.ctx, fullTargetName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrapf(err, "apply BUILD %s", fullTargetName)
		return
	}
}

func (l *listener) ExitWorkdirStmt(c *parser.WorkdirStmtContext) {
	if l.shouldSkip() {
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for WORKDIR: %v", l.stmtWords)
		return
	}
	workdirPath := l.stmtWords[0]
	l.converter.Workdir(l.ctx, workdirPath)
}

func (l *listener) ExitEntrypointStmt(c *parser.EntrypointStmtContext) {
	if l.shouldSkip() {
		return
	}
	withShell := !l.execMode
	l.converter.Entrypoint(l.ctx, l.stmtWords, withShell)
}

func (l *listener) ExitEnvStmt(c *parser.EnvStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.Env(l.ctx, l.envArgKey, l.envArgValue)
}

func (l *listener) ExitArgStmt(c *parser.ArgStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.Arg(l.ctx, l.envArgKey, l.envArgValue)
}

func (l *listener) ExitGitCloneStmt(c *parser.GitCloneStmtContext) {
	if l.shouldSkip() {
		return
	}
	fs := flag.NewFlagSet("GIT CLONE", flag.ContinueOnError)
	branch := fs.String("build-arg", "", "")
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
	err = l.converter.GitClone(l.ctx, gitURL, *branch, gitCloneDest)
	if err != nil {
		l.err = errors.Wrap(err, "git clone")
		return
	}
}

func (l *listener) ExitDockerLoadStmt(c *parser.DockerLoadStmtContext) {
	if l.shouldSkip() {
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
	err = l.converter.DockerLoad(l.ctx, fullTargetName, imageName, buildArgs.Args)
	if err != nil {
		l.err = errors.Wrap(err, "docker load")
		return
	}
}

func (l *listener) ExitDockerPullStmt(c *parser.DockerPullStmtContext) {
	if l.shouldSkip() {
		return
	}
	if len(l.stmtWords) != 1 {
		l.err = fmt.Errorf("invalid number of arguments for DOCKER PULL: %s", l.stmtWords)
		return
	}
	imageName := l.stmtWords[0]
	err := l.converter.DockerPull(l.ctx, imageName)
	if err != nil {
		l.err = errors.Wrap(err, "docker pull")
		return
	}
}

//
// Variables.

func (l *listener) EnterStmtWordsList(c *parser.StmtWordsListContext) {
	if l.shouldSkip() {
		return
	}
	l.execMode = true
}

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

func (l *listener) EnterStmtWord(c *parser.StmtWordContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = append(l.stmtWords, c.GetText())
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
