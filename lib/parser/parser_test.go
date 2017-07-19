package parser

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/ternary"
)

var parseTests = []struct {
	Input  string
	Output []Statement
	Error  string
}{
	{
		Input: "select foo; select bar;",
		Output: []Statement{
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: FieldReference{Column: Identifier{Literal: "foo"}}}}},
			}},
			SelectQuery{SelectEntity: SelectEntity{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: FieldReference{Column: Identifier{Literal: "bar"}}}}},
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
							SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewInteger(1)}}},
						},
						Operator: Token{Token: UNION, Literal: "union", Line: 1, Char: 10},
						All:      Token{Token: ALL, Literal: "all", Line: 1, Char: 16},
						RHS: SelectSet{
							LHS: SelectEntity{
								SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewInteger(2)}}},
							},
							Operator: Token{Token: INTERSECT, Literal: "intersect", Line: 1, Char: 29},
							RHS: SelectEntity{
								SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewInteger(3)}}},
							},
						},
					},
					Operator: Token{Token: EXCEPT, Literal: "except", Line: 1, Char: 48},
					RHS: SelectEntity{
						SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewInteger(4)}}},
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
						SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewInteger(1)}}},
					},
					Operator: Token{Token: UNION, Literal: "union", Line: 1, Char: 10},
					RHS: Subquery{
						Query: SelectQuery{
							SelectEntity: SelectEntity{
								SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewInteger(2)}}},
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
						Select: "select",
						Fields: []Expression{
							Field{
								Object: NewInteger(1),
								As:     "as",
								Alias:  Identifier{Literal: "a"},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []Expression{
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
						Select: "select",
						Fields: []Expression{
							Field{
								Object: FieldReference{Column: Identifier{Literal: "c1"}},
							},
						},
					},
					FromClause: FromClause{From: "from", Tables: []Expression{
						Table{Object: Stdin{Stdin: "stdin"}},
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
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Identifier{Literal: "table1"},
							},
							Table{
								Object: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("2")}}},
											FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
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
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Identifier{Literal: "table1"},
								Alias:  Identifier{Literal: "alias"},
							},
							Table{
								Object: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("2")}}},
											FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
										},
									},
								},
								Alias: Identifier{Literal: "alias2"},
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
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Identifier{Literal: "table1"},
								As:     "as",
								Alias:  Identifier{Literal: "alias"},
							},
							Table{
								Object: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("2")}}},
											FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
										},
									},
								},
								As:    "as",
								Alias: Identifier{Literal: "alias2"},
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
			"          avg() over () \n" +
			" limit 10 \n" +
			" offset 10 \n",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
					WhereClause: WhereClause{
						Where: "where",
						Filter: Comparison{
							LHS:      NewIntegerFromString("1"),
							Operator: "=",
							RHS:      NewIntegerFromString("1"),
						},
					},
					GroupByClause: GroupByClause{
						GroupBy: "group by",
						Items: []Expression{
							FieldReference{Column: Identifier{Literal: "column1"}},
							FieldReference{Column: Identifier{Literal: "column2"}},
						},
					},
					HavingClause: HavingClause{
						Having: "having",
						Filter: Comparison{
							LHS:      NewIntegerFromString("1"),
							Operator: ">",
							RHS:      NewIntegerFromString("1"),
						},
					},
				},
				OrderByClause: OrderByClause{
					OrderBy: "order by",
					Items: []Expression{
						OrderItem{Value: FieldReference{Column: Identifier{Literal: "column4"}}},
						OrderItem{Value: FieldReference{Column: Identifier{Literal: "column5"}}, Direction: Token{Token: DESC, Literal: "desc", Line: 7, Char: 19}},
						OrderItem{Value: FieldReference{Column: Identifier{Literal: "column6"}}, Direction: Token{Token: ASC, Literal: "asc", Line: 8, Char: 19}},
						OrderItem{Value: FieldReference{Column: Identifier{Literal: "column7"}}, Nulls: "nulls", Position: Token{Token: FIRST, Literal: "first", Line: 9, Char: 25}},
						OrderItem{Value: FieldReference{Column: Identifier{Literal: "column8"}}, Direction: Token{Token: DESC, Literal: "desc", Line: 10, Char: 19}, Nulls: "nulls", Position: Token{Token: LAST, Literal: "last", Line: 10, Char: 30}},
						OrderItem{Value: AnalyticFunction{
							Name: "avg",
							Over: "over",
							AnalyticClause: AnalyticClause{
								Partition:     nil,
								OrderByClause: nil,
							},
						}},
					},
				},
				LimitClause: LimitClause{
					Limit: "limit",
					Value: NewInteger(10),
				},
				OffsetClause: OffsetClause{
					Offset: "offset",
					Value:  NewInteger(10),
				},
			},
		},
	},
	{
		Input: "select 1 " +
			" from dual " +
			" limit 10 percent",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
				},
				LimitClause: LimitClause{
					Limit:   "limit",
					Value:   NewInteger(10),
					Percent: "percent",
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
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
				},
				LimitClause: LimitClause{
					Limit: "limit",
					Value: NewInteger(10),
					With:  LimitWith{With: "with", Type: Token{Token: TIES, Literal: "ties", Line: 3, Char: 16}},
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
						Select:   "select",
						Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 8},
						Fields: []Expression{
							Field{Object: AllColumns{}},
						},
					},
					FromClause: FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
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
					InlineTables: []Expression{
						InlineTable{
							Name: Identifier{Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											Field{Object: NewInteger(1)},
										},
									},
								},
							},
						},
					},
				},
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{Field{Object: AllColumns{}}},
					},
					FromClause: FromClause{
						From:   "from",
						Tables: []Expression{Table{Object: Identifier{Literal: "ct"}}},
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
					InlineTables: []Expression{
						InlineTable{
							Name: Identifier{Literal: "ct"},
							Columns: []Expression{
								Identifier{Literal: "column1"},
							},
							As: "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											Field{Object: NewInteger(1)},
										},
									},
								},
							},
						},
					},
				},
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{Field{Object: AllColumns{}}},
					},
					FromClause: FromClause{
						From:   "from",
						Tables: []Expression{Table{Object: Identifier{Literal: "ct"}}},
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
					InlineTables: []Expression{
						InlineTable{
							Name:      Identifier{Literal: "ct"},
							Recursive: Token{Token: RECURSIVE, Literal: "recursive"},
							As:        "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											NewInteger(1),
										},
									},
								},
							},
						},
						InlineTable{
							Name: Identifier{Literal: "ct2"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											NewInteger(2),
										},
									},
								},
							},
						},
					},
				},
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{Field{Object: AllColumns{}}},
					},
					FromClause: FromClause{
						From:   "from",
						Tables: []Expression{Table{Object: Identifier{Literal: "ct"}}},
					},
				},
			},
		},
	},
	{
		Input: "select ident, 'foo', 1, -1, 1.234, -1.234, true, '2010-01-01 12:00:00', null, ('bar') from dual",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: FieldReference{Column: Identifier{Literal: "ident"}}},
							Field{Object: NewString("foo")},
							Field{Object: NewIntegerFromString("1")},
							Field{Object: NewIntegerFromString("-1")},
							Field{Object: NewFloatFromString("1.234")},
							Field{Object: NewFloatFromString("-1.234")},
							Field{Object: NewTernaryFromString("true")},
							Field{Object: NewDatetimeFromString("2010-01-01 12:00:00")},
							Field{Object: NewNullFromString("null")},
							Field{Object: Parentheses{Expr: NewString("bar")}},
						},
					},
					FromClause: FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
				},
			},
		},
	},
	{
		Input: "select ident || 'foo' || 'bar'",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Concat{Items: []Expression{
								FieldReference{Column: Identifier{Literal: "ident"}},
								NewString("foo"),
								NewString("bar"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Comparison{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: "=",
								RHS:      NewInteger(1),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Comparison{
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								RHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											NewInteger(1),
											NewInteger(2),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Comparison{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: "<",
								RHS:      NewInteger(1),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Comparison{
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Operator: "<",
								RHS: RowValue{
									Value: Subquery{
										Query: SelectQuery{
											SelectEntity: SelectEntity{
												SelectClause: SelectClause{
													Select: "select",
													Fields: []Expression{
														Field{Object: NewInteger(1)},
														Field{Object: NewInteger(2)},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Is{
								Is:       "is",
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								RHS:      NewNullFromString("null"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Is{
								Is:  "is",
								LHS: FieldReference{Column: Identifier{Literal: "column1"}},
								RHS: NewTernaryFromString("true"),
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not between -10 and 10",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Between{
								Between:  "between",
								And:      "and",
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Low:      NewIntegerFromString("-10"),
								High:     NewIntegerFromString("10"),
								Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 16},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select (column1, column2) between (1, 2) and (3, 4)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Between{
								Between: "between",
								And:     "and",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Low: RowValue{
									Value: ValueList{
										Values: []Expression{
											NewInteger(1),
											NewInteger(2),
										},
									},
								},
								High: RowValue{
									Value: ValueList{
										Values: []Expression{
											NewInteger(3),
											NewInteger(4),
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
		Input: "select column1 not in (1, 2, 3)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: In{
								In:  "in",
								LHS: FieldReference{Column: Identifier{Literal: "column1"}},
								Values: RowValue{
									Value: ValueList{
										Values: []Expression{
											NewIntegerFromString("1"),
											NewIntegerFromString("2"),
											NewIntegerFromString("3"),
										},
									},
								},
								Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 16},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: In{
								In: "in",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Values: RowValueList{
									RowValues: []Expression{
										RowValue{
											Value: ValueList{
												Values: []Expression{
													NewInteger(1),
													NewInteger(2),
												},
											},
										},
										RowValue{
											Value: ValueList{
												Values: []Expression{
													NewInteger(3),
													NewInteger(4),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: In{
								In: "in",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Values: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
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
		Input: "select column1 not like 'pattern'",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Like{
								Like:     "like",
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Pattern:  String{literal: "pattern"},
								Negation: Token{Token: NOT, Literal: "not", Line: 1, Char: 16},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Any{
								Any:      "any",
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: "=",
								Values: RowValue{
									Value: Subquery{
										Query: SelectQuery{
											SelectEntity: SelectEntity{
												SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Any{
								Any: "any",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: RowValueList{
									RowValues: []Expression{
										RowValue{
											Value: ValueList{
												Values: []Expression{
													NewInteger(1),
													NewInteger(2),
												},
											},
										},
										RowValue{
											Value: ValueList{
												Values: []Expression{
													NewInteger(3),
													NewInteger(4),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Any{
								Any: "any",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: All{
								All:      "all",
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: "=",
								Values: RowValue{
									Subquery{
										Query: SelectQuery{
											SelectEntity: SelectEntity{
												SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: All{
								All: "all",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: RowValueList{
									RowValues: []Expression{
										RowValue{
											Value: ValueList{
												Values: []Expression{
													NewInteger(1),
													NewInteger(2),
												},
											},
										},
										RowValue{
											Value: ValueList{
												Values: []Expression{
													NewInteger(3),
													NewInteger(4),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: All{
								All: "all",
								LHS: RowValue{
									Value: ValueList{
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
											FieldReference{Column: Identifier{Literal: "column2"}},
										},
									},
								},
								Operator: "=",
								Values: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Exists{
								Exists: "exists",
								Query: Subquery{
									Query: SelectQuery{
										SelectEntity: SelectEntity{
											SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: int('+'),
								RHS:      NewIntegerFromString("1"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: int('-'),
								RHS:      NewIntegerFromString("1"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: int('*'),
								RHS:      NewIntegerFromString("1"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: int('/'),
								RHS:      NewIntegerFromString("1"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Arithmetic{
								LHS:      FieldReference{Column: Identifier{Literal: "column1"}},
								Operator: int('%'),
								RHS:      NewIntegerFromString("1"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Logic{
								LHS:      NewTernaryFromString("true"),
								Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 13},
								RHS:      NewTernaryFromString("false"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Logic{
								LHS:      NewTernaryFromString("true"),
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 13},
								RHS:      NewTernaryFromString("false"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Logic{
								Operator: Token{Token: NOT, Literal: "not", Line: 1, Char: 8},
								RHS:      NewTernaryFromString("false"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Logic{
								LHS:      NewTernaryFromString("true"),
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 13},
								RHS: Parentheses{
									Expr: Logic{
										LHS:      NewTernaryFromString("false"),
										Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 23},
										RHS:      NewTernaryFromString("false"),
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
		Input: "select true and true or false and not false",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Logic{
								LHS: Logic{
									LHS:      NewTernaryFromString("true"),
									Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 13},
									RHS:      NewTernaryFromString("true"),
								},
								Operator: Token{Token: OR, Literal: "or", Line: 1, Char: 22},
								RHS: Logic{
									LHS:      NewTernaryFromString("false"),
									Operator: Token{Token: AND, Literal: "and", Line: 1, Char: 31},
									RHS: Logic{
										Operator: Token{Token: NOT, Literal: "not", Line: 1, Char: 35},
										RHS:      NewTernaryFromString("false"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Variable{Name: "@var"}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: VariableSubstitution{
								Variable: Variable{Name: "@var"},
								Value:    NewInteger(1),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Case{
								Case: "case",
								End:  "end",
								When: []Expression{
									CaseWhen{
										When:      "when",
										Then:      "then",
										Condition: Ternary{literal: "true", value: ternary.TRUE},
										Result:    NewString("A"),
									},
									CaseWhen{
										When:      "when",
										Then:      "then",
										Condition: Ternary{literal: "false", value: ternary.FALSE},
										Result:    NewString("B"),
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Case{
								Case:  "case",
								End:   "end",
								Value: FieldReference{Column: Identifier{Literal: "column1"}},
								When: []Expression{
									CaseWhen{
										When:      "when",
										Then:      "then",
										Condition: NewIntegerFromString("1"),
										Result:    NewString("A"),
									},
									CaseWhen{
										When:      "when",
										Then:      "then",
										Condition: NewIntegerFromString("2"),
										Result:    NewString("B"),
									},
								},
								Else: CaseElse{
									Else:   "else",
									Result: NewString("C"),
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select count()",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Function{
								Name: "count",
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
						Select: "select",
						Fields: []Expression{
							Field{Object: Function{
								Name: "count",
								Args: []Expression{
									FieldReference{Column: Identifier{Literal: "column1"}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: AggregateFunction{
								Name: "count",
								Option: AggregateOption{
									Args: []Expression{AllColumns{}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: AggregateFunction{
								Name: "count",
								Option: AggregateOption{
									Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 14},
									Args:     []Expression{AllColumns{}},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: AggregateFunction{
								Name: "count",
								Option: AggregateOption{
									Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 14},
									Args: []Expression{
										FieldReference{Column: Identifier{Literal: "column1"}},
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
		Input: "select count(column1, column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: Function{
								Name: "count",
								Args: []Expression{
									FieldReference{Column: Identifier{Literal: "column1"}},
									FieldReference{Column: Identifier{Literal: "column2"}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select group_concat(column1 order by column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: GroupConcat{
								GroupConcat: "group_concat",
								Option:      AggregateOption{Args: []Expression{FieldReference{Column: Identifier{Literal: "column1"}}}},
								OrderBy: OrderByClause{
									OrderBy: "order by",
									Items: []Expression{
										OrderItem{Value: FieldReference{Column: Identifier{Literal: "column1"}}},
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
		Input: "select group_concat(distinct column1 order by column1)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: GroupConcat{
								GroupConcat: "group_concat",
								Option: AggregateOption{
									Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 21},
									Args:     []Expression{FieldReference{Column: Identifier{Literal: "column1"}}},
								},
								OrderBy: OrderByClause{
									OrderBy: "order by",
									Items: []Expression{
										OrderItem{Value: FieldReference{Column: Identifier{Literal: "column1"}}},
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
		Input: "select group_concat(distinct column1 separator ',')",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: GroupConcat{
								GroupConcat: "group_concat",
								Option: AggregateOption{
									Distinct: Token{Token: DISTINCT, Literal: "distinct", Line: 1, Char: 21},
									Args:     []Expression{FieldReference{Column: Identifier{Literal: "column1"}}},
								},
								SeparatorLit: "separator",
								Separator:    ",",
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select group_concat(column1 order by column1 separator ',')",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: GroupConcat{
								GroupConcat: "group_concat",
								Option: AggregateOption{
									Args: []Expression{FieldReference{Column: Identifier{Literal: "column1"}}},
								},
								OrderBy: OrderByClause{
									OrderBy: "order by",
									Items: []Expression{
										OrderItem{Value: FieldReference{Column: Identifier{Literal: "column1"}}},
									},
								},
								SeparatorLit: "separator",
								Separator:    ",",
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
						Select: "select",
						Fields: []Expression{
							Field{Object: CursorStatus{
								CursorLit: "cursor",
								Cursor:    Identifier{Literal: "cur"},
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
						Select: "select",
						Fields: []Expression{
							Field{Object: CursorStatus{
								CursorLit: "cursor",
								Cursor:    Identifier{Literal: "cur"},
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
		Input: "select rank() over (partition by column1 order by column2)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: AnalyticFunction{
								Name: "rank",
								Over: "over",
								AnalyticClause: AnalyticClause{
									Partition: Partition{
										PartitionBy: "partition by",
										Values: []Expression{
											FieldReference{Column: Identifier{Literal: "column1"}},
										},
									},
									OrderByClause: OrderByClause{
										OrderBy: "order by",
										Items: []Expression{
											OrderItem{
												Value: FieldReference{Column: Identifier{Literal: "column2"}},
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
		Input: "select avg() over ()",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: AnalyticFunction{
								Name: "avg",
								Over: "over",
								AnalyticClause: AnalyticClause{
									Partition:     nil,
									OrderByClause: nil,
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Input: "select first_value(column1) over ()",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{
						Select: "select",
						Fields: []Expression{
							Field{Object: AnalyticFunction{
								Name: "first_value",
								Args: []Expression{
									FieldReference{Column: Identifier{Literal: "column1"}},
								},
								Over: "over",
								AnalyticClause: AnalyticClause{
									Partition:     nil,
									OrderByClause: nil,
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
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
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
		Input: "select 1 from table1 inner join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									JoinType:  Token{Token: INNER, Literal: "inner", Line: 1, Char: 22},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 join table2 on table1.id = table2.id",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									Condition: JoinCondition{
										Literal: "on",
										On: Comparison{
											LHS:      FieldReference{View: Identifier{Literal: "table1"}, Column: Identifier{Literal: "id"}},
											Operator: "=",
											RHS:      FieldReference{View: Identifier{Literal: "table2"}, Column: Identifier{Literal: "id"}},
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
		Input: "select 1 from table1 natural join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									Natural:   Token{Token: NATURAL, Literal: "natural", Line: 1, Char: 22},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 left join table2 using(id)",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									Direction: Token{Token: LEFT, Literal: "left", Line: 1, Char: 22},
									Condition: JoinCondition{
										Literal: "using",
										Using: []Expression{
											Identifier{Literal: "id"},
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
		Input: "select 1 from table1 natural outer join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									Natural:   Token{Token: NATURAL, Literal: "natural", Line: 1, Char: 22},
									JoinType:  Token{Token: OUTER, Literal: "outer", Line: 1, Char: 30},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 right join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									Direction: Token{Token: RIGHT, Literal: "right", Line: 1, Char: 22},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 full join table2",
		Output: []Statement{
			SelectQuery{
				SelectEntity: SelectEntity{
					SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
					FromClause: FromClause{
						From: "from",
						Tables: []Expression{
							Table{
								Object: Join{
									Join:      "join",
									Table:     Table{Object: Identifier{Literal: "table1"}},
									JoinTable: Table{Object: Identifier{Literal: "table2"}},
									Direction: Token{Token: FULL, Literal: "full", Line: 1, Char: 22},
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
				Assignments: []Expression{
					VariableAssignment{
						Name: "@var1",
					},
					VariableAssignment{
						Name:  "@var2",
						Value: NewInteger(2),
					},
				},
			},
			VariableSubstitution{
				Variable: Variable{
					Name: "@var1",
				},
				Value: NewInteger(1),
			},
		},
	},
	{
		Input: "func('arg1', 'arg2')",
		Output: []Statement{
			Function{
				Name: "func",
				Args: []Expression{
					NewString("arg1"),
					NewString("arg2"),
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
					InlineTables: []Expression{
						InlineTable{
							Name: Identifier{Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											Field{Object: NewInteger(1)},
										},
									},
								},
							},
						},
					},
				},
				Insert: "insert",
				Into:   "into",
				Table:  Identifier{Literal: "table1"},
				Values: "values",
				ValuesList: []Expression{
					RowValue{
						Value: ValueList{
							Values: []Expression{
								NewInteger(1),
								NewString("str1"),
							},
						},
					},
					RowValue{
						Value: ValueList{
							Values: []Expression{
								NewInteger(2),
								NewString("str2"),
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "insert into table1 (column1, column2) values (1, 'str1'), (2, 'str2')",
		Output: []Statement{
			InsertQuery{
				Insert: "insert",
				Into:   "into",
				Table:  Identifier{Literal: "table1"},
				Fields: []Expression{
					FieldReference{Column: Identifier{Literal: "column1"}},
					FieldReference{Column: Identifier{Literal: "column2"}},
				},
				Values: "values",
				ValuesList: []Expression{
					RowValue{
						Value: ValueList{
							Values: []Expression{
								NewInteger(1),
								NewString("str1"),
							},
						},
					},
					RowValue{
						Value: ValueList{
							Values: []Expression{
								NewInteger(2),
								NewString("str2"),
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
				Insert: "insert",
				Into:   "into",
				Table:  Identifier{Literal: "table1"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								Field{Object: NewInteger(1)},
								Field{Object: NewInteger(2)},
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
				Insert: "insert",
				Into:   "into",
				Table:  Identifier{Literal: "table1"},
				Fields: []Expression{
					FieldReference{Column: Identifier{Literal: "column1"}},
					FieldReference{Column: Identifier{Literal: "column2"}},
				},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								Field{Object: NewInteger(1)},
								Field{Object: NewInteger(2)},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "with ct as (select 1) update table1 set column1 = 1, column2 = 2 from table1 where true",
		Output: []Statement{
			UpdateQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []Expression{
						InlineTable{
							Name: Identifier{Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											Field{Object: NewInteger(1)},
										},
									},
								},
							},
						},
					},
				},
				Update: "update",
				Tables: []Expression{
					Table{Object: Identifier{Literal: "table1"}},
				},
				Set: "set",
				SetList: []Expression{
					UpdateSet{Field: FieldReference{Column: Identifier{Literal: "column1"}}, Value: NewInteger(1)},
					UpdateSet{Field: FieldReference{Column: Identifier{Literal: "column2"}}, Value: NewInteger(2)},
				},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{Object: Identifier{Literal: "table1"}},
					},
				},
				WhereClause: WhereClause{
					Where:  "where",
					Filter: Ternary{literal: "true", value: ternary.TRUE},
				},
			},
		},
	},
	{
		Input: "with ct as (select 1) delete from table1",
		Output: []Statement{
			DeleteQuery{
				WithClause: WithClause{
					With: "with",
					InlineTables: []Expression{
						InlineTable{
							Name: Identifier{Literal: "ct"},
							As:   "as",
							Query: SelectQuery{
								SelectEntity: SelectEntity{
									SelectClause: SelectClause{
										Select: "select",
										Fields: []Expression{
											Field{Object: NewInteger(1)},
										},
									},
								},
							},
						},
					},
				},
				Delete: "delete",
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{Object: Identifier{Literal: "table1"}},
					},
				},
			},
		},
	},
	{
		Input: "delete table1 from table1 where true",
		Output: []Statement{
			DeleteQuery{
				Delete: "delete",
				Tables: []Expression{
					Table{Object: Identifier{Literal: "table1"}},
				},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{Object: Identifier{Literal: "table1"}},
					},
				},
				WhereClause: WhereClause{
					Where:  "where",
					Filter: Ternary{literal: "true", value: ternary.TRUE},
				},
			},
		},
	},
	{
		Input: "create table newtable (column1, column2)",
		Output: []Statement{
			CreateTable{
				CreateTable: "create table",
				Table:       Identifier{Literal: "newtable"},
				Fields: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
		},
	},
	{
		Input: "alter table table1 add column1",
		Output: []Statement{
			AddColumns{
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Add:        "add",
				Columns: []Expression{
					ColumnDefault{
						Column: Identifier{Literal: "column1"},
					},
				},
			},
		},
	},
	{
		Input: "alter table table1 add (column1, column2 default 1) first",
		Output: []Statement{
			AddColumns{
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Add:        "add",
				Columns: []Expression{
					ColumnDefault{
						Column: Identifier{Literal: "column1"},
					},
					ColumnDefault{
						Column:  Identifier{Literal: "column2"},
						Default: "default",
						Value:   NewInteger(1),
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
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Add:        "add",
				Columns: []Expression{
					ColumnDefault{
						Column: Identifier{Literal: "column1"},
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
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Add:        "add",
				Columns: []Expression{
					ColumnDefault{
						Column: Identifier{Literal: "column1"},
					},
				},
				Position: ColumnPosition{
					Position: Token{Token: AFTER, Literal: "after", Line: 1, Char: 32},
					Column:   FieldReference{Column: Identifier{Literal: "column2"}},
				},
			},
		},
	},
	{
		Input: "alter table table1 add column1 before column2",
		Output: []Statement{
			AddColumns{
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Add:        "add",
				Columns: []Expression{
					ColumnDefault{
						Column: Identifier{Literal: "column1"},
					},
				},
				Position: ColumnPosition{
					Position: Token{Token: BEFORE, Literal: "before", Line: 1, Char: 32},
					Column:   FieldReference{Column: Identifier{Literal: "column2"}},
				},
			},
		},
	},
	{
		Input: "alter table table1 drop column1",
		Output: []Statement{
			DropColumns{
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Drop:       "drop",
				Columns:    []Expression{FieldReference{Column: Identifier{Literal: "column1"}}},
			},
		},
	},
	{
		Input: "alter table table1 drop (column1, column2)",
		Output: []Statement{
			DropColumns{
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Drop:       "drop",
				Columns: []Expression{
					FieldReference{Column: Identifier{Literal: "column1"}},
					FieldReference{Column: Identifier{Literal: "column2"}},
				},
			},
		},
	},
	{
		Input: "alter table table1 rename column1 to column2",
		Output: []Statement{
			RenameColumn{
				AlterTable: "alter table",
				Table:      Identifier{Literal: "table1"},
				Rename:     "rename",
				Old:        FieldReference{Column: Identifier{Literal: "column1"}},
				To:         "to",
				New:        Identifier{Literal: "column2"},
			},
		},
	},
	{
		Input: "commit",
		Output: []Statement{
			TransactionControl{
				Token: COMMIT,
			},
		},
	},
	{
		Input: "rollback",
		Output: []Statement{
			TransactionControl{
				Token: ROLLBACK,
			},
		},
	},
	{
		Input: "print 'foo'",
		Output: []Statement{
			Print{
				Value: NewString("foo"),
			},
		},
	},
	{
		Input: "printf 'foo'",
		Output: []Statement{
			Printf{
				Values: []Expression{
					NewString("foo"),
				},
			},
		},
	},
	{
		Input: "source '/path/to/file.sql'",
		Output: []Statement{
			Source{
				FilePath: "/path/to/file.sql",
			},
		},
	},
	{
		Input: "set @@delimiter = ','",
		Output: []Statement{
			SetFlag{
				Name:  "@@delimiter",
				Value: NewString(","),
			},
		},
	},
	{
		Input: "declare cur cursor for select 1",
		Output: []Statement{
			CursorDeclaration{
				Cursor: Identifier{Literal: "cur"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								Field{Object: NewInteger(1)},
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
				Cursor: Identifier{Literal: "cur"},
			},
		},
	},
	{
		Input: "close cur",
		Output: []Statement{
			CloseCursor{
				Cursor: Identifier{Literal: "cur"},
			},
		},
	},
	{
		Input: "dispose cursor cur",
		Output: []Statement{
			DisposeCursor{
				Cursor: Identifier{Literal: "cur"},
			},
		},
	},
	{
		Input: "fetch cur into @var1, @var2",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Variables: []Variable{
					{Name: "@var1"},
					{Name: "@var2"},
				},
			},
		},
	},
	{
		Input: "fetch next cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: NEXT, Literal: "next", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{Name: "@var1"},
				},
			},
		},
	},
	{
		Input: "fetch prior cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: PRIOR, Literal: "prior", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{Name: "@var1"},
				},
			},
		},
	},
	{
		Input: "fetch first cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: FIRST, Literal: "first", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{Name: "@var1"},
				},
			},
		},
	},
	{
		Input: "fetch last cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: LAST, Literal: "last", Line: 1, Char: 7},
				},
				Variables: []Variable{
					{Name: "@var1"},
				},
			},
		},
	},
	{
		Input: "fetch absolute 1 cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: ABSOLUTE, Literal: "absolute", Line: 1, Char: 7},
					Number:   NewInteger(1),
				},
				Variables: []Variable{
					{Name: "@var1"},
				},
			},
		},
	},
	{
		Input: "fetch relative 1 cur into @var1",
		Output: []Statement{
			FetchCursor{
				Cursor: Identifier{Literal: "cur"},
				Position: FetchPosition{
					Position: Token{Token: RELATIVE, Literal: "relative", Line: 1, Char: 7},
					Number:   NewInteger(1),
				},
				Variables: []Variable{
					{Name: "@var1"},
				},
			},
		},
	},
	{
		Input: "declare tbl table (column1, column2)",
		Output: []Statement{
			TableDeclaration{
				Table: Identifier{Literal: "tbl"},
				Fields: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
			},
		},
	},
	{
		Input: "declare tbl table (column1, column2) for select 1, 2",
		Output: []Statement{
			TableDeclaration{
				Table: Identifier{Literal: "tbl"},
				Fields: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								Field{
									Object: NewInteger(1),
								},
								Field{
									Object: NewInteger(2),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "declare tbl table for select 1, 2",
		Output: []Statement{
			TableDeclaration{
				Table: Identifier{Literal: "tbl"},
				Query: SelectQuery{
					SelectEntity: SelectEntity{
						SelectClause: SelectClause{
							Select: "select",
							Fields: []Expression{
								Field{
									Object: NewInteger(1),
								},
								Field{
									Object: NewInteger(2),
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "dispose table tbl",
		Output: []Statement{
			DisposeTable{
				Table: Identifier{Literal: "tbl"},
			},
		},
	},
	{
		Input: "if @var1 = 1 then print 1; end if",
		Output: []Statement{
			If{
				Condition: Comparison{
					LHS:      Variable{Name: "@var1"},
					RHS:      NewInteger(1),
					Operator: "=",
				},
				Statements: []Statement{
					Print{Value: NewInteger(1)},
				},
			},
		},
	},
	{
		Input: "if @var1 = 1 then print 1; elseif @var1 = 2 then print 2; elseif @var1 = 3 then print 3; else print 4; end if",
		Output: []Statement{
			If{
				Condition: Comparison{
					LHS:      Variable{Name: "@var1"},
					RHS:      NewInteger(1),
					Operator: "=",
				},
				Statements: []Statement{
					Print{Value: NewInteger(1)},
				},
				ElseIf: []ProcExpr{
					ElseIf{
						Condition: Comparison{
							LHS:      Variable{Name: "@var1"},
							RHS:      NewInteger(2),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewInteger(2)},
						},
					},
					ElseIf{
						Condition: Comparison{
							LHS:      Variable{Name: "@var1"},
							RHS:      NewInteger(3),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewInteger(3)},
						},
					},
				},
				Else: Else{
					Statements: []Statement{
						Print{Value: NewInteger(4)},
					},
				},
			},
		},
	},
	{
		Input: "while @var1 do print @var1 end while",
		Output: []Statement{
			While{
				Condition: Variable{Name: "@var1"},
				Statements: []Statement{
					Print{Value: Variable{Name: "@var1"}},
				},
			},
		},
	},
	{
		Input: "while @var1, @var2 in cur do print @var1 end while",
		Output: []Statement{
			WhileInCursor{
				Variables: []Variable{
					{Name: "@var1"},
					{Name: "@var2"},
				},
				Cursor: Identifier{Literal: "cur"},
				Statements: []Statement{
					Print{Value: Variable{Name: "@var1"}},
				},
			},
		},
	},
	{
		Input: "exit",
		Output: []Statement{
			FlowControl{Token: EXIT},
		},
	},
	{
		Input: "while true do continue end while",
		Output: []Statement{
			While{
				Condition: Ternary{literal: "true", value: ternary.TRUE},
				Statements: []Statement{
					FlowControl{Token: CONTINUE},
				},
			},
		},
	},
	{
		Input: "while true do break end while",
		Output: []Statement{
			While{
				Condition: Ternary{literal: "true", value: ternary.TRUE},
				Statements: []Statement{
					FlowControl{Token: BREAK},
				},
			},
		},
	},
	{
		Input: "while true do if @var1 = 1 then continue; end if; end while",
		Output: []Statement{
			While{
				Condition: Ternary{literal: "true", value: ternary.TRUE},
				Statements: []Statement{
					If{
						Condition: Comparison{
							LHS:      Variable{Name: "@var1"},
							RHS:      NewInteger(1),
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
				Condition: Ternary{literal: "true", value: ternary.TRUE},
				Statements: []Statement{
					If{
						Condition: Comparison{
							LHS:      Variable{Name: "@var1"},
							RHS:      NewInteger(1),
							Operator: "=",
						},
						Statements: []Statement{
							FlowControl{Token: CONTINUE},
						},
						ElseIf: []ProcExpr{
							ElseIf{
								Condition: Comparison{
									LHS:      Variable{Name: "@var1"},
									RHS:      NewInteger(2),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: BREAK},
								},
							},
							ElseIf{
								Condition: Comparison{
									LHS:      Variable{Name: "@var1"},
									RHS:      NewInteger(3),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: EXIT},
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
		Input: "declare func1 function () as begin end",
		Output: []Statement{
			FunctionDeclaration{
				Name: Identifier{Literal: "func1"},
			},
		},
	},
	{
		Input: "declare func1 function (@arg1, @arg2) as begin " +
			"if @var1 = 1 then print 1; end if; " +
			"if @var1 = 1 then print 1; elseif @var1 = 2 then print 2; elseif @var1 = 3 then print 3; else print 4; end if; " +
			"while true do break end while; " +
			"while true do if @var1 = 1 then continue; end if; end while; " +
			"while true do if @var1 = 1 then continue; elseif @var1 = 2 then break; elseif @var1 = 3 then return; else continue; end if; end while; " +
			"return @var1; " +
			"end",
		Output: []Statement{
			FunctionDeclaration{
				Name: Identifier{Literal: "func1"},
				Parameters: []Variable{
					{Name: "@arg1"},
					{Name: "@arg2"},
				},
				Statements: []Statement{
					If{
						Condition: Comparison{
							LHS:      Variable{Name: "@var1"},
							RHS:      NewInteger(1),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewInteger(1)},
						},
					},
					If{
						Condition: Comparison{
							LHS:      Variable{Name: "@var1"},
							RHS:      NewInteger(1),
							Operator: "=",
						},
						Statements: []Statement{
							Print{Value: NewInteger(1)},
						},
						ElseIf: []ProcExpr{
							ElseIf{
								Condition: Comparison{
									LHS:      Variable{Name: "@var1"},
									RHS:      NewInteger(2),
									Operator: "=",
								},
								Statements: []Statement{
									Print{Value: NewInteger(2)},
								},
							},
							ElseIf{
								Condition: Comparison{
									LHS:      Variable{Name: "@var1"},
									RHS:      NewInteger(3),
									Operator: "=",
								},
								Statements: []Statement{
									Print{Value: NewInteger(3)},
								},
							},
						},
						Else: Else{
							Statements: []Statement{
								Print{Value: NewInteger(4)},
							},
						},
					},
					While{
						Condition: Ternary{literal: "true", value: ternary.TRUE},
						Statements: []Statement{
							FlowControl{Token: BREAK},
						},
					},
					While{
						Condition: Ternary{literal: "true", value: ternary.TRUE},
						Statements: []Statement{
							If{
								Condition: Comparison{
									LHS:      Variable{Name: "@var1"},
									RHS:      NewInteger(1),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: CONTINUE},
								},
							},
						},
					},
					While{
						Condition: Ternary{literal: "true", value: ternary.TRUE},
						Statements: []Statement{
							If{
								Condition: Comparison{
									LHS:      Variable{Name: "@var1"},
									RHS:      NewInteger(1),
									Operator: "=",
								},
								Statements: []Statement{
									FlowControl{Token: CONTINUE},
								},
								ElseIf: []ProcExpr{
									ElseIf{
										Condition: Comparison{
											LHS:      Variable{Name: "@var1"},
											RHS:      NewInteger(2),
											Operator: "=",
										},
										Statements: []Statement{
											FlowControl{Token: BREAK},
										},
									},
									ElseIf{
										Condition: Comparison{
											LHS:      Variable{Name: "@var1"},
											RHS:      NewInteger(3),
											Operator: "=",
										},
										Statements: []Statement{
											Return{Value: NewNull()},
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
					Return{
						Value: Variable{Name: "@var1"},
					},
				},
			},
		},
	},
	{
		Input: "select 'literal not terminated",
		Error: "literal not terminated [L:1 C:30]",
	},
}

func TestParse(t *testing.T) {
	for _, v := range parseTests {
		prog, err := Parse(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
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

				if !reflect.DeepEqual(parsedStmt.OrderByClause, expectStmt.OrderByClause) {
					t.Errorf("order by clause = %#v, want %#v for %q", parsedStmt.OrderByClause, expectStmt.OrderByClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.LimitClause, expectStmt.LimitClause) {
					t.Errorf("limit clause = %#v, want %#v for %q", parsedStmt.LimitClause, expectStmt.LimitClause, v.Input)
				}
			default:
				if !reflect.DeepEqual(stmt, expect) {
					t.Errorf("output = %#v, want %#v for %q", stmt, expect, v.Input)
				}
			}
		}
	}
}
