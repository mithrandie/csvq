package query

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/file"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

type StatementFlow int

const (
	Terminate StatementFlow = iota
	Error
	Exit
	Break
	Continue
	Return
)

type OperationType int

const (
	InsertQuery OperationType = iota
	UpdateQuery
	DeleteQuery
	CreateTableQuery
	AddColumnsQuery
	DropColumnsQuery
	RenameColumnQuery
)

type ExecResult struct {
	Type          OperationType
	FileInfo      *FileInfo
	OperatedCount int
}

var ViewCache = make(ViewMap, 10)
var ExecResults = make([]ExecResult, 0, 10)
var OutFile *os.File
var SelectResult = new(bytes.Buffer)

func ReleaseResources() error {
	if err := ViewCache.Clean(); err != nil {
		return err
	}
	if err := file.UnlockAll(); err != nil {
		return err
	}
	return nil
}

func ReleaseResourcesWithErrors() []error {
	var errs []error
	if es := ViewCache.CleanWithErrors(); es != nil {
		errs = append(errs, es...)
	}
	if es := file.UnlockAllWithErrors(); es != nil {
		errs = append(errs, es...)
	}
	return errs
}

func Log(log string, quiet bool) {
	if !quiet {
		cmd.WriteToStdout(log + "\n")
	}
}

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

	var results []ExecResult
	var view *View
	var views []*View
	var printstr string

	switch stmt.(type) {
	case parser.SetFlag:
		err = SetFlag(stmt.(parser.SetFlag), proc.Filter)
	case parser.ShowFlag:
		if printstr, err = ShowFlag(stmt.(parser.ShowFlag)); err == nil {
			Log(printstr, false)
		}
	case parser.VariableDeclaration:
		err = proc.Filter.Variables.Declare(stmt.(parser.VariableDeclaration), proc.Filter)
	case parser.VariableSubstitution:
		_, err = proc.Filter.Evaluate(stmt.(parser.QueryExpression))
	case parser.SetEnvVar:
		err = SetEnvVar(stmt.(parser.SetEnvVar), proc.Filter)
	case parser.UnsetEnvVar:
		err = UnsetEnvVar(stmt.(parser.UnsetEnvVar))
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
			fileInfo := &FileInfo{
				Format:             flags.Format,
				Delimiter:          flags.WriteDelimiter,
				DelimiterPositions: flags.WriteDelimiterPositions,
				Encoding:           flags.WriteEncoding,
				LineBreak:          flags.LineBreak,
				NoHeader:           flags.WithoutHeader,
				EncloseAll:         flags.EncloseAll,
				PrettyPrint:        flags.PrettyPrint,
			}

			var writer io.Writer
			if OutFile != nil {
				writer = OutFile
			} else {
				SelectResult.Reset()
				writer = SelectResult
			}
			err = EncodeView(writer, view, fileInfo)
			if err == nil {
				if OutFile != nil {
					OutFile.WriteString(cmd.GetFlags().LineBreak.Value())
				} else {
					Log(SelectResult.String(), false)
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
			results = []ExecResult{
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
			results = make([]ExecResult, len(views))
			for i, v := range views {
				results[i] = ExecResult{
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
			results = make([]ExecResult, len(views))
			for i, v := range views {
				results[i] = ExecResult{
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
			results = []ExecResult{
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
			results = []ExecResult{
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
			results = []ExecResult{
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
			results = []ExecResult{
				{
					Type:          RenameColumnQuery,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s renamed on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path), flags.Quiet)

			view.OperatedRecords = 0
		}
	case parser.SetTableAttribute:
		expr := stmt.(parser.SetTableAttribute)
		if printstr, err = SetTableAttribute(expr, proc.Filter); err == nil {
			Log(printstr, flags.Quiet)
		}
	case parser.TransactionControl:
		switch stmt.(parser.TransactionControl).Token {
		case parser.COMMIT:
			err = Commit(stmt.(parser.Expression), proc.Filter)
		case parser.ROLLBACK:
			err = Rollback(stmt.(parser.Expression), proc.Filter)
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
		if externalStatements, err = Source(stmt.(parser.Source), proc.Filter); err == nil {
			flow, err = proc.Execute(externalStatements)
		}
	case parser.Execute:
		var externalStatements []parser.Statement
		if externalStatements, err = ParseExecuteStatements(stmt.(parser.Execute), proc.Filter); err == nil {
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
		switch strings.ToUpper(trigger.Event.Literal) {
		case "ERROR":
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
		default:
			err = NewInvalidEventNameError(trigger.Event)
		}
	}

	if results != nil {
		ExecResults = append(ExecResults, results...)
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
	palette, _ := cmd.GetPalette()
	exectime := cmd.FormatNumber(time.Since(proc.MeasurementStart).Seconds(), 6, ".", ",", "")
	stats := fmt.Sprintf(palette.Render(cmd.LableEffect, "Query Execution Time: ")+"%s seconds", exectime)
	Log(stats, false)
}
