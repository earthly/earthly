// Code generated from ast/parser/EarthParser.g4 by ANTLR 4.9.1. DO NOT EDIT.

package parser // EarthParser

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 44, 640,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44,
	4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9, 49, 4,
	50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54, 4, 55,
	9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 4, 59, 9, 59, 4, 60, 9,
	60, 4, 61, 9, 61, 4, 62, 9, 62, 4, 63, 9, 63, 4, 64, 9, 64, 4, 65, 9, 65,
	4, 66, 9, 66, 3, 2, 7, 2, 134, 10, 2, 12, 2, 14, 2, 137, 11, 2, 3, 2, 5,
	2, 140, 10, 2, 3, 2, 3, 2, 3, 2, 5, 2, 145, 10, 2, 3, 2, 7, 2, 148, 10,
	2, 12, 2, 14, 2, 151, 11, 2, 3, 2, 5, 2, 154, 10, 2, 3, 2, 7, 2, 157, 10,
	2, 12, 2, 14, 2, 160, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 7, 3, 166, 10, 3,
	12, 3, 14, 3, 169, 11, 3, 3, 3, 7, 3, 172, 10, 3, 12, 3, 14, 3, 175, 11,
	3, 3, 4, 3, 4, 5, 4, 179, 10, 4, 3, 5, 3, 5, 6, 5, 183, 10, 5, 13, 5, 14,
	5, 184, 3, 5, 5, 5, 188, 10, 5, 3, 5, 3, 5, 3, 5, 6, 5, 193, 10, 5, 13,
	5, 14, 5, 194, 3, 5, 3, 5, 5, 5, 199, 10, 5, 3, 6, 3, 6, 3, 7, 3, 7, 6,
	7, 205, 10, 7, 13, 7, 14, 7, 206, 3, 7, 5, 7, 210, 10, 7, 3, 7, 3, 7, 3,
	7, 6, 7, 215, 10, 7, 13, 7, 14, 7, 216, 3, 7, 3, 7, 5, 7, 221, 10, 7, 3,
	8, 3, 8, 3, 9, 5, 9, 226, 10, 9, 3, 9, 3, 9, 6, 9, 230, 10, 9, 13, 9, 14,
	9, 231, 3, 9, 5, 9, 235, 10, 9, 3, 9, 7, 9, 238, 10, 9, 12, 9, 14, 9, 241,
	11, 9, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 247, 10, 10, 3, 11, 3, 11, 3,
	11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11,
	3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3,
	11, 3, 11, 5, 11, 274, 10, 11, 3, 12, 3, 12, 3, 12, 5, 12, 279, 10, 12,
	3, 12, 6, 12, 282, 10, 12, 13, 12, 14, 12, 283, 3, 13, 3, 13, 6, 13, 288,
	10, 13, 13, 13, 14, 13, 289, 3, 13, 5, 13, 293, 10, 13, 3, 13, 5, 13, 296,
	10, 13, 3, 13, 6, 13, 299, 10, 13, 13, 13, 14, 13, 300, 3, 13, 5, 13, 304,
	10, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 16,
	3, 16, 3, 17, 3, 17, 3, 17, 5, 17, 319, 10, 17, 3, 18, 3, 18, 6, 18, 323,
	10, 18, 13, 18, 14, 18, 324, 3, 18, 5, 18, 328, 10, 18, 3, 18, 7, 18, 331,
	10, 18, 12, 18, 14, 18, 334, 11, 18, 3, 18, 6, 18, 337, 10, 18, 13, 18,
	14, 18, 338, 3, 18, 5, 18, 342, 10, 18, 3, 18, 5, 18, 345, 10, 18, 3, 18,
	6, 18, 348, 10, 18, 13, 18, 14, 18, 349, 3, 18, 5, 18, 353, 10, 18, 3,
	18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 6, 19, 361, 10, 19, 13, 19, 14,
	19, 362, 3, 19, 5, 19, 366, 10, 19, 3, 19, 5, 19, 369, 10, 19, 3, 20, 3,
	20, 3, 21, 3, 21, 3, 21, 3, 21, 6, 21, 377, 10, 21, 13, 21, 14, 21, 378,
	3, 21, 5, 21, 382, 10, 21, 3, 21, 5, 21, 385, 10, 21, 3, 22, 3, 22, 3,
	23, 3, 23, 6, 23, 391, 10, 23, 13, 23, 14, 23, 392, 3, 23, 5, 23, 396,
	10, 23, 3, 23, 5, 23, 399, 10, 23, 3, 24, 3, 24, 3, 25, 3, 25, 3, 26, 3,
	26, 3, 27, 3, 27, 6, 27, 409, 10, 27, 13, 27, 14, 27, 410, 3, 27, 5, 27,
	414, 10, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 3, 28, 6, 28, 422, 10,
	28, 13, 28, 14, 28, 423, 3, 28, 5, 28, 427, 10, 28, 3, 28, 5, 28, 430,
	10, 28, 3, 29, 3, 29, 3, 30, 3, 30, 3, 31, 3, 31, 3, 31, 5, 31, 439, 10,
	31, 3, 32, 3, 32, 3, 32, 5, 32, 444, 10, 32, 3, 33, 3, 33, 3, 33, 5, 33,
	449, 10, 33, 3, 34, 3, 34, 3, 34, 5, 34, 454, 10, 34, 3, 35, 3, 35, 5,
	35, 458, 10, 35, 3, 36, 3, 36, 3, 36, 5, 36, 463, 10, 36, 3, 37, 3, 37,
	3, 37, 5, 37, 468, 10, 37, 3, 38, 3, 38, 3, 38, 5, 38, 473, 10, 38, 3,
	39, 3, 39, 3, 39, 5, 39, 478, 10, 39, 3, 40, 3, 40, 3, 40, 5, 40, 483,
	10, 40, 3, 41, 3, 41, 3, 41, 5, 41, 488, 10, 41, 3, 42, 3, 42, 3, 42, 5,
	42, 493, 10, 42, 3, 43, 3, 43, 3, 43, 5, 43, 498, 10, 43, 3, 44, 3, 44,
	3, 44, 5, 44, 503, 10, 44, 3, 45, 3, 45, 3, 45, 5, 45, 508, 10, 45, 3,
	46, 3, 46, 3, 46, 3, 46, 5, 46, 514, 10, 46, 3, 46, 5, 46, 517, 10, 46,
	3, 46, 5, 46, 520, 10, 46, 3, 46, 5, 46, 523, 10, 46, 3, 47, 3, 47, 3,
	47, 3, 47, 3, 47, 5, 47, 530, 10, 47, 3, 47, 3, 47, 3, 47, 5, 47, 535,
	10, 47, 3, 47, 5, 47, 538, 10, 47, 5, 47, 540, 10, 47, 3, 48, 3, 48, 5,
	48, 544, 10, 48, 3, 49, 3, 49, 3, 50, 3, 50, 5, 50, 550, 10, 50, 3, 50,
	7, 50, 553, 10, 50, 12, 50, 14, 50, 556, 11, 50, 3, 51, 3, 51, 3, 51, 3,
	51, 5, 51, 562, 10, 51, 3, 51, 3, 51, 5, 51, 566, 10, 51, 3, 51, 3, 51,
	7, 51, 570, 10, 51, 12, 51, 14, 51, 573, 11, 51, 3, 52, 3, 52, 3, 53, 3,
	53, 3, 54, 3, 54, 3, 54, 5, 54, 582, 10, 54, 3, 55, 3, 55, 3, 55, 5, 55,
	587, 10, 55, 3, 56, 3, 56, 3, 56, 5, 56, 592, 10, 56, 3, 57, 3, 57, 3,
	57, 5, 57, 597, 10, 57, 3, 58, 3, 58, 3, 58, 5, 58, 602, 10, 58, 3, 59,
	3, 59, 3, 59, 5, 59, 607, 10, 59, 3, 60, 3, 60, 3, 60, 5, 60, 612, 10,
	60, 3, 61, 3, 61, 3, 61, 5, 61, 617, 10, 61, 3, 62, 3, 62, 3, 62, 5, 62,
	622, 10, 62, 3, 63, 3, 63, 3, 64, 3, 64, 3, 65, 3, 65, 5, 65, 630, 10,
	65, 3, 65, 7, 65, 633, 10, 65, 12, 65, 14, 65, 636, 11, 65, 3, 66, 3, 66,
	3, 66, 2, 2, 67, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
	32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66,
	68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102,
	104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 2,
	2, 2, 692, 2, 135, 3, 2, 2, 2, 4, 163, 3, 2, 2, 2, 6, 178, 3, 2, 2, 2,
	8, 180, 3, 2, 2, 2, 10, 200, 3, 2, 2, 2, 12, 202, 3, 2, 2, 2, 14, 222,
	3, 2, 2, 2, 16, 225, 3, 2, 2, 2, 18, 246, 3, 2, 2, 2, 20, 273, 3, 2, 2,
	2, 22, 275, 3, 2, 2, 2, 24, 285, 3, 2, 2, 2, 26, 307, 3, 2, 2, 2, 28, 309,
	3, 2, 2, 2, 30, 313, 3, 2, 2, 2, 32, 315, 3, 2, 2, 2, 34, 320, 3, 2, 2,
	2, 36, 356, 3, 2, 2, 2, 38, 370, 3, 2, 2, 2, 40, 372, 3, 2, 2, 2, 42, 386,
	3, 2, 2, 2, 44, 388, 3, 2, 2, 2, 46, 400, 3, 2, 2, 2, 48, 402, 3, 2, 2,
	2, 50, 404, 3, 2, 2, 2, 52, 406, 3, 2, 2, 2, 54, 417, 3, 2, 2, 2, 56, 431,
	3, 2, 2, 2, 58, 433, 3, 2, 2, 2, 60, 435, 3, 2, 2, 2, 62, 440, 3, 2, 2,
	2, 64, 445, 3, 2, 2, 2, 66, 450, 3, 2, 2, 2, 68, 457, 3, 2, 2, 2, 70, 459,
	3, 2, 2, 2, 72, 464, 3, 2, 2, 2, 74, 469, 3, 2, 2, 2, 76, 474, 3, 2, 2,
	2, 78, 479, 3, 2, 2, 2, 80, 484, 3, 2, 2, 2, 82, 489, 3, 2, 2, 2, 84, 494,
	3, 2, 2, 2, 86, 499, 3, 2, 2, 2, 88, 504, 3, 2, 2, 2, 90, 509, 3, 2, 2,
	2, 92, 524, 3, 2, 2, 2, 94, 543, 3, 2, 2, 2, 96, 545, 3, 2, 2, 2, 98, 547,
	3, 2, 2, 2, 100, 557, 3, 2, 2, 2, 102, 574, 3, 2, 2, 2, 104, 576, 3, 2,
	2, 2, 106, 578, 3, 2, 2, 2, 108, 583, 3, 2, 2, 2, 110, 588, 3, 2, 2, 2,
	112, 593, 3, 2, 2, 2, 114, 598, 3, 2, 2, 2, 116, 603, 3, 2, 2, 2, 118,
	608, 3, 2, 2, 2, 120, 613, 3, 2, 2, 2, 122, 618, 3, 2, 2, 2, 124, 623,
	3, 2, 2, 2, 126, 625, 3, 2, 2, 2, 128, 627, 3, 2, 2, 2, 130, 637, 3, 2,
	2, 2, 132, 134, 7, 38, 2, 2, 133, 132, 3, 2, 2, 2, 134, 137, 3, 2, 2, 2,
	135, 133, 3, 2, 2, 2, 135, 136, 3, 2, 2, 2, 136, 139, 3, 2, 2, 2, 137,
	135, 3, 2, 2, 2, 138, 140, 5, 22, 12, 2, 139, 138, 3, 2, 2, 2, 139, 140,
	3, 2, 2, 2, 140, 144, 3, 2, 2, 2, 141, 142, 5, 16, 9, 2, 142, 143, 7, 38,
	2, 2, 143, 145, 3, 2, 2, 2, 144, 141, 3, 2, 2, 2, 144, 145, 3, 2, 2, 2,
	145, 149, 3, 2, 2, 2, 146, 148, 7, 38, 2, 2, 147, 146, 3, 2, 2, 2, 148,
	151, 3, 2, 2, 2, 149, 147, 3, 2, 2, 2, 149, 150, 3, 2, 2, 2, 150, 153,
	3, 2, 2, 2, 151, 149, 3, 2, 2, 2, 152, 154, 5, 4, 3, 2, 153, 152, 3, 2,
	2, 2, 153, 154, 3, 2, 2, 2, 154, 158, 3, 2, 2, 2, 155, 157, 7, 38, 2, 2,
	156, 155, 3, 2, 2, 2, 157, 160, 3, 2, 2, 2, 158, 156, 3, 2, 2, 2, 158,
	159, 3, 2, 2, 2, 159, 161, 3, 2, 2, 2, 160, 158, 3, 2, 2, 2, 161, 162,
	7, 2, 2, 3, 162, 3, 3, 2, 2, 2, 163, 173, 5, 6, 4, 2, 164, 166, 7, 38,
	2, 2, 165, 164, 3, 2, 2, 2, 166, 169, 3, 2, 2, 2, 167, 165, 3, 2, 2, 2,
	167, 168, 3, 2, 2, 2, 168, 170, 3, 2, 2, 2, 169, 167, 3, 2, 2, 2, 170,
	172, 5, 6, 4, 2, 171, 167, 3, 2, 2, 2, 172, 175, 3, 2, 2, 2, 173, 171,
	3, 2, 2, 2, 173, 174, 3, 2, 2, 2, 174, 5, 3, 2, 2, 2, 175, 173, 3, 2, 2,
	2, 176, 179, 5, 8, 5, 2, 177, 179, 5, 12, 7, 2, 178, 176, 3, 2, 2, 2, 178,
	177, 3, 2, 2, 2, 179, 7, 3, 2, 2, 2, 180, 182, 5, 10, 6, 2, 181, 183, 7,
	38, 2, 2, 182, 181, 3, 2, 2, 2, 183, 184, 3, 2, 2, 2, 184, 182, 3, 2, 2,
	2, 184, 185, 3, 2, 2, 2, 185, 187, 3, 2, 2, 2, 186, 188, 7, 39, 2, 2, 187,
	186, 3, 2, 2, 2, 187, 188, 3, 2, 2, 2, 188, 198, 3, 2, 2, 2, 189, 190,
	7, 3, 2, 2, 190, 192, 5, 16, 9, 2, 191, 193, 7, 38, 2, 2, 192, 191, 3,
	2, 2, 2, 193, 194, 3, 2, 2, 2, 194, 192, 3, 2, 2, 2, 194, 195, 3, 2, 2,
	2, 195, 196, 3, 2, 2, 2, 196, 197, 7, 4, 2, 2, 197, 199, 3, 2, 2, 2, 198,
	189, 3, 2, 2, 2, 198, 199, 3, 2, 2, 2, 199, 9, 3, 2, 2, 2, 200, 201, 7,
	5, 2, 2, 201, 11, 3, 2, 2, 2, 202, 204, 5, 14, 8, 2, 203, 205, 7, 38, 2,
	2, 204, 203, 3, 2, 2, 2, 205, 206, 3, 2, 2, 2, 206, 204, 3, 2, 2, 2, 206,
	207, 3, 2, 2, 2, 207, 209, 3, 2, 2, 2, 208, 210, 7, 39, 2, 2, 209, 208,
	3, 2, 2, 2, 209, 210, 3, 2, 2, 2, 210, 220, 3, 2, 2, 2, 211, 212, 7, 3,
	2, 2, 212, 214, 5, 16, 9, 2, 213, 215, 7, 38, 2, 2, 214, 213, 3, 2, 2,
	2, 215, 216, 3, 2, 2, 2, 216, 214, 3, 2, 2, 2, 216, 217, 3, 2, 2, 2, 217,
	218, 3, 2, 2, 2, 218, 219, 7, 4, 2, 2, 219, 221, 3, 2, 2, 2, 220, 211,
	3, 2, 2, 2, 220, 221, 3, 2, 2, 2, 221, 13, 3, 2, 2, 2, 222, 223, 7, 6,
	2, 2, 223, 15, 3, 2, 2, 2, 224, 226, 7, 39, 2, 2, 225, 224, 3, 2, 2, 2,
	225, 226, 3, 2, 2, 2, 226, 227, 3, 2, 2, 2, 227, 239, 5, 18, 10, 2, 228,
	230, 7, 38, 2, 2, 229, 228, 3, 2, 2, 2, 230, 231, 3, 2, 2, 2, 231, 229,
	3, 2, 2, 2, 231, 232, 3, 2, 2, 2, 232, 234, 3, 2, 2, 2, 233, 235, 7, 39,
	2, 2, 234, 233, 3, 2, 2, 2, 234, 235, 3, 2, 2, 2, 235, 236, 3, 2, 2, 2,
	236, 238, 5, 18, 10, 2, 237, 229, 3, 2, 2, 2, 238, 241, 3, 2, 2, 2, 239,
	237, 3, 2, 2, 2, 239, 240, 3, 2, 2, 2, 240, 17, 3, 2, 2, 2, 241, 239, 3,
	2, 2, 2, 242, 247, 5, 20, 11, 2, 243, 247, 5, 24, 13, 2, 244, 247, 5, 34,
	18, 2, 245, 247, 5, 52, 27, 2, 246, 242, 3, 2, 2, 2, 246, 243, 3, 2, 2,
	2, 246, 244, 3, 2, 2, 2, 246, 245, 3, 2, 2, 2, 247, 19, 3, 2, 2, 2, 248,
	274, 5, 60, 31, 2, 249, 274, 5, 62, 32, 2, 250, 274, 5, 64, 33, 2, 251,
	274, 5, 66, 34, 2, 252, 274, 5, 68, 35, 2, 253, 274, 5, 74, 38, 2, 254,
	274, 5, 76, 39, 2, 255, 274, 5, 78, 40, 2, 256, 274, 5, 80, 41, 2, 257,
	274, 5, 82, 42, 2, 258, 274, 5, 84, 43, 2, 259, 274, 5, 86, 44, 2, 260,
	274, 5, 88, 45, 2, 261, 274, 5, 90, 46, 2, 262, 274, 5, 92, 47, 2, 263,
	274, 5, 100, 51, 2, 264, 274, 5, 106, 54, 2, 265, 274, 5, 108, 55, 2, 266,
	274, 5, 110, 56, 2, 267, 274, 5, 112, 57, 2, 268, 274, 5, 114, 58, 2, 269,
	274, 5, 116, 59, 2, 270, 274, 5, 118, 60, 2, 271, 274, 5, 120, 61, 2, 272,
	274, 5, 122, 62, 2, 273, 248, 3, 2, 2, 2, 273, 249, 3, 2, 2, 2, 273, 250,
	3, 2, 2, 2, 273, 251, 3, 2, 2, 2, 273, 252, 3, 2, 2, 2, 273, 253, 3, 2,
	2, 2, 273, 254, 3, 2, 2, 2, 273, 255, 3, 2, 2, 2, 273, 256, 3, 2, 2, 2,
	273, 257, 3, 2, 2, 2, 273, 258, 3, 2, 2, 2, 273, 259, 3, 2, 2, 2, 273,
	260, 3, 2, 2, 2, 273, 261, 3, 2, 2, 2, 273, 262, 3, 2, 2, 2, 273, 263,
	3, 2, 2, 2, 273, 264, 3, 2, 2, 2, 273, 265, 3, 2, 2, 2, 273, 266, 3, 2,
	2, 2, 273, 267, 3, 2, 2, 2, 273, 268, 3, 2, 2, 2, 273, 269, 3, 2, 2, 2,
	273, 270, 3, 2, 2, 2, 273, 271, 3, 2, 2, 2, 273, 272, 3, 2, 2, 2, 274,
	21, 3, 2, 2, 2, 275, 278, 7, 33, 2, 2, 276, 277, 7, 39, 2, 2, 277, 279,
	5, 128, 65, 2, 278, 276, 3, 2, 2, 2, 278, 279, 3, 2, 2, 2, 279, 281, 3,
	2, 2, 2, 280, 282, 7, 38, 2, 2, 281, 280, 3, 2, 2, 2, 282, 283, 3, 2, 2,
	2, 283, 281, 3, 2, 2, 2, 283, 284, 3, 2, 2, 2, 284, 23, 3, 2, 2, 2, 285,
	295, 5, 28, 15, 2, 286, 288, 7, 38, 2, 2, 287, 286, 3, 2, 2, 2, 288, 289,
	3, 2, 2, 2, 289, 287, 3, 2, 2, 2, 289, 290, 3, 2, 2, 2, 290, 292, 3, 2,
	2, 2, 291, 293, 7, 39, 2, 2, 292, 291, 3, 2, 2, 2, 292, 293, 3, 2, 2, 2,
	293, 294, 3, 2, 2, 2, 294, 296, 5, 26, 14, 2, 295, 287, 3, 2, 2, 2, 295,
	296, 3, 2, 2, 2, 296, 298, 3, 2, 2, 2, 297, 299, 7, 38, 2, 2, 298, 297,
	3, 2, 2, 2, 299, 300, 3, 2, 2, 2, 300, 298, 3, 2, 2, 2, 300, 301, 3, 2,
	2, 2, 301, 303, 3, 2, 2, 2, 302, 304, 7, 39, 2, 2, 303, 302, 3, 2, 2, 2,
	303, 304, 3, 2, 2, 2, 304, 305, 3, 2, 2, 2, 305, 306, 7, 42, 2, 2, 306,
	25, 3, 2, 2, 2, 307, 308, 5, 16, 9, 2, 308, 27, 3, 2, 2, 2, 309, 310, 7,
	34, 2, 2, 310, 311, 7, 39, 2, 2, 311, 312, 5, 30, 16, 2, 312, 29, 3, 2,
	2, 2, 313, 314, 5, 32, 17, 2, 314, 31, 3, 2, 2, 2, 315, 318, 7, 35, 2,
	2, 316, 317, 7, 39, 2, 2, 317, 319, 5, 128, 65, 2, 318, 316, 3, 2, 2, 2,
	318, 319, 3, 2, 2, 2, 319, 33, 3, 2, 2, 2, 320, 332, 5, 36, 19, 2, 321,
	323, 7, 38, 2, 2, 322, 321, 3, 2, 2, 2, 323, 324, 3, 2, 2, 2, 324, 322,
	3, 2, 2, 2, 324, 325, 3, 2, 2, 2, 325, 327, 3, 2, 2, 2, 326, 328, 7, 39,
	2, 2, 327, 326, 3, 2, 2, 2, 327, 328, 3, 2, 2, 2, 328, 329, 3, 2, 2, 2,
	329, 331, 5, 40, 21, 2, 330, 322, 3, 2, 2, 2, 331, 334, 3, 2, 2, 2, 332,
	330, 3, 2, 2, 2, 332, 333, 3, 2, 2, 2, 333, 344, 3, 2, 2, 2, 334, 332,
	3, 2, 2, 2, 335, 337, 7, 38, 2, 2, 336, 335, 3, 2, 2, 2, 337, 338, 3, 2,
	2, 2, 338, 336, 3, 2, 2, 2, 338, 339, 3, 2, 2, 2, 339, 341, 3, 2, 2, 2,
	340, 342, 7, 39, 2, 2, 341, 340, 3, 2, 2, 2, 341, 342, 3, 2, 2, 2, 342,
	343, 3, 2, 2, 2, 343, 345, 5, 44, 23, 2, 344, 336, 3, 2, 2, 2, 344, 345,
	3, 2, 2, 2, 345, 347, 3, 2, 2, 2, 346, 348, 7, 38, 2, 2, 347, 346, 3, 2,
	2, 2, 348, 349, 3, 2, 2, 2, 349, 347, 3, 2, 2, 2, 349, 350, 3, 2, 2, 2,
	350, 352, 3, 2, 2, 2, 351, 353, 7, 39, 2, 2, 352, 351, 3, 2, 2, 2, 352,
	353, 3, 2, 2, 2, 353, 354, 3, 2, 2, 2, 354, 355, 7, 42, 2, 2, 355, 35,
	3, 2, 2, 2, 356, 357, 7, 36, 2, 2, 357, 358, 7, 39, 2, 2, 358, 368, 5,
	48, 25, 2, 359, 361, 7, 38, 2, 2, 360, 359, 3, 2, 2, 2, 361, 362, 3, 2,
	2, 2, 362, 360, 3, 2, 2, 2, 362, 363, 3, 2, 2, 2, 363, 365, 3, 2, 2, 2,
	364, 366, 7, 39, 2, 2, 365, 364, 3, 2, 2, 2, 365, 366, 3, 2, 2, 2, 366,
	367, 3, 2, 2, 2, 367, 369, 5, 38, 20, 2, 368, 360, 3, 2, 2, 2, 368, 369,
	3, 2, 2, 2, 369, 37, 3, 2, 2, 2, 370, 371, 5, 16, 9, 2, 371, 39, 3, 2,
	2, 2, 372, 373, 7, 41, 2, 2, 373, 374, 7, 39, 2, 2, 374, 384, 5, 50, 26,
	2, 375, 377, 7, 38, 2, 2, 376, 375, 3, 2, 2, 2, 377, 378, 3, 2, 2, 2, 378,
	376, 3, 2, 2, 2, 378, 379, 3, 2, 2, 2, 379, 381, 3, 2, 2, 2, 380, 382,
	7, 39, 2, 2, 381, 380, 3, 2, 2, 2, 381, 382, 3, 2, 2, 2, 382, 383, 3, 2,
	2, 2, 383, 385, 5, 42, 22, 2, 384, 376, 3, 2, 2, 2, 384, 385, 3, 2, 2,
	2, 385, 41, 3, 2, 2, 2, 386, 387, 5, 16, 9, 2, 387, 43, 3, 2, 2, 2, 388,
	398, 7, 40, 2, 2, 389, 391, 7, 38, 2, 2, 390, 389, 3, 2, 2, 2, 391, 392,
	3, 2, 2, 2, 392, 390, 3, 2, 2, 2, 392, 393, 3, 2, 2, 2, 393, 395, 3, 2,
	2, 2, 394, 396, 7, 39, 2, 2, 395, 394, 3, 2, 2, 2, 395, 396, 3, 2, 2, 2,
	396, 397, 3, 2, 2, 2, 397, 399, 5, 46, 24, 2, 398, 390, 3, 2, 2, 2, 398,
	399, 3, 2, 2, 2, 399, 45, 3, 2, 2, 2, 400, 401, 5, 16, 9, 2, 401, 47, 3,
	2, 2, 2, 402, 403, 5, 124, 63, 2, 403, 49, 3, 2, 2, 2, 404, 405, 5, 124,
	63, 2, 405, 51, 3, 2, 2, 2, 406, 408, 5, 54, 28, 2, 407, 409, 7, 38, 2,
	2, 408, 407, 3, 2, 2, 2, 409, 410, 3, 2, 2, 2, 410, 408, 3, 2, 2, 2, 410,
	411, 3, 2, 2, 2, 411, 413, 3, 2, 2, 2, 412, 414, 7, 39, 2, 2, 413, 412,
	3, 2, 2, 2, 413, 414, 3, 2, 2, 2, 414, 415, 3, 2, 2, 2, 415, 416, 7, 42,
	2, 2, 416, 53, 3, 2, 2, 2, 417, 418, 7, 37, 2, 2, 418, 419, 7, 39, 2, 2,
	419, 429, 5, 58, 30, 2, 420, 422, 7, 38, 2, 2, 421, 420, 3, 2, 2, 2, 422,
	423, 3, 2, 2, 2, 423, 421, 3, 2, 2, 2, 423, 424, 3, 2, 2, 2, 424, 426,
	3, 2, 2, 2, 425, 427, 7, 39, 2, 2, 426, 425, 3, 2, 2, 2, 426, 427, 3, 2,
	2, 2, 427, 428, 3, 2, 2, 2, 428, 430, 5, 56, 29, 2, 429, 421, 3, 2, 2,
	2, 429, 430, 3, 2, 2, 2, 430, 55, 3, 2, 2, 2, 431, 432, 5, 16, 9, 2, 432,
	57, 3, 2, 2, 2, 433, 434, 5, 128, 65, 2, 434, 59, 3, 2, 2, 2, 435, 438,
	7, 7, 2, 2, 436, 437, 7, 39, 2, 2, 437, 439, 5, 128, 65, 2, 438, 436, 3,
	2, 2, 2, 438, 439, 3, 2, 2, 2, 439, 61, 3, 2, 2, 2, 440, 443, 7, 8, 2,
	2, 441, 442, 7, 39, 2, 2, 442, 444, 5, 128, 65, 2, 443, 441, 3, 2, 2, 2,
	443, 444, 3, 2, 2, 2, 444, 63, 3, 2, 2, 2, 445, 448, 7, 9, 2, 2, 446, 447,
	7, 39, 2, 2, 447, 449, 5, 128, 65, 2, 448, 446, 3, 2, 2, 2, 448, 449, 3,
	2, 2, 2, 449, 65, 3, 2, 2, 2, 450, 453, 7, 10, 2, 2, 451, 452, 7, 39, 2,
	2, 452, 454, 5, 128, 65, 2, 453, 451, 3, 2, 2, 2, 453, 454, 3, 2, 2, 2,
	454, 67, 3, 2, 2, 2, 455, 458, 5, 72, 37, 2, 456, 458, 5, 70, 36, 2, 457,
	455, 3, 2, 2, 2, 457, 456, 3, 2, 2, 2, 458, 69, 3, 2, 2, 2, 459, 462, 7,
	12, 2, 2, 460, 461, 7, 39, 2, 2, 461, 463, 5, 128, 65, 2, 462, 460, 3,
	2, 2, 2, 462, 463, 3, 2, 2, 2, 463, 71, 3, 2, 2, 2, 464, 467, 7, 11, 2,
	2, 465, 466, 7, 39, 2, 2, 466, 468, 5, 128, 65, 2, 467, 465, 3, 2, 2, 2,
	467, 468, 3, 2, 2, 2, 468, 73, 3, 2, 2, 2, 469, 472, 7, 13, 2, 2, 470,
	471, 7, 39, 2, 2, 471, 473, 5, 126, 64, 2, 472, 470, 3, 2, 2, 2, 472, 473,
	3, 2, 2, 2, 473, 75, 3, 2, 2, 2, 474, 477, 7, 19, 2, 2, 475, 476, 7, 39,
	2, 2, 476, 478, 5, 128, 65, 2, 477, 475, 3, 2, 2, 2, 477, 478, 3, 2, 2,
	2, 478, 77, 3, 2, 2, 2, 479, 482, 7, 20, 2, 2, 480, 481, 7, 39, 2, 2, 481,
	483, 5, 128, 65, 2, 482, 480, 3, 2, 2, 2, 482, 483, 3, 2, 2, 2, 483, 79,
	3, 2, 2, 2, 484, 487, 7, 21, 2, 2, 485, 486, 7, 39, 2, 2, 486, 488, 5,
	128, 65, 2, 487, 485, 3, 2, 2, 2, 487, 488, 3, 2, 2, 2, 488, 81, 3, 2,
	2, 2, 489, 492, 7, 22, 2, 2, 490, 491, 7, 39, 2, 2, 491, 493, 5, 126, 64,
	2, 492, 490, 3, 2, 2, 2, 492, 493, 3, 2, 2, 2, 493, 83, 3, 2, 2, 2, 494,
	497, 7, 23, 2, 2, 495, 496, 7, 39, 2, 2, 496, 498, 5, 126, 64, 2, 497,
	495, 3, 2, 2, 2, 497, 498, 3, 2, 2, 2, 498, 85, 3, 2, 2, 2, 499, 502, 7,
	14, 2, 2, 500, 501, 7, 39, 2, 2, 501, 503, 5, 128, 65, 2, 502, 500, 3,
	2, 2, 2, 502, 503, 3, 2, 2, 2, 503, 87, 3, 2, 2, 2, 504, 507, 7, 15, 2,
	2, 505, 506, 7, 39, 2, 2, 506, 508, 5, 126, 64, 2, 507, 505, 3, 2, 2, 2,
	507, 508, 3, 2, 2, 2, 508, 89, 3, 2, 2, 2, 509, 510, 7, 16, 2, 2, 510,
	511, 7, 39, 2, 2, 511, 516, 5, 96, 49, 2, 512, 514, 7, 39, 2, 2, 513, 512,
	3, 2, 2, 2, 513, 514, 3, 2, 2, 2, 514, 515, 3, 2, 2, 2, 515, 517, 7, 44,
	2, 2, 516, 513, 3, 2, 2, 2, 516, 517, 3, 2, 2, 2, 517, 522, 3, 2, 2, 2,
	518, 520, 7, 39, 2, 2, 519, 518, 3, 2, 2, 2, 519, 520, 3, 2, 2, 2, 520,
	521, 3, 2, 2, 2, 521, 523, 5, 98, 50, 2, 522, 519, 3, 2, 2, 2, 522, 523,
	3, 2, 2, 2, 523, 91, 3, 2, 2, 2, 524, 525, 7, 17, 2, 2, 525, 526, 5, 94,
	48, 2, 526, 527, 7, 39, 2, 2, 527, 539, 5, 96, 49, 2, 528, 530, 7, 39,
	2, 2, 529, 528, 3, 2, 2, 2, 529, 530, 3, 2, 2, 2, 530, 531, 3, 2, 2, 2,
	531, 532, 7, 44, 2, 2, 532, 537, 3, 2, 2, 2, 533, 535, 7, 39, 2, 2, 534,
	533, 3, 2, 2, 2, 534, 535, 3, 2, 2, 2, 535, 536, 3, 2, 2, 2, 536, 538,
	5, 98, 50, 2, 537, 534, 3, 2, 2, 2, 537, 538, 3, 2, 2, 2, 538, 540, 3,
	2, 2, 2, 539, 529, 3, 2, 2, 2, 539, 540, 3, 2, 2, 2, 540, 93, 3, 2, 2,
	2, 541, 542, 7, 39, 2, 2, 542, 544, 5, 128, 65, 2, 543, 541, 3, 2, 2, 2,
	543, 544, 3, 2, 2, 2, 544, 95, 3, 2, 2, 2, 545, 546, 7, 43, 2, 2, 546,
	97, 3, 2, 2, 2, 547, 554, 7, 43, 2, 2, 548, 550, 7, 39, 2, 2, 549, 548,
	3, 2, 2, 2, 549, 550, 3, 2, 2, 2, 550, 551, 3, 2, 2, 2, 551, 553, 7, 43,
	2, 2, 552, 549, 3, 2, 2, 2, 553, 556, 3, 2, 2, 2, 554, 552, 3, 2, 2, 2,
	554, 555, 3, 2, 2, 2, 555, 99, 3, 2, 2, 2, 556, 554, 3, 2, 2, 2, 557, 571,
	7, 18, 2, 2, 558, 559, 7, 39, 2, 2, 559, 561, 5, 102, 52, 2, 560, 562,
	7, 39, 2, 2, 561, 560, 3, 2, 2, 2, 561, 562, 3, 2, 2, 2, 562, 563, 3, 2,
	2, 2, 563, 565, 7, 44, 2, 2, 564, 566, 7, 39, 2, 2, 565, 564, 3, 2, 2,
	2, 565, 566, 3, 2, 2, 2, 566, 567, 3, 2, 2, 2, 567, 568, 5, 104, 53, 2,
	568, 570, 3, 2, 2, 2, 569, 558, 3, 2, 2, 2, 570, 573, 3, 2, 2, 2, 571,
	569, 3, 2, 2, 2, 571, 572, 3, 2, 2, 2, 572, 101, 3, 2, 2, 2, 573, 571,
	3, 2, 2, 2, 574, 575, 7, 43, 2, 2, 575, 103, 3, 2, 2, 2, 576, 577, 7, 43,
	2, 2, 577, 105, 3, 2, 2, 2, 578, 581, 7, 24, 2, 2, 579, 580, 7, 39, 2,
	2, 580, 582, 5, 128, 65, 2, 581, 579, 3, 2, 2, 2, 581, 582, 3, 2, 2, 2,
	582, 107, 3, 2, 2, 2, 583, 586, 7, 25, 2, 2, 584, 585, 7, 39, 2, 2, 585,
	587, 5, 128, 65, 2, 586, 584, 3, 2, 2, 2, 586, 587, 3, 2, 2, 2, 587, 109,
	3, 2, 2, 2, 588, 591, 7, 26, 2, 2, 589, 590, 7, 39, 2, 2, 590, 592, 5,
	128, 65, 2, 591, 589, 3, 2, 2, 2, 591, 592, 3, 2, 2, 2, 592, 111, 3, 2,
	2, 2, 593, 596, 7, 27, 2, 2, 594, 595, 7, 39, 2, 2, 595, 597, 5, 128, 65,
	2, 596, 594, 3, 2, 2, 2, 596, 597, 3, 2, 2, 2, 597, 113, 3, 2, 2, 2, 598,
	601, 7, 28, 2, 2, 599, 600, 7, 39, 2, 2, 600, 602, 5, 128, 65, 2, 601,
	599, 3, 2, 2, 2, 601, 602, 3, 2, 2, 2, 602, 115, 3, 2, 2, 2, 603, 606,
	7, 29, 2, 2, 604, 605, 7, 39, 2, 2, 605, 607, 5, 128, 65, 2, 606, 604,
	3, 2, 2, 2, 606, 607, 3, 2, 2, 2, 607, 117, 3, 2, 2, 2, 608, 611, 7, 31,
	2, 2, 609, 610, 7, 39, 2, 2, 610, 612, 5, 128, 65, 2, 611, 609, 3, 2, 2,
	2, 611, 612, 3, 2, 2, 2, 612, 119, 3, 2, 2, 2, 613, 616, 7, 30, 2, 2, 614,
	615, 7, 39, 2, 2, 615, 617, 5, 128, 65, 2, 616, 614, 3, 2, 2, 2, 616, 617,
	3, 2, 2, 2, 617, 121, 3, 2, 2, 2, 618, 621, 7, 32, 2, 2, 619, 620, 7, 39,
	2, 2, 620, 622, 5, 128, 65, 2, 621, 619, 3, 2, 2, 2, 621, 622, 3, 2, 2,
	2, 622, 123, 3, 2, 2, 2, 623, 624, 5, 126, 64, 2, 624, 125, 3, 2, 2, 2,
	625, 626, 5, 128, 65, 2, 626, 127, 3, 2, 2, 2, 627, 634, 5, 130, 66, 2,
	628, 630, 7, 39, 2, 2, 629, 628, 3, 2, 2, 2, 629, 630, 3, 2, 2, 2, 630,
	631, 3, 2, 2, 2, 631, 633, 5, 130, 66, 2, 632, 629, 3, 2, 2, 2, 633, 636,
	3, 2, 2, 2, 634, 632, 3, 2, 2, 2, 634, 635, 3, 2, 2, 2, 635, 129, 3, 2,
	2, 2, 636, 634, 3, 2, 2, 2, 637, 638, 7, 43, 2, 2, 638, 131, 3, 2, 2, 2,
	95, 135, 139, 144, 149, 153, 158, 167, 173, 178, 184, 187, 194, 198, 206,
	209, 216, 220, 225, 231, 234, 239, 246, 273, 278, 283, 289, 292, 295, 300,
	303, 318, 324, 327, 332, 338, 341, 344, 349, 352, 362, 365, 368, 378, 381,
	384, 392, 395, 398, 410, 413, 423, 426, 429, 438, 443, 448, 453, 457, 462,
	467, 472, 477, 482, 487, 492, 497, 502, 507, 513, 516, 519, 522, 529, 534,
	537, 539, 543, 549, 554, 561, 565, 571, 581, 586, 591, 596, 601, 606, 611,
	616, 621, 629, 634,
}
var literalNames = []string{
	"", "", "", "", "", "'FROM'", "'FROM DOCKERFILE'", "'LOCALLY'", "'COPY'",
	"'SAVE ARTIFACT'", "'SAVE IMAGE'", "'RUN'", "'EXPOSE'", "'VOLUME'", "'ENV'",
	"'ARG'", "'LABEL'", "'BUILD'", "'WORKDIR'", "'USER'", "'CMD'", "'ENTRYPOINT'",
	"'GIT CLONE'", "'ADD'", "'STOPSIGNAL'", "'ONBUILD'", "'HEALTHCHECK'", "'SHELL'",
	"'DO'", "'COMMAND'", "'IMPORT'", "'VERSION'", "'WITH'", "", "", "", "",
	"", "'ELSE'", "'ELSE IF'", "'END'",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "Target", "UserCommand", "FROM", "FROM_DOCKERFILE",
	"LOCALLY", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE", "RUN", "EXPOSE", "VOLUME",
	"ENV", "ARG", "LABEL", "BUILD", "WORKDIR", "USER", "CMD", "ENTRYPOINT",
	"GIT_CLONE", "ADD", "STOPSIGNAL", "ONBUILD", "HEALTHCHECK", "SHELL", "DO",
	"COMMAND", "IMPORT", "VERSION", "WITH", "DOCKER", "IF", "FOR", "NL", "WS",
	"ELSE", "ELSE_IF", "END", "Atom", "EQUALS",
}

