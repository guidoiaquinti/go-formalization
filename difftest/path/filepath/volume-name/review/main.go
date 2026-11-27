package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct{ input, windows string }{
		{"", ""},
		{"/single/root", ""},
		{"Z:relative", "Z:"},
		{"7:/not-a-letter", "7:"},
		{"/:/trailing", "\\:"},
		{"\xff:\\name", "\xff:"},
		{"é:/wide-drive", ""},
		{"//host/share/path/..", "\\\\host\\share"},
		{"//host/../dir", ""},
		{"//../share", ""},
		{"//host/.../dir", "\\\\host\\..."},
		{"//host//tail", "\\\\host\\"},
		{"//?/uNc/host/share/path", "\\\\?\\uNc\\host\\share"},
		{"//?/UNCx/host/share", "\\\\?\\UNCx"},
		{"//?/unc", "\\\\?\\unc"},
		{"//?/unc/", "\\\\?\\unc\\"},
		{"/??/C:/item", "\\??\\C:"},
		{"//./device/child", "\\\\.\\device"},
		{"//?/../item", ""},
		{"//host/\x00\xff/dir", "\\\\host\\\x00\xff"},
	}
	var input bytes.Buffer
	wants := make([]string, 0, len(cases))
	driveLetters, includeParent, keepSlash := 0, 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, 0)
		for _, b := range []byte(c.input) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		native := ""
		switch filepath.Separator {
		case '/':
		case '\\':
			native = c.windows
		default:
			log.Fatalf("unmodeled separator %q", filepath.Separator)
		}
		if got := filepath.VolumeName(c.input); got != native {
			log.Fatalf("Go fixture %d got %q, want %q", i, got, native)
		}
		wants = append(wants, "bytes:|"+wire(c.windows))
		if len(c.input) >= 2 && c.input[1] == ':' && c.windows != "" && !asciiLetter(c.input[0]) {
			driveLetters++
		}
		if c.windows == "" && strings.HasPrefix(c.input, "//") && strings.Contains(c.input, "/../") ||
			c.windows == "" && strings.HasPrefix(c.input, "//../") {
			includeParent++
		}
		if c.windows != "" && strings.Contains(c.input[:len(c.windows)], "/") {
			keepSlash++
		}
	}
	for i, row := range check.LeanRows("godiff-path-filepath-volume-name", len(cases), &input) {
		if row != wants[i] {
			log.Fatalf("Lean fixture %d got %q, want %q", i, row, wants[i])
		}
	}
	if driveLetters != 3 || includeParent != 3 || keepSlash != 11 {
		log.Fatalf("mutation counts letters=%d parents=%d slash=%d, want 3, 3 and 11", driveLetters, includeParent, keepSlash)
	}
	fmt.Printf("filepath.VolumeName passed %d native and two-profile Lean fixtures; enforcing drive letters, allowing parent volumes and retaining slashes failed %d, %d and %d.\n", len(cases), driveLetters, includeParent, keepSlash)
}

func asciiLetter(value byte) bool {
	return 'a' <= value && value <= 'z' || 'A' <= value && value <= 'Z'
}

func wire(volume string) string {
	words := make([]string, len(volume))
	for i := range len(volume) {
		words[i] = strconv.Itoa(int(volume[i]))
	}
	return "bytes:" + strings.Join(words, ",")
}
