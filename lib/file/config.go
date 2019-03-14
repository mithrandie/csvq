package file

import "time"

const DefaultWaitTimeout = 10 * time.Second

var WaitTimeout = DefaultWaitTimeout
var RetryDelay = 10 * time.Millisecond

const (
	LockFileSuffix = ".lock"
	TempFileSuffix = ".temp"
)
