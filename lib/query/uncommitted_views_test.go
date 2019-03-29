package query

import (
	"reflect"
	"sync"
	"testing"
)

var preCreatedFileInfo = &FileInfo{Path: "pre_created.txt"}
var preUpdatedFileInfo = &FileInfo{Path: "pre_updated.txt"}

func TestUncommittedViewMap_SetForCreatedView(t *testing.T) {
	m := &UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"PRE_CREATED.TXT": {Path: "pre_created.txt"},
		},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
		},
	}

	info := &FileInfo{
		Path: "create.txt",
	}
	expect := &UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"PRE_CREATED.TXT": {Path: "pre_created.txt"},
			"CREATE.TXT":      {Path: "create.txt"},
		},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
		},
	}
	m.SetForCreatedView(info)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}

	m.SetForCreatedView(preCreatedFileInfo)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}

	m.SetForCreatedView(preUpdatedFileInfo)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}
}

func TestUncommittedViewMap_SetForUpdatedView(t *testing.T) {
	m := &UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"PRE_CREATED.TXT": {Path: "pre_created.txt"},
		},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
		},
	}

	info := &FileInfo{
		Path: "update.txt",
	}
	expect := &UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"PRE_CREATED.TXT": {Path: "pre_created.txt"},
		},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
			"UPDATE.TXT":      {Path: "update.txt"},
		},
	}
	m.SetForUpdatedView(info)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}

	m.SetForUpdatedView(preCreatedFileInfo)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}

	m.SetForUpdatedView(preUpdatedFileInfo)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}
}

func TestUncommittedViewMap_Unset(t *testing.T) {
	m := &UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"PRE_CREATED.TXT": {Path: "pre_created.txt"},
		},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
		},
	}

	expect := &UncommittedViews{
		mtx:     &sync.RWMutex{},
		Created: map[string]*FileInfo{},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
		},
	}

	m.Unset(preCreatedFileInfo)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}

	expect = &UncommittedViews{
		mtx:     &sync.RWMutex{},
		Created: map[string]*FileInfo{},
		Updated: map[string]*FileInfo{},
	}
	m.Unset(preUpdatedFileInfo)
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("map = %v, want %v", m, expect)
	}
}

func TestUncommittedViewMap_IsEmpty(t *testing.T) {
	var expect bool

	m := &UncommittedViews{
		mtx: &sync.RWMutex{},
		Created: map[string]*FileInfo{
			"PRE_CREATED.TXT": {Path: "pre_created.txt"},
		},
	}
	expect = false
	result := m.IsEmpty()
	if result != expect {
		t.Errorf("result = %t, want %t", result, expect)
	}

	m = &UncommittedViews{
		mtx: &sync.RWMutex{},
		Updated: map[string]*FileInfo{
			"PRE_UPDATED.TXT": {Path: "pre_updated.txt"},
		},
	}
	expect = false
	result = m.IsEmpty()
	if result != expect {
		t.Errorf("result = %t, want %t", result, expect)
	}

	m = &UncommittedViews{
		mtx: &sync.RWMutex{},
	}
	expect = true
	result = m.IsEmpty()
	if result != expect {
		t.Errorf("result = %t, want %t", result, expect)
	}
}
