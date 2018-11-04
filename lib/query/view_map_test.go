package query

import (
	"io/ioutil"
	"os"
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
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header:    NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
				Header:    NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		},
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
		Error: "[L:- C:-] table /path/to/table9.csv is not loaded",
	},
}

func TestTemporaryViewScopes_Get(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
		Error: "[L:- C:-] table /path/to/table9.csv is not loaded",
	},
}

func TestTemporaryViewScopes_GetWithInternalId(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
	}

	for _, v := range temporaryViewScopesGetWithInternalIdTests {
		view, err := list.GetWithInternalId(v.Path)
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
			ViewMap{
				"/PATH/TO/TABLE1.CSV": &View{
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
			ViewMap{
				"/PATH/TO/TABLE2.CSV": &View{
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
		},
	},
}

func TestTemporaryViewScopes_Set(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
	}

	for _, v := range temporaryViewScopesSetTests {
		list.Set(v.SetView)
		if !reflect.DeepEqual(list, v.Result) {
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
			ViewMap{
				"/PATH/TO/TABLE1.CSV": &View{
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
			ViewMap{
				"/PATH/TO/TABLE2.CSV": &View{
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
			},
		},
	},
}

func TestTemporaryViewScopes_Replace(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
	}

	for _, v := range temporaryViewScopesReplaceTests {
		list.Replace(v.SetView)
		if !reflect.DeepEqual(list, v.Result) {
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
			ViewMap{},
			ViewMap{
				"/PATH/TO/TABLE2.CSV": &View{
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
			},
		},
	},
	{
		Name:  "TempViewScopes Dispose Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "[L:- C:-] view /path/to/table9.csv is undeclared",
	},
}

func TestTemporaryViewScopesDispose(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
					Path:        "/path/to/table1.csv",
					Delimiter:   ',',
					IsTemporary: true,
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
		},
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
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: view = %v, want %v", v.Name, list, v.Result)
		}
	}
}

func TestTemporaryViewScopes_Store(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
					Path:             "/path/to/table1.csv",
					Delimiter:        ',',
					InitialHeader:    NewHeader("table1", []string{"column1", "column2"}),
					InitialRecordSet: RecordSet{},
				},
			},
			"/PATH/TO/TABLE2.CSV": &View{
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
					Path:             "/path/to/table2.csv",
					Delimiter:        ',',
					InitialHeader:    NewHeader("table2", []string{"column1", "column2", "column3"}),
					InitialRecordSet: RecordSet{},
				},
			},
		},
	}

	expect := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
					Path:          "/path/to/table1.csv",
					Delimiter:     ',',
					InitialHeader: NewHeader("table1", []string{"column1", "column2", "column3"}),
					InitialRecordSet: RecordSet{
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
			"/PATH/TO/TABLE2.CSV": &View{
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
					Path:             "/path/to/table2.csv",
					Delimiter:        ',',
					InitialHeader:    NewHeader("table2", []string{"column1", "column2", "column3"}),
					InitialRecordSet: RecordSet{},
				},
			},
		},
	}
	expectOut := "Commit: restore point of view \"/path/to/table1.csv\" is created.\n"

	UncommittedViews := map[string]*FileInfo{
		"/path/to/table1.csv": nil,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	list.Store(UncommittedViews)

	w.Close()
	os.Stdout = oldStdout

	log, _ := ioutil.ReadAll(r)

	if !reflect.DeepEqual(list, expect) {
		t.Errorf("Store: view = %v, want %v", list, expect)
	}

	if string(log) != expectOut {
		t.Errorf("Store: log = %s, want %s", string(log), expectOut)
	}
}

