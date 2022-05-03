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

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/json"
)

var viewLoadTests = []struct {
	Name               string
	Encoding           text.Encoding
	NoHeader           bool
	From               parser.FromClause
	ForUpdate          bool
	UseInternalId      bool
	Stdin              string
	ImportFormat       cmd.Format
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
		ImportFormat: cmd.TSV,
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
				Format:    cmd.TSV,
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
		ImportFormat: cmd.JSON,
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
				Format:    cmd.JSON,
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
		ImportFormat: cmd.JSONL,
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
				Format:    cmd.JSONL,
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
		ImportFormat: cmd.JSONL,
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
		ImportFormat: cmd.JSON,
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
				Format:     cmd.JSON,
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
		ImportFormat: cmd.JSON,
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
				Format:     cmd.JSON,
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
		ImportFormat: cmd.JSON,
		JsonQuery:    "key{",
		Error:        "json loading error: column 4: unexpected termination",
	},
	{
		Name:         "LoadView Fixed-Length Text File",
		ImportFormat: cmd.FIXED,
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
				Format:             cmd.FIXED,
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
		ImportFormat: cmd.FIXED,
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
				Format:             cmd.FIXED,
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
		ImportFormat:       cmd.FIXED,
		DelimiterPositions: []int{6, 2},
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.Identifier{Literal: "fixed_length.txt"},
				},
			},
		},
		Error: fmt.Sprintf("data parse error in file %s: invalid delimiter position: [6, 2]", GetTestFilePath("fixed_length.txt")),
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
		Error: "data parse error in file STDIN: line 1, column 8: wrong number of fields in line",
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
		Name: "LoadView TableObject From CSV File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.CSV,
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
		Name: "LoadView TableObject From TSV File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.TSV,
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
		Name: "LoadView TableObject From CSV File FormatElement Evaluate Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File FormatElement Is Not Specified",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File FormatElement is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File FormatElement Invalid Delimiter",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File 3rd Argument Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File 4th Argument Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File 5th Argument Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From CSV File Invalid Encoding Type",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Fixed-Length File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:             cmd.FIXED,
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
		Name: "LoadView TableObject From Fixed-Length File with UTF-8 BOM",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:             cmd.FIXED,
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
		Name: "LoadView TableObject From Single-Line Fixed-Length File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:             cmd.FIXED,
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
		Name: "LoadView TableObject From Fixed-Length File FormatElement Is Not Specified",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Fixed-Length File FormatElement is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Fixed-Length File Invalid Delimiter Positions",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Fixed-Length File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Json File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.JSON,
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
		Name: "LoadView TableObject From JsonH File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:     cmd.JSON,
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
		Name: "LoadView TableObject From JsonA File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:     cmd.JSON,
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
		Name: "LoadView TableObject From Json File FormatElement Is Not Specified",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Json File FormatElement is Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Json Lines File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.JSONL,
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
		Name: "LoadView TableObject From LTSV File",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.LTSV,
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
		Name: "LoadView TableObject From LTSV File Without Null",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.LTSV,
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
		Name: "LoadView TableObject From LTSV File with UTF-8 BOM",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.LTSV,
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
		Name: "LoadView TableObject From LTSV File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject Invalid Object Type",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView TableObject From Json File Arguments Length Error",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
		Name: "LoadView Json Inline Table",
		From: parser.FromClause{
			Tables: []parser.QueryExpression{
				parser.Table{
					Object: parser.TableObject{
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
				Format:    cmd.JSON,
				Delimiter: ',',
				JsonQuery: "{column1, column2}",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeTemporaryTable,
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
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
				Format:    cmd.JSON,
				Delimiter: ',',
				JsonQuery: "{}",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeTemporaryTable,
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
					Object: parser.TableObject{
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
				Format:    cmd.JSON,
				Delimiter: ',',
				JsonQuery: "{}",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeTemporaryTable,
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
					Object: parser.TableObject{
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
					Object: parser.TableObject{
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
				Format:    cmd.CSV,
				Delimiter: ',',
				JsonQuery: "",
				Encoding:  text.UTF8,
				LineBreak: text.LF,
				ViewType:  ViewTypeTemporaryTable,
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
		Error: fmt.Sprintf("data parse error in file %s: line 3, column 7: wrong number of fields in line", GetTestFilePath("table_broken.csv")),
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

func TestView_Load(t *testing.T) {
	defer func() {
		_ = TestTx.ReleaseResources()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)
		_ = TestTx.Session.SetStdin(os.Stdin)
		initFlag(TestTx.Flags)
	}()

	TestTx.Flags.Repository = TestDir
	ctx := context.Background()

	for _, v := range viewLoadTests {
		TestTx.UnlockStdin()
		_ = TestTx.cachedViews.Clean(TestTx.FileContainer)

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
		for i := range queryScope.blocks {
			if !reflect.DeepEqual(queryScope.blocks[i].temporaryTables.Keys(), v.ResultScope.blocks[i].temporaryTables.Keys()) {
				t.Errorf("%s: temp view list = %v, want %v", v.Name, queryScope.blocks[i].temporaryTables.Keys(), v.ResultScope.blocks[i].temporaryTables.Keys())
			}
		}

		if !reflect.DeepEqual(view, v.Result) {
			t.Errorf("%s: \n result = %v,\n expect = %v", v.Name, view, v.Result)
		}
	}
}

func TestNewViewFromGroupedRecord(t *testing.T) {
	fr := ReferenceRecord{
		view: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2"), value.NewString("str3")}),
				},
			},
		},
		recordIndex: 0,
		cache:       NewFieldIndexCache(10, LimitToUseFieldIndexSliceChache),
	}
	expect := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			{NewCell(value.NewInteger(1)), NewCell(value.NewInteger(1)), NewCell(value.NewString("str1"))},
			{NewCell(value.NewInteger(2)), NewCell(value.NewInteger(2)), NewCell(value.NewString("str2"))},
			{NewCell(value.NewInteger(3)), NewCell(value.NewInteger(3)), NewCell(value.NewString("str3"))},
		},
	}

	result, _ := NewViewFromGroupedRecord(context.Background(), TestTx.Flags, fr)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}
}

