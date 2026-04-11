package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// 1. Run the verification script to get the mismatches.
	// We expect verify_paramslots.go to be present in the root.
	cmd := exec.Command("go", "run", "verify_paramslots.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running verification: %v\nOutput: %s\n", err, string(output))
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		fmt.Println("No mismatches found.")
		return
	}

	for _, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}
		path := parts[0]
		sig := parts[1]
		declaredSlots := parts[2]
		expectedSlots := parts[3]

		fmt.Printf("Fixing %s in %s: %s -> %s\n", sig, path, declaredSlots, expectedSlots)

		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			continue
		}

		// Create the search and replace strings.
		// We'll search for the specific signature assignment line to avoid errors.
		searchStr := fmt.Sprintf(`ghelpers.MethodSignatures["%s"] =
		ghelpers.GMeth{ParamSlots: %s,`, sig, declaredSlots)

		replaceStr := fmt.Sprintf(`ghelpers.MethodSignatures["%s"] =
		ghelpers.GMeth{ParamSlots: %s,`, sig, expectedSlots)

		// Check if searchStr exists in the file.
		newContent := strings.Replace(string(content), searchStr, replaceStr, 1)

		// If it's not found with that exact spacing/format, we'll try a simpler version.
		if newContent == string(content) {
			searchStr = fmt.Sprintf(`ghelpers.MethodSignatures["%s"] = ghelpers.GMeth{ParamSlots: %s,`, sig, declaredSlots)
			replaceStr = fmt.Sprintf(`ghelpers.MethodSignatures["%s"] = ghelpers.GMeth{ParamSlots: %s,`, sig, expectedSlots)
			newContent = strings.Replace(string(content), searchStr, replaceStr, 1)
		}

		if newContent == string(content) {
			fmt.Printf("Could not find signature assignment in %s for %s with slots %s\n", path, sig, declaredSlots)
			continue
		}

		err = os.WriteFile(path, []byte(newContent), 0644)
		if err != nil {
			fmt.Printf("Error writing to %s: %v\n", path, err)
		}
	}
}
