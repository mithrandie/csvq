package query

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"fmt"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

var viewLoadTests = []struct {
	Name     string
	Encoding cmd.Encoding
	NoHeader bool
	From     parser.FromClause
	Stdin    string
	Filter   Filter
	Result   *View
	Error    string
}{
	{
		Name: "Dual View",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{Object: parser.Dual{}},
			},
		},
		Result: &View{
			Header: []HeaderField{{}},
			Records: []Record{
				{
					NewCell(parser.NewNull()),
				},
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList:      AliasMapList{{}},
			},
		},
	},
	{
		Name: "Load File",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1.csv"},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
					},
				},
			},
		},
	},
	{
		Name:  "Load From Stdin",
		From:  parser.FromClause{},
		Stdin: "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeaderWithoutId("stdin", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "stdin",
				Delimiter: ',',
			},
			ParentFilter: Filter{
				VariablesList: []Variables{{}},
				TempViewsList: []ViewMap{
					{
						"STDIN": nil,
					},
				},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"STDIN": "STDIN",
					},
				},
			},
		},
	},
	{
		Name: "Stdin Empty Error",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Stdin{Stdin: "stdin"},
				},
			},
		},
		Error: "[L:- C:-] stdin is empty",
	},
	{
		Name: "Load File Error",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Literal: "notexist"},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name:     "Load SJIS File",
		Encoding: cmd.SJIS,
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_sjis"},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table_sjis", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("日本語"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE_SJIS": strings.ToUpper(GetTestFilePath("table_sjis.csv")),
					},
				},
			},
		},
	},
	{
		Name:     "Load No Header File",
		NoHeader: true,
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_noheader"},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("table_noheader", []string{"c1", "c2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE_NOHEADER": strings.ToUpper(GetTestFilePath("table_noheader.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Load Multiple File",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1"},
				},
				parser.Table{
					Object: parser.Identifier{Literal: "table2"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Number: 2, FromTable: true},
				{Reference: "table2", Column: "column3", Number: 1, FromTable: true},
				{Reference: "table2", Column: "column4", Number: 2, FromTable: true},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Cross Join",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Number: 2, FromTable: true},
				{Reference: "table2", Column: "column3", Number: 1, FromTable: true},
				{Reference: "table2", Column: "column4", Number: 2, FromTable: true},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Inner Join",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Number: 2, FromTable: true},
				{Reference: "table2", Column: "column3", Number: 1, FromTable: true},
				{Reference: "table2", Column: "column4", Number: 2, FromTable: true},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Outer Join",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						Direction: parser.Token{Token: parser.LEFT, Literal: "left"},
						Condition: parser.JoinCondition{
							On: parser.Comparison{
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "column1"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column3"}},
								Operator: "=",
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Number: 2, FromTable: true},
				{Reference: "table2", Column: "column3", Number: 1, FromTable: true},
				{Reference: "table2", Column: "column4", Number: 2, FromTable: true},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewNull(),
					parser.NewNull(),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
						"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
					},
				},
			},
		},
	},
	{
		Name: "Join Left Side Table File Not Exist Error",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "notexist"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table2"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Join Right Side Table File Not Exist Error",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "notexist"},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Error: "[L:- C:-] file notexist does not exist",
	},
	{
		Name: "Load Subquery",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Subquery{
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
									Tables: []parser.Expression{
										parser.Table{Object: parser.Identifier{Literal: "table1"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "alias"},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithoutId("alias", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
			ParentFilter: Filter{
				VariablesList:    []Variables{{}},
				TempViewsList:    []ViewMap{{}},
				CursorsList:      []CursorMap{{}},
				InlineTablesList: InlineTablesList{{}},
				AliasesList: AliasMapList{
					{
						"ALIAS": "",
					},
				},
			},
		},
	},
	{
		Name: "Load CSV Parse Error",
		From: parser.FromClause{
			Tables: []parser.Expression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_broken.csv"},
				},
			},
		},
		Error: fmt.Sprintf("[L:- C:-] csv parse error in file %s: line 3, column 7: wrong number of fields in line", GetTestFilePath("table_broken.csv")),
	},
}

