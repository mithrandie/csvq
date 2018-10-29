package text

import (
	"bytes"
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/color"
	"strings"
)

const (
	DefaultLineWidth = 75
	DefaultPadding   = 1
)

type ObjectWriter struct {
	MaxWidth    int
	Padding     int
	Indent      int
	IndentWidth int

	Title1      string
	Title1Style color.Style
	Title2      string
	Title2Style color.Style

	palette *color.Palette
	buf     bytes.Buffer

	subBlock  int
	lineWidth int
	column    int
}

func NewObjectWriter() *ObjectWriter {
	palette := color.NewPalette()
	palette.Enable()

	maxWidth := DefaultLineWidth
	if cmd.Terminal != nil {
		if termw, _, err := cmd.Terminal.GetSize(); err == nil {
			maxWidth = termw
		}
	}

	return &ObjectWriter{
		MaxWidth:    maxWidth,
		Indent:      0,
		IndentWidth: 4,
		Padding:     DefaultPadding,
		palette:     palette,
		lineWidth:   0,
		column:      0,
		subBlock:    0,
	}
}

func (w *ObjectWriter) WriteColorWithoutLineBreak(s string, style color.Style) {
	w.write(s, style, true)
}

func (w *ObjectWriter) WriteColor(s string, style color.Style) {
	w.write(s, style, false)
}

func (w *ObjectWriter) write(s string, style color.Style, withoutLineBreak bool) {
	startOfLine := w.column < 1

	if startOfLine {
		width := w.leadingSpacesWidth() + w.subBlock
		w.writeToBuf(strings.Repeat(" ", width))
		w.column = width
	}

	if !withoutLineBreak && !startOfLine && !w.FitInLine(s) {
		w.NewLine()
		w.write(s, style, withoutLineBreak)
	} else {
		w.writeToBuf(w.palette.Color(s, style))
		w.column = w.column + StringWidth(s)
	}
}

func (w *ObjectWriter) writeToBuf(s string) {
	w.buf.WriteString(s)
}

func (w *ObjectWriter) leadingSpacesWidth() int {
	return w.Padding + (w.Indent * w.IndentWidth)
}

func (w *ObjectWriter) FitInLine(s string) bool {
	if w.MaxWidth-w.Padding < w.column+StringWidth(s) {
		return false
	}
	return true
}

func (w *ObjectWriter) WriteWithoutLineBreak(s string) {
	w.WriteColorWithoutLineBreak(s, color.PlainStyle)
}

func (w *ObjectWriter) Write(s string) {
	w.WriteColor(s, color.PlainStyle)
}

func (w *ObjectWriter) WriteSpaces(l int) {
	w.Write(strings.Repeat(" ", l))
}

func (w *ObjectWriter) NewLine() {
	w.buf.WriteRune('\n')
	if w.lineWidth < w.column {
		w.lineWidth = w.column
	}
	w.column = 0
}

func (w *ObjectWriter) BeginBlock() {
	w.Indent++
}

func (w *ObjectWriter) EndBlock() {
	w.Indent--
}

func (w *ObjectWriter) BeginSubBlock() {
	w.subBlock = w.column - w.leadingSpacesWidth()
}

func (w *ObjectWriter) EndSubBlock() {
	w.subBlock = 0
}

func (w *ObjectWriter) ClearBlock() {
	w.Indent = 0
}

func (w *ObjectWriter) String() string {
	var header bytes.Buffer
	if 0 < len(w.Title1) || 0 < len(w.Title2) {
		tw := StringWidth(w.Title1) + StringWidth(w.Title2)
		if 0 < len(w.Title1) && 0 < len(w.Title2) {
			tw++
		}

		hlLen := tw + 2
		if hlLen < w.lineWidth+1 {
			hlLen = w.lineWidth + 1
		}
		if hlLen < w.column+1 {
			hlLen = w.column + 1
		}
		if w.MaxWidth < hlLen {
			hlLen = w.MaxWidth
		}

		if tw < hlLen {
			header.Write(bytes.Repeat([]byte(" "), (hlLen-tw)/2))
		}
		if 0 < len(w.Title1) {
			header.WriteString(w.palette.Color(w.Title1, w.Title1Style))
		}
		if 0 < len(w.Title2) {
			header.WriteRune(' ')
			header.WriteString(w.palette.Color(w.Title2, w.Title2Style))
		}
		header.WriteRune('\n')
		header.Write(bytes.Repeat([]byte("-"), hlLen))
		header.WriteRune('\n')
	}

	return header.String() + w.buf.String()
}
