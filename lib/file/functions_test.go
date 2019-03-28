package file

import (
	"context"
	"regexp"
	"testing"
	"time"
)

func TestGetTimeoutContext(t *testing.T) {
	ctx := context.Background()
	result, cancel := GetTimeoutContext(ctx, 10*time.Second)
	tm1, ok := result.Deadline()
	if !ok {
		t.Fatalf("deadline of context does not set")
	}

	result2, _ := GetTimeoutContext(result, 100*time.Second)
	tm2, ok := result2.Deadline()
	if !ok {
		t.Fatalf("deadline of context does not set")
	}
	if !tm1.Equal(tm2) {
		t.Fatalf("deadline of context is changed")
	}

	cancel()
	result3, _ := GetTimeoutContext(result2, 100*time.Second)
	if result3.Err() == nil {
		t.Fatalf("context is not done")
	}
}

func TestRLockFilePath(t *testing.T) {
	path := GetTestFilePath("testfile.txt")
	result := RLockFilePath(path)
	expect := GetTestFilePath(".testfile.txt" + ".[0-9a-zA-Z]{12}" + RLockFileSuffix)
	r := regexp.MustCompile(expect)
	if !r.MatchString(result) {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

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
