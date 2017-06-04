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
	"'='",
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

//line parser.y:798

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
	-1, 104,
	81, 65,
	-2, 123,
	-1, 107,
	64, 92,
	-2, 125,
	-1, 112,
	29, 65,
	69, 65,
	81, 65,
	-2, 123,
	-1, 152,
	50, 125,
	54, 125,
	55, 125,
	-2, 16,
	-1, 154,
	50, 125,
	54, 125,
	55, 125,
	-2, 95,
	-1, 189,
	67, 94,
	-2, 125,
	-1, 197,
	50, 125,
	54, 125,
	55, 125,
	-2, 40,
	-1, 199,
	50, 125,
	54, 125,
	55, 125,
	-2, 84,
	-1, 205,
	64, 105,
	66, 105,
	67, 105,
	-2, 125,
}

const yyPrivate = 57344

const yyLast = 456

var yyAct = [...]int{

	154, 38, 211, 195, 153, 51, 150, 26, 180, 55,
	53, 121, 172, 68, 14, 142, 216, 207, 30, 113,
	93, 37, 71, 192, 35, 75, 215, 33, 213, 85,
	204, 42, 168, 126, 20, 191, 81, 89, 82, 83,
	84, 79, 92, 25, 77, 29, 32, 200, 31, 165,
	105, 34, 100, 99, 103, 112, 104, 92, 107, 61,
	109, 94, 95, 96, 97, 98, 78, 76, 110, 61,
	63, 36, 101, 90, 188, 102, 94, 95, 96, 97,
	98, 80, 30, 179, 11, 10, 146, 22, 125, 52,
	108, 127, 128, 92, 103, 135, 136, 137, 138, 139,
	140, 141, 122, 34, 33, 60, 61, 63, 64, 64,
	65, 11, 31, 131, 71, 92, 149, 152, 184, 145,
	30, 146, 158, 173, 124, 30, 155, 148, 147, 194,
	156, 159, 157, 164, 123, 167, 96, 97, 98, 62,
	118, 163, 162, 8, 120, 119, 177, 174, 170, 58,
	31, 178, 111, 59, 175, 31, 66, 66, 30, 117,
	30, 160, 161, 57, 183, 166, 185, 133, 67, 116,
	187, 132, 134, 17, 189, 74, 62, 169, 151, 197,
	49, 114, 199, 193, 198, 24, 30, 203, 31, 16,
	31, 205, 201, 81, 202, 82, 83, 84, 79, 10,
	110, 77, 212, 81, 19, 82, 83, 84, 197, 88,
	206, 214, 33, 143, 48, 6, 31, 6, 212, 217,
	33, 60, 61, 63, 21, 64, 65, 11, 13, 33,
	60, 61, 63, 33, 64, 65, 11, 10, 33, 60,
	61, 63, 91, 64, 65, 11, 47, 9, 1, 9,
	86, 12, 208, 5, 33, 60, 61, 63, 18, 64,
	65, 11, 54, 130, 129, 58, 50, 87, 4, 59,
	4, 171, 106, 66, 58, 44, 69, 70, 59, 57,
	28, 27, 66, 58, 67, 56, 43, 59, 57, 46,
	40, 66, 62, 67, 45, 41, 49, 57, 196, 58,
	39, 62, 67, 59, 176, 49, 115, 66, 73, 23,
	62, 72, 15, 57, 49, 100, 99, 103, 67, 7,
	92, 3, 2, 0, 0, 0, 62, 0, 209, 210,
	49, 0, 100, 99, 103, 101, 90, 92, 102, 94,
	95, 96, 97, 98, 0, 144, 100, 99, 103, 0,
	0, 92, 101, 90, 0, 102, 94, 95, 96, 97,
	98, 0, 0, 100, 99, 103, 101, 90, 92, 102,
	94, 95, 96, 97, 98, 0, 0, 190, 0, 0,
	186, 99, 103, 101, 90, 92, 102, 94, 95, 96,
	97, 98, 100, 0, 103, 0, 0, 92, 0, 0,
	101, 90, 0, 102, 94, 95, 96, 97, 98, 0,
	0, 103, 101, 90, 92, 102, 94, 95, 96, 97,
	98, 0, 0, 92, 0, 0, 0, 0, 0, 101,
	90, 0, 102, 94, 95, 96, 97, 98, 101, 90,
	0, 102, 94, 95, 96, 97, 98, 81, 0, 82,
	83, 84, 79, 181, 182, 77,
}
var yyPact = [...]int{

	73, -1000, 73, -69, -1000, -1000, -1000, 176, 193, -39,
	29, -1000, -1000, -1000, -1000, 168, 23, -1000, -58, -2,
	250, 234, -1000, 145, 250, -1000, -1, 229, -1000, -1000,
	-1000, -1000, -1000, -1000, 187, 193, 250, 281, -24, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -39, -1000, 225,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 250, -29, 250,
	-1000, -1000, 63, -1000, -1000, -1000, -1000, -25, -1000, -63,
	160, 281, -1000, 138, 127, 281, 104, 166, 64, 98,
	23, -1000, -1000, -1000, -1000, -1000, 208, -48, -1000, 281,
	250, 216, 41, 117, 250, 250, 250, 250, 250, 250,
	250, -1000, -1000, -1000, 29, 264, 22, 281, -1000, 367,
	-1000, -1000, 29, 234, 208, 149, 250, 250, 23, 94,
	64, 86, -1000, 23, -1000, -1000, -1000, -14, 281, -29,
	-29, 99, 250, -31, 250, 59, 59, 37, 37, 37,
	341, 358, -49, 100, -1000, 57, 250, 149, -1000, -1000,
	111, 119, 281, -1000, 1, 410, 23, 82, 23, 156,
	-1000, -1000, -1000, -1000, 329, 225, -1000, 281, -1000, -1000,
	-1000, 7, 22, 250, 312, -46, -1000, 53, 250, 250,
	-1000, 250, -33, 156, 23, 410, 250, -51, -1000, 281,
	250, -1000, 205, -1000, 53, -1000, -65, 295, -1000, 281,
	208, 156, -1000, 358, -1000, 281, -53, 250, -1000, -1000,
	-1000, -55, -66, -1000, -1000, -1000, 208, -1000,
}
var yyPgo = [...]int{

	0, 248, 322, 321, 267, 319, 312, 309, 308, 306,
	6, 304, 300, 0, 298, 31, 295, 294, 290, 289,
	286, 15, 285, 281, 7, 280, 8, 277, 276, 275,
	272, 271, 4, 3, 43, 2, 13, 12, 1, 266,
	5, 89, 10, 262, 9, 246, 214, 258, 253, 173,
	213, 20, 252, 67, 11, 66, 242, 228,
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
	54, 55, 55, 55, 55, 56, 56, 57, 57,
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
	1, 0, 1, 1, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -48, -46, -5, 70, -45,
	12, 11, -1, -57, 83, -6, 13, -49, -47, 11,
	73, -50, 58, -7, 17, -34, -24, -23, -25, 22,
	-38, -15, 23, 4, 80, 82, 73, -13, -38, -12,
	-18, -16, -15, -20, -29, -17, -19, -45, -46, 80,
	-39, -40, -41, -42, -43, -44, -22, 63, 49, 53,
	5, 6, 76, 7, 9, 10, 57, 68, -36, -28,
	-27, -13, 77, -8, 30, -13, -53, 45, -55, 42,
	82, 37, 39, 40, 41, -38, 21, -4, -49, -13,
	72, -56, 56, -51, 75, 76, 77, 78, 79, 52,
	51, 71, 74, 53, 80, -13, -30, -13, -15, -13,
	-40, -41, 80, 82, 21, -9, 31, 32, 36, -53,
	-55, -54, 38, 36, -34, -38, 81, -13, -13, 48,
	47, -51, 54, 50, 55, -13, -13, -13, -13, -13,
	-13, -13, -21, -50, 81, -37, 64, -21, -36, -38,
	-10, 29, -13, -32, -13, -24, 36, -54, 36, -24,
	-15, -15, -42, -44, -13, 80, -15, -13, 81, 77,
	-32, -31, -37, 66, -13, -10, -11, 35, 32, 82,
	-26, 43, 44, -24, 36, -24, 51, -32, 67, -13,
	65, 81, 69, -40, 76, -33, -14, -13, -32, -13,
	80, -24, -26, -13, 81, -13, 5, 82, -52, 33,
	34, -35, -38, 81, -33, 81, 82, -35,
}
var yyDef = [...]int{

	1, -2, 1, 137, 4, 5, 6, 9, 0, 0,
	123, 116, 2, 3, 138, 11, 0, 120, 121, 118,
	0, 0, 124, 13, 0, 10, -2, 73, 76, 77,
	70, 71, 72, 107, 0, 0, 0, 117, 27, 28,
	29, 30, 31, 32, 33, 34, 35, 36, 37, 0,
	21, 22, 23, 24, 25, 26, 64, 91, 0, 0,
	108, 109, 0, 111, 113, 114, 115, 0, 8, 103,
	88, -2, 87, 15, 0, -2, 0, -2, 129, 0,
	0, 128, 132, 133, 134, 74, 0, 0, 122, -2,
	0, 0, 125, 0, 0, 0, 0, 0, 0, 0,
	0, 135, 136, 126, -2, 125, 0, -2, 54, 62,
	110, 112, -2, 0, 0, 17, 0, 0, 0, 0,
	129, 0, 130, 0, 100, 75, 43, 44, 45, 0,
	0, 0, 0, 0, 0, 55, 56, 57, 58, 59,
	60, 61, 0, 0, 38, 93, 0, 17, 104, 89,
	19, 0, -2, 14, -2, 83, 0, 0, 0, 82,
	52, 53, 46, 47, 125, 0, 50, 51, 63, 66,
	67, 0, 106, 0, 125, 0, 7, 0, 0, 0,
	78, 0, 0, 79, 0, 83, 0, 0, 90, -2,
	0, 68, 0, 20, 0, 18, 97, -2, 96, -2,
	0, 81, 80, 48, 49, -2, 0, 0, 39, 41,
	42, 0, 101, 69, 98, 85, 0, 102,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 79, 3, 3,
	80, 81, 77, 75, 82, 76, 3, 78, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 83,
	3, 74,
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
		//line parser.y:105
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:110
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:117
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:123
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:127
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:131
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:137
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
		//line parser.y:151
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 9:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:157
		{
			yyVAL.expression = nil
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:161
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 11:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:167
		{
			yyVAL.expression = nil
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:171
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 13:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:177
		{
			yyVAL.expression = nil
		}
	case 14:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:181
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:187
		{
			yyVAL.expression = nil
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:191
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 17:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:197
		{
			yyVAL.expression = nil
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:201
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 19:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:207
		{
			yyVAL.expression = nil
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:211
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:217
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:221
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:225
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:229
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:233
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:237
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:243
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:247
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:251
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:255
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:259
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:263
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:267
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:271
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:275
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:279
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:283
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:287
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:293
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 40:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:299
		{
			yyVAL.token = Token{}
		}
	case 41:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:303
		{
			yyVAL.token = yyDollar[1].token
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:307
		{
			yyVAL.token = yyDollar[1].token
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:313
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:319
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
		//line parser.y:342
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 46:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:346
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:350
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 48:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:354
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 49:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:358
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:362
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:366
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:370
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:374
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:378
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:384
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:388
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:392
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:396
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:400
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:406
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:410
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 62:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:414
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 63:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:420
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:424
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:430
		{
			yyVAL.expression = Option{}
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:434
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:438
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 68:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:444
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 69:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:448
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 70:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:454
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:458
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:462
		{
			yyVAL.expression = Stdin{Stdin: yyDollar[1].token.Literal}
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:468
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:472
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:476
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:480
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:484
		{
			yyVAL.expression = Table{Object: Dual{Dual: yyDollar[1].token.Literal}}
		}
	case 78:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:490
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 79:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:494
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 80:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:498
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 81:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:502
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 82:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:506
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 83:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:512
		{
			yyVAL.expression = nil
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 85:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:520
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:526
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:530
		{
			yyVAL.expression = AllColumns{}
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:540
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:546
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 91:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = nil
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:556
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:562
		{
			yyVAL.expression = nil
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:566
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:576
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:582
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:586
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:596
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:602
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:606
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:616
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:622
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:626
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:632
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:638
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:644
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:648
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:655
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:659
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:666
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:672
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:678
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:684
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:690
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:696
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:700
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 120:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:706
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:712
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:716
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 123:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:722
		{
			yyVAL.token = Token{}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:726
		{
			yyVAL.token = yyDollar[1].token
		}
	case 125:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:732
		{
			yyVAL.token = Token{}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:736
		{
			yyVAL.token = yyDollar[1].token
		}
	case 127:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:742
		{
			yyVAL.token = Token{}
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:746
		{
			yyVAL.token = yyDollar[1].token
		}
	case 129:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:752
		{
			yyVAL.token = Token{}
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:756
		{
			yyVAL.token = yyDollar[1].token
		}
	case 131:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:762
		{
			yyVAL.token = Token{}
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:766
		{
			yyVAL.token = yyDollar[1].token
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:770
		{
			yyVAL.token = yyDollar[1].token
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:774
		{
			yyVAL.token = yyDollar[1].token
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:780
		{
			yyVAL.token = yyDollar[1].token
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:784
		{
			yyVAL.token = Token{Token: COMPARISON_OP, Literal: string('=')}
		}
	case 137:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:790
		{
			yyVAL.token = Token{}
		}
	case 138:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:794
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