var viewWhereTests = []struct {
	Name   string
	CPU    int
	View   *View
	Where  parser.WhereClause
	Result RecordSet
	Error  string
}{
	{
		Name: "Where",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: parser.Token{Token: '=', Literal: "="},
			},
		},
		Result: RecordSet{
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	},
	{
		Name: "Where in Multi Threading",
		CPU:  3,
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: parser.Token{Token: '=', Literal: "="},
			},
		},
		Result: RecordSet{
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	},
	{
		Name: "Where Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Where: parser.WhereClause{
			Filter: parser.Comparison{
				LHS:      parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
				RHS:      parser.NewIntegerValueFromString("2"),
				Operator: parser.Token{Token: '=', Literal: "="},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_Where(t *testing.T) {
	defer initFlag(TestTx.Flags)

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewWhereTests {
		TestTx.Flags.CPU = 1
		if v.CPU != 0 {
			TestTx.Flags.CPU = v.CPU
		}

		err := v.View.Where(ctx, scope, v.Where)
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
		if !reflect.DeepEqual(v.View.RecordSet, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.RecordSet, v.Result)
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
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewInteger(2), value.NewInteger(4)}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			isGrouped: true,
		},
	},
	{
		Name: "Group By With ColumnNumber",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(3)},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(3)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewInteger(2), value.NewInteger(4)}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
			isGrouped: true,
		},
	},
	{
		Name: "Group By Evaluation Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
					value.NewString("group1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
					value.NewString("group1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.ColumnNumber{View: parser.Identifier{Literal: "table1"}, Number: value.NewInteger(0)},
			},
		},
		Error: "field table1.0 does not exist",
	},
	{
		Name: "Group By Empty Record",
		View: &View{
			Header:    NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{},
		},
		GroupBy: parser.GroupByClause{
			Items: []parser.QueryExpression{
				parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			RecordSet: []Record{},
			isGrouped: true,
		},
	},
	{
		Name: "Group By Empty Record with No Condition",
		View: &View{
			Header:    NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{},
		},
		GroupBy: parser.GroupByClause{
			Items: nil,
		},
		Result: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					Number:      1,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					Number:      2,
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					Number:      3,
					IsFromTable: true,
				},
			},
			RecordSet: []Record{},
			isGrouped: true,
		},
	},
}

