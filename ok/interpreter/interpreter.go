package interpreter

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jesseduffield/OK/evaluator"
	"github.com/jesseduffield/OK/lexer"
	"github.com/jesseduffield/OK/object"
	"github.com/jesseduffield/OK/parser"
)

func Interpret(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	out := os.Stdout

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
