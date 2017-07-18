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
const SEPARATOR = 57438
const PARTITION = 57439
const OVER = 57440
const COMMIT = 57441
const ROLLBACK = 57442
const CONTINUE = 57443
const BREAK = 57444
const EXIT = 57445
const PRINT = 57446
const PRINTF = 57447
const SOURCE = 57448
const VAR = 57449
const COMPARISON_OP = 57450
const STRING_OP = 57451
const SUBSTITUTION_OP = 57452

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
	"SEPARATOR",
	"PARTITION",
	"OVER",
	"COMMIT",
	"ROLLBACK",
	"CONTINUE",
	"BREAK",
	"EXIT",
	"PRINT",
	"PRINTF",
	"SOURCE",
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

//line parser.y:1545

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
	-2, 80,
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	1, 1,
	77, 1,
	81, 1,
	83, 1,
	-2, 80,
	-1, 49,
	58, 60,
	59, 60,
	60, 60,
	-2, 71,
	-1, 118,
	64, 263,
	68, 263,
	69, 263,
	-2, 279,
	-1, 120,
	64, 263,
	68, 263,
	69, 263,
	-2, 199,
	-1, 155,
	4, 38,
	-2, 263,
	-1, 156,
	4, 39,
	-2, 263,
	-1, 157,
	77, 1,
	81, 1,
	83, 1,
	-2, 80,
	-1, 177,
	79, 194,
	-2, 263,
	-1, 187,
	83, 3,
	-2, 80,
	-1, 207,
	48, 265,
	50, 269,
	-2, 201,
	-1, 225,
	64, 263,
	68, 263,
	69, 263,
	-2, 187,
	-1, 235,
	64, 263,
	68, 263,
	69, 263,
	-2, 258,
	-1, 242,
	70, 0,
	108, 0,
	111, 0,
	-2, 125,
	-1, 243,
	70, 0,
	108, 0,
	111, 0,
	-2, 127,
	-1, 263,
	96, 71,
	118, 199,
	-2, 263,
	-1, 278,
	77, 3,
	81, 3,
	83, 3,
	-2, 80,
	-1, 294,
	64, 263,
	68, 263,
	69, 263,
	-2, 76,
	-1, 298,
	64, 263,
	68, 263,
	69, 263,
	-2, 116,
	-1, 311,
	50, 269,
	-2, 265,
	-1, 324,
	64, 263,
	68, 263,
	69, 263,
	-2, 66,
	-1, 346,
	83, 1,
	-2, 80,
	-1, 352,
	70, 0,
	108, 0,
	111, 0,
	-2, 136,
	-1, 360,
	96, 71,
	118, 158,
	-2, 263,
	-1, 378,
	83, 3,
	-2, 80,
	-1, 382,
	64, 263,
	68, 263,
	69, 263,
	-2, 79,
	-1, 440,
	83, 196,
	-2, 263,
	-1, 449,
	77, 1,
	81, 1,
	83, 1,
	-2, 80,
	-1, 462,
	64, 263,
	68, 263,
	69, 263,
	-2, 216,
	-1, 468,
	64, 263,
	68, 263,
	69, 263,
	-2, 70,
	-1, 477,
	64, 263,
	68, 263,
	69, 263,
	-2, 225,
	-1, 483,
	77, 1,
	81, 1,
	83, 1,
	-2, 80,
	-1, 489,
	79, 209,
	81, 209,
	83, 209,
	-2, 263,
	-1, 497,
	77, 1,
	81, 1,
	83, 1,
	-2, 19,
	-1, 528,
	83, 3,
	-2, 80,
	-1, 533,
	64, 263,
	68, 263,
	69, 263,
	-2, 185,
	-1, 559,
	77, 3,
	81, 3,
	83, 3,
	-2, 80,
}

const yyPrivate = 57344

const yyLast = 1282

