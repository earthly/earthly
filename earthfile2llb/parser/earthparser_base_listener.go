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

// EnterAsName is called when production asName is entered.
func (s *BaseEarthParserListener) EnterAsName(ctx *AsNameContext) {}

// ExitAsName is called when production asName is exited.
func (s *BaseEarthParserListener) ExitAsName(ctx *AsNameContext) {}

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

// EnterSaveFrom is called when production saveFrom is entered.
func (s *BaseEarthParserListener) EnterSaveFrom(ctx *SaveFromContext) {}

// ExitSaveFrom is called when production saveFrom is exited.
func (s *BaseEarthParserListener) ExitSaveFrom(ctx *SaveFromContext) {}

// EnterSaveTo is called when production saveTo is entered.
func (s *BaseEarthParserListener) EnterSaveTo(ctx *SaveToContext) {}

// ExitSaveTo is called when production saveTo is exited.
func (s *BaseEarthParserListener) ExitSaveTo(ctx *SaveToContext) {}

// EnterSaveAsLocalTo is called when production saveAsLocalTo is entered.
func (s *BaseEarthParserListener) EnterSaveAsLocalTo(ctx *SaveAsLocalToContext) {}

// ExitSaveAsLocalTo is called when production saveAsLocalTo is exited.
func (s *BaseEarthParserListener) ExitSaveAsLocalTo(ctx *SaveAsLocalToContext) {}

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

// EnterWorkdirPath is called when production workdirPath is entered.
func (s *BaseEarthParserListener) EnterWorkdirPath(ctx *WorkdirPathContext) {}

// ExitWorkdirPath is called when production workdirPath is exited.
func (s *BaseEarthParserListener) ExitWorkdirPath(ctx *WorkdirPathContext) {}

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

// EnterGitCloneStmt is called when production gitCloneStmt is entered.
func (s *BaseEarthParserListener) EnterGitCloneStmt(ctx *GitCloneStmtContext) {}

// ExitGitCloneStmt is called when production gitCloneStmt is exited.
func (s *BaseEarthParserListener) ExitGitCloneStmt(ctx *GitCloneStmtContext) {}

// EnterGitURL is called when production gitURL is entered.
func (s *BaseEarthParserListener) EnterGitURL(ctx *GitURLContext) {}

// ExitGitURL is called when production gitURL is exited.
func (s *BaseEarthParserListener) ExitGitURL(ctx *GitURLContext) {}

// EnterGitCloneDest is called when production gitCloneDest is entered.
func (s *BaseEarthParserListener) EnterGitCloneDest(ctx *GitCloneDestContext) {}

// ExitGitCloneDest is called when production gitCloneDest is exited.
func (s *BaseEarthParserListener) ExitGitCloneDest(ctx *GitCloneDestContext) {}

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

// EnterRunArgs is called when production runArgs is entered.
func (s *BaseEarthParserListener) EnterRunArgs(ctx *RunArgsContext) {}

// ExitRunArgs is called when production runArgs is exited.
func (s *BaseEarthParserListener) ExitRunArgs(ctx *RunArgsContext) {}

// EnterRunArgsList is called when production runArgsList is entered.
func (s *BaseEarthParserListener) EnterRunArgsList(ctx *RunArgsListContext) {}

// ExitRunArgsList is called when production runArgsList is exited.
func (s *BaseEarthParserListener) ExitRunArgsList(ctx *RunArgsListContext) {}

// EnterRunArg is called when production runArg is entered.
func (s *BaseEarthParserListener) EnterRunArg(ctx *RunArgContext) {}

// ExitRunArg is called when production runArg is exited.
func (s *BaseEarthParserListener) ExitRunArg(ctx *RunArgContext) {}

// EnterEntrypointArgs is called when production entrypointArgs is entered.
func (s *BaseEarthParserListener) EnterEntrypointArgs(ctx *EntrypointArgsContext) {}

// ExitEntrypointArgs is called when production entrypointArgs is exited.
func (s *BaseEarthParserListener) ExitEntrypointArgs(ctx *EntrypointArgsContext) {}

// EnterEntrypointArgsList is called when production entrypointArgsList is entered.
func (s *BaseEarthParserListener) EnterEntrypointArgsList(ctx *EntrypointArgsListContext) {}

// ExitEntrypointArgsList is called when production entrypointArgsList is exited.
func (s *BaseEarthParserListener) ExitEntrypointArgsList(ctx *EntrypointArgsListContext) {}

// EnterEntrypointArg is called when production entrypointArg is entered.
func (s *BaseEarthParserListener) EnterEntrypointArg(ctx *EntrypointArgContext) {}

