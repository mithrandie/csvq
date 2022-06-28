package terminal

import (
	"context"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/excmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/query"
	"github.com/mithrandie/csvq/lib/value"
)

var promptLoadConfigTests = []struct {
	Prompt                   string
	ContinuousPrompt         string
	ExpectSequence           []PromptElement
	ExpectContinuousSequence []PromptElement
	Error                    string
}{
	{
		Prompt:                   "TEST P > ",
		ContinuousPrompt:         "TEST C > ",
		ExpectSequence:           []PromptElement{{Text: "TEST P > ", Type: excmd.FixedString}},
		ExpectContinuousSequence: []PromptElement{{Text: "TEST C > ", Type: excmd.FixedString}},
	},
	{
		Prompt:                   "TEST P @ > ",
		ContinuousPrompt:         "TEST C > ",
		ExpectSequence:           nil,
		ExpectContinuousSequence: nil,
		Error:                    "prompt: invalid variable symbol",
	},
	{
		Prompt:                   "TEST P > ",
		ContinuousPrompt:         "TEST C @ > ",
		ExpectSequence:           nil,
		ExpectContinuousSequence: nil,
		Error:                    "prompt: invalid variable symbol",
	},
}

func TestPrompt_LoadConfig(t *testing.T) {
	prompt := NewPrompt(query.NewReferenceScope(TestTx))

	for _, v := range promptLoadConfigTests {
		TestTx.Environment.InteractiveShell.Prompt = v.Prompt
		TestTx.Environment.InteractiveShell.ContinuousPrompt = v.ContinuousPrompt

		err := prompt.LoadConfig()
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %s, %s", err, v.Prompt, v.ContinuousPrompt)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %s, %s", err.Error(), v.Error, v.Prompt, v.ContinuousPrompt)
			}
		} else if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %s, %s", v.Error, v.Prompt, v.ContinuousPrompt)
		}

		if !reflect.DeepEqual(prompt.sequence, v.ExpectSequence) {
			t.Errorf("sequence = %v, want %v for %s, %s", prompt.sequence, v.ExpectSequence, v.Prompt, v.ContinuousPrompt)
		}
		if !reflect.DeepEqual(prompt.continuousSequence, v.ExpectContinuousSequence) {
			t.Errorf("continuous sequence = %v, want %v for %s, %s", prompt.continuousSequence, v.ExpectContinuousSequence, v.Prompt, v.ContinuousPrompt)
		}
	}
}

var promptRenderPromptTests = []struct {
	Sequence []PromptElement
	UseColor bool
	Expect   string
	Error    string
}{
	{
		Sequence: nil,
		UseColor: false,
		Expect:   DefaultPrompt,
	},
	{
		Sequence: []PromptElement{
			{Text: "\033[32mstr\033[0m", Type: excmd.FixedString},
		},
		UseColor: true,
		Expect:   "\033[32mstr\033[0m",
	},
	{
		Sequence: []PromptElement{
			{Text: "str", Type: excmd.FixedString},
		},
		UseColor: true,
		Expect:   "\033[34mstr\033[0m",
	},
	{
		Sequence: []PromptElement{
			{Text: "\033[32mstr\033[0m", Type: excmd.FixedString},
		},
		UseColor: false,
		Expect:   "str",
	},
}

func TestPrompt_RenderPrompt(t *testing.T) {
	defer func() {
		TestTx.UseColor(false)
		initFlag(TestTx.Flags)
	}()

	prompt := NewPrompt(query.NewReferenceScope(TestTx))

	for _, v := range promptRenderPromptTests {
		TestTx.UseColor(v.UseColor)
		prompt.sequence = v.Sequence
		result, err := prompt.RenderPrompt(context.Background())

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %v", err, v.Sequence)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %v", err.Error(), v.Error, v.Sequence)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %v", v.Error, v.Sequence)
			continue
		}

		if result != v.Expect {
			t.Errorf("result = %q, want %q for %v", result, v.Expect, v.Sequence)
		}
	}
}

var promptRenderContinuousPromptTests = []struct {
	Sequence []PromptElement
	UseColor bool
	Expect   string
	Error    string
}{
	{
		Sequence: nil,
		UseColor: false,
		Expect:   DefaultContinuousPrompt,
	},
	{
		Sequence: []PromptElement{
			{Text: "\033[32mstr\033[0m", Type: excmd.FixedString},
		},
		UseColor: true,
		Expect:   "\033[32mstr\033[0m",
	},
	{
		Sequence: []PromptElement{
			{Text: "str", Type: excmd.FixedString},
		},
		UseColor: true,
		Expect:   "\033[34mstr\033[0m",
	},
	{
		Sequence: []PromptElement{
			{Text: "\033[32mstr\033[0m", Type: excmd.FixedString},
		},
		UseColor: false,
		Expect:   "str",
	},
}

