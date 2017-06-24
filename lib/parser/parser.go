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
const OFFSET = 57386
const JOIN = 57387
const INNER = 57388
const OUTER = 57389
const LEFT = 57390
const RIGHT = 57391
const FULL = 57392
const CROSS = 57393
const ON = 57394
const USING = 57395
const NATURAL = 57396
const UNION = 57397
const INTERSECT = 57398
const EXCEPT = 57399
const ALL = 57400
const ANY = 57401
const EXISTS = 57402
const IN = 57403
const AND = 57404
const OR = 57405
const NOT = 57406
const BETWEEN = 57407
const LIKE = 57408
const IS = 57409
const NULL = 57410
const DISTINCT = 57411
const WITH = 57412
const CASE = 57413
const IF = 57414
const ELSEIF = 57415
const WHILE = 57416
const WHEN = 57417
const THEN = 57418
const ELSE = 57419
const DO = 57420
const END = 57421
const DECLARE = 57422
const CURSOR = 57423
const FOR = 57424
const FETCH = 57425
const OPEN = 57426
const CLOSE = 57427
const DISPOSE = 57428
const GROUP_CONCAT = 57429
const SEPARATOR = 57430
const COMMIT = 57431
const ROLLBACK = 57432
const CONTINUE = 57433
const BREAK = 57434
const EXIT = 57435
const PRINT = 57436
const VAR = 57437
const COMPARISON_OP = 57438
const STRING_OP = 57439
const SUBSTITUTION_OP = 57440

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
	"OFFSET",
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

//line parser.y:1289

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
	55, 45,
	56, 45,
	57, 45,
	-2, 56,
	-1, 112,
	61, 213,
	65, 213,
	66, 213,
	-2, 227,
	-1, 131,
	45, 215,
	47, 219,
	-2, 148,
	-1, 165,
	55, 46,
	56, 46,
	57, 46,
	-2, 74,
	-1, 167,
	107, 109,
	-2, 211,
	-1, 169,
	75, 139,
	-2, 213,
	-1, 174,
	37, 109,
	88, 109,
	107, 109,
	-2, 211,
	-1, 191,
	61, 213,
	65, 213,
	66, 213,
	-2, 133,
	-1, 199,
	61, 213,
	65, 213,
	66, 213,
	-2, 83,
	-1, 211,
	47, 219,
	-2, 215,
	-1, 227,
	61, 213,
	65, 213,
	66, 213,
	-2, 208,
	-1, 233,
	67, 0,
	96, 0,
	99, 0,
	-2, 88,
	-1, 234,
	67, 0,
	96, 0,
	99, 0,
	-2, 89,
	-1, 267,
	61, 213,
	65, 213,
	66, 213,
	-2, 51,
	-1, 315,
	67, 0,
	96, 0,
	99, 0,
	-2, 95,
	-1, 321,
	61, 213,
	65, 213,
	66, 213,
	-2, 144,
	-1, 346,
	61, 213,
	65, 213,
	66, 213,
	-2, 166,
	-1, 374,
	79, 141,
	-2, 213,
	-1, 381,
	61, 213,
	65, 213,
	66, 213,
	-2, 55,
	-1, 399,
	61, 213,
	65, 213,
	66, 213,
	-2, 175,
	-1, 408,
	75, 156,
	77, 156,
	79, 156,
	-2, 213,
	-1, 412,
	91, 18,
	92, 18,
	-2, 1,
	-1, 415,
	61, 213,
	65, 213,
	66, 213,
	-2, 131,
}

const yyPrivate = 57344

const yyLast = 1053

