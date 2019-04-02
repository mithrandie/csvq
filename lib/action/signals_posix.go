// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package action

import (
	"os"
	"syscall"
)

var Signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
}