var yyAct = [...]int{

	40, 90, 21, 93, 21, 42, 43, 44, 45, 46,
	47, 48, 89, 37, 537, 37, 332, 173, 127, 526,
	547, 277, 63, 64, 65, 508, 420, 207, 66, 68,
	69, 70, 116, 415, 383, 400, 143, 17, 392, 17,
	344, 303, 312, 115, 299, 39, 1, 135, 215, 99,
	2, 222, 119, 206, 162, 132, 132, 310, 97, 114,
	421, 146, 144, 145, 122, 3, 295, 80, 431, 151,
	152, 153, 208, 130, 362, 154, 326, 60, 41, 174,
	331, 476, 218, 174, 315, 390, 316, 317, 318, 313,
	175, 54, 311, 174, 49, 459, 126, 457, 414, 170,
	169, 171, 395, 386, 161, 189, 189, 329, 203, 184,
	67, 106, 107, 109, 195, 110, 111, 38, 147, 191,
	192, 300, 194, 142, 131, 131, 134, 563, 196, 561,
	560, 553, 544, 201, 543, 523, 204, 132, 522, 499,
	132, 493, 159, 158, 228, 160, 164, 165, 166, 167,
	168, 492, 67, 67, 314, 41, 491, 180, 142, 21,
	158, 478, 475, 164, 165, 166, 167, 168, 471, 103,
	37, 211, 213, 104, 36, 79, 458, 112, 259, 452,
	36, 101, 106, 107, 109, 226, 110, 111, 36, 21,
	217, 105, 284, 261, 17, 426, 120, 413, 171, 356,
	37, 161, 286, 240, 228, 355, 306, 132, 67, 308,
	55, 288, 290, 319, 258, 229, 244, 301, 132, 108,
	49, 220, 221, 92, 17, 234, 78, 113, 238, 540,
	118, 269, 534, 269, 333, 336, 306, 306, 279, 159,
	158, 333, 160, 164, 165, 166, 167, 168, 112, 285,
	531, 266, 286, 514, 341, 309, 380, 370, 287, 237,
	289, 166, 167, 168, 96, 57, 338, 302, 321, 368,
	366, 230, 57, 155, 156, 171, 131, 307, 107, 109,
	21, 396, 358, 376, 377, 190, 102, 515, 379, 172,
	108, 37, 334, 126, 349, 333, 148, 470, 177, 343,
	374, 183, 340, 348, 342, 306, 539, 291, 381, 345,
	325, 95, 327, 328, 57, 17, 178, 488, 132, 179,
	436, 335, 494, 197, 404, 480, 447, 236, 558, 279,
	163, 372, 150, 228, 410, 142, 545, 225, 498, 336,
	482, 425, 306, 427, 428, 235, 439, 429, 21, 412,
	446, 394, 373, 378, 399, 241, 242, 243, 267, 37,
	363, 250, 251, 252, 253, 254, 255, 256, 403, 398,
	405, 529, 263, 186, 226, 528, 267, 423, 305, 438,
	21, 409, 530, 17, 347, 108, 430, 149, 346, 448,
	228, 37, 432, 529, 347, 565, 294, 298, 557, 306,
	450, 132, 520, 481, 505, 36, 132, 385, 337, 339,
	141, 456, 276, 324, 171, 17, 442, 110, 333, 453,
	185, 188, 306, 306, 219, 463, 469, 138, 479, 279,
	467, 472, 142, 181, 142, 465, 142, 461, 401, 228,
	512, 83, 350, 460, 352, 466, 451, 275, 274, 495,
	464, 21, 443, 455, 444, 402, 445, 36, 397, 56,
	454, 360, 37, 306, 364, 268, 271, 293, 132, 555,
	132, 496, 164, 165, 166, 167, 168, 393, 375, 112,
	336, 199, 519, 246, 407, 21, 17, 245, 247, 273,
	382, 228, 511, 272, 513, 240, 37, 487, 507, 21,
	497, 435, 485, 517, 249, 248, 408, 74, 75, 516,
	37, 137, 138, 139, 393, 200, 132, 527, 323, 541,
	17, 542, 388, 389, 506, 128, 225, 424, 263, 521,
	21, 550, 422, 351, 17, 353, 354, 59, 333, 536,
	535, 37, 58, 504, 433, 182, 518, 546, 502, 503,
	548, 552, 125, 500, 524, 556, 67, 554, 564, 365,
	440, 21, 67, 330, 239, 17, 567, 72, 73, 76,
	77, 305, 37, 562, 320, 212, 202, 129, 212, 279,
	205, 566, 55, 298, 51, 56, 52, 391, 50, 216,
	136, 532, 214, 462, 473, 474, 17, 551, 315, 117,
	316, 317, 318, 313, 468, 315, 311, 316, 317, 318,
	279, 67, 106, 107, 109, 38, 110, 111, 38, 477,
	270, 270, 315, 67, 316, 317, 318, 313, 509, 510,
	311, 484, 298, 133, 62, 393, 231, 232, 489, 525,
	486, 67, 106, 107, 109, 233, 110, 111, 38, 416,
	417, 418, 419, 391, 121, 391, 212, 391, 67, 124,
	56, 140, 56, 56, 61, 98, 94, 10, 9, 8,
	103, 7, 6, 304, 104, 5, 4, 361, 112, 176,
	264, 86, 101, 223, 298, 224, 57, 210, 270, 209,
	270, 270, 105, 538, 100, 262, 85, 84, 88, 81,
	103, 87, 82, 357, 104, 501, 533, 387, 112, 297,
	264, 296, 101, 123, 270, 367, 369, 371, 292, 198,
	108, 265, 105, 406, 92, 411, 549, 322, 53, 490,
	170, 169, 171, 16, 280, 161, 15, 91, 71, 14,
	13, 12, 270, 11, 391, 278, 0, 170, 169, 171,
	108, 265, 161, 0, 92, 260, 212, 67, 106, 107,
	109, 0, 110, 111, 38, 0, 0, 0, 0, 0,
	0, 0, 0, 159, 158, 0, 160, 164, 165, 166,
	167, 168, 0, 0, 193, 0, 0, 0, 0, 0,
	159, 158, 391, 160, 164, 165, 166, 167, 168, 0,
	257, 193, 170, 169, 171, 0, 0, 161, 270, 0,
	270, 0, 270, 0, 0, 0, 103, 0, 0, 0,
	104, 0, 0, 0, 112, 0, 0, 0, 101, 67,
	106, 107, 109, 0, 110, 111, 38, 0, 105, 212,
	0, 0, 0, 0, 212, 159, 158, 0, 160, 164,
	165, 166, 167, 168, 0, 0, 193, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 108, 227, 0, 0,
	92, 0, 0, 0, 67, 106, 107, 109, 0, 110,
	111, 38, 0, 0, 270, 0, 0, 0, 103, 0,
	0, 0, 104, 0, 0, 384, 112, 0, 0, 270,
	101, 0, 0, 0, 0, 0, 212, 0, 212, 0,
	105, 0, 0, 170, 169, 171, 0, 0, 161, 0,
	0, 0, 385, 0, 0, 170, 169, 171, 0, 0,
	161, 0, 0, 103, 0, 0, 0, 104, 108, 359,
	559, 112, 92, 0, 0, 101, 0, 270, 0, 0,
	0, 0, 0, 0, 212, 105, 159, 158, 0, 160,
	164, 165, 166, 167, 168, 170, 169, 171, 159, 158,
	161, 160, 164, 165, 166, 167, 168, 170, 169, 171,
	483, 0, 161, 108, 0, 0, 0, 92, 170, 169,
	171, 0, 449, 161, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 441, 0, 0, 0, 0, 159, 158,
	0, 160, 164, 165, 166, 167, 168, 437, 0, 0,
	159, 158, 0, 160, 164, 165, 166, 167, 168, 0,
	0, 159, 158, 0, 160, 164, 165, 166, 167, 168,
	0, 0, 0, 0, 170, 169, 171, 0, 0, 161,
	0, 38, 0, 170, 169, 171, 32, 0, 161, 0,
	0, 0, 170, 169, 171, 0, 18, 161, 0, 19,
	187, 0, 0, 0, 0, 0, 0, 157, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 159, 158, 0,
	160, 164, 165, 166, 167, 168, 159, 158, 0, 160,
	164, 165, 166, 167, 168, 159, 158, 0, 160, 164,
	165, 166, 167, 168, 36, 0, 281, 0, 30, 0,
	170, 169, 171, 0, 24, 161, 0, 28, 25, 26,
	27, 434, 169, 171, 0, 0, 161, 0, 0, 22,
	23, 282, 283, 31, 33, 34, 35, 20, 170, 0,
	171, 0, 0, 161, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 159, 158, 0, 160, 164, 165, 166,
	167, 168, 0, 0, 159, 158, 0, 160, 164, 165,
	166, 167, 168, 0, 0, 38, 0, 0, 0, 0,
	32, 159, 158, 0, 160, 164, 165, 166, 167, 168,
	18, 0, 0, 19, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 36, 0,
	29, 0, 30, 0, 0, 0, 0, 0, 24, 0,
	0, 28, 25, 26, 27, 0, 0, 0, 0, 0,
	0, 0, 0, 22, 23, 0, 0, 31, 33, 34,
	35, 20,
}
var yyPact = [...]int{

	1174, -1000, 1174, -43, -43, -43, -43, -43, -43, -43,
	-43, -1000, -1000, -1000, -1000, -1000, -1000, 569, 512, 507,
	623, -43, -43, -43, 654, 654, 654, 654, 476, 870,
	870, -43, 587, 870, 870, 649, 527, 183, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 487,
	557, 654, 619, 576, 453, 337, -1000, 331, 654, 654,
	-43, -1, 186, -1000, -1000, -1000, 302, -1000, -43, -43,
	-43, 654, -1000, -1000, -1000, -1000, 870, 870, 997, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 183,
	-1000, -1000, 106, -27, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 870, 208, 155, 870, 654, -1000, -1000, 272, -1000,
	-1000, -1000, -1000, 988, 357, -14, -1000, 174, 34, -43,
	737, -43, -1000, -5, 654, -1000, 870, 437, 474, 654,
	560, -11, 558, 148, 578, 571, 148, 363, 363, 363,
	753, -1000, 97, 197, 154, 609, -1000, 623, 870, 241,
	142, -1000, -1000, -1000, 544, 1055, 1055, 1174, 870, 870,
	870, 347, 419, 443, 870, 870, 870, 870, 870, 870,
	870, -1000, 682, 96, 654, 637, 297, 1055, 116, 116,
	425, 386, -1000, 131, 342, -1000, -1000, 1040, 654, 604,
	177, -1000, -1000, 870, -1000, 527, 190, 1055, 422, 870,
	870, 100, 654, 654, -1000, 654, 571, 35, -1000, 552,
	-1000, -1000, -1000, -1000, 148, 479, 870, -1000, 197, -1000,
	197, 197, -1000, -12, 541, 1055, -1000, -1000, -37, -1000,
	654, 204, 149, 654, -1000, 1055, 331, 654, 331, 604,
	307, 360, 51, 51, 408, 870, 116, 870, 116, 116,
	147, 147, -1000, -1000, -1000, 1083, 131, -1000, -1000, -1000,
	-1000, 87, 81, 665, 825, -1000, 279, 870, -1000, 106,
	-1000, -1000, 116, 153, 152, 140, 347, 269, 1040, -1000,
	-1000, 870, -43, -43, 271, -1000, -13, -43, -1000, -1000,
	139, 654, -1000, 870, 848, -1000, -16, 480, 1055, -1000,
	116, 654, -1000, 576, -17, 170, -41, -1000, -1000, -1000,
	410, 556, 388, 407, 148, -1000, -1000, -1000, -1000, -1000,
	654, 571, 444, 465, 1055, 368, -1000, -1000, 368, 753,
	654, 607, 79, -21, 618, 654, 497, -1000, 654, 490,
	-43, 77, -43, -43, 264, 307, 1174, 870, -1000, -1000,
	1066, -1000, 51, -1000, -1000, -1000, -1000, 460, 224, -1000,
	979, 263, 297, 870, 923, 351, 114, -1000, 114, -1000,
	114, -1000, 262, 311, -1000, 912, -1000, -1000, 1040, -1000,
	331, 61, 1055, -1000, 333, 414, 870, 339, -1000, -1000,
	-1000, -22, 58, -24, 571, 654, 870, 148, 402, 388,
	397, -1000, 148, -1000, -1000, -1000, -1000, 870, 870, -1000,
	-1000, 199, 50, -1000, 654, -1000, -1000, -1000, 654, 654,
	44, -38, 870, 43, 654, -1000, 239, -1000, -1000, 327,
	257, 317, -1000, 900, 870, 870, 635, 456, 221, -1000,
	1055, 870, 116, 38, 33, 23, -1000, 227, -43, 1040,
	255, 21, 531, -1000, -1000, -1000, 517, 116, 383, 654,
	-1000, -1000, 1055, 573, 148, 392, 148, 549, 1055, -1000,
	136, 189, -1000, -1000, -1000, 618, 654, 1055, -1000, -1000,
	331, -43, 326, 1174, 131, 20, 17, 870, 634, 1055,
	-1000, -1000, -1000, -1000, -1000, -1000, 294, 1174, 304, -1000,
	133, -1000, -1000, -1000, -1000, 116, -1000, -1000, -1000, 870,
	115, 549, 148, 573, 209, 112, -1000, -1000, -43, -1000,
	-43, -1000, -1000, -1000, 16, 14, 253, 294, 1040, 870,
	-43, 331, -1000, 1055, 654, 549, -1000, 13, 487, 428,
	209, -1000, -1000, -1000, -1000, 322, 245, 316, -1000, 860,
	-1000, 12, 11, -1000, -1000, 870, 9, -43, 319, 1040,
	-1000, -1000, -1000, -1000, -1000, -43, -1000, -1000,
}
var yyPgo = [...]int{

	0, 45, 21, 50, 745, 743, 741, 740, 739, 738,
	737, 736, 734, 733, 65, 76, 91, 728, 47, 48,
	727, 723, 18, 719, 34, 718, 36, 713, 64, 67,
	175, 196, 286, 85, 66, 711, 709, 707, 705, 441,
	702, 701, 699, 698, 697, 696, 695, 694, 44, 14,
	693, 72, 689, 27, 687, 25, 685, 683, 681, 679,
	677, 38, 17, 53, 73, 16, 51, 74, 676, 675,
	673, 41, 672, 671, 669, 60, 26, 33, 668, 667,
	68, 40, 20, 19, 3, 666, 311, 264, 58, 665,
	49, 12, 59, 1, 664, 77, 661, 54, 57, 35,
	42, 82, 659, 330, 0,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 4,
	4, 5, 5, 6, 6, 7, 7, 7, 7, 7,
	8, 8, 8, 9, 9, 9, 9, 9, 9, 9,
	10, 10, 11, 11, 11, 11, 11, 12, 12, 12,
	12, 13, 13, 13, 13, 14, 15, 15, 15, 15,
	16, 16, 17, 18, 18, 19, 19, 20, 20, 21,
	21, 22, 22, 23, 23, 23, 24, 24, 25, 25,
	26, 26, 27, 27, 28, 28, 29, 29, 29, 29,
	29, 29, 30, 30, 31, 31, 31, 31, 31, 31,
	31, 31, 31, 31, 31, 31, 31, 31, 32, 32,
	33, 33, 34, 34, 35, 35, 36, 36, 37, 37,
	37, 38, 38, 39, 40, 41, 41, 41, 41, 41,
	41, 41, 41, 41, 41, 41, 41, 41, 41, 41,
	41, 41, 41, 41, 42, 42, 42, 42, 42, 43,
	43, 43, 44, 44, 45, 45, 46, 46, 46, 47,
	47, 47, 47, 48, 48, 49, 50, 50, 51, 51,
	51, 52, 52, 53, 53, 53, 53, 53, 53, 54,
	54, 54, 54, 54, 55, 55, 55, 56, 56, 56,
	57, 57, 58, 59, 59, 60, 60, 61, 61, 62,
	62, 63, 63, 64, 64, 65, 65, 66, 66, 67,
	67, 68, 68, 68, 68, 69, 70, 71, 71, 72,
	72, 73, 74, 74, 75, 75, 76, 76, 77, 77,
	77, 77, 77, 78, 78, 79, 80, 80, 81, 81,
	82, 82, 83, 83, 84, 85, 86, 86, 87, 87,
	88, 89, 90, 91, 92, 92, 93, 94, 94, 95,
	95, 96, 96, 97, 97, 98, 98, 99, 99, 100,
	100, 100, 100, 101, 101, 102, 102, 103, 103, 104,
	104,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 2, 2, 2, 6, 3, 3, 3, 6,
	6, 9, 6, 0, 1, 1, 1, 1, 2, 2,
	5, 6, 8, 9, 7, 9, 2, 8, 9, 2,
	2, 5, 3, 3, 3, 5, 5, 4, 4, 4,
	1, 1, 3, 0, 2, 0, 2, 0, 3, 0,
	2, 0, 3, 0, 3, 4, 0, 2, 0, 2,
	0, 2, 6, 9, 1, 3, 1, 1, 1, 1,
	1, 1, 1, 3, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 3, 3, 1,
	1, 3, 1, 3, 2, 4, 1, 1, 0, 1,
	1, 1, 1, 3, 3, 3, 3, 3, 3, 4,
	4, 6, 6, 4, 6, 4, 4, 4, 6, 4,
	4, 6, 4, 2, 3, 3, 3, 3, 3, 3,
	3, 2, 3, 4, 4, 1, 1, 2, 2, 7,
	8, 7, 8, 7, 8, 2, 0, 3, 1, 2,
	3, 1, 1, 1, 1, 2, 3, 1, 1, 5,
	5, 6, 6, 4, 0, 2, 4, 1, 1, 1,
	1, 3, 5, 0, 1, 0, 2, 1, 3, 1,
	3, 1, 3, 1, 3, 1, 3, 1, 3, 4,
	2, 6, 9, 5, 8, 7, 3, 1, 3, 5,
	6, 6, 6, 8, 1, 3, 1, 3, 0, 1,
	1, 2, 2, 5, 7, 7, 4, 2, 0, 2,
	4, 2, 0, 2, 1, 1, 1, 2, 1, 2,
	1, 1, 1, 1, 1, 3, 3, 1, 3, 1,
	3, 0, 1, 0, 1, 0, 1, 0, 1, 0,
	1, 1, 1, 0, 1, 0, 1, 1, 1, 0,
	1,
}
var yyChk = [...]int{

	-1000, -1, -3, -14, -68, -69, -72, -73, -74, -78,
	-79, -5, -6, -7, -8, -11, -13, -26, 26, 29,
	107, -93, 99, 100, 84, 88, 89, 90, 87, 76,
	78, 103, 16, 104, 105, 106, 74, -91, 11, -1,
	-104, 121, -104, -104, -104, -104, -104, -104, -104, -15,
	19, 15, 17, -17, -16, 13, -39, 117, 30, 30,
	-95, -94, 11, -104, -104, -104, -84, 4, -84, -84,
	-84, -9, 91, 92, 31, 32, 93, 94, -31, -30,
	-29, -42, -40, -39, -44, -45, -58, -41, -43, -91,
	-93, -10, 117, -84, -85, -86, -87, -88, -89, -90,
	-47, 75, -32, 63, 67, 85, 5, 6, 113, 7,
	9, 10, 71, -31, -92, -91, -104, 12, -31, -62,
	-31, 5, -28, -27, -102, 25, 110, -22, 38, 20,
	-64, -51, -84, 14, -64, -18, 14, 58, 59, 60,
	-96, 73, -14, -26, -84, -84, -104, 119, 110, 85,
	30, -104, -104, -104, -84, -31, -31, 80, 109, 108,
	111, 70, -97, -103, 112, 113, 114, 115, 116, 66,
	65, 67, -31, -62, 120, 117, -59, -31, 108, 111,
	-97, -103, -39, -31, -84, -86, -87, 82, 64, 119,
	111, -104, -104, 119, -104, 119, -84, -31, -23, 44,
	41, -84, 16, 119, -84, 22, -63, -53, -51, -52,
	-54, 23, -39, 24, 14, -19, 18, -63, -101, 61,
	-101, -101, -66, -57, -56, -31, -48, 114, -84, 118,
	117, 27, 28, 36, -95, -31, 86, 117, 86, 20,
	-1, -31, -31, -31, -97, 68, 64, 69, 62, 61,
	-31, -31, -31, -31, -31, -31, -31, 118, 118, -84,
	118, -62, -46, -31, 73, 114, -67, 79, -32, 117,
	-39, -32, 68, 64, 62, 61, 70, -2, -4, -3,
	-12, 76, 101, 102, -84, -92, -91, -29, -62, -28,
	22, 117, -25, 45, -31, -34, -35, -36, -31, -48,
	21, 117, -14, -71, -70, -30, -84, -64, -84, -19,
	-98, 57, -100, 54, 119, 49, 51, 52, 53, -84,
	22, -63, -20, 39, -31, -16, -15, -16, -16, 119,
	22, 117, -65, -84, -75, 117, -84, -30, 117, -30,
	-14, -65, -14, -92, -81, -80, 81, 77, -88, -90,
	-31, -32, -31, -32, -32, 118, 118, 38, -22, 114,
	-31, -60, -67, 81, -31, -32, 117, -39, 117, -39,
	117, -39, -97, 83, -2, -31, -104, -104, 82, -104,
	117, -65, -31, -24, 47, 74, 119, -37, 42, 43,
	-33, -32, -61, -30, -18, 119, 111, 48, -98, -100,
	-99, 50, 48, -63, -84, -19, -21, 40, 41, -66,
	-84, 118, -62, 118, 119, -77, 31, 32, 33, 34,
	-76, -75, 35, -61, 37, -104, 118, -104, -104, 83,
	-81, -80, -1, -31, 65, 41, 96, 38, -22, 83,
	-31, 80, 65, -33, -33, -33, 88, 64, 78, 80,
	-2, -14, 118, -24, 46, -34, 72, 119, 118, 119,
	-19, -71, -31, -53, 48, -99, 48, -53, -31, -62,
	98, 118, -65, -30, -30, 118, 119, -31, 118, -84,
	86, 76, 83, 80, -31, -34, 5, 41, 96, -31,
	-32, 118, 118, 118, 95, -104, -2, -3, 83, 118,
	22, -38, 31, 32, -33, 21, -14, -61, -55, 55,
	56, -53, 48, -53, 117, 98, -77, -76, -14, -104,
	76, -1, 118, 118, -34, 5, -83, -82, 81, 77,
	78, 117, -33, -31, 117, -53, -55, -49, -50, 97,
	117, -104, -104, 118, 118, 83, -83, -82, -2, -31,
	-104, -14, -65, 118, -22, 41, -49, 76, 83, 80,
	118, 118, -62, 118, -104, 76, -2, -104,
}
var yyDef = [...]int{

	-2, -2, -2, 279, 279, 279, 279, 279, 279, 279,
	279, 13, 14, 15, 16, 17, 18, 0, 0, 0,
	0, 279, 279, 279, 0, 0, 0, 0, 33, 0,
	0, 279, 0, 0, 0, 0, 275, 0, 253, 2,
	5, 280, 6, 7, 8, 9, 10, 11, 12, -2,
	0, 0, 0, 63, 0, 261, 61, 80, 0, 0,
	279, 259, 257, 22, 23, 24, 0, 244, 279, 279,
	279, 0, 34, 35, 36, 37, 0, 0, 263, 94,
	95, 96, 97, 98, 99, 100, 101, 102, 103, 104,
	105, 106, 80, 92, 86, 87, 88, 89, 90, 91,
	155, 193, 263, 0, 0, 0, 245, 246, 0, 248,
	250, 251, 252, 263, 0, 104, 46, 0, -2, 279,
	-2, 279, 81, 84, 0, 276, 0, 73, 0, 0,
	0, 203, 168, 0, 0, 65, 0, 273, 273, 273,
	0, 262, 0, 0, 0, 0, 21, 0, 0, 0,
	0, 26, 27, 28, 0, -2, -2, -2, 0, 277,
	278, 263, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 264, 263, 0, 0, 0, 0, -2, 277, 278,
	0, 0, 143, 151, 0, 247, 249, -2, 0, 0,
	0, 52, 53, 0, 54, 275, 0, 256, 78, 0,
	0, 80, 0, 0, 169, 0, 65, -2, 173, 174,
	177, 178, 171, 172, 0, 67, 0, 64, 0, 274,
	0, 0, 62, 207, 190, -2, 188, 189, 92, 123,
	0, 0, 0, 0, 260, -2, 80, 0, 80, 0,
	238, 124, -2, -2, 0, 0, 0, 0, 0, 0,
	144, 145, 146, 147, 148, 149, 150, 107, 108, 93,
	152, 0, 0, -2, 0, 156, 195, 0, 126, 80,
	109, 128, 0, 0, 0, 0, 263, 0, -2, 19,
	20, 0, 279, 279, 0, 255, 254, 279, 200, 85,
	0, 0, 55, 0, -2, 72, 112, 118, -2, 117,
	0, 0, 213, 63, 217, 0, 92, 204, 170, 219,
	0, -2, 267, 0, 0, 266, 270, 271, 272, 175,
	0, 65, 69, 0, -2, 57, 60, 58, 59, 0,
	0, 0, 0, 205, 228, 0, 224, 233, 0, 0,
	279, 0, 279, 279, 0, 238, -2, 0, 129, 130,
	263, 133, -2, 137, 140, 153, 154, 0, 0, 157,
	-2, 0, 210, 0, 263, 0, 80, 135, 80, 139,
	80, 142, 0, 0, 4, 263, 49, 50, -2, 51,
	80, 0, -2, 74, 76, 0, 0, 114, 119, 120,
	211, 110, 0, 197, 65, 0, 0, 0, 0, 267,
	0, 268, 0, 202, 176, 220, 56, 0, 0, 208,
	191, 152, 0, 221, 0, 222, 229, 230, 0, 0,
	0, 226, 0, 0, 0, 25, 30, 32, 29, 0,
	0, 237, 239, 263, 0, 0, 0, 0, 0, 192,
	-2, 0, 0, 0, 0, 0, 40, 0, 279, -2,
	0, 0, 0, 75, 77, 113, 0, 0, 80, 0,
	215, 218, -2, 184, 0, 0, 0, 183, -2, 68,
	0, 153, 206, 231, 232, 228, 0, -2, 234, 235,
	80, 279, 0, -2, 131, 72, 0, 0, 0, -2,
	132, 134, 138, 141, 41, 44, 242, -2, 0, 82,
	0, 115, 121, 122, 111, 0, 214, 198, 179, 0,
	0, 180, 0, 184, 166, 0, 223, 227, 279, 42,
	279, 236, 159, 161, 72, 0, 0, 242, -2, 0,
	279, 80, 212, -2, 0, 182, 181, 0, 71, 0,
	166, 31, 43, 160, 162, 0, 0, 241, 243, 263,
	45, 0, 0, 163, 165, 0, 0, 279, 0, -2,
	83, 186, 167, 164, 47, 279, 240, 48,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 116, 3, 3,
	117, 118, 114, 112, 119, 113, 120, 115, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 121,
	3, 111,
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
	102, 103, 104, 105, 106, 107, 108, 109, 110,
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
		//line parser.y:164
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:169
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:176
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:181
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:188
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:192
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:196
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:200
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:204
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:208
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:212
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:216
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:220
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:224
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:228
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:232
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:236
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
		//line parser.y:246
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:250
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:256
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:260
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:266
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:270
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 25:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:276
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:280
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:284
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:288
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 29:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:292
		{
			yyVAL.statement = FetchCursor{Position: yyDollar[2].expression, Cursor: yyDollar[3].identifier, Variables: yyDollar[5].variables}
		}
	case 30:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:298
		{
			yyVAL.statement = TableDeclaration{Table: yyDollar[2].identifier, Fields: yyDollar[5].expressions}
		}
	case 31:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:302
		{
			yyVAL.statement = TableDeclaration{Table: yyDollar[2].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[8].expression}
		}
	case 32:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:306
		{
			yyVAL.statement = TableDeclaration{Table: yyDollar[2].identifier, Query: yyDollar[5].expression}
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:312
		{
			yyVAL.expression = nil
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:316
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:320
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:324
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:328
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:332
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token, Number: yyDollar[2].expression}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:336
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token, Number: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:342
		{
			yyVAL.expression = CursorStatus{CursorLit: yyDollar[1].token.Literal, Cursor: yyDollar[2].identifier, Is: yyDollar[3].token.Literal, Negation: yyDollar[4].token, Type: yyDollar[5].token.Token, TypeLit: yyDollar[5].token.Literal}
		}
	case 41:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:346
		{
			yyVAL.expression = CursorStatus{CursorLit: yyDollar[1].token.Literal, Cursor: yyDollar[2].identifier, Is: yyDollar[3].token.Literal, Negation: yyDollar[4].token, Type: yyDollar[6].token.Token, TypeLit: yyDollar[5].token.Literal + " " + yyDollar[6].token.Literal}
		}
	case 42:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:352
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 43:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:356
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 44:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:360
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 45:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:364
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:368
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 47:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:374
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 48:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:378
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:382
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:386
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 51:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:392
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:396
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:400
		{
			yyVAL.statement = Printf{Values: yyDollar[2].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:404
		{
			yyVAL.statement = Source{FilePath: yyDollar[2].token.Literal}
		}
	case 55:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:410
		{
			yyVAL.expression = SelectQuery{
				WithClause:    yyDollar[1].expression,
				SelectEntity:  yyDollar[2].expression,
				OrderByClause: yyDollar[3].expression,
				LimitClause:   yyDollar[4].expression,
				OffsetClause:  yyDollar[5].expression,
			}
		}
	case 56:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:422
		{
			yyVAL.expression = SelectEntity{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
			}
		}
	case 57:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:432
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 58:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:441
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:450
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:461
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:465
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:471
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:477
		{
			yyVAL.expression = nil
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:481
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:487
		{
			yyVAL.expression = nil
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:491
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 67:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:497
		{
			yyVAL.expression = nil
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:501
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 69:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:507
		{
			yyVAL.expression = nil
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:511
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 71:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:517
		{
			yyVAL.expression = nil
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:521
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 73:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:527
		{
			yyVAL.expression = nil
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:531
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 75:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:535
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
		}
	case 76:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:541
		{
			yyVAL.expression = nil
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:545
		{
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 78:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:551
		{
			yyVAL.expression = nil
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:555
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 80:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:561
		{
			yyVAL.expression = nil
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:565
		{
			yyVAL.expression = WithClause{With: yyDollar[1].token.Literal, InlineTables: yyDollar[2].expressions}
		}
	case 82:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:571
		{
			yyVAL.expression = InlineTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, As: yyDollar[3].token.Literal, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 83:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:575
		{
			yyVAL.expression = InlineTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, Columns: yyDollar[4].expressions, As: yyDollar[6].token.Literal, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:581
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:585
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:591
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:595
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:599
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:603
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:607
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 91:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:611
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:617
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:621
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:627
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:631
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:635
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:639
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:643
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:647
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:651
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:655
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:659
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:663
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:667
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:671
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:675
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:679
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:685
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:689
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:695
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:699
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:705
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:709
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 114:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:715
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 115:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:719
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:725
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:729
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 118:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:735
		{
			yyVAL.token = Token{}
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:739
		{
			yyVAL.token = yyDollar[1].token
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:743
		{
			yyVAL.token = yyDollar[1].token
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:749
		{
			yyVAL.token = yyDollar[1].token
		}
	case 122:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:753
		{
			yyVAL.token = yyDollar[1].token
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:759
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:765
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
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:788
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:796
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:800
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 129:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:804
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 130:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:808
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 131:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:812
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 132:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:816
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 133:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:820
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 134:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 135:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:828
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 136:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:832
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 137:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:836
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 138:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:840
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 139:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:844
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 140:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:848
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 141:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:852
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 142:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:856
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 143:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:860
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 144:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:866
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:870
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 146:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:874
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:878
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 148:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:882
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 149:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:888
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 150:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:892
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 151:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:896
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 152:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:902
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal}
		}
	case 153:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:906
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Args: yyDollar[3].expressions}
		}
	case 154:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:912
		{
			yyVAL.expression = AggregateFunction{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(AggregateOption)}
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:916
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:922
		{
			yyVAL.expression = AggregateOption{Args: []Expression{AllColumns{}}}
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:926
		{
			yyVAL.expression = AggregateOption{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 158:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:930
		{
			yyVAL.expression = AggregateOption{Distinct: yyDollar[1].token, Args: []Expression{yyDollar[2].expression}}
		}
	case 159:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:936
		{
			orderBy := OrderByClause{OrderBy: yyDollar[4].token.Literal + " " + yyDollar[5].token.Literal, Items: yyDollar[6].expressions}
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Args: []Expression{yyDollar[3].expression}}, OrderBy: orderBy}
		}
	case 160:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:941
		{
			orderBy := OrderByClause{OrderBy: yyDollar[5].token.Literal + " " + yyDollar[6].token.Literal, Items: yyDollar[7].expressions}
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Distinct: yyDollar[3].token, Args: []Expression{yyDollar[4].expression}}, OrderBy: orderBy}
		}
	case 161:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:946
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Args: []Expression{yyDollar[3].expression}}, OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 162:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:950
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Distinct: yyDollar[3].token, Args: []Expression{yyDollar[4].expression}}, OrderBy: yyDollar[5].expression, SeparatorLit: yyDollar[6].token.Literal, Separator: yyDollar[7].token.Literal}
		}
	case 163:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:956
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Over: yyDollar[4].token.Literal, AnalyticClause: yyDollar[6].expression.(AnalyticClause)}
		}
	case 164:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:960
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Args: yyDollar[3].expressions, Over: yyDollar[5].token.Literal, AnalyticClause: yyDollar[7].expression.(AnalyticClause)}
		}
	case 165:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:966
		{
			yyVAL.expression = AnalyticClause{Partition: yyDollar[1].expression, OrderByClause: yyDollar[2].expression}
		}
	case 166:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:972
		{
			yyVAL.expression = nil
		}
	case 167:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:976
		{
			yyVAL.expression = Partition{PartitionBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Values: yyDollar[3].expressions}
		}
	case 168:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:982
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 169:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:986
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 170:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:990
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:996
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1000
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 173:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1006
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 174:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1010
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 175:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1014
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 176:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1018
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1022
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 178:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1026
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 179:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1032
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 180:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1036
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 181:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1040
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 182:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1044
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 183:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1048
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 184:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1054
		{
			yyVAL.expression = nil
		}
	case 185:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1058
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 186:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1062
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 187:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1068
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 188:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1072
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1076
		{
			yyVAL.expression = AllColumns{}
		}
	case 190:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1082
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 191:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1086
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 192:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1092
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 193:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1098
		{
			yyVAL.expression = nil
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1102
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 195:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1108
		{
			yyVAL.expression = nil
		}
	case 196:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1112
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 197:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1118
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 198:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1122
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 199:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1128
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 200:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1132
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 201:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1138
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 202:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1142
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1148
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 204:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1152
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1158
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1162
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1168
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 208:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1172
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 209:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1178
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 210:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1182
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 211:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1188
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Values: yyDollar[5].token.Literal, ValuesList: yyDollar[6].expressions}
		}
	case 212:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:1192
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Values: yyDollar[8].token.Literal, ValuesList: yyDollar[9].expressions}
		}
	case 213:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1196
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 214:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1200
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 215:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1206
		{
			yyVAL.expression = UpdateQuery{WithClause: yyDollar[1].expression, Update: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, Set: yyDollar[4].token.Literal, SetList: yyDollar[5].expressions, FromClause: yyDollar[6].expression, WhereClause: yyDollar[7].expression}
		}
	case 216:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1212
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1218
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 218:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1222
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 219:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1228
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{WithClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 220:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1233
		{
			from := FromClause{From: yyDollar[4].token.Literal, Tables: yyDollar[5].expressions}
			yyVAL.expression = DeleteQuery{WithClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, FromClause: from, WhereClause: yyDollar[6].expression}
		}
	case 221:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1240
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 222:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1246
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 223:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1250
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1256
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 225:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1260
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1266
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 227:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1270
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 228:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1276
		{
			yyVAL.expression = nil
		}
	case 229:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1280
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 230:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1284
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 231:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1288
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 232:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1292
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 233:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1298
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 234:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1302
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 235:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1308
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 236:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1314
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 237:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1318
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 238:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1324
		{
			yyVAL.procexpr = nil
		}
	case 239:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1328
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 240:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1334
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 241:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1338
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 242:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1344
		{
			yyVAL.procexpr = nil
		}
	case 243:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1348
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 244:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1354
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 245:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1360
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 246:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1366
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 247:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1370
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 248:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1377
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 249:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1381
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 250:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1388
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 251:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1394
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 252:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1400
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 253:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1406
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 254:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1412
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 255:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1416
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 256:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1422
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 257:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1428
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 258:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1432
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 259:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1438
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 260:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1442
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 261:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1448
		{
			yyVAL.token = Token{}
		}
	case 262:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1452
		{
			yyVAL.token = yyDollar[1].token
		}
	case 263:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1458
		{
			yyVAL.token = Token{}
		}
	case 264:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1462
		{
			yyVAL.token = yyDollar[1].token
		}
	case 265:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1468
		{
			yyVAL.token = Token{}
		}
	case 266:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1472
		{
			yyVAL.token = yyDollar[1].token
		}
	case 267:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1478
		{
			yyVAL.token = Token{}
		}
	case 268:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1482
		{
			yyVAL.token = yyDollar[1].token
		}
	case 269:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1488
		{
			yyVAL.token = Token{}
		}
	case 270:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1492
		{
			yyVAL.token = yyDollar[1].token
		}
	case 271:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1496
		{
			yyVAL.token = yyDollar[1].token
		}
	case 272:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1500
		{
			yyVAL.token = yyDollar[1].token
		}
	case 273:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1506
		{
			yyVAL.token = Token{}
		}
	case 274:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1510
		{
			yyVAL.token = yyDollar[1].token
		}
	case 275:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1516
		{
			yyVAL.token = Token{}
		}
	case 276:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1520
		{
			yyVAL.token = yyDollar[1].token
		}
	case 277:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1527
		{
			yyVAL.token = yyDollar[1].token
		}
	case 278:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1531
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 279:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1537
		{
			yyVAL.token = Token{}
		}
	case 280:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1541
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
