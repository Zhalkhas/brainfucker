package ast

import "brainfucker/lexer"

type Program struct {
	Nodes []Node
}

type Node interface {
	Token() lexer.Token
}

type Command struct {
	token lexer.Token
}

func NewCommand(token lexer.Token) Command {
	return Command{token: token}
}

func (c Command) Token() lexer.Token {
	return c.token
}

type Loop struct {
	Statements []Node
}

func (l Loop) Token() lexer.Token {
	return lexer.JumpIfZero
}
