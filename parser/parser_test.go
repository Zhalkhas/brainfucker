package parser

import (
	"brainfucker/ast"
	"brainfucker/lexer"
	"reflect"
	"testing"
)

type parserTestInput struct {
	input         []lexer.TokenResult
	output        *ast.Program
	expectedError error
}

type stubLexer struct {
	res []lexer.TokenResult
}

func testParserWithInput(t *testing.T, testData parserTestInput) {
	p := NewParser(stubLexer{res: testData.input})
	actualOutput, actualErr := p.Parse()
	if actualErr != testData.expectedError {
		t.Errorf("unexpected error, expected: %+v, actual: %+v", testData.expectedError, actualErr)
		return
	}
	if !reflect.DeepEqual(actualOutput, testData.output) {
		t.Errorf("expected %v, got %v", testData.output, actualOutput)

	}
}
func (s stubLexer) Lex() <-chan lexer.TokenResult {
	res := make(chan lexer.TokenResult)
	go func() {
		defer close(res)
		for i := 0; i < len(s.res); i++ {
			res <- s.res[i]
		}
	}()
	return res

}

func TestParser(t *testing.T) {
	t.Run("basic program parse", func(t *testing.T) {
		testParserWithInput(t, parserTestInput{
			input: []lexer.TokenResult{
				lexer.NewTokenResultValue(lexer.Increment),
				lexer.NewTokenResultValue(lexer.Decrement),
				lexer.NewTokenResultValue(lexer.MoveRight),
				lexer.NewTokenResultValue(lexer.MoveLeft),
				lexer.NewTokenResultValue(lexer.Read),
				lexer.NewTokenResultValue(lexer.Write),
			},
			output: &ast.Program{Nodes: []ast.Node{
				ast.NewCommand(lexer.Increment),
				ast.NewCommand(lexer.Decrement),
				ast.NewCommand(lexer.MoveRight),
				ast.NewCommand(lexer.MoveLeft),
				ast.NewCommand(lexer.Read),
				ast.NewCommand(lexer.Write),
			}},
			expectedError: nil,
		})
	})

	t.Run("parsing loops", func(t *testing.T) {
		testParserWithInput(t, parserTestInput{
			input: []lexer.TokenResult{
				lexer.NewTokenResultValue(lexer.JumpIfZero),
				lexer.NewTokenResultValue(lexer.Increment),
				lexer.NewTokenResultValue(lexer.Decrement),
				lexer.NewTokenResultValue(lexer.MoveRight),
				lexer.NewTokenResultValue(lexer.MoveLeft),
				lexer.NewTokenResultValue(lexer.Read),
				lexer.NewTokenResultValue(lexer.Write),
				lexer.NewTokenResultValue(lexer.JumpUnlessZero),
			},
			output: &ast.Program{Nodes: []ast.Node{ast.Loop{Statements: []ast.Node{
				ast.NewCommand(lexer.Increment),
				ast.NewCommand(lexer.Decrement),
				ast.NewCommand(lexer.MoveRight),
				ast.NewCommand(lexer.MoveLeft),
				ast.NewCommand(lexer.Read),
				ast.NewCommand(lexer.Write)}},
			}},
			expectedError: nil,
		})
	})

	t.Run("parsing inner loops", func(t *testing.T) {
		testParserWithInput(t, parserTestInput{
			input: []lexer.TokenResult{
				lexer.NewTokenResultValue(lexer.JumpIfZero),
				lexer.NewTokenResultValue(lexer.JumpIfZero),
				lexer.NewTokenResultValue(lexer.Decrement),
				lexer.NewTokenResultValue(lexer.JumpUnlessZero),
				lexer.NewTokenResultValue(lexer.JumpUnlessZero),
			},
			output: &ast.Program{Nodes: []ast.Node{
				ast.Loop{Statements: []ast.Node{
					ast.Loop{Statements: []ast.Node{ast.NewCommand(lexer.Decrement)}},
				}},
			}},
			expectedError: nil,
		})
	})
}
