package interpreter

import (
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jesseduffield/OK/ok/evaluator"
	"github.com/jesseduffield/OK/ok/lexer"
	"github.com/jesseduffield/OK/ok/object"
	"github.com/jesseduffield/OK/ok/parser"
	"github.com/jesseduffield/OK/ok/quentyn"
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
	output := evaluator.New(w).Eval(program, env)
	if v, ok := output.(*object.Error); ok {
		io.WriteString(w, v.Inspect())
		io.WriteString(w, "\n")
	}

	quentynMessage := quentyn.GetQuentynMessage()
	if quentynMessage != "" {
		io.WriteString(w, "\n")
		io.WriteString(w, quentynMessage)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " Parser errors:\n")
	for _, msg := range errors {
		indentedMsg := strings.Replace(msg, "\n", "\n\t", -1)
		io.WriteString(out, "\t"+indentedMsg+"\n")
	}
}
