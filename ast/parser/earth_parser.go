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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 54, 698, 
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
	4, 66, 9, 66, 4, 67, 9, 67, 4, 68, 9, 68, 4, 69, 9, 69, 4, 70, 9, 70, 4, 
	71, 9, 71, 4, 72, 9, 72, 4, 73, 9, 73, 4, 74, 9, 74, 4, 75, 9, 75, 4, 76, 
	9, 76, 4, 77, 9, 77, 4, 78, 9, 78, 4, 79, 9, 79, 4, 80, 9, 80, 4, 81, 9, 
	81, 4, 82, 9, 82, 3, 2, 7, 2, 166, 10, 2, 12, 2, 14, 2, 169, 11, 2, 3, 
	2, 5, 2, 172, 10, 2, 3, 2, 3, 2, 3, 2, 5, 2, 177, 10, 2, 3, 2, 7, 2, 180, 
	10, 2, 12, 2, 14, 2, 183, 11, 2, 3, 2, 5, 2, 186, 10, 2, 3, 2, 7, 2, 189, 
	10, 2, 12, 2, 14, 2, 192, 11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 7, 3, 198, 10, 
	3, 12, 3, 14, 3, 201, 11, 3, 3, 3, 7, 3, 204, 10, 3, 12, 3, 14, 3, 207, 
	11, 3, 3, 4, 3, 4, 5, 4, 211, 10, 4, 3, 5, 3, 5, 6, 5, 215, 10, 5, 13, 
	5, 14, 5, 216, 3, 5, 3, 5, 7, 5, 221, 10, 5, 12, 5, 14, 5, 224, 11, 5, 
	3, 5, 5, 5, 227, 10, 5, 3, 5, 6, 5, 230, 10, 5, 13, 5, 14, 5, 231, 3, 5, 
	5, 5, 235, 10, 5, 3, 6, 3, 6, 3, 7, 3, 7, 6, 7, 241, 10, 7, 13, 7, 14, 
	7, 242, 3, 7, 3, 7, 7, 7, 247, 10, 7, 12, 7, 14, 7, 250, 11, 7, 3, 7, 3, 
	7, 6, 7, 254, 10, 7, 13, 7, 14, 7, 255, 3, 7, 3, 7, 5, 7, 260, 10, 7, 3, 
	8, 3, 8, 3, 9, 3, 9, 6, 9, 266, 10, 9, 13, 9, 14, 9, 267, 3, 9, 7, 9, 271, 
	10, 9, 12, 9, 14, 9, 274, 11, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 
	10, 5, 10, 282, 10, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 
	3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 
	11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 
	3, 11, 3, 11, 5, 11, 314, 10, 11, 3, 12, 3, 12, 3, 12, 6, 12, 319, 10, 
	12, 13, 12, 14, 12, 320, 3, 13, 3, 13, 6, 13, 325, 10, 13, 13, 13, 14, 
	13, 326, 3, 13, 5, 13, 330, 10, 13, 3, 13, 6, 13, 333, 10, 13, 13, 13, 
	14, 13, 334, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15, 3, 16, 3, 
	16, 3, 17, 3, 17, 5, 17, 348, 10, 17, 3, 18, 3, 18, 6, 18, 352, 10, 18, 
	13, 18, 14, 18, 353, 3, 18, 7, 18, 357, 10, 18, 12, 18, 14, 18, 360, 11, 
	18, 3, 18, 6, 18, 363, 10, 18, 13, 18, 14, 18, 364, 3, 18, 5, 18, 368, 
	10, 18, 3, 18, 6, 18, 371, 10, 18, 13, 18, 14, 18, 372, 3, 18, 3, 18, 3, 
	19, 3, 19, 3, 19, 6, 19, 380, 10, 19, 13, 19, 14, 19, 381, 3, 19, 5, 19, 
	385, 10, 19, 3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 6, 21, 392, 10, 21, 13, 
	21, 14, 21, 393, 3, 21, 5, 21, 397, 10, 21, 3, 22, 3, 22, 3, 23, 3, 23, 
	6, 23, 403, 10, 23, 13, 23, 14, 23, 404, 3, 23, 5, 23, 408, 10, 23, 3, 
	24, 3, 24, 3, 25, 3, 25, 3, 26, 3, 26, 3, 27, 3, 27, 6, 27, 418, 10, 27, 
	13, 27, 14, 27, 419, 3, 27, 5, 27, 423, 10, 27, 3, 27, 6, 27, 426, 10, 
	27, 13, 27, 14, 27, 427, 3, 27, 5, 27, 431, 10, 27, 3, 27, 6, 27, 434, 
	10, 27, 13, 27, 14, 27, 435, 3, 27, 3, 27, 3, 28, 3, 28, 6, 28, 442, 10, 
	28, 13, 28, 14, 28, 443, 3, 28, 5, 28, 447, 10, 28, 3, 29, 3, 29, 3, 30, 
	3, 30, 6, 30, 453, 10, 30, 13, 30, 14, 30, 454, 3, 30, 5, 30, 458, 10, 
	30, 3, 31, 3, 31, 3, 32, 3, 32, 6, 32, 464, 10, 32, 13, 32, 14, 32, 465, 
	3, 32, 5, 32, 469, 10, 32, 3, 33, 3, 33, 3, 34, 3, 34, 6, 34, 475, 10, 
	34, 13, 34, 14, 34, 476, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 6, 35, 484, 
	10, 35, 13, 35, 14, 35, 485, 3, 35, 5, 35, 489, 10, 35, 3, 36, 3, 36, 3, 
	37, 3, 37, 3, 38, 3, 38, 6, 38, 497, 10, 38, 13, 38, 14, 38, 498, 3, 38, 
	3, 38, 3, 39, 3, 39, 5, 39, 505, 10, 39, 3, 39, 6, 39, 508, 10, 39, 13, 
	39, 14, 39, 509, 3, 39, 5, 39, 513, 10, 39, 3, 40, 3, 40, 3, 41, 3, 41, 
	3, 42, 3, 42, 5, 42, 521, 10, 42, 3, 43, 3, 43, 5, 43, 525, 10, 43, 3, 
	44, 3, 44, 5, 44, 529, 10, 44, 3, 45, 3, 45, 5, 45, 533, 10, 45, 3, 46, 
	3, 46, 5, 46, 537, 10, 46, 3, 47, 3, 47, 5, 47, 541, 10, 47, 3, 48, 3, 
	48, 5, 48, 545, 10, 48, 3, 49, 3, 49, 5, 49, 549, 10, 49, 3, 50, 3, 50, 
	5, 50, 553, 10, 50, 3, 51, 3, 51, 5, 51, 557, 10, 51, 3, 52, 3, 52, 5, 
	52, 561, 10, 52, 3, 53, 3, 53, 5, 53, 565, 10, 53, 3, 54, 3, 54, 5, 54, 
	569, 10, 54, 3, 55, 3, 55, 5, 55, 573, 10, 55, 3, 56, 3, 56, 5, 56, 577, 
	10, 56, 3, 57, 3, 57, 3, 57, 5, 57, 582, 10, 57, 3, 57, 5, 57, 585, 10, 
	57, 3, 57, 5, 57, 588, 10, 57, 3, 58, 3, 58, 3, 58, 3, 58, 3, 58, 5, 58, 
	595, 10, 58, 3, 58, 5, 58, 598, 10, 58, 5, 58, 600, 10, 58, 3, 59, 5, 59, 
	603, 10, 59, 3, 60, 3, 60, 3, 61, 3, 61, 5, 61, 609, 10, 61, 3, 61, 7, 
	61, 612, 10, 61, 12, 61, 14, 61, 615, 11, 61, 3, 62, 3, 62, 3, 62, 3, 62, 
	3, 62, 7, 62, 622, 10, 62, 12, 62, 14, 62, 625, 11, 62, 3, 63, 3, 63, 3, 
	64, 3, 64, 3, 65, 3, 65, 5, 65, 633, 10, 65, 3, 66, 3, 66, 5, 66, 637, 
	10, 66, 3, 67, 3, 67, 5, 67, 641, 10, 67, 3, 68, 3, 68, 5, 68, 645, 10, 
	68, 3, 69, 3, 69, 5, 69, 649, 10, 69, 3, 70, 3, 70, 5, 70, 653, 10, 70, 
	3, 71, 3, 71, 5, 71, 657, 10, 71, 3, 72, 3, 72, 5, 72, 661, 10, 72, 3, 
	73, 3, 73, 5, 73, 665, 10, 73, 3, 74, 3, 74, 5, 74, 669, 10, 74, 3, 75, 
	3, 75, 5, 75, 673, 10, 75, 3, 76, 3, 76, 5, 76, 677, 10, 76, 3, 77, 3, 
	77, 5, 77, 681, 10, 77, 3, 78, 3, 78, 5, 78, 685, 10, 78, 3, 79, 3, 79, 
	3, 80, 3, 80, 3, 81, 6, 81, 692, 10, 81, 13, 81, 14, 81, 693, 3, 82, 3, 
	82, 3, 82, 2, 2, 83, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 
	30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 
	66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 
	102, 104, 106, 108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 
	132, 134, 136, 138, 140, 142, 144, 146, 148, 150, 152, 154, 156, 158, 160, 
	162, 2, 2, 2, 744, 2, 167, 3, 2, 2, 2, 4, 195, 3, 2, 2, 2, 6, 210, 3, 2, 
	2, 2, 8, 212, 3, 2, 2, 2, 10, 236, 3, 2, 2, 2, 12, 238, 3, 2, 2, 2, 14, 
	261, 3, 2, 2, 2, 16, 263, 3, 2, 2, 2, 18, 281, 3, 2, 2, 2, 20, 313, 3, 
	2, 2, 2, 22, 315, 3, 2, 2, 2, 24, 322, 3, 2, 2, 2, 26, 338, 3, 2, 2, 2, 
	28, 340, 3, 2, 2, 2, 30, 343, 3, 2, 2, 2, 32, 345, 3, 2, 2, 2, 34, 349, 
	3, 2, 2, 2, 36, 376, 3, 2, 2, 2, 38, 386, 3, 2, 2, 2, 40, 388, 3, 2, 2, 
	2, 42, 398, 3, 2, 2, 2, 44, 400, 3, 2, 2, 2, 46, 409, 3, 2, 2, 2, 48, 411, 
	3, 2, 2, 2, 50, 413, 3, 2, 2, 2, 52, 415, 3, 2, 2, 2, 54, 439, 3, 2, 2, 
	2, 56, 448, 3, 2, 2, 2, 58, 450, 3, 2, 2, 2, 60, 459, 3, 2, 2, 2, 62, 461, 
	3, 2, 2, 2, 64, 470, 3, 2, 2, 2, 66, 472, 3, 2, 2, 2, 68, 480, 3, 2, 2, 
	2, 70, 490, 3, 2, 2, 2, 72, 492, 3, 2, 2, 2, 74, 494, 3, 2, 2, 2, 76, 502, 
	3, 2, 2, 2, 78, 514, 3, 2, 2, 2, 80, 516, 3, 2, 2, 2, 82, 518, 3, 2, 2, 
	2, 84, 522, 3, 2, 2, 2, 86, 526, 3, 2, 2, 2, 88, 530, 3, 2, 2, 2, 90, 536, 
	3, 2, 2, 2, 92, 538, 3, 2, 2, 2, 94, 542, 3, 2, 2, 2, 96, 546, 3, 2, 2, 
	2, 98, 550, 3, 2, 2, 2, 100, 554, 3, 2, 2, 2, 102, 558, 3, 2, 2, 2, 104, 
	562, 3, 2, 2, 2, 106, 566, 3, 2, 2, 2, 108, 570, 3, 2, 2, 2, 110, 574, 
	3, 2, 2, 2, 112, 578, 3, 2, 2, 2, 114, 589, 3, 2, 2, 2, 116, 602, 3, 2, 
	2, 2, 118, 604, 3, 2, 2, 2, 120, 606, 3, 2, 2, 2, 122, 616, 3, 2, 2, 2, 
	124, 626, 3, 2, 2, 2, 126, 628, 3, 2, 2, 2, 128, 630, 3, 2, 2, 2, 130, 
	634, 3, 2, 2, 2, 132, 638, 3, 2, 2, 2, 134, 642, 3, 2, 2, 2, 136, 646, 
	3, 2, 2, 2, 138, 650, 3, 2, 2, 2, 140, 654, 3, 2, 2, 2, 142, 658, 3, 2, 
	2, 2, 144, 662, 3, 2, 2, 2, 146, 666, 3, 2, 2, 2, 148, 670, 3, 2, 2, 2, 
	150, 674, 3, 2, 2, 2, 152, 678, 3, 2, 2, 2, 154, 682, 3, 2, 2, 2, 156, 
	686, 3, 2, 2, 2, 158, 688, 3, 2, 2, 2, 160, 691, 3, 2, 2, 2, 162, 695, 
	3, 2, 2, 2, 164, 166, 7, 45, 2, 2, 165, 164, 3, 2, 2, 2, 166, 169, 3, 2, 
	2, 2, 167, 165, 3, 2, 2, 2, 167, 168, 3, 2, 2, 2, 168, 171, 3, 2, 2, 2, 
	169, 167, 3, 2, 2, 2, 170, 172, 5, 22, 12, 2, 171, 170, 3, 2, 2, 2, 171, 
	172, 3, 2, 2, 2, 172, 176, 3, 2, 2, 2, 173, 174, 5, 16, 9, 2, 174, 175, 
	7, 45, 2, 2, 175, 177, 3, 2, 2, 2, 176, 173, 3, 2, 2, 2, 176, 177, 3, 2, 
	2, 2, 177, 181, 3, 2, 2, 2, 178, 180, 7, 45, 2, 2, 179, 178, 3, 2, 2, 2, 
	180, 183, 3, 2, 2, 2, 181, 179, 3, 2, 2, 2, 181, 182, 3, 2, 2, 2, 182, 
	185, 3, 2, 2, 2, 183, 181, 3, 2, 2, 2, 184, 186, 5, 4, 3, 2, 185, 184, 
	3, 2, 2, 2, 185, 186, 3, 2, 2, 2, 186, 190, 3, 2, 2, 2, 187, 189, 7, 45, 
	2, 2, 188, 187, 3, 2, 2, 2, 189, 192, 3, 2, 2, 2, 190, 188, 3, 2, 2, 2, 
	190, 191, 3, 2, 2, 2, 191, 193, 3, 2, 2, 2, 192, 190, 3, 2, 2, 2, 193, 
	194, 7, 2, 2, 3, 194, 3, 3, 2, 2, 2, 195, 205, 5, 6, 4, 2, 196, 198, 7, 
	45, 2, 2, 197, 196, 3, 2, 2, 2, 198, 201, 3, 2, 2, 2, 199, 197, 3, 2, 2, 
	2, 199, 200, 3, 2, 2, 2, 200, 202, 3, 2, 2, 2, 201, 199, 3, 2, 2, 2, 202, 
	204, 5, 6, 4, 2, 203, 199, 3, 2, 2, 2, 204, 207, 3, 2, 2, 2, 205, 203, 
	3, 2, 2, 2, 205, 206, 3, 2, 2, 2, 206, 5, 3, 2, 2, 2, 207, 205, 3, 2, 2, 
	2, 208, 211, 5, 8, 5, 2, 209, 211, 5, 12, 7, 2, 210, 208, 3, 2, 2, 2, 210, 
	209, 3, 2, 2, 2, 211, 7, 3, 2, 2, 2, 212, 214, 5, 10, 6, 2, 213, 215, 7, 
	45, 2, 2, 214, 213, 3, 2, 2, 2, 215, 216, 3, 2, 2, 2, 216, 214, 3, 2, 2, 
	2, 216, 217, 3, 2, 2, 2, 217, 234, 3, 2, 2, 2, 218, 222, 7, 3, 2, 2, 219, 
	221, 7, 45, 2, 2, 220, 219, 3, 2, 2, 2, 221, 224, 3, 2, 2, 2, 222, 220, 
	3, 2, 2, 2, 222, 223, 3, 2, 2, 2, 223, 226, 3, 2, 2, 2, 224, 222, 3, 2, 
	2, 2, 225, 227, 5, 16, 9, 2, 226, 225, 3, 2, 2, 2, 226, 227, 3, 2, 2, 2, 
	227, 229, 3, 2, 2, 2, 228, 230, 7, 45, 2, 2, 229, 228, 3, 2, 2, 2, 230, 
	231, 3, 2, 2, 2, 231, 229, 3, 2, 2, 2, 231, 232, 3, 2, 2, 2, 232, 233, 
	3, 2, 2, 2, 233, 235, 7, 4, 2, 2, 234, 218, 3, 2, 2, 2, 234, 235, 3, 2, 
	2, 2, 235, 9, 3, 2, 2, 2, 236, 237, 7, 5, 2, 2, 237, 11, 3, 2, 2, 2, 238, 
	240, 5, 14, 8, 2, 239, 241, 7, 45, 2, 2, 240, 239, 3, 2, 2, 2, 241, 242, 
	3, 2, 2, 2, 242, 240, 3, 2, 2, 2, 242, 243, 3, 2, 2, 2, 243, 259, 3, 2, 
	2, 2, 244, 248, 7, 3, 2, 2, 245, 247, 7, 45, 2, 2, 246, 245, 3, 2, 2, 2, 
	247, 250, 3, 2, 2, 2, 248, 246, 3, 2, 2, 2, 248, 249, 3, 2, 2, 2, 249, 
	251, 3, 2, 2, 2, 250, 248, 3, 2, 2, 2, 251, 253, 5, 16, 9, 2, 252, 254, 
	7, 45, 2, 2, 253, 252, 3, 2, 2, 2, 254, 255, 3, 2, 2, 2, 255, 253, 3, 2, 
	2, 2, 255, 256, 3, 2, 2, 2, 256, 257, 3, 2, 2, 2, 257, 258, 7, 4, 2, 2, 
	258, 260, 3, 2, 2, 2, 259, 244, 3, 2, 2, 2, 259, 260, 3, 2, 2, 2, 260, 
	13, 3, 2, 2, 2, 261, 262, 7, 6, 2, 2, 262, 15, 3, 2, 2, 2, 263, 272, 5, 
	18, 10, 2, 264, 266, 7, 45, 2, 2, 265, 264, 3, 2, 2, 2, 266, 267, 3, 2, 
	2, 2, 267, 265, 3, 2, 2, 2, 267, 268, 3, 2, 2, 2, 268, 269, 3, 2, 2, 2, 
	269, 271, 5, 18, 10, 2, 270, 265, 3, 2, 2, 2, 271, 274, 3, 2, 2, 2, 272, 
	270, 3, 2, 2, 2, 272, 273, 3, 2, 2, 2, 273, 17, 3, 2, 2, 2, 274, 272, 3, 
	2, 2, 2, 275, 282, 5, 20, 11, 2, 276, 282, 5, 24, 13, 2, 277, 282, 5, 34, 
	18, 2, 278, 282, 5, 66, 34, 2, 279, 282, 5, 74, 38, 2, 280, 282, 5, 52, 
	27, 2, 281, 275, 3, 2, 2, 2, 281, 276, 3, 2, 2, 2, 281, 277, 3, 2, 2, 2, 
	281, 278, 3, 2, 2, 2, 281, 279, 3, 2, 2, 2, 281, 280, 3, 2, 2, 2, 282, 
	19, 3, 2, 2, 2, 283, 314, 5, 82, 42, 2, 284, 314, 5, 84, 43, 2, 285, 314, 
	5, 86, 44, 2, 286, 314, 5, 88, 45, 2, 287, 314, 5, 90, 46, 2, 288, 314, 
	5, 96, 49, 2, 289, 314, 5, 98, 50, 2, 290, 314, 5, 100, 51, 2, 291, 314, 
	5, 102, 52, 2, 292, 314, 5, 104, 53, 2, 293, 314, 5, 106, 54, 2, 294, 314, 
	5, 108, 55, 2, 295, 314, 5, 110, 56, 2, 296, 314, 5, 112, 57, 2, 297, 314, 
	5, 114, 58, 2, 298, 314, 5, 122, 62, 2, 299, 314, 5, 128, 65, 2, 300, 314, 
	5, 130, 66, 2, 301, 314, 5, 132, 67, 2, 302, 314, 5, 134, 68, 2, 303, 314, 
	5, 136, 69, 2, 304, 314, 5, 138, 70, 2, 305, 314, 5, 140, 71, 2, 306, 314, 
	5, 142, 72, 2, 307, 314, 5, 144, 73, 2, 308, 314, 5, 146, 74, 2, 309, 314, 
	5, 148, 75, 2, 310, 314, 5, 150, 76, 2, 311, 314, 5, 152, 77, 2, 312, 314, 
	5, 154, 78, 2, 313, 283, 3, 2, 2, 2, 313, 284, 3, 2, 2, 2, 313, 285, 3, 
	2, 2, 2, 313, 286, 3, 2, 2, 2, 313, 287, 3, 2, 2, 2, 313, 288, 3, 2, 2, 
	2, 313, 289, 3, 2, 2, 2, 313, 290, 3, 2, 2, 2, 313, 291, 3, 2, 2, 2, 313, 
	292, 3, 2, 2, 2, 313, 293, 3, 2, 2, 2, 313, 294, 3, 2, 2, 2, 313, 295, 
	3, 2, 2, 2, 313, 296, 3, 2, 2, 2, 313, 297, 3, 2, 2, 2, 313, 298, 3, 2, 
	2, 2, 313, 299, 3, 2, 2, 2, 313, 300, 3, 2, 2, 2, 313, 301, 3, 2, 2, 2, 
	313, 302, 3, 2, 2, 2, 313, 303, 3, 2, 2, 2, 313, 304, 3, 2, 2, 2, 313, 
	305, 3, 2, 2, 2, 313, 306, 3, 2, 2, 2, 313, 307, 3, 2, 2, 2, 313, 308, 
	3, 2, 2, 2, 313, 309, 3, 2, 2, 2, 313, 310, 3, 2, 2, 2, 313, 311, 3, 2, 
	2, 2, 313, 312, 3, 2, 2, 2, 314, 21, 3, 2, 2, 2, 315, 316, 7, 33, 2, 2, 
	316, 318, 5, 160, 81, 2, 317, 319, 7, 45, 2, 2, 318, 317, 3, 2, 2, 2, 319, 
	320, 3, 2, 2, 2, 320, 318, 3, 2, 2, 2, 320, 321, 3, 2, 2, 2, 321, 23, 3, 
	2, 2, 2, 322, 329, 5, 28, 15, 2, 323, 325, 7, 45, 2, 2, 324, 323, 3, 2, 
	2, 2, 325, 326, 3, 2, 2, 2, 326, 324, 3, 2, 2, 2, 326, 327, 3, 2, 2, 2, 
	327, 328, 3, 2, 2, 2, 328, 330, 5, 26, 14, 2, 329, 324, 3, 2, 2, 2, 329, 
	330, 3, 2, 2, 2, 330, 332, 3, 2, 2, 2, 331, 333, 7, 45, 2, 2, 332, 331, 
	3, 2, 2, 2, 333, 334, 3, 2, 2, 2, 334, 332, 3, 2, 2, 2, 334, 335, 3, 2, 
	2, 2, 335, 336, 3, 2, 2, 2, 336, 337, 7, 52, 2, 2, 337, 25, 3, 2, 2, 2, 
	338, 339, 5, 16, 9, 2, 339, 27, 3, 2, 2, 2, 340, 341, 7, 39, 2, 2, 341, 
	342, 5, 30, 16, 2, 342, 29, 3, 2, 2, 2, 343, 344, 5, 32, 17, 2, 344, 31, 
	3, 2, 2, 2, 345, 347, 7, 40, 2, 2, 346, 348, 5, 160, 81, 2, 347, 346, 3, 
	2, 2, 2, 347, 348, 3, 2, 2, 2, 348, 33, 3, 2, 2, 2, 349, 358, 5, 36, 19, 
	2, 350, 352, 7, 45, 2, 2, 351, 350, 3, 2, 2, 2, 352, 353, 3, 2, 2, 2, 353, 
	351, 3, 2, 2, 2, 353, 354, 3, 2, 2, 2, 354, 355, 3, 2, 2, 2, 355, 357, 
	5, 40, 21, 2, 356, 351, 3, 2, 2, 2, 357, 360, 3, 2, 2, 2, 358, 356, 3, 
	2, 2, 2, 358, 359, 3, 2, 2, 2, 359, 367, 3, 2, 2, 2, 360, 358, 3, 2, 2, 
	2, 361, 363, 7, 45, 2, 2, 362, 361, 3, 2, 2, 2, 363, 364, 3, 2, 2, 2, 364, 
	362, 3, 2, 2, 2, 364, 365, 3, 2, 2, 2, 365, 366, 3, 2, 2, 2, 366, 368, 
	5, 44, 23, 2, 367, 362, 3, 2, 2, 2, 367, 368, 3, 2, 2, 2, 368, 370, 3, 
	2, 2, 2, 369, 371, 7, 45, 2, 2, 370, 369, 3, 2, 2, 2, 371, 372, 3, 2, 2, 
	2, 372, 370, 3, 2, 2, 2, 372, 373, 3, 2, 2, 2, 373, 374, 3, 2, 2, 2, 374, 
	375, 7, 52, 2, 2, 375, 35, 3, 2, 2, 2, 376, 377, 7, 41, 2, 2, 377, 384, 
	5, 48, 25, 2, 378, 380, 7, 45, 2, 2, 379, 378, 3, 2, 2, 2, 380, 381, 3, 
	2, 2, 2, 381, 379, 3, 2, 2, 2, 381, 382, 3, 2, 2, 2, 382, 383, 3, 2, 2, 
	2, 383, 385, 5, 38, 20, 2, 384, 379, 3, 2, 2, 2, 384, 385, 3, 2, 2, 2, 
	385, 37, 3, 2, 2, 2, 386, 387, 5, 16, 9, 2, 387, 39, 3, 2, 2, 2, 388, 389, 
	7, 49, 2, 2, 389, 396, 5, 50, 26, 2, 390, 392, 7, 45, 2, 2, 391, 390, 3, 
	2, 2, 2, 392, 393, 3, 2, 2, 2, 393, 391, 3, 2, 2, 2, 393, 394, 3, 2, 2, 
	2, 394, 395, 3, 2, 2, 2, 395, 397, 5, 42, 22, 2, 396, 391, 3, 2, 2, 2, 
	396, 397, 3, 2, 2, 2, 397, 41, 3, 2, 2, 2, 398, 399, 5, 16, 9, 2, 399, 
	43, 3, 2, 2, 2, 400, 407, 7, 48, 2, 2, 401, 403, 7, 45, 2, 2, 402, 401, 
	3, 2, 2, 2, 403, 404, 3, 2, 2, 2, 404, 402, 3, 2, 2, 2, 404, 405, 3, 2, 
	2, 2, 405, 406, 3, 2, 2, 2, 406, 408, 5, 46, 24, 2, 407, 402, 3, 2, 2, 
	2, 407, 408, 3, 2, 2, 2, 408, 45, 3, 2, 2, 2, 409, 410, 5, 16, 9, 2, 410, 
	47, 3, 2, 2, 2, 411, 412, 5, 156, 79, 2, 412, 49, 3, 2, 2, 2, 413, 414, 
	5, 156, 79, 2, 414, 51, 3, 2, 2, 2, 415, 422, 5, 54, 28, 2, 416, 418, 7, 
	45, 2, 2, 417, 416, 3, 2, 2, 2, 418, 419, 3, 2, 2, 2, 419, 417, 3, 2, 2, 
	2, 419, 420, 3, 2, 2, 2, 420, 421, 3, 2, 2, 2, 421, 423, 5, 58, 30, 2, 
	422, 417, 3, 2, 2, 2, 422, 423, 3, 2, 2, 2, 423, 430, 3, 2, 2, 2, 424, 
	426, 7, 45, 2, 2, 425, 424, 3, 2, 2, 2, 426, 427, 3, 2, 2, 2, 427, 425, 
	3, 2, 2, 2, 427, 428, 3, 2, 2, 2, 428, 429, 3, 2, 2, 2, 429, 431, 5, 62, 
	32, 2, 430, 425, 3, 2, 2, 2, 430, 431, 3, 2, 2, 2, 431, 433, 3, 2, 2, 2, 
	432, 434, 7, 45, 2, 2, 433, 432, 3, 2, 2, 2, 434, 435, 3, 2, 2, 2, 435, 
	433, 3, 2, 2, 2, 435, 436, 3, 2, 2, 2, 436, 437, 3, 2, 2, 2, 437, 438, 
	7, 52, 2, 2, 438, 53, 3, 2, 2, 2, 439, 446, 7, 42, 2, 2, 440, 442, 7, 45, 
	2, 2, 441, 440, 3, 2, 2, 2, 442, 443, 3, 2, 2, 2, 443, 441, 3, 2, 2, 2, 
	443, 444, 3, 2, 2, 2, 444, 445, 3, 2, 2, 2, 445, 447, 5, 56, 29, 2, 446, 
	441, 3, 2, 2, 2, 446, 447, 3, 2, 2, 2, 447, 55, 3, 2, 2, 2, 448, 449, 5, 
	16, 9, 2, 449, 57, 3, 2, 2, 2, 450, 457, 7, 50, 2, 2, 451, 453, 7, 45, 
	2, 2, 452, 451, 3, 2, 2, 2, 453, 454, 3, 2, 2, 2, 454, 452, 3, 2, 2, 2, 
	454, 455, 3, 2, 2, 2, 455, 456, 3, 2, 2, 2, 456, 458, 5, 60, 31, 2, 457, 
	452, 3, 2, 2, 2, 457, 458, 3, 2, 2, 2, 458, 59, 3, 2, 2, 2, 459, 460, 5, 
	16, 9, 2, 460, 61, 3, 2, 2, 2, 461, 468, 7, 51, 2, 2, 462, 464, 7, 45, 
	2, 2, 463, 462, 3, 2, 2, 2, 464, 465, 3, 2, 2, 2, 465, 463, 3, 2, 2, 2, 
	465, 466, 3, 2, 2, 2, 466, 467, 3, 2, 2, 2, 467, 469, 5, 64, 33, 2, 468, 
	463, 3, 2, 2, 2, 468, 469, 3, 2, 2, 2, 469, 63, 3, 2, 2, 2, 470, 471, 5, 
	16, 9, 2, 471, 65, 3, 2, 2, 2, 472, 474, 5, 68, 35, 2, 473, 475, 7, 45, 
	2, 2, 474, 473, 3, 2, 2, 2, 475, 476, 3, 2, 2, 2, 476, 474, 3, 2, 2, 2, 
	476, 477, 3, 2, 2, 2, 477, 478, 3, 2, 2, 2, 478, 479, 7, 52, 2, 2, 479, 
	67, 3, 2, 2, 2, 480, 481, 7, 43, 2, 2, 481, 488, 5, 72, 37, 2, 482, 484, 
	7, 45, 2, 2, 483, 482, 3, 2, 2, 2, 484, 485, 3, 2, 2, 2, 485, 483, 3, 2, 
	2, 2, 485, 486, 3, 2, 2, 2, 486, 487, 3, 2, 2, 2, 487, 489, 5, 70, 36, 
	2, 488, 483, 3, 2, 2, 2, 488, 489, 3, 2, 2, 2, 489, 69, 3, 2, 2, 2, 490, 
	491, 5, 16, 9, 2, 491, 71, 3, 2, 2, 2, 492, 493, 5, 160, 81, 2, 493, 73, 
	3, 2, 2, 2, 494, 496, 5, 76, 39, 2, 495, 497, 7, 45, 2, 2, 496, 495, 3, 
	2, 2, 2, 497, 498, 3, 2, 2, 2, 498, 496, 3, 2, 2, 2, 498, 499, 3, 2, 2, 
	2, 499, 500, 3, 2, 2, 2, 500, 501, 7, 52, 2, 2, 501, 75, 3, 2, 2, 2, 502, 
	504, 7, 44, 2, 2, 503, 505, 5, 80, 41, 2, 504, 503, 3, 2, 2, 2, 504, 505, 
	3, 2, 2, 2, 505, 512, 3, 2, 2, 2, 506, 508, 7, 45, 2, 2, 507, 506, 3, 2, 
	2, 2, 508, 509, 3, 2, 2, 2, 509, 507, 3, 2, 2, 2, 509, 510, 3, 2, 2, 2, 
	510, 511, 3, 2, 2, 2, 511, 513, 5, 78, 40, 2, 512, 507, 3, 2, 2, 2, 512, 
	513, 3, 2, 2, 2, 513, 77, 3, 2, 2, 2, 514, 515, 5, 16, 9, 2, 515, 79, 3, 
	2, 2, 2, 516, 517, 5, 160, 81, 2, 517, 81, 3, 2, 2, 2, 518, 520, 7, 7, 
	2, 2, 519, 521, 5, 160, 81, 2, 520, 519, 3, 2, 2, 2, 520, 521, 3, 2, 2, 
	2, 521, 83, 3, 2, 2, 2, 522, 524, 7, 8, 2, 2, 523, 525, 5, 160, 81, 2, 
	524, 523, 3, 2, 2, 2, 524, 525, 3, 2, 2, 2, 525, 85, 3, 2, 2, 2, 526, 528, 
	7, 9, 2, 2, 527, 529, 5, 160, 81, 2, 528, 527, 3, 2, 2, 2, 528, 529, 3, 
	2, 2, 2, 529, 87, 3, 2, 2, 2, 530, 532, 7, 10, 2, 2, 531, 533, 5, 160, 
	81, 2, 532, 531, 3, 2, 2, 2, 532, 533, 3, 2, 2, 2, 533, 89, 3, 2, 2, 2, 
	534, 537, 5, 94, 48, 2, 535, 537, 5, 92, 47, 2, 536, 534, 3, 2, 2, 2, 536, 
	535, 3, 2, 2, 2, 537, 91, 3, 2, 2, 2, 538, 540, 7, 12, 2, 2, 539, 541, 
	5, 160, 81, 2, 540, 539, 3, 2, 2, 2, 540, 541, 3, 2, 2, 2, 541, 93, 3, 
	2, 2, 2, 542, 544, 7, 11, 2, 2, 543, 545, 5, 160, 81, 2, 544, 543, 3, 2, 
	2, 2, 544, 545, 3, 2, 2, 2, 545, 95, 3, 2, 2, 2, 546, 548, 7, 13, 2, 2, 
	547, 549, 5, 158, 80, 2, 548, 547, 3, 2, 2, 2, 548, 549, 3, 2, 2, 2, 549, 
	97, 3, 2, 2, 2, 550, 552, 7, 19, 2, 2, 551, 553, 5, 160, 81, 2, 552, 551, 
	3, 2, 2, 2, 552, 553, 3, 2, 2, 2, 553, 99, 3, 2, 2, 2, 554, 556, 7, 20, 
	2, 2, 555, 557, 5, 160, 81, 2, 556, 555, 3, 2, 2, 2, 556, 557, 3, 2, 2, 
	2, 557, 101, 3, 2, 2, 2, 558, 560, 7, 21, 2, 2, 559, 561, 5, 160, 81, 2, 
	560, 559, 3, 2, 2, 2, 560, 561, 3, 2, 2, 2, 561, 103, 3, 2, 2, 2, 562, 
	564, 7, 22, 2, 2, 563, 565, 5, 158, 80, 2, 564, 563, 3, 2, 2, 2, 564, 565, 
	3, 2, 2, 2, 565, 105, 3, 2, 2, 2, 566, 568, 7, 23, 2, 2, 567, 569, 5, 158, 
	80, 2, 568, 567, 3, 2, 2, 2, 568, 569, 3, 2, 2, 2, 569, 107, 3, 2, 2, 2, 
	570, 572, 7, 14, 2, 2, 571, 573, 5, 160, 81, 2, 572, 571, 3, 2, 2, 2, 572, 
	573, 3, 2, 2, 2, 573, 109, 3, 2, 2, 2, 574, 576, 7, 15, 2, 2, 575, 577, 
	5, 158, 80, 2, 576, 575, 3, 2, 2, 2, 576, 577, 3, 2, 2, 2, 577, 111, 3, 
	2, 2, 2, 578, 579, 7, 16, 2, 2, 579, 581, 5, 118, 60, 2, 580, 582, 7, 54, 
	2, 2, 581, 580, 3, 2, 2, 2, 581, 582, 3, 2, 2, 2, 582, 587, 3, 2, 2, 2, 
	583, 585, 7, 46, 2, 2, 584, 583, 3, 2, 2, 2, 584, 585, 3, 2, 2, 2, 585, 
	586, 3, 2, 2, 2, 586, 588, 5, 120, 61, 2, 587, 584, 3, 2, 2, 2, 587, 588, 
	3, 2, 2, 2, 588, 113, 3, 2, 2, 2, 589, 590, 7, 17, 2, 2, 590, 591, 5, 116, 
	59, 2, 591, 599, 5, 118, 60, 2, 592, 597, 7, 54, 2, 2, 593, 595, 7, 46, 
	2, 2, 594, 593, 3, 2, 2, 2, 594, 595, 3, 2, 2, 2, 595, 596, 3, 2, 2, 2, 
	596, 598, 5, 120, 61, 2, 597, 594, 3, 2, 2, 2, 597, 598, 3, 2, 2, 2, 598, 
	600, 3, 2, 2, 2, 599, 592, 3, 2, 2, 2, 599, 600, 3, 2, 2, 2, 600, 115, 
	3, 2, 2, 2, 601, 603, 5, 160, 81, 2, 602, 601, 3, 2, 2, 2, 602, 603, 3, 
	2, 2, 2, 603, 117, 3, 2, 2, 2, 604, 605, 7, 53, 2, 2, 605, 119, 3, 2, 2, 
	2, 606, 613, 7, 53, 2, 2, 607, 609, 7, 46, 2, 2, 608, 607, 3, 2, 2, 2, 
	608, 609, 3, 2, 2, 2, 609, 610, 3, 2, 2, 2, 610, 612, 7, 53, 2, 2, 611, 
	608, 3, 2, 2, 2, 612, 615, 3, 2, 2, 2, 613, 611, 3, 2, 2, 2, 613, 614, 
	3, 2, 2, 2, 614, 121, 3, 2, 2, 2, 615, 613, 3, 2, 2, 2, 616, 623, 7, 18, 
	2, 2, 617, 618, 5, 124, 63, 2, 618, 619, 7, 54, 2, 2, 619, 620, 5, 126, 
	64, 2, 620, 622, 3, 2, 2, 2, 621, 617, 3, 2, 2, 2, 622, 625, 3, 2, 2, 2, 
	623, 621, 3, 2, 2, 2, 623, 624, 3, 2, 2, 2, 624, 123, 3, 2, 2, 2, 625, 
	623, 3, 2, 2, 2, 626, 627, 7, 53, 2, 2, 627, 125, 3, 2, 2, 2, 628, 629, 
	7, 53, 2, 2, 629, 127, 3, 2, 2, 2, 630, 632, 7, 24, 2, 2, 631, 633, 5, 
	160, 81, 2, 632, 631, 3, 2, 2, 2, 632, 633, 3, 2, 2, 2, 633, 129, 3, 2, 
	2, 2, 634, 636, 7, 25, 2, 2, 635, 637, 5, 160, 81, 2, 636, 635, 3, 2, 2, 
	2, 636, 637, 3, 2, 2, 2, 637, 131, 3, 2, 2, 2, 638, 640, 7, 26, 2, 2, 639, 
	641, 5, 160, 81, 2, 640, 639, 3, 2, 2, 2, 640, 641, 3, 2, 2, 2, 641, 133, 
	3, 2, 2, 2, 642, 644, 7, 27, 2, 2, 643, 645, 5, 160, 81, 2, 644, 643, 3, 
	2, 2, 2, 644, 645, 3, 2, 2, 2, 645, 135, 3, 2, 2, 2, 646, 648, 7, 28, 2, 
	2, 647, 649, 5, 160, 81, 2, 648, 647, 3, 2, 2, 2, 648, 649, 3, 2, 2, 2, 
	649, 137, 3, 2, 2, 2, 650, 652, 7, 29, 2, 2, 651, 653, 5, 160, 81, 2, 652, 
	651, 3, 2, 2, 2, 652, 653, 3, 2, 2, 2, 653, 139, 3, 2, 2, 2, 654, 656, 
	7, 31, 2, 2, 655, 657, 5, 160, 81, 2, 656, 655, 3, 2, 2, 2, 656, 657, 3, 
	2, 2, 2, 657, 141, 3, 2, 2, 2, 658, 660, 7, 30, 2, 2, 659, 661, 5, 160, 
	81, 2, 660, 659, 3, 2, 2, 2, 660, 661, 3, 2, 2, 2, 661, 143, 3, 2, 2, 2, 
	662, 664, 7, 32, 2, 2, 663, 665, 5, 160, 81, 2, 664, 663, 3, 2, 2, 2, 664, 
	665, 3, 2, 2, 2, 665, 145, 3, 2, 2, 2, 666, 668, 7, 34, 2, 2, 667, 669, 
	5, 160, 81, 2, 668, 667, 3, 2, 2, 2, 668, 669, 3, 2, 2, 2, 669, 147, 3, 
	2, 2, 2, 670, 672, 7, 35, 2, 2, 671, 673, 5, 160, 81, 2, 672, 671, 3, 2, 
	2, 2, 672, 673, 3, 2, 2, 2, 673, 149, 3, 2, 2, 2, 674, 676, 7, 36, 2, 2, 
	675, 677, 5, 160, 81, 2, 676, 675, 3, 2, 2, 2, 676, 677, 3, 2, 2, 2, 677, 
	151, 3, 2, 2, 2, 678, 680, 7, 37, 2, 2, 679, 681, 5, 160, 81, 2, 680, 679, 
	3, 2, 2, 2, 680, 681, 3, 2, 2, 2, 681, 153, 3, 2, 2, 2, 682, 684, 7, 38, 
	2, 2, 683, 685, 5, 160, 81, 2, 684, 683, 3, 2, 2, 2, 684, 685, 3, 2, 2, 
	2, 685, 155, 3, 2, 2, 2, 686, 687, 5, 158, 80, 2, 687, 157, 3, 2, 2, 2, 
	688, 689, 5, 160, 81, 2, 689, 159, 3, 2, 2, 2, 690, 692, 5, 162, 82, 2, 
	691, 690, 3, 2, 2, 2, 692, 693, 3, 2, 2, 2, 693, 691, 3, 2, 2, 2, 693, 
	694, 3, 2, 2, 2, 694, 161, 3, 2, 2, 2, 695, 696, 7, 53, 2, 2, 696, 163, 
	3, 2, 2, 2, 98, 167, 171, 176, 181, 185, 190, 199, 205, 210, 216, 222, 
	226, 231, 234, 242, 248, 255, 259, 267, 272, 281, 313, 320, 326, 329, 334, 
	347, 353, 358, 364, 367, 372, 381, 384, 393, 396, 404, 407, 419, 422, 427, 
	430, 435, 443, 446, 454, 457, 465, 468, 476, 485, 488, 498, 504, 509, 512, 
	520, 524, 528, 532, 536, 540, 544, 548, 552, 556, 560, 564, 568, 572, 576, 
	581, 584, 587, 594, 597, 599, 602, 608, 613, 623, 632, 636, 640, 644, 648, 
	652, 656, 660, 664, 668, 672, 676, 680, 684, 693,
}
var literalNames = []string{
	"", "", "", "", "", "'FROM'", "'FROM DOCKERFILE'", "'LOCALLY'", "'COPY'", 
	"'SAVE ARTIFACT'", "'SAVE IMAGE'", "'RUN'", "'EXPOSE'", "'VOLUME'", "'ENV'", 
	"'ARG'", "'LABEL'", "'BUILD'", "'WORKDIR'", "'USER'", "'CMD'", "'ENTRYPOINT'", 
	"'GIT CLONE'", "'ADD'", "'STOPSIGNAL'", "'ONBUILD'", "'HEALTHCHECK'", "'SHELL'", 
	"'DO'", "'COMMAND'", "'IMPORT'", "'VERSION'", "'CACHE'", "'HOST'", "'PROJECT'", 
	"'PIPELINE'", "'TRIGGER'", "'WITH'", "", "", "", "", "", "", "", "", "'ELSE'", 
	"'ELSE IF'", "'CATCH'", "'FINALLY'", "'END'",
}
var symbolicNames = []string{
	"", "INDENT", "DEDENT", "Target", "UserCommand", "FROM", "FROM_DOCKERFILE", 
	"LOCALLY", "COPY", "SAVE_ARTIFACT", "SAVE_IMAGE", "RUN", "EXPOSE", "VOLUME", 
	"ENV", "ARG", "LABEL", "BUILD", "WORKDIR", "USER", "CMD", "ENTRYPOINT", 
	"GIT_CLONE", "ADD", "STOPSIGNAL", "ONBUILD", "HEALTHCHECK", "SHELL", "DO", 
	"COMMAND", "IMPORT", "VERSION", "CACHE", "HOST", "PROJECT", "PIPELINE", 
	"TRIGGER", "WITH", "DOCKER", "IF", "TRY", "FOR", "WAIT", "NL", "WS", "COMMENT", 
	"ELSE", "ELSE_IF", "CATCH", "FINALLY", "END", "Atom", "EQUALS",
}

