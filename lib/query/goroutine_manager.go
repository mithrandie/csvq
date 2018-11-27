package query

import (
	"sync"

	"github.com/mithrandie/csvq/lib/cmd"
)

var (
	gm    *GoroutineManager
	getGm sync.Once
)

const MinimumRequiredForParallelRoutine = 150

func GetGoroutineManager() *GoroutineManager {
	getGm.Do(func() {
		gm = &GoroutineManager{
			Count:           0,
			CountMutex:      new(sync.Mutex),
			MinimumRequired: MinimumRequiredForParallelRoutine,
		}
	})
	return gm
}

type GoroutineManager struct {
	Count           int
	CountMutex      *sync.Mutex
	MinimumRequired int
}

func (m *GoroutineManager) AssignRoutineNumber(recordLen int, minimumRequired int) int {
	number := cmd.GetFlags().CPU
	if minimumRequired < 0 {
		minimumRequired = m.MinimumRequired
	}

	m.CountMutex.Lock()
	defer m.CountMutex.Unlock()

	max := number - m.Count
	if max < 1 {
		max = 1
	}
	if max < number {
		number = max
	}
	if 1 < number && recordLen < minimumRequired {
		number = 1
	}

	if 0 < number {
		m.Count += number - 1
	}
	return number
}

func (m *GoroutineManager) Release() {
	m.CountMutex.Lock()
	if 0 < m.Count {
		m.Count--
	}
	m.CountMutex.Unlock()
}

type GoroutineTaskManager struct {
	Number int

	grCountMutex sync.Mutex
	grCount      int
	recordLen    int
	waitGroup    sync.WaitGroup
	err          error
}

func NewGoroutineTaskManager(recordLen int, minimumRequired int) *GoroutineTaskManager {
	number := GetGoroutineManager().AssignRoutineNumber(recordLen, minimumRequired)

	return &GoroutineTaskManager{
		Number:    number,
		grCount:   number - 1,
		recordLen: recordLen,
	}
}

func (m *GoroutineTaskManager) HasError() bool {
	return m.err != nil
}

func (m *GoroutineTaskManager) SetError(e error) {
	m.err = e
}

func (m *GoroutineTaskManager) Err() error {
	return m.err
}

func (m *GoroutineTaskManager) RecordRange(routineIndex int) (int, int) {
	calcLen := m.recordLen / m.Number

	var start = routineIndex * calcLen

	if m.recordLen <= start {
		return 0, 0
	}

	var end int
	if routineIndex == m.Number-1 {
		end = m.recordLen
	} else {
		end = (routineIndex + 1) * calcLen
	}
	return start, end
}

func (m *GoroutineTaskManager) Add() {
	m.waitGroup.Add(1)
}

func (m *GoroutineTaskManager) Done() {
	m.grCountMutex.Lock()
	if 0 < m.grCount {
		m.grCount--
		GetGoroutineManager().Release()
	}
	m.grCountMutex.Unlock()

	m.waitGroup.Done()
}

func (m *GoroutineTaskManager) Wait() {
	m.waitGroup.Wait()
}

func (m *GoroutineTaskManager) Run(fn func(int)) {
	for i := 0; i < m.Number; i++ {
		m.Add()
		go func(thIdx int) {
			start, end := m.RecordRange(thIdx)

			for j := start; j < end; j++ {
				fn(j)
			}

			m.Done()
		}(i)
	}
	m.Wait()
}
