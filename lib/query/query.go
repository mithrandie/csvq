package query

import (
	"context"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func FetchCursor(ctx context.Context, filter *Filter, name parser.Identifier, fetchPosition parser.FetchPosition, vars []parser.Variable) (bool, error) {
	position := parser.NEXT
	number := -1
	if !fetchPosition.Position.IsEmpty() {
		position = fetchPosition.Position.Token
		if fetchPosition.Number != nil {
			p, err := filter.Evaluate(ctx, fetchPosition.Number)
			if err != nil {
				return false, err
			}
			i := value.ToInteger(p)
			if value.IsNull(i) {
				return false, NewInvalidFetchPositionError(fetchPosition)
			}
			number = int(i.(value.Integer).Raw())
		}
	}

	primaries, err := filter.cursors.Fetch(name, position, number)
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
		_, err := filter.variables.SubstituteDirectly(v, primaries[i])
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func DeclareView(ctx context.Context, filter *Filter, expr parser.ViewDeclaration) error {
	if filter.tempViews.Exists(expr.View.Literal) {
		return NewTemporaryTableRedeclaredError(expr.View)
	}

	var view *View
	var err error

	if expr.Query != nil {
		view, err = Select(ctx, filter, expr.Query.(parser.SelectQuery))
		if err != nil {
			return err
		}

		if err := view.Header.Update(expr.View.Literal, expr.Fields); err != nil {
			if _, ok := err.(*FieldLengthNotMatchError); ok {
				return NewTemporaryTableFieldLengthError(expr.Query.(parser.SelectQuery), expr.View, len(expr.Fields))
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
		header := NewHeader(expr.View.Literal, fields)
		view = NewView(filter.tx)
		view.Header = header
		view.RecordSet = RecordSet{}
	}

	view.FileInfo = &FileInfo{
		Path:             expr.View.Literal,
		IsTemporary:      true,
		InitialHeader:    view.Header.Copy(),
		InitialRecordSet: view.RecordSet.Copy(),
	}

	filter.tempViews.Set(view)

	return err
}

func Select(ctx context.Context, parentFilter *Filter, query parser.SelectQuery) (*View, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(context.Background(), query.WithClause.(parser.WithClause)); err != nil {
			return nil, err
		}
	}

	view, err := selectEntity(ctx, filter, query.SelectEntity, query.ForUpdate)
	if err != nil {
		return nil, err
	}

	if query.OrderByClause != nil {
		if err := view.OrderBy(ctx, query.OrderByClause.(parser.OrderByClause)); err != nil {
			return nil, err
		}
	}

	if query.OffsetClause != nil {
		if err := view.Offset(ctx, query.OffsetClause.(parser.OffsetClause)); err != nil {
			return nil, err
		}
	}

	if query.LimitClause != nil {
		if err := view.Limit(ctx, query.LimitClause.(parser.LimitClause)); err != nil {
			return nil, err
		}
	}

	err = view.Fix(ctx)
	return view, err
}

func selectEntity(ctx context.Context, filter *Filter, expr parser.QueryExpression, forUpdate bool) (*View, error) {
	entity, ok := expr.(parser.SelectEntity)
	if !ok {
		return selectSet(ctx, filter, expr.(parser.SelectSet), forUpdate)
	}

	if entity.FromClause == nil {
		entity.FromClause = parser.FromClause{}
	}
	view := NewView(filter.tx)
	err := view.Load(ctx, filter, entity.FromClause.(parser.FromClause), forUpdate, false)
	if err != nil {
		return nil, err
	}

	if entity.WhereClause != nil {
		if err := view.Where(ctx, entity.WhereClause.(parser.WhereClause)); err != nil {
			return nil, err
		}
	}

	if entity.GroupByClause != nil {
		if err := view.GroupBy(ctx, entity.GroupByClause.(parser.GroupByClause)); err != nil {
			return nil, err
		}
	}

	if entity.HavingClause != nil {
		if err := view.Having(ctx, entity.HavingClause.(parser.HavingClause)); err != nil {
			return nil, err
		}
	}

	if err := view.Select(ctx, entity.SelectClause.(parser.SelectClause)); err != nil {
		return nil, err
	}

	return view, nil
}

func selectSetEntity(ctx context.Context, filter *Filter, expr parser.QueryExpression, forUpdate bool) (*View, error) {
	if subquery, ok := expr.(parser.Subquery); ok {
		return Select(ctx, filter, subquery.Query)
	}

	view, err := selectEntity(ctx, filter, expr, forUpdate)
	if err != nil {
		return nil, err
	}
	err = view.Fix(ctx)
	return view, err
}

func selectSet(ctx context.Context, filter *Filter, set parser.SelectSet, forUpdate bool) (*View, error) {
	lview, err := selectSetEntity(ctx, filter, set.LHS, forUpdate)
	if err != nil {
		return nil, err
	}

	if filter.recursiveTable != nil {
		filter.recursiveTmpView = nil
		err := selectSetForRecursion(ctx, filter, lview, set, forUpdate)
		if err != nil {
			return nil, err
		}
	} else {
		rview, err := selectSetEntity(ctx, filter, set.RHS, forUpdate)
		if err != nil {
			return nil, err
		}

		if lview.FieldLen() != rview.FieldLen() {
			return nil, NewCombinedSetFieldLengthError(set.RHS, lview.FieldLen())
		}

		switch set.Operator.Token {
		case parser.UNION:
			if err = lview.Union(ctx, rview, !set.All.IsEmpty()); err != nil {
				return nil, err
			}
		case parser.EXCEPT:
			if err = lview.Except(ctx, rview, !set.All.IsEmpty()); err != nil {
				return nil, err
			}
		case parser.INTERSECT:
			if err = lview.Intersect(ctx, rview, !set.All.IsEmpty()); err != nil {
				return nil, err
			}
		}
	}

	err = lview.SelectAllColumns(ctx)
	return lview, err
}

func selectSetForRecursion(ctx context.Context, filter *Filter, view *View, set parser.SelectSet, forUpdate bool) error {
	tmpViewName := strings.ToUpper(filter.recursiveTable.Name.Literal)

	if filter.recursiveTmpView == nil {
		err := view.Header.Update(tmpViewName, filter.recursiveTable.Fields)
		if err != nil {
			return err
		}
		filter.recursiveTmpView = view
	}

	rview, err := selectSetEntity(ctx, filter.CreateNode(), set.RHS, forUpdate)
	if err != nil {
		return err
	}
	if view.FieldLen() != rview.FieldLen() {
		return NewCombinedSetFieldLengthError(set.RHS, view.FieldLen())
	}

	if rview.RecordLen() < 1 {
		return nil
	}
	if err = rview.Header.Update(tmpViewName, filter.recursiveTable.Fields); err != nil {
		return err
	}
	filter.recursiveTmpView = rview

	switch set.Operator.Token {
	case parser.UNION:
		if err = view.Union(ctx, rview, !set.All.IsEmpty()); err != nil {
			return err
		}
	case parser.EXCEPT:
		if err = view.Except(ctx, rview, !set.All.IsEmpty()); err != nil {
			return err
		}
	case parser.INTERSECT:
		if err = view.Intersect(ctx, rview, !set.All.IsEmpty()); err != nil {
			return err
		}
	}

	return selectSetForRecursion(ctx, filter, view, set, forUpdate)
}

func Insert(ctx context.Context, parentFilter *Filter, query parser.InsertQuery) (*FileInfo, int, error) {
	filter := parentFilter.CreateNode()

	var insertRecords int

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(context.Background(), query.WithClause.(parser.WithClause)); err != nil {
			return nil, insertRecords, err
		}
	}

	fromClause := parser.FromClause{
		Tables: []parser.QueryExpression{
			query.Table,
		},
	}
	view := NewView(parentFilter.tx)
	err := view.Load(ctx, filter, fromClause, true, false)
	if err != nil {
		return nil, insertRecords, err
	}

	fields := query.Fields
	if fields == nil {
		fields = view.Header.TableColumns()
	}

	if query.ValuesList != nil {
		if insertRecords, err = view.InsertValues(ctx, fields, query.ValuesList); err != nil {
			return nil, insertRecords, err
		}
	} else {
		if insertRecords, err = view.InsertFromQuery(ctx, fields, query.Query.(parser.SelectQuery)); err != nil {
			return nil, insertRecords, err
		}
	}

	if err = view.RestoreHeaderReferences(); err != nil {
		return nil, insertRecords, err
	}
	view.Filter = nil

	if view.FileInfo.IsTemporary {
		filter.tempViews.Replace(view)
	} else {
		err = filter.tx.cachedViews.Replace(view)
	}

	return view.FileInfo, insertRecords, err
}

func Update(ctx context.Context, parentFilter *Filter, query parser.UpdateQuery) ([]*FileInfo, []int, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(context.Background(), query.WithClause.(parser.WithClause)); err != nil {
			return nil, nil, err
		}
	}

	if query.FromClause == nil {
		query.FromClause = parser.FromClause{Tables: query.Tables}
	}

	view := NewView(parentFilter.tx)
	err := view.Load(ctx, filter, query.FromClause.(parser.FromClause), true, true)
	if err != nil {
		return nil, nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(ctx, query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, nil, err
		}
	}

	viewsToUpdate := make(map[string]*View)
	updatedCount := make(map[string]int)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		fpath, err := filter.aliases.Get(table.Name())
		if err != nil {
			return nil, nil, err
		}
		viewKey := strings.ToUpper(table.Name().Literal)

		if filter.tempViews.Exists(fpath) {
			viewsToUpdate[viewKey], _ = filter.tempViews.Get(parser.Identifier{Literal: fpath})
		} else {
			viewsToUpdate[viewKey], _ = filter.tx.cachedViews.Get(parser.Identifier{Literal: fpath})
		}
		if err = viewsToUpdate[viewKey].Header.Update(table.Name().Literal, nil); err != nil {
			return nil, nil, err
		}
	}

	updatesList := make(map[string]map[int][]int)
	filterForLoop := NewFilterForSequentialEvaluation(filter, view)
	for i := range view.RecordSet {
		filterForLoop.records[0].recordIndex = i
		internalIds := make(map[string]int)

		for _, uset := range query.SetList {
			val, err := filterForLoop.Evaluate(ctx, uset.Value)
			if err != nil {
				return nil, nil, err
			}

			viewref, err := view.FieldViewName(uset.Field)
			if err != nil {
				return nil, nil, err
			}
			viewref = strings.ToUpper(viewref)

			if _, ok := viewsToUpdate[viewref]; !ok {
				return nil, nil, NewUpdateFieldNotExistError(uset.Field)
			}

			var internalId int
			if id, ok := internalIds[viewref]; ok {
				internalId = id
			} else {
				id, err := view.InternalRecordId(viewref, i)
				if err != nil {
					return nil, nil, NewUpdateValueAmbiguousError(uset.Field, uset.Value)
				}

				internalId = id
				internalIds[viewref] = internalId
			}

			fieldIdx, _ := viewsToUpdate[viewref].FieldIndex(uset.Field)
			if _, ok := updatesList[viewref]; !ok {
				updatesList[viewref] = make(map[int][]int)
			}
			if _, ok := updatesList[viewref][internalId]; !ok {
				updatesList[viewref][internalId] = []int{}
				updatedCount[viewref]++
			}
			if InIntSlice(fieldIdx, updatesList[viewref][internalId]) {
				return nil, nil, NewUpdateValueAmbiguousError(uset.Field, uset.Value)
			}
			updatesList[viewref][internalId] = append(updatesList[viewref][internalId], fieldIdx)
			viewsToUpdate[viewref].RecordSet[internalId][fieldIdx] = NewCell(val)
		}
	}

	fileInfos := make([]*FileInfo, 0)
	updateRecords := make([]int, 0)
	for k, v := range viewsToUpdate {
		if err = v.RestoreHeaderReferences(); err != nil {
			return nil, nil, err
		}

		if v.FileInfo.IsTemporary {
			filter.tempViews.Replace(v)
		} else {
			if err = filter.tx.cachedViews.Replace(v); err != nil {
				return nil, nil, err
			}
		}

		fileInfos = append(fileInfos, v.FileInfo)
		updateRecords = append(updateRecords, updatedCount[k])
	}

	return fileInfos, updateRecords, nil
}

