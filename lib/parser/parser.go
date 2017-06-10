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
const COMMIT = 57416
const ROLLBACK = 57417
const PRINT = 57418
const VAR = 57419
const COMPARISON_OP = 57420
const STRING_OP = 57421
const SUBSTITUTION_OP = 57422

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

//line parser.y:1035

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
	-1, 59,
	57, 170,
	61, 170,
	62, 170,
	-2, 151,
	-1, 95,
	43, 172,
	45, 176,
	-2, 111,
	-1, 129,
	88, 74,
	-2, 168,
	-1, 133,
	68, 104,
	-2, 170,
	-1, 138,
	36, 74,
	73, 74,
	88, 74,
	-2, 168,
	-1, 142,
	57, 170,
	61, 170,
	62, 170,
	-2, 98,
	-1, 146,
	57, 170,
	61, 170,
	62, 170,
	-2, 21,
	-1, 148,
	45, 176,
	-2, 172,
	-1, 173,
	57, 170,
	61, 170,
	62, 170,
	-2, 164,
	-1, 236,
	57, 170,
	61, 170,
	62, 170,
	-2, 107,
	-1, 246,
	57, 170,
	61, 170,
	62, 170,
	-2, 25,
	-1, 259,
	57, 170,
	61, 170,
	62, 170,
	-2, 129,
	-1, 275,
	71, 106,
	-2, 170,
	-1, 297,
	57, 170,
	61, 170,
	62, 170,
	-2, 138,
	-1, 303,
	68, 119,
	70, 119,
	71, 119,
	-2, 170,
	-1, 307,
	57, 170,
	61, 170,
	62, 170,
	-2, 49,
	-1, 310,
	57, 170,
	61, 170,
	62, 170,
	-2, 96,
}

const yyPrivate = 57344

const yyLast = 671

