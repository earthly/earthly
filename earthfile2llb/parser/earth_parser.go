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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 30, 443,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44, 9, 44,
	3, 2, 7, 2, 90, 10, 2, 12, 2, 14, 2, 93, 11, 2, 3, 2, 5, 2, 96, 10, 2,
	3, 2, 6, 2, 99, 10, 2, 13, 2, 14, 2, 100, 3, 2, 5, 2, 104, 10, 2, 3, 2,
	7, 2, 107, 10, 2, 12, 2, 14, 2, 110, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 5,
	3, 116, 10, 3, 3, 3, 6, 3, 119, 10, 3, 13, 3, 14, 3, 120, 3, 3, 3, 3, 3,
	3, 5, 3, 126, 10, 3, 7, 3, 128, 10, 3, 12, 3, 14, 3, 131, 11, 3, 3, 3,
	7, 3, 134, 10, 3, 12, 3, 14, 3, 137, 11, 3, 3, 3, 5, 3, 140, 10, 3, 3,
	4, 3, 4, 6, 4, 144, 10, 4, 13, 4, 14, 4, 145, 3, 4, 5, 4, 149, 10, 4, 3,
	4, 3, 4, 5, 4, 153, 10, 4, 3, 5, 3, 5, 3, 6, 5, 6, 158, 10, 6, 3, 6, 3,
	6, 6, 6, 162, 10, 6, 13, 6, 14, 6, 163, 3, 6, 5, 6, 167, 10, 6, 3, 6, 7,
	6, 170, 10, 6, 12, 6, 14, 6, 173, 11, 6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7,
	3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 188, 10, 7, 3, 8,
	3, 8, 3, 8, 5, 8, 193, 10, 8, 3, 9, 3, 9, 3, 9, 5, 9, 198, 10, 9, 3, 10,
	3, 10, 5, 10, 202, 10, 10, 3, 11, 3, 11, 3, 11, 7, 11, 207, 10, 11, 12,
	11, 14, 11, 210, 11, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 217,
	10, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12, 223, 10, 12, 3, 13, 3, 13, 3,
	14, 3, 14, 3, 15, 3, 15, 3, 16, 3, 16, 3, 16, 3, 16, 5, 16, 235, 10, 16,
	5, 16, 237, 10, 16, 3, 17, 3, 17, 3, 17, 7, 17, 242, 10, 17, 12, 17, 14,
	17, 245, 11, 17, 3, 17, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 18, 3, 19,
	3, 19, 3, 20, 3, 20, 3, 20, 3, 20, 5, 20, 260, 10, 20, 5, 20, 262, 10,
	20, 3, 21, 3, 21, 3, 21, 3, 21, 5, 21, 268, 10, 21, 3, 21, 5, 21, 271,
	10, 21, 3, 21, 5, 21, 274, 10, 21, 3, 21, 5, 21, 277, 10, 21, 3, 22, 3,
	22, 3, 22, 3, 22, 5, 22, 283, 10, 22, 3, 22, 3, 22, 3, 22, 5, 22, 288,
	10, 22, 3, 22, 5, 22, 291, 10, 22, 5, 22, 293, 10, 22, 3, 23, 3, 23, 3,
	23, 7, 23, 298, 10, 23, 12, 23, 14, 23, 301, 11, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 23, 3, 24, 3, 24, 3, 25, 3, 25, 3, 26, 3, 26, 3, 26, 7, 26, 315,
	10, 26, 12, 26, 14, 26, 318, 11, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26,
	3, 26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 5, 28, 334,
	10, 28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28, 340, 10, 28, 3, 29, 3, 29, 3,
	30, 3, 30, 5, 30, 346, 10, 30, 3, 30, 7, 30, 349, 10, 30, 12, 30, 14, 30,
	352, 11, 30, 3, 31, 3, 31, 5, 31, 356, 10, 31, 3, 32, 3, 32, 3, 33, 3,
	33, 3, 34, 3, 34, 5, 34, 364, 10, 34, 3, 34, 7, 34, 367, 10, 34, 12, 34,
	14, 34, 370, 11, 34, 3, 35, 3, 35, 5, 35, 374, 10, 35, 3, 35, 3, 35, 5,
	35, 378, 10, 35, 3, 35, 3, 35, 5, 35, 382, 10, 35, 3, 35, 7, 35, 385, 10,
	35, 12, 35, 14, 35, 388, 11, 35, 3, 35, 5, 35, 391, 10, 35, 5, 35, 393,
	10, 35, 3, 35, 3, 35, 3, 36, 3, 36, 3, 37, 3, 37, 3, 38, 3, 38, 5, 38,
	403, 10, 38, 3, 38, 7, 38, 406, 10, 38, 12, 38, 14, 38, 409, 11, 38, 3,
	39, 3, 39, 3, 40, 3, 40, 3, 41, 3, 41, 3, 42, 3, 42, 3, 43, 3, 43, 5, 43,
	421, 10, 43, 3, 43, 3, 43, 5, 43, 425, 10, 43, 3, 43, 3, 43, 5, 43, 429,
	10, 43, 3, 43, 6, 43, 432, 10, 43, 13, 43, 14, 43, 433, 3, 43, 5, 43, 437,
	10, 43, 3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 2, 2, 45, 2, 4, 6, 8, 10, 12,
	14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48,
	50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84,
	86, 2, 2, 2, 471, 2, 91, 3, 2, 2, 2, 4, 113, 3, 2, 2, 2, 6, 141, 3, 2,
	2, 2, 8, 154, 3, 2, 2, 2, 10, 157, 3, 2, 2, 2, 12, 187, 3, 2, 2, 2, 14,
	189, 3, 2, 2, 2, 16, 194, 3, 2, 2, 2, 18, 201, 3, 2, 2, 2, 20, 203, 3,
	2, 2, 2, 22, 211, 3, 2, 2, 2, 24, 224, 3, 2, 2, 2, 26, 226, 3, 2, 2, 2,
	28, 228, 3, 2, 2, 2, 30, 230, 3, 2, 2, 2, 32, 238, 3, 2, 2, 2, 34, 249,
	3, 2, 2, 2, 36, 253, 3, 2, 2, 2, 38, 255, 3, 2, 2, 2, 40, 263, 3, 2, 2,
	2, 42, 278, 3, 2, 2, 2, 44, 294, 3, 2, 2, 2, 46, 307, 3, 2, 2, 2, 48, 309,
	3, 2, 2, 2, 50, 311, 3, 2, 2, 2, 52, 326, 3, 2, 2, 2, 54, 330, 3, 2, 2,
	2, 56, 341, 3, 2, 2, 2, 58, 343, 3, 2, 2, 2, 60, 355, 3, 2, 2, 2, 62, 357,
	3, 2, 2, 2, 64, 359, 3, 2, 2, 2, 66, 361, 3, 2, 2, 2, 68, 371, 3, 2, 2,
	2, 70, 396, 3, 2, 2, 2, 72, 398, 3, 2, 2, 2, 74, 400, 3, 2, 2, 2, 76, 410,
	3, 2, 2, 2, 78, 412, 3, 2, 2, 2, 80, 414, 3, 2, 2, 2, 82, 416, 3, 2, 2,
	2, 84, 418, 3, 2, 2, 2, 86, 440, 3, 2, 2, 2, 88, 90, 7, 20, 2, 2, 89, 88,
	3, 2, 2, 2, 90, 93, 3, 2, 2, 2, 91, 89, 3, 2, 2, 2, 91, 92, 3, 2, 2, 2,
	92, 95, 3, 2, 2, 2, 93, 91, 3, 2, 2, 2, 94, 96, 5, 10, 6, 2, 95, 94, 3,
	2, 2, 2, 95, 96, 3, 2, 2, 2, 96, 98, 3, 2, 2, 2, 97, 99, 7, 20, 2, 2, 98,
	97, 3, 2, 2, 2, 99, 100, 3, 2, 2, 2, 100, 98, 3, 2, 2, 2, 100, 101, 3,
	2, 2, 2, 101, 103, 3, 2, 2, 2, 102, 104, 5, 4, 3, 2, 103, 102, 3, 2, 2,
	2, 103, 104, 3, 2, 2, 2, 104, 108, 3, 2, 2, 2, 105, 107, 7, 20, 2, 2, 106,
	105, 3, 2, 2, 2, 107, 110, 3, 2, 2, 2, 108, 106, 3, 2, 2, 2, 108, 109,
	3, 2, 2, 2, 109, 111, 3, 2, 2, 2, 110, 108, 3, 2, 2, 2, 111, 112, 7, 2,
	2, 3, 112, 3, 3, 2, 2, 2, 113, 115, 5, 6, 4, 2, 114, 116, 7, 21, 2, 2,
	115, 114, 3, 2, 2, 2, 115, 116, 3, 2, 2, 2, 116, 129, 3, 2, 2, 2, 117,
	119, 7, 20, 2, 2, 118, 117, 3, 2, 2, 2, 119, 120, 3, 2, 2, 2, 120, 118,
	3, 2, 2, 2, 120, 121, 3, 2, 2, 2, 121, 122, 3, 2, 2, 2, 122, 123, 7, 4,
	2, 2, 123, 125, 5, 6, 4, 2, 124, 126, 7, 21, 2, 2, 125, 124, 3, 2, 2, 2,
	125, 126, 3, 2, 2, 2, 126, 128, 3, 2, 2, 2, 127, 118, 3, 2, 2, 2, 128,
	131, 3, 2, 2, 2, 129, 127, 3, 2, 2, 2, 129, 130, 3, 2, 2, 2, 130, 135,
	3, 2, 2, 2, 131, 129, 3, 2, 2, 2, 132, 134, 7, 20, 2, 2, 133, 132, 3, 2,
	2, 2, 134, 137, 3, 2, 2, 2, 135, 133, 3, 2, 2, 2, 135, 136, 3, 2, 2, 2,
	136, 139, 3, 2, 2, 2, 137, 135, 3, 2, 2, 2, 138, 140, 7, 4, 2, 2, 139,
	138, 3, 2, 2, 2, 139, 140, 3, 2, 2, 2, 140, 5, 3, 2, 2, 2, 141, 143, 5,
	8, 5, 2, 142, 144, 7, 20, 2, 2, 143, 142, 3, 2, 2, 2, 144, 145, 3, 2, 2,
	2, 145, 143, 3, 2, 2, 2, 145, 146, 3, 2, 2, 2, 146, 148, 3, 2, 2, 2, 147,
	149, 7, 21, 2, 2, 148, 147, 3, 2, 2, 2, 148, 149, 3, 2, 2, 2, 149, 150,
	3, 2, 2, 2, 150, 152, 7, 3, 2, 2, 151, 153, 5, 10, 6, 2, 152, 151, 3, 2,
	2, 2, 152, 153, 3, 2, 2, 2, 153, 7, 3, 2, 2, 2, 154, 155, 7, 5, 2, 2, 155,
	9, 3, 2, 2, 2, 156, 158, 7, 21, 2, 2, 157, 156, 3, 2, 2, 2, 157, 158, 3,
	2, 2, 2, 158, 159, 3, 2, 2, 2, 159, 171, 5, 12, 7, 2, 160, 162, 7, 20,
	2, 2, 161, 160, 3, 2, 2, 2, 162, 163, 3, 2, 2, 2, 163, 161, 3, 2, 2, 2,
	163, 164, 3, 2, 2, 2, 164, 166, 3, 2, 2, 2, 165, 167, 7, 21, 2, 2, 166,
	165, 3, 2, 2, 2, 166, 167, 3, 2, 2, 2, 167, 168, 3, 2, 2, 2, 168, 170,
	5, 12, 7, 2, 169, 161, 3, 2, 2, 2, 170, 173, 3, 2, 2, 2, 171, 169, 3, 2,
	2, 2, 171, 172, 3, 2, 2, 2, 172, 11, 3, 2, 2, 2, 173, 171, 3, 2, 2, 2,
	174, 188, 5, 14, 8, 2, 175, 188, 5, 16, 9, 2, 176, 188, 5, 18, 10, 2, 177,
	188, 5, 30, 16, 2, 178, 188, 5, 32, 17, 2, 179, 188, 5, 34, 18, 2, 180,
	188, 5, 38, 20, 2, 181, 188, 5, 40, 21, 2, 182, 188, 5, 42, 22, 2, 183,
	188, 5, 44, 23, 2, 184, 188, 5, 50, 26, 2, 185, 188, 5, 52, 27, 2, 186,
	188, 5, 54, 28, 2, 187, 174, 3, 2, 2, 2, 187, 175, 3, 2, 2, 2, 187, 176,
	3, 2, 2, 2, 187, 177, 3, 2, 2, 2, 187, 178, 3, 2, 2, 2, 187, 179, 3, 2,
	2, 2, 187, 180, 3, 2, 2, 2, 187, 181, 3, 2, 2, 2, 187, 182, 3, 2, 2, 2,
	187, 183, 3, 2, 2, 2, 187, 184, 3, 2, 2, 2, 187, 185, 3, 2, 2, 2, 187,
	186, 3, 2, 2, 2, 188, 13, 3, 2, 2, 2, 189, 192, 7, 6, 2, 2, 190, 191, 7,
	21, 2, 2, 191, 193, 5, 66, 34, 2, 192, 190, 3, 2, 2, 2, 192, 193, 3, 2,
	2, 2, 193, 15, 3, 2, 2, 2, 194, 197, 7, 7, 2, 2, 195, 196, 7, 21, 2, 2,
	196, 198, 5, 66, 34, 2, 197, 195, 3, 2, 2, 2, 197, 198, 3, 2, 2, 2, 198,
	17, 3, 2, 2, 2, 199, 202, 5, 22, 12, 2, 200, 202, 5, 20, 11, 2, 201, 199,
	3, 2, 2, 2, 201, 200, 3, 2, 2, 2, 202, 19, 3, 2, 2, 2, 203, 208, 7, 9,
	2, 2, 204, 205, 7, 21, 2, 2, 205, 207, 5, 78, 40, 2, 206, 204, 3, 2, 2,
	2, 207, 210, 3, 2, 2, 2, 208, 206, 3, 2, 2, 2, 208, 209, 3, 2, 2, 2, 209,
	21, 3, 2, 2, 2, 210, 208, 3, 2, 2, 2, 211, 212, 7, 8, 2, 2, 212, 213, 7,
	21, 2, 2, 213, 216, 5, 24, 13, 2, 214, 215, 7, 21, 2, 2, 215, 217, 5, 26,
	14, 2, 216, 214, 3, 2, 2, 2, 216, 217, 3, 2, 2, 2, 217, 222, 3, 2, 2, 2,
	218, 219, 7, 21, 2, 2, 219, 220, 7, 26, 2, 2, 220, 221, 7, 21, 2, 2, 221,
	223, 5, 28, 15, 2, 222, 218, 3, 2, 2, 2, 222, 223, 3, 2, 2, 2, 223, 23,
	3, 2, 2, 2, 224, 225, 7, 25, 2, 2, 225, 25, 3, 2, 2, 2, 226, 227, 7, 25,
	2, 2, 227, 27, 3, 2, 2, 2, 228, 229, 7, 25, 2, 2, 229, 29, 3, 2, 2, 2,
	230, 236, 7, 10, 2, 2, 231, 234, 7, 21, 2, 2, 232, 235, 5, 66, 34, 2, 233,
	235, 5, 68, 35, 2, 234, 232, 3, 2, 2, 2, 234, 233, 3, 2, 2, 2, 235, 237,
	3, 2, 2, 2, 236, 231, 3, 2, 2, 2, 236, 237, 3, 2, 2, 2, 237, 31, 3, 2,
	2, 2, 238, 243, 7, 13, 2, 2, 239, 240, 7, 21, 2, 2, 240, 242, 5, 64, 33,
	2, 241, 239, 3, 2, 2, 2, 242, 245, 3, 2, 2, 2, 243, 241, 3, 2, 2, 2, 243,
	244, 3, 2, 2, 2, 244, 246, 3, 2, 2, 2, 245, 243, 3, 2, 2, 2, 246, 247,
	7, 21, 2, 2, 247, 248, 5, 82, 42, 2, 248, 33, 3, 2, 2, 2, 249, 250, 7,
	14, 2, 2, 250, 251, 7, 21, 2, 2, 251, 252, 5, 36, 19, 2, 252, 35, 3, 2,
	2, 2, 253, 254, 7, 25, 2, 2, 254, 37, 3, 2, 2, 2, 255, 261, 7, 15, 2, 2,
	256, 259, 7, 21, 2, 2, 257, 260, 5, 66, 34, 2, 258, 260, 5, 68, 35, 2,
	259, 257, 3, 2, 2, 2, 259, 258, 3, 2, 2, 2, 260, 262, 3, 2, 2, 2, 261,
	256, 3, 2, 2, 2, 261, 262, 3, 2, 2, 2, 262, 39, 3, 2, 2, 2, 263, 264, 7,
	11, 2, 2, 264, 265, 7, 21, 2, 2, 265, 270, 5, 72, 37, 2, 266, 268, 7, 21,
	2, 2, 267, 266, 3, 2, 2, 2, 267, 268, 3, 2, 2, 2, 268, 269, 3, 2, 2, 2,
	269, 271, 7, 29, 2, 2, 270, 267, 3, 2, 2, 2, 270, 271, 3, 2, 2, 2, 271,
	276, 3, 2, 2, 2, 272, 274, 7, 21, 2, 2, 273, 272, 3, 2, 2, 2, 273, 274,
	3, 2, 2, 2, 274, 275, 3, 2, 2, 2, 275, 277, 5, 74, 38, 2, 276, 273, 3,
	2, 2, 2, 276, 277, 3, 2, 2, 2, 277, 41, 3, 2, 2, 2, 278, 279, 7, 12, 2,
	2, 279, 280, 7, 21, 2, 2, 280, 292, 5, 72, 37, 2, 281, 283, 7, 21, 2, 2,
	282, 281, 3, 2, 2, 2, 282, 283, 3, 2, 2, 2, 283, 284, 3, 2, 2, 2, 284,
	285, 7, 29, 2, 2, 285, 290, 3, 2, 2, 2, 286, 288, 7, 21, 2, 2, 287, 286,
	3, 2, 2, 2, 287, 288, 3, 2, 2, 2, 288, 289, 3, 2, 2, 2, 289, 291, 5, 74,
	38, 2, 290, 287, 3, 2, 2, 2, 290, 291, 3, 2, 2, 2, 291, 293, 3, 2, 2, 2,
	292, 282, 3, 2, 2, 2, 292, 293, 3, 2, 2, 2, 293, 43, 3, 2, 2, 2, 294, 299,
	7, 16, 2, 2, 295, 296, 7, 21, 2, 2, 296, 298, 5, 64, 33, 2, 297, 295, 3,
	2, 2, 2, 298, 301, 3, 2, 2, 2, 299, 297, 3, 2, 2, 2, 299, 300, 3, 2, 2,
	2, 300, 302, 3, 2, 2, 2, 301, 299, 3, 2, 2, 2, 302, 303, 7, 21, 2, 2, 303,
	304, 5, 46, 24, 2, 304, 305, 7, 21, 2, 2, 305, 306, 5, 48, 25, 2, 306,
	45, 3, 2, 2, 2, 307, 308, 7, 25, 2, 2, 308, 47, 3, 2, 2, 2, 309, 310, 7,
	25, 2, 2, 310, 49, 3, 2, 2, 2, 311, 316, 7, 17, 2, 2, 312, 313, 7, 21,
	2, 2, 313, 315, 5, 64, 33, 2, 314, 312, 3, 2, 2, 2, 315, 318, 3, 2, 2,
	2, 316, 314, 3, 2, 2, 2, 316, 317, 3, 2, 2, 2, 317, 319, 3, 2, 2, 2, 318,
	316, 3, 2, 2, 2, 319, 320, 7, 21, 2, 2, 320, 321, 5, 82, 42, 2, 321, 322,
	7, 21, 2, 2, 322, 323, 7, 30, 2, 2, 323, 324, 7, 21, 2, 2, 324, 325, 5,
	76, 39, 2, 325, 51, 3, 2, 2, 2, 326, 327, 7, 18, 2, 2, 327, 328, 7, 21,
	2, 2, 328, 329, 5, 76, 39, 2, 329, 53, 3, 2, 2, 2, 330, 333, 5, 56, 29,
	2, 331, 332, 7, 21, 2, 2, 332, 334, 5, 58, 30, 2, 333, 331, 3, 2, 2, 2,
	333, 334, 3, 2, 2, 2, 334, 339, 3, 2, 2, 2, 335, 336, 7, 21, 2, 2, 336,
	340, 5, 66, 34, 2, 337, 338, 7, 21, 2, 2, 338, 340, 5, 84, 43, 2, 339,
	335, 3, 2, 2, 2, 339, 337, 3, 2, 2, 2, 339, 340, 3, 2, 2, 2, 340, 55, 3,
	2, 2, 2, 341, 342, 7, 19, 2, 2, 342, 57, 3, 2, 2, 2, 343, 350, 5, 60, 31,
	2, 344, 346, 7, 21, 2, 2, 345, 344, 3, 2, 2, 2, 345, 346, 3, 2, 2, 2, 346,
	347, 3, 2, 2, 2, 347, 349, 5, 60, 31, 2, 348, 345, 3, 2, 2, 2, 349, 352,
	3, 2, 2, 2, 350, 348, 3, 2, 2, 2, 350, 351, 3, 2, 2, 2, 351, 59, 3, 2,
	2, 2, 352, 350, 3, 2, 2, 2, 353, 356, 5, 62, 32, 2, 354, 356, 5, 64, 33,
	2, 355, 353, 3, 2, 2, 2, 355, 354, 3, 2, 2, 2, 356, 61, 3, 2, 2, 2, 357,
	358, 7, 24, 2, 2, 358, 63, 3, 2, 2, 2, 359, 360, 7, 23, 2, 2, 360, 65,
	3, 2, 2, 2, 361, 368, 5, 70, 36, 2, 362, 364, 7, 21, 2, 2, 363, 362, 3,
	2, 2, 2, 363, 364, 3, 2, 2, 2, 364, 365, 3, 2, 2, 2, 365, 367, 5, 70, 36,
	2, 366, 363, 3, 2, 2, 2, 367, 370, 3, 2, 2, 2, 368, 366, 3, 2, 2, 2, 368,
	369, 3, 2, 2, 2, 369, 67, 3, 2, 2, 2, 370, 368, 3, 2, 2, 2, 371, 373, 7,
	22, 2, 2, 372, 374, 7, 21, 2, 2, 373, 372, 3, 2, 2, 2, 373, 374, 3, 2,
	2, 2, 374, 392, 3, 2, 2, 2, 375, 386, 5, 70, 36, 2, 376, 378, 7, 21, 2,
	2, 377, 376, 3, 2, 2, 2, 377, 378, 3, 2, 2, 2, 378, 379, 3, 2, 2, 2, 379,
	381, 7, 28, 2, 2, 380, 382, 7, 21, 2, 2, 381, 380, 3, 2, 2, 2, 381, 382,
	3, 2, 2, 2, 382, 383, 3, 2, 2, 2, 383, 385, 5, 70, 36, 2, 384, 377, 3,
	2, 2, 2, 385, 388, 3, 2, 2, 2, 386, 384, 3, 2, 2, 2, 386, 387, 3, 2, 2,
	2, 387, 390, 3, 2, 2, 2, 388, 386, 3, 2, 2, 2, 389, 391, 7, 21, 2, 2, 390,
	389, 3, 2, 2, 2, 390, 391, 3, 2, 2, 2, 391, 393, 3, 2, 2, 2, 392, 375,
	3, 2, 2, 2, 392, 393, 3, 2, 2, 2, 393, 394, 3, 2, 2, 2, 394, 395, 7, 27,
	2, 2, 395, 69, 3, 2, 2, 2, 396, 397, 7, 25, 2, 2, 397, 71, 3, 2, 2, 2,
	398, 399, 7, 25, 2, 2, 399, 73, 3, 2, 2, 2, 400, 407, 7, 25, 2, 2, 401,
	403, 7, 21, 2, 2, 402, 401, 3, 2, 2, 2, 402, 403, 3, 2, 2, 2, 403, 404,
	3, 2, 2, 2, 404, 406, 7, 25, 2, 2, 405, 402, 3, 2, 2, 2, 406, 409, 3, 2,
	2, 2, 407, 405, 3, 2, 2, 2, 407, 408, 3, 2, 2, 2, 408, 75, 3, 2, 2, 2,
	409, 407, 3, 2, 2, 2, 410, 411, 7, 25, 2, 2, 411, 77, 3, 2, 2, 2, 412,
	413, 7, 25, 2, 2, 413, 79, 3, 2, 2, 2, 414, 415, 7, 25, 2, 2, 415, 81,
	3, 2, 2, 2, 416, 417, 7, 25, 2, 2, 417, 83, 3, 2, 2, 2, 418, 420, 7, 22,
	2, 2, 419, 421, 7, 21, 2, 2, 420, 419, 3, 2, 2, 2, 420, 421, 3, 2, 2, 2,
	421, 422, 3, 2, 2, 2, 422, 431, 5, 86, 44, 2, 423, 425, 7, 21, 2, 2, 424,
	423, 3, 2, 2, 2, 424, 425, 3, 2, 2, 2, 425, 426, 3, 2, 2, 2, 426, 428,
	7, 28, 2, 2, 427, 429, 7, 21, 2, 2, 428, 427, 3, 2, 2, 2, 428, 429, 3,
	2, 2, 2, 429, 430, 3, 2, 2, 2, 430, 432, 5, 86, 44, 2, 431, 424, 3, 2,
	2, 2, 432, 433, 3, 2, 2, 2, 433, 431, 3, 2, 2, 2, 433, 434, 3, 2, 2, 2,
	434, 436, 3, 2, 2, 2, 435, 437, 7, 21, 2, 2, 436, 435, 3, 2, 2, 2, 436,
	437, 3, 2, 2, 2, 437, 438, 3, 2, 2, 2, 438, 439, 7, 27, 2, 2, 439, 85,
	3, 2, 2, 2, 440, 441, 7, 25, 2, 2, 441, 87, 3, 2, 2, 2, 62, 91, 95, 100,
	103, 108, 115, 120, 125, 129, 135, 139, 145, 148, 152, 157, 163, 166, 171,
	187, 192, 197, 201, 208, 216, 222, 234, 236, 243, 259, 261, 267, 270, 273,
	276, 282, 287, 290, 292, 299, 316, 333, 339, 345, 350, 355, 363, 368, 373,
	377, 381, 386, 390, 392, 402, 407, 420, 424, 428, 433, 436,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "", "", "", "'FROM'", "'COPY'", "'SAVE ARTIFACT'", "'SAVE IMAGE'",
	"'RUN'", "'ENV'", "'ARG'", "'BUILD'", "'WORKDIR'", "'ENTRYPOINT'", "'GIT CLONE'",
	"'DOCKER LOAD'", "'DOCKER PULL'", "", "", "", "'['", "", "", "", "'AS LOCAL'",
	"']'", "','", "'='",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "Target", "FROM", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE",
	"RUN", "ENV", "ARG", "BUILD", "WORKDIR", "ENTRYPOINT", "GIT_CLONE", "DOCKER_LOAD",
	"DOCKER_PULL", "Command", "NL", "WS", "OPEN_BRACKET", "FlagKeyValue", "FlagKey",
	"Atom", "AS_LOCAL", "CLOSE_BRACKET", "COMMA", "EQUALS", "AS",
}

