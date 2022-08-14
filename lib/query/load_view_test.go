package query

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/option"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"
)

var loadViewTests = []struct {
	Name               string
	TestCache          bool
	Encoding           text.Encoding
	NoHeader           bool
	From               parser.FromClause
	ForUpdate          bool
	UseInternalId      bool
	Stdin              string
	ImportFormat       option.Format
	Delimiter          rune
	AllowUnevenFields  bool
	DelimiterPositions []int
	SingleLine         bool
	JsonQuery          string
	WithoutNull        bool
	Scope              *ReferenceScope
	Result             *View
	ResultScope        *ReferenceScope
	Error              string
}{
	{
		Name: "Dual View",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Dual{}},
			},
		},
		Result: &View{
			Header: []HeaderField{{}},
			RecordSet: []Record{
				{
					NewCell(value.NewNull()),
				},
			},
		},
	},
	{
		Name: "Dual View With Omitted FromClause",
		From: parser.FromClause{},
		Result: &View{
			Header: []HeaderField{{}},
			RecordSet: []Record{
				{
					NewCell(value.NewNull()),
				},
			},
		},
	},
	{
		Name: "LoadView File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1.csv"},
				},
			},
		},
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
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView File ForUpdate",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1.csv"},
				},
			},
		},
		ForUpdate: true,
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
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ForUpdate: true,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{{
			scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
			},
		}}, time.Time{}, nil),
	},
	{
		Name:      "LoadView from Cached View",
		TestCache: true,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1"},
				},
			},
		},
		ForUpdate: true,
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
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ForUpdate: true,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{{
			scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
			},
		}}, time.Time{}, nil),
	},
	{
		Name: "LoadView File with UTF-8 BOM",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table1_bom.csv"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1_bom", []string{"column1", "column2"}),
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
				Path:      "table1_bom.csv",
				Delimiter: ',',
				Encoding:  text.UTF8M,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1_BOM": strings.ToUpper(GetTestFilePath("table1_bom.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView with Parentheses",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Parentheses{
					Expr: parser.Table{
						Object: parser.Identifier{Literal: "table1.csv"},
					},
				},
			},
		},
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
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table1.csv",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin: "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeader("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": "STDIN",
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView From Stdin ForUpdate",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:     "column1,column2\n1,\"str1\"",
		ForUpdate: true,
		Result: &View{
			Header: NewHeader("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"T": "STDIN",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView TSV From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "column1\tcolumn2\n1\t\"str1\"",
		ImportFormat: option.TSV,
		Result: &View{
			Header: NewHeader("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Format:    option.TSV,
				Delimiter: '\t',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": "STDIN",
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView From Stdin With Internal Id",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		UseInternalId: true,
		Stdin:         "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeaderWithId("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(0, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"T": "STDIN",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "{\"key\":[{\"column1\": 1, \"column2\": \"str1\"}]}",
		ImportFormat: option.JSON,
		JsonQuery:    "key{}",
		Result: &View{
			Header: NewHeader("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewFloat(1),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Delimiter: ',',
				JsonQuery: "key{}",
				Format:    option.JSON,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": "STDIN",
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json Lines From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "{\"column1\": 1, \"column2\": \"str1\"}\n{\"column1\": 2, \"column2\": \"str2\"}",
		ImportFormat: option.JSONL,
		JsonQuery:    "",
		Result: &View{
			Header: NewHeader("t", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewFloat(1),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewFloat(2),
					value.NewString("str2"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Delimiter: ',',
				JsonQuery: "",
				Format:    option.JSONL,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": "STDIN",
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json Lines From Stdin, Json Structure Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "{\"column1\": 1, \"column2\": \"str1\"}\n\"str\"",
		ImportFormat: option.JSONL,
		JsonQuery:    "",
		Error:        "json lines must be an array of objects",
	},
	{
		Name: "LoadView JsonH From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "[{\"item1\": \"value\\u00221\",\"item2\": 1},{\"item1\": \"value2\",\"item2\": 2}]",
		ImportFormat: option.JSON,
		JsonQuery:    "{}",
		Result: &View{
			Header: NewHeader("t", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value\"1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:       "STDIN",
				Delimiter:  ',',
				JsonQuery:  "{}",
				Format:     option.JSON,
				Encoding:   text.UTF8,
				LineBreak:  text.LF,
				JsonEscape: json.HexDigits,
				ViewType:   ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"T": "STDIN",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView JsonA From Stdin",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "[{\"item1\": \"\\u0076\\u0061\\u006c\\u0075\\u0065\\u0031\",\"item2\": 1},{\"item1\": \"\\u0076\\u0061\\u006c\\u0075\\u0065\\u0032\",\"item2\": 2}]",
		ImportFormat: option.JSON,
		JsonQuery:    "{}",
		Result: &View{
			Header: NewHeader("t", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:       "STDIN",
				Delimiter:  ',',
				JsonQuery:  "{}",
				Format:     option.JSON,
				Encoding:   text.UTF8,
				LineBreak:  text.LF,
				JsonEscape: json.AllWithHexDigits,
				ViewType:   ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"T": "STDIN",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json From Stdin Json Query Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin:        "{\"key\":[{\"column1\": 1, \"column2\": \"str1\"}]}",
		ImportFormat: option.JSON,
		JsonQuery:    "key{",
		Error:        "json loading error: column 4: unexpected termination",
	},
	{
		Name:         "LoadView Fixed-Length Text File",
		ImportFormat: option.FIXED,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "fixed_length.txt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("fixed_length", []string{"column1", "__@2__"}),
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
				Path:               "fixed_length.txt",
				Delimiter:          ',',
				DelimiterPositions: []int{7, 12},
				Format:             option.FIXED,
				NoHeader:           false,
				Encoding:           text.UTF8,
				LineBreak:          text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"FIXED_LENGTH": strings.ToUpper(GetTestFilePath("fixed_length.txt")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name:         "LoadView Fixed-Length Text File NoHeader",
		NoHeader:     true,
		ImportFormat: option.FIXED,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "fixed_length.txt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("fixed_length", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("column1"),
					value.NewNull(),
				}),
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
				Path:               "fixed_length.txt",
				Delimiter:          ',',
				DelimiterPositions: []int{7, 12},
				Format:             option.FIXED,
				NoHeader:           true,
				Encoding:           text.UTF8,
				LineBreak:          text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"FIXED_LENGTH": strings.ToUpper(GetTestFilePath("fixed_length.txt")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name:               "LoadView Fixed-Length Text File Position Error",
		ImportFormat:       option.FIXED,
		DelimiterPositions: []int{6, 2},
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "fixed_length.txt"},
				},
			},
		},
		Error: fmt.Sprintf("data parse error in %s: invalid delimiter position: [6, 2]", GetTestFilePath("fixed_length.txt")),
	},
	{
		Name:  "LoadView From Stdin With Omitted FromClause",
		From:  parser.FromClause{},
		Stdin: "column1,column2\n1,\"str1\"",
		Result: &View{
			Header: NewHeader("STDIN", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "STDIN",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStdin,
			},
		},
		ResultScope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameTempTables: {
					"STDIN": &View{
						FileInfo: &FileInfo{Path: "STDIN"},
					},
				},
			},
		}, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"STDIN": "STDIN",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView From Stdin Broken CSV Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin: "column1,column2\n1\"str1\"",
		Error: "data parse error in STDIN: line 1, column 8: wrong number of fields in line",
	},
	{
		Name: "LoadView From Stdin Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Stdin: "column1,column2\n1,\"str1\"",
		Error: "table name t is a duplicate",
	},
	{
		Name: "LoadView From Stdin ForUpdate Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Stdin{}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		ForUpdate: true,
		Stdin:     "column1,column2\n1,\"str1\"",
		Error:     "table name t is a duplicate",
	},
	{
		Name: "Stdin Empty Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Stdin{},
				},
			},
		},
		Error: "STDIN is empty",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table5"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
							parser.NewTernaryValueFromString("true"),
							parser.NewTernaryValueFromString("true"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString(""),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table5.csv",
				Delimiter: ',',
				Format:    option.CSV,
				Encoding:  text.SJIS,
				LineBreak: text.LF,
				NoHeader:  true,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("table5.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From TSV File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue("\t"),
						Path:          parser.Identifier{Literal: "table3"},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "UTF8"}},
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"column5", "column6"}),
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
				Path:      "table3.tsv",
				Delimiter: '\t',
				Format:    option.TSV,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("table3.tsv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File FormatElement Evaluate Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
						Path:          parser.Identifier{Literal: "table1"},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File FormatElement Is Not Specified",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.CSV, Literal: "csv"},
						Path: parser.Identifier{Literal: "table1"},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid argument for csv: delimiter is not specified",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File FormatElement is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewNullValue(),
						Path:          parser.Identifier{Literal: "table1"},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid delimiter: NULL",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File FormatElement Invalid Delimiter",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue("invalid"),
						Path:          parser.Identifier{Literal: "table1"},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid delimiter: 'invalid'",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table5"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
							parser.NewTernaryValueFromString("true"),
							parser.NewTernaryValueFromString("true"),
							parser.NewStringValue("extra"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "table object csv takes at most 5 arguments",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File 3rd Argument Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table5"},
						Args: []parser.QueryExpression{
							parser.NewTernaryValueFromString("true"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid argument for csv: cannot be converted as a encoding value: TRUE",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File 4th Argument Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table5"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
							parser.NewStringValue("SJIS"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid argument for csv: cannot be converted as a no-header value: 'SJIS'",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File 5th Argument Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table5"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
							parser.NewTernaryValueFromString("true"),
							parser.NewStringValue("SJIS"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid argument for csv: cannot be converted as a without-null value: 'SJIS'",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From CSV File Invalid Encoding Type",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table5"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("INVALID"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid argument for csv: encoding must be one of AUTO|UTF8|UTF8M|UTF16|UTF16BE|UTF16LE|UTF16BEM|UTF16LEM|SJIS",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Fixed-Length File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
						FormatElement: parser.NewStringValue("spaces"),
						Path:          parser.Identifier{Literal: "fixed_length.txt", Quoted: true},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"column1", "__@2__"}),
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
				Path:               "fixed_length.txt",
				Delimiter:          ',',
				DelimiterPositions: []int{7, 12},
				Format:             option.FIXED,
				Encoding:           text.UTF8,
				LineBreak:          text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("fixed_length.txt")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Fixed-Length File with UTF-8 BOM",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
						FormatElement: parser.NewStringValue("spaces"),
						Path:          parser.Identifier{Literal: "fixed_length_bom.txt", Quoted: true},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"column1", "__@2__"}),
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
				Path:               "fixed_length_bom.txt",
				Delimiter:          ',',
				DelimiterPositions: []int{7, 12},
				Format:             option.FIXED,
				Encoding:           text.UTF8M,
				LineBreak:          text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("fixed_length_bom.txt")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Single-Line Fixed-Length File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
						FormatElement: parser.NewStringValue("s[1,5]"),
						Path:          parser.Identifier{Literal: "fixed_length_sl.txt", Quoted: true},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"c1", "c2"}),
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
				Path:               "fixed_length_sl.txt",
				Delimiter:          ',',
				DelimiterPositions: []int{1, 5},
				Format:             option.FIXED,
				Encoding:           text.UTF8,
				LineBreak:          text.LF,
				SingleLine:         true,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("fixed_length_sl.txt")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Fixed-Length File FormatElement Is Not Specified",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.FIXED, Literal: "fixed"},
						Path: parser.Identifier{Literal: "fixed_length.txt", Quoted: true},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid argument for fixed: delimiter positions are not specified",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Fixed-Length File FormatElement is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
						FormatElement: parser.NewNullValue(),
						Path:          parser.Identifier{Literal: "fixed_length.txt", Quoted: true},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid delimiter positions: NULL",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Fixed-Length File Invalid Delimiter Positions",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
						FormatElement: parser.NewStringValue("invalid"),
						Path:          parser.Identifier{Literal: "fixed_length.txt", Quoted: true},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "invalid delimiter positions: 'invalid'",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Fixed-Length File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.FIXED, Literal: "fixed"},
						FormatElement: parser.NewStringValue("spaces"),
						Path:          parser.Identifier{Literal: "fixed_length.txt", Quoted: true},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
							parser.NewTernaryValueFromString("true"),
							parser.NewTernaryValueFromString("true"),
							parser.NewStringValue("extra"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "table object fixed takes at most 5 arguments",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Json File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON, Literal: "json"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table.json",
				Delimiter: ',',
				JsonQuery: "{}",
				Format:    option.JSON,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"JT": strings.ToUpper(GetTestFilePath("table.json")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From JsonH File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON, Literal: "json"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table_h"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value\"1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:       "table_h.json",
				Delimiter:  ',',
				JsonQuery:  "{}",
				Format:     option.JSON,
				Encoding:   text.UTF8,
				LineBreak:  text.LF,
				JsonEscape: json.HexDigits,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"JT": strings.ToUpper(GetTestFilePath("table_h.json")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From JsonA File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON, Literal: "json"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table_a"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:       "table_a.json",
				Delimiter:  ',',
				JsonQuery:  "{}",
				Format:     option.JSON,
				Encoding:   text.UTF8,
				LineBreak:  text.LF,
				JsonEscape: json.AllWithHexDigits,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"JT": strings.ToUpper(GetTestFilePath("table_a.json")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Json File FormatElement Is Not Specified",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.JSON, Literal: "json"},
						Path: parser.Identifier{Literal: "table"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "invalid argument for json: json query is not specified",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Json File FormatElement is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON, Literal: "json"},
						FormatElement: parser.NewNullValue(),
						Path:          parser.Identifier{Literal: "table"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "invalid json query: NULL",
	},
	{
		Name: "LoadView Table Object From Json File Path Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON, Literal: "json"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "notexist"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Json Lines File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSONL, Literal: "jsonl"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table7"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table7.jsonl",
				Delimiter: ',',
				JsonQuery: "{}",
				Format:    option.JSONL,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"JT": strings.ToUpper(GetTestFilePath("table7.jsonl")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From LTSV File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.LTSV, Literal: "ltsv"},
						Path: parser.Identifier{Literal: "table6"},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"f1", "f2", "f3", "f4"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewString("value2"),
					value.NewString("value3"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("value4"),
					value.NewString("value5"),
					value.NewNull(),
					value.NewString("value6"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table6.ltsv",
				Delimiter: ',',
				Format:    option.LTSV,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("table6.ltsv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From LTSV File Without Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.LTSV, Literal: "ltsv"},
						Path: parser.Identifier{Literal: "table6"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("UTF8"),
							parser.NewTernaryValueFromString("true"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"f1", "f2", "f3", "f4"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewString("value2"),
					value.NewString("value3"),
					value.NewString(""),
				}),
				NewRecord([]value.Primary{
					value.NewString("value4"),
					value.NewString("value5"),
					value.NewString(""),
					value.NewString("value6"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table6.ltsv",
				Delimiter: ',',
				Format:    option.LTSV,
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("table6.ltsv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From LTSV File with UTF-8 BOM",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.LTSV, Literal: "ltsv"},
						Path: parser.Identifier{Literal: "table6_bom"},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"f1", "f2", "f3", "f4"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewString("value2"),
					value.NewString("value3"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("value4"),
					value.NewString("value5"),
					value.NewNull(),
					value.NewString("value6"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table6_bom.ltsv",
				Delimiter: ',',
				Format:    option.LTSV,
				Encoding:  text.UTF8M,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": strings.ToUpper(GetTestFilePath("table6_bom.ltsv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView FormatSpecifiedFunction From LTSV File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type: parser.Token{Token: parser.LTSV, Literal: "ltsv"},
						Path: parser.Identifier{Literal: "table6"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("UTF8"),
							parser.NewTernaryValueFromString("true"),
							parser.NewStringValue("extra"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "table object ltsv takes exactly 3 arguments",
	},
	{
		Name: "LoadView FormatSpecifiedFunction Invalid Object Type",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: 0, Literal: "invalid"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.Identifier{Literal: "table"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "invalid table object: invalid",
	},
	{
		Name: "LoadView FormatSpecifiedFunction From Json File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON, Literal: "json"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table"},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
						},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "table object json takes exactly 2 arguments",
	},
	{
		Name: "LoadView File Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "notexist"},
				},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "LoadView From File Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Identifier{Literal: "table2"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Error: "table name t is a duplicate",
	},
	{
		Name: "LoadView From File ForUpdate Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Identifier{Literal: "table2"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		ForUpdate: true,
		Error:     "table name t is a duplicate",
	},
	{
		Name: "LoadView From File Inline Table",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "it"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Scope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameInlineTables: {
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
			},
		}, time.Time{}, nil),
		Result: &View{
			Header: NewHeader("t", []string{"c1", "c2", "num"}),
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
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"T": "",
				},
			},
			{
				scopeNameInlineTables: {
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
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView From File Inline Table Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{Object: parser.Identifier{Literal: "it"}, Alias: parser.Identifier{Literal: "t"}},
			},
		},
		Scope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameInlineTables: {
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
			},
		}, time.Time{}, nil),
		Error: "table name t is a duplicate",
	},
	{
		Name:     "LoadView SJIS File",
		Encoding: text.SJIS,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_sjis"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table_sjis", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("日本語"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE_SJIS": strings.ToUpper(GetTestFilePath("table_sjis.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name:     "LoadView No Header File",
		NoHeader: true,
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_noheader"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table_noheader", []string{"c1", "c2"}),
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
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE_NOHEADER": strings.ToUpper(GetTestFilePath("table_noheader.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Multiple File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Cross Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("4"),
					value.NewString("str44"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Lateral Cross Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Lateral: parser.Token{Token: parser.LATERAL},
							Object: parser.Subquery{
								Query: parser.SelectQuery{
									SelectEntity: parser.SelectEntity{
										SelectClause: parser.SelectClause{
											Fields: []parser.QueryExpression{
												parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
												parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
											},
										},
										FromClause: parser.FromClause{
											Tables: []parser.QueryExpression{
												parser.Table{Object: parser.Identifier{Literal: "table2"}},
											},
										},
										WhereClause: parser.WhereClause{
											Filter: parser.Comparison{
												LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
												RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
												Operator: parser.Token{Token: '=', Literal: "="},
											},
										},
									},
								},
							},
							Alias: parser.Identifier{
								Literal: "t2",
							},
						},
						JoinType: parser.Token{Token: parser.CROSS, Literal: "cross"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "t2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "t2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Inner Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Inner Join Using Condition",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Condition: parser.JoinCondition{
							Using: []parser.QueryExpression{
								parser.Identifier{Literal: "column1"},
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Column: "column1", IsFromTable: true, IsJoinColumn: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table1b", Column: "column2b", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("str3b"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1":  strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE1B": strings.ToUpper(GetTestFilePath("table1b.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Outer Join",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewNull(),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("3"),
					value.NewString("str33"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1": strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE2": strings.ToUpper(GetTestFilePath("table2.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Outer Join Natural",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Direction: parser.Token{Token: parser.RIGHT, Literal: "right"},
						Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Column: "column1", IsFromTable: true, IsJoinColumn: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table1b", Column: "column2b", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("str3b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("4"),
					value.NewNull(),
					value.NewString("str4b"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1":  strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE1B": strings.ToUpper(GetTestFilePath("table1b.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Full Outer Join Natural",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Direction: parser.Token{Token: parser.FULL, Literal: "full"},
						Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{Column: "column1", IsFromTable: true, IsJoinColumn: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table1b", Column: "column2b", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("str3b"),
				}),
				NewRecord([]value.Primary{
					value.NewString("4"),
					value.NewNull(),
					value.NewString("str4b"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE1":  strings.ToUpper(GetTestFilePath("table1.csv")),
				"TABLE1B": strings.ToUpper(GetTestFilePath("table1b.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Incorrect LATERAL Usage Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Lateral: parser.Token{Token: parser.LATERAL},
							Object: parser.Subquery{
								Query: parser.SelectQuery{
									SelectEntity: parser.SelectEntity{
										SelectClause: parser.SelectClause{
											Fields: []parser.QueryExpression{
												parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}}},
												parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column4"}}},
											},
										},
										FromClause: parser.FromClause{
											Tables: []parser.QueryExpression{
												parser.Table{Object: parser.Identifier{Literal: "table2"}},
											},
										},
										WhereClause: parser.WhereClause{
											Filter: parser.Comparison{
												LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
												RHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
												Operator: parser.Token{Token: '=', Literal: "="},
											},
										},
									},
								},
							},
						},
						Direction: parser.Token{Token: parser.FULL, Literal: "full"},
						Natural:   parser.Token{Token: parser.NATURAL, Literal: "natural"},
					},
				},
			},
		},
		Error: "LATERAL cannot to be used in a RIGHT or FULL outer join",
	},
	{
		Name: "Join Left Side Table File Not Exist Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
		Error: "file notexist does not exist",
	},
	{
		Name: "Join Right Side Table File Not Exist Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
		Error: "file notexist does not exist",
	},
	{
		Name: "LoadView from DATA table function",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableFunction{
						Name: "data",
						Args: []parser.QueryExpression{
							parser.NewStringValue("c1,c2\n1,a\n2,b\n"),
						},
					},
					Alias: parser.Identifier{Literal: "ci"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("ci", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("a"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("b"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "",
				Format:    option.CSV,
				Delimiter: ',',
				JsonQuery: "",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStringObject,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"CI": "",
				},
			},
		}, time.Time{}, nil),
	},
	{ //TODO
		Name: "LoadView Inline Table as FormatSpecifiedFunction",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV, Literal: "csv"},
						FormatElement: parser.NewStringValue(","),
						Path: parser.TableFunction{
							Name: "inline",
							Args: []parser.QueryExpression{
								parser.NewStringValue("table5"),
							},
						},
						Args: []parser.QueryExpression{
							parser.NewStringValue("SJIS"),
							parser.NewTernaryValueFromString("true"),
							parser.NewTernaryValueFromString("true"),
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("t", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString(""),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "",
				Delimiter: ',',
				Format:    option.CSV,
				Encoding:  text.SJIS,
				LineBreak: text.LF,
				NoHeader:  true,
				ViewType:  ViewTypeInlineTable,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"T": "",
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView from Local File as URL",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Url{
						Raw: "file:./table.json",
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table.json",
				Format:    option.JSON,
				Delimiter: ',',
				JsonQuery: "",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeFile,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"JT": strings.ToUpper(GetTestFilePath("table.json")),
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json Inline Table",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewStringValue("{column1, column2}"),
						Path:          parser.NewStringValue("[{\"column1\":1, \"column2\":2},{\"column1\":3, \"column2\":4}]"),
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewFloat(1),
					value.NewFloat(2),
				}),
				NewRecord([]value.Primary{
					value.NewFloat(3),
					value.NewFloat(4),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "",
				Format:    option.JSON,
				Delimiter: ',',
				JsonQuery: "{column1, column2}",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStringObject,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"JT": "",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json Inline Table Query Evaluation Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.FieldReference{Column: parser.Identifier{Literal: "notexists"}},
						Path:          parser.NewStringValue("[{\"column1\":1, \"column2\":2},{\"column1\":3, \"column2\":4}]"),
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "field notexists does not exist",
	},
	{
		Name: "LoadView Json Inline Table JsonText Evaluation Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewStringValue("{column1, column2}"),
						Path:          parser.FieldReference{Column: parser.Identifier{Literal: "notexists"}},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "field notexists does not exist",
	},
	{
		Name: "LoadView Json Inline Table Query is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewNullValue(),
						Path:          parser.NewStringValue("[{\"column1\":1, \"column2\":2},{\"column1\":3, \"column2\":4}]"),
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "invalid json query: NULL",
	},
	{
		Name: "LoadView Json Inline Table JsonText is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewStringValue("{column1, column2}"),
						Path:          parser.NewNullValue(),
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "inline table is empty",
	},
	{
		Name: "LoadView Json Inline Table Loading Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewStringValue("{column1, column2"),
						Path:          parser.NewStringValue("[{\"column1\":1, \"column2\":2},{\"column1\":3, \"column2\":4}]"),
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "json loading error: column 17: unexpected termination",
	},
	{
		Name: "LoadView Json Inline Table From File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_INLINE, Literal: "json_inline"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("jt", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "",
				Format:    option.JSON,
				Delimiter: ',',
				JsonQuery: "{}",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeInlineTable,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"JT": "",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json Inline Table From File with No Alias",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "table"},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table", []string{"item1", "item2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("value1"),
					value.NewFloat(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("value2"),
					value.NewFloat(2),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "",
				Format:    option.JSON,
				Delimiter: ',',
				JsonQuery: "{}",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeInlineTable,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"TABLE": "",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Json Inline Table From File Path Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.JSON_TABLE, Literal: "json_table"},
						FormatElement: parser.NewStringValue("{}"),
						Path:          parser.Identifier{Literal: "notexist"},
					},
					Alias: parser.Identifier{Literal: "jt"},
				},
			},
		},
		Error: "file notexist does not exist",
	},
	{
		Name: "LoadView CSV Inline Table",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.FormatSpecifiedFunction{
						Type:          parser.Token{Token: parser.CSV_INLINE, Literal: "csv_inline"},
						FormatElement: parser.NewStringValue(","),
						Path:          parser.NewStringValue("c1,c2\n1,a\n2,b\n"),
					},
					Alias: parser.Identifier{Literal: "ci"},
				},
			},
		},
		Result: &View{
			Header: NewHeader("ci", []string{"c1", "c2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("a"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("b"),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "",
				Format:    option.CSV,
				Delimiter: ',',
				JsonQuery: "",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeStringObject,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"CI": "",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Subquery",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "table1"}},
									},
								},
							},
						},
					},
				},
			},
		},
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
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{{}}, time.Time{}, nil),
	},
	{
		Name: "LoadView Subquery with Table Name Alias",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
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
			Header: NewHeader("alias", []string{"column1", "column2"}),
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
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{
				scopeNameAliases: {
					"ALIAS": "",
				},
			},
		}, time.Time{}, nil),
	},
	{
		Name: "LoadView Subquery Duplicate Table Name Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{Object: parser.Identifier{Literal: "table1"}, Alias: parser.Identifier{Literal: "t"}},
				parser.Table{
					Object: parser.Subquery{
						Query: parser.SelectQuery{
							SelectEntity: parser.SelectEntity{
								SelectClause: parser.SelectClause{
									Fields: []parser.QueryExpression{
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
										parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
									},
								},
								FromClause: parser.FromClause{
									Tables: []parser.QueryExpression{
										parser.Table{Object: parser.Identifier{Literal: "table1"}},
									},
								},
							},
						},
					},
					Alias: parser.Identifier{Literal: "t"},
				},
			},
		},
		Error: "table name t is a duplicate",
	},
	{
		Name: "LoadView CSV Parse Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_broken.csv"},
				},
			},
		},
		Error: fmt.Sprintf("data parse error in %s: line 3, column 7: wrong number of fields in line", GetTestFilePath("table_broken.csv")),
	},
	{
		Name: "Allow Uneven Field Length",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_broken.csv"},
				},
			},
		},
		AllowUnevenFields: true,
		Result: &View{
			Header: NewHeader("table_broken", []string{"column1", "column2", "__@3__"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewNull(),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table_broken.csv",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE_BROKEN": strings.ToUpper(GetTestFilePath("table_broken.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Allow Uneven Field Length without Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "table_broken.csv"},
				},
			},
		},
		AllowUnevenFields: true,
		WithoutNull:       true,
		Result: &View{
			Header: NewHeader("table_broken", []string{"column1", "column2", "__@3__"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString(""),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString(""),
				}),
			},
			FileInfo: &FileInfo{
				Path:      "table_broken.csv",
				Delimiter: ',',
				Encoding:  text.UTF8,
				LineBreak: text.LF,
			},
		},
		ResultScope: GenerateReferenceScope(nil, []map[string]map[string]interface{}{
			{scopeNameAliases: {
				"TABLE_BROKEN": strings.ToUpper(GetTestFilePath("table_broken.csv")),
			}},
		}, time.Time{}, nil),
	},
	{
		Name: "Inner Join Join Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
								LHS:      parser.FieldReference{View: parser.Identifier{Literal: "table1"}, Column: parser.Identifier{Literal: "notexist"}},
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "column3"}},
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					},
				},
			},
		},
		Error: "field table1.notexist does not exist",
	},
	{
		Name: "Outer Join Join Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
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
								RHS:      parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.Identifier{Literal: "notexist"}},
								Operator: parser.Token{Token: '=', Literal: "="},
							},
						},
					},
				},
			},
		},
		Error: "field table2.notexist does not exist",
	},
	{
		Name: "Inner Join Using Condition Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Join{
						Table: parser.Table{
							Object: parser.Identifier{Literal: "table1"},
						},
						JoinTable: parser.Table{
							Object: parser.Identifier{Literal: "table1b"},
						},
						Condition: parser.JoinCondition{
							Using: []parser.QueryExpression{
								parser.Identifier{Literal: "notexist"},
							},
						},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestLoadView(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)
		_ = TestTx.Session.SetStdin(os.Stdin)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	ctx := context.Background()

	for _, v := range loadViewTests {
		TestTx.UnlockStdin()
		_ = TestTx.CachedViews.Clean(TestTx.FileContainer)

		_ = TestTx.Session.SetStdin(os.Stdin)
		TestTx.Flags.ImportOptions.Format = v.ImportFormat
		TestTx.Flags.ImportOptions.Delimiter = ','
		if v.Delimiter != 0 {
			TestTx.Flags.ImportOptions.Delimiter = v.Delimiter
		}
		TestTx.Flags.ImportOptions.AllowUnevenFields = v.AllowUnevenFields
		TestTx.Flags.ImportOptions.DelimiterPositions = v.DelimiterPositions
		TestTx.Flags.ImportOptions.SingleLine = v.SingleLine
		TestTx.Flags.ImportOptions.JsonQuery = v.JsonQuery
		TestTx.Flags.ImportOptions.NoHeader = v.NoHeader
		TestTx.Flags.ImportOptions.WithoutNull = v.WithoutNull
		if v.Encoding != text.AUTO {
			TestTx.Flags.ImportOptions.Encoding = v.Encoding
		} else {
			TestTx.Flags.ImportOptions.Encoding = text.UTF8
		}

		if 0 < len(v.Stdin) {
			_ = TestTx.Session.SetStdin(NewInput(strings.NewReader(v.Stdin)))
		} else {
			_ = TestTx.Session.SetStdin(nil)
		}

		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		queryScope := v.Scope.CreateNode()
		view, err := LoadView(ctx, queryScope, v.From.Tables, v.ForUpdate, v.UseInternalId)
		if v.TestCache {
			queryScope.nodes[len(queryScope.nodes)-1].Clear()
			view, err = LoadView(ctx, queryScope, v.From.Tables, v.ForUpdate, v.UseInternalId)
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
				t.Errorf("%s: FileInfo.Path = %q, want %q", v.Name, filepath.Base(view.FileInfo.Path), filepath.Base(v.Result.FileInfo.Path))
			}
			if view.FileInfo.Format != v.Result.FileInfo.Format {
				t.Errorf("%s: FileInfo.Format = %s, want %s", v.Name, view.FileInfo.Format, v.Result.FileInfo.Format)
			}
			if view.FileInfo.Delimiter != v.Result.FileInfo.Delimiter {
				t.Errorf("%s: FileInfo.Delimiter = %q, want %q", v.Name, view.FileInfo.Delimiter, v.Result.FileInfo.Delimiter)
			}
			if !reflect.DeepEqual(view.FileInfo.DelimiterPositions, v.Result.FileInfo.DelimiterPositions) {
				t.Errorf("%s: FileInfo.DelimiterPositions = %v, want %v", v.Name, view.FileInfo.DelimiterPositions, v.Result.FileInfo.DelimiterPositions)
			}
			if view.FileInfo.JsonQuery != v.Result.FileInfo.JsonQuery {
				t.Errorf("%s: FileInfo.JsonQuery = %q, want %q", v.Name, view.FileInfo.JsonQuery, v.Result.FileInfo.JsonQuery)
			}
			if view.FileInfo.Encoding != v.Result.FileInfo.Encoding {
				t.Errorf("%s: FileInfo.Encoding = %s, want %s", v.Name, view.FileInfo.Encoding, v.Result.FileInfo.Encoding)
			}
			if view.FileInfo.LineBreak != v.Result.FileInfo.LineBreak {
				t.Errorf("%s: FileInfo.LineBreak = %s, want %s", v.Name, view.FileInfo.LineBreak, v.Result.FileInfo.LineBreak)
			}
			if view.FileInfo.NoHeader != v.Result.FileInfo.NoHeader {
				t.Errorf("%s: FileInfo.NoHeader = %t, want %t", v.Name, view.FileInfo.NoHeader, v.Result.FileInfo.NoHeader)
			}
			if view.FileInfo.PrettyPrint != v.Result.FileInfo.PrettyPrint {
				t.Errorf("%s: FileInfo.PrettyPrint = %t, want %t", v.Name, view.FileInfo.PrettyPrint, v.Result.FileInfo.PrettyPrint)
			}
			if view.FileInfo.ForUpdate != v.Result.FileInfo.ForUpdate {
				t.Errorf("%s: FileInfo.ForUpdate = %t, want %t", v.Name, view.FileInfo.ForUpdate, v.Result.FileInfo.ForUpdate)
			}
			if view.FileInfo.ViewType != v.Result.FileInfo.ViewType {
				t.Errorf("%s: FileInfo.ViewType = %d, want %d", v.Name, view.FileInfo.ViewType, v.Result.FileInfo.ViewType)
			}
		}
		if view.FileInfo != nil {
			_ = TestTx.FileContainer.Close(view.FileInfo.Handler)
			view.FileInfo = nil
		}
		v.Result.FileInfo = nil

		if v.ResultScope == nil {
			v.ResultScope = NewReferenceScope(TestTx).CreateNode()
		}

		if !NodeScopeListEqual(queryScope.nodes, v.ResultScope.nodes) {
			t.Errorf("%s: node list = %v, want %v", v.Name, queryScope.nodes, v.ResultScope.nodes)
		}
		for i := range queryScope.Blocks {
			if !reflect.DeepEqual(queryScope.Blocks[i].TemporaryTables.Keys(), v.ResultScope.Blocks[i].TemporaryTables.Keys()) {
				t.Errorf("%s: temp view list = %v, want %v", v.Name, queryScope.Blocks[i].TemporaryTables.Keys(), v.ResultScope.Blocks[i].TemporaryTables.Keys())
			}
		}

		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: \n result = %v,\n expect = %v", v.Name, view, v.Result)
		}
	}
}
