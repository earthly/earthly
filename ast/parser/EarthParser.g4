parser grammar EarthParser;

options {
	tokenVocab = EarthLexer;
}

earthFile: NL* version? (stmts NL)? NL* targets? NL* EOF;

targets: targetOrUserCommand (NL* targetOrUserCommand)*;
targetOrUserCommand: target | userCommand;
target: targetHeader NL+ WS? (INDENT stmts NL+ DEDENT)?;
targetHeader: Target;
userCommand: userCommandHeader NL+ WS? (INDENT stmts NL+ DEDENT)?;
userCommandHeader: UserCommand;

stmts: WS? stmt (NL+ WS? stmt)*;

stmt:
	commandStmt
	| withStmt
	| ifStmt
	| forStmt
	| waitStmt;

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
	| userCommandStmt
	| doStmt
	| importStmt
	| cacheStmt
	| hostStmt
	| projectStmt;

// version --------------------------------------------------------------------
version: VERSION WS stmtWords NL+;

// withStmt -------------------------------------------------------------------

withStmt: withExpr (NL+ WS? withBlock)? NL+ WS? END;
withBlock: stmts;

withExpr: WITH WS withCommand;
withCommand:
	dockerCommand;

dockerCommand: DOCKER (WS stmtWords)?;

// ifStmt ---------------------------------------------------------------------

ifStmt: ifClause (NL+ WS? elseIfClause)* (NL+ WS? elseClause)? NL+ WS? END;

ifClause: IF WS ifExpr (NL+ WS? ifBlock)?;
ifBlock: stmts;
elseIfClause: ELSE_IF WS elseIfExpr (NL+ WS? elseIfBlock)?;
elseIfBlock: stmts;
elseClause: ELSE (NL+ WS? elseBlock)?;
elseBlock: stmts;

ifExpr: expr;
elseIfExpr: expr;

// forStmt --------------------------------------------------------------------

forStmt: forClause NL+ WS? END;

forClause: FOR WS forExpr (NL+ WS? forBlock)?;
forBlock: stmts;

forExpr: stmtWords;

// waitStmt --------------------------------------------------------------------

waitStmt: waitClause NL+ WS? END;
waitClause: WAIT (WS waitExpr)? (NL+ WS? waitBlock)?;
waitBlock: stmts;
waitExpr: stmtWords;

// Regular commands -----------------------------------------------------------

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
argStmt: ARG optionalFlag WS envArgKey ((WS? EQUALS) (WS? envArgValue)?)?;
optionalFlag: (WS stmtWords)?;
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
userCommandStmt: COMMAND (WS stmtWords)?;
doStmt: DO (WS stmtWords)?;
importStmt: IMPORT (WS stmtWords)?;
cacheStmt: CACHE (WS stmtWords)?;
hostStmt: HOST (WS stmtWords)?;
projectStmt: PROJECT (WS stmtWords)?;

// expr, stmtWord* ------------------------------------------------------------

expr: stmtWordsMaybeJSON;

stmtWordsMaybeJSON: stmtWords;
stmtWords: stmtWord (WS stmtWord)*;
stmtWord: Atom;
