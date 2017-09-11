package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

const (
	ERROR_MESSAGE_TEMPLATE                     = "[L:%d C:%d] %s"
	ERROR_MESSAGE_WITH_FILEPATH_TEMPLATE       = "%s [L:%d C:%d] %s"
	ERROR_MESSAGE_WITH_EMPTY_POSITION_TEMPLATE = "[L:- C:-] %s"

	ERROR_INVALID_SYNTAX                    = "syntax error: unexpected %s"
	ERROR_READ_FILE                         = "failed to read from file: %s"
	ERROR_CREATE_FILE                       = "failed to create file: %s"
	ERROR_WRITE_FILE                        = "failed to write to file: %s"
	ERROR_WRITE_FILE_IN_AUTOCOMMIT          = "[Auto-Commit] failed to write to file: %s"
	ERROR_FIELD_AMBIGUOUS                   = "field %s is ambiguous"
	ERROR_FIELD_NOT_EXIST                   = "field %s does not exist"
	ERROR_FIELD_NOT_GROUP_KEY               = "field %s is not a group key"
	ERROR_DUPLICATE_FIELD_NAME              = "field name %s is a duplicate"
	ERROR_NOT_GROUPING_RECORDS              = "function %s cannot aggregate not grouping records"
	ERROR_UNDEFINED_VARIABLE                = "variable %s is undefined"
	ERROR_VARIABLE_REDECLARED               = "variable %s is redeclared"
	ERROR_FUNCTION_NOT_EXIST                = "function %s does not exist"
	ERROR_FUNCTION_ARGUMENT_LENGTH          = "function %s takes %s"
	ERROR_FUNCTION_INVALID_ARGUMENT         = "%s for function %s"
	ERROR_UNPERMITTED_STATEMENT_FUNCTION    = "function %s cannot be used as a statement"
	ERROR_NESTED_AGGREGATE_FUNCTIONS        = "aggregate functions are nested at %s"
	ERROR_FUNCTION_REDECLARED               = "function %s is redeclared"
	ERROR_BUILT_IN_FUNCTION_DECLARED        = "function %s is a built-in function"
	ERROR_DUPLICATE_PARAMETER               = "parameter %s is a duplicate"
	ERROR_SUBQUERY_TOO_MANY_RECORDS         = "subquery returns too many records, should return only one record"
	ERROR_SUBQUERY_TOO_MANY_FIELDS          = "subquery returns too many fields, should return only one field"
	ERROR_CURSOR_REDECLARED                 = "cursor %s is redeclared"
	ERROR_UNDEFINED_CURSOR                  = "cursor %s is undefined"
	ERROR_CURSOR_CLOSED                     = "cursor %s is closed"
	ERROR_CURSOR_OPEN                       = "cursor %s is already open"
	ERROR_PSEUDO_CURSOR                     = "cursor %s is a pseudo cursor"
	ERROR_CURSOR_FETCH_LENGTH               = "fetching from cursor %s returns %s"
	ERROR_INVALID_FETCH_POSITION            = "fetching position %s is not an integer value"
	ERROR_INLINE_TABLE_REDECLARED           = "inline table %s is redeclared"
	ERROR_UNDEFINED_INLINE_TABLE            = "inline table %s is undefined"
	ERROR_INLINE_TABLE_FIELD_LENGTH         = "select query should return exactly %s for inline table %s"
	ERROR_FILE_NOT_EXIST                    = "file %s does not exist"
	ERROR_FILE_ALREADY_EXIST                = "file %s already exists"
	ERROR_FILE_UNABLE_TO_READ               = "file %s is unable to be read"
	ERROR_FILE_LOCK_TIMEOUT                 = "file %s: lock wait timeout period exceeded"
	ERROR_CSV_PARSING                       = "csv parse error in file %s: %s"
	ERROR_TABLE_FIELD_LENGTH                = "select query should return exactly %s for table %s"
	ERROR_TEMPORARY_TABLE_REDECLARED        = "temporary table %s is redeclared"
	ERROR_UNDEFINED_TEMPORARY_TABLE         = "temporary table %s is undefined"
	ERROR_TEMPORARY_TABLE_FIELD_LENGTH      = "select query should return exactly %s for temporary table %s"
	ERROR_DUPLICATE_TABLE_NAME              = "table name %s is a duplicate"
	ERROR_TABLE_NOT_LOADED                  = "table %s is not loaded"
	ERROR_STDIN_EMPTY                       = "stdin is empty"
	ERROR_ROW_VALUE_LENGTH_IN_COMPARISON    = "row value should contain exactly %s"
	ERROR_SELECT_FIELD_LENGTH_IN_COMPARISON = "select query should return exactly %s"
	ERROR_INVALID_LIMIT_PERCENTAGE          = "limit percentage %s is not a float value"
	ERROR_INVALID_LIMIT_NUMBER              = "limit number of records %s is not an integer value"
	ERROR_INVALID_OFFSET_NUMBER             = "offset number %s is not an integer value"
	ERROR_COMBINED_SET_FIELD_LENGTH         = "result set to be combined should contain exactly %s"
	ERROR_INSERT_ROW_VALUE_LENGTH           = "row value should contain exactly %s"
	ERROR_INSERT_SELECT_FIELD_LENGTH        = "select query should return exactly %s"
	ERROR_UPDATE_FIELD_NOT_EXIST            = "field %s does not exist in the tables to update"
	ERROR_UPDATE_VALUE_AMBIGUOUS            = "value %s to set in the field %s is ambiguous"
	ERROR_DELETE_TABLE_NOT_SPECIFIED        = "tables to delete records are not specified"
	ERROR_PRINTF_REPLACE_VALUE_LENGTH       = "PRINTF: %s"
	ERROR_SOURCE_INVALID_ARGUMENT           = "SOURCE: argument %s is not a string"
	ERROR_SOURCE_FILE_NOT_EXIST             = "SOURCE: file %s does not exist"
	ERROR_SOURCE_FILE_UNABLE_TO_READ        = "SOURCE: file %s is unable to read"
	ERROR_INVALID_FLAG_NAME                 = "flag name %s is invalid"
	ERROR_INVALID_FLAG_VALUE                = "SET: flag value %s for %s is invalid"
	ERROR_INTERNAL_RECORD_ID_NOT_EXIST      = "internal record id does not exist"
	ERROR_INTERNAL_RECORD_ID_EMPTY          = "internal record id is empty"
	ERROR_FIELD_LENGTH_NOT_MATCH            = "field length does not match"
	ERROR_ROW_VALUE_LENGTH_NOT_MATCH        = "row value length does not match"
	ERROR_ROW_VALUE_LENGTH_IN_LIST          = "row value length does not match at index %d"
	ERROR_FORMAT_STRING_LENGTH_NOT_MATCH    = "number of replace values does not match"
)

