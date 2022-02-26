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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 46, 652, 
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
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 3, 2, 7, 2, 138, 10, 2, 12, 2, 
	14, 2, 141, 11, 2, 3, 2, 5, 2, 144, 10, 2, 3, 2, 3, 2, 3, 2, 5, 2, 149, 
	10, 2, 3, 2, 7, 2, 152, 10, 2, 12, 2, 14, 2, 155, 11, 2, 3, 2, 5, 2, 158, 
	10, 2, 3, 2, 7, 2, 161, 10, 2, 12, 2, 14, 2, 164, 11, 2, 3, 2, 3, 2, 3, 
	3, 3, 3, 7, 3, 170, 10, 3, 12, 3, 14, 3, 173, 11, 3, 3, 3, 7, 3, 176, 10, 
	3, 12, 3, 14, 3, 179, 11, 3, 3, 4, 3, 4, 5, 4, 183, 10, 4, 3, 5, 3, 5, 
	6, 5, 187, 10, 5, 13, 5, 14, 5, 188, 3, 5, 5, 5, 192, 10, 5, 3, 5, 3, 5, 
	3, 5, 6, 5, 197, 10, 5, 13, 5, 14, 5, 198, 3, 5, 3, 5, 5, 5, 203, 10, 5, 
	3, 6, 3, 6, 3, 7, 3, 7, 6, 7, 209, 10, 7, 13, 7, 14, 7, 210, 3, 7, 5, 7, 
	214, 10, 7, 3, 7, 3, 7, 3, 7, 6, 7, 219, 10, 7, 13, 7, 14, 7, 220, 3, 7, 
	3, 7, 5, 7, 225, 10, 7, 3, 8, 3, 8, 3, 9, 5, 9, 230, 10, 9, 3, 9, 3, 9, 
	6, 9, 234, 10, 9, 13, 9, 14, 9, 235, 3, 9, 5, 9, 239, 10, 9, 3, 9, 7, 9, 
	242, 10, 9, 12, 9, 14, 9, 245, 11, 9, 3, 10, 3, 10, 3, 10, 3, 10, 5, 10, 
	251, 10, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 
	11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 
	3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 280, 10, 
	11, 3, 12, 3, 12, 3, 12, 3, 12, 6, 12, 286, 10, 12, 13, 12, 14, 12, 287, 
	3, 13, 3, 13, 6, 13, 292, 10, 13, 13, 13, 14, 13, 293, 3, 13, 5, 13, 297, 
	10, 13, 3, 13, 5, 13, 300, 10, 13, 3, 13, 6, 13, 303, 10, 13, 13, 13, 14, 
	13, 304, 3, 13, 5, 13, 308, 10, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 
	3, 15, 3, 15, 3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 5, 17, 323, 10, 
	17, 3, 18, 3, 18, 6, 18, 327, 10, 18, 13, 18, 14, 18, 328, 3, 18, 5, 18, 
	332, 10, 18, 3, 18, 7, 18, 335, 10, 18, 12, 18, 14, 18, 338, 11, 18, 3, 
	18, 6, 18, 341, 10, 18, 13, 18, 14, 18, 342, 3, 18, 5, 18, 346, 10, 18, 
	3, 18, 5, 18, 349, 10, 18, 3, 18, 6, 18, 352, 10, 18, 13, 18, 14, 18, 353, 
	3, 18, 5, 18, 357, 10, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 6, 
	19, 365, 10, 19, 13, 19, 14, 19, 366, 3, 19, 5, 19, 370, 10, 19, 3, 19, 
	5, 19, 373, 10, 19, 3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 3, 21, 6, 21, 381, 
	10, 21, 13, 21, 14, 21, 382, 3, 21, 5, 21, 386, 10, 21, 3, 21, 5, 21, 389, 
	10, 21, 3, 22, 3, 22, 3, 23, 3, 23, 6, 23, 395, 10, 23, 13, 23, 14, 23, 
	396, 3, 23, 5, 23, 400, 10, 23, 3, 23, 5, 23, 403, 10, 23, 3, 24, 3, 24, 
	3, 25, 3, 25, 3, 26, 3, 26, 3, 27, 3, 27, 6, 27, 413, 10, 27, 13, 27, 14, 
	27, 414, 3, 27, 5, 27, 418, 10, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3, 28, 
	3, 28, 6, 28, 426, 10, 28, 13, 28, 14, 28, 427, 3, 28, 5, 28, 431, 10, 
	28, 3, 28, 5, 28, 434, 10, 28, 3, 29, 3, 29, 3, 30, 3, 30, 3, 31, 3, 31, 
	3, 31, 5, 31, 443, 10, 31, 3, 32, 3, 32, 3, 32, 5, 32, 448, 10, 32, 3, 
	33, 3, 33, 3, 33, 5, 33, 453, 10, 33, 3, 34, 3, 34, 3, 34, 5, 34, 458, 
	10, 34, 3, 35, 3, 35, 5, 35, 462, 10, 35, 3, 36, 3, 36, 3, 36, 5, 36, 467, 
	10, 36, 3, 37, 3, 37, 3, 37, 5, 37, 472, 10, 37, 3, 38, 3, 38, 3, 38, 5, 
	38, 477, 10, 38, 3, 39, 3, 39, 3, 39, 5, 39, 482, 10, 39, 3, 40, 3, 40, 
	3, 40, 5, 40, 487, 10, 40, 3, 41, 3, 41, 3, 41, 5, 41, 492, 10, 41, 3, 
	42, 3, 42, 3, 42, 5, 42, 497, 10, 42, 3, 43, 3, 43, 3, 43, 5, 43, 502, 
	10, 43, 3, 44, 3, 44, 3, 44, 5, 44, 507, 10, 44, 3, 45, 3, 45, 3, 45, 5, 
	45, 512, 10, 45, 3, 46, 3, 46, 3, 46, 3, 46, 5, 46, 518, 10, 46, 3, 46, 
	5, 46, 521, 10, 46, 3, 46, 5, 46, 524, 10, 46, 3, 46, 5, 46, 527, 10, 46, 
	3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 5, 47, 534, 10, 47, 3, 47, 3, 47, 3, 
	47, 5, 47, 539, 10, 47, 3, 47, 5, 47, 542, 10, 47, 5, 47, 544, 10, 47, 
	3, 48, 3, 48, 5, 48, 548, 10, 48, 3, 49, 3, 49, 3, 50, 3, 50, 5, 50, 554, 
	10, 50, 3, 50, 7, 50, 557, 10, 50, 12, 50, 14, 50, 560, 11, 50, 3, 51, 
	3, 51, 3, 51, 3, 51, 5, 51, 566, 10, 51, 3, 51, 3, 51, 5, 51, 570, 10, 
	51, 3, 51, 3, 51, 7, 51, 574, 10, 51, 12, 51, 14, 51, 577, 11, 51, 3, 52, 
	3, 52, 3, 53, 3, 53, 3, 54, 3, 54, 3, 54, 5, 54, 586, 10, 54, 3, 55, 3, 
	55, 3, 55, 5, 55, 591, 10, 55, 3, 56, 3, 56, 3, 56, 5, 56, 596, 10, 56, 
	3, 57, 3, 57, 3, 57, 5, 57, 601, 10, 57, 3, 58, 3, 58, 3, 58, 5, 58, 606, 
	10, 58, 3, 59, 3, 59, 3, 59, 5, 59, 611, 10, 59, 3, 60, 3, 60, 3, 60, 5, 
	60, 616, 10, 60, 3, 61, 3, 61, 3, 61, 5, 61, 621, 10, 61, 3, 62, 3, 62, 
	3, 62, 5, 62, 626, 10, 62, 3, 63, 3, 63, 3, 63, 5, 63, 631, 10, 63, 3, 
	64, 3, 64, 3, 64, 5, 64, 636, 10, 64, 3, 65, 3, 65, 3, 66, 3, 66, 3, 67, 
	3, 67, 3, 67, 7, 67, 645, 10, 67, 12, 67, 14, 67, 648, 11, 67, 3, 68, 3, 
	68, 3, 68, 2, 2, 69, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 
	30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 
	66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 
	102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 
	132, 134, 2, 2, 2, 704, 2, 139, 3, 2, 2, 2, 4, 167, 3, 2, 2, 2, 6, 182, 
	3, 2, 2, 2, 8, 184, 3, 2, 2, 2, 10, 204, 3, 2, 2, 2, 12, 206, 3, 2, 2, 
	2, 14, 226, 3, 2, 2, 2, 16, 229, 3, 2, 2, 2, 18, 250, 3, 2, 2, 2, 20, 279, 
	3, 2, 2, 2, 22, 281, 3, 2, 2, 2, 24, 289, 3, 2, 2, 2, 26, 311, 3, 2, 2, 
	2, 28, 313, 3, 2, 2, 2, 30, 317, 3, 2, 2, 2, 32, 319, 3, 2, 2, 2, 34, 324, 
	3, 2, 2, 2, 36, 360, 3, 2, 2, 2, 38, 374, 3, 2, 2, 2, 40, 376, 3, 2, 2, 
	2, 42, 390, 3, 2, 2, 2, 44, 392, 3, 2, 2, 2, 46, 404, 3, 2, 2, 2, 48, 406, 
	3, 2, 2, 2, 50, 408, 3, 2, 2, 2, 52, 410, 3, 2, 2, 2, 54, 421, 3, 2, 2, 
	2, 56, 435, 3, 2, 2, 2, 58, 437, 3, 2, 2, 2, 60, 439, 3, 2, 2, 2, 62, 444, 
	3, 2, 2, 2, 64, 449, 3, 2, 2, 2, 66, 454, 3, 2, 2, 2, 68, 461, 3, 2, 2, 
	2, 70, 463, 3, 2, 2, 2, 72, 468, 3, 2, 2, 2, 74, 473, 3, 2, 2, 2, 76, 478, 
	3, 2, 2, 2, 78, 483, 3, 2, 2, 2, 80, 488, 3, 2, 2, 2, 82, 493, 3, 2, 2, 
	2, 84, 498, 3, 2, 2, 2, 86, 503, 3, 2, 2, 2, 88, 508, 3, 2, 2, 2, 90, 513, 
	3, 2, 2, 2, 92, 528, 3, 2, 2, 2, 94, 547, 3, 2, 2, 2, 96, 549, 3, 2, 2, 
	2, 98, 551, 3, 2, 2, 2, 100, 561, 3, 2, 2, 2, 102, 578, 3, 2, 2, 2, 104, 
	580, 3, 2, 2, 2, 106, 582, 3, 2, 2, 2, 108, 587, 3, 2, 2, 2, 110, 592, 
	3, 2, 2, 2, 112, 597, 3, 2, 2, 2, 114, 602, 3, 2, 2, 2, 116, 607, 3, 2, 
	2, 2, 118, 612, 3, 2, 2, 2, 120, 617, 3, 2, 2, 2, 122, 622, 3, 2, 2, 2, 
	124, 627, 3, 2, 2, 2, 126, 632, 3, 2, 2, 2, 128, 637, 3, 2, 2, 2, 130, 
	639, 3, 2, 2, 2, 132, 641, 3, 2, 2, 2, 134, 649, 3, 2, 2, 2, 136, 138, 
	7, 40, 2, 2, 137, 136, 3, 2, 2, 2, 138, 141, 3, 2, 2, 2, 139, 137, 3, 2, 
	2, 2, 139, 140, 3, 2, 2, 2, 140, 143, 3, 2, 2, 2, 141, 139, 3, 2, 2, 2, 
	142, 144, 5, 22, 12, 2, 143, 142, 3, 2, 2, 2, 143, 144, 3, 2, 2, 2, 144, 
	148, 3, 2, 2, 2, 145, 146, 5, 16, 9, 2, 146, 147, 7, 40, 2, 2, 147, 149, 
	3, 2, 2, 2, 148, 145, 3, 2, 2, 2, 148, 149, 3, 2, 2, 2, 149, 153, 3, 2, 
	2, 2, 150, 152, 7, 40, 2, 2, 151, 150, 3, 2, 2, 2, 152, 155, 3, 2, 2, 2, 
	153, 151, 3, 2, 2, 2, 153, 154, 3, 2, 2, 2, 154, 157, 3, 2, 2, 2, 155, 
	153, 3, 2, 2, 2, 156, 158, 5, 4, 3, 2, 157, 156, 3, 2, 2, 2, 157, 158, 
	3, 2, 2, 2, 158, 162, 3, 2, 2, 2, 159, 161, 7, 40, 2, 2, 160, 159, 3, 2, 
	2, 2, 161, 164, 3, 2, 2, 2, 162, 160, 3, 2, 2, 2, 162, 163, 3, 2, 2, 2, 
	163, 165, 3, 2, 2, 2, 164, 162, 3, 2, 2, 2, 165, 166, 7, 2, 2, 3, 166, 
	3, 3, 2, 2, 2, 167, 177, 5, 6, 4, 2, 168, 170, 7, 40, 2, 2, 169, 168, 3, 
	2, 2, 2, 170, 173, 3, 2, 2, 2, 171, 169, 3, 2, 2, 2, 171, 172, 3, 2, 2, 
	2, 172, 174, 3, 2, 2, 2, 173, 171, 3, 2, 2, 2, 174, 176, 5, 6, 4, 2, 175, 
	171, 3, 2, 2, 2, 176, 179, 3, 2, 2, 2, 177, 175, 3, 2, 2, 2, 177, 178, 
	3, 2, 2, 2, 178, 5, 3, 2, 2, 2, 179, 177, 3, 2, 2, 2, 180, 183, 5, 8, 5, 
	2, 181, 183, 5, 12, 7, 2, 182, 180, 3, 2, 2, 2, 182, 181, 3, 2, 2, 2, 183, 
	7, 3, 2, 2, 2, 184, 186, 5, 10, 6, 2, 185, 187, 7, 40, 2, 2, 186, 185, 
	3, 2, 2, 2, 187, 188, 3, 2, 2, 2, 188, 186, 3, 2, 2, 2, 188, 189, 3, 2, 
	2, 2, 189, 191, 3, 2, 2, 2, 190, 192, 7, 41, 2, 2, 191, 190, 3, 2, 2, 2, 
	191, 192, 3, 2, 2, 2, 192, 202, 3, 2, 2, 2, 193, 194, 7, 3, 2, 2, 194, 
	196, 5, 16, 9, 2, 195, 197, 7, 40, 2, 2, 196, 195, 3, 2, 2, 2, 197, 198, 
	3, 2, 2, 2, 198, 196, 3, 2, 2, 2, 198, 199, 3, 2, 2, 2, 199, 200, 3, 2, 
	2, 2, 200, 201, 7, 4, 2, 2, 201, 203, 3, 2, 2, 2, 202, 193, 3, 2, 2, 2, 
	202, 203, 3, 2, 2, 2, 203, 9, 3, 2, 2, 2, 204, 205, 7, 5, 2, 2, 205, 11, 
	3, 2, 2, 2, 206, 208, 5, 14, 8, 2, 207, 209, 7, 40, 2, 2, 208, 207, 3, 
	2, 2, 2, 209, 210, 3, 2, 2, 2, 210, 208, 3, 2, 2, 2, 210, 211, 3, 2, 2, 
	2, 211, 213, 3, 2, 2, 2, 212, 214, 7, 41, 2, 2, 213, 212, 3, 2, 2, 2, 213, 
	214, 3, 2, 2, 2, 214, 224, 3, 2, 2, 2, 215, 216, 7, 3, 2, 2, 216, 218, 
	5, 16, 9, 2, 217, 219, 7, 40, 2, 2, 218, 217, 3, 2, 2, 2, 219, 220, 3, 
	2, 2, 2, 220, 218, 3, 2, 2, 2, 220, 221, 3, 2, 2, 2, 221, 222, 3, 2, 2, 
	2, 222, 223, 7, 4, 2, 2, 223, 225, 3, 2, 2, 2, 224, 215, 3, 2, 2, 2, 224, 
	225, 3, 2, 2, 2, 225, 13, 3, 2, 2, 2, 226, 227, 7, 6, 2, 2, 227, 15, 3, 
	2, 2, 2, 228, 230, 7, 41, 2, 2, 229, 228, 3, 2, 2, 2, 229, 230, 3, 2, 2, 
	2, 230, 231, 3, 2, 2, 2, 231, 243, 5, 18, 10, 2, 232, 234, 7, 40, 2, 2, 
	233, 232, 3, 2, 2, 2, 234, 235, 3, 2, 2, 2, 235, 233, 3, 2, 2, 2, 235, 
	236, 3, 2, 2, 2, 236, 238, 3, 2, 2, 2, 237, 239, 7, 41, 2, 2, 238, 237, 
	3, 2, 2, 2, 238, 239, 3, 2, 2, 2, 239, 240, 3, 2, 2, 2, 240, 242, 5, 18, 
	10, 2, 241, 233, 3, 2, 2, 2, 242, 245, 3, 2, 2, 2, 243, 241, 3, 2, 2, 2, 
	243, 244, 3, 2, 2, 2, 244, 17, 3, 2, 2, 2, 245, 243, 3, 2, 2, 2, 246, 251, 
	5, 20, 11, 2, 247, 251, 5, 24, 13, 2, 248, 251, 5, 34, 18, 2, 249, 251, 
	5, 52, 27, 2, 250, 246, 3, 2, 2, 2, 250, 247, 3, 2, 2, 2, 250, 248, 3, 
	2, 2, 2, 250, 249, 3, 2, 2, 2, 251, 19, 3, 2, 2, 2, 252, 280, 5, 60, 31, 
	2, 253, 280, 5, 62, 32, 2, 254, 280, 5, 64, 33, 2, 255, 280, 5, 66, 34, 
	2, 256, 280, 5, 68, 35, 2, 257, 280, 5, 74, 38, 2, 258, 280, 5, 76, 39, 
	2, 259, 280, 5, 78, 40, 2, 260, 280, 5, 80, 41, 2, 261, 280, 5, 82, 42, 
	2, 262, 280, 5, 84, 43, 2, 263, 280, 5, 86, 44, 2, 264, 280, 5, 88, 45, 
	2, 265, 280, 5, 90, 46, 2, 266, 280, 5, 92, 47, 2, 267, 280, 5, 100, 51, 
	2, 268, 280, 5, 106, 54, 2, 269, 280, 5, 108, 55, 2, 270, 280, 5, 110, 
	56, 2, 271, 280, 5, 112, 57, 2, 272, 280, 5, 114, 58, 2, 273, 280, 5, 116, 
	59, 2, 274, 280, 5, 118, 60, 2, 275, 280, 5, 120, 61, 2, 276, 280, 5, 122, 
	62, 2, 277, 280, 5, 124, 63, 2, 278, 280, 5, 126, 64, 2, 279, 252, 3, 2, 
	2, 2, 279, 253, 3, 2, 2, 2, 279, 254, 3, 2, 2, 2, 279, 255, 3, 2, 2, 2, 
	279, 256, 3, 2, 2, 2, 279, 257, 3, 2, 2, 2, 279, 258, 3, 2, 2, 2, 279, 
	259, 3, 2, 2, 2, 279, 260, 3, 2, 2, 2, 279, 261, 3, 2, 2, 2, 279, 262, 
	3, 2, 2, 2, 279, 263, 3, 2, 2, 2, 279, 264, 3, 2, 2, 2, 279, 265, 3, 2, 
	2, 2, 279, 266, 3, 2, 2, 2, 279, 267, 3, 2, 2, 2, 279, 268, 3, 2, 2, 2, 
	279, 269, 3, 2, 2, 2, 279, 270, 3, 2, 2, 2, 279, 271, 3, 2, 2, 2, 279, 
	272, 3, 2, 2, 2, 279, 273, 3, 2, 2, 2, 279, 274, 3, 2, 2, 2, 279, 275, 
	3, 2, 2, 2, 279, 276, 3, 2, 2, 2, 279, 277, 3, 2, 2, 2, 279, 278, 3, 2, 
	2, 2, 280, 21, 3, 2, 2, 2, 281, 282, 7, 33, 2, 2, 282, 283, 7, 41, 2, 2, 
	283, 285, 5, 132, 67, 2, 284, 286, 7, 40, 2, 2, 285, 284, 3, 2, 2, 2, 286, 
	287, 3, 2, 2, 2, 287, 285, 3, 2, 2, 2, 287, 288, 3, 2, 2, 2, 288, 23, 3, 
	2, 2, 2, 289, 299, 5, 28, 15, 2, 290, 292, 7, 40, 2, 2, 291, 290, 3, 2, 
	2, 2, 292, 293, 3, 2, 2, 2, 293, 291, 3, 2, 2, 2, 293, 294, 3, 2, 2, 2, 
	294, 296, 3, 2, 2, 2, 295, 297, 7, 41, 2, 2, 296, 295, 3, 2, 2, 2, 296, 
	297, 3, 2, 2, 2, 297, 298, 3, 2, 2, 2, 298, 300, 5, 26, 14, 2, 299, 291, 
	3, 2, 2, 2, 299, 300, 3, 2, 2, 2, 300, 302, 3, 2, 2, 2, 301, 303, 7, 40, 
	2, 2, 302, 301, 3, 2, 2, 2, 303, 304, 3, 2, 2, 2, 304, 302, 3, 2, 2, 2, 
	304, 305, 3, 2, 2, 2, 305, 307, 3, 2, 2, 2, 306, 308, 7, 41, 2, 2, 307, 
	306, 3, 2, 2, 2, 307, 308, 3, 2, 2, 2, 308, 309, 3, 2, 2, 2, 309, 310, 
	7, 44, 2, 2, 310, 25, 3, 2, 2, 2, 311, 312, 5, 16, 9, 2, 312, 27, 3, 2, 
	2, 2, 313, 314, 7, 36, 2, 2, 314, 315, 7, 41, 2, 2, 315, 316, 5, 30, 16, 
	2, 316, 29, 3, 2, 2, 2, 317, 318, 5, 32, 17, 2, 318, 31, 3, 2, 2, 2, 319, 
	322, 7, 37, 2, 2, 320, 321, 7, 41, 2, 2, 321, 323, 5, 132, 67, 2, 322, 
	320, 3, 2, 2, 2, 322, 323, 3, 2, 2, 2, 323, 33, 3, 2, 2, 2, 324, 336, 5, 
	36, 19, 2, 325, 327, 7, 40, 2, 2, 326, 325, 3, 2, 2, 2, 327, 328, 3, 2, 
	2, 2, 328, 326, 3, 2, 2, 2, 328, 329, 3, 2, 2, 2, 329, 331, 3, 2, 2, 2, 
	330, 332, 7, 41, 2, 2, 331, 330, 3, 2, 2, 2, 331, 332, 3, 2, 2, 2, 332, 
	333, 3, 2, 2, 2, 333, 335, 5, 40, 21, 2, 334, 326, 3, 2, 2, 2, 335, 338, 
	3, 2, 2, 2, 336, 334, 3, 2, 2, 2, 336, 337, 3, 2, 2, 2, 337, 348, 3, 2, 
	2, 2, 338, 336, 3, 2, 2, 2, 339, 341, 7, 40, 2, 2, 340, 339, 3, 2, 2, 2, 
	341, 342, 3, 2, 2, 2, 342, 340, 3, 2, 2, 2, 342, 343, 3, 2, 2, 2, 343, 
	345, 3, 2, 2, 2, 344, 346, 7, 41, 2, 2, 345, 344, 3, 2, 2, 2, 345, 346, 
	3, 2, 2, 2, 346, 347, 3, 2, 2, 2, 347, 349, 5, 44, 23, 2, 348, 340, 3, 
	2, 2, 2, 348, 349, 3, 2, 2, 2, 349, 351, 3, 2, 2, 2, 350, 352, 7, 40, 2, 
	2, 351, 350, 3, 2, 2, 2, 352, 353, 3, 2, 2, 2, 353, 351, 3, 2, 2, 2, 353, 
	354, 3, 2, 2, 2, 354, 356, 3, 2, 2, 2, 355, 357, 7, 41, 2, 2, 356, 355, 
	3, 2, 2, 2, 356, 357, 3, 2, 2, 2, 357, 358, 3, 2, 2, 2, 358, 359, 7, 44, 
	2, 2, 359, 35, 3, 2, 2, 2, 360, 361, 7, 38, 2, 2, 361, 362, 7, 41, 2, 2, 
	362, 372, 5, 48, 25, 2, 363, 365, 7, 40, 2, 2, 364, 363, 3, 2, 2, 2, 365, 
	366, 3, 2, 2, 2, 366, 364, 3, 2, 2, 2, 366, 367, 3, 2, 2, 2, 367, 369, 
	3, 2, 2, 2, 368, 370, 7, 41, 2, 2, 369, 368, 3, 2, 2, 2, 369, 370, 3, 2, 
	2, 2, 370, 371, 3, 2, 2, 2, 371, 373, 5, 38, 20, 2, 372, 364, 3, 2, 2, 
	2, 372, 373, 3, 2, 2, 2, 373, 37, 3, 2, 2, 2, 374, 375, 5, 16, 9, 2, 375, 
	39, 3, 2, 2, 2, 376, 377, 7, 43, 2, 2, 377, 378, 7, 41, 2, 2, 378, 388, 
	5, 50, 26, 2, 379, 381, 7, 40, 2, 2, 380, 379, 3, 2, 2, 2, 381, 382, 3, 
	2, 2, 2, 382, 380, 3, 2, 2, 2, 382, 383, 3, 2, 2, 2, 383, 385, 3, 2, 2, 
	2, 384, 386, 7, 41, 2, 2, 385, 384, 3, 2, 2, 2, 385, 386, 3, 2, 2, 2, 386, 
	387, 3, 2, 2, 2, 387, 389, 5, 42, 22, 2, 388, 380, 3, 2, 2, 2, 388, 389, 
	3, 2, 2, 2, 389, 41, 3, 2, 2, 2, 390, 391, 5, 16, 9, 2, 391, 43, 3, 2, 
	2, 2, 392, 402, 7, 42, 2, 2, 393, 395, 7, 40, 2, 2, 394, 393, 3, 2, 2, 
	2, 395, 396, 3, 2, 2, 2, 396, 394, 3, 2, 2, 2, 396, 397, 3, 2, 2, 2, 397, 
	399, 3, 2, 2, 2, 398, 400, 7, 41, 2, 2, 399, 398, 3, 2, 2, 2, 399, 400, 
	3, 2, 2, 2, 400, 401, 3, 2, 2, 2, 401, 403, 5, 46, 24, 2, 402, 394, 3, 
	2, 2, 2, 402, 403, 3, 2, 2, 2, 403, 45, 3, 2, 2, 2, 404, 405, 5, 16, 9, 
	2, 405, 47, 3, 2, 2, 2, 406, 407, 5, 128, 65, 2, 407, 49, 3, 2, 2, 2, 408, 
	409, 5, 128, 65, 2, 409, 51, 3, 2, 2, 2, 410, 412, 5, 54, 28, 2, 411, 413, 
	7, 40, 2, 2, 412, 411, 3, 2, 2, 2, 413, 414, 3, 2, 2, 2, 414, 412, 3, 2, 
	2, 2, 414, 415, 3, 2, 2, 2, 415, 417, 3, 2, 2, 2, 416, 418, 7, 41, 2, 2, 
	417, 416, 3, 2, 2, 2, 417, 418, 3, 2, 2, 2, 418, 419, 3, 2, 2, 2, 419, 
	420, 7, 44, 2, 2, 420, 53, 3, 2, 2, 2, 421, 422, 7, 39, 2, 2, 422, 423, 
	7, 41, 2, 2, 423, 433, 5, 58, 30, 2, 424, 426, 7, 40, 2, 2, 425, 424, 3, 
	2, 2, 2, 426, 427, 3, 2, 2, 2, 427, 425, 3, 2, 2, 2, 427, 428, 3, 2, 2, 
	2, 428, 430, 3, 2, 2, 2, 429, 431, 7, 41, 2, 2, 430, 429, 3, 2, 2, 2, 430, 
	431, 3, 2, 2, 2, 431, 432, 3, 2, 2, 2, 432, 434, 5, 56, 29, 2, 433, 425, 
	3, 2, 2, 2, 433, 434, 3, 2, 2, 2, 434, 55, 3, 2, 2, 2, 435, 436, 5, 16, 
	9, 2, 436, 57, 3, 2, 2, 2, 437, 438, 5, 132, 67, 2, 438, 59, 3, 2, 2, 2, 
	439, 442, 7, 7, 2, 2, 440, 441, 7, 41, 2, 2, 441, 443, 5, 132, 67, 2, 442, 
	440, 3, 2, 2, 2, 442, 443, 3, 2, 2, 2, 443, 61, 3, 2, 2, 2, 444, 447, 7, 
	8, 2, 2, 445, 446, 7, 41, 2, 2, 446, 448, 5, 132, 67, 2, 447, 445, 3, 2, 
	2, 2, 447, 448, 3, 2, 2, 2, 448, 63, 3, 2, 2, 2, 449, 452, 7, 9, 2, 2, 
	450, 451, 7, 41, 2, 2, 451, 453, 5, 132, 67, 2, 452, 450, 3, 2, 2, 2, 452, 
	453, 3, 2, 2, 2, 453, 65, 3, 2, 2, 2, 454, 457, 7, 10, 2, 2, 455, 456, 
	7, 41, 2, 2, 456, 458, 5, 132, 67, 2, 457, 455, 3, 2, 2, 2, 457, 458, 3, 
	2, 2, 2, 458, 67, 3, 2, 2, 2, 459, 462, 5, 72, 37, 2, 460, 462, 5, 70, 
	36, 2, 461, 459, 3, 2, 2, 2, 461, 460, 3, 2, 2, 2, 462, 69, 3, 2, 2, 2, 
	463, 466, 7, 12, 2, 2, 464, 465, 7, 41, 2, 2, 465, 467, 5, 132, 67, 2, 
	466, 464, 3, 2, 2, 2, 466, 467, 3, 2, 2, 2, 467, 71, 3, 2, 2, 2, 468, 471, 
	7, 11, 2, 2, 469, 470, 7, 41, 2, 2, 470, 472, 5, 132, 67, 2, 471, 469, 
	3, 2, 2, 2, 471, 472, 3, 2, 2, 2, 472, 73, 3, 2, 2, 2, 473, 476, 7, 13, 
	2, 2, 474, 475, 7, 41, 2, 2, 475, 477, 5, 130, 66, 2, 476, 474, 3, 2, 2, 
	2, 476, 477, 3, 2, 2, 2, 477, 75, 3, 2, 2, 2, 478, 481, 7, 19, 2, 2, 479, 
	480, 7, 41, 2, 2, 480, 482, 5, 132, 67, 2, 481, 479, 3, 2, 2, 2, 481, 482, 
	3, 2, 2, 2, 482, 77, 3, 2, 2, 2, 483, 486, 7, 20, 2, 2, 484, 485, 7, 41, 
	2, 2, 485, 487, 5, 132, 67, 2, 486, 484, 3, 2, 2, 2, 486, 487, 3, 2, 2, 
	2, 487, 79, 3, 2, 2, 2, 488, 491, 7, 21, 2, 2, 489, 490, 7, 41, 2, 2, 490, 
	492, 5, 132, 67, 2, 491, 489, 3, 2, 2, 2, 491, 492, 3, 2, 2, 2, 492, 81, 
	3, 2, 2, 2, 493, 496, 7, 22, 2, 2, 494, 495, 7, 41, 2, 2, 495, 497, 5, 
	130, 66, 2, 496, 494, 3, 2, 2, 2, 496, 497, 3, 2, 2, 2, 497, 83, 3, 2, 
	2, 2, 498, 501, 7, 23, 2, 2, 499, 500, 7, 41, 2, 2, 500, 502, 5, 130, 66, 
	2, 501, 499, 3, 2, 2, 2, 501, 502, 3, 2, 2, 2, 502, 85, 3, 2, 2, 2, 503, 
	506, 7, 14, 2, 2, 504, 505, 7, 41, 2, 2, 505, 507, 5, 132, 67, 2, 506, 
	504, 3, 2, 2, 2, 506, 507, 3, 2, 2, 2, 507, 87, 3, 2, 2, 2, 508, 511, 7, 
	15, 2, 2, 509, 510, 7, 41, 2, 2, 510, 512, 5, 130, 66, 2, 511, 509, 3, 
	2, 2, 2, 511, 512, 3, 2, 2, 2, 512, 89, 3, 2, 2, 2, 513, 514, 7, 16, 2, 
	2, 514, 515, 7, 41, 2, 2, 515, 520, 5, 96, 49, 2, 516, 518, 7, 41, 2, 2, 
	517, 516, 3, 2, 2, 2, 517, 518, 3, 2, 2, 2, 518, 519, 3, 2, 2, 2, 519, 
	521, 7, 46, 2, 2, 520, 517, 3, 2, 2, 2, 520, 521, 3, 2, 2, 2, 521, 526, 
	3, 2, 2, 2, 522, 524, 7, 41, 2, 2, 523, 522, 3, 2, 2, 2, 523, 524, 3, 2, 
	2, 2, 524, 525, 3, 2, 2, 2, 525, 527, 5, 98, 50, 2, 526, 523, 3, 2, 2, 
	2, 526, 527, 3, 2, 2, 2, 527, 91, 3, 2, 2, 2, 528, 529, 7, 17, 2, 2, 529, 
	530, 5, 94, 48, 2, 530, 531, 7, 41, 2, 2, 531, 543, 5, 96, 49, 2, 532, 
	534, 7, 41, 2, 2, 533, 532, 3, 2, 2, 2, 533, 534, 3, 2, 2, 2, 534, 535, 
	3, 2, 2, 2, 535, 536, 7, 46, 2, 2, 536, 541, 3, 2, 2, 2, 537, 539, 7, 41, 
	2, 2, 538, 537, 3, 2, 2, 2, 538, 539, 3, 2, 2, 2, 539, 540, 3, 2, 2, 2, 
	540, 542, 5, 98, 50, 2, 541, 538, 3, 2, 2, 2, 541, 542, 3, 2, 2, 2, 542, 
	544, 3, 2, 2, 2, 543, 533, 3, 2, 2, 2, 543, 544, 3, 2, 2, 2, 544, 93, 3, 
	2, 2, 2, 545, 546, 7, 41, 2, 2, 546, 548, 5, 132, 67, 2, 547, 545, 3, 2, 
	2, 2, 547, 548, 3, 2, 2, 2, 548, 95, 3, 2, 2, 2, 549, 550, 7, 45, 2, 2, 
	550, 97, 3, 2, 2, 2, 551, 558, 7, 45, 2, 2, 552, 554, 7, 41, 2, 2, 553, 
	552, 3, 2, 2, 2, 553, 554, 3, 2, 2, 2, 554, 555, 3, 2, 2, 2, 555, 557, 
	7, 45, 2, 2, 556, 553, 3, 2, 2, 2, 557, 560, 3, 2, 2, 2, 558, 556, 3, 2, 
	2, 2, 558, 559, 3, 2, 2, 2, 559, 99, 3, 2, 2, 2, 560, 558, 3, 2, 2, 2, 
	561, 575, 7, 18, 2, 2, 562, 563, 7, 41, 2, 2, 563, 565, 5, 102, 52, 2, 
	564, 566, 7, 41, 2, 2, 565, 564, 3, 2, 2, 2, 565, 566, 3, 2, 2, 2, 566, 
	567, 3, 2, 2, 2, 567, 569, 7, 46, 2, 2, 568, 570, 7, 41, 2, 2, 569, 568, 
	3, 2, 2, 2, 569, 570, 3, 2, 2, 2, 570, 571, 3, 2, 2, 2, 571, 572, 5, 104, 
	53, 2, 572, 574, 3, 2, 2, 2, 573, 562, 3, 2, 2, 2, 574, 577, 3, 2, 2, 2, 
	575, 573, 3, 2, 2, 2, 575, 576, 3, 2, 2, 2, 576, 101, 3, 2, 2, 2, 577, 
	575, 3, 2, 2, 2, 578, 579, 7, 45, 2, 2, 579, 103, 3, 2, 2, 2, 580, 581, 
	7, 45, 2, 2, 581, 105, 3, 2, 2, 2, 582, 585, 7, 24, 2, 2, 583, 584, 7, 
	41, 2, 2, 584, 586, 5, 132, 67, 2, 585, 583, 3, 2, 2, 2, 585, 586, 3, 2, 
	2, 2, 586, 107, 3, 2, 2, 2, 587, 590, 7, 25, 2, 2, 588, 589, 7, 41, 2, 
	2, 589, 591, 5, 132, 67, 2, 590, 588, 3, 2, 2, 2, 590, 591, 3, 2, 2, 2, 
	591, 109, 3, 2, 2, 2, 592, 595, 7, 26, 2, 2, 593, 594, 7, 41, 2, 2, 594, 
	596, 5, 132, 67, 2, 595, 593, 3, 2, 2, 2, 595, 596, 3, 2, 2, 2, 596, 111, 
	3, 2, 2, 2, 597, 600, 7, 27, 2, 2, 598, 599, 7, 41, 2, 2, 599, 601, 5, 
	132, 67, 2, 600, 598, 3, 2, 2, 2, 600, 601, 3, 2, 2, 2, 601, 113, 3, 2, 
	2, 2, 602, 605, 7, 28, 2, 2, 603, 604, 7, 41, 2, 2, 604, 606, 5, 132, 67, 
	2, 605, 603, 3, 2, 2, 2, 605, 606, 3, 2, 2, 2, 606, 115, 3, 2, 2, 2, 607, 
	610, 7, 29, 2, 2, 608, 609, 7, 41, 2, 2, 609, 611, 5, 132, 67, 2, 610, 
	608, 3, 2, 2, 2, 610, 611, 3, 2, 2, 2, 611, 117, 3, 2, 2, 2, 612, 615, 
	7, 31, 2, 2, 613, 614, 7, 41, 2, 2, 614, 616, 5, 132, 67, 2, 615, 613, 
	3, 2, 2, 2, 615, 616, 3, 2, 2, 2, 616, 119, 3, 2, 2, 2, 617, 620, 7, 30, 
	2, 2, 618, 619, 7, 41, 2, 2, 619, 621, 5, 132, 67, 2, 620, 618, 3, 2, 2, 
	2, 620, 621, 3, 2, 2, 2, 621, 121, 3, 2, 2, 2, 622, 625, 7, 32, 2, 2, 623, 
	624, 7, 41, 2, 2, 624, 626, 5, 132, 67, 2, 625, 623, 3, 2, 2, 2, 625, 626, 
	3, 2, 2, 2, 626, 123, 3, 2, 2, 2, 627, 630, 7, 34, 2, 2, 628, 629, 7, 41, 
	2, 2, 629, 631, 5, 132, 67, 2, 630, 628, 3, 2, 2, 2, 630, 631, 3, 2, 2, 
	2, 631, 125, 3, 2, 2, 2, 632, 635, 7, 35, 2, 2, 633, 634, 7, 41, 2, 2, 
	634, 636, 5, 132, 67, 2, 635, 633, 3, 2, 2, 2, 635, 636, 3, 2, 2, 2, 636, 
	127, 3, 2, 2, 2, 637, 638, 5, 130, 66, 2, 638, 129, 3, 2, 2, 2, 639, 640, 
	5, 132, 67, 2, 640, 131, 3, 2, 2, 2, 641, 646, 5, 134, 68, 2, 642, 643, 
	7, 41, 2, 2, 643, 645, 5, 134, 68, 2, 644, 642, 3, 2, 2, 2, 645, 648, 3, 
	2, 2, 2, 646, 644, 3, 2, 2, 2, 646, 647, 3, 2, 2, 2, 647, 133, 3, 2, 2, 
	2, 648, 646, 3, 2, 2, 2, 649, 650, 7, 45, 2, 2, 650, 135, 3, 2, 2, 2, 95, 
	139, 143, 148, 153, 157, 162, 171, 177, 182, 188, 191, 198, 202, 210, 213, 
	220, 224, 229, 235, 238, 243, 250, 279, 287, 293, 296, 299, 304, 307, 322, 
	328, 331, 336, 342, 345, 348, 353, 356, 366, 369, 372, 382, 385, 388, 396, 
	399, 402, 414, 417, 427, 430, 433, 442, 447, 452, 457, 461, 466, 471, 476, 
	481, 486, 491, 496, 501, 506, 511, 517, 520, 523, 526, 533, 538, 541, 543, 
	547, 553, 558, 565, 569, 575, 585, 590, 595, 600, 605, 610, 615, 620, 625, 
	630, 635, 646,
}
var literalNames = []string{
	"", "", "", "", "", "'FROM'", "'FROM DOCKERFILE'", "'LOCALLY'", "'COPY'", 
	"'SAVE ARTIFACT'", "'SAVE IMAGE'", "'RUN'", "'EXPOSE'", "'VOLUME'", "'ENV'", 
	"'ARG'", "'LABEL'", "'BUILD'", "'WORKDIR'", "'USER'", "'CMD'", "'ENTRYPOINT'", 
	"'GIT CLONE'", "'ADD'", "'STOPSIGNAL'", "'ONBUILD'", "'HEALTHCHECK'", "'SHELL'", 
	"'DO'", "'COMMAND'", "'IMPORT'", "'VERSION'", "'CACHE'", "'HOST'", "'WITH'", 
	"", "", "", "", "", "'ELSE'", "'ELSE IF'", "'END'",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "Target", "UserCommand", "FROM", "FROM_DOCKERFILE", 
	"LOCALLY", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE", "RUN", "EXPOSE", "VOLUME", 
	"ENV", "ARG", "LABEL", "BUILD", "WORKDIR", "USER", "CMD", "ENTRYPOINT", 
	"GIT_CLONE", "ADD", "STOPSIGNAL", "ONBUILD", "HEALTHCHECK", "SHELL", "DO", 
	"COMMAND", "IMPORT", "VERSION", "CACHE", "HOST", "WITH", "DOCKER", "IF", 
	"FOR", "NL", "WS", "ELSE", "ELSE_IF", "END", "Atom", "EQUALS",
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
	"shellStmt", "userCommandStmt", "doStmt", "importStmt", "cacheStmt", "hostStmt", 
	"expr", "stmtWordsMaybeJSON", "stmtWords", "stmtWord",
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
	EarthParserWITH = 34
	EarthParserDOCKER = 35
	EarthParserIF = 36
	EarthParserFOR = 37
	EarthParserNL = 38
	EarthParserWS = 39
	EarthParserELSE = 40
	EarthParserELSE_IF = 41
	EarthParserEND = 42
	EarthParserAtom = 43
	EarthParserEQUALS = 44
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
	EarthParserRULE_forStmt = 25
	EarthParserRULE_forClause = 26
	EarthParserRULE_forBlock = 27
	EarthParserRULE_forExpr = 28
	EarthParserRULE_fromStmt = 29
	EarthParserRULE_fromDockerfileStmt = 30
	EarthParserRULE_locallyStmt = 31
	EarthParserRULE_copyStmt = 32
	EarthParserRULE_saveStmt = 33
	EarthParserRULE_saveImage = 34
	EarthParserRULE_saveArtifact = 35
	EarthParserRULE_runStmt = 36
	EarthParserRULE_buildStmt = 37
	EarthParserRULE_workdirStmt = 38
	EarthParserRULE_userStmt = 39
	EarthParserRULE_cmdStmt = 40
	EarthParserRULE_entrypointStmt = 41
	EarthParserRULE_exposeStmt = 42
	EarthParserRULE_volumeStmt = 43
	EarthParserRULE_envStmt = 44
	EarthParserRULE_argStmt = 45
	EarthParserRULE_optionalFlag = 46
	EarthParserRULE_envArgKey = 47
	EarthParserRULE_envArgValue = 48
	EarthParserRULE_labelStmt = 49
	EarthParserRULE_labelKey = 50
	EarthParserRULE_labelValue = 51
	EarthParserRULE_gitCloneStmt = 52
	EarthParserRULE_addStmt = 53
	EarthParserRULE_stopsignalStmt = 54
	EarthParserRULE_onbuildStmt = 55
	EarthParserRULE_healthcheckStmt = 56
	EarthParserRULE_shellStmt = 57
	EarthParserRULE_userCommandStmt = 58
	EarthParserRULE_doStmt = 59
	EarthParserRULE_importStmt = 60
	EarthParserRULE_cacheStmt = 61
	EarthParserRULE_hostStmt = 62
	EarthParserRULE_expr = 63
	EarthParserRULE_stmtWordsMaybeJSON = 64
	EarthParserRULE_stmtWords = 65
	EarthParserRULE_stmtWord = 66
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
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(134)
				p.Match(EarthParserNL)
			}


		}
		p.SetState(139)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserVERSION {
		{
			p.SetState(140)
			p.Version()
		}

	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << EarthParserFROM) | (1 << EarthParserFROM_DOCKERFILE) | (1 << EarthParserLOCALLY) | (1 << EarthParserCOPY) | (1 << EarthParserSAVE_ARTIFACT) | (1 << EarthParserSAVE_IMAGE) | (1 << EarthParserRUN) | (1 << EarthParserEXPOSE) | (1 << EarthParserVOLUME) | (1 << EarthParserENV) | (1 << EarthParserARG) | (1 << EarthParserLABEL) | (1 << EarthParserBUILD) | (1 << EarthParserWORKDIR) | (1 << EarthParserUSER) | (1 << EarthParserCMD) | (1 << EarthParserENTRYPOINT) | (1 << EarthParserGIT_CLONE) | (1 << EarthParserADD) | (1 << EarthParserSTOPSIGNAL) | (1 << EarthParserONBUILD) | (1 << EarthParserHEALTHCHECK) | (1 << EarthParserSHELL) | (1 << EarthParserDO) | (1 << EarthParserCOMMAND) | (1 << EarthParserIMPORT))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (EarthParserCACHE - 32)) | (1 << (EarthParserHOST - 32)) | (1 << (EarthParserWITH - 32)) | (1 << (EarthParserIF - 32)) | (1 << (EarthParserFOR - 32)) | (1 << (EarthParserWS - 32)))) != 0) {
		{
			p.SetState(143)
			p.Stmts()
		}
		{
			p.SetState(144)
			p.Match(EarthParserNL)
		}

	}
	p.SetState(151)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(148)
				p.Match(EarthParserNL)
			}


		}
		p.SetState(153)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}
	p.SetState(155)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserTarget || _la == EarthParserUserCommand {
		{
			p.SetState(154)
			p.Targets()
		}

	}
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserNL {
		{
			p.SetState(157)
			p.Match(EarthParserNL)
		}


		p.SetState(162)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(163)
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
		p.SetState(165)
		p.TargetOrUserCommand()
	}
	p.SetState(175)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(169)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for _la == EarthParserNL {
				{
					p.SetState(166)
					p.Match(EarthParserNL)
				}


				p.SetState(171)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(172)
				p.TargetOrUserCommand()
			}


		}
		p.SetState(177)
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

	p.SetState(180)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserTarget:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(178)
			p.Target()
		}


	case EarthParserUserCommand:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(179)
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
		p.SetState(182)
		p.TargetHeader()
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(183)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(186)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}
	p.SetState(189)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(188)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(200)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserINDENT {
		{
			p.SetState(191)
			p.Match(EarthParserINDENT)
		}
		{
			p.SetState(192)
			p.Stmts()
		}
		p.SetState(194)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(193)
				p.Match(EarthParserNL)
			}


			p.SetState(196)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(198)
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
		p.SetState(202)
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
		p.SetState(204)
		p.UserCommandHeader()
	}
	p.SetState(206)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(205)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(208)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext())
	}
	p.SetState(211)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(210)
			p.Match(EarthParserWS)
		}

	}
	p.SetState(222)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserINDENT {
		{
			p.SetState(213)
			p.Match(EarthParserINDENT)
		}
		{
			p.SetState(214)
			p.Stmts()
		}
		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(215)
				p.Match(EarthParserNL)
			}


			p.SetState(218)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(220)
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
		p.SetState(224)
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
	p.SetState(227)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(226)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(229)
		p.Stmt()
	}
	p.SetState(241)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(231)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(230)
					p.Match(EarthParserNL)
				}


				p.SetState(233)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			p.SetState(236)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			if _la == EarthParserWS {
				{
					p.SetState(235)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(238)
				p.Stmt()
			}


		}
		p.SetState(243)
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

	p.SetState(248)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM, EarthParserFROM_DOCKERFILE, EarthParserLOCALLY, EarthParserCOPY, EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE, EarthParserRUN, EarthParserEXPOSE, EarthParserVOLUME, EarthParserENV, EarthParserARG, EarthParserLABEL, EarthParserBUILD, EarthParserWORKDIR, EarthParserUSER, EarthParserCMD, EarthParserENTRYPOINT, EarthParserGIT_CLONE, EarthParserADD, EarthParserSTOPSIGNAL, EarthParserONBUILD, EarthParserHEALTHCHECK, EarthParserSHELL, EarthParserDO, EarthParserCOMMAND, EarthParserIMPORT, EarthParserCACHE, EarthParserHOST:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(244)
			p.CommandStmt()
		}


	case EarthParserWITH:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(245)
			p.WithStmt()
		}


	case EarthParserIF:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(246)
			p.IfStmt()
		}


	case EarthParserFOR:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(247)
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

