// Code generated from earthfile2llb/parser/EarthParser.g4 by ANTLR 4.8. DO NOT EDIT.

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 32, 563,
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
	9, 55, 4, 56, 9, 56, 3, 2, 7, 2, 114, 10, 2, 12, 2, 14, 2, 117, 11, 2,
	3, 2, 5, 2, 120, 10, 2, 3, 2, 6, 2, 123, 10, 2, 13, 2, 14, 2, 124, 3, 2,
	5, 2, 128, 10, 2, 3, 2, 7, 2, 131, 10, 2, 12, 2, 14, 2, 134, 11, 2, 3,
	2, 3, 2, 3, 3, 3, 3, 5, 3, 140, 10, 3, 3, 3, 6, 3, 143, 10, 3, 13, 3, 14,
	3, 144, 3, 3, 3, 3, 3, 3, 5, 3, 150, 10, 3, 7, 3, 152, 10, 3, 12, 3, 14,
	3, 155, 11, 3, 3, 3, 7, 3, 158, 10, 3, 12, 3, 14, 3, 161, 11, 3, 3, 3,
	5, 3, 164, 10, 3, 3, 4, 3, 4, 6, 4, 168, 10, 4, 13, 4, 14, 4, 169, 3, 4,
	5, 4, 173, 10, 4, 3, 4, 3, 4, 5, 4, 177, 10, 4, 3, 5, 3, 5, 3, 6, 5, 6,
	182, 10, 6, 3, 6, 3, 6, 6, 6, 186, 10, 6, 13, 6, 14, 6, 187, 3, 6, 5, 6,
	191, 10, 6, 3, 6, 7, 6, 194, 10, 6, 12, 6, 14, 6, 197, 11, 6, 3, 7, 3,
	7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 5,
	7, 212, 10, 7, 3, 8, 3, 8, 3, 8, 7, 8, 217, 10, 8, 12, 8, 14, 8, 220, 11,
	8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 5, 8, 227, 10, 8, 3, 9, 3, 9, 3, 9, 3,
	9, 3, 9, 5, 9, 234, 10, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3,
	10, 3, 10, 3, 11, 3, 11, 3, 11, 7, 11, 247, 10, 11, 12, 11, 14, 11, 250,
	11, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 3, 12,
	3, 12, 5, 12, 262, 10, 12, 3, 13, 3, 13, 3, 13, 7, 13, 267, 10, 13, 12,
	13, 14, 13, 270, 11, 13, 3, 14, 3, 14, 5, 14, 274, 10, 14, 3, 15, 3, 15,
	3, 15, 7, 15, 279, 10, 15, 12, 15, 14, 15, 282, 11, 15, 3, 16, 3, 16, 3,
	16, 3, 16, 3, 16, 5, 16, 289, 10, 16, 3, 16, 3, 16, 3, 16, 3, 16, 5, 16,
	295, 10, 16, 3, 17, 3, 17, 3, 17, 7, 17, 300, 10, 17, 12, 17, 14, 17, 303,
	11, 17, 3, 17, 3, 17, 3, 17, 5, 17, 308, 10, 17, 3, 18, 3, 18, 3, 18, 7,
	18, 313, 10, 18, 12, 18, 14, 18, 316, 11, 18, 3, 18, 3, 18, 3, 18, 3, 19,
	3, 19, 3, 19, 3, 19, 3, 20, 3, 20, 3, 20, 3, 20, 5, 20, 329, 10, 20, 3,
	21, 3, 21, 3, 21, 3, 21, 5, 21, 335, 10, 21, 3, 21, 5, 21, 338, 10, 21,
	3, 21, 5, 21, 341, 10, 21, 3, 21, 5, 21, 344, 10, 21, 3, 22, 3, 22, 3,
	22, 3, 22, 5, 22, 350, 10, 22, 3, 22, 3, 22, 3, 22, 5, 22, 355, 10, 22,
	3, 22, 5, 22, 358, 10, 22, 5, 22, 360, 10, 22, 3, 23, 3, 23, 3, 23, 7,
	23, 365, 10, 23, 12, 23, 14, 23, 368, 11, 23, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 24, 3, 24, 3, 24, 7, 24, 378, 10, 24, 12, 24, 14, 24, 381, 11,
	24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25,
	3, 25, 3, 26, 3, 26, 3, 26, 5, 26, 397, 10, 26, 3, 26, 3, 26, 3, 26, 3,
	26, 5, 26, 403, 10, 26, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 7, 28, 410,
	10, 28, 12, 28, 14, 28, 413, 11, 28, 3, 29, 3, 29, 5, 29, 417, 10, 29,
	3, 29, 3, 29, 5, 29, 421, 10, 29, 3, 29, 3, 29, 5, 29, 425, 10, 29, 3,
	29, 6, 29, 428, 10, 29, 13, 29, 14, 29, 429, 3, 29, 5, 29, 433, 10, 29,
	3, 29, 3, 29, 3, 30, 3, 30, 3, 31, 3, 31, 3, 31, 7, 31, 442, 10, 31, 12,
	31, 14, 31, 445, 11, 31, 3, 32, 3, 32, 5, 32, 449, 10, 32, 3, 32, 3, 32,
	5, 32, 453, 10, 32, 3, 32, 3, 32, 5, 32, 457, 10, 32, 3, 32, 6, 32, 460,
	10, 32, 13, 32, 14, 32, 461, 3, 32, 5, 32, 465, 10, 32, 3, 32, 3, 32, 3,
	33, 3, 33, 3, 34, 3, 34, 5, 34, 473, 10, 34, 3, 34, 7, 34, 476, 10, 34,
	12, 34, 14, 34, 479, 11, 34, 3, 35, 3, 35, 5, 35, 483, 10, 35, 3, 36, 3,
	36, 3, 37, 3, 37, 3, 38, 3, 38, 5, 38, 491, 10, 38, 3, 38, 7, 38, 494,
	10, 38, 12, 38, 14, 38, 497, 11, 38, 3, 39, 3, 39, 3, 40, 3, 40, 3, 41,
	3, 41, 5, 41, 505, 10, 41, 3, 41, 7, 41, 508, 10, 41, 12, 41, 14, 41, 511,
	11, 41, 3, 42, 3, 42, 3, 43, 3, 43, 3, 44, 3, 44, 3, 45, 3, 45, 3, 46,
	3, 46, 3, 47, 3, 47, 3, 48, 3, 48, 3, 49, 3, 49, 3, 50, 3, 50, 3, 51, 3,
	51, 3, 52, 3, 52, 3, 53, 3, 53, 3, 54, 3, 54, 3, 55, 3, 55, 5, 55, 541,
	10, 55, 3, 55, 3, 55, 5, 55, 545, 10, 55, 3, 55, 3, 55, 5, 55, 549, 10,
	55, 3, 55, 6, 55, 552, 10, 55, 13, 55, 14, 55, 553, 3, 55, 5, 55, 557,
	10, 55, 3, 55, 3, 55, 3, 56, 3, 56, 3, 56, 2, 2, 57, 2, 4, 6, 8, 10, 12,
	14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48,
	50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84,
	86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 108, 110, 2, 2, 2, 589,
	2, 115, 3, 2, 2, 2, 4, 137, 3, 2, 2, 2, 6, 165, 3, 2, 2, 2, 8, 178, 3,
	2, 2, 2, 10, 181, 3, 2, 2, 2, 12, 211, 3, 2, 2, 2, 14, 213, 3, 2, 2, 2,
	16, 228, 3, 2, 2, 2, 18, 235, 3, 2, 2, 2, 20, 243, 3, 2, 2, 2, 22, 261,
	3, 2, 2, 2, 24, 263, 3, 2, 2, 2, 26, 273, 3, 2, 2, 2, 28, 275, 3, 2, 2,
	2, 30, 283, 3, 2, 2, 2, 32, 296, 3, 2, 2, 2, 34, 309, 3, 2, 2, 2, 36, 320,
	3, 2, 2, 2, 38, 324, 3, 2, 2, 2, 40, 330, 3, 2, 2, 2, 42, 345, 3, 2, 2,
	2, 44, 361, 3, 2, 2, 2, 46, 374, 3, 2, 2, 2, 48, 389, 3, 2, 2, 2, 50, 393,
	3, 2, 2, 2, 52, 404, 3, 2, 2, 2, 54, 406, 3, 2, 2, 2, 56, 414, 3, 2, 2,
	2, 58, 436, 3, 2, 2, 2, 60, 438, 3, 2, 2, 2, 62, 446, 3, 2, 2, 2, 64, 468,
	3, 2, 2, 2, 66, 470, 3, 2, 2, 2, 68, 482, 3, 2, 2, 2, 70, 484, 3, 2, 2,
	2, 72, 486, 3, 2, 2, 2, 74, 488, 3, 2, 2, 2, 76, 498, 3, 2, 2, 2, 78, 500,
	3, 2, 2, 2, 80, 502, 3, 2, 2, 2, 82, 512, 3, 2, 2, 2, 84, 514, 3, 2, 2,
	2, 86, 516, 3, 2, 2, 2, 88, 518, 3, 2, 2, 2, 90, 520, 3, 2, 2, 2, 92, 522,
	3, 2, 2, 2, 94, 524, 3, 2, 2, 2, 96, 526, 3, 2, 2, 2, 98, 528, 3, 2, 2,
	2, 100, 530, 3, 2, 2, 2, 102, 532, 3, 2, 2, 2, 104, 534, 3, 2, 2, 2, 106,
	536, 3, 2, 2, 2, 108, 538, 3, 2, 2, 2, 110, 560, 3, 2, 2, 2, 112, 114,
	7, 20, 2, 2, 113, 112, 3, 2, 2, 2, 114, 117, 3, 2, 2, 2, 115, 113, 3, 2,
	2, 2, 115, 116, 3, 2, 2, 2, 116, 119, 3, 2, 2, 2, 117, 115, 3, 2, 2, 2,
	118, 120, 5, 10, 6, 2, 119, 118, 3, 2, 2, 2, 119, 120, 3, 2, 2, 2, 120,
	122, 3, 2, 2, 2, 121, 123, 7, 20, 2, 2, 122, 121, 3, 2, 2, 2, 123, 124,
	3, 2, 2, 2, 124, 122, 3, 2, 2, 2, 124, 125, 3, 2, 2, 2, 125, 127, 3, 2,
	2, 2, 126, 128, 5, 4, 3, 2, 127, 126, 3, 2, 2, 2, 127, 128, 3, 2, 2, 2,
	128, 132, 3, 2, 2, 2, 129, 131, 7, 20, 2, 2, 130, 129, 3, 2, 2, 2, 131,
	134, 3, 2, 2, 2, 132, 130, 3, 2, 2, 2, 132, 133, 3, 2, 2, 2, 133, 135,
	3, 2, 2, 2, 134, 132, 3, 2, 2, 2, 135, 136, 7, 2, 2, 3, 136, 3, 3, 2, 2,
	2, 137, 139, 5, 6, 4, 2, 138, 140, 7, 21, 2, 2, 139, 138, 3, 2, 2, 2, 139,
	140, 3, 2, 2, 2, 140, 153, 3, 2, 2, 2, 141, 143, 7, 20, 2, 2, 142, 141,
	3, 2, 2, 2, 143, 144, 3, 2, 2, 2, 144, 142, 3, 2, 2, 2, 144, 145, 3, 2,
	2, 2, 145, 146, 3, 2, 2, 2, 146, 147, 7, 4, 2, 2, 147, 149, 5, 6, 4, 2,
	148, 150, 7, 21, 2, 2, 149, 148, 3, 2, 2, 2, 149, 150, 3, 2, 2, 2, 150,
	152, 3, 2, 2, 2, 151, 142, 3, 2, 2, 2, 152, 155, 3, 2, 2, 2, 153, 151,
	3, 2, 2, 2, 153, 154, 3, 2, 2, 2, 154, 159, 3, 2, 2, 2, 155, 153, 3, 2,
	2, 2, 156, 158, 7, 20, 2, 2, 157, 156, 3, 2, 2, 2, 158, 161, 3, 2, 2, 2,
	159, 157, 3, 2, 2, 2, 159, 160, 3, 2, 2, 2, 160, 163, 3, 2, 2, 2, 161,
	159, 3, 2, 2, 2, 162, 164, 7, 4, 2, 2, 163, 162, 3, 2, 2, 2, 163, 164,
	3, 2, 2, 2, 164, 5, 3, 2, 2, 2, 165, 167, 5, 8, 5, 2, 166, 168, 7, 20,
	2, 2, 167, 166, 3, 2, 2, 2, 168, 169, 3, 2, 2, 2, 169, 167, 3, 2, 2, 2,
	169, 170, 3, 2, 2, 2, 170, 172, 3, 2, 2, 2, 171, 173, 7, 21, 2, 2, 172,
	171, 3, 2, 2, 2, 172, 173, 3, 2, 2, 2, 173, 174, 3, 2, 2, 2, 174, 176,
	7, 3, 2, 2, 175, 177, 5, 10, 6, 2, 176, 175, 3, 2, 2, 2, 176, 177, 3, 2,
	2, 2, 177, 7, 3, 2, 2, 2, 178, 179, 7, 5, 2, 2, 179, 9, 3, 2, 2, 2, 180,
	182, 7, 21, 2, 2, 181, 180, 3, 2, 2, 2, 181, 182, 3, 2, 2, 2, 182, 183,
	3, 2, 2, 2, 183, 195, 5, 12, 7, 2, 184, 186, 7, 20, 2, 2, 185, 184, 3,
	2, 2, 2, 186, 187, 3, 2, 2, 2, 187, 185, 3, 2, 2, 2, 187, 188, 3, 2, 2,
	2, 188, 190, 3, 2, 2, 2, 189, 191, 7, 21, 2, 2, 190, 189, 3, 2, 2, 2, 190,
	191, 3, 2, 2, 2, 191, 192, 3, 2, 2, 2, 192, 194, 5, 12, 7, 2, 193, 185,
	3, 2, 2, 2, 194, 197, 3, 2, 2, 2, 195, 193, 3, 2, 2, 2, 195, 196, 3, 2,
	2, 2, 196, 11, 3, 2, 2, 2, 197, 195, 3, 2, 2, 2, 198, 212, 5, 14, 8, 2,
	199, 212, 5, 16, 9, 2, 200, 212, 5, 26, 14, 2, 201, 212, 5, 32, 17, 2,
	202, 212, 5, 34, 18, 2, 203, 212, 5, 36, 19, 2, 204, 212, 5, 38, 20, 2,
	205, 212, 5, 40, 21, 2, 206, 212, 5, 42, 22, 2, 207, 212, 5, 44, 23, 2,
	208, 212, 5, 46, 24, 2, 209, 212, 5, 48, 25, 2, 210, 212, 5, 50, 26, 2,
	211, 198, 3, 2, 2, 2, 211, 199, 3, 2, 2, 2, 211, 200, 3, 2, 2, 2, 211,
	201, 3, 2, 2, 2, 211, 202, 3, 2, 2, 2, 211, 203, 3, 2, 2, 2, 211, 204,
	3, 2, 2, 2, 211, 205, 3, 2, 2, 2, 211, 206, 3, 2, 2, 2, 211, 207, 3, 2,
	2, 2, 211, 208, 3, 2, 2, 2, 211, 209, 3, 2, 2, 2, 211, 210, 3, 2, 2, 2,
	212, 13, 3, 2, 2, 2, 213, 218, 7, 6, 2, 2, 214, 215, 7, 21, 2, 2, 215,
	217, 5, 72, 37, 2, 216, 214, 3, 2, 2, 2, 217, 220, 3, 2, 2, 2, 218, 216,
	3, 2, 2, 2, 218, 219, 3, 2, 2, 2, 219, 221, 3, 2, 2, 2, 220, 218, 3, 2,
	2, 2, 221, 222, 7, 21, 2, 2, 222, 226, 5, 88, 45, 2, 223, 224, 7, 21, 2,
	2, 224, 225, 7, 28, 2, 2, 225, 227, 5, 86, 44, 2, 226, 223, 3, 2, 2, 2,
	226, 227, 3, 2, 2, 2, 227, 15, 3, 2, 2, 2, 228, 229, 7, 7, 2, 2, 229, 233,
	7, 21, 2, 2, 230, 234, 5, 18, 10, 2, 231, 234, 5, 20, 11, 2, 232, 234,
	5, 22, 12, 2, 233, 230, 3, 2, 2, 2, 233, 231, 3, 2, 2, 2, 233, 232, 3,
	2, 2, 2, 234, 17, 3, 2, 2, 2, 235, 236, 7, 24, 2, 2, 236, 237, 7, 21, 2,
	2, 237, 238, 5, 88, 45, 2, 238, 239, 7, 21, 2, 2, 239, 240, 5, 76, 39,
	2, 240, 241, 7, 21, 2, 2, 241, 242, 5, 74, 38, 2, 242, 19, 3, 2, 2, 2,
	243, 248, 7, 23, 2, 2, 244, 245, 7, 21, 2, 2, 245, 247, 5, 72, 37, 2, 246,
	244, 3, 2, 2, 2, 247, 250, 3, 2, 2, 2, 248, 246, 3, 2, 2, 2, 248, 249,
	3, 2, 2, 2, 249, 251, 3, 2, 2, 2, 250, 248, 3, 2, 2, 2, 251, 252, 7, 21,
	2, 2, 252, 253, 5, 96, 49, 2, 253, 254, 7, 21, 2, 2, 254, 255, 5, 84, 43,
	2, 255, 21, 3, 2, 2, 2, 256, 257, 5, 24, 13, 2, 257, 258, 7, 21, 2, 2,
	258, 259, 5, 84, 43, 2, 259, 262, 3, 2, 2, 2, 260, 262, 5, 108, 55, 2,
	261, 256, 3, 2, 2, 2, 261, 260, 3, 2, 2, 2, 262, 23, 3, 2, 2, 2, 263, 268,
	5, 82, 42, 2, 264, 265, 7, 21, 2, 2, 265, 267, 5, 82, 42, 2, 266, 264,
	3, 2, 2, 2, 267, 270, 3, 2, 2, 2, 268, 266, 3, 2, 2, 2, 268, 269, 3, 2,
	2, 2, 269, 25, 3, 2, 2, 2, 270, 268, 3, 2, 2, 2, 271, 274, 5, 30, 16, 2,
	272, 274, 5, 28, 15, 2, 273, 271, 3, 2, 2, 2, 273, 272, 3, 2, 2, 2, 274,
	27, 3, 2, 2, 2, 275, 280, 7, 9, 2, 2, 276, 277, 7, 21, 2, 2, 277, 279,
	5, 90, 46, 2, 278, 276, 3, 2, 2, 2, 279, 282, 3, 2, 2, 2, 280, 278, 3,
	2, 2, 2, 280, 281, 3, 2, 2, 2, 281, 29, 3, 2, 2, 2, 282, 280, 3, 2, 2,
	2, 283, 284, 7, 8, 2, 2, 284, 285, 7, 21, 2, 2, 285, 288, 5, 98, 50, 2,
	286, 287, 7, 21, 2, 2, 287, 289, 5, 100, 51, 2, 288, 286, 3, 2, 2, 2, 288,
	289, 3, 2, 2, 2, 289, 294, 3, 2, 2, 2, 290, 291, 7, 21, 2, 2, 291, 292,
	7, 29, 2, 2, 292, 293, 7, 21, 2, 2, 293, 295, 5, 102, 52, 2, 294, 290,
	3, 2, 2, 2, 294, 295, 3, 2, 2, 2, 295, 31, 3, 2, 2, 2, 296, 301, 7, 10,
	2, 2, 297, 298, 7, 21, 2, 2, 298, 300, 5, 68, 35, 2, 299, 297, 3, 2, 2,
	2, 300, 303, 3, 2, 2, 2, 301, 299, 3, 2, 2, 2, 301, 302, 3, 2, 2, 2, 302,
	304, 3, 2, 2, 2, 303, 301, 3, 2, 2, 2, 304, 307, 7, 21, 2, 2, 305, 308,
	5, 54, 28, 2, 306, 308, 5, 56, 29, 2, 307, 305, 3, 2, 2, 2, 307, 306, 3,
	2, 2, 2, 308, 33, 3, 2, 2, 2, 309, 314, 7, 13, 2, 2, 310, 311, 7, 21, 2,
	2, 311, 313, 5, 72, 37, 2, 312, 310, 3, 2, 2, 2, 313, 316, 3, 2, 2, 2,
	314, 312, 3, 2, 2, 2, 314, 315, 3, 2, 2, 2, 315, 317, 3, 2, 2, 2, 316,
	314, 3, 2, 2, 2, 317, 318, 7, 21, 2, 2, 318, 319, 5, 94, 48, 2, 319, 35,
	3, 2, 2, 2, 320, 321, 7, 14, 2, 2, 321, 322, 7, 21, 2, 2, 322, 323, 5,
	104, 53, 2, 323, 37, 3, 2, 2, 2, 324, 325, 7, 15, 2, 2, 325, 328, 7, 21,
	2, 2, 326, 329, 5, 60, 31, 2, 327, 329, 5, 62, 32, 2, 328, 326, 3, 2, 2,
	2, 328, 327, 3, 2, 2, 2, 329, 39, 3, 2, 2, 2, 330, 331, 7, 11, 2, 2, 331,
	332, 7, 21, 2, 2, 332, 337, 5, 78, 40, 2, 333, 335, 7, 21, 2, 2, 334, 333,
	3, 2, 2, 2, 334, 335, 3, 2, 2, 2, 335, 336, 3, 2, 2, 2, 336, 338, 7, 32,
	2, 2, 337, 334, 3, 2, 2, 2, 337, 338, 3, 2, 2, 2, 338, 343, 3, 2, 2, 2,
	339, 341, 7, 21, 2, 2, 340, 339, 3, 2, 2, 2, 340, 341, 3, 2, 2, 2, 341,
	342, 3, 2, 2, 2, 342, 344, 5, 80, 41, 2, 343, 340, 3, 2, 2, 2, 343, 344,
	3, 2, 2, 2, 344, 41, 3, 2, 2, 2, 345, 346, 7, 12, 2, 2, 346, 347, 7, 21,
	2, 2, 347, 359, 5, 78, 40, 2, 348, 350, 7, 21, 2, 2, 349, 348, 3, 2, 2,
	2, 349, 350, 3, 2, 2, 2, 350, 351, 3, 2, 2, 2, 351, 352, 7, 32, 2, 2, 352,
	357, 3, 2, 2, 2, 353, 355, 7, 21, 2, 2, 354, 353, 3, 2, 2, 2, 354, 355,
	3, 2, 2, 2, 355, 356, 3, 2, 2, 2, 356, 358, 5, 80, 41, 2, 357, 354, 3,
	2, 2, 2, 357, 358, 3, 2, 2, 2, 358, 360, 3, 2, 2, 2, 359, 349, 3, 2, 2,
	2, 359, 360, 3, 2, 2, 2, 360, 43, 3, 2, 2, 2, 361, 366, 7, 16, 2, 2, 362,
	363, 7, 21, 2, 2, 363, 365, 5, 72, 37, 2, 364, 362, 3, 2, 2, 2, 365, 368,
	3, 2, 2, 2, 366, 364, 3, 2, 2, 2, 366, 367, 3, 2, 2, 2, 367, 369, 3, 2,
	2, 2, 368, 366, 3, 2, 2, 2, 369, 370, 7, 21, 2, 2, 370, 371, 5, 106, 54,
	2, 371, 372, 7, 21, 2, 2, 372, 373, 5, 84, 43, 2, 373, 45, 3, 2, 2, 2,
	374, 379, 7, 17, 2, 2, 375, 376, 7, 21, 2, 2, 376, 378, 5, 72, 37, 2, 377,
	375, 3, 2, 2, 2, 378, 381, 3, 2, 2, 2, 379, 377, 3, 2, 2, 2, 379, 380,
	3, 2, 2, 2, 380, 382, 3, 2, 2, 2, 381, 379, 3, 2, 2, 2, 382, 383, 7, 21,
	2, 2, 383, 384, 5, 94, 48, 2, 384, 385, 7, 21, 2, 2, 385, 386, 7, 28, 2,
	2, 386, 387, 7, 21, 2, 2, 387, 388, 5, 88, 45, 2, 388, 47, 3, 2, 2, 2,
	389, 390, 7, 18, 2, 2, 390, 391, 7, 21, 2, 2, 391, 392, 5, 88, 45, 2, 392,
	49, 3, 2, 2, 2, 393, 396, 5, 52, 27, 2, 394, 395, 7, 21, 2, 2, 395, 397,
	5, 66, 34, 2, 396, 394, 3, 2, 2, 2, 396, 397, 3, 2, 2, 2, 397, 402, 3,
	2, 2, 2, 398, 399, 7, 21, 2, 2, 399, 403, 5, 74, 38, 2, 400, 401, 7, 21,
	2, 2, 401, 403, 5, 108, 55, 2, 402, 398, 3, 2, 2, 2, 402, 400, 3, 2, 2,
	2, 402, 403, 3, 2, 2, 2, 403, 51, 3, 2, 2, 2, 404, 405, 7, 19, 2, 2, 405,
	53, 3, 2, 2, 2, 406, 411, 5, 58, 30, 2, 407, 408, 7, 21, 2, 2, 408, 410,
	5, 58, 30, 2, 409, 407, 3, 2, 2, 2, 410, 413, 3, 2, 2, 2, 411, 409, 3,
	2, 2, 2, 411, 412, 3, 2, 2, 2, 412, 55, 3, 2, 2, 2, 413, 411, 3, 2, 2,
	2, 414, 416, 7, 22, 2, 2, 415, 417, 7, 21, 2, 2, 416, 415, 3, 2, 2, 2,
	416, 417, 3, 2, 2, 2, 417, 418, 3, 2, 2, 2, 418, 427, 5, 58, 30, 2, 419,
	421, 7, 21, 2, 2, 420, 419, 3, 2, 2, 2, 420, 421, 3, 2, 2, 2, 421, 422,
	3, 2, 2, 2, 422, 424, 7, 31, 2, 2, 423, 425, 7, 21, 2, 2, 424, 423, 3,
	2, 2, 2, 424, 425, 3, 2, 2, 2, 425, 426, 3, 2, 2, 2, 426, 428, 5, 58, 30,
	2, 427, 420, 3, 2, 2, 2, 428, 429, 3, 2, 2, 2, 429, 427, 3, 2, 2, 2, 429,
	430, 3, 2, 2, 2, 430, 432, 3, 2, 2, 2, 431, 433, 7, 21, 2, 2, 432, 431,
	3, 2, 2, 2, 432, 433, 3, 2, 2, 2, 433, 434, 3, 2, 2, 2, 434, 435, 7, 30,
	2, 2, 435, 57, 3, 2, 2, 2, 436, 437, 7, 27, 2, 2, 437, 59, 3, 2, 2, 2,
	438, 443, 5, 64, 33, 2, 439, 440, 7, 21, 2, 2, 440, 442, 5, 64, 33, 2,
	441, 439, 3, 2, 2, 2, 442, 445, 3, 2, 2, 2, 443, 441, 3, 2, 2, 2, 443,
	444, 3, 2, 2, 2, 444, 61, 3, 2, 2, 2, 445, 443, 3, 2, 2, 2, 446, 448, 7,
	22, 2, 2, 447, 449, 7, 21, 2, 2, 448, 447, 3, 2, 2, 2, 448, 449, 3, 2,
	2, 2, 449, 450, 3, 2, 2, 2, 450, 459, 5, 64, 33, 2, 451, 453, 7, 21, 2,
	2, 452, 451, 3, 2, 2, 2, 452, 453, 3, 2, 2, 2, 453, 454, 3, 2, 2, 2, 454,
	456, 7, 31, 2, 2, 455, 457, 7, 21, 2, 2, 456, 455, 3, 2, 2, 2, 456, 457,
	3, 2, 2, 2, 457, 458, 3, 2, 2, 2, 458, 460, 5, 64, 33, 2, 459, 452, 3,
	2, 2, 2, 460, 461, 3, 2, 2, 2, 461, 459, 3, 2, 2, 2, 461, 462, 3, 2, 2,
	2, 462, 464, 3, 2, 2, 2, 463, 465, 7, 21, 2, 2, 464, 463, 3, 2, 2, 2, 464,
	465, 3, 2, 2, 2, 465, 466, 3, 2, 2, 2, 466, 467, 7, 30, 2, 2, 467, 63,
	3, 2, 2, 2, 468, 469, 7, 27, 2, 2, 469, 65, 3, 2, 2, 2, 470, 477, 5, 68,
	35, 2, 471, 473, 7, 21, 2, 2, 472, 471, 3, 2, 2, 2, 472, 473, 3, 2, 2,
	2, 473, 474, 3, 2, 2, 2, 474, 476, 5, 68, 35, 2, 475, 472, 3, 2, 2, 2,
	476, 479, 3, 2, 2, 2, 477, 475, 3, 2, 2, 2, 477, 478, 3, 2, 2, 2, 478,
	67, 3, 2, 2, 2, 479, 477, 3, 2, 2, 2, 480, 483, 5, 70, 36, 2, 481, 483,
	5, 72, 37, 2, 482, 480, 3, 2, 2, 2, 482, 481, 3, 2, 2, 2, 483, 69, 3, 2,
	2, 2, 484, 485, 7, 26, 2, 2, 485, 71, 3, 2, 2, 2, 486, 487, 7, 25, 2, 2,
	487, 73, 3, 2, 2, 2, 488, 495, 5, 76, 39, 2, 489, 491, 7, 21, 2, 2, 490,
	489, 3, 2, 2, 2, 490, 491, 3, 2, 2, 2, 491, 492, 3, 2, 2, 2, 492, 494,
	5, 76, 39, 2, 493, 490, 3, 2, 2, 2, 494, 497, 3, 2, 2, 2, 495, 493, 3,
	2, 2, 2, 495, 496, 3, 2, 2, 2, 496, 75, 3, 2, 2, 2, 497, 495, 3, 2, 2,
	2, 498, 499, 7, 27, 2, 2, 499, 77, 3, 2, 2, 2, 500, 501, 7, 27, 2, 2, 501,
	79, 3, 2, 2, 2, 502, 509, 7, 27, 2, 2, 503, 505, 7, 21, 2, 2, 504, 503,
	3, 2, 2, 2, 504, 505, 3, 2, 2, 2, 505, 506, 3, 2, 2, 2, 506, 508, 7, 27,
	2, 2, 507, 504, 3, 2, 2, 2, 508, 511, 3, 2, 2, 2, 509, 507, 3, 2, 2, 2,
	509, 510, 3, 2, 2, 2, 510, 81, 3, 2, 2, 2, 511, 509, 3, 2, 2, 2, 512, 513,
	7, 27, 2, 2, 513, 83, 3, 2, 2, 2, 514, 515, 7, 27, 2, 2, 515, 85, 3, 2,
	2, 2, 516, 517, 7, 27, 2, 2, 517, 87, 3, 2, 2, 2, 518, 519, 7, 27, 2, 2,
	519, 89, 3, 2, 2, 2, 520, 521, 7, 27, 2, 2, 521, 91, 3, 2, 2, 2, 522, 523,
	7, 27, 2, 2, 523, 93, 3, 2, 2, 2, 524, 525, 7, 27, 2, 2, 525, 95, 3, 2,
	2, 2, 526, 527, 7, 27, 2, 2, 527, 97, 3, 2, 2, 2, 528, 529, 7, 27, 2, 2,
	529, 99, 3, 2, 2, 2, 530, 531, 7, 27, 2, 2, 531, 101, 3, 2, 2, 2, 532,
	533, 7, 27, 2, 2, 533, 103, 3, 2, 2, 2, 534, 535, 7, 27, 2, 2, 535, 105,
	3, 2, 2, 2, 536, 537, 7, 27, 2, 2, 537, 107, 3, 2, 2, 2, 538, 540, 7, 22,
	2, 2, 539, 541, 7, 21, 2, 2, 540, 539, 3, 2, 2, 2, 540, 541, 3, 2, 2, 2,
	541, 542, 3, 2, 2, 2, 542, 551, 5, 110, 56, 2, 543, 545, 7, 21, 2, 2, 544,
	543, 3, 2, 2, 2, 544, 545, 3, 2, 2, 2, 545, 546, 3, 2, 2, 2, 546, 548,
	7, 31, 2, 2, 547, 549, 7, 21, 2, 2, 548, 547, 3, 2, 2, 2, 548, 549, 3,
	2, 2, 2, 549, 550, 3, 2, 2, 2, 550, 552, 5, 110, 56, 2, 551, 544, 3, 2,
	2, 2, 552, 553, 3, 2, 2, 2, 553, 551, 3, 2, 2, 2, 553, 554, 3, 2, 2, 2,
	554, 556, 3, 2, 2, 2, 555, 557, 7, 21, 2, 2, 556, 555, 3, 2, 2, 2, 556,
	557, 3, 2, 2, 2, 557, 558, 3, 2, 2, 2, 558, 559, 7, 30, 2, 2, 559, 109,
	3, 2, 2, 2, 560, 561, 7, 27, 2, 2, 561, 111, 3, 2, 2, 2, 71, 115, 119,
	124, 127, 132, 139, 144, 149, 153, 159, 163, 169, 172, 176, 181, 187, 190,
	195, 211, 218, 226, 233, 248, 261, 268, 273, 280, 288, 294, 301, 307, 314,
	328, 334, 337, 340, 343, 349, 354, 357, 359, 366, 379, 396, 402, 411, 416,
	420, 424, 429, 432, 443, 448, 452, 456, 461, 464, 472, 477, 482, 490, 495,
	504, 509, 540, 544, 548, 553, 556,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "'FROM'", "'COPY'", "'SAVE ARTIFACT'", "'SAVE IMAGE'",
	"'RUN'", "'ENV'", "'ARG'", "'BUILD'", "'WORKDIR'", "'ENTRYPOINT'", "'GIT CLONE'",
	"'DOCKER LOAD'", "'DOCKER PULL'", "", "", "", "'['", "'--artifact'", "'--from'",
	"", "", "", "'AS'", "'AS LOCAL'", "']'", "','", "'='",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "Target", "FROM", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE",
	"RUN", "ENV", "ARG", "BUILD", "WORKDIR", "ENTRYPOINT", "GIT_CLONE", "DOCKER_LOAD",
	"DOCKER_PULL", "Command", "NL", "WS", "OPEN_BRACKET", "FLAG_ARTIFACT",
	"FLAG_FROM", "FlagKeyValue", "FlagKey", "Atom", "AS", "AS_LOCAL", "CLOSE_BRACKET",
	"COMMA", "EQUALS",
}

