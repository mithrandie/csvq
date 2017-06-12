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
const FLAG = 57354
const SELECT = 57355
const FROM = 57356
const UPDATE = 57357
const SET = 57358
const DELETE = 57359
const WHERE = 57360
const INSERT = 57361
const INTO = 57362
const VALUES = 57363
const AS = 57364
const DUAL = 57365
const STDIN = 57366
const CREATE = 57367
const ADD = 57368
const DROP = 57369
const ALTER = 57370
const TABLE = 57371
const FIRST = 57372
const LAST = 57373
const AFTER = 57374
const BEFORE = 57375
const DEFAULT = 57376
const RENAME = 57377
const TO = 57378
const ORDER = 57379
const GROUP = 57380
const HAVING = 57381
const BY = 57382
const ASC = 57383
const DESC = 57384
const LIMIT = 57385
const JOIN = 57386
const INNER = 57387
const OUTER = 57388
const LEFT = 57389
const RIGHT = 57390
const FULL = 57391
const CROSS = 57392
const ON = 57393
const USING = 57394
const NATURAL = 57395
const UNION = 57396
const ALL = 57397
const ANY = 57398
const EXISTS = 57399
const IN = 57400
const AND = 57401
const OR = 57402
const NOT = 57403
const BETWEEN = 57404
const LIKE = 57405
const IS = 57406
const NULL = 57407
const DISTINCT = 57408
const WITH = 57409
const CASE = 57410
const WHEN = 57411
const THEN = 57412
const ELSE = 57413
const END = 57414
const GROUP_CONCAT = 57415
const SEPARATOR = 57416
const COMMIT = 57417
const ROLLBACK = 57418
const PRINT = 57419
const VAR = 57420
const COMPARISON_OP = 57421
const STRING_OP = 57422
const SUBSTITUTION_OP = 57423

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
	"FLAG",
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
	"COMMIT",
	"ROLLBACK",
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
	"'.'",
	"'('",
	"')'",
	"','",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:1068

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
	-1, 62,
	58, 176,
	62, 176,
	63, 176,
	-2, 156,
	-1, 100,
	44, 178,
	46, 182,
	-2, 116,
	-1, 137,
	90, 77,
	-2, 174,
	-1, 139,
	69, 107,
	-2, 176,
	-1, 144,
	37, 77,
	74, 77,
	90, 77,
	-2, 174,
	-1, 149,
	58, 176,
	62, 176,
	63, 176,
	-2, 101,
	-1, 153,
	58, 176,
	62, 176,
	63, 176,
	-2, 22,
	-1, 155,
	46, 182,
	-2, 178,
	-1, 181,
	58, 176,
	62, 176,
	63, 176,
	-2, 170,
	-1, 247,
	58, 176,
	62, 176,
	63, 176,
	-2, 112,
	-1, 257,
	58, 176,
	62, 176,
	63, 176,
	-2, 26,
	-1, 270,
	58, 176,
	62, 176,
	63, 176,
	-2, 134,
	-1, 287,
	72, 109,
	-2, 176,
	-1, 310,
	58, 176,
	62, 176,
	63, 176,
	-2, 143,
	-1, 316,
	69, 124,
	71, 124,
	72, 124,
	-2, 176,
	-1, 320,
	58, 176,
	62, 176,
	63, 176,
	-2, 52,
	-1, 323,
	58, 176,
	62, 176,
	63, 176,
	-2, 99,
}

const yyPrivate = 57344

const yyLast = 732

