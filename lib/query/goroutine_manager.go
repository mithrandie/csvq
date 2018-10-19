package query

import (
	"sync"

	"github.com/mithrandie/csvq/lib/cmd"
)

var goroutineCount int
var goroutineCountMutex sync.Mutex

const MinimumRequiredForParallelRoutine = 150

func useParallelRoutine(recordLen int, minimumRequired int) (int, int) {
	cpu := cmd.GetFlags().CPU

	goroutineCountMutex.Lock()
	defer goroutineCountMutex.Unlock()

	max := cpu - goroutineCount
	if max < 1 {
		max = 1
	}
	if max < cpu {
		cpu = max
	}
	if 1 < cpu && (recordLen < minimumRequired || recordLen < cpu) {
		cpu = 1
	}

	addGRCount := cpu - 1
	if 0 < addGRCount {
		goroutineCount += addGRCount
	}
	return cpu, addGRCount
}

func decreamentGoroutineCount() {
	goroutineCountMutex.Lock()
	goroutineCount--
	goroutineCountMutex.Unlock()
}

type GoroutineManager struct {
	CPU int

	grCountMutex sync.Mutex
	grCount      int
	recordLen    int
	waitGroup    sync.WaitGroup
	errMutex     sync.Mutex
	err          error
}

func NewGoroutineManager(recordLen int, minimumRequired int) *GoroutineManager {
	cpu, grCount := useParallelRoutine(recordLen, minimumRequired)

	return &GoroutineManager{
		CPU:       cpu,
		grCount:   grCount,
		recordLen: recordLen,
	}
}

func (m *GoroutineManager) HasError() bool {
	return m.err != nil
}

func (m *GoroutineManager) SetError(e error) {
	m.errMutex.Lock()
	m.err = e
	m.errMutex.Unlock()
}

func (m *GoroutineManager) Error() error {
	return m.err
}

func (m *GoroutineManager) RecordRange(routineIndex int) (int, int) {
	return RecordRange(routineIndex, m.recordLen, m.CPU)
}

func (m *GoroutineManager) Add() {
	m.waitGroup.Add(1)
}

func (m *GoroutineManager) Done() {
	m.grCountMutex.Lock()
	if 0 < m.grCount {
		m.grCount--
		decreamentGoroutineCount()
	}
	m.grCountMutex.Unlock()

	m.waitGroup.Done()
}

func (m *GoroutineManager) Wait() {
	m.waitGroup.Wait()
}
