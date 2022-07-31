package query

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var viewMapExistsTests = []struct {
	Name   string
	Path   string
	Result bool
}{
	{
		Name:   "ViewMap Exists",
		Path:   strings.ToUpper("/path/to/table1.csv"),
		Result: true,
	},
	{
		Name:   "ViewMap Exists Not Exist",
		Path:   strings.ToUpper("/path/to/notexist.csv"),
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
	Path   string
	Result *View
	Error  string
}{
	{
		Name: "ViewMap Get",
		Path: strings.ToUpper("/path/to/table1.csv"),
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
		Path:  strings.ToUpper("/path/to/table2.csv"),
		Error: "table not loaded",
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
	Path   string
	Result *View
	Error  string
}{
	{
		Name: "ViewMap GetWithInternalId",
		Path: strings.ToUpper("/path/to/table1.csv"),
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
		Path:  strings.ToUpper("/path/to/table2.csv"),
		Error: "table not loaded",
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
	OK     bool
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
		OK: true,
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Not Temporary Table",
		Table: parser.Identifier{Literal: "/path/to/table2.csv"},
		OK:    false,
	},
	{
		Name:  "ViewMap DisposeTemporaryTable Undeclared Error",
		Table: parser.Identifier{Literal: "/path/to/undef.csv"},
		OK:    false,
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
		ok := viewMap.DisposeTemporaryTable(v.Table)
		if ok != v.OK {
			t.Errorf("%s: result = %t, want %t", v.Name, ok, v.OK)
		}
		if ok && v.OK && !SyncMapEqual(viewMap, v.Result) {
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
	view, _ := m.Load("BENCH_VIEW")
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
		_, _ = viewMapGetWithInternalIdBench.GetWithInternalId(context.Background(), strings.ToUpper("BENCH_VIEW"), TestTx.Flags)
	}
}
