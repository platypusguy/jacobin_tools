package main

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const ExpectedMagicNumber = 0x4A4D

// Walk the base JMOD file and count the bytes for all classes/*.class entries
func main() {

	// Form the full path to the jmods directory
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		log.Fatal("os.GetEnv failed to find JAVA_HOME")
	}
	dirPath := javaHome + string(os.PathSeparator) + "jmods"

	// Open jmods directory
	dirOpened, err := os.Open(dirPath)
	if err != nil {
		log.Fatal("os.Open(jmods directory) failed", dirPath, err)
	}

	// Get all the file entries
	names, err := dirOpened.Readdirnames(0) // get all entries
	if err != nil {
		log.Fatal("Readdirnames(jmods directory) failed", dirPath, err)
	}

	// For each jmods file, process it
	for index := range names {
		name := names[index]
		fullPath := filepath.Join(dirPath, name)
		processJmodsFile(name, fullPath)
	}

}

func processJmodsFile(baseName string, fullPath string) {
	// Open the jmod file
	_, err := os.Open(fullPath)
	if err != nil {
		log.Fatal("os.Open failed:", fullPath, err)
	}

	// Read entire file contents
	jmodBytes, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatal("os.ReadFile failed:", fullPath, err)
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
	if fileMagicNumber != ExpectedMagicNumber {
		log.Fatal("fileMagicNumber != ExpectedMagicNumber:", fullPath, err)
	}

	// Skip over the jmod header so that it is recognized as a ZIP file
	offsetReader := bytes.NewReader(jmodBytes[4:])

	// Prepare the reader for the zip archive
	zipReader, err := zip.NewReader(offsetReader, int64(len(jmodBytes)-4))
	if err != nil {
		log.Fatal("zip.NewReader failed:", fullPath, err)
	}

	// For each file entry within the zip reader, process it
	countFiles := 0
	for _, fileEntry := range zipReader.File {

		// Has the right prefix and suffix?
		if !strings.HasPrefix(fileEntry.Name, "classes/") {
			continue
		}
		if !strings.HasSuffix(fileEntry.Name, ".class") {
			continue
		}

		// Remove the "classes/" prefix.
		classFileName := strings.Replace(fileEntry.Name, "classes/", "", 1)

		// Form array splut = [ "a", "b", etc ] based on a/b/etc
		splut := strings.Split(classFileName, string(os.PathSeparator))
		lenSplut := len(splut)
		if lenSplut < 2 {
			continue
		}

		className := splut[0]
		for ii := 1; ii < lenSplut; ii++ {
			className += string(os.PathSeparator) + splut[ii]
		}

		// Log this a.b.c
		fmt.Printf("%s , %s\n", baseName, className)

		// Add byte count to total
		countFiles += 1
	}

}
