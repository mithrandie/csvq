package json

import (
	"github.com/mithrandie/csvq/lib/value"
	"reflect"
	"testing"
)

var loadValueTests = []struct {
	Query  string
	Json   string
	Expect value.Primary
	Error  string
}{
	{
		Query:  "key",
		Json:   "{\"key\":\"value\"}",
		Expect: value.NewString("value"),
	},
	{
		Query: "'key",
		Json:  "{\"key\":\"value\"}",
		Error: "column 4: string not terminated",
	},
	{
		Query: "key",
		Json:  "{\"key\":\"valu",
		Error: "line 1, column 12: unexpected termination",
	},
}

func TestLoadValue(t *testing.T) {
	for _, v := range loadValueTests {
		result, err := LoadValue(v.Query, v.Json)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q, %q", err.Error(), v.Query, v.Json)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q, %q", err, v.Error, v.Query, v.Json)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q, %q", v.Error, v.Query, v.Json)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q, %q", result, v.Expect, v.Query, v.Json)
		}
	}
}

var loadRowValueTests = []struct {
	Query  string
	Json   string
	Expect []value.Primary
	Error  string
}{
	{
		Query: "key[]",
		Json:  "{\"key\":[1, 2, 3]}",
		Expect: []value.Primary{
			value.NewInteger(1),
			value.NewInteger(2),
			value.NewInteger(3),
		},
	},
	{
		Query: "notexist[]",
		Json:  "{\"key\":[{\"key2\":2, \"key3\": 3}]}",
		Error: "json value does not exists for \"notexist[]\"",
	},
	{
		Query: "key[]",
		Json:  "{\"key\":\"value\"}",
		Error: "json value must be an array",
	},
}

func TestLoadRowValue(t *testing.T) {
	for _, v := range loadRowValueTests {
		result, err := LoadArray(v.Query, v.Json)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q, %q", err.Error(), v.Query, v.Json)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q, %q", err, v.Error, v.Query, v.Json)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q, %q", v.Error, v.Query, v.Json)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q, %q", result, v.Expect, v.Query, v.Json)
		}
	}
}

var loadTableTests = []struct {
	Query        string
	Json         string
	ExpectHeader []string
	ExpectValues [][]value.Primary
	EscapeType   EscapeType
	Error        string
}{
	{
		Query:        "key{}",
		Json:         "{\"key\":[{\"key2\":2, \"key3\": 3}, {\"key2\":4, \"key3\": 5}]}",
		ExpectHeader: []string{"key2", "key3"},
		ExpectValues: [][]value.Primary{
			{
				value.NewInteger(2),
				value.NewInteger(3),
			},
			{
				value.NewInteger(4),
				value.NewInteger(5),
			},
		},
	},
	{
		Query: "notexist{}",
		Json:  "{\"key\":[{\"key2\":2, \"key3\": 3}]}",
		Error: "json value does not exists for \"notexist{}\"",
	},
	{
		Query: "key{}",
		Json:  "{\"key\":[{\"key2\":2, \"key3\": 3}, {\"key2\":4, key3: 5}]}",
		Error: "line 1, column 43: unexpected token \"key\"",
	},
}

func TestLoadTable(t *testing.T) {
	for _, v := range loadTableTests {
		header, values, et, err := LoadTable(v.Query, v.Json)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q, %q", err.Error(), v.Query, v.Json)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q, %q", err, v.Error, v.Query, v.Json)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q, %q", v.Error, v.Query, v.Json)
			continue
		}
		if !reflect.DeepEqual(header, v.ExpectHeader) {
			t.Errorf("header = %#v, want %#v for %q, %q", header, v.ExpectHeader, v.Query, v.Json)
		}
		if !reflect.DeepEqual(values, v.ExpectValues) {
			t.Errorf("values = %#v, want %#v for %q, %q", values, v.ExpectValues, v.Query, v.Json)
		}
		if et != v.EscapeType {
			t.Errorf("escape type = %d, want %d for %q, %q", et, v.EscapeType, v.Query, v.Json)
		}
	}
}

