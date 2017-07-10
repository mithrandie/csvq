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
const RECURSIVE = 57367
const CREATE = 57368
const ADD = 57369
const DROP = 57370
const ALTER = 57371
const TABLE = 57372
const FIRST = 57373
const LAST = 57374
const AFTER = 57375
const BEFORE = 57376
const DEFAULT = 57377
const RENAME = 57378
const TO = 57379
const ORDER = 57380
const GROUP = 57381
const HAVING = 57382
const BY = 57383
const ASC = 57384
const DESC = 57385
const LIMIT = 57386
const OFFSET = 57387
const TIES = 57388
const PERCENT = 57389
const JOIN = 57390
const INNER = 57391
const OUTER = 57392
const LEFT = 57393
const RIGHT = 57394
const FULL = 57395
const CROSS = 57396
const ON = 57397
const USING = 57398
const NATURAL = 57399
const UNION = 57400
const INTERSECT = 57401
const EXCEPT = 57402
const ALL = 57403
const ANY = 57404
const EXISTS = 57405
const IN = 57406
const AND = 57407
const OR = 57408
const NOT = 57409
const BETWEEN = 57410
const LIKE = 57411
const IS = 57412
const NULL = 57413
const NULLS = 57414
const DISTINCT = 57415
const WITH = 57416
const CASE = 57417
const IF = 57418
const ELSEIF = 57419
const WHILE = 57420
const WHEN = 57421
const THEN = 57422
const ELSE = 57423
const DO = 57424
const END = 57425
const DECLARE = 57426
const CURSOR = 57427
const FOR = 57428
const FETCH = 57429
const OPEN = 57430
const CLOSE = 57431
const DISPOSE = 57432
const NEXT = 57433
const PRIOR = 57434
const ABSOLUTE = 57435
const RELATIVE = 57436
const RANGE = 57437
const GROUP_CONCAT = 57438
const SEPARATOR = 57439
const PARTITION = 57440
const OVER = 57441
const COMMIT = 57442
const ROLLBACK = 57443
const CONTINUE = 57444
const BREAK = 57445
const EXIT = 57446
const PRINT = 57447
const VAR = 57448
const COMPARISON_OP = 57449
const STRING_OP = 57450
const SUBSTITUTION_OP = 57451

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
	"RECURSIVE",
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
	"NEXT",
	"PRIOR",
	"ABSOLUTE",
	"RELATIVE",
	"RANGE",
	"GROUP_CONCAT",
	"SEPARATOR",
	"PARTITION",
	"OVER",
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
	"'('",
	"')'",
	"','",
	"'.'",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:1489

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
	-1, 0,
	1, 1,
	-2, 74,
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	1, 1,
	77, 1,
	81, 1,
	83, 1,
	-2, 74,
	-1, 46,
	58, 54,
	59, 54,
	60, 54,
	-2, 65,
	-1, 115,
	64, 251,
	68, 251,
	69, 251,
	-2, 267,
	-1, 148,
	4, 34,
	-2, 251,
	-1, 149,
	4, 35,
	-2, 251,
	-1, 150,
	77, 1,
	81, 1,
	83, 1,
	-2, 74,
	-1, 168,
	117, 147,
	-2, 249,
	-1, 170,
	79, 182,
	-2, 251,
	-1, 180,
	38, 147,
	97, 147,
	117, 147,
	-2, 249,
	-1, 181,
	83, 3,
	-2, 74,
	-1, 198,
	48, 253,
	50, 257,
	-2, 189,
	-1, 216,
	64, 251,
	68, 251,
	69, 251,
	-2, 175,
	-1, 226,
	64, 251,
	68, 251,
	69, 251,
	-2, 246,
	-1, 231,
	70, 0,
	107, 0,
	110, 0,
	-2, 118,
	-1, 232,
	70, 0,
	107, 0,
	110, 0,
	-2, 120,
	-1, 265,
	77, 3,
	81, 3,
	83, 3,
	-2, 74,
	-1, 280,
	64, 251,
	68, 251,
	69, 251,
	-2, 70,
	-1, 284,
	64, 251,
	68, 251,
	69, 251,
	-2, 109,
	-1, 297,
	50, 257,
	-2, 253,
	-1, 310,
	64, 251,
	68, 251,
	69, 251,
	-2, 60,
	-1, 317,
	117, 147,
	-2, 249,
	-1, 330,
	83, 1,
	-2, 74,
	-1, 336,
	70, 0,
	107, 0,
	110, 0,
	-2, 129,
	-1, 340,
	64, 251,
	68, 251,
	69, 251,
	-2, 187,
	-1, 362,
	83, 3,
	-2, 74,
	-1, 366,
	64, 251,
	68, 251,
	69, 251,
	-2, 73,
	-1, 417,
	83, 184,
	-2, 251,
	-1, 428,
	77, 1,
	81, 1,
	83, 1,
	-2, 74,
	-1, 441,
	64, 251,
	68, 251,
	69, 251,
	-2, 204,
	-1, 447,
	64, 251,
	68, 251,
	69, 251,
	-2, 64,
	-1, 455,
	64, 251,
	68, 251,
	69, 251,
	-2, 213,
	-1, 460,
	77, 1,
	81, 1,
	83, 1,
	-2, 74,
	-1, 462,
	79, 197,
	81, 197,
	83, 197,
	-2, 251,
	-1, 471,
	77, 1,
	81, 1,
	83, 1,
	-2, 18,
	-1, 497,
	83, 3,
	-2, 74,
	-1, 502,
	64, 251,
	68, 251,
	69, 251,
	-2, 173,
	-1, 521,
	77, 3,
	81, 3,
	83, 3,
	-2, 74,
}