var ruleNames = []string{
	"earthFile", "targets", "target", "targetHeader", "stmts", "stmt", "fromStmt",
	"copyStmt", "copyArgsFrom", "copyArgsArtifact", "copyArgsClassical", "copySrcs",
	"saveStmt", "saveImage", "saveArtifact", "runStmt", "buildStmt", "workdirStmt",
	"entrypointStmt", "envStmt", "argStmt", "gitCloneStmt", "dockerLoadStmt",
	"dockerPullStmt", "genericCommand", "commandName", "runArgs", "runArgsList",
	"runArg", "entrypointArgs", "entrypointArgsList", "entrypointArg", "flags",
	"flag", "flagKey", "flagKeyValue", "stmtWords", "stmtWord", "envArgKey",
	"envArgValue", "copySrc", "copyDest", "asName", "imageName", "saveImageName",
	"targetName", "fullTargetName", "artifactName", "saveFrom", "saveTo", "saveAsLocalTo",
	"workdirPath", "gitURL", "argsList", "arg",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type EarthParser struct {
	*antlr.BaseParser
}

func NewEarthParser(input antlr.TokenStream) *EarthParser {
	this := new(EarthParser)

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
	EarthParserEOF           = antlr.TokenEOF
	EarthParserINDENT        = 1
	EarthParserDEDENT        = 2
	EarthParserTarget        = 3
	EarthParserFROM          = 4
	EarthParserCOPY          = 5
	EarthParserSAVE_ARTIFACT = 6
	EarthParserSAVE_IMAGE    = 7
	EarthParserRUN           = 8
	EarthParserENV           = 9
	EarthParserARG           = 10
	EarthParserBUILD         = 11
	EarthParserWORKDIR       = 12
	EarthParserENTRYPOINT    = 13
	EarthParserGIT_CLONE     = 14
	EarthParserDOCKER_LOAD   = 15
	EarthParserDOCKER_PULL   = 16
	EarthParserCommand       = 17
	EarthParserNL            = 18
	EarthParserWS            = 19
	EarthParserOPEN_BRACKET  = 20
	EarthParserFLAG_ARTIFACT = 21
	EarthParserFLAG_FROM     = 22
	EarthParserFlagKeyValue  = 23
	EarthParserFlagKey       = 24
	EarthParserAtom          = 25
	EarthParserAS            = 26
	EarthParserAS_LOCAL      = 27
	EarthParserCLOSE_BRACKET = 28
	EarthParserCOMMA         = 29
	EarthParserEQUALS        = 30
)

// EarthParser rules.
const (
	EarthParserRULE_earthFile          = 0
	EarthParserRULE_targets            = 1
	EarthParserRULE_target             = 2
	EarthParserRULE_targetHeader       = 3
	EarthParserRULE_stmts              = 4
	EarthParserRULE_stmt               = 5
	EarthParserRULE_fromStmt           = 6
	EarthParserRULE_copyStmt           = 7
	EarthParserRULE_copyArgsFrom       = 8
	EarthParserRULE_copyArgsArtifact   = 9
	EarthParserRULE_copyArgsClassical  = 10
	EarthParserRULE_copySrcs           = 11
	EarthParserRULE_saveStmt           = 12
	EarthParserRULE_saveImage          = 13
	EarthParserRULE_saveArtifact       = 14
	EarthParserRULE_runStmt            = 15
	EarthParserRULE_buildStmt          = 16
	EarthParserRULE_workdirStmt        = 17
	EarthParserRULE_entrypointStmt     = 18
	EarthParserRULE_envStmt            = 19
	EarthParserRULE_argStmt            = 20
	EarthParserRULE_gitCloneStmt       = 21
	EarthParserRULE_dockerLoadStmt     = 22
	EarthParserRULE_dockerPullStmt     = 23
	EarthParserRULE_genericCommand     = 24
	EarthParserRULE_commandName        = 25
	EarthParserRULE_runArgs            = 26
	EarthParserRULE_runArgsList        = 27
	EarthParserRULE_runArg             = 28
	EarthParserRULE_entrypointArgs     = 29
	EarthParserRULE_entrypointArgsList = 30
	EarthParserRULE_entrypointArg      = 31
	EarthParserRULE_flags              = 32
	EarthParserRULE_flag               = 33
	EarthParserRULE_flagKey            = 34
	EarthParserRULE_flagKeyValue       = 35
	EarthParserRULE_stmtWords          = 36
	EarthParserRULE_stmtWord           = 37
	EarthParserRULE_envArgKey          = 38
	EarthParserRULE_envArgValue        = 39
	EarthParserRULE_copySrc            = 40
	EarthParserRULE_copyDest           = 41
	EarthParserRULE_asName             = 42
	EarthParserRULE_imageName          = 43
	EarthParserRULE_saveImageName      = 44
	EarthParserRULE_targetName         = 45
	EarthParserRULE_fullTargetName     = 46
	EarthParserRULE_artifactName       = 47
	EarthParserRULE_saveFrom           = 48
	EarthParserRULE_saveTo             = 49
	EarthParserRULE_saveAsLocalTo      = 50
	EarthParserRULE_workdirPath        = 51
	EarthParserRULE_gitURL             = 52
	EarthParserRULE_argsList           = 53
	EarthParserRULE_arg                = 54
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
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(110)
				p.Match(EarthParserNL)
			}

		}
		p.SetState(115)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<EarthParserFROM)|(1<<EarthParserCOPY)|(1<<EarthParserSAVE_ARTIFACT)|(1<<EarthParserSAVE_IMAGE)|(1<<EarthParserRUN)|(1<<EarthParserENV)|(1<<EarthParserARG)|(1<<EarthParserBUILD)|(1<<EarthParserWORKDIR)|(1<<EarthParserENTRYPOINT)|(1<<EarthParserGIT_CLONE)|(1<<EarthParserDOCKER_LOAD)|(1<<EarthParserDOCKER_PULL)|(1<<EarthParserCommand)|(1<<EarthParserWS))) != 0 {
		{
			p.SetState(116)
			p.Stmts()
		}

	}
	p.SetState(120)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(119)
				p.Match(EarthParserNL)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(122)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}
	p.SetState(125)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserTarget {
		{
			p.SetState(124)
			p.Targets()
		}

	}
	p.SetState(130)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == EarthParserNL {
		{
			p.SetState(127)
			p.Match(EarthParserNL)
		}

		p.SetState(132)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(133)
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

func (s *TargetsContext) AllTarget() []ITargetContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITargetContext)(nil)).Elem())
	var tst = make([]ITargetContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITargetContext)
		}
	}

	return tst
}

