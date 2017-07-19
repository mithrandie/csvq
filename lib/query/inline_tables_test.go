package query

import (
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"reflect"
)

var inlineTablesSetTests = []struct {
	Name   string
	Expr   parser.InlineTable
	Result InlineTables
	Error  string
}{
	{
		Name: "InlineTables Set",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
				parser.Identifier{Literal: "num"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
							parser.Field{Object: parser.NewInteger(1)},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Result: InlineTables{
			"IT": &View{
				Header: NewHeaderWithoutId("it", []string{"c1", "c2", "num"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str3"),
						parser.NewInteger(1),
					}),
				},
			},
		},
	},
	{
		Name: "InlineTables Set Recursive Table",
		Expr: parser.InlineTable{
			Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
			Name:      parser.Identifier{Literal: "it_recursive"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "n"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectSet{
					LHS: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{Object: parser.NewInteger(1)},
							},
						},
					},
					Operator: parser.Token{Token: parser.UNION, Literal: "union"},
					RHS: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.Expression{
								parser.Field{
									Object: parser.Arithmetic{
										LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
										RHS:      parser.NewInteger(1),
										Operator: '+',
									},
								},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.Expression{
								parser.Table{Object: parser.Identifier{Literal: "it_recursive"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
								RHS:      parser.NewInteger(3),
								Operator: "<",
							},
						},
					},
				},
			},
		},
		Result: InlineTables{
			"IT": &View{
				Header: NewHeaderWithoutId("it", []string{"c1", "c2", "num"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str3"),
						parser.NewInteger(1),
					}),
				},
			},
			"IT_RECURSIVE": &View{
				Header: []HeaderField{
					{
						Reference: "it_recursive",
						Column:    "n",
						FromTable: true,
					},
				},
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(2),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(3),
					}),
				},
			},
		},
	},
	{
		Name: "InlineTables Set Duplicate Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "inline table it already exists",
	},
	{
		Name: "InlineTables Set Query Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it2"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "InlineTables Set Field Length Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it2"},
			Columns: []parser.Expression{
				parser.Identifier{Literal: "c1"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.Expression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.Expression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "view it2: field length does not match",
	},
}

func TestInlineTables_Set(t *testing.T) {
	it := InlineTables{}

	for _, v := range inlineTablesSetTests {
		err := it.Set(v.Expr, NewFilter([]Variables{{}}))
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
		if !reflect.DeepEqual(it, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, it, v.Result)
		}
	}
}

var inlineTablesGetTests = []struct {
	Name      string
	TableName string
	Result    *View
	Error     string
}{
	{
		Name:      "InlineTables Get",
		TableName: "it",
		Result: &View{
			Header: NewHeaderWithoutId("it", []string{"c1", "c2", "num"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewInteger(1),
				}),
			},
		},
	},
	{
		Name:      "InlineTables Get Not Exist Error",
		TableName: "notexist",
		Error:     "inline table notexist does not exist",
	},
}

func TestInlineTables_Get(t *testing.T) {
	it := InlineTables{
		"IT": &View{
			Header: NewHeaderWithoutId("it", []string{"c1", "c2", "num"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewInteger(1),
				}),
			},
		},
	}

	for _, v := range inlineTablesGetTests {
		ret, err := it.Get(v.TableName)
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
		if !reflect.DeepEqual(ret, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, ret, v.Result)
		}
	}
}

var inlineTablesLoadTests = []struct {
	Name   string
	Expr   parser.WithClause
	Result InlineTables
	Error  string
}{
	{
		Name: "InlineTables Load",
		Expr: parser.WithClause{
			With: "with",
			InlineTables: []parser.Expression{
				parser.InlineTable{
					Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
					Name:      parser.Identifier{Literal: "it_recursive"},
					Columns: []parser.Expression{
						parser.Identifier{Literal: "n"},
					},
					As: "as",
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectSet{
							LHS: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.Expression{
										parser.Field{Object: parser.NewInteger(1)},
									},
								},
							},
							Operator: parser.Token{Token: parser.UNION, Literal: "union"},
							RHS: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.Expression{
										parser.Field{
											Object: parser.Arithmetic{
												LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
												RHS:      parser.NewInteger(1),
												Operator: '+',
											},
										},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.Expression{
										parser.Table{Object: parser.Identifier{Literal: "it_recursive"}},
									},
								},
								WhereClause: parser.WhereClause{
									Filter: parser.Comparison{
										LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
										RHS:      parser.NewInteger(3),
										Operator: "<",
									},
								},
							},
						},
					},
				},
			},
		},
		Result: InlineTables{
			"IT": &View{
				Header: NewHeaderWithoutId("it", []string{"c1", "c2", "num"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str3"),
						parser.NewInteger(1),
					}),
				},
			},
			"IT_RECURSIVE": &View{
				Header: []HeaderField{
					{
						Reference: "it_recursive",
						Column:    "n",
						FromTable: true,
					},
				},
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(2),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewInteger(3),
					}),
				},
			},
		},
	},
	{
		Name: "InlineTables Load Set Error",
		Expr: parser.WithClause{
			With: "with",
			InlineTables: []parser.Expression{
				parser.InlineTable{
					Name: parser.Identifier{Literal: "it"},
					Columns: []parser.Expression{
						parser.Identifier{Literal: "c1"},
						parser.Identifier{Literal: "c2"},
					},
					As: "as",
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.Expression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								From: "from",
								Tables: []parser.Expression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
		},
		Error: "inline table it already exists",
	},
}

func TestInlineTables_Load(t *testing.T) {
	for _, v := range inlineTablesLoadTests {
		it := InlineTables{
			"IT": &View{
				Header: NewHeaderWithoutId("it", []string{"c1", "c2", "num"}),
				Records: []Record{
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
						parser.NewInteger(1),
					}),
					NewRecordWithoutId([]parser.Primary{
						parser.NewString("3"),
						parser.NewString("str3"),
						parser.NewInteger(1),
					}),
				},
			},
		}

		err := it.Load(v.Expr, NewFilter([]Variables{{}}))
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
		if !reflect.DeepEqual(it, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, it, v.Result)
		}
	}
}