func TestPrompt_RenderContinuousPrompt(t *testing.T) {
	defer func() {
		TestTx.UseColor(false)
		initFlag(TestTx.Flags)
	}()

	prompt := NewPrompt(query.NewReferenceScope(TestTx))

	for _, v := range promptRenderContinuousPromptTests {
		TestTx.UseColor(v.UseColor)
		prompt.continuousSequence = v.Sequence
		result, err := prompt.RenderContinuousPrompt(context.Background())

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %v", err, v.Sequence)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %v", err.Error(), v.Error, v.Sequence)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %v", v.Error, v.Sequence)
			continue
		}

		if result != v.Expect {
			t.Errorf("result = %q, want %q for %v", result, v.Expect, v.Sequence)
		}
	}
}

var promptRenderTests = []struct {
	Input  []PromptElement
	Expect string
	Error  string
}{
	{
		Input:  nil,
		Expect: "",
	},
	{
		Input: []PromptElement{
			{Text: "str", Type: excmd.FixedString},
			{Text: "var", Type: excmd.Variable},
			{Text: "CSVQ_TEST_ENV", Type: excmd.EnvironmentVariable},
			{Text: "VERSION", Type: excmd.RuntimeInformation},
			{Text: "@var", Type: excmd.CsvqExpression},
		},
		Expect: "strabcfoov1.0.0abc",
	},
	{
		Input: []PromptElement{
			{Text: "notexist", Type: excmd.Variable},
		},
		Error: "prompt: variable @notexist is undeclared",
	},
	{
		Input: []PromptElement{
			{Text: "NOTEXIST", Type: excmd.RuntimeInformation},
		},
		Error: "prompt: @#NOTEXIST is an unknown runtime information",
	},
	{
		Input: []PromptElement{
			{Text: "invalid invalid", Type: excmd.CsvqExpression},
		},
		Error: "prompt: syntax error: unexpected token \"invalid\"",
	},
	{
		Input: []PromptElement{
			{Text: "", Type: excmd.CsvqExpression},
		},
		Expect: "",
	},
	{
		Input: []PromptElement{
			{Text: "print 1;", Type: excmd.CsvqExpression},
		},
		Error: "prompt: print 1;: cannot evaluate as a value",
	},
	{
		Input: []PromptElement{
			{Text: "1;2", Type: excmd.CsvqExpression},
		},
		Error: "prompt: 1;2: cannot evaluate as a value",
	},
	{
		Input: []PromptElement{
			{Text: "@notexist", Type: excmd.CsvqExpression},
		},
		Error: "prompt: variable @notexist is undeclared",
	},
}

func TestPrompt_Render(t *testing.T) {
	scope := query.NewReferenceScope(TestTx)
	_ = scope.DeclareVariableDirectly(parser.Variable{Name: "var"}, value.NewString("abc"))
	prompt := NewPrompt(scope)

	for _, v := range promptRenderTests {
		result, err := prompt.Render(context.Background(), v.Input)

		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("unexpected error %q for %v", err, v.Input)
			} else if err.Error() != v.Error {
				t.Errorf("error %q, want error %q for %v", err.Error(), v.Error, v.Input)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("no error, want error %q for %v", v.Error, v.Input)
			continue
		}

		if result != v.Expect {
			t.Errorf("result = %q, want %q for %v", result, v.Expect, v.Input)
		}
	}
}

var promptStripEscapeSequenceTests = []struct {
	Input  string
	Expect string
}{
	{
		Input:  "\u001b[34;1m/path/to/working/directory \u001b[33;1m(Uncommitted:2)\u001b[34;1m >\u001b[0m ",
		Expect: "/path/to/working/directory (Uncommitted:2) > ",
	},
}

func TestPrompt_StripEscapeSequence(t *testing.T) {
	prompt := NewPrompt(query.NewReferenceScope(TestTx))

	for _, v := range promptStripEscapeSequenceTests {
		result := prompt.StripEscapeSequence(v.Input)

		if result != v.Expect {
			t.Errorf("result = %q, want %q for %v", result, v.Expect, v.Input)
		}
	}
}
