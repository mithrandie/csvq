package color

import (
	"testing"
)

var paletteColorTests = []struct {
	Text     string
	Style    Style
	UseStyle bool
	Expect   string
}{
	{
		Text:     "abc",
		Style:    PlainStyle,
		UseStyle: true,
		Expect:   "\033[0mabc\033[0m",
	},
	{
		Text:     "abc",
		Style:    FieldLableStyle,
		UseStyle: true,
		Expect:   "\033[34;1mabc\033[0m",
	},
	{
		Text:     "abc",
		Style:    999,
		UseStyle: true,
		Expect:   "\033[0mabc\033[0m",
	},
}

func TestPalette_Color(t *testing.T) {
	p := NewPalette()

	for _, v := range paletteColorTests {
		p.useStyle = v.UseStyle

		result := p.Color(v.Text, v.Style)
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %q, %#v", result, v.Expect, v.Text, v.Style)
		}
	}

	style := paletteColorTests[0]
	p.Disable()
	result := p.Color(style.Text, style.Style)
	expect := "abc"
	if result != "abc" {
		t.Errorf("result = %q, want %q for %q, %#v", result, expect, style.Text, style.Style)
	}

	p.Enable()
	result = p.Color(style.Text, style.Style)
	if result != style.Expect {
		t.Errorf("result = %q, want %q for %q, %#v", result, style.Expect, style.Text, style.Style)
	}
}
