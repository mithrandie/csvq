package parser

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/value"
)

var parseTests = []struct {
	Input       string
	ForPrepared bool
	Output      []Statement
	SourceFile  string
	HolderNum   int
	Error       string
	ErrorLine   int
	ErrorChar   int
	ErrorFile   string
}{
	{
		Input: "select foo; select bar;",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "foo"}}}}},
			}},
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 13}, Select: "select", Fields: []QueryExpression{Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 20}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "bar"}}}}},
			}},
		},
	},
	{
		Input: "select 1 union all select 2 intersect select 3 except select 4",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectSet{
					LHS: SelectSet{
						LHS: SelectEntity{
							SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
						},
						Operator: Token{Token: UNION, Literal: "union", Line: 1, Char: 10},
						All:      Token{Token: ALL, Literal: "all", Line: 1, Char: 16},
						RHS: SelectSet{
							LHS: SelectEntity{
								SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 20}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("2")}}},
							},
							Operator: Token{Token: INTERSECT, Literal: "intersect", Line: 1, Char: 29},
							RHS: SelectEntity{
								SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 39}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("3")}}},
							},
						},
					},
					Operator: Token{Token: EXCEPT, Literal: "except", Line: 1, Char: 48},
					RHS: SelectEntity{
						SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 55}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("4")}}},
					},
				},
			},
		},
	},
	{
		Input: "select 1 union (select 2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectSet{
					LHS: SelectEntity{
						SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					},
					Operator: Token{Token: UNION, Literal: "union", Line: 1, Char: 10},
					RHS: Subquery{
						BaseExpr: &BaseExpr{line: 1, char: 16},
						Query: SelectQuery{
							SelectEntity: SelectEntity{
								SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 17}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("2")}}},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 as a from dual",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: NewIntegerValueFromString("1"),
								As:     "as",
								Alias:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "a"},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{Object: Dual{Dual: "dual"}},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from stdin",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{Object: Stdin{BaseExpr: &BaseExpr{line: 1, char: 16}, Stdin: "stdin"}},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from fixed('[1, 2, 3]', `fixed_length.dat`) fl",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: TableObject{
								BaseExpr:      &BaseExpr{line: 1, char: 16},
								Type:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "fixed"},
								FormatElement: NewStringValue("[1, 2, 3]"),
								Path:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 35}, Literal: "fixed_length.dat", Quoted: true},
							},
							Alias: Identifier{BaseExpr: &BaseExpr{line: 1, char: 55}, Literal: "fl"},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from csv(',', `table.csv`, 'utf8', null)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: TableObject{
								BaseExpr:      &BaseExpr{line: 1, char: 16},
								Type:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "csv"},
								FormatElement: NewStringValue(","),
								Path:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "table.csv", Quoted: true},
								Args:          []QueryExpression{NewStringValue("utf8"), NewNullValueFromString("null")},
							},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from json('{}', `table.txt`) as t",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: TableObject{
								BaseExpr:      &BaseExpr{line: 1, char: 16},
								Type:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "json"},
								FormatElement: NewStringValue("{}"),
								Path:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 27}, Literal: "table.txt", Quoted: true},
							},
							As:    "as",
							Alias: Identifier{BaseExpr: &BaseExpr{line: 1, char: 43}, Literal: "t"},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from ltsv(`table.ltsv`)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: TableObject{
								BaseExpr: &BaseExpr{line: 1, char: 16},
								Type:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "ltsv"},
								Path:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 21}, Literal: "table.ltsv", Quoted: true},
							},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from ltsv(`table.ltsv`, 'utf8')",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: TableObject{
								BaseExpr: &BaseExpr{line: 1, char: 16},
								Type:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "ltsv"},
								Path:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 21}, Literal: "table.ltsv", Quoted: true},
								Args:     []QueryExpression{NewStringValue("utf8")},
							},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from json_table('key', `table.json`)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: JsonQuery{
								BaseExpr:  &BaseExpr{line: 1, char: 16},
								JsonQuery: "json_table",
								Query:     NewStringValue("key"),
								JsonText:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "table.json", Quoted: true},
							},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from json_table('key', '{\"key2\":1}') jt",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: JsonQuery{
								BaseExpr:  &BaseExpr{line: 1, char: 16},
								JsonQuery: "json_table",
								Query:     NewStringValue("key"),
								JsonText:  NewStringValue("{\"key2\":1}"),
							},
							Alias: Identifier{BaseExpr: &BaseExpr{line: 1, char: 48}, Literal: "jt"},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select c1 from json_table('key', '{\"key2\":1}') as jt",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{
								Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{
						Table{
							Object: JsonQuery{
								BaseExpr:  &BaseExpr{line: 1, char: 16},
								JsonQuery: "json_table",
								Query:     NewStringValue("key"),
								JsonText:  NewStringValue("{\"key2\":1}"),
							},
							As:    "as",
							Alias: Identifier{BaseExpr: &BaseExpr{line: 1, char: 51}, Literal: "jt"},
						},
					}},
				},
			},
		},
	},
	{
		Input: "select 1 from table1, (select 2 from dual)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"},
							},
							Table{
								Object: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 23},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 24}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("2")}}},
											FromClause:   FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 alias, (select 2 from dual) alias2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"},
								Alias:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 22}, Literal: "alias"},
							},
							Table{
								Object: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 29},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 30}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("2")}}},
											FromClause:   FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
										},
									},
								},
								Alias: Identifier{BaseExpr: &BaseExpr{line: 1, char: 50}, Literal: "alias2"},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 as alias, (select 2 from dual) as alias2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"},
								As:     "as",
								Alias:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "alias"},
							},
							Table{
								Object: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 32},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 33}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("2")}}},
											FromClause:   FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
										},
									},
								},
								As:    "as",
								Alias: Identifier{BaseExpr: &BaseExpr{line: 1, char: 56}, Literal: "alias2"},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 \r\n" +
			" from dual \n" +
			" where 1 = 1 \n" +
			" group by column1, column2 \n" +
			" having 1 > 1 \n" +
			" order by column4, \n" +
			"          column5 desc, \n" +
			"          column6 asc, \n" +
			"          column7 nulls first, \n" +
			"          column8 desc nulls last, \n" +
			"          rank() over () \n" +
			" limit 10 \n" +
			" offset 10 \n",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause:   FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
					WhereClause: WhereClause{
						Where: "where",
						Filter: Comparison{
							LHS:      NewIntegerValueFromString("1"),
							Operator: "=",
							RHS:      NewIntegerValueFromString("1"),
						},
					},
					GroupByClause: GroupByClause{
						GroupBy: "group by",
						Items: []QueryExpression{
							FieldReference{BaseExpr: &BaseExpr{line: 4, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 4, char: 11}, Literal: "column1"}},
							FieldReference{BaseExpr: &BaseExpr{line: 4, char: 20}, Column: Identifier{BaseExpr: &BaseExpr{line: 4, char: 20}, Literal: "column2"}},
						},
					},
					HavingClause: HavingClause{
						Having: "having",
						Filter: Comparison{
							LHS:      NewIntegerValueFromString("1"),
							Operator: ">",
							RHS:      NewIntegerValueFromString("1"),
						},
					},
				},
				OrderByClause: OrderByClause{
					OrderBy: "order by",
					Items: []QueryExpression{
						OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 6, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 6, char: 11}, Literal: "column4"}}},
						OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 7, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 7, char: 11}, Literal: "column5"}}, Direction: Token{Token: DESC, Literal: "desc", Line: 7, Char: 19}},
						OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 8, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 8, char: 11}, Literal: "column6"}}, Direction: Token{Token: ASC, Literal: "asc", Line: 8, Char: 19}},
						OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 9, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 9, char: 11}, Literal: "column7"}}, Nulls: "nulls", Position: Token{Token: FIRST, Literal: "first", Line: 9, Char: 25}},
						OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 10, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 10, char: 11}, Literal: "column8"}}, Direction: Token{Token: DESC, Literal: "desc", Line: 10, Char: 19}, Nulls: "nulls", Position: Token{Token: LAST, Literal: "last", Line: 10, Char: 30}},
						OrderItem{Value: AnalyticFunction{
							BaseExpr: &BaseExpr{line: 11, char: 11},
							Name:     "rank",
							Over:     "over",
							AnalyticClause: AnalyticClause{
								PartitionClause: nil,
								OrderByClause:   nil,
							},
						}},
					},
				},
				LimitClause: LimitClause{
					BaseExpr: &BaseExpr{line: 12, char: 2},
					Limit:    "limit",
					Value:    NewIntegerValueFromString("10"),
				},
				OffsetClause: OffsetClause{
					BaseExpr: &BaseExpr{line: 13, char: 2},
					Offset:   "offset",
					Value:    NewIntegerValueFromString("10"),
				},
			},
		},
	},
	{
		Input: "select 1 \n" +
			" from dual \n" +
			" limit 10 percent",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause:   FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
				},
				LimitClause: LimitClause{
					BaseExpr: &BaseExpr{line: 3, char: 2},
					Limit:    "limit",
					Value:    NewIntegerValueFromString("10"),
					Percent:  "percent",
				},
			},
		},
	},
	{
		Input: "select 1 \n" +
			" from dual \n" +
			" limit 10 with ties",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause:   FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
				},
				LimitClause: LimitClause{
					BaseExpr: &BaseExpr{line: 3, char: 2},
					Limit:    "limit",
					Value:    NewIntegerValueFromString("10"),
					With:     LimitWith{With: "with", Type: Token{Token: TIES, Literal: "ties", Line: 3, Char: 16}},
				},
			},
		},
	},
	{
		Input: "select distinct * from dual",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 8},
						Fields: []QueryExpression{
							Field{Object: AllColumns{BaseExpr: &BaseExpr{line: 1, char: 17}}},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
				},
			},
		},
	},
	{
		Input: "with ct as (select 1) select * from ct",
		Output: []Statement{
			SelectQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []QueryExpression{
						InlineTable{
							Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 13},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("1")},
										},
									},
								},
							},
						},
					},
				},
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 23},
						Select:   "select",
						Fields:   []QueryExpression{Field{Object: AllColumns{BaseExpr: &BaseExpr{line: 1, char: 30}}}},
					},
					FromClause: FromClause{
						From:   "from",
						Tables: []QueryExpression{Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 37}, Literal: "ct"}}},
					},
				},
			},
		},
	},
	{
		Input: "with ct (column1) as (select 1) select * from ct",
		Output: []Statement{
			SelectQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []QueryExpression{
						InlineTable{
							Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "ct"},
							Fields: []QueryExpression{
								Identifier{BaseExpr: &BaseExpr{line: 1, char: 10}, Literal: "column1"},
							},
							As: "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 23},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("1")},
										},
									},
								},
							},
						},
					},
				},
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 33},
						Select:   "select",
						Fields:   []QueryExpression{Field{Object: AllColumns{BaseExpr: &BaseExpr{line: 1, char: 40}}}},
					},
					FromClause: FromClause{
						From:   "from",
						Tables: []QueryExpression{Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 47}, Literal: "ct"}}},
					},
				},
			},
		},
	},
	{
		Input: "with recursive ct as (select 1), ct2 as (select 2) select * from ct",
		Output: []Statement{
			SelectQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []QueryExpression{
						InlineTable{
							Name:      Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "ct"},
							Recursive: Token{Token: RECURSIVE, Literal: "recursive", Line: 1, Char: 6},
							As:        "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 23},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("1")},
										},
									},
								},
							},
						},
						InlineTable{
							Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "ct2"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 42},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("2")},
										},
									},
								},
							},
						},
					},
				},
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 52},
						Select:   "select",
						Fields:   []QueryExpression{Field{Object: AllColumns{BaseExpr: &BaseExpr{line: 1, char: 59}}}},
					},
					FromClause: FromClause{
						From:   "from",
						Tables: []QueryExpression{Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 66}, Literal: "ct"}}},
					},
				},
			},
		},
	},
	{
		Input: "select ident, tbl.3, 'foo', 1, 1.234, true, '2010-01-01 12:00:00', null, ('bar') from dual",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "ident"}}},
							Field{Object: ColumnNumber{BaseExpr: &BaseExpr{line: 1, char: 15}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "tbl"}, Number: value.NewInteger(3)}},
							Field{Object: NewStringValue("foo")},
							Field{Object: NewIntegerValueFromString("1")},
							Field{Object: NewFloatValueFromString("1.234")},
							Field{Object: NewTernaryValueFromString("true")},
							Field{Object: NewDatetimeValueFromString("2010-01-01 12:00:00", nil)},
							Field{Object: NewNullValueFromString("null")},
							Field{Object: Parentheses{Expr: NewStringValue("bar")}},
						},
					},
					FromClause: FromClause{From: "from", Tables: []QueryExpression{Table{Object: Dual{Dual: "dual"}}}},
				},
			},
		},
	},
	{
		Input: "select foo, \n" +
			" bar.foo, \n" +
			" stdin.foo, \n" +
			" bar.3, \n" +
			" stdin.3",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1},
					Select: "select",
					Fields: []QueryExpression{
						Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "foo"}}},
						Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 2, char: 2}, View: Identifier{BaseExpr: &BaseExpr{line: 2, char: 2}, Literal: "bar"}, Column: Identifier{BaseExpr: &BaseExpr{line: 2, char: 6}, Literal: "foo"}}},
						Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 3, char: 2}, View: Identifier{BaseExpr: &BaseExpr{line: 3, char: 2}, Literal: "stdin"}, Column: Identifier{BaseExpr: &BaseExpr{line: 3, char: 8}, Literal: "foo"}}},
						Field{Object: ColumnNumber{BaseExpr: &BaseExpr{line: 4, char: 2}, View: Identifier{BaseExpr: &BaseExpr{line: 4, char: 2}, Literal: "bar"}, Number: value.NewInteger(3)}},
						Field{Object: ColumnNumber{BaseExpr: &BaseExpr{line: 5, char: 2}, View: Identifier{BaseExpr: &BaseExpr{line: 5, char: 2}, Literal: "stdin"}, Number: value.NewInteger(3)}},
					},
				},
			}},
		},
	},
	{
		Input: "select ident || 'foo' || 'bar'",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Concat{Items: []QueryExpression{
								FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "ident"}},
								NewStringValue("foo"),
								NewStringValue("bar"),
							}}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 = 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Comparison{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: "=",
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) = (1, 2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Comparison{
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								RHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 29},
									Value: ValueList{
										Values: []QueryExpression{
											NewIntegerValueFromString("1"),
											NewIntegerValueFromString("2"),
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 < 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Comparison{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: "<",
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) < (select 1, 2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Comparison{
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Operator: "<",
								RHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 29},
									Value: Subquery{
										BaseExpr: &BaseExpr{line: 1, char: 29},
										Query: SelectQuery{
											SelectEntity: SelectEntity{
												SelectClause: SelectClause{
													BaseExpr: &BaseExpr{line: 1, char: 30},
													Select:   "select",
													Fields: []QueryExpression{
														Field{Object: NewIntegerValueFromString("1")},
														Field{Object: NewIntegerValueFromString("2")},
													},
												},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 is not null",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Is{
								Is:       "is",
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								RHS:      NewNullValueFromString("null"),
								Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 19},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 is true",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Is{
								Is:  "is",
								LHS: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								RHS: NewTernaryValueFromString("true"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not between -10 and +10 or column2 between 20 and 30",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS: Between{
									Between: "between",
									And:     "and",
									LHS:     FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
									Low: UnaryArithmetic{
										Operand:  NewIntegerValueFromString("10"),
										Operator: Token{Token: '-', Literal: "-", Line: 1, Char: 28},
									},
									High: UnaryArithmetic{
										Operand:  NewIntegerValueFromString("10"),
										Operator: Token{Token: '+', Literal: "+", Line: 1, Char: 36},
									},
									Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 16},
								},
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 40},
								RHS: Between{
									Between: "between",
									And:     "and",
									LHS:     FieldReference{BaseExpr: &BaseExpr{line: 1, char: 43}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 43}, Literal: "column2"}},
									Low:     NewIntegerValueFromString("20"),
									High:    NewIntegerValueFromString("30"),
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) not between (1, 2) and (3, 4) and (column3, column4) between (5, 6) and (7, 8)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS: Between{
									Between: "between",
									And:     "and",
									LHS: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 8},
										Value: ValueList{
											Values: []QueryExpression{
												FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
												FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
											},
										},
									},
									Low: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 39},
										Value: ValueList{
											Values: []QueryExpression{
												NewIntegerValueFromString("1"),
												NewIntegerValueFromString("2"),
											},
										},
									},
									High: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 50},
										Value: ValueList{
											Values: []QueryExpression{
												NewIntegerValueFromString("3"),
												NewIntegerValueFromString("4"),
											},
										},
									},
									Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 27},
								},
								RHS: Between{
									Between: "between",
									And:     "and",
									LHS: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 61},
										Value: ValueList{
											Values: []QueryExpression{
												FieldReference{BaseExpr: &BaseExpr{line: 1, char: 62}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 62}, Literal: "column3"}},
												FieldReference{BaseExpr: &BaseExpr{line: 1, char: 71}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 71}, Literal: "column4"}},
											},
										},
									},
									Low: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 88},
										Value: ValueList{
											Values: []QueryExpression{
												NewIntegerValueFromString("5"),
												NewIntegerValueFromString("6"),
											},
										},
									},
									High: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 99},
										Value: ValueList{
											Values: []QueryExpression{
												NewIntegerValueFromString("7"),
												NewIntegerValueFromString("8"),
											},
										},
									},
								},
								Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 57},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not in (1, 2, 3) and column2 in (4, 5, 6)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS: In{
									In:  "in",
									LHS: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
									Values: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 23},
										Value: ValueList{
											Values: []QueryExpression{
												NewIntegerValueFromString("1"),
												NewIntegerValueFromString("2"),
												NewIntegerValueFromString("3"),
											},
										},
									},
									Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 16},
								},
								RHS: In{
									In:  "in",
									LHS: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 37}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 37}, Literal: "column2"}},
									Values: RowValue{
										BaseExpr: &BaseExpr{line: 1, char: 48},
										Value: ValueList{
											Values: []QueryExpression{
												NewIntegerValueFromString("4"),
												NewIntegerValueFromString("5"),
												NewIntegerValueFromString("6"),
											},
										},
									},
								},
								Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 33},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 in json_row('key', '{\"key\":1}')",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: In{
								In:  "in",
								LHS: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Values: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 19},
									Value: JsonQuery{
										JsonQuery: "json_row",
										Query:     NewStringValue("key"),
										JsonText:  NewStringValue("{\"key\":1}"),
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) not in ((1, 2), (3, 4))",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: In{
								In: "in",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Values: RowValueList{
									RowValues: []QueryExpression{
										RowValue{
											BaseExpr: &BaseExpr{line: 1, char: 35},
											Value: ValueList{
												Values: []QueryExpression{
													NewIntegerValueFromString("1"),
													NewIntegerValueFromString("2"),
												},
											},
										},
										RowValue{
											BaseExpr: &BaseExpr{line: 1, char: 43},
											Value: ValueList{
												Values: []QueryExpression{
													NewIntegerValueFromString("3"),
													NewIntegerValueFromString("4"),
												},
											},
										},
									},
								},
								Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 27},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) in (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: In{
								In: "in",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Values: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 30},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 31}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) in json_row('key', '{\"key\":1}')",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: In{
								In: "in",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Values: JsonQuery{
									BaseExpr:  &BaseExpr{line: 1, char: 30},
									JsonQuery: "json_row",
									Query:     NewStringValue("key"),
									JsonText:  NewStringValue("{\"key\":1}"),
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not like 'pattern1' and column2 like 'pattern2'",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS: Like{
									Like:     "like",
									LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
									Pattern:  NewStringValue("pattern1"),
									Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 16},
								},
								Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 36},
								RHS: Like{
									Like:    "like",
									LHS:     FieldReference{BaseExpr: &BaseExpr{line: 1, char: 40}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 40}, Literal: "column2"}},
									Pattern: NewStringValue("pattern2"),
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 like 'pattern1' or column2 not like 'pattern2'",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS: Like{
									Like:    "like",
									LHS:     FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
									Pattern: NewStringValue("pattern1"),
								},
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 32},
								RHS: Like{
									Like:     "like",
									LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 35}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 35}, Literal: "column2"}},
									Pattern:  NewStringValue("pattern2"),
									Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 43},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 = any (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Any{
								Any:      "any",
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: "=",
								Values: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 22},
									Value: Subquery{
										BaseExpr: &BaseExpr{line: 1, char: 22},
										Query: SelectQuery{
											SelectEntity: SelectEntity{
												SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 23}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) = any ((1, 2), (3, 4))",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Any{
								Any: "any",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: RowValueList{
									RowValues: []QueryExpression{
										RowValue{
											BaseExpr: &BaseExpr{line: 1, char: 34},
											Value: ValueList{
												Values: []QueryExpression{
													NewIntegerValueFromString("1"),
													NewIntegerValueFromString("2"),
												},
											},
										},
										RowValue{
											BaseExpr: &BaseExpr{line: 1, char: 42},
											Value: ValueList{
												Values: []QueryExpression{
													NewIntegerValueFromString("3"),
													NewIntegerValueFromString("4"),
												},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) = any (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Any{
								Any: "any",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 33},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 34}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 = all (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: All{
								All:      "all",
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: "=",
								Values: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 22},
									Value: Subquery{
										BaseExpr: &BaseExpr{line: 1, char: 22},
										Query: SelectQuery{
											SelectEntity: SelectEntity{
												SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 23}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) = all ((1, 2), (3, 4))",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: All{
								All: "all",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: RowValueList{
									RowValues: []QueryExpression{
										RowValue{
											BaseExpr: &BaseExpr{line: 1, char: 34},
											Value: ValueList{
												Values: []QueryExpression{
													NewIntegerValueFromString("1"),
													NewIntegerValueFromString("2"),
												},
											},
										},
										RowValue{
											BaseExpr: &BaseExpr{line: 1, char: 42},
											Value: ValueList{
												Values: []QueryExpression{
													NewIntegerValueFromString("3"),
													NewIntegerValueFromString("4"),
												},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) = all (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: All{
								All: "all",
								LHS: RowValue{
									BaseExpr: &BaseExpr{line: 1, char: 8},
									Value: ValueList{
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 9}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "column1"}},
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 18}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 33},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 34}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select exists (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Exists{
								Exists: "exists",
								Query: Subquery{
									BaseExpr: &BaseExpr{line: 1, char: 15},
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 16}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 + 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: int('+'),
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 - 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: int('-'),
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 * 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: int('*'),
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 / 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: int('/'),
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 % 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "column1"}},
								Operator: int('%'),
								RHS:      NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select true and false",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS:      NewTernaryValueFromString("true"),
								Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 13},
								RHS:      NewTernaryValueFromString("false"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select true or false",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS:      NewTernaryValueFromString("true"),
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 13},
								RHS:      NewTernaryValueFromString("false"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select not false",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: UnaryLogic{
								Operator: Token{Token: NOT, Literal: "not", Line: 1, Char: 8},
								Operand:  NewTernaryValueFromString("false"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select true or (false and false)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS:      NewTernaryValueFromString("true"),
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 13},
								RHS: Parentheses{
									Expr: Logic{
										LHS:      NewTernaryValueFromString("false"),
										Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 23},
										RHS:      NewTernaryValueFromString("false"),
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select true and true or !false and not false",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Logic{
								LHS: Logic{
									LHS:      NewTernaryValueFromString("true"),
									Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 13},
									RHS:      NewTernaryValueFromString("true"),
								},
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 22},
								RHS: Logic{
									LHS: UnaryLogic{
										Operator: Token{Token: '!', Literal: "!", Line: 1, Char: 25},
										Operand:  NewTernaryValueFromString("false"),
									},
									Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 32},
									RHS: UnaryLogic{
										Operator: Token{Token: NOT, Literal: "not", Line: 1, Char: 36},
										Operand:  NewTernaryValueFromString("false"),
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select @var",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Variable{BaseExpr: &BaseExpr{line: 1, char: 8}, Name: "var"}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select @var := 1",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: VariableSubstitution{
								Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 8}, Name: "var"},
								Value:    NewIntegerValueFromString("1"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select case when true then 'A' when false then 'B' end",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: CaseExpr{
								Case: "case",
								End:  "end",
								When: []QueryExpression{
									CaseExprWhen{
										When:      "when",
										Then:      "then",
										Condition: NewTernaryValueFromString("true"),
										Result:    NewStringValue("A"),
									},
									CaseExprWhen{
										When:      "when",
										Then:      "then",
										Condition: NewTernaryValueFromString("false"),
										Result:    NewStringValue("B"),
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select case column1 when 1 then 'A' when 2 then 'B' else 'C' end",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: CaseExpr{
								Case:  "case",
								End:   "end",
								Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 13}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "column1"}},
								When: []QueryExpression{
									CaseExprWhen{
										When:      "when",
										Then:      "then",
										Condition: NewIntegerValueFromString("1"),
										Result:    NewStringValue("A"),
									},
									CaseExprWhen{
										When:      "when",
										Then:      "then",
										Condition: NewIntegerValueFromString("2"),
										Result:    NewStringValue("B"),
									},
								},
								Else: CaseExprElse{
									Else:   "else",
									Result: NewStringValue("C"),
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select now()",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Function{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "now",
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select trim(column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Function{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "trim",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 13}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "column1"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select trim(column1, column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Function{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "trim",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 13}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "column1"}},
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 22}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 22}, Literal: "column2"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select json_object()",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Function{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "json_object",
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select json_object(column1, column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Function{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "json_object",
								Args: []QueryExpression{
									Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 20}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "column1"}}},
									Field{Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 29}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 29}, Literal: "column2"}}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select if(column1, column2, column3)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: Function{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "if",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 11}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 11}, Literal: "column1"}},
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 20}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "column2"}},
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 29}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 29}, Literal: "column3"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select aggfunc(distinct column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AggregateFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "aggfunc",
								Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 16},
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 25}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "column1"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count(*)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AggregateFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "count",
								Args: []QueryExpression{
									AllColumns{BaseExpr: &BaseExpr{line: 1, char: 14}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count(distinct *)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AggregateFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "count",
								Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 14},
								Args: []QueryExpression{
									AllColumns{BaseExpr: &BaseExpr{line: 1, char: 23}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count(column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AggregateFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "count",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 14}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "column1"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count(distinct column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AggregateFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "count",
								Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 14},
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 23}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 23}, Literal: "column1"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select listagg(column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: ListFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "listagg",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 16}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "column1"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select listagg(distinct column1, ',')",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: ListFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "listagg",
								Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 16},
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 25}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "column1"}},
									NewStringValue(","),
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select listagg(distinct column1) within group (order by column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: ListFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "listagg",
								Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 16},
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 25}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "column1"}},
								},
								WithinGroup: "within group",
								OrderBy: OrderByClause{
									OrderBy: "order by",
									Items: []QueryExpression{
										OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 57}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 57}, Literal: "column1"}}},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select listagg(column1, ',') within group (order by column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: ListFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "listagg",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 16}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "column1"}},
									NewStringValue(","),
								},
								WithinGroup: "within group",
								OrderBy: OrderByClause{
									OrderBy: "order by",
									Items: []QueryExpression{
										OrderItem{Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 53}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 53}, Literal: "column1"}}},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select cursor cur is not open",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: CursorStatus{
								CursorLit: "cursor",
								Cursor:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "cur"},
								Is:        "is",
								Negation:  Token{Token: NOT, Literal: "not", Line: 1, Char: 22},
								Type:      OPEN,
								TypeLit:   "open",
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select cursor cur is not in range",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: CursorStatus{
								CursorLit: "cursor",
								Cursor:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "cur"},
								Is:        "is",
								Negation:  Token{Token: NOT, Literal: "not", Line: 1, Char: 22},
								Type:      RANGE,
								TypeLit:   "in range",
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select cursor cur count",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: CursorAttrebute{
								CursorLit: "cursor",
								Cursor:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "cur"},
								Attrebute: Token{Token: COUNT, Literal: "count", Line: 1, Char: 19},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 38}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 38}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 55}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 55}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (order by column2 rows current row)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (order by column2 rows unbounded preceding)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: PRECEDING,
											Unbounded: true,
											Literal:   "unbounded preceding",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (order by column2 rows 1 preceding)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: PRECEDING,
											Offset:    1,
											Literal:   "1 preceding",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (order by column2 rows between unbounded preceding and 1 following)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: PRECEDING,
											Unbounded: true,
											Literal:   "unbounded preceding",
										},
										FrameHigh: WindowFramePosition{
											Direction: FOLLOWING,
											Offset:    1,
											Literal:   "1 following",
										},
										Between: "between",
										And:     "and",
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (order by column2 rows between 1 preceding and unbounded following)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: PRECEDING,
											Offset:    1,
											Literal:   "1 preceding",
										},
										FrameHigh: WindowFramePosition{
											Direction: FOLLOWING,
											Unbounded: true,
											Literal:   "unbounded following",
										},
										Between: "between",
										And:     "and",
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select userfunc() over (order by column2 rows between current row and unbounded following)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "userfunc",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
										FrameHigh: WindowFramePosition{
											Direction: FOLLOWING,
											Unbounded: true,
											Literal:   "unbounded following",
										},
										Between: "between",
										And:     "and",
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select f(column1) over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "f",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 10}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 10}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 38}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 38}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 55}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 55}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select f(distinct column1) over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "f",
								Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 10},
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 19}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 19}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 47}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 47}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 64}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 64}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select min(column1) over (partition by column1 order by column2 rows current row)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "min",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 12}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 12}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 40}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 40}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 57}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 57}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count(column1) over (partition by column1 order by column2 rows current row)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "count",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 14}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 42}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 42}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 59}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 59}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count(*) over (partition by column1 order by column2 rows current row)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "count",
								Args: []QueryExpression{
									AllColumns{BaseExpr: &BaseExpr{line: 1, char: 14}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 36}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 36}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 53}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 53}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select listagg(column1) over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "listagg",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 16}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 44}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 44}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 61}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 61}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select listagg(column1, ',') over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "listagg",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 16}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "column1"}},
									NewStringValue(","),
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 49}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 49}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 66}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 66}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select rank() over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "rank",
								Over:     "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 34}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 51}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 51}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select first_value(column1) over (partition by column1 order by column2 rows current row)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "first_value",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 20}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 48}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 48}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 65}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 65}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select first_value(column1) ignore nulls over (partition by column1 order by column2 rows current row)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "first_value",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 20}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "column1"}},
								},
								IgnoreNulls:    true,
								IgnoreNullsLit: "ignore nulls",
								Over:           "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 61}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 61}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 78}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 78}, Literal: "column2"}},
											},
										},
									},
									WindowingClause: WindowingClause{
										Rows: "rows",
										FrameLow: WindowFramePosition{
											Direction: CURRENT,
											Literal:   "current row",
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select lag(column1) over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "lag",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 12}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 12}, Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 40}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 40}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 57}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 57}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select lag(column1) ignore nulls over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						BaseExpr: &BaseExpr{line: 1, char: 1},
						Select:   "select",
						Fields: []QueryExpression{
							Field{Object: AnalyticFunction{
								BaseExpr: &BaseExpr{line: 1, char: 8},
								Name:     "lag",
								Args: []QueryExpression{
									FieldReference{BaseExpr: &BaseExpr{line: 1, char: 12}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 12}, Literal: "column1"}},
								},
								IgnoreNulls:    true,
								IgnoreNullsLit: "ignore nulls",
								Over:           "over",
								AnalyticClause: AnalyticClause{
									PartitionClause: PartitionClause{
										PartitionBy: "partition by",
										Values: []QueryExpression{
											FieldReference{BaseExpr: &BaseExpr{line: 1, char: 53}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 53}, Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []QueryExpression{
											OrderItem{
												Value: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 70}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 70}, Literal: "column2"}},
											},
										},
									},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 cross join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 33}, Literal: "table2"}},
									JoinType:  Token{Token: CROSS, Literal: "cross", Line: 1, Char: 22},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 cross join table2 cross join table3",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join: "join",
									Table: Table{
										Object: Join{
											Join:      "join",
											Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
											JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 33}, Literal: "table2"}},
											JoinType:  Token{Token: CROSS, Literal: "cross", Line: 1, Char: 22},
										},
									},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 51}, Literal: "table3"}},
									JoinType:  Token{Token: CROSS, Literal: "cross", Line: 1, Char: 40},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 join table2 on table1.id = table2.id inner join table3 on table1.id = table3.id",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join: "join",
									Table: Table{
										Object: Join{
											Join:      "join",
											Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
											JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 27}, Literal: "table2"}},
											Condition: JoinCondition{
												Literal: "on",
												On: Comparison{
													LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 37}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 37}, Literal: "table1"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 44}, Literal: "id"}},
													Operator: "=",
													RHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 49}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 49}, Literal: "table2"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 56}, Literal: "id"}},
												},
											},
										},
									},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 70}, Literal: "table3"}},
									Condition: JoinCondition{
										Literal: "on",
										On: Comparison{
											LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 80}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 80}, Literal: "table1"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 87}, Literal: "id"}},
											Operator: "=",
											RHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 92}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 92}, Literal: "table3"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 99}, Literal: "id"}},
										},
									},
									JoinType: Token{Token: INNER, Literal: "inner", Line: 1, Char: 59},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 inner join table2 on table1.id = table2.id",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 33}, Literal: "table2"}},
									Condition: JoinCondition{
										Literal: "on",
										On: Comparison{
											LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 43}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 43}, Literal: "table1"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 50}, Literal: "id"}},
											Operator: "=",
											RHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 55}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 55}, Literal: "table2"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 62}, Literal: "id"}},
										},
									},
									JoinType: Token{Token: INNER, Literal: "inner", Line: 1, Char: 22},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 natural join table2 natural join table3",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join: "join",
									Table: Table{
										Object: Join{
											Join:      "join",
											Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
											JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 35}, Literal: "table2"}},
											Natural:   Token{Token: NATURAL, Literal: "natural", Line: 1, Char: 22},
										},
									},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 55}, Literal: "table3"}},
									Natural:   Token{Token: NATURAL, Literal: "natural", Line: 1, Char: 42},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 left join table2 using(id) left join table3 using(id)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join: "join",
									Table: Table{
										Object: Join{
											Join:      "join",
											Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
											JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 32}, Literal: "table2"}},
											Direction: Token{Token: LEFT, Literal: "left", Line: 1, Char: 22},
											Condition: JoinCondition{
												Literal: "using",
												Using: []QueryExpression{
													Identifier{BaseExpr: &BaseExpr{line: 1, char: 45}, Literal: "id"},
												},
											},
										},
									},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 59}, Literal: "table3"}},
									Direction: Token{Token: LEFT, Literal: "left", Line: 1, Char: 49},
									Condition: JoinCondition{
										Literal: "using",
										Using: []QueryExpression{
											Identifier{BaseExpr: &BaseExpr{line: 1, char: 72}, Literal: "id"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 right outer join table2 using(id)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 39}, Literal: "table2"}},
									Direction: Token{Token: RIGHT, Literal: "right", Line: 1, Char: 22},
									JoinType:  Token{Token: OUTER, Literal: "outer", Line: 1, Char: 28},
									Condition: JoinCondition{
										Literal: "using",
										Using: []QueryExpression{
											Identifier{BaseExpr: &BaseExpr{line: 1, char: 52}, Literal: "id"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 natural right join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 41}, Literal: "table2"}},
									Natural:   Token{Token: NATURAL, Literal: "natural", Line: 1, Char: 22},
									Direction: Token{Token: RIGHT, Literal: "right", Line: 1, Char: 30},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 full join table2 on table1.id = table2.id full join table3 on table3.id = table1.id",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join: "join",
									Table: Table{
										Object: Join{
											Join:      "join",
											Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
											JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 32}, Literal: "table2"}},
											Direction: Token{Token: FULL, Literal: "full", Line: 1, Char: 22},
											Condition: JoinCondition{
												Literal: "on",
												On: Comparison{
													LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 42}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 42}, Literal: "table1"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 49}, Literal: "id"}},
													Operator: "=",
													RHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 54}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 54}, Literal: "table2"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 61}, Literal: "id"}},
												},
											},
										},
									},
									JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 74}, Literal: "table3"}},
									Direction: Token{Token: FULL, Literal: "full", Line: 1, Char: 64},
									Condition: JoinCondition{
										Literal: "on",
										On: Comparison{
											LHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 84}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 84}, Literal: "table3"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 91}, Literal: "id"}},
											Operator: "=",
											RHS:      FieldReference{BaseExpr: &BaseExpr{line: 1, char: 96}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 96}, Literal: "table1"}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 103}, Literal: "id"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 cross join (table2 cross join table3)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1}, Select: "select", Fields: []QueryExpression{Field{Object: NewIntegerValueFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []QueryExpression{
							Table{
								Object: Join{
									Join:  "join",
									Table: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 15}, Literal: "table1"}},
									JoinTable: Parentheses{Expr: Table{
										Object: Join{
											Join:      "join",
											Table:     Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "table2"}},
											JoinTable: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 52}, Literal: "table3"}},
											JoinType:  Token{Token: CROSS, Literal: "cross", Line: 1, Char: 41},
										},
									}},
									JoinType: Token{Token: CROSS, Literal: "cross", Line: 1, Char: 22},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "var @var1, @var2 := 2; @var1 := 1;",
		Output: []Statement{
			VariableDeclaration{
				Assignments: []VariableAssignment{
					{
						Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 5}, Name: "var1"},
					},
					{
						Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 12}, Name: "var2"},
						Value:    NewIntegerValueFromString("2"),
					},
				},
			},
			VariableSubstitution{
				Variable: Variable{
					BaseExpr: &BaseExpr{line: 1, char: 24},
					Name:     "var1",
				},
				Value: NewIntegerValueFromString("1"),
			},
		},
	},
	{
		Input: "declare @var1 := 1",
		Output: []Statement{
			VariableDeclaration{
				Assignments: []VariableAssignment{
					{
						Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 9}, Name: "var1"},
						Value:    NewIntegerValueFromString("1"),
					},
				},
			},
		},
	},
	{
		Input: "dispose @var1",
		Output: []Statement{
			DisposeVariable{
				Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 9}, Name: "var1"},
			},
		},
	},
	{
		Input: "set @%var = ident",
		Output: []Statement{
			SetEnvVar{
				EnvVar: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 5}, Name: "var"},
				Value:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "ident"},
			},
		},
	},
	{
		Input: "set @%var = 1",
		Output: []Statement{
			SetEnvVar{
				EnvVar: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 5}, Name: "var"},
				Value:  NewIntegerValueFromString("1"),
			},
		},
	},
	{
		Input: "set @%var to ident",
		Output: []Statement{
			SetEnvVar{
				EnvVar: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 5}, Name: "var"},
				Value:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "ident"},
			},
		},
	},
	{
		Input: "set @%var to 1",
		Output: []Statement{
			SetEnvVar{
				EnvVar: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 5}, Name: "var"},
				Value:  NewIntegerValueFromString("1"),
			},
		},
	},
	{
		Input: "unset @%var",
		Output: []Statement{
			UnsetEnvVar{
				EnvVar: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 7}, Name: "var"},
			},
		},
	},
	{
		Input: "func('arg1', 'arg2')",
		Output: []Statement{
			Function{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "func",
				Args: []QueryExpression{
					NewStringValue("arg1"),
					NewStringValue("arg2"),
				},
			},
		},
	},
	{
		Input: "with ct as (select 1) insert into table1 values (1, 'str1'), (2, 'str2')",
		Output: []Statement{
			InsertQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []QueryExpression{
						InlineTable{
							Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 13},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("1")},
										},
									},
								},
							},
						},
					},
				},
				Table: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 35}, Literal: "table1"}},
				ValuesList: []QueryExpression{
					RowValue{
						BaseExpr: &BaseExpr{line: 1, char: 49},
						Value: ValueList{
							Values: []QueryExpression{
								NewIntegerValueFromString("1"),
								NewStringValue("str1"),
							},
						},
					},
					RowValue{
						BaseExpr: &BaseExpr{line: 1, char: 62},
						Value: ValueList{
							Values: []QueryExpression{
								NewIntegerValueFromString("2"),
								NewStringValue("str2"),
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "insert into table1 (column1, column2, table1.3) values (1, 'str1'), (2, 'str2')",
		Output: []Statement{
			InsertQuery{
				Table: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"}},
				Fields: []QueryExpression{
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 21}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 21}, Literal: "column1"}},
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 30}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 30}, Literal: "column2"}},
					ColumnNumber{BaseExpr: &BaseExpr{line: 1, char: 39}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 39}, Literal: "table1"}, Number: value.NewInteger(3)},
				},
				ValuesList: []QueryExpression{
					RowValue{
						BaseExpr: &BaseExpr{line: 1, char: 56},
						Value: ValueList{
							Values: []QueryExpression{
								NewIntegerValueFromString("1"),
								NewStringValue("str1"),
							},
						},
					},
					RowValue{
						BaseExpr: &BaseExpr{line: 1, char: 69},
						Value: ValueList{
							Values: []QueryExpression{
								NewIntegerValueFromString("2"),
								NewStringValue("str2"),
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "insert into table1 select 1, 2",
		Output: []Statement{
			InsertQuery{
				Table: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"}},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 20},
							Select:   "select",
							Fields: []QueryExpression{
								Field{Object: NewIntegerValueFromString("1")},
								Field{Object: NewIntegerValueFromString("2")},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "insert into table1 (column1, column2) select 1, 2",
		Output: []Statement{
			InsertQuery{
				Table: Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"}},
				Fields: []QueryExpression{
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 21}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 21}, Literal: "column1"}},
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 30}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 30}, Literal: "column2"}},
				},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 39},
							Select:   "select",
							Fields: []QueryExpression{
								Field{Object: NewIntegerValueFromString("1")},
								Field{Object: NewIntegerValueFromString("2")},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "with ct as (select 1) update table1 set column1 = 1, column2 = 2, table1.3 = 3 from table1 where true",
		Output: []Statement{
			UpdateQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []QueryExpression{
						InlineTable{
							Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 13},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("1")},
										},
									},
								},
							},
						},
					},
				},
				Tables: []QueryExpression{
					Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 30}, Literal: "table1"}},
				},
				SetList: []UpdateSet{
					{Field: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 41}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 41}, Literal: "column1"}}, Value: NewIntegerValueFromString("1")},
					{Field: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 54}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 54}, Literal: "column2"}}, Value: NewIntegerValueFromString("2")},
					{Field: ColumnNumber{BaseExpr: &BaseExpr{line: 1, char: 67}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 67}, Literal: "table1"}, Number: value.NewInteger(3)}, Value: NewIntegerValueFromString("3")},
				},
				FromClause: FromClause{
					From: "from",
					Tables: []QueryExpression{
						Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 85}, Literal: "table1"}},
					},
				},
				WhereClause: WhereClause{
					Where:  "where",
					Filter: NewTernaryValueFromString("true"),
				},
			},
		},
	},
	{
		Input: "update csv(',', table1) set column1 = 1, column2 = 2, table1.3 = 3 where true",
		Output: []Statement{
			UpdateQuery{
				Tables: []QueryExpression{
					Table{Object: TableObject{
						BaseExpr:      &BaseExpr{line: 1, char: 8},
						Type:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "csv"},
						FormatElement: NewStringValue(","),
						Path:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 17}, Literal: "table1"},
					}},
				},
				SetList: []UpdateSet{
					{Field: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 29}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 29}, Literal: "column1"}}, Value: NewIntegerValueFromString("1")},
					{Field: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 42}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 42}, Literal: "column2"}}, Value: NewIntegerValueFromString("2")},
					{Field: ColumnNumber{BaseExpr: &BaseExpr{line: 1, char: 55}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 55}, Literal: "table1"}, Number: value.NewInteger(3)}, Value: NewIntegerValueFromString("3")},
				},
				WhereClause: WhereClause{
					Where:  "where",
					Filter: NewTernaryValueFromString("true"),
				},
			},
		},
	},
	{
		Input: "with ct as (select 1) delete from table1",
		Output: []Statement{
			DeleteQuery{
				BaseExpr: &BaseExpr{line: 1, char: 23},
				WithClause: WithClause{
					With: "with",
					InlineTables: []QueryExpression{
						InlineTable{
							Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										BaseExpr: &BaseExpr{line: 1, char: 13},
										Select:   "select",
										Fields: []QueryExpression{
											Field{Object: NewIntegerValueFromString("1")},
										},
									},
								},
							},
						},
					},
				},
				FromClause: FromClause{
					From: "from",
					Tables: []QueryExpression{
						Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 35}, Literal: "table1"}},
					},
				},
			},
		},
	},
	{
		Input: "delete table1 from table1 where true",
		Output: []Statement{
			DeleteQuery{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Tables: []QueryExpression{
					Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "table1"}},
				},
				FromClause: FromClause{
					From: "from",
					Tables: []QueryExpression{
						Table{Object: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "table1"}},
					},
				},
				WhereClause: WhereClause{
					Where:  "where",
					Filter: NewTernaryValueFromString("true"),
				},
			},
		},
	},
	{
		Input: "create table newtable (column1, column2)",
		Output: []Statement{
			CreateTable{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "newtable"},
				Fields: []QueryExpression{
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 33}, Literal: "column2"},
				},
			},
		},
	},
	{
		Input: "create table newtable (column1, column2) select 1, 2",
		Output: []Statement{
			CreateTable{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "newtable"},
				Fields: []QueryExpression{
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 33}, Literal: "column2"},
				},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 42},
							Select:   "select",
							Fields: []QueryExpression{
								Field{
									Object: NewIntegerValueFromString("1"),
								},
								Field{
									Object: NewIntegerValueFromString("2"),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "create table newtable select 1, 2",
		Output: []Statement{
			CreateTable{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "newtable"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 23},
							Select:   "select",
							Fields: []QueryExpression{
								Field{
									Object: NewIntegerValueFromString("1"),
								},
								Field{
									Object: NewIntegerValueFromString("2"),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "create table newtable (column1, column2) as select 1, 2",
		Output: []Statement{
			CreateTable{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "newtable"},
				Fields: []QueryExpression{
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 33}, Literal: "column2"},
				},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 45},
							Select:   "select",
							Fields: []QueryExpression{
								Field{
									Object: NewIntegerValueFromString("1"),
								},
								Field{
									Object: NewIntegerValueFromString("2"),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "create table newtable as select 1, 2",
		Output: []Statement{
			CreateTable{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "newtable"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 26},
							Select:   "select",
							Fields: []QueryExpression{
								Field{
									Object: NewIntegerValueFromString("1"),
								},
								Field{
									Object: NewIntegerValueFromString("2"),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "alter table table1 add column1",
		Output: []Statement{
			AddColumns{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []ColumnDefault{
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					},
				},
			},
		},
	},
	{
		Input: "alter table csv(',', table1) add column1",
		Output: []Statement{
			AddColumns{
				Table: TableObject{
					BaseExpr:      &BaseExpr{line: 1, char: 13},
					Type:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "csv"},
					FormatElement: NewStringValue(","),
					Path:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 22}, Literal: "table1"},
				},
				Columns: []ColumnDefault{
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column1"},
					},
				},
			},
		},
	},
	{
		Input: "alter table table1 add (column1, column2 default 1) first",
		Output: []Statement{
			AddColumns{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []ColumnDefault{
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "column1"},
					},
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "column2"},
						Value:  NewIntegerValueFromString("1"),
					},
				},
				Position: ColumnPosition{
					Position: Token{Token: FIRST, Literal: "first", Line: 1, Char: 53},
				},
			},
		},
	},
	{
		Input: "alter table table1 add column1 last",
		Output: []Statement{
			AddColumns{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []ColumnDefault{
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					},
				},
				Position: ColumnPosition{
					Position: Token{Token: LAST, Literal: "last", Line: 1, Char: 32},
				},
			},
		},
	},
	{
		Input: "alter table table1 add column1 after column2",
		Output: []Statement{
			AddColumns{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []ColumnDefault{
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					},
				},
				Position: ColumnPosition{
					Position: Token{Token: AFTER, Literal: "after", Line: 1, Char: 32},
					Column:   FieldReference{BaseExpr: &BaseExpr{line: 1, char: 38}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 38}, Literal: "column2"}},
				},
			},
		},
	},
	{
		Input: "alter table table1 add column1 before column2",
		Output: []Statement{
			AddColumns{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []ColumnDefault{
					{
						Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "column1"},
					},
				},
				Position: ColumnPosition{
					Position: Token{Token: BEFORE, Literal: "before", Line: 1, Char: 32},
					Column:   FieldReference{BaseExpr: &BaseExpr{line: 1, char: 39}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 39}, Literal: "column2"}},
				},
			},
		},
	},
	{
		Input: "alter table table1 drop column1",
		Output: []Statement{
			DropColumns{
				Table:   Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []QueryExpression{FieldReference{BaseExpr: &BaseExpr{line: 1, char: 25}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 25}, Literal: "column1"}}},
			},
		},
	},
	{
		Input: "alter table table1 drop (column1, column2, table1.3)",
		Output: []Statement{
			DropColumns{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Columns: []QueryExpression{
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 26}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 26}, Literal: "column1"}},
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 35}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 35}, Literal: "column2"}},
					ColumnNumber{BaseExpr: &BaseExpr{line: 1, char: 44}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 44}, Literal: "table1"}, Number: value.NewInteger(3)},
				},
			},
		},
	},
	{
		Input: "alter table table1 rename column1 to column2",
		Output: []Statement{
			RenameColumn{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Old:   FieldReference{BaseExpr: &BaseExpr{line: 1, char: 27}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 27}, Literal: "column1"}},
				New:   Identifier{BaseExpr: &BaseExpr{line: 1, char: 38}, Literal: "column2"},
			},
		},
	},
	{
		Input: "alter table table1 rename table1.3 to column2",
		Output: []Statement{
			RenameColumn{
				Table: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Old:   ColumnNumber{BaseExpr: &BaseExpr{line: 1, char: 27}, View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 27}, Literal: "table1"}, Number: value.NewInteger(3)},
				New:   Identifier{BaseExpr: &BaseExpr{line: 1, char: 39}, Literal: "column2"},
			},
		},
	},
	{
		Input: "alter table table1 set format to 'json'",
		Output: []Statement{
			SetTableAttribute{
				BaseExpr:  &BaseExpr{line: 1, char: 1},
				Table:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Attribute: Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "format"},
				Value:     NewStringValue("json"),
			},
		},
	},
	{
		Input: "alter table table1 set format to json",
		Output: []Statement{
			SetTableAttribute{
				BaseExpr:  &BaseExpr{line: 1, char: 1},
				Table:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "table1"},
				Attribute: Identifier{BaseExpr: &BaseExpr{line: 1, char: 24}, Literal: "format"},
				Value:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 34}, Literal: "json"},
			},
		},
	},
	{
		Input: "commit",
		Output: []Statement{
			TransactionControl{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Token:    COMMIT,
			},
		},
	},
	{
		Input: "rollback",
		Output: []Statement{
			TransactionControl{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Token:    ROLLBACK,
			},
		},
	},
	{
		Input: "echo 'foo'",
		Output: []Statement{
			Echo{
				Value: NewStringValue("foo"),
			},
		},
	},
	{
		Input: "print 'foo'",
		Output: []Statement{
			Print{
				Value: NewStringValue("foo"),
			},
		},
	},
	{
		Input: "printf 'foo'",
		Output: []Statement{
			Printf{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Format:   NewStringValue("foo"),
			},
		},
	},
	{
		Input: "printf 'foo', 'bar'",
		Output: []Statement{
			Printf{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Format:   NewStringValue("foo"),
				Values: []QueryExpression{
					NewStringValue("bar"),
				},
			},
		},
	},
	{
		Input: "printf 'foo' using 'bar'",
		Output: []Statement{
			Printf{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Format:   NewStringValue("foo"),
				Values: []QueryExpression{
					NewStringValue("bar"),
				},
			},
		},
	},
	{
		Input: "source `/path/to/file.sql`",
		Output: []Statement{
			Source{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				FilePath: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "/path/to/file.sql", Quoted: true},
			},
		},
	},
	{
		Input: "source '/path/to/file.sql'",
		Output: []Statement{
			Source{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				FilePath: NewStringValue("/path/to/file.sql"),
			},
		},
	},
	{
		Input: "execute 'select 1'",
		Output: []Statement{
			Execute{
				BaseExpr:   &BaseExpr{line: 1, char: 1},
				Statements: NewStringValue("select 1"),
			},
		},
	},
	{
		Input: "execute 'select %s' using 'a'",
		Output: []Statement{
			Execute{
				BaseExpr:   &BaseExpr{line: 1, char: 1},
				Statements: NewStringValue("select %s"),
				Values: []QueryExpression{
					NewStringValue("a"),
				},
			},
		},
	},
	{
		Input: "syntax",
		Output: []Statement{
			Syntax{
				BaseExpr: &BaseExpr{line: 1, char: 1},
			},
		},
	},
	{
		Input: "syntax foo",
		Output: []Statement{
			Syntax{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Keywords: []QueryExpression{
					FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "foo"}},
				},
			},
		},
	},
	{
		Input: "chdir `dirpath`",
		Output: []Statement{
			Chdir{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				DirPath:  Identifier{BaseExpr: &BaseExpr{line: 1, char: 7}, Literal: "dirpath", Quoted: true},
			},
		},
	},
	{
		Input: "chdir 'dirpath'",
		Output: []Statement{
			Chdir{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				DirPath:  NewStringValue("dirpath"),
			},
		},
	},
	{
		Input: "pwd",
		Output: []Statement{
			Pwd{
				BaseExpr: &BaseExpr{line: 1, char: 1},
			},
		},
	},
	{
		Input: "reload config",
		Output: []Statement{
			Reload{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Type:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "config"},
			},
		},
	},
	{
		Input: "set @@delimiter = ','",
		Output: []Statement{
			SetFlag{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "delimiter",
				Value:    NewStringValue(","),
			},
		},
	},
	{
		Input: "set @@encoding = sjis",
		Output: []Statement{
			SetFlag{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "encoding",
				Value:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "sjis"},
			},
		},
	},
	{
		Input: "set @@delimiter to ','",
		Output: []Statement{
			SetFlag{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "delimiter",
				Value:    NewStringValue(","),
			},
		},
	},
	{
		Input: "set @@encoding to sjis",
		Output: []Statement{
			SetFlag{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "encoding",
				Value:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 19}, Literal: "sjis"},
			},
		},
	},
	{
		Input: "add '%Y%m%d' to @@datetime_format",
		Output: []Statement{
			AddFlagElement{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "datetime_format",
				Value:    NewStringValue("%Y%m%d"),
			},
		},
	},
	{
		Input: "remove '%Y%m%d' from @@datetime_format",
		Output: []Statement{
			RemoveFlagElement{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "datetime_format",
				Value:    NewStringValue("%Y%m%d"),
			},
		},
	},
	{
		Input: "show @@delimiter",
		Output: []Statement{
			ShowFlag{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     "delimiter",
			},
		},
	},
	{
		Input: "show tables",
		Output: []Statement{
			ShowObjects{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Type:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "tables"},
			},
		},
	},
	{
		Input: "show fields from table1",
		Output: []Statement{
			ShowFields{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Type:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "fields"},
				Table:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "table1"},
			},
		},
	},
	{
		Input: "show fields from csv(',', table1)",
		Output: []Statement{
			ShowFields{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Type:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "fields"},
				Table: TableObject{
					BaseExpr:      &BaseExpr{line: 1, char: 18},
					Type:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "csv"},
					FormatElement: NewStringValue(","),
					Path:          Identifier{BaseExpr: &BaseExpr{line: 1, char: 27}, Literal: "table1"},
				},
			},
		},
	},
	{
		Input: "trigger error",
		Output: []Statement{
			Trigger{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Event:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "error"},
			},
		},
	},
	{
		Input: "trigger error 'user error'",
		Output: []Statement{
			Trigger{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Event:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "error"},
				Message:  NewStringValue("user error"),
			},
		},
	},
	{
		Input: "trigger error 300 'user error'",
		Output: []Statement{
			Trigger{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Event:    Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "error"},
				Message:  NewStringValue("user error"),
				Code:     value.NewInteger(300),
			},
		},
	},
	{
		Input: "declare cur cursor for select 1",
		Output: []Statement{
			CursorDeclaration{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "cur"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 24},
							Select:   "select",
							Fields: []QueryExpression{
								Field{Object: NewIntegerValueFromString("1")},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "open cur",
		Output: []Statement{
			OpenCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 6}, Literal: "cur"},
			},
		},
	},
	{
		Input: "close cur",
		Output: []Statement{
			CloseCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 7}, Literal: "cur"},
			},
		},
	},
	{
		Input: "dispose cursor cur",
		Output: []Statement{
			DisposeCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "cur"},
			},
		},
	},
	{
		Input: "fetch cur into @var1, @var2",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 7}, Literal: "cur"},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 16}, Name: "var1"},
					{BaseExpr: &BaseExpr{line: 1, char: 23}, Name: "var2"},
				},
			},
		},
	},
	{
		Input: "fetch next cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 12}, Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: NEXT, Literal: "next", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 21}, Name: "var1"},
				},
			},
		},
	},
	{
		Input: "fetch prior cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: PRIOR, Literal: "prior", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 22}, Name: "var1"},
				},
			},
		},
	},
	{
		Input: "fetch first cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 13}, Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: FIRST, Literal: "first", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 22}, Name: "var1"},
				},
			},
		},
	},
	{
		Input: "fetch last cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 12}, Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: LAST, Literal: "last", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 21}, Name: "var1"},
				},
			},
		},
	},
	{
		Input: "fetch absolute 1 cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "cur"},
				Position: FetchPosition{
					BaseExpr: &BaseExpr{line: 1, char: 7},
					Position: Token{Token: ABSOLUTE, Literal: "absolute", Line: 1, Char: 7},
					Number:   NewIntegerValueFromString("1"),
				},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 27}, Name: "var1"},
				},
			},
		},
	},
	{
		Input: "fetch relative 1 cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "cur"},
				Position: FetchPosition{
					BaseExpr: &BaseExpr{line: 1, char: 7},
					Position: Token{Token: RELATIVE, Literal: "relative", Line: 1, Char: 7},
					Number:   NewIntegerValueFromString("1"),
				},
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 27}, Name: "var1"},
				},
			},
		},
	},
	{
		Input: "declare tbl view (column1, column2)",
		Output: []Statement{
			ViewDeclaration{
				View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "tbl"},
				Fields: []QueryExpression{
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 19}, Literal: "column1"},
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 28}, Literal: "column2"},
				},
			},
		},
	},
	{
		Input: "declare tbl view (column1, column2) as select 1, 2",
		Output: []Statement{
			ViewDeclaration{
				View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "tbl"},
				Fields: []QueryExpression{
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 19}, Literal: "column1"},
					Identifier{BaseExpr: &BaseExpr{line: 1, char: 28}, Literal: "column2"},
				},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 40},
							Select:   "select",
							Fields: []QueryExpression{
								Field{
									Object: NewIntegerValueFromString("1"),
								},
								Field{
									Object: NewIntegerValueFromString("2"),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "declare tbl view as select 1, 2",
		Output: []Statement{
			ViewDeclaration{
				View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "tbl"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							BaseExpr: &BaseExpr{line: 1, char: 21},
							Select:   "select",
							Fields: []QueryExpression{
								Field{
									Object: NewIntegerValueFromString("1"),
								},
								Field{
									Object: NewIntegerValueFromString("2"),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "dispose view tbl",
		Output: []Statement{
			DisposeView{
				View: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "tbl"},
			},
		},
	},
	{
		Input: "prepare stmt from 'select :val'",
		Output: []Statement{
			StatementPreparation{
				Name:      Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "stmt"},
				Statement: value.NewString("select :val"),
			},
		},
	},
	{
		Input: "execute stmt",
		Output: []Statement{
			ExecuteStatement{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "stmt"},
			},
		},
	},
	{
		Input: "execute stmt using 'a', 1 as val",
		Output: []Statement{
			ExecuteStatement{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Name:     Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "stmt"},
				Values: []ReplaceValue{
					{Value: NewStringValue("a")},
					{Value: NewIntegerValueFromString("1"), Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 30}, Literal: "val"}},
				},
			},
		},
	},
	{
		Input: "dispose prepare stmt",
		Output: []Statement{
			DisposeStatement{
				Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 17}, Literal: "stmt"},
			},
		},
	},
	{
		Input: "if @var1 = 1 then print 1; end if",
		Output: []Statement{
			If{
				Condition: Comparison{
					LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 4}, Name: "var1"},
					RHS:      NewIntegerValueFromString("1"),
					Operator: "=",
				},
				Statements: []Statement{
					Print{Value: NewIntegerValueFromString("1")},
				},
			},
		},
	},
	{
		Input: "if @var1 = 1 then print 1; elseif @var1 = 2 then print 2; elseif @var1 = 3 then print 3; else print 4; end if",
		Output: []Statement{
			If{
				Condition: Comparison{
					LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 4}, Name: "var1"},
					RHS:      NewIntegerValueFromString("1"),
					Operator: "=",
				},
				Statements: []Statement{
					Print{Value: NewIntegerValueFromString("1")},
				},
				ElseIf: []ElseIf{
					{
						Condition: Comparison{
							LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 35}, Name: "var1"},
							RHS:      NewIntegerValueFromString("2"),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewIntegerValueFromString("2")},
						},
					},
					{
						Condition: Comparison{
							LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 66}, Name: "var1"},
							RHS:      NewIntegerValueFromString("3"),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewIntegerValueFromString("3")},
						},
					},
				},
				Else: Else{
					Statements: []Statement{
						Print{Value: NewIntegerValueFromString("4")},
					},
				},
			},
		},
	},
	{
		Input: "while @var1 do print @var1; end while",
		Output: []Statement{
			While{
				Condition: Variable{BaseExpr: &BaseExpr{line: 1, char: 7}, Name: "var1"},
				Statements: []Statement{
					Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 22}, Name: "var1"}},
				},
			},
		},
	},
	{
		Input: "while @var1 in cur do print @var1; end while",
		Output: []Statement{
			WhileInCursor{
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 7}, Name: "var1"},
				},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 16}, Literal: "cur"},
				Statements: []Statement{
					Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 29}, Name: "var1"}},
				},
			},
		},
	},
	{
		Input: "while @var1, @var2 in cur do print @var1; end while",
		Output: []Statement{
			WhileInCursor{
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 7}, Name: "var1"},
					{BaseExpr: &BaseExpr{line: 1, char: 14}, Name: "var2"},
				},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 23}, Literal: "cur"},
				Statements: []Statement{
					Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 36}, Name: "var1"}},
				},
			},
		},
	},
	{
		Input: "while var @var1 in cur do print @var1; end while",
		Output: []Statement{
			WhileInCursor{
				WithDeclaration: true,
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 11}, Name: "var1"},
				},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 20}, Literal: "cur"},
				Statements: []Statement{
					Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 33}, Name: "var1"}},
				},
			},
		},
	},
	{
		Input: "while declare @var1, @var2 in cur do print @var1; end while",
		Output: []Statement{
			WhileInCursor{
				WithDeclaration: true,
				Variables: []Variable{
					{BaseExpr: &BaseExpr{line: 1, char: 15}, Name: "var1"},
					{BaseExpr: &BaseExpr{line: 1, char: 22}, Name: "var2"},
				},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 31}, Literal: "cur"},
				Statements: []Statement{
					Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 44}, Name: "var1"}},
				},
			},
		},
	},
	{
		Input: "case when true then print @var1; when false then print @var2; end case",
		Output: []Statement{
			Case{
				When: []CaseWhen{
					{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 27}, Name: "var1"}},
						},
					},
					{
						Condition: NewTernaryValueFromString("false"),
						Statements: []Statement{
							Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 56}, Name: "var2"}},
						},
					},
				},
			},
		},
	},
	{
		Input: "case when true then print @var1; when false then print @var2; else print @var3; end case",
		Output: []Statement{
			Case{
				When: []CaseWhen{
					{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 27}, Name: "var1"}},
						},
					},
					{
						Condition: NewTernaryValueFromString("false"),
						Statements: []Statement{
							Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 56}, Name: "var2"}},
						},
					},
				},
				Else: CaseElse{
					Statements: []Statement{
						Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 74}, Name: "var3"}},
					},
				},
			},
		},
	},
	{
		Input: "exit",
		Output: []Statement{
			Exit{},
		},
	},
	{
		Input: "exit 1",
		Output: []Statement{
			Exit{Code: value.NewIntegerFromString("1")},
		},
	},
	{
		Input: "while true do print @var1; continue; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 21}, Name: "var1"}},
					FlowControl{Token: CONTINUE},
				},
			},
		},
	},
	{
		Input: "while true do break; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					FlowControl{Token: BREAK},
				},
			},
		},
	},
	{
		Input: "while true do exit; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					Exit{},
				},
			},
		},
	},
	{
		Input: "while true do if @var1 = 1 then continue; end if; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					If{
						Condition: Comparison{
							LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 18}, Name: "var1"},
							RHS:      NewIntegerValueFromString("1"),
							Operator: "=",
						},
						Statements: []Statement{
							FlowControl{Token: CONTINUE},
						},
					},
				},
			},
		},
	},
	{
		Input: "while true do if @var1 = 1 then continue; elseif @var1 = 2 then break; elseif @var1 = 3 then exit; else continue; end if; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					If{
						Condition: Comparison{
							LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 18}, Name: "var1"},
							RHS:      NewIntegerValueFromString("1"),
							Operator: "=",
						},
						Statements: []Statement{
							FlowControl{Token: CONTINUE},
						},
						ElseIf: []ElseIf{
							{
								Condition: Comparison{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 50}, Name: "var1"},
									RHS:      NewIntegerValueFromString("2"),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: BREAK},
								},
							},
							{
								Condition: Comparison{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 79}, Name: "var1"},
									RHS:      NewIntegerValueFromString("3"),
									Operator: "=",
								},
								Statements: []Statement{
									Exit{},
								},
							},
						},
						Else: Else{
							Statements: []Statement{
								FlowControl{Token: CONTINUE},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "while true do case when true then print @var1; when false then continue; end case; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					Case{
						When: []CaseWhen{
							{
								Condition: NewTernaryValueFromString("true"),
								Statements: []Statement{
									Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 41}, Name: "var1"}},
								},
							},
							{
								Condition: NewTernaryValueFromString("false"),
								Statements: []Statement{
									FlowControl{Token: CONTINUE},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "while true do case when true then print @var1; when false then exit; else continue; end case; end while",
		Output: []Statement{
			While{
				Condition: NewTernaryValueFromString("true"),
				Statements: []Statement{
					Case{
						When: []CaseWhen{
							{
								Condition: NewTernaryValueFromString("true"),
								Statements: []Statement{
									Print{Value: Variable{BaseExpr: &BaseExpr{line: 1, char: 41}, Name: "var1"}},
								},
							},
							{
								Condition: NewTernaryValueFromString("false"),
								Statements: []Statement{
									Exit{},
								},
							},
						},
						Else: CaseElse{
							Statements: []Statement{
								FlowControl{Token: CONTINUE},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "declare func1 function () as begin end",
		Output: []Statement{
			FunctionDeclaration{
				Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "func1"},
			},
		},
	},
	{
		Input: "declare func1 function (@arg1 default 0, @arg2 default 1) as begin end",
		Output: []Statement{
			FunctionDeclaration{
				Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "func1"},
				Parameters: []VariableAssignment{
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 25}, Name: "arg1"}, Value: NewIntegerValueFromString("0")},
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 42}, Name: "arg2"}, Value: NewIntegerValueFromString("1")},
				},
			},
		},
	},
	{
		Input: "declare func1 function (@arg1, @arg2 default 0) as begin \n" +
			"if @var1 = 1 then print 1; end if; \n" +
			"if @var1 = 1 then print 1; elseif @var1 = 2 then print 2; elseif @var1 = 3 then print 3; else print 4; end if; \n" +
			"while true do break; end while; \n" +
			"while true do if @var1 = 1 then continue; end if; end while; \n" +
			"while true do if @var1 = 1 then continue; elseif @var1 = 2 then break; elseif @var1 = 3 then return; else continue; end if; end while; \n" +
			"while @var1 in cur do print @var1; end while; \n" +
			"while @var1, @var2 in cur do print @var1; end while; \n" +
			"case when true then print @var1; when false then print @var2; end case; \n" +
			"case when true then print @var1; when false then return; else return; end case; \n" +
			"while true do case when true then print @var1; when false then continue; end case; end while; \n" +
			"while true do case when true then print @var1; when false then return; else continue; end case; end while; \n" +
			"return; \n" +
			"return @var1; \n" +
			"end",
		Output: []Statement{
			FunctionDeclaration{
				Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "func1"},
				Parameters: []VariableAssignment{
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 25}, Name: "arg1"}},
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 32}, Name: "arg2"}, Value: NewIntegerValueFromString("0")},
				},
				Statements: []Statement{
					If{
						Condition: Comparison{
							LHS:      Variable{BaseExpr: &BaseExpr{line: 2, char: 4}, Name: "var1"},
							RHS:      NewIntegerValueFromString("1"),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewIntegerValueFromString("1")},
						},
					},
					If{
						Condition: Comparison{
							LHS:      Variable{BaseExpr: &BaseExpr{line: 3, char: 4}, Name: "var1"},
							RHS:      NewIntegerValueFromString("1"),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewIntegerValueFromString("1")},
						},
						ElseIf: []ElseIf{
							{
								Condition: Comparison{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 3, char: 35}, Name: "var1"},
									RHS:      NewIntegerValueFromString("2"),
									Operator: "=",
								},
								Statements: []Statement{
									Print{Value: NewIntegerValueFromString("2")},
								},
							},
							{
								Condition: Comparison{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 3, char: 66}, Name: "var1"},
									RHS:      NewIntegerValueFromString("3"),
									Operator: "=",
								},
								Statements: []Statement{
									Print{Value: NewIntegerValueFromString("3")},
								},
							},
						},
						Else: Else{
							Statements: []Statement{
								Print{Value: NewIntegerValueFromString("4")},
							},
						},
					},
					While{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							FlowControl{Token: BREAK},
						},
					},
					While{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							If{
								Condition: Comparison{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 5, char: 18}, Name: "var1"},
									RHS:      NewIntegerValueFromString("1"),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: CONTINUE},
								},
							},
						},
					},
					While{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							If{
								Condition: Comparison{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 6, char: 18}, Name: "var1"},
									RHS:      NewIntegerValueFromString("1"),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: CONTINUE},
								},
								ElseIf: []ElseIf{
									{
										Condition: Comparison{
											LHS:      Variable{BaseExpr: &BaseExpr{line: 6, char: 50}, Name: "var1"},
											RHS:      NewIntegerValueFromString("2"),
											Operator: "=",
										},
										Statements: []Statement{
											FlowControl{Token: BREAK},
										},
									},
									{
										Condition: Comparison{
											LHS:      Variable{BaseExpr: &BaseExpr{line: 6, char: 79}, Name: "var1"},
											RHS:      NewIntegerValueFromString("3"),
											Operator: "=",
										},
										Statements: []Statement{
											Return{Value: NewNullValue()},
										},
									},
								},
								Else: Else{
									Statements: []Statement{
										FlowControl{Token: CONTINUE},
									},
								},
							},
						},
					},
					WhileInCursor{
						Variables: []Variable{
							{BaseExpr: &BaseExpr{line: 7, char: 7}, Name: "var1"},
						},
						Cursor: Identifier{BaseExpr: &BaseExpr{line: 7, char: 16}, Literal: "cur"},
						Statements: []Statement{
							Print{Value: Variable{BaseExpr: &BaseExpr{line: 7, char: 29}, Name: "var1"}},
						},
					},
					WhileInCursor{
						Variables: []Variable{
							{BaseExpr: &BaseExpr{line: 8, char: 7}, Name: "var1"},
							{BaseExpr: &BaseExpr{line: 8, char: 14}, Name: "var2"},
						},
						Cursor: Identifier{BaseExpr: &BaseExpr{line: 8, char: 23}, Literal: "cur"},
						Statements: []Statement{
							Print{Value: Variable{BaseExpr: &BaseExpr{line: 8, char: 36}, Name: "var1"}},
						},
					},
					Case{
						When: []CaseWhen{
							{
								Condition: NewTernaryValueFromString("true"),
								Statements: []Statement{
									Print{Value: Variable{BaseExpr: &BaseExpr{line: 9, char: 27}, Name: "var1"}},
								},
							},
							{
								Condition: NewTernaryValueFromString("false"),
								Statements: []Statement{
									Print{Value: Variable{BaseExpr: &BaseExpr{line: 9, char: 56}, Name: "var2"}},
								},
							},
						},
					},
					Case{
						When: []CaseWhen{
							{
								Condition: NewTernaryValueFromString("true"),
								Statements: []Statement{
									Print{Value: Variable{BaseExpr: &BaseExpr{line: 10, char: 27}, Name: "var1"}},
								},
							},
							{
								Condition: NewTernaryValueFromString("false"),
								Statements: []Statement{
									Return{Value: NewNullValue()},
								},
							},
						},
						Else: CaseElse{
							Statements: []Statement{
								Return{Value: NewNullValue()},
							},
						},
					},
					While{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							Case{
								When: []CaseWhen{
									{
										Condition: NewTernaryValueFromString("true"),
										Statements: []Statement{
											Print{Value: Variable{BaseExpr: &BaseExpr{line: 11, char: 41}, Name: "var1"}},
										},
									},
									{
										Condition: NewTernaryValueFromString("false"),
										Statements: []Statement{
											FlowControl{Token: CONTINUE},
										},
									},
								},
							},
						},
					},
					While{
						Condition: NewTernaryValueFromString("true"),
						Statements: []Statement{
							Case{
								When: []CaseWhen{
									{
										Condition: NewTernaryValueFromString("true"),
										Statements: []Statement{
											Print{Value: Variable{BaseExpr: &BaseExpr{line: 12, char: 41}, Name: "var1"}},
										},
									},
									{
										Condition: NewTernaryValueFromString("false"),
										Statements: []Statement{
											Return{Value: NewNullValue()},
										},
									},
								},
								Else: CaseElse{
									Statements: []Statement{
										FlowControl{Token: CONTINUE},
									},
								},
							},
						},
					},
					Return{
						Value: NewNullValue(),
					},
					Return{
						Value: Variable{BaseExpr: &BaseExpr{line: 14, char: 8}, Name: "var1"},
					},
				},
			},
		},
	},
	{
		Input: "declare aggfunc aggregate (cur) as begin end",
		Output: []Statement{
			AggregateDeclaration{
				Name:   Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "aggfunc"},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 28}, Literal: "cur"},
			},
		},
	},
	{
		Input: "declare aggfunc aggregate (cur, @var1) as begin end",
		Output: []Statement{
			AggregateDeclaration{
				Name:   Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "aggfunc"},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 28}, Literal: "cur"},
				Parameters: []VariableAssignment{
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 33}, Name: "var1"}},
				},
			},
		},
	},
	{
		Input: "declare aggfunc aggregate (cur, @var1, @var2) as begin end",
		Output: []Statement{
			AggregateDeclaration{
				Name:   Identifier{BaseExpr: &BaseExpr{line: 1, char: 9}, Literal: "aggfunc"},
				Cursor: Identifier{BaseExpr: &BaseExpr{line: 1, char: 28}, Literal: "cur"},
				Parameters: []VariableAssignment{
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 33}, Name: "var1"}},
					{Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 40}, Name: "var2"}},
				},
			},
		},
	},
	{
		Input: "dispose function userfunc",
		Output: []Statement{
			DisposeFunction{
				Name: Identifier{BaseExpr: &BaseExpr{line: 1, char: 18}, Literal: "userfunc"},
			},
		},
	},
	{
		Input: "select @var1 := @var2 + @var3",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: VariableSubstitution{
								Variable: Variable{BaseExpr: &BaseExpr{line: 1, char: 8}, Name: "var1"},
								Value: Arithmetic{
									LHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 17}, Name: "var2"},
									Operator: int('+'),
									RHS:      Variable{BaseExpr: &BaseExpr{line: 1, char: 25}, Name: "var3"},
								},
							},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select @%var",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 8}, Name: "var"},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select @%`var`",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: EnvironmentVariable{BaseExpr: &BaseExpr{line: 1, char: 8}, Name: "var", Quoted: true},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select @#var",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: RuntimeInformation{BaseExpr: &BaseExpr{line: 1, char: 8}, Name: "var"},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select ties",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "ties"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select nulls",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "nulls"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select tables",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "tables"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select views",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "views"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select cursors",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "cursors"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select functions",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "functions"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select rows",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "rows"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select csv",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "csv"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select json",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "json"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select fixed",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "fixed"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select ltsv",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "ltsv"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "select fields",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "fields"}},
						},
					},
				},
			}},
		},
	},
	{
		Input: "'abc'",
		Output: []Statement{
			NewStringValue("abc"),
		},
	},
	{
		Input: "(if(column1, column2, column3))",
		Output: []Statement{
			Parentheses{
				Expr: Function{
					BaseExpr: &BaseExpr{line: 1, char: 2},
					Name:     "if",
					Args: []QueryExpression{
						FieldReference{BaseExpr: &BaseExpr{line: 1, char: 5}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 5}, Literal: "column1"}},
						FieldReference{BaseExpr: &BaseExpr{line: 1, char: 14}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 14}, Literal: "column2"}},
						FieldReference{BaseExpr: &BaseExpr{line: 1, char: 23}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 23}, Literal: "column3"}},
					},
				},
			},
		},
	},
	{
		Input: "select c1;\n$echo foo;",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{
					BaseExpr: &BaseExpr{line: 1, char: 1},
					Select:   "select",
					Fields: []QueryExpression{
						Field{
							Object: FieldReference{BaseExpr: &BaseExpr{line: 1, char: 8}, Column: Identifier{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "c1"}},
						},
					},
				},
			}},
			ExternalCommand{
				BaseExpr: &BaseExpr{line: 2, char: 1},
				Command:  "echo foo",
			},
		},
	},
	{
		Input: "$",
		Output: []Statement{
			ExternalCommand{
				BaseExpr: &BaseExpr{line: 1, char: 1},
				Command:  "",
			},
		},
	},
	{
		Input:       "select ?, :val",
		ForPrepared: true,
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{BaseExpr: &BaseExpr{line: 1, char: 1},
					Select: "select",
					Fields: []QueryExpression{
						Field{Object: Placeholder{BaseExpr: &BaseExpr{line: 1, char: 8}, Literal: "?", Ordinal: 1}},
						Field{Object: Placeholder{BaseExpr: &BaseExpr{line: 1, char: 11}, Literal: ":val", Ordinal: 2, Name: "val"}},
					},
				},
			}},
		},
		HolderNum: 2,
	},
	{
		Input:     "select 1 = 1 = 1",
		Error:     "syntax error: unexpected token \"=\"",
		ErrorLine: 1,
		ErrorChar: 14,
	},
	{
		Input:     "select 1 < 2 < 3",
		Error:     "syntax error: unexpected token \"<\"",
		ErrorLine: 1,
		ErrorChar: 14,
	},
	{
		Input:     "select 'literal not terminated",
		Error:     "literal not terminated",
		ErrorLine: 1,
		ErrorChar: 30,
	},
	{
		Input:      "select select",
		SourceFile: GetTestFilePath("dummy.sql"),
		Error:      "syntax error: unexpected token \"select\"",
		ErrorLine:  1,
		ErrorChar:  8,
		ErrorFile:  GetTestFilePath("dummy.sql"),
	},
	{
		Input:      "print 'foo' 'bar'",
		SourceFile: GetTestFilePath("dummy.sql"),
		Error:      "syntax error: unexpected token \"bar\"",
		ErrorLine:  1,
		ErrorChar:  13,
		ErrorFile:  GetTestFilePath("dummy.sql"),
	},
	{
		Input:      "print !=",
		SourceFile: GetTestFilePath("dummy.sql"),
		Error:      "syntax error: unexpected token \"!=\"",
		ErrorLine:  1,
		ErrorChar:  7,
		ErrorFile:  GetTestFilePath("dummy.sql"),
	},
}

