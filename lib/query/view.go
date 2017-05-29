package query

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/csv"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type FileInfo struct {
	Path      string
	Delimiter rune
}

func NewFileInfo(filename string, repository string, delimiter rune) (*FileInfo, error) {
	filepath := filename
	if !path.IsAbs(filepath) {
		filepath = path.Join(repository, filepath)
	}

	var info os.FileInfo
	var err error

	if info, err = os.Stat(filepath); err != nil {
		if info, err = os.Stat(filepath + cmd.CSV_EXT); err == nil {
			filepath = filepath + cmd.CSV_EXT
		} else if info, err = os.Stat(filepath + cmd.TSV_EXT); err == nil {
			filepath = filepath + cmd.TSV_EXT
		} else {
			return nil, errors.New(fmt.Sprintf("file %s does not exist", filename))
		}
	}

	if info.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is a directory", filename))
	}

	if delimiter == cmd.UNDEF {
		if strings.ToUpper(path.Ext(filepath)) == strings.ToUpper(cmd.TSV_EXT) {
			delimiter = '\t'
		} else {
			delimiter = ','
		}
	}

	return &FileInfo{
		Path:      filepath,
		Delimiter: delimiter,
	}, nil
}

type View struct {
	Header  Header
	Records []Record

	FileInfo *FileInfo

	selectFields []int
	isGrouped    bool
	parentFilter Filter

	filteredIndices []int

	sortIndices    []int
	sortDirections []int
}

func NewView(clause parser.FromClause, parentFilter Filter) (*View, error) {
	if len(clause.Tables) < 2 {
		if _, ok := clause.Tables[0].(parser.Dual); ok {
			return NewDualView(), nil
		}
	}

	views := make([]*View, len(clause.Tables))
	for i, v := range clause.Tables {
		view, err := loadView(v.(parser.Table), parentFilter)
		if err != nil {
			return nil, err
		}
		views[i] = view
	}

	joinedView := views[0]

	for i := 1; i < len(views); i++ {
		joinedView = CrossJoin(joinedView, views[i])
	}

	if parentFilter != nil {
		joinedView.parentFilter = parentFilter
	}
	return joinedView, nil
}

func loadView(table parser.Table, parentFilter Filter) (*View, error) {
	var view *View
	var err error

	switch table.Object.(type) {
	case parser.Identifier:
		file := table.Object.(parser.Identifier)
		view, err = loadViewFromFile(file.Literal, table.Name())
	case parser.Join:
		join := table.Object.(parser.Join)
		view1, err := loadView(join.Table, parentFilter)
		if err != nil {
			return nil, err
		}
		view2, err := loadView(join.JoinTable, parentFilter)
		if err != nil {
			return nil, err
		}

		condition := ParseJoinCondition(join, view1, view2)

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
			view = CrossJoin(view1, view2)
		case parser.INNER:
			view, err = InnerJoin(view1, view2, condition, parentFilter)
		case parser.OUTER:
			view, err = OuterJoin(view1, view2, condition, join.Direction.Token, parentFilter)
		}
	case parser.Subquery:
		subquery := table.Object.(parser.Subquery)
		view, err = executeSelect(subquery.Query, parentFilter)
	}

	if err != nil {
		return nil, err
	}
	return view, nil
}

