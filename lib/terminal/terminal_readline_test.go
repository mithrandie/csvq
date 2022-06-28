//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || windows

package terminal

import (
	"path/filepath"
	"testing"
)

var historyFilePathTests = []struct {
	Filename    string
	Expect      string
	JoinHomeDir bool
	Error       string
}{
	{
		Filename:    ".zsh_history",
		Expect:      ".zsh_history",
		JoinHomeDir: true,
	},
	{
		Filename:    filepath.Join("~", ".zsh_history"),
		Expect:      ".zsh_history",
		JoinHomeDir: true,
	},
	{
		Filename: filepath.Join(TestDir, "zsh_history"),
		Expect:   filepath.Join(TestDir, "zsh_history"),
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

		expect := v.Expect
		if v.JoinHomeDir {
			expect = filepath.Join(HomeDir, expect)
		}
		if result != expect {
			t.Errorf("filepath = %q, want %q for %q", result, expect, v.Filename)
		}
	}
}