func TestTemporaryViewScopes_Restore(t *testing.T) {
	list := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
					Path:             "/path/to/table1.csv",
					Delimiter:        ',',
					InitialHeader:    NewHeader("table1", []string{"column1", "column2"}),
					InitialRecordSet: RecordSet{},
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
					Path:          "/path/to/table2.csv",
					Delimiter:     ',',
					InitialHeader: NewHeader("table2", []string{"column1", "column2"}),
					InitialRecordSet: []Record{
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
	}

	expect := TemporaryViewScopes{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header:    NewHeader("table1", []string{"column1", "column2"}),
				RecordSet: []Record{},
				FileInfo: &FileInfo{
					Path:             "/path/to/table1.csv",
					Delimiter:        ',',
					InitialHeader:    NewHeader("table1", []string{"column1", "column2"}),
					InitialRecordSet: RecordSet{},
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
					Path:          "/path/to/table2.csv",
					Delimiter:     ',',
					InitialHeader: NewHeader("table2", []string{"column1", "column2"}),
					InitialRecordSet: []Record{
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
	}
	expectOut := "Rollback: view \"/path/to/table1.csv\" is restored.\nRollback: view \"/path/to/table2.csv\" is restored.\n"

	UncommittedViews := map[string]*FileInfo{
		"/path/to/table1.csv": nil,
		"/path/to/table2.csv": nil,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	list.Restore(UncommittedViews)

	w.Close()
	os.Stdout = oldStdout

	log, _ := ioutil.ReadAll(r)

	if !reflect.DeepEqual(list, expect) {
		t.Errorf("Restore: view = %v, want %v", list, expect)
	}

	if string(log) != expectOut {
		t.Errorf("Restore: log = %s, want %s", string(log), expectOut)
	}
}

func TestTemporaryViewScopes_All(t *testing.T) {
	scopes := TemporaryViewScopes{
		ViewMap{
			"VIEW1": &View{
				FileInfo: &FileInfo{
					Path:        "view1",
					IsTemporary: true,
				},
				Header: NewHeader("view1", []string{"column1", "column2"}),
			},
		},
		ViewMap{
			"VIEW1": &View{
				FileInfo: &FileInfo{
					Path: "view1",
				},
				Header: NewHeader("view1", []string{"column1", "column2", "column3"}),
			},
			"VIEW2": &View{
				FileInfo: &FileInfo{
					Path:        "view2",
					IsTemporary: true,
				},
				Header: NewHeader("view2", []string{"column1", "column2"}),
			},
		},
	}

	expect := ViewMap{
		"VIEW1": &View{
			FileInfo: &FileInfo{
				Path:        "view1",
				IsTemporary: true,
			},
			Header: NewHeader("view1", []string{"column1", "column2"}),
		},
		"VIEW2": &View{
			FileInfo: &FileInfo{
				Path:        "view2",
				IsTemporary: true,
			},
			Header: NewHeader("view2", []string{"column1", "column2"}),
		},
	}

	list := scopes.All()
	if !reflect.DeepEqual(list, expect) {
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
	viewMap := ViewMap{
		"/PATH/TO/TABLE1.CSV": &View{
			Header:    NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
	}

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
		Error: "[L:- C:-] table /path/to/table2.csv is not loaded",
	},
}

func TestViewMap_Get(t *testing.T) {
	viewMap := ViewMap{
		"/PATH/TO/TABLE1.CSV": &View{
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
	}

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
		Error: "[L:- C:-] table /path/to/table2.csv is not loaded",
	},
}

func TestViewMap_GetWithInternalId(t *testing.T) {
	viewMap := ViewMap{
		"/PATH/TO/TABLE1.CSV": &View{
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
	}

	for _, v := range viewMapGetWithInternalIdTests {
		view, err := viewMap.GetWithInternalId(v.Path)
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
		Result: ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
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
	},
}

func TestViewMap_Set(t *testing.T) {
	viewMap := ViewMap{}

	for _, v := range viewMapSetTests {
		viewMap.Set(v.SetView)
		if !reflect.DeepEqual(viewMap, v.Result) {
			t.Errorf("%s: map = %v, want %v", v.Name, viewMap, v.Result)
		}
	}
}

var viewMapReplaceTests = []struct {
	Name    string
	SetView *View
	Result  ViewMap
	Error   string
}{
	{
		Name: "ViewMap Replace",
		SetView: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
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
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
		Result: ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
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
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		},
	},
	{
		Name: "ViewMap Replace Not Loaded Error",
		SetView: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
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
		Error: "[L:- C:-] table /path/to/table2.csv is not loaded",
	},
}

func TestViewMap_Replace(t *testing.T) {
	viewMap := ViewMap{
		"/PATH/TO/TABLE1.CSV": &View{
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
	}

	for _, v := range viewMapReplaceTests {
		err := viewMap.Replace(v.SetView)
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
		if !reflect.DeepEqual(viewMap, v.Result) {
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
		Result: ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
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
		},
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Not Temporary Table",
		Table: parser.Identifier{Literal: "/path/to/table2.csv"},
		Error: "[L:- C:-] view /path/to/table2.csv is undeclared",
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Undeclared Error",
		Table: parser.Identifier{Literal: "/path/to/undef.csv"},
		Error: "[L:- C:-] view /path/to/undef.csv is undeclared",
	},
}

func TestViewMap_DisposeTemporaryTable(t *testing.T) {
	viewMap := ViewMap{
		"/PATH/TO/TABLE1.CSV": &View{
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
				Path:        "/path/to/table1.csv",
				Delimiter:   ',',
				IsTemporary: true,
			},
		},
		"/PATH/TO/TABLE2.CSV": &View{
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
	}

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
		if !reflect.DeepEqual(viewMap, v.Result) {
			t.Errorf("%s: map = %v, want %v", v.Name, viewMap, v.Result)
		}
	}
}

func TestViewMap_Clear(t *testing.T) {
	viewMap := ViewMap{
		"/PATH/TO/TABLE1.CSV": &View{
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
				Path:        "/path/to/table1.csv",
				Delimiter:   ',',
				IsTemporary: true,
			},
		},
		"/PATH/TO/TABLE2.CSV": &View{
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
	}

	expect := ViewMap{}

	viewMap.Clean()
	if !reflect.DeepEqual(viewMap, expect) {
		t.Errorf("result = %v, want %v", viewMap, expect)
	}
}

var viewMapGetWithInternalIdBench = generateViewMapGetWithInternalIdBenchViewMap()

func generateViewMapGetWithInternalIdBenchViewMap() ViewMap {
	m := ViewMap{
		"BENCH_VIEW": &View{
			Header: NewHeader("bench_view", []string{"c1", "c2", "c3", "c4"}),
		},
	}
	view := m["BENCH_VIEW"]
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
		viewMapGetWithInternalIdBench.GetWithInternalId(parser.Identifier{Literal: "BENCH_VIEW"})
	}
}
