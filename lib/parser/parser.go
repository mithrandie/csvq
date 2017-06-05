//line parser.y:2
package parser

import __yyfmt__ "fmt"

//line parser.y:2
//line parser.y:5
type yySymType struct {
	yys         int
	program     []Statement
	statement   Statement
	expression  Expression
	expressions []Expression
	primary     Primary
	identifier  Identifier
	text        String
	integer     Integer
	float       Float
	ternary     Ternary
	datetime    Datetime
	null        Null
	variable    Variable
	token       Token
}

const IDENTIFIER = 57346
const STRING = 57347
const INTEGER = 57348
const FLOAT = 57349
const BOOLEAN = 57350
const TERNARY = 57351
const DATETIME = 57352
const VARIABLE = 57353
const SELECT = 57354
const FROM = 57355
const UPDATE = 57356
const SET = 57357
const DELETE = 57358
const WHERE = 57359
const INSERT = 57360
const INTO = 57361
const VALUES = 57362
const AS = 57363
const DUAL = 57364
const STDIN = 57365
const CREATE = 57366
const ADD = 57367
const DROP = 57368
const ALTER = 57369
const TABLE = 57370
const FIRST = 57371
const LAST = 57372
const AFTER = 57373
const BEFORE = 57374
const DEFAULT = 57375
const RENAME = 57376
const TO = 57377
const ORDER = 57378
const GROUP = 57379
const HAVING = 57380
const BY = 57381
const ASC = 57382
const DESC = 57383
const LIMIT = 57384
const JOIN = 57385
const INNER = 57386
const OUTER = 57387
const LEFT = 57388
const RIGHT = 57389
const FULL = 57390
const CROSS = 57391
const ON = 57392
const USING = 57393
const NATURAL = 57394
const UNION = 57395
const ALL = 57396
const ANY = 57397
const EXISTS = 57398
const IN = 57399
const AND = 57400
const OR = 57401
const NOT = 57402
const BETWEEN = 57403
const LIKE = 57404
const IS = 57405
const NULL = 57406
const DISTINCT = 57407
const WITH = 57408
const CASE = 57409
const WHEN = 57410
const THEN = 57411
const ELSE = 57412
const END = 57413
const GROUP_CONCAT = 57414
const SEPARATOR = 57415
const PRINT = 57416
const VAR = 57417
const COMPARISON_OP = 57418
const STRING_OP = 57419
const SUBSTITUTION_OP = 57420

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENTIFIER",
	"STRING",
	"INTEGER",
	"FLOAT",
	"BOOLEAN",
	"TERNARY",
	"DATETIME",
	"VARIABLE",
	"SELECT",
	"FROM",
	"UPDATE",
	"SET",
	"DELETE",
	"WHERE",
	"INSERT",
	"INTO",
	"VALUES",
	"AS",
	"DUAL",
	"STDIN",
	"CREATE",
	"ADD",
	"DROP",
	"ALTER",
	"TABLE",
	"FIRST",
	"LAST",
	"AFTER",
	"BEFORE",
	"DEFAULT",
	"RENAME",
	"TO",
	"ORDER",
	"GROUP",
	"HAVING",
	"BY",
	"ASC",
	"DESC",
	"LIMIT",
	"JOIN",
	"INNER",
	"OUTER",
	"LEFT",
	"RIGHT",
	"FULL",
	"CROSS",
	"ON",
	"USING",
	"NATURAL",
	"UNION",
	"ALL",
	"ANY",
	"EXISTS",
	"IN",
	"AND",
	"OR",
	"NOT",
	"BETWEEN",
	"LIKE",
	"IS",
	"NULL",
	"DISTINCT",
	"WITH",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"END",
	"GROUP_CONCAT",
	"SEPARATOR",
	"PRINT",
	"VAR",
	"COMPARISON_OP",
	"STRING_OP",
	"SUBSTITUTION_OP",
	"'='",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'('",
	"')'",
	"','",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:1012

func SetDebugLevel(level int, verbose bool) {
	yyDebug = level
	yyErrorVerbose = verbose
}

