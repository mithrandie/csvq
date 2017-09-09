package query

import (
	"fmt"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
	"github.com/mithrandie/csvq/lib/value"
)

type Procedure struct {
	Filter    *Filter
	ReturnVal value.Primary
}

func NewProcedure() *Procedure {
	return &Procedure{
		Filter: NewEmptyFilter(),
	}
}

func (proc *Procedure) NewChildProcedure() *Procedure {
	return &Procedure{
		Filter: proc.Filter.CreateChildScope(),
	}
}

func (proc *Procedure) ExecuteChild(statements []parser.Statement) (StatementFlow, error) {
	child := proc.NewChildProcedure()
	f, err := child.Execute(statements)
	return f, err
}

func (proc *Procedure) Execute(statements []parser.Statement) (StatementFlow, error) {
	flow := TERMINATE

	for _, stmt := range statements {
		f, err := proc.ExecuteStatement(stmt)
		if err != nil {
			return f, err
		}
		if f != TERMINATE {
			flow = f
			break
		}
	}
	return flow, nil
}

func (proc *Procedure) ExecuteStatement(stmt parser.Statement) (StatementFlow, error) {
	flow := TERMINATE

	var err error

	var results []Result
	var view *View
	var views []*View
	var printstr string

	switch stmt.(type) {
	case parser.SetFlag:
		err = SetFlag(stmt.(parser.SetFlag))
	case parser.ShowFlag:
		if printstr, err = ShowFlag(stmt.(parser.ShowFlag)); err == nil {
			Log(printstr, false)
		}
	case parser.VariableDeclaration:
		err = proc.Filter.Variables.Declare(stmt.(parser.VariableDeclaration), proc.Filter)
	case parser.VariableSubstitution:
		_, err = proc.Filter.Evaluate(stmt.(parser.QueryExpression))
	case parser.DisposeVariable:
		err = proc.Filter.Variables.Dispose(stmt.(parser.DisposeVariable).Variable)
	case parser.CursorDeclaration:
		err = proc.Filter.Cursors.Declare(stmt.(parser.CursorDeclaration))
	case parser.OpenCursor:
		err = proc.Filter.Cursors.Open(stmt.(parser.OpenCursor).Cursor, proc.Filter)
	case parser.CloseCursor:
		err = proc.Filter.Cursors.Close(stmt.(parser.CloseCursor).Cursor)
	case parser.DisposeCursor:
		err = proc.Filter.Cursors.Dispose(stmt.(parser.DisposeCursor).Cursor)
	case parser.FetchCursor:
		fetch := stmt.(parser.FetchCursor)
		_, err = FetchCursor(fetch.Cursor, fetch.Position, fetch.Variables, proc.Filter)
	case parser.TableDeclaration:
		err = DeclareTable(stmt.(parser.TableDeclaration), proc.Filter)
	case parser.DisposeTable:
		err = proc.Filter.TempViews.Dispose(stmt.(parser.DisposeTable).Table)
	case parser.FunctionDeclaration:
		err = proc.Filter.Functions.Declare(stmt.(parser.FunctionDeclaration))
	case parser.AggregateDeclaration:
		err = proc.Filter.Functions.DeclareAggregate(stmt.(parser.AggregateDeclaration))
	case parser.SelectQuery:
		if view, err = Select(stmt.(parser.SelectQuery), proc.Filter); err == nil {
			flags := cmd.GetFlags()
			var viewstr string
			var lineBreak = cmd.LF
			if 0 < len(flags.OutFile) {
				lineBreak = flags.LineBreak
			}
			viewstr, err = EncodeView(view, flags.Format, flags.WriteDelimiter, flags.WithoutHeader, flags.WriteEncoding, lineBreak)
			if err == nil {
				if 0 < len(flags.OutFile) {
					AddSelectLog(viewstr)
				} else {
					Log(viewstr, false)
				}
			}
		}
	case parser.InsertQuery:
		if view, err = Insert(stmt.(parser.InsertQuery), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          INSERT,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedRecords,
				},
			}
			Log(fmt.Sprintf("%s inserted on %q.", FormatCount(view.OperatedRecords, "record"), view.FileInfo.Path), cmd.GetFlags().Quiet)

			view.OperatedRecords = 0
		}
	case parser.UpdateQuery:
		if views, err = Update(stmt.(parser.UpdateQuery), proc.Filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          UPDATE,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				Log(fmt.Sprintf("%s updated on %q.", FormatCount(v.OperatedRecords, "record"), v.FileInfo.Path), cmd.GetFlags().Quiet)

				v.OperatedRecords = 0
			}
		}
	case parser.DeleteQuery:
		if views, err = Delete(stmt.(parser.DeleteQuery), proc.Filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          DELETE,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				Log(fmt.Sprintf("%s deleted on %q.", FormatCount(v.OperatedRecords, "record"), v.FileInfo.Path), cmd.GetFlags().Quiet)

				v.OperatedRecords = 0
			}
		}
	case parser.CreateTable:
		if view, err = CreateTable(stmt.(parser.CreateTable), proc.Filter); err == nil {
			results = []Result{
				{
					Type:     CREATE_TABLE,
					FileInfo: view.FileInfo,
				},
			}
			Log(fmt.Sprintf("file %q is created.", view.FileInfo.Path), cmd.GetFlags().Quiet)

			view.OperatedRecords = 0
		}
	case parser.AddColumns:
		if view, err = AddColumns(stmt.(parser.AddColumns), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          ADD_COLUMNS,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s added on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), cmd.GetFlags().Quiet)

			view.OperatedRecords = 0
		}
	case parser.DropColumns:
		if view, err = DropColumns(stmt.(parser.DropColumns), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          DROP_COLUMNS,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s dropped on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), cmd.GetFlags().Quiet)

			view.OperatedRecords = 0
		}
	case parser.RenameColumn:
		if view, err = RenameColumn(stmt.(parser.RenameColumn), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          RENAME_COLUMN,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s renamed on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), cmd.GetFlags().Quiet)

			view.OperatedRecords = 0
		}
	case parser.TransactionControl:
		switch stmt.(parser.TransactionControl).Token {
		case parser.COMMIT:
			err = proc.Commit(stmt.(parser.Expression))
		case parser.ROLLBACK:
			proc.Rollback()
		}
	case parser.FlowControl:
		switch stmt.(parser.FlowControl).Token {
		case parser.CONTINUE:
			flow = CONTINUE
		case parser.BREAK:
			flow = BREAK
		}
	case parser.Exit:
		ex := stmt.(parser.Exit)
		code := 0
		if ex.Code != nil {
			code = int(ex.Code.(value.Integer).Raw())
		}
		if 0 < code {
			flow = ERROR
			err = NewExit(code)
		} else {
			flow = EXIT
		}
	case parser.Return:
		var ret value.Primary
		if ret, err = proc.Filter.Evaluate(stmt.(parser.Return).Value); err == nil {
			proc.ReturnVal = ret
			flow = RETURN
		}
	case parser.If:
		flow, err = proc.IfStmt(stmt.(parser.If))
	case parser.Case:
		flow, err = proc.Case(stmt.(parser.Case))
	case parser.While:
		flow, err = proc.While(stmt.(parser.While))
	case parser.WhileInCursor:
		flow, err = proc.WhileInCursor(stmt.(parser.WhileInCursor))
	case parser.Print:
		if printstr, err = Print(stmt.(parser.Print), proc.Filter); err == nil {
			Log(printstr, false)
		}
	case parser.Function:
		_, err = proc.Filter.Evaluate(stmt.(parser.Function))
	case parser.Printf:
		if printstr, err = Printf(stmt.(parser.Printf), proc.Filter); err == nil {
			Log(printstr, false)
		}
	case parser.Source:
		var externalStatements []parser.Statement
		source := stmt.(parser.Source)
		if externalStatements, err = Source(source, proc.Filter); err == nil {
			flow, err = proc.Execute(externalStatements)
		}
	case parser.Trigger:
		trigger := stmt.(parser.Trigger)
		switch trigger.Token {
		case parser.ERROR:
			var message string
			if trigger.Message != nil {
				if pt, ok := trigger.Message.(parser.PrimitiveType); ok && trigger.Code == nil && pt.IsInteger() {
					trigger.Code = pt.Value
				} else {
					var p value.Primary
					if p, err = proc.Filter.Evaluate(trigger.Message); err == nil {
						if s := value.ToString(p); !value.IsNull(s) {
							message = s.(value.String).Raw()
						}
					}
				}
			}
			if err == nil {
				err = NewUserTriggeredError(trigger, message)
			}
		}
	}

	if results != nil {
		Results = append(Results, results...)
	}

	if err != nil {
		flow = ERROR
	}
	return flow, err
}

