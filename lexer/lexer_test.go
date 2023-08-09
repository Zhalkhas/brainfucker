package lexer

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

type lexerTestInput struct {
	input         string
	output        []Token
	exceptedError error
}

func testLexerWithInput(t *testing.T, testData lexerTestInput) {
	l := NewLexer(bytes.NewBuffer([]byte(testData.input)))
	actual, actualErr := l.Lex()
	if actualErr != testData.exceptedError {
		t.Errorf("unexpected error: %v", actualErr)
		return
	}
	if !reflect.DeepEqual(testData.output, actual) {
		t.Errorf("expected %v, got %v", testData.output, actual)
	}
}

func TestLexer(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		testLexerWithInput(t, lexerTestInput{
			input:  "",
			output: []Token{},
		})
	})

	t.Run("valid input", func(t *testing.T) {
		testLexerWithInput(t, lexerTestInput{
			input:  "+-<>.,[]",
			output: []Token{Increment, Decrement, MoveLeft, MoveRight, Write, Read, JumpIfZero, JumpUnlessZero},
		})
	})

	t.Run("ignore unknown tokens", func(t *testing.T) {
		input := []byte("+-<>.,[]")
		// inserting random runes to random positions of input
		for i := 0; i < 100; i++ {
			b := byte(rand.Int() % 255)
			if _, ok := tokens[Token(b)]; ok {
				continue
			}
			input = append(input, 0)
			copy(input[i+1:], input[i:])
			input[i] = b
		}
		// passing new generated input to lexer
		testLexerWithInput(t, lexerTestInput{
			input:  string(input),
			output: []Token{Increment, Decrement, MoveLeft, MoveRight, Write, Read, JumpIfZero, JumpUnlessZero},
		})
	})

	t.Run("throwing io.EOF in lexer", func(t *testing.T) {
		l := NewLexer(stubRuneReader{runeReaderFunc: func() (r rune, size int, err error) {
			return 0, 0, io.EOF
		}})
		expected := make([]Token, 0)
		actual, actualError := l.Lex()
		if actualError != nil {
			t.Errorf("unexpected error: %v", actualError)
			return
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("throwing error in lexer", func(t *testing.T) {
		err := fmt.Errorf("test error")
		l := NewLexer(stubRuneReader{runeReaderFunc: func() (rune, int, error) {
			return 0, 0, err
		}})
		expected := make([]Token, 0)
		actual, actualErr := l.Lex()
		if actualErr != err {
			t.Errorf("unexpected error: %v", actualErr)
			return
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

type stubRuneReader struct {
	runeReaderFunc func() (r rune, size int, err error)
}

func (s stubRuneReader) ReadRune() (r rune, size int, err error) {
	return s.runeReaderFunc()
}
