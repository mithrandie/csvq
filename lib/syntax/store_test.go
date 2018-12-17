package syntax

import (
	"reflect"
	"testing"
)

var testSyntax = []Expression{
	{
		Label: "SELECT Statement",
		Grammar: []Definition{
			{
				Name: "select_statement",
				Group: []Grammar{
					{Option{Link("with_clause")}, Link("select_query")},
				},
			},
			{
				Name: "select_query",
				Group: []Grammar{
					{Link("select_entity"), Option{Link("order_by_clause")}, Option{Link("limit_clause")}, Option{Link("offset_clause")}},
				},
			},
		},
		Children: []Expression{
			{
				Label: "WITH Clause",
				Grammar: []Definition{
					{
						Name: "with_clause",
						Group: []Grammar{
							{Keyword("WITH"), ContinuousOption{Link("common_table_expression")}},
						},
					},
					{
						Name: "common_table_expression",
						Group: []Grammar{
							{Option{Keyword("RECURSIVE")}, Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Keyword("AS"), Parentheses{Link("select_query")}},
						},
					},
				},
			},
			{
				Label: "SELECT Clause",
				Grammar: []Definition{
					{
						Name: "select_clause",
						Group: []Grammar{
							{Keyword("SELECT"), Option{Keyword("DISTINCT")}, ContinuousOption{Link("field")}},
						},
					},
					{
						Name: "field",
						Group: []Grammar{
							{Link("value")},
							{Link("value"), Keyword("AS"), Identifier("alias")},
						},
					},
				},
			},
		},
	},
	{
		Label: "INSERT Statement",
		Grammar: []Definition{
			{
				Name: "insert_statement",
				Group: []Grammar{
					{Option{Link("with_clause")}, Link("insert_query")},
				},
			},
			{
				Name: "insert_query",
				Group: []Grammar{
					{Keyword("INSERT"), Keyword("INTO"), Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Keyword("VALUES"), ContinuousOption{Link("row_value")}},
					{Keyword("INSERT"), Keyword("INTO"), Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Link("select_query")},
				},
			},
		},
	},
	{
		Label: "Operators",
		Children: []Expression{
			{
				Label: "Operator Precedence",
				Description: Description{
					Template: "The following table list operators from highest precedence to lowest.",
				},
			},
			{
				Label: "String Operators",
				Grammar: []Definition{
					{
						Name: "concatenation",
						Group: []Grammar{
							{Link("value"), Keyword("||"), Link("value")},
						},
					},
				},
			},
		},
	},
}

