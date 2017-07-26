package query

import (
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
	RETURN
)

type StatementType int

const (
	INSERT StatementType = iota
	UPDATE
	DELETE
	CREATE_TABLE
	ADD_COLUMNS
	DROP_COLUMNS
	RENAME_COLUMN
)

type Result struct {
	Type          StatementType
	FileInfo      *FileInfo
	OperatedCount int
}

var ViewCache = ViewMap{}
var Results = []Result{}
var Logs = []string{}
var SelectLogs = []string{}

func AddLog(log string) {
	Logs = append(Logs, log)
}

func AddSelectLog(log string) {
	SelectLogs = append(SelectLogs, log)
}

func ReadLog() string {
	if len(Logs) < 1 {
		return ""
	}
	return strings.Join(Logs, "\n") + "\n"
}

func ReadSelectLog() string {
	if len(SelectLogs) < 1 {
		return ""
	}
	lb := cmd.GetFlags().LineBreak
	return strings.Join(SelectLogs, lb.Value()) + lb.Value()
}

func Execute(input string, sourceFile string) (string, string, error) {
	statements, err := parser.Parse(input, sourceFile)
	if err != nil {
		syntaxErr := err.(*parser.SyntaxError)
		return "", "", NewSyntaxError(syntaxErr.Message, syntaxErr.Line, syntaxErr.Char, syntaxErr.SourceFile)
	}

	Init()

	proc := NewProcedure()
	flow, err := proc.Execute(statements)

	if flow == TERMINATE {
		err = proc.Commit(nil)
	}

	return ReadLog(), ReadSelectLog(), err
}

func Init() {
	DefineAnalyticFunctions()
}

func FetchCursor(name parser.Identifier, fetchPosition parser.Expression, vars []parser.Variable, filter Filter) (bool, error) {
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
				return false, NewInvalidFetchPositionError(fp)
			}
			number = int(i.(parser.Integer).Value())
		}
	}

	primaries, err := filter.CursorsList.Fetch(name, position, number)
	if err != nil {
		return false, err
	}
	if primaries == nil {
		return false, nil
	}
	if len(vars) != len(primaries) {
		return false, NewCursorFetchLengthError(name, len(primaries))
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
	if filter.TempViewsList.Exists(expr.Table.Literal) {
		return NewTemporaryTableRedeclaredError(expr.Table)
	}

	var view *View
	var err error

	if expr.Query != nil {
		view, err = Select(expr.Query.(parser.SelectQuery), filter)
		if err != nil {
			return err
		}

		if err := view.Header.Update(expr.Table.Literal, expr.Fields); err != nil {
			if _, ok := err.(*FieldLengthNotMatchError); ok {
				return NewTemporaryTableFieldLengthError(expr.Query.(parser.SelectQuery), expr.Table, len(expr.Fields))
			}
			return err
		}
	} else {
		fields := make([]string, len(expr.Fields))
		for i, v := range expr.Fields {
			f, _ := v.(parser.Identifier)
			if InStrSliceWithCaseInsensitive(f.Literal, fields) {
				return NewDuplicateFieldNameError(f)
			}
			fields[i] = f.Literal
		}
		header := NewHeaderWithoutId(expr.Table.Literal, fields)
		view = &View{
			Header:  header,
			Records: Records{},
		}
	}

	view.FileInfo = &FileInfo{
		Path:           expr.Table.Literal,
		Temporary:      true,
		InitialRecords: view.Records.Copy(),
	}

	filter.TempViewsList.Set(view)

	return err
}

func Select(query parser.SelectQuery, parentFilter Filter) (*View, error) {
	filter := parentFilter.CreateNode()

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
			return nil, NewCombinedSetFieldLengthError(set.RHS, lview.FieldLen())
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
		err := view.Header.Update(tmpViewName, filter.RecursiveTable.Fields)
		if err != nil {
			return err
		}
		filter.RecursiveTmpView = view
	}

	rview, err := selectSetEntity(set.RHS, filter.CreateNode())
	if err != nil {
		return err
	}
	if view.FieldLen() != rview.FieldLen() {
		return NewCombinedSetFieldLengthError(set.RHS, view.FieldLen())
	}

	if rview.RecordLen() < 1 {
		return nil
	}
	rview.Header.Update(tmpViewName, filter.RecursiveTable.Fields)
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

