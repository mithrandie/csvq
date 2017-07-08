package query

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

const STDIN_VIRTUAL_FILE_PATH = ";;__STDIN__;;"

type ViewMap struct {
	views map[string]*View
	alias map[string]string
}

func NewViewMap() *ViewMap {
	return &ViewMap{
		views: make(map[string]*View),
		alias: make(map[string]string),
	}
}

func (m *ViewMap) Exists(fpath string) (string, bool) {
	if _, ok := m.views[fpath]; ok {
		return fpath, true
	}
	if substance, ok := m.alias[fpath]; ok {
		if _, ok := m.views[substance]; ok {
			return substance, true
		}
	}
	return "", false
}

func (m *ViewMap) HasAlias(alias string) (string, bool) {
	if fpath, ok := m.alias[alias]; ok {
		return fpath, true
	}
	return "", false
}

func (m *ViewMap) Get(fpath string) (*View, error) {
	if pt, ok := m.Exists(fpath); ok {
		return m.views[pt].Copy(), nil
	}
	return nil, errors.New(fmt.Sprintf("file %s is not loaded", fpath))
}

func (m *ViewMap) GetWithInternalId(fpath string) (*View, error) {
	if pt, ok := m.Exists(fpath); ok {
		ret := m.views[pt].Copy()

		if 0 < ret.FieldLen() {
			ret.Header = MergeHeader(NewHeader(ret.Header[0].Reference, []string{}), ret.Header)

			for i, v := range ret.Records {
				ret.Records[i] = append(Record{NewCell(parser.NewInteger(int64(i)))}, v...)
			}
		}

		return ret, nil

	}
	return nil, errors.New(fmt.Sprintf("file %s is not loaded", fpath))
}

func (m *ViewMap) Set(view *View, alias string) error {
	if view.FileInfo == nil || len(view.FileInfo.Path) < 1 {
		return errors.New("view cache failed")
	}
	if _, ok := m.alias[alias]; ok {
		return errors.New(fmt.Sprintf("table name %s is duplicated", alias))
	}
	m.views[view.FileInfo.Path] = view.Copy()
	m.alias[alias] = view.FileInfo.Path
	return nil
}

func (m *ViewMap) SetAlias(alias string, fpath string) error {
	if _, ok := m.alias[alias]; ok {
		return errors.New(fmt.Sprintf("table name %s is duplicated", alias))
	}
	m.alias[alias] = fpath
	return nil
}

func (m *ViewMap) Replace(view *View) error {
	if pt, ok := m.Exists(view.FileInfo.Path); ok {
		m.views[pt] = view
	}
	return errors.New(fmt.Sprintf("file %s is not loaded", view.FileInfo.Path))
}

func (m *ViewMap) Clear() {
	for k := range m.views {
		delete(m.views, k)
	}
	for k := range m.alias {
		delete(m.alias, k)
	}
}

func (m *ViewMap) ClearAliases() {
	for k := range m.alias {
		delete(m.alias, k)
	}
}

type FileInfo struct {
	Path      string
	Delimiter rune
	NoHeader  bool
	Encoding  cmd.Encoding
	LineBreak cmd.LineBreak
}

func NewFileInfo(filename string, repository string, delimiter rune) (*FileInfo, error) {
	fpath := filename
	if !path.IsAbs(fpath) {
		fpath = path.Join(repository, fpath)
	}

	var info os.FileInfo
	var err error

	if info, err = os.Stat(fpath); err != nil {
		if info, err = os.Stat(fpath + cmd.CSV_EXT); err == nil {
			fpath = fpath + cmd.CSV_EXT
		} else if info, err = os.Stat(fpath + cmd.TSV_EXT); err == nil {
			fpath = fpath + cmd.TSV_EXT
		} else {
			return nil, errors.New(fmt.Sprintf("file %s does not exist", filename))
		}
	}

	fpath, err = filepath.Abs(fpath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("file %s does not exist", filename))
	}

	if info.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is a directory", filename))
	}

	if delimiter == cmd.UNDEF {
		if strings.EqualFold(path.Ext(fpath), cmd.TSV_EXT) {
			delimiter = '\t'
		} else {
			delimiter = ','
		}
	}

	return &FileInfo{
		Path:      fpath,
		Delimiter: delimiter,
	}, nil
}

