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
	| projectStmt stmtWords?;

// expr, stmtWord* ------------------------------------------------------------

expr: stmtWordsMaybeJSON;

stmtWordsMaybeJSON: stmtWords;
stmtWords: stmtWord+;
stmtWord: Atom;
