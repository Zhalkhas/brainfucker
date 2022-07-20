package brainfucker

import (
	"bufio"
	"io"
)

type Compiler struct {
	buffSize     int
	input        []rune
	currInputPos int
}

func NewCompiler(dataCellsBufferSize ...int) *Compiler {
	buffSize := 30000
	if len(dataCellsBufferSize) != 0 {
		buffSize = dataCellsBufferSize[0]
	}
	return &Compiler{buffSize: buffSize}
}

func (c *Compiler) Compile(input io.Reader) error {
	bufio.NewReader(input).ReadByte()
	return nil
}