func TestView_Load(t *testing.T) {
	tf := cmd.GetFlags()
	tf.Repository = TestDir

	for _, v := range viewLoadTests {
		ViewCache.Clear()

		tf.Delimiter = cmd.UNDEF
		tf.NoHeader = v.NoHeader
		if v.Encoding != "" {
			tf.Encoding = v.Encoding
		} else {
			tf.Encoding = cmd.UTF8
		}

		var oldStdin *os.File
		if 0 < len(v.Stdin) {
			oldStdin = os.Stdin
			r, w, _ := os.Pipe()
			w.WriteString(v.Stdin)
			w.Close()
			os.Stdin = r
		}

		view := NewView()
		if reflect.DeepEqual(v.Filter, Filter{}) {
			v.Filter = NewEmptyFilter()
		}
		err := view.Load(v.From, v.Filter.CreateNode())

		if 0 < len(v.Stdin) {
			os.Stdin = oldStdin
		}

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

		if v.Result.FileInfo != nil {
			if filepath.Base(view.FileInfo.Path) != filepath.Base(v.Result.FileInfo.Path) {
				t.Errorf("%s: filepath = %q, want %q", v.Name, filepath.Base(view.FileInfo.Path), filepath.Base(v.Result.FileInfo.Path))
			}
			if view.FileInfo.Delimiter != v.Result.FileInfo.Delimiter {
				t.Errorf("%s: delimiter = %q, want %q", v.Name, view.FileInfo.Delimiter, v.Result.FileInfo.Delimiter)
			}
		}
		view.FileInfo = nil
		v.Result.FileInfo = nil

		if !reflect.DeepEqual(view.ParentFilter.AliasesList, v.Result.ParentFilter.AliasesList) {
			t.Errorf("%s: alias list = %q, want %q", v.Name, view.ParentFilter.AliasesList, v.Result.ParentFilter.AliasesList)
		}
		for i, tviews := range v.Result.ParentFilter.TempViewsList {
			resultKeys := []string{}
			for key := range tviews {
				resultKeys = append(resultKeys, key)
			}
			viewKeys := []string{}
			for key := range view.ParentFilter.TempViewsList[i] {
				viewKeys = append(viewKeys, key)
			}
			if !reflect.DeepEqual(resultKeys, viewKeys) {
				t.Errorf("%s: temp view list = %q, want %q", v.Name, view.ParentFilter.TempViewsList, v.Result.ParentFilter.TempViewsList)
			}
		}

		view.ParentFilter = NewEmptyFilter()
		v.Result.ParentFilter = NewEmptyFilter()

		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, view, v.Result)
		}
	}
}

func TestNewViewFromGroupedRecord(t *testing.T) {
	fr := FilterRecord{
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				{
					NewGroupCell([]parser.Primary{parser.NewInteger(1), parser.NewInteger(2), parser.NewInteger(3)}),
					NewGroupCell([]parser.Primary{parser.NewInteger(1), parser.NewInteger(2), parser.NewInteger(3)}),
					NewGroupCell([]parser.Primary{parser.NewString("str1"), parser.NewString("str2"), parser.NewString("str3")}),
				},
			},
		},
		RecordIndex: 0,
	}
	expect := &View{
		Header: NewHeader("table1", []string{"column1", "column2"}),
		Records: []Record{
			{NewCell(parser.NewInteger(1)), NewCell(parser.NewInteger(1)), NewCell(parser.NewString("str1"))},
			{NewCell(parser.NewInteger(2)), NewCell(parser.NewInteger(2)), NewCell(parser.NewString("str2"))},
			{NewCell(parser.NewInteger(3)), NewCell(parser.NewInteger(3)), NewCell(parser.NewString("str3"))},
		},
	}

	result := NewViewFromGroupedRecord(fr)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

var viewWhereTests = []struct {
	Name   string
	View   *View
	Where  parser.WhereClause
	Result []int
	Error  string
}{
	{
		Name: "Where",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewInteger(2),
				Operator: "=",
			},
		},
		Result: []int{1},
	},
	{
		Name: "Where Filter Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				RHS:      parser.NewInteger(2),
				Operator: "=",
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_Where(t *testing.T) {
	for _, v := range viewWhereTests {
		err := v.View.Where(v.Where)
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
		if !reflect.DeepEqual(v.View.filteredIndices, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.filteredIndices, v.Result)
		}
	}
}

