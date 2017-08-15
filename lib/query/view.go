package query

import (
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type View struct {
	Header  Header
	Records Records

	FileInfo *FileInfo

	selectFields []int
	selectLabels []string
	isGrouped    bool
	Filter       Filter

	filteredIndices []int

	sortIndices       []int
	sortDirections    []int
	sortNullPositions []int

	offset int

	OperatedRecords int
	OperatedFields  int

	UseInternalId bool
}

func NewView() *View {
	return &View{
		UseInternalId: false,
	}
}

func (view *View) Load(clause parser.FromClause, filter Filter) error {
	if clause.Tables == nil {
		var obj parser.Expression
		if IsReadableFromStdin() {
			obj = parser.Stdin{Stdin: "stdin"}
		} else {
			obj = parser.Dual{}
		}
		clause.Tables = []parser.Expression{parser.Table{Object: obj}}
	}

	views := make([]*View, len(clause.Tables))
	for i, v := range clause.Tables {
		loaded, err := loadView(v.(parser.Table), filter, view.UseInternalId)
		if err != nil {
			return err
		}
		views[i] = loaded
	}

	view.Header = views[0].Header
	view.Records = views[0].Records
	view.FileInfo = views[0].FileInfo

	for i := 1; i < len(views); i++ {
		CrossJoin(view, views[i])
	}

	view.Filter = filter
	return nil
}

func (view *View) LoadFromTableIdentifier(table parser.Expression, filter Filter) error {
	fromClause := parser.FromClause{
		Tables: []parser.Expression{
			parser.Table{Object: table},
		},
	}

	return view.Load(fromClause, filter)
}

func loadView(table parser.Table, filter Filter, useInternalId bool) (*View, error) {
	var view *View
	var err error

	switch table.Object.(type) {
	case parser.Dual:
		view = loadDualView()
	case parser.Stdin:
		delimiter := cmd.GetFlags().Delimiter
		if delimiter == cmd.UNDEF {
			delimiter = ','
		}
		fileInfo := &FileInfo{
			Path:      table.Object.String(),
			Delimiter: delimiter,
			Temporary: true,
		}

		if !filter.TempViewsList[len(filter.TempViewsList)-1].Exists(fileInfo.Path) {
			if !IsReadableFromStdin() {
				return nil, NewStdinEmptyError(table.Object.(parser.Stdin))
			}

			file := os.Stdin
			defer file.Close()

			loadView, err := loadViewFromFile(file, fileInfo)
			if err != nil {
				return nil, err
			}
			loadView.FileInfo.InitialRecords = loadView.Records.Copy()
			filter.TempViewsList[len(filter.TempViewsList)-1].Set(loadView)
		}
		if err = filter.AliasesList.Add(table.Name(), fileInfo.Path); err != nil {
			return nil, err
		}

		pathIdent := parser.Identifier{Literal: table.Object.String()}
		if useInternalId {
			view, _ = filter.TempViewsList[len(filter.TempViewsList)-1].GetWithInternalId(pathIdent)
		} else {
			view, _ = filter.TempViewsList[len(filter.TempViewsList)-1].Get(pathIdent)
		}
		if !strings.EqualFold(table.Object.String(), table.Name().Literal) {
			view.Header.Update(table.Name().Literal, nil)
		}
	case parser.Identifier:
		tableIdentifier := table.Object.(parser.Identifier)
		if filter.RecursiveTable != nil && strings.EqualFold(tableIdentifier.Literal, filter.RecursiveTable.Name.Literal) && filter.RecursiveTmpView != nil {
			view = filter.RecursiveTmpView
			if !strings.EqualFold(filter.RecursiveTable.Name.Literal, table.Name().Literal) {
				view.Header.Update(table.Name().Literal, nil)
			}
		} else if ct, err := filter.InlineTablesList.Get(tableIdentifier); err == nil {
			if err = filter.AliasesList.Add(table.Name(), ""); err != nil {
				return nil, err
			}
			view = ct
			if !strings.EqualFold(tableIdentifier.Literal, table.Name().Literal) {
				view.Header.Update(table.Name().Literal, nil)
			}
		} else {
			var fileInfo *FileInfo
			var commonTableName string

			if filter.TempViewsList.Exists(tableIdentifier.Literal) {
				fileInfo = &FileInfo{
					Path: tableIdentifier.Literal,
				}

				commonTableName = parser.FormatTableName(fileInfo.Path)

				pathIdent := parser.Identifier{Literal: fileInfo.Path}
				if useInternalId {
					view, _ = filter.TempViewsList.GetWithInternalId(pathIdent)
				} else {
					view, _ = filter.TempViewsList.Get(pathIdent)
				}
			} else {
				flags := cmd.GetFlags()

				fileInfo, err = NewFileInfoForCreate(tableIdentifier, flags.Repository, flags.Delimiter)
				if err != nil {
					return nil, err
				}

				if !ViewCache.Exists(fileInfo.Path) {
					fileInfo, err = NewFileInfo(tableIdentifier, flags.Repository, flags.Delimiter)
					if err != nil {
						return nil, err
					}

					if !ViewCache.Exists(fileInfo.Path) {
						file, err := os.Open(fileInfo.Path)
						if err != nil {
							return nil, NewReadFileError(tableIdentifier, err.Error())
						}
						defer file.Close()
						loadView, err := loadViewFromFile(file, fileInfo)
						if err != nil {
							return nil, NewCsvParsingError(tableIdentifier, fileInfo.Path, err.Error())
						}
						ViewCache.Set(loadView)
					}
				}
				commonTableName = parser.FormatTableName(fileInfo.Path)

				pathIdent := parser.Identifier{Literal: fileInfo.Path}
				if useInternalId {
					view, _ = ViewCache.GetWithInternalId(pathIdent)
				} else {
					view, _ = ViewCache.Get(pathIdent)
				}
			}

			if err = filter.AliasesList.Add(table.Name(), fileInfo.Path); err != nil {
				return nil, err
			}

			if !strings.EqualFold(commonTableName, table.Name().Literal) {
				view.Header.Update(table.Name().Literal, nil)
			}
		}
	case parser.Join:
		join := table.Object.(parser.Join)
		view, err = loadView(join.Table, filter, useInternalId)
		if err != nil {
			return nil, err
		}
		view2, err := loadView(join.JoinTable, filter, useInternalId)
		if err != nil {
			return nil, err
		}

		condition := ParseJoinCondition(join, view, view2)

		joinType := join.JoinType.Token
		if join.JoinType.IsEmpty() {
			if join.Direction.IsEmpty() {
				joinType = parser.INNER
			} else {
				joinType = parser.OUTER
			}
		}

		switch joinType {
		case parser.CROSS:
			CrossJoin(view, view2)
		case parser.INNER:
			if err = InnerJoin(view, view2, condition, filter); err != nil {
				return nil, err
			}
		case parser.OUTER:
			if err = OuterJoin(view, view2, condition, join.Direction.Token, filter); err != nil {
				return nil, err
			}
		}
	case parser.Subquery:
		subquery := table.Object.(parser.Subquery)
		view, err = Select(subquery.Query, filter)
		if table.Alias != nil {
			if err = filter.AliasesList.Add(table.Alias.(parser.Identifier), ""); err != nil {
				return nil, err
			}
		}
		if err == nil {
			view.Header.Update(table.Name().Literal, nil)
		}
	}

	return view, err
}

func loadViewFromFile(file *os.File, fileInfo *FileInfo) (*View, error) {
	flags := cmd.GetFlags()

	r := cmd.GetReader(file, flags.Encoding)

	view := new(View)

	reader := csv.NewReader(r)
	reader.Delimiter = fileInfo.Delimiter
	reader.WithoutNull = flags.WithoutNull

	var err error
	var header []string
	if !flags.NoHeader {
		header, err = reader.ReadHeader()
		if err != nil && err != csv.EOF {
			return nil, err
		}
	}

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	view.Records = make([]Record, len(rows))
	for i, v := range rows {
		view.Records[i] = NewRecord(v)
	}

	if header == nil {
		header = make([]string, reader.FieldsPerRecords)
		for i := 0; i < reader.FieldsPerRecords; i++ {
			header[i] = "c" + strconv.Itoa(i+1)
		}
	}

	fileInfo.NoHeader = flags.NoHeader
	fileInfo.Encoding = flags.Encoding
	fileInfo.LineBreak = reader.LineBreak
	if fileInfo.LineBreak == "" {
		fileInfo.LineBreak = flags.LineBreak
	}

	view.Header = NewHeader(parser.FormatTableName(fileInfo.Path), header)
	view.FileInfo = fileInfo
	return view, nil
}

func loadDualView() *View {
	view := View{
		Header:  NewDualHeader(),
		Records: make([]Record, 1),
	}
	view.Records[0] = NewEmptyRecord(1)
	return &view
}

func NewViewFromGroupedRecord(filterRecord FilterRecord) *View {
	view := new(View)
	view.Header = filterRecord.View.Header
	record := filterRecord.View.Records[filterRecord.RecordIndex]

	view.Records = make([]Record, record.GroupLen())
	for i := 0; i < record.GroupLen(); i++ {
		view.Records[i] = make(Record, view.FieldLen())
		for j, cell := range record {
			grpIdx := i
			if cell.Len() < 2 {
				grpIdx = 0
			}
			view.Records[i][j] = NewCell(cell.GroupedPrimary(grpIdx))
		}
	}

	return view
}

func (view *View) Where(clause parser.WhereClause) error {
	indices, err := view.filter(clause.Filter)
	if err != nil {
		return err
	}

	view.filteredIndices = indices
	return nil
}

func (view *View) filter(condition parser.Expression) ([]int, error) {
	indices := []int{}
	for i := range view.Records {
		filter := NewFilterForRecord(view, i, view.Filter)
		primary, err := filter.Evaluate(condition)
		if err != nil {
			return nil, err
		}
		if primary.Ternary() == ternary.TRUE {
			indices = append(indices, i)
		}
	}
	return indices, nil
}

func (view *View) Extract() {
	records := make([]Record, len(view.filteredIndices))
	for i, idx := range view.filteredIndices {
		records[i] = view.Records[idx]
	}
	view.Records = records
	view.filteredIndices = nil
}

func (view *View) GroupBy(clause parser.GroupByClause) error {
	return view.group(clause.Items)
}

func (view *View) group(items []parser.Expression) error {
	if len(view.Records) < 1 {
		return nil
	}

	type group struct {
		primaries []parser.Primary
		indices   []int
	}

	var match = func(g group, keys []parser.Primary) bool {
		for i, primary := range g.primaries {
			if EquivalentTo(primary, keys[i]) != ternary.TRUE {
				return false
			}
		}
		return true
	}

	var groups []group

	for i := 0; i < view.RecordLen(); i++ {
		filter := NewFilterForRecord(view, i, view.Filter)

		keys := make([]parser.Primary, len(items))
		for j, item := range items {
			key, err := filter.Evaluate(item)
			if err != nil {
				return err
			}
			keys[j] = key
		}

		newGr := true
		for j, gr := range groups {
			if match(gr, keys) {
				groups[j].indices = append(groups[j].indices, i)
				newGr = false
				break
			}
		}
		if newGr {
			gr := group{
				primaries: keys,
				indices:   []int{i},
			}
			groups = append(groups, gr)
		}
	}

	records := make([]Record, len(groups))
	for i, gr := range groups {
		record := make(Record, view.FieldLen())
		for j := 0; j < view.FieldLen(); j++ {
			primaries := make([]parser.Primary, len(gr.indices))
			for k, idx := range gr.indices {
				primaries[k] = view.Records[idx][j].Primary()
			}
			record[j] = NewGroupCell(primaries)
		}

		records[i] = record
	}

	view.Records = records
	view.isGrouped = true
	for _, item := range items {
		if fieldRef, ok := item.(parser.FieldReference); ok {
			idx, _ := view.FieldIndex(fieldRef)
			view.Header[idx].IsGroupKey = true
		}
	}
	return nil
}

func (view *View) Having(clause parser.HavingClause) error {
	indices, err := view.filter(clause.Filter)
	if err != nil {
		if _, ok := err.(*NotGroupingRecordsError); ok {
			view.group(nil)
			indices, err = view.filter(clause.Filter)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	view.filteredIndices = indices
	return nil
}

func (view *View) Select(clause parser.SelectClause) error {
	var parseAllColumns = func(view *View, fields []parser.Expression) []parser.Expression {
		insertIdx := -1

		for i, field := range fields {
			if _, ok := field.(parser.Field).Object.(parser.AllColumns); ok {
				insertIdx = i
				break
			}
		}

		if insertIdx < 0 {
			return fields
		}

		columns := view.Header.TableColumns()
		insertLen := len(columns)
		insert := make([]parser.Expression, insertLen)
		for i, c := range columns {
			insert[i] = parser.Field{
				Object: c,
			}
		}

		list := make([]parser.Expression, len(fields)-1+insertLen)
		for i, field := range fields {
			switch {
			case i == insertIdx:
				continue
			case i < insertIdx:
				list[i] = field
			default:
				list[i+insertLen-1] = field
			}
		}
		for i, field := range insert {
			list[i+insertIdx] = field
		}

		return list
	}

	var evalFields = func(view *View, fields []parser.Expression) error {
		view.selectFields = make([]int, len(fields))
		view.selectLabels = make([]string, len(fields))
		for i, f := range fields {
			field := f.(parser.Field)
			label := field.Name()
			idx, err := view.evalColumn(field.Object, field.Object.String(), label)
			if err != nil {
				return err
			}
			view.selectFields[i] = idx
			view.selectLabels[i] = label
		}
		return nil
	}

	fields := parseAllColumns(view, clause.Fields)

	origFieldLen := view.FieldLen()
	err := evalFields(view, fields)
	if err != nil {
		if _, ok := err.(*NotGroupingRecordsError); ok {
			view.Header = view.Header[:origFieldLen]
			for i := range view.Records {
				view.Records[i] = view.Records[i][:origFieldLen]
			}

			view.group(nil)
			err = evalFields(view, fields)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if clause.IsDistinct() {
		records := make(Records, view.RecordLen())
		for i, v := range view.Records {
			record := make(Record, len(view.selectFields))
			for j, idx := range view.selectFields {
				record[j] = v[idx]
			}
			records[i] = record
		}

		hfields := NewEmptyHeader(len(view.selectFields))
		for i, idx := range view.selectFields {
			hfields[i] = view.Header[idx]
			view.selectFields[i] = i
		}

		view.Header = hfields
		view.Records = records

		view.Distinct()
	}

	return nil
}

func (view *View) Distinct() {
	distinguished := Records{}

	for _, record := range view.Records {
		if !distinguished.Contains(record) {
			distinguished = append(distinguished, record)
		}
	}

	view.Records = distinguished
}

func (view *View) SelectAllColumns() error {
	selectClause := parser.SelectClause{
		Fields: []parser.Expression{
			parser.Field{Object: parser.AllColumns{}},
		},
	}
	return view.Select(selectClause)
}

func (view *View) OrderBy(clause parser.OrderByClause) error {
	view.sortIndices = make([]int, len(clause.Items))
	view.sortDirections = make([]int, len(clause.Items))
	view.sortNullPositions = make([]int, len(clause.Items))

	for i, v := range clause.Items {
		oi := v.(parser.OrderItem)
		idx, err := view.evalColumn(oi.Value, oi.Value.String(), "")
		if err != nil {
			return err
		}
		view.sortIndices[i] = idx

		if oi.Direction.IsEmpty() {
			view.sortDirections[i] = parser.ASC
		} else {
			view.sortDirections[i] = oi.Direction.Token
		}

		if oi.Position.IsEmpty() {
			switch view.sortDirections[i] {
			case parser.ASC:
				view.sortNullPositions[i] = parser.FIRST
			default: //parser.DESC
				view.sortNullPositions[i] = parser.LAST
			}
		} else {
			view.sortNullPositions[i] = oi.Position.Token
		}
	}

	sort.Sort(view)
	return nil
}

func (view *View) evalColumn(obj parser.Expression, column string, alias string) (idx int, err error) {
	fieldInTable := false

	if fr, ok := obj.(parser.FieldReference); ok {
		if idx, err = view.FieldIndex(fr); err == nil {
			if view.isGrouped && view.Header[idx].FromTable && !view.Header[idx].IsGroupKey {
				err = NewFieldNotGroupKeyError(fr)
				return
			}
			fieldInTable = true
		}
	} else if cn, ok := obj.(parser.ColumnNumber); ok {
		if idx, err = view.Header.ContainsNumber(cn); err == nil {
			if view.isGrouped && view.Header[idx].FromTable && !view.Header[idx].IsGroupKey {
				err = NewFieldNumberNotGroupKeyError(cn)
				return
			}
			fieldInTable = true
		}
	}

	if !fieldInTable {
		idx, err = view.Header.ContainsObject(obj)
		if err != nil {
			err = nil

			if analyticFunction, ok := obj.(parser.AnalyticFunction); ok {
				err = view.evalAnalyticFunction(analyticFunction)
				if err != nil {
					return
				}
			} else {
				filter := NewFilterForSequentialEvaluation(view, view.Filter)
				for i := range view.Records {
					var primary parser.Primary
					filter.Records[0].RecordIndex = i

					primary, err = filter.Evaluate(obj)
					if err != nil {
						return
					}
					view.Records[i] = append(view.Records[i], NewCell(primary))
				}
			}
			view.Header, idx = AddHeaderField(view.Header, column, alias)
		}
	}

	if 0 < len(alias) {
		if !strings.EqualFold(view.Header[idx].Column, alias) && !InStrSliceWithCaseInsensitive(alias, view.Header[idx].Aliases) {
			view.Header[idx].Aliases = append(view.Header[idx].Aliases, alias)
		}
	}

	return
}

func (view *View) evalAnalyticFunction(expr parser.AnalyticFunction) error {
	name := strings.ToUpper(expr.Name)
	if _, ok := AggregateFunctions[name]; !ok {
		if _, ok := AnalyticFunctions[name]; !ok {
			if udfn, err := view.Filter.FunctionsList.Get(expr, expr.Name); err != nil || !udfn.IsAggregate {
				return NewFunctionNotExistError(expr, expr.Name)
			}
		}
	}

	if expr.AnalyticClause.OrderByClause != nil {
		err := view.OrderBy(expr.AnalyticClause.OrderByClause.(parser.OrderByClause))
		if err != nil {
			return err
		}
	}

	return Analyze(view, expr)
}

func (view *View) Offset(clause parser.OffsetClause) error {
	value, err := view.Filter.Evaluate(clause.Value)
	if err != nil {
		return err
	}
	number := parser.PrimaryToInteger(value)
	if parser.IsNull(number) {
		return NewInvalidOffsetNumberError(clause)
	}
	view.offset = int(number.(parser.Integer).Value())
	if view.offset < 0 {
		view.offset = 0
	}

	if view.RecordLen() <= view.offset {
		view.Records = Records{}
	} else {
		view.Records = view.Records[view.offset:]
		records := make(Records, len(view.Records))
		copy(records, view.Records)
		view.Records = records
	}
	return nil
}

func (view *View) Limit(clause parser.LimitClause) error {
	value, err := view.Filter.Evaluate(clause.Value)
	if err != nil {
		return err
	}

	var limit int
	if clause.IsPercentage() {
		number := parser.PrimaryToFloat(value)
		if parser.IsNull(number) {
			return NewInvalidLimitPercentageError(clause)
		}
		percentage := number.(parser.Float).Value()
		if 100 < percentage {
			limit = 100
		} else if percentage < 0 {
			limit = 0
		} else {
			limit = int(math.Ceil(float64(view.RecordLen()+view.offset) * percentage / 100))
		}
	} else {
		number := parser.PrimaryToInteger(value)
		if parser.IsNull(number) {
			return NewInvalidLimitNumberError(clause)
		}
		limit = int(number.(parser.Integer).Value())
		if limit < 0 {
			limit = 0
		}
	}

	if view.RecordLen() <= limit {
		return nil
	}

	if clause.IsWithTies() && 0 < len(view.sortIndices) {
		bottomRecord := view.Records[limit-1]
		for limit < view.RecordLen() {
			if !bottomRecord.Match(view.Records[limit], view.sortIndices) {
				break
			}
			limit++
		}
	}

	view.Records = view.Records[:limit]
	records := make(Records, view.RecordLen())
	copy(records, view.Records)
	view.Records = records
	return nil
}

func (view *View) InsertValues(fields []parser.Expression, list []parser.Expression, filter Filter) error {
	valuesList := make([][]parser.Primary, len(list))

	for i, item := range list {
		rv := item.(parser.RowValue)
		values, err := filter.evalRowValue(rv)
		if err != nil {
			return err
		}
		if len(fields) != len(values) {
			return NewInsertRowValueLengthError(rv, len(fields))
		}

		valuesList[i] = values
	}

	return view.insert(fields, valuesList)
}

func (view *View) InsertFromQuery(fields []parser.Expression, query parser.SelectQuery, filter Filter) error {
	insertView, err := Select(query, filter)
	if err != nil {
		return err
	}
	if len(fields) != insertView.FieldLen() {
		return NewInsertSelectFieldLengthError(query, len(fields))
	}

	valuesList := make([][]parser.Primary, insertView.RecordLen())

	for i, record := range insertView.Records {
		values := make([]parser.Primary, insertView.FieldLen())
		for j, cell := range record {
			values[j] = cell.Primary()
		}
		valuesList[i] = values
	}

	return view.insert(fields, valuesList)
}

func (view *View) insert(fields []parser.Expression, valuesList [][]parser.Primary) error {
	var valueIndex = func(i int, list []int) int {
		for j, v := range list {
			if i == v {
				return j
			}
		}
		return -1
	}

	fieldIndices, err := view.FieldIndices(fields)
	if err != nil {
		return err
	}

	records := make([]Record, len(valuesList))
	for i, values := range valuesList {
		record := make(Record, view.FieldLen())
		for j := 0; j < view.FieldLen(); j++ {
			idx := valueIndex(j, fieldIndices)
			if idx < 0 {
				record[j] = NewCell(parser.NewNull())
			} else {
				record[j] = NewCell(values[idx])
			}
		}
		records[i] = record
	}

	view.Records = append(view.Records, records...)
	view.OperatedRecords = len(valuesList)
	return nil
}

func (view *View) Fix() {
	hfields := NewEmptyHeader(len(view.selectFields))
	records := make([]Record, view.RecordLen())

	for i, v := range view.Records {
		record := make(Record, len(view.selectFields))
		for j, idx := range view.selectFields {
			if 1 < v.GroupLen() {
				record[j] = NewCell(v[idx].Primary())
			} else {
				record[j] = v[idx]
			}
		}

		records[i] = record
	}

	colNumber := 0
	for i, idx := range view.selectFields {
		colNumber++

		hfields[i] = view.Header[idx]
		hfields[i].Aliases = nil
		hfields[i].Number = colNumber
		hfields[i].FromTable = true
		hfields[i].IsGroupKey = false

		if 0 < len(view.selectLabels) {
			hfields[i].Column = view.selectLabels[i]
		}
	}

	view.Header = hfields
	view.Records = records
	view.selectFields = nil
	view.selectLabels = nil
	view.isGrouped = false
	view.Filter = Filter{}
	view.sortIndices = nil
	view.sortDirections = nil
	view.sortNullPositions = nil
	view.offset = 0
}

func (view *View) Union(calcView *View, all bool) {
	view.Records = append(view.Records, calcView.Records...)
	view.FileInfo = nil

	if !all {
		view.Distinct()
	}
}

func (view *View) Except(calcView *View, all bool) {
	indices := []int{}
	for i, record := range view.Records {
		if !calcView.Records.Contains(record) {
			indices = append(indices, i)
		}
	}
	view.filteredIndices = indices
	view.Extract()

	view.FileInfo = nil

	if !all {
		view.Distinct()
	}
}

func (view *View) Intersect(calcView *View, all bool) {
	indices := []int{}
	for i, record := range view.Records {
		if calcView.Records.Contains(record) {
			indices = append(indices, i)
		}
	}
	view.filteredIndices = indices
	view.Extract()

	view.FileInfo = nil

	if !all {
		view.Distinct()
	}
}

func (view *View) ListValuesForAggregateFunctions(expr parser.Expression, arg parser.Expression, distinct bool, filter Filter) ([]parser.Primary, error) {
	list := make([]parser.Primary, view.RecordLen())
	f := NewFilterForSequentialEvaluation(view, filter)
	for i := 0; i < view.RecordLen(); i++ {
		f.Records[0].RecordIndex = i
		p, err := f.Evaluate(arg)
		if err != nil {
			if _, ok := err.(*NotGroupingRecordsError); ok {
				err = NewNestedAggregateFunctionsError(expr)
			}
			return nil, err
		}
		list[i] = p
	}

	if distinct {
		list = Distinguish(list)
	}

	return list, nil
}

func (view *View) ListValuesForAnalyticFunctions(fn parser.AnalyticFunction, partitionItems PartitionItemList) ([]parser.Primary, error) {
	list := make([]parser.Primary, len(partitionItems))
	f := NewFilterForSequentialEvaluation(view, view.Filter)
	for i, item := range partitionItems {
		f.Records[0].RecordIndex = item.RecordIndex
		value, err := f.Evaluate(fn.Args[0])
		if err != nil {
			return nil, err
		}
		list[i] = value
	}

	if fn.IsDistinct() {
		list = Distinguish(list)
	}

	return list, nil
}

func (view *View) RestoreHeaderReferences() {
	view.Header.Update(parser.FormatTableName(view.FileInfo.Path), nil)
}

func (view *View) FieldIndex(fieldRef parser.Expression) (int, error) {
	if number, ok := fieldRef.(parser.ColumnNumber); ok {
		return view.Header.ContainsNumber(number)
	}
	return view.Header.Contains(fieldRef.(parser.FieldReference))
}

func (view *View) FieldIndices(fields []parser.Expression) ([]int, error) {
	indices := make([]int, len(fields))
	for i, v := range fields {
		idx, err := view.FieldIndex(v)
		if err != nil {
			return nil, err
		}
		indices[i] = idx
	}
	return indices, nil
}

func (view *View) FieldViewName(fieldRef parser.Expression) (string, error) {
	idx, err := view.FieldIndex(fieldRef)
	if err != nil {
		return "", err
	}
	return view.Header[idx].View, nil
}

func (view *View) InternalRecordId(ref string, recordIndex int) (int, error) {
	fieldRef := parser.FieldReference{
		View:   parser.Identifier{Literal: ref},
		Column: parser.Identifier{Literal: INTERNAL_ID_COLUMN},
	}
	idx, err := view.Header.Contains(fieldRef)
	if err != nil {
		return -1, NewInternalRecordIdNotExistError()
	}
	internalId, ok := view.Records[recordIndex][idx].Primary().(parser.Integer)
	if !ok {
		return -1, NewInternalRecordIdEmptyError()
	}
	return int(internalId.Value()), nil
}

func (view *View) FieldLen() int {
	return view.Header.Len()
}

func (view *View) RecordLen() int {
	return view.Len()
}

func (view *View) Len() int {
	return len(view.Records)
}

func (view *View) Swap(i, j int) {
	view.Records[i], view.Records[j] = view.Records[j], view.Records[i]
}

func (view *View) Less(i, j int) bool {
	for k, idx := range view.sortIndices {
		pi := view.Records[i][idx].Primary()
		pj := view.Records[j][idx].Primary()

		t := EquivalentTo(pi, pj)
		switch t {
		case ternary.TRUE:
			continue
		case ternary.UNKNOWN:
			if parser.IsNull(pi) {
				if view.sortNullPositions[k] == parser.FIRST {
					return true
				} else {
					return false
				}
			} else if parser.IsNull(pj) {
				if view.sortNullPositions[k] == parser.FIRST {
					return false
				} else {
					return true
				}
			} else {
				continue
			}
		case ternary.FALSE:
			t = LessThan(pi, pj)
		}

		if view.sortDirections[k] == parser.ASC {
			return t == ternary.TRUE
		} else {
			return t == ternary.FALSE
		}
	}
	return false
}

func (view *View) Rollback() {
	view.Records = view.FileInfo.InitialRecords.Copy()
}

func (view *View) Copy() *View {
	header := view.Header.Copy()
	records := view.Records.Copy()

	return &View{
		Header:   header,
		Records:  records,
		FileInfo: view.FileInfo,
	}
}
