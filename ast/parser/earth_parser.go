// Code generated from ast/parser/EarthParser.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // EarthParser

import (
	"fmt"
	"strconv"
  "sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}


type EarthParser struct {
	*antlr.BaseParser
}

var earthparserParserStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  literalNames           []string
  symbolicNames          []string
  ruleNames              []string
  predictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func earthparserParserInit() {
  staticData := &earthparserParserStaticData
  staticData.literalNames = []string{
    "", "", "", "", "", "'FROM'", "'FROM DOCKERFILE'", "'LOCALLY'", "'COPY'", 
    "'SAVE ARTIFACT'", "'SAVE IMAGE'", "'RUN'", "'EXPOSE'", "'VOLUME'", 
    "'ENV'", "'ARG'", "'LABEL'", "'BUILD'", "'WORKDIR'", "'USER'", "'CMD'", 
    "'ENTRYPOINT'", "'GIT CLONE'", "'ADD'", "'STOPSIGNAL'", "'ONBUILD'", 
    "'HEALTHCHECK'", "'SHELL'", "'DO'", "'COMMAND'", "'IMPORT'", "'VERSION'", 
    "'CACHE'", "'HOST'", "'PROJECT'", "'PIPELINE'", "'TRIGGER'", "'WITH'", 
    "", "", "", "", "", "", "", "", "'ELSE'", "'ELSE IF'", "'CATCH'", "'FINALLY'", 
    "'END'",
  }
  staticData.symbolicNames = []string{
    "", "INDENT", "DEDENT", "Target", "UserCommand", "FROM", "FROM_DOCKERFILE", 
    "LOCALLY", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE", "RUN", "EXPOSE", "VOLUME", 
    "ENV", "ARG", "LABEL", "BUILD", "WORKDIR", "USER", "CMD", "ENTRYPOINT", 
    "GIT_CLONE", "ADD", "STOPSIGNAL", "ONBUILD", "HEALTHCHECK", "SHELL", 
    "DO", "COMMAND", "IMPORT", "VERSION", "CACHE", "HOST", "PROJECT", "PIPELINE", 
    "TRIGGER", "WITH", "DOCKER", "IF", "TRY", "FOR", "WAIT", "NL", "WS", 
    "COMMENT", "ELSE", "ELSE_IF", "CATCH", "FINALLY", "END", "Atom", "EQUALS",
  }
  staticData.ruleNames = []string{
    "earthFile", "targets", "targetOrUserCommand", "target", "targetHeader", 
    "userCommand", "userCommandHeader", "stmts", "stmt", "commandStmt", 
    "version", "withStmt", "withBlock", "withExpr", "withCommand", "dockerCommand", 
    "ifStmt", "ifClause", "ifBlock", "elseIfClause", "elseIfBlock", "elseClause", 
    "elseBlock", "ifExpr", "elseIfExpr", "tryStmt", "tryClause", "tryBlock", 
    "catchClause", "catchBlock", "finallyClause", "finallyBlock", "forStmt", 
    "forClause", "forBlock", "forExpr", "waitStmt", "waitClause", "waitBlock", 
    "waitExpr", "fromStmt", "fromDockerfileStmt", "locallyStmt", "copyStmt", 
    "saveStmt", "saveImage", "saveArtifact", "runStmt", "buildStmt", "workdirStmt", 
    "userStmt", "cmdStmt", "entrypointStmt", "exposeStmt", "volumeStmt", 
    "envStmt", "argStmt", "optionalFlag", "envArgKey", "envArgValue", "labelStmt", 
    "labelKey", "labelValue", "gitCloneStmt", "addStmt", "stopsignalStmt", 
    "onbuildStmt", "healthcheckStmt", "shellStmt", "userCommandStmt", "doStmt", 
    "importStmt", "cacheStmt", "hostStmt", "projectStmt", "pipelineStmt", 
    "triggerStmt", "expr", "stmtWordsMaybeJSON", "stmtWords", "stmtWord",
  }
  staticData.predictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 52, 696, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
	4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 
	10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 
	2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2, 
	21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26, 
	7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7, 
	31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36, 
	2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2, 
	42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47, 
	7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7, 
	52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57, 
	2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7, 62, 2, 
	63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67, 2, 68, 
	7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7, 71, 2, 72, 7, 72, 2, 73, 7, 
	73, 2, 74, 7, 74, 2, 75, 7, 75, 2, 76, 7, 76, 2, 77, 7, 77, 2, 78, 7, 78, 
	2, 79, 7, 79, 2, 80, 7, 80, 1, 0, 5, 0, 164, 8, 0, 10, 0, 12, 0, 167, 9, 
	0, 1, 0, 3, 0, 170, 8, 0, 1, 0, 1, 0, 1, 0, 3, 0, 175, 8, 0, 1, 0, 5, 0, 
	178, 8, 0, 10, 0, 12, 0, 181, 9, 0, 1, 0, 3, 0, 184, 8, 0, 1, 0, 5, 0, 
	187, 8, 0, 10, 0, 12, 0, 190, 9, 0, 1, 0, 1, 0, 1, 1, 1, 1, 5, 1, 196, 
	8, 1, 10, 1, 12, 1, 199, 9, 1, 1, 1, 5, 1, 202, 8, 1, 10, 1, 12, 1, 205, 
	9, 1, 1, 2, 1, 2, 3, 2, 209, 8, 2, 1, 3, 1, 3, 4, 3, 213, 8, 3, 11, 3, 
	12, 3, 214, 1, 3, 1, 3, 5, 3, 219, 8, 3, 10, 3, 12, 3, 222, 9, 3, 1, 3, 
	3, 3, 225, 8, 3, 1, 3, 4, 3, 228, 8, 3, 11, 3, 12, 3, 229, 1, 3, 3, 3, 
	233, 8, 3, 1, 4, 1, 4, 1, 5, 1, 5, 4, 5, 239, 8, 5, 11, 5, 12, 5, 240, 
	1, 5, 1, 5, 5, 5, 245, 8, 5, 10, 5, 12, 5, 248, 9, 5, 1, 5, 1, 5, 4, 5, 
	252, 8, 5, 11, 5, 12, 5, 253, 1, 5, 1, 5, 3, 5, 258, 8, 5, 1, 6, 1, 6, 
	1, 7, 1, 7, 4, 7, 264, 8, 7, 11, 7, 12, 7, 265, 1, 7, 5, 7, 269, 8, 7, 
	10, 7, 12, 7, 272, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 280, 
	8, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 
	1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 
	1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 312, 8, 9, 1, 10, 1, 10, 
	1, 10, 4, 10, 317, 8, 10, 11, 10, 12, 10, 318, 1, 11, 1, 11, 4, 11, 323, 
	8, 11, 11, 11, 12, 11, 324, 1, 11, 3, 11, 328, 8, 11, 1, 11, 4, 11, 331, 
	8, 11, 11, 11, 12, 11, 332, 1, 11, 1, 11, 1, 12, 1, 12, 1, 13, 1, 13, 1, 
	13, 1, 14, 1, 14, 1, 15, 1, 15, 3, 15, 346, 8, 15, 1, 16, 1, 16, 4, 16, 
	350, 8, 16, 11, 16, 12, 16, 351, 1, 16, 5, 16, 355, 8, 16, 10, 16, 12, 
	16, 358, 9, 16, 1, 16, 4, 16, 361, 8, 16, 11, 16, 12, 16, 362, 1, 16, 3, 
	16, 366, 8, 16, 1, 16, 4, 16, 369, 8, 16, 11, 16, 12, 16, 370, 1, 16, 1, 
	16, 1, 17, 1, 17, 1, 17, 4, 17, 378, 8, 17, 11, 17, 12, 17, 379, 1, 17, 
	3, 17, 383, 8, 17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 4, 19, 390, 8, 19, 
	11, 19, 12, 19, 391, 1, 19, 3, 19, 395, 8, 19, 1, 20, 1, 20, 1, 21, 1, 
	21, 4, 21, 401, 8, 21, 11, 21, 12, 21, 402, 1, 21, 3, 21, 406, 8, 21, 1, 
	22, 1, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 4, 25, 416, 8, 25, 
	11, 25, 12, 25, 417, 1, 25, 3, 25, 421, 8, 25, 1, 25, 4, 25, 424, 8, 25, 
	11, 25, 12, 25, 425, 1, 25, 3, 25, 429, 8, 25, 1, 25, 4, 25, 432, 8, 25, 
	11, 25, 12, 25, 433, 1, 25, 1, 25, 1, 26, 1, 26, 4, 26, 440, 8, 26, 11, 
	26, 12, 26, 441, 1, 26, 3, 26, 445, 8, 26, 1, 27, 1, 27, 1, 28, 1, 28, 
	4, 28, 451, 8, 28, 11, 28, 12, 28, 452, 1, 28, 3, 28, 456, 8, 28, 1, 29, 
	1, 29, 1, 30, 1, 30, 4, 30, 462, 8, 30, 11, 30, 12, 30, 463, 1, 30, 3, 
	30, 467, 8, 30, 1, 31, 1, 31, 1, 32, 1, 32, 4, 32, 473, 8, 32, 11, 32, 
	12, 32, 474, 1, 32, 1, 32, 1, 33, 1, 33, 1, 33, 4, 33, 482, 8, 33, 11, 
	33, 12, 33, 483, 1, 33, 3, 33, 487, 8, 33, 1, 34, 1, 34, 1, 35, 1, 35, 
	1, 36, 1, 36, 4, 36, 495, 8, 36, 11, 36, 12, 36, 496, 1, 36, 1, 36, 1, 
	37, 1, 37, 3, 37, 503, 8, 37, 1, 37, 4, 37, 506, 8, 37, 11, 37, 12, 37, 
	507, 1, 37, 3, 37, 511, 8, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 1, 40, 
	3, 40, 519, 8, 40, 1, 41, 1, 41, 3, 41, 523, 8, 41, 1, 42, 1, 42, 3, 42, 
	527, 8, 42, 1, 43, 1, 43, 3, 43, 531, 8, 43, 1, 44, 1, 44, 3, 44, 535, 
	8, 44, 1, 45, 1, 45, 3, 45, 539, 8, 45, 1, 46, 1, 46, 3, 46, 543, 8, 46, 
	1, 47, 1, 47, 3, 47, 547, 8, 47, 1, 48, 1, 48, 3, 48, 551, 8, 48, 1, 49, 
	1, 49, 3, 49, 555, 8, 49, 1, 50, 1, 50, 3, 50, 559, 8, 50, 1, 51, 1, 51, 
	3, 51, 563, 8, 51, 1, 52, 1, 52, 3, 52, 567, 8, 52, 1, 53, 1, 53, 3, 53, 
	571, 8, 53, 1, 54, 1, 54, 3, 54, 575, 8, 54, 1, 55, 1, 55, 1, 55, 3, 55, 
	580, 8, 55, 1, 55, 3, 55, 583, 8, 55, 1, 55, 3, 55, 586, 8, 55, 1, 56, 
	1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 593, 8, 56, 1, 56, 3, 56, 596, 8, 56, 
	3, 56, 598, 8, 56, 1, 57, 3, 57, 601, 8, 57, 1, 58, 1, 58, 1, 59, 1, 59, 
	3, 59, 607, 8, 59, 1, 59, 5, 59, 610, 8, 59, 10, 59, 12, 59, 613, 9, 59, 
	1, 60, 1, 60, 1, 60, 1, 60, 1, 60, 5, 60, 620, 8, 60, 10, 60, 12, 60, 623, 
	9, 60, 1, 61, 1, 61, 1, 62, 1, 62, 1, 63, 1, 63, 3, 63, 631, 8, 63, 1, 
	64, 1, 64, 3, 64, 635, 8, 64, 1, 65, 1, 65, 3, 65, 639, 8, 65, 1, 66, 1, 
	66, 3, 66, 643, 8, 66, 1, 67, 1, 67, 3, 67, 647, 8, 67, 1, 68, 1, 68, 3, 
	68, 651, 8, 68, 1, 69, 1, 69, 3, 69, 655, 8, 69, 1, 70, 1, 70, 3, 70, 659, 
	8, 70, 1, 71, 1, 71, 3, 71, 663, 8, 71, 1, 72, 1, 72, 3, 72, 667, 8, 72, 
	1, 73, 1, 73, 3, 73, 671, 8, 73, 1, 74, 1, 74, 3, 74, 675, 8, 74, 1, 75, 
	1, 75, 3, 75, 679, 8, 75, 1, 76, 1, 76, 3, 76, 683, 8, 76, 1, 77, 1, 77, 
	1, 78, 1, 78, 1, 79, 4, 79, 690, 8, 79, 11, 79, 12, 79, 691, 1, 80, 1, 
	80, 1, 80, 0, 0, 81, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 
	28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 
	64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 
	100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 
	130, 132, 134, 136, 138, 140, 142, 144, 146, 148, 150, 152, 154, 156, 158, 
	160, 0, 0, 742, 0, 165, 1, 0, 0, 0, 2, 193, 1, 0, 0, 0, 4, 208, 1, 0, 0, 
	0, 6, 210, 1, 0, 0, 0, 8, 234, 1, 0, 0, 0, 10, 236, 1, 0, 0, 0, 12, 259, 
	1, 0, 0, 0, 14, 261, 1, 0, 0, 0, 16, 279, 1, 0, 0, 0, 18, 311, 1, 0, 0, 
	0, 20, 313, 1, 0, 0, 0, 22, 320, 1, 0, 0, 0, 24, 336, 1, 0, 0, 0, 26, 338, 
	1, 0, 0, 0, 28, 341, 1, 0, 0, 0, 30, 343, 1, 0, 0, 0, 32, 347, 1, 0, 0, 
	0, 34, 374, 1, 0, 0, 0, 36, 384, 1, 0, 0, 0, 38, 386, 1, 0, 0, 0, 40, 396, 
	1, 0, 0, 0, 42, 398, 1, 0, 0, 0, 44, 407, 1, 0, 0, 0, 46, 409, 1, 0, 0, 
	0, 48, 411, 1, 0, 0, 0, 50, 413, 1, 0, 0, 0, 52, 437, 1, 0, 0, 0, 54, 446, 
	1, 0, 0, 0, 56, 448, 1, 0, 0, 0, 58, 457, 1, 0, 0, 0, 60, 459, 1, 0, 0, 
	0, 62, 468, 1, 0, 0, 0, 64, 470, 1, 0, 0, 0, 66, 478, 1, 0, 0, 0, 68, 488, 
	1, 0, 0, 0, 70, 490, 1, 0, 0, 0, 72, 492, 1, 0, 0, 0, 74, 500, 1, 0, 0, 
	0, 76, 512, 1, 0, 0, 0, 78, 514, 1, 0, 0, 0, 80, 516, 1, 0, 0, 0, 82, 520, 
	1, 0, 0, 0, 84, 524, 1, 0, 0, 0, 86, 528, 1, 0, 0, 0, 88, 534, 1, 0, 0, 
	0, 90, 536, 1, 0, 0, 0, 92, 540, 1, 0, 0, 0, 94, 544, 1, 0, 0, 0, 96, 548, 
	1, 0, 0, 0, 98, 552, 1, 0, 0, 0, 100, 556, 1, 0, 0, 0, 102, 560, 1, 0, 
	0, 0, 104, 564, 1, 0, 0, 0, 106, 568, 1, 0, 0, 0, 108, 572, 1, 0, 0, 0, 
	110, 576, 1, 0, 0, 0, 112, 587, 1, 0, 0, 0, 114, 600, 1, 0, 0, 0, 116, 
	602, 1, 0, 0, 0, 118, 604, 1, 0, 0, 0, 120, 614, 1, 0, 0, 0, 122, 624, 
	1, 0, 0, 0, 124, 626, 1, 0, 0, 0, 126, 628, 1, 0, 0, 0, 128, 632, 1, 0, 
	0, 0, 130, 636, 1, 0, 0, 0, 132, 640, 1, 0, 0, 0, 134, 644, 1, 0, 0, 0, 
	136, 648, 1, 0, 0, 0, 138, 652, 1, 0, 0, 0, 140, 656, 1, 0, 0, 0, 142, 
	660, 1, 0, 0, 0, 144, 664, 1, 0, 0, 0, 146, 668, 1, 0, 0, 0, 148, 672, 
	1, 0, 0, 0, 150, 676, 1, 0, 0, 0, 152, 680, 1, 0, 0, 0, 154, 684, 1, 0, 
	0, 0, 156, 686, 1, 0, 0, 0, 158, 689, 1, 0, 0, 0, 160, 693, 1, 0, 0, 0, 
	162, 164, 5, 43, 0, 0, 163, 162, 1, 0, 0, 0, 164, 167, 1, 0, 0, 0, 165, 
	163, 1, 0, 0, 0, 165, 166, 1, 0, 0, 0, 166, 169, 1, 0, 0, 0, 167, 165, 
	1, 0, 0, 0, 168, 170, 3, 20, 10, 0, 169, 168, 1, 0, 0, 0, 169, 170, 1, 
	0, 0, 0, 170, 174, 1, 0, 0, 0, 171, 172, 3, 14, 7, 0, 172, 173, 5, 43, 
	0, 0, 173, 175, 1, 0, 0, 0, 174, 171, 1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 
	175, 179, 1, 0, 0, 0, 176, 178, 5, 43, 0, 0, 177, 176, 1, 0, 0, 0, 178, 
	181, 1, 0, 0, 0, 179, 177, 1, 0, 0, 0, 179, 180, 1, 0, 0, 0, 180, 183, 
	1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 182, 184, 3, 2, 1, 0, 183, 182, 1, 0, 
	0, 0, 183, 184, 1, 0, 0, 0, 184, 188, 1, 0, 0, 0, 185, 187, 5, 43, 0, 0, 
	186, 185, 1, 0, 0, 0, 187, 190, 1, 0, 0, 0, 188, 186, 1, 0, 0, 0, 188, 
	189, 1, 0, 0, 0, 189, 191, 1, 0, 0, 0, 190, 188, 1, 0, 0, 0, 191, 192, 
	5, 0, 0, 1, 192, 1, 1, 0, 0, 0, 193, 203, 3, 4, 2, 0, 194, 196, 5, 43, 
	0, 0, 195, 194, 1, 0, 0, 0, 196, 199, 1, 0, 0, 0, 197, 195, 1, 0, 0, 0, 
	197, 198, 1, 0, 0, 0, 198, 200, 1, 0, 0, 0, 199, 197, 1, 0, 0, 0, 200, 
	202, 3, 4, 2, 0, 201, 197, 1, 0, 0, 0, 202, 205, 1, 0, 0, 0, 203, 201, 
	1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 3, 1, 0, 0, 0, 205, 203, 1, 0, 0, 
	0, 206, 209, 3, 6, 3, 0, 207, 209, 3, 10, 5, 0, 208, 206, 1, 0, 0, 0, 208, 
	207, 1, 0, 0, 0, 209, 5, 1, 0, 0, 0, 210, 212, 3, 8, 4, 0, 211, 213, 5, 
	43, 0, 0, 212, 211, 1, 0, 0, 0, 213, 214, 1, 0, 0, 0, 214, 212, 1, 0, 0, 
	0, 214, 215, 1, 0, 0, 0, 215, 232, 1, 0, 0, 0, 216, 220, 5, 1, 0, 0, 217, 
	219, 5, 43, 0, 0, 218, 217, 1, 0, 0, 0, 219, 222, 1, 0, 0, 0, 220, 218, 
	1, 0, 0, 0, 220, 221, 1, 0, 0, 0, 221, 224, 1, 0, 0, 0, 222, 220, 1, 0, 
	0, 0, 223, 225, 3, 14, 7, 0, 224, 223, 1, 0, 0, 0, 224, 225, 1, 0, 0, 0, 
	225, 227, 1, 0, 0, 0, 226, 228, 5, 43, 0, 0, 227, 226, 1, 0, 0, 0, 228, 
	229, 1, 0, 0, 0, 229, 227, 1, 0, 0, 0, 229, 230, 1, 0, 0, 0, 230, 231, 
	1, 0, 0, 0, 231, 233, 5, 2, 0, 0, 232, 216, 1, 0, 0, 0, 232, 233, 1, 0, 
	0, 0, 233, 7, 1, 0, 0, 0, 234, 235, 5, 3, 0, 0, 235, 9, 1, 0, 0, 0, 236, 
	238, 3, 12, 6, 0, 237, 239, 5, 43, 0, 0, 238, 237, 1, 0, 0, 0, 239, 240, 
	1, 0, 0, 0, 240, 238, 1, 0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 257, 1, 0, 
	0, 0, 242, 246, 5, 1, 0, 0, 243, 245, 5, 43, 0, 0, 244, 243, 1, 0, 0, 0, 
	245, 248, 1, 0, 0, 0, 246, 244, 1, 0, 0, 0, 246, 247, 1, 0, 0, 0, 247, 
	249, 1, 0, 0, 0, 248, 246, 1, 0, 0, 0, 249, 251, 3, 14, 7, 0, 250, 252, 
	5, 43, 0, 0, 251, 250, 1, 0, 0, 0, 252, 253, 1, 0, 0, 0, 253, 251, 1, 0, 
	0, 0, 253, 254, 1, 0, 0, 0, 254, 255, 1, 0, 0, 0, 255, 256, 5, 2, 0, 0, 
	256, 258, 1, 0, 0, 0, 257, 242, 1, 0, 0, 0, 257, 258, 1, 0, 0, 0, 258, 
	11, 1, 0, 0, 0, 259, 260, 5, 4, 0, 0, 260, 13, 1, 0, 0, 0, 261, 270, 3, 
	16, 8, 0, 262, 264, 5, 43, 0, 0, 263, 262, 1, 0, 0, 0, 264, 265, 1, 0, 
	0, 0, 265, 263, 1, 0, 0, 0, 265, 266, 1, 0, 0, 0, 266, 267, 1, 0, 0, 0, 
	267, 269, 3, 16, 8, 0, 268, 263, 1, 0, 0, 0, 269, 272, 1, 0, 0, 0, 270, 
	268, 1, 0, 0, 0, 270, 271, 1, 0, 0, 0, 271, 15, 1, 0, 0, 0, 272, 270, 1, 
	0, 0, 0, 273, 280, 3, 18, 9, 0, 274, 280, 3, 22, 11, 0, 275, 280, 3, 32, 
	16, 0, 276, 280, 3, 64, 32, 0, 277, 280, 3, 72, 36, 0, 278, 280, 3, 50, 
	25, 0, 279, 273, 1, 0, 0, 0, 279, 274, 1, 0, 0, 0, 279, 275, 1, 0, 0, 0, 
	279, 276, 1, 0, 0, 0, 279, 277, 1, 0, 0, 0, 279, 278, 1, 0, 0, 0, 280, 
	17, 1, 0, 0, 0, 281, 312, 3, 80, 40, 0, 282, 312, 3, 82, 41, 0, 283, 312, 
	3, 84, 42, 0, 284, 312, 3, 86, 43, 0, 285, 312, 3, 88, 44, 0, 286, 312, 
	3, 94, 47, 0, 287, 312, 3, 96, 48, 0, 288, 312, 3, 98, 49, 0, 289, 312, 
	3, 100, 50, 0, 290, 312, 3, 102, 51, 0, 291, 312, 3, 104, 52, 0, 292, 312, 
	3, 106, 53, 0, 293, 312, 3, 108, 54, 0, 294, 312, 3, 110, 55, 0, 295, 312, 
	3, 112, 56, 0, 296, 312, 3, 120, 60, 0, 297, 312, 3, 126, 63, 0, 298, 312, 
	3, 128, 64, 0, 299, 312, 3, 130, 65, 0, 300, 312, 3, 132, 66, 0, 301, 312, 
	3, 134, 67, 0, 302, 312, 3, 136, 68, 0, 303, 312, 3, 138, 69, 0, 304, 312, 
	3, 140, 70, 0, 305, 312, 3, 142, 71, 0, 306, 312, 3, 144, 72, 0, 307, 312, 
	3, 146, 73, 0, 308, 312, 3, 148, 74, 0, 309, 312, 3, 150, 75, 0, 310, 312, 
	3, 152, 76, 0, 311, 281, 1, 0, 0, 0, 311, 282, 1, 0, 0, 0, 311, 283, 1, 
	0, 0, 0, 311, 284, 1, 0, 0, 0, 311, 285, 1, 0, 0, 0, 311, 286, 1, 0, 0, 
	0, 311, 287, 1, 0, 0, 0, 311, 288, 1, 0, 0, 0, 311, 289, 1, 0, 0, 0, 311, 
	290, 1, 0, 0, 0, 311, 291, 1, 0, 0, 0, 311, 292, 1, 0, 0, 0, 311, 293, 
	1, 0, 0, 0, 311, 294, 1, 0, 0, 0, 311, 295, 1, 0, 0, 0, 311, 296, 1, 0, 
	0, 0, 311, 297, 1, 0, 0, 0, 311, 298, 1, 0, 0, 0, 311, 299, 1, 0, 0, 0, 
	311, 300, 1, 0, 0, 0, 311, 301, 1, 0, 0, 0, 311, 302, 1, 0, 0, 0, 311, 
	303, 1, 0, 0, 0, 311, 304, 1, 0, 0, 0, 311, 305, 1, 0, 0, 0, 311, 306, 
	1, 0, 0, 0, 311, 307, 1, 0, 0, 0, 311, 308, 1, 0, 0, 0, 311, 309, 1, 0, 
	0, 0, 311, 310, 1, 0, 0, 0, 312, 19, 1, 0, 0, 0, 313, 314, 5, 31, 0, 0, 
	314, 316, 3, 158, 79, 0, 315, 317, 5, 43, 0, 0, 316, 315, 1, 0, 0, 0, 317, 
	318, 1, 0, 0, 0, 318, 316, 1, 0, 0, 0, 318, 319, 1, 0, 0, 0, 319, 21, 1, 
	0, 0, 0, 320, 327, 3, 26, 13, 0, 321, 323, 5, 43, 0, 0, 322, 321, 1, 0, 
	0, 0, 323, 324, 1, 0, 0, 0, 324, 322, 1, 0, 0, 0, 324, 325, 1, 0, 0, 0, 
	325, 326, 1, 0, 0, 0, 326, 328, 3, 24, 12, 0, 327, 322, 1, 0, 0, 0, 327, 
	328, 1, 0, 0, 0, 328, 330, 1, 0, 0, 0, 329, 331, 5, 43, 0, 0, 330, 329, 
	1, 0, 0, 0, 331, 332, 1, 0, 0, 0, 332, 330, 1, 0, 0, 0, 332, 333, 1, 0, 
	0, 0, 333, 334, 1, 0, 0, 0, 334, 335, 5, 50, 0, 0, 335, 23, 1, 0, 0, 0, 
	336, 337, 3, 14, 7, 0, 337, 25, 1, 0, 0, 0, 338, 339, 5, 37, 0, 0, 339, 
	340, 3, 28, 14, 0, 340, 27, 1, 0, 0, 0, 341, 342, 3, 30, 15, 0, 342, 29, 
	1, 0, 0, 0, 343, 345, 5, 38, 0, 0, 344, 346, 3, 158, 79, 0, 345, 344, 1, 
	0, 0, 0, 345, 346, 1, 0, 0, 0, 346, 31, 1, 0, 0, 0, 347, 356, 3, 34, 17, 
	0, 348, 350, 5, 43, 0, 0, 349, 348, 1, 0, 0, 0, 350, 351, 1, 0, 0, 0, 351, 
	349, 1, 0, 0, 0, 351, 352, 1, 0, 0, 0, 352, 353, 1, 0, 0, 0, 353, 355, 
	3, 38, 19, 0, 354, 349, 1, 0, 0, 0, 355, 358, 1, 0, 0, 0, 356, 354, 1, 
	0, 0, 0, 356, 357, 1, 0, 0, 0, 357, 365, 1, 0, 0, 0, 358, 356, 1, 0, 0, 
	0, 359, 361, 5, 43, 0, 0, 360, 359, 1, 0, 0, 0, 361, 362, 1, 0, 0, 0, 362, 
	360, 1, 0, 0, 0, 362, 363, 1, 0, 0, 0, 363, 364, 1, 0, 0, 0, 364, 366, 
	3, 42, 21, 0, 365, 360, 1, 0, 0, 0, 365, 366, 1, 0, 0, 0, 366, 368, 1, 
	0, 0, 0, 367, 369, 5, 43, 0, 0, 368, 367, 1, 0, 0, 0, 369, 370, 1, 0, 0, 
	0, 370, 368, 1, 0, 0, 0, 370, 371, 1, 0, 0, 0, 371, 372, 1, 0, 0, 0, 372, 
	373, 5, 50, 0, 0, 373, 33, 1, 0, 0, 0, 374, 375, 5, 39, 0, 0, 375, 382, 
	3, 46, 23, 0, 376, 378, 5, 43, 0, 0, 377, 376, 1, 0, 0, 0, 378, 379, 1, 
	0, 0, 0, 379, 377, 1, 0, 0, 0, 379, 380, 1, 0, 0, 0, 380, 381, 1, 0, 0, 
	0, 381, 383, 3, 36, 18, 0, 382, 377, 1, 0, 0, 0, 382, 383, 1, 0, 0, 0, 
	383, 35, 1, 0, 0, 0, 384, 385, 3, 14, 7, 0, 385, 37, 1, 0, 0, 0, 386, 387, 
	5, 47, 0, 0, 387, 394, 3, 48, 24, 0, 388, 390, 5, 43, 0, 0, 389, 388, 1, 
	0, 0, 0, 390, 391, 1, 0, 0, 0, 391, 389, 1, 0, 0, 0, 391, 392, 1, 0, 0, 
	0, 392, 393, 1, 0, 0, 0, 393, 395, 3, 40, 20, 0, 394, 389, 1, 0, 0, 0, 
	394, 395, 1, 0, 0, 0, 395, 39, 1, 0, 0, 0, 396, 397, 3, 14, 7, 0, 397, 
	41, 1, 0, 0, 0, 398, 405, 5, 46, 0, 0, 399, 401, 5, 43, 0, 0, 400, 399, 
	1, 0, 0, 0, 401, 402, 1, 0, 0, 0, 402, 400, 1, 0, 0, 0, 402, 403, 1, 0, 
	0, 0, 403, 404, 1, 0, 0, 0, 404, 406, 3, 44, 22, 0, 405, 400, 1, 0, 0, 
	0, 405, 406, 1, 0, 0, 0, 406, 43, 1, 0, 0, 0, 407, 408, 3, 14, 7, 0, 408, 
	45, 1, 0, 0, 0, 409, 410, 3, 154, 77, 0, 410, 47, 1, 0, 0, 0, 411, 412, 
	3, 154, 77, 0, 412, 49, 1, 0, 0, 0, 413, 420, 3, 52, 26, 0, 414, 416, 5, 
	43, 0, 0, 415, 414, 1, 0, 0, 0, 416, 417, 1, 0, 0, 0, 417, 415, 1, 0, 0, 
	0, 417, 418, 1, 0, 0, 0, 418, 419, 1, 0, 0, 0, 419, 421, 3, 56, 28, 0, 
	420, 415, 1, 0, 0, 0, 420, 421, 1, 0, 0, 0, 421, 428, 1, 0, 0, 0, 422, 
	424, 5, 43, 0, 0, 423, 422, 1, 0, 0, 0, 424, 425, 1, 0, 0, 0, 425, 423, 
	1, 0, 0, 0, 425, 426, 1, 0, 0, 0, 426, 427, 1, 0, 0, 0, 427, 429, 3, 60, 
	30, 0, 428, 423, 1, 0, 0, 0, 428, 429, 1, 0, 0, 0, 429, 431, 1, 0, 0, 0, 
	430, 432, 5, 43, 0, 0, 431, 430, 1, 0, 0, 0, 432, 433, 1, 0, 0, 0, 433, 
	431, 1, 0, 0, 0, 433, 434, 1, 0, 0, 0, 434, 435, 1, 0, 0, 0, 435, 436, 
	5, 50, 0, 0, 436, 51, 1, 0, 0, 0, 437, 444, 5, 40, 0, 0, 438, 440, 5, 43, 
	0, 0, 439, 438, 1, 0, 0, 0, 440, 441, 1, 0, 0, 0, 441, 439, 1, 0, 0, 0, 
	441, 442, 1, 0, 0, 0, 442, 443, 1, 0, 0, 0, 443, 445, 3, 54, 27, 0, 444, 
	439, 1, 0, 0, 0, 444, 445, 1, 0, 0, 0, 445, 53, 1, 0, 0, 0, 446, 447, 3, 
	14, 7, 0, 447, 55, 1, 0, 0, 0, 448, 455, 5, 48, 0, 0, 449, 451, 5, 43, 
	0, 0, 450, 449, 1, 0, 0, 0, 451, 452, 1, 0, 0, 0, 452, 450, 1, 0, 0, 0, 
	452, 453, 1, 0, 0, 0, 453, 454, 1, 0, 0, 0, 454, 456, 3, 58, 29, 0, 455, 
	450, 1, 0, 0, 0, 455, 456, 1, 0, 0, 0, 456, 57, 1, 0, 0, 0, 457, 458, 3, 
	14, 7, 0, 458, 59, 1, 0, 0, 0, 459, 466, 5, 49, 0, 0, 460, 462, 5, 43, 
	0, 0, 461, 460, 1, 0, 0, 0, 462, 463, 1, 0, 0, 0, 463, 461, 1, 0, 0, 0, 
	463, 464, 1, 0, 0, 0, 464, 465, 1, 0, 0, 0, 465, 467, 3, 62, 31, 0, 466, 
	461, 1, 0, 0, 0, 466, 467, 1, 0, 0, 0, 467, 61, 1, 0, 0, 0, 468, 469, 3, 
	14, 7, 0, 469, 63, 1, 0, 0, 0, 470, 472, 3, 66, 33, 0, 471, 473, 5, 43, 
	0, 0, 472, 471, 1, 0, 0, 0, 473, 474, 1, 0, 0, 0, 474, 472, 1, 0, 0, 0, 
	474, 475, 1, 0, 0, 0, 475, 476, 1, 0, 0, 0, 476, 477, 5, 50, 0, 0, 477, 
	65, 1, 0, 0, 0, 478, 479, 5, 41, 0, 0, 479, 486, 3, 70, 35, 0, 480, 482, 
	5, 43, 0, 0, 481, 480, 1, 0, 0, 0, 482, 483, 1, 0, 0, 0, 483, 481, 1, 0, 
	0, 0, 483, 484, 1, 0, 0, 0, 484, 485, 1, 0, 0, 0, 485, 487, 3, 68, 34, 
	0, 486, 481, 1, 0, 0, 0, 486, 487, 1, 0, 0, 0, 487, 67, 1, 0, 0, 0, 488, 
	489, 3, 14, 7, 0, 489, 69, 1, 0, 0, 0, 490, 491, 3, 158, 79, 0, 491, 71, 
	1, 0, 0, 0, 492, 494, 3, 74, 37, 0, 493, 495, 5, 43, 0, 0, 494, 493, 1, 
	0, 0, 0, 495, 496, 1, 0, 0, 0, 496, 494, 1, 0, 0, 0, 496, 497, 1, 0, 0, 
	0, 497, 498, 1, 0, 0, 0, 498, 499, 5, 50, 0, 0, 499, 73, 1, 0, 0, 0, 500, 
	502, 5, 42, 0, 0, 501, 503, 3, 78, 39, 0, 502, 501, 1, 0, 0, 0, 502, 503, 
	1, 0, 0, 0, 503, 510, 1, 0, 0, 0, 504, 506, 5, 43, 0, 0, 505, 504, 1, 0, 
	0, 0, 506, 507, 1, 0, 0, 0, 507, 505, 1, 0, 0, 0, 507, 508, 1, 0, 0, 0, 
	508, 509, 1, 0, 0, 0, 509, 511, 3, 76, 38, 0, 510, 505, 1, 0, 0, 0, 510, 
	511, 1, 0, 0, 0, 511, 75, 1, 0, 0, 0, 512, 513, 3, 14, 7, 0, 513, 77, 1, 
	0, 0, 0, 514, 515, 3, 158, 79, 0, 515, 79, 1, 0, 0, 0, 516, 518, 5, 5, 
	0, 0, 517, 519, 3, 158, 79, 0, 518, 517, 1, 0, 0, 0, 518, 519, 1, 0, 0, 
	0, 519, 81, 1, 0, 0, 0, 520, 522, 5, 6, 0, 0, 521, 523, 3, 158, 79, 0, 
	522, 521, 1, 0, 0, 0, 522, 523, 1, 0, 0, 0, 523, 83, 1, 0, 0, 0, 524, 526, 
	5, 7, 0, 0, 525, 527, 3, 158, 79, 0, 526, 525, 1, 0, 0, 0, 526, 527, 1, 
	0, 0, 0, 527, 85, 1, 0, 0, 0, 528, 530, 5, 8, 0, 0, 529, 531, 3, 158, 79, 
	0, 530, 529, 1, 0, 0, 0, 530, 531, 1, 0, 0, 0, 531, 87, 1, 0, 0, 0, 532, 
	535, 3, 92, 46, 0, 533, 535, 3, 90, 45, 0, 534, 532, 1, 0, 0, 0, 534, 533, 
	1, 0, 0, 0, 535, 89, 1, 0, 0, 0, 536, 538, 5, 10, 0, 0, 537, 539, 3, 158, 
	79, 0, 538, 537, 1, 0, 0, 0, 538, 539, 1, 0, 0, 0, 539, 91, 1, 0, 0, 0, 
	540, 542, 5, 9, 0, 0, 541, 543, 3, 158, 79, 0, 542, 541, 1, 0, 0, 0, 542, 
	543, 1, 0, 0, 0, 543, 93, 1, 0, 0, 0, 544, 546, 5, 11, 0, 0, 545, 547, 
	3, 156, 78, 0, 546, 545, 1, 0, 0, 0, 546, 547, 1, 0, 0, 0, 547, 95, 1, 
	0, 0, 0, 548, 550, 5, 17, 0, 0, 549, 551, 3, 158, 79, 0, 550, 549, 1, 0, 
	0, 0, 550, 551, 1, 0, 0, 0, 551, 97, 1, 0, 0, 0, 552, 554, 5, 18, 0, 0, 
	553, 555, 3, 158, 79, 0, 554, 553, 1, 0, 0, 0, 554, 555, 1, 0, 0, 0, 555, 
	99, 1, 0, 0, 0, 556, 558, 5, 19, 0, 0, 557, 559, 3, 158, 79, 0, 558, 557, 
	1, 0, 0, 0, 558, 559, 1, 0, 0, 0, 559, 101, 1, 0, 0, 0, 560, 562, 5, 20, 
	0, 0, 561, 563, 3, 156, 78, 0, 562, 561, 1, 0, 0, 0, 562, 563, 1, 0, 0, 
	0, 563, 103, 1, 0, 0, 0, 564, 566, 5, 21, 0, 0, 565, 567, 3, 156, 78, 0, 
	566, 565, 1, 0, 0, 0, 566, 567, 1, 0, 0, 0, 567, 105, 1, 0, 0, 0, 568, 
	570, 5, 12, 0, 0, 569, 571, 3, 158, 79, 0, 570, 569, 1, 0, 0, 0, 570, 571, 
	1, 0, 0, 0, 571, 107, 1, 0, 0, 0, 572, 574, 5, 13, 0, 0, 573, 575, 3, 156, 
	78, 0, 574, 573, 1, 0, 0, 0, 574, 575, 1, 0, 0, 0, 575, 109, 1, 0, 0, 0, 
	576, 577, 5, 14, 0, 0, 577, 579, 3, 116, 58, 0, 578, 580, 5, 52, 0, 0, 
	579, 578, 1, 0, 0, 0, 579, 580, 1, 0, 0, 0, 580, 585, 1, 0, 0, 0, 581, 
	583, 5, 44, 0, 0, 582, 581, 1, 0, 0, 0, 582, 583, 1, 0, 0, 0, 583, 584, 
	1, 0, 0, 0, 584, 586, 3, 118, 59, 0, 585, 582, 1, 0, 0, 0, 585, 586, 1, 
	0, 0, 0, 586, 111, 1, 0, 0, 0, 587, 588, 5, 15, 0, 0, 588, 589, 3, 114, 
	57, 0, 589, 597, 3, 116, 58, 0, 590, 595, 5, 52, 0, 0, 591, 593, 5, 44, 
	0, 0, 592, 591, 1, 0, 0, 0, 592, 593, 1, 0, 0, 0, 593, 594, 1, 0, 0, 0, 
	594, 596, 3, 118, 59, 0, 595, 592, 1, 0, 0, 0, 595, 596, 1, 0, 0, 0, 596, 
	598, 1, 0, 0, 0, 597, 590, 1, 0, 0, 0, 597, 598, 1, 0, 0, 0, 598, 113, 
	1, 0, 0, 0, 599, 601, 3, 158, 79, 0, 600, 599, 1, 0, 0, 0, 600, 601, 1, 
	0, 0, 0, 601, 115, 1, 0, 0, 0, 602, 603, 5, 51, 0, 0, 603, 117, 1, 0, 0, 
	0, 604, 611, 5, 51, 0, 0, 605, 607, 5, 44, 0, 0, 606, 605, 1, 0, 0, 0, 
	606, 607, 1, 0, 0, 0, 607, 608, 1, 0, 0, 0, 608, 610, 5, 51, 0, 0, 609, 
	606, 1, 0, 0, 0, 610, 613, 1, 0, 0, 0, 611, 609, 1, 0, 0, 0, 611, 612, 
	1, 0, 0, 0, 612, 119, 1, 0, 0, 0, 613, 611, 1, 0, 0, 0, 614, 621, 5, 16, 
	0, 0, 615, 616, 3, 122, 61, 0, 616, 617, 5, 52, 0, 0, 617, 618, 3, 124, 
	62, 0, 618, 620, 1, 0, 0, 0, 619, 615, 1, 0, 0, 0, 620, 623, 1, 0, 0, 0, 
	621, 619, 1, 0, 0, 0, 621, 622, 1, 0, 0, 0, 622, 121, 1, 0, 0, 0, 623, 
	621, 1, 0, 0, 0, 624, 625, 5, 51, 0, 0, 625, 123, 1, 0, 0, 0, 626, 627, 
	5, 51, 0, 0, 627, 125, 1, 0, 0, 0, 628, 630, 5, 22, 0, 0, 629, 631, 3, 
	158, 79, 0, 630, 629, 1, 0, 0, 0, 630, 631, 1, 0, 0, 0, 631, 127, 1, 0, 
	0, 0, 632, 634, 5, 23, 0, 0, 633, 635, 3, 158, 79, 0, 634, 633, 1, 0, 0, 
	0, 634, 635, 1, 0, 0, 0, 635, 129, 1, 0, 0, 0, 636, 638, 5, 24, 0, 0, 637, 
	639, 3, 158, 79, 0, 638, 637, 1, 0, 0, 0, 638, 639, 1, 0, 0, 0, 639, 131, 
	1, 0, 0, 0, 640, 642, 5, 25, 0, 0, 641, 643, 3, 158, 79, 0, 642, 641, 1, 
	0, 0, 0, 642, 643, 1, 0, 0, 0, 643, 133, 1, 0, 0, 0, 644, 646, 5, 26, 0, 
	0, 645, 647, 3, 158, 79, 0, 646, 645, 1, 0, 0, 0, 646, 647, 1, 0, 0, 0, 
	647, 135, 1, 0, 0, 0, 648, 650, 5, 27, 0, 0, 649, 651, 3, 158, 79, 0, 650, 
	649, 1, 0, 0, 0, 650, 651, 1, 0, 0, 0, 651, 137, 1, 0, 0, 0, 652, 654, 
	5, 29, 0, 0, 653, 655, 3, 158, 79, 0, 654, 653, 1, 0, 0, 0, 654, 655, 1, 
	0, 0, 0, 655, 139, 1, 0, 0, 0, 656, 658, 5, 28, 0, 0, 657, 659, 3, 158, 
	79, 0, 658, 657, 1, 0, 0, 0, 658, 659, 1, 0, 0, 0, 659, 141, 1, 0, 0, 0, 
	660, 662, 5, 30, 0, 0, 661, 663, 3, 158, 79, 0, 662, 661, 1, 0, 0, 0, 662, 
	663, 1, 0, 0, 0, 663, 143, 1, 0, 0, 0, 664, 666, 5, 32, 0, 0, 665, 667, 
	3, 158, 79, 0, 666, 665, 1, 0, 0, 0, 666, 667, 1, 0, 0, 0, 667, 145, 1, 
	0, 0, 0, 668, 670, 5, 33, 0, 0, 669, 671, 3, 158, 79, 0, 670, 669, 1, 0, 
	0, 0, 670, 671, 1, 0, 0, 0, 671, 147, 1, 0, 0, 0, 672, 674, 5, 34, 0, 0, 
	673, 675, 3, 158, 79, 0, 674, 673, 1, 0, 0, 0, 674, 675, 1, 0, 0, 0, 675, 
	149, 1, 0, 0, 0, 676, 678, 5, 35, 0, 0, 677, 679, 3, 158, 79, 0, 678, 677, 
	1, 0, 0, 0, 678, 679, 1, 0, 0, 0, 679, 151, 1, 0, 0, 0, 680, 682, 5, 36, 
	0, 0, 681, 683, 3, 158, 79, 0, 682, 681, 1, 0, 0, 0, 682, 683, 1, 0, 0, 
	0, 683, 153, 1, 0, 0, 0, 684, 685, 3, 156, 78, 0, 685, 155, 1, 0, 0, 0, 
	686, 687, 3, 158, 79, 0, 687, 157, 1, 0, 0, 0, 688, 690, 3, 160, 80, 0, 
	689, 688, 1, 0, 0, 0, 690, 691, 1, 0, 0, 0, 691, 689, 1, 0, 0, 0, 691, 
	692, 1, 0, 0, 0, 692, 159, 1, 0, 0, 0, 693, 694, 5, 51, 0, 0, 694, 161, 
	1, 0, 0, 0, 96, 165, 169, 174, 179, 183, 188, 197, 203, 208, 214, 220, 
	224, 229, 232, 240, 246, 253, 257, 265, 270, 279, 311, 318, 324, 327, 332, 
	345, 351, 356, 362, 365, 370, 379, 382, 391, 394, 402, 405, 417, 420, 425, 
	428, 433, 441, 444, 452, 455, 463, 466, 474, 483, 486, 496, 502, 507, 510, 
	518, 522, 526, 530, 534, 538, 542, 546, 550, 554, 558, 562, 566, 570, 574, 
	579, 582, 585, 592, 595, 597, 600, 606, 611, 621, 630, 634, 638, 642, 646, 
	650, 654, 658, 662, 666, 670, 674, 678, 682, 691,
}
  deserializer := antlr.NewATNDeserializer(nil)
  staticData.atn = deserializer.Deserialize(staticData.serializedATN)
  atn := staticData.atn
  staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
  decisionToDFA := staticData.decisionToDFA
  for index, state := range atn.DecisionToState {
    decisionToDFA[index] = antlr.NewDFA(state, index)
  }
}

