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
const VAR = 57447
const COMPARISON_OP = 57448
const STRING_OP = 57449
const SUBSTITUTION_OP = 57450

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

//line parser.y:1537

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
	64, 261,
	68, 261,
	69, 261,
	-2, 277,
	-1, 150,
	4, 38,
	-2, 261,
	-1, 151,
	4, 39,
	-2, 261,
	-1, 152,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 172,
	79, 192,
	-2, 261,
	-1, 182,
	83, 3,
	-2, 78,
	-1, 199,
	48, 263,
	50, 267,
	-2, 199,
	-1, 217,
	64, 261,
	68, 261,
	69, 261,
	-2, 185,
	-1, 227,
	64, 261,
	68, 261,
	69, 261,
	-2, 256,
	-1, 234,
	70, 0,
	106, 0,
	109, 0,
	-2, 123,
	-1, 235,
	70, 0,
	106, 0,
	109, 0,
	-2, 125,
	-1, 256,
	96, 69,
	116, 197,
	-2, 261,
	-1, 271,
	77, 3,
	81, 3,
	83, 3,
	-2, 78,
	-1, 286,
	64, 261,
	68, 261,
	69, 261,
	-2, 74,
	-1, 290,
	64, 261,
	68, 261,
	69, 261,
	-2, 114,
	-1, 303,
	50, 267,
	-2, 263,
	-1, 316,
	64, 261,
	68, 261,
	69, 261,
	-2, 64,
	-1, 338,
	83, 1,
	-2, 78,
	-1, 344,
	70, 0,
	106, 0,
	109, 0,
	-2, 134,
	-1, 348,
	64, 261,
	68, 261,
	69, 261,
	-2, 197,
	-1, 354,
	96, 69,
	116, 156,
	-2, 261,
	-1, 372,
	83, 3,
	-2, 78,
	-1, 376,
	64, 261,
	68, 261,
	69, 261,
	-2, 77,
	-1, 434,
	83, 194,
	-2, 261,
	-1, 443,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 456,
	64, 261,
	68, 261,
	69, 261,
	-2, 214,
	-1, 462,
	64, 261,
	68, 261,
	69, 261,
	-2, 68,
	-1, 471,
	64, 261,
	68, 261,
	69, 261,
	-2, 223,
	-1, 477,
	77, 1,
	81, 1,
	83, 1,
	-2, 78,
	-1, 483,
	79, 207,
	81, 207,
	83, 207,
	-2, 261,
	-1, 491,
	77, 1,
	81, 1,
	83, 1,
	-2, 19,
	-1, 522,
	83, 3,
	-2, 78,
	-1, 527,
	64, 261,
	68, 261,
	69, 261,
	-2, 183,
	-1, 553,
	77, 3,
	81, 3,
	83, 3,
	-2, 78,
}

const yyPrivate = 57344

const yyLast = 1206

