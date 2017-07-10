package parser

import (
	"errors"
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
		l.err = errors.New(fmt.Sprintf("%s: unexpected %s [L:%d C:%d]", e, l.token.Literal, l.token.Line, l.token.Char))
	} else {
		l.err = errors.New(fmt.Sprintf("%s [L:%d C:%d]", e, l.token.Line, l.token.Char))
	}
}

type Token struct {
	Token   int
	Literal string
	Quoted  bool
	Line    int
	Char    int
}

func (t *Token) IsEmpty() bool {
	return len(t.Literal) < 1
}
