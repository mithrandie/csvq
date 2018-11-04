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
	ObjectStyle
	AttributeStyle
	IdentifierStyle
	ValueStyle
	EmphasisStyle
)

type TextStyle struct {
	Color int
	Bold  bool
}

type Palette struct {
	textStyles []*TextStyle
	useStyle   bool
}

func NewPalette() *Palette {
	return &Palette{
		textStyles: []*TextStyle{
			{Color: PlainColor, Bold: false},
			{Color: FGBlue, Bold: true},
			{Color: FGMagenta, Bold: false},
			{Color: FGGreen, Bold: false},
			{Color: FGYellow, Bold: true},
			{Color: FGYellow, Bold: false},
			{Color: FGCyan, Bold: false},
			{Color: FGBrightBlack, Bold: false},
			{Color: FGGreen, Bold: true},
			{Color: FGYellow, Bold: false},
			{Color: FGCyan, Bold: true},
			{Color: FGBlue, Bold: true},
			{Color: FGRed, Bold: true},
		},
		useStyle: false,
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

	idx := int(style)
	if len(p.textStyles) <= idx {
		idx = 0
	}

	return p.color(s, p.textStyles[idx])
}

func (p *Palette) color(s string, style *TextStyle) string {
	return Colorize(s, style.Color, style.Bold)
}
