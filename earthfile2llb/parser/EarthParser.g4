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
	| fromDockerfileStmt
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
	| dockerLoadStmt
	| dockerPullStmt
	| addStmt
	| stopsignalStmt
	| onbuildStmt
	| healthcheckStmt
	| shellStmt
	| withDockerStmt
	| endStmt
	| genericCommandStmt;

fromStmt: FROM (WS stmtWords)?;

fromDockerfileStmt: FROM_DOCKERFILE (WS stmtWords)?;

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

dockerLoadStmt: DOCKER_LOAD (WS stmtWords)?;

dockerPullStmt: DOCKER_PULL (WS stmtWords)?;

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
