package compiler

import (
	"brainfucker/ast"
	"brainfucker/lexer"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Compiler struct {
	module  *ir.Module
	main    *ir.Block
	buf     *ir.Global
	currIdx *ir.Global
}

func NewCompiler() *Compiler {
	m := ir.NewModule()
	return &Compiler{
		module:  ir.NewModule(),
		buf:     m.NewGlobalDef("buf", constant.NewArray(types.NewArray(30000, types.I8))),
		currIdx: m.NewGlobal("currIdx", types.NewInt(16)),
	}
}

func (c Compiler) Compile(program ast.Program) (*ir.Module, error) {
	for _, node := range program.Nodes {
		switch node.(type) {
		case ast.Command:
			cmd := node.(ast.Command)
			if err := c.handleCommand(cmd); err != nil {
				return nil, err
			}
		case ast.Loop:
			loop := node.(ast.Loop)
			if err := c.handleLoop(loop); err != nil {
				return nil, err
			}
		}
	}
	return c.module, nil
}

func (c Compiler) handleCommand(command ast.Command) error {
	switch command.Token() {
	case lexer.Increment:

	}
	return nil
}

func (c Compiler) handleLoop(loop ast.Loop) error {
	return nil
}
