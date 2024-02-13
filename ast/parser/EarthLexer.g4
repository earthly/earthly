lexer grammar EarthLexer;

tokens {
	INDENT,
	DEDENT
}

channels {
    WHITESPACE_CHANNEL,
    COMMENTS_CHANNEL
}

Target: [a-z] ([a-zA-Z0-9.] | '-')* ':' -> pushMode(RECIPE);
UserCommand: [A-Z] ([A-Z0-9._])* ':' -> pushMode(RECIPE);
Function: [A-Z] ([A-Z0-9._])* ':' -> pushMode(RECIPE);

FROM: 'FROM' -> pushMode(COMMAND_ARGS);
FROM_DOCKERFILE: 'FROM DOCKERFILE' -> pushMode(COMMAND_ARGS);
LOCALLY: 'LOCALLY' -> pushMode(COMMAND_ARGS);
COPY: 'COPY' -> pushMode(COMMAND_ARGS);
SAVE_ARTIFACT: 'SAVE ARTIFACT' -> pushMode(COMMAND_ARGS);
SAVE_IMAGE: 'SAVE IMAGE' -> pushMode(COMMAND_ARGS);
RUN: 'RUN' -> pushMode(COMMAND_ARGS);
EXPOSE: 'EXPOSE' -> pushMode(COMMAND_ARGS);
VOLUME: 'VOLUME' -> pushMode(COMMAND_ARGS);
ENV: 'ENV' -> pushMode(COMMAND_ARGS_KEY_VALUE);
ARG: 'ARG' -> pushMode(COMMAND_ARGS_KEY_VALUE);
SET: 'SET' -> pushMode(COMMAND_ARGS_KEY_VALUE);
LET: 'LET' -> pushMode(COMMAND_ARGS_KEY_VALUE);
LABEL: 'LABEL' -> pushMode(COMMAND_ARGS_KEY_VALUE_LABEL);
BUILD: 'BUILD' -> pushMode(COMMAND_ARGS);
WORKDIR: 'WORKDIR' -> pushMode(COMMAND_ARGS);
USER: 'USER' -> pushMode(COMMAND_ARGS);
CMD: 'CMD' -> pushMode(COMMAND_ARGS);
ENTRYPOINT: 'ENTRYPOINT' -> pushMode(COMMAND_ARGS);
GIT_CLONE: 'GIT CLONE' -> pushMode(COMMAND_ARGS);
ADD: 'ADD' -> pushMode(COMMAND_ARGS);
STOPSIGNAL: 'STOPSIGNAL' -> pushMode(COMMAND_ARGS);
ONBUILD: 'ONBUILD' -> pushMode(COMMAND_ARGS);
HEALTHCHECK: 'HEALTHCHECK' -> pushMode(COMMAND_ARGS);
SHELL: 'SHELL' -> pushMode(COMMAND_ARGS);
DO: 'DO' -> pushMode(COMMAND_ARGS);
COMMAND: 'COMMAND' -> pushMode(COMMAND_ARGS);
FUNCTION: 'FUNCTION' -> pushMode(COMMAND_ARGS);
IMPORT: 'IMPORT' -> pushMode(COMMAND_ARGS);
VERSION: 'VERSION' -> pushMode(COMMAND_ARGS);
CACHE: 'CACHE' -> pushMode(COMMAND_ARGS);
HOST: 'HOST' -> pushMode(COMMAND_ARGS);
PROJECT: 'PROJECT' -> pushMode(COMMAND_ARGS);
PIPELINE: 'PIPELINE' -> pushMode(COMMAND_ARGS);
TRIGGER: 'TRIGGER' -> pushMode(COMMAND_ARGS);

WITH: 'WITH';
DOCKER: 'DOCKER' -> pushMode(BLOCK), pushMode(COMMAND_ARGS);
IF: 'IF' -> pushMode(BLOCK), pushMode(COMMAND_ARGS);
TRY: 'TRY' -> pushMode(BLOCK), pushMode(COMMAND_ARGS);
FOR: 'FOR' -> pushMode(BLOCK), pushMode(COMMAND_ARGS);
WAIT: 'WAIT' -> pushMode(BLOCK), pushMode(COMMAND_ARGS);

