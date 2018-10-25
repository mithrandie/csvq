package query

import (
	"fmt"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

type Procedure struct {
	Filter           *Filter
	ReturnVal        value.Primary
	MeasurementStart time.Time
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
	return child.Execute(statements)
}

func (proc *Procedure) Execute(statements []parser.Statement) (StatementFlow, error) {
	flow := Terminate

	for _, stmt := range statements {
		f, err := proc.ExecuteStatement(stmt)
		if err != nil {
			return f, err
		}
		if f != Terminate {
			flow = f
			break
		}
	}
	return flow, nil
}

func (proc *Procedure) ExecuteStatement(stmt parser.Statement) (StatementFlow, error) {
	flags := cmd.GetFlags()
	flow := Terminate

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
	case parser.ViewDeclaration:
		err = DeclareView(stmt.(parser.ViewDeclaration), proc.Filter)
	case parser.DisposeView:
		err = proc.Filter.TempViews.Dispose(stmt.(parser.DisposeView).View)
	case parser.FunctionDeclaration:
		err = proc.Filter.Functions.Declare(stmt.(parser.FunctionDeclaration))
	case parser.DisposeFunction:
		err = proc.Filter.Functions.Dispose(stmt.(parser.DisposeFunction).Name)
	case parser.AggregateDeclaration:
		err = proc.Filter.Functions.DeclareAggregate(stmt.(parser.AggregateDeclaration))
	case parser.SelectQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}
		if view, err = Select(stmt.(parser.SelectQuery), proc.Filter); err == nil {
			var viewstr string
			fileInfo := &FileInfo{
				Format:             flags.Format,
				Delimiter:          flags.WriteDelimiter,
				DelimiterPositions: flags.WriteDelimiterPositions,
				Encoding:           flags.WriteEncoding,
				LineBreak:          flags.LineBreak,
				NoHeader:           flags.WithoutHeader,
				PrettyPrint:        flags.PrettyPrint,
			}
			viewstr, err = EncodeView(view, fileInfo)
			if err == nil {
				if 0 < len(flags.OutFile) {
					AddSelectLog(viewstr)
				} else {
					Log(viewstr, false)
				}
			}
		}
		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.InsertQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}
		if view, err = Insert(stmt.(parser.InsertQuery), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          InsertQuery,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedRecords,
				},
			}
			Log(fmt.Sprintf("%s inserted on %q.", FormatCount(view.OperatedRecords, "record"), view.FileInfo.Path), flags.Quiet)

			view.OperatedRecords = 0
		}
		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.UpdateQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}
		if views, err = Update(stmt.(parser.UpdateQuery), proc.Filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          UpdateQuery,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				Log(fmt.Sprintf("%s updated on %q.", FormatCount(v.OperatedRecords, "record"), v.FileInfo.Path), flags.Quiet)

				v.OperatedRecords = 0
			}
		}
		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.DeleteQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}
		if views, err = Delete(stmt.(parser.DeleteQuery), proc.Filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          DeleteQuery,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				Log(fmt.Sprintf("%s deleted on %q.", FormatCount(v.OperatedRecords, "record"), v.FileInfo.Path), flags.Quiet)

				v.OperatedRecords = 0
			}
		}
		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.CreateTable:
		if view, err = CreateTable(stmt.(parser.CreateTable), proc.Filter); err == nil {
			results = []Result{
				{
					Type:     CreateTableQuery,
					FileInfo: view.FileInfo,
				},
			}
			Log(fmt.Sprintf("file %q is created.", view.FileInfo.Path), flags.Quiet)

			view.OperatedRecords = 0
		}
	case parser.AddColumns:
		if view, err = AddColumns(stmt.(parser.AddColumns), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          AddColumnsQuery,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s added on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), flags.Quiet)

			view.OperatedRecords = 0
		}
	case parser.DropColumns:
		if view, err = DropColumns(stmt.(parser.DropColumns), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          DropColumnsQuery,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s dropped on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), flags.Quiet)

			view.OperatedRecords = 0
		}
	case parser.RenameColumn:
		if view, err = RenameColumn(stmt.(parser.RenameColumn), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          RenameColumnQuery,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s renamed on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), flags.Quiet)

			view.OperatedRecords = 0
		}
	case parser.TransactionControl:
		switch stmt.(parser.TransactionControl).Token {
		case parser.COMMIT:
			err = Commit(stmt.(parser.Expression), proc.Filter)
		case parser.ROLLBACK:
			Rollback(proc.Filter)
		}
	case parser.FlowControl:
		switch stmt.(parser.FlowControl).Token {
		case parser.CONTINUE:
			flow = Continue
		case parser.BREAK:
			flow = Break
		}
	case parser.Exit:
		ex := stmt.(parser.Exit)
		code := 0
		if ex.Code != nil {
			code = int(ex.Code.(value.Integer).Raw())
		}
		if 0 < code {
			flow = Error
			err = NewForcedExit(code)
		} else {
			flow = Exit
		}
	case parser.Return:
		var ret value.Primary
		if ret, err = proc.Filter.Evaluate(stmt.(parser.Return).Value); err == nil {
			proc.ReturnVal = ret
			flow = Return
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
	case parser.ShowObjects:
		if printstr, err = ShowObjects(stmt.(parser.ShowObjects), proc.Filter); err == nil {
			Log(printstr, false)
		}
	case parser.ShowFields:
		if printstr, err = ShowFields(stmt.(parser.ShowFields), proc.Filter); err == nil {
			Log(printstr, false)
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
		flow = Error
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
			return Error, err
		}
		if p.Ternary() == ternary.TRUE {
			return proc.ExecuteChild(v.Statements)
		}
	}

	if stmt.Else.Statements != nil {
		return proc.ExecuteChild(stmt.Else.Statements)
	}
	return Terminate, nil
}

