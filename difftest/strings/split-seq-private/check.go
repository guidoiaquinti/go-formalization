package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s, sep              string
	save, first, second int
}

func main() {
	file, err := os.CreateTemp("", "golean-split-seq-private-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	fixture := json.NewEncoder(buffer)
	var cases []input
	add := func(s, sep string) {
		for _, save := range []int{math.MinInt, -3, -1, 0, 1, len(sep), len(sep) + 1, len(s), len(s) + 1, math.MaxInt} {
			for _, budgets := range [][2]int{{0, 0}, {1, 1}, {2, 1}, {1, 2}, {len(s) + 2, 1}} {
				cases = append(cases, input{s, sep, save, budgets[0], budgets[1]})
			}
		}
	}
	for _, s := range []string{"", "a", "aa", "a--b--", "---x", "é🚀", "\xff\xed\xa0\x80", "a\x00a"} {
		for _, sep := range []string{"", "a", "--", "-", "é", "\xa9", "\x00", "missing"} {
			add(s, sep)
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		add(raw+"a"+raw, raw)
	}
	compare("boundaries", cases, fixture)
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
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_SPLIT_SEQ_PRIVATE_CASES", file.Name()))
	const data = "difftest/strings/split-seq-private/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanSplitSeqPrivate", data+"bridge_test.go", data+"split_seq_private_test.go"))
}

func sample(rng *rand.Rand, index int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	start := rng.Intn(len(raw) + 1)
	end := start + rng.Intn(min(6, len(raw)-start)+1)
	sep := string(raw[start:end])
	s := string(raw)
	if rng.Intn(3) == 0 {
		sep = strings.Repeat("a", rng.Intn(4))
	}
	if index%3 == 0 {
		s = sep + s + sep + sep
	}
	save := rng.Intn(len(s)+7) - 3
	if index%4 == 0 {
		save = []int{math.MinInt, -1, 0, len(sep), math.MaxInt}[rng.Intn(5)]
	}
	return input{s, sep, save, rng.Intn(len(s) + 2), rng.Intn(len(s) + 2)}
}

// The Split fragments determine original match offsets. In each resumed
// batch the current fragment stays unconsumed when returning false or
// when the signed endpoint would fail slicing.
func expected(c input) string {
	fragments := strings.Split(c.s, c.sep)
	cursor, pos := 0, 0
	run := func(budget int) string {
		var parts []string
		for len(parts) < budget && cursor < len(fragments) {
			part := fragments[cursor]
			final := c.sep != "" && cursor == len(fragments)-1
			if c.sep != "" && !final {
				remaining := len(c.s) - pos
				if c.save < -len(part) || c.save > remaining-len(part) {
					return "p:" + strings.Join(parts, ";")
				}
				part = c.s[pos : pos+len(part)+c.save]
			}
			parts = append(parts, wire(part))
			if len(parts) == budget || final {
				break
			}
			pos += len(fragments[cursor]) + len(c.sep)
			cursor++
		}
		return "o:" + strings.Join(parts, ";")
	}
	return strings.Join([]string{run(c.first), run(c.second), run(len(c.s) + 2), run(1)}, "|")
}

func wire(part string) string {
	raw := []byte{'f'}
	for i := range len(part) {
		if i != 0 {
			raw = append(raw, ',')
		}
		raw = strconv.AppendUint(raw, uint64(part[i]), 10)
	}
	return string(raw)
}

func compare(label string, cases []input, fixture *json.Encoder) {
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&rows, "%d %d %d %d", len(c.sep), c.save, c.first, c.second)
		for _, raw := range []string{c.sep, c.s} {
			for i := range len(raw) {
				fmt.Fprintf(&rows, " %d", raw[i])
			}
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-strings-split-seq-private", len(cases), &rows) {
		c := cases[i]
		want := expected(c)
		if got != want {
			log.Fatalf("%s row %d input=%x sep=%x save=%d: Lean=%s want=%s", label, i, c.s, c.sep, c.save, got, want)
		}
		must(fixture.Encode(struct {
			Input, Sep          []byte
			Save, First, Second int
			Want                string
		}{[]byte(c.s), []byte(c.sep), c.save, c.first, c.second, want}))
	}
	fmt.Printf("strings.splitSeq %s: %d Lean trace inputs agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