func isReadableFromStdin() bool {
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeNamedPipe != 0 || 0 < fi.Size()) {
		return true
	}
	return false
}

type View struct {
	Header  Header
	Records Records

	FileInfo *FileInfo

	selectFields []int
	isGrouped    bool
	ParentFilter Filter

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

func (view *View) Load(clause parser.FromClause, parentFilter Filter) error {
	if clause.Tables == nil {
		var obj parser.Expression
		if isReadableFromStdin() {
			obj = parser.Stdin{Stdin: "stdin"}
		} else {
			obj = parser.Dual{}
		}
		clause.Tables = []parser.Expression{parser.Table{Object: obj}}
	}

	views := make([]*View, len(clause.Tables))
	for i, v := range clause.Tables {
		loaded, err := loadView(v.(parser.Table), parentFilter, view.UseInternalId)
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

	view.ParentFilter = parentFilter
	return nil
}

func (view *View) LoadFromIdentifier(table parser.Identifier) error {
	return view.LoadFromIdentifierWithCommonTables(table, Filter{})
}

func (view *View) LoadFromIdentifierWithCommonTables(table parser.Identifier, parentFilter Filter) error {
	fromClause := parser.FromClause{
		Tables: []parser.Expression{
			parser.Table{Object: table},
		},
	}

	return view.Load(fromClause, parentFilter)
}

func loadView(table parser.Table, parentFilter Filter, useInternalId bool) (*View, error) {
	var view *View
	var err error

	switch table.Object.(type) {
	case parser.Dual:
		view = loadDualView()
	case parser.Stdin:
		if !isReadableFromStdin() {
			return nil, errors.New("stdin is empty")
		}

		delimiter := cmd.GetFlags().Delimiter
		if delimiter == cmd.UNDEF {
			delimiter = ','
		}
		fileInfo := &FileInfo{
			Path:      STDIN_VIRTUAL_FILE_PATH,
			Delimiter: delimiter,
		}

		if _, ok := ViewCache.Exists(fileInfo.Path); !ok {
			file := os.Stdin
			defer file.Close()
			if err := loadViewFromFile(file, fileInfo, table.Name()); err != nil {
				return nil, err
			}
		} else {
			if _, ok := ViewCache.HasAlias(table.Name()); !ok {
				ViewCache.SetAlias(table.Name(), fileInfo.Path)
			}
		}
		if useInternalId {
			view, _ = ViewCache.GetWithInternalId(fileInfo.Path)
		} else {
			view, _ = ViewCache.Get(fileInfo.Path)
		}
	case parser.Identifier:
		tableIdentifier := table.Object.(parser.Identifier).Literal
		if strings.EqualFold(tableIdentifier, parentFilter.RecursiveTable.Name.Literal) && parentFilter.RecursiveTmpView != nil {
			view = parentFilter.RecursiveTmpView
			if parentFilter.RecursiveTable.Name.Literal != table.Name() {
				view.UpdateHeader(table.Name(), nil)
			}
		} else if ct, err := parentFilter.CommonTables.Get(tableIdentifier); err == nil {
			view = ct
			if tableIdentifier != table.Name() {
				view.UpdateHeader(table.Name(), nil)
			}
		} else if _, err := parentFilter.CommonTables.Get(table.Name()); err == nil {
			return nil, errors.New(fmt.Sprintf("table name %s is duplicated", table.Name()))
		} else {

			flags := cmd.GetFlags()

			fileInfo, err := NewFileInfo(tableIdentifier, flags.Repository, flags.Delimiter)
			if err != nil {
				return nil, err
			}

			commonTableName := parser.FormatTableName(fileInfo.Path)

			if _, ok := ViewCache.Exists(fileInfo.Path); !ok {
				file, err := os.Open(fileInfo.Path)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				if err := loadViewFromFile(file, fileInfo, commonTableName); err != nil {
					return nil, err
				}
			}
			if _, ok := ViewCache.HasAlias(table.Name()); !ok {
				ViewCache.SetAlias(table.Name(), fileInfo.Path)
			}

			if useInternalId {
				view, _ = ViewCache.GetWithInternalId(fileInfo.Path)
			} else {
				view, _ = ViewCache.Get(fileInfo.Path)
			}
			if commonTableName != table.Name() {
				view.UpdateHeader(table.Name(), nil)
			}
		}
	case parser.Join:
		join := table.Object.(parser.Join)
		view, err = loadView(join.Table, parentFilter, useInternalId)
		if err != nil {
			return nil, err
		}
		view2, err := loadView(join.JoinTable, parentFilter, useInternalId)
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
			err = InnerJoin(view, view2, condition, parentFilter)
		case parser.OUTER:
			err = OuterJoin(view, view2, condition, join.Direction.Token, parentFilter)
		}
	case parser.Subquery:
		subquery := table.Object.(parser.Subquery)
		view, err = SelectAsSubquery(subquery.Query, parentFilter)
		if err == nil {
			view.UpdateHeader(table.Name(), nil)
		}
	}

	return view, err
}

func loadViewFromFile(file *os.File, fileInfo *FileInfo, reference string) error {
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
			return err
		}
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	view.Records = make([]Record, len(records))
	for i, v := range records {
		view.Records[i] = NewRecordWithoutId(v)
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

	view.Header = NewHeaderWithoutId(reference, header)
	view.FileInfo = fileInfo
	ViewCache.Set(view, reference)
	return nil
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
			view.Records[i][j] = NewCell(cell.GroupedPrimary(i))
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
		filter := NewFilterForRecord(view, i, view.ParentFilter)
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
		filter := NewFilterForRecord(view, i, view.ParentFilter)

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
		if _, ok := err.(*NotGroupedError); ok {
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
		insert := make([]parser.Expression, len(columns))
		for i, c := range columns {
			insert[i] = parser.Field{
				Object: c,
			}
		}

		return append(append(fields[:insertIdx], insert...), fields[insertIdx+1:]...)
	}

	var evalFields = func(view *View, fields []parser.Expression) error {
		view.selectFields = make([]int, len(fields))
		for i, f := range fields {
			field := f.(parser.Field)
			idx, err := view.evalColumn(field.Object, field.Object.String(), field.Name())
			if err != nil {
				return err
			}
			view.selectFields[i] = idx
		}
		return nil
	}

	fields := parseAllColumns(view, clause.Fields)
	err := evalFields(view, fields)
	if err != nil {
		if _, ok := err.(*NotGroupedError); ok {
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
	switch obj.(type) {
	case parser.FieldReference:
		idx, err = view.FieldIndex(obj.(parser.FieldReference))
	default:
		idx, err = view.Header.ContainsObject(obj)
		if err != nil {
			err = nil

			if analyticFunction, ok := obj.(parser.AnalyticFunction); ok {
				err = view.evalAnalyticFunction(analyticFunction)
				if err != nil {
					return
				}
			} else {
				filter := NewFilterForRecord(view, 0, view.ParentFilter)
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
	return
}

func (view *View) evalAnalyticFunction(expr parser.AnalyticFunction) error {
	DefineAnalyticFunctions()

	name := strings.ToUpper(expr.Name)
	fn, ok := AnalyticFunctions[name]
	if !ok {
		return errors.New(fmt.Sprintf("function %s does not exist", expr.Name))
	}

	if expr.Option.IsDistinct() {
		return errors.New(fmt.Sprintf("syntax error: unexpected %s", expr.Option.Distinct.Literal))
	}

	if expr.AnalyticClause.OrderByClause != nil {
		err := view.OrderBy(expr.AnalyticClause.OrderByClause.(parser.OrderByClause))
		if err != nil {
			return err
		}
	}

	return fn(view, expr.Option.Args, expr.AnalyticClause)
}

func (view *View) Offset(clause parser.OffsetClause) error {
	var filter Filter
	value, err := filter.Evaluate(clause.Value)
	if err != nil {
		return err
	}
	number := parser.PrimaryToInteger(value)
	if parser.IsNull(number) {
		return errors.New("offset number is not an integer")
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
	var filter Filter
	value, err := filter.Evaluate(clause.Value)
	if err != nil {
		return err
	}

	var limit int
	if clause.IsPercentage() {
		number := parser.PrimaryToFloat(value)
		if parser.IsNull(number) {
			return errors.New("limit percentage is not a float value")
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
			return errors.New("limit number of records is not an integer value")
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
		values, err := filter.evalRowValue(item.(parser.RowValue))
		if err != nil {
			return err
		}
		if len(fields) != len(values) {
			return errors.New("field length does not match value length")
		}

		valuesList[i] = values
	}

	return view.insert(fields, valuesList)
}

func (view *View) InsertFromQuery(fields []parser.Expression, query parser.SelectQuery, filter Filter) error {
	insertView, err := SelectAsSubquery(query, filter)
	if err != nil {
		return err
	}
	if len(fields) != insertView.FieldLen() {
		return errors.New("field length does not match value length")
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

	for i, idx := range view.selectFields {
		hfields[i] = view.Header[idx]
		hfields[i].FromTable = true
		hfields[i].IsGroupKey = false
	}

	view.Header = hfields
	view.Records = records
	view.selectFields = []int(nil)
	view.isGrouped = false
	view.ParentFilter = Filter{}
	view.sortIndices = []int(nil)
	view.sortDirections = []int(nil)
	view.sortNullPositions = []int(nil)
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

func (view *View) FieldIndex(fieldRef parser.FieldReference) (int, error) {
	return view.Header.Contains(fieldRef)
}

func (view *View) FieldIndices(fields []parser.Expression) ([]int, error) {
	indices := make([]int, len(fields))
	for i, v := range fields {
		idx, err := view.FieldIndex(v.(parser.FieldReference))
		if err != nil {
			return nil, err
		}
		indices[i] = idx
	}
	return indices, nil
}

func (view *View) FieldViewName(fieldRef parser.FieldReference) (string, error) {
	idx, err := view.FieldIndex(fieldRef)
	if err != nil {
		return "", err
	}
	return view.Header[idx].Reference, nil
}

func (view *View) InternalRecordId(ref string, recordIndex int) (int, error) {
	fieldRef := parser.FieldReference{
		View:   parser.Identifier{Literal: ref},
		Column: parser.Identifier{Literal: INTERNAL_ID_COLUMN},
	}
	idx, err := view.Header.Contains(fieldRef)
	if err != nil {
		return -1, errors.New("internal record id does not exist")
	}
	internalId, ok := view.Records[recordIndex][idx].Primary().(parser.Integer)
	if !ok {
		return -1, errors.New("internal record id is empty")
	}
	return int(internalId.Value()), nil
}

func (view *View) UpdateHeader(reference string, fields []parser.Expression) error {
	if fields != nil && len(fields) != view.FieldLen() {
		return errors.New(fmt.Sprintf("common table %s: field length does not match", reference))
	}

	for i := range view.Header {
		view.Header[i].Reference = reference
		if fields != nil {
			view.Header[i].Column = fields[i].(parser.Identifier).Literal
		} else if 0 < len(view.Header[i].Alias) {
			view.Header[i].Column = view.Header[i].Alias
		}
		view.Header[i].Alias = ""
	}
	return nil
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

func (view *View) Copy() *View {
	header := view.Header.Copy()

	records := make([]Record, view.RecordLen())
	for i, v := range view.Records {
		records[i] = v.Copy()
	}

	return &View{
		Header:   header,
		Records:  records,
		FileInfo: view.FileInfo,
	}
}
