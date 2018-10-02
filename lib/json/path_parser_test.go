package json

import (
	"reflect"
	"testing"
)

var parsePathTests = []struct {
	Input  string
	Expect PathExpression
	Error  string
}{
	{
		Input:  "",
		Expect: ObjectPath{},
	},
	{
		Input: "abc",
		Expect: ObjectPath{
			Name: "abc",
		},
	},
	{
		Input: "abc\\def\\",
		Expect: ObjectPath{
			Name: "abc\\def\\",
		},
	},
	{
		Input: "abc.d\\.ef",
		Expect: ObjectPath{
			Name: "abc",
			Child: ObjectPath{
				Name: "d.ef",
			},
		},
	},
	{
		Input: "abc.",
		Error: "unexpected termination",
	},
	{
		Input: "abc..",
		Error: "unexpected token \".\"",
	},
}

func TestParsePath(t *testing.T) {
	for _, v := range parsePathTests {
		result, err := ParsePath(v.Input)
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
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", result, v.Expect, v.Input)
		}
	}
}
