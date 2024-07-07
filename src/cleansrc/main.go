/*

cleansrc utility - For a given input source code file, either show any invalid characters without modifying the file or fix it,
depending on the command line options. Multiple input files can be specified as a list and/or wild-card specifications
on the command line.

Valid source code characters: printable, \n, \r, or \t.

*/

package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// Show help and then exit to the O/S
func showHelp() {
	suffix := filepath.Base(os.Args[0])
	fmt.Printf("\nUsage:  %s  [-h]  [-s <one-character> | -d]  [-u]  [-v]  Input file(s)...\n\nwhere\n", suffix)
	fmt.Printf("\t-h : This display.\n")
	fmt.Printf("\t-s <one-character>: Substitute the specified character for invalid source characters.\n")
	fmt.Printf("\t     CAUTION: This causes a rewrite of the input source code file.\n")
	fmt.Printf("\t     Examples:\n")
	fmt.Printf("\t\t-s A --> The letter A is the substitute.\n")
	fmt.Printf("\t\t-s space --> The space character (' ') is the substitute.\n")
	fmt.Printf("\t\t-s tab --> The tab character ('\\t') is the substitute.\n")
	fmt.Printf("\t\t-s nl --> The newline character ('\\n') is the substitute.\n")
	fmt.Printf("\t\t-s cr --> The carriage return character ('\\r') is the substitute.\n")
	fmt.Printf("\t\t-s 0x23 --> The pound sign character ('#') is the substitute.\n")
	fmt.Printf("\t-d : Instead of replacing invalid source characters, delete them.\n")
	fmt.Printf("\t-u : Update the source file.\n")
	fmt.Printf("\t-v : Verbose output.\n")
	fmt.Printf("\nExit codes:\n")
	fmt.Printf("\t0\tNormal completion.\n")
	fmt.Printf("\t1\tSomething went wrong during execution.\n\n")
	os.Exit(1)
}

// Is the argument byte valid?
func isValidSrcByte(argByte byte) bool {
	if unicode.IsPrint(rune(argByte)) {
		return true
	}
	if argByte == '\n' || argByte == '\r' || argByte == '\t' {
		return true
	}
	return false
}

