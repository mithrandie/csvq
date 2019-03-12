package file

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	fileForRead := GetTestFilePath("open.txt")
	fileForUpdate := GetTestFilePath("update.txt")
	fileForCreate := GetTestFilePath("create.txt")

	ctx := context.Background()

	rh, err := NewHandlerForRead(ctx, fileForCreate)
	if err == nil {
		_ = rh.Close()
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	rh, err = NewHandlerForRead(ctx, fileForRead)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = rh.Close()

	uh, err := NewHandlerForUpdate(ctx, fileForCreate)
	if err == nil {
		_ = uh.Close()
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err := NewHandlerForCreate(fileForRead)
	if err == nil {
		_ = ch.Close()
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err = NewHandlerForCreate(fileForCreate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}

	if ch.FileForRead().Name() != fileForCreate {
		_ = uh.Close()
		t.Fatalf("filename to read = %q, expect %q", ch.FileForRead().Name(), fileForCreate)
	}

	if ch.FileForUpdate().Name() != fileForCreate {
		_ = uh.Close()
		t.Fatalf("filename to update = %q, expect %q", ch.FileForUpdate().Name(), fileForCreate)
	}

	rh, err = NewHandlerForRead(ctx, fileForCreate)
	if err == nil {
		_ = rh.Close()
		_ = ch.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		_ = ch.Close()
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	ch2, err := NewHandlerForCreate(fileForCreate)
	if err == nil {
		_ = ch.Close()
		_ = ch2.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*IOError); !ok {
		_ = ch.Close()
		t.Fatalf("error = %#v, want IOError", err)
	}

	_ = ch.Commit()

	rh, err = NewHandlerForRead(ctx, fileForCreate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = rh.Close()

	uh, err = NewHandlerForUpdate(ctx, fileForUpdate)
	if err != nil {
		_ = uh.Close()
		t.Fatalf("error = %#v, expect no error", err)
	}

	if uh.FileForRead().Name() != fileForUpdate {
		_ = uh.Close()
		t.Fatalf("filename to read = %q, expect %q", uh.FileForRead().Name(), fileForUpdate)
	}

	if uh.FileForUpdate().Name() != TempFilePath(fileForUpdate) {
		_ = uh.Close()
		t.Fatalf("filename to update = %q, expect %q", uh.FileForUpdate().Name(), TempFilePath(fileForUpdate))
	}

	rh, err = NewHandlerForRead(ctx, fileForUpdate)
	if err == nil {
		_ = rh.Close()
		_ = uh.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		_ = uh.Close()
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	uh2, err := NewHandlerForUpdate(ctx, fileForUpdate)
	if err == nil {
		_ = uh2.Close()
		_ = uh.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		_ = uh.Close()
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_ = uh.Commit()

	rh, err = NewHandlerForRead(ctx, fileForUpdate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = rh.Close()
}
