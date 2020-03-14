// Code generated from earthfile2llb/parser/EarthParser.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // EarthParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseEarthParserListener is a complete listener for a parse tree produced by EarthParser.
type BaseEarthParserListener struct{}

var _ EarthParserListener = &BaseEarthParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseEarthParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseEarthParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseEarthParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseEarthParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterEarthFile is called when production earthFile is entered.
func (s *BaseEarthParserListener) EnterEarthFile(ctx *EarthFileContext) {}

// ExitEarthFile is called when production earthFile is exited.
func (s *BaseEarthParserListener) ExitEarthFile(ctx *EarthFileContext) {}

// EnterTargets is called when production targets is entered.
func (s *BaseEarthParserListener) EnterTargets(ctx *TargetsContext) {}

// ExitTargets is called when production targets is exited.
func (s *BaseEarthParserListener) ExitTargets(ctx *TargetsContext) {}

// EnterTarget is called when production target is entered.
func (s *BaseEarthParserListener) EnterTarget(ctx *TargetContext) {}

// ExitTarget is called when production target is exited.
func (s *BaseEarthParserListener) ExitTarget(ctx *TargetContext) {}

// EnterTargetHeader is called when production targetHeader is entered.
func (s *BaseEarthParserListener) EnterTargetHeader(ctx *TargetHeaderContext) {}

// ExitTargetHeader is called when production targetHeader is exited.
func (s *BaseEarthParserListener) ExitTargetHeader(ctx *TargetHeaderContext) {}

// EnterStmts is called when production stmts is entered.
func (s *BaseEarthParserListener) EnterStmts(ctx *StmtsContext) {}

// ExitStmts is called when production stmts is exited.
func (s *BaseEarthParserListener) ExitStmts(ctx *StmtsContext) {}

// EnterStmt is called when production stmt is entered.
func (s *BaseEarthParserListener) EnterStmt(ctx *StmtContext) {}

// ExitStmt is called when production stmt is exited.
func (s *BaseEarthParserListener) ExitStmt(ctx *StmtContext) {}

// EnterFromStmt is called when production fromStmt is entered.
func (s *BaseEarthParserListener) EnterFromStmt(ctx *FromStmtContext) {}

// ExitFromStmt is called when production fromStmt is exited.
func (s *BaseEarthParserListener) ExitFromStmt(ctx *FromStmtContext) {}

// EnterCopyStmt is called when production copyStmt is entered.
func (s *BaseEarthParserListener) EnterCopyStmt(ctx *CopyStmtContext) {}

// ExitCopyStmt is called when production copyStmt is exited.
func (s *BaseEarthParserListener) ExitCopyStmt(ctx *CopyStmtContext) {}

// EnterSaveStmt is called when production saveStmt is entered.
func (s *BaseEarthParserListener) EnterSaveStmt(ctx *SaveStmtContext) {}

// ExitSaveStmt is called when production saveStmt is exited.
func (s *BaseEarthParserListener) ExitSaveStmt(ctx *SaveStmtContext) {}

// EnterSaveImage is called when production saveImage is entered.
func (s *BaseEarthParserListener) EnterSaveImage(ctx *SaveImageContext) {}

// ExitSaveImage is called when production saveImage is exited.
func (s *BaseEarthParserListener) ExitSaveImage(ctx *SaveImageContext) {}

// EnterSaveArtifact is called when production saveArtifact is entered.
func (s *BaseEarthParserListener) EnterSaveArtifact(ctx *SaveArtifactContext) {}

// ExitSaveArtifact is called when production saveArtifact is exited.
func (s *BaseEarthParserListener) ExitSaveArtifact(ctx *SaveArtifactContext) {}

// EnterRunStmt is called when production runStmt is entered.
func (s *BaseEarthParserListener) EnterRunStmt(ctx *RunStmtContext) {}

// ExitRunStmt is called when production runStmt is exited.
func (s *BaseEarthParserListener) ExitRunStmt(ctx *RunStmtContext) {}

// EnterBuildStmt is called when production buildStmt is entered.
func (s *BaseEarthParserListener) EnterBuildStmt(ctx *BuildStmtContext) {}

// ExitBuildStmt is called when production buildStmt is exited.
func (s *BaseEarthParserListener) ExitBuildStmt(ctx *BuildStmtContext) {}

// EnterWorkdirStmt is called when production workdirStmt is entered.
func (s *BaseEarthParserListener) EnterWorkdirStmt(ctx *WorkdirStmtContext) {}

// ExitWorkdirStmt is called when production workdirStmt is exited.
func (s *BaseEarthParserListener) ExitWorkdirStmt(ctx *WorkdirStmtContext) {}