var ruleNames = []string{
	"earthFile", "targets", "target", "targetHeader", "stmts", "stmt", "fromStmt",
	"copyStmt", "saveStmt", "saveImage", "saveArtifact", "saveFrom", "saveTo",
	"saveAsLocalTo", "runStmt", "buildStmt", "workdirStmt", "workdirPath",
	"entrypointStmt", "envStmt", "argStmt", "gitCloneStmt", "gitURL", "gitCloneDest",
	"dockerLoadStmt", "dockerPullStmt", "genericCommand", "commandName", "flags",
	"flag", "flagKey", "flagKeyValue", "stmtWords", "stmtWordsList", "stmtWord",
	"envArgKey", "envArgValue", "imageName", "saveImageName", "targetName",
	"fullTargetName", "argsList", "arg",
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
	EarthParserFlagKeyValue  = 21
	EarthParserFlagKey       = 22
	EarthParserAtom          = 23
	EarthParserAS_LOCAL      = 24
	EarthParserCLOSE_BRACKET = 25
	EarthParserCOMMA         = 26
	EarthParserEQUALS        = 27
	EarthParserAS            = 28
)

// EarthParser rules.
const (
	EarthParserRULE_earthFile      = 0
	EarthParserRULE_targets        = 1
	EarthParserRULE_target         = 2
	EarthParserRULE_targetHeader   = 3
	EarthParserRULE_stmts          = 4
	EarthParserRULE_stmt           = 5
	EarthParserRULE_fromStmt       = 6
	EarthParserRULE_copyStmt       = 7
	EarthParserRULE_saveStmt       = 8
	EarthParserRULE_saveImage      = 9
	EarthParserRULE_saveArtifact   = 10
	EarthParserRULE_saveFrom       = 11
	EarthParserRULE_saveTo         = 12
	EarthParserRULE_saveAsLocalTo  = 13
	EarthParserRULE_runStmt        = 14
	EarthParserRULE_buildStmt      = 15
	EarthParserRULE_workdirStmt    = 16
	EarthParserRULE_workdirPath    = 17
	EarthParserRULE_entrypointStmt = 18
	EarthParserRULE_envStmt        = 19
	EarthParserRULE_argStmt        = 20
	EarthParserRULE_gitCloneStmt   = 21
	EarthParserRULE_gitURL         = 22
	EarthParserRULE_gitCloneDest   = 23
	EarthParserRULE_dockerLoadStmt = 24
	EarthParserRULE_dockerPullStmt = 25
	EarthParserRULE_genericCommand = 26
	EarthParserRULE_commandName    = 27
	EarthParserRULE_flags          = 28
	EarthParserRULE_flag           = 29
	EarthParserRULE_flagKey        = 30
	EarthParserRULE_flagKeyValue   = 31
	EarthParserRULE_stmtWords      = 32
	EarthParserRULE_stmtWordsList  = 33
	EarthParserRULE_stmtWord       = 34
	EarthParserRULE_envArgKey      = 35
	EarthParserRULE_envArgValue    = 36
	EarthParserRULE_imageName      = 37
	EarthParserRULE_saveImageName  = 38
	EarthParserRULE_targetName     = 39
	EarthParserRULE_fullTargetName = 40
	EarthParserRULE_argsList       = 41
	EarthParserRULE_arg            = 42
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
	p.SetState(89)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(86)
				p.Match(EarthParserNL)
			}

		}
		p.SetState(91)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<EarthParserFROM)|(1<<EarthParserCOPY)|(1<<EarthParserSAVE_ARTIFACT)|(1<<EarthParserSAVE_IMAGE)|(1<<EarthParserRUN)|(1<<EarthParserENV)|(1<<EarthParserARG)|(1<<EarthParserBUILD)|(1<<EarthParserWORKDIR)|(1<<EarthParserENTRYPOINT)|(1<<EarthParserGIT_CLONE)|(1<<EarthParserDOCKER_LOAD)|(1<<EarthParserDOCKER_PULL)|(1<<EarthParserCommand)|(1<<EarthParserWS))) != 0 {
		{
			p.SetState(92)
			p.Stmts()
		}

	}
	p.SetState(96)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(95)
				p.Match(EarthParserNL)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(98)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext())
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserTarget {
		{
			p.SetState(100)
			p.Targets()
		}

	}
	p.SetState(106)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == EarthParserNL {
		{
			p.SetState(103)
			p.Match(EarthParserNL)
		}

		p.SetState(108)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(109)
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
		p.SetState(111)
		p.Target()
	}
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(112)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(127)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(116)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(115)
					p.Match(EarthParserNL)
				}

				p.SetState(118)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(120)
				p.Match(EarthParserDEDENT)
			}
			{
				p.SetState(121)
				p.Target()
			}
			p.SetState(123)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(122)
					p.Match(EarthParserWS)
				}

			}

		}
		p.SetState(129)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
	}
	p.SetState(133)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(130)
				p.Match(EarthParserNL)
			}

		}
		p.SetState(135)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserDEDENT {
		{
			p.SetState(136)
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
		p.SetState(139)
		p.TargetHeader()
	}
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(140)
			p.Match(EarthParserNL)
		}

		p.SetState(143)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(145)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(148)
		p.Match(EarthParserINDENT)
	}
	p.SetState(150)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(149)
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
		p.SetState(152)
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
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(154)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(157)
		p.Stmt()
	}
	p.SetState(169)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(159)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(158)
					p.Match(EarthParserNL)
				}

				p.SetState(161)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			p.SetState(164)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(163)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(166)
				p.Stmt()
			}

		}
		p.SetState(171)
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

	p.SetState(185)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(172)
			p.FromStmt()
		}

	case EarthParserCOPY:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(173)
			p.CopyStmt()
		}

	case EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(174)
			p.SaveStmt()
		}

	case EarthParserRUN:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(175)
			p.RunStmt()
		}

	case EarthParserBUILD:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(176)
			p.BuildStmt()
		}

	case EarthParserWORKDIR:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(177)
			p.WorkdirStmt()
		}

	case EarthParserENTRYPOINT:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(178)
			p.EntrypointStmt()
		}

	case EarthParserENV:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(179)
			p.EnvStmt()
		}

	case EarthParserARG:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(180)
			p.ArgStmt()
		}

	case EarthParserGIT_CLONE:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(181)
			p.GitCloneStmt()
		}

	case EarthParserDOCKER_LOAD:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(182)
			p.DockerLoadStmt()
		}

	case EarthParserDOCKER_PULL:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(183)
			p.DockerPullStmt()
		}

	case EarthParserCommand:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(184)
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

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(187)
		p.Match(EarthParserFROM)
	}
	p.SetState(190)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(188)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(189)
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
		p.SetState(192)
		p.Match(EarthParserCOPY)
	}
	p.SetState(195)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(193)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(194)
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
	p.EnterRule(localctx, 16, EarthParserRULE_saveStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(199)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserSAVE_ARTIFACT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(197)
			p.SaveArtifact()
		}

	case EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(198)
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
	p.EnterRule(localctx, 18, EarthParserRULE_saveImage)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(201)
		p.Match(EarthParserSAVE_IMAGE)
	}
	p.SetState(206)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(202)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(203)
				p.SaveImageName()
			}

		}
		p.SetState(208)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 20, EarthParserRULE_saveArtifact)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(209)
		p.Match(EarthParserSAVE_ARTIFACT)
	}
	{
		p.SetState(210)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(211)
		p.SaveFrom()
	}
	p.SetState(214)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(212)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(213)
			p.SaveTo()
		}

	}
	p.SetState(220)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(216)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(217)
			p.Match(EarthParserAS_LOCAL)
		}
		{
			p.SetState(218)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(219)
			p.SaveAsLocalTo()
		}

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
	p.EnterRule(localctx, 22, EarthParserRULE_saveFrom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(222)
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
	p.EnterRule(localctx, 24, EarthParserRULE_saveTo)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(224)
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
	p.EnterRule(localctx, 26, EarthParserRULE_saveAsLocalTo)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserAtom)
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

