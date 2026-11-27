package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text string
	cut  byte
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-trim-right-byte-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_STRINGS_TRIM_RIGHT_BYTE_CASES", path); err != nil {
		return err
	}
	var edges []input
	for cut := range 256 {
		edges = append(edges, input{"", byte(cut)}, input{string([]byte{byte(cut), byte(cut) ^ 1}), byte(cut)})
		for other := range 256 {
			edges = append(edges, input{string([]byte{byte(cut), byte(other), byte(cut), byte(cut)}), byte(cut)})
		}
	}
	batch := func(label string, cases []input) error {
		var modelInput, fixture bytes.Buffer
		expected := make([]string, len(cases))
		for i, c := range cases {
			row := strconv.AppendUint(nil, uint64(c.cut), 10)
			for j := range len(c.text) {
				row = strconv.AppendUint(append(row, ' '), uint64(c.text[j]), 10)
			}
			fmt.Fprintln(&modelInput, string(row))

			same := strings.Repeat(string([]byte{c.cut}), len(c.text))
			suffix := sort.Search(len(c.text), func(i int) bool { return c.text[len(c.text)-(i+1):] != same[:i+1] })
			want := c.text[:len(c.text)-suffix]
			fmt.Fprintf(&fixture, "%d\t%x\t%x\n", c.cut, c.text, want)
			expected[i] = describe(want)
		}
		for i, got := range check.LeanRows("godiff-strings-trim-right-byte", len(cases), &modelInput) {
			if got != expected[i] {
				return fmt.Errorf("%s trimRightByte(%x,%02x): Lean=%s want=%s", label, cases[i].text, cases[i].cut, got, expected[i])
			}
		}
		if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
			return err
		}
		if err := check.TestInPackage("strings", "TestGoLeanTrimRightByte",
			"difftest/strings/trim-right-byte/testdata/bridge_test.go",
			"difftest/strings/trim-right-byte/testdata/trim_right_byte_test.go"); err != nil {
			return err
		}
		fmt.Printf("strings.trimRightByte %s: %d inputs agree.\n", label, len(cases))
		return nil
	}
	if err := batch("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			cut := byte(r.Intn(256))
			raw := make([]byte, r.Intn(65))
			r.Read(raw)
			suffix := strings.Repeat(string([]byte{cut}), r.Intn(32))
			cases[i] = input{string(raw) + suffix, cut}
		}
		if err := batch(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	return nil
}

func describe(text string) string {
	row := []byte("bytes:")
	for i := range len(text) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
