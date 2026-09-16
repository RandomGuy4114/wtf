// --------------------------------------------------------------
// wtf / Checkers / File / Main
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Check dispatches to the appropriate file-type checker based on the file's
// extension (.js -> CheckJSErrors, .py -> CheckPythonErrors, .go ->
// CheckGoErrors, .rs -> CheckRustErrors).
func Check(filePath string, output []string) {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".js":
		CheckJSErrors(filePath)
	case ".py":
		CheckPythonErrors(filePath)
	case ".go":
		CheckGoErrors(filePath)
	case ".rs":
		CheckRustErrors(filePath)
	default:
		fmt.Printf("No checker available for file type %q\n", filepath.Ext(filePath))
	}
}
