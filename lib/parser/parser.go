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

//line parser.y:1508

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
	-2, 78,
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	1, 1,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 47,
	58, 58,
	59, 58,
	60, 58,
	-2, 69,
	-1, 116,
	64, 255,
	68, 255,
	69, 255,
	-2, 271,
	-1, 150,
	4, 38,
	-2, 255,
	-1, 151,
	4, 39,
	-2, 255,
	-1, 152,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 170,
	117, 151,
	-2, 253,
	-1, 172,
	79, 186,
	-2, 255,
	-1, 182,
	38, 151,
	97, 151,
	117, 151,
	-2, 253,
	-1, 183,
	83, 3,
	-2, 78,
	-1, 200,
	48, 257,
	50, 261,
	-2, 193,
	-1, 218,
	64, 255,
	68, 255,
	69, 255,
	-2, 179,
	-1, 228,
	64, 255,
	68, 255,
	69, 255,
	-2, 250,
	-1, 235,
	70, 0,
	107, 0,
	110, 0,
	-2, 122,
	-1, 236,
	70, 0,
	107, 0,
	110, 0,
	-2, 124,
	-1, 269,
	77, 3,
	81, 3,
	83, 3,
	-2, 78,
	-1, 284,
	64, 255,
	68, 255,
	69, 255,
	-2, 74,
	-1, 288,
	64, 255,
	68, 255,
	69, 255,
	-2, 113,
	-1, 301,
	50, 261,
	-2, 257,
	-1, 314,
	64, 255,
	68, 255,
	69, 255,
	-2, 64,
	-1, 321,
	117, 151,
	-2, 253,
	-1, 336,
	83, 1,
	-2, 78,
	-1, 342,
	70, 0,
	107, 0,
	110, 0,
	-2, 133,
	-1, 346,
	64, 255,
	68, 255,
	69, 255,
	-2, 191,
	-1, 368,
	83, 3,
	-2, 78,
	-1, 372,
	64, 255,
	68, 255,
	69, 255,
	-2, 77,
	-1, 425,
	83, 188,
	-2, 255,
	-1, 436,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 449,
	64, 255,
	68, 255,
	69, 255,
	-2, 208,
	-1, 455,
	64, 255,
	68, 255,
	69, 255,
	-2, 68,
	-1, 463,
	64, 255,
	68, 255,
	69, 255,
	-2, 217,
	-1, 469,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 471,
	79, 201,
	81, 201,
	83, 201,
	-2, 255,
	-1, 480,
	77, 1,
	81, 1,
	83, 1,
	-2, 19,
	-1, 507,
	83, 3,
	-2, 78,
	-1, 512,
	64, 255,
	68, 255,
	69, 255,
	-2, 177,
	-1, 532,
	77, 3,
	81, 3,
	83, 3,
	-2, 78,
}

const yyPrivate = 57344

const yyLast = 1104

