package query

import (
	"fmt"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
	"strconv"
)

const (
	ErrorMessageTemplate                  = "[L:%d C:%d] %s"
	ErrorMessageWithFilepathTemplate      = "%s [L:%d C:%d] %s"
	ErrorMessageWithEmptyPositionTemplate = "[L:- C:-] %s"

	ErrorInvalidSyntax                        = "syntax error: unexpected %s"
	ErrorReadFile                             = "failed to read from file: %s"
	ErrorCreateFile                           = "failed to create file: %s"
	ErrorWriteFile                            = "failed to write to file: %s"
	ErrorWriteFileInAutoCommit                = "[Auto-Commit] failed to write to file: %s"
	ErrorFieldAmbiguous                       = "field %s is ambiguous"
	ErrorFieldNotExist                        = "field %s does not exist"
	ErrorFieldNotGroupKey                     = "field %s is not a group key"
	ErrorDuplicateFieldName                   = "field name %s is a duplicate"
	ErrorNotGroupingRecords                   = "function %s cannot aggregate not grouping records"
	ErrorUndeclaredVariable                   = "variable %s is undeclared"
	ErrorVariableRedeclared                   = "variable %s is redeclared"
	ErrorFunctionNotExist                     = "function %s does not exist"
	ErrorFunctionArgumentsLength              = "function %s takes %s"
	ErrorFunctionInvalidArgument              = "%s for function %s"
	ErrorUnpermittedStatementFunction         = "function %s cannot be used as a statement"
	ErrorNestedAggregateFunctions             = "aggregate functions are nested at %s"
	ErrorFunctionRedeclared                   = "function %s is redeclared"
	ErrorBuiltInFunctionDeclared              = "function %s is a built-in function"
	ErrorDuplicateParameter                   = "parameter %s is a duplicate"
	ErrorSubqueryTooManyRecords               = "subquery returns too many records, should return only one record"
	ErrorSubqueryTooManyFields                = "subquery returns too many fields, should return only one field"
	ErrorJsonQueryTooManyRecords              = "json query returns too many records, should return only one record"
	ErrorJsonQuery                            = "json query error: %s"
	ErrorJsonQueryEmpty                       = "json query is empty"
	ErrorJsonTableEmpty                       = "json table is empty"
	ErrorTableObjectInvalidObject             = "invalid table object: %s"
	ErrorTableObjectInvalidDelimiter          = "invalid delimiter: %s"
	ErrorTableObjectInvalidDelimiterPositions = "invalid delimiter positions: %s"
	ErrorTableObjectInvalidJsonQuery          = "invalid json query: %s"
	ErrorTableObjectMultipleRead              = "file %s has already been loaded"
	ErrorTableObjectArgumentsLength           = "table object %s takes at most %d arguments"
	ErrorTableObjectJsonArgumentsLength       = "table object %s takes exactly %d arguments"
	ErrorTableObjectInvalidArgument           = "invalid argument for %s: %s"
	ErrorCursorRedeclared                     = "cursor %s is redeclared"
	ErrorUndeclaredCursor                     = "cursor %s is undeclared"
	ErrorCursorClosed                         = "cursor %s is closed"
	errorCursorOpen                           = "cursor %s is already open"
	ErrorPseudoCursor                         = "cursor %s is a pseudo cursor"
	ErrorCursorFetchLength                    = "fetching from cursor %s returns %s"
	ErrorInvalidFetchPosition                 = "fetching position %s is not an integer value"
	ErrorInlineTableRedefined                 = "inline table %s is redefined"
	ErrorUndefinedInlineTable                 = "inline table %s is undefined"
	ErrorInlineTableFieldLength               = "select query should return exactly %s for inline table %s"
	ErrorFileNotExist                         = "file %s does not exist"
	ErrorFileAlreadyExist                     = "file %s already exists"
	ErrorFileUnableToRead                     = "file %s is unable to be read"
	ErrorFileLockTimeout                      = "file %s: lock wait timeout period exceeded"
	ErrorFileNameAmbiguous                    = "filename %s is ambiguous"
	ErrorDataParsing                          = "data parse error in file %s: %s"
	ErrorTableFieldLength                     = "select query should return exactly %s for table %s"
	ErrorTemporaryTableRedeclared             = "view %s is redeclared"
	ErrorUndeclaredTemporaryTable             = "view %s is undeclared"
	ErrorTemporaryTableFieldLength            = "select query should return exactly %s for view %s"
	ErrorDuplicateTableName                   = "table name %s is a duplicate"
	ErrorTableNotLoaded                       = "table %s is not loaded"
	ErrorStdinEmpty                           = "stdin is empty"
	ErrorRowValueLengthInComparison           = "row value should contain exactly %s"
	ErrorFieldLengthInComparison              = "select query should return exactly %s"
	ErrorInvalidLimitPercentage               = "limit percentage %s is not a float value"
	ErrorInvalidLimitNumber                   = "limit number of records %s is not an integer value"
	ErrorInvalidOffsetNumber                  = "offset number %s is not an integer value"
	ErrorCombinedSetFieldLength               = "result set to be combined should contain exactly %s"
	ErrorInsertRowValueLength                 = "row value should contain exactly %s"
	ErrorInsertSelectFieldLength              = "select query should return exactly %s"
	ErrorUpdateFieldNotExist                  = "field %s does not exist in the tables to update"
	ErrorUpdateValueAmbiguous                 = "value %s to set in the field %s is ambiguous"
	ErrorDeleteTableNotSpecified              = "tables to delete records are not specified"
	ErrorShowInvalidObjectType                = "object type %s is invalid"
	ErrorPrintfReplaceValueLength             = "%s"
	ErrorSourceInvalidArgument                = "argument %s is not a string"
	ErrorSourceFileNotExist                   = "file %s does not exist"
	ErrorSourceFileUnableToRead               = "file %s is unable to read"
	ErrorInvalidFlagName                      = "flag %s does not exist"
	ErrorFlagValueNowAllowedFormat            = "%s for %s is not allowed"
	ErrorInvalidFlagValue                     = "%s"
	ErrorNotTable                             = "%s is not a table that has attributes"
	ErrorInvalidTableAttributeName            = "table attribute %s does not exist"
	ErrorTableAttributeValueNotAllowedFormat  = "%s for %s is not allowed"
	ErrorInvalidTableAttributeValue           = "%s"
	ErrorInvalidEventName                     = "%s is an unknown event"
	ErrorInternalRecordIdNotExist             = "internal record id does not exist"
	ErrorInternalRecordIdEmpty                = "internal record id is empty"
	ErrorFieldLengthNotMatch                  = "field length does not match"
	ErrorRowValueLengthInList                 = "row value length does not match at index %d"
	ErrorFormatStringLengthNotMatch           = "number of replace values does not match"
)

