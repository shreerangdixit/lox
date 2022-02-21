package main

import (
	"flag"
	"fmt"
	"github.com/shreerangdixit/lox/runner"
	"os"
)

var Version = "<NOT SET>"
var BuildDate = "<NOT SET>"
var BuildOS = "<NOT SET>"
var BuildHost = "<NOT SET>"
var BuildArch = "<NOT SET>"
var BuildKernelVersion = "<NOT SET>"

var flagVer bool

func init() {
	flag.BoolVar(&flagVer, "v", false, "Display version/build info")
}

func main() {
	flag.Parse()

	if flagVer {
		fmt.Fprint(os.Stdout, buildInfo())
		os.Exit(0)
	} else if len(os.Args) > 1 {
		err := runner.RunFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	} else {
		runner.StartREPL(buildInfo())
	}
}

func buildInfo() string {
	info := ""
	info += fmt.Sprintf("Version: %s\n", Version)
	info += "Build Info:\n"
	info += fmt.Sprintf("  Date: %s\n", BuildDate)
	info += fmt.Sprintf("  OS: %s\n", BuildOS)
	info += fmt.Sprintf("  Host: %s\n", BuildHost)
	info += fmt.Sprintf("  Arch: %s\n", BuildArch)
	info += fmt.Sprintf("  Kernel Version: %s\n", BuildKernelVersion)
	return info
}
