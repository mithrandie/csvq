package output

import "fmt"

func Write(file string, s string) {
	if len(file) < 1 {
		writeStdout(s)
	}
}

func writeStdout(s string) {
	fmt.Print(s)
}
