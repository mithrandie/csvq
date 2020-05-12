package query

import (
	"context"
	"errors"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/csvq/lib/file"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

var errTableNotLoaded = errors.New("table not loaded")

type ViewMap struct {
	*SyncMap
}

func NewViewMap() ViewMap {
	return ViewMap{
		NewSyncMap(),
	}
}

func (m ViewMap) IsEmpty() bool {
	return m.SyncMap == nil
}

func (m ViewMap) Store(fpath string, view *View) {
	m.store(strings.ToUpper(fpath), view)
}

func (m ViewMap) LoadDirect(name string) (interface{}, bool) {
	return m.load(strings.ToUpper(name))
}

func (m ViewMap) Load(fpath string) (*View, bool) {
	if v, ok := m.load(strings.ToUpper(fpath)); ok {
		return v.(*View), true
	}
	return nil, false
}

func (m ViewMap) Delete(fpath string) {
	m.delete(strings.ToUpper(fpath))
}

func (m ViewMap) Exists(fpath string) bool {
	return m.exists(strings.ToUpper(fpath))
}

func (m ViewMap) Get(fpath parser.Identifier) (*View, error) {
	if view, ok := m.Load(fpath.Literal); ok {
		return view.Copy(), nil
	}
	return nil, errTableNotLoaded
}

func (m ViewMap) GetWithInternalId(ctx context.Context, fpath parser.Identifier, flags *cmd.Flags) (*View, error) {
	if view, ok := m.Load(fpath.Literal); ok {
		ret := view.Copy()

		ret.Header = NewHeaderWithId(ret.Header[0].View, []string{}).Merge(ret.Header)

		if err := NewGoroutineTaskManager(ret.RecordLen(), -1, flags.CPU).Run(ctx, func(index int) error {
			record := make(Record, len(ret.RecordSet[index])+1)
			record[0] = NewCell(value.NewInteger(int64(index)))
			for i := 0; i < len(ret.RecordSet[index]); i++ {
				record[i+1] = ret.RecordSet[index][i]
			}
			ret.RecordSet[index] = record
			return nil
		}); err != nil {
			return nil, err
		}

		return ret, nil
	}
	return nil, errTableNotLoaded
}

func (m ViewMap) Set(view *View) {
	if view.FileInfo != nil {
		m.Store(view.FileInfo.Path, view)
	}
}

func (m ViewMap) DisposeTemporaryTable(table parser.QueryExpression) bool {
	var tableName string
	if e, ok := table.(parser.Stdin); ok {
		tableName = e.String()
	} else {
		tableName = table.(parser.Identifier).Literal
	}

	if v, ok := m.Load(tableName); ok && !v.FileInfo.IsFile() {
		m.Delete(tableName)
		return true
	}
	return false
}

func (m ViewMap) Dispose(container *file.Container, name string) error {
	if view, ok := m.Load(name); ok {
		if view.FileInfo.Handler != nil {
			if err := container.Close(view.FileInfo.Handler); err != nil {
				return err
			}
		}
		m.Delete(name)
	}
	return nil
}

func (m ViewMap) Clean(container *file.Container) error {
	keys := m.Keys()
	for _, k := range keys {
		if err := m.Dispose(container, k); err != nil {
			return err
		}
	}
	return nil
}

func (m ViewMap) CleanWithErrors(container *file.Container) error {
	keys := m.Keys()
	var errs []error
	for _, k := range keys {
		if view, ok := m.Load(k); ok {
			if err := container.CloseWithErrors(view.FileInfo.Handler); err != nil {
				errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
			}
			m.Delete(k)
		}
	}

	return file.NewForcedUnlockError(errs)
}
