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
const CREATE = 57365
const DROP = 57366
const ALTER = 57367
const TABLE = 57368
const COLUMN = 57369
const ORDER = 57370
const GROUP = 57371
const HAVING = 57372
const BY = 57373
const ASC = 57374
const DESC = 57375
const LIMIT = 57376
const JOIN = 57377
const INNER = 57378
const OUTER = 57379
const LEFT = 57380
const RIGHT = 57381
const FULL = 57382
const CROSS = 57383
const ON = 57384
const USING = 57385
const NATURAL = 57386
const UNION = 57387
const ALL = 57388
const ANY = 57389
const EXISTS = 57390
const IN = 57391
const AND = 57392
const OR = 57393
const NOT = 57394
const BETWEEN = 57395
const LIKE = 57396
const IS = 57397
const NULL = 57398
const DISTINCT = 57399
const WITH = 57400
const TRUE = 57401
const FALSE = 57402
const UNKNOWN = 57403
const CASE = 57404
const WHEN = 57405
const THEN = 57406
const ELSE = 57407
const END = 57408
const GROUP_CONCAT = 57409
const SEPARATOR = 57410
const VAR = 57411
const COMPARISON_OP = 57412
const STRING_OP = 57413
const SUBSTITUTION_OP = 57414

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

//line parser.y:784

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
	-1, 27,
	35, 127,
	37, 131,
	-2, 99,
	-1, 69,
	49, 125,
	53, 125,
	54, 125,
	-2, 86,
	-1, 73,
	49, 125,
	53, 125,
	54, 125,
	-2, 13,
	-1, 75,
	37, 131,
	-2, 127,
	-1, 89,
	49, 125,
	53, 125,
	54, 125,
	-2, 119,
	-1, 102,
	79, 66,
	-2, 123,
	-1, 105,
	63, 92,
	-2, 125,
	-1, 110,
	28, 66,
	68, 66,
	79, 66,
	-2, 123,
	-1, 151,
	49, 125,
	53, 125,
	54, 125,
	-2, 17,
	-1, 153,
	49, 125,
	53, 125,
	54, 125,
	-2, 95,
	-1, 188,
	66, 94,
	-2, 125,
	-1, 196,
	49, 125,
	53, 125,
	54, 125,
	-2, 41,
	-1, 198,
	49, 125,
	53, 125,
	54, 125,
	-2, 84,
	-1, 204,
	63, 105,
	65, 105,
	66, 105,
	-2, 125,
}

const yyPrivate = 57344

const yyLast = 485

