package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	edges := []string{"", "é", "before🚀after", "\x00\x7f\x80\x01"}
	for first := range 256 {
		for second := range 256 {
			edges = append(edges, string([]byte{byte(first), byte(second)}))
		}
	}
	ascii := make([]byte, 128)
	for i := range ascii {
		ascii[i] = byte(i)
	}
	for stop := range len(ascii) + 1 {
		edges = append(edges, string(ascii[:stop]), string(ascii[:stop])+"\xff"+string(ascii[stop:]))
	}

	seeds := []int64{17, 2718, 65537}
	cases := append([]string(nil), edges...)
	for _, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		for i := range 20_000 {
			cases = append(cases, sample(rng, i))
		}
	}
	var leanRows, fixture bytes.Buffer
	expected := make([]string, len(cases))
	for i, input := range cases {
		row := []byte{'0'}
		for j := range len(input) {
			row = strconv.AppendUint(append(row, ' '), uint64(input[j]), 10)
		}
		fmt.Fprintln(&leanRows, string(row))
		expected[i] = specification(input)
		fmt.Fprintf(&fixture, "x%x %s\n", input, expected[i])
	}
	for i, got := range check.LeanRows("godiff-strings-make-ascii-set", len(cases), &leanRows) {
		if got != expected[i] {
			return fmt.Errorf("input %q: Lean=%s want=%s", cases[i], got, expected[i])
		}
	}
	dir, err := os.MkdirTemp("", "golean-make-ascii-set-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "cases")
	if err := os.WriteFile(path, fixture.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.Setenv("GOLEAN_STRINGS_MAKE_ASCII_SET_CASES", path); err != nil {
		return err
	}
	const data = "difftest/strings/make-ascii-set/testdata/"
	if err := check.TestInPackage("strings", "TestGoLeanMakeASCIISet", data+"bridge_test.go", data+"make_ascii_set_test.go"); err != nil {
		return err
	}
	fmt.Printf("strings.makeASCIISet: %d boundaries and %d samples agree, seeds %v.\n", len(edges), len(cases)-len(edges), seeds)
	return nil
}

func specification(input string) string {
	cut := strings.IndexFunc(input, func(r rune) bool { return r >= 128 })
	valid := cut < 0
	if !valid {
		input = input[:cut]
	}
	result := strconv.AppendBool(nil, valid)
	result = append(result, '|')
	for i := range 256 {
		digit := byte('0')
		if strings.IndexByte(input, byte(i)) >= 0 {
			digit = '1'
		}
		result = append(result, digit)
	}
	return string(result)
}

func sample(rng *rand.Rand, i int) string {
	raw := make([]byte, rng.Intn(129))
	rng.Read(raw)
	if i%3 != 0 {
		for j := range raw {
			raw[j] &= 127
		}
		if i%3 == 2 && len(raw) != 0 {
			raw[rng.Intn(len(raw))] |= 128
		}
	}
	return string(raw)
}