var ruleNames = []string{
	"earthFile", "targets", "targetOrUserCommand", "target", "targetHeader",
	"userCommand", "userCommandHeader", "stmts", "stmt", "commandStmt", "version",
	"withStmt", "withBlock", "withExpr", "withCommand", "dockerCommand", "ifStmt",
	"ifClause", "ifBlock", "elseIfClause", "elseIfBlock", "elseClause", "elseBlock",
	"ifExpr", "elseIfExpr", "forStmt", "forClause", "forBlock", "forExpr",
	"fromStmt", "fromDockerfileStmt", "locallyStmt", "copyStmt", "saveStmt",
	"saveImage", "saveArtifact", "runStmt", "buildStmt", "workdirStmt", "userStmt",
	"cmdStmt", "entrypointStmt", "exposeStmt", "volumeStmt", "envStmt", "argStmt",
	"optionalFlag", "envArgKey", "envArgValue", "labelStmt", "labelKey", "labelValue",
	"gitCloneStmt", "addStmt", "stopsignalStmt", "onbuildStmt", "healthcheckStmt",
	"shellStmt", "userCommandStmt", "doStmt", "importStmt", "expr", "stmtWordsMaybeJSON",
	"stmtWords", "stmtWord",
}

type EarthParser struct {
	*antlr.BaseParser
}

