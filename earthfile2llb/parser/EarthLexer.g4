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
EXPOSE: 'EXPOSE' -> pushMode(COMMAND_ARGS);
VOLUME: 'VOLUME' -> pushMode(COMMAND_ARGS);
ENV: 'ENV' -> pushMode(COMMAND_ARGS_KEY_VALUE);
ARG: 'ARG' -> pushMode(COMMAND_ARGS_KEY_VALUE);
LABEL: 'LABEL' -> pushMode(COMMAND_ARGS_KEY_VALUE_LABEL);
BUILD: 'BUILD' -> pushMode(COMMAND_ARGS);
WORKDIR: 'WORKDIR' -> pushMode(COMMAND_ARGS);
USER: 'USER' -> pushMode(COMMAND_ARGS);
CMD: 'CMD' -> pushMode(COMMAND_ARGS);
ENTRYPOINT: 'ENTRYPOINT' -> pushMode(COMMAND_ARGS);
GIT_CLONE: 'GIT CLONE' -> pushMode(COMMAND_ARGS);
DOCKER_LOAD: 'DOCKER LOAD' -> pushMode(COMMAND_ARGS);
DOCKER_PULL: 'DOCKER PULL' -> pushMode(COMMAND_ARGS);
ADD: 'ADD' -> pushMode(COMMAND_ARGS);
STOPSIGNAL: 'STOPSIGNAL' -> pushMode(COMMAND_ARGS);
ONBUILD: 'ONBUILD' -> pushMode(COMMAND_ARGS);
HEALTHCHECK: 'HEALTHCHECK' -> pushMode(COMMAND_ARGS);
SHELL: 'SHELL' -> pushMode(COMMAND_ARGS);
Command: [A-Z]+ -> pushMode(COMMAND_ARGS);

NL: WS? COMMENT? CRLF;
WS: [ \t] ([ \t] | LC)*;
fragment CRLF: ('\r' | '\n' | '\r\n');
fragment COMMENT: '#' (~[\r\n])*;
fragment LC: '\\' [ \t]* CRLF;

mode RECIPE;

// Note: RECIPE mode is popped via golang code, when DEDENT occurs.

Target_R: Target -> type(Target), pushMode(RECIPE);

FROM_R: FROM -> type(FROM), pushMode(COMMAND_ARGS);
COPY_R: COPY -> type(COPY), pushMode(COMMAND_ARGS);
SAVE_ARTIFACT_R: SAVE_ARTIFACT -> type(SAVE_ARTIFACT), pushMode(COMMAND_ARGS);
SAVE_IMAGE_R: SAVE_IMAGE -> type(SAVE_IMAGE), pushMode(COMMAND_ARGS);
RUN_R: RUN -> type(RUN), pushMode(COMMAND_ARGS);
EXPOSE_R: EXPOSE -> type(EXPOSE), pushMode(COMMAND_ARGS);
VOLUME_R: VOLUME -> type(VOLUME), pushMode(COMMAND_ARGS);
ENV_R: ENV -> type(ENV), pushMode(COMMAND_ARGS_KEY_VALUE);
ARG_R: ARG -> type(ARG), pushMode(COMMAND_ARGS_KEY_VALUE);
LABEL_R: LABEL -> type(LABEL), pushMode(COMMAND_ARGS_KEY_VALUE_LABEL);
BUILD_R: BUILD -> type(BUILD), pushMode(COMMAND_ARGS);
WORKDIR_R: WORKDIR -> type(WORKDIR), pushMode(COMMAND_ARGS);
USER_R: USER -> type(USER), pushMode(COMMAND_ARGS);
CMD_R: CMD -> type(CMD), pushMode(COMMAND_ARGS);
ENTRYPOINT_R: ENTRYPOINT -> type(ENTRYPOINT), pushMode(COMMAND_ARGS);
GIT_CLONE_R: GIT_CLONE -> type(GIT_CLONE), pushMode(COMMAND_ARGS);
DOCKER_LOAD_R: DOCKER_LOAD -> type(DOCKER_LOAD), pushMode(COMMAND_ARGS);
DOCKER_PULL_R: DOCKER_PULL -> type(DOCKER_PULL), pushMode(COMMAND_ARGS);
ADD_R: ADD -> type(ADD), pushMode(COMMAND_ARGS);
STOPSIGNAL_R: STOPSIGNAL -> type(STOPSIGNAL), pushMode(COMMAND_ARGS);
ONBUILD_R: ONBUILD -> type(ONBUILD), pushMode(COMMAND_ARGS);
HEALTHCHECK_R: HEALTHCHECK -> type(HEALTHCHECK), pushMode(COMMAND_ARGS);
SHELL_R: SHELL -> type(SHELL), pushMode(COMMAND_ARGS);
Command_R: Command -> type(Command), pushMode(COMMAND_ARGS);

NL_R: NL -> type(NL);
WS_R: WS -> type(WS);

mode COMMAND_ARGS;

Atom: (RegularAtomPart | QuotedAtomPart)+;
fragment QuotedAtomPart: ('"' (~'"' | '\\"')* '"');
fragment RegularAtomPart: ~([ \t\r\n\\"]) | EscapedAtomPart;
fragment EscapedAtomPart: ('\\' .) | (LC [ \t]*);

// Note; Comments not allowed in command lines.
NL_C: WS? CRLF -> type(NL), popMode;
WS_C: WS -> type(WS);

mode COMMAND_ARGS_KEY_VALUE;

// Switch mode after '=' (may contain '=' as part of value after that).
EQUALS: '=' -> mode(COMMAND_ARGS);

// Similar Atom, but don't allow '=' as part of it, unless it's in quotes.
Atom_CAKV: (RegularAtomPart_CAKV | QuotedAtomPart)+ -> type(Atom);
fragment RegularAtomPart_CAKV: ~([ \t\r\n"=]) | EscapedAtomPart;

// Note; Comments not allowed in command lines.
NL_CAKV: WS? CRLF -> type(NL), popMode;
WS_CAKV: WS -> type(WS);

mode COMMAND_ARGS_KEY_VALUE_LABEL;

EQUALS_L: '=' -> type(EQUALS);

// Similar Atom, but don't allow '=' as part of it, unless it's in quotes.
Atom_CAKVL: Atom_CAKV -> type(Atom);

// Note; Comments not allowed in command lines.
NL_CAKVL: NL_CAKV -> type(NL), popMode;
WS_CAKVL: WS_CAKV -> type(WS);
