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
const NULLS = 57411
const DISTINCT = 57412
const WITH = 57413
const CASE = 57414
const IF = 57415
const ELSEIF = 57416
const WHILE = 57417
const WHEN = 57418
const THEN = 57419
const ELSE = 57420
const DO = 57421
const END = 57422
const DECLARE = 57423
const CURSOR = 57424
const FOR = 57425
const FETCH = 57426
const OPEN = 57427
const CLOSE = 57428
const DISPOSE = 57429
const GROUP_CONCAT = 57430
const SEPARATOR = 57431
const COMMIT = 57432
const ROLLBACK = 57433
const CONTINUE = 57434
const BREAK = 57435
const EXIT = 57436
const PRINT = 57437
const VAR = 57438
const COMPARISON_OP = 57439
const STRING_OP = 57440
const SUBSTITUTION_OP = 57441

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

//line parser.y:1340

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
	-1, 113,
	61, 225,
	65, 225,
	66, 225,
	-2, 239,
	-1, 132,
	45, 227,
	47, 231,
	-2, 163,
	-1, 167,
	55, 46,
	56, 46,
	57, 46,
	-2, 74,
	-1, 169,
	108, 126,
	-2, 223,
	-1, 171,
	76, 156,
	-2, 225,
	-1, 180,
	37, 126,
	89, 126,
	108, 126,
	-2, 223,
	-1, 197,
	61, 225,
	65, 225,
	66, 225,
	-2, 150,
	-1, 205,
	61, 225,
	65, 225,
	66, 225,
	-2, 90,
	-1, 217,
	47, 231,
	-2, 227,
	-1, 233,
	61, 225,
	65, 225,
	66, 225,
	-2, 220,
	-1, 239,
	67, 0,
	97, 0,
	100, 0,
	-2, 97,
	-1, 240,
	67, 0,
	97, 0,
	100, 0,
	-2, 99,
	-1, 283,
	61, 225,
	65, 225,
	66, 225,
	-2, 51,
	-1, 329,
	67, 0,
	97, 0,
	100, 0,
	-2, 108,
	-1, 333,
	61, 225,
	65, 225,
	66, 225,
	-2, 161,
	-1, 368,
	61, 225,
	65, 225,
	66, 225,
	-2, 178,
	-1, 394,
	80, 158,
	-2, 225,
	-1, 398,
	108, 83,
	109, 83,
	-2, 46,
	-1, 406,
	61, 225,
	65, 225,
	66, 225,
	-2, 55,
	-1, 426,
	61, 225,
	65, 225,
	66, 225,
	-2, 187,
	-1, 433,
	76, 171,
	78, 171,
	80, 171,
	-2, 225,
	-1, 441,
	92, 18,
	93, 18,
	-2, 1,
	-1, 444,
	61, 225,
	65, 225,
	66, 225,
	-2, 148,
}

const yyPrivate = 57344

const yyLast = 1081

