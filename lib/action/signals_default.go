// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris,!windows

package action

import (
	"os"
)

var Signals = []os.Signal{
	os.Interrupt,
}
