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
const IF = 57411
const ELSEIF = 57412
const WHILE = 57413
const WHEN = 57414
const THEN = 57415
const ELSE = 57416
const DO = 57417
const END = 57418
const DECLARE = 57419
const CURSOR = 57420
const FOR = 57421
const FETCH = 57422
const OPEN = 57423
const CLOSE = 57424
const DISPOSE = 57425
const GROUP_CONCAT = 57426
const SEPARATOR = 57427
const COMMIT = 57428
const ROLLBACK = 57429
const CONTINUE = 57430
const BREAK = 57431
const EXIT = 57432
const PRINT = 57433
const VAR = 57434
const COMPARISON_OP = 57435
const STRING_OP = 57436
const SUBSTITUTION_OP = 57437

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

//line parser.y:1212

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
	-1, 108,
	58, 204,
	62, 204,
	63, 204,
	-2, 216,
	-1, 115,
	44, 206,
	46, 210,
	-2, 139,
	-1, 158,
	104, 100,
	-2, 202,
	-1, 160,
	72, 130,
	-2, 204,
	-1, 165,
	37, 100,
	85, 100,
	104, 100,
	-2, 202,
	-1, 174,
	58, 204,
	62, 204,
	63, 204,
	-2, 124,
	-1, 179,
	58, 204,
	62, 204,
	63, 204,
	-2, 45,
	-1, 181,
	46, 210,
	-2, 206,
	-1, 207,
	58, 204,
	62, 204,
	63, 204,
	-2, 199,
	-1, 292,
	58, 204,
	62, 204,
	63, 204,
	-2, 135,
	-1, 309,
	58, 204,
	62, 204,
	63, 204,
	-2, 49,
	-1, 322,
	58, 204,
	62, 204,
	63, 204,
	-2, 157,
	-1, 345,
	76, 132,
	-2, 204,
	-1, 371,
	58, 204,
	62, 204,
	63, 204,
	-2, 166,
	-1, 380,
	72, 147,
	74, 147,
	76, 147,
	-2, 204,
	-1, 384,
	58, 204,
	62, 204,
	63, 204,
	-2, 75,
	-1, 387,
	88, 18,
	89, 18,
	-2, 1,
	-1, 391,
	58, 204,
	62, 204,
	63, 204,
	-2, 122,
}

const yyPrivate = 57344

const yyLast = 1012

