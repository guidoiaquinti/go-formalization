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
	cases := []struct{ input, windows string }{
		{"", ""},
		{".", "."},
		{"a/../b//", "a/../b//"},
		{"C:\\a\\..\\b\\", "C:/a/../b/"},
		{"\\\\server\\share\\.\\x", "//server/share/./x"},
		{"\\?\\UNC\\server\\share", "/?/UNC/server/share"},
		{"x\\\\/y", "x///y"},
		{"\\\x00\xffa\xe2\x82\\", "/\x00\xffa\xe2\x82/"},
		{"é\\🚀", "é/🚀"},
		{"lookalike：＼／", "lookalike：＼／"},
	}
	var input bytes.Buffer
	wants := make([]string, 0, len(cases))
	clean, sanitize := 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, 0)
		for _, b := range []byte(c.input) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		want := c.input
		switch filepath.Separator {
		case '/':
		case '\\':
			want = c.windows
		default:
			log.Fatalf("unmodeled native separator %q", filepath.Separator)
		}
		if got := filepath.ToSlash(c.input); got != want {
			log.Fatalf("Go fixture %d returned %x, want %x", i, got, want)
		}
		if path.Clean(c.windows) != c.windows {
			clean++
		}
		if strings.ToValidUTF8(c.windows, "�") != c.windows {
			sanitize++
		}
		wants = append(wants, wire(c.input)+"|"+wire(c.windows))
	}
	for i, row := range check.LeanRows("godiff-path-filepath-to-slash", len(cases), &input) {
		if row != wants[i] {
			log.Fatalf("Lean fixture %d returned %q, want %q", i, row, wants[i])
		}
	}
	if clean != 6 || sanitize != 1 {
		log.Fatalf("mutation counts clean=%d sanitize=%d, want 6 and 1", clean, sanitize)
	}
	fmt.Printf("filepath.ToSlash passed %d native and two-profile Lean fixtures; cleaning and sanitizing malformed bytes failed %d and %d.\n", len(cases), clean, sanitize)
}

func wire(path string) string {
	words := make([]string, len(path))
	for i := range len(path) {
		words[i] = strconv.Itoa(int(path[i]))
	}
	return strings.Join(words, ",")
}
