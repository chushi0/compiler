package lexer

import (
	"bufio"
	"os"
)

// IO
// 与文件交互的接口
type IO struct {
	Filename        string        // 文件名
	File            *os.File      // 文件
	BufferReader    *bufio.Reader // 缓冲输入流
	RuneBuffer      []rune        // 字符缓存（缓存越大，回滚时更有可能减小io）
	RuneBufferIndex int           // 当前读取的字符缓存位置
	RuneValidCount  int           // 字符缓存有效大小
	IgnoreLF        bool          // 忽略一次`\n`

	Line   int   // 当前行号
	Column int   // 当前列
	Offset int64 // 当前偏移（字节，相对于文件）

	SaveLine     int   // 保存的行号
	SaveColumn   int   // 保存的列
	SaveOffset   int64 // 保存的偏移
	SaveIndex    int   // 保存的字符缓存读取位置
	SaveIgnoreLF bool  // 保存是否忽略`\n`
}

// 有穷自动状态机
// 包含NFA和DFA
type FiniteAutomaton struct {
	StateCount     int            // 状态数
	JumpTables     [][]*JumpMap   // 转移函数
	AcceptStates   []int          // 接受的状态数
	AcceptStateTag map[int]string // 接受的状态标签
}

type JumpMap struct {
	// 字符范围开始与字符范围结束均为 0 时，表示无字符跳转
	// 仅在 NFA 中有效
	RuneStart rune // 字符范围开始（包含）
	RuneEnd   rune // 字符范围结束（不包含）
	Target    int  // 跳转到的状态
}

// 词法分析器
type Lexer struct {
	Io    *IO              // 文件 IO 接口
	DFA   *FiniteAutomaton // DFA
	State int              // DFA 当前状态
}

// 词法分析器读到的内容
type Token struct {
	RawValue string // 原始值
	State    int    // 接受的状态
	Tag      string // tag
}
