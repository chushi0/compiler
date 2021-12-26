package parser

func (fa *LRFinateAutomaton) Parse(reader LRReader) Symbol {
	ctx := &LRContext{
		StateStack:  make([]int, 0),
		SymbolStack: make([]Symbol, 0),
		Reader:      reader,
	}
	ctx.StateStack = append(ctx.StateStack, 0)
	ctx.SymbolStack = append(ctx.SymbolStack, EOF)
	for {
		token := reader.LookupNextToken()
		if token == nil {
			token = EOF
		}
		action := fa.ActionTables[ctx.StateStack[len(ctx.StateStack)-1]][token.Name]
		switch action.Type {
		case LRActionType_Accept:
			return ctx.SymbolStack[len(ctx.SymbolStack)-1]
		case LRActionType_Shift:
			ctx.StateStack = append(ctx.StateStack, action.Index)
			ctx.SymbolStack = append(ctx.SymbolStack, token)
			reader.SkipToken()
		case LRActionType_Reduce:
			count := len(fa.Grammar.Productions[action.Index]) - 1
			reduceFunc := fa.ReduceFunc[action.Index]
			symbol := reduceFunc(ctx.SymbolStack[len(ctx.SymbolStack)-count:])
			ctx.SymbolStack = append(ctx.SymbolStack[:len(ctx.SymbolStack)], symbol)
			ctx.StateStack = append(ctx.StateStack[:len(ctx.StateStack)-count], fa.GotoTables[ctx.StateStack[len(ctx.StateStack)-count-1]][symbol.SymbolName()])
		case LRActionType_Error:
			errFunc := fa.ErrorProcessFunc[ctx.StateStack[len(ctx.StateStack)-1]]
			errFunc(ctx)
		}
	}
}
