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
		r      rune
		prefix string
		suffix string
	}{
		{0, "", "\xc0\x80"},
		{'A', "\xff", "\xc1\x81"},
		{0x7f, "!", "\xc1\xbf"},
		{0x80, "", "\xc2\x80"},
		{0x800, "\x00", "\xe0\xa0\x80"},
		{0x10000, "x", "\xf0\x90\x80\x80"},
		{-1, "x", "\xef\xbf\xbd"},
		{0xd800, "", "\xef\xbf\xbd"},
		{0x110000, "\xff", "\xef\xbf\xbd"},
	}
	var rows bytes.Buffer
	var asciiFailures, droppedInvalid int
	for _, f := range fixtures {
		want := f.prefix + f.suffix
		fmt.Fprintf(&rows, "%d x%x %x\n", f.r, f.prefix, want)
		if string(utf8.AppendRune([]byte(f.prefix), f.r)) != want {
			asciiFailures++
		}
		if !utf8.ValidRune(f.r) && f.prefix != want {
			droppedInvalid++
		}
	}
	if asciiFailures != 3 || droppedInvalid != 3 {
		log.Fatalf("public ASCII fast path failed %d, skip-invalid mutant failed %d", asciiFailures, droppedInvalid)
	}
	dir, err := os.MkdirTemp("", "golean-append-non-ascii-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_UTF8_APPEND_NON_ASCII_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/unicode/utf8/append-rune-non-ascii/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanAppendRuneNonASCII", overlay+"append_test.go", overlay+"bridge_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private Go helper passed %d fixtures; both mutations failed three\n", len(fixtures))
}
