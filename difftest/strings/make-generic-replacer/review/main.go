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
	pairs          []string
	alphabet, trie string
}

func mapping(c fixture) string {
	values := make([]string, 256)
	for i := range values {
		rank := strings.IndexByte(c.alphabet, byte(i))
		if rank < 0 {
			rank = len(c.alphabet)
		}
		values[i] = strconv.Itoa(rank)
	}
	return strconv.Itoa(len(c.alphabet)) + ":" + strings.Join(values, ",") + "|" + c.trie
}

func main() {
	cases := []fixture{
		{nil, "", "T0:{}"},
		{[]string{"", "first", "", "later"}, "", "T4:102,105,114,115,116{}"},
		{[]string{"abc", "~"}, "abc", "T0:{97=E0:[98,99|L2:126]}"},
		{[]string{"abc", "Z", "a", "X", "ab", "Y"}, "abc", "T0:{97=E4:88[98|E2:89[99|L6:90]]}"},
		{[]string{"ab", "X", "ac", "Y", "z", "Q"}, "abcz", "T0:{97=T0:{98=L6:88|99=L4:89}|122=L2:81}"},
		{[]string{"a", "", "a", "A"}, "a", "T0:{97=L4:}"},
		{[]string{"\xff\x00", "\x01\xfe"}, "\x00\xff", "T0:{255=E0:[0|L2:1,254]}"},
	}
	alphabet := make([]byte, 256)
	var full []string
	children := make([]string, 256)
	for i := range 256 {
		alphabet[i] = byte(i)
		full = append(full, string([]byte{byte(255 - i)}), "")
		children[i] = fmt.Sprintf("%d=L%d:", i, 2*(i+1))
	}
	cases = append(cases, fixture{full, string(alphabet), "T0:{" + strings.Join(children, "|") + "}"})
	f, err := os.CreateTemp("", "golean-review-make-generic-")
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
		fmt.Fprintln(&model)
		must(enc.Encode(struct {
			Pairs [][]byte
			Want  string
		}{pairs, mapping(c)}))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_MAKE_GENERIC_REPLACER_CASES", f.Name()))
	const data = "difftest/strings/make-generic-replacer/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanMakeGenericReplacer", data+"bridge_test.go", data+"constructor_test.go"))
	for i, got := range check.LeanRows("godiff-strings-make-generic-replacer", len(cases), &model) {
		if got != mapping(cases[i]) {
			log.Fatalf("Lean snapshot differed at row %d", i)
		}
	}
	wrongOrder := cases[6]
	wrongOrder.alphabet = "\xff\x00"
	if mapping(wrongOrder) == mapping(cases[6]) {
		log.Fatal("insertion-order mapping survived")
	}
	wrongDuplicate := cases[5]
	wrongDuplicate.trie = "T0:{97=L2:65}"
	if mapping(wrongDuplicate) == mapping(cases[5]) {
		log.Fatal("last duplicate won")
	}
	wrongCompressed := cases[2]
	wrongCompressed.trie = "T0:{97=E0:[98|E0:[99|L2:126]]}"
	if mapping(wrongCompressed) == mapping(cases[2]) {
		log.Fatal("uncompressed edges survived")
	}
	fmt.Printf("makeGenericReplacer passed %d independent native and Lean snapshots; mapping order, duplicate precedence and compression mutations detected\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
