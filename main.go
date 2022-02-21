package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shreerangdixit/lox/build"
	"github.com/shreerangdixit/lox/runner"
)

var flagVer bool

func init() {
	flag.BoolVar(&flagVer, "v", false, "Display version/build info")
}

func main() {
	flag.Parse()

	if flagVer {
		fmt.Fprint(os.Stdout, build.Info)
		os.Exit(0)
	} else if len(os.Args) > 1 {
		err := runner.RunFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	} else {
		runner.StartREPL()
	}
}
