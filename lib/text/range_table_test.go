package text

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"testing"
)

func TestStringWidth(t *testing.T) {
	s := "日本語abc"
	expect := 9

	result := StringWidth(s)
	if result != expect {
		t.Errorf("string width = %d, want %d", result, expect)
	}

	s = "日本語\033[33mab\033[0mc"
	expect = 9

	result = StringWidth(s)
	if result != expect {
		t.Errorf("string width = %d, want %d", result, expect)
	}
}

func TestRuneByteSize(t *testing.T) {
	r := '日'
	expect := 3

	result := RuneByteSize(r, cmd.UTF8)
	if result != expect {
		t.Errorf("byte size = %d, want %d for %s in %s", result, expect, string(r), cmd.UTF8)
	}

	r = 'a'
	expect = 1

	result = RuneByteSize(r, cmd.UTF8)
	if result != expect {
		t.Errorf("byte size = %d, want %d for %s in %s", result, expect, string(r), cmd.UTF8)
	}

	r = '日'
	expect = 2

	result = RuneByteSize(r, cmd.SJIS)
	if result != expect {
		t.Errorf("byte size = %d, want %d for %s in %s", result, expect, string(r), cmd.SJIS)
	}

	r = 'a'
	expect = 1

	result = RuneByteSize(r, cmd.SJIS)
	if result != expect {
		t.Errorf("byte size = %d, want %d for %s in %s", result, expect, string(r), cmd.SJIS)
	}
}

func TestByteSize(t *testing.T) {
	s := "日本語abc"
	expect := 12
	result := ByteSize(s, cmd.UTF8)
	if result != expect {
		t.Errorf("byte size = %d, want %d for %s in %s", result, expect, s, cmd.UTF8)
	}

	expect = 9
	result = ByteSize(s, cmd.SJIS)
	if result != expect {
		t.Errorf("byte size = %d, want %d for %s in %s", result, expect, s, cmd.SJIS)
	}
}
