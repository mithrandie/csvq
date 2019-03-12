package file

import "time"

var WaitTimeout = 10 * time.Second
var RetryDelay = 10 * time.Millisecond

const (
	LockFileSuffix = ".lock"
	TempFileSuffix = ".temp"
)