func Parse(s string) ([]Statement, error) {
	l := new(Lexer)
	l.Init(s)
	yyParse(l)
	return l.program, l.err
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 53,
	57, 166,
	61, 166,
	62, 166,
	-2, 147,
	-1, 89,
	43, 168,
	45, 172,
	-2, 109,
	-1, 123,
	86, 72,
	-2, 164,
	-1, 127,
	68, 102,
	-2, 166,
	-1, 132,
	36, 72,
	73, 72,
	86, 72,
	-2, 164,
	-1, 136,
	57, 166,
	61, 166,
	62, 166,
	-2, 96,
	-1, 140,
	57, 166,
	61, 166,
	62, 166,
	-2, 19,
	-1, 142,
	45, 172,
	-2, 168,
	-1, 167,
	57, 166,
	61, 166,
	62, 166,
	-2, 160,
	-1, 230,
	57, 166,
	61, 166,
	62, 166,
	-2, 105,
	-1, 240,
	57, 166,
	61, 166,
	62, 166,
	-2, 23,
	-1, 253,
	57, 166,
	61, 166,
	62, 166,
	-2, 127,
	-1, 269,
	71, 104,
	-2, 166,
	-1, 291,
	57, 166,
	61, 166,
	62, 166,
	-2, 136,
	-1, 297,
	68, 117,
	70, 117,
	71, 117,
	-2, 166,
	-1, 301,
	57, 166,
	61, 166,
	62, 166,
	-2, 47,
	-1, 304,
	57, 166,
	61, 166,
	62, 166,
	-2, 94,
}

const yyPrivate = 57344

const yyLast = 634

