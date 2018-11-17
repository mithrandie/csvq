package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestGetConfigDirFilePath(t *testing.T) {
	pwd, _ := os.Getwd()
	home, _ := homedir.Dir()

	filename := "file.txt"
	expect := []string{
		filepath.Join(home, string(HiddenPrefix)+filename),
		filepath.Join(home, string(HiddenPrefix)+CSVQConfigDir, filename),
		filepath.Join(home, ConfigDir, CSVQConfigDir, filename),
		filepath.Join(pwd, filename),
	}
	result := GetSpecialFilePath(filename)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}

	filename = ""
	expect = []string{
		home,
		filepath.Join(home, string(HiddenPrefix)+CSVQConfigDir),
		filepath.Join(home, ConfigDir, CSVQConfigDir),
		filepath.Join(pwd),
	}
	result = GetSpecialFilePath(filename)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}

	s := os.ExpandEnv("$HOME/foo")
	fmt.Println(s)
	vars := os.Environ()
	for _, v := range vars {
		fmt.Println(v)
	}
}
