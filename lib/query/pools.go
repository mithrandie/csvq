package query

import (
	"sync"
)

var variableMapPool = sync.Pool{
	New: func() interface{} {
		return NewVariableMap()
	},
}

func GetVariableMap() VariableMap {
	t := variableMapPool.Get().(VariableMap)
	t.Clear()
	return t
}

func PutVariableMap(t VariableMap) {
	variableMapPool.Put(t)
}

var temporaryViewMapPool = sync.Pool{
	New: func() interface{} {
		return NewViewMap()
	},
}

func GetTemporaryViewMap() ViewMap {
	t := temporaryViewMapPool.Get().(ViewMap)
	t.Clear()
	return t
}

func PutTemporaryViewMap(t ViewMap) {
	temporaryViewMapPool.Put(t)
}

var cursorMapPool = sync.Pool{
	New: func() interface{} {
		return make(CursorMap)
	},
}

func GetCursorMap() CursorMap {
	t := cursorMapPool.Get().(CursorMap)
	t.Clear()
	return t
}

func PutCursorMap(t CursorMap) {
	cursorMapPool.Put(t)
}

var userDefinedFunctionMapPool = sync.Pool{
	New: func() interface{} {
		return make(UserDefinedFunctionMap)
	},
}

func GetUserDefinedFunctionMap() UserDefinedFunctionMap {
	t := userDefinedFunctionMapPool.Get().(UserDefinedFunctionMap)
	t.Clear()
	return t
}

func PutUserDefinedFunctionMap(t UserDefinedFunctionMap) {
	userDefinedFunctionMapPool.Put(t)
}

var inlineTableMapPool = sync.Pool{
	New: func() interface{} {
		return make(InlineTableMap)
	},
}

func GetInlineTableMap() InlineTableMap {
	t := inlineTableMapPool.Get().(InlineTableMap)
	t.Clear()
	return t
}

func PutInlineTableMap(t InlineTableMap) {
	inlineTableMapPool.Put(t)
}

var aliasMapPool = sync.Pool{
	New: func() interface{} {
		return make(AliasMap)
	},
}

func GetAliasMap() AliasMap {
	t := aliasMapPool.Get().(AliasMap)
	t.Clear()
	return t
}

func PutAliasMap(t AliasMap) {
	aliasMapPool.Put(t)
}
