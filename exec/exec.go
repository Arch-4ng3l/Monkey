package exec

import (
	"fmt"

	"github.com/Arch-4ng3l/Monkey/eval"
	"github.com/Arch-4ng3l/Monkey/lexer"
	"github.com/Arch-4ng3l/Monkey/object"
	"github.com/Arch-4ng3l/Monkey/parser"
)

func ExecCode(code string) string {
	//oldStdout := os.Stdout

	//r, w, _ := os.Pipe()

	fmt.Println("start")
	//os.Stdout = w
	l := lexer.NewLexer(code)
	p := parser.NewParser(l)
	fmt.Println("start Parsing")
	program := p.ParseProgram()

	env := object.NewEnv()
	//outC := make(chan string)

	//go func() {
	//var buf bytes.Buffer
	//io.Copy(&buf, r)
	//outC <- buf.String()
	//}()

	fmt.Println("start Eval")
	eval.Eval(program, env)

	fmt.Println("Done eval")
	//w.Close()

	//os.Stdout = oldStdout
	//output := <-outC
	output := ""

	return output
}
