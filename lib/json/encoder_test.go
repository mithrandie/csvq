package json

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"testing"
)

var encoderEncodeTests = []struct {
	Input       Structure
	Escape      EscapeType
	PrettyPrint bool
	LineBreak   cmd.LineBreak
	Expect      string
}{
	{
		Input:       String("abc"),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "\"abc\"",
	},
	{
		Input:       Number(-1.234),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "-1.234",
	},
	{
		Input:       Boolean(true),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "true",
	},
	{
		Input:       Boolean(false),
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "false",
	},
	{
		Input:       Null{},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "null",
	},
	{
		Input: Array{
			String("value1"),
			String("value2"),
			String("value3"),
		},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "[\"value1\",\"value2\",\"value3\"]",
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
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "{\"key1\":\"value1\",\"key2\":\"value2\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key\"1",
					Value: String("value\"1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "{\"key\\\"1\":\"value\\\"1\",\"key2\":\"value2\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key\"1",
					Value: String("value\"1"),
				},
				{
					Key:   "key2",
					Value: String("value2"),
				},
			},
		},
		Escape:      HexDigits,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "{\"key\\u00221\":\"value\\u00221\",\"key2\":\"value2\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key\"1",
					Value: String("value\"1"),
				},
			},
		},
		Escape:      AllWithHexDigits,
		PrettyPrint: false,
		LineBreak:   cmd.LF,
		Expect:      "{\"\\u006b\\u0065\\u0079\\u0022\\u0031\":\"\\u0076\\u0061\\u006c\\u0075\\u0065\\u0022\\u0031\"}",
	},
	{
		Input: Object{
			Members: []ObjectMember{
				{
					Key:   "key1",
					Value: String("value1"),
				},
				{
					Key: "key2",
					Value: Array{
						Object{
							Members: []ObjectMember{
								{
									Key:   "akey1",
									Value: Boolean(true),
								},
								{
									Key:   "akey2",
									Value: Null{},
								},
							},
						},
						Object{
							Members: []ObjectMember{
								{
									Key:   "akey1",
									Value: Number(-2.3e-6),
								},
								{
									Key: "akey2",
									Value: Array{
										String("A"),
										String("B"),
										String("C"),
									},
								},
							},
						},
					},
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   cmd.LF,
		Expect: "{\n" +
			"  \"key1\": \"value1\",\n" +
			"  \"key2\": [\n" +
			"    {\n" +
			"      \"akey1\": true,\n" +
			"      \"akey2\": null\n" +
			"    },\n" +
			"    {\n" +
			"      \"akey1\": -0.0000023,\n" +
			"      \"akey2\": [\n" +
			"        \"A\",\n" +
			"        \"B\",\n" +
			"        \"C\"\n" +
			"      ]\n" +
			"    }\n" +
			"  ]\n" +
			"}",
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
					Value: String("[1, 2, 3]"),
				},
			},
		},
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   cmd.LF,
		Expect: "{\n" +
			"  \"key1\": \"value1\",\n" +
			"  \"key2\": [\n" +
			"    1,\n" +
			"    2,\n" +
			"    3\n" +
			"  ]\n" +
			"}",
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
		Escape:      Backslash,
		PrettyPrint: true,
		LineBreak:   cmd.CRLF,
		Expect: "{\r\n" +
			"  \"key1\": \"value1\",\r\n" +
			"  \"key2\": \"value2\"\r\n" +
			"}",
	},
}

func TestEncoder_Encode(t *testing.T) {
	for _, v := range encoderEncodeTests {
		e := NewEncoder()

		e.EscapeType = v.Escape
		e.PrettyPrint = v.PrettyPrint
		e.LineBreak = v.LineBreak

		result := e.Encode(v.Input)
		if result != v.Expect {
			t.Errorf("result = %q, want %q for EscapeType:%d PrettyPrint:%t Input:%#v", result, v.Expect, v.Escape, v.PrettyPrint, v.Input)
		}
	}
}
