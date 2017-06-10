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
	"'('",
	"')'",
	"','",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:1046

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
	58, 172,
	62, 172,
	63, 172,
	-2, 152,
	-1, 99,
	44, 174,
	46, 178,
	-2, 112,
	-1, 133,
	89, 75,
	-2, 170,
	-1, 137,
	69, 105,
	-2, 172,
	-1, 142,
	37, 75,
	74, 75,
	89, 75,
	-2, 170,
	-1, 147,
	58, 172,
	62, 172,
	63, 172,
	-2, 99,
	-1, 151,
	58, 172,
	62, 172,
	63, 172,
	-2, 22,
	-1, 153,
	46, 178,
	-2, 174,
	-1, 178,
	58, 172,
	62, 172,
	63, 172,
	-2, 166,
	-1, 242,
	58, 172,
	62, 172,
	63, 172,
	-2, 108,
	-1, 252,
	58, 172,
	62, 172,
	63, 172,
	-2, 26,
	-1, 265,
	58, 172,
	62, 172,
	63, 172,
	-2, 130,
	-1, 281,
	72, 107,
	-2, 172,
	-1, 303,
	58, 172,
	62, 172,
	63, 172,
	-2, 139,
	-1, 309,
	69, 120,
	71, 120,
	72, 120,
	-2, 172,
	-1, 313,
	58, 172,
	62, 172,
	63, 172,
	-2, 50,
	-1, 316,
	58, 172,
	62, 172,
	63, 172,
	-2, 97,
}

const yyPrivate = 57344

const yyLast = 678

