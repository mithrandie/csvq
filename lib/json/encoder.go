package json

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/color"
	"strings"
)

const IndentSpaces = 2

type Style int

const (
	ObjectKeyStyle Style = iota
	NumberStyle
	StringStyle
	BooleanStyle
	NullStyle
)

type TextStyle struct {
	Color int
	Bold  bool
}

type Palette struct {
	ObjectKey TextStyle
	Number    TextStyle
	String    TextStyle
	Boolean   TextStyle
	Null      TextStyle
	Plain     TextStyle

	useStyle bool
}

func NewPalette() *Palette {
	return &Palette{
		ObjectKey: TextStyle{Color: color.FGBlue, Bold: true},
		Number:    TextStyle{Color: color.FGMagenta, Bold: false},
		String:    TextStyle{Color: color.FGGreen, Bold: false},
		Boolean:   TextStyle{Color: color.FGYellow, Bold: true},
		Null:      TextStyle{Color: color.FGBrightBlack, Bold: false},
		Plain:     TextStyle{Color: color.PlainColor, Bold: false},
		useStyle:  false,
	}
}

func (p *Palette) Enable() {
	p.useStyle = true
}

func (p *Palette) Disable() {
	p.useStyle = false
}

func (p *Palette) Color(s string, style Style) string {
	var textStyle TextStyle

	if p.useStyle {
		switch style {
		case ObjectKeyStyle:
			textStyle = p.ObjectKey
		case NumberStyle:
			textStyle = p.Number
		case StringStyle:
			textStyle = p.String
		case BooleanStyle:
			textStyle = p.Boolean
		case NullStyle:
			textStyle = p.Null
		}
	} else {
		textStyle = p.Plain
	}

	return p.color(s, textStyle)
}

func (p *Palette) color(s string, style TextStyle) string {
	return color.Colorize(s, style.Color, style.Bold)
}

type Encoder struct {
	EscapeType  EscapeType
	PrettyPrint bool
	LineBreak   cmd.LineBreak

	nameSeparator string
	lineBreak     string
	palette       *Palette

	decoder *Decoder
}

func NewEncoder() *Encoder {
	return &Encoder{
		EscapeType:    Backslash,
		PrettyPrint:   false,
		LineBreak:     cmd.LF,
		nameSeparator: string(NameSeparator),
		decoder:       NewDecoder(),
		palette:       NewPalette(),
	}
}

func (e *Encoder) Encode(structure Structure, useStyle bool) string {
	if e.PrettyPrint {
		e.lineBreak = e.LineBreak.Value()
		e.nameSeparator = string(NameSeparator) + " "
	} else {
		e.lineBreak = ""
		e.nameSeparator = string(NameSeparator)
	}
	if useStyle {
		e.palette.Enable()
	} else {
		e.palette.Disable()
	}

	return e.encodeStructure(structure, 0)
}

func (e *Encoder) encodeStructure(structure Structure, depth int) string {
	var indent string
	var elementIndent string
	if e.PrettyPrint {
		indent = strings.Repeat(" ", IndentSpaces*depth)
		elementIndent = strings.Repeat(" ", IndentSpaces*(depth+1))
	}

	var encoded string

	switch structure.(type) {
	case Object:
		obj := structure.(Object)
		strs := make([]string, 0, obj.Len())
		for _, member := range obj.Members {
			strs = append(
				strs,
				elementIndent+
					e.palette.Color(Quote(e.escape(member.Key)), ObjectKeyStyle)+
					e.nameSeparator+
					e.encodeStructure(member.Value, depth+1),
			)
		}
		encoded = string(BeginObject) +
			e.lineBreak +
			strings.Join(strs[:], string(ValueSeparator)+e.lineBreak) +
			e.lineBreak +
			indent + string(EndObject)
	case Array:
		array := structure.(Array)
		strs := make([]string, 0, len(array))
		for _, v := range array {
			strs = append(strs, elementIndent+e.encodeStructure(v, depth+1))
		}
		encoded = string(BeginArray) +
			e.lineBreak +
			strings.Join(strs[:], string(ValueSeparator)+e.lineBreak) +
			e.lineBreak +
			indent + string(EndArray)
	case Number:
		encoded = e.palette.Color(structure.String(), NumberStyle)
	case String:
		str := string(structure.(String))
		decoded, err := e.decoder.Decode(str)
		if err == nil {
			encoded = e.encodeStructure(decoded, depth)
		} else {
			encoded = e.palette.Color(Quote(e.escape(str)), StringStyle)
		}
	case Boolean:
		encoded = e.palette.Color(structure.String(), BooleanStyle)
	case Null:
		encoded = e.palette.Color(structure.String(), NullStyle)
	}

	return encoded
}

func (e *Encoder) escape(s string) string {
	var escaped string

	switch e.EscapeType {
	case AllWithHexDigits:
		escaped = EscapeAll(s)
	case HexDigits:
		escaped = EscapeWithHexDigits(s)
	default:
		escaped = Escape(s)
	}

	return escaped
}