// EarthParserInit initializes any static state used to implement EarthParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewEarthParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EarthParserInit() {
  staticData := &earthparserParserStaticData
  staticData.once.Do(earthparserParserInit)
}

// NewEarthParser produces a new parser instance for the optional input antlr.TokenStream.
func NewEarthParser(input antlr.TokenStream) *EarthParser {
	EarthParserInit()
	this := new(EarthParser)
	this.BaseParser = antlr.NewBaseParser(input)
  staticData := &earthparserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "EarthParser.g4"

	return this
}


// EarthParser tokens.
const (
	EarthParserEOF = antlr.TokenEOF
	EarthParserINDENT = 1
	EarthParserDEDENT = 2
	EarthParserTarget = 3
	EarthParserUserCommand = 4
	EarthParserFROM = 5
	EarthParserFROM_DOCKERFILE = 6
	EarthParserLOCALLY = 7
	EarthParserCOPY = 8
	EarthParserSAVE_ARTIFACT = 9
	EarthParserSAVE_IMAGE = 10
	EarthParserRUN = 11
	EarthParserEXPOSE = 12
	EarthParserVOLUME = 13
	EarthParserENV = 14
	EarthParserARG = 15
	EarthParserLABEL = 16
	EarthParserBUILD = 17
	EarthParserWORKDIR = 18
	EarthParserUSER = 19
	EarthParserCMD = 20
	EarthParserENTRYPOINT = 21
	EarthParserGIT_CLONE = 22
	EarthParserADD = 23
	EarthParserSTOPSIGNAL = 24
	EarthParserONBUILD = 25
	EarthParserHEALTHCHECK = 26
	EarthParserSHELL = 27
	EarthParserDO = 28
	EarthParserCOMMAND = 29
	EarthParserIMPORT = 30
	EarthParserVERSION = 31
	EarthParserCACHE = 32
	EarthParserHOST = 33
	EarthParserPROJECT = 34
	EarthParserPIPELINE = 35
	EarthParserTRIGGER = 36
	EarthParserWITH = 37
	EarthParserDOCKER = 38
	EarthParserIF = 39
	EarthParserTRY = 40
	EarthParserFOR = 41
	EarthParserWAIT = 42
	EarthParserNL = 43
	EarthParserWS = 44
	EarthParserCOMMENT = 45
	EarthParserELSE = 46
	EarthParserELSE_IF = 47
	EarthParserCATCH = 48
	EarthParserFINALLY = 49
	EarthParserEND = 50
	EarthParserAtom = 51
	EarthParserEQUALS = 52
)