var ruleNames = []string{
	"earthFile", "targets", "targetOrUserCommand", "target", "targetHeader", 
	"userCommand", "userCommandHeader", "stmts", "stmt", "commandStmt", "version", 
	"withStmt", "withBlock", "withExpr", "withCommand", "dockerCommand", "ifStmt", 
	"ifClause", "ifBlock", "elseIfClause", "elseIfBlock", "elseClause", "elseBlock", 
	"ifExpr", "elseIfExpr", "tryStmt", "tryClause", "tryBlock", "catchClause", 
	"catchBlock", "finallyClause", "finallyBlock", "forStmt", "forClause", 
	"forBlock", "forExpr", "waitStmt", "waitClause", "waitBlock", "waitExpr", 
	"fromStmt", "fromDockerfileStmt", "locallyStmt", "copyStmt", "saveStmt", 
	"saveImage", "saveArtifact", "runStmt", "buildStmt", "workdirStmt", "userStmt", 
	"cmdStmt", "entrypointStmt", "exposeStmt", "volumeStmt", "envStmt", "argStmt", 
	"optionalFlag", "envArgKey", "envArgValue", "labelStmt", "labelKey", "labelValue", 
	"gitCloneStmt", "addStmt", "stopsignalStmt", "onbuildStmt", "healthcheckStmt", 
	"shellStmt", "userCommandStmt", "doStmt", "importStmt", "cacheStmt", "hostStmt", 
	"projectStmt", "pipelineStmt", "triggerStmt", "expr", "stmtWordsMaybeJSON", 
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


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << EarthParserFROM) | (1 << EarthParserFROM_DOCKERFILE) | (1 << EarthParserLOCALLY) | (1 << EarthParserCOPY) | (1 << EarthParserSAVE_ARTIFACT) | (1 << EarthParserSAVE_IMAGE) | (1 << EarthParserRUN) | (1 << EarthParserEXPOSE) | (1 << EarthParserVOLUME) | (1 << EarthParserENV) | (1 << EarthParserARG) | (1 << EarthParserLABEL) | (1 << EarthParserBUILD) | (1 << EarthParserWORKDIR) | (1 << EarthParserUSER) | (1 << EarthParserCMD) | (1 << EarthParserENTRYPOINT) | (1 << EarthParserGIT_CLONE) | (1 << EarthParserADD) | (1 << EarthParserSTOPSIGNAL) | (1 << EarthParserONBUILD) | (1 << EarthParserHEALTHCHECK) | (1 << EarthParserSHELL) | (1 << EarthParserDO) | (1 << EarthParserCOMMAND) | (1 << EarthParserIMPORT))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (EarthParserCACHE - 32)) | (1 << (EarthParserHOST - 32)) | (1 << (EarthParserPROJECT - 32)) | (1 << (EarthParserPIPELINE - 32)) | (1 << (EarthParserTRIGGER - 32)) | (1 << (EarthParserWITH - 32)) | (1 << (EarthParserIF - 32)) | (1 << (EarthParserTRY - 32)) | (1 << (EarthParserFOR - 32)) | (1 << (EarthParserWAIT - 32)))) != 0) {
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

func (s *TargetContext) INDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserINDENT, 0)
}

