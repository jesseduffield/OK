package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/jesseduffield/OK/interpreter"
	"github.com/jesseduffield/OK/repl"
)

func main() {
	if len(os.Args) == 1 {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the OK? programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	} else {
		filename := os.Args[1]
		interpreter.Interpret(filename)
	}
}
