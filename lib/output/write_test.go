package output

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	file := "test.txt"
	s := "test"

	Write(file, s)

	_, err := os.Stat(file)
	if err != nil {
		t.Errorf("file %q does not get created", file)
	}

	if err := os.Remove(file); err != nil {
		t.Errorf("unexpected error %q", err)
	}
}

func ExampleWrite() {
	Write("", "write test")
	//Output:
	//write test
}
