package file

import "testing"

func TestLockFilePath(t *testing.T) {
	path := GetTestFilePath("testfile.txt")
	result := LockFilePath(path)
	expect := GetTestFilePath(".testfile.txt" + LockFileSuffix)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

func TestTempFilePath(t *testing.T) {
	path := GetTestFilePath("testfile.txt")
	result := TempFilePath(path)
	expect := GetTestFilePath(".testfile.txt" + TempFileSuffix)
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}