// NewEarthParser produces a new parser instance for the optional input antlr.TokenStream.
//
// The *EarthParser instance produced may be reused by calling the SetInputStream method.
// The initial parser configuration is expensive to construct, and the object is not thread-safe;
// however, if used within a Golang sync.Pool, the construction cost amortizes well and the
// objects can be used in a thread-safe manner.
func NewEarthParser(input antlr.TokenStream) *EarthParser {
	this := new(EarthParser)
	deserializer := antlr.NewATNDeserializer(nil)
	deserializedATN := deserializer.DeserializeFromUInt16(parserATN)
	decisionToDFA := make([]*antlr.DFA, len(deserializedATN.DecisionToState))
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "EarthParser.g4"

	return this
}

// EarthParser tokens.
const (
	EarthParserEOF             = antlr.TokenEOF
	EarthParserINDENT          = 1
	EarthParserDEDENT          = 2
	EarthParserTarget          = 3
	EarthParserUserCommand     = 4
	EarthParserFROM            = 5
	EarthParserFROM_DOCKERFILE = 6
	EarthParserLOCALLY         = 7
	EarthParserCOPY            = 8
	EarthParserSAVE_ARTIFACT   = 9
	EarthParserSAVE_IMAGE      = 10
	EarthParserRUN             = 11
	EarthParserEXPOSE          = 12
	EarthParserVOLUME          = 13
	EarthParserENV             = 14
	EarthParserARG             = 15
	EarthParserLABEL           = 16
	EarthParserBUILD           = 17
	EarthParserWORKDIR         = 18
	EarthParserUSER            = 19
	EarthParserCMD             = 20
	EarthParserENTRYPOINT      = 21
	EarthParserGIT_CLONE       = 22
	EarthParserADD             = 23
	EarthParserSTOPSIGNAL      = 24
	EarthParserONBUILD         = 25
	EarthParserHEALTHCHECK     = 26
	EarthParserSHELL           = 27
	EarthParserDO              = 28
	EarthParserCOMMAND         = 29
	EarthParserIMPORT          = 30
	EarthParserVERSION         = 31
	EarthParserWITH            = 32
	EarthParserDOCKER          = 33
	EarthParserIF              = 34
	EarthParserFOR             = 35
	EarthParserNL              = 36
	EarthParserWS              = 37
	EarthParserELSE            = 38
	EarthParserELSE_IF         = 39
	EarthParserEND             = 40
	EarthParserAtom            = 41
	EarthParserEQUALS          = 42
)

