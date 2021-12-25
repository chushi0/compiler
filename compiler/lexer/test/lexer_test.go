package lexer_test

import (
	"testing"

	"github.com/chushi0/compiler/compiler/lexer"
	"github.com/chushi0/compiler/compiler/types"
)

func TestLexer(t *testing.T) {
	identityFa, err := lexer.NewFinateAutomatonFromRegexp([]rune("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9])*"))
	if err != nil {
		t.Fatal(err)
	}
	numerFa, err := lexer.NewFinateAutomatonFromRegexp([]rune("([0-9]|[1-9][0-9]*)([]|\\.[0-9]*)([]|(E|e)[0-9]*)"))
	if err != nil {
		t.Fatal(err)
	}
	symbolFa, err := lexer.NewFinateAutomatonFromRegexp([]rune("+|-|\\*|/|=|\\(|\\)|;"))
	if err != nil {
		t.Fatal(err)
	}

	identityFa.SetAcceptTag("identify")
	numerFa.SetAcceptTag("number")
	symbolFa.SetAcceptTag("symbol")

	errContainer := types.NewErrorContainer()
	file, err := lexer.NewIOFromFile("lexer_test_1.txt")
	if err != nil {
		t.Fatal(err)
	}
	dfa := identityFa.MergeOr(numerFa).MergeOr(symbolFa).AsDFA()
	lex := &lexer.Lexer{
		ErrorContainer: *errContainer,
		Io:             file,
		DFA:            dfa,
	}

	dump := dfa.Dump()
	t.Log(dump)

	order := []string{
		"identify", // var
		"identify", // a
		"symbol",   // =
		"number",   // 5.1
		"symbol",   // ;
		"identify", // var
		"identify", // b
		"symbol",   // =
		"number",   // 7
		"symbol",   // ;
		"identify", // print
		"symbol",   // (
		"identify", // a
		"symbol",   // +
		"identify", // b
		"symbol",   // -
		"number",   // 0.15e
		"symbol",   // )
		"symbol",   // ;
	}

	for i := 0; i < len(order); i++ {
		token := lex.NextToken()
		if token == nil {
			t.FailNow()
		}
		if token.Tag != order[i] {
			t.Fail()
		}
	}
	if lex.NextToken() != nil {
		t.FailNow()
	}
	if len(lex.ErrorContainer.Errors) != 0 {
		t.Fail()
	}
}
