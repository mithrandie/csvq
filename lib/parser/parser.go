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
const GROUP_CONCAT = 57433
const SEPARATOR = 57434
const PARTITION = 57435
const OVER = 57436
const COMMIT = 57437
const ROLLBACK = 57438
const CONTINUE = 57439
const BREAK = 57440
const EXIT = 57441
const PRINT = 57442
const VAR = 57443
const COMPARISON_OP = 57444
const STRING_OP = 57445
const SUBSTITUTION_OP = 57446

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

//line parser.y:1442

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
	-2, 65,
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	1, 1,
	77, 1,
	81, 1,
	83, 1,
	-2, 65,
	-1, 46,
	58, 45,
	59, 45,
	60, 45,
	-2, 56,
	-1, 107,
	64, 241,
	68, 241,
	69, 241,
	-2, 257,
	-1, 140,
	77, 1,
	81, 1,
	83, 1,
	-2, 65,
	-1, 158,
	112, 137,
	-2, 239,
	-1, 160,
	79, 172,
	-2, 241,
	-1, 169,
	38, 137,
	92, 137,
	112, 137,
	-2, 239,
	-1, 170,
	83, 3,
	-2, 65,
	-1, 187,
	48, 243,
	50, 247,
	-2, 179,
	-1, 205,
	64, 241,
	68, 241,
	69, 241,
	-2, 165,
	-1, 215,
	64, 241,
	68, 241,
	69, 241,
	-2, 236,
	-1, 221,
	70, 0,
	102, 0,
	105, 0,
	-2, 108,
	-1, 222,
	70, 0,
	102, 0,
	105, 0,
	-2, 110,
	-1, 254,
	77, 3,
	81, 3,
	83, 3,
	-2, 65,
	-1, 268,
	64, 241,
	68, 241,
	69, 241,
	-2, 61,
	-1, 272,
	64, 241,
	68, 241,
	69, 241,
	-2, 99,
	-1, 285,
	50, 247,
	-2, 243,
	-1, 298,
	64, 241,
	68, 241,
	69, 241,
	-2, 51,
	-1, 305,
	112, 137,
	-2, 239,
	-1, 318,
	83, 1,
	-2, 65,
	-1, 324,
	70, 0,
	102, 0,
	105, 0,
	-2, 119,
	-1, 328,
	64, 241,
	68, 241,
	69, 241,
	-2, 177,
	-1, 349,
	83, 3,
	-2, 65,
	-1, 353,
	64, 241,
	68, 241,
	69, 241,
	-2, 64,
	-1, 403,
	83, 174,
	-2, 241,
	-1, 412,
	77, 1,
	81, 1,
	83, 1,
	-2, 65,
	-1, 425,
	64, 241,
	68, 241,
	69, 241,
	-2, 194,
	-1, 431,
	64, 241,
	68, 241,
	69, 241,
	-2, 55,
	-1, 439,
	64, 241,
	68, 241,
	69, 241,
	-2, 203,
	-1, 444,
	77, 1,
	81, 1,
	83, 1,
	-2, 65,
	-1, 446,
	79, 187,
	81, 187,
	83, 187,
	-2, 241,
	-1, 454,
	77, 1,
	81, 1,
	83, 1,
	-2, 18,
	-1, 480,
	83, 3,
	-2, 65,
	-1, 485,
	64, 241,
	68, 241,
	69, 241,
	-2, 163,
	-1, 504,
	77, 3,
	81, 3,
	83, 3,
	-2, 65,
}

const yyPrivate = 57344

const yyLast = 1077