type ForcedExit struct {
	Code int
}

func NewForcedExit(code int) error {
	return &ForcedExit{
		Code: code,
	}
}

func (e ForcedExit) Error() string {
	return ""
}

func (e ForcedExit) GetCode() int {
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
		return fmt.Sprintf(ErrorMessageWithEmptyPositionTemplate, e.Message)
	}
	if 0 < len(e.SourceFile) {
		return fmt.Sprintf(ErrorMessageWithFilepathTemplate, e.SourceFile, e.Line, e.Char, e.Message)
	}
	return fmt.Sprintf(ErrorMessageTemplate, e.Line, e.Char, e.Message)
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
		NewBaseError(expr, fmt.Sprintf(ErrorInvalidSyntax, expr)),
	}
}

type ReadFileError struct {
	*BaseError
}

func NewReadFileError(expr parser.Expression, message string) error {
	return &ReadFileError{
		NewBaseError(expr, fmt.Sprintf(ErrorReadFile, message)),
	}
}

type CreateFileError struct {
	*BaseError
}

func NewCreateFileError(expr parser.Expression, message string) error {
	return &CreateFileError{
		NewBaseError(expr, fmt.Sprintf(ErrorCreateFile, message)),
	}
}

type WriteFileError struct {
	*BaseError
}

func NewWriteFileError(expr parser.Expression, message string) error {
	return &WriteFileError{
		NewBaseError(expr, fmt.Sprintf(ErrorWriteFile, message)),
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
		Message: fmt.Sprintf(ErrorWriteFileInAutoCommit, message),
	}
}