var viewGroupByTests = []struct {
	Name       string
	View       *View
	GroupBy    parser.GroupByClause
	Result     *View
	IsGrouped  bool
	GroupItems []string
	Error      string
}{
	{
		Name: "Group By",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewString("group1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("group2"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
					parser.NewString("group1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("4"),
					parser.NewString("str4"),
					parser.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.Expression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    INTERNAL_ID_COLUMN,
				},
				{
					Reference: "table1",
					Column:    "column1",
					Number:    1,
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					Number:    2,
					FromTable: true,
				},
				{
					Reference:  "table1",
					Column:     "column3",
					Number:     3,
					FromTable:  true,
					IsGroupKey: true,
				},
			},
			Records: []Record{
				{
					NewGroupCell([]parser.Primary{parser.NewInteger(1), parser.NewInteger(3)}),
					NewGroupCell([]parser.Primary{parser.NewString("1"), parser.NewString("3")}),
					NewGroupCell([]parser.Primary{parser.NewString("str1"), parser.NewString("str3")}),
					NewGroupCell([]parser.Primary{parser.NewString("group1"), parser.NewString("group1")}),
				},
				{
					NewGroupCell([]parser.Primary{parser.NewInteger(2), parser.NewInteger(4)}),
					NewGroupCell([]parser.Primary{parser.NewString("2"), parser.NewString("4")}),
					NewGroupCell([]parser.Primary{parser.NewString("str2"), parser.NewString("str4")}),
					NewGroupCell([]parser.Primary{parser.NewString("group2"), parser.NewString("group2")}),
				},
			},
			isGrouped: true,
		},
	},
}

func TestView_GroupBy(t *testing.T) {
	for _, v := range viewGroupByTests {
		err := v.View.GroupBy(v.GroupBy)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View, v.Result)
		}
	}
}

var viewHavingTests = []struct {
	Name   string
	View   *View
	Having parser.HavingClause
	Result []int
	Record Records
	Error  string
}{
	{
		Name: "Having",
		View: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    INTERNAL_ID_COLUMN,
				},
				{
					Reference: "table1",
					Column:    "column1",
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					FromTable: true,
				},
				{
					Reference:  "table1",
					Column:     "column3",
					FromTable:  true,
					IsGroupKey: true,
				},
			},
			isGrouped: true,
			Records: []Record{
				{
					NewGroupCell([]parser.Primary{parser.NewString("1"), parser.NewString("3")}),
					NewGroupCell([]parser.Primary{parser.NewString("1"), parser.NewString("3")}),
					NewGroupCell([]parser.Primary{parser.NewString("str1"), parser.NewString("str3")}),
					NewGroupCell([]parser.Primary{parser.NewString("group1"), parser.NewString("group1")}),
				},
				{
					NewGroupCell([]parser.Primary{parser.NewString("2"), parser.NewString("4")}),
					NewGroupCell([]parser.Primary{parser.NewString("2"), parser.NewString("4")}),
					NewGroupCell([]parser.Primary{parser.NewString("str2"), parser.NewString("str4")}),
					NewGroupCell([]parser.Primary{parser.NewString("group2"), parser.NewString("group2")}),
				},
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.Function{
					Name: "sum",
					Args: []parser.Expression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				},
				RHS:      parser.NewInteger(5),
				Operator: ">",
			},
		},
		Result: []int{1},
	},
	{
		Name: "Having Filter Error",
		View: &View{
			Header: []HeaderField{
				{
					Reference: "table1",
					Column:    INTERNAL_ID_COLUMN,
				},
				{
					Reference: "table1",
					Column:    "column1",
					FromTable: true,
				},
				{
					Reference: "table1",
					Column:    "column2",
					FromTable: true,
				},
				{
					Reference:  "table1",
					Column:     "column3",
					FromTable:  true,
					IsGroupKey: true,
				},
			},
			isGrouped: true,
			Records: []Record{
				{
					NewGroupCell([]parser.Primary{parser.NewString("1"), parser.NewString("3")}),
					NewGroupCell([]parser.Primary{parser.NewString("1"), parser.NewString("3")}),
					NewGroupCell([]parser.Primary{parser.NewString("str1"), parser.NewString("str3")}),
					NewGroupCell([]parser.Primary{parser.NewString("group1"), parser.NewString("group1")}),
				},
				{
					NewGroupCell([]parser.Primary{parser.NewString("2"), parser.NewString("4")}),
					NewGroupCell([]parser.Primary{parser.NewString("2"), parser.NewString("4")}),
					NewGroupCell([]parser.Primary{parser.NewString("str2"), parser.NewString("str4")}),
					NewGroupCell([]parser.Primary{parser.NewString("group2"), parser.NewString("group2")}),
				},
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.Function{
					Name: "sum",
					Args: []parser.Expression{parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
				},
				RHS:      parser.NewInteger(5),
				Operator: ">",
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Having Not Grouped",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2", "column3"}),
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewString("group2"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("4"),
					parser.NewString("str4"),
					parser.NewString("group2"),
				}),
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.Function{
					Name: "sum",
					Args: []parser.Expression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				},
				RHS:      parser.NewInteger(5),
				Operator: ">",
			},
		},
		Result: []int{0},
		Record: []Record{
			{
				NewGroupCell([]parser.Primary{parser.NewInteger(1), parser.NewInteger(2)}),
				NewGroupCell([]parser.Primary{parser.NewString("2"), parser.NewString("4")}),
				NewGroupCell([]parser.Primary{parser.NewString("str2"), parser.NewString("str4")}),
				NewGroupCell([]parser.Primary{parser.NewString("group2"), parser.NewString("group2")}),
			},
		},
	},
}

