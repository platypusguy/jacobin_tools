/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Show help and then exit to the O/S
func showHelp() {
	suffix := filepath.Base(os.Args[0])
	fmt.Printf("\nUsage:  %s  [-h]  {Input GOB file}\n\nwhere\n", suffix)
	fmt.Printf("\t-h : This display.\n")
	fmt.Printf("\nExit codes:\n")
	fmt.Printf("\t0\tNormal completion.\n")
	fmt.Printf("\t1\tSomething went wrong during execution or help requested.\n\n")
	os.Exit(1)
}

// read the gob file of the JDK jmod file and print it as a CSV file to stdout
// normally, the output is redirected to a CSV file

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
	if countArgs != 2 {
		fmt.Printf("*** ERROR: Too many command line arguments (%d), should be only 1\n", countArgs-1)
		showHelp()
	}

	// Get the full path of the gob file.
	gobFile, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Printf("*** ERROR: Input file (%s) cannot be accessed: %s\n", os.Args[1], err.Error())
		showHelp()
	}

	// Read the entire gob file
	data, err := os.ReadFile(gobFile)
	if err != nil {
		fmt.Printf("*** ERROR: os.ReadFile(%s) failed, reason: %s\n", gobFile, err.Error())
		os.Exit(1)
	}

	// Unmarshal the data into a struct
	var dataStruct map[string]string

	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&dataStruct)
	if err != nil {
		fmt.Printf("*** ERROR: gob.NewDecoder(%s) failed, reason: %s\n", gobFile, err.Error())
		os.Exit(1)
	}

	s := fmt.Sprintf("%v", dataStruct)
	s = strings.ReplaceAll(s, " ", "\n")
	s = strings.ReplaceAll(s, ":", ",")

	// Print the data
	fmt.Println(s)
}