type FieldAmbiguousError struct {
	*BaseError
}

func NewFieldAmbiguousError(field parser.QueryExpression) error {
	return &FieldAmbiguousError{
		NewBaseError(field, fmt.Sprintf(ErrorFieldAmbiguous, field)),
	}
}

type FieldNotExistError struct {
	*BaseError
}

func NewFieldNotExistError(field parser.QueryExpression) error {
	return &FieldNotExistError{
		NewBaseError(field, fmt.Sprintf(ErrorFieldNotExist, field)),
	}
}

type FieldNotGroupKeyError struct {
	*BaseError
}

func NewFieldNotGroupKeyError(field parser.QueryExpression) error {
	return &FieldNotGroupKeyError{
		NewBaseError(field, fmt.Sprintf(ErrorFieldNotGroupKey, field)),
	}
}

type DuplicateFieldNameError struct {
	*BaseError
}

func NewDuplicateFieldNameError(fieldName parser.Identifier) error {
	return &DuplicateFieldNameError{
		NewBaseError(fieldName, fmt.Sprintf(ErrorDuplicateFieldName, fieldName)),
	}
}

type NotGroupingRecordsError struct {
	*BaseError
}

func NewNotGroupingRecordsError(expr parser.QueryExpression, funcname string) error {
	return &NotGroupingRecordsError{
		NewBaseError(expr, fmt.Sprintf(ErrorNotGroupingRecords, funcname)),
	}
}

type UndeclaredVariableError struct {
	*BaseError
}

func NewUndeclaredVariableError(expr parser.Variable) error {
	return &UndeclaredVariableError{
		NewBaseError(expr, fmt.Sprintf(ErrorUndeclaredVariable, expr)),
	}
}

type VariableRedeclaredError struct {
	*BaseError
}

func NewVariableRedeclaredError(expr parser.Variable) error {
	return &VariableRedeclaredError{
		NewBaseError(expr, fmt.Sprintf(ErrorVariableRedeclared, expr)),
	}
}

type FunctionNotExistError struct {
	*BaseError
}

func NewFunctionNotExistError(expr parser.QueryExpression, funcname string) error {
	return &FunctionNotExistError{
		NewBaseError(expr, fmt.Sprintf(ErrorFunctionNotExist, funcname)),
	}
}

type FunctionArgumentLengthError struct {
	*BaseError
}

func NewFunctionArgumentLengthError(expr parser.QueryExpression, funcname string, argslen []int) error {
	var argstr string
	if 1 < len(argslen) {
		first := argslen[0]
		last := argslen[len(argslen)-1]
		lastarg := FormatCount(last, "argument")
		if len(argslen) == 2 {
			argstr = strconv.Itoa(first) + " or " + lastarg
		} else {
			argstr = strconv.Itoa(first) + " to " + lastarg
		}
	} else {
		argstr = FormatCount(argslen[0], "argument")
		if 0 < argslen[0] {
			argstr = "exactly " + argstr
		}
	}
	return &FunctionArgumentLengthError{
		NewBaseError(expr, fmt.Sprintf(ErrorFunctionArgumentsLength, funcname, argstr)),
	}
}

func NewFunctionArgumentLengthErrorWithCustomArgs(expr parser.QueryExpression, funcname string, argstr string) error {
	return &FunctionArgumentLengthError{
		NewBaseError(expr, fmt.Sprintf(ErrorFunctionArgumentsLength, funcname, argstr)),
	}
}

type FunctionInvalidArgumentError struct {
	*BaseError
}

func NewFunctionInvalidArgumentError(function parser.QueryExpression, funcname string, message string) error {
	return &FunctionInvalidArgumentError{
		NewBaseError(function, fmt.Sprintf(ErrorFunctionInvalidArgument, message, funcname)),
	}
}

type UnpermittedStatementFunctionError struct {
	*BaseError
}

func NewUnpermittedStatementFunctionError(expr parser.QueryExpression, funcname string) error {
	return &UnpermittedStatementFunctionError{
		NewBaseError(expr, fmt.Sprintf(ErrorUnpermittedStatementFunction, funcname)),
	}
}

