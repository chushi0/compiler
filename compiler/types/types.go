package types

type ErrorContainer struct {
	Fatal    []*Error
	Errors   []*Error
	Warnings []*Error
}

type ErrorType uint

type Error struct {
	Type   ErrorType // 错误类型
	File   string    // 所属文件
	Line   int       // 行号
	Column int       // 列号
	Detail string    // 详细信息
}

func NewErrorContainer() *ErrorContainer {
	return &ErrorContainer{
		Errors:   make([]*Error, 0),
		Warnings: make([]*Error, 0),
	}
}