func Delete(ctx context.Context, parentFilter *Filter, query parser.DeleteQuery) ([]*FileInfo, []int, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(context.Background(), query.WithClause.(parser.WithClause)); err != nil {
			return nil, nil, err
		}
	}

	fromClause := query.FromClause
	if query.Tables == nil {
		if 1 < len(fromClause.Tables) {
			return nil, nil, NewDeleteTableNotSpecifiedError(query)
		}
		table := fromClause.Tables[0].(parser.Table)
		switch table.Object.(type) {
		case parser.Identifier, parser.TableObject, parser.Stdin:
			query.Tables = fromClause.Tables
		default:
			return nil, nil, NewDeleteTableNotSpecifiedError(query)
		}
	}

	view := NewView(parentFilter.tx)
	err := view.Load(ctx, filter, query.FromClause, true, true)
	if err != nil {
		return nil, nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(ctx, query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, nil, err
		}
	}

	viewsToDelete := make(map[string]*View)
	deletedIndices := make(map[string]map[int]bool)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		fpath, err := filter.aliases.Get(table.Name())
		if err != nil {
			return nil, nil, err
		}

		viewKey := strings.ToUpper(table.Name().Literal)
		if filter.tempViews.Exists(fpath) {
			viewsToDelete[viewKey], _ = filter.tempViews.Get(parser.Identifier{Literal: fpath})
		} else {
			viewsToDelete[viewKey], _ = filter.tx.cachedViews.Get(parser.Identifier{Literal: fpath})
		}
		if err = viewsToDelete[viewKey].Header.Update(table.Name().Literal, nil); err != nil {
			return nil, nil, err
		}
		deletedIndices[viewKey] = make(map[int]bool)
	}

	for i := range view.RecordSet {
		for viewref := range viewsToDelete {
			internalId, err := view.InternalRecordId(viewref, i)
			if err != nil {
				continue
			}
			if !deletedIndices[viewref][internalId] {
				deletedIndices[viewref][internalId] = true
			}
		}
	}

	fileInfos := make([]*FileInfo, 0)
	deletedCounts := make([]int, 0)
	for k, v := range viewsToDelete {
		records := make(RecordSet, 0, v.RecordLen()-len(deletedIndices[k]))
		for i, record := range v.RecordSet {
			if !deletedIndices[k][i] {
				records = append(records, record)
			}
		}
		v.RecordSet = records

		if err = v.RestoreHeaderReferences(); err != nil {
			return nil, nil, err
		}

		if v.FileInfo.IsTemporary {
			filter.tempViews.Replace(v)
		} else {
			if err = filter.tx.cachedViews.Replace(v); err != nil {
				return nil, nil, err
			}
		}

		fileInfos = append(fileInfos, v.FileInfo)
		deletedCounts = append(deletedCounts, len(deletedIndices[k]))
	}

	return fileInfos, deletedCounts, nil
}

