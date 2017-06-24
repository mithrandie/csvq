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
const JOIN = 57386
const INNER = 57387
const OUTER = 57388
const LEFT = 57389
const RIGHT = 57390
const FULL = 57391
const CROSS = 57392
const ON = 57393
const USING = 57394
const NATURAL = 57395
const UNION = 57396
const INTERSECT = 57397
const EXCEPT = 57398
const ALL = 57399
const ANY = 57400
const EXISTS = 57401
const IN = 57402
const AND = 57403
const OR = 57404
const NOT = 57405
const BETWEEN = 57406
const LIKE = 57407
const IS = 57408
const NULL = 57409
const DISTINCT = 57410
const WITH = 57411
const CASE = 57412
const IF = 57413
const ELSEIF = 57414
const WHILE = 57415
const WHEN = 57416
const THEN = 57417
const ELSE = 57418
const DO = 57419
const END = 57420
const DECLARE = 57421
const CURSOR = 57422
const FOR = 57423
const FETCH = 57424
const OPEN = 57425
const CLOSE = 57426
const DISPOSE = 57427
const GROUP_CONCAT = 57428
const SEPARATOR = 57429
const COMMIT = 57430
const ROLLBACK = 57431
const CONTINUE = 57432
const BREAK = 57433
const EXIT = 57434
const PRINT = 57435
const VAR = 57436
const COMPARISON_OP = 57437
const STRING_OP = 57438
const SUBSTITUTION_OP = 57439

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

//line parser.y:1266

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
	-1, 112,
	60, 209,
	64, 209,
	65, 209,
	-2, 223,
	-1, 131,
	44, 211,
	46, 215,
	-2, 144,
	-1, 168,
	106, 105,
	-2, 207,
	-1, 170,
	74, 135,
	-2, 209,
	-1, 175,
	37, 105,
	87, 105,
	106, 105,
	-2, 207,
	-1, 188,
	60, 209,
	64, 209,
	65, 209,
	-2, 129,
	-1, 197,
	60, 209,
	64, 209,
	65, 209,
	-2, 79,
	-1, 209,
	46, 215,
	-2, 211,
	-1, 225,
	60, 209,
	64, 209,
	65, 209,
	-2, 204,
	-1, 231,
	66, 0,
	95, 0,
	98, 0,
	-2, 84,
	-1, 232,
	66, 0,
	95, 0,
	98, 0,
	-2, 85,
	-1, 266,
	60, 209,
	64, 209,
	65, 209,
	-2, 49,
	-1, 309,
	66, 0,
	95, 0,
	98, 0,
	-2, 91,
	-1, 315,
	60, 209,
	64, 209,
	65, 209,
	-2, 140,
	-1, 340,
	60, 209,
	64, 209,
	65, 209,
	-2, 162,
	-1, 368,
	78, 137,
	-2, 209,
	-1, 375,
	60, 209,
	64, 209,
	65, 209,
	-2, 53,
	-1, 393,
	60, 209,
	64, 209,
	65, 209,
	-2, 171,
	-1, 402,
	74, 152,
	76, 152,
	78, 152,
	-2, 209,
	-1, 406,
	90, 18,
	91, 18,
	-2, 1,
	-1, 409,
	60, 209,
	64, 209,
	65, 209,
	-2, 127,
}

const yyPrivate = 57344

const yyLast = 930