type Exit struct {
	Code int
}

func NewExit(code int) error {
	return &Exit{
		Code: code,
	}
}

func (e Exit) Error() string {
	return ""
}

func (e Exit) GetCode() int {
	return e.Code
}

type AppError interface {
	Error() string
	ErrorMessage() string
	GetCode() int
}

type BaseError struct {
	SourceFile string
	Line       int
	Char       int
	Message    string
	Code       int
}

func (e BaseError) Error() string {
	if e.Line < 1 {
		return fmt.Sprintf(ERROR_MESSAGE_WITH_EMPTY_POSITION_TEMPLATE, e.Message)
	}
	if 0 < len(e.SourceFile) {
		return fmt.Sprintf(ERROR_MESSAGE_WITH_FILEPATH_TEMPLATE, e.SourceFile, e.Line, e.Char, e.Message)
	}
	return fmt.Sprintf(ERROR_MESSAGE_TEMPLATE, e.Line, e.Char, e.Message)
}

func (e BaseError) ErrorMessage() string {
	return e.Message
}

func (e BaseError) GetCode() int {
	return e.Code
}

func NewBaseError(expr parser.Expression, message string) *BaseError {
	return NewBaseErrorWithCode(expr, message, 1)
}

func NewBaseErrorWithCode(expr parser.Expression, message string, code int) *BaseError {
	var sourceFile string
	var line int
	var char int
	if expr.HasParseInfo() {
		sourceFile = expr.SourceFile()
		line = expr.Line()
		char = expr.Char()
	}

	return &BaseError{
		SourceFile: sourceFile,
		Line:       line,
		Char:       char,
		Message:    message,
		Code:       code,
	}
}

