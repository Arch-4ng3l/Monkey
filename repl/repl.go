package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Arch-4ng3l/Monkey/eval"
	"github.com/Arch-4ng3l/Monkey/lexer"
	"github.com/Arch-4ng3l/Monkey/object"
	"github.com/Arch-4ng3l/Monkey/parser"
	"github.com/TwiN/go-color"
)

const PROMPT = " ➡ "
const MONKEY_FACE = `            
+----------------------------------------------------+
|     Welcome To The Monkey Programming Language     |
+----------------------------------------------------+
|                                                    |
|   +--------------------------------------------+   |
|   |                    __,__                   |   |
|   |           .--.  .-"     "-.  .--.          |   |
|   |          / .. \/  .-. .-.  \/ .. \         |   |
|   |         | |  '|  /   Y   \  |'  | |        |   |
|   |         | \   \  \ 0 | 0 /  /   / |        |   |
|   |          \ '- ,\.-"""""""-./, -' /         |   |
|   |           ''-' /_   ^ ^   _\ '-''          |   |
|   |               |  \._   _./  |              |   |
|   |               \   \ '~' /   /              |   |
|   |                '._ '-=-' _.'               |   |
|   |                   '-----'                  |   |
|   +--------------------------------------------+   |
|                                                    |
+----------------------------------------------------+
`

func Sart(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Fprintf(out, "%s%s%s%s", color.Green, color.Bold, MONKEY_FACE, color.Reset)
	env := object.NewEnv()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		if errs := p.Errors(); len(errs) != 0 {
			for _, err := range errs {
				fmt.Fprintf(out, "ERROR: %s\n", err)
			}
		}
		eval.Eval(program, env)
	}
}