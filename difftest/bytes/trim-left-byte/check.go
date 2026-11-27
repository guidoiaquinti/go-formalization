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
	dir, err := os.MkdirTemp("", "golean-bytes-trim-left-byte-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.Setenv("GOLEAN_BYTES_TRIM_LEFT_BYTE_CASES", path); err != nil {
		return err
	}
	var edges []input
	for cut := range 256 {
		edges = append(edges, input{"", byte(cut), -1}, input{"", byte(cut), 0}, input{"", byte(cut), 5})
		for other := range 256 {
			text := string([]byte{byte(cut), byte(cut), byte(other), byte(cut)})
			edges = append(edges, input{text, byte(cut), len(text) + (cut+other)%8})
		}
	}

	batch := func(label string, cases []input) error {
		var modelInput, fixture bytes.Buffer
		expected := make([]string, len(cases))
		for i, c := range cases {
			row := strconv.AppendUint(nil, uint64(c.cut), 10)
			row = strconv.AppendInt(append(row, ' '), int64(max(0, c.capacity-len(c.text))), 10)
			for j := range len(c.text) {
				row = strconv.AppendUint(append(row, ' '), uint64(c.text[j]), 10)
			}
			fmt.Fprintln(&modelInput, string(row))

			same := strings.Repeat(string([]byte{c.cut}), len(c.text))
			prefix := sort.Search(len(c.text), func(i int) bool {
				return c.text[:i+1] != same[:i+1]
			})
			want := c.text[prefix:]
			capacity := c.capacity - prefix
			if want == "" {
				capacity = 0
			}
			fmt.Fprintf(&fixture, "%d\t%d\t%x\t%d\t%x\n", c.cut, c.capacity, c.text, capacity, want)
			if want == "" {
				expected[i] = "nil"
			} else {
				expected[i] = describe(want, capacity)
			}
		}
		for i, got := range check.LeanRows("godiff-bytes-trim-left-byte", len(cases), &modelInput) {
			if got != expected[i] {
				return fmt.Errorf("%s trimLeftByte(%x,%02x): Lean=%s want=%s", label, cases[i].text, cases[i].cut, got, expected[i])
			}
		}
		if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
			return err
		}
		if err := check.TestInPackage("bytes", "TestGoLeanTrimLeftByte",
			"difftest/bytes/trim-left-byte/testdata/bridge_test.go",
			"difftest/bytes/trim-left-byte/testdata/trim_left_byte_test.go"); err != nil {
			return err
		}
		fmt.Printf("bytes.trimLeftByte %s: %d inputs agree.\n", label, len(cases))
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
			prefix := strings.Repeat(string([]byte{cut}), r.Intn(32))
			text := prefix + string(raw)
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
