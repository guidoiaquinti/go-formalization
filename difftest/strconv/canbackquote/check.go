package main

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func decimalBytes(s string) string {
	fields := make([]string, len(s))
	for i := range len(s) {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return strings.Join(fields, " ")
}

func main() {
	edges := []string{"", "a\tb", "a\nb", "a\rb", "a`b", "\u0080", "\u007f", "\ufeff", "a\ufeffb", "\ufffd", "\U0010ffff", "\xc0\x80", "\xed\xa0\x80", "\xf4\x90\x80\x80", "\xe2\x82", "\xf0\x9f\x98\x80"}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}))
	}
	check.CompareSamples("godiff-canbackquote", edges, func(r *rand.Rand, i int) string {
		if i%3 == 0 {
			return "prefix" + edges[r.Intn(len(edges))] + "suffix"
		}
		if i%3 == 1 {
			var s strings.Builder
			for n := r.Intn(12); n > 0; n-- {
				s.WriteRune(rune(r.Intn(utf8.MaxRune + 1)))
			}
			return s.String()
		}
		raw := make([]byte, r.Intn(33))
		r.Read(raw)
		return string(raw)
	}, func(s string) string {
		return "0 " + decimalBytes(s)
	}, func(s string) string {
		return strconv.FormatBool(strconv.CanBackquote(s))
	}, func(s string) string {
		ok := utf8.ValidString(s) && !strings.ContainsAny(s, "`\u007f\ufeff") &&
			strings.IndexFunc(s, func(c rune) bool { return c < ' ' && c != '\t' }) == -1
		return strconv.FormatBool(ok)
	})
}
