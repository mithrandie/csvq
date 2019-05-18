package excmd

import (
	"reflect"
	"testing"
)

var argsSplitterScanTests = []struct {
	Args   string
	Expect []string
	Error  string
}{
	{
		Args: "cmd -opt arg1 'arg  2' 'arg\\'3' \"arg\\\\4\" ",
		Expect: []string{
			"cmd",
			"-opt",
			"arg1",
			"'arg  2'",
			"'arg\\'3'",
			"\"arg\\\\4\"",
		},
	},
	{
		Args: "cmd arg1 arg2",
		Expect: []string{
			"cmd",
			"arg1",
			"arg2",
		},
	},
	{
		Args: "cmd arg1 'arg 2'",
		Expect: []string{
			"cmd",
			"arg1",
			"'arg 2'",
		},
	},
	{
		Args: "cmd arg1 'arg''2'",
		Expect: []string{
			"cmd",
			"arg1",
			"'arg''2'",
		},
	},
	{
		Args: "cmd arg1 @var",
		Expect: []string{
			"cmd",
			"arg1",
			"@var",
		},
	},
	{
		Args: "cmd arg1 @%var",
		Expect: []string{
			"cmd",
			"arg1",
			"@%var",
		},
	},
	{
		Args: "cmd arg1 @%`var \\`2`",
		Expect: []string{
			"cmd",
			"arg1",
			"@%`var \\`2`",
		},
	},
	{
		Args: "cmd arg1 @#var",
		Expect: []string{
			"cmd",
			"arg1",
			"@#var",
		},
	},
	{
		Args: "cmd arg1 ${'arg''2'}",
		Expect: []string{
			"cmd",
			"arg1",
			"${'arg''2'}",
		},
	},
	{
		Args: "sh -c ${format('echo %s | wc', @%HOME)}",
		Expect: []string{
			"sh",
			"-c",
			"${format('echo %s | wc', @%HOME)}",
		},
	},
	{
		Args: "sh -c ${format('echo \"${\\\\\\}\"')}",
		Expect: []string{
			"sh",
			"-c",
			"${format('echo \"${\\\\}\"')}",
		},
	},
	{
		Args:  "cmd arg1 'arg2",
		Error: "string not terminated",
	},
	{
		Args:  "sh -c ${format('echo %s | wc', @%HOME)",
		Error: "command not terminated",
	},
	{
		Args:  "sh -c $format('echo %s | wc', @%HOME)",
		Error: "invalid command symbol",
	},
	{
		Args:  "cmd @%`var",
		Error: "environment variable not terminated",
	},
}

func TestArgsSplitter_Scan(t *testing.T) {
	for _, v := range argsSplitterScanTests {
		splitter := new(ArgsSplitter).Init(v.Args)
		var args = make([]string, 0)
		for splitter.Scan() {
			args = append(args, splitter.Text())
		}
		err := splitter.Err()

		if err != nil {
			if v.Error == "" {
				t.Errorf("unexpected error %q for %q", err.Error(), v.Args)
			} else if v.Error != err.Error() {
				t.Errorf("error %q, want error %q for %q", err.Error(), v.Error, v.Args)
			}
			continue
		}
		if v.Error != "" {
			t.Errorf("no error, want error %q for %q", v.Error, v.Args)
			continue
		}

		if !reflect.DeepEqual(args, v.Expect) {
			t.Errorf("result = %#v, want %#v for %q", args, v.Expect, v.Args)
		}
	}
}
