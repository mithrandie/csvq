package query

import (
	"context"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var temporaryViewScopesExistsTests = []struct {
	Name   string
	Path   string
	Result bool
}{
	{
		Name:   "TempViewScopes Exists",
		Path:   "/path/to/table1.csv",
		Result: true,
	},
	{
		Name:   "TempViewScopes Exists Not Exist",
		Path:   "/path/to/notexist.csv",
		Result: false,
	},
}

func TestTemporaryViewScopes_Exists(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
				Header:    NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header:    NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		}),
	}

	for _, v := range temporaryViewScopesExistsTests {
		result := list.Exists(v.Path)
		if result != v.Result {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Result)
		}
	}
}

var temporaryViewScopesGetTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "TempViewScopes Get",
		Path: parser.Identifier{Literal: "/path/to/table2.csv"},
		Result: &View{
			Header: NewHeader("table2", []string{"column1", "column2"}),
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
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	},
	{
		Name:  "TempViewScopes Get Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "table /path/to/table9.csv is not loaded",
	},
}

func TestTemporaryViewScopes_Get(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
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
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		}),
	}

	for _, v := range temporaryViewScopesGetTests {
		view, err := list.Get(v.Path)
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

var temporaryViewScopesGetWithInternalIdTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "TempViewScopes GetWithInternalId",
		Path: parser.Identifier{Literal: "/path/to/table2.csv"},
		Result: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column2"}),
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
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	},
	{
		Name:  "TempViewScopes GetWithInternalId Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "table /path/to/table9.csv is not loaded",
	},
}

func TestTemporaryViewScopes_GetWithInternalId(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
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
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		}),
	}

	for _, v := range temporaryViewScopesGetWithInternalIdTests {
		view, err := list.GetWithInternalId(context.Background(), v.Path, TestTx.Flags)
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

var temporaryViewScopesSetTests = []struct {
	Name    string
	SetView *View
	Result  TemporaryViewScopes
}{
	{
		Name: "TempViewScopes Set",
		SetView: &View{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
		Result: TemporaryViewScopes{
			GenerateViewMap([]*View{
				{
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
					},
					FileInfo: &FileInfo{
						Path:      "/path/to/table1.csv",
						Delimiter: ',',
					},
				},
			}),
			GenerateViewMap([]*View{
				{
					Header: NewHeader("table2", []string{"column1", "column2"}),
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
						Path:      "/path/to/table2.csv",
						Delimiter: ',',
					},
				},
			}),
		},
	},
}

func TestTemporaryViewScopes_Set(t *testing.T) {
	list := TemporaryViewScopes{
		NewViewMap(),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		}),
	}

	for _, v := range temporaryViewScopesSetTests {
		list.Set(v.SetView)
		if !SyncMapListEqual(TempViewScopesToSyncMapList(list), TempViewScopesToSyncMapList(v.Result)) {
			t.Errorf("%s: map = %v, want %v", v.Name, list, v.Result)
		}
	}
}

var temporaryViewScopesReplaceTests = []struct {
	Name    string
	SetView *View
	Result  TemporaryViewScopes
	Error   string
}{
	{
		Name: "TempViewScopes Replace",
		SetView: &View{
			Header: NewHeader("table2", []string{"column1", "column2"}),
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
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
		Result: TemporaryViewScopes{
			GenerateViewMap([]*View{
				{
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
					},
					FileInfo: &FileInfo{
						Path:      "/path/to/table1.csv",
						Delimiter: ',',
					},
				},
			}),
			GenerateViewMap([]*View{
				{
					Header: NewHeader("table2", []string{"column1", "column2"}),
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
						Path:      "/path/to/table2.csv",
						Delimiter: ',',
					},
				},
			}),
		},
	},
}

func TestTemporaryViewScopes_Replace(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
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
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		}),
	}

	for _, v := range temporaryViewScopesReplaceTests {
		list.Replace(v.SetView)
		if !SyncMapListEqual(TempViewScopesToSyncMapList(list), TempViewScopesToSyncMapList(v.Result)) {
			t.Errorf("%s: map = %v, want %v", v.Name, list, v.Result)
		}
	}
}

var temporaryViewScopesDisposeTests = []struct {
	Name   string
	Path   parser.Identifier
	Result TemporaryViewScopes
	Error  string
}{
	{
		Name: "TempViewScopes Dispose",
		Path: parser.Identifier{Literal: "/path/to/table1.csv"},
		Result: TemporaryViewScopes{
			NewViewMap(),
			GenerateViewMap([]*View{
				{
					Header: NewHeader("table2", []string{"column1", "column2"}),
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
						Path:      "/path/to/table2.csv",
						Delimiter: ',',
					},
				},
			}),
		},
	},
	{
		Name:  "TempViewScopes Dispose Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "view /path/to/table9.csv is undeclared",
	},
}

