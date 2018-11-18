package file

import "time"

var WaitTimeout = 30.0
var RetryInterval = 50 * time.Millisecond

const (
	LockFileSuffix = ".lock"
	TempFileSuffix = ".temp"
)