/*
Main function
*/
func main() {

	var Args []string
	var flagSubstitute = false
	var flagDelete = false
	var flagUpdate = false
	var flagVerbose = false
	var nilByte = byte(0x00)
	var replByte = byte('?')
	var err error
	var inPathSet []string
	var inPath string
	modFileCount := 0
	errFileCount := 0

	// Parse command line arguments.
	for _, singleVar := range os.Args[1:] {
		Args = append(Args, singleVar)
	}
	if len(Args) < 1 {
		showHelp()
	}
	for ii := 0; ii < len(Args); ii++ {
		if !strings.HasPrefix(Args[ii], "-") {
			inPath, err = filepath.Abs(Args[ii])
			if err != nil {
				fmt.Printf(fmt.Sprintf("Input file (%s) cannot be accessed: %s", Args[ii], err.Error()))
				showHelp()
			}
			inPathSet = append(inPathSet, inPath)
			continue
		}
		switch Args[ii] {
		case "-d":
			flagDelete = true
		case "-h":
			showHelp()
		case "-s":
			flagSubstitute = true
			ii += 1
			if Args[ii] == "space" {
				replByte = ' '
				break
			}
			if Args[ii] == "tab" {
				replByte = '\t'
				break
			}
			if Args[ii] == "nl" {
				replByte = '\n'
				break
			}
			if Args[ii] == "cr" {
				replByte = '\r'
				break
			}
			if strings.HasPrefix(Args[ii], "0x") {
				// Hexstring for one character: 0xn or 0xnn
				if len(Args[ii]) > 4 {
					fmt.Printf(fmt.Sprintf("-s %s value too large for one character", Args[ii]))
					showHelp()
				}
				decoded, err := hex.DecodeString(Args[ii][2:])
				if err != nil {
					fmt.Printf(fmt.Sprintf("-s %s is not hexidecimal (0xH or 0xHH where H is from 0123456789ABCDEFabcdef)", Args[ii]))
					showHelp()
				}
				replByte = decoded[0]
				if replByte == nilByte {
					fmt.Printf("-s 0x00 ---> invalid source characters will not be replaced but discarded (squeezed out)")
				}
				break
			}
			if len(Args[ii]) == 1 {
				replByte = Args[ii][0]
				break
			}
			fmt.Printf("Use -s to specify only ONE character")
			showHelp()
		case "-u":
			flagUpdate = true
		case "-v":
			flagVerbose = true
		default:
			fmt.Printf(fmt.Sprintf("Unrecognizable command option: %s", Args[ii]))
			showHelp()
		}
	}

	// Check for contradiction.
	if flagDelete && flagSubstitute {
		fmt.Printf("-s and -d together makes no sense")
		showHelp()
	}

	// Make sure that at least one input file path specification is present.
	if len(inPathSet) < 1 {
		fmt.Printf("At least one input file specification is required")
		showHelp()
	}

	// For each input file path in the set, process it.
	for _, inPath := range inPathSet {

		// Read entire input file contents.
		inBytes, err := os.ReadFile(inPath)
		if err != nil {
			errFileCount++
			errMsg := fmt.Sprintf("os.ReadFile(%s) failed: %s", inPath, err.Error())
			fmt.Printf(errMsg)
			continue
		}
		if flagVerbose {
			fmt.Printf("Read %d input characters from %s", len(inBytes), inPath)
		}

		// Copy inBytes to outBytes.
		// Check outBytes for source code validity.
		// If a byte deviates, make a substitution or discard it.
		var outBytes = inBytes
		numTotalBytes := len(inBytes)
		numInvalidBytes := 0
		lineNum := 1
		offset := 0
		for ii := 0; ii < numTotalBytes; ii++ {
			if outBytes[ii] == '\n' {
				lineNum += 1
				offset = 0
				continue
			}
			if !isValidSrcByte(outBytes[ii]) {
				// Not a valid source code byte.
				nastyOne := outBytes[ii]
				numInvalidBytes += 1
				if flagDelete {
					// Discard byte.
					tempBytes := inBytes[:ii]
					outBytes = append(tempBytes, outBytes[ii+1:]...)
					numTotalBytes -= 1
					fmt.Printf("\t%s Line %d, offset %d: discarded 0x%x\n", inPath, lineNum, offset, nastyOne)
				} else {
					// Replace byte.
					outBytes[ii] = replByte
					fmt.Printf("\t%s Line %d, offset %d: replaced 0x%x with 0x%x\n", inPath, lineNum, offset, nastyOne, replByte)
				}
			}
			offset += 1
		}

		// If clean input file, continue to next file.
		if numInvalidBytes < 1 {
			fmt.Printf("Clean input source code file: %s", inPath)
			continue
		}
		if flagVerbose {
			// Report how many invalid source file characters were detected.
			fmt.Printf("Detected %d invalid source file character(s) in %s", numInvalidBytes, inPath)
		}

		// If not rewriting the input file, continue to next file.
		if !flagSubstitute && !flagDelete {
			continue
		}

		// Replace input file contents with outBytes.
		if flagUpdate {
			err = os.WriteFile(inPath, outBytes, 0666)
			if err != nil {
				errFileCount++
				errMsg := fmt.Sprintf("os.WriteFile (%s) failed: %s", inPath, err.Error())
				fmt.Printf(errMsg)
				continue
			}
			modFileCount += 1
			if flagVerbose {
				fmt.Printf("Wrote %d output characters to %s", len(outBytes), inPath)
			}
		}
	}

	// Done with input source files.
	// Return 1 to O/S if any errors were diagnosed.
	// Else, return 0.
	if errFileCount > 0 {
		fmt.Printf("%d input file(s) had I/O errors", errFileCount)
		os.Exit(1)
	}
	if modFileCount > 0 {
		fmt.Printf("%d input file(s) were modified", modFileCount)
	} else {
		fmt.Printf("No input file(s) were modified")
	}
	os.Exit(0)

}
