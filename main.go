// This is the Clear programming language
// This only exists for me to create an interpreted language following along with "Writing an Interpreter in Go" by Thorsten Ball
// The end goal for this is to implement all the code from the book, have it be as industry-standard as possible for me, and prepare me for my actual goal: writing a compiler
// This is the precursor to my project that will follow along with Thorsten Ball's sequel "Writing a Compiler in Go"

package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/ajtroup1/clearv2/repl"
)

func main() {
	// Retreives current user's name. Not necessary at all, but hey
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Clear programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	// Initiate the REPL to execute commands in Clear
	repl.Start(os.Stdin, os.Stdout)
}
