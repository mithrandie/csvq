package json

import (
	"reflect"
	"testing"
)

var parseQueryTests = []struct {
	Input  string
	Expect QueryExpression
	Error  string
}{
	{
		Input:  "",
		Expect: nil,
	},
	{
		Input: "abc",
		Expect: Element{
			Label: "abc",
		},
	},
	{
		Input: "`abc\\`def`",
		Expect: Element{
			Label: "abc`def",
		},
	},
	{
		Input: "[0]",
		Expect: ArrayItem{
			Index: 0,
		},
	},
	{
		Input:  "[]",
		Expect: RowValueExpr{},
	},
	{
		Input:  "{}",
		Expect: TableExpr{},
	},
	{
		Input: "abc.def",
		Expect: Element{
			Label: "abc",
			Child: Element{
				Label: "def",
			},
		},
	},
	{
		Input: "abc[1]",
		Expect: Element{
			Label: "abc",
			Child: ArrayItem{
				Index: 1,
			},
		},
	},
	{
		Input: "abc[]",
		Expect: Element{
			Label: "abc",
			Child: RowValueExpr{},
		},
	},
	{
		Input: "abc{}",
		Expect: Element{
			Label: "abc",
			Child: TableExpr{},
		},
	},
	{
		Input: "abc[2].def",
		Expect: Element{
			Label: "abc",
			Child: ArrayItem{
				Index: 2,
				Child: Element{
					Label: "def",
				},
			},
		},
	},
	{
		Input: "abc[2][1]",
		Expect: Element{
			Label: "abc",
			Child: ArrayItem{
				Index: 2,
				Child: ArrayItem{
					Index: 1,
				},
			},
		},
	},
	{
		Input: "abc[22][]",
		Expect: Element{
			Label: "abc",
			Child: ArrayItem{
				Index: 22,
				Child: RowValueExpr{},
			},
		},
	},
	{
		Input: "abc[2]{}",
		Expect: Element{
			Label: "abc",
			Child: ArrayItem{
				Index: 2,
				Child: TableExpr{},
			},
		},
	},
	{
		Input: "abc[].def",
		Expect: Element{
			Label: "abc",
			Child: RowValueExpr{
				Child: Element{
					Label: "def",
				},
			},
		},
	},
	{
		Input: "abc[][1]",
		Expect: Element{
			Label: "abc",
			Child: RowValueExpr{
				Child: ArrayItem{
					Index: 1,
				},
			},
		},
	},
	{
		Input: "{abc}",
		Expect: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{Label: "abc"},
				},
			},
		},
	},
	{
		Input: "{abc, def as alias}",
		Expect: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{Label: "abc"},
				},
				{
					Element: Element{Label: "def"},
					Alias:   "alias",
				},
			},
		},
	},
	{
		Input: "{abc.def, ghi[2]}",
		Expect: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{
						Label: "abc",
						Child: Element{Label: "def"},
					},
				},
				{
					Element: Element{
						Label: "ghi",
						Child: ArrayItem{
							Index: 2,
						},
					},
				},
			},
		},
	},
	{
		Input: "abc.\"def\".ghi",
		Expect: Element{
			Label: "abc",
			Child: Element{
				Label: "def",
				Child: Element{
					Label: "ghi",
				},
			},
		},
	},
	{
		Input: "abc def",
		Error: "column 5: unexpected token \"def\"",
	},
	{
		Input: "abc[].def{}",
		Error: "column 10: unexpected token \"{\"",
	},
	{
		Input: "abc[].def[]",
		Error: "column 11: unexpected token \"]\"",
	},
	{
		Input: "abc{}.def",
		Error: "column 6: unexpected token \".\"",
	},
	{
		Input: "abc{def[]}",
		Error: "column 9: unexpected token \"]\"",
	},
	{
		Input: "abc{def{}}",
		Error: "column 8: unexpected token \"{\"",
	},
	{
		Input: "abc{def{}}",
		Error: "column 8: unexpected token \"{\"",
	},
	{
		Input: "`abc",
		Error: "column 4: string not terminated",
	},
	{
		Input: "abc[",
		Error: "column 4: unexpected termination",
	},
}

func TestParseQuery(t *testing.T) {
	for _, v := range parseQueryTests {
		result, err := ParseQuery(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", result, v.Expect, v.Input)
		}
	}
}