func TestTemporaryViewScopesDispose(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
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
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
					ViewType:  ViewTypeTemporaryTable,
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		}),
	}

	for _, v := range temporaryViewScopesDisposeTests {
		err := list.Dispose(v.Path)
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
		if !SyncMapListEqual(TempViewScopesToSyncMapList(list), TempViewScopesToSyncMapList(v.Result)) {
			t.Errorf("%s: view = %v, want %v", v.Name, list, v.Result)
		}
	}
}

func TestTemporaryViewScopes_Store(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
				RecordSet: RecordSet{
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
					Path:                  "/path/to/table1.csv",
					Delimiter:             ',',
					restorePointHeader:    NewHeader("table1", []string{"column1", "column2"}),
					restorePointRecordSet: RecordSet{},
				},
			},
			{
				Header: NewHeader("table2", []string{"column1", "column2", "column3"}),
				RecordSet: RecordSet{
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
					Path:                  "/path/to/table2.csv",
					Delimiter:             ',',
					restorePointHeader:    NewHeader("table2", []string{"column1", "column2", "column3"}),
					restorePointRecordSet: RecordSet{},
				},
			},
		}),
	}

	expect := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
				RecordSet: RecordSet{
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
					Path:               "/path/to/table1.csv",
					Delimiter:          ',',
					restorePointHeader: NewHeader("table1", []string{"column1", "column2", "column3"}),
					restorePointRecordSet: RecordSet{
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
			{
				Header: NewHeader("table2", []string{"column1", "column2", "column3"}),
				RecordSet: RecordSet{
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
					Path:                  "/path/to/table2.csv",
					Delimiter:             ',',
					restorePointHeader:    NewHeader("table2", []string{"column1", "column2", "column3"}),
					restorePointRecordSet: RecordSet{},
				},
			},
		}),
	}
	expectOut := []string{"Commit: restore point of view \"/path/to/table1.csv\" is created.\n"}

	UncommittedViews := map[string]*FileInfo{
		"/PATH/TO/TABLE1.CSV": nil,
	}

	log := list.Store(TestTx.Session, UncommittedViews)

	if !SyncMapListEqual(TempViewScopesToSyncMapList(list), TempViewScopesToSyncMapList(expect)) {
		t.Errorf("Store: view = %v, want %v", list, expect)
	}

	if reflect.DeepEqual(log, expectOut) {
		t.Errorf("Store: log = %s, want %s", log, expectOut)
	}
}

func TestTemporaryViewScopes_Restore(t *testing.T) {
	list := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
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
					Path:                  "/path/to/table1.csv",
					Delimiter:             ',',
					restorePointHeader:    NewHeader("table1", []string{"column1", "column2"}),
					restorePointRecordSet: RecordSet{},
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:               "/path/to/table2.csv",
					Delimiter:          ',',
					restorePointHeader: NewHeader("table2", []string{"column1", "column2"}),
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
		}),
	}

	expect := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
				Header:    NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{},
				FileInfo: &FileInfo{
					Path:                  "/path/to/table1.csv",
					Delimiter:             ',',
					restorePointHeader:    NewHeader("table1", []string{"column1", "column2"}),
					restorePointRecordSet: RecordSet{},
				},
			},
		}),
		GenerateViewMap([]*View{
			{
				Header: NewHeader("table2", []string{"column1", "column2"}),
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
					Path:               "/path/to/table2.csv",
					Delimiter:          ',',
					restorePointHeader: NewHeader("table2", []string{"column1", "column2"}),
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
		}),
	}
	expectOut := []string{"Rollback: view \"/path/to/table1.csv\" is restored.\nRollback: view \"/path/to/table2.csv\" is restored.\n"}

	UncommittedViews := map[string]*FileInfo{
		"/PATH/TO/TABLE1.CSV": nil,
		"/PATH/TO/TABLE2.CSV": nil,
	}

	log := list.Restore(UncommittedViews)

	if !SyncMapListEqual(TempViewScopesToSyncMapList(list), TempViewScopesToSyncMapList(expect)) {
		t.Errorf("Restore: view = %v, want %v", list, expect)
	}

	if reflect.DeepEqual(log, expectOut) {
		t.Errorf("Restore: log = %s, want %s", log, expectOut)
	}
}

func TestTemporaryViewScopes_All(t *testing.T) {
	scopes := TemporaryViewScopes{
		GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path:     "view1",
					ViewType: ViewTypeTemporaryTable,
				},
				Header: NewHeader("view1", []string{"column1", "column2"}),
			},
		}),
		GenerateViewMap([]*View{
			{
				FileInfo: &FileInfo{
					Path: "view1",
				},
				Header: NewHeader("view1", []string{"column1", "column2", "column3"}),
			},
			{
				FileInfo: &FileInfo{
					Path:     "view2",
					ViewType: ViewTypeTemporaryTable,
				},
				Header: NewHeader("view2", []string{"column1", "column2"}),
			},
		}),
	}

	expect := GenerateViewMap([]*View{
		{
			FileInfo: &FileInfo{
				Path:     "view1",
				ViewType: ViewTypeTemporaryTable,
			},
			Header: NewHeader("view1", []string{"column1", "column2"}),
		},
		{
			FileInfo: &FileInfo{
				Path:     "view2",
				ViewType: ViewTypeTemporaryTable,
			},
			Header: NewHeader("view2", []string{"column1", "column2"}),
		},
	})

	list := scopes.All()
	if !SyncMapEqual(list, expect) {
		t.Errorf("List: list = %v, want %v", list, expect)
	}
}

