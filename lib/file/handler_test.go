package file

import (
	"testing"
)

func TestHandler(t *testing.T) {
	fileForRead := GetTestFilePath("open.txt")
	fileForUpdate := GetTestFilePath("update.txt")
	fileForCreate := GetTestFilePath("create.txt")

	rh, err := NewHandlerForRead(fileForCreate)
	if err == nil {
		rh.Close()
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	rh, err = NewHandlerForRead(fileForRead)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	rh.Close()

	uh, err := NewHandlerForUpdate(fileForCreate)
	if err == nil {
		uh.Close()
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err := NewHandlerForCreate(fileForRead)
	if err == nil {
		ch.Close()
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
		uh.Close()
		t.Fatalf("filename to read = %q, expect %q", ch.FileForRead().Name(), fileForCreate)
	}

	if ch.FileForUpdate().Name() != fileForCreate {
		uh.Close()
		t.Fatalf("filename to update = %q, expect %q", ch.FileForUpdate().Name(), fileForCreate)
	}

	rh, err = NewHandlerForRead(fileForCreate)
	if err == nil {
		rh.Close()
		ch.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		ch.Close()
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	ch2, err := NewHandlerForCreate(fileForCreate)
	if err == nil {
		ch.Close()
		ch2.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*IOError); !ok {
		ch.Close()
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch.Commit()

	rh, err = NewHandlerForRead(fileForCreate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	rh.Close()

	uh, err = NewHandlerForUpdate(fileForUpdate)
	if err != nil {
		uh.Close()
		t.Fatalf("error = %#v, expect no error", err)
	}

	if uh.FileForRead().Name() != fileForUpdate {
		uh.Close()
		t.Fatalf("filename to read = %q, expect %q", uh.FileForRead().Name(), fileForUpdate)
	}

	if uh.FileForUpdate().Name() != TempFilePath(fileForUpdate) {
		uh.Close()
		t.Fatalf("filename to update = %q, expect %q", uh.FileForUpdate().Name(), TempFilePath(fileForUpdate))
	}

	rh, err = NewHandlerForRead(fileForUpdate)
	if err == nil {
		rh.Close()
		uh.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		uh.Close()
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	uh2, err := NewHandlerForUpdate(fileForUpdate)
	if err == nil {
		uh2.Close()
		uh.Close()
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		uh.Close()
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	uh.Commit()

	rh, err = NewHandlerForRead(fileForUpdate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	rh.Close()
}
