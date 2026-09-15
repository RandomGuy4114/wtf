// --------------------------------------------------------------
// wtf / Version -- Version information for the wtf tool
//  By: FormalBlaze
// --------------------------------------------------------------

package main

import (
	"fmt"
	"runtime"
)

const (
	version = "1.0.0"
)

func printVersion() {
	fmt.Printf("wtf version: %s\n", version)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}