NL: [ \t]* (EOF | CRLF);
WS: [ \t] ([ \t] | LC)* -> channel(WHITESPACE_CHANNEL);
COMMENT: [ \t]* '#' (~[\r\n])* -> channel(COMMENTS_CHANNEL);
fragment CRLF: ('\r' | '\n' | '\r\n');

// TODO: figure out if adding COMMENT explicitly is necessary.
fragment NL_NOLC: ([ \t]* CRLF | [ \t]* COMMENT);
fragment LC: '\\' NL_NOLC+;

// ----------------------------------------------------------------------------

mode RECIPE;

// Note: RECIPE mode is popped via golang code, when DEDENT occurs.

Target_R: Target -> type(Target);
UserCommand_R: UserCommand -> type(UserCommand);
Function_R: Function -> type(Function);

FROM_R: FROM -> type(FROM), pushMode(COMMAND_ARGS);
FROM_DOCKERFILE_R: FROM_DOCKERFILE -> type(FROM_DOCKERFILE), pushMode(COMMAND_ARGS);
LOCALLY_R: LOCALLY -> type(LOCALLY), pushMode(COMMAND_ARGS);
COPY_R: COPY -> type(COPY), pushMode(COMMAND_ARGS);
SAVE_ARTIFACT_R: SAVE_ARTIFACT -> type(SAVE_ARTIFACT), pushMode(COMMAND_ARGS);
SAVE_IMAGE_R: SAVE_IMAGE -> type(SAVE_IMAGE), pushMode(COMMAND_ARGS);
RUN_R: RUN -> type(RUN), pushMode(COMMAND_ARGS);
EXPOSE_R: EXPOSE -> type(EXPOSE), pushMode(COMMAND_ARGS);
VOLUME_R: VOLUME -> type(VOLUME), pushMode(COMMAND_ARGS);
ENV_R: ENV -> type(ENV), pushMode(COMMAND_ARGS_KEY_VALUE);
ARG_R: ARG -> type(ARG), pushMode(COMMAND_ARGS_KEY_VALUE);
SET_R: SET -> type(SET), pushMode(COMMAND_ARGS_KEY_VALUE);
LET_R: LET -> type(LET), pushMode(COMMAND_ARGS_KEY_VALUE);
LABEL_R: LABEL -> type(LABEL), pushMode(COMMAND_ARGS_KEY_VALUE_LABEL);
BUILD_R: BUILD -> type(BUILD), pushMode(COMMAND_ARGS);
WORKDIR_R: WORKDIR -> type(WORKDIR), pushMode(COMMAND_ARGS);
USER_R: USER -> type(USER), pushMode(COMMAND_ARGS);
CMD_R: CMD -> type(CMD), pushMode(COMMAND_ARGS);
ENTRYPOINT_R: ENTRYPOINT -> type(ENTRYPOINT), pushMode(COMMAND_ARGS);
GIT_CLONE_R: GIT_CLONE -> type(GIT_CLONE), pushMode(COMMAND_ARGS);
ADD_R: ADD -> type(ADD), pushMode(COMMAND_ARGS);
STOPSIGNAL_R: STOPSIGNAL -> type(STOPSIGNAL), pushMode(COMMAND_ARGS);
ONBUILD_R: ONBUILD -> type(ONBUILD), pushMode(COMMAND_ARGS);
HEALTHCHECK_R: HEALTHCHECK -> type(HEALTHCHECK), pushMode(COMMAND_ARGS);
SHELL_R: SHELL -> type(SHELL), pushMode(COMMAND_ARGS);
DO_R: DO -> type(DO), pushMode(COMMAND_ARGS);
COMMAND_R: COMMAND -> type(COMMAND), pushMode(COMMAND_ARGS);
FUNCTION_R: FUNCTION -> type(FUNCTION), pushMode(COMMAND_ARGS);
IMPORT_R: IMPORT -> type(IMPORT), pushMode(COMMAND_ARGS);
CACHE_R: CACHE -> type(CACHE), pushMode(COMMAND_ARGS);
HOST_R: HOST -> type(HOST), pushMode(COMMAND_ARGS);
PIPELINE_R: PIPELINE -> type(PIPELINE), pushMode(COMMAND_ARGS);
TRIGGER_R: TRIGGER -> type(TRIGGER), pushMode(COMMAND_ARGS);

