package cmd

import (
	"sync"

	"github.com/mithrandie/go-text/color"
	"github.com/mithrandie/go-text/json"
)

const (
	NoEffect         = ""
	LableEffect      = "label"
	NumberEffect     = "number"
	StringEffect     = "string"
	BooleanEffect    = "boolean"
	TernaryEffect    = "ternary"
	DatetimeEffect   = "datetime"
	NullEffect       = "null"
	ObjectEffect     = "object"
	AttributeEffect  = "attribute"
	IdentifierEffect = "identifier"
	ValueEffect      = "value"
	EmphasisEffect   = "emphasis"
	PromptEffect     = "prompt"
)

var (
	palette    *color.Palette
	getPalette sync.Once
)

func GetPalette() *color.Palette {
	getPalette.Do(func() {
		label := color.NewEffector()
		label.SetEffect(color.Bold)
		label.SetFGColor(color.Blue)

		num := color.NewEffector()
		num.SetFGColor(color.Magenta)

		str := color.NewEffector()
		str.SetFGColor(color.Green)

		b := color.NewEffector()
		b.SetFGColor(color.Yellow)
		b.SetEffect(color.Bold)

		t := color.NewEffector()
		t.SetFGColor(color.Yellow)

		dt := color.NewEffector()
		dt.SetFGColor(color.Cyan)

		n := color.NewEffector()
		n.SetFGColor(color.BrightBlack)

		obj := color.NewEffector()
		obj.SetFGColor(color.Green)
		obj.SetEffect(color.Bold)

		attr := color.NewEffector()
		attr.SetFGColor(color.Yellow)

		ident := color.NewEffector()
		ident.SetFGColor(color.Cyan)
		ident.SetEffect(color.Bold)

		val := color.NewEffector()
		val.SetFGColor(color.Blue)
		val.SetEffect(color.Bold)

		emphasis := color.NewEffector()
		emphasis.SetFGColor(color.Red)
		emphasis.SetEffect(color.Bold)

		prompt := color.NewEffector()
		prompt.SetFGColor(color.Blue)

		palette = color.NewPalette()
		palette.SetEffector(LableEffect, label)
		palette.SetEffector(NumberEffect, num)
		palette.SetEffector(StringEffect, str)
		palette.SetEffector(BooleanEffect, b)
		palette.SetEffector(TernaryEffect, t)
		palette.SetEffector(DatetimeEffect, dt)
		palette.SetEffector(NullEffect, n)
		palette.SetEffector(ObjectEffect, obj)
		palette.SetEffector(AttributeEffect, attr)
		palette.SetEffector(IdentifierEffect, ident)
		palette.SetEffector(ValueEffect, val)
		palette.SetEffector(EmphasisEffect, emphasis)
		palette.SetEffector(PromptEffect, prompt)
		palette.SetEffector(json.ObjectKeyEffect, label)
		palette.SetEffector(json.NumberEffect, num)
		palette.SetEffector(json.StringEffect, str)
		palette.SetEffector(json.BooleanEffect, b)
		palette.SetEffector(json.NullEffect, n)
	})

	return palette
}
