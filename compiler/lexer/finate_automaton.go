package lexer

import (
	"sort"

	"github.com/chushi0/compiler/utils/set"
	utilslice "github.com/chushi0/compiler/utils/util_slice"
)

// 检查是否是ε转换
func (r *RuneRange) isEpsilon() bool {
	return r.RuneStart == 0 && r.RuneEnd == 0
}

// 检查是否包含指定集合
func (r *RuneRange) contains(o *RuneRange) bool {
	return r.RuneStart <= o.RuneStart && r.RuneEnd >= o.RuneEnd
}

// 检查是否相交
func (r *RuneRange) isIntersect(o *RuneRange) bool {
	// 若两者相等，只能是连续
	return r.RuneEnd > o.RuneStart && o.RuneEnd > r.RuneStart
}

// 分割相交的区域
func (r *RuneRange) splitWith(o *RuneRange) []*RuneRange {
	var r1, r2, r3 RuneRange
	if r.RuneStart < o.RuneStart {
		r1.RuneStart = r.RuneStart
		r1.RuneEnd = o.RuneStart
		r2.RuneStart = o.RuneStart
	} else if r.RuneStart > o.RuneStart {
		r1.RuneStart = o.RuneStart
		r1.RuneEnd = r.RuneStart
		r2.RuneStart = r.RuneStart
	} else {
		r2.RuneStart = r.RuneStart
	}
	if r.RuneEnd < o.RuneEnd {
		r3.RuneStart = r.RuneEnd
		r3.RuneEnd = o.RuneEnd
		r2.RuneEnd = r.RuneEnd
	} else if r.RuneEnd > o.RuneEnd {
		r3.RuneStart = o.RuneEnd
		r3.RuneEnd = r.RuneEnd
		r2.RuneEnd = o.RuneEnd
	} else {
		r2.RuneEnd = r.RuneEnd
	}

	result := make([]*RuneRange, 0)
	if !r1.isEpsilon() {
		result = append(result, &r1)
	}
	result = append(result, &r2)
	if !r3.isEpsilon() {
		result = append(result, &r3)
	}
	return result
}

// 两个非确定有穷自动机作或运算 a|b
// 合并初始状态
func (fa *FiniteAutomaton) MergeOr(o *FiniteAutomaton) *FiniteAutomaton {
	stateCount := fa.StateCount + o.StateCount - 1
	result := &FiniteAutomaton{
		StateCount:     stateCount,
		JumpTables:     make([][]*JumpMap, stateCount),
		AcceptStateTag: make(map[int]string),
	}

	// 填充 JumpTables
	for i := 0; i < fa.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneRange: jumpMap.RuneRange,
				Target:    jumpMap.Target,
			})
		}
		result.JumpTables[i] = jumpTable
	}
	for _, jumpMap := range o.JumpTables[0] {
		result.JumpTables[0] = append(result.JumpTables[0], &JumpMap{
			RuneRange: jumpMap.RuneRange,
			Target:    jumpMap.Target + fa.StateCount - 1,
		})
	}
	for i := 1; i < o.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneRange: jumpMap.RuneRange,
				Target:    jumpMap.Target + fa.StateCount - 1,
			})
		}
		result.JumpTables[i+fa.StateCount-1] = jumpTable
	}

	// 填充 AcceptStates
	result.AcceptStates = fa.AcceptStates.Clone()
	for state := range o.AcceptStates {
		if state == 0 {
			result.AcceptStates.Put(0)
		}
		result.AcceptStates.Put(state + fa.StateCount - 1)
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

// 两个非确定有穷自动机作连接运算 a+b
// 将前者的接受状态连接后者的初始状态
func (fa *FiniteAutomaton) MergeConnect(o *FiniteAutomaton) *FiniteAutomaton {
	stateCount := fa.StateCount + o.StateCount
	result := &FiniteAutomaton{
		StateCount:     stateCount,
		JumpTables:     make([][]*JumpMap, stateCount),
		AcceptStateTag: make(map[int]string),
	}

	// 填充 JumpTables
	for i := 0; i < fa.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneRange: jumpMap.RuneRange,
				Target:    jumpMap.Target,
			})
		}
		result.JumpTables[i] = jumpTable
	}
	for i := 0; i < o.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneRange: jumpMap.RuneRange,
				Target:    jumpMap.Target + fa.StateCount,
			})
		}
		result.JumpTables[i+fa.StateCount] = jumpTable
	}
	for state := range fa.AcceptStates {
		result.JumpTables[state] = append(result.JumpTables[state], &JumpMap{
			RuneRange: RuneRange{
				RuneStart: 0,
				RuneEnd:   0,
			},
			Target: fa.StateCount,
		})
	}

	// 填充 AcceptStates
	for state := range o.AcceptStates {
		result.AcceptStates.Put(state + fa.StateCount)
	}

	// 填充 AcceptStateTag
	for state := range o.AcceptStates {
		result.AcceptStateTag[state+fa.StateCount] = o.AcceptStateTag[state]
	}
	return result
}