func TestView_GroupBy(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewGroupByTests {
		err := v.View.GroupBy(ctx, scope, v.GroupBy)
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
			t.Errorf("%s: result = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewHavingTests = []struct {
	Name   string
	View   *View
	Having parser.HavingClause
	Result RecordSet
	Error  string
}{
	{
		Name: "Having",
		View: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			isGrouped: true,
			RecordSet: RecordSet{
				{
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Result: RecordSet{
			{
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
				NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
			},
		},
	},
	{
		Name: "Having Filter Error",
		View: &View{
			Header: []HeaderField{
				{
					View:   "table1",
					Column: InternalIdColumn,
				},
				{
					View:        "table1",
					Column:      "column1",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column2",
					IsFromTable: true,
				},
				{
					View:        "table1",
					Column:      "column3",
					IsFromTable: true,
					IsGroupKey:  true,
				},
			},
			isGrouped: true,
			RecordSet: RecordSet{
				{
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("3")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str3")}),
					NewGroupCell([]value.Primary{value.NewString("group1"), value.NewString("group1")}),
				},
				{
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
					NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
					NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
				},
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "Having Not Grouped",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Result: RecordSet{
			{
				NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
				NewGroupCell([]value.Primary{value.NewString("2"), value.NewString("4")}),
				NewGroupCell([]value.Primary{value.NewString("str2"), value.NewString("str4")}),
				NewGroupCell([]value.Primary{value.NewString("group2"), value.NewString("group2")}),
			},
		},
	},
	{
		Name: "Having All RecordSet Filter Error",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: RecordSet{
				NewRecordWithId(1, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
					value.NewString("group2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("4"),
					value.NewString("str4"),
					value.NewString("group2"),
				}),
			},
		},
		Having: parser.HavingClause{
			Filter: parser.Comparison{
				LHS: parser.AggregateFunction{
					Name:     "sum",
					Distinct: parser.Token{},
					Args: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
				RHS:      parser.NewIntegerValueFromString("5"),
				Operator: parser.Token{Token: parser.COMPARISON_OP, Literal: ">"},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_Having(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewHavingTests {
		err := v.View.Having(ctx, scope, v.Having)
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
		if !reflect.DeepEqual(v.View.RecordSet, v.Result) {
			t.Errorf("%s: result = %s, want %s", v.Name, v.View.RecordSet, v.Result)
		}
	}
}

var viewSelectTests = []struct {
	Name   string
	View   *View
	Scope  *ReferenceScope
	Select parser.SelectClause
	Result *View
	Error  string
}{
	{
		Name: "Select",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2"}},
				parser.Field{Object: parser.AllColumns{}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2a"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21a"}},
				parser.Field{Object: parser.PrimitiveType{
					Literal: "2012-01-01",
					Value:   value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Aliases: []string{"c2", "c2a"}, Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Aliases: []string{"t21", "t21a"}, Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
				{Column: "1", Aliases: []string{"a"}},
				{Column: "2012-01-01T00:00:00Z"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
			},
			selectFields: []int{2, 1, 2, 4, 5, 6, 2, 4, 4, 7},
		},
	},
	{
		Name: "Select using Table Wildcard",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2"}},
				parser.Field{Object: parser.FieldReference{View: parser.Identifier{Literal: "table2"}, Column: parser.AllColumns{}}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}, Alias: parser.Identifier{Literal: "c2a"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21"}},
				parser.Field{Object: parser.ColumnNumber{View: parser.Identifier{Literal: "table2"}, Number: value.NewInteger(1)}, Alias: parser.Identifier{Literal: "t21a"}},
				parser.Field{Object: parser.PrimitiveType{
					Literal: "2012-01-01",
					Value:   value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Aliases: []string{"c2", "c2a"}, Number: 2, IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", Aliases: []string{"t21", "t21a"}, Number: 1, IsFromTable: true},
				{View: "table2", Column: "column4", Number: 2, IsFromTable: true},
				{Column: "1", Aliases: []string{"a"}},
				{Column: "2012-01-01T00:00:00Z"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("1"),
					value.NewString("str44"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
					value.NewInteger(1),
					value.NewDatetime(time.Date(2012, 1, 1, 0, 0, 0, 0, GetTestLocation())),
				}),
			},
			selectFields: []int{2, 4, 5, 6, 2, 4, 4, 7},
		},
	},
	{
		Name: "Select Distinct",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table2", Column: InternalIdColumn},
				{View: "table2", Column: "column3", IsFromTable: true},
				{View: "table2", Column: "column4", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(2),
					value.NewString("3"),
					value.NewString("str33"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
					value.NewInteger(3),
					value.NewString("4"),
					value.NewString("str44"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
					value.NewInteger(1),
					value.NewString("2"),
					value.NewString("str22"),
				}),
			},
		},
		Select: parser.SelectClause{
			Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.NewIntegerValueFromString("1"), Alias: parser.Identifier{Literal: "a"}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{Column: "1", Aliases: []string{"a"}},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewInteger(1),
				}),
			},
			selectFields: []int{0, 1},
		},
	},
	{
		Name: "Select Aggregate Function",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{
					Object: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Column: "SUM(column1)"},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("2")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2")}),
					NewCell(value.NewFloat(3)),
				},
			},
			selectFields: []int{3},
		},
	},
	{
		Name: "Select Aggregate Function Not Group Key Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Error: "field column2 is not a group key",
	},
	{
		Name: "Select Aggregate Function All RecordSet Lazy Evaluation",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.NewIntegerValueFromString("1")},
				parser.Field{
					Object: parser.Arithmetic{
						LHS: parser.AggregateFunction{
							Name:     "sum",
							Distinct: parser.Token{},
							Args: []parser.QueryExpression{
								parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
							},
						},
						RHS:      parser.NewIntegerValueFromString("1"),
						Operator: parser.Token{Token: '+', Literal: "+"},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Column: "1"},
				{Column: "SUM(column1) + 1"},
			},
			RecordSet: []Record{
				{
					NewGroupCell([]value.Primary{value.NewInteger(1), value.NewInteger(2)}),
					NewGroupCell([]value.Primary{value.NewString("1"), value.NewString("2")}),
					NewGroupCell([]value.Primary{value.NewString("str1"), value.NewString("str2")}),
					NewCell(value.NewInteger(1)),
					NewCell(value.NewFloat(4)),
				},
			},
			selectFields: []int{3, 4},
		},
	},
	{
		Name: "Select Analytic Function",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
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
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{Column: "ROW_NUMBER() OVER (PARTITION BY column1 ORDER BY column2)", Aliases: []string{"rownum"}},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
					value.NewInteger(3),
				}),
			},
			selectFields: []int{0, 1, 2},
		},
	},
	{
		Name: "Select Analytic Function Not Exist Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "notexist",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
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
		Error: "function notexist does not exist",
	},
	{
		Name: "Select Analytic Function Partition Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Select Analytic Function Order Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "row_number",
						AnalyticClause: parser.AnalyticClause{
							PartitionClause: parser.PartitionClause{
								Values: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
								},
							},
							OrderByClause: parser.OrderByClause{
								Items: []parser.QueryExpression{
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
		Error: "field notexist does not exist",
	},
	{
		Name: "Select User Defined Analytic Function",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERAGGFUNC": &UserDefinedFunction{
						Name:         parser.Identifier{Literal: "useraggfunc"},
						IsAggregate:  true,
						Cursor:       parser.Identifier{Literal: "list"},
						RequiredArgs: 0,
						Statements: []parser.Statement{
							parser.VariableDeclaration{
								Assignments: []parser.VariableAssignment{
									{
										Variable: parser.Variable{Name: "value"},
									},
									{
										Variable: parser.Variable{Name: "fetch"},
									},
								},
							},
							parser.WhileInCursor{
								Variables: []parser.Variable{
									{Name: "fetch"},
								},
								Cursor: parser.Identifier{Literal: "list"},
								Statements: []parser.Statement{
									parser.If{
										Condition: parser.Is{
											LHS: parser.Variable{Name: "value"},
											RHS: parser.NewNullValue(),
										},
										Statements: []parser.Statement{
											parser.VariableSubstitution{
												Variable: parser.Variable{Name: "value"},
												Value:    parser.Variable{Name: "fetch"},
											},
											parser.FlowControl{Token: parser.CONTINUE},
										},
									},
									parser.VariableSubstitution{
										Variable: parser.Variable{Name: "value"},
										Value: parser.Arithmetic{
											LHS:      parser.Variable{Name: "value"},
											RHS:      parser.Variable{Name: "fetch"},
											Operator: parser.Token{Token: '+', Literal: "+"},
										},
									},
								},
							},

							parser.Return{
								Value: parser.Variable{Name: "value"},
							},
						},
					},
				},
			},
		}, nil, time.Time{}, nil),
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}},
				parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
				parser.Field{
					Object: parser.AnalyticFunction{
						Name: "useraggfunc",
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
						},
						AnalyticClause: parser.AnalyticClause{},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
				{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
				{Column: "USERAGGFUNC(column2) OVER ()"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(2),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(3),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(5),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("a"),
					value.NewInteger(1),
					value.NewInteger(15),
				}),
				NewRecord([]value.Primary{
					value.NewString("b"),
					value.NewInteger(4),
					value.NewInteger(15),
				}),
			},
			selectFields: []int{0, 1, 2},
		},
	},
	{
		Name: "Select Aggregate Empty Rows",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{Object: parser.AggregateFunction{Name: "count", Args: []parser.QueryExpression{parser.AllColumns{}}}},
				parser.Field{Object: parser.AggregateFunction{Name: "sum", Args: []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}}}},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Column: "COUNT(*)"},
				{Column: "SUM(column1)"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewNull(),
					value.NewInteger(0),
					value.NewNull(),
				}),
			},
			selectFields: []int{2, 3},
		},
	},
	{
		Name: "Select Compound Function with Aggregate Empty Rows",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
		Select: parser.SelectClause{
			Fields: []parser.QueryExpression{
				parser.Field{
					Object: parser.Function{Name: "coalesce",
						Args: []parser.QueryExpression{
							parser.AggregateFunction{Name: "sum", Args: []parser.QueryExpression{parser.FieldReference{Column: parser.Identifier{Literal: "column1"}}}},
							parser.NewIntegerValue(0),
						},
					},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{Column: "COALESCE(SUM(column1), 0)"},
			},
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewNull(),
					value.NewInteger(0),
				}),
			},
			selectFields: []int{2},
		},
	},
}

