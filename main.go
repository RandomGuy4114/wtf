// --------------------------------------------------------------
// wtf / Main
//  By: FormalBlaze
// --------------------------------------------------------------

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"wtf/checkers"
	"wtf/checkers/file"
)

func main() {
	defaultCommand := os.Getenv("COMMAND")
	defaultFile := os.Getenv("FILE")
	logFile := os.Getenv("LOG_FILE")
	commandPtr := flag.String("r", defaultCommand, "Command to run")
	filePtr := flag.String("f", defaultFile, "File to check for errors (required)")
	versionPtr := flag.Bool("v", false, "Print version information")
	logFilePtr := flag.String("l", logFile, "Log file path (optional)")
	flag.Parse()

	checkers.SetLogFile(*logFilePtr)

	if *versionPtr {
		printVersion()
		return
	}

	if *filePtr != "" {
		info, err := os.Stat(*filePtr)
		if err != nil || info.IsDir() {
			log.Fatalf("-f must point to an existing file: %s", *filePtr)
		}
	}

	if *commandPtr == "" {
		if *filePtr == "" {
			log.Fatalf("either -r (command to run) or -f (file to check) must be provided")
		}
		content, err := os.ReadFile(*filePtr)
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}
		file.Check(*filePtr, strings.Split(string(content), "\n"))
		return
	}

	parts := strings.Fields(*commandPtr)
	cmd := exec.Command(parts[0], parts[1:]...)

	out, err := cmd.CombinedOutput()
	checkers.LogMessage(fmt.Sprintf("ran %q, output: %s", *commandPtr, string(out)))
	if err != nil {
		log.Printf("Error found, diagnosing.")
		println(string(out))
		if *filePtr != "" {
			file.Check(*filePtr, strings.Split(string(out), "\n"))
		}
		if !checkers.JudgeError(string(out) + " " + err.Error()) {
			log.Printf("No known fix found.")
			checkers.LogMessage("no known fix found")
		}
	} else {
		log.Printf("Command executed successfully.")
		checkers.LogMessage(fmt.Sprintf("ran %q successfully", *commandPtr))
	}
}