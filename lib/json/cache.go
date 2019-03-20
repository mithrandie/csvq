package json

import (
	"strings"
	"sync"
)

var Path = NewPathMap()
var Query = NewQueryMap()

type PathMap struct {
	m *sync.Map
}

func NewPathMap() *PathMap {
	return &PathMap{
		m: &sync.Map{},
	}
}

func (pmap *PathMap) Store(key string, value PathExpression) {
	pmap.m.Store(key, value)
}

func (pmap *PathMap) Load(key string) (PathExpression, bool) {
	v, ok := pmap.m.Load(key)
	if ok && v != nil {
		return v.(PathExpression), ok
	}
	return nil, ok
}

func (pmap *PathMap) Parse(s string) (PathExpression, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return nil, nil
	}

	if e, ok := pmap.Load(s); ok {
		return e, nil
	}
	e, err := ParsePath(s)
	if err != nil || e == nil {
		return nil, err
	}
	pmap.Store(s, e)
	return e, nil
}

type QueryMap struct {
	m *sync.Map
}

func NewQueryMap() *QueryMap {
	return &QueryMap{
		m: &sync.Map{},
	}
}

func (qmap *QueryMap) Store(key string, value QueryExpression) {
	qmap.m.Store(key, value)
}

func (qmap *QueryMap) Load(key string) (QueryExpression, bool) {
	v, ok := qmap.m.Load(key)
	if ok && v != nil {
		return v.(QueryExpression), ok
	}
	return nil, ok
}
func (qmap QueryMap) Parse(s string) (QueryExpression, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return nil, nil
	}

	if e, ok := qmap.Load(s); ok {
		return e, nil
	}
	e, err := ParseQuery(s)
	if err != nil || e == nil {
		return nil, err
	}
	qmap.Store(s, e)
	return e, nil
}
