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
	primary     Primary
	identifier  Identifier
	text        String
	integer     Integer
	float       Float
	ternary     Ternary
	datetime    Datetime
	null        Null
	variable    Variable
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
const SELECT = 57354
const FROM = 57355
const UPDATE = 57356
const SET = 57357
const DELETE = 57358
const WHERE = 57359
const INSERT = 57360
const INTO = 57361
const VALUES = 57362
const AS = 57363
const DUAL = 57364
const STDIN = 57365
const CREATE = 57366
const DROP = 57367
const ALTER = 57368
const TABLE = 57369
const COLUMN = 57370
const ORDER = 57371
const GROUP = 57372
const HAVING = 57373
const BY = 57374
const ASC = 57375
const DESC = 57376
const LIMIT = 57377
const JOIN = 57378
const INNER = 57379
const OUTER = 57380
const LEFT = 57381
const RIGHT = 57382
const FULL = 57383
const CROSS = 57384
const ON = 57385
const USING = 57386
const NATURAL = 57387
const UNION = 57388
const ALL = 57389
const ANY = 57390
const EXISTS = 57391
const IN = 57392
const AND = 57393
const OR = 57394
const NOT = 57395
const BETWEEN = 57396
const LIKE = 57397
const IS = 57398
const NULL = 57399
const DISTINCT = 57400
const WITH = 57401
const TRUE = 57402
const FALSE = 57403
const UNKNOWN = 57404
const CASE = 57405
const WHEN = 57406
const THEN = 57407
const ELSE = 57408
const END = 57409
const GROUP_CONCAT = 57410
const SEPARATOR = 57411
const VAR = 57412
const COMPARISON_OP = 57413
const STRING_OP = 57414
const SUBSTITUTION_OP = 57415

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
	"DROP",
	"ALTER",
	"TABLE",
	"COLUMN",
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
	"TRUE",
	"FALSE",
	"UNKNOWN",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"END",
	"GROUP_CONCAT",
	"SEPARATOR",
	"VAR",
	"COMPARISON_OP",
	"STRING_OP",
	"SUBSTITUTION_OP",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'('",
	"')'",
	"','",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:787

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
	-1, 26,
	36, 127,
	38, 131,
	-2, 99,
	-1, 71,
	50, 125,
	54, 125,
	55, 125,
	-2, 86,
	-1, 75,
	50, 125,
	54, 125,
	55, 125,
	-2, 12,
	-1, 77,
	38, 131,
	-2, 127,
	-1, 89,
	50, 125,
	54, 125,
	55, 125,
	-2, 119,
	-1, 102,
	80, 65,
	-2, 123,
	-1, 105,
	64, 92,
	-2, 125,
	-1, 110,
	29, 65,
	69, 65,
	80, 65,
	-2, 123,
	-1, 150,
	50, 125,
	54, 125,
	55, 125,
	-2, 16,
	-1, 152,
	50, 125,
	54, 125,
	55, 125,
	-2, 95,
	-1, 187,
	67, 94,
	-2, 125,
	-1, 195,
	50, 125,
	54, 125,
	55, 125,
	-2, 40,
	-1, 197,
	50, 125,
	54, 125,
	55, 125,
	-2, 84,
	-1, 203,
	64, 105,
	66, 105,
	67, 105,
	-2, 125,
}

const yyPrivate = 57344

const yyLast = 439

