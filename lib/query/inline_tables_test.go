package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var inlineTableNodesSetTests = []struct {
	Name   string
	Expr   parser.InlineTable
	Result InlineTableNodes
	Error  string
}{
	{
		Name: "InlineTableNodes Set",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
				parser.Identifier{Literal: "num"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
							parser.Field{Object: parser.NewIntegerValue(1)},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Result: InlineTableNodes{
			InlineTableMap{
				"IT": &View{
					Header: NewHeader("it", []string{"c1", "c2", "num"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
							value.NewInteger(1),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
							value.NewInteger(1),
						}),
						NewRecord([]value.Primary{
							value.NewString("3"),
							value.NewString("str3"),
							value.NewInteger(1),
						}),
					},
					Tx: TestTx,
				},
			},
			InlineTableMap{
				"IT2": &View{
					Header: NewHeader("it2", []string{"c1", "c2", "num"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
							value.NewInteger(1),
						}),
					},
					Tx: TestTx,
				},
			},
		},
	},
}

func TestInlineTableNodes_Set(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	list := InlineTableNodes{
		{},
		InlineTableMap{
			"IT2": &View{
				Header: NewHeader("it2", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
				},
				Tx: TestTx,
			},
		},
	}

	_ = TestTx.Flags.SetRepository(TestDataDir)

	for _, v := range inlineTableNodesSetTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		err := list.Set(context.Background(), NewFilter(TestTx), v.Expr)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, list, v.Result)
		}
	}
}

var inlineTableNodesGetTests = []struct {
	Name      string
	TableName parser.Identifier
	Result    *View
	Error     string
}{
	{
		Name:      "InlineTableNodes Get",
		TableName: parser.Identifier{Literal: "it2"},
		Result: &View{
			Header: NewHeader("it2", []string{"c1", "c2", "num"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
				}),
			},
		},
	},
	{
		Name:      "InlineTableNodes Get Undefined Error",
		TableName: parser.Identifier{Literal: "notexist"},
		Error:     "inline table notexist is undefined",
	},
}

func TestInlineTableNodes_Get(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
	}()

	list := InlineTableNodes{
		InlineTableMap{
			"IT": &View{
				Header: NewHeader("it", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
						value.NewInteger(1),
					}),
				},
			},
		},
		InlineTableMap{
			"IT2": &View{
				Header: NewHeader("it2", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
				},
			},
		},
	}

	for _, v := range inlineTableNodesGetTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		view, err := list.Get(v.TableName)
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
		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, view, v.Result)
		}
	}
}

var inlineTableNodesLoadTests = []struct {
	Name   string
	Expr   parser.WithClause
	Result InlineTableNodes
	Error  string
}{
	{
		Name: "InlineTableNodes Load",
		Expr: parser.WithClause{
			With: "with",
			InlineTables: []parser.QueryExpression{
				parser.InlineTable{
					Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
					Name:      parser.Identifier{Literal: "it_recursive"},
					Fields: []parser.QueryExpression{
						parser.Identifier{Literal: "n"},
					},
					As: "as",
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectSet{
							LHS: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.NewIntegerValueFromString("1")},
									},
								},
							},
							Operator: parser.Token{Token: parser.UNION, Literal: "union"},
							RHS: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Select: "select",
									Fields: []parser.QueryExpression{
										parser.Field{
											Object: parser.Arithmetic{
												LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
												RHS:      parser.NewIntegerValueFromString("1"),
												Operator: '+',
											},
										},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "it_recursive"}, Alias: parser.Identifier{Literal: "t"}},
									},
								},
								WhereClause: parser.WhereClause{
									Filter: parser.Comparison{
										LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
										RHS:      parser.NewIntegerValueFromString("3"),
										Operator: "<",
									},
								},
							},
						},
					},
				},
			},
		},
		Result: InlineTableNodes{
			InlineTableMap{
				"IT": &View{
					Header: NewHeader("it", []string{"c1", "c2", "num"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
							value.NewInteger(1),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
							value.NewInteger(1),
						}),
						NewRecord([]value.Primary{
							value.NewString("3"),
							value.NewString("str3"),
							value.NewInteger(1),
						}),
					},
					Tx: TestTx,
				},
				"IT_RECURSIVE": &View{
					Header: []HeaderField{
						{
							View:        "it_recursive",
							Column:      "n",
							Number:      1,
							IsFromTable: true,
						},
					},
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewInteger(1),
						}),
						NewRecord([]value.Primary{
							value.NewInteger(2),
						}),
						NewRecord([]value.Primary{
							value.NewInteger(3),
						}),
					},
					Tx: TestTx,
				},
			},
			InlineTableMap{
				"IT2": &View{
					Header: NewHeader("it2", []string{"c1", "c2", "num"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
							value.NewInteger(1),
						}),
					},
					Tx: TestTx,
				},
			},
		},
	},
	{
		Name: "InlineTableNodes Load Set Error",
		Expr: parser.WithClause{
			With: "with",
			InlineTables: []parser.QueryExpression{
				parser.InlineTable{
					Name: parser.Identifier{Literal: "it"},
					Fields: []parser.QueryExpression{
						parser.Identifier{Literal: "c1"},
						parser.Identifier{Literal: "c2"},
					},
					As: "as",
					Query: parser.SelectQuery{
						SelectEntity: parser.SelectEntity{
							SelectClause: parser.SelectClause{
								Select: "select",
								Fields: []parser.QueryExpression{
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
									parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
								},
							},
							FromClause: parser.FromClause{
								From: "from",
								Tables: []parser.QueryExpression{
									parser.Table{Object: parser.Identifier{Literal: "table1"}},
								},
							},
						},
					},
				},
			},
		},
		Error: "inline table it is redefined",
	},
}

