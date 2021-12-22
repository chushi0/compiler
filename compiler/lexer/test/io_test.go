package lexer_test

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/chushi0/compiler/compiler/lexer"
)

const TestFileName = "io_test_file.txt"

// 测试正常情况下是否可以正确读取文件
func TestIO_NormalReadFile(t *testing.T) {
	content, err := ioutil.ReadFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	fileContent := strings.ReplaceAll(string(content), "\r\n", "\n")
	pio, err := lexer.NewIOFromFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer pio.Close()
	rns := make([]rune, 0)
	for {
		rn, err := pio.ReadChar()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}
		rns = append(rns, rn)
	}

	if fileContent != string(rns) {
		t.Fail()
	}
}

// 测试正常情况下的行号、列号
func TestIO_LineColumn(t *testing.T) {
	pio, err := lexer.NewIOFromFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer pio.Close()

	// 读取15个字符后，Line=1，Column=15
	for i := 0; i < 15; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	if pio.Line != 1 || pio.Column != 15 {
		t.Fail()
	}

	// 再读取15个字符后，Line=2，Column=3
	for i := 0; i < 15; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	if pio.Line != 2 || pio.Column != 3 {
		t.Fail()
	}

	// 再读取30个字符后，Line=3，Column=5
	for i := 0; i < 30; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	if pio.Line != 3 || pio.Column != 6 {
		t.Fail()
	}
}

// 测试回滚操作
func TestIO_SaveRestore(t *testing.T) {
	pio, err := lexer.NewIOFromFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer pio.Close()

	// 读30个字符
	for i := 0; i < 30; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 保存
	pio.Save()
	// 再读30个
	for i := 0; i < 30; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 恢复
	pio.Restore()
	// 这时应该在 Line=2，Column=3
	if pio.Line != 2 || pio.Column != 3 {
		t.Fail()
	}
	// 读取的下一个字符应该是 'D'
	rn, err := pio.ReadChar()
	if err != nil {
		t.Fatal(err)
	}
	if rn != 'D' {
		t.Fail()
	}
}

// 测试回滚操作
// 保存后触发io，恢复时不应该触发
func TestIO_SaveRestore2(t *testing.T) {
	pio, err := lexer.NewIOFromFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer pio.Close()

	// 读30个字符
	for i := 0; i < 30; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 保存
	pio.Save()
	// 再读3960个
	// 触发io，但恢复时不会重新读取
	for i := 0; i < 3960; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 检查触发io
	if pio.SaveIndex != 0 {
		t.Fail()
	}
	// 恢复
	pio.Restore()
	// 这时应该在 Line=2，Column=3
	if pio.Line != 2 || pio.Column != 3 {
		t.Fail()
	}
	// 缓冲区有效
	if pio.RuneValidCount == 0 {
		t.Fail()
	}
	// 读取的下一个字符应该是 'D'
	rn, err := pio.ReadChar()
	if err != nil {
		t.Fatal(err)
	}
	if rn != 'D' {
		t.Fail()
	}
}

// 测试回滚操作
// 保存后触发两次io，恢复时触发io
func TestIO_SaveRestore3(t *testing.T) {
	pio, err := lexer.NewIOFromFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer pio.Close()

	// 读30个字符
	for i := 0; i < 30; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 保存
	pio.Save()
	// 再读4000个
	// 触发io，但恢复时不会重新读取
	for i := 0; i < 4000; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 检查触发io
	if pio.SaveIndex != -1 {
		t.Fail()
	}
	// 恢复
	pio.Restore()
	// 这时应该在 Line=2，Column=3
	if pio.Line != 2 || pio.Column != 3 {
		t.Fail()
	}
	// 缓冲区无效
	if pio.RuneValidCount != 0 {
		t.Fail()
	}
	// 读取的下一个字符应该是 'D'
	rn, err := pio.ReadChar()
	if err != nil {
		t.Fatal(err)
	}
	if rn != 'D' {
		t.Fail()
	}
}

// 忽略换行回滚测试
func TestIO_CRLF(t *testing.T) {
	pio, err := lexer.NewIOFromFile(TestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer pio.Close()

	// 读27个字符
	for i := 0; i < 27; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 此时应该等待换行符
	if !pio.IgnoreLF {
		t.Fail()
	}
	// 保存
	pio.Save()
	// 再读40个
	for i := 0; i < 40; i++ {
		_, err := pio.ReadChar()
		if err != nil {
			t.Fatal(err)
		}
	}
	// 恢复
	pio.Restore()
	// 这时应该在 Line=2，Column=0
	if pio.Line != 2 || pio.Column != 0 {
		t.Fail()
	}
	// 此时应该等待换行符
	if !pio.IgnoreLF {
		t.Fail()
	}
}
