package query

import (
	"os"
	"path"
	"testing"

	"github.com/mithrandie/csvq/lib/cmd"
)

func TestMain(m *testing.M) {
	setup()
	r := m.Run()
	teardown()
	os.Exit(r)
}

func setup() {
	flags := cmd.GetFlags()
	flags.Location = "America/Los_Angeles"
	flags.Now = "2012-02-03 09:18:15"

	dir, _ := os.Getwd()
	r, _ := os.Open(path.Join(dir, "..", "..", "testdata", "csv", "empty.txt"))
	os.Stdin = r
}

func teardown() {}
