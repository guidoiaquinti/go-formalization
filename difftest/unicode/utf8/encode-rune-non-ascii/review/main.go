package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixtures := []struct {
		r      rune
		length int
		prefix string
	}{
		{'A', 1, "\xc1\x81"},
		{'A', 2, "\xc1\x81"},
		{0, 3, "\xc0\x80"},
		{0x80, 2, "\xc2\x80"},
		{0x800, 2, "\xe0\xa0\x80"},
		{0xd800, 2, "\xef\xbf\xbd"},
		{-1, 3, "\xef\xbf\xbd"},
		{0x10ffff, 3, "\xf4\x8f\xbf\xbf"},
		{0x10ffff, 4, "\xf4\x8f\xbf\xbf"},
	}
	var rows bytes.Buffer
	var asciiFailures, partialFailures int
	for _, f := range fixtures {
		dst := bytes.Repeat([]byte{'_'}, f.length)
		want := bytes.Clone(dst)
		status := "bounds"
		if len(f.prefix) <= f.length {
			copy(want, f.prefix)
			status = strconv.Itoa(len(f.prefix))
		} else {
			partial := bytes.Clone(dst)
			copy(partial, f.prefix)
			if !bytes.Equal(partial, want) {
				partialFailures++
			}
		}
		if f.r < utf8.RuneSelf && f.r >= 0 {
			public := bytes.Clone(dst)
			width := utf8.EncodeRune(public, f.r)
			if strconv.Itoa(width) != status || !bytes.Equal(public, want) {
				asciiFailures++
			}
		}
		fmt.Fprintf(&rows, "%d x%x %s x%x\n", f.r, dst, status, want)
	}
	if asciiFailures != 3 || partialFailures != 4 {
		log.Fatalf("public fast path failed %d, partial-write mutant failed %d", asciiFailures, partialFailures)
	}
	dir, err := os.MkdirTemp("", "golean-encode-non-ascii-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_UTF8_ENCODE_NON_ASCII_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/unicode/utf8/encode-rune-non-ascii/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanEncodeRuneNonASCII", overlay+"bridge_test.go", overlay+"encode_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private Go helper passed %d fixtures; mutants failed %d and %d\n", len(fixtures), asciiFailures, partialFailures)
}
