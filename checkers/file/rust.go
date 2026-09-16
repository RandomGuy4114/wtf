// --------------------------------------------------------------
// wtf / Checkers / File / Rust
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CheckRustErrors(path string) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve path: %v\n", err)
		return nil, err
	}

	tmpBinary := filepath.Join(os.TempDir(), fmt.Sprintf("wtf-rust-%d", os.Getpid()))
	defer os.Remove(tmpBinary)

	compileCmd := exec.Command("rustc", absPath, "-o", tmpBinary)
	compileOutput, err := compileCmd.CombinedOutput()
	fmt.Println(string(compileOutput))

	if err != nil {
		diagnose(string(compileOutput), absPath, "modules/json/rust.json", "Rust")
		return compileOutput, err
	}

	runCmd := exec.Command(tmpBinary)
	output, err := runCmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		diagnose(string(output), absPath, "modules/json/rust.json", "Rust")
		return output, err
	}

	return output, nil
}
