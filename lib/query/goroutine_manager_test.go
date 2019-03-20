package query

import (
	"testing"
)

var goroutineManagerAssignRoutineNumberTests = []struct {
	Name                          string
	PresetCPU                     int
	DefaultMinimumRequiredPerCore int
	PresetCount                   int
	RecordLen                     int
	MinimumRequired               int
	Expect                        int
	ExpectCount                   int
}{
	{
		Name:                          "First Assign",
		PresetCPU:                     4,
		DefaultMinimumRequiredPerCore: 40,
		PresetCount:                   0,
		RecordLen:                     10000,
		MinimumRequired:               -1,
		Expect:                        4,
		ExpectCount:                   3,
	},
	{
		Name:                          "Already Assigned",
		PresetCPU:                     4,
		DefaultMinimumRequiredPerCore: 40,
		PresetCount:                   4,
		RecordLen:                     10000,
		MinimumRequired:               -1,
		Expect:                        1,
		ExpectCount:                   4,
	},
	{
		Name:                          "Use One CPU Core",
		PresetCPU:                     4,
		DefaultMinimumRequiredPerCore: 40,
		PresetCount:                   0,
		RecordLen:                     10,
		MinimumRequired:               -1,
		Expect:                        1,
		ExpectCount:                   0,
	},
	{
		Name:                          "Use Two CPU Cores",
		PresetCPU:                     4,
		DefaultMinimumRequiredPerCore: 40,
		PresetCount:                   0,
		RecordLen:                     90,
		MinimumRequired:               -1,
		Expect:                        2,
		ExpectCount:                   1,
	},
}

func TestGoroutineManager_AssignRoutineNumber(t *testing.T) {
	gm := GetGoroutineManager()

	for _, v := range goroutineManagerAssignRoutineNumberTests {
		gm.Count = v.PresetCount
		gm.MinimumRequiredPerCore = v.DefaultMinimumRequiredPerCore

		result := gm.AssignRoutineNumber(v.RecordLen, v.MinimumRequired, v.PresetCPU)
		if result != v.Expect {
			t.Errorf("%s: result = %d, want %d", v.Name, result, v.Expect)
		}
		if gm.Count != v.ExpectCount {
			t.Errorf("%s: count = %d, want %d", v.Name, gm.Count, v.ExpectCount)
		}
	}

	gm.Count = 0
	gm.MinimumRequiredPerCore = MinimumRequiredPerCPUCore
}

var goroutineTaskManagerRecordRangeTests = []struct {
	Name            string
	NumberOfRoutine int
	LecordLen       int
	RoutineIndex    int
	ExpectStart     int
	ExpectEnd       int
}{
	{
		Name:            "First Routine",
		NumberOfRoutine: 4,
		LecordLen:       999,
		RoutineIndex:    0,
		ExpectStart:     0,
		ExpectEnd:       249,
	},
	{
		Name:            "Second Routine",
		NumberOfRoutine: 4,
		LecordLen:       999,
		RoutineIndex:    1,
		ExpectStart:     249,
		ExpectEnd:       498,
	},
	{
		Name:            "Third Routine",
		NumberOfRoutine: 4,
		LecordLen:       999,
		RoutineIndex:    2,
		ExpectStart:     498,
		ExpectEnd:       747,
	},
	{
		Name:            "Fourth Routine",
		NumberOfRoutine: 4,
		LecordLen:       999,
		RoutineIndex:    3,
		ExpectStart:     747,
		ExpectEnd:       999,
	},
	{
		Name:            "Fifth Routine",
		NumberOfRoutine: 4,
		LecordLen:       999,
		RoutineIndex:    5,
		ExpectStart:     0,
		ExpectEnd:       0,
	},
}

func TestGoroutineTaskManager_RecordRange(t *testing.T) {
	gtm := &GoroutineTaskManager{}

	for _, v := range goroutineTaskManagerRecordRangeTests {
		gtm.Number = v.NumberOfRoutine
		gtm.recordLen = v.LecordLen

		start, end := gtm.RecordRange(v.RoutineIndex)
		if start != v.ExpectStart || end != v.ExpectEnd {
			t.Errorf("%s: result = %d, %d, want %d, %d", v.Name, start, end, v.ExpectStart, v.ExpectEnd)
		}
	}
}
