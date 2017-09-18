package file

import (
	"fmt"
	"testing"
)

func TestOpenToRead(t *testing.T) {
	LockFiles = make(LockFileContainer)
	defer UnlockAll()

	expectErr := fmt.Sprintf("file %s: lock waiting time exceeded", GetTestFilePath("locked_by_other.txt"))
	_, err := OpenToRead(GetTestFilePath("locked_by_other.txt"))
	if err == nil {
		t.Fatalf("no error, want error for lock timeout")
	}
	if err.Error() != expectErr {
		t.Fatalf("error = %q, want error %q", err.Error(), expectErr)
	}

	fp, err := OpenToRead(GetTestFilePath("notexist.txt"))
	defer Close(fp)
	if err == nil {
		t.Fatalf("no error, want error for open error")
	}

	fp, err = OpenToRead(GetTestFilePath("open.txt"))
	defer Close(fp)
	if err != nil {
		t.Fatalf("error = %q, want no error", err)
	}
}

func TestOpenToUpdate(t *testing.T) {
	LockFiles = make(LockFileContainer)
	defer UnlockAll()

	expectErr := fmt.Sprintf("file %s: lock waiting time exceeded", GetTestFilePath("locked_by_other.txt"))
	_, err := OpenToUpdate(GetTestFilePath("locked_by_other.txt"))
	if err == nil {
		t.Fatalf("no error, want error for lock timeout")
	}
	if err.Error() != expectErr {
		t.Fatalf("error = %q, want error %q", err.Error(), expectErr)
	}

	fp, err := OpenToUpdate(GetTestFilePath("notexist.txt"))
	defer Close(fp)
	if err == nil {
		t.Fatalf("no error, want error for open error")
	}

	fp, err = OpenToUpdate(GetTestFilePath("update.txt"))
	defer Close(fp)
	if err != nil {
		t.Fatalf("error = %q, want no error", err)
	}

	if len(LockFiles) != 1 {
		t.Fatalf("exactly one lock file must exist")
	}
}

func TestCanRead(t *testing.T) {
	LockFiles = make(LockFileContainer)
	defer UnlockAll()

	path := GetTestFilePath("canread.txt")

	b := CanRead(path)
	if !b {
		t.Fatalf("can read = %t, want %t", b, !b)
	}

	LockWithTimeout(path)
	b = CanRead(path)
	if !b {
		t.Fatalf("can read = %t, want %t", b, !b)
	}

	b = CanRead(GetTestFilePath("locked_by_other.txt"))
	if b {
		t.Fatalf("can read = %t, want %t", b, !b)
	}
}

func TestLockWithTimeout(t *testing.T) {
	LockFiles = make(LockFileContainer)
	defer UnlockAll()

	path := GetTestFilePath("trylock.txt")

	err := LockWithTimeout(path)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(LockFiles) != 1 {
		t.Fatalf("exactly one lock file must exist")
	}

	err = LockWithTimeout(path)
	if err == nil {
		t.Fatalf("no error, want error for duplicate lock")
	}
}

func TestTryLock(t *testing.T) {
	LockFiles = make(LockFileContainer)
	defer UnlockAll()

	path := GetTestFilePath("trylock.txt")

	err := TryLock(path)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(LockFiles) != 1 {
		t.Fatalf("exactly one lock file must exist")
	}

	err = TryLock(path)
	if err == nil {
		t.Fatalf("no error, want error for duplicate lock")
	}
}

func TestLockFilePath(t *testing.T) {
	path := GetTestFilePath("testfile.txt")
	result := LockFilePath(path)
	expect := GetTestFilePath(".testfile.txt.lock")
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}
