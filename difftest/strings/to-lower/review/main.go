package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct{ input, want string }{
		{"", ""},
		{"AZaz!\x00", "azaz!\x00"},
		{"\u0130\u212a", "ik"},
		{"\U00010400", "\U00010428"},
		{"ΟΣ", "οσ"},
		{"ẞ", "ß"},
		{"ǅ", "ǆ"},
		{"\xffA\x80", "�a�"},
		{"A\xe2\x82", "a��"},
		{"�", "�"},
	}
	var rows bytes.Buffer
	pureASCII, preserveMalformed := 0, 0
	for i, c := range cases {
		if got := strings.ToLower(c.input); got != c.want {
			log.Fatalf("Go fixture %d: %q, want %q", i, got, c.want)
		}
		fmt.Fprint(&rows, 1)
		for j := range len(c.input) {
			fmt.Fprint(&rows, " ", c.input[j])
		}
		fmt.Fprintln(&rows)
		if strings.Map(func(r rune) rune {
			if 'A' <= r && r <= 'Z' {
				return r + 'a' - 'A'
			}
			return r
		}, c.input) != c.want {
			pureASCII++
		}
		if lowerPreservingInvalid(c.input) != c.want {
			preserveMalformed++
		}
	}
	for i, got := range check.LeanRows("godiff-strings-to-lower", len(cases), &rows) {
		bytes := []byte(cases[i].want)
		fields := make([]string, len(bytes))
		for j, b := range bytes {
			fields[j] = strconv.Itoa(int(b))
		}
		if want := "bytes:" + strings.Join(fields, ","); got != want {
			log.Fatalf("Lean fixture %d: %q, want %q", i, got, want)
		}
	}
	if pureASCII != 5 || preserveMalformed != 2 {
		log.Fatalf("mutation counts ASCII=%d invalid=%d, want 5 and 2", pureASCII, preserveMalformed)
	}
	fmt.Printf("strings.ToLower passed %d fixed Go and Lean fixtures; ASCII-only and malformed-preserving mutants failed %d and %d.\n",
		len(cases), pureASCII, preserveMalformed)
}

func lowerPreservingInvalid(s string) string {
	var out strings.Builder
	start := 0
	for next := range s {
		if next > start {
			writeOriginalOrLower(&out, s[start:next])
		}
		start = next
	}
	if start < len(s) {
		writeOriginalOrLower(&out, s[start:])
	}
	return out.String()
}

func writeOriginalOrLower(out *strings.Builder, unit string) {
	for _, r := range unit {
		if r == unicode.ReplacementChar && len(unit) == 1 {
			out.WriteString(unit)
		} else {
			out.WriteRune(unicode.ToLower(r))
		}
	}
}
