package parser

import (
	"fmt"
)

type Lexer struct {
	Scanner
	program []Statement
	token   Token
	err     error
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok, err := l.Scan()
	if err != nil {
		l.Error(err.Error())
	}

	lval.token = tok
	l.token = lval.token
	return tok.Token
}

func (l *Lexer) Error(e string) {
	if 0 < l.token.Token {
		var lit string
		if TOKEN_FROM <= l.token.Token && l.token.Token <= KEYWORD_TO {
			lit = TokenLiteral(l.token.Token)
		} else {
			lit = l.token.Literal
		}

		l.err = NewSyntaxError(fmt.Sprintf("%s: unexpected %s", e, lit), l.token)
	} else {
		l.err = NewSyntaxError(fmt.Sprintf("%s", e), l.token)
	}
}

type Token struct {
	Token      int
	Literal    string
	Quoted     bool
	Line       int
	Char       int
	SourceFile string
}

func (t *Token) IsEmpty() bool {
	return len(t.Literal) < 1
}

type SyntaxError struct {
	SourceFile string
	Line       int
	Char       int
	Message    string
}

func (e SyntaxError) Error() string {
	return e.Message
}

func NewSyntaxError(message string, token Token) error {
	return &SyntaxError{
		SourceFile: token.SourceFile,
		Line:       token.Line,
		Char:       token.Char,
		Message:    message,
	}
}
