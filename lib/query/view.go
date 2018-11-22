package query

import (
	"bytes"
	gojson "encoding/json"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/json"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text"
	"github.com/mithrandie/go-text/csv"
	"github.com/mithrandie/go-text/fixedlen"
	txjson "github.com/mithrandie/go-text/json"
	"github.com/mithrandie/ternary"
)

type RecordReader interface {
	Read() ([]text.RawText, error)
}

type View struct {
	Header    Header
	RecordSet RecordSet
	FileInfo  *FileInfo

	Filter *Filter

	selectFields []int
	selectLabels []string
	isGrouped    bool

	comparisonKeysInEachRecord []string
	sortValuesInEachCell       [][]*SortValue
	sortValuesInEachRecord     []SortValues
	sortDirections             []int
	sortNullPositions          []int

	offset int

	UseInternalId bool
	ForUpdate     bool
}

func NewView() *View {
	return &View{
		UseInternalId: false,
	}
}

func (view *View) Load(clause parser.FromClause, filter *Filter) error {
	if clause.Tables == nil {
		var obj parser.QueryExpression
		if cmd.IsReadableFromPipeOrRedirection() {
			obj = parser.Stdin{Stdin: "stdin"}
		} else {
			obj = parser.Dual{}
		}
		clause.Tables = []parser.QueryExpression{parser.Table{Object: obj}}
	}

	views := make([]*View, len(clause.Tables))
	for i, v := range clause.Tables {
		loaded, err := loadView(v, filter, view.UseInternalId, view.ForUpdate)
		if err != nil {
			return err
		}
		views[i] = loaded
	}

	view.Header = views[0].Header
	view.RecordSet = views[0].RecordSet
	view.FileInfo = views[0].FileInfo

	for i := 1; i < len(views); i++ {
		CrossJoin(view, views[i])
	}

	view.Filter = filter
	return nil
}

func (view *View) LoadFromTableIdentifier(table parser.QueryExpression, filter *Filter) error {
	fromClause := parser.FromClause{
		Tables: []parser.QueryExpression{
			parser.Table{Object: table},
		},
	}

	return view.Load(fromClause, filter)
}

