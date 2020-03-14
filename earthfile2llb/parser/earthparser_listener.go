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

	// EnterEnvArgKey is called when entering the envArgKey production.
	EnterEnvArgKey(c *EnvArgKeyContext)

	// EnterEnvArgValue is called when entering the envArgValue production.
	EnterEnvArgValue(c *EnvArgValueContext)

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

	// EnterStmtWords is called when entering the stmtWords production.
	EnterStmtWords(c *StmtWordsContext)

	// EnterStmtWordsList is called when entering the stmtWordsList production.
	EnterStmtWordsList(c *StmtWordsListContext)

	// EnterStmtWord is called when entering the stmtWord production.
	EnterStmtWord(c *StmtWordContext)

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

	// ExitEnvArgKey is called when exiting the envArgKey production.
	ExitEnvArgKey(c *EnvArgKeyContext)

	// ExitEnvArgValue is called when exiting the envArgValue production.
	ExitEnvArgValue(c *EnvArgValueContext)

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

	// ExitStmtWords is called when exiting the stmtWords production.
	ExitStmtWords(c *StmtWordsContext)

	// ExitStmtWordsList is called when exiting the stmtWordsList production.
	ExitStmtWordsList(c *StmtWordsListContext)

	// ExitStmtWord is called when exiting the stmtWord production.
	ExitStmtWord(c *StmtWordContext)
}