type NestedAggregateFunctionsError struct {
	*BaseError
}

func NewNestedAggregateFunctionsError(expr parser.QueryExpression) error {
	return &NestedAggregateFunctionsError{
		NewBaseError(expr, fmt.Sprintf(ErrorNestedAggregateFunctions, expr)),
	}
}

type FunctionRedeclaredError struct {
	*BaseError
}

func NewFunctionRedeclaredError(expr parser.Identifier) error {
	return &FunctionRedeclaredError{
		NewBaseError(expr, fmt.Sprintf(ErrorFunctionRedeclared, expr.Literal)),
	}
}

type BuiltInFunctionDeclaredError struct {
	*BaseError
}

func NewBuiltInFunctionDeclaredError(expr parser.Identifier) error {
	return &BuiltInFunctionDeclaredError{
		NewBaseError(expr, fmt.Sprintf(ErrorBuiltInFunctionDeclared, expr.Literal)),
	}
}

type DuplicateParameterError struct {
	*BaseError
}

func NewDuplicateParameterError(expr parser.Variable) error {
	return &DuplicateParameterError{
		NewBaseError(expr, fmt.Sprintf(ErrorDuplicateParameter, expr.Name)),
	}
}

type SubqueryTooManyRecordsError struct {
	*BaseError
}

func NewSubqueryTooManyRecordsError(expr parser.Subquery) error {
	return &SubqueryTooManyRecordsError{
		NewBaseError(expr, ErrorSubqueryTooManyRecords),
	}
}

type SubqueryTooManyFieldsError struct {
	*BaseError
}

func NewSubqueryTooManyFieldsError(expr parser.Subquery) error {
	return &SubqueryTooManyFieldsError{
		NewBaseError(expr, ErrorSubqueryTooManyFields),
	}
}

type JsonQueryTooManyRecordsError struct {
	*BaseError
}

func NewJsonQueryTooManyRecordsError(expr parser.JsonQuery) error {
	return &JsonQueryTooManyRecordsError{
		NewBaseError(expr, ErrorJsonQueryTooManyRecords),
	}
}

type JsonQueryError struct {
	*BaseError
}

func NewJsonQueryError(expr parser.JsonQuery, message string) error {
	return &JsonQueryError{
		NewBaseError(expr, fmt.Sprintf(ErrorJsonQuery, message)),
	}
}

type JsonQueryEmptyError struct {
	*BaseError
}

func NewJsonQueryEmptyError(expr parser.JsonQuery) error {
	return &JsonQueryEmptyError{
		NewBaseError(expr, ErrorJsonQueryEmpty),
	}
}

type JsonTableEmptyError struct {
	*BaseError
}

func NewJsonTableEmptyError(expr parser.JsonQuery) error {
	return &JsonTableEmptyError{
		NewBaseError(expr, ErrorJsonTableEmpty),
	}
}

type TableObjectInvalidObjectError struct {
	*BaseError
}

func NewTableObjectInvalidObjectError(expr parser.TableObject, objectName string) error {
	return &TableObjectInvalidObjectError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectInvalidObject, objectName)),
	}
}

type TableObjectInvalidDelimiterError struct {
	*BaseError
}

func NewTableObjectInvalidDelimiterError(expr parser.TableObject, delimiter string) error {
	return &TableObjectInvalidObjectError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectInvalidDelimiter, delimiter)),
	}
}

type TableObjectInvalidDelimiterPositionsError struct {
	*BaseError
}

func NewTableObjectInvalidDelimiterPositionsError(expr parser.TableObject, positions string) error {
	return &TableObjectInvalidObjectError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectInvalidDelimiterPositions, positions)),
	}
}

type TableObjectInvalidJsonQueryError struct {
	*BaseError
}

func NewTableObjectInvalidJsonQueryError(expr parser.TableObject, jsonQuery string) error {
	return &TableObjectInvalidObjectError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectInvalidJsonQuery, jsonQuery)),
	}
}

type TableObjectMultipleReadError struct {
	*BaseError
}

