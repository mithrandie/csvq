package json

import (
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"reflect"
	"testing"
)

var convertToValueTests = []struct {
	Input  Structure
	Expect value.Primary
}{
	{
		Input:  Number(-2.34),
		Expect: value.NewFloat(-2.34),
	},
	{
		Input:  Number(234),
		Expect: value.NewInteger(234),
	},
	{
		Input:  String("abc"),
		Expect: value.NewString("abc"),
	},
	{
		Input:  Boolean(false),
		Expect: value.NewBoolean(false),
	},
	{
		Input:  Null{},
		Expect: value.NewNull(),
	},
	{
		Input: Array{
			String("abc"),
			String("def"),
		},
		Expect: value.NewString("[\"abc\",\"def\"]"),
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Expect: value.NewString("{\"key1\":\"value1\",\"key2\":\"value2\"}"),
	},
}

func TestConvertToValue(t *testing.T) {
	for _, v := range convertToValueTests {
		result := ConvertToValue(v.Input)
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %#v", result, v.Expect, v.Input)
		}
	}
}

var convertToArrayTests = []struct {
	Input  Array
	Expect []value.Primary
}{
	{
		Input: Array{
			String("elem1"),
			String("elem2"),
		},
		Expect: []value.Primary{
			value.NewString("elem1"),
			value.NewString("elem2"),
		},
	},
}

func TestConvertToArray(t *testing.T) {
	for _, v := range convertToArrayTests {
		rowValue := ConvertToArray(v.Input)
		if !reflect.DeepEqual(rowValue, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", rowValue, v.Expect, v.Input)
		}
	}
}

var convertToTableValueTests = []struct {
	Input        Array
	ExpectHeader []string
	ExpectRows   [][]value.Primary
	Error        string
}{
	{
		Input: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key1",
						Value: Number(1),
					},
					{
						Key:   "key2",
						Value: Number(2),
					},
					{
						Key:   "key3",
						Value: Number(3),
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "key2",
						Value: Number(22),
					},
					{
						Key:   "key3",
						Value: Number(23),
					},
					{
						Key:   "key2",
						Value: Number(24),
					},
				},
			},
		},
		ExpectHeader: []string{
			"key1",
			"key2",
			"key3",
		},
		ExpectRows: [][]value.Primary{
			{
				value.NewInteger(1),
				value.NewInteger(2),
				value.NewInteger(3),
			},
			{
				value.NewNull(),
				value.NewInteger(22),
				value.NewInteger(23),
			},
		},
	},
	{
		Input: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "key1",
						Value: Number(1),
					},
					{
						Key:   "key2",
						Value: Number(2),
					},
					{
						Key:   "key3",
						Value: Number(3),
					},
				},
			},
			String("abc"),
		},
		Error: "rows loaded from json must be objects",
	},
}

func TestConvertToTableValue(t *testing.T) {
	for _, v := range convertToTableValueTests {
		header, rows, err := ConvertToTableValue(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Input)
			continue
		}
		if !reflect.DeepEqual(header, v.ExpectHeader) {
			t.Errorf("header = %#v, want %#v for %q", header, v.ExpectHeader, v.Input)
		}
		if !reflect.DeepEqual(rows, v.ExpectRows) {
			t.Errorf("rows = %#v, want %#v for %q", rows, v.ExpectRows, v.Input)
		}
	}
}

