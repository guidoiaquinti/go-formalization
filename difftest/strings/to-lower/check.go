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

const executable = "godiff-strings-to-lower"

func main() {
	const shard = 65_536
	edges := []string{"", "Hello, World", "\x00ABCXYZaz09\x7f", "ȺȾİK", "ΣΣσς", "ΟΣ", "Straßeẞ",
		"é É 🚀 �", "\xff", "\xe2\x82", "\xc0\x80", "\xed\xa0\x80", "\xf4\x90\x80\x80",
		"\xffA�\x80K", "AİȺZ", strings.Repeat("AZ.\x00", 128)}
	for i := range 256 {
		raw := string([]byte{byte(i)})
		edges = append(edges, raw, raw+"A�"+raw)
	}
	var nonidentity int
	for first := 0; first <= unicode.MaxRune; first += shard {
		last := min(first+shard-1, unicode.MaxRune)
		var input bytes.Buffer
		for r := first; r <= last; r++ {
			fmt.Fprintln(&input, 0, r)
		}
		for i, got := range check.LeanRows(executable, last-first+1, &input) {
			r := rune(first + i)
			want := unicode.ToLower(r)
			if got != strconv.Itoa(int(want)) {
				log.Fatalf("ToLower(%U): Lean=%s Go=%U", r, got, want)
			}
			if want != r {
				edges = append(edges, string(r), "ABC"+string(r)+"\xffZ�")
				nonidentity++
			}
		}
	}
	fmt.Printf("Unicode simple lowercase: %d exact lookups, %d nonidentity mappings.\n",
		unicode.MaxRune+1, nonidentity)
	check.CompareSamples(executable, edges, sample, encode, func(input string) string {
		return describe(strings.ToLower(input))
	}, func(input string) string {
		var out strings.Builder
		for _, r := range input {
			out.WriteRune(unicode.ToLower(r))
		}
		return describe(out.String())
	})
}

func sample(rng *rand.Rand, index int) string {
	tokens := []string{"a", "b", "Ǆ", "Ⱥ", "İ", "Σ", "K", "🚀", "�",
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

func describe(output string) string {
	row := []byte("bytes:")
	for i := range len(output) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(output[i]), 10)
	}
	return string(row)
}