var yyAct = [...]int{

	41, 419, 290, 428, 315, 43, 44, 45, 46, 47,
	48, 49, 382, 2, 353, 254, 348, 314, 273, 85,
	37, 276, 37, 284, 67, 68, 69, 131, 300, 185,
	40, 1, 201, 195, 110, 361, 75, 107, 50, 94,
	112, 210, 181, 92, 354, 130, 208, 113, 249, 317,
	108, 90, 109, 77, 86, 23, 155, 23, 64, 42,
	132, 56, 392, 371, 347, 142, 115, 88, 337, 334,
	178, 279, 146, 147, 148, 269, 178, 267, 143, 57,
	57, 61, 370, 38, 3, 16, 58, 58, 127, 432,
	59, 198, 165, 418, 70, 71, 72, 73, 74, 400,
	394, 170, 391, 172, 378, 76, 163, 162, 164, 135,
	137, 154, 346, 180, 213, 336, 214, 215, 216, 211,
	184, 188, 209, 125, 312, 247, 128, 410, 58, 197,
	140, 141, 151, 59, 275, 157, 158, 159, 160, 161,
	152, 151, 138, 153, 157, 158, 159, 160, 161, 225,
	167, 59, 366, 173, 167, 168, 230, 231, 232, 307,
	183, 220, 239, 240, 241, 242, 243, 244, 245, 228,
	37, 190, 166, 207, 91, 199, 212, 159, 160, 161,
	175, 229, 58, 280, 100, 219, 179, 266, 57, 205,
	256, 138, 115, 144, 204, 58, 37, 206, 228, 226,
	227, 217, 224, 145, 434, 23, 192, 58, 193, 194,
	200, 233, 157, 158, 159, 160, 161, 426, 407, 251,
	397, 100, 102, 166, 253, 326, 80, 367, 299, 262,
	359, 23, 203, 263, 296, 248, 321, 252, 422, 306,
	303, 309, 421, 173, 302, 261, 252, 423, 318, 278,
	372, 283, 293, 422, 303, 437, 282, 319, 287, 433,
	324, 325, 289, 323, 327, 301, 292, 204, 416, 256,
	396, 322, 188, 305, 197, 37, 174, 304, 191, 103,
	58, 117, 164, 177, 121, 340, 288, 136, 291, 294,
	204, 204, 320, 335, 238, 237, 235, 331, 52, 358,
	234, 236, 120, 333, 285, 277, 386, 343, 363, 341,
	23, 298, 339, 344, 345, 101, 342, 286, 356, 281,
	119, 338, 37, 368, 171, 365, 54, 330, 295, 297,
	360, 124, 329, 362, 375, 265, 332, 105, 54, 357,
	256, 136, 374, 51, 52, 53, 37, 204, 376, 58,
	388, 250, 355, 377, 58, 122, 123, 23, 63, 381,
	393, 294, 221, 222, 204, 268, 136, 62, 59, 399,
	385, 223, 387, 404, 402, 213, 59, 214, 215, 216,
	211, 23, 149, 209, 401, 203, 218, 406, 409, 405,
	116, 38, 166, 37, 129, 55, 126, 415, 408, 379,
	412, 182, 277, 114, 229, 204, 59, 414, 413, 420,
	58, 139, 58, 424, 411, 291, 60, 425, 37, 204,
	204, 380, 427, 38, 431, 395, 37, 430, 23, 417,
	349, 350, 351, 352, 436, 256, 111, 429, 439, 136,
	39, 37, 66, 277, 403, 163, 162, 164, 59, 256,
	154, 438, 156, 23, 58, 37, 270, 389, 390, 65,
	294, 23, 308, 93, 310, 311, 59, 99, 100, 102,
	89, 103, 104, 39, 10, 38, 23, 9, 291, 152,
	151, 8, 153, 157, 158, 159, 160, 161, 7, 6,
	23, 202, 42, 5, 59, 99, 100, 102, 274, 103,
	104, 39, 213, 4, 214, 215, 216, 316, 136, 169,
	82, 186, 213, 136, 214, 215, 216, 211, 383, 384,
	209, 97, 99, 100, 102, 98, 103, 104, 187, 105,
	134, 133, 96, 59, 99, 100, 102, 95, 103, 104,
	39, 81, 84, 78, 83, 79, 196, 118, 106, 97,
	328, 264, 36, 98, 15, 257, 14, 105, 13, 12,
	96, 11, 101, 163, 162, 164, 255, 87, 154, 136,
	0, 136, 59, 99, 100, 102, 106, 103, 104, 39,
	0, 0, 0, 0, 105, 0, 0, 0, 97, 0,
	101, 189, 98, 0, 0, 87, 105, 152, 151, 96,
	153, 157, 158, 159, 160, 161, 0, 0, 246, 0,
	0, 0, 0, 136, 0, 106, 0, 101, 0, 0,
	271, 272, 0, 0, 0, 0, 0, 97, 0, 101,
	313, 98, 0, 0, 87, 105, 0, 0, 96, 0,
	163, 162, 164, 0, 39, 154, 38, 0, 18, 34,
	19, 0, 17, 0, 106, 163, 162, 164, 20, 0,
	154, 21, 0, 0, 0, 0, 0, 0, 101, 435,
	0, 0, 0, 87, 152, 151, 0, 153, 157, 158,
	159, 160, 161, 0, 0, 0, 0, 0, 0, 152,
	151, 0, 153, 157, 158, 159, 160, 161, 0, 0,
	163, 162, 164, 0, 258, 154, 32, 0, 0, 0,
	0, 0, 26, 0, 398, 30, 27, 28, 29, 0,
	0, 24, 25, 259, 260, 33, 35, 22, 0, 0,
	0, 163, 162, 164, 152, 151, 154, 153, 157, 158,
	159, 160, 161, 0, 0, 373, 163, 162, 164, 0,
	0, 154, 0, 0, 0, 0, 0, 0, 0, 0,
	369, 0, 163, 162, 164, 152, 151, 154, 153, 157,
	158, 159, 160, 161, 163, 162, 164, 0, 176, 154,
	152, 151, 0, 153, 157, 158, 159, 160, 161, 0,
	0, 0, 0, 163, 162, 164, 152, 151, 154, 153,
	157, 158, 159, 160, 161, 0, 0, 150, 152, 151,
	0, 153, 157, 158, 159, 160, 161, 364, 162, 164,
	0, 0, 154, 0, 0, 0, 0, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 39, 0, 38, 0,
	18, 34, 19, 163, 17, 164, 0, 0, 154, 0,
	20, 152, 151, 21, 153, 157, 158, 159, 160, 161,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 152, 151, 0,
	153, 157, 158, 159, 160, 161, 0, 0, 0, 164,
	0, 0, 154, 0, 0, 0, 31, 0, 32, 0,
	0, 0, 0, 0, 26, 0, 0, 30, 27, 28,
	29, 0, 0, 24, 25, 0, 0, 33, 35, 22,
	0, 152, 151, 0, 153, 157, 158, 159, 160, 161,
}
var yyPact = [...]int{

	825, -1000, 825, -49, -49, -49, -49, -49, -49, -49,
	-49, -1000, -1000, -1000, -1000, -1000, 289, 375, 444, 402,
	338, 329, 431, -49, -49, -49, 444, 444, 444, 444,
	444, 568, 568, -49, 424, 568, 389, 95, 213, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	277, 227, 227, 227, 291, 444, 380, -19, 372, -1000,
	86, 397, 444, 444, -49, -29, 96, -1000, -1000, -1000,
	123, -49, -49, -49, 362, 732, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 95, -1000, 462, 50, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 568, 37, 568, -1000,
	-1000, 215, -1000, -1000, -1000, -1000, 75, 701, 223, -31,
	-1000, 88, 384, 383, 86, 568, 490, -1000, -1000, 178,
	410, -1000, 410, 410, 568, 70, 444, 444, -1000, 444,
	383, 69, -1000, 364, -1000, -1000, -1000, -1000, 410, 86,
	56, 336, -1000, 431, 568, 118, -1000, -1000, -1000, 429,
	825, 568, 568, 568, 219, 236, 237, 568, 568, 568,
	568, 568, 568, 568, -1000, 502, 19, 444, 213, 163,
	713, -1000, 826, -1000, -1000, 213, 633, 444, 429, 517,
	-1000, 297, 568, -1000, 713, -1000, -30, 343, 713, -1000,
	-1000, 178, 243, -1000, 243, -1000, -32, 579, 29, 444,
	-1000, 389, -36, 85, 46, -1000, -1000, -1000, 275, 457,
	258, 273, 86, -1000, -1000, -1000, -1000, -1000, 444, 383,
	444, 147, 129, 444, -1000, 713, 410, -49, -37, 168,
	113, 36, 36, 270, 568, 54, 568, 37, 37, 76,
	76, -1000, -1000, -1000, 782, 826, -1000, -1000, -1000, 18,
	529, 172, 568, 301, 158, 633, -1000, -1000, 568, -49,
	-49, 148, -1000, -49, 293, 287, 713, 490, 444, 568,
	-1000, -1000, -1000, -1000, -38, 568, 9, -39, 383, 444,
	568, 86, 272, 258, 269, -1000, 86, -1000, -1000, -1000,
	6, -43, 400, 444, 318, -1000, 444, 303, -49, -1000,
	152, 168, 825, 568, -1000, -1000, 756, 462, -1000, 36,
	-1000, -1000, -1000, -1000, -1000, 45, 149, 163, 568, 685,
	-24, 177, -1000, 670, -1000, -1000, 633, -1000, -1000, 568,
	568, -1000, -1000, -1000, 29, -2, 378, 444, -1000, -1000,
	713, 467, 86, 262, 86, 330, -1000, 444, -1000, -1000,
	-1000, 444, 444, -4, -45, 568, -6, 444, -1000, 199,
	142, 182, -1000, 639, 568, -7, 568, -1000, 713, 568,
	-1000, 439, -49, 633, 140, 713, -1000, -1000, -1000, 29,
	-1000, -1000, -1000, 568, 22, 330, 86, 467, -1000, -1000,
	-1000, 400, 444, 713, -1000, -1000, -49, 197, 825, 826,
	-1000, -1000, 713, -13, -1000, 166, 825, 174, -1000, 713,
	444, 330, -1000, -1000, -1000, -1000, -49, -1000, -1000, 139,
	166, 633, 568, -49, -17, -1000, 188, 126, 181, -1000,
	594, -1000, -1000, -49, 184, 633, -1000, -49, -1000, -1000,
}
var yyPgo = [...]int{

	0, 30, 15, 13, 566, 561, 559, 558, 556, 555,
	554, 84, 85, 552, 47, 42, 551, 550, 38, 547,
	53, 105, 4, 546, 226, 545, 544, 543, 542, 541,
	48, 537, 60, 531, 27, 530, 12, 528, 511, 510,
	509, 507, 21, 17, 33, 45, 61, 2, 29, 49,
	503, 498, 18, 493, 491, 32, 489, 488, 481, 44,
	14, 16, 477, 474, 35, 28, 3, 1, 67, 470,
	51, 174, 43, 463, 39, 19, 50, 54, 459, 58,
	351, 56, 456, 46, 23, 41, 302, 452, 0,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4,
	5, 5, 6, 6, 7, 7, 7, 7, 7, 8,
	8, 8, 8, 8, 9, 9, 9, 9, 10, 10,
	11, 12, 12, 12, 12, 13, 14, 14, 15, 15,
	16, 16, 17, 17, 18, 18, 19, 19, 20, 20,
	20, 20, 20, 20, 21, 21, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 23, 82,
	82, 82, 24, 25, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 27, 27, 27, 27, 27,
	28, 28, 28, 29, 29, 30, 30, 30, 31, 31,
	32, 32, 32, 33, 33, 34, 34, 34, 34, 34,
	34, 35, 35, 35, 35, 35, 36, 36, 36, 37,
	37, 38, 38, 39, 40, 40, 41, 41, 42, 42,
	43, 43, 44, 44, 45, 45, 46, 46, 47, 47,
	48, 48, 49, 49, 50, 50, 50, 50, 51, 52,
	52, 53, 54, 55, 55, 56, 56, 57, 58, 58,
	59, 59, 60, 60, 61, 61, 61, 61, 61, 62,
	62, 63, 64, 64, 65, 65, 66, 66, 67, 67,
	68, 69, 70, 70, 71, 71, 72, 73, 74, 75,
	76, 76, 77, 78, 78, 79, 79, 80, 80, 81,
	81, 83, 83, 84, 84, 85, 85, 85, 85, 86,
	86, 87, 87, 88, 88,
}
var yyR2 = [...]int{

	0, 0, 2, 0, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
	3, 2, 2, 2, 6, 3, 3, 3, 5, 8,
	9, 7, 9, 2, 8, 9, 2, 2, 5, 3,
	3, 5, 4, 4, 4, 3, 0, 2, 0, 2,
	0, 3, 0, 2, 0, 3, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 3, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 3, 2, 0,
	1, 1, 3, 3, 3, 3, 4, 4, 6, 6,
	4, 4, 4, 4, 2, 3, 3, 3, 3, 3,
	3, 3, 2, 4, 1, 0, 2, 2, 5, 7,
	1, 2, 3, 1, 1, 1, 1, 2, 3, 1,
	1, 5, 5, 6, 6, 4, 0, 2, 4, 1,
	1, 1, 3, 5, 0, 1, 0, 2, 1, 3,
	1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
	1, 3, 4, 2, 5, 8, 4, 7, 3, 1,
	3, 6, 3, 1, 3, 4, 5, 6, 6, 8,
	1, 3, 1, 3, 0, 1, 1, 2, 2, 5,
	7, 7, 4, 2, 0, 2, 4, 2, 0, 2,
	1, 1, 1, 2, 1, 2, 1, 1, 1, 1,
	1, 3, 3, 1, 3, 1, 3, 0, 1, 0,
	1, 0, 1, 0, 1, 0, 1, 1, 1, 0,
	1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -3, -11, -50, -53, -56, -57, -58, -62,
	-63, -5, -6, -7, -8, -10, -12, 19, 15, 17,
	25, 28, 94, -77, 88, 89, 79, 83, 84, 85,
	82, 71, 73, 92, 16, 93, -13, -75, 13, 11,
	-1, -88, 108, -88, -88, -88, -88, -88, -88, -88,
	-18, 54, 55, 56, 37, 20, -46, -32, -68, 4,
	14, -46, 29, 29, -79, -78, 11, -88, -88, -88,
	-68, -68, -68, -68, -68, -22, -21, -20, -27, -25,
	-24, -29, -39, -26, -28, -75, -77, 105, -68, -69,
	-70, -71, -72, -73, -74, -31, 70, 59, 63, 5,
	6, 100, 7, 9, 10, 67, 86, -22, -76, -75,
	-88, 12, -22, -14, 14, 97, -80, 68, -19, 43,
	-86, 57, -86, -86, 40, -68, 16, 107, -68, 22,
	-45, -34, -32, -33, -35, 23, -24, 24, 105, 14,
	-68, -68, -88, 107, 97, 80, -88, -88, -88, 20,
	75, 96, 95, 98, 66, -81, -87, 99, 100, 101,
	102, 103, 62, 61, 63, -22, -11, 104, 105, -40,
	-22, -24, -22, -70, -71, 105, 77, 60, 107, 98,
	-88, -15, 18, -45, -22, -48, -38, -37, -22, 101,
	-70, 100, -12, -12, -12, -44, -23, -22, 21, 105,
	-11, -55, -54, -21, -68, -46, -68, -15, -83, 53,
	-85, 50, 107, 45, 47, 48, 49, -68, 22, -45,
	105, 26, 27, 35, -79, -22, 81, -76, -75, -1,
	-22, -22, -22, -81, 64, 60, 65, 58, 57, -22,
	-22, -22, -22, -22, -22, -22, 106, 106, -68, -30,
	-80, -49, 74, -30, -2, -4, -3, -9, 71, 90,
	91, -68, -76, -20, -16, 38, -22, 107, 22, 107,
	-82, 41, 42, -52, -51, 105, -42, -21, -14, 107,
	98, 44, -83, -85, -84, 46, 44, -45, -68, -15,
	-47, -68, -59, 105, -68, -21, 105, -21, -11, -88,
	-65, -64, 76, 72, -72, -74, -22, 105, -24, -22,
	-24, -24, 106, 101, -43, -22, -41, -49, 76, -22,
	-18, 78, -2, -22, -88, -88, 77, -88, -17, 39,
	40, -48, -68, -44, 107, -43, 106, 107, -15, -55,
	-22, -34, 44, -84, 44, -34, 106, 107, -61, 30,
	31, 32, 33, -60, -59, 34, -42, 36, -88, 78,
	-65, -64, -1, -22, 61, -43, 107, 78, -22, 75,
	106, 87, 73, 75, -2, -22, -43, -52, 106, 21,
	-11, -42, -36, 51, 52, -34, 44, -34, -47, -21,
	-21, 106, 107, -22, 106, -68, 71, 78, 75, -22,
	106, -43, -22, 5, -88, -2, -3, 78, -52, -22,
	105, -34, -36, -61, -60, -88, 71, -1, 106, -67,
	-66, 76, 72, 73, -47, -88, 78, -67, -66, -2,
	-22, -88, 106, 71, 78, 75, -88, 71, -2, -88,
}
var yyDef = [...]int{

	1, -2, 1, 223, 223, 223, 223, 223, 223, 223,
	223, 13, 14, 15, 16, 17, 54, 0, 0, 0,
	0, 0, 0, 223, 223, 223, 0, 0, 0, 0,
	0, 0, 0, 223, 0, 0, 46, 0, 207, 199,
	2, 5, 224, 6, 7, 8, 9, 10, 11, 12,
	56, 219, 219, 219, 0, 0, 0, 146, 110, 190,
	0, 0, 0, 0, 223, 205, 203, 21, 22, 23,
	0, 223, 223, 223, 0, 209, 66, 67, 68, 69,
	70, 71, 72, 73, 74, 75, 76, 0, 64, 58,
	59, 60, 61, 62, 63, 104, 134, 0, 0, 191,
	192, 0, 194, 196, 197, 198, 0, 209, 0, 75,
	33, 0, -2, 48, 0, 0, 0, 208, 40, 0,
	0, 220, 0, 0, 0, 0, 0, 0, 111, 0,
	48, -2, 115, 116, 119, 120, 113, 114, 0, 0,
	0, 0, 20, 0, 0, 0, 25, 26, 27, 0,
	1, 0, 221, 222, 209, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 210, 209, 0, 0, -2, 0,
	-2, 94, 102, 193, 195, -2, 3, 0, 0, 0,
	39, 50, 0, 47, 202, 45, 150, 131, -2, 130,
	57, 0, 42, 43, 44, 55, 142, -2, 0, 0,
	156, 46, 163, 0, 64, 147, 112, 165, 0, -2,
	213, 0, 0, 212, 216, 217, 218, 117, 0, 48,
	0, 0, 0, 0, 206, -2, 0, 223, 200, 184,
	83, -2, -2, 0, 0, 0, 0, 0, 0, 95,
	96, 97, 98, 99, 100, 101, 77, 82, 65, 0,
	0, 136, 0, 54, 0, 3, 18, 19, 0, 223,
	223, 0, 201, 223, 52, 0, -2, 0, 0, 0,
	78, 80, 81, 154, 159, 0, 0, 138, 48, 0,
	0, 0, 0, 213, 0, 214, 0, 145, 118, 166,
	0, 148, 174, 0, 170, 179, 0, 0, 223, 28,
	0, 184, 1, 0, 86, 87, 209, 0, 90, -2,
	92, 93, 103, 106, 107, -2, 0, 153, 0, 209,
	0, 0, 4, 209, 36, 37, 3, 38, 41, 0,
	0, 151, 132, 143, 0, 0, 0, 0, 161, 164,
	-2, 126, 0, 0, 0, 125, 167, 0, 168, 175,
	176, 0, 0, 0, 172, 0, 0, 0, 24, 0,
	0, 183, 185, 209, 0, 0, 0, 133, -2, 0,
	108, 0, 223, 1, 0, -2, 51, 160, 158, 0,
	157, 139, 121, 0, 0, 122, 0, 126, 149, 177,
	178, 174, 0, -2, 180, 181, 223, 0, 1, 88,
	89, 141, -2, 0, 31, 188, -2, 0, 155, -2,
	0, 124, 123, 169, 173, 29, 223, 182, 109, 0,
	188, 3, 0, 223, 0, 30, 0, 0, 187, 189,
	209, 32, 128, 223, 0, 3, 34, 223, 186, 35,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 103, 3, 3,
	105, 106, 101, 99, 107, 100, 104, 102, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 108,
	3, 98,
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
	92, 93, 94, 95, 96, 97,
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
		//line parser.y:146
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:151
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:158
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:163
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:170
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:174
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:178
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:182
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:186
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:190
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:194
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:198
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:202
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:206
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:210
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:214
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:218
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:224
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:228
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:234
		{
			yyVAL.statement = VariableDeclaration{Assignments: yyDollar[2].expressions}
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:238
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:244
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:248
		{
			yyVAL.statement = TransactionControl{Token: yyDollar[1].token.Token}
		}
	case 24:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:254
		{
			yyVAL.statement = CursorDeclaration{Cursor: yyDollar[2].identifier, Query: yyDollar[5].expression.(SelectQuery)}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:258
		{
			yyVAL.statement = OpenCursor{Cursor: yyDollar[2].identifier}
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:262
		{
			yyVAL.statement = CloseCursor{Cursor: yyDollar[2].identifier}
		}
	case 27:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:266
		{
			yyVAL.statement = DisposeCursor{Cursor: yyDollar[2].identifier}
		}
	case 28:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:270
		{
			yyVAL.statement = FetchCursor{Cursor: yyDollar[2].identifier, Variables: yyDollar[4].variables}
		}
	case 29:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:276
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 30:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:280
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 31:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:284
		{
			yyVAL.statement = While{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}
		}
	case 32:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:288
		{
			yyVAL.statement = WhileInCursor{Variables: yyDollar[2].variables, Cursor: yyDollar[4].identifier, Statements: yyDollar[6].program}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:292
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 34:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:298
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, Else: yyDollar[5].procexpr}
		}
	case 35:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.y:302
		{
			yyVAL.statement = If{Condition: yyDollar[2].expression, Statements: yyDollar[4].program, ElseIf: yyDollar[5].procexprs, Else: yyDollar[6].procexpr}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:306
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:310
		{
			yyVAL.statement = FlowControl{Token: yyDollar[1].token.Token}
		}
	case 38:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:316
		{
			yyVAL.statement = SetFlag{Name: yyDollar[2].token.Literal, Value: yyDollar[4].primary}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:320
		{
			yyVAL.statement = Print{Value: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:326
		{
			yyVAL.expression = SelectQuery{
				SelectEntity:  yyDollar[1].expression,
				OrderByClause: yyDollar[2].expression,
				LimitClause:   yyDollar[3].expression,
			}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:336
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
		//line parser.y:346
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
		//line parser.y:355
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
		//line parser.y:364
		{
			yyVAL.expression = SelectSet{
				LHS:      yyDollar[1].expression,
				Operator: yyDollar[2].token,
				All:      yyDollar[3].token,
				RHS:      yyDollar[4].expression,
			}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:375
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 46:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:381
		{
			yyVAL.expression = nil
		}
	case 47:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:385
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:391
		{
			yyVAL.expression = nil
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:401
		{
			yyVAL.expression = nil
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:405
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 52:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:411
		{
			yyVAL.expression = nil
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:415
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:421
		{
			yyVAL.expression = nil
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:425
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:431
		{
			yyVAL.expression = nil
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:435
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:441
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:445
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:449
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:453
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:457
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:461
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:467
		{
			yyVAL.expression = FieldReference{Column: yyDollar[1].identifier}
		}
	case 65:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:471
		{
			yyVAL.expression = FieldReference{View: yyDollar[1].identifier, Column: yyDollar[3].identifier}
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:477
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:481
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:485
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:489
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:493
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:497
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:501
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:505
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:509
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:513
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:517
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:521
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 78:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:527
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 79:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:533
		{
			yyVAL.token = Token{}
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:537
		{
			yyVAL.token = yyDollar[1].token
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:541
		{
			yyVAL.token = yyDollar[1].token
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:547
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:553
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
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:580
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: Token{Token: COMPARISON_OP, Literal: "="}, RHS: yyDollar[3].expression}
		}
	case 86:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:584
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 87:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:588
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 88:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 89:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 90:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:600
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 91:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:604
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 92:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:608
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 93:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:616
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:622
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:626
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:630
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:634
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:638
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:644
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:648
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:652
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 103:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:658
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:662
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 105:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:668
		{
			yyVAL.expression = Option{}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:672
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:676
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 108:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:682
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 109:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:686
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:692
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:696
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 112:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:700
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:706
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:710
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:716
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:720
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 117:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:724
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 118:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:728
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:732
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:736
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 121:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:742
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 122:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:746
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 123:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:750
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 124:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:754
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 125:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:758
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 126:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:764
		{
			yyVAL.expression = nil
		}
	case 127:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:768
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 128:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:772
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 129:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:778
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:782
		{
			yyVAL.expression = AllColumns{}
		}
	case 131:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:788
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 132:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:792
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 133:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:798
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 134:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:804
		{
			yyVAL.expression = nil
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:808
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 136:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:814
		{
			yyVAL.expression = nil
		}
	case 137:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:818
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 138:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:824
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 139:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:828
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:834
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 141:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:838
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 142:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:844
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 143:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:848
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:854
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 145:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:858
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:864
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 147:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:868
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:874
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 149:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:878
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 150:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:884
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 151:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:888
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 152:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:894
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 153:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:898
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 154:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:904
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Values: yyDollar[4].token.Literal, ValuesList: yyDollar[5].expressions}
		}
	case 155:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:908
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Values: yyDollar[7].token.Literal, ValuesList: yyDollar[8].expressions}
		}
	case 156:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:912
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Query: yyDollar[4].expression.(SelectQuery)}
		}
	case 157:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:916
		{
			yyVAL.expression = InsertQuery{Insert: yyDollar[1].token.Literal, Into: yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions, Query: yyDollar[7].expression.(SelectQuery)}
		}
	case 158:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:922
		{
			yyVAL.expression = InsertValues{Values: yyDollar[2].expressions}
		}
	case 159:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:928
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 160:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:932
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 161:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:938
		{
			yyVAL.expression = UpdateQuery{Update: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, Set: yyDollar[3].token.Literal, SetList: yyDollar[4].expressions, FromClause: yyDollar[5].expression, WhereClause: yyDollar[6].expression}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:944
		{
			yyVAL.expression = UpdateSet{Field: yyDollar[1].expression.(FieldReference), Value: yyDollar[3].expression}
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:950
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 164:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:954
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 165:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:960
		{
			from := FromClause{From: yyDollar[2].token.Literal, Tables: yyDollar[3].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, FromClause: from, WhereClause: yyDollar[4].expression}
		}
	case 166:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:965
		{
			from := FromClause{From: yyDollar[3].token.Literal, Tables: yyDollar[4].expressions}
			yyVAL.expression = DeleteQuery{Delete: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions, FromClause: from, WhereClause: yyDollar[5].expression}
		}
	case 167:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:972
		{
			yyVAL.expression = CreateTable{CreateTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Fields: yyDollar[5].expressions}
		}
	case 168:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:978
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}, Position: yyDollar[6].expression}
		}
	case 169:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.y:982
		{
			yyVAL.expression = AddColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Add: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions, Position: yyDollar[8].expression}
		}
	case 170:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:988
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier}
		}
	case 171:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:992
		{
			yyVAL.expression = ColumnDefault{Column: yyDollar[1].identifier, Default: yyDollar[2].token.Literal, Value: yyDollar[3].expression}
		}
	case 172:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:998
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 173:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1002
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 174:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1008
		{
			yyVAL.expression = nil
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1012
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1016
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token}
		}
	case 177:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1020
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 178:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1024
		{
			yyVAL.expression = ColumnPosition{Position: yyDollar[1].token, Column: yyDollar[2].expression}
		}
	case 179:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:1030
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: []Expression{yyDollar[5].expression}}
		}
	case 180:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1034
		{
			yyVAL.expression = DropColumns{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Drop: yyDollar[4].token.Literal, Columns: yyDollar[6].expressions}
		}
	case 181:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:1040
		{
			yyVAL.expression = RenameColumn{AlterTable: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Table: yyDollar[3].identifier, Rename: yyDollar[4].token.Literal, Old: yyDollar[5].expression.(FieldReference), To: yyDollar[6].token.Literal, New: yyDollar[7].identifier}
		}
	case 182:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1046
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 183:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1050
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 184:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1056
		{
			yyVAL.procexpr = nil
		}
	case 185:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1060
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 186:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:1066
		{
			yyVAL.procexprs = []ProcExpr{ElseIf{Condition: yyDollar[2].expression, Statements: yyDollar[4].program}}
		}
	case 187:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1070
		{
			yyVAL.procexprs = append(yyDollar[1].procexprs, yyDollar[2].procexprs...)
		}
	case 188:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1076
		{
			yyVAL.procexpr = nil
		}
	case 189:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1080
		{
			yyVAL.procexpr = Else{Statements: yyDollar[2].program}
		}
	case 190:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1086
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 191:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1092
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 192:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1098
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 193:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1102
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1109
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 195:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:1113
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 196:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1120
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 197:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1126
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 198:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1132
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 199:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1138
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1144
		{
			yyVAL.variables = []Variable{yyDollar[1].variable}
		}
	case 201:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1148
		{
			yyVAL.variables = append([]Variable{yyDollar[1].variable}, yyDollar[3].variables...)
		}
	case 202:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1154
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1160
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 204:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1164
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1170
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 206:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:1174
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 207:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1180
		{
			yyVAL.token = Token{}
		}
	case 208:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1184
		{
			yyVAL.token = yyDollar[1].token
		}
	case 209:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1190
		{
			yyVAL.token = Token{}
		}
	case 210:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1194
		{
			yyVAL.token = yyDollar[1].token
		}
	case 211:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1200
		{
			yyVAL.token = Token{}
		}
	case 212:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1204
		{
			yyVAL.token = yyDollar[1].token
		}
	case 213:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1210
		{
			yyVAL.token = Token{}
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1214
		{
			yyVAL.token = yyDollar[1].token
		}
	case 215:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1220
		{
			yyVAL.token = Token{}
		}
	case 216:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1224
		{
			yyVAL.token = yyDollar[1].token
		}
	case 217:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1228
		{
			yyVAL.token = yyDollar[1].token
		}
	case 218:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1232
		{
			yyVAL.token = yyDollar[1].token
		}
	case 219:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1238
		{
			yyVAL.token = Token{}
		}
	case 220:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1242
		{
			yyVAL.token = yyDollar[1].token
		}
	case 221:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1248
		{
			yyVAL.token = yyDollar[1].token
		}
	case 222:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1252
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 223:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:1258
		{
			yyVAL.token = Token{}
		}
	case 224:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:1262
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
