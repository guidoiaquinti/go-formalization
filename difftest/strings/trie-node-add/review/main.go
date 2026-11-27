package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		pairs      []string
		key, value string
		priority   int
		want       string
	}{
		{nil, "", "R", 7, "T7:82{}"},
		{[]string{"", "X"}, "", "Y", 99, "T2:88{}"},
		{[]string{"abc", "X"}, "abc", "Y", 99, "T0:{97=E0:[98,99|L2:88]}"},
		{[]string{"abc", "X"}, "ab", "Y", 7, "T0:{97=E0:[98|E7:89[99|L2:88]]}"},
		{[]string{"abc", "X"}, "abd", "Y", 7, "T0:{97=E0:[98|T0:{99=L2:88|100=L7:89}]}"},
		{[]string{"abc", "X"}, "ax", "Y", 7, "T0:{97=T0:{98=E0:[99|L2:88]|120=L7:89}}"},
		{[]string{"a", "X"}, "abc", "Y", 7, "T0:{97=E2:88[98,99|L7:89]}"},
		{[]string{"ab", "X", "ax", "Y"}, "ad", "", 7, "T0:{97=T0:{98=L4:88|100=L7:|120=L2:89}}"},
		{[]string{"\xff\x00", "X"}, "\xff\x80", "\x00", 7, "T0:{255=T0:{0=L2:88|128=L7:0}}"},
	}
	f, err := os.CreateTemp("", "golean-review-trie-add-")
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
		fmt.Fprintf(&model, "%d %d %d %d", len(c.pairs)/2, c.priority, len(c.key), len(c.value))
		add(c.key)
		add(c.value)
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
			Pairs      [][]byte
			Key, Value []byte
			Priority   int
			Want       string
		}{pairs, []byte(c.key), []byte(c.value), c.priority, c.want}))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_TRIE_NODE_ADD_CASES", f.Name()))
	const dir = "difftest/strings/trie-node-add/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanTrieNodeAdd", dir+"bridge_test.go", dir+"add_test.go"))
	for i, got := range check.LeanRows("godiff-strings-trie-node-add", len(cases), &model) {
		if got != cases[i].want {
			log.Fatalf("row %d: got %q want %q", i, got, cases[i].want)
		}
	}
	if cases[2].want == "T0:{97=E0:[98,99|L99:89]}" || cases[3].want == "T0:{97=E0:[98,99|L2:88]}" {
		log.Fatal("duplicate overwrite or swallowed terminal mutation survived")
	}
	fmt.Printf("trieNode.add passed %d independent native and Lean insertion snapshots\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
