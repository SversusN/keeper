package main

import (
	"fmt"
	"runtime"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	fmt.Printf("%s %s %s\n", buildVersion, buildDate, runtime.Version())
}