func (proc *Procedure) IfStmt(stmt parser.If) (StatementFlow, error) {
	stmts := make([]parser.ElseIf, 0, len(stmt.ElseIf)+1)
	stmts = append(stmts, parser.ElseIf{
		Condition:  stmt.Condition,
		Statements: stmt.Statements,
	})
	for _, v := range stmt.ElseIf {
		stmts = append(stmts, v)
	}

	for _, v := range stmts {
		p, err := proc.Filter.Evaluate(v.Condition)
		if err != nil {
			return ERROR, err
		}
		if p.Ternary() == ternary.TRUE {
			return proc.ExecuteChild(v.Statements)
		}
	}

	if stmt.Else.Statements != nil {
		return proc.ExecuteChild(stmt.Else.Statements)
	}
	return TERMINATE, nil
}

func (proc *Procedure) Case(stmt parser.Case) (StatementFlow, error) {
	var val value.Primary
	var err error
	if stmt.Value != nil {
		val, err = proc.Filter.Evaluate(stmt.Value)
		if err != nil {
			return ERROR, err
		}
	}

	for _, when := range stmt.When {
		var t ternary.Value

		cond, err := proc.Filter.Evaluate(when.Condition)
		if err != nil {
			return ERROR, err
		}

		if val == nil {
			t = cond.Ternary()
		} else {
			t = value.Equal(val, cond)
		}

		if t == ternary.TRUE {
			return proc.ExecuteChild(when.Statements)
		}
	}

	if stmt.Else.Statements == nil {
		return TERMINATE, nil
	}
	return proc.ExecuteChild(stmt.Else.Statements)
}

