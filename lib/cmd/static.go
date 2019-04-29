package cmd

import (
	"math/rand"
	"sync"
	"time"
)

var (
	TestTime time.Time // For Tests
	location = time.Local

	random  *rand.Rand
	getRand sync.Once
)

func GetRand() *rand.Rand {
	getRand.Do(func() {
		random = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return random
}

func GetLocation() *time.Location {
	return location
}

func Now() time.Time {
	if !TestTime.IsZero() {
		return TestTime
	}
	return time.Now().In(location)
}
