package query

import (
	"bytes"
	"context"
	gojson "encoding/json"
	"fmt"
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
	"github.com/mithrandie/go-text/ltsv"
	"github.com/mithrandie/ternary"
)

const fileLoadingPreparedRecordSetCap = 300
const fileLoadingBuffer = 300

type ApplyView struct {
	View     *View
	JoinExpr parser.Join
}

type RecordReader interface {
	Read() ([]text.RawText, error)
}

type View struct {
	Header    Header
	RecordSet RecordSet
	FileInfo  *FileInfo

	selectFields []int
	selectLabels []string
	isGrouped    bool

	comparisonKeysInEachRecord []string
	sortValuesInEachCell       [][]*SortValue
	sortValuesInEachRecord     []SortValues
	sortDirections             []int
	sortNullPositions          []int

	tempRecord Record

	offset int
}

func NewView() *View {
	return &View{}
}

func LoadView(ctx context.Context, scope *ReferenceScope, tables []parser.QueryExpression, forUpdate bool, useInternalId bool) (*View, error) {
	if tables == nil {
		var obj parser.QueryExpression
		if scope.Tx.Session.CanReadStdin {
			obj = parser.Stdin{}
		} else {
			obj = parser.Dual{}
		}
		tables = []parser.QueryExpression{parser.Table{Object: obj}}
	}

	table := tables[0]

	for i := 1; i < len(tables); i++ {
		table = parser.Table{
			Object: parser.Join{
				Table:     table,
				JoinTable: tables[i],
				JoinType:  parser.Token{Token: parser.CROSS},
			},
		}
	}

	view, err := loadView(ctx, scope, table, forUpdate, useInternalId)
	return view, err
}

func LoadViewFromTableIdentifier(ctx context.Context, scope *ReferenceScope, table parser.QueryExpression, forUpdate bool, useInternalId bool) (*View, error) {
	tables := []parser.QueryExpression{
		parser.Table{Object: table},
	}

	return LoadView(ctx, scope, tables, forUpdate, useInternalId)
}

func loadView(ctx context.Context, scope *ReferenceScope, tableExpr parser.QueryExpression, forUpdate bool, useInternalId bool) (view *View, err error) {
	if parentheses, ok := tableExpr.(parser.Parentheses); ok {
		return loadView(ctx, scope, parentheses.Expr, forUpdate, useInternalId)
	}

	table := tableExpr.(parser.Table)
	tableName := table.Name()

	switch table.Object.(type) {
	case parser.Dual:
		view = loadDualView()
	case parser.TableObject:
		tableObject := table.Object.(parser.TableObject)
		options := scope.Tx.Flags.ImportOptions.Copy()

		var felem value.Primary
		if tableObject.FormatElement != nil {
			p, err := Evaluate(ctx, scope, tableObject.FormatElement)
			if err != nil {
				return nil, err
			}
			felem = value.ToString(p)
			defer value.Discard(felem)
		}

		encodingIdx := 0
		noHeaderIdx := 1
		withoutNullIdx := 2

		switch tableObject.Type.Token {
		case parser.CSV:
			if felem == nil {
				return nil, NewTableObjectInvalidArgumentError(tableObject, "delimiter is not specified")
			}
			if value.IsNull(felem) {
				return nil, NewTableObjectInvalidDelimiterError(tableObject, tableObject.FormatElement.String())
			}
			s := felem.(*value.String).Raw()
			d := []rune(s)
			if 1 != len(d) {
				return nil, NewTableObjectInvalidDelimiterError(tableObject, tableObject.FormatElement.String())
			}
			if 3 < len(tableObject.Args) {
				return nil, NewTableObjectArgumentsLengthError(tableObject, 5)
			}
			options.Delimiter = d[0]
			if options.Delimiter == '\t' {
				options.Format = cmd.TSV
			} else {
				options.Format = cmd.CSV
			}
		case parser.FIXED:
			if felem == nil {
				return nil, NewTableObjectInvalidArgumentError(tableObject, "delimiter positions are not specified")
			}
			if value.IsNull(felem) {
				return nil, NewTableObjectInvalidDelimiterPositionsError(tableObject, tableObject.FormatElement.String())
			}
			s := felem.(*value.String).Raw()

			var positions []int
			if !strings.EqualFold("SPACES", s) {
				if strings.HasPrefix(s, "s[") || strings.HasPrefix(s, "S[") {
					options.SingleLine = true
					s = s[1:]
				}
				err = gojson.Unmarshal([]byte(s), &positions)
				if err != nil {
					return nil, NewTableObjectInvalidDelimiterPositionsError(tableObject, tableObject.FormatElement.String())
				}
			}
			if 3 < len(tableObject.Args) {
				return nil, NewTableObjectArgumentsLengthError(tableObject, 5)
			}
			options.DelimiterPositions = positions
			options.Format = cmd.FIXED
		case parser.JSON:
			if felem == nil {
				return nil, NewTableObjectInvalidArgumentError(tableObject, "json query is not specified")
			}
			if value.IsNull(felem) {
				return nil, NewTableObjectInvalidJsonQueryError(tableObject, tableObject.FormatElement.String())
			}
			if 0 < len(tableObject.Args) {
				return nil, NewTableObjectJsonArgumentsLengthError(tableObject, 2)
			}
			options.JsonQuery = felem.(*value.String).Raw()
			options.Format = cmd.JSON
			options.Encoding = text.UTF8
		case parser.LTSV:
			if 2 < len(tableObject.Args) {
				return nil, NewTableObjectJsonArgumentsLengthError(tableObject, 3)
			}
			options.Format = cmd.LTSV
			withoutNullIdx, noHeaderIdx = noHeaderIdx, withoutNullIdx
		default:
			return nil, NewInvalidTableObjectError(tableObject, tableObject.Type.Literal)
		}

		args := make([]value.Primary, 3)
		defer func() {
			for i := range args {
				if args[i] != nil {
					value.Discard(args[i])
				}
			}
		}()

		for i, a := range tableObject.Args {
			if pt, ok := a.(parser.PrimitiveType); ok && value.IsNull(pt.Value) {
				continue
			}

			var p value.Primary = value.NewNull()
			if fr, ok := a.(parser.FieldReference); ok {
				a = parser.NewStringValue(fr.Column.Literal)
			}
			if pv, err := Evaluate(ctx, scope, a); err == nil {
				p = pv
			}

			switch i {
			case encodingIdx:
				v := value.ToString(p)
				if !value.IsNull(v) {
					args[i] = v
				} else {
					return nil, NewTableObjectInvalidArgumentError(tableObject, fmt.Sprintf("cannot be converted as a encoding value: %s", tableObject.Args[encodingIdx].String()))
				}
			case noHeaderIdx:
				v := value.ToBoolean(p)
				if !value.IsNull(v) {
					args[i] = v
				} else {
					return nil, NewTableObjectInvalidArgumentError(tableObject, fmt.Sprintf("cannot be converted as a no-header value: %s", tableObject.Args[noHeaderIdx].String()))
				}
			case withoutNullIdx:
				v := value.ToBoolean(p)
				if !value.IsNull(v) {
					args[i] = v
				} else {
					return nil, NewTableObjectInvalidArgumentError(tableObject, fmt.Sprintf("cannot be converted as a without-null value: %s", tableObject.Args[withoutNullIdx].String()))
				}
			}
		}

		if args[encodingIdx] != nil {
			if options.Encoding, err = cmd.ParseEncoding(args[0].(*value.String).Raw()); err != nil {
				return nil, NewTableObjectInvalidArgumentError(tableObject, err.Error())
			}
		}
		if args[noHeaderIdx] != nil {
			options.NoHeader = args[noHeaderIdx].(*value.Boolean).Raw()
		}
		if args[withoutNullIdx] != nil {
			options.WithoutNull = args[withoutNullIdx].(*value.Boolean).Raw()
		}

		view, err = loadObject(
			ctx,
			scope,
			table.Object.(parser.TableObject).Path,
			tableName,
			forUpdate,
			useInternalId,
			options,
		)
		if err != nil {
			return nil, err
		}

	case parser.Identifier, parser.Stdin:
		options := scope.Tx.Flags.ImportOptions.Copy()
		options.Format = cmd.AutoSelect

		view, err = loadObject(
			ctx,
			scope,
			table.Object,
			tableName,
			forUpdate,
			useInternalId,
			options,
		)
		if err != nil {
			return nil, err
		}
	case parser.Join:
		join := table.Object.(parser.Join)
		view, err = loadView(ctx, scope, join.Table, forUpdate, useInternalId)
		if err != nil {
			return nil, err
		}

		if t, ok := join.JoinTable.(parser.Table); ok && !t.Lateral.IsEmpty() {
			switch join.Direction.Token {
			case parser.RIGHT, parser.FULL:
				return nil, NewIncorrectLateralUsageError(t)
			}

			joinTableName := t.Name()

			subquery := t.Object.(parser.Subquery)
			var hfields Header
			resultSetList := make([]RecordSet, view.RecordLen())

			if err := EvaluateSequentially(ctx, scope, view, func(seqScope *ReferenceScope, rIdx int) error {
				appliedView, err := Select(ctx, seqScope, subquery.Query)
				if err != nil {
					return err
				}

				if 0 < len(joinTableName.Literal) {
					if err = appliedView.Header.Update(joinTableName.Literal, nil); err != nil {
						return err
					}
				}

				calcView := NewView()
				calcView.Header = view.Header.Copy()
				calcView.RecordSet = RecordSet{view.RecordSet[rIdx].Copy()}
				if err = joinViews(ctx, scope, calcView, appliedView, join); err != nil {
					return err
				}

				if rIdx == 0 {
					hfields = calcView.Header
				}
				resultSetList[rIdx] = calcView.RecordSet
				return nil
			}); err != nil {
				return nil, err
			}

			resultSet := make(RecordSet, 0, view.RecordLen())
			for i := range resultSetList {
				resultSet = append(resultSet, resultSetList[i]...)
			}

			view.Header = hfields
			view.RecordSet = resultSet

		} else {
			joinView, err := loadView(ctx, scope, join.JoinTable, forUpdate, useInternalId)
			if err != nil {
				return nil, err
			}

			if err = joinViews(ctx, scope, view, joinView, join); err != nil {
				return nil, err
			}
		}

	case parser.JsonQuery:
		if 0 < len(tableName.Literal) {
			if err := scope.AddAlias(tableName, ""); err != nil {
				return nil, err
			}
		}

		jsonQuery := table.Object.(parser.JsonQuery)

		queryValue, err := Evaluate(ctx, scope, jsonQuery.Query)
		if err != nil {
			return nil, err
		}
		jq := value.ToString(queryValue)
		if value.IsNull(jq) {
			return nil, NewEmptyJsonQueryError(jsonQuery)
		}
		jqStr := jq.(*value.String).Raw()
		value.Discard(jq)

		var reader io.Reader

		if jsonPath, ok := jsonQuery.JsonText.(parser.Identifier); ok {
			fpath, err := SearchJsonFilePath(jsonPath, scope.Tx.Flags.Repository)
			if err != nil {
				return nil, err
			}

			h, err := file.NewHandlerForRead(ctx, scope.Tx.FileContainer, fpath, scope.Tx.WaitTimeout, scope.Tx.RetryDelay)
			if err != nil {
				jsonPath.Literal = fpath
				return nil, ConvertFileHandlerError(err, jsonPath)
			}
			defer func() {
				err = appendCompositeError(err, scope.Tx.FileContainer.Close(h))
			}()
			reader = h.File()
		} else {
			jsonTextValue, err := Evaluate(ctx, scope, jsonQuery.JsonText)
			if err != nil {
				return nil, err
			}
			jsonText := value.ToString(jsonTextValue)

			if value.IsNull(jsonText) {
				return nil, NewEmptyJsonTableError(jsonQuery)
			}

			reader = strings.NewReader(jsonText.(*value.String).Raw())
			value.Discard(jsonText)
		}

		fileInfo := &FileInfo{
			Path:      tableName.Literal,
			Format:    cmd.JSON,
			JsonQuery: jqStr,
			Encoding:  text.UTF8,
			LineBreak: scope.Tx.Flags.ExportOptions.LineBreak,
			ViewType:  ViewTypeTemporaryTable,
		}

		view, err = loadViewFromJsonFile(reader, fileInfo, jsonQuery)
		if err != nil {
			if _, ok := err.(Error); !ok {
				err = NewLoadJsonError(jsonQuery, err.Error())
			}
			return nil, err
		}

	case parser.Subquery:
		subquery := table.Object.(parser.Subquery)
		view, err = Select(ctx, scope, subquery.Query)
		if err != nil {
			return nil, err
		}

		if 0 < len(tableName.Literal) {
			if err := scope.AddAlias(tableName, ""); err != nil {
				return nil, err
			}

			if err = view.Header.Update(tableName.Literal, nil); err != nil {
				return nil, err
			}
		}
	}

	return view, err
}