type UserTriggeredError struct {
	*BaseError
}

func NewUserTriggeredError(expr parser.Trigger, message string) error {
	code := 1
	if expr.Code != nil {
		code = int(expr.Code.(value.Integer).Raw())
	}

	return &UserTriggeredError{
		NewBaseErrorWithCode(expr, message, code),
	}
}

type SyntaxError struct {
	*BaseError
}

func NewSyntaxError(message string, line int, char int, sourceFile string) error {
	return &SyntaxError{
		&BaseError{
			SourceFile: sourceFile,
			Line:       line,
			Char:       char,
			Message:    message,
			Code:       1,
		},
	}
}

func NewSyntaxErrorFromExpr(expr parser.QueryExpression) error {
	return &SyntaxError{
		NewBaseError(expr, fmt.Sprintf(ERROR_INVALID_SYNTAX, expr)),
	}
}

type ReadFileError struct {
	*BaseError
}

func NewReadFileError(expr parser.Expression, message string) error {
	return &ReadFileError{
		NewBaseError(expr, fmt.Sprintf(ERROR_READ_FILE, message)),
	}
}

type CreateFileError struct {
	*BaseError
}

func NewCreateFileError(expr parser.Expression, message string) error {
	return &CreateFileError{
		NewBaseError(expr, fmt.Sprintf(ERROR_CREATE_FILE, message)),
	}
}

type WriteFileError struct {
	*BaseError
}

func NewWriteFileError(expr parser.Expression, message string) error {
	return &WriteFileError{
		NewBaseError(expr, fmt.Sprintf(ERROR_WRITE_FILE, message)),
	}
}

type AutoCommitError struct {
	Message string
}

func (e AutoCommitError) Error() string {
	return e.Message
}

func NewAutoCommitError(message string) error {
	return &AutoCommitError{
		Message: fmt.Sprintf(ERROR_WRITE_FILE_IN_AUTOCOMMIT, message),
	}
}

type FieldAmbiguousError struct {
	*BaseError
}

func NewFieldAmbiguousError(field parser.QueryExpression) error {
	return &FieldAmbiguousError{
		NewBaseError(field, fmt.Sprintf(ERROR_FIELD_AMBIGUOUS, field)),
	}
}

type FieldNotExistError struct {
	*BaseError
}

func NewFieldNotExistError(field parser.QueryExpression) error {
	return &FieldNotExistError{
		NewBaseError(field, fmt.Sprintf(ERROR_FIELD_NOT_EXIST, field)),
	}
}

type FieldNotGroupKeyError struct {
	*BaseError
}

func NewFieldNotGroupKeyError(field parser.QueryExpression) error {
	return &FieldNotGroupKeyError{
		NewBaseError(field, fmt.Sprintf(ERROR_FIELD_NOT_GROUP_KEY, field)),
	}
}

type DuplicateFieldNameError struct {
	*BaseError
}

func NewDuplicateFieldNameError(fieldName parser.Identifier) error {
	return &DuplicateFieldNameError{
		NewBaseError(fieldName, fmt.Sprintf(ERROR_DUPLICATE_FIELD_NAME, fieldName)),
	}
}

type NotGroupingRecordsError struct {
	*BaseError
}

func NewNotGroupingRecordsError(expr parser.QueryExpression, funcname string) error {
	return &NotGroupingRecordsError{
		NewBaseError(expr, fmt.Sprintf(ERROR_NOT_GROUPING_RECORDS, funcname)),
	}
}