var yyAct = [...]int{

	236, 212, 305, 266, 282, 241, 261, 73, 161, 64,
	204, 235, 209, 139, 77, 75, 44, 149, 267, 147,
	238, 92, 94, 189, 95, 131, 3, 59, 3, 55,
	118, 96, 47, 152, 31, 153, 154, 155, 150, 318,
	296, 148, 256, 117, 253, 278, 215, 196, 112, 60,
	48, 48, 52, 105, 324, 100, 317, 50, 301, 114,
	277, 100, 119, 120, 121, 122, 123, 298, 49, 49,
	295, 27, 130, 289, 108, 99, 101, 50, 151, 158,
	133, 260, 135, 255, 117, 233, 50, 192, 211, 311,
	134, 142, 136, 230, 146, 49, 103, 102, 168, 106,
	138, 49, 129, 110, 111, 121, 122, 123, 83, 85,
	216, 74, 58, 113, 173, 83, 174, 175, 274, 100,
	182, 183, 184, 185, 186, 187, 188, 91, 194, 160,
	166, 194, 167, 239, 117, 128, 86, 48, 164, 205,
	102, 286, 172, 251, 249, 180, 159, 156, 178, 179,
	181, 206, 281, 193, 163, 49, 201, 165, 279, 49,
	223, 100, 195, 200, 199, 145, 203, 242, 202, 220,
	218, 270, 268, 152, 207, 153, 154, 155, 214, 54,
	229, 53, 232, 197, 46, 84, 225, 226, 219, 217,
	231, 88, 309, 228, 227, 240, 137, 142, 50, 190,
	246, 49, 93, 104, 245, 50, 45, 208, 109, 213,
	243, 100, 247, 57, 250, 157, 100, 259, 213, 221,
	222, 224, 107, 254, 258, 269, 248, 90, 169, 170,
	152, 252, 153, 154, 155, 150, 257, 171, 148, 50,
	275, 27, 272, 262, 263, 264, 265, 244, 51, 290,
	27, 49, 304, 50, 70, 12, 49, 12, 292, 100,
	1, 100, 116, 29, 319, 163, 288, 11, 56, 297,
	221, 76, 300, 213, 285, 72, 287, 303, 15, 14,
	307, 291, 13, 10, 310, 302, 9, 69, 23, 308,
	23, 8, 313, 7, 6, 162, 100, 5, 210, 49,
	316, 49, 315, 314, 4, 237, 213, 132, 66, 140,
	141, 312, 98, 322, 293, 294, 97, 136, 78, 307,
	299, 323, 50, 82, 83, 85, 65, 86, 87, 28,
	68, 62, 125, 124, 128, 67, 49, 117, 63, 306,
	61, 280, 198, 144, 16, 2, 221, 0, 0, 0,
	0, 0, 126, 115, 0, 127, 119, 120, 121, 122,
	123, 213, 0, 273, 50, 82, 83, 85, 0, 86,
	87, 28, 177, 176, 80, 0, 0, 0, 81, 0,
	0, 0, 88, 0, 0, 79, 50, 82, 83, 85,
	89, 86, 87, 28, 27, 0, 0, 0, 0, 0,
	0, 84, 0, 0, 0, 71, 0, 50, 82, 83,
	85, 0, 86, 87, 28, 0, 80, 0, 0, 0,
	81, 0, 0, 0, 88, 0, 0, 79, 50, 82,
	83, 85, 89, 86, 87, 28, 0, 0, 80, 0,
	0, 0, 81, 84, 143, 0, 88, 71, 0, 79,
	0, 0, 0, 0, 89, 0, 0, 0, 0, 80,
	0, 0, 0, 81, 0, 84, 0, 88, 0, 71,
	79, 0, 0, 0, 0, 89, 0, 0, 0, 0,
	80, 0, 0, 0, 81, 0, 84, 234, 88, 0,
	71, 79, 125, 124, 128, 0, 89, 117, 0, 0,
	0, 0, 0, 0, 0, 320, 321, 84, 0, 0,
	0, 71, 126, 115, 0, 127, 119, 120, 121, 122,
	123, 0, 191, 125, 124, 128, 0, 0, 117, 0,
	0, 0, 0, 0, 125, 124, 128, 0, 0, 117,
	0, 0, 0, 126, 115, 276, 127, 119, 120, 121,
	122, 123, 0, 0, 126, 115, 0, 127, 119, 120,
	121, 122, 123, 125, 124, 128, 0, 0, 117, 0,
	0, 0, 271, 124, 128, 0, 0, 117, 0, 0,
	0, 0, 0, 126, 115, 0, 127, 119, 120, 121,
	122, 123, 126, 115, 0, 127, 119, 120, 121, 122,
	123, 125, 0, 128, 28, 27, 117, 18, 0, 19,
	0, 17, 128, 0, 0, 117, 0, 20, 0, 0,
	21, 126, 115, 0, 127, 119, 120, 121, 122, 123,
	126, 115, 0, 127, 119, 120, 121, 122, 123, 152,
	0, 153, 154, 155, 150, 283, 284, 148, 30, 0,
	0, 0, 0, 32, 33, 34, 35, 36, 37, 38,
	39, 40, 41, 42, 43, 0, 0, 24, 25, 26,
	22,
}
var yyPact = [...]int{

	593, -1000, 593, -56, -56, -56, -56, -56, -56, -56,
	-56, -56, -56, -56, -56, -56, 193, 165, 249, 235,
	153, 151, 202, 32, -1000, -1000, 424, 62, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 185, 53, 249, 188, -36, 201,
	-1000, 53, 195, 249, 249, -1000, -41, 33, 424, 505,
	15, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 32,
	-1000, 382, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 424,
	10, 424, -1000, -1000, 102, -1000, -1000, -1000, -1000, 13,
	360, -1000, 128, 424, -1000, -11, -1000, 194, -1000, -1000,
	-1000, -1000, 238, 59, 249, 249, -1000, 249, 185, 53,
	11, 203, 202, 424, 505, 424, 318, 75, 88, 424,
	424, 424, 424, 424, 424, 424, -1000, -1000, -1000, 62,
	434, -1, 60, 505, -1000, 552, -1000, -1000, 62, -1000,
	-42, 162, 505, -1000, 126, 124, 505, 113, 129, 94,
	108, 53, -1000, -1000, -1000, -1000, -1000, 249, 1, 249,
	-1000, 193, -43, 29, -1000, -1000, -1000, 185, 249, 82,
	73, 249, -1000, 505, -20, 505, 10, 10, 127, 424,
	6, 424, 21, 21, 71, 71, 71, 543, 552, -3,
	403, -1000, -1000, 63, 424, 131, 360, 249, 131, 424,
	424, 53, 101, 94, 100, -1000, 53, -1000, -1000, -1000,
	-45, 424, -5, -47, 185, 249, 424, -1000, -7, 214,
	249, 139, -1000, 249, 136, -1000, -1000, -1000, -1000, 514,
	382, -1000, 505, -1000, -1000, -1000, 274, 47, 60, 424,
	476, -28, 119, -1000, -1000, 110, 505, -1000, 595, 53,
	98, 53, 186, 1, -15, 229, 249, -1000, -1000, 505,
	-1000, -1000, -1000, -1000, 249, 249, -18, -49, 424, -21,
	249, 424, -30, 424, -1000, 505, 424, -1000, 247, 424,
	-1000, 109, -1000, 424, 2, 186, 53, 595, -1000, -1000,
	1, -1000, -1000, -1000, -1000, 214, 249, 505, -1000, -1000,
	552, -1000, -1000, 505, -32, -1000, -50, 465, -1000, 109,
	505, 249, 186, -1000, -1000, -1000, -1000, -1000, 424, -1000,
	-1000, -1000, -34, -1000, -1000,
}
var yyPgo = [...]int{

	0, 260, 345, 25, 344, 16, 21, 343, 342, 5,
	341, 340, 0, 339, 9, 338, 335, 331, 330, 326,
	23, 318, 31, 316, 24, 312, 4, 310, 309, 308,
	307, 305, 11, 2, 22, 32, 1, 13, 20, 304,
	298, 12, 297, 295, 8, 294, 293, 291, 18, 3,
	6, 286, 283, 282, 279, 278, 49, 275, 7, 111,
	15, 271, 14, 287, 254, 268, 267, 29, 199, 30,
	264, 19, 10, 17, 262, 648,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 3, 4, 5, 5,
	6, 6, 7, 7, 8, 8, 9, 9, 10, 10,
	11, 11, 11, 11, 11, 11, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 13, 70,
	70, 70, 14, 15, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 17, 17, 17, 17, 17, 18,
	18, 18, 19, 19, 20, 20, 20, 21, 21, 22,
	22, 22, 23, 23, 24, 24, 24, 24, 24, 24,
	25, 25, 25, 25, 25, 26, 26, 26, 27, 27,
	28, 28, 29, 30, 30, 31, 31, 32, 32, 33,
	33, 34, 34, 35, 35, 36, 36, 37, 37, 38,
	38, 39, 39, 39, 39, 40, 41, 41, 42, 43,
	44, 44, 45, 45, 46, 47, 47, 48, 48, 49,
	49, 50, 50, 50, 50, 50, 51, 51, 52, 53,
	54, 55, 56, 57, 58, 58, 59, 59, 60, 61,
	62, 63, 64, 65, 65, 66, 67, 67, 68, 68,
	69, 69, 71, 71, 72, 72, 73, 73, 73, 73,
	74, 74, 75, 75,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 7, 3, 0, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 3, 2, 0,
	1, 1, 3, 3, 3, 4, 4, 6, 6, 4,
	4, 4, 4, 2, 3, 3, 3, 3, 3, 3,
	3, 2, 4, 1, 0, 2, 2, 5, 7, 1,
	2, 3, 1, 1, 1, 1, 2, 3, 1, 1,
	5, 5, 6, 6, 4, 0, 2, 4, 1, 1,
	1, 3, 5, 0, 1, 0, 2, 1, 3, 1,
	3, 1, 3, 1, 3, 1, 3, 1, 3, 4,
	2, 5, 8, 4, 7, 3, 1, 3, 6, 3,
	1, 3, 4, 5, 6, 6, 8, 1, 3, 1,
	3, 0, 1, 1, 2, 2, 5, 7, 7, 1,
	1, 2, 1, 1, 1, 2, 1, 2, 1, 1,
	1, 1, 3, 1, 3, 2, 1, 3, 0, 1,
	0, 1, 0, 1, 0, 1, 0, 1, 1, 1,
	1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -39, -42, -45, -46, -47, -51,
	-52, -66, -64, -53, -54, -55, -4, 18, 14, 16,
	24, 27, 77, -63, 74, 75, 76, 12, 11, -1,
	-75, 90, -75, -75, -75, -75, -75, -75, -75, -75,
	-75, -75, -75, -75, -5, 13, 19, -35, -22, -56,
	4, 13, -35, 28, 28, -67, -65, 11, 80, -12,
	-56, -11, -17, -15, -14, -19, -29, -16, -18, -63,
	-64, 87, -57, -58, -59, -60, -61, -62, -21, 67,
	56, 60, 5, 6, 83, 7, 9, 10, 64, 72,
	-68, 65, -6, 17, -34, -24, -22, -23, -25, 22,
	-14, 23, 87, -56, 15, 89, -56, 21, -34, 13,
	-56, -56, 89, 80, -12, 79, -74, 63, -69, 82,
	83, 84, 85, 86, 59, 58, 78, 81, 60, 87,
	-12, -3, -30, -12, -14, -12, -58, -59, 87, -37,
	-28, -27, -12, 84, -7, 37, -12, -71, 52, -73,
	49, 89, 44, 46, 47, 48, -56, 21, 20, 87,
	-3, -44, -43, -56, -35, -56, -6, -34, 87, 25,
	26, 34, -67, -12, -12, -12, 55, 54, -69, 61,
	57, 62, -12, -12, -12, -12, -12, -12, -12, -20,
	-68, 88, 88, -38, 68, -20, 89, 21, -8, 38,
	39, 43, -71, -73, -72, 45, 43, -34, -56, -41,
	-40, 87, -36, -56, -5, 89, 81, -6, -36, -48,
	87, -56, -56, 87, -56, -14, -14, -60, -62, -12,
	87, -14, -12, 88, 84, -32, -12, -31, -38, 70,
	-12, -9, 36, -37, -56, -9, -12, -32, -24, 43,
	-72, 43, -24, 89, -32, 88, 89, -6, -44, -12,
	88, -50, 29, 30, 31, 32, -49, -48, 33, -36,
	35, 58, -32, 89, 71, -12, 69, 88, 73, 39,
	-10, 42, -26, 50, 51, -24, 43, -24, -41, 88,
	20, -3, -36, -56, -56, 88, 89, -12, 88, -56,
	-12, 88, -32, -12, 5, -33, -13, -12, -58, 83,
	-12, 87, -24, -26, -41, -50, -49, 88, 89, -70,
	40, 41, -36, -33, 88,
}
var yyDef = [...]int{

	1, -2, 1, 182, 182, 182, 182, 182, 182, 182,
	182, 182, 182, 182, 182, 182, 18, 0, 0, 0,
	0, 0, 0, 0, 149, 150, 0, 168, 161, 2,
	3, 183, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 20, 0, 0, 0, 113, 79,
	152, 0, 0, 0, 0, 165, 166, 163, 0, -2,
	36, 37, 38, 39, 40, 41, 42, 43, 44, 45,
	46, 0, 30, 31, 32, 33, 34, 35, 73, 103,
	0, 0, 153, 154, 0, 156, 158, 159, 160, 0,
	0, 169, 22, 0, 19, -2, 84, 85, 88, 89,
	82, 83, 0, 0, 0, 0, 80, 0, 20, 0,
	0, 0, 0, 0, 162, 0, 0, 170, 0, 0,
	0, 0, 0, 0, 0, 0, 180, 181, 171, -2,
	170, 0, 0, -2, 63, 71, 155, 157, -2, 17,
	117, 100, -2, 99, 24, 0, -2, 0, -2, 174,
	0, 0, 173, 177, 178, 179, 86, 0, 0, 0,
	123, 18, 130, 0, 114, 81, 132, 20, 0, 0,
	0, 0, 167, -2, 53, 54, 0, 0, 0, 0,
	0, 0, 64, 65, 66, 67, 68, 69, 70, 0,
	0, 47, 52, 105, 0, 26, 0, 0, 26, 0,
	0, 0, 0, 174, 0, 175, 0, 112, 87, 121,
	126, 0, 0, 115, 20, 0, 0, 133, 0, 141,
	0, 137, 146, 0, 0, 61, 62, 55, 56, 170,
	0, 59, 60, 72, 75, 76, -2, 0, 120, 0,
	170, 0, 0, 118, 101, 28, -2, 23, 95, 0,
	0, 0, 94, 0, 0, 0, 0, 128, 131, -2,
	134, 135, 142, 143, 0, 0, 0, 139, 0, 0,
	0, 0, 0, 0, 102, -2, 0, 77, 0, 0,
	16, 0, 90, 0, 0, 91, 0, 95, 127, 125,
	0, 124, 116, 144, 145, 141, 0, -2, 147, 148,
	57, 58, 108, -2, 0, 27, 109, -2, 29, 0,
	-2, 0, 93, 92, 122, 136, 140, 78, 0, 48,
	50, 51, 0, 110, 97,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 86, 3, 3,
	87, 88, 84, 82, 89, 83, 3, 85, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 90,
	3, 81,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 79, 80,
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
		//line parser.y:124
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:129
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:136
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:140
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:144
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:148
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:152
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:156
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:160
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:164
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:168
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:172
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:176
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:180
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:184
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 16:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:190
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
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:204
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 18:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:210
		{
			yyVAL.expression = nil
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:214
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:220
		{
			yyVAL.expression = nil
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:224
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 22:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:230
		{
			yyVAL.expression = nil
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:234
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 24:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:240
		{
			yyVAL.expression = nil
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:244
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 26:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:250
		{
			yyVAL.expression = nil
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:254
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 28:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:260
		{
			yyVAL.expression = nil
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:264
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:270
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:274
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:278
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:282
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:286
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:290
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:296
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:300
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:304
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:308
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:312
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:316
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:320
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:324
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:328
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:332
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:336
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:340
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:346
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 49:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:352
		{
			yyVAL.token = Token{}
		}
	case 50:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:356
		{
			yyVAL.token = yyDollar[1].token
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:360
		{
			yyVAL.token = yyDollar[1].token
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:366
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:372
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
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:399
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 56:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:403
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 57:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:407
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 58:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:411
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:415
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:419
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 61:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:423
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:427
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:431
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:437
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:441
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:445
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:449
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:453
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:459
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:463
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:467
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 72:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:473
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:477
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:483
		{
			yyVAL.expression = Option{}
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:487
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:491
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 77:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:497
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 78:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:501
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:507
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:511
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:515
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:521
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:525
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:531
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:535
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:539
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:543
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:547
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:551
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:557
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 91:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:561
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 92:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:565
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 93:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:569
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 94:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:573
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 95:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:579
		{
			yyVAL.expression = nil
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:583
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:587
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:593
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:597
		{
			yyVAL.expression = AllColumns{}
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:603
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:607
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 102:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:613
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 103:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expression = nil
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:623
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 105:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:629
		{
			yyVAL.expression = nil
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:633
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:639
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:643
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:649
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:653
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:659
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:663
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:669
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 114:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:673
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:679
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:683
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:689
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 118:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:693
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 119:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:699
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 120:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:703
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 121:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:709
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 122:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:713
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 123:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:717
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 124:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:721
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:727
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:733
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:737
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 128:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:743
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 129:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:749
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].identifier, Value: yyDollar[3].expression}
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:755
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 131:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:759
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 132:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:765
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 133:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:770
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 134:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:777
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 135:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:783
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 136:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:787
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:793
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 138:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:797
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 139:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:803
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 140:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:807
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 141:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:813
		{
			yyVAL.expression = nil
		}
	case 142:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:817
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:821
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 144:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:825
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].identifier}
		}
	case 145:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:829
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].identifier}
		}
	case 146:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:835
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].identifier}}
		}
	case 147:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:839
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 148:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:845
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].identifier, To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 149:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:851
		{
			yyVAL.expression = Commit{Literal: yyDollar[1].token.Literal}
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expression = Rollback{Literal: yyDollar[1].token.Literal}
		}
	case 151:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:863
		{
			yyVAL.expression = Print{Print: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:869
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 153:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:875
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:881
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 155:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:885
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:892
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:896
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:903
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:909
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:915
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:921
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:927
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:933
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:937
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 165:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:943
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:949
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 167:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:953
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 168:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:959
		{
			yyVAL.token = Token{}
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:963
		{
			yyVAL.token = yyDollar[1].token
		}
	case 170:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:969
		{
			yyVAL.token = Token{}
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:973
		{
			yyVAL.token = yyDollar[1].token
		}
	case 172:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:979
		{
			yyVAL.token = Token{}
		}
	case 173:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:983
		{
			yyVAL.token = yyDollar[1].token
		}
	case 174:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:989
		{
			yyVAL.token = Token{}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:993
		{
			yyVAL.token = yyDollar[1].token
		}
	case 176:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:999
		{
			yyVAL.token = Token{}
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1003
		{
			yyVAL.token = yyDollar[1].token
		}
	case 178:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1007
		{
			yyVAL.token = yyDollar[1].token
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1011
		{
			yyVAL.token = yyDollar[1].token
		}
	case 180:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1017
		{
			yyVAL.token = yyDollar[1].token
		}
	case 181:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1021
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 182:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1027
		{
			yyVAL.token = Token{}
		}
	case 183:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1031
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
