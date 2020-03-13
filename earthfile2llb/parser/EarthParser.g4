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
	fromStmt
	| copyStmt
	| saveStmt
	| runStmt
	| buildStmt
	| workdirStmt
	| entrypointStmt
	| envStmt
	| argStmt
	| gitCloneStmt
	| dockerLoadStmt
	| dockerPullStmt
	| genericCommand;

fromStmt: FROM (WS flagKeyValue)* WS imageName (WS AS asName)?;
asName: Atom;

copyStmt: COPY WS stmtWords;

saveStmt: saveArtifact | saveImage;
saveImage: SAVE_IMAGE (WS saveImageName)*;
saveArtifact:
	SAVE_ARTIFACT WS saveFrom (WS saveTo)? (WS AS_LOCAL WS saveAsLocalTo)?;
saveFrom: Atom;
saveTo: Atom;
saveAsLocalTo: Atom;

runStmt: RUN (WS flag)* WS (runArgs | runArgsList);

buildStmt: BUILD (WS flagKeyValue)* WS fullTargetName;

workdirStmt: WORKDIR WS workdirPath;
workdirPath: Atom;

entrypointStmt: ENTRYPOINT WS (entrypointArgs | entrypointArgsList);

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

runArgs: runArg (WS runArg)*;
runArgsList:
	OPEN_BRACKET WS? runArg (WS? COMMA WS? runArg)+ WS? CLOSE_BRACKET;
runArg: Atom;

entrypointArgs: entrypointArg (WS entrypointArg)*;
entrypointArgsList:
	OPEN_BRACKET WS? entrypointArg (WS? COMMA WS? entrypointArg)+ WS? CLOSE_BRACKET;
entrypointArg: Atom;

flags: flag (WS? flag)*;
flag: flagKey | flagKeyValue;
flagKey: FlagKey;
flagKeyValue: FlagKeyValue;

stmtWords: stmtWord (WS? stmtWord)*;
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