func NewTableObjectMultipleReadError(tableIdentifier parser.Identifier) error {
	return &TableObjectMultipleReadError{
		NewBaseError(tableIdentifier, fmt.Sprintf(ErrorTableObjectMultipleRead, tableIdentifier)),
	}
}

type TableObjectArgumentsLengthError struct {
	*BaseError
}

func NewTableObjectArgumentsLengthError(expr parser.TableObject, argLen int) error {
	return &TableObjectArgumentsLengthError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectArgumentsLength, expr.Type.Literal, argLen)),
	}
}

type TableObjectJsonArgumentsLengthError struct {
	*BaseError
}

func NewTableObjectJsonArgumentsLengthError(expr parser.TableObject, argLen int) error {
	return &TableObjectJsonArgumentsLengthError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectJsonArgumentsLength, expr.Type.Literal, argLen)),
	}
}

type TableObjectInvalidArgumentError struct {
	*BaseError
}

func NewTableObjectInvalidArgumentError(expr parser.TableObject, message string) error {
	return &TableObjectInvalidArgumentError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableObjectInvalidArgument, expr.Type.Literal, message)),
	}
}

type CursorRedeclaredError struct {
	*BaseError
}

func NewCursorRedeclaredError(cursor parser.Identifier) error {
	return &CursorRedeclaredError{
		NewBaseError(cursor, fmt.Sprintf(ErrorCursorRedeclared, cursor)),
	}
}

type UndeclaredCursorError struct {
	*BaseError
}

func NewUndeclaredCursorError(cursor parser.Identifier) error {
	return &UndeclaredCursorError{
		NewBaseError(cursor, fmt.Sprintf(ErrorUndeclaredCursor, cursor)),
	}
}

type CursorClosedError struct {
	*BaseError
}

func NewCursorClosedError(cursor parser.Identifier) error {
	return &CursorClosedError{
		NewBaseError(cursor, fmt.Sprintf(ErrorCursorClosed, cursor)),
	}
}

type CursorOpenError struct {
	*BaseError
}

func NewCursorOpenError(cursor parser.Identifier) error {
	return &CursorOpenError{
		NewBaseError(cursor, fmt.Sprintf(errorCursorOpen, cursor)),
	}
}

type PseudoCursorError struct {
	*BaseError
}

func NewPseudoCursorError(cursor parser.Identifier) error {
	return &PseudoCursorError{
		NewBaseError(cursor, fmt.Sprintf(ErrorPseudoCursor, cursor)),
	}
}

type CursorFetchLengthError struct {
	*BaseError
}

func NewCursorFetchLengthError(cursor parser.Identifier, returnLen int) error {
	return &CursorFetchLengthError{
		NewBaseError(cursor, fmt.Sprintf(ErrorCursorFetchLength, cursor, FormatCount(returnLen, "value"))),
	}
}

type InvalidFetchPositionError struct {
	*BaseError
}

func NewInvalidFetchPositionError(position parser.FetchPosition) error {
	return &InvalidFetchPositionError{
		NewBaseError(position, fmt.Sprintf(ErrorInvalidFetchPosition, position.Number)),
	}
}

type InLineTableRedefinedError struct {
	*BaseError
}

func NewInLineTableRedefinedError(table parser.Identifier) error {
	return &InLineTableRedefinedError{
		NewBaseError(table, fmt.Sprintf(ErrorInlineTableRedefined, table)),
	}
}

type UndefinedInLineTableError struct {
	*BaseError
}

func NewUndefinedInLineTableError(table parser.Identifier) error {
	return &UndefinedInLineTableError{
		NewBaseError(table, fmt.Sprintf(ErrorUndefinedInlineTable, table)),
	}
}

type InlineTableFieldLengthError struct {
	*BaseError
}

func NewInlineTableFieldLengthError(query parser.SelectQuery, table parser.Identifier, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &InlineTableFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ErrorInlineTableFieldLength, FormatCount(fieldLen, "field"), table)),
	}
}

type FileNotExistError struct {
	*BaseError
}

func NewFileNotExistError(file parser.Identifier) error {
	return &FileNotExistError{
		NewBaseError(file, fmt.Sprintf(ErrorFileNotExist, file)),
	}
}

