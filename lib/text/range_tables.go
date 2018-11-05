package text

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"unicode"
)

var fullWidthTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x1100, 0x11ff, 1}, //Hangul Jamo
		{0x2460, 0x24ff, 1}, //Enclosed Alphanumerics
		{0x2500, 0x257f, 1}, //Box Drawing
		{0x2e80, 0x2eff, 1}, //CJK Radicals Supplement
		{0x2f00, 0x2fdf, 1}, //CJK Radicals
		{0x2ff0, 0x2fff, 1}, //Ideographic Description Characters
		{0x3000, 0x303e, 1}, //CJK Symbols and Punctuation
		{0x3040, 0x309f, 1}, //Hiragana
		{0x30a0, 0x30ff, 1}, //Katakana
		{0x3100, 0x312f, 1}, //Bopomofo
		{0x3130, 0x318f, 1}, //Hangul Compatibility Jamo
		{0x3190, 0x319f, 1}, //Ideographic Annotations
		{0x31a0, 0x31bf, 1}, //Bopomofo Extended
		{0x31c0, 0x31ef, 1}, //CJK Strokes
		{0x31f0, 0x31ff, 1}, //Phonetic extensions for Ainu
		{0x3200, 0x32ff, 1}, //Enclosed CJK Letters and Months
		{0x3300, 0x33ff, 1}, //CJK Compatibility
		{0x3400, 0x4dbf, 1}, //CJK Unified Ideographs Extension A
		{0x4e00, 0x9fff, 1}, //CJK Unified Ideographs
		{0xa960, 0xa97f, 1}, //Hangul Jamo Extended A
		{0xac00, 0xd7af, 1}, //Hangul Syllables
		{0xd7b0, 0xd7ff, 1}, //Hangul Jamo Extended B
		{0xf900, 0xfaff, 1}, //CJK Compatibility Ideographs
		{0xff01, 0xff60, 1}, //FullWidth ASCII variants
		{0xffe0, 0xffe6, 1}, //FullWidth Symbol variants
	},
	R32: []unicode.Range32{
		{0x1b000, 0x1b0ff, 1}, //Historic Kana
		{0x1f100, 0x1f1ff, 1}, //Enclosed Alphanumeric Supplement
		{0x1f200, 0x1f2ff, 1}, //Enclosed Ideographic Supplement
		{0x20000, 0x2a6df, 1}, //CJK Unified Ideographs Extension B
		{0x2a700, 0x2b73f, 1}, //CJK Unified Ideographs Extension C
		{0x2b740, 0x2b81f, 1}, //CJK Unified Ideographs Extension D
		{0x2b820, 0x2ceaf, 1}, //CJK Unified Ideographs Extension E
		{0x2f800, 0x2fa1f, 1}, //CJK Compatibility Ideographs Supplement
		{0xe0100, 0xe01ef, 1}, //Variation Selectors Supplement
	},
}

var zeroWidthTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0300, 0x036f, 1}, //Combining Diacritical Marks
		{0x0591, 0x05af, 1}, //Hebrew Cantillation Marks
		{0x05b0, 0x05bd, 1}, //Hebrew Points
		{0x05bf, 0x05bf, 1}, //Hebrew Points
		{0x05c1, 0x05c2, 1}, //Hebrew Points
		{0x05c4, 0x05c5, 1}, //Hebrew Points
		{0x05c7, 0x05c7, 1}, //Hebrew Points
		{0x064b, 0x0652, 1}, //Arabic Tashkil from ISO 8859-6
		{0x0653, 0x065f, 1}, //Arabic Combining Marks
		{0x0670, 0x0670, 1}, //Arabic Tashkil
		{0x08a0, 0x08ff, 1}, //Arabic Extended-A
		{0x2028, 0x202f, 1}, //Format Characters
		{0xfbb2, 0xfbc1, 1}, //Arabic pedagogical symbols
		{0xfeff, 0xfeff, 1}, //Arabic Zero Width No-Break Space
	},
}

var rightToLeftTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0590, 0x05ff, 1}, //Hebrew
		{0x0600, 0x06ff, 1}, //Arabic
		{0x0700, 0x074f, 1}, //Syriac
		{0x0750, 0x077f, 1}, //Arabic Supplement
		{0x0860, 0x086f, 1}, //Syriac Supplement
		{0x08a0, 0x08ff, 1}, //Arabic Extended-A
		{0x200f, 0x200f, 1}, //Right-To-Left Mark
		{0x202b, 0x202b, 1}, //Right-To-Left Embedding
		{0x202e, 0x202e, 1}, //Right-To-Left Override
		{0xfb50, 0xfdff, 1}, //Arabic Presentation Forms-A
		{0xfe70, 0xfeff, 1}, //Arabic Presentation Forms-B
	},
	R32: []unicode.Range32{
		{0x1ee00, 0x1eeff, 1}, //Arabic Mathematical Alphabetic Symbols
	},
}

var sjisSingleByteTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0020, 0x007e, 1}, //ASCII
		{0x00a5, 0x00a5, 1},
		{0x203e, 0x203e, 1},
		{0xff61, 0xff9f, 1}, //Half Width Katakana
	},
}

func StringWidth(s string) int {
	l := 0

	inEscSeq := false // Ignore ANSI Escape Sequence
	for _, r := range s {
		if inEscSeq {
			if unicode.IsLetter(r) {
				inEscSeq = false
			}
		} else if r == 27 {
			inEscSeq = true
		} else {
			l = l + RuneWidth(r)
		}
	}
	return l
}

func RuneWidth(r rune) int {
	switch {
	case unicode.In(r, fullWidthTable):
		return 2
	case unicode.In(r, zeroWidthTable) || unicode.IsControl(r):
		return 0
	}
	return 1
}

func RuneByteSize(r rune, encoding cmd.Encoding) int {
	switch encoding {
	case cmd.SJIS:
		return SJISRuneByteSize(r)
	default:
		return len(string(r))
	}
}

func SJISRuneByteSize(r rune) int {
	switch {
	case unicode.In(r, sjisSingleByteTable) || unicode.IsControl(r):
		return 1
	}
	return 2
}

func ByteSize(s string, encoding cmd.Encoding) int {
	size := 0
	switch encoding {
	case cmd.UTF8:
		size = len(s)
	default:
		for _, c := range s {
			size = size + RuneByteSize(c, encoding)
		}
	}
	return size
}
