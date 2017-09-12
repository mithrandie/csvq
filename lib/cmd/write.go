package cmd

import (
	"bufio"
	"os"

	"github.com/mithrandie/go-file"
)

var Terminal VirtualTerminal

func ToStdout(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}
	return CreateFile("", s)
}

func CreateFile(filename string, s string) error {
	var fp *os.File
	var err error

	if len(filename) < 1 {
		fp = os.Stdout
	} else {
		fp, err = file.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close(fp)
	}

	w := bufio.NewWriter(fp)
	_, err = w.WriteString(s)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func UpdateFile(fp *os.File, s string) error {
	defer file.Close(fp)

	fp.Truncate(0)
	fp.Seek(0, 0)

	w := bufio.NewWriter(fp)
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	w.Flush()

	return nil
}

func TryCreateFile(filename string) error {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	fp.Close()
	os.Remove(filename)

	return nil
}
