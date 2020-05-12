package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/csvq/lib/parser"
)

func TestPreparedStatementMap_Prepare(t *testing.T) {
	m := NewPreparedStatementMap()

	expr := parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "stmt"},
		Statement: value.NewString("select 1"),
	}

	expect := GenerateStatementMap([]*PreparedStatement{
		{
			Name:            "stmt",
			StatementString: "select 1",
			Statements: []parser.Statement{
				parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
							Fields: []parser.QueryExpression{
								parser.Field{
									Object: parser.NewIntegerValueFromString("1"),
								},
							},
						},
					},
				},
			},
			HolderNumber: 0,
		},
	})

	err := m.Prepare(TestTx.Flags, expr)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	} else {
		if !SyncMapEqual(m, expect) {
			t.Errorf("result = %v, want %v", m, expect)
		}
	}

	expectErr := "statement stmt is a duplicate"
	err = m.Prepare(TestTx.Flags, expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else {
		if err.Error() != expectErr {
			t.Errorf("error = %q, want error %q", err.Error(), expectErr)
		}
	}

	expr = parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "stmt2"},
		Statement: value.NewString("select from"),
	}
	expectErr = "prepare stmt2 [L:1 C:8] syntax error: unexpected token \"from\""
	err = m.Prepare(TestTx.Flags, expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else {
		if err.Error() != expectErr {
			t.Errorf("error = %q, want error %q", err.Error(), expectErr)
		}
	}
}

func TestPreparedStatementMap_Get(t *testing.T) {
	m := GenerateStatementMap([]*PreparedStatement{
		{
			Name: "stmt",
			Statements: []parser.Statement{
				parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
							Fields: []parser.QueryExpression{
								parser.Field{
									Object: parser.NewIntegerValueFromString("1"),
								},
							},
						},
					},
				},
			},
		},
	})

	name := parser.Identifier{Literal: "stmt"}
	expect := &PreparedStatement{
		Name: "stmt",
		Statements: []parser.Statement{
			parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
						Fields: []parser.QueryExpression{
							parser.Field{
								Object: parser.NewIntegerValueFromString("1"),
							},
						},
					},
				},
			},
		},
	}

	stmt, err := m.Get(name)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	} else {
		if !reflect.DeepEqual(stmt, expect) {
			t.Errorf("result = %v, want %v", stmt, expect)
		}
	}

	name = parser.Identifier{Literal: "notexist"}
	expectErr := "statement notexist does not exist"

	_, err = m.Get(name)
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error %q, want error %q", err.Error(), expectErr)
	}
}

func TestPreparedStatementMap_Dispose(t *testing.T) {
	m := GenerateStatementMap([]*PreparedStatement{
		{
			Name: "stmt",
			Statements: []parser.Statement{
				parser.SelectQuery{
					SelectEntity: parser.SelectEntity{
						SelectClause: parser.SelectClause{
							BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
							Fields: []parser.QueryExpression{
								parser.Field{
									Object: parser.NewIntegerValueFromString("1"),
								},
							},
						},
					},
				},
			},
		},
	})

	expr := parser.DisposeStatement{
		Name: parser.Identifier{Literal: "stmt"},
	}

	expect := NewPreparedStatementMap()

	err := m.Dispose(expr)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	} else {
		if !SyncMapEqual(m, expect) {
			t.Errorf("result = %v, want %v", m, expect)
		}
	}

	expectErr := "statement stmt does not exist"
	err = m.Dispose(expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error %q, want error %q", err.Error(), expectErr)
	}
}

func TestNewPreparedStatement(t *testing.T) {
	expr := parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "stmt"},
		Statement: value.NewString("select 1"),
	}
	expect := &PreparedStatement{
		Name:            "stmt",
		StatementString: "select 1",
		Statements: []parser.Statement{
			parser.SelectQuery{
				SelectEntity: parser.SelectEntity{
					SelectClause: parser.SelectClause{
						BaseExpr: parser.NewBaseExpr(parser.Token{Line: 1, Char: 1, SourceFile: "stmt"}),
						Fields: []parser.QueryExpression{
							parser.Field{
								Object: parser.NewIntegerValueFromString("1"),
							},
						},
					},
				},
			},
		},
		HolderNumber: 0,
	}

	result, err := NewPreparedStatement(TestTx.Flags, expr)
	if err != nil {
		t.Errorf("error %q, want no error", err.Error())
	} else {
		if !reflect.DeepEqual(result, expect) {
			t.Errorf("result = %v, want %v", result, expect)
		}
	}

	expr = parser.StatementPreparation{
		Name:      parser.Identifier{Literal: "stmt"},
		Statement: value.NewString("select from"),
	}
	expectErr := "prepare stmt [L:1 C:8] syntax error: unexpected token \"from\""

	_, err = NewPreparedStatement(TestTx.Flags, expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectErr)
	} else if err.Error() != expectErr {
		t.Errorf("error %q, want error %q", err.Error(), expectErr)
	}
}

func TestNewReplaceValues(t *testing.T) {
	values := []parser.ReplaceValue{
		{Value: parser.NewIntegerValueFromString("1")},
		{Value: parser.NewStringValue("a"), Name: parser.Identifier{Literal: "val"}},
	}
	expect := &ReplaceValues{
		Values: []parser.QueryExpression{
			parser.NewIntegerValueFromString("1"),
			parser.NewStringValue("a"),
		},
		Names: map[string]int{
			"val": 1,
		},
	}

	result := NewReplaceValues(values)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}
}