func loadView(tableExpr parser.QueryExpression, filter *Filter, useInternalId bool, forUpdate bool) (*View, error) {
	if parentheses, ok := tableExpr.(parser.Parentheses); ok {
		return loadView(parentheses.Expr, filter, useInternalId, forUpdate)
	}

	table := tableExpr.(parser.Table)

	var view *View
	var err error

	switch table.Object.(type) {
	case parser.Dual:
		view = loadDualView()
	case parser.Stdin:
		flags := cmd.GetFlags()
		fileInfo := &FileInfo{
			Path:               table.Object.String(),
			Format:             flags.SelectImportFormat(),
			Delimiter:          flags.Delimiter,
			DelimiterPositions: flags.DelimiterPositions,
			JsonQuery:          flags.JsonQuery,
			Encoding:           flags.Encoding,
			LineBreak:          flags.LineBreak,
			NoHeader:           flags.NoHeader,
			EncloseAll:         flags.EncloseAll,
			JsonEscape:         flags.JsonEscape,
			IsTemporary:        true,
		}

		if !filter.TempViews[len(filter.TempViews)-1].Exists(fileInfo.Path) {
			if !cmd.IsReadableFromPipeOrRedirection() {
				return nil, NewStdinEmptyError(table.Object.(parser.Stdin))
			}

			var loadView *View

			if fileInfo.Format != cmd.JSON {
				fp := os.Stdin
				defer fp.Close()

				loadView, err = loadViewFromFile(fp, fileInfo, flags.WithoutNull)
				if err != nil {
					return nil, NewDataParsingError(table.Object, fileInfo.Path, err.Error())
				}
			} else {
				fileInfo.Encoding = text.UTF8

				buf, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return nil, NewReadFileError(table.Object.(parser.Stdin), err.Error())
				}

				headerLabels, rows, escapeType, err := json.LoadTable(fileInfo.JsonQuery, string(buf))
				if err != nil {
					return nil, NewJsonQueryError(parser.JsonQuery{BaseExpr: table.Object.GetBaseExpr()}, err.Error())
				}

				records := make([]Record, 0, len(rows))
				for _, row := range rows {
					records = append(records, NewRecord(row))
				}

				fileInfo.JsonEscape = escapeType

				loadView = NewView()
				loadView.Header = NewHeader(parser.FormatTableName(fileInfo.Path), headerLabels)
				loadView.RecordSet = records
				loadView.FileInfo = fileInfo
			}

			loadView.FileInfo.InitialHeader = loadView.Header.Copy()
			loadView.FileInfo.InitialRecordSet = loadView.RecordSet.Copy()
			filter.TempViews[len(filter.TempViews)-1].Set(loadView)
		}
		if err = filter.Aliases.Add(table.Name(), fileInfo.Path); err != nil {
			return nil, err
		}

		pathIdent := parser.Identifier{Literal: table.Object.String()}
		if useInternalId {
			view, _ = filter.TempViews[len(filter.TempViews)-1].GetWithInternalId(pathIdent)
		} else {
			view, _ = filter.TempViews[len(filter.TempViews)-1].Get(pathIdent)
		}
		if !strings.EqualFold(table.Object.String(), table.Name().Literal) {
			view.Header.Update(table.Name().Literal, nil)
		}
	case parser.TableObject:
		tableObject := table.Object.(parser.TableObject)

		flags := cmd.GetFlags()
		importFormat := flags.SelectImportFormat()
		delimiter := flags.Delimiter
		delimiterPositions := flags.DelimiterPositions
		jsonQuery := flags.JsonQuery
		encoding := flags.Encoding
		noHeader := flags.NoHeader
		withoutNull := flags.WithoutNull

		felem, err := filter.Evaluate(tableObject.FormatElement)
		if err != nil {
			return nil, err
		}
		felemStr := value.ToString(felem)

		switch strings.ToUpper(tableObject.Type.Literal) {
		case cmd.CSV.String():
			if value.IsNull(felemStr) {
				return nil, NewTableObjectInvalidDelimiterError(tableObject, tableObject.FormatElement.String())
			}
			s := cmd.UnescapeString(felemStr.(value.String).Raw())
			d := []rune(s)
			if 1 != len(d) {
				return nil, NewTableObjectInvalidDelimiterError(tableObject, tableObject.FormatElement.String())
			}
			if 3 < len(tableObject.Args) {
				return nil, NewTableObjectArgumentsLengthError(tableObject, 5)
			}
			delimiter = d[0]
			if delimiter == '\t' {
				importFormat = cmd.TSV
			} else {
				importFormat = cmd.CSV
			}
		case cmd.FIXED.String():
			if value.IsNull(felemStr) {
				return nil, NewTableObjectInvalidDelimiterPositionsError(tableObject, tableObject.FormatElement.String())
			}
			s := felemStr.(value.String).Raw()

			var positions []int
			if !strings.EqualFold("SPACES", s) {
				err = gojson.Unmarshal([]byte(s), &positions)
				if err != nil {
					return nil, NewTableObjectInvalidDelimiterPositionsError(tableObject, tableObject.FormatElement.String())
				}
			}
			if 3 < len(tableObject.Args) {
				return nil, NewTableObjectArgumentsLengthError(tableObject, 5)
			}
			delimiterPositions = positions
			importFormat = cmd.FIXED
		case cmd.JSON.String():
			if value.IsNull(felemStr) {
				return nil, NewTableObjectInvalidJsonQueryError(tableObject, tableObject.FormatElement.String())
			}
			if 0 < len(tableObject.Args) {
				return nil, NewTableObjectJsonArgumentsLengthError(tableObject, 2)
			}
			jsonQuery = felemStr.(value.String).Raw()
			importFormat = cmd.JSON
			encoding = text.UTF8
		default:
			return nil, NewTableObjectInvalidObjectError(tableObject, tableObject.Type.Literal)
		}

		args := make([]value.Primary, 3)
		for i, a := range tableObject.Args {
			p, err := filter.Evaluate(a)
			if err != nil {
				if appErr, ok := err.(AppError); ok {
					err = NewTableObjectInvalidArgumentError(tableObject, appErr.ErrorMessage())
				}
				return nil, err
			}
			switch i {
			case 0:
				v := value.ToString(p)
				if !value.IsNull(v) {
					args[i] = v
				}
			default:
				v := value.ToBoolean(p)
				if !value.IsNull(v) {
					args[i] = v
				}
			}
		}

		if args[0] != nil {
			if encoding, err = cmd.ParseEncoding(args[0].(value.String).Raw()); err != nil {
				return nil, NewTableObjectInvalidArgumentError(tableObject, err.Error())
			}
		}
		if args[1] != nil {
			noHeader = args[1].(value.Boolean).Raw()
		}
		if args[2] != nil {
			withoutNull = args[2].(value.Boolean).Raw()
		}

		view, err = loadObject(
			table.Object.(parser.TableObject).Path,
			table.Name(),
			filter,
			useInternalId,
			forUpdate,
			importFormat,
			delimiter,
			delimiterPositions,
			jsonQuery,
			encoding,
			flags.LineBreak,
			noHeader,
			flags.EncloseAll,
			flags.JsonEscape,
			withoutNull,
		)
		if err != nil {
			return nil, err
		}

	case parser.Identifier:
		flags := cmd.GetFlags()

		view, err = loadObject(
			table.Object.(parser.Identifier),
			table.Name(),
			filter,
			useInternalId,
			forUpdate,
			cmd.AutoSelect,
			flags.Delimiter,
			flags.DelimiterPositions,
			flags.JsonQuery,
			flags.Encoding,
			flags.LineBreak,
			flags.NoHeader,
			flags.EncloseAll,
			flags.JsonEscape,
			flags.WithoutNull,
		)
		if err != nil {
			return nil, err
		}
	case parser.Join:
		join := table.Object.(parser.Join)
		view, err = loadView(join.Table, filter, useInternalId, forUpdate)
		if err != nil {
			return nil, err
		}
		view2, err := loadView(join.JoinTable, filter, useInternalId, forUpdate)
		if err != nil {
			return nil, err
		}

		condition, includeFields, excludeFields, err := ParseJoinCondition(join, view, view2)
		if err != nil {
			return nil, err
		}

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

		includeIndices := make([]int, 0, len(includeFields))
		excludeIndices := make([]int, 0, len(includeFields))
		if includeFields != nil {
			for i := range includeFields {
				idx, _ := view.Header.Contains(includeFields[i])
				includeIndices = append(includeIndices, idx)

				idx, _ = view.Header.Contains(excludeFields[i])
				excludeIndices = append(excludeIndices, idx)
			}

			fieldIndices := make([]int, 0, view.FieldLen())
			header := make(Header, 0, view.FieldLen()-len(excludeIndices))
			for _, idx := range includeIndices {
				view.Header[idx].View = ""
				view.Header[idx].Number = 0
				view.Header[idx].IsJoinColumn = true
				header = append(header, view.Header[idx])
				fieldIndices = append(fieldIndices, idx)
			}
			for i := range view.Header {
				if InIntSlice(i, excludeIndices) || InIntSlice(i, includeIndices) {
					continue
				}
				header = append(header, view.Header[i])
				fieldIndices = append(fieldIndices, i)
			}
			view.Header = header

			gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
			for i := 0; i < gm.CPU; i++ {
				gm.Add()
				go func(thIdx int) {
					start, end := gm.RecordRange(thIdx)

					for i := start; i < end; i++ {
						record := make(Record, len(fieldIndices))
						for j, idx := range fieldIndices {
							record[j] = view.RecordSet[i][idx]
						}
						view.RecordSet[i] = record
					}

					gm.Done()
				}(i)
			}
			gm.Wait()
		}

	case parser.JsonQuery:
		jsonQuery := table.Object.(parser.JsonQuery)
		alias := table.Name().Literal

		queryValue, err := filter.Evaluate(jsonQuery.Query)
		if err != nil {
			return nil, err
		}
		queryValue = value.ToString(queryValue)

		if value.IsNull(queryValue) {
			return nil, NewJsonQueryEmptyError(jsonQuery)
		}

		var reader io.Reader

		if jsonPath, ok := jsonQuery.JsonText.(parser.Identifier); ok {
			fpath, err := SearchJsonFilePath(jsonPath, cmd.GetFlags().Repository)
			if err != nil {
				return nil, err
			}

			h, err := file.NewHandlerForRead(fpath)
			if err != nil {
				if _, ok := err.(*file.TimeoutError); ok {
					return nil, NewFileLockTimeoutError(jsonPath, fpath)
				}
				return nil, NewReadFileError(jsonPath, err.Error())
			}
			defer h.Close()
			reader = h.FileForRead()
		} else {
			jsonTextValue, err := filter.Evaluate(jsonQuery.JsonText)
			if err != nil {
				return nil, err
			}
			jsonTextValue = value.ToString(jsonTextValue)

			if value.IsNull(jsonTextValue) {
				return nil, NewJsonTableEmptyError(jsonQuery)
			}

			reader = strings.NewReader(jsonTextValue.(value.String).Raw())
		}

		fileInfo := &FileInfo{
			Path:        alias,
			Format:      cmd.JSON,
			JsonQuery:   queryValue.(value.String).Raw(),
			Encoding:    text.UTF8,
			LineBreak:   cmd.GetFlags().LineBreak,
			IsTemporary: true,
		}

		view, err = loadViewFromJsonFile(reader, fileInfo)
		if err != nil {
			return nil, NewJsonQueryError(jsonQuery, err.Error())
		}

		if err = filter.Aliases.Add(table.Name(), ""); err != nil {
			return nil, err
		}

	case parser.Subquery:
		subquery := table.Object.(parser.Subquery)
		view, err = Select(subquery.Query, filter)
		if err != nil {
			return nil, err
		}

		view.Header.Update(table.Name().Literal, nil)

		if err = filter.Aliases.Add(table.Name(), ""); err != nil {
			return nil, err
		}
	}

	return view, err
}