type UndefinedVariableError struct {
	*BaseError
}

func NewUndefinedVariableError(expr parser.Variable) error {
	return &UndefinedVariableError{
		NewBaseError(expr, fmt.Sprintf(ERROR_UNDEFINED_VARIABLE, expr)),
	}
}

type VariableRedeclaredError struct {
	*BaseError
}

func NewVariableRedeclaredError(expr parser.Variable) error {
	return &VariableRedeclaredError{
		NewBaseError(expr, fmt.Sprintf(ERROR_VARIABLE_REDECLARED, expr)),
	}
}

type FunctionNotExistError struct {
	*BaseError
}

func NewFunctionNotExistError(expr parser.QueryExpression, funcname string) error {
	return &FunctionNotExistError{
		NewBaseError(expr, fmt.Sprintf(ERROR_FUNCTION_NOT_EXIST, funcname)),
	}
}

type FunctionArgumentLengthError struct {
	*BaseError
}

func NewFunctionArgumentLengthError(expr parser.QueryExpression, funcname string, argslen []int) error {
	var argstr string
	if 1 < len(argslen) {
		lastarg := FormatCount(argslen[len(argslen)-1], "argument")
		strs := make([]string, len(argslen))
		for i := 0; i < len(argslen)-1; i++ {
			strs[i] = strconv.Itoa(argslen[i])
		}
		strs[len(argslen)-1] = lastarg
		argstr = strings.Join(strs, " or ")
	} else {
		argstr = FormatCount(argslen[0], "argument")
		if 0 < argslen[0] {
			argstr = "exactly " + argstr
		}
	}
	return &FunctionArgumentLengthError{
		NewBaseError(expr, fmt.Sprintf(ERROR_FUNCTION_ARGUMENT_LENGTH, funcname, argstr)),
	}
}

func NewFunctionArgumentLengthErrorWithCustomArgs(expr parser.QueryExpression, funcname string, argstr string) error {
	return &FunctionArgumentLengthError{
		NewBaseError(expr, fmt.Sprintf(ERROR_FUNCTION_ARGUMENT_LENGTH, funcname, argstr)),
	}
}

type FunctionInvalidArgumentError struct {
	*BaseError
}

func NewFunctionInvalidArgumentError(function parser.QueryExpression, funcname string, message string) error {
	return &FunctionInvalidArgumentError{
		NewBaseError(function, fmt.Sprintf(ERROR_FUNCTION_INVALID_ARGUMENT, message, funcname)),
	}
}

type UnpermittedStatementFunctionError struct {
	*BaseError
}

func NewUnpermittedStatementFunctionError(expr parser.QueryExpression, funcname string) error {
	return &UnpermittedStatementFunctionError{
		NewBaseError(expr, fmt.Sprintf(ERROR_UNPERMITTED_STATEMENT_FUNCTION, funcname)),
	}
}

type NestedAggregateFunctionsError struct {
	*BaseError
}

func NewNestedAggregateFunctionsError(expr parser.QueryExpression) error {
	return &NestedAggregateFunctionsError{
		NewBaseError(expr, fmt.Sprintf(ERROR_NESTED_AGGREGATE_FUNCTIONS, expr)),
	}
}

type FunctionRedeclaredError struct {
	*BaseError
}

func NewFunctionRedeclaredError(expr parser.Identifier) error {
	return &FunctionRedeclaredError{
		NewBaseError(expr, fmt.Sprintf(ERROR_FUNCTION_REDECLARED, expr.Literal)),
	}
}

type BuiltInFunctionDeclaredError struct {
	*BaseError
}

func NewBuiltInFunctionDeclaredError(expr parser.Identifier) error {
	return &BuiltInFunctionDeclaredError{
		NewBaseError(expr, fmt.Sprintf(ERROR_BUILT_IN_FUNCTION_DECLARED, expr.Literal)),
	}
}

