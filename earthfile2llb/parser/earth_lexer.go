// Code generated from earthfile2llb/parser/EarthLexer.g4 by ANTLR 4.8. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 29, 533,
	8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4,
	4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10,
	4, 11, 9, 11, 4, 12, 9, 12, 4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4,
	16, 9, 16, 4, 17, 9, 17, 4, 18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21,
	9, 21, 4, 22, 9, 22, 4, 23, 9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9,
	26, 4, 27, 9, 27, 4, 28, 9, 28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31,
	4, 32, 9, 32, 4, 33, 9, 33, 4, 34, 9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4,
	37, 9, 37, 4, 38, 9, 38, 4, 39, 9, 39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42,
	9, 42, 4, 43, 9, 43, 4, 44, 9, 44, 4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9,
	47, 4, 48, 9, 48, 4, 49, 9, 49, 4, 50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52,
	4, 53, 9, 53, 4, 54, 9, 54, 4, 55, 9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4,
	58, 9, 58, 4, 59, 9, 59, 4, 60, 9, 60, 4, 61, 9, 61, 4, 62, 9, 62, 3, 2,
	6, 2, 132, 10, 2, 13, 2, 14, 2, 133, 3, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5,
	3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6,
	3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 8,
	3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10,
	3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 3,
	11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12,
	3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3,
	13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13,
	3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3,
	14, 3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15,
	3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 16, 6, 16, 273, 10,
	16, 13, 16, 14, 16, 274, 3, 16, 3, 16, 3, 17, 5, 17, 280, 10, 17, 3, 17,
	5, 17, 283, 10, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 7, 18, 290, 10,
	18, 12, 18, 14, 18, 293, 11, 18, 3, 18, 6, 18, 296, 10, 18, 13, 18, 14,
	18, 297, 3, 19, 3, 19, 3, 19, 5, 19, 303, 10, 19, 3, 20, 3, 20, 7, 20,
	307, 10, 20, 12, 20, 14, 20, 310, 11, 20, 3, 21, 3, 21, 3, 21, 3, 21, 3,
	21, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3,
	26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28,
	3, 28, 3, 28, 3, 28, 3, 28, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 30, 3,
	30, 3, 30, 3, 30, 3, 30, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 32, 3, 32,
	3, 32, 3, 32, 3, 32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 34, 3, 34, 3,
	34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 36, 3, 36, 3, 36,
	3, 36, 3, 37, 3, 37, 3, 37, 3, 37, 3, 38, 3, 38, 3, 38, 3, 38, 3, 39, 3,
	39, 3, 39, 3, 39, 3, 40, 3, 40, 3, 40, 3, 40, 5, 40, 407, 10, 40, 3, 41,
	3, 41, 5, 41, 411, 10, 41, 3, 41, 3, 41, 7, 41, 415, 10, 41, 12, 41, 14,
	41, 418, 11, 41, 3, 41, 3, 41, 3, 42, 3, 42, 3, 42, 3, 42, 7, 42, 426,
	10, 42, 12, 42, 14, 42, 429, 11, 42, 3, 42, 3, 42, 3, 43, 3, 43, 3, 44,
	3, 44, 3, 45, 5, 45, 438, 10, 45, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 3,
	46, 3, 46, 3, 46, 3, 46, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 48, 3, 48,
	3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 49, 3, 49, 3, 49, 3,
	49, 3, 50, 5, 50, 468, 10, 50, 3, 50, 3, 50, 3, 50, 3, 50, 3, 50, 3, 51,
	3, 51, 3, 51, 3, 51, 3, 52, 3, 52, 3, 52, 3, 52, 3, 53, 3, 53, 3, 54, 3,
	54, 3, 54, 3, 54, 3, 55, 5, 55, 490, 10, 55, 3, 55, 3, 55, 3, 55, 3, 55,
	3, 55, 3, 56, 3, 56, 3, 56, 3, 56, 3, 57, 3, 57, 3, 57, 3, 57, 3, 58, 3,
	58, 5, 58, 507, 10, 58, 3, 58, 3, 58, 7, 58, 511, 10, 58, 12, 58, 14, 58,
	514, 11, 58, 3, 58, 3, 58, 3, 59, 3, 59, 3, 60, 3, 60, 3, 61, 5, 61, 523,
	10, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 61, 3, 62, 3, 62, 3, 62, 3, 62,
	2, 2, 63, 8, 5, 10, 6, 12, 7, 14, 8, 16, 9, 18, 10, 20, 11, 22, 12, 24,
	13, 26, 14, 28, 15, 30, 16, 32, 17, 34, 18, 36, 19, 38, 20, 40, 21, 42,
	2, 44, 2, 46, 2, 48, 2, 50, 2, 52, 2, 54, 2, 56, 2, 58, 2, 60, 2, 62, 2,
	64, 2, 66, 2, 68, 2, 70, 2, 72, 2, 74, 2, 76, 2, 78, 2, 80, 22, 82, 23,
	84, 24, 86, 25, 88, 2, 90, 2, 92, 2, 94, 2, 96, 2, 98, 2, 100, 26, 102,
	2, 104, 2, 106, 2, 108, 27, 110, 28, 112, 2, 114, 2, 116, 2, 118, 29, 120,
	2, 122, 2, 124, 2, 126, 2, 128, 2, 8, 2, 3, 4, 5, 6, 7, 11, 6, 2, 47, 48,
	50, 59, 67, 92, 99, 124, 3, 2, 67, 92, 4, 2, 11, 11, 34, 34, 4, 2, 12,
	12, 15, 15, 3, 2, 36, 36, 7, 2, 11, 12, 15, 15, 34, 34, 36, 36, 93, 93,
	6, 2, 11, 12, 15, 15, 34, 34, 36, 36, 8, 2, 11, 12, 15, 15, 34, 34, 36,
	36, 63, 63, 93, 93, 7, 2, 11, 12, 15, 15, 34, 34, 36, 36, 63, 63, 2, 542,
	2, 8, 3, 2, 2, 2, 2, 10, 3, 2, 2, 2, 2, 12, 3, 2, 2, 2, 2, 14, 3, 2, 2,
	2, 2, 16, 3, 2, 2, 2, 2, 18, 3, 2, 2, 2, 2, 20, 3, 2, 2, 2, 2, 22, 3, 2,
	2, 2, 2, 24, 3, 2, 2, 2, 2, 26, 3, 2, 2, 2, 2, 28, 3, 2, 2, 2, 2, 30, 3,
	2, 2, 2, 2, 32, 3, 2, 2, 2, 2, 34, 3, 2, 2, 2, 2, 36, 3, 2, 2, 2, 2, 38,
	3, 2, 2, 2, 2, 40, 3, 2, 2, 2, 3, 46, 3, 2, 2, 2, 3, 48, 3, 2, 2, 2, 3,
	50, 3, 2, 2, 2, 3, 52, 3, 2, 2, 2, 3, 54, 3, 2, 2, 2, 3, 56, 3, 2, 2, 2,
	3, 58, 3, 2, 2, 2, 3, 60, 3, 2, 2, 2, 3, 62, 3, 2, 2, 2, 3, 64, 3, 2, 2,
	2, 3, 66, 3, 2, 2, 2, 3, 68, 3, 2, 2, 2, 3, 70, 3, 2, 2, 2, 3, 72, 3, 2,
	2, 2, 3, 74, 3, 2, 2, 2, 3, 76, 3, 2, 2, 2, 3, 78, 3, 2, 2, 2, 4, 80, 3,
	2, 2, 2, 4, 82, 3, 2, 2, 2, 4, 84, 3, 2, 2, 2, 4, 86, 3, 2, 2, 2, 4, 94,
	3, 2, 2, 2, 4, 96, 3, 2, 2, 2, 5, 98, 3, 2, 2, 2, 5, 100, 3, 2, 2, 2, 5,
	102, 3, 2, 2, 2, 5, 104, 3, 2, 2, 2, 5, 106, 3, 2, 2, 2, 6, 108, 3, 2,
	2, 2, 6, 110, 3, 2, 2, 2, 6, 112, 3, 2, 2, 2, 6, 114, 3, 2, 2, 2, 6, 116,
	3, 2, 2, 2, 7, 118, 3, 2, 2, 2, 7, 120, 3, 2, 2, 2, 7, 126, 3, 2, 2, 2,
	7, 128, 3, 2, 2, 2, 8, 131, 3, 2, 2, 2, 10, 139, 3, 2, 2, 2, 12, 146, 3,
	2, 2, 2, 14, 153, 3, 2, 2, 2, 16, 169, 3, 2, 2, 2, 18, 182, 3, 2, 2, 2,
	20, 188, 3, 2, 2, 2, 22, 194, 3, 2, 2, 2, 24, 200, 3, 2, 2, 2, 26, 208,
	3, 2, 2, 2, 28, 218, 3, 2, 2, 2, 30, 231, 3, 2, 2, 2, 32, 243, 3, 2, 2,
	2, 34, 257, 3, 2, 2, 2, 36, 272, 3, 2, 2, 2, 38, 279, 3, 2, 2, 2, 40, 295,
	3, 2, 2, 2, 42, 302, 3, 2, 2, 2, 44, 304, 3, 2, 2, 2, 46, 311, 3, 2, 2,
	2, 48, 316, 3, 2, 2, 2, 50, 321, 3, 2, 2, 2, 52, 326, 3, 2, 2, 2, 54, 331,
	3, 2, 2, 2, 56, 336, 3, 2, 2, 2, 58, 341, 3, 2, 2, 2, 60, 346, 3, 2, 2,
	2, 62, 351, 3, 2, 2, 2, 64, 356, 3, 2, 2, 2, 66, 361, 3, 2, 2, 2, 68, 366,
	3, 2, 2, 2, 70, 371, 3, 2, 2, 2, 72, 376, 3, 2, 2, 2, 74, 381, 3, 2, 2,
	2, 76, 386, 3, 2, 2, 2, 78, 390, 3, 2, 2, 2, 80, 394, 3, 2, 2, 2, 82, 398,
	3, 2, 2, 2, 84, 402, 3, 2, 2, 2, 86, 410, 3, 2, 2, 2, 88, 421, 3, 2, 2,
	2, 90, 432, 3, 2, 2, 2, 92, 434, 3, 2, 2, 2, 94, 437, 3, 2, 2, 2, 96, 444,
	3, 2, 2, 2, 98, 448, 3, 2, 2, 2, 100, 453, 3, 2, 2, 2, 102, 462, 3, 2,
	2, 2, 104, 467, 3, 2, 2, 2, 106, 474, 3, 2, 2, 2, 108, 478, 3, 2, 2, 2,
	110, 482, 3, 2, 2, 2, 112, 484, 3, 2, 2, 2, 114, 489, 3, 2, 2, 2, 116,
	496, 3, 2, 2, 2, 118, 500, 3, 2, 2, 2, 120, 506, 3, 2, 2, 2, 122, 517,
	3, 2, 2, 2, 124, 519, 3, 2, 2, 2, 126, 522, 3, 2, 2, 2, 128, 529, 3, 2,
	2, 2, 130, 132, 9, 2, 2, 2, 131, 130, 3, 2, 2, 2, 132, 133, 3, 2, 2, 2,
	133, 131, 3, 2, 2, 2, 133, 134, 3, 2, 2, 2, 134, 135, 3, 2, 2, 2, 135,
	136, 7, 60, 2, 2, 136, 137, 3, 2, 2, 2, 137, 138, 8, 2, 2, 2, 138, 9, 3,
	2, 2, 2, 139, 140, 7, 72, 2, 2, 140, 141, 7, 84, 2, 2, 141, 142, 7, 81,
	2, 2, 142, 143, 7, 79, 2, 2, 143, 144, 3, 2, 2, 2, 144, 145, 8, 3, 3, 2,
	145, 11, 3, 2, 2, 2, 146, 147, 7, 69, 2, 2, 147, 148, 7, 81, 2, 2, 148,
	149, 7, 82, 2, 2, 149, 150, 7, 91, 2, 2, 150, 151, 3, 2, 2, 2, 151, 152,
	8, 4, 3, 2, 152, 13, 3, 2, 2, 2, 153, 154, 7, 85, 2, 2, 154, 155, 7, 67,
	2, 2, 155, 156, 7, 88, 2, 2, 156, 157, 7, 71, 2, 2, 157, 158, 7, 34, 2,
	2, 158, 159, 7, 67, 2, 2, 159, 160, 7, 84, 2, 2, 160, 161, 7, 86, 2, 2,
	161, 162, 7, 75, 2, 2, 162, 163, 7, 72, 2, 2, 163, 164, 7, 67, 2, 2, 164,
	165, 7, 69, 2, 2, 165, 166, 7, 86, 2, 2, 166, 167, 3, 2, 2, 2, 167, 168,
	8, 5, 4, 2, 168, 15, 3, 2, 2, 2, 169, 170, 7, 85, 2, 2, 170, 171, 7, 67,
	2, 2, 171, 172, 7, 88, 2, 2, 172, 173, 7, 71, 2, 2, 173, 174, 7, 34, 2,
	2, 174, 175, 7, 75, 2, 2, 175, 176, 7, 79, 2, 2, 176, 177, 7, 67, 2, 2,
	177, 178, 7, 73, 2, 2, 178, 179, 7, 71, 2, 2, 179, 180, 3, 2, 2, 2, 180,
	181, 8, 6, 4, 2, 181, 17, 3, 2, 2, 2, 182, 183, 7, 84, 2, 2, 183, 184,
	7, 87, 2, 2, 184, 185, 7, 80, 2, 2, 185, 186, 3, 2, 2, 2, 186, 187, 8,
	7, 3, 2, 187, 19, 3, 2, 2, 2, 188, 189, 7, 71, 2, 2, 189, 190, 7, 80, 2,
	2, 190, 191, 7, 88, 2, 2, 191, 192, 3, 2, 2, 2, 192, 193, 8, 8, 5, 2, 193,
	21, 3, 2, 2, 2, 194, 195, 7, 67, 2, 2, 195, 196, 7, 84, 2, 2, 196, 197,
	7, 73, 2, 2, 197, 198, 3, 2, 2, 2, 198, 199, 8, 9, 5, 2, 199, 23, 3, 2,
	2, 2, 200, 201, 7, 68, 2, 2, 201, 202, 7, 87, 2, 2, 202, 203, 7, 75, 2,
	2, 203, 204, 7, 78, 2, 2, 204, 205, 7, 70, 2, 2, 205, 206, 3, 2, 2, 2,
	206, 207, 8, 10, 4, 2, 207, 25, 3, 2, 2, 2, 208, 209, 7, 89, 2, 2, 209,
	210, 7, 81, 2, 2, 210, 211, 7, 84, 2, 2, 211, 212, 7, 77, 2, 2, 212, 213,
	7, 70, 2, 2, 213, 214, 7, 75, 2, 2, 214, 215, 7, 84, 2, 2, 215, 216, 3,
	2, 2, 2, 216, 217, 8, 11, 4, 2, 217, 27, 3, 2, 2, 2, 218, 219, 7, 71, 2,
	2, 219, 220, 7, 80, 2, 2, 220, 221, 7, 86, 2, 2, 221, 222, 7, 84, 2, 2,
	222, 223, 7, 91, 2, 2, 223, 224, 7, 82, 2, 2, 224, 225, 7, 81, 2, 2, 225,
	226, 7, 75, 2, 2, 226, 227, 7, 80, 2, 2, 227, 228, 7, 86, 2, 2, 228, 229,
	3, 2, 2, 2, 229, 230, 8, 12, 3, 2, 230, 29, 3, 2, 2, 2, 231, 232, 7, 73,
	2, 2, 232, 233, 7, 75, 2, 2, 233, 234, 7, 86, 2, 2, 234, 235, 7, 34, 2,
	2, 235, 236, 7, 69, 2, 2, 236, 237, 7, 78, 2, 2, 237, 238, 7, 81, 2, 2,
	238, 239, 7, 80, 2, 2, 239, 240, 7, 71, 2, 2, 240, 241, 3, 2, 2, 2, 241,
	242, 8, 13, 4, 2, 242, 31, 3, 2, 2, 2, 243, 244, 7, 70, 2, 2, 244, 245,
	7, 81, 2, 2, 245, 246, 7, 69, 2, 2, 246, 247, 7, 77, 2, 2, 247, 248, 7,
	71, 2, 2, 248, 249, 7, 84, 2, 2, 249, 250, 7, 34, 2, 2, 250, 251, 7, 78,
	2, 2, 251, 252, 7, 81, 2, 2, 252, 253, 7, 67, 2, 2, 253, 254, 7, 70, 2,
	2, 254, 255, 3, 2, 2, 2, 255, 256, 8, 14, 4, 2, 256, 33, 3, 2, 2, 2, 257,
	258, 7, 70, 2, 2, 258, 259, 7, 81, 2, 2, 259, 260, 7, 69, 2, 2, 260, 261,
	7, 77, 2, 2, 261, 262, 7, 71, 2, 2, 262, 263, 7, 84, 2, 2, 263, 264, 7,
	34, 2, 2, 264, 265, 7, 82, 2, 2, 265, 266, 7, 87, 2, 2, 266, 267, 7, 78,
	2, 2, 267, 268, 7, 78, 2, 2, 268, 269, 3, 2, 2, 2, 269, 270, 8, 15, 4,
	2, 270, 35, 3, 2, 2, 2, 271, 273, 9, 3, 2, 2, 272, 271, 3, 2, 2, 2, 273,
	274, 3, 2, 2, 2, 274, 272, 3, 2, 2, 2, 274, 275, 3, 2, 2, 2, 275, 276,
	3, 2, 2, 2, 276, 277, 8, 16, 4, 2, 277, 37, 3, 2, 2, 2, 278, 280, 5, 40,
	18, 2, 279, 278, 3, 2, 2, 2, 279, 280, 3, 2, 2, 2, 280, 282, 3, 2, 2, 2,
	281, 283, 5, 44, 20, 2, 282, 281, 3, 2, 2, 2, 282, 283, 3, 2, 2, 2, 283,
	284, 3, 2, 2, 2, 284, 285, 5, 42, 19, 2, 285, 39, 3, 2, 2, 2, 286, 296,
	9, 4, 2, 2, 287, 291, 7, 94, 2, 2, 288, 290, 9, 4, 2, 2, 289, 288, 3, 2,
	2, 2, 290, 293, 3, 2, 2, 2, 291, 289, 3, 2, 2, 2, 291, 292, 3, 2, 2, 2,
	292, 294, 3, 2, 2, 2, 293, 291, 3, 2, 2, 2, 294, 296, 5, 42, 19, 2, 295,
	286, 3, 2, 2, 2, 295, 287, 3, 2, 2, 2, 296, 297, 3, 2, 2, 2, 297, 295,
	3, 2, 2, 2, 297, 298, 3, 2, 2, 2, 298, 41, 3, 2, 2, 2, 299, 303, 9, 5,
	2, 2, 300, 301, 7, 15, 2, 2, 301, 303, 7, 12, 2, 2, 302, 299, 3, 2, 2,
	2, 302, 300, 3, 2, 2, 2, 303, 43, 3, 2, 2, 2, 304, 308, 7, 37, 2, 2, 305,
	307, 10, 5, 2, 2, 306, 305, 3, 2, 2, 2, 307, 310, 3, 2, 2, 2, 308, 306,
	3, 2, 2, 2, 308, 309, 3, 2, 2, 2, 309, 45, 3, 2, 2, 2, 310, 308, 3, 2,
	2, 2, 311, 312, 5, 8, 2, 2, 312, 313, 3, 2, 2, 2, 313, 314, 8, 21, 6, 2,
	314, 315, 8, 21, 2, 2, 315, 47, 3, 2, 2, 2, 316, 317, 5, 10, 3, 2, 317,
	318, 3, 2, 2, 2, 318, 319, 8, 22, 7, 2, 319, 320, 8, 22, 3, 2, 320, 49,
	3, 2, 2, 2, 321, 322, 5, 12, 4, 2, 322, 323, 3, 2, 2, 2, 323, 324, 8, 23,
	8, 2, 324, 325, 8, 23, 3, 2, 325, 51, 3, 2, 2, 2, 326, 327, 5, 14, 5, 2,
	327, 328, 3, 2, 2, 2, 328, 329, 8, 24, 9, 2, 329, 330, 8, 24, 4, 2, 330,
	53, 3, 2, 2, 2, 331, 332, 5, 16, 6, 2, 332, 333, 3, 2, 2, 2, 333, 334,
	8, 25, 10, 2, 334, 335, 8, 25, 4, 2, 335, 55, 3, 2, 2, 2, 336, 337, 5,
	18, 7, 2, 337, 338, 3, 2, 2, 2, 338, 339, 8, 26, 11, 2, 339, 340, 8, 26,
	3, 2, 340, 57, 3, 2, 2, 2, 341, 342, 5, 20, 8, 2, 342, 343, 3, 2, 2, 2,
	343, 344, 8, 27, 12, 2, 344, 345, 8, 27, 5, 2, 345, 59, 3, 2, 2, 2, 346,
	347, 5, 22, 9, 2, 347, 348, 3, 2, 2, 2, 348, 349, 8, 28, 13, 2, 349, 350,
	8, 28, 5, 2, 350, 61, 3, 2, 2, 2, 351, 352, 5, 24, 10, 2, 352, 353, 3,
	2, 2, 2, 353, 354, 8, 29, 14, 2, 354, 355, 8, 29, 4, 2, 355, 63, 3, 2,
	2, 2, 356, 357, 5, 26, 11, 2, 357, 358, 3, 2, 2, 2, 358, 359, 8, 30, 15,
	2, 359, 360, 8, 30, 4, 2, 360, 65, 3, 2, 2, 2, 361, 362, 5, 28, 12, 2,
	362, 363, 3, 2, 2, 2, 363, 364, 8, 31, 16, 2, 364, 365, 8, 31, 3, 2, 365,
	67, 3, 2, 2, 2, 366, 367, 5, 30, 13, 2, 367, 368, 3, 2, 2, 2, 368, 369,
	8, 32, 17, 2, 369, 370, 8, 32, 4, 2, 370, 69, 3, 2, 2, 2, 371, 372, 5,
	32, 14, 2, 372, 373, 3, 2, 2, 2, 373, 374, 8, 33, 18, 2, 374, 375, 8, 33,
	4, 2, 375, 71, 3, 2, 2, 2, 376, 377, 5, 34, 15, 2, 377, 378, 3, 2, 2, 2,
	378, 379, 8, 34, 19, 2, 379, 380, 8, 34, 4, 2, 380, 73, 3, 2, 2, 2, 381,
	382, 5, 36, 16, 2, 382, 383, 3, 2, 2, 2, 383, 384, 8, 35, 20, 2, 384, 385,
	8, 35, 4, 2, 385, 75, 3, 2, 2, 2, 386, 387, 5, 38, 17, 2, 387, 388, 3,
	2, 2, 2, 388, 389, 8, 36, 21, 2, 389, 77, 3, 2, 2, 2, 390, 391, 5, 40,
	18, 2, 391, 392, 3, 2, 2, 2, 392, 393, 8, 37, 22, 2, 393, 79, 3, 2, 2,
	2, 394, 395, 7, 93, 2, 2, 395, 396, 3, 2, 2, 2, 396, 397, 8, 38, 23, 2,
	397, 81, 3, 2, 2, 2, 398, 399, 5, 84, 40, 2, 399, 400, 7, 63, 2, 2, 400,
	401, 5, 86, 41, 2, 401, 83, 3, 2, 2, 2, 402, 403, 7, 47, 2, 2, 403, 404,
	7, 47, 2, 2, 404, 406, 3, 2, 2, 2, 405, 407, 5, 86, 41, 2, 406, 405, 3,
	2, 2, 2, 406, 407, 3, 2, 2, 2, 407, 85, 3, 2, 2, 2, 408, 411, 5, 90, 43,
	2, 409, 411, 5, 88, 42, 2, 410, 408, 3, 2, 2, 2, 410, 409, 3, 2, 2, 2,
	411, 416, 3, 2, 2, 2, 412, 415, 5, 92, 44, 2, 413, 415, 5, 88, 42, 2, 414,
	412, 3, 2, 2, 2, 414, 413, 3, 2, 2, 2, 415, 418, 3, 2, 2, 2, 416, 414,
	3, 2, 2, 2, 416, 417, 3, 2, 2, 2, 417, 419, 3, 2, 2, 2, 418, 416, 3, 2,
	2, 2, 419, 420, 8, 41, 24, 2, 420, 87, 3, 2, 2, 2, 421, 427, 7, 36, 2,
	2, 422, 426, 10, 6, 2, 2, 423, 424, 7, 94, 2, 2, 424, 426, 7, 36, 2, 2,
	425, 422, 3, 2, 2, 2, 425, 423, 3, 2, 2, 2, 426, 429, 3, 2, 2, 2, 427,
	425, 3, 2, 2, 2, 427, 428, 3, 2, 2, 2, 428, 430, 3, 2, 2, 2, 429, 427,
	3, 2, 2, 2, 430, 431, 7, 36, 2, 2, 431, 89, 3, 2, 2, 2, 432, 433, 10, 7,
	2, 2, 433, 91, 3, 2, 2, 2, 434, 435, 10, 8, 2, 2, 435, 93, 3, 2, 2, 2,
	436, 438, 5, 40, 18, 2, 437, 436, 3, 2, 2, 2, 437, 438, 3, 2, 2, 2, 438,
	439, 3, 2, 2, 2, 439, 440, 5, 42, 19, 2, 440, 441, 3, 2, 2, 2, 441, 442,
	8, 45, 21, 2, 442, 443, 8, 45, 25, 2, 443, 95, 3, 2, 2, 2, 444, 445, 5,
	40, 18, 2, 445, 446, 3, 2, 2, 2, 446, 447, 8, 46, 22, 2, 447, 97, 3, 2,
	2, 2, 448, 449, 7, 93, 2, 2, 449, 450, 3, 2, 2, 2, 450, 451, 8, 47, 26,
	2, 451, 452, 8, 47, 23, 2, 452, 99, 3, 2, 2, 2, 453, 454, 7, 67, 2, 2,
	454, 455, 7, 85, 2, 2, 455, 456, 7, 34, 2, 2, 456, 457, 7, 78, 2, 2, 457,
	458, 7, 81, 2, 2, 458, 459, 7, 69, 2, 2, 459, 460, 7, 67, 2, 2, 460, 461,
	7, 78, 2, 2, 461, 101, 3, 2, 2, 2, 462, 463, 5, 86, 41, 2, 463, 464, 3,
	2, 2, 2, 464, 465, 8, 49, 27, 2, 465, 103, 3, 2, 2, 2, 466, 468, 5, 40,
	18, 2, 467, 466, 3, 2, 2, 2, 467, 468, 3, 2, 2, 2, 468, 469, 3, 2, 2, 2,
	469, 470, 5, 42, 19, 2, 470, 471, 3, 2, 2, 2, 471, 472, 8, 50, 21, 2, 472,
	473, 8, 50, 25, 2, 473, 105, 3, 2, 2, 2, 474, 475, 5, 40, 18, 2, 475, 476,
	3, 2, 2, 2, 476, 477, 8, 51, 22, 2, 477, 107, 3, 2, 2, 2, 478, 479, 7,
	95, 2, 2, 479, 480, 3, 2, 2, 2, 480, 481, 8, 52, 25, 2, 481, 109, 3, 2,
	2, 2, 482, 483, 7, 46, 2, 2, 483, 111, 3, 2, 2, 2, 484, 485, 5, 88, 42,
	2, 485, 486, 3, 2, 2, 2, 486, 487, 8, 54, 27, 2, 487, 113, 3, 2, 2, 2,
	488, 490, 5, 40, 18, 2, 489, 488, 3, 2, 2, 2, 489, 490, 3, 2, 2, 2, 490,
	491, 3, 2, 2, 2, 491, 492, 5, 42, 19, 2, 492, 493, 3, 2, 2, 2, 493, 494,
	8, 55, 21, 2, 494, 495, 8, 55, 25, 2, 495, 115, 3, 2, 2, 2, 496, 497, 5,
	40, 18, 2, 497, 498, 3, 2, 2, 2, 498, 499, 8, 56, 22, 2, 499, 117, 3, 2,
	2, 2, 500, 501, 7, 63, 2, 2, 501, 502, 3, 2, 2, 2, 502, 503, 8, 57, 24,
	2, 503, 119, 3, 2, 2, 2, 504, 507, 5, 122, 59, 2, 505, 507, 5, 88, 42,
	2, 506, 504, 3, 2, 2, 2, 506, 505, 3, 2, 2, 2, 507, 512, 3, 2, 2, 2, 508,
	511, 5, 124, 60, 2, 509, 511, 5, 88, 42, 2, 510, 508, 3, 2, 2, 2, 510,
	509, 3, 2, 2, 2, 511, 514, 3, 2, 2, 2, 512, 510, 3, 2, 2, 2, 512, 513,
	3, 2, 2, 2, 513, 515, 3, 2, 2, 2, 514, 512, 3, 2, 2, 2, 515, 516, 8, 58,
	27, 2, 516, 121, 3, 2, 2, 2, 517, 518, 10, 9, 2, 2, 518, 123, 3, 2, 2,
	2, 519, 520, 10, 10, 2, 2, 520, 125, 3, 2, 2, 2, 521, 523, 5, 40, 18, 2,
	522, 521, 3, 2, 2, 2, 522, 523, 3, 2, 2, 2, 523, 524, 3, 2, 2, 2, 524,
	525, 5, 42, 19, 2, 525, 526, 3, 2, 2, 2, 526, 527, 8, 61, 21, 2, 527, 528,
	8, 61, 25, 2, 528, 127, 3, 2, 2, 2, 529, 530, 5, 40, 18, 2, 530, 531, 3,
	2, 2, 2, 531, 532, 8, 62, 22, 2, 532, 129, 3, 2, 2, 2, 31, 2, 3, 4, 5,
	6, 7, 131, 133, 274, 279, 282, 291, 295, 297, 302, 308, 406, 410, 414,
	416, 425, 427, 437, 467, 489, 506, 510, 512, 522, 28, 7, 3, 2, 7, 5, 2,
	7, 4, 2, 7, 7, 2, 9, 5, 2, 9, 6, 2, 9, 7, 2, 9, 8, 2, 9, 9, 2, 9, 10, 2,
	9, 11, 2, 9, 12, 2, 9, 13, 2, 9, 14, 2, 9, 15, 2, 9, 16, 2, 9, 17, 2, 9,
	18, 2, 9, 19, 2, 9, 20, 2, 9, 21, 2, 7, 6, 2, 4, 5, 2, 6, 2, 2, 9, 22,
	2, 9, 25, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE", "RECIPE", "COMMAND_ARGS", "COMMAND_ARGS_ATOMS_ONLY", "COMMAND_BRACKETS",
	"COMMAND_ARGS_KEY_VALUE",
}