var yyAct = [...]int{

	80, 40, 296, 40, 434, 359, 44, 425, 354, 320,
	282, 46, 47, 48, 49, 50, 51, 52, 388, 85,
	38, 306, 38, 188, 279, 180, 290, 131, 43, 1,
	67, 68, 69, 203, 53, 113, 197, 94, 92, 212,
	110, 360, 2, 40, 210, 77, 367, 37, 323, 154,
	108, 255, 109, 64, 132, 130, 183, 250, 45, 118,
	398, 136, 16, 56, 353, 377, 343, 340, 177, 177,
	285, 141, 275, 57, 57, 3, 272, 59, 145, 146,
	147, 142, 127, 61, 376, 438, 39, 39, 165, 215,
	424, 216, 217, 218, 213, 385, 406, 211, 170, 400,
	397, 321, 384, 352, 59, 86, 23, 162, 23, 163,
	342, 318, 153, 193, 59, 136, 416, 59, 121, 179,
	162, 161, 163, 135, 137, 153, 40, 156, 157, 158,
	159, 160, 39, 75, 107, 166, 167, 112, 281, 136,
	200, 151, 150, 42, 152, 156, 157, 158, 159, 160,
	40, 214, 313, 222, 151, 150, 209, 152, 156, 157,
	158, 159, 160, 121, 158, 159, 160, 45, 230, 38,
	182, 174, 163, 185, 186, 153, 40, 166, 231, 42,
	42, 286, 57, 178, 40, 118, 40, 40, 91, 164,
	90, 207, 39, 143, 221, 38, 226, 230, 169, 229,
	171, 202, 40, 235, 151, 150, 42, 152, 156, 157,
	158, 159, 160, 228, 144, 136, 302, 252, 257, 299,
	187, 191, 440, 76, 264, 201, 199, 432, 263, 40,
	413, 268, 254, 270, 271, 403, 305, 373, 314, 284,
	316, 317, 365, 327, 332, 227, 269, 295, 269, 269,
	253, 289, 232, 233, 234, 23, 288, 40, 241, 242,
	243, 244, 245, 246, 247, 298, 429, 330, 331, 378,
	293, 333, 428, 311, 310, 428, 38, 121, 307, 427,
	309, 23, 309, 267, 308, 42, 100, 102, 136, 326,
	173, 341, 172, 136, 443, 253, 337, 324, 439, 257,
	422, 99, 100, 102, 304, 103, 104, 402, 328, 40,
	344, 364, 339, 362, 165, 347, 349, 120, 163, 345,
	351, 237, 176, 371, 103, 236, 238, 184, 38, 366,
	240, 239, 116, 40, 115, 116, 117, 368, 312, 88,
	315, 291, 392, 40, 350, 348, 382, 292, 287, 136,
	205, 136, 38, 195, 387, 325, 394, 123, 58, 58,
	336, 329, 23, 124, 105, 383, 70, 71, 72, 73,
	74, 335, 266, 54, 191, 257, 391, 199, 393, 363,
	40, 101, 407, 105, 380, 410, 361, 251, 346, 121,
	63, 62, 273, 136, 148, 125, 55, 101, 128, 38,
	58, 181, 139, 140, 420, 40, 419, 126, 231, 421,
	414, 369, 418, 40, 23, 114, 426, 59, 386, 430,
	417, 59, 412, 138, 38, 283, 374, 119, 40, 431,
	59, 411, 38, 423, 433, 220, 437, 381, 23, 129,
	60, 215, 40, 216, 217, 218, 442, 38, 301, 303,
	445, 111, 150, 41, 58, 156, 157, 158, 159, 160,
	66, 38, 409, 399, 274, 196, 206, 58, 59, 208,
	257, 155, 405, 219, 223, 224, 276, 408, 58, 435,
	65, 93, 89, 225, 257, 23, 355, 356, 357, 358,
	10, 415, 215, 444, 216, 217, 218, 213, 389, 390,
	211, 9, 8, 7, 6, 204, 249, 5, 280, 205,
	23, 4, 322, 168, 82, 41, 262, 39, 23, 18,
	34, 19, 189, 17, 190, 134, 283, 133, 95, 20,
	436, 81, 21, 23, 215, 84, 216, 217, 218, 213,
	78, 206, 211, 83, 79, 198, 194, 23, 122, 334,
	265, 36, 15, 258, 58, 14, 13, 12, 11, 256,
	294, 0, 297, 300, 206, 206, 0, 283, 0, 0,
	0, 0, 0, 0, 0, 0, 259, 0, 32, 0,
	0, 395, 396, 0, 26, 0, 0, 30, 27, 28,
	29, 0, 0, 24, 25, 260, 261, 33, 35, 22,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	42, 0, 0, 338, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 206, 0, 58, 162, 161,
	163, 0, 58, 153, 0, 0, 0, 0, 0, 300,
	0, 41, 206, 39, 0, 18, 34, 19, 0, 17,
	0, 0, 0, 0, 0, 20, 0, 0, 21, 0,
	0, 0, 151, 150, 0, 152, 156, 157, 158, 159,
	160, 0, 0, 0, 372, 59, 99, 100, 102, 0,
	103, 104, 41, 206, 39, 0, 0, 0, 58, 0,
	58, 0, 0, 297, 0, 0, 0, 206, 206, 0,
	0, 0, 31, 401, 32, 0, 0, 0, 0, 0,
	26, 0, 0, 30, 27, 28, 29, 0, 0, 24,
	25, 0, 0, 33, 35, 22, 0, 0, 0, 0,
	0, 97, 58, 0, 0, 98, 42, 0, 300, 105,
	0, 0, 96, 59, 99, 100, 102, 0, 103, 104,
	41, 0, 0, 0, 0, 0, 297, 0, 106, 0,
	59, 99, 100, 102, 0, 103, 104, 41, 0, 0,
	0, 0, 101, 0, 0, 0, 0, 87, 0, 0,
	0, 0, 59, 99, 100, 102, 0, 103, 104, 41,
	0, 0, 277, 278, 0, 0, 0, 0, 0, 97,
	0, 0, 0, 98, 0, 0, 0, 105, 0, 0,
	96, 0, 0, 162, 161, 163, 97, 0, 153, 0,
	98, 0, 0, 0, 105, 0, 106, 96, 0, 162,
	161, 163, 0, 0, 153, 0, 0, 0, 97, 0,
	101, 192, 98, 106, 0, 87, 105, 151, 150, 96,
	152, 156, 157, 158, 159, 160, 0, 101, 319, 0,
	0, 0, 87, 151, 150, 106, 152, 156, 157, 158,
	159, 160, 0, 0, 248, 162, 161, 163, 0, 101,
	153, 0, 0, 0, 87, 162, 161, 163, 0, 441,
	153, 0, 0, 0, 0, 0, 0, 0, 0, 404,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 151,
	150, 0, 152, 156, 157, 158, 159, 160, 0, 151,
	150, 0, 152, 156, 157, 158, 159, 160, 162, 161,
	163, 0, 0, 153, 0, 0, 0, 0, 162, 161,
	163, 0, 379, 153, 0, 0, 0, 0, 162, 161,
	163, 0, 375, 153, 0, 0, 0, 0, 0, 0,
	0, 0, 151, 150, 175, 152, 156, 157, 158, 159,
	160, 0, 151, 150, 0, 152, 156, 157, 158, 159,
	160, 0, 151, 150, 0, 152, 156, 157, 158, 159,
	160, 162, 161, 163, 0, 0, 153, 0, 0, 0,
	0, 162, 161, 163, 0, 149, 153, 0, 0, 0,
	370, 161, 163, 0, 0, 153, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 151, 150, 0, 152, 156,
	157, 158, 159, 160, 0, 151, 150, 0, 152, 156,
	157, 158, 159, 160, 151, 150, 0, 152, 156, 157,
	158, 159, 160,
}
var yyPact = [...]int{

	630, -1000, 630, -51, -51, -51, -51, -51, -51, -51,
	-51, -1000, -1000, -1000, -1000, -1000, 336, 376, 464, 426,
	362, 361, 449, -51, -51, -51, 464, 464, 464, 464,
	464, 778, 778, -51, 439, 778, 401, 279, 87, 248,
	-1000, -1000, 179, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 314, 323, 464, 391, -26, 417, -1000,
	100, 409, 464, 464, -51, -27, 95, -1000, -1000, -1000,
	133, -51, -51, -51, 374, 929, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 87, -1000, 671, 30, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 778, 37, 778, -1000,
	-1000, 280, -1000, -1000, -1000, -1000, 65, 886, 261, -39,
	-1000, 84, 58, 383, 100, 269, 269, 269, 778, 739,
	-1000, 6, 309, 459, 778, 119, 464, 464, -1000, 464,
	383, 43, -1000, 413, -1000, -1000, -1000, -1000, 100, 47,
	448, -1000, 449, 778, 131, -1000, -1000, -1000, 442, 630,
	778, 778, 778, 254, 260, 272, 778, 778, 778, 778,
	778, 778, 778, -1000, 767, -1000, 464, 248, 175, 939,
	-1000, 108, -1000, -1000, 248, 504, 464, 442, 296, -1000,
	334, 778, -1000, 179, -1000, 179, 179, 939, -1000, -32,
	370, 939, -1000, -1000, -1000, 458, -1000, -1000, -36, 751,
	32, 73, -1000, 401, -38, 82, 72, -1000, -1000, -1000,
	303, 395, 294, 302, 100, -1000, -1000, -1000, -1000, -1000,
	464, 383, 464, 113, 110, 464, -1000, 939, 179, -51,
	-40, 207, 27, 355, 355, 315, 778, 46, 778, 37,
	37, 62, 62, -1000, -1000, -1000, 45, 108, -1000, -1000,
	4, 756, 220, 778, 336, 164, 504, -1000, -1000, 778,
	-51, -51, 166, -1000, -51, 332, 320, 939, 276, -1000,
	-1000, 276, 739, 464, -1000, 778, -1000, -1000, -1000, -1000,
	-41, 778, 3, -42, 383, 464, 778, 100, 300, 294,
	299, -1000, 100, -1000, -1000, -1000, -4, -44, 456, 464,
	352, -1000, 464, 343, -51, -1000, 163, 207, 630, 778,
	-1000, -1000, 948, 671, -1000, 355, -1000, -1000, -1000, -1000,
	-1000, 566, 158, 175, 778, 876, -23, 195, -1000, 866,
	-1000, -1000, 504, -1000, -1000, 778, 778, -1000, -1000, -1000,
	32, -5, 74, 464, -1000, -1000, 939, 446, 100, 297,
	100, 488, -1000, 464, -1000, -1000, -1000, 464, 464, -7,
	-48, 778, -8, 464, -1000, 235, 156, 209, -1000, 823,
	778, -11, 778, -1000, 939, 778, -1000, 457, -51, 504,
	151, 939, -1000, -1000, -1000, 32, -1000, -1000, -1000, 778,
	10, 488, 100, 446, -1000, -1000, -1000, 456, 464, 939,
	-1000, -1000, -51, 228, 630, 108, -1000, -1000, 939, -17,
	-1000, 202, 630, 192, -1000, 939, 464, 488, -1000, -1000,
	-1000, -1000, -51, -1000, -1000, 148, 202, 504, 778, -51,
	-22, -1000, 226, 143, 199, -1000, 813, -1000, -1000, -51,
	222, 504, -1000, -51, -1000, -1000,
}
var yyPgo = [...]int{

	0, 28, 51, 42, 559, 558, 557, 556, 555, 553,
	552, 75, 62, 47, 551, 35, 25, 550, 549, 34,
	548, 546, 45, 223, 101, 545, 0, 544, 543, 540,
	535, 531, 57, 528, 54, 527, 27, 525, 18, 524,
	522, 514, 513, 512, 10, 9, 36, 55, 63, 2,
	23, 48, 511, 508, 24, 507, 505, 33, 504, 503,
	502, 41, 5, 8, 501, 490, 46, 21, 4, 7,
	339, 482, 190, 188, 38, 481, 37, 19, 50, 105,
	480, 53, 387, 49, 476, 44, 26, 39, 56, 471,
	6,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10,
	11, 12, 12, 12, 12, 13, 13, 14, 15, 15,
	16, 16, 17, 17, 18, 18, 19, 19, 20, 20,
	21, 21, 22, 22, 22, 22, 22, 22, 23, 23,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 25, 84, 84, 84, 26, 27, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 29,
	29, 29, 29, 29, 30, 30, 30, 31, 31, 32,
	32, 32, 33, 33, 34, 34, 34, 35, 35, 36,
	36, 36, 36, 36, 36, 37, 37, 37, 37, 37,
	38, 38, 38, 39, 39, 40, 40, 41, 42, 42,
	43, 43, 44, 44, 45, 45, 46, 46, 47, 47,
	48, 48, 49, 49, 50, 50, 51, 51, 52, 52,
	52, 52, 53, 54, 54, 55, 56, 57, 57, 58,
	58, 59, 60, 60, 61, 61, 62, 62, 63, 63,
	63, 63, 63, 64, 64, 65, 66, 66, 67, 67,
	68, 68, 69, 69, 70, 71, 72, 72, 73, 73,
	74, 75, 76, 77, 78, 78, 79, 80, 80, 81,
	81, 82, 82, 83, 83, 85, 85, 86, 86, 87,
	87, 87, 87, 88, 88, 89, 89, 90, 90,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 5, 8,
	9, 7, 9, 2, 8, 9, 2, 2, 5, 3,
	4, 5, 4, 4, 4, 1, 1, 3, 0, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 2,
	0, 2, 1, 1, 1, 1, 1, 1, 1, 3,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 2, 0, 1, 1, 3, 3, 3, 3,
	4, 4, 6, 6, 4, 4, 4, 4, 2, 3,
	3, 3, 3, 3, 3, 3, 2, 4, 1, 0,
	2, 2, 5, 7, 1, 2, 3, 1, 1, 1,
	1, 2, 3, 1, 1, 5, 5, 6, 6, 4,
	0, 2, 4, 1, 1, 1, 3, 5, 0, 1,
	0, 2, 1, 3, 1, 3, 1, 3, 1, 3,
	1, 3, 1, 3, 1, 3, 4, 2, 5, 8,
	4, 7, 3, 1, 3, 6, 3, 1, 3, 4,
	5, 6, 6, 8, 1, 3, 1, 3, 0, 1,
	1, 2, 2, 5, 7, 7, 4, 2, 0, 2,
	4, 2, 0, 2, 1, 1, 1, 2, 1, 2,
	1, 1, 1, 1, 1, 3, 3, 1, 3, 1,
	3, 0, 1, 0, 1, 0, 1, 0, 1, 0,
	1, 1, 1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -52, -55, -58, -59, -60, -64,
	-65, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 95, -79, 89, 90, 80, 84, 85, 86,
	83, 72, 74, 93, 16, 94, -14, -13, -77, 13,
	-26, 11, 106, -1, -90, 109, -90, -90, -90, -90,
	-90, -90, -90, -19, 37, 20, -48, -34, -70, 4,
	14, -48, 29, 29, -81, -80, 11, -90, -90, -90,
	-70, -70, -70, -70, -70, -24, -23, -22, -29, -27,
	-26, -31, -41, -28, -30, -77, -79, 106, -70, -71,
	-72, -73, -74, -75, -76, -33, 71, 60, 64, 5,
	6, 101, 7, 9, 10, 68, 87, -24, -78, -77,
	-90, 12, -24, -15, 14, 55, 56, 57, 98, -82,
	69, -11, -20, 43, 40, -70, 16, 108, -70, 22,
	-47, -36, -34, -35, -37, 23, -26, 24, 14, -70,
	-70, -90, 108, 98, 81, -90, -90, -90, 20, 76,
	97, 96, 99, 67, -83, -89, 100, 101, 102, 103,
	104, 63, 62, 64, -24, -26, 105, 106, -42, -24,
	-26, -24, -72, -73, 106, 78, 61, 108, 99, -90,
	-16, 18, -47, -88, 58, -88, -88, -24, -50, -40,
	-39, -24, 102, 107, -21, 44, 6, -46, -25, -24,
	21, 106, -11, -57, -56, -23, -70, -48, -70, -16,
	-85, 54, -87, 51, 108, 46, 48, 49, 50, -70,
	22, -47, 106, 26, 27, 35, -81, -24, 82, -78,
	-77, -1, -24, -24, -24, -83, 65, 61, 66, 59,
	58, -24, -24, -24, -24, -24, -24, -24, 107, -70,
	-32, -82, -51, 75, -32, -2, -4, -3, -9, 72,
	91, 92, -70, -78, -22, -17, 38, -24, -13, -12,
	-13, -13, 108, 22, 6, 108, -84, 41, 42, -54,
	-53, 106, -44, -23, -15, 108, 99, 45, -85, -87,
	-86, 47, 45, -47, -70, -16, -49, -70, -61, 106,
	-70, -23, 106, -23, -11, -90, -67, -66, 77, 73,
	-74, -76, -24, 106, -26, -24, -26, -26, 107, 102,
	-45, -24, -43, -51, 77, -24, -19, 79, -2, -24,
	-90, -90, 78, -90, -18, 39, 40, -50, -70, -46,
	108, -45, 107, 108, -16, -57, -24, -36, 45, -86,
	45, -36, 107, 108, -63, 30, 31, 32, 33, -62,
	-61, 34, -44, 36, -90, 79, -67, -66, -1, -24,
	62, -45, 108, 79, -24, 76, 107, 88, 74, 76,
	-2, -24, -45, -54, 107, 21, -11, -44, -38, 52,
	53, -36, 45, -36, -49, -23, -23, 107, 108, -24,
	107, -70, 72, 79, 76, -24, 107, -45, -24, 5,
	-90, -2, -3, 79, -54, -24, 106, -36, -38, -63,
	-62, -90, 72, -1, 107, -69, -68, 77, 73, 74,
	-49, -90, 79, -69, -68, -2, -24, -90, 107, 72,
	79, 76, -90, 72, -2, -90,
}
var yyDef = [...]int{

	1, -2, 1, 227, 227, 227, 227, 227, 227, 227,
	227, 13, 14, 15, 16, 17, -2, 0, 0, 0,
	0, 0, 0, 227, 227, 227, 0, 0, 0, 0,
	0, 0, 0, 227, 0, 0, 48, 0, 0, 211,
	46, 203, 0, 2, 5, 228, 6, 7, 8, 9,
	10, 11, 12, 58, 0, 0, 0, 150, 114, 194,
	0, 0, 0, 0, 227, 209, 207, 21, 22, 23,
	0, 227, 227, 227, 0, 213, 70, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 80, 0, 68, 62,
	63, 64, 65, 66, 67, 108, 138, 0, 0, 195,
	196, 0, 198, 200, 201, 202, 0, 213, 0, 79,
	33, 0, -2, 50, 0, 223, 223, 223, 0, 0,
	212, 0, 60, 0, 0, 0, 0, 0, 115, 0,
	50, -2, 119, 120, 123, 124, 117, 118, 0, 0,
	0, 20, 0, 0, 0, 25, 26, 27, 0, 1,
	0, 225, 226, 213, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 214, 213, -2, 0, -2, 0, -2,
	98, 106, 197, 199, -2, 3, 0, 0, 0, 39,
	52, 0, 49, 0, 224, 0, 0, 206, 47, 154,
	135, -2, 134, 86, 40, 0, 59, 57, 146, -2,
	0, 0, 160, 48, 167, 0, 68, 151, 116, 169,
	0, -2, 217, 0, 0, 216, 220, 221, 222, 121,
	0, 50, 0, 0, 0, 0, 210, -2, 0, 227,
	204, 188, 87, -2, -2, 0, 0, 0, 0, 0,
	0, 99, 100, 101, 102, 103, 104, 105, 81, 69,
	0, 0, 140, 0, 56, 0, 3, 18, 19, 0,
	227, 227, 0, 205, 227, 54, 0, -2, 42, 45,
	43, 44, 0, 0, 61, 0, 82, 84, 85, 158,
	163, 0, 0, 142, 50, 0, 0, 0, 0, 217,
	0, 218, 0, 149, 122, 170, 0, 152, 178, 0,
	174, 183, 0, 0, 227, 28, 0, 188, 1, 0,
	90, 91, 213, 0, 94, -2, 96, 97, 107, 110,
	111, -2, 0, 157, 0, 213, 0, 0, 4, 213,
	36, 37, 3, 38, 41, 0, 0, 155, 136, 147,
	0, 0, 0, 0, 165, 168, -2, 130, 0, 0,
	0, 129, 171, 0, 172, 179, 180, 0, 0, 0,
	176, 0, 0, 0, 24, 0, 0, 187, 189, 213,
	0, 0, 0, 137, -2, 0, 112, 0, 227, 1,
	0, -2, 53, 164, 162, 0, 161, 143, 125, 0,
	0, 126, 0, 130, 153, 181, 182, 178, 0, -2,
	184, 185, 227, 0, 1, 92, 93, 145, -2, 0,
	31, 192, -2, 0, 159, -2, 0, 128, 127, 173,
	177, 29, 227, 186, 113, 0, 192, 3, 0, 227,
	0, 30, 0, 0, 191, 193, 213, 32, 132, 227,
	0, 3, 34, 227, 190, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 104, 3, 3,
	106, 107, 102, 100, 108, 101, 105, 103, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 109,
	3, 99,
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
	92, 93, 94, 95, 96, 97, 98,
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
		//line parser.y:148
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:153
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:160
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:165
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:172
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:176
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:180
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:184
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:188
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:192
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:196
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:200
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:204
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:208
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:212
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:216
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:220
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:226
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:230
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:236
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:240
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:246
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:250
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:256
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:260
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:264
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:268
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:272
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:278
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:282
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:286
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:290
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:294
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:300
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:304
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:308
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:312
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:318
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:322
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:328
		{
			yyVAL.expression = SelectQuery{
				SelectEntity:  yyDollar[1].expression,
				OrderByClause: yyDollar[2].expression,
				LimitClause:   yyDollar[3].expression,
				OffsetClause:  yyDollar[4].expression,
			}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:339
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
		//line parser.y:349
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
		//line parser.y:358
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
		//line parser.y:367
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
		//line parser.y:378
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:382
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:388
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:394
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:398
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:404
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:408
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:414
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:418
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:424
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:428
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:434
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:438
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 58:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:444
		{
			yyVAL.expression = nil
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:448
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: StrToInt64(yyDollar[2].token.Literal)}
		}
	case 60:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:454
		{
			yyVAL.expression = nil
		}
	case 61:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:458
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Number: StrToInt64(yyDollar[2].token.Literal)}
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:464
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:468
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:472
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:476
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:480
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:484
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:490
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:494
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:500
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:504
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:508
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:512
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:528
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:544
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:550
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 83:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:556
		{
			yyVAL.token = Token{}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:560
		{
			yyVAL.token = yyDollar[1].token
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:564
		{
			yyVAL.token = yyDollar[1].token
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:570
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:576
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
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:599
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:603
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 90:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:607
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 91:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:611
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 92:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:615
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 93:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 94:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:623
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 95:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:627
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 96:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:631
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:635
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:639
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:645
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:649
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:653
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:657
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:661
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:667
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:671
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:675
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 107:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:681
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:685
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 109:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:691
		{
			yyVAL.expression = Option{}
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:695
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:699
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 112:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:705
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 113:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:709
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:715
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 115:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:719
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:723
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:729
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:733
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:739
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:743
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 121:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:747
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:751
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 123:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:755
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:759
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 125:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:765
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 126:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:769
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 127:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:773
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 128:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:777
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 129:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:781
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 130:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:787
		{
			yyVAL.expression = nil
		}
	case 131:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:791
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 132:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:795
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:801
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:805
		{
			yyVAL.expression = AllColumns{}
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:811
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:815
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 137:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:821
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 138:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:827
		{
			yyVAL.expression = nil
		}
	case 139:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:831
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 140:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:837
		{
			yyVAL.expression = nil
		}
	case 141:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:841
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 142:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:847
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:851
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:861
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:867
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:871
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:877
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 149:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:881
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:887
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 151:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:891
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:897
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 153:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:901
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:907
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 155:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:911
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 156:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:917
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:921
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 158:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:927
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 159:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:931
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 160:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:935
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 161:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:939
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:945
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:951
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:955
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 165:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:961
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 166:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:967
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:973
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:977
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 169:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:983
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 170:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:988
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 171:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:995
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 172:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1001
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 173:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1005
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 174:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1011
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 175:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1015
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1021
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 177:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1025
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 178:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1031
		{
			yyVAL.expression = nil
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1035
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 180:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1039
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 181:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1043
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 182:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1047
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 183:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1053
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 184:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1057
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 185:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1063
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 186:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1069
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 187:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1073
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 188:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1079
		{
			yyVAL.procexpr = nil
		}
	case 189:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1083
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 190:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1089
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 191:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1093
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 192:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1099
		{
			yyVAL.procexpr = nil
		}
	case 193:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1103
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1109
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1115
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 196:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1121
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 197:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1125
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 198:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1132
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 199:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1136
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1143
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 201:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1149
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1155
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1161
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 204:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1167
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 205:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1171
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1177
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1183
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 208:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1187
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 209:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1193
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 210:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1197
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 211:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1203
		{
			yyVAL.token = Token{}
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1207
		{
			yyVAL.token = yyDollar[1].token
		}
	case 213:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1213
		{
			yyVAL.token = Token{}
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1217
		{
			yyVAL.token = yyDollar[1].token
		}
	case 215:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1223
		{
			yyVAL.token = Token{}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1227
		{
			yyVAL.token = yyDollar[1].token
		}
	case 217:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1233
		{
			yyVAL.token = Token{}
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1237
		{
			yyVAL.token = yyDollar[1].token
		}
	case 219:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1243
		{
			yyVAL.token = Token{}
		}
	case 220:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1247
		{
			yyVAL.token = yyDollar[1].token
		}
	case 221:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1251
		{
			yyVAL.token = yyDollar[1].token
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1255
		{
			yyVAL.token = yyDollar[1].token
		}
	case 223:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1261
		{
			yyVAL.token = Token{}
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1265
		{
			yyVAL.token = yyDollar[1].token
		}
	case 225:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1271
		{
			yyVAL.token = yyDollar[1].token
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1275
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 227:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1281
		{
			yyVAL.token = Token{}
		}
	case 228:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1285
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
