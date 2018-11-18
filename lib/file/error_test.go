package file

import (
	"reflect"
	"testing"

	"github.com/mithrandie/go-file"
)

var parseErrorTests = []struct {
	Error       error
	ExpectError error
	Message     string
}{
	{
		Error:       file.NewIOError("io error"),
		ExpectError: NewIOError("io error"),
		Message:     "io error",
	},
	{
		Error:       file.NewLockError("lock error"),
		ExpectError: NewLockError("lock error"),
		Message:     "lock error",
	},
	{
		Error:       file.NewTimeoutError("filepath"),
		ExpectError: &TimeoutError{message: "file filepath: lock waiting time exceeded"},
		Message:     "file filepath: lock waiting time exceeded",
	},
	{
		Error:       NewTimeoutError("filepath"),
		ExpectError: &TimeoutError{message: "file filepath: lock waiting time exceeded"},
		Message:     "file filepath: lock waiting time exceeded",
	},
}

func TestParseError(t *testing.T) {
	for _, v := range parseErrorTests {
		e := ParseError(v.Error)
		if !reflect.DeepEqual(e, v.ExpectError) {
			t.Errorf("result = %#v, want %#v for %#v", e, v.ExpectError, v.Error)
		}
		if e.Error() != v.Message {
			t.Errorf("message = %q, want %q for %#v", e.Error(), v.Message, v.Error)
		}
	}
}
