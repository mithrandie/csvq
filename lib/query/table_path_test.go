package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var parseTableNameTests = []struct {
	Table  parser.Table
	Result string
	Error  string
}{
	{
		Table: parser.Table{
			Object: parser.Identifier{Literal: "table.csv"},
			As:     parser.Token{Token: parser.AS, Literal: "as"},
			Alias:  parser.Identifier{Literal: "alias"},
		},
		Result: "alias",
	},
	{
		Table: parser.Table{
			Object: parser.Identifier{Literal: "/path/to/table.csv"},
		},
		Result: "table",
	},
	{
		Table: parser.Table{
			Object: parser.Stdin{},
		},
		Result: "STDIN",
	},
	{
		Table: parser.Table{
			Object: parser.Subquery{
				Query: parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Fields: []parser.QueryExpression{
								parser.NewIntegerValueFromString("1"),
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.QueryExpression{parser.Dual{}},
						},
					},
				},
			},
		},
		Result: "",
	},
	{
		Table: parser.Table{
			Object: parser.TableObject{
				Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
				FormatElement: parser.NewStringValue("[1, 2, 3]"),
				Path:          parser.Identifier{Literal: "fixed_length.dat", Quoted: true},
				Args:          nil,
			},
		},
		Result: "fixed_length",
	},
	{
		Table: parser.Table{
			Object: parser.TableFunction{
				Name: "file",
				Args: []parser.QueryExpression{parser.NewStringValue("table.csv")},
			},
		},
		Result: "table",
	},
	{
		Table: parser.Table{
			Object: parser.TableFunction{
				Name: "file",
			},
		},
		Error: "function FILE takes exactly 1 argument",
	},
}

func TestParseTableName(t *testing.T) {
	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range parseTableNameTests {
		name, err := ParseTableName(ctx, scope, v.Table)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Table.String(), err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Table.String(), err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Table.String(), v.Error)
			continue
		}

		if name.Literal != v.Result {
			t.Errorf("%s: \n result = %q,\n expect = %q", v.Table.String(), name.Literal, v.Result)
		}
	}
}

var convertTableFunctionTests = []struct {
	Name          string
	TableFunction parser.TableFunction
	Result        parser.QueryExpression
	Error         string
}{
	{
		Name: "Convert FILE function",
		TableFunction: parser.TableFunction{
			Name: "file",
			Args: []parser.QueryExpression{
				parser.NewStringValue("./table.csv"),
			},
		},
		Result: parser.Identifier{
			Literal: "./table.csv",
		},
	},
	{
		Name: "Argument Length Error in FILE function",
		TableFunction: parser.TableFunction{
			Name: "file",
			Args: []parser.QueryExpression{
				parser.NewStringValue("./table.csv"),
				parser.NewStringValue("./table.csv"),
			},
		},
		Error: "function FILE takes exactly 1 argument",
	},
	{
		Name: "First argument must be a string in FILE function",
		TableFunction: parser.TableFunction{
			Name: "file",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
			},
		},
		Error: "the first argument must be a string for function FILE",
	},
	{
		Name: "Convert INLINE function",
		TableFunction: parser.TableFunction{
			Name: "inline",
			Args: []parser.QueryExpression{
				parser.NewStringValue("./table.csv"),
			},
		},
		Result: parser.Identifier{
			Literal: "./table.csv",
		},
	},
	{
		Name: "Convert URL function",
		TableFunction: parser.TableFunction{
			Name: "url",
			Args: []parser.QueryExpression{
				parser.NewStringValue("https://example.com"),
			},
		},
		Result: parser.Url{
			Raw: "https://example.com",
		},
	},
	{
		Name: "First argument must be a string in URL function",
		TableFunction: parser.TableFunction{
			Name: "url",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
			},
		},
		Error: "the first argument must be a string for function URL",
	},
	{
		Name: "Convert DATA function",
		TableFunction: parser.TableFunction{
			Name: "data",
			Args: []parser.QueryExpression{
				parser.NewStringValue("1,a\n2,b"),
			},
		},
		Result: DataObject{
			Raw: "1,a\n2,b",
		},
	},
	{
		Name: "First argument must be a string in DATA function",
		TableFunction: parser.TableFunction{
			Name: "data",
			Args: []parser.QueryExpression{
				parser.NewNullValue(),
			},
		},
		Error: "the first argument must be a string for function DATA",
	},
	{
		Name: "Invalid Function Name Error",
		TableFunction: parser.TableFunction{
			Name: "invalid",
			Args: []parser.QueryExpression{
				parser.NewStringValue("1,a\n2,b"),
			},
		},
		Error: "function INVALID does not exist",
	},
}

func TestConvertTableFunction(t *testing.T) {
	ctx := context.Background()
	scope := NewReferenceScope(TestTx)

	for _, v := range convertTableFunctionTests {
		expr, err := ConvertTableFunction(ctx, scope, v.TableFunction)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		if !reflect.DeepEqual(expr, v.Result) {
			t.Errorf("%s: \n result = %v,\n expect = %v", v.Name, expr, v.Result)
		}
	}
}

var convertUrlExprTests = []struct {
	Name   string
	Url    parser.Url
	Result parser.QueryExpression
	Error  string
}{
	{
		Name: "Convert URL to HttpObject",
		Url: parser.Url{
			Raw: "https://example.com/おなかすいた/csv",
		},
		Result: HttpObject{
			URL: "https://example.com/%E3%81%8A%E3%81%AA%E3%81%8B%E3%81%99%E3%81%84%E3%81%9F/csv",
		},
	},
	{
		Name: "Convert Absolute File Path to Identifier",
		Url: parser.Url{
			Raw: "file:///home/my dir/data/foo.csv",
		},
		Result: parser.Identifier{
			Literal: "/home/my dir/data/foo.csv",
		},
	},
	{
		Name: "Convert Absolute File Path with URL Encoding to Identifier",
		Url: parser.Url{
			Raw: "file:///home/my%20dir/data/foo.csv",
		},
		Result: parser.Identifier{
			Literal: "/home/my dir/data/foo.csv",
		},
	},
	{
		Name: "Convert Relative File Path to Identifier",
		Url: parser.Url{
			Raw: "file:./data/foo.csv",
		},
		Result: parser.Identifier{
			Literal: "./data/foo.csv",
		},
	},
	{
		Name: "Unsupported URL Scheme",
		Url: parser.Url{
			Raw: "invalid:./data/foo.csv",
		},
		Error: "url scheme invalid is not supported",
	},
}

func TestConvertUrlExpr(t *testing.T) {
	for _, v := range convertUrlExprTests {
		expr, err := ConvertUrlExpr(v.Url)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}

		if !reflect.DeepEqual(expr, v.Result) {
			t.Errorf("%s: \n result = %v,\n expect = %v", v.Name, expr, v.Result)
		}
	}
}
