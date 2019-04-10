package file

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mithrandie/go-file/v2"
)

type OpenType int

const (
	ForRead OpenType = iota
	ForCreate
	ForUpdate
)

type mngFileType int

const (
	fileTypeRLock mngFileType = iota
	fileTypeLock
	fileTypeTemp
)

var mngFileTypeLit = map[mngFileType]string{
	fileTypeRLock: "read lock",
	fileTypeLock:  "lock",
	fileTypeTemp:  "temporary",
}

func (t mngFileType) String() string {
	return mngFileTypeLit[t]
}

type mngFile struct {
	path string
	fp   *os.File
}

func newMngFile(path string, fp *os.File) *mngFile {
	return &mngFile{
		path: path,
		fp:   fp,
	}
}

func (m *mngFile) close() error {
	if m != nil {
		if m.fp != nil {
			if err := file.Close(m.fp); err != nil {
				return err
			}
		}

		if Exists(m.path) {
			if err := os.Remove(m.path); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *mngFile) closeWithErrors() []error {
	var errs []error
	if m != nil {
		if m.fp != nil {
			if err := file.Close(m.fp); err != nil {
				errs = append(errs, err)
			}
		}

		if Exists(m.path) {
			if err := os.Remove(m.path); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

type Handler struct {
	path string
	fp   *os.File

	openType OpenType

	rlockFile *mngFile
	lockFile  *mngFile
	tempFile  *mngFile

	closed bool
}

func NewHandlerWithoutLock(ctx context.Context, container *Container, path string, defaultWaitTimeout time.Duration, retryDelay time.Duration) (*Handler, error) {
	tctx, cancel := GetTimeoutContext(ctx, defaultWaitTimeout)
	defer cancel()

	h := &Handler{
		path:     path,
		openType: ForRead,
	}

	if !Exists(h.path) {
		return h, NewNotExistError(fmt.Sprintf("file %s does not exist", h.path))
	}

	fp, err := file.OpenToReadContext(tctx, retryDelay, h.path)
	if err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	h.fp = fp

	if err := container.Add(h.path, h); err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	return h, nil
}

func NewHandlerForRead(ctx context.Context, container *Container, path string, defaultWaitTimeout time.Duration, retryDelay time.Duration) (*Handler, error) {
	tctx, cancel := GetTimeoutContext(ctx, defaultWaitTimeout)
	defer cancel()

	h := &Handler{
		path:     path,
		openType: ForRead,
	}

	if !Exists(h.path) {
		return h, NewNotExistError(fmt.Sprintf("file %s does not exist", h.path))
	}

	if err := h.CreateManagementFileContext(tctx, retryDelay, fileTypeRLock); err != nil {
		return h, closeIsolatedHandler(h, err)
	}

	fp, err := file.OpenToReadContext(tctx, retryDelay, h.path)
	if err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	h.fp = fp

	if err := container.Add(h.path, h); err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	return h, nil
}

func NewHandlerForCreate(container *Container, path string) (*Handler, error) {
	h := &Handler{
		path:     path,
		openType: ForCreate,
	}

	if Exists(h.path) {
		return h, NewAlreadyExistError(fmt.Sprintf("file %s already exists", h.path))
	}

	if err := h.tryCreateManagementFile(fileTypeLock); err != nil {
		return h, closeIsolatedHandler(h, err)
	}

	fp, err := file.Create(h.path)
	if err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	h.fp = fp

	if err := container.Add(h.path, h); err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	return h, nil
}

func NewHandlerForUpdate(ctx context.Context, container *Container, path string, defaultWaitTimeout time.Duration, retryDelay time.Duration) (*Handler, error) {
	tctx, cancel := GetTimeoutContext(ctx, defaultWaitTimeout)
	defer cancel()

	h := &Handler{
		path:     path,
		openType: ForUpdate,
	}

	if !Exists(h.path) {
		return h, NewNotExistError(fmt.Sprintf("file %s does not exist", h.path))
	}

	if err := h.CreateManagementFileContext(tctx, retryDelay, fileTypeLock); err != nil {
		return h, closeIsolatedHandler(h, err)
	}

	fp, err := file.OpenToUpdateContext(tctx, retryDelay, path)
	if err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	h.fp = fp

	if err := h.CreateManagementFileContext(tctx, retryDelay, fileTypeTemp); err != nil {
		return h, closeIsolatedHandler(h, err)
	}

	if err := container.Add(h.path, h); err != nil {
		return h, closeIsolatedHandler(h, err)
	}
	return h, nil
}

func closeIsolatedHandler(h *Handler, err error) error {
	return NewCompositeError(ParseError(err), h.closeWithErrors())
}

func (h *Handler) Path() string {
	return h.path
}

func (h *Handler) File() *os.File {
	return h.fp
}

func (h *Handler) FileForUpdate() (*os.File, error) {
	switch h.openType {
	case ForUpdate:
		return h.tempFile.fp, nil
	case ForCreate:
		return h.fp, nil
	}
	return nil, fmt.Errorf("file %s cannot be updated", h.path)
}

func (h *Handler) close() error {
	if h.closed {
		return nil
	}

	if h.fp != nil {
		if err := file.Close(h.fp); err != nil {
			return err
		}
		h.fp = nil
	}

	if h.openType == ForCreate && Exists(h.path) {
		if err := os.Remove(h.path); err != nil {
			return err
		}
	}

	if err := h.tempFile.close(); err != nil {
		return err
	}
	h.tempFile = nil

	if err := h.lockFile.close(); err != nil {
		return err
	}
	h.lockFile = nil

	if err := h.rlockFile.close(); err != nil {
		return err
	}
	h.rlockFile = nil

	h.closed = true
	return nil
}

func (h *Handler) commit() error {
	if h.closed {
		return nil
	}

	if h.fp != nil {
		if err := file.Close(h.fp); err != nil {
			return err
		}
		h.fp = nil
	}

	if h.openType == ForUpdate {
		if h.tempFile.fp != nil {
			if err := file.Close(h.tempFile.fp); err != nil {
				return err
			}
			h.tempFile.fp = nil
		}

		if Exists(h.path) {
			if err := os.Remove(h.path); err != nil {
				return err
			}
		}

		if err := os.Rename(h.tempFile.path, h.path); err != nil {
			return err
		}
	} else {
		if err := h.tempFile.close(); err != nil {
			return err
		}
		h.tempFile = nil
	}

	if err := h.lockFile.close(); err != nil {
		return err
	}
	h.lockFile = nil

	if err := h.rlockFile.close(); err != nil {
		return err
	}
	h.rlockFile = nil

	h.closed = true
	return nil
}

func (h *Handler) closeWithErrors() error {
	if h.closed {
		return nil
	}

	var errs []error

	if h.fp != nil {
		if err := file.Close(h.fp); err != nil {
			errs = append(errs, err)
		} else {
			h.fp = nil
		}
	}

	if h.openType == ForCreate && Exists(h.path) {
		if err := os.Remove(h.path); err != nil {
			errs = append(errs, err)
		}
	}

	if cerrs := h.tempFile.closeWithErrors(); cerrs != nil {
		errs = append(errs, cerrs...)
	} else {
		h.tempFile = nil
	}

	if cerrs := h.lockFile.closeWithErrors(); cerrs != nil {
		errs = append(errs, cerrs...)
	} else {
		h.lockFile = nil
	}

	if cerrs := h.rlockFile.closeWithErrors(); cerrs != nil {
		errs = append(errs, cerrs...)
	} else {
		h.rlockFile = nil
	}

	return NewForcedUnlockError(errs)
}

func (h *Handler) CreateManagementFileContext(ctx context.Context, retryDelay time.Duration, fileType mngFileType) error {
	if ctx.Err() != nil {
		if ctx.Err() == context.Canceled {
			return NewContextCanceled()
		}
		return NewContextDone(ctx.Err().Error())
	}

	for {
		if err := h.tryCreateManagementFile(fileType); err != nil {
			if _, ok := err.(*LockError); !ok {
				return err
			}
		} else {
			return nil
		}

		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				return NewContextCanceled()
			}
			return NewTimeoutError(h.path)
		case <-time.After(retryDelay):
			// try again
		}
	}
}

func (h *Handler) tryCreateManagementFile(fileType mngFileType) error {
	if len(h.path) < 1 {
		return NewLockError("filename not specified")
	}

	switch fileType {
	case fileTypeLock:
		return h.tryCreateLockFile()
	case fileTypeTemp:
		return h.tryCreateTempFile()
	default: //fileTypeRLock
		return h.tryCreateRLockFile()
	}
}

func (h *Handler) tryCreateRLockFile() (err error) {
	if h.rlockFile != nil {
		return NewLockError(fmt.Sprintf("%s file for %s is already created", fileTypeRLock, h.path))
	}

	lockFilePath := LockFilePath(h.path)
	if LockExists(h.path) {
		return NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeRLock, h.path))
	}

	lfp, err := file.Create(lockFilePath)
	if err != nil {
		return NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeRLock, h.path))
	}
	lockFile := newMngFile(lockFilePath, lfp)
	defer func() {
		err = NewCompositeError(err, lockFile.close())
	}()

	filePath := RLockFilePath(h.path)
	fp, e := file.Create(filePath)
	if e != nil {
		err = NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeRLock, h.path))
		return
	}

	h.rlockFile = newMngFile(filePath, fp)
	return
}

func (h *Handler) tryCreateLockFile() error {
	if h.lockFile != nil {
		return NewLockError(fmt.Sprintf("%s file for %s is already created", fileTypeLock, h.path))
	}

	filePath := LockFilePath(h.path)
	if LockExists(h.path) || RLockExists(h.path) {
		return NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeLock, h.path))
	}

	fp, err := file.Create(filePath)
	if err != nil {
		return NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeLock, h.path))
	}
	h.lockFile = newMngFile(filePath, fp)

	if RLockExists(h.path) {
		return NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeLock, h.path))
	}

	return nil
}

func (h *Handler) tryCreateTempFile() error {
	if h.tempFile != nil {
		return NewLockError(fmt.Sprintf("%s file for %s is already created", fileTypeTemp, h.path))
	}

	filePath := TempFilePath(h.path)
	fp, err := file.Create(filePath)
	if err != nil {
		return NewLockError(fmt.Sprintf("failed to create %s file for %q", fileTypeTemp, h.path))
	}

	h.tempFile = newMngFile(filePath, fp)
	return nil
}
