package query

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/mithrandie/go-text"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestFieldIndexCache_Get(t *testing.T) {
	cache := &FieldIndexCache{
		limitToUseSlice: 2,
		m:               nil,
		exprs:           []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
		indices:         []int{1},
	}

	_, ok := cache.Get(parser.FieldReference{Column: parser.Identifier{Literal: "c9"}})
	if ok {
		t.Error("Get() is succeeded, want to be failed")
	}

	expect := 1
	idx, ok := cache.Get(parser.FieldReference{Column: parser.Identifier{Literal: "c1"}})
	if !ok {
		t.Error("Get() is failed, want to be succeeded")
	}
	if idx != expect {
		t.Errorf("result = %d, want %d", idx, expect)
	}

	cache = &FieldIndexCache{
		limitToUseSlice: 2,
		m: map[parser.QueryExpression]int{
			parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}: 1,
			parser.FieldReference{Column: parser.Identifier{Literal: "c2"}}: 2,
			parser.FieldReference{Column: parser.Identifier{Literal: "c3"}}: 3,
		},
		exprs:   nil,
		indices: nil,
	}

	_, ok = cache.Get(parser.FieldReference{Column: parser.Identifier{Literal: "c9"}})
	if ok {
		t.Error("Get() is succeeded, want to be failed")
	}

	expect = 2
	idx, ok = cache.Get(parser.FieldReference{Column: parser.Identifier{Literal: "c2"}})
	if !ok {
		t.Error("Get() is failed, want to be succeeded")
	}
	if idx != expect {
		t.Errorf("result = %d, want %d", idx, expect)
	}
}

func TestFieldIndexCache_Add(t *testing.T) {
	cache := NewFieldIndexCache(1, 2)

	cache.Add(parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}, 1)
	expect := &FieldIndexCache{
		limitToUseSlice: 2,
		m:               nil,
		exprs:           []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
		indices:         []int{1},
	}

	if !reflect.DeepEqual(cache, expect) {
		t.Errorf("cache = %v, want %v", cache, expect)
	}

	cache.Add(parser.FieldReference{Column: parser.Identifier{Literal: "c2"}}, 2)
	cache.Add(parser.FieldReference{Column: parser.Identifier{Literal: "c3"}}, 3)
	expect = &FieldIndexCache{
		limitToUseSlice: 2,
		m: map[parser.QueryExpression]int{
			parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}: 1,
			parser.FieldReference{Column: parser.Identifier{Literal: "c2"}}: 2,
			parser.FieldReference{Column: parser.Identifier{Literal: "c3"}}: 3,
		},
		exprs:   nil,
		indices: nil,
	}

	if !reflect.DeepEqual(cache, expect) {
		t.Errorf("cache = %v, want %v", cache, expect)
	}
}

var testVariablesReferenceScope = GenerateReferenceScope([]map[string]map[string]interface{}{
	{
		scopeNameVariables: {
			"var1": value.NewInteger(1),
		},
	},
	{
		scopeNameVariables: {
			"var1": value.NewInteger(2),
		},
	},
}, nil, time.Time{}, nil)

var testTemporaryTablesReferenceScope = GenerateReferenceScope([]map[string]map[string]interface{}{
	{
		scopeNameTempTables: {
			"TEMPTABLE1": &View{
				Header: NewHeader("temptable1", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "temptable1",
					Delimiter: ',',
					ViewType:  ViewTypeTemporaryTable,
				},
			},
		},
	},
	{
		scopeNameTempTables: {
			"TEMPTABLE2": &View{
				Header: NewHeader("temptable2", []string{"column1", "column2"}),
				RecordSet: []Record{
					NewRecord([]value.Primary{
						value.NewString("1"),
						value.NewString("str1"),
					}),
					NewRecord([]value.Primary{
						value.NewString("2"),
						value.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "temptable2",
					Delimiter: ',',
					ViewType:  ViewTypeTemporaryTable,
				},
			},
		},
	},
}, nil, time.Time{}, nil)

var referenceScopeGetVariableTests = []struct {
	Name   string
	Expr   parser.Variable
	Result value.Primary
	Error  string
}{
	{
		Name:   "ReferenceScope GetVariable",
		Expr:   parser.Variable{Name: "var1"},
		Result: value.NewInteger(1),
	},
	{
		Name:  "ReferenceScope GetVariable Undeclared Error",
		Expr:  parser.Variable{Name: "undef"},
		Error: "variable @undef is undeclared",
	},
}

func TestReferenceScope_GetVariable(t *testing.T) {

	for _, v := range referenceScopeGetVariableTests {
		result, err := testVariablesReferenceScope.GetVariable(v.Expr)
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
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, result, v.Result)
		}
	}
}

