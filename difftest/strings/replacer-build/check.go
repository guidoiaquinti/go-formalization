package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	file, err := os.CreateTemp("", "golean-replacer-build-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	fixture := json.NewEncoder(buffer)
	compare("boundaries", boundaries(), fixture)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([][]string, 5000)
			for i := range cases {
				cases[i] = sample(rng, i)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixture)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_REPLACER_BUILD_CASES", file.Name()))
	const data = "difftest/strings/replacer-build/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanReplacerBuild", data+"bridge_test.go", data+"build_test.go"))
}

func boundaries() [][]string {
	cases := [][]string{
		nil, {"a", "A"}, {"a", ""}, {"a", "A", "b", "B", "a", "xy"}, {"a", "", "b", "B", "a", "C"}, {"", ""}, {"", "", "", "unused"}, {"abc", "one"}, {"ab", "one", "abc", "two"},
		{"abc", "two", "ab", "one"}, {"abc", "one", "abd", "two"}, {"abc", "one", "adc", "two"},
		{"", "root", "a", "first", "abc", "deep", "a", "ignored", "abd", "next", "b", "B"},
		{"\xff\x00", "\x00", "\xff", "", "\xe2\x82", "🚀", "é", "\xff"},
	}
	all := make([]string, 0, 512)
	for b := range 256 {
		one, other := string([]byte{byte(b)}), string([]byte{byte(255 - b)})
		cases = append(cases, []string{one + "ab", "", one + "abc", other, one + "d", "s", one, one})
		all = append(all, one, other)
	}
	cases = append(cases, all, append(append([]string{}, all...), "\xff", "", "\x00", "long"))
	var alphabet [256]byte
	for i := range alphabet {
		alphabet[i] = byte(i)
	}
	cases = append(cases, []string{string(alphabet[:]), "whole"}, []string{"f" + string(alphabet[:]), "deep", "f", ""})
	return cases
}

func sample(rng *rand.Rand, index int) []string {
	alphabet := [...]int{2, 7, 256}[index%3]
	raw := func(max int) string {
		data := make([]byte, rng.Intn(max+1))
		for i := range data {
			data[i] = byte(rng.Intn(alphabet))
		}
		return string(data)
	}
	var pairs []string
	for range rng.Intn(10) {
		pattern := raw(15)
		if index%4 == 0 {
			pattern = string([]byte{byte(rng.Intn(alphabet))})
		}
		if index%4 == 0 {
			value := string([]byte{byte(rng.Intn(alphabet))})
			if index%8 == 0 && rng.Intn(3) == 0 {
				value = raw(13)
			}
			pairs = append(pairs, pattern, value)
			continue
		}
		if len(pairs) > 0 {
			old := pairs[2*rng.Intn(len(pairs)/2)]
			switch rng.Intn(4) {
			case 0:
				pattern = old[:rng.Intn(len(old)+1)]
			case 1:
				pattern = old + pattern
			}
		}
		pairs = append(pairs, pattern, raw(13))
	}
	return pairs
}

// Classify by byte width and derive each output entry by first occurrence.
// The toReplace sequence is obtained by sorting keys on their final occurrence.
func expected(pairs []string) string {
	if len(pairs) == 2 && len(pairs[0]) > 1 {
		return "single:" + wire(pairs[0]) + "|" + wire(pairs[1])
	}
	var used [256]bool
	var first [256]*string
	var last [256]int
	for byte := range last {
		last[byte] = -1
	}
	multiKey, multiValue := false, false
	for i := 0; i < len(pairs); i += 2 {
		key, value := pairs[i], pairs[i+1]
		multiKey = multiKey || len(key) != 1
		multiValue = multiValue || len(value) != 1
		for j := range len(key) {
			used[key[j]] = true
		}
		if len(key) == 1 {
			b := key[0]
			if first[b] == nil {
				first[b] = &pairs[i+1]
			}
			last[b] = i
		}
	}
	if multiKey {
		var size int
		for _, present := range used {
			if present {
				size++
			}
		}
		var mapping [256]byte
		rank := 0
		for b, present := range used {
			mapping[b] = byte(size)
			if present {
				mapping[b] = byte(rank)
				rank++
			}
		}
		return fmt.Sprintf("generic:%d:%s", size, wire(string(mapping[:])))
	}
	if !multiValue {
		var mapping [256]byte
		for b, value := range first {
			mapping[b] = byte(b)
			if value != nil {
				mapping[b] = (*value)[0]
			}
		}
		return "bytes:" + wire(string(mapping[:]))
	}
	var order []byte
	values := make([]string, 256)
	for b, value := range first {
		values[b] = "-"
		if value != nil {
			values[b] = "+" + wire(*value)
			order = append(order, byte(b))
		}
	}
	sort.Slice(order, func(i, j int) bool { return last[order[i]] > last[order[j]] })
	return "byteStrings:" + wire(string(order)) + "|" + strings.Join(values, "|")
}

func compare(label string, cases [][]string, fixture *json.Encoder) {
	var input bytes.Buffer
	for _, pairs := range cases {
		fmt.Fprint(&input, len(pairs)/2)
		for i := 0; i < len(pairs); i += 2 {
			fmt.Fprintf(&input, " %d %d", len(pairs[i]), len(pairs[i+1]))
			for _, raw := range pairs[i : i+2] {
				for j := range len(raw) {
					fmt.Fprintf(&input, " %d", raw[j])
				}
			}
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-strings-replacer-build", len(cases), &input) {
		pairs := cases[i]
		want := expected(pairs)
		if got != want {
			log.Fatalf("%s row %d pairs=%q: Lean=%s want=%s", label, i, pairs, got, want)
		}
		rawPairs := make([][]byte, len(pairs))
		for j, raw := range pairs {
			rawPairs[j] = []byte(raw)
		}
		must(fixture.Encode(struct {
			Pairs [][]byte
			Want  string
		}{rawPairs, want}))
	}
	fmt.Printf("Replacer.build %s: %d Lean snapshots agree.\n", label, len(cases))
}

func wire(raw string) string {
	var out []byte
	for i := range len(raw) {
		if i != 0 {
			out = append(out, ',')
		}
		out = strconv.AppendUint(out, uint64(raw[i]), 10)
	}
	return string(out)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
