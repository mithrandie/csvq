package file

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	fileForRead := GetTestFilePath("open.txt")
	fileForUpdate := GetTestFilePath("update.txt")
	fileForCreate := GetTestFilePath("create.txt")

	doneCtx, cancel := context.WithCancel(context.Background())
	cancel()

	ctx := context.Background()
	container := NewContainer()
	defer func() {
		if err := container.CloseAllWithErrors(); err != nil {
			t.Log(err)
		}
	}()

	rh, err := NewHandlerForRead(doneCtx, container, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want ContextIsDone")
	}
	if _, ok := err.(*ContextIsDone); !ok {
		t.Fatalf("error = %#v, want ContextIsDone", err)
	}
	_ = container.Close(rh)

	rh, err = NewHandlerForRead(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	rh, err = NewHandlerForRead(ctx, container, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_, err = rh.FileForUpdate()
	if err == nil {
		t.Fatalf("no error, want error")
	}

	_, err = NewHandlerForRead(ctx, container, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want error")
	}
	_ = container.Close(rh)

	uh, err := NewHandlerForUpdate(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err := NewHandlerForCreate(container, fileForRead)
	if err == nil {
		t.Fatalf("no error, want IOError")
	}
	if _, ok := err.(*IOError); !ok {
		t.Fatalf("error = %#v, want IOError", err)
	}

	ch, err = NewHandlerForCreate(container, fileForCreate)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}

	if ch.File().Name() != fileForCreate {
		t.Fatalf("filename to read = %q, expect %q", ch.File().Name(), fileForCreate)
	}

	fp, err := ch.FileForUpdate()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if fp.Name() != fileForCreate {
		t.Fatalf("filename to update = %q, expect %q", fp.Name(), fileForCreate)
	}

	rh, err = NewHandlerForRead(ctx, container, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_, err = NewHandlerForCreate(container, fileForCreate)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*IOError); !ok {
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
		t.Fatalf("error = %#v, expect no error", err)
	}

	if uh.File().Name() != fileForUpdate {
		t.Fatalf("filename to read = %q, expect %q", uh.File().Name(), fileForUpdate)
	}

	fp, err = uh.FileForUpdate()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if fp.Name() != TempFilePath(fileForUpdate) {
		t.Fatalf("filename to update = %q, expect %q", fp.Name(), TempFilePath(fileForUpdate))
	}

	rh, err = NewHandlerForRead(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_, err = NewHandlerForUpdate(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_ = container.Commit(uh)

	rh, err = NewHandlerForRead(ctx, container, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
}