func (s *TargetsContext) Target(i int) ITargetContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITargetContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITargetContext)
}

func (s *TargetsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *TargetsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *TargetsContext) AllDEDENT() []antlr.TerminalNode {
	return s.GetTokens(EarthParserDEDENT)
}

func (s *TargetsContext) DEDENT(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserDEDENT, i)
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
		p.SetState(135)
		p.Target()
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(136)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(151)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(140)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(139)
					p.Match(EarthParserNL)
				}

				p.SetState(142)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(144)
				p.Match(EarthParserDEDENT)
			}
			{
				p.SetState(145)
				p.Target()
			}
			p.SetState(147)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(146)
					p.Match(EarthParserWS)
				}

			}

		}
		p.SetState(153)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
	}
	p.SetState(157)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(154)
				p.Match(EarthParserNL)
			}

		}
		p.SetState(159)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}
	p.SetState(161)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserDEDENT {
		{
			p.SetState(160)
			p.Match(EarthParserDEDENT)
		}

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

func (s *TargetContext) INDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserINDENT, 0)
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

func (s *TargetContext) Stmts() IStmtsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
	localctx = NewTargetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, EarthParserRULE_target)
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
		p.SetState(163)
		p.TargetHeader()
	}
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(164)
			p.Match(EarthParserNL)
		}

		p.SetState(167)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(170)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(169)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(172)
		p.Match(EarthParserINDENT)
	}
	p.SetState(174)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(173)
			p.Stmts()
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
	p.EnterRule(localctx, 6, EarthParserRULE_targetHeader)

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
		p.SetState(176)
		p.Match(EarthParserTarget)
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
	p.EnterRule(localctx, 8, EarthParserRULE_stmts)
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
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(178)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(181)
		p.Stmt()
	}
	p.SetState(193)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(183)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(182)
					p.Match(EarthParserNL)
				}

				p.SetState(185)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			p.SetState(188)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(187)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(190)
				p.Stmt()
			}

		}
		p.SetState(195)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext())
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

