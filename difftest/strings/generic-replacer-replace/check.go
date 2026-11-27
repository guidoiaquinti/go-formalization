package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	pairs []string
	text  string
}

func main() {
	file, err := os.CreateTemp("", "golean-generic-replacer-replace-")
	must(err)
	defer os.Remove(file.Name())
	writer := bufio.NewWriter(file)
	fixture := json.NewEncoder(writer)
	compare("boundaries", boundaries(), fixture)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5000)
			for i := range cases {
				cases[i] = sample(rng, i)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixture)
		}
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_GENERIC_REPLACER_REPLACE_CASES", file.Name()))
	const data = "difftest/strings/generic-replacer-replace/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanGenericReplacerReplace", data+"bridge_test.go", data+"replace_test.go"))
}

func boundaries() []input {
	var cases []input
	for _, pairs := range [][]string{
		nil, {"", ""}, {"", "|"}, {"", "|", "a", "X"}, {"a", "X", "", "|"},
		{"", "first", "", "second"}, {"ab", "L", "a", "S"}, {"a", "", "ab", "L"},
		{"a", "first", "a", "second"}, {"\x80", "\xff\x00", "\xff", "", "\xe2\x82", "!"},
		{"ǆ", "z", "", "."}, {"aa", "A", "aba", "B", "abaaba", "a"},
	} {
		for _, text := range []string{"", "a", "ab", "abc", "abba", "xabab", "aaaaa", "ababa", "z", "ǆ🚀", "\xff\x80\xe2\x82", "\x00\xff"} {
			cases = append(cases, input{pairs, text})
		}
	}
	for b := range 256 {
		key := string([]byte{byte(b)})
		text := key + "!" + key + key
		cases = append(cases, input{[]string{"", "", key + key, "\x00\xff", key, ""}, text},
			input{[]string{key, "L", key + key, "S", "", "\x00"}, text})
	}
	return cases
}

func sample(rng *rand.Rand, index int) input {
	alphabet := [...]int{2, 7, 256}[index%3]
	raw := func(limit int) string {
		value := make([]byte, rng.Intn(limit+1))
		for i := range value {
			value[i] = byte(rng.Intn(alphabet))
		}
		return string(value)
	}
	var pairs []string
	for range rng.Intn(12) {
		key := raw(12)
		if len(pairs) > 0 {
			base := pairs[2*rng.Intn(len(pairs)/2)]
			switch rng.Intn(4) {
			case 0:
				key = base[:rng.Intn(len(base)+1)]
			case 1:
				key = base + key
			}
		}
		pairs = append(pairs, key, raw(17))
	}
	text := raw(41)
	if len(pairs) > 0 && index%3 != 0 {
		key := pairs[2*rng.Intn(len(pairs)/2)]
		pos := rng.Intn(len(text) + 1)
		text = text[:pos] + key + text[pos:] + key
	}
	return input{pairs, text}
}

// Find the first eligible pair in constructor order and write directly
// to the result, consuming raw bytes or emitting one empty hit per offset.
func expected(c input) string {
	var out strings.Builder
	for position := 0; position <= len(c.text); {
		emittedEmpty := false
		for {
			match := -1
			for i := 0; i < len(c.pairs); i += 2 {
				key := c.pairs[i]
				if (!emittedEmpty || key != "") && strings.HasPrefix(c.text[position:], key) {
					match = i
					break
				}
			}
			if match < 0 {
				if position < len(c.text) {
					out.WriteByte(c.text[position])
				}
				position++
				break
			}
			out.WriteString(c.pairs[match+1])
			if width := len(c.pairs[match]); width != 0 {
				position += width
				break
			}
			emittedEmpty = true
		}
	}
	return out.String()
}

func compare(label string, cases []input, fixture *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d", len(c.pairs)/2)
		for i := 0; i < len(c.pairs); i++ {
			raw := c.pairs[i]
			if i%2 == 0 {
				fmt.Fprintf(&wire, " %d %d", len(raw), len(c.pairs[i+1]))
			}
			for j := range len(raw) {
				fmt.Fprintf(&wire, " %d", raw[j])
			}
		}
		for j := range len(c.text) {
			fmt.Fprintf(&wire, " %d", c.text[j])
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-strings-generic-replacer-replace", len(cases), &wire) {
		c := cases[i]
		want := describe(expected(c))
		if got != want {
			log.Fatalf("%s row %d pairs=%q text=%x: Lean=%s want=%s", label, i, c.pairs, c.text, got, want)
		}
		rawPairs := make([][]byte, len(c.pairs))
		for j, field := range c.pairs {
			rawPairs[j] = []byte(field)
		}
		must(fixture.Encode(struct {
			OldNew [][]byte
			Text   []byte
			Want   string
		}{rawPairs, []byte(c.text), want}))
	}
	fmt.Printf("genericReplacer.Replace %s: %d Lean results agree.\n", label, len(cases))
}

func describe(raw string) string {
	row := []byte("bytes:")
	for i := range len(raw) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(raw[i]), 10)
	}
	return string(row)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