type FileAlreadyExistError struct {
	*BaseError
}

func NewFileAlreadyExistError(file parser.Identifier) error {
	return &FileAlreadyExistError{
		NewBaseError(file, fmt.Sprintf(ErrorFileAlreadyExist, file)),
	}
}

type FileUnableToReadError struct {
	*BaseError
}

func NewFileUnableToReadError(file parser.Identifier) error {
	return &FileUnableToReadError{
		NewBaseError(file, fmt.Sprintf(ErrorFileUnableToRead, file)),
	}
}

type FileLockTimeoutError struct {
	*BaseError
}

func NewFileLockTimeoutError(file parser.Identifier, path string) error {
	return &FileLockTimeoutError{
		NewBaseError(file, fmt.Sprintf(ErrorFileLockTimeout, path)),
	}
}

type FileNameAmbiguousError struct {
	*BaseError
}

func NewFileNameAmbiguousError(file parser.Identifier) error {
	return &FileNameAmbiguousError{
		NewBaseError(file, fmt.Sprintf(ErrorFileNameAmbiguous, file)),
	}
}

type DataParsingError struct {
	*BaseError
}

func NewDataParsingError(file parser.QueryExpression, filepath string, message string) error {
	return &DataParsingError{
		NewBaseError(file, fmt.Sprintf(ErrorDataParsing, filepath, message)),
	}
}

type TableFieldLengthError struct {
	*BaseError
}

func NewTableFieldLengthError(query parser.SelectQuery, table parser.Identifier, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &TableFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ErrorTableFieldLength, FormatCount(fieldLen, "field"), table)),
	}
}

type TemporaryTableRedeclaredError struct {
	*BaseError
}

func NewTemporaryTableRedeclaredError(table parser.Identifier) error {
	return &TemporaryTableRedeclaredError{
		NewBaseError(table, fmt.Sprintf(ErrorTemporaryTableRedeclared, table)),
	}
}

type UndeclaredTemporaryTableError struct {
	*BaseError
}

func NewUndeclaredTemporaryTableError(table parser.Identifier) error {
	return &UndeclaredTemporaryTableError{
		NewBaseError(table, fmt.Sprintf(ErrorUndeclaredTemporaryTable, table)),
	}
}

type TemporaryTableFieldLengthError struct {
	*BaseError
}

func NewTemporaryTableFieldLengthError(query parser.SelectQuery, table parser.Identifier, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &TemporaryTableFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ErrorTemporaryTableFieldLength, FormatCount(fieldLen, "field"), table)),
	}
}

type DuplicateTableNameError struct {
	*BaseError
}

func NewDuplicateTableNameError(table parser.Identifier) error {
	return &DuplicateTableNameError{
		NewBaseError(table, fmt.Sprintf(ErrorDuplicateTableName, table)),
	}
}

type TableNotLoadedError struct {
	*BaseError
}

func NewTableNotLoadedError(table parser.Identifier) error {
	return &TableNotLoadedError{
		NewBaseError(table, fmt.Sprintf(ErrorTableNotLoaded, table)),
	}
}

type StdinEmptyError struct {
	*BaseError
}

func NewStdinEmptyError(stdin parser.Stdin) error {
	return &StdinEmptyError{
		NewBaseError(stdin, ErrorStdinEmpty),
	}
}

type RowValueLengthInComparisonError struct {
	*BaseError
}

func NewRowValueLengthInComparisonError(expr parser.QueryExpression, valueLen int) error {
	return &RowValueLengthInComparisonError{
		NewBaseError(expr, fmt.Sprintf(ErrorRowValueLengthInComparison, FormatCount(valueLen, "value"))),
	}
}

type SelectFieldLengthInComparisonError struct {
	*BaseError
}

func NewSelectFieldLengthInComparisonError(query parser.Subquery, valueLen int) error {
	return &SelectFieldLengthInComparisonError{
		NewBaseError(query, fmt.Sprintf(ErrorFieldLengthInComparison, FormatCount(valueLen, "field"))),
	}
}

type InvalidLimitPercentageError struct {
	*BaseError
}

