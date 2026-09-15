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
	logo = `
           __  ____
 _      __/ /_/ __/
| | /| / / __/ /_  
| |/ |/ / /_/ __/  
|__/|__/\__/_/     

`
)

const (
	version = "1.0.0"
)

func printVersion() {
	fmt.Printf(logo)
	fmt.Printf("wtf version: %s\n", version)
	fmt.Printf("Go version: %s\n", runtime.Version())
}