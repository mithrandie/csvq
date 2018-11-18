package file

import (
	"errors"
	"fmt"
	"strings"
)

var container = make(map[string]*Handler)

func addToContainer(path string, handler *Handler) error {
	key := strings.ToUpper(path)
	if _, ok := container[key]; ok {
		return errors.New(fmt.Sprintf("file %s already opened", path))
	}
	container[key] = handler
	return nil
}

func removeFromContainer(path string) {
	key := strings.ToUpper(path)
	if _, ok := container[key]; ok {
		delete(container, key)
	}
}

func UnlockAll() error {
	for k := range container {
		if err := container[k].Close(); err != nil {
			return err
		}
		delete(container, k)
	}
	return nil
}

func UnlockAllWithErrors() error {
	var errs []error
	for k := range container {
		if err := container[k].CloseWithErrors(); err != nil {
			errs = append(errs, err.(*ForcedUnlockError).Errors...)
		}
		delete(container, k)
	}

	if errs != nil {
		return NewForcedUnlockError(errs)
	}
	return nil
}
