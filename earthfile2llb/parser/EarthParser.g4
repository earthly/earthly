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
	| saveStmt
	| runStmt  // Updated.
	| buildStmt
	| workdirStmt
	| entrypointStmt  // Updated.
	| envStmt
	| argStmt
	| gitCloneStmt
	| dockerLoadStmt
	| dockerPullStmt
	| genericCommand;

fromStmt: FROM (WS stmtWords)?;

copyStmt: COPY (WS stmtWords)?;

saveStmt: saveArtifact | saveImage;
saveImage: SAVE_IMAGE (WS saveImageName)*;
saveArtifact:
	SAVE_ARTIFACT WS saveFrom (WS saveTo)? (WS AS_LOCAL WS saveAsLocalTo)?;
saveFrom: Atom;
saveTo: Atom;
saveAsLocalTo: Atom;

runStmt: RUN (WS (stmtWords | stmtWordsList))?;

buildStmt: BUILD (WS flagKeyValue)* WS fullTargetName;

workdirStmt: WORKDIR WS workdirPath;
workdirPath: Atom;

entrypointStmt: ENTRYPOINT (WS (stmtWords | stmtWordsList))?;

envStmt: ENV WS envArgKey (WS? EQUALS)? (WS? envArgValue)?;

argStmt: ARG WS envArgKey ((WS? EQUALS) (WS? envArgValue)?)?;

gitCloneStmt: GIT_CLONE (WS flagKeyValue)* WS gitURL WS gitCloneDest;
gitURL: Atom;
gitCloneDest: Atom;

dockerLoadStmt:
	DOCKER_LOAD (WS flagKeyValue)* WS fullTargetName WS AS WS imageName;

dockerPullStmt: DOCKER_PULL WS imageName;

genericCommand:
	commandName (WS flags)? (WS stmtWords | WS argsList)?;
commandName: Command;

flags: flag (WS? flag)*;
flag: flagKey | flagKeyValue;
flagKey: FlagKey;
flagKeyValue: FlagKeyValue;

stmtWords: stmtWord (WS? stmtWord)*;
stmtWordsList:
	OPEN_BRACKET WS? (stmtWord (WS? COMMA WS? stmtWord)* WS?)? CLOSE_BRACKET;
stmtWord: Atom;

envArgKey: Atom;
envArgValue: Atom (WS? Atom)*;

imageName: Atom;
saveImageName: Atom;
targetName: Atom;
fullTargetName: Atom;

argsList:
	OPEN_BRACKET WS? arg (WS? COMMA WS? arg)+ WS? CLOSE_BRACKET;
arg: Atom;
