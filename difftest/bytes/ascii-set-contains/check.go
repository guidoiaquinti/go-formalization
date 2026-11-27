package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
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
	var clear, full [256]bool
	for i := range full {
		full[i] = true
	}
	boundaries := [][256]bool{clear, full}
	for b := range 256 {
		single, missing := clear, full
		single[b], missing[b] = true, false
		boundaries = append(boundaries, single, missing)
	}
	seeds := []int64{17, 2718, 65537}
	var leanInput, fixture bytes.Buffer
	want := []string{strings.Repeat("p", 256)}
	fmt.Fprintln(&leanInput, 0)
	fmt.Fprintln(&fixture, want[0])
	add := func(table [256]bool) {
		var words [4]uint64
		entries := make([]byte, len(table))
		for i, member := range table {
			entries[i] = '0'
			if member {
				entries[i] = '1'
				words[i/64] |= uint64(1) << (i % 64)
			}
		}
		want = append(want, string(entries))
		fmt.Fprintln(&leanInput, 1, words[0], words[1], words[2], words[3])
		fmt.Fprintln(&fixture, string(entries))
	}
	for _, table := range boundaries {
		add(table)
	}
	const tablesPerSeed = 80
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for range tablesPerSeed {
			var table [256]bool
			for i := range table {
				table[i] = rng.Intn(2) == 1
			}
			add(table)
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-ascii-set-contains", len(want), &leanInput) {
		if got != want[i] {
			return fmt.Errorf("table %d: Lean=%s want=%s", i, got, want[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-ascii-set-contains-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_BYTES_ASCII_SET_CONTAINS_CASES", path); err != nil {
		return err
	}
	const data = "difftest/bytes/ascii-set-contains/testdata/"
	if err := check.TestInPackage("bytes", "TestGoLeanASCIISetContains", data+"bridge_test.go", data+"contains_test.go"); err != nil {
		return err
	}
	fmt.Printf("bytes.asciiSet.contains: all 256 bytes for %d boundary tables, a nil receiver and %d sampled tables, seeds %v.\n",
		len(boundaries), len(seeds)*tablesPerSeed, seeds)
	return nil
}
