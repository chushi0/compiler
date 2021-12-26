package parser_test

import (
	"testing"

	"github.com/chushi0/compiler/compiler/parser"
)

type SimpleInput struct {
	Inputs []string
	Index  int
}

func (arr *SimpleInput) LookupNextToken() *parser.Token {
	if arr.Index >= len(arr.Inputs) {
		return nil
	}
	return &parser.Token{
		Name: arr.Inputs[arr.Index],
	}
}

func (arr *SimpleInput) SkipToken() {
	arr.Index++
}

func TestLR1(t *testing.T) {
	grammar := &parser.Grammar{
		Productions: []parser.Production{
			{"S", "B", "B"},
			{"B", "a", "B"},
			{"B", "b"},
		},
		Start: "S",
	}
	actionTable := []map[string]*parser.LRAction{
		{
			"a": &parser.LRAction{
				Type:  parser.LRActionType_Shift,
				Index: 3,
			},
			"b": &parser.LRAction{
				Type:  parser.LRActionType_Shift,
				Index: 4,
			},
		},
		{
			"$": &parser.LRAction{
				Type: parser.LRActionType_Accept,
			},
		},
		{
			"a": &parser.LRAction{
				Type:  parser.LRActionType_Shift,
				Index: 3,
			},
			"b": &parser.LRAction{
				Type:  parser.LRActionType_Shift,
				Index: 4,
			},
		},
		{
			"a": &parser.LRAction{
				Type:  parser.LRActionType_Shift,
				Index: 3,
			},
			"b": &parser.LRAction{
				Type:  parser.LRActionType_Shift,
				Index: 4,
			},
		},
		{
			"a": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 2,
			},
			"b": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 2,
			},
			"$": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 2,
			},
		},
		{
			"a": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 0,
			},
			"b": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 0,
			},
			"$": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 0,
			},
		},
		{
			"a": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 1,
			},
			"b": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 1,
			},
			"$": &parser.LRAction{
				Type:  parser.LRActionType_Reduce,
				Index: 1,
			},
		},
	}
	gotoTable := []map[string]int{
		{
			"S": 1,
			"B": 2,
		},
		{},
		{
			"B": 5,
		},
		{
			"B": 6,
		},
		{},
		{},
		{},
	}
	reduceFunction := map[int]func([]parser.Symbol) parser.Symbol{
		0: func(s []parser.Symbol) parser.Symbol {
			t.Log("S->BB")
			return &parser.Token{
				Name: "S",
			}
		},
		1: func(s []parser.Symbol) parser.Symbol {
			t.Log("B->aB")
			return &parser.Token{
				Name: "B",
			}
		},
		2: func(s []parser.Symbol) parser.Symbol {
			t.Log("B->b")
			return &parser.Token{
				Name: "B",
			}
		},
	}

	lrfa := parser.LRFinateAutomaton{
		Grammar:          grammar,
		StateCount:       7,
		ActionTables:     actionTable,
		GotoTables:       gotoTable,
		ReduceFunc:       reduceFunction,
		ErrorProcessFunc: map[int]func(ctx *parser.LRContext){},
	}

	input := &SimpleInput{
		Inputs: []string{
			"b",
			"a",
			"b",
		},
		Index: 0,
	}

	result := lrfa.Parse(input)
	t.Log(result)
	t.Fail()
}
