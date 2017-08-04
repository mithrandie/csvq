package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

var temporaryViewMapListExistsTests = []struct {
	Name   string
	Path   string
	Result bool
}{
	{
		Name:   "TemporaryViewMapList Exists",
		Path:   "/path/to/table1.csv",
		Result: true,
	},
	{
		Name:   "TemporaryViewMapList Exists Not Exist",
		Path:   "/path/to/notexist.csv",
		Result: false,
	},
}

func TestTemporaryViewMapList_Exists(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header:  NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
				Header:  NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
				},
			},
		},
	}

	for _, v := range temporaryViewMapListExistsTests {
		result := list.Exists(v.Path)
		if result != v.Result {
			t.Errorf("%s: result = %t, want %t", v.Name, result, v.Result)
		}
	}
}

var temporaryViewMapListGetTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "TemporaryViewMapList Get",
		Path: parser.Identifier{Literal: "/path/to/table2.csv"},
		Result: &View{
			Header: NewHeader("table2", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	},
	{
		Name:  "TemporaryViewMapList Get Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "[L:- C:-] table /path/to/table9.csv is not loaded",
	},
}

func TestTemporaryViewMapList_Get(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
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
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		},
	}

	for _, v := range temporaryViewMapListGetTests {
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

var temporaryViewMapListGetWithInternalIdTests = []struct {
	Name   string
	Path   parser.Identifier
	Result *View
	Error  string
}{
	{
		Name: "TemporaryViewMapList GetWithInternalId",
		Path: parser.Identifier{Literal: "/path/to/table2.csv"},
		Result: &View{
			Header: NewHeaderWithId("table2", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithId(0, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithId(1, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	},
	{
		Name:  "TemporaryViewMapList GetWithInternalId Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "[L:- C:-] table /path/to/table9.csv is not loaded",
	},
}

func TestTemporaryViewMapList_GetWithInternalId(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
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
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		},
	}

	for _, v := range temporaryViewMapListGetWithInternalIdTests {
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

var temporaryViewMapListSetTests = []struct {
	Name    string
	SetView *View
	Result  TemporaryViewMapList
}{
	{
		Name: "TemporaryViewMapList Set",
		SetView: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
			},
		},
		Result: TemporaryViewMapList{
			ViewMap{
				"/PATH/TO/TABLE1.CSV": &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
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
					Records: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
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

func TestTemporaryViewMapList_Set(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
				Header: NewHeader("table2", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		},
	}

	for _, v := range temporaryViewMapListSetTests {
		list.Set(v.SetView)
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: map = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var temporaryViewMapListReplaceTests = []struct {
	Name    string
	SetView *View
	Result  TemporaryViewMapList
	Error   string
}{
	{
		Name: "TemporaryViewMapList Replace",
		SetView: &View{
			Header: NewHeader("table2", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("updated"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("updated"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
		Result: TemporaryViewMapList{
			ViewMap{
				"/PATH/TO/TABLE1.CSV": &View{
					Header: NewHeader("table1", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
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
					Records: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("updated"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("updated"),
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

func TestTemporaryViewMapList_Replace(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
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
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		},
	}

	for _, v := range temporaryViewMapListReplaceTests {
		list.Replace(v.SetView)
		if !reflect.DeepEqual(list, v.Result) {
			t.Errorf("%s: map = %s, want %s", v.Name, list, v.Result)
		}
	}
}

var temporaryViewMapListDisposeTests = []struct {
	Name   string
	Path   parser.Identifier
	Result TemporaryViewMapList
	Error  string
}{
	{
		Name: "TemporaryViewMapList Dispose",
		Path: parser.Identifier{Literal: "/path/to/table1.csv"},
		Result: TemporaryViewMapList{
			ViewMap{},
			ViewMap{
				"/PATH/TO/TABLE2.CSV": &View{
					Header: NewHeader("table2", []string{"column1", "column2"}),
					Records: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("updated"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("updated"),
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
		Name:  "TemporaryViewMapList Dispose Not Loaded Error",
		Path:  parser.Identifier{Literal: "/path/to/table9.csv"},
		Error: "[L:- C:-] temporary table /path/to/table9.csv is undefined",
	},
}

func TestTemporaryViewMapList_Dispose(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table1.csv",
					Delimiter: ',',
					Temporary: true,
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
				Header: NewHeader("table2", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("updated"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("updated"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
				},
			},
		},
	}

	for _, v := range temporaryViewMapListDisposeTests {
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

func TestTemporaryViewMapList_Rollback(t *testing.T) {
	list := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header: NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:           "/path/to/table1.csv",
					Delimiter:      ',',
					InitialRecords: Records{},
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
				Header: NewHeader("table2", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("updated"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("updated"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
					InitialRecords: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
						}),
					},
				},
			},
		},
	}

	expect := TemporaryViewMapList{
		ViewMap{
			"/PATH/TO/TABLE1.CSV": &View{
				Header:  NewHeader("table1", []string{"column1", "column2"}),
				Records: []Record{},
				FileInfo: &FileInfo{
					Path:           "/path/to/table1.csv",
					Delimiter:      ',',
					InitialRecords: Records{},
				},
			},
		},
		ViewMap{
			"/PATH/TO/TABLE2.CSV": &View{
				Header: NewHeader("table2", []string{"column1", "column2"}),
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
					}),
				},
				FileInfo: &FileInfo{
					Path:      "/path/to/table2.csv",
					Delimiter: ',',
					InitialRecords: []Record{
						NewRecord([]parser.Primary{
							parser.NewString("1"),
							parser.NewString("str1"),
						}),
						NewRecord([]parser.Primary{
							parser.NewString("2"),
							parser.NewString("str2"),
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
			Header:  NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{},
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
			Records: []Record{
				NewRecordWithId(0, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithId(1, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("updated"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("updated"),
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
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("updated"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("updated"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("updated"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("updated"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
				Records: []Record{
					NewRecord([]parser.Primary{
						parser.NewString("1"),
						parser.NewString("str1"),
					}),
					NewRecord([]parser.Primary{
						parser.NewString("2"),
						parser.NewString("str2"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
				Temporary: true,
			},
		},
		"/PATH/TO/TABLE2.CSV": &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
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
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table1.csv",
				Delimiter: ',',
				Temporary: true,
			},
		},
		"/PATH/TO/TABLE2.CSV": &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "/path/to/table2.csv",
				Delimiter: ',',
			},
		},
	}

	expect := ViewMap{}

	viewMap.Clear()
	if !reflect.DeepEqual(viewMap, expect) {
		t.Errorf("result = %s, want %s", viewMap, expect)
	}
}
