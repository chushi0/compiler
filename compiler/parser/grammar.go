package parser

func init() {
	EOF = &Token{
		Name: "$",
	}
}

func (p Production) TypeAt(i int) UnitType {
	rn := p[i][0]
	if rn >= 'a' && rn <= 'z' {
		return UnitType_Terminal
	}
	return UnitType_Nonterminal
}

func (t *Token) SymbolName() string {
	return t.Name
}
