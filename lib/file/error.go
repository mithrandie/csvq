package file

import (
	"fmt"
	"strings"

	"github.com/mithrandie/go-file/v2"
)

func ParseError(err error) error {
	switch err.(type) {
	case *file.IOError:
		return NewIOError(err.Error())
	case *file.LockError:
		return NewLockError(err.Error())
	case *file.TimeoutError:
		return &TimeoutError{
			message: err.Error(),
		}
	case *file.ContextIsDone:
		return NewContextIsDone(err.Error())
	default:
		return err
	}
}

type IOError struct {
	message string
}

func NewIOError(message string) error {
	return &IOError{
		message: message,
	}
}

func (e IOError) Error() string {
	return e.message
}

type LockError struct {
	message string
}

func NewLockError(message string) error {
	return &LockError{
		message: message,
	}
}

func (e LockError) Error() string {
	return e.message
}

type TimeoutError struct {
	message string
}

func NewTimeoutError(path string) error {
	return &TimeoutError{
		message: fmt.Sprintf("file %s: lock waiting time exceeded", path),
	}
}

func (e TimeoutError) Error() string {
	return e.message
}

type ContextIsDone struct {
	message string
}

func NewContextIsDone(message string) error {
	return &ContextIsDone{
		message: message,
	}
}

func (e ContextIsDone) Error() string {
	return e.message
}

type ForcedUnlockError struct {
	Errors []error
}

func NewForcedUnlockError(errs []error) error {
	if errs == nil {
		return nil
	}

	return &ForcedUnlockError{
		Errors: errs,
	}
}

func (e ForcedUnlockError) Error() string {
	list := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		list = append(list, err.Error())
	}
	return strings.Join(list, "\n  ")
}

type CompositeError struct {
	message string
}

func NewCompositeError(err1 error, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}

	return &CompositeError{
		message: err1.Error() + "\n  " + err2.Error(),
	}
}

func (e CompositeError) Error() string {
	return e.message
}
