package query

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"sync"
)

var parallelProcessing bool
var parallelProcessingMutex = new(sync.Mutex)

func useParallelRoutine(recordLen int, minimumRequired int) (int, bool) {
	parallelProcessingMutex.Lock()
	defer parallelProcessingMutex.Unlock()

	cpu := cmd.GetFlags().CPU
	if 2 < cpu {
		cpu = cpu - 1
	}
	if parallelProcessing || recordLen < minimumRequired || recordLen < cpu {
		cpu = 1
	}

	if cpu < 2 {
		return cpu, false
	}

	parallelProcessing = true
	return cpu, true
}

func closeParallelRoutine() {
	parallelProcessingMutex.Lock()
	defer parallelProcessingMutex.Unlock()

	parallelProcessing = false
}

type GoroutineManager struct {
	waitGroup   sync.WaitGroup
	recordLen   int
	cpu         int
	useParallel bool
}

func NewGoroutineManager(recordLen int, minimumRequired int) *GoroutineManager {
	cpu, useParallel := useParallelRoutine(recordLen, minimumRequired)

	return &GoroutineManager{
		waitGroup:   sync.WaitGroup{},
		recordLen:   recordLen,
		cpu:         cpu,
		useParallel: useParallel,
	}
}

func (m *GoroutineManager) CPU() int {
	return m.cpu
}

func (m *GoroutineManager) RecordRange(routineIndex int) (int, int) {
	return RecordRange(routineIndex, m.recordLen, m.cpu)
}

func (m *GoroutineManager) Add() {
	m.waitGroup.Add(1)
}

func (m *GoroutineManager) Done() {
	m.waitGroup.Done()
}

func (m *GoroutineManager) Wait() {
	m.waitGroup.Wait()
	if m.useParallel {
		closeParallelRoutine()
	}
}
