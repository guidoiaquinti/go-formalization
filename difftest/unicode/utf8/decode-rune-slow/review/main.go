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
		data []byte
		r    rune
		n    int
	}{
		{nil, utf8.RuneError, 0},
		{[]byte{}, utf8.RuneError, 0},
		{[]byte("\x00\xff"), 0, 1},
		{[]byte("\x80"), utf8.RuneError, 1},
		{[]byte("\xc2\x80"), 0x80, 2},
		{[]byte("\xc0\x80"), utf8.RuneError, 1},
		{[]byte("\xef\xbf\xbd"), utf8.RuneError, 3},
		{[]byte("\xf0\x90\x80"), utf8.RuneError, 1},
		{[]byte("\xf4\x8f\xbf\xbf"), 0x10ffff, 4},
	}
	var rows bytes.Buffer
	var dropWidthFailures, overlongFailures int
	for _, f := range fixtures {
		encoded := "nil"
		if f.data != nil {
			encoded = fmt.Sprintf("x%x", f.data)
		}
		fmt.Fprintf(&rows, "%s %d %d\n", encoded, f.r, f.n)
		r, n := utf8.DecodeRune(f.data)
		if r == utf8.RuneError && n > 1 {
			n = 1
		}
		if r != f.r || n != f.n {
			dropWidthFailures++
		}
		if bytes.Equal(f.data, []byte{0xc0, 0x80}) {
			r, n = 0, 2
		} else {
			r, n = utf8.DecodeRune(f.data)
		}
		if r != f.r || n != f.n {
			overlongFailures++
		}
	}
	if dropWidthFailures != 1 || overlongFailures != 1 {
		log.Fatalf("drop replacement width failed %d, accept overlong NUL failed %d", dropWidthFailures, overlongFailures)
	}
	dir, err := os.MkdirTemp("", "golean-decode-slice-slow-review-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, rows.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("GOLEAN_UTF8_DECODE_SLICE_SLOW_CASES", path); err != nil {
		log.Fatal(err)
	}
	const overlay = "difftest/unicode/utf8/decode-rune-slow/testdata/"
	if err := check.TestInPackage("unicode/utf8", "TestGoLeanDecodeRuneSlow", overlay+"bridge_test.go", overlay+"decode_test.go"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("private slow decoder passed %d fixtures; both mutants failed one\n", len(fixtures))
}