var referenceScopeSubstituteVariableTests = []struct {
	Name        string
	Expr        parser.VariableSubstitution
	ResultScope *ReferenceScope
	Result      value.Primary
	Error       string
}{
	{
		Name: "ReferenceScope SubstituteVariable",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var1"},
			Value:    parser.NewIntegerValue(3),
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameVariables: {
					"var1": value.NewInteger(3),
				},
			},
			{
				scopeNameVariables: {
					"var1": value.NewInteger(2),
				},
			},
		}, nil, time.Time{}, nil),
		Result: value.NewInteger(3),
	},
	{
		Name: "ReferenceScope SubstituteVariable Undeclared Error",
		Expr: parser.VariableSubstitution{
			Variable: parser.Variable{Name: "var2"},
			Value:    parser.NewIntegerValue(3),
		},
		Error: "variable @var2 is undeclared",
	},
}

func TestReferenceScope_SubstituteVariable(t *testing.T) {
	for _, v := range referenceScopeSubstituteVariableTests {
		result, err := testVariablesReferenceScope.SubstituteVariable(context.Background(), v.Expr)
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
		if !BlockScopeListEqual(testVariablesReferenceScope.Blocks, v.ResultScope.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, testVariablesReferenceScope.Blocks, v.ResultScope.Blocks)
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}
	}
}

var referenceScopeDisposeVariableTests = []struct {
	Name        string
	Expr        parser.Variable
	ResultScope *ReferenceScope
	Error       string
}{
	{
		Name: "ReferenceScope DisposeVariable",
		Expr: parser.Variable{Name: "var1"},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{},
			{
				scopeNameVariables: {
					"var1": value.NewInteger(2),
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:  "ReferenceScope DisposeVariable Undeclared Error",
		Expr:  parser.Variable{Name: "undef"},
		Error: "variable @undef is undeclared",
	},
}

func TestReferenceScope_DisposeVariable(t *testing.T) {
	for _, v := range referenceScopeDisposeVariableTests {
		err := testVariablesReferenceScope.DisposeVariable(v.Expr)
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
		if !BlockScopeListEqual(testVariablesReferenceScope.Blocks, v.ResultScope.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, testVariablesReferenceScope.Blocks, v.ResultScope.Blocks)
		}
	}
}

var referenceScopeTemporaryTableExistsTests = []struct {
	Name   string
	Path   string
	Result bool
}{
	{
		Name:   "ReferenceScope TemporaryTableExists",
		Path:   "temptable2",
		Result: true,
	},
	{
		Name:   "ReferenceScope TemporaryTableExists Not Exist",
		Path:   "notexist",
		Result: false,
	},
}

func TestReferenceScope_TemporaryTableExists(t *testing.T) {

	for _, v := range referenceScopeTemporaryTableExistsTests {
		result := testTemporaryTablesReferenceScope.TemporaryTableExists(v.Path)
		if result != v.Result {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Result)
		}
	}
}

var referenceScopeGetTemporaryTableTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "ReferenceScope GetTemporaryTable",
		Path: parser.Identifier{Literal: "temptable2"},
		Result: &View{
			Header: NewHeader("temptable2", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "temptable2",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
	},
	{
		Name:  "ReferenceScope GetTemporaryTable Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "view /path/to/table9.csv is undeclared",
	},
}

