package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// parseSignature parses a Java method descriptor and returns the number of parameter slots.
// According to the latest rule: 1 slot for long (J) and 1 slot for double (D).
func parseSignature(sig string) (int, error) {
	idxStart := strings.Index(sig, "(")
	idxEnd := strings.Index(sig, ")")
	if idxStart == -1 || idxEnd == -1 || idxEnd < idxStart {
		return 0, fmt.Errorf("invalid descriptor: %s", sig)
	}
	params := sig[idxStart+1 : idxEnd]
	count := 0
	inRef := false
	for i := 0; i < len(params); i++ {
		if inRef {
			if params[i] == ';' {
				inRef = false
				count++
			}
			continue
		}
		switch params[i] {
		case '[':
			// Array type: skip '[' and process the base type as one slot.
			continue
		case 'L':
			inRef = true
		default:
			// Primitives: B, C, D, F, I, J, S, Z
			// Per instruction: J (long) and D (double) are 1 slot each.
			count++
		}
	}
	return count, nil
}

func main() {
	// Matches ghelpers.MethodSignatures["..."] = ghelpers.GMeth{ParamSlots: X, ...}
	// or ghelpers.MethodSignatures[types.ClassName...] = ghelpers.GMeth{ParamSlots: X, ...}
	// We'll use a regex to find these assignments.
	re := regexp.MustCompile(`ghelpers\.MethodSignatures\[(.*?)\]\s*=\s*(?:ghelpers\.)?GMeth\{\s*ParamSlots:\s*(\d+)`)

	err := filepath.WalkDir("src/gfunction", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, m := range matches {
			key := m[1]
			declaredSlots := m[2]

			// The key might be a literal string or a constant.
			// If it's a literal string like "java/lang/Object.<init>()V", we can parse it.
			if strings.HasPrefix(key, "\"") && strings.HasSuffix(key, "\"") {
				sig := key[1 : len(key)-1]
				expectedSlots, err := parseSignature(sig)
				if err != nil {
					continue
				}

				if fmt.Sprintf("%d", expectedSlots) != declaredSlots {
					fmt.Printf("%s|%s|%s|%d\n", path, sig, declaredSlots, expectedSlots)
				}
			}
			// Handling constants like types.ClassNameBigDecimal + ".abs()Ljava/math/BigDecimal;"
			// would require more complex parsing, but for the basic string literals, this works.
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking files: %v\n", err)
		os.Exit(1)
	}
}
