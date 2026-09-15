// --------------------------------------------------------------
// wtf / Checkers / File / JS
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
		diagnoseJSError(string(checkOutput), absPath)
		return checkOutput, err
	}

	cmd := exec.Command("node", absPath)

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		diagnoseJSError(string(output), absPath)
		return output, err
	}

	return output, nil
}

func diagnoseJSError(output, path string) {
	data, err := checkers.LoadJSONData("modules/json/js.json")
	if err != nil {
		fmt.Printf("Failed to load js.json: %v\n", err)
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
				if loc := nearestLocation(locs, idx); loc != "" {
					fmt.Println(color.YellowString("Location: %s", loc))
				}
				fmt.Println(divider())
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println(color.RedString("No known fix found for this JavaScript error."))
		if loc := nearestLocation(locs, 0); loc != "" {
			fmt.Println(color.YellowString("Location: %s", loc))
		}
	}
}
