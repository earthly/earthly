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

	imageName      string
	saveImageNames []string
	asName         string
	fullTargetName string
	saveFrom       string
	saveTo         string
	saveAsLocalTo  string
	workdirPath    string
	envArgKey      string
	envArgValue    string
	gitURL         string
	gitCloneDest   string
	flagKeyValues  []string

	isListWithBrackets bool
	stmtWords          []string

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

func (l *listener) EnterFromStmt(c *parser.FromStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.flagKeyValues = nil
	l.imageName = ""
	l.asName = ""
}

func (l *listener) ExitFromStmt(c *parser.FromStmtContext) {
	if l.shouldSkip() {
		return
	}
	buildArgs, err := parseBuildArgFlags(l.flagKeyValues)
	if err != nil {
		l.err = errors.Wrap(err, "parse build arg flags")
		return
	}
	err = l.converter.From(l.ctx, l.imageName, buildArgs)
	if err != nil {
		l.err = errors.Wrapf(err, "apply FROM %s", l.imageName)
		return
	}
}

func (l *listener) EnterCopyStmt(c *parser.CopyStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = nil
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
		l.err = fmt.Errorf("Not enough COPY arguments %v", l.stmtWords)
		return
	}
	if *from != "" && *isArtifactCopy {
		l.err = fmt.Errorf("Invalid COPY flags %v: . The flags --from and --artifact cannot both be specified at the same time", l.stmtWords)
		return
	}
	if *from != "" {
		l.err = errors.New("COPY --from not implemented. Use COPY --artifact instead")
		return
	}
	if *isArtifactCopy {
		if fs.NArg() != 2 {
			l.err = errors.New("More than 2 COPY arguments not yet supported for --artifact")
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
			l.err = fmt.Errorf("Build args not supported for non --artifact case %v", l.stmtWords)
			return
		}
		srcs := fs.Args()[:fs.NArg()-1]
		dest := fs.Arg(fs.NArg() - 1)
		l.converter.CopyClassical(l.ctx, srcs, dest, *isDirCopy)
	}
}

func (l *listener) EnterRunStmt(c *parser.RunStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = nil
	l.isListWithBrackets = false
}