func Insert(query parser.InsertQuery, parentFilter Filter) (*View, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	view := NewView()
	err := view.LoadFromTableIdentifier(query.Table, filter)
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

	if view.FileInfo.Temporary {
		filter.TempViewsList.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view, nil
}

func Update(query parser.UpdateQuery, parentFilter Filter) ([]*View, error) {
	filter := parentFilter.CreateNode()

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
		fpath, err := filter.AliasesList.Get(table.Name())
		if err != nil {
			return nil, err
		}
		viewKey := strings.ToUpper(table.Name().Literal)

		if parentFilter.TempViewsList.Exists(fpath) {
			viewsToUpdate[viewKey], _ = parentFilter.TempViewsList.Get(parser.Identifier{Literal: fpath})
		} else {
			viewsToUpdate[viewKey], _ = ViewCache.Get(parser.Identifier{Literal: fpath})
		}
		viewsToUpdate[viewKey].Header.Update(table.Name().Literal, nil)
		updatedIndices[viewKey] = []int{}
	}

	filterForLoop := NewFilterForSequentialEvaluation(view, filter)
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
			viewref = strings.ToUpper(viewref)

			if _, ok := viewsToUpdate[viewref]; !ok {
				return nil, NewUpdateFieldNotExistError(uset.Field)
			}

			internalId, err := view.InternalRecordId(viewref, i)
			if err != nil {
				return nil, NewUpdateValueAmbiguousError(uset.Field, uset.Value)
			}

			if InIntSlice(internalId, updatedIndices[viewref]) {
				return nil, NewUpdateValueAmbiguousError(uset.Field, uset.Value)
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
		v.Header.Update(parser.FormatTableName(v.FileInfo.Path), nil)
		v.OperatedRecords = len(updatedIndices[k])

		if v.FileInfo.Temporary {
			filter.TempViewsList.Replace(v)
		} else {
			ViewCache.Replace(v)
		}

		views = append(views, v)
	}

	return views, nil
}

func Delete(query parser.DeleteQuery, parentFilter Filter) ([]*View, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	fromClause := query.FromClause.(parser.FromClause)
	if query.Tables == nil {
		if 1 < len(fromClause.Tables) {
			return nil, NewDeleteTableNotSpecifiedError(query)
		}
		table := fromClause.Tables[0].(parser.Table)
		if _, ok := table.Object.(parser.Identifier); !ok {
			if _, ok := table.Object.(parser.Stdin); !ok {
				return nil, NewDeleteTableNotSpecifiedError(query)
			}
		}
		query.Tables = fromClause.Tables
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
		fpath, err := filter.AliasesList.Get(table.Name())
		if err != nil {
			return nil, err
		}

		viewKey := strings.ToUpper(table.Name().Literal)
		if parentFilter.TempViewsList.Exists(fpath) {
			viewsToDelete[viewKey], _ = parentFilter.TempViewsList.Get(parser.Identifier{Literal: fpath})
		} else {
			viewsToDelete[viewKey], _ = ViewCache.Get(parser.Identifier{Literal: fpath})
		}
		viewsToDelete[viewKey].Header.Update(table.Name().Literal, nil)
		deletedIndices[viewKey] = []int{}
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
		v.Header.Update(parser.FormatTableName(v.FileInfo.Path), nil)
		v.OperatedRecords = len(deletedIndices[k])

		if v.FileInfo.Temporary {
			filter.TempViewsList.Replace(v)
		} else {
			ViewCache.Replace(v)
		}

		views = append(views, v)
	}

	return views, nil
}

func CreateTable(query parser.CreateTable) (*View, error) {
	fields := make([]string, len(query.Fields))
	for i, v := range query.Fields {
		f, _ := v.(parser.Identifier)
		if InStrSliceWithCaseInsensitive(f.Literal, fields) {
			return nil, NewDuplicateFieldNameError(f)
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

	ViewCache.Set(view)

	return view, nil
}

func AddColumns(query parser.AddColumns, parentFilter Filter) (*View, error) {
	filter := parentFilter.CreateNode()

	if query.Position == nil {
		query.Position = parser.ColumnPosition{
			Position: parser.Token{Token: parser.LAST, Literal: parser.TokenLiteral(parser.LAST)},
		}
	}

	view := NewView()
	err := view.LoadFromTableIdentifier(query.Table, filter)
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
		if InStrSliceWithCaseInsensitive(col.Column.Literal, columnNames) || InStrSliceWithCaseInsensitive(col.Column.Literal, fields) {
			return nil, NewDuplicateFieldNameError(col.Column)
		}
		fields[i] = col.Column.Literal
		defaults[i] = col.Value
	}
	newFieldLen := view.FieldLen() + len(query.Columns)

	addHeader := NewHeaderWithoutId(parser.FormatTableName(view.FileInfo.Path), fields)
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
	colNumber := 0
	for i := range header {
		colNumber++
		header[i].Number = colNumber
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

	if view.FileInfo.Temporary {
		filter.TempViewsList.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view, nil
}

func DropColumns(query parser.DropColumns, parentFilter Filter) (*View, error) {
	filter := parentFilter.CreateNode()

	view := NewView()
	err := view.LoadFromTableIdentifier(query.Table, filter)
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

	if view.FileInfo.Temporary {
		filter.TempViewsList.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view, nil

}

func RenameColumn(query parser.RenameColumn, parentFilter Filter) (*View, error) {
	filter := parentFilter.CreateNode()

	view := NewView()
	err := view.LoadFromTableIdentifier(query.Table, filter)
	if err != nil {
		return nil, err
	}

	columnNames := view.Header.TableColumnNames()
	if InStrSliceWithCaseInsensitive(query.New.Literal, columnNames) {
		return nil, NewDuplicateFieldNameError(query.New)
	}

	idx, err := view.FieldIndex(query.Old)
	if err != nil {
		return nil, err
	}

	view.Header[idx].Column = query.New.Literal
	view.OperatedFields = 1
	view.ParentFilter = Filter{}

	if view.FileInfo.Temporary {
		filter.TempViewsList.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view, nil
}