var extractTests = []struct {
	Query  QueryExpression
	Data   Structure
	Expect Structure
	Error  string
}{
	{
		Query: nil,
		Data: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
				{
					Key: "obj",
					Value: Object{
						Members: []ObjectMember{
							{
								Key:   "key2",
								Value: String("value2"),
							},
						},
					},
				},
			},
		},
		Expect: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
				{
					Key: "obj",
					Value: Object{
						Members: []ObjectMember{
							{
								Key:   "key2",
								Value: String("value2"),
							},
						},
					},
				},
			},
		},
	},
	{
		Query: Element{
			Label: "obj",
			Child: Element{
				Label: "key2",
			},
		},
		Data: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
				{
					Key: "obj",
					Value: Object{
						Members: []ObjectMember{
							{
								Key:   "key2",
								Value: String("value2"),
							},
						},
					},
				},
			},
		},
		Expect: String("value2"),
	},
	{
		Query: Element{
			Label: "notexist",
		},
		Data: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
			},
		},
		Expect: Null{},
	},
	{
		Query: Element{
			Label: "notexist",
		},
		Data:   String("value"),
		Expect: Null{},
	},
	{
		Query: ArrayItem{
			Index: 1,
			Child: Element{
				Label: "key2",
			},
		},
		Data: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value"),
					},
					{
						Key:   "key2",
						Value: String("value2"),
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value10"),
					},
					{
						Key:   "key2",
						Value: String("value12"),
					},
				},
			},
		},
		Expect: String("value12"),
	},
	{
		Query: ArrayItem{
			Index: 0,
		},
		Data: Array{
			String("value1"),
			String("value2"),
		},
		Expect: String("value1"),
	},
	{
		Query: ArrayItem{
			Index: 3,
		},
		Data: Array{
			String("value1"),
			String("value2"),
		},
		Expect: Null{},
	},
	{
		Query: ArrayItem{
			Index: 3,
		},
		Data:   String("value1"),
		Expect: Null{},
	},
	{
		Query: RowValueExpr{},
		Data: Array{
			String("value1"),
			String("value2"),
		},
		Expect: Array{
			String("value1"),
			String("value2"),
		},
	},
	{
		Query: RowValueExpr{
			Child: Element{Label: "key"},
		},
		Data: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value"),
					},
					{
						Key:   "key2",
						Value: String("value2"),
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value10"),
					},
					{
						Key:   "key2",
						Value: String("value12"),
					},
				},
			},
		},
		Expect: Array{
			String("value"),
			String("value10"),
		},
	},
	{
		Query: RowValueExpr{},
		Data:  String("value1"),
		Error: "json value must be an array",
	},
	{
		Query: RowValueExpr{
			Child: RowValueExpr{},
		},
		Data: Array{
			String("value1"),
			String("value2"),
		},
		Error: "json value must be an array",
	},
	{
		Query: TableExpr{},
		Data: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value"),
					},
					{
						Key:   "key2",
						Value: String("value2"),
					},
				},
			},
		},
	},
	{
		Query: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{
						Label: "key2",
						Child: Element{Label: "key3"},
					},
					Alias: "key3",
				},
				{
					Element: Element{Label: "key"},
				},
			},
		},
		Data: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
				{
					Key: "key2",
					Value: Object{
						Members: []ObjectMember{
							{
								Key:   "key3",
								Value: String("value3"),
							},
							{
								Key:   "key4",
								Value: String("value4"),
							},
						},
					},
				},
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key3",
						Value: String("value3"),
					},
					{
						Key:   "key",
						Value: String("value"),
					},
				},
			},
		},
	},
	{
		Query: TableExpr{},
		Data: Array{
			Object{
				Members: []ObjectMember{
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value3"),
								},
								{
									Key:   "key4",
									Value: String("value4"),
								},
							},
						},
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value11"),
					},
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value13"),
								},
								{
									Key:   "key4",
									Value: String("value14"),
								},
							},
						},
					},
					{
						Key:   "key9",
						Value: String("value19"),
					},
				},
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value3"),
								},
								{
									Key:   "key4",
									Value: String("value4"),
								},
							},
						},
					},
					{
						Key:   "key",
						Value: Null{},
					},
					{
						Key:   "key9",
						Value: Null{},
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value13"),
								},
								{
									Key:   "key4",
									Value: String("value14"),
								},
							},
						},
					},
					{
						Key:   "key",
						Value: String("value11"),
					},
					{
						Key:   "key9",
						Value: String("value19"),
					},
				},
			},
		},
	},
	{
		Query: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{
						Label: "key2",
						Child: Element{Label: "key3"},
					},
					Alias: "key3",
				},
				{
					Element: Element{Label: "key"},
				},
				{
					Element: Element{Label: "key5"},
				},
			},
		},
		Data: Array{
			Object{
				Members: []ObjectMember{
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value3"),
								},
								{
									Key:   "key4",
									Value: String("value4"),
								},
							},
						},
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value11"),
					},
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value13"),
								},
								{
									Key:   "key4",
									Value: String("value14"),
								},
							},
						},
					},
					{
						Key:   "key9",
						Value: String("value19"),
					},
				},
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key3",
						Value: String("value3"),
					},
					{
						Key:   "key",
						Value: Null{},
					},
					{
						Key:   "key5",
						Value: Null{},
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key3",
						Value: String("value13"),
					},
					{
						Key:   "key",
						Value: String("value11"),
					},
					{
						Key:   "key5",
						Value: Null{},
					},
				},
			},
		},
	},
	{
		Query: TableExpr{},
		Data: Array{
			String("value1"),
			String("value2"),
		},
		Error: "all elements in array must be objects",
	},
	{
		Query: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{
						Label: "key2",
						Child: RowValueExpr{},
					},
				},
				{
					Element: Element{Label: "key"},
				},
			},
		},
		Data: Array{
			Object{
				Members: []ObjectMember{
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value3"),
								},
								{
									Key:   "key4",
									Value: String("value4"),
								},
							},
						},
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key",
						Value: String("value11"),
					},
					{
						Key: "key2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "key3",
									Value: String("value13"),
								},
								{
									Key:   "key4",
									Value: String("value14"),
								},
							},
						},
					},
					{
						Key:   "key9",
						Value: String("value19"),
					},
				},
			},
		},
		Error: "json value must be an array",
	},
	{
		Query: TableExpr{
			Fields: []FieldExpr{
				{
					Element: Element{
						Label: "key2",
						Child: RowValueExpr{},
					},
				},
				{
					Element: Element{Label: "key"},
				},
			},
		},
		Data: Object{
			Members: []ObjectMember{
				{
					Key: "key2",
					Value: Object{
						Members: []ObjectMember{
							{
								Key:   "key3",
								Value: String("value3"),
							},
							{
								Key:   "key4",
								Value: String("value4"),
							},
						},
					},
				},
			},
		},
		Error: "json value must be an array",
	},
	{
		Query: TableExpr{},
		Data:  String("value1"),
		Error: "json value must be an array or object",
	},
	{
		Query: FieldExpr{},
		Data:  String("value1"),
		Error: "invalid expression",
	},
}

func TestExtract(t *testing.T) {
	for _, v := range extractTests {
		result, err := Extract(v.Query, v.Data)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q, %q", err.Error(), v.Query, v.Data)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q, %q", err, v.Error, v.Query, v.Data)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q, %q", v.Error, v.Query, v.Data)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q, %q", result, v.Expect, v.Query, v.Data)
		}
	}
}
