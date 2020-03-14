parser grammar EarthParser;

options {
	tokenVocab = EarthLexer;
}

earthFile: NL* stmts? NL+ targets? NL* EOF;

targets: target WS? (NL+ DEDENT target WS?)* NL* DEDENT?;
target: targetHeader NL+ WS? INDENT stmts?;
targetHeader: Target;

stmts: WS? stmt (NL+ WS? stmt)*;

stmt:
	fromStmt  // Updated.
	| copyStmt  // Updated.
	| saveStmt  // Updated.
	| runStmt  // Updated.
	| buildStmt  // Updated.
	| workdirStmt // Updated.
	| entrypointStmt  // Updated.
	| envStmt
	| argStmt
	| gitCloneStmt  // Updated.
	| dockerLoadStmt  // Updated.
	| dockerPullStmt  // Updated.
	| genericCommand;  // Updated.

fromStmt: FROM (WS stmtWords)?;

copyStmt: COPY (WS stmtWords)?;

saveStmt: saveArtifact | saveImage;
saveImage: SAVE_IMAGE (WS stmtWords)?;
saveArtifact: SAVE_ARTIFACT (WS stmtWords)?;

runStmt: RUN (WS (stmtWords | stmtWordsList))?;

buildStmt: BUILD (WS stmtWords)?;

workdirStmt: WORKDIR (WS stmtWords)?;

entrypointStmt: ENTRYPOINT (WS (stmtWords | stmtWordsList))?;

envStmt: ENV WS envArgKey (WS? EQUALS)? (WS? envArgValue)?;
argStmt: ARG WS envArgKey ((WS? EQUALS) (WS? envArgValue)?)?;
envArgKey: Atom;
envArgValue: Atom (WS? Atom)*;

gitCloneStmt: GIT_CLONE (WS stmtWords)?;

dockerLoadStmt: DOCKER_LOAD (WS stmtWords)?;

dockerPullStmt: DOCKER_PULL (WS stmtWords)?;

genericCommand: commandName (WS stmtWords)?;
commandName: Command;

stmtWords: stmtWord (WS? stmtWord)*;
stmtWordsList:
	OPEN_BRACKET WS? (stmtWord (WS? COMMA WS? stmtWord)* WS?)? CLOSE_BRACKET;
stmtWord: Atom;