func TestView_Select(t *testing.T) {
	ctx := context.Background()
	for _, v := range viewSelectTests {
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}
		err := v.View.Select(ctx, v.Scope, v.Select)
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
			t.Errorf("%s: header = %v, want %v", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result.RecordSet) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.RecordSet, v.Result.RecordSet)
		}
		if !reflect.DeepEqual(v.View.selectFields, v.Result.selectFields) {
			t.Errorf("%s: select indices = %v, want %v", v.Name, v.View.selectFields, v.Result.selectFields)
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
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table1", Column: "column3", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
				},
				parser.OrderItem{
					Value:     parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					Direction: parser.Token{Token: parser.DESC, Literal: "desc"},
				},
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column3"}},
				},
				parser.OrderItem{
					Value: parser.NewIntegerValueFromString("1"),
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
				{View: "table1", Column: "column3", IsFromTable: true},
				{Column: "1"},
			},
			RecordSet: []Record{
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
					value.NewInteger(1),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
					value.NewInteger(1),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
					value.NewInteger(1),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
					value.NewInteger(1),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
					value.NewInteger(1),
				}),
			},
		},
	},
	{
		Name: "Order By with Cached SortValues",
		View: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("2"), TestTx.Flags), nil},
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2", "column3"}),
			RecordSet: []Record{
				NewRecordWithId(5, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
					value.NewString("4"),
				}),
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewString("2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewString("3"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("4"),
					value.NewString("3"),
				}),
			},
			sortValuesInEachCell: [][]*SortValue{
				{nil, nil, NewSortValue(value.NewString("2"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("3"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
				{nil, nil, NewSortValue(value.NewString("4"), TestTx.Flags), nil},
			},
		},
	},
	{
		Name: "Order By With Null Positions",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value:         parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
					NullsPosition: parser.Token{Token: parser.LAST, Literal: "last"},
				},
				parser.OrderItem{
					Value:         parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					NullsPosition: parser.Token{Token: parser.FIRST, Literal: "first"},
				},
			},
		},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(5, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewNull(),
					value.NewString("2"),
				}),
			},
		},
	},
	{
		Name: "Order By Record Extend Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewNull(),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("2"),
				}),
			},
		},
		OrderBy: parser.OrderByClause{
			Items: []parser.QueryExpression{
				parser.OrderItem{
					Value: parser.AggregateFunction{
						Name:     "sum",
						Distinct: parser.Token{},
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
		},
		Error: "function sum cannot aggregate not grouping records",
	},
}

func TestView_OrderBy(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewOrderByTests {
		err := v.View.OrderBy(ctx, scope, v.OrderBy)
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
			t.Errorf("%s: header = %v, want %v", v.Name, v.View.Header, v.Result.Header)
		}
		if !reflect.DeepEqual(v.View.RecordSet, v.Result.RecordSet) {
			t.Errorf("%s: records = %s, want %s", v.Name, v.View.RecordSet, v.Result.RecordSet)
		}
	}
}

