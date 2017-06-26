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

//line parser.y:1325

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
	61, 222,
	65, 222,
	66, 222,
	-2, 236,
	-1, 132,
	45, 224,
	47, 228,
	-2, 160,
	-1, 167,
	55, 46,
	56, 46,
	57, 46,
	-2, 74,
	-1, 169,
	107, 121,
	-2, 220,
	-1, 171,
	75, 151,
	-2, 222,
	-1, 180,
	37, 121,
	88, 121,
	107, 121,
	-2, 220,
	-1, 197,
	61, 222,
	65, 222,
	66, 222,
	-2, 145,
	-1, 205,
	61, 222,
	65, 222,
	66, 222,
	-2, 87,
	-1, 217,
	47, 228,
	-2, 224,
	-1, 233,
	61, 222,
	65, 222,
	66, 222,
	-2, 217,
	-1, 239,
	67, 0,
	96, 0,
	99, 0,
	-2, 92,
	-1, 240,
	67, 0,
	96, 0,
	99, 0,
	-2, 94,
	-1, 283,
	61, 222,
	65, 222,
	66, 222,
	-2, 51,
	-1, 329,
	67, 0,
	96, 0,
	99, 0,
	-2, 103,
	-1, 333,
	61, 222,
	65, 222,
	66, 222,
	-2, 156,
	-1, 367,
	61, 222,
	65, 222,
	66, 222,
	-2, 175,
	-1, 393,
	79, 153,
	-2, 222,
	-1, 397,
	107, 83,
	108, 83,
	-2, 46,
	-1, 405,
	61, 222,
	65, 222,
	66, 222,
	-2, 55,
	-1, 422,
	61, 222,
	65, 222,
	66, 222,
	-2, 184,
	-1, 429,
	75, 168,
	77, 168,
	79, 168,
	-2, 222,
	-1, 437,
	91, 18,
	92, 18,
	-2, 1,
	-1, 440,
	61, 222,
	65, 222,
	66, 222,
	-2, 143,
}

const yyPrivate = 57344

const yyLast = 1132

