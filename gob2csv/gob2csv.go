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
	"io/ioutil"
	"os"
	"strings"
)

// read the gob file of the JDK jmod file and print it as a CSV file to stdout
// normally, the output is redirected to a CSV file

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Converts Jacobin gob file to csv file, displayed to stdout")
		fmt.Println("Usage: main filename.gob")
		return
	}

	for range os.Args {
		fmt.Println(os.Args)
	}

	gobFile := os.Args[1]
	// Open the .gob file
	file, err := os.Open(gobFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Read the contents of the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the data into a struct
	var dataStruct map[string]string

	err = gob.NewDecoder(bytes.NewBuffer(data)).Decode(&dataStruct)
	if err != nil {
		fmt.Println(err)
		return
	}

	s := fmt.Sprintf("%v", dataStruct)
	s = strings.ReplaceAll(s, " ", "\n")
	s = strings.ReplaceAll(s, ":", ",")

	// Print the data
	fmt.Println(s)
}