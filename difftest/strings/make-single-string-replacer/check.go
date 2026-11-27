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

type input struct{ pattern, value string }

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := os.MkdirTemp("", "golean-strings-make-single-string-replacer-")
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
	if err := os.Setenv("GOLEAN_STRINGS_MAKE_SINGLE_STRING_REPLACER_CASES", path); err != nil {
		return err
	}
	fixture := bufio.NewWriter(file)

	patterns := []string{"", "a", "mississi", "abcxxxabc", "ababaaba", "🚀é🚀", "\xff\x00\xff",
		"\xe2\x82\xff", strings.Repeat("a", 128), strings.Repeat("ab", 64), strings.Repeat("\xff\x00", 64)}
	values := []string{"", "replacement", "�", "🚀", "\x00", "\xff\x00", "\xe2\x82", strings.Repeat("替", 64)}
	var edges []input
	for _, pattern := range patterns {
		for _, value := range values {
			edges = append(edges, input{pattern, value})
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		edges = append(edges, input{raw, "\xff" + raw}, input{"\xff" + raw + "\xff", raw})
	}
	for width := 0; width <= 9; width++ {
		for bits := range 1 << width {
			pattern := make([]byte, width)
			for i := range pattern {
				pattern[i] = byte(bits>>i&1) * 255
			}
			edges = append(edges, input{string(pattern), "\x00\xff"})
		}
	}
	if err := compare(fixture, "boundaries", edges); err != nil {
		return err
	}
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for batch := range 4 {
			cases := make([]input, 5000)
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
	const fixtures = "difftest/strings/make-single-string-replacer/testdata/"
	return check.TestInPackage("strings", "TestGoLeanMakeSingleStringReplacer",
		fixtures+"bridge_test.go", fixtures+"make_single_string_replacer_test.go")
}

func sample(rng *rand.Rand, index int) input {
	pattern := make([]byte, rng.Intn(57))
	alphabet := [...]int{2, 5, 31, 256}
	for i := range pattern {
		pattern[i] = byte(rng.Intn(alphabet[index%len(alphabet)]))
	}
	if index%8 == 0 {
		copy(pattern[len(pattern)/2:], pattern[:len(pattern)/2])
	}
	value := make([]byte, rng.Intn(57))
	rng.Read(value)
	return input{string(pattern), string(value)}
}

// Find the first candidate frame shift that can still match the known suffix.
func expected(pattern string) ([256]int, []int) {
	var bad [256]int
	n := len(pattern)
	for b := range bad {
		bad[b] = n
		for distance := 1; distance < n; distance++ {
			if pattern[n-1-distance] == byte(b) {
				bad[b] = distance
				break
			}
		}
	}
	good := make([]int, n)
	for i := range good {
		for shift := 1; shift <= n; shift++ {
			if shift <= i {
				if pattern[i-shift] == pattern[i] || pattern[i+1-shift:n-shift] != pattern[i+1:] {
					continue
				}
			} else if pattern[:n-shift] != pattern[shift:] {
				continue
			}
			good[i] = n - 1 - i + shift
			break
		}
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

func compare(fixture *bufio.Writer, label string, cases []input) error {
	var modelInput bytes.Buffer
	wants := make([]string, len(cases))
	for i, c := range cases {
		row := strconv.AppendInt(nil, int64(len(c.pattern)), 10)
		for _, raw := range []string{c.pattern, c.value} {
			for j := range len(raw) {
				row = strconv.AppendUint(append(row, ' '), uint64(raw[j]), 10)
			}
		}
		fmt.Fprintln(&modelInput, string(row))
		bad, good := expected(c.pattern)
		table := appendBytes(nil, c.pattern)
		table = appendBytes(append(table, ':'), c.value)
		table = appendNumbers(append(table, ':'), bad[:])
		table = appendNumbers(append(table, ':'), good)
		wants[i] = string(table)
		if _, err := fmt.Fprintf(fixture, "%x\t%x\t%s\n", c.pattern, c.value, table); err != nil {
			return err
		}
	}
	for i, got := range check.LeanRows("godiff-strings-make-single-string-replacer", len(cases), &modelInput) {
		if got != wants[i] {
			return fmt.Errorf("%s makeSingleStringReplacer(%x,%x): Lean=%s want=%s",
				label, cases[i].pattern, cases[i].value, got, wants[i])
		}
	}
	fmt.Printf("strings.makeSingleStringReplacer %s: %d Lean snapshots agree.\n", label, len(cases))
	return nil
}

func appendBytes(encoded []byte, raw string) []byte {
	values := make([]int, len(raw))
	for i := range len(raw) {
		values[i] = int(raw[i])
	}
	return appendNumbers(encoded, values)
}
