package json

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

var Path = PathMap{}

type QueryMap map[string]QueryExpression

func (m QueryMap) Parse(s string) (QueryExpression, error) {
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

var Query = QueryMap{}
