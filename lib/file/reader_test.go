package file

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

type readerReadResult struct {
	Result string
	EOF    bool
}

var readerReadTests = []struct {
	Source      string
	HeaderLen   int
	ReadByteLen int
	Results     []readerReadResult
}{
	{
		Source:      "abcdefg",
		HeaderLen:   3,
		ReadByteLen: 2,
		Results: []readerReadResult{
			{
				Result: "ab",
				EOF:    false,
			},
			{
				Result: "cd",
				EOF:    false,
			},
			{
				Result: "ef",
				EOF:    false,
			},
			{
				Result: "g",
				EOF:    false,
			},
			{
				Result: "",
				EOF:    true,
			},
		},
	},
	{
		Source:      "abcdef",
		HeaderLen:   3,
		ReadByteLen: 2,
		Results: []readerReadResult{
			{
				Result: "ab",
				EOF:    false,
			},
			{
				Result: "cd",
				EOF:    false,
			},
			{
				Result: "ef",
				EOF:    false,
			},
			{
				Result: "",
				EOF:    true,
			},
		},
	},
	{
		Source:      "abcd",
		HeaderLen:   2,
		ReadByteLen: 2,
		Results: []readerReadResult{
			{
				Result: "ab",
				EOF:    false,
			},
			{
				Result: "cd",
				EOF:    false,
			},
			{
				Result: "",
				EOF:    true,
			},
		},
	},
	{
		Source:      "ab",
		HeaderLen:   10,
		ReadByteLen: 2,
		Results: []readerReadResult{
			{
				Result: "ab",
				EOF:    false,
			},
			{
				Result: "",
				EOF:    true,
			},
		},
	},
	{
		Source:      "",
		HeaderLen:   2,
		ReadByteLen: 2,
		Results: []readerReadResult{
			{
				Result: "",
				EOF:    true,
			},
		},
	},
}

func TestReader_Read(t *testing.T) {
	for _, v := range readerReadTests {
		p := make([]byte, v.ReadByteLen)
		reader, _ := NewReader(strings.NewReader(v.Source), v.HeaderLen)

		for i, expect := range v.Results {
			n, err := reader.Read(p)

			if err != nil && err != io.EOF {
				t.Errorf("%s - %d: unexpected error %q", v.Source, i, err)
				continue
			}

			if err == io.EOF && !expect.EOF {
				t.Errorf("%s - %d: error is nil, want io.EOF", v.Source, i)
			}

			result := string(p[:n])
			if result != expect.Result {
				t.Errorf("%s - %d: result = %q, want %q", v.Source, i, result, expect.Result)
			}
		}
	}
}

func TestReader_HeadBytes(t *testing.T) {
	source := strings.NewReader("abcdefg")
	reader, _ := NewReader(source, 2)
	headBytes, _ := reader.HeadBytes()

	expect := "ab"
	result, _ := io.ReadAll(headBytes)
	if string(result) != expect {
		t.Errorf("result = %s, want %s", result, expect)
	}
}

var readerSizeTests = []struct {
	Name   string
	Reader io.Reader
	Result int64
}{
	{
		Name:   "*strings.Reader",
		Reader: strings.NewReader("abcdefg"),
		Result: 7,
	},
	{
		Name:   "*bytes.Reader",
		Reader: bytes.NewReader([]byte("abcdefg")),
		Result: 7,
	},
}

func TestReader_Size(t *testing.T) {
	for _, v := range readerSizeTests {
		reader, _ := NewReader(v.Reader, 2)
		result := reader.Size()
		if result != v.Result {
			t.Errorf("%s: result = %d, want %d", v.Name, result, v.Result)
		}
	}
}