func TestReferenceScope_GetTemporaryTable(t *testing.T) {
	for _, v := range referenceScopeGetTemporaryTableTests {
		view, err := testTemporaryTablesReferenceScope.GetTemporaryTable(v.Path)
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
			t.Errorf("%s: view = %v, want %v", v.Name, view, v.Result)
		}
	}
}

var referenceScopeGetTemporaryTableWithInternalIdTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "ReferenceScope GetTemporaryTableWithInternalId",
		Path: parser.Identifier{Literal: "temptable2"},
		Result: &View{
			Header: NewHeaderWithId("temptable2", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(0, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "temptable2",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
	},
	{
		Name:  "ReferenceScope GetTemporaryTableWithInternalId Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "view /path/to/table9.csv is undeclared",
	},
}

func TestTemporaryViewScopes_GetWithInternalId(t *testing.T) {
	for _, v := range referenceScopeGetTemporaryTableWithInternalIdTests {
		view, err := testTemporaryTablesReferenceScope.GetTemporaryTableWithInternalId(context.Background(), v.Path, TestTx.Flags)
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
			t.Errorf("%s: view = %v, want %v", v.Name, view, v.Result)
		}
	}
}

var referenceScopeSetTemporaryTableTests = []struct {
	Name    string
	SetView *View
	Result  *ReferenceScope
}{
	{
		Name: "ReferenceScope SetTemporaryTable",
		SetView: &View{
			Header: NewHeader("tempview3", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "tempview3",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
		Result: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"TEMPTABLE1": &View{
						Header: NewHeader("temptable1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "temptable1",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
					"TEMPTABLE3": &View{
						Header: NewHeader("tempview3", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "tempview3",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
			{
				scopeNameTempTables: {
					"TEMPTABLE2": &View{
						Header: NewHeader("temptable2", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "temptable2",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
}

func TestReferenceScope_SetTemporaryTable(t *testing.T) {
	for _, v := range referenceScopeSetTemporaryTableTests {
		testTemporaryTablesReferenceScope.SetTemporaryTable(v.SetView)
		if !BlockScopeListEqual(testTemporaryTablesReferenceScope.Blocks, v.Result.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, testTemporaryTablesReferenceScope.Blocks, v.Result.Blocks)
		}
	}
}

var referenceScopeReplaceTemporaryTableTests = []struct {
	Name    string
	SetView *View
	Result  *ReferenceScope
	Error   string
}{
	{
		Name: "ReferenceScope ReplaceTemporaryTable",
		SetView: &View{
			Header: NewHeader("temptable2", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("updated"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("updated"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "temptable2",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
		Result: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"TEMPTABLE1": &View{
						Header: NewHeader("temptable1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "temptable1",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
					"TEMPTABLE3": &View{
						Header: NewHeader("tempview3", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "tempview3",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
			{
				scopeNameTempTables: {
					"TEMPTABLE2": &View{
						Header: NewHeader("temptable2", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("updated"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("updated"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "temptable2",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
}

func TestReferenceScope_ReplaceTemporaryTable(t *testing.T) {
	for _, v := range referenceScopeReplaceTemporaryTableTests {
		testTemporaryTablesReferenceScope.ReplaceTemporaryTable(v.SetView)
		if !BlockScopeListEqual(testTemporaryTablesReferenceScope.Blocks, v.Result.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, testTemporaryTablesReferenceScope.Blocks, v.Result.Blocks)
		}
	}
}

var referenceScopesDisposeTemporaryTableTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *ReferenceScope
	Error  string
}{
	{
		Name: "ReferenceScope DisposeTemporaryTable",
		Path: parser.Identifier{Literal: "temptable1"},
		Result: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"TEMPTABLE3": &View{
						Header: NewHeader("tempview3", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "tempview3",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
			{
				scopeNameTempTables: {
					"TEMPTABLE2": &View{
						Header: NewHeader("temptable2", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("updated"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("updated"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      "temptable2",
							Delimiter: ',',
							ViewType:  ViewTypeTemporaryTable,
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:  "ReferenceScope DisposeTemporaryTable Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "view /path/to/table9.csv is undeclared",
	},
}

func TestReferenceScope_DisposeTemporaryTable(t *testing.T) {
	for _, v := range referenceScopesDisposeTemporaryTableTests {
		err := testTemporaryTablesReferenceScope.DisposeTemporaryTable(v.Path)
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
		if !BlockScopeListEqual(testTemporaryTablesReferenceScope.Blocks, v.Result.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, testTemporaryTablesReferenceScope.Blocks, v.Result.Blocks)
		}
	}
}

func TestReferenceScope_StoreTemporaryTable(t *testing.T) {
	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
				"TEMPTABLE3": &View{
					Header: NewHeader("tempview3", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      "tempview3",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
		{
			scopeNameTempTables: {
				"TEMPTABLE2": &View{
					Header: NewHeader("temptable2", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("updated"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("updated"),
						}),
					},
					FileInfo: &FileInfo{
						Path:                  "temptable2",
						Delimiter:             ',',
						ViewType:              ViewTypeTemporaryTable,
						restorePointHeader:    NewHeader("table2", []string{"column1", "column2"}),
						restorePointRecordSet: RecordSet{},
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	expect := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
				"TEMPTABLE3": &View{
					Header: NewHeader("tempview3", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      "tempview3",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
		{
			scopeNameTempTables: {
				"TEMPTABLE2": &View{
					Header: NewHeader("temptable2", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("updated"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("updated"),
						}),
					},
					FileInfo: &FileInfo{
						Path:               "temptable2",
						Delimiter:          ',',
						ViewType:           ViewTypeTemporaryTable,
						restorePointHeader: NewHeader("temptable2", []string{"column1", "column2"}),
						restorePointRecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("updated"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("updated"),
							}),
						},
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	expectOut := []string{"Commit: restore point of view \"temptable2\" is created.\n"}

	UncommittedViews := map[string]*FileInfo{
		"TEMPTABLE2": nil,
	}

	log := scope.StoreTemporaryTable(TestTx.Session, UncommittedViews)

	if !BlockScopeListEqual(scope.Blocks, expect.Blocks) {
		t.Errorf("Store: blocks = %v, want %v", scope.Blocks, expect.Blocks)
	}

	if reflect.DeepEqual(log, expectOut) {
		t.Errorf("Store: log = %s, want %s", log, expectOut)
	}
}

func TestReferenceScope_RestoreTemporaryTable(t *testing.T) {
	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
				"TEMPTABLE3": &View{
					Header: NewHeader("tempview3", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      "tempview3",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
		{
			scopeNameTempTables: {
				"TEMPTABLE2": &View{
					Header: NewHeader("temptable2", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("updated"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("updated"),
						}),
					},
					FileInfo: &FileInfo{
						Path:               "temptable2",
						Delimiter:          ',',
						ViewType:           ViewTypeTemporaryTable,
						restorePointHeader: NewHeader("temptable2", []string{"column1", "column2"}),
						restorePointRecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
					},
				},
			},
		},
	}, nil, time.Time{}, nil)

	expect := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameTempTables: {
				"TEMPTABLE3": &View{
					Header: NewHeader("tempview3", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:      "tempview3",
						Delimiter: ',',
						ViewType:  ViewTypeTemporaryTable,
					},
				},
			},
		},
		{
			scopeNameTempTables: {
				"TEMPTABLE2": &View{
					Header: NewHeader("temptable2", []string{"column1", "column2"}),
					RecordSet: []Record{
						NewRecord([]value.Primary{
							value.NewString("1"),
							value.NewString("str1"),
						}),
						NewRecord([]value.Primary{
							value.NewString("2"),
							value.NewString("str2"),
						}),
					},
					FileInfo: &FileInfo{
						Path:               "temptable2",
						Delimiter:          ',',
						ViewType:           ViewTypeTemporaryTable,
						restorePointHeader: NewHeader("temptable2", []string{"column1", "column2"}),
						restorePointRecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
						},
					},
				},
			},
		},
	}, nil, time.Time{}, nil)
	expectOut := []string{"Rollback: view \"tempview2\" is restored.\n"}

	UncommittedViews := map[string]*FileInfo{
		"TEMPTABLE2": nil,
	}

	log := scope.RestoreTemporaryTable(UncommittedViews)

	if !BlockScopeListEqual(scope.Blocks, expect.Blocks) {
		t.Errorf("Restore: blocks = %v, want %v", scope.Blocks, expect.Blocks)
	}

	if reflect.DeepEqual(log, expectOut) {
		t.Errorf("Restore: log = %s, want %s", log, expectOut)
	}
}

var referenceScopeDisposeCursorTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  *ReferenceScope
	Error   string
}{
	{
		Name:    "ReferenceScope DisposeCursor",
		CurName: parser.Identifier{Literal: "cur"},
		Result: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameCursors: {
					"PCUR": &Cursor{
						view: &View{
							Header: NewHeader("", []string{"c1"}),
							RecordSet: RecordSet{
								NewRecord([]value.Primary{value.NewInteger(1)}),
								NewRecord([]value.Primary{value.NewInteger(2)}),
							},
						},
						index:    -1,
						isPseudo: true,
						mtx:      &sync.Mutex{},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "ReferenceScope DisposeCursor Pseudo Cursor Error",
		CurName: parser.Identifier{Literal: "pcur"},
		Error:   "cursor pcur is a pseudo cursor",
	},
	{
		Name:    "ReferenceScope DisposeCursor Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "cursor notexist is undeclared",
	},
}

func TestReferenceScope_DisposeCursor(t *testing.T) {
	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameCursors: {
				"PCUR": &Cursor{
					view: &View{
						Header: NewHeader("", []string{"c1"}),
						RecordSet: RecordSet{
							NewRecord([]value.Primary{value.NewInteger(1)}),
							NewRecord([]value.Primary{value.NewInteger(2)}),
						},
					},
					index:    -1,
					isPseudo: true,
					mtx:      &sync.Mutex{},
				},
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
					mtx:   &sync.Mutex{},
				},
			},
		},
	}, nil, time.Time{}, nil)

	for _, v := range referenceScopeDisposeCursorTests {
		err := scope.DisposeCursor(v.CurName)
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
		if !BlockScopeListEqual(scope.Blocks, v.Result.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, scope.Blocks, v.Result.Blocks)
		}
	}
}

var referenceScopeOpenCursorTests = []struct {
	Name      string
	CurName   parser.Identifier
	CurValues []parser.ReplaceValue
	Result    *ReferenceScope
	Error     string
}{
	{
		Name:    "ReferenceScope OpenCursor",
		CurName: parser.Identifier{Literal: "cur"},
		Result: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameCursors: {
					"PCUR": &Cursor{
						view: &View{
							Header: NewHeader("", []string{"c1"}),
							RecordSet: RecordSet{
								NewRecord([]value.Primary{value.NewInteger(1)}),
								NewRecord([]value.Primary{value.NewInteger(2)}),
							},
						},
						index:    -1,
						isPseudo: true,
						mtx:      &sync.Mutex{},
					},
					"CUR": &Cursor{
						query: selectQueryForCursorTest,
						view: &View{
							Header: NewHeader("table1", []string{"column1", "column2"}),
							RecordSet: []Record{
								NewRecord([]value.Primary{
									value.NewString("1"),
									value.NewString("str1"),
								}),
								NewRecord([]value.Primary{
									value.NewString("2"),
									value.NewString("str2"),
								}),
								NewRecord([]value.Primary{
									value.NewString("3"),
									value.NewString("str3"),
								}),
							},
							FileInfo: &FileInfo{
								Path:      GetTestFilePath("table1.csv"),
								Delimiter: ',',
								NoHeader:  false,
								Encoding:  text.UTF8,
								LineBreak: text.LF,
							},
						},
						index: -1,
						mtx:   &sync.Mutex{},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "ReferenceScope OpenCursor Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "cursor notexist is undeclared",
	},
	{
		Name:    "ReferenceScope OpenCursor Open Error",
		CurName: parser.Identifier{Literal: "cur"},
		Error:   "cursor cur is already open",
	},
	{
		Name:    "ReferenceScope OpenCursor Pseudo Cursor Error",
		CurName: parser.Identifier{Literal: "pcur"},
		Error:   "cursor pcur is a pseudo cursor",
	},
}

func TestReferenceScope_OpenCursor(t *testing.T) {
	defer func() {
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameCursors: {
				"PCUR": &Cursor{
					view: &View{
						Header: NewHeader("", []string{"c1"}),
						RecordSet: RecordSet{
							NewRecord([]value.Primary{value.NewInteger(1)}),
							NewRecord([]value.Primary{value.NewInteger(2)}),
						},
					},
					index:    -1,
					isPseudo: true,
					mtx:      &sync.Mutex{},
				},
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
					mtx:   &sync.Mutex{},
				},
			},
		},
	}, nil, time.Time{}, nil)

	for _, v := range referenceScopeOpenCursorTests {
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)

		err := scope.OpenCursor(context.Background(), v.CurName, v.CurValues)
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
		if !BlockScopeListEqual(scope.Blocks, v.Result.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, scope.Blocks, v.Result.Blocks)
		}
	}
}

var referenceScopeCloseCursorTests = []struct {
	Name    string
	CurName parser.Identifier
	Result  *ReferenceScope
	Error   string
}{
	{
		Name:    "ReferenceScope CloseCursor",
		CurName: parser.Identifier{Literal: "cur"},
		Result: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameCursors: {
					"PCUR": &Cursor{
						view: &View{
							Header: NewHeader("", []string{"c1"}),
							RecordSet: RecordSet{
								NewRecord([]value.Primary{value.NewInteger(1)}),
								NewRecord([]value.Primary{value.NewInteger(2)}),
							},
						},
						index:    -1,
						isPseudo: true,
						mtx:      &sync.Mutex{},
					},
					"CUR": &Cursor{
						query: selectQueryForCursorTest,
						mtx:   &sync.Mutex{},
					},
				},
			},
		}, nil, time.Time{}, nil),
	},
	{
		Name:    "ReferenceScope CloseCursor Pseudo Cursor Error",
		CurName: parser.Identifier{Literal: "pcur"},
		Error:   "cursor pcur is a pseudo cursor",
	},
	{
		Name:    "ReferenceScope CloseCursor Undeclared Error",
		CurName: parser.Identifier{Literal: "notexist"},
		Error:   "cursor notexist is undeclared",
	},
}

func TestReferenceScope_CloseCursor(t *testing.T) {
	defer func() {
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir

	scope := GenerateReferenceScope([]map[string]map[string]interface{}{
		{
			scopeNameCursors: {
				"PCUR": &Cursor{
					view: &View{
						Header: NewHeader("", []string{"c1"}),
						RecordSet: RecordSet{
							NewRecord([]value.Primary{value.NewInteger(1)}),
							NewRecord([]value.Primary{value.NewInteger(2)}),
						},
					},
					index:    -1,
					isPseudo: true,
					mtx:      &sync.Mutex{},
				},
				"CUR": &Cursor{
					query: selectQueryForCursorTest,
					view: &View{
						Header: NewHeader("table1", []string{"column1", "column2"}),
						RecordSet: []Record{
							NewRecord([]value.Primary{
								value.NewString("1"),
								value.NewString("str1"),
							}),
							NewRecord([]value.Primary{
								value.NewString("2"),
								value.NewString("str2"),
							}),
							NewRecord([]value.Primary{
								value.NewString("3"),
								value.NewString("str3"),
							}),
						},
						FileInfo: &FileInfo{
							Path:      GetTestFilePath("table1.csv"),
							Delimiter: ',',
							NoHeader:  false,
							Encoding:  text.UTF8,
							LineBreak: text.LF,
						},
					},
					index: -1,
					mtx:   &sync.Mutex{},
				},
			},
		},
	}, nil, time.Time{}, nil)

	for _, v := range referenceScopeCloseCursorTests {
		err := scope.CloseCursor(v.CurName)
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
		if !BlockScopeListEqual(scope.Blocks, v.Result.Blocks) {
			t.Errorf("%s: blocks = %v, want %v", v.Name, scope.Blocks, v.Result.Blocks)
		}
	}
}