var viewMapExistsTests = []struct {
	Name   string
	Path   string
	Result bool
}{
	{
		Name:   "ViewMap Exists",
		Path:   "/path/to/table1.csv",
		Result: true,
	},
	{
		Name:   "ViewMap Exists Not Exist",
		Path:   "/path/to/notexist.csv",
		Result: false,
	},
}

func TestViewMap_Exists(t *testing.T) {
	viewMap := GenerateViewMap([]*View{
		{
			Header:    NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
	})

	for _, v := range viewMapExistsTests {
		result := viewMap.Exists(v.Path)
		if result != v.Result {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Result)
		}
	}
}

var viewMapGetTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "ViewMap Get",
		Path: parser.Identifier{Literal: "/path/to/table1.csv"},
		Result: &View{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
	},
	{
		Name:  "ViewMap Get Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table2.csv"},
		Error: "table /path/to/table2.csv is not loaded",
	},
}

func TestViewMap_Get(t *testing.T) {
	viewMap := GenerateViewMap([]*View{
		{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
	})

	for _, v := range viewMapGetTests {
		view, err := viewMap.Get(v.Path)
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

var viewMapGetWithInternalIdTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "ViewMap GetWithInternalId",
		Path: parser.Identifier{Literal: "/path/to/table1.csv"},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
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
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
	},
	{
		Name:  "ViewMap GetWithInternalId Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table2.csv"},
		Error: "table /path/to/table2.csv is not loaded",
	},
}

func TestViewMap_GetWithInternalId(t *testing.T) {
	viewMap := GenerateViewMap([]*View{
		{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
	})

	for _, v := range viewMapGetWithInternalIdTests {
		view, err := viewMap.GetWithInternalId(context.Background(), v.Path, TestTx.Flags)
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

var viewMapSetTests = []struct {
	Name    string
	SetView *View
	Result  ViewMap
}{
	{
		Name: "ViewMap Set",
		SetView: &View{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
		Result: GenerateViewMap([]*View{
			{
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
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		}),
	},
}

func TestViewMap_Set(t *testing.T) {
	viewMap := NewViewMap()

	for _, v := range viewMapSetTests {
		viewMap.Set(v.SetView)
		if !SyncMapEqual(viewMap, v.Result) {
			t.Errorf("%s: map = %v, want %v", v.Name, viewMap, v.Result)
		}
	}
}

var viewMapDisposeTemporaryTable = []struct {
	Name   string
	Table  parser.Identifier
	Result ViewMap
	Error  string
}{
	{
		Name:  "ViewMap DisposeTemporaryTable",
		Table: parser.Identifier{Literal: "/path/to/table1.csv"},
		Result: GenerateViewMap([]*View{
			{
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
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		}),
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Not Temporary Table",
		Table: parser.Identifier{Literal: "/path/to/table2.csv"},
		Error: "view /path/to/table2.csv is undeclared",
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Undeclared Error",
		Table: parser.Identifier{Literal: "/path/to/undef.csv"},
		Error: "view /path/to/undef.csv is undeclared",
	},
}

func TestViewMap_DisposeTemporaryTable(t *testing.T) {
	viewMap := GenerateViewMap([]*View{
		{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
		{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	})

	for _, v := range viewMapDisposeTemporaryTable {
		err := viewMap.DisposeTemporaryTable(v.Table)
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
		if !SyncMapEqual(viewMap, v.Result) {
			t.Errorf("%s: map = %v, want %v", v.Name, viewMap, v.Result)
		}
	}
}

func TestViewMap_Clear(t *testing.T) {
	viewMap := GenerateViewMap([]*View{
		{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
				ViewType:  ViewTypeTemporaryTable,
			},
		},
		{
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
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	})

	expect := NewViewMap()

	_ = viewMap.Clean(TestTx.FileContainer)
	if !SyncMapEqual(viewMap, expect) {
		t.Errorf("result = %v, want %v", viewMap, expect)
	}
}

var viewMapGetWithInternalIdBench = generateViewMapGetWithInternalIdBenchViewMap()

func generateViewMapGetWithInternalIdBenchViewMap() ViewMap {
	m := GenerateViewMap([]*View{{
		Header: NewHeader("bench_view", []string{"c1", "c2", "c3", "c4"}),
		FileInfo: &FileInfo{
			Path: "bench_view",
		},
	}})
	view, _ := m.Load("bench_view")
	view.RecordSet = make(RecordSet, 10000)
	for i := 0; i < 10000; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
			value.NewInteger(4),
		})
	}
	return m
}

func BenchmarkViewMap_GetWithInternalId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = viewMapGetWithInternalIdBench.GetWithInternalId(context.Background(), parser.Identifier{Literal: "BENCH_VIEW"}, TestTx.Flags)
	}
}
