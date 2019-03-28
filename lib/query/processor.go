package query

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/excmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

type StatementFlow int

const (
	Terminate StatementFlow = iota
	TerminateWithError
	Exit
	Break
	Continue
	Return
)

const StoringResultsContextKey = "store_query_results"
const StatementReplaceValuesContextKey = "statement_replace_values"

func ContextForStoringResults(ctx context.Context) context.Context {
	return context.WithValue(ctx, StoringResultsContextKey, true)
}

func ContextForPreparedStatement(ctx context.Context, values *ReplaceValues) context.Context {
	return context.WithValue(ctx, StatementReplaceValuesContextKey, values)
}

type Processor struct {
	Tx     *Transaction
	Filter *Filter

	storeResults bool

	returnVal        value.Primary
	measurementStart time.Time
}

func NewProcessor(tx *Transaction) *Processor {
	return NewProcessorWithFilter(tx, NewFilter(tx))
}

func NewProcessorWithFilter(tx *Transaction, filter *Filter) *Processor {
	return &Processor{
		Tx:     tx,
		Filter: filter,
	}
}

func (proc *Processor) NewChildProcessor() *Processor {
	return &Processor{
		Tx:     proc.Tx,
		Filter: proc.Filter.CreateChildScope(),
	}
}

func (proc *Processor) Execute(ctx context.Context, statements []parser.Statement) (StatementFlow, error) {
	if v := ctx.Value(StoringResultsContextKey); v != nil {
		if b, ok := v.(bool); ok && b {
			proc.storeResults = true
		}
	}

	proc.Tx.SelectedViews = nil
	proc.Tx.AffectedRows = 0

	flow, err := proc.execute(ctx, statements)
	if err == nil && flow == Terminate && proc.Tx.AutoCommit {
		err = proc.AutoCommit()
	}
	return flow, err
}

