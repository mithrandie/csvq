package cmd

import (
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
	ConfigDir     = ".config"
	CSVQConfigDir = "csvq"
	EnvFileName   = "csvq_env.json"
	RCFileName    = "csvqrc"

	HiddenPrefix = '.'
)

var (
	environment *Environment
	getEnv      sync.Once
)

type Environment struct {
	EnvironmentVariables map[string]string   `json:"environment_variables"`
	Palette              color.PaletteConfig `json:"palette"`
}

func (e *Environment) Merge(e2 *Environment) {
	for k, v := range e2.EnvironmentVariables {
		e.EnvironmentVariables[k] = v
	}

	for k, v := range e2.Palette.Effectors {
		e.Palette.Effectors[k] = v
	}
}

func LoadEnvironment() error {
	var err error

	environment = &Environment{}
	if err = json.Unmarshal([]byte(defaultEnvJson), environment); err != nil {
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
	var appendToList = func(list []string, fpath string) []string {
		if len(fpath) < 1 {
			return list
		}
		for _, v := range list {
			if fpath == v {
				return list
			}
		}
		return append(list, fpath)
	}

	files := make([]string, 0, 4)
	files = appendToList(files, GetHomeDirFilePath(filename))
	files = appendToList(files, GetCSVQConfigDirFilePath(filename))
	files = appendToList(files, GetConfigDirFilePath(filename))
	files = appendToList(files, GetCurrentDirFilePath(filename))
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