func (s *RunStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *RunStmtContext) StmtWordsList() IStmtWordsListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordsListContext)
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
	p.EnterRule(localctx, 28, EarthParserRULE_runStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(228)
		p.Match(EarthParserRUN)
	}
	p.SetState(234)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(229)
			p.Match(EarthParserWS)
		}
		p.SetState(232)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case EarthParserAtom:
			{
				p.SetState(230)
				p.StmtWords()
			}

		case EarthParserOPEN_BRACKET:
			{
				p.SetState(231)
				p.StmtWordsList()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.EnterRule(localctx, 30, EarthParserRULE_buildStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserBUILD)
	}
	p.SetState(241)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(237)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(238)
				p.FlagKeyValue()
			}

		}
		p.SetState(243)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())
	}
	{
		p.SetState(244)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(245)
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
	p.EnterRule(localctx, 32, EarthParserRULE_workdirStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(247)
		p.Match(EarthParserWORKDIR)
	}
	{
		p.SetState(248)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(249)
		p.WorkdirPath()
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
	p.EnterRule(localctx, 34, EarthParserRULE_workdirPath)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(251)
		p.Match(EarthParserAtom)
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

func (s *EntrypointStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordsContext)
}

func (s *EntrypointStmtContext) StmtWordsList() IStmtWordsListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStmtWordsListContext)
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
		p.SetState(253)
		p.Match(EarthParserENTRYPOINT)
	}
	p.SetState(259)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(254)
			p.Match(EarthParserWS)
		}
		p.SetState(257)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case EarthParserAtom:
			{
				p.SetState(255)
				p.StmtWords()
			}

		case EarthParserOPEN_BRACKET:
			{
				p.SetState(256)
				p.StmtWordsList()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
		p.SetState(261)
		p.Match(EarthParserENV)
	}
	{
		p.SetState(262)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(263)
		p.EnvArgKey()
	}
	p.SetState(268)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) == 1 {
		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(264)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(267)
			p.Match(EarthParserEQUALS)
		}

	}
	p.SetState(274)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
		p.SetState(271)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(270)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(273)
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
		p.SetState(276)
		p.Match(EarthParserARG)
	}
	{
		p.SetState(277)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(278)
		p.EnvArgKey()
	}
	p.SetState(290)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 37, p.GetParserRuleContext()) == 1 {
		p.SetState(280)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(279)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(282)
			p.Match(EarthParserEQUALS)
		}

		p.SetState(288)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext()) == 1 {
			p.SetState(285)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(284)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(287)
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

func (s *GitCloneStmtContext) GitCloneDest() IGitCloneDestContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IGitCloneDestContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IGitCloneDestContext)
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
		p.SetState(292)
		p.Match(EarthParserGIT_CLONE)
	}
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(293)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(294)
				p.FlagKeyValue()
			}

		}
		p.SetState(299)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext())
	}
	{
		p.SetState(300)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(301)
		p.GitURL()
	}
	{
		p.SetState(302)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(303)
		p.GitCloneDest()
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
	p.EnterRule(localctx, 44, EarthParserRULE_gitURL)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserAtom)
	}

	return localctx
}