var viewExtendRecordCapacity = []struct {
	Name   string
	View   *View
	Scope  *ReferenceScope
	Exprs  []parser.QueryExpression
	Result int
	Error  string
}{
	{
		Name: "ExtendRecordCapacity",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
			isGrouped: true,
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "arg1"}},
						},
						IsAggregate: true,
					},
				},
			},
		}, nil, time.Time{}, nil),
		Exprs: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.Function{
				Name: "userfunc",
				Args: []parser.QueryExpression{
					parser.NewIntegerValueFromString("1"),
				},
			},
			parser.AggregateFunction{
				Name:     "avg",
				Distinct: parser.Token{},
				Args: []parser.QueryExpression{
					parser.AggregateFunction{
						Name: "avg",
						Args: []parser.QueryExpression{
							parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
						},
					},
				},
			},
			parser.ListFunction{
				Name:     "listagg",
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.QueryExpression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					parser.NewStringValue(","),
				},
				OrderBy: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
			},
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.Arithmetic{
								LHS:      parser.NewIntegerValueFromString("1"),
								RHS:      parser.NewIntegerValueFromString("2"),
								Operator: parser.Token{Token: '+', Literal: "+"},
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.Arithmetic{
									LHS:      parser.NewIntegerValueFromString("3"),
									RHS:      parser.NewIntegerValueFromString("4"),
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
							},
						},
					},
				},
			},
			parser.Arithmetic{
				LHS:      parser.NewIntegerValueFromString("5"),
				RHS:      parser.NewIntegerValueFromString("6"),
				Operator: parser.Token{Token: '+', Literal: "+"},
			},
		},
		Result: 9,
	},
	{
		Name: "ExtendRecordCapacity UserDefinedFunction Not Grouped Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Scope: GenerateReferenceScope([]map[string]map[string]interface{}{
			{
				scopeNameFunctions: {
					"USERFUNC": &UserDefinedFunction{
						Name: parser.Identifier{Literal: "userfunc"},
						Parameters: []parser.Variable{
							{Name: "arg1"},
						},
						RequiredArgs: 1,
						Statements: []parser.Statement{
							parser.Return{Value: parser.Variable{Name: "arg1"}},
						},
						IsAggregate: true,
					},
				},
			},
		}, nil, time.Time{}, nil),
		Exprs: []parser.QueryExpression{
			parser.Function{
				Name: "userfunc",
				Args: []parser.QueryExpression{
					parser.NewIntegerValueFromString("1"),
				},
			},
		},
		Error: "function userfunc cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity AggregateFunction Not Grouped Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.AggregateFunction{
				Name: "avg",
				Args: []parser.QueryExpression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
				},
			},
		},
		Error: "function avg cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity ListAgg Not Grouped Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.ListFunction{
				Name:     "listagg",
				Distinct: parser.Token{Token: parser.DISTINCT, Literal: "distinct"},
				Args: []parser.QueryExpression{
					parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
					parser.NewStringValue(","),
				},
				OrderBy: parser.OrderByClause{
					Items: []parser.QueryExpression{
						parser.OrderItem{Value: parser.FieldReference{Column: parser.Identifier{Literal: "column2"}}},
					},
				},
			},
		},
		Error: "function listagg cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity AnalyticFunction Partition Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.AggregateFunction{
								Name: "avg",
								Args: []parser.QueryExpression{
									parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
								},
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.Arithmetic{
									LHS:      parser.NewIntegerValueFromString("3"),
									RHS:      parser.NewIntegerValueFromString("4"),
									Operator: parser.Token{Token: '+', Literal: "+"},
								},
							},
						},
					},
				},
			},
		},
		Error: "function avg cannot aggregate not grouping records",
	},
	{
		Name: "ExtendRecordCapacity AnalyticFunction OrderBy Value Error",
		View: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: RecordSet{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewInteger(2),
				}),
			},
		},
		Exprs: []parser.QueryExpression{
			parser.AnalyticFunction{
				Name: "rank",
				AnalyticClause: parser.AnalyticClause{
					PartitionClause: parser.PartitionClause{
						Values: []parser.QueryExpression{
							parser.Arithmetic{
								LHS:      parser.NewIntegerValueFromString("1"),
								RHS:      parser.NewIntegerValueFromString("2"),
								Operator: parser.Token{Token: '+', Literal: "+"},
							},
						},
					},
					OrderByClause: parser.OrderByClause{
						Items: []parser.QueryExpression{
							parser.OrderItem{
								Value: parser.AggregateFunction{
									Name: "avg",
									Args: []parser.QueryExpression{
										parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
									},
								},
							},
						},
					},
				},
			},
		},
		Error: "function avg cannot aggregate not grouping records",
	},
}