func (s *TargetContext) DEDENT() antlr.TerminalNode {
	return s.GetToken(EarthParserDEDENT, 0)
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


		if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << EarthParserFROM) | (1 << EarthParserFROM_DOCKERFILE) | (1 << EarthParserLOCALLY) | (1 << EarthParserCOPY) | (1 << EarthParserSAVE_ARTIFACT) | (1 << EarthParserSAVE_IMAGE) | (1 << EarthParserRUN) | (1 << EarthParserEXPOSE) | (1 << EarthParserVOLUME) | (1 << EarthParserENV) | (1 << EarthParserARG) | (1 << EarthParserLABEL) | (1 << EarthParserBUILD) | (1 << EarthParserWORKDIR) | (1 << EarthParserUSER) | (1 << EarthParserCMD) | (1 << EarthParserENTRYPOINT) | (1 << EarthParserGIT_CLONE) | (1 << EarthParserADD) | (1 << EarthParserSTOPSIGNAL) | (1 << EarthParserONBUILD) | (1 << EarthParserHEALTHCHECK) | (1 << EarthParserSHELL) | (1 << EarthParserDO) | (1 << EarthParserCOMMAND) | (1 << EarthParserIMPORT))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (EarthParserCACHE - 32)) | (1 << (EarthParserHOST - 32)) | (1 << (EarthParserPROJECT - 32)) | (1 << (EarthParserPIPELINE - 32)) | (1 << (EarthParserTRIGGER - 32)) | (1 << (EarthParserWITH - 32)) | (1 << (EarthParserIF - 32)) | (1 << (EarthParserTRY - 32)) | (1 << (EarthParserFOR - 32)) | (1 << (EarthParserWAIT - 32)))) != 0) {
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

func (s *StmtContext) WaitStmt() IWaitStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitStmtContext)
}

func (s *StmtContext) TryStmt() ITryStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITryStmtContext)(nil)).Elem(), 0)

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

func (s *CommandStmtContext) ProjectStmt() IProjectStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IProjectStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IProjectStmtContext)
}

func (s *CommandStmtContext) PipelineStmt() IPipelineStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPipelineStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPipelineStmtContext)
}

func (s *CommandStmtContext) TriggerStmt() ITriggerStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITriggerStmtContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITryClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITryClauseContext)
}

func (s *TryStmtContext) END() antlr.TerminalNode {
	return s.GetToken(EarthParserEND, 0)
}

func (s *TryStmtContext) CatchClause() ICatchClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICatchClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICatchClauseContext)
}

func (s *TryStmtContext) FinallyClause() IFinallyClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFinallyClauseContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITryBlockContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICatchBlockContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFinallyBlockContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitClauseContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IWaitExprContext)
}

func (s *WaitClauseContext) WaitBlock() IWaitBlockContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWaitBlockContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStmtWordsContext)(nil)).Elem(), 0)

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