func TestInlineTableNodes_Load(t *testing.T) {
	list := InlineTableNodes{
		InlineTableMap{
			"IT": &View{
				Header: NewHeader("it", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
						value.NewInteger(1),
					}),
				},
				Tx: TestTx,
			},
		},
		InlineTableMap{
			"IT2": &View{
				Header: NewHeader("it2", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
				},
				Tx: TestTx,
			},
		},
	}

	for _, v := range inlineTableNodesLoadTests {
		err := list.Load(context.Background(), NewFilter(TestTx), v.Expr)
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, list, v.Result)
		}
	}
}

var inlineTableMapSetTests = []struct {
	Name   string
	Expr   parser.InlineTable
	Result InlineTableMap
	Error  string
}{
	{
		Name: "InlineTableMap Set",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
				parser.Identifier{Literal: "num"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
							parser.Field{Object: parser.NewIntegerValueFromString("1")},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Result: InlineTableMap{
			"IT": &View{
				Header: NewHeader("it", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
						value.NewInteger(1),
					}),
				},
				Tx: TestTx,
			},
		},
	},
	{
		Name: "InlineTableMap Set Recursive Table",
		Expr: parser.InlineTable{
			Recursive: parser.Token{Token: parser.RECURSIVE, Literal: "recursive"},
			Name:      parser.Identifier{Literal: "it_recursive"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "n"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectSet{
					LHS: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{Object: parser.NewIntegerValueFromString("1")},
							},
						},
					},
					Operator: parser.Token{Token: parser.UNION, Literal: "union"},
					RHS: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							Select: "select",
							Fields: []parser.QueryExpression{
								parser.Field{
									Object: parser.Arithmetic{
										LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
										RHS:      parser.NewIntegerValueFromString("1"),
										Operator: '+',
									},
								},
							},
						},
						FromClause: parser.FromClause{
							Tables: []parser.QueryExpression{
								parser.Table{Object: parser.Identifier{Literal: "it_recursive"}},
							},
						},
						WhereClause: parser.WhereClause{
							Filter: parser.Comparison{
								LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "n"}},
								RHS:      parser.NewIntegerValueFromString("3"),
								Operator: "<",
							},
						},
					},
				},
			},
		},
		Result: InlineTableMap{
			"IT": &View{
				Header: NewHeader("it", []string{"c1", "c2", "num"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewString("3"),
						value.NewString("str3"),
						value.NewInteger(1),
					}),
				},
				Tx: TestTx,
			},
			"IT_RECURSIVE": &View{
				Header: []HeaderField{
					{
						View:        "it_recursive",
						Column:      "n",
						Number:      1,
						IsFromTable: true,
					},
				},
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewInteger(1),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(2),
					}),
					NewRecord([]value.Primary{
						value.NewInteger(3),
					}),
				},
				Tx: TestTx,
			},
		},
	},
	{
		Name: "InlineTableMap Set Redefined Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "inline table it is redefined",
	},
	{
		Name: "InlineTableMap Set Query Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it2"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c2"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "InlineTableMap Set Field Length Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it2"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "c1"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "select query should return exactly 1 field for inline table it2",
	},
	{
		Name: "InlineTableMap Set Duplicate Field Name Error",
		Expr: parser.InlineTable{
			Name: parser.Identifier{Literal: "it2"},
			Fields: []parser.QueryExpression{
				parser.Identifier{Literal: "c1"},
				parser.Identifier{Literal: "c1"},
			},
			As: "as",
			Query: parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						Select: "select",
						Fields: []parser.QueryExpression{
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
							parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
						},
					},
					FromClause: parser.FromClause{
						From: "from",
						Tables: []parser.QueryExpression{
							parser.Table{Object: parser.Identifier{Literal: "table1"}},
						},
					},
				},
			},
		},
		Error: "field name c1 is a duplicate",
	},
}

func TestInlineTableMap_Set(t *testing.T) {
	defer func() {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDataDir

	it := InlineTableMap{}

	for _, v := range inlineTableMapSetTests {
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		err := it.Set(context.Background(), NewFilter(TestTx), v.Expr)
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
			t.Errorf("%s: result = %v, want %v", v.Name, it, v.Result)
		}
	}
}

var inlineTableMapGetTests = []struct {
	Name      string
	TableName parser.Identifier
	Result    *View
	Error     string
}{
	{
		Name:      "InlineTableMap Get",
		TableName: parser.Identifier{Literal: "it"},
		Result: &View{
			Header: NewHeader("it", []string{"c1", "c2", "num"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewInteger(1),
				}),
			},
		},
	},
	{
		Name:      "InlineTableMap Get Undefined Error",
		TableName: parser.Identifier{Literal: "notexist"},
		Error:     "inline table notexist is undefined",
	},
}

func TestInlineTableMap_Get(t *testing.T) {
	it := InlineTableMap{
		"IT": &View{
			Header: NewHeader("it", []string{"c1", "c2", "num"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewInteger(1),
				}),
			},
		},
	}

	for _, v := range inlineTableMapGetTests {
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
			t.Errorf("%s: result = %v, want %v", v.Name, ret, v.Result)
		}
	}
}