var yyAct = [...]int{

	40, 406, 382, 292, 86, 42, 43, 44, 45, 46,
	47, 48, 416, 330, 267, 355, 81, 37, 291, 37,
	82, 23, 2, 23, 63, 64, 65, 115, 325, 261,
	258, 277, 194, 235, 106, 71, 103, 39, 1, 108,
	253, 297, 171, 90, 112, 3, 88, 331, 182, 105,
	338, 104, 49, 230, 114, 73, 180, 294, 144, 60,
	116, 132, 52, 41, 111, 402, 370, 324, 136, 137,
	138, 143, 55, 185, 168, 186, 187, 188, 183, 53,
	53, 181, 57, 55, 95, 96, 98, 155, 99, 100,
	38, 119, 121, 348, 319, 316, 160, 168, 162, 264,
	245, 133, 163, 125, 145, 146, 147, 148, 149, 170,
	36, 128, 347, 174, 420, 176, 401, 179, 191, 378,
	372, 369, 362, 323, 318, 289, 228, 55, 260, 156,
	55, 157, 158, 184, 215, 214, 93, 392, 207, 286,
	94, 122, 202, 165, 101, 212, 213, 92, 157, 220,
	221, 222, 223, 224, 225, 226, 210, 37, 96, 98,
	265, 23, 143, 102, 96, 72, 87, 169, 156, 193,
	111, 122, 134, 200, 208, 422, 84, 97, 211, 135,
	414, 388, 83, 37, 201, 210, 53, 23, 198, 237,
	375, 209, 344, 206, 336, 54, 54, 147, 148, 149,
	192, 299, 216, 66, 67, 68, 69, 70, 409, 280,
	276, 304, 408, 279, 233, 233, 295, 232, 410, 234,
	243, 285, 350, 288, 409, 244, 273, 54, 123, 270,
	252, 126, 280, 54, 425, 130, 131, 296, 251, 256,
	421, 302, 303, 301, 399, 305, 266, 263, 374, 174,
	97, 269, 309, 37, 275, 99, 390, 23, 76, 237,
	284, 110, 278, 283, 164, 143, 154, 218, 310, 322,
	300, 217, 219, 167, 254, 359, 335, 314, 311, 317,
	312, 255, 250, 315, 340, 354, 178, 349, 306, 308,
	196, 249, 248, 313, 189, 298, 37, 321, 231, 345,
	23, 197, 54, 333, 199, 342, 54, 332, 320, 120,
	337, 101, 334, 203, 204, 120, 185, 339, 186, 187,
	188, 37, 205, 59, 58, 23, 246, 237, 326, 327,
	328, 329, 156, 36, 229, 109, 371, 55, 352, 366,
	358, 363, 360, 139, 242, 377, 51, 361, 55, 365,
	380, 385, 161, 384, 113, 190, 124, 50, 262, 389,
	391, 54, 379, 55, 364, 129, 127, 257, 37, 197,
	272, 274, 23, 56, 387, 398, 394, 36, 107, 268,
	271, 197, 197, 38, 397, 386, 62, 393, 120, 211,
	381, 55, 142, 37, 395, 163, 403, 23, 396, 407,
	412, 61, 89, 85, 37, 413, 384, 411, 23, 415,
	10, 419, 9, 418, 400, 8, 154, 7, 6, 143,
	195, 5, 424, 307, 259, 37, 427, 54, 4, 23,
	196, 237, 54, 293, 159, 78, 172, 173, 118, 262,
	37, 197, 417, 120, 23, 117, 237, 271, 152, 141,
	197, 153, 145, 146, 147, 148, 149, 426, 91, 77,
	80, 74, 79, 75, 95, 96, 98, 383, 99, 100,
	353, 247, 177, 281, 282, 16, 15, 287, 151, 150,
	154, 238, 14, 143, 13, 262, 12, 11, 236, 54,
	0, 54, 0, 0, 367, 368, 197, 0, 0, 0,
	0, 268, 0, 0, 0, 197, 197, 0, 0, 120,
	0, 373, 152, 141, 120, 153, 145, 146, 147, 148,
	149, 151, 150, 154, 101, 41, 143, 0, 55, 95,
	96, 98, 0, 99, 100, 38, 54, 36, 0, 0,
	0, 0, 0, 55, 95, 96, 98, 271, 99, 100,
	38, 0, 0, 0, 0, 152, 141, 97, 153, 145,
	146, 147, 148, 149, 0, 0, 0, 343, 0, 268,
	0, 120, 185, 120, 186, 187, 188, 183, 356, 357,
	181, 93, 0, 0, 0, 94, 0, 0, 0, 101,
	0, 0, 92, 0, 0, 0, 93, 0, 0, 0,
	94, 0, 0, 0, 101, 0, 0, 92, 102, 55,
	95, 96, 98, 0, 99, 100, 38, 0, 120, 151,
	150, 154, 97, 102, 143, 0, 0, 83, 0, 55,
	95, 96, 98, 0, 99, 100, 38, 97, 175, 0,
	0, 185, 83, 186, 187, 188, 183, 0, 0, 181,
	0, 0, 0, 152, 141, 0, 153, 145, 146, 147,
	148, 149, 93, 0, 227, 0, 94, 0, 0, 0,
	101, 0, 0, 92, 0, 0, 0, 404, 405, 0,
	0, 0, 93, 0, 0, 0, 94, 0, 0, 102,
	101, 0, 0, 92, 0, 151, 150, 154, 0, 0,
	143, 0, 0, 97, 290, 0, 0, 0, 83, 102,
	38, 0, 36, 0, 18, 34, 19, 0, 17, 0,
	0, 0, 0, 97, 20, 0, 0, 21, 83, 152,
	141, 0, 153, 145, 146, 147, 148, 149, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 151,
	150, 154, 0, 0, 143, 0, 0, 0, 239, 0,
	32, 0, 0, 423, 0, 0, 26, 0, 0, 30,
	27, 28, 29, 0, 0, 24, 25, 240, 241, 33,
	35, 22, 0, 152, 141, 0, 153, 145, 146, 147,
	148, 149, 151, 150, 154, 0, 0, 143, 0, 0,
	0, 0, 151, 150, 154, 0, 376, 143, 0, 0,
	0, 0, 0, 0, 0, 0, 351, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 152, 141, 0, 153,
	145, 146, 147, 148, 149, 0, 152, 141, 0, 153,
	145, 146, 147, 148, 149, 151, 150, 154, 0, 0,
	143, 0, 0, 0, 0, 151, 150, 154, 0, 346,
	143, 0, 0, 0, 151, 150, 154, 0, 0, 143,
	0, 166, 0, 0, 0, 0, 0, 0, 140, 152,
	141, 0, 153, 145, 146, 147, 148, 149, 0, 152,
	141, 0, 153, 145, 146, 147, 148, 149, 152, 141,
	0, 153, 145, 146, 147, 148, 149, 151, 150, 154,
	0, 0, 143, 0, 0, 0, 341, 150, 154, 0,
	38, 143, 36, 0, 18, 34, 19, 151, 17, 154,
	0, 0, 143, 0, 20, 0, 0, 21, 0, 0,
	0, 152, 141, 0, 153, 145, 146, 147, 148, 149,
	152, 141, 0, 153, 145, 146, 147, 148, 149, 0,
	0, 152, 141, 0, 153, 145, 146, 147, 148, 149,
	0, 0, 0, 0, 0, 0, 0, 0, 31, 0,
	32, 0, 0, 0, 0, 0, 26, 0, 0, 30,
	27, 28, 29, 0, 0, 24, 25, 0, 0, 33,
	35, 22,
}
var yyPact = [...]int{

	919, -1000, 919, -43, -43, -43, -43, -43, -43, -43,
	-43, -1000, -1000, -1000, -1000, -1000, 343, 326, 387, 359,
	295, 294, 375, -43, -43, -43, 387, 387, 387, 387,
	387, 625, 625, -43, 366, 625, 195, 75, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 336,
	68, 387, 340, -2, 344, -1000, 68, 351, 387, 387,
	-43, -4, 77, -1000, -1000, -1000, 101, -43, -43, -43,
	323, 815, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 75, -1000, 524, 29, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 625, 38, 625, -1000, -1000, 152, -1000, -1000,
	-1000, -1000, 40, 806, 215, -31, -1000, 71, 419, 539,
	-1000, 625, 248, 625, -1000, 28, -1000, 333, -1000, -1000,
	-1000, -1000, 364, 97, 387, 387, -1000, 387, 336, 68,
	39, 287, -1000, 375, 625, 95, -1000, -1000, -1000, 372,
	919, 625, 79, 205, 209, 625, 625, 625, 625, 625,
	625, 625, -1000, -1000, -1000, 560, 22, 387, 195, 143,
	858, -1000, 355, -1000, -1000, 195, 699, 387, 372, 459,
	-1000, -1000, -5, 304, 858, -1000, 858, 253, 251, 858,
	238, 271, 228, 237, 68, -1000, -1000, -1000, -1000, -1000,
	387, 25, 387, -1000, 343, -6, 64, 46, -1000, -1000,
	-1000, 336, 387, 126, 123, 387, -1000, 858, 364, -43,
	-8, 139, 7, 858, 38, 38, 246, 625, 36, 625,
	98, 98, 201, 201, 201, 878, 355, -1000, -1000, -1000,
	21, 605, 142, 625, 258, 125, 699, -1000, -1000, 625,
	-43, -43, 136, -1000, -43, 539, 387, 258, 625, 625,
	68, 236, 228, 233, -1000, 68, -1000, -1000, -1000, -10,
	625, 20, -11, 336, 387, 625, -1000, 19, -38, 298,
	387, 273, -1000, 387, 276, -43, -1000, 118, 139, 919,
	625, -1000, -1000, -1000, -1000, 867, 524, -1000, 858, -1000,
	-1000, -1000, 462, 116, 143, 625, 796, 8, 247, 151,
	-1000, 753, -1000, -1000, 699, -1000, -1000, -1000, 242, 858,
	-1000, 527, 68, 231, 68, 596, 25, 18, 320, 387,
	-1000, -1000, 858, -1000, 387, -1000, -1000, -1000, 387, 387,
	17, -39, 625, 16, 387, -1000, 179, 114, 162, -1000,
	743, 625, 15, 625, -1000, 858, 625, -1000, 385, 625,
	-43, 699, 105, -1000, 158, -1000, 625, 34, 596, 68,
	527, -1000, -1000, 25, -1000, -1000, -1000, -1000, -1000, 298,
	387, 858, -1000, -1000, -43, 175, 919, 355, -1000, -1000,
	858, 12, -1000, -40, 636, -1000, 138, 919, 147, -1000,
	158, 858, 387, 596, -1000, -1000, -1000, -1000, -1000, -43,
	-1000, -1000, 625, -1000, -1000, -1000, 104, 138, 699, 625,
	-43, 10, -1000, -1000, 171, 99, 154, -1000, 700, -1000,
	-1000, -43, 165, 699, -1000, -43, -1000, -1000,
}
var yyPgo = [...]int{

	0, 37, 33, 22, 488, 487, 486, 484, 482, 481,
	476, 45, 475, 52, 44, 472, 471, 41, 470, 55,
	165, 3, 467, 258, 463, 462, 461, 460, 459, 53,
	458, 60, 445, 27, 438, 15, 437, 436, 435, 434,
	433, 29, 18, 2, 54, 62, 14, 42, 57, 428,
	424, 30, 421, 420, 32, 418, 417, 415, 47, 13,
	28, 412, 410, 50, 31, 12, 1, 176, 403, 4,
	166, 46, 402, 43, 16, 51, 20, 401, 59, 298,
	58, 396, 56, 40, 48, 392, 0,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10,
	11, 12, 13, 13, 14, 14, 15, 15, 16, 16,
	17, 17, 18, 18, 19, 19, 19, 19, 19, 19,
	20, 20, 21, 21, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 22, 81, 81, 81, 23, 24,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	26, 26, 26, 26, 26, 27, 27, 27, 28, 28,
	29, 29, 29, 30, 30, 31, 31, 31, 32, 32,
	33, 33, 33, 33, 33, 33, 34, 34, 34, 34,
	34, 35, 35, 35, 36, 36, 37, 37, 38, 39,
	39, 40, 40, 41, 41, 42, 42, 43, 43, 44,
	44, 45, 45, 46, 46, 47, 47, 48, 48, 49,
	49, 49, 49, 50, 51, 51, 52, 53, 54, 54,
	55, 55, 56, 57, 57, 58, 58, 59, 59, 60,
	60, 60, 60, 60, 61, 61, 62, 63, 63, 64,
	64, 65, 65, 66, 66, 67, 68, 69, 69, 70,
	70, 71, 72, 73, 74, 75, 75, 76, 77, 77,
	78, 78, 79, 79, 80, 80, 82, 82, 83, 83,
	84, 84, 84, 84, 85, 85, 86, 86,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 5, 8,
	9, 7, 9, 2, 8, 9, 2, 2, 5, 3,
	7, 3, 0, 2, 0, 2, 0, 3, 0, 2,
	0, 3, 0, 2, 1, 1, 1, 1, 1, 1,
	1, 3, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 2, 0, 1, 1, 3, 3,
	3, 4, 4, 6, 6, 4, 4, 4, 4, 2,
	3, 3, 3, 3, 3, 3, 3, 2, 4, 1,
	0, 2, 2, 5, 7, 1, 2, 3, 1, 1,
	1, 1, 2, 3, 1, 1, 5, 5, 6, 6,
	4, 0, 2, 4, 1, 1, 1, 3, 5, 0,
	1, 0, 2, 1, 3, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 1, 3, 4, 2, 5,
	8, 4, 7, 3, 1, 3, 6, 3, 1, 3,
	4, 5, 6, 6, 8, 1, 3, 1, 3, 0,
	1, 1, 2, 2, 5, 7, 7, 4, 2, 0,
	2, 4, 2, 0, 2, 1, 1, 1, 2, 1,
	2, 1, 1, 1, 1, 1, 3, 3, 1, 3,
	1, 3, 0, 1, 0, 1, 0, 1, 0, 1,
	0, 1, 1, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -49, -52, -55, -56, -57, -61,
	-62, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 92, -76, 86, 87, 77, 81, 82, 83,
	80, 69, 71, 90, 16, 91, 13, -74, 11, -1,
	-86, 106, -86, -86, -86, -86, -86, -86, -86, -13,
	14, 20, -45, -31, -67, 4, 14, -45, 29, 29,
	-78, -77, 11, -86, -86, -86, -67, -67, -67, -67,
	-67, -21, -20, -19, -26, -24, -23, -28, -38, -25,
	-27, -74, -76, 103, -67, -68, -69, -70, -71, -72,
	-73, -30, 68, 57, 61, 5, 6, 98, 7, 9,
	10, 65, 84, -21, -75, -74, -86, 12, -21, -79,
	66, 95, -14, 18, -44, -33, -31, -32, -34, 23,
	-23, 24, 103, -67, 16, 105, -67, 22, -44, 14,
	-67, -67, -86, 105, 95, 78, -86, -86, -86, 20,
	73, 94, -85, 64, -80, 97, 98, 99, 100, 101,
	60, 59, 93, 96, 61, -21, -11, 102, 103, -39,
	-21, -23, -21, -69, -70, 103, 75, 58, 105, 96,
	-86, -47, -37, -36, -21, 99, -21, -15, 38, -21,
	-82, 53, -84, 50, 105, 45, 47, 48, 49, -67,
	22, 21, 103, -11, -54, -53, -20, -67, -45, -67,
	-14, -44, 103, 26, 27, 35, -78, -21, 79, -75,
	-74, -1, -21, -21, 56, 55, -80, 62, 58, 63,
	-21, -21, -21, -21, -21, -21, -21, 104, 104, -67,
	-29, -79, -48, 72, -29, -2, -4, -3, -9, 69,
	88, 89, -67, -75, -19, 105, 22, -16, 39, 40,
	44, -82, -84, -83, 46, 44, -44, -67, -51, -50,
	103, -41, -20, -13, 105, 96, -14, -46, -67, -58,
	103, -67, -20, 103, -20, -11, -86, -64, -63, 74,
	70, -23, -23, -71, -73, -21, 103, -23, -21, 104,
	99, -42, -21, -40, -48, 74, -21, -17, 37, 76,
	-2, -21, -86, -86, 75, -86, -47, -67, -17, -21,
	-42, -33, 44, -83, 44, -33, 105, -42, 104, 105,
	-14, -54, -21, 104, 105, -60, 30, 31, 32, 33,
	-59, -58, 34, -41, 36, -86, 76, -64, -63, -1,
	-21, 59, -42, 105, 76, -21, 73, 104, 85, 40,
	71, 73, -2, -18, 43, -35, 51, 52, -33, 44,
	-33, -51, 104, 21, -11, -41, -46, -20, -20, 104,
	105, -21, 104, -67, 69, 76, 73, -21, 104, -42,
	-21, 5, -43, -22, -21, -86, -2, -3, 76, -69,
	98, -21, 103, -33, -35, -51, -60, -59, -86, 69,
	-1, 104, 105, -81, 41, 42, -66, -65, 74, 70,
	71, -46, -86, -43, 76, -66, -65, -2, -21, -86,
	104, 69, 76, 73, -86, 69, -2, -86,
}
var yyDef = [...]int{

	1, -2, 1, 216, 216, 216, 216, 216, 216, 216,
	216, 13, 14, 15, 16, 17, 42, 0, 0, 0,
	0, 0, 0, 216, 216, 216, 0, 0, 0, 0,
	0, 0, 0, 216, 0, 0, 202, 0, 194, 2,
	5, 217, 6, 7, 8, 9, 10, 11, 12, 44,
	0, 0, 0, 141, 105, 185, 0, 0, 0, 0,
	216, 200, 198, 21, 22, 23, 0, 216, 216, 216,
	0, 204, 62, 63, 64, 65, 66, 67, 68, 69,
	70, 71, 72, 0, 60, 54, 55, 56, 57, 58,
	59, 99, 129, 0, 0, 186, 187, 0, 189, 191,
	192, 193, 0, 204, 0, 71, 33, 0, -2, 0,
	203, 0, 46, 0, 43, -2, 110, 111, 114, 115,
	108, 109, 0, 0, 0, 0, 106, 0, 44, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	1, 0, 0, 204, 0, 0, 0, 0, 0, 0,
	0, 0, 214, 215, 205, 204, 0, 0, -2, 0,
	-2, 89, 97, 188, 190, -2, 3, 0, 0, 0,
	39, 41, 145, 126, -2, 125, 197, 48, 0, -2,
	0, -2, 208, 0, 0, 207, 211, 212, 213, 112,
	0, 0, 0, 151, 42, 158, 0, 60, 142, 107,
	160, 44, 0, 0, 0, 0, 201, -2, 0, 216,
	195, 179, 79, 80, 0, 0, 0, 0, 0, 0,
	90, 91, 92, 93, 94, 95, 96, 73, 78, 61,
	0, 0, 131, 0, 50, 0, 3, 18, 19, 0,
	216, 216, 0, 196, 216, 0, 0, 50, 0, 0,
	0, 0, 208, 0, 209, 0, 140, 113, 149, 154,
	0, 0, 133, 44, 0, 0, 161, 0, 143, 169,
	0, 165, 174, 0, 0, 216, 28, 0, 179, 1,
	0, 87, 88, 81, 82, 204, 0, 85, 86, 98,
	101, 102, -2, 0, 148, 0, 204, 0, 0, 0,
	4, 204, 36, 37, 3, 38, 146, 127, 52, -2,
	47, 121, 0, 0, 0, 120, 0, 0, 0, 0,
	156, 159, -2, 162, 0, 163, 170, 171, 0, 0,
	0, 167, 0, 0, 0, 24, 0, 0, 178, 180,
	204, 0, 0, 0, 128, -2, 0, 103, 0, 0,
	216, 1, 0, 40, 0, 116, 0, 0, 117, 0,
	121, 155, 153, 0, 152, 134, 144, 172, 173, 169,
	0, -2, 175, 176, 216, 0, 1, 83, 84, 136,
	-2, 0, 51, 137, -2, 31, 183, -2, 0, 53,
	0, -2, 0, 119, 118, 150, 164, 168, 29, 216,
	177, 104, 0, 74, 76, 77, 0, 183, 3, 0,
	216, 0, 30, 138, 0, 0, 182, 184, 204, 32,
	123, 216, 0, 3, 34, 216, 181, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 101, 3, 3,
	103, 104, 99, 97, 105, 98, 102, 100, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 106,
	3, 96,
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
	92, 93, 94, 95,
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
		//line parser.y:141
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:146
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:153
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:158
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:165
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:169
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:173
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:177
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:181
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:185
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:189
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:193
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:197
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:201
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:205
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:209
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:213
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:219
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:223
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:229
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:233
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:239
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:243
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:249
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:253
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:257
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:261
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:265
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:271
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:275
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:279
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:283
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:287
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:293
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:297
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:301
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:305
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:311
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:315
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:321
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
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:335
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 42:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = nil
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 44:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:351
		{
			yyVAL.expression = nil
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:355
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 46:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:361
		{
			yyVAL.expression = nil
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:365
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:371
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:375
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:381
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:385
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:391
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:401
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:405
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:409
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:413
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:417
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:421
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:427
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:431
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:437
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:441
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:445
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:449
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:453
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:457
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:461
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:465
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:469
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:473
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:477
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:481
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:487
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 75:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:493
		{
			yyVAL.token = Token{}
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:497
		{
			yyVAL.token = yyDollar[1].token
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:501
		{
			yyVAL.token = yyDollar[1].token
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:507
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:513
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
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 81:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 82:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:544
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 83:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:548
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 84:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 85:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 86:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:560
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 87:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:564
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 88:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:568
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:578
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:582
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:586
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:590
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:594
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:600
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:604
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:608
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 98:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:614
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:618
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 100:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:624
		{
			yyVAL.expression = Option{}
		}
	case 101:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:628
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:632
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 103:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:638
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 104:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:642
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:648
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:652
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:656
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:662
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:666
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:672
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:676
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:680
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:684
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:688
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:692
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 116:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:698
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 117:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:702
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 118:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:706
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 119:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:710
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 120:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:714
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 121:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:720
		{
			yyVAL.expression = nil
		}
	case 122:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:724
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 123:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:728
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:734
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 125:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:738
		{
			yyVAL.expression = AllColumns{}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:744
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:748
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 128:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:754
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 129:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:760
		{
			yyVAL.expression = nil
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:764
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 131:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:770
		{
			yyVAL.expression = nil
		}
	case 132:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:774
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:780
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 134:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:784
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:790
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:794
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:800
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 138:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:804
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 139:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:810
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 140:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:814
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 141:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:820
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:830
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 144:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:834
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 145:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:840
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 146:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:844
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 147:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:850
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 148:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:854
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 149:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:860
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 150:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:864
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 151:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:868
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 152:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:872
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 153:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:878
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:884
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 155:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:888
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 156:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:894
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 157:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:900
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:906
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 159:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:910
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 160:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:916
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 161:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:921
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 162:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:928
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 163:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:934
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 164:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:938
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:944
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 166:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:948
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:954
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:958
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 169:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:964
		{
			yyVAL.expression = nil
		}
	case 170:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:968
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:972
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 172:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:976
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 173:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:980
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 174:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:986
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 175:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:990
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 176:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:996
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 177:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1002
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 178:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1006
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 179:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1012
		{
			yyVAL.procexpr = nil
		}
	case 180:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1016
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 181:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1022
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 182:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1026
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 183:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1032
		{
			yyVAL.procexpr = nil
		}
	case 184:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1036
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1042
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 186:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1048
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 187:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1054
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 188:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1058
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1065
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 190:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1069
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1076
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 192:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1082
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 193:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1088
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1094
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1100
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 196:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1104
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 197:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1110
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 198:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1116
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 199:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1120
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1126
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 201:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1130
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 202:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1136
		{
			yyVAL.token = Token{}
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1140
		{
			yyVAL.token = yyDollar[1].token
		}
	case 204:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1146
		{
			yyVAL.token = Token{}
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1150
		{
			yyVAL.token = yyDollar[1].token
		}
	case 206:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1156
		{
			yyVAL.token = Token{}
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1160
		{
			yyVAL.token = yyDollar[1].token
		}
	case 208:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1166
		{
			yyVAL.token = Token{}
		}
	case 209:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1170
		{
			yyVAL.token = yyDollar[1].token
		}
	case 210:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1176
		{
			yyVAL.token = Token{}
		}
	case 211:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1180
		{
			yyVAL.token = yyDollar[1].token
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1184
		{
			yyVAL.token = yyDollar[1].token
		}
	case 213:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1188
		{
			yyVAL.token = yyDollar[1].token
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1194
		{
			yyVAL.token = yyDollar[1].token
		}
	case 215:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1198
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 216:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1204
		{
			yyVAL.token = Token{}
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1208
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
