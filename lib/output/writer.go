package output

import (
	"bufio"
	"os"
)

func Write(file string, s string) error {
	var fp *os.File
	var err error

	if len(file) < 1 {
		fp = os.Stdout
	} else {
		fp, err = os.Create(file)
		if err != nil {
			return err
		}
	}

	defer fp.Close()

	writer := bufio.NewWriter(fp)
	_, err = writer.WriteString(s)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}
