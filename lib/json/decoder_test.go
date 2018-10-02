package json

import (
	"fmt"
	"reflect"
	"testing"
)

var decoderDecodeTests = []struct {
	Input  string
	Expect Structure
	Error  string
}{
	{
		Input:  "  ",
		Expect: nil,
	},
	{
		Input:  " { } ",
		Expect: Object{},
	},
	{
		Input: "{\"key\":\"value\", \"key2\":\"value2\"}",
		Expect: Object{
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
	{
		Input: "{\"key\":\"value\", \"key2\":\"value2\", \"key2\":\"value3\", \"\": \"value4\"}",
		Expect: Object{
			Members: []ObjectMember{
				{
					Key:   "key",
					Value: String("value"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
				{
					Key:   "key2",
					Value: String("value3"),
				},
				{
					Key:   "",
					Value: String("value4"),
				},
			},
		},
	},
	{
		Input:  "[]",
		Expect: Array{},
	},
	{
		Input: "[1, -2.345, -3e3, 4.5e-3, \"ab\\u005cc\\\"de\", false, true, null]",
		Expect: Array{
			Number(1),
			Number(-2.345),
			Number(-3000),
			Number(0.0045),
			String("ab\\c\"de"),
			Boolean(false),
			Boolean(true),
			Null{},
		},
	},
	{
		Input: "{\"key\":{\"child\":\"value\"}}",
		Expect: Object{
			Members: []ObjectMember{
				{
					Key: "key",
					Value: Object{
						Members: []ObjectMember{
							{
								Key:   "child",
								Value: String("value"),
							},
						},
					},
				},
			},
		},
	},
	{
		Input: "{\"key\":[1, 2, 3]}",
		Expect: Object{
			Members: []ObjectMember{
				{
					Key: "key",
					Value: Array{
						Number(1),
						Number(2),
						Number(3),
					},
				},
			},
		},
	},
	{
		Input: "[1, 2, {\"key\":{\"child\":\"value\"}}, 3]",
		Expect: Array{
			Number(1),
			Number(2),
			Object{
				Members: []ObjectMember{
					{
						Key: "key",
						Value: Object{
							Members: []ObjectMember{
								{
									Key:   "child",
									Value: String("value"),
								},
							},
						},
					},
				},
			},
			Number(3),
		},
	},
	{
		Input: "[1, \"abc\", true], []",
		Error: "line 1, column 17: unexpected token \",\"",
	},
	{
		Input: "1",
		Error: "line 1, column 1: unexpected token \"1\"",
	},
	{
		Input: "{\"key\":\"value\", ",
		Error: "line 1, column 16: unexpected termination",
	},
	{
		Input: "{\"key\":\"val\r\nue\", ",
		Error: "line 2, column 5: unexpected termination",
	},
	{
		Input: "{\"key\":\"value }",
		Error: "line 1, column 15: unexpected termination",
	},
	{
		Input: "[1, 1e+500]",
		Error: "line 1, column 11: could not convert \"1e+500\" into float64",
	},
	{
		Input: "[1, -a]",
		Error: "line 1, column 7: invalid number",
	},
	{
		Input: "[1, -1.a]",
		Error: "line 1, column 9: invalid number",
	},
	{
		Input: "[1, -1.1e+a]",
		Error: "line 1, column 12: invalid number",
	},
}

func TestDecoder_Decode(t *testing.T) {
	for _, v := range decoderDecodeTests {
		d := NewDecoder()

		value, err := d.Decode(v.Input)
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
			fmt.Printf("%#v\n", value.(Array)[1].String())
			continue
		}
		if !reflect.DeepEqual(value, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", value, v.Expect, v.Input)
		}
	}
}
