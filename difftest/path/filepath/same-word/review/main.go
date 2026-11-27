package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cases := []struct {
		left, right   string
		unix, windows bool
	}{
		{"", "", true, true},
		{"A\x00", "a\x00", false, true},
		{"\x00", "", false, false},
		{"K", "K", false, true},
		{"ſΣ", "Sς", false, true},
		{"Σ", "ς", false, true},
		{"ẞ", "ß", false, true},
		{"ß", "SS", false, false},
		{"é", "e\u0301", false, false},
		{"C:/Item", "c:\\Item", false, false},
		{"a/../B", "B", false, false},
		{"\xffa", "�A", false, true},
		{"\xff\x80", "�", false, false},
		{"\xff\x80", "��", false, true},
		{"\xff\x80", "\xff\x80", true, true},
		{"[", "{", false, false},
	}
	var input, fixtures bytes.Buffer
	wants := make([]string, 0, len(cases))
	byteEquality, lowercase, slash := 0, 0, 0
	for _, c := range cases {
		fmt.Fprint(&input, len(c.left))
		for _, part := range []string{c.left, c.right} {
			for _, b := range []byte(part) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		native := c.unix
		if filepath.Separator == '\\' {
			native = c.windows
		}
		fmt.Fprintf(&fixtures, "x%x x%x %t\n", c.left, c.right, native)
		wants = append(wants, strconv.FormatBool(c.unix)+"|"+strconv.FormatBool(c.windows))
		if c.unix != c.windows {
			byteEquality++
		}
		if (strings.ToLower(c.left) == strings.ToLower(c.right)) != c.windows {
			lowercase++
		}
		if strings.EqualFold(strings.ReplaceAll(c.left, "\\", "/"),
			strings.ReplaceAll(c.right, "\\", "/")) != c.windows {
			slash++
		}
	}
	for i, row := range check.LeanRows("godiff-path-filepath-same-word", len(cases), &input) {
		if row != wants[i] {
			return fmt.Errorf("Lean fixture %d got %q, want %q", i, row, wants[i])
		}
	}
	if byteEquality != 7 || lowercase != 2 || slash != 1 {
		return fmt.Errorf("mutation counts bytes=%d lowercase=%d slash=%d, want 7, 2, 1", byteEquality, lowercase, slash)
	}
	dir, err := os.MkdirTemp("", "golean-filepath-same-word-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixtures.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_SAME_WORD_CASES", file); err != nil {
		return err
	}
	const data = "difftest/path/filepath/same-word/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanSameWord", data+"bridge_test.go", data+"same_word_test.go"); err != nil {
		return err
	}
	fmt.Printf("filepath.sameWord passed %d fixed private Go and dual-profile Lean cases; raw equality, lowercase and slash normalization failed %d, %d and %d.\n",
		len(cases), byteEquality, lowercase, slash)
	return nil
}
