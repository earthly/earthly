// Code generated from ast/parser/EarthParser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // EarthParser

import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

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

// EnterTargetOrUserCommand is called when production targetOrUserCommand is entered.
func (s *BaseEarthParserListener) EnterTargetOrUserCommand(ctx *TargetOrUserCommandContext) {}

// ExitTargetOrUserCommand is called when production targetOrUserCommand is exited.
func (s *BaseEarthParserListener) ExitTargetOrUserCommand(ctx *TargetOrUserCommandContext) {}

// EnterTarget is called when production target is entered.
func (s *BaseEarthParserListener) EnterTarget(ctx *TargetContext) {}

// ExitTarget is called when production target is exited.
func (s *BaseEarthParserListener) ExitTarget(ctx *TargetContext) {}

// EnterTargetHeader is called when production targetHeader is entered.
func (s *BaseEarthParserListener) EnterTargetHeader(ctx *TargetHeaderContext) {}

// ExitTargetHeader is called when production targetHeader is exited.
func (s *BaseEarthParserListener) ExitTargetHeader(ctx *TargetHeaderContext) {}

// EnterUserCommand is called when production userCommand is entered.
func (s *BaseEarthParserListener) EnterUserCommand(ctx *UserCommandContext) {}

// ExitUserCommand is called when production userCommand is exited.
func (s *BaseEarthParserListener) ExitUserCommand(ctx *UserCommandContext) {}

// EnterUserCommandHeader is called when production userCommandHeader is entered.
func (s *BaseEarthParserListener) EnterUserCommandHeader(ctx *UserCommandHeaderContext) {}

// ExitUserCommandHeader is called when production userCommandHeader is exited.
func (s *BaseEarthParserListener) ExitUserCommandHeader(ctx *UserCommandHeaderContext) {}

// EnterFunction is called when production function is entered.
func (s *BaseEarthParserListener) EnterFunction(ctx *FunctionContext) {}

// ExitFunction is called when production function is exited.
func (s *BaseEarthParserListener) ExitFunction(ctx *FunctionContext) {}

// EnterFunctionHeader is called when production functionHeader is entered.
func (s *BaseEarthParserListener) EnterFunctionHeader(ctx *FunctionHeaderContext) {}

// ExitFunctionHeader is called when production functionHeader is exited.
func (s *BaseEarthParserListener) ExitFunctionHeader(ctx *FunctionHeaderContext) {}

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

// EnterVersion is called when production version is entered.
func (s *BaseEarthParserListener) EnterVersion(ctx *VersionContext) {}

// ExitVersion is called when production version is exited.
func (s *BaseEarthParserListener) ExitVersion(ctx *VersionContext) {}

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

// EnterTryStmt is called when production tryStmt is entered.
func (s *BaseEarthParserListener) EnterTryStmt(ctx *TryStmtContext) {}

// ExitTryStmt is called when production tryStmt is exited.
func (s *BaseEarthParserListener) ExitTryStmt(ctx *TryStmtContext) {}

// EnterTryClause is called when production tryClause is entered.
func (s *BaseEarthParserListener) EnterTryClause(ctx *TryClauseContext) {}

// ExitTryClause is called when production tryClause is exited.
func (s *BaseEarthParserListener) ExitTryClause(ctx *TryClauseContext) {}

// EnterTryBlock is called when production tryBlock is entered.
func (s *BaseEarthParserListener) EnterTryBlock(ctx *TryBlockContext) {}

// ExitTryBlock is called when production tryBlock is exited.
func (s *BaseEarthParserListener) ExitTryBlock(ctx *TryBlockContext) {}

// EnterCatchClause is called when production catchClause is entered.
func (s *BaseEarthParserListener) EnterCatchClause(ctx *CatchClauseContext) {}

// ExitCatchClause is called when production catchClause is exited.
func (s *BaseEarthParserListener) ExitCatchClause(ctx *CatchClauseContext) {}

// EnterCatchBlock is called when production catchBlock is entered.
func (s *BaseEarthParserListener) EnterCatchBlock(ctx *CatchBlockContext) {}

// ExitCatchBlock is called when production catchBlock is exited.
func (s *BaseEarthParserListener) ExitCatchBlock(ctx *CatchBlockContext) {}

// EnterFinallyClause is called when production finallyClause is entered.
func (s *BaseEarthParserListener) EnterFinallyClause(ctx *FinallyClauseContext) {}

// ExitFinallyClause is called when production finallyClause is exited.
func (s *BaseEarthParserListener) ExitFinallyClause(ctx *FinallyClauseContext) {}

// EnterFinallyBlock is called when production finallyBlock is entered.
func (s *BaseEarthParserListener) EnterFinallyBlock(ctx *FinallyBlockContext) {}

// ExitFinallyBlock is called when production finallyBlock is exited.
func (s *BaseEarthParserListener) ExitFinallyBlock(ctx *FinallyBlockContext) {}

// EnterForStmt is called when production forStmt is entered.
func (s *BaseEarthParserListener) EnterForStmt(ctx *ForStmtContext) {}

// ExitForStmt is called when production forStmt is exited.
func (s *BaseEarthParserListener) ExitForStmt(ctx *ForStmtContext) {}

// EnterForClause is called when production forClause is entered.
func (s *BaseEarthParserListener) EnterForClause(ctx *ForClauseContext) {}

