package color

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	defer teardown()

	setup()
	return m.Run()
}

func setup() {
	UseEscapeSequences = true
}

func teardown() {
	UseEscapeSequences = false
}
