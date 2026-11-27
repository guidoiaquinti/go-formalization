package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type match struct {
	value string
	n     int
	found bool
}
type witness struct {
	pairs  []string
	text   string
	ignore bool
	want   match
}

func wire(m match) string {
	tag := "miss"
	if m.found {
		tag = "hit"
	}
	values := make([]string, len(m.value))
	for i := range len(m.value) {
		values[i] = strconv.Itoa(int(m.value[i]))
	}
	return fmt.Sprintf("%s:%d:%s", tag, m.n, strings.Join(values, ","))
}

func main() {
	cases := []witness{
		{nil, "", false, match{}},
		{[]string{"", "R", "ab", "two", "a", "one"}, "abZ", false, match{"R", 0, true}},
		{[]string{"", "R", "ab", "two", "a", "one"}, "abZ", true, match{"two", 2, true}},
		{[]string{"a", "one", "ab", "two"}, "abZ", false, match{"one", 1, true}},
		{[]string{"ab", "two", "a", "one"}, "abZ", false, match{"two", 2, true}},
		{[]string{"ab", "first", "ab", "second", "a", "short"}, "ab", true, match{"first", 2, true}},
		{[]string{"ab", "", "a", "short"}, "ab", false, match{"", 2, true}},
		{[]string{"abcdef", "long", "ab", "AB", "axy", "AXY"}, "abc", false, match{"AB", 2, true}},
		{[]string{"abcdef", "long", "axy", "AXY"}, "abc", false, match{}},
		{[]string{"abc", "A", "x\xff", "B"}, "qabc", false, match{}},
		{[]string{"\xff\x00", "\x80\xff", "\xff", "F"}, "\xff\x00!", false, match{"\x80\xff", 2, true}},
		{[]string{"", "root", "abc", "a"}, "", true, match{}},
	}
	var full []string
	for i := range 256 {
		full = append(full, string([]byte{byte(i), 0xff}), string([]byte{byte(i)}))
	}
	cases = append(cases, witness{full, "\xff\xff!", true, match{"\xff", 2, true}}, witness{full, "\xff\x80", true, match{}})
	file, err := os.CreateTemp("", "golean-review-generic-lookup-")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	var model bytes.Buffer
	raw := func(s string) {
		for i := range len(s) {
			fmt.Fprintf(&model, " %d", s[i])
		}
	}
	for _, c := range cases {
		flag := 0
		if c.ignore {
			flag = 1
		}
		fmt.Fprintf(&model, "%d %d", flag, len(c.pairs)/2)
		oldnew := make([][]byte, len(c.pairs))
		for i := 0; i < len(c.pairs); i += 2 {
			key, value := c.pairs[i], c.pairs[i+1]
			fmt.Fprintf(&model, " %d %d", len(key), len(value))
			raw(key)
			raw(value)
			oldnew[i], oldnew[i+1] = []byte(key), []byte(value)
		}
		raw(c.text)
		fmt.Fprintln(&model)
		must(enc.Encode(struct {
			OldNew     [][]byte
			Text       []byte
			IgnoreRoot bool
			Want       string
		}{oldnew, []byte(c.text), c.ignore, wire(c.want)}))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_GENERIC_REPLACER_LOOKUP_CASES", file.Name()))
	const fixtures = "difftest/strings/generic-replacer-lookup/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanGenericReplacerLookup", fixtures+"bridge_test.go", fixtures+"lookup_test.go"))
	for i, got := range check.LeanRows("godiff-strings-generic-replacer-lookup", len(cases), &model) {
		if got != wire(cases[i].want) {
			log.Fatalf("Lean row %d: got %q want %q", i, got, wire(cases[i].want))
		}
	}
	longestBad, ignoredRootBad := 0, 0
	for _, c := range cases {
		longest, neverSkipRoot := match{}, match{}
		for i := 0; i < len(c.pairs); i += 2 {
			key, value := c.pairs[i], c.pairs[i+1]
			if strings.HasPrefix(c.text, key) {
				if !neverSkipRoot.found {
					neverSkipRoot = match{value, len(key), true}
				}
				if (key != "" || !c.ignore) && (!longest.found || len(key) > longest.n) {
					longest = match{value, len(key), true}
				}
			}
		}
		if longest != c.want {
			longestBad++
		}
		if neverSkipRoot != c.want {
			ignoredRootBad++
		}
	}
	if longestBad == 0 || ignoredRootBad == 0 {
		log.Fatal("plausible mutation survived")
	}
	fmt.Printf("genericReplacer.lookup passed %d independent native and Lean fixtures; longest-first failed %d, root-always-eligible failed %d\n", len(cases), longestBad, ignoredRootBad)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