// ExitEntrypointArg is called when production entrypointArg is exited.
func (s *BaseEarthParserListener) ExitEntrypointArg(ctx *EntrypointArgContext) {}

// EnterFlags is called when production flags is entered.
func (s *BaseEarthParserListener) EnterFlags(ctx *FlagsContext) {}

// ExitFlags is called when production flags is exited.
func (s *BaseEarthParserListener) ExitFlags(ctx *FlagsContext) {}

// EnterFlag is called when production flag is entered.
func (s *BaseEarthParserListener) EnterFlag(ctx *FlagContext) {}

// ExitFlag is called when production flag is exited.
func (s *BaseEarthParserListener) ExitFlag(ctx *FlagContext) {}

// EnterFlagKey is called when production flagKey is entered.
func (s *BaseEarthParserListener) EnterFlagKey(ctx *FlagKeyContext) {}

// ExitFlagKey is called when production flagKey is exited.
func (s *BaseEarthParserListener) ExitFlagKey(ctx *FlagKeyContext) {}

// EnterFlagKeyValue is called when production flagKeyValue is entered.
func (s *BaseEarthParserListener) EnterFlagKeyValue(ctx *FlagKeyValueContext) {}

// ExitFlagKeyValue is called when production flagKeyValue is exited.
func (s *BaseEarthParserListener) ExitFlagKeyValue(ctx *FlagKeyValueContext) {}

// EnterStmtWords is called when production stmtWords is entered.
func (s *BaseEarthParserListener) EnterStmtWords(ctx *StmtWordsContext) {}

// ExitStmtWords is called when production stmtWords is exited.
func (s *BaseEarthParserListener) ExitStmtWords(ctx *StmtWordsContext) {}

// EnterStmtWord is called when production stmtWord is entered.
func (s *BaseEarthParserListener) EnterStmtWord(ctx *StmtWordContext) {}

// ExitStmtWord is called when production stmtWord is exited.
func (s *BaseEarthParserListener) ExitStmtWord(ctx *StmtWordContext) {}

// EnterEnvArgKey is called when production envArgKey is entered.
func (s *BaseEarthParserListener) EnterEnvArgKey(ctx *EnvArgKeyContext) {}

// ExitEnvArgKey is called when production envArgKey is exited.
func (s *BaseEarthParserListener) ExitEnvArgKey(ctx *EnvArgKeyContext) {}

// EnterEnvArgValue is called when production envArgValue is entered.
func (s *BaseEarthParserListener) EnterEnvArgValue(ctx *EnvArgValueContext) {}

// ExitEnvArgValue is called when production envArgValue is exited.
func (s *BaseEarthParserListener) ExitEnvArgValue(ctx *EnvArgValueContext) {}

// EnterImageName is called when production imageName is entered.
func (s *BaseEarthParserListener) EnterImageName(ctx *ImageNameContext) {}

// ExitImageName is called when production imageName is exited.
func (s *BaseEarthParserListener) ExitImageName(ctx *ImageNameContext) {}

// EnterSaveImageName is called when production saveImageName is entered.
func (s *BaseEarthParserListener) EnterSaveImageName(ctx *SaveImageNameContext) {}

// ExitSaveImageName is called when production saveImageName is exited.
func (s *BaseEarthParserListener) ExitSaveImageName(ctx *SaveImageNameContext) {}

// EnterTargetName is called when production targetName is entered.
func (s *BaseEarthParserListener) EnterTargetName(ctx *TargetNameContext) {}

// ExitTargetName is called when production targetName is exited.
func (s *BaseEarthParserListener) ExitTargetName(ctx *TargetNameContext) {}

// EnterFullTargetName is called when production fullTargetName is entered.
func (s *BaseEarthParserListener) EnterFullTargetName(ctx *FullTargetNameContext) {}

// ExitFullTargetName is called when production fullTargetName is exited.
func (s *BaseEarthParserListener) ExitFullTargetName(ctx *FullTargetNameContext) {}

// EnterArgsList is called when production argsList is entered.
func (s *BaseEarthParserListener) EnterArgsList(ctx *ArgsListContext) {}

// ExitArgsList is called when production argsList is exited.
func (s *BaseEarthParserListener) ExitArgsList(ctx *ArgsListContext) {}

// EnterArg is called when production arg is entered.
func (s *BaseEarthParserListener) EnterArg(ctx *ArgContext) {}

// ExitArg is called when production arg is exited.
func (s *BaseEarthParserListener) ExitArg(ctx *ArgContext) {}
