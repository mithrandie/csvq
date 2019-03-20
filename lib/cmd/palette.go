package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/mithrandie/csvq/lib/file"

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
	palette     *color.Palette
	loadPalette sync.Once
)

func GetPalette() *color.Palette {
	if palette == nil {
		env, err := NewEnvironment(context.Background(), file.DefaultWaitTimeout, file.DefaultRetryDelay)
		if err != nil {
			println(err.Error())
		}
		if err := LoadPalette(env); err != nil {
			println(err.Error())
		}
	}
	return palette
}

func LoadPalette(env *Environment) (err error) {
	loadPalette.Do(func() {
		p, err := color.GeneratePalette(env.Palette)
		if err != nil {
			err = errors.New(fmt.Sprintf("palette configuration error: %s", err.Error()))
			return
		}

		palette = p
	})
	return
}

func Error(s string) string {
	if p := GetPalette(); p != nil {
		return p.Render(ErrorEffect, s)
	}
	return s
}

func Warn(s string) string {
	if p := GetPalette(); p != nil {
		return p.Render(WarnEffect, s)
	}
	return s
}

func Notice(s string) string {
	if p := GetPalette(); p != nil {
		return p.Render(NoticeEffect, s)
	}
	return s
}
