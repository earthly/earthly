// Code generated from ast/parser/EarthParser.g4 by ANTLR 4.8. DO NOT EDIT.

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

// EnterCommandStmt is called when production commandStmt is entered.
func (s *BaseEarthParserListener) EnterCommandStmt(ctx *CommandStmtContext) {}

// ExitCommandStmt is called when production commandStmt is exited.
func (s *BaseEarthParserListener) ExitCommandStmt(ctx *CommandStmtContext) {}

// EnterWithStmt is called when production withStmt is entered.
func (s *BaseEarthParserListener) EnterWithStmt(ctx *WithStmtContext) {}

// ExitWithStmt is called when production withStmt is exited.
func (s *BaseEarthParserListener) ExitWithStmt(ctx *WithStmtContext) {}

// EnterWithBlock is called when production withBlock is entered.
func (s *BaseEarthParserListener) EnterWithBlock(ctx *WithBlockContext) {}

// ExitWithBlock is called when production withBlock is exited.
func (s *BaseEarthParserListener) ExitWithBlock(ctx *WithBlockContext) {}

// EnterWithExpr is called when production withExpr is entered.
func (s *BaseEarthParserListener) EnterWithExpr(ctx *WithExprContext) {}

// ExitWithExpr is called when production withExpr is exited.
func (s *BaseEarthParserListener) ExitWithExpr(ctx *WithExprContext) {}

// EnterWithCommand is called when production withCommand is entered.
func (s *BaseEarthParserListener) EnterWithCommand(ctx *WithCommandContext) {}

// ExitWithCommand is called when production withCommand is exited.
func (s *BaseEarthParserListener) ExitWithCommand(ctx *WithCommandContext) {}

// EnterDockerCommand is called when production dockerCommand is entered.
func (s *BaseEarthParserListener) EnterDockerCommand(ctx *DockerCommandContext) {}

// ExitDockerCommand is called when production dockerCommand is exited.
func (s *BaseEarthParserListener) ExitDockerCommand(ctx *DockerCommandContext) {}

// EnterIfStmt is called when production ifStmt is entered.
func (s *BaseEarthParserListener) EnterIfStmt(ctx *IfStmtContext) {}

// ExitIfStmt is called when production ifStmt is exited.
func (s *BaseEarthParserListener) ExitIfStmt(ctx *IfStmtContext) {}

// EnterIfClause is called when production ifClause is entered.
func (s *BaseEarthParserListener) EnterIfClause(ctx *IfClauseContext) {}

// ExitIfClause is called when production ifClause is exited.
func (s *BaseEarthParserListener) ExitIfClause(ctx *IfClauseContext) {}

// EnterIfBlock is called when production ifBlock is entered.
func (s *BaseEarthParserListener) EnterIfBlock(ctx *IfBlockContext) {}

// ExitIfBlock is called when production ifBlock is exited.
func (s *BaseEarthParserListener) ExitIfBlock(ctx *IfBlockContext) {}

// EnterElseIfClause is called when production elseIfClause is entered.
func (s *BaseEarthParserListener) EnterElseIfClause(ctx *ElseIfClauseContext) {}

// ExitElseIfClause is called when production elseIfClause is exited.
func (s *BaseEarthParserListener) ExitElseIfClause(ctx *ElseIfClauseContext) {}

// EnterElseIfBlock is called when production elseIfBlock is entered.
func (s *BaseEarthParserListener) EnterElseIfBlock(ctx *ElseIfBlockContext) {}

// ExitElseIfBlock is called when production elseIfBlock is exited.
func (s *BaseEarthParserListener) ExitElseIfBlock(ctx *ElseIfBlockContext) {}

// EnterElseClause is called when production elseClause is entered.
func (s *BaseEarthParserListener) EnterElseClause(ctx *ElseClauseContext) {}

// ExitElseClause is called when production elseClause is exited.
func (s *BaseEarthParserListener) ExitElseClause(ctx *ElseClauseContext) {}

// EnterElseBlock is called when production elseBlock is entered.
func (s *BaseEarthParserListener) EnterElseBlock(ctx *ElseBlockContext) {}

// ExitElseBlock is called when production elseBlock is exited.
func (s *BaseEarthParserListener) ExitElseBlock(ctx *ElseBlockContext) {}

// EnterIfExpr is called when production ifExpr is entered.
func (s *BaseEarthParserListener) EnterIfExpr(ctx *IfExprContext) {}

// ExitIfExpr is called when production ifExpr is exited.
func (s *BaseEarthParserListener) ExitIfExpr(ctx *IfExprContext) {}

// EnterElseIfExpr is called when production elseIfExpr is entered.
func (s *BaseEarthParserListener) EnterElseIfExpr(ctx *ElseIfExprContext) {}

// ExitElseIfExpr is called when production elseIfExpr is exited.
func (s *BaseEarthParserListener) ExitElseIfExpr(ctx *ElseIfExprContext) {}

// EnterFromStmt is called when production fromStmt is entered.
func (s *BaseEarthParserListener) EnterFromStmt(ctx *FromStmtContext) {}

// ExitFromStmt is called when production fromStmt is exited.
func (s *BaseEarthParserListener) ExitFromStmt(ctx *FromStmtContext) {}

// EnterFromDockerfileStmt is called when production fromDockerfileStmt is entered.
func (s *BaseEarthParserListener) EnterFromDockerfileStmt(ctx *FromDockerfileStmtContext) {}