func TestView_Having(t *testing.T) {
	for _, v := range viewHavingTests {
		err := v.View.Having(v.Having)
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
		if !reflect.DeepEqual(v.View.filteredIndices, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.filteredIndices, v.Result)
		}
		if v.Record != nil {
			if !reflect.DeepEqual(v.View.Records, v.Record) {
				t.Errorf("%s: result = %s, want %s", v.Name, v.View.Records, v.Record)
			}
		}
	}
}

var viewSelectTests = []struct {
	Name   string
	View   *View
	Select parser.SelectClause
	Result *View
	Error  string
}{
	{
		Name: "Select",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Number: 2, FromTable: true},
				{Reference: "table2", Column: INTERNAL_ID_COLUMN},
				{Reference: "table2", Column: "column3", Number: 1, FromTable: true},
				{Reference: "table2", Column: "column4", Number: 2, FromTable: true},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(2),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(3),
					parser.NewString("1"),
					parser.NewString("str44"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2"}},
				parser.Field{Object: parser.AllColumns{}},
				parser.Field{Object: parser.NewInteger(1), Alias: parser.Identifier{Literal: "a"}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2a"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: parser.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: parser.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21a"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Alias: "c2", Number: 2, FromTable: true},
				{Reference: "table2", Column: INTERNAL_ID_COLUMN},
				{Reference: "table2", Column: "column3", Alias: "t21", Number: 1, FromTable: true},
				{Reference: "table2", Column: "column4", Number: 2, FromTable: true},
				{Column: "1", Alias: "a"},
				{Column: "", Alias: "c2a"},
				{Column: "", Alias: "t21a"},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str22"),
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewString("2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(2),
					parser.NewString("3"),
					parser.NewString("str33"),
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewString("3"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(3),
					parser.NewString("1"),
					parser.NewString("str44"),
					parser.NewInteger(1),
					parser.NewString("str1"),
					parser.NewString("1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str22"),
					parser.NewInteger(1),
					parser.NewString("str2"),
					parser.NewString("2"),
				}),
			},
			selectFields: []int{2, 1, 2, 4, 5, 6, 7, 4, 8},
		},
	},
	{
		Name: "Select Distinct",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table2", Column: INTERNAL_ID_COLUMN},
				{Reference: "table2", Column: "column3", FromTable: true},
				{Reference: "table2", Column: "column4", FromTable: true},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(2),
					parser.NewString("3"),
					parser.NewString("str33"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
					parser.NewInteger(3),
					parser.NewString("4"),
					parser.NewString("str44"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("2"),
					parser.NewString("str2"),
					parser.NewInteger(1),
					parser.NewString("2"),
					parser.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.NewInteger(1), Alias: parser.Identifier{Literal: "a"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", FromTable: true},
				{Column: "1", Alias: "a"},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("1"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("2"),
					parser.NewInteger(1),
				}),
			},
			selectFields: []int{0, 1},
		},
	},
	{
		Name: "Select Aggregate Function",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{
					Object: parser.Function{
						Name: "sum",
						Args: []parser.Expression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Column: "sum(column1)", Alias: "sum(column1)"},
			},
			Records: []Record{
				{
					NewGroupCell([]parser.Primary{parser.NewInteger(1), parser.NewInteger(2)}),
					NewGroupCell([]parser.Primary{parser.NewString("1"), parser.NewString("2")}),
					NewGroupCell([]parser.Primary{parser.NewString("str1"), parser.NewString("str2")}),
					NewCell(parser.NewInteger(3)),
				},
			},
			selectFields: []int{3},
		},
	},
	{
		Name: "Select Aggregate Function Not Group Key Error",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.Function{
						Name: "sum",
						Args: []parser.Expression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
					},
				},
			},
		},
		Error: "[L:- C:-] field column2 is not a group key",
	},
	{
		Name: "Select Analytic Function",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							Partition: parser.Partition{
								PartitionBy: "partition by",
								Values: []parser.Expression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.Expression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: "column1", Number: 1, FromTable: true},
				{Reference: "table1", Column: "column2", Number: 2, FromTable: true},
				{Column: "row_number() over (partition by column1 order by column2)", Alias: "rownum"},
			},
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
					parser.NewInteger(3),
				}),
			},
			selectFields: []int{0, 1, 2},
		},
	},
	{
		Name: "Select Analytic Function Not Exist Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "notexist",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							Partition: parser.Partition{
								PartitionBy: "partition by",
								Values: []parser.Expression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.Expression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "[L:- C:-] function notexist does not exist",
	},
	{
		Name: "Select Analytic Function Order Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							Partition: parser.Partition{
								PartitionBy: "partition by",
								Values: []parser.Expression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.Expression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "Select Analytic Function Partition Value Error",
		View: &View{
			Header: NewHeaderWithoutId("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(2),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(3),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(5),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("a"),
					parser.NewInteger(1),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewString("b"),
					parser.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.Expression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						Over: "over",
						AnalyticClause: parser.AnalyticClause{
							Partition: parser.Partition{
								PartitionBy: "partition by",
								Values: []parser.Expression{
									parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								OrderBy: "order by",
								Items: []parser.Expression{
									parser.OrderItem{
										Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "rownum"},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_Select(t *testing.T) {
	DefineAnalyticFunctions()

	for _, v := range viewSelectTests {
		err := v.View.Select(v.Select)
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
		if !reflect.DeepEqual(v.View.Header, v.Result.Header) {
			t.Errorf("%s: header = %s, want %s", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.Records, v.Result.Records) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.Records, v.Result.Records)
		}
		if !reflect.DeepEqual(v.View.selectFields, v.Result.selectFields) {
			t.Errorf("%s: select indices = %s, want %s", v.Name, v.View.selectFields, v.Result.selectFields)
		}
	}
}

var viewOrderByTests = []struct {
	Name    string
	View    *View
	OrderBy parser.OrderByClause
	Result  *View
	Error   string
}{
	{
		Name: "Order By",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table1", Column: "column3", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("3"),
					parser.NewString("2"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("4"),
					parser.NewString("3"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("4"),
					parser.NewString("3"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("3"),
					parser.NewNull(),
				}),
				NewRecord(5, []parser.Primary{
					parser.NewNull(),
					parser.NewString("2"),
					parser.NewString("4"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.Expression{
				parser.OrderItem{
					Value: parser.Identifier{Literal: "column1"},
				},
				parser.OrderItem{
					Value:     parser.Identifier{Literal: "column2"},
					Direction: parser.Token{Token: parser.DESC, Literal: "desc"},
				},
				parser.OrderItem{
					Value: parser.Identifier{Literal: "column3"},
				},
				parser.OrderItem{
					Value: parser.NewInteger(1),
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
				{Reference: "table1", Column: "column3", FromTable: true},
				{Column: "1"},
			},
			Records: []Record{
				NewRecord(5, []parser.Primary{
					parser.NewNull(),
					parser.NewString("2"),
					parser.NewString("4"),
					parser.NewInteger(1),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("4"),
					parser.NewString("3"),
					parser.NewInteger(1),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("4"),
					parser.NewString("3"),
					parser.NewInteger(1),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("3"),
					parser.NewNull(),
					parser.NewInteger(1),
				}),
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("3"),
					parser.NewString("2"),
					parser.NewInteger(1),
				}),
			},
		},
	},
	{
		Name: "Order By With Null Positions",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("2"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("2"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewNull(),
					parser.NewString("2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("1"),
					parser.NewNull(),
				}),
				NewRecord(5, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("2"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.Expression{
				parser.OrderItem{
					Value:    parser.Identifier{Literal: "column1"},
					Position: parser.Token{Token: parser.LAST, Literal: "last"},
				},
				parser.OrderItem{
					Value:    parser.Identifier{Literal: "column2"},
					Position: parser.Token{Token: parser.FIRST, Literal: "first"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(4, []parser.Primary{
					parser.NewString("1"),
					parser.NewNull(),
				}),
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("2"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("2"),
				}),
				NewRecord(5, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("2"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewNull(),
					parser.NewString("2"),
				}),
			},
		},
	},
}

func TestView_OrderBy(t *testing.T) {
	for _, v := range viewOrderByTests {
		err := v.View.OrderBy(v.OrderBy)
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
		if !reflect.DeepEqual(v.View.Header, v.Result.Header) {
			t.Errorf("%s: header = %s, want %s", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.Records, v.Result.Records) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.Records, v.Result.Records)
		}
	}
}

var viewLimitTests = []struct {
	Name   string
	View   *View
	Limit  parser.LimitClause
	Result *View
	Error  string
}{
	{
		Name: "Limit",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewInteger(2)},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Limit With Ties",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
			sortIndices: []int{1, 2},
		},
		Limit: parser.LimitClause{Value: parser.NewInteger(2), With: parser.LimitWith{Type: parser.Token{Token: parser.TIES}}},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
			sortIndices: []int{1, 2},
		},
	},
	{
		Name: "Limit By Percentage",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
			offset: 1,
		},
		Limit: parser.LimitClause{Value: parser.NewFloat(50.5), Percent: "percent"},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			offset: 1,
		},
	},
	{
		Name: "Limit By Over 100 Percentage",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewFloat(150), Percent: "percent"},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
	},
	{
		Name: "Limit By Negative Percentage",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("3"),
					parser.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewFloat(-10), Percent: "percent"},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{},
		},
	},
	{
		Name: "Limit Greater Than Records",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewInteger(5)},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Limit Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.Variable{Name: "notexist"}},
		Error: "[L:- C:-] variable notexist is undefined",
	},
	{
		Name: "Limit Value Error",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewString("str")},
		Error: "[L:- C:-] limit number of records 'str' is not an integer value",
	},
	{
		Name: "Limit Negative Value",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewInteger(-1)},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{},
		},
	},
	{
		Name: "Limit By Percentage Value Error",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewString("str"), Percent: "percent"},
		Error: "[L:- C:-] limit percentage 'str' is not a float value",
	},
}

