package lexer

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

func collectChan(ch <-chan TokenResult) []TokenResult {
	res := make([]TokenResult, 0)
	for tr := range ch {
		res = append(res, tr)
	}
	return res
}

type lexerTestInput struct {
	input  string
	output []TokenResult
}

func testLexerWithInput(t *testing.T, testData lexerTestInput) {
	l := NewLexer(bytes.NewBuffer([]byte(testData.input)))
	actual := collectChan(l.Lex())
	if !reflect.DeepEqual(testData.output, actual) {
		t.Errorf("expected %v, got %v", testData.output, actual)
	}
}

func TestLexer(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		testLexerWithInput(t, lexerTestInput{
			input:  "",
			output: []TokenResult{},
		})
	})

	t.Run("valid input", func(t *testing.T) {
		testLexerWithInput(t, lexerTestInput{
			input:  "+-<>.,[]",
			output: []TokenResult{NewTokenResultValue(Increment), NewTokenResultValue(Decrement), NewTokenResultValue(MoveLeft), NewTokenResultValue(MoveRight), NewTokenResultValue(Write), NewTokenResultValue(Read), NewTokenResultValue(JumpIfZero), NewTokenResultValue(JumpUnlessZero)},
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
			output: []TokenResult{NewTokenResultValue(Increment), NewTokenResultValue(Decrement), NewTokenResultValue(MoveLeft), NewTokenResultValue(MoveRight), NewTokenResultValue(Write), NewTokenResultValue(Read), NewTokenResultValue(JumpIfZero), NewTokenResultValue(JumpUnlessZero)},
		})
	})

	t.Run("throwing io.EOF in lexer", func(t *testing.T) {
		l := NewLexer(stubRuneReader{runeReaderFunc: func() (r rune, size int, err error) {
			return 0, 0, io.EOF
		}})
		expected := make([]TokenResult, 0)
		actual := collectChan(l.Lex())
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("throwing error in lexer", func(t *testing.T) {
		err := fmt.Errorf("test error")
		l := NewLexer(stubRuneReader{runeReaderFunc: func() (rune, int, error) {
			return 0, 0, err
		}})
		expected := []TokenResult{NewTokenResultError(err)}
		actual := collectChan(l.Lex())
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
