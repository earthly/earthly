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
    "'ENV'", "'ARG'", "'SET'", "'LABEL'", "'BUILD'", "'WORKDIR'", "'USER'", 
    "'CMD'", "'ENTRYPOINT'", "'GIT CLONE'", "'ADD'", "'STOPSIGNAL'", "'ONBUILD'", 
    "'HEALTHCHECK'", "'SHELL'", "'DO'", "'COMMAND'", "'IMPORT'", "'VERSION'", 
    "'CACHE'", "'HOST'", "'PROJECT'", "'PIPELINE'", "'TRIGGER'", "'WITH'", 
    "", "", "", "", "", "", "", "", "'ELSE'", "'ELSE IF'", "'CATCH'", "'FINALLY'", 
    "'END'",
  }
  staticData.symbolicNames = []string{
    "", "INDENT", "DEDENT", "Target", "UserCommand", "FROM", "FROM_DOCKERFILE", 
    "LOCALLY", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE", "RUN", "EXPOSE", "VOLUME", 
    "ENV", "ARG", "SET", "LABEL", "BUILD", "WORKDIR", "USER", "CMD", "ENTRYPOINT", 
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
    "envStmt", "argStmt", "setStmt", "optionalFlag", "envArgKey", "envArgValue", 
    "labelStmt", "labelKey", "labelValue", "gitCloneStmt", "addStmt", "stopsignalStmt", 
    "onbuildStmt", "healthcheckStmt", "shellStmt", "userCommandStmt", "doStmt", 
    "importStmt", "cacheStmt", "hostStmt", "projectStmt", "pipelineStmt", 
    "triggerStmt", "expr", "stmtWordsMaybeJSON", "stmtWords", "stmtWord",
  }
  staticData.predictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 53, 707, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
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
	2, 79, 7, 79, 2, 80, 7, 80, 2, 81, 7, 81, 1, 0, 5, 0, 166, 8, 0, 10, 0, 
	12, 0, 169, 9, 0, 1, 0, 3, 0, 172, 8, 0, 1, 0, 1, 0, 1, 0, 3, 0, 177, 8, 
	0, 1, 0, 5, 0, 180, 8, 0, 10, 0, 12, 0, 183, 9, 0, 1, 0, 3, 0, 186, 8, 
	0, 1, 0, 5, 0, 189, 8, 0, 10, 0, 12, 0, 192, 9, 0, 1, 0, 1, 0, 1, 1, 1, 
	1, 5, 1, 198, 8, 1, 10, 1, 12, 1, 201, 9, 1, 1, 1, 5, 1, 204, 8, 1, 10, 
	1, 12, 1, 207, 9, 1, 1, 2, 1, 2, 3, 2, 211, 8, 2, 1, 3, 1, 3, 4, 3, 215, 
	8, 3, 11, 3, 12, 3, 216, 1, 3, 1, 3, 5, 3, 221, 8, 3, 10, 3, 12, 3, 224, 
	9, 3, 1, 3, 3, 3, 227, 8, 3, 1, 3, 4, 3, 230, 8, 3, 11, 3, 12, 3, 231, 
	1, 3, 3, 3, 235, 8, 3, 1, 4, 1, 4, 1, 5, 1, 5, 4, 5, 241, 8, 5, 11, 5, 
	12, 5, 242, 1, 5, 1, 5, 5, 5, 247, 8, 5, 10, 5, 12, 5, 250, 9, 5, 1, 5, 
	1, 5, 4, 5, 254, 8, 5, 11, 5, 12, 5, 255, 1, 5, 1, 5, 3, 5, 260, 8, 5, 
	1, 6, 1, 6, 1, 7, 1, 7, 4, 7, 266, 8, 7, 11, 7, 12, 7, 267, 1, 7, 5, 7, 
	271, 8, 7, 10, 7, 12, 7, 274, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 
	3, 8, 282, 8, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 
	1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 
	1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 315, 
	8, 9, 1, 10, 1, 10, 1, 10, 4, 10, 320, 8, 10, 11, 10, 12, 10, 321, 1, 11, 
	1, 11, 4, 11, 326, 8, 11, 11, 11, 12, 11, 327, 1, 11, 3, 11, 331, 8, 11, 
	1, 11, 4, 11, 334, 8, 11, 11, 11, 12, 11, 335, 1, 11, 1, 11, 1, 12, 1, 
	12, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 15, 1, 15, 3, 15, 349, 8, 15, 
	1, 16, 1, 16, 4, 16, 353, 8, 16, 11, 16, 12, 16, 354, 1, 16, 5, 16, 358, 
	8, 16, 10, 16, 12, 16, 361, 9, 16, 1, 16, 4, 16, 364, 8, 16, 11, 16, 12, 
	16, 365, 1, 16, 3, 16, 369, 8, 16, 1, 16, 4, 16, 372, 8, 16, 11, 16, 12, 
	16, 373, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 4, 17, 381, 8, 17, 11, 17, 
	12, 17, 382, 1, 17, 3, 17, 386, 8, 17, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 
	4, 19, 393, 8, 19, 11, 19, 12, 19, 394, 1, 19, 3, 19, 398, 8, 19, 1, 20, 
	1, 20, 1, 21, 1, 21, 4, 21, 404, 8, 21, 11, 21, 12, 21, 405, 1, 21, 3, 
	21, 409, 8, 21, 1, 22, 1, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 
	4, 25, 419, 8, 25, 11, 25, 12, 25, 420, 1, 25, 3, 25, 424, 8, 25, 1, 25, 
	4, 25, 427, 8, 25, 11, 25, 12, 25, 428, 1, 25, 3, 25, 432, 8, 25, 1, 25, 
	4, 25, 435, 8, 25, 11, 25, 12, 25, 436, 1, 25, 1, 25, 1, 26, 1, 26, 4, 
	26, 443, 8, 26, 11, 26, 12, 26, 444, 1, 26, 3, 26, 448, 8, 26, 1, 27, 1, 
	27, 1, 28, 1, 28, 4, 28, 454, 8, 28, 11, 28, 12, 28, 455, 1, 28, 3, 28, 
	459, 8, 28, 1, 29, 1, 29, 1, 30, 1, 30, 4, 30, 465, 8, 30, 11, 30, 12, 
	30, 466, 1, 30, 3, 30, 470, 8, 30, 1, 31, 1, 31, 1, 32, 1, 32, 4, 32, 476, 
	8, 32, 11, 32, 12, 32, 477, 1, 32, 1, 32, 1, 33, 1, 33, 1, 33, 4, 33, 485, 
	8, 33, 11, 33, 12, 33, 486, 1, 33, 3, 33, 490, 8, 33, 1, 34, 1, 34, 1, 
	35, 1, 35, 1, 36, 1, 36, 4, 36, 498, 8, 36, 11, 36, 12, 36, 499, 1, 36, 
	1, 36, 1, 37, 1, 37, 3, 37, 506, 8, 37, 1, 37, 4, 37, 509, 8, 37, 11, 37, 
	12, 37, 510, 1, 37, 3, 37, 514, 8, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 
	1, 40, 3, 40, 522, 8, 40, 1, 41, 1, 41, 3, 41, 526, 8, 41, 1, 42, 1, 42, 
	3, 42, 530, 8, 42, 1, 43, 1, 43, 3, 43, 534, 8, 43, 1, 44, 1, 44, 3, 44, 
	538, 8, 44, 1, 45, 1, 45, 3, 45, 542, 8, 45, 1, 46, 1, 46, 3, 46, 546, 
	8, 46, 1, 47, 1, 47, 3, 47, 550, 8, 47, 1, 48, 1, 48, 3, 48, 554, 8, 48, 
	1, 49, 1, 49, 3, 49, 558, 8, 49, 1, 50, 1, 50, 3, 50, 562, 8, 50, 1, 51, 
	1, 51, 3, 51, 566, 8, 51, 1, 52, 1, 52, 3, 52, 570, 8, 52, 1, 53, 1, 53, 
	3, 53, 574, 8, 53, 1, 54, 1, 54, 3, 54, 578, 8, 54, 1, 55, 1, 55, 1, 55, 
	3, 55, 583, 8, 55, 1, 55, 3, 55, 586, 8, 55, 1, 55, 3, 55, 589, 8, 55, 
	1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 596, 8, 56, 1, 56, 3, 56, 599, 
	8, 56, 3, 56, 601, 8, 56, 1, 57, 1, 57, 1, 57, 1, 57, 3, 57, 607, 8, 57, 
	1, 57, 1, 57, 1, 58, 3, 58, 612, 8, 58, 1, 59, 1, 59, 1, 60, 1, 60, 3, 
	60, 618, 8, 60, 1, 60, 5, 60, 621, 8, 60, 10, 60, 12, 60, 624, 9, 60, 1, 
	61, 1, 61, 1, 61, 1, 61, 1, 61, 5, 61, 631, 8, 61, 10, 61, 12, 61, 634, 
	9, 61, 1, 62, 1, 62, 1, 63, 1, 63, 1, 64, 1, 64, 3, 64, 642, 8, 64, 1, 
	65, 1, 65, 3, 65, 646, 8, 65, 1, 66, 1, 66, 3, 66, 650, 8, 66, 1, 67, 1, 
	67, 3, 67, 654, 8, 67, 1, 68, 1, 68, 3, 68, 658, 8, 68, 1, 69, 1, 69, 3, 
	69, 662, 8, 69, 1, 70, 1, 70, 3, 70, 666, 8, 70, 1, 71, 1, 71, 3, 71, 670, 
	8, 71, 1, 72, 1, 72, 3, 72, 674, 8, 72, 1, 73, 1, 73, 3, 73, 678, 8, 73, 
	1, 74, 1, 74, 3, 74, 682, 8, 74, 1, 75, 1, 75, 3, 75, 686, 8, 75, 1, 76, 
	1, 76, 3, 76, 690, 8, 76, 1, 77, 1, 77, 3, 77, 694, 8, 77, 1, 78, 1, 78, 
	1, 79, 1, 79, 1, 80, 4, 80, 701, 8, 80, 11, 80, 12, 80, 702, 1, 81, 1, 
	81, 1, 81, 0, 0, 82, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 
	28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 
	64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 
	100, 102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 
	130, 132, 134, 136, 138, 140, 142, 144, 146, 148, 150, 152, 154, 156, 158, 
	160, 162, 0, 0, 754, 0, 167, 1, 0, 0, 0, 2, 195, 1, 0, 0, 0, 4, 210, 1, 
	0, 0, 0, 6, 212, 1, 0, 0, 0, 8, 236, 1, 0, 0, 0, 10, 238, 1, 0, 0, 0, 12, 
	261, 1, 0, 0, 0, 14, 263, 1, 0, 0, 0, 16, 281, 1, 0, 0, 0, 18, 314, 1, 
	0, 0, 0, 20, 316, 1, 0, 0, 0, 22, 323, 1, 0, 0, 0, 24, 339, 1, 0, 0, 0, 
	26, 341, 1, 0, 0, 0, 28, 344, 1, 0, 0, 0, 30, 346, 1, 0, 0, 0, 32, 350, 
	1, 0, 0, 0, 34, 377, 1, 0, 0, 0, 36, 387, 1, 0, 0, 0, 38, 389, 1, 0, 0, 
	0, 40, 399, 1, 0, 0, 0, 42, 401, 1, 0, 0, 0, 44, 410, 1, 0, 0, 0, 46, 412, 
	1, 0, 0, 0, 48, 414, 1, 0, 0, 0, 50, 416, 1, 0, 0, 0, 52, 440, 1, 0, 0, 
	0, 54, 449, 1, 0, 0, 0, 56, 451, 1, 0, 0, 0, 58, 460, 1, 0, 0, 0, 60, 462, 
	1, 0, 0, 0, 62, 471, 1, 0, 0, 0, 64, 473, 1, 0, 0, 0, 66, 481, 1, 0, 0, 
	0, 68, 491, 1, 0, 0, 0, 70, 493, 1, 0, 0, 0, 72, 495, 1, 0, 0, 0, 74, 503, 
	1, 0, 0, 0, 76, 515, 1, 0, 0, 0, 78, 517, 1, 0, 0, 0, 80, 519, 1, 0, 0, 
	0, 82, 523, 1, 0, 0, 0, 84, 527, 1, 0, 0, 0, 86, 531, 1, 0, 0, 0, 88, 537, 
	1, 0, 0, 0, 90, 539, 1, 0, 0, 0, 92, 543, 1, 0, 0, 0, 94, 547, 1, 0, 0, 
	0, 96, 551, 1, 0, 0, 0, 98, 555, 1, 0, 0, 0, 100, 559, 1, 0, 0, 0, 102, 
	563, 1, 0, 0, 0, 104, 567, 1, 0, 0, 0, 106, 571, 1, 0, 0, 0, 108, 575, 
	1, 0, 0, 0, 110, 579, 1, 0, 0, 0, 112, 590, 1, 0, 0, 0, 114, 602, 1, 0, 
	0, 0, 116, 611, 1, 0, 0, 0, 118, 613, 1, 0, 0, 0, 120, 615, 1, 0, 0, 0, 
	122, 625, 1, 0, 0, 0, 124, 635, 1, 0, 0, 0, 126, 637, 1, 0, 0, 0, 128, 
	639, 1, 0, 0, 0, 130, 643, 1, 0, 0, 0, 132, 647, 1, 0, 0, 0, 134, 651, 
	1, 0, 0, 0, 136, 655, 1, 0, 0, 0, 138, 659, 1, 0, 0, 0, 140, 663, 1, 0, 
	0, 0, 142, 667, 1, 0, 0, 0, 144, 671, 1, 0, 0, 0, 146, 675, 1, 0, 0, 0, 
	148, 679, 1, 0, 0, 0, 150, 683, 1, 0, 0, 0, 152, 687, 1, 0, 0, 0, 154, 
	691, 1, 0, 0, 0, 156, 695, 1, 0, 0, 0, 158, 697, 1, 0, 0, 0, 160, 700, 
	1, 0, 0, 0, 162, 704, 1, 0, 0, 0, 164, 166, 5, 44, 0, 0, 165, 164, 1, 0, 
	0, 0, 166, 169, 1, 0, 0, 0, 167, 165, 1, 0, 0, 0, 167, 168, 1, 0, 0, 0, 
	168, 171, 1, 0, 0, 0, 169, 167, 1, 0, 0, 0, 170, 172, 3, 20, 10, 0, 171, 
	170, 1, 0, 0, 0, 171, 172, 1, 0, 0, 0, 172, 176, 1, 0, 0, 0, 173, 174, 
	3, 14, 7, 0, 174, 175, 5, 44, 0, 0, 175, 177, 1, 0, 0, 0, 176, 173, 1, 
	0, 0, 0, 176, 177, 1, 0, 0, 0, 177, 181, 1, 0, 0, 0, 178, 180, 5, 44, 0, 
	0, 179, 178, 1, 0, 0, 0, 180, 183, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 181, 
	182, 1, 0, 0, 0, 182, 185, 1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 184, 186, 
	3, 2, 1, 0, 185, 184, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 190, 1, 0, 
	0, 0, 187, 189, 5, 44, 0, 0, 188, 187, 1, 0, 0, 0, 189, 192, 1, 0, 0, 0, 
	190, 188, 1, 0, 0, 0, 190, 191, 1, 0, 0, 0, 191, 193, 1, 0, 0, 0, 192, 
	190, 1, 0, 0, 0, 193, 194, 5, 0, 0, 1, 194, 1, 1, 0, 0, 0, 195, 205, 3, 
	4, 2, 0, 196, 198, 5, 44, 0, 0, 197, 196, 1, 0, 0, 0, 198, 201, 1, 0, 0, 
	0, 199, 197, 1, 0, 0, 0, 199, 200, 1, 0, 0, 0, 200, 202, 1, 0, 0, 0, 201, 
	199, 1, 0, 0, 0, 202, 204, 3, 4, 2, 0, 203, 199, 1, 0, 0, 0, 204, 207, 
	1, 0, 0, 0, 205, 203, 1, 0, 0, 0, 205, 206, 1, 0, 0, 0, 206, 3, 1, 0, 0, 
	0, 207, 205, 1, 0, 0, 0, 208, 211, 3, 6, 3, 0, 209, 211, 3, 10, 5, 0, 210, 
	208, 1, 0, 0, 0, 210, 209, 1, 0, 0, 0, 211, 5, 1, 0, 0, 0, 212, 214, 3, 
	8, 4, 0, 213, 215, 5, 44, 0, 0, 214, 213, 1, 0, 0, 0, 215, 216, 1, 0, 0, 
	0, 216, 214, 1, 0, 0, 0, 216, 217, 1, 0, 0, 0, 217, 234, 1, 0, 0, 0, 218, 
	222, 5, 1, 0, 0, 219, 221, 5, 44, 0, 0, 220, 219, 1, 0, 0, 0, 221, 224, 
	1, 0, 0, 0, 222, 220, 1, 0, 0, 0, 222, 223, 1, 0, 0, 0, 223, 226, 1, 0, 
	0, 0, 224, 222, 1, 0, 0, 0, 225, 227, 3, 14, 7, 0, 226, 225, 1, 0, 0, 0, 
	226, 227, 1, 0, 0, 0, 227, 229, 1, 0, 0, 0, 228, 230, 5, 44, 0, 0, 229, 
	228, 1, 0, 0, 0, 230, 231, 1, 0, 0, 0, 231, 229, 1, 0, 0, 0, 231, 232, 
	1, 0, 0, 0, 232, 233, 1, 0, 0, 0, 233, 235, 5, 2, 0, 0, 234, 218, 1, 0, 
	0, 0, 234, 235, 1, 0, 0, 0, 235, 7, 1, 0, 0, 0, 236, 237, 5, 3, 0, 0, 237, 
	9, 1, 0, 0, 0, 238, 240, 3, 12, 6, 0, 239, 241, 5, 44, 0, 0, 240, 239, 
	1, 0, 0, 0, 241, 242, 1, 0, 0, 0, 242, 240, 1, 0, 0, 0, 242, 243, 1, 0, 
	0, 0, 243, 259, 1, 0, 0, 0, 244, 248, 5, 1, 0, 0, 245, 247, 5, 44, 0, 0, 
	246, 245, 1, 0, 0, 0, 247, 250, 1, 0, 0, 0, 248, 246, 1, 0, 0, 0, 248, 
	249, 1, 0, 0, 0, 249, 251, 1, 0, 0, 0, 250, 248, 1, 0, 0, 0, 251, 253, 
	3, 14, 7, 0, 252, 254, 5, 44, 0, 0, 253, 252, 1, 0, 0, 0, 254, 255, 1, 
	0, 0, 0, 255, 253, 1, 0, 0, 0, 255, 256, 1, 0, 0, 0, 256, 257, 1, 0, 0, 
	0, 257, 258, 5, 2, 0, 0, 258, 260, 1, 0, 0, 0, 259, 244, 1, 0, 0, 0, 259, 
	260, 1, 0, 0, 0, 260, 11, 1, 0, 0, 0, 261, 262, 5, 4, 0, 0, 262, 13, 1, 
	0, 0, 0, 263, 272, 3, 16, 8, 0, 264, 266, 5, 44, 0, 0, 265, 264, 1, 0, 
	0, 0, 266, 267, 1, 0, 0, 0, 267, 265, 1, 0, 0, 0, 267, 268, 1, 0, 0, 0, 
	268, 269, 1, 0, 0, 0, 269, 271, 3, 16, 8, 0, 270, 265, 1, 0, 0, 0, 271, 
	274, 1, 0, 0, 0, 272, 270, 1, 0, 0, 0, 272, 273, 1, 0, 0, 0, 273, 15, 1, 
	0, 0, 0, 274, 272, 1, 0, 0, 0, 275, 282, 3, 18, 9, 0, 276, 282, 3, 22, 
	11, 0, 277, 282, 3, 32, 16, 0, 278, 282, 3, 64, 32, 0, 279, 282, 3, 72, 
	36, 0, 280, 282, 3, 50, 25, 0, 281, 275, 1, 0, 0, 0, 281, 276, 1, 0, 0, 
	0, 281, 277, 1, 0, 0, 0, 281, 278, 1, 0, 0, 0, 281, 279, 1, 0, 0, 0, 281, 
	280, 1, 0, 0, 0, 282, 17, 1, 0, 0, 0, 283, 315, 3, 80, 40, 0, 284, 315, 
	3, 82, 41, 0, 285, 315, 3, 84, 42, 0, 286, 315, 3, 86, 43, 0, 287, 315, 
	3, 88, 44, 0, 288, 315, 3, 94, 47, 0, 289, 315, 3, 96, 48, 0, 290, 315, 
	3, 98, 49, 0, 291, 315, 3, 100, 50, 0, 292, 315, 3, 102, 51, 0, 293, 315, 
	3, 104, 52, 0, 294, 315, 3, 106, 53, 0, 295, 315, 3, 108, 54, 0, 296, 315, 
	3, 110, 55, 0, 297, 315, 3, 112, 56, 0, 298, 315, 3, 114, 57, 0, 299, 315, 
	3, 122, 61, 0, 300, 315, 3, 128, 64, 0, 301, 315, 3, 130, 65, 0, 302, 315, 
	3, 132, 66, 0, 303, 315, 3, 134, 67, 0, 304, 315, 3, 136, 68, 0, 305, 315, 
	3, 138, 69, 0, 306, 315, 3, 140, 70, 0, 307, 315, 3, 142, 71, 0, 308, 315, 
	3, 144, 72, 0, 309, 315, 3, 146, 73, 0, 310, 315, 3, 148, 74, 0, 311, 315, 
	3, 150, 75, 0, 312, 315, 3, 152, 76, 0, 313, 315, 3, 154, 77, 0, 314, 283, 
	1, 0, 0, 0, 314, 284, 1, 0, 0, 0, 314, 285, 1, 0, 0, 0, 314, 286, 1, 0, 
	0, 0, 314, 287, 1, 0, 0, 0, 314, 288, 1, 0, 0, 0, 314, 289, 1, 0, 0, 0, 
	314, 290, 1, 0, 0, 0, 314, 291, 1, 0, 0, 0, 314, 292, 1, 0, 0, 0, 314, 
	293, 1, 0, 0, 0, 314, 294, 1, 0, 0, 0, 314, 295, 1, 0, 0, 0, 314, 296, 
	1, 0, 0, 0, 314, 297, 1, 0, 0, 0, 314, 298, 1, 0, 0, 0, 314, 299, 1, 0, 
	0, 0, 314, 300, 1, 0, 0, 0, 314, 301, 1, 0, 0, 0, 314, 302, 1, 0, 0, 0, 
	314, 303, 1, 0, 0, 0, 314, 304, 1, 0, 0, 0, 314, 305, 1, 0, 0, 0, 314, 
	306, 1, 0, 0, 0, 314, 307, 1, 0, 0, 0, 314, 308, 1, 0, 0, 0, 314, 309, 
	1, 0, 0, 0, 314, 310, 1, 0, 0, 0, 314, 311, 1, 0, 0, 0, 314, 312, 1, 0, 
	0, 0, 314, 313, 1, 0, 0, 0, 315, 19, 1, 0, 0, 0, 316, 317, 5, 32, 0, 0, 
	317, 319, 3, 160, 80, 0, 318, 320, 5, 44, 0, 0, 319, 318, 1, 0, 0, 0, 320, 
	321, 1, 0, 0, 0, 321, 319, 1, 0, 0, 0, 321, 322, 1, 0, 0, 0, 322, 21, 1, 
	0, 0, 0, 323, 330, 3, 26, 13, 0, 324, 326, 5, 44, 0, 0, 325, 324, 1, 0, 
	0, 0, 326, 327, 1, 0, 0, 0, 327, 325, 1, 0, 0, 0, 327, 328, 1, 0, 0, 0, 
	328, 329, 1, 0, 0, 0, 329, 331, 3, 24, 12, 0, 330, 325, 1, 0, 0, 0, 330, 
	331, 1, 0, 0, 0, 331, 333, 1, 0, 0, 0, 332, 334, 5, 44, 0, 0, 333, 332, 
	1, 0, 0, 0, 334, 335, 1, 0, 0, 0, 335, 333, 1, 0, 0, 0, 335, 336, 1, 0, 
	0, 0, 336, 337, 1, 0, 0, 0, 337, 338, 5, 51, 0, 0, 338, 23, 1, 0, 0, 0, 
	339, 340, 3, 14, 7, 0, 340, 25, 1, 0, 0, 0, 341, 342, 5, 38, 0, 0, 342, 
	343, 3, 28, 14, 0, 343, 27, 1, 0, 0, 0, 344, 345, 3, 30, 15, 0, 345, 29, 
	1, 0, 0, 0, 346, 348, 5, 39, 0, 0, 347, 349, 3, 160, 80, 0, 348, 347, 1, 
	0, 0, 0, 348, 349, 1, 0, 0, 0, 349, 31, 1, 0, 0, 0, 350, 359, 3, 34, 17, 
	0, 351, 353, 5, 44, 0, 0, 352, 351, 1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 
	352, 1, 0, 0, 0, 354, 355, 1, 0, 0, 0, 355, 356, 1, 0, 0, 0, 356, 358, 
	3, 38, 19, 0, 357, 352, 1, 0, 0, 0, 358, 361, 1, 0, 0, 0, 359, 357, 1, 
	0, 0, 0, 359, 360, 1, 0, 0, 0, 360, 368, 1, 0, 0, 0, 361, 359, 1, 0, 0, 
	0, 362, 364, 5, 44, 0, 0, 363, 362, 1, 0, 0, 0, 364, 365, 1, 0, 0, 0, 365, 
	363, 1, 0, 0, 0, 365, 366, 1, 0, 0, 0, 366, 367, 1, 0, 0, 0, 367, 369, 
	3, 42, 21, 0, 368, 363, 1, 0, 0, 0, 368, 369, 1, 0, 0, 0, 369, 371, 1, 
	0, 0, 0, 370, 372, 5, 44, 0, 0, 371, 370, 1, 0, 0, 0, 372, 373, 1, 0, 0, 
	0, 373, 371, 1, 0, 0, 0, 373, 374, 1, 0, 0, 0, 374, 375, 1, 0, 0, 0, 375, 
	376, 5, 51, 0, 0, 376, 33, 1, 0, 0, 0, 377, 378, 5, 40, 0, 0, 378, 385, 
	3, 46, 23, 0, 379, 381, 5, 44, 0, 0, 380, 379, 1, 0, 0, 0, 381, 382, 1, 
	0, 0, 0, 382, 380, 1, 0, 0, 0, 382, 383, 1, 0, 0, 0, 383, 384, 1, 0, 0, 
	0, 384, 386, 3, 36, 18, 0, 385, 380, 1, 0, 0, 0, 385, 386, 1, 0, 0, 0, 
	386, 35, 1, 0, 0, 0, 387, 388, 3, 14, 7, 0, 388, 37, 1, 0, 0, 0, 389, 390, 
	5, 48, 0, 0, 390, 397, 3, 48, 24, 0, 391, 393, 5, 44, 0, 0, 392, 391, 1, 
	0, 0, 0, 393, 394, 1, 0, 0, 0, 394, 392, 1, 0, 0, 0, 394, 395, 1, 0, 0, 
	0, 395, 396, 1, 0, 0, 0, 396, 398, 3, 40, 20, 0, 397, 392, 1, 0, 0, 0, 
	397, 398, 1, 0, 0, 0, 398, 39, 1, 0, 0, 0, 399, 400, 3, 14, 7, 0, 400, 
	41, 1, 0, 0, 0, 401, 408, 5, 47, 0, 0, 402, 404, 5, 44, 0, 0, 403, 402, 
	1, 0, 0, 0, 404, 405, 1, 0, 0, 0, 405, 403, 1, 0, 0, 0, 405, 406, 1, 0, 
	0, 0, 406, 407, 1, 0, 0, 0, 407, 409, 3, 44, 22, 0, 408, 403, 1, 0, 0, 
	0, 408, 409, 1, 0, 0, 0, 409, 43, 1, 0, 0, 0, 410, 411, 3, 14, 7, 0, 411, 
	45, 1, 0, 0, 0, 412, 413, 3, 156, 78, 0, 413, 47, 1, 0, 0, 0, 414, 415, 
	3, 156, 78, 0, 415, 49, 1, 0, 0, 0, 416, 423, 3, 52, 26, 0, 417, 419, 5, 
	44, 0, 0, 418, 417, 1, 0, 0, 0, 419, 420, 1, 0, 0, 0, 420, 418, 1, 0, 0, 
	0, 420, 421, 1, 0, 0, 0, 421, 422, 1, 0, 0, 0, 422, 424, 3, 56, 28, 0, 
	423, 418, 1, 0, 0, 0, 423, 424, 1, 0, 0, 0, 424, 431, 1, 0, 0, 0, 425, 
	427, 5, 44, 0, 0, 426, 425, 1, 0, 0, 0, 427, 428, 1, 0, 0, 0, 428, 426, 
	1, 0, 0, 0, 428, 429, 1, 0, 0, 0, 429, 430, 1, 0, 0, 0, 430, 432, 3, 60, 
	30, 0, 431, 426, 1, 0, 0, 0, 431, 432, 1, 0, 0, 0, 432, 434, 1, 0, 0, 0, 
	433, 435, 5, 44, 0, 0, 434, 433, 1, 0, 0, 0, 435, 436, 1, 0, 0, 0, 436, 
	434, 1, 0, 0, 0, 436, 437, 1, 0, 0, 0, 437, 438, 1, 0, 0, 0, 438, 439, 
	5, 51, 0, 0, 439, 51, 1, 0, 0, 0, 440, 447, 5, 41, 0, 0, 441, 443, 5, 44, 
	0, 0, 442, 441, 1, 0, 0, 0, 443, 444, 1, 0, 0, 0, 444, 442, 1, 0, 0, 0, 
	444, 445, 1, 0, 0, 0, 445, 446, 1, 0, 0, 0, 446, 448, 3, 54, 27, 0, 447, 
	442, 1, 0, 0, 0, 447, 448, 1, 0, 0, 0, 448, 53, 1, 0, 0, 0, 449, 450, 3, 
	14, 7, 0, 450, 55, 1, 0, 0, 0, 451, 458, 5, 49, 0, 0, 452, 454, 5, 44, 
	0, 0, 453, 452, 1, 0, 0, 0, 454, 455, 1, 0, 0, 0, 455, 453, 1, 0, 0, 0, 
	455, 456, 1, 0, 0, 0, 456, 457, 1, 0, 0, 0, 457, 459, 3, 58, 29, 0, 458, 
	453, 1, 0, 0, 0, 458, 459, 1, 0, 0, 0, 459, 57, 1, 0, 0, 0, 460, 461, 3, 
	14, 7, 0, 461, 59, 1, 0, 0, 0, 462, 469, 5, 50, 0, 0, 463, 465, 5, 44, 
	0, 0, 464, 463, 1, 0, 0, 0, 465, 466, 1, 0, 0, 0, 466, 464, 1, 0, 0, 0, 
	466, 467, 1, 0, 0, 0, 467, 468, 1, 0, 0, 0, 468, 470, 3, 62, 31, 0, 469, 
	464, 1, 0, 0, 0, 469, 470, 1, 0, 0, 0, 470, 61, 1, 0, 0, 0, 471, 472, 3, 
	14, 7, 0, 472, 63, 1, 0, 0, 0, 473, 475, 3, 66, 33, 0, 474, 476, 5, 44, 
	0, 0, 475, 474, 1, 0, 0, 0, 476, 477, 1, 0, 0, 0, 477, 475, 1, 0, 0, 0, 
	477, 478, 1, 0, 0, 0, 478, 479, 1, 0, 0, 0, 479, 480, 5, 51, 0, 0, 480, 
	65, 1, 0, 0, 0, 481, 482, 5, 42, 0, 0, 482, 489, 3, 70, 35, 0, 483, 485, 
	5, 44, 0, 0, 484, 483, 1, 0, 0, 0, 485, 486, 1, 0, 0, 0, 486, 484, 1, 0, 
	0, 0, 486, 487, 1, 0, 0, 0, 487, 488, 1, 0, 0, 0, 488, 490, 3, 68, 34, 
	0, 489, 484, 1, 0, 0, 0, 489, 490, 1, 0, 0, 0, 490, 67, 1, 0, 0, 0, 491, 
	492, 3, 14, 7, 0, 492, 69, 1, 0, 0, 0, 493, 494, 3, 160, 80, 0, 494, 71, 
	1, 0, 0, 0, 495, 497, 3, 74, 37, 0, 496, 498, 5, 44, 0, 0, 497, 496, 1, 
	0, 0, 0, 498, 499, 1, 0, 0, 0, 499, 497, 1, 0, 0, 0, 499, 500, 1, 0, 0, 
	0, 500, 501, 1, 0, 0, 0, 501, 502, 5, 51, 0, 0, 502, 73, 1, 0, 0, 0, 503, 
	505, 5, 43, 0, 0, 504, 506, 3, 78, 39, 0, 505, 504, 1, 0, 0, 0, 505, 506, 
	1, 0, 0, 0, 506, 513, 1, 0, 0, 0, 507, 509, 5, 44, 0, 0, 508, 507, 1, 0, 
	0, 0, 509, 510, 1, 0, 0, 0, 510, 508, 1, 0, 0, 0, 510, 511, 1, 0, 0, 0, 
	511, 512, 1, 0, 0, 0, 512, 514, 3, 76, 38, 0, 513, 508, 1, 0, 0, 0, 513, 
	514, 1, 0, 0, 0, 514, 75, 1, 0, 0, 0, 515, 516, 3, 14, 7, 0, 516, 77, 1, 
	0, 0, 0, 517, 518, 3, 160, 80, 0, 518, 79, 1, 0, 0, 0, 519, 521, 5, 5, 
	0, 0, 520, 522, 3, 160, 80, 0, 521, 520, 1, 0, 0, 0, 521, 522, 1, 0, 0, 
	0, 522, 81, 1, 0, 0, 0, 523, 525, 5, 6, 0, 0, 524, 526, 3, 160, 80, 0, 
	525, 524, 1, 0, 0, 0, 525, 526, 1, 0, 0, 0, 526, 83, 1, 0, 0, 0, 527, 529, 
	5, 7, 0, 0, 528, 530, 3, 160, 80, 0, 529, 528, 1, 0, 0, 0, 529, 530, 1, 
	0, 0, 0, 530, 85, 1, 0, 0, 0, 531, 533, 5, 8, 0, 0, 532, 534, 3, 160, 80, 
	0, 533, 532, 1, 0, 0, 0, 533, 534, 1, 0, 0, 0, 534, 87, 1, 0, 0, 0, 535, 
	538, 3, 92, 46, 0, 536, 538, 3, 90, 45, 0, 537, 535, 1, 0, 0, 0, 537, 536, 
	1, 0, 0, 0, 538, 89, 1, 0, 0, 0, 539, 541, 5, 10, 0, 0, 540, 542, 3, 160, 
	80, 0, 541, 540, 1, 0, 0, 0, 541, 542, 1, 0, 0, 0, 542, 91, 1, 0, 0, 0, 
	543, 545, 5, 9, 0, 0, 544, 546, 3, 160, 80, 0, 545, 544, 1, 0, 0, 0, 545, 
	546, 1, 0, 0, 0, 546, 93, 1, 0, 0, 0, 547, 549, 5, 11, 0, 0, 548, 550, 
	3, 158, 79, 0, 549, 548, 1, 0, 0, 0, 549, 550, 1, 0, 0, 0, 550, 95, 1, 
	0, 0, 0, 551, 553, 5, 18, 0, 0, 552, 554, 3, 160, 80, 0, 553, 552, 1, 0, 
	0, 0, 553, 554, 1, 0, 0, 0, 554, 97, 1, 0, 0, 0, 555, 557, 5, 19, 0, 0, 
	556, 558, 3, 160, 80, 0, 557, 556, 1, 0, 0, 0, 557, 558, 1, 0, 0, 0, 558, 
	99, 1, 0, 0, 0, 559, 561, 5, 20, 0, 0, 560, 562, 3, 160, 80, 0, 561, 560, 
	1, 0, 0, 0, 561, 562, 1, 0, 0, 0, 562, 101, 1, 0, 0, 0, 563, 565, 5, 21, 
	0, 0, 564, 566, 3, 158, 79, 0, 565, 564, 1, 0, 0, 0, 565, 566, 1, 0, 0, 
	0, 566, 103, 1, 0, 0, 0, 567, 569, 5, 22, 0, 0, 568, 570, 3, 158, 79, 0, 
	569, 568, 1, 0, 0, 0, 569, 570, 1, 0, 0, 0, 570, 105, 1, 0, 0, 0, 571, 
	573, 5, 12, 0, 0, 572, 574, 3, 160, 80, 0, 573, 572, 1, 0, 0, 0, 573, 574, 
	1, 0, 0, 0, 574, 107, 1, 0, 0, 0, 575, 577, 5, 13, 0, 0, 576, 578, 3, 158, 
	79, 0, 577, 576, 1, 0, 0, 0, 577, 578, 1, 0, 0, 0, 578, 109, 1, 0, 0, 0, 
	579, 580, 5, 14, 0, 0, 580, 582, 3, 118, 59, 0, 581, 583, 5, 53, 0, 0, 
	582, 581, 1, 0, 0, 0, 582, 583, 1, 0, 0, 0, 583, 588, 1, 0, 0, 0, 584, 
	586, 5, 45, 0, 0, 585, 584, 1, 0, 0, 0, 585, 586, 1, 0, 0, 0, 586, 587, 
	1, 0, 0, 0, 587, 589, 3, 120, 60, 0, 588, 585, 1, 0, 0, 0, 588, 589, 1, 
	0, 0, 0, 589, 111, 1, 0, 0, 0, 590, 591, 5, 15, 0, 0, 591, 592, 3, 116, 
	58, 0, 592, 600, 3, 118, 59, 0, 593, 598, 5, 53, 0, 0, 594, 596, 5, 45, 
	0, 0, 595, 594, 1, 0, 0, 0, 595, 596, 1, 0, 0, 0, 596, 597, 1, 0, 0, 0, 
	597, 599, 3, 120, 60, 0, 598, 595, 1, 0, 0, 0, 598, 599, 1, 0, 0, 0, 599, 
	601, 1, 0, 0, 0, 600, 593, 1, 0, 0, 0, 600, 601, 1, 0, 0, 0, 601, 113, 
	1, 0, 0, 0, 602, 603, 5, 16, 0, 0, 603, 604, 3, 118, 59, 0, 604, 606, 5, 
	53, 0, 0, 605, 607, 5, 45, 0, 0, 606, 605, 1, 0, 0, 0, 606, 607, 1, 0, 
	0, 0, 607, 608, 1, 0, 0, 0, 608, 609, 3, 120, 60, 0, 609, 115, 1, 0, 0, 
	0, 610, 612, 3, 160, 80, 0, 611, 610, 1, 0, 0, 0, 611, 612, 1, 0, 0, 0, 
	612, 117, 1, 0, 0, 0, 613, 614, 5, 52, 0, 0, 614, 119, 1, 0, 0, 0, 615, 
	622, 5, 52, 0, 0, 616, 618, 5, 45, 0, 0, 617, 616, 1, 0, 0, 0, 617, 618, 
	1, 0, 0, 0, 618, 619, 1, 0, 0, 0, 619, 621, 5, 52, 0, 0, 620, 617, 1, 0, 
	0, 0, 621, 624, 1, 0, 0, 0, 622, 620, 1, 0, 0, 0, 622, 623, 1, 0, 0, 0, 
	623, 121, 1, 0, 0, 0, 624, 622, 1, 0, 0, 0, 625, 632, 5, 17, 0, 0, 626, 
	627, 3, 124, 62, 0, 627, 628, 5, 53, 0, 0, 628, 629, 3, 126, 63, 0, 629, 
	631, 1, 0, 0, 0, 630, 626, 1, 0, 0, 0, 631, 634, 1, 0, 0, 0, 632, 630, 
	1, 0, 0, 0, 632, 633, 1, 0, 0, 0, 633, 123, 1, 0, 0, 0, 634, 632, 1, 0, 
	0, 0, 635, 636, 5, 52, 0, 0, 636, 125, 1, 0, 0, 0, 637, 638, 5, 52, 0, 
	0, 638, 127, 1, 0, 0, 0, 639, 641, 5, 23, 0, 0, 640, 642, 3, 160, 80, 0, 
	641, 640, 1, 0, 0, 0, 641, 642, 1, 0, 0, 0, 642, 129, 1, 0, 0, 0, 643, 
	645, 5, 24, 0, 0, 644, 646, 3, 160, 80, 0, 645, 644, 1, 0, 0, 0, 645, 646, 
	1, 0, 0, 0, 646, 131, 1, 0, 0, 0, 647, 649, 5, 25, 0, 0, 648, 650, 3, 160, 
	80, 0, 649, 648, 1, 0, 0, 0, 649, 650, 1, 0, 0, 0, 650, 133, 1, 0, 0, 0, 
	651, 653, 5, 26, 0, 0, 652, 654, 3, 160, 80, 0, 653, 652, 1, 0, 0, 0, 653, 
	654, 1, 0, 0, 0, 654, 135, 1, 0, 0, 0, 655, 657, 5, 27, 0, 0, 656, 658, 
	3, 160, 80, 0, 657, 656, 1, 0, 0, 0, 657, 658, 1, 0, 0, 0, 658, 137, 1, 
	0, 0, 0, 659, 661, 5, 28, 0, 0, 660, 662, 3, 160, 80, 0, 661, 660, 1, 0, 
	0, 0, 661, 662, 1, 0, 0, 0, 662, 139, 1, 0, 0, 0, 663, 665, 5, 30, 0, 0, 
	664, 666, 3, 160, 80, 0, 665, 664, 1, 0, 0, 0, 665, 666, 1, 0, 0, 0, 666, 
	141, 1, 0, 0, 0, 667, 669, 5, 29, 0, 0, 668, 670, 3, 160, 80, 0, 669, 668, 
	1, 0, 0, 0, 669, 670, 1, 0, 0, 0, 670, 143, 1, 0, 0, 0, 671, 673, 5, 31, 
	0, 0, 672, 674, 3, 160, 80, 0, 673, 672, 1, 0, 0, 0, 673, 674, 1, 0, 0, 
	0, 674, 145, 1, 0, 0, 0, 675, 677, 5, 33, 0, 0, 676, 678, 3, 160, 80, 0, 
	677, 676, 1, 0, 0, 0, 677, 678, 1, 0, 0, 0, 678, 147, 1, 0, 0, 0, 679, 
	681, 5, 34, 0, 0, 680, 682, 3, 160, 80, 0, 681, 680, 1, 0, 0, 0, 681, 682, 
	1, 0, 0, 0, 682, 149, 1, 0, 0, 0, 683, 685, 5, 35, 0, 0, 684, 686, 3, 160, 
	80, 0, 685, 684, 1, 0, 0, 0, 685, 686, 1, 0, 0, 0, 686, 151, 1, 0, 0, 0, 
	687, 689, 5, 36, 0, 0, 688, 690, 3, 160, 80, 0, 689, 688, 1, 0, 0, 0, 689, 
	690, 1, 0, 0, 0, 690, 153, 1, 0, 0, 0, 691, 693, 5, 37, 0, 0, 692, 694, 
	3, 160, 80, 0, 693, 692, 1, 0, 0, 0, 693, 694, 1, 0, 0, 0, 694, 155, 1, 
	0, 0, 0, 695, 696, 3, 158, 79, 0, 696, 157, 1, 0, 0, 0, 697, 698, 3, 160, 
	80, 0, 698, 159, 1, 0, 0, 0, 699, 701, 3, 162, 81, 0, 700, 699, 1, 0, 0, 
	0, 701, 702, 1, 0, 0, 0, 702, 700, 1, 0, 0, 0, 702, 703, 1, 0, 0, 0, 703, 
	161, 1, 0, 0, 0, 704, 705, 5, 52, 0, 0, 705, 163, 1, 0, 0, 0, 97, 167, 
	171, 176, 181, 185, 190, 199, 205, 210, 216, 222, 226, 231, 234, 242, 248, 
	255, 259, 267, 272, 281, 314, 321, 327, 330, 335, 348, 354, 359, 365, 368, 
	373, 382, 385, 394, 397, 405, 408, 420, 423, 428, 431, 436, 444, 447, 455, 
	458, 466, 469, 477, 486, 489, 499, 505, 510, 513, 521, 525, 529, 533, 537, 
	541, 545, 549, 553, 557, 561, 565, 569, 573, 577, 582, 585, 588, 595, 598, 
	600, 606, 611, 617, 622, 632, 641, 645, 649, 653, 657, 661, 665, 669, 673, 
	677, 681, 685, 689, 693, 702,
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
	EarthParserSET = 16
	EarthParserLABEL = 17
	EarthParserBUILD = 18
	EarthParserWORKDIR = 19
	EarthParserUSER = 20
	EarthParserCMD = 21
	EarthParserENTRYPOINT = 22
	EarthParserGIT_CLONE = 23
	EarthParserADD = 24
	EarthParserSTOPSIGNAL = 25
	EarthParserONBUILD = 26
	EarthParserHEALTHCHECK = 27
	EarthParserSHELL = 28
	EarthParserDO = 29
	EarthParserCOMMAND = 30
	EarthParserIMPORT = 31
	EarthParserVERSION = 32
	EarthParserCACHE = 33
	EarthParserHOST = 34
	EarthParserPROJECT = 35
	EarthParserPIPELINE = 36
	EarthParserTRIGGER = 37
	EarthParserWITH = 38
	EarthParserDOCKER = 39
	EarthParserIF = 40
	EarthParserTRY = 41
	EarthParserFOR = 42
	EarthParserWAIT = 43
	EarthParserNL = 44
	EarthParserWS = 45
	EarthParserCOMMENT = 46
	EarthParserELSE = 47
	EarthParserELSE_IF = 48
	EarthParserCATCH = 49
	EarthParserFINALLY = 50
	EarthParserEND = 51
	EarthParserAtom = 52
	EarthParserEQUALS = 53
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
	EarthParserRULE_setStmt = 57
	EarthParserRULE_optionalFlag = 58
	EarthParserRULE_envArgKey = 59
	EarthParserRULE_envArgValue = 60
	EarthParserRULE_labelStmt = 61
	EarthParserRULE_labelKey = 62
	EarthParserRULE_labelValue = 63
	EarthParserRULE_gitCloneStmt = 64
	EarthParserRULE_addStmt = 65
	EarthParserRULE_stopsignalStmt = 66
	EarthParserRULE_onbuildStmt = 67
	EarthParserRULE_healthcheckStmt = 68
	EarthParserRULE_shellStmt = 69
	EarthParserRULE_userCommandStmt = 70
	EarthParserRULE_doStmt = 71
	EarthParserRULE_importStmt = 72
	EarthParserRULE_cacheStmt = 73
	EarthParserRULE_hostStmt = 74
	EarthParserRULE_projectStmt = 75
	EarthParserRULE_pipelineStmt = 76
	EarthParserRULE_triggerStmt = 77
	EarthParserRULE_expr = 78
	EarthParserRULE_stmtWordsMaybeJSON = 79
	EarthParserRULE_stmtWords = 80
	EarthParserRULE_stmtWord = 81
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
	p.SetState(167)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(164)
				p.Match(EarthParserNL)
			}


		}
		p.SetState(169)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserVERSION {
		{
			p.SetState(170)
			p.Version()
		}

	}
	p.SetState(176)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 17038135263200) != 0) {
		{
			p.SetState(173)
			p.Stmts()
		}
		{
			p.SetState(174)
			p.Match(EarthParserNL)
		}

	}
	p.SetState(181)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(178)
				p.Match(EarthParserNL)
			}


		}
		p.SetState(183)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}
	p.SetState(185)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserTarget || _la == EarthParserUserCommand {
		{
			p.SetState(184)
			p.Targets()
		}

	}
	p.SetState(190)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserNL {
		{
			p.SetState(187)
			p.Match(EarthParserNL)
		}


		p.SetState(192)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(193)
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
		p.SetState(195)
		p.TargetOrUserCommand()
	}
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(199)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for _la == EarthParserNL {
				{
					p.SetState(196)
					p.Match(EarthParserNL)
				}


				p.SetState(201)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(202)
				p.TargetOrUserCommand()
			}


		}
		p.SetState(207)
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

	p.SetState(210)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserTarget:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(208)
			p.Target()
		}


	case EarthParserUserCommand:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(209)
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
		p.SetState(212)
		p.TargetHeader()
	}
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(213)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}
	p.SetState(234)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserINDENT {
		{
			p.SetState(218)
			p.Match(EarthParserINDENT)
		}
		p.SetState(222)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(219)
					p.Match(EarthParserNL)
				}


			}
			p.SetState(224)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
		}
		p.SetState(226)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if ((int64(_la) & ^0x3f) == 0 && ((int64(1) << _la) & 17038135263200) != 0) {
			{
				p.SetState(225)
				p.Stmts()
			}

		}
		p.SetState(229)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(228)
				p.Match(EarthParserNL)
			}


			p.SetState(231)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(233)
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
		p.SetState(236)
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
		p.SetState(238)
		p.UserCommandHeader()
	}
	p.SetState(240)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(239)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(242)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext())
	}
	p.SetState(259)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserINDENT {
		{
			p.SetState(244)
			p.Match(EarthParserINDENT)
		}
		p.SetState(248)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == EarthParserNL {
			{
				p.SetState(245)
				p.Match(EarthParserNL)
			}


			p.SetState(250)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(251)
			p.Stmts()
		}
		p.SetState(253)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(252)
				p.Match(EarthParserNL)
			}


			p.SetState(255)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(257)
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
		p.SetState(261)
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
		p.SetState(263)
		p.Stmt()
	}
	p.SetState(272)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(265)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(264)
					p.Match(EarthParserNL)
				}


				p.SetState(267)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(269)
				p.Stmt()
			}


		}
		p.SetState(274)
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

	p.SetState(281)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM, EarthParserFROM_DOCKERFILE, EarthParserLOCALLY, EarthParserCOPY, EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE, EarthParserRUN, EarthParserEXPOSE, EarthParserVOLUME, EarthParserENV, EarthParserARG, EarthParserSET, EarthParserLABEL, EarthParserBUILD, EarthParserWORKDIR, EarthParserUSER, EarthParserCMD, EarthParserENTRYPOINT, EarthParserGIT_CLONE, EarthParserADD, EarthParserSTOPSIGNAL, EarthParserONBUILD, EarthParserHEALTHCHECK, EarthParserSHELL, EarthParserDO, EarthParserCOMMAND, EarthParserIMPORT, EarthParserCACHE, EarthParserHOST, EarthParserPROJECT, EarthParserPIPELINE, EarthParserTRIGGER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(275)
			p.CommandStmt()
		}


	case EarthParserWITH:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(276)
			p.WithStmt()
		}


	case EarthParserIF:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(277)
			p.IfStmt()
		}


	case EarthParserFOR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(278)
			p.ForStmt()
		}


	case EarthParserWAIT:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(279)
			p.WaitStmt()
		}


	case EarthParserTRY:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(280)
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
	SetStmt() ISetStmtContext
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

