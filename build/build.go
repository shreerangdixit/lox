package build

import (
	"fmt"
)

var version = "<NOT SET>"
var date = "<NOT SET>"
var os = "<NOT SET>"
var host = "<NOT SET>"
var arch = "<NOT SET>"
var kernelVersion = "<NOT SET>"

var Info *BuildInfo

func init() {
	Info = &BuildInfo{
		Version:            version,
		BuildDate:          date,
		BuildOS:            os,
		BuildHost:          host,
		BuildArch:          arch,
		BuildKernelVersion: kernelVersion,
	}
}

type BuildInfo struct {
	Version            string
	BuildDate          string
	BuildOS            string
	BuildHost          string
	BuildArch          string
	BuildKernelVersion string
}

func (b *BuildInfo) String() string {
	info := ""
	info += fmt.Sprintf("Version: %s\n", b.Version)
	info += "Build Info:\n"
	info += fmt.Sprintf("  Date: %s\n", b.BuildDate)
	info += fmt.Sprintf("  OS: %s\n", b.BuildOS)
	info += fmt.Sprintf("  Host: %s\n", b.BuildHost)
	info += fmt.Sprintf("  Arch: %s\n", b.BuildArch)
	info += fmt.Sprintf("  Kernel Version: %s\n", b.BuildKernelVersion)
	return info
}
