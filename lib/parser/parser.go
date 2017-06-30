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
const TIES = 57387
const PERCENT = 57388
const JOIN = 57389
const INNER = 57390
const OUTER = 57391
const LEFT = 57392
const RIGHT = 57393
const FULL = 57394
const CROSS = 57395
const ON = 57396
const USING = 57397
const NATURAL = 57398
const UNION = 57399
const INTERSECT = 57400
const EXCEPT = 57401
const ALL = 57402
const ANY = 57403
const EXISTS = 57404
const IN = 57405
const AND = 57406
const OR = 57407
const NOT = 57408
const BETWEEN = 57409
const LIKE = 57410
const IS = 57411
const NULL = 57412
const NULLS = 57413
const DISTINCT = 57414
const WITH = 57415
const CASE = 57416
const IF = 57417
const ELSEIF = 57418
const WHILE = 57419
const WHEN = 57420
const THEN = 57421
const ELSE = 57422
const DO = 57423
const END = 57424
const DECLARE = 57425
const CURSOR = 57426
const FOR = 57427
const FETCH = 57428
const OPEN = 57429
const CLOSE = 57430
const DISPOSE = 57431
const GROUP_CONCAT = 57432
const SEPARATOR = 57433
const COMMIT = 57434
const ROLLBACK = 57435
const CONTINUE = 57436
const BREAK = 57437
const EXIT = 57438
const PRINT = 57439
const VAR = 57440
const COMPARISON_OP = 57441
const STRING_OP = 57442
const SUBSTITUTION_OP = 57443

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
	"TIES",
	"PERCENT",
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
	"NULLS",
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

//line parser.y:1355

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
	57, 45,
	58, 45,
	59, 45,
	-2, 56,
	-1, 113,
	63, 228,
	67, 228,
	68, 228,
	-2, 242,
	-1, 132,
	47, 230,
	49, 234,
	-2, 166,
	-1, 167,
	57, 46,
	58, 46,
	59, 46,
	-2, 77,
	-1, 169,
	110, 129,
	-2, 226,
	-1, 171,
	78, 159,
	-2, 228,
	-1, 180,
	37, 129,
	91, 129,
	110, 129,
	-2, 226,
	-1, 197,
	63, 228,
	67, 228,
	68, 228,
	-2, 153,
	-1, 202,
	63, 228,
	67, 228,
	68, 228,
	-2, 61,
	-1, 205,
	63, 228,
	67, 228,
	68, 228,
	-2, 93,
	-1, 217,
	49, 234,
	-2, 230,
	-1, 233,
	63, 228,
	67, 228,
	68, 228,
	-2, 223,
	-1, 239,
	69, 0,
	99, 0,
	102, 0,
	-2, 100,
	-1, 240,
	69, 0,
	99, 0,
	102, 0,
	-2, 102,
	-1, 283,
	63, 228,
	67, 228,
	68, 228,
	-2, 51,
	-1, 290,
	63, 228,
	67, 228,
	68, 228,
	-2, 64,
	-1, 332,
	69, 0,
	99, 0,
	102, 0,
	-2, 111,
	-1, 336,
	63, 228,
	67, 228,
	68, 228,
	-2, 164,
	-1, 373,
	63, 228,
	67, 228,
	68, 228,
	-2, 181,
	-1, 399,
	82, 161,
	-2, 228,
	-1, 403,
	110, 86,
	111, 86,
	-2, 46,
	-1, 411,
	63, 228,
	67, 228,
	68, 228,
	-2, 55,
	-1, 431,
	63, 228,
	67, 228,
	68, 228,
	-2, 190,
	-1, 438,
	78, 174,
	80, 174,
	82, 174,
	-2, 228,
	-1, 446,
	94, 18,
	95, 18,
	-2, 1,
	-1, 449,
	63, 228,
	67, 228,
	68, 228,
	-2, 151,
}

const yyPrivate = 57344

const yyLast = 999

