package query

import (
	"errors"
	"fmt"
)

type IdentificationError struct {
	Identifier string
	Err        error
}

func (e IdentificationError) Error() string {
	return fmt.Sprintf("identifier = %s: %s", e.Identifier, e.Err)
}

type NotGroupedError struct {
	Function string
	Err      error
}

func (e NotGroupedError) Error() string {
	return fmt.Sprintf("function %s: %s", e.Function, e.Err)
}

var (
	ErrFieldAmbiguous = errors.New("field is ambiguous")
	ErrFieldNotExist  = errors.New("field does not exist")
	ErrNotGrouped     = errors.New("records are not grouped")
)
