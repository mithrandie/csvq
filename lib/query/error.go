package query

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/parser"
)

type IdentificationError struct {
	Field parser.FieldReference
	Err   error
}

func (e IdentificationError) Error() string {
	return fmt.Sprintf(e.Err.Error(), e.Field.String())
}

type NotGroupedError struct {
	Function string
	Err      error
}

func (e NotGroupedError) Error() string {
	return fmt.Sprintf("function %s: %s", e.Function, e.Err)
}

var (
	ErrFieldAmbiguous = errors.New("field %s is ambiguous")
	ErrFieldNotExist  = errors.New("field %s does not exist")
	ErrNotGrouped     = errors.New("records are not grouped")
)
