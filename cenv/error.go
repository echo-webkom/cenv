package cenv

import (
	"errors"
	"fmt"
)

type longError struct {
	s   string
	err bool
}

func (err *longError) Add(s string) {
	err.s += fmt.Sprintf("cenv: %s\n", s)
	err.err = true
}

func (err *longError) Error() error {
	if err.err {
		return errors.New(err.s)
	}
	return nil
}
