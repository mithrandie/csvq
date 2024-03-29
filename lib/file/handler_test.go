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

	h, err := container.CreateHandlerWithoutLock(ctx, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want NotExistError")
	}
	if _, ok := err.(*NotExistError); !ok {
		t.Fatalf("error = %#v, want NotExistError", err)
	}
	_ = container.Close(h)

	h, err = container.CreateHandlerWithoutLock(ctx, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("unexpected error %#v", err)
	}
	_ = container.Close(h)

	rh, err := container.CreateHandlerForRead(doneCtx, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want ContextCanceled")
	}
	if _, ok := err.(*ContextCanceled); !ok {
		t.Fatalf("error = %#v, want ContextCanceled", err)
	}
	_ = container.Close(rh)

	rh, err = container.CreateHandlerForRead(ctx, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want NotExistError")
	}
	if _, ok := err.(*NotExistError); !ok {
		t.Fatalf("error = %#v, want NotExistError", err)
	}

	rh, err = container.CreateHandlerForRead(ctx, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_, err = rh.FileForUpdate()
	if err == nil {
		t.Fatalf("no error, want error")
	}

	_, err = container.CreateHandlerForRead(ctx, fileForRead, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want error")
	}
	_ = container.Close(rh)

	uh, err := container.CreateHandlerForUpdate(ctx, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want NotExistError")
	}
	if _, ok := err.(*NotExistError); !ok {
		t.Fatalf("error = %#v, want NotExistError", err)
	}

	ch, err := container.CreateHandlerForCreate(fileForRead)
	if err == nil {
		t.Fatalf("no error, want AlreadyExistError")
	}
	if _, ok := err.(*AlreadyExistError); !ok {
		t.Fatalf("error = %#v, want AlreadyExistError", err)
	}

	ch, err = container.CreateHandlerForCreate(fileForCreate)
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

	rh, err = container.CreateHandlerForRead(ctx, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_, err = container.CreateHandlerForCreate(fileForCreate)
	if err == nil {
		t.Fatalf("no error, want AlreadyExistError")
	}
	if _, ok := err.(*AlreadyExistError); !ok {
		t.Fatalf("error = %#v, want AlreadyExistError", err)
	}

	_ = container.Commit(ch)

	rh, err = container.CreateHandlerForRead(ctx, fileForCreate, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
	_ = container.Close(rh)

	uh, err = container.CreateHandlerForUpdate(ctx, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
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

	rh, err = container.CreateHandlerForRead(ctx, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_, err = container.CreateHandlerForUpdate(ctx, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err == nil {
		t.Fatalf("no error, want TimeoutError")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("error = %#v, want TimeoutError", err)
	}

	_ = container.Commit(uh)

	rh, err = container.CreateHandlerForRead(ctx, fileForUpdate, waitTimeoutForTests, retryDelayForTests)
	if err != nil {
		t.Fatalf("error = %#v, expect no error", err)
	}
}
