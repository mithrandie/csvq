package file

import "time"

const DefaultWaitTimeout = 10 * time.Second
const DefaultRetryDelay = 10 * time.Millisecond

const (
	LockFileSuffix = ".lock"
	TempFileSuffix = ".temp"
)
