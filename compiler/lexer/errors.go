package lexer

import "fmt"

type IOError struct {
	Original error
}

func (e *IOError) Error() string {
	return fmt.Sprintf("io exception: %v", e.Original.Error())
}

func (e *IOError) Unwrap() error {
	return e.Original
}
