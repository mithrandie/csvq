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
	procexpr    ProcExpr
	procexprs   []ProcExpr
	primary     Primary
	identifier  Identifier
	text        String
	integer     Integer
	float       Float
	ternary     Ternary
	datetime    Datetime
	null        Null
	variable    Variable
	variables   []Variable
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
const INTERSECT = 57397
const EXCEPT = 57398
const ALL = 57399
const ANY = 57400
const EXISTS = 57401
const IN = 57402
const AND = 57403
const OR = 57404
const NOT = 57405
const BETWEEN = 57406
const LIKE = 57407
const IS = 57408
const NULL = 57409
const DISTINCT = 57410
const WITH = 57411
const CASE = 57412
const IF = 57413
const ELSEIF = 57414
const WHILE = 57415
const WHEN = 57416
const THEN = 57417
const ELSE = 57418
const DO = 57419
const END = 57420
const DECLARE = 57421
const CURSOR = 57422
const FOR = 57423
const FETCH = 57424
const OPEN = 57425
const CLOSE = 57426
const DISPOSE = 57427
const GROUP_CONCAT = 57428
const SEPARATOR = 57429
const COMMIT = 57430
const ROLLBACK = 57431
const CONTINUE = 57432
const BREAK = 57433
const EXIT = 57434
const PRINT = 57435
const VAR = 57436
const COMPARISON_OP = 57437
const STRING_OP = 57438
const SUBSTITUTION_OP = 57439

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
	"INTERSECT",
	"EXCEPT",
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
	"IF",
	"ELSEIF",
	"WHILE",
	"WHEN",
	"THEN",
	"ELSE",
	"DO",
	"END",
	"DECLARE",
	"CURSOR",
	"FOR",
	"FETCH",
	"OPEN",
	"CLOSE",
	"DISPOSE",
	"GROUP_CONCAT",
	"SEPARATOR",
	"COMMIT",
	"ROLLBACK",
	"CONTINUE",
	"BREAK",
	"EXIT",
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

//line parser.y:1277

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
	-1, 16,
	54, 45,
	55, 45,
	56, 45,
	-2, 56,
	-1, 112,
	60, 211,
	64, 211,
	65, 211,
	-2, 225,
	-1, 131,
	44, 213,
	46, 217,
	-2, 146,
	-1, 165,
	54, 46,
	55, 46,
	56, 46,
	-2, 72,
	-1, 167,
	106, 107,
	-2, 209,
	-1, 169,
	74, 137,
	-2, 211,
	-1, 174,
	37, 107,
	87, 107,
	106, 107,
	-2, 209,
	-1, 191,
	60, 211,
	64, 211,
	65, 211,
	-2, 131,
	-1, 198,
	60, 211,
	64, 211,
	65, 211,
	-2, 81,
	-1, 210,
	46, 217,
	-2, 213,
	-1, 226,
	60, 211,
	64, 211,
	65, 211,
	-2, 206,
	-1, 232,
	66, 0,
	95, 0,
	98, 0,
	-2, 86,
	-1, 233,
	66, 0,
	95, 0,
	98, 0,
	-2, 87,
	-1, 266,
	60, 211,
	64, 211,
	65, 211,
	-2, 51,
	-1, 313,
	66, 0,
	95, 0,
	98, 0,
	-2, 93,
	-1, 319,
	60, 211,
	64, 211,
	65, 211,
	-2, 142,
	-1, 344,
	60, 211,
	64, 211,
	65, 211,
	-2, 164,
	-1, 372,
	78, 139,
	-2, 211,
	-1, 379,
	60, 211,
	64, 211,
	65, 211,
	-2, 55,
	-1, 397,
	60, 211,
	64, 211,
	65, 211,
	-2, 173,
	-1, 406,
	74, 154,
	76, 154,
	78, 154,
	-2, 211,
	-1, 410,
	90, 18,
	91, 18,
	-2, 1,
	-1, 413,
	60, 211,
	64, 211,
	65, 211,
	-2, 129,
}