func NewInvalidLimitPercentageError(clause parser.LimitClause) error {
	return &InvalidLimitPercentageError{
		NewBaseError(clause, fmt.Sprintf(ErrorInvalidLimitPercentage, clause.Value)),
	}
}

type InvalidLimitNumberError struct {
	*BaseError
}

func NewInvalidLimitNumberError(clause parser.LimitClause) error {
	return &InvalidLimitNumberError{
		NewBaseError(clause, fmt.Sprintf(ErrorInvalidLimitNumber, clause.Value)),
	}
}

type InvalidOffsetNumberError struct {
	*BaseError
}

func NewInvalidOffsetNumberError(clause parser.OffsetClause) error {
	return &InvalidOffsetNumberError{
		NewBaseError(clause, fmt.Sprintf(ErrorInvalidOffsetNumber, clause.Value)),
	}
}

type CombinedSetFieldLengthError struct {
	*BaseError
}

func NewCombinedSetFieldLengthError(selectEntity parser.QueryExpression, fieldLen int) error {
	selectClause := searchSelectClauseInSelectEntity(selectEntity)

	return &CombinedSetFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ErrorCombinedSetFieldLength, FormatCount(fieldLen, "field"))),
	}
}

type InsertRowValueLengthError struct {
	*BaseError
}

func NewInsertRowValueLengthError(rowValue parser.RowValue, valueLen int) error {
	return &InsertRowValueLengthError{
		NewBaseError(rowValue, fmt.Sprintf(ErrorInsertRowValueLength, FormatCount(valueLen, "value"))),
	}
}

type InsertSelectFieldLengthError struct {
	*BaseError
}

func NewInsertSelectFieldLengthError(query parser.SelectQuery, fieldLen int) error {
	selectClause := searchSelectClause(query)

	return &InsertSelectFieldLengthError{
		NewBaseError(selectClause, fmt.Sprintf(ErrorInsertSelectFieldLength, FormatCount(fieldLen, "field"))),
	}
}

type UpdateFieldNotExistError struct {
	*BaseError
}

func NewUpdateFieldNotExistError(field parser.QueryExpression) error {
	return &UpdateFieldNotExistError{
		NewBaseError(field, fmt.Sprintf(ErrorUpdateFieldNotExist, field)),
	}
}

type UpdateValueAmbiguousError struct {
	*BaseError
}

func NewUpdateValueAmbiguousError(field parser.QueryExpression, value parser.QueryExpression) error {
	return &UpdateValueAmbiguousError{
		NewBaseError(field, fmt.Sprintf(ErrorUpdateValueAmbiguous, value, field)),
	}
}

type DeleteTableNotSpecifiedError struct {
	*BaseError
}

func NewDeleteTableNotSpecifiedError(query parser.DeleteQuery) error {
	return &DeleteTableNotSpecifiedError{
		NewBaseError(query, ErrorDeleteTableNotSpecified),
	}
}

type ShowInvalidObjectTypeError struct {
	*BaseError
}

func NewShowInvalidObjectTypeError(expr parser.Expression, objectType string) error {
	return &ShowInvalidObjectTypeError{
		NewBaseError(expr, fmt.Sprintf(ErrorShowInvalidObjectType, objectType)),
	}
}

type PrintfReplaceValueLengthError struct {
	*BaseError
}

func NewPrintfReplaceValueLengthError(printf parser.Printf, message string) error {
	return &PrintfReplaceValueLengthError{
		NewBaseError(printf, fmt.Sprintf(ErrorPrintfReplaceValueLength, message)),
	}
}

type SourceInvalidArgumentError struct {
	*BaseError
}

func NewSourceInvalidArgumentError(source parser.Source, arg parser.QueryExpression) error {
	return &SourceInvalidArgumentError{
		NewBaseError(source, fmt.Sprintf(ErrorSourceInvalidArgument, arg)),
	}
}

type SourceFileNotExistError struct {
	*BaseError
}

func NewSourceFileNotExistError(source parser.Source, fpath string) error {
	return &SourceFileNotExistError{
		NewBaseError(source, fmt.Sprintf(ErrorSourceFileNotExist, fpath)),
	}
}

type SourceFileUnableToReadError struct {
	*BaseError
}

