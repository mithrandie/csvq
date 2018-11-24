package query

import "strings"

type UncommittedViewMap struct {
	Created map[string]*FileInfo
	Updated map[string]*FileInfo
}

func NewUncommittedViewMap() *UncommittedViewMap {
	return &UncommittedViewMap{
		Created: make(map[string]*FileInfo),
		Updated: make(map[string]*FileInfo),
	}
}

func (m *UncommittedViewMap) SetForCreatedView(fileInfo *FileInfo) {
	ufpath := strings.ToUpper(fileInfo.Path)

	if _, ok := m.Created[ufpath]; !ok {
		if _, ok := m.Updated[ufpath]; !ok {
			m.Created[ufpath] = fileInfo
		}
	}
}

func (m *UncommittedViewMap) SetForUpdatedView(fileInfo *FileInfo) {
	ufpath := strings.ToUpper(fileInfo.Path)

	if _, ok := m.Created[ufpath]; !ok {
		if _, ok := m.Updated[ufpath]; !ok {
			m.Updated[ufpath] = fileInfo
		}
	}
}

func (m *UncommittedViewMap) Unset(fileInfo *FileInfo) {
	ufpath := strings.ToUpper(fileInfo.Path)

	if _, ok := m.Updated[ufpath]; ok {
		delete(m.Updated, ufpath)
		return
	}

	if _, ok := m.Created[ufpath]; ok {
		delete(m.Created, ufpath)
	}
}

func (m *UncommittedViewMap) Clean() {
	for k := range m.Updated {
		delete(m.Updated, k)
	}
	for k := range m.Created {
		delete(m.Created, k)
	}
}

func (m *UncommittedViewMap) UncommittedFiles() (map[string]*FileInfo, map[string]*FileInfo) {
	var createdFiles = make(map[string]*FileInfo)
	var updatedFiles = make(map[string]*FileInfo)

	for k, v := range m.Created {
		if !v.IsTemporary {
			createdFiles[k] = v
		}
	}
	for k, v := range m.Updated {
		if !v.IsTemporary {
			updatedFiles[k] = v
		}
	}

	return createdFiles, updatedFiles
}

func (m *UncommittedViewMap) UncommittedTempViews() map[string]*FileInfo {
	var updatedViews = map[string]*FileInfo{}

	for k, v := range m.Updated {
		if v.IsTemporary {
			updatedViews[k] = v
		}
	}

	return updatedViews
}

func (m *UncommittedViewMap) IsEmpty() bool {
	if 0 < len(m.Created) {
		return false
	}
	if 0 < len(m.Updated) {
		return false
	}
	return true
}

func (m *UncommittedViewMap) CountCreatedTables() int {
	return len(m.Created)
}

func (m *UncommittedViewMap) CountUpdatedTables() int {
	cnt := 0
	for _, v := range m.Updated {
		if !v.IsTemporary {
			cnt++
		}
	}
	return cnt
}

func (m *UncommittedViewMap) CountUpdatedViews() int {
	cnt := 0
	for _, v := range m.Updated {
		if v.IsTemporary {
			cnt++
		}
	}
	return cnt
}