func CreateTable(ctx context.Context, parentFilter *Filter, query parser.CreateTable) (*FileInfo, error) {
	filter := parentFilter.CreateNode()

	var view *View
	var err error

	flags := parentFilter.tx.Flags
	fileInfo, err := NewFileInfoForCreate(query.Table, flags.Repository, flags.WriteDelimiter, flags.WriteEncoding)
	if err != nil {
		return nil, err
	}
	h, err := file.NewHandlerForCreate(filter.tx.FileContainer, fileInfo.Path)
	if err != nil {
		return nil, NewFileAlreadyExistError(query.Table)
	}
	fileInfo.Handler = h

	fileInfo.LineBreak = flags.LineBreak
	fileInfo.EncloseAll = flags.EncloseAll
	fileInfo.NoHeader = flags.WithoutHeader
	fileInfo.PrettyPrint = flags.PrettyPrint
	fileInfo.ForUpdate = true

	if query.Query != nil {
		view, err = Select(ctx, filter, query.Query.(parser.SelectQuery))
		if err != nil {
			return nil, AppendCompositeError(err, filter.tx.FileContainer.Close(fileInfo.Handler))
		}

		if err = view.Header.Update(parser.FormatTableName(fileInfo.Path), query.Fields); err != nil {
			if _, ok := err.(*FieldLengthNotMatchError); ok {
				err = NewTableFieldLengthError(query.Query.(parser.SelectQuery), query.Table, len(query.Fields))
			}
			return nil, AppendCompositeError(err, filter.tx.FileContainer.Close(fileInfo.Handler))
		}
	} else {
		fields := make([]string, len(query.Fields))
		for i, v := range query.Fields {
			f, _ := v.(parser.Identifier)
			if InStrSliceWithCaseInsensitive(f.Literal, fields) {
				err = NewDuplicateFieldNameError(f)
				return nil, AppendCompositeError(err, filter.tx.FileContainer.Close(fileInfo.Handler))
			}
			fields[i] = f.Literal
		}
		header := NewHeader(parser.FormatTableName(fileInfo.Path), fields)
		view = &View{
			Header:    header,
			RecordSet: RecordSet{},
		}
	}

	view.FileInfo = fileInfo

	filter.tx.cachedViews.Set(view)

	return view.FileInfo, nil
}