const yyPrivate = 57344

const yyLast = 1046

var yyAct = [...]int{

	80, 40, 432, 40, 423, 277, 44, 294, 357, 304,
	352, 46, 47, 48, 49, 50, 51, 52, 43, 1,
	86, 23, 386, 23, 196, 180, 288, 131, 188, 280,
	67, 68, 69, 202, 53, 113, 94, 92, 365, 211,
	110, 2, 358, 40, 318, 37, 209, 90, 77, 249,
	321, 154, 64, 45, 254, 108, 183, 56, 396, 351,
	130, 136, 132, 150, 59, 118, 156, 157, 158, 159,
	160, 141, 341, 85, 38, 177, 38, 61, 145, 146,
	147, 57, 57, 135, 137, 39, 375, 338, 165, 177,
	283, 16, 273, 383, 271, 142, 127, 436, 170, 422,
	404, 319, 398, 395, 3, 374, 109, 156, 157, 158,
	159, 160, 382, 350, 59, 136, 39, 340, 316, 179,
	162, 161, 163, 39, 199, 153, 40, 193, 166, 167,
	414, 279, 42, 75, 107, 59, 311, 112, 221, 136,
	214, 59, 215, 216, 217, 212, 174, 121, 210, 172,
	40, 166, 39, 91, 151, 150, 208, 152, 156, 157,
	158, 159, 160, 284, 178, 42, 118, 45, 230, 227,
	23, 194, 143, 185, 186, 182, 40, 42, 158, 159,
	160, 144, 438, 430, 40, 206, 40, 40, 100, 164,
	57, 411, 121, 100, 102, 225, 23, 401, 169, 220,
	171, 40, 213, 371, 228, 234, 363, 325, 200, 330,
	252, 427, 426, 426, 136, 42, 425, 256, 376, 251,
	187, 191, 229, 38, 253, 307, 198, 263, 40, 267,
	201, 269, 270, 262, 441, 303, 300, 312, 282, 314,
	315, 437, 297, 172, 42, 226, 293, 420, 400, 38,
	287, 229, 231, 232, 233, 173, 40, 286, 240, 241,
	242, 243, 244, 245, 246, 296, 328, 329, 103, 305,
	331, 309, 308, 163, 291, 268, 23, 268, 268, 252,
	307, 322, 195, 266, 306, 120, 136, 101, 324, 236,
	176, 136, 184, 235, 237, 239, 238, 256, 337, 116,
	335, 76, 115, 116, 117, 121, 289, 40, 342, 362,
	326, 390, 165, 345, 347, 364, 348, 343, 349, 214,
	346, 215, 216, 217, 339, 366, 105, 23, 290, 38,
	360, 40, 302, 285, 123, 334, 124, 310, 88, 313,
	333, 40, 250, 265, 381, 54, 361, 136, 359, 136,
	63, 23, 222, 223, 323, 148, 369, 58, 58, 392,
	327, 224, 62, 59, 272, 70, 71, 72, 73, 74,
	59, 385, 256, 191, 389, 198, 391, 55, 40, 380,
	38, 219, 119, 408, 126, 378, 344, 181, 129, 412,
	59, 136, 114, 138, 125, 111, 230, 128, 23, 58,
	60, 139, 140, 40, 38, 418, 417, 419, 41, 367,
	66, 40, 424, 407, 416, 405, 121, 59, 415, 410,
	155, 421, 428, 23, 372, 274, 40, 429, 204, 431,
	65, 23, 409, 93, 435, 379, 353, 354, 355, 356,
	40, 89, 10, 9, 440, 384, 23, 8, 443, 7,
	6, 38, 203, 58, 214, 5, 215, 216, 217, 212,
	23, 397, 210, 278, 4, 205, 58, 256, 207, 320,
	403, 168, 218, 82, 189, 406, 38, 58, 190, 134,
	433, 256, 133, 95, 38, 81, 84, 78, 83, 413,
	79, 197, 122, 332, 442, 264, 36, 15, 257, 38,
	14, 13, 281, 12, 11, 248, 255, 162, 161, 163,
	0, 0, 153, 38, 0, 261, 41, 0, 39, 0,
	18, 34, 19, 0, 17, 299, 301, 0, 434, 0,
	20, 0, 162, 21, 163, 0, 0, 153, 0, 205,
	0, 151, 150, 0, 152, 156, 157, 158, 159, 160,
	0, 0, 58, 370, 0, 0, 0, 0, 292, 0,
	295, 298, 205, 205, 0, 0, 151, 150, 0, 152,
	156, 157, 158, 159, 160, 0, 258, 0, 32, 0,
	163, 0, 0, 153, 26, 204, 0, 30, 27, 28,
	29, 0, 0, 24, 25, 259, 260, 33, 35, 22,
	0, 0, 281, 99, 100, 102, 0, 103, 104, 0,
	42, 336, 151, 150, 0, 152, 156, 157, 158, 159,
	160, 0, 205, 0, 58, 0, 0, 0, 0, 58,
	0, 0, 0, 0, 0, 0, 298, 0, 0, 205,
	0, 0, 41, 281, 39, 0, 18, 34, 19, 0,
	17, 0, 0, 0, 0, 0, 20, 393, 394, 21,
	0, 0, 0, 0, 0, 105, 0, 59, 99, 100,
	102, 0, 103, 104, 41, 0, 39, 0, 0, 0,
	205, 0, 0, 0, 0, 58, 0, 58, 0, 0,
	295, 0, 0, 0, 205, 205, 0, 0, 101, 0,
	399, 0, 31, 0, 32, 0, 0, 0, 0, 0,
	26, 0, 0, 30, 27, 28, 29, 0, 0, 24,
	25, 0, 97, 33, 35, 22, 98, 0, 0, 58,
	105, 0, 0, 96, 0, 298, 42, 59, 99, 100,
	102, 0, 103, 104, 41, 0, 0, 0, 0, 106,
	0, 0, 0, 295, 59, 99, 100, 102, 0, 103,
	104, 41, 0, 101, 0, 0, 0, 214, 87, 215,
	216, 217, 212, 387, 388, 210, 59, 99, 100, 102,
	0, 103, 104, 41, 0, 0, 275, 276, 0, 0,
	0, 0, 97, 0, 0, 0, 98, 0, 0, 0,
	105, 0, 0, 96, 0, 0, 162, 161, 163, 97,
	0, 153, 0, 98, 0, 0, 0, 105, 0, 106,
	96, 0, 162, 161, 163, 0, 0, 153, 0, 0,
	0, 97, 0, 101, 192, 98, 106, 0, 87, 105,
	151, 150, 96, 152, 156, 157, 158, 159, 160, 0,
	101, 317, 0, 0, 0, 87, 151, 150, 106, 152,
	156, 157, 158, 159, 160, 0, 0, 247, 162, 161,
	163, 0, 101, 153, 0, 0, 0, 87, 162, 161,
	163, 0, 439, 153, 0, 0, 0, 0, 0, 0,
	0, 0, 402, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 151, 150, 0, 152, 156, 157, 158, 159,
	160, 0, 151, 150, 0, 152, 156, 157, 158, 159,
	160, 162, 161, 163, 0, 0, 153, 0, 0, 0,
	0, 162, 161, 163, 0, 377, 153, 0, 0, 0,
	0, 162, 161, 163, 0, 373, 153, 0, 0, 0,
	0, 0, 0, 0, 0, 151, 150, 175, 152, 156,
	157, 158, 159, 160, 0, 151, 150, 0, 152, 156,
	157, 158, 159, 160, 0, 151, 150, 0, 152, 156,
	157, 158, 159, 160, 162, 161, 163, 0, 0, 153,
	0, 0, 0, 0, 162, 161, 163, 0, 149, 153,
	0, 0, 0, 368, 161, 163, 0, 0, 153, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 151, 150,
	0, 152, 156, 157, 158, 159, 160, 0, 151, 150,
	0, 152, 156, 157, 158, 159, 160, 151, 150, 0,
	152, 156, 157, 158, 159, 160,
}
var yyPact = [...]int{

	631, -1000, 631, -55, -55, -55, -55, -55, -55, -55,
	-55, -1000, -1000, -1000, -1000, -1000, 308, 357, 413, 386,
	333, 321, 399, -55, -55, -55, 413, 413, 413, 413,
	413, 772, 772, -55, 383, 772, 378, 248, 69, 217,
	-1000, -1000, 139, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 291, 296, 413, 368, -11, 366, -1000,
	60, 379, 413, 413, -55, -12, 75, -1000, -1000, -1000,
	101, -55, -55, -55, 335, 923, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 69, -1000, 663, 24, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 772, 27, 772, -1000,
	-1000, 187, -1000, -1000, -1000, -1000, 41, 880, 230, -32,
	-1000, 66, 59, 369, 60, 235, 235, 235, 772, 733,
	-1000, 21, -1000, 182, 772, 103, 413, 413, -1000, 413,
	369, 95, -1000, 359, -1000, -1000, -1000, -1000, 60, 33,
	326, -1000, 399, 772, 88, -1000, -1000, -1000, 397, 631,
	772, 772, 772, 210, 229, 238, 772, 772, 772, 772,
	772, 772, 772, -1000, 761, -1000, 413, 217, 136, 933,
	-1000, 517, -1000, -1000, 217, 505, 413, 397, 598, -1000,
	305, 772, -1000, 139, -1000, 139, 139, 933, -1000, -13,
	342, 933, -1000, -1000, -1000, 182, -1000, -15, 745, 26,
	110, -1000, 378, -17, 65, 47, -1000, -1000, -1000, 289,
	274, 260, 284, 60, -1000, -1000, -1000, -1000, -1000, 413,
	369, 413, 137, 131, 413, -1000, 933, 139, -55, -18,
	208, 8, -33, -33, 259, 772, 31, 772, 27, 27,
	77, 77, -1000, -1000, -1000, 471, 517, -1000, -1000, 12,
	750, 205, 772, 308, 129, 505, -1000, -1000, 772, -55,
	-55, 132, -1000, -55, 301, 295, 933, 244, -1000, -1000,
	244, 733, 413, 772, -1000, -1000, -1000, -1000, -20, 772,
	11, -35, 369, 413, 772, 60, 276, 260, 272, -1000,
	60, -1000, -1000, -1000, 7, -48, 406, 413, 314, -1000,
	413, 310, -55, -1000, 128, 208, 631, 772, -1000, -1000,
	942, 663, -1000, -33, -1000, -1000, -1000, -1000, -1000, 446,
	125, 136, 772, 870, -1, 145, -1000, 860, -1000, -1000,
	505, -1000, -1000, 772, 772, -1000, -1000, -1000, 26, 6,
	72, 413, -1000, -1000, 933, 722, 60, 267, 60, 409,
	-1000, 413, -1000, -1000, -1000, 413, 413, -3, -49, 772,
	-4, 413, -1000, 177, 119, 153, -1000, 817, 772, -6,
	772, -1000, 933, 772, -1000, 408, -55, 505, 113, 933,
	-1000, -1000, -1000, 26, -1000, -1000, -1000, 772, 25, 409,
	60, 722, -1000, -1000, -1000, 406, 413, 933, -1000, -1000,
	-55, 176, 631, 517, -1000, -1000, 933, -7, -1000, 140,
	631, 138, -1000, 933, 413, 409, -1000, -1000, -1000, -1000,
	-55, -1000, -1000, 105, 140, 505, 772, -55, -9, -1000,
	170, 104, 141, -1000, 807, -1000, -1000, -55, 163, 505,
	-1000, -55, -1000, -1000,
}
var yyPgo = [...]int{

	0, 18, 54, 41, 506, 504, 503, 501, 500, 498,
	497, 104, 91, 45, 496, 35, 25, 495, 493, 34,
	492, 48, 301, 101, 491, 0, 490, 488, 487, 486,
	485, 49, 483, 62, 482, 27, 479, 22, 478, 474,
	473, 471, 469, 29, 44, 24, 60, 57, 7, 28,
	50, 464, 463, 5, 455, 452, 33, 450, 449, 447,
	42, 8, 10, 443, 442, 38, 9, 2, 4, 338,
	441, 47, 153, 37, 433, 36, 73, 55, 20, 430,
	52, 342, 51, 425, 46, 26, 39, 56, 420, 6,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10,
	11, 12, 12, 12, 12, 13, 13, 14, 15, 15,
	16, 16, 17, 17, 18, 18, 19, 19, 20, 20,
	21, 21, 21, 21, 21, 21, 22, 22, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 23, 23,
	24, 83, 83, 83, 25, 26, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 28, 28, 28,
	28, 28, 29, 29, 29, 30, 30, 31, 31, 31,
	32, 32, 33, 33, 33, 34, 34, 35, 35, 35,
	35, 35, 35, 36, 36, 36, 36, 36, 37, 37,
	37, 38, 38, 39, 39, 40, 41, 41, 42, 42,
	43, 43, 44, 44, 45, 45, 46, 46, 47, 47,
	48, 48, 49, 49, 50, 50, 51, 51, 51, 51,
	52, 53, 53, 54, 55, 56, 56, 57, 57, 58,
	59, 59, 60, 60, 61, 61, 62, 62, 62, 62,
	62, 63, 63, 64, 65, 65, 66, 66, 67, 67,
	68, 68, 69, 70, 71, 71, 72, 72, 73, 74,
	75, 76, 77, 77, 78, 79, 79, 80, 80, 81,
	81, 82, 82, 84, 84, 85, 85, 86, 86, 86,
	86, 87, 87, 88, 88, 89, 89,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 5, 8,
	9, 7, 9, 2, 8, 9, 2, 2, 5, 3,
	3, 5, 4, 4, 4, 1, 1, 3, 0, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 2,
	1, 1, 1, 1, 1, 1, 1, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	2, 0, 1, 1, 3, 3, 3, 3, 4, 4,
	6, 6, 4, 4, 4, 4, 2, 3, 3, 3,
	3, 3, 3, 3, 2, 4, 1, 0, 2, 2,
	5, 7, 1, 2, 3, 1, 1, 1, 1, 2,
	3, 1, 1, 5, 5, 6, 6, 4, 0, 2,
	4, 1, 1, 1, 3, 5, 0, 1, 0, 2,
	1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
	1, 3, 1, 3, 4, 2, 5, 8, 4, 7,
	3, 1, 3, 6, 3, 1, 3, 4, 5, 6,
	6, 8, 1, 3, 1, 3, 0, 1, 1, 2,
	2, 5, 7, 7, 4, 2, 0, 2, 4, 2,
	0, 2, 1, 1, 1, 2, 1, 2, 1, 1,
	1, 1, 1, 3, 3, 1, 3, 1, 3, 0,
	1, 0, 1, 0, 1, 0, 1, 0, 1, 1,
	1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -51, -54, -57, -58, -59, -63,
	-64, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 94, -78, 88, 89, 79, 83, 84, 85,
	82, 71, 73, 92, 16, 93, -14, -13, -76, 13,
	-25, 11, 105, -1, -89, 108, -89, -89, -89, -89,
	-89, -89, -89, -19, 37, 20, -47, -33, -69, 4,
	14, -47, 29, 29, -80, -79, 11, -89, -89, -89,
	-69, -69, -69, -69, -69, -23, -22, -21, -28, -26,
	-25, -30, -40, -27, -29, -76, -78, 105, -69, -70,
	-71, -72, -73, -74, -75, -32, 70, 59, 63, 5,
	6, 100, 7, 9, 10, 67, 86, -23, -77, -76,
	-89, 12, -23, -15, 14, 54, 55, 56, 97, -81,
	68, -11, -20, 43, 40, -69, 16, 107, -69, 22,
	-46, -35, -33, -34, -36, 23, -25, 24, 14, -69,
	-69, -89, 107, 97, 80, -89, -89, -89, 20, 75,
	96, 95, 98, 66, -82, -88, 99, 100, 101, 102,
	103, 62, 61, 63, -23, -25, 104, 105, -41, -23,
	-25, -23, -71, -72, 105, 77, 60, 107, 98, -89,
	-16, 18, -46, -87, 57, -87, -87, -23, -49, -39,
	-38, -23, 101, 106, -71, 100, -45, -24, -23, 21,
	105, -11, -56, -55, -22, -69, -47, -69, -16, -84,
	53, -86, 50, 107, 45, 47, 48, 49, -69, 22,
	-46, 105, 26, 27, 35, -80, -23, 81, -77, -76,
	-1, -23, -23, -23, -82, 64, 60, 65, 58, 57,
	-23, -23, -23, -23, -23, -23, -23, 106, -69, -31,
	-81, -50, 74, -31, -2, -4, -3, -9, 71, 90,
	91, -69, -77, -21, -17, 38, -23, -13, -12, -13,
	-13, 107, 22, 107, -83, 41, 42, -53, -52, 105,
	-43, -22, -15, 107, 98, 44, -84, -86, -85, 46,
	44, -46, -69, -16, -48, -69, -60, 105, -69, -22,
	105, -22, -11, -89, -66, -65, 76, 72, -73, -75,
	-23, 105, -25, -23, -25, -25, 106, 101, -44, -23,
	-42, -50, 76, -23, -19, 78, -2, -23, -89, -89,
	77, -89, -18, 39, 40, -49, -69, -45, 107, -44,
	106, 107, -16, -56, -23, -35, 44, -85, 44, -35,
	106, 107, -62, 30, 31, 32, 33, -61, -60, 34,
	-43, 36, -89, 78, -66, -65, -1, -23, 61, -44,
	107, 78, -23, 75, 106, 87, 73, 75, -2, -23,
	-44, -53, 106, 21, -11, -43, -37, 51, 52, -35,
	44, -35, -48, -22, -22, 106, 107, -23, 106, -69,
	71, 78, 75, -23, 106, -44, -23, 5, -89, -2,
	-3, 78, -53, -23, 105, -35, -37, -62, -61, -89,
	71, -1, 106, -68, -67, 76, 72, 73, -48, -89,
	78, -68, -67, -2, -23, -89, 106, 71, 78, 75,
	-89, 71, -2, -89,
}
var yyDef = [...]int{

	1, -2, 1, 225, 225, 225, 225, 225, 225, 225,
	225, 13, 14, 15, 16, 17, -2, 0, 0, 0,
	0, 0, 0, 225, 225, 225, 0, 0, 0, 0,
	0, 0, 0, 225, 0, 0, 48, 0, 0, 209,
	46, 201, 0, 2, 5, 226, 6, 7, 8, 9,
	10, 11, 12, 58, 0, 0, 0, 148, 112, 192,
	0, 0, 0, 0, 225, 207, 205, 21, 22, 23,
	0, 225, 225, 225, 0, 211, 68, 69, 70, 71,
	72, 73, 74, 75, 76, 77, 78, 0, 66, 60,
	61, 62, 63, 64, 65, 106, 136, 0, 0, 193,
	194, 0, 196, 198, 199, 200, 0, 211, 0, 77,
	33, 0, -2, 50, 0, 221, 221, 221, 0, 0,
	210, 0, 40, 0, 0, 0, 0, 0, 113, 0,
	50, -2, 117, 118, 121, 122, 115, 116, 0, 0,
	0, 20, 0, 0, 0, 25, 26, 27, 0, 1,
	0, 223, 224, 211, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 212, 211, -2, 0, -2, 0, -2,
	96, 104, 195, 197, -2, 3, 0, 0, 0, 39,
	52, 0, 49, 0, 222, 0, 0, 204, 47, 152,
	133, -2, 132, 84, 59, 0, 57, 144, -2, 0,
	0, 158, 48, 165, 0, 66, 149, 114, 167, 0,
	-2, 215, 0, 0, 214, 218, 219, 220, 119, 0,
	50, 0, 0, 0, 0, 208, -2, 0, 225, 202,
	186, 85, -2, -2, 0, 0, 0, 0, 0, 0,
	97, 98, 99, 100, 101, 102, 103, 79, 67, 0,
	0, 138, 0, 56, 0, 3, 18, 19, 0, 225,
	225, 0, 203, 225, 54, 0, -2, 42, 45, 43,
	44, 0, 0, 0, 80, 82, 83, 156, 161, 0,
	0, 140, 50, 0, 0, 0, 0, 215, 0, 216,
	0, 147, 120, 168, 0, 150, 176, 0, 172, 181,
	0, 0, 225, 28, 0, 186, 1, 0, 88, 89,
	211, 0, 92, -2, 94, 95, 105, 108, 109, -2,
	0, 155, 0, 211, 0, 0, 4, 211, 36, 37,
	3, 38, 41, 0, 0, 153, 134, 145, 0, 0,
	0, 0, 163, 166, -2, 128, 0, 0, 0, 127,
	169, 0, 170, 177, 178, 0, 0, 0, 174, 0,
	0, 0, 24, 0, 0, 185, 187, 211, 0, 0,
	0, 135, -2, 0, 110, 0, 225, 1, 0, -2,
	53, 162, 160, 0, 159, 141, 123, 0, 0, 124,
	0, 128, 151, 179, 180, 176, 0, -2, 182, 183,
	225, 0, 1, 90, 91, 143, -2, 0, 31, 190,
	-2, 0, 157, -2, 0, 126, 125, 171, 175, 29,
	225, 184, 111, 0, 190, 3, 0, 225, 0, 30,
	0, 0, 189, 191, 211, 32, 130, 225, 0, 3,
	34, 225, 188, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 103, 3, 3,
	105, 106, 101, 99, 107, 100, 104, 102, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 108,
	3, 98,
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
	82, 83, 84, 85, 86, 87, 88, 89, 90, 91,
	92, 93, 94, 95, 96, 97,
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
		//line parser.y:147
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:152
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:159
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:164
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:171
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:175
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:179
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:183
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:187
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:191
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:195
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:199
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:203
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:207
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:211
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:215
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:219
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:225
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:229
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:235
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:239
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:245
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:249
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:255
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:259
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:263
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:267
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:271
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:277
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:281
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:285
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:289
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:293
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:299
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:303
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:307
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:311
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:317
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:321
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:327
		{
			yyVAL.expression = SelectQuery{
				SelectEntity:  yyDollar[1].expression,
				OrderByClause: yyDollar[2].expression,
				LimitClause:   yyDollar[3].expression,
			}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:337
		{
			yyVAL.expression = SelectEntity{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
			}
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:347
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 43:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:356
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:365
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:376
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:380
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:386
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:392
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:396
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:402
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:406
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:412
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:422
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:426
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:432
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:436
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 58:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:442
		{
			yyVAL.expression = nil
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:446
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:452
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:456
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:460
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:464
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:468
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:472
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:478
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:482
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:488
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:492
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:496
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:500
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:504
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:508
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:512
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:528
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:538
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 81:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:544
		{
			yyVAL.token = Token{}
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:548
		{
			yyVAL.token = yyDollar[1].token
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.token = yyDollar[1].token
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:558
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:564
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
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:587
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:591
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 88:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:595
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 89:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:599
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 90:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:603
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 91:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:607
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 92:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:611
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 93:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:615
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 94:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 95:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:623
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:627
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:633
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:637
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:641
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:645
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:649
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:655
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:659
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:663
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:669
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:673
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 107:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:679
		{
			yyVAL.expression = Option{}
		}
	case 108:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:683
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:687
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 110:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:693
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 111:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:697
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:703
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 113:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:707
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 114:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:711
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:717
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:721
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:727
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:731
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 119:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:735
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:739
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:743
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 122:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:747
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 123:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:753
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 124:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:757
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 125:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:761
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 126:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:765
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 127:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:769
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 128:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:775
		{
			yyVAL.expression = nil
		}
	case 129:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:779
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 130:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:783
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:789
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:793
		{
			yyVAL.expression = AllColumns{}
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:799
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 134:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:803
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 135:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:809
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 136:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:815
		{
			yyVAL.expression = nil
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:819
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 138:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:825
		{
			yyVAL.expression = nil
		}
	case 139:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:829
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:835
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 141:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:839
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 142:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:845
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:849
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:855
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:859
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:865
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:869
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:875
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 149:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:879
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:885
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 151:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:889
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:895
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 153:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:899
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 154:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:905
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 155:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:909
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 156:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:915
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 157:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:919
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 158:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:923
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 159:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:927
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 160:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:933
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:939
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:943
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 163:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:949
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:955
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:961
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 166:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:965
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 167:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:971
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 168:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:976
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 169:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:983
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 170:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:989
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 171:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:993
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:999
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 173:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1003
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 174:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1009
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 175:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1013
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 176:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1019
		{
			yyVAL.expression = nil
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1023
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 178:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1027
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 179:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1031
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 180:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1035
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 181:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1041
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 182:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1045
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 183:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1051
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 184:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1057
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 185:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1061
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 186:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1067
		{
			yyVAL.procexpr = nil
		}
	case 187:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1071
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 188:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1077
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 189:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1081
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 190:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1087
		{
			yyVAL.procexpr = nil
		}
	case 191:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1091
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 192:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1097
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 193:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1103
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1109
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 195:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1113
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 196:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1120
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 197:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1124
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 198:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1131
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 199:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1137
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1143
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 201:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1149
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1155
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 203:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1159
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 204:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1165
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1171
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1175
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1181
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 208:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1185
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 209:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1191
		{
			yyVAL.token = Token{}
		}
	case 210:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1195
		{
			yyVAL.token = yyDollar[1].token
		}
	case 211:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1201
		{
			yyVAL.token = Token{}
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1205
		{
			yyVAL.token = yyDollar[1].token
		}
	case 213:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1211
		{
			yyVAL.token = Token{}
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1215
		{
			yyVAL.token = yyDollar[1].token
		}
	case 215:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1221
		{
			yyVAL.token = Token{}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1225
		{
			yyVAL.token = yyDollar[1].token
		}
	case 217:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1231
		{
			yyVAL.token = Token{}
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1235
		{
			yyVAL.token = yyDollar[1].token
		}
	case 219:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1239
		{
			yyVAL.token = yyDollar[1].token
		}
	case 220:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1243
		{
			yyVAL.token = yyDollar[1].token
		}
	case 221:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1249
		{
			yyVAL.token = Token{}
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1253
		{
			yyVAL.token = yyDollar[1].token
		}
	case 223:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1259
		{
			yyVAL.token = yyDollar[1].token
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1263
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 225:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1269
		{
			yyVAL.token = Token{}
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1273
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
