package query

import (
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

func TestTemporaryViewScopesExists(t *testing.T) {
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

func TestTemporaryViewScopesGet(t *testing.T) {
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
			t.Errorf("%s: view = %s, want %s", v.Name, view, v.Result)
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

func TestTemporaryViewScopesGetWithInternalId(t *testing.T) {
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
			t.Errorf("%s: view = %s, want %s", v.Name, view, v.Result)
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

func TestTemporaryViewScopesSet(t *testing.T) {
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
			t.Errorf("%s: map = %s, want %s", v.Name, list, v.Result)
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

func TestTemporaryViewScopesReplace(t *testing.T) {
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
			t.Errorf("%s: map = %s, want %s", v.Name, list, v.Result)
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
		Error: "[L:- C:-] temporary table /path/to/table9.csv is undefined",
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
			t.Errorf("%s: view = %s, want %s", v.Name, list, v.Result)
		}
	}
}

func TestTemporaryViewScopesRollback(t *testing.T) {
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

	list.Rollback()
	if !reflect.DeepEqual(list, expect) {
		t.Errorf("Rollback: view = %s, want %s", list, expect)
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
			t.Errorf("%s: view = %s, want %s", v.Name, view, v.Result)
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
			t.Errorf("%s: view = %s, want %s", v.Name, view, v.Result)
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
			t.Errorf("%s: map = %s, want %s", v.Name, viewMap, v.Result)
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
			t.Errorf("%s: map = %s, want %s", v.Name, viewMap, v.Result)
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
		Error: "[L:- C:-] temporary table /path/to/table2.csv is undefined",
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Undefined Error",
		Table: parser.Identifier{Literal: "/path/to/undef.csv"},
		Error: "[L:- C:-] temporary table /path/to/undef.csv is undefined",
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
			t.Errorf("%s: map = %s, want %s", v.Name, viewMap, v.Result)
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
		t.Errorf("result = %s, want %s", viewMap, expect)
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