// EarthParser rules.
const (
	EarthParserRULE_earthFile = 0
	EarthParserRULE_targets = 1
	EarthParserRULE_targetOrUserCommand = 2
	EarthParserRULE_target = 3
	EarthParserRULE_targetHeader = 4
	EarthParserRULE_userCommand = 5
	EarthParserRULE_userCommandHeader = 6
	EarthParserRULE_stmts = 7
	EarthParserRULE_stmt = 8
	EarthParserRULE_commandStmt = 9
	EarthParserRULE_version = 10
	EarthParserRULE_withStmt = 11
	EarthParserRULE_withBlock = 12
	EarthParserRULE_withExpr = 13
	EarthParserRULE_withCommand = 14
	EarthParserRULE_dockerCommand = 15
	EarthParserRULE_ifStmt = 16
	EarthParserRULE_ifClause = 17
	EarthParserRULE_ifBlock = 18
	EarthParserRULE_elseIfClause = 19
	EarthParserRULE_elseIfBlock = 20
	EarthParserRULE_elseClause = 21
	EarthParserRULE_elseBlock = 22
	EarthParserRULE_ifExpr = 23
	EarthParserRULE_elseIfExpr = 24
	EarthParserRULE_tryStmt = 25
	EarthParserRULE_tryClause = 26
	EarthParserRULE_tryBlock = 27
	EarthParserRULE_catchClause = 28
	EarthParserRULE_catchBlock = 29
	EarthParserRULE_finallyClause = 30
	EarthParserRULE_finallyBlock = 31
	EarthParserRULE_forStmt = 32
	EarthParserRULE_forClause = 33
	EarthParserRULE_forBlock = 34
	EarthParserRULE_forExpr = 35
	EarthParserRULE_waitStmt = 36
	EarthParserRULE_waitClause = 37
	EarthParserRULE_waitBlock = 38
	EarthParserRULE_waitExpr = 39
	EarthParserRULE_fromStmt = 40
	EarthParserRULE_fromDockerfileStmt = 41
	EarthParserRULE_locallyStmt = 42
	EarthParserRULE_copyStmt = 43
	EarthParserRULE_saveStmt = 44
	EarthParserRULE_saveImage = 45
	EarthParserRULE_saveArtifact = 46
	EarthParserRULE_runStmt = 47
	EarthParserRULE_buildStmt = 48
	EarthParserRULE_workdirStmt = 49
	EarthParserRULE_userStmt = 50
	EarthParserRULE_cmdStmt = 51
	EarthParserRULE_entrypointStmt = 52
	EarthParserRULE_exposeStmt = 53
	EarthParserRULE_volumeStmt = 54
	EarthParserRULE_envStmt = 55
	EarthParserRULE_argStmt = 56
	EarthParserRULE_optionalFlag = 57
	EarthParserRULE_envArgKey = 58
	EarthParserRULE_envArgValue = 59
	EarthParserRULE_labelStmt = 60
	EarthParserRULE_labelKey = 61
	EarthParserRULE_labelValue = 62
	EarthParserRULE_gitCloneStmt = 63
	EarthParserRULE_addStmt = 64
	EarthParserRULE_stopsignalStmt = 65
	EarthParserRULE_onbuildStmt = 66
	EarthParserRULE_healthcheckStmt = 67
	EarthParserRULE_shellStmt = 68
	EarthParserRULE_userCommandStmt = 69
	EarthParserRULE_doStmt = 70
	EarthParserRULE_importStmt = 71
	EarthParserRULE_cacheStmt = 72
	EarthParserRULE_hostStmt = 73
	EarthParserRULE_projectStmt = 74
	EarthParserRULE_pipelineStmt = 75
	EarthParserRULE_triggerStmt = 76
	EarthParserRULE_expr = 77
	EarthParserRULE_stmtWordsMaybeJSON = 78
	EarthParserRULE_stmtWords = 79
	EarthParserRULE_stmtWord = 80
)

// IEarthFileContext is an interface to support dynamic dispatch.
type IEarthFileContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EOF() antlr.TerminalNode
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode
	Version() IVersionContext
	Stmts() IStmtsContext
	Targets() ITargetsContext

	// IsEarthFileContext differentiates from other interfaces.
	IsEarthFileContext()
}

type EarthFileContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEarthFileContext() *EarthFileContext {
	var p = new(EarthFileContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_earthFile
	return p
}

func (*EarthFileContext) IsEarthFileContext() {}

func NewEarthFileContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EarthFileContext {
	var p = new(EarthFileContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_earthFile

	return p
}

func (s *EarthFileContext) GetParser() antlr.Parser { return s.parser }

func (s *EarthFileContext) EOF() antlr.TerminalNode {
	return s.GetToken(EarthParserEOF, 0)
}

func (s *EarthFileContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *EarthFileContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *EarthFileContext) Version() IVersionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVersionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVersionContext)
}

func (s *EarthFileContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *EarthFileContext) Targets() ITargetsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITargetsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITargetsContext)
}

func (s *EarthFileContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EarthFileContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *EarthFileContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEarthFile(s)
	}
}

func (s *EarthFileContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEarthFile(s)
	}
}




func (p *EarthParser) EarthFile() (localctx IEarthFileContext) {
	this := p
	_ = this

	localctx = NewEarthFileContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EarthParserRULE_earthFile)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(162)
				p.Match(EarthParserNL)
			}


		}
		p.SetState(167)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(169)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserVERSION {
		{
			p.SetState(168)
			p.Version()
		}

	}
	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 8519067631584) != 0) {
		{
			p.SetState(171)
			p.Stmts()
		}
		{
			p.SetState(172)
			p.Match(EarthParserNL)
		}

	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(176)
				p.Match(EarthParserNL)
			}


		}
		p.SetState(181)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserTarget || _la == EarthParserUserCommand {
		{
			p.SetState(182)
			p.Targets()
		}

	}
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserNL {
		{
			p.SetState(185)
			p.Match(EarthParserNL)
		}


		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(191)
		p.Match(EarthParserEOF)
	}



	return localctx
}


// ITargetsContext is an interface to support dynamic dispatch.
type ITargetsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTargetOrUserCommand() []ITargetOrUserCommandContext
	TargetOrUserCommand(i int) ITargetOrUserCommandContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsTargetsContext differentiates from other interfaces.
	IsTargetsContext()
}

type TargetsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetsContext() *TargetsContext {
	var p = new(TargetsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_targets
	return p
}

func (*TargetsContext) IsTargetsContext() {}

func NewTargetsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetsContext {
	var p = new(TargetsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_targets

	return p
}

func (s *TargetsContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetsContext) AllTargetOrUserCommand() []ITargetOrUserCommandContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITargetOrUserCommandContext); ok {
			len++
		}
	}

	tst := make([]ITargetOrUserCommandContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITargetOrUserCommandContext); ok {
			tst[i] = t.(ITargetOrUserCommandContext)
			i++
		}
	}

	return tst
}

func (s *TargetsContext) TargetOrUserCommand(i int) ITargetOrUserCommandContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITargetOrUserCommandContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITargetOrUserCommandContext)
}

func (s *TargetsContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *TargetsContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *TargetsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TargetsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTargets(s)
	}
}

func (s *TargetsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTargets(s)
	}
}




func (p *EarthParser) Targets() (localctx ITargetsContext) {
	this := p
	_ = this

	localctx = NewTargetsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, EarthParserRULE_targets)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(193)
		p.TargetOrUserCommand()
	}
	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(197)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for _la == EarthParserNL {
				{
					p.SetState(194)
					p.Match(EarthParserNL)
				}


				p.SetState(199)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(200)
				p.TargetOrUserCommand()
			}


		}
		p.SetState(205)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())
	}



	return localctx
}


// ITargetOrUserCommandContext is an interface to support dynamic dispatch.
type ITargetOrUserCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Target() ITargetContext
	UserCommand() IUserCommandContext

	// IsTargetOrUserCommandContext differentiates from other interfaces.
	IsTargetOrUserCommandContext()
}

type TargetOrUserCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetOrUserCommandContext() *TargetOrUserCommandContext {
	var p = new(TargetOrUserCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_targetOrUserCommand
	return p
}

func (*TargetOrUserCommandContext) IsTargetOrUserCommandContext() {}

func NewTargetOrUserCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetOrUserCommandContext {
	var p = new(TargetOrUserCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_targetOrUserCommand

	return p
}

func (s *TargetOrUserCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetOrUserCommandContext) Target() ITargetContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITargetContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITargetContext)
}

func (s *TargetOrUserCommandContext) UserCommand() IUserCommandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUserCommandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUserCommandContext)
}

func (s *TargetOrUserCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetOrUserCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TargetOrUserCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTargetOrUserCommand(s)
	}
}

func (s *TargetOrUserCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTargetOrUserCommand(s)
	}
}




func (p *EarthParser) TargetOrUserCommand() (localctx ITargetOrUserCommandContext) {
	this := p
	_ = this

	localctx = NewTargetOrUserCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, EarthParserRULE_targetOrUserCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(208)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserTarget:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(206)
			p.Target()
		}


	case EarthParserUserCommand:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(207)
			p.UserCommand()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}


	return localctx
}


// ITargetContext is an interface to support dynamic dispatch.
type ITargetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TargetHeader() ITargetHeaderContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode
	INDENT() antlr.TerminalNode
	DEDENT() antlr.TerminalNode
	Stmts() IStmtsContext

	// IsTargetContext differentiates from other interfaces.
	IsTargetContext()
}

type TargetContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetContext() *TargetContext {
	var p = new(TargetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_target
	return p
}

func (*TargetContext) IsTargetContext() {}

func NewTargetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetContext {
	var p = new(TargetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_target

	return p
}

func (s *TargetContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetContext) TargetHeader() ITargetHeaderContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITargetHeaderContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITargetHeaderContext)
}

func (s *TargetContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *TargetContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *TargetContext) INDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserINDENT, 0)
}

func (s *TargetContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserDEDENT, 0)
}

func (s *TargetContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *TargetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TargetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTarget(s)
	}
}

func (s *TargetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTarget(s)
	}
}




func (p *EarthParser) Target() (localctx ITargetContext) {
	this := p
	_ = this

	localctx = NewTargetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, EarthParserRULE_target)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(210)
		p.TargetHeader()
	}
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(211)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(214)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}
	p.SetState(232)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserINDENT {
		{
			p.SetState(216)
			p.Match(EarthParserINDENT)
		}
		p.SetState(220)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(217)
					p.Match(EarthParserNL)
				}


			}
			p.SetState(222)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
		}
		p.SetState(224)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 8519067631584) != 0) {
			{
				p.SetState(223)
				p.Stmts()
			}

		}
		p.SetState(227)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(226)
				p.Match(EarthParserNL)
			}


			p.SetState(229)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(231)
			p.Match(EarthParserDEDENT)
		}

	}



	return localctx
}


// ITargetHeaderContext is an interface to support dynamic dispatch.
type ITargetHeaderContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Target() antlr.TerminalNode

	// IsTargetHeaderContext differentiates from other interfaces.
	IsTargetHeaderContext()
}

type TargetHeaderContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetHeaderContext() *TargetHeaderContext {
	var p = new(TargetHeaderContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_targetHeader
	return p
}

func (*TargetHeaderContext) IsTargetHeaderContext() {}

func NewTargetHeaderContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetHeaderContext {
	var p = new(TargetHeaderContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_targetHeader

	return p
}

func (s *TargetHeaderContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetHeaderContext) Target() antlr.TerminalNode {
	return s.GetToken(EarthParserTarget, 0)
}

func (s *TargetHeaderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetHeaderContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TargetHeaderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTargetHeader(s)
	}
}

func (s *TargetHeaderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTargetHeader(s)
	}
}




func (p *EarthParser) TargetHeader() (localctx ITargetHeaderContext) {
	this := p
	_ = this

	localctx = NewTargetHeaderContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, EarthParserRULE_targetHeader)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(234)
		p.Match(EarthParserTarget)
	}



	return localctx
}


// IUserCommandContext is an interface to support dynamic dispatch.
type IUserCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UserCommandHeader() IUserCommandHeaderContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode
	INDENT() antlr.TerminalNode
	Stmts() IStmtsContext
	DEDENT() antlr.TerminalNode

	// IsUserCommandContext differentiates from other interfaces.
	IsUserCommandContext()
}

type UserCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUserCommandContext() *UserCommandContext {
	var p = new(UserCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_userCommand
	return p
}

func (*UserCommandContext) IsUserCommandContext() {}

func NewUserCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UserCommandContext {
	var p = new(UserCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_userCommand

	return p
}

func (s *UserCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *UserCommandContext) UserCommandHeader() IUserCommandHeaderContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUserCommandHeaderContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUserCommandHeaderContext)
}

func (s *UserCommandContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *UserCommandContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *UserCommandContext) INDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserINDENT, 0)
}

func (s *UserCommandContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *UserCommandContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserDEDENT, 0)
}

func (s *UserCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UserCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *UserCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterUserCommand(s)
	}
}

func (s *UserCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitUserCommand(s)
	}
}




func (p *EarthParser) UserCommand() (localctx IUserCommandContext) {
	this := p
	_ = this

	localctx = NewUserCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, EarthParserRULE_userCommand)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(236)
		p.UserCommandHeader()
	}
	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(237)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(240)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext())
	}
	p.SetState(257)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserINDENT {
		{
			p.SetState(242)
			p.Match(EarthParserINDENT)
		}
		p.SetState(246)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == EarthParserNL {
			{
				p.SetState(243)
				p.Match(EarthParserNL)
			}


			p.SetState(248)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(249)
			p.Stmts()
		}
		p.SetState(251)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(250)
				p.Match(EarthParserNL)
			}


			p.SetState(253)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(255)
			p.Match(EarthParserDEDENT)
		}

	}



	return localctx
}


// IUserCommandHeaderContext is an interface to support dynamic dispatch.
type IUserCommandHeaderContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UserCommand() antlr.TerminalNode

	// IsUserCommandHeaderContext differentiates from other interfaces.
	IsUserCommandHeaderContext()
}

type UserCommandHeaderContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUserCommandHeaderContext() *UserCommandHeaderContext {
	var p = new(UserCommandHeaderContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_userCommandHeader
	return p
}

func (*UserCommandHeaderContext) IsUserCommandHeaderContext() {}

func NewUserCommandHeaderContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UserCommandHeaderContext {
	var p = new(UserCommandHeaderContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_userCommandHeader

	return p
}

func (s *UserCommandHeaderContext) GetParser() antlr.Parser { return s.parser }

func (s *UserCommandHeaderContext) UserCommand() antlr.TerminalNode {
	return s.GetToken(EarthParserUserCommand, 0)
}

func (s *UserCommandHeaderContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UserCommandHeaderContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *UserCommandHeaderContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterUserCommandHeader(s)
	}
}

func (s *UserCommandHeaderContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitUserCommandHeader(s)
	}
}




func (p *EarthParser) UserCommandHeader() (localctx IUserCommandHeaderContext) {
	this := p
	_ = this

	localctx = NewUserCommandHeaderContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, EarthParserRULE_userCommandHeader)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(259)
		p.Match(EarthParserUserCommand)
	}



	return localctx
}


// IStmtsContext is an interface to support dynamic dispatch.
type IStmtsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStmt() []IStmtContext
	Stmt(i int) IStmtContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsStmtsContext differentiates from other interfaces.
	IsStmtsContext()
}

type StmtsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtsContext() *StmtsContext {
	var p = new(StmtsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stmts
	return p
}

func (*StmtsContext) IsStmtsContext() {}

func NewStmtsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtsContext {
	var p = new(StmtsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stmts

	return p
}

func (s *StmtsContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtsContext) AllStmt() []IStmtContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStmtContext); ok {
			len++
		}
	}

	tst := make([]IStmtContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStmtContext); ok {
			tst[i] = t.(IStmtContext)
			i++
		}
	}

	return tst
}

func (s *StmtsContext) Stmt(i int) IStmtContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtContext)
}

func (s *StmtsContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *StmtsContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *StmtsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StmtsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStmts(s)
	}
}

func (s *StmtsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStmts(s)
	}
}




func (p *EarthParser) Stmts() (localctx IStmtsContext) {
	this := p
	_ = this

	localctx = NewStmtsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, EarthParserRULE_stmts)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(261)
		p.Stmt()
	}
	p.SetState(270)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(263)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(262)
					p.Match(EarthParserNL)
				}


				p.SetState(265)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(267)
				p.Stmt()
			}


		}
		p.SetState(272)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())
	}



	return localctx
}


// IStmtContext is an interface to support dynamic dispatch.
type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CommandStmt() ICommandStmtContext
	WithStmt() IWithStmtContext
	IfStmt() IIfStmtContext
	ForStmt() IForStmtContext
	WaitStmt() IWaitStmtContext
	TryStmt() ITryStmtContext

	// IsStmtContext differentiates from other interfaces.
	IsStmtContext()
}

type StmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtContext() *StmtContext {
	var p = new(StmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stmt
	return p
}

func (*StmtContext) IsStmtContext() {}

func NewStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtContext {
	var p = new(StmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stmt

	return p
}

func (s *StmtContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtContext) CommandStmt() ICommandStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommandStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommandStmtContext)
}

func (s *StmtContext) WithStmt() IWithStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWithStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWithStmtContext)
}

func (s *StmtContext) IfStmt() IIfStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfStmtContext)
}

func (s *StmtContext) ForStmt() IForStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForStmtContext)
}

func (s *StmtContext) WaitStmt() IWaitStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitStmtContext)
}

func (s *StmtContext) TryStmt() ITryStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITryStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITryStmtContext)
}

func (s *StmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStmt(s)
	}
}

func (s *StmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStmt(s)
	}
}




func (p *EarthParser) Stmt() (localctx IStmtContext) {
	this := p
	_ = this

	localctx = NewStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, EarthParserRULE_stmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(279)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM, EarthParserFROM_DOCKERFILE, EarthParserLOCALLY, EarthParserCOPY, EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE, EarthParserRUN, EarthParserEXPOSE, EarthParserVOLUME, EarthParserENV, EarthParserARG, EarthParserLABEL, EarthParserBUILD, EarthParserWORKDIR, EarthParserUSER, EarthParserCMD, EarthParserENTRYPOINT, EarthParserGIT_CLONE, EarthParserADD, EarthParserSTOPSIGNAL, EarthParserONBUILD, EarthParserHEALTHCHECK, EarthParserSHELL, EarthParserDO, EarthParserCOMMAND, EarthParserIMPORT, EarthParserCACHE, EarthParserHOST, EarthParserPROJECT, EarthParserPIPELINE, EarthParserTRIGGER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(273)
			p.CommandStmt()
		}


	case EarthParserWITH:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(274)
			p.WithStmt()
		}


	case EarthParserIF:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(275)
			p.IfStmt()
		}


	case EarthParserFOR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(276)
			p.ForStmt()
		}


	case EarthParserWAIT:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(277)
			p.WaitStmt()
		}


	case EarthParserTRY:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(278)
			p.TryStmt()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}


	return localctx
}


// ICommandStmtContext is an interface to support dynamic dispatch.
type ICommandStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FromStmt() IFromStmtContext
	FromDockerfileStmt() IFromDockerfileStmtContext
	LocallyStmt() ILocallyStmtContext
	CopyStmt() ICopyStmtContext
	SaveStmt() ISaveStmtContext
	RunStmt() IRunStmtContext
	BuildStmt() IBuildStmtContext
	WorkdirStmt() IWorkdirStmtContext
	UserStmt() IUserStmtContext
	CmdStmt() ICmdStmtContext
	EntrypointStmt() IEntrypointStmtContext
	ExposeStmt() IExposeStmtContext
	VolumeStmt() IVolumeStmtContext
	EnvStmt() IEnvStmtContext
	ArgStmt() IArgStmtContext
	LabelStmt() ILabelStmtContext
	GitCloneStmt() IGitCloneStmtContext
	AddStmt() IAddStmtContext
	StopsignalStmt() IStopsignalStmtContext
	OnbuildStmt() IOnbuildStmtContext
	HealthcheckStmt() IHealthcheckStmtContext
	ShellStmt() IShellStmtContext
	UserCommandStmt() IUserCommandStmtContext
	DoStmt() IDoStmtContext
	ImportStmt() IImportStmtContext
	CacheStmt() ICacheStmtContext
	HostStmt() IHostStmtContext
	ProjectStmt() IProjectStmtContext
	PipelineStmt() IPipelineStmtContext
	TriggerStmt() ITriggerStmtContext

	// IsCommandStmtContext differentiates from other interfaces.
	IsCommandStmtContext()
}

type CommandStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandStmtContext() *CommandStmtContext {
	var p = new(CommandStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_commandStmt
	return p
}

func (*CommandStmtContext) IsCommandStmtContext() {}

func NewCommandStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandStmtContext {
	var p = new(CommandStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_commandStmt

	return p
}

func (s *CommandStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandStmtContext) FromStmt() IFromStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromStmtContext)
}

func (s *CommandStmtContext) FromDockerfileStmt() IFromDockerfileStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromDockerfileStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromDockerfileStmtContext)
}

func (s *CommandStmtContext) LocallyStmt() ILocallyStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILocallyStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILocallyStmtContext)
}

func (s *CommandStmtContext) CopyStmt() ICopyStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICopyStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICopyStmtContext)
}

func (s *CommandStmtContext) SaveStmt() ISaveStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISaveStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISaveStmtContext)
}

func (s *CommandStmtContext) RunStmt() IRunStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRunStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRunStmtContext)
}

func (s *CommandStmtContext) BuildStmt() IBuildStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBuildStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBuildStmtContext)
}

func (s *CommandStmtContext) WorkdirStmt() IWorkdirStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWorkdirStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWorkdirStmtContext)
}

func (s *CommandStmtContext) UserStmt() IUserStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUserStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUserStmtContext)
}

func (s *CommandStmtContext) CmdStmt() ICmdStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICmdStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICmdStmtContext)
}

func (s *CommandStmtContext) EntrypointStmt() IEntrypointStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEntrypointStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEntrypointStmtContext)
}

func (s *CommandStmtContext) ExposeStmt() IExposeStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExposeStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExposeStmtContext)
}

func (s *CommandStmtContext) VolumeStmt() IVolumeStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVolumeStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVolumeStmtContext)
}

func (s *CommandStmtContext) EnvStmt() IEnvStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnvStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnvStmtContext)
}

func (s *CommandStmtContext) ArgStmt() IArgStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgStmtContext)
}

func (s *CommandStmtContext) LabelStmt() ILabelStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILabelStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILabelStmtContext)
}

func (s *CommandStmtContext) GitCloneStmt() IGitCloneStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGitCloneStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGitCloneStmtContext)
}

func (s *CommandStmtContext) AddStmt() IAddStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAddStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAddStmtContext)
}

func (s *CommandStmtContext) StopsignalStmt() IStopsignalStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStopsignalStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStopsignalStmtContext)
}

func (s *CommandStmtContext) OnbuildStmt() IOnbuildStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOnbuildStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOnbuildStmtContext)
}

func (s *CommandStmtContext) HealthcheckStmt() IHealthcheckStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHealthcheckStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHealthcheckStmtContext)
}

func (s *CommandStmtContext) ShellStmt() IShellStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShellStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShellStmtContext)
}

func (s *CommandStmtContext) UserCommandStmt() IUserCommandStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUserCommandStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUserCommandStmtContext)
}

func (s *CommandStmtContext) DoStmt() IDoStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDoStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDoStmtContext)
}

func (s *CommandStmtContext) ImportStmt() IImportStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IImportStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IImportStmtContext)
}

func (s *CommandStmtContext) CacheStmt() ICacheStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICacheStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICacheStmtContext)
}

func (s *CommandStmtContext) HostStmt() IHostStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHostStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHostStmtContext)
}

func (s *CommandStmtContext) ProjectStmt() IProjectStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IProjectStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IProjectStmtContext)
}

func (s *CommandStmtContext) PipelineStmt() IPipelineStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPipelineStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPipelineStmtContext)
}

func (s *CommandStmtContext) TriggerStmt() ITriggerStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITriggerStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITriggerStmtContext)
}

func (s *CommandStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CommandStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCommandStmt(s)
	}
}

func (s *CommandStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCommandStmt(s)
	}
}




func (p *EarthParser) CommandStmt() (localctx ICommandStmtContext) {
	this := p
	_ = this

	localctx = NewCommandStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, EarthParserRULE_commandStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(311)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(281)
			p.FromStmt()
		}


	case EarthParserFROM_DOCKERFILE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(282)
			p.FromDockerfileStmt()
		}


	case EarthParserLOCALLY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(283)
			p.LocallyStmt()
		}


	case EarthParserCOPY:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(284)
			p.CopyStmt()
		}


	case EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(285)
			p.SaveStmt()
		}


	case EarthParserRUN:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(286)
			p.RunStmt()
		}


	case EarthParserBUILD:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(287)
			p.BuildStmt()
		}


	case EarthParserWORKDIR:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(288)
			p.WorkdirStmt()
		}


	case EarthParserUSER:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(289)
			p.UserStmt()
		}


	case EarthParserCMD:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(290)
			p.CmdStmt()
		}


	case EarthParserENTRYPOINT:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(291)
			p.EntrypointStmt()
		}


	case EarthParserEXPOSE:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(292)
			p.ExposeStmt()
		}


	case EarthParserVOLUME:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(293)
			p.VolumeStmt()
		}


	case EarthParserENV:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(294)
			p.EnvStmt()
		}


	case EarthParserARG:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(295)
			p.ArgStmt()
		}


	case EarthParserLABEL:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(296)
			p.LabelStmt()
		}


	case EarthParserGIT_CLONE:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(297)
			p.GitCloneStmt()
		}


	case EarthParserADD:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(298)
			p.AddStmt()
		}


	case EarthParserSTOPSIGNAL:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(299)
			p.StopsignalStmt()
		}


	case EarthParserONBUILD:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(300)
			p.OnbuildStmt()
		}


	case EarthParserHEALTHCHECK:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(301)
			p.HealthcheckStmt()
		}


	case EarthParserSHELL:
		p.EnterOuterAlt(localctx, 22)
		{
			p.SetState(302)
			p.ShellStmt()
		}


	case EarthParserCOMMAND:
		p.EnterOuterAlt(localctx, 23)
		{
			p.SetState(303)
			p.UserCommandStmt()
		}


	case EarthParserDO:
		p.EnterOuterAlt(localctx, 24)
		{
			p.SetState(304)
			p.DoStmt()
		}


	case EarthParserIMPORT:
		p.EnterOuterAlt(localctx, 25)
		{
			p.SetState(305)
			p.ImportStmt()
		}


	case EarthParserCACHE:
		p.EnterOuterAlt(localctx, 26)
		{
			p.SetState(306)
			p.CacheStmt()
		}


	case EarthParserHOST:
		p.EnterOuterAlt(localctx, 27)
		{
			p.SetState(307)
			p.HostStmt()
		}


	case EarthParserPROJECT:
		p.EnterOuterAlt(localctx, 28)
		{
			p.SetState(308)
			p.ProjectStmt()
		}


	case EarthParserPIPELINE:
		p.EnterOuterAlt(localctx, 29)
		{
			p.SetState(309)
			p.PipelineStmt()
		}


	case EarthParserTRIGGER:
		p.EnterOuterAlt(localctx, 30)
		{
			p.SetState(310)
			p.TriggerStmt()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}


	return localctx
}


// IVersionContext is an interface to support dynamic dispatch.
type IVersionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VERSION() antlr.TerminalNode
	StmtWords() IStmtWordsContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsVersionContext differentiates from other interfaces.
	IsVersionContext()
}

type VersionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVersionContext() *VersionContext {
	var p = new(VersionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_version
	return p
}

func (*VersionContext) IsVersionContext() {}

func NewVersionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VersionContext {
	var p = new(VersionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_version

	return p
}

func (s *VersionContext) GetParser() antlr.Parser { return s.parser }

func (s *VersionContext) VERSION() antlr.TerminalNode {
	return s.GetToken(EarthParserVERSION, 0)
}

func (s *VersionContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *VersionContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *VersionContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *VersionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VersionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *VersionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterVersion(s)
	}
}

func (s *VersionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitVersion(s)
	}
}




func (p *EarthParser) Version() (localctx IVersionContext) {
	this := p
	_ = this

	localctx = NewVersionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, EarthParserRULE_version)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(313)
		p.Match(EarthParserVERSION)
	}
	{
		p.SetState(314)
		p.StmtWords()
	}
	p.SetState(316)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(315)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(318)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())
	}



	return localctx
}


// IWithStmtContext is an interface to support dynamic dispatch.
type IWithStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WithExpr() IWithExprContext
	END() antlr.TerminalNode
	WithBlock() IWithBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsWithStmtContext differentiates from other interfaces.
	IsWithStmtContext()
}

type WithStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithStmtContext() *WithStmtContext {
	var p = new(WithStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_withStmt
	return p
}

func (*WithStmtContext) IsWithStmtContext() {}

func NewWithStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithStmtContext {
	var p = new(WithStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_withStmt

	return p
}

func (s *WithStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *WithStmtContext) WithExpr() IWithExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWithExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWithExprContext)
}

func (s *WithStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *WithStmtContext) WithBlock() IWithBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWithBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWithBlockContext)
}

func (s *WithStmtContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *WithStmtContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *WithStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WithStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWithStmt(s)
	}
}

func (s *WithStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWithStmt(s)
	}
}




func (p *EarthParser) WithStmt() (localctx IWithStmtContext) {
	this := p
	_ = this

	localctx = NewWithStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, EarthParserRULE_withStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(320)
		p.WithExpr()
	}
	p.SetState(327)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
		p.SetState(322)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(321)
				p.Match(EarthParserNL)
			}


			p.SetState(324)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(326)
			p.WithBlock()
		}


	}
	p.SetState(330)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(329)
			p.Match(EarthParserNL)
		}


		p.SetState(332)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(334)
		p.Match(EarthParserEND)
	}



	return localctx
}


// IWithBlockContext is an interface to support dynamic dispatch.
type IWithBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsWithBlockContext differentiates from other interfaces.
	IsWithBlockContext()
}

type WithBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithBlockContext() *WithBlockContext {
	var p = new(WithBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_withBlock
	return p
}

func (*WithBlockContext) IsWithBlockContext() {}

func NewWithBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithBlockContext {
	var p = new(WithBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_withBlock

	return p
}

func (s *WithBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *WithBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *WithBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WithBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWithBlock(s)
	}
}

func (s *WithBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWithBlock(s)
	}
}




func (p *EarthParser) WithBlock() (localctx IWithBlockContext) {
	this := p
	_ = this

	localctx = NewWithBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, EarthParserRULE_withBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(336)
		p.Stmts()
	}



	return localctx
}


// IWithExprContext is an interface to support dynamic dispatch.
type IWithExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WITH() antlr.TerminalNode
	WithCommand() IWithCommandContext

	// IsWithExprContext differentiates from other interfaces.
	IsWithExprContext()
}

type WithExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithExprContext() *WithExprContext {
	var p = new(WithExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_withExpr
	return p
}

func (*WithExprContext) IsWithExprContext() {}

func NewWithExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithExprContext {
	var p = new(WithExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_withExpr

	return p
}

func (s *WithExprContext) GetParser() antlr.Parser { return s.parser }

func (s *WithExprContext) WITH() antlr.TerminalNode {
	return s.GetToken(EarthParserWITH, 0)
}

func (s *WithExprContext) WithCommand() IWithCommandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWithCommandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWithCommandContext)
}

func (s *WithExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WithExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWithExpr(s)
	}
}

func (s *WithExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWithExpr(s)
	}
}




func (p *EarthParser) WithExpr() (localctx IWithExprContext) {
	this := p
	_ = this

	localctx = NewWithExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, EarthParserRULE_withExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(338)
		p.Match(EarthParserWITH)
	}
	{
		p.SetState(339)
		p.WithCommand()
	}



	return localctx
}


// IWithCommandContext is an interface to support dynamic dispatch.
type IWithCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DockerCommand() IDockerCommandContext

	// IsWithCommandContext differentiates from other interfaces.
	IsWithCommandContext()
}

type WithCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithCommandContext() *WithCommandContext {
	var p = new(WithCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_withCommand
	return p
}

func (*WithCommandContext) IsWithCommandContext() {}

func NewWithCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithCommandContext {
	var p = new(WithCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_withCommand

	return p
}

func (s *WithCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *WithCommandContext) DockerCommand() IDockerCommandContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDockerCommandContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDockerCommandContext)
}

func (s *WithCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WithCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWithCommand(s)
	}
}

func (s *WithCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWithCommand(s)
	}
}




func (p *EarthParser) WithCommand() (localctx IWithCommandContext) {
	this := p
	_ = this

	localctx = NewWithCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, EarthParserRULE_withCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(341)
		p.DockerCommand()
	}



	return localctx
}


// IDockerCommandContext is an interface to support dynamic dispatch.
type IDockerCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOCKER() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsDockerCommandContext differentiates from other interfaces.
	IsDockerCommandContext()
}

type DockerCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDockerCommandContext() *DockerCommandContext {
	var p = new(DockerCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_dockerCommand
	return p
}

func (*DockerCommandContext) IsDockerCommandContext() {}

func NewDockerCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DockerCommandContext {
	var p = new(DockerCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_dockerCommand

	return p
}

func (s *DockerCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *DockerCommandContext) DOCKER() antlr.TerminalNode {
	return s.GetToken(EarthParserDOCKER, 0)
}

func (s *DockerCommandContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *DockerCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DockerCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DockerCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterDockerCommand(s)
	}
}

func (s *DockerCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitDockerCommand(s)
	}
}




func (p *EarthParser) DockerCommand() (localctx IDockerCommandContext) {
	this := p
	_ = this

	localctx = NewDockerCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, EarthParserRULE_dockerCommand)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(343)
		p.Match(EarthParserDOCKER)
	}
	p.SetState(345)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(344)
			p.StmtWords()
		}

	}



	return localctx
}


// IIfStmtContext is an interface to support dynamic dispatch.
type IIfStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IfClause() IIfClauseContext
	END() antlr.TerminalNode
	AllElseIfClause() []IElseIfClauseContext
	ElseIfClause(i int) IElseIfClauseContext
	ElseClause() IElseClauseContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsIfStmtContext differentiates from other interfaces.
	IsIfStmtContext()
}

type IfStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIfStmtContext() *IfStmtContext {
	var p = new(IfStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_ifStmt
	return p
}

func (*IfStmtContext) IsIfStmtContext() {}

func NewIfStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfStmtContext {
	var p = new(IfStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_ifStmt

	return p
}

func (s *IfStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *IfStmtContext) IfClause() IIfClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfClauseContext)
}

func (s *IfStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *IfStmtContext) AllElseIfClause() []IElseIfClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IElseIfClauseContext); ok {
			len++
		}
	}

	tst := make([]IElseIfClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IElseIfClauseContext); ok {
			tst[i] = t.(IElseIfClauseContext)
			i++
		}
	}

	return tst
}

func (s *IfStmtContext) ElseIfClause(i int) IElseIfClauseContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseIfClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseIfClauseContext)
}

func (s *IfStmtContext) ElseClause() IElseClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseClauseContext)
}

func (s *IfStmtContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *IfStmtContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *IfStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IfStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterIfStmt(s)
	}
}

func (s *IfStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitIfStmt(s)
	}
}




func (p *EarthParser) IfStmt() (localctx IIfStmtContext) {
	this := p
	_ = this

	localctx = NewIfStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, EarthParserRULE_ifStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(347)
		p.IfClause()
	}
	p.SetState(356)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(349)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(348)
					p.Match(EarthParserNL)
				}


				p.SetState(351)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(353)
				p.ElseIfClause()
			}


		}
		p.SetState(358)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())
	}
	p.SetState(365)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) == 1 {
		p.SetState(360)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(359)
				p.Match(EarthParserNL)
			}


			p.SetState(362)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(364)
			p.ElseClause()
		}


	}
	p.SetState(368)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(367)
			p.Match(EarthParserNL)
		}


		p.SetState(370)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(372)
		p.Match(EarthParserEND)
	}



	return localctx
}


// IIfClauseContext is an interface to support dynamic dispatch.
type IIfClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IF() antlr.TerminalNode
	IfExpr() IIfExprContext
	IfBlock() IIfBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsIfClauseContext differentiates from other interfaces.
	IsIfClauseContext()
}

type IfClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIfClauseContext() *IfClauseContext {
	var p = new(IfClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_ifClause
	return p
}

func (*IfClauseContext) IsIfClauseContext() {}

func NewIfClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfClauseContext {
	var p = new(IfClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_ifClause

	return p
}

func (s *IfClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *IfClauseContext) IF() antlr.TerminalNode {
	return s.GetToken(EarthParserIF, 0)
}

func (s *IfClauseContext) IfExpr() IIfExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfExprContext)
}

func (s *IfClauseContext) IfBlock() IIfBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIfBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIfBlockContext)
}

func (s *IfClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *IfClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *IfClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IfClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterIfClause(s)
	}
}

func (s *IfClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitIfClause(s)
	}
}




func (p *EarthParser) IfClause() (localctx IIfClauseContext) {
	this := p
	_ = this

	localctx = NewIfClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, EarthParserRULE_ifClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(374)
		p.Match(EarthParserIF)
	}
	{
		p.SetState(375)
		p.IfExpr()
	}
	p.SetState(382)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
		p.SetState(377)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(376)
				p.Match(EarthParserNL)
			}


			p.SetState(379)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(381)
			p.IfBlock()
		}


	}



	return localctx
}


// IIfBlockContext is an interface to support dynamic dispatch.
type IIfBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsIfBlockContext differentiates from other interfaces.
	IsIfBlockContext()
}

type IfBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIfBlockContext() *IfBlockContext {
	var p = new(IfBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_ifBlock
	return p
}

func (*IfBlockContext) IsIfBlockContext() {}

func NewIfBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfBlockContext {
	var p = new(IfBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_ifBlock

	return p
}

func (s *IfBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *IfBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *IfBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IfBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterIfBlock(s)
	}
}

func (s *IfBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitIfBlock(s)
	}
}




func (p *EarthParser) IfBlock() (localctx IIfBlockContext) {
	this := p
	_ = this

	localctx = NewIfBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, EarthParserRULE_ifBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(384)
		p.Stmts()
	}



	return localctx
}


// IElseIfClauseContext is an interface to support dynamic dispatch.
type IElseIfClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ELSE_IF() antlr.TerminalNode
	ElseIfExpr() IElseIfExprContext
	ElseIfBlock() IElseIfBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsElseIfClauseContext differentiates from other interfaces.
	IsElseIfClauseContext()
}

type ElseIfClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElseIfClauseContext() *ElseIfClauseContext {
	var p = new(ElseIfClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_elseIfClause
	return p
}

func (*ElseIfClauseContext) IsElseIfClauseContext() {}

func NewElseIfClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseIfClauseContext {
	var p = new(ElseIfClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_elseIfClause

	return p
}

func (s *ElseIfClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ElseIfClauseContext) ELSE_IF() antlr.TerminalNode {
	return s.GetToken(EarthParserELSE_IF, 0)
}

func (s *ElseIfClauseContext) ElseIfExpr() IElseIfExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseIfExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseIfExprContext)
}

func (s *ElseIfClauseContext) ElseIfBlock() IElseIfBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseIfBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseIfBlockContext)
}

func (s *ElseIfClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *ElseIfClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *ElseIfClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElseIfClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ElseIfClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterElseIfClause(s)
	}
}

func (s *ElseIfClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitElseIfClause(s)
	}
}




func (p *EarthParser) ElseIfClause() (localctx IElseIfClauseContext) {
	this := p
	_ = this

	localctx = NewElseIfClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, EarthParserRULE_elseIfClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(386)
		p.Match(EarthParserELSE_IF)
	}
	{
		p.SetState(387)
		p.ElseIfExpr()
	}
	p.SetState(394)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext()) == 1 {
		p.SetState(389)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(388)
				p.Match(EarthParserNL)
			}


			p.SetState(391)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(393)
			p.ElseIfBlock()
		}


	}



	return localctx
}


// IElseIfBlockContext is an interface to support dynamic dispatch.
type IElseIfBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsElseIfBlockContext differentiates from other interfaces.
	IsElseIfBlockContext()
}

type ElseIfBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElseIfBlockContext() *ElseIfBlockContext {
	var p = new(ElseIfBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_elseIfBlock
	return p
}

func (*ElseIfBlockContext) IsElseIfBlockContext() {}

func NewElseIfBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseIfBlockContext {
	var p = new(ElseIfBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_elseIfBlock

	return p
}

func (s *ElseIfBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *ElseIfBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *ElseIfBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElseIfBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ElseIfBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterElseIfBlock(s)
	}
}

func (s *ElseIfBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitElseIfBlock(s)
	}
}




func (p *EarthParser) ElseIfBlock() (localctx IElseIfBlockContext) {
	this := p
	_ = this

	localctx = NewElseIfBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, EarthParserRULE_elseIfBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(396)
		p.Stmts()
	}



	return localctx
}


// IElseClauseContext is an interface to support dynamic dispatch.
type IElseClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ELSE() antlr.TerminalNode
	ElseBlock() IElseBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsElseClauseContext differentiates from other interfaces.
	IsElseClauseContext()
}

type ElseClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElseClauseContext() *ElseClauseContext {
	var p = new(ElseClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_elseClause
	return p
}

func (*ElseClauseContext) IsElseClauseContext() {}

func NewElseClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseClauseContext {
	var p = new(ElseClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_elseClause

	return p
}

func (s *ElseClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ElseClauseContext) ELSE() antlr.TerminalNode {
	return s.GetToken(EarthParserELSE, 0)
}

func (s *ElseClauseContext) ElseBlock() IElseBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseBlockContext)
}

func (s *ElseClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *ElseClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *ElseClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElseClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ElseClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterElseClause(s)
	}
}

func (s *ElseClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitElseClause(s)
	}
}




func (p *EarthParser) ElseClause() (localctx IElseClauseContext) {
	this := p
	_ = this

	localctx = NewElseClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, EarthParserRULE_elseClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(398)
		p.Match(EarthParserELSE)
	}
	p.SetState(405)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext()) == 1 {
		p.SetState(400)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(399)
				p.Match(EarthParserNL)
			}


			p.SetState(402)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(404)
			p.ElseBlock()
		}


	}



	return localctx
}


// IElseBlockContext is an interface to support dynamic dispatch.
type IElseBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsElseBlockContext differentiates from other interfaces.
	IsElseBlockContext()
}

type ElseBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElseBlockContext() *ElseBlockContext {
	var p = new(ElseBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_elseBlock
	return p
}

func (*ElseBlockContext) IsElseBlockContext() {}

func NewElseBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseBlockContext {
	var p = new(ElseBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_elseBlock

	return p
}

func (s *ElseBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *ElseBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *ElseBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElseBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ElseBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterElseBlock(s)
	}
}

func (s *ElseBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitElseBlock(s)
	}
}




func (p *EarthParser) ElseBlock() (localctx IElseBlockContext) {
	this := p
	_ = this

	localctx = NewElseBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, EarthParserRULE_elseBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(407)
		p.Stmts()
	}



	return localctx
}


// IIfExprContext is an interface to support dynamic dispatch.
type IIfExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext

	// IsIfExprContext differentiates from other interfaces.
	IsIfExprContext()
}

type IfExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIfExprContext() *IfExprContext {
	var p = new(IfExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_ifExpr
	return p
}

func (*IfExprContext) IsIfExprContext() {}

func NewIfExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfExprContext {
	var p = new(IfExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_ifExpr

	return p
}

func (s *IfExprContext) GetParser() antlr.Parser { return s.parser }

func (s *IfExprContext) Expr() IExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *IfExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IfExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IfExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterIfExpr(s)
	}
}

func (s *IfExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitIfExpr(s)
	}
}




func (p *EarthParser) IfExpr() (localctx IIfExprContext) {
	this := p
	_ = this

	localctx = NewIfExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, EarthParserRULE_ifExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(409)
		p.Expr()
	}



	return localctx
}


// IElseIfExprContext is an interface to support dynamic dispatch.
type IElseIfExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expr() IExprContext

	// IsElseIfExprContext differentiates from other interfaces.
	IsElseIfExprContext()
}

type ElseIfExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyElseIfExprContext() *ElseIfExprContext {
	var p = new(ElseIfExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_elseIfExpr
	return p
}

func (*ElseIfExprContext) IsElseIfExprContext() {}

func NewElseIfExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ElseIfExprContext {
	var p = new(ElseIfExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_elseIfExpr

	return p
}

func (s *ElseIfExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ElseIfExprContext) Expr() IExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ElseIfExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ElseIfExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ElseIfExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterElseIfExpr(s)
	}
}

func (s *ElseIfExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitElseIfExpr(s)
	}
}




func (p *EarthParser) ElseIfExpr() (localctx IElseIfExprContext) {
	this := p
	_ = this

	localctx = NewElseIfExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, EarthParserRULE_elseIfExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(411)
		p.Expr()
	}



	return localctx
}


// ITryStmtContext is an interface to support dynamic dispatch.
type ITryStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TryClause() ITryClauseContext
	END() antlr.TerminalNode
	CatchClause() ICatchClauseContext
	FinallyClause() IFinallyClauseContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsTryStmtContext differentiates from other interfaces.
	IsTryStmtContext()
}

type TryStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTryStmtContext() *TryStmtContext {
	var p = new(TryStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_tryStmt
	return p
}

func (*TryStmtContext) IsTryStmtContext() {}

func NewTryStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TryStmtContext {
	var p = new(TryStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_tryStmt

	return p
}

func (s *TryStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *TryStmtContext) TryClause() ITryClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITryClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITryClauseContext)
}

func (s *TryStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *TryStmtContext) CatchClause() ICatchClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICatchClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICatchClauseContext)
}

func (s *TryStmtContext) FinallyClause() IFinallyClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFinallyClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFinallyClauseContext)
}

func (s *TryStmtContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *TryStmtContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *TryStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TryStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TryStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTryStmt(s)
	}
}

func (s *TryStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTryStmt(s)
	}
}




func (p *EarthParser) TryStmt() (localctx ITryStmtContext) {
	this := p
	_ = this

	localctx = NewTryStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, EarthParserRULE_tryStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(413)
		p.TryClause()
	}
	p.SetState(420)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext()) == 1 {
		p.SetState(415)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(414)
				p.Match(EarthParserNL)
			}


			p.SetState(417)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(419)
			p.CatchClause()
		}


	}
	p.SetState(428)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) == 1 {
		p.SetState(423)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(422)
				p.Match(EarthParserNL)
			}


			p.SetState(425)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(427)
			p.FinallyClause()
		}


	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(430)
			p.Match(EarthParserNL)
		}


		p.SetState(433)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(435)
		p.Match(EarthParserEND)
	}



	return localctx
}


// ITryClauseContext is an interface to support dynamic dispatch.
type ITryClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRY() antlr.TerminalNode
	TryBlock() ITryBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsTryClauseContext differentiates from other interfaces.
	IsTryClauseContext()
}

type TryClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTryClauseContext() *TryClauseContext {
	var p = new(TryClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_tryClause
	return p
}

func (*TryClauseContext) IsTryClauseContext() {}

func NewTryClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TryClauseContext {
	var p = new(TryClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_tryClause

	return p
}

func (s *TryClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *TryClauseContext) TRY() antlr.TerminalNode {
	return s.GetToken(EarthParserTRY, 0)
}

func (s *TryClauseContext) TryBlock() ITryBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITryBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITryBlockContext)
}

func (s *TryClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *TryClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *TryClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TryClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TryClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTryClause(s)
	}
}

func (s *TryClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTryClause(s)
	}
}




func (p *EarthParser) TryClause() (localctx ITryClauseContext) {
	this := p
	_ = this

	localctx = NewTryClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, EarthParserRULE_tryClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(437)
		p.Match(EarthParserTRY)
	}
	p.SetState(444)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
		p.SetState(439)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(438)
				p.Match(EarthParserNL)
			}


			p.SetState(441)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(443)
			p.TryBlock()
		}


	}



	return localctx
}


// ITryBlockContext is an interface to support dynamic dispatch.
type ITryBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsTryBlockContext differentiates from other interfaces.
	IsTryBlockContext()
}

type TryBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTryBlockContext() *TryBlockContext {
	var p = new(TryBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_tryBlock
	return p
}

func (*TryBlockContext) IsTryBlockContext() {}

func NewTryBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TryBlockContext {
	var p = new(TryBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_tryBlock

	return p
}

func (s *TryBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *TryBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *TryBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TryBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TryBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTryBlock(s)
	}
}

func (s *TryBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTryBlock(s)
	}
}




func (p *EarthParser) TryBlock() (localctx ITryBlockContext) {
	this := p
	_ = this

	localctx = NewTryBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, EarthParserRULE_tryBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(446)
		p.Stmts()
	}



	return localctx
}


// ICatchClauseContext is an interface to support dynamic dispatch.
type ICatchClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CATCH() antlr.TerminalNode
	CatchBlock() ICatchBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsCatchClauseContext differentiates from other interfaces.
	IsCatchClauseContext()
}

type CatchClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCatchClauseContext() *CatchClauseContext {
	var p = new(CatchClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_catchClause
	return p
}

func (*CatchClauseContext) IsCatchClauseContext() {}

func NewCatchClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CatchClauseContext {
	var p = new(CatchClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_catchClause

	return p
}

func (s *CatchClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *CatchClauseContext) CATCH() antlr.TerminalNode {
	return s.GetToken(EarthParserCATCH, 0)
}

func (s *CatchClauseContext) CatchBlock() ICatchBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICatchBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICatchBlockContext)
}

func (s *CatchClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *CatchClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *CatchClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CatchClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CatchClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCatchClause(s)
	}
}

func (s *CatchClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCatchClause(s)
	}
}




func (p *EarthParser) CatchClause() (localctx ICatchClauseContext) {
	this := p
	_ = this

	localctx = NewCatchClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, EarthParserRULE_catchClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(448)
		p.Match(EarthParserCATCH)
	}
	p.SetState(455)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
		p.SetState(450)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(449)
				p.Match(EarthParserNL)
			}


			p.SetState(452)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(454)
			p.CatchBlock()
		}


	}



	return localctx
}


// ICatchBlockContext is an interface to support dynamic dispatch.
type ICatchBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsCatchBlockContext differentiates from other interfaces.
	IsCatchBlockContext()
}

type CatchBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCatchBlockContext() *CatchBlockContext {
	var p = new(CatchBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_catchBlock
	return p
}

func (*CatchBlockContext) IsCatchBlockContext() {}

func NewCatchBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CatchBlockContext {
	var p = new(CatchBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_catchBlock

	return p
}

func (s *CatchBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *CatchBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *CatchBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CatchBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CatchBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCatchBlock(s)
	}
}

func (s *CatchBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCatchBlock(s)
	}
}




func (p *EarthParser) CatchBlock() (localctx ICatchBlockContext) {
	this := p
	_ = this

	localctx = NewCatchBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, EarthParserRULE_catchBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(457)
		p.Stmts()
	}



	return localctx
}


// IFinallyClauseContext is an interface to support dynamic dispatch.
type IFinallyClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FINALLY() antlr.TerminalNode
	FinallyBlock() IFinallyBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsFinallyClauseContext differentiates from other interfaces.
	IsFinallyClauseContext()
}

type FinallyClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFinallyClauseContext() *FinallyClauseContext {
	var p = new(FinallyClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_finallyClause
	return p
}

func (*FinallyClauseContext) IsFinallyClauseContext() {}

func NewFinallyClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FinallyClauseContext {
	var p = new(FinallyClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_finallyClause

	return p
}

func (s *FinallyClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FinallyClauseContext) FINALLY() antlr.TerminalNode {
	return s.GetToken(EarthParserFINALLY, 0)
}

func (s *FinallyClauseContext) FinallyBlock() IFinallyBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFinallyBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFinallyBlockContext)
}

func (s *FinallyClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *FinallyClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *FinallyClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FinallyClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FinallyClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFinallyClause(s)
	}
}

func (s *FinallyClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFinallyClause(s)
	}
}