var yyAct = [...]int{

	230, 299, 206, 67, 260, 276, 203, 255, 125, 3,
	86, 3, 229, 155, 133, 71, 261, 235, 88, 198,
	69, 38, 183, 53, 143, 141, 232, 112, 49, 90,
	41, 89, 27, 119, 118, 122, 312, 58, 111, 290,
	250, 247, 209, 190, 106, 99, 42, 42, 46, 318,
	54, 120, 109, 108, 121, 113, 114, 115, 116, 117,
	272, 23, 267, 311, 102, 295, 124, 43, 43, 152,
	292, 289, 283, 271, 127, 254, 129, 94, 249, 227,
	44, 44, 130, 94, 111, 136, 186, 205, 140, 305,
	43, 97, 224, 96, 100, 162, 43, 132, 104, 105,
	123, 113, 114, 115, 116, 117, 154, 77, 167, 210,
	168, 169, 128, 160, 176, 177, 178, 179, 180, 181,
	182, 68, 161, 146, 111, 147, 148, 149, 144, 42,
	158, 142, 52, 107, 153, 166, 268, 188, 85, 172,
	44, 94, 150, 115, 116, 117, 80, 77, 79, 157,
	43, 111, 159, 187, 43, 189, 122, 199, 93, 95,
	280, 217, 214, 245, 201, 212, 145, 197, 196, 243,
	200, 188, 211, 233, 223, 195, 226, 208, 273, 194,
	213, 174, 303, 94, 275, 173, 175, 193, 222, 234,
	139, 136, 236, 221, 240, 146, 43, 147, 148, 149,
	131, 82, 202, 264, 207, 237, 262, 241, 219, 220,
	239, 253, 225, 207, 215, 216, 218, 244, 248, 251,
	263, 96, 78, 252, 48, 47, 44, 242, 184, 191,
	163, 164, 246, 94, 269, 40, 23, 266, 94, 165,
	44, 98, 238, 151, 284, 87, 43, 256, 257, 258,
	259, 43, 84, 286, 282, 44, 39, 101, 285, 103,
	157, 23, 298, 291, 45, 215, 294, 51, 207, 44,
	1, 297, 110, 25, 301, 279, 313, 281, 304, 302,
	296, 94, 11, 94, 50, 64, 12, 307, 12, 63,
	21, 308, 21, 70, 43, 310, 43, 309, 66, 13,
	10, 207, 9, 8, 7, 6, 156, 130, 316, 287,
	288, 5, 306, 301, 317, 293, 204, 4, 94, 44,
	76, 77, 79, 231, 80, 81, 24, 126, 60, 134,
	146, 43, 147, 148, 149, 144, 277, 278, 142, 135,
	92, 215, 44, 76, 77, 79, 91, 80, 81, 24,
	23, 72, 44, 76, 77, 79, 207, 80, 81, 24,
	59, 24, 23, 62, 16, 56, 17, 61, 15, 171,
	170, 74, 57, 300, 18, 75, 55, 19, 274, 82,
	192, 138, 73, 44, 76, 77, 79, 83, 80, 81,
	24, 14, 2, 0, 74, 0, 78, 0, 75, 0,
	65, 0, 82, 0, 74, 73, 0, 0, 75, 0,
	83, 0, 82, 0, 0, 73, 0, 0, 0, 78,
	83, 0, 0, 65, 22, 20, 0, 0, 0, 78,
	137, 0, 0, 65, 0, 74, 0, 0, 0, 75,
	0, 0, 0, 82, 0, 0, 73, 44, 76, 77,
	79, 83, 80, 81, 24, 0, 0, 119, 118, 122,
	78, 228, 111, 146, 65, 147, 148, 149, 144, 0,
	0, 142, 0, 0, 0, 120, 109, 0, 121, 113,
	114, 115, 116, 117, 0, 185, 0, 0, 0, 0,
	0, 0, 0, 314, 315, 0, 0, 0, 0, 74,
	0, 0, 0, 75, 0, 0, 0, 82, 0, 0,
	73, 119, 118, 122, 0, 83, 111, 0, 0, 0,
	119, 118, 122, 0, 78, 111, 0, 0, 65, 120,
	109, 270, 121, 113, 114, 115, 116, 117, 120, 109,
	0, 121, 113, 114, 115, 116, 117, 119, 118, 122,
	0, 0, 111, 0, 0, 0, 265, 118, 122, 0,
	0, 111, 0, 0, 0, 120, 109, 0, 121, 113,
	114, 115, 116, 117, 120, 109, 0, 121, 113, 114,
	115, 116, 117, 119, 0, 122, 0, 0, 111, 0,
	0, 0, 0, 0, 122, 0, 0, 111, 0, 0,
	0, 120, 109, 0, 121, 113, 114, 115, 116, 117,
	120, 109, 0, 121, 113, 114, 115, 116, 117, 26,
	0, 0, 0, 0, 28, 29, 30, 31, 32, 33,
	34, 35, 36, 37,
}
var yyPact = [...]int{

	350, -1000, 350, -56, -56, -56, -56, -56, -56, -56,
	-56, -56, -56, -56, 243, 216, 265, 251, 197, 196,
	256, 54, 443, 73, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 228, 136,
	265, 226, -42, 236, -1000, 136, 246, 265, 265, -1000,
	-43, 55, 443, 489, 15, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 54, -1000, 338, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 443, 8, 443, -1000, -1000, 141, -1000,
	-1000, -1000, -1000, 12, 348, -1000, 153, 443, -1000, 79,
	-1000, 222, -1000, -1000, -1000, -1000, 249, 49, 265, 265,
	-1000, 265, 228, 136, 10, 205, 256, 443, 489, 443,
	315, 96, 124, 443, 443, 443, 443, 443, 443, 443,
	-1000, -1000, -1000, 73, 399, 0, 69, 489, -1000, 534,
	-1000, -1000, 73, -1000, -44, 208, 489, -1000, 149, 140,
	489, 132, 151, 112, 127, 136, -1000, -1000, -1000, -1000,
	-1000, 265, 2, 265, -1000, 243, -45, 30, -1000, -1000,
	-1000, 228, 265, 77, 76, 265, -1000, 489, 21, 489,
	8, 8, 137, 443, 7, 443, 61, 61, 88, 88,
	88, 525, 534, -7, 379, -1000, -1000, 103, 443, 156,
	348, 265, 156, 443, 443, 136, 126, 112, 120, -1000,
	136, -1000, -1000, -1000, -46, 443, -8, -47, 228, 265,
	443, -1000, -11, 218, 265, 173, -1000, 265, 168, -1000,
	-1000, -1000, -1000, 498, 338, -1000, 489, -1000, -1000, -1000,
	-25, 65, 69, 443, 462, -13, 139, -1000, -1000, 142,
	489, -1000, 286, 136, 117, 136, 419, 2, -14, 224,
	265, -1000, -1000, 489, -1000, -1000, -1000, -1000, 265, 265,
	-15, -48, 443, -16, 265, 443, -21, 443, -1000, 489,
	443, -1000, 257, 443, -1000, 101, -1000, 443, 4, 419,
	136, 286, -1000, -1000, 2, -1000, -1000, -1000, -1000, 218,
	265, 489, -1000, -1000, 534, -1000, -1000, 489, -23, -1000,
	-51, 453, -1000, 101, 489, 265, 419, -1000, -1000, -1000,
	-1000, -1000, 443, -1000, -1000, -1000, -37, -1000, -1000,
}
var yyPgo = [...]int{

	0, 270, 392, 8, 391, 21, 10, 381, 380, 17,
	378, 376, 0, 373, 37, 372, 367, 365, 363, 360,
	22, 351, 29, 346, 31, 340, 5, 339, 329, 328,
	327, 323, 12, 1, 18, 30, 2, 14, 26, 317,
	316, 6, 311, 306, 13, 305, 304, 303, 16, 4,
	7, 302, 300, 299, 50, 298, 3, 121, 20, 293,
	15, 289, 285, 284, 282, 28, 228, 27, 276, 25,
	19, 24, 272, 619,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 3, 4, 5, 5, 6, 6,
	7, 7, 8, 8, 9, 9, 10, 10, 11, 11,
	11, 11, 11, 11, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 13, 68, 68, 68,
	14, 15, 16, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 17, 17, 17, 17, 17, 18, 18, 18,
	19, 19, 20, 20, 20, 21, 21, 22, 22, 22,
	23, 23, 24, 24, 24, 24, 24, 24, 25, 25,
	25, 25, 25, 26, 26, 26, 27, 27, 28, 28,
	29, 30, 30, 31, 31, 32, 32, 33, 33, 34,
	34, 35, 35, 36, 36, 37, 37, 38, 38, 39,
	39, 39, 39, 40, 41, 41, 42, 43, 44, 44,
	45, 45, 46, 47, 47, 48, 48, 49, 49, 50,
	50, 50, 50, 50, 51, 51, 52, 53, 54, 55,
	56, 56, 57, 57, 58, 59, 60, 61, 62, 63,
	63, 64, 65, 65, 66, 66, 67, 67, 69, 69,
	70, 70, 71, 71, 71, 71, 72, 72, 73, 73,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 7, 3, 0, 2, 0, 2,
	0, 3, 0, 2, 0, 3, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 3, 2, 0, 1, 1,
	3, 3, 3, 4, 4, 6, 6, 4, 4, 4,
	4, 2, 3, 3, 3, 3, 3, 3, 3, 2,
	4, 1, 0, 2, 2, 5, 7, 1, 2, 3,
	1, 1, 1, 1, 2, 3, 1, 1, 5, 5,
	6, 6, 4, 0, 2, 4, 1, 1, 1, 3,
	5, 0, 1, 0, 2, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 1, 3, 4, 2, 5,
	8, 4, 7, 3, 1, 3, 6, 3, 1, 3,
	4, 5, 6, 6, 8, 1, 3, 1, 3, 0,
	1, 1, 2, 2, 5, 7, 7, 2, 1, 1,
	1, 2, 1, 2, 1, 1, 1, 1, 3, 1,
	3, 2, 1, 3, 0, 1, 0, 1, 0, 1,
	0, 1, 0, 1, 1, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -39, -42, -45, -46, -47, -51,
	-52, -64, -62, -53, -4, 18, 14, 16, 24, 27,
	75, -61, 74, 12, 11, -1, -73, 88, -73, -73,
	-73, -73, -73, -73, -73, -73, -73, -73, -5, 13,
	19, -35, -22, -54, 4, 13, -35, 28, 28, -65,
	-63, 11, 78, -12, -54, -11, -17, -15, -14, -19,
	-29, -16, -18, -61, -62, 85, -55, -56, -57, -58,
	-59, -60, -21, 67, 56, 60, 5, 6, 81, 7,
	9, 10, 64, 72, -66, 65, -6, 17, -34, -24,
	-22, -23, -25, 22, -14, 23, 85, -54, 15, 87,
	-54, 21, -34, 13, -54, -54, 87, 78, -12, 77,
	-72, 63, -67, 80, 81, 82, 83, 84, 59, 58,
	76, 79, 60, 85, -12, -3, -30, -12, -14, -12,
	-56, -57, 85, -37, -28, -27, -12, 82, -7, 37,
	-12, -69, 52, -71, 49, 87, 44, 46, 47, 48,
	-54, 21, 20, 85, -3, -44, -43, -54, -35, -54,
	-6, -34, 85, 25, 26, 34, -65, -12, -12, -12,
	55, 54, -67, 61, 57, 62, -12, -12, -12, -12,
	-12, -12, -12, -20, -66, 86, 86, -38, 68, -20,
	87, 21, -8, 38, 39, 43, -69, -71, -70, 45,
	43, -34, -54, -41, -40, 85, -36, -54, -5, 87,
	79, -6, -36, -48, 85, -54, -54, 85, -54, -14,
	-14, -58, -60, -12, 85, -14, -12, 86, 82, -32,
	-12, -31, -38, 70, -12, -9, 36, -37, -54, -9,
	-12, -32, -24, 43, -70, 43, -24, 87, -32, 86,
	87, -6, -44, -12, 86, -50, 29, 30, 31, 32,
	-49, -48, 33, -36, 35, 58, -32, 87, 71, -12,
	69, 86, 73, 39, -10, 42, -26, 50, 51, -24,
	43, -24, -41, 86, 20, -3, -36, -54, -54, 86,
	87, -12, 86, -54, -12, 86, -32, -12, 5, -33,
	-13, -12, -56, 81, -12, 85, -24, -26, -41, -50,
	-49, 86, 87, -68, 40, 41, -36, -33, 86,
}
var yyDef = [...]int{

	1, -2, 1, 178, 178, 178, 178, 178, 178, 178,
	178, 178, 178, 178, 16, 0, 0, 0, 0, 0,
	0, 0, 0, 164, 157, 2, 3, 179, 4, 5,
	6, 7, 8, 9, 10, 11, 12, 13, 18, 0,
	0, 0, 111, 77, 148, 0, 0, 0, 0, 161,
	162, 159, 0, -2, 34, 35, 36, 37, 38, 39,
	40, 41, 42, 43, 44, 0, 28, 29, 30, 31,
	32, 33, 71, 101, 0, 0, 149, 150, 0, 152,
	154, 155, 156, 0, 0, 165, 20, 0, 17, -2,
	82, 83, 86, 87, 80, 81, 0, 0, 0, 0,
	78, 0, 18, 0, 0, 0, 0, 0, 158, 0,
	0, 166, 0, 0, 0, 0, 0, 0, 0, 0,
	176, 177, 167, -2, 166, 0, 0, -2, 61, 69,
	151, 153, -2, 15, 115, 98, -2, 97, 22, 0,
	-2, 0, -2, 170, 0, 0, 169, 173, 174, 175,
	84, 0, 0, 0, 121, 16, 128, 0, 112, 79,
	130, 18, 0, 0, 0, 0, 163, -2, 51, 52,
	0, 0, 0, 0, 0, 0, 62, 63, 64, 65,
	66, 67, 68, 0, 0, 45, 50, 103, 0, 24,
	0, 0, 24, 0, 0, 0, 0, 170, 0, 171,
	0, 110, 85, 119, 124, 0, 0, 113, 18, 0,
	0, 131, 0, 139, 0, 135, 144, 0, 0, 59,
	60, 53, 54, 166, 0, 57, 58, 70, 73, 74,
	-2, 0, 118, 0, 166, 0, 0, 116, 99, 26,
	-2, 21, 93, 0, 0, 0, 92, 0, 0, 0,
	0, 126, 129, -2, 132, 133, 140, 141, 0, 0,
	0, 137, 0, 0, 0, 0, 0, 0, 100, -2,
	0, 75, 0, 0, 14, 0, 88, 0, 0, 89,
	0, 93, 125, 123, 0, 122, 114, 142, 143, 139,
	0, -2, 145, 146, 55, 56, 106, -2, 0, 25,
	107, -2, 27, 0, -2, 0, 91, 90, 120, 134,
	138, 76, 0, 46, 48, 49, 0, 108, 95,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 84, 3, 3,
	85, 86, 82, 80, 87, 81, 3, 83, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 88,
	3, 79,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:121
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:126
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:133
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:137
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:141
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:145
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:149
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:153
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:157
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:161
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:165
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:169
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:173
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 14:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:179
		{
			yyVAL.expression = SelectQuery{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
				OrderByClause: yyDollar[6].expression,
				LimitClause:   yyDollar[7].expression,
			}
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:193
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 16:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:199
		{
			yyVAL.expression = nil
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:203
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 18:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:209
		{
			yyVAL.expression = nil
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:213
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:219
		{
			yyVAL.expression = nil
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:223
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 22:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:229
		{
			yyVAL.expression = nil
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:233
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 24:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:239
		{
			yyVAL.expression = nil
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:243
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 26:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:249
		{
			yyVAL.expression = nil
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:253
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:259
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:263
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:267
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:271
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:275
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:279
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:285
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:289
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:293
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:297
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:301
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:305
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:309
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:313
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:317
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:321
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:325
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:329
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:335
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 47:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:341
		{
			yyVAL.token = Token{}
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:345
		{
			yyVAL.token = yyDollar[1].token
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:349
		{
			yyVAL.token = yyDollar[1].token
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:355
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:361
		{
			var item1 []Expression
			var item2 []Expression

			c1, ok := yyDollar[1].expression.(Concat)
			if ok {
				item1 = c1.Items
			} else {
				item1 = []Expression{yyDollar[1].expression}
			}

			c2, ok := yyDollar[3].expression.(Concat)
			if ok {
				item2 = c2.Items
			} else {
				item2 = []Expression{yyDollar[3].expression}
			}

			yyVAL.expression = Concat{Items: append(item1, item2...)}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:384
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:388
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 54:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:392
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 55:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:396
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 56:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:400
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 57:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:404
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 58:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:408
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:412
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 61:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:426
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:430
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:434
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:438
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:442
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:448
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:452
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:456
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 70:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:462
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:466
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 72:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:472
		{
			yyVAL.expression = Option{}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:476
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:480
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 75:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:486
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 76:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:490
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:496
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:500
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:504
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:510
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:514
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:528
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 88:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:546
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 89:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:550
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 90:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:554
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 91:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:558
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 92:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:562
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:568
		{
			yyVAL.expression = nil
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 95:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:582
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:586
		{
			yyVAL.expression = AllColumns{}
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 100:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:602
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 101:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:608
		{
			yyVAL.expression = nil
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 103:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:618
		{
			yyVAL.expression = nil
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:622
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:628
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:632
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:638
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:642
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:648
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:652
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:658
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:662
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:668
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 114:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:672
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:678
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:682
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 117:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:688
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 118:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:692
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 119:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:698
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 120:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:702
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 121:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:706
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 122:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:710
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:716
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:722
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:726
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 126:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:732
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:738
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].identifier, Value: yyDollar[3].expression}
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:744
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 129:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:748
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 130:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:754
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 131:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:759
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 132:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:766
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 133:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:772
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 134:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:776
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:782
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:786
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 138:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:796
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 139:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:802
		{
			yyVAL.expression = nil
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:806
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 141:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:810
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 142:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:814
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].identifier}
		}
	case 143:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:818
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].identifier}
		}
	case 144:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].identifier}}
		}
	case 145:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:828
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 146:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:834
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].identifier, To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 147:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:840
		{
			yyVAL.expression = Print{Print: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:846
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 149:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:852
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:858
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 151:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:862
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:869
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:873
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:880
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:886
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:892
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 157:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:898
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 158:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:904
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:910
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 160:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:914
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 161:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:920
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:926
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 163:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:930
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 164:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:936
		{
			yyVAL.token = Token{}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:940
		{
			yyVAL.token = yyDollar[1].token
		}
	case 166:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:946
		{
			yyVAL.token = Token{}
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:950
		{
			yyVAL.token = yyDollar[1].token
		}
	case 168:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:956
		{
			yyVAL.token = Token{}
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:960
		{
			yyVAL.token = yyDollar[1].token
		}
	case 170:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:966
		{
			yyVAL.token = Token{}
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:970
		{
			yyVAL.token = yyDollar[1].token
		}
	case 172:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:976
		{
			yyVAL.token = Token{}
		}
	case 173:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:980
		{
			yyVAL.token = yyDollar[1].token
		}
	case 174:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:984
		{
			yyVAL.token = yyDollar[1].token
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:988
		{
			yyVAL.token = yyDollar[1].token
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:994
		{
			yyVAL.token = yyDollar[1].token
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:998
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 178:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1004
		{
			yyVAL.token = Token{}
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1008
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