func joinViews(ctx context.Context, scope *ReferenceScope, view *View, joinView *View, join parser.Join) error {
	condition, includeFields, excludeFields, err := ParseJoinCondition(join, view, joinView)
	if err != nil {
		return err
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
		if err = CrossJoin(ctx, scope, view, joinView); err != nil {
			return err
		}
	case parser.INNER:
		if err = InnerJoin(ctx, scope, view, joinView, condition); err != nil {
			return err
		}
	case parser.OUTER:
		if err = OuterJoin(ctx, scope, view, joinView, condition, join.Direction.Token); err != nil {
			return err
		}
	}

	if includeFields != nil {
		includeIndices := NewUintPool(len(includeFields), LimitToUseUintSlicePool)
		excludeIndices := NewUintPool(view.FieldLen()-len(includeFields), LimitToUseUintSlicePool)
		alternatives := make(map[int]int)

		for i := range includeFields {
			idx, _ := view.Header.SearchIndex(includeFields[i])
			includeIndices.Add(uint(idx))

			eidx, _ := view.Header.SearchIndex(excludeFields[i])
			excludeIndices.Add(uint(eidx))

			alternatives[idx] = eidx
		}

		fieldIndices := make([]int, 0, view.FieldLen()-excludeIndices.Len())
		header := make(Header, 0, view.FieldLen()-excludeIndices.Len())
		_ = includeIndices.Range(func(_ int, fidx uint) error {
			view.Header[fidx].View = ""
			view.Header[fidx].Number = 0
			view.Header[fidx].IsJoinColumn = true
			header = append(header, view.Header[fidx])
			fieldIndices = append(fieldIndices, int(fidx))
			return nil
		})
		for i := range view.Header {
			if excludeIndices.Exists(uint(i)) || includeIndices.Exists(uint(i)) {
				continue
			}
			header = append(header, view.Header[i])
			fieldIndices = append(fieldIndices, i)
		}
		view.Header = header
		fieldLen := len(fieldIndices)

		if err = NewGoroutineTaskManager(view.RecordLen(), -1, scope.Tx.Flags.CPU).Run(ctx, func(index int) error {
			record := make(Record, fieldLen)
			for i, idx := range fieldIndices {
				if includeIndices.Exists(uint(idx)) && value.IsNull(view.RecordSet[index][idx][0]) {
					record[i] = view.RecordSet[index][alternatives[idx]]
				} else {
					record[i] = view.RecordSet[index][idx]
				}
			}
			view.RecordSet[index] = record
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func loadStdin(ctx context.Context, scope *ReferenceScope, fileInfo *FileInfo, stdin parser.Stdin, tableName parser.Identifier, forUpdate bool, useInternalId bool) (*View, error) {
	scope.Tx.viewLoadingMutex.Lock()
	defer scope.Tx.viewLoadingMutex.Unlock()

	var err error
	view, ok := scope.Global().temporaryTables.Load(stdin.String())
	if !ok || (forUpdate && !view.FileInfo.ForUpdate) {
		if forUpdate {
			if err = scope.Tx.LockStdinContext(ctx); err != nil {
				return nil, err
			}
		} else {
			if err = scope.Tx.RLockStdinContext(ctx); err != nil {
				return nil, err
			}
			defer scope.Tx.RUnlockStdin()
		}
		view, err = scope.Tx.Session.GetStdinView(ctx, scope.Tx.Flags, fileInfo, stdin)
		if err != nil {
			return nil, err
		}
		scope.Global().temporaryTables.Set(view)
	}

	pathIdent := parser.Identifier{Literal: stdin.String()}
	if useInternalId {
		if view, err = scope.Global().temporaryTables.GetWithInternalId(ctx, pathIdent, scope.Tx.Flags); err != nil {
			if err == errTableNotLoaded {
				err = NewUndeclaredTemporaryTableError(pathIdent)
			}
			return nil, err
		}
	} else {
		if view, err = scope.Global().temporaryTables.Get(pathIdent); err != nil {
			if err == errTableNotLoaded {
				err = NewUndeclaredTemporaryTableError(pathIdent)
			}
			return nil, err
		}
	}

	if err = scope.AddAlias(tableName, view.FileInfo.Path); err != nil {
		return nil, err
	}

	if !strings.EqualFold(stdin.String(), tableName.Literal) {
		if err = view.Header.Update(tableName.Literal, nil); err != nil {
			return nil, err
		}
	}

	return view, nil
}

func loadObject(
	ctx context.Context,
	scope *ReferenceScope,
	tableExpr parser.QueryExpression,
	tableName parser.Identifier,
	forUpdate bool,
	useInternalId bool,
	options cmd.ImportOptions,
) (*View, error) {
	if stdin, ok := tableExpr.(parser.Stdin); ok {
		if options.Format == cmd.AutoSelect {
			options.Format = scope.Tx.Flags.ImportOptions.Format
		}

		if options.Format == cmd.TSV {
			options.Delimiter = '\t'
		}

		fileInfo := &FileInfo{
			Path:               stdin.String(),
			Format:             options.Format,
			Delimiter:          options.Delimiter,
			DelimiterPositions: options.DelimiterPositions,
			SingleLine:         options.SingleLine,
			JsonQuery:          options.JsonQuery,
			Encoding:           options.Encoding,
			LineBreak:          scope.Tx.Flags.ExportOptions.LineBreak,
			NoHeader:           options.NoHeader,
			ViewType:           ViewTypeStdin,
		}
		return loadStdin(ctx, scope, fileInfo, stdin, tableName, forUpdate, useInternalId)
	}

	tableIdentifier := tableExpr.(parser.Identifier)

	if scope.RecursiveTable != nil && strings.EqualFold(tableIdentifier.Literal, scope.RecursiveTable.Name.Literal) && scope.RecursiveTmpView != nil {
		view := scope.RecursiveTmpView.Copy()
		if !strings.EqualFold(scope.RecursiveTable.Name.Literal, tableName.Literal) {
			if err := view.Header.Update(tableName.Literal, nil); err != nil {
				return nil, err
			}
		}
		return view, nil
	}

	if scope.InlineTableExists(tableIdentifier) {
		if err := scope.AddAlias(tableName, ""); err != nil {
			return nil, err
		}

		view, _ := scope.GetInlineTable(tableIdentifier)
		if !strings.EqualFold(tableIdentifier.Literal, tableName.Literal) {
			if err := view.Header.Update(tableName.Literal, nil); err != nil {
				return nil, err
			}
		}
		return view, nil
	}

	filePath := tableIdentifier.Literal
	if scope.TemporaryTableExists(filePath) {
		if err := scope.AddAlias(tableName, filePath); err != nil {
			return nil, err
		}

		var view *View
		var err error
		pathIdent := parser.Identifier{Literal: filePath}
		if useInternalId {
			if view, err = scope.GetTemporaryTableWithInternalId(ctx, pathIdent, scope.Tx.Flags); err != nil {
				return nil, err
			}
		} else {
			if view, err = scope.GetTemporaryTable(pathIdent); err != nil {
				return nil, err
			}
		}

		if !strings.EqualFold(parser.FormatTableName(filePath), tableName.Literal) {
			if err := view.Header.Update(tableName.Literal, nil); err != nil {
				return nil, err
			}
		}

		return view, nil
	}

	filePath, err := cacheViewFromFile(
		ctx,
		scope,
		tableIdentifier,
		forUpdate,
		options,
	)
	if err != nil {
		return nil, err
	}

	var view *View
	pathIdent := parser.Identifier{Literal: filePath}
	if useInternalId {
		if view, err = scope.Tx.cachedViews.GetWithInternalId(ctx, pathIdent, scope.Tx.Flags); err != nil {
			if err == errTableNotLoaded {
				err = NewTableNotLoadedError(pathIdent)
			}
			return nil, err
		}
	} else {
		if view, err = scope.Tx.cachedViews.Get(pathIdent); err != nil {
			return nil, NewTableNotLoadedError(pathIdent)
		}
	}

	if err = scope.AddAlias(tableName, filePath); err != nil {
		return nil, err
	}

	if !strings.EqualFold(parser.FormatTableName(filePath), tableName.Literal) {
		if err = view.Header.Update(tableName.Literal, nil); err != nil {
			return nil, err
		}
	}
	return view, nil
}

func cacheViewFromFile(
	ctx context.Context,
	scope *ReferenceScope,
	tableIdentifier parser.Identifier,
	forUpdate bool,
	options cmd.ImportOptions,
) (string, error) {
	scope.Tx.viewLoadingMutex.Lock()
	defer scope.Tx.viewLoadingMutex.Unlock()

	filePath, cacheExists := scope.LoadFilePath(tableIdentifier.Literal)
	if !cacheExists {
		p, err := CreateFilePath(tableIdentifier, scope.Tx.Flags.Repository)
		if err != nil {
			return filePath, NewIOError(tableIdentifier, err.Error())
		}
		filePath = p
	}

	view, ok := scope.Tx.cachedViews.Load(filePath)
	if !ok || (forUpdate && !view.FileInfo.ForUpdate) {
		fileInfo, err := NewFileInfo(tableIdentifier, scope.Tx.Flags.Repository, options, scope.Tx.Flags.ImportOptions.Format)
		if err != nil {
			return filePath, err
		}
		filePath = fileInfo.Path

		view, ok = scope.Tx.cachedViews.Load(filePath)
		if !ok || (forUpdate && !view.FileInfo.ForUpdate) {
			fileInfo.DelimiterPositions = options.DelimiterPositions
			fileInfo.SingleLine = options.SingleLine
			fileInfo.JsonQuery = cmd.TrimSpace(options.JsonQuery)
			fileInfo.LineBreak = scope.Tx.Flags.ExportOptions.LineBreak
			fileInfo.NoHeader = options.NoHeader
			fileInfo.EncloseAll = scope.Tx.Flags.ExportOptions.EncloseAll
			fileInfo.JsonEscape = scope.Tx.Flags.ExportOptions.JsonEscape

			if ok {
				fileInfo = view.FileInfo
			}

			if err = scope.Tx.cachedViews.Dispose(scope.Tx.FileContainer, fileInfo.Path); err != nil {
				return filePath, err
			}

			var fp *os.File
			if forUpdate {
				h, err := file.NewHandlerForUpdate(ctx, scope.Tx.FileContainer, fileInfo.Path, scope.Tx.WaitTimeout, scope.Tx.RetryDelay)
				if err != nil {
					tableIdentifier.Literal = fileInfo.Path
					return filePath, ConvertFileHandlerError(err, tableIdentifier)
				}
				fileInfo.Handler = h
				fp = h.File()
			} else {
				h, err := file.NewHandlerForRead(ctx, scope.Tx.FileContainer, fileInfo.Path, scope.Tx.WaitTimeout, scope.Tx.RetryDelay)
				if err != nil {
					tableIdentifier.Literal = fileInfo.Path
					return filePath, ConvertFileHandlerError(err, tableIdentifier)
				}
				defer func() {
					err = appendCompositeError(err, scope.Tx.FileContainer.Close(h))
				}()
				fp = h.File()
			}

			loadView, err := loadViewFromFile(ctx, scope.Tx.Flags, fp, fileInfo, options.WithoutNull, tableIdentifier)
			if err != nil {
				if _, ok := err.(Error); !ok {
					err = NewDataParsingError(tableIdentifier, fileInfo.Path, err.Error())
				}
				return filePath, appendCompositeError(err, scope.Tx.FileContainer.Close(fileInfo.Handler))
			}
			loadView.FileInfo.ForUpdate = forUpdate
			scope.Tx.cachedViews.Set(loadView)
		}
	}
	if !cacheExists {
		scope.StoreFilePath(tableIdentifier.Literal, filePath)
	}
	return filePath, nil
}

func loadViewFromFile(ctx context.Context, flags *cmd.Flags, fp io.ReadSeeker, fileInfo *FileInfo, withoutNull bool, expr parser.QueryExpression) (*View, error) {
	switch fileInfo.Format {
	case cmd.FIXED:
		return loadViewFromFixedLengthTextFile(ctx, fp, fileInfo, withoutNull, expr)
	case cmd.LTSV:
		return loadViewFromLTSVFile(ctx, flags, fp, fileInfo, withoutNull, expr)
	case cmd.JSON:
		return loadViewFromJsonFile(fp, fileInfo, expr)
	}
	return loadViewFromCSVFile(ctx, fp, fileInfo, withoutNull, expr)
}

func loadViewFromFixedLengthTextFile(ctx context.Context, fp io.ReadSeeker, fileInfo *FileInfo, withoutNull bool, expr parser.QueryExpression) (*View, error) {
	enc, err := text.DetectInSpecifiedEncoding(fp, fileInfo.Encoding)
	if err != nil {
		return nil, NewCannotDetectFileEncodingError(expr)
	}
	fileInfo.Encoding = enc

	var r io.Reader

	if fileInfo.DelimiterPositions == nil {
		data, err := ioutil.ReadAll(fp)
		if err != nil {
			return nil, NewIOError(expr, err.Error())
		}
		br := bytes.NewReader(data)

		d, err := fixedlen.NewDelimiter(br, fileInfo.Encoding)
		if err != nil {
			return nil, err
		}
		d.NoHeader = fileInfo.NoHeader
		d.Encoding = fileInfo.Encoding
		fileInfo.DelimiterPositions, err = d.Delimit()
		if err != nil {
			return nil, err
		}

		if _, err = br.Seek(0, io.SeekStart); err != nil {
			return nil, NewSystemError(err.Error())
		}
		r = br
	} else {
		r = fp
	}

	reader, err := fixedlen.NewReader(r, fileInfo.DelimiterPositions, fileInfo.Encoding)
	if err != nil {
		return nil, err
	}
	reader.WithoutNull = withoutNull
	reader.Encoding = fileInfo.Encoding
	reader.SingleLine = fileInfo.SingleLine

	var header []string
	if !fileInfo.NoHeader && !fileInfo.SingleLine {
		header, err = reader.ReadHeader()
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	records, err := readRecordSet(ctx, reader, fileSize(fp))
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

func loadViewFromCSVFile(ctx context.Context, fp io.ReadSeeker, fileInfo *FileInfo, withoutNull bool, expr parser.QueryExpression) (*View, error) {
	enc, err := text.DetectInSpecifiedEncoding(fp, fileInfo.Encoding)
	if err != nil {
		return nil, NewCannotDetectFileEncodingError(expr)
	}
	fileInfo.Encoding = enc

	reader, err := csv.NewReader(fp, fileInfo.Encoding)
	if err != nil {
		return nil, err
	}
	reader.Delimiter = fileInfo.Delimiter
	reader.WithoutNull = withoutNull

	var header []string
	if !fileInfo.NoHeader {
		header, err = reader.ReadHeader()
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	records, err := readRecordSet(ctx, reader, fileSize(fp))
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

func loadViewFromLTSVFile(ctx context.Context, flags *cmd.Flags, fp io.ReadSeeker, fileInfo *FileInfo, withoutNull bool, expr parser.QueryExpression) (*View, error) {
	enc, err := text.DetectInSpecifiedEncoding(fp, fileInfo.Encoding)
	if err != nil {
		return nil, NewCannotDetectFileEncodingError(expr)
	}
	fileInfo.Encoding = enc

	reader, err := ltsv.NewReader(fp, fileInfo.Encoding)
	if err != nil {
		return nil, NewIOError(expr, err.Error())
	}
	reader.WithoutNull = withoutNull

	records, err := readRecordSet(ctx, reader, fileSize(fp))
	if err != nil {
		return nil, err
	}

	header := reader.Header.Fields()
	if err = NewGoroutineTaskManager(len(records), -1, flags.CPU).Run(ctx, func(index int) error {
		for j := len(records[index]); j < len(header); j++ {
			if withoutNull {
				records[index] = append(records[index], NewCell(value.NewString("")))
			} else {
				records[index] = append(records[index], NewCell(value.NewNull()))
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if reader.DetectedLineBreak != "" {
		fileInfo.LineBreak = reader.DetectedLineBreak
	}

	view := NewView()
	view.Header = NewHeader(parser.FormatTableName(fileInfo.Path), header)
	view.RecordSet = records
	view.FileInfo = fileInfo
	return view, nil
}

func fileSize(fp io.ReadSeeker) int64 {
	if f, ok := fp.(*os.File); ok {
		if fi, err := f.Stat(); err == nil {
			return fi.Size()
		}
	}
	return 0
}

func readRecordSet(ctx context.Context, reader RecordReader, fileSize int64) (RecordSet, error) {
	var err error
	recordSet := make(RecordSet, 0, fileLoadingPreparedRecordSetCap)
	rowch := make(chan []text.RawText, fileLoadingBuffer)
	pos := 0

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for {
			row, ok := <-rowch
			if !ok {
				break
			}
			record := make(Record, len(row))
			for i, v := range row {
				if v == nil {
					record[i] = NewCell(value.NewNull())
				} else {
					record[i] = NewCell(value.NewString(string(v)))
				}
			}

			if 0 < fileSize && len(recordSet) == fileLoadingPreparedRecordSetCap && int64(pos) < fileSize {
				l := int((float64(fileSize) / float64(pos)) * fileLoadingPreparedRecordSetCap * 1.2)
				newSet := make(RecordSet, fileLoadingPreparedRecordSetCap, l)
				copy(newSet, recordSet)
				recordSet = newSet
			}

			recordSet = append(recordSet, record)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		i := 0
		for {
			if i&15 == 0 && ctx.Err() != nil {
				err = ConvertContextError(ctx.Err())
				break
			}

			row, e := reader.Read()
			if e == io.EOF {
				break
			}
			if e != nil {
				err = e
				break
			}

			if 0 < fileSize && i < fileLoadingPreparedRecordSetCap {
				for j := range row {
					pos += len(row[j])
				}
			}

			rowch <- row
			i++
		}
		close(rowch)
		wg.Done()
	}()

	wg.Wait()

	return recordSet, err
}

func loadViewFromJsonFile(fp io.Reader, fileInfo *FileInfo, expr parser.QueryExpression) (*View, error) {
	jsonText, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, NewIOError(expr, err.Error())
	}

	headerLabels, rows, escapeType, err := json.LoadTable(fileInfo.JsonQuery, string(jsonText))
	if err != nil {
		return nil, NewLoadJsonError(expr, err.Error())
	}

	records := make(RecordSet, len(rows))
	for i := range rows {
		records[i] = NewRecord(rows[i])
	}

	fileInfo.Encoding = text.UTF8
	fileInfo.JsonEscape = escapeType

	view := NewView()
	view.Header = NewHeader(parser.FormatTableName(fileInfo.Path), headerLabels)
	view.RecordSet = records
	view.FileInfo = fileInfo
	return view, nil
}

func loadDualView() *View {
	return &View{
		Header:    NewEmptyHeader(1),
		RecordSet: RecordSet{NewEmptyRecord(1)},
	}
}

func NewViewFromGroupedRecord(ctx context.Context, flags *cmd.Flags, referenceRecor ReferenceRecord) (*View, error) {
	view := NewView()
	view.Header = referenceRecor.view.Header
	record := referenceRecor.view.RecordSet[referenceRecor.recordIndex]

	view.RecordSet = make(RecordSet, record.GroupLen())

	if err := NewGoroutineTaskManager(record.GroupLen(), -1, flags.CPU).Run(ctx, func(index int) error {
		view.RecordSet[index] = make(Record, view.FieldLen())
		for j := range record {
			grpIdx := index
			if len(record[j]) < 2 {
				grpIdx = 0
			}
			view.RecordSet[index][j] = record[j][grpIdx : grpIdx+1]
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return view, nil
}

func (view *View) Where(ctx context.Context, scope *ReferenceScope, clause parser.WhereClause) error {
	return view.filter(ctx, scope, clause.Filter)
}

func (view *View) filter(ctx context.Context, scope *ReferenceScope, condition parser.QueryExpression) error {
	results := make([]bool, view.RecordLen())

	if err := EvaluateSequentially(ctx, scope, view, func(seqScope *ReferenceScope, rIdx int) error {
		primary, e := Evaluate(ctx, seqScope, condition)
		if e != nil {
			return e
		}

		if primary.Ternary() == ternary.TRUE {
			results[rIdx] = true
		}
		return nil
	}); err != nil {
		return err
	}

	newIdx := 0
	for i, ok := range results {
		if ok {
			if i != newIdx {
				view.RecordSet[newIdx] = view.RecordSet[i]
			}
			newIdx++
		}
	}

	view.RecordSet = view.RecordSet[:newIdx]
	return nil
}

func (view *View) GroupBy(ctx context.Context, scope *ReferenceScope, clause parser.GroupByClause) error {
	return view.group(ctx, scope, clause.Items)
}

func (view *View) group(ctx context.Context, scope *ReferenceScope, items []parser.QueryExpression) error {
	if items == nil {
		return view.groupAll(ctx, scope.Tx.Flags)
	}

	gm := NewGoroutineTaskManager(view.RecordLen(), -1, scope.Tx.Flags.CPU)
	groupsList := make([]map[string][]int, gm.Number)
	groupKeyCnt := make(map[string]int, 20)
	groupKeys := make([]string, 0, 20)
	mtx := &sync.Mutex{}

	var grpFn = func(thIdx int) {
		start, end := gm.RecordRange(thIdx)
		seqScope := scope.CreateScopeForSequentialEvaluation(view)
		groups := make(map[string][]int, 20)

	GroupKeyLoop:
		for i := start; i < end; i++ {
			if gm.HasError() {
				break GroupKeyLoop
			}
			if i&15 == 0 && ctx.Err() != nil {
				break GroupKeyLoop
			}

			seqScope.Records[0].recordIndex = i

			values := make([]value.Primary, len(items))
			for i, item := range items {
				p, e := Evaluate(ctx, seqScope, item)
				if e != nil {
					gm.SetError(e)
					break GroupKeyLoop
				}
				values[i] = p
			}
			keyBuf := GetComparisonKeysBuf()
			SerializeComparisonKeys(keyBuf, values, seqScope.Tx.Flags)
			key := keyBuf.String()
			PutComparisonkeysBuf(keyBuf)

			if _, ok := groups[key]; ok {
				groups[key] = append(groups[key], i)
			} else {
				groups[key] = make([]int, 0, view.RecordLen()/18)
				groups[key] = append(groups[key], i)
				mtx.Lock()
				if _, ok := groupKeyCnt[key]; !ok {
					groupKeyCnt[key] = 0
					groupKeys = append(groupKeys, key)
				}
				mtx.Unlock()
			}
		}

		groupsList[thIdx] = groups

		if 1 < gm.Number {
			gm.Done()
		}
	}

	if 1 < gm.Number {
		for i := 0; i < gm.Number; i++ {
			gm.Add()
			go grpFn(i)
		}
		gm.Wait()
	} else {
		grpFn(0)
	}

	if gm.HasError() {
		return gm.Err()
	}
	if ctx.Err() != nil {
		return ConvertContextError(ctx.Err())
	}

	for i := range groupsList {
		for k := range groupsList[i] {
			groupKeyCnt[k] = groupKeyCnt[k] + len(groupsList[i][k])
		}
	}

	records := make(RecordSet, len(groupKeys))
	calcCnt := view.RecordLen() * len(groupKeys)
	minReq := -1
	if MinimumRequiredPerCPUCore < calcCnt {
		minReq = int(math.Ceil(float64(len(groupKeys)) / (math.Floor(float64(calcCnt) / MinimumRequiredPerCPUCore))))
	}
	if err := NewGoroutineTaskManager(len(groupKeys), minReq, scope.Tx.Flags.CPU).Run(ctx, func(gIdx int) error {
		record := make(Record, view.FieldLen())

		for i := 0; i < view.FieldLen(); i++ {
			primaries := make(Cell, groupKeyCnt[groupKeys[gIdx]])
			pos := 0
			for j := range groupsList {
				if indices, ok := groupsList[j][groupKeys[gIdx]]; ok {
					for k := range indices {
						primaries[pos+k] = view.RecordSet[indices[k]][i][0]
					}
					pos += len(indices)
				}
			}
			record[i] = primaries
		}

		records[gIdx] = record
		return nil
	}); err != nil {
		return err
	}

	view.RecordSet = records
	view.isGrouped = true
	for _, item := range items {
		switch item.(type) {
		case parser.FieldReference, parser.ColumnNumber:
			idx, _ := view.Header.SearchIndex(item)
			view.Header[idx].IsGroupKey = true
		}
	}
	return nil
}

func (view *View) groupAll(ctx context.Context, flags *cmd.Flags) error {
	if 0 < view.RecordLen() {
		record := make(Record, view.FieldLen())

		calcCnt := view.RecordLen() * view.FieldLen()
		minReq := -1
		if MinimumRequiredPerCPUCore < calcCnt {
			minReq = int(math.Ceil(float64(view.FieldLen()) / (math.Floor(float64(calcCnt) / MinimumRequiredPerCPUCore))))
		}
		if err := NewGoroutineTaskManager(view.FieldLen(), minReq, flags.CPU).Run(ctx, func(fIdx int) error {
			primaries := make(Cell, len(view.RecordSet))
			for j := range view.RecordSet {
				primaries[j] = view.RecordSet[j][fIdx][0]
			}
			record[fIdx] = primaries
			return nil
		}); err != nil {
			return err
		}

		view.RecordSet = view.RecordSet[:1]
		view.RecordSet[0] = record
	}

	view.isGrouped = true
	return nil
}

func (view *View) Having(ctx context.Context, scope *ReferenceScope, clause parser.HavingClause) error {
	err := view.filter(ctx, scope, clause.Filter)
	if err != nil {
		if _, ok := err.(*NotGroupingRecordsError); ok {
			if err = view.group(ctx, scope, nil); err != nil {
				return err
			}
			if err = view.filter(ctx, scope, clause.Filter); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (view *View) Select(ctx context.Context, scope *ReferenceScope, clause parser.SelectClause) error {
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
		if err := view.ExtendRecordCapacity(ctx, scope, fieldsObjects); err != nil {
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
			idx, err := view.evalColumn(ctx, scope, field.Object, alias)
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

			if err = view.group(ctx, scope, nil); err != nil {
				return err
			}

			if err = evalFields(view, fields); err != nil {
				return err
			}

			if view.tempRecord != nil {
				view.RecordSet = append(view.RecordSet, view.tempRecord)
				view.tempRecord = nil
			}
		} else {
			return err
		}
	}

	if clause.IsDistinct() {
		if err = view.GenerateComparisonKeys(ctx, scope.Tx.Flags); err != nil {
			return err
		}
		records := make(RecordSet, 0, 40)
		values := make(map[string]bool, 40)
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

func (view *View) GenerateComparisonKeys(ctx context.Context, flags *cmd.Flags) error {
	view.comparisonKeysInEachRecord = make([]string, view.RecordLen())

	return NewGoroutineTaskManager(view.RecordLen(), -1, flags.CPU).Run(ctx, func(index int) error {
		flags := flags
		buf := GetComparisonKeysBuf()
		var primaries []value.Primary = nil

		if view.selectFields != nil {
			primaries = make([]value.Primary, len(view.selectFields))
			for j, idx := range view.selectFields {
				primaries[j] = view.RecordSet[index][idx][0]
			}
		} else {
			primaries = make([]value.Primary, view.FieldLen())
			for j := range view.RecordSet[index] {
				primaries[j] = view.RecordSet[index][j][0]
			}
		}
		SerializeComparisonKeys(buf, primaries, flags)

		view.comparisonKeysInEachRecord[index] = buf.String()
		PutComparisonkeysBuf(buf)
		return nil
	})
}

func (view *View) SelectAllColumns(ctx context.Context, scope *ReferenceScope) error {
	selectClause := parser.SelectClause{
		Fields: []parser.QueryExpression{
			parser.Field{Object: parser.AllColumns{}},
		},
	}
	return view.Select(ctx, scope, selectClause)
}

func (view *View) OrderBy(ctx context.Context, scope *ReferenceScope, clause parser.OrderByClause) error {
	orderValues := make([]parser.QueryExpression, len(clause.Items))
	for i, item := range clause.Items {
		orderValues[i] = item.(parser.OrderItem).Value
	}
	if err := view.ExtendRecordCapacity(ctx, scope, orderValues); err != nil {
		return err
	}

	sortIndices := make([]int, len(clause.Items))
	for i, v := range clause.Items {
		oi := v.(parser.OrderItem)
		idx, err := view.evalColumn(ctx, scope, oi.Value, "")
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

		if oi.NullsPosition.IsEmpty() {
			switch view.sortDirections[i] {
			case parser.ASC:
				view.sortNullPositions[i] = parser.FIRST
			default: //parser.DESC
				view.sortNullPositions[i] = parser.LAST
			}
		} else {
			view.sortNullPositions[i] = oi.NullsPosition.Token
		}
	}

	if err := NewGoroutineTaskManager(view.RecordLen(), -1, scope.Tx.Flags.CPU).Run(ctx, func(index int) error {
		if view.sortValuesInEachCell != nil && view.sortValuesInEachCell[index] == nil {
			view.sortValuesInEachCell[index] = make([]*SortValue, cap(view.RecordSet[index]))
		}

		sortValues := make(SortValues, len(sortIndices))
		for j, idx := range sortIndices {
			if view.sortValuesInEachCell != nil && idx < len(view.sortValuesInEachCell[index]) && view.sortValuesInEachCell[index][idx] != nil {
				sortValues[j] = view.sortValuesInEachCell[index][idx]
			} else {
				sortValues[j] = NewSortValue(view.RecordSet[index][idx][0], scope.Tx.Flags)
				if view.sortValuesInEachCell != nil && idx < len(view.sortValuesInEachCell[index]) {
					view.sortValuesInEachCell[index][idx] = sortValues[j]
				}
			}
		}
		view.sortValuesInEachRecord[index] = sortValues
		return nil
	}); err != nil {
		return err
	}

	sort.Sort(view)
	return nil
}

func (view *View) additionalColumns(ctx context.Context, scope *ReferenceScope, expr parser.QueryExpression) (map[string]bool, error) {
	m := make(map[string]bool, 5)

	switch expr.(type) {
	case parser.FieldReference, parser.ColumnNumber:
		return nil, nil
	case parser.Function:
		if udfn, err := scope.GetFunction(expr, expr.(parser.Function).Name); err == nil {
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
				columns, err := view.additionalColumns(ctx, scope, pvalue)
				if err != nil {
					return nil, err
				}
				for k := range columns {
					if _, ok := m[k]; !ok {
						m[k] = true
					}
				}
			}
		}
		if ovalues != nil {
			for _, v := range ovalues {
				item := v.(parser.OrderItem)
				columns, err := view.additionalColumns(ctx, scope, item.Value)
				if err != nil {
					return nil, err
				}
				for k := range columns {
					if _, ok := m[k]; !ok {
						m[k] = true
					}
				}
			}
		}
	}

	if _, ok := view.Header.ContainsObject(expr); !ok {
		s := expr.String()
		if _, ok := m[s]; !ok {
			m[s] = true
		}
	}

	return m, nil
}

func (view *View) ExtendRecordCapacity(ctx context.Context, scope *ReferenceScope, exprs []parser.QueryExpression) error {
	additions := make(map[string]bool, 5)
	for _, expr := range exprs {
		columns, err := view.additionalColumns(ctx, scope, expr)
		if err != nil {
			return err
		}
		for k := range columns {
			if _, ok := additions[k]; !ok {
				additions[k] = true
			}
		}
	}

	currentLen := view.FieldLen()
	fieldCap := currentLen + len(additions)

	if 0 < view.RecordLen() && fieldCap <= cap(view.RecordSet[0]) {
		return nil
	}

	return NewGoroutineTaskManager(view.RecordLen(), -1, scope.Tx.Flags.CPU).Run(ctx, func(index int) error {
		record := make(Record, currentLen, fieldCap)
		copy(record, view.RecordSet[index])
		view.RecordSet[index] = record
		return nil
	})
}

func (view *View) evalColumn(ctx context.Context, scope *ReferenceScope, obj parser.QueryExpression, alias string) (idx int, err error) {
	idx, ok := view.Header.ContainsObject(obj)
	if ok {
		rScope := scope.CreateScopeForRecordEvaluation(view, -1)
		if _, err = Evaluate(ctx, rScope, obj); err != nil {
			return
		}
	} else {
		if analyticFunction, ok := obj.(parser.AnalyticFunction); ok {
			err = view.evalAnalyticFunction(ctx, scope, analyticFunction)
			if err != nil {
				return
			}
		} else if view.RecordLen() < 1 {
			if view.tempRecord == nil {
				view.tempRecord = NewEmptyRecord(view.FieldLen())
			}

			rScope := scope.CreateScopeForRecordEvaluation(view, -1)
			primary, e := Evaluate(ctx, rScope, obj)
			if e != nil {
				err = e
				return
			}
			view.tempRecord = append(view.tempRecord, NewCell(primary))
		} else {
			if err = EvaluateSequentially(ctx, scope, view, func(seqScope *ReferenceScope, rIdx int) error {
				primary, e := Evaluate(ctx, seqScope, obj)
				if e != nil {
					return e
				}

				view.RecordSet[rIdx] = append(view.RecordSet[rIdx], NewCell(primary))
				return nil
			}); err != nil {
				return
			}
		}
		view.Header, idx = AddHeaderField(view.Header, parser.FormatFieldIdentifier(obj), alias)
	}

	if 0 < len(alias) {
		if !strings.EqualFold(view.Header[idx].Column, alias) && !InStrSliceWithCaseInsensitive(alias, view.Header[idx].Aliases) {
			view.Header[idx].Aliases = append(view.Header[idx].Aliases, alias)
		}
	}

	return
}

func (view *View) evalAnalyticFunction(ctx context.Context, scope *ReferenceScope, expr parser.AnalyticFunction) error {
	name := strings.ToUpper(expr.Name)
	if _, ok := AggregateFunctions[name]; !ok {
		if _, ok := AnalyticFunctions[name]; !ok {
			if udfn, err := scope.GetFunction(expr, expr.Name); err != nil || !udfn.IsAggregate {
				return NewFunctionNotExistError(expr, expr.Name)
			}
		}
	}

	var partitionIndices []int
	if expr.AnalyticClause.PartitionClause != nil {
		partitionExprs := expr.AnalyticClause.PartitionValues()

		partitionIndices = make([]int, len(partitionExprs))
		for i, pexpr := range partitionExprs {
			idx, err := view.evalColumn(ctx, scope, pexpr, "")
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
		err := view.OrderBy(ctx, scope, expr.AnalyticClause.OrderByClause.(parser.OrderByClause))
		if err != nil {
			return err
		}
	}

	err := Analyze(ctx, scope, view, expr, partitionIndices)

	view.sortValuesInEachRecord = nil
	view.sortDirections = nil
	view.sortNullPositions = nil

	return err
}

func (view *View) Offset(ctx context.Context, scope *ReferenceScope, clause parser.OffsetClause) error {
	val, err := Evaluate(ctx, scope, clause.Value)
	if err != nil {
		return err
	}
	number := value.ToInteger(val)
	if value.IsNull(number) {
		return NewInvalidOffsetNumberError(clause)
	}
	view.offset = int(number.(*value.Integer).Raw())
	value.Discard(number)

	if view.offset < 0 {
		view.offset = 0
	}

	if view.RecordLen() <= view.offset {
		view.RecordSet = RecordSet{}
	} else {
		newSet := view.RecordSet[view.offset:]
		view.RecordSet = view.RecordSet[:len(newSet)]
		for i := range newSet {
			view.RecordSet[i] = newSet[i]
		}
	}
	return nil
}

func (view *View) Limit(ctx context.Context, scope *ReferenceScope, clause parser.LimitClause) error {
	val, err := Evaluate(ctx, scope, clause.Value)
	if err != nil {
		return err
	}

	var limit int
	if clause.Percentage() {
		number := value.ToFloat(val)
		if value.IsNull(number) {
			return NewInvalidLimitPercentageError(clause)
		}
		percentage := number.(*value.Float).Raw()
		value.Discard(number)

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
		limit = int(number.(*value.Integer).Raw())
		value.Discard(number)

		if limit < 0 {
			limit = 0
		}
	}

	if view.RecordLen() <= limit {
		return nil
	}

	if clause.WithTies() && view.sortValuesInEachRecord != nil {
		bottomSortValues := view.sortValuesInEachRecord[limit-1]
		for limit < view.RecordLen() {
			if !bottomSortValues.EquivalentTo(view.sortValuesInEachRecord[limit]) {
				break
			}
			limit++
		}
	}

	view.RecordSet = view.RecordSet[:limit]
	return nil
}

func (view *View) InsertValues(ctx context.Context, scope *ReferenceScope, fields []parser.QueryExpression, list []parser.QueryExpression) (int, error) {
	recordValues, err := view.convertListToRecordValues(ctx, scope, fields, list)
	if err != nil {
		return 0, err
	}
	return view.insert(ctx, fields, recordValues)
}

func (view *View) InsertFromQuery(ctx context.Context, scope *ReferenceScope, fields []parser.QueryExpression, query parser.SelectQuery) (int, error) {
	recordValues, err := view.convertResultSetToRecordValues(ctx, scope, fields, query)
	if err != nil {
		return 0, err
	}
	return view.insert(ctx, fields, recordValues)
}

func (view *View) ReplaceValues(ctx context.Context, scope *ReferenceScope, fields []parser.QueryExpression, list []parser.QueryExpression, keys []parser.QueryExpression) (int, error) {
	recordValues, err := view.convertListToRecordValues(ctx, scope, fields, list)
	if err != nil {
		return 0, err
	}
	return view.replace(ctx, scope.Tx.Flags, fields, recordValues, keys)
}

func (view *View) ReplaceFromQuery(ctx context.Context, scope *ReferenceScope, fields []parser.QueryExpression, query parser.SelectQuery, keys []parser.QueryExpression) (int, error) {
	recordValues, err := view.convertResultSetToRecordValues(ctx, scope, fields, query)
	if err != nil {
		return 0, err
	}
	return view.replace(ctx, scope.Tx.Flags, fields, recordValues, keys)
}

func (view *View) convertListToRecordValues(ctx context.Context, scope *ReferenceScope, fields []parser.QueryExpression, list []parser.QueryExpression) ([][]value.Primary, error) {
	recordValues := make([][]value.Primary, len(list))
	for i, item := range list {
		if ctx.Err() != nil {
			return nil, ConvertContextError(ctx.Err())
		}

		rv := item.(parser.RowValue)
		values, err := EvalRowValue(ctx, scope, rv)
		if err != nil {
			return recordValues, err
		}
		if len(fields) != len(values) {
			return recordValues, NewInsertRowValueLengthError(rv, len(fields))
		}

		recordValues[i] = values
	}
	return recordValues, nil
}

func (view *View) convertResultSetToRecordValues(ctx context.Context, scope *ReferenceScope, fields []parser.QueryExpression, query parser.SelectQuery) ([][]value.Primary, error) {
	selectedView, err := Select(ctx, scope, query)
	if err != nil {
		return nil, err
	}
	if len(fields) != selectedView.FieldLen() {
		return nil, NewInsertSelectFieldLengthError(query, len(fields))
	}

	recordValues := make([][]value.Primary, selectedView.RecordLen())
	for i, record := range selectedView.RecordSet {
		if ctx.Err() != nil {
			return nil, ConvertContextError(ctx.Err())
		}

		values := make([]value.Primary, selectedView.FieldLen())
		for j, cell := range record {
			values[j] = cell[0]
		}
		recordValues[i] = values
	}
	return recordValues, nil
}

func (view *View) convertRecordValuesToRecordSet(ctx context.Context, fields []parser.QueryExpression, recordValues [][]value.Primary) (RecordSet, error) {
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
		return nil, err
	}

	recordIndices := make([]int, view.FieldLen())
	for i := 0; i < view.FieldLen(); i++ {
		recordIndices[i] = valueIndex(i, fieldIndices)
	}

	records := make(RecordSet, len(recordValues))
	for i, values := range recordValues {
		if ctx.Err() != nil {
			return nil, ConvertContextError(ctx.Err())
		}

		record := make(Record, view.FieldLen())
		for j := 0; j < view.FieldLen(); j++ {
			if recordIndices[j] < 0 {
				record[j] = NewCell(value.NewNull())
			} else {
				record[j] = NewCell(values[recordIndices[j]])
			}
		}
		records[i] = record
	}
	return records, nil
}

func (view *View) insert(ctx context.Context, fields []parser.QueryExpression, recordValues [][]value.Primary) (int, error) {
	records, err := view.convertRecordValuesToRecordSet(ctx, fields, recordValues)
	if err != nil {
		return 0, err
	}

	view.RecordSet = view.RecordSet.Merge(records)
	return len(recordValues), nil
}

func (view *View) replace(ctx context.Context, flags *cmd.Flags, fields []parser.QueryExpression, recordValues [][]value.Primary, keys []parser.QueryExpression) (int, error) {
	fieldIndices, err := view.FieldIndices(fields)
	if err != nil {
		return 0, err
	}
	fieldIndicesMap := make(map[uint]bool, len(fieldIndices))
	for _, v := range fieldIndices {
		fieldIndicesMap[uint(v)] = true
	}

	keyIndices, err := view.FieldIndices(keys)
	if err != nil {
		return 0, err
	}
	keyIndicesMap := make(map[uint]bool, len(keyIndices))
	for _, v := range keyIndices {
		keyIndicesMap[uint(v)] = true
	}

	for idx, i := range keyIndices {
		if _, ok := fieldIndicesMap[uint(i)]; !ok {
			return 0, NewReplaceKeyNotSetError(keys[idx])
		}
	}
	updateIndices := make([]int, 0, len(fieldIndices)-len(keyIndices))
	for _, i := range fieldIndices {
		if _, ok := keyIndicesMap[uint(i)]; !ok {
			updateIndices = append(updateIndices, i)
		}
	}

	records, err := view.convertRecordValuesToRecordSet(ctx, fields, recordValues)
	if err != nil {
		return 0, err
	}

	sortValuesInEachRecord := make([]SortValues, view.RecordLen())
	if err := NewGoroutineTaskManager(view.RecordLen(), -1, flags.CPU).Run(ctx, func(index int) error {
		sortValues := make(SortValues, len(keyIndices))
		for j, idx := range keyIndices {
			sortValues[j] = NewSortValue(view.RecordSet[index][idx][0], flags)
		}
		sortValuesInEachRecord[index] = sortValues
		return nil
	}); err != nil {
		return 0, err
	}

	sortValuesInInsertRecords := make([]SortValues, view.RecordLen())
	if err := NewGoroutineTaskManager(len(records), -1, flags.CPU).Run(ctx, func(index int) error {
		sortValues := make(SortValues, len(keyIndices))
		for j, idx := range keyIndices {
			sortValues[j] = NewSortValue(records[index][idx][0], flags)
		}
		sortValuesInInsertRecords[index] = sortValues
		return nil
	}); err != nil {
		return 0, err
	}

	replacedRecord := make(map[int]bool, len(records))
	replacedCount := 0
	for i := range records {
		replacedRecord[i] = false
	}
	replaceMtx := &sync.Mutex{}
	var replaced = func(idx int) {
		replaceMtx.Lock()
		replacedRecord[idx] = true
		replacedCount++
		replaceMtx.Unlock()
	}
	if err := NewGoroutineTaskManager(view.RecordLen(), -1, flags.CPU).Run(ctx, func(index int) error {
		for j, rsv := range sortValuesInInsertRecords {
			if sortValuesInEachRecord[index].EquivalentTo(rsv) {
				for _, fidx := range updateIndices {
					view.RecordSet[index][fidx] = records[j][fidx]
				}
				replaced(j)
				break
			}
		}
		return nil
	}); err != nil {
		return 0, err
	}

	insertRecords := make(RecordSet, 0, len(records))
	for i, isReplaced := range replacedRecord {
		if !isReplaced {
			insertRecords = append(insertRecords, records[i])
		}
	}
	view.RecordSet = view.RecordSet.Merge(insertRecords)
	return len(insertRecords) + replacedCount, nil
}

func (view *View) Fix(ctx context.Context, flags *cmd.Flags) error {
	fieldLen := len(view.selectFields)
	resize := false
	if fieldLen != view.FieldLen() {
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
		if err := NewGoroutineTaskManager(view.RecordLen(), -1, flags.CPU).Run(ctx, func(index int) error {
			record := make(Record, fieldLen)
			for j, idx := range view.selectFields {
				record[j] = view.RecordSet[index][idx][:1]
			}

			if len(view.RecordSet[index]) < fieldLen {
				view.RecordSet[index] = make(Record, fieldLen)
			} else if fieldLen < len(view.RecordSet[index]) {
				view.RecordSet[index] = view.RecordSet[index][:fieldLen]
			}
			for i := range record {
				view.RecordSet[index][i] = record[i]
			}
			return nil
		}); err != nil {
			return err
		}
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
	view.selectFields = nil
	view.selectLabels = nil
	view.isGrouped = false
	view.comparisonKeysInEachRecord = nil
	view.sortValuesInEachCell = nil
	view.sortValuesInEachRecord = nil
	view.sortDirections = nil
	view.sortNullPositions = nil
	view.offset = 0
	view.tempRecord = nil
	return nil
}

func (view *View) Union(ctx context.Context, flags *cmd.Flags, calcView *View, all bool) (err error) {
	view.RecordSet = view.RecordSet.Merge(calcView.RecordSet)
	view.FileInfo = nil

	if !all {
		if err = view.GenerateComparisonKeys(ctx, flags); err != nil {
			return err
		}

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
	return
}

func (view *View) Except(ctx context.Context, flags *cmd.Flags, calcView *View, all bool) (err error) {
	if err = view.GenerateComparisonKeys(ctx, flags); err != nil {
		return err
	}
	if err = calcView.GenerateComparisonKeys(ctx, flags); err != nil {
		return err
	}

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
	return
}

func (view *View) Intersect(ctx context.Context, flags *cmd.Flags, calcView *View, all bool) (err error) {
	if err = view.GenerateComparisonKeys(ctx, flags); err != nil {
		return err
	}
	if err = calcView.GenerateComparisonKeys(ctx, flags); err != nil {
		return err
	}

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
	return
}

func (view *View) ListValuesForAggregateFunctions(ctx context.Context, scope *ReferenceScope, expr parser.QueryExpression, arg parser.QueryExpression, distinct bool) ([]value.Primary, error) {
	list := make([]value.Primary, view.RecordLen())

	if err := EvaluateSequentially(ctx, scope, view, func(sqlScope *ReferenceScope, rIdx int) error {
		p, e := Evaluate(ctx, sqlScope, arg)
		if e != nil {
			if _, ok := e.(*NotGroupingRecordsError); ok {
				e = NewNestedAggregateFunctionsError(expr)
			}
			return e
		}
		list[rIdx] = p
		return nil
	}); err != nil {
		return nil, err
	}

	if distinct {
		list = Distinguish(list, scope.Tx.Flags)
	}

	return list, nil
}

func (view *View) RestoreHeaderReferences() error {
	return view.Header.Update(parser.FormatTableName(view.FileInfo.Path), nil)
}

func (view *View) FieldIndex(fieldRef parser.QueryExpression) (int, error) {
	idx, err := view.Header.SearchIndex(fieldRef)
	if err != nil {
		if err == errFieldAmbiguous {
			err = NewFieldAmbiguousError(fieldRef)
		} else {
			err = NewFieldNotExistError(fieldRef)
		}
	}
	return idx, err
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
	internalId, ok := view.RecordSet[recordIndex][idx][0].(*value.Integer)
	if !ok {
		return -1, NewInternalRecordIdEmptyError()
	}
	return int(internalId.Raw()), nil
}

func (view *View) CreateRestorePoint() {
	view.FileInfo.restorePointRecordSet = view.RecordSet.Copy()
	view.FileInfo.restorePointHeader = view.Header.Copy()
}

func (view *View) Restore() {
	view.RecordSet = view.FileInfo.restorePointRecordSet.Copy()
	view.Header = view.FileInfo.restorePointHeader.Copy()
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
	}
}
