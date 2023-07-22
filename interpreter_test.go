package brainfucker

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

type testInput struct {
	progInput          string
	program            string
	progOutputExpected string
	dataCellsExpected  []byte
}

func testWithInput(t *testing.T, testData testInput) {
	out := bytes.NewBuffer([]byte{})
	i := NewInterpreter(bufio.NewReader(strings.NewReader(testData.program)), strings.NewReader(testData.progInput), out)
	if err := i.Run(); err != nil {
		t.Errorf("error running interpreter: %+v\n", err)
		return
	}

	if testData.dataCellsExpected != nil {
		assertCells(testData.dataCellsExpected, i.buffer(), t)
	}

	actual := string(out.Bytes())
	if actual != testData.progOutputExpected {
		t.Errorf("error: program output values does not match, expected %+v, actual %+v", testData.progOutputExpected, actual)
	}
}

func assertCells(expected, actual []byte, t *testing.T) {
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("error: interpreter data state does not match, expected %+v, actual %+v", expected, actual[:len(expected)])
			return
		}
	}
}

func TestRunInterpreter(t *testing.T) {
	t.Run("basic brainfuck program", func(t *testing.T) {
		testWithInput(t, testInput{
			progInput:          "",
			program:            "++>++>+++<-",
			progOutputExpected: "",
			dataCellsExpected:  []byte{2, 1, 3},
		})
	})

	t.Run("test program output", func(t *testing.T) {
		testWithInput(t, testInput{
			progInput:          "",
			program:            "+++++.>++.>+++.",
			progOutputExpected: string([]byte{5, 2, 3}),
			dataCellsExpected:  []byte{5, 2, 3},
		})
	})

	t.Run("multiple values output test", func(t *testing.T) {
		testWithInput(t, testInput{
			progInput:          "",
			program:            "+++++.>++.>+++.",
			progOutputExpected: string([]byte{5, 2, 3}),
			dataCellsExpected:  []byte{5, 2, 3},
		})
	})

	t.Run("test program input", func(t *testing.T) {
		testWithInput(t, testInput{
			progInput:          "3",
			program:            ",++",
			progOutputExpected: "",
			dataCellsExpected:  []byte{"5"[0]},
		})
	})

	t.Run("hello world brainfuck program", func(t *testing.T) {
		testWithInput(t, testInput{
			progInput:          "",
			program:            ">++++++++[<+++++++++>-]<.>++++[<+++++++>-]<+.+++++++..+++.>>++++++[<+++++++>-]<++.------------.>++++++[<+++++++++>-]<+.<.+++.------.--------.>>>++++[<++++++++>-]<+.",
			progOutputExpected: "Hello, World!",
			dataCellsExpected:  nil,
		})
	})
}
