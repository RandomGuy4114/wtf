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

// Check dispatches to the appropriate file-type checker based on the
// file's extension (e.g. .js -> CheckJSErrors, .py -> CheckPythonErrors).
func Check(filePath string, output []string) {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".js":
		CheckJSErrors(filePath)
	case ".py":
		CheckPythonErrors(filePath)
	default:
		fmt.Printf("No checker available for file type %q\n", filepath.Ext(filePath))
	}
}
