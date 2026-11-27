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

const executable = "godiff-bytes-simple-case"

func main() {
	const shard = 65_536
	edges := []string{"", "Hello, World", "\x00ABCXYZaz09\x7f", "ⱥɐİK", "ΣΣσς", "ος", "Straßeẞﬃ",
		"é É 🚀 �", "\xff", "\xe2\x82", "\xc0\x80", "\xed\xa0\x80", "\xf4\x90\x80\x80",
		"\xffA�\x80K", "AİȺZ", strings.Repeat("AZ.\x00", 128)}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, raw, raw+"A�"+raw)
	}
	var nonidentity [3]int
	for first := 0; first <= unicode.MaxRune; first += shard {
		last := min(first+shard-1, unicode.MaxRune)
		var input bytes.Buffer
		for r := first; r <= last; r++ {
			fmt.Fprintln(&input, 0, r)
		}
		for i, got := range check.LeanRows(executable, last-first+1, &input) {
			r := rune(first + i)
			upper, lower, title := unicode.ToUpper(r), unicode.ToLower(r), unicode.ToTitle(r)
			want := fmt.Sprintf("%d,%d,%d", upper, lower, title)
			if got != want {
				log.Fatalf("case mappings for %U: Lean=%s Go=%s", r, got, want)
			}
			for k, value := range []rune{upper, lower, title} {
				if value != r {
					nonidentity[k]++
				}
			}
			if upper != r || lower != r || title != r {
				edges = append(edges, string(r), "Abc"+string(r)+"\xffZ�")
			}
		}
	}
	fmt.Printf("Unicode simple upper/lower/title: %d exact lookups per mapping, %v nonidentity.\n",
		unicode.MaxRune+1, nonidentity)
	check.CompareSamples(executable, edges, sample, encode, actual, expected)

}

func sample(rng *rand.Rand, index int) string {
	tokens := []string{"a", "b", "Ǆ", "ⱥ", "İ", "ς", "K", "🚀", "�",
		"\x00", "\xff", "\xe2\x82", "\xed\xa0\x80", "\xc0\x80"}
	var out strings.Builder
	for range rng.Intn(40) {
		switch index % 4 {
		case 0:
			out.WriteByte(byte(rng.Intn(128)))
		case 1:
			out.WriteRune(rune(rng.Intn(unicode.MaxRune + 1)))
		case 2:
			out.WriteByte(byte(rng.Intn(256)))
		default:
			out.WriteString(tokens[rng.Intn(len(tokens))])
		}
	}
	return out.String()
}

func encode(input string) string {
	row := []byte("1")
	for i := range len(input) {
		row = strconv.AppendUint(append(row, ' '), uint64(input[i]), 10)
	}
	return string(row)
}

func describe(output string, isNil bool, capacity string) string {
	tag := "slice"
	if isNil {
		tag = "nil"
	}
	row := []byte(tag + ":" + capacity + ":")
	for i := range len(output) {
		if i > 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(output[i]), 10)
	}
	return string(row)
}

// The exact case specifies the clipped ASCII path or an allocation that
// never grows. The growth case includes copy-append allocator rounding.
func expectedCapacity(input, output string, kind int) string {
	ascii := true
	for i := range len(input) {
		if input[i] >= 128 {
			ascii = false
			break
		}
	}
	if kind < 2 && ascii {
		if input == output && len(input) > 0 {
			return "growth"
		}
		return strconv.Itoa(len(input))
	}
	if len(output) > len(input) {
		return "growth"
	}
	return strconv.Itoa(len(input))
}

func actual(input string) string {
	backing := make([]byte, len(input)+2)
	copy(backing, input)
	backing[len(input)], backing[len(input)+1] = 0xa5, 0x59
	source := backing[:len(input)]
	if input == "" {
		source = nil
	}
	var results []string
	for kind, convert := range []func([]byte) []byte{bytes.ToUpper, bytes.ToLower, bytes.ToTitle} {
		out := convert(source)
		text := string(out)
		capacity := expectedCapacity(input, text, kind)
		if cap(out) < len(input) || (capacity != "growth" && strconv.Itoa(cap(out)) != capacity) {
			log.Fatalf("case %d cap(%q→%q)=%d want %s", kind, input, text, cap(out), capacity)
		}
		results = append(results, describe(text, out == nil, capacity))
		for i := range out {
			out[i] ^= 0xff
		}
		if string(backing[:len(input)]) != input || backing[len(input)] != 0xa5 || backing[len(input)+1] != 0x59 {
			log.Fatal("case conversion shares source backing")
		}
		if input == "" {
			empty := convert(backing[:0])
			if empty == nil || cap(empty) != 0 || len(empty) != 0 {
				log.Fatal("nonnil empty conversion mismatch")
			}
		}
	}
	return strings.Join(results, "|")
}

func expected(input string) string {
	var results []string
	for kind, convert := range []func(rune) rune{unicode.ToUpper, unicode.ToLower, unicode.ToTitle} {
		var out strings.Builder
		for _, r := range input {
			out.WriteRune(convert(r))
		}
		output := out.String()
		results = append(results, describe(output, false, expectedCapacity(input, output, kind)))
	}
	return strings.Join(results, "|")
}