func loadObject(
	tableIdentifier parser.Identifier,
	tableName parser.Identifier,
	filter *Filter,
	useInternalId bool,
	forUpdate bool,
	importFormat cmd.Format,
	delimiter rune,
	delimiterPositions []int,
	jsonQuery string,
	encoding text.Encoding,
	lineBreak text.LineBreak,
	noHeader bool,
	encloseAll bool,
	jsonEscape txjson.EscapeType,
	withoutNull bool,
) (*View, error) {
	var view *View

	if filter.RecursiveTable != nil && strings.EqualFold(tableIdentifier.Literal, filter.RecursiveTable.Name.Literal) && filter.RecursiveTmpView != nil {
		view = filter.RecursiveTmpView
		if !strings.EqualFold(filter.RecursiveTable.Name.Literal, tableName.Literal) {
			view.Header.Update(tableName.Literal, nil)
		}
	} else if it, err := filter.InlineTables.Get(tableIdentifier); err == nil {
		if err = filter.Aliases.Add(tableName, ""); err != nil {
			return nil, err
		}
		view = it
		if tableIdentifier.Literal != tableName.Literal {
			view.Header.Update(tableName.Literal, nil)
		}
	} else {
		var filePath string
		var commonTableName string

		filePath = tableIdentifier.Literal
		if filter.TempViews.Exists(filePath) {
			commonTableName = parser.FormatTableName(filePath)

			pathIdent := parser.Identifier{Literal: filePath}
			if useInternalId {
				view, _ = filter.TempViews.GetWithInternalId(pathIdent)
			} else {
				view, _ = filter.TempViews.Get(pathIdent)
			}
		} else {
			filePath, err = CreateFilePath(tableIdentifier, cmd.GetFlags().Repository)
			if err != nil {
				return nil, err
			}

			if !ViewCache.Exists(filePath) {
				fileInfo, err := NewFileInfo(tableIdentifier, cmd.GetFlags().Repository, importFormat, delimiter, encoding)
				if err != nil {
					return nil, err
				}
				filePath = fileInfo.Path

				fileInfo.DelimiterPositions = delimiterPositions
				fileInfo.JsonQuery = strings.TrimSpace(jsonQuery)
				fileInfo.LineBreak = lineBreak
				fileInfo.NoHeader = noHeader
				fileInfo.EncloseAll = encloseAll
				fileInfo.JsonEscape = jsonEscape

				if !ViewCache.Exists(fileInfo.Path) || (forUpdate && !ViewCache[strings.ToUpper(fileInfo.Path)].ForUpdate) {
					ViewCache.Dispose(fileInfo.Path)

					var fp *os.File
					if forUpdate {
						h, err := file.NewHandlerForUpdate(fileInfo.Path)
						if err != nil {
							if _, ok := err.(*file.TimeoutError); ok {
								return nil, NewFileLockTimeoutError(tableIdentifier, fileInfo.Path)
							}
							return nil, NewReadFileError(tableIdentifier, err.Error())
						}
						fileInfo.Handler = h
						fp = h.FileForRead()
					} else {
						h, err := file.NewHandlerForRead(fileInfo.Path)
						if err != nil {
							if _, ok := err.(*file.TimeoutError); ok {
								return nil, NewFileLockTimeoutError(tableIdentifier, fileInfo.Path)
							}
							return nil, NewReadFileError(tableIdentifier, err.Error())
						}
						defer h.Close()
						fp = h.FileForRead()
					}

					loadView, err := loadViewFromFile(fp, fileInfo, withoutNull)
					if err != nil {
						fileInfo.Close()
						return nil, NewDataParsingError(tableIdentifier, fileInfo.Path, err.Error())
					}
					loadView.ForUpdate = forUpdate
					ViewCache.Set(loadView)
				}
			}
			commonTableName = parser.FormatTableName(filePath)

			pathIdent := parser.Identifier{Literal: filePath}
			if useInternalId {
				view, _ = ViewCache.GetWithInternalId(pathIdent)
			} else {
				view, _ = ViewCache.Get(pathIdent)
			}
		}

		if err = filter.Aliases.Add(tableName, filePath); err != nil {
			return nil, err
		}

		if !strings.EqualFold(commonTableName, tableName.Literal) {
			view.Header.Update(tableName.Literal, nil)
		}
	}
	return view, nil
}

func loadViewFromFile(fp *os.File, fileInfo *FileInfo, withoutNull bool) (*View, error) {
	switch fileInfo.Format {
	case cmd.FIXED:
		return loadViewFromFixedLengthTextFile(fp, fileInfo, withoutNull)
	case cmd.JSON:
		return loadViewFromJsonFile(fp, fileInfo)
	}
	return loadViewFromCSVFile(fp, fileInfo, withoutNull)
}

