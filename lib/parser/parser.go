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
const PARTITION = 57434
const OVER = 57435
const COMMIT = 57436
const ROLLBACK = 57437
const CONTINUE = 57438
const BREAK = 57439
const EXIT = 57440
const PRINT = 57441
const VAR = 57442
const COMPARISON_OP = 57443
const STRING_OP = 57444
const SUBSTITUTION_OP = 57445

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

//line parser.y:1395

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
	63, 235,
	67, 235,
	68, 235,
	-2, 249,
	-1, 132,
	47, 237,
	49, 241,
	-2, 173,
	-1, 167,
	57, 46,
	58, 46,
	59, 46,
	-2, 77,
	-1, 169,
	112, 131,
	-2, 233,
	-1, 171,
	78, 166,
	-2, 235,
	-1, 180,
	37, 131,
	91, 131,
	112, 131,
	-2, 233,
	-1, 197,
	63, 235,
	67, 235,
	68, 235,
	-2, 159,
	-1, 204,
	63, 235,
	67, 235,
	68, 235,
	-2, 61,
	-1, 208,
	63, 235,
	67, 235,
	68, 235,
	-2, 93,
	-1, 221,
	49, 241,
	-2, 237,
	-1, 237,
	63, 235,
	67, 235,
	68, 235,
	-2, 230,
	-1, 243,
	69, 0,
	101, 0,
	104, 0,
	-2, 102,
	-1, 244,
	69, 0,
	101, 0,
	104, 0,
	-2, 104,
	-1, 287,
	63, 235,
	67, 235,
	68, 235,
	-2, 51,
	-1, 294,
	112, 131,
	-2, 233,
	-1, 295,
	63, 235,
	67, 235,
	68, 235,
	-2, 64,
	-1, 337,
	69, 0,
	101, 0,
	104, 0,
	-2, 113,
	-1, 341,
	63, 235,
	67, 235,
	68, 235,
	-2, 171,
	-1, 379,
	63, 235,
	67, 235,
	68, 235,
	-2, 188,
	-1, 405,
	82, 168,
	-2, 235,
	-1, 409,
	112, 86,
	113, 86,
	-2, 46,
	-1, 417,
	63, 235,
	67, 235,
	68, 235,
	-2, 55,
	-1, 438,
	63, 235,
	67, 235,
	68, 235,
	-2, 197,
	-1, 445,
	78, 181,
	80, 181,
	82, 181,
	-2, 235,
	-1, 453,
	96, 18,
	97, 18,
	-2, 1,
	-1, 457,
	63, 235,
	67, 235,
	68, 235,
	-2, 157,
}

const yyPrivate = 57344

const yyLast = 1142

