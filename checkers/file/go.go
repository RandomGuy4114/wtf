// --------------------------------------------------------------
// wtf / Checkers / File / Go
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func CheckGoErrors(path string) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve path: %v\n", err)
		return nil, err
	}

	vetCmd := exec.Command("go", "vet", absPath)
	vetOutput, err := vetCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(vetOutput))
		diagnose(string(vetOutput), absPath, "modules/json/go.json", "Go")
		return vetOutput, err
	}

	runCmd := exec.Command("go", "run", absPath)
	output, err := runCmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		diagnose(string(output), absPath, "modules/json/go.json", "Go")
		return output, err
	}

	return output, nil
}