func TestView_Limit(t *testing.T) {
	for _, v := range viewLimitTests {
		err := v.View.Limit(v.Limit)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: view = %s, want %s", v.Name, v.View, v.Result)
		}
	}
}

var viewOffsetTests = []struct {
	Name   string
	View   *View
	Offset parser.OffsetClause
	Result *View
	Error  string
}{
	{
		Name: "Offset",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewInteger(3)},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
			offset: 3,
		},
	},
	{
		Name: "Offset Equal To Record Length",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewInteger(4)},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{},
			offset:  4,
		},
	},
	{
		Name: "Offset Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.Variable{Name: "notexist"}},
		Error:  "[L:- C:-] variable notexist is undefined",
	},
	{
		Name: "Offset Value Error",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(3, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(4, []parser.Primary{
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewString("str")},
		Error:  "[L:- C:-] offset number 'str' is not an integer value",
	},
	{
		Name: "Offset Negative Number",
		View: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewInteger(-3)},
		Result: &View{
			Header: []HeaderField{
				{Reference: "table1", Column: INTERNAL_ID_COLUMN},
				{Reference: "table1", Column: "column1", FromTable: true},
				{Reference: "table1", Column: "column2", FromTable: true},
			},
			Records: []Record{
				NewRecord(1, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecord(2, []parser.Primary{
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
			},
			offset: 0,
		},
	},
}

func TestView_Offset(t *testing.T) {
	for _, v := range viewOffsetTests {
		err := v.View.Offset(v.Offset)
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
		if !reflect.DeepEqual(v.View, v.Result) {
			t.Errorf("%s: view = %s, want %s", v.Name, v.View, v.Result)
		}
	}
}

var viewInsertValuesTests = []struct {
	Name       string
	Fields     []parser.Expression
	ValuesList []parser.Expression
	Result     *View
	Error      string
}{
	{
		Name: "InsertValues",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.Expression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(3),
					},
				},
			},
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(4),
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewNull(),
					parser.NewInteger(3),
					parser.NewNull(),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewNull(),
					parser.NewInteger(4),
					parser.NewNull(),
				}),
			},
			OperatedRecords: 2,
		},
	},
	{
		Name: "InsertValues Field Length Does Not Match Error",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		ValuesList: []parser.Expression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] row value should contain exactly 2 values",
	},
	{
		Name: "InsertValues Value Evaluation Error",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.Expression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
	{
		Name: "InsertValues Field Does Not Exist Error",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		ValuesList: []parser.Expression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.Expression{
						parser.NewInteger(3),
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_InsertValues(t *testing.T) {
	view := &View{
		Header: NewHeader("table1", []string{"column1", "column2"}),
		Records: []Record{
			NewRecord(1, []parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecord(2, []parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	for _, v := range viewInsertValuesTests {
		err := view.InsertValues(v.Fields, v.ValuesList, Filter{})
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
			t.Errorf("%s: result = %q, want %q", v.Name, view, v.Result)
		}
	}
}

var viewInsertFromQueryTests = []struct {
	Name   string
	Fields []parser.Expression
	Query  parser.SelectQuery
	Result *View
	Error  string
}{
	{
		Name: "InsertFromQuery",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.NewInteger(3)},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			Records: []Record{
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(1),
					parser.NewString("1"),
					parser.NewString("str1"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewInteger(2),
					parser.NewString("2"),
					parser.NewString("str2"),
				}),
				NewRecordWithoutId([]parser.Primary{
					parser.NewNull(),
					parser.NewInteger(3),
					parser.NewNull(),
				}),
			},
			OperatedRecords: 1,
		},
	},
	{
		Name: "InsertFromQuery Field Lenght Does Not Match Error",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.NewInteger(3)},
					},
				},
			},
		},
		Error: "[L:- C:-] select query should return exactly 2 fields",
	},
	{
		Name: "Insert Values Query Exuecution Error",
		Fields: []parser.Expression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.Expression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
					},
				},
			},
		},
		Error: "[L:- C:-] field notexist does not exist",
	},
}

