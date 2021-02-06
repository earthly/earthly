package ast

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/earthly/earthly/ast/parser"
	"github.com/earthly/earthly/ast/spec"
)

var _ parser.EarthParserListener = &listener{}

type contextFrame int

const (
	contextFrameEarthfile = iota
	contextFrameTarget
	contextFrameStatement
	contextFrameWith
)

type listener struct {
	*parser.BaseEarthParserListener

	ef        *spec.Earthfile
	target    *spec.Target
	block     spec.Block
	statement *spec.Statement
	command   *spec.Command
	with      *spec.WithStatement

	contextStack []contextFrame

	ctx             context.Context
	filePath        string
	enableSourceMap bool

	envArgKey   string
	envArgValue string
	labelKeys   []string
	labelValues []string

	err error
}

func newListener(ctx context.Context, filePath string, enableSourceMap bool) *listener {
	return &listener{
		ctx:             ctx,
		filePath:        filePath,
		enableSourceMap: enableSourceMap,
		ef:              new(spec.Earthfile),
		contextStack:    []contextFrame{contextFrameEarthfile},
	}
}

func (l *listener) Err() error {
	if l.currentContext() != contextFrameEarthfile && l.err == nil {
		return errors.New("parsing did not finish")
	}
	return l.err
}

func (l *listener) Earthfile() spec.Earthfile {
	return *l.ef
}

func (l *listener) currentContext() contextFrame {
	return l.contextStack[len(l.contextStack)-1]
}

func (l *listener) pushContext(frame contextFrame) {
	l.contextStack = append(l.contextStack, frame)
}

func (l *listener) popContext(frame contextFrame) contextFrame {
	if l.currentContext() != frame {
		panic(fmt.Sprintf("unexpected frame %v when expecting frame %v", l.currentContext(), frame))
	}
	l.contextStack = l.contextStack[:len(l.contextStack)-1]
	return l.currentContext()
}

// Target ---------------------------------------------------------------------

func (l *listener) EnterTarget(c *parser.TargetContext) {
	l.target = new(spec.Target)
	if l.enableSourceMap {
		l.target.SourceLocation = &spec.SourceLocation{
			File:        l.filePath,
			StartLine:   c.GetStart().GetLine(),
			StartColumn: c.GetStart().GetColumn(),
			EndLine:     c.GetStop().GetLine(),
			EndColumn:   c.GetStop().GetColumn(),
		}
	}
	l.pushContext(contextFrameTarget)
}

func (l *listener) EnterTargetHeader(c *parser.TargetHeaderContext) {
	l.target.Name = strings.TrimSuffix(c.GetText(), ":")
}

func (l *listener) ExitTarget(c *parser.TargetContext) {
	l.ef.Targets = append(l.ef.Targets, *l.target)
	l.target = nil
	l.popContext(contextFrameTarget)
}

// Block ----------------------------------------------------------------------

func (l *listener) EnterStmts(c *parser.StmtsContext) {
	l.block = []spec.Statement{}
}

func (l *listener) ExitStmts(c *parser.StmtsContext) {
	switch l.currentContext() {
	case contextFrameEarthfile:
		l.ef.BaseRecipe = l.block
	case contextFrameTarget:
		l.target.Recipe = l.block
	case contextFrameWith:
		l.with.Body = l.block
	default:
		panic(fmt.Sprintf("unhandled block for context %v", l.currentContext()))
	}
	l.block = nil
}

// Statement ------------------------------------------------------------------

func (l *listener) EnterStmt(c *parser.StmtContext) {
	l.statement = new(spec.Statement)
	if l.enableSourceMap {
		l.statement.SourceLocation = &spec.SourceLocation{
			File:        l.filePath,
			StartLine:   c.GetStart().GetLine(),
			StartColumn: c.GetStart().GetColumn(),
			EndLine:     c.GetStop().GetLine(),
			EndColumn:   c.GetStop().GetColumn(),
		}
	}
	l.pushContext(contextFrameStatement)
}

func (l *listener) ExitStmt(c *parser.StmtContext) {
	l.block = append(l.block, *l.statement)
	l.statement = nil
	l.popContext(contextFrameStatement)
}

// Command --------------------------------------------------------------------

func (l *listener) EnterCommand(c *parser.CommandContext) {
	l.command = new(spec.Command)
	if l.enableSourceMap {
		l.command.SourceLocation = &spec.SourceLocation{
			File:        l.filePath,
			StartLine:   c.GetStart().GetLine(),
			StartColumn: c.GetStart().GetColumn(),
			EndLine:     c.GetStop().GetLine(),
			EndColumn:   c.GetStop().GetColumn(),
		}
	}
}

func (l *listener) ExitCommand(c *parser.CommandContext) {
	switch l.currentContext() {
	case contextFrameStatement:
		l.statement.Command = l.command
	default:
		panic(fmt.Sprintf("unhandled command for context %v", l.currentContext()))
	}
	l.command = nil
}

// Individual commands --------------------------------------------------------

func (l *listener) EnterFromStmt(c *parser.FromStmtContext) {
	l.command.Name = "FROM"
}

func (l *listener) EnterFromDockerfileStmt(c *parser.FromDockerfileStmtContext) {
	l.command.Name = "FROM DOCKERFILE"
}

func (l *listener) EnterLocallyStmt(c *parser.LocallyStmtContext) {
	l.command.Name = "LOCALLY"
}

func (l *listener) EnterCopyStmt(c *parser.CopyStmtContext) {
	l.command.Name = "COPY"
}

func (l *listener) EnterRunStmt(c *parser.RunStmtContext) {
	l.command.Name = "RUN"
}