const yyPrivate = 57344

const yyLast = 1065

var yyAct = [...]int{

	37, 2, 403, 318, 198, 39, 40, 41, 42, 43,
	44, 45, 264, 121, 85, 34, 495, 34, 510, 384,
	166, 60, 61, 62, 89, 289, 482, 36, 1, 86,
	20, 113, 20, 398, 250, 374, 206, 376, 281, 367,
	328, 213, 129, 412, 112, 298, 217, 111, 63, 65,
	66, 67, 95, 197, 155, 296, 93, 199, 140, 404,
	16, 312, 124, 116, 57, 345, 144, 145, 146, 38,
	167, 209, 317, 126, 126, 167, 454, 438, 46, 51,
	138, 139, 168, 436, 77, 167, 3, 301, 426, 302,
	303, 304, 299, 147, 151, 297, 397, 157, 158, 159,
	160, 161, 340, 379, 120, 64, 125, 125, 425, 52,
	370, 183, 128, 183, 315, 137, 185, 52, 286, 48,
	194, 49, 186, 47, 202, 204, 177, 141, 524, 523,
	522, 75, 110, 494, 473, 115, 157, 158, 159, 160,
	161, 136, 276, 187, 466, 465, 464, 456, 192, 137,
	453, 195, 126, 173, 449, 126, 300, 437, 431, 219,
	33, 396, 341, 248, 220, 34, 159, 160, 161, 380,
	506, 33, 76, 64, 503, 136, 148, 149, 229, 255,
	20, 500, 364, 266, 208, 64, 353, 351, 349, 221,
	180, 165, 249, 54, 103, 105, 34, 184, 273, 46,
	170, 164, 255, 176, 211, 212, 225, 271, 98, 233,
	92, 20, 54, 287, 120, 263, 219, 54, 292, 126,
	54, 294, 142, 188, 91, 305, 488, 156, 518, 227,
	126, 272, 467, 143, 295, 252, 277, 216, 285, 520,
	508, 171, 472, 273, 172, 226, 319, 322, 292, 292,
	275, 459, 125, 137, 230, 231, 232, 293, 416, 307,
	239, 240, 241, 242, 243, 244, 245, 266, 339, 274,
	360, 361, 343, 329, 424, 363, 327, 356, 358, 288,
	34, 365, 320, 410, 357, 324, 333, 362, 137, 311,
	332, 313, 314, 280, 284, 20, 253, 321, 423, 498,
	104, 331, 319, 497, 253, 330, 346, 499, 427, 498,
	310, 331, 292, 528, 326, 179, 137, 355, 519, 492,
	458, 479, 33, 369, 135, 126, 174, 408, 409, 178,
	435, 388, 378, 80, 262, 164, 106, 334, 419, 336,
	219, 394, 136, 383, 389, 34, 322, 259, 182, 292,
	53, 258, 395, 382, 387, 210, 347, 393, 413, 385,
	20, 132, 406, 235, 266, 433, 291, 234, 236, 486,
	411, 359, 261, 260, 33, 429, 251, 34, 238, 237,
	254, 257, 366, 445, 443, 420, 442, 421, 386, 422,
	381, 446, 20, 279, 190, 219, 323, 325, 108, 372,
	373, 450, 526, 444, 292, 440, 126, 392, 432, 434,
	137, 126, 137, 448, 137, 439, 191, 285, 216, 102,
	103, 105, 319, 106, 107, 137, 292, 292, 469, 134,
	471, 391, 457, 175, 414, 309, 136, 122, 136, 407,
	136, 470, 405, 34, 335, 56, 337, 338, 485, 417,
	487, 430, 131, 132, 133, 55, 229, 490, 20, 491,
	377, 203, 119, 292, 203, 222, 223, 348, 126, 474,
	126, 53, 478, 284, 224, 34, 481, 316, 301, 322,
	302, 303, 304, 441, 64, 108, 34, 489, 493, 496,
	20, 504, 207, 507, 447, 375, 228, 377, 137, 266,
	513, 20, 306, 71, 72, 256, 256, 515, 455, 123,
	511, 126, 34, 509, 505, 501, 476, 477, 461, 64,
	527, 462, 193, 266, 480, 64, 104, 20, 319, 531,
	130, 525, 205, 468, 529, 127, 34, 196, 114, 203,
	35, 59, 64, 53, 118, 53, 53, 530, 58, 94,
	90, 20, 291, 399, 400, 401, 402, 10, 375, 9,
	375, 137, 375, 69, 70, 73, 74, 8, 7, 256,
	6, 256, 256, 290, 451, 452, 163, 162, 164, 5,
	301, 154, 302, 303, 304, 299, 502, 514, 297, 4,
	344, 169, 256, 350, 352, 354, 82, 64, 102, 103,
	105, 512, 106, 107, 35, 163, 162, 164, 214, 215,
	154, 377, 201, 200, 517, 516, 96, 81, 152, 151,
	256, 153, 157, 158, 159, 160, 161, 163, 463, 164,
	84, 38, 154, 78, 203, 83, 79, 64, 102, 103,
	105, 475, 106, 107, 35, 375, 371, 152, 151, 283,
	153, 157, 158, 159, 160, 161, 99, 246, 247, 282,
	100, 117, 278, 189, 108, 390, 308, 50, 97, 152,
	151, 15, 153, 157, 158, 159, 160, 161, 101, 267,
	14, 87, 68, 256, 13, 256, 12, 256, 375, 109,
	11, 265, 64, 102, 103, 105, 99, 106, 107, 35,
	100, 0, 0, 0, 108, 104, 218, 33, 97, 88,
	0, 0, 0, 0, 0, 203, 0, 0, 101, 0,
	203, 164, 0, 0, 154, 0, 0, 0, 0, 109,
	0, 64, 102, 103, 105, 0, 106, 107, 35, 0,
	0, 0, 0, 0, 0, 104, 0, 0, 0, 88,
	0, 99, 0, 256, 0, 100, 0, 0, 0, 108,
	0, 152, 151, 97, 153, 157, 158, 159, 160, 161,
	256, 0, 0, 101, 0, 0, 0, 203, 0, 203,
	0, 0, 0, 0, 109, 0, 0, 163, 162, 164,
	99, 368, 154, 0, 100, 0, 0, 0, 108, 0,
	104, 342, 97, 0, 88, 0, 0, 0, 0, 163,
	162, 164, 101, 256, 154, 0, 0, 0, 369, 0,
	203, 0, 0, 109, 0, 0, 163, 162, 164, 152,
	151, 154, 153, 157, 158, 159, 160, 161, 0, 104,
	247, 521, 0, 88, 0, 163, 162, 164, 0, 0,
	154, 152, 151, 0, 153, 157, 158, 159, 160, 161,
	460, 0, 163, 162, 164, 0, 0, 154, 152, 151,
	0, 153, 157, 158, 159, 160, 161, 428, 0, 0,
	0, 163, 162, 164, 0, 0, 154, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 418, 163, 162, 164,
	0, 0, 154, 0, 152, 151, 0, 153, 157, 158,
	159, 160, 161, 0, 181, 0, 0, 163, 162, 164,
	0, 0, 154, 152, 151, 0, 153, 157, 158, 159,
	160, 161, 150, 163, 162, 164, 35, 0, 154, 152,
	151, 31, 153, 157, 158, 159, 160, 161, 415, 162,
	164, 17, 0, 154, 18, 0, 0, 0, 0, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 0, 35,
	0, 0, 0, 0, 31, 152, 151, 0, 153, 157,
	158, 159, 160, 161, 17, 0, 0, 18, 0, 0,
	152, 151, 0, 153, 157, 158, 159, 160, 161, 33,
	0, 268, 0, 29, 0, 0, 0, 0, 0, 23,
	0, 0, 27, 24, 25, 26, 301, 0, 302, 303,
	304, 299, 483, 484, 297, 21, 22, 269, 270, 30,
	32, 19, 33, 0, 28, 0, 29, 0, 0, 0,
	0, 0, 23, 0, 0, 27, 24, 25, 26, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 21, 22,
	0, 0, 30, 32, 19,
}
var yyPact = [...]int{

	958, -1000, 958, -51, -51, -51, -51, -51, -51, -51,
	-51, -1000, -1000, -1000, -1000, -1000, 104, 425, 415, 530,
	-51, -51, -51, 538, 538, 538, 538, 472, 727, 727,
	-51, 526, 727, 437, 105, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 399, 489, 538, 521,
	516, 394, 251, -1000, 248, 538, 538, -51, 9, 113,
	-1000, -1000, -1000, 148, -1000, -51, -51, -51, 538, -1000,
	-1000, -1000, -1000, 727, 727, 852, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 105, -1000, -1000, 633, -34,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 727, 134, 77,
	727, 538, -1000, -1000, 188, -1000, -1000, -1000, -1000, 74,
	832, 284, -5, -1000, 87, 511, -1000, 4, 538, -1000,
	727, 350, 375, 538, 506, 2, 515, 101, 518, 474,
	101, 294, 294, 294, 593, -1000, 47, 96, 73, 438,
	-1000, 530, 727, 143, -1000, -1000, -1000, 476, 868, 868,
	958, 727, 727, 727, 268, 299, 317, 727, 727, 727,
	727, 727, 727, 727, -1000, 540, 46, 538, 251, 217,
	868, 63, 63, 283, 311, -1000, 654, 264, -1000, -1000,
	251, 925, 538, 529, 414, -1000, 437, 120, 868, 348,
	727, 727, 97, 538, 538, -1000, 538, 474, 38, -1000,
	480, -1000, -1000, -1000, -1000, 101, 396, 727, -1000, 96,
	-1000, 96, 96, -1000, -4, 455, 868, -1000, -1000, -44,
	-1000, 538, 181, 169, 538, -1000, 868, 248, 529, 224,
	25, -14, -14, 327, 727, 63, 727, 63, 63, 53,
	53, -1000, -1000, -1000, 562, 654, -1000, 727, -1000, -1000,
	45, 688, 225, 727, -1000, 633, -1000, -1000, 63, 72,
	71, 70, 268, 399, 201, 925, -1000, -1000, 727, -51,
	-51, 205, -1000, -7, -51, -1000, 66, 538, -1000, 727,
	744, -1000, -8, 357, 868, -1000, 63, 538, -1000, 516,
	-15, 59, -49, -1000, -1000, -1000, 342, 429, 309, 340,
	101, -1000, -1000, -1000, -1000, -1000, 538, 474, 391, 366,
	868, 302, -1000, -1000, 302, 593, 538, 251, 44, -22,
	522, 538, 407, -1000, 538, 402, -51, -51, 200, 224,
	958, 727, -1000, -1000, 883, -1000, -14, -1000, -1000, -1000,
	722, -1000, -1000, -1000, 175, 217, 727, 816, 273, 86,
	-1000, 86, -1000, 86, -1000, 210, -9, 230, -1000, 797,
	-1000, -1000, 925, -1000, 248, 41, 868, -1000, 249, 319,
	727, 258, -1000, -1000, -1000, -35, 40, -41, 474, 538,
	727, 101, 336, 309, 335, -1000, 101, -1000, -1000, -1000,
	-1000, 727, 727, -1000, -1000, 37, -1000, 538, -1000, -1000,
	-1000, 538, 538, 33, -42, 727, 30, 538, -1000, -1000,
	244, 168, 234, -1000, 780, 727, -1000, 868, 727, 63,
	29, 28, 27, -1000, 137, -1000, 528, -51, 925, 159,
	17, 447, -1000, -1000, -1000, 485, 63, 300, 538, -1000,
	-1000, 868, 967, 101, 321, 101, 531, 868, -1000, 127,
	-1000, -1000, -1000, 522, 538, 868, -1000, -1000, -51, 243,
	958, 654, 868, -1000, -1000, -1000, -1000, -1000, 16, -1000,
	222, 958, 229, -1000, 65, -1000, -1000, -1000, -1000, 63,
	-1000, -1000, -1000, 727, 58, 531, 101, 967, 54, -1000,
	-1000, -1000, -51, -1000, -1000, 157, 222, 925, 727, -51,
	248, -1000, 868, 538, 531, -1000, 130, -1000, 242, 156,
	232, -1000, 761, -1000, 13, 12, 11, 399, 361, -51,
	237, 925, -1000, -1000, -1000, -1000, 727, -1000, -51, -1000,
	-1000, -1000,
}
var yyPgo = [...]int{

	0, 27, 12, 1, 691, 690, 686, 684, 682, 681,
	680, 679, 671, 86, 61, 79, 667, 42, 36, 666,
	665, 13, 663, 39, 662, 60, 661, 63, 84, 172,
	102, 208, 35, 38, 659, 649, 646, 641, 333, 636,
	635, 633, 630, 617, 34, 616, 46, 615, 614, 57,
	613, 4, 612, 26, 609, 608, 596, 591, 590, 37,
	20, 53, 62, 3, 41, 65, 589, 579, 573, 25,
	570, 568, 567, 59, 2, 33, 559, 557, 43, 40,
	18, 16, 24, 550, 224, 210, 56, 549, 52, 14,
	47, 29, 548, 64, 376, 54, 55, 19, 45, 71,
	544, 227, 0,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 8, 8, 9, 9, 10, 10,
	10, 10, 10, 11, 11, 11, 11, 12, 12, 13,
	14, 14, 14, 14, 15, 15, 16, 17, 17, 18,
	18, 19, 19, 20, 20, 21, 21, 22, 22, 22,
	23, 23, 24, 24, 25, 25, 26, 26, 27, 27,
	28, 28, 28, 28, 28, 28, 29, 29, 30, 30,
	30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
	30, 31, 31, 32, 32, 33, 33, 34, 34, 35,
	35, 36, 36, 36, 37, 37, 38, 39, 40, 40,
	40, 40, 40, 40, 40, 40, 40, 40, 40, 40,
	40, 40, 40, 40, 40, 40, 40, 41, 41, 41,
	41, 41, 42, 42, 42, 43, 43, 44, 44, 44,
	45, 45, 46, 47, 48, 48, 49, 49, 49, 50,
	50, 51, 51, 51, 51, 51, 51, 52, 52, 52,
	52, 52, 53, 53, 53, 54, 54, 54, 55, 55,
	56, 57, 57, 58, 58, 59, 59, 60, 60, 61,
	61, 62, 62, 63, 63, 64, 64, 65, 65, 66,
	66, 66, 66, 67, 68, 69, 69, 70, 70, 71,
	72, 72, 73, 73, 74, 74, 75, 75, 75, 75,
	75, 76, 76, 77, 78, 78, 79, 79, 80, 80,
	81, 81, 82, 83, 84, 84, 85, 85, 86, 87,
	88, 89, 90, 90, 91, 92, 92, 93, 93, 94,
	94, 95, 95, 96, 96, 97, 97, 98, 98, 98,
	98, 99, 99, 100, 100, 101, 101, 102, 102,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 6, 0,
	1, 1, 1, 1, 2, 2, 5, 6, 8, 9,
	7, 9, 2, 8, 9, 2, 2, 5, 3, 5,
	5, 4, 4, 4, 1, 1, 3, 0, 2, 0,
	2, 0, 3, 0, 2, 0, 3, 0, 3, 4,
	0, 2, 0, 2, 0, 2, 6, 9, 1, 3,
	1, 1, 1, 1, 1, 1, 1, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	3, 3, 1, 1, 3, 1, 3, 2, 4, 1,
	1, 0, 1, 1, 1, 1, 3, 3, 3, 3,
	3, 3, 4, 4, 6, 6, 4, 6, 4, 4,
	4, 6, 4, 4, 6, 4, 2, 3, 3, 3,
	3, 3, 3, 3, 2, 4, 1, 0, 2, 2,
	5, 7, 8, 2, 0, 3, 1, 2, 3, 1,
	1, 1, 1, 2, 3, 1, 1, 5, 5, 6,
	6, 4, 0, 2, 4, 1, 1, 1, 1, 3,
	5, 0, 1, 0, 2, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 1, 3, 4, 2, 6,
	9, 5, 8, 7, 3, 1, 3, 5, 6, 6,
	6, 8, 1, 3, 1, 3, 0, 1, 1, 2,
	2, 5, 7, 7, 4, 2, 0, 2, 4, 2,
	0, 2, 1, 1, 1, 2, 1, 2, 1, 1,
	1, 1, 1, 3, 3, 1, 3, 1, 3, 0,
	1, 0, 1, 0, 1, 0, 1, 0, 1, 1,
	1, 0, 1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -13, -66, -67, -70, -71, -72, -76,
	-77, -5, -6, -7, -10, -12, -25, 26, 29, 106,
	-91, 100, 101, 84, 88, 89, 90, 87, 76, 78,
	104, 16, 105, 74, -89, 11, -1, -102, 120, -102,
	-102, -102, -102, -102, -102, -102, -14, 19, 15, 17,
	-16, -15, 13, -38, 116, 30, 30, -93, -92, 11,
	-102, -102, -102, -82, 4, -82, -82, -82, -8, 91,
	92, 31, 32, 93, 94, -30, -29, -28, -41, -39,
	-38, -43, -56, -40, -42, -89, -91, -9, 116, -82,
	-83, -84, -85, -86, -87, -88, -45, 75, -31, 63,
	67, 85, 5, 6, 112, 7, 9, 10, 71, 96,
	-30, -90, -89, -102, 12, -30, -27, -26, -100, 25,
	109, -21, 38, 20, -62, -49, -82, 14, -62, -17,
	14, 58, 59, 60, -94, 73, -13, -25, -82, -82,
	-102, 118, 109, 85, -102, -102, -102, -82, -30, -30,
	80, 108, 107, 110, 70, -95, -101, 111, 112, 113,
	114, 115, 66, 65, 67, -30, -60, 119, 116, -57,
	-30, 107, 110, -95, -101, -38, -30, -82, -84, -85,
	116, 82, 64, 118, 110, -102, 118, -82, -30, -22,
	44, 41, -82, 16, 118, -82, 22, -61, -51, -49,
	-50, -52, 23, -38, 24, 14, -18, 18, -61, -99,
	61, -99, -99, -64, -55, -54, -30, -46, 113, -82,
	117, 116, 27, 28, 36, -93, -30, 86, 20, -1,
	-30, -30, -30, -95, 68, 64, 69, 62, 61, -30,
	-30, -30, -30, -30, -30, -30, 117, 118, 117, -82,
	-44, -94, -65, 79, -31, 116, -38, -31, 68, 64,
	62, 61, 70, -44, -2, -4, -3, -11, 76, 102,
	103, -82, -90, -89, -28, -27, 22, 116, -24, 45,
	-30, -33, -34, -35, -30, -46, 21, 116, -13, -69,
	-68, -29, -82, -62, -82, -18, -96, 57, -98, 54,
	118, 49, 51, 52, 53, -82, 22, -61, -19, 39,
	-30, -15, -14, -15, -15, 118, 22, 116, -63, -82,
	-73, 116, -82, -29, 116, -29, -13, -90, -79, -78,
	81, 77, -86, -88, -30, -31, -30, -31, -31, -60,
	-30, 117, 113, -60, -58, -65, 81, -30, -31, 116,
	-38, 116, -38, 116, -38, -95, -21, 83, -2, -30,
	-102, -102, 82, -102, 116, -63, -30, -23, 47, 74,
	118, -36, 42, 43, -32, -31, -59, -29, -17, 118,
	110, 48, -96, -98, -97, 50, 48, -61, -82, -18,
	-20, 40, 41, -64, -82, -44, 117, 118, -75, 31,
	32, 33, 34, -74, -73, 35, -59, 37, -102, -102,
	83, -79, -78, -1, -30, 65, 83, -30, 80, 65,
	-32, -32, -32, 88, 64, 117, 97, 78, 80, -2,
	-13, 117, -23, 46, -33, 72, 118, 117, 118, -18,
	-69, -30, -51, 48, -97, 48, -51, -30, -60, 117,
	-63, -29, -29, 117, 118, -30, 117, -82, 76, 83,
	80, -30, -30, -31, 117, 117, 117, 95, 5, -102,
	-2, -3, 83, 117, 22, -37, 31, 32, -32, 21,
	-13, -59, -53, 55, 56, -51, 48, -51, 99, -75,
	-74, -102, 76, -1, 117, -81, -80, 81, 77, 78,
	116, -32, -30, 116, -51, -53, 116, -102, 83, -81,
	-80, -2, -30, -102, -13, -63, -47, -48, 98, 76,
	83, 80, 117, 117, 117, -21, 41, -102, 76, -2,
	-60, -102,
}
var yyDef = [...]int{

	-2, -2, -2, 267, 267, 267, 267, 267, 267, 267,
	267, 13, 14, 15, 16, 17, 0, 0, 0, 0,
	267, 267, 267, 0, 0, 0, 0, 29, 0, 0,
	267, 0, 0, 263, 0, 241, 2, 5, 268, 6,
	7, 8, 9, 10, 11, 12, -2, 0, 0, 0,
	57, 0, 249, 55, 74, 0, 0, 267, 247, 245,
	21, 22, 23, 0, 232, 267, 267, 267, 0, 30,
	31, 32, 33, 0, 0, 251, 88, 89, 90, 91,
	92, 93, 94, 95, 96, 97, 98, 99, 74, 86,
	80, 81, 82, 83, 84, 85, 146, 181, 251, 0,
	0, 0, 233, 234, 0, 236, 238, 239, 240, 0,
	251, 0, 97, 42, 0, -2, 75, 78, 0, 264,
	0, 67, 0, 0, 0, 191, 156, 0, 0, 59,
	0, 261, 261, 261, 0, 250, 0, 0, 0, 0,
	20, 0, 0, 0, 25, 26, 27, 0, -2, -2,
	-2, 0, 265, 266, 251, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 252, 251, 0, 0, -2, 0,
	-2, 265, 266, 0, 0, 136, 144, 0, 235, 237,
	-2, -2, 0, 0, 0, 48, 263, 0, 244, 72,
	0, 0, 74, 0, 0, 157, 0, 59, -2, 161,
	162, 165, 166, 159, 160, 0, 61, 0, 58, 0,
	262, 0, 0, 56, 195, 178, -2, 176, 177, 86,
	116, 0, 0, 0, 0, 248, -2, 74, 0, 226,
	117, -2, -2, 0, 0, 0, 0, 0, 0, 137,
	138, 139, 140, 141, 142, 143, 100, 0, 101, 87,
	0, 0, 183, 0, 119, 74, 102, 121, 0, 0,
	0, 0, 251, 65, 0, -2, 18, 19, 0, 267,
	267, 0, 243, 242, 267, 79, 0, 0, 49, 0,
	-2, 66, 105, 111, -2, 110, 0, 0, 201, 57,
	205, 0, 86, 192, 158, 207, 0, -2, 255, 0,
	0, 254, 258, 259, 260, 163, 0, 59, 63, 0,
	-2, 51, 54, 52, 53, 0, 0, -2, 0, 193,
	216, 0, 212, 221, 0, 0, 267, 267, 0, 226,
	-2, 0, 122, 123, 251, 126, -2, 130, 133, 188,
	-2, 145, 148, 149, 0, 198, 0, 251, 0, 74,
	128, 74, 132, 74, 135, 0, 0, 0, 4, 251,
	45, 46, -2, 47, 74, 0, -2, 68, 70, 0,
	0, 107, 112, 113, 199, 103, 0, 185, 59, 0,
	0, 0, 0, 255, 0, 256, 0, 190, 164, 208,
	50, 0, 0, 196, 179, 0, 209, 0, 210, 217,
	218, 0, 0, 0, 214, 0, 0, 0, 24, 28,
	0, 0, 225, 227, 251, 0, 180, -2, 0, 0,
	0, 0, 0, 36, 0, 150, 0, 267, -2, 0,
	0, 0, 69, 71, 106, 0, 0, 74, 0, 203,
	206, -2, 172, 0, 0, 0, 171, -2, 62, 145,
	194, 219, 220, 216, 0, -2, 222, 223, 267, 0,
	-2, 124, -2, 125, 127, 131, 134, 37, 0, 40,
	230, -2, 0, 76, 0, 108, 114, 115, 104, 0,
	202, 186, 167, 0, 0, 168, 0, 172, 0, 211,
	215, 38, 267, 224, 151, 0, 230, -2, 0, 267,
	74, 200, -2, 0, 170, 169, 154, 39, 0, 0,
	229, 231, 251, 41, 0, 0, 0, 65, 0, 267,
	0, -2, 77, 174, 152, 153, 0, 43, 267, 228,
	155, 44,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 115, 3, 3,
	116, 117, 113, 111, 118, 112, 119, 114, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 120,
	3, 110,
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
	102, 103, 104, 105, 106, 107, 108, 109,
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
		//line parser.y:162
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:167
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:174
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:179
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:186
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:190
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:194
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:198
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:202
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:206
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:210
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:214
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:218
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:222
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:226
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:230
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:234
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:240
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:244
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:250
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:254
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:260
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:264
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:270
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:274
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:278
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:282
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:286
		{
			yyVAL.statement = FetchCursor{Position: yyDollar[2].expression, Cursor: yyDollar[3].identifier, Variables: yyDollar[5].variables}
		}
	case 29:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:292
		{
			yyVAL.expression = nil
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:296
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:300
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:304
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:308
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:312
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token, Number: yyDollar[2].expression}
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:316
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token, Number: yyDollar[2].expression}
		}
	case 36:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:322
		{
			yyVAL.expression = CursorStatus{CursorLit: yyDollar[1].token.Literal, Cursor: yyDollar[2].identifier, Is: yyDollar[3].token.Literal, Negation: yyDollar[4].token, Type: yyDollar[5].token.Token, TypeLit: yyDollar[5].token.Literal}
		}
	case 37:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:326
		{
			yyVAL.expression = CursorStatus{CursorLit: yyDollar[1].token.Literal, Cursor: yyDollar[2].identifier, Is: yyDollar[3].token.Literal, Negation: yyDollar[4].token, Type: yyDollar[6].token.Token, TypeLit: yyDollar[5].token.Literal + " " + yyDollar[6].token.Literal}
		}
	case 38:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:332
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 39:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:336
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 40:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:340
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 41:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:344
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 42:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:348
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 43:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:354
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 44:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:358
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 45:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:362
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:366
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 47:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:372
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:376
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 49:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:382
		{
			yyVAL.expression = SelectQuery{
				CommonTableClause: yyDollar[1].expression,
				SelectEntity:      yyDollar[2].expression,
				OrderByClause:     yyDollar[3].expression,
				LimitClause:       yyDollar[4].expression,
				OffsetClause:      yyDollar[5].expression,
			}
		}
	case 50:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:394
		{
			yyVAL.expression = SelectEntity{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
			}
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:404
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:413
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:422
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:433
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:437
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:443
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 57:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:449
		{
			yyVAL.expression = nil
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:453
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 59:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:459
		{
			yyVAL.expression = nil
		}
	case 60:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:463
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:469
		{
			yyVAL.expression = nil
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:473
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:479
		{
			yyVAL.expression = nil
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:483
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:489
		{
			yyVAL.expression = nil
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:493
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 67:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:499
		{
			yyVAL.expression = nil
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:503
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 69:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:507
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
		}
	case 70:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:513
		{
			yyVAL.expression = nil
		}
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:517
		{
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 72:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:523
		{
			yyVAL.expression = nil
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:527
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 74:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:533
		{
			yyVAL.expression = nil
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:537
		{
			yyVAL.expression = CommonTableClause{With: yyDollar[1].token.Literal, CommonTables: yyDollar[2].expressions}
		}
	case 76:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:543
		{
			yyVAL.expression = CommonTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, As: yyDollar[3].token.Literal, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 77:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:547
		{
			yyVAL.expression = CommonTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, Columns: yyDollar[4].expressions, As: yyDollar[6].token.Literal, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:553
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:557
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:563
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:567
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:571
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:575
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:579
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:583
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:589
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:593
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:599
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:603
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:607
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:611
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:615
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:623
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:627
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:631
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:635
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:639
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:643
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:647
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:653
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:657
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:663
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:667
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:673
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:677
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:683
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 108:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:687
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:693
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:697
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 111:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:703
		{
			yyVAL.token = Token{}
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:707
		{
			yyVAL.token = yyDollar[1].token
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:711
		{
			yyVAL.token = yyDollar[1].token
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:717
		{
			yyVAL.token = yyDollar[1].token
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:721
		{
			yyVAL.token = yyDollar[1].token
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:727
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:733
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
	case 118:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:756
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:760
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:764
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:768
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 122:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:772
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 123:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:776
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 124:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:780
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 125:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:784
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 126:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:788
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 127:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 128:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:796
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 129:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:800
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 130:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:804
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 131:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:808
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 132:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:812
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 133:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:816
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 134:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:820
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 135:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 136:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:828
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 137:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:834
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 138:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:838
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 139:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:842
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 140:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:846
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 141:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:850
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:856
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:860
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 144:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:864
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 145:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:870
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:874
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 147:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:880
		{
			yyVAL.expression = Option{}
		}
	case 148:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:884
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 149:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:888
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 150:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:894
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 151:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:898
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 152:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:904
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option), Over: yyDollar[5].token.Literal, AnalyticClause: yyDollar[7].expression.(AnalyticClause)}
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:910
		{
			yyVAL.expression = AnalyticClause{Partition: yyDollar[1].expression, OrderByClause: yyDollar[2].expression}
		}
	case 154:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:916
		{
			yyVAL.expression = nil
		}
	case 155:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:920
		{
			yyVAL.expression = Partition{PartitionBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Values: yyDollar[3].expressions}
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:926
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:930
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 158:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:934
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:940
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:944
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:950
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:954
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 163:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:958
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:962
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:966
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:970
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 167:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:976
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 168:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:980
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 169:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:984
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 170:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:988
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 171:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:992
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 172:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:998
		{
			yyVAL.expression = nil
		}
	case 173:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1002
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 174:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1006
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1012
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1016
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1020
		{
			yyVAL.expression = AllColumns{}
		}
	case 178:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1026
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 179:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1030
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 180:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1036
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 181:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1042
		{
			yyVAL.expression = nil
		}
	case 182:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1046
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 183:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1052
		{
			yyVAL.expression = nil
		}
	case 184:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1056
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1062
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 186:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1066
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 187:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1072
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 188:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1076
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1082
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 190:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1086
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1092
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 192:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1096
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 193:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1102
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 194:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1106
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1112
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 196:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1116
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 197:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1122
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 198:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1126
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 199:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1132
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Values: yyDollar[5].token.Literal, ValuesList: yyDollar[6].expressions}
		}
	case 200:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:1136
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Values: yyDollar[8].token.Literal, ValuesList: yyDollar[9].expressions}
		}
	case 201:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1140
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 202:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1144
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 203:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1150
		{
			yyVAL.expression = UpdateQuery{CommonTableClause: yyDollar[1].expression, Update: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, Set: yyDollar[4].token.Literal, SetList: yyDollar[5].expressions, FromClause: yyDollar[6].expression, WhereClause: yyDollar[7].expression}
		}
	case 204:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1156
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1162
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1166
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 207:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1172
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{CommonTableClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 208:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1177
		{
			from := FromClause{From: yyDollar[4].token.Literal, Tables: yyDollar[5].expressions}
			yyVAL.expression = DeleteQuery{CommonTableClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, FromClause: from, WhereClause: yyDollar[6].expression}
		}
	case 209:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1184
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 210:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1190
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 211:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1194
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1200
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 213:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1204
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1210
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 215:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1214
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 216:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1220
		{
			yyVAL.expression = nil
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1224
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1228
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 219:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1232
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 220:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1236
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 221:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1242
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 222:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1246
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 223:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1252
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 224:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1258
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 225:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1262
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 226:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1268
		{
			yyVAL.procexpr = nil
		}
	case 227:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1272
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 228:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1278
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 229:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1282
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 230:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1288
		{
			yyVAL.procexpr = nil
		}
	case 231:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1292
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 232:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1298
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 233:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1304
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 234:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1310
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 235:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1314
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 236:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1321
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 237:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1325
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 238:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1332
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 239:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1338
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 240:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1344
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 241:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1350
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 242:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1356
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 243:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1360
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 244:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1366
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 245:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1372
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 246:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1376
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 247:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1382
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 248:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1386
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 249:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1392
		{
			yyVAL.token = Token{}
		}
	case 250:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1396
		{
			yyVAL.token = yyDollar[1].token
		}
	case 251:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1402
		{
			yyVAL.token = Token{}
		}
	case 252:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1406
		{
			yyVAL.token = yyDollar[1].token
		}
	case 253:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1412
		{
			yyVAL.token = Token{}
		}
	case 254:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1416
		{
			yyVAL.token = yyDollar[1].token
		}
	case 255:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1422
		{
			yyVAL.token = Token{}
		}
	case 256:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1426
		{
			yyVAL.token = yyDollar[1].token
		}
	case 257:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1432
		{
			yyVAL.token = Token{}
		}
	case 258:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1436
		{
			yyVAL.token = yyDollar[1].token
		}
	case 259:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1440
		{
			yyVAL.token = yyDollar[1].token
		}
	case 260:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1444
		{
			yyVAL.token = yyDollar[1].token
		}
	case 261:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1450
		{
			yyVAL.token = Token{}
		}
	case 262:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1454
		{
			yyVAL.token = yyDollar[1].token
		}
	case 263:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1460
		{
			yyVAL.token = Token{}
		}
	case 264:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1464
		{
			yyVAL.token = yyDollar[1].token
		}
	case 265:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1471
		{
			yyVAL.token = yyDollar[1].token
		}
	case 266:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1475
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 267:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1481
		{
			yyVAL.token = Token{}
		}
	case 268:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1485
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
