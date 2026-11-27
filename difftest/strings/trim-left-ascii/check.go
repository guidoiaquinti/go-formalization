package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	text     []byte
	members  []byte
	nilTable bool
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-trim-left-ascii-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_STRINGS_TRIM_LEFT_ASCII_CASES", path); err != nil {
		return err
	}
	edges := []input{{}, {nilTable: true}, {text: []byte{0xff}, nilTable: true}}
	full := make([]byte, 256)
	for value := range 256 {
		full[value] = byte(value)
		for other := range 256 {
			edges = append(edges, input{members: []byte{byte(value)}, text: []byte{byte(value), byte(other), byte(value)}})
		}
	}
	edges = append(edges, input{text: full, members: full}, input{text: full})
	batch := func(label string, cases []input) error {
		var modelInput, fixture bytes.Buffer
		expected := make([]string, len(cases))
		for i, c := range cases {
			size, set := len(c.members), fmt.Sprintf("%x", c.members)
			if c.nilTable {
				size, set = -1, "-"
			}
			row := strconv.AppendInt(nil, int64(size), 10)
			for _, raw := range [][]byte{c.members, c.text} {
				for _, b := range raw {
					row = strconv.AppendUint(append(row, ' '), uint64(b), 10)
				}
			}
			fmt.Fprintln(&modelInput, string(row))

			if c.nilTable && len(c.text) != 0 {
				fmt.Fprintf(&fixture, "%s\t%x\tpanic\n", set, c.text)
				expected[i] = "panic"
				continue
			}
			start := slices.IndexFunc(c.text, func(b byte) bool { return !slices.Contains(c.members, b) })
			var want []byte
			if start >= 0 {
				want = c.text[start:]
			}
			fmt.Fprintf(&fixture, "%s\t%x\t%x\n", set, c.text, want)
			expected[i] = describe(want)
		}
		for i, got := range check.LeanRows("godiff-strings-trim-left-ascii", len(cases), &modelInput) {
			if got != expected[i] {
				return fmt.Errorf("%s trimLeftASCII(%x,%x): Lean=%s want=%s", label, cases[i].text, cases[i].members, got, expected[i])
			}
		}
		if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
			return err
		}
		const fixtures = "difftest/strings/trim-left-ascii/testdata/"
		if err := check.TestInPackage("strings", "TestGoLeanTrimLeftASCII", fixtures+"bridge_test.go", fixtures+"trim_left_ascii_test.go"); err != nil {
			return err
		}
		fmt.Printf("strings.trimLeftASCII %s: %d inputs agree.\n", label, len(cases))
		return nil
	}
	if err := batch("boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]input, 20_000)
		for i := range cases {
			members := make([]byte, r.Intn(257))
			r.Read(members)
			text := make([]byte, r.Intn(65))
			r.Read(text)
			if len(members) != 0 {
				for j, count := 0, r.Intn(len(text)+1); j < count; j++ {
					text[j] = members[r.Intn(len(members))]
				}
			}
			if i%101 == 0 {
				cases[i] = input{text: text, nilTable: true}
			} else {
				cases[i] = input{text: text, members: members}
			}
		}
		if err := batch(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	return nil
}

func describe(text []byte) string {
	row := []byte("bytes:")
	for i, b := range text {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(b), 10)
	}
	return string(row)
}
