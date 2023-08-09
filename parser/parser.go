package parser

import (
	"brainfucker/ast"
	"brainfucker/lexer"
	"errors"
	"fmt"
)

var ErrInvalidLoop = errors.New("invalid loop parens")

type Parser struct {
	lexer lexer.Lexer
}

func NewParser(lexer lexer.Lexer) Parser {
	return Parser{lexer: lexer}
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := ast.Program{Nodes: []ast.Node{}}
	tokens, err := p.lexer.Lex()
	if err != nil {
		return &program, fmt.Errorf("error in lexing program: %w", err)
	}
	currTokenIndex := 0
	for currTokenIndex < len(tokens) {
		token := tokens[currTokenIndex]
		switch token {
		case lexer.JumpIfZero:
			var loop ast.Loop
			var err error
			loop, currTokenIndex, err = p.parseLoop(tokens, currTokenIndex+1)
			if err != nil {
				return nil, err
			}
			program.Nodes = append(program.Nodes, loop)
		default:
			program.Nodes = append(program.Nodes, ast.NewCommand(token))
			currTokenIndex++
		}
	}
	return &program, nil
}

func (p *Parser) parseLoop(tokens []lexer.Token, currTokenIdx int) (ast.Loop, int, error) {
	l := ast.Loop{}
	for currTokenIdx < len(tokens) {
		token := tokens[currTokenIdx]
		switch token {
		case lexer.JumpIfZero:
			var loop ast.Loop
			var err error
			loop, currTokenIdx, err = p.parseLoop(tokens, currTokenIdx+1)
			if err != nil {
				return l, currTokenIdx, fmt.Errorf("error parsing nested loop: %w", err)

			}
			l.Statements = append(l.Statements, loop)
		case lexer.JumpUnlessZero:
			return l, currTokenIdx + 1, nil
		default:
			l.Statements = append(l.Statements, ast.NewCommand(token))
			currTokenIdx++
		}
	}
	return ast.Loop{}, currTokenIdx, ErrInvalidLoop
}