type DuplicateParameterError struct {
	*BaseError
}

func NewDuplicateParameterError(expr parser.Variable) error {
	return &DuplicateParameterError{
		NewBaseError(expr, fmt.Sprintf(ERROR_DUPLICATE_PARAMETER, expr.Name)),
	}
}

type SubqueryTooManyRecordsError struct {
	*BaseError
}

func NewSubqueryTooManyRecordsError(expr parser.Subquery) error {
	return &SubqueryTooManyRecordsError{
		NewBaseError(expr, ERROR_SUBQUERY_TOO_MANY_RECORDS),
	}
}

type SubqueryTooManyFieldsError struct {
	*BaseError
}

func NewSubqueryTooManyFieldsError(expr parser.Subquery) error {
	return &SubqueryTooManyFieldsError{
		NewBaseError(expr, ERROR_SUBQUERY_TOO_MANY_FIELDS),
	}
}

type CursorRedeclaredError struct {
	*BaseError
}

func NewCursorRedeclaredError(cursor parser.Identifier) error {
	return &CursorRedeclaredError{
		NewBaseError(cursor, fmt.Sprintf(ERROR_CURSOR_REDECLARED, cursor)),
	}
}

type UndefinedCursorError struct {
	*BaseError
}

func NewUndefinedCursorError(cursor parser.Identifier) error {
	return &UndefinedCursorError{
		NewBaseError(cursor, fmt.Sprintf(ERROR_UNDEFINED_CURSOR, cursor)),
	}
}

type CursorClosedError struct {
	*BaseError
}

func NewCursorClosedError(cursor parser.Identifier) error {
	return &CursorClosedError{
		NewBaseError(cursor, fmt.Sprintf(ERROR_CURSOR_CLOSED, cursor)),
	}
}

type CursorOpenError struct {
	*BaseError
}

func NewCursorOpenError(cursor parser.Identifier) error {
	return &CursorOpenError{
		NewBaseError(cursor, fmt.Sprintf(ERROR_CURSOR_OPEN, cursor)),
	}
}

type PseudoCursorError struct {
	*BaseError
}

func NewPseudoCursorError(cursor parser.Identifier) error {
	return &PseudoCursorError{
		NewBaseError(cursor, fmt.Sprintf(ERROR_PSEUDO_CURSOR, cursor)),
	}
}

type CursorFetchLengthError struct {
	*BaseError
}

func NewCursorFetchLengthError(cursor parser.Identifier, returnLen int) error {
	return &CursorFetchLengthError{
		NewBaseError(cursor, fmt.Sprintf(ERROR_CURSOR_FETCH_LENGTH, cursor, FormatCount(returnLen, "value"))),
	}
}

type InvalidFetchPositionError struct {
	*BaseError
}

func NewInvalidFetchPositionError(position parser.FetchPosition) error {
	return &InvalidFetchPositionError{
		NewBaseError(position, fmt.Sprintf(ERROR_INVALID_FETCH_POSITION, position.Number)),
	}
}

type InLineTableRedeclaredError struct {
	*BaseError
}

func NewInLineTableRedeclaredError(table parser.Identifier) error {
	return &InLineTableRedeclaredError{
		NewBaseError(table, fmt.Sprintf(ERROR_INLINE_TABLE_REDECLARED, table)),
	}
}

type UndefinedInLineTableError struct {
	*BaseError
}

func NewUndefinedInLineTableError(table parser.Identifier) error {
	return &UndefinedInLineTableError{
		NewBaseError(table, fmt.Sprintf(ERROR_UNDEFINED_INLINE_TABLE, table)),
	}
}

type InlineTableFieldLengthError struct {
	*BaseError
}

func NewInlineTableFieldLengthError(query parser.SelectQuery, table parser.Identifier, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &InlineTableFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ERROR_INLINE_TABLE_FIELD_LENGTH, FormatCount(fieldLen, "field"), table)),
	}
}

type FileNotExistError struct {
	*BaseError
}

