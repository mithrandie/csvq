package query

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
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
	PRINT
)

type Result struct {
	Type  StatementType
	View  *View
	Count int
	Log   string
}

func Execute(input string) ([]Result, error) {
	results := []Result{}

	parser.SetDebugLevel(0, true)
	program, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	for _, stmt := range program {
		Variable.ClearAutoIncrement()
		ViewCache.Clear()

		switch stmt.(type) {
		case parser.VariableDeclaration:
			if err := Variable.Decrare(stmt.(parser.VariableDeclaration), nil); err != nil {
				return nil, err
			}
		case parser.VariableSubstitution:
			if _, err := Variable.Substitute(stmt.(parser.VariableSubstitution), nil); err != nil {
				return nil, err
			}
		case parser.SelectQuery:
			view, err := Select(stmt.(parser.SelectQuery), nil)
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  SELECT,
				View:  view,
				Count: view.RecordLen(),
			})
		case parser.InsertQuery:
			view, err := Insert(stmt.(parser.InsertQuery))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  INSERT,
				View:  view,
				Count: view.OperatedRecords,
			})
		case parser.UpdateQuery:
			views, err := Update(stmt.(parser.UpdateQuery))
			if err != nil {
				return nil, err
			}
			for _, view := range views {
				results = append(results, Result{
					Type:  UPDATE,
					View:  view,
					Count: view.OperatedRecords,
				})
			}
		case parser.DeleteQuery:
			views, err := Delete(stmt.(parser.DeleteQuery))
			if err != nil {
				return nil, err
			}
			for _, view := range views {
				results = append(results, Result{
					Type:  DELETE,
					View:  view,
					Count: view.OperatedRecords,
				})
			}
		case parser.CreateTable:
			view, err := CreateTable(stmt.(parser.CreateTable))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type: CREATE_TABLE,
				View: view,
			})
		case parser.AddColumns:
			view, err := AddColumns(stmt.(parser.AddColumns))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  ADD_COLUMNS,
				View:  view,
				Count: view.OperatedFields,
			})
		case parser.DropColumns:
			view, err := DropColumns(stmt.(parser.DropColumns))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  DROP_COLUMNS,
				View:  view,
				Count: view.OperatedFields,
			})
		case parser.RenameColumn:
			view, err := RenameColumn(stmt.(parser.RenameColumn))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type:  RENAME_COLUMN,
				View:  view,
				Count: view.OperatedFields,
			})
		case parser.Print:
			log, err := Print(stmt.(parser.Print))
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				Type: PRINT,
				Log:  log,
			})
		}
	}

	return results, nil
}

