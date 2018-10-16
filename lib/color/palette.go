package color

type Style int

const (
	PlainStyle Style = iota
	FieldLableStyle
	NumberStyle
	StringStyle
	BooleanStyle
	TernaryStyle
	DatetimeStyle
	NullStyle
)

type TextStyle struct {
	Color int
	Bold  bool
}

type Palette struct {
	Plain      TextStyle
	FieldLabel TextStyle
	Number     TextStyle
	String     TextStyle
	Boolean    TextStyle
	Ternary    TextStyle
	Datetime   TextStyle
	Null       TextStyle

	useStyle bool
}

func NewPalette() *Palette {
	return &Palette{
		Plain:      TextStyle{Color: PlainColor, Bold: false},
		FieldLabel: TextStyle{Color: FGBlue, Bold: true},
		Number:     TextStyle{Color: FGMagenta, Bold: false},
		String:     TextStyle{Color: FGGreen, Bold: false},
		Boolean:    TextStyle{Color: FGYellow, Bold: true},
		Ternary:    TextStyle{Color: FGYellow, Bold: false},
		Datetime:   TextStyle{Color: FGCyan, Bold: false},
		Null:       TextStyle{Color: FGBrightBlack, Bold: false},
		useStyle:   false,
	}
}

func (p *Palette) Enable() {
	p.useStyle = true
}

func (p *Palette) Disable() {
	p.useStyle = false
}

func (p *Palette) Color(s string, style Style) string {
	if !p.useStyle {
		return s
	}

	var textStyle TextStyle

	switch style {
	case PlainColor:
		textStyle = p.Plain
	case FieldLableStyle:
		textStyle = p.FieldLabel
	case NumberStyle:
		textStyle = p.Number
	case StringStyle:
		textStyle = p.String
	case BooleanStyle:
		textStyle = p.Boolean
	case TernaryStyle:
		textStyle = p.Ternary
	case DatetimeStyle:
		textStyle = p.Datetime
	case NullStyle:
		textStyle = p.Null
	}

	return p.color(s, textStyle)
}

func (p *Palette) color(s string, style TextStyle) string {
	return Colorize(s, style.Color, style.Bold)
}
