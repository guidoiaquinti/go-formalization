package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
		input         string
		unix, windows bool
	}{
		{"", false, false},
		{"]", false, false},
		{"a\x00\xff]", false, false},
		{"[", true, true},
		{"a[", true, true},
		{"?", true, true},
		{"*", true, true},
		{"\\", true, false},
		{"a\\b", true, false},
		{"a\\*", true, true},
		{"a\\[", true, true},
		{"\xff?\xc3", true, true},
		{"＊？［＼", false, false},
	}
	var leanInput, fixtures bytes.Buffer
	rightBracket, missedBackslash, escape := 0, 0, 0
	for _, c := range cases {
		fmt.Fprint(&leanInput, 0)
		for _, b := range []byte(c.input) {
			fmt.Fprint(&leanInput, " ", b)
		}
		fmt.Fprintln(&leanInput)
		native := c.unix
		if runtime.GOOS == "windows" {
			native = c.windows
		}
		fmt.Fprintf(&fixtures, "x%x %t\n", c.input, native)
		if (c.windows || strings.Contains(c.input, "]")) != c.windows {
			rightBracket++
		}
		if c.unix != c.windows {
			missedBackslash++
		}
		if hasUnescaped(c.input) != c.windows {
			escape++
		}
	}
	for i, row := range check.LeanRows("godiff-path-filepath-has-meta", len(cases), &leanInput) {
		want := strconv.FormatBool(cases[i].unix) + "|" + strconv.FormatBool(cases[i].windows)
		if row != want {
			return fmt.Errorf("Lean fixture %d got %q, want %q", i, row, want)
		}
	}
	if rightBracket != 2 || missedBackslash != 2 || escape != 2 {
		return fmt.Errorf("mutation counts right bracket=%d missed backslash=%d escape=%d, want 2 each", rightBracket, missedBackslash, escape)
	}
	dir, err := os.MkdirTemp("", "golean-filepath-has-meta-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, fixtures.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_FILEPATH_HAS_META_CASES", file); err != nil {
		return err
	}
	const path = "difftest/path/filepath/has-meta/testdata/"
	if err := check.TestInPackage("path/filepath", "TestGoLeanHasMeta", path+"bridge_test.go", path+"has_meta_test.go"); err != nil {
		return err
	}
	fmt.Printf("filepath.hasMeta passed %d private native and dual-profile Lean fixtures; including ], dropping backslash, and interpreting escapes failed %d, %d and %d.\n", len(cases), rightBracket, missedBackslash, escape)
	return nil
}

func hasUnescaped(input string) bool {
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' {
			i++
			continue
		}
		if strings.IndexByte("*?[", input[i]) >= 0 {
			return true
		}
	}
	return false
}
