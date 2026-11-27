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
	text     string
	cut      byte
	capacity int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-bytes-trim-right-byte-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BYTES_TRIM_RIGHT_BYTE_CASES", path); err != nil {
		return err
	}
	var edges []input
	for cut := range 256 {
		edges = append(edges, input{"", byte(cut), -1}, input{"", byte(cut), 0}, input{"", byte(cut), 5})
		for other := range 256 {
			text := string([]byte{byte(cut), byte(other), byte(cut), byte(cut)})
			edges = append(edges, input{text, byte(cut), len(text) + (cut+other)%8})
		}
	}

	batch := func(label string, cases []input) error {
		var modelInput, fixture bytes.Buffer
		expected := make([]string, len(cases))
		for i, c := range cases {
			row := strconv.AppendUint(nil, uint64(c.cut), 10)
			row = strconv.AppendInt(append(row, ' '), int64(max(0, c.capacity)), 10)
			if c.capacity < 0 {
				row = append(row, " 1"...)
			} else {
				row = append(row, " 0"...)
			}
			for j := range len(c.text) {
				row = strconv.AppendUint(append(row, ' '), uint64(c.text[j]), 10)
			}
			fmt.Fprintln(&modelInput, string(row))

			same := strings.Repeat(string([]byte{c.cut}), len(c.text))
			suffix := sort.Search(len(c.text), func(i int) bool {
				offset := len(c.text) - i - 1
				return c.text[offset:] != same[offset:]
			})
			want := c.text[:len(c.text)-suffix]
			isNil := c.capacity < 0
			fmt.Fprintf(&fixture, "%d\t%d\t%x\t%t\t%x\n", c.cut, c.capacity, c.text, isNil, want)
			if isNil {
				expected[i] = "nil"
			} else {
				expected[i] = describe(want, c.capacity)
			}
		}
		for i, got := range check.LeanRows("godiff-bytes-trim-right-byte", len(cases), &modelInput) {
			if got != expected[i] {
				return fmt.Errorf("%s trimRightByte(%x,%02x): Lean=%s want=%s", label, cases[i].text, cases[i].cut, got, expected[i])
			}
		}
		if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
			return err
		}
		if err := check.TestInPackage("bytes", "TestGoLeanTrimRightByte",
			"difftest/bytes/trim-right-byte/testdata/bridge_test.go",
			"difftest/bytes/trim-right-byte/testdata/trim_right_byte_test.go"); err != nil {
			return err
		}
		fmt.Printf("bytes.trimRightByte %s: %d inputs agree.\n", label, len(cases))
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
			text := string(raw) + suffix
			cases[i] = input{text, cut, len(text) + r.Intn(12)}
		}
		if err := batch(fmt.Sprintf("seed=%d", seed), cases); err != nil {
			return err
		}
	}
	return nil
}

func describe(text string, capacity int) string {
	row := strconv.AppendInt(nil, int64(capacity), 10)
	row = append(row, ':')
	for i := range len(text) {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(text[i]), 10)
	}
	return string(row)
}
