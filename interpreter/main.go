package interpreter

import (
	"fmt"
	"learning/interpreter/repl"
	"os"
	"os/user"
)

func Run() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!", user.Username)
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}