func (s *StmtContext) FromStmt() IFromStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFromStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFromStmtContext)
}

func (s *StmtContext) CopyStmt() ICopyStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyStmtContext)
}

func (s *StmtContext) SaveStmt() ISaveStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveStmtContext)
}

func (s *StmtContext) RunStmt() IRunStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRunStmtContext)
}

func (s *StmtContext) BuildStmt() IBuildStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBuildStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBuildStmtContext)
}

func (s *StmtContext) WorkdirStmt() IWorkdirStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWorkdirStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWorkdirStmtContext)
}

func (s *StmtContext) EntrypointStmt() IEntrypointStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntrypointStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEntrypointStmtContext)
}

func (s *StmtContext) EnvStmt() IEnvStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEnvStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEnvStmtContext)
}

func (s *StmtContext) ArgStmt() IArgStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgStmtContext)
}

func (s *StmtContext) GitCloneStmt() IGitCloneStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGitCloneStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGitCloneStmtContext)
}

func (s *StmtContext) DockerLoadStmt() IDockerLoadStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDockerLoadStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDockerLoadStmtContext)
}

func (s *StmtContext) DockerPullStmt() IDockerPullStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDockerPullStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDockerPullStmtContext)
}

func (s *StmtContext) GenericCommand() IGenericCommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGenericCommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGenericCommandContext)
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
	p.EnterRule(localctx, 10, EarthParserRULE_stmt)

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

	p.SetState(209)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(196)
			p.FromStmt()
		}

	case EarthParserCOPY:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(197)
			p.CopyStmt()
		}

	case EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(198)
			p.SaveStmt()
		}

	case EarthParserRUN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(199)
			p.RunStmt()
		}

	case EarthParserBUILD:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(200)
			p.BuildStmt()
		}

	case EarthParserWORKDIR:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(201)
			p.WorkdirStmt()
		}

	case EarthParserENTRYPOINT:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(202)
			p.EntrypointStmt()
		}

	case EarthParserENV:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(203)
			p.EnvStmt()
		}

	case EarthParserARG:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(204)
			p.ArgStmt()
		}

	case EarthParserGIT_CLONE:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(205)
			p.GitCloneStmt()
		}

	case EarthParserDOCKER_LOAD:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(206)
			p.DockerLoadStmt()
		}

	case EarthParserDOCKER_PULL:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(207)
			p.DockerPullStmt()
		}

	case EarthParserCommand:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(208)
			p.GenericCommand()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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

func (s *FromStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *FromStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *FromStmtContext) ImageName() IImageNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImageNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImageNameContext)
}

func (s *FromStmtContext) AllFlagKeyValue() []IFlagKeyValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem())
	var tst = make([]IFlagKeyValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagKeyValueContext)
		}
	}

	return tst
}

func (s *FromStmtContext) FlagKeyValue(i int) IFlagKeyValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyValueContext)
}

func (s *FromStmtContext) AS() antlr.TerminalNode {
	return s.GetToken(EarthParserAS, 0)
}

func (s *FromStmtContext) AsName() IAsNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAsNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAsNameContext)
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
	p.EnterRule(localctx, 12, EarthParserRULE_fromStmt)

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
		p.SetState(211)
		p.Match(EarthParserFROM)
	}
	p.SetState(216)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(212)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(213)
				p.FlagKeyValue()
			}

		}
		p.SetState(218)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())
	}
	{
		p.SetState(219)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(220)
		p.ImageName()
	}
	p.SetState(224)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(221)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(222)
			p.Match(EarthParserAS)
		}
		{
			p.SetState(223)
			p.AsName()
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

func (s *CopyStmtContext) CopyArgsFrom() ICopyArgsFromContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyArgsFromContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyArgsFromContext)
}

func (s *CopyStmtContext) CopyArgsArtifact() ICopyArgsArtifactContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyArgsArtifactContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyArgsArtifactContext)
}

func (s *CopyStmtContext) CopyArgsClassical() ICopyArgsClassicalContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyArgsClassicalContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyArgsClassicalContext)
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
	p.EnterRule(localctx, 14, EarthParserRULE_copyStmt)

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
		p.SetState(226)
		p.Match(EarthParserCOPY)
	}
	{
		p.SetState(227)
		p.Match(EarthParserWS)
	}
	p.SetState(231)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFLAG_FROM:
		{
			p.SetState(228)
			p.CopyArgsFrom()
		}

	case EarthParserFLAG_ARTIFACT:
		{
			p.SetState(229)
			p.CopyArgsArtifact()
		}

	case EarthParserOPEN_BRACKET, EarthParserAtom:
		{
			p.SetState(230)
			p.CopyArgsClassical()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ICopyArgsFromContext is an interface to support dynamic dispatch.
type ICopyArgsFromContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCopyArgsFromContext differentiates from other interfaces.
	IsCopyArgsFromContext()
}

type CopyArgsFromContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopyArgsFromContext() *CopyArgsFromContext {
	var p = new(CopyArgsFromContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copyArgsFrom
	return p
}

func (*CopyArgsFromContext) IsCopyArgsFromContext() {}

func NewCopyArgsFromContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopyArgsFromContext {
	var p = new(CopyArgsFromContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copyArgsFrom

	return p
}

func (s *CopyArgsFromContext) GetParser() antlr.Parser { return s.parser }

func (s *CopyArgsFromContext) FLAG_FROM() antlr.TerminalNode {
	return s.GetToken(EarthParserFLAG_FROM, 0)
}

func (s *CopyArgsFromContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *CopyArgsFromContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *CopyArgsFromContext) ImageName() IImageNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImageNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImageNameContext)
}

func (s *CopyArgsFromContext) StmtWord() IStmtWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordContext)
}

func (s *CopyArgsFromContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *CopyArgsFromContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopyArgsFromContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CopyArgsFromContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopyArgsFrom(s)
	}
}

func (s *CopyArgsFromContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopyArgsFrom(s)
	}
}

func (p *EarthParser) CopyArgsFrom() (localctx ICopyArgsFromContext) {
	localctx = NewCopyArgsFromContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, EarthParserRULE_copyArgsFrom)

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
		p.SetState(233)
		p.Match(EarthParserFLAG_FROM)
	}
	{
		p.SetState(234)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(235)
		p.ImageName()
	}
	{
		p.SetState(236)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(237)
		p.StmtWord()
	}
	{
		p.SetState(238)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(239)
		p.StmtWords()
	}

	return localctx
}

// ICopyArgsArtifactContext is an interface to support dynamic dispatch.
type ICopyArgsArtifactContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCopyArgsArtifactContext differentiates from other interfaces.
	IsCopyArgsArtifactContext()
}

type CopyArgsArtifactContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopyArgsArtifactContext() *CopyArgsArtifactContext {
	var p = new(CopyArgsArtifactContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copyArgsArtifact
	return p
}

func (*CopyArgsArtifactContext) IsCopyArgsArtifactContext() {}

func NewCopyArgsArtifactContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopyArgsArtifactContext {
	var p = new(CopyArgsArtifactContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copyArgsArtifact

	return p
}

func (s *CopyArgsArtifactContext) GetParser() antlr.Parser { return s.parser }

func (s *CopyArgsArtifactContext) FLAG_ARTIFACT() antlr.TerminalNode {
	return s.GetToken(EarthParserFLAG_ARTIFACT, 0)
}

func (s *CopyArgsArtifactContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *CopyArgsArtifactContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *CopyArgsArtifactContext) ArtifactName() IArtifactNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArtifactNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArtifactNameContext)
}

func (s *CopyArgsArtifactContext) CopyDest() ICopyDestContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyDestContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyDestContext)
}

func (s *CopyArgsArtifactContext) AllFlagKeyValue() []IFlagKeyValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem())
	var tst = make([]IFlagKeyValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagKeyValueContext)
		}
	}

	return tst
}

func (s *CopyArgsArtifactContext) FlagKeyValue(i int) IFlagKeyValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyValueContext)
}

func (s *CopyArgsArtifactContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopyArgsArtifactContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CopyArgsArtifactContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopyArgsArtifact(s)
	}
}

func (s *CopyArgsArtifactContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopyArgsArtifact(s)
	}
}

func (p *EarthParser) CopyArgsArtifact() (localctx ICopyArgsArtifactContext) {
	localctx = NewCopyArgsArtifactContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, EarthParserRULE_copyArgsArtifact)

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
		p.SetState(241)
		p.Match(EarthParserFLAG_ARTIFACT)
	}
	p.SetState(246)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(242)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(243)
				p.FlagKeyValue()
			}

		}
		p.SetState(248)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())
	}
	{
		p.SetState(249)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(250)
		p.ArtifactName()
	}
	{
		p.SetState(251)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(252)
		p.CopyDest()
	}

	return localctx
}

// ICopyArgsClassicalContext is an interface to support dynamic dispatch.
type ICopyArgsClassicalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCopyArgsClassicalContext differentiates from other interfaces.
	IsCopyArgsClassicalContext()
}

type CopyArgsClassicalContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopyArgsClassicalContext() *CopyArgsClassicalContext {
	var p = new(CopyArgsClassicalContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copyArgsClassical
	return p
}

func (*CopyArgsClassicalContext) IsCopyArgsClassicalContext() {}

func NewCopyArgsClassicalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopyArgsClassicalContext {
	var p = new(CopyArgsClassicalContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copyArgsClassical

	return p
}

func (s *CopyArgsClassicalContext) GetParser() antlr.Parser { return s.parser }

func (s *CopyArgsClassicalContext) CopySrcs() ICopySrcsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopySrcsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopySrcsContext)
}

func (s *CopyArgsClassicalContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *CopyArgsClassicalContext) CopyDest() ICopyDestContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyDestContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyDestContext)
}

func (s *CopyArgsClassicalContext) ArgsList() IArgsListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgsListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgsListContext)
}

func (s *CopyArgsClassicalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopyArgsClassicalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CopyArgsClassicalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopyArgsClassical(s)
	}
}

func (s *CopyArgsClassicalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopyArgsClassical(s)
	}
}

func (p *EarthParser) CopyArgsClassical() (localctx ICopyArgsClassicalContext) {
	localctx = NewCopyArgsClassicalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, EarthParserRULE_copyArgsClassical)

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

	p.SetState(259)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserAtom:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(254)
			p.CopySrcs()
		}
		{
			p.SetState(255)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(256)
			p.CopyDest()
		}

	case EarthParserOPEN_BRACKET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(258)
			p.ArgsList()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ICopySrcsContext is an interface to support dynamic dispatch.
type ICopySrcsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCopySrcsContext differentiates from other interfaces.
	IsCopySrcsContext()
}

type CopySrcsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopySrcsContext() *CopySrcsContext {
	var p = new(CopySrcsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copySrcs
	return p
}

func (*CopySrcsContext) IsCopySrcsContext() {}

func NewCopySrcsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopySrcsContext {
	var p = new(CopySrcsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copySrcs

	return p
}

func (s *CopySrcsContext) GetParser() antlr.Parser { return s.parser }

func (s *CopySrcsContext) AllCopySrc() []ICopySrcContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICopySrcContext)(nil)).Elem())
	var tst = make([]ICopySrcContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICopySrcContext)
		}
	}

	return tst
}

func (s *CopySrcsContext) CopySrc(i int) ICopySrcContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopySrcContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICopySrcContext)
}

func (s *CopySrcsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *CopySrcsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *CopySrcsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopySrcsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CopySrcsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopySrcs(s)
	}
}

func (s *CopySrcsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopySrcs(s)
	}
}

func (p *EarthParser) CopySrcs() (localctx ICopySrcsContext) {
	localctx = NewCopySrcsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, EarthParserRULE_copySrcs)

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
		p.CopySrc()
	}
	p.SetState(266)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(262)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(263)
				p.CopySrc()
			}

		}
		p.SetState(268)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 24, EarthParserRULE_saveStmt)

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
	case EarthParserSAVE_ARTIFACT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(269)
			p.SaveArtifact()
		}

	case EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(270)
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

func (s *SaveImageContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *SaveImageContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *SaveImageContext) AllSaveImageName() []ISaveImageNameContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISaveImageNameContext)(nil)).Elem())
	var tst = make([]ISaveImageNameContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISaveImageNameContext)
		}
	}

	return tst
}

func (s *SaveImageContext) SaveImageName(i int) ISaveImageNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveImageNameContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISaveImageNameContext)
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
	p.EnterRule(localctx, 26, EarthParserRULE_saveImage)

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
		p.Match(EarthParserSAVE_IMAGE)
	}
	p.SetState(278)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(274)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(275)
				p.SaveImageName()
			}

		}
		p.SetState(280)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext())
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

func (s *SaveArtifactContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *SaveArtifactContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *SaveArtifactContext) SaveFrom() ISaveFromContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveFromContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveFromContext)
}

func (s *SaveArtifactContext) SaveTo() ISaveToContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveToContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveToContext)
}

func (s *SaveArtifactContext) AS_LOCAL() antlr.TerminalNode {
	return s.GetToken(EarthParserAS_LOCAL, 0)
}

func (s *SaveArtifactContext) SaveAsLocalTo() ISaveAsLocalToContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveAsLocalToContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveAsLocalToContext)
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
	p.EnterRule(localctx, 28, EarthParserRULE_saveArtifact)

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
		p.SetState(281)
		p.Match(EarthParserSAVE_ARTIFACT)
	}
	{
		p.SetState(282)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(283)
		p.SaveFrom()
	}
	p.SetState(286)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(284)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(285)
			p.SaveTo()
		}

	}
	p.SetState(292)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(288)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(289)
			p.Match(EarthParserAS_LOCAL)
		}
		{
			p.SetState(290)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(291)
			p.SaveAsLocalTo()
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

func (s *RunStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *RunStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *RunStmtContext) RunArgs() IRunArgsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunArgsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRunArgsContext)
}

func (s *RunStmtContext) RunArgsList() IRunArgsListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunArgsListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRunArgsListContext)
}

func (s *RunStmtContext) AllFlag() []IFlagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagContext)(nil)).Elem())
	var tst = make([]IFlagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagContext)
		}
	}

	return tst
}

func (s *RunStmtContext) Flag(i int) IFlagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagContext)
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
	p.EnterRule(localctx, 30, EarthParserRULE_runStmt)

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
		p.SetState(294)
		p.Match(EarthParserRUN)
	}
	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(295)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(296)
				p.Flag()
			}

		}
		p.SetState(301)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())
	}
	{
		p.SetState(302)
		p.Match(EarthParserWS)
	}
	p.SetState(305)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserAtom:
		{
			p.SetState(303)
			p.RunArgs()
		}

	case EarthParserOPEN_BRACKET:
		{
			p.SetState(304)
			p.RunArgsList()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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

func (s *BuildStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *BuildStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *BuildStmtContext) FullTargetName() IFullTargetNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFullTargetNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFullTargetNameContext)
}

func (s *BuildStmtContext) AllFlagKeyValue() []IFlagKeyValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem())
	var tst = make([]IFlagKeyValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagKeyValueContext)
		}
	}

	return tst
}

func (s *BuildStmtContext) FlagKeyValue(i int) IFlagKeyValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyValueContext)
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
	p.EnterRule(localctx, 32, EarthParserRULE_buildStmt)

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
		p.SetState(307)
		p.Match(EarthParserBUILD)
	}
	p.SetState(312)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(308)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(309)
				p.FlagKeyValue()
			}

		}
		p.SetState(314)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext())
	}
	{
		p.SetState(315)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(316)
		p.FullTargetName()
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

func (s *WorkdirStmtContext) WorkdirPath() IWorkdirPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWorkdirPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWorkdirPathContext)
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
	p.EnterRule(localctx, 34, EarthParserRULE_workdirStmt)

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
		p.SetState(318)
		p.Match(EarthParserWORKDIR)
	}
	{
		p.SetState(319)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(320)
		p.WorkdirPath()
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

func (s *EntrypointStmtContext) EntrypointArgs() IEntrypointArgsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntrypointArgsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEntrypointArgsContext)
}