var yyAct = [...]int{

	80, 40, 392, 40, 53, 166, 44, 467, 319, 477,
	387, 46, 47, 48, 49, 50, 51, 52, 329, 305,
	43, 1, 313, 198, 205, 85, 38, 275, 38, 186,
	67, 68, 69, 427, 213, 194, 2, 94, 303, 92,
	111, 296, 400, 40, 393, 132, 262, 77, 222, 86,
	23, 37, 23, 220, 109, 346, 131, 64, 110, 155,
	45, 137, 114, 59, 16, 133, 56, 163, 162, 164,
	189, 142, 154, 437, 386, 413, 376, 119, 146, 147,
	148, 374, 136, 138, 57, 57, 61, 183, 167, 39,
	183, 308, 59, 163, 162, 164, 412, 424, 154, 176,
	299, 39, 292, 143, 152, 151, 128, 153, 157, 158,
	159, 160, 161, 488, 484, 39, 137, 45, 466, 449,
	185, 225, 448, 226, 227, 228, 223, 40, 447, 221,
	152, 151, 439, 153, 157, 158, 159, 160, 161, 436,
	137, 258, 259, 157, 158, 159, 160, 161, 419, 209,
	472, 40, 385, 375, 39, 342, 59, 174, 260, 59,
	39, 219, 210, 168, 294, 76, 201, 3, 168, 169,
	42, 241, 188, 268, 268, 240, 38, 458, 267, 163,
	162, 164, 40, 354, 154, 352, 224, 42, 191, 192,
	40, 350, 40, 40, 57, 217, 231, 232, 180, 42,
	23, 236, 42, 168, 239, 101, 103, 38, 309, 240,
	122, 268, 40, 267, 245, 184, 152, 151, 277, 153,
	157, 158, 159, 160, 161, 137, 264, 274, 259, 164,
	151, 23, 284, 157, 158, 159, 160, 161, 283, 40,
	119, 288, 144, 290, 291, 91, 328, 90, 268, 156,
	268, 268, 211, 455, 289, 122, 289, 289, 42, 483,
	238, 318, 145, 325, 172, 340, 322, 173, 167, 344,
	312, 268, 351, 353, 355, 311, 307, 40, 321, 356,
	486, 316, 475, 334, 330, 333, 454, 360, 361, 442,
	404, 363, 398, 215, 212, 357, 362, 341, 159, 160,
	161, 265, 38, 470, 358, 102, 332, 469, 471, 414,
	331, 137, 265, 277, 347, 470, 137, 332, 492, 485,
	464, 441, 298, 209, 372, 121, 23, 373, 367, 75,
	108, 104, 40, 113, 397, 382, 271, 377, 164, 370,
	270, 369, 407, 378, 182, 395, 190, 175, 179, 399,
	178, 409, 401, 409, 247, 409, 380, 38, 246, 248,
	117, 384, 97, 40, 273, 272, 250, 249, 116, 117,
	118, 371, 418, 314, 431, 268, 40, 306, 383, 122,
	381, 23, 137, 315, 137, 165, 310, 203, 38, 408,
	416, 410, 106, 411, 171, 433, 426, 177, 124, 277,
	324, 326, 301, 302, 490, 366, 327, 225, 268, 226,
	227, 228, 23, 423, 125, 365, 40, 193, 197, 286,
	54, 451, 204, 208, 396, 268, 394, 430, 263, 432,
	233, 234, 137, 421, 422, 122, 241, 63, 62, 235,
	462, 38, 237, 452, 40, 293, 149, 461, 463, 242,
	243, 244, 453, 88, 40, 251, 252, 253, 254, 255,
	256, 257, 468, 456, 465, 23, 460, 473, 120, 38,
	40, 474, 58, 58, 215, 55, 476, 459, 480, 38,
	70, 71, 72, 73, 74, 287, 59, 489, 40, 59,
	187, 306, 491, 23, 127, 38, 494, 478, 59, 495,
	115, 295, 139, 23, 230, 112, 277, 130, 60, 126,
	41, 66, 129, 38, 58, 493, 140, 141, 122, 23,
	122, 450, 122, 59, 277, 59, 100, 101, 103, 65,
	104, 105, 41, 93, 89, 266, 269, 23, 388, 389,
	390, 391, 306, 425, 335, 10, 337, 9, 8, 7,
	100, 101, 103, 6, 104, 105, 434, 435, 214, 5,
	4, 345, 225, 348, 226, 227, 228, 223, 170, 58,
	221, 82, 195, 304, 200, 196, 135, 359, 134, 200,
	482, 216, 58, 98, 218, 481, 95, 99, 229, 81,
	197, 106, 84, 58, 78, 96, 225, 208, 226, 227,
	228, 223, 428, 429, 221, 83, 79, 379, 420, 300,
	336, 107, 338, 339, 207, 106, 206, 202, 123, 364,
	285, 36, 261, 15, 278, 14, 13, 102, 199, 163,
	402, 164, 87, 349, 154, 12, 282, 11, 276, 0,
	0, 0, 0, 0, 41, 405, 39, 0, 18, 34,
	19, 102, 17, 0, 0, 0, 0, 0, 20, 0,
	0, 21, 0, 417, 0, 216, 152, 151, 0, 153,
	157, 158, 159, 160, 161, 0, 0, 0, 58, 0,
	0, 0, 0, 0, 317, 0, 320, 323, 216, 216,
	0, 0, 438, 0, 0, 0, 0, 0, 0, 0,
	0, 444, 0, 0, 445, 0, 0, 0, 279, 0,
	32, 0, 0, 304, 0, 304, 26, 304, 0, 30,
	27, 28, 29, 0, 0, 0, 457, 24, 25, 280,
	281, 33, 35, 22, 41, 0, 39, 304, 18, 34,
	19, 0, 17, 0, 42, 0, 200, 368, 20, 0,
	0, 21, 0, 200, 0, 0, 0, 0, 0, 0,
	0, 0, 216, 0, 58, 0, 0, 0, 479, 58,
	446, 0, 0, 0, 0, 0, 323, 0, 0, 216,
	0, 0, 0, 0, 0, 0, 0, 304, 0, 0,
	59, 100, 101, 103, 0, 104, 105, 41, 31, 39,
	32, 0, 0, 0, 0, 0, 26, 0, 0, 30,
	27, 28, 29, 0, 0, 0, 0, 24, 25, 0,
	0, 33, 35, 22, 0, 0, 0, 0, 0, 0,
	216, 0, 0, 0, 42, 58, 0, 58, 0, 0,
	320, 0, 0, 0, 216, 216, 0, 0, 98, 0,
	440, 0, 99, 0, 0, 0, 106, 0, 0, 0,
	96, 59, 100, 101, 103, 0, 104, 105, 41, 0,
	0, 0, 59, 100, 101, 103, 107, 104, 105, 41,
	0, 0, 0, 0, 0, 58, 0, 0, 0, 0,
	297, 323, 102, 0, 0, 0, 0, 87, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 163, 162,
	164, 0, 320, 154, 0, 0, 0, 298, 0, 98,
	0, 0, 0, 99, 0, 0, 0, 106, 0, 0,
	98, 96, 0, 0, 99, 0, 0, 0, 106, 0,
	0, 0, 96, 0, 0, 152, 151, 107, 153, 157,
	158, 159, 160, 161, 163, 162, 164, 0, 107, 154,
	0, 0, 0, 102, 343, 163, 162, 164, 87, 487,
	154, 0, 0, 0, 102, 0, 0, 0, 0, 87,
	443, 0, 0, 0, 0, 163, 162, 164, 0, 0,
	154, 152, 151, 0, 153, 157, 158, 159, 160, 161,
	415, 0, 152, 151, 0, 153, 157, 158, 159, 160,
	161, 0, 0, 0, 0, 0, 163, 162, 164, 0,
	0, 154, 152, 151, 0, 153, 157, 158, 159, 160,
	161, 406, 163, 162, 164, 0, 0, 154, 0, 0,
	0, 163, 162, 164, 0, 0, 154, 0, 0, 181,
	163, 162, 164, 152, 151, 154, 153, 157, 158, 159,
	160, 161, 403, 162, 164, 150, 0, 154, 0, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 152, 151,
	0, 153, 157, 158, 159, 160, 161, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 0, 0, 164, 152,
	151, 154, 153, 157, 158, 159, 160, 161, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 152, 151, 0, 153, 157, 158, 159,
	160, 161,
}
var yyPact = [...]int{

	723, -1000, 723, -54, -54, -54, -54, -54, -54, -54,
	-54, -1000, -1000, -1000, -1000, -1000, 383, 455, 519, 494,
	409, 408, 500, -54, -54, -54, 519, 519, 519, 519,
	519, 868, 868, -54, 493, 868, 486, 311, 137, 253,
	-1000, -1000, 147, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 355, 374, 519, 478, -7, 485, -1000,
	59, 488, 519, 519, -54, -10, 139, -1000, -1000, -1000,
	178, -54, -54, -54, 426, 986, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 137, -1000, 786, 58, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 868, 163, 91, 868,
	-1000, -1000, 199, -1000, -1000, -1000, -1000, 87, 968, 281,
	-26, -1000, 111, 3, 472, 59, 286, 286, 286, 868,
	521, -1000, 54, 343, 868, 868, 141, 519, 519, -1000,
	519, 472, 73, -1000, 482, -1000, -1000, -1000, -1000, 59,
	86, 404, -1000, 500, 868, 175, -1000, -1000, -1000, 499,
	723, 868, 868, 868, 272, 291, 306, 868, 868, 868,
	868, 868, 868, 868, -1000, 29, 46, -1000, 519, 253,
	223, 977, 67, 67, 273, 304, -1000, 1032, -1000, -1000,
	253, 633, 519, 499, 545, -1000, 381, 868, -1000, 147,
	-1000, 147, 147, 977, -1000, -11, 423, 977, -1000, -1000,
	53, -1000, -1000, 868, 844, -1000, -13, 361, 977, -1000,
	67, 88, -1000, 486, -22, 104, 93, -1000, -1000, -1000,
	339, 359, 324, 336, 59, -1000, -1000, -1000, -1000, -1000,
	519, 472, 519, 155, 152, 519, -1000, 977, 147, -54,
	-23, 230, 38, 128, 128, 322, 868, 67, 868, 67,
	67, 191, 191, -1000, -1000, -1000, 565, 1032, -1000, 868,
	-1000, -1000, 43, 857, 234, 868, -1000, 786, -1000, -1000,
	67, 80, 74, 72, 383, 213, 633, -1000, -1000, 868,
	-54, -54, 215, -1000, -54, 376, 365, 977, 302, -1000,
	-1000, 302, 521, 519, 253, 977, -1000, 249, 326, 868,
	256, -1000, -1000, -1000, -32, 41, -37, 472, 519, 868,
	59, 333, 324, 331, -1000, 59, -1000, -1000, -1000, 40,
	-39, 508, 519, 392, -1000, 519, 388, -54, -1000, 210,
	230, 723, 868, -1000, -1000, 998, -1000, 128, -1000, -1000,
	-1000, 115, -1000, -1000, -1000, 208, 223, 868, 952, 278,
	102, -1000, 102, -1000, 102, -1000, -16, 232, -1000, 921,
	-1000, -1000, 633, -1000, -1000, 868, 868, -1000, -1000, 36,
	-1000, -1000, -1000, 403, 67, 76, 519, -1000, -1000, 977,
	548, 59, 327, 59, 514, -1000, 519, -1000, -1000, -1000,
	519, 519, 27, -40, 868, 20, 519, -1000, 246, 207,
	241, -1000, 901, 868, -1000, 977, 868, 67, 16, -1000,
	10, 7, -1000, 516, -54, 633, 204, 977, -1000, 160,
	-1000, -1000, -1000, -1000, 67, -1000, -1000, -1000, 868, 66,
	514, 59, 548, -1000, -1000, -1000, 508, 519, 977, -1000,
	-1000, -54, 245, 723, 1032, 977, -1000, -1000, -1000, -1000,
	6, -1000, 227, 723, 231, 39, -1000, 977, 519, 514,
	-1000, -1000, -1000, -1000, -54, -1000, -1000, 200, 227, 633,
	868, -54, 167, 2, -1000, 244, 198, 239, -1000, 890,
	-1000, 1, 383, 364, -1000, -54, 243, 633, -1000, -1000,
	868, -1000, -54, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 20, 27, 36, 638, 637, 635, 626, 625, 624,
	623, 167, 64, 51, 621, 62, 29, 620, 619, 4,
	618, 41, 617, 47, 165, 297, 362, 38, 24, 616,
	614, 609, 608, 0, 606, 605, 594, 592, 589, 46,
	586, 23, 585, 580, 65, 578, 45, 576, 33, 575,
	572, 571, 568, 561, 19, 5, 56, 66, 8, 35,
	55, 560, 559, 558, 34, 553, 549, 548, 44, 2,
	10, 547, 545, 42, 18, 9, 7, 453, 534, 247,
	245, 39, 533, 37, 25, 54, 49, 529, 57, 428,
	59, 53, 22, 48, 70, 249, 6,
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
	28, 29, 29, 30, 30, 31, 31, 31, 32, 32,
	33, 34, 35, 35, 35, 35, 35, 35, 35, 35,
	35, 35, 35, 35, 35, 35, 35, 35, 35, 35,
	35, 36, 36, 36, 36, 36, 37, 37, 37, 38,
	38, 39, 39, 39, 40, 40, 41, 42, 43, 43,
	44, 44, 44, 45, 45, 46, 46, 46, 46, 46,
	46, 47, 47, 47, 47, 47, 48, 48, 48, 49,
	49, 49, 50, 50, 51, 52, 52, 53, 53, 54,
	54, 55, 55, 56, 56, 57, 57, 58, 58, 59,
	59, 60, 60, 61, 61, 61, 61, 62, 63, 64,
	64, 65, 65, 66, 67, 67, 68, 68, 69, 69,
	70, 70, 70, 70, 70, 71, 71, 72, 73, 73,
	74, 74, 75, 75, 76, 76, 77, 78, 79, 79,
	80, 80, 81, 82, 83, 84, 85, 85, 86, 87,
	87, 88, 88, 89, 89, 90, 90, 91, 91, 92,
	92, 93, 93, 93, 93, 94, 94, 95, 95, 96,
	96,
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
	3, 2, 4, 1, 1, 0, 1, 1, 1, 1,
	3, 3, 3, 3, 3, 3, 4, 4, 6, 6,
	4, 6, 4, 4, 4, 6, 4, 4, 6, 4,
	2, 3, 3, 3, 3, 3, 3, 3, 2, 4,
	1, 0, 2, 2, 5, 7, 8, 2, 0, 3,
	1, 2, 3, 1, 1, 1, 1, 2, 3, 1,
	1, 5, 5, 6, 6, 4, 0, 2, 4, 1,
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

	-1000, -1, -3, -11, -61, -62, -65, -66, -67, -71,
	-72, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 100, -86, 94, 95, 83, 87, 88, 89,
	86, 75, 77, 98, 16, 99, -14, -13, -84, 13,
	-33, 11, 111, -1, -96, 114, -96, -96, -96, -96,
	-96, -96, -96, -19, 37, 20, -57, -44, -77, 4,
	14, -57, 29, 29, -88, -87, 11, -96, -96, -96,
	-77, -77, -77, -77, -77, -25, -24, -23, -36, -34,
	-33, -38, -51, -35, -37, -84, -86, 111, -77, -78,
	-79, -80, -81, -82, -83, -40, 74, -26, 62, 66,
	5, 6, 106, 7, 9, 10, 70, 90, -25, -85,
	-84, -96, 12, -25, -15, 14, 57, 58, 59, 103,
	-89, 72, -11, -20, 43, 40, -77, 16, 113, -77,
	22, -56, -46, -44, -45, -47, 23, -33, 24, 14,
	-77, -77, -96, 113, 103, 84, -96, -96, -96, 20,
	79, 102, 101, 104, 69, -90, -95, 105, 106, 107,
	108, 109, 65, 64, 66, -25, -55, -33, 110, 111,
	-52, -25, 101, 104, -90, -95, -33, -25, -79, -80,
	111, 81, 63, 113, 104, -96, -16, 18, -56, -94,
	60, -94, -94, -25, -59, -50, -49, -25, -41, 107,
	-77, 112, -22, 44, -25, -28, -29, -30, -25, -41,
	21, 111, -11, -64, -63, -24, -77, -57, -77, -16,
	-91, 56, -93, 53, 113, 48, 50, 51, 52, -77,
	22, -56, 111, 26, 27, 35, -88, -25, 85, -85,
	-84, -1, -25, -25, -25, -90, 67, 63, 68, 61,
	60, -25, -25, -25, -25, -25, -25, -25, 112, 113,
	112, -77, -39, -89, -60, 78, -26, 111, -33, -26,
	67, 63, 61, 60, -39, -2, -4, -3, -9, 75,
	96, 97, -77, -85, -23, -17, 38, -25, -13, -12,
	-13, -13, 113, 22, 111, -25, -21, 46, 73, 113,
	-31, 41, 42, -27, -26, -54, -24, -15, 113, 104,
	47, -91, -93, -92, 49, 47, -56, -77, -16, -58,
	-77, -68, 111, -77, -24, 111, -24, -11, -96, -74,
	-73, 80, 76, -81, -83, -25, -26, -25, -26, -26,
	-55, -25, 112, 107, -55, -53, -60, 80, -25, -26,
	111, -33, 111, -33, 111, -33, -19, 82, -2, -25,
	-96, -96, 81, -96, -18, 39, 40, -59, -77, -39,
	-21, 45, -28, 71, 113, 112, 113, -16, -64, -25,
	-46, 47, -92, 47, -46, 112, 113, -70, 30, 31,
	32, 33, -69, -68, 34, -54, 36, -96, 82, -74,
	-73, -1, -25, 64, 82, -25, 79, 64, -27, -33,
	-27, -27, 112, 91, 77, 79, -2, -25, -55, 112,
	-32, 30, 31, -27, 21, -11, -54, -48, 54, 55,
	-46, 47, -46, -58, -24, -24, 112, 113, -25, 112,
	-77, 75, 82, 79, -25, -25, -26, 112, 112, 112,
	5, -96, -2, -3, 82, 93, -27, -25, 111, -46,
	-48, -70, -69, -96, 75, -1, 112, -76, -75, 80,
	76, 77, 111, -58, -96, 82, -76, -75, -2, -25,
	-96, -42, -43, 92, 112, 75, 82, 79, 112, -19,
	40, -96, 75, -2, -55, -96,
}
var yyDef = [...]int{

	1, -2, 1, 249, 249, 249, 249, 249, 249, 249,
	249, 13, 14, 15, 16, 17, -2, 0, 0, 0,
	0, 0, 0, 249, 249, 249, 0, 0, 0, 0,
	0, 0, 0, 249, 0, 0, 48, 0, 0, 233,
	46, 225, 0, 2, 5, 250, 6, 7, 8, 9,
	10, 11, 12, 58, 0, 0, 0, 175, 140, 216,
	0, 0, 0, 0, 249, 231, 229, 21, 22, 23,
	0, 249, 249, 249, 0, 235, 73, 74, 75, 76,
	77, 78, 79, 80, 81, 82, 83, 0, 71, 65,
	66, 67, 68, 69, 70, 130, 165, 235, 0, 0,
	217, 218, 0, 220, 222, 223, 224, 0, 235, 0,
	82, 33, 0, -2, 50, 0, 245, 245, 245, 0,
	0, 234, 0, 63, 0, 0, 0, 0, 0, 141,
	0, 50, -2, 145, 146, 149, 150, 143, 144, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	1, 0, 247, 248, 235, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 236, 235, 0, -2, 0, -2,
	0, -2, 247, 248, 0, 0, 120, 128, 219, 221,
	-2, 3, 0, 0, 0, 39, 52, 0, 49, 0,
	246, 0, 0, 228, 47, 179, 162, -2, 160, 161,
	71, 100, 40, 0, -2, 57, 89, 95, -2, 94,
	0, 0, 185, 48, 189, 0, 71, 176, 142, 191,
	0, -2, 239, 0, 0, 238, 242, 243, 244, 147,
	0, 50, 0, 0, 0, 0, 232, -2, 0, 249,
	226, 210, 101, -2, -2, 0, 0, 0, 0, 0,
	0, 121, 122, 123, 124, 125, 126, 127, 84, 0,
	85, 72, 0, 0, 167, 0, 103, 0, 86, 105,
	0, 0, 0, 0, 56, 0, 3, 18, 19, 0,
	249, 249, 0, 227, 249, 54, 0, -2, 42, 45,
	43, 44, 0, 0, -2, -2, 59, 61, 0, 0,
	91, 96, 97, 183, 87, 0, 169, 50, 0, 0,
	0, 0, 239, 0, 240, 0, 174, 148, 192, 0,
	177, 200, 0, 196, 205, 0, 0, 249, 28, 0,
	210, 1, 0, 106, 107, 235, 110, -2, 114, 117,
	172, -2, 129, 132, 133, 0, 182, 0, 235, 0,
	0, 112, 0, 116, 0, 119, 0, 0, 4, 235,
	36, 37, 3, 38, 41, 0, 0, 180, 163, 0,
	60, 62, 90, 0, 0, 0, 0, 187, 190, -2,
	156, 0, 0, 0, 155, 193, 0, 194, 201, 202,
	0, 0, 0, 198, 0, 0, 0, 24, 0, 0,
	209, 211, 235, 0, 164, -2, 0, 0, 0, -2,
	0, 0, 134, 0, 249, 1, 0, -2, 53, 129,
	92, 98, 99, 88, 0, 186, 170, 151, 0, 0,
	152, 0, 156, 178, 203, 204, 200, 0, -2, 206,
	207, 249, 0, 1, 108, -2, 109, 111, 115, 118,
	0, 31, 214, -2, 0, 0, 184, -2, 0, 154,
	153, 195, 199, 29, 249, 208, 135, 0, 214, 3,
	0, 249, 138, 0, 30, 0, 0, 213, 215, 235,
	32, 0, 56, 0, 158, 249, 0, 3, 136, 137,
	0, 34, 249, 212, 139, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 109, 3, 3,
	111, 112, 107, 105, 113, 106, 110, 108, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 114,
	3, 104,
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
	102, 103,
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
		//line parser.y:154
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:159
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:166
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:171
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:178
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:182
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:186
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:190
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:194
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:198
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:202
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:206
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:210
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:214
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:218
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:222
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:226
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:232
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:236
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:242
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:246
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:252
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:256
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:262
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:266
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:270
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:274
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:278
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:284
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:288
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:292
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:296
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:300
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:306
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:310
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:314
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:318
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:324
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:328
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:334
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
		//line parser.y:345
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
		//line parser.y:355
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
		//line parser.y:364
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
		//line parser.y:373
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
		//line parser.y:384
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:388
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:394
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:400
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:404
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:410
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:414
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:424
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:430
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:434
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:440
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:444
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 58:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:450
		{
			yyVAL.expression = nil
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:454
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:458
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:464
		{
			yyVAL.expression = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:468
		{
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:474
		{
			yyVAL.expression = nil
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:478
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:484
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:488
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:492
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:496
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:500
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:504
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:510
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:514
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:528
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 78:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:544
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:548
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:560
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:564
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:570
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:574
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:580
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:584
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:590
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:594
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:600
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 92:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:604
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:610
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:614
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 95:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:620
		{
			yyVAL.token = Token{}
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:624
		{
			yyVAL.token = yyDollar[1].token
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:628
		{
			yyVAL.token = yyDollar[1].token
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:634
		{
			yyVAL.token = yyDollar[1].token
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:638
		{
			yyVAL.token = yyDollar[1].token
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:644
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:650
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
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:673
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:677
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:681
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:685
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 106:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:689
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 107:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:693
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 108:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:697
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 109:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:701
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 110:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:705
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 111:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:709
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 112:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:713
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 113:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:717
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 114:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:721
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 115:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:725
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 116:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:729
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 117:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:733
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 118:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:737
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 119:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:741
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 120:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:745
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:751
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:755
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:759
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:763
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:767
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:773
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:777
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 128:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:781
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 129:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:787
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:791
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 131:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:797
		{
			yyVAL.expression = Option{}
		}
	case 132:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:801
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 133:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:805
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 134:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:811
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 135:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:815
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 136:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:821
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option), Over: yyDollar[5].token.Literal, AnalyticClause: yyDollar[7].expression.(AnalyticClause)}
		}
	case 137:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:827
		{
			yyVAL.expression = AnalyticClause{Partition: yyDollar[1].expression, OrderByClause: yyDollar[2].expression}
		}
	case 138:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:833
		{
			yyVAL.expression = nil
		}
	case 139:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:837
		{
			yyVAL.expression = Partition{PartitionBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Values: yyDollar[3].expressions}
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:843
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 141:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:847
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:851
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:861
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 145:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:867
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:871
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 147:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:875
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 148:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:879
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 149:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:883
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:887
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 151:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:893
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 152:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:897
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 153:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:901
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 154:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:905
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 155:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:909
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 156:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:915
		{
			yyVAL.expression = nil
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:919
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 158:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:923
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:929
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:933
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:937
		{
			yyVAL.expression = AllColumns{}
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:943
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 163:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:947
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 164:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:953
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 165:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:959
		{
			yyVAL.expression = nil
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:963
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 167:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:969
		{
			yyVAL.expression = nil
		}
	case 168:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:973
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:979
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 170:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:983
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:989
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 172:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:993
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 173:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:999
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 174:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1003
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1009
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 176:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1013
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1019
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 178:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1023
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1029
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 180:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1033
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 181:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1039
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 182:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1043
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 183:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1049
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 184:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1053
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 185:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1057
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 186:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1061
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 187:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1067
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 188:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1073
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1079
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 190:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1083
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 191:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1089
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 192:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1094
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 193:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1101
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 194:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1107
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 195:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1111
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 196:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1117
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 197:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1121
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 198:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1127
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 199:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1131
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 200:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1137
		{
			yyVAL.expression = nil
		}
	case 201:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1141
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1145
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 203:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1149
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 204:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1153
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 205:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1159
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 206:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1163
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 207:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1169
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 208:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1175
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 209:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1179
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 210:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1185
		{
			yyVAL.procexpr = nil
		}
	case 211:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1189
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 212:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1195
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 213:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1199
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 214:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1205
		{
			yyVAL.procexpr = nil
		}
	case 215:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1209
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1215
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1221
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1227
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 219:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1231
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 220:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1238
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 221:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1242
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1249
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 223:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1255
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1261
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 225:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1267
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1273
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 227:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1277
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 228:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1283
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 229:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1289
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 230:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1293
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 231:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1299
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 232:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1303
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 233:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1309
		{
			yyVAL.token = Token{}
		}
	case 234:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1313
		{
			yyVAL.token = yyDollar[1].token
		}
	case 235:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1319
		{
			yyVAL.token = Token{}
		}
	case 236:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1323
		{
			yyVAL.token = yyDollar[1].token
		}
	case 237:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1329
		{
			yyVAL.token = Token{}
		}
	case 238:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1333
		{
			yyVAL.token = yyDollar[1].token
		}
	case 239:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1339
		{
			yyVAL.token = Token{}
		}
	case 240:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1343
		{
			yyVAL.token = yyDollar[1].token
		}
	case 241:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1349
		{
			yyVAL.token = Token{}
		}
	case 242:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1353
		{
			yyVAL.token = yyDollar[1].token
		}
	case 243:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1357
		{
			yyVAL.token = yyDollar[1].token
		}
	case 244:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1361
		{
			yyVAL.token = yyDollar[1].token
		}
	case 245:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1367
		{
			yyVAL.token = Token{}
		}
	case 246:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1371
		{
			yyVAL.token = yyDollar[1].token
		}
	case 247:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1377
		{
			yyVAL.token = yyDollar[1].token
		}
	case 248:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1381
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 249:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1387
		{
			yyVAL.token = Token{}
		}
	case 250:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1391
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
