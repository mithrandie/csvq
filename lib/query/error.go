package query

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/parser"
)

type FieldAmbiguousError struct {
	Field parser.FieldReference
	Err   error
}

func (e FieldAmbiguousError) Error() string {
	return fmt.Sprintf(e.Err.Error(), e.Field.String())
}

func NewFieldAmbiguousError(field parser.FieldReference) error {
	return &FieldAmbiguousError{
		Field: field,
		Err:   ErrFieldAmbiguous,
	}
}

type FieldNotExistError struct {
	Field parser.FieldReference
	Err   error
}

func (e FieldNotExistError) Error() string {
	return fmt.Sprintf(e.Err.Error(), e.Field.String())
}

func NewFieldNotExistError(field parser.FieldReference) error {
	return &FieldNotExistError{
		Field: field,
		Err:   ErrFieldNotExist,
	}
}

type NotGroupedError struct {
	Function string
	Err      error
}

func (e NotGroupedError) Error() string {
	return fmt.Sprintf("aggregate function %s: %s", e.Function, e.Err)
}

func NewNotGroupedErr(funcName string) error {
	return &NotGroupedError{
		Function: funcName,
		Err:      ErrNotGrouped,
	}
}

type UndefinedVariableError struct {
	Name string
	Err  error
}

func (e UndefinedVariableError) Error() string {
	return fmt.Sprintf(e.Err.Error(), e.Name)
}

func NewUndefinedVariableError(name string) error {
	return &UndefinedVariableError{
		Name: name,
		Err:  ErrUndefinedVariable,
	}
}

type RedeclaredVariableError struct {
	Name string
	Err  error
}

func (e RedeclaredVariableError) Error() string {
	return fmt.Sprintf(e.Err.Error(), e.Name)
}

func NewRedeclaredVariableError(name string) error {
	return &RedeclaredVariableError{
		Name: name,
		Err:  ErrRedeclaredVariable,
	}
}

var (
	ErrFieldAmbiguous     = errors.New("field %s is ambiguous")
	ErrFieldNotExist      = errors.New("field %s does not exist")
	ErrNotGrouped         = errors.New("records are not grouped")
	ErrUndefinedVariable  = errors.New("variable %s is undefined")
	ErrRedeclaredVariable = errors.New("variable %s is redeclared")
)