var yyAct = [...]int{

	290, 38, 531, 409, 502, 520, 40, 41, 42, 43,
	44, 45, 46, 324, 122, 414, 88, 21, 168, 21,
	270, 2, 384, 61, 62, 63, 541, 287, 199, 386,
	76, 111, 394, 114, 116, 304, 377, 336, 207, 214,
	295, 138, 17, 291, 17, 130, 198, 302, 157, 112,
	97, 425, 95, 52, 200, 356, 415, 39, 37, 1,
	141, 125, 78, 210, 117, 58, 169, 470, 146, 147,
	148, 318, 453, 451, 408, 150, 151, 87, 35, 323,
	35, 557, 169, 165, 164, 166, 389, 380, 156, 47,
	307, 167, 308, 309, 310, 305, 170, 121, 303, 169,
	172, 65, 184, 178, 126, 126, 184, 321, 113, 195,
	153, 187, 129, 159, 160, 161, 162, 163, 186, 142,
	203, 205, 189, 555, 154, 153, 554, 155, 159, 160,
	161, 162, 163, 547, 53, 292, 217, 39, 159, 160,
	161, 162, 163, 538, 227, 537, 517, 516, 493, 175,
	3, 487, 486, 485, 233, 234, 235, 65, 306, 472,
	242, 243, 244, 245, 246, 247, 248, 469, 465, 21,
	452, 256, 165, 164, 166, 91, 446, 156, 209, 218,
	420, 407, 34, 350, 349, 65, 251, 553, 34, 254,
	282, 221, 286, 534, 17, 528, 262, 212, 213, 21,
	64, 66, 67, 68, 272, 236, 137, 525, 226, 316,
	47, 232, 55, 154, 153, 230, 155, 159, 160, 161,
	162, 163, 94, 262, 17, 127, 127, 259, 508, 293,
	35, 374, 139, 140, 278, 364, 55, 301, 342, 362,
	344, 137, 360, 333, 229, 149, 222, 55, 280, 390,
	126, 348, 281, 313, 77, 185, 121, 299, 354, 166,
	35, 358, 279, 348, 317, 143, 319, 320, 330, 347,
	53, 352, 49, 158, 50, 369, 48, 370, 371, 179,
	326, 335, 373, 283, 337, 509, 376, 341, 21, 340,
	105, 107, 368, 272, 464, 188, 327, 375, 173, 93,
	193, 174, 533, 196, 127, 482, 430, 127, 488, 279,
	474, 220, 228, 17, 161, 162, 163, 552, 366, 441,
	539, 145, 217, 492, 256, 476, 104, 105, 107, 181,
	108, 109, 433, 423, 419, 367, 421, 422, 372, 393,
	427, 388, 406, 440, 294, 252, 523, 260, 524, 35,
	522, 392, 399, 397, 260, 21, 357, 100, 434, 277,
	417, 403, 339, 559, 442, 218, 338, 523, 220, 432,
	298, 127, 55, 300, 176, 424, 144, 311, 339, 332,
	17, 334, 127, 437, 551, 438, 514, 439, 475, 21,
	499, 456, 110, 444, 272, 106, 34, 426, 325, 328,
	298, 298, 462, 348, 379, 325, 180, 136, 449, 450,
	108, 269, 166, 137, 17, 447, 35, 471, 436, 238,
	457, 463, 466, 237, 239, 461, 459, 454, 183, 478,
	455, 266, 106, 268, 267, 265, 483, 241, 240, 132,
	133, 134, 133, 34, 489, 211, 395, 506, 460, 297,
	35, 307, 458, 308, 309, 310, 396, 479, 81, 325,
	21, 391, 448, 285, 490, 491, 382, 383, 549, 298,
	191, 481, 110, 510, 498, 429, 54, 513, 402, 329,
	331, 192, 127, 501, 401, 17, 511, 505, 398, 507,
	315, 123, 223, 224, 21, 418, 416, 220, 404, 496,
	497, 225, 232, 328, 527, 57, 298, 56, 21, 518,
	120, 137, 530, 137, 535, 137, 536, 521, 494, 17,
	65, 35, 526, 65, 543, 445, 544, 540, 322, 231,
	124, 261, 264, 17, 194, 529, 515, 550, 312, 21,
	208, 197, 546, 542, 272, 131, 206, 548, 387, 115,
	348, 36, 65, 558, 60, 35, 220, 72, 73, 519,
	177, 561, 128, 480, 17, 298, 65, 127, 556, 35,
	21, 119, 127, 135, 560, 272, 410, 411, 412, 413,
	351, 59, 96, 92, 325, 387, 10, 204, 298, 298,
	204, 9, 8, 7, 473, 17, 343, 54, 345, 346,
	35, 6, 296, 500, 5, 220, 4, 165, 164, 166,
	355, 171, 156, 84, 215, 216, 202, 70, 71, 74,
	75, 201, 532, 359, 98, 512, 255, 83, 82, 298,
	86, 35, 263, 263, 127, 307, 127, 308, 309, 310,
	305, 503, 504, 303, 297, 79, 328, 85, 154, 153,
	385, 155, 159, 160, 161, 162, 163, 220, 307, 250,
	308, 309, 310, 305, 80, 204, 303, 467, 468, 54,
	166, 54, 54, 156, 495, 381, 545, 289, 65, 104,
	105, 107, 127, 108, 109, 36, 288, 118, 284, 65,
	104, 105, 107, 190, 108, 109, 36, 263, 400, 263,
	263, 314, 51, 16, 325, 273, 15, 89, 387, 154,
	153, 69, 155, 159, 160, 161, 162, 163, 385, 14,
	385, 13, 385, 12, 263, 361, 363, 365, 11, 271,
	165, 164, 166, 0, 0, 156, 0, 101, 0, 0,
	0, 102, 0, 0, 0, 110, 0, 257, 101, 99,
	0, 263, 102, 0, 0, 0, 110, 0, 257, 103,
	99, 0, 0, 0, 0, 204, 0, 0, 0, 0,
	103, 154, 153, 0, 155, 159, 160, 161, 162, 163,
	0, 249, 250, 0, 0, 106, 258, 0, 0, 90,
	405, 0, 0, 0, 484, 0, 106, 258, 0, 0,
	90, 253, 0, 0, 0, 0, 0, 0, 0, 385,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 263,
	0, 263, 0, 263, 165, 164, 166, 0, 0, 156,
	65, 104, 105, 107, 0, 108, 109, 36, 0, 0,
	0, 65, 104, 105, 107, 0, 108, 109, 36, 0,
	204, 0, 0, 0, 0, 204, 0, 385, 0, 0,
	0, 0, 0, 0, 0, 154, 153, 0, 155, 159,
	160, 161, 162, 163, 0, 0, 250, 0, 0, 0,
	65, 104, 105, 107, 0, 108, 109, 36, 0, 101,
	0, 0, 0, 102, 0, 263, 0, 110, 0, 0,
	101, 99, 0, 0, 102, 0, 0, 0, 110, 0,
	263, 103, 99, 378, 0, 0, 0, 204, 0, 204,
	0, 0, 103, 65, 104, 105, 107, 0, 108, 109,
	36, 165, 164, 166, 0, 0, 156, 106, 219, 101,
	379, 90, 0, 102, 0, 0, 0, 110, 106, 353,
	34, 99, 90, 0, 165, 164, 166, 0, 263, 156,
	0, 103, 0, 0, 0, 204, 165, 164, 166, 477,
	0, 156, 154, 153, 0, 155, 159, 160, 161, 162,
	163, 443, 101, 0, 0, 0, 102, 106, 0, 0,
	110, 90, 0, 0, 99, 154, 153, 0, 155, 159,
	160, 161, 162, 163, 103, 0, 0, 154, 153, 431,
	155, 159, 160, 161, 162, 163, 165, 164, 166, 0,
	0, 156, 0, 0, 0, 0, 0, 165, 164, 166,
	106, 435, 156, 0, 90, 0, 165, 164, 166, 0,
	0, 156, 0, 0, 182, 0, 0, 0, 0, 0,
	0, 165, 164, 166, 0, 0, 156, 154, 153, 0,
	155, 159, 160, 161, 162, 163, 152, 0, 154, 153,
	0, 155, 159, 160, 161, 162, 163, 154, 153, 0,
	155, 159, 160, 161, 162, 163, 165, 164, 166, 0,
	0, 156, 154, 153, 0, 155, 159, 160, 161, 162,
	163, 428, 164, 166, 36, 0, 156, 0, 0, 32,
	165, 36, 166, 0, 0, 156, 32, 0, 0, 18,
	0, 0, 19, 0, 0, 0, 18, 154, 153, 19,
	155, 159, 160, 161, 162, 163, 0, 0, 0, 0,
	0, 0, 154, 153, 0, 155, 159, 160, 161, 162,
	163, 154, 153, 0, 155, 159, 160, 161, 162, 163,
	0, 0, 0, 0, 0, 0, 0, 34, 0, 274,
	0, 30, 0, 0, 34, 0, 29, 24, 30, 0,
	28, 25, 26, 27, 24, 0, 0, 28, 25, 26,
	27, 0, 22, 23, 275, 276, 31, 33, 20, 22,
	23, 0, 0, 31, 33, 20,
}
var yyPact = [...]int{

	1100, -1000, 1100, -62, -62, -62, -62, -62, -62, -62,
	-62, -1000, -1000, -1000, -1000, -1000, -1000, 257, 477, 475,
	543, -62, -62, -62, 562, 562, 562, 562, 526, 919,
	919, -62, 537, 919, 485, 148, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 453, 510, 562,
	548, 531, 381, 334, -1000, 322, 562, 562, -62, 2,
	157, -1000, -1000, -1000, 291, -1000, -62, -62, -62, 562,
	-1000, -1000, -1000, -1000, 919, 919, 986, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 148, -1000, -1000,
	876, -19, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 919,
	192, 132, 919, 562, -1000, -1000, 284, -1000, -1000, -1000,
	-1000, 962, 364, -11, -1000, 146, 18, -1000, -6, 562,
	-1000, 919, 426, 440, 562, 518, -8, 519, 97, 532,
	522, 97, 384, 384, 384, 826, -1000, 75, 121, 131,
	465, -1000, 543, 919, 226, 129, -1000, -1000, -1000, 509,
	1021, 1021, 1100, 919, 919, 919, 345, 355, 376, 919,
	919, 919, 919, 919, 919, 919, -1000, 665, 70, 562,
	685, 268, 1021, 81, 81, 367, 372, -1000, 603, 341,
	-1000, -1000, 1093, 562, 540, 321, -1000, 485, 168, 1021,
	418, 919, 919, 114, 562, 562, -1000, 562, 522, 41,
	-1000, 516, -1000, -1000, -1000, -1000, 97, 451, 919, -1000,
	121, -1000, 121, 121, -1000, -10, 506, 1021, -1000, -1000,
	-36, -1000, 562, 181, 153, 562, -1000, 1021, 322, 562,
	322, 540, 285, 28, 3, 3, 401, 919, 81, 919,
	81, 81, 202, 202, -1000, -1000, -1000, 1045, 603, -1000,
	919, -1000, -1000, -1000, 68, 67, 542, 837, -1000, 275,
	919, -1000, 876, -1000, -1000, 81, 127, 124, 120, 345,
	252, 1093, -1000, -1000, 919, -62, -62, 256, -1000, -15,
	-62, -1000, 116, 562, -1000, 919, 866, -1000, -30, 424,
	1021, -1000, 81, 562, -1000, 531, -31, 140, -52, -1000,
	-1000, -1000, 413, 402, 396, 408, 97, -1000, -1000, -1000,
	-1000, -1000, 562, 522, 444, 437, 1021, 383, -1000, -1000,
	383, 826, 562, 674, 65, -43, 545, 562, 461, -1000,
	562, 458, -62, 64, -62, -62, 250, 285, 1100, 919,
	-1000, -1000, 1036, -1000, 3, -1000, -1000, -1000, 759, -1000,
	-1000, 434, 210, -1000, 971, 249, 268, 919, 951, 353,
	108, -1000, 108, -1000, 108, -1000, 255, 286, -1000, 901,
	-1000, -1000, 1093, -1000, 322, 60, 1021, -1000, 330, 416,
	919, 337, -1000, -1000, -1000, -44, 54, -45, 522, 562,
	919, 97, 404, 396, 400, -1000, 97, -1000, -1000, -1000,
	-1000, 919, 919, -1000, -1000, 196, 52, -1000, 562, -1000,
	-1000, -1000, 562, 562, 51, -50, 919, 43, 562, -1000,
	224, -1000, -1000, 312, 242, 301, -1000, 889, 919, 919,
	558, 430, 209, -1000, 1021, 919, 81, 37, 36, 35,
	-1000, 213, -62, 1093, 240, 32, 496, -1000, -1000, -1000,
	468, 81, 369, 562, -1000, -1000, 1021, 586, 97, 399,
	97, 609, 1021, -1000, 113, 187, -1000, -1000, -1000, 545,
	562, 1021, -1000, -1000, 322, -62, 310, 1100, 603, 31,
	30, 919, 554, 1021, -1000, -1000, -1000, -1000, -1000, -1000,
	269, 1100, 270, -1000, 92, -1000, -1000, -1000, -1000, 81,
	-1000, -1000, -1000, 919, 80, 609, 97, 586, 205, 78,
	-1000, -1000, -62, -1000, -62, -1000, -1000, -1000, 29, 27,
	237, 269, 1093, 919, -62, 322, -1000, 1021, 562, 609,
	-1000, 17, 453, 427, 205, -1000, -1000, -1000, -1000, 308,
	234, 290, -1000, 107, -1000, 10, 7, -1000, -1000, 919,
	-35, -62, 287, 1093, -1000, -1000, -1000, -1000, -1000, -62,
	-1000, -1000,
}
var yyPgo = [...]int{

	0, 58, 20, 21, 729, 728, 723, 721, 719, 711,
	707, 706, 705, 703, 150, 71, 53, 702, 45, 38,
	701, 698, 14, 693, 36, 688, 41, 687, 64, 62,
	254, 0, 357, 22, 27, 686, 677, 675, 674, 458,
	664, 647, 645, 630, 628, 627, 626, 624, 43, 2,
	622, 54, 621, 28, 616, 4, 615, 614, 613, 611,
	610, 29, 18, 46, 61, 13, 39, 55, 606, 604,
	602, 40, 601, 593, 592, 56, 15, 3, 591, 586,
	51, 37, 26, 5, 175, 583, 299, 222, 52, 582,
	50, 77, 49, 16, 581, 65, 573, 48, 47, 32,
	35, 63, 571, 273, 1,
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
	31, 31, 31, 31, 31, 31, 32, 32, 33, 33,
	34, 34, 35, 35, 36, 36, 37, 37, 37, 38,
	38, 39, 40, 41, 41, 41, 41, 41, 41, 41,
	41, 41, 41, 41, 41, 41, 41, 41, 41, 41,
	41, 41, 42, 42, 42, 42, 42, 43, 43, 43,
	44, 44, 45, 45, 46, 46, 46, 47, 47, 47,
	47, 48, 48, 49, 50, 50, 51, 51, 51, 52,
	52, 53, 53, 53, 53, 53, 53, 54, 54, 54,
	54, 54, 55, 55, 55, 56, 56, 56, 57, 57,
	58, 59, 59, 60, 60, 61, 61, 62, 62, 63,
	63, 64, 64, 65, 65, 66, 66, 67, 67, 68,
	68, 68, 68, 69, 70, 71, 71, 72, 72, 73,
	74, 74, 75, 75, 76, 76, 77, 77, 77, 77,
	77, 78, 78, 79, 80, 80, 81, 81, 82, 82,
	83, 83, 84, 85, 86, 86, 87, 87, 88, 89,
	90, 91, 92, 92, 93, 94, 94, 95, 95, 96,
	96, 97, 97, 98, 98, 99, 99, 100, 100, 100,
	100, 101, 101, 102, 102, 103, 103, 104, 104,
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
	1, 1, 1, 1, 1, 3, 3, 1, 1, 3,
	1, 3, 2, 4, 1, 1, 0, 1, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 4, 4, 6,
	6, 4, 6, 4, 4, 4, 6, 4, 4, 6,
	4, 2, 3, 3, 3, 3, 3, 3, 3, 2,
	3, 4, 4, 1, 1, 2, 2, 7, 8, 7,
	8, 7, 8, 2, 0, 3, 1, 2, 3, 1,
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

	-1000, -1, -3, -14, -68, -69, -72, -73, -74, -78,
	-79, -5, -6, -7, -8, -11, -13, -26, 26, 29,
	105, -93, 99, 100, 84, 88, 89, 90, 87, 76,
	78, 103, 16, 104, 74, -91, 11, -1, -104, 119,
	-104, -104, -104, -104, -104, -104, -104, -15, 19, 15,
	17, -17, -16, 13, -39, 115, 30, 30, -95, -94,
	11, -104, -104, -104, -84, 4, -84, -84, -84, -9,
	91, 92, 31, 32, 93, 94, -31, -30, -29, -42,
	-40, -39, -44, -45, -58, -41, -43, -91, -93, -10,
	115, -84, -85, -86, -87, -88, -89, -90, -47, 75,
	-32, 63, 67, 85, 5, 6, 111, 7, 9, 10,
	71, -31, -92, -91, -104, 12, -31, -28, -27, -102,
	25, 108, -22, 38, 20, -64, -51, -84, 14, -64,
	-18, 14, 58, 59, 60, -96, 73, -14, -26, -84,
	-84, -104, 117, 108, 85, 30, -104, -104, -104, -84,
	-31, -31, 80, 107, 106, 109, 70, -97, -103, 110,
	111, 112, 113, 114, 66, 65, 67, -31, -62, 118,
	115, -59, -31, 106, 109, -97, -103, -39, -31, -84,
	-86, -87, 82, 64, 117, 109, -104, 117, -84, -31,
	-23, 44, 41, -84, 16, 117, -84, 22, -63, -53,
	-51, -52, -54, 23, -39, 24, 14, -19, 18, -63,
	-101, 61, -101, -101, -66, -57, -56, -31, -48, 112,
	-84, 116, 115, 27, 28, 36, -95, -31, 86, 115,
	86, 20, -1, -31, -31, -31, -97, 68, 64, 69,
	62, 61, -31, -31, -31, -31, -31, -31, -31, 116,
	117, 116, -84, 116, -62, -46, -31, 73, 112, -67,
	79, -32, 115, -39, -32, 68, 64, 62, 61, 70,
	-2, -4, -3, -12, 76, 101, 102, -84, -92, -91,
	-29, -28, 22, 115, -25, 45, -31, -34, -35, -36,
	-31, -48, 21, 115, -14, -71, -70, -30, -84, -64,
	-84, -19, -98, 57, -100, 54, 117, 49, 51, 52,
	53, -84, 22, -63, -20, 39, -31, -16, -15, -16,
	-16, 117, 22, 115, -65, -84, -75, 115, -84, -30,
	115, -30, -14, -65, -14, -92, -81, -80, 81, 77,
	-88, -90, -31, -32, -31, -32, -32, -62, -31, 116,
	116, 38, -22, 112, -31, -60, -67, 81, -31, -32,
	115, -39, 115, -39, 115, -39, -97, 83, -2, -31,
	-104, -104, 82, -104, 115, -65, -31, -24, 47, 74,
	117, -37, 42, 43, -33, -32, -61, -30, -18, 117,
	109, 48, -98, -100, -99, 50, 48, -63, -84, -19,
	-21, 40, 41, -66, -84, 116, -62, 116, 117, -77,
	31, 32, 33, 34, -76, -75, 35, -61, 37, -104,
	116, -104, -104, 83, -81, -80, -1, -31, 65, 41,
	96, 38, -22, 83, -31, 80, 65, -33, -33, -33,
	88, 64, 78, 80, -2, -14, 116, -24, 46, -34,
	72, 117, 116, 117, -19, -71, -31, -53, 48, -99,
	48, -53, -31, -62, 98, 116, -65, -30, -30, 116,
	117, -31, 116, -84, 86, 76, 83, 80, -31, -34,
	5, 41, 96, -31, -32, 116, 116, 116, 95, -104,
	-2, -3, 83, 116, 22, -38, 31, 32, -33, 21,
	-14, -61, -55, 55, 56, -53, 48, -53, 115, 98,
	-77, -76, -14, -104, 76, -1, 116, 116, -34, 5,
	-83, -82, 81, 77, 78, 115, -33, -31, 115, -53,
	-55, -49, -50, 97, 115, -104, -104, 116, 116, 83,
	-83, -82, -2, -31, -104, -14, -65, 116, -22, 41,
	-49, 76, 83, 80, 116, 116, -62, 116, -104, 76,
	-2, -104,
}
var yyDef = [...]int{

	-2, -2, -2, 277, 277, 277, 277, 277, 277, 277,
	277, 13, 14, 15, 16, 17, 18, 0, 0, 0,
	0, 277, 277, 277, 0, 0, 0, 0, 33, 0,
	0, 277, 0, 0, 273, 0, 251, 2, 5, 278,
	6, 7, 8, 9, 10, 11, 12, -2, 0, 0,
	0, 61, 0, 259, 59, 78, 0, 0, 277, 257,
	255, 22, 23, 24, 0, 242, 277, 277, 277, 0,
	34, 35, 36, 37, 0, 0, 261, 92, 93, 94,
	95, 96, 97, 98, 99, 100, 101, 102, 103, 104,
	78, 90, 84, 85, 86, 87, 88, 89, 153, 191,
	261, 0, 0, 0, 243, 244, 0, 246, 248, 249,
	250, 261, 0, 102, 46, 0, -2, 79, 82, 0,
	274, 0, 71, 0, 0, 0, 201, 166, 0, 0,
	63, 0, 271, 271, 271, 0, 260, 0, 0, 0,
	0, 21, 0, 0, 0, 0, 26, 27, 28, 0,
	-2, -2, -2, 0, 275, 276, 261, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 262, 261, 0, 0,
	0, 0, -2, 275, 276, 0, 0, 141, 149, 0,
	245, 247, -2, 0, 0, 0, 52, 273, 0, 254,
	76, 0, 0, 78, 0, 0, 167, 0, 63, -2,
	171, 172, 175, 176, 169, 170, 0, 65, 0, 62,
	0, 272, 0, 0, 60, 205, 188, -2, 186, 187,
	90, 121, 0, 0, 0, 0, 258, -2, 78, 0,
	78, 0, 236, 122, -2, -2, 0, 0, 0, 0,
	0, 0, 142, 143, 144, 145, 146, 147, 148, 105,
	0, 106, 91, 150, 0, 0, -2, 0, 154, 193,
	0, 124, 78, 107, 126, 0, 0, 0, 0, 261,
	0, -2, 19, 20, 0, 277, 277, 0, 253, 252,
	277, 83, 0, 0, 53, 0, -2, 70, 110, 116,
	-2, 115, 0, 0, 211, 61, 215, 0, 90, 202,
	168, 217, 0, -2, 265, 0, 0, 264, 268, 269,
	270, 173, 0, 63, 67, 0, -2, 55, 58, 56,
	57, 0, 0, 0, 0, 203, 226, 0, 222, 231,
	0, 0, 277, 0, 277, 277, 0, 236, -2, 0,
	127, 128, 261, 131, -2, 135, 138, 198, -2, 151,
	152, 0, 0, 155, -2, 0, 208, 0, 261, 0,
	78, 133, 78, 137, 78, 140, 0, 0, 4, 261,
	49, 50, -2, 51, 78, 0, -2, 72, 74, 0,
	0, 112, 117, 118, 209, 108, 0, 195, 63, 0,
	0, 0, 0, 265, 0, 266, 0, 200, 174, 218,
	54, 0, 0, 206, 189, 150, 0, 219, 0, 220,
	227, 228, 0, 0, 0, 224, 0, 0, 0, 25,
	30, 32, 29, 0, 0, 235, 237, 261, 0, 0,
	0, 0, 0, 190, -2, 0, 0, 0, 0, 0,
	40, 0, 277, -2, 0, 0, 0, 73, 75, 111,
	0, 0, 78, 0, 213, 216, -2, 182, 0, 0,
	0, 181, -2, 66, 0, 151, 204, 229, 230, 226,
	0, -2, 232, 233, 78, 277, 0, -2, 129, 70,
	0, 0, 0, -2, 130, 132, 136, 139, 41, 44,
	240, -2, 0, 80, 0, 113, 119, 120, 109, 0,
	212, 196, 177, 0, 0, 178, 0, 182, 164, 0,
	221, 225, 277, 42, 277, 234, 157, 159, 70, 0,
	0, 240, -2, 0, 277, 78, 210, -2, 0, 180,
	179, 0, 69, 0, 164, 31, 43, 158, 160, 0,
	0, 239, 241, 261, 45, 0, 0, 161, 163, 0,
	0, 277, 0, -2, 81, 184, 165, 162, 47, 277,
	238, 48,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 114, 3, 3,
	115, 116, 112, 110, 117, 111, 118, 113, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 119,
	3, 109,
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
	102, 103, 104, 105, 106, 107, 108,
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
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:402
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
		//line parser.y:414
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
		//line parser.y:424
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
		//line parser.y:433
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
		//line parser.y:442
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
		//line parser.y:453
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:457
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:463
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:469
		{
			yyVAL.expression = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:473
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
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
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
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
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 67:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:499
		{
			yyVAL.expression = nil
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:503
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 69:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:509
		{
			yyVAL.expression = nil
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:513
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 71:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:519
		{
			yyVAL.expression = nil
		}
	case 72:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:523
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 73:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:527
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
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
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 76:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:543
		{
			yyVAL.expression = nil
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:547
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 78:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:553
		{
			yyVAL.expression = nil
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:557
		{
			yyVAL.expression = WithClause{With: yyDollar[1].token.Literal, InlineTables: yyDollar[2].expressions}
		}
	case 80:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:563
		{
			yyVAL.expression = InlineTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, As: yyDollar[3].token.Literal, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 81:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:567
		{
			yyVAL.expression = InlineTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, Columns: yyDollar[4].expressions, As: yyDollar[6].token.Literal, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:573
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:577
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:583
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:587
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:591
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:595
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:599
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:603
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:609
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:613
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:623
		{
			yyVAL.expression = yyDollar[1].primary
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
			yyVAL.expression = yyDollar[1].expression
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
			yyVAL.expression = yyDollar[1].variable
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
			yyVAL.expression = yyDollar[1].expression
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:671
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:677
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:681
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:687
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:691
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:697
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:701
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:707
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 113:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:711
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:717
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:721
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 116:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:727
		{
			yyVAL.token = Token{}
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:731
		{
			yyVAL.token = yyDollar[1].token
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:735
		{
			yyVAL.token = yyDollar[1].token
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:741
		{
			yyVAL.token = yyDollar[1].token
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:745
		{
			yyVAL.token = yyDollar[1].token
		}
	case 121:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:751
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:757
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
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:780
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:784
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, RHS: yyDollar[3].expression}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:788
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: "=", RHS: yyDollar[3].expression}
		}
	case 127:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:796
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 128:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:800
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 129:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:804
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 130:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:808
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 131:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:812
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 132:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:816
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 133:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:820
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 134:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 135:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:828
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 136:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:832
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 137:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:836
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 138:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:840
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 139:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:844
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 140:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:848
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token.Literal, Values: yyDollar[4].expression}
		}
	case 141:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:852
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 142:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:858
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:862
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 144:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:866
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:870
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 146:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:874
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:880
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 148:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:884
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 149:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:888
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 150:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:894
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal}
		}
	case 151:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:898
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Args: yyDollar[3].expressions}
		}
	case 152:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:904
		{
			yyVAL.expression = AggregateFunction{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(AggregateOption)}
		}
	case 153:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:908
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:914
		{
			yyVAL.expression = AggregateOption{Args: []Expression{AllColumns{}}}
		}
	case 155:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:918
		{
			yyVAL.expression = AggregateOption{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 156:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:922
		{
			yyVAL.expression = AggregateOption{Distinct: yyDollar[1].token, Args: []Expression{yyDollar[2].expression}}
		}
	case 157:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:928
		{
			orderBy := OrderByClause{OrderBy: yyDollar[4].token.Literal + " " + yyDollar[5].token.Literal, Items: yyDollar[6].expressions}
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Args: []Expression{yyDollar[3].expression}}, OrderBy: orderBy}
		}
	case 158:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:933
		{
			orderBy := OrderByClause{OrderBy: yyDollar[5].token.Literal + " " + yyDollar[6].token.Literal, Items: yyDollar[7].expressions}
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Distinct: yyDollar[3].token, Args: []Expression{yyDollar[4].expression}}, OrderBy: orderBy}
		}
	case 159:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:938
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Args: []Expression{yyDollar[3].expression}}, OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 160:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:942
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].identifier.Literal, Option: AggregateOption{Distinct: yyDollar[3].token, Args: []Expression{yyDollar[4].expression}}, OrderBy: yyDollar[5].expression, SeparatorLit: yyDollar[6].token.Literal, Separator: yyDollar[7].token.Literal}
		}
	case 161:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:948
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Over: yyDollar[4].token.Literal, AnalyticClause: yyDollar[6].expression.(AnalyticClause)}
		}
	case 162:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:952
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Args: yyDollar[3].expressions, Over: yyDollar[5].token.Literal, AnalyticClause: yyDollar[7].expression.(AnalyticClause)}
		}
	case 163:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:958
		{
			yyVAL.expression = AnalyticClause{Partition: yyDollar[1].expression, OrderByClause: yyDollar[2].expression}
		}
	case 164:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:964
		{
			yyVAL.expression = nil
		}
	case 165:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:968
		{
			yyVAL.expression = Partition{PartitionBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Values: yyDollar[3].expressions}
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:974
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 167:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:978
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 168:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:982
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 169:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:988
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 170:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:992
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:998
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1002
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 173:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1006
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 174:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1010
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1014
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1018
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 177:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1024
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 178:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1028
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 179:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1032
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 180:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1036
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 181:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1040
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 182:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1046
		{
			yyVAL.expression = nil
		}
	case 183:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1050
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 184:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1054
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1060
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 186:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1064
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 187:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1068
		{
			yyVAL.expression = AllColumns{}
		}
	case 188:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1074
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 189:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1078
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token.Literal, Alias: yyDollar[3].identifier}
		}
	case 190:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1084
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 191:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1090
		{
			yyVAL.expression = nil
		}
	case 192:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1094
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 193:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1100
		{
			yyVAL.expression = nil
		}
	case 194:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1104
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1110
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 196:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1114
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 197:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1120
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 198:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1124
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 199:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1130
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 200:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1134
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 201:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1140
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 202:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1144
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1150
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 204:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1154
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1160
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1164
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 207:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1170
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 208:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1174
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 209:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1180
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Values: yyDollar[5].token.Literal, ValuesList: yyDollar[6].expressions}
		}
	case 210:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:1184
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Values: yyDollar[8].token.Literal, ValuesList: yyDollar[9].expressions}
		}
	case 211:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1188
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 212:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1192
		{
			yyVAL.expression = InsertQuery{WithClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 213:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1198
		{
			yyVAL.expression = UpdateQuery{WithClause: yyDollar[1].expression, Update: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, Set: yyDollar[4].token.Literal, SetList: yyDollar[5].expressions, FromClause: yyDollar[6].expression, WhereClause: yyDollar[7].expression}
		}
	case 214:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1204
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 215:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1210
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 216:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1214
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 217:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1220
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{WithClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 218:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1225
		{
			from := FromClause{From: yyDollar[4].token.Literal, Tables: yyDollar[5].expressions}
			yyVAL.expression = DeleteQuery{WithClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, FromClause: from, WhereClause: yyDollar[6].expression}
		}
	case 219:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1232
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 220:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1238
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 221:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1242
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1248
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 223:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1252
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1258
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 225:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1262
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 226:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1268
		{
			yyVAL.expression = nil
		}
	case 227:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1272
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 228:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1276
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 229:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1280
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 230:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1284
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 231:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1290
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 232:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1294
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 233:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1300
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 234:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1306
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 235:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1310
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 236:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1316
		{
			yyVAL.procexpr = nil
		}
	case 237:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1320
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 238:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1326
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 239:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1330
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 240:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1336
		{
			yyVAL.procexpr = nil
		}
	case 241:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1340
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 242:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1346
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 243:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1352
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 244:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1358
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 245:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1362
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 246:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1369
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 247:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1373
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 248:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1380
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 249:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1386
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 250:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1392
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 251:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1398
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 252:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1404
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 253:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1408
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 254:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1414
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 255:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1420
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 256:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1424
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 257:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1430
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 258:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1434
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 259:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1440
		{
			yyVAL.token = Token{}
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
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1470
		{
			yyVAL.token = Token{}
		}
	case 266:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1474
		{
			yyVAL.token = yyDollar[1].token
		}
	case 267:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1480
		{
			yyVAL.token = Token{}
		}
	case 268:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1484
		{
			yyVAL.token = yyDollar[1].token
		}
	case 269:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1488
		{
			yyVAL.token = yyDollar[1].token
		}
	case 270:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1492
		{
			yyVAL.token = yyDollar[1].token
		}
	case 271:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1498
		{
			yyVAL.token = Token{}
		}
	case 272:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1502
		{
			yyVAL.token = yyDollar[1].token
		}
	case 273:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1508
		{
			yyVAL.token = Token{}
		}
	case 274:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1512
		{
			yyVAL.token = yyDollar[1].token
		}
	case 275:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1519
		{
			yyVAL.token = yyDollar[1].token
		}
	case 276:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1523
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 277:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1529
		{
			yyVAL.token = Token{}
		}
	case 278:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1533
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