// EarthParser rules.
const (
	EarthParserRULE_earthFile           = 0
	EarthParserRULE_targets             = 1
	EarthParserRULE_targetOrUserCommand = 2
	EarthParserRULE_target              = 3
	EarthParserRULE_targetHeader        = 4
	EarthParserRULE_userCommand         = 5
	EarthParserRULE_userCommandHeader   = 6
	EarthParserRULE_stmts               = 7
	EarthParserRULE_stmt                = 8
	EarthParserRULE_commandStmt         = 9
	EarthParserRULE_version             = 10
	EarthParserRULE_withStmt            = 11
	EarthParserRULE_withBlock           = 12
	EarthParserRULE_withExpr            = 13
	EarthParserRULE_withCommand         = 14
	EarthParserRULE_dockerCommand       = 15
	EarthParserRULE_ifStmt              = 16
	EarthParserRULE_ifClause            = 17
	EarthParserRULE_ifBlock             = 18
	EarthParserRULE_elseIfClause        = 19
	EarthParserRULE_elseIfBlock         = 20
	EarthParserRULE_elseClause          = 21
	EarthParserRULE_elseBlock           = 22
	EarthParserRULE_ifExpr              = 23
	EarthParserRULE_elseIfExpr          = 24
	EarthParserRULE_forStmt             = 25
	EarthParserRULE_forClause           = 26
	EarthParserRULE_forBlock            = 27
	EarthParserRULE_forExpr             = 28
	EarthParserRULE_fromStmt            = 29
	EarthParserRULE_fromDockerfileStmt  = 30
	EarthParserRULE_locallyStmt         = 31
	EarthParserRULE_copyStmt            = 32
	EarthParserRULE_saveStmt            = 33
	EarthParserRULE_saveImage           = 34
	EarthParserRULE_saveArtifact        = 35
	EarthParserRULE_runStmt             = 36
	EarthParserRULE_buildStmt           = 37
	EarthParserRULE_workdirStmt         = 38
	EarthParserRULE_userStmt            = 39
	EarthParserRULE_cmdStmt             = 40
	EarthParserRULE_entrypointStmt      = 41
	EarthParserRULE_exposeStmt          = 42
	EarthParserRULE_volumeStmt          = 43
	EarthParserRULE_envStmt             = 44
	EarthParserRULE_argStmt             = 45
	EarthParserRULE_optionalFlag        = 46
	EarthParserRULE_envArgKey           = 47
	EarthParserRULE_envArgValue         = 48
	EarthParserRULE_labelStmt           = 49
	EarthParserRULE_labelKey            = 50
	EarthParserRULE_labelValue          = 51
	EarthParserRULE_gitCloneStmt        = 52
	EarthParserRULE_addStmt             = 53
	EarthParserRULE_stopsignalStmt      = 54
	EarthParserRULE_onbuildStmt         = 55
	EarthParserRULE_healthcheckStmt     = 56
	EarthParserRULE_shellStmt           = 57
	EarthParserRULE_userCommandStmt     = 58
	EarthParserRULE_doStmt              = 59
	EarthParserRULE_importStmt          = 60
	EarthParserRULE_expr                = 61
	EarthParserRULE_stmtWordsMaybeJSON  = 62
	EarthParserRULE_stmtWords           = 63
	EarthParserRULE_stmtWord            = 64
)

