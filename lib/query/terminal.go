package query

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/excmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/go-text/color"
)

const (
	TerminalPrompt           string = "csvq > "
	TerminalContinuousPrompt string = "     > "
)

type VirtualTerminal interface {
	ReadLine() (string, error)
	Write(string) error
	WriteError(string) error
	SetPrompt(ctx context.Context)
	SetContinuousPrompt(ctx context.Context)
	SaveHistory(string) error
	Teardown() error
	GetSize() (int, int, error)
	ReloadConfig() error
	UpdateCompleter()
}

type PromptEvaluationError struct {
	Message string
}

func NewPromptEvaluationError(message string) error {
	return &PromptEvaluationError{
		Message: message,
	}
}

func (e PromptEvaluationError) Error() string {
	return fmt.Sprintf("prompt: %s", e.Message)
}

type PromptElement struct {
	Text string
	Type excmd.ElementType
}

type Prompt struct {
	filter             *Filter
	palette            *color.Palette
	sequence           []PromptElement
	continuousSequence []PromptElement

	buf bytes.Buffer
}

func NewPrompt(filter *Filter, palette *color.Palette) *Prompt {
	return &Prompt{
		filter:  filter,
		palette: palette,
	}
}

func (p *Prompt) LoadConfig() error {
	p.sequence = nil
	p.continuousSequence = nil

	scanner := new(excmd.ArgumentScanner)

	scanner.Init(p.filter.tx.Environment.InteractiveShell.Prompt)
	for scanner.Scan() {
		p.sequence = append(p.sequence, PromptElement{
			Text: scanner.Text(),
			Type: scanner.ElementType(),
		})
	}
	if err := scanner.Err(); err != nil {
		p.sequence = nil
		p.continuousSequence = nil
		return NewPromptEvaluationError(err.Error())
	}

	scanner.Init(p.filter.tx.Environment.InteractiveShell.ContinuousPrompt)
	for scanner.Scan() {
		p.continuousSequence = append(p.continuousSequence, PromptElement{
			Text: scanner.Text(),
			Type: scanner.ElementType(),
		})
	}
	if err := scanner.Err(); err != nil {
		p.sequence = nil
		p.continuousSequence = nil
		return NewPromptEvaluationError(err.Error())
	}
	return nil
}

func (p *Prompt) RenderPrompt(ctx context.Context) (string, error) {
	s, err := p.Render(ctx, p.sequence)
	if err != nil || len(s) < 1 {
		s = TerminalPrompt
	}
	if p.filter.tx.Flags.Color {
		if strings.IndexByte(s, 0x1b) < 0 {
			s = p.palette.Render(cmd.PromptEffect, s)
		}
	} else {
		s = p.StripEscapeSequence(s)
	}
	return s, err
}

func (p *Prompt) RenderContinuousPrompt(ctx context.Context) (string, error) {
	s, err := p.Render(ctx, p.continuousSequence)
	if err != nil || len(s) < 1 {
		s = TerminalContinuousPrompt
	}
	if p.filter.tx.Flags.Color {
		if strings.IndexByte(s, 0x1b) < 0 {
			s = p.palette.Render(cmd.PromptEffect, s)
		}
	} else {
		s = p.StripEscapeSequence(s)
	}
	return s, err
}

func (p *Prompt) Render(ctx context.Context, sequence []PromptElement) (string, error) {
	p.buf.Reset()
	var err error

	for _, element := range sequence {
		switch element.Type {
		case excmd.FixedString:
			p.buf.WriteString(element.Text)
		case excmd.Variable:
			if err = p.evaluate(ctx, parser.Variable{Name: element.Text}); err != nil {
				return "", err
			}
		case excmd.EnvironmentVariable:
			if err = p.evaluate(ctx, parser.EnvironmentVariable{Name: element.Text}); err != nil {
				return "", err
			}
		case excmd.RuntimeInformation:
			if err = p.evaluate(ctx, parser.RuntimeInformation{Name: element.Text}); err != nil {
				return "", err
			}
		case excmd.CsvqExpression:
			if 0 < len(element.Text) {
				command := element.Text
				statements, err := parser.Parse(command, "", p.filter.tx.Flags.DatetimeFormat)
				if err != nil {
					syntaxErr := err.(*parser.SyntaxError)
					return "", NewPromptEvaluationError(syntaxErr.Message)
				}

				switch len(statements) {
				case 1:
					expr, ok := statements[0].(parser.QueryExpression)
					if !ok {
						return "", NewPromptEvaluationError(fmt.Sprintf(ErrorInvalidValue, command))
					}
					if err = p.evaluate(ctx, expr); err != nil {
						return "", err
					}
				default:
					return "", NewPromptEvaluationError(fmt.Sprintf(ErrorInvalidValue, command))
				}
			}
		}
	}

	return p.buf.String(), nil
}

func (p *Prompt) evaluate(ctx context.Context, expr parser.QueryExpression) error {
	val, err := p.filter.Evaluate(ctx, expr)
	if err != nil {
		if ae, ok := err.(AppError); ok {
			err = NewPromptEvaluationError(ae.ErrorMessage())
		}
		return err
	}
	s, _ := NewStringFormatter().Format("%s", []value.Primary{val})
	p.buf.WriteString(s)
	if err != nil {
		err = NewPromptEvaluationError(err.Error())
	}
	return err
}

func (p *Prompt) StripEscapeSequence(s string) string {
	buf := new(bytes.Buffer)

	inEscSeq := false
	for _, r := range s {
		if inEscSeq {
			if unicode.IsLetter(r) {
				inEscSeq = false
			}
		} else if r == 27 {
			inEscSeq = true
		} else {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}
