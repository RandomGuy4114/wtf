// --------------------------------------------------------------
// wtf / Checkers / File / Python
//  By: FormalBlaze
// --------------------------------------------------------------

package file

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"wtf/checkers"
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
		diagnosePythonError(string(output), absPath)
		return output, err
	}

	return output, nil
}

func diagnosePythonError(output, path string) {
	data, err := checkers.LoadJSONData("modules/json/python.json")
	if err != nil {
		fmt.Printf("Failed to load python.json: %v\n", err)
		return
	}

	lowerOutput := strings.ToLower(output)
	locs := extractLocations(output, path)
	fmt.Println(divider())
	found := false
	for _, e := range data.Errors {
		for _, s := range e.Strings {
			if idx := strings.Index(lowerOutput, strings.ToLower(s)); idx != -1 {
				CombinedString := fmt.Sprintf("%s: %s\nFix: %s\n", data.Name, e.Message, e.Fix)
				fmt.Print(color.RedString(CombinedString))
				loc := nearestLocation(locs, idx)
				if loc != "" {
					fmt.Println(color.YellowString("Location: %s", loc))
				}
				fmt.Println(divider())
				checkers.LogMessage(fmt.Sprintf("%s: %s | Fix: %s | Location: %s", data.Name, e.Message, e.Fix, loc))
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println(color.RedString("No known fix found for this Python error."))
		if loc := nearestLocation(locs, 0); loc != "" {
			fmt.Println(color.YellowString("Location: %s", loc))
		}
		checkers.LogMessage("no known fix found for Python error in " + path)
	}
}
