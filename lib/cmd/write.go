package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ToStdout(s string) error {
	return CreateFile("", s)
}

func CreateFile(filename string, s string) error {
	var fp *os.File
	var err error

	if len(filename) < 1 {
		fp = os.Stdout
	} else {
		if _, err := os.Stat(filename); err == nil {
			return errors.New(fmt.Sprintf("file %s already exists", filename))
		}
		fp, err = os.Create(filename)
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

func UpdateFile(filename string, s string) error {
	var fp *os.File
	var err error

	fp, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
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

func TryCreateFile(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return errors.New(fmt.Sprintf("file %s already exists", filename))
	}

	fp, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		fp.Close()
		os.Remove(filename)
	}()

	return nil
}

func TryOpenFileToWrite(filename string) error {
	fp, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer fp.Close()
	return nil
}
