package main

import (
	"brainfucker"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

func main() {
	var args []string
	if len(os.Args) == 1 {
		args = []string{"repl"}
	} else {
		args = os.Args[1:]
	}

	cmd := args[0]

	if len(args) > 1 {
		args = args[1:]
	}

	switch cmd {
	case "run":
		handleRun(args...)
		break
	case "compile":
		break
	default:
		handleRepl(args...)
	}
}

func handleRepl(args ...string) {
	flagSet := flag.NewFlagSet("repl", flag.ExitOnError)
	var cellsBuffSize int
	flagSet.IntVar(&cellsBuffSize, "buff-size", 30000, "sets cells buffer size")
	err := flagSet.Parse(args)
	if err != nil {
		panic(err)
		return
	}
	interpreter := brainfucker.NewInterpreter(os.Stdin, os.Stdout, cellsBuffSize)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Bytes()
		fmt.Println("input", string(input))
		err := interpreter.Run(bytes.NewReader(input))
		if err != nil {
			panic(err)
			return
		}
	}
}

func handleRun(args ...string) {

}
