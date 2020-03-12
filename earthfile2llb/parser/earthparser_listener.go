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

	// EnterCopyArgsFrom is called when entering the copyArgsFrom production.
	EnterCopyArgsFrom(c *CopyArgsFromContext)

	// EnterCopyArgsArtifact is called when entering the copyArgsArtifact production.
	EnterCopyArgsArtifact(c *CopyArgsArtifactContext)

	// EnterCopyArgsClassical is called when entering the copyArgsClassical production.
	EnterCopyArgsClassical(c *CopyArgsClassicalContext)

	// EnterCopySrcs is called when entering the copySrcs production.
	EnterCopySrcs(c *CopySrcsContext)

	// EnterSaveStmt is called when entering the saveStmt production.
	EnterSaveStmt(c *SaveStmtContext)

	// EnterSaveImage is called when entering the saveImage production.
	EnterSaveImage(c *SaveImageContext)

	// EnterSaveArtifact is called when entering the saveArtifact production.
	EnterSaveArtifact(c *SaveArtifactContext)

	// EnterRunStmt is called when entering the runStmt production.
	EnterRunStmt(c *RunStmtContext)

	// EnterBuildStmt is called when entering the buildStmt production.
	EnterBuildStmt(c *BuildStmtContext)

	// EnterWorkdirStmt is called when entering the workdirStmt production.
	EnterWorkdirStmt(c *WorkdirStmtContext)

	// EnterEntrypointStmt is called when entering the entrypointStmt production.
	EnterEntrypointStmt(c *EntrypointStmtContext)

	// EnterEnvStmt is called when entering the envStmt production.
	EnterEnvStmt(c *EnvStmtContext)

	// EnterArgStmt is called when entering the argStmt production.
	EnterArgStmt(c *ArgStmtContext)

	// EnterGitCloneStmt is called when entering the gitCloneStmt production.
	EnterGitCloneStmt(c *GitCloneStmtContext)

	// EnterDockerLoadStmt is called when entering the dockerLoadStmt production.
	EnterDockerLoadStmt(c *DockerLoadStmtContext)

	// EnterDockerPullStmt is called when entering the dockerPullStmt production.
	EnterDockerPullStmt(c *DockerPullStmtContext)

	// EnterGenericCommand is called when entering the genericCommand production.
	EnterGenericCommand(c *GenericCommandContext)

	// EnterCommandName is called when entering the commandName production.
	EnterCommandName(c *CommandNameContext)

	// EnterRunArgs is called when entering the runArgs production.
	EnterRunArgs(c *RunArgsContext)

	// EnterRunArgsList is called when entering the runArgsList production.
	EnterRunArgsList(c *RunArgsListContext)

	// EnterRunArg is called when entering the runArg production.
	EnterRunArg(c *RunArgContext)

	// EnterEntrypointArgs is called when entering the entrypointArgs production.
	EnterEntrypointArgs(c *EntrypointArgsContext)

	// EnterEntrypointArgsList is called when entering the entrypointArgsList production.
	EnterEntrypointArgsList(c *EntrypointArgsListContext)

	// EnterEntrypointArg is called when entering the entrypointArg production.
	EnterEntrypointArg(c *EntrypointArgContext)

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

	// EnterStmtWord is called when entering the stmtWord production.
	EnterStmtWord(c *StmtWordContext)

	// EnterEnvArgKey is called when entering the envArgKey production.
	EnterEnvArgKey(c *EnvArgKeyContext)

	// EnterEnvArgValue is called when entering the envArgValue production.
	EnterEnvArgValue(c *EnvArgValueContext)

	// EnterCopySrc is called when entering the copySrc production.
	EnterCopySrc(c *CopySrcContext)

	// EnterCopyDest is called when entering the copyDest production.
	EnterCopyDest(c *CopyDestContext)

	// EnterAsName is called when entering the asName production.
	EnterAsName(c *AsNameContext)

	// EnterImageName is called when entering the imageName production.
	EnterImageName(c *ImageNameContext)

	// EnterSaveImageName is called when entering the saveImageName production.
	EnterSaveImageName(c *SaveImageNameContext)

	// EnterTargetName is called when entering the targetName production.
	EnterTargetName(c *TargetNameContext)

	// EnterFullTargetName is called when entering the fullTargetName production.
	EnterFullTargetName(c *FullTargetNameContext)

	// EnterArtifactName is called when entering the artifactName production.
	EnterArtifactName(c *ArtifactNameContext)

	// EnterSaveFrom is called when entering the saveFrom production.
	EnterSaveFrom(c *SaveFromContext)

	// EnterSaveTo is called when entering the saveTo production.
	EnterSaveTo(c *SaveToContext)

	// EnterSaveAsLocalTo is called when entering the saveAsLocalTo production.
	EnterSaveAsLocalTo(c *SaveAsLocalToContext)

	// EnterWorkdirPath is called when entering the workdirPath production.
	EnterWorkdirPath(c *WorkdirPathContext)

	// EnterGitURL is called when entering the gitURL production.
	EnterGitURL(c *GitURLContext)

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

	// ExitCopyArgsFrom is called when exiting the copyArgsFrom production.
	ExitCopyArgsFrom(c *CopyArgsFromContext)

	// ExitCopyArgsArtifact is called when exiting the copyArgsArtifact production.
	ExitCopyArgsArtifact(c *CopyArgsArtifactContext)

	// ExitCopyArgsClassical is called when exiting the copyArgsClassical production.
	ExitCopyArgsClassical(c *CopyArgsClassicalContext)

	// ExitCopySrcs is called when exiting the copySrcs production.
	ExitCopySrcs(c *CopySrcsContext)

	// ExitSaveStmt is called when exiting the saveStmt production.
	ExitSaveStmt(c *SaveStmtContext)

	// ExitSaveImage is called when exiting the saveImage production.
	ExitSaveImage(c *SaveImageContext)

	// ExitSaveArtifact is called when exiting the saveArtifact production.
	ExitSaveArtifact(c *SaveArtifactContext)

	// ExitRunStmt is called when exiting the runStmt production.
	ExitRunStmt(c *RunStmtContext)

	// ExitBuildStmt is called when exiting the buildStmt production.
	ExitBuildStmt(c *BuildStmtContext)

	// ExitWorkdirStmt is called when exiting the workdirStmt production.
	ExitWorkdirStmt(c *WorkdirStmtContext)

	// ExitEntrypointStmt is called when exiting the entrypointStmt production.
	ExitEntrypointStmt(c *EntrypointStmtContext)

	// ExitEnvStmt is called when exiting the envStmt production.
	ExitEnvStmt(c *EnvStmtContext)

	// ExitArgStmt is called when exiting the argStmt production.
	ExitArgStmt(c *ArgStmtContext)

	// ExitGitCloneStmt is called when exiting the gitCloneStmt production.
	ExitGitCloneStmt(c *GitCloneStmtContext)

	// ExitDockerLoadStmt is called when exiting the dockerLoadStmt production.
	ExitDockerLoadStmt(c *DockerLoadStmtContext)

	// ExitDockerPullStmt is called when exiting the dockerPullStmt production.
	ExitDockerPullStmt(c *DockerPullStmtContext)

	// ExitGenericCommand is called when exiting the genericCommand production.
	ExitGenericCommand(c *GenericCommandContext)

	// ExitCommandName is called when exiting the commandName production.
	ExitCommandName(c *CommandNameContext)

	// ExitRunArgs is called when exiting the runArgs production.
	ExitRunArgs(c *RunArgsContext)

	// ExitRunArgsList is called when exiting the runArgsList production.
	ExitRunArgsList(c *RunArgsListContext)

	// ExitRunArg is called when exiting the runArg production.
	ExitRunArg(c *RunArgContext)

	// ExitEntrypointArgs is called when exiting the entrypointArgs production.
	ExitEntrypointArgs(c *EntrypointArgsContext)

	// ExitEntrypointArgsList is called when exiting the entrypointArgsList production.
	ExitEntrypointArgsList(c *EntrypointArgsListContext)

	// ExitEntrypointArg is called when exiting the entrypointArg production.
	ExitEntrypointArg(c *EntrypointArgContext)

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

	// ExitStmtWord is called when exiting the stmtWord production.
	ExitStmtWord(c *StmtWordContext)

	// ExitEnvArgKey is called when exiting the envArgKey production.
	ExitEnvArgKey(c *EnvArgKeyContext)

	// ExitEnvArgValue is called when exiting the envArgValue production.
	ExitEnvArgValue(c *EnvArgValueContext)

	// ExitCopySrc is called when exiting the copySrc production.
	ExitCopySrc(c *CopySrcContext)

	// ExitCopyDest is called when exiting the copyDest production.
	ExitCopyDest(c *CopyDestContext)

	// ExitAsName is called when exiting the asName production.
	ExitAsName(c *AsNameContext)

	// ExitImageName is called when exiting the imageName production.
	ExitImageName(c *ImageNameContext)

	// ExitSaveImageName is called when exiting the saveImageName production.
	ExitSaveImageName(c *SaveImageNameContext)

	// ExitTargetName is called when exiting the targetName production.
	ExitTargetName(c *TargetNameContext)

	// ExitFullTargetName is called when exiting the fullTargetName production.
	ExitFullTargetName(c *FullTargetNameContext)

	// ExitArtifactName is called when exiting the artifactName production.
	ExitArtifactName(c *ArtifactNameContext)

	// ExitSaveFrom is called when exiting the saveFrom production.
	ExitSaveFrom(c *SaveFromContext)

	// ExitSaveTo is called when exiting the saveTo production.
	ExitSaveTo(c *SaveToContext)

	// ExitSaveAsLocalTo is called when exiting the saveAsLocalTo production.
	ExitSaveAsLocalTo(c *SaveAsLocalToContext)

	// ExitWorkdirPath is called when exiting the workdirPath production.
	ExitWorkdirPath(c *WorkdirPathContext)

	// ExitGitURL is called when exiting the gitURL production.
	ExitGitURL(c *GitURLContext)

	// ExitArgsList is called when exiting the argsList production.
	ExitArgsList(c *ArgsListContext)

	// ExitArg is called when exiting the arg production.
	ExitArg(c *ArgContext)
}
