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
	ErrorSystemError = 1000

	//IO Error
	ErrorIOError          = 2000
	ErrorReadFile         = 2001
	ErrorWriteFile        = 2002
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

	//Application Error
	ErrorFieldAmbiguous                       = 16001
	ErrorFieldNotExist                        = 16002
	ErrorFieldNotGroupKey                     = 16003
	ErrorDuplicateFieldName                   = 16004
	ErrorNotGroupingRecords                   = 16005
	ErrorUndeclaredVariable                   = 16006
	ErrorVariableRedeclared                   = 16007
	ErrorFunctionNotExist                     = 16008
	ErrorFunctionArgumentsLength              = 16009
	ErrorFunctionInvalidArgument              = 16010
	ErrorFunctionRedeclared                   = 16011
	ErrorBuiltInFunctionDeclared              = 16012
	ErrorDuplicateParameter                   = 16013
	ErrorSubqueryTooManyRecords               = 16014
	ErrorSubqueryTooManyFields                = 16015
	ErrorJsonQueryTooManyRecords              = 16016
	ErrorLoadJson                             = 16017
	ErrorEmptyJsonQuery                       = 16018
	ErrorEmptyJsonTable                       = 16019
	ErrorInvalidTableObject                   = 16020
	ErrorTableObjectInvalidDelimiter          = 16021
	ErrorTableObjectInvalidDelimiterPositions = 16022
	ErrorTableObjectInvalidJsonQuery          = 16023
	ErrorTableObjectArgumentsLength           = 16024
	ErrorTableObjectJsonArgumentsLength       = 16025
	ErrorTableObjectInvalidArgument           = 16026
	ErrorCursorRedeclared                     = 16027
	ErrorUndeclaredCursor                     = 16028
	ErrorCursorClosed                         = 16029
	ErrorCursorOpen                           = 16030
	ErrorPseudoCursor                         = 16031
	ErrorCursorFetchLength                    = 16032
	ErrorInvalidFetchPosition                 = 16033
	ErrorInlineTableRedefined                 = 16034
	ErrorUndefinedInlineTable                 = 16035
	ErrorInlineTableFieldLength               = 16036
	ErrorFileNameAmbiguous                    = 16037
	ErrorDataParsing                          = 16038
	ErrorTableFieldLength                     = 16039
	ErrorTemporaryTableRedeclared             = 16040
	ErrorUndeclaredTemporaryTable             = 16041
	ErrorTemporaryTableFieldLength            = 16042
	ErrorDuplicateTableName                   = 16043
	ErrorTableNotLoaded                       = 16044
	ErrorStdinEmpty                           = 16045
	ErrorRowValueLengthInComparison           = 16046
	ErrorFieldLengthInComparison              = 16047
	ErrorInvalidLimitPercentage               = 16048
	ErrorInvalidLimitNumber                   = 16049
	ErrorInvalidOffsetNumber                  = 16050
	ErrorCombinedSetFieldLength               = 16051
	ErrorInsertRowValueLength                 = 16052
	ErrorInsertSelectFieldLength              = 16053
	ErrorUpdateFieldNotExist                  = 16054
	ErrorUpdateValueAmbiguous                 = 16055
	ErrorDeleteTableNotSpecified              = 16056
	ErrorShowInvalidObjectType                = 16057
	ErrorReplaceValueLength                   = 16058
	ErrorSourceInvalidFilePath                = 16059
	ErrorInvalidFlagName                      = 16060
	ErrorFlagValueNowAllowedFormat            = 16061
	ErrorInvalidFlagValue                     = 16062
	ErrorAddFlagNotSupportedName              = 16063
	ErrorRemoveFlagNotSupportedName           = 16064
	ErrorInvalidFlagValueToBeRemoved          = 16065
	ErrorInvalidRuntimeInformation            = 16066
	ErrorNotTable                             = 16067
	ErrorInvalidTableAttributeName            = 16068
	ErrorTableAttributeValueNotAllowedFormat  = 16069
	ErrorInvalidTableAttributeValue           = 16070
	ErrorInvalidEventName                     = 16071
	ErrorInternalRecordIdNotExist             = 16072
	ErrorInternalRecordIdEmpty                = 16073
	ErrorFieldLengthNotMatch                  = 16074
	ErrorRowValueLengthInList                 = 16075
	ErrorFormatStringLengthNotMatch           = 16076
	ErrorUnknownFormatPlaceholder             = 16077
	ErrorFormatUnexpectedTermination          = 16078
	ErrorExternalCommand                      = 16079
	ErrorInvalidReloadType                    = 16080
	ErrorLoadConfiguration                    = 16081

	//User Triggered Error
	ErrorExit          = 32000
	ErrorUserTriggered = 32001
)