func (p *EarthParser) FinallyClause() (localctx IFinallyClauseContext) {
	this := p
	_ = this

	localctx = NewFinallyClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, EarthParserRULE_finallyClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(459)
		p.Match(EarthParserFINALLY)
	}
	p.SetState(466)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext()) == 1 {
		p.SetState(461)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(460)
				p.Match(EarthParserNL)
			}


			p.SetState(463)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(465)
			p.FinallyBlock()
		}


	}



	return localctx
}


// IFinallyBlockContext is an interface to support dynamic dispatch.
type IFinallyBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsFinallyBlockContext differentiates from other interfaces.
	IsFinallyBlockContext()
}

type FinallyBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFinallyBlockContext() *FinallyBlockContext {
	var p = new(FinallyBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_finallyBlock
	return p
}

func (*FinallyBlockContext) IsFinallyBlockContext() {}

func NewFinallyBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FinallyBlockContext {
	var p = new(FinallyBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_finallyBlock

	return p
}

func (s *FinallyBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *FinallyBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *FinallyBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FinallyBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FinallyBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFinallyBlock(s)
	}
}

func (s *FinallyBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFinallyBlock(s)
	}
}




func (p *EarthParser) FinallyBlock() (localctx IFinallyBlockContext) {
	this := p
	_ = this

	localctx = NewFinallyBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, EarthParserRULE_finallyBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(468)
		p.Stmts()
	}



	return localctx
}


// IForStmtContext is an interface to support dynamic dispatch.
type IForStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ForClause() IForClauseContext
	END() antlr.TerminalNode
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsForStmtContext differentiates from other interfaces.
	IsForStmtContext()
}

type ForStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForStmtContext() *ForStmtContext {
	var p = new(ForStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_forStmt
	return p
}

func (*ForStmtContext) IsForStmtContext() {}

func NewForStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForStmtContext {
	var p = new(ForStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_forStmt

	return p
}

func (s *ForStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ForStmtContext) ForClause() IForClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForClauseContext)
}

func (s *ForStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *ForStmtContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *ForStmtContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *ForStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ForStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterForStmt(s)
	}
}

func (s *ForStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitForStmt(s)
	}
}




func (p *EarthParser) ForStmt() (localctx IForStmtContext) {
	this := p
	_ = this

	localctx = NewForStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, EarthParserRULE_forStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(470)
		p.ForClause()
	}
	p.SetState(472)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(471)
			p.Match(EarthParserNL)
		}


		p.SetState(474)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(476)
		p.Match(EarthParserEND)
	}



	return localctx
}


// IForClauseContext is an interface to support dynamic dispatch.
type IForClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	ForExpr() IForExprContext
	ForBlock() IForBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsForClauseContext differentiates from other interfaces.
	IsForClauseContext()
}

type ForClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForClauseContext() *ForClauseContext {
	var p = new(ForClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_forClause
	return p
}

func (*ForClauseContext) IsForClauseContext() {}

func NewForClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForClauseContext {
	var p = new(ForClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_forClause

	return p
}

func (s *ForClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *ForClauseContext) FOR() antlr.TerminalNode {
	return s.GetToken(EarthParserFOR, 0)
}

func (s *ForClauseContext) ForExpr() IForExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForExprContext)
}

func (s *ForClauseContext) ForBlock() IForBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForBlockContext)
}

func (s *ForClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *ForClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *ForClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ForClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterForClause(s)
	}
}

func (s *ForClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitForClause(s)
	}
}




func (p *EarthParser) ForClause() (localctx IForClauseContext) {
	this := p
	_ = this

	localctx = NewForClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, EarthParserRULE_forClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(478)
		p.Match(EarthParserFOR)
	}
	{
		p.SetState(479)
		p.ForExpr()
	}
	p.SetState(486)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) == 1 {
		p.SetState(481)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(480)
				p.Match(EarthParserNL)
			}


			p.SetState(483)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(485)
			p.ForBlock()
		}


	}



	return localctx
}


// IForBlockContext is an interface to support dynamic dispatch.
type IForBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsForBlockContext differentiates from other interfaces.
	IsForBlockContext()
}

type ForBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForBlockContext() *ForBlockContext {
	var p = new(ForBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_forBlock
	return p
}

func (*ForBlockContext) IsForBlockContext() {}

func NewForBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForBlockContext {
	var p = new(ForBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_forBlock

	return p
}

func (s *ForBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *ForBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *ForBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ForBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterForBlock(s)
	}
}

func (s *ForBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitForBlock(s)
	}
}




func (p *EarthParser) ForBlock() (localctx IForBlockContext) {
	this := p
	_ = this

	localctx = NewForBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, EarthParserRULE_forBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(488)
		p.Stmts()
	}



	return localctx
}


// IForExprContext is an interface to support dynamic dispatch.
type IForExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StmtWords() IStmtWordsContext

	// IsForExprContext differentiates from other interfaces.
	IsForExprContext()
}

type ForExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyForExprContext() *ForExprContext {
	var p = new(ForExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_forExpr
	return p
}

func (*ForExprContext) IsForExprContext() {}

func NewForExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForExprContext {
	var p = new(ForExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_forExpr

	return p
}

func (s *ForExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ForExprContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *ForExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ForExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterForExpr(s)
	}
}

func (s *ForExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitForExpr(s)
	}
}




func (p *EarthParser) ForExpr() (localctx IForExprContext) {
	this := p
	_ = this

	localctx = NewForExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, EarthParserRULE_forExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(490)
		p.StmtWords()
	}



	return localctx
}


// IWaitStmtContext is an interface to support dynamic dispatch.
type IWaitStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WaitClause() IWaitClauseContext
	END() antlr.TerminalNode
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsWaitStmtContext differentiates from other interfaces.
	IsWaitStmtContext()
}

type WaitStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitStmtContext() *WaitStmtContext {
	var p = new(WaitStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_waitStmt
	return p
}

func (*WaitStmtContext) IsWaitStmtContext() {}

func NewWaitStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitStmtContext {
	var p = new(WaitStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_waitStmt

	return p
}

func (s *WaitStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitStmtContext) WaitClause() IWaitClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitClauseContext)
}

func (s *WaitStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *WaitStmtContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *WaitStmtContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *WaitStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WaitStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWaitStmt(s)
	}
}

func (s *WaitStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWaitStmt(s)
	}
}




func (p *EarthParser) WaitStmt() (localctx IWaitStmtContext) {
	this := p
	_ = this

	localctx = NewWaitStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, EarthParserRULE_waitStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(492)
		p.WaitClause()
	}
	p.SetState(494)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(493)
			p.Match(EarthParserNL)
		}


		p.SetState(496)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(498)
		p.Match(EarthParserEND)
	}



	return localctx
}


// IWaitClauseContext is an interface to support dynamic dispatch.
type IWaitClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WAIT() antlr.TerminalNode
	WaitExpr() IWaitExprContext
	WaitBlock() IWaitBlockContext
	AllNL() []antlr.TerminalNode
	NL(i int) antlr.TerminalNode

	// IsWaitClauseContext differentiates from other interfaces.
	IsWaitClauseContext()
}

type WaitClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitClauseContext() *WaitClauseContext {
	var p = new(WaitClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_waitClause
	return p
}

func (*WaitClauseContext) IsWaitClauseContext() {}

func NewWaitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitClauseContext {
	var p = new(WaitClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_waitClause

	return p
}

func (s *WaitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitClauseContext) WAIT() antlr.TerminalNode {
	return s.GetToken(EarthParserWAIT, 0)
}

func (s *WaitClauseContext) WaitExpr() IWaitExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitExprContext)
}

func (s *WaitClauseContext) WaitBlock() IWaitBlockContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWaitBlockContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWaitBlockContext)
}

func (s *WaitClauseContext) AllNL() []antlr.TerminalNode {
	return s.GetTokens(EarthParserNL)
}

func (s *WaitClauseContext) NL(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserNL, i)
}

func (s *WaitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WaitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWaitClause(s)
	}
}

func (s *WaitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWaitClause(s)
	}
}




func (p *EarthParser) WaitClause() (localctx IWaitClauseContext) {
	this := p
	_ = this

	localctx = NewWaitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, EarthParserRULE_waitClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(500)
		p.Match(EarthParserWAIT)
	}
	p.SetState(502)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(501)
			p.WaitExpr()
		}

	}
	p.SetState(510)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) == 1 {
		p.SetState(505)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(504)
				p.Match(EarthParserNL)
			}


			p.SetState(507)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(509)
			p.WaitBlock()
		}


	}



	return localctx
}


// IWaitBlockContext is an interface to support dynamic dispatch.
type IWaitBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Stmts() IStmtsContext

	// IsWaitBlockContext differentiates from other interfaces.
	IsWaitBlockContext()
}

type WaitBlockContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitBlockContext() *WaitBlockContext {
	var p = new(WaitBlockContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_waitBlock
	return p
}

func (*WaitBlockContext) IsWaitBlockContext() {}

func NewWaitBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitBlockContext {
	var p = new(WaitBlockContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_waitBlock

	return p
}

func (s *WaitBlockContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitBlockContext) Stmts() IStmtsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *WaitBlockContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WaitBlockContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWaitBlock(s)
	}
}

func (s *WaitBlockContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWaitBlock(s)
	}
}




func (p *EarthParser) WaitBlock() (localctx IWaitBlockContext) {
	this := p
	_ = this

	localctx = NewWaitBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, EarthParserRULE_waitBlock)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(512)
		p.Stmts()
	}



	return localctx
}


// IWaitExprContext is an interface to support dynamic dispatch.
type IWaitExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StmtWords() IStmtWordsContext

	// IsWaitExprContext differentiates from other interfaces.
	IsWaitExprContext()
}

type WaitExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWaitExprContext() *WaitExprContext {
	var p = new(WaitExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_waitExpr
	return p
}

func (*WaitExprContext) IsWaitExprContext() {}

func NewWaitExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WaitExprContext {
	var p = new(WaitExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_waitExpr

	return p
}

func (s *WaitExprContext) GetParser() antlr.Parser { return s.parser }

func (s *WaitExprContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *WaitExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WaitExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WaitExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWaitExpr(s)
	}
}

func (s *WaitExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWaitExpr(s)
	}
}




func (p *EarthParser) WaitExpr() (localctx IWaitExprContext) {
	this := p
	_ = this

	localctx = NewWaitExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, EarthParserRULE_waitExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(514)
		p.StmtWords()
	}



	return localctx
}


// IFromStmtContext is an interface to support dynamic dispatch.
type IFromStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FROM() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsFromStmtContext differentiates from other interfaces.
	IsFromStmtContext()
}

type FromStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromStmtContext() *FromStmtContext {
	var p = new(FromStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_fromStmt
	return p
}

func (*FromStmtContext) IsFromStmtContext() {}

func NewFromStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromStmtContext {
	var p = new(FromStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_fromStmt

	return p
}

func (s *FromStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *FromStmtContext) FROM() antlr.TerminalNode {
	return s.GetToken(EarthParserFROM, 0)
}

func (s *FromStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *FromStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FromStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFromStmt(s)
	}
}

func (s *FromStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFromStmt(s)
	}
}




func (p *EarthParser) FromStmt() (localctx IFromStmtContext) {
	this := p
	_ = this

	localctx = NewFromStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, EarthParserRULE_fromStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(516)
		p.Match(EarthParserFROM)
	}
	p.SetState(518)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(517)
			p.StmtWords()
		}

	}



	return localctx
}


// IFromDockerfileStmtContext is an interface to support dynamic dispatch.
type IFromDockerfileStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FROM_DOCKERFILE() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsFromDockerfileStmtContext differentiates from other interfaces.
	IsFromDockerfileStmtContext()
}

type FromDockerfileStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromDockerfileStmtContext() *FromDockerfileStmtContext {
	var p = new(FromDockerfileStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_fromDockerfileStmt
	return p
}

func (*FromDockerfileStmtContext) IsFromDockerfileStmtContext() {}

func NewFromDockerfileStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromDockerfileStmtContext {
	var p = new(FromDockerfileStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_fromDockerfileStmt

	return p
}

func (s *FromDockerfileStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *FromDockerfileStmtContext) FROM_DOCKERFILE() antlr.TerminalNode {
	return s.GetToken(EarthParserFROM_DOCKERFILE, 0)
}

func (s *FromDockerfileStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *FromDockerfileStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromDockerfileStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FromDockerfileStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFromDockerfileStmt(s)
	}
}

func (s *FromDockerfileStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFromDockerfileStmt(s)
	}
}




func (p *EarthParser) FromDockerfileStmt() (localctx IFromDockerfileStmtContext) {
	this := p
	_ = this

	localctx = NewFromDockerfileStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, EarthParserRULE_fromDockerfileStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(520)
		p.Match(EarthParserFROM_DOCKERFILE)
	}
	p.SetState(522)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(521)
			p.StmtWords()
		}

	}



	return localctx
}


// ILocallyStmtContext is an interface to support dynamic dispatch.
type ILocallyStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LOCALLY() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsLocallyStmtContext differentiates from other interfaces.
	IsLocallyStmtContext()
}

type LocallyStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLocallyStmtContext() *LocallyStmtContext {
	var p = new(LocallyStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_locallyStmt
	return p
}

func (*LocallyStmtContext) IsLocallyStmtContext() {}

func NewLocallyStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LocallyStmtContext {
	var p = new(LocallyStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_locallyStmt

	return p
}

func (s *LocallyStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *LocallyStmtContext) LOCALLY() antlr.TerminalNode {
	return s.GetToken(EarthParserLOCALLY, 0)
}

func (s *LocallyStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *LocallyStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LocallyStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LocallyStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterLocallyStmt(s)
	}
}

func (s *LocallyStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitLocallyStmt(s)
	}
}




func (p *EarthParser) LocallyStmt() (localctx ILocallyStmtContext) {
	this := p
	_ = this

	localctx = NewLocallyStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, EarthParserRULE_locallyStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(524)
		p.Match(EarthParserLOCALLY)
	}
	p.SetState(526)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(525)
			p.StmtWords()
		}

	}



	return localctx
}


// ICopyStmtContext is an interface to support dynamic dispatch.
type ICopyStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COPY() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsCopyStmtContext differentiates from other interfaces.
	IsCopyStmtContext()
}

type CopyStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopyStmtContext() *CopyStmtContext {
	var p = new(CopyStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copyStmt
	return p
}

func (*CopyStmtContext) IsCopyStmtContext() {}

func NewCopyStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopyStmtContext {
	var p = new(CopyStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copyStmt

	return p
}

func (s *CopyStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CopyStmtContext) COPY() antlr.TerminalNode {
	return s.GetToken(EarthParserCOPY, 0)
}

func (s *CopyStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *CopyStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopyStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CopyStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopyStmt(s)
	}
}

func (s *CopyStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopyStmt(s)
	}
}




func (p *EarthParser) CopyStmt() (localctx ICopyStmtContext) {
	this := p
	_ = this

	localctx = NewCopyStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, EarthParserRULE_copyStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(528)
		p.Match(EarthParserCOPY)
	}
	p.SetState(530)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(529)
			p.StmtWords()
		}

	}



	return localctx
}


// ISaveStmtContext is an interface to support dynamic dispatch.
type ISaveStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SaveArtifact() ISaveArtifactContext
	SaveImage() ISaveImageContext

	// IsSaveStmtContext differentiates from other interfaces.
	IsSaveStmtContext()
}

type SaveStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveStmtContext() *SaveStmtContext {
	var p = new(SaveStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveStmt
	return p
}

func (*SaveStmtContext) IsSaveStmtContext() {}

func NewSaveStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveStmtContext {
	var p = new(SaveStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveStmt

	return p
}

func (s *SaveStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveStmtContext) SaveArtifact() ISaveArtifactContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISaveArtifactContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISaveArtifactContext)
}

func (s *SaveStmtContext) SaveImage() ISaveImageContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISaveImageContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISaveImageContext)
}

func (s *SaveStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SaveStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveStmt(s)
	}
}

func (s *SaveStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveStmt(s)
	}
}




func (p *EarthParser) SaveStmt() (localctx ISaveStmtContext) {
	this := p
	_ = this

	localctx = NewSaveStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, EarthParserRULE_saveStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(534)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserSAVE_ARTIFACT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(532)
			p.SaveArtifact()
		}


	case EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(533)
			p.SaveImage()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}


	return localctx
}


// ISaveImageContext is an interface to support dynamic dispatch.
type ISaveImageContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SAVE_IMAGE() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsSaveImageContext differentiates from other interfaces.
	IsSaveImageContext()
}

type SaveImageContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveImageContext() *SaveImageContext {
	var p = new(SaveImageContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveImage
	return p
}

func (*SaveImageContext) IsSaveImageContext() {}

func NewSaveImageContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveImageContext {
	var p = new(SaveImageContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveImage

	return p
}

func (s *SaveImageContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveImageContext) SAVE_IMAGE() antlr.TerminalNode {
	return s.GetToken(EarthParserSAVE_IMAGE, 0)
}

func (s *SaveImageContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *SaveImageContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveImageContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SaveImageContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveImage(s)
	}
}

func (s *SaveImageContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveImage(s)
	}
}




func (p *EarthParser) SaveImage() (localctx ISaveImageContext) {
	this := p
	_ = this

	localctx = NewSaveImageContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, EarthParserRULE_saveImage)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(536)
		p.Match(EarthParserSAVE_IMAGE)
	}
	p.SetState(538)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(537)
			p.StmtWords()
		}

	}



	return localctx
}


// ISaveArtifactContext is an interface to support dynamic dispatch.
type ISaveArtifactContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SAVE_ARTIFACT() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsSaveArtifactContext differentiates from other interfaces.
	IsSaveArtifactContext()
}

type SaveArtifactContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveArtifactContext() *SaveArtifactContext {
	var p = new(SaveArtifactContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveArtifact
	return p
}

func (*SaveArtifactContext) IsSaveArtifactContext() {}

func NewSaveArtifactContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveArtifactContext {
	var p = new(SaveArtifactContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveArtifact

	return p
}

func (s *SaveArtifactContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveArtifactContext) SAVE_ARTIFACT() antlr.TerminalNode {
	return s.GetToken(EarthParserSAVE_ARTIFACT, 0)
}

func (s *SaveArtifactContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *SaveArtifactContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveArtifactContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SaveArtifactContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveArtifact(s)
	}
}

func (s *SaveArtifactContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveArtifact(s)
	}
}




func (p *EarthParser) SaveArtifact() (localctx ISaveArtifactContext) {
	this := p
	_ = this

	localctx = NewSaveArtifactContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, EarthParserRULE_saveArtifact)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(540)
		p.Match(EarthParserSAVE_ARTIFACT)
	}
	p.SetState(542)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(541)
			p.StmtWords()
		}

	}



	return localctx
}


// IRunStmtContext is an interface to support dynamic dispatch.
type IRunStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RUN() antlr.TerminalNode
	StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext

	// IsRunStmtContext differentiates from other interfaces.
	IsRunStmtContext()
}

type RunStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRunStmtContext() *RunStmtContext {
	var p = new(RunStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_runStmt
	return p
}

func (*RunStmtContext) IsRunStmtContext() {}

func NewRunStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RunStmtContext {
	var p = new(RunStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_runStmt

	return p
}

func (s *RunStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *RunStmtContext) RUN() antlr.TerminalNode {
	return s.GetToken(EarthParserRUN, 0)
}

func (s *RunStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsMaybeJSONContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsMaybeJSONContext)
}

func (s *RunStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RunStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *RunStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterRunStmt(s)
	}
}

func (s *RunStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitRunStmt(s)
	}
}




func (p *EarthParser) RunStmt() (localctx IRunStmtContext) {
	this := p
	_ = this

	localctx = NewRunStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, EarthParserRULE_runStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(544)
		p.Match(EarthParserRUN)
	}
	p.SetState(546)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(545)
			p.StmtWordsMaybeJSON()
		}

	}



	return localctx
}


// IBuildStmtContext is an interface to support dynamic dispatch.
type IBuildStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BUILD() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsBuildStmtContext differentiates from other interfaces.
	IsBuildStmtContext()
}

type BuildStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBuildStmtContext() *BuildStmtContext {
	var p = new(BuildStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_buildStmt
	return p
}

func (*BuildStmtContext) IsBuildStmtContext() {}

func NewBuildStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BuildStmtContext {
	var p = new(BuildStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_buildStmt

	return p
}

func (s *BuildStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *BuildStmtContext) BUILD() antlr.TerminalNode {
	return s.GetToken(EarthParserBUILD, 0)
}

func (s *BuildStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *BuildStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BuildStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BuildStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterBuildStmt(s)
	}
}

func (s *BuildStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitBuildStmt(s)
	}
}




