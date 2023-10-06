package main

import (
	"fmt"
	"os"

	"github.com/Arch-4ng3l/Monkey/exec"
	"github.com/Arch-4ng3l/Monkey/repl"
)

func main() {
	input := os.Stdin

	if len(os.Args) != 1 {
		fileName := os.Args[1]
		content, _ := os.ReadFile(fileName)
		output := exec.ExecCodeWithInterpreter(string(content))
		fmt.Println("output " + output)
	} else {
		repl.StartComp(input, os.Stdout)
	}

}
