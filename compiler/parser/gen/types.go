package gen

import (
	"github.com/chushi0/compiler/compiler/parser"
	"github.com/chushi0/compiler/utils/set"
)

// 动作
// 设置为数组方便表示冲突
// （移入-归约冲突和归约-归约冲突）
type LRAction []*parser.LRAction

// 项目
type Item struct {
	Production parser.Production // 产生式
	State      int               // 状态
	Lookahead  set.StringSet     // 展望符
}

// 跳转表
type JumpTable struct {
	Items set.HashSet         // 项目集
	Jump  map[string]LRAction // 跳转信息（遇到非终结符或终结符后进入哪个状态）
}
