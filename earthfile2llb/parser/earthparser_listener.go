// Code generated from earthfile2llb/parser/EarthParser.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // EarthParser

import "github.com/antlr/antlr4/runtime/Go/antlr"

// EarthParserListener is a complete listener for a parse tree produced by EarthParser.
type EarthParserListener interface {
	antlr.ParseTreeListener

	// EnterEarthFile is called when entering the earthFile production.
	EnterEarthFile(c *EarthFileContext)

	// EnterTargets is called when entering the targets production.
	EnterTargets(c *TargetsContext)

	// EnterTarget is called when entering the target production.
	EnterTarget(c *TargetContext)

	// EnterTargetHeader is called when entering the targetHeader production.
	EnterTargetHeader(c *TargetHeaderContext)

	// EnterStmts is called when entering the stmts production.
	EnterStmts(c *StmtsContext)

	// EnterStmt is called when entering the stmt production.
	EnterStmt(c *StmtContext)

	// EnterFromStmt is called when entering the fromStmt production.
	EnterFromStmt(c *FromStmtContext)

	// EnterCopyStmt is called when entering the copyStmt production.
	EnterCopyStmt(c *CopyStmtContext)

	// EnterSaveStmt is called when entering the saveStmt production.
	EnterSaveStmt(c *SaveStmtContext)

	// EnterSaveImage is called when entering the saveImage production.
	EnterSaveImage(c *SaveImageContext)

	// EnterSaveArtifact is called when entering the saveArtifact production.
	EnterSaveArtifact(c *SaveArtifactContext)

	// EnterSaveFrom is called when entering the saveFrom production.
	EnterSaveFrom(c *SaveFromContext)

	// EnterSaveTo is called when entering the saveTo production.
	EnterSaveTo(c *SaveToContext)

	// EnterSaveAsLocalTo is called when entering the saveAsLocalTo production.
	EnterSaveAsLocalTo(c *SaveAsLocalToContext)

	// EnterRunStmt is called when entering the runStmt production.
	EnterRunStmt(c *RunStmtContext)

	// EnterBuildStmt is called when entering the buildStmt production.
	EnterBuildStmt(c *BuildStmtContext)

	// EnterWorkdirStmt is called when entering the workdirStmt production.
	EnterWorkdirStmt(c *WorkdirStmtContext)

	// EnterWorkdirPath is called when entering the workdirPath production.
	EnterWorkdirPath(c *WorkdirPathContext)

	// EnterEntrypointStmt is called when entering the entrypointStmt production.
	EnterEntrypointStmt(c *EntrypointStmtContext)

	// EnterEnvStmt is called when entering the envStmt production.
	EnterEnvStmt(c *EnvStmtContext)

	// EnterArgStmt is called when entering the argStmt production.
	EnterArgStmt(c *ArgStmtContext)

	// EnterGitCloneStmt is called when entering the gitCloneStmt production.
	EnterGitCloneStmt(c *GitCloneStmtContext)

	// EnterGitURL is called when entering the gitURL production.
	EnterGitURL(c *GitURLContext)

	// EnterGitCloneDest is called when entering the gitCloneDest production.
	EnterGitCloneDest(c *GitCloneDestContext)

	// EnterDockerLoadStmt is called when entering the dockerLoadStmt production.
	EnterDockerLoadStmt(c *DockerLoadStmtContext)

	// EnterDockerPullStmt is called when entering the dockerPullStmt production.
	EnterDockerPullStmt(c *DockerPullStmtContext)

	// EnterGenericCommand is called when entering the genericCommand production.
	EnterGenericCommand(c *GenericCommandContext)

	// EnterCommandName is called when entering the commandName production.
	EnterCommandName(c *CommandNameContext)

	// EnterFlags is called when entering the flags production.
	EnterFlags(c *FlagsContext)

	// EnterFlag is called when entering the flag production.
	EnterFlag(c *FlagContext)

	// EnterFlagKey is called when entering the flagKey production.
	EnterFlagKey(c *FlagKeyContext)

	// EnterFlagKeyValue is called when entering the flagKeyValue production.
	EnterFlagKeyValue(c *FlagKeyValueContext)

	// EnterStmtWords is called when entering the stmtWords production.
	EnterStmtWords(c *StmtWordsContext)

	// EnterStmtWordsList is called when entering the stmtWordsList production.
	EnterStmtWordsList(c *StmtWordsListContext)

	// EnterStmtWord is called when entering the stmtWord production.
	EnterStmtWord(c *StmtWordContext)

	// EnterEnvArgKey is called when entering the envArgKey production.
	EnterEnvArgKey(c *EnvArgKeyContext)

	// EnterEnvArgValue is called when entering the envArgValue production.
	EnterEnvArgValue(c *EnvArgValueContext)

	// EnterImageName is called when entering the imageName production.
	EnterImageName(c *ImageNameContext)

	// EnterSaveImageName is called when entering the saveImageName production.
	EnterSaveImageName(c *SaveImageNameContext)

	// EnterTargetName is called when entering the targetName production.
	EnterTargetName(c *TargetNameContext)

	// EnterFullTargetName is called when entering the fullTargetName production.
	EnterFullTargetName(c *FullTargetNameContext)

	// EnterArgsList is called when entering the argsList production.
	EnterArgsList(c *ArgsListContext)

	// EnterArg is called when entering the arg production.
	EnterArg(c *ArgContext)

	// ExitEarthFile is called when exiting the earthFile production.
	ExitEarthFile(c *EarthFileContext)

	// ExitTargets is called when exiting the targets production.
	ExitTargets(c *TargetsContext)

	// ExitTarget is called when exiting the target production.
	ExitTarget(c *TargetContext)

	// ExitTargetHeader is called when exiting the targetHeader production.
	ExitTargetHeader(c *TargetHeaderContext)

	// ExitStmts is called when exiting the stmts production.
	ExitStmts(c *StmtsContext)

	// ExitStmt is called when exiting the stmt production.
	ExitStmt(c *StmtContext)

	// ExitFromStmt is called when exiting the fromStmt production.
	ExitFromStmt(c *FromStmtContext)

	// ExitCopyStmt is called when exiting the copyStmt production.
	ExitCopyStmt(c *CopyStmtContext)

	// ExitSaveStmt is called when exiting the saveStmt production.
	ExitSaveStmt(c *SaveStmtContext)

	// ExitSaveImage is called when exiting the saveImage production.
	ExitSaveImage(c *SaveImageContext)

	// ExitSaveArtifact is called when exiting the saveArtifact production.
	ExitSaveArtifact(c *SaveArtifactContext)

	// ExitSaveFrom is called when exiting the saveFrom production.
	ExitSaveFrom(c *SaveFromContext)

	// ExitSaveTo is called when exiting the saveTo production.
	ExitSaveTo(c *SaveToContext)

	// ExitSaveAsLocalTo is called when exiting the saveAsLocalTo production.
	ExitSaveAsLocalTo(c *SaveAsLocalToContext)

	// ExitRunStmt is called when exiting the runStmt production.
	ExitRunStmt(c *RunStmtContext)

	// ExitBuildStmt is called when exiting the buildStmt production.
	ExitBuildStmt(c *BuildStmtContext)

	// ExitWorkdirStmt is called when exiting the workdirStmt production.
	ExitWorkdirStmt(c *WorkdirStmtContext)

	// ExitWorkdirPath is called when exiting the workdirPath production.
	ExitWorkdirPath(c *WorkdirPathContext)

	// ExitEntrypointStmt is called when exiting the entrypointStmt production.
	ExitEntrypointStmt(c *EntrypointStmtContext)

	// ExitEnvStmt is called when exiting the envStmt production.
	ExitEnvStmt(c *EnvStmtContext)

	// ExitArgStmt is called when exiting the argStmt production.
	ExitArgStmt(c *ArgStmtContext)

	// ExitGitCloneStmt is called when exiting the gitCloneStmt production.
	ExitGitCloneStmt(c *GitCloneStmtContext)

	// ExitGitURL is called when exiting the gitURL production.
	ExitGitURL(c *GitURLContext)

	// ExitGitCloneDest is called when exiting the gitCloneDest production.
	ExitGitCloneDest(c *GitCloneDestContext)

	// ExitDockerLoadStmt is called when exiting the dockerLoadStmt production.
	ExitDockerLoadStmt(c *DockerLoadStmtContext)

	// ExitDockerPullStmt is called when exiting the dockerPullStmt production.
	ExitDockerPullStmt(c *DockerPullStmtContext)

	// ExitGenericCommand is called when exiting the genericCommand production.
	ExitGenericCommand(c *GenericCommandContext)

	// ExitCommandName is called when exiting the commandName production.
	ExitCommandName(c *CommandNameContext)

	// ExitFlags is called when exiting the flags production.
	ExitFlags(c *FlagsContext)

	// ExitFlag is called when exiting the flag production.
	ExitFlag(c *FlagContext)

	// ExitFlagKey is called when exiting the flagKey production.
	ExitFlagKey(c *FlagKeyContext)

	// ExitFlagKeyValue is called when exiting the flagKeyValue production.
	ExitFlagKeyValue(c *FlagKeyValueContext)

	// ExitStmtWords is called when exiting the stmtWords production.
	ExitStmtWords(c *StmtWordsContext)

	// ExitStmtWordsList is called when exiting the stmtWordsList production.
	ExitStmtWordsList(c *StmtWordsListContext)

	// ExitStmtWord is called when exiting the stmtWord production.
	ExitStmtWord(c *StmtWordContext)

	// ExitEnvArgKey is called when exiting the envArgKey production.
	ExitEnvArgKey(c *EnvArgKeyContext)

	// ExitEnvArgValue is called when exiting the envArgValue production.
	ExitEnvArgValue(c *EnvArgValueContext)

	// ExitImageName is called when exiting the imageName production.
	ExitImageName(c *ImageNameContext)

	// ExitSaveImageName is called when exiting the saveImageName production.
	ExitSaveImageName(c *SaveImageNameContext)

	// ExitTargetName is called when exiting the targetName production.
	ExitTargetName(c *TargetNameContext)

	// ExitFullTargetName is called when exiting the fullTargetName production.
	ExitFullTargetName(c *FullTargetNameContext)

	// ExitArgsList is called when exiting the argsList production.
	ExitArgsList(c *ArgsListContext)

	// ExitArg is called when exiting the arg production.
	ExitArg(c *ArgContext)
}