// IEarthFileContext is an interface to support dynamic dispatch.
type IEarthFileContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVersionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVersionContext)
}

func (s *EarthFileContext) Stmts() IStmtsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *EarthFileContext) Targets() ITargetsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITargetsContext)(nil)).Elem(), 0)

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
	p.SetState(133)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(130)
				p.Match(EarthParserNL)
			}

		}
		p.SetState(135)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserVERSION {
		{
			p.SetState(136)
			p.Version()
		}

	}
	p.SetState(142)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if (((_la)&-(0x1f+1)) == 0 && ((int64(1)<<uint(_la))&((int64(1)<<EarthParserFROM)|(int64(1)<<EarthParserFROM_DOCKERFILE)|(int64(1)<<EarthParserLOCALLY)|(int64(1)<<EarthParserCOPY)|(int64(1)<<EarthParserSAVE_ARTIFACT)|(int64(1)<<EarthParserSAVE_IMAGE)|(int64(1)<<EarthParserRUN)|(int64(1)<<EarthParserEXPOSE)|(int64(1)<<EarthParserVOLUME)|(int64(1)<<EarthParserENV)|(int64(1)<<EarthParserARG)|(int64(1)<<EarthParserLABEL)|(int64(1)<<EarthParserBUILD)|(int64(1)<<EarthParserWORKDIR)|(int64(1)<<EarthParserUSER)|(int64(1)<<EarthParserCMD)|(int64(1)<<EarthParserENTRYPOINT)|(int64(1)<<EarthParserGIT_CLONE)|(int64(1)<<EarthParserADD)|(int64(1)<<EarthParserSTOPSIGNAL)|(int64(1)<<EarthParserONBUILD)|(int64(1)<<EarthParserHEALTHCHECK)|(int64(1)<<EarthParserSHELL)|(int64(1)<<EarthParserDO)|(int64(1)<<EarthParserCOMMAND)|(int64(1)<<EarthParserIMPORT))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((int64(1)<<uint((_la-32)))&((int64(1)<<(EarthParserWITH-32))|(int64(1)<<(EarthParserIF-32))|(int64(1)<<(EarthParserFOR-32))|(int64(1)<<(EarthParserWS-32)))) != 0) {
		{
			p.SetState(139)
			p.Stmts()
		}
		{
			p.SetState(140)
			p.Match(EarthParserNL)
		}

	}
	p.SetState(147)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(144)
				p.Match(EarthParserNL)
			}

		}
		p.SetState(149)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}
	p.SetState(151)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserTarget || _la == EarthParserUserCommand {
		{
			p.SetState(150)
			p.Targets()
		}

	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == EarthParserNL {
		{
			p.SetState(153)
			p.Match(EarthParserNL)
		}

		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(159)
		p.Match(EarthParserEOF)
	}

	return localctx
}