func TestView_ExtendRecordCapacity(t *testing.T) {
	ctx := context.Background()
	for _, v := range viewExtendRecordCapacity {
		if v.Scope == nil {
			v.Scope = NewReferenceScope(TestTx)
		}

		err := v.View.ExtendRecordCapacity(ctx, v.Scope, v.Exprs)
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
		if cap(v.View.RecordSet[0]) != v.Result {
			t.Errorf("%s: record capacity = %d, want %d", v.Name, cap(v.View.RecordSet[0]), v.Result)
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
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("2")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Limit With Ties",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			sortValuesInEachRecord: []SortValues{
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 2},
					&SortValue{Type: StringType, String: "str2"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 3},
					&SortValue{Type: StringType, String: "str3"},
				},
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("2"), Restriction: parser.Token{Token: parser.TIES}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			sortValuesInEachRecord: []SortValues{
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 1},
					&SortValue{Type: StringType, String: "str1"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 2},
					&SortValue{Type: StringType, String: "str2"},
				},
				{
					&SortValue{Type: IntegerType, Integer: 3},
					&SortValue{Type: StringType, String: "str3"},
				},
			},
		},
	},
	{
		Name: "Limit By Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
			offset: 1,
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(50.5), Unit: parser.Token{Token: parser.PERCENT}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			offset: 1,
		},
	},
	{
		Name: "Limit By Over 100 Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(150), Unit: parser.Token{Token: parser.PERCENT}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
	},
	{
		Name: "Limit By Negative Percentage",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("3"),
					value.NewString("str3"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewFloatValue(-10), Unit: parser.Token{Token: parser.PERCENT}},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
	},
	{
		Name: "Limit Greater Than RecordSet",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValueFromString("5")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
	},
	{
		Name: "Limit Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.Variable{Name: "notexist"}},
		Error: "variable @notexist is undeclared",
	},
	{
		Name: "Limit Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewStringValue("str")},
		Error: "limit number of records 'str' is not an integer value",
	},
	{
		Name: "Limit Negative Value",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewIntegerValue(-1)},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
		},
	},
	{
		Name: "Limit By Percentage Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Limit: parser.LimitClause{Value: parser.NewStringValue("str"), Unit: parser.Token{Token: parser.PERCENT}},
		Error: "limit percentage 'str' is not a float value",
	},
}

func TestView_Limit(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewLimitTests {
		err := v.View.Limit(ctx, scope, v.Limit)
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
			t.Errorf("%s: view = %v, want %v", v.Name, v.View, v.Result)
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
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValueFromString("3")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
			offset: 3,
		},
	},
	{
		Name: "Offset Equal To Record Length",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValueFromString("4")},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{},
			offset:    4,
		},
	},
	{
		Name: "Offset Evaluate Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.Variable{Name: "notexist"}},
		Error:  "variable @notexist is undeclared",
	},
	{
		Name: "Offset Value Error",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(3, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(4, []value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewStringValue("str")},
		Error:  "offset number 'str' is not an integer value",
	},
	{
		Name: "Offset Negative Number",
		View: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
		},
		Offset: parser.OffsetClause{Value: parser.NewIntegerValue(-3)},
		Result: &View{
			Header: []HeaderField{
				{View: "table1", Column: InternalIdColumn},
				{View: "table1", Column: "column1", IsFromTable: true},
				{View: "table1", Column: "column2", IsFromTable: true},
			},
			RecordSet: []Record{
				NewRecordWithId(1, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecordWithId(2, []value.Primary{
					value.NewString("1"),
					value.NewString("str1"),
				}),
			},
			offset: 0,
		},
	},
}

func TestView_Offset(t *testing.T) {
	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewOffsetTests {
		err := v.View.Offset(ctx, scope, v.Offset)
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
			t.Errorf("%s: view = %v, want %v", v.Name, v.View, v.Result)
		}
	}
}

var viewInsertValuesTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	ValuesList  []parser.QueryExpression
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "InsertValues",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("4"),
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(3),
					value.NewNull(),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(4),
					value.NewNull(),
				}),
			},
		},
		UpdateCount: 2,
	},
	{
		Name: "InsertValues Field Length Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "InsertValues Value Evaluation Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "InsertValues Field Does Not Exist Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_InsertValues(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewInsertValuesTests {
		cnt, err := view.InsertValues(ctx, scope, v.Fields, v.ValuesList)
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
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

	}
}

var viewInsertFromQueryTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	Query       parser.SelectQuery
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "InsertFromQuery",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewInteger(1),
					value.NewString("1"),
					value.NewString("str1"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(2),
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewNull(),
					value.NewInteger(3),
					value.NewNull(),
				}),
			},
		},
		UpdateCount: 1,
	},
	{
		Name: "InsertFromQuery Field Lenght Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Error: "select query should return exactly 2 fields",
	},
	{
		Name: "InsertFromQuery Exuecution Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
}

func TestView_InsertFromQuery(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewInsertFromQueryTests {
		cnt, err := view.InsertFromQuery(ctx, scope, v.Fields, v.Query)
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
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}
	}
}

var viewReplaceValuesTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	Keys        []parser.QueryExpression
	ValuesList  []parser.QueryExpression
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "ReplaceValues",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("1"),
						parser.NewStringValue("str3"),
					},
				},
			},
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("4"),
						parser.NewStringValue("str4"),
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
				NewRecord([]value.Primary{
					value.NewInteger(4),
					value.NewString("str4"),
				}),
			},
		},
		UpdateCount: 2,
	},
	{
		Name: "ReplaceValues Field Length Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "row value should contain exactly 2 values",
	},
	{
		Name: "ReplaceValues Value Evaluation Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ReplaceValues Field Does Not Exist Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ReplaceValues Key Does Not Exist Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "field notexist does not exist",
	},
	{
		Name: "ReplaceValues Key Not Set Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		ValuesList: []parser.QueryExpression{
			parser.RowValue{
				Value: parser.ValueList{
					Values: []parser.QueryExpression{
						parser.NewIntegerValueFromString("3"),
					},
				},
			},
		},
		Error: "replace Key column2 is not set",
	},
}

func TestView_ReplaceValues(t *testing.T) {
	view := &View{
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
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewReplaceValuesTests {
		cnt, err := view.ReplaceValues(ctx, scope, v.Fields, v.ValuesList, v.Keys)
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
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}

	}
}

var viewReplaceFromQueryTests = []struct {
	Name        string
	Fields      []parser.QueryExpression
	Keys        []parser.QueryExpression
	Query       parser.SelectQuery
	Result      *View
	UpdateCount int
	Error       string
}{
	{
		Name: "ReplaceFromQuery",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Keys: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("1")},
						parser.Field{Object: parser.NewStringValue("str3")},
					},
				},
			},
		},
		Result: &View{
			Header: NewHeader("table1", []string{"column1", "column2"}),
			RecordSet: []Record{
				NewRecord([]value.Primary{
					value.NewString("1"),
					value.NewString("str3"),
				}),
				NewRecord([]value.Primary{
					value.NewString("2"),
					value.NewString("str2"),
				}),
			},
		},
		UpdateCount: 1,
	},
	{
		Name: "ReplaceFromQuery Field Lenght Does Not Match Error",
		Fields: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		},
		Query: parser.SelectQuery{
			SelectEntity: parser.SelectEntity{
				SelectClause: parser.SelectClause{
					Fields: []parser.QueryExpression{
						parser.Field{Object: parser.NewIntegerValueFromString("3")},
					},
				},
			},
		},
		Error: "select query should return exactly 2 fields",
	},
}

func TestView_ReplaceFromQuery(t *testing.T) {
	view := &View{
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
	}

	scope := NewReferenceScope(TestTx)
	ctx := context.Background()
	for _, v := range viewReplaceFromQueryTests {
		cnt, err := view.ReplaceFromQuery(ctx, scope, v.Fields, v.Query, v.Keys)
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
		if cnt != v.UpdateCount {
			t.Errorf("%s: update count = %d, want %d", v.Name, cnt, v.UpdateCount)
		}
	}
}

func TestView_Fix(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: InternalIdColumn},
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int{2},
	}
	expect := &View{
		Header: NewHeader("table1", []string{"column2"}),
		RecordSet: []Record{
			NewRecord([]value.Primary{
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("str1"),
			}),
		},
		selectFields: []int(nil),
	}

	_ = view.Fix(context.Background(), TestTx.Flags)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("fix: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: InternalIdColumn},
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: []Record{
			NewRecordWithId(1, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecordWithId(2, []value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int{2, 2, 2, 2},
	}
	expect = &View{
		Header: NewHeader("table1", []string{"column2", "column2", "column2", "column2"}),
		RecordSet: []Record{
			NewRecord([]value.Primary{
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
				value.NewString("str1"),
			}),
		},
		selectFields: []int(nil),
	}

	_ = view.Fix(context.Background(), TestTx.Flags)
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("fix: view = %v, want %v", view, expect)
	}
}

func TestView_Union(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
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
	}

	ctx := context.Background()
	err := view.Union(ctx, TestTx.Flags, calcView, false)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
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
	}

	err = view.Union(ctx, TestTx.Flags, calcView, true)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("union all: view = %v, want %v", view, expect)
	}
}

func TestView_Except(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
	}

	ctx := context.Background()
	err := view.Except(ctx, TestTx.Flags, calcView, false)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
			}),
		},
	}

	err = view.Except(ctx, TestTx.Flags, calcView, true)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("except all: view = %v, want %v", view, expect)
	}
}

