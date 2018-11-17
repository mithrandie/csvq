// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package cmd

import (
	"path/filepath"
	"testing"
)

var historyFilePathTests = []struct {
	Filename string
	Expect   string
	Error    string
}{
	{
		Filename: ".zsh_history",
		Expect:   filepath.Join(HomeDir, ".zsh_history"),
	},
	{
		Filename: "~/.zsh_history",
		Expect:   filepath.Join(HomeDir, ".zsh_history"),
	},
	{
		Filename: "/var/zsh_history",
		Expect:   filepath.Join(HomeDir, "/var/zsh_history"),
	},
}

func TestHistoryFilePath(t *testing.T) {
	for _, v := range historyFilePathTests {
		result, err := HistoryFilePath(v.Filename)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %q", err, v.Filename)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %q", err, v.Error, v.Filename)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %q", v.Error, v.Filename)
			continue
		}
		if result != v.Expect {
			t.Errorf("filepath = %q, want %q for %q", result, v.Expect, v.Filename)
		}
	}
}