func loadViewFromFixedLengthTextFile(fp *os.File, fileInfo *FileInfo, withoutNull bool) (*View, error) {
	var err error

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)

	if fileInfo.DelimiterPositions == nil {
		d := fixedlen.NewDelimiter(r, fileInfo.Encoding)
		d.NoHeader = fileInfo.NoHeader
		d.Encoding = fileInfo.Encoding
		fileInfo.DelimiterPositions, err = d.Delimit()
		if err != nil {
			return nil, err
		}
	}

	r.Seek(0, io.SeekStart)
	reader := fixedlen.NewReader(r, fileInfo.DelimiterPositions, fileInfo.Encoding)
	reader.WithoutNull = withoutNull
	reader.Encoding = fileInfo.Encoding

	var header []string
	if !fileInfo.NoHeader {
		header, err = reader.ReadHeader()
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	records, err := readRecordSet(reader)
	if err != nil {
		return nil, err
	}

	if header == nil {
		header = make([]string, len(fileInfo.DelimiterPositions))
		for i := 0; i < len(fileInfo.DelimiterPositions); i++ {
			header[i] = "c" + strconv.Itoa(i+1)
		}
	}

	if reader.DetectedLineBreak != "" {
		fileInfo.LineBreak = reader.DetectedLineBreak
	}

	view := NewView()
	view.Header = NewHeaderWithAutofill(parser.FormatTableName(fileInfo.Path), header)
	view.RecordSet = records
	view.FileInfo = fileInfo
	return view, nil
}

func loadViewFromCSVFile(fp *os.File, fileInfo *FileInfo, withoutNull bool) (*View, error) {
	reader := csv.NewReader(fp, fileInfo.Encoding)
	reader.Delimiter = fileInfo.Delimiter
	reader.WithoutNull = withoutNull

	var err error
	var header []string
	if !fileInfo.NoHeader {
		header, err = reader.ReadHeader()
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	records, err := readRecordSet(reader)
	if err != nil {
		return nil, err
	}

	if header == nil {
		header = make([]string, reader.FieldsPerRecord)
		for i := 0; i < reader.FieldsPerRecord; i++ {
			header[i] = "c" + strconv.Itoa(i+1)
		}
	}

	if reader.DetectedLineBreak != "" {
		fileInfo.LineBreak = reader.DetectedLineBreak
	}
	fileInfo.EncloseAll = reader.EnclosedAll

	view := NewView()
	view.Header = NewHeader(parser.FormatTableName(fileInfo.Path), header)
	view.RecordSet = records
	view.FileInfo = fileInfo
	return view, nil
}

func readRecordSet(reader RecordReader) (RecordSet, error) {
	var err error
	records := make(RecordSet, 0, 1000)
	rowch := make(chan []text.RawText, 1000)
	fieldch := make(chan []value.Primary, 1000)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for {
			primaries, ok := <-fieldch
			if !ok {
				break
			}
			records = append(records, NewRecord(primaries))
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			row, ok := <-rowch
			if !ok {
				break
			}
			fields := make([]value.Primary, len(row))
			for i, v := range row {
				if v == nil {
					fields[i] = value.NewNull()
				} else {
					fields[i] = value.NewString(string(v))
				}
			}
			fieldch <- fields
		}
		close(fieldch)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			record, e := reader.Read()
			if e == io.EOF {
				break
			}
			if e != nil {
				err = e
				break
			}
			rowch <- record
		}
		close(rowch)
		wg.Done()
	}()

	wg.Wait()

	return records, err
}

func loadViewFromJsonFile(fp io.Reader, fileInfo *FileInfo) (*View, error) {
	jsonText, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	headerLabels, rows, escapeType, err := json.LoadTable(fileInfo.JsonQuery, string(jsonText))
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0, len(rows))
	for _, row := range rows {
		records = append(records, NewRecord(row))
	}

	fileInfo.JsonEscape = escapeType

	view := NewView()
	view.Header = NewHeader(parser.FormatTableName(fileInfo.Path), headerLabels)
	view.RecordSet = records
	view.FileInfo = fileInfo
	return view, nil
}

func loadDualView() *View {
	view := View{
		Header:    NewDualHeader(),
		RecordSet: make([]Record, 1),
	}
	view.RecordSet[0] = NewEmptyRecord(1)
	return &view
}

func NewViewFromGroupedRecord(filterRecord FilterRecord) *View {
	view := new(View)
	view.Header = filterRecord.View.Header
	record := filterRecord.View.RecordSet[filterRecord.RecordIndex]

	view.RecordSet = make([]Record, record.GroupLen())
	for i := 0; i < record.GroupLen(); i++ {
		view.RecordSet[i] = make(Record, view.FieldLen())
		for j, cell := range record {
			grpIdx := i
			if cell.Len() < 2 {
				grpIdx = 0
			}
			view.RecordSet[i][j] = NewCell(cell.GroupedValue(grpIdx))
		}
	}

	view.Filter = filterRecord.View.Filter

	return view
}

func (view *View) Where(clause parser.WhereClause) error {
	return view.filter(clause.Filter)
}

func (view *View) filter(condition parser.QueryExpression) error {
	results := make([]bool, view.RecordLen())

	gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
	for i := 0; i < gm.CPU; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)
			filter := NewFilterForSequentialEvaluation(view, view.Filter)

		FilterLoop:
			for i := start; i < end; i++ {
				if gm.HasError() {
					break FilterLoop
				}

				filter.Records[0].RecordIndex = i
				primary, e := filter.Evaluate(condition)
				if e != nil {
					gm.SetError(e)
					break FilterLoop
				}
				if primary.Ternary() == ternary.TRUE {
					results[i] = true
				}
			}

			gm.Done()
		}(i)
	}
	gm.Wait()

	if gm.HasError() {
		return gm.Error()
	}

	records := make(RecordSet, 0, len(results))
	for i, ok := range results {
		if ok {
			records = append(records, view.RecordSet[i])
		}
	}

	view.RecordSet = make(RecordSet, len(records))
	copy(view.RecordSet, records)
	return nil
}

func (view *View) GroupBy(clause parser.GroupByClause) error {
	return view.group(clause.Items)
}