func (s *CommandStmtContext) CacheStmt() ICacheStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICacheStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICacheStmtContext)
}

func (s *CommandStmtContext) HostStmt() IHostStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHostStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHostStmtContext)
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

	p.SetState(277)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserFROM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(250)
			p.FromStmt()
		}


	case EarthParserFROM_DOCKERFILE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(251)
			p.FromDockerfileStmt()
		}


	case EarthParserLOCALLY:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(252)
			p.LocallyStmt()
		}


	case EarthParserCOPY:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(253)
			p.CopyStmt()
		}


	case EarthParserSAVE_ARTIFACT, EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(254)
			p.SaveStmt()
		}


	case EarthParserRUN:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(255)
			p.RunStmt()
		}


	case EarthParserBUILD:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(256)
			p.BuildStmt()
		}


	case EarthParserWORKDIR:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(257)
			p.WorkdirStmt()
		}


	case EarthParserUSER:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(258)
			p.UserStmt()
		}


	case EarthParserCMD:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(259)
			p.CmdStmt()
		}


	case EarthParserENTRYPOINT:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(260)
			p.EntrypointStmt()
		}


	case EarthParserEXPOSE:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(261)
			p.ExposeStmt()
		}


	case EarthParserVOLUME:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(262)
			p.VolumeStmt()
		}


	case EarthParserENV:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(263)
			p.EnvStmt()
		}


	case EarthParserARG:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(264)
			p.ArgStmt()
		}


	case EarthParserLABEL:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(265)
			p.LabelStmt()
		}


	case EarthParserGIT_CLONE:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(266)
			p.GitCloneStmt()
		}


	case EarthParserADD:
		p.EnterOuterAlt(localctx, 18)
		{
			p.SetState(267)
			p.AddStmt()
		}


	case EarthParserSTOPSIGNAL:
		p.EnterOuterAlt(localctx, 19)
		{
			p.SetState(268)
			p.StopsignalStmt()
		}


	case EarthParserONBUILD:
		p.EnterOuterAlt(localctx, 20)
		{
			p.SetState(269)
			p.OnbuildStmt()
		}


	case EarthParserHEALTHCHECK:
		p.EnterOuterAlt(localctx, 21)
		{
			p.SetState(270)
			p.HealthcheckStmt()
		}


	case EarthParserSHELL:
		p.EnterOuterAlt(localctx, 22)
		{
			p.SetState(271)
			p.ShellStmt()
		}


	case EarthParserCOMMAND:
		p.EnterOuterAlt(localctx, 23)
		{
			p.SetState(272)
			p.UserCommandStmt()
		}


	case EarthParserDO:
		p.EnterOuterAlt(localctx, 24)
		{
			p.SetState(273)
			p.DoStmt()
		}


	case EarthParserIMPORT:
		p.EnterOuterAlt(localctx, 25)
		{
			p.SetState(274)
			p.ImportStmt()
		}


	case EarthParserCACHE:
		p.EnterOuterAlt(localctx, 26)
		{
			p.SetState(275)
			p.CacheStmt()
		}


	case EarthParserHOST:
		p.EnterOuterAlt(localctx, 27)
		{
			p.SetState(276)
			p.HostStmt()
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

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(279)
		p.Match(EarthParserVERSION)
	}
	{
		p.SetState(280)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(281)
		p.StmtWords()
	}
	p.SetState(283)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
				{
					p.SetState(282)
					p.Match(EarthParserNL)
				}




		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(285)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext())
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
		p.SetState(287)
		p.WithExpr()
	}
	p.SetState(297)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(288)
				p.Match(EarthParserNL)
			}


			p.SetState(291)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(294)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(293)
				p.Match(EarthParserWS)
			}


		}
		{
			p.SetState(296)
			p.WithBlock()
		}


	}
	p.SetState(300)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(299)
			p.Match(EarthParserNL)
		}


		p.SetState(302)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(305)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(304)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(307)
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
		p.SetState(309)
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
		p.SetState(311)
		p.Match(EarthParserWITH)
	}
	{
		p.SetState(312)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(313)
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
		p.SetState(315)
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
		p.SetState(317)
		p.Match(EarthParserDOCKER)
	}
	p.SetState(320)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(318)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(319)
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
		p.SetState(322)
		p.IfClause()
	}
	p.SetState(334)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(324)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			for ok := true; ok; ok = _la == EarthParserNL {
				{
					p.SetState(323)
					p.Match(EarthParserNL)
				}


				p.SetState(326)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			p.SetState(329)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			if _la == EarthParserWS {
				{
					p.SetState(328)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(331)
				p.ElseIfClause()
			}


		}
		p.SetState(336)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext())
	}
	p.SetState(346)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext()) == 1 {
		p.SetState(338)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(337)
				p.Match(EarthParserNL)
			}


			p.SetState(340)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
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
			p.ElseClause()
		}


	}
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
	p.SetState(354)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(353)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(356)
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
		p.SetState(358)
		p.Match(EarthParserIF)
	}
	{
		p.SetState(359)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(360)
		p.IfExpr()
	}
	p.SetState(370)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) == 1 {
		p.SetState(362)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(361)
				p.Match(EarthParserNL)
			}


			p.SetState(364)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(367)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(366)
				p.Match(EarthParserWS)
			}


		}
		{
			p.SetState(369)
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
		p.SetState(372)
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
		p.SetState(374)
		p.Match(EarthParserELSE_IF)
	}
	{
		p.SetState(375)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(376)
		p.ElseIfExpr()
	}
	p.SetState(386)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) == 1 {
		p.SetState(378)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for ok := true; ok; ok = _la == EarthParserNL {
			{
				p.SetState(377)
				p.Match(EarthParserNL)
			}


			p.SetState(380)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		p.SetState(383)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(382)
				p.Match(EarthParserWS)
			}


		}
		{
			p.SetState(385)
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
		p.SetState(388)
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
		p.SetState(390)
		p.Match(EarthParserELSE)
	}
	p.SetState(400)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) == 1 {
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
		p.SetState(397)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(396)
				p.Match(EarthParserWS)
			}


		}
		{
			p.SetState(399)
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
		p.SetState(402)
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
		p.SetState(404)
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
		p.SetState(406)
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
		p.SetState(408)
		p.ForClause()
	}
	p.SetState(410)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for ok := true; ok; ok = _la == EarthParserNL {
		{
			p.SetState(409)
			p.Match(EarthParserNL)
		}


		p.SetState(412)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(415)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(414)
			p.Match(EarthParserWS)
		}

	}
	{
		p.SetState(417)
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
		p.SetState(419)
		p.Match(EarthParserFOR)
	}
	{
		p.SetState(420)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(421)
		p.ForExpr()
	}
	p.SetState(431)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) == 1 {
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
		p.SetState(428)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 50, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(427)
				p.Match(EarthParserWS)
			}


		}
		{
			p.SetState(430)
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
		p.SetState(433)
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
		p.SetState(435)
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
		p.SetState(437)
		p.Match(EarthParserFROM)
	}
	p.SetState(440)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(438)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(439)
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
		p.SetState(442)
		p.Match(EarthParserFROM_DOCKERFILE)
	}
	p.SetState(445)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(443)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(444)
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
		p.SetState(447)
		p.Match(EarthParserLOCALLY)
	}
	p.SetState(450)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(448)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(449)
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
		p.SetState(452)
		p.Match(EarthParserCOPY)
	}
	p.SetState(455)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(453)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(454)
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

	p.SetState(459)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case EarthParserSAVE_ARTIFACT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(457)
			p.SaveArtifact()
		}


	case EarthParserSAVE_IMAGE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(458)
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
		p.SetState(461)
		p.Match(EarthParserSAVE_IMAGE)
	}
	p.SetState(464)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(462)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(463)
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
		p.SetState(466)
		p.Match(EarthParserSAVE_ARTIFACT)
	}
	p.SetState(469)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(467)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(468)
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
		p.SetState(471)
		p.Match(EarthParserRUN)
	}
	p.SetState(474)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(472)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(473)
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
		p.SetState(476)
		p.Match(EarthParserBUILD)
	}
	p.SetState(479)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(477)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(478)
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
		p.SetState(481)
		p.Match(EarthParserWORKDIR)
	}
	p.SetState(484)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(482)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(483)
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
		p.SetState(486)
		p.Match(EarthParserUSER)
	}
	p.SetState(489)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(487)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(488)
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
		p.SetState(491)
		p.Match(EarthParserCMD)
	}
	p.SetState(494)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(492)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(493)
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
		p.SetState(496)
		p.Match(EarthParserENTRYPOINT)
	}
	p.SetState(499)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(497)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(498)
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
		p.SetState(501)
		p.Match(EarthParserEXPOSE)
	}
	p.SetState(504)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(502)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(503)
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
		p.SetState(506)
		p.Match(EarthParserVOLUME)
	}
	p.SetState(509)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(507)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(508)
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
		p.SetState(511)
		p.Match(EarthParserENV)
	}
	{
		p.SetState(512)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(513)
		p.EnvArgKey()
	}
	p.SetState(518)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 68, p.GetParserRuleContext()) == 1 {
		p.SetState(515)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(514)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(517)
			p.Match(EarthParserEQUALS)
		}


	}
	p.SetState(524)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(521)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(520)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(523)
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
		p.SetState(526)
		p.Match(EarthParserARG)
	}
	{
		p.SetState(527)
		p.OptionalFlag()
	}
	{
		p.SetState(528)
		p.Match(EarthParserWS)
	}
	{
		p.SetState(529)
		p.EnvArgKey()
	}
	p.SetState(541)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS || _la == EarthParserEQUALS {
		p.SetState(531)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(530)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(533)
			p.Match(EarthParserEQUALS)
		}

		p.SetState(539)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS || _la == EarthParserAtom {
			p.SetState(536)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)


			if _la == EarthParserWS {
				{
					p.SetState(535)
					p.Match(EarthParserWS)
				}

			}
			{
				p.SetState(538)
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
	p.SetState(545)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 75, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(543)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(544)
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
		p.SetState(547)
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
		p.SetState(549)
		p.Match(EarthParserAtom)
	}
	p.SetState(556)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserWS || _la == EarthParserAtom {
		p.SetState(551)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(550)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(553)
			p.Match(EarthParserAtom)
		}


		p.SetState(558)
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
		p.SetState(559)
		p.Match(EarthParserLABEL)
	}
	p.SetState(573)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == EarthParserWS {
		{
			p.SetState(560)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(561)
			p.LabelKey()
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
			p.Match(EarthParserEQUALS)
		}
		p.SetState(567)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == EarthParserWS {
			{
				p.SetState(566)
				p.Match(EarthParserWS)
			}

		}
		{
			p.SetState(569)
			p.LabelValue()
		}


		p.SetState(575)
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
		p.SetState(576)
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
		p.SetState(578)
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
		p.SetState(580)
		p.Match(EarthParserGIT_CLONE)
	}
	p.SetState(583)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(581)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(582)
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
		p.SetState(585)
		p.Match(EarthParserADD)
	}
	p.SetState(588)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(586)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(587)
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
		p.SetState(590)
		p.Match(EarthParserSTOPSIGNAL)
	}
	p.SetState(593)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(591)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(592)
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
		p.SetState(595)
		p.Match(EarthParserONBUILD)
	}
	p.SetState(598)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(596)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(597)
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
		p.SetState(600)
		p.Match(EarthParserHEALTHCHECK)
	}
	p.SetState(603)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(601)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(602)
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
		p.SetState(605)
		p.Match(EarthParserSHELL)
	}
	p.SetState(608)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(606)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(607)
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
		p.SetState(610)
		p.Match(EarthParserCOMMAND)
	}
	p.SetState(613)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(611)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(612)
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
		p.SetState(615)
		p.Match(EarthParserDO)
	}
	p.SetState(618)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(616)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(617)
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
		p.SetState(620)
		p.Match(EarthParserIMPORT)
	}
	p.SetState(623)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(621)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(622)
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