func NewSourceFileUnableToReadError(source parser.Source, fpath string) error {
	return &SourceFileUnableToReadError{
		NewBaseError(source, fmt.Sprintf(ErrorSourceFileUnableToRead, fpath)),
	}
}

type InvalidFlagNameError struct {
	*BaseError
}

func NewInvalidFlagNameError(expr parser.Expression, name string) error {
	return &InvalidFlagNameError{
		NewBaseError(expr, fmt.Sprintf(ErrorInvalidFlagName, name)),
	}
}

type FlagValueNotAllowedFormatError struct {
	*BaseError
}

func NewFlagValueNotAllowedFormatError(setFlag parser.SetFlag) error {
	return &FlagValueNotAllowedFormatError{
		NewBaseError(setFlag, fmt.Sprintf(ErrorFlagValueNowAllowedFormat, setFlag.Value, setFlag.Name)),
	}
}

type InvalidFlagValueError struct {
	*BaseError
}

func NewInvalidFlagValueError(setFlag parser.SetFlag, message string) error {
	return &InvalidFlagValueError{
		NewBaseError(setFlag, fmt.Sprintf(ErrorInvalidFlagValue, message)),
	}
}

type NotTableError struct {
	*BaseError
}

func NewNotTableError(expr parser.QueryExpression) error {
	return &NotTableError{
		NewBaseError(expr, fmt.Sprintf(ErrorNotTable, expr)),
	}
}

type InvalidTableAttributeNameError struct {
	*BaseError
}

func NewInvalidTableAttributeNameError(expr parser.Identifier) error {
	return &InvalidTableAttributeNameError{
		NewBaseError(expr, fmt.Sprintf(ErrorInvalidTableAttributeName, expr)),
	}
}

type TableAttributeValueNotAllowedFormatError struct {
	*BaseError
}

func NewTableAttributeValueNotAllowedFormatError(expr parser.SetTableAttribute) error {
	return &TableAttributeValueNotAllowedFormatError{
		NewBaseError(expr, fmt.Sprintf(ErrorTableAttributeValueNotAllowedFormat, expr.Value, expr.Attribute)),
	}
}

type InvalidTableAttributeValueError struct {
	*BaseError
}

func NewInvalidTableAttributeValueError(expr parser.SetTableAttribute, message string) error {
	return &InvalidTableAttributeValueError{
		NewBaseError(expr, fmt.Sprintf(ErrorInvalidTableAttributeValue, message)),
	}
}

type InvalidEventNameError struct {
	*BaseError
}

func NewInvalidEventNameError(expr parser.Identifier) error {
	return &InvalidEventNameError{
		NewBaseError(expr, fmt.Sprintf(ErrorInvalidEventName, expr)),
	}
}

type InternalRecordIdNotExistError struct {
	*BaseError
}

func NewInternalRecordIdNotExistError() error {
	return &InternalRecordIdNotExistError{
		NewBaseError(parser.NewNullValue(), ErrorInternalRecordIdNotExist),
	}
}

type InternalRecordIdEmptyError struct {
	*BaseError
}

func NewInternalRecordIdEmptyError() error {
	return &InternalRecordIdEmptyError{
		NewBaseError(parser.NewNullValue(), ErrorInternalRecordIdEmpty),
	}
}

type FieldLengthNotMatchError struct {
	*BaseError
}

func NewFieldLengthNotMatchError() error {
	return &FieldLengthNotMatchError{
		NewBaseError(parser.NewNullValue(), ErrorFieldLengthNotMatch),
	}
}

type RowValueLengthInListError struct {
	*BaseError
	Index int
}

func NewRowValueLengthInListError(i int) error {
	return &RowValueLengthInListError{
		BaseError: NewBaseError(parser.NewNullValue(), fmt.Sprintf(ErrorRowValueLengthInList, i)),
		Index:     i,
	}
}

type FormatStringLengthNotMatchError struct {
	*BaseError
}

func NewFormatStringLengthNotMatchError() error {
	return &FormatStringLengthNotMatchError{
		BaseError: NewBaseError(parser.NewNullValue(), ErrorFormatStringLengthNotMatch),
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
