package json

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text/json"
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
	EscapeType   json.EscapeType
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
	Data   json.Structure
	Expect json.Structure
	Error  string
}{
	{
		Query: nil,
		Data: json.Object{
			Members: []json.ObjectMember{
				{
					Key:   "key",
					Value: json.String("value"),
				},
				{
					Key: "obj",
					Value: json.Object{
						Members: []json.ObjectMember{
							{
								Key:   "key2",
								Value: json.String("value2"),
							},
						},
					},
				},
			},
		},
		Expect: json.Object{
			Members: []json.ObjectMember{
				{
					Key:   "key",
					Value: json.String("value"),
				},
				{
					Key: "obj",
					Value: json.Object{
						Members: []json.ObjectMember{
							{
								Key:   "key2",
								Value: json.String("value2"),
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
		Data: json.Object{
			Members: []json.ObjectMember{
				{
					Key:   "key",
					Value: json.String("value"),
				},
				{
					Key: "obj",
					Value: json.Object{
						Members: []json.ObjectMember{
							{
								Key:   "key2",
								Value: json.String("value2"),
							},
						},
					},
				},
			},
		},
		Expect: json.String("value2"),
	},
	{
		Query: Element{
			Label: "notexist",
		},
		Data: json.Object{
			Members: []json.ObjectMember{
				{
					Key:   "key",
					Value: json.String("value"),
				},
			},
		},
		Expect: json.Null{},
	},
	{
		Query: Element{
			Label: "notexist",
		},
		Data:   json.String("value"),
		Expect: json.Null{},
	},
	{
		Query: ArrayItem{
			Index: 1,
			Child: Element{
				Label: "key2",
			},
		},
		Data: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value10"),
					},
					{
						Key:   "key2",
						Value: json.String("value12"),
					},
				},
			},
		},
		Expect: json.String("value12"),
	},
	{
		Query: ArrayItem{
			Index: 0,
		},
		Data: json.Array{
			json.String("value1"),
			json.String("value2"),
		},
		Expect: json.String("value1"),
	},
	{
		Query: ArrayItem{
			Index: 3,
		},
		Data: json.Array{
			json.String("value1"),
			json.String("value2"),
		},
		Expect: json.Null{},
	},
	{
		Query: ArrayItem{
			Index: 3,
		},
		Data:   json.String("value1"),
		Expect: json.Null{},
	},
	{
		Query: RowValueExpr{},
		Data: json.Array{
			json.String("value1"),
			json.String("value2"),
		},
		Expect: json.Array{
			json.String("value1"),
			json.String("value2"),
		},
	},
	{
		Query: RowValueExpr{
			Child: Element{Label: "key"},
		},
		Data: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value10"),
					},
					{
						Key:   "key2",
						Value: json.String("value12"),
					},
				},
			},
		},
		Expect: json.Array{
			json.String("value"),
			json.String("value10"),
		},
	},
	{
		Query: RowValueExpr{},
		Data:  json.String("value1"),
		Error: "json value must be an array",
	},
	{
		Query: RowValueExpr{
			Child: RowValueExpr{},
		},
		Data: json.Array{
			json.String("value1"),
			json.String("value2"),
		},
		Error: "json value must be an array",
	},
	{
		Query: TableExpr{},
		Data: json.Object{
			Members: []json.ObjectMember{
				{
					Key:   "key",
					Value: json.String("value"),
				},
				{
					Key:   "key2",
					Value: json.String("value2"),
				},
			},
		},
		Expect: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value"),
					},
					{
						Key:   "key2",
						Value: json.String("value2"),
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
		Data: json.Object{
			Members: []json.ObjectMember{
				{
					Key:   "key",
					Value: json.String("value"),
				},
				{
					Key: "key2",
					Value: json.Object{
						Members: []json.ObjectMember{
							{
								Key:   "key3",
								Value: json.String("value3"),
							},
							{
								Key:   "key4",
								Value: json.String("value4"),
							},
						},
					},
				},
			},
		},
		Expect: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key3",
						Value: json.String("value3"),
					},
					{
						Key:   "key",
						Value: json.String("value"),
					},
				},
			},
		},
	},
	{
		Query: TableExpr{},
		Data: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value3"),
								},
								{
									Key:   "key4",
									Value: json.String("value4"),
								},
							},
						},
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value11"),
					},
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value13"),
								},
								{
									Key:   "key4",
									Value: json.String("value14"),
								},
							},
						},
					},
					{
						Key:   "key9",
						Value: json.String("value19"),
					},
				},
			},
		},
		Expect: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value3"),
								},
								{
									Key:   "key4",
									Value: json.String("value4"),
								},
							},
						},
					},
					{
						Key:   "key",
						Value: json.Null{},
					},
					{
						Key:   "key9",
						Value: json.Null{},
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value13"),
								},
								{
									Key:   "key4",
									Value: json.String("value14"),
								},
							},
						},
					},
					{
						Key:   "key",
						Value: json.String("value11"),
					},
					{
						Key:   "key9",
						Value: json.String("value19"),
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
		Data: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value3"),
								},
								{
									Key:   "key4",
									Value: json.String("value4"),
								},
							},
						},
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value11"),
					},
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value13"),
								},
								{
									Key:   "key4",
									Value: json.String("value14"),
								},
							},
						},
					},
					{
						Key:   "key9",
						Value: json.String("value19"),
					},
				},
			},
		},
		Expect: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key3",
						Value: json.String("value3"),
					},
					{
						Key:   "key",
						Value: json.Null{},
					},
					{
						Key:   "key5",
						Value: json.Null{},
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key3",
						Value: json.String("value13"),
					},
					{
						Key:   "key",
						Value: json.String("value11"),
					},
					{
						Key:   "key5",
						Value: json.Null{},
					},
				},
			},
		},
	},
	{
		Query: TableExpr{},
		Data: json.Array{
			json.String("value1"),
			json.String("value2"),
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
		Data: json.Array{
			json.Object{
				Members: []json.ObjectMember{
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value3"),
								},
								{
									Key:   "key4",
									Value: json.String("value4"),
								},
							},
						},
					},
				},
			},
			json.Object{
				Members: []json.ObjectMember{
					{
						Key:   "key",
						Value: json.String("value11"),
					},
					{
						Key: "key2",
						Value: json.Object{
							Members: []json.ObjectMember{
								{
									Key:   "key3",
									Value: json.String("value13"),
								},
								{
									Key:   "key4",
									Value: json.String("value14"),
								},
							},
						},
					},
					{
						Key:   "key9",
						Value: json.String("value19"),
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
		Data: json.Object{
			Members: []json.ObjectMember{
				{
					Key: "key2",
					Value: json.Object{
						Members: []json.ObjectMember{
							{
								Key:   "key3",
								Value: json.String("value3"),
							},
							{
								Key:   "key4",
								Value: json.String("value4"),
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
		Data:  json.String("value1"),
		Error: "json value must be an array or object",
	},
	{
		Query: FieldExpr{},
		Data:  json.String("value1"),
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
