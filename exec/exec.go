package exec

import (
	"bytes"
	"io"
	"os"

	"github.com/Arch-4ng3l/Monkey/eval"
	"github.com/Arch-4ng3l/Monkey/lexer"
	"github.com/Arch-4ng3l/Monkey/object"
	"github.com/Arch-4ng3l/Monkey/parser"
)

func ExecFile(fileName string) {

	b, err := os.ReadFile(fileName)

	if err != nil {
		return
	}

	l := lexer.NewLexer(string(b))
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnv()

	eval.Eval(program, env)
}

func ExecCode(code string) string {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w
	l := lexer.NewLexer(code)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	env := object.NewEnv()
	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	eval.Eval(program, env)

	w.Close()

	os.Stdout = oldStdout
	output := <-outC

	return output
}
