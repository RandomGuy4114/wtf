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

	"golang.org/x/term"
)

type location struct {
	text string
	idx  int
}

// extractLocations finds every "file:line[:col]" location in a checker's
// output that points at path itself, each tagged with where it appears in
// output so it can be matched up against the error text it belongs to.
func extractLocations(output, path string) []location {
	base := filepath.Base(path)
	stackRegexp := regexp.MustCompile(`\(([^()\s]*` + regexp.QuoteMeta(base) + `:\d+:\d+)\)`)
	fileLineRegexp := regexp.MustCompile(`(?m)^(.+` + regexp.QuoteMeta(filepath.Ext(path)) + `):(\d+)$`)

	var locs []location
	for _, m := range stackRegexp.FindAllStringSubmatchIndex(output, -1) {
		locs = append(locs, location{text: output[m[2]:m[3]], idx: m[2]})
	}
	if len(locs) == 0 {
		if m := fileLineRegexp.FindStringSubmatchIndex(output); m != nil {
			locs = append(locs, location{
				text: fmt.Sprintf("%s:%s", output[m[2]:m[3]], output[m[4]:m[5]]),
				idx:  m[0],
			})
		}
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