func TestView_InsertFromQuery(t *testing.T) {
	view := &View{
		Header: NewHeader("table1", []string{"column1", "column2"}),
		Records: []Record{
			NewRecord(1, []parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecord(2, []parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	for _, v := range viewInsertFromQueryTests {
		err := view.InsertFromQuery(v.Fields, v.Query, Filter{})
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
			t.Errorf("%s: result = %q, want %q", v.Name, view, v.Result)
		}
	}
}

func TestView_Fix(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: INTERNAL_ID_COLUMN},
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: []Record{
			NewRecord(1, []parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecord(2, []parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
		},
		selectFields: []int{2},
	}
	expect := &View{
		Header: NewHeaderWithoutId("table1", []string{"column2"}),
		Records: []Record{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("str1"),
			}),
		},
		selectFields: []int(nil),
	}

	view.Fix()
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("fix: view = %s, want %s", view, expect)
	}
}

func TestView_Union(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{Reference: "table2", Column: "column3", FromTable: true},
			{Reference: "table2", Column: "column4", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("3"),
				parser.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("3"),
				parser.NewString("str3"),
			}),
		},
	}

	view.Union(calcView, false)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union: view = %s, want %s", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("3"),
				parser.NewString("str3"),
			}),
		},
	}

	view.Union(calcView, true)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union all: view = %s, want %s", view, expect)
	}
}

