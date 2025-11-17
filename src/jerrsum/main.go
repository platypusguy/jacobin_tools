package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Show help and then exit to the O/S
func showHelp() {
	arg0 := filepath.Base(os.Args[0])
	fmt.Printf("\nUsage:  %s  [-h]  {Input summary report}  {output summary report}\n\nwhere\n", arg0)
	fmt.Printf("\nExit codes:\n")
	fmt.Printf("\t0\tNormal completion.\n")
	fmt.Printf("\t1\tSomething went wrong during execution or help requested.\n\n")
	os.Exit(1)
}

var headerRe = regexp.MustCompile(`^===== (.+?) =====$`)

type row struct {
	category string
	test     string
	reason   string
}

func WriteFailuresCSV(inputPath, outputPath string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	scanner := bufio.NewScanner(in)

	var (
		currentCategory string
		rows            []row
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect category section
		if m := headerRe.FindStringSubmatch(line); m != nil {
			currentCategory = m[1]
			continue
		}

		// End of section
		if strings.HasPrefix(line, "--- Total") {
			currentCategory = ""
			continue
		}

		// Skip irrelevant lines
		if currentCategory == "" || line == "" {
			continue
		}

		// Parse: "name: reason"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		testName := strings.TrimSpace(parts[0])
		reason := strings.TrimSpace(parts[1])

		rows = append(rows, row{
			category: currentCategory,
			test:     testName,
			reason:   reason,
		})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// ---- Sorting ----
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].category == rows[j].category {
			return rows[i].test < rows[j].test
		}
		return rows[i].category < rows[j].category
	})

	// ---- Write CSV ----
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	w := csv.NewWriter(out)
	defer w.Flush()

	// Heading row
	w.Write([]string{"ERROR CATEGORY", "TEST CASE", "TRIGGERING MESSAGE"})

	for _, r := range rows {
		// Let csv.Writer handle quoting, separators, and escapes
		if err := w.Write([]string{r.category, r.test, r.reason}); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Count of command-line arguments = 1 for the program name + input from operator.
	countArgs := len(os.Args)

	// Help requested?
	if countArgs < 2 {
		showHelp()
	}
	if countArgs > 1 && os.Args[1] == "-h" {
		showHelp()
	}

	// Make sure there is only 1 command-line argument from the operator.
	if countArgs != 3 {
		fmt.Printf("*** ERROR: Wrong number of command line arguments (%d), should be exactly 2\n", countArgs-1)
		showHelp()
	}

	err := WriteFailuresCSV(os.Args[1], os.Args[2])
	if err != nil {
		println("*** ERROR: ", err.Error())
	}
}
