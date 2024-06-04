parser grammar EarthParser;

options {
	tokenVocab = EarthLexer;
}

earthFile: NL* version? (stmts NL)? NL* targets? NL* EOF;

targets: targetOrUserCommand (NL* targetOrUserCommand)*;
targetOrUserCommand: target | userCommand;
target: targetHeader NL+ (INDENT NL* stmts? NL+ DEDENT)?;
targetHeader: Target;
userCommand: userCommandHeader NL+ (INDENT NL* stmts NL+ DEDENT)?;
userCommandHeader: UserCommand;
function: functionHeader NL+ (INDENT NL* stmts NL+ DEDENT)?;
functionHeader: Function;

stmts: stmt (NL+ stmt)*;

stmt:
	commandStmt
	| withStmt
	| ifStmt
	| forStmt
	| waitStmt
	| tryStmt;

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
	| setStmt
	| letStmt
	| labelStmt
	| gitCloneStmt
	| addStmt
	| stopsignalStmt
	| onbuildStmt
	| healthcheckStmt
	| shellStmt
	| userCommandStmt
	| functionStmt
	| doStmt
	| importStmt
	| cacheStmt
	| hostStmt
	| projectStmt;

// version --------------------------------------------------------------------
version: VERSION stmtWords NL+;

// withStmt -------------------------------------------------------------------

withStmt: withExpr (NL+ withBlock)? NL+ END;
withBlock: stmts;

withExpr: WITH withCommand;
withCommand:
	dockerCommand;

dockerCommand: DOCKER stmtWords?;

// ifStmt ---------------------------------------------------------------------

ifStmt: ifClause (NL+ elseIfClause)* (NL+ elseClause)? NL+ END;

ifClause: IF ifExpr (NL+ ifBlock)?;
ifBlock: stmts;
elseIfClause: ELSE_IF elseIfExpr (NL+ elseIfBlock)?;
elseIfBlock: stmts;
elseClause: ELSE (NL+ elseBlock)?;
elseBlock: stmts;

ifExpr: expr;
elseIfExpr: expr;

// tryStmt ---------------------------------------------------------------------

tryStmt: tryClause (NL+ catchClause)? (NL+ finallyClause)? NL+ END;

tryClause: TRY (NL+ tryBlock)?;
tryBlock: stmts;
catchClause: CATCH (NL+ catchBlock)?;
catchBlock: stmts;
finallyClause: FINALLY (NL+ finallyBlock)?;
finallyBlock: stmts;


// forStmt --------------------------------------------------------------------

forStmt: forClause NL+ END;

forClause: FOR forExpr (NL+ forBlock)?;
forBlock: stmts;

forExpr: stmtWords;

// waitStmt --------------------------------------------------------------------

waitStmt: waitClause NL+ END;
waitClause: WAIT waitExpr? (NL+ waitBlock)?;
waitBlock: stmts;
waitExpr: stmtWords;

// Regular commands -----------------------------------------------------------

fromStmt: FROM stmtWords?;

fromDockerfileStmt: FROM_DOCKERFILE stmtWords?;

locallyStmt: LOCALLY stmtWords?;

copyStmt: COPY stmtWords?;

saveStmt: saveArtifact | saveImage;
saveImage: SAVE_IMAGE stmtWords?;
saveArtifact: SAVE_ARTIFACT stmtWords?;

runStmt: RUN stmtWordsMaybeJSON?;

buildStmt: BUILD stmtWords?;

workdirStmt: WORKDIR stmtWords?;

userStmt: USER stmtWords?;

cmdStmt: CMD stmtWordsMaybeJSON?;

entrypointStmt: ENTRYPOINT stmtWordsMaybeJSON?;

exposeStmt: EXPOSE stmtWords?;

volumeStmt: VOLUME stmtWordsMaybeJSON?;

envStmt: ENV envArgKey EQUALS? (WS? envArgValue)?;
argStmt: ARG optionalFlag envArgKey (EQUALS (WS? envArgValue)?)?;
setStmt: SET envArgKey EQUALS WS? envArgValue;
letStmt: LET optionalFlag envArgKey EQUALS WS? envArgValue;
optionalFlag: stmtWords?;
envArgKey: Atom;
envArgValue: Atom (WS? Atom)*;

labelStmt: LABEL (labelKey EQUALS labelValue)*;
labelKey: Atom;
labelValue: Atom;

gitCloneStmt: GIT_CLONE stmtWords?;

addStmt: ADD stmtWords?;
stopsignalStmt: STOPSIGNAL stmtWords?;
onbuildStmt: ONBUILD stmtWords?;
healthcheckStmt: HEALTHCHECK stmtWords?;
shellStmt: SHELL stmtWords?;
userCommandStmt: COMMAND stmtWords?;
functionStmt: FUNCTION stmtWords?;
doStmt: DO stmtWords?;
importStmt: IMPORT stmtWords?;
cacheStmt: CACHE stmtWords?;
hostStmt: HOST stmtWords?;
projectStmt: PROJECT stmtWords?;

// expr, stmtWord* ------------------------------------------------------------

expr: stmtWordsMaybeJSON;

stmtWordsMaybeJSON: stmtWords;
stmtWords: stmtWord+;
stmtWord: Atom;
