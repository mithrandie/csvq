package parser

import (
	"errors"
)

type Lexer struct {
	Scanner
	program []Statement
	err     error
}

func (l *Lexer) Lex(lval *yySymType) int {
	token, literal, quoted, err := l.Scan()
	if err != nil {
		l.Error(err.Error())
	}

	lval.token = Token{Token: token, Literal: literal, Quoted: quoted}
	return token
}

func (l *Lexer) Error(e string) {
	l.err = errors.New(e)
}

type Token struct {
	Token   int
	Literal string
	Quoted  bool
}

func (t *Token) IsEmpty() bool {
	return len(t.Literal) < 1
}
