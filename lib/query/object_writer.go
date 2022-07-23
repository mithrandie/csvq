package query

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/mithrandie/csvq/lib/option"
	"github.com/mithrandie/go-text/color"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	DefaultLineWidth = 75
	DefaultPadding   = 1
)

type ObjectWriter struct {
	Flags   *option.Flags
	Palette *color.Palette

	MaxWidth    int
	Padding     int
	Indent      int
	IndentWidth int

	Title1       string
	Title1Effect string
	Title2       string
	Title2Effect string

	buf bytes.Buffer

	subBlock  int
	lineWidth int
	Column    int
}

func NewObjectWriter(tx *Transaction) *ObjectWriter {
	maxWidth := DefaultLineWidth
	if tx.Session.Terminal() != nil {
		if termw, _, err := tx.Session.Terminal().GetSize(); err == nil {
			maxWidth = termw
		}
	} else {
		if w, _, err := terminal.GetSize(int(tx.Session.ScreenFd())); err == nil {
			maxWidth = w
		}
	}

	return &ObjectWriter{
		Flags:       tx.Flags,
		Palette:     tx.Palette,
		MaxWidth:    maxWidth,
		Indent:      0,
		IndentWidth: 4,
		Padding:     DefaultPadding,
		lineWidth:   0,
		Column:      0,
		subBlock:    0,
	}
}

func (w *ObjectWriter) Clear() {
	w.Title1 = ""
	w.Title1Effect = ""
	w.Title2 = ""
	w.Title2Effect = ""
	w.lineWidth = 0
	w.Column = 0
	w.subBlock = 0
	w.buf.Reset()
}

func (w *ObjectWriter) WriteColorWithoutLineBreak(s string, effect string) {
	w.write(s, effect, true)
}

func (w *ObjectWriter) WriteColor(s string, effect string) {
	w.write(s, effect, false)
}

func (w *ObjectWriter) write(s string, effect string, withoutLineBreak bool) {
	startOfLine := w.Column < 1

	if startOfLine {
		width := w.LeadingSpacesWidth() + w.subBlock
		w.writeToBuf(strings.Repeat(" ", width))
		w.Column = width
	}

	if !withoutLineBreak && !startOfLine && !w.FitInLine(s) {
		w.NewLine()
		w.write(s, effect, withoutLineBreak)
	} else {
		if w.Palette == nil {
			w.writeToBuf(s)
		} else {
			w.writeToBuf(w.Palette.Render(effect, s))
		}
		w.Column = w.Column + option.TextWidth(s, w.Flags)
	}
}

func (w *ObjectWriter) writeToBuf(s string) {
	w.buf.WriteString(s)
}

func (w *ObjectWriter) LeadingSpacesWidth() int {
	return w.Padding + (w.Indent * w.IndentWidth)
}

func (w *ObjectWriter) FitInLine(s string) bool {
	if w.MaxWidth-(w.Padding*2)-1 < w.Column+option.TextWidth(s, w.Flags) {
		return false
	}
	return true
}

func (w *ObjectWriter) WriteWithoutLineBreak(s string) {
	w.WriteColorWithoutLineBreak(s, option.NoEffect)
}

func (w *ObjectWriter) Write(s string) {
	w.WriteColor(s, option.NoEffect)
}

func (w *ObjectWriter) WriteWithAutoLineBreak(s string) {
	w.writeWithAutoLineBreak(s, false, true)
}

func (w *ObjectWriter) WriteWithAutoLineBreakWithContinueMark(s string) {
	w.writeWithAutoLineBreak(s, true, false)
}

func (w *ObjectWriter) writeWithAutoLineBreak(s string, useContinueMark bool, useBlock bool) {
	continueMark := ""
	if useContinueMark {
		continueMark = "\\"
	}

	scanner := bufio.NewScanner(strings.NewReader(s))
	firstLine := true
	blockQuote := false
	preformatted := false
	for scanner.Scan() {
		if blockQuote {
			w.EndBlock()
			blockQuote = false
		}

		line := scanner.Text()
		if useBlock && option.TrimSpace(line) == "```" {
			preformatted = !preformatted
			continue
		} else {
			if firstLine {
				firstLine = false
			} else {
				w.NewLine()
			}
		}

		if preformatted {
			w.Write(line)
			continue
		}

		wscanner := bufio.NewScanner(strings.NewReader(line))
		wscanner.Split(bufio.ScanWords)
		lineHead := true

		for wscanner.Scan() {
			word := wscanner.Text()
			if lineHead {
				if useBlock && blockQuote == false && word == ">" {
					blockQuote = true
					w.BeginBlock()
					continue
				}

				lineHead = false
			} else {
				if !w.FitInLine(" " + word + continueMark) {
					w.Write(continueMark)
					w.NewLine()
				} else {
					word = " " + word
				}
			}

			w.Write(word)
		}
	}

	if blockQuote {
		w.EndBlock()
	}
}

func (w *ObjectWriter) WriteSpaces(l int) {
	w.Write(strings.Repeat(" ", l))
}

func (w *ObjectWriter) NewLine() {
	w.buf.WriteRune('\n')
	if w.lineWidth < w.Column {
		w.lineWidth = w.Column
	}
	w.Column = 0
}

func (w *ObjectWriter) BeginBlock() {
	w.Indent++
}

func (w *ObjectWriter) EndBlock() {
	w.Indent--
}

func (w *ObjectWriter) BeginSubBlock() {
	w.subBlock = w.Column - w.LeadingSpacesWidth()
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
		tw := option.TextWidth(w.Title1, w.Flags) + option.TextWidth(w.Title2, w.Flags)
		if 0 < len(w.Title1) && 0 < len(w.Title2) {
			tw++
		}

		hlLen := tw + 2
		if hlLen < w.lineWidth+1 {
			hlLen = w.lineWidth + 1
		}
		if hlLen < w.Column+1 {
			hlLen = w.Column + 1
		}
		if w.MaxWidth < hlLen {
			hlLen = w.MaxWidth
		}

		if tw < hlLen {
			header.Write(bytes.Repeat([]byte(" "), (hlLen-tw)/2))
		}
		if 0 < len(w.Title1) {
			if w.Palette == nil {
				header.WriteString(w.Title1)
			} else {
				header.WriteString(w.Palette.Render(w.Title1Effect, w.Title1))
			}
		}
		if 0 < len(w.Title2) {
			header.WriteRune(' ')
			if w.Palette == nil {
				header.WriteString(w.Title2)
			} else {
				header.WriteString(w.Palette.Render(w.Title2Effect, w.Title2))
			}
		}
		header.WriteRune('\n')
		header.Write(bytes.Repeat([]byte("-"), hlLen))
		header.WriteRune('\n')
	}

	return header.String() + w.buf.String()
}
