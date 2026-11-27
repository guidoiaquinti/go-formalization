package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const executable = "godiff-bytes-fold-runes"

type input struct{ left, right string }

func main() {
	const shard = 65_536
	var nonidentity, longestOrbit int
	var edges []input
	for first := 0; first <= unicode.MaxRune; first += shard {
		last := min(first+shard-1, unicode.MaxRune)
		var lookups bytes.Buffer
		for r := first; r <= last; r++ {
			fmt.Fprintln(&lookups, 0, r)
		}
		for i, got := range check.LeanRows(executable, last-first+1, &lookups) {
			r := rune(first + i)
			next := unicode.SimpleFold(r)
			if got != strconv.Itoa(int(next)) {
				log.Fatalf("SimpleFold(%U): Lean=%s Go=%U", r, got, next)
			}
			if next == r {
				continue
			}
			nonidentity++
			reachable := 1
			for candidate := next; candidate != r; candidate = unicode.SimpleFold(candidate) {
				edges = append(edges, input{string(r), string(candidate)})
				reachable++
				if reachable > 4 {
					log.Fatalf("fold orbit for %U exceeds model bound", r)
				}
			}
			longestOrbit = max(longestOrbit, reachable)
			edges = append(edges, input{string(r), "�"})
		}
	}
	fmt.Printf("SimpleFold: %d code points agree, %d nonidentity mappings, maximum orbit %d.\n",
		unicode.MaxRune+1, nonidentity, longestOrbit)

	for i := range 1 << 16 {
		edges = append(edges, input{string([]byte{byte(i >> 8)}), string([]byte{byte(i)})})
	}
	words := []string{"", "ab\x00Z", "aB\x00z", "KſΣ", "KSς", "kSσ", "ẞ", "ß", "SS",
		"İ", "i", "ı", "I", "ﬃ", "ffi", "µ", "μ", "Μ", "🚀", "�", "\xff", "\x80",
		"\xff\x80", "��", "\xe2\x82", "\xe2A", "�a", "\xed\xa0\x80A", "���a",
		"\xf4\x90\x80\x80", strings.Repeat("AZ", 128), strings.Repeat("az", 128)}
	for _, left := range words {
		for _, right := range words {
			edges = append(edges, input{left, right})
		}
	}
	check.CompareSamples(executable, edges, sample, encode, actual, expected)

}

func renderRunes(values []rune, capacity int) string {
	tag := "slice"
	if values == nil {
		tag = "nil"
	}
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = strconv.Itoa(int(value))
	}
	return fmt.Sprintf("%s:%d:%s", tag, capacity, strings.Join(fields, ","))
}

func actual(c input) string {
	left, right := []byte(c.left), []byte(c.right)
	// Both nil and nonnil empty inputs occur in the comparison.
	if c.left == "" && len(c.right)%2 == 0 {
		left = nil
	}
	if c.right == "" && len(c.left)%2 == 0 {
		right = nil
	}
	fold := bytes.EqualFold(left, right)
	decoded := func(data []byte) string {
		values := bytes.Runes(data)
		result := renderRunes(values, cap(values))
		for i := range values {
			values[i] = -1
		}
		return result
	}
	a, b := decoded(left), decoded(right)
	if string(left) != c.left || string(right) != c.right {
		panic("changed original bytes")
	}
	return fmt.Sprintf("%t|%s|%s", fold, a, b)
}

func expected(c input) string {
	left, right := []rune(c.left), []rune(c.right)
	return fmt.Sprintf("%t|%s|%s", canonical(c.left) == canonical(c.right),
		renderRunes(left, len(left)), renderRunes(right, len(right)))
}

// Build a canonical rune sequence rather than walking the source's
// pairwise ordered fold loop. Go range replaces each malformed byte.
func canonical(raw string) string {
	var output strings.Builder
	for _, r := range raw {
		least := r
		for next := unicode.SimpleFold(r); next != r; next = unicode.SimpleFold(next) {
			least = min(least, next)
		}
		output.WriteRune(least)
	}
	return output.String()
}

func encode(c input) string {
	row := fmt.Appendf(nil, "1 %d", len(c.left))
	for _, raw := range []string{c.left, c.right} {
		for i := range len(raw) {
			row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
		}
	}
	return string(row)
}

func sample(rng *rand.Rand, index int) input {
	tokens := []string{"a", "A", "s", "ſ", "k", "K", "é", "É", "σ", "Σ", "ς", "�", "\xff", "\xe2\x82", "🚀"}
	var left, right strings.Builder
	for range rng.Intn(22) {
		if rng.Intn(3) == 0 {
			left.WriteByte(byte(rng.Intn(256)))
		} else {
			left.WriteString(tokens[rng.Intn(len(tokens))])
		}
	}
	for _, r := range left.String() {
		if index%4 == 0 && rng.Intn(4) == 0 {
			right.WriteByte(byte(rng.Intn(256)))
		} else if r == 0xfffd && rng.Intn(2) == 0 {
			right.WriteByte(0xff)
		} else {
			for range rng.Intn(4) {
				r = unicode.SimpleFold(r)
			}
			right.WriteRune(r)
		}
	}
	if index%5 == 0 {
		right.WriteString(tokens[rng.Intn(len(tokens))])
	}
	return input{left.String(), right.String()}
}