var yyAct = [...]int{

	153, 36, 210, 194, 27, 49, 179, 149, 53, 141,
	152, 40, 51, 66, 171, 93, 119, 14, 28, 215,
	206, 35, 69, 92, 111, 73, 33, 214, 29, 191,
	83, 85, 212, 203, 79, 89, 80, 81, 82, 77,
	190, 92, 75, 96, 97, 98, 167, 199, 103, 100,
	99, 101, 125, 31, 92, 31, 105, 26, 107, 94,
	95, 96, 97, 98, 20, 164, 108, 32, 106, 91,
	90, 25, 94, 95, 96, 97, 98, 59, 78, 178,
	28, 110, 102, 59, 61, 76, 123, 74, 124, 34,
	29, 126, 127, 187, 145, 134, 135, 136, 137, 138,
	139, 140, 145, 62, 172, 50, 11, 10, 130, 22,
	92, 101, 69, 120, 148, 151, 183, 157, 28, 144,
	146, 154, 155, 28, 121, 147, 158, 32, 29, 32,
	177, 116, 163, 29, 166, 156, 122, 176, 115, 162,
	159, 160, 114, 161, 165, 193, 173, 72, 150, 17,
	64, 60, 112, 169, 174, 16, 24, 28, 10, 28,
	182, 118, 184, 117, 8, 19, 109, 29, 142, 29,
	31, 31, 132, 188, 205, 186, 131, 133, 196, 21,
	13, 198, 192, 88, 207, 28, 202, 86, 200, 197,
	204, 201, 5, 92, 79, 29, 80, 81, 82, 108,
	18, 211, 46, 6, 52, 6, 31, 196, 91, 90,
	213, 94, 95, 96, 97, 98, 48, 211, 216, 31,
	58, 59, 61, 84, 62, 63, 11, 170, 31, 58,
	59, 61, 104, 62, 63, 11, 10, 79, 42, 80,
	81, 82, 77, 180, 181, 75, 67, 68, 31, 58,
	59, 61, 30, 62, 63, 11, 45, 9, 54, 9,
	41, 129, 128, 56, 44, 87, 4, 57, 4, 38,
	1, 64, 56, 12, 43, 39, 57, 55, 195, 37,
	64, 175, 65, 113, 71, 23, 55, 15, 7, 60,
	3, 65, 56, 47, 2, 0, 57, 0, 60, 0,
	64, 0, 47, 0, 0, 0, 55, 31, 58, 59,
	61, 65, 62, 63, 11, 100, 99, 101, 60, 168,
	92, 0, 47, 0, 31, 58, 59, 61, 0, 62,
	63, 11, 0, 0, 0, 91, 90, 0, 94, 95,
	96, 97, 98, 0, 143, 79, 0, 80, 81, 82,
	77, 56, 0, 75, 0, 57, 0, 0, 0, 64,
	0, 0, 0, 0, 0, 55, 208, 209, 56, 0,
	65, 0, 57, 0, 0, 0, 64, 60, 70, 0,
	0, 47, 55, 0, 100, 99, 101, 65, 0, 92,
	0, 0, 0, 0, 60, 0, 0, 0, 47, 0,
	0, 100, 99, 101, 91, 90, 92, 94, 95, 96,
	97, 98, 100, 99, 101, 189, 0, 92, 0, 0,
	0, 91, 90, 0, 94, 95, 96, 97, 98, 185,
	99, 101, 91, 90, 92, 94, 95, 96, 97, 98,
	100, 0, 101, 0, 0, 92, 0, 0, 0, 91,
	90, 0, 94, 95, 96, 97, 98, 0, 0, 101,
	91, 90, 92, 94, 95, 96, 97, 98, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 91, 90, 0,
	94, 95, 96, 97, 98,
}
var yyPact = [...]int{

	95, -1000, 95, -64, -1000, -1000, -1000, 142, 154, -8,
	52, -1000, -1000, -1000, -1000, 139, 49, -1000, -54, 17,
	320, 303, -1000, 118, 320, -1000, -1000, -2, 202, 166,
	-1000, -1000, 146, 154, 320, 362, 4, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -8, -1000, 224, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 320, -11, 320, -1000, -1000,
	77, -1000, -1000, -1000, -1000, 3, -1000, -56, 131, 362,
	-1000, 112, 107, 362, 96, 158, 76, 89, 51, -1000,
	-1000, -1000, -1000, -1000, 167, -1000, 167, -27, -1000, 362,
	320, 215, 59, 123, 320, 320, 320, 320, 320, 320,
	320, -1000, 52, 265, 31, 362, -1000, 138, -1000, -1000,
	52, 303, 167, 120, 320, 320, 51, 87, 76, 82,
	-1000, 51, -1000, -1000, -1000, -1000, -14, -14, -11, -11,
	94, 320, -13, 320, -32, -32, 55, 55, 55, 390,
	407, -33, 244, -1000, 39, 320, 120, -1000, -1000, 103,
	99, 362, -1000, -1, 201, 51, 81, 51, 309, -1000,
	-1000, -1000, -1000, 379, 224, -1000, 362, -1000, -1000, -1000,
	27, 31, 320, 351, -39, -1000, 71, 320, 320, -1000,
	320, -31, 309, 51, 201, 320, -46, -1000, 362, 320,
	-1000, 169, -1000, 71, -1000, -60, 334, -1000, 362, 167,
	309, -1000, 407, -1000, 362, -47, 320, -1000, -1000, -1000,
	-52, -61, -1000, -1000, -1000, 167, -1000,
}
var yyPgo = [...]int{

	0, 270, 294, 290, 265, 288, 287, 285, 284, 283,
	7, 281, 279, 0, 278, 11, 275, 274, 269, 264,
	260, 9, 258, 4, 252, 6, 247, 246, 238, 232,
	227, 10, 3, 57, 2, 13, 14, 1, 216, 5,
	105, 12, 204, 8, 256, 202, 200, 192, 149, 168,
	15, 184, 87, 16, 85, 180,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 3, 3, 4, 5, 6,
	6, 6, 7, 7, 8, 8, 9, 9, 10, 10,
	11, 11, 12, 12, 12, 12, 12, 12, 13, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	14, 51, 51, 51, 15, 16, 17, 17, 17, 17,
	17, 17, 17, 17, 17, 17, 18, 18, 18, 18,
	18, 19, 19, 19, 20, 20, 21, 21, 21, 22,
	22, 23, 23, 23, 23, 23, 23, 23, 24, 24,
	24, 24, 24, 25, 25, 25, 26, 26, 27, 27,
	28, 29, 29, 30, 30, 31, 31, 32, 32, 33,
	33, 34, 34, 35, 35, 36, 36, 37, 38, 39,
	39, 40, 40, 41, 42, 43, 44, 45, 46, 46,
	47, 48, 48, 49, 49, 50, 50, 52, 52, 53,
	53, 54, 54, 54, 54, 55, 55,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 1, 1, 1, 7, 3, 0,
	2, 2, 0, 2, 0, 3, 0, 2, 0, 3,
	0, 2, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	2, 0, 1, 1, 3, 3, 3, 4, 4, 6,
	6, 4, 4, 4, 4, 2, 3, 3, 3, 3,
	3, 3, 3, 2, 4, 1, 0, 2, 2, 5,
	7, 1, 2, 3, 1, 2, 3, 1, 5, 5,
	6, 6, 4, 0, 2, 4, 1, 1, 1, 3,
	5, 0, 1, 0, 2, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 4, 2, 1, 1, 1,
	2, 1, 2, 1, 1, 1, 1, 3, 1, 3,
	2, 1, 3, 0, 1, 0, 1, 0, 1, 0,
	1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -47, -45, -5, 69, -44,
	12, 11, -1, -55, 81, -6, 13, -48, -46, 11,
	72, -49, 57, -7, 17, 22, -33, -23, -37, -15,
	-24, 4, 78, 80, 72, -13, -37, -12, -18, -16,
	-15, -20, -28, -17, -19, -44, -45, 78, -38, -39,
	-40, -41, -42, -43, -22, 62, 48, 52, 5, 6,
	74, 7, 9, 10, 56, 67, -35, -27, -26, -13,
	75, -8, 29, -13, -52, 44, -54, 41, 80, 36,
	38, 39, 40, -37, 21, -37, 21, -4, -48, -13,
	71, 70, 55, -50, 73, 74, 75, 76, 77, 51,
	50, 52, 78, -13, -29, -13, -15, -13, -39, -40,
	78, 80, 21, -9, 30, 31, 35, -52, -54, -53,
	37, 35, -33, -37, -37, 79, -13, -13, 47, 46,
	-50, 53, 49, 54, -13, -13, -13, -13, -13, -13,
	-13, -21, -49, 79, -36, 63, -21, -35, -37, -10,
	28, -13, -31, -13, -23, 35, -53, 35, -23, -15,
	-15, -41, -43, -13, 78, -15, -13, 79, 75, -31,
	-30, -36, 65, -13, -10, -11, 34, 31, 80, -25,
	42, 43, -23, 35, -23, 50, -31, 66, -13, 64,
	79, 68, -39, 74, -32, -14, -13, -31, -13, 78,
	-23, -25, -13, 79, -13, 5, 80, -51, 32, 33,
	-34, -37, 79, -32, 79, 80, -34,
}
var yyDef = [...]int{

	1, -2, 1, 135, 4, 5, 6, 9, 0, 0,
	123, 116, 2, 3, 136, 12, 0, 120, 121, 118,
	0, 0, 124, 14, 0, 10, 11, -2, 71, 74,
	77, 107, 0, 0, 0, 117, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 0, 22, 23,
	24, 25, 26, 27, 65, 91, 0, 0, 108, 109,
	0, 111, 113, 114, 115, 0, 8, 103, 88, -2,
	87, 16, 0, -2, 0, -2, 129, 0, 0, 128,
	132, 133, 134, 72, 0, 75, 0, 0, 122, -2,
	0, 0, 125, 0, 0, 0, 0, 0, 0, 0,
	0, 126, -2, 125, 0, -2, 55, 63, 110, 112,
	-2, 0, 0, 18, 0, 0, 0, 0, 129, 0,
	130, 0, 100, 73, 76, 44, 45, 46, 0, 0,
	0, 0, 0, 0, 56, 57, 58, 59, 60, 61,
	62, 0, 0, 39, 93, 0, 18, 104, 89, 20,
	0, -2, 15, -2, 83, 0, 0, 0, 82, 53,
	54, 47, 48, 125, 0, 51, 52, 64, 67, 68,
	0, 106, 0, 125, 0, 7, 0, 0, 0, 78,
	0, 0, 79, 0, 83, 0, 0, 90, -2, 0,
	69, 0, 21, 0, 19, 97, -2, 96, -2, 0,
	81, 80, 49, 50, -2, 0, 0, 40, 42, 43,
	0, 101, 70, 98, 85, 0, 102,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 77, 3, 3,
	78, 79, 75, 73, 80, 74, 3, 76, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 81,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71,
	72,
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
		//line parser.y:103
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:108
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:115
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:121
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:125
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:129
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 7:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:135
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
		//line parser.y:149
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 9:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:155
		{
			yyVAL.expression = nil
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:159
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: []Expression{Dual{Dual: yyDollar[2].token.Literal}}}
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:163
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:169
		{
			yyVAL.expression = nil
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:173
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:179
		{
			yyVAL.expression = nil
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:183
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 16:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:189
		{
			yyVAL.expression = nil
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:193
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 18:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:199
		{
			yyVAL.expression = nil
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:203
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:209
		{
			yyVAL.expression = nil
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:213
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:219
		{
			yyVAL.primary = yyDollar[1].text
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:223
		{
			yyVAL.primary = yyDollar[1].integer
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:227
		{
			yyVAL.primary = yyDollar[1].float
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:231
		{
			yyVAL.primary = yyDollar[1].ternary
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:235
		{
			yyVAL.primary = yyDollar[1].datetime
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:239
		{
			yyVAL.primary = yyDollar[1].null
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:245
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:249
		{
			yyVAL.expression = yyDollar[1].primary
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:253
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:257
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:261
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 33:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:265
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:269
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:273
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:277
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:281
		{
			yyVAL.expression = yyDollar[1].variable
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:285
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:289
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 40:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:295
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 41:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:301
		{
			yyVAL.token = Token{}
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:305
		{
			yyVAL.token = yyDollar[1].token
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:309
		{
			yyVAL.token = yyDollar[1].token
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:315
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:321
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
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:344
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:348
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 48:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:352
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 49:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:356
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 50:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:360
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:364
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:368
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:372
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 54:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:376
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:380
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:386
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:390
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 58:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:394
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:398
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:402
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:408
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:412
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 63:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:416
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 64:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:422
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:426
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 66:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:432
		{
			yyVAL.expression = Option{}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:436
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:440
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 69:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:446
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 70:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:450
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].token.Literal}
		}
	case 71:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:456
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:460
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:464
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:468
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:472
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:476
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:480
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 78:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:486
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 79:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:490
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 80:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:494
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 81:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:498
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 82:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:502
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 83:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:508
		{
			yyVAL.expression = nil
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:512
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 85:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:516
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:522
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:526
		{
			yyVAL.expression = AllColumns{}
		}
	case 88:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:532
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:536
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:542
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 91:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:548
		{
			yyVAL.expression = nil
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:552
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 93:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:558
		{
			yyVAL.expression = nil
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:562
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:568
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:572
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:578
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:582
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:588
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:592
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:598
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:602
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:608
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:612
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:618
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:622
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:628
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:634
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:640
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:644
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:651
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:655
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:662
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:668
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:674
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:680
		{
			yyVAL.variable = Variable{Name: yyDollar[1].token.Literal}
		}
	case 117:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:686
		{
			yyVAL.expression = VariableSubstitution{Variable: yyDollar[1].variable, Value: yyDollar[3].expression}
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:692
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal}
		}
	case 119:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:696
		{
			yyVAL.expression = VariableAssignment{Name: yyDollar[1].token.Literal, Value: yyDollar[3].expression}
		}
	case 120:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:702
		{
			yyVAL.expression = VariableDeclaration{Var: yyDollar[1].token.Literal, Assignments: yyDollar[2].expressions}
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:708
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:712
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 123:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:718
		{
			yyVAL.token = Token{}
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:722
		{
			yyVAL.token = yyDollar[1].token
		}
	case 125:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:728
		{
			yyVAL.token = Token{}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:732
		{
			yyVAL.token = yyDollar[1].token
		}
	case 127:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:738
		{
			yyVAL.token = Token{}
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:742
		{
			yyVAL.token = yyDollar[1].token
		}
	case 129:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:748
		{
			yyVAL.token = Token{}
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:752
		{
			yyVAL.token = yyDollar[1].token
		}
	case 131:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:758
		{
			yyVAL.token = Token{}
		}
	case 132:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:762
		{
			yyVAL.token = yyDollar[1].token
		}
	case 133:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:766
		{
			yyVAL.token = yyDollar[1].token
		}
	case 134:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:770
		{
			yyVAL.token = yyDollar[1].token
		}
	case 135:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:776
		{
			yyVAL.token = Token{}
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:780
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
