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
	container := NewContainer()

	rh, err := NewHandlerForRead(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		_ = container.Close(rh)
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	rh, err = NewHandlerForRead(ctx, container, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = container.Close(rh)

	uh, err := NewHandlerForUpdate(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		_ = container.Close(uh)
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err := NewHandlerForCreate(container, fileForRead)
	if err == nil {
		_ = container.Close(ch)
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err = NewHandlerForCreate(container, fileForCreate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}

	if ch.FileForRead().Name() != fileForCreate {
		_ = container.Close(uh)
		t.Fatalf("filename to read = %q, expect %q", ch.FileForRead().Name(), fileForCreate)
	}

	if ch.FileForUpdate().Name() != fileForCreate {
		_ = container.Close(uh)
		t.Fatalf("filename to update = %q, expect %q", ch.FileForUpdate().Name(), fileForCreate)
	}

	rh, err = NewHandlerForRead(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		_ = container.Close(rh)
		_ = container.Close(ch)
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		_ = container.Close(ch)
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	ch2, err := NewHandlerForCreate(container, fileForCreate)
	if err == nil {
		_ = container.Close(ch)
		_ = container.Close(ch2)
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*IOError); !ok {
		_ = container.Close(ch)
		t.Fatalf("error = %#v, want IOError", err)
	}

	_ = container.Commit(ch)

	rh, err = NewHandlerForRead(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = container.Close(rh)

	uh, err = NewHandlerForUpdate(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		_ = container.Close(uh)
		t.Fatalf("error = %#v, expect no error", err)
	}

	if uh.FileForRead().Name() != fileForUpdate {
		_ = container.Close(uh)
		t.Fatalf("filename to read = %q, expect %q", uh.FileForRead().Name(), fileForUpdate)
	}

	if uh.FileForUpdate().Name() != TempFilePath(fileForUpdate) {
		_ = container.Close(uh)
		t.Fatalf("filename to update = %q, expect %q", uh.FileForUpdate().Name(), TempFilePath(fileForUpdate))
	}

	rh, err = NewHandlerForRead(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		_ = container.Close(rh)
		_ = container.Close(uh)
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		_ = container.Close(uh)
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	uh2, err := NewHandlerForUpdate(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		_ = container.Close(uh2)
		_ = container.Close(uh)
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		_ = container.Close(uh)
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_ = container.Commit(uh)

	rh, err = NewHandlerForRead(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = container.Close(rh)
}
