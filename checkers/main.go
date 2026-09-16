// --------------------------------------------------------------
// wtf / Checkers / Main
//  By: FormalBlaze
// --------------------------------------------------------------

package checkers

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed modules
var modulesFS embed.FS

var (
	pathsOnce sync.Once
	paths     Paths
	pathsErr  error
)

func loadPaths() (Paths, error) {
	pathsOnce.Do(func() {
		f, err := modulesFS.ReadFile("modules/PATHS.yml")
		if err != nil {
			pathsErr = err
			return
		}
		pathsErr = yaml.Unmarshal(f, &paths)
	})
	return paths, pathsErr
}

type JSONData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Errors      []struct {
		Type    string   `json:"type"`
		Strings []string `json:"strings"`
		Message string   `json:"message"`
		Fix     string   `json:"fix"`
	} `json:"errors"`
}

type Paths struct {
	Paths map[string]string `yaml:"paths"`
}

func LoadJSONData(filePath string) (*JSONData, error) {
	file, err := modulesFS.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	var jsonData JSONData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON data: %w", err)
	}

	return &jsonData, nil
}

// JudgeError scans output against every known error module and prints the
// fix for the first matching error. It returns true if a fix was found.
func JudgeError(output string) bool {
	if output == "" {
		log.Fatalf("No Error Given")
	}

	p, err := loadPaths()
	if err != nil {
		log.Fatal(err)
	}

	lowerOutput := strings.ToLower(output)

	for _, modulePath := range p.Paths {
		data, err := LoadJSONData(path.Join("modules", modulePath))
		if err != nil {
			log.Printf("skipping module %q: %v", modulePath, err)
			continue
		}
		for _, e := range data.Errors {
			for _, s := range e.Strings {
				if strings.Contains(lowerOutput, strings.ToLower(s)) {
					fmt.Println(divider())
					fmt.Printf("%s: %s\nFix: %s\n", data.Name, e.Message, e.Fix)
					fmt.Println(divider())
					LogMessage(fmt.Sprintf("%s: %s | Fix: %s", data.Name, e.Message, e.Fix))
					return true
				}
			}
		}
	}

	return false
}
