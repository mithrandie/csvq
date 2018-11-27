package query

import (
	"fmt"
	"io"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func FetchCursor(name parser.Identifier, fetchPosition parser.FetchPosition, vars []parser.Variable, filter *Filter) (bool, error) {
	position := parser.NEXT
	number := -1
	if !fetchPosition.Position.IsEmpty() {
		position = fetchPosition.Position.Token
		if fetchPosition.Number != nil {
			p, err := filter.Evaluate(fetchPosition.Number)
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

	primaries, err := filter.Cursors.Fetch(name, position, number)
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
		_, err := filter.Variables.SubstituteDirectly(v, primaries[i])
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func DeclareView(expr parser.ViewDeclaration, filter *Filter) error {
	if filter.TempViews.Exists(expr.View.Literal) {
		return NewTemporaryTableRedeclaredError(expr.View)
	}

	var view *View
	var err error

	if expr.Query != nil {
		view, err = Select(expr.Query.(parser.SelectQuery), filter)
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
		view = &View{
			Header:    header,
			RecordSet: RecordSet{},
		}
	}

	view.FileInfo = &FileInfo{
		Path:             expr.View.Literal,
		IsTemporary:      true,
		InitialHeader:    view.Header.Copy(),
		InitialRecordSet: view.RecordSet.Copy(),
	}

	filter.TempViews.Set(view)

	return err
}

func Select(query parser.SelectQuery, parentFilter *Filter) (*View, error) {
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

func selectEntity(expr parser.QueryExpression, filter *Filter) (*View, error) {
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
	}

	if err := view.Select(entity.SelectClause.(parser.SelectClause)); err != nil {
		return nil, err
	}

	return view, nil
}

func selectSetEntity(expr parser.QueryExpression, filter *Filter) (*View, error) {
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

func selectSet(set parser.SelectSet, filter *Filter) (*View, error) {
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

func selectSetForRecursion(view *View, set parser.SelectSet, filter *Filter) error {
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

func Insert(query parser.InsertQuery, parentFilter *Filter) (*FileInfo, int, error) {
	filter := parentFilter.CreateNode()

	var insertRecords int

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, insertRecords, err
		}
	}

	fromClause := parser.FromClause{
		Tables: []parser.QueryExpression{
			query.Table,
		},
	}
	view := NewView()
	view.ForUpdate = true
	err := view.Load(fromClause, filter)
	if err != nil {
		return nil, insertRecords, err
	}

	fields := query.Fields
	if fields == nil {
		fields = view.Header.TableColumns()
	}

	if query.ValuesList != nil {
		if insertRecords, err = view.InsertValues(fields, query.ValuesList); err != nil {
			return nil, insertRecords, err
		}
	} else {
		if insertRecords, err = view.InsertFromQuery(fields, query.Query.(parser.SelectQuery)); err != nil {
			return nil, insertRecords, err
		}
	}

	view.RestoreHeaderReferences()
	view.Filter = nil

	if view.FileInfo.IsTemporary {
		filter.TempViews.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view.FileInfo, insertRecords, nil
}

func Update(query parser.UpdateQuery, parentFilter *Filter) ([]*FileInfo, []int, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, nil, err
		}
	}

	if query.FromClause == nil {
		query.FromClause = parser.FromClause{Tables: query.Tables}
	}

	view := NewView()
	view.ForUpdate = true
	view.UseInternalId = true
	err := view.Load(query.FromClause.(parser.FromClause), filter)
	if err != nil {
		return nil, nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, nil, err
		}
	}

	viewsToUpdate := make(map[string]*View)
	updatedCount := make(map[string]int)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		fpath, err := filter.Aliases.Get(table.Name())
		if err != nil {
			return nil, nil, err
		}
		viewKey := strings.ToUpper(table.Name().Literal)

		if filter.TempViews.Exists(fpath) {
			viewsToUpdate[viewKey], _ = filter.TempViews.Get(parser.Identifier{Literal: fpath})
		} else {
			viewsToUpdate[viewKey], _ = ViewCache.Get(parser.Identifier{Literal: fpath})
		}
		viewsToUpdate[viewKey].Header.Update(table.Name().Literal, nil)
	}

	updatesList := make(map[string]map[int][]int)
	filterForLoop := NewFilterForSequentialEvaluation(view, filter)
	for i := range view.RecordSet {
		filterForLoop.Records[0].RecordIndex = i
		internalIds := make(map[string]int)

		for _, uset := range query.SetList {
			val, err := filterForLoop.Evaluate(uset.Value)
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
		v.RestoreHeaderReferences()

		if v.FileInfo.IsTemporary {
			filter.TempViews.Replace(v)
		} else {
			ViewCache.Replace(v)
		}

		fileInfos = append(fileInfos, v.FileInfo)
		updateRecords = append(updateRecords, updatedCount[k])
	}

	return fileInfos, updateRecords, nil
}

func Delete(query parser.DeleteQuery, parentFilter *Filter) ([]*FileInfo, []int, error) {
	filter := parentFilter.CreateNode()

	if query.WithClause != nil {
		if err := filter.LoadInlineTable(query.WithClause.(parser.WithClause)); err != nil {
			return nil, nil, err
		}
	}

	fromClause := query.FromClause
	if query.Tables == nil {
		if 1 < len(fromClause.Tables) {
			return nil, nil, NewDeleteTableNotSpecifiedError(query)
		}
		table := fromClause.Tables[0].(parser.Table)
		if _, ok := table.Object.(parser.Identifier); !ok {
			if _, ok := table.Object.(parser.Stdin); !ok {
				return nil, nil, NewDeleteTableNotSpecifiedError(query)
			}
		}
		query.Tables = fromClause.Tables
	}

	view := NewView()
	view.UseInternalId = true
	view.ForUpdate = true
	err := view.Load(query.FromClause, filter)
	if err != nil {
		return nil, nil, err
	}

	if query.WhereClause != nil {
		if err := view.Where(query.WhereClause.(parser.WhereClause)); err != nil {
			return nil, nil, err
		}
	}

	viewsToDelete := make(map[string]*View)
	deletedIndices := make(map[string]map[int]bool)
	for _, v := range query.Tables {
		table := v.(parser.Table)
		fpath, err := filter.Aliases.Get(table.Name())
		if err != nil {
			return nil, nil, err
		}

		viewKey := strings.ToUpper(table.Name().Literal)
		if filter.TempViews.Exists(fpath) {
			viewsToDelete[viewKey], _ = filter.TempViews.Get(parser.Identifier{Literal: fpath})
		} else {
			viewsToDelete[viewKey], _ = ViewCache.Get(parser.Identifier{Literal: fpath})
		}
		viewsToDelete[viewKey].Header.Update(table.Name().Literal, nil)
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

		v.RestoreHeaderReferences()

		if v.FileInfo.IsTemporary {
			filter.TempViews.Replace(v)
		} else {
			ViewCache.Replace(v)
		}

		fileInfos = append(fileInfos, v.FileInfo)
		deletedCounts = append(deletedCounts, len(deletedIndices[k]))
	}

	return fileInfos, deletedCounts, nil
}

func CreateTable(query parser.CreateTable, parentFilter *Filter) (*FileInfo, error) {
	filter := parentFilter.CreateNode()

	var view *View
	var err error

	flags := cmd.GetFlags()
	fileInfo, err := NewFileInfoForCreate(query.Table, flags.Repository, flags.WriteDelimiter, flags.WriteEncoding)
	if err != nil {
		return nil, err
	}
	h, err := file.NewHandlerForCreate(fileInfo.Path)
	if err != nil {
		return nil, NewFileAlreadyExistError(query.Table)
	}
	fileInfo.Handler = h

	fileInfo.LineBreak = flags.LineBreak
	fileInfo.EncloseAll = flags.EncloseAll
	fileInfo.NoHeader = flags.WithoutHeader
	fileInfo.PrettyPrint = flags.PrettyPrint

	if query.Query != nil {
		view, err = Select(query.Query.(parser.SelectQuery), filter)
		if err != nil {
			fileInfo.Close()
			return nil, err
		}

		if err = view.Header.Update(parser.FormatTableName(fileInfo.Path), query.Fields); err != nil {
			fileInfo.Close()
			if _, ok := err.(*FieldLengthNotMatchError); ok {
				return nil, NewTableFieldLengthError(query.Query.(parser.SelectQuery), query.Table, len(query.Fields))
			}
			return nil, err
		}
	} else {
		fields := make([]string, len(query.Fields))
		for i, v := range query.Fields {
			f, _ := v.(parser.Identifier)
			if InStrSliceWithCaseInsensitive(f.Literal, fields) {
				fileInfo.Close()
				return nil, NewDuplicateFieldNameError(f)
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
	view.ForUpdate = true

	ViewCache.Set(view)

	return view.FileInfo, nil
}

func AddColumns(query parser.AddColumns, parentFilter *Filter) (*FileInfo, int, error) {
	filter := parentFilter.CreateNode()

	if query.Position == nil {
		query.Position = parser.ColumnPosition{
			Position: parser.Token{Token: parser.LAST, Literal: parser.TokenLiteral(parser.LAST)},
		}
	}

	view := NewView()
	view.ForUpdate = true
	err := view.LoadFromTableIdentifier(query.Table, filter)
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

	err = NewFilterForSequentialEvaluation(view, filter).EvaluateSequentially(defaults, func(f *Filter, startIdx int) error {
		idx := f.currentIndex() + startIdx

		record := make(Record, newFieldLen)
		for i, cell := range view.RecordSet[idx] {
			var idx int
			if i < insertPos {
				idx = i
			} else {
				idx = i + len(fields)
			}
			record[idx] = cell
		}

		for j, v := range defaults {
			if v == nil {
				v = parser.NewNullValue()
			}
			val, e := f.Evaluate(v)
			if e != nil {
				return e
			}
			record[j+insertPos] = NewCell(val)
		}
		records[idx] = record
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	view.Header = header
	view.RecordSet = records
	view.Filter = nil

	if view.FileInfo.IsTemporary {
		filter.TempViews.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view.FileInfo, len(fields), nil
}

func DropColumns(query parser.DropColumns, parentFilter *Filter) (*FileInfo, int, error) {
	filter := parentFilter.CreateNode()

	view := NewView()
	view.ForUpdate = true
	err := view.LoadFromTableIdentifier(query.Table, filter)
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

	view.Fix()

	if view.FileInfo.IsTemporary {
		filter.TempViews.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view.FileInfo, len(dropIndices), nil

}

func RenameColumn(query parser.RenameColumn, parentFilter *Filter) (*FileInfo, error) {
	filter := parentFilter.CreateNode()

	view := NewView()
	view.ForUpdate = true
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
	view.Filter = nil

	if view.FileInfo.IsTemporary {
		filter.TempViews.Replace(view)
	} else {
		ViewCache.Replace(view)
	}

	return view.FileInfo, nil
}

func SetTableAttribute(query parser.SetTableAttribute, parentFilter *Filter) (*FileInfo, string, error) {
	var log string
	filter := parentFilter.CreateNode()

	view := NewView()
	view.ForUpdate = true
	err := view.LoadFromTableIdentifier(query.Table, filter)
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
		p, err = filter.Evaluate(query.Value)
		if err != nil {
			return nil, log, err
		}
	}

	fileInfo := view.FileInfo
	attr := strings.ToUpper(query.Attribute.Literal)
	switch attr {
	case TableDelimiter, TableFormat, TableEncoding, TableLineBreak, TableJsonEscape:
		s := value.ToString(p)
		if value.IsNull(s) {
			return nil, log, NewTableAttributeValueNotAllowedFormatError(query)
		}
		switch attr {
		case TableDelimiter:
			err = fileInfo.SetDelimiter(s.(value.String).Raw())
		case TableFormat:
			err = fileInfo.SetFormat(s.(value.String).Raw())
		case TableEncoding:
			err = fileInfo.SetEncoding(s.(value.String).Raw())
		case TableLineBreak:
			err = fileInfo.SetLineBreak(s.(value.String).Raw())
		case TableJsonEscape:
			err = fileInfo.SetJsonEscape(s.(value.String).Raw())
		}
	case TableHeader, TableEncloseAll, TablePrettyPring:
		b := value.ToBoolean(p)
		if value.IsNull(b) {
			return nil, log, NewTableAttributeValueNotAllowedFormatError(query)
		}
		switch attr {
		case TableHeader:
			err = fileInfo.SetNoHeader(!b.(value.Boolean).Raw())
		case TableEncloseAll:
			err = fileInfo.SetEncloseAll(b.(value.Boolean).Raw())
		case TablePrettyPring:
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

	w := NewObjectWriter()
	w.WriteColorWithoutLineBreak("Path: ", cmd.LableEffect)
	w.WriteColorWithoutLineBreak(fileInfo.Path, cmd.ObjectEffect)
	w.NewLine()
	writeTableAttribute(w, fileInfo)
	w.NewLine()

	w.Title1 = "Attributes Updated in"
	w.Title2 = query.Table.(parser.Identifier).Literal
	w.Title2Effect = cmd.IdentifierEffect
	log = "\n" + w.String() + "\n"

	ViewCache.Replace(view)

	return view.FileInfo, log, nil
}

func Commit(expr parser.Expression, filter *Filter) error {
	createdFiles, updatedFiles := UncommittedViews.UncommittedFiles()

	createFileInfo := make([]*FileInfo, 0, len(createdFiles))
	updateFileInfo := make([]*FileInfo, 0, len(updatedFiles))

	if 0 < len(createdFiles) {
		for _, fileinfo := range createdFiles {
			view, _ := ViewCache.Get(parser.Identifier{Literal: fileinfo.Path})

			fp := view.FileInfo.Handler.FileForUpdate()
			fp.Truncate(0)
			fp.Seek(0, io.SeekStart)

			err := EncodeView(fp, view, fileinfo)
			if err != nil {
				return NewCommitError(expr, err.Error())
			}
			createFileInfo = append(createFileInfo, view.FileInfo)
		}
	}

	if 0 < len(updatedFiles) {
		for _, fileinfo := range updatedFiles {
			view, _ := ViewCache.Get(parser.Identifier{Literal: fileinfo.Path})

			fp := view.FileInfo.Handler.FileForUpdate()
			fp.Truncate(0)
			fp.Seek(0, io.SeekStart)

			if err := EncodeView(fp, view, fileinfo); err != nil {
				return NewCommitError(expr, err.Error())
			}

			updateFileInfo = append(updateFileInfo, view.FileInfo)
		}
	}

	for _, f := range createFileInfo {
		if err := f.Commit(); err != nil {
			return NewCommitError(expr, err.Error())
		}
		UncommittedViews.Unset(f)
		Log(cmd.Notice(fmt.Sprintf("Commit: file %q is created.", f.Path)), cmd.GetFlags().Quiet)
	}
	for _, f := range updateFileInfo {
		if err := f.Commit(); err != nil {
			return NewCommitError(expr, err.Error())
		}
		UncommittedViews.Unset(f)
		Log(cmd.Notice(fmt.Sprintf("Commit: file %q is updated.", f.Path)), cmd.GetFlags().Quiet)
	}

	filter.TempViews.Store(UncommittedViews.UncommittedTempViews())
	UncommittedViews.Clean()
	if err := ReleaseResources(); err != nil {
		return NewCommitError(expr, err.Error())
	}
	return nil
}

func Rollback(expr parser.Expression, filter *Filter) error {
	createdFiles, updatedFiles := UncommittedViews.UncommittedFiles()

	if 0 < len(createdFiles) {
		for _, fileinfo := range createdFiles {
			Log(cmd.Notice(fmt.Sprintf("Rollback: file %q is deleted.", fileinfo.Path)), cmd.GetFlags().Quiet)
		}
	}

	if 0 < len(updatedFiles) {
		for _, fileinfo := range updatedFiles {
			Log(cmd.Notice(fmt.Sprintf("Rollback: file %q is restored.", fileinfo.Path)), cmd.GetFlags().Quiet)
		}
	}

	if filter != nil {
		filter.TempViews.Restore(UncommittedViews.UncommittedTempViews())
	}
	UncommittedViews.Clean()
	if err := ReleaseResources(); err != nil {
		return NewRollbackError(expr, err.Error())
	}
	return nil
}