func TestParse(t *testing.T) {
	for _, v := range parseTests {
		prog, holderNum, err := Parse(v.Input, v.SourceFile, nil, v.ForPrepared)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}

			syntaxErr := err.(*SyntaxError)
			if syntaxErr.Line != v.ErrorLine {
				t.Errorf("error line %d, want error line %d for %q", syntaxErr.Line, v.ErrorLine, v.Input)
			}
			if syntaxErr.Char != v.ErrorChar {
				t.Errorf("error char %d, want error char %d for %q", syntaxErr.Char, v.ErrorChar, v.Input)
			}
			if syntaxErr.SourceFile != v.ErrorFile {
				t.Errorf("error file %s, want error file %s for %q", syntaxErr.SourceFile, v.ErrorFile, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}

		if v.ForPrepared {
			if holderNum != v.HolderNum {
				t.Errorf("holder number = %d, want %d for %q", holderNum, v.HolderNum, v.Input)
			}
		}

		if len(v.Output) != len(prog) {
			t.Errorf("parsed program has %d statement(s), want %d statement(s) for %q", len(prog), len(v.Output), v.Input)
			continue
		}

		for i, stmt := range prog {
			expect := v.Output[i]

			stmtType := reflect.TypeOf(stmt).Name()
			expectType := reflect.TypeOf(expect).Name()

			if stmtType != expectType {
				t.Errorf("statement type is %q, want %q for %q", stmtType, expectType, v.Input)
				continue
			}

			switch stmtType {
			case "SelectQuery":
				expectStmt := expect.(SelectQuery)
				parsedStmt := stmt.(SelectQuery)

				if entity, ok := parsedStmt.SelectEntity.(SelectEntity); ok {
					expectEntity, ok := expectStmt.SelectEntity.(SelectEntity)
					if !ok {
						t.Errorf("entity = %#v, want %#v for %q", entity, expectEntity, v.Input)
					}

					if !reflect.DeepEqual(entity.SelectClause, expectEntity.SelectClause) {
						t.Errorf("select clause = %#v, want %#v for %q", entity.SelectClause, expectEntity.SelectClause, v.Input)
					}
					if !reflect.DeepEqual(entity.FromClause, expectEntity.FromClause) {
						t.Errorf("from clause = %#v, want %#v for %q", entity.FromClause, expectEntity.FromClause, v.Input)
					}
					if !reflect.DeepEqual(entity.WhereClause, expectEntity.WhereClause) {
						t.Errorf("where clause = %#v, want %#v for %q", entity.WhereClause, expectEntity.WhereClause, v.Input)
					}
					if !reflect.DeepEqual(entity.GroupByClause, expectEntity.GroupByClause) {
						t.Errorf("group by clause = %#v, want %#v for %q", entity.GroupByClause, expectEntity.GroupByClause, v.Input)
					}
					if !reflect.DeepEqual(entity.HavingClause, expectEntity.HavingClause) {
						t.Errorf("having clause = %#v, want %#v for %q", entity.HavingClause, expectEntity.HavingClause, v.Input)
					}
				} else if set, ok := parsedStmt.SelectEntity.(SelectSet); ok {
					expectSet, ok := expectStmt.SelectEntity.(SelectSet)
					if !ok {
						t.Errorf("set = %#v, want %#v for %q", set, expectSet, v.Input)
					}

					if !reflect.DeepEqual(set, expectSet) {
						t.Errorf("set = %#v, want %#v for %q", set, expectSet, v.Input)
					}
				}

				if !reflect.DeepEqual(parsedStmt.WithClause, expectStmt.WithClause) {
					t.Errorf("with clause = %#v, want %#v for %q", parsedStmt.WithClause, expectStmt.WithClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.OrderByClause, expectStmt.OrderByClause) {
					t.Errorf("order by clause = %#v, want %#v for %q", parsedStmt.OrderByClause, expectStmt.OrderByClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.LimitClause, expectStmt.LimitClause) {
					t.Errorf("limit clause = %#v, want %#v for %q", parsedStmt.LimitClause, expectStmt.LimitClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.OffsetClause, expectStmt.OffsetClause) {
					t.Errorf("offset clause = %#v, want %#v for %q", parsedStmt.OffsetClause, expectStmt.OffsetClause, v.Input)
				}
			default:
				if !reflect.DeepEqual(stmt, expect) {
					t.Errorf("output = %#v, want %#v for %q", stmt, expect, v.Input)
				}
			}
		}
	}
}
