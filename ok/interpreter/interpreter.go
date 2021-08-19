package interpreter

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/jesseduffield/OK/ok/evaluator"
	"github.com/jesseduffield/OK/ok/lexer"
	"github.com/jesseduffield/OK/ok/object"
	"github.com/jesseduffield/OK/ok/parser"
)

func Interpret(r io.Reader, w io.Writer) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(w, p.Errors())
		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.New(w).Eval(program, env)
	if evaluated != nil {
		io.WriteString(w, evaluated.Inspect())
		io.WriteString(w, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
