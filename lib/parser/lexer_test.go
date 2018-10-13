package parser

import "testing"

func TestToken_IsEmpty(t *testing.T) {
	token := Token{}
	if !token.IsEmpty() {
		t.Error("Token.Empty() is false, want true for empty token")
	}
}

func TestLexer_Error(t *testing.T) {
	lexer := Lexer{
		token: Token{
			Token:   SELECT,
			Literal: "select",
		},
	}
	message := "syntax error"

	expect := "syntax error: unexpected token \"select\""
	lexer.Error(message)
	if lexer.err.Error() != expect {
		t.Errorf("error message = %s, want %s for token %v", lexer.err.Error(), expect, lexer.token)
	}

	lexer = Lexer{
		token: Token{
			Token:   AGGREGATE_FUNCTION,
			Literal: "min",
		},
	}
	expect = "syntax error: unexpected token \"min\""
	lexer.Error(message)
	if lexer.err.Error() != expect {
		t.Errorf("error message = %s, want %s for token %v", lexer.err.Error(), expect, lexer.token)
	}

	lexer = Lexer{
		token: Token{
			Token:   SUBSTITUTION_OP,
			Literal: ":=",
		},
	}
	expect = "syntax error: unexpected token \":=\""
	lexer.Error(message)
	if lexer.err.Error() != expect {
		t.Errorf("error message = %s, want %s for token %v", lexer.err.Error(), expect, lexer.token)
	}

	lexer = Lexer{
		token: Token{
			Token: -1,
		},
	}
	expect = "syntax error: unexpected termination"
	lexer.Error(message)
	if lexer.err.Error() != expect {
		t.Errorf("error message = %s, want %s for token %v", lexer.err.Error(), expect, lexer.token)
	}
}