func (view *View) group(items []parser.QueryExpression) error {
	if items == nil {
		return view.groupAll()
	}

	keys := make([]string, view.RecordLen())

	gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
	for i := 0; i < gm.CPU; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)

			filter := NewFilterForSequentialEvaluation(view, view.Filter)
			values := make([]value.Primary, len(items))

		GroupLoop:
			for i := start; i < end; i++ {
				if gm.HasError() {
					break GroupLoop
				}

				filter.Records[0].RecordIndex = i
				for j, item := range items {
					p, e := filter.Evaluate(item)
					if e != nil {
						gm.SetError(e)
						break GroupLoop
					}
					values[j] = p
				}
				keys[i] = SerializeComparisonKeys(values)
			}

			gm.Done()
		}(i)
	}
	gm.Wait()

	if gm.HasError() {
		return gm.Error()
	}

	groups := make(map[string][]int)
	groupKeys := make([]string, 0)
	for i, key := range keys {
		if _, ok := groups[key]; ok {
			groups[key] = append(groups[key], i)
		} else {
			groups[key] = []int{i}
			groupKeys = append(groupKeys, key)
		}
	}

	records := make(RecordSet, len(groupKeys))
	for i, groupKey := range groupKeys {
		record := make(Record, view.FieldLen())
		indices := groups[groupKey]

		for j := 0; j < view.FieldLen(); j++ {
			primaries := make([]value.Primary, len(indices))
			for k, idx := range indices {
				primaries[k] = view.RecordSet[idx][j].Value()
			}
			record[j] = NewGroupCell(primaries)
		}

		records[i] = record
	}

	view.RecordSet = records
	view.isGrouped = true
	for _, item := range items {
		switch item.(type) {
		case parser.FieldReference, parser.ColumnNumber:
			idx, _ := view.FieldIndex(item)
			view.Header[idx].IsGroupKey = true
		}
	}
	return nil
}

func (view *View) groupAll() error {
	if 0 < view.RecordLen() {
		records := make(RecordSet, 1)
		record := make(Record, view.FieldLen())
		for i := 0; i < view.FieldLen(); i++ {
			primaries := make([]value.Primary, len(view.RecordSet))
			for j := range view.RecordSet {
				primaries[j] = view.RecordSet[j][i].Value()
			}
			record[i] = NewGroupCell(primaries)
		}
		records[0] = record
		view.RecordSet = records
	}

	view.isGrouped = true
	return nil
}