func (s *EntrypointStmtContext) EntrypointArgsList() IEntrypointArgsListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntrypointArgsListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEntrypointArgsListContext)
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
	p.EnterRule(localctx, 36, EarthParserRULE_entrypointStmt)

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
		p.SetState(322)
		p.Match(EarthParserENTRYPOINT)
	}
	{
		p.SetState(323)
		p.Match(EarthParserWS)
	}
	p.SetState(326)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserAtom:
		{
			p.SetState(324)
			p.EntrypointArgs()
		}

	case EarthParserOPEN_BRACKET:
		{
			p.SetState(325)
			p.EntrypointArgsList()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 38, EarthParserRULE_envStmt)
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
		p.SetState(328)
		p.Match(EarthParserENV)
	}
	{
		p.SetState(329)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(330)
		p.EnvArgKey()
	}
	p.SetState(335)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext()) == 1 {
		p.SetState(332)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(331)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(334)
			p.Match(EarthParserEQUALS)
		}

	}
	p.SetState(341)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext()) == 1 {
		p.SetState(338)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(337)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(340)
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
	p.EnterRule(localctx, 40, EarthParserRULE_argStmt)
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
		p.Match(EarthParserARG)
	}
	{
		p.SetState(344)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(345)
		p.EnvArgKey()
	}
	p.SetState(357)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) == 1 {
		p.SetState(347)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(346)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(349)
			p.Match(EarthParserEQUALS)
		}

		p.SetState(355)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext()) == 1 {
			p.SetState(352)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(351)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(354)
				p.EnvArgValue()
			}

		}

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

func (s *GitCloneStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *GitCloneStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *GitCloneStmtContext) GitURL() IGitURLContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGitURLContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGitURLContext)
}

func (s *GitCloneStmtContext) CopyDest() ICopyDestContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICopyDestContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICopyDestContext)
}

func (s *GitCloneStmtContext) AllFlagKeyValue() []IFlagKeyValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem())
	var tst = make([]IFlagKeyValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagKeyValueContext)
		}
	}

	return tst
}

func (s *GitCloneStmtContext) FlagKeyValue(i int) IFlagKeyValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyValueContext)
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
	p.EnterRule(localctx, 42, EarthParserRULE_gitCloneStmt)

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
		p.SetState(359)
		p.Match(EarthParserGIT_CLONE)
	}
	p.SetState(364)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(360)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(361)
				p.FlagKeyValue()
			}

		}
		p.SetState(366)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())
	}
	{
		p.SetState(367)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(368)
		p.GitURL()
	}
	{
		p.SetState(369)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(370)
		p.CopyDest()
	}

	return localctx
}

// IDockerLoadStmtContext is an interface to support dynamic dispatch.
type IDockerLoadStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDockerLoadStmtContext differentiates from other interfaces.
	IsDockerLoadStmtContext()
}

type DockerLoadStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDockerLoadStmtContext() *DockerLoadStmtContext {
	var p = new(DockerLoadStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_dockerLoadStmt
	return p
}

func (*DockerLoadStmtContext) IsDockerLoadStmtContext() {}

func NewDockerLoadStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DockerLoadStmtContext {
	var p = new(DockerLoadStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_dockerLoadStmt

	return p
}

func (s *DockerLoadStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *DockerLoadStmtContext) DOCKER_LOAD() antlr.TerminalNode {
	return s.GetToken(EarthParserDOCKER_LOAD, 0)
}

func (s *DockerLoadStmtContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *DockerLoadStmtContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *DockerLoadStmtContext) FullTargetName() IFullTargetNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFullTargetNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFullTargetNameContext)
}

func (s *DockerLoadStmtContext) AS() antlr.TerminalNode {
	return s.GetToken(EarthParserAS, 0)
}

func (s *DockerLoadStmtContext) ImageName() IImageNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImageNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImageNameContext)
}

func (s *DockerLoadStmtContext) AllFlagKeyValue() []IFlagKeyValueContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem())
	var tst = make([]IFlagKeyValueContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagKeyValueContext)
		}
	}

	return tst
}

func (s *DockerLoadStmtContext) FlagKeyValue(i int) IFlagKeyValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyValueContext)
}

func (s *DockerLoadStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DockerLoadStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DockerLoadStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterDockerLoadStmt(s)
	}
}

func (s *DockerLoadStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitDockerLoadStmt(s)
	}
}

func (p *EarthParser) DockerLoadStmt() (localctx IDockerLoadStmtContext) {
	localctx = NewDockerLoadStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, EarthParserRULE_dockerLoadStmt)

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
		p.SetState(372)
		p.Match(EarthParserDOCKER_LOAD)
	}
	p.SetState(377)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(373)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(374)
				p.FlagKeyValue()
			}

		}
		p.SetState(379)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())
	}
	{
		p.SetState(380)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(381)
		p.FullTargetName()
	}
	{
		p.SetState(382)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(383)
		p.Match(EarthParserAS)
	}
	{
		p.SetState(384)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(385)
		p.ImageName()
	}

	return localctx
}

// IDockerPullStmtContext is an interface to support dynamic dispatch.
type IDockerPullStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDockerPullStmtContext differentiates from other interfaces.
	IsDockerPullStmtContext()
}

type DockerPullStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDockerPullStmtContext() *DockerPullStmtContext {
	var p = new(DockerPullStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_dockerPullStmt
	return p
}

func (*DockerPullStmtContext) IsDockerPullStmtContext() {}

func NewDockerPullStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DockerPullStmtContext {
	var p = new(DockerPullStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_dockerPullStmt

	return p
}

func (s *DockerPullStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *DockerPullStmtContext) DOCKER_PULL() antlr.TerminalNode {
	return s.GetToken(EarthParserDOCKER_PULL, 0)
}

func (s *DockerPullStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *DockerPullStmtContext) ImageName() IImageNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImageNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImageNameContext)
}

func (s *DockerPullStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DockerPullStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DockerPullStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterDockerPullStmt(s)
	}
}

func (s *DockerPullStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitDockerPullStmt(s)
	}
}

func (p *EarthParser) DockerPullStmt() (localctx IDockerPullStmtContext) {
	localctx = NewDockerPullStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, EarthParserRULE_dockerPullStmt)

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
		p.Match(EarthParserDOCKER_PULL)
	}
	{
		p.SetState(388)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(389)
		p.ImageName()
	}

	return localctx
}

// IGenericCommandContext is an interface to support dynamic dispatch.
type IGenericCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGenericCommandContext differentiates from other interfaces.
	IsGenericCommandContext()
}

type GenericCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGenericCommandContext() *GenericCommandContext {
	var p = new(GenericCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_genericCommand
	return p
}

func (*GenericCommandContext) IsGenericCommandContext() {}

func NewGenericCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericCommandContext {
	var p = new(GenericCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_genericCommand

	return p
}

func (s *GenericCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *GenericCommandContext) CommandName() ICommandNameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandNameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommandNameContext)
}

func (s *GenericCommandContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *GenericCommandContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *GenericCommandContext) Flags() IFlagsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFlagsContext)
}

func (s *GenericCommandContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *GenericCommandContext) ArgsList() IArgsListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgsListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IArgsListContext)
}

func (s *GenericCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GenericCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GenericCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterGenericCommand(s)
	}
}

func (s *GenericCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitGenericCommand(s)
	}
}

func (p *EarthParser) GenericCommand() (localctx IGenericCommandContext) {
	localctx = NewGenericCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, EarthParserRULE_genericCommand)

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
		p.SetState(391)
		p.CommandName()
	}
	p.SetState(394)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(392)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(393)
			p.Flags()
		}

	}
	p.SetState(400)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(396)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(397)
			p.StmtWords()
		}

	} else if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) == 2 {
		{
			p.SetState(398)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(399)
			p.ArgsList()
		}

	}

	return localctx
}

// ICommandNameContext is an interface to support dynamic dispatch.
type ICommandNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommandNameContext differentiates from other interfaces.
	IsCommandNameContext()
}

type CommandNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandNameContext() *CommandNameContext {
	var p = new(CommandNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_commandName
	return p
}

func (*CommandNameContext) IsCommandNameContext() {}

func NewCommandNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandNameContext {
	var p = new(CommandNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_commandName

	return p
}

func (s *CommandNameContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandNameContext) Command() antlr.TerminalNode {
	return s.GetToken(EarthParserCommand, 0)
}

func (s *CommandNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCommandName(s)
	}
}

func (s *CommandNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCommandName(s)
	}
}

func (p *EarthParser) CommandName() (localctx ICommandNameContext) {
	localctx = NewCommandNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, EarthParserRULE_commandName)

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
		p.Match(EarthParserCommand)
	}

	return localctx
}

// IRunArgsContext is an interface to support dynamic dispatch.
type IRunArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRunArgsContext differentiates from other interfaces.
	IsRunArgsContext()
}

type RunArgsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRunArgsContext() *RunArgsContext {
	var p = new(RunArgsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_runArgs
	return p
}

func (*RunArgsContext) IsRunArgsContext() {}

func NewRunArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RunArgsContext {
	var p = new(RunArgsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_runArgs

	return p
}

func (s *RunArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *RunArgsContext) AllRunArg() []IRunArgContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IRunArgContext)(nil)).Elem())
	var tst = make([]IRunArgContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IRunArgContext)
		}
	}

	return tst
}

func (s *RunArgsContext) RunArg(i int) IRunArgContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunArgContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IRunArgContext)
}

func (s *RunArgsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *RunArgsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *RunArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RunArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RunArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterRunArgs(s)
	}
}

func (s *RunArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitRunArgs(s)
	}
}

func (p *EarthParser) RunArgs() (localctx IRunArgsContext) {
	localctx = NewRunArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, EarthParserRULE_runArgs)

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
		p.SetState(404)
		p.RunArg()
	}
	p.SetState(409)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(405)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(406)
				p.RunArg()
			}

		}
		p.SetState(411)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext())
	}

	return localctx
}

// IRunArgsListContext is an interface to support dynamic dispatch.
type IRunArgsListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRunArgsListContext differentiates from other interfaces.
	IsRunArgsListContext()
}

type RunArgsListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRunArgsListContext() *RunArgsListContext {
	var p = new(RunArgsListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_runArgsList
	return p
}

func (*RunArgsListContext) IsRunArgsListContext() {}

func NewRunArgsListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RunArgsListContext {
	var p = new(RunArgsListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_runArgsList

	return p
}

func (s *RunArgsListContext) GetParser() antlr.Parser { return s.parser }

func (s *RunArgsListContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserOPEN_BRACKET, 0)
}

func (s *RunArgsListContext) AllRunArg() []IRunArgContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IRunArgContext)(nil)).Elem())
	var tst = make([]IRunArgContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IRunArgContext)
		}
	}

	return tst
}

func (s *RunArgsListContext) RunArg(i int) IRunArgContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRunArgContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IRunArgContext)
}

func (s *RunArgsListContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserCLOSE_BRACKET, 0)
}

func (s *RunArgsListContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *RunArgsListContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *RunArgsListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EarthParserCOMMA)
}

func (s *RunArgsListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserCOMMA, i)
}

func (s *RunArgsListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RunArgsListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RunArgsListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterRunArgsList(s)
	}
}

func (s *RunArgsListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitRunArgsList(s)
	}
}

func (p *EarthParser) RunArgsList() (localctx IRunArgsListContext) {
	localctx = NewRunArgsListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, EarthParserRULE_runArgsList)
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
		p.SetState(412)
		p.Match(EarthParserOPEN_BRACKET)
	}
	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(413)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(416)
		p.RunArg()
	}
	p.SetState(425)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(418)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(417)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(420)
				p.Match(EarthParserCOMMA)
			}
			p.SetState(422)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(421)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(424)
				p.RunArg()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(427)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 49, p.GetParserRuleContext())
	}
	p.SetState(430)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(429)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(432)
		p.Match(EarthParserCLOSE_BRACKET)
	}

	return localctx
}

// IRunArgContext is an interface to support dynamic dispatch.
type IRunArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRunArgContext differentiates from other interfaces.
	IsRunArgContext()
}

type RunArgContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRunArgContext() *RunArgContext {
	var p = new(RunArgContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_runArg
	return p
}

func (*RunArgContext) IsRunArgContext() {}

func NewRunArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RunArgContext {
	var p = new(RunArgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_runArg

	return p
}

func (s *RunArgContext) GetParser() antlr.Parser { return s.parser }

func (s *RunArgContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *RunArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RunArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RunArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterRunArg(s)
	}
}

func (s *RunArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitRunArg(s)
	}
}

