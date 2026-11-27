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

type key struct {
	raw, value string
	priority   int
}

type input struct {
	pairs      []string
	key, value string
	priority   int
}

func main() {
	file, err := os.CreateTemp("", "golean-trie-node-add-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	fixture := json.NewEncoder(buffer)
	var edges []input
	for _, table := range boundaries() {
		for _, match := range []string{"", "a", "ab", "abc", "ac", "abcd", "\x00", "\xff\x00", "é"} {
			edges = append(edges, input{table, match, "\xff\x00", 1}, input{table, match, "", 99})
		}
	}
	compare("boundaries", edges, fixture)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5000)
			for i := range cases {
				pairs := sample(rng, i)
				pattern := sample(rng, i)
				added := ""
				if len(pattern) > 0 {
					added = pattern[0]
				}
				if len(pairs) > 0 {
					base := pairs[2*rng.Intn(len(pairs)/2)]
					switch i % 4 {
					case 0:
						added = base
					case 1:
						added = base[:rng.Intn(len(base)+1)]
					case 2:
						added = base + added
					}
				}
				cases[i] = input{pairs, added, string([]byte{byte(rng.Intn(256)), 0}), 1 + rng.Intn(500)}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixture)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_TRIE_NODE_ADD_CASES", file.Name()))
	const data = "difftest/strings/trie-node-add/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanTrieNodeAdd", data+"bridge_test.go", data+"add_test.go"))
}

func boundaries() [][]string {
	cases := [][]string{
		nil, {"", ""}, {"", "", "", "unused"}, {"abc", "one"}, {"ab", "one", "abc", "two"},
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
	cases = append(cases, all)
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

// Build a canonical snapshot directly from the complete key set. Compress
// an edge by taking the longest prefix shared by all unfinished suffixes.
func expected(c input) string {
	pairs := c.pairs
	alphabet := make(map[byte]bool)
	seen := make(map[string]bool)
	var keys []key
	for i := 0; i < len(pairs); i += 2 {
		pattern := pairs[i]
		if seen[pattern] {
			continue
		}
		seen[pattern] = true
		for j := range len(pattern) {
			alphabet[pattern[j]] = true
		}
		keys = append(keys, key{pattern, pairs[i+1], len(pairs) - i})
	}
	if !seen[c.key] {
		keys = append(keys, key{c.key, c.value, c.priority})
		for i := range len(c.key) {
			alphabet[c.key[i]] = true
		}
	}
	var sorted []byte
	for b := range alphabet {
		sorted = append(sorted, b)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	var snapshot func([]key, bool) string
	snapshot = func(entries []key, root bool) string {
		terminal := "0:"
		var proper []key
		for _, entry := range entries {
			if entry.raw == "" {
				terminal = strconv.Itoa(entry.priority) + ":" + wire(entry.value)
			} else {
				proper = append(proper, entry)
			}
		}
		if len(proper) == 0 && !root {
			return "L" + terminal
		}
		if !root {
			prefix := proper[0].raw
			for _, entry := range proper[1:] {
				for !strings.HasPrefix(entry.raw, prefix) {
					prefix = prefix[:len(prefix)-1]
				}
			}
			if prefix != "" {
				next := make([]key, len(proper))
				for i, entry := range proper {
					entry.raw = entry.raw[len(prefix):]
					next[i] = entry
				}
				return "E" + terminal + "[" + wire(prefix) + "|" + snapshot(next, false) + "]"
			}
		}
		groups := make(map[byte][]key)
		for _, entry := range proper {
			first := entry.raw[0]
			entry.raw = entry.raw[1:]
			groups[first] = append(groups[first], entry)
		}
		var children []string
		for _, b := range sorted {
			if group := groups[b]; len(group) > 0 {
				children = append(children, strconv.Itoa(int(b))+"="+snapshot(group, false))
			}
		}
		return "T" + terminal + "{" + strings.Join(children, "|") + "}"
	}
	return snapshot(keys, true)
}

func compare(label string, cases []input, fixture *json.Encoder) {
	var input bytes.Buffer
	for _, c := range cases {
		pairs := c.pairs
		fmt.Fprintf(&input, "%d %d %d %d", len(pairs)/2, c.priority, len(c.key), len(c.value))
		for _, raw := range []string{c.key, c.value} {
			for i := range len(raw) {
				fmt.Fprintf(&input, " %d", raw[i])
			}
		}
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
	for i, got := range check.LeanRows("godiff-strings-trie-node-add", len(cases), &input) {
		c := cases[i]
		pairs := c.pairs
		want := expected(c)
		if got != want {
			log.Fatalf("%s row %d pairs=%q key=%q value=%q priority=%d: Lean=%s want=%s", label, i, pairs, c.key, c.value, c.priority, got, want)
		}
		rawPairs := make([][]byte, len(pairs))
		for j, raw := range pairs {
			rawPairs[j] = []byte(raw)
		}
		must(fixture.Encode(struct {
			Pairs      [][]byte
			Key, Value []byte
			Priority   int
			Want       string
		}{rawPairs, []byte(c.key), []byte(c.value), c.priority, want}))
	}
	fmt.Printf("trieNode.add %s: %d Lean snapshots agree.\n", label, len(cases))
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
