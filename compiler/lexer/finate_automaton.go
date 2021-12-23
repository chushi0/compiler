package lexer

// 两个非确定有穷自动机作或运算 a|b
// 合并初始状态
func (fa *FiniteAutomaton) MergeOr(o *FiniteAutomaton) *FiniteAutomaton {
	stateCount := fa.StateCount + o.StateCount - 1
	result := &FiniteAutomaton{
		StateCount:     stateCount,
		JumpTables:     make([][]*JumpMap, stateCount),
		AcceptStates:   make([]int, 0),
		AcceptStateTag: make(map[int]string),
	}

	// 填充 JumpTables
	for i := 0; i < fa.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneStart: jumpMap.RuneStart,
				RuneEnd:   jumpMap.RuneEnd,
				Target:    jumpMap.Target,
			})
		}
		result.JumpTables[i] = jumpTable
	}
	for _, jumpMap := range o.JumpTables[0] {
		result.JumpTables[0] = append(result.JumpTables[0], &JumpMap{
			RuneStart: jumpMap.RuneStart,
			RuneEnd:   jumpMap.RuneEnd,
			Target:    jumpMap.Target + fa.StateCount - 1,
		})
	}
	for i := 1; i < o.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneStart: jumpMap.RuneStart,
				RuneEnd:   jumpMap.RuneEnd,
				Target:    jumpMap.Target + fa.StateCount - 1,
			})
		}
		result.JumpTables[i+fa.StateCount-1] = jumpTable
	}

	// 填充 AcceptStates
	result.AcceptStates = append(result.AcceptStates, fa.AcceptStates...)
	if o.AcceptStates[0] == 0 {
		if fa.AcceptStates[0] == 0 {
			result.AcceptStates = append([]int{0}, result.AcceptStates...)
			for _, state := range o.AcceptStates[1:] {
				result.AcceptStates = append(result.AcceptStates, state+fa.StateCount-1)
			}
		}
	} else {
		for _, state := range o.AcceptStates {
			result.AcceptStates = append(result.AcceptStates, state+fa.StateCount-1)
		}
	}

	// 填充 Tags
	for state, tag := range fa.AcceptStateTag {
		result.AcceptStateTag[state] = tag
	}
	for state, tag := range o.AcceptStateTag {
		if state == 0 {
			result.AcceptStateTag[0] = tag
		} else {
			result.AcceptStateTag[state+fa.StateCount-1] = tag
		}
	}

	return result
}