var yyAct = [...]int{

	38, 2, 322, 90, 122, 40, 41, 42, 43, 44,
	45, 46, 87, 21, 168, 21, 268, 491, 380, 404,
	505, 409, 61, 62, 63, 521, 382, 285, 64, 66,
	67, 68, 114, 390, 86, 35, 208, 35, 37, 1,
	293, 334, 200, 219, 215, 302, 157, 373, 199, 300,
	96, 94, 130, 127, 127, 138, 17, 254, 17, 141,
	139, 140, 410, 112, 420, 113, 52, 146, 147, 148,
	78, 201, 351, 149, 58, 211, 125, 316, 39, 169,
	117, 462, 305, 446, 306, 307, 308, 303, 321, 170,
	301, 169, 169, 444, 121, 47, 434, 166, 403, 346,
	156, 385, 99, 185, 376, 185, 179, 165, 164, 166,
	319, 196, 156, 188, 142, 535, 433, 187, 534, 533,
	504, 126, 126, 189, 290, 482, 475, 129, 194, 76,
	111, 197, 127, 116, 474, 127, 473, 154, 153, 221,
	155, 159, 160, 161, 162, 163, 175, 53, 464, 154,
	153, 304, 155, 159, 160, 161, 162, 163, 461, 386,
	457, 445, 39, 153, 439, 21, 159, 160, 161, 162,
	163, 415, 516, 253, 150, 151, 402, 34, 347, 65,
	210, 161, 162, 163, 252, 270, 222, 35, 275, 167,
	34, 233, 513, 65, 65, 259, 21, 221, 172, 296,
	127, 178, 298, 237, 510, 370, 309, 359, 17, 213,
	214, 127, 204, 206, 77, 231, 47, 227, 35, 291,
	277, 190, 159, 160, 161, 162, 163, 323, 326, 296,
	296, 357, 259, 331, 323, 218, 299, 289, 355, 17,
	267, 280, 223, 228, 256, 230, 182, 55, 186, 276,
	55, 121, 143, 234, 235, 236, 311, 278, 158, 243,
	244, 245, 246, 247, 248, 249, 345, 277, 126, 279,
	349, 270, 362, 297, 366, 367, 258, 261, 315, 369,
	317, 318, 21, 497, 371, 323, 364, 324, 339, 338,
	529, 328, 284, 288, 476, 296, 333, 53, 335, 49,
	466, 50, 229, 48, 35, 55, 325, 145, 127, 314,
	531, 519, 93, 361, 394, 432, 166, 481, 468, 104,
	106, 92, 424, 221, 400, 17, 3, 418, 363, 326,
	368, 414, 296, 416, 417, 281, 257, 508, 340, 431,
	342, 507, 341, 509, 343, 344, 384, 389, 395, 21,
	337, 388, 435, 393, 336, 412, 173, 353, 176, 174,
	508, 257, 144, 352, 399, 354, 337, 539, 530, 502,
	270, 35, 365, 467, 428, 421, 429, 419, 430, 401,
	221, 21, 137, 372, 34, 437, 488, 443, 375, 296,
	81, 127, 17, 381, 136, 266, 127, 107, 263, 166,
	55, 427, 262, 35, 442, 184, 458, 323, 54, 212,
	295, 296, 296, 456, 265, 264, 137, 465, 181, 218,
	289, 447, 440, 452, 17, 105, 448, 180, 239, 133,
	450, 495, 238, 240, 391, 454, 478, 422, 480, 34,
	327, 329, 242, 241, 132, 133, 134, 453, 255, 21,
	296, 451, 425, 479, 392, 127, 387, 127, 381, 109,
	381, 441, 381, 487, 283, 192, 326, 537, 501, 378,
	379, 35, 398, 490, 193, 233, 288, 397, 313, 123,
	413, 498, 21, 411, 499, 305, 449, 306, 307, 308,
	57, 177, 17, 21, 494, 56, 496, 455, 120, 127,
	65, 517, 135, 518, 35, 506, 383, 511, 503, 270,
	524, 463, 485, 486, 515, 35, 526, 323, 310, 205,
	21, 292, 205, 470, 522, 17, 471, 520, 65, 54,
	472, 538, 483, 536, 270, 320, 17, 232, 514, 124,
	542, 195, 35, 383, 209, 21, 198, 381, 131, 540,
	224, 225, 541, 207, 115, 36, 330, 60, 332, 226,
	477, 65, 65, 17, 260, 260, 119, 35, 65, 103,
	104, 106, 128, 107, 108, 36, 59, 65, 103, 104,
	106, 95, 107, 108, 36, 91, 137, 10, 17, 9,
	8, 381, 512, 405, 406, 407, 408, 7, 205, 6,
	295, 294, 54, 5, 54, 54, 4, 350, 523, 171,
	83, 216, 65, 103, 104, 106, 217, 107, 108, 36,
	203, 202, 459, 460, 528, 527, 97, 100, 82, 85,
	260, 101, 260, 260, 79, 109, 100, 84, 80, 98,
	101, 484, 377, 287, 109, 286, 118, 34, 98, 102,
	282, 191, 396, 260, 356, 358, 360, 312, 102, 51,
	110, 383, 16, 165, 164, 166, 271, 15, 156, 110,
	88, 100, 69, 14, 13, 101, 105, 220, 12, 109,
	89, 260, 137, 98, 137, 105, 137, 11, 269, 89,
	0, 0, 0, 102, 165, 205, 166, 438, 0, 156,
	0, 0, 0, 0, 110, 154, 153, 0, 155, 159,
	160, 161, 162, 163, 0, 250, 251, 165, 164, 166,
	105, 348, 156, 0, 89, 0, 0, 0, 65, 103,
	104, 106, 0, 107, 108, 36, 154, 153, 0, 155,
	159, 160, 161, 162, 163, 0, 260, 0, 260, 0,
	260, 103, 104, 106, 0, 107, 108, 0, 0, 154,
	153, 374, 155, 159, 160, 161, 162, 163, 0, 0,
	251, 0, 489, 0, 0, 0, 0, 0, 205, 165,
	164, 166, 0, 205, 156, 0, 0, 100, 375, 0,
	0, 101, 0, 500, 0, 109, 0, 0, 305, 98,
	306, 307, 308, 303, 492, 493, 301, 0, 0, 102,
	72, 73, 0, 0, 0, 0, 0, 109, 260, 0,
	110, 154, 153, 0, 155, 159, 160, 161, 162, 163,
	0, 0, 0, 0, 0, 260, 105, 525, 0, 0,
	89, 0, 205, 0, 205, 165, 164, 166, 0, 0,
	156, 0, 0, 0, 165, 164, 166, 0, 105, 156,
	532, 0, 0, 165, 164, 166, 0, 0, 156, 469,
	70, 71, 74, 75, 0, 0, 0, 0, 436, 260,
	0, 0, 0, 0, 0, 0, 205, 154, 153, 0,
	155, 159, 160, 161, 162, 163, 154, 153, 0, 155,
	159, 160, 161, 162, 163, 154, 153, 0, 155, 159,
	160, 161, 162, 163, 165, 164, 166, 0, 0, 156,
	0, 0, 0, 165, 164, 166, 0, 0, 156, 426,
	0, 0, 165, 164, 166, 0, 0, 156, 0, 0,
	183, 165, 164, 166, 0, 0, 156, 152, 0, 0,
	423, 164, 166, 0, 0, 156, 154, 153, 0, 155,
	159, 160, 161, 162, 163, 154, 153, 0, 155, 159,
	160, 161, 162, 163, 154, 153, 0, 155, 159, 160,
	161, 162, 163, 154, 153, 0, 155, 159, 160, 161,
	162, 163, 154, 153, 0, 155, 159, 160, 161, 162,
	163, 36, 0, 0, 0, 0, 32, 305, 36, 306,
	307, 308, 303, 32, 0, 301, 18, 0, 0, 19,
	0, 0, 0, 18, 0, 0, 19, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 34, 0, 272, 0, 30, 0,
	0, 34, 0, 29, 24, 30, 0, 28, 25, 26,
	27, 24, 0, 0, 28, 25, 26, 27, 0, 0,
	22, 23, 273, 274, 31, 33, 20, 22, 23, 0,
	0, 31, 33, 20,
}
var yyPact = [...]int{

	997, -1000, 997, -42, -42, -42, -42, -42, -42, -42,
	-42, -1000, -1000, -1000, -1000, -1000, -1000, 284, 465, 460,
	546, -42, -42, -42, 557, 557, 557, 557, 779, 724,
	724, -42, 542, 724, 473, 142, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 441, 519, 557,
	558, 534, 386, 321, -1000, 310, 557, 557, -42, -4,
	143, -1000, -1000, -1000, 277, -1000, -42, -42, -42, 557,
	-1000, -1000, -1000, -1000, 724, 724, 867, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 142, -1000, -1000, 573,
	-27, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 724, 249,
	131, 724, 557, -1000, -1000, 313, -1000, -1000, -1000, -1000,
	130, 858, 341, -15, -1000, 138, 42, -1000, -5, 557,
	-1000, 724, 421, 433, 557, 525, -7, 524, 189, 539,
	526, 189, 348, 348, 348, 564, -1000, 69, 134, 126,
	523, -1000, 546, 724, 216, 129, -1000, -1000, -1000, 517,
	876, 876, 997, 724, 724, 724, 332, 364, 381, 724,
	724, 724, 724, 724, 724, 724, -1000, 598, 67, 557,
	321, 257, 876, 79, 79, 334, 353, -1000, 30, 325,
	-1000, -1000, 321, 990, 557, 544, 746, -1000, 473, 219,
	876, 419, 724, 724, 103, 557, 557, -1000, 557, 526,
	33, -1000, 496, -1000, -1000, -1000, -1000, 189, 439, 724,
	-1000, 134, -1000, 134, 134, -1000, -8, 513, 876, -1000,
	-1000, -28, -1000, 557, 190, 175, 557, -1000, 876, 310,
	557, 310, 544, 273, 111, 55, 55, 388, 724, 79,
	724, 79, 79, 68, 68, -1000, -1000, -1000, 629, 30,
	-1000, 724, -1000, -1000, 61, 608, 282, 724, -1000, 573,
	-1000, -1000, 79, 122, 115, 91, 332, 441, 245, 990,
	-1000, -1000, 724, -42, -42, 248, -1000, -13, -42, -1000,
	89, 557, -1000, 724, 714, -1000, -14, 427, 876, -1000,
	79, 557, -1000, 534, -17, 49, -40, -1000, -1000, -1000,
	408, 436, 384, 406, 189, -1000, -1000, -1000, -1000, -1000,
	557, 526, 437, 431, 876, 370, -1000, -1000, 370, 564,
	557, 321, 59, -20, 562, 557, 448, -1000, 557, 443,
	-42, 54, -42, -42, 244, 273, 997, 724, -1000, -1000,
	885, -1000, 55, -1000, -1000, -1000, 652, -1000, -1000, -1000,
	239, 257, 724, 849, 336, 116, -1000, 116, -1000, 116,
	-1000, 251, -1, 274, -1000, 798, -1000, -1000, 990, -1000,
	310, 47, 876, -1000, 314, 415, 724, 315, -1000, -1000,
	-1000, -25, 44, -35, 526, 557, 724, 189, 403, 384,
	399, -1000, 189, -1000, -1000, -1000, -1000, 724, 724, -1000,
	-1000, 43, -1000, 557, -1000, -1000, -1000, 557, 557, 41,
	-37, 724, 31, 557, -1000, 214, -1000, -1000, 297, 235,
	289, -1000, 789, 724, -1000, 876, 724, 79, 19, 17,
	9, -1000, 199, -1000, 555, -42, 990, 234, 8, 510,
	-1000, -1000, -1000, 481, 79, 365, 557, -1000, -1000, 876,
	749, 189, 383, 189, 958, 876, -1000, 184, -1000, -1000,
	-1000, 562, 557, 876, -1000, -1000, 310, -42, 293, 997,
	30, 876, -1000, -1000, -1000, -1000, -1000, 3, -1000, 260,
	997, 265, -1000, 88, -1000, -1000, -1000, -1000, 79, -1000,
	-1000, -1000, 724, 76, 958, 189, 749, 56, -1000, -1000,
	-42, -1000, -42, -1000, -1000, 228, 260, 990, 724, -42,
	310, -1000, 876, 557, 958, -1000, 192, -1000, -1000, 292,
	227, 283, -1000, 780, -1000, 2, 1, -2, 441, 426,
	-42, 291, 990, -1000, -1000, -1000, -1000, 724, -1000, -42,
	-1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 38, 16, 1, 688, 687, 678, 674, 673, 672,
	670, 667, 666, 662, 326, 77, 66, 659, 52, 36,
	657, 652, 4, 651, 47, 650, 55, 646, 80, 70,
	214, 99, 102, 18, 27, 645, 643, 642, 641, 390,
	638, 637, 634, 629, 628, 57, 626, 43, 625, 624,
	71, 621, 42, 620, 17, 616, 611, 610, 609, 607,
	26, 14, 48, 76, 2, 44, 72, 606, 603, 601,
	40, 599, 597, 590, 62, 21, 19, 589, 587, 64,
	41, 25, 20, 3, 585, 321, 312, 51, 581, 50,
	34, 63, 12, 576, 74, 448, 46, 49, 33, 45,
	75, 566, 258, 0,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 4,
	4, 5, 5, 6, 6, 7, 7, 7, 7, 7,
	8, 8, 8, 9, 9, 9, 9, 9, 9, 9,
	10, 10, 11, 11, 11, 11, 11, 12, 12, 12,
	12, 13, 13, 14, 15, 15, 15, 15, 16, 16,
	17, 18, 18, 19, 19, 20, 20, 21, 21, 22,
	22, 23, 23, 23, 24, 24, 25, 25, 26, 26,
	27, 27, 28, 28, 29, 29, 29, 29, 29, 29,
	30, 30, 31, 31, 31, 31, 31, 31, 31, 31,
	31, 31, 31, 31, 31, 32, 32, 33, 33, 34,
	34, 35, 35, 36, 36, 37, 37, 37, 38, 38,
	39, 40, 41, 41, 41, 41, 41, 41, 41, 41,
	41, 41, 41, 41, 41, 41, 41, 41, 41, 41,
	41, 42, 42, 42, 42, 42, 43, 43, 43, 44,
	44, 45, 45, 45, 46, 46, 47, 48, 49, 49,
	50, 50, 50, 51, 51, 52, 52, 52, 52, 52,
	52, 53, 53, 53, 53, 53, 54, 54, 54, 55,
	55, 55, 56, 56, 57, 58, 58, 59, 59, 60,
	60, 61, 61, 62, 62, 63, 63, 64, 64, 65,
	65, 66, 66, 67, 67, 67, 67, 68, 69, 70,
	70, 71, 71, 72, 73, 73, 74, 74, 75, 75,
	76, 76, 76, 76, 76, 77, 77, 78, 79, 79,
	80, 80, 81, 81, 82, 82, 83, 84, 85, 85,
	86, 86, 87, 88, 89, 90, 91, 91, 92, 93,
	93, 94, 94, 95, 95, 96, 96, 97, 97, 98,
	98, 99, 99, 99, 99, 100, 100, 101, 101, 102,
	102, 103, 103,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	1, 3, 2, 2, 2, 6, 3, 3, 3, 6,
	6, 9, 6, 0, 1, 1, 1, 1, 2, 2,
	5, 6, 8, 9, 7, 9, 2, 8, 9, 2,
	2, 5, 3, 5, 5, 4, 4, 4, 1, 1,
	3, 0, 2, 0, 2, 0, 3, 0, 2, 0,
	3, 0, 3, 4, 0, 2, 0, 2, 0, 2,
	6, 9, 1, 3, 1, 1, 1, 1, 1, 1,
	1, 3, 1, 1, 1, 1, 1, 1, 1, 1,
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
	3, 4, 2, 6, 9, 5, 8, 7, 3, 1,
	3, 5, 6, 6, 6, 8, 1, 3, 1, 3,
	0, 1, 1, 2, 2, 5, 7, 7, 4, 2,
	0, 2, 4, 2, 0, 2, 1, 1, 1, 2,
	1, 2, 1, 1, 1, 1, 1, 3, 3, 1,
	3, 1, 3, 0, 1, 0, 1, 0, 1, 0,
	1, 0, 1, 1, 1, 0, 1, 0, 1, 1,
	1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -14, -67, -68, -71, -72, -73, -77,
	-78, -5, -6, -7, -8, -11, -13, -26, 26, 29,
	106, -92, 100, 101, 84, 88, 89, 90, 87, 76,
	78, 104, 16, 105, 74, -90, 11, -1, -103, 120,
	-103, -103, -103, -103, -103, -103, -103, -15, 19, 15,
	17, -17, -16, 13, -39, 116, 30, 30, -94, -93,
	11, -103, -103, -103, -83, 4, -83, -83, -83, -9,
	91, 92, 31, 32, 93, 94, -31, -30, -29, -42,
	-40, -39, -44, -57, -41, -43, -90, -92, -10, 116,
	-83, -84, -85, -86, -87, -88, -89, -46, 75, -32,
	63, 67, 85, 5, 6, 112, 7, 9, 10, 71,
	96, -31, -91, -90, -103, 12, -31, -28, -27, -101,
	25, 109, -22, 38, 20, -63, -50, -83, 14, -63,
	-18, 14, 58, 59, 60, -95, 73, -14, -26, -83,
	-83, -103, 118, 109, 85, 30, -103, -103, -103, -83,
	-31, -31, 80, 108, 107, 110, 70, -96, -102, 111,
	112, 113, 114, 115, 66, 65, 67, -31, -61, 119,
	116, -58, -31, 107, 110, -96, -102, -39, -31, -83,
	-85, -86, 116, 82, 64, 118, 110, -103, 118, -83,
	-31, -23, 44, 41, -83, 16, 118, -83, 22, -62,
	-52, -50, -51, -53, 23, -39, 24, 14, -19, 18,
	-62, -100, 61, -100, -100, -65, -56, -55, -31, -47,
	113, -83, 117, 116, 27, 28, 36, -94, -31, 86,
	116, 86, 20, -1, -31, -31, -31, -96, 68, 64,
	69, 62, 61, -31, -31, -31, -31, -31, -31, -31,
	117, 118, 117, -83, -45, -95, -66, 79, -32, 116,
	-39, -32, 68, 64, 62, 61, 70, -45, -2, -4,
	-3, -12, 76, 102, 103, -83, -91, -90, -29, -28,
	22, 116, -25, 45, -31, -34, -35, -36, -31, -47,
	21, 116, -14, -70, -69, -30, -83, -63, -83, -19,
	-97, 57, -99, 54, 118, 49, 51, 52, 53, -83,
	22, -62, -20, 39, -31, -16, -15, -16, -16, 118,
	22, 116, -64, -83, -74, 116, -83, -30, 116, -30,
	-14, -64, -14, -91, -80, -79, 81, 77, -87, -89,
	-31, -32, -31, -32, -32, -61, -31, 117, 113, -61,
	-59, -66, 81, -31, -32, 116, -39, 116, -39, 116,
	-39, -96, -22, 83, -2, -31, -103, -103, 82, -103,
	116, -64, -31, -24, 47, 74, 118, -37, 42, 43,
	-33, -32, -60, -30, -18, 118, 110, 48, -97, -99,
	-98, 50, 48, -62, -83, -19, -21, 40, 41, -65,
	-83, -45, 117, 118, -76, 31, 32, 33, 34, -75,
	-74, 35, -60, 37, -103, 117, -103, -103, 83, -80,
	-79, -1, -31, 65, 83, -31, 80, 65, -33, -33,
	-33, 88, 64, 117, 97, 78, 80, -2, -14, 117,
	-24, 46, -34, 72, 118, 117, 118, -19, -70, -31,
	-52, 48, -98, 48, -52, -31, -61, 117, -64, -30,
	-30, 117, 118, -31, 117, -83, 86, 76, 83, 80,
	-31, -31, -32, 117, 117, 117, 95, 5, -103, -2,
	-3, 83, 117, 22, -38, 31, 32, -33, 21, -14,
	-60, -54, 55, 56, -52, 48, -52, 99, -76, -75,
	-14, -103, 76, -1, 117, -82, -81, 81, 77, 78,
	116, -33, -31, 116, -52, -54, 116, -103, -103, 83,
	-82, -81, -2, -31, -103, -14, -64, -48, -49, 98,
	76, 83, 80, 117, 117, 117, -22, 41, -103, 76,
	-2, -61, -103,
}
var yyDef = [...]int{

	-2, -2, -2, 271, 271, 271, 271, 271, 271, 271,
	271, 13, 14, 15, 16, 17, 18, 0, 0, 0,
	0, 271, 271, 271, 0, 0, 0, 0, 33, 0,
	0, 271, 0, 0, 267, 0, 245, 2, 5, 272,
	6, 7, 8, 9, 10, 11, 12, -2, 0, 0,
	0, 61, 0, 253, 59, 78, 0, 0, 271, 251,
	249, 22, 23, 24, 0, 236, 271, 271, 271, 0,
	34, 35, 36, 37, 0, 0, 255, 92, 93, 94,
	95, 96, 97, 98, 99, 100, 101, 102, 103, 78,
	90, 84, 85, 86, 87, 88, 89, 150, 185, 255,
	0, 0, 0, 237, 238, 0, 240, 242, 243, 244,
	0, 255, 0, 101, 46, 0, -2, 79, 82, 0,
	268, 0, 71, 0, 0, 0, 195, 160, 0, 0,
	63, 0, 265, 265, 265, 0, 254, 0, 0, 0,
	0, 21, 0, 0, 0, 0, 26, 27, 28, 0,
	-2, -2, -2, 0, 269, 270, 255, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 256, 255, 0, 0,
	-2, 0, -2, 269, 270, 0, 0, 140, 148, 0,
	239, 241, -2, -2, 0, 0, 0, 52, 267, 0,
	248, 76, 0, 0, 78, 0, 0, 161, 0, 63,
	-2, 165, 166, 169, 170, 163, 164, 0, 65, 0,
	62, 0, 266, 0, 0, 60, 199, 182, -2, 180,
	181, 90, 120, 0, 0, 0, 0, 252, -2, 78,
	0, 78, 0, 230, 121, -2, -2, 0, 0, 0,
	0, 0, 0, 141, 142, 143, 144, 145, 146, 147,
	104, 0, 105, 91, 0, 0, 187, 0, 123, 78,
	106, 125, 0, 0, 0, 0, 255, 69, 0, -2,
	19, 20, 0, 271, 271, 0, 247, 246, 271, 83,
	0, 0, 53, 0, -2, 70, 109, 115, -2, 114,
	0, 0, 205, 61, 209, 0, 90, 196, 162, 211,
	0, -2, 259, 0, 0, 258, 262, 263, 264, 167,
	0, 63, 67, 0, -2, 55, 58, 56, 57, 0,
	0, -2, 0, 197, 220, 0, 216, 225, 0, 0,
	271, 0, 271, 271, 0, 230, -2, 0, 126, 127,
	255, 130, -2, 134, 137, 192, -2, 149, 152, 153,
	0, 202, 0, 255, 0, 78, 132, 78, 136, 78,
	139, 0, 0, 0, 4, 255, 49, 50, -2, 51,
	78, 0, -2, 72, 74, 0, 0, 111, 116, 117,
	203, 107, 0, 189, 63, 0, 0, 0, 0, 259,
	0, 260, 0, 194, 168, 212, 54, 0, 0, 200,
	183, 0, 213, 0, 214, 221, 222, 0, 0, 0,
	218, 0, 0, 0, 25, 30, 32, 29, 0, 0,
	229, 231, 255, 0, 184, -2, 0, 0, 0, 0,
	0, 40, 0, 154, 0, 271, -2, 0, 0, 0,
	73, 75, 110, 0, 0, 78, 0, 207, 210, -2,
	176, 0, 0, 0, 175, -2, 66, 149, 198, 223,
	224, 220, 0, -2, 226, 227, 78, 271, 0, -2,
	128, -2, 129, 131, 135, 138, 41, 0, 44, 234,
	-2, 0, 80, 0, 112, 118, 119, 108, 0, 206,
	190, 171, 0, 0, 172, 0, 176, 0, 215, 219,
	271, 42, 271, 228, 155, 0, 234, -2, 0, 271,
	78, 204, -2, 0, 174, 173, 158, 31, 43, 0,
	0, 233, 235, 255, 45, 0, 0, 0, 69, 0,
	271, 0, -2, 81, 178, 156, 157, 0, 47, 271,
	232, 159, 48,
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
		//line parser.y:163
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:168
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:175
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:180
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:187
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:191
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:195
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:199
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:203
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:207
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:211
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:215
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:219
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:223
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:227
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:231
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:235
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:239
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:245
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:249
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:255
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:259
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:265
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:269
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 25:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:275
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:279
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:283
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:287
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 29:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:291
		{
			yyVAL.statement = FetchCursor{Position: yyDollar[2].expression, Cursor: yyDollar[3].identifier, Variables: yyDollar[5].variables}
		}
	case 30:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:297
		{
			yyVAL.statement = TableDeclaration{Table: yyDollar[2].identifier, Fields: yyDollar[5].expressions}
		}
	case 31:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:301
		{
			yyVAL.statement = TableDeclaration{Table: yyDollar[2].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[8].expression}
		}
	case 32:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:305
		{
			yyVAL.statement = TableDeclaration{Table: yyDollar[2].identifier, Query: yyDollar[5].expression}
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:311
		{
			yyVAL.expression = nil
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:315
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:319
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:323
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:327
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token}
		}
	case 38:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:331
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token, Number: yyDollar[2].expression}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:335
		{
			yyVAL.expression = FetchPosition{Position: yyDollar[1].token, Number: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = CursorStatus{CursorLit: yyDollar[1].token.Literal, Cursor: yyDollar[2].identifier, Is: yyDollar[3].token.Literal, Negation: yyDollar[4].token, Type: yyDollar[5].token.Token, TypeLit: yyDollar[5].token.Literal}
		}
	case 41:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = CursorStatus{CursorLit: yyDollar[1].token.Literal, Cursor: yyDollar[2].identifier, Is: yyDollar[3].token.Literal, Negation: yyDollar[4].token, Type: yyDollar[6].token.Token, TypeLit: yyDollar[5].token.Literal + " " + yyDollar[6].token.Literal}
		}
	case 42:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:351
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 43:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:355
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 44:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:359
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 45:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:363
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:367
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 47:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:373
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 48:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:377
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:381
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:385
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 51:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:391
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:395
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 53:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:401
		{
			yyVAL.expression = SelectQuery{
				WithClause:    yyDollar[1].expression,
				SelectEntity:  yyDollar[2].expression,
				OrderByClause: yyDollar[3].expression,
				LimitClause:   yyDollar[4].expression,
				OffsetClause:  yyDollar[5].expression,
			}
		}
	case 54:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:413
		{
			yyVAL.expression = SelectEntity{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
			}
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:423
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 56:
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
	case 57:
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
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:452
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:456
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:462
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:468
		{
			yyVAL.expression = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:472
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:478
		{
			yyVAL.expression = nil
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:482
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:488
		{
			yyVAL.expression = nil
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:492
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 67:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:498
		{
			yyVAL.expression = nil
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:502
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 69:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:508
		{
			yyVAL.expression = nil
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:512
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 71:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:518
		{
			yyVAL.expression = nil
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:522
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 73:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:526
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
		}
	case 74:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = nil
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 76:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:542
		{
			yyVAL.expression = nil
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:546
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 78:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = nil
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = WithClause{With: yyDollar[1].token.Literal, InlineTables: yyDollar[2].expressions}
		}
	case 80:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:562
		{
			yyVAL.expression = InlineTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, As: yyDollar[3].token.Literal, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 81:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:566
		{
			yyVAL.expression = InlineTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, Columns: yyDollar[4].expressions, As: yyDollar[6].token.Literal, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:582
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:586
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:590
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:594
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:598
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:602
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:608
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:618
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:622
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:626
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:630
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:634
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:638
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:642
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:646
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:650
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:654
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:658
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:662
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:666
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:672
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:676
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:682
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:686
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:692
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:696
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:702
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 112:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:706
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:712
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:716
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 115:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:722
		{
			yyVAL.token = Token{}
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:726
		{
			yyVAL.token = yyDollar[1].token
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:730
		{
			yyVAL.token = yyDollar[1].token
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:736
		{
			yyVAL.token = yyDollar[1].token
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:740
		{
			yyVAL.token = yyDollar[1].token
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:746
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:752
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
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:775
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:779
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:783
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:787
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 126:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:791
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 127:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:795
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 128:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:799
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 129:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:803
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 130:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:807
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 131:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:811
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 132:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:815
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 133:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:819
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 134:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:823
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 135:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:827
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 136:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:831
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 137:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:835
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 138:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:839
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 139:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:843
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 140:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:847
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 141:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:853
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:861
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 144:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:865
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:869
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 146:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:875
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:879
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 148:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:883
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 149:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:889
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:893
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 151:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:899
		{
			yyVAL.expression = Option{}
		}
	case 152:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:903
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:907
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 154:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:913
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 155:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:917
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 156:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:923
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option), Over: yyDollar[5].token.Literal, AnalyticClause: yyDollar[7].expression.(AnalyticClause)}
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:929
		{
			yyVAL.expression = AnalyticClause{Partition: yyDollar[1].expression, OrderByClause: yyDollar[2].expression}
		}
	case 158:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:935
		{
			yyVAL.expression = nil
		}
	case 159:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:939
		{
			yyVAL.expression = Partition{PartitionBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Values: yyDollar[3].expressions}
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:945
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 161:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:949
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:953
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:959
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 164:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:963
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:969
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:973
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 167:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:977
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:981
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:985
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 170:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:989
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 171:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:995
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 172:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:999
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 173:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1003
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 174:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1007
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 175:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1011
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 176:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1017
		{
			yyVAL.expression = nil
		}
	case 177:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1021
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 178:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1025
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1031
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 180:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1035
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 181:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1039
		{
			yyVAL.expression = AllColumns{}
		}
	case 182:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1045
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 183:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1049
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 184:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1055
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 185:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1061
		{
			yyVAL.expression = nil
		}
	case 186:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1065
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 187:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1071
		{
			yyVAL.expression = nil
		}
	case 188:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1075
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1081
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 190:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1085
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1091
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 192:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1095
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 193:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1101
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 194:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1105
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1111
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 196:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1115
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 197:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1121
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 198:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1125
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 199:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1131
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 200:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1135
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 201:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1141
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 202:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1145
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 203:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1151
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Values: yyDollar[5].token.Literal, ValuesList: yyDollar[6].expressions}
		}
	case 204:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:1155
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Values: yyDollar[8].token.Literal, ValuesList: yyDollar[9].expressions}
		}
	case 205:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1159
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 206:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1163
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 207:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1169
		{
			yyVAL.expression = UpdateQuery{WithClause: yyDollar[1].expression, Update: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, Set: yyDollar[4].token.Literal, SetList: yyDollar[5].expressions, FromClause: yyDollar[6].expression, WhereClause: yyDollar[7].expression}
		}
	case 208:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1175
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 209:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1181
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 210:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1185
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 211:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1191
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{WithClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 212:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1196
		{
			from := FromClause{From: yyDollar[4].token.Literal, Tables: yyDollar[5].expressions}
			yyVAL.expression = DeleteQuery{WithClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, FromClause: from, WhereClause: yyDollar[6].expression}
		}
	case 213:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1203
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 214:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1209
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 215:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1213
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1219
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 217:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1223
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1229
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 219:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1233
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 220:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1239
		{
			yyVAL.expression = nil
		}
	case 221:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1243
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1247
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 223:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1251
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 224:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1255
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 225:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1261
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 226:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1265
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 227:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1271
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 228:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1277
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 229:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1281
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 230:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1287
		{
			yyVAL.procexpr = nil
		}
	case 231:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1291
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 232:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1297
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 233:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1301
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 234:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1307
		{
			yyVAL.procexpr = nil
		}
	case 235:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1311
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 236:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1317
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 237:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1323
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 238:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1329
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 239:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1333
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 240:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1340
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 241:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1344
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 242:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1351
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 243:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1357
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 244:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1363
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 245:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1369
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 246:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1375
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 247:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1379
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 248:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1385
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 249:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1391
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 250:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1395
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 251:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1401
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 252:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1405
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 253:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1411
		{
			yyVAL.token = Token{}
		}
	case 254:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1415
		{
			yyVAL.token = yyDollar[1].token
		}
	case 255:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1421
		{
			yyVAL.token = Token{}
		}
	case 256:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1425
		{
			yyVAL.token = yyDollar[1].token
		}
	case 257:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1431
		{
			yyVAL.token = Token{}
		}
	case 258:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1435
		{
			yyVAL.token = yyDollar[1].token
		}
	case 259:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1441
		{
			yyVAL.token = Token{}
		}
	case 260:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1445
		{
			yyVAL.token = yyDollar[1].token
		}
	case 261:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1451
		{
			yyVAL.token = Token{}
		}
	case 262:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1455
		{
			yyVAL.token = yyDollar[1].token
		}
	case 263:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1459
		{
			yyVAL.token = yyDollar[1].token
		}
	case 264:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1463
		{
			yyVAL.token = yyDollar[1].token
		}
	case 265:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1469
		{
			yyVAL.token = Token{}
		}
	case 266:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1473
		{
			yyVAL.token = yyDollar[1].token
		}
	case 267:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1479
		{
			yyVAL.token = Token{}
		}
	case 268:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1483
		{
			yyVAL.token = yyDollar[1].token
		}
	case 269:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1490
		{
			yyVAL.token = yyDollar[1].token
		}
	case 270:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1494
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 271:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1500
		{
			yyVAL.token = Token{}
		}
	case 272:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1504
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
