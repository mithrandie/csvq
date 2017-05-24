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
const COMPARISON_OP = 57408
const STRING_OP = 57409

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

//line parser.y:719

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
	34, 116,
	36, 120,
	-2, 96,
	-1, 29,
	48, 114,
	52, 114,
	53, 114,
	-2, 58,
	-1, 59,
	36, 120,
	-2, 116,
	-1, 86,
	49, 58,
	50, 58,
	-2, 114,
	-1, 90,
	73, 63,
	-2, 112,
}

const yyPrivate = 57344

const yyLast = 390

var yyAct = [...]int{

	41, 193, 29, 177, 153, 43, 27, 18, 47, 45,
	163, 156, 19, 36, 102, 12, 9, 197, 79, 189,
	67, 69, 57, 72, 136, 20, 85, 84, 196, 78,
	76, 80, 81, 82, 83, 86, 182, 84, 162, 87,
	186, 77, 76, 80, 81, 82, 83, 89, 126, 84,
	151, 92, 108, 88, 82, 83, 172, 94, 84, 24,
	22, 148, 63, 19, 64, 65, 66, 61, 23, 106,
	59, 107, 90, 22, 110, 85, 20, 16, 78, 113,
	114, 111, 112, 121, 122, 123, 124, 125, 131, 50,
	77, 76, 80, 81, 82, 83, 132, 117, 84, 17,
	19, 137, 62, 130, 135, 19, 129, 138, 80, 81,
	82, 83, 142, 20, 84, 131, 140, 157, 20, 75,
	74, 147, 13, 150, 75, 74, 146, 145, 23, 143,
	144, 60, 109, 149, 50, 52, 154, 85, 158, 53,
	19, 23, 19, 75, 74, 171, 75, 166, 127, 168,
	103, 137, 176, 20, 58, 20, 167, 174, 141, 44,
	173, 139, 105, 104, 179, 137, 175, 99, 19, 75,
	74, 181, 185, 170, 160, 183, 187, 188, 154, 184,
	161, 20, 94, 194, 119, 54, 98, 180, 118, 120,
	97, 101, 179, 195, 56, 134, 73, 51, 194, 198,
	22, 49, 50, 52, 15, 53, 46, 22, 49, 50,
	52, 95, 53, 46, 100, 63, 22, 64, 65, 66,
	22, 49, 50, 52, 22, 53, 46, 6, 22, 11,
	6, 8, 70, 22, 49, 50, 52, 190, 53, 46,
	68, 71, 4, 39, 4, 1, 42, 40, 7, 155,
	39, 54, 91, 38, 40, 25, 26, 48, 54, 21,
	128, 37, 31, 39, 48, 51, 152, 40, 32, 34,
	30, 54, 51, 28, 35, 32, 39, 48, 178, 33,
	40, 159, 133, 96, 54, 51, 191, 192, 32, 55,
	48, 22, 49, 50, 52, 14, 53, 46, 51, 10,
	5, 32, 3, 22, 49, 50, 52, 2, 53, 46,
	6, 0, 22, 49, 50, 52, 169, 53, 46, 0,
	0, 0, 76, 80, 81, 82, 83, 0, 0, 84,
	0, 0, 116, 115, 76, 80, 81, 82, 83, 0,
	0, 84, 54, 76, 80, 81, 82, 83, 48, 126,
	84, 0, 0, 0, 54, 0, 51, 0, 0, 93,
	48, 0, 63, 54, 64, 65, 66, 61, 51, 48,
	59, 93, 76, 80, 81, 82, 83, 51, 0, 84,
	93, 63, 0, 64, 65, 66, 61, 164, 165, 59,
}
var yyPact = [...]int{

	219, -1000, 219, -60, -1000, 217, 66, -1000, -1000, -1000,
	188, 56, 203, -1000, 166, 229, -1000, -1000, 27, 220,
	212, -1000, -1000, 219, -1000, -52, 176, 120, -1000, 24,
	-1000, -1000, 216, -1000, -1000, -1000, -1000, -1000, -1000, -4,
	229, 0, -1000, -1000, -1000, -1000, -1000, -1000, 308, -1000,
	-1000, 128, -1000, -1000, -1000, 161, 156, 120, 133, 180,
	114, 129, 69, -1000, -1000, -1000, -1000, -1000, 224, -1000,
	224, -21, 203, 224, 229, 229, 308, 287, 86, 136,
	308, 308, 308, 308, 308, -1000, -25, 75, -1000, -1000,
	66, 26, 305, 299, -1000, -1000, 168, 229, 308, 69,
	127, 114, 124, -1000, 69, -1000, -1000, -1000, -1000, -1000,
	-1000, 97, -1000, 40, 305, -4, -4, 130, 308, -11,
	308, -16, -16, -47, -47, 305, -1000, -1000, -23, 196,
	53, 229, 276, 141, 150, 120, -1000, -37, 346, 69,
	122, 69, 327, -1000, -1000, -1000, -1000, 267, 299, -1000,
	305, -1000, -1000, -1000, 70, -9, 26, 308, 94, -1000,
	83, 308, 308, -1000, 229, -36, 327, 69, 346, 308,
	-33, 229, -1000, 305, 308, -1000, 83, -1000, -56, 255,
	-1000, 120, 224, 327, -1000, 305, -1000, -1000, 305, 308,
	-1000, -1000, -1000, -45, -58, -1000, -1000, 224, -1000,
}
var yyPgo = [...]int{

	0, 245, 307, 302, 241, 300, 299, 295, 289, 283,
	282, 281, 279, 2, 278, 13, 274, 270, 269, 262,
	6, 261, 260, 7, 259, 10, 256, 255, 253, 252,
	249, 24, 4, 3, 99, 1, 59, 11, 0, 246,
	5, 159, 9, 8, 15, 18, 237, 154, 14, 131,
	231,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 3, 4, 5, 6, 6, 6,
	7, 7, 8, 8, 9, 9, 10, 10, 11, 11,
	12, 12, 12, 12, 12, 12, 12, 13, 13, 13,
	13, 13, 13, 13, 14, 46, 46, 46, 15, 16,
	17, 17, 17, 17, 17, 17, 17, 17, 17, 17,
	18, 18, 18, 18, 18, 19, 19, 19, 20, 20,
	20, 20, 21, 22, 22, 22, 23, 23, 23, 23,
	23, 23, 23, 24, 24, 24, 24, 24, 25, 25,
	25, 26, 26, 27, 27, 28, 29, 29, 30, 30,
	31, 31, 32, 32, 33, 33, 34, 34, 35, 35,
	36, 36, 37, 37, 38, 39, 40, 40, 41, 41,
	42, 43, 44, 44, 45, 45, 47, 47, 48, 48,
	49, 49, 49, 49, 50, 50,
}
var yyR2 = [...]int{

	0, 0, 2, 2, 1, 7, 3, 0, 2, 2,
	0, 2, 0, 3, 0, 2, 0, 3, 0, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 2, 0, 1, 1, 3, 3,
	3, 4, 4, 6, 6, 4, 4, 4, 4, 2,
	3, 3, 3, 3, 3, 3, 3, 2, 1, 1,
	1, 3, 4, 0, 2, 2, 1, 2, 3, 1,
	2, 3, 1, 5, 5, 6, 6, 4, 0, 2,
	4, 1, 1, 1, 3, 5, 0, 1, 0, 2,
	1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
	1, 3, 4, 2, 1, 1, 1, 2, 1, 2,
	1, 1, 0, 1, 0, 1, 0, 1, 0, 1,
	0, 1, 1, 1, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, 11, -1, -50, 76,
	-6, 12, -44, 56, -7, 16, 21, -34, -23, -38,
	-15, -24, 4, 72, -36, -27, -26, -20, 70, -13,
	-17, -19, 72, -12, -18, -16, -15, -21, -28, 47,
	51, -38, -39, -40, -41, -42, 10, -43, 61, 5,
	6, 69, 7, 9, 55, -8, 28, -20, -47, 43,
	-49, 40, 75, 35, 37, 38, 39, -38, 20, -38,
	20, -4, 75, 20, 50, 49, 67, 66, 54, -45,
	68, 69, 70, 71, 74, 51, -13, -20, -15, -20,
	72, -29, -13, 72, -40, -41, -9, 29, 30, 34,
	-47, -49, -48, 36, 34, -34, -38, -38, 73, -36,
	-38, -20, -20, -13, -13, 46, 45, -45, 52, 48,
	53, -13, -13, -13, -13, -13, 73, 73, -22, -44,
	-37, 62, -13, -10, 27, -20, -31, -13, -23, 34,
	-48, 34, -23, -15, -15, -42, -43, -13, 72, -15,
	-13, 73, 70, -32, -20, -30, -37, 64, -20, -11,
	33, 30, 75, -25, 41, 42, -23, 34, -23, 49,
	-31, 75, 65, -13, 63, -40, 69, -33, -14, -13,
	-31, -20, 72, -23, -25, -13, 73, -32, -13, 75,
	-46, 31, 32, -35, -38, -33, 73, 75, -35,
}
var yyDef = [...]int{

	1, -2, 1, 124, 4, 7, 112, 2, 3, 125,
	10, 0, 0, 113, 12, 0, 8, 9, -2, 66,
	69, 72, 104, 0, 6, 100, 83, 81, 82, -2,
	59, 60, 0, 27, 28, 29, 30, 31, 32, 0,
	0, 20, 21, 22, 23, 24, 25, 26, 86, 105,
	106, 0, 108, 110, 111, 14, 0, 11, 0, -2,
	118, 0, 0, 117, 121, 122, 123, 67, 0, 70,
	0, 0, 0, 0, 0, 0, 0, 0, 114, 0,
	0, 0, 0, 0, 0, 115, -2, 0, 49, 57,
	-2, 0, 87, 0, 107, 109, 16, 0, 0, 0,
	0, 118, 0, 119, 0, 97, 68, 71, 38, 101,
	84, 55, 56, 39, 40, 0, 0, 0, 0, 0,
	0, 50, 51, 52, 53, 54, 33, 61, 0, 0,
	88, 0, 0, 18, 0, 15, 13, 90, 78, 0,
	0, 0, 77, 47, 48, 41, 42, 0, 0, 45,
	46, 62, 64, 65, 92, 0, 103, 0, 0, 5,
	0, 0, 0, 73, 0, 0, 74, 0, 78, 0,
	0, 0, 85, 89, 0, 19, 0, 17, 94, 35,
	91, 79, 0, 76, 75, 43, 44, 93, 102, 0,
	34, 36, 37, 0, 98, 95, 80, 0, 99,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 74, 3, 3,
	72, 73, 70, 68, 75, 69, 3, 71, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 76,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67,
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
		//line parser.y:92
		{
			yyVAL.program = nil
			yylex.(*Lexer).program = yyVAL.program
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:97
		{
			yyVAL.program = append([]Statement{yyDollar[1].statement}, yyDollar[2].program...)
			yylex.(*Lexer).program = yyVAL.program
		}
	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:104
		{
			yyVAL.statement = yyDollar[1].expression
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:110
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 5:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.y:116
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
		//line parser.y:130
		{
			yyVAL.expression = SelectClause{Select: yyDollar[1].token.Literal, Distinct: yyDollar[2].token, Fields: yyDollar[3].expressions}
		}
	case 7:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:136
		{
			yyVAL.expression = nil
		}
	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:140
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: []Expression{Dual{Dual: yyDollar[2].token.Literal}}}
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:144
		{
			yyVAL.expression = FromClause{From: yyDollar[1].token.Literal, Tables: yyDollar[2].expressions}
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:150
		{
			yyVAL.expression = nil
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:154
		{
			yyVAL.expression = WhereClause{Where: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:160
		{
			yyVAL.expression = nil
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:164
		{
			yyVAL.expression = GroupByClause{GroupBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:170
		{
			yyVAL.expression = nil
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:174
		{
			yyVAL.expression = HavingClause{Having: yyDollar[1].token.Literal, Filter: yyDollar[2].expression}
		}
	case 16:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:180
		{
			yyVAL.expression = nil
		}
	case 17:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:184
		{
			yyVAL.expression = OrderByClause{OrderBy: yyDollar[1].token.Literal + " " + yyDollar[2].token.Literal, Items: yyDollar[3].expressions}
		}
	case 18:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:190
		{
			yyVAL.expression = nil
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:194
		{
			yyVAL.expression = LimitClause{Limit: yyDollar[1].token.Literal, Number: yyDollar[2].integer.Value()}
		}
	case 20:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:200
		{
			yyVAL.expression = yyDollar[1].identifier
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:204
		{
			yyVAL.expression = yyDollar[1].text
		}
	case 22:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:208
		{
			yyVAL.expression = yyDollar[1].integer
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:212
		{
			yyVAL.expression = yyDollar[1].float
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:216
		{
			yyVAL.expression = yyDollar[1].ternary
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:220
		{
			yyVAL.expression = NewDatetimeFromString(yyDollar[1].token.Literal)
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:224
		{
			yyVAL.expression = yyDollar[1].null
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:230
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:234
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:238
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:242
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:246
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:250
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 33:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:254
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 34:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:260
		{
			yyVAL.expression = OrderItem{Item: yyDollar[1].expression, Direction: yyDollar[2].token}
		}
	case 35:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:266
		{
			yyVAL.token = Token{}
		}
	case 36:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:270
		{
			yyVAL.token = yyDollar[1].token
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:274
		{
			yyVAL.token = yyDollar[1].token
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:280
		{
			yyVAL.expression = Subquery{Query: yyDollar[2].expression.(SelectQuery)}
		}
	case 39:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:286
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
		//line parser.y:309
		{
			yyVAL.expression = Comparison{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:313
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].ternary, Negation: yyDollar[3].token}
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:317
		{
			yyVAL.expression = Is{Is: yyDollar[2].token.Literal, LHS: yyDollar[1].expression, RHS: yyDollar[4].null, Negation: yyDollar[3].token}
		}
	case 43:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:321
		{
			yyVAL.expression = Between{Between: yyDollar[3].token.Literal, And: yyDollar[5].token.Literal, LHS: yyDollar[1].expression, Low: yyDollar[4].expression, High: yyDollar[6].expression, Negation: yyDollar[2].token}
		}
	case 44:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:325
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, List: yyDollar[5].expressions, Negation: yyDollar[2].token}
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:329
		{
			yyVAL.expression = In{In: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Query: yyDollar[4].expression.(Subquery), Negation: yyDollar[2].token}
		}
	case 46:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:333
		{
			yyVAL.expression = Like{Like: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Pattern: yyDollar[4].expression, Negation: yyDollar[2].token}
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:337
		{
			yyVAL.expression = Any{Any: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 48:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:341
		{
			yyVAL.expression = All{All: yyDollar[3].token.Literal, LHS: yyDollar[1].expression, Operator: yyDollar[2].token, Query: yyDollar[4].expression.(Subquery)}
		}
	case 49:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:345
		{
			yyVAL.expression = Exists{Exists: yyDollar[1].token.Literal, Query: yyDollar[2].expression.(Subquery)}
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:351
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('+'), RHS: yyDollar[3].expression}
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:355
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('-'), RHS: yyDollar[3].expression}
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:359
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('*'), RHS: yyDollar[3].expression}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:363
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('/'), RHS: yyDollar[3].expression}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:367
		{
			yyVAL.expression = Arithmetic{LHS: yyDollar[1].expression, Operator: int('%'), RHS: yyDollar[3].expression}
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:373
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:377
		{
			yyVAL.expression = Logic{LHS: yyDollar[1].expression, Operator: yyDollar[2].token, RHS: yyDollar[3].expression}
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:381
		{
			yyVAL.expression = Logic{LHS: nil, Operator: yyDollar[1].token, RHS: yyDollar[2].expression}
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:387
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:391
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 60:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:395
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:399
		{
			yyVAL.expression = Parentheses{Expr: yyDollar[2].expression}
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:405
		{
			yyVAL.expression = Function{Name: yyDollar[1].identifier.Literal, Option: yyDollar[3].expression.(Option)}
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:411
		{
			yyVAL.expression = Option{}
		}
	case 64:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:415
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: []Expression{AllColumns{}}}
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:419
		{
			yyVAL.expression = Option{Distinct: yyDollar[1].token, Args: yyDollar[2].expressions}
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:425
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier}
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:429
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, Alias: yyDollar[2].identifier}
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:433
		{
			yyVAL.expression = Table{Object: yyDollar[1].identifier, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:437
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:441
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, Alias: yyDollar[2].identifier}
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:445
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:449
		{
			yyVAL.expression = Table{Object: yyDollar[1].expression}
		}
	case 73:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:455
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: yyDollar[5].expression}
		}
	case 74:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:459
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[3].token, Condition: nil}
		}
	case 75:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:463
		{
			yyVAL.expression = Join{Join: yyDollar[4].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[5].expression.(Table), Natural: Token{}, JoinType: yyDollar[3].token, Direction: yyDollar[2].token, Condition: yyDollar[6].expression}
		}
	case 76:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.y:467
		{
			yyVAL.expression = Join{Join: yyDollar[5].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[6].expression.(Table), Natural: yyDollar[2].token, JoinType: yyDollar[4].token, Direction: yyDollar[3].token, Condition: nil}
		}
	case 77:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:471
		{
			yyVAL.expression = Join{Join: yyDollar[3].token.Literal, Table: yyDollar[1].expression.(Table), JoinTable: yyDollar[4].expression.(Table), Natural: Token{}, JoinType: yyDollar[2].token, Condition: nil}
		}
	case 78:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:477
		{
			yyVAL.expression = nil
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:481
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, On: yyDollar[2].expression}
		}
	case 80:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:485
		{
			yyVAL.expression = JoinCondition{Literal: yyDollar[1].token.Literal, Using: yyDollar[3].expressions}
		}
	case 81:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:491
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 82:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:495
		{
			yyVAL.expression = AllColumns{}
		}
	case 83:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:501
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression}
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:505
		{
			yyVAL.expression = Field{Object: yyDollar[1].expression, As: yyDollar[2].token, Alias: yyDollar[3].identifier}
		}
	case 85:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.y:511
		{
			yyVAL.expression = Case{Case: yyDollar[1].token.Literal, End: yyDollar[5].token.Literal, Value: yyDollar[2].expression, When: yyDollar[3].expressions, Else: yyDollar[4].expression}
		}
	case 86:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:517
		{
			yyVAL.expression = nil
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:521
		{
			yyVAL.expression = yyDollar[1].expression
		}
	case 88:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:527
		{
			yyVAL.expression = nil
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:531
		{
			yyVAL.expression = CaseElse{Else: yyDollar[1].token.Literal, Result: yyDollar[2].expression}
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:537
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:541
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 92:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:547
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:551
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 94:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:557
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:561
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 96:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:567
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 97:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:571
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 98:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:577
		{
			yyVAL.expressions = []Expression{yyDollar[1].identifier}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:581
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].identifier}, yyDollar[3].expressions...)
		}
	case 100:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:587
		{
			yyVAL.expressions = []Expression{yyDollar[1].expression}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.y:591
		{
			yyVAL.expressions = append([]Expression{yyDollar[1].expression}, yyDollar[3].expressions...)
		}
	case 102:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.y:597
		{
			yyVAL.expressions = []Expression{CaseWhen{When: yyDollar[1].token.Literal, Then: yyDollar[3].token.Literal, Condition: yyDollar[2].expression, Result: yyDollar[4].expression}}
		}
	case 103:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:601
		{
			yyVAL.expressions = append(yyDollar[1].expressions, yyDollar[2].expressions...)
		}
	case 104:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:607
		{
			yyVAL.identifier = Identifier{Literal: yyDollar[1].token.Literal, Quoted: yyDollar[1].token.Quoted}
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:613
		{
			yyVAL.text = NewString(yyDollar[1].token.Literal)
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:619
		{
			yyVAL.integer = NewIntegerFromString(yyDollar[1].token.Literal)
		}
	case 107:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:623
		{
			i := yyDollar[2].integer.Value() * -1
			yyVAL.integer = NewInteger(i)
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:630
		{
			yyVAL.float = NewFloatFromString(yyDollar[1].token.Literal)
		}
	case 109:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.y:634
		{
			f := yyDollar[2].float.Value() * -1
			yyVAL.float = NewFloat(f)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:641
		{
			yyVAL.ternary = NewTernaryFromString(yyDollar[1].token.Literal)
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:647
		{
			yyVAL.null = NewNullFromString(yyDollar[1].token.Literal)
		}
	case 112:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:653
		{
			yyVAL.token = Token{}
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:657
		{
			yyVAL.token = yyDollar[1].token
		}
	case 114:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:663
		{
			yyVAL.token = Token{}
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:667
		{
			yyVAL.token = yyDollar[1].token
		}
	case 116:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:673
		{
			yyVAL.token = Token{}
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:677
		{
			yyVAL.token = yyDollar[1].token
		}
	case 118:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:683
		{
			yyVAL.token = Token{}
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:687
		{
			yyVAL.token = yyDollar[1].token
		}
	case 120:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:693
		{
			yyVAL.token = Token{}
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:697
		{
			yyVAL.token = yyDollar[1].token
		}
	case 122:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:701
		{
			yyVAL.token = yyDollar[1].token
		}
	case 123:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:705
		{
			yyVAL.token = yyDollar[1].token
		}
	case 124:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.y:711
		{
			yyVAL.token = Token{}
		}
	case 125:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.y:715
		{
			yyVAL.token = Token{Token: ';', Literal: string(';')}
		}
	}
	goto yystack /* stack new state and value */
}
