package cenv

import (
	"errors"
	"fmt"
)

type longError struct {
	s   string
	err bool
}

func (err *longError) Add(e error) {
	err.s += fmt.Sprintf("cenv: %s\n", e.Error())
	err.err = true
}

func (err *longError) AddMany(errs longError) {
	err.s += errs.s
	if errs.err {
		err.err = true
	}
}

func (err *longError) Error() error {
	if err.err {
		return errors.New(err.s)
	}
	return nil
}