// IGitCloneDestContext is an interface to support dynamic dispatch.
type IGitCloneDestContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGitCloneDestContext differentiates from other interfaces.
	IsGitCloneDestContext()
}

type GitCloneDestContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGitCloneDestContext() *GitCloneDestContext {
	var p = new(GitCloneDestContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_gitCloneDest
	return p
}

func (*GitCloneDestContext) IsGitCloneDestContext() {}

func NewGitCloneDestContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GitCloneDestContext {
	var p = new(GitCloneDestContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_gitCloneDest

	return p
}

func (s *GitCloneDestContext) GetParser() antlr.Parser { return s.parser }

func (s *GitCloneDestContext) Atom() antlr.TerminalNode {
	return s.GetToken(EarthParserAtom, 0)
}

func (s *GitCloneDestContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GitCloneDestContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GitCloneDestContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterGitCloneDest(s)
	}
}

func (s *GitCloneDestContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitGitCloneDest(s)
	}
}

func (p *EarthParser) GitCloneDest() (localctx IGitCloneDestContext) {
	localctx = NewGitCloneDestContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, EarthParserRULE_gitCloneDest)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserAtom)
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
	p.EnterRule(localctx, 48, EarthParserRULE_dockerLoadStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(309)
		p.Match(EarthParserDOCKER_LOAD)
	}
	p.SetState(314)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(310)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(311)
				p.FlagKeyValue()
			}

		}
		p.SetState(316)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())
	}
	{
		p.SetState(317)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(318)
		p.FullTargetName()
	}
	{
		p.SetState(319)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(320)
		p.Match(EarthParserAS)
	}
	{
		p.SetState(321)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(322)
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
	p.EnterRule(localctx, 50, EarthParserRULE_dockerPullStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(324)
		p.Match(EarthParserDOCKER_PULL)
	}
	{
		p.SetState(325)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(326)
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
	p.EnterRule(localctx, 52, EarthParserRULE_genericCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.CommandName()
	}
	p.SetState(331)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(329)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(330)
			p.Flags()
		}

	}
	p.SetState(337)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(333)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(334)
			p.StmtWords()
		}

	} else if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext()) == 2 {
		{
			p.SetState(335)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(336)
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
	p.EnterRule(localctx, 54, EarthParserRULE_commandName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserCommand)
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
	p.EnterRule(localctx, 56, EarthParserRULE_flags)
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
		p.SetState(341)
		p.Flag()
	}
	p.SetState(348)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(343)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(342)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(345)
				p.Flag()
			}

		}
		p.SetState(350)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 58, EarthParserRULE_flag)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(353)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFlagKey:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(351)
			p.FlagKey()
		}

	case EarthParserFlagKeyValue:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(352)
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
	p.EnterRule(localctx, 60, EarthParserRULE_flagKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(355)
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
	p.EnterRule(localctx, 62, EarthParserRULE_flagKeyValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(357)
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
	p.EnterRule(localctx, 64, EarthParserRULE_stmtWords)
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
		p.SetState(359)
		p.StmtWord()
	}
	p.SetState(366)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(361)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(360)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(363)
				p.StmtWord()
			}

		}
		p.SetState(368)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext())
	}

	return localctx
}

