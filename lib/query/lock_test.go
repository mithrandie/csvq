package query

import (
	"fmt"
	"github.com/mithrandie/csvq/lib/parser"
	"testing"
	"time"
)

func TestLockFileContainer_Lock(t *testing.T) {
	locks := NewFileLockContainer()
	locks.WaitTimeout = 0.1
	locks.RetryInterval = 10 * time.Millisecond

	path := GetTestFilePath("lockfilecontainer_lock.txt")
	err := locks.LockWithTimeout(parser.Identifier{Literal: "lockfilecontainer_lock.txt"}, path)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	err = locks.LockWithTimeout(parser.Identifier{Literal: "lockfilecontainer_lock.txt"}, path)
	if err == nil {
		t.Errorf("no error, want error for duplicate lock")
	} else {
		expectedErr := fmt.Sprintf("[L:- C:-] file %s: lock wait timeout period exceeded", GetTestFilePath("lockfilecontainer_lock.txt"))
		if err.Error() != expectedErr {
			t.Errorf("error = %s, want error %s", err, expectedErr)
		}
	}
}

func TestFileLockContainer_TryLock(t *testing.T) {
	path := GetTestFilePath("lockfilecontainer_trylock.txt")
	locks := NewFileLockContainer()

	err := locks.TryLock(path)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	err = locks.TryLock(path)
	if err == nil {
		t.Errorf("no error, want error for duplicate lock")
	}
}

func TestFileLockContainer_Unlock(t *testing.T) {
	path := GetTestFilePath("lockfilecontainer_unlock.txt")
	locks := NewFileLockContainer()
	locks.TryLock(path)

	err := locks.Unlock(path)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	err = locks.Unlock(path)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestFileLockContainer_UnlockAll(t *testing.T) {
	path1 := GetTestFilePath("lockfilecontainer_unlock1.txt")
	path2 := GetTestFilePath("lockfilecontainer_unlock2.txt")
	locks := NewFileLockContainer()
	locks.TryLock(path1)
	locks.TryLock(path2)

	err := locks.UnlockAll()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if 0 < len(locks.FileMap) {
		t.Error("file map is not deleted")
	}

	locks.FileMap["DUMMY"] = "notexist.txt"
	err = locks.UnlockAll()
	if err == nil {
		t.Errorf("no error, want error for locks that does not exist")
	}
}

func TestFileLockPath(t *testing.T) {
	path := GetTestFilePath("testfile.txt")
	result := FileLockPath(path)
	expect := GetTestFilePath(".testfile.txt.lock")
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}
