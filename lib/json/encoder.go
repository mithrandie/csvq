package json

import (
	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/color"
	"strings"
)

const IndentSpaces = 2

type Encoder struct {
	EscapeType  EscapeType
	PrettyPrint bool
	LineBreak   cmd.LineBreak

	nameSeparator string
	lineBreak     string

	decoder *Decoder
}

func NewEncoder() *Encoder {
	return &Encoder{
		EscapeType:    Backslash,
		PrettyPrint:   false,
		LineBreak:     cmd.LF,
		nameSeparator: string(NameSeparator),
		decoder:       NewDecoder(),
	}
}

func (e *Encoder) Encode(structure Structure) string {
	if e.PrettyPrint {
		e.lineBreak = e.LineBreak.Value()
		e.nameSeparator = string(NameSeparator) + " "
	} else {
		e.lineBreak = ""
		e.nameSeparator = string(NameSeparator)
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
					color.BlueB(Quote(e.escape(member.Key)))+
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
		encoded = color.Magenta(structure.String())
	case String:
		str := string(structure.(String))
		decoded, err := e.decoder.Decode(str)
		if err == nil {
			encoded = e.encodeStructure(decoded, depth)
		} else {
			encoded = color.Green(Quote(e.escape(str)))
		}
	case Boolean:
		encoded = color.YellowB(structure.String())
	case Null:
		encoded = color.BrightBlack(structure.String())
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