var yyAct = [...]int{

	80, 40, 311, 40, 454, 463, 44, 2, 381, 43,
	1, 46, 47, 48, 49, 50, 51, 52, 376, 209,
	305, 203, 186, 415, 86, 23, 297, 23, 271, 321,
	67, 68, 69, 389, 166, 53, 131, 295, 194, 94,
	111, 218, 97, 40, 85, 38, 382, 38, 114, 92,
	109, 77, 258, 132, 216, 64, 338, 3, 155, 45,
	425, 137, 37, 119, 39, 59, 375, 133, 365, 56,
	16, 142, 412, 183, 402, 363, 183, 110, 146, 147,
	148, 300, 291, 189, 136, 138, 57, 57, 167, 61,
	163, 162, 164, 401, 39, 154, 288, 143, 128, 176,
	122, 221, 206, 222, 223, 224, 219, 151, 467, 217,
	157, 158, 159, 160, 161, 59, 137, 59, 453, 437,
	185, 436, 435, 427, 39, 152, 151, 40, 153, 157,
	158, 159, 160, 161, 424, 39, 254, 255, 374, 364,
	137, 334, 59, 39, 256, 122, 157, 158, 159, 160,
	161, 40, 188, 445, 215, 199, 174, 263, 42, 346,
	237, 344, 168, 169, 220, 342, 228, 180, 42, 159,
	160, 161, 42, 264, 264, 23, 227, 168, 333, 91,
	301, 119, 40, 184, 208, 144, 76, 156, 207, 273,
	40, 234, 40, 40, 236, 38, 57, 145, 213, 232,
	235, 191, 192, 469, 101, 103, 23, 264, 40, 90,
	75, 108, 461, 241, 113, 262, 265, 442, 42, 430,
	317, 137, 393, 387, 349, 354, 38, 260, 236, 263,
	261, 261, 339, 270, 279, 40, 280, 42, 458, 457,
	457, 164, 320, 456, 264, 314, 264, 264, 403, 296,
	310, 472, 284, 324, 286, 287, 468, 308, 299, 304,
	285, 451, 285, 285, 167, 122, 165, 264, 343, 345,
	347, 322, 303, 40, 172, 171, 313, 173, 177, 429,
	273, 326, 179, 352, 353, 175, 328, 355, 330, 331,
	332, 325, 319, 121, 336, 362, 164, 23, 193, 197,
	102, 350, 396, 137, 205, 182, 348, 324, 137, 341,
	104, 323, 178, 361, 211, 267, 190, 38, 117, 266,
	367, 122, 366, 233, 40, 371, 386, 359, 269, 268,
	238, 239, 240, 390, 246, 245, 247, 248, 249, 250,
	251, 252, 253, 398, 384, 398, 243, 398, 23, 419,
	242, 244, 388, 306, 372, 40, 369, 116, 117, 118,
	370, 373, 273, 307, 264, 40, 283, 302, 38, 106,
	124, 137, 221, 137, 222, 223, 224, 201, 421, 23,
	397, 358, 399, 405, 400, 296, 125, 296, 259, 296,
	357, 282, 414, 407, 298, 54, 385, 264, 383, 38,
	122, 411, 122, 63, 122, 40, 296, 229, 230, 62,
	439, 289, 441, 264, 237, 149, 231, 316, 318, 59,
	137, 327, 413, 329, 418, 55, 420, 59, 120, 23,
	409, 410, 40, 440, 449, 187, 450, 226, 127, 434,
	340, 452, 40, 448, 447, 130, 455, 115, 459, 38,
	443, 59, 139, 112, 351, 296, 23, 40, 460, 41,
	462, 60, 66, 88, 273, 466, 23, 197, 290, 202,
	205, 40, 438, 446, 59, 471, 38, 65, 273, 474,
	368, 23, 58, 58, 93, 464, 38, 211, 89, 10,
	70, 71, 72, 73, 74, 23, 9, 8, 7, 473,
	6, 38, 210, 391, 298, 100, 101, 103, 5, 104,
	105, 377, 378, 379, 380, 38, 4, 337, 394, 126,
	170, 82, 129, 195, 58, 196, 140, 141, 135, 134,
	221, 95, 222, 223, 224, 219, 406, 41, 217, 39,
	81, 18, 34, 19, 164, 17, 84, 154, 78, 83,
	79, 20, 298, 221, 21, 222, 223, 224, 219, 416,
	417, 217, 426, 408, 292, 204, 422, 423, 106, 200,
	123, 432, 356, 281, 433, 36, 15, 152, 151, 58,
	153, 157, 158, 159, 160, 161, 274, 14, 13, 12,
	11, 212, 58, 272, 214, 444, 0, 0, 225, 275,
	0, 32, 102, 58, 0, 0, 0, 26, 0, 0,
	30, 27, 28, 29, 0, 0, 24, 25, 276, 277,
	33, 35, 22, 0, 163, 162, 164, 0, 0, 154,
	0, 0, 257, 42, 0, 41, 465, 39, 0, 18,
	34, 19, 0, 17, 0, 0, 278, 0, 0, 20,
	0, 0, 21, 0, 0, 0, 0, 0, 0, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 0, 0,
	0, 212, 45, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 58, 0, 0, 0, 0, 0,
	309, 0, 312, 315, 212, 212, 0, 31, 0, 32,
	0, 163, 162, 164, 0, 26, 154, 0, 30, 27,
	28, 29, 0, 0, 24, 25, 0, 0, 33, 35,
	22, 0, 0, 59, 100, 101, 103, 0, 104, 105,
	41, 42, 0, 0, 0, 0, 152, 151, 0, 153,
	157, 158, 159, 160, 161, 0, 0, 0, 255, 59,
	100, 101, 103, 360, 104, 105, 41, 0, 39, 0,
	0, 0, 0, 0, 212, 0, 58, 59, 100, 101,
	103, 58, 104, 105, 41, 0, 0, 0, 315, 98,
	0, 212, 0, 99, 0, 0, 0, 106, 0, 0,
	0, 96, 0, 0, 59, 100, 101, 103, 0, 104,
	105, 41, 0, 0, 0, 98, 0, 107, 0, 99,
	0, 0, 0, 106, 0, 0, 0, 96, 0, 0,
	0, 102, 198, 98, 0, 0, 87, 99, 0, 212,
	0, 106, 0, 107, 58, 96, 58, 0, 0, 312,
	0, 0, 0, 212, 212, 0, 0, 102, 0, 428,
	98, 107, 87, 0, 99, 0, 0, 0, 106, 293,
	294, 0, 96, 0, 0, 102, 335, 0, 0, 0,
	87, 0, 0, 0, 0, 0, 0, 0, 107, 0,
	163, 162, 164, 58, 0, 154, 0, 0, 0, 315,
	0, 0, 102, 163, 162, 164, 0, 87, 154, 0,
	0, 0, 163, 162, 164, 0, 0, 154, 470, 312,
	0, 0, 0, 0, 0, 152, 151, 431, 153, 157,
	158, 159, 160, 161, 0, 0, 0, 0, 152, 151,
	0, 153, 157, 158, 159, 160, 161, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 163, 162, 164, 0,
	0, 154, 0, 0, 0, 163, 162, 164, 0, 0,
	154, 404, 0, 0, 0, 0, 163, 162, 164, 0,
	395, 154, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 152, 151, 181, 153, 157, 158, 159, 160, 161,
	152, 151, 0, 153, 157, 158, 159, 160, 161, 0,
	0, 152, 151, 0, 153, 157, 158, 159, 160, 161,
	163, 162, 164, 0, 0, 154, 0, 0, 0, 163,
	162, 164, 0, 0, 154, 150, 0, 0, 392, 162,
	164, 0, 0, 154, 0, 0, 0, 163, 0, 164,
	0, 0, 154, 0, 0, 152, 151, 0, 153, 157,
	158, 159, 160, 161, 152, 151, 0, 153, 157, 158,
	159, 160, 161, 152, 151, 0, 153, 157, 158, 159,
	160, 161, 152, 151, 0, 153, 157, 158, 159, 160,
	161,
}
var yyPact = [...]int{

	624, -1000, 624, -51, -51, -51, -51, -51, -51, -51,
	-51, -1000, -1000, -1000, -1000, -1000, 358, 405, 470, 447,
	380, 374, 451, -51, -51, -51, 470, 470, 470, 470,
	470, 790, 790, -51, 441, 790, 433, 302, 82, 223,
	-1000, -1000, 130, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 327, 346, 470, 422, -11, 423, -1000,
	61, 438, 470, 470, -51, -12, 86, -1000, -1000, -1000,
	115, -51, -51, -51, 395, 948, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 82, -1000, 745, 56, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 790, 177, 65, 790,
	-1000, -1000, 198, -1000, -1000, -1000, -1000, 60, 904, 244,
	-36, -1000, 83, 562, 417, 61, 258, 258, 258, 790,
	719, -1000, 47, 333, 463, 790, 81, 470, 470, -1000,
	470, 417, 55, -1000, 415, -1000, -1000, -1000, -1000, 61,
	59, 381, -1000, 451, 790, 108, -1000, -1000, -1000, 448,
	624, 790, 790, 790, 232, 285, 276, 790, 790, 790,
	790, 790, 790, 790, -1000, 28, 36, -1000, 470, 223,
	155, 957, 50, 50, 254, 270, -1000, 480, -1000, -1000,
	223, 526, 470, 448, 500, -1000, 353, 790, -1000, 130,
	-1000, 130, 130, 957, -1000, -13, 389, 957, -1000, -1000,
	-1000, 462, -1000, -1000, -27, 818, 50, 111, -1000, 433,
	-28, 80, 71, -1000, -1000, -1000, 322, 326, 306, 318,
	61, -1000, -1000, -1000, -1000, -1000, 470, 417, 470, 138,
	113, 470, -1000, 957, 130, -51, -33, 233, 45, 9,
	9, 301, 790, 50, 790, 50, 50, 66, 66, -1000,
	-1000, -1000, 975, 480, -1000, 790, -1000, -1000, 33, 763,
	154, 790, -1000, 745, -1000, -1000, 50, 58, 54, 52,
	358, 144, 526, -1000, -1000, 790, -51, -51, 146, -1000,
	-51, 351, 341, 957, 262, -1000, -1000, 262, 719, 470,
	-1000, 790, 226, -1000, -1000, -1000, -34, 31, -41, 417,
	470, 790, 61, 315, 306, 309, -1000, 61, -1000, -1000,
	-1000, 30, -43, 481, 470, 364, -1000, 470, 360, -51,
	-1000, 143, 233, 624, 790, -1000, -1000, 966, -1000, 9,
	-1000, -1000, -1000, 639, -1000, -1000, -1000, 142, 155, 790,
	893, 240, 122, -1000, 122, -1000, 122, -1000, -15, 173,
	-1000, 884, -1000, -1000, 526, -1000, -1000, 790, 790, -1000,
	-1000, -1000, 400, 50, 51, 470, -1000, -1000, 957, 507,
	61, 304, 61, 484, -1000, 470, -1000, -1000, -1000, 470,
	470, 26, -49, 790, 15, 470, -1000, 206, 139, 179,
	-1000, 840, 790, -1000, 957, 790, 50, 14, -1000, 13,
	11, -1000, 467, -51, 526, 137, 957, -1000, -1000, -1000,
	-1000, -1000, 50, -1000, -1000, -1000, 790, 46, 484, 61,
	507, -1000, -1000, -1000, 481, 470, 957, -1000, -1000, -51,
	188, 624, 480, 957, -1000, -1000, -1000, -1000, 10, -1000,
	165, 624, 163, -1000, 957, 470, 484, -1000, -1000, -1000,
	-1000, -51, -1000, -1000, 132, 165, 526, 790, -51, 0,
	-1000, 183, 123, 166, -1000, 831, -1000, -1000, -51, 178,
	526, -1000, -51, -1000, -1000,
}
var yyPgo = [...]int{

	0, 9, 28, 7, 593, 590, 589, 588, 587, 586,
	576, 57, 70, 62, 575, 48, 22, 573, 572, 35,
	570, 569, 51, 186, 178, 42, 37, 21, 565, 564,
	563, 0, 550, 549, 548, 546, 540, 52, 531, 67,
	529, 53, 528, 23, 525, 523, 521, 520, 517, 26,
	34, 36, 69, 2, 38, 56, 516, 508, 502, 19,
	500, 498, 497, 46, 8, 18, 496, 489, 33, 29,
	5, 4, 463, 488, 209, 179, 49, 484, 39, 44,
	50, 24, 477, 55, 388, 58, 54, 20, 41, 83,
	187, 6,
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
	24, 24, 25, 25, 26, 26, 27, 27, 28, 28,
	29, 29, 29, 30, 30, 31, 32, 33, 33, 33,
	33, 33, 33, 33, 33, 33, 33, 33, 33, 33,
	33, 33, 33, 33, 33, 33, 34, 34, 34, 34,
	34, 35, 35, 35, 36, 36, 37, 37, 37, 38,
	38, 39, 39, 39, 40, 40, 41, 41, 41, 41,
	41, 41, 42, 42, 42, 42, 42, 43, 43, 43,
	44, 44, 45, 45, 46, 47, 47, 48, 48, 49,
	49, 50, 50, 51, 51, 52, 52, 53, 53, 54,
	54, 55, 55, 56, 56, 56, 56, 57, 58, 59,
	59, 60, 60, 61, 62, 62, 63, 63, 64, 64,
	65, 65, 65, 65, 65, 66, 66, 67, 68, 68,
	69, 69, 70, 70, 71, 71, 72, 73, 74, 74,
	75, 75, 76, 77, 78, 79, 80, 80, 81, 82,
	82, 83, 83, 84, 84, 85, 85, 86, 86, 87,
	87, 88, 88, 88, 88, 89, 89, 90, 90, 91,
	91,
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
	1, 3, 3, 1, 1, 3, 1, 3, 2, 4,
	0, 1, 1, 1, 1, 3, 3, 3, 3, 3,
	3, 4, 4, 6, 6, 4, 6, 4, 4, 4,
	6, 4, 4, 6, 4, 2, 3, 3, 3, 3,
	3, 3, 3, 2, 4, 1, 0, 2, 2, 5,
	7, 1, 2, 3, 1, 1, 1, 1, 2, 3,
	1, 1, 5, 5, 6, 6, 4, 0, 2, 4,
	1, 1, 1, 3, 5, 0, 1, 0, 2, 1,
	3, 1, 3, 1, 3, 1, 3, 1, 3, 1,
	3, 4, 2, 5, 8, 4, 7, 6, 3, 1,
	3, 4, 5, 6, 6, 8, 1, 3, 1, 3,
	0, 1, 1, 2, 2, 5, 7, 7, 4, 2,
	0, 2, 4, 2, 0, 2, 1, 1, 1, 2,
	1, 2, 1, 1, 1, 1, 1, 3, 3, 1,
	3, 1, 3, 0, 1, 0, 1, 0, 1, 0,
	1, 0, 1, 1, 1, 0, 1, 1, 1, 0,
	1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -56, -57, -60, -61, -62, -66,
	-67, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 96, -81, 90, 91, 81, 85, 86, 87,
	84, 73, 75, 94, 16, 95, -14, -13, -79, 13,
	-31, 11, 107, -1, -91, 110, -91, -91, -91, -91,
	-91, -91, -91, -19, 37, 20, -52, -39, -72, 4,
	14, -52, 29, 29, -83, -82, 11, -91, -91, -91,
	-72, -72, -72, -72, -72, -24, -23, -22, -34, -32,
	-31, -36, -46, -33, -35, -79, -81, 107, -72, -73,
	-74, -75, -76, -77, -78, -38, 72, -25, 60, 64,
	5, 6, 102, 7, 9, 10, 68, 88, -24, -80,
	-79, -91, 12, -24, -15, 14, 55, 56, 57, 99,
	-84, 70, -11, -20, 43, 40, -72, 16, 109, -72,
	22, -51, -41, -39, -40, -42, 23, -31, 24, 14,
	-72, -72, -91, 109, 99, 82, -91, -91, -91, 20,
	77, 98, 97, 100, 67, -85, -90, 101, 102, 103,
	104, 105, 63, 62, 64, -24, -50, -31, 106, 107,
	-47, -24, 97, 100, -85, -90, -31, -24, -74, -75,
	107, 79, 61, 109, 100, -91, -16, 18, -51, -89,
	58, -89, -89, -24, -54, -45, -44, -24, 103, 108,
	-21, 44, 6, -27, -28, -24, 21, 107, -11, -59,
	-58, -23, -72, -52, -72, -16, -86, 54, -88, 51,
	109, 46, 48, 49, 50, -72, 22, -51, 107, 26,
	27, 35, -83, -24, 83, -80, -79, -1, -24, -24,
	-24, -85, 65, 61, 66, 59, 58, -24, -24, -24,
	-24, -24, -24, -24, 108, 109, 108, -72, -37, -84,
	-55, 76, -25, 107, -31, -25, 65, 61, 59, 58,
	-37, -2, -4, -3, -9, 73, 92, 93, -72, -80,
	-22, -17, 38, -24, -13, -12, -13, -13, 109, 22,
	6, 109, -29, 41, 42, -26, -25, -49, -23, -15,
	109, 100, 45, -86, -88, -87, 47, 45, -51, -72,
	-16, -53, -72, -63, 107, -72, -23, 107, -23, -11,
	-91, -69, -68, 78, 74, -76, -78, -24, -25, -24,
	-25, -25, -50, -24, 108, 103, -50, -48, -55, 78,
	-24, -25, 107, -31, 107, -31, 107, -31, -19, 80,
	-2, -24, -91, -91, 79, -91, -18, 39, 40, -54,
	-72, -27, 69, 109, 108, 109, -16, -59, -24, -41,
	45, -87, 45, -41, 108, 109, -65, 30, 31, 32,
	33, -64, -63, 34, -49, 36, -91, 80, -69, -68,
	-1, -24, 62, 80, -24, 77, 62, -26, -31, -26,
	-26, 108, 89, 75, 77, -2, -24, -50, -30, 30,
	31, -26, 21, -11, -49, -43, 52, 53, -41, 45,
	-41, -53, -23, -23, 108, 109, -24, 108, -72, 73,
	80, 77, -24, -24, -25, 108, 108, 108, 5, -91,
	-2, -3, 80, -26, -24, 107, -41, -43, -65, -64,
	-91, 73, -1, 108, -71, -70, 78, 74, 75, -53,
	-91, 80, -71, -70, -2, -24, -91, 108, 73, 80,
	77, -91, 73, -2, -91,
}
var yyDef = [...]int{

	1, -2, 1, 239, 239, 239, 239, 239, 239, 239,
	239, 13, 14, 15, 16, 17, -2, 0, 0, 0,
	0, 0, 0, 239, 239, 239, 0, 0, 0, 0,
	0, 0, 0, 239, 0, 0, 48, 0, 0, 223,
	46, 215, 0, 2, 5, 240, 6, 7, 8, 9,
	10, 11, 12, 58, 0, 0, 0, 165, 131, 206,
	0, 0, 0, 0, 239, 221, 219, 21, 22, 23,
	0, 239, 239, 239, 0, 225, 70, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 80, 0, 68, 62,
	63, 64, 65, 66, 67, 125, 155, 225, 0, 0,
	207, 208, 0, 210, 212, 213, 214, 0, 225, 0,
	79, 33, 0, -2, 50, 0, 235, 235, 235, 0,
	0, 224, 0, 60, 0, 0, 0, 0, 0, 132,
	0, 50, -2, 136, 137, 140, 141, 134, 135, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	1, 0, 237, 238, 225, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 226, 225, 0, -2, 0, -2,
	0, -2, 237, 238, 0, 0, 115, 123, 209, 211,
	-2, 3, 0, 0, 0, 39, 52, 0, 49, 0,
	236, 0, 0, 218, 47, 169, 152, -2, 151, 95,
	40, 0, 59, 57, 86, -2, 0, 0, 175, 48,
	179, 0, 68, 166, 133, 181, 0, -2, 229, 0,
	0, 228, 232, 233, 234, 138, 0, 50, 0, 0,
	0, 0, 222, -2, 0, 239, 216, 200, 96, -2,
	-2, 0, 0, 0, 0, 0, 0, 116, 117, 118,
	119, 120, 121, 122, 81, 0, 82, 69, 0, 0,
	157, 0, 98, 0, 83, 100, 0, 0, 0, 0,
	56, 0, 3, 18, 19, 0, 239, 239, 0, 217,
	239, 54, 0, -2, 42, 45, 43, 44, 0, 0,
	61, 0, 88, 91, 92, 173, 84, 0, 159, 50,
	0, 0, 0, 0, 229, 0, 230, 0, 164, 139,
	182, 0, 167, 190, 0, 186, 195, 0, 0, 239,
	28, 0, 200, 1, 0, 101, 102, 225, 105, -2,
	109, 112, 162, -2, 124, 127, 128, 0, 172, 0,
	225, 0, 0, 107, 0, 111, 0, 114, 0, 0,
	4, 225, 36, 37, 3, 38, 41, 0, 0, 170,
	153, 87, 0, 0, 0, 0, 177, 180, -2, 147,
	0, 0, 0, 146, 183, 0, 184, 191, 192, 0,
	0, 0, 188, 0, 0, 0, 24, 0, 0, 199,
	201, 225, 0, 154, -2, 0, 0, 0, -2, 0,
	0, 129, 0, 239, 1, 0, -2, 53, 89, 93,
	94, 85, 0, 176, 160, 142, 0, 0, 143, 0,
	147, 168, 193, 194, 190, 0, -2, 196, 197, 239,
	0, 1, 103, -2, 104, 106, 110, 113, 0, 31,
	204, -2, 0, 174, -2, 0, 145, 144, 185, 189,
	29, 239, 198, 130, 0, 204, 3, 0, 239, 0,
	30, 0, 0, 203, 205, 225, 32, 149, 239, 0,
	3, 34, 239, 202, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 105, 3, 3,
	107, 108, 103, 101, 109, 102, 106, 104, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 110,
	3, 100,
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
	92, 93, 94, 95, 96, 97, 98, 99,
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
		//line parser.y:149
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:154
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:161
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:166
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:173
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:177
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:181
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:185
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:189
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:193
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:197
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:201
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:205
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:209
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:213
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:217
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:221
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:227
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:231
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:237
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:241
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:247
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:251
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:257
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:261
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:265
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:269
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:273
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:279
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:283
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:287
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:291
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:295
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:301
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:305
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:309
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:313
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:319
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:323
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:329
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
		//line parser.y:340
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
		//line parser.y:350
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
		//line parser.y:359
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
		//line parser.y:368
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
		//line parser.y:379
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:383
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:389
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:399
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:405
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:409
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:415
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:419
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:425
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:429
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:435
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:439
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 58:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:445
		{
			yyVAL.expression = nil
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:449
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: StrToInt64(yyDollar[2].token.Literal)}
		}
	case 60:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:455
		{
			yyVAL.expression = nil
		}
	case 61:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:459
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Number: StrToInt64(yyDollar[2].token.Literal)}
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:465
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:469
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:473
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:477
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:481
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:485
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:491
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:495
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:501
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:505
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:509
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:513
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:517
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:521
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:525
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:529
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:533
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:537
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:541
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:545
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:551
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:555
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:561
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:565
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:571
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:575
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:581
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 89:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:585
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 90:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:591
		{
			yyVAL.token = Token{}
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:595
		{
			yyVAL.token = yyDollar[1].token
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:599
		{
			yyVAL.token = yyDollar[1].token
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:605
		{
			yyVAL.token = yyDollar[1].token
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:609
		{
			yyVAL.token = yyDollar[1].token
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:615
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:621
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
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:644
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:648
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:652
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:656
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 101:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:660
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 102:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:664
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 103:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:668
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 104:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:672
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:676
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 106:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:680
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 107:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:684
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 108:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:688
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 109:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:692
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 110:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:696
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 111:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:700
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 112:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:704
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 113:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:708
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 114:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:712
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 115:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:716
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:722
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:726
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 118:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:730
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:734
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:738
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:744
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:748
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 123:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:752
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 124:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:758
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 125:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:762
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 126:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:768
		{
			yyVAL.expression = Option{}
		}
	case 127:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:772
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 128:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:776
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 129:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:782
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 130:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:786
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 132:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:796
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 133:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:800
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:806
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:810
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:816
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:820
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 138:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 139:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:828
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:832
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 141:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:836
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 142:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:842
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 143:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:846
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 144:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:850
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 145:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:854
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 146:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:858
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 147:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:864
		{
			yyVAL.expression = nil
		}
	case 148:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:868
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 149:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:872
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:878
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 151:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:882
		{
			yyVAL.expression = AllColumns{}
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:888
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 153:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:892
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 154:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:898
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 155:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:904
		{
			yyVAL.expression = nil
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:908
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 157:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:914
		{
			yyVAL.expression = nil
		}
	case 158:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:918
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:924
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 160:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:928
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:934
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:938
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:944
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:948
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:954
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 166:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:958
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:964
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:968
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:974
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 170:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:978
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 171:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:984
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 172:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:988
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 173:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:994
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 174:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:998
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 175:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1002
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 176:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1006
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 177:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1012
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 178:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1018
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1024
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 180:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1028
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 181:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1034
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 182:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1039
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 183:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1046
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 184:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1052
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 185:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1056
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 186:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1062
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 187:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1066
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 188:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1072
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 189:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1076
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 190:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1082
		{
			yyVAL.expression = nil
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1086
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 192:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1090
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 193:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1094
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 194:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1098
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 195:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1104
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 196:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1108
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 197:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1114
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 198:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1120
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 199:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1124
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 200:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1130
		{
			yyVAL.procexpr = nil
		}
	case 201:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1134
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 202:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1140
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 203:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1144
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 204:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1150
		{
			yyVAL.procexpr = nil
		}
	case 205:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1154
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 206:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1160
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1166
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 208:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1172
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 209:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1176
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 210:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1183
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 211:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1187
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1194
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 213:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1200
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1206
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 215:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1212
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1218
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 217:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1222
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 218:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1228
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 219:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1234
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 220:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1238
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 221:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1244
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 222:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1248
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 223:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1254
		{
			yyVAL.token = Token{}
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1258
		{
			yyVAL.token = yyDollar[1].token
		}
	case 225:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1264
		{
			yyVAL.token = Token{}
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1268
		{
			yyVAL.token = yyDollar[1].token
		}
	case 227:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1274
		{
			yyVAL.token = Token{}
		}
	case 228:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1278
		{
			yyVAL.token = yyDollar[1].token
		}
	case 229:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1284
		{
			yyVAL.token = Token{}
		}
	case 230:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1288
		{
			yyVAL.token = yyDollar[1].token
		}
	case 231:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1294
		{
			yyVAL.token = Token{}
		}
	case 232:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1298
		{
			yyVAL.token = yyDollar[1].token
		}
	case 233:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1302
		{
			yyVAL.token = yyDollar[1].token
		}
	case 234:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1306
		{
			yyVAL.token = yyDollar[1].token
		}
	case 235:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1312
		{
			yyVAL.token = Token{}
		}
	case 236:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1316
		{
			yyVAL.token = yyDollar[1].token
		}
	case 237:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1322
		{
			yyVAL.token = yyDollar[1].token
		}
	case 238:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1326
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 239:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1332
		{
			yyVAL.token = Token{}
		}
	case 240:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1336
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
