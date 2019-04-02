package query

const (
	ReturnCodeSystemError               = 1
	ReturnCodeIOError                   = 2
	ReturnCodeContextIsDone             = 4
	ReturnCodeSyntaxError               = 8
	ReturnCodeApplicationError          = 16
	ReturnCodeDefaultUserTriggeredError = 32
)

const (
	//System Error
	ErrorSystemError     = 1000
	ErrorExternalCommand = 1100

	//IO Error
	ErrorIO               = 2000
	ErrorCommit           = 2100
	ErrorRollback         = 2101
	ErrorInvalidPath      = 2200
	ErrorFileNotExist     = 2201
	ErrorFileAlreadyExist = 2202
	ErrorFileUnableToRead = 2203

	//Context Error
	ErrorContextIsDone   = 4000
	ErrorFileLockTimeout = 4001

	//Syntax Error
	ErrorSyntaxError                  = 8000
	ErrorInvalidValueExpression       = 8001
	ErrorUnpermittedFunctionStatement = 8002
	ErrorNestedAggregateFunctions     = 8003
	ErrorPreparedStatementSyntaxError = 8004

	//Application Error
	ErrorFieldAmbiguous                       = 16001
	ErrorFieldNotExist                        = 16002
	ErrorFieldNotGroupKey                     = 16003
	ErrorDuplicateFieldName                   = 16004
	ErrorNotGroupingRecords                   = 16101
	ErrorUndeclaredVariable                   = 16201
	ErrorVariableRedeclared                   = 16202
	ErrorFunctionNotExist                     = 16301
	ErrorFunctionArgumentsLength              = 16302
	ErrorFunctionInvalidArgument              = 16303
	ErrorFunctionRedeclared                   = 16401
	ErrorBuiltInFunctionDeclared              = 16402
	ErrorDuplicateParameter                   = 16403
	ErrorSubqueryTooManyRecords               = 16501
	ErrorSubqueryTooManyFields                = 16502
	ErrorJsonQueryTooManyRecords              = 16601
	ErrorLoadJson                             = 16602
	ErrorEmptyJsonQuery                       = 16603
	ErrorEmptyJsonTable                       = 16701
	ErrorInvalidTableObject                   = 16801
	ErrorTableObjectInvalidDelimiter          = 16802
	ErrorTableObjectInvalidDelimiterPositions = 16803
	ErrorTableObjectInvalidJsonQuery          = 16804
	ErrorTableObjectArgumentsLength           = 16805
	ErrorTableObjectJsonArgumentsLength       = 16806
	ErrorTableObjectInvalidArgument           = 16807
	ErrorCursorRedeclared                     = 16901
	ErrorUndeclaredCursor                     = 16902
	ErrorCursorClosed                         = 16903
	ErrorCursorOpen                           = 16904
	ErrorInvalidCursorStatement               = 16905
	ErrorPseudoCursor                         = 16906
	ErrorCursorFetchLength                    = 16907
	ErrorInvalidFetchPosition                 = 16908
	ErrorInlineTableRedefined                 = 17001
	ErrorUndefinedInlineTable                 = 17002
	ErrorInlineTableFieldLength               = 17003
	ErrorFileNameAmbiguous                    = 17101
	ErrorDataParsing                          = 17201
	ErrorTableFieldLength                     = 17301
	ErrorTemporaryTableRedeclared             = 17401
	ErrorUndeclaredTemporaryTable             = 17402
	ErrorTemporaryTableFieldLength            = 17403
	ErrorDuplicateTableName                   = 17501
	ErrorTableNotLoaded                       = 17502
	ErrorStdinEmpty                           = 17503
	ErrorRowValueLengthInComparison           = 17601
	ErrorFieldLengthInComparison              = 17602
	ErrorInvalidLimitPercentage               = 17701
	ErrorInvalidLimitNumber                   = 17702
	ErrorInvalidOffsetNumber                  = 17801
	ErrorCombinedSetFieldLength               = 17901
	ErrorInsertRowValueLength                 = 18001
	ErrorInsertSelectFieldLength              = 18002
	ErrorUpdateFieldNotExist                  = 19001
	ErrorUpdateValueAmbiguous                 = 19002
	ErrorDeleteTableNotSpecified              = 19101
	ErrorShowInvalidObjectType                = 19201
	ErrorReplaceValueLength                   = 19301
	ErrorSourceInvalidFilePath                = 19401
	ErrorInvalidFlagName                      = 19501
	ErrorFlagValueNowAllowedFormat            = 19502
	ErrorInvalidFlagValue                     = 19503
	ErrorAddFlagNotSupportedName              = 19601
	ErrorRemoveFlagNotSupportedName           = 19602
	ErrorInvalidFlagValueToBeRemoved          = 19603
	ErrorInvalidRuntimeInformation            = 19701
	ErrorNotTable                             = 19801
	ErrorInvalidTableAttributeName            = 19802
	ErrorTableAttributeValueNotAllowedFormat  = 19803
	ErrorInvalidTableAttributeValue           = 19804
	ErrorInvalidEventName                     = 19901
	ErrorInternalRecordIdNotExist             = 20001
	ErrorInternalRecordIdEmpty                = 20002
	ErrorFieldLengthNotMatch                  = 20101
	ErrorRowValueLengthInList                 = 20201
	ErrorFormatStringLengthNotMatch           = 20301
	ErrorUnknownFormatPlaceholder             = 20302
	ErrorFormatUnexpectedTermination          = 20303
	ErrorInvalidReloadType                    = 20401
	ErrorLoadConfiguration                    = 20501
	ErrorDuplicateStatementName               = 20601
	ErrorStatementNotExist                    = 20602
	ErrorStatementReplaceValueNotSpecified    = 20603

	ErrorIncorrectCommandUsage = 30000

	//User Triggered Error
	ErrorExit          = 32000
	ErrorUserTriggered = 32001
)
