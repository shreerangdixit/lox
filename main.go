package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shreerangdixit/yeti/build"
	"github.com/shreerangdixit/yeti/run"
)

var flagVer bool

func init() {
	flag.BoolVar(&flagVer, "v", false, "Display version/build info")
}

func main() {
	flag.Parse()

	if flagVer {
		fmt.Println(build.Info)
		os.Exit(0)
	} else if len(os.Args) > 1 {
		err := run.File(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	} else {
		run.REPL()
	}
}
