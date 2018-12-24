package query

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/excmd"
	"github.com/mithrandie/csvq/lib/file"
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

var Version string
var ViewCache = make(ViewMap, 10)
var UncommittedViews = NewUncommittedViewMap()

var Formatter = NewStringFormatter()

func ReleaseResources() error {
	if err := ViewCache.Clean(); err != nil {
		return err
	}
	if err := file.UnlockAll(); err != nil {
		return err
	}
	return nil
}

func ReleaseResourcesWithErrors() error {
	var errs []error
	if err := ViewCache.CleanWithErrors(); err != nil {
		errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
	}
	if err := file.UnlockAllWithErrors(); err != nil {
		errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
	}

	if errs != nil {
		return file.NewForcedUnlockError(errs)
	}
	return nil
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

	var printstr string

	switch stmt.(type) {
	case parser.SetFlag:
		err = SetFlag(stmt.(parser.SetFlag), proc.Filter)
	case parser.AddFlagElement:
		err = AddFlagElement(stmt.(parser.AddFlagElement), proc.Filter)
	case parser.RemoveFlagElement:
		err = RemoveFlagElement(stmt.(parser.RemoveFlagElement), proc.Filter)
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

		view, e := Select(stmt.(parser.SelectQuery), proc.Filter)
		if e == nil {
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
				writer = Stdout
			}
			err = EncodeView(writer, view, fileInfo)
			if err == nil {
				writer.Write([]byte(cmd.GetFlags().LineBreak.Value()))
			} else if _, ok := err.(*EmptyResultSetError); ok {
				err = nil
			}
		} else {
			err = e
		}

		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.InsertQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}

		fileInfo, cnt, e := Insert(stmt.(parser.InsertQuery), proc.Filter)
		if e == nil {
			if 0 < cnt {
				UncommittedViews.SetForUpdatedView(fileInfo)
			}
			Log(fmt.Sprintf("%s inserted on %q.", FormatCount(cnt, "record"), fileInfo.Path), flags.Quiet)
		} else {
			err = e
		}

		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.UpdateQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}

		infos, cnts, e := Update(stmt.(parser.UpdateQuery), proc.Filter)
		if e == nil {
			for i, info := range infos {
				if 0 < cnts[i] {
					UncommittedViews.SetForUpdatedView(info)
				}
				Log(fmt.Sprintf("%s updated on %q.", FormatCount(cnts[i], "record"), info.Path), flags.Quiet)
			}
		} else {
			err = e
		}

		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.DeleteQuery:
		if flags.Stats {
			proc.MeasurementStart = time.Now()
		}

		infos, cnts, e := Delete(stmt.(parser.DeleteQuery), proc.Filter)
		if e == nil {
			for i, info := range infos {
				if 0 < cnts[i] {
					UncommittedViews.SetForUpdatedView(info)
				}
				Log(fmt.Sprintf("%s deleted on %q.", FormatCount(cnts[i], "record"), info.Path), flags.Quiet)
			}
		} else {
			err = e
		}

		if flags.Stats {
			proc.showExecutionTime()
		}
	case parser.CreateTable:
		info, e := CreateTable(stmt.(parser.CreateTable), proc.Filter)
		if e == nil {
			UncommittedViews.SetForCreatedView(info)
			Log(fmt.Sprintf("file %q is created.", info.Path), flags.Quiet)
		} else {
			err = e
		}
	case parser.AddColumns:
		info, cnt, e := AddColumns(stmt.(parser.AddColumns), proc.Filter)
		if e == nil {
			UncommittedViews.SetForUpdatedView(info)
			Log(fmt.Sprintf("%s added on %q.", FormatCount(cnt, "field"), info.Path), flags.Quiet)
		} else {
			err = e
		}
	case parser.DropColumns:
		info, cnt, e := DropColumns(stmt.(parser.DropColumns), proc.Filter)
		if e == nil {
			UncommittedViews.SetForUpdatedView(info)
			Log(fmt.Sprintf("%s dropped on %q.", FormatCount(cnt, "field"), info.Path), flags.Quiet)
		} else {
			err = e
		}
	case parser.RenameColumn:
		info, e := RenameColumn(stmt.(parser.RenameColumn), proc.Filter)
		if e == nil {
			UncommittedViews.SetForUpdatedView(info)
			Log(fmt.Sprintf("%s renamed on %q.", FormatCount(1, "field"), info.Path), flags.Quiet)
		} else {
			err = e
		}
	case parser.SetTableAttribute:
		expr := stmt.(parser.SetTableAttribute)
		info, log, e := SetTableAttribute(expr, proc.Filter)
		if e == nil {
			UncommittedViews.SetForUpdatedView(info)
			Log(log, flags.Quiet)
		} else {
			if unchanged, ok := e.(*TableAttributeUnchangedError); ok {
				Log(fmt.Sprintf("Table attributes of %s remain unchanged.", unchanged.Path), flags.Quiet)
			} else {
				err = e
			}
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
	case parser.Echo:
		if printstr, err = Echo(stmt.(parser.Echo), proc.Filter); err == nil {
			Log(printstr, false)
		}
	case parser.Print:
		if printstr, err = Print(stmt.(parser.Print), proc.Filter); err == nil {
			Log(printstr, false)
		}
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
	case parser.Chdir:
		err = Chdir(stmt.(parser.Chdir), proc.Filter)
	case parser.Pwd:
		var dirpath string
		dirpath, err = Pwd(stmt.(parser.Pwd))
		if err == nil {
			Log(dirpath, false)
		}
	case parser.Reload:
		err = Reload(stmt.(parser.Reload))
	case parser.ShowObjects:
		if printstr, err = ShowObjects(stmt.(parser.ShowObjects), proc.Filter); err == nil {
			Log(printstr, false)
		}
	case parser.ShowFields:
		if printstr, err = ShowFields(stmt.(parser.ShowFields), proc.Filter); err == nil {
			Log(printstr, false)
		}
	case parser.Syntax:
		printstr = Syntax(stmt.(parser.Syntax), proc.Filter)
		Log(printstr, false)
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
	case parser.ExternalCommand:
		err = proc.ExecExternalCommand(stmt.(parser.ExternalCommand))
	default:
		if expr, ok := stmt.(parser.QueryExpression); ok {
			_, err = proc.Filter.Evaluate(expr)
		}
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

func (proc *Procedure) ExecExternalCommand(stmt parser.ExternalCommand) error {
	splitter := new(excmd.ArgsSplitter).Init(stmt.Command)
	var argStrs = make([]string, 0, 8)
	for splitter.Scan() {
		argStrs = append(argStrs, splitter.Text())
	}
	err := splitter.Err()
	if err != nil {
		return NewExternalCommandError(stmt, err.Error())
	}

	args := make([]string, 0, len(argStrs))
	for _, argStr := range argStrs {
		arg, err := proc.Filter.EvaluateEmbeddedString(argStr)
		if err != nil {
			if appErr, ok := err.(AppError); ok {
				err = NewExternalCommandError(stmt, appErr.ErrorMessage())
			} else {
				err = NewExternalCommandError(stmt, err.Error())
			}
			return err
		}
		args = append(args, arg)
	}

	if len(args) < 1 {
		return nil
	}

	c := exec.Command(args[0], args[1:]...)
	c.Stdin = Stdin
	c.Stdout = Stdout
	c.Stderr = Stderr

	err = c.Run()
	if err != nil {
		err = NewExternalCommandError(stmt, err.Error())
	}
	return err
}

func (proc *Procedure) showExecutionTime() {
	palette, _ := cmd.GetPalette()
	exectime := cmd.FormatNumber(time.Since(proc.MeasurementStart).Seconds(), 6, ".", ",", "")
	stats := fmt.Sprintf(palette.Render(cmd.LableEffect, "Query Execution Time: ")+"%s seconds", exectime)
	Log(stats, false)
}
