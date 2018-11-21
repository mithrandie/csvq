package cmd

func WriteToStdout(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	_, err := Stdout.Write([]byte(s))
	return err
}

func WriteToStdErr(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	_, err := Stderr.Write([]byte(s))
	return err
}