func Select(query parser.SelectQuery, parentFilter Filter) (*View, error) {
	if query.FromClause == nil {
		query.FromClause = parser.FromClause{}
	}
	view := NewView()
	err := view.Load(query.FromClause.(parser.FromClause), parentFilter)
	if err != nil {
		return nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if query.GroupByClause != nil {
		if err := view.GroupBy(query.GroupByClause.(parser.GroupByClause)); err != nil {
			return nil, err
		}
	}

	if query.HavingClause != nil {
		if err := view.Having(query.HavingClause.(parser.HavingClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if err := view.Select(query.SelectClause.(parser.SelectClause)); err != nil {
		return nil, err
	}

	if query.OrderByClause != nil {
		if err := view.OrderBy(query.OrderByClause.(parser.OrderByClause)); err != nil {
			return nil, err
		}
	}

	if query.LimitClause != nil {
		view.Limit(query.LimitClause.(parser.LimitClause))
	}

	view.Fix()

	return view, nil
}

func Insert(query parser.InsertQuery) (*View, error) {
	view := NewView()
	view.UseCache = false
	view.UseInternalId = false
	err := view.LoadFromIdentifier(query.Table)
	if err != nil {
		return nil, err
	}

	fields := query.Fields
	if fields == nil {
		fields = view.Header.TableColumns()
	}

	if query.ValuesList != nil {
		if err := view.InsertValues(fields, query.ValuesList); err != nil {
			return nil, err
		}
	} else {
		if err := view.InsertFromQuery(fields, query.Query.(parser.SelectQuery)); err != nil {
			return nil, err
		}
	}

	return view, nil
}

func Update(query parser.UpdateQuery) ([]*View, error) {
	if query.FromClause == nil {
		query.FromClause = parser.FromClause{Tables: query.Tables}
	}

	view := NewView()
	err := view.Load(query.FromClause.(parser.FromClause), nil)
	if err != nil {
		return nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	viewsToUpdate := make(map[string]*View)
	updatedIndices := make(map[string][]int)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		if viewsToUpdate[table.Name()], err = ViewCache.Get(table.Name()); err != nil {
			return nil, err
		}
		updatedIndices[table.Name()] = []int{}
	}

	for i := range view.Records {
		var filter Filter = []FilterRecord{{View: view, RecordIndex: i}}

		for _, v := range query.SetList {
			uset := v.(parser.UpdateSet)

			value, err := filter.Evaluate(uset.Value)
			if err != nil {
				return nil, err
			}

			viewref, err := view.FieldRef(uset.Field)
			if err != nil {
				return nil, err
			}

			internalId, err := view.InternalRecordId(viewref, i)
			if err != nil {
				return nil, errors.New("record to update is ambiguous")
			}

			if InIntSlice(internalId, updatedIndices[viewref]) {
				return nil, errors.New("record to update is ambiguous")
			}

			fieldIdx, _ := viewsToUpdate[viewref].FieldIndex(uset.Field)

			viewsToUpdate[viewref].Records[internalId][fieldIdx] = NewCell(value)
			updatedIndices[viewref] = append(updatedIndices[viewref], internalId)
		}
	}

	views := []*View{}
	for k, v := range viewsToUpdate {
		if err := v.SelectAllColumns(); err != nil {
			return nil, err
		}

		v.Fix()
		v.OperatedRecords = len(updatedIndices[k])
		views = append(views, v)
	}

	return views, nil
}

func Delete(query parser.DeleteQuery) ([]*View, error) {
	fromClause := query.FromClause.(parser.FromClause)
	if query.Tables == nil {
		table := fromClause.Tables[0].(parser.Table)
		if _, ok := table.Object.(parser.Identifier); !ok || 1 < len(fromClause.Tables) {
			return nil, errors.New("update file is not specified")
		}
		query.Tables = []parser.Expression{table}
	}

	view := NewView()
	err := view.Load(query.FromClause.(parser.FromClause), nil)
	if err != nil {
		return nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	viewsToDelete := make(map[string]*View)
	deletedIndices := make(map[string][]int)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		if viewsToDelete[table.Name()], err = ViewCache.Get(table.Name()); err != nil {
			return nil, err
		}
		deletedIndices[table.Name()] = []int{}
	}

	for i := range view.Records {
		for viewref := range viewsToDelete {
			internalId, err := view.InternalRecordId(viewref, i)
			if err != nil {
				continue
			}
			if InIntSlice(internalId, deletedIndices[viewref]) {
				continue
			}
			deletedIndices[viewref] = append(deletedIndices[viewref], internalId)
		}
	}

	views := []*View{}
	for k, v := range viewsToDelete {
		filterdIndices := []int{}
		for i := range v.Records {
			if !InIntSlice(i, deletedIndices[k]) {
				filterdIndices = append(filterdIndices, i)
			}
		}
		v.filteredIndices = filterdIndices
		v.Extract()

		if err := v.SelectAllColumns(); err != nil {
			return nil, err
		}

		v.Fix()
		v.OperatedRecords = len(deletedIndices[k])
		views = append(views, v)
	}

	return views, nil
}

func CreateTable(query parser.CreateTable) (*View, error) {
	fields := make([]string, len(query.Fields))
	for i, v := range query.Fields {
		f, _ := v.(parser.Identifier)
		if InStrSlice(f.Literal, fields) {
			return nil, errors.New(fmt.Sprintf("field %s is duplicate", f))
		}
		fields[i] = f.Literal
	}

	flags := cmd.GetFlags()
	filepath := query.Table.Literal
	if !path.IsAbs(filepath) {
		filepath = path.Join(flags.Repository, filepath)
	}
	delimiter := flags.Delimiter
	if delimiter == cmd.UNDEF {
		if strings.EqualFold(path.Ext(filepath), cmd.TSV_EXT) {
			delimiter = '\t'
		} else {
			delimiter = ','
		}
	}

	header := NewHeaderWithoutId(parser.FormatTableName(query.Table.Literal), fields)
	view := &View{
		Header: header,
		FileInfo: &FileInfo{
			Path:      filepath,
			Delimiter: delimiter,
		},
	}
	return view, nil
}

func AddColumns(query parser.AddColumns) (*View, error) {
	if query.Position == nil {
		query.Position = parser.ColumnPosition{
			Position: parser.Token{Token: parser.LAST, Literal: parser.TokenLiteral(parser.LAST)},
		}
	}

	view := NewView()
	view.UseCache = false
	view.UseInternalId = false
	err := view.LoadFromIdentifier(query.Table)
	if err != nil {
		return nil, err
	}

	var insertPos int
	pos, _ := query.Position.(parser.ColumnPosition)
	switch pos.Position.Token {
	case parser.FIRST:
		insertPos = 0
	case parser.LAST:
		insertPos = view.FieldLen()
	default:
		idx, err := view.FieldIndex(pos.Column.(parser.Identifier))
		if err != nil {
			return nil, err
		}
		switch pos.Position.Token {
		case parser.BEFORE:
			insertPos = idx
		default: //parser.AFTER
			insertPos = idx + 1
		}
	}

	columnNames := view.Header.TableColumnNames()
	fields := make([]string, len(query.Columns))
	defaults := make([]parser.Expression, len(query.Columns))
	for i, v := range query.Columns {
		col := v.(parser.ColumnDefault)
		if InStrSlice(col.Column.Literal, columnNames) || InStrSlice(col.Column.Literal, fields) {
			return nil, errors.New(fmt.Sprintf("field %s is duplicate", col.Column))
		}
		fields[i] = col.Column.Literal
		defaults[i] = col.Value
	}
	newFieldLen := view.FieldLen() + len(query.Columns)

	addHeader := NewHeaderWithoutId(parser.FormatTableName(query.Table.Literal), fields)
	header := make(Header, newFieldLen)
	for i, v := range view.Header {
		var idx int
		if i < insertPos {
			idx = i
		} else {
			idx = i + len(fields)
		}
		header[idx] = v
	}
	for i, v := range addHeader {
		header[i+insertPos] = v
	}

	records := make([]Record, view.RecordLen())
	for i, v := range view.Records {
		record := make(Record, newFieldLen)
		for j, cell := range v {
			var idx int
			if j < insertPos {
				idx = j
			} else {
				idx = j + len(fields)
			}
			record[idx] = cell
		}

		var filter Filter = []FilterRecord{{View: view, RecordIndex: i}}
		for j, v := range defaults {
			if v == nil {
				v = parser.NewNull()
			}
			val, err := filter.Evaluate(v)
			if err != nil {
				return nil, err
			}
			record[j+insertPos] = NewCell(val)
		}

		records[i] = record
	}

	view.Header = header
	view.Records = records
	view.OperatedFields = len(fields)

	return view, nil
}

func DropColumns(query parser.DropColumns) (*View, error) {
	view := NewView()
	view.UseCache = false
	view.UseInternalId = false
	err := view.LoadFromIdentifier(query.Table)
	if err != nil {
		return nil, err
	}

	dropIndices := make([]int, len(query.Columns))
	for i, v := range query.Columns {
		idx, err := view.FieldIndex(v.(parser.Identifier))
		if err != nil {
			return nil, err
		}
		dropIndices[i] = idx
	}

	view.selectFields = []int{}
	for i := 0; i < view.FieldLen(); i++ {
		if view.Header[i].FromTable && !InIntSlice(i, dropIndices) {
			view.selectFields = append(view.selectFields, i)
		}
	}

	view.Fix()
	view.OperatedFields = len(dropIndices)

	return view, nil

}

func RenameColumn(query parser.RenameColumn) (*View, error) {
	view := NewView()
	view.UseCache = false
	view.UseInternalId = false
	err := view.LoadFromIdentifier(query.Table)
	if err != nil {
		return nil, err
	}

	columnNames := view.Header.TableColumnNames()
	if InStrSlice(query.New.Literal, columnNames) {
		return nil, errors.New(fmt.Sprintf("field %s is duplicate", query.New))
	}

	idx, err := view.FieldIndex(query.Old)
	if err != nil {
		return nil, err
	}

	view.Header[idx].Column = query.New.Literal
	view.OperatedFields = 1

	return view, nil
}

func Print(query parser.Print) (string, error) {
	var filter Filter
	p, err := filter.Evaluate(query.Value)
	if err != nil {
		return "", err
	}
	return p.String(), err
}