func (view *View) Having(clause parser.HavingClause) error {
	err := view.filter(clause.Filter)
	if err != nil {
		if _, ok := err.(*NotGroupingRecordsError); ok {
			view.group(nil)
			err = view.filter(clause.Filter)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (view *View) Select(clause parser.SelectClause) error {
	var parseAllColumns = func(view *View, fields []parser.QueryExpression) []parser.QueryExpression {
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
		insert := make([]parser.QueryExpression, insertLen)
		for i, c := range columns {
			insert[i] = parser.Field{
				Object: c,
			}
		}

		list := make([]parser.QueryExpression, len(fields)-1+insertLen)
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

	var evalFields = func(view *View, fields []parser.QueryExpression) error {
		fieldsObjects := make([]parser.QueryExpression, len(fields))
		for i, f := range fields {
			fieldsObjects[i] = f.(parser.Field).Object
		}
		if err := view.ExtendRecordCapacity(fieldsObjects); err != nil {
			return err
		}

		view.selectFields = make([]int, len(fields))
		view.selectLabels = make([]string, len(fields))
		for i, f := range fields {
			field := f.(parser.Field)
			alias := ""
			if field.Alias != nil {
				alias = field.Alias.(parser.Identifier).Literal
			}
			idx, err := view.evalColumn(field.Object, alias)
			if err != nil {
				return err
			}
			view.selectFields[i] = idx
			view.selectLabels[i] = field.Name()
		}
		return nil
	}

	fields := parseAllColumns(view, clause.Fields)

	origFieldLen := view.FieldLen()
	err := evalFields(view, fields)
	if err != nil {
		if _, ok := err.(*NotGroupingRecordsError); ok {
			view.Header = view.Header[:origFieldLen]
			if 0 < view.RecordLen() && view.FieldLen() < len(view.RecordSet[0]) {
				for i := range view.RecordSet {
					view.RecordSet[i] = view.RecordSet[i][:origFieldLen]
				}
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
		view.GenerateComparisonKeys()
		records := make(RecordSet, 0, view.RecordLen())
		values := make(map[string]bool)
		for i, v := range view.RecordSet {
			if !values[view.comparisonKeysInEachRecord[i]] {
				values[view.comparisonKeysInEachRecord[i]] = true

				record := make(Record, len(view.selectFields))
				for j, idx := range view.selectFields {
					record[j] = v[idx]
				}
				records = append(records, record)
			}
		}

		hfields := NewEmptyHeader(len(view.selectFields))
		for i, idx := range view.selectFields {
			hfields[i] = view.Header[idx]
			view.selectFields[i] = i
		}

		view.Header = hfields
		view.RecordSet = records
		view.comparisonKeysInEachRecord = nil
		view.sortValuesInEachCell = nil
	}

	return nil
}

func (view *View) GenerateComparisonKeys() {
	view.comparisonKeysInEachRecord = make([]string, view.RecordLen())

	gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
	for i := 0; i < gm.CPU; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)

			var primaries []value.Primary
			if view.selectFields != nil {
				primaries = make([]value.Primary, len(view.selectFields))
			}

			for i := start; i < end; i++ {
				if view.selectFields != nil {
					for j, idx := range view.selectFields {
						primaries[j] = view.RecordSet[i][idx].Value()
					}
					view.comparisonKeysInEachRecord[i] = SerializeComparisonKeys(primaries)
				} else {
					view.comparisonKeysInEachRecord[i] = view.RecordSet[i].SerializeComparisonKeys()
				}
			}

			gm.Done()
		}(i)
	}
	gm.Wait()
}

func (view *View) SelectAllColumns() error {
	selectClause := parser.SelectClause{
		Fields: []parser.QueryExpression{
			parser.Field{Object: parser.AllColumns{}},
		},
	}
	return view.Select(selectClause)
}

func (view *View) OrderBy(clause parser.OrderByClause) error {
	orderValues := make([]parser.QueryExpression, len(clause.Items))
	for i, item := range clause.Items {
		orderValues[i] = item.(parser.OrderItem).Value
	}
	if err := view.ExtendRecordCapacity(orderValues); err != nil {
		return err
	}

	sortIndices := make([]int, len(clause.Items))
	for i, v := range clause.Items {
		oi := v.(parser.OrderItem)
		idx, err := view.evalColumn(oi.Value, "")
		if err != nil {
			return err
		}
		sortIndices[i] = idx
	}

	view.sortValuesInEachRecord = make([]SortValues, view.RecordLen())
	view.sortDirections = make([]int, len(clause.Items))
	view.sortNullPositions = make([]int, len(clause.Items))

	for i, v := range clause.Items {
		oi := v.(parser.OrderItem)
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

	gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
	for i := 0; i < gm.CPU; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)

			for i := start; i < end; i++ {
				if view.sortValuesInEachCell != nil && view.sortValuesInEachCell[i] == nil {
					view.sortValuesInEachCell[i] = make([]*SortValue, cap(view.RecordSet[i]))
				}

				sortValues := make(SortValues, len(sortIndices))
				for j, idx := range sortIndices {
					if view.sortValuesInEachCell != nil && idx < len(view.sortValuesInEachCell[i]) && view.sortValuesInEachCell[i][idx] != nil {
						sortValues[j] = view.sortValuesInEachCell[i][idx]
					} else {
						sortValues[j] = NewSortValue(view.RecordSet[i][idx].Value())
						if view.sortValuesInEachCell != nil && idx < len(view.sortValuesInEachCell[i]) {
							view.sortValuesInEachCell[i][idx] = sortValues[j]
						}
					}
				}
				view.sortValuesInEachRecord[i] = sortValues
			}

			gm.Done()
		}(i)
	}
	gm.Wait()

	sort.Sort(view)
	return nil
}

func (view *View) additionalColumns(expr parser.QueryExpression) ([]string, error) {
	list := make([]string, 0)

	switch expr.(type) {
	case parser.FieldReference, parser.ColumnNumber:
		return nil, nil
	case parser.Function:
		if udfn, err := view.Filter.Functions.Get(expr, expr.(parser.Function).Name); err == nil {
			if udfn.IsAggregate && !view.isGrouped {
				return nil, NewNotGroupingRecordsError(expr, expr.(parser.Function).Name)
			}
		}
	case parser.AggregateFunction:
		if !view.isGrouped {
			return nil, NewNotGroupingRecordsError(expr, expr.(parser.AggregateFunction).Name)
		}
	case parser.ListFunction:
		if !view.isGrouped {
			return nil, NewNotGroupingRecordsError(expr, expr.(parser.ListFunction).Name)
		}
	case parser.AnalyticFunction:
		fn := expr.(parser.AnalyticFunction)
		pvalues := fn.AnalyticClause.PartitionValues()
		ovalues := []parser.QueryExpression(nil)
		if fn.AnalyticClause.OrderByClause != nil {
			ovalues = fn.AnalyticClause.OrderByClause.(parser.OrderByClause).Items
		}

		if pvalues != nil {
			for _, pvalue := range pvalues {
				columns, err := view.additionalColumns(pvalue)
				if err != nil {
					return nil, err
				}
				for _, s := range columns {
					if !InStrSliceWithCaseInsensitive(s, list) {
						list = append(list, s)
					}
				}
			}
		}
		if ovalues != nil {
			for _, v := range ovalues {
				item := v.(parser.OrderItem)
				columns, err := view.additionalColumns(item.Value)
				if err != nil {
					return nil, err
				}
				for _, s := range columns {
					if !InStrSliceWithCaseInsensitive(s, list) {
						list = append(list, s)
					}
				}
			}
		}
	}

	if _, err := view.Header.ContainsObject(expr); err != nil {
		s := expr.String()
		if !InStrSliceWithCaseInsensitive(s, list) {
			list = append(list, s)
		}
	}

	return list, nil
}

func (view *View) ExtendRecordCapacity(exprs []parser.QueryExpression) error {
	additions := make([]string, 0)
	for _, expr := range exprs {
		columns, err := view.additionalColumns(expr)
		if err != nil {
			return err
		}
		for _, s := range columns {
			if !InStrSliceWithCaseInsensitive(s, additions) {
				additions = append(additions, s)
			}
		}
	}

	currentLen := view.FieldLen()
	fieldCap := currentLen + len(additions)

	if 0 < view.RecordLen() && fieldCap <= cap(view.RecordSet[0]) {
		return nil
	}

	gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
	for i := 0; i < gm.CPU; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)
			for i := start; i < end; i++ {
				record := make(Record, currentLen, fieldCap)
				copy(record, view.RecordSet[i])
				view.RecordSet[i] = record
			}
			gm.Done()
		}(i)
	}
	gm.Wait()
	return nil
}

func (view *View) evalColumn(obj parser.QueryExpression, alias string) (idx int, err error) {
	switch obj.(type) {
	case parser.FieldReference, parser.ColumnNumber:
		if idx, err = view.FieldIndex(obj); err != nil {
			return
		}
		if view.isGrouped && view.Header[idx].IsFromTable && !view.Header[idx].IsGroupKey {
			err = NewFieldNotGroupKeyError(obj)
			return
		}
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
				gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
				for i := 0; i < gm.CPU; i++ {
					gm.Add()
					go func(thIdx int) {
						start, end := gm.RecordRange(thIdx)
						filter := NewFilterForSequentialEvaluation(view, view.Filter)

					EvalColumnLoop:
						for i := start; i < end; i++ {
							if gm.HasError() {
								break EvalColumnLoop
							}

							filter.Records[0].RecordIndex = i

							primary, e := filter.Evaluate(obj)
							if e != nil {
								gm.SetError(e)
								break EvalColumnLoop
							}
							view.RecordSet[i] = append(view.RecordSet[i], NewCell(primary))
						}

						gm.Done()
					}(i)
				}
				gm.Wait()

				if gm.HasError() {
					err = gm.Error()
					return
				}
			}
			view.Header, idx = AddHeaderField(view.Header, parser.FormatFieldIdentifier(obj), alias)
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
			if udfn, err := view.Filter.Functions.Get(expr, expr.Name); err != nil || !udfn.IsAggregate {
				return NewFunctionNotExistError(expr, expr.Name)
			}
		}
	}

	var partitionIndices []int
	if expr.AnalyticClause.PartitionClause != nil {
		partitionExprs := expr.AnalyticClause.PartitionValues()

		partitionIndices = make([]int, len(partitionExprs))
		for i, pexpr := range partitionExprs {
			idx, err := view.evalColumn(pexpr, "")
			if err != nil {
				return err
			}
			partitionIndices[i] = idx
		}
	}

	if view.sortValuesInEachCell == nil {
		view.sortValuesInEachCell = make([][]*SortValue, view.RecordLen())
	}

	if expr.AnalyticClause.OrderByClause != nil {
		err := view.OrderBy(expr.AnalyticClause.OrderByClause.(parser.OrderByClause))
		if err != nil {
			return err
		}
	}

	err := Analyze(view, expr, partitionIndices)

	view.sortValuesInEachRecord = nil
	view.sortDirections = nil
	view.sortNullPositions = nil

	return err
}

