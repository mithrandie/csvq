package query

import (
	"errors"
	"fmt"
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
		return errors.New("duplicate alias")
	}
	m.views[view.FileInfo.Path] = view.Copy()
	m.alias[alias] = view.FileInfo.Path
	return nil
}

func (m *ViewMap) SetAlias(alias string, fpath string) error {
	if _, ok := m.alias[alias]; ok {
		return errors.New("duplicate alias")
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
	parentFilter Filter

	filteredIndices []int

	sortIndices       []int
	sortDirections    []int
	sortNullPositions []int

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
		view, err := loadView(v.(parser.Table), parentFilter, view.UseInternalId)
		if err != nil {
			return err
		}
		views[i] = view
	}

	view.Header = views[0].Header
	view.Records = views[0].Records
	view.FileInfo = views[0].FileInfo

	for i := 1; i < len(views); i++ {
		CrossJoin(view, views[i])
	}

	if parentFilter != nil {
		view.parentFilter = parentFilter
	}
	return nil
}

func (view *View) LoadFromIdentifier(table parser.Identifier) error {
	fromClause := parser.FromClause{
		Tables: []parser.Expression{
			parser.Table{Object: table},
		},
	}
	var filter Filter

	return view.Load(fromClause, filter)
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
			Path:      "__stdin",
			Delimiter: delimiter,
		}

		if _, ok := ViewCache.Exists(fileInfo.Path); !ok {
			file := os.Stdin
			defer file.Close()
			err = loadViewFromFile(file, fileInfo, table.Name())
		} else {
			if _, ok := ViewCache.HasAlias(table.Name()); !ok {
				ViewCache.SetAlias(table.Name(), fileInfo.Path)
			}
		}
		if err == nil {
			if useInternalId {
				view, _ = ViewCache.GetWithInternalId(fileInfo.Path)
			} else {
				view, _ = ViewCache.Get(fileInfo.Path)
			}
		}
	case parser.Identifier:
		flags := cmd.GetFlags()
		fileInfo, err := NewFileInfo(table.Object.(parser.Identifier).Literal, flags.Repository, flags.Delimiter)
		if err != nil {
			return nil, err
		}

		if _, ok := ViewCache.Exists(fileInfo.Path); !ok {
			file, err := os.Open(fileInfo.Path)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			err = loadViewFromFile(file, fileInfo, table.Name())
		} else {
			if _, ok := ViewCache.HasAlias(table.Name()); !ok {
				ViewCache.SetAlias(table.Name(), fileInfo.Path)
			}
		}
		if err == nil {
			if useInternalId {
				view, _ = ViewCache.GetWithInternalId(fileInfo.Path)
			} else {
				view, _ = ViewCache.Get(fileInfo.Path)
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
		view, err = Select(subquery.Query, parentFilter)
		if err == nil {
			for i := range view.Header {
				view.Header[i].Reference = table.Name()
			}
		}
	}

	if err != nil {
		return nil, err
	}
	return view, nil
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

func NewViewFromGroupedRecord(fr FilterRecord) *View {
	view := new(View)
	view.Header = fr.View.Header
	record := fr.View.Records[fr.RecordIndex]

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
		var filter Filter = append([]FilterRecord{{View: view, RecordIndex: i}}, view.parentFilter...)
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
		var filter Filter = append([]FilterRecord{{View: view, RecordIndex: i}}, view.parentFilter...)

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

	var evalFields = func(view *View, fields []parser.Expression) ([]Record, error) {
		records := make([]Record, view.RecordLen())
		for i := range view.Records {
			var record Record
			var filter Filter = append([]FilterRecord{{View: view, RecordIndex: i}}, view.parentFilter...)

			for _, f := range fields {
				field := f.(parser.Field)
				primary, err := filter.Evaluate(field.Object)
				if err != nil {
					return nil, err
				}
				if _, ok := field.Object.(parser.FieldReference); !ok {
					record = append(record, NewCell(primary))
				}
			}
			records[i] = record
		}
		return records, nil
	}

	fields := parseAllColumns(view, clause.Fields)
	records, err := evalFields(view, fields)
	if err != nil {
		if _, ok := err.(*NotGroupedError); ok {
			view.group(nil)
			records, err = evalFields(view, fields)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	view.selectFields = make([]int, len(fields))
	for i, f := range fields {
		field := f.(parser.Field)
		if fieldRef, ok := field.Object.(parser.FieldReference); ok {
			idx, err := view.Header.Contains(fieldRef)
			if err != nil {
				return err
			}
			view.selectFields[i] = idx
		} else {
			view.Header, view.selectFields[i] = AddHeaderField(view.Header, field.Name())
		}
	}

	for i := range view.Records {
		if 0 < len(records[i]) {
			view.Records[i] = append(view.Records[i], records[i]...)
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
		switch oi.Item.(type) {
		case parser.FieldReference:
			idx, err := view.FieldIndex(oi.Item.(parser.FieldReference))
			if err != nil {
				return err
			}
			view.sortIndices[i] = idx
		default:
			idx, err := view.Header.ContainsAlias(oi.Item.String())
			if err != nil {
				for j := range view.Records {
					var filter Filter = append([]FilterRecord{{View: view, RecordIndex: j}}, view.parentFilter...)

					primary, err := filter.Evaluate(oi.Item)
					if err != nil {
						return err
					}
					view.Records[j] = append(view.Records[j], NewCell(primary))
				}
				view.Header, idx = AddHeaderField(view.Header, oi.Item.String())
			}
			view.sortIndices[i] = idx
		}

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

func (view *View) Offset(clause parser.OffsetClause) {
	if int64(len(view.Records)) <= clause.Number {
		view.Records = Records{}
	} else {
		view.Records = view.Records[clause.Number:]
		records := make(Records, len(view.Records))
		copy(records, view.Records)
		view.Records = records
	}
}

func (view *View) Limit(clause parser.LimitClause) {
	if clause.Number < int64(len(view.Records)) {
		view.Records = view.Records[:clause.Number]
		records := make(Records, len(view.Records))
		copy(records, view.Records)
		view.Records = records
	}
}

func (view *View) InsertValues(fields []parser.Expression, list []parser.Expression) error {
	var filter Filter
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

func (view *View) InsertFromQuery(fields []parser.Expression, query parser.SelectQuery) error {
	insertView, err := Select(query, nil)
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
	view.parentFilter = Filter(nil)
	view.sortIndices = []int(nil)
	view.sortDirections = []int(nil)
	view.sortNullPositions = []int(nil)
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
