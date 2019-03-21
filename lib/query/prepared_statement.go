package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type PreparedStatementMap map[string]*PreparedStatement

func (m PreparedStatementMap) Prepare(filter *Filter, expr parser.StatementPreparation) error {
	uname := strings.ToUpper(expr.Name.Literal)
	if _, ok := m[uname]; ok {
		return NewDuplicateStatementNameError(expr.Name)
	}
	stmt, err := NewPreparedStatement(filter, expr)
	if err != nil {
		return err
	}
	m[uname] = stmt
	return nil
}

func (m PreparedStatementMap) Get(name parser.Identifier) (*PreparedStatement, error) {
	uname := strings.ToUpper(name.Literal)
	if stmt, ok := m[uname]; ok {
		return stmt, nil
	}
	return nil, NewStatementNotExistError(name)
}

func (m PreparedStatementMap) Dispose(expr parser.DisposeStatement) error {
	uname := strings.ToUpper(expr.Name.Literal)
	if _, ok := m[uname]; ok {
		delete(m, uname)
		return nil
	}
	return NewStatementNotExistError(expr.Name)
}

type PreparedStatement struct {
	Name            string
	StatementString string
	Statements      []parser.Statement
	HolderNumber    int
}

func NewPreparedStatement(filter *Filter, expr parser.StatementPreparation) (*PreparedStatement, error) {
	statements, holderNum, err := parser.Parse(expr.Statement.Raw(), expr.Name.Literal, filter.tx.Flags.DatetimeFormat, true)
	if err != nil {
		return nil, NewPreparedStatementSyntaxError(err.(*parser.SyntaxError))
	}

	return &PreparedStatement{
		Name:            expr.Name.Literal,
		StatementString: expr.Statement.Raw(),
		Statements:      statements,
		HolderNumber:    holderNum,
	}, nil
}

type ReplaceValues struct {
	Values []parser.QueryExpression
	Names  map[string]int
}

func NewReplaceValues(replace []parser.ReplaceValue) *ReplaceValues {
	values := make([]parser.QueryExpression, 0, len(replace))
	names := make(map[string]int, len(replace))

	for i := range replace {
		if 0 < len(replace[i].Name.Literal) {
			names[replace[i].Name.Literal] = i
		}
		values = append(values, replace[i].Value)
	}

	return &ReplaceValues{
		Values: values,
		Names:  names,
	}
}