var yyAct = [...]int{

	80, 40, 450, 40, 311, 459, 375, 44, 76, 380,
	305, 321, 46, 47, 48, 49, 50, 51, 52, 297,
	411, 2, 186, 94, 209, 132, 166, 203, 271, 194,
	92, 67, 68, 69, 53, 43, 1, 388, 218, 216,
	131, 111, 114, 40, 381, 258, 109, 86, 23, 77,
	23, 338, 189, 64, 155, 133, 56, 45, 59, 119,
	221, 137, 222, 223, 224, 219, 421, 374, 217, 183,
	364, 362, 142, 183, 57, 57, 61, 136, 138, 146,
	147, 148, 401, 300, 291, 288, 151, 59, 167, 157,
	158, 159, 160, 161, 39, 39, 39, 37, 143, 176,
	128, 400, 408, 206, 463, 16, 449, 433, 163, 162,
	164, 432, 431, 154, 423, 420, 137, 39, 373, 363,
	334, 185, 220, 59, 39, 59, 256, 40, 85, 38,
	199, 38, 168, 169, 441, 295, 211, 263, 346, 344,
	137, 342, 152, 151, 228, 153, 157, 158, 159, 160,
	161, 40, 174, 180, 215, 45, 188, 159, 160, 161,
	42, 110, 157, 158, 159, 160, 161, 42, 168, 164,
	191, 192, 301, 264, 264, 184, 101, 103, 91, 156,
	227, 234, 40, 119, 57, 213, 237, 42, 207, 42,
	40, 144, 40, 40, 145, 465, 235, 232, 23, 90,
	354, 172, 457, 273, 173, 438, 426, 264, 40, 241,
	263, 392, 386, 349, 453, 324, 298, 42, 452, 323,
	261, 137, 260, 454, 402, 317, 270, 314, 453, 23,
	279, 261, 324, 339, 280, 40, 468, 464, 447, 316,
	318, 425, 121, 320, 264, 104, 264, 264, 164, 243,
	310, 267, 299, 242, 244, 266, 304, 303, 395, 182,
	333, 308, 97, 190, 167, 326, 117, 264, 343, 345,
	347, 102, 325, 40, 313, 322, 306, 175, 236, 38,
	415, 179, 332, 371, 352, 353, 336, 284, 355, 286,
	287, 369, 75, 108, 273, 285, 113, 285, 285, 269,
	268, 350, 178, 137, 106, 348, 246, 245, 137, 211,
	38, 307, 236, 302, 201, 370, 124, 358, 359, 361,
	23, 125, 365, 357, 40, 366, 298, 385, 368, 116,
	117, 118, 282, 372, 387, 259, 221, 383, 222, 223,
	224, 219, 54, 397, 217, 397, 384, 397, 165, 221,
	382, 222, 223, 224, 63, 40, 62, 171, 289, 389,
	177, 149, 55, 264, 40, 229, 230, 187, 127, 115,
	137, 23, 137, 298, 231, 120, 273, 139, 112, 417,
	193, 197, 59, 404, 410, 406, 205, 418, 419, 376,
	377, 378, 379, 59, 41, 414, 264, 416, 59, 66,
	226, 38, 23, 290, 40, 233, 202, 434, 60, 264,
	435, 130, 238, 239, 240, 59, 137, 292, 247, 248,
	249, 250, 251, 252, 253, 437, 65, 444, 40, 93,
	88, 445, 436, 446, 89, 262, 265, 443, 40, 237,
	10, 442, 451, 9, 8, 7, 455, 6, 283, 58,
	58, 23, 38, 40, 458, 456, 210, 70, 71, 72,
	73, 74, 462, 448, 5, 4, 337, 40, 170, 296,
	82, 195, 467, 196, 273, 23, 470, 135, 396, 134,
	398, 460, 399, 38, 95, 23, 126, 3, 273, 129,
	81, 58, 84, 140, 141, 469, 78, 83, 407, 79,
	23, 204, 200, 327, 123, 329, 328, 356, 330, 331,
	281, 36, 15, 274, 23, 100, 101, 103, 14, 104,
	105, 13, 340, 12, 11, 272, 0, 0, 0, 341,
	122, 0, 38, 0, 0, 221, 351, 222, 223, 224,
	219, 412, 413, 217, 439, 0, 58, 0, 0, 197,
	0, 0, 205, 0, 0, 0, 38, 0, 212, 58,
	0, 214, 367, 0, 0, 225, 38, 0, 0, 0,
	58, 0, 0, 0, 0, 122, 0, 0, 106, 0,
	0, 38, 163, 162, 164, 390, 0, 154, 0, 0,
	0, 0, 0, 0, 0, 38, 466, 0, 0, 257,
	393, 0, 0, 0, 0, 296, 0, 296, 0, 296,
	0, 102, 0, 278, 208, 0, 152, 151, 405, 153,
	157, 158, 159, 160, 161, 296, 0, 0, 41, 0,
	39, 0, 18, 34, 19, 0, 17, 0, 212, 0,
	0, 0, 20, 422, 0, 21, 0, 0, 0, 0,
	0, 58, 428, 0, 0, 429, 0, 309, 430, 312,
	315, 212, 212, 0, 0, 0, 0, 0, 0, 0,
	0, 296, 0, 440, 0, 0, 0, 0, 0, 59,
	100, 101, 103, 0, 104, 105, 41, 0, 0, 275,
	0, 32, 0, 0, 0, 122, 0, 26, 0, 0,
	30, 27, 28, 29, 0, 0, 24, 25, 276, 277,
	33, 35, 22, 0, 461, 0, 0, 0, 0, 0,
	360, 0, 319, 42, 0, 0, 0, 0, 0, 0,
	0, 212, 0, 58, 0, 98, 0, 0, 58, 99,
	0, 0, 0, 106, 0, 315, 96, 0, 212, 0,
	0, 122, 163, 162, 164, 0, 0, 154, 0, 0,
	0, 0, 107, 0, 0, 0, 41, 0, 39, 0,
	18, 34, 19, 0, 17, 0, 102, 198, 0, 0,
	20, 87, 0, 21, 0, 0, 152, 151, 0, 153,
	157, 158, 159, 160, 161, 212, 0, 254, 255, 0,
	58, 0, 58, 0, 0, 312, 0, 0, 0, 212,
	212, 0, 0, 0, 0, 424, 0, 59, 100, 101,
	103, 0, 104, 105, 41, 0, 39, 31, 0, 32,
	122, 0, 122, 0, 122, 26, 0, 0, 30, 27,
	28, 29, 0, 0, 24, 25, 58, 0, 33, 35,
	22, 409, 315, 163, 162, 164, 0, 0, 154, 0,
	0, 42, 59, 100, 101, 103, 0, 104, 105, 41,
	0, 0, 312, 98, 0, 0, 0, 99, 0, 0,
	0, 106, 0, 0, 96, 0, 0, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 0, 0, 0, 255,
	107, 59, 100, 101, 103, 0, 104, 105, 41, 0,
	0, 293, 294, 0, 102, 0, 0, 0, 98, 87,
	0, 0, 99, 0, 0, 0, 106, 0, 0, 96,
	0, 0, 163, 162, 164, 0, 0, 154, 0, 0,
	0, 163, 162, 164, 0, 107, 154, 0, 0, 0,
	0, 0, 0, 0, 0, 427, 0, 98, 0, 102,
	335, 99, 0, 0, 87, 106, 152, 151, 96, 153,
	157, 158, 159, 160, 161, 152, 151, 0, 153, 157,
	158, 159, 160, 161, 107, 163, 162, 164, 0, 0,
	154, 0, 0, 0, 0, 163, 162, 164, 102, 403,
	154, 0, 0, 87, 0, 163, 162, 164, 0, 394,
	154, 0, 0, 0, 0, 0, 0, 0, 0, 152,
	151, 181, 153, 157, 158, 159, 160, 161, 0, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 0, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 163, 162,
	164, 0, 0, 154, 0, 0, 0, 0, 163, 162,
	164, 0, 150, 154, 0, 0, 0, 391, 162, 164,
	0, 0, 154, 0, 0, 0, 163, 0, 164, 0,
	0, 154, 152, 151, 0, 153, 157, 158, 159, 160,
	161, 164, 152, 151, 154, 153, 157, 158, 159, 160,
	161, 152, 151, 0, 153, 157, 158, 159, 160, 161,
	152, 151, 0, 153, 157, 158, 159, 160, 161, 0,
	0, 0, 0, 152, 151, 0, 153, 157, 158, 159,
	160, 161,
}
var yyPact = [...]int{

	755, -1000, 755, -52, -52, -52, -52, -52, -52, -52,
	-52, -1000, -1000, -1000, -1000, -1000, 305, 342, 411, 394,
	327, 325, 388, -52, -52, -52, 411, 411, 411, 411,
	411, 897, 897, -52, 366, 897, 355, 274, 85, 173,
	-1000, -1000, 111, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 273, 281, 411, 352, -8, 389, -1000,
	54, 363, 411, 411, -52, -10, 93, -1000, -1000, -1000,
	113, -52, -52, -52, 341, 986, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 85, -1000, 813, 27, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 897, 105, 61, 897,
	-1000, -1000, 170, -1000, -1000, -1000, -1000, 47, 943, 198,
	-39, -1000, 76, 46, 349, 54, 205, 205, 205, 897,
	675, -1000, 23, 270, 400, 897, 82, 411, 411, -1000,
	411, 349, 14, -1000, 378, -1000, -1000, -1000, -1000, 54,
	38, 339, -1000, 388, 897, 99, -1000, -1000, -1000, 383,
	755, 897, 897, 897, 184, 188, 248, 897, 897, 897,
	897, 897, 897, 897, -1000, 690, 19, -1000, 411, 173,
	145, 996, 31, 31, 190, 241, -1000, 1027, -1000, -1000,
	173, 617, 411, 383, 510, -1000, 294, 897, -1000, 111,
	-1000, 111, 111, 996, -1000, -23, 336, 996, -1000, -1000,
	-1000, 397, -1000, -1000, -24, 870, 31, 83, -1000, 355,
	-25, 73, 63, -1000, -1000, -1000, 268, 303, 229, 266,
	54, -1000, -1000, -1000, -1000, -1000, 411, 349, 411, 121,
	119, 411, -1000, 996, 111, -52, -35, 142, 62, -11,
	-11, 236, 897, 31, 897, 31, 31, 55, 55, -1000,
	-1000, -1000, 1014, 1027, -1000, 897, -1000, -1000, 13, 858,
	156, 897, -1000, 813, -1000, -1000, 31, 35, 33, 32,
	305, 134, 617, -1000, -1000, 897, -52, -52, 122, -1000,
	-52, 284, 277, 996, 210, -1000, -1000, 210, 675, 411,
	-1000, 897, -1000, -1000, -1000, -1000, -37, 12, -38, 349,
	411, 897, 54, 246, 229, 238, -1000, 54, -1000, -1000,
	-1000, 11, -41, 359, 411, 316, -1000, 411, 310, -52,
	-1000, 133, 142, 755, 897, -1000, -1000, 1005, -1000, -11,
	-1000, -1000, -1000, 791, -1000, -1000, -1000, 132, 145, 897,
	933, 196, 104, -1000, 104, -1000, 104, -1000, -6, 150,
	-1000, 923, -1000, -1000, 617, -1000, -1000, 897, 897, -1000,
	-1000, -1000, 31, 81, 411, -1000, -1000, 996, 489, 54,
	235, 54, 290, -1000, 411, -1000, -1000, -1000, 411, 411,
	8, -42, 897, 7, 411, -1000, 169, 127, 159, -1000,
	879, 897, -1000, 996, 897, 31, 5, -1000, 4, 0,
	-1000, 402, -52, 617, 126, 996, -1000, -1000, 31, -1000,
	-1000, -1000, 897, 28, 290, 54, 489, -1000, -1000, -1000,
	359, 411, 996, -1000, -1000, -52, 166, 755, 1027, 996,
	-1000, -1000, -1000, -1000, -1, -1000, 141, 755, 149, -1000,
	996, 411, 290, -1000, -1000, -1000, -1000, -52, -1000, -1000,
	123, 141, 617, 897, -52, -3, -1000, 165, 116, 155,
	-1000, 520, -1000, -1000, -52, 164, 617, -1000, -52, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 35, 28, 21, 525, 524, 523, 521, 518, 513,
	512, 487, 105, 97, 511, 42, 22, 510, 507, 34,
	504, 502, 49, 8, 260, 262, 135, 501, 0, 499,
	497, 496, 492, 490, 45, 484, 55, 479, 25, 477,
	20, 473, 471, 470, 468, 466, 19, 26, 27, 40,
	56, 4, 29, 51, 465, 464, 456, 24, 447, 445,
	444, 44, 9, 6, 443, 440, 37, 11, 5, 2,
	430, 434, 199, 178, 30, 429, 23, 128, 46, 47,
	426, 53, 335, 54, 417, 39, 10, 38, 52, 179,
	7,
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
	24, 24, 25, 25, 26, 26, 27, 84, 84, 84,
	28, 29, 30, 30, 30, 30, 30, 30, 30, 30,
	30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
	30, 31, 31, 31, 31, 31, 32, 32, 32, 33,
	33, 34, 34, 34, 35, 35, 36, 36, 36, 37,
	37, 38, 38, 38, 38, 38, 38, 39, 39, 39,
	39, 39, 40, 40, 40, 41, 41, 42, 42, 43,
	44, 44, 45, 45, 46, 46, 47, 47, 48, 48,
	49, 49, 50, 50, 51, 51, 52, 52, 53, 53,
	54, 54, 54, 54, 55, 56, 57, 57, 58, 58,
	59, 60, 60, 61, 61, 62, 62, 63, 63, 63,
	63, 63, 64, 64, 65, 66, 66, 67, 67, 68,
	68, 69, 69, 70, 71, 72, 72, 73, 73, 74,
	75, 76, 77, 78, 78, 79, 80, 80, 81, 81,
	82, 82, 83, 83, 85, 85, 86, 86, 87, 87,
	87, 87, 88, 88, 89, 89, 90, 90,
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
	1, 3, 3, 1, 1, 3, 2, 0, 1, 1,
	3, 3, 3, 3, 3, 3, 4, 4, 6, 6,
	4, 6, 4, 4, 4, 6, 4, 4, 6, 4,
	2, 3, 3, 3, 3, 3, 3, 3, 2, 4,
	1, 0, 2, 2, 5, 7, 1, 2, 3, 1,
	1, 1, 1, 2, 3, 1, 1, 5, 5, 6,
	6, 4, 0, 2, 4, 1, 1, 1, 3, 5,
	0, 1, 0, 2, 1, 3, 1, 3, 1, 3,
	1, 3, 1, 3, 1, 3, 1, 3, 4, 2,
	5, 8, 4, 7, 6, 3, 1, 3, 4, 5,
	6, 6, 8, 1, 3, 1, 3, 0, 1, 1,
	2, 2, 5, 7, 7, 4, 2, 0, 2, 4,
	2, 0, 2, 1, 1, 1, 2, 1, 2, 1,
	1, 1, 1, 1, 3, 3, 1, 3, 1, 3,
	0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
	1, 1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -54, -55, -58, -59, -60, -64,
	-65, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 95, -79, 89, 90, 80, 84, 85, 86,
	83, 72, 74, 93, 16, 94, -14, -13, -77, 13,
	-28, 11, 106, -1, -90, 109, -90, -90, -90, -90,
	-90, -90, -90, -19, 37, 20, -50, -36, -70, 4,
	14, -50, 29, 29, -81, -80, 11, -90, -90, -90,
	-70, -70, -70, -70, -70, -24, -23, -22, -31, -29,
	-28, -33, -43, -30, -32, -77, -79, 106, -70, -71,
	-72, -73, -74, -75, -76, -35, 71, -25, 60, 64,
	5, 6, 101, 7, 9, 10, 68, 87, -24, -78,
	-77, -90, 12, -24, -15, 14, 55, 56, 57, 98,
	-82, 69, -11, -20, 43, 40, -70, 16, 108, -70,
	22, -49, -38, -36, -37, -39, 23, -28, 24, 14,
	-70, -70, -90, 108, 98, 81, -90, -90, -90, 20,
	76, 97, 96, 99, 67, -83, -89, 100, 101, 102,
	103, 104, 63, 62, 64, -24, -47, -28, 105, 106,
	-44, -24, 96, 99, -83, -89, -28, -24, -72, -73,
	106, 78, 61, 108, 99, -90, -16, 18, -49, -88,
	58, -88, -88, -24, -52, -42, -41, -24, 102, 107,
	-21, 44, 6, -48, -27, -24, 21, 106, -11, -57,
	-56, -23, -70, -50, -70, -16, -85, 54, -87, 51,
	108, 46, 48, 49, 50, -70, 22, -49, 106, 26,
	27, 35, -81, -24, 82, -78, -77, -1, -24, -24,
	-24, -83, 65, 61, 66, 59, 58, -24, -24, -24,
	-24, -24, -24, -24, 107, 108, 107, -70, -34, -82,
	-53, 75, -25, 106, -28, -25, 65, 61, 59, 58,
	-34, -2, -4, -3, -9, 72, 91, 92, -70, -78,
	-22, -17, 38, -24, -13, -12, -13, -13, 108, 22,
	6, 108, -84, 41, 42, -26, -25, -46, -23, -15,
	108, 99, 45, -85, -87, -86, 47, 45, -49, -70,
	-16, -51, -70, -61, 106, -70, -23, 106, -23, -11,
	-90, -67, -66, 77, 73, -74, -76, -24, -25, -24,
	-25, -25, -47, -24, 107, 102, -47, -45, -53, 77,
	-24, -25, 106, -28, 106, -28, 106, -28, -19, 79,
	-2, -24, -90, -90, 78, -90, -18, 39, 40, -52,
	-70, -48, 108, 107, 108, -16, -57, -24, -38, 45,
	-86, 45, -38, 107, 108, -63, 30, 31, 32, 33,
	-62, -61, 34, -46, 36, -90, 79, -67, -66, -1,
	-24, 62, 79, -24, 76, 62, -26, -28, -26, -26,
	107, 88, 74, 76, -2, -24, -47, -26, 21, -11,
	-46, -40, 52, 53, -38, 45, -38, -51, -23, -23,
	107, 108, -24, 107, -70, 72, 79, 76, -24, -24,
	-25, 107, 107, 107, 5, -90, -2, -3, 79, -26,
	-24, 106, -38, -40, -63, -62, -90, 72, -1, 107,
	-69, -68, 77, 73, 74, -51, -90, 79, -69, -68,
	-2, -24, -90, 107, 72, 79, 76, -90, 72, -2,
	-90,
}
var yyDef = [...]int{

	1, -2, 1, 236, 236, 236, 236, 236, 236, 236,
	236, 13, 14, 15, 16, 17, -2, 0, 0, 0,
	0, 0, 0, 236, 236, 236, 0, 0, 0, 0,
	0, 0, 0, 236, 0, 0, 48, 0, 0, 220,
	46, 212, 0, 2, 5, 237, 6, 7, 8, 9,
	10, 11, 12, 58, 0, 0, 0, 162, 126, 203,
	0, 0, 0, 0, 236, 218, 216, 21, 22, 23,
	0, 236, 236, 236, 0, 222, 70, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 80, 0, 68, 62,
	63, 64, 65, 66, 67, 120, 150, 222, 0, 0,
	204, 205, 0, 207, 209, 210, 211, 0, 222, 0,
	79, 33, 0, -2, 50, 0, 232, 232, 232, 0,
	0, 221, 0, 60, 0, 0, 0, 0, 0, 127,
	0, 50, -2, 131, 132, 135, 136, 129, 130, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	1, 0, 234, 235, 222, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 223, 222, 0, -2, 0, -2,
	0, -2, 234, 235, 0, 0, 110, 118, 206, 208,
	-2, 3, 0, 0, 0, 39, 52, 0, 49, 0,
	233, 0, 0, 215, 47, 166, 147, -2, 146, 90,
	40, 0, 59, 57, 158, -2, 0, 0, 172, 48,
	176, 0, 68, 163, 128, 178, 0, -2, 226, 0,
	0, 225, 229, 230, 231, 133, 0, 50, 0, 0,
	0, 0, 219, -2, 0, 236, 213, 197, 91, -2,
	-2, 0, 0, 0, 0, 0, 0, 111, 112, 113,
	114, 115, 116, 117, 81, 0, 82, 69, 0, 0,
	152, 0, 93, 0, 83, 95, 0, 0, 0, 0,
	56, 0, 3, 18, 19, 0, 236, 236, 0, 214,
	236, 54, 0, -2, 42, 45, 43, 44, 0, 0,
	61, 0, 86, 88, 89, 170, 84, 0, 154, 50,
	0, 0, 0, 0, 226, 0, 227, 0, 161, 134,
	179, 0, 164, 187, 0, 183, 192, 0, 0, 236,
	28, 0, 197, 1, 0, 96, 97, 222, 100, -2,
	104, 107, 157, -2, 119, 122, 123, 0, 169, 0,
	222, 0, 0, 102, 0, 106, 0, 109, 0, 0,
	4, 222, 36, 37, 3, 38, 41, 0, 0, 167,
	148, 159, 0, 0, 0, 174, 177, -2, 142, 0,
	0, 0, 141, 180, 0, 181, 188, 189, 0, 0,
	0, 185, 0, 0, 0, 24, 0, 0, 196, 198,
	222, 0, 149, -2, 0, 0, 0, -2, 0, 0,
	124, 0, 236, 1, 0, -2, 53, 85, 0, 173,
	155, 137, 0, 0, 138, 0, 142, 165, 190, 191,
	187, 0, -2, 193, 194, 236, 0, 1, 98, -2,
	99, 101, 105, 108, 0, 31, 201, -2, 0, 171,
	-2, 0, 140, 139, 182, 186, 29, 236, 195, 125,
	0, 201, 3, 0, 236, 0, 30, 0, 0, 200,
	202, 222, 32, 144, 236, 0, 3, 34, 236, 199,
	35,
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
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:550
		{
			yyVAL.expression = RowValue{Value: ValueList{Values: yyDollar[2].expressions}}
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:554
		{
			yyVAL.expression = RowValue{Value: yyDollar[1].expression}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:560
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:564
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:570
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 87:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:576
		{
			yyVAL.token = Token{}
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:580
		{
			yyVAL.token = yyDollar[1].token
		}
	case 89:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:584
		{
			yyVAL.token = yyDollar[1].token
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:590
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:596
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
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:623
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:627
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:631
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 96:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:635
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:639
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 98:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:643
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 99:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:647
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 100:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:651
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 101:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:655
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: RowValueList{RowValues: yyDollar[5].expressions}, Negation: yyDollar[2].token}
		}
	case 102:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:659
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Values: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 103:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:663
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 104:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:667
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 105:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:671
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 106:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:675
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 107:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:679
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 108:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:683
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: RowValueList{RowValues: yyDollar[5].expressions}}
		}
	case 109:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:687
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Values: yyDollar[4].expression}
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:691
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 111:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:697
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:701
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:705
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 114:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:709
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 115:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:713
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:719
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:723
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 118:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:727
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 119:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:733
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:737
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 121:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:743
		{
			yyVAL.expression = Option{}
		}
	case 122:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:747
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 123:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:751
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 124:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:757
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 125:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:761
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:767
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 127:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:771
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:775
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 129:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:781
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:785
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:791
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:795
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 133:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:799
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 134:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:803
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:807
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:811
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 137:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:817
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 138:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:821
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 139:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:825
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 140:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:829
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 141:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:833
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 142:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:839
		{
			yyVAL.expression = nil
		}
	case 143:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:843
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 144:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:847
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 145:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:853
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:857
		{
			yyVAL.expression = AllColumns{}
		}
	case 147:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:863
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 148:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:867
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 149:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:873
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 150:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:879
		{
			yyVAL.expression = nil
		}
	case 151:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:883
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 152:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:889
		{
			yyVAL.expression = nil
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:893
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:899
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 155:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:903
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 156:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:909
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 157:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:913
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:919
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 159:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:923
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 160:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:929
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 161:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:933
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
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
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 165:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:953
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
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
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:969
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 169:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:973
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 170:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:979
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 171:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:983
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 172:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:987
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 173:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:991
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 174:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:997
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 175:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1003
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1009
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 177:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1013
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 178:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1019
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 179:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1024
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 180:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1031
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 181:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:1037
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 182:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:1041
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 183:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1047
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 184:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1051
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1057
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 186:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1061
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 187:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1067
		{
			yyVAL.expression = nil
		}
	case 188:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1071
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 189:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1075
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 190:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1079
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 191:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1083
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 192:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1089
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 193:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1093
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 194:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1099
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 195:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1105
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 196:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1109
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 197:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1115
		{
			yyVAL.procexpr = nil
		}
	case 198:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1119
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 199:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1125
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 200:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1129
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 201:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1135
		{
			yyVAL.procexpr = nil
		}
	case 202:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1139
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1145
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 204:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1151
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1157
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 206:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1161
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 207:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1168
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 208:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1172
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 209:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1179
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 210:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1185
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 211:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1191
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1197
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 213:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1203
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 214:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1207
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 215:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1213
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1219
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 217:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1223
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
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
			yyVAL.token = Token{}
		}
	case 221:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1243
		{
			yyVAL.token = yyDollar[1].token
		}
	case 222:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1249
		{
			yyVAL.token = Token{}
		}
	case 223:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1253
		{
			yyVAL.token = yyDollar[1].token
		}
	case 224:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1259
		{
			yyVAL.token = Token{}
		}
	case 225:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1263
		{
			yyVAL.token = yyDollar[1].token
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
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1287
		{
			yyVAL.token = yyDollar[1].token
		}
	case 231:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1291
		{
			yyVAL.token = yyDollar[1].token
		}
	case 232:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1297
		{
			yyVAL.token = Token{}
		}
	case 233:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1301
		{
			yyVAL.token = yyDollar[1].token
		}
	case 234:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1307
		{
			yyVAL.token = yyDollar[1].token
		}
	case 235:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1311
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 236:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1317
		{
			yyVAL.token = Token{}
		}
	case 237:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1321
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
