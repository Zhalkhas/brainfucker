package brainfucker

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestRunInterpreter(t *testing.T) {
	t.Run("basic brainfuck program", func(t *testing.T) {
		progInput := ""
		prog := "++>++>+++<-"
		i := NewInterpreter(strings.NewReader(progInput), os.Stdout)
		err := i.Run(strings.NewReader(prog))
		if err != nil {
			t.Error(err)
			return
		}
		assertCells([]byte{2, 1, 3}, i.dataCells, t)
	})

	t.Run("test program output", func(t *testing.T) {
		progInput := ""
		prog := "+++++.>++.>+++."
		expected := []byte{5}
		bytesWriter := bytes.NewBuffer([]byte{})
		i := NewInterpreter(strings.NewReader(progInput), bytesWriter)
		err := i.Run(strings.NewReader(prog))
		if err != nil {
			t.Errorf("%+v\n", err)
			return
		}

		bytesBuf := make([]byte, len(expected))
		_, err = bytesWriter.Read(bytesBuf)
		if err != nil {
			t.Errorf("%+v\n", err)
			return
		}

		assertCells(expected, i.dataCells, t)
		assertCells(expected, bytesBuf, t)
	})

	t.Run("multiple values output test", func(t *testing.T) {
		progInput := ""
		prog := "+++++.>++.>+++."
		expected := []byte{5, 2, 3}
		bytesWriter := bytes.NewBuffer([]byte{})
		i := NewInterpreter(strings.NewReader(progInput), bytesWriter)
		err := i.Run(strings.NewReader(prog))
		if err != nil {
			t.Errorf("%+v\n", err)
			return
		}

		bytesBuf := make([]byte, len(expected))
		_, err = bytesWriter.Read(bytesBuf)
		if err != nil {
			t.Errorf("%+v\n", err)
			return
		}

		assertCells(expected, i.dataCells, t)
		assertCells(expected, bytesBuf, t)
	})

	t.Run("test program input", func(t *testing.T) {
		progInput := ""
		prog := ",++"
		expected := "5"
		bytesWriter := bytes.NewBufferString("")
		i := NewInterpreter(strings.NewReader(progInput), bytesWriter)
		err := i.Run(strings.NewReader(prog))
		if err != nil {
			t.Errorf("%+v\n", err)
			return
		}
		assertCells([]byte{5}, i.dataCells, t)
		bytesBuf := make([]byte, 1)
		_, err = bytesWriter.Read(bytesBuf)
		if err != nil {
			t.Error(err)
		}
		actual := string(bytesBuf)
		if actual != fmt.Sprint(expected) {
			t.Errorf("error: values does not match, expected %+v, actual %+v", expected, actual)
		}
	})

	t.Run("hello world brainfuck program", func(t *testing.T) {
		progInput := ""
		prog := ">++++++++[<+++++++++>-]<.>++++[<+++++++>-]<+.+++++++..+++.>>++++++[<+++++++>-]<++.------------.>++++++[<+++++++++>-]<+.<.+++.------.--------.>>>++++[<++++++++>-]<+."
		buffWriter := bytes.NewBufferString("")
		i := NewInterpreter(strings.NewReader(progInput), buffWriter)
		err := i.Run(strings.NewReader(prog))
		if err != nil {
			t.Error(err)
			return
		}
		expected := "Hello, World!"
		actual := buffWriter.String()
		if actual != expected {
			t.Errorf("error: values does not match, expected %+v, actual %+v", expected, actual)
		}
	})
}

func assertCells(expected, actual []byte, t *testing.T) {
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("error: values does not match, expected %+v, actual %+v", expected, actual[:len(expected)])
			return
		}
	}
}
