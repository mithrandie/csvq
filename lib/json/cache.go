package json

import "strings"

var Path = PathMap{}
var Query = QueryMap{}

type PathMap map[string]PathExpression

func (m PathMap) Parse(s string) (PathExpression, error) {
	if e, ok := m[s]; ok {
		return e, nil
	}
	e, err := ParsePath(s)
	if err != nil {
		return nil, err
	}
	m[s] = e
	return e, nil
}

type QueryMap map[string]QueryExpression

func (m QueryMap) Parse(s string) (QueryExpression, error) {
	s = strings.TrimSpace(s)

	if e, ok := m[s]; ok {
		return e, nil
	}
	e, err := ParseQuery(s)
	if err != nil {
		return nil, err
	}
	m[s] = e
	return e, nil
}
