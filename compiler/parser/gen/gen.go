package gen

import (
	"github.com/chushi0/compiler/compiler/parser"
	"github.com/chushi0/compiler/utils/set"
)

func GenerateLRFinateAutomaton(grammar *parser.Grammar) *parser.LRFinateAutomaton {
	// 抽取新的开始符号
	startProduction := parser.Production{
		"LR_Start_Grammar",
		grammar.Start,
	}
	grammar.Productions = append(grammar.Productions, startProduction)
	grammar.Start = startProduction[0]
	// 统计非终结符
	nonTerminals := set.NewStringSet()
	for _, production := range grammar.Productions {
		nonTerminals.Put(production[0])
	}
	// 计算非终结符的 First 集
	firstSet := make(map[string]set.StringSet)
	for nonTerminal := range nonTerminals {
		firstSet[nonTerminal] = set.NewStringSet()
	}
	needRecompute := true
	for needRecompute {
		needRecompute = false
		for _, production := range grammar.Productions {
			mergeEmpty := true
			for i := 1; i < len(production) && mergeEmpty; i++ {
				mergeEmpty = false
				if production.TypeAt(i) == parser.UnitType_Terminal {
					if !firstSet[production[0]].Contains(production[i]) {
						needRecompute = true
						firstSet[production[0]].Put(production[i])
					}
					break
				}
				for first := range firstSet[production[i]] {
					if first == "" {
						mergeEmpty = true
						continue
					}
					if !firstSet[production[0]].Contains(first) {
						needRecompute = true
						firstSet[production[0]].Put(first)
					}
				}
			}
			if mergeEmpty && !firstSet[production[0]].Contains("") {
				needRecompute = true
				firstSet[production[0]].Put("")
			}
		}
	}

	// 计算各状态及跳转表
	jumpTables := make([]*JumpTable, 0)
	jumpTables = append(jumpTables, &JumpTable{
		Items: closure(grammar, &Item{
			Production: startProduction,
			State:      0,
			Lookahead:  set.NewStringSet("$"),
		}),
		Jump: make(map[string]LRAction),
	})
	for i := 0; i < len(jumpTables); i++ {
		jumpTables[i].Items.Foreach(func(sk set.SetKey) {
			item := sk.(*Item)
			// 归约（不管）
			if item.State+1 == len(item.Production) {
				return
			}
			// 移入、归约操作
			// 展望符不变
			newItem := &Item{
				Production: item.Production,
				State:      item.State + 1,
				Lookahead:  item.Lookahead.Clone(),
			}
			itemClosure := closure(grammar, newItem)
			index := -1
			for i, jt := range jumpTables {
				if jt.Items.Equals(itemClosure) {
					index = i
				}
			}
			if index == -1 {
				jumpTables = append(jumpTables, &JumpTable{
					Items: itemClosure,
					Jump:  make(map[string]LRAction),
				})
				index = len(jumpTables) - 1
			}
			actType := parser.LRActionType_Shift
			if item.Production.TypeAt(item.State+1) == parser.UnitType_Nonterminal {
				actType = parser.LRActionType_Reduce
				if item.Production[0] == startProduction[0] {
					actType = parser.LRActionType_Accept
				}
			}
			if _, ok := jumpTables[i].Jump[item.Production[item.State+1]]; !ok {
				jumpTables[i].Jump[item.Production[item.State+1]] = make(LRAction, 0)
			}
			jumpTables[i].Jump[item.Production[item.State+1]] = append(jumpTables[i].Jump[item.Production[item.State+1]], &parser.LRAction{
				Type:  actType,
				Index: index,
			})
		})

	}

	return nil
}

func closure(grammar *parser.Grammar, item *Item) set.HashSet {
	return nil
}