func (s *CommandStmtContext) SetStmt() ISetStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetStmtContext)
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

	p.SetState(314)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(283)
			p.FromStmt()
		}


	case EarthParserFROM_DOCKERFILE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(284)
			p.FromDockerfileStmt()
		}


	case EarthParserLOCALLY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(285)
			p.LocallyStmt()
		}


	case EarthParserCOPY:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(286)
			p.CopyStmt()
		}


	case EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(287)
			p.SaveStmt()
		}


	case EarthParserRUN:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(288)
			p.RunStmt()
		}


	case EarthParserBUILD:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(289)
			p.BuildStmt()
		}


	case EarthParserWORKDIR:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(290)
			p.WorkdirStmt()
		}


	case EarthParserUSER:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(291)
			p.UserStmt()
		}


	case EarthParserCMD:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(292)
			p.CmdStmt()
		}


	case EarthParserENTRYPOINT:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(293)
			p.EntrypointStmt()
		}


	case EarthParserEXPOSE:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(294)
			p.ExposeStmt()
		}


	case EarthParserVOLUME:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(295)
			p.VolumeStmt()
		}


	case EarthParserENV:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(296)
			p.EnvStmt()
		}


	case EarthParserARG:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(297)
			p.ArgStmt()
		}


	case EarthParserSET:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(298)
			p.SetStmt()
		}


	case EarthParserLABEL:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(299)
			p.LabelStmt()
		}


	case EarthParserGIT_CLONE:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(300)
			p.GitCloneStmt()
		}


	case EarthParserADD:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(301)
			p.AddStmt()
		}


	case EarthParserSTOPSIGNAL:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(302)
			p.StopsignalStmt()
		}


	case EarthParserONBUILD:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(303)
			p.OnbuildStmt()
		}


	case EarthParserHEALTHCHECK:
		p.EnterOuterAlt(localctx, 22)
		{
			p.SetState(304)
			p.HealthcheckStmt()
		}


	case EarthParserSHELL:
		p.EnterOuterAlt(localctx, 23)
		{
			p.SetState(305)
			p.ShellStmt()
		}


	case EarthParserCOMMAND:
		p.EnterOuterAlt(localctx, 24)
		{
			p.SetState(306)
			p.UserCommandStmt()
		}


	case EarthParserDO:
		p.EnterOuterAlt(localctx, 25)
		{
			p.SetState(307)
			p.DoStmt()
		}


	case EarthParserIMPORT:
		p.EnterOuterAlt(localctx, 26)
		{
			p.SetState(308)
			p.ImportStmt()
		}


	case EarthParserCACHE:
		p.EnterOuterAlt(localctx, 27)
		{
			p.SetState(309)
			p.CacheStmt()
		}


	case EarthParserHOST:
		p.EnterOuterAlt(localctx, 28)
		{
			p.SetState(310)
			p.HostStmt()
		}


	case EarthParserPROJECT:
		p.EnterOuterAlt(localctx, 29)
		{
			p.SetState(311)
			p.ProjectStmt()
		}


	case EarthParserPIPELINE:
		p.EnterOuterAlt(localctx, 30)
		{
			p.SetState(312)
			p.PipelineStmt()
		}


	case EarthParserTRIGGER:
		p.EnterOuterAlt(localctx, 31)
		{
			p.SetState(313)
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
		p.SetState(316)
		p.Match(EarthParserVERSION)
	}
	{
		p.SetState(317)
		p.StmtWords()
	}
	p.SetState(319)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(318)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(321)
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
		p.SetState(323)
		p.WithExpr()
	}
	p.SetState(330)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
		p.SetState(325)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(324)
				p.Match(EarthParserNL)
			}


			p.SetState(327)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(329)
			p.WithBlock()
		}


	}
	p.SetState(333)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(332)
			p.Match(EarthParserNL)
		}


		p.SetState(335)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(337)
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
		p.SetState(339)
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
		p.SetState(341)
		p.Match(EarthParserWITH)
	}
	{
		p.SetState(342)
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
		p.SetState(344)
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
		p.SetState(346)
		p.Match(EarthParserDOCKER)
	}
	p.SetState(348)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(347)
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
		p.SetState(350)
		p.IfClause()
	}
	p.SetState(359)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(352)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(351)
					p.Match(EarthParserNL)
				}


				p.SetState(354)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(356)
				p.ElseIfClause()
			}


		}
		p.SetState(361)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())
	}
	p.SetState(368)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) == 1 {
		p.SetState(363)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(362)
				p.Match(EarthParserNL)
			}


			p.SetState(365)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(367)
			p.ElseClause()
		}


	}
	p.SetState(371)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(370)
			p.Match(EarthParserNL)
		}


		p.SetState(373)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(375)
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
		p.SetState(377)
		p.Match(EarthParserIF)
	}
	{
		p.SetState(378)
		p.IfExpr()
	}
	p.SetState(385)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
		p.SetState(380)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(379)
				p.Match(EarthParserNL)
			}


			p.SetState(382)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(384)
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
		p.SetState(387)
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
		p.SetState(389)
		p.Match(EarthParserELSE_IF)
	}
	{
		p.SetState(390)
		p.ElseIfExpr()
	}
	p.SetState(397)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext()) == 1 {
		p.SetState(392)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(391)
				p.Match(EarthParserNL)
			}


			p.SetState(394)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(396)
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
		p.SetState(399)
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
		p.SetState(401)
		p.Match(EarthParserELSE)
	}
	p.SetState(408)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext()) == 1 {
		p.SetState(403)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(402)
				p.Match(EarthParserNL)
			}


			p.SetState(405)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(407)
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
		p.SetState(410)
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
		p.SetState(412)
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
		p.SetState(414)
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
		p.SetState(416)
		p.TryClause()
	}
	p.SetState(423)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext()) == 1 {
		p.SetState(418)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(417)
				p.Match(EarthParserNL)
			}


			p.SetState(420)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(422)
			p.CatchClause()
		}


	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) == 1 {
		p.SetState(426)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(425)
				p.Match(EarthParserNL)
			}


			p.SetState(428)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(430)
			p.FinallyClause()
		}


	}
	p.SetState(434)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(433)
			p.Match(EarthParserNL)
		}


		p.SetState(436)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(438)
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
		p.SetState(440)
		p.Match(EarthParserTRY)
	}
	p.SetState(447)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
		p.SetState(442)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(441)
				p.Match(EarthParserNL)
			}


			p.SetState(444)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(446)
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
		p.SetState(449)
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
		p.SetState(451)
		p.Match(EarthParserCATCH)
	}
	p.SetState(458)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
		p.SetState(453)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(452)
				p.Match(EarthParserNL)
			}


			p.SetState(455)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(457)
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
		p.SetState(460)
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
		p.SetState(462)
		p.Match(EarthParserFINALLY)
	}
	p.SetState(469)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext()) == 1 {
		p.SetState(464)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(463)
				p.Match(EarthParserNL)
			}


			p.SetState(466)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(468)
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
		p.SetState(471)
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
		p.SetState(473)
		p.ForClause()
	}
	p.SetState(475)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(474)
			p.Match(EarthParserNL)
		}


		p.SetState(477)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(479)
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
		p.SetState(481)
		p.Match(EarthParserFOR)
	}
	{
		p.SetState(482)
		p.ForExpr()
	}
	p.SetState(489)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) == 1 {
		p.SetState(484)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(483)
				p.Match(EarthParserNL)
			}


			p.SetState(486)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(488)
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
		p.SetState(491)
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
		p.SetState(493)
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
		p.SetState(495)
		p.WaitClause()
	}
	p.SetState(497)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(496)
			p.Match(EarthParserNL)
		}


		p.SetState(499)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(501)
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
		p.SetState(503)
		p.Match(EarthParserWAIT)
	}
	p.SetState(505)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(504)
			p.WaitExpr()
		}

	}
	p.SetState(513)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) == 1 {
		p.SetState(508)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(507)
				p.Match(EarthParserNL)
			}


			p.SetState(510)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(512)
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
		p.SetState(515)
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
		p.SetState(517)
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
		p.SetState(519)
		p.Match(EarthParserFROM)
	}
	p.SetState(521)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(520)
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
		p.SetState(523)
		p.Match(EarthParserFROM_DOCKERFILE)
	}
	p.SetState(525)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(524)
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
		p.SetState(527)
		p.Match(EarthParserLOCALLY)
	}
	p.SetState(529)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(528)
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
		p.SetState(531)
		p.Match(EarthParserCOPY)
	}
	p.SetState(533)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(532)
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

	p.SetState(537)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserSAVE_ARTIFACT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(535)
			p.SaveArtifact()
		}


	case EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(536)
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
		p.SetState(539)
		p.Match(EarthParserSAVE_IMAGE)
	}
	p.SetState(541)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(540)
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
		p.SetState(543)
		p.Match(EarthParserSAVE_ARTIFACT)
	}
	p.SetState(545)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(544)
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
		p.SetState(547)
		p.Match(EarthParserRUN)
	}
	p.SetState(549)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(548)
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
		p.SetState(551)
		p.Match(EarthParserBUILD)
	}
	p.SetState(553)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(552)
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
		p.SetState(555)
		p.Match(EarthParserWORKDIR)
	}
	p.SetState(557)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(556)
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
		p.SetState(559)
		p.Match(EarthParserUSER)
	}
	p.SetState(561)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(560)
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
		p.SetState(563)
		p.Match(EarthParserCMD)
	}
	p.SetState(565)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(564)
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
		p.SetState(567)
		p.Match(EarthParserENTRYPOINT)
	}
	p.SetState(569)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(568)
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
		p.SetState(571)
		p.Match(EarthParserEXPOSE)
	}
	p.SetState(573)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(572)
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
		p.SetState(575)
		p.Match(EarthParserVOLUME)
	}
	p.SetState(577)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(576)
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
		p.SetState(579)
		p.Match(EarthParserENV)
	}
	{
		p.SetState(580)
		p.EnvArgKey()
	}
	p.SetState(582)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserEQUALS {
		{
			p.SetState(581)
			p.Match(EarthParserEQUALS)
		}

	}
	p.SetState(588)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(585)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(584)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(587)
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
		p.SetState(590)
		p.Match(EarthParserARG)
	}
	{
		p.SetState(591)
		p.OptionalFlag()
	}
	{
		p.SetState(592)
		p.EnvArgKey()
	}
	p.SetState(600)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserEQUALS {
		{
			p.SetState(593)
			p.Match(EarthParserEQUALS)
		}
		p.SetState(598)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS || _la == EarthParserAtom {
			p.SetState(595)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			if _la == EarthParserWS {
				{
					p.SetState(594)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(597)
				p.EnvArgValue()
			}

		}

	}



	return localctx
}


