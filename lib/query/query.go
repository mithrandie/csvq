package query

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
)

type StatementFlow int

const (
	TERMINATE StatementFlow = iota
	ERROR
	EXIT
	BREAK
	CONTINUE
)

var ViewCache = NewViewMap()
var Cursors = CursorMap{}

func Execute(input string) (string, error) {
	statements, err := parser.Parse(input)
	if err != nil {
		return "", err
	}

	proc := NewProcedure()
	flow, err := proc.Execute(statements)

	if flow == TERMINATE {
		err = proc.Commit()
	}

	return proc.Log(), err
}

func FetchCursor(name string, fetchPosition parser.Expression, vars []parser.Variable, filter Filter) (bool, error) {
	position := parser.NEXT
	number := -1
	if fetchPosition != nil {
		fp := fetchPosition.(parser.FetchPosition)
		position = fp.Position.Token
		if fp.Number != nil {
			p, err := filter.Evaluate(fp.Number)
			if err != nil {
				return false, err
			}
			i := parser.PrimaryToInteger(p)
			if parser.IsNull(i) {
				return false, errors.New(fmt.Sprintf("fetch position %s is not a integer", fp.Number))
			}
			number = int(i.(parser.Integer).Value())
		}
	}

	primaries, err := Cursors.Fetch(name, position, number)
	if err != nil {
		return false, err
	}
	if primaries == nil {
		return false, nil
	}
	if len(vars) != len(primaries) {
		return false, errors.New(fmt.Sprintf("cursor %s field length does not match variables number", name))
	}

	for i, v := range vars {
		substitution := parser.VariableSubstitution{
			Variable: v,
			Value:    primaries[i],
		}
		_, err := filter.VariablesList[0].Substitute(substitution, filter)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func DeclareTable(expr parser.TableDeclaration, filter Filter) error {
	if _, ok := ViewCache.Exists(expr.Table.Literal); ok {
		return errors.New(fmt.Sprintf("table %s already exists", expr.Table.Literal))
	}

	var view *View
	var err error

	if expr.Query != nil {
		view, err = Select(expr.Query.(parser.SelectQuery), filter)
		if err != nil {
			return err
		}

		if err := view.UpdateHeader(expr.Table.Literal, expr.Fields); err != nil {
			return err
		}
	} else {
		fields := make([]string, len(expr.Fields))
		for i, v := range expr.Fields {
			f, _ := v.(parser.Identifier)
			if InStrSlice(f.Literal, fields) {
				return errors.New(fmt.Sprintf("field %s is a duplicate", f))
			}
			fields[i] = f.Literal
		}
		header := NewHeaderWithoutId(expr.Table.Literal, fields)
		view = &View{
			Header: header,
		}
	}

	view.FileInfo = &FileInfo{
		Path:      expr.Table.Literal,
		Temporary: true,
	}

	ViewCache.Set(view, expr.Table.Literal)

	return err
}

func Select(query parser.SelectQuery, parentFilter Filter) (*View, error) {
	filter := parentFilter.Copy()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	view, err := selectEntity(query.SelectEntity, filter)
	if err != nil {
		return nil, err
	}

	if query.OrderByClause != nil {
		if err := view.OrderBy(query.OrderByClause.(parser.OrderByClause)); err != nil {
			return nil, err
		}
	}

	if query.OffsetClause != nil {
		if err := view.Offset(query.OffsetClause.(parser.OffsetClause)); err != nil {
			return nil, err
		}
	}

	if query.LimitClause != nil {
		if err := view.Limit(query.LimitClause.(parser.LimitClause)); err != nil {
			return nil, err
		}
	}

	view.Fix()

	return view, nil
}

func selectEntity(expr parser.Expression, filter Filter) (*View, error) {
	entity, ok := expr.(parser.SelectEntity)
	if !ok {
		return selectSet(expr.(parser.SelectSet), filter)
	}

	if entity.FromClause == nil {
		entity.FromClause = parser.FromClause{}
	}
	view := NewView()
	err := view.Load(entity.FromClause.(parser.FromClause), filter)
	if err != nil {
		return nil, err
	}

	if entity.WhereClause != nil {
		if err := view.Where(entity.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if entity.GroupByClause != nil {
		if err := view.GroupBy(entity.GroupByClause.(parser.GroupByClause)); err != nil {
			return nil, err
		}
	}

	if entity.HavingClause != nil {
		if err := view.Having(entity.HavingClause.(parser.HavingClause)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if err := view.Select(entity.SelectClause.(parser.SelectClause)); err != nil {
		return nil, err
	}

	return view, nil
}

func selectSetEntity(expr parser.Expression, filter Filter) (*View, error) {
	if subquery, ok := expr.(parser.Subquery); ok {
		return Select(subquery.Query, filter)
	}

	view, err := selectEntity(expr, filter)
	if err != nil {
		return nil, err
	}
	view.Fix()
	return view, nil
}

func selectSet(set parser.SelectSet, filter Filter) (*View, error) {
	lview, err := selectSetEntity(set.LHS, filter)
	if err != nil {
		return nil, err
	}

	if filter.RecursiveTable != nil {
		filter.RecursiveTmpView = nil
		err := selectSetForRecursion(lview, set, filter)
		if err != nil {
			return nil, err
		}
	} else {
		rview, err := selectSetEntity(set.RHS, filter)
		if err != nil {
			return nil, err
		}

		if lview.FieldLen() != rview.FieldLen() {
			return nil, errors.New(fmt.Sprintf("%s: field length does not match", parser.TokenLiteral(set.Operator.Token)))
		}

		switch set.Operator.Token {
		case parser.UNION:
			lview.Union(rview, !set.All.IsEmpty())
		case parser.EXCEPT:
			lview.Except(rview, !set.All.IsEmpty())
		case parser.INTERSECT:
			lview.Intersect(rview, !set.All.IsEmpty())
		}
	}

	lview.SelectAllColumns()

	return lview, nil
}

func selectSetForRecursion(view *View, set parser.SelectSet, filter Filter) error {
	tmpViewName := strings.ToUpper(filter.RecursiveTable.Name.Literal)

	if filter.RecursiveTmpView == nil {
		err := view.UpdateHeader(tmpViewName, filter.RecursiveTable.Columns)
		if err != nil {
			return err
		}
		filter.RecursiveTmpView = view
	}

	rview, err := selectSetEntity(set.RHS, filter)
	if err != nil {
		return err
	}
	if view.FieldLen() != rview.FieldLen() {
		return errors.New(fmt.Sprintf("%s: field length does not match", parser.TokenLiteral(set.Operator.Token)))
	}

	if rview.RecordLen() < 1 {
		return nil
	}
	rview.UpdateHeader(tmpViewName, filter.RecursiveTable.Columns)
	filter.RecursiveTmpView = rview

	switch set.Operator.Token {
	case parser.UNION:
		view.Union(rview, !set.All.IsEmpty())
	case parser.EXCEPT:
		view.Except(rview, !set.All.IsEmpty())
	case parser.INTERSECT:
		view.Intersect(rview, !set.All.IsEmpty())
	}

	return selectSetForRecursion(view, set, filter)
}

func Insert(query parser.InsertQuery, filter Filter) (*View, error) {
	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	view := NewView()
	err := view.LoadFromIdentifier(query.Table, filter)
	if err != nil {
		return nil, err
	}

	fields := query.Fields
	if fields == nil {
		fields = view.Header.TableColumns()
	}

	if query.ValuesList != nil {
		if err := view.InsertValues(fields, query.ValuesList, filter); err != nil {
			return nil, err
		}
	} else {
		if err := view.InsertFromQuery(fields, query.Query.(parser.SelectQuery), filter); err != nil {
			return nil, err
		}
	}

	view.ParentFilter = Filter{}

	ViewCache.Replace(view)

	return view, nil
}

func Update(query parser.UpdateQuery, filter Filter) ([]*View, error) {
	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	if query.FromClause == nil {
		query.FromClause = parser.FromClause{Tables: query.Tables}
	}

	view := NewView()
	view.UseInternalId = true
	err := view.Load(query.FromClause.(parser.FromClause), filter)
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
		viewsToUpdate[table.Name()].UpdateHeader(table.Name(), nil)
		updatedIndices[table.Name()] = []int{}
	}

	filterForLoop := NewFilterForLoop(view, filter)
	for i := range view.Records {
		filterForLoop.Records[0].RecordIndex = i

		for _, v := range query.SetList {
			uset := v.(parser.UpdateSet)

			value, err := filterForLoop.Evaluate(uset.Value)
			if err != nil {
				return nil, err
			}

			viewref, err := view.FieldViewName(uset.Field)
			if err != nil {
				return nil, err
			}
			if _, ok := viewsToUpdate[viewref]; !ok {
				return nil, errors.New(fmt.Sprintf("table %s is not specified in tables to update", viewref))
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

		ViewCache.Replace(v)

		views = append(views, v)
	}

	return views, nil
}

func Delete(query parser.DeleteQuery, filter Filter) ([]*View, error) {
	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	fromClause := query.FromClause.(parser.FromClause)
	if query.Tables == nil {
		table := fromClause.Tables[0].(parser.Table)
		if _, ok := table.Object.(parser.Identifier); !ok || 1 < len(fromClause.Tables) {
			return nil, errors.New("update file is not specified")
		}
		query.Tables = []parser.Expression{table}
	}

	view := NewView()
	view.UseInternalId = true
	err := view.Load(query.FromClause.(parser.FromClause), filter)
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
		viewsToDelete[table.Name()].UpdateHeader(table.Name(), nil)
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

		ViewCache.Replace(v)

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
	fpath := query.Table.Literal
	if !filepath.IsAbs(fpath) {
		fpath = filepath.Join(flags.Repository, fpath)
	}
	delimiter := flags.Delimiter
	if delimiter == cmd.UNDEF {
		if strings.EqualFold(filepath.Ext(fpath), cmd.TSV_EXT) {
			delimiter = '\t'
		} else {
			delimiter = ','
		}
	}

	header := NewHeaderWithoutId(parser.FormatTableName(query.Table.Literal), fields)
	view := &View{
		Header: header,
		FileInfo: &FileInfo{
			Path:      fpath,
			Delimiter: delimiter,
			NoHeader:  false,
			Encoding:  flags.Encoding,
			LineBreak: flags.LineBreak,
		},
	}

	ViewCache.Set(view, parser.FormatTableName(view.FileInfo.Path))

	return view, nil
}

func AddColumns(query parser.AddColumns, filter Filter) (*View, error) {
	if query.Position == nil {
		query.Position = parser.ColumnPosition{
			Position: parser.Token{Token: parser.LAST, Literal: parser.TokenLiteral(parser.LAST)},
		}
	}

	view := NewView()
	err := view.LoadFromIdentifier(query.Table, filter)
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
		idx, err := view.FieldIndex(pos.Column.(parser.FieldReference))
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
	filter.Records = append(filter.Records, FilterRecord{
		View: view,
	})
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

		filter.Records[0].RecordIndex = i
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
	view.ParentFilter = Filter{}

	ViewCache.Replace(view)

	return view, nil
}

func DropColumns(query parser.DropColumns) (*View, error) {
	view := NewView()
	err := view.LoadFromIdentifier(query.Table, NewFilter([]Variables{{}}))
	if err != nil {
		return nil, err
	}

	dropIndices := make([]int, len(query.Columns))
	for i, v := range query.Columns {
		idx, err := view.FieldIndex(v.(parser.FieldReference))
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

	ViewCache.Replace(view)

	return view, nil

}

func RenameColumn(query parser.RenameColumn) (*View, error) {
	view := NewView()
	err := view.LoadFromIdentifier(query.Table, NewFilter([]Variables{{}}))
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
	view.ParentFilter = Filter{}

	ViewCache.Replace(view)

	return view, nil
}
