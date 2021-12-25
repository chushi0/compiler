package lexer_test

import (
	"errors"
	"testing"

	"github.com/chushi0/compiler/compiler/lexer"
)

func TestFinateAutomaton1(t *testing.T) {
	fa, err := lexer.NewFinateAutomatonFromRegexp([]rune("hello"))
	if err != nil {
		t.FailNow()
	}
	fa = fa.AsDFA()
	input := "hello"
	state := 0
	for _, r := range input {
		next, err := fa.NextState(state, r)
		if err != nil {
			t.FailNow()
		}
		state = next
	}
	if !fa.AcceptStates.Contains(state) {
		t.Fail()
	}
}

func TestFinateAutomaton2(t *testing.T) {
	fa, err := lexer.NewFinateAutomatonFromRegexp([]rune("[0-9]*"))
	if err != nil {
		t.FailNow()
	}
	fa = fa.AsDFA()
	inputs := []string{
		"123456789",
		"111111",
		"",
		"5",
	}
	for _, input := range inputs {
		state := 0
		for _, r := range input {
			next, err := fa.NextState(state, r)
			if err != nil {
				t.Fail()
				goto next
			}
			state = next
		}
		if !fa.AcceptStates.Contains(state) {
			t.Fail()
		}

	next:
		continue
	}
}

func TestFinateAutomaton3(t *testing.T) {
	fa, err := lexer.NewFinateAutomatonFromRegexp([]rune("([a-z][0-9])*5"))
	if err != nil {
		t.FailNow()
	}
	fa = fa.AsDFA()
	inputs := []string{
		"a1b2c35",
		"5",
		"a55",
	}
	for _, input := range inputs {
		state := 0
		for _, r := range input {
			next, err := fa.NextState(state, r)
			if err != nil {
				t.Fail()
				goto next
			}
			state = next
		}
		if !fa.AcceptStates.Contains(state) {
			t.Fail()
		}

	next:
		continue
	}
}

func TestFinateAutomaton4(t *testing.T) {
	fa, err := lexer.NewFinateAutomatonFromRegexp([]rune(".*"))
	if err != nil {
		t.FailNow()
	}
	fa = fa.AsDFA()
	inputs := []string{
		"75275",
		"5",
		"",
		"中文",
		"にほん",
		"emoji😊",
	}
	for _, input := range inputs {
		state := 0
		for _, r := range input {
			next, err := fa.NextState(state, r)
			if err != nil {
				t.Fail()
				goto next
			}
			state = next
		}
		if !fa.AcceptStates.Contains(state) {
			t.Fail()
		}

	next:
		continue
	}
}

func TestFinateAutomaton5(t *testing.T) {
	errRegexps := []string{
		"[",
		"\\",
		"|",
		"*",
		"5||9",
		"(489",
		"a)",
		"ad(fasga(agds)(62)",
		"\\u556*",
	}

	for _, regexp := range errRegexps {
		_, err := lexer.NewFinateAutomatonFromRegexp([]rune(regexp))
		if err != nil && errors.Is(err, lexer.ErrorRegexpParse) {
			continue
		}
		t.Fail()
	}
}