func AddColumns(ctx context.Context, parentFilter *Filter, query parser.AddColumns) (*FileInfo, int, error) {
	filter := parentFilter.CreateNode()

	if query.Position == nil {
		query.Position = parser.ColumnPosition{
			Position: parser.Token{Token: parser.LAST, Literal: parser.TokenLiteral(parser.LAST)},
		}
	}

	view := NewView(parentFilter.tx)
	err := view.LoadFromTableIdentifier(ctx, filter, query.Table, true, false)
	if err != nil {
		return nil, 0, err
	}

	var insertPos int
	pos, _ := query.Position.(parser.ColumnPosition)
	switch pos.Position.Token {
	case parser.FIRST:
		insertPos = 0
	case parser.LAST:
		insertPos = view.FieldLen()
	default:
		idx, err := view.FieldIndex(pos.Column)
		if err != nil {
			return nil, 0, err
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
	defaults := make([]parser.QueryExpression, len(query.Columns))
	for i, coldef := range query.Columns {
		if InStrSliceWithCaseInsensitive(coldef.Column.Literal, columnNames) || InStrSliceWithCaseInsensitive(coldef.Column.Literal, fields) {
			return nil, 0, NewDuplicateFieldNameError(coldef.Column)
		}
		fields[i] = coldef.Column.Literal
		defaults[i] = coldef.Value
	}
	newFieldLen := view.FieldLen() + len(query.Columns)

	addHeader := NewHeader(parser.FormatTableName(view.FileInfo.Path), fields)
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

	records := make(RecordSet, view.RecordLen())

	err = NewFilterForSequentialEvaluation(filter, view).EvaluateSequentially(ctx, func(f *Filter, rIdx int) error {
		record := make(Record, newFieldLen)
		for i, cell := range view.RecordSet[rIdx] {
			var cellIdx int
			if i < insertPos {
				cellIdx = i
			} else {
				cellIdx = i + len(fields)
			}
			record[cellIdx] = cell
		}

		for i, v := range defaults {
			if v == nil {
				v = parser.NewNullValue()
			}
			val, e := f.Evaluate(ctx, v)
			if e != nil {
				return e
			}
			record[i+insertPos] = NewCell(val)
		}
		records[rIdx] = record
		return nil
	}, defaults)
	if err != nil {
		return nil, 0, err
	}

	view.Header = header
	view.RecordSet = records
	view.Filter = nil

	if view.FileInfo.IsTemporary {
		filter.tempViews.Replace(view)
	} else {
		err = filter.tx.cachedViews.Replace(view)
	}

	return view.FileInfo, len(fields), err
}

func DropColumns(ctx context.Context, parentFilter *Filter, query parser.DropColumns) (*FileInfo, int, error) {
	filter := parentFilter.CreateNode()

	view := NewView(parentFilter.tx)
	err := view.LoadFromTableIdentifier(ctx, filter, query.Table, true, false)
	if err != nil {
		return nil, 0, err
	}

	dropIndices := make([]int, len(query.Columns))
	for i, v := range query.Columns {
		idx, err := view.FieldIndex(v)
		if err != nil {
			return nil, 0, err
		}
		dropIndices[i] = idx
	}

	view.selectFields = []int{}
	for i := 0; i < view.FieldLen(); i++ {
		if view.Header[i].IsFromTable && !InIntSlice(i, dropIndices) {
			view.selectFields = append(view.selectFields, i)
		}
	}

	if err = view.Fix(ctx); err != nil {
		return nil, 0, err
	}

	if view.FileInfo.IsTemporary {
		filter.tempViews.Replace(view)
	} else {
		err = filter.tx.cachedViews.Replace(view)
	}

	return view.FileInfo, len(dropIndices), err

}

func RenameColumn(ctx context.Context, parentFilter *Filter, query parser.RenameColumn) (*FileInfo, error) {
	filter := parentFilter.CreateNode()

	view := NewView(parentFilter.tx)
	err := view.LoadFromTableIdentifier(ctx, filter, query.Table, true, false)
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
	view.Filter = nil

	if view.FileInfo.IsTemporary {
		filter.tempViews.Replace(view)
	} else {
		err = filter.tx.cachedViews.Replace(view)
	}

	return view.FileInfo, err
}

func SetTableAttribute(ctx context.Context, parentFilter *Filter, query parser.SetTableAttribute) (*FileInfo, string, error) {
	var log string
	filter := parentFilter.CreateNode()

	view := NewView(parentFilter.tx)
	err := view.LoadFromTableIdentifier(ctx, filter, query.Table, true, false)
	if err != nil {
		return nil, log, err
	}
	if view.FileInfo.IsTemporary {
		return nil, log, NewNotTableError(query.Table)
	}

	var p value.Primary
	if ident, ok := query.Value.(parser.Identifier); ok {
		p = value.NewString(ident.Literal)
	} else {
		p, err = filter.Evaluate(ctx, query.Value)
		if err != nil {
			return nil, log, err
		}
	}

	fileInfo := view.FileInfo
	attr := strings.ToUpper(query.Attribute.Literal)
	switch attr {
	case TableDelimiter, TableDelimiterPositions, TableFormat, TableEncoding, TableLineBreak, TableJsonEscape:
		s := value.ToString(p)
		if value.IsNull(s) {
			return nil, log, NewTableAttributeValueNotAllowedFormatError(query)
		}
		switch attr {
		case TableDelimiter:
			err = fileInfo.SetDelimiter(s.(value.String).Raw())
		case TableDelimiterPositions:
			err = fileInfo.SetDelimiterPositions(s.(value.String).Raw())
		case TableFormat:
			err = fileInfo.SetFormat(s.(value.String).Raw())
		case TableEncoding:
			err = fileInfo.SetEncoding(s.(value.String).Raw())
		case TableLineBreak:
			err = fileInfo.SetLineBreak(s.(value.String).Raw())
		case TableJsonEscape:
			err = fileInfo.SetJsonEscape(s.(value.String).Raw())
		}
	case TableHeader, TableEncloseAll, TablePrettyPrint:
		b := value.ToBoolean(p)
		if value.IsNull(b) {
			return nil, log, NewTableAttributeValueNotAllowedFormatError(query)
		}
		switch attr {
		case TableHeader:
			err = fileInfo.SetNoHeader(!b.(value.Boolean).Raw())
		case TableEncloseAll:
			err = fileInfo.SetEncloseAll(b.(value.Boolean).Raw())
		case TablePrettyPrint:
			err = fileInfo.SetPrettyPrint(b.(value.Boolean).Raw())
		}
	default:
		return nil, log, NewInvalidTableAttributeNameError(query.Attribute)
	}

	if err != nil {
		if _, ok := err.(*TableAttributeUnchangedError); ok {
			return nil, log, err
		}
		return nil, log, NewInvalidTableAttributeValueError(query, err.Error())
	}

	w := NewObjectWriter(filter.tx)
	w.WriteColorWithoutLineBreak("Path: ", cmd.LableEffect)
	w.WriteColorWithoutLineBreak(fileInfo.Path, cmd.ObjectEffect)
	w.NewLine()
	writeTableAttribute(w, parentFilter.tx.Flags, fileInfo)
	w.NewLine()

	w.Title1 = "Attributes Updated in"
	if i, ok := query.Table.(parser.Identifier); ok {
		w.Title2 = i.Literal
	} else if to, ok := query.Table.(parser.TableObject); ok {
		w.Title2 = to.Path.Literal
	}
	w.Title2Effect = cmd.IdentifierEffect
	log = "\n" + w.String() + "\n"

	err = filter.tx.cachedViews.Replace(view)
	return view.FileInfo, log, err
}