var yyAct = [...]int{

	152, 38, 209, 193, 26, 51, 55, 178, 148, 53,
	119, 140, 68, 93, 170, 33, 14, 214, 30, 205,
	101, 37, 71, 92, 151, 75, 111, 35, 213, 85,
	211, 202, 25, 29, 32, 42, 166, 89, 91, 90,
	124, 94, 95, 96, 97, 98, 100, 99, 101, 198,
	103, 92, 31, 190, 100, 99, 101, 163, 105, 92,
	107, 34, 110, 61, 189, 102, 91, 90, 108, 94,
	95, 96, 97, 98, 91, 90, 177, 94, 95, 96,
	97, 98, 30, 142, 100, 99, 101, 78, 123, 92,
	34, 125, 126, 76, 106, 133, 134, 135, 136, 137,
	138, 139, 20, 36, 91, 90, 129, 94, 95, 96,
	97, 98, 71, 122, 147, 150, 31, 186, 30, 143,
	144, 153, 145, 30, 146, 64, 157, 92, 144, 155,
	171, 162, 192, 165, 61, 63, 161, 22, 81, 160,
	82, 83, 84, 79, 52, 172, 77, 96, 97, 98,
	92, 101, 31, 120, 173, 175, 30, 31, 30, 181,
	114, 183, 182, 158, 159, 118, 168, 164, 11, 10,
	156, 117, 187, 66, 131, 92, 154, 195, 130, 132,
	197, 191, 80, 121, 30, 201, 116, 199, 185, 203,
	31, 200, 31, 94, 95, 96, 97, 98, 108, 74,
	210, 17, 196, 62, 176, 115, 195, 109, 81, 212,
	82, 83, 84, 149, 112, 33, 210, 215, 31, 33,
	60, 61, 63, 92, 64, 65, 11, 8, 33, 60,
	61, 63, 86, 64, 65, 11, 10, 88, 91, 90,
	24, 94, 95, 96, 97, 98, 16, 10, 33, 60,
	61, 63, 19, 64, 65, 11, 141, 204, 33, 13,
	206, 5, 128, 127, 58, 18, 1, 21, 59, 12,
	54, 50, 66, 58, 169, 104, 44, 59, 57, 48,
	6, 66, 6, 67, 69, 47, 9, 57, 9, 70,
	62, 28, 67, 58, 49, 87, 4, 59, 4, 62,
	27, 66, 56, 49, 43, 46, 40, 57, 33, 60,
	61, 63, 67, 64, 65, 11, 45, 41, 194, 62,
	167, 207, 208, 49, 39, 33, 60, 61, 63, 174,
	64, 65, 11, 113, 73, 23, 15, 7, 3, 100,
	99, 101, 2, 81, 92, 82, 83, 84, 79, 179,
	180, 77, 0, 58, 0, 0, 0, 59, 0, 91,
	90, 66, 94, 95, 96, 97, 98, 57, 0, 0,
	58, 0, 67, 0, 59, 0, 0, 0, 66, 62,
	72, 0, 0, 49, 57, 100, 99, 101, 0, 67,
	92, 0, 0, 184, 99, 101, 62, 0, 92, 188,
	49, 0, 100, 0, 101, 91, 90, 92, 94, 95,
	96, 97, 98, 91, 90, 0, 94, 95, 96, 97,
	98, 0, 91, 90, 0, 94, 95, 96, 97, 98,
	81, 0, 82, 83, 84, 79, 0, 0, 77,
}
var yyPact = [...]int{

	157, -1000, 157, -66, -1000, -1000, -1000, 233, 241, 29,
	79, -1000, -1000, -1000, -1000, 223, 11, -1000, -54, 30,
	321, 304, -1000, 169, 321, -1000, 101, 211, -1000, -1000,
	-1000, -1000, -1000, -1000, 235, 241, 321, 33, -14, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 29, -1000, 224,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 321, -18, 321,
	-1000, -1000, 128, -1000, -1000, -1000, -1000, -17, -1000, -55,
	193, 33, -1000, 129, 173, 33, 150, 171, 115, 147,
	11, -1000, -1000, -1000, -1000, -1000, 254, -40, -1000, 33,
	321, 215, 98, 124, 321, 321, 321, 321, 321, 321,
	321, -1000, 79, 3, 56, 33, -1000, 167, -1000, -1000,
	79, 304, 254, 184, 321, 321, 11, 140, 115, 134,
	-1000, 11, -1000, -1000, -1000, 119, 119, -18, -18, 116,
	321, -22, 321, 71, 71, 94, 94, 94, 351, -33,
	-44, 244, -1000, 64, 321, 184, -1000, -1000, 120, 172,
	33, -1000, -5, 306, 11, 126, 11, 393, -1000, -1000,
	-1000, -1000, 342, 224, -1000, 33, -1000, -1000, -1000, 50,
	56, 321, 334, -16, -1000, 57, 321, 321, -1000, 321,
	-30, 393, 11, 306, 321, -49, -1000, 33, 321, -1000,
	252, -1000, 57, -1000, -62, 288, -1000, 33, 254, 393,
	-1000, -33, -1000, 33, -50, 321, -1000, -1000, -1000, -52,
	-64, -1000, -1000, -1000, 254, -1000,
}
var yyPgo = [...]int{

	0, 266, 342, 338, 295, 337, 336, 335, 334, 333,
	8, 329, 324, 0, 318, 35, 317, 316, 306, 305,
	304, 11, 302, 300, 4, 291, 7, 289, 284, 276,
	275, 274, 24, 3, 32, 2, 12, 14, 1, 271,
	5, 144, 9, 270, 6, 285, 279, 265, 261, 201,
	256, 13, 260, 93, 10, 87, 259,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 3, 3, 4, 5, 6,
	6, 7, 7, 8, 8, 9, 9, 10, 10, 11,
	11, 12, 12, 12, 12, 12, 12, 13, 13, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 14,
	52, 52, 52, 15, 16, 17, 17, 17, 17, 17,
	17, 17, 17, 17, 17, 18, 18, 18, 18, 18,
	19, 19, 19, 20, 20, 21, 21, 21, 22, 22,
	23, 23, 23, 24, 24, 24, 24, 24, 25, 25,
	25, 25, 25, 26, 26, 26, 27, 27, 28, 28,
	29, 30, 30, 31, 31, 32, 32, 33, 33, 34,
	34, 35, 35, 36, 36, 37, 37, 38, 39, 40,
	40, 41, 41, 42, 43, 44, 45, 46, 47, 47,
	48, 49, 49, 50, 50, 51, 51, 53, 53, 54,
	54, 55, 55, 55, 55, 56, 56,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 1, 1, 1, 7, 3, 0,
	2, 0, 2, 0, 3, 0, 2, 0, 3, 0,
	2, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 3, 2,
	0, 1, 1, 3, 3, 3, 4, 4, 6, 6,
	4, 4, 4, 4, 2, 3, 3, 3, 3, 3,
	3, 3, 2, 4, 1, 0, 2, 2, 5, 7,
	1, 1, 1, 1, 2, 3, 1, 1, 5, 5,
	6, 6, 4, 0, 2, 4, 1, 1, 1, 3,
	5, 0, 1, 0, 2, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 4, 2, 1, 1, 1,
	2, 1, 2, 1, 1, 1, 1, 3, 1, 3,
	2, 1, 3, 0, 1, 0, 1, 0, 1, 0,
	1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -48, -46, -5, 70, -45,
	12, 11, -1, -56, 82, -6, 13, -49, -47, 11,
	73, -50, 58, -7, 17, -34, -24, -23, -25, 22,
	-38, -15, 23, 4, 79, 81, 73, -13, -38, -12,
	-18, -16, -15, -20, -29, -17, -19, -45, -46, 79,
	-39, -40, -41, -42, -43, -44, -22, 63, 49, 53,
	5, 6, 75, 7, 9, 10, 57, 68, -36, -28,
	-27, -13, 76, -8, 30, -13, -53, 45, -55, 42,
	81, 37, 39, 40, 41, -38, 21, -4, -49, -13,
	72, 71, 56, -51, 74, 75, 76, 77, 78, 52,
	51, 53, 79, -13, -30, -13, -15, -13, -40, -41,
	79, 81, 21, -9, 31, 32, 36, -53, -55, -54,
	38, 36, -34, -38, 80, -13, -13, 48, 47, -51,
	54, 50, 55, -13, -13, -13, -13, -13, -13, -13,
	-21, -50, 80, -37, 64, -21, -36, -38, -10, 29,
	-13, -32, -13, -24, 36, -54, 36, -24, -15, -15,
	-42, -44, -13, 79, -15, -13, 80, 76, -32, -31,
	-37, 66, -13, -10, -11, 35, 32, 81, -26, 43,
	44, -24, 36, -24, 51, -32, 67, -13, 65, 80,
	69, -40, 75, -33, -14, -13, -32, -13, 79, -24,
	-26, -13, 80, -13, 5, 81, -52, 33, 34, -35,
	-38, 80, -33, 80, 81, -35,
}
var yyDef = [...]int{

	1, -2, 1, 135, 4, 5, 6, 9, 0, 0,
	123, 116, 2, 3, 136, 11, 0, 120, 121, 118,
	0, 0, 124, 13, 0, 10, -2, 73, 76, 77,
	70, 71, 72, 107, 0, 0, 0, 117, 27, 28,
	29, 30, 31, 32, 33, 34, 35, 36, 37, 0,
	21, 22, 23, 24, 25, 26, 64, 91, 0, 0,
	108, 109, 0, 111, 113, 114, 115, 0, 8, 103,
	88, -2, 87, 15, 0, -2, 0, -2, 129, 0,
	0, 128, 132, 133, 134, 74, 0, 0, 122, -2,
	0, 0, 125, 0, 0, 0, 0, 0, 0, 0,
	0, 126, -2, 125, 0, -2, 54, 62, 110, 112,
	-2, 0, 0, 17, 0, 0, 0, 0, 129, 0,
	130, 0, 100, 75, 43, 44, 45, 0, 0, 0,
	0, 0, 0, 55, 56, 57, 58, 59, 60, 61,
	0, 0, 38, 93, 0, 17, 104, 89, 19, 0,
	-2, 14, -2, 83, 0, 0, 0, 82, 52, 53,
	46, 47, 125, 0, 50, 51, 63, 66, 67, 0,
	106, 0, 125, 0, 7, 0, 0, 0, 78, 0,
	0, 79, 0, 83, 0, 0, 90, -2, 0, 68,
	0, 20, 0, 18, 97, -2, 96, -2, 0, 81,
	80, 48, 49, -2, 0, 0, 39, 41, 42, 0,
	101, 69, 98, 85, 0, 102,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 78, 3, 3,
	79, 80, 76, 74, 81, 75, 3, 77, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 82,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72, 73,
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
		//line parser.y:104
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:109
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:116
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:122
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:126
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:130
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:136
		{
			yyVAL.expression = SelectQuery{
				SelectClause:  yyDollar[1].expression,
				FromClause:    yyDollar[2].expression,
				WhereClause:   yyDollar[3].expression,
				GroupByClause: yyDollar[4].expression,
				HavingClause:  yyDollar[5].expression,
				OrderByClause: yyDollar[6].expression,
				LimitClause:   yyDollar[7].expression,
			}
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:150
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 9:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:156
		{
			yyVAL.expression = nil
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:160
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 11:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:166
		{
			yyVAL.expression = nil
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:170
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 13:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:176
		{
			yyVAL.expression = nil
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:180
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:186
		{
			yyVAL.expression = nil
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:190
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 17:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:196
		{
			yyVAL.expression = nil
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:200
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:206
		{
			yyVAL.expression = nil
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:210
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:216
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:220
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:224
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:228
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:232
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:236
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:242
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:246
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:250
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:254
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:258
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:262
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:266
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:270
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:274
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:278
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:282
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:286
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:292
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 40:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:298
		{
			yyVAL.token = Token{}
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:302
		{
			yyVAL.token = yyDollar[1].token
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:306
		{
			yyVAL.token = yyDollar[1].token
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:312
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:318
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
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 46:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:349
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 48:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:353
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 49:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:357
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:361
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:365
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:369
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:373
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:377
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:383
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:387
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:391
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:399
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:405
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:409
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:413
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 63:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:419
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:423
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:429
		{
			yyVAL.expression = Option{}
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:433
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:437
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 68:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:443
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 69:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:447
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:453
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:457
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:461
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:467
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:471
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:475
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:479
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:483
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 78:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:489
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 79:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:493
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 80:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:497
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 81:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:501
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 82:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:505
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 83:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:511
		{
			yyVAL.expression = nil
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:515
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 85:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:519
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:525
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:529
		{
			yyVAL.expression = AllColumns{}
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:535
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:539
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:545
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 91:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:551
		{
			yyVAL.expression = nil
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:555
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:561
		{
			yyVAL.expression = nil
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:565
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:571
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:575
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:581
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:585
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:591
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:595
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:601
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:605
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:611
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:615
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:621
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:625
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:631
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:637
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:643
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:647
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:654
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:658
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:665
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:671
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:677
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:683
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:689
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:695
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:699
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 120:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:705
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:711
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:715
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 123:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:721
		{
			yyVAL.token = Token{}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:725
		{
			yyVAL.token = yyDollar[1].token
		}
	case 125:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:731
		{
			yyVAL.token = Token{}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:735
		{
			yyVAL.token = yyDollar[1].token
		}
	case 127:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:741
		{
			yyVAL.token = Token{}
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:745
		{
			yyVAL.token = yyDollar[1].token
		}
	case 129:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:751
		{
			yyVAL.token = Token{}
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:755
		{
			yyVAL.token = yyDollar[1].token
		}
	case 131:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:761
		{
			yyVAL.token = Token{}
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:765
		{
			yyVAL.token = yyDollar[1].token
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:769
		{
			yyVAL.token = yyDollar[1].token
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:773
		{
			yyVAL.token = yyDollar[1].token
		}
	case 135:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:779
		{
			yyVAL.token = Token{}
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:783
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