var yyAct = [...]int{

	75, 318, 228, 77, 278, 294, 219, 273, 63, 214,
	252, 168, 135, 3, 146, 3, 222, 246, 81, 100,
	52, 52, 79, 247, 97, 279, 47, 156, 154, 200,
	64, 99, 249, 123, 33, 331, 159, 58, 160, 161,
	162, 157, 53, 309, 155, 272, 101, 50, 290, 52,
	108, 62, 122, 111, 267, 52, 264, 115, 116, 225,
	206, 104, 106, 117, 289, 337, 51, 51, 55, 110,
	330, 124, 125, 126, 127, 128, 314, 311, 122, 308,
	301, 67, 158, 271, 266, 119, 113, 29, 244, 198,
	136, 137, 142, 221, 324, 165, 53, 53, 134, 126,
	127, 128, 241, 163, 107, 176, 144, 139, 78, 141,
	171, 52, 136, 173, 87, 52, 87, 89, 170, 149,
	61, 167, 153, 226, 145, 118, 286, 107, 203, 203,
	105, 250, 96, 90, 122, 188, 105, 199, 174, 187,
	189, 133, 181, 215, 182, 183, 175, 298, 190, 191,
	192, 193, 194, 195, 196, 180, 186, 51, 172, 52,
	262, 260, 216, 166, 211, 218, 140, 171, 293, 291,
	210, 202, 209, 152, 204, 223, 205, 229, 232, 171,
	171, 234, 231, 213, 212, 253, 282, 233, 235, 92,
	217, 280, 322, 57, 88, 224, 105, 143, 177, 178,
	227, 207, 56, 230, 53, 239, 53, 179, 255, 238,
	49, 240, 52, 243, 98, 109, 29, 52, 201, 256,
	48, 254, 164, 261, 112, 114, 171, 251, 258, 94,
	149, 259, 232, 257, 170, 171, 263, 269, 60, 265,
	105, 317, 159, 223, 160, 161, 162, 53, 95, 268,
	270, 281, 29, 86, 87, 89, 121, 90, 91, 284,
	302, 52, 53, 52, 332, 11, 236, 237, 171, 59,
	242, 300, 54, 229, 287, 305, 223, 171, 171, 303,
	297, 80, 299, 312, 304, 306, 307, 159, 76, 160,
	161, 162, 157, 105, 16, 155, 15, 321, 105, 52,
	14, 13, 10, 315, 310, 326, 9, 313, 1, 327,
	232, 31, 316, 92, 329, 320, 328, 8, 325, 323,
	274, 275, 276, 277, 7, 229, 142, 335, 73, 12,
	6, 12, 88, 336, 53, 86, 87, 89, 169, 90,
	91, 30, 105, 5, 105, 130, 129, 133, 220, 4,
	122, 53, 86, 87, 89, 320, 90, 91, 30, 248,
	29, 72, 24, 138, 24, 131, 120, 69, 132, 124,
	125, 126, 127, 128, 147, 148, 103, 285, 102, 82,
	105, 68, 71, 65, 70, 185, 184, 84, 66, 319,
	292, 85, 208, 151, 17, 92, 2, 0, 83, 0,
	0, 0, 0, 93, 84, 0, 0, 0, 85, 0,
	0, 0, 92, 0, 88, 83, 0, 0, 0, 74,
	93, 53, 86, 87, 89, 0, 90, 91, 30, 0,
	0, 88, 0, 0, 0, 0, 74, 0, 0, 53,
	86, 87, 89, 0, 90, 91, 30, 133, 0, 0,
	122, 0, 53, 86, 87, 89, 0, 90, 91, 30,
	0, 0, 0, 0, 0, 131, 120, 0, 132, 124,
	125, 126, 127, 128, 84, 0, 0, 0, 85, 0,
	0, 0, 92, 0, 0, 83, 0, 0, 0, 0,
	93, 0, 84, 0, 0, 0, 85, 0, 0, 0,
	92, 88, 150, 83, 0, 84, 74, 0, 93, 85,
	0, 0, 0, 92, 0, 0, 83, 0, 0, 88,
	245, 93, 0, 0, 74, 130, 129, 133, 0, 0,
	122, 0, 88, 0, 0, 0, 0, 74, 0, 333,
	334, 0, 0, 0, 0, 131, 120, 0, 132, 124,
	125, 126, 127, 128, 0, 0, 197, 130, 129, 133,
	0, 0, 122, 0, 0, 0, 0, 0, 130, 129,
	133, 0, 0, 122, 0, 0, 0, 131, 120, 288,
	132, 124, 125, 126, 127, 128, 0, 0, 131, 120,
	0, 132, 124, 125, 126, 127, 128, 130, 129, 133,
	0, 0, 122, 0, 0, 0, 283, 129, 133, 0,
	0, 122, 0, 0, 0, 0, 0, 131, 120, 0,
	132, 124, 125, 126, 127, 128, 131, 120, 0, 132,
	124, 125, 126, 127, 128, 130, 0, 133, 0, 159,
	122, 160, 161, 162, 157, 295, 296, 155, 0, 0,
	0, 0, 0, 0, 0, 131, 120, 0, 132, 124,
	125, 126, 127, 128, 30, 0, 29, 0, 19, 28,
	20, 0, 18, 0, 0, 0, 0, 32, 21, 0,
	0, 22, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 25, 26,
	27, 23,
}
var yyPact = [...]int{

	653, -1000, 653, -58, -58, -58, -58, -58, -58, -58,
	-58, -58, -58, -58, -58, -58, -58, 206, 190, 243,
	258, 173, 164, 227, 39, -1000, -1000, 448, 217, 66,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 196, 38, 243,
	199, -22, 202, -1000, 38, 211, 243, 243, -1000, -28,
	44, 448, 538, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 39, -1000, 347, 2, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 448, 15, 448, -1000, -1000, 110, -1000,
	-1000, -1000, -1000, 17, 42, 417, -1000, 135, 448, -1000,
	-9, -1000, 200, -1000, -1000, -1000, -1000, 203, 74, 243,
	243, -1000, 243, 196, 38, 16, 172, 227, 448, 538,
	448, 330, 80, 77, 448, 448, 448, 448, 448, 448,
	448, -1000, -1000, -1000, 466, -1, 243, 66, 59, 538,
	-1000, 386, -1000, -1000, 66, 248, -1000, -31, 179, 538,
	-1000, 133, 130, 538, 120, 197, 97, 118, 38, -1000,
	-1000, -1000, -1000, -1000, 243, 4, 243, -1000, 206, -32,
	41, 24, -1000, -1000, -1000, 196, 243, 93, 92, 243,
	-1000, 538, -12, 538, 15, 15, 124, 448, 13, 448,
	14, 14, 70, 70, 70, 576, 386, -1000, -1000, -1000,
	-2, 435, 60, 448, 148, -1000, 417, 243, 148, 448,
	448, 38, 117, 97, 116, -1000, 38, -1000, -1000, -1000,
	-35, 448, -6, -37, 196, 243, 448, -1000, -7, -46,
	290, 243, 157, -1000, 243, 150, -1000, -1000, -1000, -1000,
	547, 347, -1000, 538, -1000, -1000, -1000, 286, 54, 59,
	448, 509, -26, 129, -1000, -1000, 125, 538, -1000, 594,
	38, 103, 38, 242, 4, -10, 239, 243, -1000, -1000,
	538, -1000, 243, -1000, -1000, -1000, 243, 243, -11, -48,
	448, -13, 243, 448, -14, 448, -1000, 538, 448, -1000,
	236, 448, -1000, 108, -1000, 448, 5, 242, 38, 594,
	-1000, -1000, 4, -1000, -1000, -1000, -1000, -1000, 290, 243,
	538, -1000, -1000, 386, -1000, -1000, 538, -20, -1000, -56,
	498, -1000, 108, 538, 243, 242, -1000, -1000, -1000, -1000,
	-1000, 448, -1000, -1000, -1000, -25, -1000, -1000,
}
var yyPgo = [...]int{

	0, 308, 396, 12, 394, 26, 24, 393, 392, 10,
	390, 30, 8, 23, 389, 81, 388, 384, 383, 382,
	381, 29, 379, 46, 378, 19, 376, 5, 375, 374,
	367, 363, 359, 16, 17, 1, 31, 47, 2, 14,
	32, 349, 348, 6, 343, 338, 11, 330, 324, 317,
	25, 4, 7, 306, 302, 301, 300, 296, 294, 0,
	288, 3, 108, 22, 281, 18, 361, 328, 269, 265,
	37, 218, 33, 264, 28, 9, 27, 256, 677,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 3, 4, 5,
	5, 6, 6, 7, 7, 8, 8, 9, 9, 10,
	10, 11, 11, 11, 11, 11, 11, 12, 12, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	13, 14, 73, 73, 73, 15, 16, 17, 17, 17,
	17, 17, 17, 17, 17, 17, 17, 18, 18, 18,
	18, 18, 19, 19, 19, 20, 20, 21, 21, 21,
	22, 22, 23, 23, 23, 24, 24, 25, 25, 25,
	25, 25, 25, 26, 26, 26, 26, 26, 27, 27,
	27, 28, 28, 29, 29, 30, 31, 31, 32, 32,
	33, 33, 34, 34, 35, 35, 36, 36, 37, 37,
	38, 38, 39, 39, 40, 40, 41, 41, 41, 41,
	42, 43, 43, 44, 45, 46, 46, 47, 47, 48,
	49, 49, 50, 50, 51, 51, 52, 52, 52, 52,
	52, 53, 53, 54, 55, 56, 57, 58, 59, 60,
	61, 61, 62, 62, 63, 64, 65, 66, 67, 68,
	68, 69, 70, 70, 71, 71, 72, 72, 74, 74,
	75, 75, 76, 76, 76, 76, 77, 77, 78, 78,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 7, 3, 0,
	2, 0, 2, 0, 3, 0, 2, 0, 3, 0,
	2, 1, 1, 1, 1, 1, 1, 1, 3, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 0, 1, 1, 3, 3, 3, 4, 4,
	6, 6, 4, 4, 4, 4, 2, 3, 3, 3,
	3, 3, 3, 3, 2, 4, 1, 0, 2, 2,
	5, 7, 1, 2, 3, 1, 1, 1, 1, 2,
	3, 1, 1, 5, 5, 6, 6, 4, 0, 2,
	4, 1, 1, 1, 3, 5, 0, 1, 0, 2,
	1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
	1, 3, 1, 3, 4, 2, 5, 8, 4, 7,
	3, 1, 3, 6, 3, 1, 3, 4, 5, 6,
	6, 8, 1, 3, 1, 3, 0, 1, 1, 2,
	2, 5, 7, 7, 1, 1, 2, 4, 1, 1,
	1, 2, 1, 2, 1, 1, 1, 1, 3, 1,
	3, 2, 1, 3, 0, 1, 0, 1, 0, 1,
	0, 1, 0, 1, 1, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -41, -44, -47, -48, -49, -53,
	-54, -69, -67, -55, -56, -57, -58, -4, 19, 15,
	17, 25, 28, 78, -66, 75, 76, 77, 16, 13,
	11, -1, -78, 92, -78, -78, -78, -78, -78, -78,
	-78, -78, -78, -78, -78, -78, -78, -5, 14, 20,
	-37, -23, -59, 4, 14, -37, 29, 29, -70, -68,
	11, 81, -13, -12, -11, -18, -16, -15, -20, -30,
	-17, -19, -66, -67, 89, -59, -60, -61, -62, -63,
	-64, -65, -22, 68, 57, 61, 5, 6, 84, 7,
	9, 10, 65, 73, 12, -71, 66, -6, 18, -36,
	-25, -23, -24, -26, 23, -15, 24, 89, -59, 16,
	91, -59, 22, -36, 14, -59, -59, 91, 81, -13,
	80, -77, 64, -72, 83, 84, 85, 86, 87, 60,
	59, 79, 82, 61, -13, -3, 88, 89, -31, -13,
	-15, -13, -61, -62, 89, 82, -39, -29, -28, -13,
	85, -7, 38, -13, -74, 53, -76, 50, 91, 45,
	47, 48, 49, -59, 22, 21, 89, -3, -46, -45,
	-12, -59, -37, -59, -6, -36, 89, 26, 27, 35,
	-70, -13, -13, -13, 56, 55, -72, 62, 58, 63,
	-13, -13, -13, -13, -13, -13, -13, 90, 90, -59,
	-21, -71, -40, 69, -21, -11, 91, 22, -8, 39,
	40, 44, -74, -76, -75, 46, 44, -36, -59, -43,
	-42, 89, -33, -12, -5, 91, 82, -6, -38, -59,
	-50, 89, -59, -12, 89, -12, -15, -15, -63, -65,
	-13, 89, -15, -13, 90, 85, -34, -13, -32, -40,
	71, -13, -9, 37, -39, -59, -9, -13, -34, -25,
	44, -75, 44, -25, 91, -34, 90, 91, -6, -46,
	-13, 90, 91, -52, 30, 31, 32, 33, -51, -50,
	34, -33, 36, 59, -34, 91, 72, -13, 70, 90,
	74, 40, -10, 43, -27, 51, 52, -25, 44, -25,
	-43, 90, 21, -3, -33, -38, -12, -12, 90, 91,
	-13, 90, -59, -13, 90, -34, -13, 5, -35, -14,
	-13, -61, 84, -13, 89, -25, -27, -43, -52, -51,
	90, 91, -73, 41, 42, -38, -35, 90,
}
var yyDef = [...]int{

	1, -2, 1, 188, 188, 188, 188, 188, 188, 188,
	188, 188, 188, 188, 188, 188, 188, 19, 0, 0,
	0, 0, 0, 0, 0, 154, 155, 0, 0, 174,
	167, 2, 3, 189, 4, 5, 6, 7, 8, 9,
	10, 11, 12, 13, 14, 15, 16, 21, 0, 0,
	0, 118, 82, 158, 0, 0, 0, 0, 171, 172,
	169, 0, -2, 39, 40, 41, 42, 43, 44, 45,
	46, 47, 48, 49, 0, 37, 31, 32, 33, 34,
	35, 36, 76, 106, 0, 0, 159, 160, 0, 162,
	164, 165, 166, 0, 0, 0, 175, 23, 0, 20,
	-2, 87, 88, 91, 92, 85, 86, 0, 0, 0,
	0, 83, 0, 21, 0, 0, 0, 0, 0, 168,
	0, 0, 176, 0, 0, 0, 0, 0, 0, 0,
	0, 186, 187, 177, 176, 0, 0, -2, 0, -2,
	66, 74, 161, 163, -2, 0, 18, 122, 103, -2,
	102, 25, 0, -2, 0, -2, 180, 0, 0, 179,
	183, 184, 185, 89, 0, 0, 0, 128, 19, 135,
	0, 37, 119, 84, 137, 21, 0, 0, 0, 0,
	173, -2, 56, 57, 0, 0, 0, 0, 0, 0,
	67, 68, 69, 70, 71, 72, 73, 50, 55, 38,
	0, 0, 108, 0, 27, 157, 0, 0, 27, 0,
	0, 0, 0, 180, 0, 181, 0, 117, 90, 126,
	131, 0, 0, 110, 21, 0, 0, 138, 0, 120,
	146, 0, 142, 151, 0, 0, 64, 65, 58, 59,
	176, 0, 62, 63, 75, 78, 79, -2, 0, 125,
	0, 176, 0, 0, 123, 104, 29, -2, 24, 98,
	0, 0, 0, 97, 0, 0, 0, 0, 133, 136,
	-2, 139, 0, 140, 147, 148, 0, 0, 0, 144,
	0, 0, 0, 0, 0, 0, 105, -2, 0, 80,
	0, 0, 17, 0, 93, 0, 0, 94, 0, 98,
	132, 130, 0, 129, 111, 121, 149, 150, 146, 0,
	-2, 152, 153, 60, 61, 113, -2, 0, 28, 114,
	-2, 30, 0, -2, 0, 96, 95, 127, 141, 145,
	81, 0, 51, 53, 54, 0, 115, 100,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 87, 3, 3,
	89, 90, 85, 83, 91, 84, 88, 86, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 92,
	3, 82,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80, 81,
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
		//line parser.y:127
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:132
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:139
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:143
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:147
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:151
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:155
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:159
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:163
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:167
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:171
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:175
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:179
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:183
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:187
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:191
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 17:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:197
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
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:211
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:217
		{
			yyVAL.expression = nil
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:221
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:227
		{
			yyVAL.expression = nil
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:231
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 23:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:237
		{
			yyVAL.expression = nil
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:241
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 25:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:247
		{
			yyVAL.expression = nil
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:251
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 27:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:257
		{
			yyVAL.expression = nil
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:261
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:267
		{
			yyVAL.expression = nil
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:271
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:277
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:281
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:285
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:289
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:293
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:297
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:303
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:307
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:313
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:317
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:321
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:325
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:329
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:333
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:337
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:349
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:353
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:357
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:363
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:369
		{
			yyVAL.token = Token{}
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:373
		{
			yyVAL.token = yyDollar[1].token
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:377
		{
			yyVAL.token = yyDollar[1].token
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:383
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:389
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
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:412
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 58:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 60:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:424
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 61:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:428
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:432
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 63:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:436
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 64:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:440
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 65:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:444
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:448
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:454
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:458
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:462
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:466
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:470
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:476
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:480
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:484
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 75:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:490
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:494
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:500
		{
			yyVAL.expression = Option{}
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:504
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:508
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 80:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:514
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 81:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:518
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:528
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:538
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:542
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:548
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:560
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:564
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:568
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 93:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:574
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 94:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:578
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 95:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:582
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 96:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:586
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:590
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 98:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expression = nil
		}
	case 99:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:600
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 100:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:604
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:610
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:614
		{
			yyVAL.expression = AllColumns{}
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:620
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:624
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 105:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:630
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 106:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:636
		{
			yyVAL.expression = nil
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:640
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 108:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:646
		{
			yyVAL.expression = nil
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:650
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:656
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:660
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:666
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:670
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:676
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 115:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:680
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:686
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:690
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:696
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:700
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:706
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:710
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 122:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:716
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:720
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 124:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:726
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 125:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:730
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 126:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:736
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 127:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:740
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 128:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:744
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 129:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:748
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 130:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:754
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:760
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 132:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:764
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 133:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:770
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 134:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:776
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:782
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:786
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 137:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:792
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 138:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:797
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 139:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:804
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 140:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:810
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 141:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:814
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 142:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:820
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:830
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:834
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 146:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:840
		{
			yyVAL.expression = nil
		}
	case 147:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:844
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:848
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 149:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:852
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 150:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:856
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 151:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:862
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 152:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:866
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 153:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:872
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:878
		{
			yyVAL.expression = Commit{Literal: yyDollar[1].token.Literal}
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:884
		{
			yyVAL.expression = Rollback{Literal: yyDollar[1].token.Literal}
		}
	case 156:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:890
		{
			yyVAL.expression = Print{Print: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 157:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:896
		{
			yyVAL.expression = SetFlag{Set: yyDollar[1].token.Literal, Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:902
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:908
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:914
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 161:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:918
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:925
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 163:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:929
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 164:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:936
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:942
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:948
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:954
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:960
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:966
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 170:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:970
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 171:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:976
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:982
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 173:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:986
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 174:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:992
		{
			yyVAL.token = Token{}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:996
		{
			yyVAL.token = yyDollar[1].token
		}
	case 176:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1002
		{
			yyVAL.token = Token{}
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1006
		{
			yyVAL.token = yyDollar[1].token
		}
	case 178:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1012
		{
			yyVAL.token = Token{}
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1016
		{
			yyVAL.token = yyDollar[1].token
		}
	case 180:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1022
		{
			yyVAL.token = Token{}
		}
	case 181:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1026
		{
			yyVAL.token = yyDollar[1].token
		}
	case 182:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1032
		{
			yyVAL.token = Token{}
		}
	case 183:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1036
		{
			yyVAL.token = yyDollar[1].token
		}
	case 184:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1040
		{
			yyVAL.token = yyDollar[1].token
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1044
		{
			yyVAL.token = yyDollar[1].token
		}
	case 186:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1050
		{
			yyVAL.token = yyDollar[1].token
		}
	case 187:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1054
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 188:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1060
		{
			yyVAL.token = Token{}
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1064
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
