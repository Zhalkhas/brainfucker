package parser

import (
	"brainfucker/ast"
	"brainfucker/lexer"
	"reflect"
	"testing"
)

type parserTestInput struct {
	input         []lexer.Token
	output        *ast.Program
	expectedError error
}

type stubLexer struct {
	res       []lexer.Token
	stubError error
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
func (s stubLexer) Lex() ([]lexer.Token, error) {
	return s.res, s.stubError

}

func TestParser(t *testing.T) {
	t.Run("basic program parse", func(t *testing.T) {
		testParserWithInput(t, parserTestInput{
			input: []lexer.Token{lexer.Increment,
				lexer.Decrement,
				lexer.MoveRight,
				lexer.MoveLeft,
				lexer.Read,
				lexer.Write,
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
			input: []lexer.Token{
				lexer.JumpIfZero,
				lexer.Increment,
				lexer.Decrement,
				lexer.MoveRight,
				lexer.MoveLeft,
				lexer.Read,
				lexer.Write,
				lexer.JumpUnlessZero,
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
			input: []lexer.Token{
				lexer.JumpIfZero,
				lexer.JumpIfZero,
				lexer.Decrement,
				lexer.JumpUnlessZero,
				lexer.JumpUnlessZero,
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
