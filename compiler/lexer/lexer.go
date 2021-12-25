package lexer

import (
	"errors"
	"fmt"
	"io"

	"github.com/chushi0/compiler/compiler/types"
)

// 读取下一个 Token
func (lexer *Lexer) NextToken() *Token {
scan_start:
	// 初始化
	err := lexer.clearSpace()
	if err != nil {
		lexer.ErrorContainer.Fatal = append(lexer.ErrorContainer.Fatal, &types.Error{
			Type:   types.ErrorType_SystemFileError,
			File:   lexer.Io.Filename,
			Detail: fmt.Sprintf("unexcepted error: %s", err.Error()),
		})
		return nil
	}
	token := &Token{
		Line:   lexer.Io.Line,
		Column: lexer.Io.Column,
		File:   lexer.Io.Filename,
		State:  -1,
	}
	rawVal := ""
	state := 0

	// 循环读取字符
	// 贪心模式匹配尽可能多的字符
	for {
		rn, err := lexer.Io.ReadChar()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			lexer.ErrorContainer.Fatal = append(lexer.ErrorContainer.Fatal, &types.Error{
				Type:   types.ErrorType_SystemFileError,
				File:   lexer.Io.Filename,
				Detail: fmt.Sprintf("unexcepted error: %s", err.Error()),
			})
			return nil
		}
		state, err := lexer.DFA.NextState(state, rn)
		if err != nil {
			break
		}
		rawVal += string(rn)
		if lexer.DFA.AcceptStates.Contains(state) {
			lexer.Io.Save()
			token.RawValue = rawVal
			token.State = state
			token.Tag = lexer.DFA.AcceptStateTag[state]
		}
	}

	// 曾经匹配成功，恢复io后返回
	if token.State != -1 {
		lexer.Io.Restore()
		return token
	}

	// 未匹配成功过，报错后继续扫描
	lexer.ErrorContainer.Errors = append(lexer.ErrorContainer.Errors, &types.Error{
		Type:   types.ErrorType_UnexpectedToken,
		File:   lexer.Io.Filename,
		Line:   token.Line,
		Column: token.Column,
		Detail: fmt.Sprintf("unexpected token: %s", rawVal),
	})
	goto scan_start
}

// 清除空白字符
func (lexer *Lexer) clearSpace() error {
	for {
		rn, err := lexer.Io.Lookup()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if !isSpace(rn) {
			return nil
		}
		_, err = lexer.Io.ReadChar()
		if err != nil {
			return err
		}
	}
}

func isSpace(rn rune) bool {
	switch rn {
	case '\t', '\r', '\n', '\v', '\f', ' ':
		return true
	default:
		return false
	}
}
