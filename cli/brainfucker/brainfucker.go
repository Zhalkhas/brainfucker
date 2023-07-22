package main

import (
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
	}
}

func handleRun(args ...string) {

}
