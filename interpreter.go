package brainfucker

import (
	"bufio"
	"fmt"
	"io"
)

type Interpreter struct {
	input        []rune
	currInputPos int
	dataCells    []byte
	currCellPos  int
	loopStart    int
	progInput    io.Reader
	progOutput   io.Writer
}

func NewInterpreter(input io.Reader, output io.Writer, dataCellsBufferSize ...int) *Interpreter {
	bufSize := 30000
	if len(dataCellsBufferSize) != 0 {
		bufSize = dataCellsBufferSize[0]
	}
	dataCells := make([]byte, bufSize)
	return &Interpreter{dataCells: dataCells, currCellPos: 0, currInputPos: 0, progInput: input, progOutput: output}
}

func (i *Interpreter) Reset() {
	i.input = make([]rune, 0)
	i.currInputPos = 0
	i.dataCells = make([]byte, len(i.dataCells))
	i.currCellPos = 0
	i.loopStart = 0
}

func (i *Interpreter) readInput(reader io.Reader) error {
	s := bufio.NewReader(reader)
	input, err := s.ReadString(0)
	if err != nil && err != io.EOF {
		return err
	}
	i.input = []rune(input)
	return nil
}

func (i *Interpreter) currentRune() (rune, error) {
	if i.currInputPos >= len(i.input) {
		return 0, io.EOF
	}
	return i.input[i.currInputPos], nil
}

func (i *Interpreter) currentCell() byte {
	if i.currCellPos < 0 {
		i.currCellPos = 0
	}
	return i.dataCells[i.currCellPos]
}

func (i *Interpreter) Run(input io.Reader) error {
	err := i.readInput(input)
	if err != nil {
		return err
	}

	for i.currInputPos < len(i.input) {
		curr, err := i.currentRune()
		if err == io.EOF {
			break
		}
		switch curr {
		case IncrementData:
			i.dataCells[i.currCellPos]++
			break
		case DecrementData:
			i.dataCells[i.currCellPos]--
			break
		case IncrementPointer:
			i.currCellPos++
			break
		case DecrementPointer:
			i.currCellPos--
			break
		case OutputData:
			_, err := i.progOutput.Write([]byte{i.currentCell()})
			if err != nil {
				return err
			}
			break
		case InputData:
			_, err := fmt.Fscanf(i.progInput, "%c", &i.dataCells[i.currCellPos])
			if err != nil {
				return err
			}
			break
		case StartLoop:
			if i.currentCell() == 0 {
				for {
					i.currInputPos++
					curr, err := i.currentRune()
					if err != nil {
						return err
					}
					if curr == EndLoop {
						break
					}
				}
			} else {
				i.loopStart = i.currInputPos
			}
			break
		case EndLoop:
			if i.currentCell() != 0 {
				i.currInputPos = i.loopStart
			} else {
				i.loopStart = 0
			}
			break
		}
		i.currInputPos++
	}
	return nil
}