// IStmtWordsListContext is an interface to support dynamic dispatch.
type IStmtWordsListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStmtWordsListContext differentiates from other interfaces.
	IsStmtWordsListContext()
}

type StmtWordsListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStmtWordsListContext() *StmtWordsListContext {
	var p = new(StmtWordsListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = EarthParserRULE_stmtWordsList
	return p
}

func (*StmtWordsListContext) IsStmtWordsListContext() {}

func NewStmtWordsListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StmtWordsListContext {
	var p = new(StmtWordsListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = EarthParserRULE_stmtWordsList

	return p
}

func (s *StmtWordsListContext) GetParser() antlr.Parser { return s.parser }

func (s *StmtWordsListContext) OPEN_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserOPEN_BRACKET, 0)
}

func (s *StmtWordsListContext) CLOSE_BRACKET() antlr.TerminalNode {
	return s.GetToken(EarthParserCLOSE_BRACKET, 0)
}

func (s *StmtWordsListContext) AllWS() []antlr.TerminalNode {
	return s.GetTokens(EarthParserWS)
}

func (s *StmtWordsListContext) WS(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserWS, i)
}

func (s *StmtWordsListContext) AllStmtWord() []IStmtWordContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStmtWordContext)(nil)).Elem())
	var tst = make([]IStmtWordContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStmtWordContext)
		}
	}

	return tst
}