var yyAct = [...]int{

	80, 40, 468, 40, 381, 459, 44, 420, 314, 43,
	1, 46, 47, 48, 49, 50, 51, 52, 386, 300,
	186, 209, 53, 308, 94, 271, 76, 132, 2, 166,
	67, 68, 69, 324, 394, 194, 218, 203, 216, 114,
	111, 291, 77, 40, 92, 109, 37, 387, 85, 38,
	258, 38, 341, 64, 45, 86, 23, 189, 23, 155,
	41, 137, 39, 59, 18, 34, 19, 133, 17, 131,
	164, 142, 56, 154, 20, 119, 430, 21, 146, 147,
	148, 110, 136, 138, 380, 183, 57, 57, 167, 370,
	368, 183, 61, 221, 407, 222, 223, 224, 219, 176,
	303, 217, 16, 152, 151, 97, 153, 157, 158, 159,
	160, 161, 294, 406, 450, 39, 137, 288, 143, 59,
	185, 128, 472, 417, 275, 458, 32, 40, 39, 442,
	441, 440, 26, 432, 429, 30, 27, 28, 29, 379,
	137, 24, 25, 276, 277, 33, 35, 22, 369, 39,
	337, 40, 215, 256, 211, 336, 220, 174, 42, 263,
	237, 199, 349, 298, 347, 39, 59, 345, 42, 59,
	168, 169, 3, 264, 264, 191, 192, 157, 158, 159,
	160, 161, 40, 304, 228, 188, 39, 75, 108, 180,
	40, 113, 40, 40, 206, 235, 57, 232, 236, 38,
	42, 213, 159, 160, 161, 184, 23, 264, 40, 227,
	273, 42, 168, 151, 241, 122, 157, 158, 159, 160,
	161, 137, 119, 260, 42, 144, 234, 280, 156, 279,
	38, 270, 236, 145, 301, 40, 284, 23, 286, 287,
	474, 91, 323, 165, 264, 263, 264, 264, 313, 302,
	466, 90, 171, 447, 307, 177, 306, 319, 321, 435,
	122, 42, 101, 103, 167, 398, 329, 264, 346, 348,
	350, 320, 325, 40, 317, 193, 197, 316, 262, 265,
	202, 205, 207, 355, 356, 335, 328, 358, 392, 339,
	311, 352, 285, 351, 285, 285, 357, 261, 353, 208,
	233, 273, 261, 463, 342, 408, 137, 238, 239, 240,
	462, 137, 299, 247, 248, 249, 250, 251, 252, 253,
	164, 38, 327, 371, 362, 372, 175, 40, 23, 391,
	211, 376, 366, 374, 364, 462, 395, 327, 378, 461,
	389, 326, 477, 283, 179, 473, 403, 301, 403, 331,
	403, 333, 334, 172, 178, 456, 173, 290, 40, 393,
	102, 434, 100, 101, 103, 88, 104, 105, 293, 264,
	40, 121, 344, 367, 164, 38, 137, 267, 137, 104,
	122, 266, 23, 410, 58, 58, 273, 401, 182, 426,
	419, 412, 70, 71, 72, 73, 74, 301, 330, 190,
	332, 117, 264, 423, 243, 425, 38, 322, 242, 244,
	40, 427, 428, 23, 309, 444, 424, 343, 264, 237,
	377, 126, 269, 268, 129, 137, 58, 106, 140, 141,
	375, 354, 310, 452, 453, 445, 122, 40, 446, 305,
	106, 455, 246, 245, 197, 365, 457, 40, 460, 454,
	205, 299, 451, 299, 124, 299, 201, 361, 38, 464,
	373, 102, 40, 465, 125, 23, 467, 116, 117, 118,
	471, 360, 259, 282, 299, 221, 40, 222, 223, 224,
	476, 58, 54, 396, 479, 38, 390, 469, 229, 230,
	273, 388, 23, 212, 58, 38, 214, 231, 399, 63,
	225, 478, 23, 62, 273, 58, 187, 439, 289, 402,
	38, 404, 120, 405, 414, 415, 411, 23, 122, 149,
	122, 55, 122, 299, 38, 163, 162, 164, 127, 59,
	154, 23, 416, 59, 257, 41, 115, 39, 112, 18,
	34, 19, 418, 17, 431, 59, 41, 226, 278, 20,
	139, 130, 21, 437, 66, 60, 438, 443, 59, 65,
	152, 151, 93, 153, 157, 158, 159, 160, 161, 89,
	10, 254, 255, 212, 9, 8, 7, 449, 6, 210,
	221, 448, 222, 223, 224, 219, 58, 5, 217, 4,
	340, 170, 312, 82, 315, 318, 212, 212, 195, 31,
	196, 32, 135, 163, 162, 164, 134, 26, 154, 95,
	30, 27, 28, 29, 81, 84, 24, 25, 470, 78,
	33, 35, 22, 163, 162, 164, 83, 79, 154, 59,
	100, 101, 103, 42, 104, 105, 41, 413, 152, 151,
	295, 153, 157, 158, 159, 160, 161, 59, 100, 101,
	103, 45, 104, 105, 41, 363, 39, 204, 152, 151,
	200, 153, 157, 158, 159, 160, 161, 123, 359, 212,
	255, 58, 382, 383, 384, 385, 58, 281, 36, 15,
	274, 14, 13, 318, 12, 11, 212, 98, 272, 0,
	0, 99, 0, 0, 0, 106, 0, 0, 0, 96,
	0, 59, 100, 101, 103, 98, 104, 105, 41, 99,
	0, 0, 0, 106, 0, 107, 0, 96, 59, 100,
	101, 103, 0, 104, 105, 41, 0, 0, 0, 102,
	198, 0, 0, 107, 87, 0, 212, 0, 0, 0,
	0, 58, 0, 58, 0, 0, 315, 102, 0, 0,
	212, 212, 87, 0, 292, 0, 433, 0, 0, 98,
	0, 0, 0, 99, 0, 0, 0, 106, 0, 296,
	297, 96, 163, 162, 164, 0, 98, 154, 0, 0,
	99, 293, 0, 0, 106, 0, 0, 107, 96, 0,
	58, 0, 163, 162, 164, 0, 318, 154, 0, 0,
	0, 102, 338, 0, 107, 0, 87, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 315, 0, 102, 0,
	0, 0, 0, 87, 163, 162, 164, 152, 151, 154,
	153, 157, 158, 159, 160, 161, 163, 162, 164, 475,
	221, 154, 222, 223, 224, 219, 421, 422, 217, 0,
	0, 436, 0, 163, 162, 164, 0, 0, 154, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 409, 0,
	0, 152, 151, 0, 153, 157, 158, 159, 160, 161,
	0, 0, 163, 162, 164, 0, 0, 154, 152, 151,
	0, 153, 157, 158, 159, 160, 161, 400, 163, 162,
	164, 0, 0, 154, 0, 0, 0, 0, 0, 0,
	0, 163, 162, 164, 0, 181, 154, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 150, 163, 162, 164,
	0, 0, 154, 152, 151, 0, 153, 157, 158, 159,
	160, 161, 0, 397, 162, 164, 152, 151, 154, 153,
	157, 158, 159, 160, 161, 163, 0, 164, 0, 0,
	154, 0, 152, 151, 0, 153, 157, 158, 159, 160,
	161, 0, 0, 0, 0, 0, 0, 0, 152, 151,
	0, 153, 157, 158, 159, 160, 161, 0, 0, 0,
	152, 151, 0, 153, 157, 158, 159, 160, 161,
}
var yyPact = [...]int{

	524, -1000, 524, -58, -58, -58, -58, -58, -58, -58,
	-58, -1000, -1000, -1000, -1000, -1000, 445, 501, 554, 541,
	474, 470, 543, -58, -58, -58, 554, 554, 554, 554,
	554, 714, 714, -58, 526, 714, 522, 410, 121, 299,
	-1000, -1000, 152, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 411, 424, 554, 512, 10, 529, -1000,
	59, 536, 554, 554, -58, 7, 124, -1000, -1000, -1000,
	149, -58, -58, -58, 499, 847, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 121, -1000, 643, 62, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 714, 254, 91, 714,
	-1000, -1000, 256, -1000, -1000, -1000, -1000, 80, 834, 325,
	-26, -1000, 103, 539, 488, 59, 339, 339, 339, 714,
	625, -1000, 51, 412, 714, 714, 173, 554, 554, -1000,
	554, 488, 45, -1000, 525, -1000, -1000, -1000, -1000, 59,
	75, 462, -1000, 543, 714, 141, -1000, -1000, -1000, 535,
	524, 714, 714, 714, 308, 341, 382, 714, 714, 714,
	714, 714, 714, 714, -1000, 461, 43, -1000, 554, 299,
	219, 863, 50, 50, 314, 362, -1000, 4, -1000, -1000,
	299, 49, 554, 535, 357, -1000, 435, 714, -1000, 152,
	-1000, 152, 152, 863, -1000, 6, 486, 863, -1000, -1000,
	-1000, 714, 708, -1000, 1, 728, 50, 115, -1000, 522,
	-11, 81, 104, -1000, -1000, -1000, 392, 427, 365, 385,
	59, -1000, -1000, -1000, -1000, -1000, 554, 488, 554, 165,
	162, 554, -1000, 863, 152, -58, -20, 261, 74, 113,
	113, 370, 714, 50, 714, 50, 50, 97, 97, -1000,
	-1000, -1000, 891, 4, -1000, 714, -1000, -1000, 40, 697,
	224, 714, -1000, 643, -1000, -1000, 50, 58, 55, 53,
	445, 209, 49, -1000, -1000, 714, -58, -58, 215, -1000,
	-58, 432, 417, 863, 343, -1000, -1000, 343, 625, 554,
	863, -1000, 295, 400, 714, 302, -1000, -1000, -1000, -21,
	38, -22, 488, 554, 714, 59, 383, 365, 373, -1000,
	59, -1000, -1000, -1000, 29, -27, 642, 554, 457, -1000,
	554, 450, -58, -1000, 206, 261, 524, 714, -1000, -1000,
	879, -1000, 113, -1000, -1000, -1000, 559, -1000, -1000, -1000,
	183, 219, 714, 818, 323, 136, -1000, 136, -1000, 136,
	-1000, 3, 228, -1000, 789, -1000, -1000, 49, -1000, -1000,
	714, 714, -1000, -1000, -1000, -1000, -1000, 484, 50, 102,
	554, -1000, -1000, 863, 792, 59, 369, 59, 532, -1000,
	554, -1000, -1000, -1000, 554, 554, 24, -35, 714, 23,
	554, -1000, 286, 177, 246, -1000, 772, 714, -1000, 863,
	714, 50, 21, -1000, 20, 19, -1000, 552, -58, 49,
	171, 863, -1000, -1000, -1000, -1000, -1000, 50, -1000, -1000,
	-1000, 714, 5, 532, 59, 792, -1000, -1000, -1000, 642,
	554, 863, -1000, -1000, -58, 280, 524, 4, 863, -1000,
	-1000, -1000, -1000, 15, -1000, 259, 524, 226, -1000, 863,
	554, 532, -1000, -1000, -1000, -1000, -58, -1000, -1000, 168,
	259, 49, 714, -58, 12, -1000, 270, 158, 234, -1000,
	760, -1000, -1000, -58, 267, 49, -1000, -58, -1000, -1000,
}
var yyPgo = [...]int{

	0, 9, 25, 28, 688, 685, 684, 682, 681, 680,
	679, 172, 102, 46, 678, 39, 20, 677, 668, 22,
	667, 41, 660, 42, 26, 155, 105, 163, 37, 657,
	640, 637, 0, 627, 626, 619, 615, 614, 50, 609,
	67, 606, 27, 602, 7, 600, 598, 593, 591, 590,
	19, 29, 69, 72, 8, 35, 52, 589, 587, 579,
	21, 578, 576, 575, 47, 18, 4, 574, 570, 34,
	33, 2, 5, 365, 569, 251, 241, 44, 562, 24,
	48, 45, 55, 559, 53, 472, 59, 38, 23, 36,
	57, 228, 6,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10,
	11, 12, 12, 12, 12, 13, 13, 14, 15, 15,
	16, 16, 17, 17, 18, 18, 19, 19, 20, 20,
	20, 21, 21, 22, 22, 23, 23, 23, 23, 23,
	23, 24, 24, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 26, 26, 27, 27, 28,
	28, 29, 29, 30, 30, 30, 31, 31, 32, 33,
	34, 34, 34, 34, 34, 34, 34, 34, 34, 34,
	34, 34, 34, 34, 34, 34, 34, 34, 34, 35,
	35, 35, 35, 35, 36, 36, 36, 37, 37, 38,
	38, 38, 39, 39, 40, 40, 40, 41, 41, 42,
	42, 42, 42, 42, 42, 43, 43, 43, 43, 43,
	44, 44, 44, 45, 45, 46, 46, 47, 48, 48,
	49, 49, 50, 50, 51, 51, 52, 52, 53, 53,
	54, 54, 55, 55, 56, 56, 57, 57, 57, 57,
	58, 59, 60, 60, 61, 61, 62, 63, 63, 64,
	64, 65, 65, 66, 66, 66, 66, 66, 67, 67,
	68, 69, 69, 70, 70, 71, 71, 72, 72, 73,
	74, 75, 75, 76, 76, 77, 78, 79, 80, 81,
	81, 82, 83, 83, 84, 84, 85, 85, 86, 86,
	87, 87, 88, 88, 89, 89, 89, 89, 90, 90,
	91, 91, 92, 92,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 5, 8,
	9, 7, 9, 2, 8, 9, 2, 2, 5, 3,
	4, 5, 4, 4, 4, 1, 1, 3, 0, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 3,
	4, 0, 2, 0, 2, 1, 1, 1, 1, 1,
	1, 1, 3, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 3, 3, 1, 1, 3, 1,
	3, 2, 4, 0, 1, 1, 1, 1, 3, 3,
	3, 3, 3, 3, 4, 4, 6, 6, 4, 6,
	4, 4, 4, 6, 4, 4, 6, 4, 2, 3,
	3, 3, 3, 3, 3, 3, 2, 4, 1, 0,
	2, 2, 5, 7, 1, 2, 3, 1, 1, 1,
	1, 2, 3, 1, 1, 5, 5, 6, 6, 4,
	0, 2, 4, 1, 1, 1, 3, 5, 0, 1,
	0, 2, 1, 3, 1, 3, 1, 3, 1, 3,
	1, 3, 1, 3, 4, 2, 5, 8, 4, 7,
	6, 3, 1, 3, 4, 5, 6, 6, 8, 1,
	3, 1, 3, 0, 1, 1, 2, 2, 5, 7,
	7, 4, 2, 0, 2, 4, 2, 0, 2, 1,
	1, 1, 2, 1, 2, 1, 1, 1, 1, 1,
	3, 3, 1, 3, 1, 3, 0, 1, 0, 1,
	0, 1, 0, 1, 0, 1, 1, 1, 0, 1,
	1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -57, -58, -61, -62, -63, -67,
	-68, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 98, -82, 92, 93, 83, 87, 88, 89,
	86, 75, 77, 96, 16, 97, -14, -13, -80, 13,
	-32, 11, 109, -1, -92, 112, -92, -92, -92, -92,
	-92, -92, -92, -19, 37, 20, -53, -40, -73, 4,
	14, -53, 29, 29, -84, -83, 11, -92, -92, -92,
	-73, -73, -73, -73, -73, -25, -24, -23, -35, -33,
	-32, -37, -47, -34, -36, -80, -82, 109, -73, -74,
	-75, -76, -77, -78, -79, -39, 74, -26, 62, 66,
	5, 6, 104, 7, 9, 10, 70, 90, -25, -81,
	-80, -92, 12, -25, -15, 14, 57, 58, 59, 101,
	-85, 72, -11, -20, 43, 40, -73, 16, 111, -73,
	22, -52, -42, -40, -41, -43, 23, -32, 24, 14,
	-73, -73, -92, 111, 101, 84, -92, -92, -92, 20,
	79, 100, 99, 102, 69, -86, -91, 103, 104, 105,
	106, 107, 65, 64, 66, -25, -51, -32, 108, 109,
	-48, -25, 99, 102, -86, -91, -32, -25, -75, -76,
	109, 81, 63, 111, 102, -92, -16, 18, -52, -90,
	60, -90, -90, -25, -55, -46, -45, -25, 105, 110,
	-22, 44, -25, -28, -29, -25, 21, 109, -11, -60,
	-59, -24, -73, -53, -73, -16, -87, 56, -89, 53,
	111, 48, 50, 51, 52, -73, 22, -52, 109, 26,
	27, 35, -84, -25, 85, -81, -80, -1, -25, -25,
	-25, -86, 67, 63, 68, 61, 60, -25, -25, -25,
	-25, -25, -25, -25, 110, 111, 110, -73, -38, -85,
	-56, 78, -26, 109, -32, -26, 67, 63, 61, 60,
	-38, -2, -4, -3, -9, 75, 94, 95, -73, -81,
	-23, -17, 38, -25, -13, -12, -13, -13, 111, 22,
	-25, -21, 46, 73, 111, -30, 41, 42, -27, -26,
	-50, -24, -15, 111, 102, 47, -87, -89, -88, 49,
	47, -52, -73, -16, -54, -73, -64, 109, -73, -24,
	109, -24, -11, -92, -70, -69, 80, 76, -77, -79,
	-25, -26, -25, -26, -26, -51, -25, 110, 105, -51,
	-49, -56, 80, -25, -26, 109, -32, 109, -32, 109,
	-32, -19, 82, -2, -25, -92, -92, 81, -92, -18,
	39, 40, -55, -73, -21, 45, -28, 71, 111, 110,
	111, -16, -60, -25, -42, 47, -88, 47, -42, 110,
	111, -66, 30, 31, 32, 33, -65, -64, 34, -50,
	36, -92, 82, -70, -69, -1, -25, 64, 82, -25,
	79, 64, -27, -32, -27, -27, 110, 91, 77, 79,
	-2, -25, -51, -31, 30, 31, -27, 21, -11, -50,
	-44, 54, 55, -42, 47, -42, -54, -24, -24, 110,
	111, -25, 110, -73, 75, 82, 79, -25, -25, -26,
	110, 110, 110, 5, -92, -2, -3, 82, -27, -25,
	109, -42, -44, -66, -65, -92, 75, -1, 110, -72,
	-71, 80, 76, 77, -54, -92, 82, -72, -71, -2,
	-25, -92, 110, 75, 82, 79, -92, 75, -2, -92,
}
var yyDef = [...]int{

	1, -2, 1, 242, 242, 242, 242, 242, 242, 242,
	242, 13, 14, 15, 16, 17, -2, 0, 0, 0,
	0, 0, 0, 242, 242, 242, 0, 0, 0, 0,
	0, 0, 0, 242, 0, 0, 48, 0, 0, 226,
	46, 218, 0, 2, 5, 243, 6, 7, 8, 9,
	10, 11, 12, 58, 0, 0, 0, 168, 134, 209,
	0, 0, 0, 0, 242, 224, 222, 21, 22, 23,
	0, 242, 242, 242, 0, 228, 73, 74, 75, 76,
	77, 78, 79, 80, 81, 82, 83, 0, 71, 65,
	66, 67, 68, 69, 70, 128, 158, 228, 0, 0,
	210, 211, 0, 213, 215, 216, 217, 0, 228, 0,
	82, 33, 0, -2, 50, 0, 238, 238, 238, 0,
	0, 227, 0, 63, 0, 0, 0, 0, 0, 135,
	0, 50, -2, 139, 140, 143, 144, 137, 138, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	1, 0, 240, 241, 228, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 229, 228, 0, -2, 0, -2,
	0, -2, 240, 241, 0, 0, 118, 126, 212, 214,
	-2, 3, 0, 0, 0, 39, 52, 0, 49, 0,
	239, 0, 0, 221, 47, 172, 155, -2, 154, 98,
	40, 0, -2, 57, 89, -2, 0, 0, 178, 48,
	182, 0, 71, 169, 136, 184, 0, -2, 232, 0,
	0, 231, 235, 236, 237, 141, 0, 50, 0, 0,
	0, 0, 225, -2, 0, 242, 219, 203, 99, -2,
	-2, 0, 0, 0, 0, 0, 0, 119, 120, 121,
	122, 123, 124, 125, 84, 0, 85, 72, 0, 0,
	160, 0, 101, 0, 86, 103, 0, 0, 0, 0,
	56, 0, 3, 18, 19, 0, 242, 242, 0, 220,
	242, 54, 0, -2, 42, 45, 43, 44, 0, 0,
	-2, 59, 61, 0, 0, 91, 94, 95, 176, 87,
	0, 162, 50, 0, 0, 0, 0, 232, 0, 233,
	0, 167, 142, 185, 0, 170, 193, 0, 189, 198,
	0, 0, 242, 28, 0, 203, 1, 0, 104, 105,
	228, 108, -2, 112, 115, 165, -2, 127, 130, 131,
	0, 175, 0, 228, 0, 0, 110, 0, 114, 0,
	117, 0, 0, 4, 228, 36, 37, 3, 38, 41,
	0, 0, 173, 156, 60, 62, 90, 0, 0, 0,
	0, 180, 183, -2, 150, 0, 0, 0, 149, 186,
	0, 187, 194, 195, 0, 0, 0, 191, 0, 0,
	0, 24, 0, 0, 202, 204, 228, 0, 157, -2,
	0, 0, 0, -2, 0, 0, 132, 0, 242, 1,
	0, -2, 53, 92, 96, 97, 88, 0, 179, 163,
	145, 0, 0, 146, 0, 150, 171, 196, 197, 193,
	0, -2, 199, 200, 242, 0, 1, 106, -2, 107,
	109, 113, 116, 0, 31, 207, -2, 0, 177, -2,
	0, 148, 147, 188, 192, 29, 242, 201, 133, 0,
	207, 3, 0, 242, 0, 30, 0, 0, 206, 208,
	228, 32, 152, 242, 0, 3, 34, 242, 205, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 107, 3, 3,
	109, 110, 105, 103, 111, 104, 108, 106, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 112,
	3, 102,
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
	92, 93, 94, 95, 96, 97, 98, 99, 100, 101,
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
		//line parser.y:150
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:155
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:162
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:167
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:174
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:178
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:182
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:186
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:190
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:194
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:198
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:202
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:206
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:210
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:214
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:218
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:222
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:228
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:232
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:238
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:242
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:248
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:252
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:258
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:262
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:266
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:270
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:274
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:280
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:284
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:288
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:292
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:296
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:302
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:306
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:310
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:314
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:320
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:324
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:330
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
		//line parser.y:341
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
		//line parser.y:351
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
		//line parser.y:360
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
		//line parser.y:369
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
		//line parser.y:380
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:384
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:390
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:396
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:400
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:406
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:410
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:426
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:430
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:436
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:440
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 58:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:446
		{
			yyVAL.expression = nil
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:450
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:454
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:460
		{
			yyVAL.expression = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:464
		{
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:470
		{
			yyVAL.expression = nil
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:474
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:480
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:484
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:488
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:492
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:496
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:500
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:506
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:510
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:528
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:544
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:548
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:560
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:566
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:570
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:580
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:586
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:590
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 92:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:600
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:606
		{
			yyVAL.token = Token{}
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:610
		{
			yyVAL.token = yyDollar[1].token
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:614
		{
			yyVAL.token = yyDollar[1].token
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:620
		{
			yyVAL.token = yyDollar[1].token
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:624
		{
			yyVAL.token = yyDollar[1].token
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:630
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:636
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
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:659
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:663
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:667
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:671
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 104:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:675
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:679
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 106:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:683
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 107:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:687
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 108:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:691
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 109:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:695
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 110:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:699
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 111:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:703
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 112:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:707
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 113:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:711
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 114:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:715
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 115:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:719
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 116:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:723
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 117:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:727
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 118:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:731
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:737
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:741
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:745
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:749
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:753
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:759
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:763
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 126:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:767
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 127:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:773
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:777
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 129:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:783
		{
			yyVAL.expression = Option{}
		}
	case 130:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:787
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 131:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:791
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 132:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:797
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 133:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:801
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:807
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 135:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:811
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:815
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:821
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 138:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:825
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 139:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:831
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:835
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 141:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:839
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:843
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:847
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:851
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 145:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 146:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:861
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 147:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:865
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 148:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:869
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 149:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:873
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 150:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:879
		{
			yyVAL.expression = nil
		}
	case 151:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:883
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 152:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:887
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 153:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:893
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:897
		{
			yyVAL.expression = AllColumns{}
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:903
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 156:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:907
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 157:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:913
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 158:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:919
		{
			yyVAL.expression = nil
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:923
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 160:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:929
		{
			yyVAL.expression = nil
		}
	case 161:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:933
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:939
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 163:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:943
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 164:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:949
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 165:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:953
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:959
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 167:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:963
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 168:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:969
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 169:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:973
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 170:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:979
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 171:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:983
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:989
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 173:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:993
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 174:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:999
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 175:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1003
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 176:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1009
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 177:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1013
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 178:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1017
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 179:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1021
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 180:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1027
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 181:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1033
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 182:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1039
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 183:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1043
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 184:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1049
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 185:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1054
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 186:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1061
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 187:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1067
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 188:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1071
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1077
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 190:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1081
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1087
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 192:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1091
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 193:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1097
		{
			yyVAL.expression = nil
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1101
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1105
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 196:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1109
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 197:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1113
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 198:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1119
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 199:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1123
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 200:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1129
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 201:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1135
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 202:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1139
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 203:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1145
		{
			yyVAL.procexpr = nil
		}
	case 204:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1149
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 205:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1155
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 206:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1159
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 207:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1165
		{
			yyVAL.procexpr = nil
		}
	case 208:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1169
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 209:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1175
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 210:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1181
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 211:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1187
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 212:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1191
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 213:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1198
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 214:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1202
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 215:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1209
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1215
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1221
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1227
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 219:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1233
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 220:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1237
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 221:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1243
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1249
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 223:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1253
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1259
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 225:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1263
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 226:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1269
		{
			yyVAL.token = Token{}
		}
	case 227:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1273
		{
			yyVAL.token = yyDollar[1].token
		}
	case 228:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1279
		{
			yyVAL.token = Token{}
		}
	case 229:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1283
		{
			yyVAL.token = yyDollar[1].token
		}
	case 230:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1289
		{
			yyVAL.token = Token{}
		}
	case 231:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1293
		{
			yyVAL.token = yyDollar[1].token
		}
	case 232:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1299
		{
			yyVAL.token = Token{}
		}
	case 233:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1303
		{
			yyVAL.token = yyDollar[1].token
		}
	case 234:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1309
		{
			yyVAL.token = Token{}
		}
	case 235:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1313
		{
			yyVAL.token = yyDollar[1].token
		}
	case 236:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1317
		{
			yyVAL.token = yyDollar[1].token
		}
	case 237:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1321
		{
			yyVAL.token = yyDollar[1].token
		}
	case 238:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1327
		{
			yyVAL.token = Token{}
		}
	case 239:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1331
		{
			yyVAL.token = yyDollar[1].token
		}
	case 240:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1337
		{
			yyVAL.token = yyDollar[1].token
		}
	case 241:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1341
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 242:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1347
		{
			yyVAL.token = Token{}
		}
	case 243:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1351
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
