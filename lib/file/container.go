package file

import (
	"errors"
	"fmt"
	"strings"
)

type Container struct {
	m map[string]*Handler
}

func NewContainer() *Container {
	return &Container{
		m: make(map[string]*Handler),
	}
}

func (c *Container) Keys() []string {
	l := make([]string, 0, len(c.m))
	for k := range c.m {
		l = append(l, k)
	}
	return l
}

func (c *Container) Add(path string, handler *Handler) error {
	key := strings.ToUpper(path)
	if _, ok := c.m[key]; ok {
		return errors.New(fmt.Sprintf("file %s already opened", path))
	}
	c.m[key] = handler
	return nil
}

func (c *Container) Remove(path string) {
	key := strings.ToUpper(path)
	if _, ok := c.m[key]; ok {
		delete(c.m, key)
	}
}

func (c *Container) Close(h *Handler) error {
	if h == nil {
		return nil
	}

	key := strings.ToUpper(h.Path())
	if _, ok := c.m[key]; ok {
		if err := c.m[key].close(); err != nil {
			return err
		}
		c.Remove(h.Path())
	}
	return nil
}

func (c *Container) Commit(h *Handler) error {
	if h == nil {
		return nil
	}

	key := strings.ToUpper(h.Path())
	if _, ok := c.m[key]; ok {
		if err := c.m[key].commit(); err != nil {
			return err
		}
		c.Remove(h.Path())
	}
	return nil
}

func (c *Container) CloseWithErrors(h *Handler) (err error) {
	if h == nil {
		return nil
	}

	key := strings.ToUpper(h.Path())
	if _, ok := c.m[key]; ok {
		err = c.m[key].closeWithErrors()
		c.Remove(h.Path())
	}
	return
}

func (c *Container) CloseAll() error {
	for k := range c.m {
		if err := c.Close(c.m[k]); err != nil {
			return err
		}
	}
	return nil
}

func (c *Container) CloseAllWithErrors() error {
	var errs []error
	for k := range c.m {
		if err := c.CloseWithErrors(c.m[k]); err != nil {
			errs = append(errs, err.(*ForcedUnlockError).Errors...)
		}
	}

	return NewForcedUnlockError(errs)
}
