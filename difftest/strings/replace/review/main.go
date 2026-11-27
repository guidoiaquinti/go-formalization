package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		input, old, replacement string
		n                       int
		want                    string
	}{
		{"", "", "!", -1, "!"},
		{"ab", "", "!", 0, "ab"},
		{"é", "", "!", -3, "!é!"},
		{"é\xff", "", "!", 2, "!é!\xff"},
		{"\xe2\x82�", "", "\xff", math.MinInt, "\xff\xe2\xff\x82\xff�\xff"},
		{"🚀a", "", "X", 1, "X🚀a"},
		{"ababa", "aba", "X", -3, "Xba"},
		{"aaaaa", "aa", "X", -1, "XXa"},
		{"aaaaa", "aa", "X", 1, "Xaaa"},
		{"aaaaa", "aa", "X", 0, "aaaaa"},
		{"aaaaa", "aa", "", math.MaxInt, "a"},
		{"a\x00a", "b", "X", -1, "a\x00a"},
		{"a\x00a", "a", "a", -1, "a\x00a"},
		{"é\xffé", "\xa9\xff\xc3", "X", -2, "\xc3X\xa9"},
		{"éé", "\xa9", "", -1, "\xc3\xc3"},
		{"\x00a\x00", "\x00", "🚀", -2, "🚀a🚀"},
	}
	var input bytes.Buffer
	zero, overlap, byteGap := 0, 0, 0
	wants := make([]string, 0, len(cases))
	for i, c := range cases {
		fmt.Fprint(&input, c.n, " ", len(c.old), " ", len(c.replacement))
		for _, part := range []string{c.old, c.replacement, c.input} {
			for _, b := range []byte(part) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		if got := strings.Replace(c.input, c.old, c.replacement, c.n); got != c.want {
			log.Fatalf("Go fixture %d returned %x, want %x", i, got, c.want)
		}
		if c.n == 0 && strings.ReplaceAll(c.input, c.old, c.replacement) != c.want {
			zero++
		}
		if c.old != "" && c.n < 0 && overlapped(c.input, c.old, c.replacement) != c.want {
			overlap++
		}
		if c.old == "" && c.n < 0 && perByte(c.input, c.replacement) != c.want {
			byteGap++
		}
		fields := make([]string, len(c.want))
		for j := range len(c.want) {
			fields[j] = strconv.Itoa(int(c.want[j]))
		}
		wants = append(wants, "bytes:"+strings.Join(fields, ","))
	}
	for i, row := range check.LeanRows("godiff-strings-replace", len(cases), &input) {
		if row != wants[i] {
			log.Fatalf("Lean fixture %d returned %s, want %s", i, row, wants[i])
		}
	}
	if zero != 2 || overlap != 2 || byteGap != 2 {
		log.Fatalf("mutation counts zero=%d overlap=%d byte-boundary=%d, want 2, 2 and 2", zero, overlap, byteGap)
	}
	fmt.Printf("strings.Replace passed %d fixed Go and Lean fixtures; unlimited-zero, overlapping and byte boundary mutants failed %d, %d and %d.\n",
		len(cases), zero, overlap, byteGap)
}

func overlapped(input, old, replacement string) string {
	var result strings.Builder
	copied := 0
	for i := range len(input) {
		if strings.HasPrefix(input[i:], old) {
			if i >= copied {
				result.WriteString(input[copied:i])
			}
			result.WriteString(replacement)
			copied = i + len(old)
		}
	}
	result.WriteString(input[copied:])
	return result.String()
}

func perByte(input, replacement string) string {
	var result strings.Builder
	result.WriteString(replacement)
	for i := range len(input) {
		result.WriteByte(input[i])
		result.WriteString(replacement)
	}
	return result.String()
}