func (view *View) Offset(clause parser.OffsetClause) error {
	val, err := view.Filter.Evaluate(clause.Value)
	if err != nil {
		return err
	}
	number := value.ToInteger(val)
	if value.IsNull(number) {
		return NewInvalidOffsetNumberError(clause)
	}
	view.offset = int(number.(value.Integer).Raw())
	if view.offset < 0 {
		view.offset = 0
	}

	if view.RecordLen() <= view.offset {
		view.RecordSet = RecordSet{}
	} else {
		view.RecordSet = view.RecordSet[view.offset:]
		records := make(RecordSet, len(view.RecordSet))
		copy(records, view.RecordSet)
		view.RecordSet = records
	}
	return nil
}

func (view *View) Limit(clause parser.LimitClause) error {
	val, err := view.Filter.Evaluate(clause.Value)
	if err != nil {
		return err
	}

	var limit int
	if clause.IsPercentage() {
		number := value.ToFloat(val)
		if value.IsNull(number) {
			return NewInvalidLimitPercentageError(clause)
		}
		percentage := number.(value.Float).Raw()
		if 100 < percentage {
			limit = 100
		} else if percentage < 0 {
			limit = 0
		} else {
			limit = int(math.Ceil(float64(view.RecordLen()+view.offset) * percentage / 100))
		}
	} else {
		number := value.ToInteger(val)
		if value.IsNull(number) {
			return NewInvalidLimitNumberError(clause)
		}
		limit = int(number.(value.Integer).Raw())
		if limit < 0 {
			limit = 0
		}
	}

	if view.RecordLen() <= limit {
		return nil
	}

	if clause.IsWithTies() && view.sortValuesInEachRecord != nil {
		bottomSortValues := view.sortValuesInEachRecord[limit-1]
		for limit < view.RecordLen() {
			if !bottomSortValues.EquivalentTo(view.sortValuesInEachRecord[limit]) {
				break
			}
			limit++
		}
	}

	view.RecordSet = view.RecordSet[:limit]
	records := make(RecordSet, view.RecordLen())
	copy(records, view.RecordSet)
	view.RecordSet = records
	return nil
}

func (view *View) InsertValues(fields []parser.QueryExpression, list []parser.QueryExpression) (int, error) {
	valuesList := make([][]value.Primary, len(list))

	for i, item := range list {
		rv := item.(parser.RowValue)
		values, err := view.Filter.evalRowValue(rv)
		if err != nil {
			return 0, err
		}
		if len(fields) != len(values) {
			return 0, NewInsertRowValueLengthError(rv, len(fields))
		}

		valuesList[i] = values
	}

	return view.insert(fields, valuesList)
}

func (view *View) InsertFromQuery(fields []parser.QueryExpression, query parser.SelectQuery) (int, error) {
	insertView, err := Select(query, view.Filter)
	if err != nil {
		return 0, err
	}
	if len(fields) != insertView.FieldLen() {
		return 0, NewInsertSelectFieldLengthError(query, len(fields))
	}

	valuesList := make([][]value.Primary, insertView.RecordLen())

	for i, record := range insertView.RecordSet {
		values := make([]value.Primary, insertView.FieldLen())
		for j, cell := range record {
			values[j] = cell.Value()
		}
		valuesList[i] = values
	}

	return view.insert(fields, valuesList)
}

func (view *View) insert(fields []parser.QueryExpression, valuesList [][]value.Primary) (int, error) {
	var valueIndex = func(i int, list []int) int {
		for j, v := range list {
			if i == v {
				return j
			}
		}
		return -1
	}

	var insertRecords int

	fieldIndices, err := view.FieldIndices(fields)
	if err != nil {
		return insertRecords, err
	}

	records := make([]Record, len(valuesList))
	for i, values := range valuesList {
		record := make(Record, view.FieldLen())
		for j := 0; j < view.FieldLen(); j++ {
			idx := valueIndex(j, fieldIndices)
			if idx < 0 {
				record[j] = NewCell(value.NewNull())
			} else {
				record[j] = NewCell(values[idx])
			}
		}
		records[i] = record
	}

	view.RecordSet = append(view.RecordSet, records...)
	return len(valuesList), nil
}

func (view *View) Fix() {
	resize := false
	if len(view.selectFields) < view.FieldLen() {
		resize = true
	} else {
		for i := 0; i < view.FieldLen(); i++ {
			if view.selectFields[i] != i {
				resize = true
				break
			}
		}
	}

	if resize {
		gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
		for i := 0; i < gm.CPU; i++ {
			gm.Add()
			go func(thIdx int) {
				start, end := gm.RecordRange(thIdx)

				for i := start; i < end; i++ {
					record := make(Record, len(view.selectFields))
					for j, idx := range view.selectFields {
						if 1 < view.RecordSet[i].GroupLen() {
							record[j] = NewCell(view.RecordSet[i][idx].Value())
						} else {
							record[j] = view.RecordSet[i][idx]
						}
					}
					view.RecordSet[i] = record
				}

				gm.Done()
			}(i)
		}
		gm.Wait()
	}

	hfields := NewEmptyHeader(len(view.selectFields))

	colNumber := 0
	for i, idx := range view.selectFields {
		colNumber++

		hfields[i] = view.Header[idx]
		hfields[i].Aliases = nil
		hfields[i].Number = colNumber
		hfields[i].IsFromTable = true
		hfields[i].IsJoinColumn = false
		hfields[i].IsGroupKey = false

		if 0 < len(view.selectLabels) {
			hfields[i].Column = view.selectLabels[i]
		}
	}

	view.Header = hfields
	view.Filter = nil
	view.selectFields = nil
	view.selectLabels = nil
	view.isGrouped = false
	view.comparisonKeysInEachRecord = nil
	view.sortValuesInEachCell = nil
	view.sortValuesInEachRecord = nil
	view.sortDirections = nil
	view.sortNullPositions = nil
	view.offset = 0
}