// ExitFromDockerfileStmt is called when production fromDockerfileStmt is exited.
func (s *BaseEarthParserListener) ExitFromDockerfileStmt(ctx *FromDockerfileStmtContext) {}

// EnterLocallyStmt is called when production locallyStmt is entered.
func (s *BaseEarthParserListener) EnterLocallyStmt(ctx *LocallyStmtContext) {}

// ExitLocallyStmt is called when production locallyStmt is exited.
func (s *BaseEarthParserListener) ExitLocallyStmt(ctx *LocallyStmtContext) {}

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

// EnterUserStmt is called when production userStmt is entered.
func (s *BaseEarthParserListener) EnterUserStmt(ctx *UserStmtContext) {}

// ExitUserStmt is called when production userStmt is exited.
func (s *BaseEarthParserListener) ExitUserStmt(ctx *UserStmtContext) {}

// EnterCmdStmt is called when production cmdStmt is entered.
func (s *BaseEarthParserListener) EnterCmdStmt(ctx *CmdStmtContext) {}

// ExitCmdStmt is called when production cmdStmt is exited.
func (s *BaseEarthParserListener) ExitCmdStmt(ctx *CmdStmtContext) {}

// EnterEntrypointStmt is called when production entrypointStmt is entered.
func (s *BaseEarthParserListener) EnterEntrypointStmt(ctx *EntrypointStmtContext) {}

// ExitEntrypointStmt is called when production entrypointStmt is exited.
func (s *BaseEarthParserListener) ExitEntrypointStmt(ctx *EntrypointStmtContext) {}

// EnterExposeStmt is called when production exposeStmt is entered.
func (s *BaseEarthParserListener) EnterExposeStmt(ctx *ExposeStmtContext) {}

// ExitExposeStmt is called when production exposeStmt is exited.
func (s *BaseEarthParserListener) ExitExposeStmt(ctx *ExposeStmtContext) {}

// EnterVolumeStmt is called when production volumeStmt is entered.
func (s *BaseEarthParserListener) EnterVolumeStmt(ctx *VolumeStmtContext) {}

// ExitVolumeStmt is called when production volumeStmt is exited.
func (s *BaseEarthParserListener) ExitVolumeStmt(ctx *VolumeStmtContext) {}

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

// EnterLabelStmt is called when production labelStmt is entered.
func (s *BaseEarthParserListener) EnterLabelStmt(ctx *LabelStmtContext) {}

// ExitLabelStmt is called when production labelStmt is exited.
func (s *BaseEarthParserListener) ExitLabelStmt(ctx *LabelStmtContext) {}

// EnterLabelKey is called when production labelKey is entered.
func (s *BaseEarthParserListener) EnterLabelKey(ctx *LabelKeyContext) {}

// ExitLabelKey is called when production labelKey is exited.
func (s *BaseEarthParserListener) ExitLabelKey(ctx *LabelKeyContext) {}

// EnterLabelValue is called when production labelValue is entered.
func (s *BaseEarthParserListener) EnterLabelValue(ctx *LabelValueContext) {}

// ExitLabelValue is called when production labelValue is exited.
func (s *BaseEarthParserListener) ExitLabelValue(ctx *LabelValueContext) {}

// EnterGitCloneStmt is called when production gitCloneStmt is entered.
func (s *BaseEarthParserListener) EnterGitCloneStmt(ctx *GitCloneStmtContext) {}

// ExitGitCloneStmt is called when production gitCloneStmt is exited.
func (s *BaseEarthParserListener) ExitGitCloneStmt(ctx *GitCloneStmtContext) {}

// EnterAddStmt is called when production addStmt is entered.
func (s *BaseEarthParserListener) EnterAddStmt(ctx *AddStmtContext) {}

// ExitAddStmt is called when production addStmt is exited.
func (s *BaseEarthParserListener) ExitAddStmt(ctx *AddStmtContext) {}

// EnterStopsignalStmt is called when production stopsignalStmt is entered.
func (s *BaseEarthParserListener) EnterStopsignalStmt(ctx *StopsignalStmtContext) {}

// ExitStopsignalStmt is called when production stopsignalStmt is exited.
func (s *BaseEarthParserListener) ExitStopsignalStmt(ctx *StopsignalStmtContext) {}

// EnterOnbuildStmt is called when production onbuildStmt is entered.
func (s *BaseEarthParserListener) EnterOnbuildStmt(ctx *OnbuildStmtContext) {}

// ExitOnbuildStmt is called when production onbuildStmt is exited.
func (s *BaseEarthParserListener) ExitOnbuildStmt(ctx *OnbuildStmtContext) {}

// EnterHealthcheckStmt is called when production healthcheckStmt is entered.
func (s *BaseEarthParserListener) EnterHealthcheckStmt(ctx *HealthcheckStmtContext) {}

// ExitHealthcheckStmt is called when production healthcheckStmt is exited.
func (s *BaseEarthParserListener) ExitHealthcheckStmt(ctx *HealthcheckStmtContext) {}

// EnterShellStmt is called when production shellStmt is entered.
func (s *BaseEarthParserListener) EnterShellStmt(ctx *ShellStmtContext) {}

// ExitShellStmt is called when production shellStmt is exited.
func (s *BaseEarthParserListener) ExitShellStmt(ctx *ShellStmtContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseEarthParserListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseEarthParserListener) ExitExpr(ctx *ExprContext) {}

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
