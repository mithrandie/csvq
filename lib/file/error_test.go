package file

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mithrandie/go-file/v2"
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

func TestForcedUnlockError_Error(t *testing.T) {
	errs := []error{
		errors.New("err1"),
		errors.New("err2"),
	}

	err := NewForcedUnlockError(errs)
	expect := "err1\n  err2"
	if err == nil {
		t.Fatalf("no error, want error %q", expect)
	}
	if err.Error() != expect {
		t.Fatalf("err = %q, want error %q", err.Error(), expect)
	}

	err = NewForcedUnlockError(nil)
	if err != nil {
		t.Fatalf("error = %q, want no error", err.Error())
	}
}

func TestCompositeError_Error(t *testing.T) {
	err1 := NewForcedUnlockError([]error{
		errors.New("ferr1"),
		errors.New("ferr2"),
	})
	err2 := errors.New("err2")

	err := NewCompositeError(err1, err2)
	expect := "ferr1\n  ferr2\n  err2"
	if err == nil {
		t.Fatalf("no error, want error %q", expect)
	}
	if err.Error() != expect {
		t.Fatalf("err = %q, want error %q", err.Error(), expect)
	}

	err = NewCompositeError(err1, nil)
	expect = "ferr1\n  ferr2"
	if err == nil {
		t.Fatalf("no error, want error %q", expect)
	}
	if err.Error() != expect {
		t.Fatalf("err = %q, want error %q", err.Error(), expect)
	}

	err = NewCompositeError(nil, err2)
	expect = "err2"
	if err == nil {
		t.Fatalf("no error, want error %q", expect)
	}
	if err.Error() != expect {
		t.Fatalf("err = %q, want error %q", err.Error(), expect)
	}

	err = NewCompositeError(nil, nil)
	if err != nil {
		t.Fatalf("error = %q, want no error", err.Error())
	}
}
