package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var clear, full, sparse, high, alternating [256]bool
	for i := range full {
		full[i] = true
		alternating[i] = i%2 == 0
	}
	for _, i := range []int{0, 63, 64, 127, 128, 190, 192, 255} {
		sparse[i] = true
	}
	for _, i := range []int{128, 129, 255} {
		high[i] = true
	}
	var input, data bytes.Buffer
	want := []string{strings.Repeat("p", 256)}
	fmt.Fprintln(&input, 0)
	fmt.Fprintln(&data, want[0])
	var maskFailures, nextFailures int
	for _, table := range [][256]bool{clear, full, sparse, high, alternating} {
		var words [4]uint64
		row := make([]byte, len(table))
		for key, member := range table {
			row[key] = '0'
			if member {
				row[key] = '1'
				words[key/64] |= uint64(1) << (key % 64)
			}
			if table[key&127] != member {
				maskFailures++
			}
			if table[(key+1)%256] != member {
				nextFailures++
			}
		}
		want = append(want, string(row))
		fmt.Fprintln(&input, 1, words[0], words[1], words[2], words[3])
		fmt.Fprintln(&data, string(row))
	}
	for i, row := range check.LeanRows("godiff-bytes-ascii-set-contains", len(want), &input) {
		if row != want[i] {
			return fmt.Errorf("table %d: Lean=%s want=%s", i, row, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-ascii-set-contains-review-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, data.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_ASCII_SET_CONTAINS_CASES", path); err != nil {
		return err
	}
	const fixtures = "difftest/bytes/ascii-set-contains/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanASCIISetContains", fixtures+"bridge_test.go", fixtures+"contains_test.go"); err != nil {
		return err
	}
	nilFailures := len(want[0])
	if maskFailures != 5 || nextFailures != 270 || nilFailures != 256 {
		return fmt.Errorf("masking failed %d, advancing index failed %d, nil-as-false failed %d", maskFailures, nextFailures, nilFailures)
	}
	fmt.Printf("asciiSet.contains passed %d lookups; masking failed %d, advancing index failed %d, nil-as-false failed %d\n",
		256*len(want), maskFailures, nextFailures, nilFailures)
	return nil
}