var yyAct = [...]int{

	37, 113, 306, 465, 156, 39, 40, 41, 42, 43,
	44, 45, 2, 390, 82, 253, 79, 34, 478, 34,
	371, 60, 61, 62, 493, 363, 195, 36, 1, 91,
	277, 105, 80, 20, 316, 20, 206, 269, 63, 65,
	66, 67, 68, 186, 385, 187, 104, 286, 284, 202,
	354, 240, 121, 88, 398, 86, 188, 391, 132, 116,
	361, 108, 103, 118, 118, 71, 136, 137, 138, 333,
	130, 131, 57, 38, 305, 145, 157, 157, 300, 153,
	152, 154, 51, 141, 144, 198, 147, 148, 149, 150,
	151, 438, 422, 112, 16, 46, 289, 420, 290, 291,
	292, 287, 172, 384, 285, 117, 117, 64, 174, 120,
	153, 152, 154, 366, 357, 144, 142, 141, 410, 143,
	147, 148, 149, 150, 151, 176, 191, 193, 158, 38,
	181, 157, 172, 184, 118, 303, 183, 118, 409, 507,
	274, 208, 70, 175, 133, 506, 505, 142, 141, 129,
	143, 147, 148, 149, 150, 151, 218, 34, 237, 33,
	288, 147, 148, 149, 150, 151, 197, 163, 219, 477,
	456, 450, 239, 20, 449, 52, 129, 48, 448, 49,
	440, 47, 437, 255, 433, 421, 260, 34, 415, 218,
	383, 244, 247, 33, 329, 208, 245, 280, 118, 64,
	282, 238, 217, 20, 293, 64, 214, 52, 46, 118,
	200, 201, 264, 283, 54, 367, 209, 273, 315, 489,
	223, 252, 149, 150, 151, 307, 310, 280, 280, 242,
	275, 486, 245, 483, 351, 261, 341, 263, 295, 262,
	117, 339, 327, 281, 337, 210, 331, 169, 54, 95,
	97, 153, 152, 154, 343, 323, 144, 325, 326, 347,
	348, 173, 112, 350, 134, 471, 501, 255, 352, 308,
	345, 34, 154, 54, 317, 216, 129, 321, 336, 320,
	307, 299, 85, 301, 302, 146, 503, 20, 142, 141,
	280, 143, 147, 148, 149, 150, 151, 135, 236, 237,
	491, 265, 455, 118, 362, 54, 312, 161, 3, 375,
	162, 129, 309, 84, 328, 395, 443, 402, 208, 381,
	396, 344, 376, 349, 310, 279, 243, 280, 334, 243,
	365, 482, 374, 370, 369, 34, 411, 481, 393, 481,
	129, 480, 319, 69, 102, 319, 399, 107, 511, 318,
	96, 20, 397, 380, 502, 311, 313, 382, 475, 442,
	33, 356, 255, 128, 462, 413, 34, 362, 127, 362,
	419, 362, 208, 154, 249, 171, 98, 164, 248, 168,
	405, 280, 20, 118, 432, 251, 250, 434, 118, 199,
	128, 428, 423, 124, 273, 418, 155, 424, 406, 307,
	407, 241, 408, 280, 280, 160, 416, 372, 166, 441,
	167, 225, 452, 74, 426, 224, 226, 33, 364, 430,
	228, 227, 123, 124, 125, 454, 469, 177, 453, 34,
	53, 429, 129, 427, 129, 447, 129, 280, 100, 373,
	219, 205, 118, 474, 118, 20, 129, 368, 464, 215,
	362, 417, 473, 310, 126, 364, 220, 221, 222, 267,
	179, 34, 229, 230, 231, 232, 233, 234, 235, 359,
	360, 34, 476, 468, 488, 470, 490, 20, 479, 509,
	379, 461, 472, 496, 118, 180, 378, 20, 297, 498,
	276, 114, 362, 255, 268, 272, 494, 34, 492, 211,
	212, 307, 508, 510, 394, 392, 165, 56, 213, 279,
	55, 298, 514, 20, 513, 487, 129, 255, 459, 460,
	512, 34, 64, 484, 289, 314, 290, 291, 292, 111,
	457, 435, 436, 192, 304, 139, 192, 20, 196, 322,
	294, 324, 64, 53, 355, 386, 387, 388, 389, 115,
	64, 94, 95, 97, 128, 98, 99, 35, 335, 182,
	185, 122, 153, 152, 154, 364, 194, 144, 106, 64,
	35, 356, 346, 59, 451, 246, 246, 64, 129, 119,
	110, 58, 353, 87, 154, 83, 10, 144, 64, 94,
	95, 97, 9, 98, 99, 35, 8, 7, 6, 142,
	141, 278, 143, 147, 148, 149, 150, 151, 192, 92,
	5, 4, 53, 93, 53, 53, 332, 100, 205, 142,
	141, 90, 143, 147, 148, 149, 150, 151, 159, 76,
	203, 94, 95, 97, 400, 98, 99, 101, 204, 246,
	190, 246, 246, 189, 500, 499, 128, 92, 128, 403,
	128, 93, 89, 96, 207, 100, 75, 81, 33, 90,
	414, 78, 246, 338, 340, 342, 289, 72, 290, 291,
	292, 287, 272, 77, 285, 101, 73, 458, 358, 271,
	270, 109, 425, 266, 153, 152, 154, 178, 246, 144,
	377, 96, 296, 431, 50, 81, 15, 100, 64, 94,
	95, 97, 192, 98, 99, 35, 256, 439, 64, 94,
	95, 97, 14, 98, 99, 35, 445, 13, 12, 446,
	11, 142, 141, 254, 143, 147, 148, 149, 150, 151,
	463, 0, 289, 96, 290, 291, 292, 287, 466, 467,
	285, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 246, 0, 246, 0, 246, 0, 92, 0, 0,
	0, 93, 0, 0, 0, 100, 0, 92, 0, 90,
	0, 93, 0, 0, 0, 100, 0, 0, 0, 90,
	0, 485, 192, 0, 0, 101, 0, 192, 153, 152,
	154, 0, 497, 144, 0, 101, 495, 0, 0, 0,
	0, 96, 330, 504, 0, 81, 0, 153, 152, 154,
	0, 96, 144, 0, 0, 81, 0, 0, 0, 246,
	0, 0, 444, 0, 0, 142, 141, 0, 143, 147,
	148, 149, 150, 151, 246, 0, 0, 0, 0, 0,
	0, 192, 0, 192, 142, 141, 0, 143, 147, 148,
	149, 150, 151, 153, 152, 154, 0, 0, 144, 0,
	0, 0, 153, 152, 154, 0, 0, 144, 412, 0,
	0, 0, 0, 153, 152, 154, 246, 404, 144, 0,
	0, 0, 0, 192, 153, 152, 154, 0, 140, 144,
	142, 141, 0, 143, 147, 148, 149, 150, 151, 142,
	141, 170, 143, 147, 148, 149, 150, 151, 0, 0,
	142, 141, 0, 143, 147, 148, 149, 150, 151, 0,
	0, 142, 141, 0, 143, 147, 148, 149, 150, 151,
	401, 152, 154, 35, 0, 144, 0, 0, 31, 153,
	0, 154, 0, 0, 144, 0, 0, 0, 17, 0,
	0, 18, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 142, 141, 0,
	143, 147, 148, 149, 150, 151, 142, 141, 0, 143,
	147, 148, 149, 150, 151, 0, 35, 0, 0, 0,
	0, 31, 0, 0, 0, 0, 33, 0, 257, 0,
	29, 17, 0, 0, 18, 0, 23, 0, 0, 27,
	24, 25, 26, 0, 0, 0, 0, 21, 22, 258,
	259, 30, 32, 19, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 33,
	0, 28, 0, 29, 0, 0, 0, 0, 0, 23,
	0, 0, 27, 24, 25, 26, 0, 0, 0, 0,
	21, 22, 0, 0, 30, 32, 19,
}
var yyPact = [...]int{

	975, -1000, 975, -42, -42, -42, -42, -42, -42, -42,
	-42, -1000, -1000, -1000, -1000, -1000, 162, 480, 477, 562,
	-42, -42, -42, 573, 573, 573, 573, 573, 704, 704,
	-42, 556, 704, 504, 158, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 453, 529, 573, 565,
	547, 364, 295, -1000, 286, 573, 573, -42, 31, 160,
	-1000, -1000, -1000, 212, -1000, -42, -42, -42, 515, 808,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 158,
	-1000, 584, 17, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	704, 205, 137, 704, -1000, -1000, 243, -1000, -1000, -1000,
	-1000, 136, 819, 311, -11, -1000, 156, 14, -1000, 30,
	573, -1000, 704, 416, 444, 573, 543, 23, 538, 103,
	552, 520, 103, 328, 328, 328, 546, -1000, 104, 194,
	134, 472, -1000, 562, 704, 189, -1000, -1000, -1000, 559,
	975, 704, 704, 704, 306, 347, 359, 704, 704, 704,
	704, 704, 704, 704, -1000, 186, 89, 573, 295, 250,
	619, 121, 121, 310, 324, -1000, 517, -1000, -1000, 295,
	922, 573, 559, 626, -1000, 504, 190, 619, 414, 704,
	704, 119, 573, 573, -1000, 573, 520, 47, -1000, 518,
	-1000, -1000, -1000, -1000, 103, 449, 704, -1000, 194, -1000,
	194, 194, -1000, 22, 512, 619, -1000, -1000, -37, -1000,
	573, 201, 195, 573, -1000, 619, 286, -42, 19, 268,
	55, -20, -20, 367, 704, 121, 704, 121, 121, 114,
	114, -1000, -1000, -1000, 874, 517, -1000, 704, -1000, -1000,
	82, 694, 247, 704, -1000, 584, -1000, -1000, 121, 133,
	130, 125, 453, 238, 922, -1000, -1000, 704, -42, -42,
	241, -1000, -42, -1000, 123, 573, -1000, 704, 497, -1000,
	1, 427, 619, -1000, 121, 573, -1000, 547, 0, 110,
	-38, -1000, -1000, -1000, 399, 475, 357, 391, 103, -1000,
	-1000, -1000, -1000, -1000, 573, 520, 446, 439, 619, 334,
	-1000, -1000, 334, 546, 573, 295, 78, -10, 514, 573,
	470, -1000, 573, 467, -42, -1000, 237, 268, 975, 704,
	-1000, -1000, 865, -1000, -20, -1000, -1000, -1000, 45, -1000,
	-1000, -1000, 234, 250, 704, 797, 315, 85, -1000, 85,
	-1000, 85, -1000, 26, 258, -1000, 788, -1000, -1000, 922,
	-1000, 286, 76, 619, -1000, 287, 405, 704, 298, -1000,
	-1000, -1000, -16, 73, -21, 520, 573, 704, 103, 385,
	357, 383, -1000, 103, -1000, -1000, -1000, -1000, 704, 704,
	-1000, -1000, 72, -1000, 573, -1000, -1000, -1000, 573, 573,
	70, -22, 704, 68, 573, -1000, 283, 233, 265, -1000,
	742, 704, -1000, 619, 704, 121, 66, 62, 59, -1000,
	569, -42, 922, 219, 58, 508, -1000, -1000, -1000, 487,
	121, 343, 573, -1000, -1000, 619, 683, 103, 378, 103,
	617, 619, -1000, 171, -1000, -1000, -1000, 514, 573, 619,
	-1000, -1000, -42, 282, 975, 517, 619, -1000, -1000, -1000,
	-1000, 57, -1000, 260, 975, 253, -1000, 122, -1000, -1000,
	-1000, -1000, 121, -1000, -1000, -1000, 704, 120, 617, 103,
	683, 108, -1000, -1000, -1000, -42, -1000, -1000, 217, 260,
	922, 704, -42, 286, -1000, 619, 573, 617, -1000, 173,
	-1000, 278, 203, 262, -1000, 723, -1000, 34, 33, 27,
	453, 438, -42, 272, 922, -1000, -1000, -1000, -1000, 704,
	-1000, -42, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 27, 15, 12, 723, 720, 718, 717, 712, 706,
	696, 308, 78, 82, 694, 52, 26, 692, 690, 1,
	687, 50, 683, 94, 681, 61, 65, 142, 314, 29,
	60, 37, 680, 679, 678, 677, 413, 676, 673, 667,
	661, 656, 51, 652, 36, 645, 644, 56, 643, 45,
	640, 3, 638, 630, 629, 628, 616, 25, 4, 43,
	59, 2, 49, 69, 611, 610, 601, 30, 598, 597,
	596, 57, 13, 44, 592, 586, 54, 34, 24, 18,
	14, 585, 313, 282, 55, 583, 53, 16, 62, 32,
	581, 72, 401, 75, 48, 20, 47, 85, 580, 285,
	0,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10,
	11, 12, 12, 12, 12, 13, 13, 14, 15, 15,
	16, 16, 17, 17, 18, 18, 19, 19, 20, 20,
	20, 21, 21, 22, 22, 23, 23, 24, 24, 25,
	25, 26, 26, 26, 26, 26, 26, 27, 27, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 29, 29, 30, 30, 31, 31, 32, 32, 33,
	33, 34, 34, 34, 35, 35, 36, 37, 38, 38,
	38, 38, 38, 38, 38, 38, 38, 38, 38, 38,
	38, 38, 38, 38, 38, 38, 38, 39, 39, 39,
	39, 39, 40, 40, 40, 41, 41, 42, 42, 42,
	43, 43, 44, 45, 46, 46, 47, 47, 47, 48,
	48, 49, 49, 49, 49, 49, 49, 50, 50, 50,
	50, 50, 51, 51, 51, 52, 52, 52, 53, 53,
	54, 55, 55, 56, 56, 57, 57, 58, 58, 59,
	59, 60, 60, 61, 61, 62, 62, 63, 63, 64,
	64, 64, 64, 65, 66, 67, 67, 68, 68, 69,
	70, 70, 71, 71, 72, 72, 73, 73, 73, 73,
	73, 74, 74, 75, 76, 76, 77, 77, 78, 78,
	79, 79, 80, 81, 82, 82, 83, 83, 84, 85,
	86, 87, 88, 88, 89, 90, 90, 91, 91, 92,
	92, 93, 93, 94, 94, 95, 95, 96, 96, 96,
	96, 97, 97, 98, 98, 99, 99, 100, 100,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 5, 8,
	9, 7, 9, 2, 8, 9, 2, 2, 5, 3,
	5, 5, 4, 4, 4, 1, 1, 3, 0, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 3,
	4, 0, 2, 0, 2, 0, 2, 6, 9, 1,
	3, 1, 1, 1, 1, 1, 1, 1, 3, 1,
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

	-1000, -1, -3, -11, -64, -65, -68, -69, -70, -74,
	-75, -5, -6, -7, -8, -10, -23, 26, 29, 101,
	-89, 95, 96, 84, 88, 89, 90, 87, 76, 78,
	99, 16, 100, 74, -87, 11, -1, -100, 115, -100,
	-100, -100, -100, -100, -100, -100, -12, 19, 15, 17,
	-14, -13, 13, -36, 111, 30, 30, -91, -90, 11,
	-100, -100, -100, -80, 4, -80, -80, -80, -80, -28,
	-27, -26, -39, -37, -36, -41, -54, -38, -40, -87,
	-89, 111, -80, -81, -82, -83, -84, -85, -86, -43,
	75, -29, 63, 67, 5, 6, 107, 7, 9, 10,
	71, 91, -28, -88, -87, -100, 12, -28, -25, -24,
	-98, 25, 104, -19, 38, 20, -60, -47, -80, 14,
	-60, -15, 14, 58, 59, 60, -92, 73, -11, -23,
	-80, -80, -100, 113, 104, 85, -100, -100, -100, 20,
	80, 103, 102, 105, 70, -93, -99, 106, 107, 108,
	109, 110, 66, 65, 67, -28, -58, 114, 111, -55,
	-28, 102, 105, -93, -99, -36, -28, -82, -83, 111,
	82, 64, 113, 105, -100, 113, -80, -28, -20, 44,
	41, -80, 16, 113, -80, 22, -59, -49, -47, -48,
	-50, 23, -36, 24, 14, -16, 18, -59, -97, 61,
	-97, -97, -62, -53, -52, -28, -44, 108, -80, 112,
	111, 27, 28, 36, -91, -28, 86, -88, -87, -1,
	-28, -28, -28, -93, 68, 64, 69, 62, 61, -28,
	-28, -28, -28, -28, -28, -28, 112, 113, 112, -80,
	-42, -92, -63, 79, -29, 111, -36, -29, 68, 64,
	62, 61, -42, -2, -4, -3, -9, 76, 97, 98,
	-80, -88, -26, -25, 22, 111, -22, 45, -28, -31,
	-32, -33, -28, -44, 21, 111, -11, -67, -66, -27,
	-80, -60, -80, -16, -94, 57, -96, 54, 113, 49,
	51, 52, 53, -80, 22, -59, -17, 39, -28, -13,
	-12, -13, -13, 113, 22, 111, -61, -80, -71, 111,
	-80, -27, 111, -27, -11, -100, -77, -76, 81, 77,
	-84, -86, -28, -29, -28, -29, -29, -58, -28, 112,
	108, -58, -56, -63, 81, -28, -29, 111, -36, 111,
	-36, 111, -36, -19, 83, -2, -28, -100, -100, 82,
	-100, 111, -61, -28, -21, 47, 74, 113, -34, 42,
	43, -30, -29, -57, -27, -15, 113, 105, 48, -94,
	-96, -95, 50, 48, -59, -80, -16, -18, 40, 41,
	-62, -80, -42, 112, 113, -73, 31, 32, 33, 34,
	-72, -71, 35, -57, 37, -100, 83, -77, -76, -1,
	-28, 65, 83, -28, 80, 65, -30, -30, -30, 112,
	92, 78, 80, -2, -11, 112, -21, 46, -31, 72,
	113, 112, 113, -16, -67, -28, -49, 48, -95, 48,
	-49, -28, -58, 112, -61, -27, -27, 112, 113, -28,
	112, -80, 76, 83, 80, -28, -28, -29, 112, 112,
	112, 5, -100, -2, -3, 83, 112, 22, -35, 31,
	32, -30, 21, -11, -57, -51, 55, 56, -49, 48,
	-49, 94, -73, -72, -100, 76, -1, 112, -79, -78,
	81, 77, 78, 111, -30, -28, 111, -49, -51, 111,
	-100, 83, -79, -78, -2, -28, -100, -11, -61, -45,
	-46, 93, 76, 83, 80, 112, 112, 112, -19, 41,
	-100, 76, -2, -58, -100,
}
var yyDef = [...]int{

	-2, -2, -2, 257, 257, 257, 257, 257, 257, 257,
	257, 13, 14, 15, 16, 17, 0, 0, 0, 0,
	257, 257, 257, 0, 0, 0, 0, 0, 0, 0,
	257, 0, 0, 253, 0, 231, 2, 5, 258, 6,
	7, 8, 9, 10, 11, 12, -2, 0, 0, 0,
	48, 0, 239, 46, 65, 0, 0, 257, 237, 235,
	21, 22, 23, 0, 222, 257, 257, 257, 0, 241,
	79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
	89, 65, 77, 71, 72, 73, 74, 75, 76, 136,
	171, 241, 0, 0, 223, 224, 0, 226, 228, 229,
	230, 0, 241, 0, 88, 33, 0, -2, 66, 69,
	0, 254, 0, 58, 0, 0, 0, 181, 146, 0,
	0, 50, 0, 251, 251, 251, 0, 240, 0, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	-2, 0, 255, 256, 241, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 242, 241, 0, 0, -2, 0,
	-2, 255, 256, 0, 0, 126, 134, 225, 227, -2,
	-2, 0, 0, 0, 39, 253, 0, 234, 63, 0,
	0, 65, 0, 0, 147, 0, 50, -2, 151, 152,
	155, 156, 149, 150, 0, 52, 0, 49, 0, 252,
	0, 0, 47, 185, 168, -2, 166, 167, 77, 106,
	0, 0, 0, 0, 238, -2, 65, 257, 232, 216,
	107, -2, -2, 0, 0, 0, 0, 0, 0, 127,
	128, 129, 130, 131, 132, 133, 90, 0, 91, 78,
	0, 0, 173, 0, 109, 65, 92, 111, 0, 0,
	0, 0, 56, 0, -2, 18, 19, 0, 257, 257,
	0, 233, 257, 70, 0, 0, 40, 0, -2, 57,
	95, 101, -2, 100, 0, 0, 191, 48, 195, 0,
	77, 182, 148, 197, 0, -2, 245, 0, 0, 244,
	248, 249, 250, 153, 0, 50, 54, 0, -2, 42,
	45, 43, 44, 0, 0, -2, 0, 183, 206, 0,
	202, 211, 0, 0, 257, 28, 0, 216, -2, 0,
	112, 113, 241, 116, -2, 120, 123, 178, -2, 135,
	138, 139, 0, 188, 0, 241, 0, 65, 118, 65,
	122, 65, 125, 0, 0, 4, 241, 36, 37, -2,
	38, 65, 0, -2, 59, 61, 0, 0, 97, 102,
	103, 189, 93, 0, 175, 50, 0, 0, 0, 0,
	245, 0, 246, 0, 180, 154, 198, 41, 0, 0,
	186, 169, 0, 199, 0, 200, 207, 208, 0, 0,
	0, 204, 0, 0, 0, 24, 0, 0, 215, 217,
	241, 0, 170, -2, 0, 0, 0, 0, 0, 140,
	0, 257, -2, 0, 0, 0, 60, 62, 96, 0,
	0, 65, 0, 193, 196, -2, 162, 0, 0, 0,
	161, -2, 53, 135, 184, 209, 210, 206, 0, -2,
	212, 213, 257, 0, -2, 114, -2, 115, 117, 121,
	124, 0, 31, 220, -2, 0, 67, 0, 98, 104,
	105, 94, 0, 192, 176, 157, 0, 0, 158, 0,
	162, 0, 201, 205, 29, 257, 214, 141, 0, 220,
	-2, 0, 257, 65, 190, -2, 0, 160, 159, 144,
	30, 0, 0, 219, 221, 241, 32, 0, 0, 0,
	56, 0, 257, 0, -2, 68, 164, 142, 143, 0,
	34, 257, 218, 145, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 110, 3, 3,
	111, 112, 108, 106, 113, 107, 114, 109, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 115,
	3, 105,
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
	102, 103, 104,
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
		//line parser.y:159
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:164
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:171
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:176
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:183
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:187
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:191
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:195
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:199
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:203
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:207
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:211
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:215
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:219
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:223
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:227
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:231
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:237
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:241
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:247
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:251
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:257
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:261
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:267
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:271
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:275
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:279
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:283
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:289
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:293
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:297
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:301
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:305
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:311
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:315
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:319
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:323
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:329
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:333
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:339
		{
			yyVAL.expression = SelectQuery{
				CommonTableClause: yyDollar[1].expression,
				SelectEntity:      yyDollar[2].expression,
				OrderByClause:     yyDollar[3].expression,
				LimitClause:       yyDollar[4].expression,
				OffsetClause:      yyDollar[5].expression,
			}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:351
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
		//line parser.y:361
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
		//line parser.y:370
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
		//line parser.y:379
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
		//line parser.y:390
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:394
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:400
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:406
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:410
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:426
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:430
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:436
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:440
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:446
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:450
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 58:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:456
		{
			yyVAL.expression = nil
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:460
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, With: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:464
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Value: yyDollar[2].expression, Percent: yyDollar[3].token.Literal, With: yyDollar[4].expression}
		}
	case 61:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:470
		{
			yyVAL.expression = nil
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:474
		{
			yyVAL.expression = LimitWith{With: yyDollar[1].token.Literal, Type: yyDollar[2].token}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:480
		{
			yyVAL.expression = nil
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:484
		{
			yyVAL.expression = OffsetClause{Offset: yyDollar[1].token.Literal, Value: yyDollar[2].expression}
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:490
		{
			yyVAL.expression = nil
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:494
		{
			yyVAL.expression = CommonTableClause{With: yyDollar[1].token.Literal, CommonTables: yyDollar[2].expressions}
		}
	case 67:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:500
		{
			yyVAL.expression = CommonTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, As: yyDollar[3].token.Literal, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 68:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:504
		{
			yyVAL.expression = CommonTable{Recursive: yyDollar[1].token, Name: yyDollar[2].identifier, Columns: yyDollar[4].expressions, As: yyDollar[6].token.Literal, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:510
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:514
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:520
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:524
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:528
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:532
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:540
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:546
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:550
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:560
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:564
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:568
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:580
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:584
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:588
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:600
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:606
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:610
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:616
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:620
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:626
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:630
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:636
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 98:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:640
		{
			yyVAL.expression = OrderItem{Value: yyDollar[1].expression, Direction: yyDollar[2].token, Nulls: yyDollar[3].token.Literal, Position: yyDollar[4].token}
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
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:656
		{
			yyVAL.token = Token{}
		}
	case 102:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:660
		{
			yyVAL.token = yyDollar[1].token
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:664
		{
			yyVAL.token = yyDollar[1].token
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:670
		{
			yyVAL.token = yyDollar[1].token
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:674
		{
			yyVAL.token = yyDollar[1].token
		}
	case 106:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:680
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 107:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:686
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
	case 108:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:709
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:713
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 110:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:717
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:721
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 112:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:725
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 113:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:729
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 114:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:733
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 115:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:737
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 116:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:741
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 117:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:745
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 118:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:749
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 119:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:753
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 120:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:757
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 121:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:761
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 122:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:765
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 123:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:769
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 124:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:773
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 125:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:777
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 126:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:781
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:787
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:791
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 129:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:795
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 130:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:799
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 131:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:803
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 132:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:809
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 133:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:813
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 134:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:817
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 135:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:823
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:827
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 137:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:833
		{
			yyVAL.expression = Option{}
		}
	case 138:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:837
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 139:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:841
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 140:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:847
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 141:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:851
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 142:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expression = AnalyticFunction{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option), Over: yyDollar[5].token.Literal, AnalyticClause: yyDollar[7].expression.(AnalyticClause)}
		}
	case 143:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:863
		{
			yyVAL.expression = AnalyticClause{Partition: yyDollar[1].expression, OrderByClause: yyDollar[2].expression}
		}
	case 144:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:869
		{
			yyVAL.expression = nil
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:873
		{
			yyVAL.expression = Partition{PartitionBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Values: yyDollar[3].expressions}
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:879
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 147:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:883
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 148:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:887
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 149:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:893
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:897
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 151:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:903
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:907
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:911
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 154:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:915
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 155:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:919
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:923
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 157:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:929
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 158:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:933
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 159:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:937
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 160:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:941
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 161:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:945
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 162:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:951
		{
			yyVAL.expression = nil
		}
	case 163:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:955
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 164:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:959
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:965
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:969
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 167:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:973
		{
			yyVAL.expression = AllColumns{}
		}
	case 168:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:979
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 169:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:983
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 170:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:989
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 171:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:995
		{
			yyVAL.expression = nil
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:999
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 173:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1005
		{
			yyVAL.expression = nil
		}
	case 174:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1009
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1015
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 176:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1019
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 177:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1025
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 178:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1029
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 179:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1035
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 180:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1039
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 181:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1045
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 182:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1049
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 183:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1055
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 184:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1059
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1065
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 186:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1069
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 187:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1075
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 188:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1079
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 189:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1085
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Values: yyDollar[5].token.Literal, ValuesList: yyDollar[6].expressions}
		}
	case 190:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:1089
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Values: yyDollar[8].token.Literal, ValuesList: yyDollar[9].expressions}
		}
	case 191:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1093
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 192:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1097
		{
			yyVAL.expression = InsertQuery{CommonTableClause: yyDollar[1].expression, Insert: yyDollar[2].token.Literal, Into: yyDollar[3].token.Literal, Table: yyDollar[4].identifier, Fields: yyDollar[6].expressions, Query: yyDollar[8].expression.(SelectQuery)}
		}
	case 193:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1103
		{
			yyVAL.expression = UpdateQuery{CommonTableClause: yyDollar[1].expression, Update: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, Set: yyDollar[4].token.Literal, SetList: yyDollar[5].expressions, FromClause: yyDollar[6].expression, WhereClause: yyDollar[7].expression}
		}
	case 194:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1109
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1115
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 196:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1119
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 197:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1125
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{CommonTableClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 198:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1130
		{
			from := FromClause{From: yyDollar[4].token.Literal, Tables: yyDollar[5].expressions}
			yyVAL.expression = DeleteQuery{CommonTableClause: yyDollar[1].expression, Delete: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions, FromClause: from, WhereClause: yyDollar[6].expression}
		}
	case 199:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1137
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 200:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1143
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 201:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1147
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1153
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 203:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1157
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 204:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1163
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 205:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1167
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 206:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1173
		{
			yyVAL.expression = nil
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1177
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 208:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1181
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 209:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1185
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 210:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1189
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 211:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1195
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 212:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1199
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 213:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1205
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 214:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1211
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 215:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1215
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 216:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1221
		{
			yyVAL.procexpr = nil
		}
	case 217:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1225
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 218:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1231
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 219:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1235
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 220:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1241
		{
			yyVAL.procexpr = nil
		}
	case 221:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1245
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1251
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 223:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1257
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1263
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 225:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1267
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 226:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1274
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 227:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1278
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 228:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1285
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 229:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1291
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 230:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1297
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 231:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1303
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 232:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1309
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 233:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1313
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 234:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1319
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 235:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1325
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 236:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1329
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 237:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1335
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 238:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1339
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 239:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1345
		{
			yyVAL.token = Token{}
		}
	case 240:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1349
		{
			yyVAL.token = yyDollar[1].token
		}
	case 241:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1355
		{
			yyVAL.token = Token{}
		}
	case 242:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1359
		{
			yyVAL.token = yyDollar[1].token
		}
	case 243:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1365
		{
			yyVAL.token = Token{}
		}
	case 244:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1369
		{
			yyVAL.token = yyDollar[1].token
		}
	case 245:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1375
		{
			yyVAL.token = Token{}
		}
	case 246:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1379
		{
			yyVAL.token = yyDollar[1].token
		}
	case 247:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1385
		{
			yyVAL.token = Token{}
		}
	case 248:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1389
		{
			yyVAL.token = yyDollar[1].token
		}
	case 249:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1393
		{
			yyVAL.token = yyDollar[1].token
		}
	case 250:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1397
		{
			yyVAL.token = yyDollar[1].token
		}
	case 251:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1403
		{
			yyVAL.token = Token{}
		}
	case 252:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1407
		{
			yyVAL.token = yyDollar[1].token
		}
	case 253:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1413
		{
			yyVAL.token = Token{}
		}
	case 254:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1417
		{
			yyVAL.token = yyDollar[1].token
		}
	case 255:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1424
		{
			yyVAL.token = yyDollar[1].token
		}
	case 256:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1428
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 257:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1434
		{
			yyVAL.token = Token{}
		}
	case 258:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1438
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