func NewFileNotExistError(file parser.Identifier) error {
	return &FileNotExistError{
		NewBaseError(file, fmt.Sprintf(ERROR_FILE_NOT_EXIST, file)),
	}
}

type FileAlreadyExistError struct {
	*BaseError
}

func NewFileAlreadyExistError(file parser.Identifier) error {
	return &FileAlreadyExistError{
		NewBaseError(file, fmt.Sprintf(ERROR_FILE_ALREADY_EXIST, file)),
	}
}

type FileUnableToReadError struct {
	*BaseError
}

func NewFileUnableToReadError(file parser.Identifier) error {
	return &FileUnableToReadError{
		NewBaseError(file, fmt.Sprintf(ERROR_FILE_UNABLE_TO_READ, file)),
	}
}

type FileLockTimeoutError struct {
	*BaseError
}

func NewFileLockTimeoutError(file parser.Identifier, path string) error {
	return &FileLockTimeoutError{
		NewBaseError(file, fmt.Sprintf(ERROR_FILE_LOCK_TIMEOUT, path)),
	}
}

type CsvParsingError struct {
	*BaseError
}

func NewCsvParsingError(file parser.QueryExpression, filepath string, message string) error {
	return &CsvParsingError{
		NewBaseError(file, fmt.Sprintf(ERROR_CSV_PARSING, filepath, message)),
	}
}

type TableFieldLengthError struct {
	*BaseError
}

func NewTableFieldLengthError(query parser.SelectQuery, table parser.Identifier, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &TableFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ERROR_TABLE_FIELD_LENGTH, FormatCount(fieldLen, "field"), table)),
	}
}

type TemporaryTableRedeclaredError struct {
	*BaseError
}

func NewTemporaryTableRedeclaredError(table parser.Identifier) error {
	return &TemporaryTableRedeclaredError{
		NewBaseError(table, fmt.Sprintf(ERROR_TEMPORARY_TABLE_REDECLARED, table)),
	}
}

type UndefinedTemporaryTableError struct {
	*BaseError
}

func NewUndefinedTemporaryTableError(table parser.Identifier) error {
	return &UndefinedTemporaryTableError{
		NewBaseError(table, fmt.Sprintf(ERROR_UNDEFINED_TEMPORARY_TABLE, table)),
	}
}

type TemporaryTableFieldLengthError struct {
	*BaseError
}

func NewTemporaryTableFieldLengthError(query parser.SelectQuery, table parser.Identifier, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &TemporaryTableFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ERROR_TEMPORARY_TABLE_FIELD_LENGTH, FormatCount(fieldLen, "field"), table)),
	}
}

type DuplicateTableNameError struct {
	*BaseError
}

func NewDuplicateTableNameError(table parser.Identifier) error {
	return &DuplicateTableNameError{
		NewBaseError(table, fmt.Sprintf(ERROR_DUPLICATE_TABLE_NAME, table)),
	}
}

type TableNotLoadedError struct {
	*BaseError
}

func NewTableNotLoadedError(table parser.Identifier) error {
	return &TableNotLoadedError{
		NewBaseError(table, fmt.Sprintf(ERROR_TABLE_NOT_LOADED, table)),
	}
}

type StdinEmptyError struct {
	*BaseError
}

func NewStdinEmptyError(stdin parser.Stdin) error {
	return &StdinEmptyError{
		NewBaseError(stdin, ERROR_STDIN_EMPTY),
	}
}

type RowValueLengthInComparisonError struct {
	*BaseError
}

func NewRowValueLengthInComparisonError(expr parser.QueryExpression, valueLen int) error {
	return &RowValueLengthInComparisonError{
		NewBaseError(expr, fmt.Sprintf(ERROR_ROW_VALUE_LENGTH_IN_COMPARISON, FormatCount(valueLen, "value"))),
	}
}

type SelectFieldLengthInComparisonError struct {
	*BaseError
}