func (l *listener) ExitRunStmt(c *parser.RunStmtContext) {
	if l.shouldSkip() {
		return
	}
	if len(l.stmtWords) < 1 {
		l.err = errors.New("Not enough arguments for RUN")
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
	withShell := !l.isListWithBrackets
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

func (l *listener) EnterSaveArtifact(c *parser.SaveArtifactContext) {
	if l.shouldSkip() {
		return
	}
	l.saveFrom = ""
	l.saveTo = ""
	l.saveAsLocalTo = ""
}

func (l *listener) ExitSaveArtifact(c *parser.SaveArtifactContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.SaveArtifact(l.ctx, l.saveFrom, l.saveTo, l.saveAsLocalTo)
}

func (l *listener) EnterSaveImage(c *parser.SaveImageContext) {
	if l.shouldSkip() {
		return
	}
	if l.saveImageExists {
		l.err = fmt.Errorf(
			"More than one SAVE IMAGE statement per target not allowed: %s", c.GetText())
		return
	}
	l.saveImageExists = true

	l.saveImageNames = nil
}

func (l *listener) ExitSaveImage(c *parser.SaveImageContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.SaveImage(l.ctx, l.saveImageNames)
}

func (l *listener) EnterBuildStmt(c *parser.BuildStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.fullTargetName = ""
	l.flagKeyValues = nil
}

func (l *listener) ExitBuildStmt(c *parser.BuildStmtContext) {
	if l.shouldSkip() {
		return
	}
	buildArgs, err := parseBuildArgFlags(l.flagKeyValues)
	if err != nil {
		l.err = errors.Wrap(err, "parse build arg flags")
		return
	}
	_, err = l.converter.Build(l.ctx, l.fullTargetName, buildArgs)
	if err != nil {
		l.err = errors.Wrapf(err, "apply BUILD %s", l.fullTargetName)
		return
	}
}

func (l *listener) EnterWorkdirStmt(c *parser.WorkdirStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.workdirPath = ""
}

func (l *listener) ExitWorkdirStmt(c *parser.WorkdirStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.Workdir(l.ctx, l.workdirPath)
}

func (l *listener) EnterEntrypointStmt(c *parser.EntrypointStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.stmtWords = nil
	l.isListWithBrackets = false
}

func (l *listener) ExitEntrypointStmt(c *parser.EntrypointStmtContext) {
	if l.shouldSkip() {
		return
	}
	withShell := !l.isListWithBrackets
	l.converter.Entrypoint(l.ctx, l.stmtWords, withShell)
}

func (l *listener) EnterEnvStmt(c *parser.EnvStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.envArgKey = ""
	l.envArgValue = ""
}

func (l *listener) ExitEnvStmt(c *parser.EnvStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.Env(l.ctx, l.envArgKey, l.envArgValue)
}

func (l *listener) EnterArgStmt(c *parser.ArgStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.envArgKey = ""
	l.envArgValue = ""
}

func (l *listener) ExitArgStmt(c *parser.ArgStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.converter.Arg(l.ctx, l.envArgKey, l.envArgValue)
}

func (l *listener) EnterGitCloneStmt(c *parser.GitCloneStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.gitURL = ""
	l.flagKeyValues = nil
	l.gitCloneDest = ""
}

func (l *listener) ExitGitCloneStmt(c *parser.GitCloneStmtContext) {
	if l.shouldSkip() {
		return
	}
	branch, err := parseGitCloneFlags(l.flagKeyValues)
	if err != nil {
		l.err = errors.Wrap(err, "parse git clone flags")
		return
	}
	err = l.converter.GitClone(l.ctx, l.gitURL, branch, l.gitCloneDest)
	if err != nil {
		l.err = errors.Wrap(err, "git clone")
		return
	}
}

func (l *listener) EnterDockerLoadStmt(c *parser.DockerLoadStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.fullTargetName = ""
	l.imageName = ""
	l.flagKeyValues = nil
}

func (l *listener) ExitDockerLoadStmt(c *parser.DockerLoadStmtContext) {
	if l.shouldSkip() {
		return
	}
	buildArgs, err := parseBuildArgFlags(l.flagKeyValues)
	if err != nil {
		l.err = errors.Wrap(err, "parse build arg flags")
		return
	}
	err = l.converter.DockerLoad(l.ctx, l.fullTargetName, l.imageName, buildArgs)
	if err != nil {
		l.err = errors.Wrap(err, "docker load")
		return
	}
}

func (l *listener) EnterDockerPullStmt(c *parser.DockerPullStmtContext) {
	if l.shouldSkip() {
		return
	}
	l.imageName = ""
}

func (l *listener) ExitDockerPullStmt(c *parser.DockerPullStmtContext) {
	if l.shouldSkip() {
		return
	}
	err := l.converter.DockerPull(l.ctx, l.imageName)
	if err != nil {
		l.err = errors.Wrap(err, "docker pull")
		return
	}
}

//
// Variables.

func (l *listener) EnterImageName(c *parser.ImageNameContext) {
	if l.shouldSkip() {
		return
	}
	l.imageName = c.GetText()
}

func (l *listener) EnterSaveImageName(c *parser.SaveImageNameContext) {
	if l.shouldSkip() {
		return
	}
	l.saveImageNames = append(l.saveImageNames, c.GetText())
}

func (l *listener) EnterAsName(c *parser.AsNameContext) {
	if l.shouldSkip() {
		return
	}
	l.asName = c.GetText()
}

func (l *listener) EnterStmtWordsList(c *parser.StmtWordsListContext) {
	if l.shouldSkip() {
		return
	}
	l.isListWithBrackets = true
}

func (l *listener) EnterSaveFrom(c *parser.SaveFromContext) {
	if l.shouldSkip() {
		return
	}
	l.saveFrom = c.GetText()
}

func (l *listener) EnterSaveTo(c *parser.SaveToContext) {
	if l.shouldSkip() {
		return
	}
	l.saveTo = c.GetText()
}

func (l *listener) EnterSaveAsLocalTo(c *parser.SaveAsLocalToContext) {
	if l.shouldSkip() {
		return
	}
	l.saveAsLocalTo = c.GetText()
}

func (l *listener) EnterFullTargetName(c *parser.FullTargetNameContext) {
	if l.shouldSkip() {
		return
	}
	l.fullTargetName = c.GetText()
}

func (l *listener) EnterWorkdirPath(c *parser.WorkdirPathContext) {
	if l.shouldSkip() {
		return
	}
	l.workdirPath = c.GetText()
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

func (l *listener) EnterFlagKeyValue(c *parser.FlagKeyValueContext) {
	if l.shouldSkip() {
		return
	}
	l.flagKeyValues = append(l.flagKeyValues, c.GetText())
}

func (l *listener) EnterFlagKey(c *parser.FlagKeyContext) {
	if l.shouldSkip() {
		return
	}
	l.flagKeyValues = append(l.flagKeyValues, c.GetText())
}

func (l *listener) EnterGitURL(c *parser.GitURLContext) {
	if l.shouldSkip() {
		return
	}
	l.gitURL = c.GetText()
}

func (l *listener) EnterGitCloneDest(c *parser.GitCloneDestContext) {
	if l.shouldSkip() {
		return
	}
	l.gitCloneDest = c.GetText()
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

func parseBuildArgFlags(flagKeyValues []string) ([]string, error) {
	var out []string
	for _, flag := range flagKeyValues {
		split := strings.SplitN(flag, "=", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("Invalid flag format %s", flag)
		}
		if split[0] != "--build-arg" {
			return nil, fmt.Errorf("Invalid flag %s", split[0])
		}
		out = append(out, split[1])
	}
	return out, nil
}

func parseGitCloneFlags(flagKeyValues []string) (string, error) {
	branch := ""
	for _, flag := range flagKeyValues {
		split := strings.SplitN(flag, "=", 2)
		if len(split) != 2 {
			return "", fmt.Errorf("Invalid flag format %s", flag)
		}
		if split[0] != "--branch" {
			return "", fmt.Errorf("Invalid flag %s", split[0])
		}
		branch = split[1]
	}
	return branch, nil
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
