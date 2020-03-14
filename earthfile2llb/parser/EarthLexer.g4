lexer grammar EarthLexer;

tokens {
	INDENT,
	DEDENT
}

Target: ([a-zA-Z0-9.] | '-')+ ':' -> pushMode(RECIPE);

FROM: 'FROM' -> pushMode(COMMAND_ARGS);
COPY: 'COPY' -> pushMode(COMMAND_ARGS);
SAVE_ARTIFACT: 'SAVE ARTIFACT' -> pushMode(COMMAND_ARGS);
SAVE_IMAGE: 'SAVE IMAGE' -> pushMode(COMMAND_ARGS);
RUN: 'RUN' -> pushMode(COMMAND_ARGS);
ENV: 'ENV' -> pushMode(COMMAND_ARGS_KEY_VALUE);
ARG: 'ARG' -> pushMode(COMMAND_ARGS_KEY_VALUE);
BUILD: 'BUILD' -> pushMode(COMMAND_ARGS);
WORKDIR: 'WORKDIR' -> pushMode(COMMAND_ARGS);
ENTRYPOINT: 'ENTRYPOINT' -> pushMode(COMMAND_ARGS);
GIT_CLONE: 'GIT CLONE' -> pushMode(COMMAND_ARGS);
DOCKER_LOAD:
	'DOCKER LOAD' -> pushMode(COMMAND_ARGS);
DOCKER_PULL:
	'DOCKER PULL' -> pushMode(COMMAND_ARGS);
Command: [A-Z]+ -> pushMode(COMMAND_ARGS);

NL: WS? COMMENT? CRLF;
WS: ([ \t] | ('\\' [ \t]* CRLF))+;
fragment CRLF: ('\r' | '\n' | '\r\n');
fragment COMMENT: '#' (~[\r\n])*;

mode RECIPE;

// Note: RECIPE mode is popped via golang code, when DEDENT occurs.

Target_R: Target -> type(Target), pushMode(RECIPE);

FROM_R:
	FROM -> type(FROM), pushMode(COMMAND_ARGS);
COPY_R: COPY -> type(COPY), pushMode(COMMAND_ARGS);
SAVE_ARTIFACT_R: SAVE_ARTIFACT -> type(SAVE_ARTIFACT), pushMode(COMMAND_ARGS);
SAVE_IMAGE_R: SAVE_IMAGE -> type(SAVE_IMAGE), pushMode(COMMAND_ARGS);
RUN_R: RUN -> type(RUN), pushMode(COMMAND_ARGS);
ENV_R: ENV -> type(ENV), pushMode(COMMAND_ARGS_KEY_VALUE);
ARG_R: ARG -> type(ARG), pushMode(COMMAND_ARGS_KEY_VALUE);
BUILD_R:
	BUILD -> type(BUILD), pushMode(COMMAND_ARGS);
WORKDIR_R:
	WORKDIR -> type(WORKDIR), pushMode(COMMAND_ARGS);
ENTRYPOINT_R:
	ENTRYPOINT -> type(ENTRYPOINT), pushMode(COMMAND_ARGS);
GIT_CLONE_R:
	GIT_CLONE -> type(GIT_CLONE), pushMode(COMMAND_ARGS);
DOCKER_LOAD_R:
	DOCKER_LOAD -> type(DOCKER_LOAD), pushMode(COMMAND_ARGS);
DOCKER_PULL_R:
	DOCKER_PULL -> type(DOCKER_PULL), pushMode(COMMAND_ARGS);
Command_R: Command -> type(Command), pushMode(COMMAND_ARGS);

NL_R: NL -> type(NL);
WS_R: WS -> type(WS);

mode COMMAND_ARGS;

Atom: (NonWSNLQuote | QuotedAtom)+;
fragment QuotedAtom: ('"' (~'"' | '\\"')* '"');
fragment NonWSNLQuote: ~([ \t\r\n"]);

// Note; Comments not allowed in command lines.
NL_C: WS? CRLF -> type(NL), popMode;
WS_C: WS -> type(WS);

mode COMMAND_ARGS_KEY_VALUE;

// Switch mode after '=' (may contain '=' as part of value after that).
EQUALS: '=' -> mode(COMMAND_ARGS);

// Similar Atom, but don't allow '=' as part of it, unless it's in quotes.
Atom_CAKV: (NonWSNLQuote_CAKV | QuotedAtom)+ -> type(Atom);
fragment NonWSNLQuote_CAKV: ~([ \t\r\n"=]);

// Note; Comments not allowed in command lines.
NL_CAKV: WS? CRLF -> type(NL), popMode;
WS_CAKC: WS -> type(WS);