func TestView_Except(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{Reference: "table2", Column: "column3", FromTable: true},
			{Reference: "table2", Column: "column4", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("3"),
				parser.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
		},
	}

	view.Except(calcView, false)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except: view = %s, want %s", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
		},
	}

	view.Except(calcView, true)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except all: view = %s, want %s", view, expect)
	}
}

func TestView_Intersect(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{Reference: "table2", Column: "column3", FromTable: true},
			{Reference: "table2", Column: "column4", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("3"),
				parser.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	view.Intersect(calcView, false)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect: view = %s, want %s", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("1"),
				parser.NewString("str1"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
		Records: Records{
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
			NewRecordWithoutId([]parser.Primary{
				parser.NewString("2"),
				parser.NewString("str2"),
			}),
		},
	}

	view.Intersect(calcView, true)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect all: view = %s, want %s", view, expect)
	}
}

func TestView_FieldIndex(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
	}
	fieldRef := parser.FieldReference{
		Column: parser.Identifier{Literal: "column1"},
	}
	expect := 0

	idx, _ := view.FieldIndex(fieldRef)
	if idx != expect {
		t.Errorf("field index = %d, want %d", idx, expect)
	}
}

func TestView_FieldIndices(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table1", Column: "column2", FromTable: true},
		},
	}
	fields := []parser.Expression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
	}
	expect := []int{1, 0}

	indices, _ := view.FieldIndices(fields)
	if !reflect.DeepEqual(indices, expect) {
		t.Errorf("field indices = %d, want %d", indices, expect)
	}
}

