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

type fixture struct {
	pairs      []string
	text, want string
}

func wire(s string) string {
	raw := make([]string, len(s))
	for i := range len(s) {
		raw[i] = strconv.Itoa(int(s[i]))
	}
	return "bytes:" + strings.Join(raw, ",")
}

func main() {
	cases := []fixture{
		{nil, "", ""},
		{nil, "\xff\x00", "\xff\x00"},
		{[]string{"", "X"}, "", "X"},
		{[]string{"", "X"}, "é", "X\xc3X\xa9X"},
		{[]string{"", "X"}, "\xe2\x82", "X\xe2X\x82X"},
		{[]string{"", "X", "\xff", "Y"}, "\xff", "XYX"},
		{[]string{"\xff", "Y", "", "X"}, "\xff", "YX"},
		{[]string{"a", "short", "ab", "long"}, "ab", "shortb"},
		{[]string{"ab", "long", "a", "short"}, "ab", "long"},
		{[]string{"aa", "R"}, "aaa", "Ra"},
		{[]string{"a", "", "a", "second"}, "a", ""},
		{[]string{"a", "b", "b", "c"}, "ab", "bc"},
		{[]string{"\xff\x00", "\x80\xff"}, "x\xff\x00", "x\x80\xff"},
	}
	f, err := os.CreateTemp("", "golean-review-generic-replace-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var model bytes.Buffer
	add := func(s string) {
		for i := range len(s) {
			fmt.Fprintf(&model, " %d", s[i])
		}
	}
	for _, c := range cases {
		fmt.Fprint(&model, len(c.pairs)/2)
		pairs := make([][]byte, len(c.pairs))
		for i := 0; i < len(c.pairs); i += 2 {
			key, value := c.pairs[i], c.pairs[i+1]
			fmt.Fprintf(&model, " %d %d", len(key), len(value))
			add(key)
			add(value)
			pairs[i], pairs[i+1] = []byte(key), []byte(value)
		}
		add(c.text)
		fmt.Fprintln(&model)
		must(enc.Encode(struct {
			OldNew [][]byte
			Text   []byte
			Want   string
		}{pairs, []byte(c.text), wire(c.want)}))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_GENERIC_REPLACER_REPLACE_CASES", f.Name()))
	const data = "difftest/strings/generic-replacer-replace/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanGenericReplacerReplace", data+"bridge_test.go", data+"replace_test.go"))
	for i, got := range check.LeanRows("godiff-strings-generic-replacer-replace", len(cases), &model) {
		if got != wire(cases[i].want) {
			log.Fatalf("Lean row %d: got %q want %q", i, got, wire(cases[i].want))
		}
	}
	runeBoundary := cases[3]
	if strings.ReplaceAll(runeBoundary.text, "", "X") == runeBoundary.want {
		log.Fatal("rune-boundary empty replacement survived")
	}
	ordered := cases[7]
	if strings.ReplaceAll(ordered.text, "ab", "long") == ordered.want {
		log.Fatal("longest-first survived")
	}
	literal := cases[11]
	if strings.ReplaceAll(strings.ReplaceAll(literal.text, "a", "b"), "b", "c") == literal.want {
		log.Fatal("cascading replacement survived")
	}
	fmt.Printf("genericReplacer.Replace passed %d independent native and Lean results; rune-boundary, longest-first and cascading mutations detected\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
