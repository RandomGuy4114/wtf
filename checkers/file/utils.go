// --------------------------------------------------------------
// wtf / Checkers / File / Utils -- Not meant to be used as a
// checker itself, but provides utility functions for checkers.
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
	"wtf/checkers"
)

type location struct {
	text string
	idx  int
}

// extractLocations finds every "file:line[:col]" location in a checker's
// output that points at path itself, each tagged with where it appears in
// output so it can be matched up against the error text it belongs to.
// It matches the location regardless of surrounding punctuation, since
// every language formats it differently: node wraps it in parens
// ("at foo (file.js:1:13)"), go prints it bare ("file.go:8:14: message"),
// and rustc/rust panics trail it with a colon ("panicked at file.rs:3:22:").
func extractLocations(output, path string) []location {
	base := filepath.Base(path)
	locRegexp := regexp.MustCompile(`([^\s()]*` + regexp.QuoteMeta(base) + `:\d+(?::\d+)?)`)

	var locs []location
	for _, m := range locRegexp.FindAllStringSubmatchIndex(output, -1) {
		locs = append(locs, location{text: output[m[2]:m[3]], idx: m[2]})
	}
	return locs
}

// nearestLocation returns the location text closest to pos in the output,
// so each matched error is paired with the location it actually came from.
func nearestLocation(locs []location, pos int) string {
	if len(locs) == 0 {
		return ""
	}
	best := locs[0]
	bestDist := abs(best.idx - pos)
	for _, l := range locs[1:] {
		if d := abs(l.idx - pos); d < bestDist {
			best, bestDist = l, d
		}
	}
	return best.text
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func divider() string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80 // fallback if not a terminal (e.g. piped output)
	}
	return strings.Repeat("-", width)
}

// diagnose loads modulePath's JSON rules and matches them against output,
// printing (and logging) a fix plus the nearest location for every match.
// langName is only used for the "no known fix found" message. It's the
// shared second half of every per-language checker (CheckJSErrors,
// CheckPythonErrors, etc.) — they differ only in how they produce output.
func diagnose(output, path, modulePath, langName string) {
	data, err := checkers.LoadJSONData(modulePath)
	if err != nil {
		fmt.Printf("Failed to load %s: %v\n", modulePath, err)
		return
	}

	lowerOutput := strings.ToLower(output)
	locs := extractLocations(output, path)
	fmt.Println(divider())
	found := false
	for _, e := range data.Errors {
		for _, s := range e.Strings {
			idx := strings.Index(lowerOutput, strings.ToLower(s))
			if idx == -1 {
				continue
			}

			loc := nearestLocation(locs, idx)
			fmt.Print(color.RedString("%s: %s\nFix: %s\n", data.Name, e.Message, e.Fix))
			if loc != "" {
				fmt.Println(color.YellowString("Location: %s", loc))
			}
			fmt.Println(divider())
			checkers.LogMessage(fmt.Sprintf("%s: %s | Fix: %s | Location: %s", data.Name, e.Message, e.Fix, loc))
			found = true
			break
		}
	}

	if !found {
		fmt.Println(color.RedString("No known fix found for this %s error.", langName))
		if loc := nearestLocation(locs, 0); loc != "" {
			fmt.Println(color.YellowString("Location: %s", loc))
		}
		checkers.LogMessage(fmt.Sprintf("no known fix found for %s error in %s", langName, path))
	}
}