func TestView_FieldViewName(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table2", Column: "column2", FromTable: true},
		},
	}
	fieldRef := parser.FieldReference{
		Column: parser.Identifier{Literal: "column1"},
	}
	expect := "table1"

	ref, _ := view.FieldViewName(fieldRef)
	if ref != expect {
		t.Errorf("field reference = %s, want %s", ref, expect)
	}
}

func TestView_InternalRecordId(t *testing.T) {
	view := &View{
		Header: NewHeader("table1", []string{"column1", "column2"}),
		Records: []Record{
			NewRecord(0, []parser.Primary{parser.NewInteger(1), parser.NewString("str1")}),
			NewRecord(1, []parser.Primary{parser.NewInteger(2), parser.NewString("str2")}),
			NewRecord(2, []parser.Primary{parser.NewInteger(3), parser.NewString("str3")}),
		},
	}
	ref := "table1"
	recordIndex := 1
	expect := 1

	id, _ := view.InternalRecordId(ref, recordIndex)
	if id != expect {
		t.Errorf("field internal id = %d, want %d", id, expect)
	}

	view.Records[1][0] = NewCell(parser.NewNull())
	expectError := "[L:- C:-] internal record id is empty"
	_, err := view.InternalRecordId(ref, recordIndex)
	if err.Error() != expectError {
		t.Errorf("error = %q, want error %q", err, expectError)
	}

	view = &View{
		Header: []HeaderField{
			{Reference: "table1", Column: "column1", FromTable: true},
			{Reference: "table2", Column: "column2", FromTable: true},
		},
	}
	expectError = "[L:- C:-] internal record id does not exist"
	_, err = view.InternalRecordId(ref, recordIndex)
	if err.Error() != expectError {
		t.Errorf("error = %q, want error %q", err, expectError)
	}
}
