package parser

import (
	"fmt"
	"github.com/mithrandie/csvq/lib/ternary"
	"reflect"
	"testing"
)

var parseTests = []struct {
	Input  string
	Output []Statement
}{
	{
		Input: "select foo; select bar;",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: Identifier{Literal: "foo"}}}},
			},
			SelectQuery{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: Identifier{Literal: "bar"}}}},
			},
		},
	},
	{
		Input: "select 1 as a from dual",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{
							Object: NewInteger(1),
							As:     Token{Token: AS, Literal: "as"},
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
	{
		Input: "select c1 from stdin",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{
							Object: Identifier{Literal: "c1"},
						},
					},
				},
				FromClause: FromClause{From: "from", Tables: []Expression{
					Table{Object: Stdin{Stdin: "stdin"}},
				}},
			},
		},
	},
	{
		Input: "select 1 from table1, (select 2 from dual)",
		Output: []Statement{
			SelectQuery{
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
	{
		Input: "select 1 from table1 alias, (select 2 from dual) alias2",
		Output: []Statement{
			SelectQuery{
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
									SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("2")}}},
									FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
								},
							},
							Alias: Identifier{Literal: "alias2"},
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
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Identifier{Literal: "table1"},
							As:     Token{Token: AS, Literal: "as"},
							Alias:  Identifier{Literal: "alias"},
						},
						Table{
							Object: Subquery{
								Query: SelectQuery{
									SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("2")}}},
									FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
								},
							},
							As:    Token{Token: AS, Literal: "as"},
							Alias: Identifier{Literal: "alias2"},
						},
					},
				},
			},
		},
	},
	{
		Input: "select 1 " +
			" from dual " +
			" where 1 = 1" +
			" group by column1, column2 " +
			" having 1 > 1 " +
			" order by column4, column5 desc, column6 asc " +
			" limit 10 ",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause:   FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
				WhereClause: WhereClause{
					Where: "where",
					Filter: Comparison{
						LHS:      NewIntegerFromString("1"),
						Operator: Token{Token: COMPARISON_OP, Literal: "="},
						RHS:      NewIntegerFromString("1"),
					},
				},
				GroupByClause: GroupByClause{
					GroupBy: "group by",
					Items: []Expression{
						Identifier{Literal: "column1"},
						Identifier{Literal: "column2"},
					},
				},
				HavingClause: HavingClause{
					Having: "having",
					Filter: Comparison{
						LHS:      NewIntegerFromString("1"),
						Operator: Token{Token: COMPARISON_OP, Literal: ">"},
						RHS:      NewIntegerFromString("1"),
					},
				},
				OrderByClause: OrderByClause{
					OrderBy: "order by",
					Items: []Expression{
						OrderItem{Item: Identifier{Literal: "column4"}},
						OrderItem{Item: Identifier{Literal: "column5"}, Direction: Token{Token: DESC, Literal: "desc"}},
						OrderItem{Item: Identifier{Literal: "column6"}, Direction: Token{Token: ASC, Literal: "asc"}},
					},
				},
				LimitClause: LimitClause{
					Limit:  "limit",
					Number: 10,
				},
			},
		},
	},
	{
		Input: "select distinct * from dual",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select:   "select",
					Distinct: Token{Token: DISTINCT, Literal: "distinct"},
					Fields: []Expression{
						Field{Object: AllColumns{}},
					},
				},
				FromClause: FromClause{From: "from", Tables: []Expression{Table{Object: Dual{Dual: "dual"}}}},
			},
		},
	},
	{
		Input: "select ident, 'foo', 1, -1, 1.234, -1.234, true, '2010-01-01 12:00:00', null, ('bar') from dual",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Identifier{Literal: "ident"}},
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
	{
		Input: "select ident || 'foo' || 'bar'",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Concat{Items: []Expression{
							Identifier{Literal: "ident"},
							NewString("foo"),
							NewString("bar"),
						}}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 = 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Comparison{
							LHS:      Identifier{Literal: "column1"},
							Operator: Token{Token: COMPARISON_OP, Literal: "="},
							RHS:      NewInteger(1),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 < 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Comparison{
							LHS:      Identifier{Literal: "column1"},
							Operator: Token{Token: COMPARISON_OP, Literal: "<"},
							RHS:      NewInteger(1),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 is not null",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Is{
							Is:       "is",
							LHS:      Identifier{Literal: "column1"},
							RHS:      NewNullFromString("null"),
							Negation: Token{Token: NOT, Literal: "not"},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 is true",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Is{
							Is:  "is",
							LHS: Identifier{Literal: "column1"},
							RHS: NewTernaryFromString("true"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not between -10 and 10",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Between{
							Between:  "between",
							And:      "and",
							LHS:      Identifier{Literal: "column1"},
							Low:      NewIntegerFromString("-10"),
							High:     NewIntegerFromString("10"),
							Negation: Token{Token: NOT, Literal: "not"},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not in (1, 2, 3)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: In{
							In:  "in",
							LHS: Identifier{Literal: "column1"},
							List: []Expression{
								NewIntegerFromString("1"),
								NewIntegerFromString("2"),
								NewIntegerFromString("3"),
							},
							Negation: Token{Token: NOT, Literal: "not"},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 in (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: In{
							In:  "in",
							LHS: Identifier{Literal: "column1"},
							Query: Subquery{
								Query: SelectQuery{
									SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 not like 'pattern'",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Like{
							Like:     "like",
							LHS:      Identifier{Literal: "column1"},
							Pattern:  String{literal: "pattern"},
							Negation: Token{Token: NOT, Literal: "not"},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 = any (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Any{
							Any:      "any",
							LHS:      Identifier{Literal: "column1"},
							Operator: Token{Token: COMPARISON_OP, Literal: "="},
							Query: Subquery{
								Query: SelectQuery{
									SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 = all (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: All{
							All:      "all",
							LHS:      Identifier{Literal: "column1"},
							Operator: Token{Token: COMPARISON_OP, Literal: "="},
							Query: Subquery{
								Query: SelectQuery{
									SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select exists (select 1)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Exists{
							Exists: "exists",
							Query: Subquery{
								Query: SelectQuery{
									SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 + 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Arithmetic{
							LHS:      Identifier{Literal: "column1"},
							Operator: int('+'),
							RHS:      NewIntegerFromString("1"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 - 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Arithmetic{
							LHS:      Identifier{Literal: "column1"},
							Operator: int('-'),
							RHS:      NewIntegerFromString("1"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 * 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Arithmetic{
							LHS:      Identifier{Literal: "column1"},
							Operator: int('*'),
							RHS:      NewIntegerFromString("1"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 / 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Arithmetic{
							LHS:      Identifier{Literal: "column1"},
							Operator: int('/'),
							RHS:      NewIntegerFromString("1"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select column1 % 1",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Arithmetic{
							LHS:      Identifier{Literal: "column1"},
							Operator: int('%'),
							RHS:      NewIntegerFromString("1"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select true and false",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Logic{
							LHS:      NewTernaryFromString("true"),
							Operator: Token{Token: AND, Literal: "and"},
							RHS:      NewTernaryFromString("false"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select true or false",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Logic{
							LHS:      NewTernaryFromString("true"),
							Operator: Token{Token: OR, Literal: "or"},
							RHS:      NewTernaryFromString("false"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select not false",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Logic{
							Operator: Token{Token: NOT, Literal: "not"},
							RHS:      NewTernaryFromString("false"),
						}},
					},
				},
			},
		},
	},
	{
		Input: "select true or (false and false)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Logic{
							LHS:      NewTernaryFromString("true"),
							Operator: Token{Token: OR, Literal: "or"},
							RHS: Parentheses{
								Expr: Logic{
									LHS:      NewTernaryFromString("false"),
									Operator: Token{Token: AND, Literal: "and"},
									RHS:      NewTernaryFromString("false"),
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select true and true or false and not false",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Logic{
							LHS: Logic{
								LHS:      NewTernaryFromString("true"),
								Operator: Token{Token: AND, Literal: "and"},
								RHS:      NewTernaryFromString("true"),
							},
							Operator: Token{Token: OR, Literal: "or"},
							RHS: Logic{
								LHS:      NewTernaryFromString("false"),
								Operator: Token{Token: AND, Literal: "and"},
								RHS: Logic{
									Operator: Token{Token: NOT, Literal: "not"},
									RHS:      NewTernaryFromString("false"),
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select @var",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Variable{Name: "@var"}},
					},
				},
			},
		},
	},
	{
		Input: "select @var := 1",
		Output: []Statement{
			SelectQuery{
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
	{
		Input: "select case when true then 'A' when false then 'B' end",
		Output: []Statement{
			SelectQuery{
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
	{
		Input: "select case column1 when 1 then 'A' when 2 then 'B' else 'C' end",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Case{
							Case:  "case",
							End:   "end",
							Value: Identifier{Literal: "column1"},
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
	{
		Input: "select count()",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Function{
							Name:   "count",
							Option: Option{},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select count(distinct *)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Function{
							Name: "count",
							Option: Option{
								Distinct: Token{Token: DISTINCT, Literal: "distinct"},
								Args:     []Expression{AllColumns{}},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select count(column1, column2)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: Function{
							Name: "count",
							Option: Option{
								Args: []Expression{
									Identifier{Literal: "column1"},
									Identifier{Literal: "column2"},
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select group_concat(column1 order by column1)",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: GroupConcat{
							GroupConcat: "group_concat",
							Option:      Option{Args: []Expression{Identifier{Literal: "column1"}}},
							OrderBy: OrderByClause{
								OrderBy: "order by",
								Items: []Expression{
									OrderItem{Item: Identifier{Literal: "column1"}},
								},
							},
						}},
					},
				},
			},
		},
	},
	{
		Input: "select group_concat(column1 separator ',')",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{
					Select: "select",
					Fields: []Expression{
						Field{Object: GroupConcat{
							GroupConcat:  "group_concat",
							Option:       Option{Args: []Expression{Identifier{Literal: "column1"}}},
							SeparatorLit: "separator",
							Separator:    ",",
						}},
					},
				},
			},
		},
	},
	{
		Input: "select 1 from table1 cross join table2",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								JoinType:  Token{Token: CROSS, Literal: "cross"},
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
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								JoinType:  Token{Token: INNER, Literal: "inner"},
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
										LHS:      Identifier{Literal: "table1.id"},
										Operator: Token{Token: COMPARISON_OP, Literal: "="},
										RHS:      Identifier{Literal: "table2.id"},
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
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								Natural:   Token{Token: NATURAL, Literal: "natural"},
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
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								Direction: Token{Token: LEFT, Literal: "left"},
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
	{
		Input: "select 1 from table1 natural outer join table2",
		Output: []Statement{
			SelectQuery{
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								Natural:   Token{Token: NATURAL, Literal: "natural"},
								JoinType:  Token{Token: OUTER, Literal: "outer"},
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
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								Direction: Token{Token: RIGHT, Literal: "right"},
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
				SelectClause: SelectClause{Select: "select", Fields: []Expression{Field{Object: NewIntegerFromString("1")}}},
				FromClause: FromClause{
					From: "from",
					Tables: []Expression{
						Table{
							Object: Join{
								Join:      "join",
								Table:     Table{Object: Identifier{Literal: "table1"}},
								JoinTable: Table{Object: Identifier{Literal: "table2"}},
								Direction: Token{Token: FULL, Literal: "full"},
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
				Var: "var",
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
		Input: "insert into table1 values (1, 'str1'), (2, 'str2')",
		Output: []Statement{
			InsertQuery{
				Insert: "insert",
				Into:   "into",
				Table:  Identifier{Literal: "table1"},
				Values: "values",
				ValuesList: []Expression{
					InsertValues{
						Values: []Expression{
							NewInteger(1),
							NewString("str1"),
						},
					},
					InsertValues{
						Values: []Expression{
							NewInteger(2),
							NewString("str2"),
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
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
				Values: "values",
				ValuesList: []Expression{
					InsertValues{
						Values: []Expression{
							NewInteger(1),
							NewString("str1"),
						},
					},
					InsertValues{
						Values: []Expression{
							NewInteger(2),
							NewString("str2"),
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
	{
		Input: "insert into table1 (column1, column2) select 1, 2",
		Output: []Statement{
			InsertQuery{
				Insert: "insert",
				Into:   "into",
				Table:  Identifier{Literal: "table1"},
				Fields: []Expression{
					Identifier{Literal: "column1"},
					Identifier{Literal: "column2"},
				},
				Query: SelectQuery{
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
	{
		Input: "update table1 set column1 = 1, column2 = 2 from table1 where true",
		Output: []Statement{
			UpdateQuery{
				Update: "update",
				Tables: []Expression{
					Table{Object: Identifier{Literal: "table1"}},
				},
				Set: "set",
				SetList: []Expression{
					UpdateSet{Field: Identifier{Literal: "column1"}, Value: NewInteger(1)},
					UpdateSet{Field: Identifier{Literal: "column2"}, Value: NewInteger(2)},
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
		Input: "delete from table1",
		Output: []Statement{
			DeleteQuery{
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
					Position: Token{Token: FIRST, Literal: "first"},
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
					Position: Token{Token: LAST, Literal: "last"},
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
					Position: Token{Token: AFTER, Literal: "after"},
					Column:   Identifier{Literal: "column2"},
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
					Position: Token{Token: BEFORE, Literal: "before"},
					Column:   Identifier{Literal: "column2"},
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
				Columns:    []Expression{Identifier{Literal: "column1"}},
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
				Columns:    []Expression{Identifier{Literal: "column1"}, Identifier{Literal: "column2"}},
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
				Old:        Identifier{Literal: "column1"},
				To:         "to",
				New:        Identifier{Literal: "column2"},
			},
		},
	},
	{
		Input: "commit",
		Output: []Statement{
			Commit{
				Literal: "commit",
			},
		},
	},
	{
		Input: "rollback",
		Output: []Statement{
			Rollback{
				Literal: "rollback",
			},
		},
	},
	{
		Input: "print 'foo'",
		Output: []Statement{
			Print{
				Print: "print",
				Value: NewString("foo"),
			},
		},
	},
}

func TestParse(t *testing.T) {
	SetDebugLevel(0, true)

	errorQuery := "select 'literal not teriinated"
	_, err := Parse(errorQuery)
	if err == nil {
		t.Errorf("no error, want an error for %q", errorQuery)
	}

	for _, v := range parseTests {
		prog, err := Parse(v.Input)
		if err != nil {
			t.Errorf("unexpected error %q at %q", err.Error(), v.Input)
			return
		}

		if len(v.Output) != len(prog) {
			t.Errorf("parsed program has %d statement(s), want %d statement(s) for %q", len(prog), len(v.Output), v.Input)
			return
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

				if !reflect.DeepEqual(parsedStmt.SelectClause, expectStmt.SelectClause) {
					t.Errorf("select clause = %#v, want %#v for %q", parsedStmt.SelectClause, expectStmt.SelectClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.FromClause, expectStmt.FromClause) {
					t.Errorf("from clause = %#v, want %#v for %q", parsedStmt.FromClause, expectStmt.FromClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.WhereClause, expectStmt.WhereClause) {
					t.Errorf("where clause = %#v, want %#v for %q", parsedStmt.WhereClause, expectStmt.WhereClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.GroupByClause, expectStmt.GroupByClause) {
					t.Errorf("group by clause = %#v, want %#v for %q", parsedStmt.GroupByClause, expectStmt.GroupByClause, v.Input)
				}
				if !reflect.DeepEqual(parsedStmt.HavingClause, expectStmt.HavingClause) {
					t.Errorf("having clause = %#v, want %#v for %q", parsedStmt.HavingClause, expectStmt.HavingClause, v.Input)
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

func ExampleSetDebugLevel() {
	SetDebugLevel(0, false)
	_, err := Parse("select select")
	fmt.Println(err)

	SetDebugLevel(0, true)
	_, err = Parse("select select")
	fmt.Println(err)
	//Output:
	//syntax error
	//syntax error: unexpected SELECT
}
