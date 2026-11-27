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
	oldnew     []string
	text       string
	ignoreRoot bool
}

func main() {
	file, err := os.CreateTemp("", "golean-generic-replacer-lookup-")
	must(err)
	defer os.Remove(file.Name())
	writer := bufio.NewWriter(file)
	rows := json.NewEncoder(writer)
	compare("boundaries", boundaries(), rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cases[i] = sample(rng, i)
		}
		compare(fmt.Sprintf("seed=%d", seed), cases, rows)
	}
	must(writer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_GENERIC_REPLACER_LOOKUP_CASES", file.Name()))
	const data = "difftest/strings/generic-replacer-lookup/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanGenericReplacerLookup", data+"bridge_test.go", data+"lookup_test.go"))
}

func boundaries() []input {
	tables := [][]string{
		nil, {"", ""}, {"", "first", "", "second"}, {"", "root", "a", "A"},
		{"ab", "long", "a", "short", "", "root"}, {"a", "short", "ab", "long", "", "root"},
		{"a", "first", "a", "second", "abc", "long"}, {"abcdef", "AB", "abcxyz", "XYZ", "a", "A"},
		{"\x00\xff", "\xff\x00", "\x00", "", "ǆ🚀", "\x80", "\xe2\x82", "\xff"},
		{"abcdef", "AB", "abz", "ABZ", "a\x00", "A0", "abd", "ABD"},
	}
	allBytes := make([]string, 0, 512)
	for b := range 256 {
		allBytes = append(allBytes, string([]byte{byte(b), 0xff}), string([]byte{byte(255 - b)}))
	}
	var cases []input
	for _, table := range tables {
		for _, text := range []string{"", "a", "ab", "abc", "abcd", "abcdef", "abcxyz", "abcxy", "abd", "abzZ",
			"\x00\xff\x80", "ǆ🚀z", "\xe2\x82Z", "\xff", "a\x00z", "zabc"} {
			cases = append(cases, input{table, text, false}, input{table, text, true})
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		cases = append(cases, input{allBytes, raw + "\xff!", false}, input{allBytes, raw + "\x00", true})
	}
	return cases
}

func sample(rng *rand.Rand, index int) input {
	const duplicate, truncate, prolong, prepend = 0, 1, 2, 3
	alphabet := [...]int{2, 8, 256}[index%3]
	raw := func(max int) string {
		text := make([]byte, rng.Intn(max+1))
		for i := range text {
			text[i] = byte(rng.Intn(alphabet))
		}
		return string(text)
	}
	var oldnew []string
	for range rng.Intn(16) {
		key := raw(22)
		if len(oldnew) > 0 {
			base := oldnew[2*rng.Intn(len(oldnew)/2)]
			switch rng.Intn(5) {
			case duplicate:
				key = base
			case truncate:
				key = base[:rng.Intn(len(base)+1)]
			case prolong:
				key = base + key
			case prepend:
				key = key + base
			}
		}
		oldnew = append(oldnew, key, raw(11))
	}
	text := raw(35)
	if index%2 == 0 && len(oldnew) != 0 {
		key := oldnew[2*rng.Intn(len(oldnew)/2)]
		if index%4 != 0 {
			key = key[:rng.Intn(len(key)+1)]
		}
		text = key + text
	}
	return input{oldnew, text, rng.Intn(2) != 0}
}

// Constructor priorities strictly decrease in old/new order. The first
// matching pair therefore determines the returned value and byte length.
func expected(c input) string {
	for i := 0; i < len(c.oldnew); i += 2 {
		key, value := c.oldnew[i], c.oldnew[i+1]
		if !(key == "" && c.ignoreRoot) && strings.HasPrefix(c.text, key) {
			fields := make([]string, len(value))
			for j := range len(value) {
				fields[j] = strconv.Itoa(int(value[j]))
			}
			return fmt.Sprintf("hit:%d:%s", len(key), strings.Join(fields, ","))
		}
	}
	return "miss:0:"
}

func compare(label string, cases []input, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		flag := 0
		if c.ignoreRoot {
			flag = 1
		}
		fmt.Fprintf(&wire, "%d %d", flag, len(c.oldnew)/2)
		for i := 0; i < len(c.oldnew); i += 2 {
			key, value := c.oldnew[i], c.oldnew[i+1]
			fmt.Fprintf(&wire, " %d %d", len(key), len(value))
			writeBytes(&wire, key)
			writeBytes(&wire, value)
		}
		writeBytes(&wire, c.text)
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-strings-generic-replacer-lookup", len(cases), &wire) {
		c := cases[i]
		want := expected(c)
		if got != want {
			log.Fatalf("%s row %d pairs=%q input=%x ignoreRoot=%t: Lean=%s want=%s", label, i, c.oldnew, c.text, c.ignoreRoot, got, want)
		}
		oldnew := make([][]byte, len(c.oldnew))
		for j, raw := range c.oldnew {
			oldnew[j] = []byte(raw)
		}
		must(rows.Encode(struct {
			OldNew     [][]byte
			Text       []byte
			IgnoreRoot bool
			Want       string
		}{oldnew, []byte(c.text), c.ignoreRoot, want}))
	}
	fmt.Printf("genericReplacer.lookup %s: %d Lean rows agree.\n", label, len(cases))
}

func writeBytes(wire *bytes.Buffer, value string) {
	for i := range len(value) {
		fmt.Fprintf(wire, " %d", value[i])
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
