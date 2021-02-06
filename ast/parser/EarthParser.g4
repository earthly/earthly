parser grammar EarthParser;

options {
	tokenVocab = EarthLexer;
}

earthFile: NL* (stmts NL)? NL* targets? NL* EOF;

targets: target WS? (NL+ DEDENT target WS?)* NL* DEDENT?;
target: targetHeader NL+ WS? INDENT stmts?;
targetHeader: Target;

stmts: WS? stmt (NL+ WS? stmt)*;

stmt:
	commandStmt
	| withDockerStmt
	| endStmt;

commandStmt:
	fromStmt
	| fromDockerfileStmt
	| locallyStmt
	| copyStmt
	| saveStmt
	| runStmt
	| buildStmt
	| workdirStmt
	| userStmt
	| cmdStmt
	| entrypointStmt
	| exposeStmt
	| volumeStmt
	| envStmt
	| argStmt
	| labelStmt
	| gitCloneStmt
	| addStmt
	| stopsignalStmt
	| onbuildStmt
	| healthcheckStmt
	| shellStmt
	| genericCommandStmt;

fromStmt: FROM (WS stmtWords)?;

fromDockerfileStmt: FROM_DOCKERFILE (WS stmtWords)?;

locallyStmt: LOCALLY (WS stmtWords)?;

copyStmt: COPY (WS stmtWords)?;

saveStmt: saveArtifact | saveImage;
saveImage: SAVE_IMAGE (WS stmtWords)?;
saveArtifact: SAVE_ARTIFACT (WS stmtWords)?;

runStmt: RUN (WS stmtWordsMaybeJSON)?;

buildStmt: BUILD (WS stmtWords)?;

workdirStmt: WORKDIR (WS stmtWords)?;

userStmt: USER (WS stmtWords)?;

cmdStmt: CMD (WS stmtWordsMaybeJSON)?;

entrypointStmt: ENTRYPOINT (WS stmtWordsMaybeJSON)?;

exposeStmt: EXPOSE (WS stmtWords)?;

volumeStmt: VOLUME (WS stmtWordsMaybeJSON)?;

envStmt: ENV WS envArgKey (WS? EQUALS)? (WS? envArgValue)?;
argStmt: ARG WS envArgKey ((WS? EQUALS) (WS? envArgValue)?)?;
envArgKey: Atom;
envArgValue: Atom (WS? Atom)*;

labelStmt: LABEL (WS labelKey WS? EQUALS WS? labelValue)*;
labelKey: Atom;
labelValue: Atom;

gitCloneStmt: GIT_CLONE (WS stmtWords)?;

addStmt: ADD (WS stmtWords)?;
stopsignalStmt: STOPSIGNAL (WS stmtWords)?;
onbuildStmt: ONBUILD (WS stmtWords)?;
healthcheckStmt: HEALTHCHECK (WS stmtWords)?;
shellStmt: SHELL (WS stmtWords)?;

withDockerStmt: WITH_DOCKER (WS stmtWords)?;
endStmt: END (WS stmtWords)?;

genericCommandStmt: commandName (WS stmtWords)?;
commandName: Command;

stmtWordsMaybeJSON: stmtWords;
stmtWords: stmtWord (WS? stmtWord)*;
stmtWord: Atom;