var yyAct = [...]int{

	242, 311, 218, 76, 267, 272, 288, 215, 166, 210,
	96, 80, 99, 247, 135, 3, 273, 3, 144, 47,
	154, 241, 98, 78, 152, 64, 244, 194, 62, 58,
	122, 33, 100, 50, 324, 302, 262, 259, 221, 63,
	284, 202, 116, 129, 128, 132, 109, 330, 121, 323,
	307, 67, 51, 51, 55, 283, 304, 301, 295, 52,
	52, 53, 118, 130, 119, 29, 131, 123, 124, 125,
	126, 127, 266, 163, 279, 134, 261, 112, 239, 197,
	103, 105, 53, 137, 53, 139, 217, 317, 52, 107,
	132, 140, 110, 121, 52, 147, 114, 115, 151, 157,
	104, 158, 159, 160, 155, 236, 104, 153, 130, 119,
	121, 131, 123, 124, 125, 126, 127, 77, 178, 106,
	179, 180, 165, 171, 187, 188, 189, 190, 191, 192,
	193, 125, 126, 127, 86, 138, 172, 121, 173, 142,
	164, 161, 51, 169, 156, 106, 177, 133, 168, 52,
	222, 170, 183, 52, 143, 61, 123, 124, 125, 126,
	127, 280, 117, 198, 199, 104, 229, 95, 226, 201,
	200, 86, 88, 199, 209, 245, 224, 185, 208, 213,
	89, 184, 186, 223, 121, 235, 220, 238, 132, 211,
	157, 225, 158, 159, 160, 234, 52, 292, 257, 255,
	246, 212, 214, 147, 219, 141, 252, 233, 104, 285,
	207, 287, 315, 219, 227, 228, 230, 206, 251, 256,
	254, 249, 205, 265, 150, 258, 248, 276, 253, 274,
	264, 263, 275, 231, 232, 57, 91, 237, 56, 260,
	29, 174, 175, 250, 195, 53, 281, 52, 296, 87,
	176, 203, 52, 268, 269, 270, 271, 49, 278, 104,
	97, 168, 108, 162, 104, 298, 227, 294, 291, 219,
	293, 53, 48, 113, 94, 303, 297, 157, 306, 158,
	159, 160, 155, 309, 53, 153, 313, 29, 93, 111,
	316, 314, 60, 310, 54, 52, 53, 52, 120, 325,
	319, 308, 219, 11, 320, 318, 321, 104, 322, 104,
	299, 300, 73, 12, 1, 12, 305, 31, 59, 140,
	328, 72, 24, 79, 24, 313, 329, 75, 53, 85,
	86, 88, 52, 89, 90, 30, 16, 15, 14, 13,
	10, 9, 227, 8, 104, 53, 85, 86, 88, 7,
	89, 90, 30, 6, 29, 167, 5, 219, 53, 85,
	86, 88, 216, 89, 90, 30, 4, 53, 85, 86,
	88, 243, 89, 90, 30, 136, 69, 145, 146, 182,
	181, 83, 102, 101, 81, 84, 68, 71, 65, 91,
	70, 66, 82, 312, 286, 204, 149, 92, 83, 17,
	2, 0, 84, 0, 0, 0, 91, 0, 87, 82,
	0, 83, 74, 0, 92, 84, 0, 0, 0, 91,
	83, 0, 82, 0, 84, 87, 0, 92, 91, 74,
	0, 82, 0, 0, 0, 0, 92, 0, 87, 148,
	0, 0, 74, 0, 129, 128, 132, 87, 240, 121,
	0, 74, 53, 85, 86, 88, 0, 89, 90, 30,
	0, 0, 0, 0, 130, 119, 0, 131, 123, 124,
	125, 126, 127, 0, 196, 0, 0, 0, 0, 0,
	0, 326, 327, 157, 0, 158, 159, 160, 155, 289,
	290, 153, 0, 0, 0, 0, 0, 0, 0, 129,
	128, 132, 0, 0, 121, 83, 0, 0, 0, 84,
	0, 0, 0, 91, 0, 0, 82, 0, 0, 130,
	119, 92, 131, 123, 124, 125, 126, 127, 129, 128,
	132, 0, 87, 121, 0, 0, 74, 0, 0, 282,
	129, 128, 132, 0, 0, 121, 0, 0, 130, 119,
	0, 131, 123, 124, 125, 126, 127, 277, 128, 132,
	130, 119, 121, 131, 123, 124, 125, 126, 127, 129,
	0, 132, 0, 0, 121, 0, 0, 130, 119, 0,
	131, 123, 124, 125, 126, 127, 0, 0, 0, 130,
	119, 0, 131, 123, 124, 125, 126, 127, 85, 86,
	88, 0, 89, 90, 0, 30, 0, 29, 0, 19,
	28, 20, 0, 18, 0, 0, 0, 0, 32, 21,
	0, 0, 22, 34, 35, 36, 37, 38, 39, 40,
	41, 42, 43, 44, 45, 46, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 91, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 25,
	26, 27, 23, 0, 0, 0, 0, 87,
}
var yyPact = [...]int{

	594, -1000, 594, -60, -60, -60, -60, -60, -60, -60,
	-60, -60, -60, -60, -60, -60, -60, 258, 237, 292,
	280, 209, 206, 281, 74, -1000, -1000, 448, 276, 101,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 242, 57, 292,
	246, -44, 267, -1000, 57, 259, 292, 292, -1000, -48,
	81, 448, 481, 59, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 74, -1000, 341, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 448, 31, 448, -1000, -1000, 165, -1000, -1000,
	-1000, -1000, 51, 72, 354, -1000, 186, 448, -1000, 54,
	-1000, 241, -1000, -1000, -1000, -1000, 274, 52, 292, 292,
	-1000, 292, 242, 57, 50, 215, 281, 448, 481, 448,
	324, 127, 119, 448, 448, 448, 448, 448, 448, 448,
	-1000, -1000, -1000, 101, 385, -10, 95, 481, -1000, 29,
	-1000, -1000, 101, 593, -1000, -49, 229, 481, -1000, 183,
	177, 481, 166, 145, 143, 157, 57, -1000, -1000, -1000,
	-1000, -1000, 292, -2, 292, -1000, 258, -52, 68, -1000,
	-1000, -1000, 242, 292, 80, 78, 292, -1000, 481, 73,
	481, 31, 31, 171, 448, 17, 448, 46, 46, 120,
	120, 120, 510, 29, -11, 363, -1000, -1000, 104, 448,
	189, -1000, 354, 292, 189, 448, 448, 57, 155, 143,
	154, -1000, 57, -1000, -1000, -1000, -53, 448, -13, -54,
	242, 292, 448, -1000, -17, 223, 292, 195, -1000, 292,
	191, -1000, -1000, -1000, -1000, 498, 341, -1000, 481, -1000,
	-1000, -1000, -16, 89, 95, 448, 469, -34, 169, -1000,
	-1000, 168, 481, -1000, 438, 57, 153, 57, 232, -2,
	-31, 227, 292, -1000, -1000, 481, -1000, -1000, -1000, -1000,
	292, 292, -32, -55, 448, -33, 292, 448, -39, 448,
	-1000, 481, 448, -1000, 288, 448, -1000, 128, -1000, 448,
	-1, 232, 57, 438, -1000, -1000, -2, -1000, -1000, -1000,
	-1000, 223, 292, 481, -1000, -1000, 29, -1000, -1000, 481,
	-40, -1000, -56, 440, -1000, 128, 481, 292, 232, -1000,
	-1000, -1000, -1000, -1000, 448, -1000, -1000, -1000, -42, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 314, 400, 14, 399, 19, 10, 396, 395, 13,
	394, 25, 0, 393, 51, 391, 390, 388, 387, 386,
	27, 384, 32, 383, 12, 382, 6, 378, 377, 376,
	375, 371, 21, 1, 22, 33, 2, 18, 26, 366,
	362, 7, 356, 355, 8, 353, 349, 343, 16, 5,
	4, 341, 340, 339, 338, 337, 336, 39, 327, 3,
	117, 23, 323, 11, 321, 312, 318, 303, 29, 244,
	30, 299, 24, 9, 20, 298, 618,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 3, 4, 5,
	5, 6, 6, 7, 7, 8, 8, 9, 9, 10,
	10, 11, 11, 11, 11, 11, 11, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 13,
	71, 71, 71, 14, 15, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 17, 17, 17, 17, 17,
	18, 18, 18, 19, 19, 20, 20, 20, 21, 21,
	22, 22, 22, 23, 23, 24, 24, 24, 24, 24,
	24, 25, 25, 25, 25, 25, 26, 26, 26, 27,
	27, 28, 28, 29, 30, 30, 31, 31, 32, 32,
	33, 33, 34, 34, 35, 35, 36, 36, 37, 37,
	38, 38, 39, 39, 39, 39, 40, 41, 41, 42,
	43, 44, 44, 45, 45, 46, 47, 47, 48, 48,
	49, 49, 50, 50, 50, 50, 50, 51, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 59, 60, 60,
	61, 62, 63, 64, 65, 66, 66, 67, 68, 68,
	69, 69, 70, 70, 72, 72, 73, 73, 74, 74,
	74, 74, 75, 75, 76, 76,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 7, 3, 0,
	2, 0, 2, 0, 3, 0, 2, 0, 3, 0,
	2, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 3, 2,
	0, 1, 1, 3, 3, 3, 4, 4, 6, 6,
	4, 4, 4, 4, 2, 3, 3, 3, 3, 3,
	3, 3, 2, 4, 1, 0, 2, 2, 5, 7,
	1, 2, 3, 1, 1, 1, 1, 2, 3, 1,
	1, 5, 5, 6, 6, 4, 0, 2, 4, 1,
	1, 1, 3, 5, 0, 1, 0, 2, 1, 3,
	1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
	4, 2, 5, 8, 4, 7, 3, 1, 3, 6,
	3, 1, 3, 4, 5, 6, 6, 8, 1, 3,
	1, 3, 0, 1, 1, 2, 2, 5, 7, 7,
	1, 1, 2, 4, 1, 1, 1, 2, 1, 2,
	1, 1, 1, 1, 3, 1, 3, 2, 1, 3,
	0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
	1, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -39, -42, -45, -46, -47, -51,
	-52, -67, -65, -53, -54, -55, -56, -4, 19, 15,
	17, 25, 28, 78, -64, 75, 76, 77, 16, 13,
	11, -1, -76, 91, -76, -76, -76, -76, -76, -76,
	-76, -76, -76, -76, -76, -76, -76, -5, 14, 20,
	-35, -22, -57, 4, 14, -35, 29, 29, -68, -66,
	11, 81, -12, -57, -11, -17, -15, -14, -19, -29,
	-16, -18, -64, -65, 88, -58, -59, -60, -61, -62,
	-63, -21, 68, 57, 61, 5, 6, 84, 7, 9,
	10, 65, 73, 12, -69, 66, -6, 18, -34, -24,
	-22, -23, -25, 23, -14, 24, 88, -57, 16, 90,
	-57, 22, -34, 14, -57, -57, 90, 81, -12, 80,
	-75, 64, -70, 83, 84, 85, 86, 87, 60, 59,
	79, 82, 61, 88, -12, -3, -30, -12, -14, -12,
	-59, -60, 88, 82, -37, -28, -27, -12, 85, -7,
	38, -12, -72, 53, -74, 50, 90, 45, 47, 48,
	49, -57, 22, 21, 88, -3, -44, -43, -57, -35,
	-57, -6, -34, 88, 26, 27, 35, -68, -12, -12,
	-12, 56, 55, -70, 62, 58, 63, -12, -12, -12,
	-12, -12, -12, -12, -20, -69, 89, 89, -38, 69,
	-20, -11, 90, 22, -8, 39, 40, 44, -72, -74,
	-73, 46, 44, -34, -57, -41, -40, 88, -36, -57,
	-5, 90, 82, -6, -36, -48, 88, -57, -57, 88,
	-57, -14, -14, -61, -63, -12, 88, -14, -12, 89,
	85, -32, -12, -31, -38, 71, -12, -9, 37, -37,
	-57, -9, -12, -32, -24, 44, -73, 44, -24, 90,
	-32, 89, 90, -6, -44, -12, 89, -50, 30, 31,
	32, 33, -49, -48, 34, -36, 36, 59, -32, 90,
	72, -12, 70, 89, 74, 40, -10, 43, -26, 51,
	52, -24, 44, -24, -41, 89, 21, -3, -36, -57,
	-57, 89, 90, -12, 89, -57, -12, 89, -32, -12,
	5, -33, -13, -12, -59, 84, -12, 88, -24, -26,
	-41, -50, -49, 89, 90, -71, 41, 42, -36, -33,
	89,
}
var yyDef = [...]int{

	1, -2, 1, 184, 184, 184, 184, 184, 184, 184,
	184, 184, 184, 184, 184, 184, 184, 19, 0, 0,
	0, 0, 0, 0, 0, 150, 151, 0, 0, 170,
	163, 2, 3, 185, 4, 5, 6, 7, 8, 9,
	10, 11, 12, 13, 14, 15, 16, 21, 0, 0,
	0, 114, 80, 154, 0, 0, 0, 0, 167, 168,
	165, 0, -2, 37, 38, 39, 40, 41, 42, 43,
	44, 45, 46, 47, 0, 31, 32, 33, 34, 35,
	36, 74, 104, 0, 0, 155, 156, 0, 158, 160,
	161, 162, 0, 0, 0, 171, 23, 0, 20, -2,
	85, 86, 89, 90, 83, 84, 0, 0, 0, 0,
	81, 0, 21, 0, 0, 0, 0, 0, 164, 0,
	0, 172, 0, 0, 0, 0, 0, 0, 0, 0,
	182, 183, 173, -2, 172, 0, 0, -2, 64, 72,
	157, 159, -2, 0, 18, 118, 101, -2, 100, 25,
	0, -2, 0, -2, 176, 0, 0, 175, 179, 180,
	181, 87, 0, 0, 0, 124, 19, 131, 0, 115,
	82, 133, 21, 0, 0, 0, 0, 169, -2, 54,
	55, 0, 0, 0, 0, 0, 0, 65, 66, 67,
	68, 69, 70, 71, 0, 0, 48, 53, 106, 0,
	27, 153, 0, 0, 27, 0, 0, 0, 0, 176,
	0, 177, 0, 113, 88, 122, 127, 0, 0, 116,
	21, 0, 0, 134, 0, 142, 0, 138, 147, 0,
	0, 62, 63, 56, 57, 172, 0, 60, 61, 73,
	76, 77, -2, 0, 121, 0, 172, 0, 0, 119,
	102, 29, -2, 24, 96, 0, 0, 0, 95, 0,
	0, 0, 0, 129, 132, -2, 135, 136, 143, 144,
	0, 0, 0, 140, 0, 0, 0, 0, 0, 0,
	103, -2, 0, 78, 0, 0, 17, 0, 91, 0,
	0, 92, 0, 96, 128, 126, 0, 125, 117, 145,
	146, 142, 0, -2, 148, 149, 58, 59, 109, -2,
	0, 28, 110, -2, 30, 0, -2, 0, 94, 93,
	123, 137, 141, 79, 0, 49, 51, 52, 0, 111,
	98,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 87, 3, 3,
	88, 89, 85, 83, 90, 84, 3, 86, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 91,
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
		//line parser.y:125
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:130
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:137
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:141
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:145
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:149
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:153
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:157
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:161
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:165
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:169
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:173
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:177
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:181
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:185
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:189
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 17:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:195
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
		//line parser.y:209
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:215
		{
			yyVAL.expression = nil
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:219
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:225
		{
			yyVAL.expression = nil
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:229
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 23:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:235
		{
			yyVAL.expression = nil
		}
	case 24:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:239
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 25:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:245
		{
			yyVAL.expression = nil
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:249
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 27:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:255
		{
			yyVAL.expression = nil
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:259
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:265
		{
			yyVAL.expression = nil
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:269
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:275
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:279
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:283
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:287
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:291
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:295
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:301
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:305
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:309
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:313
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:317
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:321
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:325
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:329
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:333
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:337
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 47:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:351
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:357
		{
			yyVAL.token = Token{}
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:361
		{
			yyVAL.token = yyDollar[1].token
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:365
		{
			yyVAL.token = yyDollar[1].token
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:371
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:377
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
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:400
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:404
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 57:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:408
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 58:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:412
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 59:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 61:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:424
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:428
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 63:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:432
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:436
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:442
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:446
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:450
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:454
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:458
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:464
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:468
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:472
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 73:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:478
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:482
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 75:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:488
		{
			yyVAL.expression = Option{}
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:492
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:496
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 78:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:502
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 79:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:506
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:512
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:526
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:530
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:544
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:548
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 91:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:562
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 92:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:566
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 93:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:570
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 94:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:574
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 95:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:578
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 96:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:584
		{
			yyVAL.expression = nil
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:588
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 98:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:598
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:602
		{
			yyVAL.expression = AllColumns{}
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:608
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 103:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:618
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 104:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:624
		{
			yyVAL.expression = nil
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:628
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 106:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:634
		{
			yyVAL.expression = nil
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:638
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:644
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:648
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:654
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:658
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:664
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:668
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:674
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 115:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:678
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:684
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:688
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:694
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:698
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 120:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:704
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 121:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:708
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 122:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:714
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 123:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:718
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 124:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:722
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 125:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:726
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:732
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 127:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:738
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:742
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 129:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:748
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 130:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:754
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].identifier, Value: yyDollar[3].expression}
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
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:770
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 134:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:775
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 135:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:782
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 136:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:788
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 137:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 138:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:798
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 139:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:802
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:808
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 141:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:812
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 142:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:818
		{
			yyVAL.expression = nil
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:822
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:826
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 145:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:830
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].identifier}
		}
	case 146:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:834
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].identifier}
		}
	case 147:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:840
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].identifier}}
		}
	case 148:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:844
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 149:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:850
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].identifier, To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:856
		{
			yyVAL.expression = Commit{Literal: yyDollar[1].token.Literal}
		}
	case 151:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:862
		{
			yyVAL.expression = Rollback{Literal: yyDollar[1].token.Literal}
		}
	case 152:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:868
		{
			yyVAL.expression = Print{Print: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 153:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:874
		{
			yyVAL.expression = SetFlag{Set: yyDollar[1].token.Literal, Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:880
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:886
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:892
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:896
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:903
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 159:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:907
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:914
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:920
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:926
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:932
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:938
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:944
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 166:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:948
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 167:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:954
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 168:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:960
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 169:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:964
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 170:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:970
		{
			yyVAL.token = Token{}
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:974
		{
			yyVAL.token = yyDollar[1].token
		}
	case 172:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:980
		{
			yyVAL.token = Token{}
		}
	case 173:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:984
		{
			yyVAL.token = yyDollar[1].token
		}
	case 174:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:990
		{
			yyVAL.token = Token{}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:994
		{
			yyVAL.token = yyDollar[1].token
		}
	case 176:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1000
		{
			yyVAL.token = Token{}
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1004
		{
			yyVAL.token = yyDollar[1].token
		}
	case 178:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1010
		{
			yyVAL.token = Token{}
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1014
		{
			yyVAL.token = yyDollar[1].token
		}
	case 180:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1018
		{
			yyVAL.token = yyDollar[1].token
		}
	case 181:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1022
		{
			yyVAL.token = yyDollar[1].token
		}
	case 182:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1028
		{
			yyVAL.token = yyDollar[1].token
		}
	case 183:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1032
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 184:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1038
		{
			yyVAL.token = Token{}
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1042
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
