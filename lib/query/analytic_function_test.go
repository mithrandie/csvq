package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

type analyticFunctionTest struct {
	Name   string
	View   *View
	Args   []parser.Expression
	Clause parser.AnalyticClause
	Result *View
	Error  string
}

func testAnalyticFunction(t *testing.T, f func(*View, []parser.Expression, parser.AnalyticClause) error, tests []analyticFunctionTest) {
	for _, v := range tests {
		ViewCache.Clear()
		err := f(v.View, v.Args, v.Clause)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: result = %q, want %q", v.Name, v.View, v.Result)
		}
	}
}

var rowNumberTests = []analyticFunctionTest{
	{
		Name: "RowNumber",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
					parser.NewInteger(4),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
					parser.NewInteger(5),
				}),
			},
		},
	},
	{
		Name: "RowNumber with Partition",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
					parser.NewInteger(3),
				}),
			},
		},
	},
	{
		Name: "RowNumber Arguments Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Args: []parser.Expression{
			parser.NewInteger(1),
		},
		Error: "analytic function ROW_NUMBER takes no argument",
	},
	{
		Name: "RowNumber Partition Value Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestRowNumber(t *testing.T) {
	testAnalyticFunction(t, RowNumber, rowNumberTests)
}

var rankTests = []analyticFunctionTest{
	{
		Name: "Rank",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
			},
		},
	},
	{
		Name: "Rank Arguments Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Args: []parser.Expression{
			parser.NewInteger(1),
		},
		Error: "analytic function RANK takes no argument",
	},
	{
		Name: "Rank Partition Value Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Rank Order Value Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestRank(t *testing.T) {
	testAnalyticFunction(t, Rank, rankTests)
}

var denseRankTests = []analyticFunctionTest{
	{
		Name: "DenseRank",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
			},
		},
	},
	{
		Name: "DenseRank Arguments Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Args: []parser.Expression{
			parser.NewInteger(1),
		},
		Error: "analytic function DENSE_RANK takes no argument",
	},
	{
		Name: "DenseRank Partition Value Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "DenseRank Order Value Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
			},
		},
		Clause: parser.AnalyticClause{
			Partition: parser.Partition{
				Values: []parser.Expression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
			},
			OrderByClause: parser.OrderByClause{
				Items: []parser.Expression{
					parser.OrderItem{
						Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestDenseRank(t *testing.T) {
	testAnalyticFunction(t, DenseRank, denseRankTests)
}
