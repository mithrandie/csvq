package query

import (
	"fmt"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type StatementType int

const (
	SELECT StatementType = iota
	INSERT
	UPDATE
	DELETE
	CREATE_TABLE
	ADD_COLUMNS
	DROP_COLUMNS
	RENAME_COLUMN
)

type Result struct {
	Type          StatementType
	View          *View
	FileInfo      *FileInfo
	OperatedCount int
}

type Procedure struct {
	VariablesList []Variables
	Results       []Result
	Logs          []string
}

func NewProcedure() *Procedure {
	return &Procedure{
		VariablesList: []Variables{{}},
		Results:       []Result{},
		Logs:          []string{},
	}
}

func (proc *Procedure) NewChildProcedure() *Procedure {
	return &Procedure{
		VariablesList: append([]Variables{{}}, proc.VariablesList...),
		Results:       []Result{},
		Logs:          []string{},
	}
}

func (proc *Procedure) AddLog(log string) {
	proc.Logs = append(proc.Logs, log)
}

func (proc *Procedure) Log() string {
	if len(proc.Logs) < 1 {
		return ""
	}
	lb := cmd.GetFlags().LineBreak
	return strings.Join(proc.Logs, lb.Value()) + lb.Value()
}

func (proc *Procedure) ExecuteChild(statements []parser.Statement) (StatementFlow, error) {
	child := proc.NewChildProcedure()
	f, err := child.Execute(statements)
	proc.Logs = append(proc.Logs, child.Logs...)
	proc.Results = append(proc.Results, child.Results...)
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

	ViewCache.ClearAliases()

	var err error

	var results []Result
	var view *View
	var views []*View
	var printstr string

	filter := NewFilter(proc.VariablesList)

	switch stmt.(type) {
	case parser.SetFlag:
		err = SetFlag(stmt.(parser.SetFlag))
	case parser.VariableDeclaration:
		err = proc.VariablesList[0].Decrare(stmt.(parser.VariableDeclaration), filter)
	case parser.VariableSubstitution:
		_, err = filter.Evaluate(stmt.(parser.Expression))
	case parser.CursorDeclaration:
		decl := stmt.(parser.CursorDeclaration)
		err = Cursors.Add(decl.Cursor.Literal, decl.Query)
	case parser.OpenCursor:
		err = Cursors.Open(stmt.(parser.OpenCursor).Cursor.Literal, filter)
	case parser.CloseCursor:
		err = Cursors.Close(stmt.(parser.CloseCursor).Cursor.Literal)
	case parser.DisposeCursor:
		Cursors.Dispose(stmt.(parser.DisposeCursor).Cursor.Literal)
	case parser.FetchCursor:
		fetch := stmt.(parser.FetchCursor)
		_, err = FetchCursor(fetch.Cursor.Literal, fetch.Position, fetch.Variables, filter)
	case parser.TableDeclaration:
		err = DeclareTable(stmt.(parser.TableDeclaration), filter)
	case parser.SelectQuery:
		if view, err = Select(stmt.(parser.SelectQuery), filter); err == nil {
			results = []Result{
				{
					Type: SELECT,
					View: view,
				},
			}
		}
	case parser.InsertQuery:
		if view, err = Insert(stmt.(parser.InsertQuery), filter); err == nil {
			results = []Result{
				{
					Type:          INSERT,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedRecords,
				},
			}
			proc.AddLog(fmt.Sprintf("%s inserted on %q.", proc.formatCount(view.OperatedRecords, "record"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.UpdateQuery:
		if views, err = Update(stmt.(parser.UpdateQuery), filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          UPDATE,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				proc.AddLog(fmt.Sprintf("%s updated on %q.", proc.formatCount(v.OperatedRecords, "record"), v.FileInfo.Path))

				v.OperatedRecords = 0
			}
		}
	case parser.DeleteQuery:
		if views, err = Delete(stmt.(parser.DeleteQuery), filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          DELETE,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				proc.AddLog(fmt.Sprintf("%s deleted on %q.", proc.formatCount(v.OperatedRecords, "record"), v.FileInfo.Path))

				v.OperatedRecords = 0
			}
		}
	case parser.CreateTable:
		if view, err = CreateTable(stmt.(parser.CreateTable)); err == nil {
			results = []Result{
				{
					Type:     CREATE_TABLE,
					FileInfo: view.FileInfo,
				},
			}
			proc.AddLog(fmt.Sprintf("file %q is created.", view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.AddColumns:
		if view, err = AddColumns(stmt.(parser.AddColumns), filter); err == nil {
			results = []Result{
				{
					Type:          ADD_COLUMNS,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			proc.AddLog(fmt.Sprintf("%s added on %q.", proc.formatCount(view.OperatedFields, "field"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.DropColumns:
		if view, err = DropColumns(stmt.(parser.DropColumns)); err == nil {
			results = []Result{
				{
					Type:          DROP_COLUMNS,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			proc.AddLog(fmt.Sprintf("%s dropped on %q.", proc.formatCount(view.OperatedFields, "field"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.RenameColumn:
		if view, err = RenameColumn(stmt.(parser.RenameColumn)); err == nil {
			results = []Result{
				{
					Type:          RENAME_COLUMN,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			proc.AddLog(fmt.Sprintf("%s renamed on %q.", proc.formatCount(view.OperatedFields, "field"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.TransactionControl:
		switch stmt.(parser.TransactionControl).Token {
		case parser.COMMIT:
			err = proc.Commit()
		case parser.ROLLBACK:
			proc.Rollback()
		}
	case parser.FlowControl:
		switch stmt.(parser.FlowControl).Token {
		case parser.CONTINUE:
			flow = CONTINUE
		case parser.BREAK:
			flow = BREAK
		case parser.EXIT:
			flow = EXIT
		}
	case parser.If:
		flow, err = proc.IfStmt(stmt.(parser.If))
	case parser.While:
		flow, err = proc.While(stmt.(parser.While))
	case parser.WhileInCursor:
		flow, err = proc.WhileInCursor(stmt.(parser.WhileInCursor))
	case parser.Print:
		if printstr, err = Print(stmt.(parser.Print), filter); err == nil {
			proc.AddLog(printstr)
		}
	case parser.Printf:
		if printstr, err = Printf(stmt.(parser.Printf), filter); err == nil {
			proc.AddLog(printstr)
		}
	case parser.Source:
		var externalStatements []parser.Statement
		if externalStatements, err = Source(stmt.(parser.Source)); err == nil {
			flow, err = proc.Execute(externalStatements)
		}
	}

	if results != nil {
		proc.Results = append(proc.Results, results...)
	}

	if err != nil {
		flow = ERROR
	}
	return flow, err
}

func (proc *Procedure) IfStmt(stmt parser.If) (StatementFlow, error) {
	stmts := make([]parser.ElseIf, len(stmt.ElseIf)+1)
	stmts[0] = parser.ElseIf{
		Condition:  stmt.Condition,
		Statements: stmt.Statements,
	}
	for i, v := range stmt.ElseIf {
		stmts[i+1] = v.(parser.ElseIf)
	}

	filter := NewFilter(proc.VariablesList)
	for _, v := range stmts {
		p, err := filter.Evaluate(v.Condition)
		if err != nil {
			return ERROR, err
		}
		if p.Ternary() == ternary.TRUE {
			return proc.ExecuteChild(v.Statements)
		}
	}

	if stmt.Else != nil {
		return proc.ExecuteChild(stmt.Else.(parser.Else).Statements)
	}
	return TERMINATE, nil
}

func (proc *Procedure) While(stmt parser.While) (StatementFlow, error) {
	filter := NewFilter(proc.VariablesList)

	for {
		p, err := filter.Evaluate(stmt.Condition)
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
	filter := NewFilter(proc.VariablesList)

	for {
		success, err := FetchCursor(stmt.Cursor.Literal, nil, stmt.Variables, filter)
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

func (proc *Procedure) formatCount(i int, obj string) string {
	var s string
	if i == 0 {
		s = fmt.Sprintf("no %s", obj)
	} else if i == 1 {
		s = fmt.Sprintf("%d %s", i, obj)
	} else {
		s = fmt.Sprintf("%d %ss", i, obj)
	}
	return s
}

func (proc *Procedure) Commit() error {
	flags := cmd.GetFlags()

	var createFiles = map[string]*FileInfo{}
	var updateFiles = map[string]*FileInfo{}

	for _, result := range proc.Results {
		if result.View != nil {
			//SELECT
			viewstr, err := EncodeView(result.View, flags.Format, flags.WriteDelimiter, flags.WithoutHeader, flags.WriteEncoding, flags.LineBreak)
			if err != nil {
				return err
			}
			proc.AddLog(viewstr)
		} else if result.FileInfo != nil {
			//CREATE or UPDATE
			switch result.Type {
			case CREATE_TABLE:
				createFiles[result.FileInfo.Path] = result.FileInfo
			default:
				if !result.FileInfo.Temporary && 0 < result.OperatedCount {
					if _, ok := createFiles[result.FileInfo.Path]; !ok {
						if _, ok := updateFiles[result.FileInfo.Path]; !ok {
							updateFiles[result.FileInfo.Path] = result.FileInfo
						}
					}
				}
			}
		}
	}

	var modified bool

	if 0 < len(createFiles) {
		for pt, fi := range createFiles {
			view, _ := ViewCache.Get(pt)
			viewstr, err := EncodeView(view, cmd.CSV, fi.Delimiter, false, fi.Encoding, fi.LineBreak)
			if err != nil {
				return err
			}

			if err = cmd.CreateFile(pt, viewstr); err != nil {
				return err
			}
			proc.AddLog(fmt.Sprintf("Commit: file %q is created.", pt))
			if !modified {
				modified = true
			}
		}
	}

	if 0 < len(updateFiles) {
		for pt, fi := range updateFiles {
			view, _ := ViewCache.Get(pt)
			viewstr, err := EncodeView(view, cmd.CSV, fi.Delimiter, fi.NoHeader, fi.Encoding, fi.LineBreak)
			if err != nil {
				return err
			}

			if err = cmd.UpdateFile(pt, viewstr); err != nil {
				return err
			}
			proc.AddLog(fmt.Sprintf("Commit: file %q is updated.", pt))
			if !modified {
				modified = true
			}
		}
	}

	proc.Results = []Result{}
	ViewCache.Clear()

	return nil
}

func (proc *Procedure) Rollback() {
	proc.Results = []Result{}
	ViewCache.Clear()

	proc.AddLog("Rolled back.")
	return
}
