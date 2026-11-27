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

func table(mapping func(byte) string) string {
	out := make([]string, 256)
	for i := range out {
		out[i] = mapping(byte(i))
	}
	return strings.Join(out, ",")
}
func identity(b byte) string { return strconv.Itoa(int(b)) }
func generic(alphabet string) string {
	return fmt.Sprintf("generic:%d:", len(alphabet)) + table(func(b byte) string {
		i := strings.IndexByte(alphabet, b)
		if i < 0 {
			i = len(alphabet)
		}
		return strconv.Itoa(i)
	})
}
func byteStrings(keys string, entries map[byte]string) string {
	out := make([]string, 256)
	for i := range out {
		out[i] = "-"
		if field, ok := entries[byte(i)]; ok {
			out[i] = "+" + field
		}
	}
	return "byteStrings:" + keys + "|" + strings.Join(out, "|")
}
func main() {
	cases := []struct {
		pairs []string
		want  string
	}{
		{nil, "bytes:" + table(identity)},
		{[]string{"é", "\xff"}, "single:195,169|255"},
		{[]string{"\xff\x00", ""}, "single:255,0|"},
		{[]string{"", "hello"}, generic("")},
		{[]string{"ab", "x", "a", "y"}, generic("ab")},
		{[]string{"\xff", "A", "\xff", "B"}, "bytes:" + table(func(b byte) string {
			if b == 255 {
				return "65"
			}
			return identity(b)
		})},
		{[]string{"a", "", "b", "wide", "a", "q", "\xff", "t", "b", "u"},
			byteStrings("98,255,97", map[byte]string{'a': "", 'b': "119,105,100,101", 255: "116"})},
		{[]string{"a", "b", "a", ""}, byteStrings("97", map[byte]string{'a': "98"})},
	}
	f, err := os.CreateTemp("", "golean-review-replacer-build-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&model, len(c.pairs)/2)
		pairs := make([][]byte, len(c.pairs))
		for i := 0; i < len(c.pairs); i += 2 {
			fmt.Fprintf(&model, " %d %d", len(c.pairs[i]), len(c.pairs[i+1]))
			for j := i; j < i+2; j++ {
				pairs[j] = []byte(c.pairs[j])
				for _, b := range pairs[j] {
					fmt.Fprintf(&model, " %d", b)
				}
			}
		}
		fmt.Fprintln(&model)
		must(enc.Encode(struct {
			Pairs [][]byte
			Want  string
		}{pairs, c.want}))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_REPLACER_BUILD_CASES", f.Name()))
	const data = "difftest/strings/replacer-build/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanReplacerBuild", data+"bridge_test.go", data+"build_test.go"))
	for i, got := range check.LeanRows("godiff-strings-replacer-build", len(cases), &model) {
		if got != cases[i].want {
			log.Fatalf("Lean snapshot differed at row %d", i)
		}
	}
	if generic("ab") == cases[1].want ||
		byteStrings("97,98,255", map[byte]string{'a': "", 'b': "119,105,100,101", 255: "116"}) == cases[6].want ||
		byteStrings("97", map[byte]string{'a': ""}) == cases[7].want {
		log.Fatal("implementation or duplicate mutation survived")
	}
	fmt.Printf("Replacer.build passed %d independent native and Lean snapshots\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
