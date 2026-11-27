package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixtures := []struct {
		s string
		r rune
		n int
	}{
		{"", utf8.RuneError, 0},
		{"\x00\xff", 0, 1},
		{"\x7f", 0x7f, 1},
		{"\xff", utf8.RuneError, 1},
		{"\xc2\xa2", 0xa2, 2},
		{"\xc0\x80", utf8.RuneError, 1},
		{"\xe0\x9f\xbf", utf8.RuneError, 1},
		{"\xed\xa0\x80", utf8.RuneError, 1},
		{"\xe2\x82", utf8.RuneError, 1},
		{"\xef\xbf\xbd", utf8.RuneError, 3},
		{"\xf4\x8f\xbf\xbf", 0x10ffff, 4},
	}
	var rows bytes.Buffer
	var rejectASCIIFailures, collapseReplacementFailures int
	for _, f := range fixtures {
		fmt.Fprintf(&rows, "x%x %d %d\n", f.s, f.r, f.n)
		r, n := utf8.DecodeRuneInString(f.s)
		if len(f.s) != 0 && f.s[0] < utf8.RuneSelf {
			r = utf8.RuneError
		}
		if r != f.r || n != f.n {
			rejectASCIIFailures++
		}
		r, n = utf8.DecodeRuneInString(f.s)
		if r == utf8.RuneError && n > 1 {
			n = 1
		}
		if r != f.r || n != f.n {
			collapseReplacementFailures++
		}
	}
	if rejectASCIIFailures != 2 || collapseReplacementFailures != 1 {
		log.Fatalf("reject ASCII failed %d, collapse replacement width failed %d", rejectASCIIFailures, collapseReplacementFailures)
	}
	dir, err := os.MkdirTemp("", "golean-decode-string-slow-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_UTF8_DECODE_STRING_SLOW_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/unicode/utf8/decode-rune-in-string-slow/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanDecodeRuneInStringSlow", overlay+"bridge_test.go", overlay+"decode_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private slow decoder passed %d fixtures; mutants failed %d and %d\n", len(fixtures), rejectASCIIFailures, collapseReplacementFailures)
}