var storeSearchTests = []struct {
	Keys   []string
	Expect []Expression
}{
	{
		Keys: nil,
		Expect: []Expression{
			{
				Label: "SELECT Statement",
				Grammar: []Definition{
					{
						Name: "select_statement",
						Group: []Grammar{
							{Option{Link("with_clause")}, Link("select_query")},
						},
					},
					{
						Name: "select_query",
						Group: []Grammar{
							{Link("select_entity"), Option{Link("order_by_clause")}, Option{Link("limit_clause")}, Option{Link("offset_clause")}},
						},
					},
				},
				Children: []Expression{
					{
						Label: "WITH Clause",
						Grammar: []Definition{
							{
								Name: "with_clause",
								Group: []Grammar{
									{Keyword("WITH"), ContinuousOption{Link("common_table_expression")}},
								},
							},
							{
								Name: "common_table_expression",
								Group: []Grammar{
									{Option{Keyword("RECURSIVE")}, Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Keyword("AS"), Parentheses{Link("select_query")}},
								},
							},
						},
					},
					{
						Label: "SELECT Clause",
						Grammar: []Definition{
							{
								Name: "select_clause",
								Group: []Grammar{
									{Keyword("SELECT"), Option{Keyword("DISTINCT")}, ContinuousOption{Link("field")}},
								},
							},
							{
								Name: "field",
								Group: []Grammar{
									{Link("value")},
									{Link("value"), Keyword("AS"), Identifier("alias")},
								},
							},
						},
					},
				},
			},
			{
				Label: "INSERT Statement",
				Grammar: []Definition{
					{
						Name: "insert_statement",
						Group: []Grammar{
							{Option{Link("with_clause")}, Link("insert_query")},
						},
					},
					{
						Name: "insert_query",
						Group: []Grammar{
							{Keyword("INSERT"), Keyword("INTO"), Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Keyword("VALUES"), ContinuousOption{Link("row_value")}},
							{Keyword("INSERT"), Keyword("INTO"), Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Link("select_query")},
						},
					},
				},
			},
			{
				Label: "Operators",
				Children: []Expression{
					{
						Label: "Operator Precedence",
						Description: Description{
							Template: "The following table list operators from highest precedence to lowest.",
						},
					},
					{
						Label: "String Operators",
						Grammar: []Definition{
							{
								Name: "concatenation",
								Group: []Grammar{
									{Link("value"), Keyword("||"), Link("value")},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Keys: []string{"select"},
		Expect: []Expression{
			{
				Label: "SELECT Statement",
				Grammar: []Definition{
					{
						Name: "select_statement",
						Group: []Grammar{
							{Option{Link("with_clause")}, Link("select_query")},
						},
					},
					{
						Name: "select_query",
						Group: []Grammar{
							{Link("select_entity"), Option{Link("order_by_clause")}, Option{Link("limit_clause")}, Option{Link("offset_clause")}},
						},
					},
				},
				Children: []Expression{
					{
						Label: "WITH Clause",
						Grammar: []Definition{
							{
								Name: "with_clause",
								Group: []Grammar{
									{Keyword("WITH"), ContinuousOption{Link("common_table_expression")}},
								},
							},
							{
								Name: "common_table_expression",
								Group: []Grammar{
									{Option{Keyword("RECURSIVE")}, Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Keyword("AS"), Parentheses{Link("select_query")}},
								},
							},
						},
					},
					{
						Label: "SELECT Clause",
						Grammar: []Definition{
							{
								Name: "select_clause",
								Group: []Grammar{
									{Keyword("SELECT"), Option{Keyword("DISTINCT")}, ContinuousOption{Link("field")}},
								},
							},
							{
								Name: "field",
								Group: []Grammar{
									{Link("value")},
									{Link("value"), Keyword("AS"), Identifier("alias")},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Keys: []string{"clause"},
		Expect: []Expression{
			{
				Label: "WITH Clause",
				Grammar: []Definition{
					{
						Name: "with_clause",
						Group: []Grammar{
							{Keyword("WITH"), ContinuousOption{Link("common_table_expression")}},
						},
					},
					{
						Name: "common_table_expression",
						Group: []Grammar{
							{Option{Keyword("RECURSIVE")}, Identifier("table_name"), Option{Parentheses{ContinuousOption{Identifier("column_name")}}}, Keyword("AS"), Parentheses{Link("select_query")}},
						},
					},
				},
			},
			{
				Label: "SELECT Clause",
				Grammar: []Definition{
					{
						Name: "select_clause",
						Group: []Grammar{
							{Keyword("SELECT"), Option{Keyword("DISTINCT")}, ContinuousOption{Link("field")}},
						},
					},
					{
						Name: "field",
						Group: []Grammar{
							{Link("value")},
							{Link("value"), Keyword("AS"), Identifier("alias")},
						},
					},
				},
			},
		},
	},
	{
		Keys: []string{"field"},
		Expect: []Expression{
			{
				Label: "SELECT Clause",
				Grammar: []Definition{
					{
						Name: "field",
						Group: []Grammar{
							{Link("value")},
							{Link("value"), Keyword("AS"), Identifier("alias")},
						},
					},
				},
			},
		},
	},
	{
		Keys: []string{"operator prec"},
		Expect: []Expression{
			{
				Label: "Operator Precedence",
				Description: Description{
					Template: "The following table list operators from highest precedence to lowest.",
				},
			},
		},
	},
}

func TestStore_Search(t *testing.T) {
	store := NewStore()
	store.Syntax = testSyntax

	for _, v := range storeSearchTests {
		result := store.Search(v.Keys)
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %v", result, v.Expect, v.Keys)
		}
	}
}