func (s *StmtWordsListContext) StmtWord(i int) IStmtWordContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStmtWordContext)
}

func (s *StmtWordsListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EarthParserCOMMA)
}

func (s *StmtWordsListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EarthParserCOMMA, i)
}

func (s *StmtWordsListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StmtWordsListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StmtWordsListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.EnterStmtWordsList(s)
	}
}

func (s *StmtWordsListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EarthParserListener); ok {
		listenerT.ExitStmtWordsList(s)
	}
}

func (p *EarthParser) StmtWordsList() (localctx IStmtWordsListContext) {
	localctx = NewStmtWordsListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, EarthParserRULE_stmtWordsList)
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
		p.SetState(369)
		p.Match(EarthParserOPEN_BRACKET)
	}
	p.SetState(371)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(370)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(390)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserAtom {
		{
			p.SetState(373)
			p.StmtWord()
		}
		p.SetState(384)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(375)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == EarthParserWS {
					{
						p.SetState(374)
						p.Match(EarthParserWS)
					}

				}
				{
					p.SetState(377)
					p.Match(EarthParserCOMMA)
				}
				p.SetState(379)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == EarthParserWS {
					{
						p.SetState(378)
						p.Match(EarthParserWS)
					}

				}
				{
					p.SetState(381)
					p.StmtWord()
				}

			}
			p.SetState(386)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext())
		}
		p.SetState(388)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == EarthParserWS {
			{
				p.SetState(387)
				p.Match(EarthParserWS)
			}

		}

	}
	{
		p.SetState(392)
		p.Match(EarthParserCLOSE_BRACKET)
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
	p.EnterRule(localctx, 68, EarthParserRULE_stmtWord)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(394)
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
	p.EnterRule(localctx, 70, EarthParserRULE_envArgKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.EnterRule(localctx, 72, EarthParserRULE_envArgValue)
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
		p.SetState(398)
		p.Match(EarthParserAtom)
	}
	p.SetState(405)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 54, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(400)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(399)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(402)
				p.Match(EarthParserAtom)
			}

		}
		p.SetState(407)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 54, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 74, EarthParserRULE_imageName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(408)
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
	p.EnterRule(localctx, 76, EarthParserRULE_saveImageName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.EnterRule(localctx, 78, EarthParserRULE_targetName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.EnterRule(localctx, 80, EarthParserRULE_fullTargetName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.EnterRule(localctx, 82, EarthParserRULE_argsList)
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
		p.SetState(416)
		p.Match(EarthParserOPEN_BRACKET)
	}
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
		p.Arg()
	}
	p.SetState(429)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
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
				p.Match(EarthParserCOMMA)
			}
			p.SetState(426)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == EarthParserWS {
				{
					p.SetState(425)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(428)
				p.Arg()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(431)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 58, p.GetParserRuleContext())
	}
	p.SetState(434)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == EarthParserWS {
		{
			p.SetState(433)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(436)
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
	p.EnterRule(localctx, 84, EarthParserRULE_arg)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserAtom)
	}

	return localctx
}