func (proc *Procedure) While(stmt parser.While) (StatementFlow, error) {
	for {
		p, err := proc.Filter.Evaluate(stmt.Condition)
		if err != nil {
			return ERROR, err
		}
		if p.Ternary() != ternary.TRUE {
			break
		}
		f, err := proc.ExecuteChild(stmt.Statements)
		if err != nil {
			return ERROR, err
		}

		if f == BREAK {
			return TERMINATE, nil
		}
		if f == EXIT {
			return EXIT, nil
		}
	}
	return TERMINATE, nil
}

func (proc *Procedure) WhileInCursor(stmt parser.WhileInCursor) (StatementFlow, error) {
	fetchPosition := parser.FetchPosition{
		Position: parser.Token{Token: parser.NEXT},
	}
	for {
		success, err := FetchCursor(stmt.Cursor, fetchPosition, stmt.Variables, proc.Filter)
		if err != nil {
			return ERROR, err
		}
		if !success {
			break
		}

		f, err := proc.ExecuteChild(stmt.Statements)
		if err != nil {
			return ERROR, err
		}

		if f == BREAK {
			return TERMINATE, nil
		}
		if f == EXIT {
			return EXIT, nil
		}
	}

	return TERMINATE, nil
}

func (proc *Procedure) Commit(expr parser.Expression) error {
	var createFiles = map[string]*FileInfo{}
	var updateFiles = map[string]*FileInfo{}

	for _, result := range Results {
		if result.FileInfo != nil {
			switch result.Type {
			case CREATE_TABLE:
				createFiles[result.FileInfo.Path] = result.FileInfo
			default:
				if !result.FileInfo.IsTemporary && 0 < result.OperatedCount {
					if _, ok := createFiles[result.FileInfo.Path]; !ok {
						if _, ok := updateFiles[result.FileInfo.Path]; !ok {
							updateFiles[result.FileInfo.Path] = result.FileInfo
						}
					}
				}
			}
		}
	}

	if 0 < len(createFiles) {
		for filename, fileinfo := range createFiles {
			view, _ := ViewCache.Get(parser.Identifier{Literal: filename})
			viewstr, err := EncodeView(view, cmd.CSV, fileinfo.Delimiter, false, fileinfo.Encoding, fileinfo.LineBreak)
			if err != nil {
				return err
			}

			if err = cmd.CreateFile(filename, viewstr); err != nil {
				if expr == nil {
					return NewAutoCommitError(err.Error())
				}
				return NewWriteFileError(expr, err.Error())
			}
			Log(fmt.Sprintf("Commit: file %q is created.", filename), cmd.GetFlags().Quiet)
		}
	}

	if 0 < len(updateFiles) {
		for filename, fileinfo := range updateFiles {
			view, _ := ViewCache.Get(parser.Identifier{Literal: filename})
			viewstr, err := EncodeView(view, cmd.CSV, fileinfo.Delimiter, fileinfo.NoHeader, fileinfo.Encoding, fileinfo.LineBreak)
			if err != nil {
				return err
			}

			if err = cmd.UpdateFile(fileinfo.File, viewstr); err != nil {
				if expr == nil {
					return NewAutoCommitError(err.Error())
				}
				return NewWriteFileError(expr, err.Error())
			}
			Log(fmt.Sprintf("Commit: file %q is updated.", filename), cmd.GetFlags().Quiet)
		}
	}

	Results = []Result{}
	ReleaseResources()
	if expr != nil {
		proc.Filter.TempViews.Store()
	}

	return nil
}

func (proc *Procedure) Rollback() {
	Results = []Result{}
	ReleaseResources()
	proc.Filter.TempViews.Restore()

	Log("Rolled back.", cmd.GetFlags().Quiet)
	return
}