func NewSelectFieldLengthInComparisonError(query parser.Subquery, valueLen int) error {
	return &SelectFieldLengthInComparisonError{
		NewBaseError(query, fmt.Sprintf(ERROR_SELECT_FIELD_LENGTH_IN_COMPARISON, FormatCount(valueLen, "field"))),
	}
}

type InvalidLimitPercentageError struct {
	*BaseError
}

func NewInvalidLimitPercentageError(clause parser.LimitClause) error {
	return &InvalidLimitPercentageError{
		NewBaseError(clause, fmt.Sprintf(ERROR_INVALID_LIMIT_PERCENTAGE, clause.Value)),
	}
}

type InvalidLimitNumberError struct {
	*BaseError
}

func NewInvalidLimitNumberError(clause parser.LimitClause) error {
	return &InvalidLimitNumberError{
		NewBaseError(clause, fmt.Sprintf(ERROR_INVALID_LIMIT_NUMBER, clause.Value)),
	}
}

type InvalidOffsetNumberError struct {
	*BaseError
}

func NewInvalidOffsetNumberError(clause parser.OffsetClause) error {
	return &InvalidOffsetNumberError{
		NewBaseError(clause, fmt.Sprintf(ERROR_INVALID_OFFSET_NUMBER, clause.Value)),
	}
}

type CombinedSetFieldLengthError struct {
	*BaseError
}

func NewCombinedSetFieldLengthError(selectEntity parser.QueryExpression, fieldLen int) error {
	selectClause := searchSelectClauseInSelectEntity(selectEntity)

	return &CombinedSetFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ERROR_COMBINED_SET_FIELD_LENGTH, FormatCount(fieldLen, "field"))),
	}
}

type InsertRowValueLengthError struct {
	*BaseError
}

func NewInsertRowValueLengthError(rowValue parser.RowValue, valueLen int) error {
	return &InsertRowValueLengthError{
		NewBaseError(rowValue, fmt.Sprintf(ERROR_INSERT_ROW_VALUE_LENGTH, FormatCount(valueLen, "value"))),
	}
}

type InsertSelectFieldLengthError struct {
	*BaseError
}

func NewInsertSelectFieldLengthError(query parser.SelectQuery, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &InsertSelectFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ERROR_INSERT_SELECT_FIELD_LENGTH, FormatCount(fieldLen, "field"))),
	}
}

type UpdateFieldNotExistError struct {
	*BaseError
}

func NewUpdateFieldNotExistError(field parser.QueryExpression) error {
	return &UpdateFieldNotExistError{
		NewBaseError(field, fmt.Sprintf(ERROR_UPDATE_FIELD_NOT_EXIST, field)),
	}
}

type UpdateValueAmbiguousError struct {
	*BaseError
}

func NewUpdateValueAmbiguousError(field parser.QueryExpression, value parser.QueryExpression) error {
	return &UpdateValueAmbiguousError{
		NewBaseError(field, fmt.Sprintf(ERROR_UPDATE_VALUE_AMBIGUOUS, value, field)),
	}
}

type DeleteTableNotSpecifiedError struct {
	*BaseError
}

func NewDeleteTableNotSpecifiedError(query parser.DeleteQuery) error {
	return &DeleteTableNotSpecifiedError{
		NewBaseError(query, ERROR_DELETE_TABLE_NOT_SPECIFIED),
	}
}

type PrintfReplaceValueLengthError struct {
	*BaseError
}

func NewPrintfReplaceValueLengthError(printf parser.Printf, message string) error {
	return &PrintfReplaceValueLengthError{
		NewBaseError(printf, fmt.Sprintf(ERROR_PRINTF_REPLACE_VALUE_LENGTH, message)),
	}
}

type SourceInvalidArgumentError struct {
	*BaseError
}

func NewSourceInvalidArgumentError(source parser.Source, arg parser.QueryExpression) error {
	return &SourceInvalidArgumentError{
		NewBaseError(source, fmt.Sprintf(ERROR_SOURCE_INVALID_ARGUMENT, arg)),
	}
}

