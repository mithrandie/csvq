package value

import (
	"reflect"
	"testing"
)

var fieldToPrimaryTests = []struct {
	Field  Field
	Expect Primary
}{
	{
		Field:  nil,
		Expect: NewNull(),
	},
	{
		Field:  NewField(""),
		Expect: NewString(""),
	},
	{
		Field:  NewField("abc"),
		Expect: NewString("abc"),
	},
}

func TestField_ToPrimary(t *testing.T) {
	for _, v := range fieldToPrimaryTests {
		result := v.Field.ToPrimary()
		if !reflect.DeepEqual(result, v.Expect) {
			t.Errorf("result = %#v, want %#v for %#v", result, v.Expect, v.Field)
		}
	}
}
