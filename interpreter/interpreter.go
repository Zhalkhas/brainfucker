package interpreter

import (
	"brainfucker/ast"
	"brainfucker/lexer"
	"bufio"
	"fmt"
	"io"
	"log"
)

type Interpreter struct {
	input  io.ByteReader
	output io.ByteWriter
}

func NewInterpreter(input io.Reader, output io.Writer) Interpreter {
	return Interpreter{
		input:  bufio.NewReader(input),
		output: bufio.NewWriter(output),
	}
}

func (i Interpreter) Eval(tape *Tape, node ast.Node) error {
	switch node.(type) {
	case ast.Command:
		if err := i.evalCommand(tape, node.(ast.Command)); err != nil {
			return fmt.Errorf("error evaluating command:%w", err)
		}
	case ast.Loop:
		if err := i.evalLoop(tape, node.(ast.Loop)); err != nil {
			return fmt.Errorf("error evaluating loop:%w", err)
		}
	}
	return nil
}

func (i Interpreter) evalLoop(tape *Tape, loop ast.Loop) error {
	for tape.Curr() != 0 {
		for _, stmt := range loop.Statements {
			if err := i.Eval(tape, stmt); err != nil {
				return err
			}
		}
	}
	return nil
}
func (i Interpreter) evalCommand(tape *Tape, cmd ast.Command) error {
	switch cmd.Token() {
	case lexer.Increment:
		tape.WriteCurr(tape.Curr() + 1)
	case lexer.Decrement:
		tape.WriteCurr(tape.Curr() - 1)
	case lexer.MoveRight:
		tape.GoRight()
	case lexer.MoveLeft:
		tape.GoLeft()
	case lexer.Read:
		b, err := i.input.ReadByte()
		if err != nil {
			return fmt.Errorf("error reading from input: %w\n", err)
		}
		tape.WriteCurr(b)
	case lexer.Write:
		log.Println("writing", string([]byte{tape.Curr()}))
		if err := i.output.WriteByte(tape.Curr()); err != nil {
			return fmt.Errorf("error writing to output: %w\n", err)
		}
	}
	return nil
}
