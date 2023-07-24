package parser

import (
	"brainfucker/ast"
	"brainfucker/lexer"
	"fmt"
	"log"
)

var ErrInvalidLoop = fmt.Errorf("invalid loop parens")

type Parser struct {
	lexer lexer.Lexer
}

func NewParser(lexer lexer.Lexer) Parser {
	return Parser{lexer: lexer}
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := ast.Program{Nodes: []ast.Node{}}
	log.Println("starting parser")
	for {
		result, ok := <-p.lexer.Lex()
		log.Println("receiving token", result)
		if !ok {
			log.Println("cannot read from chan")
		}
		if err := result.Error(); err != nil {
			return nil, err
		}

		token := result.Token()
		if token == lexer.NilToken {
			return nil, fmt.Errorf("invalid token result from lexer: both token and error are nil")
		}

		log.Println("token in parser", token)
		switch token {
		case lexer.JumpIfZero:
			for i := 0; i < 10; i++ {
				fmt.Println("prikol", <-p.lexer.Lex())
			}
			fmt.Println("prikol", <-p.lexer.Lex())
			loop, err := p.parseLoop()
			if err != nil {
				return nil, err
			}
			program.Nodes = append(program.Nodes, loop)
		default:
			program.Nodes = append(program.Nodes, ast.NewCommand(token))
		}
	}
	return &program, nil
}

var parseLoopCalls = 0

func (p *Parser) parseLoop() (ast.Loop, error) {
	parseLoopCalls++
	fmt.Println("parseLoopCalls", parseLoopCalls)
	l := ast.Loop{}
	for result := range p.lexer.Lex() {
		if err := result.Error(); err != nil {
			return ast.Loop{}, err
		}

		token := result.Token()
		if token == lexer.NilToken {
			return ast.Loop{}, fmt.Errorf("invalid token result from lexer: both token and error are nil")
		}

		log.Println("token in parseLoop", token)
		switch token {
		case lexer.JumpIfZero:
			loop, err := p.parseLoop()
			if err != nil {
				return ast.Loop{}, fmt.Errorf("error parsing nested loop: %w", err)

			}
			l.Statements = append(l.Statements, loop)
		case lexer.JumpUnlessZero:
			if len(l.Statements) == 0 {
				return ast.Loop{}, nil
			}
			return l, nil
		default:
			l.Statements = append(l.Statements, ast.NewCommand(token))
		}
	}
	return ast.Loop{}, ErrInvalidLoop
}

type tokenStack []ast.Node

func newTokenStack() tokenStack {
	return tokenStack{}
}

func (s *tokenStack) Push(val ast.Node) {
	*s = append(*s, val)
}
func (s *tokenStack) Pop() ast.Node {
	l := len(*s)
	if l > 0 {
		last := (*s)[l-1]
		*s = (*s)[:l-1]
		return last
	}
	return nil
}

func (s *tokenStack) Peek() ast.Node {
	l := len(*s)
	if l > 0 {
		return (*s)[l-1]
	}
	return nil
}
