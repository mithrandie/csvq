package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type AliasMapList []AliasMap

func (list AliasMapList) Add(alias parser.Identifier, path string) error {
	return list[0].Add(alias, path)
}

func (list AliasMapList) Get(alias parser.Identifier) (path string, err error) {
	for _, m := range list {
		if path, err = m.Get(alias); err == nil {
			return
		}
	}
	err = NewTableNotLoadedError(alias)
	return
}

func (list AliasMapList) CreateNode() AliasMapList {
	node := make(AliasMapList, len(list)+1)
	node[0] = AliasMap{}
	for i := 0; i < len(list); i++ {
		node[i+1] = list[i]
	}
	return node
}

type AliasMap map[string]string

func (m AliasMap) Add(alias parser.Identifier, path string) error {
	uname := strings.ToUpper(alias.Literal)
	if _, ok := m[uname]; ok {
		return NewDuplicateTableNameError(alias)
	}
	m[uname] = strings.ToUpper(path)
	return nil
}

func (m AliasMap) Get(alias parser.Identifier) (string, error) {
	uname := strings.ToUpper(alias.Literal)
	if fpath, ok := m[uname]; ok {
		if len(fpath) < 1 {
			return "", NewTableNotLoadedError(alias)
		}
		return fpath, nil
	}
	return "", NewTableNotLoadedError(alias)
}
