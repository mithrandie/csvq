package query

const returnCodeBaseSignal = 128
const errorSignalBase = 91280

const (
	ReturnCodeApplicationError          = 1
	ReturnCodeIncorrectUsage            = 2
	ReturnCodeSyntaxError               = 4
	ReturnCodeContextIsDone             = 8
	ReturnCodeIOError                   = 16
	ReturnCodeSystemError               = 32
	ReturnCodeDefaultUserTriggeredError = 64
)

const (
	//Application Error
	ErrorFieldAmbiguous                       = 101
	ErrorFieldNotExist                        = 102
	ErrorFieldNotGroupKey                     = 103
	ErrorDuplicateFieldName                   = 104
	ErrorNotGroupingRecords                   = 201
	ErrorUndeclaredVariable                   = 301
	ErrorVariableRedeclared                   = 302
	ErrorFunctionNotExist                     = 401
	ErrorFunctionArgumentsLength              = 402
	ErrorFunctionInvalidArgument              = 403
	ErrorFunctionRedeclared                   = 501
	ErrorBuiltInFunctionDeclared              = 502
	ErrorDuplicateParameter                   = 503
	ErrorSubqueryTooManyRecords               = 601
	ErrorSubqueryTooManyFields                = 602
	ErrorJsonQueryTooManyRecords              = 701
	ErrorLoadJson                             = 702
	ErrorEmptyJsonQuery                       = 703
	ErrorEmptyJsonTable                       = 801
	ErrorInvalidTableObject                   = 901
	ErrorTableObjectInvalidDelimiter          = 902
	ErrorTableObjectInvalidDelimiterPositions = 903
	ErrorTableObjectInvalidJsonQuery          = 904
	ErrorTableObjectArgumentsLength           = 905
	ErrorTableObjectJsonArgumentsLength       = 906
	ErrorTableObjectInvalidArgument           = 907
	ErrorCursorRedeclared                     = 1001
	ErrorUndeclaredCursor                     = 1002
	ErrorCursorClosed                         = 1003
	ErrorCursorOpen                           = 1004
	ErrorInvalidCursorStatement               = 1005
	ErrorPseudoCursor                         = 1006
	ErrorCursorFetchLength                    = 1007
	ErrorInvalidFetchPosition                 = 1008
	ErrorInlineTableRedefined                 = 1101
	ErrorUndefinedInlineTable                 = 1102
	ErrorInlineTableFieldLength               = 1103
	ErrorFileNameAmbiguous                    = 1201
	ErrorDataParsing                          = 1301
	ErrorTableFieldLength                     = 1401
	ErrorTemporaryTableRedeclared             = 1501
	ErrorUndeclaredTemporaryTable             = 1502
	ErrorTemporaryTableFieldLength            = 1503
	ErrorDuplicateTableName                   = 1601
	ErrorTableNotLoaded                       = 1602
	ErrorStdinEmpty                           = 1603
	ErrorRowValueLengthInComparison           = 1701
	ErrorFieldLengthInComparison              = 1702
	ErrorInvalidLimitPercentage               = 1801
	ErrorInvalidLimitNumber                   = 1802
	ErrorInvalidOffsetNumber                  = 1901
	ErrorCombinedSetFieldLength               = 2001
	ErrorInsertRowValueLength                 = 2101
	ErrorInsertSelectFieldLength              = 2102
	ErrorUpdateFieldNotExist                  = 2201
	ErrorUpdateValueAmbiguous                 = 2202
	ErrorDeleteTableNotSpecified              = 2301
	ErrorShowInvalidObjectType                = 2401
	ErrorReplaceValueLength                   = 2501
	ErrorSourceInvalidFilePath                = 2601
	ErrorInvalidFlagName                      = 2701
	ErrorFlagValueNowAllowedFormat            = 2702
	ErrorInvalidFlagValue                     = 2703
	ErrorAddFlagNotSupportedName              = 2801
	ErrorRemoveFlagNotSupportedName           = 2802
	ErrorInvalidFlagValueToBeRemoved          = 2803
	ErrorInvalidRuntimeInformation            = 2901
	ErrorNotTable                             = 3001
	ErrorInvalidTableAttributeName            = 3002
	ErrorTableAttributeValueNotAllowedFormat  = 3003
	ErrorInvalidTableAttributeValue           = 3004
	ErrorInvalidEventName                     = 3101
	ErrorInternalRecordIdNotExist             = 3201
	ErrorInternalRecordIdEmpty                = 3202
	ErrorFieldLengthNotMatch                  = 3301
	ErrorRowValueLengthInList                 = 3401
	ErrorFormatStringLengthNotMatch           = 3501
	ErrorUnknownFormatPlaceholder             = 3502
	ErrorFormatUnexpectedTermination          = 3503
	ErrorInvalidReloadType                    = 3601
	ErrorLoadConfiguration                    = 3701
	ErrorDuplicateStatementName               = 3801
	ErrorStatementNotExist                    = 3802
	ErrorStatementReplaceValueNotSpecified    = 3803

	//Incorrect Command Usage
	ErrorIncorrectCommandUsage = 90020

	//Syntax Error
	ErrorSyntaxError                  = 90040
	ErrorInvalidValueExpression       = 90041
	ErrorUnpermittedFunctionStatement = 90042
	ErrorNestedAggregateFunctions     = 90043
	ErrorPreparedStatementSyntaxError = 90044

	//Context Error
	ErrorContextIsDone   = 90080
	ErrorFileLockTimeout = 90081

	//IO Error
	ErrorIO               = 90160
	ErrorCommit           = 90171
	ErrorRollback         = 90172
	ErrorInvalidPath      = 90180
	ErrorFileNotExist     = 90181
	ErrorFileAlreadyExist = 90182
	ErrorFileUnableToRead = 90183

	//System Error
	ErrorSystemError     = 90320
	ErrorExternalCommand = 30330

	//User Triggered Error
	ErrorExit          = 90640
	ErrorUserTriggered = 90650
)
