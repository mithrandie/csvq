package query

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

type TemporaryViewScopes []ViewMap

func (list TemporaryViewScopes) Exists(name string) bool {
	for _, m := range list {
		if m.Exists(name) {
			return true
		}
	}
	return false
}

func (list TemporaryViewScopes) Get(name parser.Identifier) (*View, error) {
	for _, m := range list {
		if view, err := m.Get(name); err == nil {
			return view, nil
		}
	}
	return nil, NewTableNotLoadedError(name)
}

func (list TemporaryViewScopes) GetWithInternalId(name parser.Identifier) (*View, error) {
	for _, m := range list {
		if view, err := m.GetWithInternalId(name); err == nil {
			return view, nil
		}
	}
	return nil, NewTableNotLoadedError(name)
}

func (list TemporaryViewScopes) Set(view *View) {
	list[0].Set(view)
}

func (list TemporaryViewScopes) Replace(view *View) {
	for _, m := range list {
		if err := m.Replace(view); err == nil {
			return
		}
	}
}

func (list TemporaryViewScopes) Dispose(name parser.Identifier) error {
	for _, m := range list {
		if err := m.DisposeTemporaryTable(name); err == nil {
			return nil
		}
	}
	return NewUndeclaredTemporaryTableError(name)
}

func (list TemporaryViewScopes) Store(uncomittedViews map[string]*FileInfo) {
	for _, m := range list {
		for _, view := range m {
			if _, ok := uncomittedViews[view.FileInfo.Path]; ok {
				view.FileInfo.InitialRecordSet = view.RecordSet.Copy()
				view.FileInfo.InitialHeader = view.Header.Copy()
				Log(cmd.Notice(fmt.Sprintf("Commit: restore point of view %q is created.", view.FileInfo.Path)), cmd.GetFlags().Quiet)
			}
		}
	}
}

func (list TemporaryViewScopes) Restore(uncomittedViews map[string]*FileInfo) {
	for _, m := range list {
		for _, view := range m {
			if _, ok := uncomittedViews[view.FileInfo.Path]; ok {
				view.RecordSet = view.FileInfo.InitialRecordSet.Copy()
				view.Header = view.FileInfo.InitialHeader.Copy()
				Log(cmd.Notice(fmt.Sprintf("Rollback: view %q is restored.", view.FileInfo.Path)), cmd.GetFlags().Quiet)
			}
		}
	}
}

func (list TemporaryViewScopes) All() ViewMap {
	all := make(ViewMap, 10)

	for _, m := range list {
		for key, view := range m {
			if !view.FileInfo.IsTemporary {
				continue
			}
			if _, ok := all[key]; !ok {
				all[key] = view
			}
		}
	}
	return all
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

		ret.Header = MergeHeader(NewHeaderWithId(ret.Header[0].View, []string{}), ret.Header)
		fieldLen := ret.FieldLen()

		gm := NewGoroutineManager(ret.RecordLen(), MinimumRequiredForParallelRoutine)
		for i := 0; i < gm.CPU; i++ {
			gm.Add()
			go func(thIdx int) {
				start, end := gm.RecordRange(thIdx)

				for i := start; i < end; i++ {
					record := make(Record, fieldLen)
					record[0] = NewCell(value.NewInteger(int64(i)))
					for j, cell := range ret.RecordSet[i] {
						record[j+1] = cell
					}
					ret.RecordSet[i] = record
				}
				gm.Done()
			}(i)
		}
		gm.Wait()

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

func (m ViewMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m ViewMap) SortedKeys() []string {
	keys := m.Keys()
	sort.Strings(keys)
	return keys
}

func (m ViewMap) DisposeTemporaryTable(table parser.Identifier) error {
	uname := strings.ToUpper(table.Literal)
	if v, ok := m[uname]; ok {
		if v.FileInfo.IsTemporary {
			delete(m, uname)
			return nil
		} else {
			return NewUndeclaredTemporaryTableError(table)
		}
	}
	return NewUndeclaredTemporaryTableError(table)
}

func (m ViewMap) Dispose(name string) error {
	uname := strings.ToUpper(name)
	if _, ok := m[uname]; ok {
		if err := m[uname].FileInfo.Close(); err != nil {
			return err
		}
		delete(m, uname)
	}
	return nil
}

func (m ViewMap) Clean() error {
	for k := range m {
		if err := m.Dispose(k); err != nil {
			return err
		}
	}
	return nil
}

func (m ViewMap) CleanWithErrors() []error {
	var errs []error
	for k := range m {
		if _, ok := m[k]; ok {
			if es := m[k].FileInfo.CloseWithErrors(); es != nil {
				errs = append(errs, es...)
			}
			delete(m, k)
		}
	}
	return errs
}