func (proc *Processor) execute(ctx context.Context, statements []parser.Statement) (StatementFlow, error) {
	flow := Terminate

	for _, stmt := range statements {
		f, err := proc.ExecuteStatement(ctx, stmt)
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

func (proc *Processor) executeChild(ctx context.Context, statements []parser.Statement) (StatementFlow, error) {
	child := proc.NewChildProcessor()
	flow, err := child.execute(ctx, statements)
	if child.returnVal != nil {
		proc.returnVal = child.returnVal
	}
	return flow, err
}

func (proc *Processor) ExecuteStatement(ctx context.Context, stmt parser.Statement) (StatementFlow, error) {
	if ctx.Err() != nil {
		return TerminateWithError, NewContextIsDone(ctx.Err().Error())
	}

	flow := Terminate

	var err error

	var printstr string

	switch stmt.(type) {
	case parser.SetFlag:
		err = SetFlag(ctx, proc.Filter, stmt.(parser.SetFlag))
	case parser.AddFlagElement:
		err = AddFlagElement(ctx, proc.Filter, stmt.(parser.AddFlagElement))
	case parser.RemoveFlagElement:
		err = RemoveFlagElement(ctx, proc.Filter, stmt.(parser.RemoveFlagElement))
	case parser.ShowFlag:
		if printstr, err = ShowFlag(proc.Tx.Flags, stmt.(parser.ShowFlag)); err == nil {
			proc.Log(printstr, false)
		}
	case parser.VariableDeclaration:
		err = proc.Filter.variables.Declare(ctx, proc.Filter, stmt.(parser.VariableDeclaration))
	case parser.VariableSubstitution:
		_, err = proc.Filter.Evaluate(ctx, stmt.(parser.QueryExpression))
	case parser.SetEnvVar:
		err = SetEnvVar(ctx, proc.Filter, stmt.(parser.SetEnvVar))
	case parser.UnsetEnvVar:
		err = UnsetEnvVar(stmt.(parser.UnsetEnvVar))
	case parser.DisposeVariable:
		err = proc.Filter.variables.Dispose(stmt.(parser.DisposeVariable).Variable)
	case parser.CursorDeclaration:
		err = proc.Filter.cursors.Declare(stmt.(parser.CursorDeclaration))
	case parser.OpenCursor:
		openCur := stmt.(parser.OpenCursor)
		err = proc.Filter.cursors.Open(ctx, proc.Filter, openCur.Cursor, openCur.Values)
	case parser.CloseCursor:
		err = proc.Filter.cursors.Close(stmt.(parser.CloseCursor).Cursor)
	case parser.DisposeCursor:
		err = proc.Filter.cursors.Dispose(stmt.(parser.DisposeCursor).Cursor)
	case parser.FetchCursor:
		fetch := stmt.(parser.FetchCursor)
		_, err = FetchCursor(ctx, proc.Filter, fetch.Cursor, fetch.Position, fetch.Variables)
	case parser.ViewDeclaration:
		err = DeclareView(ctx, proc.Filter, stmt.(parser.ViewDeclaration))
	case parser.DisposeView:
		err = proc.Filter.tempViews.Dispose(stmt.(parser.DisposeView).View)
	case parser.FunctionDeclaration:
		err = proc.Filter.functions.Declare(stmt.(parser.FunctionDeclaration))
	case parser.DisposeFunction:
		err = proc.Filter.functions.Dispose(stmt.(parser.DisposeFunction).Name)
	case parser.AggregateDeclaration:
		err = proc.Filter.functions.DeclareAggregate(stmt.(parser.AggregateDeclaration))
	case parser.StatementPreparation:
		err = proc.Tx.PreparedStatements.Prepare(proc.Filter, stmt.(parser.StatementPreparation))
	case parser.ExecuteStatement:
		execStmt := stmt.(parser.ExecuteStatement)
		prepared, e := proc.Tx.PreparedStatements.Get(execStmt.Name)
		if e != nil {
			err = e
		} else {
			flow, err = proc.execute(ContextForPreparedStatement(ctx, NewReplaceValues(execStmt.Values)), prepared.Statements)
		}
	case parser.DisposeStatement:
		err = proc.Tx.PreparedStatements.Dispose(stmt.(parser.DisposeStatement))
	case parser.SelectQuery:
		if proc.Tx.Flags.Stats {
			proc.measurementStart = time.Now()
		}

		view, e := Select(ctx, proc.Filter, stmt.(parser.SelectQuery))
		if e == nil {
			if proc.storeResults {
				proc.Tx.SelectedViews = append(proc.Tx.SelectedViews, view)

			} else {
				fileInfo := &FileInfo{
					Format:             proc.Tx.Flags.Format,
					Delimiter:          proc.Tx.Flags.WriteDelimiter,
					DelimiterPositions: proc.Tx.Flags.WriteDelimiterPositions,
					Encoding:           proc.Tx.Flags.WriteEncoding,
					LineBreak:          proc.Tx.Flags.LineBreak,
					NoHeader:           proc.Tx.Flags.WithoutHeader,
					EncloseAll:         proc.Tx.Flags.EncloseAll,
					PrettyPrint:        proc.Tx.Flags.PrettyPrint,
					SingleLine:         proc.Tx.Flags.WriteAsSingleLine,
				}

				var writer io.Writer
				if proc.Tx.Session.OutFile != nil {
					writer = proc.Tx.Session.OutFile
				} else {
					writer = proc.Tx.Session.Stdout
				}
				warnmsg, e := EncodeView(writer, view, fileInfo, proc.Tx.Flags)

				if e != nil {
					if _, ok := e.(*EmptyResultSetError); ok {
						if 0 < len(warnmsg) {
							proc.LogWarn(warnmsg, proc.Tx.Flags.Quiet)
						}
					} else {
						err = e
					}
				} else if !(proc.Tx.Session.OutFile != nil && fileInfo.Format == cmd.FIXED && fileInfo.SingleLine) {
					_, err = writer.Write([]byte(proc.Tx.Flags.LineBreak.Value()))
				}
			}
		} else {
			err = e
		}

		if proc.Tx.Flags.Stats {
			proc.showExecutionTime()
		}
	case parser.InsertQuery:
		if proc.Tx.Flags.Stats {
			proc.measurementStart = time.Now()
		}

		fileInfo, cnt, e := Insert(ctx, proc.Filter, stmt.(parser.InsertQuery))
		if e == nil {
			if 0 < cnt {
				proc.Tx.uncommittedViews.SetForUpdatedView(fileInfo)
			}
			proc.Log(fmt.Sprintf("%s inserted on %q.", FormatCount(cnt, "record"), fileInfo.Path), proc.Tx.Flags.Quiet)
			if proc.storeResults {
				proc.Tx.AffectedRows = cnt
			}
		} else {
			err = e
		}

		if proc.Tx.Flags.Stats {
			proc.showExecutionTime()
		}
	case parser.UpdateQuery:
		if proc.Tx.Flags.Stats {
			proc.measurementStart = time.Now()
		}

		infos, cnts, e := Update(ctx, proc.Filter, stmt.(parser.UpdateQuery))
		if e == nil {
			cntTotal := 0
			for i, info := range infos {
				if 0 < cnts[i] {
					proc.Tx.uncommittedViews.SetForUpdatedView(info)
					cntTotal += cnts[i]
				}
				proc.Log(fmt.Sprintf("%s updated on %q.", FormatCount(cnts[i], "record"), info.Path), proc.Tx.Flags.Quiet)
			}
			if proc.storeResults {
				proc.Tx.AffectedRows = cntTotal
			}
		} else {
			err = e
		}

		if proc.Tx.Flags.Stats {
			proc.showExecutionTime()
		}
	case parser.DeleteQuery:
		if proc.Tx.Flags.Stats {
			proc.measurementStart = time.Now()
		}

		infos, cnts, e := Delete(ctx, proc.Filter, stmt.(parser.DeleteQuery))
		if e == nil {
			cntTotal := 0
			for i, info := range infos {
				if 0 < cnts[i] {
					proc.Tx.uncommittedViews.SetForUpdatedView(info)
					cntTotal += cnts[i]
				}
				proc.Log(fmt.Sprintf("%s deleted on %q.", FormatCount(cnts[i], "record"), info.Path), proc.Tx.Flags.Quiet)
			}
			if proc.storeResults {
				proc.Tx.AffectedRows = cntTotal
			}
		} else {
			err = e
		}

		if proc.Tx.Flags.Stats {
			proc.showExecutionTime()
		}
	case parser.CreateTable:
		info, e := CreateTable(ctx, proc.Filter, stmt.(parser.CreateTable))
		if e == nil {
			proc.Tx.uncommittedViews.SetForCreatedView(info)
			proc.Log(fmt.Sprintf("file %q is created.", info.Path), proc.Tx.Flags.Quiet)
		} else {
			err = e
		}
	case parser.AddColumns:
		info, cnt, e := AddColumns(ctx, proc.Filter, stmt.(parser.AddColumns))
		if e == nil {
			proc.Tx.uncommittedViews.SetForUpdatedView(info)
			proc.Log(fmt.Sprintf("%s added on %q.", FormatCount(cnt, "field"), info.Path), proc.Tx.Flags.Quiet)
		} else {
			err = e
		}
	case parser.DropColumns:
		info, cnt, e := DropColumns(ctx, proc.Filter, stmt.(parser.DropColumns))
		if e == nil {
			proc.Tx.uncommittedViews.SetForUpdatedView(info)
			proc.Log(fmt.Sprintf("%s dropped on %q.", FormatCount(cnt, "field"), info.Path), proc.Tx.Flags.Quiet)
		} else {
			err = e
		}
	case parser.RenameColumn:
		info, e := RenameColumn(ctx, proc.Filter, stmt.(parser.RenameColumn))
		if e == nil {
			proc.Tx.uncommittedViews.SetForUpdatedView(info)
			proc.Log(fmt.Sprintf("%s renamed on %q.", FormatCount(1, "field"), info.Path), proc.Tx.Flags.Quiet)
		} else {
			err = e
		}
	case parser.SetTableAttribute:
		expr := stmt.(parser.SetTableAttribute)
		info, log, e := SetTableAttribute(ctx, proc.Filter, expr)
		if e == nil {
			proc.Tx.uncommittedViews.SetForUpdatedView(info)
			proc.Log(log, proc.Tx.Flags.Quiet)
		} else {
			if unchanged, ok := e.(*TableAttributeUnchangedError); ok {
				proc.Log(fmt.Sprintf("Table attributes of %s remain unchanged.", unchanged.Path), proc.Tx.Flags.Quiet)
			} else {
				err = e
			}
		}
	case parser.TransactionControl:
		switch stmt.(parser.TransactionControl).Token {
		case parser.COMMIT:
			err = proc.Commit(stmt.(parser.Expression))
		case parser.ROLLBACK:
			err = proc.Rollback(stmt.(parser.Expression))
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
			flow = TerminateWithError
			err = NewForcedExit(code)
		} else {
			flow = Exit
		}
	case parser.Return:
		var ret value.Primary
		if ret, err = proc.Filter.Evaluate(ctx, stmt.(parser.Return).Value); err == nil {
			proc.returnVal = ret
			flow = Return
		}
	case parser.If:
		flow, err = proc.IfStmt(ctx, stmt.(parser.If))
	case parser.Case:
		flow, err = proc.Case(ctx, stmt.(parser.Case))
	case parser.While:
		flow, err = proc.While(ctx, stmt.(parser.While))
	case parser.WhileInCursor:
		flow, err = proc.WhileInCursor(ctx, stmt.(parser.WhileInCursor))
	case parser.Echo:
		if printstr, err = Echo(ctx, proc.Filter, stmt.(parser.Echo)); err == nil {
			proc.Log(printstr, false)
		}
	case parser.Print:
		if printstr, err = Print(ctx, proc.Filter, stmt.(parser.Print)); err == nil {
			proc.Log(printstr, false)
		}
	case parser.Printf:
		if printstr, err = Printf(ctx, proc.Filter, stmt.(parser.Printf)); err == nil {
			proc.Log(printstr, false)
		}
	case parser.Source:
		var externalStatements []parser.Statement
		if externalStatements, err = Source(ctx, proc.Filter, stmt.(parser.Source)); err == nil {
			flow, err = proc.execute(ctx, externalStatements)
		}
	case parser.Execute:
		var externalStatements []parser.Statement
		if externalStatements, err = ParseExecuteStatements(ctx, proc.Filter, stmt.(parser.Execute)); err == nil {
			flow, err = proc.execute(ctx, externalStatements)
		}
	case parser.Chdir:
		err = Chdir(ctx, proc.Filter, stmt.(parser.Chdir))
	case parser.Pwd:
		var dirpath string
		dirpath, err = Pwd(stmt.(parser.Pwd))
		if err == nil {
			proc.Log(dirpath, false)
		}
	case parser.Reload:
		err = Reload(ctx, proc.Tx, stmt.(parser.Reload))
	case parser.ShowObjects:
		if printstr, err = ShowObjects(proc.Filter, stmt.(parser.ShowObjects)); err == nil {
			proc.Log(printstr, false)
		}
	case parser.ShowFields:
		if printstr, err = ShowFields(ctx, proc.Filter, stmt.(parser.ShowFields)); err == nil {
			proc.Log(printstr, false)
		}
	case parser.Syntax:
		printstr = Syntax(ctx, proc.Filter, stmt.(parser.Syntax))
		proc.Log(printstr, false)
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
					if p, err = proc.Filter.Evaluate(ctx, trigger.Message); err == nil {
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
		err = proc.ExecExternalCommand(ctx, stmt.(parser.ExternalCommand))
	default:
		if expr, ok := stmt.(parser.QueryExpression); ok {
			_, err = proc.Filter.Evaluate(ctx, expr)
		}
	}

	if err != nil {
		flow = TerminateWithError
	}
	return flow, err
}

func (proc *Processor) IfStmt(ctx context.Context, stmt parser.If) (StatementFlow, error) {
	stmts := make([]parser.ElseIf, 0, len(stmt.ElseIf)+1)
	stmts = append(stmts, parser.ElseIf{
		Condition:  stmt.Condition,
		Statements: stmt.Statements,
	})
	for _, v := range stmt.ElseIf {
		stmts = append(stmts, v)
	}

	for _, v := range stmts {
		p, err := proc.Filter.Evaluate(ctx, v.Condition)
		if err != nil {
			return TerminateWithError, err
		}
		if p.Ternary() == ternary.TRUE {
			return proc.executeChild(ctx, v.Statements)
		}
	}

	if stmt.Else.Statements != nil {
		return proc.executeChild(ctx, stmt.Else.Statements)
	}
	return Terminate, nil
}

func (proc *Processor) Case(ctx context.Context, stmt parser.Case) (StatementFlow, error) {
	var val value.Primary
	var err error
	if stmt.Value != nil {
		val, err = proc.Filter.Evaluate(ctx, stmt.Value)
		if err != nil {
			return TerminateWithError, err
		}
	}

	for _, when := range stmt.When {
		var t ternary.Value

		cond, err := proc.Filter.Evaluate(ctx, when.Condition)
		if err != nil {
			return TerminateWithError, err
		}

		if val == nil {
			t = cond.Ternary()
		} else {
			t = value.Equal(val, cond, proc.Tx.Flags.DatetimeFormat)
		}

		if t == ternary.TRUE {
			return proc.executeChild(ctx, when.Statements)
		}
	}

	if stmt.Else.Statements == nil {
		return Terminate, nil
	}
	return proc.executeChild(ctx, stmt.Else.Statements)
}

func (proc *Processor) While(ctx context.Context, stmt parser.While) (StatementFlow, error) {
	childProc := proc.NewChildProcessor()

	for {
		childProc.Filter.ResetCurrentScope()
		p, err := proc.Filter.Evaluate(ctx, stmt.Condition)
		if err != nil {
			return TerminateWithError, err
		}
		if p.Ternary() != ternary.TRUE {
			break
		}

		f, err := childProc.execute(ctx, stmt.Statements)
		if err != nil {
			return TerminateWithError, err
		}

		switch f {
		case Break:
			return Terminate, nil
		case Exit:
			return Exit, nil
		case Return:
			proc.returnVal = childProc.returnVal
			return Return, nil
		}
	}
	return Terminate, nil
}

func (proc *Processor) WhileInCursor(ctx context.Context, stmt parser.WhileInCursor) (StatementFlow, error) {
	fetchPosition := parser.FetchPosition{
		Position: parser.Token{Token: parser.NEXT},
	}

	childProc := proc.NewChildProcessor()
	for {
		childProc.Filter.ResetCurrentScope()
		if stmt.WithDeclaration {
			assigns := make([]parser.VariableAssignment, len(stmt.Variables))
			for i, v := range stmt.Variables {
				assigns[i] = parser.VariableAssignment{Variable: v}
			}
			decl := parser.VariableDeclaration{Assignments: assigns}
			if err := childProc.Filter.variables.Declare(ctx, childProc.Filter, decl); err != nil {
				return TerminateWithError, err
			}
		}

		success, err := FetchCursor(ctx, childProc.Filter, stmt.Cursor, fetchPosition, stmt.Variables)
		if err != nil {
			return TerminateWithError, err
		}
		if !success {
			break
		}

		f, err := childProc.execute(ctx, stmt.Statements)
		if err != nil {
			return TerminateWithError, err
		}

		switch f {
		case Break:
			return Terminate, nil
		case Exit:
			return Exit, nil
		case Return:
			proc.returnVal = childProc.returnVal
			return Return, nil
		}
	}

	return Terminate, nil
}

func (proc *Processor) ExecExternalCommand(ctx context.Context, stmt parser.ExternalCommand) error {
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
		arg, err := proc.Filter.EvaluateEmbeddedString(ctx, argStr)
		if err != nil {
			if appErr, ok := err.(Error); ok {
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
	c.Stdin = proc.Tx.Session.Stdin
	c.Stdout = proc.Tx.Session.Stdout
	c.Stderr = proc.Tx.Session.Stderr

	err = c.Run()
	if err != nil {
		err = NewExternalCommandError(stmt, err.Error())
	}
	return err
}

func (proc *Processor) showExecutionTime() {
	palette := cmd.GetPalette()
	exectime := cmd.FormatNumber(time.Since(proc.measurementStart).Seconds(), 6, ".", ",", "")
	stats := fmt.Sprintf(palette.Render(cmd.LableEffect, "Query Execution Time: ")+"%s seconds", exectime)
	proc.Log(stats, false)
}

func (proc *Processor) Log(log string, quiet bool) {
	proc.Tx.Session.Log(log, quiet)
}

func (proc *Processor) LogNotice(log string, quiet bool) {
	proc.Tx.Session.LogNotice(log, quiet)
}

func (proc *Processor) LogWarn(log string, quiet bool) {
	proc.Tx.Session.LogWarn(log, quiet)
}

func (proc *Processor) LogError(log string) {
	proc.Tx.Session.LogError(log)
}

func (proc *Processor) AutoCommit() error {
	return proc.Commit(nil)
}

func (proc *Processor) Commit(expr parser.Expression) error {
	return proc.Tx.Commit(proc.Filter, expr)
}

func (proc *Processor) AutoRollback() error {
	return proc.Rollback(nil)
}

func (proc *Processor) Rollback(expr parser.Expression) error {
	return proc.Tx.Rollback(proc.Filter, expr)
}

func (proc *Processor) ReleaseResources() error {
	return proc.Tx.ReleaseResources()
}

func (proc *Processor) ReleaseResourcesWithErrors() error {
	return proc.Tx.ReleaseResourcesWithErrors()
}
