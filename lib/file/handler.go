package file

import (
	"fmt"
	"os"
	"time"

	"github.com/mithrandie/go-file"
)

type OpenType int

const (
	ForRead OpenType = iota
	ForCreate
	ForUpdate
)

type Handler struct {
	path string
	fp   *os.File

	openType OpenType

	lockFilePath string
	lockFileFp   *os.File

	tempFilePath string
	tempFp       *os.File

	closed bool
}

func NewHandlerForRead(path string) (*Handler, error) {
	h := &Handler{
		path:     path,
		openType: ForRead,
	}

	if err := h.PrepareToRead(); err != nil {
		return h, err
	}

	fp, err := file.OpenToReadWithTimeout(h.path)
	if err != nil {
		return h, ParseError(err)
	}
	h.fp = fp

	return h, nil
}

func NewHandlerForCreate(path string) (*Handler, error) {
	h := &Handler{
		path:     path,
		openType: ForCreate,
	}

	if Exists(h.path) {
		return h, NewIOError(fmt.Sprintf("file %s already exists", h.path))
	}

	if err := h.TryCreateLockFile(); err != nil {
		return h, err
	}

	fp, err := file.Create(h.path)
	if err != nil {
		return h, ParseError(err)
	}
	h.fp = fp

	if err := addToContainer(h.path, h); err != nil {
		return h, err
	}
	return h, nil
}

func NewHandlerForUpdate(path string) (*Handler, error) {
	h := &Handler{
		path:     path,
		openType: ForUpdate,
	}

	if !Exists(h.path) {
		return h, NewIOError(fmt.Sprintf("file %s does not exist", h.path))
	}

	if err := h.TryCreateLockFileWithTimeout(); err != nil {
		return h, err
	}

	fp, err := file.OpenToUpdateWithTimeout(path)
	if err != nil {
		h.Close()
		return h, ParseError(err)
	}
	h.fp = fp

	if err := h.TryCreateTempFile(); err != nil {
		return h, err
	}

	if err := addToContainer(h.path, h); err != nil {
		return h, err
	}
	return h, nil
}

func (h *Handler) Path() string {
	return h.path
}

func (h *Handler) FileForRead() *os.File {
	return h.fp
}

func (h *Handler) FileForUpdate() *os.File {
	if h.openType == ForUpdate {
		return h.tempFp
	}
	return h.fp
}

func (h *Handler) Close() error {
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

	if h.tempFp != nil {
		if err := file.Close(h.tempFp); err != nil {
			return err
		}
		h.tempFp = nil
	}

	if Exists(h.tempFilePath) {
		if err := os.Remove(h.tempFilePath); err != nil {
			return err
		}
	}

	if h.lockFileFp != nil {
		if err := file.Close(h.lockFileFp); err != nil {
			return err
		}
		h.lockFileFp = nil
	}

	if Exists(h.lockFilePath) {
		if err := os.Remove(h.lockFilePath); err != nil {
			return err
		}
	}

	h.closed = true
	removeFromContainer(h.path)
	return nil
}

func (h *Handler) Commit() error {
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
		if h.tempFp != nil {
			if err := file.Close(h.tempFp); err != nil {
				return err
			}
			h.tempFp = nil
		}

		if Exists(h.path) {
			if err := os.Remove(h.path); err != nil {
				return err
			}
		}

		if err := os.Rename(h.tempFilePath, h.path); err != nil {
			return err
		}
	} else {
		if h.tempFp != nil {
			if err := file.Close(h.tempFp); err != nil {
				return err
			}
			h.tempFp = nil
		}

		if Exists(h.tempFilePath) {
			if err := os.Remove(h.tempFilePath); err != nil {
				return err
			}
		}
	}

	if h.lockFileFp != nil {
		if err := file.Close(h.lockFileFp); err != nil {
			return err
		}
		h.lockFileFp = nil
	}

	if Exists(h.lockFilePath) {
		if err := os.Remove(h.lockFilePath); err != nil {
			return err
		}
	}

	h.closed = true
	removeFromContainer(h.path)
	return nil
}

func (h *Handler) CloseWithErrors() []error {
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

	if h.tempFp != nil {
		if err := file.Close(h.tempFp); err != nil {
			errs = append(errs, err)
		} else {
			h.tempFp = nil
		}
	}

	if Exists(h.tempFilePath) {
		if err := os.Remove(h.tempFilePath); err != nil {
			errs = append(errs, err)
		}
	}

	if h.lockFileFp != nil {
		if err := file.Close(h.lockFileFp); err != nil {
			errs = append(errs, err)
		} else {
			h.lockFileFp = nil
		}
	}

	if Exists(h.lockFilePath) {
		if err := os.Remove(h.lockFilePath); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (h *Handler) TryCreateLockFileWithTimeout() error {
	var start time.Time

	for {
		if start.IsZero() {
			start = time.Now()
		} else if time.Since(start).Seconds() > WaitTimeout {
			return NewTimeoutError(h.path)
		}

		if err := h.TryCreateLockFile(); err == nil {
			break
		}
		time.Sleep(RetryInterval)
	}
	return nil
}

func (h *Handler) TryCreateLockFile() error {
	if len(h.path) < 1 {
		return NewLockError("filename not specified")
	}
	if h.lockFileFp != nil {
		return NewLockError(fmt.Sprintf("lock file for %s is already created", h.path))
	}

	lockFilePath := LockFilePath(h.path)
	fp, err := file.Create(lockFilePath)
	if err != nil {
		return NewLockError(fmt.Sprintf("unable to create lock file for %q", h.path))
	}

	h.lockFilePath = lockFilePath
	h.lockFileFp = fp
	return nil
}

func (h *Handler) TryCreateTempFile() error {
	if len(h.path) < 1 {
		return NewLockError("filename not specified")
	}
	if h.tempFp != nil {
		return NewLockError(fmt.Sprintf("temporary file for %s is already created", h.path))
	}

	tempFilePath := TempFilePath(h.path)
	fp, err := file.Create(tempFilePath)
	if err != nil {
		return NewLockError(fmt.Sprintf("unable to create temporary file for %q", h.path))
	}

	h.tempFilePath = tempFilePath
	h.tempFp = fp
	return nil
}

func (h *Handler) PrepareToRead() error {
	if !Exists(h.path) {
		return NewIOError(fmt.Sprintf("file %s does not exist", h.path))
	}

	lockFilePath := LockFilePath(h.path)

	var start time.Time

	for {
		if start.IsZero() {
			start = time.Now()
		} else if time.Since(start).Seconds() > WaitTimeout {
			return NewTimeoutError(h.path)
		}

		if _, err := os.Stat(lockFilePath); err != nil {
			break
		}

		time.Sleep(RetryInterval)
	}

	return nil
}
