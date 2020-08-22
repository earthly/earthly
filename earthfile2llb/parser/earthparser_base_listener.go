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

// EnterFromDockerfileStmt is called when production fromDockerfileStmt is entered.
func (s *BaseEarthParserListener) EnterFromDockerfileStmt(ctx *FromDockerfileStmtContext) {}

// ExitFromDockerfileStmt is called when production fromDockerfileStmt is exited.
func (s *BaseEarthParserListener) ExitFromDockerfileStmt(ctx *FromDockerfileStmtContext) {}

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

// EnterDockerLoadStmt is called when production dockerLoadStmt is entered.
func (s *BaseEarthParserListener) EnterDockerLoadStmt(ctx *DockerLoadStmtContext) {}

// ExitDockerLoadStmt is called when production dockerLoadStmt is exited.
func (s *BaseEarthParserListener) ExitDockerLoadStmt(ctx *DockerLoadStmtContext) {}

// EnterDockerPullStmt is called when production dockerPullStmt is entered.
func (s *BaseEarthParserListener) EnterDockerPullStmt(ctx *DockerPullStmtContext) {}

// ExitDockerPullStmt is called when production dockerPullStmt is exited.
func (s *BaseEarthParserListener) ExitDockerPullStmt(ctx *DockerPullStmtContext) {}

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

// EnterWithDockerStmt is called when production withDockerStmt is entered.
func (s *BaseEarthParserListener) EnterWithDockerStmt(ctx *WithDockerStmtContext) {}

// ExitWithDockerStmt is called when production withDockerStmt is exited.
func (s *BaseEarthParserListener) ExitWithDockerStmt(ctx *WithDockerStmtContext) {}

// EnterEndStmt is called when production endStmt is entered.
func (s *BaseEarthParserListener) EnterEndStmt(ctx *EndStmtContext) {}

// ExitEndStmt is called when production endStmt is exited.
func (s *BaseEarthParserListener) ExitEndStmt(ctx *EndStmtContext) {}

// EnterGenericCommandStmt is called when production genericCommandStmt is entered.
func (s *BaseEarthParserListener) EnterGenericCommandStmt(ctx *GenericCommandStmtContext) {}

// ExitGenericCommandStmt is called when production genericCommandStmt is exited.
func (s *BaseEarthParserListener) ExitGenericCommandStmt(ctx *GenericCommandStmtContext) {}

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