WITH_R: WITH -> type(WITH);
DOCKER_R: DOCKER -> type(DOCKER), pushMode(BLOCK), pushMode(COMMAND_ARGS);
IF_R: IF -> type(IF), pushMode(BLOCK), pushMode(COMMAND_ARGS);
TRY_R: TRY -> type(TRY), pushMode(BLOCK), pushMode(COMMAND_ARGS);
FOR_R: FOR -> type(FOR), pushMode(BLOCK), pushMode(COMMAND_ARGS);
WAIT_R: WAIT -> type(WAIT), pushMode(BLOCK), pushMode(COMMAND_ARGS);

NL_R: NL -> type(NL);
WS_R: WS -> type(WS), channel(WHITESPACE_CHANNEL);
COMMENT_R: COMMENT -> type(COMMENT), channel(COMMENTS_CHANNEL);

// ----------------------------------------------------------------------------

mode BLOCK;

FROM_B: FROM -> type(FROM), pushMode(COMMAND_ARGS);
FROM_DOCKERFILE_B: FROM_DOCKERFILE -> type(FROM_DOCKERFILE), pushMode(COMMAND_ARGS);
LOCALLY_B: LOCALLY -> type(LOCALLY), pushMode(COMMAND_ARGS);
COPY_B: COPY -> type(COPY), pushMode(COMMAND_ARGS);
SAVE_ARTIFACT_B: SAVE_ARTIFACT -> type(SAVE_ARTIFACT), pushMode(COMMAND_ARGS);
SAVE_IMAGE_B: SAVE_IMAGE -> type(SAVE_IMAGE), pushMode(COMMAND_ARGS);
RUN_B: RUN -> type(RUN), pushMode(COMMAND_ARGS);
EXPOSE_B: EXPOSE -> type(EXPOSE), pushMode(COMMAND_ARGS);
VOLUME_B: VOLUME -> type(VOLUME), pushMode(COMMAND_ARGS);
ENV_B: ENV -> type(ENV), pushMode(COMMAND_ARGS_KEY_VALUE);
ARG_B: ARG -> type(ARG), pushMode(COMMAND_ARGS_KEY_VALUE);
SET_B: SET -> type(SET), pushMode(COMMAND_ARGS_KEY_VALUE);
LET_B: LET -> type(LET), pushMode(COMMAND_ARGS_KEY_VALUE);
LABEL_B: LABEL -> type(LABEL), pushMode(COMMAND_ARGS_KEY_VALUE_LABEL);
BUILD_B: BUILD -> type(BUILD), pushMode(COMMAND_ARGS);
WORKDIR_B: WORKDIR -> type(WORKDIR), pushMode(COMMAND_ARGS);
USER_B: USER -> type(USER), pushMode(COMMAND_ARGS);
CMD_B: CMD -> type(CMD), pushMode(COMMAND_ARGS);
ENTRYPOINT_B: ENTRYPOINT -> type(ENTRYPOINT), pushMode(COMMAND_ARGS);
GIT_CLONE_B: GIT_CLONE -> type(GIT_CLONE), pushMode(COMMAND_ARGS);
ADD_B: ADD -> type(ADD), pushMode(COMMAND_ARGS);
STOPSIGNAL_B: STOPSIGNAL -> type(STOPSIGNAL), pushMode(COMMAND_ARGS);
ONBUILD_B: ONBUILD -> type(ONBUILD), pushMode(COMMAND_ARGS);
HEALTHCHECK_B: HEALTHCHECK -> type(HEALTHCHECK), pushMode(COMMAND_ARGS);
SHELL_B: SHELL -> type(SHELL), pushMode(COMMAND_ARGS);
DO_B: DO -> type(DO), pushMode(COMMAND_ARGS);
COMMAND_B: COMMAND -> type(COMMAND), pushMode(COMMAND_ARGS);
FUNCTION_B: FUNCTION -> type(FUNCTION), pushMode(COMMAND_ARGS);
IMPORT_B: IMPORT -> type(IMPORT), pushMode(COMMAND_ARGS);
CACHE_B: CACHE -> type(CACHE), pushMode(COMMAND_ARGS);
HOST_B: HOST -> type(HOST), pushMode(COMMAND_ARGS);

