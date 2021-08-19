package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/jesseduffield/OK/ok/interpreter"
	"github.com/jesseduffield/OK/ok/repl"
)

func main() {
	if len(os.Args) == 1 {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Hello %s! This is the OK? programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	} else {
		filename := os.Args[1]

		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

		interpreter.Interpret(f, os.Stdout)
	}
}
