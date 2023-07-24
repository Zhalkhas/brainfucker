package main

import (
	"brainfucker/interpreter"
	"brainfucker/lexer"
	"brainfucker/parser"
	"bufio"
	"log"
	"os"
)

func main() {
	var args []string
	if len(os.Args) == 1 {
		args = []string{"repl"}
	} else {
		args = os.Args[1:]
	}
	log.Println("well cum to brainfuck")
	cmd := args[0]
	log.Println("args", args)
	if len(args) > 1 {
		args = args[1:]
	}

	switch cmd {
	case "run":
		handleRun(args...)
	}
}

func handleRun(args ...string) {
	name := args[0]
	prog, err := os.Open(name)
	if err != nil {
		log.Panicln("error opening source file:", err)
	}
	log.Println("running", name)
	p := parser.NewParser(lexer.NewLexer(bufio.NewReader(prog)))
	ast, err := p.Parse()
	if err != nil {
		log.Panicln("error parsing source file:", err)
	}
	log.Println("parsed", name)
	i := interpreter.NewInterpreter(os.Stdin, os.Stdout)
	tape := interpreter.NewTape(30000)
	log.Println("starting interpreting", name)
	for _, node := range ast.Nodes {
		if err := i.Eval(&tape, node); err != nil {
			log.Panicln("error executing code:", err)
		}
	}
	log.Println("interpreting finished")
}
