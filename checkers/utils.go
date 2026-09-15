// --------------------------------------------------------------
// wtf / Checkers / Utils -- Not meant to be used as a
// checker itself, but provides utility functions for checkers.
//  By: FormalBlaze
// --------------------------------------------------------------

package checkers

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

func divider() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80 // fallback if not a terminal (e.g. piped output)
	}
	return strings.Repeat("-", width)
}

var (
	logFileMu   sync.Mutex
	logFilePath string
)

// SetLogFile sets the path that LogMessage appends to. An empty path
// disables logging (the default).
func SetLogFile(path string) {
	logFileMu.Lock()
	defer logFileMu.Unlock()
	logFilePath = path
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "wtf Log started at %s\n", time.Now().Format(time.RFC3339))
}

// LogMessage appends a timestamped message to the log file configured via
// SetLogFile. It's a no-op if no log file has been set.
func LogMessage(message string) {
	logFileMu.Lock()
	path := logFilePath
	logFileMu.Unlock()

	if path == "" {
		return
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write log file: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "[%s] %s\n", time.Now().Format(time.RFC3339), message)
}