func (proc *Procedure) Case(stmt parser.Case) (StatementFlow, error) {
	var val value.Primary
	var err error
	if stmt.Value != nil {
		val, err = proc.Filter.Evaluate(stmt.Value)
		if err != nil {
			return Error, err
		}
	}

	for _, when := range stmt.When {
		var t ternary.Value

		cond, err := proc.Filter.Evaluate(when.Condition)
		if err != nil {
			return Error, err
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
		return Terminate, nil
	}
	return proc.ExecuteChild(stmt.Else.Statements)
}

func (proc *Procedure) While(stmt parser.While) (StatementFlow, error) {
	childProc := proc.NewChildProcedure()

	for {
		childProc.Filter.ResetCurrentScope()
		p, err := proc.Filter.Evaluate(stmt.Condition)
		if err != nil {
			return Error, err
		}
		if p.Ternary() != ternary.TRUE {
			break
		}

		f, err := childProc.Execute(stmt.Statements)
		if err != nil {
			return Error, err
		}

		if f == Break {
			return Terminate, nil
		}
		if f == Exit {
			return Exit, nil
		}
	}
	return Terminate, nil
}

func (proc *Procedure) WhileInCursor(stmt parser.WhileInCursor) (StatementFlow, error) {
	fetchPosition := parser.FetchPosition{
		Position: parser.Token{Token: parser.NEXT},
	}

	childProc := proc.NewChildProcedure()
	for {
		childProc.Filter.ResetCurrentScope()
		if stmt.WithDeclaration {
			assigns := make([]parser.VariableAssignment, len(stmt.Variables))
			for i, v := range stmt.Variables {
				assigns[i] = parser.VariableAssignment{Variable: v}
			}
			decl := parser.VariableDeclaration{Assignments: assigns}
			childProc.Filter.Variables.Declare(decl, childProc.Filter)
		}

		success, err := FetchCursor(stmt.Cursor, fetchPosition, stmt.Variables, childProc.Filter)
		if err != nil {
			return Error, err
		}
		if !success {
			break
		}

		f, err := childProc.Execute(stmt.Statements)
		if err != nil {
			return Error, err
		}

		if f == Break {
			return Terminate, nil
		}
		if f == Exit {
			return Exit, nil
		}
	}

	return Terminate, nil
}

func (proc *Procedure) showExecutionTime() {
	exectime := cmd.HumarizeNumber(fmt.Sprintf("%f", time.Since(proc.MeasurementStart).Seconds()))
	stats := fmt.Sprintf("Query Execution Time: %s seconds", exectime)
	Log(stats, false)
}
