package parser

import (
	"testing"
)

type scanResult struct {
	Token   int
	Literal string
	Quoted  bool
}

var scanTests = []struct {
	Name   string
	Input  string
	Output []scanResult
	Error  string
}{
	{
		Name:  "Identifier",
		Input: "identifier",
		Output: []scanResult{
			{
				Token:   IDENTIFIER,
				Literal: "identifier",
			},
		},
	},
	{
		Name:  "QuotedIdentifier",
		Input: "`identifier`",
		Output: []scanResult{
			{
				Token:   IDENTIFIER,
				Literal: "identifier",
				Quoted:  true,
			},
		},
	},
	{
		Name:  "QuotedString",
		Input: "\"string\\\"\"",
		Output: []scanResult{
			{
				Token:   STRING,
				Literal: "string\"",
			},
		},
	},
	{
		Name:  "QuotedString(Single-Quote)",
		Input: "'strin\\'g string'",
		Output: []scanResult{
			{
				Token:   STRING,
				Literal: "strin'g string",
			},
		},
	},
	{
		Name:  "Integer",
		Input: "1",
		Output: []scanResult{
			{
				Token:   INTEGER,
				Literal: "1",
			},
		},
	},
	{
		Name:  "Float",
		Input: "1.234",
		Output: []scanResult{
			{
				Token:   FLOAT,
				Literal: "1.234",
			},
		},
	},
	{
		Name:  "Ternary",
		Input: "true",
		Output: []scanResult{
			{
				Token:   TERNARY,
				Literal: "true",
			},
		},
	},
	{
		Name:  "Datetime",
		Input: "\"2012-05-21 12:00:00\"",
		Output: []scanResult{
			{
				Token:   DATETIME,
				Literal: "2012-05-21 12:00:00",
			},
		},
	},
	{
		Name:  "Datetime(RFC3339)",
		Input: "\"2012-05-21T12:00:00-12:00\"",
		Output: []scanResult{
			{
				Token:   DATETIME,
				Literal: "2012-05-21T12:00:00-12:00",
			},
		},
	},
	{
		Name:  "Variable",
		Input: "@var",
		Output: []scanResult{
			{
				Token:   VARIABLE,
				Literal: "@var",
			},
		},
	},
	{
		Name:  "ComparisonOperator",
		Input: "=",
		Output: []scanResult{
			{
				Token:   COMPARISON_OP,
				Literal: "=",
			},
		},
	},
	{
		Name:  "StringOperator",
		Input: "||",
		Output: []scanResult{
			{
				Token:   STRING_OP,
				Literal: "||",
			},
		},
	},
	{
		Name:  "SubstitutionOperator",
		Input: ":=",
		Output: []scanResult{
			{
				Token:   SUBSTITUTION_OP,
				Literal: ":=",
			},
		},
	},
	{
		Name:  "UncategorizedOperator",
		Input: "====",
		Output: []scanResult{
			{
				Token:   UNCATEGORIZED,
				Literal: "====",
			},
		},
	},
	{
		Name:  "Keyword",
		Input: "select",
		Output: []scanResult{
			{
				Token:   SELECT,
				Literal: "select",
			},
		},
	},
	{
		Name:  "PassThrough",
		Input: ",",
		Output: []scanResult{
			{
				Token:   int(','),
				Literal: ",",
			},
		},
	},
	{
		Name:  "Statement",
		Input: "identifier   'string', \n 1-2",
		Output: []scanResult{
			{
				Token:   IDENTIFIER,
				Literal: "identifier",
			},
			{
				Token:   STRING,
				Literal: "string",
			},
			{
				Token:   int(','),
				Literal: ",",
			},
			{
				Token:   INTEGER,
				Literal: "1",
			},
			{
				Token:   int('-'),
				Literal: "-",
			},
			{
				Token:   INTEGER,
				Literal: "2",
			},
		},
	},
	{
		Name:  "LiteralNotTerminatedError",
		Input: "\"string",
		Error: "literal not terminated",
	},
}

func TestScan(t *testing.T) {
	for _, v := range scanTests {
		s := new(Scanner).Init(v.Input)

		tokenCount := 0
		for {
			token, literal, quoted, err := s.Scan()
			tokenCount++

			if err != nil {
				if v.Error == "" {
					t.Errorf("%s, token %d: unexpected error %q", v.Name, tokenCount, err.Error())
				} else if v.Error != err.Error() {
					t.Errorf("%s, token %d: error %q, want error %q", v.Name, tokenCount, err.Error(), v.Error)
				}
				break
			}
			if v.Error != "" {
				t.Errorf("%s, token %d: no error, want error %q", v.Name, tokenCount, v.Error)
				break
			}

			if token == EOF {
				break
			}

			if len(v.Output) < tokenCount {
				break
			}
			expect := v.Output[tokenCount-1]
			if token != expect.Token {
				t.Errorf("%s, token %d: token = %s, want %s", v.Name, tokenCount, TokenLiteral(token), TokenLiteral(expect.Token))
			}
			if literal != expect.Literal {
				t.Errorf("%s, token %d: literal = %q, want %q", v.Name, tokenCount, literal, expect.Literal)
			}
			if quoted != expect.Quoted {
				t.Errorf("%s, token %d: quoted = %t, want %t", v.Name, tokenCount, quoted, expect.Quoted)
			}
		}

		tokenCount--
		if tokenCount != len(v.Output) {
			t.Errorf("%s: scan %d token(s) in a statement, want %d token(s)", v.Name, tokenCount, len(v.Output))
		}
	}
}

var tokenLiteralTests = map[int]string{
	SELECT: "SELECT",
	43:     "+",
}

func TestTokenLiteral(t *testing.T) {
	for k, v := range tokenLiteralTests {
		n := TokenLiteral(k)
		if n != v {
			t.Errorf("token literal = %q, want %q for %d", n, v, k)
		}
	}
}