func loadViewFromFile(filename string, reference string) (*View, error) {
	flags := cmd.GetFlags()

	fileInfo, err := NewFileInfo(filename, flags.Repository, flags.Delimiter)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fileInfo.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := cmd.GetReader(f, flags.Encoding)

	view := new(View)

	reader := csv.NewReader(r)
	reader.Delimiter = fileInfo.Delimiter
	reader.WithoutNull = flags.WithoutNull

	var header []string
	if !flags.NoHeader {
		header, err = reader.ReadHeader()
		if err != nil {
			return nil, err
		}
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	view.Records = make([]Record, len(records))
	for i, v := range records {
		view.Records[i] = NewRecord(v)
	}

	if header == nil {
		header = make([]string, reader.FieldsPerRecords)
		for i := 0; i < reader.FieldsPerRecords; i++ {
			header[i] = "c" + strconv.Itoa(i+1)
		}
	}
	view.Header = NewHeader(reference, header)

	view.FileInfo = fileInfo

	return view, nil
}

func NewDualView() *View {
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
		if ident, ok := item.(parser.Identifier); ok {
			ref, field, _ := ident.FieldRef()
			idx, _ := view.Header.Contains(ref, field)
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
				if _, ok := field.Object.(parser.Identifier); !ok {
					record = append(record, NewCell(primary))
				}
			}
			records[i] = record
		}
		return records, nil
	}

	var isDuplicate = func(records []Record, record Record) bool {
		for _, w := range records {
			if reflect.DeepEqual(record, w) {
				return true
			}
		}
		return false
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
	addIndex := view.FieldLen()
	for i, f := range fields {
		field := f.(parser.Field)
		if ident, ok := field.Object.(parser.Identifier); ok {
			ref, field, _ := ident.FieldRef()
			idx, _ := view.Header.Contains(ref, field)
			view.selectFields[i] = idx
		} else {
			view.Header = AddHeaderField(view.Header, field.Name())
			view.selectFields[i] = addIndex
			addIndex++
		}
	}

	for i := range view.Records {
		if 0 < len(records[i]) {
			view.Records[i] = append(view.Records[i], records[i]...)
		}
	}

	if clause.IsDistinct() {
		hfields := NewEmptyHeader(len(view.selectFields))
		distinguished := []Record{}

		for _, v := range view.Records {
			record := make(Record, len(view.selectFields))
			for i, idx := range view.selectFields {
				record[i] = v[idx]
			}

			if !isDuplicate(distinguished, record) {
				distinguished = append(distinguished, record)
			}
		}

		for i, idx := range view.selectFields {
			hfields[i] = view.Header[idx]
			view.selectFields[i] = i
		}

		view.Header = hfields
		view.Records = distinguished
	}

	return nil
}

func (view *View) OrderBy(clause parser.OrderByClause) error {
	view.sortIndices = []int{}

	for _, v := range clause.Items {
		oi := v.(parser.OrderItem)
		switch oi.Item.(type) {
		case parser.Identifier:
			item := oi.Item.(parser.Identifier)
			ref, column, err := item.FieldRef()
			if err != nil {
				return err
			}
			idx, err := view.Header.Contains(ref, column)
			if err != nil {
				return err
			}
			view.sortIndices = append(view.sortIndices, idx)
		default:
			idx, err := view.Header.Contains("", oi.String())
			if err != nil {
				for i := range view.Records {
					var filter Filter = append([]FilterRecord{{View: view, RecordIndex: i}}, view.parentFilter...)

					primary, err := filter.Evaluate(oi.Item)
					if err != nil {
						return err
					}
					view.Records[i] = append(view.Records[i], NewCell(primary))
				}
				view.Header = AddHeaderField(view.Header, oi.String())
				idx = view.FieldLen() - 1
			}
			view.sortIndices = append(view.sortIndices, idx)
		}
		view.sortDirections = append(view.sortDirections, oi.Direction.Token)
	}

	direction := parser.ASC
	for i := len(view.sortDirections) - 1; i >= 0; i-- {
		if view.sortDirections[i] == parser.ASC || view.sortDirections[i] == parser.DESC {
			direction = view.sortDirections[i]
		} else {
			view.sortDirections[i] = direction
		}
	}

	sort.Sort(view)
	return nil
}

func (view *View) Limit(clause parser.LimitClause) {
	if clause.Number < int64(len(view.Records)) {
		view.Records = view.Records[:clause.Number]
	}
}

func (view *View) Fix() {
	hfields := NewEmptyHeader(len(view.selectFields))
	records := make([]Record, view.RecordLen())

	for i, v := range view.Records {
		record := make(Record, len(view.selectFields))
		for j, idx := range view.selectFields {
			record[j] = v[idx]
		}

		records[i] = record
	}

	for i, idx := range view.selectFields {
		hfields[i] = view.Header[idx]
	}

	view.Header = hfields
	view.Records = records
	view.selectFields = []int(nil)
	view.isGrouped = false
	view.parentFilter = Filter(nil)
	view.sortIndices = []int(nil)
	view.sortDirections = []int(nil)
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
				t = ternary.TRUE
			} else if parser.IsNull(pj) {
				t = ternary.FALSE
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