// ExitForClause is called when production forClause is exited.
func (s *BaseEarthParserListener) ExitForClause(ctx *ForClauseContext) {}

// EnterForBlock is called when production forBlock is entered.
func (s *BaseEarthParserListener) EnterForBlock(ctx *ForBlockContext) {}

// ExitForBlock is called when production forBlock is exited.
func (s *BaseEarthParserListener) ExitForBlock(ctx *ForBlockContext) {}

// EnterForExpr is called when production forExpr is entered.
func (s *BaseEarthParserListener) EnterForExpr(ctx *ForExprContext) {}

// ExitForExpr is called when production forExpr is exited.
func (s *BaseEarthParserListener) ExitForExpr(ctx *ForExprContext) {}

// EnterWaitStmt is called when production waitStmt is entered.
func (s *BaseEarthParserListener) EnterWaitStmt(ctx *WaitStmtContext) {}

// ExitWaitStmt is called when production waitStmt is exited.
func (s *BaseEarthParserListener) ExitWaitStmt(ctx *WaitStmtContext) {}

// EnterWaitClause is called when production waitClause is entered.
func (s *BaseEarthParserListener) EnterWaitClause(ctx *WaitClauseContext) {}

// ExitWaitClause is called when production waitClause is exited.
func (s *BaseEarthParserListener) ExitWaitClause(ctx *WaitClauseContext) {}

// EnterWaitBlock is called when production waitBlock is entered.
func (s *BaseEarthParserListener) EnterWaitBlock(ctx *WaitBlockContext) {}

// ExitWaitBlock is called when production waitBlock is exited.
func (s *BaseEarthParserListener) ExitWaitBlock(ctx *WaitBlockContext) {}

// EnterWaitExpr is called when production waitExpr is entered.
func (s *BaseEarthParserListener) EnterWaitExpr(ctx *WaitExprContext) {}

// ExitWaitExpr is called when production waitExpr is exited.
func (s *BaseEarthParserListener) ExitWaitExpr(ctx *WaitExprContext) {}

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

// EnterSetStmt is called when production setStmt is entered.
func (s *BaseEarthParserListener) EnterSetStmt(ctx *SetStmtContext) {}

// ExitSetStmt is called when production setStmt is exited.
func (s *BaseEarthParserListener) ExitSetStmt(ctx *SetStmtContext) {}

// EnterLetStmt is called when production letStmt is entered.
func (s *BaseEarthParserListener) EnterLetStmt(ctx *LetStmtContext) {}

// ExitLetStmt is called when production letStmt is exited.
func (s *BaseEarthParserListener) ExitLetStmt(ctx *LetStmtContext) {}

// EnterOptionalFlag is called when production optionalFlag is entered.
func (s *BaseEarthParserListener) EnterOptionalFlag(ctx *OptionalFlagContext) {}

// ExitOptionalFlag is called when production optionalFlag is exited.
func (s *BaseEarthParserListener) ExitOptionalFlag(ctx *OptionalFlagContext) {}

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

// EnterUserCommandStmt is called when production userCommandStmt is entered.
func (s *BaseEarthParserListener) EnterUserCommandStmt(ctx *UserCommandStmtContext) {}

// ExitUserCommandStmt is called when production userCommandStmt is exited.
func (s *BaseEarthParserListener) ExitUserCommandStmt(ctx *UserCommandStmtContext) {}

// EnterFunctionStmt is called when production functionStmt is entered.
func (s *BaseEarthParserListener) EnterFunctionStmt(ctx *FunctionStmtContext) {}

// ExitFunctionStmt is called when production functionStmt is exited.
func (s *BaseEarthParserListener) ExitFunctionStmt(ctx *FunctionStmtContext) {}

// EnterDoStmt is called when production doStmt is entered.
func (s *BaseEarthParserListener) EnterDoStmt(ctx *DoStmtContext) {}

// ExitDoStmt is called when production doStmt is exited.
func (s *BaseEarthParserListener) ExitDoStmt(ctx *DoStmtContext) {}

// EnterImportStmt is called when production importStmt is entered.
func (s *BaseEarthParserListener) EnterImportStmt(ctx *ImportStmtContext) {}

// ExitImportStmt is called when production importStmt is exited.
func (s *BaseEarthParserListener) ExitImportStmt(ctx *ImportStmtContext) {}

// EnterCacheStmt is called when production cacheStmt is entered.
func (s *BaseEarthParserListener) EnterCacheStmt(ctx *CacheStmtContext) {}

// ExitCacheStmt is called when production cacheStmt is exited.
func (s *BaseEarthParserListener) ExitCacheStmt(ctx *CacheStmtContext) {}

// EnterHostStmt is called when production hostStmt is entered.
func (s *BaseEarthParserListener) EnterHostStmt(ctx *HostStmtContext) {}

// ExitHostStmt is called when production hostStmt is exited.
func (s *BaseEarthParserListener) ExitHostStmt(ctx *HostStmtContext) {}

// EnterProjectStmt is called when production projectStmt is entered.
func (s *BaseEarthParserListener) EnterProjectStmt(ctx *ProjectStmtContext) {}

// ExitProjectStmt is called when production projectStmt is exited.
func (s *BaseEarthParserListener) ExitProjectStmt(ctx *ProjectStmtContext) {}

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
