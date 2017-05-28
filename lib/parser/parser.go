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
	identifier  Identifier
	text        String
	integer     Integer
	float       Float
	ternary     Ternary
	datetime    Datetime
	null        Null
	token       Token
}

const IDENTIFIER = 57346
const STRING = 57347
const INTEGER = 57348
const FLOAT = 57349
const BOOLEAN = 57350
const TERNARY = 57351
const DATETIME = 57352
const SELECT = 57353
const FROM = 57354
const UPDATE = 57355
const SET = 57356
const DELETE = 57357
const WHERE = 57358
const INSERT = 57359
const INTO = 57360
const VALUES = 57361
const AS = 57362
const DUAL = 57363
const CREATE = 57364
const DROP = 57365
const ALTER = 57366
const TABLE = 57367
const COLUMN = 57368
const ORDER = 57369
const GROUP = 57370
const HAVING = 57371
const BY = 57372
const ASC = 57373
const DESC = 57374
const LIMIT = 57375
const JOIN = 57376
const INNER = 57377
const OUTER = 57378
const LEFT = 57379
const RIGHT = 57380
const FULL = 57381
const CROSS = 57382
const ON = 57383
const USING = 57384
const NATURAL = 57385
const UNION = 57386
const ALL = 57387
const ANY = 57388
const EXISTS = 57389
const IN = 57390
const AND = 57391
const OR = 57392
const NOT = 57393
const BETWEEN = 57394
const LIKE = 57395
const IS = 57396
const NULL = 57397
const DISTINCT = 57398
const WITH = 57399
const TRUE = 57400
const FALSE = 57401
const UNKNOWN = 57402
const CASE = 57403
const WHEN = 57404
const THEN = 57405
const ELSE = 57406
const END = 57407
const GROUP_CONCAT = 57408
const SEPARATOR = 57409
const COMPARISON_OP = 57410
const STRING_OP = 57411

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
	"COMPARISON_OP",
	"STRING_OP",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'('",
	"')'",
	"'%'",
	"','",
	"';'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.y:743

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
	-1, 18,
	34, 120,
	36, 124,
	-2, 99,
	-1, 29,
	48, 118,
	52, 118,
	53, 118,
	-2, 58,
	-1, 62,
	36, 124,
	-2, 120,
	-1, 89,
	49, 58,
	50, 58,
	-2, 118,
	-1, 93,
	75, 64,
	-2, 116,
	-1, 99,
	27, 64,
	67, 64,
	75, 64,
	-2, 116,
}

const yyPrivate = 57344

const yyLast = 394