// ITargetsContext is an interface to support dynamic dispatch.
type ITargetsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITargetOrUserCommandContext)(nil)).Elem())
	var tst = make([]ITargetOrUserCommandContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITargetOrUserCommandContext)
		}
	}

	return tst
}

func (s *TargetsContext) TargetOrUserCommand(i int) ITargetOrUserCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITargetOrUserCommandContext)(nil)).Elem(), i)

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
		p.SetState(161)
		p.TargetOrUserCommand()
	}
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(165)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for _la == EarthParserNL {
				{
					p.SetState(162)
					p.Match(EarthParserNL)
				}

				p.SetState(167)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(168)
				p.TargetOrUserCommand()
			}

		}
		p.SetState(173)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITargetContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITargetContext)
}

func (s *TargetOrUserCommandContext) UserCommand() IUserCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUserCommandContext)(nil)).Elem(), 0)

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

	p.SetState(176)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserTarget:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(174)
			p.Target()
		}

	case EarthParserUserCommand:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(175)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITargetHeaderContext)(nil)).Elem(), 0)

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

func (s *TargetContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *TargetContext) INDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserINDENT, 0)
}

func (s *TargetContext) Stmts() IStmtsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtsContext)
}

func (s *TargetContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserDEDENT, 0)
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
		p.SetState(178)
		p.TargetHeader()
	}
	p.SetState(180)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(179)
				p.Match(EarthParserNL)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(182)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}
	p.SetState(185)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(184)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(196)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserINDENT {
		{
			p.SetState(187)
			p.Match(EarthParserINDENT)
		}
		{
			p.SetState(188)
			p.Stmts()
		}
		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(189)
				p.Match(EarthParserNL)
			}

			p.SetState(192)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(194)
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
		p.SetState(198)
		p.Match(EarthParserTarget)
	}

	return localctx
}

// IUserCommandContext is an interface to support dynamic dispatch.
type IUserCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUserCommandHeaderContext)(nil)).Elem(), 0)

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

func (s *UserCommandContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *UserCommandContext) INDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserINDENT, 0)
}

func (s *UserCommandContext) Stmts() IStmtsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
		p.SetState(200)
		p.UserCommandHeader()
	}
	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(201)
				p.Match(EarthParserNL)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(204)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())
	}
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(206)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(218)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserINDENT {
		{
			p.SetState(209)
			p.Match(EarthParserINDENT)
		}
		{
			p.SetState(210)
			p.Stmts()
		}
		p.SetState(212)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(211)
				p.Match(EarthParserNL)
			}

			p.SetState(214)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(216)
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
		p.SetState(220)
		p.Match(EarthParserUserCommand)
	}

	return localctx
}

// IStmtsContext is an interface to support dynamic dispatch.
type IStmtsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStmtContext)(nil)).Elem())
	var tst = make([]IStmtContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStmtContext)
		}
	}

	return tst
}

func (s *StmtsContext) Stmt(i int) IStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStmtContext)
}

func (s *StmtsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *StmtsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
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
	p.SetState(223)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(222)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(225)
		p.Stmt()
	}
	p.SetState(237)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
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
			p.SetState(232)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(231)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(234)
				p.Stmt()
			}

		}
		p.SetState(239)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext())
	}

	return localctx
}

// IStmtContext is an interface to support dynamic dispatch.
type IStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommandStmtContext)
}

func (s *StmtContext) WithStmt() IWithStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithStmtContext)
}

func (s *StmtContext) IfStmt() IIfStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIfStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIfStmtContext)
}

func (s *StmtContext) ForStmt() IForStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForStmtContext)
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

	p.SetState(244)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM, EarthParserFROM_DOCKERFILE, EarthParserLOCALLY, EarthParserCOPY, EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE, EarthParserRUN, EarthParserEXPOSE, EarthParserVOLUME, EarthParserENV, EarthParserARG, EarthParserLABEL, EarthParserBUILD, EarthParserWORKDIR, EarthParserUSER, EarthParserCMD, EarthParserENTRYPOINT, EarthParserGIT_CLONE, EarthParserADD, EarthParserSTOPSIGNAL, EarthParserONBUILD, EarthParserHEALTHCHECK, EarthParserSHELL, EarthParserDO, EarthParserCOMMAND, EarthParserIMPORT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(240)
			p.CommandStmt()
		}

	case EarthParserWITH:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(241)
			p.WithStmt()
		}

	case EarthParserIF:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(242)
			p.IfStmt()
		}

	case EarthParserFOR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(243)
			p.ForStmt()
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromStmtContext)
}

func (s *CommandStmtContext) FromDockerfileStmt() IFromDockerfileStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromDockerfileStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromDockerfileStmtContext)
}

func (s *CommandStmtContext) LocallyStmt() ILocallyStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILocallyStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILocallyStmtContext)
}

func (s *CommandStmtContext) CopyStmt() ICopyStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyStmtContext)
}

func (s *CommandStmtContext) SaveStmt() ISaveStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveStmtContext)
}

func (s *CommandStmtContext) RunStmt() IRunStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRunStmtContext)
}

func (s *CommandStmtContext) BuildStmt() IBuildStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBuildStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBuildStmtContext)
}

func (s *CommandStmtContext) WorkdirStmt() IWorkdirStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWorkdirStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWorkdirStmtContext)
}

func (s *CommandStmtContext) UserStmt() IUserStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUserStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUserStmtContext)
}

func (s *CommandStmtContext) CmdStmt() ICmdStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICmdStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICmdStmtContext)
}

func (s *CommandStmtContext) EntrypointStmt() IEntrypointStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntrypointStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEntrypointStmtContext)
}

func (s *CommandStmtContext) ExposeStmt() IExposeStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExposeStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExposeStmtContext)
}

func (s *CommandStmtContext) VolumeStmt() IVolumeStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVolumeStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVolumeStmtContext)
}

func (s *CommandStmtContext) EnvStmt() IEnvStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEnvStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEnvStmtContext)
}

func (s *CommandStmtContext) ArgStmt() IArgStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgStmtContext)
}

func (s *CommandStmtContext) LabelStmt() ILabelStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILabelStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILabelStmtContext)
}

func (s *CommandStmtContext) GitCloneStmt() IGitCloneStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGitCloneStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGitCloneStmtContext)
}

func (s *CommandStmtContext) AddStmt() IAddStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAddStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAddStmtContext)
}

func (s *CommandStmtContext) StopsignalStmt() IStopsignalStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStopsignalStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStopsignalStmtContext)
}

func (s *CommandStmtContext) OnbuildStmt() IOnbuildStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOnbuildStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOnbuildStmtContext)
}

func (s *CommandStmtContext) HealthcheckStmt() IHealthcheckStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHealthcheckStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHealthcheckStmtContext)
}

func (s *CommandStmtContext) ShellStmt() IShellStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShellStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShellStmtContext)
}

func (s *CommandStmtContext) UserCommandStmt() IUserCommandStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUserCommandStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUserCommandStmtContext)
}

func (s *CommandStmtContext) DoStmt() IDoStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDoStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDoStmtContext)
}

func (s *CommandStmtContext) ImportStmt() IImportStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImportStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportStmtContext)
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

	p.SetState(271)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(246)
			p.FromStmt()
		}

	case EarthParserFROM_DOCKERFILE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(247)
			p.FromDockerfileStmt()
		}

	case EarthParserLOCALLY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(248)
			p.LocallyStmt()
		}

	case EarthParserCOPY:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(249)
			p.CopyStmt()
		}

	case EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(250)
			p.SaveStmt()
		}

	case EarthParserRUN:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(251)
			p.RunStmt()
		}

	case EarthParserBUILD:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(252)
			p.BuildStmt()
		}

	case EarthParserWORKDIR:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(253)
			p.WorkdirStmt()
		}

	case EarthParserUSER:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(254)
			p.UserStmt()
		}

	case EarthParserCMD:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(255)
			p.CmdStmt()
		}

	case EarthParserENTRYPOINT:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(256)
			p.EntrypointStmt()
		}

	case EarthParserEXPOSE:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(257)
			p.ExposeStmt()
		}

	case EarthParserVOLUME:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(258)
			p.VolumeStmt()
		}

	case EarthParserENV:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(259)
			p.EnvStmt()
		}

	case EarthParserARG:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(260)
			p.ArgStmt()
		}

	case EarthParserLABEL:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(261)
			p.LabelStmt()
		}

	case EarthParserGIT_CLONE:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(262)
			p.GitCloneStmt()
		}

	case EarthParserADD:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(263)
			p.AddStmt()
		}

	case EarthParserSTOPSIGNAL:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(264)
			p.StopsignalStmt()
		}

	case EarthParserONBUILD:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(265)
			p.OnbuildStmt()
		}

	case EarthParserHEALTHCHECK:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(266)
			p.HealthcheckStmt()
		}

	case EarthParserSHELL:
		p.EnterOuterAlt(localctx, 22)
		{
			p.SetState(267)
			p.ShellStmt()
		}

	case EarthParserCOMMAND:
		p.EnterOuterAlt(localctx, 23)
		{
			p.SetState(268)
			p.UserCommandStmt()
		}

	case EarthParserDO:
		p.EnterOuterAlt(localctx, 24)
		{
			p.SetState(269)
			p.DoStmt()
		}

	case EarthParserIMPORT:
		p.EnterOuterAlt(localctx, 25)
		{
			p.SetState(270)
			p.ImportStmt()
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

func (s *VersionContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *VersionContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewVersionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, EarthParserRULE_version)
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
		p.SetState(273)
		p.Match(EarthParserVERSION)
	}
	p.SetState(276)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(274)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(275)
			p.StmtWords()
		}

	}
	p.SetState(279)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(278)
				p.Match(EarthParserNL)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(281)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext())
	}

	return localctx
}