func (p *EarthParser) RunArg() (localctx IRunArgContext) {
	localctx = NewRunArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, EarthParserRULE_runArg)

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
		p.SetState(434)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IEntrypointArgsContext is an interface to support dynamic dispatch.
type IEntrypointArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEntrypointArgsContext differentiates from other interfaces.
	IsEntrypointArgsContext()
}

type EntrypointArgsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEntrypointArgsContext() *EntrypointArgsContext {
	var p = new(EntrypointArgsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_entrypointArgs
	return p
}

func (*EntrypointArgsContext) IsEntrypointArgsContext() {}

func NewEntrypointArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EntrypointArgsContext {
	var p = new(EntrypointArgsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_entrypointArgs

	return p
}

func (s *EntrypointArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *EntrypointArgsContext) AllEntrypointArg() []IEntrypointArgContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEntrypointArgContext)(nil)).Elem())
	var tst = make([]IEntrypointArgContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEntrypointArgContext)
		}
	}

	return tst
}

func (s *EntrypointArgsContext) EntrypointArg(i int) IEntrypointArgContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntrypointArgContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEntrypointArgContext)
}

func (s *EntrypointArgsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *EntrypointArgsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *EntrypointArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EntrypointArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EntrypointArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEntrypointArgs(s)
	}
}

func (s *EntrypointArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEntrypointArgs(s)
	}
}

func (p *EarthParser) EntrypointArgs() (localctx IEntrypointArgsContext) {
	localctx = NewEntrypointArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, EarthParserRULE_entrypointArgs)

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
		p.SetState(436)
		p.EntrypointArg()
	}
	p.SetState(441)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(437)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(438)
				p.EntrypointArg()
			}

		}
		p.SetState(443)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext())
	}

	return localctx
}

// IEntrypointArgsListContext is an interface to support dynamic dispatch.
type IEntrypointArgsListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEntrypointArgsListContext differentiates from other interfaces.
	IsEntrypointArgsListContext()
}

type EntrypointArgsListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEntrypointArgsListContext() *EntrypointArgsListContext {
	var p = new(EntrypointArgsListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_entrypointArgsList
	return p
}

func (*EntrypointArgsListContext) IsEntrypointArgsListContext() {}

func NewEntrypointArgsListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EntrypointArgsListContext {
	var p = new(EntrypointArgsListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_entrypointArgsList

	return p
}

func (s *EntrypointArgsListContext) GetParser() antlr.Parser { return s.parser }

func (s *EntrypointArgsListContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserOPEN_BRACKET, 0)
}

func (s *EntrypointArgsListContext) AllEntrypointArg() []IEntrypointArgContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEntrypointArgContext)(nil)).Elem())
	var tst = make([]IEntrypointArgContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEntrypointArgContext)
		}
	}

	return tst
}

func (s *EntrypointArgsListContext) EntrypointArg(i int) IEntrypointArgContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEntrypointArgContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEntrypointArgContext)
}

func (s *EntrypointArgsListContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserCLOSE_BRACKET, 0)
}

func (s *EntrypointArgsListContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *EntrypointArgsListContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *EntrypointArgsListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EarthParserCOMMA)
}

func (s *EntrypointArgsListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserCOMMA, i)
}

func (s *EntrypointArgsListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EntrypointArgsListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EntrypointArgsListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEntrypointArgsList(s)
	}
}

func (s *EntrypointArgsListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEntrypointArgsList(s)
	}
}

func (p *EarthParser) EntrypointArgsList() (localctx IEntrypointArgsListContext) {
	localctx = NewEntrypointArgsListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, EarthParserRULE_entrypointArgsList)
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
		p.SetState(444)
		p.Match(EarthParserOPEN_BRACKET)
	}
	p.SetState(446)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(445)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(448)
		p.EntrypointArg()
	}
	p.SetState(457)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(450)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(449)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(452)
				p.Match(EarthParserCOMMA)
			}
			p.SetState(454)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(453)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(456)
				p.EntrypointArg()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(459)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext())
	}
	p.SetState(462)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(461)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(464)
		p.Match(EarthParserCLOSE_BRACKET)
	}

	return localctx
}

// IEntrypointArgContext is an interface to support dynamic dispatch.
type IEntrypointArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEntrypointArgContext differentiates from other interfaces.
	IsEntrypointArgContext()
}

type EntrypointArgContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEntrypointArgContext() *EntrypointArgContext {
	var p = new(EntrypointArgContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_entrypointArg
	return p
}

func (*EntrypointArgContext) IsEntrypointArgContext() {}

func NewEntrypointArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EntrypointArgContext {
	var p = new(EntrypointArgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_entrypointArg

	return p
}

func (s *EntrypointArgContext) GetParser() antlr.Parser { return s.parser }

func (s *EntrypointArgContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *EntrypointArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EntrypointArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EntrypointArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterEntrypointArg(s)
	}
}

func (s *EntrypointArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitEntrypointArg(s)
	}
}

func (p *EarthParser) EntrypointArg() (localctx IEntrypointArgContext) {
	localctx = NewEntrypointArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, EarthParserRULE_entrypointArg)

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
		p.SetState(466)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IFlagsContext is an interface to support dynamic dispatch.
type IFlagsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFlagsContext differentiates from other interfaces.
	IsFlagsContext()
}

type FlagsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFlagsContext() *FlagsContext {
	var p = new(FlagsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_flags
	return p
}

func (*FlagsContext) IsFlagsContext() {}

func NewFlagsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FlagsContext {
	var p = new(FlagsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_flags

	return p
}

func (s *FlagsContext) GetParser() antlr.Parser { return s.parser }

func (s *FlagsContext) AllFlag() []IFlagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFlagContext)(nil)).Elem())
	var tst = make([]IFlagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFlagContext)
		}
	}

	return tst
}

func (s *FlagsContext) Flag(i int) IFlagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFlagContext)
}

func (s *FlagsContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *FlagsContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *FlagsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FlagsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FlagsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFlags(s)
	}
}

func (s *FlagsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFlags(s)
	}
}

func (p *EarthParser) Flags() (localctx IFlagsContext) {
	localctx = NewFlagsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, EarthParserRULE_flags)
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
		p.SetState(468)
		p.Flag()
	}
	p.SetState(475)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 58, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(470)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(469)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(472)
				p.Flag()
			}

		}
		p.SetState(477)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 58, p.GetParserRuleContext())
	}

	return localctx
}

// IFlagContext is an interface to support dynamic dispatch.
type IFlagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFlagContext differentiates from other interfaces.
	IsFlagContext()
}

type FlagContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFlagContext() *FlagContext {
	var p = new(FlagContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_flag
	return p
}

func (*FlagContext) IsFlagContext() {}

func NewFlagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FlagContext {
	var p = new(FlagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_flag

	return p
}

func (s *FlagContext) GetParser() antlr.Parser { return s.parser }

func (s *FlagContext) FlagKey() IFlagKeyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyContext)
}

func (s *FlagContext) FlagKeyValue() IFlagKeyValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFlagKeyValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFlagKeyValueContext)
}

func (s *FlagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FlagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FlagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFlag(s)
	}
}

func (s *FlagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFlag(s)
	}
}

func (p *EarthParser) Flag() (localctx IFlagContext) {
	localctx = NewFlagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, EarthParserRULE_flag)

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

	p.SetState(480)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFlagKey:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(478)
			p.FlagKey()
		}

	case EarthParserFlagKeyValue:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(479)
			p.FlagKeyValue()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFlagKeyContext is an interface to support dynamic dispatch.
type IFlagKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFlagKeyContext differentiates from other interfaces.
	IsFlagKeyContext()
}

type FlagKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFlagKeyContext() *FlagKeyContext {
	var p = new(FlagKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_flagKey
	return p
}

func (*FlagKeyContext) IsFlagKeyContext() {}

func NewFlagKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FlagKeyContext {
	var p = new(FlagKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_flagKey

	return p
}

func (s *FlagKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *FlagKeyContext) FlagKey() antlr.TerminalNode {
	return s.GetToken(EarthParserFlagKey, 0)
}

func (s *FlagKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FlagKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FlagKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFlagKey(s)
	}
}

func (s *FlagKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFlagKey(s)
	}
}

func (p *EarthParser) FlagKey() (localctx IFlagKeyContext) {
	localctx = NewFlagKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, EarthParserRULE_flagKey)

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
		p.Match(EarthParserFlagKey)
	}

	return localctx
}

// IFlagKeyValueContext is an interface to support dynamic dispatch.
type IFlagKeyValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFlagKeyValueContext differentiates from other interfaces.
	IsFlagKeyValueContext()
}

type FlagKeyValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFlagKeyValueContext() *FlagKeyValueContext {
	var p = new(FlagKeyValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_flagKeyValue
	return p
}

func (*FlagKeyValueContext) IsFlagKeyValueContext() {}

func NewFlagKeyValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FlagKeyValueContext {
	var p = new(FlagKeyValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_flagKeyValue

	return p
}

func (s *FlagKeyValueContext) GetParser() antlr.Parser { return s.parser }

func (s *FlagKeyValueContext) FlagKeyValue() antlr.TerminalNode {
	return s.GetToken(EarthParserFlagKeyValue, 0)
}

func (s *FlagKeyValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FlagKeyValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FlagKeyValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFlagKeyValue(s)
	}
}

func (s *FlagKeyValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFlagKeyValue(s)
	}
}

func (p *EarthParser) FlagKeyValue() (localctx IFlagKeyValueContext) {
	localctx = NewFlagKeyValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, EarthParserRULE_flagKeyValue)

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
		p.SetState(484)
		p.Match(EarthParserFlagKeyValue)
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
	p.EnterRule(localctx, 72, EarthParserRULE_stmtWords)
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
		p.SetState(486)
		p.StmtWord()
	}
	p.SetState(493)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 61, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(488)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(487)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(490)
				p.StmtWord()
			}

		}
		p.SetState(495)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 61, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 74, EarthParserRULE_stmtWord)

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
		p.SetState(496)
		p.Match(EarthParserAtom)
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
	p.EnterRule(localctx, 76, EarthParserRULE_envArgKey)

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
		p.SetState(498)
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
	p.EnterRule(localctx, 78, EarthParserRULE_envArgValue)
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
		p.SetState(500)
		p.Match(EarthParserAtom)
	}
	p.SetState(507)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 63, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(502)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(501)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(504)
				p.Match(EarthParserAtom)
			}

		}
		p.SetState(509)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 63, p.GetParserRuleContext())
	}

	return localctx
}

// ICopySrcContext is an interface to support dynamic dispatch.
type ICopySrcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCopySrcContext differentiates from other interfaces.
	IsCopySrcContext()
}

type CopySrcContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopySrcContext() *CopySrcContext {
	var p = new(CopySrcContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copySrc
	return p
}

func (*CopySrcContext) IsCopySrcContext() {}

func NewCopySrcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopySrcContext {
	var p = new(CopySrcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copySrc

	return p
}

func (s *CopySrcContext) GetParser() antlr.Parser { return s.parser }

func (s *CopySrcContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *CopySrcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopySrcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CopySrcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopySrc(s)
	}
}

func (s *CopySrcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopySrc(s)
	}
}

func (p *EarthParser) CopySrc() (localctx ICopySrcContext) {
	localctx = NewCopySrcContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, EarthParserRULE_copySrc)

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
		p.SetState(510)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ICopyDestContext is an interface to support dynamic dispatch.
type ICopyDestContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCopyDestContext differentiates from other interfaces.
	IsCopyDestContext()
}

type CopyDestContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCopyDestContext() *CopyDestContext {
	var p = new(CopyDestContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_copyDest
	return p
}

func (*CopyDestContext) IsCopyDestContext() {}

func NewCopyDestContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CopyDestContext {
	var p = new(CopyDestContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_copyDest

	return p
}

func (s *CopyDestContext) GetParser() antlr.Parser { return s.parser }

func (s *CopyDestContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *CopyDestContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CopyDestContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CopyDestContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterCopyDest(s)
	}
}

func (s *CopyDestContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitCopyDest(s)
	}
}