type SourceFileNotExistError struct {
	*BaseError
}

func NewSourceFileNotExistError(source parser.Source, fpath string) error {
	return &SourceFileNotExistError{
		NewBaseError(source, fmt.Sprintf(ERROR_SOURCE_FILE_NOT_EXIST, fpath)),
	}
}

type SourceFileUnableToReadError struct {
	*BaseError
}

func NewSourceFileUnableToReadError(source parser.Source, fpath string) error {
	return &SourceFileUnableToReadError{
		NewBaseError(source, fmt.Sprintf(ERROR_SOURCE_FILE_UNABLE_TO_READ, fpath)),
	}
}

type InvalidFlagNameError struct {
	*BaseError
}

func NewInvalidFlagNameError(expr parser.Expression, name string) error {
	return &InvalidFlagNameError{
		NewBaseError(expr, fmt.Sprintf(ERROR_INVALID_FLAG_NAME, name)),
	}
}

type InvalidFlagValueError struct {
	*BaseError
}

func NewInvalidFlagValueError(setFlag parser.SetFlag) error {
	return &InvalidFlagValueError{
		NewBaseError(setFlag, fmt.Sprintf(ERROR_INVALID_FLAG_VALUE, setFlag.Value, setFlag.Name)),
	}
}

type InternalRecordIdNotExistError struct {
	*BaseError
}

func NewInternalRecordIdNotExistError() error {
	return &InternalRecordIdNotExistError{
		NewBaseError(parser.NewNullValue(), ERROR_INTERNAL_RECORD_ID_NOT_EXIST),
	}
}

type InternalRecordIdEmptyError struct {
	*BaseError
}

func NewInternalRecordIdEmptyError() error {
	return &InternalRecordIdEmptyError{
		NewBaseError(parser.NewNullValue(), ERROR_INTERNAL_RECORD_ID_EMPTY),
	}
}

type FieldLengthNotMatchError struct {
	*BaseError
}

func NewFieldLengthNotMatchError() error {
	return &FieldLengthNotMatchError{
		NewBaseError(parser.NewNullValue(), ERROR_FIELD_LENGTH_NOT_MATCH),
	}
}

type RowValueLengthNotMatchError struct {
	*BaseError
}

func NewRowValueLengthNotMatchError() error {
	return &RowValueLengthNotMatchError{
		NewBaseError(parser.NewNullValue(), ERROR_ROW_VALUE_LENGTH_NOT_MATCH),
	}
}

type RowValueLengthInListError struct {
	*BaseError
	Index int
}

func NewRowValueLengthInListError(i int) error {
	return &RowValueLengthInListError{
		BaseError: NewBaseError(parser.NewNullValue(), fmt.Sprintf(ERROR_ROW_VALUE_LENGTH_IN_LIST, i)),
		Index:     i,
	}
}

type FormatStringLengthNotMatchError struct {
	*BaseError
}

func NewFormatStringLengthNotMatchError() error {
	return &FormatStringLengthNotMatchError{
		BaseError: NewBaseError(parser.NewNullValue(), ERROR_FORMAT_STRING_LENGTH_NOT_MATCH),
	}
}

func searchSelectClause(query parser.SelectQuery) parser.SelectClause {
	return searchSelectClauseInSelectEntity(query.SelectEntity)
}

func searchSelectClauseInSelectEntity(selectEntity parser.QueryExpression) parser.SelectClause {
	if entity, ok := selectEntity.(parser.SelectEntity); ok {
		return entity.SelectClause.(parser.SelectClause)
	}
	return searchSelectClauseInSelectSetEntity(selectEntity.(parser.SelectSet).LHS)
}

func searchSelectClauseInSelectSetEntity(selectSetEntity parser.QueryExpression) parser.SelectClause {
	if subquery, ok := selectSetEntity.(parser.Subquery); ok {
		return searchSelectClause(subquery.Query)
	}
	return searchSelectClauseInSelectEntity(selectSetEntity)
}