// ISetStmtContext is an interface to support dynamic dispatch.
type ISetStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SET() antlr.TerminalNode
	EnvArgKey() IEnvArgKeyContext
	EQUALS() antlr.TerminalNode
	EnvArgValue() IEnvArgValueContext
	WS() antlr.TerminalNode

	// IsSetStmtContext differentiates from other interfaces.
	IsSetStmtContext()
}

type SetStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetStmtContext() *SetStmtContext {
	var p = new(SetStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_setStmt
	return p
}

func (*SetStmtContext) IsSetStmtContext() {}

func NewSetStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetStmtContext {
	var p = new(SetStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_setStmt

	return p
}

func (s *SetStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *SetStmtContext) SET() antlr.TerminalNode {
	return s.GetToken(EarthParserSET, 0)
}

func (s *SetStmtContext) EnvArgKey() IEnvArgKeyContext {
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

func (s *SetStmtContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(EarthParserEQUALS, 0)
}

func (s *SetStmtContext) EnvArgValue() IEnvArgValueContext {
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

func (s *SetStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *SetStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SetStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSetStmt(s)
	}
}

func (s *SetStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSetStmt(s)
	}
}




func (p *EarthParser) SetStmt() (localctx ISetStmtContext) {
	this := p
	_ = this

	localctx = NewSetStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, EarthParserRULE_setStmt)
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
		p.SetState(602)
		p.Match(EarthParserSET)
	}
	{
		p.SetState(603)
		p.EnvArgKey()
	}
	{
		p.SetState(604)
		p.Match(EarthParserEQUALS)
	}
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
		p.EnvArgValue()
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
	p.EnterRule(localctx, 116, EarthParserRULE_optionalFlag)

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
	p.SetState(611)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 78, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(610)
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
	p.EnterRule(localctx, 118, EarthParserRULE_envArgKey)

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
		p.SetState(613)
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
	p.EnterRule(localctx, 120, EarthParserRULE_envArgValue)
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
		p.SetState(615)
		p.Match(EarthParserAtom)
	}
	p.SetState(622)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(617)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(616)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(619)
			p.Match(EarthParserAtom)
		}


		p.SetState(624)
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
	p.EnterRule(localctx, 122, EarthParserRULE_labelStmt)
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
		p.SetState(625)
		p.Match(EarthParserLABEL)
	}
	p.SetState(632)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserAtom {
		{
			p.SetState(626)
			p.LabelKey()
		}
		{
			p.SetState(627)
			p.Match(EarthParserEQUALS)
		}
		{
			p.SetState(628)
			p.LabelValue()
		}


		p.SetState(634)
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
	p.EnterRule(localctx, 124, EarthParserRULE_labelKey)

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
		p.SetState(635)
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
	p.EnterRule(localctx, 126, EarthParserRULE_labelValue)

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
		p.SetState(637)
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
	p.EnterRule(localctx, 128, EarthParserRULE_gitCloneStmt)
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
		p.SetState(639)
		p.Match(EarthParserGIT_CLONE)
	}
	p.SetState(641)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(640)
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
	p.EnterRule(localctx, 130, EarthParserRULE_addStmt)
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
		p.SetState(643)
		p.Match(EarthParserADD)
	}
	p.SetState(645)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(644)
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
	p.EnterRule(localctx, 132, EarthParserRULE_stopsignalStmt)
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
		p.SetState(647)
		p.Match(EarthParserSTOPSIGNAL)
	}
	p.SetState(649)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(648)
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
	p.EnterRule(localctx, 134, EarthParserRULE_onbuildStmt)
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
		p.SetState(651)
		p.Match(EarthParserONBUILD)
	}
	p.SetState(653)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(652)
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
	p.EnterRule(localctx, 136, EarthParserRULE_healthcheckStmt)
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
		p.SetState(655)
		p.Match(EarthParserHEALTHCHECK)
	}
	p.SetState(657)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(656)
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
	p.EnterRule(localctx, 138, EarthParserRULE_shellStmt)
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
		p.SetState(659)
		p.Match(EarthParserSHELL)
	}
	p.SetState(661)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(660)
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
	p.EnterRule(localctx, 140, EarthParserRULE_userCommandStmt)
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
		p.SetState(663)
		p.Match(EarthParserCOMMAND)
	}
	p.SetState(665)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(664)
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
	p.EnterRule(localctx, 142, EarthParserRULE_doStmt)
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
		p.SetState(667)
		p.Match(EarthParserDO)
	}
	p.SetState(669)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(668)
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
	p.EnterRule(localctx, 144, EarthParserRULE_importStmt)
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
		p.SetState(671)
		p.Match(EarthParserIMPORT)
	}
	p.SetState(673)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(672)
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
	p.EnterRule(localctx, 146, EarthParserRULE_cacheStmt)
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
		p.SetState(675)
		p.Match(EarthParserCACHE)
	}
	p.SetState(677)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(676)
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
	p.EnterRule(localctx, 148, EarthParserRULE_hostStmt)
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
		p.SetState(679)
		p.Match(EarthParserHOST)
	}
	p.SetState(681)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(680)
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
	p.EnterRule(localctx, 150, EarthParserRULE_projectStmt)
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
		p.SetState(683)
		p.Match(EarthParserPROJECT)
	}
	p.SetState(685)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(684)
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
	p.EnterRule(localctx, 152, EarthParserRULE_pipelineStmt)
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
		p.SetState(687)
		p.Match(EarthParserPIPELINE)
	}
	p.SetState(689)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(688)
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
	p.EnterRule(localctx, 154, EarthParserRULE_triggerStmt)
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
		p.SetState(691)
		p.Match(EarthParserTRIGGER)
	}
	p.SetState(693)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserAtom {
		{
			p.SetState(692)
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
	p.EnterRule(localctx, 156, EarthParserRULE_expr)

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
		p.SetState(695)
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
	p.EnterRule(localctx, 158, EarthParserRULE_stmtWordsMaybeJSON)

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
		p.SetState(697)
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
	p.EnterRule(localctx, 160, EarthParserRULE_stmtWords)

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
	p.SetState(700)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(699)
					p.StmtWord()
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(702)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 96, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 162, EarthParserRULE_stmtWord)

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
		p.SetState(704)
		p.Match(EarthParserAtom)
	}



	return localctx
}


