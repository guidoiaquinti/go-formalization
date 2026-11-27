package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	pairs        []string
	text         string
	stringWriter bool
	script       [][2]int
	want         string
}

func main() {
	short := [][2]int{{0, 0}, {1, 0}}
	cases := []fixture{
		{nil, "", false, nil, "0:0:"},
		{nil, "raw\xff", false, nil, "4:0:B,114,97,119,255"},
		{[]string{"", "R", "a", "K"}, "a", true, nil, "3:0:S;S,82;S;S,75;S;S,82"},
		{[]string{"a", "K", "", "R"}, "a", true, nil, "2:0:S;S,75;S;S,82"},
		{[]string{"a", ""}, "aa", false, [][2]int{{-1, 0}, {-1, 2}}, "0:2:B;B"},
		{[]string{"b", "BB"}, "abz", true, [][2]int{{1, 1}}, "1:1:S,97"},
		{[]string{"b", "BB"}, "abz", true, [][2]int{{-1, 0}, {1, 2}}, "2:2:S,97;S,66,66"},
		{[]string{"b", "BB"}, "abz", true, short, "2:0:S,97;S,66,66;S,122"},
		{[]string{"b", "BB"}, "abz", true, [][2]int{{-1, 0}, {-1, 0}, {0, 1}}, "3:1:S,97;S,66,66;S,122"},
		{[]string{"", "R"}, "\xe2\x82", true, nil, "5:0:S;S,82;S,226;S,82;S,130;S,82"},
	}
	f, err := os.CreateTemp("", "golean-review-generic-write-")
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
		mode := 0
		if c.stringWriter {
			mode = 1
		}
		fmt.Fprintf(&model, "%d %d %d", mode, len(c.script), len(c.pairs)/2)
		for _, reply := range c.script {
			fmt.Fprintf(&model, " %d %d", reply[0], reply[1])
		}
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
			Pairs        [][]byte
			Text         []byte
			StringWriter bool
			Script       [][2]int
			Want         string
		}{pairs, []byte(c.text), c.stringWriter, c.script, c.want}))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_GENERIC_REPLACER_WRITE_STRING_CASES", f.Name()))
	const data = "difftest/strings/generic-replacer-write-string/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanGenericReplacerWriteString", data+"bridge_test.go", data+"write_string_test.go"))
	for i, got := range check.LeanRows("godiff-strings-generic-replacer-write-string", len(cases), &model) {
		if got != cases[i].want {
			log.Fatalf("Lean row %d: got %q want %q", i, got, cases[i].want)
		}
	}
	droppedEmpty, forcedBytes := 0, 0
	for _, c := range cases {
		parts := strings.SplitN(c.want, ":", 3)
		calls := strings.Split(parts[2], ";")
		kept := make([]string, 0, len(calls))
		for _, call := range calls {
			if call != "S" && call != "B" {
				kept = append(kept, call)
			}
		}
		if strings.Join(kept, ";") != parts[2] {
			droppedEmpty++
		}
		if strings.ReplaceAll(parts[2], "S", "B") != parts[2] {
			forcedBytes++
		}
	}
	if droppedEmpty == 0 || forcedBytes == 0 || cases[7].want == "0:0:S,97" {
		log.Fatal("plausible mutation survived")
	}
	fmt.Printf("genericReplacer.WriteString passed %d independent native and Lean traces; removed-empty failed %d, forced-byte dispatch failed %d\n", len(cases), droppedEmpty, forcedBytes)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