var lexerLiteralNames = []string{
	"", "", "", "", "'FROM'", "'COPY'", "'SAVE ARTIFACT'", "'SAVE IMAGE'",
	"'RUN'", "'ENV'", "'ARG'", "'BUILD'", "'WORKDIR'", "'ENTRYPOINT'", "'GIT CLONE'",
	"'DOCKER LOAD'", "'DOCKER PULL'", "", "", "", "'['", "", "", "", "'AS LOCAL'",
	"']'", "','", "'='",
}

var lexerSymbolicNames = []string{
	"", "INDENT", "DEDENT", "Target", "FROM", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE",
	"RUN", "ENV", "ARG", "BUILD", "WORKDIR", "ENTRYPOINT", "GIT_CLONE", "DOCKER_LOAD",
	"DOCKER_PULL", "Command", "NL", "WS", "OPEN_BRACKET", "FlagKeyValue", "FlagKey",
	"Atom", "AS_LOCAL", "CLOSE_BRACKET", "COMMA", "EQUALS",
}

var lexerRuleNames = []string{
	"Target", "FROM", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE", "RUN", "ENV",
	"ARG", "BUILD", "WORKDIR", "ENTRYPOINT", "GIT_CLONE", "DOCKER_LOAD", "DOCKER_PULL",
	"Command", "NL", "WS", "CRLF", "COMMENT", "Target_R", "FROM_R", "COPY_R",
	"SAVE_ARTIFACT_R", "SAVE_IMAGE_R", "RUN_R", "ENV_R", "ARG_R", "BUILD_R",
	"WORKDIR_R", "ENTRYPOINT_R", "GIT_CLONE_R", "DOCKER_LOAD_R", "DOCKER_PULL_R",
	"Command_R", "NL_R", "WS_R", "OPEN_BRACKET", "FlagKeyValue", "FlagKey",
	"Atom", "QuotedAtom", "NonWSNLQuoteBracket", "NonWSNLQuote", "NL_C", "WS_C",
	"OPEN_BRACKET_CAAO", "AS_LOCAL", "Atom_CAAO", "NL_CAAO", "WS_CAAO", "CLOSE_BRACKET",
	"COMMA", "Atom_CB", "NL_CB", "WS_CB", "EQUALS", "Atom_CAKV", "NonWSNLQuoteBracket_CAKV",
	"NonWSNLQuote_CAKV", "NL_CAKV", "WS_CAKC",
}

type EarthLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewEarthLexer(input antlr.CharStream) *EarthLexer {

	l := new(EarthLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "EarthLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// EarthLexer tokens.
const (
	EarthLexerINDENT        = 1
	EarthLexerDEDENT        = 2
	EarthLexerTarget        = 3
	EarthLexerFROM          = 4
	EarthLexerCOPY          = 5
	EarthLexerSAVE_ARTIFACT = 6
	EarthLexerSAVE_IMAGE    = 7
	EarthLexerRUN           = 8
	EarthLexerENV           = 9
	EarthLexerARG           = 10
	EarthLexerBUILD         = 11
	EarthLexerWORKDIR       = 12
	EarthLexerENTRYPOINT    = 13
	EarthLexerGIT_CLONE     = 14
	EarthLexerDOCKER_LOAD   = 15
	EarthLexerDOCKER_PULL   = 16
	EarthLexerCommand       = 17
	EarthLexerNL            = 18
	EarthLexerWS            = 19
	EarthLexerOPEN_BRACKET  = 20
	EarthLexerFlagKeyValue  = 21
	EarthLexerFlagKey       = 22
	EarthLexerAtom          = 23
	EarthLexerAS_LOCAL      = 24
	EarthLexerCLOSE_BRACKET = 25
	EarthLexerCOMMA         = 26
	EarthLexerEQUALS        = 27
)

// EarthLexer modes.
const (
	EarthLexerRECIPE = iota + 1
	EarthLexerCOMMAND_ARGS
	EarthLexerCOMMAND_ARGS_ATOMS_ONLY
	EarthLexerCOMMAND_BRACKETS
	EarthLexerCOMMAND_ARGS_KEY_VALUE
)
