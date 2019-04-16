package query

import (
	"errors"
	"testing"
)

var appendCompositeErrorTests = []struct {
	Err1   error
	Err2   error
	Expect string
}{
	{
		Err1:   nil,
		Err2:   errors.New("error2"),
		Expect: "error2",
	},
	{
		Err1:   errors.New("error1"),
		Err2:   nil,
		Expect: "error1",
	},
	{
		Err1: errors.New("error1"),
		Err2: errors.New("error2"),
		Expect: "composite error:" +
			"\n  [System Error] error1" +
			"\n  [System Error] error2",
	},
}

func TestAppendCompositeError(t *testing.T) {
	for _, v := range appendCompositeErrorTests {
		result := appendCompositeError(v.Err1, v.Err2)
		if result == nil || result.Error() != v.Expect {
			t.Errorf("result = %s, want %s for %v, %v", result, v.Expect, v.Err1, v.Err2)
		}
	}
}
