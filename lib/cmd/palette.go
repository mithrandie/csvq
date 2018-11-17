package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mithrandie/go-text/color"
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
	ErrorEffect      = "error"
	WarnEffect       = "warn"
	NoticeEffect     = "notice"
)

var (
	palette    *color.Palette
	getPalette sync.Once
)

func GetPalette() (*color.Palette, error) {
	var err error

	getPalette.Do(func() {
		var env *Environment
		env, err = GetEnvironment()
		if err != nil {
			return
		}

		palette, err = color.GeneratePalette(env.Palette)
		if err != nil {
			err = errors.New(fmt.Sprintf("palette configuration error: %s", err.Error()))
		}
	})

	return palette, err
}

func Error(s string) string {
	if p, err := GetPalette(); err == nil && p != nil {
		return p.Render(ErrorEffect, s)
	}
	return s
}

func Warn(s string) string {
	if p, err := GetPalette(); err == nil && p != nil {
		return p.Render(WarnEffect, s)
	}
	return s
}

func Notice(s string) string {
	if p, err := GetPalette(); err == nil && p != nil {
		return p.Render(NoticeEffect, s)
	}
	return s
}