func TestView_Intersect(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	calcView := &View{
		Header: []HeaderField{
			{View: "table2", Column: "column3", IsFromTable: true},
			{View: "table2", Column: "column4", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("3"),
				value.NewString("str3"),
			}),
		},
	}

	expect := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	ctx := context.Background()
	err := view.Intersect(ctx, TestTx.Flags, calcView, false)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect: view = %v, want %v", view, expect)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("1"),
				value.NewString("str1"),
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
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	expect = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
		RecordSet: RecordSet{
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
			NewRecord([]value.Primary{
				value.NewString("2"),
				value.NewString("str2"),
			}),
		},
	}

	err = view.Intersect(ctx, TestTx.Flags, calcView, true)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if !reflect.DeepEqual(view, expect) {
		t.Errorf("intersect all: view = %v, want %v", view, expect)
	}
}

func TestView_FieldIndex(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", Number: 1, IsFromTable: true},
			{View: "table1", Column: "column2", Number: 2, IsFromTable: true},
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

	columnNum := parser.ColumnNumber{
		View:   parser.Identifier{Literal: "table1"},
		Number: value.NewInteger(2),
	}
	expect = 1

	idx, _ = view.FieldIndex(columnNum)
	if idx != expect {
		t.Errorf("field index = %d, want %d", idx, expect)
	}
}

func TestView_FieldIndices(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table1", Column: "column2", IsFromTable: true},
		},
	}
	fields := []parser.QueryExpression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "column1"}},
	}
	expect := []int{1, 0}

	indices, _ := view.FieldIndices(fields)
	if !reflect.DeepEqual(indices, expect) {
		t.Errorf("field indices = %v, want %v", indices, expect)
	}

	fields = []parser.QueryExpression{
		parser.FieldReference{Column: parser.Identifier{Literal: "column2"}},
		parser.FieldReference{Column: parser.Identifier{Literal: "notexist"}},
	}
	expectErr := "field notexist does not exist"
	_, err := view.FieldIndices(fields)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func TestView_FieldViewName(t *testing.T) {
	view := &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table2", Column: "column2", IsFromTable: true},
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

	fieldRef = parser.FieldReference{
		Column: parser.Identifier{Literal: "notexist"},
	}
	expectErr := "field notexist does not exist"
	_, err := view.FieldViewName(fieldRef)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func TestView_InternalRecordId(t *testing.T) {
	view := &View{
		Header: NewHeaderWithId("table1", []string{"column1", "column2"}),
		RecordSet: []Record{
			NewRecordWithId(0, []value.Primary{value.NewInteger(1), value.NewString("str1")}),
			NewRecordWithId(1, []value.Primary{value.NewInteger(2), value.NewString("str2")}),
			NewRecordWithId(2, []value.Primary{value.NewInteger(3), value.NewString("str3")}),
		},
	}
	ref := "table1"
	recordIndex := 1
	expect := 1

	id, _ := view.InternalRecordId(ref, recordIndex)
	if id != expect {
		t.Errorf("field internal id = %d, want %d", id, expect)
	}

	view.RecordSet[1][0] = NewCell(value.NewNull())
	expectErr := "internal record id is empty"
	_, err := view.InternalRecordId(ref, recordIndex)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}

	view = &View{
		Header: []HeaderField{
			{View: "table1", Column: "column1", IsFromTable: true},
			{View: "table2", Column: "column2", IsFromTable: true},
		},
	}
	expectErr = "internal record id does not exist"
	_, err = view.InternalRecordId(ref, recordIndex)
	if err == nil {
		t.Errorf("no error, want error %s", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error = %s, want %s", err, expectErr)
	}
}

func BenchmarkView_GroupBy(b *testing.B) {
	view := &View{
		Header:    NewHeader("t", []string{"c1", "c2", "c3"}),
		RecordSet: make(RecordSet, 10000),
	}
	for i := int64(0); i < 10000; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{
			value.NewInteger(i),
			value.NewString(randomStr(1)),
			value.NewString(randomStr(1)),
		})
	}

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)
	clause := parser.GroupByClause{
		Items: []parser.QueryExpression{
			parser.FieldReference{Column: parser.Identifier{Literal: "c2"}},
			parser.FieldReference{Column: parser.Identifier{Literal: "c3"}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := &View{
			Header:    view.Header.Copy(),
			RecordSet: view.RecordSet.Copy(),
		}
		_ = v.GroupBy(ctx, scope, clause)
	}
}

func BenchmarkView_SelectDistinct(b *testing.B) {
	view := &View{
		Header:    NewHeader("t", []string{"c1", "c2", "c3"}),
		RecordSet: make(RecordSet, 10000),
	}
	for i := int64(0); i < 10000; i++ {
		view.RecordSet[i] = NewRecord([]value.Primary{
			value.NewInteger(i),
			value.NewString(randomStr(1)),
			value.NewString(randomStr(1)),
		})
	}

	ctx := context.Background()
	scope := NewReferenceScope(TestTx)
	clause := parser.SelectClause{
		Distinct: parser.Token{Token: parser.DISTINCT},
		Fields: []parser.QueryExpression{
			parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c1"}}},
			parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c2"}}},
			parser.Field{Object: parser.FieldReference{Column: parser.Identifier{Literal: "c3"}}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := &View{
			Header:    view.Header.Copy(),
			RecordSet: view.RecordSet.Copy(),
		}
		_ = v.Select(ctx, scope, clause)
	}
}
