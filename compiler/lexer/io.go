package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

func NewIOFromFile(filepath string) (*IO, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file \"%s\": %w", filepath, &IOError{Original: err})
	}
	return &IO{
		Filename:        filepath,
		File:            file,
		BufferReader:    bufio.NewReader(file),
		RuneBuffer:      make([]rune, 4096),
		RuneBufferIndex: 0,
		RuneValidCount:  0,
		IgnoreLF:        false,
		Line:            1,
		Column:          0,
		Offset:          0,
	}, nil
}

func (pio *IO) Close() error {
	if pio.File == nil {
		return nil
	}
	err := pio.File.Close()
	if err != nil {
		return fmt.Errorf("close io fail: %w", err)
	}
	pio.File = nil
	return nil
}

func (pio *IO) ReadChar() (rune, error) {
	if pio.RuneBufferIndex >= pio.RuneValidCount {
		err := pio.fillBuffer()
		if err != nil {
			return 0, err
		}
	}
	rn := pio.RuneBuffer[pio.RuneBufferIndex]
	pio.RuneBufferIndex++
	pio.Offset += int64(utf8.RuneLen(rn))
	if rn == '\r' {
		pio.IgnoreLF = true
		pio.Line += 1
		pio.Column = 0
		return '\n', nil
	}
	if rn == '\n' {
		if pio.IgnoreLF {
			pio.IgnoreLF = false
			return pio.ReadChar()
		}
		pio.Line += 1
		pio.Column = 0
	} else {
		pio.IgnoreLF = false
		pio.Column += 1
	}
	return rn, nil
}

func (pio *IO) Lookup() (rune, error) {
	if pio.RuneBufferIndex >= pio.RuneValidCount {
		err := pio.fillBuffer()
		if err != nil {
			return 0, err
		}
	}
	rn := pio.RuneBuffer[pio.RuneBufferIndex]
	if rn == '\r' {
		return '\n', nil
	}
	return rn, nil
}

func (pio *IO) Save() {
	pio.SaveLine = pio.Line
	pio.SaveColumn = pio.Column
	pio.SaveOffset = pio.Offset
	pio.SaveIndex = pio.RuneBufferIndex
	pio.SaveIgnoreLF = pio.IgnoreLF
}

func (pio *IO) Restore() error {
	if pio.SaveIndex >= 0 {
		pio.RuneBufferIndex = pio.SaveIndex
		pio.Line = pio.SaveLine
		pio.Column = pio.SaveColumn
		pio.Offset = pio.SaveOffset
		pio.IgnoreLF = pio.SaveIgnoreLF
		return nil
	}
	_, err := pio.File.Seek(pio.SaveOffset, 0)
	if err != nil {
		return &IOError{Original: err}
	}
	pio.BufferReader.Reset(pio.File)
	pio.Line = pio.SaveLine
	pio.Column = pio.SaveColumn
	pio.Offset = pio.SaveOffset
	pio.IgnoreLF = pio.SaveIgnoreLF
	pio.RuneBufferIndex = 0
	pio.RuneValidCount = 0
	return nil
}

func (pio *IO) fillBuffer() error {
	fillStart := 0
	if pio.SaveIndex == 0 {
		pio.SaveIndex = -1
	}
	if pio.SaveIndex > 0 {
		fillStart = pio.RuneValidCount - pio.SaveIndex
		for i := 0; i < fillStart; i++ {
			pio.RuneBuffer[i] = pio.RuneBuffer[pio.SaveIndex+i]
		}
		pio.SaveIndex = 0
	}
	i := fillStart
	for i < len(pio.RuneBuffer) {
		rn, _, err := pio.BufferReader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		pio.RuneBuffer[i] = rn
		i++
	}
	pio.RuneValidCount = i
	if pio.RuneValidCount == 0 {
		return io.EOF
	}
	pio.RuneBufferIndex = fillStart
	return nil
}
