package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/mithrandie/csvq/lib/cmd"

	"github.com/mithrandie/csvq/lib/file"

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

func (list TemporaryViewScopes) GetWithInternalId(ctx context.Context, name parser.Identifier, flags *cmd.Flags) (*View, error) {
	for _, m := range list {
		if view, err := m.GetWithInternalId(ctx, name, flags); err == nil {
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
		if m.Exists(view.FileInfo.Path) {
			m.Set(view)
			return
		}
	}
}

func (list TemporaryViewScopes) Dispose(name parser.QueryExpression) error {
	for _, m := range list {
		if err := m.DisposeTemporaryTable(name); err == nil {
			return nil
		}
	}
	return NewUndeclaredTemporaryTableError(name)
}

func (list TemporaryViewScopes) Store(session *Session, uncomittedViews map[string]*FileInfo) []string {
	msglist := make([]string, 0, len(uncomittedViews))
	for _, m := range list {
		m.Range(func(key, value interface{}) bool {
			if _, ok := uncomittedViews[key.(string)]; ok {
				view := value.(*View)

				if view.FileInfo.IsStdin() {
					session.updateStdinView(view.Copy())
				} else {
					view.CreateRestorePoint()
				}
				msglist = append(msglist, fmt.Sprintf("Commit: restore point of view %q is created.", view.FileInfo.Path))
			}
			return true
		})
	}
	return msglist
}

func (list TemporaryViewScopes) Restore(uncomittedViews map[string]*FileInfo) []string {
	msglist := make([]string, 0, len(uncomittedViews))
	for _, m := range list {
		m.Range(func(key, value interface{}) bool {
			if _, ok := uncomittedViews[key.(string)]; ok {
				view := value.(*View)

				if view.FileInfo.IsStdin() {
					m.Delete(view.FileInfo.Path)
				} else {
					view.Restore()
				}
				msglist = append(msglist, fmt.Sprintf("Rollback: view %q is restored.", view.FileInfo.Path))
			}
			return true
		})
	}
	return msglist
}

func (list TemporaryViewScopes) All() ViewMap {
	all := NewViewMap()

	for _, m := range list {
		m.Range(func(key, value interface{}) bool {
			if !value.(*View).FileInfo.IsFile() {
				k := key.(string)
				if !all.Exists(k) {
					all.Store(k, value.(*View))
				}
			}
			return true
		})
	}
	return all
}

type ViewMap struct {
	*SyncMap
}

func NewViewMap() ViewMap {
	return ViewMap{
		NewSyncMap(),
	}
}

func (m ViewMap) Store(fpath string, view *View) {
	m.store(strings.ToUpper(fpath), view)
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
	return nil, NewTableNotLoadedError(fpath)
}

func (m ViewMap) GetWithInternalId(ctx context.Context, fpath parser.Identifier, flags *cmd.Flags) (*View, error) {
	if view, ok := m.Load(fpath.Literal); ok {
		ret := view.Copy()

		ret.Header = MergeHeader(NewHeaderWithId(ret.Header[0].View, []string{}), ret.Header)

		if err := NewGoroutineTaskManager(ret.RecordLen(), -1, flags.CPU).Run(ctx, func(index int) error {
			ret.RecordSet[index] = append(Record{NewCell(value.NewInteger(int64(index)))}, ret.RecordSet[index]...)
			return nil
		}); err != nil {
			return nil, err
		}

		return ret, nil
	}
	return nil, NewTableNotLoadedError(fpath)
}

func (m ViewMap) Set(view *View) {
	if view.FileInfo != nil {
		m.Store(view.FileInfo.Path, view)
	}
}

func (m ViewMap) DisposeTemporaryTable(table parser.QueryExpression) error {
	var tableName string
	if e, ok := table.(parser.Stdin); ok {
		tableName = e.Stdin
	} else {
		tableName = table.(parser.Identifier).Literal
	}

	if v, ok := m.Load(tableName); ok {
		if !v.FileInfo.IsFile() {
			m.Delete(tableName)
			return nil
		} else {
			return NewUndeclaredTemporaryTableError(table)
		}
	}
	return NewUndeclaredTemporaryTableError(table)
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
