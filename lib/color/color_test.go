package color

import (
	"testing"
)

func TestStripEscapeSequence(t *testing.T) {
	s := "\033[34ma\033[0m"
	expect := "a"
	result := StripEscapeSequence(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}

	s = "\033[34;1ma\033[0m"
	expect = "a"
	result = StripEscapeSequence(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestColorize(t *testing.T) {
	s := "a"
	expect := "\033[34ma\033[0m"
	result := Colorize(s, FGBlue, false)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}

	s = "a"
	expect = "\033[34;1ma\033[0m"
	result = Colorize(s, FGBlue, true)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}

	UseEscapeSequences = false
	s = "a"
	expect = "a"
	result = Colorize(s, FGBlue, true)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}

	UseEscapeSequences = true
}

func TestBlack(t *testing.T) {
	s := "a"
	expect := "\033[30ma\033[0m"
	result := Black(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestRed(t *testing.T) {
	s := "a"
	expect := "\033[31ma\033[0m"
	result := Red(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestGreen(t *testing.T) {
	s := "a"
	expect := "\033[32ma\033[0m"
	result := Green(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestYellow(t *testing.T) {
	s := "a"
	expect := "\033[33ma\033[0m"
	result := Yellow(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBlue(t *testing.T) {
	s := "a"
	expect := "\033[34ma\033[0m"
	result := Blue(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestMagenta(t *testing.T) {
	s := "a"
	expect := "\033[35ma\033[0m"
	result := Magenta(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestCyan(t *testing.T) {
	s := "a"
	expect := "\033[36ma\033[0m"
	result := Cyan(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestWhite(t *testing.T) {
	s := "a"
	expect := "\033[37ma\033[0m"
	result := White(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightBlack(t *testing.T) {
	s := "a"
	expect := "\033[90ma\033[0m"
	result := BrightBlack(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightRed(t *testing.T) {
	s := "a"
	expect := "\033[91ma\033[0m"
	result := BrightRed(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightGreen(t *testing.T) {
	s := "a"
	expect := "\033[92ma\033[0m"
	result := BrightGreen(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightYellow(t *testing.T) {
	s := "a"
	expect := "\033[93ma\033[0m"
	result := BrightYellow(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightBlue(t *testing.T) {
	s := "a"
	expect := "\033[94ma\033[0m"
	result := BrightBlue(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightMagenta(t *testing.T) {
	s := "a"
	expect := "\033[95ma\033[0m"
	result := BrightMagenta(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightCyan(t *testing.T) {
	s := "a"
	expect := "\033[96ma\033[0m"
	result := BrightCyan(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightWhite(t *testing.T) {
	s := "a"
	expect := "\033[97ma\033[0m"
	result := BrightWhite(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBlackB(t *testing.T) {
	s := "a"
	expect := "\033[30;1ma\033[0m"
	result := BlackB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestRedB(t *testing.T) {
	s := "a"
	expect := "\033[31;1ma\033[0m"
	result := RedB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestGreenB(t *testing.T) {
	s := "a"
	expect := "\033[32;1ma\033[0m"
	result := GreenB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestYellowB(t *testing.T) {
	s := "a"
	expect := "\033[33;1ma\033[0m"
	result := YellowB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBlueB(t *testing.T) {
	s := "a"
	expect := "\033[34;1ma\033[0m"
	result := BlueB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestMagentaB(t *testing.T) {
	s := "a"
	expect := "\033[35;1ma\033[0m"
	result := MagentaB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestCyanB(t *testing.T) {
	s := "a"
	expect := "\033[36;1ma\033[0m"
	result := CyanB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestWhiteB(t *testing.T) {
	s := "a"
	expect := "\033[37;1ma\033[0m"
	result := WhiteB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightBlackB(t *testing.T) {
	s := "a"
	expect := "\033[90;1ma\033[0m"
	result := BrightBlackB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightRedB(t *testing.T) {
	s := "a"
	expect := "\033[91;1ma\033[0m"
	result := BrightRedB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightGreenB(t *testing.T) {
	s := "a"
	expect := "\033[92;1ma\033[0m"
	result := BrightGreenB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightYellowB(t *testing.T) {
	s := "a"
	expect := "\033[93;1ma\033[0m"
	result := BrightYellowB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightBlueB(t *testing.T) {
	s := "a"
	expect := "\033[94;1ma\033[0m"
	result := BrightBlueB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightMagentaB(t *testing.T) {
	s := "a"
	expect := "\033[95;1ma\033[0m"
	result := BrightMagentaB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightCyanB(t *testing.T) {
	s := "a"
	expect := "\033[96;1ma\033[0m"
	result := BrightCyanB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestBrightWhiteB(t *testing.T) {
	s := "a"
	expect := "\033[97;1ma\033[0m"
	result := BrightWhiteB(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestError(t *testing.T) {
	s := "a"
	expect := RedB(s)
	result := Error(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestWarn(t *testing.T) {
	s := "a"
	expect := YellowB(s)
	result := Warn(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}

func TestInfo(t *testing.T) {
	s := "a"
	expect := GreenB(s)
	result := Info(s)

	if result != expect {
		t.Errorf("result = %q, want %q for %s", result, expect, s)
	}
}