// IWithStmtContext is an interface to support dynamic dispatch.
type IWithStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWithExprContext)
}

func (s *WithStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *WithStmtContext) WithBlock() IWithBlockContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithBlockContext)(nil)).Elem(), 0)

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

func (s *WithStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *WithStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
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
		p.SetState(283)
		p.WithExpr()
	}
	p.SetState(293)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		p.SetState(285)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(284)
				p.Match(EarthParserNL)
			}

			p.SetState(287)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(290)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(289)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(292)
			p.WithBlock()
		}

	}
	p.SetState(296)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(295)
			p.Match(EarthParserNL)
		}

		p.SetState(298)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(301)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(300)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(303)
		p.Match(EarthParserEND)
	}

	return localctx
}

// IWithBlockContext is an interface to support dynamic dispatch.
type IWithBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
		p.SetState(305)
		p.Stmts()
	}

	return localctx
}

// IWithExprContext is an interface to support dynamic dispatch.
type IWithExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *WithExprContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *WithExprContext) WithCommand() IWithCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWithCommandContext)(nil)).Elem(), 0)

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
		p.SetState(307)
		p.Match(EarthParserWITH)
	}
	{
		p.SetState(308)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(309)
		p.WithCommand()
	}

	return localctx
}

// IWithCommandContext is an interface to support dynamic dispatch.
type IWithCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDockerCommandContext)(nil)).Elem(), 0)

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
		p.SetState(311)
		p.DockerCommand()
	}

	return localctx
}

// IDockerCommandContext is an interface to support dynamic dispatch.
type IDockerCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *DockerCommandContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *DockerCommandContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
		p.SetState(313)
		p.Match(EarthParserDOCKER)
	}
	p.SetState(316)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(314)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(315)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIfClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIfClauseContext)
}

func (s *IfStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *IfStmtContext) AllElseIfClause() []IElseIfClauseContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IElseIfClauseContext)(nil)).Elem())
	var tst = make([]IElseIfClauseContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IElseIfClauseContext)
		}
	}

	return tst
}

func (s *IfStmtContext) ElseIfClause(i int) IElseIfClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IElseIfClauseContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IElseIfClauseContext)
}

func (s *IfStmtContext) ElseClause() IElseClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IElseClauseContext)(nil)).Elem(), 0)

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

func (s *IfStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *IfStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
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
		p.SetState(318)
		p.IfClause()
	}
	p.SetState(330)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(320)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(319)
					p.Match(EarthParserNL)
				}

				p.SetState(322)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			p.SetState(325)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(324)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(327)
				p.ElseIfClause()
			}

		}
		p.SetState(332)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext())
	}
	p.SetState(342)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext()) == 1 {
		p.SetState(334)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(333)
				p.Match(EarthParserNL)
			}

			p.SetState(336)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(339)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(338)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(341)
			p.ElseClause()
		}

	}
	p.SetState(345)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(344)
			p.Match(EarthParserNL)
		}

		p.SetState(347)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(349)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(352)
		p.Match(EarthParserEND)
	}

	return localctx
}

// IIfClauseContext is an interface to support dynamic dispatch.
type IIfClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *IfClauseContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *IfClauseContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *IfClauseContext) IfExpr() IIfExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIfExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIfExprContext)
}

func (s *IfClauseContext) IfBlock() IIfBlockContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIfBlockContext)(nil)).Elem(), 0)

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
		p.SetState(354)
		p.Match(EarthParserIF)
	}
	{
		p.SetState(355)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(356)
		p.IfExpr()
	}
	p.SetState(366)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) == 1 {
		p.SetState(358)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(357)
				p.Match(EarthParserNL)
			}

			p.SetState(360)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(363)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(362)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(365)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
		p.SetState(368)
		p.Stmts()
	}

	return localctx
}

// IElseIfClauseContext is an interface to support dynamic dispatch.
type IElseIfClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *ElseIfClauseContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *ElseIfClauseContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *ElseIfClauseContext) ElseIfExpr() IElseIfExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IElseIfExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IElseIfExprContext)
}

func (s *ElseIfClauseContext) ElseIfBlock() IElseIfBlockContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IElseIfBlockContext)(nil)).Elem(), 0)

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
		p.SetState(370)
		p.Match(EarthParserELSE_IF)
	}
	{
		p.SetState(371)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(372)
		p.ElseIfExpr()
	}
	p.SetState(382)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
		p.SetState(374)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(373)
				p.Match(EarthParserNL)
			}

			p.SetState(376)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(379)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(378)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(381)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
		p.SetState(384)
		p.Stmts()
	}

	return localctx
}

// IElseClauseContext is an interface to support dynamic dispatch.
type IElseClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IElseBlockContext)(nil)).Elem(), 0)

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

func (s *ElseClauseContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
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
		p.SetState(386)
		p.Match(EarthParserELSE)
	}
	p.SetState(396)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext()) == 1 {
		p.SetState(388)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(387)
				p.Match(EarthParserNL)
			}

			p.SetState(390)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(393)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(392)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(395)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
		p.SetState(398)
		p.Stmts()
	}

	return localctx
}

// IIfExprContext is an interface to support dynamic dispatch.
type IIfExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

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
		p.SetState(400)
		p.Expr()
	}

	return localctx
}

// IElseIfExprContext is an interface to support dynamic dispatch.
type IElseIfExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

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
		p.SetState(402)
		p.Expr()
	}

	return localctx
}

// IForStmtContext is an interface to support dynamic dispatch.
type IForStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForClauseContext)(nil)).Elem(), 0)

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

func (s *ForStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
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
	localctx = NewForStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, EarthParserRULE_forStmt)
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
		p.SetState(404)
		p.ForClause()
	}
	p.SetState(406)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(405)
			p.Match(EarthParserNL)
		}

		p.SetState(408)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(411)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(410)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(413)
		p.Match(EarthParserEND)
	}

	return localctx
}

// IForClauseContext is an interface to support dynamic dispatch.
type IForClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *ForClauseContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *ForClauseContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *ForClauseContext) ForExpr() IForExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IForExprContext)
}

func (s *ForClauseContext) ForBlock() IForBlockContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IForBlockContext)(nil)).Elem(), 0)

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
	localctx = NewForClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, EarthParserRULE_forClause)
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
		p.SetState(415)
		p.Match(EarthParserFOR)
	}
	{
		p.SetState(416)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(417)
		p.ForExpr()
	}
	p.SetState(427)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext()) == 1 {
		p.SetState(419)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(418)
				p.Match(EarthParserNL)
			}

			p.SetState(421)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(424)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(423)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(426)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
	localctx = NewForBlockContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, EarthParserRULE_forBlock)

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
		p.SetState(429)
		p.Stmts()
	}

	return localctx
}

// IForExprContext is an interface to support dynamic dispatch.
type IForExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewForExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, EarthParserRULE_forExpr)

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
		p.SetState(431)
		p.StmtWords()
	}

	return localctx
}

// IFromStmtContext is an interface to support dynamic dispatch.
type IFromStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *FromStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *FromStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewFromStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, EarthParserRULE_fromStmt)
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
		p.SetState(433)
		p.Match(EarthParserFROM)
	}
	p.SetState(436)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(434)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(435)
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

func (s *FromDockerfileStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *FromDockerfileStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewFromDockerfileStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, EarthParserRULE_fromDockerfileStmt)
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
		p.SetState(438)
		p.Match(EarthParserFROM_DOCKERFILE)
	}
	p.SetState(441)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(439)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(440)
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

func (s *LocallyStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *LocallyStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewLocallyStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, EarthParserRULE_locallyStmt)
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
		p.SetState(443)
		p.Match(EarthParserLOCALLY)
	}
	p.SetState(446)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(444)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(445)
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

func (s *CopyStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *CopyStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewCopyStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, EarthParserRULE_copyStmt)
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
		p.Match(EarthParserCOPY)
	}
	p.SetState(451)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(449)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(450)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveArtifactContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveArtifactContext)
}

func (s *SaveStmtContext) SaveImage() ISaveImageContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveImageContext)(nil)).Elem(), 0)

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
	localctx = NewSaveStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, EarthParserRULE_saveStmt)

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

	p.SetState(455)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserSAVE_ARTIFACT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(453)
			p.SaveArtifact()
		}

	case EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(454)
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

func (s *SaveImageContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *SaveImageContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewSaveImageContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, EarthParserRULE_saveImage)
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
		p.SetState(457)
		p.Match(EarthParserSAVE_IMAGE)
	}
	p.SetState(460)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(458)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(459)
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

