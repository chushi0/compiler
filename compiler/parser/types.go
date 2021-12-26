package parser

import (
	"github.com/chushi0/compiler/compiler/lexer"
)

// 语法分析器输入 token
type Token struct {
	Name  string      // 所属终结符类型
	Token lexer.Token // 词法分析器输出 token
}

type UnitType uint

const (
	UnitType_Terminal    UnitType = 1 // 终结符
	UnitType_Nonterminal UnitType = 2 // 非终结符
)

// 产生式
// 第一个为左部，其它的为右部
type Production []string

// 文法
type Grammar struct {
	Productions []Production // 所有文法
	Start       string       // 开始符号
}

// LR 分析器
type LRFinateAutomaton struct {
	Grammar          *Grammar                      // 产生式
	StateCount       int                           // 状态数
	ActionTables     []map[string]*LRAction        // Action 跳转表
	GotoTables       []map[string]int              // Goto 跳转表
	ReduceFunc       map[int]func([]Symbol) Symbol // 归约函数
	ErrorProcessFunc map[int]func(ctx *LRContext)  // 错误处理函数
}

type LRContext struct {
	StateStack  []int    // 状态栈
	SymbolStack []Symbol // 符号栈
	Reader      LRReader // 输入
}

type LRActionType uint

const (
	LRActionType_Accept LRActionType = 1 // 接受
	LRActionType_Shift  LRActionType = 2 // 移入
	LRActionType_Reduce LRActionType = 3 // 归约
	LRActionType_Error  LRActionType = 4 // 错误处理
)

type LRAction struct {
	Type  LRActionType // 动作类型
	Index int          // 参数
}

var (
	EOF *Token
)

type Symbol interface {
	SymbolName() string
}

type LRReader interface {
	LookupNextToken() *Token // 检查下一个 token
	SkipToken()              // 跳过 token
}
