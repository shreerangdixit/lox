package main

import (
	"fmt"
	"github.com/shreerangdixit/lox/runner"
	"os"
)

var Version = "<NOT SET>"

func main() {
	if len(os.Args) > 1 {
		err := runner.RunFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	} else {
		runner.StartREPL(Version)
	}
}
