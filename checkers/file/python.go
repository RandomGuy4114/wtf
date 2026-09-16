// --------------------------------------------------------------
// wtf / Checkers / File / Python
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func CheckPythonErrors(path string) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve path: %v\n", err)
		return nil, err
	}

	cmd := exec.Command("python3", absPath)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		diagnose(string(output), absPath, "modules/json/python.json", "Python")
		return output, err
	}

	return output, nil
}