func (s *CacheStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *CacheStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewCacheStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 122, EarthParserRULE_cacheStmt)
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
		p.Match(EarthParserCACHE)
	}
	p.SetState(628)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(626)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(627)
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

func (s *HostStmtContext) WS() antlr.TerminalNode {
	return s.GetToken(EarthParserWS, 0)
}

func (s *HostStmtContext) StmtWords() IStmtWordsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	localctx = NewHostStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, EarthParserRULE_hostStmt)
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
		p.SetState(630)
		p.Match(EarthParserHOST)
	}
	p.SetState(633)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == EarthParserWS {
		{
			p.SetState(631)
			p.Match(EarthParserWS)
		}
		{
			p.SetState(632)
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
	p.EnterRule(localctx, 126, EarthParserRULE_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.EnterRule(localctx, 128, EarthParserRULE_stmtWordsMaybeJSON)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
	p.EnterRule(localctx, 130, EarthParserRULE_stmtWords)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.SetState(639)
		p.StmtWord()
	}
	p.SetState(644)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 92, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(640)
				p.Match(EarthParserWS)
			}
			{
				p.SetState(641)
				p.StmtWord()
			}


		}
		p.SetState(646)
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
	p.EnterRule(localctx, 132, EarthParserRULE_stmtWord)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
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
		p.Match(EarthParserAtom)
	}



	return localctx
}


