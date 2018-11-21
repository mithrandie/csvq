package cmd

func WriteToStdout(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	_, err := Stdout.Write([]byte(s))
	return err
}

func WriteToStdoutWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return WriteToStdout(s)
}

func WriteToStderr(s string) error {
	if Terminal != nil {
		return Terminal.Write(s)
	}

	_, err := Stderr.Write([]byte(s))
	return err
}

func WriteToStderrWithLineBreak(s string) error {
	if 0 < len(s) && s[len(s)-1] != '\n' {
		s = s + "\n"
	}
	return WriteToStderr(s)
}