func (p *EarthParser) BuildStmt() (localctx IBuildStmtContext) {
	this := p
	_ = this

	localctx = NewBuildStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, EarthParserRULE_buildStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(548)
		p.Match(EarthParserBUILD)
	}
	p.SetState(550)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(549)
			p.StmtWords()
		}

	}



	return localctx
}


// IWorkdirStmtContext is an interface to support dynamic dispatch.
type IWorkdirStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WORKDIR() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsWorkdirStmtContext differentiates from other interfaces.
	IsWorkdirStmtContext()
}

type WorkdirStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWorkdirStmtContext() *WorkdirStmtContext {
	var p = new(WorkdirStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_workdirStmt
	return p
}

func (*WorkdirStmtContext) IsWorkdirStmtContext() {}

func NewWorkdirStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WorkdirStmtContext {
	var p = new(WorkdirStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_workdirStmt

	return p
}

func (s *WorkdirStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *WorkdirStmtContext) WORKDIR() antlr.TerminalNode {
	return s.GetToken(EarthParserWORKDIR, 0)
}

func (s *WorkdirStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *WorkdirStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WorkdirStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WorkdirStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWorkdirStmt(s)
	}
}

func (s *WorkdirStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWorkdirStmt(s)
	}
}




func (p *EarthParser) WorkdirStmt() (localctx IWorkdirStmtContext) {
	this := p
	_ = this

	localctx = NewWorkdirStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, EarthParserRULE_workdirStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(552)
		p.Match(EarthParserWORKDIR)
	}
	p.SetState(554)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(553)
			p.StmtWords()
		}

	}



	return localctx
}


// IUserStmtContext is an interface to support dynamic dispatch.
type IUserStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	USER() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsUserStmtContext differentiates from other interfaces.
	IsUserStmtContext()
}

type UserStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUserStmtContext() *UserStmtContext {
	var p = new(UserStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_userStmt
	return p
}

func (*UserStmtContext) IsUserStmtContext() {}

func NewUserStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UserStmtContext {
	var p = new(UserStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_userStmt

	return p
}

func (s *UserStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *UserStmtContext) USER() antlr.TerminalNode {
	return s.GetToken(EarthParserUSER, 0)
}

func (s *UserStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *UserStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UserStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *UserStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterUserStmt(s)
	}
}

func (s *UserStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitUserStmt(s)
	}
}




func (p *EarthParser) UserStmt() (localctx IUserStmtContext) {
	this := p
	_ = this

	localctx = NewUserStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, EarthParserRULE_userStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(556)
		p.Match(EarthParserUSER)
	}
	p.SetState(558)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(557)
			p.StmtWords()
		}

	}



	return localctx
}


// ICmdStmtContext is an interface to support dynamic dispatch.
type ICmdStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CMD() antlr.TerminalNode
	StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext

	// IsCmdStmtContext differentiates from other interfaces.
	IsCmdStmtContext()
}

type CmdStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCmdStmtContext() *CmdStmtContext {
	var p = new(CmdStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_cmdStmt
	return p
}

func (*CmdStmtContext) IsCmdStmtContext() {}

func NewCmdStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CmdStmtContext {
	var p = new(CmdStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_cmdStmt

	return p
}

func (s *CmdStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CmdStmtContext) CMD() antlr.TerminalNode {
	return s.GetToken(EarthParserCMD, 0)
}

func (s *CmdStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsMaybeJSONContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsMaybeJSONContext)
}

func (s *CmdStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CmdStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CmdStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCmdStmt(s)
	}
}

func (s *CmdStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCmdStmt(s)
	}
}




func (p *EarthParser) CmdStmt() (localctx ICmdStmtContext) {
	this := p
	_ = this

	localctx = NewCmdStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, EarthParserRULE_cmdStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(560)
		p.Match(EarthParserCMD)
	}
	p.SetState(562)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(561)
			p.StmtWordsMaybeJSON()
		}

	}



	return localctx
}


// IEntrypointStmtContext is an interface to support dynamic dispatch.
type IEntrypointStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENTRYPOINT() antlr.TerminalNode
	StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext

	// IsEntrypointStmtContext differentiates from other interfaces.
	IsEntrypointStmtContext()
}

type EntrypointStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEntrypointStmtContext() *EntrypointStmtContext {
	var p = new(EntrypointStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_entrypointStmt
	return p
}

func (*EntrypointStmtContext) IsEntrypointStmtContext() {}

func NewEntrypointStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EntrypointStmtContext {
	var p = new(EntrypointStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_entrypointStmt

	return p
}

func (s *EntrypointStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *EntrypointStmtContext) ENTRYPOINT() antlr.TerminalNode {
	return s.GetToken(EarthParserENTRYPOINT, 0)
}

func (s *EntrypointStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsMaybeJSONContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsMaybeJSONContext)
}

func (s *EntrypointStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EntrypointStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *EntrypointStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEntrypointStmt(s)
	}
}

func (s *EntrypointStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEntrypointStmt(s)
	}
}




func (p *EarthParser) EntrypointStmt() (localctx IEntrypointStmtContext) {
	this := p
	_ = this

	localctx = NewEntrypointStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, EarthParserRULE_entrypointStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(564)
		p.Match(EarthParserENTRYPOINT)
	}
	p.SetState(566)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(565)
			p.StmtWordsMaybeJSON()
		}

	}



	return localctx
}


// IExposeStmtContext is an interface to support dynamic dispatch.
type IExposeStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXPOSE() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsExposeStmtContext differentiates from other interfaces.
	IsExposeStmtContext()
}

type ExposeStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExposeStmtContext() *ExposeStmtContext {
	var p = new(ExposeStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_exposeStmt
	return p
}

func (*ExposeStmtContext) IsExposeStmtContext() {}

func NewExposeStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExposeStmtContext {
	var p = new(ExposeStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_exposeStmt

	return p
}

func (s *ExposeStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ExposeStmtContext) EXPOSE() antlr.TerminalNode {
	return s.GetToken(EarthParserEXPOSE, 0)
}

func (s *ExposeStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *ExposeStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExposeStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ExposeStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterExposeStmt(s)
	}
}

func (s *ExposeStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitExposeStmt(s)
	}
}




func (p *EarthParser) ExposeStmt() (localctx IExposeStmtContext) {
	this := p
	_ = this

	localctx = NewExposeStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, EarthParserRULE_exposeStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(568)
		p.Match(EarthParserEXPOSE)
	}
	p.SetState(570)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(569)
			p.StmtWords()
		}

	}



	return localctx
}


// IVolumeStmtContext is an interface to support dynamic dispatch.
type IVolumeStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VOLUME() antlr.TerminalNode
	StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext

	// IsVolumeStmtContext differentiates from other interfaces.
	IsVolumeStmtContext()
}

type VolumeStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVolumeStmtContext() *VolumeStmtContext {
	var p = new(VolumeStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_volumeStmt
	return p
}

func (*VolumeStmtContext) IsVolumeStmtContext() {}

func NewVolumeStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VolumeStmtContext {
	var p = new(VolumeStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_volumeStmt

	return p
}

func (s *VolumeStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *VolumeStmtContext) VOLUME() antlr.TerminalNode {
	return s.GetToken(EarthParserVOLUME, 0)
}

func (s *VolumeStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsMaybeJSONContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsMaybeJSONContext)
}

func (s *VolumeStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VolumeStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *VolumeStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterVolumeStmt(s)
	}
}

func (s *VolumeStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitVolumeStmt(s)
	}
}




func (p *EarthParser) VolumeStmt() (localctx IVolumeStmtContext) {
	this := p
	_ = this

	localctx = NewVolumeStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, EarthParserRULE_volumeStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(572)
		p.Match(EarthParserVOLUME)
	}
	p.SetState(574)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(573)
			p.StmtWordsMaybeJSON()
		}

	}



	return localctx
}


// IEnvStmtContext is an interface to support dynamic dispatch.
type IEnvStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENV() antlr.TerminalNode
	EnvArgKey() IEnvArgKeyContext
	EQUALS() antlr.TerminalNode
	EnvArgValue() IEnvArgValueContext
	WS() antlr.TerminalNode

	// IsEnvStmtContext differentiates from other interfaces.
	IsEnvStmtContext()
}

type EnvStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEnvStmtContext() *EnvStmtContext {
	var p = new(EnvStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_envStmt
	return p
}

func (*EnvStmtContext) IsEnvStmtContext() {}

func NewEnvStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnvStmtContext {
	var p = new(EnvStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_envStmt

	return p
}

func (s *EnvStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *EnvStmtContext) ENV() antlr.TerminalNode {
	return s.GetToken(EarthParserENV, 0)
}

func (s *EnvStmtContext) EnvArgKey() IEnvArgKeyContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnvArgKeyContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnvArgKeyContext)
}

func (s *EnvStmtContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(EarthParserEQUALS, 0)
}

func (s *EnvStmtContext) EnvArgValue() IEnvArgValueContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnvArgValueContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnvArgValueContext)
}

func (s *EnvStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *EnvStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnvStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *EnvStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEnvStmt(s)
	}
}

func (s *EnvStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEnvStmt(s)
	}
}




func (p *EarthParser) EnvStmt() (localctx IEnvStmtContext) {
	this := p
	_ = this

	localctx = NewEnvStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, EarthParserRULE_envStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(576)
		p.Match(EarthParserENV)
	}
	{
		p.SetState(577)
		p.EnvArgKey()
	}
	p.SetState(579)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserEQUALS {
		{
			p.SetState(578)
			p.Match(EarthParserEQUALS)
		}

	}
	p.SetState(585)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(582)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(581)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(584)
			p.EnvArgValue()
		}

	}



	return localctx
}


// IArgStmtContext is an interface to support dynamic dispatch.
type IArgStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARG() antlr.TerminalNode
	OptionalFlag() IOptionalFlagContext
	EnvArgKey() IEnvArgKeyContext
	EQUALS() antlr.TerminalNode
	EnvArgValue() IEnvArgValueContext
	WS() antlr.TerminalNode

	// IsArgStmtContext differentiates from other interfaces.
	IsArgStmtContext()
}

type ArgStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgStmtContext() *ArgStmtContext {
	var p = new(ArgStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_argStmt
	return p
}

func (*ArgStmtContext) IsArgStmtContext() {}

func NewArgStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgStmtContext {
	var p = new(ArgStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_argStmt

	return p
}

func (s *ArgStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgStmtContext) ARG() antlr.TerminalNode {
	return s.GetToken(EarthParserARG, 0)
}

func (s *ArgStmtContext) OptionalFlag() IOptionalFlagContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOptionalFlagContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOptionalFlagContext)
}

func (s *ArgStmtContext) EnvArgKey() IEnvArgKeyContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnvArgKeyContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnvArgKeyContext)
}

func (s *ArgStmtContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(EarthParserEQUALS, 0)
}

func (s *ArgStmtContext) EnvArgValue() IEnvArgValueContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEnvArgValueContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEnvArgValueContext)
}

func (s *ArgStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *ArgStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ArgStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterArgStmt(s)
	}
}

func (s *ArgStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitArgStmt(s)
	}
}




func (p *EarthParser) ArgStmt() (localctx IArgStmtContext) {
	this := p
	_ = this

	localctx = NewArgStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, EarthParserRULE_argStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(587)
		p.Match(EarthParserARG)
	}
	{
		p.SetState(588)
		p.OptionalFlag()
	}
	{
		p.SetState(589)
		p.EnvArgKey()
	}
	p.SetState(597)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserEQUALS {
		{
			p.SetState(590)
			p.Match(EarthParserEQUALS)
		}
		p.SetState(595)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS || _la == EarthParserAtom {
			p.SetState(592)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			if _la == EarthParserWS {
				{
					p.SetState(591)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(594)
				p.EnvArgValue()
			}

		}

	}



	return localctx
}


// IOptionalFlagContext is an interface to support dynamic dispatch.
type IOptionalFlagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StmtWords() IStmtWordsContext

	// IsOptionalFlagContext differentiates from other interfaces.
	IsOptionalFlagContext()
}

type OptionalFlagContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOptionalFlagContext() *OptionalFlagContext {
	var p = new(OptionalFlagContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_optionalFlag
	return p
}

func (*OptionalFlagContext) IsOptionalFlagContext() {}

func NewOptionalFlagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptionalFlagContext {
	var p = new(OptionalFlagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_optionalFlag

	return p
}

func (s *OptionalFlagContext) GetParser() antlr.Parser { return s.parser }

func (s *OptionalFlagContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *OptionalFlagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OptionalFlagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OptionalFlagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterOptionalFlag(s)
	}
}

func (s *OptionalFlagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitOptionalFlag(s)
	}
}




func (p *EarthParser) OptionalFlag() (localctx IOptionalFlagContext) {
	this := p
	_ = this

	localctx = NewOptionalFlagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, EarthParserRULE_optionalFlag)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(600)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 77, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(599)
			p.StmtWords()
		}


	}



	return localctx
}


// IEnvArgKeyContext is an interface to support dynamic dispatch.
type IEnvArgKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Atom() antlr.TerminalNode

	// IsEnvArgKeyContext differentiates from other interfaces.
	IsEnvArgKeyContext()
}

type EnvArgKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEnvArgKeyContext() *EnvArgKeyContext {
	var p = new(EnvArgKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_envArgKey
	return p
}

func (*EnvArgKeyContext) IsEnvArgKeyContext() {}

func NewEnvArgKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnvArgKeyContext {
	var p = new(EnvArgKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_envArgKey

	return p
}

func (s *EnvArgKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *EnvArgKeyContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *EnvArgKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnvArgKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *EnvArgKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEnvArgKey(s)
	}
}

func (s *EnvArgKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEnvArgKey(s)
	}
}




func (p *EarthParser) EnvArgKey() (localctx IEnvArgKeyContext) {
	this := p
	_ = this

	localctx = NewEnvArgKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, EarthParserRULE_envArgKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(602)
		p.Match(EarthParserAtom)
	}



	return localctx
}


// IEnvArgValueContext is an interface to support dynamic dispatch.
type IEnvArgValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAtom() []antlr.TerminalNode
	Atom(i int) antlr.TerminalNode
	AllWS() []antlr.TerminalNode
	WS(i int) antlr.TerminalNode

	// IsEnvArgValueContext differentiates from other interfaces.
	IsEnvArgValueContext()
}

type EnvArgValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEnvArgValueContext() *EnvArgValueContext {
	var p = new(EnvArgValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_envArgValue
	return p
}

func (*EnvArgValueContext) IsEnvArgValueContext() {}

func NewEnvArgValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnvArgValueContext {
	var p = new(EnvArgValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_envArgValue

	return p
}

func (s *EnvArgValueContext) GetParser() antlr.Parser { return s.parser }

func (s *EnvArgValueContext) AllAtom() []antlr.TerminalNode {
	return s.GetTokens(EarthParserAtom)
}

func (s *EnvArgValueContext) Atom(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, i)
}

func (s *EnvArgValueContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *EnvArgValueContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *EnvArgValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EnvArgValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *EnvArgValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEnvArgValue(s)
	}
}

func (s *EnvArgValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEnvArgValue(s)
	}
}




func (p *EarthParser) EnvArgValue() (localctx IEnvArgValueContext) {
	this := p
	_ = this

	localctx = NewEnvArgValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, EarthParserRULE_envArgValue)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(604)
		p.Match(EarthParserAtom)
	}
	p.SetState(611)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(606)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(605)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(608)
			p.Match(EarthParserAtom)
		}


		p.SetState(613)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// ILabelStmtContext is an interface to support dynamic dispatch.
type ILabelStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LABEL() antlr.TerminalNode
	AllLabelKey() []ILabelKeyContext
	LabelKey(i int) ILabelKeyContext
	AllEQUALS() []antlr.TerminalNode
	EQUALS(i int) antlr.TerminalNode
	AllLabelValue() []ILabelValueContext
	LabelValue(i int) ILabelValueContext

	// IsLabelStmtContext differentiates from other interfaces.
	IsLabelStmtContext()
}

type LabelStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelStmtContext() *LabelStmtContext {
	var p = new(LabelStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_labelStmt
	return p
}

func (*LabelStmtContext) IsLabelStmtContext() {}

func NewLabelStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelStmtContext {
	var p = new(LabelStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_labelStmt

	return p
}

func (s *LabelStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelStmtContext) LABEL() antlr.TerminalNode {
	return s.GetToken(EarthParserLABEL, 0)
}

func (s *LabelStmtContext) AllLabelKey() []ILabelKeyContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILabelKeyContext); ok {
			len++
		}
	}

	tst := make([]ILabelKeyContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILabelKeyContext); ok {
			tst[i] = t.(ILabelKeyContext)
			i++
		}
	}

	return tst
}

func (s *LabelStmtContext) LabelKey(i int) ILabelKeyContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILabelKeyContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILabelKeyContext)
}

func (s *LabelStmtContext) AllEQUALS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserEQUALS)
}

func (s *LabelStmtContext) EQUALS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserEQUALS, i)
}

func (s *LabelStmtContext) AllLabelValue() []ILabelValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILabelValueContext); ok {
			len++
		}
	}

	tst := make([]ILabelValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILabelValueContext); ok {
			tst[i] = t.(ILabelValueContext)
			i++
		}
	}

	return tst
}

func (s *LabelStmtContext) LabelValue(i int) ILabelValueContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILabelValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILabelValueContext)
}

func (s *LabelStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LabelStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterLabelStmt(s)
	}
}

func (s *LabelStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitLabelStmt(s)
	}
}




func (p *EarthParser) LabelStmt() (localctx ILabelStmtContext) {
	this := p
	_ = this

	localctx = NewLabelStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, EarthParserRULE_labelStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(614)
		p.Match(EarthParserLABEL)
	}
	p.SetState(621)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserAtom {
		{
			p.SetState(615)
			p.LabelKey()
		}
		{
			p.SetState(616)
			p.Match(EarthParserEQUALS)
		}
		{
			p.SetState(617)
			p.LabelValue()
		}


		p.SetState(623)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// ILabelKeyContext is an interface to support dynamic dispatch.
type ILabelKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Atom() antlr.TerminalNode

	// IsLabelKeyContext differentiates from other interfaces.
	IsLabelKeyContext()
}

type LabelKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelKeyContext() *LabelKeyContext {
	var p = new(LabelKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_labelKey
	return p
}

func (*LabelKeyContext) IsLabelKeyContext() {}

func NewLabelKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelKeyContext {
	var p = new(LabelKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_labelKey

	return p
}

func (s *LabelKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelKeyContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *LabelKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LabelKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterLabelKey(s)
	}
}

func (s *LabelKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitLabelKey(s)
	}
}




func (p *EarthParser) LabelKey() (localctx ILabelKeyContext) {
	this := p
	_ = this

	localctx = NewLabelKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, EarthParserRULE_labelKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(624)
		p.Match(EarthParserAtom)
	}



	return localctx
}


// ILabelValueContext is an interface to support dynamic dispatch.
type ILabelValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Atom() antlr.TerminalNode

	// IsLabelValueContext differentiates from other interfaces.
	IsLabelValueContext()
}

type LabelValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLabelValueContext() *LabelValueContext {
	var p = new(LabelValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_labelValue
	return p
}

func (*LabelValueContext) IsLabelValueContext() {}

func NewLabelValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LabelValueContext {
	var p = new(LabelValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_labelValue

	return p
}

func (s *LabelValueContext) GetParser() antlr.Parser { return s.parser }

func (s *LabelValueContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *LabelValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LabelValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LabelValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterLabelValue(s)
	}
}

func (s *LabelValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitLabelValue(s)
	}
}




func (p *EarthParser) LabelValue() (localctx ILabelValueContext) {
	this := p
	_ = this

	localctx = NewLabelValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, EarthParserRULE_labelValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(626)
		p.Match(EarthParserAtom)
	}



	return localctx
}


// IGitCloneStmtContext is an interface to support dynamic dispatch.
type IGitCloneStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	GIT_CLONE() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsGitCloneStmtContext differentiates from other interfaces.
	IsGitCloneStmtContext()
}

type GitCloneStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGitCloneStmtContext() *GitCloneStmtContext {
	var p = new(GitCloneStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_gitCloneStmt
	return p
}

func (*GitCloneStmtContext) IsGitCloneStmtContext() {}

func NewGitCloneStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GitCloneStmtContext {
	var p = new(GitCloneStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_gitCloneStmt

	return p
}

func (s *GitCloneStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *GitCloneStmtContext) GIT_CLONE() antlr.TerminalNode {
	return s.GetToken(EarthParserGIT_CLONE, 0)
}

func (s *GitCloneStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *GitCloneStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GitCloneStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GitCloneStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterGitCloneStmt(s)
	}
}

func (s *GitCloneStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitGitCloneStmt(s)
	}
}




func (p *EarthParser) GitCloneStmt() (localctx IGitCloneStmtContext) {
	this := p
	_ = this

	localctx = NewGitCloneStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, EarthParserRULE_gitCloneStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(628)
		p.Match(EarthParserGIT_CLONE)
	}
	p.SetState(630)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(629)
			p.StmtWords()
		}

	}



	return localctx
}


// IAddStmtContext is an interface to support dynamic dispatch.
type IAddStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ADD() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsAddStmtContext differentiates from other interfaces.
	IsAddStmtContext()
}

type AddStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAddStmtContext() *AddStmtContext {
	var p = new(AddStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_addStmt
	return p
}

func (*AddStmtContext) IsAddStmtContext() {}

func NewAddStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AddStmtContext {
	var p = new(AddStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_addStmt

	return p
}

func (s *AddStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *AddStmtContext) ADD() antlr.TerminalNode {
	return s.GetToken(EarthParserADD, 0)
}

func (s *AddStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *AddStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AddStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *AddStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterAddStmt(s)
	}
}

func (s *AddStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitAddStmt(s)
	}
}




func (p *EarthParser) AddStmt() (localctx IAddStmtContext) {
	this := p
	_ = this

	localctx = NewAddStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, EarthParserRULE_addStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(632)
		p.Match(EarthParserADD)
	}
	p.SetState(634)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(633)
			p.StmtWords()
		}

	}



	return localctx
}


// IStopsignalStmtContext is an interface to support dynamic dispatch.
type IStopsignalStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STOPSIGNAL() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsStopsignalStmtContext differentiates from other interfaces.
	IsStopsignalStmtContext()
}

type StopsignalStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStopsignalStmtContext() *StopsignalStmtContext {
	var p = new(StopsignalStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stopsignalStmt
	return p
}

func (*StopsignalStmtContext) IsStopsignalStmtContext() {}

func NewStopsignalStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StopsignalStmtContext {
	var p = new(StopsignalStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stopsignalStmt

	return p
}

func (s *StopsignalStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *StopsignalStmtContext) STOPSIGNAL() antlr.TerminalNode {
	return s.GetToken(EarthParserSTOPSIGNAL, 0)
}

func (s *StopsignalStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *StopsignalStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StopsignalStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StopsignalStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStopsignalStmt(s)
	}
}

func (s *StopsignalStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStopsignalStmt(s)
	}
}




func (p *EarthParser) StopsignalStmt() (localctx IStopsignalStmtContext) {
	this := p
	_ = this

	localctx = NewStopsignalStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, EarthParserRULE_stopsignalStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(636)
		p.Match(EarthParserSTOPSIGNAL)
	}
	p.SetState(638)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(637)
			p.StmtWords()
		}

	}



	return localctx
}


// IOnbuildStmtContext is an interface to support dynamic dispatch.
type IOnbuildStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ONBUILD() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsOnbuildStmtContext differentiates from other interfaces.
	IsOnbuildStmtContext()
}

type OnbuildStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOnbuildStmtContext() *OnbuildStmtContext {
	var p = new(OnbuildStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_onbuildStmt
	return p
}

func (*OnbuildStmtContext) IsOnbuildStmtContext() {}

func NewOnbuildStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OnbuildStmtContext {
	var p = new(OnbuildStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_onbuildStmt

	return p
}

func (s *OnbuildStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *OnbuildStmtContext) ONBUILD() antlr.TerminalNode {
	return s.GetToken(EarthParserONBUILD, 0)
}

func (s *OnbuildStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *OnbuildStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OnbuildStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OnbuildStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterOnbuildStmt(s)
	}
}

func (s *OnbuildStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitOnbuildStmt(s)
	}
}




func (p *EarthParser) OnbuildStmt() (localctx IOnbuildStmtContext) {
	this := p
	_ = this

	localctx = NewOnbuildStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, EarthParserRULE_onbuildStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(640)
		p.Match(EarthParserONBUILD)
	}
	p.SetState(642)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(641)
			p.StmtWords()
		}

	}



	return localctx
}


// IHealthcheckStmtContext is an interface to support dynamic dispatch.
type IHealthcheckStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HEALTHCHECK() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsHealthcheckStmtContext differentiates from other interfaces.
	IsHealthcheckStmtContext()
}

type HealthcheckStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHealthcheckStmtContext() *HealthcheckStmtContext {
	var p = new(HealthcheckStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_healthcheckStmt
	return p
}

func (*HealthcheckStmtContext) IsHealthcheckStmtContext() {}

func NewHealthcheckStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HealthcheckStmtContext {
	var p = new(HealthcheckStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_healthcheckStmt

	return p
}

func (s *HealthcheckStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *HealthcheckStmtContext) HEALTHCHECK() antlr.TerminalNode {
	return s.GetToken(EarthParserHEALTHCHECK, 0)
}

func (s *HealthcheckStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *HealthcheckStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HealthcheckStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *HealthcheckStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterHealthcheckStmt(s)
	}
}

func (s *HealthcheckStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitHealthcheckStmt(s)
	}
}




func (p *EarthParser) HealthcheckStmt() (localctx IHealthcheckStmtContext) {
	this := p
	_ = this

	localctx = NewHealthcheckStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, EarthParserRULE_healthcheckStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(644)
		p.Match(EarthParserHEALTHCHECK)
	}
	p.SetState(646)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(645)
			p.StmtWords()
		}

	}



	return localctx
}


// IShellStmtContext is an interface to support dynamic dispatch.
type IShellStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SHELL() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsShellStmtContext differentiates from other interfaces.
	IsShellStmtContext()
}

type ShellStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShellStmtContext() *ShellStmtContext {
	var p = new(ShellStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_shellStmt
	return p
}

func (*ShellStmtContext) IsShellStmtContext() {}

func NewShellStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShellStmtContext {
	var p = new(ShellStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_shellStmt

	return p
}

func (s *ShellStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShellStmtContext) SHELL() antlr.TerminalNode {
	return s.GetToken(EarthParserSHELL, 0)
}

func (s *ShellStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *ShellStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShellStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShellStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterShellStmt(s)
	}
}

func (s *ShellStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitShellStmt(s)
	}
}




func (p *EarthParser) ShellStmt() (localctx IShellStmtContext) {
	this := p
	_ = this

	localctx = NewShellStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, EarthParserRULE_shellStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(648)
		p.Match(EarthParserSHELL)
	}
	p.SetState(650)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(649)
			p.StmtWords()
		}

	}



	return localctx
}


// IUserCommandStmtContext is an interface to support dynamic dispatch.
type IUserCommandStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMAND() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsUserCommandStmtContext differentiates from other interfaces.
	IsUserCommandStmtContext()
}

type UserCommandStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUserCommandStmtContext() *UserCommandStmtContext {
	var p = new(UserCommandStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_userCommandStmt
	return p
}

func (*UserCommandStmtContext) IsUserCommandStmtContext() {}

func NewUserCommandStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UserCommandStmtContext {
	var p = new(UserCommandStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_userCommandStmt

	return p
}

func (s *UserCommandStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *UserCommandStmtContext) COMMAND() antlr.TerminalNode {
	return s.GetToken(EarthParserCOMMAND, 0)
}

func (s *UserCommandStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *UserCommandStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UserCommandStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *UserCommandStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterUserCommandStmt(s)
	}
}

func (s *UserCommandStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitUserCommandStmt(s)
	}
}




func (p *EarthParser) UserCommandStmt() (localctx IUserCommandStmtContext) {
	this := p
	_ = this

	localctx = NewUserCommandStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, EarthParserRULE_userCommandStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(652)
		p.Match(EarthParserCOMMAND)
	}
	p.SetState(654)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(653)
			p.StmtWords()
		}

	}



	return localctx
}


// IDoStmtContext is an interface to support dynamic dispatch.
type IDoStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DO() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsDoStmtContext differentiates from other interfaces.
	IsDoStmtContext()
}

type DoStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDoStmtContext() *DoStmtContext {
	var p = new(DoStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_doStmt
	return p
}

func (*DoStmtContext) IsDoStmtContext() {}

func NewDoStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DoStmtContext {
	var p = new(DoStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_doStmt

	return p
}

func (s *DoStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *DoStmtContext) DO() antlr.TerminalNode {
	return s.GetToken(EarthParserDO, 0)
}

func (s *DoStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *DoStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DoStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DoStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterDoStmt(s)
	}
}

func (s *DoStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitDoStmt(s)
	}
}




func (p *EarthParser) DoStmt() (localctx IDoStmtContext) {
	this := p
	_ = this

	localctx = NewDoStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, EarthParserRULE_doStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(656)
		p.Match(EarthParserDO)
	}
	p.SetState(658)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(657)
			p.StmtWords()
		}

	}



	return localctx
}


// IImportStmtContext is an interface to support dynamic dispatch.
type IImportStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IMPORT() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsImportStmtContext differentiates from other interfaces.
	IsImportStmtContext()
}

type ImportStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportStmtContext() *ImportStmtContext {
	var p = new(ImportStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_importStmt
	return p
}

func (*ImportStmtContext) IsImportStmtContext() {}

func NewImportStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportStmtContext {
	var p = new(ImportStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_importStmt

	return p
}

func (s *ImportStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportStmtContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(EarthParserIMPORT, 0)
}

func (s *ImportStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *ImportStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ImportStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterImportStmt(s)
	}
}

func (s *ImportStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitImportStmt(s)
	}
}




func (p *EarthParser) ImportStmt() (localctx IImportStmtContext) {
	this := p
	_ = this

	localctx = NewImportStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, EarthParserRULE_importStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(660)
		p.Match(EarthParserIMPORT)
	}
	p.SetState(662)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(661)
			p.StmtWords()
		}

	}



	return localctx
}


// ICacheStmtContext is an interface to support dynamic dispatch.
type ICacheStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CACHE() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsCacheStmtContext differentiates from other interfaces.
	IsCacheStmtContext()
}

type CacheStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCacheStmtContext() *CacheStmtContext {
	var p = new(CacheStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_cacheStmt
	return p
}

func (*CacheStmtContext) IsCacheStmtContext() {}

func NewCacheStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CacheStmtContext {
	var p = new(CacheStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_cacheStmt

	return p
}

func (s *CacheStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CacheStmtContext) CACHE() antlr.TerminalNode {
	return s.GetToken(EarthParserCACHE, 0)
}

func (s *CacheStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *CacheStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CacheStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CacheStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCacheStmt(s)
	}
}

func (s *CacheStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCacheStmt(s)
	}
}




func (p *EarthParser) CacheStmt() (localctx ICacheStmtContext) {
	this := p
	_ = this

	localctx = NewCacheStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 144, EarthParserRULE_cacheStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(664)
		p.Match(EarthParserCACHE)
	}
	p.SetState(666)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(665)
			p.StmtWords()
		}

	}



	return localctx
}


// IHostStmtContext is an interface to support dynamic dispatch.
type IHostStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HOST() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsHostStmtContext differentiates from other interfaces.
	IsHostStmtContext()
}

type HostStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHostStmtContext() *HostStmtContext {
	var p = new(HostStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_hostStmt
	return p
}

func (*HostStmtContext) IsHostStmtContext() {}

func NewHostStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HostStmtContext {
	var p = new(HostStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_hostStmt

	return p
}

func (s *HostStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *HostStmtContext) HOST() antlr.TerminalNode {
	return s.GetToken(EarthParserHOST, 0)
}

func (s *HostStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *HostStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HostStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *HostStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterHostStmt(s)
	}
}

func (s *HostStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitHostStmt(s)
	}
}




func (p *EarthParser) HostStmt() (localctx IHostStmtContext) {
	this := p
	_ = this

	localctx = NewHostStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 146, EarthParserRULE_hostStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(668)
		p.Match(EarthParserHOST)
	}
	p.SetState(670)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(669)
			p.StmtWords()
		}

	}



	return localctx
}


// IProjectStmtContext is an interface to support dynamic dispatch.
type IProjectStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PROJECT() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsProjectStmtContext differentiates from other interfaces.
	IsProjectStmtContext()
}

type ProjectStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProjectStmtContext() *ProjectStmtContext {
	var p = new(ProjectStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_projectStmt
	return p
}

func (*ProjectStmtContext) IsProjectStmtContext() {}

func NewProjectStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProjectStmtContext {
	var p = new(ProjectStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_projectStmt

	return p
}

func (s *ProjectStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ProjectStmtContext) PROJECT() antlr.TerminalNode {
	return s.GetToken(EarthParserPROJECT, 0)
}

func (s *ProjectStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *ProjectStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProjectStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ProjectStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterProjectStmt(s)
	}
}

func (s *ProjectStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitProjectStmt(s)
	}
}




func (p *EarthParser) ProjectStmt() (localctx IProjectStmtContext) {
	this := p
	_ = this

	localctx = NewProjectStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 148, EarthParserRULE_projectStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(672)
		p.Match(EarthParserPROJECT)
	}
	p.SetState(674)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(673)
			p.StmtWords()
		}

	}



	return localctx
}


// IPipelineStmtContext is an interface to support dynamic dispatch.
type IPipelineStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PIPELINE() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsPipelineStmtContext differentiates from other interfaces.
	IsPipelineStmtContext()
}

type PipelineStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPipelineStmtContext() *PipelineStmtContext {
	var p = new(PipelineStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_pipelineStmt
	return p
}

func (*PipelineStmtContext) IsPipelineStmtContext() {}

func NewPipelineStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PipelineStmtContext {
	var p = new(PipelineStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_pipelineStmt

	return p
}

func (s *PipelineStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *PipelineStmtContext) PIPELINE() antlr.TerminalNode {
	return s.GetToken(EarthParserPIPELINE, 0)
}

func (s *PipelineStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *PipelineStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PipelineStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PipelineStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterPipelineStmt(s)
	}
}

func (s *PipelineStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitPipelineStmt(s)
	}
}




func (p *EarthParser) PipelineStmt() (localctx IPipelineStmtContext) {
	this := p
	_ = this

	localctx = NewPipelineStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 150, EarthParserRULE_pipelineStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(676)
		p.Match(EarthParserPIPELINE)
	}
	p.SetState(678)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(677)
			p.StmtWords()
		}

	}



	return localctx
}


// ITriggerStmtContext is an interface to support dynamic dispatch.
type ITriggerStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRIGGER() antlr.TerminalNode
	StmtWords() IStmtWordsContext

	// IsTriggerStmtContext differentiates from other interfaces.
	IsTriggerStmtContext()
}

type TriggerStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTriggerStmtContext() *TriggerStmtContext {
	var p = new(TriggerStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_triggerStmt
	return p
}

func (*TriggerStmtContext) IsTriggerStmtContext() {}

func NewTriggerStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TriggerStmtContext {
	var p = new(TriggerStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_triggerStmt

	return p
}

func (s *TriggerStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *TriggerStmtContext) TRIGGER() antlr.TerminalNode {
	return s.GetToken(EarthParserTRIGGER, 0)
}

func (s *TriggerStmtContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *TriggerStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TriggerStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TriggerStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTriggerStmt(s)
	}
}

func (s *TriggerStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTriggerStmt(s)
	}
}




func (p *EarthParser) TriggerStmt() (localctx ITriggerStmtContext) {
	this := p
	_ = this

	localctx = NewTriggerStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 152, EarthParserRULE_triggerStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(680)
		p.Match(EarthParserTRIGGER)
	}
	p.SetState(682)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(681)
			p.StmtWords()
		}

	}



	return localctx
}


// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsMaybeJSONContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsMaybeJSONContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitExpr(s)
	}
}




func (p *EarthParser) Expr() (localctx IExprContext) {
	this := p
	_ = this

	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 154, EarthParserRULE_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(684)
		p.StmtWordsMaybeJSON()
	}



	return localctx
}


// IStmtWordsMaybeJSONContext is an interface to support dynamic dispatch.
type IStmtWordsMaybeJSONContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	StmtWords() IStmtWordsContext

	// IsStmtWordsMaybeJSONContext differentiates from other interfaces.
	IsStmtWordsMaybeJSONContext()
}

type StmtWordsMaybeJSONContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtWordsMaybeJSONContext() *StmtWordsMaybeJSONContext {
	var p = new(StmtWordsMaybeJSONContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stmtWordsMaybeJSON
	return p
}

func (*StmtWordsMaybeJSONContext) IsStmtWordsMaybeJSONContext() {}

func NewStmtWordsMaybeJSONContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtWordsMaybeJSONContext {
	var p = new(StmtWordsMaybeJSONContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stmtWordsMaybeJSON

	return p
}

func (s *StmtWordsMaybeJSONContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtWordsMaybeJSONContext) StmtWords() IStmtWordsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *StmtWordsMaybeJSONContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtWordsMaybeJSONContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StmtWordsMaybeJSONContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStmtWordsMaybeJSON(s)
	}
}

func (s *StmtWordsMaybeJSONContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStmtWordsMaybeJSON(s)
	}
}




func (p *EarthParser) StmtWordsMaybeJSON() (localctx IStmtWordsMaybeJSONContext) {
	this := p
	_ = this

	localctx = NewStmtWordsMaybeJSONContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 156, EarthParserRULE_stmtWordsMaybeJSON)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(686)
		p.StmtWords()
	}



	return localctx
}


// IStmtWordsContext is an interface to support dynamic dispatch.
type IStmtWordsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllStmtWord() []IStmtWordContext
	StmtWord(i int) IStmtWordContext

	// IsStmtWordsContext differentiates from other interfaces.
	IsStmtWordsContext()
}

type StmtWordsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtWordsContext() *StmtWordsContext {
	var p = new(StmtWordsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stmtWords
	return p
}

func (*StmtWordsContext) IsStmtWordsContext() {}

func NewStmtWordsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtWordsContext {
	var p = new(StmtWordsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stmtWords

	return p
}

func (s *StmtWordsContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtWordsContext) AllStmtWord() []IStmtWordContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStmtWordContext); ok {
			len++
		}
	}

	tst := make([]IStmtWordContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStmtWordContext); ok {
			tst[i] = t.(IStmtWordContext)
			i++
		}
	}

	return tst
}

func (s *StmtWordsContext) StmtWord(i int) IStmtWordContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStmtWordContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStmtWordContext)
}

func (s *StmtWordsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtWordsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StmtWordsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStmtWords(s)
	}
}

func (s *StmtWordsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStmtWords(s)
	}
}




func (p *EarthParser) StmtWords() (localctx IStmtWordsContext) {
	this := p
	_ = this

	localctx = NewStmtWordsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 158, EarthParserRULE_stmtWords)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(689)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(688)
					p.StmtWord()
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(691)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 95, p.GetParserRuleContext())
	}



	return localctx
}


// IStmtWordContext is an interface to support dynamic dispatch.
type IStmtWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Atom() antlr.TerminalNode

	// IsStmtWordContext differentiates from other interfaces.
	IsStmtWordContext()
}

type StmtWordContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtWordContext() *StmtWordContext {
	var p = new(StmtWordContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stmtWord
	return p
}

func (*StmtWordContext) IsStmtWordContext() {}

func NewStmtWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtWordContext {
	var p = new(StmtWordContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stmtWord

	return p
}

func (s *StmtWordContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtWordContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *StmtWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StmtWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStmtWord(s)
	}
}

func (s *StmtWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStmtWord(s)
	}
}




func (p *EarthParser) StmtWord() (localctx IStmtWordContext) {
	this := p
	_ = this

	localctx = NewStmtWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 160, EarthParserRULE_stmtWord)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(693)
		p.Match(EarthParserAtom)
	}



	return localctx
}


