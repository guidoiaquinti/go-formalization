package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	table *[256]byte
	text  string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-byte-replacer-replace-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := os.Setenv("GOLEAN_STRINGS_BYTE_REPLACER_REPLACE_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)
	identity := new([256]byte)
	for i := range identity {
		identity[i] = byte(i)
	}
	edges := []input{{nil, ""}, {nil, "\xff"}, {nil, "a🚀\x00"}, {identity, ""}, {identity, "a🚀\x00"}}
	for shift := range 256 {
		table := new([256]byte)
		distinct := new([256]byte)
		for i := range table {
			table[i] = byte(i + shift)
			distinct[i] = byte(shift)
		}
		// Across the shifted tables, each possible input/output byte pair
		// occurs once. The constant tables check convergent mappings.
		edges = append(edges, input{table, string(identity[:])}, input{distinct, string(identity[:])})
	}
	for _, value := range []string{"", "\xffa\x80", "�🚀", "aaaa", strings.Repeat("a", 256)} {
		chain := new([256]byte)
		*chain = *identity
		chain['a'], chain[0xff], chain[0x80] = 0xff, 0x80, 'a'
		edges = append(edges, input{chain, value})
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5_000)
			for i := range cases {
				text := make([]byte, rng.Intn(56))
				rng.Read(text)
				if i%23 == 0 {
					cases[i] = input{nil, string(text)}
					continue
				}
				table := new([256]byte)
				*table = *identity
				changes := []int{0, 1, 16, 256}[i%4]
				for range changes {
					table[byte(rng.Uint32())] = byte(rng.Uint32())
				}
				cases[i] = input{table, string(text)}
			}
			if err := compare(fixture, fmt.Sprintf("seed=%d block=%d", seed, block), cases); err != nil {
				return err
			}
		}
	}
	if err := fixture.Flush(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	const fixtures = "difftest/strings/byte-replacer-replace/testdata/"
	return check.TestInPackage("strings", "TestGoLeanByteReplacerReplace",
		fixtures+"bridge_test.go", fixtures+"byte_replacer_replace_test.go")
}

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		var table []byte
		row := []byte("0")
		oracle := ""
		if c.table != nil {
			row[0] = '1'
			table = c.table[:]
			for _, value := range table {
				row = strconv.AppendUint(append(row, ' '), uint64(value), 10)
			}
		}
		for j := range len(c.text) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.text[j]), 10)
		}
		want := []byte("bytes:")
		switch {
		case table != nil:
			result := make([]byte, len(c.text))
			for j := range result {
				result[j] = table[c.text[j]]
				if j != 0 {
					want = append(want, ',')
				}
				want = strconv.AppendUint(want, uint64(result[j]), 10)
			}
			oracle = hex.EncodeToString(result)
		case c.text != "":
			want = []byte("panic")
			oracle = "panic"
		}
		fmt.Fprintln(&modelInput, string(row))
		wants[i] = string(want)
		if _, err := fmt.Fprintf(fixture, "%x\t%x\t%s\n", table, c.text, oracle); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-byte-replacer-replace", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s byteReplacer.Replace(%x): Lean=%s want=%s", label, cases[i].text, got, wants[i])
		}
	}
	fmt.Printf("strings.byteReplacer.Replace %s: %d Lean results agree.\n", label, len(cases))
	return nil
}