var yyAct = [...]int{

	41, 202, 29, 185, 158, 43, 27, 18, 138, 42,
	169, 161, 19, 36, 47, 106, 45, 9, 132, 141,
	70, 72, 60, 207, 88, 20, 198, 81, 82, 79,
	83, 84, 85, 86, 75, 89, 87, 168, 87, 90,
	206, 80, 79, 83, 84, 85, 86, 92, 130, 87,
	78, 77, 95, 91, 83, 84, 85, 86, 97, 24,
	87, 204, 22, 50, 51, 53, 19, 54, 55, 6,
	78, 77, 110, 182, 111, 88, 131, 114, 81, 20,
	194, 181, 117, 118, 115, 116, 125, 126, 127, 128,
	129, 22, 80, 79, 83, 84, 85, 86, 177, 136,
	87, 156, 112, 190, 19, 142, 134, 153, 140, 19,
	121, 143, 23, 56, 99, 93, 147, 20, 137, 49,
	17, 145, 20, 54, 57, 152, 178, 155, 135, 52,
	78, 77, 96, 148, 149, 113, 151, 154, 150, 13,
	159, 135, 163, 162, 180, 19, 164, 19, 88, 85,
	86, 22, 172, 87, 174, 51, 142, 78, 20, 63,
	20, 23, 51, 53, 107, 179, 44, 173, 16, 56,
	187, 142, 183, 176, 19, 123, 61, 189, 193, 122,
	124, 191, 195, 196, 159, 192, 109, 20, 188, 146,
	97, 203, 197, 22, 50, 51, 53, 144, 54, 55,
	108, 187, 205, 79, 83, 84, 85, 86, 203, 208,
	87, 22, 50, 51, 53, 103, 54, 55, 166, 98,
	184, 23, 105, 22, 50, 51, 53, 52, 54, 55,
	6, 79, 83, 84, 85, 86, 39, 130, 87, 104,
	40, 78, 77, 66, 56, 67, 68, 69, 64, 101,
	49, 62, 167, 102, 39, 57, 59, 139, 40, 22,
	52, 157, 56, 32, 76, 15, 39, 22, 49, 11,
	40, 6, 50, 57, 56, 73, 200, 201, 52, 28,
	49, 32, 22, 71, 8, 57, 22, 50, 51, 53,
	52, 54, 55, 32, 1, 133, 199, 7, 22, 50,
	51, 53, 12, 54, 55, 46, 160, 22, 50, 51,
	53, 94, 54, 55, 79, 83, 84, 85, 86, 38,
	66, 87, 67, 68, 69, 64, 25, 26, 62, 39,
	21, 48, 66, 40, 67, 68, 69, 56, 37, 120,
	119, 74, 4, 49, 4, 31, 34, 30, 57, 56,
	35, 186, 33, 52, 165, 49, 32, 100, 56, 175,
	57, 58, 65, 14, 49, 52, 10, 5, 96, 57,
	3, 2, 0, 0, 52, 0, 0, 96, 0, 79,
	83, 84, 85, 86, 0, 66, 87, 67, 68, 69,
	64, 170, 171, 62,
}
var yyPact = [...]int{

	260, -1000, 260, -61, -1000, 257, 83, -1000, -1000, -1000,
	249, 147, 207, -1000, 228, 282, -1000, -1000, 285, 263,
	255, -1000, -1000, 260, -1000, -43, 244, 192, -1000, 24,
	-1000, -1000, 219, -1000, -1000, -1000, -1000, -1000, -1000, 38,
	282, 41, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 303,
	-1000, -1000, 156, -1000, -1000, -1000, -1000, 40, 220, 223,
	192, 181, 297, 128, 166, 87, -1000, -1000, -1000, -1000,
	-1000, 278, -1000, 278, 27, 207, 278, 282, 282, 303,
	294, 97, 127, 303, 303, 303, 303, 303, -1000, -27,
	1, -1000, -1000, 83, 66, 134, 58, -1000, -1000, 83,
	230, 282, 303, 87, 163, 128, 155, -1000, 87, -1000,
	-1000, -1000, -1000, -1000, -1000, 108, -1000, -16, 134, 38,
	38, 114, 303, 33, 303, 77, 77, -38, -38, 134,
	-1000, -1000, 26, 189, 79, 282, 162, 230, 185, 222,
	192, -1000, -40, 350, 87, 133, 87, 208, -1000, -1000,
	-1000, -1000, 310, 58, -1000, 134, -1000, -1000, -1000, 21,
	61, 66, 303, 81, 6, -1000, 149, 303, 303, -1000,
	282, 29, 208, 87, 350, 303, 5, 282, -1000, 134,
	303, -1000, 267, -1000, 149, -1000, -51, 245, -1000, 192,
	278, 208, -1000, 134, -1000, -1000, 134, -14, 303, -1000,
	-1000, -1000, -35, -54, -1000, -1000, -1000, 278, -1000,
}
var yyPgo = [...]int{

	0, 294, 371, 370, 341, 367, 366, 363, 361, 357,
	8, 354, 352, 2, 351, 13, 350, 347, 346, 345,
	6, 338, 18, 331, 7, 330, 10, 327, 326, 319,
	311, 306, 19, 4, 3, 120, 1, 59, 11, 0,
	9, 5, 166, 16, 305, 14, 295, 28, 296, 176,
	15, 159, 284,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 4, 5, 6, 6, 6,
	7, 7, 8, 8, 9, 9, 10, 10, 11, 11,
	12, 12, 12, 12, 12, 12, 12, 13, 13, 13,
	13, 13, 13, 13, 14, 48, 48, 48, 15, 16,
	17, 17, 17, 17, 17, 17, 17, 17, 17, 17,
	18, 18, 18, 18, 18, 19, 19, 19, 20, 20,
	20, 20, 21, 21, 22, 22, 22, 23, 23, 24,
	24, 24, 24, 24, 24, 24, 25, 25, 25, 25,
	25, 26, 26, 26, 27, 27, 28, 28, 29, 30,
	30, 31, 31, 32, 32, 33, 33, 34, 34, 35,
	35, 36, 36, 37, 37, 38, 38, 39, 40, 41,
	41, 42, 42, 43, 44, 45, 46, 46, 47, 47,
	49, 49, 50, 50, 51, 51, 51, 51, 52, 52,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 1, 7, 3, 0, 2, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 2, 0, 1, 1, 3, 3,
	3, 4, 4, 6, 6, 4, 4, 4, 4, 2,
	3, 3, 3, 3, 3, 3, 3, 2, 1, 1,
	1, 3, 4, 1, 0, 2, 2, 5, 7, 1,
	2, 3, 1, 2, 3, 1, 5, 5, 6, 6,
	4, 0, 2, 4, 1, 1, 1, 3, 5, 0,
	1, 0, 2, 1, 3, 1, 3, 1, 3, 1,
	3, 1, 3, 1, 3, 4, 2, 1, 1, 1,
	2, 1, 2, 1, 1, 1, 0, 1, 0, 1,
	0, 1, 0, 1, 0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, 11, -1, -52, 78,
	-6, 12, -46, 56, -7, 16, 21, -35, -24, -39,
	-15, -25, 4, 74, -37, -28, -27, -20, 72, -13,
	-17, -19, 74, -12, -18, -16, -15, -21, -29, 47,
	51, -39, -40, -41, -42, -43, -44, -45, -23, 61,
	5, 6, 71, 7, 9, 10, 55, 66, -8, 28,
	-20, -49, 43, -51, 40, 77, 35, 37, 38, 39,
	-39, 20, -39, 20, -4, 77, 20, 50, 49, 69,
	68, 54, -47, 70, 71, 72, 73, 76, 51, -13,
	-20, -15, -20, 74, -30, -13, 74, -41, -42, 74,
	-9, 29, 30, 34, -49, -51, -50, 36, 34, -35,
	-39, -39, 75, -37, -39, -20, -20, -13, -13, 46,
	45, -47, 52, 48, 53, -13, -13, -13, -13, -13,
	75, 75, -22, -46, -38, 62, -13, -22, -10, 27,
	-20, -32, -13, -24, 34, -50, 34, -24, -15, -15,
	-43, -45, -13, 74, -15, -13, 75, 72, -33, -20,
	-31, -38, 64, -20, -10, -11, 33, 30, 77, -26,
	41, 42, -24, 34, -24, 49, -32, 77, 65, -13,
	63, 75, 67, -41, 71, -34, -14, -13, -32, -20,
	74, -24, -26, -13, 75, -33, -13, -40, 77, -48,
	31, 32, -36, -39, 75, -34, 75, 77, -36,
}
var yyDef = [...]int{

	1, -2, 1, 128, 4, 7, 116, 2, 3, 129,
	10, 0, 0, 117, 12, 0, 8, 9, -2, 69,
	72, 75, 107, 0, 6, 103, 86, 84, 85, -2,
	59, 60, 0, 27, 28, 29, 30, 31, 32, 0,
	0, 20, 21, 22, 23, 24, 25, 26, 63, 89,
	108, 109, 0, 111, 113, 114, 115, 0, 14, 0,
	11, 0, -2, 122, 0, 0, 121, 125, 126, 127,
	70, 0, 73, 0, 0, 0, 0, 0, 0, 0,
	0, 118, 0, 0, 0, 0, 0, 0, 119, -2,
	0, 49, 57, -2, 0, 90, 0, 110, 112, -2,
	16, 0, 0, 0, 0, 122, 0, 123, 0, 100,
	71, 74, 38, 104, 87, 55, 56, 39, 40, 0,
	0, 0, 0, 0, 0, 50, 51, 52, 53, 54,
	33, 61, 0, 0, 91, 0, 0, 16, 18, 0,
	15, 13, 93, 81, 0, 0, 0, 80, 47, 48,
	41, 42, 0, 0, 45, 46, 62, 65, 66, 95,
	0, 106, 0, 0, 0, 5, 0, 0, 0, 76,
	0, 0, 77, 0, 81, 0, 0, 0, 88, 92,
	0, 67, 0, 19, 0, 17, 97, 35, 94, 82,
	0, 79, 78, 43, 44, 96, 105, 0, 0, 34,
	36, 37, 0, 101, 68, 98, 83, 0, 102,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 76, 3, 3,
	74, 75, 72, 70, 77, 71, 3, 73, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 78,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69,
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
		//line parser.y:96
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:101
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:108
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:114
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:120
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
	case 6:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:134
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:140
		{
			yyVAL.expression = nil
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:144
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: []Expression{Dual{Dual: yyDollar[2].token.Literal}}}
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:148
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:154
		{
			yyVAL.expression = nil
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:158
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:164
		{
			yyVAL.expression = nil
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:168
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:174
		{
			yyVAL.expression = nil
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:178
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 16:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:184
		{
			yyVAL.expression = nil
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:188
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 18:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:194
		{
			yyVAL.expression = nil
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:198
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:204
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:208
		{
			yyVAL.expression = yyDollar[1].text
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:212
		{
			yyVAL.expression = yyDollar[1].integer
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:216
		{
			yyVAL.expression = yyDollar[1].float
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:220
		{
			yyVAL.expression = yyDollar[1].ternary
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:224
		{
			yyVAL.expression = yyDollar[1].datetime
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:228
		{
			yyVAL.expression = yyDollar[1].null
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:234
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:238
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:242
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:246
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:250
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:254
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:258
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:264
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 35:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:270
		{
			yyVAL.token = Token{}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:274
		{
			yyVAL.token = yyDollar[1].token
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:278
		{
			yyVAL.token = yyDollar[1].token
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:284
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:290
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
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:313
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:317
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:321
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 43:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:325
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 44:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:329
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:333
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 46:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:337
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 48:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:349
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:355
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:359
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:363
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:367
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:371
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:377
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:381
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:385
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:391
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:399
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:403
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:409
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:413
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 64:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:419
		{
			yyVAL.expression = Option{}
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:423
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:427
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 67:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:433
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression}
		}
	case 68:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:437
		{
			yyVAL.expression = GroupConcat{GroupConcat: yyDollar[1].token.Literal, Option: yyDollar[3].expression.(Option), OrderBy: yyDollar[4].expression, SeparatorLit: yyDollar[5].token.Literal, Separator: yyDollar[6].text.Value()}
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:443
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:447
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:451
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:455
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:459
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 74:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:463
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:467
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 76:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:473
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 77:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:477
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 78:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:481
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 79:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:485
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 80:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:489
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 81:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:495
		{
			yyVAL.expression = nil
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:499
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 83:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:503
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:509
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 85:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:513
		{
			yyVAL.expression = AllColumns{}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:519
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:523
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 88:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:529
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 89:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:535
		{
			yyVAL.expression = nil
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:539
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 91:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:545
		{
			yyVAL.expression = nil
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:549
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:555
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:559
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:565
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:569
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 97:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:575
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:579
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 99:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:585
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:589
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 101:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:595
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:599
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 103:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:605
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 104:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:609
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 105:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:615
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:619
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:625
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:631
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:637
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 110:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:641
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:648
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 112:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:652
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:659
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:665
		{
			yyVAL.datetime = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:671
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 116:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:677
		{
			yyVAL.token = Token{}
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:681
		{
			yyVAL.token = yyDollar[1].token
		}
	case 118:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:687
		{
			yyVAL.token = Token{}
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:691
		{
			yyVAL.token = yyDollar[1].token
		}
	case 120:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:697
		{
			yyVAL.token = Token{}
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:701
		{
			yyVAL.token = yyDollar[1].token
		}
	case 122:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:707
		{
			yyVAL.token = Token{}
		}
	case 123:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:711
		{
			yyVAL.token = yyDollar[1].token
		}
	case 124:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:717
		{
			yyVAL.token = Token{}
		}
	case 125:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:721
		{
			yyVAL.token = yyDollar[1].token
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:725
		{
			yyVAL.token = yyDollar[1].token
		}
	case 127:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:729
		{
			yyVAL.token = yyDollar[1].token
		}
	case 128:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:735
		{
			yyVAL.token = Token{}
		}
	case 129:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:739
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
