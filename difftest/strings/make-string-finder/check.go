package main

import (
	"bufio"
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
	dir, err := os.MkdirTemp("", "golean-strings-make-string-finder-")
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
	if err := os.Setenv("GOLEAN_STRINGS_MAKE_STRING_FINDER_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)

	edges := []string{"mississi", "abcxxxabc", "ababaaba", "🚀é🚀", "\xff\x00\xff", "\xe2\x82\xff",
		strings.Repeat("a", 128), strings.Repeat("ab", 64), strings.Repeat("\xff\x00", 64)}
	for width := 0; width <= 10; width++ {
		for bits := range 1 << width {
			pattern := make([]byte, width)
			for i := range pattern {
				pattern[i] = byte(bits>>i&1) * 255
			}
			edges = append(edges, string(pattern))
		}
	}
	for b := range 256 {
		edges = append(edges, string([]byte{byte(b)}), string([]byte{byte(b), 0xff, byte(b), 0}))
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for batch := range 4 {
			cases := make([]string, 5000)
			for i := range cases {
				cases[i] = sample(rng, i)
			}
			if err := compare(fixture, fmt.Sprintf("seed=%d batch=%d", seed, batch), cases); err != nil {
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
	const fixtures = "difftest/strings/make-string-finder/testdata/"
	return check.TestInPackage("strings", "TestGoLeanMakeStringFinder",
		fixtures+"bridge_test.go", fixtures+"make_string_finder_test.go")
}

func sample(rng *rand.Rand, index int) string {
	pattern := make([]byte, rng.Intn(57))
	alphabet := [...]int{2, 5, 31, 256}
	for i := range pattern {
		pattern[i] = byte(rng.Intn(alphabet[index%len(alphabet)]))
	}
	if index%8 == 0 {
		copy(pattern[len(pattern)/2:], pattern[:len(pattern)/2])
	}
	return string(pattern)
}

// For each matched suffix, enumerate previous alignments from nearest to
// farthest. If none has a different leading byte, align the largest border
// that can fit into the known suffix.
func expected(pattern string) ([256]int, []int) {
	var bad [256]int
	n := len(pattern)
	for b := range bad {
		bad[b] = n
		if n > 0 {
			if i := strings.LastIndexByte(pattern[:n-1], byte(b)); i >= 0 {
				bad[b] = n - 1 - i
			}
		}
	}
	good := make([]int, n)
	for i := range good {
		suffix := pattern[i+1:]
		matched := len(suffix)
		shift := 0
		for d := 1; d <= i; d++ {
			if pattern[i-d] != pattern[i] && pattern[i+1-d:n-d] == suffix {
				shift = d
				break
			}
		}
		if shift == 0 {
			for border := matched; border >= 0; border-- {
				if pattern[:border] == pattern[n-border:] {
					shift = n - border
					break
				}
			}
		}
		good[i] = matched + shift
	}
	return bad, good
}

func appendNumbers(row []byte, values []int) []byte {
	for i, value := range values {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(value), 10)
	}
	return row
}

func compare(fixture *bufio.Writer, label string, cases []string) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, pattern := range cases {
		row := strconv.AppendInt(nil, int64(len(pattern)), 10)
		patternValues := make([]int, len(pattern))
		for j := range len(pattern) {
			patternValues[j] = int(pattern[j])
			row = strconv.AppendUint(append(row, ' '), uint64(pattern[j]), 10)
		}
		fmt.Fprintln(&modelInput, string(row))
		bad, good := expected(pattern)
		table := appendNumbers(nil, patternValues)
		table = appendNumbers(append(table, ':'), bad[:])
		table = appendNumbers(append(table, ':'), good)
		wants[i] = string(table)
		if _, err := fmt.Fprintf(fixture, "%x\t%s\n", pattern, table); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-make-string-finder", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s makeStringFinder(%x): Lean=%s want=%s", label, cases[i], got, wants[i])
		}
	}
	fmt.Printf("strings.makeStringFinder %s: %d Lean tables agree.\n", label, len(cases))
	return nil
}