WITH_B: WITH -> type(WITH);
DOCKER_B: DOCKER -> type(DOCKER), pushMode(BLOCK), pushMode(COMMAND_ARGS);
IF_B: IF -> type(IF), pushMode(BLOCK), pushMode(COMMAND_ARGS);
ELSE: 'ELSE' -> pushMode(COMMAND_ARGS);
ELSE_IF: 'ELSE IF' -> pushMode(COMMAND_ARGS);
TRY_B: TRY -> type(TRY), pushMode(BLOCK), pushMode(COMMAND_ARGS);
CATCH: 'CATCH' -> pushMode(COMMAND_ARGS);
FINALLY: 'FINALLY' -> pushMode(COMMAND_ARGS);
FOR_B: FOR -> type(FOR), pushMode(BLOCK), pushMode(COMMAND_ARGS);
WAIT_B: WAIT -> type(WAIT), pushMode(BLOCK);
END: 'END' -> popMode, pushMode(COMMAND_ARGS);

NL_B: NL -> type(NL);
WS_B: WS -> type(WS), channel(WHITESPACE_CHANNEL);
COMMENT_B: COMMENT -> type(COMMENT), channel(COMMENTS_CHANNEL);

// ----------------------------------------------------------------------------

mode COMMAND_ARGS;

Atom: (RegularAtomPart | DoubleQuotedAtomPart | SingleQuotedAtomPart | ShellAtomPart)+;
fragment DoubleQuotedAtomPart: '"' (ShellAtomPart | ~('"' | '\\') | ('\\' .))* '"';
fragment SingleQuotedAtomPart: '\'' (~('\'' | '\\') | ('\\' .))* '\'';
fragment ShellAtomPart: '$(' (~([ \t\r\n\\"')]) | ('\\' .) | DoubleQuotedAtomPart | SingleQuotedAtomPart | ShellAtomPart | WS)+ ')';

fragment RegularAtomPart: ~([ \t\r\n\\"']) | EscapedAtomPart;
fragment EscapedAtomPart: ('\\' .) | (LC [ \t]*);

NL_C: NL -> type(NL), popMode;
WS_C: WS -> type(WS), channel(WHITESPACE_CHANNEL);
COMMENT_C: COMMENT -> type(COMMENT), channel(COMMENTS_CHANNEL);

// ----------------------------------------------------------------------------

mode COMMAND_ARGS_KEY_VALUE;

// Switch mode after '=' (may contain '=' as part of value after that).
EQUALS: '=' -> mode(COMMAND_ARGS_KEY_VALUE_ASSIGNMENT);

// Similar Atom, but don't allow '=' as part of it, unless it's in quotes.
Atom_CAKV: (RegularAtomPart_CAKV | DoubleQuotedAtomPart | SingleQuotedAtomPart | ShellAtomPart)+ -> type(Atom);
fragment RegularAtomPart_CAKV: ~([ \t\r\n"=\\]) | EscapedAtomPart;

NL_CAKV: NL -> type(NL), popMode;
WS_CAKV: WS -> type(WS), channel(WHITESPACE_CHANNEL);
COMMENT_CAKV: COMMENT -> type(COMMENT), channel(COMMENTS_CHANNEL);

// ----------------------------------------------------------------------------

mode COMMAND_ARGS_KEY_VALUE_ASSIGNMENT;

// Like COMMAND_ARGS, but include WS tokens so the whitespace
// gets added back to the value when we call 'GetText()' in the
// listener.

Atom_CAKVA: Atom -> type(Atom);
NL_CAKVA: NL -> type(NL), popMode;
WS_CAKVA: WS -> type(WS);
COMMENT_CAKVA: COMMENT -> type(COMMENT), channel(COMMENTS_CHANNEL);

// ----------------------------------------------------------------------------

mode COMMAND_ARGS_KEY_VALUE_LABEL;

EQUALS_L: '=' -> type(EQUALS);

// Similar Atom, but don't allow '=' as part of it, unless it's in quotes.
Atom_CAKVL: Atom_CAKV -> type(Atom);

NL_CAKVL: NL_CAKV -> type(NL), popMode;
WS_CAKVL: WS_CAKV -> type(WS), channel(WHITESPACE_CHANNEL);
COMMENT_CAKVL: COMMENT -> type(COMMENT), channel(COMMENTS_CHANNEL);
