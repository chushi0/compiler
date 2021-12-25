package lexer

import (
	"errors"
	"fmt"
)

var (
	ErrorFinateAutomatonInput = errors.New("finate automaton not accpet this input")
	ErrorRegexpParse          = errors.New("parse fail")
)

type IOError struct {
	Original error
}

func (e *IOError) Error() string {
	return fmt.Sprintf("io exception: %v", e.Original.Error())
}

func (e *IOError) Unwrap() error {
	return e.Original
}