// EnterEntrypointStmt is called when production entrypointStmt is entered.
func (s *BaseEarthParserListener) EnterEntrypointStmt(ctx *EntrypointStmtContext) {}

// ExitEntrypointStmt is called when production entrypointStmt is exited.
func (s *BaseEarthParserListener) ExitEntrypointStmt(ctx *EntrypointStmtContext) {}

// EnterEnvStmt is called when production envStmt is entered.
func (s *BaseEarthParserListener) EnterEnvStmt(ctx *EnvStmtContext) {}

// ExitEnvStmt is called when production envStmt is exited.
func (s *BaseEarthParserListener) ExitEnvStmt(ctx *EnvStmtContext) {}

// EnterArgStmt is called when production argStmt is entered.
func (s *BaseEarthParserListener) EnterArgStmt(ctx *ArgStmtContext) {}

// ExitArgStmt is called when production argStmt is exited.
func (s *BaseEarthParserListener) ExitArgStmt(ctx *ArgStmtContext) {}

// EnterEnvArgKey is called when production envArgKey is entered.
func (s *BaseEarthParserListener) EnterEnvArgKey(ctx *EnvArgKeyContext) {}

// ExitEnvArgKey is called when production envArgKey is exited.
func (s *BaseEarthParserListener) ExitEnvArgKey(ctx *EnvArgKeyContext) {}

// EnterEnvArgValue is called when production envArgValue is entered.
func (s *BaseEarthParserListener) EnterEnvArgValue(ctx *EnvArgValueContext) {}

// ExitEnvArgValue is called when production envArgValue is exited.
func (s *BaseEarthParserListener) ExitEnvArgValue(ctx *EnvArgValueContext) {}

// EnterGitCloneStmt is called when production gitCloneStmt is entered.
func (s *BaseEarthParserListener) EnterGitCloneStmt(ctx *GitCloneStmtContext) {}

// ExitGitCloneStmt is called when production gitCloneStmt is exited.
func (s *BaseEarthParserListener) ExitGitCloneStmt(ctx *GitCloneStmtContext) {}

// EnterDockerLoadStmt is called when production dockerLoadStmt is entered.
func (s *BaseEarthParserListener) EnterDockerLoadStmt(ctx *DockerLoadStmtContext) {}

// ExitDockerLoadStmt is called when production dockerLoadStmt is exited.
func (s *BaseEarthParserListener) ExitDockerLoadStmt(ctx *DockerLoadStmtContext) {}

// EnterDockerPullStmt is called when production dockerPullStmt is entered.
func (s *BaseEarthParserListener) EnterDockerPullStmt(ctx *DockerPullStmtContext) {}

// ExitDockerPullStmt is called when production dockerPullStmt is exited.
func (s *BaseEarthParserListener) ExitDockerPullStmt(ctx *DockerPullStmtContext) {}

// EnterGenericCommand is called when production genericCommand is entered.
func (s *BaseEarthParserListener) EnterGenericCommand(ctx *GenericCommandContext) {}

// ExitGenericCommand is called when production genericCommand is exited.
func (s *BaseEarthParserListener) ExitGenericCommand(ctx *GenericCommandContext) {}

// EnterCommandName is called when production commandName is entered.
func (s *BaseEarthParserListener) EnterCommandName(ctx *CommandNameContext) {}

// ExitCommandName is called when production commandName is exited.
func (s *BaseEarthParserListener) ExitCommandName(ctx *CommandNameContext) {}

// EnterStmtWordsMaybeJSON is called when production stmtWordsMaybeJSON is entered.
func (s *BaseEarthParserListener) EnterStmtWordsMaybeJSON(ctx *StmtWordsMaybeJSONContext) {}

// ExitStmtWordsMaybeJSON is called when production stmtWordsMaybeJSON is exited.
func (s *BaseEarthParserListener) ExitStmtWordsMaybeJSON(ctx *StmtWordsMaybeJSONContext) {}

// EnterStmtWords is called when production stmtWords is entered.
func (s *BaseEarthParserListener) EnterStmtWords(ctx *StmtWordsContext) {}

// ExitStmtWords is called when production stmtWords is exited.
func (s *BaseEarthParserListener) ExitStmtWords(ctx *StmtWordsContext) {}

// EnterStmtWord is called when production stmtWord is entered.
func (s *BaseEarthParserListener) EnterStmtWord(ctx *StmtWordContext) {}

// ExitStmtWord is called when production stmtWord is exited.
func (s *BaseEarthParserListener) ExitStmtWord(ctx *StmtWordContext) {}
