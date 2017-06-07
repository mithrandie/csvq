package output

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ToStdout(s string) error {
	return Create("", s)
}

func Create(file string, s string) error {
	var fp *os.File
	var err error

	if len(file) < 1 {
		fp = os.Stdout
	} else {
		if _, err := os.Stat(file); err == nil {
			return errors.New(fmt.Sprintf("file %s already exists", file))
		}
		fp, err = os.Create(file)
		if err != nil {
			return err
		}
		defer fp.Close()
	}

	w := bufio.NewWriter(fp)
	_, err = w.WriteString(s)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func Update(file string, s string) error {
	var fp *os.File
	var err error

	fp, err = os.OpenFile(file, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer fp.Close()

	w := bufio.NewWriter(fp)
	_, err = w.WriteString(s)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}