func (s *SaveArtifactContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *SaveArtifactContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewSaveArtifactContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, EarthParserRULE_saveArtifact)
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
		p.Match(EarthParserSAVE_ARTIFACT)
	}
	p.SetState(465)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(463)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(464)
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

func (s *RunStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *RunStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsMaybeJSONContext)(nil)).Elem(), 0)

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
	localctx = NewRunStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, EarthParserRULE_runStmt)
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
		p.SetState(467)
		p.Match(EarthParserRUN)
	}
	p.SetState(470)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(468)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(469)
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

func (s *BuildStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *BuildStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewBuildStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, EarthParserRULE_buildStmt)
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
		p.SetState(472)
		p.Match(EarthParserBUILD)
	}
	p.SetState(475)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(473)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(474)
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

func (s *WorkdirStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *WorkdirStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewWorkdirStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, EarthParserRULE_workdirStmt)
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
		p.SetState(477)
		p.Match(EarthParserWORKDIR)
	}
	p.SetState(480)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(478)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(479)
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

func (s *UserStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *UserStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewUserStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, EarthParserRULE_userStmt)
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
		p.SetState(482)
		p.Match(EarthParserUSER)
	}
	p.SetState(485)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(483)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(484)
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

func (s *CmdStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *CmdStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsMaybeJSONContext)(nil)).Elem(), 0)

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
	localctx = NewCmdStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, EarthParserRULE_cmdStmt)
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
		p.SetState(487)
		p.Match(EarthParserCMD)
	}
	p.SetState(490)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(488)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(489)
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

func (s *EntrypointStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *EntrypointStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsMaybeJSONContext)(nil)).Elem(), 0)

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
	localctx = NewEntrypointStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, EarthParserRULE_entrypointStmt)
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
		p.Match(EarthParserENTRYPOINT)
	}
	p.SetState(495)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(493)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(494)
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

func (s *ExposeStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *ExposeStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewExposeStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, EarthParserRULE_exposeStmt)
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
		p.SetState(497)
		p.Match(EarthParserEXPOSE)
	}
	p.SetState(500)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(498)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(499)
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

func (s *VolumeStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *VolumeStmtContext) StmtWordsMaybeJSON() IStmtWordsMaybeJSONContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsMaybeJSONContext)(nil)).Elem(), 0)

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
	localctx = NewVolumeStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, EarthParserRULE_volumeStmt)
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
		p.SetState(502)
		p.Match(EarthParserVOLUME)
	}
	p.SetState(505)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(503)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(504)
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

func (s *EnvStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *EnvStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *EnvStmtContext) EnvArgKey() IEnvArgKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEnvArgKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEnvArgKeyContext)
}

func (s *EnvStmtContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(EarthParserEQUALS, 0)
}

func (s *EnvStmtContext) EnvArgValue() IEnvArgValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEnvArgValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEnvArgValueContext)
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
	localctx = NewEnvStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, EarthParserRULE_envStmt)
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
		p.SetState(507)
		p.Match(EarthParserENV)
	}
	{
		p.SetState(508)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(509)
		p.EnvArgKey()
	}
	p.SetState(514)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 69, p.GetParserRuleContext()) == 1 {
		p.SetState(511)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(510)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(513)
			p.Match(EarthParserEQUALS)
		}

	}
	p.SetState(520)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(517)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(516)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(519)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOptionalFlagContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOptionalFlagContext)
}

func (s *ArgStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *ArgStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *ArgStmtContext) EnvArgKey() IEnvArgKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEnvArgKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEnvArgKeyContext)
}

func (s *ArgStmtContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(EarthParserEQUALS, 0)
}

func (s *ArgStmtContext) EnvArgValue() IEnvArgValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEnvArgValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEnvArgValueContext)
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
	localctx = NewArgStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, EarthParserRULE_argStmt)
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
		p.SetState(522)
		p.Match(EarthParserARG)
	}
	{
		p.SetState(523)
		p.OptionalFlag()
	}
	{
		p.SetState(524)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(525)
		p.EnvArgKey()
	}
	p.SetState(537)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS || _la == EarthParserEQUALS {
		p.SetState(527)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(526)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(529)
			p.Match(EarthParserEQUALS)
		}

		p.SetState(535)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS || _la == EarthParserAtom {
			p.SetState(532)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(531)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(534)
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

func (s *OptionalFlagContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *OptionalFlagContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewOptionalFlagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, EarthParserRULE_optionalFlag)

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
	p.SetState(541)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 76, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(539)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(540)
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
	localctx = NewEnvArgKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, EarthParserRULE_envArgKey)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IEnvArgValueContext is an interface to support dynamic dispatch.
type IEnvArgValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	localctx = NewEnvArgValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, EarthParserRULE_envArgValue)
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
		p.SetState(545)
		p.Match(EarthParserAtom)
	}
	p.SetState(552)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(547)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(546)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(549)
			p.Match(EarthParserAtom)
		}

		p.SetState(554)
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

func (s *LabelStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *LabelStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *LabelStmtContext) AllLabelKey() []ILabelKeyContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILabelKeyContext)(nil)).Elem())
	var tst = make([]ILabelKeyContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILabelKeyContext)
		}
	}

	return tst
}

func (s *LabelStmtContext) LabelKey(i int) ILabelKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILabelKeyContext)(nil)).Elem(), i)

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILabelValueContext)(nil)).Elem())
	var tst = make([]ILabelValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILabelValueContext)
		}
	}

	return tst
}

func (s *LabelStmtContext) LabelValue(i int) ILabelValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILabelValueContext)(nil)).Elem(), i)

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
	localctx = NewLabelStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, EarthParserRULE_labelStmt)
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
		p.Match(EarthParserLABEL)
	}
	p.SetState(569)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == EarthParserWS {
		{
			p.SetState(556)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(557)
			p.LabelKey()
		}
		p.SetState(559)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(558)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(561)
			p.Match(EarthParserEQUALS)
		}
		p.SetState(563)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(562)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(565)
			p.LabelValue()
		}

		p.SetState(571)
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
	localctx = NewLabelKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, EarthParserRULE_labelKey)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ILabelValueContext is an interface to support dynamic dispatch.
type ILabelValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	localctx = NewLabelValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, EarthParserRULE_labelValue)

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
		p.SetState(574)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IGitCloneStmtContext is an interface to support dynamic dispatch.
type IGitCloneStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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

func (s *GitCloneStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *GitCloneStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewGitCloneStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, EarthParserRULE_gitCloneStmt)
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
		p.Match(EarthParserGIT_CLONE)
	}
	p.SetState(579)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(577)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(578)
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

func (s *AddStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *AddStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewAddStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, EarthParserRULE_addStmt)
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
		p.SetState(581)
		p.Match(EarthParserADD)
	}
	p.SetState(584)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(582)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(583)
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

func (s *StopsignalStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *StopsignalStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewStopsignalStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, EarthParserRULE_stopsignalStmt)
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
		p.SetState(586)
		p.Match(EarthParserSTOPSIGNAL)
	}
	p.SetState(589)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(587)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(588)
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

func (s *OnbuildStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *OnbuildStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewOnbuildStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, EarthParserRULE_onbuildStmt)
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
		p.SetState(591)
		p.Match(EarthParserONBUILD)
	}
	p.SetState(594)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(592)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(593)
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

func (s *HealthcheckStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *HealthcheckStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewHealthcheckStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, EarthParserRULE_healthcheckStmt)
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
		p.SetState(596)
		p.Match(EarthParserHEALTHCHECK)
	}
	p.SetState(599)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(597)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(598)
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

func (s *ShellStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *ShellStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewShellStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, EarthParserRULE_shellStmt)
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
		p.SetState(601)
		p.Match(EarthParserSHELL)
	}
	p.SetState(604)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(602)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(603)
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

func (s *UserCommandStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *UserCommandStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewUserCommandStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, EarthParserRULE_userCommandStmt)
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
		p.SetState(606)
		p.Match(EarthParserCOMMAND)
	}
	p.SetState(609)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(607)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(608)
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

func (s *DoStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *DoStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewDoStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, EarthParserRULE_doStmt)
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
		p.SetState(611)
		p.Match(EarthParserDO)
	}
	p.SetState(614)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(612)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(613)
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

func (s *ImportStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *ImportStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewImportStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, EarthParserRULE_importStmt)
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
		p.SetState(616)
		p.Match(EarthParserIMPORT)
	}
	p.SetState(619)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(617)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(618)
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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsMaybeJSONContext)(nil)).Elem(), 0)

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
	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, EarthParserRULE_expr)

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
		p.SetState(621)
		p.StmtWordsMaybeJSON()
	}

	return localctx
}

// IStmtWordsMaybeJSONContext is an interface to support dynamic dispatch.
type IStmtWordsMaybeJSONContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewStmtWordsMaybeJSONContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, EarthParserRULE_stmtWordsMaybeJSON)

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
		p.SetState(623)
		p.StmtWords()
	}

	return localctx
}

// IStmtWordsContext is an interface to support dynamic dispatch.
type IStmtWordsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStmtWordContext)(nil)).Elem())
	var tst = make([]IStmtWordContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStmtWordContext)
		}
	}

	return tst
}

func (s *StmtWordsContext) StmtWord(i int) IStmtWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStmtWordContext)
}

func (s *StmtWordsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *StmtWordsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
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
	localctx = NewStmtWordsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, EarthParserRULE_stmtWords)
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
		p.SetState(625)
		p.StmtWord()
	}
	p.SetState(632)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 92, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(627)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(626)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(629)
				p.StmtWord()
			}

		}
		p.SetState(634)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 92, p.GetParserRuleContext())
	}

	return localctx
}

// IStmtWordContext is an interface to support dynamic dispatch.
type IStmtWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

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
	localctx = NewStmtWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, EarthParserRULE_stmtWord)

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