func (l *listener) EnterSaveArtifact(c *parser.SaveArtifactContext) {
	l.command.Name = "SAVE ARTIFACT"
}

func (l *listener) EnterSaveImage(c *parser.SaveImageContext) {
	l.command.Name = "SAVE IMAGE"
}

func (l *listener) EnterBuildStmt(c *parser.BuildStmtContext) {
	l.command.Name = "BUILD"
}

func (l *listener) EnterWorkdirStmt(c *parser.WorkdirStmtContext) {
	l.command.Name = "WORKDIR"
}

func (l *listener) EnterUserStmt(c *parser.UserStmtContext) {
	l.command.Name = "USER"
}

func (l *listener) EnterCmdStmt(c *parser.CmdStmtContext) {
	l.command.Name = "CMD"
}

func (l *listener) EnterEntrypointStmt(c *parser.EntrypointStmtContext) {
	l.command.Name = "ENTRYPOINT"
}

func (l *listener) EnterExposeStmt(c *parser.ExposeStmtContext) {
	l.command.Name = "EXPOSE"
}

func (l *listener) EnterVolumeStmt(c *parser.VolumeStmtContext) {
	l.command.Name = "VOLUME"
}

func (l *listener) EnterEnvStmt(c *parser.EnvStmtContext) {
	l.command.Name = "ENV"
}

func (l *listener) EnterArgStmt(c *parser.ArgStmtContext) {
	l.command.Name = "ARG"
}

func (l *listener) EnterLabelStmt(c *parser.LabelStmtContext) {
	l.command.Name = "LABEL"
}

func (l *listener) EnterGitCloneStmt(c *parser.GitCloneStmtContext) {
	l.command.Name = "GIT CLONE"
}

func (l *listener) EnterHealthcheckStmt(c *parser.HealthcheckStmtContext) {
	l.command.Name = "HEALTHCHECK"
}

func (l *listener) EnterAddStmt(c *parser.AddStmtContext) {
	l.command.Name = "ADD"
}

func (l *listener) EnterStopsignalStmt(c *parser.StopsignalStmtContext) {
	l.command.Name = "STOP SIGNAL"
}

func (l *listener) EnterOnbuildStmt(c *parser.OnbuildStmtContext) {
	l.command.Name = "ONBUILD"
}

func (l *listener) EnterShellStmt(c *parser.ShellStmtContext) {
	l.command.Name = "SHELL"
}

func (l *listener) EnterGenericCommandStmt(c *parser.GenericCommandStmtContext) {
	// Set in EnterCommandName.
}

func (l *listener) EnterCommandName(c *parser.CommandNameContext) {
	l.command.Name = c.GetText()
}

// With -----------------------------------------------------------------------

func (l *listener) EnterWithDockerStmt(c *parser.WithDockerStmtContext) {
	// TODO: Reuse EnterCommand.
	l.with = new(spec.WithStatement)
	l.command = new(spec.Command)
	if l.enableSourceMap {
		l.command.SourceLocation = &spec.SourceLocation{
			File:        l.filePath,
			StartLine:   c.GetStart().GetLine(),
			StartColumn: c.GetStart().GetColumn(),
			EndLine:     c.GetStop().GetLine(),
			EndColumn:   c.GetStop().GetColumn(),
		}
	}
	l.pushContext(contextFrameWith)
}

func (l *listener) ExitWithDockerStmt(c *parser.WithDockerStmtContext) {
	l.with.Command = *l.command
	l.command = nil
}

func (l *listener) ExitEndStmt(c *parser.EndStmtContext) {
	l.statement.With = l.with
	l.with = nil
	l.popContext(contextFrameWith)
}

// EnvArgKey, EnvArgValue, LabelKey, LabelValue -------------------------------

func (l *listener) EnterEnvArgKey(c *parser.EnvArgKeyContext) {
	err := checkEnvVarName(c.GetText())
	if err != nil {
		l.err = err
		return
	}
	l.command.Args = append(l.command.Args, c.GetText())
}

func (l *listener) EnterEnvArgValue(c *parser.EnvArgValueContext) {
	l.command.Args = append(l.command.Args, "=", c.GetText())
}

func (l *listener) EnterLabelKey(c *parser.LabelKeyContext) {
	l.command.Args = append(l.command.Args, c.GetText())
}

func (l *listener) EnterLabelValue(c *parser.LabelValueContext) {
	l.command.Args = append(l.command.Args, "=", c.GetText())
}

// StmtWord -------------------------------------------------------------------

func (l *listener) ExitStmtWordsMaybeJSON(c *parser.StmtWordsMaybeJSONContext) {
	// Try to parse as JSON. If parse works, override the already collected stmtWords.
	var words []string
	err := json.Unmarshal([]byte(c.GetText()), &words)
	if err == nil {
		l.command.Args = words
		l.command.ExecMode = true
	}
}

func (l *listener) EnterStmtWord(c *parser.StmtWordContext) {
	l.command.Args = append(l.command.Args, replaceEscape(c.GetText()))
}

// ----------------------------------------------------------------------------

var envVarNameRegexp = regexp.MustCompile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")

func checkEnvVarName(str string) error {
	itMatch := envVarNameRegexp.MatchString(str)
	if !itMatch {
		return fmt.Errorf("invalid env key definition %s", str)
	}
	return nil
}

var lineContinuationRegexp = regexp.MustCompile("\\\\(\\n|(\\r\\n))[\\t ]*")

func replaceEscape(str string) string {
	return lineContinuationRegexp.ReplaceAllString(str, "")
}
