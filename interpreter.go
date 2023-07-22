package brainfucker

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type Interpreter struct {
	input      *bufio.Reader
	buf        Buffer
	progInput  io.ByteReader
	progOutput io.ByteWriter
}

func NewInterpreter(
	input *bufio.Reader,
	progInput io.ByteReader,
	progOutput io.ByteWriter,
	dataCellsBufferSize ...int,
) *Interpreter {
	bufSize := 30000
	if len(dataCellsBufferSize) != 0 {
		bufSize = dataCellsBufferSize[0]
	}
	return &Interpreter{
		input:      input,
		buf:        NewBuffer(bufSize),
		progInput:  progInput,
		progOutput: progOutput,
	}
}

func (i *Interpreter) buffer() []byte {
	return i.buf.raw()
}

func (i *Interpreter) currentSymbol() (byte, error) {
	b, err := i.input.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("error reading program: %w\n", err)
	}
	return b, nil
}

func (i *Interpreter) peekCurrentSymbol() (byte, error) {
	b, err := i.input.Peek(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil

}

func (i *Interpreter) currentCell() (byte, error) {
	return i.buf.Current()
}

func (i *Interpreter) Run() error {
	var curr byte
	var err error

	for curr, err = i.currentSymbol(); err == nil; curr, err = i.currentSymbol() {
		switch curr {
		case IncrementData:
			if err := i.buf.IncrementCurrent(); err != nil {
				return err
			}
		case DecrementData:
			if err := i.buf.DecrementCurrent(); err != nil {
				return err
			}
		case IncrementPointer:
			if err := i.buf.GoRight(); err != nil {
				return err
			}
		case DecrementPointer:
			if err := i.buf.GoLeft(); err != nil {
				return err
			}
		case OutputData:
			curr, err := i.buf.Current()
			if err != nil {
				return err
			}
			if err := i.progOutput.WriteByte(curr); err != nil {
				return fmt.Errorf("error writing program output: %w\n", err)
			}
		case InputData:
			b, err := i.progInput.ReadByte()
			if err != nil {
				return fmt.Errorf("error reading program input: %w\n", err)
			}
			if err := i.buf.WriteCurrent(b); err != nil {
				return err
			}
		case StartLoop:
			curr, err := i.currentCell()
			if err != nil {
				return err
			}
			if curr == 0 {
				for {
					curr, err := i.currentSymbol()
					if err != nil {
						return err
					}
					if curr == EndLoop {
						break
					}
				}
			}
		case EndLoop:
			curr, err := i.currentCell()
			if err != nil {
				return err
			}
			if curr != 0 {
				var s byte
				var err error
				for s, err = i.peekCurrentSymbol(); err == nil && s != StartLoop; s, err = i.peekCurrentSymbol() {
					if err := i.input.UnreadByte(); err != nil {
						return err
					}
				}
				if err != nil {
					return err
				}
			}

		}
	}
	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}