func (view *View) Union(calcView *View, all bool) {
	view.RecordSet = append(view.RecordSet, calcView.RecordSet...)
	view.FileInfo = nil

	if !all {
		view.GenerateComparisonKeys()

		records := make(RecordSet, 0, view.RecordLen())
		values := make(map[string]bool)

		for i, key := range view.comparisonKeysInEachRecord {
			if !values[key] {
				values[key] = true
				records = append(records, view.RecordSet[i])
			}
		}

		view.RecordSet = records
		view.comparisonKeysInEachRecord = nil
	}
}

func (view *View) Except(calcView *View, all bool) {
	view.GenerateComparisonKeys()
	calcView.GenerateComparisonKeys()

	keys := make(map[string]bool)
	for _, key := range calcView.comparisonKeysInEachRecord {
		if !keys[key] {
			keys[key] = true
		}
	}

	distinctKeys := make(map[string]bool)
	records := make(RecordSet, 0, view.RecordLen())
	for i, key := range view.comparisonKeysInEachRecord {
		if !keys[key] {
			if !all {
				if distinctKeys[key] {
					continue
				}
				distinctKeys[key] = true
			}
			records = append(records, view.RecordSet[i])
		}
	}
	view.RecordSet = records
	view.FileInfo = nil
	view.comparisonKeysInEachRecord = nil
}

func (view *View) Intersect(calcView *View, all bool) {
	view.GenerateComparisonKeys()
	calcView.GenerateComparisonKeys()

	keys := make(map[string]bool)
	for _, key := range calcView.comparisonKeysInEachRecord {
		if !keys[key] {
			keys[key] = true
		}
	}

	distinctKeys := make(map[string]bool)
	records := make(RecordSet, 0, view.RecordLen())
	for i, key := range view.comparisonKeysInEachRecord {
		if _, ok := keys[key]; ok {
			if !all {
				if distinctKeys[key] {
					continue
				}
				distinctKeys[key] = true
			}
			records = append(records, view.RecordSet[i])
		}
	}
	view.RecordSet = records
	view.FileInfo = nil
	view.comparisonKeysInEachRecord = nil
}

func (view *View) ListValuesForAggregateFunctions(expr parser.QueryExpression, arg parser.QueryExpression, distinct bool, filter *Filter) ([]value.Primary, error) {
	list := make([]value.Primary, view.RecordLen())

	gm := NewGoroutineManager(view.RecordLen(), MinimumRequiredForParallelRoutine)
	for i := 0; i < gm.CPU; i++ {
		gm.Add()
		go func(thIdx int) {
			start, end := gm.RecordRange(thIdx)
			filter := NewFilterForSequentialEvaluation(view, filter)

		ListAggregateFunctionLoop:
			for i := start; i < end; i++ {
				if gm.HasError() {
					break ListAggregateFunctionLoop
				}

				filter.Records[0].RecordIndex = i
				p, e := filter.Evaluate(arg)
				if e != nil {
					if _, ok := e.(*NotGroupingRecordsError); ok {
						gm.SetError(NewNestedAggregateFunctionsError(expr))
					} else {
						gm.SetError(e)
					}
					break ListAggregateFunctionLoop
				}
				list[i] = p
			}

			gm.Done()
		}(i)
	}
	gm.Wait()

	if gm.HasError() {
		return nil, gm.Error()
	}

	if distinct {
		list = Distinguish(list)
	}

	return list, nil
}

func (view *View) RestoreHeaderReferences() {
	view.Header.Update(parser.FormatTableName(view.FileInfo.Path), nil)
}

func (view *View) FieldIndex(fieldRef parser.QueryExpression) (int, error) {
	if number, ok := fieldRef.(parser.ColumnNumber); ok {
		return view.Header.ContainsNumber(number)
	}
	return view.Header.Contains(fieldRef.(parser.FieldReference))
}

func (view *View) FieldIndices(fields []parser.QueryExpression) ([]int, error) {
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

func (view *View) FieldViewName(fieldRef parser.QueryExpression) (string, error) {
	idx, err := view.FieldIndex(fieldRef)
	if err != nil {
		return "", err
	}
	return view.Header[idx].View, nil
}

func (view *View) InternalRecordId(ref string, recordIndex int) (int, error) {
	idx, err := view.Header.ContainsInternalId(ref)
	if err != nil {
		return -1, NewInternalRecordIdNotExistError()
	}
	internalId, ok := view.RecordSet[recordIndex][idx].Value().(value.Integer)
	if !ok {
		return -1, NewInternalRecordIdEmptyError()
	}
	return int(internalId.Raw()), nil
}

func (view *View) FieldLen() int {
	return view.Header.Len()
}

func (view *View) RecordLen() int {
	return view.Len()
}

func (view *View) Len() int {
	return len(view.RecordSet)
}

func (view *View) Swap(i, j int) {
	view.RecordSet[i], view.RecordSet[j] = view.RecordSet[j], view.RecordSet[i]
	view.sortValuesInEachRecord[i], view.sortValuesInEachRecord[j] = view.sortValuesInEachRecord[j], view.sortValuesInEachRecord[i]
	if view.sortValuesInEachCell != nil {
		view.sortValuesInEachCell[i], view.sortValuesInEachCell[j] = view.sortValuesInEachCell[j], view.sortValuesInEachCell[i]
	}
}

func (view *View) Less(i, j int) bool {
	return view.sortValuesInEachRecord[i].Less(view.sortValuesInEachRecord[j], view.sortDirections, view.sortNullPositions)
}

func (view *View) Copy() *View {
	header := view.Header.Copy()
	records := view.RecordSet.Copy()

	return &View{
		Header:    header,
		RecordSet: records,
		FileInfo:  view.FileInfo,
		ForUpdate: view.ForUpdate,
	}
}