func (p *EarthParser) CopyDest() (localctx ICopyDestContext) {
	localctx = NewCopyDestContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, EarthParserRULE_copyDest)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IAsNameContext is an interface to support dynamic dispatch.
type IAsNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAsNameContext differentiates from other interfaces.
	IsAsNameContext()
}

type AsNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAsNameContext() *AsNameContext {
	var p = new(AsNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_asName
	return p
}

func (*AsNameContext) IsAsNameContext() {}

func NewAsNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AsNameContext {
	var p = new(AsNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_asName

	return p
}

func (s *AsNameContext) GetParser() antlr.Parser { return s.parser }

func (s *AsNameContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *AsNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AsNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AsNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterAsName(s)
	}
}

func (s *AsNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitAsName(s)
	}
}

func (p *EarthParser) AsName() (localctx IAsNameContext) {
	localctx = NewAsNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, EarthParserRULE_asName)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IImageNameContext is an interface to support dynamic dispatch.
type IImageNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImageNameContext differentiates from other interfaces.
	IsImageNameContext()
}

type ImageNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImageNameContext() *ImageNameContext {
	var p = new(ImageNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_imageName
	return p
}

func (*ImageNameContext) IsImageNameContext() {}

func NewImageNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImageNameContext {
	var p = new(ImageNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_imageName

	return p
}

func (s *ImageNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ImageNameContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *ImageNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImageNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImageNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterImageName(s)
	}
}

func (s *ImageNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitImageName(s)
	}
}

func (p *EarthParser) ImageName() (localctx IImageNameContext) {
	localctx = NewImageNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, EarthParserRULE_imageName)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ISaveImageNameContext is an interface to support dynamic dispatch.
type ISaveImageNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSaveImageNameContext differentiates from other interfaces.
	IsSaveImageNameContext()
}

type SaveImageNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveImageNameContext() *SaveImageNameContext {
	var p = new(SaveImageNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveImageName
	return p
}

func (*SaveImageNameContext) IsSaveImageNameContext() {}

func NewSaveImageNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveImageNameContext {
	var p = new(SaveImageNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveImageName

	return p
}

func (s *SaveImageNameContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveImageNameContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *SaveImageNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveImageNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SaveImageNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveImageName(s)
	}
}

func (s *SaveImageNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveImageName(s)
	}
}

func (p *EarthParser) SaveImageName() (localctx ISaveImageNameContext) {
	localctx = NewSaveImageNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, EarthParserRULE_saveImageName)

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
		p.SetState(518)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ITargetNameContext is an interface to support dynamic dispatch.
type ITargetNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTargetNameContext differentiates from other interfaces.
	IsTargetNameContext()
}

type TargetNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetNameContext() *TargetNameContext {
	var p = new(TargetNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_targetName
	return p
}

func (*TargetNameContext) IsTargetNameContext() {}

func NewTargetNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetNameContext {
	var p = new(TargetNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_targetName

	return p
}

func (s *TargetNameContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetNameContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *TargetNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TargetNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterTargetName(s)
	}
}

func (s *TargetNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitTargetName(s)
	}
}

func (p *EarthParser) TargetName() (localctx ITargetNameContext) {
	localctx = NewTargetNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, EarthParserRULE_targetName)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IFullTargetNameContext is an interface to support dynamic dispatch.
type IFullTargetNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFullTargetNameContext differentiates from other interfaces.
	IsFullTargetNameContext()
}

type FullTargetNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFullTargetNameContext() *FullTargetNameContext {
	var p = new(FullTargetNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_fullTargetName
	return p
}

func (*FullTargetNameContext) IsFullTargetNameContext() {}

func NewFullTargetNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FullTargetNameContext {
	var p = new(FullTargetNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_fullTargetName

	return p
}

func (s *FullTargetNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FullTargetNameContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *FullTargetNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FullTargetNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FullTargetNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterFullTargetName(s)
	}
}

func (s *FullTargetNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitFullTargetName(s)
	}
}

func (p *EarthParser) FullTargetName() (localctx IFullTargetNameContext) {
	localctx = NewFullTargetNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, EarthParserRULE_fullTargetName)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IArtifactNameContext is an interface to support dynamic dispatch.
type IArtifactNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArtifactNameContext differentiates from other interfaces.
	IsArtifactNameContext()
}

type ArtifactNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArtifactNameContext() *ArtifactNameContext {
	var p = new(ArtifactNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_artifactName
	return p
}

func (*ArtifactNameContext) IsArtifactNameContext() {}

func NewArtifactNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArtifactNameContext {
	var p = new(ArtifactNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_artifactName

	return p
}

func (s *ArtifactNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ArtifactNameContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *ArtifactNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArtifactNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArtifactNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterArtifactName(s)
	}
}

func (s *ArtifactNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitArtifactName(s)
	}
}

func (p *EarthParser) ArtifactName() (localctx IArtifactNameContext) {
	localctx = NewArtifactNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, EarthParserRULE_artifactName)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ISaveFromContext is an interface to support dynamic dispatch.
type ISaveFromContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSaveFromContext differentiates from other interfaces.
	IsSaveFromContext()
}

type SaveFromContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveFromContext() *SaveFromContext {
	var p = new(SaveFromContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveFrom
	return p
}

func (*SaveFromContext) IsSaveFromContext() {}

func NewSaveFromContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveFromContext {
	var p = new(SaveFromContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveFrom

	return p
}

func (s *SaveFromContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveFromContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *SaveFromContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveFromContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SaveFromContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveFrom(s)
	}
}

func (s *SaveFromContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveFrom(s)
	}
}

func (p *EarthParser) SaveFrom() (localctx ISaveFromContext) {
	localctx = NewSaveFromContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, EarthParserRULE_saveFrom)

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
		p.SetState(526)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ISaveToContext is an interface to support dynamic dispatch.
type ISaveToContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSaveToContext differentiates from other interfaces.
	IsSaveToContext()
}

type SaveToContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveToContext() *SaveToContext {
	var p = new(SaveToContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveTo
	return p
}

func (*SaveToContext) IsSaveToContext() {}

func NewSaveToContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveToContext {
	var p = new(SaveToContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveTo

	return p
}

func (s *SaveToContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveToContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *SaveToContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveToContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SaveToContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveTo(s)
	}
}

func (s *SaveToContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveTo(s)
	}
}

func (p *EarthParser) SaveTo() (localctx ISaveToContext) {
	localctx = NewSaveToContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, EarthParserRULE_saveTo)

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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// ISaveAsLocalToContext is an interface to support dynamic dispatch.
type ISaveAsLocalToContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSaveAsLocalToContext differentiates from other interfaces.
	IsSaveAsLocalToContext()
}

type SaveAsLocalToContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveAsLocalToContext() *SaveAsLocalToContext {
	var p = new(SaveAsLocalToContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_saveAsLocalTo
	return p
}

func (*SaveAsLocalToContext) IsSaveAsLocalToContext() {}

func NewSaveAsLocalToContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveAsLocalToContext {
	var p = new(SaveAsLocalToContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_saveAsLocalTo

	return p
}

func (s *SaveAsLocalToContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveAsLocalToContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *SaveAsLocalToContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveAsLocalToContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SaveAsLocalToContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterSaveAsLocalTo(s)
	}
}

func (s *SaveAsLocalToContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitSaveAsLocalTo(s)
	}
}

func (p *EarthParser) SaveAsLocalTo() (localctx ISaveAsLocalToContext) {
	localctx = NewSaveAsLocalToContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, EarthParserRULE_saveAsLocalTo)

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
		p.SetState(530)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IWorkdirPathContext is an interface to support dynamic dispatch.
type IWorkdirPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWorkdirPathContext differentiates from other interfaces.
	IsWorkdirPathContext()
}

type WorkdirPathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWorkdirPathContext() *WorkdirPathContext {
	var p = new(WorkdirPathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_workdirPath
	return p
}

func (*WorkdirPathContext) IsWorkdirPathContext() {}

func NewWorkdirPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WorkdirPathContext {
	var p = new(WorkdirPathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_workdirPath

	return p
}

func (s *WorkdirPathContext) GetParser() antlr.Parser { return s.parser }

func (s *WorkdirPathContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *WorkdirPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WorkdirPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WorkdirPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterWorkdirPath(s)
	}
}

func (s *WorkdirPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitWorkdirPath(s)
	}
}

func (p *EarthParser) WorkdirPath() (localctx IWorkdirPathContext) {
	localctx = NewWorkdirPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, EarthParserRULE_workdirPath)

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
		p.SetState(532)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IGitURLContext is an interface to support dynamic dispatch.
type IGitURLContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGitURLContext differentiates from other interfaces.
	IsGitURLContext()
}

type GitURLContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGitURLContext() *GitURLContext {
	var p = new(GitURLContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_gitURL
	return p
}

func (*GitURLContext) IsGitURLContext() {}

func NewGitURLContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GitURLContext {
	var p = new(GitURLContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_gitURL

	return p
}

func (s *GitURLContext) GetParser() antlr.Parser { return s.parser }

func (s *GitURLContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *GitURLContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GitURLContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GitURLContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterGitURL(s)
	}
}

func (s *GitURLContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitGitURL(s)
	}
}

func (p *EarthParser) GitURL() (localctx IGitURLContext) {
	localctx = NewGitURLContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, EarthParserRULE_gitURL)

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
		p.SetState(534)
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IArgsListContext is an interface to support dynamic dispatch.
type IArgsListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArgsListContext differentiates from other interfaces.
	IsArgsListContext()
}

type ArgsListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgsListContext() *ArgsListContext {
	var p = new(ArgsListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_argsList
	return p
}

func (*ArgsListContext) IsArgsListContext() {}

func NewArgsListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgsListContext {
	var p = new(ArgsListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_argsList

	return p
}

func (s *ArgsListContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgsListContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserOPEN_BRACKET, 0)
}

func (s *ArgsListContext) AllArg() []IArgContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IArgContext)(nil)).Elem())
	var tst = make([]IArgContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IArgContext)
		}
	}

	return tst
}

func (s *ArgsListContext) Arg(i int) IArgContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IArgContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IArgContext)
}

func (s *ArgsListContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserCLOSE_BRACKET, 0)
}

func (s *ArgsListContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *ArgsListContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *ArgsListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EarthParserCOMMA)
}

func (s *ArgsListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserCOMMA, i)
}

func (s *ArgsListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgsListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgsListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterArgsList(s)
	}
}

func (s *ArgsListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitArgsList(s)
	}
}

func (p *EarthParser) ArgsList() (localctx IArgsListContext) {
	localctx = NewArgsListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, EarthParserRULE_argsList)
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
		p.SetState(536)
		p.Match(EarthParserOPEN_BRACKET)
	}
	p.SetState(538)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(537)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(540)
		p.Arg()
	}
	p.SetState(549)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(542)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(541)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(544)
				p.Match(EarthParserCOMMA)
			}
			p.SetState(546)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(545)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(548)
				p.Arg()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(551)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 67, p.GetParserRuleContext())
	}
	p.SetState(554)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(553)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(556)
		p.Match(EarthParserCLOSE_BRACKET)
	}

	return localctx
}

// IArgContext is an interface to support dynamic dispatch.
type IArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArgContext differentiates from other interfaces.
	IsArgContext()
}

type ArgContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgContext() *ArgContext {
	var p = new(ArgContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_arg
	return p
}

func (*ArgContext) IsArgContext() {}

func NewArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgContext {
	var p = new(ArgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_arg

	return p
}

func (s *ArgContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *ArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterArg(s)
	}
}

func (s *ArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitArg(s)
	}
}

func (p *EarthParser) Arg() (localctx IArgContext) {
	localctx = NewArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, EarthParserRULE_arg)

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
		p.SetState(558)
		p.Match(EarthParserAtom)
	}

	return localctx
}
