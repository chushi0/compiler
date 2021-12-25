package types

const (
	// 严重错误：10000~19999
	// 编译错误：20000~29999
	// 编译警告：30000~39999

	// 文件读取错误
	ErrorType_SystemFileError ErrorType = 10001

	// 词法分析器无法解析出正确的单词
	ErrorType_UnexpectedToken ErrorType = 20001
)
