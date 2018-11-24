package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/go-text/color"
)

const (
	ConfigDir              = ".config"
	CSVQConfigDir          = "csvq"
	EnvFileName            = "csvq_env.json"
	PreloadCommandFileName = "csvqrc"

	HiddenPrefix = '.'
)

var (
	environment *Environment
	getEnv      sync.Once
)

type Environment struct {
	DatetimeFormat       []string            `json:"datetime_format"`
	InteractiveShell     InteractiveShell    `json:"interactive_shell"`
	EnvironmentVariables map[string]string   `json:"environment_variables"`
	Palette              color.PaletteConfig `json:"palette"`
}

func (e *Environment) Merge(e2 *Environment) {
	for _, f := range e2.DatetimeFormat {
		e.DatetimeFormat = AppendStrIfNotExist(e.DatetimeFormat, f)
	}

	if 0 < len(e2.InteractiveShell.HistoryFile) {
		e.InteractiveShell.HistoryFile = e2.InteractiveShell.HistoryFile
	}

	if e2.InteractiveShell.HistoryLimit != 0 {
		e.InteractiveShell.HistoryLimit = e2.InteractiveShell.HistoryLimit
	}

	if 0 < len(e2.InteractiveShell.Prompt) {
		e.InteractiveShell.Prompt = e2.InteractiveShell.Prompt
	}

	if 0 < len(e2.InteractiveShell.ContinuousPrompt) {
		e.InteractiveShell.ContinuousPrompt = e2.InteractiveShell.ContinuousPrompt
	}

	for k, v := range e2.EnvironmentVariables {
		e.EnvironmentVariables[k] = v
	}

	for k, v := range e2.Palette.Effectors {
		e.Palette.Effectors[k] = v
	}
}

type InteractiveShell struct {
	HistoryFile      string `json:"history_file"`
	HistoryLimit     int    `json:"history_limit"`
	Prompt           string `json:"prompt"`
	ContinuousPrompt string `json:"continuous_prompt"`
}

func LoadEnvironment() error {
	var err error

	environment = &Environment{}
	if err = json.Unmarshal([]byte(DefaultEnvJson), environment); err != nil {
		return errors.New(fmt.Sprintf("`json syntax error: %s", err.Error()))
	}

	handlers := make([]*file.Handler, 0, 4)
	defer func() {
		for _, h := range handlers {
			h.Close()
		}
	}()

	files := GetSpecialFilePath(EnvFileName)
	for _, fpath := range files {
		if !file.Exists(fpath) {
			continue
		}

		var h *file.Handler
		var buf []byte

		h, err = file.NewHandlerForRead(fpath)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to load %q: %s", fpath, err.Error()))
		}
		handlers = append(handlers, h)

		buf, err = ioutil.ReadAll(h.FileForRead())
		if err != nil {
			return errors.New(fmt.Sprintf("failed to load %q: %s", fpath, err.Error()))
		}
		buf = bytes.TrimSuffix(buf, []byte{0x00})
		userDefinedEnv := &Environment{}
		if err = json.Unmarshal(buf, userDefinedEnv); err != nil {
			return errors.New(fmt.Sprintf("failed to load %q: %s", fpath, err.Error()))
		}

		environment.Merge(userDefinedEnv)
	}

	for k, v := range environment.EnvironmentVariables {
		os.Setenv(k, v)
	}

	return nil
}

func GetEnvironment() (*Environment, error) {
	var err error

	getEnv.Do(func() {
		err = LoadEnvironment()
	})

	return environment, err
}

func GetSpecialFilePath(filename string) []string {
	files := make([]string, 0, 4)
	files = AppendStrIfNotExist(files, GetHomeDirFilePath(filename))
	files = AppendStrIfNotExist(files, GetCSVQConfigDirFilePath(filename))
	files = AppendStrIfNotExist(files, GetConfigDirFilePath(filename))
	files = AppendStrIfNotExist(files, GetCurrentDirFilePath(filename))
	return files
}

func GetHomeDirFilePath(filename string) string {
	home, err := homedir.Dir()
	if err != nil {
		return filename
	}

	if 0 < len(filename) && filename[0] != HiddenPrefix {
		filename = string(HiddenPrefix) + filename
	}

	return filepath.Join(home, filename)
}

func GetCSVQConfigDirFilePath(filename string) string {
	home, err := homedir.Dir()
	if err != nil {
		return filename
	}

	return filepath.Join(home, string(HiddenPrefix)+CSVQConfigDir, filename)
}

func GetConfigDirFilePath(filename string) string {
	home, err := homedir.Dir()
	if err != nil {
		return filename
	}

	return filepath.Join(home, ConfigDir, CSVQConfigDir, filename)
}

func GetCurrentDirFilePath(filename string) string {
	if !filepath.IsAbs(filename) {
		if abs, err := filepath.Abs(filename); err == nil {
			filename = abs
		}
	}
	return filename
}
