package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type TemporaryViewMapList []ViewMap

func (list TemporaryViewMapList) Exists(name string) bool {
	for _, m := range list {
		if m.Exists(name) {
			return true
		}
	}
	return false
}

func (list TemporaryViewMapList) Get(name parser.Identifier) (*View, error) {
	for _, m := range list {
		if view, err := m.Get(name); err == nil {
			return view, nil
		}
	}
	return nil, NewTableNotLoadedError(name)
}

func (list TemporaryViewMapList) GetWithInternalId(name parser.Identifier) (*View, error) {
	for _, m := range list {
		if view, err := m.GetWithInternalId(name); err == nil {
			return view, nil
		}
	}
	return nil, NewTableNotLoadedError(name)
}

func (list TemporaryViewMapList) Set(view *View) {
	list[0].Set(view)
}

func (list TemporaryViewMapList) Replace(view *View) {
	for _, m := range list {
		if err := m.Replace(view); err == nil {
			return
		}
	}
}

func (list TemporaryViewMapList) Dispose(name parser.Identifier) error {
	for _, m := range list {
		if err := m.DisposeTemporaryTable(name); err == nil {
			return nil
		}
	}
	return NewUndefinedTemporaryTableError(name)
}

func (list TemporaryViewMapList) Rollback() {
	for _, m := range list {
		for _, view := range m {
			view.Rollback()
		}
	}
}

type ViewMap map[string]*View

func (m ViewMap) Exists(fpath string) bool {
	ufpath := strings.ToUpper(fpath)
	if _, ok := m[ufpath]; ok {
		return true
	}
	return false
}

func (m ViewMap) Get(fpath parser.Identifier) (*View, error) {
	ufpath := strings.ToUpper(fpath.Literal)
	if view, ok := m[ufpath]; ok {
		return view.Copy(), nil
	}
	return nil, NewTableNotLoadedError(fpath)
}

func (m ViewMap) GetWithInternalId(fpath parser.Identifier) (*View, error) {
	ufpath := strings.ToUpper(fpath.Literal)
	if view, ok := m[ufpath]; ok {
		ret := view.Copy()

		if 0 < ret.FieldLen() {
			ret.Header = MergeHeader(NewHeaderWithId(ret.Header[0].View, []string{}), ret.Header)

			for i, v := range ret.Records {
				ret.Records[i] = append(Record{NewCell(parser.NewInteger(int64(i)))}, v...)
			}
		}

		return ret, nil
	}
	return nil, NewTableNotLoadedError(fpath)
}

func (m ViewMap) Set(view *View) {
	if view.FileInfo != nil {
		m[strings.ToUpper(view.FileInfo.Path)] = view
	}
}

func (m ViewMap) Replace(view *View) error {
	ufpath := strings.ToUpper(view.FileInfo.Path)
	if ok := m.Exists(ufpath); ok {
		m[ufpath] = view
		return nil
	}
	return NewTableNotLoadedError(parser.Identifier{Literal: view.FileInfo.Path})
}

func (m ViewMap) DisposeTemporaryTable(table parser.Identifier) error {
	uname := strings.ToUpper(table.Literal)
	if v, ok := m[uname]; ok {
		if v.FileInfo.Temporary {
			delete(m, uname)
			return nil
		} else {
			return NewUndefinedTemporaryTableError(table)
		}
	}
	return NewUndefinedTemporaryTableError(table)
}

func (m ViewMap) Clear() {
	for k := range m {
		delete(m, k)
	}
}
