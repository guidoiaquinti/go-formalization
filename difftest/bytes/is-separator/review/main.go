package main

import (
	"bytes"
	"fmt"
	"log"
	"maps"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	spaces := []rune{0x85, 0xa0, 0x1680, 0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200a, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000}
	fixtures := map[rune]bool{}
	for r := rune(0); r <= 127; r++ {
		fixtures[r] = !('0' <= r && r <= '9' || 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || r == '_')
	}
	for _, r := range spaces {
		fixtures[r-1], fixtures[r], fixtures[r+1] = false, false, false
	}
	for _, r := range spaces {
		fixtures[r] = true
	}
	fixtures[-1], fixtures[math.MinInt32] = true, true
	for _, r := range []rune{0x180e, 0x200b, 0xd800, 0xfffd, '。', 'é', '٣', '🚀', 0x110000, math.MaxInt32} {
		fixtures[r] = false
	}
	var inputs, expected bytes.Buffer
	want := make([]string, 0, len(fixtures))
	underscore, asciiOnly, unsigned := 0, 0, 0
	for _, r := range slices.Sorted(maps.Keys(fixtures)) {
		value := fixtures[r]
		fmt.Fprintln(&inputs, r)
		fmt.Fprintf(&expected, "%d\t%t\n", r, value)
		want = append(want, strconv.FormatBool(value))
		if (value || r == '_') != value {
			underscore++
		}
		if (value && r <= 127) != value {
			asciiOnly++
		}
		if (value && r >= 0) != value {
			unsigned++
		}
	}
	if underscore != 1 || asciiOnly != 19 || unsigned != 2 {
		return fmt.Errorf("mutation counts: underscore=%d, ASCII-only=%d, unsigned=%d", underscore, asciiOnly, unsigned)
	}
	for i, row := range check.LeanRows("godiff-bytes-is-separator", len(want), &inputs) {
		if row != want[i] {
			return fmt.Errorf("Lean row %d = %s, want %s", i, row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-bytes-separator-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "cases")
	if err := os.WriteFile(file, expected.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_IS_SEPARATOR_CASES", file); err != nil {
		return err
	}
	const path = "difftest/bytes/is-separator/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanIsSeparator", path+"bridge_test.go", path+"is_separator_test.go"); err != nil {
		return err
	}
	fmt.Printf("bytes.isSeparator passed %d fixed Lean and private Go fixtures; underscore, ASCII-only and unsigned mutants failed %d, %d and %d.\n",
		len(want), underscore, asciiOnly, unsigned)
	return nil
}
