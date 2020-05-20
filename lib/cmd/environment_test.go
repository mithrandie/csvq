package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestGetConfigDirFilePath(t *testing.T) {
	oldConfigHome := os.Getenv(XDGConfigHomeEnvName)
	defer func() {
		_ = os.Setenv(XDGConfigHomeEnvName, oldConfigHome)
	}()

	pwd, _ := os.Getwd()
	home, _ := homedir.Dir()
	xdgConfigHome := filepath.Join("home", "mithrandie")
	_ = os.Setenv(XDGConfigHomeEnvName, xdgConfigHome)

	filename := "file.txt"
	expect := []string{
		filepath.Join(xdgConfigHome, CSVQConfigDir, filename),
		filepath.Join(home, string(HiddenPrefix)+filename),
		filepath.Join(home, string(HiddenPrefix)+CSVQConfigDir, filename),
		filepath.Join(pwd, filename),
	}
	result := GetSpecialFilePath(filename)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}

	filename = ""
	expect = []string{
		filepath.Join(xdgConfigHome, CSVQConfigDir),
		home,
		filepath.Join(home, string(HiddenPrefix)+CSVQConfigDir),
		filepath.Join(pwd),
	}
	result = GetSpecialFilePath(filename)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}

	_ = os.Unsetenv(XDGConfigHomeEnvName)
	expect = []string{
		filepath.Join(home, DefaultXDGConfigDir, CSVQConfigDir, filename),
		filepath.Join(home, string(HiddenPrefix)+filename),
		filepath.Join(home, string(HiddenPrefix)+CSVQConfigDir, filename),
		filepath.Join(pwd, filename),
	}
	result = GetSpecialFilePath(filename)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, want %v", result, expect)
	}
}
