package main

import (
	"fmt"
	"log"
	"runtime"
	"unicode/utf8"
)

func encode(dst []byte, value rune) (width int, panicked bool) {
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(runtime.Error); !ok {
				panic(p)
			}
			panicked = true
		}
	}()
	width = utf8.EncodeRune(dst, value)
	return
}

func main() {
	fixtures := []struct {
		rune   rune
		length int
		bytes  string
	}{
		{'x', 0, "____"},
		{'x', 4, "x___"},
		{0x80, 1, "____"},
		{0x80, 4, "\xc2\x80__"},
		{0x800, 2, "____"},
		{0xd800, 2, "____"},
		{0xd800, 3, "\xef\xbf\xbd_"},
		{0x10ffff, 3, "____"},
		{0x10ffff, 4, "\xf4\x8f\xbf\xbf"},
	}
	var partialFailures, wrongSizeFailures int
	for _, f := range fixtures {
		dst := []byte("____")
		encoded := string(f.rune)
		wantPanic := f.length < len(encoded)
		wantWidth := len(encoded)
		if wantPanic {
			wantWidth = 0
		}
		n, panicked := encode(dst[:f.length], f.rune)
		if panicked != wantPanic || n != wantWidth || string(dst) != f.bytes {
			log.Fatalf("EncodeRune(U+%04X, len=%d) returned (%d, %v, % x)", f.rune, f.length, n, panicked, dst)
		}
		if wantPanic {
			partial := []byte("____")
			copy(partial[:f.length], encoded)
			if string(partial) != f.bytes {
				partialFailures++
			}
		} else if wantWidth != f.length {
			wrongSizeFailures++
		}
	}
	if partialFailures != 4 || wrongSizeFailures != 2 {
		log.Fatalf("partial write mutant failed %d, return buffer length mutant failed %d", partialFailures, wrongSizeFailures)
	}
	fmt.Printf("%d fixtures passed, mutants failed %d and %d\n", len(fixtures), partialFailures, wrongSizeFailures)
}