var convertTableValueToJsonStructureTests = []struct {
	Fields []string
	Rows   [][]value.Primary
	Expect Structure
	Error  string
}{
	{
		Fields: []string{
			"column1",
			"column2",
		},
		Rows: [][]value.Primary{
			{
				value.NewString("a"),
				value.NewInteger(1),
			},
			{
				value.NewString("b"),
				value.NewInteger(2),
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "column1",
						Value: String("a"),
					},
					{
						Key:   "column2",
						Value: Number(1),
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "column1",
						Value: String("b"),
					},
					{
						Key:   "column2",
						Value: Number(2),
					},
				},
			},
		},
	},
	{
		Fields: []string{
			"column1",
			"column2.child1.child11",
			"column2.child2.child22",
		},
		Rows: [][]value.Primary{
			{
				value.NewString("a"),
				value.NewInteger(1),
				value.NewInteger(11),
			},
			{
				value.NewString("b"),
				value.NewInteger(2),
				value.NewInteger(22),
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "column1",
						Value: String("a"),
					},
					{
						Key: "column2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key: "child1",
									Value: Object{
										Members: []ObjectMember{
											{
												Key:   "child11",
												Value: Number(1),
											},
										},
									},
								},
								{
									Key: "child2",
									Value: Object{
										Members: []ObjectMember{
											{
												Key:   "child22",
												Value: Number(11),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Object{
				Members: []ObjectMember{
					{
						Key:   "column1",
						Value: String("b"),
					},
					{
						Key: "column2",
						Value: Object{
							Members: []ObjectMember{
								{
									Key: "child1",
									Value: Object{
										Members: []ObjectMember{
											{
												Key:   "child11",
												Value: Number(2),
											},
										},
									},
								},
								{
									Key: "child2",
									Value: Object{
										Members: []ObjectMember{
											{
												Key:   "child22",
												Value: Number(22),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Fields: []string{
			"string",
			"integer",
			"float",
			"boolean",
			"ternary",
			"ternary2",
			"datetime",
			"null",
		},
		Rows: [][]value.Primary{
			{
				value.NewString("abc"),
				value.NewInteger(1),
				value.NewFloat(1.1),
				value.NewBoolean(false),
				value.NewTernary(ternary.TRUE),
				value.NewTernary(ternary.UNKNOWN),
				value.NewDatetimeFromString("2012-02-02 22:22:22 -07:00"),
				value.NewNull(),
			},
		},
		Expect: Array{
			Object{
				Members: []ObjectMember{
					{
						Key:   "string",
						Value: String("abc"),
					},
					{
						Key:   "integer",
						Value: Number(1),
					},
					{
						Key:   "float",
						Value: Number(1.1),
					},
					{
						Key:   "boolean",
						Value: Boolean(false),
					},
					{
						Key:   "ternary",
						Value: Boolean(true),
					},
					{
						Key:   "ternary2",
						Value: Null{},
					},
					{
						Key:   "datetime",
						Value: String("2012-02-02T22:22:22-07:00"),
					},
					{
						Key:   "null",
						Value: Null{},
					},
				},
			},
		},
	},
	{
		Fields: []string{
			"column1",
			"column2",
		},
		Rows: [][]value.Primary{
			{
				value.NewString("a"),
				value.NewInteger(1),
			},
			{
				value.NewString("b"),
			},
		},
		Error: "field length does not match",
	},
	{
		Fields: []string{
			"column1",
			"column2..",
		},
		Rows: [][]value.Primary{
			{
				value.NewString("a"),
				value.NewInteger(1),
			},
			{
				value.NewString("b"),
				value.NewInteger(2),
			},
		},
		Error: "unexpected token \".\" at column 9 in \"column2..\"",
	},
}

func TestConvertTableValueToJsonStructure(t *testing.T) {
	for _, v := range convertTableValueToJsonStructureTests {
		result, err := ConvertTableValueToJsonStructure(v.Fields, v.Rows)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %s, %s", err.Error(), v.Fields, v.Rows)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %s, %s", err.Error(), v.Error, v.Fields, v.Rows)
			}
			continue
		}
		if 0 < len(v.Error) {
			//fmt.Printf("%#v\n", result)
			t.Errorf("no error, want error %q for %s, %s", v.Error, v.Fields, v.Rows)
			continue
		}
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %s, %s", result, v.Expect, v.Fields, v.Rows)
		}
	}
}