// 非确定有穷自动机作克林闭包 a*
// 将接受状态连接初始状态，然后将初始状态设置为接受状态（清除其他接受状态）
func (fa *FiniteAutomaton) MergeKleene() *FiniteAutomaton {
	result := &FiniteAutomaton{
		StateCount:     fa.StateCount,
		JumpTables:     make([][]*JumpMap, fa.StateCount),
		AcceptStates:   set.NewIntSet(0),
		AcceptStateTag: make(map[int]string),
	}
	for i := 0; i < fa.StateCount; i++ {
		jumpTable := make([]*JumpMap, 0)
		for _, jumpMap := range fa.JumpTables[i] {
			jumpTable = append(jumpTable, &JumpMap{
				RuneRange: jumpMap.RuneRange,
				Target:    jumpMap.Target,
			})
		}
		result.JumpTables[i] = jumpTable
	}
	for state := range fa.AcceptStates {
		result.JumpTables[state] = append(result.JumpTables[state], &JumpMap{
			RuneRange: RuneRange{
				RuneStart: 0,
				RuneEnd:   0,
			},
			Target: 0,
		})
	}
	for state := range fa.AcceptStates {
		result.AcceptStateTag[0] = fa.AcceptStateTag[state]
		break
	}
	return result
}

// 为接受状态设置标记
func (fa *FiniteAutomaton) SetAcceptTag(tag string) *FiniteAutomaton {
	for state := range fa.AcceptStates {
		fa.AcceptStateTag[state] = tag
	}
	return fa
}

// 单字符匹配的有穷自动机
func NewFinateAutomaton(runeRange *RuneRange) *FiniteAutomaton {
	return &FiniteAutomaton{
		StateCount: 2,
		JumpTables: [][]*JumpMap{
			{
				&JumpMap{
					RuneRange: *runeRange,
					Target:    1,
				},
			},
		},
		AcceptStates: set.NewIntSet(1),
		AcceptStateTag: map[int]string{
			1: "",
		},
	}
}

// 能够从 NFA 的指定状态只通过ε转换到达的 NFA 状态的集合
func (fa *FiniteAutomaton) closure(state int) set.IntSet {
	res := set.NewIntSet()
	for state > 0 {
		res.Put(state)
		table := fa.JumpTables[state]
		state = -1
		for _, jumpMap := range table {
			if jumpMap.isEpsilon() {
				state = jumpMap.Target
				break
			}
		}
	}
	return res
}

func (fa *FiniteAutomaton) closureSet(states set.IntSet) set.IntSet {
	result := set.NewIntSet()
	for state := range states {
		result.Merge(fa.closure(state))
	}
	return result
}

// 分割所有字符范围
func (fa *FiniteAutomaton) splitRange() []*RuneRange {
	ranges := make([]*RuneRange, 0)
	for _, jumpTable := range fa.JumpTables {
		for _, jumpMap := range jumpTable {
			if jumpMap.isEpsilon() {
				continue
			}
			done := false
			for i, rng := range ranges {
				if jumpMap.RuneRange == *rng {
					done = true
					break
				}
				if jumpMap.isIntersect(rng) {
					last := len(ranges) - 1
					ranges[i], ranges[last] = ranges[last], ranges[i]
					ranges = append(ranges[:last], jumpMap.splitWith(ranges[last])...)
					done = true
					break
				}
			}
			if !done {
				ranges = append(ranges, &RuneRange{
					RuneStart: jumpMap.RuneStart,
					RuneEnd:   jumpMap.RuneEnd,
				})
			}
		}
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].RuneStart < ranges[j].RuneStart
	})
	return ranges
}

// NFA 转 DFA
func (fa *FiniteAutomaton) AsDFA() *FiniteAutomaton {
	result := &FiniteAutomaton{
		JumpTables:     make([][]*JumpMap, 0),
		AcceptStates:   set.NewIntSet(),
		AcceptStateTag: make(map[int]string),
	}
	// 拆分字符范围
	ranges := fa.splitRange()
	// 状态
	states := make([]set.IntSet, 0)
	states = append(states, fa.closure(0))
	for i := 0; i < len(states); i++ {
		curState := states[i]
		// 检查当前状态集是否可以接受
		for state := range curState {
			if fa.AcceptStates.Contains(state) {
				result.AcceptStates.Put(state)
				result.AcceptStateTag[i] = fa.AcceptStateTag[state]
				break
			}
		}
		// 计算转移函数和更多状态
		jumpTable := make([]*JumpMap, 0)
		for _, rng := range ranges {
			moveTo := set.NewIntSet()
			for state := range curState {
				for _, jumpMap := range fa.JumpTables[state] {
					if jumpMap.isEpsilon() {
						continue
					}
					if jumpMap.contains(rng) {
						moveTo.Put(jumpMap.Target)
					}
				}
			}
			if len(moveTo) == 0 {
				continue
			}
			moveTo = fa.closureSet(moveTo)
			index := utilslice.LinearSearch(len(states), func(i int) bool {
				return states[i].Equals(moveTo)
			})
			if index == -1 {
				states = append(states, moveTo)
				index = len(states) - 1
			}
			jumpTable = append(jumpTable, &JumpMap{
				RuneRange: *rng,
				Target:    index,
			})
		}
		result.JumpTables = append(result.JumpTables, jumpTable)
	}
	result.StateCount = len(states)
	return result
}
