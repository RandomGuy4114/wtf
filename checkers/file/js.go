// --------------------------------------------------------------
// wtf / Checkers / File / JS
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func CheckJSErrors(path string) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve path: %v\n", err)
		return nil, err
	}

	checkCmd := exec.Command("node", "--check", absPath)
	checkOutput, err := checkCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(checkOutput))
		diagnose(string(checkOutput), absPath, "modules/json/js.json", "JavaScript")
		return checkOutput, err
	}

	cmd := exec.Command("node", absPath)

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		diagnose(string(output), absPath, "modules/json/js.json", "JavaScript")
		return output, err
	}

	return output, nil
}
