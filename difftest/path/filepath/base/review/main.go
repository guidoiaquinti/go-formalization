package main

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct{ input, unix, windows string }{
		{"", ".", "."},
		{"///", "/", "\\"},
		{"\\\\", "\\\\", "\\"},
		{"7:", "7:", "\\"},
		{"7:tail", "7:tail", "tail"},
		{"a\\b\\", "a\\b\\", "b"},
		{"C:\\a/b\\", "b\\", "b"},
		{"//?/UNC/host/share//", "share", "\\"},
		{"//?/unc/host/share/leaf//", "leaf", "leaf"},
		{"//host/../", "..", ".."},
		{"//host/.../", "...", "\\"},
		{"//host/.././", ".", "."},
		{"a/../", "..", ".."},
		{"a/b/./", ".", "."},
		{"a/\xff\x00//", "\xff\x00", "\xff\x00"},
		{"C:/a/..\\", "..\\", ".."},
	}
	var input bytes.Buffer
	wants := make([]string, 0, 2*len(cases))
	wrongBackslash, cleaned, driveLetters := 0, 0, 0
	for i, c := range cases {
		for profile, result := range []string{c.unix, c.windows} {
			fmt.Fprint(&input, profile)
			for _, b := range []byte(c.input) {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
			wants = append(wants, wire(result))
		}
		native := c.unix
		switch filepath.Separator {
		case '/':
		case '\\':
			native = c.windows
		default:
			log.Fatalf("unmodeled separator")
		}
		if got := filepath.Base(c.input); got != native {
			log.Fatalf("Go fixture %d got %q, want %q", i, got, native)
		}
		if path.Base(strings.ReplaceAll(c.input, "\\", "/")) != c.unix {
			wrongBackslash++
		}
		if path.Base(path.Clean(c.input)) != c.unix {
			cleaned++
		}
		if strings.HasPrefix(c.input, "7:") && c.input != c.windows {
			driveLetters++
		}
	}
	for i, row := range check.LeanRows("godiff-path-filepath-base", len(wants), &input) {
		if row != wants[i] {
			log.Fatalf("Lean fixture %d got %q, want %q", i, row, wants[i])
		}
	}
	if wrongBackslash != 4 || cleaned != 4 || driveLetters != 2 {
		log.Fatalf("mutation counts backslash=%d cleaned=%d letters=%d, want 4, 4 and 2", wrongBackslash, cleaned, driveLetters)
	}
	fmt.Printf("filepath.Base passed %d native and %d Lean fixtures; treating backslash everywhere, cleaning and requiring drive letters failed %d, %d and %d.\n", len(cases), len(wants), wrongBackslash, cleaned, driveLetters)
}

func wire(value string) string {
	words := make([]string, len(value))
	for i := range len(value) {
		words[i] = strconv.Itoa(int(value[i]))
	}
	return "bytes:" + strings.Join(words, ",")
}
